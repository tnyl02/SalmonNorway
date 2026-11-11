// Package equipment คือฟีเจอร์ 001-equipment-status (UC-01 ตรวจสอบสถานะอุปกรณ์)
// เป็นฟีเจอร์อ่านอย่างเดียว ไม่มีการเปลี่ยนสถานะอุปกรณ์จากที่นี่ (CON-DATA-02)
package equipment

import "time"

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

// รายละเอียดอุปกรณ์รายชิ้น 1 ชิ้น แยก struct ตามบทบาท (FR-EQS-04, FR-EQS-07, DOM-PRIV-01, plan.md ข้อ 6.1)
// ฟิลด์ที่บทบาทนั้นไม่มีสิทธิ์เห็นต้องไม่อยู่ใน struct เลย กันข้อมูลหลุดโดยไม่ตั้งใจ
//   - dueAt (กำหนดคืน วันที่และเวลา) มีเฉพาะสถานะ "ถูกยืม" (CR-002, DOM-LOAN-01)
//   - isOverdue = ถูกยืมและเวลาปัจจุบันเลยกำหนดคืนแล้ว เป็นป้าย ไม่ใช่สถานะที่ 4 สถานะยังเป็น borrowed (CON-DATA-01)
//   - ปุ่มส่งคำขอยืม (T-34 รอ Q-16) และชื่อผู้ยืมของเจ้าหน้าที่ (T-36 รอ Q-20) ยังไม่มี

// StudentItemDTO คือสิ่งที่นักศึกษาได้ เมื่อเลยกำหนดคืนแล้วไม่ส่ง dueAt ส่งแค่ป้าย (FR-EQS-04 ข้อ ข, CR-001)
type StudentItemDTO struct {
	AssetCode string     `json:"assetCode"`
	Status    Status     `json:"status"`
	DueAt     *time.Time `json:"dueAt,omitempty"`
	IsOverdue bool       `json:"isOverdue,omitempty"`
}

// StaffItemDTO คือสิ่งที่เจ้าหน้าที่ได้ เห็นทั้งกำหนดคืนและป้ายเลยกำหนดคืน (FR-EQS-04 ข้อ ข)
type StaffItemDTO struct {
	AssetCode string     `json:"assetCode"`
	Status    Status     `json:"status"`
	DueAt     *time.Time `json:"dueAt,omitempty"`
	IsOverdue bool       `json:"isOverdue,omitempty"`
}

// ItemDetailDTO คือชื่อเดิมของรายละเอียดจาก T-10 ใช้อ่าน JSON ใน test (มีฟิลด์ครบทุกตัวที่ส่งได้)
type ItemDetailDTO = StaffItemDTO

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
