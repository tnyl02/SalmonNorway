// ค่าที่หน้าจอใช้ซึ่งมาจาก spec รวมไว้ที่เดียว (AGENTS.md กฎข้อ 9)
// ต้องตรงกับ backend/internal/config/config.go เสมอ หลังบ้านเป็นผู้ตรวจจริง หน้าจอตรวจซ้ำเพื่อแจ้งผู้ใช้ก่อนส่ง

// เขตเวลาที่ใช้แสดงวันที่และเวลา (specs/000-shared/glossary.md)
export const TIMEZONE = 'Asia/Bangkok'

// ชื่อสถานะอุปกรณ์ที่แสดงบนหน้าจอ 3 ค่า (CON-DATA-01, specs/000-shared/data-model.md)
export const STATUS_LABELS = {
  available: 'ว่าง',
  borrowed: 'ถูกยืม',
  unavailable: 'ไม่พร้อมใช้งาน',
}

// ---- 006-return (specs/006-return/plan.md ข้อ 5 ตาราง "ค่าที่ตั้งได้") ----

// สภาพตอนรับคืน 2 ค่า ไม่มี "สูญหาย" (FR-RET-02, glossary "สภาพ (ตอนรับคืน)")
// value ตรงกับที่หลังบ้านรับ / label คือคำที่แสดงบนหน้าจอ
export const RETURN_CONDITIONS = [
  { value: 'normal', label: 'ปกติ' },
  { value: 'damaged', label: 'ชำรุด' },
]
