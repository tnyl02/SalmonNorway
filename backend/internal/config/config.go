// Package config รวมค่าที่ตั้งได้ของระบบไว้ที่เดียว (AGENTS.md กฎข้อ 9)
//
// มีทั้งค่าของสภาพแวดล้อม (port, การเชื่อมฐานข้อมูล) ซึ่งทีมเลือกเอง ไม่ได้มาจาก spec
// และค่าที่มาจาก spec ซึ่งทุกตัวต้องมีคอมเมนต์อ้าง ID กำกับ ห้ามฝังตัวเลขเหล่านี้กระจายในโค้ด
package config

import (
	"os"
	"time"
	_ "time/tzdata" // ฝังข้อมูลเขตเวลาไว้ในโปรแกรม ใช้ Asia/Bangkok ได้แม้ container ไม่มี zoneinfo
)

// TimeZone คือเขตเวลาที่ใช้เทียบ "วันที่และเวลาปัจจุบัน" (glossary, FR-EQS-04 ข้อ ข, DOM-LOAN-01)
const TimeZone = "Asia/Bangkok"

// ---------- ค่าของฟีเจอร์ 002 ส่งคำขอยืม (plan.md ข้อ 5 ตาราง "ค่าที่ตั้งได้") ----------

// LoanNormalReturnHour คือชั่วโมงที่การยืมแบบปกติต้องคืนภายในวันเดียวกัน (DOM-LOAN-01, FR-BRQ-03 ข้อ ค)
const LoanNormalReturnHour = 16

// OverdueBlockUntilNextDayStart บอกว่าผู้ขอที่เคยคืนของเกินกำหนด ส่งคำขอใหม่ได้เมื่อถึงต้นวันถัดไป
// ไม่ใช่ทันทีที่คืนและไม่ใช่ครบ 24 ชั่วโมง เทียบด้วย TimeZone (FR-BRQ-06 ข้อ ข, ทีมตัดสิน 2569-11-11)
const OverdueBlockUntilNextDayStart = true

// AttachmentMaxBytes คือขนาดสูงสุดของใบบันทึกข้อความที่แนบกับการยืมค้างคืน (FR-BRQ-10 ข้อ ก)
const AttachmentMaxBytes = 10 << 20 // 10MB

// StatusRequested คือสถานะเดียวที่ฟีเจอร์ 002 ตั้งได้ (FR-BRQ-04, data-model.md v9)
const StatusRequested = "requested"

// PurposeOther คือค่าวัตถุประสงค์ที่บังคับให้ผู้ขอพิมพ์ระบุเพิ่ม (FR-BRQ-02 ข้อ ก)
const PurposeOther = "other"

// AuditActionRequest คือค่า action ที่ใช้บันทึก audit log ของการส่งคำขอ ใช้ค่าเดียวทั้งยืมปกติและค้างคืน
// (FR-BRQ-08, DM-07, Q-BRQ-07)
const AuditActionRequest = "request"

// AttachmentAllowedTypes คือชนิดไฟล์ที่รับเป็นใบบันทึกข้อความ (FR-BRQ-10 ข้อ ก)
// spec เขียนว่า "PDF หรือรูปถ่าย" ทีมยืนยันชนิดรูปเป็น jpeg และ png เมื่อ 2569-11-11
var AttachmentAllowedTypes = []string{"application/pdf", "image/jpeg", "image/png"}

// PurposeValues คือวัตถุประสงค์การยืมทั้ง 3 ค่า (FR-BRQ-02, DM-06)
// ข้อความภาษาไทยที่แสดงบนหน้าจออยู่ที่ frontend/src/config.js
var PurposeValues = []string{"study", "activity", PurposeOther}

// RequestStatuses คือสถานะของคำขอยืมทั้ง 8 ค่าตาม data-model.md v9
// ฟีเจอร์ 002 ตั้งได้เฉพาะ StatusRequested ส่วนค่าที่เหลือเป็นของ 003 ถึง 006
var RequestStatuses = []string{
	StatusRequested,
	"approved_by_officer",
	"pending_advisor_approval",
	"approved",
	"rejected",
	"handed_over",
	"returned",
	"cancelled",
}

// ---------- ค่าของฟีเจอร์ 004 อนุมัติคำขอยืมปกติ (specs/004-approve-request/plan.md ข้อ 5, T-02) ----------
// เวลาตัด 16:00 น. ของวันที่เริ่มยืม (FR-APV-01 ค, Q-APV-11) ใช้ LoanNormalReturnHour ของ 002 ร่วม ไม่ประกาศซ้ำ
// สถานะของคำขอใช้ RequestStatuses / StatusRequested ของ 002 ร่วม

// StatusApproved และ StatusRejected คือสถานะของคำขอที่ 004 ตั้งได้ (FR-APV-01 ค ง, data-model.md v9 แถว 4, 5)
const (
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

// ApprovalEmailDeadline คือเวลาที่ต้องส่ง Email แจ้งผลให้ผู้ขอหลังบันทึกผล (FR-APV-02)
const ApprovalEmailDeadline = 60 * time.Second

// ApprovalEmailRetryCount คือจำนวนครั้งที่ส่งซ้ำอัตโนมัติ **หลังครั้งแรก** รวมครั้งแรกเป็น 4 ครั้ง (FR-APV-02 ง, Q-APV-17)
const ApprovalEmailRetryCount = 3

// ApprovalEmailRetryInterval คือระยะห่างระหว่างการส่งซ้ำแต่ละครั้ง (FR-APV-02 ง)
const ApprovalEmailRetryInterval = 1 * time.Minute

// ApprovalAutoRejectHour คือเวลาที่ปฏิเสธคำขอยืมปกติที่ยังค้างอัตโนมัติ = 00:00 น. ของวันถัดจากวันที่ส่ง
// เทียบด้วย TimeZone (FR-APV-01 ฉ, ASM-05, Q-APV-15)
const ApprovalAutoRejectHour = 0

// ApprovalSupersededReason คือเหตุผลของคำขอที่ถูกปฏิเสธตามเมื่อคำขออื่นของชิ้นเดียวกันได้รับอนุมัติ (FR-APV-01 จ)
const ApprovalSupersededReason = "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว"

// ApprovalAutoRejectReason คือเหตุผลของคำขอที่ถูกปฏิเสธอัตโนมัติตอนเที่ยงคืน (FR-APV-01 ฉ, FR-APV-02 ค)
const ApprovalAutoRejectReason = "คำขอยังคงค้างอยู่และไม่ได้รับการพิจารณาภายในเวลาที่กำหนด"

// ApprovalPickupLocation คือสถานที่รับอุปกรณ์ที่แสดงใน Email ผลอนุมัติ (FR-APV-02 ก)
const ApprovalPickupLocation = "คณะวิทยาศาสตร์ ตึก 1 ชั้น 2 ห้องอุปกรณ์ภาควิชาคอมพิวเตอร์"

// ApprovalSystemActorID คือผู้ดำเนินการของการปฏิเสธอัตโนมัติ (FR-APV-01 ฉ "ผู้ดำเนินการคือ System", plan.md ข้อ 3)
const ApprovalSystemActorID = "system"

// ค่า action ของ audit log ที่ 004 บันทึก (RQ-25, FR-APV-01 ฌ) ต้องตรงกับ CHECK ใน migration 0008
const (
	AuditActionApprove    = "approve"
	AuditActionReject     = "reject"
	AuditActionAutoReject = "auto_reject"
)

// ชนิดของผลการพิจารณาและของ Email (FR-APV-02 ก ข ค) ต้องตรงกับ CHECK ของ approval_decision / approval_email ใน migration 0008
const (
	ApprovalKindApproved     = "approved"
	ApprovalKindRejected     = "rejected"
	ApprovalKindAutoRejected = "auto_rejected"
)

// ApprovalEmailSubjects คือหัวข้อ Email แต่ละแบบ (FR-APV-02 ก ข ค)
var ApprovalEmailSubjects = map[string]string{
	ApprovalKindApproved:     "แจ้งผลการอนุมัติการยืมอุปกรณ์ภาคคอมพิวเตอร์",
	ApprovalKindRejected:     "แจ้งผลการปฏิเสธการยืมอุปกรณ์ภาคคอมพิวเตอร์",
	ApprovalKindAutoRejected: "แจ้งปฏิเสธคำขอยืมอุปกรณ์อัตโนมัติ",
}

// ApprovalEmailTemplates คือข้อมูลที่ Email แต่ละแบบต้องมี เรียงตามที่ spec เขียน (FR-APV-02 ก ข ค)
// ใช้ชื่อหัวข้อข้อมูลตามคำใน spec / ข้อความตายตัวของ Email ปฏิเสธอัตโนมัติอยู่ใน ApprovalAutoRejectEmail
// รูปแบบการจัดวางในเนื้อ Email เป็นงานของ T-14 (ทีมเลือกเอง ไม่ได้มาจาก spec)
var ApprovalEmailTemplates = map[string][]string{
	// FR-APV-02 (ก) ยืมปกติไม่แสดงจำนวนอุปกรณ์
	ApprovalKindApproved: {
		"ชื่อผู้ขอ", "รายการอุปกรณ์ที่ขอยืม", "วัตถุประสงค์", "วันเวลาที่ต้องการยืม",
		"วันเวลากำหนดคืน", "ชื่อผู้อนุมัติ", "สถานที่รับ",
	},
	// FR-APV-02 (ข) รวมคำขอที่ถูกปฏิเสธตาม FR-APV-01 (จ)
	ApprovalKindRejected: {
		"ชื่อผู้ขอ", "รายการอุปกรณ์ที่ขอยืม", "วันที่ยื่นคำขอ", "เหตุผลในการปฏิเสธ", "ชื่อผู้พิจารณา",
	},
	// FR-APV-02 (ค)
	ApprovalKindAutoRejected: {
		"เลขที่คำขอ", "ชื่อผู้ขอ", "รายการอุปกรณ์ที่ขอยืม", "วันที่ยื่นคำขอ", "สถานะ",
		"เหตุผล", "เวลาที่ดำเนินการ", "แนวทาง",
	},
}

// ApprovalAutoRejectEmail คือค่าตายตัวใน Email ปฏิเสธอัตโนมัติ (FR-APV-02 ค)
var ApprovalAutoRejectEmail = struct {
	Status   string
	DoneAt   string
	Guidance string
}{
	Status:   "ปฏิเสธอัตโนมัติ",
	DoneAt:   "00:00 น.",
	Guidance: "หากยังต้องการยืมอุปกรณ์ สามารถยื่นคำขอใหม่ได้",
}

// ApprovalMessages คือข้อความ error ภาษาไทยของ 004 ตามรหัสใน CON-TECH-03 (FR-APV-01 ซ, ASM-07, Q-APV-18)
// หน้าจอใช้ APPROVAL_MESSAGES ใน frontend/src/config.js ซึ่งต้องตรงกับค่านี้
var ApprovalMessages = map[string]string{
	"ITEM_NOT_AVAILABLE":           "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้",
	"REQUEST_CANCELLED":            "คำขอนี้ถูกผู้ขอยกเลิกแล้ว",
	"REQUEST_ALREADY_DECIDED":      "คำขอนี้ได้รับการพิจารณาแล้ว",
	"REQUEST_DECIDED_CONCURRENTLY": "คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น กรุณารีเฟรชหน้า",
	"NOT_FIRST_IN_QUEUE":           "ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน",
	"APPROVAL_CUTOFF_PASSED":       "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้",
	"VALIDATION_FAILED":            "กรุณาระบุเหตุผลในการปฏิเสธ", // ช่องเดียวที่ 004 ตรวจคือเหตุผลการปฏิเสธ (FR-APV-01 ง)
	"FORBIDDEN_ROLE":               "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้",
	"FORBIDDEN_DEPARTMENT":         "คำขอนี้เป็นของอุปกรณ์ภาควิชาอื่น คุณไม่มีสิทธิ์พิจารณา",
}

// ---------- จบค่าของฟีเจอร์ 004 ----------

// Config คือค่าที่โปรแกรมหลังบ้านใช้ตอนเริ่มทำงาน
type Config struct {
	// HTTPAddr คือที่อยู่ที่หลังบ้านรอรับ request
	// ทีมเลือกเอง: port 8000 ตรงกับ proxy ใน frontend/vite.config.js ของ template รายวิชา
	HTTPAddr string

	// DatabaseURL คือที่อยู่ PostgreSQL (CON-TECH-01)
	// ค่าเริ่มต้นตรงกับ deploy/docker-compose.yml สำหรับเครื่องพัฒนาเท่านั้น
	DatabaseURL string

	// AttachmentDir คือโฟลเดอร์ที่เก็บไฟล์ใบบันทึกข้อความของการยืมค้างคืน (ASM-BRQ-07)
	// เก็บเป็นไฟล์บนเครื่องที่รันระบบ ไม่ได้เก็บในฐานข้อมูล
	// ค่าเริ่มต้นเป็นของเครื่องพัฒนา ไม่ใช่ requirement เปลี่ยนด้วย ATTACHMENT_DIR
	AttachmentDir string
}

// Load อ่านค่าจาก environment ถ้าไม่ตั้งไว้ใช้ค่าเริ่มต้นสำหรับเครื่องพัฒนา
func Load() Config {
	return Config{
		HTTPAddr:      getenv("HTTP_ADDR", ":8000"),
		DatabaseURL:   getenv("DATABASE_URL", "postgres://salmon:salmon@localhost:5432/salmon?sslmode=disable"),
		AttachmentDir: getenv("ATTACHMENT_DIR", "./var/attachments"),
	}
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// ---- 005-handover (plan.md ข้อ 5 ของ 005, T-01) ----

// BorrowerAuthTimeout คือเวลาที่รอ Auth API ยืนยันตัวตนผู้ยืม เกินนี้ถือว่าเชื่อมต่อไม่สำเร็จ
// แล้วใช้ผลการตรวจบัตรแทน (FR-HND-02, ASM-HND-02, IF-IDP-01)
const BorrowerAuthTimeout = 10 * time.Second

// HandoverMinPhotos คือจำนวนรูปสภาพอุปกรณ์ขั้นต่ำที่ต้องแนบตอนส่งมอบ (FR-HND-06, ASM-HND-07)
const HandoverMinPhotos = 2

// HandoverPhotoMaxBytes คือขนาดไฟล์รูปสูงสุด 10 MB (FR-HND-06, ASM-HND-07)
// นับต่อรูปหรือรวมทุกรูป รอ Q-HND-10 (ใช้ใน T-28)
const HandoverPhotoMaxBytes = 10 << 20

// HandoverPhotoTypes คือชนิดไฟล์รูปที่รับ JPG และ PNG (FR-HND-06, ASM-HND-07)
// คืน slice ใหม่ทุกครั้ง กันโค้ดที่เรียกใช้เผลอแก้ค่ากลาง
func HandoverPhotoTypes() []string {
	return []string{"image/jpeg", "image/png"}
}

// PhotoDir คือโฟลเดอร์รากของไฟล์รูปบน Docker volume `photos` (plan.md ของ 005 Q-P1: ทีมเลือกเอง ไม่ได้มาจาก spec)
// 005 ใช้โฟลเดอร์ย่อย handover/ ฟีเจอร์อื่น (เช่น 006) ใช้โฟลเดอร์ข้างกันได้ เปลี่ยนได้ด้วย env PHOTO_DIR
func PhotoDir() string {
	return getenv("PHOTO_DIR", "/data/photos")
}

// HandoverReadyStatuses คือสถานะของคำขอยืมที่ส่งมอบได้ = `approved` (FR-HND-01, data-model.md v9, Q-HND-05 ปิดแล้ว)
// คืน slice ใหม่ทุกครั้ง กันโค้ดที่เรียกใช้เผลอแก้ค่ากลาง
func HandoverReadyStatuses() []string {
	return []string{"approved"}
}
