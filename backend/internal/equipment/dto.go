// Package equipment คือฟีเจอร์ 001-equipment-status (UC-01 ตรวจสอบสถานะอุปกรณ์)
// เป็นฟีเจอร์อ่านอย่างเดียว ไม่มีการเปลี่ยนสถานะอุปกรณ์จากที่นี่ (CON-DATA-02)
package equipment

// Status คือสถานะของอุปกรณ์รายชิ้น มี 3 ค่าเท่านั้น (CON-DATA-01, CR-002)
// ค่าตรงกับ CHECK ของคอลัมน์ equipment_item.status ใน db/migrations/0001_equipment_read_model.up.sql
type Status string

const (
	StatusAvailable   Status = "available"   // ว่าง
	StatusBorrowed    Status = "borrowed"    // ถูกยืม
	StatusUnavailable Status = "unavailable" // ไม่พร้อมใช้งาน
)

// ModelRowDTO คือ 1 แถวในหน้ารายการอุปกรณ์ 1 แถว = 1 รุ่น (FR-EQS-01, Q-13)
// availableCount นับเฉพาะชิ้นที่สถานะ "ว่าง" ชิ้นที่มีคำขอยืมค้างยังนับเป็นว่าง (DOM-REQ-01)
type ModelRowDTO struct {
	ModelID        int    `json:"modelId"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	AvailableCount int    `json:"availableCount"`
	TotalCount     int    `json:"totalCount"`
}

// ListDTO คือรูปแบบคำตอบของ endpoint ที่คืนรายการ (plan.md ข้อ 5: ไม่พบ = items ว่าง ไม่ใช่ error)
type ListDTO[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

// ErrorDTO คือรูปแบบคำตอบเมื่อเกิดข้อผิดพลาด (plan.md ข้อ 5 Error contract)
type ErrorDTO struct {
	Code      string `json:"code"`
	Retryable bool   `json:"retryable,omitempty"`
}
