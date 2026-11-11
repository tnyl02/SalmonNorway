package borrowrequest

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"salmonnorway/backend/internal/auth"
)

// ErrForbiddenDepartment เกิดเมื่อผู้ขอพยายามส่งคำขอของอุปกรณ์ต่างภาควิชา (DOM-ACCESS-01, FR-BRQ-07)
// ตอบโดยไม่ส่งข้อมูลอุปกรณ์กลับไปเลย
// (ประกาศไว้ที่นี่เพราะ model.go ไม่อยู่ในช่องไฟล์ที่แตะของ T-10 ย้ายไปรวมกับ error ตัวอื่นได้ภายหลัง)
var ErrForbiddenDepartment = errors.New("อุปกรณ์ชิ้นนี้ไม่ใช่ของภาควิชาผู้ขอ")

// ErrDuplicateRequest เกิดเมื่อผู้ขอคนเดิมส่งคำขออุปกรณ์ชิ้นเดิมซ้ำ ขณะที่คำขอเก่ายังไม่จบ
// คือสถานะยังเป็น requested, approved_by_officer, pending_advisor_approval หรือ approved
// (FR-BRQ-05 ข, CR-004) คำขอที่ถูกยกเลิกหรือถูกปฏิเสธไปแล้วไม่นับ จึงส่งใหม่ได้
// (ประกาศไว้ที่นี่เพราะ model.go ไม่อยู่ในช่องไฟล์ที่แตะของ T-11 ย้ายไปรวมกับ error ตัวอื่นได้ภายหลัง)
var ErrDuplicateRequest = errors.New("มีคำขอของอุปกรณ์ชิ้นนี้ที่ยังไม่จบอยู่แล้ว")

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
// คืน ErrForbiddenDepartment เมื่ออุปกรณ์เป็นของภาควิชาอื่น (DOM-ACCESS-01, FR-BRQ-07)
// คืน ErrItemNotAvailable เมื่อไม่มีอุปกรณ์รหัสนั้น หรือสถานะไม่ใช่ "ว่าง" แล้ว (CON-DATA-01)
func (r *Repository) Create(ctx context.Context, requesterID string, in NewRequest, status string) (Request, error) {
	// ตรวจขอบเขตภาควิชาก่อนทุกอย่าง เพื่อไม่ให้ผู้ขอต่างภาควิชารู้แม้แต่สถานะของอุปกรณ์ (FR-BRQ-07)
	// **ลำดับการตรวจเมื่อข้ามภาควิชาและอุปกรณ์ไม่ว่างพร้อมกัน ยังรอ Q-BRQ-15 (ดู T-21)**
	// ระหว่างรอคำตอบเลือกทางที่ปลอดภัยกว่าคือไม่เปิดเผยอะไรเลย
	if err := r.checkDepartment(ctx, in.AssetCode); err != nil {
		return Request{}, err
	}

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
	if isDuplicateRequest(err) {
		// index กันคำขอซ้ำไว้ที่ฐานข้อมูล จึงกันได้แม้สองคำขอมาถึงพร้อมกัน (FR-BRQ-05 ข, AC-BRQ-11)
		return Request{}, ErrDuplicateRequest
	}
	if err != nil {
		return Request{}, err
	}
	return out, nil
}

// uniqueOpenRequestIndex คือ partial unique index ที่ migration 0003 สร้างไว้ (T-01)
// ครอบเฉพาะคำขอที่ยังไม่จบของคู่ (requester_id, item_id) ตามรายการสถานะใน FR-BRQ-05 ข
const uniqueOpenRequestIndex = "uq_request_open_per_item"

// uniqueViolationCode คือรหัส SQLSTATE ของ PostgreSQL สำหรับ unique violation
// เป็นค่าคงที่ของฐานข้อมูล ไม่ใช่ตัวเลขที่มาจาก spec จึงไม่ต้องอยู่ใน config (AGENTS.md กฎข้อ 9)
const uniqueViolationCode = "23505"

// isDuplicateRequest บอกว่า error จากการ INSERT เป็นการชน index กันคำขอซ้ำหรือไม่ (FR-BRQ-05 ข)
// ตรวจชื่อ constraint ด้วย เพื่อไม่เหมาเอาว่า unique violation ตัวอื่นในตารางนี้คือคำขอซ้ำ
func isDuplicateRequest(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	return pgErr.Code == uniqueViolationCode && pgErr.ConstraintName == uniqueOpenRequestIndex
}

// nullIfEmpty เก็บ NULL แทนสตริงว่าง เพื่อให้คอลัมน์ purpose_other ว่างจริงเมื่อไม่ได้เลือก "อื่น ๆ" (FR-BRQ-02)
func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// departmentSQL อ่านภาควิชาเจ้าของอุปกรณ์จากรุ่นของอุปกรณ์ชิ้นนั้น (DM-03, DM-04)
const departmentSQL = `
SELECT m.department_id
FROM equipment_item i
JOIN equipment_model m ON m.id = i.model_id
WHERE i.asset_code = $1`

// checkDepartment ปฏิเสธคำขอของอุปกรณ์ต่างภาควิชา (FR-BRQ-07, DOM-ACCESS-01)
// ภาควิชาของผู้ขอมาจากผลยืนยันตัวตนเท่านั้น (IF-IDP-01, DM-08) ไม่ใช่ค่าที่ส่งมาใน request
// รหัสครุภัณฑ์ที่ไม่มีอยู่จริงตอบเป็น ErrItemNotAvailable เหมือนกรณีไม่ว่าง เพื่อไม่เปิดเผยว่ามีรหัสนั้นไหม
func (r *Repository) checkDepartment(ctx context.Context, assetCode string) error {
	p, ok := auth.FromContext(ctx)
	if !ok || !p.Complete() {
		return ErrNoPrincipal // ไม่รู้ว่าผู้ขออยู่ภาควิชาไหน ปฏิเสธ ไม่เดา (plan.md ข้อ 6.1 ของ 001)
	}

	var departmentID int
	err := r.db.QueryRow(ctx, departmentSQL, assetCode).Scan(&departmentID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrItemNotAvailable
	}
	if err != nil {
		return err
	}
	if departmentID != p.DepartmentID {
		return ErrForbiddenDepartment
	}
	return nil
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
