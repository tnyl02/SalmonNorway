package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-11-1 ถึง TC-BRQ-11-4 (AC-BRQ-11, FR-BRQ-05 ข)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/config"
)

// errorCodeOf อ่านรหัส error จากคำตอบตามรูปแบบกลาง {code, message} (CON-TECH-03)
func errorCodeOf(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var got ErrorDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v body=%s", err, w.Body.String())
	}
	return got.Code
}

// insertExistingRequest ใส่คำขอที่มีอยู่ก่อนแล้วลงฐานข้อมูลตรง ๆ
// ใช้เตรียมสถานะที่ 002 ตั้งเองไม่ได้ (approved ของ 004, cancelled ของ 003, rejected ของ 004)
func insertExistingRequest(t *testing.T, pool *pgxpool.Pool, requesterID, assetCode, status string) {
	t.Helper()
	start := bangkokToday(9, 0)
	if _, err := pool.Exec(context.Background(), `
		INSERT INTO borrow_request
			(requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
		SELECT $1, 'ผู้ขอทดสอบ', '660710000', id, $2, 'study', $3, $4, $5
		FROM equipment_item WHERE asset_code = $6`,
		requesterID, config.LoanTypeNormal, start, bangkokToday(15, 0), status, assetCode); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM borrow_request WHERE requester_id = $1`, requesterID)
	})
}

// TC-BRQ-11-1 (ทางปกติ) คำขอเดิมของตนสถานะ requested แล้วส่งซ้ำ
func TestTC_BRQ_11_1_DuplicateRequested(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-11-1"

	// Given นักศึกษาคนหนึ่งมีคำขอของอุปกรณ์ชิ้นหนึ่ง สถานะ requested
	insertTempItem(t, assetCode)
	decodeCreated(t, postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode)))

	// When ส่งคำขอของชิ้นเดิมซ้ำ
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 409 DUPLICATE_REQUEST
	if w.Code != http.StatusConflict {
		t.Fatalf("(1) ส่งซ้ำต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if code := errorCodeOf(t, w); code != "DUPLICATE_REQUEST" {
		t.Errorf("(1) รหัสต้องเป็น DUPLICATE_REQUEST ตาม CON-TECH-03 แต่ได้ %q", code)
	}
	// Then (2) มีคำขอของคนนี้บนชิ้นนี้เพียง 1 แถว
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(2) ต้องมีคำขอเพียง 1 ใบ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-11-2 (ขอบ) คำขอเดิมสถานะ approved ซึ่ง AC ระบุว่ายังนับเป็นคำขอค้าง
func TestTC_BRQ_11_2_DuplicateApproved(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-11-2"

	// Given คำขอเดิมของ student-001 ถูกอนุมัติแล้ว (สถานะที่ 004 เป็นผู้ตั้ง)
	insertTempItem(t, assetCode)
	insertExistingRequest(t, pool, "student-001", assetCode, config.StatusApproved)

	// When ส่งคำขอของชิ้นเดิมซ้ำ
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 409 DUPLICATE_REQUEST
	if w.Code != http.StatusConflict {
		t.Fatalf("(1) ส่งซ้ำขณะคำขอเดิมถูกอนุมัติต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if code := errorCodeOf(t, w); code != "DUPLICATE_REQUEST" {
		t.Errorf("(1) รหัสต้องเป็น DUPLICATE_REQUEST แต่ได้ %q", code)
	}
	// Then (2) ยังมีเพียง 1 แถว
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(2) ต้องมีคำขอเพียง 1 ใบ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-11-3 (ทางผิด) ทำให้เงื่อนไข "ยังค้าง" ไม่เป็นจริง คือคำขอเดิมจบไปแล้ว
func TestTC_BRQ_11_3_AfterCancelledOrRejected(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)

	for _, tc := range []struct {
		name      string
		assetCode string
		status    string
	}{
		{"คำขอเดิมถูกยกเลิก", "TC-BRQ-11-3a", config.StatusCancelled},
		{"คำขอเดิมถูกปฏิเสธ", "TC-BRQ-11-3b", config.StatusRejected},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Given คำขอเดิมของชิ้นนี้จบไปแล้ว
			insertTempItem(t, tc.assetCode)
			insertExistingRequest(t, pool, "student-001", tc.assetCode, tc.status)

			// When ส่งคำขอของชิ้นเดิมใหม่
			w := postRequest(t, srv, tokenStudentDept1, bodyFor(tc.assetCode))

			// Then (1) ได้ 201
			if w.Code != http.StatusCreated {
				t.Fatalf("(1) คำขอเดิมจบแล้วต้องส่งใหม่ได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
			}
			decodeCreated(t, w)
			// Then (2) มีแถวใหม่เพิ่ม รวมของเดิมเป็น 2 ใบ
			if n := countRequestsOnItem(t, pool, "student-001", tc.assetCode); n != 2 {
				t.Errorf("(2) ต้องมีคำขอ 2 ใบ (ใบเดิมที่จบแล้วกับใบใหม่) แต่พบ %d ใบ", n)
			}
		})
	}
}

// TC-BRQ-11-4 (ทางผิด) ผู้ส่งซ้ำเป็นคนละคนกับเจ้าของคำขอเดิม
func TestTC_BRQ_11_4_OtherStudentNotDuplicate(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-11-4"

	// Given อุปกรณ์ชิ้นนี้มีคำขอค้างของนักศึกษาอีกคนอยู่แล้ว
	insertTempItem(t, assetCode)
	insertExistingRequest(t, pool, "student-อื่น", assetCode, config.StatusRequested)

	// When student-001 ส่งคำขอของชิ้นเดิม
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 201 เพราะกฎคำขอซ้ำนับเฉพาะคำขอของตนเอง (FR-BRQ-05 ข, CR-004)
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) คำขอของคนอื่นต้องไม่ทำให้ถูกปฏิเสธ ต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(1) student-001 ต้องมีคำขอของชิ้นนี้ 1 ใบ แต่พบ %d ใบ", n)
	}
}
