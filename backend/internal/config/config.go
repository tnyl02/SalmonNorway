// Package config รวมค่าที่ตั้งได้ของระบบไว้ที่เดียว (AGENTS.md กฎข้อ 9)
//
// มีทั้งค่าของสภาพแวดล้อม (port, การเชื่อมฐานข้อมูล) ซึ่งทีมเลือกเอง ไม่ได้มาจาก spec
// และค่าที่มาจาก spec ซึ่งทุกตัวต้องมีคอมเมนต์อ้าง ID กำกับ ห้ามฝังตัวเลขเหล่านี้กระจายในโค้ด
package config

import (
	"os"
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
