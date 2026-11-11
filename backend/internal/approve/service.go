package approve

import (
	"context"
	"errors"
	"strings"
	"time"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// ErrNoPrincipal เกิดเมื่อไม่มีผู้ใช้ที่ยืนยันตัวตนแล้วใน context (IF-IDP-01) ปกติ auth.Middleware กันไว้ก่อนแล้ว
var ErrNoPrincipal = errors.New("ไม่มีผู้ใช้ที่ยืนยันตัวตนแล้วใน context")

// ErrForbiddenRole เกิดเมื่อผู้ใช้ไม่ใช่เจ้าหน้าที่ ผู้ใช้ UC-12 คือเจ้าหน้าที่เท่านั้น (ASM-01, AC-APV-05, CON-TECH-03 FORBIDDEN_ROLE)
var ErrForbiddenRole = errors.New("บทบาทนี้ไม่มีสิทธิ์พิจารณาคำขอ")

// Service รวมกฎของฟีเจอร์ 004 ตรวจทุกอย่างที่ฝั่ง server (plan.md ข้อ 1)
type Service struct {
	repo *Repository
	now  func() time.Time // นาฬิกาที่ inject ได้ ใช้ทดสอบเวลา 16:00, 00:00 (plan.md ข้อ 2, ข้อ 6)
}

// NewService สร้าง Service ที่ใช้นาฬิกาจริง
func NewService(repo *Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// staff คืนเจ้าหน้าที่ที่เรียกอยู่ ถ้าไม่ใช่เจ้าหน้าที่ปฏิเสธ (ASM-01, NFR-SEC-01, AC-APV-05)
func staff(ctx context.Context) (auth.Principal, error) {
	p, ok := auth.FromContext(ctx)
	if !ok {
		return auth.Principal{}, ErrNoPrincipal
	}
	if p.Role != auth.RoleStaff {
		return auth.Principal{}, ErrForbiddenRole
	}
	return p, nil
}

// ListPending คืนคำขอยืมปกติที่รอพิจารณาของภาควิชาเดียวกับเจ้าหน้าที่ และเวลาที่โหลด (FR-APV-01 ก, ASM-06)
// loadedAt ให้หน้าจอส่งกลับมาตอนพิจารณา เพื่อแยก "พิจารณาแล้ว" กับ "กดพร้อมกัน" (plan.md ข้อ 4, T-23)
func (s *Service) ListPending(ctx context.Context) (items []PendingRequest, loadedAt time.Time, err error) {
	p, err := staff(ctx)
	if err != nil {
		return nil, time.Time{}, err
	}
	loadedAt = s.now()
	items, err = s.repo.ListPending(ctx, p.DepartmentID)
	if err != nil {
		return nil, time.Time{}, err
	}
	return items, loadedAt, nil
}

// ErrRequestNotFound เกิดเมื่อไม่มีคำขอรหัสนี้ หรือไม่ใช่คำขอยืมแบบปกติที่ 004 พิจารณา (FR-APV-01, Q-APV-01)
var ErrRequestNotFound = errors.New("ไม่พบคำขอยืมแบบปกติรหัสนี้")

// ErrNotRequested เกิดเมื่อคำขอไม่ได้อยู่ในสถานะ requested แล้ว อนุมัติหรือปฏิเสธได้เฉพาะ requested (FR-APV-01 ค ง)
// T-23 แยกกรณีนี้เป็นถูกยกเลิก / พิจารณาแล้ว / กดพร้อมกัน ตาม FR-APV-01 (ซ)
var ErrNotRequested = errors.New("คำขอนี้ไม่ได้อยู่ในสถานะรออนุมัติ")

// Approve อนุมัติคำขอยืมแบบปกติ 1 ใบ ใน transaction เดียว (FR-APV-01 ค ฌ, AC-APV-01)
//   - ล็อกคำขอทุกใบของอุปกรณ์รายชิ้นนั้นก่อน (FR-APV-01 ช, plan.md ข้อ 7.3)
//   - คำขอเป็น approved บันทึก approval_decision และ audit_log action approve (RQ-25)
//   - ไม่เปลี่ยนสถานะอุปกรณ์ ยังเป็น "ว่าง" จนกว่า 005 ส่งมอบ (CON-DATA-01, data-model.md แถว UC-12)
//
// การตรวจคำขอแรก (T-17), อุปกรณ์ว่าง (T-19), เวลาตัด 16:00 (T-21), ภาควิชา (T-25),
// ปฏิเสธคำขออื่นตาม (T-16) และ Email (T-15) เพิ่มใน task ของแต่ละเรื่อง
func (s *Service) Approve(ctx context.Context, requestID int) (Result, error) {
	p, err := staff(ctx)
	if err != nil {
		return Result{}, err
	}
	var out Result
	err = s.repo.inTx(ctx, func(tx *Repository) error {
		target, _, ok, err := tx.lockItemRequests(ctx, requestID)
		if err != nil {
			return err
		}
		if !ok || target.LoanType != config.LoanTypeNormal {
			return ErrRequestNotFound
		}
		if target.Status != config.StatusRequested {
			return ErrNotRequested
		}
		if err := tx.setStatus(ctx, target.ID, config.StatusApproved); err != nil {
			return err
		}
		at, err := tx.insertDecision(ctx, decision{
			RequestID:         target.ID,
			Result:            config.ApprovalKindApproved,
			ActorID:           p.UserID,
			ActorDepartmentID: p.DepartmentID,
		})
		if err != nil {
			return err
		}
		if err := tx.insertAudit(ctx, p.UserID, string(p.Role), config.AuditActionApprove, target.AssetCode, target.ID); err != nil {
			return err
		}
		out = Result{RequestID: target.ID, Status: config.StatusApproved, DecidedAt: at, RejectedBy: []int{}}
		return nil
	})
	return out, err
}

// ErrReasonRequired เกิดเมื่อปฏิเสธโดยไม่ระบุเหตุผล (FR-APV-01 ง, AC-APV-06, CON-TECH-03 VALIDATION_FAILED ช่อง reason)
var ErrReasonRequired = errors.New("ต้องระบุเหตุผลในการปฏิเสธ")

// Reject ปฏิเสธคำขอยืมแบบปกติ 1 ใบพร้อมเหตุผลที่เจ้าหน้าที่ระบุ ใน transaction เดียว (FR-APV-01 ง ฌ, AC-APV-02)
//   - คำขอเป็น rejected บันทึกเหตุผลใน approval_decision และ audit_log action reject (RQ-25)
//   - ไม่เปลี่ยนสถานะอุปกรณ์ (CON-DATA-01)
//
// การตรวจภาควิชา (T-25), Email (T-15) และข้อความกรณีพิจารณาแล้ว (T-23) เพิ่มใน task ของแต่ละเรื่อง
func (s *Service) Reject(ctx context.Context, requestID int, reason string) (Result, error) {
	p, err := staff(ctx)
	if err != nil {
		return Result{}, err
	}
	// FR-APV-01 (ง)(ซ), AC-APV-06: ต้องระบุเหตุผล ว่างหรือมีแต่ช่องว่างไม่นับ ไม่เปิด transaction เลย สถานะจึงคงเดิม
	if strings.TrimSpace(reason) == "" {
		return Result{}, ErrReasonRequired
	}
	var out Result
	err = s.repo.inTx(ctx, func(tx *Repository) error {
		target, _, ok, err := tx.lockItemRequests(ctx, requestID)
		if err != nil {
			return err
		}
		if !ok || target.LoanType != config.LoanTypeNormal {
			return ErrRequestNotFound
		}
		if target.Status != config.StatusRequested {
			return ErrNotRequested
		}
		if err := tx.setStatus(ctx, target.ID, config.StatusRejected); err != nil {
			return err
		}
		at, err := tx.insertDecision(ctx, decision{
			RequestID:         target.ID,
			Result:            config.ApprovalKindRejected,
			Reason:            reason,
			ActorID:           p.UserID,
			ActorDepartmentID: p.DepartmentID,
		})
		if err != nil {
			return err
		}
		if err := tx.insertAudit(ctx, p.UserID, string(p.Role), config.AuditActionReject, target.AssetCode, target.ID); err != nil {
			return err
		}
		out = Result{RequestID: target.ID, Status: config.StatusRejected, DecidedAt: at, RejectedBy: []int{}}
		return nil
	})
	return out, err
}
