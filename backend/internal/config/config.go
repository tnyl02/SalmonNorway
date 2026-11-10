// Package config รวมค่าที่ตั้งได้ของระบบไว้ที่เดียว (AGENTS.md กฎข้อ 9)
//
// ตอนนี้มีเฉพาะค่าของสภาพแวดล้อม (port, การเชื่อมฐานข้อมูล) ซึ่งทีมเลือกเอง ไม่ได้มาจาก spec
// ค่าที่มาจาก spec (เช่น เขตเวลา Asia/Bangkok, จำนวนต่อหน้า) ให้ task ที่ใช้เพิ่มที่นี่ พร้อมคอมเมนต์อ้าง ID
package config

import "os"

// Config คือค่าที่โปรแกรมหลังบ้านใช้ตอนเริ่มทำงาน
type Config struct {
	// HTTPAddr คือที่อยู่ที่หลังบ้านรอรับ request
	// ทีมเลือกเอง: port 8000 ตรงกับ proxy ใน frontend/vite.config.js ของ template รายวิชา
	HTTPAddr string

	// DatabaseURL คือที่อยู่ PostgreSQL (CON-TECH-01)
	// ค่าเริ่มต้นตรงกับ deploy/docker-compose.yml สำหรับเครื่องพัฒนาเท่านั้น
	DatabaseURL string
}

// ---- ค่าที่มาจาก spec ของ 006-return (specs/006-return/plan.md ข้อ 5 ตาราง "ค่าที่ตั้งได้") ----

// Timezone คือเขตเวลาที่ใช้เทียบ "วันที่ปัจจุบัน / เวลาปัจจุบัน / วันเดียวกัน" (specs/000-shared/glossary.md)
const Timezone = "Asia/Bangkok"

const (
	// ReturnPhotoMin คือจำนวนรูปน้อยสุดที่ต้องแนบทุกครั้งที่บันทึกการรับคืน (FR-RET-05, ASM-06)
	ReturnPhotoMin = 1
	// ReturnPhotoMax คือจำนวนรูปมากสุดต่อการรับคืน 1 ครั้ง (FR-RET-05, ASM-06)
	ReturnPhotoMax = 5
	// ReturnPhotoMaxBytes คือขนาดสูงสุดต่อรูป "5 MB" (FR-RET-05, ASM-06)
	// ตีความเป็น 5 × 1024 × 1024 ไบต์ ทีมเลือกเอง ไม่ได้มาจาก spec (plan.md Q-P2)
	ReturnPhotoMaxBytes = 5 * 1024 * 1024
	// returnFormOverheadBytes คือที่เผื่อให้ฟิลด์อื่นของ form (assetCode, condition, หัว multipart)
	// 1 MB ทีมเลือกเอง ไม่ได้มาจาก spec (plan.md ข้อ 5)
	returnFormOverheadBytes = 1 * 1024 * 1024
	// ReturnRequestMaxBytes คือเพดานขนาด request ของการบันทึกรับคืน = รูปครบจำนวนสูงสุด + ที่เผื่อ (FR-RET-05)
	ReturnRequestMaxBytes = ReturnPhotoMax*ReturnPhotoMaxBytes + returnFormOverheadBytes
)

// ReturnPhotoTypes คือชนิดรูปที่รับ JPG หรือ PNG (FR-RET-05, ASM-06) ตรงกับ CHECK ใน 0006_return.up.sql
var ReturnPhotoTypes = []string{"image/jpeg", "image/png"}

// ReturnConditionStatus คือสภาพตอนรับคืน 2 ค่า และสถานะอุปกรณ์หลังบันทึกการรับคืนสำเร็จ
// ปกติ -> ว่าง (FR-RET-03, RQ-13) / ชำรุด -> ไม่พร้อมใช้งาน ไม่ผ่าน "ว่าง" (FR-RET-04, ASM-05)
// สูญหายไม่ใช่สภาพ (FR-RET-02) ค่าสถานะตาม CON-DATA-01 ใน specs/000-shared/data-model.md
var ReturnConditionStatus = map[string]string{
	"normal":  "available",
	"damaged": "unavailable",
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
