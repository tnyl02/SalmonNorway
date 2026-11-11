// ค่าที่ตั้งได้ของหน้าจอรวมไว้ที่เดียว (AGENTS.md กฎข้อ 9)

// ---- 005-handover (plan.md ข้อ 5 ของ 005, T-01) ----
// วางไว้ต้นไฟล์ ไม่ต่อท้าย เพื่อไม่ชนกับฟีเจอร์อื่นที่เพิ่มค่าต่อท้ายไฟล์นี้ตอน merge
// ใช้ตรวจรูปที่หน้าจอก่อนส่งเพื่อความสะดวกเท่านั้น หลังบ้านตรวจซ้ำเสมอ (backend/internal/config/config.go)

// จำนวนรูปสภาพอุปกรณ์ขั้นต่ำตอนส่งมอบ (FR-HND-06, ASM-HND-07)
export const HANDOVER_MIN_PHOTOS = 2

// ขนาดไฟล์รูปสูงสุด 10 MB (FR-HND-06, ASM-HND-07) ต่อรูปหรือรวม รอ Q-HND-10
export const HANDOVER_PHOTO_MAX_BYTES = 10 * 1024 * 1024

// ชนิดไฟล์รูปที่รับ JPG และ PNG (FR-HND-06, ASM-HND-07)
export const HANDOVER_PHOTO_TYPES = ['image/jpeg', 'image/png']
// ---- จบ 005-handover ----

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

// ---------- ค่าของฟีเจอร์ 002 ส่งคำขอยืม (plan.md ข้อ 5 ตาราง "ค่าที่ตั้งได้") ----------

// ชั่วโมงที่การยืมแบบปกติต้องคืนภายในวันเดียวกัน ใช้แสดงเงื่อนไขและเตือนก่อนส่ง (DOM-LOAN-01, FR-BRQ-03 ข้อ ค, FR-BRQ-10 ข้อ ข)
export const LOAN_NORMAL_RETURN_HOUR = 16

// ขนาดสูงสุดของใบบันทึกข้อความที่แนบกับการยืมค้างคืน (FR-BRQ-10 ข้อ ก) ต้องตรงกับ AttachmentMaxBytes ใน backend
export const ATTACHMENT_MAX_BYTES = 10 * 1024 * 1024

// ชนิดไฟล์ที่รับเป็นใบบันทึกข้อความ (FR-BRQ-10 ข้อ ก) spec เขียนว่า "PDF หรือรูปถ่าย" ทีมยืนยันชนิดรูปเมื่อ 2569-11-11
export const ATTACHMENT_ALLOWED_TYPES = ['application/pdf', 'image/jpeg', 'image/png']

// วัตถุประสงค์การยืม 3 ค่า พร้อมข้อความที่แสดงบนหน้าจอ (FR-BRQ-02, DM-06)
// ค่า value ต้องตรงกับ PurposeValues ใน backend/internal/config/config.go
export const PURPOSE_OPTIONS = [
  { value: 'study', label: 'เพื่อการเรียน' },
  { value: 'activity', label: 'เพื่อทำกิจกรรม' },
  { value: 'other', label: 'อื่น ๆ' },
]

// ค่าวัตถุประสงค์ที่บังคับให้ผู้ขอพิมพ์ระบุเพิ่ม (FR-BRQ-02 ข้อ ก)
export const PURPOSE_OTHER = 'other'
