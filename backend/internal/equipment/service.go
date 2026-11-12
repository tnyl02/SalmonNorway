package equipment

import (
	"context"
	"errors"
	"time"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// Service คือกฎของฟีเจอร์ที่อยู่ระหว่าง handler กับ repository
// กฎเลยกำหนดคืนและการแยกข้อมูลตามบทบาทคำนวณที่นี่ ฝั่ง server ไม่พึ่งการซ่อนที่หน้าจอ (plan.md ข้อ 5, 6.1)
type Service struct {
	repo *Repository
	now  func() time.Time // นาฬิกา test ตั้งเวลาเองได้
	loc  *time.Location   // เขตเวลา Asia/Bangkok (config.TimeZone)
}

// NewService สร้าง Service ใช้นาฬิกาจริงและเขตเวลาจาก config
func NewService(repo *Repository) *Service {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		panic("load time zone " + config.TimeZone + ": " + err.Error()) // ฝัง tzdata ไว้แล้ว ไม่ควรเกิด
	}
	return &Service{repo: repo, now: time.Now, loc: loc}
}

// ListModels คืนรายการรุ่นพร้อมจำนวนว่าง/ทั้งหมด ตามตัวกรอง (FR-EQS-01, FR-EQS-02)
func (s *Service) ListModels(ctx context.Context, f Filter) ([]ModelRowDTO, error) {
	return s.repo.ListModels(ctx, f)
}

// ListItems คืนอุปกรณ์รายชิ้นของรุ่นที่เลือก ตามตัวกรองเดียวกัน (FR-EQS-01, FR-EQS-02)
func (s *Service) ListItems(ctx context.Context, modelID int, f Filter) ([]ItemRow, error) {
	return s.repo.ListItems(ctx, modelID, f)
}

// GetItem คืนรายละเอียดอุปกรณ์รายชิ้นตามบทบาทของผู้ใช้ใน context (FR-EQS-04, DOM-PRIV-01)
//   - กำหนดคืนส่งเฉพาะสถานะ "ถูกยืม" (plan.md ข้อ 5)
//   - เลยกำหนดคืน = วันที่และเวลาปัจจุบันตาม Asia/Bangkok > กำหนดคืน (FR-EQS-04 ข้อ ข, CR-002)
//   - นักศึกษา: เลยกำหนดแล้วไม่ส่งกำหนดคืน ส่งแค่ป้าย / เจ้าหน้าที่: ได้ทั้งคู่ (CR-001)
//
// ผู้ใช้ที่ไม่ใช่เจ้าหน้าที่ (รวมกรณีไม่มีผู้ใช้ใน context) ได้ข้อมูลแบบนักศึกษา ให้สิทธิ์น้อยที่สุดไว้ก่อน
func (s *Service) GetItem(ctx context.Context, assetCode string) (any, bool, error) {
	rec, found, err := s.repo.GetItem(ctx, assetCode)
	if err != nil || !found {
		return nil, found, err
	}

	var dueAt *time.Time
	overdue := false
	if rec.Status == StatusBorrowed && rec.DueAt != nil {
		dueAt = rec.DueAt
		overdue = s.now().In(s.loc).After(rec.DueAt.In(s.loc))
	}

	if p, ok := auth.FromContext(ctx); ok && p.Role == auth.RoleStaff {
		return StaffItemDTO{AssetCode: rec.AssetCode, Status: rec.Status, DueAt: dueAt, IsOverdue: overdue}, true, nil
	}
	student := StudentItemDTO{AssetCode: rec.AssetCode, Status: rec.Status, IsOverdue: overdue}
	if !overdue {
		student.DueAt = dueAt
	}
	return student, true, nil
}

// ErrStaffOnly หมายถึงผู้ใช้ไม่ใช่เจ้าหน้าที่ จึงไม่มีสิทธิ์ดูข้อมูลนี้ (FR-EQS-07) handler ตอบ 403 FORBIDDEN_ROLE
var ErrStaffOnly = errors.New("staff only")

// Summary คืนจำนวนอุปกรณ์แยกตามสถานะ ให้เฉพาะเจ้าหน้าที่ (FR-EQS-07 "นักศึกษาไม่เห็น", plan.md ข้อ 6.1 ตรวจสิทธิ์ที่ server)
// ไม่มีผู้ใช้ใน context ถือว่าไม่ใช่เจ้าหน้าที่ (ให้สิทธิ์น้อยที่สุด)
func (s *Service) Summary(ctx context.Context) (StatusSummaryDTO, error) {
	if p, ok := auth.FromContext(ctx); !ok || p.Role != auth.RoleStaff {
		return StatusSummaryDTO{}, ErrStaffOnly
	}
	return s.repo.CountByStatus(ctx)
}
