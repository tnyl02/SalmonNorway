// test จากตาราง specs/004-approve-request/test-cases.md แถว TC-APV-01-x สถานะ "ใช้ได้" ห้ามแก้ให้ผ่าน
// ถ้าไม่ผ่าน ต้องหาว่าโค้ดผิด หรือ AC ผิด (ถ้า AC ผิด แก้ spec ก่อน)
package approve

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/auth"
)

// tcServer ประกอบ handler กับ auth.Middleware แบบเดียวกับ main.go (IF-IDP-01)
// clock ไม่ nil = ใช้นาฬิกาจำลองของ service (ตารางบอกว่า test ที่ขึ้นกับเวลาใช้นาฬิกาจำลอง)
func tcServer(tx pgx.Tx, clock *fakeClock) http.Handler {
	svc := NewService(NewRepository(tx))
	if clock != nil {
		svc.now = clock.Now
	}
	return auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes())
}

// tcApprove เรียก POST /v1/approvals/requests/{id}/approve พร้อม loadedAt (plan ข้อ 4)
func tcApprove(srv http.Handler, token string, id int) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"loadedAt":%q}`, time.Now().UTC().Format(time.RFC3339))
	return postJSON(srv, token, fmt.Sprintf("%srequests/%d/approve", PathPrefix, id), body)
}

// tcRequestStatus คืน borrow_request.status ของคำขอ
func tcRequestStatus(t *testing.T, tx pgx.Tx, id int) string {
	t.Helper()
	var s string
	if err := tx.QueryRow(context.Background(), `SELECT status FROM borrow_request WHERE id = $1`, id).Scan(&s); err != nil {
		t.Fatal(err)
	}
	return s
}

// tcItemStatus คืน equipment_item.status ของอุปกรณ์รายชิ้น
func tcItemStatus(t *testing.T, tx pgx.Tx, assetCode string) string {
	t.Helper()
	var s string
	if err := tx.QueryRow(context.Background(), `SELECT status FROM equipment_item WHERE asset_code = $1`, assetCode).Scan(&s); err != nil {
		t.Fatal(err)
	}
	return s
}

// tcStartAt คืนเวลาเริ่มยืมของคำขอ ใช้ตั้งนาฬิกาจำลองเทียบกับ 16:00 น. ของวันที่เริ่มยืม
func tcStartAt(t *testing.T, tx pgx.Tx, id int) time.Time {
	t.Helper()
	var at time.Time
	if err := tx.QueryRow(context.Background(), `SELECT start_at FROM borrow_request WHERE id = $1`, id).Scan(&at); err != nil {
		t.Fatal(err)
	}
	return at
}

// tcError อ่าน body {code, message} (CON-TECH-03)
func tcError(t *testing.T, w *httptest.ResponseRecorder) ErrorDTO {
	t.Helper()
	var e ErrorDTO
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatalf("body ต้องเป็น {code, message} แต่ได้ %s", w.Body.String())
	}
	return e
}

// TC-APV-01-1 (AC-APV-01 ปกติ)
func TestTC_APV_01_1_ApproveFirstRequest(t *testing.T) {
	// Given staff-001 ยืนยันตัวตนแล้ว คำขอยืมปกติ requested ของ student-001 บน AP-TEST-0002 ("ว่าง") เป็นคำขอแรกของชิ้นนั้น
	//       เวลาก่อน 16:00 น. ของวันที่เริ่มยืม
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	clock := &fakeClock{t: approvalCutoff(tcStartAt(t, tx, id)).Add(-time.Hour)}
	srv := tcServer(tx, clock)

	// When staff-001 เรียก POST /v1/approvals/requests/{id}/approve พร้อม loadedAt
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 200 คำตอบมี requestId ของคำขอนั้นและ status approved
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var got ResultDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.RequestID != id || got.Status != "approved" {
		t.Errorf("(1) คำตอบต้องเป็น requestId %d status approved แต่ได้ %+v", id, got)
	}
	// Then (2) ในฐานข้อมูล borrow_request.status เป็น approved
	if s := tcRequestStatus(t, tx, id); s != "approved" {
		t.Errorf("(2) borrow_request.status ต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (3) equipment_item.status ของ AP-TEST-0002 ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0002"); s != "available" {
		t.Errorf("(3) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
	// Then (4) มีแถว approval_decision ของคำขอนั้น result approved ผู้ดำเนินการ staff-001
	var result, actor string
	if err := tx.QueryRow(context.Background(),
		`SELECT result, actor_id FROM approval_decision WHERE request_id = $1`, id).Scan(&result, &actor); err != nil {
		t.Fatalf("(4) ต้องมีแถว approval_decision ของคำขอนั้น: %v", err)
	}
	if result != "approved" || actor != "staff-001" {
		t.Errorf("(4) approval_decision ต้องเป็น approved โดย staff-001 แต่ได้ %q โดย %q", result, actor)
	}
}

// TC-APV-01-4 (AC-APV-01 ขอบ)
func TestTC_APV_01_4_ApproveAt1559(t *testing.T) {
	// Given เหมือน TC-APV-01-1 แต่ตั้งนาฬิกาจำลองเป็น 15:59:59 ของวันที่เริ่มยืม
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	cutoff := approvalCutoff(tcStartAt(t, tx, id))
	clock := &fakeClock{t: cutoff.Add(-time.Second)}
	if lt := localTime(clock.t); lt.Hour() != 15 || lt.Minute() != 59 || lt.Second() != 59 {
		t.Fatalf("นาฬิกาจำลองต้องเป็น 15:59:59 เวลาไทย แต่ได้ %v", lt)
	}
	srv := tcServer(tx, clock)

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) อนุมัติสำเร็จ สถานะคำขอเป็น approved
	if w.Code != http.StatusOK {
		t.Errorf("(1) 15:59:59 ต้องอนุมัติได้ แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, id); s != "approved" {
		t.Errorf("(1) สถานะคำขอต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (2) equipment_item.status ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0002"); s != "available" {
		t.Errorf("(2) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
}

// TC-APV-01-5 (AC-APV-01 ขอบ)
func TestTC_APV_01_5_ApproveAt1600(t *testing.T) {
	// Given เหมือน TC-APV-01-1 แต่ตั้งนาฬิกาจำลองเป็น 16:00:00 ตรงของวันที่เริ่มยืม
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	clock := &fakeClock{t: approvalCutoff(tcStartAt(t, tx, id))}
	if lt := localTime(clock.t); lt.Hour() != 16 || lt.Minute() != 0 || lt.Second() != 0 {
		t.Fatalf("นาฬิกาจำลองต้องเป็น 16:00:00 เวลาไทย แต่ได้ %v", lt)
	}
	srv := tcServer(tx, clock)

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 409 APPROVAL_CUTOFF_PASSED ข้อความ "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้" (Q-APV-26)
	if w.Code != http.StatusConflict {
		t.Errorf("(1) 16:00:00 ตรงต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	} else if e := tcError(t, w); e.Code != "APPROVAL_CUTOFF_PASSED" || e.Message != "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้" {
		t.Errorf("(1) ต้องได้ APPROVAL_CUTOFF_PASSED พร้อมข้อความตาม FR-APV-01 (ซ) แต่ได้ %+v", e)
	}
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	// Then (3) equipment_item.status ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0002"); s != "available" {
		t.Errorf("(3) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
}

// TC-APV-01-6 (AC-APV-01 ขอบ)
func TestTC_APV_01_6_NextBecomesFirst(t *testing.T) {
	// Given AP-TEST-0001 คำขอ A (คำขอแรก) ถูกเจ้าหน้าที่ปฏิเสธไปแล้ว เหลือ B (ยืมปกติ) และ C (ค้างคืน) เป็น requested
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	clock := &fakeClock{t: approvalCutoff(tcStartAt(t, tx, b)).Add(-time.Hour)}
	srv := tcServer(tx, clock)
	rw := postJSON(srv, tokenStaffDept1, fmt.Sprintf("%srequests/%d/reject", PathPrefix, a), `{"reason":"อุปกรณ์ต้องใช้สอบ"}`)
	if rw.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: ปฏิเสธ A ต้องได้ 200 แต่ได้ %d body=%s", rw.Code, rw.Body.String())
	}
	if tcRequestStatus(t, tx, b) != "requested" || tcRequestStatus(t, tx, c) != "requested" {
		t.Fatal("เตรียมข้อมูล: B และ C ต้องยังเป็น requested")
	}

	// When staff-001 อนุมัติ B
	w := tcApprove(srv, tokenStaffDept1, b)

	// Then (1) อนุมัติสำเร็จ B เป็น approved เพราะคำขอถัดไปกลายเป็นคำขอแรก
	if w.Code != http.StatusOK {
		t.Errorf("(1) อนุมัติ B ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, b); s != "approved" {
		t.Errorf("(1) B ต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (2) equipment_item.status ของ AP-TEST-0001 ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0001"); s != "available" {
		t.Errorf("(2) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
	// Then (3) C เป็น rejected
	if s := tcRequestStatus(t, tx, c); s != "rejected" {
		t.Errorf("(3) C ต้องเป็น rejected แต่ได้ %q", s)
	}
}

// TC-APV-01-7 (AC-APV-01 ขอบ)
func TestTC_APV_01_7_OvernightDoesNotBlock(t *testing.T) {
	// Given อุปกรณ์ชิ้นหนึ่งที่ "ว่าง" (AP-TEST-0002) มีคำขอค้างคืน requested ที่ส่งก่อน แล้วจึงมีคำขอยืมปกติ requested
	tx := devTx(t)
	normal := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	var overnight int
	if err := tx.QueryRow(context.Background(), `
		INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose,
		                            start_at, end_at, status, submitted_at)
		SELECT 'student-003', 'นักศึกษาตัวอย่าง ค', 'student-003', id, 'overnight', 'activity',
		       now(), now() + interval '1 day', 'requested', now() - interval '2 hours'
		FROM equipment_item WHERE asset_code = 'AP-TEST-0002'
		RETURNING id`).Scan(&overnight); err != nil {
		t.Fatal(err)
	}
	clock := &fakeClock{t: approvalCutoff(tcStartAt(t, tx, normal)).Add(-time.Hour)}
	srv := tcServer(tx, clock)

	// When staff-001 อนุมัติคำขอยืมปกตินั้น
	w := tcApprove(srv, tokenStaffDept1, normal)

	// Then (1) อนุมัติสำเร็จ เพราะคำขอค้างคืนไม่ขวางคิว
	if w.Code != http.StatusOK {
		t.Errorf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, normal); s != "approved" {
		t.Errorf("(1) คำขอยืมปกติต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (2) equipment_item.status ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0002"); s != "available" {
		t.Errorf("(2) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
	// Then (3) คำขอค้างคืนเป็น rejected
	if s := tcRequestStatus(t, tx, overnight); s != "rejected" {
		t.Errorf("(3) คำขอค้างคืนต้องเป็น rejected แต่ได้ %q", s)
	}
}

// TC-APV-01-8 (AC-APV-01 ผิด)
func TestTC_APV_01_8_NoTokenNotApproved(t *testing.T) {
	// Given ไม่มี token คำขอเดียวกับ TC-APV-01-1
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, nil)

	// When เรียก POST .../approve
	w := tcApprove(srv, "", id)

	// Then (1) ได้ 401 UNAUTHENTICATED
	if w.Code != http.StatusUnauthorized {
		t.Errorf("(1) ต้องได้ 401 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if e := tcError(t, w); e.Code != "UNAUTHENTICATED" {
		t.Errorf("(1) code ต้องเป็น UNAUTHENTICATED แต่ได้ %+v", e)
	}
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-01-9 (AC-APV-01 ผิด)
func TestTC_APV_01_9_NotFirstRejected(t *testing.T) {
	// Given คำขอ B ของ AP-TEST-0001 (A ยังเป็น requested)
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	clock := &fakeClock{t: approvalCutoff(tcStartAt(t, tx, b)).Add(-time.Hour)}
	srv := tcServer(tx, clock)

	// When staff-001 อนุมัติ B
	w := tcApprove(srv, tokenStaffDept1, b)

	// Then (1) ได้ 409 NOT_FIRST_IN_QUEUE ข้อความ "ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน"
	if w.Code != http.StatusConflict {
		t.Errorf("(1) ต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if e := tcError(t, w); e.Code != "NOT_FIRST_IN_QUEUE" || e.Message != "ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน" {
		t.Errorf("(1) ต้องได้ NOT_FIRST_IN_QUEUE พร้อมข้อความตาม FR-APV-01 (ซ) แต่ได้ %+v", e)
	}
	// Then (2) A และ B ยังเป็น requested
	if s := tcRequestStatus(t, tx, a); s != "requested" {
		t.Errorf("(2) A ต้องยังเป็น requested แต่ได้ %q", s)
	}
	if s := tcRequestStatus(t, tx, b); s != "requested" {
		t.Errorf("(2) B ต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-01-10 (AC-APV-01 ผิด)
func TestTC_APV_01_10_ItemNotAvailable(t *testing.T) {
	// Given คำขอแรกของ AP-TEST-0004 ("ถูกยืม")
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0004", "student-001")
	clock := &fakeClock{t: approvalCutoff(tcStartAt(t, tx, id)).Add(-time.Hour)}
	srv := tcServer(tx, clock)

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 409 ITEM_NOT_AVAILABLE ข้อความ "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้"
	if w.Code != http.StatusConflict {
		t.Errorf("(1) ต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if e := tcError(t, w); e.Code != "ITEM_NOT_AVAILABLE" || e.Message != "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้" {
		t.Errorf("(1) ต้องได้ ITEM_NOT_AVAILABLE พร้อมข้อความตาม FR-APV-01 (ซ) แต่ได้ %+v", e)
	}
	// Then (2) คำขอยังเป็น requested และ equipment_item.status ยังเป็น borrowed
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) คำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	if s := tcItemStatus(t, tx, "AP-TEST-0004"); s != "borrowed" {
		t.Errorf("(2) equipment_item.status ต้องยังเป็น borrowed แต่ได้ %q", s)
	}
}

// TC-APV-01-11 (AC-APV-01 ผิด)
func TestTC_APV_01_11_CancelledNotApproved(t *testing.T) {
	// Given คำขอ cancelled ของ AP-TEST-0005
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0005", "student-002")
	srv := tcServer(tx, nil)

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 409 REQUEST_CANCELLED ข้อความ "คำขอนี้ถูกผู้ขอยกเลิกแล้ว"
	if w.Code != http.StatusConflict {
		t.Errorf("(1) ต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if e := tcError(t, w); e.Code != "REQUEST_CANCELLED" || e.Message != "คำขอนี้ถูกผู้ขอยกเลิกแล้ว" {
		t.Errorf("(1) ต้องได้ REQUEST_CANCELLED พร้อมข้อความตาม FR-APV-01 (ซ) แต่ได้ %+v", e)
	}
	// Then (2) คำขอยังเป็น cancelled
	if s := tcRequestStatus(t, tx, id); s != "cancelled" {
		t.Errorf("(2) คำขอต้องยังเป็น cancelled แต่ได้ %q", s)
	}
}

// TC-APV-01-12 (AC-APV-01 ผิด) แก้ตาม spec v3 CR-009 2569-11-13
func TestTC_APV_01_12_OvernightWithoutAdvisor(t *testing.T) {
	// Given คำขอค้างคืน C ของ AP-TEST-0001 และไม่ส่ง advisorName
	tx := devTx(t)
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	srv := tcServer(tx, nil)

	// When staff-001 อนุมัติ C โดยไม่ระบุชื่ออาจารย์ (tcApprove ส่งแค่ loadedAt)
	w := tcApprove(srv, tokenStaffDept1, c)

	// Then (1) ได้ 422 VALIDATION_FAILED fields ["advisorName"] ข้อความ "กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ"
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("(1) ต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var e ValidationErrorDTO
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if e.Code != "VALIDATION_FAILED" || e.Message != "กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ" || len(e.Fields) != 1 || e.Fields[0] != "advisorName" {
		t.Errorf("(1) ต้องได้ VALIDATION_FAILED ช่อง advisorName พร้อมข้อความ \"กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ\" แต่ได้ %+v", e)
	}
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}
