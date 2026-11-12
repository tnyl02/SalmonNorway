package cancelrequest

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// ErrStudentOnly คือผู้ใช้ไม่ใช่นักศึกษา ตอบ 403 FORBIDDEN_ROLE (FR-CXL-03, ASM-13, CR-006)
var ErrStudentOnly = errors.New("student only")

// ErrReasonRequired คือเหตุผลว่างหรือมีแต่ช่องว่าง ตอบ 422 VALIDATION_FAILED (FR-CXL-08, CR-006)
var ErrReasonRequired = errors.New("reason required")

// ErrReasonTooLong คือเหตุผลยาวเกิน config.CancelReasonMaxLen ตอบ 422 VALIDATION_FAILED (FR-CXL-08, Q-14, CR-006)
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

// ErrRequestNotFound คือไม่มีคำขอนี้ หรือเป็นคำขอของผู้อื่น ตอบ 404 REQUEST_NOT_FOUND (FR-CXL-03, AC-CXL-03)
// ทั้ง 2 กรณีตอบเหมือนกัน ไม่บอกว่ามีคำขอนี้อยู่หรือไม่ เพื่อไม่เปิดเผยข้อมูลคำขอของผู้อื่น
var ErrRequestNotFound = errors.New("request not found")

// ErrNotCancelled คือคำขอของตนที่สถานะยกเลิกไม่ได้ (นอกจาก approved) ตอบ 409 NOT_CANCELLABLE (FR-CXL-02, T-17)
var ErrNotCancelled = errors.New("not cancelled")

// statusApproved คือค่าในฐานข้อมูลของสถานะ Approved (CON-CXL-02 ตารางแมปชื่อสถานะ)
const statusApproved = "approved"

// Service คือกฎของฟีเจอร์ 003
type Service struct {
	repo      *Repository
	dbTimeout time.Duration // เวลารอฐานข้อมูลตอนยกเลิก (ASM-06) test เปลี่ยนเป็นค่าสั้นได้
	notifier  Notifier      // ส่ง Email แจ้งเจ้าหน้าที่หลังยกเลิกสำเร็จ (FR-CXL-09, IF-NOT-01)
}

// NewService สร้าง Service ใช้เวลารอฐานข้อมูลจาก config.CancelDBTimeout (ASM-06)
// และใช้ LogNotifier (Email จำลอง ASM-10) เปลี่ยนตัวส่งได้ด้วย WithNotifier
func NewService(repo *Repository) *Service {
	return &Service{repo: repo, dbTimeout: config.CancelDBTimeout, notifier: LogNotifier{}}
}

// WithNotifier เปลี่ยนตัวส่ง Email (เช่นตัวส่งจริงเมื่อได้คำตอบ Q-15 หรือตัวปลอมใน test)
func (s *Service) WithNotifier(n Notifier) *Service {
	s.notifier = n
	return s
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
	// ฐานข้อมูลไม่ตอบภายใน dbTimeout -> error -> handler ตอบ 503 และ rollback สถานะคงเดิม (FR-CXL-01, ASM-06)
	ctx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()
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
		switch {
		case !found:
			return CancelResultDTO{}, ErrRequestNotFound // ไม่มี หรือไม่ใช่ของผู้ขอคนนี้ (FR-CXL-03)
		case status == statusApproved:
			return CancelResultDTO{}, &StatusError{Err: ErrAlreadyApproved, Status: status} // แนบ status ใน 409 (CR-006)
		default:
			return CancelResultDTO{}, &StatusError{Err: ErrNotCancelled, Status: status}
		}
	}
	// อ่านข้อมูลสำหรับ Email ใน transaction เดียวกัน แต่ส่งหลัง commit เท่านั้น (FR-CXL-09, plan.md ข้อ 4 ขั้น 4)
	notice, err := s.repo.CancelledNotice(ctx, tx, requestID)
	if err != nil {
		return CancelResultDTO{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return CancelResultDTO{}, err
	}
	notice.CancelledBy, notice.Reason = p.UserID, reason
	_ = s.notify(context.WithoutCancel(ctx), notice) // ผลของการส่งใช้บันทึก audit ใน T-32 ไม่กระทบผลการยกเลิก
	return CancelResultDTO{RequestID: requestID, Status: "cancelled"}, nil
}

// notify ส่ง Email แจ้งเจ้าหน้าที่หลังยกเลิกสำเร็จ แบบรอผล (FR-CXL-09, ทีมตัดสิน 2569-11-12)
// การยกเลิก commit ไปแล้ว ส่งไม่สำเร็จจึงไม่ย้อนการยกเลิก คืน error ให้ T-32 บันทึก audit notify error
// ภาควิชาที่ไม่มี Email เจ้าหน้าที่นับเป็นส่งไม่สำเร็จ (ErrNoRecipients ทีมตัดสิน 2569-11-12) ไม่เรียก Notifier
func (s *Service) notify(ctx context.Context, n CancelNotice) error {
	err := ErrNoRecipients
	if len(n.To) > 0 {
		err = s.notifier.NotifyCancelled(ctx, n)
	}
	if err != nil {
		log.Printf("cancelrequest: ส่ง Email แจ้งการยกเลิกคำขอ %d ไม่สำเร็จ: %v", n.RequestID, err)
	}
	return err
}
