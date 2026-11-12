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
// 003 ใช้ค่าเดียวกันแสดงเวลาที่ส่งคำขอ (FR-CXL-06)
const TimeZone = "Asia/Bangkok"

// ---- 003-cancel-request (plan.md ข้อ 5 ของ 003, T-03) ----

// CancellableStatuses คือสถานะของคำขอที่นักศึกษายังยกเลิกได้ (CON-CXL-02, FR-CXL-01, ASM-04)
// คืน slice ใหม่ทุกครั้ง กันโค้ดที่เรียกใช้เผลอแก้ค่ากลาง
func CancellableStatuses() []string {
	return []string{"requested", "approved_by_officer", "pending_advisor_approval"}
}

// CancelReasonMaxLen คือความยาวสูงสุดของเหตุผลการยกเลิก นับตาม code point (FR-CXL-08, Q-14, ASM-13)
const CancelReasonMaxLen = 200

// CancelPerfP95 คือเวลาตอบสนอง p95 สูงสุดของการยกเลิกและการเปิดหน้ารายการ (NFR-PERF-CXL-01) ใช้เป็น threshold ของ load test
const CancelPerfP95 = 15 * time.Second

// CancelDBTimeout คือเวลารอฐานข้อมูลตอนยกเลิกก่อนตอบ 503 (ASM-06)
// ตัวเลขทีมเลือกเอง ไม่ได้มาจาก spec (ตัดสิน 2569-11-11 ตอนทำ T-03) ต้องน้อยกว่า CancelPerfP95
const CancelDBTimeout = 5 * time.Second

// NotifyTimeout คือเวลารอส่ง Email แจ้งเจ้าหน้าที่หลังยกเลิก (FR-CXL-09, ASM-10, NFR-PERF-CXL-01 ไม่นับเวลาส่ง Email)
// ส่งไม่ทันนับเป็นส่งไม่สำเร็จ การยกเลิกยังสำเร็จ ตัวเลขทีมเลือกเอง ไม่ได้มาจาก spec (rtm F-04 ตัดสิน 2569-11-12, T-34)
const NotifyTimeout = 5 * time.Second

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
