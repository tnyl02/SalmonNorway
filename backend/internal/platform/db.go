// Package platform คือของกลางที่ทุกฟีเจอร์ใช้ร่วม: การเชื่อมฐานข้อมูลและ HTTP router
package platform

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pinger คือสิ่งที่ตรวจได้ว่าฐานข้อมูลยังตอบอยู่ไหม ใช้กับ /healthz และแทนด้วยตัวจำลองตอน test
type Pinger interface {
	Ping(ctx context.Context) error
}

// OpenDB สร้าง connection pool ไปยัง PostgreSQL (CON-TECH-01)
// pool ยังไม่ต่อจริงจนกว่าจะมีการใช้งาน หลังบ้านจึงเปิดได้แม้ฐานข้อมูลยังไม่พร้อม แล้ว /healthz จะบอกสถานะเอง
func OpenDB(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, databaseURL)
}
