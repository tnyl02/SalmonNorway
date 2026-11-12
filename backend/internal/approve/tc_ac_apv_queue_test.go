// test จากตาราง specs/004-approve-request/test-cases.md แถว TC-APV-03-x, 07-x, 08-x, 09-x, 16-x, 17-x สถานะ "ใช้ได้" ห้ามแก้ให้ผ่าน
// ถ้าไม่ผ่าน ต้องหาว่าโค้ดผิด หรือ AC ผิด (ถ้า AC ผิด แก้ spec ก่อน)
// ใช้ helper จาก tc_ac_apv_01_test.go, tc_ac_apv_02_06_11_test.go และ tc_ac_apv_05_12_test.go
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
)

const tcMsgSuperseded = "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว"
const tcMsgNotFirst = "ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน"
const tcMsgItemNotAvailable = "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้"
const tcMsgAlreadyDecided = "คำขอนี้ได้รับการพิจารณาแล้ว"
const tcMsgCutoff = "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้"

// tcAddRequest เพิ่มคำขอยืมปกติ requested ของผู้ขออีกคนบนอุปกรณ์ชิ้นเดียวกับคำขอ like ช่วงเวลาเดียวกัน ส่งเมื่อ ago ก่อน
func tcAddRequest(t *testing.T, tx pgx.Tx, like int, requester string, ago time.Duration) int {
	t.Helper()
	var id int
	if err := tx.QueryRow(context.Background(), `
		INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose,
		                            start_at, end_at, status, submitted_at)
		SELECT $2, 'นักศึกษาเพิ่ม '||$2, $2, r.item_id, 'normal', 'study', r.start_at, r.end_at, 'requested', $3
		FROM borrow_request r WHERE r.id = $1
		RETURNING id`, like, requester, time.Now().Add(-ago)).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

// tcSetStatus เปลี่ยนสถานะคำขอโดยตรง (เตรียมข้อมูล)
func tcSetStatus(t *testing.T, tx pgx.Tx, id int, status string) {
	t.Helper()
	q := `UPDATE borrow_request SET status = $2 WHERE id = $1`
	if status == "cancelled" {
		q = `UPDATE borrow_request SET status = $2, cancel_reason = 'ไม่ต้องใช้แล้ว', cancelled_at = now() WHERE id = $1`
	}
	if _, err := tx.Exec(context.Background(), q, id, status); err != nil {
		t.Fatal(err)
	}
}

// tcSetItemStatus เปลี่ยนสถานะอุปกรณ์รายชิ้นโดยตรง (เตรียมข้อมูล)
func tcSetItemStatus(t *testing.T, tx pgx.Tx, asset, status string) {
	t.Helper()
	if _, err := tx.Exec(context.Background(), `UPDATE equipment_item SET status = $2 WHERE asset_code = $1`, asset, status); err != nil {
		t.Fatal(err)
	}
}

// tcDecision คืนผล เหตุผล และผู้ดำเนินการของ approval_decision (ok = false ถ้าไม่มี)
func tcDecision(t *testing.T, tx pgx.Tx, id int) (result, reason, actor string, decidedAt time.Time, ok bool) {
	t.Helper()
	var r *string
	err := tx.QueryRow(context.Background(),
		`SELECT result, reason, actor_id, decided_at FROM approval_decision WHERE request_id = $1`, id).Scan(&result, &r, &actor, &decidedAt)
	if err == pgx.ErrNoRows {
		return "", "", "", time.Time{}, false
	}
	if err != nil {
		t.Fatal(err)
	}
	if r != nil {
		reason = *r
	}
	return result, reason, actor, decidedAt, true
}

// tcEmailKind คืน kind ของ approval_email ของคำขอ ("" ถ้าไม่มี)
func tcEmailKind(t *testing.T, tx pgx.Tx, id int) string {
	t.Helper()
	var kind string
	err := tx.QueryRow(context.Background(), `SELECT kind FROM approval_email WHERE request_id = $1`, id).Scan(&kind)
	if err == pgx.ErrNoRows {
		return ""
	}
	if err != nil {
		t.Fatal(err)
	}
	return kind
}

// tcApproveResult อ่านคำตอบการอนุมัติที่สำเร็จ
func tcApproveResult(t *testing.T, w *httptest.ResponseRecorder) ResultDTO {
	t.Helper()
	var got ResultDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("อ่านคำตอบไม่ได้: %s", w.Body.String())
	}
	return got
}

func tcContains(ids []int, id int) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

// tcAt คืนเวลา hh:mm:ss เวลาไทยของวันที่เริ่มยืมของคำขอ
func tcAt(t *testing.T, tx pgx.Tx, id int, h, m, s int) time.Time {
	lt := localTime(tcStartAt(t, tx, id))
	return time.Date(lt.Year(), lt.Month(), lt.Day(), h, m, s, 0, lt.Location())
}

// tcApproveWithLoaded อนุมัติพร้อม loadedAt ที่ระบุ
func tcApproveWithLoaded(srv http.Handler, id int, loadedAt time.Time) *httptest.ResponseRecorder {
	body := fmt.Sprintf(`{"loadedAt":%q}`, loadedAt.UTC().Format(time.RFC3339))
	return postJSON(srv, tokenStaffDept1, fmt.Sprintf("%srequests/%d/approve", PathPrefix, id), body)
}

// ---------- AC-APV-03 ----------

// TC-APV-03-1 (AC-APV-03 ปกติ)
func TestTC_APV_03_1_ApproveSupersedesOthers(t *testing.T) {
	// Given AP-TEST-0001 มี A และ B เป็น requested เวลาก่อน 16:00 น.
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติ A
	w := tcApprove(srv, tokenStaffDept1, a)

	// Then (1) ได้ 200 และ rejectedRequestIds มีรหัสของ B
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if got := tcApproveResult(t, w); !tcContains(got.RejectedRequestIDs, b) {
		t.Errorf("(1) rejectedRequestIds ต้องมี B (%d) แต่ได้ %v", b, got.RejectedRequestIDs)
	}
	// Then (2) A เป็น approved (3) B เป็น rejected
	if s := tcRequestStatus(t, tx, a); s != "approved" {
		t.Errorf("(2) A ต้องเป็น approved แต่ได้ %q", s)
	}
	if s := tcRequestStatus(t, tx, b); s != "rejected" {
		t.Errorf("(3) B ต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (4) approval_decision ของ B result rejected reason ตามที่กำหนด ผู้ดำเนินการ staff-001
	if res, reason, actor, _, ok := tcDecision(t, tx, b); !ok || res != "rejected" || reason != tcMsgSuperseded || actor != "staff-001" {
		t.Errorf("(4) B ต้องเป็น rejected / %q / staff-001 แต่ได้ %q / %q / %q (มีแถว %v)", tcMsgSuperseded, res, reason, actor, ok)
	}
	// Then (5) มีแถว approval_email ของ B แบบ rejected
	if k := tcEmailKind(t, tx, b); k != "rejected" {
		t.Errorf("(5) Email ของ B ต้องเป็นแบบ rejected แต่ได้ %q", k)
	}
	// Then (6) equipment_item.status ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0001"); s != "available" {
		t.Errorf("(6) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
}

// TC-APV-03-3 (AC-APV-03 ขอบ)
func TestTC_APV_03_3_SupersedesAllOthers(t *testing.T) {
	// Given AP-TEST-0001 มีคำขอยืมปกติ A, B และ D (ส่งหลัง B)
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	d := tcAddRequest(t, tx, a, "student-005", 5*time.Minute)
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติ A
	w := tcApprove(srv, tokenStaffDept1, a)
	if w.Code != http.StatusOK {
		t.Fatalf("อนุมัติ A ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}

	// Then (1) A เป็น approved
	if s := tcRequestStatus(t, tx, a); s != "approved" {
		t.Errorf("(1) A ต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (2) B และ D เป็น rejected ทั้งคู่ เหตุผลตามที่กำหนด
	for name, id := range map[string]int{"B": b, "D": d} {
		if s := tcRequestStatus(t, tx, id); s != "rejected" {
			t.Errorf("(2) %s ต้องเป็น rejected แต่ได้ %q", name, s)
		}
		if _, reason, _, _, _ := tcDecision(t, tx, id); reason != tcMsgSuperseded {
			t.Errorf("(2) เหตุผลของ %s ต้องเป็น %q แต่ได้ %q", name, tcMsgSuperseded, reason)
		}
	}
	// Then (3) rejectedRequestIds มีทั้ง B และ D
	got := tcApproveResult(t, w)
	if !tcContains(got.RejectedRequestIDs, b) || !tcContains(got.RejectedRequestIDs, d) {
		t.Errorf("(3) rejectedRequestIds ต้องมี B และ D แต่ได้ %v", got.RejectedRequestIDs)
	}
}

// TC-APV-03-4 (AC-APV-03 ขอบ)
func TestTC_APV_03_4_NoOtherRequests(t *testing.T) {
	// Given AP-TEST-0002 มีคำขอ requested ของ student-001 คำขอเดียว
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, id))

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 200 และ rejectedRequestIds เป็นอาร์เรย์ว่าง []
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if string(raw["rejectedRequestIds"]) != "[]" {
		t.Errorf("(1) rejectedRequestIds ต้องเป็น [] แต่ได้ %s", raw["rejectedRequestIds"])
	}
	// Then (2) ไม่มีแถว approval_decision ผล rejected เกิดขึ้น
	var n int
	if err := tx.QueryRow(context.Background(), `SELECT count(*) FROM approval_decision WHERE result = 'rejected'`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("(2) ต้องไม่มีผล rejected แต่มี %d แถว", n)
	}
}

// TC-APV-03-5 (AC-APV-03 ขอบ)
func TestTC_APV_03_5_CancelledOtherUntouched(t *testing.T) {
	// Given AP-TEST-0001 มี A (requested) และ B ที่ผู้ขอยกเลิกไปแล้ว
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	tcSetStatus(t, tx, b, "cancelled")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติ A
	w := tcApprove(srv, tokenStaffDept1, a)
	if w.Code != http.StatusOK {
		t.Fatalf("อนุมัติ A ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}

	// Then (1) A เป็น approved
	if s := tcRequestStatus(t, tx, a); s != "approved" {
		t.Errorf("(1) A ต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (2) B ยังเป็น cancelled และไม่มีแถว approval_decision ของ B
	if s := tcRequestStatus(t, tx, b); s != "cancelled" {
		t.Errorf("(2) B ต้องยังเป็น cancelled แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, b); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision ของ B แต่มี %d", n)
	}
	// Then (3) B ไม่อยู่ใน rejectedRequestIds
	if got := tcApproveResult(t, w); tcContains(got.RejectedRequestIDs, b) {
		t.Errorf("(3) B ต้องไม่อยู่ใน rejectedRequestIds แต่ได้ %v", got.RejectedRequestIDs)
	}
}

// TC-APV-03-6 (AC-APV-03 ขอบ)
func TestTC_APV_03_6_OtherItemsUntouched(t *testing.T) {
	// Given ข้อมูล approve_dev.sql (AP-TEST-0002, 0003, 0004 มีคำขอ requested ของตัวเอง)
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	others := []int{
		requestIDOf(t, tx, "AP-TEST-0002", "student-001"),
		requestIDOf(t, tx, "AP-TEST-0003", "student-002"),
		requestIDOf(t, tx, "AP-TEST-0004", "student-001"),
	}
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติ A ของ AP-TEST-0001
	w := tcApprove(srv, tokenStaffDept1, a)
	if w.Code != http.StatusOK {
		t.Fatalf("อนุมัติ A ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	got := tcApproveResult(t, w)

	for _, id := range others {
		// Then (1) คำขอของ AP-TEST-0002, 0003, 0004 ยังเป็น requested ทั้งหมด
		if s := tcRequestStatus(t, tx, id); s != "requested" {
			t.Errorf("(1) คำขอ %d ของชิ้นอื่นต้องยังเป็น requested แต่ได้ %q", id, s)
		}
		// Then (2) rejectedRequestIds ไม่มีคำขอของชิ้นอื่น
		if tcContains(got.RejectedRequestIDs, id) {
			t.Errorf("(2) rejectedRequestIds ต้องไม่มีคำขอ %d ของชิ้นอื่น", id)
		}
	}
}

// TC-APV-03-7 (AC-APV-03 ผิด)
func TestTC_APV_03_7_NoSupersedeWhenApproveFails(t *testing.T) {
	// Given AP-TEST-0004 ("ถูกยืม") มีคำขอแรก A และเพิ่มคำขอ B (ผู้ขออีกคน)
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0004", "student-001")
	b := tcAddRequest(t, tx, a, "student-005", 5*time.Minute)
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติ A
	w := tcApprove(srv, tokenStaffDept1, a)

	// Then (1) ได้ 409 ITEM_NOT_AVAILABLE
	tcExpectError(t, "(1)", w, http.StatusConflict, "ITEM_NOT_AVAILABLE", tcMsgItemNotAvailable)
	// Then (2) A และ B ยังเป็น requested
	for name, id := range map[string]int{"A": a, "B": b} {
		if s := tcRequestStatus(t, tx, id); s != "requested" {
			t.Errorf("(2) %s ต้องยังเป็น requested แต่ได้ %q", name, s)
		}
	}
	// Then (3) ไม่มีแถว approval_decision และ approval_email ของ B
	if n := tcDecisionCount(t, tx, b); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_decision ของ B แต่มี %d", n)
	}
	if n := tcEmailCount(t, tx, b); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_email ของ B แต่มี %d", n)
	}
}

// ---------- AC-APV-07 ----------

// TC-APV-07-3 (AC-APV-07 ขอบ)
func TestTC_APV_07_3_PastCutoffFirstStillBlocks(t *testing.T) {
	// Given AP-TEST-0001 มี A และ B เป็น requested ตั้งนาฬิกาเป็น 16:30 น. ของวันที่เริ่มยืมของ A
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	srv := tcServer(tx, &fakeClock{t: tcAt(t, tx, a, 16, 30, 0)})

	// When staff-001 อนุมัติ B
	w := tcApprove(srv, tokenStaffDept1, b)

	// Then (1) ได้ 409 NOT_FIRST_IN_QUEUE
	tcExpectError(t, "(1)", w, http.StatusConflict, "NOT_FIRST_IN_QUEUE", tcMsgNotFirst)
	// Then (2) A และ B ยังเป็น requested
	for name, id := range map[string]int{"A": a, "B": b} {
		if s := tcRequestStatus(t, tx, id); s != "requested" {
			t.Errorf("(2) %s ต้องยังเป็น requested แต่ได้ %q", name, s)
		}
	}
}

// TC-APV-07-4 (AC-APV-07 ขอบ)
func TestTC_APV_07_4_CancelledFirstFreesQueue(t *testing.T) {
	// Given AP-TEST-0001: A ถูกผู้ขอยกเลิก เหลือ B เป็น requested
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	tcSetStatus(t, tx, a, "cancelled")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, b))

	// When staff-001 อนุมัติ B
	w := tcApprove(srv, tokenStaffDept1, b)

	// Then (1) ได้ 200 B เป็น approved
	if w.Code != http.StatusOK {
		t.Errorf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, b); s != "approved" {
		t.Errorf("(1) B ต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (2) A ยังเป็น cancelled
	if s := tcRequestStatus(t, tx, a); s != "cancelled" {
		t.Errorf("(2) A ต้องยังเป็น cancelled แต่ได้ %q", s)
	}
}

// TC-APV-07-5 (AC-APV-07 ขอบ)
func TestTC_APV_07_5_ThirdInQueueBlocked(t *testing.T) {
	// Given AP-TEST-0001 มีคำขอยืมปกติ A, B, D (D ส่งหลังสุด)
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	d := tcAddRequest(t, tx, a, "student-005", 5*time.Minute)
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติ D
	w := tcApprove(srv, tokenStaffDept1, d)

	// Then (1) ได้ 409 NOT_FIRST_IN_QUEUE ข้อความตามที่กำหนด
	tcExpectError(t, "(1)", w, http.StatusConflict, "NOT_FIRST_IN_QUEUE", tcMsgNotFirst)
	// Then (2) A, B, D ยังเป็น requested
	for name, id := range map[string]int{"A": a, "B": b, "D": d} {
		if s := tcRequestStatus(t, tx, id); s != "requested" {
			t.Errorf("(2) %s ต้องยังเป็น requested แต่ได้ %q", name, s)
		}
	}
}

// TC-APV-07-6 (AC-APV-07 ผิด)
func TestTC_APV_07_6_SameTimeLowerIdIsFirst(t *testing.T) {
	// Given A และ B ของ AP-TEST-0001 submitted_at เท่ากันพอดี
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	if _, err := tx.Exec(context.Background(),
		`UPDATE borrow_request SET submitted_at = date_trunc('second', now()) - interval '30 minutes' WHERE id IN ($1, $2)`, a, b); err != nil {
		t.Fatal(err)
	}
	low, high := a, b
	if high < low {
		low, high = high, low
	}
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติคำขอที่รหัสมากกว่า แล้วอนุมัติคำขอที่รหัสน้อยกว่า
	wHigh := tcApprove(srv, tokenStaffDept1, high)
	wLow := tcApprove(srv, tokenStaffDept1, low)

	// Then (1) รหัสมากกว่าได้ 409 NOT_FIRST_IN_QUEUE
	tcExpectError(t, "(1)", wHigh, http.StatusConflict, "NOT_FIRST_IN_QUEUE", tcMsgNotFirst)
	// Then (2) รหัสน้อยกว่าได้ 200 approved
	if wLow.Code != http.StatusOK {
		t.Errorf("(2) รหัสน้อยกว่าต้องได้ 200 แต่ได้ %d body=%s", wLow.Code, wLow.Body.String())
	}
	if s := tcRequestStatus(t, tx, low); s != "approved" {
		t.Errorf("(2) รหัสน้อยกว่าต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (3) คำขอที่รหัสมากกว่าเป็น rejected ตาม
	if s := tcRequestStatus(t, tx, high); s != "rejected" {
		t.Errorf("(3) รหัสมากกว่าต้องเป็น rejected แต่ได้ %q", s)
	}
}

// ---------- AC-APV-08 ----------

// TC-APV-08-1 (AC-APV-08 ปกติ)
func TestTC_APV_08_1_UnavailableItem(t *testing.T) {
	// Given AP-TEST-0002 เปลี่ยนเป็น unavailable คำขอแรก requested ของ student-001
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	tcSetItemStatus(t, tx, "AP-TEST-0002", "unavailable")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, id))

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 409 ITEM_NOT_AVAILABLE ข้อความตามที่กำหนด
	tcExpectError(t, "(1)", w, http.StatusConflict, "ITEM_NOT_AVAILABLE", tcMsgItemNotAvailable)
	// Then (2) คำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) คำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	// Then (3) equipment_item.status ยังเป็น unavailable
	if s := tcItemStatus(t, tx, "AP-TEST-0002"); s != "unavailable" {
		t.Errorf("(3) equipment_item.status ต้องยังเป็น unavailable แต่ได้ %q", s)
	}
}

// TC-APV-08-3 (AC-APV-08 ขอบ)
func TestTC_APV_08_3_ItemAvailableAgain(t *testing.T) {
	// Given AP-TEST-0004 เปลี่ยนจาก borrowed เป็น available คำขอแรก requested ของ student-001
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0004", "student-001")
	tcSetItemStatus(t, tx, "AP-TEST-0004", "available")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, id))

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 200 คำขอเป็น approved
	if w.Code != http.StatusOK {
		t.Errorf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, id); s != "approved" {
		t.Errorf("(1) คำขอต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (2) equipment_item.status ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0004"); s != "available" {
		t.Errorf("(2) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
}

// TC-APV-08-4 (AC-APV-08 ผิด)
func TestTC_APV_08_4_NotAvailableAndNotFirst(t *testing.T) {
	// Given AP-TEST-0004 ("ถูกยืม") มีคำขอแรก A และเพิ่มคำขอ B
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0004", "student-001")
	b := tcAddRequest(t, tx, a, "student-005", 5*time.Minute)
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติ B
	w := tcApprove(srv, tokenStaffDept1, b)

	// Then (1) ได้ 409 NOT_FIRST_IN_QUEUE เพราะตรวจคำขอแรกก่อนสถานะอุปกรณ์ (Q-APV-35)
	tcExpectError(t, "(1)", w, http.StatusConflict, "NOT_FIRST_IN_QUEUE", tcMsgNotFirst)
	// Then (2) A และ B ยังเป็น requested และ equipment_item.status ยังเป็น borrowed
	for name, id := range map[string]int{"A": a, "B": b} {
		if s := tcRequestStatus(t, tx, id); s != "requested" {
			t.Errorf("(2) %s ต้องยังเป็น requested แต่ได้ %q", name, s)
		}
	}
	if s := tcItemStatus(t, tx, "AP-TEST-0004"); s != "borrowed" {
		t.Errorf("(2) equipment_item.status ต้องยังเป็น borrowed แต่ได้ %q", s)
	}
}

// ---------- AC-APV-09 ----------

// TC-APV-09-1 (AC-APV-09 ปกติ)
func TestTC_APV_09_1_ApproveAlreadyApproved(t *testing.T) {
	// Given คำขอของ AP-TEST-0002 ที่ staff-001 อนุมัติไปแล้ว และ loadedAt หลัง decided_at
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, nil)
	if w := tcApprove(srv, tokenStaffDept1, id); w.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: อนุมัติต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	_, _, _, decidedAt, _ := tcDecision(t, tx, id)

	// When staff-001 อนุมัติคำขอนั้นอีกครั้ง
	w := tcApproveWithLoaded(srv, id, decidedAt.Add(time.Minute))

	// Then (1) ได้ 409 REQUEST_ALREADY_DECIDED ข้อความ "คำขอนี้ได้รับการพิจารณาแล้ว"
	tcExpectError(t, "(1)", w, http.StatusConflict, "REQUEST_ALREADY_DECIDED", tcMsgAlreadyDecided)
	// Then (2) คำขอยังเป็น approved และมี approval_decision แถวเดียว
	if s := tcRequestStatus(t, tx, id); s != "approved" {
		t.Errorf("(2) คำขอต้องยังเป็น approved แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 1 {
		t.Errorf("(2) ต้องมี approval_decision แถวเดียว แต่มี %d", n)
	}
}

// TC-APV-09-2 (AC-APV-09 ปกติ)
func TestTC_APV_09_2_DecideAlreadyRejected(t *testing.T) {
	// Given คำขอของ AP-TEST-0003 ที่ staff-001 ปฏิเสธไปแล้วด้วยเหตุผล "อุปกรณ์ต้องใช้สอบ" (loadedAt หลัง decided_at)
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, id))
	if w := tcReject(srv, tokenStaffDept1, id, "อุปกรณ์ต้องใช้สอบ", time.Now()); w.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: ปฏิเสธต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	_, _, _, decidedAt, _ := tcDecision(t, tx, id)
	loaded := decidedAt.Add(time.Minute)

	// When staff-001 อนุมัติ แล้วปฏิเสธอีกครั้งด้วยเหตุผลใหม่
	wa := tcApproveWithLoaded(srv, id, loaded)
	wr := tcReject(srv, tokenStaffDept1, id, "เหตุผลใหม่", loaded)

	// Then (1) ได้ 409 REQUEST_ALREADY_DECIDED ทั้ง 2 ครั้ง
	tcExpectError(t, "(1) approve", wa, http.StatusConflict, "REQUEST_ALREADY_DECIDED", tcMsgAlreadyDecided)
	tcExpectError(t, "(1) reject", wr, http.StatusConflict, "REQUEST_ALREADY_DECIDED", tcMsgAlreadyDecided)
	// Then (2) คำขอยังเป็น rejected เหตุผลยังเป็น "อุปกรณ์ต้องใช้สอบ" มี approval_decision แถวเดียว
	if s := tcRequestStatus(t, tx, id); s != "rejected" {
		t.Errorf("(2) คำขอต้องยังเป็น rejected แต่ได้ %q", s)
	}
	if _, reason, _, _, _ := tcDecision(t, tx, id); reason != "อุปกรณ์ต้องใช้สอบ" {
		t.Errorf("(2) เหตุผลต้องยังเป็น \"อุปกรณ์ต้องใช้สอบ\" แต่ได้ %q", reason)
	}
	if n := tcDecisionCount(t, tx, id); n != 1 {
		t.Errorf("(2) ต้องมี approval_decision แถวเดียว แต่มี %d", n)
	}
}

// TC-APV-09-4 (AC-APV-09 ขอบ)
func TestTC_APV_09_4_AutoRejectedNotDecidable(t *testing.T) {
	// Given คำขอยืมปกติที่ถูกปฏิเสธอัตโนมัติตอนเที่ยงคืนแล้ว (ผล auto_rejected ผู้ดำเนินการ System)
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	setSubmittedAt(t, tx, id, autoRejectCutoff(time.Now()).Add(-3*time.Hour)) // ส่งเมื่อวาน
	job := &AutoRejecter{db: tx, now: time.Now}
	if _, err := job.RejectStale(context.Background()); err != nil {
		t.Fatal(err)
	}
	res, _, _, decidedAt, ok := tcDecision(t, tx, id)
	if !ok || res != "auto_rejected" {
		t.Fatalf("เตรียมข้อมูล: คำขอต้องถูกปฏิเสธอัตโนมัติ แต่ได้ %q (มีแถว %v)", res, ok)
	}
	srv := tcServer(tx, nil)

	// When staff-001 อนุมัติคำขอนั้น โดย loadedAt หลังเวลาที่ปฏิเสธอัตโนมัติ
	w := tcApproveWithLoaded(srv, id, decidedAt.Add(time.Minute))

	// Then (1) ได้ 409 REQUEST_ALREADY_DECIDED
	tcExpectError(t, "(1)", w, http.StatusConflict, "REQUEST_ALREADY_DECIDED", tcMsgAlreadyDecided)
	// Then (2) คำขอยังเป็น rejected
	if s := tcRequestStatus(t, tx, id); s != "rejected" {
		t.Errorf("(2) คำขอต้องยังเป็น rejected แต่ได้ %q", s)
	}
}

// TC-APV-09-5 (AC-APV-09 ขอบ)
func TestTC_APV_09_5_SupersededNotDecidable(t *testing.T) {
	// Given B ของ AP-TEST-0001 ถูกปฏิเสธตามเพราะ A ได้รับอนุมัติ (loadedAt หลังเวลาที่ A ได้รับอนุมัติ)
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))
	if w := tcApprove(srv, tokenStaffDept1, a); w.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: อนุมัติ A ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	_, _, _, decidedAt, _ := tcDecision(t, tx, a)

	// When staff-001 ปฏิเสธ B ด้วยเหตุผล
	w := tcReject(srv, tokenStaffDept1, b, "อุปกรณ์ต้องใช้สอบ", decidedAt.Add(time.Minute))

	// Then (1) ได้ 409 REQUEST_ALREADY_DECIDED
	tcExpectError(t, "(1)", w, http.StatusConflict, "REQUEST_ALREADY_DECIDED", tcMsgAlreadyDecided)
	// Then (2) B ยังเป็น rejected เหตุผลยังเป็น "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว"
	if s := tcRequestStatus(t, tx, b); s != "rejected" {
		t.Errorf("(2) B ต้องยังเป็น rejected แต่ได้ %q", s)
	}
	if _, reason, _, _, _ := tcDecision(t, tx, b); reason != tcMsgSuperseded {
		t.Errorf("(2) เหตุผลของ B ต้องยังเป็น %q แต่ได้ %q", tcMsgSuperseded, reason)
	}
}

// TC-APV-09-6 (AC-APV-09 ขอบ)
func TestTC_APV_09_6_HandedOverNotDecidable(t *testing.T) {
	// Given คำขอของ AP-TEST-0002 สถานะ handed_over (005 ส่งมอบแล้ว)
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	tcSetStatus(t, tx, id, "handed_over")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, id))

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 409 REQUEST_ALREADY_DECIDED ข้อความ "คำขอนี้ได้รับการพิจารณาแล้ว"
	tcExpectError(t, "(1)", w, http.StatusConflict, "REQUEST_ALREADY_DECIDED", tcMsgAlreadyDecided)
	// Then (2) คำขอยังเป็น handed_over
	if s := tcRequestStatus(t, tx, id); s != "handed_over" {
		t.Errorf("(2) คำขอต้องยังเป็น handed_over แต่ได้ %q", s)
	}
}

// ---------- AC-APV-16 ----------

// TC-APV-16-1 (AC-APV-16 ปกติ)
func TestTC_APV_16_1_ApproveAfterCutoff(t *testing.T) {
	// Given คำขอแรก requested ของ AP-TEST-0002 ("ว่าง") ตั้งนาฬิกาเป็น 16:30 น. เวลาไทยของวันที่เริ่มยืม
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, &fakeClock{t: tcAt(t, tx, id, 16, 30, 0)})

	// When staff-001 เรียก POST .../approve
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 409 APPROVAL_CUTOFF_PASSED ข้อความ "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้"
	tcExpectError(t, "(1)", w, http.StatusConflict, "APPROVAL_CUTOFF_PASSED", tcMsgCutoff)
	// Then (2) คำขอยังเป็น requested และไม่มีแถว approval_decision
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) คำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision แต่มี %d", n)
	}
}

// TC-APV-16-3 (AC-APV-16 ขอบ)
func TestTC_APV_16_3_CutoffIsStartDayNotToday(t *testing.T) {
	// Given คำขอแรก requested ของ AP-TEST-0002 ที่วันที่เริ่มยืมเป็นพรุ่งนี้ ตั้งนาฬิกาเป็น 17:00 น. ของวันนี้
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	today := tcAt(t, tx, id, 17, 0, 0) // วันที่เริ่มยืมเดิม = วันนี้
	if _, err := tx.Exec(context.Background(),
		`UPDATE borrow_request SET start_at = start_at + interval '1 day', end_at = end_at + interval '1 day' WHERE id = $1`, id); err != nil {
		t.Fatal(err)
	}
	if !localTime(tcStartAt(t, tx, id)).After(today) {
		t.Fatal("เตรียมข้อมูล: วันที่เริ่มยืมต้องเป็นพรุ่งนี้")
	}
	srv := tcServer(tx, &fakeClock{t: today})

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 200 คำขอเป็น approved
	if w.Code != http.StatusOK {
		t.Errorf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, id); s != "approved" {
		t.Errorf("(1) คำขอต้องเป็น approved แต่ได้ %q", s)
	}
}

// TC-APV-16-4 (AC-APV-16 ขอบ)
func TestTC_APV_16_4_CutoffInBangkokTime(t *testing.T) {
	// Given คำขอแรก requested ของ AP-TEST-0002 เริ่มยืม 09:00 น. เวลาไทย
	// When  staff-001 อนุมัติตอน 09:00 UTC (= 16:00 น.) และอีกกรณี 08:59:59 UTC (= 15:59:59 น.)
	run := func(t *testing.T, h, m, s int) (pgx.Tx, int, *httptest.ResponseRecorder) {
		tx := devTx(t)
		id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
		start := localTime(tcStartAt(t, tx, id))
		if start.Hour() != 9 || start.Minute() != 0 {
			t.Fatalf("เตรียมข้อมูล: ต้องเริ่มยืม 09:00 น. เวลาไทย แต่ได้ %v", start)
		}
		at := time.Date(start.Year(), start.Month(), start.Day(), h, m, s, 0, time.UTC)
		return tx, id, tcApprove(tcServer(tx, &fakeClock{t: at}), tokenStaffDept1, id)
	}

	// Then (1) 09:00 UTC ได้ 409 APPROVAL_CUTOFF_PASSED
	t.Run("0900UTC", func(t *testing.T) {
		_, _, w := run(t, 9, 0, 0)
		tcExpectError(t, "(1)", w, http.StatusConflict, "APPROVAL_CUTOFF_PASSED", tcMsgCutoff)
	})
	// Then (2) 08:59:59 UTC ได้ 200 approved
	t.Run("085959UTC", func(t *testing.T) {
		tx, id, w := run(t, 8, 59, 59)
		if w.Code != http.StatusOK {
			t.Errorf("(2) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
		}
		if s := tcRequestStatus(t, tx, id); s != "approved" {
			t.Errorf("(2) คำขอต้องเป็น approved แต่ได้ %q", s)
		}
	})
}

// TC-APV-16-5 (AC-APV-16 ผิด)
func TestTC_APV_16_5_NotAvailableAndPastCutoff(t *testing.T) {
	// Given AP-TEST-0004 ("ถูกยืม") คำขอแรก requested ตั้งนาฬิกาเป็น 16:30 น. ของวันที่เริ่มยืม
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0004", "student-001")
	srv := tcServer(tx, &fakeClock{t: tcAt(t, tx, id, 16, 30, 0)})

	// When staff-001 อนุมัติคำขอนั้น
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 409 ITEM_NOT_AVAILABLE เพราะตรวจสถานะอุปกรณ์ก่อนเวลา 16:00 น. (Q-APV-36)
	tcExpectError(t, "(1)", w, http.StatusConflict, "ITEM_NOT_AVAILABLE", tcMsgItemNotAvailable)
	// Then (2) คำขอยังเป็น requested และไม่มีแถว approval_decision
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) คำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision แต่มี %d", n)
	}
}

// ---------- AC-APV-17 ----------

// TC-APV-17-1 (AC-APV-17 ปกติ)
func TestTC_APV_17_1_OvernightSupersededWithEmail(t *testing.T) {
	// Given AP-TEST-0001 มีคำขอยืมปกติ A (ส่งก่อน) และคำขอค้างคืน C (student-003) เป็น requested
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติ A
	if w := tcApprove(srv, tokenStaffDept1, a); w.Code != http.StatusOK {
		t.Fatalf("อนุมัติ A ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}

	// Then (1) C เป็น rejected
	if s := tcRequestStatus(t, tx, c); s != "rejected" {
		t.Errorf("(1) C ต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (2) approval_decision ของ C เหตุผลตามที่กำหนด ผู้ดำเนินการ staff-001
	if _, reason, actor, _, ok := tcDecision(t, tx, c); !ok || reason != tcMsgSuperseded || actor != "staff-001" {
		t.Errorf("(2) C ต้องมีเหตุผล %q โดย staff-001 แต่ได้ %q โดย %q (มีแถว %v)", tcMsgSuperseded, reason, actor, ok)
	}
	// Then (3) มีแถว approval_email ของ C แบบ rejected
	if k := tcEmailKind(t, tx, c); k != "rejected" {
		t.Errorf("(3) Email ของ C ต้องเป็นแบบ rejected แต่ได้ %q", k)
	}
}

// tcOvernightAt ตั้งเวลาที่ส่งของคำขอค้างคืน C เป็น 2 ชั่วโมงก่อนเที่ยงคืนของคืนนี้ (วันที่ D = วันนี้) คืน C และเที่ยงคืนหลังวันที่ D
func tcOvernightAt(t *testing.T, tx pgx.Tx) (int, time.Time) {
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	midnight := autoRejectCutoff(time.Now()).AddDate(0, 0, 1)
	setSubmittedAt(t, tx, c, midnight.Add(-2*time.Hour))
	return c, midnight
}

// TC-APV-17-2 (AC-APV-17 ปกติ)
func TestTC_APV_17_2_MidnightSkipsOvernight(t *testing.T) {
	// Given C ส่งเมื่อวันที่ D และยังเป็น requested ตั้งนาฬิกาเป็น 00:00 น. หลังวันที่ D
	tx := devTx(t)
	c, midnight := tcOvernightAt(t, tx)

	// When งานปฏิเสธอัตโนมัติทำงาน (RejectStale)
	job := &AutoRejecter{db: tx, now: func() time.Time { return midnight }}
	if _, err := job.RejectStale(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Then (1) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(1) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
	// Then (2) ไม่มีแถว approval_decision และ approval_email ของ C
	if n := tcDecisionCount(t, tx, c); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision ของ C แต่มี %d", n)
	}
	if n := tcEmailCount(t, tx, c); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_email ของ C แต่มี %d", n)
	}
}

// TC-APV-17-3 (AC-APV-17 ขอบ)
func TestTC_APV_17_3_OvernightStaysAfterDays(t *testing.T) {
	// Given C ส่งเมื่อวันที่ D ยังเป็น requested ตั้งนาฬิกาเป็น 3 วันหลังวันที่ D
	tx := devTx(t)
	c, midnight := tcOvernightAt(t, tx)

	// When งานปฏิเสธอัตโนมัติทำงาน
	job := &AutoRejecter{db: tx, now: func() time.Time { return midnight.AddDate(0, 0, 3) }}
	if _, err := job.RejectStale(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Then (1) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(1) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-17-4 (AC-APV-17 ขอบ)
func TestTC_APV_17_4_CatchUpSkipsOvernight(t *testing.T) {
	// Given C และคำขอยืมปกติ E (AP-TEST-0003) ส่งวันที่ D ระบบไม่ได้ทำงานตอนเที่ยงคืน แล้วเริ่มทำงานหลัง 00:00 น.
	tx := devTx(t)
	c, midnight := tcOvernightAt(t, tx)
	e := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	setSubmittedAt(t, tx, e, midnight.Add(-3*time.Hour))

	// When งานปฏิเสธอัตโนมัติไล่ปฏิเสธตอนเริ่มทำงาน (CatchUp)
	job := &AutoRejecter{db: tx, now: func() time.Time { return midnight.Add(9*time.Hour + 30*time.Minute) }}
	job.CatchUp(context.Background())

	// Then (1) E เป็น rejected ผล auto_rejected
	if s := tcRequestStatus(t, tx, e); s != "rejected" {
		t.Errorf("(1) E ต้องเป็น rejected แต่ได้ %q", s)
	}
	if res, _, _, _, _ := tcDecision(t, tx, e); res != "auto_rejected" {
		t.Errorf("(1) ผลของ E ต้องเป็น auto_rejected แต่ได้ %q", res)
	}
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-17-5 (AC-APV-17 ผิด)
func TestTC_APV_17_5_FailedApproveKeepsOvernight(t *testing.T) {
	// Given ตั้งนาฬิกาเป็น 16:30 น. ของวันที่เริ่มยืมของ A (C ค้างคืนยังเป็น requested)
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	srv := tcServer(tx, &fakeClock{t: tcAt(t, tx, a, 16, 30, 0)})

	// When staff-001 อนุมัติ A
	w := tcApprove(srv, tokenStaffDept1, a)

	// Then (1) ได้ 409 APPROVAL_CUTOFF_PASSED
	tcExpectError(t, "(1)", w, http.StatusConflict, "APPROVAL_CUTOFF_PASSED", tcMsgCutoff)
	// Then (2) C ยังเป็น requested และไม่มีแถว approval_decision ของ C
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, c); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision ของ C แต่มี %d", n)
	}
}
