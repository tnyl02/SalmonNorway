package approve

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"salmonnorway/backend/internal/config"
)

// Querier คือส่วนของ pgxpool.Pool หรือ pgx.Tx ที่ repository ใช้
// รับ pgx.Tx ได้ด้วย เพื่อให้การพิจารณาทำใน transaction เดียว (plan.md ข้อ 1, T-07)
// Begin บน pgx.Tx คือ savepoint test จึงครอบทุกอย่างใน transaction ที่ยกเลิกทิ้งได้
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
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

// inTx ทำงานทั้งหมดใน transaction เดียว สำเร็จทั้งหมดหรือไม่บันทึกเลย (plan.md ข้อ 1)
func (r *Repository) inTx(ctx context.Context, fn func(tx *Repository) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }() // ไม่มีผลถ้า commit แล้ว
	if err := fn(NewRepository(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// lockItemRequests ล็อกคำขอทุกใบ (ทุกประเภทการยืม) ของอุปกรณ์รายชิ้นเดียวกับคำขอที่ระบุ ด้วย SELECT ... FOR UPDATE
// คืนคำขอที่ระบุ และคำขอทั้งหมดของชิ้นนั้นเรียงตามเวลาที่ส่ง (plan.md ข้อ 7.3, FR-APV-01 ช, AC-APV-10)
// ล็อกเฉพาะแถว borrow_request (FOR UPDATE OF r) ไม่ล็อกตารางอุปกรณ์ของ 001
// ok = false ถ้าไม่มีคำขอรหัสนี้
func (r *Repository) lockItemRequests(ctx context.Context, requestID int) (target requestRow, all []requestRow, ok bool, err error) {
	const q = `
SELECT r.id, r.item_id, r.loan_type, r.status, r.start_at, r.submitted_at,
       i.asset_code, i.status, m.department_id
FROM borrow_request r
JOIN equipment_item  i ON i.id = r.item_id
JOIN equipment_model m ON m.id = i.model_id
WHERE r.item_id = (SELECT item_id FROM borrow_request WHERE id = $1)
ORDER BY r.submitted_at, r.id
FOR UPDATE OF r`
	rows, err := r.db.Query(ctx, q, requestID)
	if err != nil {
		return requestRow{}, nil, false, err
	}
	all, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (requestRow, error) {
		var x requestRow
		err := row.Scan(&x.ID, &x.ItemID, &x.LoanType, &x.Status, &x.StartAt, &x.SubmittedAt,
			&x.AssetCode, &x.ItemStatus, &x.DepartmentID)
		return x, err
	})
	if err != nil {
		return requestRow{}, nil, false, err
	}
	for _, x := range all {
		if x.ID == requestID {
			return x, all, true, nil
		}
	}
	return requestRow{}, all, false, nil
}

// setStatus เปลี่ยนสถานะของคำขอ 004 ตั้งได้เฉพาะ approved / rejected (data-model.md v9)
// ไม่แตะ equipment_item.status (CON-DATA-01, data-model.md แถว UC-12)
func (r *Repository) setStatus(ctx context.Context, requestID int, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE borrow_request SET status = $2 WHERE id = $1`, requestID, status)
	return err
}

// insertDecision บันทึกผลการพิจารณา 1 แถวต่อคำขอ = ประวัติตาม RQ-25 (FR-APV-01 ฌ) คืนเวลาที่บันทึก
func (r *Repository) insertDecision(ctx context.Context, d decision) (time.Time, error) {
	var reason, actorName *string
	if d.Reason != "" {
		reason = &d.Reason
	}
	if d.ActorName != "" {
		actorName = &d.ActorName
	}
	var at time.Time
	err := r.db.QueryRow(ctx, `
INSERT INTO approval_decision (request_id, result, reason, actor_id, actor_name, actor_department_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING decided_at`,
		d.RequestID, d.Result, reason, d.ActorID, actorName, d.ActorDepartmentID).Scan(&at)
	return at, err
}

// insertAudit บันทึก audit log ของการพิจารณา (RQ-25, DM-07) ค่า action ตาม migration 0008
func (r *Repository) insertAudit(ctx context.Context, actorID, actorRole, action, assetCode string, requestID int) error {
	_, err := r.db.Exec(ctx, `
INSERT INTO audit_log (actor_id, actor_role, action, asset_code, request_id, outcome)
VALUES ($1, $2, $3, $4, $5, 'ok')`,
		actorID, actorRole, action, assetCode, requestID)
	return err
}

// insertEmail ใส่ Email แจ้งผล 1 ฉบับลงคิว (outbox) ใน transaction เดียวกับผลการพิจารณา (plan.md ข้อ 1, FR-APV-02)
// ถ้าบันทึกผลไม่สำเร็จ Email ก็ไม่ถูกสร้าง และถ้าบันทึกผลสำเร็จ Email ไม่หาย แม้ระบบดับก่อนส่ง
func (r *Repository) insertEmail(ctx context.Context, requestID int, kind string, now time.Time) error {
	_, err := r.db.Exec(ctx, `
INSERT INTO approval_email (request_id, kind, status, next_attempt_at, created_at)
VALUES ($1, $2, $3, $4, $4)`,
		requestID, kind, config.ApprovalEmailPending, now)
	return err
}

// decisionTime คืนเวลาที่บันทึกผลการพิจารณาของคำขอ ok = false ถ้ายังไม่มีผล (เช่น สถานะถูกตั้งโดยฟีเจอร์อื่น)
// ใช้แยก "พิจารณาแล้ว" กับ "กดพร้อมกัน" (FR-APV-01 ซ, T-23)
func (r *Repository) decisionTime(ctx context.Context, requestID int) (time.Time, bool, error) {
	var at time.Time
	err := r.db.QueryRow(ctx, `SELECT decided_at FROM approval_decision WHERE request_id = $1`, requestID).Scan(&at)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, err
	}
	return at, true, nil
}

// ListFailedEmails คืน Email สถานะ failed ของคำขอที่อุปกรณ์เป็นของภาควิชาที่ระบุ รวมของการปฏิเสธอัตโนมัติ (System)
// เรียงตามเวลาที่ลองส่งครั้งล่าสุด (FR-APV-02 ง, ASM-06, ASM-08, T-33)
func (r *Repository) ListFailedEmails(ctx context.Context, departmentID int) ([]FailedEmail, error) {
	const q = `
SELECT e.id, e.request_id, e.kind, r.requester_name, m.name, i.asset_code,
       e.attempts, coalesce(e.last_error, ''), e.last_attempt_at
FROM approval_email e
JOIN borrow_request  r ON r.id = e.request_id
JOIN equipment_item  i ON i.id = r.item_id
JOIN equipment_model m ON m.id = i.model_id
WHERE e.status = $1 AND m.department_id = $2
ORDER BY e.last_attempt_at, e.id`
	rows, err := r.db.Query(ctx, q, config.ApprovalEmailFailed, departmentID)
	if err != nil {
		return nil, err
	}
	out, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (FailedEmail, error) {
		var f FailedEmail
		err := row.Scan(&f.ID, &f.RequestID, &f.Kind, &f.RequesterName, &f.ModelName, &f.AssetCode,
			&f.Attempts, &f.LastError, &f.LastAttemptAt)
		return f, err
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		out = []FailedEmail{}
	}
	return out, nil
}

// lockEmail ล็อก Email 1 ฉบับ คืนสถานะการส่งและภาควิชาเจ้าของอุปกรณ์ ok = false ถ้าไม่มี (T-33)
func (r *Repository) lockEmail(ctx context.Context, emailID int) (status string, departmentID int, ok bool, err error) {
	err = r.db.QueryRow(ctx, `
SELECT e.status, m.department_id
FROM approval_email e
JOIN borrow_request  r ON r.id = e.request_id
JOIN equipment_item  i ON i.id = r.item_id
JOIN equipment_model m ON m.id = i.model_id
WHERE e.id = $1
FOR UPDATE OF e`, emailID).Scan(&status, &departmentID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", 0, false, nil
	}
	return status, departmentID, err == nil, err
}

// requeueEmail ส่ง Email ที่ failed กลับเข้าคิว ตัวส่งเบื้องหลังส่งในรอบถัดไป (FR-APV-02 ง "กดส่งซ้ำได้โดยไม่ต้องพิจารณาใหม่")
// เริ่มนับครั้งใหม่ จึงได้ส่งซ้ำอัตโนมัติอีกชุดถ้ายังล้ม ไม่แตะ approval_decision และ borrow_request
func (r *Repository) requeueEmail(ctx context.Context, emailID int, now time.Time) error {
	_, err := r.db.Exec(ctx, `
UPDATE approval_email SET status = $2, attempts = 0, next_attempt_at = $3 WHERE id = $1`,
		emailID, config.ApprovalEmailPending, now)
	return err
}
