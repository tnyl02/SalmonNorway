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
