// Package returns คือฟีเจอร์ 006-return บันทึกการรับคืนอุปกรณ์ (UC-09) ตาม specs/006-return/plan.md
// เจ้าหน้าที่เท่านั้นที่ใช้ได้ และเห็นเฉพาะอุปกรณ์ของภาควิชาตนเอง (FR-RET-01, ASM-03)
package returns

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrForbiddenDepartment หมายถึงอุปกรณ์หรือการยืมเป็นของภาควิชาอื่น เจ้าหน้าที่ทำรายการไม่ได้ (FR-RET-01, ASM-03)
var ErrForbiddenDepartment = errors.New("equipment belongs to another department")

// querier คือสิ่งที่รัน QueryRow ได้ ใช้ได้ทั้ง pool และ transaction
type querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// checkItemDepartment ตรวจว่าอุปกรณ์ assetCode เป็นของภาควิชา departmentID (ASM-03)
// ตรวจก่อนเรื่องการยืมค้างเสมอ: อุปกรณ์ภาควิชาอื่นได้ ErrForbiddenDepartment ไม่ว่าจะถูกยืมอยู่หรือไม่
// ไม่มีอุปกรณ์รหัสนี้ -> ErrNoOpenLoan (ไม่มีการยืมค้างให้คืน)
func checkItemDepartment(ctx context.Context, q querier, assetCode string, departmentID int) error {
	var dept int
	err := q.QueryRow(ctx, `
		SELECT m.department_id FROM equipment_item i JOIN equipment_model m ON m.id = i.model_id
		WHERE i.asset_code = $1`, assetCode).Scan(&dept)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNoOpenLoan
	}
	if err != nil {
		return err
	}
	if dept != departmentID {
		return ErrForbiddenDepartment
	}
	return nil
}

// ErrNoOpenLoan หมายถึงอุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน หรือสถานะไม่ใช่ "ถูกยืม" (FR-RET-01)
var ErrNoOpenLoan = errors.New("no open loan for this asset")

// Repository คือจุดเดียวที่ฟีเจอร์นี้อ่านเขียนตาราง loan, equipment_item, audit_log (CON-TECH-01)
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository สร้าง Repository จาก connection pool ของ PostgreSQL (CON-TECH-01)
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// ReturnRecord คือผลของการบันทึกรับคืน 1 ครั้ง (FR-RET-02, FR-RET-03)
type ReturnRecord struct {
	LoanID       int
	AssetCode    string
	ReturnedAt   time.Time
	Condition    string
	ReturnedLate bool
	ItemStatus   string
}

// RecordReturn ปิดการยืม และเปลี่ยนสถานะอุปกรณ์ ใน transaction เดียว (plan.md ข้อ 4)
//   - เวลาคืน = now() ของฐานข้อมูลตอนบันทึก ไม่รับจาก client (FR-RET-02, ASM-04)
//   - สถานะเปลี่ยนครั้งเดียวจาก "ถูกยืม" ไปเป็น newStatus (FR-RET-03, CON-DATA-01)
//   - returned_late = now() > due_at ต้องใส่พร้อมกันเพราะ CHECK loan_return_fields_complete (FR-RET-06, test ใน T-13)
//
// ล็อกแถวด้วย FOR UPDATE: ถ้าบันทึกพร้อมกัน คนที่สองจะไม่พบการยืมที่ยังไม่คืน ได้ ErrNoOpenLoan
// อุปกรณ์ต้องเป็นของภาควิชา departmentID ของเจ้าหน้าที่ ไม่งั้นได้ ErrForbiddenDepartment และไม่แตะข้อมูล (ASM-03)
// audit ของการบันทึกสำเร็จเขียนใน transaction เดียวกัน: ไม่มีการรับคืนที่ไม่มี audit (FR-RET-07, DOM-AUDIT-01)
func (r *Repository) RecordReturn(ctx context.Context, assetCode, condition, newStatus string, departmentID int, audit AuditEntry) (ReturnRecord, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return ReturnRecord{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // rollback หลัง commit แล้วไม่มีผล

	if err := checkItemDepartment(ctx, tx, assetCode, departmentID); err != nil {
		return ReturnRecord{}, err
	}

	rec := ReturnRecord{AssetCode: assetCode, Condition: condition, ItemStatus: newStatus}
	var itemID int
	err = tx.QueryRow(ctx, `
		SELECT l.id, i.id
		FROM loan l
		JOIN equipment_item i ON i.id = l.item_id
		WHERE i.asset_code = $1 AND l.returned_at IS NULL AND i.status = 'borrowed'
		FOR UPDATE OF l, i`, assetCode).Scan(&rec.LoanID, &itemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ReturnRecord{}, ErrNoOpenLoan
	}
	if err != nil {
		return ReturnRecord{}, err
	}

	err = tx.QueryRow(ctx, `
		UPDATE loan
		SET returned_at = now(), return_condition = $2, returned_late = (now() > due_at)
		WHERE id = $1
		RETURNING returned_at, returned_late`, rec.LoanID, condition).Scan(&rec.ReturnedAt, &rec.ReturnedLate)
	if err != nil {
		return ReturnRecord{}, err
	}

	if _, err = tx.Exec(ctx,
		`UPDATE equipment_item SET status = $2, updated_at = now() WHERE id = $1`, itemID, newStatus); err != nil {
		return ReturnRecord{}, err
	}

	audit.Action, audit.AssetCode, audit.Outcome = actionReturn, assetCode, "ok"
	if err := insertAudit(ctx, tx, audit); err != nil {
		return ReturnRecord{}, err
	}

	return rec, tx.Commit(ctx)
}

// OpenLoan คือการยืมที่ยังไม่คืน 1 รายการ ตามรูปแบบใน plan.md ข้อ 4 (FR-RET-01)
type OpenLoan struct {
	LoanID     int       `json:"loanId"`
	AssetCode  string    `json:"assetCode"`
	ModelName  string    `json:"modelName"`
	BorrowerID string    `json:"borrowerId"`
	LoanType   string    `json:"loanType"`
	DueAt      time.Time `json:"dueAt"`
	IsOverdue  bool      `json:"isOverdue"`
}

// FindOpenLoan หาการยืมที่ยังไม่คืนของอุปกรณ์ assetCode (FR-RET-01(ก), CON-DATA-03 รหัสครุภัณฑ์ไม่ซ้ำ จึงได้ไม่เกิน 1 รายการ)
// isOverdue = now() > due_at คำนวณที่ฐานข้อมูล ใช้ due_at ที่ 005 บันทึก ไม่คำนวณ 16:00 เอง (ASM-01)
// ไม่พบ หรือสถานะไม่ใช่ "ถูกยืม" -> ErrNoOpenLoan / อุปกรณ์ภาควิชาอื่น -> ErrForbiddenDepartment (ASM-03)
func (r *Repository) FindOpenLoan(ctx context.Context, assetCode string, departmentID int) (OpenLoan, error) {
	if err := checkItemDepartment(ctx, r.pool, assetCode, departmentID); err != nil {
		return OpenLoan{}, err
	}
	var ol OpenLoan
	err := r.pool.QueryRow(ctx, `
		SELECT l.id, i.asset_code, m.name, l.borrower_id, l.loan_type, l.due_at, now() > l.due_at
		FROM loan l
		JOIN equipment_item i ON i.id = l.item_id
		JOIN equipment_model m ON m.id = i.model_id
		WHERE i.asset_code = $1 AND l.returned_at IS NULL AND i.status = 'borrowed'`, assetCode).
		Scan(&ol.LoanID, &ol.AssetCode, &ol.ModelName, &ol.BorrowerID, &ol.LoanType, &ol.DueAt, &ol.IsOverdue)
	if errors.Is(err, pgx.ErrNoRows) {
		return OpenLoan{}, ErrNoOpenLoan
	}
	return ol, err
}

// ErrLoanNotFound หมายถึงไม่มีการยืมรหัสนี้ (FR-RET-06 หน้ารายละเอียดการยืม)
var ErrLoanNotFound = errors.New("loan not found")

// LoanDetail คือรายละเอียดการยืม 1 รายการตาม plan.md ข้อ 4 (FR-RET-06)
// ยังไม่คืน: returnedAt, condition, returnedLate เป็น null
type LoanDetail struct {
	LoanID       int        `json:"loanId"`
	AssetCode    string     `json:"assetCode"`
	BorrowerID   string     `json:"borrowerId"`
	DueAt        time.Time  `json:"dueAt"`
	ReturnedAt   *time.Time `json:"returnedAt"`
	Condition    *string    `json:"condition"`
	ReturnedLate *bool      `json:"returnedLate"`
}

// GetLoan อ่านรายละเอียดการยืมจากรหัสการยืม (FR-RET-06) ใช้ returned_late ที่บันทึกตอนรับคืน ไม่คำนวณใหม่
// ไม่พบ -> ErrLoanNotFound / การยืมของอุปกรณ์ภาควิชาอื่น -> ErrForbiddenDepartment (ASM-03)
func (r *Repository) GetLoan(ctx context.Context, loanID, departmentID int) (LoanDetail, error) {
	var d LoanDetail
	var dept int
	err := r.pool.QueryRow(ctx, `
		SELECT l.id, i.asset_code, l.borrower_id, l.due_at, l.returned_at, l.return_condition, l.returned_late, m.department_id
		FROM loan l
		JOIN equipment_item i ON i.id = l.item_id
		JOIN equipment_model m ON m.id = i.model_id
		WHERE l.id = $1`, loanID).
		Scan(&d.LoanID, &d.AssetCode, &d.BorrowerID, &d.DueAt, &d.ReturnedAt, &d.Condition, &d.ReturnedLate, &dept)
	if errors.Is(err, pgx.ErrNoRows) {
		return LoanDetail{}, ErrLoanNotFound
	}
	if err != nil {
		return LoanDetail{}, err
	}
	if dept != departmentID {
		return LoanDetail{}, ErrForbiddenDepartment
	}
	return d, nil
}

// ListOpenLoans คืนการยืมที่ยังไม่คืนทั้งหมดของภาควิชา departmentID เรียงตามกำหนดคืน (FR-RET-01(ข), ASM-03)
// กรองภาควิชาที่ชั้น repository ทุก query (plan.md ข้อ 4) / ไม่แบ่งหน้า (plan.md Q-P1)
// isOverdue = now() > due_at ใช้ due_at ที่ 005 บันทึก (FR-RET-06, ASM-01)
func (r *Repository) ListOpenLoans(ctx context.Context, departmentID int) ([]OpenLoan, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT l.id, i.asset_code, m.name, l.borrower_id, l.loan_type, l.due_at, now() > l.due_at
		FROM loan l
		JOIN equipment_item i ON i.id = l.item_id
		JOIN equipment_model m ON m.id = i.model_id
		WHERE l.returned_at IS NULL AND i.status = 'borrowed' AND m.department_id = $1
		ORDER BY l.due_at, l.id`, departmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	loans := []OpenLoan{}
	for rows.Next() {
		var ol OpenLoan
		if err := rows.Scan(&ol.LoanID, &ol.AssetCode, &ol.ModelName, &ol.BorrowerID, &ol.LoanType, &ol.DueAt, &ol.IsOverdue); err != nil {
			return nil, err
		}
		loans = append(loans, ol)
	}
	return loans, rows.Err()
}
