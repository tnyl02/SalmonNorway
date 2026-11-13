package approve

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

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

// ErrRequestNotFound เกิดเมื่อไม่มีคำขอรหัสนี้ (FR-APV-01 ซ, Q-APV-25)
// v3: การอนุมัติรับคำขอค้างคืนแล้ว (T-44) การปฏิเสธคำขอค้างคืนยังได้ error นี้จนกว่าจะทำ T-48
var ErrRequestNotFound = errors.New("ไม่พบคำขอรหัสนี้")

// คำขอไม่ได้อยู่ในสถานะ requested แล้ว อนุมัติหรือปฏิเสธไม่ได้ แยก 3 กรณีตาม FR-APV-01 (ซ) (T-23, AC-APV-09, AC-APV-10)
var (
	// ErrRequestCancelled คำขอถูกผู้ขอยกเลิกแล้ว (REQUEST_CANCELLED)
	ErrRequestCancelled = errors.New("คำขอนี้ถูกผู้ขอยกเลิกแล้ว")
	// ErrAlreadyDecided คำขอถูกพิจารณาไปแล้วก่อนหน้าจอของผู้เรียกโหลด (REQUEST_ALREADY_DECIDED)
	ErrAlreadyDecided = errors.New("คำขอนี้ได้รับการพิจารณาแล้ว")
	// ErrDecidedConcurrently คำขอเพิ่งถูกพิจารณาโดยผู้อื่นหลังหน้าจอของผู้เรียกโหลด กดพร้อมกัน (REQUEST_DECIDED_CONCURRENTLY, FR-APV-01 ช)
	ErrDecidedConcurrently = errors.New("คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น")
)

// notRequestedError บอกว่าคำขอที่ไม่ใช่ requested แล้วเป็นกรณีไหน (FR-APV-01 ซ, plan.md ข้อ 4)
// แยก "พิจารณาแล้ว" กับ "กดพร้อมกัน" ด้วยการเทียบ approval_decision.decided_at กับ loadedAt ที่หน้าจอส่งมา
// (ทีมเลือกเอง ไม่ได้มาจาก spec) loadedAt ว่าง (ไม่ได้ส่งมา) ถือเป็นพิจารณาแล้ว
func notRequestedError(ctx context.Context, tx *Repository, target requestRow, loadedAt time.Time) error {
	if target.Status == config.StatusCancelled {
		return ErrRequestCancelled
	}
	decidedAt, ok, err := tx.decisionTime(ctx, target.ID)
	if err != nil {
		return err
	}
	if ok && !loadedAt.IsZero() && decidedAt.After(loadedAt) {
		return ErrDecidedConcurrently
	}
	return ErrAlreadyDecided
}

// ErrAdvisorNameRequired เกิดเมื่ออนุมัติคำขอค้างคืนโดยไม่ระบุชื่ออาจารย์ ว่างหรือมีแต่ช่องว่าง
// (FR-APV-03 ช, Q-APV-44, VALIDATION_FAILED ช่อง advisorName)
var ErrAdvisorNameRequired = errors.New("ต้องระบุชื่ออาจารย์ที่เซ็นอนุมัติ")

// ErrAdvisorNameTooLong เกิดเมื่อชื่ออาจารย์ยาวเกิน config.ApprovalAdvisorNameMaxLength นับตาม code point
// (FR-APV-03 ช, Q-APV-44, VALIDATION_FAILED ช่อง advisorName)
var ErrAdvisorNameTooLong = errors.New("ชื่ออาจารย์ยาวเกินกำหนด")

// Approve อนุมัติคำขอโดยไม่มีชื่ออาจารย์ ใช้กับคำขอยืมปกติ (FR-APV-01 ค) คำขอค้างคืนได้ ErrAdvisorNameRequired
func (s *Service) Approve(ctx context.Context, requestID int, loadedAt time.Time) (Result, error) {
	return s.ApproveWithAdvisor(ctx, requestID, loadedAt, "")
}

// ApproveWithAdvisor อนุมัติคำขอยืม 1 ใบ ใน transaction เดียว แยกกฎตาม loan_type (plan.md v2 ข้อ 1)
//
// ยืมปกติ (FR-APV-01 ค ฌ, AC-APV-01) advisorName ไม่ใช้ ส่งมาก็ไม่บันทึก
//
// ค้างคืน (FR-APV-03 ข ค ง ช, AC-APV-22) ต้องมีชื่ออาจารย์ คำขอเป็น approved ทันที บันทึก advisor_name
// และ approval_step 3 แถว (approved_by_officer -> pending_advisor_approval -> approved) ไม่ตรวจคิวและเวลาตัด
//
// ตรวจตามลำดับ: ภาควิชา -> ชื่ออาจารย์ (เฉพาะค้างคืน ทีมตอบ 2569-11-13 ตอน T-44 แบบเดียวกับเหตุผลของการปฏิเสธ Q-APV-34)
// -> สถานะคำขอ -> คำขอแรก (เฉพาะยืมปกติ) -> อุปกรณ์ "ว่าง" (ทั้ง 2 ประเภท Q-APV-42) -> เวลาตัด (เฉพาะยืมปกติ)
//
// รายละเอียดของยืมปกติ (เดิม):
//   - ล็อกคำขอทุกใบของอุปกรณ์รายชิ้นนั้นก่อน (FR-APV-01 ช, plan.md ข้อ 7.3)
//   - คำขอเป็น approved บันทึก approval_decision และ audit_log action approve (RQ-25)
//   - ไม่เปลี่ยนสถานะอุปกรณ์ ยังเป็น "ว่าง" จนกว่า 005 ส่งมอบ (CON-DATA-01, data-model.md แถว UC-12)
//   - คำขออื่นที่ยังเป็น requested ของอุปกรณ์ชิ้นเดียวกัน ทั้งยืมปกติและค้างคืน ถูกปฏิเสธตามทันที (FR-APV-01 จ, Q-APV-14, T-16)
//
// ตรวจตามลำดับ: ภาควิชา (ASM-06, T-25) -> สถานะคำขอ (T-23) -> คำขอแรก (FR-APV-01 ข, T-17)
// -> อุปกรณ์ "ว่าง" (FR-APV-01 ค, T-19) -> เวลาตัด 16:00 น. ของวันที่เริ่มยืม (FR-APV-01 ค, T-21)
func (s *Service) ApproveWithAdvisor(ctx context.Context, requestID int, loadedAt time.Time, advisorName string) (Result, error) {
	p, err := staff(ctx)
	if err != nil {
		return Result{}, err
	}
	var out Result
	err = s.repo.inTx(ctx, func(tx *Repository) error {
		target, all, ok, err := tx.lockItemRequests(ctx, requestID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrRequestNotFound
		}
		overnight := target.LoanType == config.LoanTypeOvernight // FR-APV-03, CR-009
		// ASM-06, IF-IDP-01: พิจารณาได้เฉพาะคำขอของอุปกรณ์ภาควิชาเดียวกับเจ้าหน้าที่ (AC-APV-12, T-25)
		// ตรวจก่อนเรื่องอื่น ผู้ที่อยู่ภาควิชาอื่นจึงไม่รู้สถานะหรือข้อมูลใดของคำขอนั้น
		if target.DepartmentID != p.DepartmentID {
			return ErrForbiddenDepartment
		}
		// FR-APV-03 (ข)(ช), Q-APV-44: ค้างคืนต้องระบุชื่ออาจารย์ ว่างหรือมีแต่ช่องว่างไม่นับ ไม่เกิน 100 code point
		// ตรวจก่อนสถานะ (ทีมตอบ 2569-11-13 ตอน T-44) ยังไม่เขียนอะไร สถานะจึงคงเดิม
		if overnight {
			if strings.TrimSpace(advisorName) == "" {
				return ErrAdvisorNameRequired
			}
			if utf8.RuneCountInString(advisorName) > config.ApprovalAdvisorNameMaxLength {
				return ErrAdvisorNameTooLong
			}
		} else {
			advisorName = "" // ยืมปกติไม่มีชื่ออาจารย์ (plan.md v2 ข้อ 4)
		}
		if target.Status != config.StatusRequested {
			return notRequestedError(ctx, tx, target, loadedAt) // FR-APV-01 (ซ), T-23
		}
		// FR-APV-01 (ข), ASM-04: อนุมัติได้เฉพาะคำขอแรก (ส่งก่อนสุด) ที่ยังเป็น requested ของชิ้นนั้น ข้ามลำดับไม่ได้
		// นับเฉพาะยืมปกติ คำขอค้างคืนไม่ขวางคิว (Q-APV-23) all เรียงตาม submitted_at, id แล้ว (lockItemRequests)
		// คำขอค้างคืนไม่มีคิว (FR-APV-03 ค)
		if !overnight {
			if first, ok := firstInQueue(all); !ok || first.ID != target.ID {
				return ErrNotFirst
			}
		}
		// FR-APV-01 (ค), CON-DATA-01: อนุมัติได้เมื่ออุปกรณ์รายชิ้นมีสถานะ "ว่าง" เท่านั้น (AC-APV-08, T-19)
		// ค้างคืนก็เช่นกัน (FR-APV-03 ข, Q-APV-42)
		if target.ItemStatus != config.ItemStatusAvailable {
			return ErrItemNotAvailable
		}
		// FR-APV-01 (ค), Q-APV-11, DOM-LOAN-01: อนุมัติไม่ได้เมื่อเลย 16:00 น. ของวันที่เริ่มยืม (AC-APV-16, T-21)
		// "เลย 16:00 น." นับตั้งแต่ 16:00:00 น. ตรง อนุมัติได้ถึง 15:59:59 น. (Q-APV-26, TC-APV-01-5)
		// คำขอนั้นยังเป็นคำขอแรก เจ้าหน้าที่ปฏิเสธเองหรือรอเที่ยงคืน (Q-APV-16)
		// คำขอค้างคืนไม่มีเวลาตัด (FR-APV-03 ง)
		if !overnight && !s.now().Before(approvalCutoff(target.StartAt)) {
			return ErrApprovalCutoff
		}
		if err := tx.setStatus(ctx, target.ID, config.StatusApproved); err != nil {
			return err
		}
		at, err := tx.insertDecision(ctx, decision{
			RequestID:         target.ID,
			Result:            config.ApprovalKindApproved,
			ActorID:           p.UserID,
			ActorDepartmentID: p.DepartmentID,
			AdvisorName:       advisorName,
		})
		if err != nil {
			return err
		}
		// FR-APV-03 (ข): ประวัติสถานะระหว่างทางของคำขอค้างคืนที่ผ่านในครั้งเดียว (approval_step, plan.md v2 ข้อ 3)
		if overnight {
			for _, st := range []string{config.StatusApprovedByOfficer, config.StatusPendingAdvisorApproval, config.StatusApproved} {
				if err := tx.insertStep(ctx, target.ID, st, p.UserID); err != nil {
					return err
				}
			}
		}
		if err := tx.insertAudit(ctx, p.UserID, string(p.Role), config.AuditActionApprove, target.AssetCode, target.ID); err != nil {
			return err
		}
		// FR-APV-02: Email ผลอนุมัติลงคิวใน transaction เดียวกัน ตัวส่งเบื้องหลังส่งภายใน 60 วินาที (T-15)
		if err := tx.insertEmail(ctx, target.ID, config.ApprovalKindApproved, s.now()); err != nil {
			return err
		}
		rejected, err := s.rejectSuperseded(ctx, tx, p, target, all)
		if err != nil {
			return err
		}
		out = Result{RequestID: target.ID, Status: config.StatusApproved, DecidedAt: at, RejectedBy: rejected}
		return nil
	})
	return out, err
}

// ErrForbiddenDepartment เกิดเมื่อเจ้าหน้าที่พิจารณาคำขอของอุปกรณ์ภาควิชาอื่น (ASM-06, AC-APV-12, FORBIDDEN_DEPARTMENT)
var ErrForbiddenDepartment = errors.New("คำขอนี้เป็นของอุปกรณ์ภาควิชาอื่น")

// ErrApprovalCutoff เกิดเมื่ออนุมัติหลัง 16:00 น. ของวันที่เริ่มยืม (FR-APV-01 ค ซ, AC-APV-16, APPROVAL_CUTOFF_PASSED)
// ข้อความภายในไม่มีตัวเลข ข้อความที่ผู้ใช้เห็นมาจาก config (rtm.md F-06 ทีมตัดสินแก้โค้ด 2569-11-13)
var ErrApprovalCutoff = errors.New("เลยเวลาตัดของวันที่เริ่มยืมแล้ว อนุมัติไม่ได้")

// approvalCutoff คือเวลาตัด = config.LoanNormalReturnHour (16:00 น.) ของวันที่เริ่มยืม ตามเวลาไทย (glossary)
// ใช้ค่าของ 002 ร่วม ไม่ประกาศ 16 ซ้ำ (AGENTS.md กฎข้อ 9, DOM-LOAN-01)
func approvalCutoff(startAt time.Time) time.Time {
	lt := localTime(startAt)
	return time.Date(lt.Year(), lt.Month(), lt.Day(), config.LoanNormalReturnHour, 0, 0, 0, lt.Location())
}

// ErrItemNotAvailable เกิดเมื่ออุปกรณ์รายชิ้นไม่ได้อยู่ในสถานะ "ว่าง" (FR-APV-01 ค ซ, AC-APV-08, ITEM_NOT_AVAILABLE)
var ErrItemNotAvailable = errors.New("อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว")

// ErrNotFirst เกิดเมื่ออนุมัติคำขอที่ไม่ใช่คำขอแรกของอุปกรณ์ชิ้นนั้น (FR-APV-01 ข ซ, AC-APV-07, NOT_FIRST_IN_QUEUE)
var ErrNotFirst = errors.New("ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน")

// firstInQueue คืนคำขอแรกของชิ้นนั้น = คำขอยืมปกติ requested ที่ส่งก่อนสุด (FR-APV-01 ข, Q-APV-23)
// all ต้องเรียงตาม submitted_at, id แล้ว แบบเดียวกับ is_first ในหน้ารายการ (ListPending)
func firstInQueue(all []requestRow) (requestRow, bool) {
	for _, r := range all {
		if r.Status == config.StatusRequested && r.LoanType == config.LoanTypeNormal {
			return r, true
		}
	}
	return requestRow{}, false
}

// rejectSuperseded ปฏิเสธคำขออื่นที่ยังเป็น requested ของอุปกรณ์ชิ้นเดียวกับคำขอที่เพิ่งอนุมัติ (FR-APV-01 จ, Q-SH-06)
//   - รวมคำขอค้างคืนด้วย (Q-APV-14) แถวทั้งหมดถูกล็อกไว้แล้วใน lockItemRequests
//   - เหตุผล config.ApprovalSupersededReason ผู้ดำเนินการคือเจ้าหน้าที่ที่กดอนุมัติ
//   - แต่ละใบมี approval_decision, audit_log action reject และ Email ผลปฏิเสธ (FR-APV-02 ข, RQ-25)
//
// คืนรหัสคำขอที่ถูกปฏิเสธตาม ให้หน้าจอเอาออกจากรายการ (plan.md ข้อ 4)
func (s *Service) rejectSuperseded(ctx context.Context, tx *Repository, p auth.Principal, approved requestRow, all []requestRow) ([]int, error) {
	rejected := []int{}
	for _, other := range all {
		if other.ID == approved.ID || other.Status != config.StatusRequested {
			continue
		}
		if err := tx.setStatus(ctx, other.ID, config.StatusRejected); err != nil {
			return nil, err
		}
		if _, err := tx.insertDecision(ctx, decision{
			RequestID:         other.ID,
			Result:            config.ApprovalKindRejected,
			Reason:            config.ApprovalSupersededReason,
			ActorID:           p.UserID,
			ActorDepartmentID: p.DepartmentID,
		}); err != nil {
			return nil, err
		}
		if err := tx.insertAudit(ctx, p.UserID, string(p.Role), config.AuditActionReject, other.AssetCode, other.ID); err != nil {
			return nil, err
		}
		if err := tx.insertEmail(ctx, other.ID, config.ApprovalKindRejected, s.now()); err != nil {
			return nil, err
		}
		rejected = append(rejected, other.ID)
	}
	return rejected, nil
}

// ErrReasonRequired เกิดเมื่อปฏิเสธโดยไม่ระบุเหตุผล (FR-APV-01 ง, AC-APV-06, CON-TECH-03 VALIDATION_FAILED ช่อง reason)
var ErrReasonRequired = errors.New("ต้องระบุเหตุผลในการปฏิเสธ")

// ErrReasonTooLong เกิดเมื่อเหตุผลการปฏิเสธยาวเกิน config.ApprovalRejectReasonMaxLength (FR-APV-01 ง, Q-APV-28, VALIDATION_FAILED ช่อง reason)
var ErrReasonTooLong = errors.New("เหตุผลในการปฏิเสธยาวเกินกำหนด")

// Reject ปฏิเสธคำขอยืมแบบปกติ 1 ใบพร้อมเหตุผลที่เจ้าหน้าที่ระบุ ใน transaction เดียว (FR-APV-01 ง ฌ, AC-APV-02)
//   - คำขอเป็น rejected บันทึกเหตุผลใน approval_decision และ audit_log action reject (RQ-25)
//   - ไม่เปลี่ยนสถานะอุปกรณ์ (CON-DATA-01)
//
// พิจารณาได้เฉพาะอุปกรณ์ภาควิชาเดียวกับเจ้าหน้าที่ (ASM-06, T-25)
func (s *Service) Reject(ctx context.Context, requestID int, reason string, loadedAt time.Time) (Result, error) {
	p, err := staff(ctx)
	if err != nil {
		return Result{}, err
	}
	// ลำดับการตรวจ (Q-APV-29, Q-APV-34): หาคำขอ -> ภาควิชา -> เหตุผล -> สถานะ
	// ตรวจเหตุผลหลังภาควิชา เจ้าหน้าที่ภาควิชาอื่นจึงได้ FORBIDDEN_DEPARTMENT เสมอ และก่อนสถานะ คำขอที่ยกเลิกแล้วจึงได้ VALIDATION_FAILED
	var out Result
	err = s.repo.inTx(ctx, func(tx *Repository) error {
		target, _, ok, err := tx.lockItemRequests(ctx, requestID)
		if err != nil {
			return err
		}
		if !ok || target.LoanType != config.LoanTypeNormal {
			return ErrRequestNotFound
		}
		// ASM-06, IF-IDP-01: พิจารณาได้เฉพาะคำขอของอุปกรณ์ภาควิชาเดียวกับเจ้าหน้าที่ (AC-APV-12, T-25)
		// ตรวจก่อนเรื่องอื่น ผู้ที่อยู่ภาควิชาอื่นจึงไม่รู้สถานะหรือข้อมูลใดของคำขอนั้น
		if target.DepartmentID != p.DepartmentID {
			return ErrForbiddenDepartment
		}
		// FR-APV-01 (ง)(ซ), AC-APV-06: ต้องระบุเหตุผล ว่างหรือมีแต่ช่องว่างไม่นับ ยังไม่เขียนอะไร สถานะจึงคงเดิม
		if strings.TrimSpace(reason) == "" {
			return ErrReasonRequired
		}
		// FR-APV-01 (ง), Q-APV-28: ไม่เกิน ApprovalRejectReasonMaxLength นับตาม code point แบบเดียวกับ 003
		if utf8.RuneCountInString(reason) > config.ApprovalRejectReasonMaxLength {
			return ErrReasonTooLong
		}
		if target.Status != config.StatusRequested {
			return notRequestedError(ctx, tx, target, loadedAt) // FR-APV-01 (ซ), T-23
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
		// FR-APV-02 (ข): Email ผลปฏิเสธลงคิวใน transaction เดียวกัน (T-15)
		if err := tx.insertEmail(ctx, target.ID, config.ApprovalKindRejected, s.now()); err != nil {
			return err
		}
		out = Result{RequestID: target.ID, Status: config.StatusRejected, DecidedAt: at, RejectedBy: []int{}}
		return nil
	})
	return out, err
}

// ListFailedEmails คืนรายการ "Email ที่ส่งไม่สำเร็จ" ของภาควิชาเจ้าหน้าที่ (FR-APV-02 ง, ASM-08, AC-APV-14)
func (s *Service) ListFailedEmails(ctx context.Context) ([]FailedEmail, error) {
	p, err := staff(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.ListFailedEmails(ctx, p.DepartmentID)
}

// ErrEmailNotFound เกิดเมื่อไม่มี Email รหัสนี้ (ตอบด้วย REQUEST_NOT_FOUND ที่ลงทะเบียนแล้ว Q-APV-25)
var ErrEmailNotFound = errors.New("ไม่พบ Email นี้")

// Resend ส่ง Email ที่ส่งไม่สำเร็จกลับเข้าคิวให้ตัวส่งเบื้องหลังส่งภายในรอบถัดไป โดยไม่พิจารณาใหม่ (FR-APV-02 ง, AC-APV-14, T-33)
//   - เฉพาะ Email ของอุปกรณ์ภาควิชาเดียวกับเจ้าหน้าที่ (ASM-06) ไม่อย่างนั้น ErrForbiddenDepartment
//   - Email ที่ไม่ได้ failed (ส่งแล้วหรืออยู่ในคิวแล้ว) ไม่ทำอะไร คืนสถานะปัจจุบัน กดซ้ำได้ไม่เกิด Email ซ้ำ
func (s *Service) Resend(ctx context.Context, emailID int) (string, error) {
	p, err := staff(ctx)
	if err != nil {
		return "", err
	}
	var status string
	err = s.repo.inTx(ctx, func(tx *Repository) error {
		st, dept, ok, err := tx.lockEmail(ctx, emailID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrEmailNotFound
		}
		if dept != p.DepartmentID {
			return ErrForbiddenDepartment
		}
		status = st
		if st != config.ApprovalEmailFailed {
			return nil
		}
		status = config.ApprovalEmailPending
		return tx.requeueEmail(ctx, emailID, s.now())
	})
	return status, err
}
