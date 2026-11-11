// ค่าที่ตั้งได้ของหน้าจอรวมไว้ที่เดียว (AGENTS.md กฎข้อ 9)

// token ของผู้ใช้จำลองที่หน้าจอส่งไปกับทุก request (IF-IDP-01 ใช้ระบบยืนยันตัวตนจำลองไปก่อน รอ Q-SH-03)
// ทีมเลือกเอง ไม่ได้มาจาก spec (ตัดสิน 2569-11-10 ตอนทำ T-06) ใช้จนกว่าจะมีหน้าล็อกอินจำลองใน T-21
// เปลี่ยนผู้ใช้ได้ด้วย VITE_MOCK_TOKEN เช่น `VITE_MOCK_TOKEN=mock-staff-001 npm run dev`
// token ที่ใช้ได้ดูใน backend/internal/auth/mock_provider.go
export const MOCK_TOKEN = import.meta.env.VITE_MOCK_TOKEN ?? 'mock-student-001'

// เขตเวลาที่ใช้แสดงวันที่และเวลา (glossary: เวลาปัจจุบัน / วันเดียวกัน เทียบด้วย Asia/Bangkok, DOM-LOAN-01)
// 003 ใช้ค่าเดียวกันแสดงเวลาที่ส่งคำขอ (FR-CXL-06)
export const TIME_ZONE = 'Asia/Bangkok'

// ---- 003-cancel-request (plan.md ข้อ 5 ของ 003, T-03) ----

// ความยาวสูงสุดของเหตุผลการยกเลิก นับตาม code point เหมือนหลังบ้าน (FR-CXL-08, Q-14, ASM-13)
export const CANCEL_REASON_MAX_LEN = 200

// ข้อความเมื่อไม่มีคำขอที่รออนุมัติ (FR-CXL-06)
export const MSG_NO_PENDING = 'ไม่มีคำขอที่รออนุมัติ'

// ข้อความเมื่อโหลดรายการไม่สำเร็จ และป้ายปุ่มโหลดใหม่ (FR-CXL-06, ASM-13)
export const MSG_LOAD_FAILED = 'โหลดรายการไม่สำเร็จ กรุณาลองใหม่อีกครั้ง'
export const LABEL_RETRY = 'ลองใหม่'

// ข้อความเมื่อยกเลิกสำเร็จ (FR-CXL-01, ASM-13)
export const MSG_CANCEL_SUCCESS = (modelName, assetCode) => `ยกเลิกคำขอยืม ${modelName} (${assetCode}) เรียบร้อยแล้ว`

// ข้อความเมื่อยกเลิกไม่สำเร็จเพราะระบบไม่ตอบสนอง (503) ลองใหม่ได้ (FR-CXL-01, ASM-06, AC-CXL-13)
// คำทีมเลือก 2569-11-12 ตอนทำ T-24 รูปแบบเดียวกับ MSG_LOAD_FAILED
export const MSG_CANCEL_FAILED = 'ยกเลิกไม่สำเร็จ กรุณาลองใหม่อีกครั้ง'

// ข้อความเมื่อคำขอได้รับการอนุมัติแล้ว (FR-CXL-02 คำตามที่เขียนใน spec ทีมยืนยัน 2569-11-11)
export const MSG_ALREADY_APPROVED = 'ยกเลิกไม่ได้ เพราะคำขอได้รับการอนุมัติแล้ว'

// ข้อความเมื่อคำขอ Rejected, Handed Over, Returned หรือยกเลิกไปแล้ว (FR-CXL-02, Q-03)
export const MSG_NOT_CANCELLABLE = 'คำขอนี้ยกเลิกไม่ได้'

// ข้อความเมื่อไม่มีคำขอนั้น หรือเป็นคำขอของผู้อื่น (FR-CXL-03, Q-03)
export const MSG_REQUEST_NOT_FOUND = 'ไม่พบคำขอ'
