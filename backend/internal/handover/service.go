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
	identityCardOnly = "card_only" // Auth API ไม่ตอบภายใน BorrowerAuthTimeout ใช้ผลตรวจบัตรแทน (FR-HND-02, AC-HND-03)
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
	StaffID            string       // เจ้าหน้าที่ผู้บันทึก (CON-HND-01)
	BorrowerCredential string       // ผู้ยืมล็อกอินบนเครื่องเจ้าหน้าที่ (ASM-HND-09) หลังบ้านยืนยันซ้ำเอง
	BorrowerFirstName  string       // เจ้าหน้าที่พิมพ์ตามบัตร (ASM-HND-10)
	BorrowerLastName   string       // เจ้าหน้าที่พิมพ์ตามบัตร (ASM-HND-10)
	CardChecked        bool         // FR-HND-03 ต้องเป็น true
	Photos             []PhotoInput // FR-HND-06 รูปสภาพอุปกรณ์
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
	repo       *Repository
	requesters RowQuerier // อ่าน requester_id ของคำขอ ใช้เมื่อ Auth API ไม่ตอบ (FR-HND-02, AC-HND-03)
	verifier   BorrowerVerifier
	photos     *PhotoStore      // FR-HND-06 ที่เก็บไฟล์รูป
	now        func() time.Time // นาฬิกา inject ได้ ใช้ทดสอบเวลา 10:00 และ 16:30 โดยไม่ต้องรอจริง (FR-HND-05, plan.md ข้อ 6)
}

// NewService สร้าง Service now = nil ใช้ time.Now
func NewService(repo *Repository, requesters RowQuerier, verifier BorrowerVerifier, photos *PhotoStore, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{repo: repo, requesters: requesters, verifier: verifier, photos: photos, now: now}
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
	// FR-HND-06: รูปอย่างน้อย 2 รูป JPG/PNG ถ้าไม่ครบ ไม่บันทึกการส่งมอบ
	photos, err := validatePhotos(in.Photos)
	if err != nil {
		return HandoverOutput{}, err
	}

	// FR-HND-02: หลังบ้านยืนยันผู้ยืมเอง ไม่เชื่อผลที่หน้าจอส่งมา (plan.md v1.4)
	borrowerID, method, err := s.identifyBorrower(ctx, in)
	if err != nil {
		return HandoverOutput{}, err
	}

	now := s.now() // FR-HND-05 เวลาส่งมอบจริง
	rec := HandoverRecord{
		RequestID:         in.RequestID,
		BorrowerUserID:    borrowerID,
		BorrowerStudentID: borrowerID, // ASM-HND-09 รหัสนักศึกษา = รหัสผู้ใช้จากการล็อกอิน (card_only = requester_id)
		BorrowerFirstName: strings.TrimSpace(in.BorrowerFirstName),
		BorrowerLastName:  strings.TrimSpace(in.BorrowerLastName),
		IdentityMethod:    method,
		CardChecked:       in.CardChecked,
		HandedOverBy:      in.StaffID,
		HandedOverAt:      now,
		DueAt:             normalDueAt(now),
	}
	res, err := s.repo.CreateHandover(ctx, rec, config.HandoverReadyStatuses(), photoSaver{store: s.photos, photos: photos})
	if err != nil {
		return HandoverOutput{}, err
	}
	return HandoverOutput{LoanID: res.LoanID, AssetCode: res.AssetCode, HandedOverAt: rec.HandedOverAt, DueAt: rec.DueAt}, nil
}

// identifyBorrower คืนรหัสผู้ยืมและวิธียืนยันตัวตน (FR-HND-02, FR-HND-03, IF-IDP-01)
// Auth API ตอบทันเวลา = รหัสจาก Auth API และ auth_api (ผู้ส่งคำขอหรือไม่ repository ตรวจต่อ)
// Auth API ไม่ตอบภายใน BorrowerAuthTimeout = ใช้ผลตรวจบัตร (AC-HND-03) บันทึกด้วย card_only
// และใช้ requester_id ของคำขอเป็นรหัสผู้ยืม เพราะเจ้าหน้าที่ตรวจบัตรเทียบกับผู้ส่งคำขอที่หน้าจอแสดง
// (ทีมตัดสิน 2569-11-12 ตอน /implement T-13 บันทึกใน prompt-log ยังไม่มีใน spec)
func (s *Service) identifyBorrower(ctx context.Context, in HandoverInput) (userID, method string, err error) {
	userID, err = s.verifier.VerifyBorrower(ctx, in.BorrowerCredential)
	if err == nil {
		return userID, identityAuthAPI, nil
	}
	if !errors.Is(err, ErrBorrowerAuthUnavailable) {
		return "", "", err
	}
	requesterID, found, err := GetRequester(ctx, s.requesters, in.RequestID)
	if err != nil {
		return "", "", err
	}
	if !found {
		return "", "", ErrRequestNotFound
	}
	return requesterID, identityCardOnly, nil
}

// normalDueAt คือกำหนดคืนของการยืมแบบปกติ = LoanNormalReturnHour:00 (ค่าของ 002 ใน config.go) ของวันที่ส่งมอบตามเขตเวลา TimeZone (FR-HND-05, DOM-LOAN-01)
// ส่งมอบหลัง LoanNormalReturnHour:00 กำหนดคืน = LoanNormalReturnHour:00 ของวันถัดไป (FR-HND-05, Q-HND-03 ปิด 2569-11-12, T-25)
// ส่งมอบตรง 16:00:00 พอดียังนับเป็นวันเดียวกัน เพราะ spec เขียนว่า "หลัง 16:00"
func normalDueAt(handedOverAt time.Time) time.Time {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		loc = time.UTC
	}
	t := handedOverAt.In(loc)
	due := time.Date(t.Year(), t.Month(), t.Day(), config.LoanNormalReturnHour, 0, 0, 0, loc)
	if t.After(due) {
		due = due.AddDate(0, 0, 1)
	}
	return due
}
