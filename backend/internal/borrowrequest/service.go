package borrowrequest

import (
	"context"
	"errors"
	"strings"
	"time"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// Service คือกฎของฟีเจอร์ที่อยู่ระหว่าง handler กับ repository
// ทุกกฎบังคับที่ฝั่ง server ไม่พึ่งการซ่อนปุ่มที่หน้าจอ (plan.md ข้อ 1)
type Service struct {
	repo  *Repository
	clock Clock          // นาฬิกาที่ฉีดได้ ใช้เทียบกฎคืนเกินกำหนด (FR-BRQ-06, T-12)
	loc   *time.Location // เขตเวลา Asia/Bangkok (config.TimeZone, glossary)
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
	return &Service{repo: repo, clock: clock, loc: loc}
}

// ErrHasOverdueLoan เกิดเมื่อผู้ขอมีอุปกรณ์ที่คืนเกินกำหนดค้างอยู่ จึงส่งคำขอใหม่ไม่ได้
// (FR-BRQ-06, CON-TECH-03 รหัส HAS_OVERDUE_LOAN)
var ErrHasOverdueLoan = errors.New("มีอุปกรณ์ที่คืนเกินกำหนดค้างอยู่")

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

// ItemStatus คืนสถานะล่าสุดของอุปกรณ์ ใช้บอกผู้ขอเมื่อส่งคำขอไม่ได้ (FR-BRQ-01)
func (s *Service) ItemStatus(ctx context.Context, assetCode string) (string, error) {
	return s.repo.ItemStatus(ctx, assetCode)
}
