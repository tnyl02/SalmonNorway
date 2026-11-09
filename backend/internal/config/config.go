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
