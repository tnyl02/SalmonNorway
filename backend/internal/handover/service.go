package handover

import (
	"context"
	"errors"
	"strings"
	"time"

	"salmonnorway/backend/internal/config"
)

// ค่า identity_method ในตาราง handover (FR-HND-02)
const (
	identityAuthAPI  = "auth_api"  // ยืนยันผ่าน Auth API สำเร็จ
	identityCardOnly = "card_only" // Auth API ไม่ตอบ ใช้ผลตรวจบัตรแทน (T-13)
)

var (
	// ErrCardNotChecked หมายถึงเจ้าหน้าที่ไม่ได้ยืนยันว่าตรวจบัตรแล้ว (FR-HND-03, ASM-HND-01)
	ErrCardNotChecked = errors.New("card not checked")
)

// MissingFieldsError หมายถึงข้อมูลที่ต้องกรอกไม่ครบ Fields คือชื่อช่องที่ขาด (CON-TECH-03 VALIDATION_FAILED)
type MissingFieldsError struct {
	Fields []string
}

func (e *MissingFieldsError) Error() string { return "missing fields: " + strings.Join(e.Fields, ", ") }

// HandoverInput คือข้อมูลที่เจ้าหน้าที่ส่งมาตอนบันทึกการส่งมอบ (plan.md v1.4 ข้อ 4)
type HandoverInput struct {
	RequestID          int
	StaffID            string // เจ้าหน้าที่ผู้บันทึก (CON-HND-01)
	BorrowerCredential string // ผู้ยืมล็อกอินบนเครื่องเจ้าหน้าที่ (ASM-HND-09) หลังบ้านยืนยันซ้ำเอง
	BorrowerFirstName  string // เจ้าหน้าที่พิมพ์ตามบัตร (ASM-HND-10)
	BorrowerLastName   string // เจ้าหน้าที่พิมพ์ตามบัตร (ASM-HND-10)
	CardChecked        bool   // FR-HND-03 ต้องเป็น true
}

// HandoverOutput คือผลของการส่งมอบที่สำเร็จ (plan.md ข้อ 4 out 201)
type HandoverOutput struct {
	LoanID       int       `json:"loanId"`
	AssetCode    string    `json:"assetCode"`
	HandedOverAt time.Time `json:"handedOverAt"`
	DueAt        time.Time `json:"dueAt"`
}

// Service ตรวจกฎของการส่งมอบก่อนส่งให้ repository บันทึก (FR-HND-02 ถึง 07)
type Service struct {
	repo     *Repository
	verifier BorrowerVerifier
}

// NewService สร้าง Service
func NewService(repo *Repository, verifier BorrowerVerifier) *Service {
	return &Service{repo: repo, verifier: verifier}
}

// Handover ตรวจข้อมูล ยืนยันผู้ยืม แล้วบันทึกการส่งมอบ (FR-HND-02, 03, 05, 07)
func (s *Service) Handover(ctx context.Context, in HandoverInput) (HandoverOutput, error) {
	var missing []string
	if strings.TrimSpace(in.BorrowerCredential) == "" {
		missing = append(missing, "borrowerCredential")
	}
	if strings.TrimSpace(in.BorrowerFirstName) == "" {
		missing = append(missing, "borrowerFirstName")
	}
	if strings.TrimSpace(in.BorrowerLastName) == "" {
		missing = append(missing, "borrowerLastName")
	}
	if len(missing) > 0 {
		return HandoverOutput{}, &MissingFieldsError{Fields: missing}
	}
	// FR-HND-03: ตรวจบัตรทุกครั้ง แม้ Auth API ยืนยันสำเร็จ (ASM-HND-01)
	if !in.CardChecked {
		return HandoverOutput{}, ErrCardNotChecked
	}

	// FR-HND-02: หลังบ้านยืนยันผู้ยืมเอง ไม่เชื่อผลที่หน้าจอส่งมา (plan.md v1.4)
	borrowerID, err := s.verifier.VerifyBorrower(ctx, in.BorrowerCredential)
	if err != nil {
		return HandoverOutput{}, err
	}

	now := time.Now()
	rec := HandoverRecord{
		RequestID:         in.RequestID,
		BorrowerUserID:    borrowerID,
		BorrowerStudentID: borrowerID, // ASM-HND-09 รหัสนักศึกษา = รหัสผู้ใช้จากการล็อกอิน
		BorrowerFirstName: strings.TrimSpace(in.BorrowerFirstName),
		BorrowerLastName:  strings.TrimSpace(in.BorrowerLastName),
		IdentityMethod:    identityAuthAPI,
		CardChecked:       in.CardChecked,
		HandedOverBy:      in.StaffID,
		HandedOverAt:      now,
		DueAt:             normalDueAt(now),
	}
	res, err := s.repo.CreateHandover(ctx, rec, config.HandoverReadyStatuses())
	if err != nil {
		return HandoverOutput{}, err
	}
	return HandoverOutput{LoanID: res.LoanID, AssetCode: res.AssetCode, HandedOverAt: rec.HandedOverAt, DueAt: rec.DueAt}, nil
}

// normalDueAt คือกำหนดคืนของการยืมแบบปกติ = LoanNormalReturnHour:00 ของวันที่ส่งมอบตามเขตเวลา TimeZone (FR-HND-05, DOM-LOAN-01)
// กรณีส่งมอบหลัง 16:00 ยังรอ Q-HND-03 (T-25)
func normalDueAt(handedOverAt time.Time) time.Time {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		loc = time.UTC
	}
	t := handedOverAt.In(loc)
	return time.Date(t.Year(), t.Month(), t.Day(), config.LoanNormalReturnHour, 0, 0, 0, loc)
}
