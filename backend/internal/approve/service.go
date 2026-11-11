package approve

import (
	"context"
	"errors"
	"time"

	"salmonnorway/backend/internal/auth"
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
