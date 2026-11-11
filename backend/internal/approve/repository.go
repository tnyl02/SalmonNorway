package approve

import (
	"context"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/config"
)

// Querier คือส่วนของ pgxpool.Pool หรือ pgx.Tx ที่ repository ใช้
// รับ pgx.Tx ได้ด้วย เพื่อให้การพิจารณาทำใน transaction เดียว (plan.md ข้อ 1, T-07)
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Repository อ่านและเขียนข้อมูลคำขอยืมของฟีเจอร์ 004 (CON-TECH-01)
// ทุก query ใช้ parameter ($1, $2, ...) ห้ามต่อสตริง SQL จากค่าที่ผู้ใช้ส่งมา
type Repository struct {
	db Querier
}

// NewRepository สร้าง Repository จาก connection pool หรือ transaction
func NewRepository(db Querier) *Repository {
	return &Repository{db: db}
}

// ListPending คืนคำขอยืมแบบปกติสถานะ requested ของอุปกรณ์ภาควิชาที่ระบุ เรียงตามเวลาที่ส่งจากเก่าไปใหม่
// พร้อมบอกว่าแถวไหนเป็นคำขอแรกของอุปกรณ์รายชิ้นนั้น (FR-APV-01 ก ข, DOM-REQ-01, ASM-06)
//
// คำขอค้างคืนไม่ขึ้นในรายการ (Out of scope รอ Q-APV-01) และไม่นับในลำดับคำขอแรก (Q-APV-23)
// ภาควิชาเทียบกับภาควิชาเจ้าของอุปกรณ์ (equipment_model.department_id, DM-03) ไม่ใช่ภาควิชาของผู้ขอ
func (r *Repository) ListPending(ctx context.Context, departmentID int) ([]PendingRequest, error) {
	const q = `
SELECT r.id, r.requester_name, r.student_code, m.name, i.asset_code,
       r.purpose, coalesce(r.purpose_other, ''), r.start_at, r.end_at, r.submitted_at,
       row_number() OVER (PARTITION BY r.item_id ORDER BY r.submitted_at, r.id) = 1 AS is_first
FROM borrow_request r
JOIN equipment_item  i ON i.id = r.item_id
JOIN equipment_model m ON m.id = i.model_id
WHERE r.status = $1
  AND r.loan_type = $2
  AND m.department_id = $3
ORDER BY r.submitted_at, r.id`

	rows, err := r.db.Query(ctx, q, config.StatusRequested, config.LoanTypeNormal, departmentID)
	if err != nil {
		return nil, err
	}
	out, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (PendingRequest, error) {
		var p PendingRequest
		err := row.Scan(&p.ID, &p.RequesterName, &p.StudentCode, &p.ModelName, &p.AssetCode,
			&p.Purpose, &p.PurposeOther, &p.StartAt, &p.EndAt, &p.SubmittedAt, &p.IsFirst)
		return p, err
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		out = []PendingRequest{} // ไม่มีคำขอ = รายการว่าง ไม่ใช่ null
	}
	return out, nil
}
