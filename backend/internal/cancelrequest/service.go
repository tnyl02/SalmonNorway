package cancelrequest

import (
	"context"
	"errors"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// ErrStudentOnly คือผู้ใช้ไม่ใช่นักศึกษา ตอบ 403 STUDENT_ONLY (FR-CXL-03, ASM-13)
var ErrStudentOnly = errors.New("student only")

// ErrNotCancelled คือไม่มีแถวที่ยกเลิกได้ (UPDATE ไม่เจอแถว) ตอนนี้ตอบ 409 NOT_CANCELLABLE ทุกกรณี
// การแยกสาเหตุเป็น 404 REQUEST_NOT_FOUND / 409 ALREADY_APPROVED / 409 NOT_CANCELLABLE อยู่ใน T-15, T-17, T-19
var ErrNotCancelled = errors.New("not cancelled")

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
// การตรวจเหตุผลอยู่ใน T-09 การบันทึก audit_log ใน transaction เดียวกันอยู่ใน T-31 การส่ง Email หลัง commit อยู่ใน T-27
func (s *Service) Cancel(ctx context.Context, p auth.Principal, requestID int, reason string) (CancelResultDTO, error) {
	if p.Role != auth.RoleStudent {
		return CancelResultDTO{}, ErrStudentOnly
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
		return CancelResultDTO{}, ErrNotCancelled
	}
	if err := tx.Commit(ctx); err != nil {
		return CancelResultDTO{}, err
	}
	return CancelResultDTO{RequestID: requestID, Status: "cancelled"}, nil
}
