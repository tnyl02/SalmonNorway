package borrowrequest

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
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
	// Query ใช้กับคำสั่งที่คืนได้หลายแถว เช่น รายการยืมที่คืนแล้ว (FR-BRQ-06 ข้อ ข, T-13)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	// Exec ใช้กับคำสั่งที่ไม่คืนแถว เช่น บันทึก path ของไฟล์แนบ (FR-BRQ-10 ข้อ ก, T-14)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
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
	// ลำดับนี้คือคำตอบของ Q-BRQ-15 (ทีมตัดสิน 2569-11-12): ตรวจภาควิชาก่อนเสมอ
	// คนต่างภาควิชาได้ 403 โดยไม่รู้สถานะอุปกรณ์ FR-BRQ-07 ชนะ FR-BRQ-01 เมื่อสองข้อนี้ชนกัน
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

// earliestOpenLoanSQL หากำหนดคืนที่เร็วที่สุดของอุปกรณ์ที่ผู้ขอยังไม่คืน เฉพาะอุปกรณ์ของภาควิชาที่ระบุ
// (DM-05 ตาราง loan ของ 001, FR-BRQ-06 ข้อ ค, DOM-ACCESS-01)
// ใบที่กำหนดคืนเร็วที่สุดเกินกำหนดก็พอให้ปฏิเสธแล้ว จึงอ่านค่าเดียวแทนการดึงทุกใบ
const earliestOpenLoanSQL = `
SELECT min(l.due_at)
FROM loan l
JOIN equipment_item i ON i.id = l.item_id
JOIN equipment_model m ON m.id = i.model_id
WHERE l.borrower_id = $1
  AND l.returned_at IS NULL
  AND m.department_id = $2`

// EarliestOpenLoanDue คืนกำหนดคืนที่เร็วที่สุดของการยืมที่ผู้ขอยังไม่คืน (FR-BRQ-06)
// found = false แปลว่าไม่มีของค้างในภาควิชาตนเอง
// ฟีเจอร์นี้อ่านตาราง loan อย่างเดียว ไม่เขียนทับงานของ 005 และ 006 (DOM-REQ-01)
func (r *Repository) EarliestOpenLoanDue(ctx context.Context, borrowerID string, departmentID int) (time.Time, bool, error) {
	var dueAt *time.Time // NULL เมื่อไม่มีใบที่ยังไม่คืน
	if err := r.db.QueryRow(ctx, earliestOpenLoanSQL, borrowerID, departmentID).Scan(&dueAt); err != nil {
		return time.Time{}, false, err
	}
	if dueAt == nil {
		return time.Time{}, false, nil
	}
	return *dueAt, true, nil
}

// ReturnedLoan คือการยืมที่คืนแล้ว 1 ใบ เก็บเท่าที่กฎ FR-BRQ-06 ข้อ ข ต้องใช้
// กำหนดคืนใช้ตัดสินว่าใบนั้น "คืนช้า" หรือไม่ ส่วนเวลาที่คืนใช้ตัดสินว่าปลดล็อกหรือยัง
type ReturnedLoan struct {
	DueAt      time.Time
	ReturnedAt time.Time
}

// returnedSinceSQL อ่านการยืมที่คืนแล้วตั้งแต่เวลาที่กำหนด เฉพาะอุปกรณ์ของภาควิชาที่ระบุ
// (DM-05 ตาราง loan ของ 001, FR-BRQ-06 ข้อ ค, DOM-ACCESS-01)
// เวลาเริ่มนับมาจาก service เพราะกฎ "00:00 น. ของวันถัดไป" เป็นกฎของฟีเจอร์ ไม่ใช่ของฐานข้อมูล
const returnedSinceSQL = `
SELECT l.due_at, l.returned_at
FROM loan l
JOIN equipment_item i ON i.id = l.item_id
JOIN equipment_model m ON m.id = i.model_id
WHERE l.borrower_id = $1
  AND m.department_id = $2
  AND l.returned_at IS NOT NULL
  AND l.returned_at >= $3
ORDER BY l.returned_at DESC`

// ReturnedSince คืนการยืมที่คืนแล้วตั้งแต่เวลา since เป็นต้นมา (FR-BRQ-06 ข้อ ข)
// อ่านตาราง loan อย่างเดียวเหมือน T-12 ไม่เขียนทับงานของ 006 ที่เป็นผู้บันทึกการคืน
func (r *Repository) ReturnedSince(ctx context.Context, borrowerID string, departmentID int, since time.Time) ([]ReturnedLoan, error) {
	rows, err := r.db.Query(ctx, returnedSinceSQL, borrowerID, departmentID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ReturnedLoan
	for rows.Next() {
		var l ReturnedLoan
		if err := rows.Scan(&l.DueAt, &l.ReturnedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

// setAttachmentPathSQL บันทึกที่อยู่ของไฟล์ใบบันทึกข้อความลงคำขอใบนั้น (FR-BRQ-10 ข้อ ก)
// เก็บเฉพาะ path ตัวไฟล์อยู่บนเครื่องที่รันระบบ ไม่ได้อยู่ในฐานข้อมูล (ASM-BRQ-07)
const setAttachmentPathSQL = `UPDATE borrow_request SET attachment_path = $2 WHERE id = $1`

// SetAttachmentPath ผูกไฟล์ที่เก็บแล้วเข้ากับคำขอ (FR-BRQ-10 ข้อ ก)
func (r *Repository) SetAttachmentPath(ctx context.Context, requestID int, path string) error {
	_, err := r.db.Exec(ctx, setAttachmentPathSQL, requestID, path)
	return err
}

// deleteRequestSQL ลบคำขอที่สร้างไว้แล้วแต่เก็บไฟล์แนบไม่สำเร็จ
// การยืมค้างคืนต้องมีไฟล์แนบตั้งแต่ตอนส่ง การทิ้งคำขอที่ไม่มีไฟล์ไว้จึงเป็นคำขอค้างที่ไม่ควรมี (FR-BRQ-09)
const deleteRequestSQL = `DELETE FROM borrow_request WHERE id = $1`

// Delete ลบคำขอใบที่ระบุ ใช้เฉพาะตอนย้อนกลับภายใน request เดียวกัน
// ไม่ใช่การยกเลิกคำขอตามความต้องการของผู้ใช้ ซึ่งเป็นของ 003 (UC-15)
func (r *Repository) Delete(ctx context.Context, requestID int) error {
	_, err := r.db.Exec(ctx, deleteRequestSQL, requestID)
	return err
}

// PendingRequest คือคำขอที่ยังไม่จบ 1 ใบของอุปกรณ์ชิ้นหนึ่ง ใช้สร้างลำดับคำขอให้ 001 (FR-BRQ-05 ข้อ ก, CR-005)
type PendingRequest struct {
	RequestID     int
	SubmittedAt   time.Time // ใช้เรียงลำดับคำขอที่มาก่อน (DOM-REQ-01)
	RequesterName string    // ชื่อที่ผู้ขอพิมพ์เองในฟอร์ม ระบบไม่ได้ตรวจกับระบบยืนยันตัวตน (FR-BRQ-03 ข้อ ก)
}

// pendingByItemSQL อ่านคำขอที่ยังไม่จบของอุปกรณ์ชิ้นนั้น เรียงตามเวลาที่ส่ง (FR-BRQ-05 ข้อ ก, DOM-REQ-01)
// รายการสถานะมาจาก config.OpenRequestStatuses ที่เดียว ไม่เขียนซ้ำใน SQL (AGENTS.md กฎข้อ 9)
// เรียงด้วย id ต่อท้ายกันกรณีเวลาที่ส่งเท่ากันพอดี เพื่อให้ลำดับคงที่ทุกครั้งที่เรียก
const pendingByItemSQL = `
SELECT r.id, r.submitted_at, r.requester_name
FROM borrow_request r
JOIN equipment_item i ON i.id = r.item_id
WHERE i.asset_code = $1 AND r.status = ANY($2)
ORDER BY r.submitted_at, r.id`

// PendingByItem คืนคำขอที่ยังไม่จบของอุปกรณ์ชิ้นนั้นเรียงตามเวลาที่ส่ง (FR-BRQ-05 ข้อ ก)
// รหัสครุภัณฑ์ที่ไม่มีอยู่จริงคืนรายการว่าง เหมือนกรณีไม่มีคำขอค้าง (ทีมตัดสิน 2569-11-12 ตามสัญญาใน CR-005)
func (r *Repository) PendingByItem(ctx context.Context, assetCode string) ([]PendingRequest, error) {
	rows, err := r.db.Query(ctx, pendingByItemSQL, assetCode, config.OpenRequestStatuses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []PendingRequest{} // ไม่ใช่ nil เพื่อให้ JSON ออกมาเป็น [] ไม่ใช่ null (สัญญา API ใน CR-005)
	for rows.Next() {
		var p PendingRequest
		if err := rows.Scan(&p.RequestID, &p.SubmittedAt, &p.RequesterName); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// RequestAttachment คือข้อมูลที่ต้องใช้ตัดสินสิทธิ์เปิดดูไฟล์แนบ (FR-BRQ-10 ข้อ ก, DOM-PRIV-01)
type RequestAttachment struct {
	RequesterID string // เจ้าของคำขอ ใช้เทียบกับผู้ที่ยืนยันตัวตนมา (IF-IDP-01)
	Path        string // ที่อยู่ไฟล์บนเครื่องที่รันระบบ (ASM-BRQ-07)
}

// attachmentByIDSQL อ่านเจ้าของคำขอและที่อยู่ไฟล์แนบของคำขอใบนั้น (FR-BRQ-10 ข้อ ก)
const attachmentByIDSQL = `SELECT requester_id, attachment_path FROM borrow_request WHERE id = $1`

// AttachmentOf คืนเจ้าของคำขอและที่อยู่ไฟล์แนบ
// found = false ทั้งกรณีที่ไม่มีคำขอใบนั้น และกรณีที่คำขอใบนั้นไม่ได้แนบไฟล์ไว้
// ผู้เรียกจึงแยกสองกรณีนี้ไม่ได้ ซึ่งตั้งใจให้เป็นแบบนั้น จะได้ไม่บอกคนนอกว่าเลขคำขอใดมีอยู่จริง (DOM-PRIV-01)
func (r *Repository) AttachmentOf(ctx context.Context, requestID int) (RequestAttachment, bool, error) {
	var out RequestAttachment
	var path *string // NULL เมื่อคำขอนั้นไม่มีไฟล์แนบ
	err := r.db.QueryRow(ctx, attachmentByIDSQL, requestID).Scan(&out.RequesterID, &path)
	if errors.Is(err, pgx.ErrNoRows) {
		return RequestAttachment{}, false, nil
	}
	if err != nil {
		return RequestAttachment{}, false, err
	}
	if path == nil {
		return RequestAttachment{}, false, nil
	}
	out.Path = *path
	return out, true, nil
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
