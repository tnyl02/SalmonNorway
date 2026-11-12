package cancelrequest

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/audit"
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

// ค่าใน audit_log ของการยกเลิก (FR-CXL-05 ตาราง, ASM-12, DM-07)
const (
	auditActionCancel    = "cancel"
	auditActionNotify    = "notify" // ส่ง Email ไม่สำเร็จ (FR-CXL-09, migration 0005)
	causeNotOwner        = "not_owner"
	causeNotFound        = "not_found"
	causeAlreadyApproved = "already_approved"
	causeNotCancellable  = "not_cancellable"
	causeReasonInvalid   = "reason_invalid"
)

// Cancel ยกเลิกคำขอ requestID ของนักศึกษาที่เข้าสู่ระบบพร้อมเหตุผล (FR-CXL-01, FR-CXL-04, FR-CXL-08, CON-CXL-02)
// ผู้ขอมาจาก Principal เสมอ (IF-IDP-01) ไม่เปลี่ยนสถานะอุปกรณ์ (CON-CXL-01)
// ทุกครั้งที่นักศึกษาพยายามยกเลิก บันทึก audit_log ใน transaction เดียวกับการยกเลิก (FR-CXL-05, T-31)
// บันทึกไม่ได้ = error -> 503 และไม่ยกเลิก เจ้าหน้าที่ไม่บันทึก เพราะ FR-CXL-05 กำหนดเฉพาะนักศึกษา
// การส่ง Email หลัง commit อยู่ใน T-27
func (s *Service) Cancel(ctx context.Context, p auth.Principal, requestID int, reason string) (CancelResultDTO, error) {
	if p.Role != auth.RoleStudent {
		return CancelResultDTO{}, ErrStudentOnly
	}
	// ฐานข้อมูลไม่ตอบภายใน dbTimeout -> error -> handler ตอบ 503 และ rollback สถานะคงเดิม (FR-CXL-01, ASM-06)
	ctx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()
	tx, err := s.repo.db.Begin(ctx)
	if err != nil {
		return CancelResultDTO{}, err
	}
	defer tx.Rollback(ctx) // ไม่มีผลถ้า commit ไปแล้ว

	// ตรวจเหตุผลก่อนแตะคำขอ (FR-CXL-08) ไม่ผ่านบันทึก reason_invalid แล้วตอบ 422
	if verr := ValidateReason(reason); verr != nil {
		if err := s.auditRejected(ctx, tx, p, requestID, causeReasonInvalid); err != nil {
			return CancelResultDTO{}, err
		}
		return CancelResultDTO{}, verr
	}

	ok, err := s.repo.Cancel(ctx, tx, requestID, p.UserID, reason, config.CancellableStatuses())
	if err != nil {
		return CancelResultDTO{}, err
	}
	if !ok {
		// แยกสาเหตุจากสถานะ ณ ตอนนี้ ใน transaction เดียวกัน (FR-CXL-02, FR-CXL-04)
		owner, status, found, err := s.repo.RequestState(ctx, tx, requestID)
		if err != nil {
			return CancelResultDTO{}, err
		}
		var cause string
		var result error
		switch {
		case !found:
			cause, result = causeNotFound, ErrRequestNotFound
		case owner != p.UserID:
			cause, result = causeNotOwner, ErrRequestNotFound // ตอบเหมือนไม่มีคำขอ ไม่บอกว่ามีอยู่ (FR-CXL-03)
		case status == statusApproved:
			cause, result = causeAlreadyApproved, &StatusError{Err: ErrAlreadyApproved, Status: status} // แนบ status ใน 409 (CR-006)
		default:
			cause, result = causeNotCancellable, &StatusError{Err: ErrNotCancelled, Status: status}
		}
		if err := s.auditRejected(ctx, tx, p, requestID, cause); err != nil {
			return CancelResultDTO{}, err
		}
		return CancelResultDTO{}, result
	}
	// ยกเลิกสำเร็จ: บันทึก ok พร้อมความยาวเหตุผล นับตาม code point แบบ FR-CXL-08
	if _, err := audit.Write(ctx, tx, audit.Entry{
		ActorID: p.UserID, ActorRole: string(p.Role), Action: auditActionCancel,
		Query: map[string]any{"reason_len": utf8.RuneCountInString(reason)}, RequestID: &requestID, Outcome: audit.OutcomeOK,
	}); err != nil {
		return CancelResultDTO{}, err
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
	// ส่ง Email หลัง commit ส่งไม่สำเร็จไม่ย้อนการยกเลิก แต่บันทึก notify error (FR-CXL-09, FR-CXL-05, T-32)
	if err := s.notify(context.WithoutCancel(ctx), notice); err != nil {
		s.auditNotifyFailed(context.WithoutCancel(ctx), p, requestID)
	}
	return CancelResultDTO{RequestID: requestID, Status: "cancelled"}, nil
}

// auditNotifyFailed บันทึกการส่ง Email ไม่สำเร็จ ตามตาราง FR-CXL-05: action notify, outcome error, query {"for":"cancel"}
// รวมกรณีภาควิชาไม่มี Email เจ้าหน้าที่ (ErrNoRecipients ทีมตัดสิน 2569-11-12)
// การยกเลิก commit ไปแล้ว เขียนนอก transaction และรอไม่เกิน dbTimeout (ASM-06)
// บันทึกไม่สำเร็จแค่ log ไว้ ไม่เปลี่ยนผลที่ตอบนักศึกษา เพราะยกเลิกสำเร็จแล้ว (FR-CXL-09 ทีมเลือกเอง)
func (s *Service) auditNotifyFailed(ctx context.Context, p auth.Principal, requestID int) {
	ctx, cancel := context.WithTimeout(ctx, s.dbTimeout)
	defer cancel()
	if _, err := audit.Write(ctx, s.repo.db, audit.Entry{
		ActorID: p.UserID, ActorRole: string(p.Role), Action: auditActionNotify,
		Query: map[string]any{"for": auditActionCancel}, RequestID: &requestID, Outcome: audit.OutcomeError,
	}); err != nil {
		log.Printf("cancelrequest: บันทึก audit notify error ของคำขอ %d ไม่สำเร็จ: %v", requestID, err)
	}
}

// auditRejected บันทึกการยกเลิกที่ไม่สำเร็จ (outcome forbidden) แล้ว commit (FR-CXL-05 ตาราง)
// ใน transaction นี้ไม่มีการเปลี่ยนสถานะคำขอ commit จึงเก็บแค่แถว audit
// คำขอที่มีอยู่จริงใส่ request_id ไม่มีคำขอนี้ใส่ request_id ว่างและเก็บรหัสที่ส่งมาใน query (FK ของ 0003 อ้างคำขอที่ไม่มีไม่ได้)
// เหตุผลผิดกับคำขอที่ไม่มีอยู่ใช้รูปแบบเดียวกับ not_found (ทีมเลือกเอง ตาราง FR-CXL-05 ไม่ได้แยกกรณีนี้)
func (s *Service) auditRejected(ctx context.Context, tx pgx.Tx, p auth.Principal, requestID int, cause string) error {
	exists := cause != causeNotFound
	if cause == causeReasonInvalid {
		_, _, found, err := s.repo.RequestState(ctx, tx, requestID)
		if err != nil {
			return err
		}
		exists = found
	}
	query := map[string]any{"cause": cause}
	var reqID *int
	if exists {
		reqID = &requestID
	} else {
		query["id"] = requestID
	}
	if _, err := audit.Write(ctx, tx, audit.Entry{
		ActorID: p.UserID, ActorRole: string(p.Role), Action: auditActionCancel,
		Query: query, RequestID: reqID, Outcome: audit.OutcomeForbidden,
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// notify ส่ง Email แจ้งเจ้าหน้าที่หลังยกเลิกสำเร็จ แบบรอผล (FR-CXL-09, ทีมตัดสิน 2569-11-12)
// การยกเลิก commit ไปแล้ว ส่งไม่สำเร็จจึงไม่ย้อนการยกเลิก คืน error ให้ Cancel บันทึก audit notify error (T-32)
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
