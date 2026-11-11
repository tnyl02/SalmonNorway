package cancelrequest

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// ErrStudentOnly คือผู้ใช้ไม่ใช่นักศึกษา ตอบ 403 STUDENT_ONLY (FR-CXL-03, ASM-13)
var ErrStudentOnly = errors.New("student only")

// ErrReasonRequired คือเหตุผลว่างหรือมีแต่ช่องว่าง ตอบ 400 REASON_REQUIRED (FR-CXL-08)
var ErrReasonRequired = errors.New("reason required")

// ErrReasonTooLong คือเหตุผลยาวเกิน config.CancelReasonMaxLen ตอบ 400 REASON_TOO_LONG (FR-CXL-08, Q-14)
var ErrReasonTooLong = errors.New("reason too long")

// ValidateReason ตรวจเหตุผลการยกเลิก (FR-CXL-08, ASM-13)
// ว่างหรือมีแต่ช่องว่าง = ErrReasonRequired, ยาวเกินนับตาม code point ของข้อความที่ส่งมา = ErrReasonTooLong
// นับแบบเดียวกับตัวนับ "x/200" บนหน้าจอ (countReason ใน CancelDialog.jsx นับข้อความเดิมรวมช่องว่าง)
func ValidateReason(reason string) error {
	if strings.TrimSpace(reason) == "" {
		return ErrReasonRequired
	}
	if utf8.RuneCountInString(reason) > config.CancelReasonMaxLen {
		return ErrReasonTooLong
	}
	return nil
}

// ErrAlreadyApproved คือคำขอได้รับการอนุมัติแล้ว ตอบ 409 ALREADY_APPROVED (FR-CXL-02, FR-CXL-04, AC-CXL-02, AC-CXL-04)
var ErrAlreadyApproved = errors.New("already approved")

// ErrNotCancelled คือ UPDATE ไม่เจอแถวด้วยสาเหตุอื่นนอกจากอนุมัติแล้ว ตอนนี้ตอบ 409 NOT_CANCELLABLE
// การแยก 404 REQUEST_NOT_FOUND (ไม่มีหรือของผู้อื่น) อยู่ใน T-19 สถานะอื่นที่ยกเลิกไม่ได้อยู่ใน T-17
var ErrNotCancelled = errors.New("not cancelled")

// statusApproved คือค่าในฐานข้อมูลของสถานะ Approved (CON-CXL-02 ตารางแมปชื่อสถานะ)
const statusApproved = "approved"

// Service คือกฎของฟีเจอร์ 003
type Service struct {
	repo *Repository
}

// NewService สร้าง Service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ListMyPending คืนคำขอที่ยังยกเลิกได้ของนักศึกษาที่เข้าสู่ระบบ (FR-CXL-06, FR-CXL-03, CON-CXL-02)
// ผู้ขอมาจาก Principal เสมอ ไม่รับรหัสผู้ใช้จาก request (IF-IDP-01) เจ้าหน้าที่ได้ ErrStudentOnly
func (s *Service) ListMyPending(ctx context.Context, p auth.Principal) ([]PendingRequestDTO, error) {
	if p.Role != auth.RoleStudent {
		return nil, ErrStudentOnly
	}
	rows, err := s.repo.ListPending(ctx, p.UserID, config.CancellableStatuses())
	if err != nil {
		return nil, err
	}
	out := make([]PendingRequestDTO, 0, len(rows))
	for _, r := range rows {
		out = append(out, PendingRequestDTO{
			RequestID:   r.RequestID,
			ModelName:   r.ModelName,
			AssetCode:   r.AssetCode,
			SubmittedAt: r.SubmittedAt.UTC(),
		})
	}
	return out, nil
}

// Cancel ยกเลิกคำขอ requestID ของนักศึกษาที่เข้าสู่ระบบพร้อมเหตุผล (FR-CXL-01, FR-CXL-04, FR-CXL-08, CON-CXL-02)
// ผู้ขอมาจาก Principal เสมอ (IF-IDP-01) ไม่เปลี่ยนสถานะอุปกรณ์ (CON-CXL-01)
// ตรวจเหตุผลก่อนแตะฐานข้อมูล (FR-CXL-08) การบันทึก audit_log ใน transaction เดียวกันอยู่ใน T-31 การส่ง Email หลัง commit อยู่ใน T-27
func (s *Service) Cancel(ctx context.Context, p auth.Principal, requestID int, reason string) (CancelResultDTO, error) {
	if p.Role != auth.RoleStudent {
		return CancelResultDTO{}, ErrStudentOnly
	}
	if err := ValidateReason(reason); err != nil {
		return CancelResultDTO{}, err
	}
	tx, err := s.repo.db.Begin(ctx)
	if err != nil {
		return CancelResultDTO{}, err
	}
	defer tx.Rollback(ctx) // ไม่มีผลถ้า commit ไปแล้ว

	ok, err := s.repo.Cancel(ctx, tx, requestID, p.UserID, reason, config.CancellableStatuses())
	if err != nil {
		return CancelResultDTO{}, err
	}
	if !ok {
		// แยกสาเหตุจากสถานะ ณ ตอนนี้ ใน transaction เดียวกัน (FR-CXL-02, FR-CXL-04)
		status, found, err := s.repo.OwnStatus(ctx, tx, requestID, p.UserID)
		if err != nil {
			return CancelResultDTO{}, err
		}
		if found && status == statusApproved {
			return CancelResultDTO{}, ErrAlreadyApproved
		}
		return CancelResultDTO{}, ErrNotCancelled
	}
	if err := tx.Commit(ctx); err != nil {
		return CancelResultDTO{}, err
	}
	return CancelResultDTO{RequestID: requestID, Status: "cancelled"}, nil
}
