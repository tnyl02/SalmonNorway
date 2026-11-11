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
export const TIME_ZONE = 'Asia/Bangkok'
