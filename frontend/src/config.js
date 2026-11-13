// ค่าที่ตั้งได้ของหน้าจอรวมไว้ที่เดียว (AGENTS.md กฎข้อ 9)

// ---- 005-handover (plan.md ข้อ 5 ของ 005, T-01) ----
// วางไว้ต้นไฟล์ ไม่ต่อท้าย เพื่อไม่ชนกับฟีเจอร์อื่นที่เพิ่มค่าต่อท้ายไฟล์นี้ตอน merge
// ใช้ตรวจรูปที่หน้าจอก่อนส่งเพื่อความสะดวกเท่านั้น หลังบ้านตรวจซ้ำเสมอ (backend/internal/config/config.go)

// จำนวนรูปสภาพอุปกรณ์ขั้นต่ำตอนส่งมอบ (FR-HND-06, ASM-HND-07)
export const HANDOVER_MIN_PHOTOS = 2

// ขนาดไฟล์รูปสูงสุด 10 MB ต่อรูป ไม่จำกัดจำนวนรูป (FR-HND-06, ASM-HND-07, Q-HND-10 ปิด 2569-11-12)
export const HANDOVER_PHOTO_MAX_BYTES = 10 * 1024 * 1024

// ชนิดไฟล์รูปที่รับ JPG และ PNG (FR-HND-06, ASM-HND-07)
export const HANDOVER_PHOTO_TYPES = ['image/jpeg', 'image/png']
// ---- จบ 005-handover ----

// ---- 004-approve-request (specs/004-approve-request/plan.md ข้อ 5, T-02) ----
// วางไว้ต้นไฟล์ต่อจาก 005 เพื่อไม่ชนกับฟีเจอร์อื่นที่เพิ่มค่าต่อท้ายไฟล์ตอน merge

// ข้อความ error ภาษาไทยของหน้าจออนุมัติ เลือกตามรหัสใน CON-TECH-03 ไม่ใช่ตามข้อความ (FR-APV-01 ซ, ASM-07, Q-APV-18)
// ต้องตรงกับ ApprovalMessages ใน backend/internal/config/config.go
export const APPROVAL_MESSAGES = {
  ITEM_NOT_AVAILABLE: 'อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้',
  REQUEST_CANCELLED: 'คำขอนี้ถูกผู้ขอยกเลิกแล้ว',
  REQUEST_ALREADY_DECIDED: 'คำขอนี้ได้รับการพิจารณาแล้ว',
  REQUEST_DECIDED_CONCURRENTLY: 'คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น กรุณารีเฟรชหน้า',
  NOT_FIRST_IN_QUEUE: 'ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน',
  APPROVAL_CUTOFF_PASSED: 'เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้',
  VALIDATION_FAILED: 'กรุณาระบุเหตุผลในการปฏิเสธ',
  FORBIDDEN_ROLE: 'บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้',
  FORBIDDEN_DEPARTMENT: 'คำขอนี้เป็นของอุปกรณ์ภาควิชาอื่น คุณไม่มีสิทธิ์พิจารณา',
  UPSTREAM_UNAVAILABLE: 'ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง',
  REQUEST_NOT_FOUND: 'ไม่พบคำขอนี้',
}
// เหตุผลการปฏิเสธยาวไม่เกินเท่านี้ นับตาม code point แบบเดียวกับ 003 (FR-APV-01 ง, Q-APV-28)
// ต้องตรงกับ ApprovalRejectReasonMaxLength ใน backend/internal/config/config.go
export const APPROVAL_REASON_MAX_LENGTH = 500
// ข้อความเมื่อเหตุผลยาวเกิน (FR-APV-01 ซ, Q-APV-28) ต้องตรงกับ ApprovalReasonTooLongMessage ของหลังบ้าน
export const APPROVAL_REASON_TOO_LONG = 'เหตุผลต้องไม่เกิน 500 ตัวอักษร'
// ชื่ออาจารย์ที่เซ็นอนุมัติคำขอค้างคืน ยาวไม่เกินเท่านี้ นับตาม code point (FR-APV-03 ช, Q-APV-44)
// ต้องตรงกับ ApprovalAdvisorNameMaxLength ใน backend/internal/config/config.go
export const APPROVAL_ADVISOR_NAME_MAX_LENGTH = 100
// ข้อความเมื่อไม่ระบุชื่ออาจารย์ (FR-APV-01 ซ, Q-APV-44) ต้องตรงกับ ApprovalAdvisorNameRequiredMessage ของหลังบ้าน
export const APPROVAL_ADVISOR_NAME_REQUIRED = 'กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ'
// ข้อความเมื่อชื่ออาจารย์ยาวเกิน (FR-APV-01 ซ, Q-APV-44) ต้องตรงกับ ApprovalAdvisorNameTooLongMessage ของหลังบ้าน
export const APPROVAL_ADVISOR_NAME_TOO_LONG = 'ชื่ออาจารย์ต้องไม่เกิน 100 ตัวอักษร'
// ข้อความเมื่อไม่มีคำขอที่รอพิจารณาในหน้ารายการ (FR-APV-01 ก, Q-APV-31) ข้อความเดียวกับ MSG_NO_PENDING ของ 003
export const APPROVAL_EMPTY_LIST = 'ไม่มีคำขอที่รออนุมัติ'
// แบบของ Email ที่แสดงในรายการ "Email ที่ส่งไม่สำเร็จ" (FR-APV-02 ก ข ค, T-34) ค่า key ตรงกับ kind ของหลังบ้าน
// ข้อความทีมออกแบบเอง ไม่ได้มาจาก spec
export const APPROVAL_EMAIL_KIND_LABELS = {
  approved: 'ผลอนุมัติ',
  rejected: 'ผลปฏิเสธ',
  auto_rejected: 'ปฏิเสธอัตโนมัติ',
}
// ---- จบ 004-approve-request ----

// ---- 006-return (specs/006-return/plan.md ข้อ 5 ตาราง "ค่าที่ตั้งได้") ----
// วางไว้ต้นไฟล์ต่อจาก 004 เพื่อไม่ชนกับฟีเจอร์อื่นที่เพิ่มค่าต่อท้ายไฟล์ตอน merge / เขตเวลาใช้ TIME_ZONE ร่วมกับฟีเจอร์อื่น

// ชื่อสถานะอุปกรณ์ที่แสดงบนหน้าจอ 3 ค่า (CON-DATA-01, specs/000-shared/data-model.md)
export const STATUS_LABELS = {
  available: 'ว่าง',
  borrowed: 'ถูกยืม',
  unavailable: 'ไม่พร้อมใช้งาน',
}

// สภาพตอนรับคืน 2 ค่า ไม่มี "สูญหาย" (FR-RET-02, glossary "สภาพ (ตอนรับคืน)")
// value ตรงกับ ReturnConditionStatus ใน backend/internal/config/config.go / label คือคำที่แสดงบนหน้าจอ
export const RETURN_CONDITIONS = [
  { value: 'normal', label: 'ปกติ' },
  { value: 'damaged', label: 'ชำรุด' },
]
// ---- จบ 006-return ----

// token ของผู้ใช้จำลองที่หน้าจอส่งไปกับทุก request (IF-IDP-01 ใช้ระบบยืนยันตัวตนจำลองไปก่อน รอ Q-SH-03)
// ทีมเลือกเอง ไม่ได้มาจาก spec (ตัดสิน 2569-11-10 ตอนทำ T-06) ใช้จนกว่าจะมีหน้าล็อกอินจำลองใน T-21
// เปลี่ยนผู้ใช้ได้ด้วย VITE_MOCK_TOKEN เช่น `VITE_MOCK_TOKEN=mock-staff-001 npm run dev`
// token ที่ใช้ได้ดูใน backend/internal/auth/mock_provider.go
export const MOCK_TOKEN = import.meta.env.VITE_MOCK_TOKEN ?? 'mock-student-001'

// เขตเวลาที่ใช้แสดงวันที่และเวลา (glossary: เวลาปัจจุบัน / วันเดียวกัน เทียบด้วย Asia/Bangkok, DOM-LOAN-01)
// 003 ใช้ค่าเดียวกันแสดงเวลาที่ส่งคำขอ (FR-CXL-06)
export const TIME_ZONE = 'Asia/Bangkok'

// ---- 001-equipment-status ----

// เวลาที่หน้าจอรอหลังบ้านของ 001 (รายการรุ่น รายชิ้น รายละเอียด) เกินนี้ถือว่าดึงข้อมูลไม่สำเร็จ แสดงข้อผิดพลาดพร้อมปุ่มลองใหม่
// (FR-EQS-08, AC-EQS-08, Q-08 ทีมตัดสิน 2569-11-12) ต้องมากกว่า EquipmentDBTimeout (5 วินาที) ใน backend/internal/config/config.go
// เพื่อให้ได้ 503 พร้อม message จากหลังบ้านก่อน ใช้เฉพาะ request ของ 001 ไม่ใช้กับ 002, 003 (ทีมตัดสิน 2569-11-12 ตอนทำ T-26)
export const EQUIPMENT_REQUEST_TIMEOUT_MS = 6000

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

// ---- ข้อความบนหน้าจอของ 002 (T-19) ----
// Q-BRQ-14 ปิดแล้ว 2569-11-12: spec ไม่มีหมวดหน้าจอและไม่มี mockup เจ้าของฟีเจอร์เป็นผู้กำหนดข้อความ
// ข้อความที่ FR สั่งให้แสดงรวมไว้ที่นี่ที่เดียว เพื่อให้ตามรอยกลับไปหา FR ได้และ test อ้างค่าเดียวกัน

// เงื่อนไขของการยืมแบบปกติ ต้องแสดงก่อนส่ง (FR-BRQ-10 ข้อ ข, AC-BRQ-10, DOM-LOAN-01)
export const MSG_LOAN_NORMAL_CONDITIONS = `ยืมแบบปกติต้องคืนภายใน ${LOAN_NORMAL_RETURN_HOUR}:00 น. ของวันเดียวกัน และใช้งานภายในพื้นที่มหาวิทยาลัย`

// เงื่อนไขไฟล์แนบของการยืมค้างคืน ต้องบอกชนิดและขนาดที่รับได้ (FR-BRQ-10 ข้อ ก, AC-BRQ-10)
export const MSG_ATTACHMENT_HINT = `แนบใบบันทึกข้อความ: PDF หรือรูปถ่าย ขนาดไม่เกิน ${ATTACHMENT_MAX_BYTES / (1024 * 1024)} MB`

// ผลของการส่งคำขอ (FR-BRQ-04 สถานะเริ่มต้นคือรออนุมัติ)
export const MSG_SUBMIT_SUCCESS = (requestId) => `ส่งคำขอยืมเรียบร้อยแล้ว หมายเลขคำขอ ${requestId} สถานะ รออนุมัติ`

// ส่งไม่สำเร็จเพราะระบบไม่ตอบสนอง ใช้เมื่อหลังบ้านไม่ได้ส่งข้อความมา (FR-BRQ-09) รูปแบบเดียวกับ MSG_CANCEL_FAILED ของ 003
export const MSG_SUBMIT_FAILED = 'ส่งคำขอไม่สำเร็จ กรุณาลองใหม่อีกครั้ง'

// ป้ายปุ่มส่งคำขอ ใช้คำเดียวกับปุ่มของ 001 ที่พาผู้ใช้มาหน้านี้ (FR-EQS-06, T-34 ของ 001)
export const LABEL_SUBMIT_REQUEST = 'ส่งคำขอยืม'

// ---- ข้อความของแต่ละรหัส error บนหน้าฟอร์มคำขอ (T-20) ----
// หน้าจอตัดสินใจจาก "รหัส" ไม่ใช่จากข้อความที่หลังบ้านส่งมา (CON-TECH-03)
// ครบทุกรหัสที่ฟีเจอร์นี้ใช้ตามตาราง Error contract ใน plan.md ข้อ 4
// ข้อความตรงกับที่หลังบ้านตอบใน backend/internal/borrowrequest/handler.go เพื่อไม่ให้ผู้ใช้เห็นคนละคำกัน
export const BORROW_REQUEST_MESSAGES = {
  UNAUTHENTICATED: 'กรุณาเข้าสู่ระบบก่อนส่งคำขอ', // 401 IF-IDP-01
  PRINCIPAL_INCOMPLETE: 'บัญชีของคุณยังไม่มีข้อมูลภาควิชา กรุณาติดต่อเจ้าหน้าที่', // 403 IF-IDP-01
  FORBIDDEN_DEPARTMENT: 'ส่งคำขอยืมได้เฉพาะอุปกรณ์ของภาควิชาตนเอง', // 403 FR-BRQ-07
  ITEM_NOT_AVAILABLE: 'อุปกรณ์ชิ้นนี้ไม่อยู่ในสถานะว่างแล้ว', // 409 FR-BRQ-01 (แสดงสถานะล่าสุดต่อท้าย)
  DUPLICATE_REQUEST: 'คุณมีคำขอของอุปกรณ์ชิ้นนี้ค้างอยู่แล้ว', // 409 FR-BRQ-05 ข
  HAS_OVERDUE_LOAN: 'คุณมีอุปกรณ์ที่คืนเกินกำหนดค้างอยู่ กรุณาคืนอุปกรณ์ก่อนส่งคำขอใหม่', // 409 FR-BRQ-06
  VALIDATION_FAILED: 'กรุณาแก้ไขข้อมูลให้ครบถ้วนก่อนส่งคำขอ', // 422 FR-BRQ-03, FR-BRQ-10
  UPSTREAM_UNAVAILABLE: MSG_SUBMIT_FAILED, // 503 FR-BRQ-09 ใช้ข้อความเดียวกับตอนส่งไม่สำเร็จ
}

// ชื่อช่องภาษาไทยของแต่ละค่าที่หลังบ้านส่งกลับมาใน fields ของ VALIDATION_FAILED (FR-BRQ-03, AC-BRQ-03)
// ชื่อ key ตรงกับ CreateRequestDTO และ validate.go ของหลังบ้าน
export const FIELD_LABELS = {
  requesterName: 'ชื่อ-นามสกุล',
  studentCode: 'รหัสนักศึกษา',
  startAt: 'วันและเวลาที่เริ่มยืม',
  endAt: 'วันและเวลาที่คืน',
  purpose: 'วัตถุประสงค์การยืม',
  purposeOther: 'ระบุวัตถุประสงค์',
  attachment: 'ใบบันทึกข้อความ',
}

// นำหน้ารายชื่อช่องที่ต้องแก้ (FR-BRQ-03 "ระบุช่องที่ยังขาด", AC-BRQ-03, AC-BRQ-10)
export const MSG_FIELDS_TO_FIX = 'ช่องที่ต้องแก้:'

// แจ้งผู้ขอเรื่องคำขอที่ส่งซ้อนกับของผู้อื่น แสดงคู่กับข้อความสำเร็จทุกครั้ง (FR-BRQ-05, T-22)
// ทีมตัดสิน 2569-11-13 ให้เป็นข้อความคงที่ ไม่บอกว่ามีคำขออื่นอยู่จริงหรือไม่ และไม่บอกลำดับ
// เพราะ AC-BRQ-05 ระบุว่าหน้าจอของ 002 ไม่แสดงข้อมูลคำขอของผู้อื่นให้นักศึกษาเห็น (DOM-PRIV-01)
// ส่วน "จะได้รับแจ้งเมื่อมีผล" มาจากคำตอบของ Q-SH-06 ที่ว่าต้องแจ้งผู้ขอที่ถูกปฏิเสธ (RQ-21 การแจ้งเป็นงานของ 004)
export const MSG_STACKED_NOTICE =
  'อุปกรณ์ชิ้นนี้อาจมีผู้อื่นขอยืมด้วย เจ้าหน้าที่จะพิจารณาตามลำดับการส่งคำขอ และคุณจะได้รับแจ้งเมื่อมีผล'
