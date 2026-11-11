package cancelrequest

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// Querier คือส่วนของ pgxpool.Pool ที่ repository ใช้ (แทนด้วยตัวจำลองได้)
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Repository อ่านและเขียนตาราง borrow_request ของ 002 (DM-06) เฉพาะส่วนที่ 003 ใช้
// ทุก query ใช้ parameter ไม่ต่อสตริง SQL และไม่แตะตาราง equipment_item (CON-CXL-01, CON-DATA-01)
type Repository struct {
	db Querier
}

// NewRepository สร้าง Repository จาก connection pool
func NewRepository(db Querier) *Repository {
	return &Repository{db: db}
}

// PendingRow คือคำขอที่ยังยกเลิกได้ 1 แถว พร้อมชื่อรุ่นและรหัสครุภัณฑ์ (FR-CXL-06)
type PendingRow struct {
	RequestID   int
	ModelName   string
	AssetCode   string
	SubmittedAt time.Time
}

// listPendingSQL ดึงคำขอของผู้ขอคนเดียวที่สถานะยกเลิกได้ (FR-CXL-03, CON-CXL-02)
// เรียงจากส่งล่าสุดก่อน เวลาเท่ากันเรียงรหัสคำขอมากไปน้อย (FR-CXL-06, ASM-13) ใช้ index idx_request_owner_status ของ 0003
const listPendingSQL = `
SELECT r.id, m.name, i.asset_code, r.submitted_at
FROM borrow_request r
JOIN equipment_item  i ON i.id = r.item_id
JOIN equipment_model m ON m.id = i.model_id
WHERE r.requester_id = $1 AND r.status = ANY($2::text[])
ORDER BY r.submitted_at DESC, r.id DESC`

// ListPending คืนคำขอที่ยังยกเลิกได้ของ requesterID (FR-CXL-06)
func (r *Repository) ListPending(ctx context.Context, requesterID string, statuses []string) ([]PendingRow, error) {
	rows, err := r.db.Query(ctx, listPendingSQL, requesterID, statuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PendingRow
	for rows.Next() {
		var p PendingRow
		if err := rows.Scan(&p.RequestID, &p.ModelName, &p.AssetCode, &p.SubmittedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
