package borrowrequest

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"salmonnorway/backend/internal/audit"
	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// Service คือกฎของฟีเจอร์ที่อยู่ระหว่าง handler กับ repository
// ทุกกฎบังคับที่ฝั่ง server ไม่พึ่งการซ่อนปุ่มที่หน้าจอ (plan.md ข้อ 1)
type Service struct {
	repo        *Repository
	clock       Clock            // นาฬิกาที่ฉีดได้ ใช้เทียบกฎคืนเกินกำหนด (FR-BRQ-06, T-12)
	loc         *time.Location   // เขตเวลา Asia/Bangkok (config.TimeZone, glossary)
	attachments *AttachmentStore // ที่เก็บไฟล์ใบบันทึกข้อความ (FR-BRQ-10 ข้อ ก, ASM-BRQ-07)
}

// NewService สร้าง Service ใช้นาฬิกาจริงและเขตเวลาจาก config
func NewService(repo *Repository) *Service {
	return NewServiceWithClock(repo, SystemClock{})
}

// NewServiceWithClock สร้าง Service ด้วยนาฬิกาที่กำหนดเอง ใช้ใน test ที่ต้องตรึงเวลา (AC-BRQ-06)
func NewServiceWithClock(repo *Repository, clock Clock) *Service {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		panic("load time zone " + config.TimeZone + ": " + err.Error()) // ฝัง tzdata ไว้แล้ว ไม่ควรเกิด
	}
	// โฟลเดอร์ไฟล์แนบอ่านจาก config ที่เดียว เปลี่ยนด้วย ATTACHMENT_DIR ได้ (ASM-BRQ-07, plan.md ข้อ 5)
	return &Service{repo: repo, clock: clock, loc: loc,
		attachments: NewAttachmentStore(config.Load().AttachmentDir)}
}

// ErrHasOverdueLoan เกิดเมื่อผู้ขอมีอุปกรณ์ที่คืนเกินกำหนดค้างอยู่ จึงส่งคำขอใหม่ไม่ได้
// (FR-BRQ-06, CON-TECH-03 รหัส HAS_OVERDUE_LOAN)
var ErrHasOverdueLoan = errors.New("มีอุปกรณ์ที่คืนเกินกำหนดค้างอยู่")

// ลำดับการตรวจเมื่อเงื่อนไขซ้อนกัน (Q-BRQ-15 ปิดแล้ว 2569-11-12, T-21)
// ผู้ขอต่างภาควิชาที่ขอยืมอุปกรณ์ซึ่งไม่ว่างอยู่แล้ว ต้องได้ ErrForbiddenDepartment ไม่ใช่ ErrItemNotAvailable
// เพราะ FR-BRQ-07 ห้ามส่งข้อมูลอุปกรณ์กลับไปเลย ส่วน FR-BRQ-01 ที่ให้แจ้งสถานะล่าสุด
// ใช้กับผู้ขอที่มีสิทธิ์เห็นอุปกรณ์ชิ้นนั้นเท่านั้น การตรวจจริงอยู่ที่ repository.Create ซึ่งเรียก checkDepartment ก่อน
//
// Create สร้างคำขอยืมใบใหม่ด้วยสถานะเริ่มต้น "รออนุมัติ" (config.StatusRequested)
// บันทึกผู้ขอจากรหัสผู้ใช้ที่ยืนยันตัวตนแล้วเท่านั้น ไม่ใช่รหัสที่พิมพ์ในฟอร์ม (FR-BRQ-04, ASM-BRQ-06)
// และต้องไม่เปลี่ยนสถานะอุปกรณ์ อุปกรณ์ที่มีคำขอค้างยังเป็น "ว่าง" (DOM-REQ-01)
//
// ส่งคำขอซ้อนกับคำขอของ**ผู้อื่น**บนอุปกรณ์ชิ้นเดียวกันได้ตามปกติ (FR-BRQ-05 ก, DOM-REQ-01)
// แต่ถ้าเป็นผู้ขอ**คนเดิม**ที่คำขอเก่ายังไม่จบ จะคืน ErrDuplicateRequest และไม่สร้างคำขอใหม่
// (FR-BRQ-05 ข, AC-BRQ-11) การกันซ้ำบังคับด้วย unique index ที่ฐานข้อมูล ไม่ใช่การอ่านมาตรวจก่อนเขียน
// จึงไม่มีช่องให้สองคำขอที่มาถึงพร้อมกันผ่านไปได้ทั้งคู่
func (s *Service) Create(ctx context.Context, in NewRequest) (Request, error) {
	p, ok := auth.FromContext(ctx)
	if !ok || !p.Complete() {
		return Request{}, ErrNoPrincipal
	}

	// ผู้ขอที่มีของคืนเกินกำหนดค้างอยู่ ส่งคำขอใหม่ไม่ได้ และต้องไม่สร้างคำขอ (FR-BRQ-06, AC-BRQ-06)
	// ตรวจก่อนแตะอุปกรณ์ได้ เพราะเป็นเรื่องของตัวผู้ขอเอง คำตอบไม่เปิดเผยข้อมูลอุปกรณ์ชิ้นใดเลย (DOM-ACCESS-01)
	overdue, err := s.hasOverdueLoan(ctx, p)
	if err != nil {
		return Request{}, err
	}
	if overdue {
		return Request{}, ErrHasOverdueLoan
	}

	req, err := s.repo.Create(ctx, p.UserID, normalizePurpose(in), config.StatusRequested)
	if err != nil {
		return Request{}, err
	}

	// เวลาที่ส่งเทียบด้วยเขตเวลา Asia/Bangkok ตาม glossary (FR-BRQ-04)
	req.SubmittedAt = req.SubmittedAt.In(s.loc)
	return req, nil
}

// CreateWithAttachment สร้างคำขอแล้วเก็บไฟล์ใบบันทึกข้อความของการยืมค้างคืน (FR-BRQ-10 ข้อ ก)
// up เป็น nil ได้ กรณีนั้นเท่ากับเรียก Create ตรง ๆ เช่นการยืมปกติที่ไม่ได้แนบไฟล์
//
// ชื่อไฟล์ต้องมีรหัสคำขออยู่ด้วยตาม plan.md ข้อ 3 จึงต้องสร้างคำขอก่อนแล้วค่อยเก็บไฟล์
// ถ้าเก็บไฟล์หรือบันทึก path ไม่สำเร็จ จะลบคำขอใบนั้นทิ้ง ไม่ปล่อยให้มีคำขอค้างที่ไม่มีไฟล์แนบ (FR-BRQ-09)
func (s *Service) CreateWithAttachment(ctx context.Context, in NewRequest, up *Upload) (Request, error) {
	req, err := s.createRequest(ctx, in, up)
	s.auditRequest(ctx, in.AssetCode, err) // บันทึกทุกผลลัพธ์ ทั้งสำเร็จและถูกปฏิเสธ (FR-BRQ-08)
	return req, err
}

// createRequest คืองานหลักของ CreateWithAttachment แยกออกมาเพื่อให้บันทึก audit ได้ที่เดียวจบ
func (s *Service) createRequest(ctx context.Context, in NewRequest, up *Upload) (Request, error) {
	req, err := s.Create(ctx, in)
	if err != nil || up == nil {
		return req, err
	}

	path, err := s.attachments.Save(req.ID, *up)
	if err != nil {
		_ = s.repo.Delete(ctx, req.ID)
		return Request{}, err
	}
	if err := s.repo.SetAttachmentPath(ctx, req.ID, path); err != nil {
		_ = s.attachments.Remove(path)
		_ = s.repo.Delete(ctx, req.ID)
		return Request{}, err
	}
	return req, nil
}

// ErrAttachmentNotFound เกิดเมื่อไม่มีคำขอใบนั้น หรือคำขอใบนั้นไม่ได้แนบไฟล์ไว้
// รวมสองกรณีไว้ด้วยกัน เพื่อไม่บอกคนนอกว่าเลขคำขอใดมีอยู่จริง (DOM-PRIV-01, FR-BRQ-10 ข้อ ก)
var ErrAttachmentNotFound = errors.New("ไม่พบคำขอหรือคำขอนั้นไม่มีไฟล์แนบ")

// ErrNotRequestOwner เกิดเมื่อผู้ใช้ไม่ใช่เจ้าหน้าที่และไม่ใช่เจ้าของคำขอ (FR-BRQ-10 ข้อ ก, DOM-PRIV-01)
var ErrNotRequestOwner = errors.New("เปิดดูไฟล์แนบได้เฉพาะเจ้าหน้าที่และเจ้าของคำขอ")

// OpenAttachment คืนที่อยู่ไฟล์แนบเมื่อผู้ใช้มีสิทธิ์เปิดดู (FR-BRQ-10 ข้อ ก)
// เปิดดูได้เฉพาะเจ้าหน้าที่กับเจ้าของคำขอ นักศึกษาคนอื่นดูไม่ได้ (ASM-BRQ-07, DOM-PRIV-01)
// เจ้าของคำขอคือรหัสผู้ใช้จากการยืนยันตัวตน ไม่ใช่รหัสนักศึกษาที่พิมพ์ในฟอร์ม (IF-IDP-01, ASM-BRQ-06)
func (s *Service) OpenAttachment(ctx context.Context, requestID int) (string, error) {
	p, ok := auth.FromContext(ctx)
	if !ok || !p.Complete() {
		return "", ErrNoPrincipal
	}

	att, found, err := s.repo.AttachmentOf(ctx, requestID)
	if err != nil {
		return "", err
	}
	if !found {
		return "", ErrAttachmentNotFound
	}
	if p.Role != auth.RoleStaff && att.RequesterID != p.UserID {
		return "", ErrNotRequestOwner
	}
	return att.Path, nil
}

// auditRequest บันทึก audit log ของการส่งคำขอ 1 ครั้ง ทั้งกรณีสำเร็จและถูกปฏิเสธ (FR-BRQ-08, RQ-25, DM-07)
// บันทึกผู้ขอ เวลา (ฐานข้อมูลใส่ให้) รหัสครุภัณฑ์ และผลลัพธ์ ตามที่ AC-BRQ-08 กำหนด
// ใช้ตาราง audit_log ของกลางผ่าน audit.Write และใช้ค่า action จาก config ที่เดียว (AGENTS.md กฎข้อ 9)
//
// ไม่บันทึกเมื่อไม่รู้ว่าผู้ขอเป็นใคร เพราะ audit ต้องระบุผู้เข้าถึงเสมอ (DM-07)
//
// **เขียนแบบ best-effort (ทีมตัดสิน 2569-11-12)** ถ้าบันทึก audit ไม่สำเร็จ ผู้ขอยังได้คำตอบเดิม
// แต่ขึ้นข้อความที่เครื่องให้เห็น เหตุผล: ถ้าล้มคำขอที่สร้างสำเร็จไปแล้วเพราะ log ไม่ติด
// ต้องย้อนลบคำขอและไฟล์แนบ ซึ่งขัด FR-BRQ-09 ที่ห้ามทิ้งคำขอค้าง
// หมายเหตุ: 001 เลือก fail closed (ตอบ 503 เมื่อเขียน audit ไม่ได้) ระบบจึงยังมีสองแบบ
// ต้องยกเข้าที่ประชุมให้เลือกมาตรฐานเดียวกันทั้งระบบ
func (s *Service) auditRequest(ctx context.Context, assetCode string, cause error) {
	p, ok := auth.FromContext(ctx)
	if !ok || !p.Complete() {
		return
	}

	// ไม่ผูก RequestID ไว้กับแถว audit แม้คำขอจะสร้างสำเร็จ
	// เพราะ audit_log เป็น append-only (trigger ของ migration 0002) แต่คำขอถูกลบได้
	// เช่นทางย้อนกลับตอนเก็บไฟล์แนบไม่สำเร็จใน T-14 ถ้ามี foreign key ชี้อยู่ การลบจะทำไม่ได้เลย
	// AC-BRQ-08 ต้องการแค่ผู้ขอ เวลา รหัสครุภัณฑ์ และผลลัพธ์ ซึ่ง 4 ช่องนี้ครบแล้ว
	// (DM-07 เขียนว่า request_id ว่างได้ ถ้าทีมอยากให้ผูกด้วย ต้องยอมรับว่าคำขอนั้นลบไม่ได้อีก)
	entry := audit.Entry{
		ActorID:   p.UserID,
		ActorRole: string(p.Role),
		Action:    config.AuditActionRequest,
		AssetCode: assetCode,
		Outcome:   auditOutcome(cause),
	}

	if _, err := audit.Write(ctx, s.repo.db, entry); err != nil {
		log.Printf("borrowrequest audit: เขียน audit log ไม่สำเร็จ actor=%s asset=%s: %v", p.UserID, assetCode, err)
	}
}

// ErrValidationRejected ใช้บอก audit ว่าคำขอถูกปฏิเสธตั้งแต่การตรวจข้อมูลที่ handler (FR-BRQ-03, FR-BRQ-10)
// ไม่ได้ถูกส่งกลับไปให้ผู้ขอ handler ตอบ 422 VALIDATION_FAILED พร้อมรายชื่อช่องที่ผิดเองอยู่แล้ว
var ErrValidationRejected = errors.New("คำขอถูกปฏิเสธตั้งแต่การตรวจข้อมูล")

// AuditValidationRejected บันทึก audit ของคำขอที่ถูกปฏิเสธก่อนถึง service (FR-BRQ-08, T-24)
// FR-BRQ-08 สั่งให้บันทึก "ทุกครั้งที่มีการส่งคำขอ ทั้งกรณีสำเร็จและถูกปฏิเสธ"
// คำขอที่ข้อมูลไม่ครบหรือไฟล์ผิดกติกาก็นับเป็นการส่งคำขอที่ถูกปฏิเสธ (ทีมตัดสิน 2569-11-12)
func (s *Service) AuditValidationRejected(ctx context.Context, assetCode string) {
	s.auditRequest(ctx, assetCode, ErrValidationRejected)
}

// auditOutcome แปลงผลของการส่งคำขอเป็นค่า outcome 1 ใน 3 ค่าที่ audit_log รับ (DM-07, migration 0002)
//   - สำเร็จ = ok
//   - ถูกปฏิเสธด้วยกฎของฟีเจอร์ (ต่างภาควิชา อุปกรณ์ไม่ว่าง คำขอซ้ำ มีของเกินกำหนด) = forbidden
//   - ระบบขัดข้อง เช่น ฐานข้อมูลหรือที่เก็บไฟล์ไม่ตอบ = error
//
// ตารางมีแค่ 3 ค่า การแยก "ถูกปฏิเสธตามกฎ" ออกจาก "ไม่มีสิทธิ์" ต้องเพิ่มค่าใหม่ใน CHECK ซึ่งต้องทำ CR
func auditOutcome(cause error) string {
	switch {
	case cause == nil:
		return audit.OutcomeOK
	case errors.Is(cause, ErrForbiddenDepartment),
		errors.Is(cause, ErrItemNotAvailable),
		errors.Is(cause, ErrDuplicateRequest),
		errors.Is(cause, ErrHasOverdueLoan),
		errors.Is(cause, ErrValidationRejected):
		return audit.OutcomeForbidden
	default:
		return audit.OutcomeError
	}
}

// hasOverdueLoan บอกว่าผู้ขอยังถูกกฎคืนเกินกำหนดห้ามส่งคำขออยู่หรือไม่ (FR-BRQ-06)
// นับเฉพาะอุปกรณ์ของภาควิชาผู้ขอเอง ของค้างในภาควิชาอื่นไม่นับ (FR-BRQ-06 ข้อ ค, DOM-ACCESS-01)
// มี 2 กรณีที่ห้ามส่ง
//   - ข้อ ก ยังไม่คืน และเลยเวลา 16:00 น. ของกำหนดคืนไปแล้ว
//   - ข้อ ข คืนไปแล้ว แต่คืนช้า และยังไม่ถึง 00:00 น. ของวันถัดไปจากวันที่คืน
func (s *Service) hasOverdueLoan(ctx context.Context, p auth.Principal) (bool, error) {
	now := s.clock.Now()

	// ข้อ ก ใบที่ยังไม่คืน
	dueAt, found, err := s.repo.EarliestOpenLoanDue(ctx, p.UserID, p.DepartmentID)
	if err != nil {
		return false, err
	}
	if found && now.After(overdueDeadline(dueAt, s.loc)) {
		return true, nil
	}

	if !config.OverdueBlockUntilNextDayStart {
		return false, nil // ทีมปิดกฎข้อ ข ไว้ คืนแล้วส่งคำขอได้ทันที
	}
	return s.returnedLateToday(ctx, p, now)
}

// returnedLateToday บอกว่าผู้ขอคืนของช้าไปแล้วในวันนี้หรือไม่ (FR-BRQ-06 ข้อ ข, AC-BRQ-12)
// ทีมตัดสิน 2569-11-11 ว่าปลดล็อกเมื่อถึง 00:00 น. ของวันถัดไป ไม่ใช่ทันทีที่คืนและไม่ใช่ครบ 24 ชั่วโมง
// เทียบว่า "วันนี้" คือวันไหนด้วยเขตเวลา Asia/Bangkok ตาม glossary การคืนเมื่อวานจึงปลดล็อกแล้วเสมอ
//
// "คืนช้า" ใช้เส้นแบ่งเดียวกับข้อ ก คือคืนหลัง 16:00 น. ของกำหนดคืน (overdueDeadline)
// คืนตรงเวลาไม่เข้ากฎนี้ เพราะ FR-BRQ-06 ทั้งข้อพูดถึงผู้ที่คืนเกินกำหนดเท่านั้น
func (s *Service) returnedLateToday(ctx context.Context, p auth.Principal, now time.Time) (bool, error) {
	returned, err := s.repo.ReturnedSince(ctx, p.UserID, p.DepartmentID, startOfDay(now, s.loc))
	if err != nil {
		return false, err
	}
	for _, l := range returned {
		if l.ReturnedAt.After(overdueDeadline(l.DueAt, s.loc)) {
			return true, nil
		}
	}
	return false, nil
}

// startOfDay คืนเวลา 00:00 น. ของวันเดียวกับ t ตามเขตเวลาที่ให้ (FR-BRQ-06 ข้อ ข)
func startOfDay(t time.Time, loc *time.Location) time.Time {
	d := t.In(loc)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
}

// normalizePurpose จัดค่าวัตถุประสงค์ก่อนบันทึก (FR-BRQ-02)
// ถ้าไม่ได้เลือก "อื่น ๆ" ข้อความที่พิมพ์ระบุจะไม่ถูกเก็บ ตาม FR-BRQ-02 ข้อ ข ที่ว่าไม่ต้องพิมพ์อะไรเพิ่ม
// กันไม่ให้มีข้อความค้างในฐานข้อมูลจากการที่ผู้ขอพิมพ์ไว้แล้วเปลี่ยนใจไปเลือกค่าอื่น
func normalizePurpose(in NewRequest) NewRequest {
	in.Purpose = strings.TrimSpace(in.Purpose)
	in.PurposeOther = strings.TrimSpace(in.PurposeOther)
	if in.Purpose != config.PurposeOther {
		in.PurposeOther = ""
	}
	return in
}

// ItemNotAvailableError คือ ErrItemNotAvailable ที่พกสถานะล่าสุดของอุปกรณ์มาด้วย (FR-BRQ-01, AC-BRQ-01)
// handler ใช้ errors.As ดึงสถานะไปแจ้งผู้ขอว่าอุปกรณ์เป็นอะไรอยู่ (CON-DATA-01 มี 3 ค่า)
//
// รองรับ errors.Is(err, ErrItemNotAvailable) ด้วย โค้ดเดิมที่เทียบกับ ErrItemNotAvailable จึงใช้ได้เหมือนเดิม
// (ประกาศไว้ที่นี่เพราะ model.go ไม่อยู่ในช่องไฟล์ที่แตะของ T-09 ย้ายไปรวมกับ error ตัวอื่นได้ภายหลัง)
type ItemNotAvailableError struct {
	AssetCode string
	Status    string // ค่าสถานะล่าสุดตาม CON-DATA-01 (available / borrowed / unavailable)
}

func (e *ItemNotAvailableError) Error() string {
	return "อุปกรณ์ " + e.AssetCode + " อยู่ในสถานะ " + e.Status + " จึงส่งคำขอไม่ได้"
}

// Is ทำให้ errors.Is(err, ErrItemNotAvailable) เป็นจริง
func (e *ItemNotAvailableError) Is(target error) bool { return target == ErrItemNotAvailable }

// ลบ Service.ItemStatus ออกเมื่อ 2569-11-13 ตาม F-08 ใน rtm.md
// เป็น wrapper ที่ไม่มีผู้เรียกในโค้ดจริง การตรวจสถานะอุปกรณ์เกิดในคำสั่งเดียวกับ INSERT
// ที่ Repository.Create อยู่แล้ว (FR-BRQ-01) ส่วน Repository.ItemStatus ยังใช้อยู่จึงไม่ลบ
