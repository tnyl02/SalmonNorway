package handover

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/config"
)

// insertItem เพิ่มอุปกรณ์รายชิ้นของรุ่น 1 ใน transaction ของ test คืน id
func insertItem(t *testing.T, tx pgx.Tx, assetCode, status string) int {
	t.Helper()
	var id int
	if err := tx.QueryRow(context.Background(),
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, $2) RETURNING id`,
		assetCode, status).Scan(&id); err != nil {
		t.Fatalf("เพิ่มอุปกรณ์ไม่สำเร็จ: %v", err)
	}
	return id
}

// txRepo ใช้ transaction ของ test แทน pool: Begin ใน pgx.Tx คือ savepoint ทุกอย่างจึงถูกยกเลิกเมื่อ test จบ
type txRepo struct{ tx pgx.Tx }

func (r txRepo) Begin(ctx context.Context) (pgx.Tx, error) { return r.tx.Begin(ctx) }

func sampleRecord(requestID int) HandoverRecord {
	now := time.Now()
	return HandoverRecord{
		RequestID:         requestID,
		BorrowerUserID:    "student-001",
		BorrowerStudentID: "student-001",
		BorrowerFirstName: "ชื่อ",
		BorrowerLastName:  "นามสกุล",
		IdentityMethod:    "auth_api",
		CardChecked:       true,
		HandedOverBy:      "staff-001",
		HandedOverAt:      now,
		DueAt:             now.Add(time.Hour),
	}
}

// T-06 เสร็จเมื่อ: ล็อกอุปกรณ์ สร้าง Loan ที่ returned_at ว่าง และแถว handover คู่กัน อุปกรณ์เป็น borrowed
// คำขอเป็น handed_over ทั้งหมดใน transaction เดียว ไม่แตะคำขออื่น (FR-HND-04, FR-HND-07, DOM-REQ-01)
func TestRepository_CreateHandover(t *testing.T) {
	conn := testDB(t)
	ctx := context.Background()
	ready := config.HandoverReadyStatuses()

	t.Run("ส่งมอบสำเร็จ", func(t *testing.T) {
		tx := rollbackTx(t, conn)
		itemID := insertItem(t, tx, "T06-OK-0001", "available")
		requestID := insertRequest(t, tx, "student-001", itemID, "approved")
		otherID := insertRequest(t, tx, "student-002", itemID, "rejected") // คำขออื่นของชิ้นเดียวกัน

		res, err := NewRepository(txRepo{tx}).CreateHandover(ctx, sampleRecord(requestID), ready, nil)
		if err != nil {
			t.Fatalf("CreateHandover error = %v", err)
		}
		if res.ItemID != itemID || res.AssetCode != "T06-OK-0001" || res.LoanID == 0 {
			t.Fatalf("ผลลัพธ์ = %+v", res)
		}

		var loanItem int
		var borrower string
		var returnedAt *time.Time
		if err := tx.QueryRow(ctx, `SELECT item_id, borrower_id, returned_at FROM loan WHERE id = $1`, res.LoanID).
			Scan(&loanItem, &borrower, &returnedAt); err != nil {
			t.Fatal(err)
		}
		if loanItem != itemID || borrower != "student-001" || returnedAt != nil {
			t.Errorf("loan = item %d borrower %q returned_at %v", loanItem, borrower, returnedAt)
		}

		var hoRequest int
		var method string
		if err := tx.QueryRow(ctx, `SELECT request_id, identity_method FROM handover WHERE loan_id = $1`, res.LoanID).
			Scan(&hoRequest, &method); err != nil {
			t.Fatalf("ไม่พบแถว handover: %v", err)
		}
		if hoRequest != requestID || method != "auth_api" {
			t.Errorf("handover = request %d method %q", hoRequest, method)
		}

		assertItemStatus(t, tx, itemID, "borrowed")
		assertRequestStatus(t, tx, requestID, "handed_over")
		assertRequestStatus(t, tx, otherID, "rejected") // ASM-HND-08 ไม่แตะคำขออื่น
	})

	t.Run("อุปกรณ์ไม่ว่าง ไม่มีอะไรเปลี่ยน", func(t *testing.T) {
		tx := rollbackTx(t, conn)
		itemID := insertItem(t, tx, "T06-NA-0001", "unavailable")
		requestID := insertRequest(t, tx, "student-001", itemID, "approved")

		_, err := NewRepository(txRepo{tx}).CreateHandover(ctx, sampleRecord(requestID), ready, nil)
		if !errors.Is(err, ErrItemNotAvailable) {
			t.Fatalf("error = %v ต้องเป็น ErrItemNotAvailable", err)
		}
		assertNoLoan(t, tx, itemID)
		assertItemStatus(t, tx, itemID, "unavailable")
		assertRequestStatus(t, tx, requestID, "approved")
	})

	t.Run("คำขอไม่อยู่ในสถานะที่ส่งมอบได้ ไม่มีอะไรเปลี่ยน", func(t *testing.T) {
		tx := rollbackTx(t, conn)
		itemID := insertItem(t, tx, "T06-NR-0001", "available")
		requestID := insertRequest(t, tx, "student-001", itemID, "requested")

		_, err := NewRepository(txRepo{tx}).CreateHandover(ctx, sampleRecord(requestID), ready, nil)
		if !errors.Is(err, ErrRequestNotReady) {
			t.Fatalf("error = %v ต้องเป็น ErrRequestNotReady", err)
		}
		assertNoLoan(t, tx, itemID)
		assertItemStatus(t, tx, itemID, "available")
		assertRequestStatus(t, tx, requestID, "requested")
	})

	t.Run("ไม่มีคำขอ", func(t *testing.T) {
		tx := rollbackTx(t, conn)
		_, err := NewRepository(txRepo{tx}).CreateHandover(ctx, sampleRecord(-1), ready, nil)
		if !errors.Is(err, ErrRequestNotFound) {
			t.Fatalf("error = %v ต้องเป็น ErrRequestNotFound", err)
		}
	})
}

func assertItemStatus(t *testing.T, tx pgx.Tx, itemID int, want string) {
	t.Helper()
	var got string
	if err := tx.QueryRow(context.Background(), `SELECT status FROM equipment_item WHERE id = $1`, itemID).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น %q", got, want)
	}
}

func assertRequestStatus(t *testing.T, tx pgx.Tx, requestID int, want string) {
	t.Helper()
	var got string
	if err := tx.QueryRow(context.Background(), `SELECT status FROM borrow_request WHERE id = $1`, requestID).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("สถานะคำขอ %d = %q ต้องเป็น %q", requestID, got, want)
	}
}

func assertNoLoan(t *testing.T, tx pgx.Tx, itemID int) {
	t.Helper()
	var n int
	if err := tx.QueryRow(context.Background(), `SELECT count(*) FROM loan WHERE item_id = $1`, itemID).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("มี Loan %d รายการ ต้องไม่มี", n)
	}
}

// T-18 เสร็จเมื่อ: เจ้าหน้าที่ 2 คนส่งมอบอุปกรณ์ชิ้นเดียวกันพร้อมกัน สำเร็จได้ 1 อีกคนได้ ErrItemNotAvailable
// และมี Loan 1 แถว (FR-HND-04 ยืนยัน SELECT ... FOR UPDATE, DOM-REQ-01)
// ใช้ connection pool จริงและ commit จริง เพราะ transaction ที่ยกเลิกทิ้งทดสอบการแย่งกันข้าม connection ไม่ได้ ข้อมูลที่สร้างลบทิ้งตอนจบ
func TestHandover_ConcurrentSameItem(t *testing.T) {
	testDB(t) // เตรียมฐานข้อมูลและ dbURL
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close) // ลงทะเบียนก่อน จึงทำงานหลังการลบข้อมูล (Cleanup ทำงานย้อนลำดับ)

	// คำขออนุมัติแล้ว 2 รายการของอุปกรณ์ชิ้นเดียวกัน คนละผู้ขอ (ข้อมูลตั้งเองเพื่อบังคับให้แย่งกัน)
	var itemID, reqA, reqB int
	if err := pool.QueryRow(ctx,
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ('T18-CC-0001', 1, 'available') RETURNING id`,
	).Scan(&itemID); err != nil {
		t.Fatal(err)
	}
	const insertReq = `INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
		VALUES ($1, 'ผู้ขอทดสอบ', $1, $2, 'normal', 'study', now(), now(), 'approved') RETURNING id`
	if err := pool.QueryRow(ctx, insertReq, "student-001", itemID).Scan(&reqA); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, insertReq, "student-002", itemID).Scan(&reqB); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `DELETE FROM handover WHERE request_id IN ($1, $2)`, reqA, reqB); err != nil {
			t.Errorf("ลบข้อมูลทดสอบไม่สำเร็จ: %v", err)
		}
		if _, err := pool.Exec(ctx, `DELETE FROM loan WHERE item_id = $1`, itemID); err != nil {
			t.Errorf("ลบข้อมูลทดสอบไม่สำเร็จ: %v", err)
		}
		// audit_log เป็น append-only และอ้าง borrow_request (T-26) คำขอที่มี audit log แล้วลบไม่ได้ ตั้งเป็น cancelled แทน
		for _, q := range []string{
			`UPDATE borrow_request SET status = 'cancelled'
			 WHERE id IN ($1, $2) AND EXISTS (SELECT 1 FROM audit_log a WHERE a.request_id = borrow_request.id)`,
			`DELETE FROM borrow_request
			 WHERE id IN ($1, $2) AND NOT EXISTS (SELECT 1 FROM audit_log a WHERE a.request_id = borrow_request.id)`,
		} {
			if _, err := pool.Exec(ctx, q, reqA, reqB); err != nil {
				t.Errorf("ลบข้อมูลทดสอบไม่สำเร็จ: %v", err)
			}
		}
		if _, err := pool.Exec(ctx,
			`DELETE FROM equipment_item WHERE id = $1 AND NOT EXISTS (SELECT 1 FROM borrow_request WHERE item_id = $1)`, itemID); err != nil {
			t.Errorf("ลบข้อมูลทดสอบไม่สำเร็จ: %v", err)
		}
	})

	repo := NewRepository(pool)
	ready := config.HandoverReadyStatuses()
	start := make(chan struct{})
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for requestID, requester := range map[int]string{reqA: "student-001", reqB: "student-002"} {
		rec := sampleRecord(requestID)
		rec.BorrowerUserID, rec.BorrowerStudentID = requester, requester // ผู้ยืมตรงกับผู้ขอของแต่ละคำขอ (FR-HND-02)
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := repo.CreateHandover(ctx, rec, ready, nil)
			errs <- err
		}()
	}
	close(start)
	wg.Wait()
	close(errs)

	var ok, notAvailable int
	for err := range errs {
		switch {
		case err == nil:
			ok++
		case errors.Is(err, ErrItemNotAvailable):
			notAvailable++
		default:
			t.Errorf("error ที่ไม่คาดคิด: %v", err)
		}
	}
	if ok != 1 || notAvailable != 1 {
		t.Fatalf("สำเร็จ %d ครั้ง ได้ ErrItemNotAvailable %d ครั้ง ต้องเป็น 1 และ 1", ok, notAvailable)
	}

	var loans int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM loan WHERE item_id = $1`, itemID).Scan(&loans); err != nil {
		t.Fatal(err)
	}
	if loans != 1 {
		t.Errorf("มี Loan %d แถว ต้องเป็น 1", loans)
	}
}
