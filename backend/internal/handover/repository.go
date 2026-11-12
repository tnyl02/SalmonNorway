package handover

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/audit"
	"salmonnorway/backend/internal/auth"
)

// auditActionHandover คือค่า action ของ audit_log สำหรับการส่งมอบ (DM-07, FR-HND-08, Q-HND-06 ปิด 2569-11-12)
// อยู่ใน CHECK ด้วย migration 0009 ของ 005
const auditActionHandover = "handover"

// ค่าสถานะในฐานข้อมูลที่ 005 อ่านและเขียน (CON-DATA-01, data-model.md ตารางสถานะของคำขอยืม)
const (
	itemStatusAvailable     = "available"   // ว่าง
	itemStatusBorrowed      = "borrowed"    // ถูกยืม
	requestStatusHandedOver = "handed_over" // ส่งมอบแล้ว (ตั้งค่าโดย UC-08 ตาม data-model.md)
)

var (
	// ErrRequestNotFound หมายถึงไม่มีคำขอยืมรหัสนี้
	ErrRequestNotFound = errors.New("borrow request not found")
	// ErrRequestNotReady หมายถึงคำขอไม่อยู่ในสถานะที่ส่งมอบได้ เช่น ยังไม่อนุมัติ หรือส่งมอบไปแล้ว (FR-HND-01)
	ErrRequestNotReady = errors.New("borrow request not ready for handover")
	// ErrBorrowerMismatch หมายถึงผู้ยืมที่ยืนยันตัวตนไม่ใช่ผู้ส่งคำขอ (FR-HND-02, ASM-HND-04)
	ErrBorrowerMismatch = errors.New("borrower is not the requester")
	// ErrItemNotAvailable หมายถึงอุปกรณ์รายชิ้นในคำขอไม่ได้มีสถานะ "ว่าง" (FR-HND-04, CON-DATA-01)
	ErrItemNotAvailable = errors.New("equipment item not available")
)

// PhotoSaver เขียนและลบรูปของการยืม 1 รายการ (FR-HND-06) แยกจาก repository เพื่อให้ test ใช้ตัวจำลองได้
type PhotoSaver interface {
	SavePhotos(loanID int) ([]StoredPhoto, error)
	RemovePhotos(loanID int)
}

// TxBeginner คือส่วนของ pgxpool.Pool (หรือ pgx.Conn) ที่ repository ใช้เปิด transaction
type TxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// Repository บันทึกการส่งมอบลงฐานข้อมูล (FR-HND-04, FR-HND-07, CON-TECH-01)
type Repository struct {
	db TxBeginner
}

// NewRepository สร้าง repository จาก connection pool
func NewRepository(db TxBeginner) *Repository {
	return &Repository{db: db}
}

// HandoverRecord คือข้อมูลการส่งมอบที่ service ตรวจและเตรียมไว้แล้ว (FR-HND-02, 03, 05)
type HandoverRecord struct {
	RequestID         int
	BorrowerUserID    string    // รหัสผู้ใช้ที่ยืนยันตัวตนแล้ว เก็บใน loan.borrower_id (DM-05)
	BorrowerStudentID string    // FR-HND-03
	BorrowerFirstName string    // FR-HND-03 เจ้าหน้าที่พิมพ์ตามบัตร (ASM-HND-10)
	BorrowerLastName  string    // FR-HND-03 เจ้าหน้าที่พิมพ์ตามบัตร (ASM-HND-10)
	IdentityMethod    string    // FR-HND-02: auth_api หรือ card_only
	CardChecked       bool      // FR-HND-03
	HandedOverBy      string    // รหัสเจ้าหน้าที่ (CON-HND-01)
	HandedOverAt      time.Time // FR-HND-05: เวลาส่งมอบจริง
	DueAt             time.Time // FR-HND-05: กำหนดคืนที่ service คำนวณ
}

// HandoverResult คือผลของการส่งมอบที่สำเร็จ
type HandoverResult struct {
	LoanID    int
	ItemID    int
	AssetCode string // CON-DATA-03: ระบุอุปกรณ์ด้วยรหัสครุภัณฑ์
}

// CreateHandover บันทึกการส่งมอบใน transaction เดียว (FR-HND-04, FR-HND-07, DOM-REQ-01)
// ล็อกแถวอุปกรณ์ด้วย SELECT ... FOR UPDATE กันเจ้าหน้าที่ 2 คนส่งมอบชิ้นเดียวกันพร้อมกัน
// ตรวจว่าอุปกรณ์ "ว่าง" และคำขออยู่ใน readyStatuses แล้วสร้าง Loan + แถว handover
// เปลี่ยนอุปกรณ์เป็น "ถูกยืม" และคำขอเป็น "ส่งมอบแล้ว" ไม่แตะคำขออื่นของอุปกรณ์ชิ้นเดียวกัน (ASM-HND-08)
// ถ้าไม่ผ่านข้อใด transaction ถูกยกเลิก ไม่มีอะไรเปลี่ยน
// photos (ว่างได้) เขียนรูปลงดิสก์หลังสร้าง Loan แล้วบันทึกแถว handover_photo ใน transaction เดียวกัน
// ถ้าขั้นใดหลังเขียนรูปล้มเหลว ลบรูปทิ้ง ไม่ให้มีไฟล์ค้าง (FR-HND-06, plan.md ข้อ 4)
func (r *Repository) CreateHandover(ctx context.Context, rec HandoverRecord, readyStatuses []string, photos PhotoSaver) (HandoverResult, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return HandoverResult{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // หลัง Commit แล้ว Rollback ไม่มีผล

	var res HandoverResult
	err = tx.QueryRow(ctx, `SELECT item_id FROM borrow_request WHERE id = $1`, rec.RequestID).Scan(&res.ItemID)
	if errors.Is(err, pgx.ErrNoRows) {
		return HandoverResult{}, ErrRequestNotFound
	}
	if err != nil {
		return HandoverResult{}, err
	}

	// FR-HND-04: ล็อกอุปกรณ์ก่อน แล้วตรวจว่า "ว่าง"
	var itemStatus string
	if err := tx.QueryRow(ctx,
		`SELECT status, asset_code FROM equipment_item WHERE id = $1 FOR UPDATE`, res.ItemID,
	).Scan(&itemStatus, &res.AssetCode); err != nil {
		return HandoverResult{}, err
	}
	if itemStatus != itemStatusAvailable {
		return HandoverResult{}, ErrItemNotAvailable
	}

	// FR-HND-01: คำขอต้องยังอยู่ในสถานะที่ส่งมอบได้ (ล็อกไว้กันการเปลี่ยนสถานะพร้อมกัน)
	var requestStatus, loanType, requesterID string
	if err := tx.QueryRow(ctx,
		`SELECT status, loan_type, requester_id FROM borrow_request WHERE id = $1 FOR UPDATE`, rec.RequestID,
	).Scan(&requestStatus, &loanType, &requesterID); err != nil {
		return HandoverResult{}, err
	}
	if !slices.Contains(readyStatuses, requestStatus) {
		return HandoverResult{}, ErrRequestNotReady
	}
	// FR-HND-02: ผู้ยืมที่ยืนยันตัวตนต้องเป็นผู้ส่งคำขอ (requester_id คือรหัสจากการล็อกอิน ASM-BRQ-06 ของ 002)
	if rec.BorrowerUserID != requesterID {
		return HandoverResult{}, ErrBorrowerMismatch
	}

	// FR-HND-07: สร้าง Loan ด้วยคอลัมน์เดิมของ 001 (DM-05) ไม่แก้โครงตาราง loan
	if err := tx.QueryRow(ctx,
		`INSERT INTO loan (item_id, borrower_id, loan_type, due_at) VALUES ($1, $2, $3, $4) RETURNING id`,
		res.ItemID, rec.BorrowerUserID, loanType, rec.DueAt,
	).Scan(&res.LoanID); err != nil {
		return HandoverResult{}, fmt.Errorf("insert loan: %w", err)
	}

	// FR-HND-02, 03, 05: ข้อมูลการส่งมอบอยู่ในตาราง handover ของ 005
	if _, err := tx.Exec(ctx,
		`INSERT INTO handover (loan_id, request_id, handed_over_at, handed_over_by,
		                       borrower_student_id, borrower_first_name, borrower_last_name, identity_method, card_checked)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		res.LoanID, rec.RequestID, rec.HandedOverAt, rec.HandedOverBy,
		rec.BorrowerStudentID, rec.BorrowerFirstName, rec.BorrowerLastName, rec.IdentityMethod, rec.CardChecked,
	); err != nil {
		return HandoverResult{}, fmt.Errorf("insert handover: %w", err)
	}

	// FR-HND-07, CON-DATA-01, DOM-REQ-01: อุปกรณ์เปลี่ยนจาก "ว่าง" เป็น "ถูกยืม" ที่ฟีเจอร์นี้เท่านั้น
	if _, err := tx.Exec(ctx,
		`UPDATE equipment_item SET status = $2, updated_at = now() WHERE id = $1`, res.ItemID, itemStatusBorrowed,
	); err != nil {
		return HandoverResult{}, fmt.Errorf("update item: %w", err)
	}

	// FR-HND-07: เปลี่ยนเฉพาะคำขอนี้เป็น "ส่งมอบแล้ว" (ASM-HND-08 ไม่แตะคำขออื่น)
	if _, err := tx.Exec(ctx,
		`UPDATE borrow_request SET status = $2 WHERE id = $1`, rec.RequestID, requestStatusHandedOver,
	); err != nil {
		return HandoverResult{}, fmt.Errorf("update request: %w", err)
	}

	// FR-HND-08, DOM-AUDIT-01: audit log ของการส่งมอบสำเร็จ เขียนใน transaction เดียวกัน (writer.go ของ 001 รับ tx ได้)
	// ถ้าเขียน audit ไม่ได้ การส่งมอบจะไม่ถูกบันทึก จึงไม่มีการส่งมอบที่สำเร็จโดยไม่มี audit log
	requestID := rec.RequestID
	if _, err := audit.Write(ctx, tx, audit.Entry{
		ActorID:   rec.HandedOverBy,
		ActorRole: string(auth.RoleStaff), // CON-HND-01 ผู้บันทึกการส่งมอบเป็นเจ้าหน้าที่เท่านั้น
		Action:    auditActionHandover,
		AssetCode: res.AssetCode, // CON-DATA-03
		RequestID: &requestID,
		Outcome:   audit.OutcomeOK,
	}); err != nil {
		return HandoverResult{}, fmt.Errorf("write audit log: %w", err)
	}

	// FR-HND-06: เขียนรูปลงดิสก์ก่อน commit แล้วบันทึก path ผูกกับการยืมนี้
	if photos != nil {
		stored, err := photos.SavePhotos(res.LoanID)
		if err != nil {
			return HandoverResult{}, fmt.Errorf("save photos: %w", err)
		}
		for _, p := range stored {
			if _, err := tx.Exec(ctx,
				`INSERT INTO handover_photo (loan_id, file_path, content_type, size_bytes) VALUES ($1, $2, $3, $4)`,
				res.LoanID, p.Path, p.ContentType, p.Size,
			); err != nil {
				photos.RemovePhotos(res.LoanID)
				return HandoverResult{}, fmt.Errorf("insert handover_photo: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		if photos != nil {
			photos.RemovePhotos(res.LoanID)
		}
		return HandoverResult{}, err
	}
	return res, nil
}

// Querier คือส่วนของ pgxpool.Pool ที่ใช้อ่านข้อมูล
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

// ReadyRequestDTO คือคำขอ 1 รายการในรายการรอส่งมอบ (FR-HND-01, plan.md ข้อ 4)
type ReadyRequestDTO struct {
	RequestID          int       `json:"requestId"`
	AssetCode          string    `json:"assetCode"` // CON-DATA-03
	ModelName          string    `json:"modelName"`
	RequesterStudentID string    `json:"requesterStudentId"` // รหัสผู้ใช้ของผู้ส่งคำขอ (borrow_request.requester_id)
	SubmittedAt        time.Time `json:"submittedAt"`
}

// ListReadyRequests คืนคำขอยืมทั้งหมดที่อยู่ในสถานะที่ส่งมอบได้ (อนุมัติแล้วและยังไม่ส่งมอบ) (FR-HND-01, ASM-HND-05)
// เรียงตามเวลาที่ส่งคำขอ เก่าก่อน (ลำดับการแสดงทีมเลือกเอง ไม่ได้มาจาก spec)
func ListReadyRequests(ctx context.Context, q Querier, readyStatuses []string) ([]ReadyRequestDTO, error) {
	rows, err := q.Query(ctx, `
SELECT r.id, i.asset_code, m.name, r.requester_id, r.submitted_at
FROM borrow_request r
JOIN equipment_item i ON i.id = r.item_id
JOIN equipment_model m ON m.id = i.model_id
WHERE r.status = ANY($1)
ORDER BY r.submitted_at, r.id`, readyStatuses)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (ReadyRequestDTO, error) {
		var d ReadyRequestDTO
		err := row.Scan(&d.RequestID, &d.AssetCode, &d.ModelName, &d.RequesterStudentID, &d.SubmittedAt)
		return d, err
	})
}

// RowQuerier คือส่วนของ pgxpool.Pool ที่ใช้อ่านแถวเดียว
type RowQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// GetRequester คืนรหัสผู้ใช้ของผู้ส่งคำขอ (borrow_request.requester_id) ใช้เทียบกับผู้ยืมที่ยืนยันตัวตน (FR-HND-02, ASM-HND-04)
// found = false ถ้าไม่มีคำขอรหัสนี้
func GetRequester(ctx context.Context, q RowQuerier, requestID int) (requesterID string, found bool, err error) {
	err = q.QueryRow(ctx, `SELECT requester_id FROM borrow_request WHERE id = $1`, requestID).Scan(&requesterID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return requesterID, true, nil
}
