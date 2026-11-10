package equipment

import "context"

// Service คือกฎของฟีเจอร์ที่อยู่ระหว่าง handler กับ repository
// ตอนนี้ส่งต่ออย่างเดียว กฎเรื่องเลยกำหนดคืนและการแยกข้อมูลตามบทบาทจะเพิ่มใน T-12 และ T-17
type Service struct {
	repo *Repository
}

// NewService สร้าง Service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ListModels คืนรายการรุ่นพร้อมจำนวนว่าง/ทั้งหมด ตามตัวกรอง (FR-EQS-01, FR-EQS-02)
func (s *Service) ListModels(ctx context.Context, f Filter) ([]ModelRowDTO, error) {
	return s.repo.ListModels(ctx, f)
}

// ListItems คืนอุปกรณ์รายชิ้นของรุ่นที่เลือก ตามตัวกรองเดียวกัน (FR-EQS-01, FR-EQS-02)
func (s *Service) ListItems(ctx context.Context, modelID int, f Filter) ([]ItemRow, error) {
	return s.repo.ListItems(ctx, modelID, f)
}
