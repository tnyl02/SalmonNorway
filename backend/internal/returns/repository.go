// Package returns คือฟีเจอร์ 006-return บันทึกการรับคืนอุปกรณ์ (UC-09) ตาม specs/006-return/plan.md
// เจ้าหน้าที่เท่านั้นที่ใช้ได้ และเห็นเฉพาะอุปกรณ์ของภาควิชาตนเอง (FR-RET-01, ASM-03)
package returns

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository คือจุดเดียวที่ฟีเจอร์นี้อ่านเขียนตาราง loan, equipment_item, loan_return_photo, audit_log (CON-TECH-01)
// query ของแต่ละ endpoint จะเพิ่มใน task ของ endpoint นั้น (T-04, T-05, T-14, T-17)
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository สร้าง Repository จาก connection pool ของ PostgreSQL (CON-TECH-01)
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
