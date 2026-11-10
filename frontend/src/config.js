// ค่าที่ตั้งได้ของหน้าจอรวมไว้ที่เดียว (AGENTS.md กฎข้อ 9)

// token ของผู้ใช้จำลองที่หน้าจอส่งไปกับทุก request (IF-IDP-01 ใช้ระบบยืนยันตัวตนจำลองไปก่อน รอ Q-SH-03)
// ทีมเลือกเอง ไม่ได้มาจาก spec (ตัดสิน 2569-11-10 ตอนทำ T-06) ใช้จนกว่าจะมีหน้าล็อกอินจำลองใน T-21
// เปลี่ยนผู้ใช้ได้ด้วย VITE_MOCK_TOKEN เช่น `VITE_MOCK_TOKEN=mock-staff-001 npm run dev`
// token ที่ใช้ได้ดูใน backend/internal/auth/mock_provider.go
export const MOCK_TOKEN = import.meta.env.VITE_MOCK_TOKEN ?? 'mock-student-001'

// เขตเวลาที่ใช้แสดงวันที่และเวลา (glossary: เวลาปัจจุบัน / วันเดียวกัน เทียบด้วย Asia/Bangkok, DOM-LOAN-01)
export const TIME_ZONE = 'Asia/Bangkok'
