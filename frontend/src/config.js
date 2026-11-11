// ค่าที่หน้าจอใช้ซึ่งมาจาก spec รวมไว้ที่เดียว (AGENTS.md กฎข้อ 9)
// ต้องตรงกับ backend/internal/config/config.go เสมอ หลังบ้านเป็นผู้ตรวจจริง หน้าจอตรวจซ้ำเพื่อแจ้งผู้ใช้ก่อนส่ง

// ---- 006-return (specs/006-return/plan.md ข้อ 5 ตาราง "ค่าที่ตั้งได้") ----

// จำนวนรูปน้อยสุดที่ต้องแนบทุกครั้งที่บันทึกการรับคืน (FR-RET-05, ASM-06)
export const RETURN_PHOTO_MIN = 1

// จำนวนรูปมากสุดต่อการรับคืน 1 ครั้ง (FR-RET-05, ASM-06)
export const RETURN_PHOTO_MAX = 5

// ขนาดสูงสุดต่อรูป "5 MB" = 5 × 1024 × 1024 ไบต์ (FR-RET-05, ASM-06, ทีมเลือกการตีความเอง plan.md Q-P2)
export const RETURN_PHOTO_MAX_BYTES = 5 * 1024 * 1024

// ชนิดรูปที่รับ JPG หรือ PNG (FR-RET-05, ASM-06)
export const RETURN_PHOTO_TYPES = ['image/jpeg', 'image/png']

// สภาพตอนรับคืน 2 ค่า ไม่มี "สูญหาย" (FR-RET-02, glossary "สภาพ (ตอนรับคืน)")
// value ตรงกับที่หลังบ้านรับ / label คือคำที่แสดงบนหน้าจอ
export const RETURN_CONDITIONS = [
  { value: 'normal', label: 'ปกติ' },
  { value: 'damaged', label: 'ชำรุด' },
]
