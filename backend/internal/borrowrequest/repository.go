package borrowrequest

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// Querier คือส่วนของ pgxpool.Pool ที่ repository นี้ใช้
type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Repository เขียนตาราง borrow_request (DM-06) และอ่าน equipment_item เพื่อตรวจสถานะ
// ทุก query ใช้ parameter ($1, $2, ...) ห้ามต่อสตริง SQL (plan.md ข้อ 6.3 ของ 001)
type Repository struct {
	db Querier
}

// NewRepository สร้าง Repository จาก connection pool
func NewRepository(db Querier) *Repository {
	return &Repository{db: db}
}

// insertSQL สร้างคำขอจากรหัสครุภัณฑ์ โดยคัดเฉพาะอุปกรณ์ที่สถานะ "ว่าง" (CON-DATA-01, FR-BRQ-01)
//
// ใช้ INSERT ... SELECT คำสั่งเดียว การอ่านสถานะกับการเขียนคำขอจึงอยู่ในคำสั่งเดียวกัน
// ไม่มีช่องให้สถานะเปลี่ยนระหว่างกลาง และ **ไม่มีคำสั่งใดแตะ equipment_item** (DOM-REQ-01)
const insertSQL = `
INSERT INTO borrow_request
	(requester_id, requester_name, student_code, item_id,
	 loan_type, purpose, purpose_other, start_at, end_at, status)
SELECT $1, $2, $3, i.id, $4, $5, $6, $7, $8, $9
FROM equipment_item i
WHERE i.asset_code = $10 AND i.status = 'available'
RETURNING id, submitted_at`

// Create บันทึกคำขอยืมใบใหม่ด้วยสถานะเริ่มต้นที่ service ส่งมา (FR-BRQ-04)
// คืน ErrItemNotAvailable เมื่อไม่มีอุปกรณ์รหัสนั้น หรือสถานะไม่ใช่ "ว่าง" แล้ว (CON-DATA-01)
func (r *Repository) Create(ctx context.Context, requesterID string, in NewRequest, status string) (Request, error) {
	out := Request{
		RequesterID: requesterID,
		AssetCode:   in.AssetCode,
		LoanType:    in.LoanType,
		Purpose:     in.Purpose,
		Status:      status,
	}
	err := r.db.QueryRow(ctx, insertSQL,
		requesterID, in.RequesterName, in.StudentCode,
		in.LoanType, in.Purpose, nullIfEmpty(in.PurposeOther),
		in.StartAt, in.EndAt, status, in.AssetCode,
	).Scan(&out.ID, &out.SubmittedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		// SELECT ไม่ได้แถว แปลว่าไม่มีอุปกรณ์รหัสนี้ หรือสถานะไม่ใช่ "ว่าง" แล้ว
		// อ่านสถานะล่าสุดกลับไปด้วย เพื่อให้แจ้งผู้ขอได้ว่าทำไมส่งไม่ได้ (FR-BRQ-01, AC-BRQ-01)
		status, statusErr := r.ItemStatus(ctx, in.AssetCode)
		if statusErr != nil {
			// ไม่มีรหัสครุภัณฑ์นี้อยู่จริง ตอบแบบเดียวกับกรณีไม่ว่าง
			// เพื่อไม่เปิดเผยว่ารหัสนั้นมีอยู่ในระบบหรือไม่ (DOM-ACCESS-01 ในภาพรวม)
			return Request{}, ErrItemNotAvailable
		}
		return Request{}, &ItemNotAvailableError{AssetCode: in.AssetCode, Status: status}
	}
	if err != nil {
		return Request{}, err
	}
	return out, nil
}

// nullIfEmpty เก็บ NULL แทนสตริงว่าง เพื่อให้คอลัมน์ purpose_other ว่างจริงเมื่อไม่ได้เลือก "อื่น ๆ" (FR-BRQ-02)
func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// statusSQL อ่านสถานะล่าสุดของอุปกรณ์ตามรหัสครุภัณฑ์ (CON-DATA-01, CON-DATA-03)
const statusSQL = `SELECT status FROM equipment_item WHERE asset_code = $1`

// ItemStatus คืนสถานะล่าสุดของอุปกรณ์ ใช้ตอนบอกผู้ขอว่าทำไมส่งคำขอไม่ได้ (FR-BRQ-01)
// ไม่พบรหัสครุภัณฑ์ถือเป็น ErrItemNotAvailable เพื่อไม่เปิดเผยว่ามีรหัสนั้นอยู่จริงหรือไม่
func (r *Repository) ItemStatus(ctx context.Context, assetCode string) (string, error) {
	var status string
	err := r.db.QueryRow(ctx, statusSQL, assetCode).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrItemNotAvailable
	}
	if err != nil {
		return "", err
	}
	return status, nil
}
