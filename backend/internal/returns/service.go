package returns

import (
	"context"
	"errors"

	"salmonnorway/backend/internal/config"
)

// ErrUnknownCondition หมายถึงไม่ได้ระบุสภาพ หรือสภาพไม่ใช่ ปกติ / ชำรุด (FR-RET-02)
// T-04 แค่ป้องกันไม่ให้บันทึกค่าที่ไม่รู้จัก ข้อความแจ้งและ test ของ AC-RET-04 อยู่ใน T-09
var ErrUnknownCondition = errors.New("condition must be normal or damaged")

// Service คือกฎของการรับคืน แยกจาก HTTP และ SQL (plan.md ข้อ 4)
type Service struct {
	repo *Repository
}

// NewService สร้าง Service ของฟีเจอร์รับคืน
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// RecordReturn บันทึกการรับคืนของอุปกรณ์ assetCode (FR-RET-02, FR-RET-03, FR-RET-05)
// สถานะหลังคืนมาจาก config.ReturnConditionStatus ที่เดียว ไม่ฝังในโค้ด (AGENTS.md กฎข้อ 9)
func (s *Service) RecordReturn(ctx context.Context, assetCode, condition string, photos []Photo) (ReturnRecord, error) {
	newStatus, ok := config.ReturnConditionStatus[condition]
	if !ok {
		return ReturnRecord{}, ErrUnknownCondition
	}
	return s.repo.RecordReturn(ctx, assetCode, condition, newStatus, photos)
}
