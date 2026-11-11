// Package config รวมค่าที่ตั้งได้ของระบบไว้ที่เดียว (AGENTS.md กฎข้อ 9)
//
// ตอนนี้มีเฉพาะค่าของสภาพแวดล้อม (port, การเชื่อมฐานข้อมูล) ซึ่งทีมเลือกเอง ไม่ได้มาจาก spec
// ค่าที่มาจาก spec (เช่น เขตเวลา Asia/Bangkok, จำนวนต่อหน้า) ให้ task ที่ใช้เพิ่มที่นี่ พร้อมคอมเมนต์อ้าง ID
package config

import (
	"os"
	"time"
	_ "time/tzdata" // ฝังข้อมูลเขตเวลาไว้ในโปรแกรม ใช้ Asia/Bangkok ได้แม้ container ไม่มี zoneinfo
)

// TimeZone คือเขตเวลาที่ใช้เทียบ "วันที่และเวลาปัจจุบัน" (glossary, FR-EQS-04 ข้อ ข, DOM-LOAN-01)
const TimeZone = "Asia/Bangkok"

// Config คือค่าที่โปรแกรมหลังบ้านใช้ตอนเริ่มทำงาน
type Config struct {
	// HTTPAddr คือที่อยู่ที่หลังบ้านรอรับ request
	// ทีมเลือกเอง: port 8000 ตรงกับ proxy ใน frontend/vite.config.js ของ template รายวิชา
	HTTPAddr string

	// DatabaseURL คือที่อยู่ PostgreSQL (CON-TECH-01)
	// ค่าเริ่มต้นตรงกับ deploy/docker-compose.yml สำหรับเครื่องพัฒนาเท่านั้น
	DatabaseURL string
}

// Load อ่านค่าจาก environment ถ้าไม่ตั้งไว้ใช้ค่าเริ่มต้นสำหรับเครื่องพัฒนา
func Load() Config {
	return Config{
		HTTPAddr:    getenv("HTTP_ADDR", ":8000"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://salmon:salmon@localhost:5432/salmon?sslmode=disable"),
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

// HandoverDueHour คือชั่วโมงของกำหนดคืนการยืมแบบปกติ = 16:00 น. ที่ 005 ใช้ตอนส่งมอบ (DOM-LOAN-01, FR-HND-05) เทียบด้วยเขตเวลา TimeZone
// ค่าเดียวกับ LoanNormalReturnHour ของ 002 ที่ยังไม่ merge เข้า main ใช้ชื่อแยกชั่วคราวเพื่อไม่ให้ประกาศชื่อซ้ำตอน merge
// เมื่อ 002 merge แล้ว ให้ 005 เปลี่ยนไปใช้ LoanNormalReturnHour แล้วลบค่านี้ (AGENTS.md กฎข้อ 9 ค่าเดียวที่เดียว)
const HandoverDueHour = 16

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
