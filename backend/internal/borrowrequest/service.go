package borrowrequest

import (
	"context"
	"time"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// Service คือกฎของฟีเจอร์ที่อยู่ระหว่าง handler กับ repository
// ทุกกฎบังคับที่ฝั่ง server ไม่พึ่งการซ่อนปุ่มที่หน้าจอ (plan.md ข้อ 1)
type Service struct {
	repo *Repository
	now  func() time.Time // นาฬิกา test ตั้งเวลาเองได้ (ใช้เต็มที่ใน T-12, T-13)
	loc  *time.Location   // เขตเวลา Asia/Bangkok (config.TimeZone, glossary)
}

// NewService สร้าง Service ใช้นาฬิกาจริงและเขตเวลาจาก config
func NewService(repo *Repository) *Service {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		panic("load time zone " + config.TimeZone + ": " + err.Error()) // ฝัง tzdata ไว้แล้ว ไม่ควรเกิด
	}
	return &Service{repo: repo, now: time.Now, loc: loc}
}

// Create สร้างคำขอยืมใบใหม่ด้วยสถานะเริ่มต้น "รออนุมัติ" (config.StatusRequested)
// บันทึกผู้ขอจากรหัสผู้ใช้ที่ยืนยันตัวตนแล้วเท่านั้น ไม่ใช่รหัสที่พิมพ์ในฟอร์ม (FR-BRQ-04, ASM-BRQ-06)
// และต้องไม่เปลี่ยนสถานะอุปกรณ์ อุปกรณ์ที่มีคำขอค้างยังเป็น "ว่าง" (DOM-REQ-01)
func (s *Service) Create(ctx context.Context, in NewRequest) (Request, error) {
	p, ok := auth.FromContext(ctx)
	if !ok || !p.Complete() {
		return Request{}, ErrNoPrincipal
	}

	req, err := s.repo.Create(ctx, p.UserID, in, config.StatusRequested)
	if err != nil {
		return Request{}, err
	}

	// เวลาที่ส่งเทียบด้วยเขตเวลา Asia/Bangkok ตาม glossary (FR-BRQ-04)
	req.SubmittedAt = req.SubmittedAt.In(s.loc)
	return req, nil
}

// ItemStatus คืนสถานะล่าสุดของอุปกรณ์ ใช้บอกผู้ขอเมื่อส่งคำขอไม่ได้ (FR-BRQ-01)
func (s *Service) ItemStatus(ctx context.Context, assetCode string) (string, error) {
	return s.repo.ItemStatus(ctx, assetCode)
}
