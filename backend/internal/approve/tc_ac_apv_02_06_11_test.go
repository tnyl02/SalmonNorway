// test จากตาราง specs/004-approve-request/test-cases.md แถว TC-APV-02-x, TC-APV-06-x, TC-APV-11-x สถานะ "ใช้ได้" ห้ามแก้ให้ผ่าน
// ถ้าไม่ผ่าน ต้องหาว่าโค้ดผิด หรือ AC ผิด (ถ้า AC ผิด แก้ spec ก่อน)
// ใช้ helper tcServer, tcApprove, tcRequestStatus, tcItemStatus, tcStartAt, tcError จาก tc_ac_apv_01_test.go
package approve

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
)

// tcReject เรียก POST /v1/approvals/requests/{id}/reject พร้อม reason และ loadedAt (plan ข้อ 4)
func tcReject(srv http.Handler, token string, id int, reason string, loadedAt time.Time) *httptest.ResponseRecorder {
	body, _ := json.Marshal(map[string]string{"reason": reason, "loadedAt": loadedAt.UTC().Format(time.RFC3339)})
	return postJSON(srv, token, fmt.Sprintf("%srequests/%d/reject", PathPrefix, id), string(body))
}

// tcDecisionCount นับแถว approval_decision ของคำขอ
func tcDecisionCount(t *testing.T, tx pgx.Tx, id int) int {
	t.Helper()
	var n int
	if err := tx.QueryRow(context.Background(), `SELECT count(*) FROM approval_decision WHERE request_id = $1`, id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// tcList เรียก GET /v1/approvals/requests แล้วอ่านคำตอบ
func tcList(t *testing.T, srv http.Handler, token string) (*httptest.ResponseRecorder, PendingListDTO) {
	t.Helper()
	w := getJSON(srv, token, PathPrefix+"requests")
	var got PendingListDTO
	if w.Code == http.StatusOK {
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
	}
	return w, got
}

// tcHasItems บอกว่า body มีช่อง items หรือไม่
func tcHasItems(w *httptest.ResponseRecorder) bool {
	var raw map[string]any
	if json.Unmarshal(w.Body.Bytes(), &raw) != nil {
		return false
	}
	_, ok := raw["items"]
	return ok
}

// tcBeforeCutoff คืนนาฬิกาจำลอง 1 ชั่วโมงก่อน 16:00 น. ของวันที่เริ่มยืมของคำขอ
func tcBeforeCutoff(t *testing.T, tx pgx.Tx, id int) *fakeClock {
	return &fakeClock{t: approvalCutoff(tcStartAt(t, tx, id)).Add(-time.Hour)}
}

// ---------- AC-APV-02 ----------

// TC-APV-02-1 (AC-APV-02 ปกติ)
func TestTC_APV_02_1_RejectWithReason(t *testing.T) {
	// Given staff-001 ยืนยันตัวตนแล้ว คำขอยืมปกติ requested ของ student-002 บน AP-TEST-0003 ("ว่าง")
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServer(tx, nil)

	// When staff-001 เรียก POST .../reject พร้อม reason "อุปกรณ์ต้องใช้สอบ" และ loadedAt
	w := tcReject(srv, tokenStaffDept1, id, "อุปกรณ์ต้องใช้สอบ", time.Now())

	// Then (1) ได้ 200 คำตอบมี requestId ของคำขอนั้นและ status rejected
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var got ResultDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.RequestID != id || got.Status != "rejected" {
		t.Errorf("(1) คำตอบต้องเป็น requestId %d status rejected แต่ได้ %+v", id, got)
	}
	// Then (2) ในฐานข้อมูล borrow_request.status เป็น rejected
	if s := tcRequestStatus(t, tx, id); s != "rejected" {
		t.Errorf("(2) borrow_request.status ต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (3) equipment_item.status ของ AP-TEST-0003 ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0003"); s != "available" {
		t.Errorf("(3) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
	// Then (4) มีแถว approval_decision result rejected reason ตรงตัว ผู้ดำเนินการ staff-001
	var result, reason, actor string
	if err := tx.QueryRow(context.Background(),
		`SELECT result, reason, actor_id FROM approval_decision WHERE request_id = $1`, id).Scan(&result, &reason, &actor); err != nil {
		t.Fatalf("(4) ต้องมีแถว approval_decision ของคำขอนั้น: %v", err)
	}
	if result != "rejected" || reason != "อุปกรณ์ต้องใช้สอบ" || actor != "staff-001" {
		t.Errorf("(4) ต้องเป็น rejected / อุปกรณ์ต้องใช้สอบ / staff-001 แต่ได้ %q / %q / %q", result, reason, actor)
	}
}

// TC-APV-02-4 (AC-APV-02 ขอบ)
func TestTC_APV_02_4_RejectBorrowedItemUnchanged(t *testing.T) {
	// Given คำขอแรก requested ของ AP-TEST-0004 ("ถูกยืม")
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0004", "student-001")
	srv := tcServer(tx, nil)

	// When staff-001 ปฏิเสธพร้อมเหตุผล
	w := tcReject(srv, tokenStaffDept1, id, "อุปกรณ์ยังไม่ถูกคืน", time.Now())

	// Then (1) ได้ 200 คำขอเป็น rejected
	if w.Code != http.StatusOK {
		t.Errorf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, id); s != "rejected" {
		t.Errorf("(1) คำขอต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (2) equipment_item.status ยังเป็น borrowed
	if s := tcItemStatus(t, tx, "AP-TEST-0004"); s != "borrowed" {
		t.Errorf("(2) equipment_item.status ต้องยังเป็น borrowed แต่ได้ %q", s)
	}
}

// TC-APV-02-5 (AC-APV-02 ขอบ)
func TestTC_APV_02_5_RejectAfterCutoff(t *testing.T) {
	// Given คำขอแรก A ของ AP-TEST-0001 ตั้งนาฬิกาจำลองเป็น 16:00:00 ของวันที่เริ่มยืม
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	clock := &fakeClock{t: approvalCutoff(tcStartAt(t, tx, a))}
	if lt := localTime(clock.t); lt.Hour() != 16 || lt.Minute() != 0 || lt.Second() != 0 {
		t.Fatalf("นาฬิกาจำลองต้องเป็น 16:00:00 เวลาไทย แต่ได้ %v", lt)
	}
	srv := tcServer(tx, clock)

	// When staff-001 ปฏิเสธ A พร้อมเหตุผล
	w := tcReject(srv, tokenStaffDept1, a, "เลยเวลาที่อนุมัติได้", clock.t)

	// Then (1) ได้ 200 A เป็น rejected
	if w.Code != http.StatusOK {
		t.Errorf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, a); s != "rejected" {
		t.Errorf("(1) A ต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (2) equipment_item.status ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0001"); s != "available" {
		t.Errorf("(2) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
}

// TC-APV-02-6 (AC-APV-02 ขอบ)
func TestTC_APV_02_6_RejectFirstNextBecomesFirst(t *testing.T) {
	// Given AP-TEST-0001 มี A (คำขอแรก), B (ยืมปกติ) และ C (ค้างคืน) เป็น requested
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 ปฏิเสธ A พร้อมเหตุผล แล้วเรียก GET /v1/approvals/requests
	if w := tcReject(srv, tokenStaffDept1, a, "อุปกรณ์ต้องใช้สอบ", time.Now()); w.Code != http.StatusOK {
		t.Fatalf("ปฏิเสธ A ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	lw, list := tcList(t, srv, tokenStaffDept1)
	if lw.Code != http.StatusOK {
		t.Fatalf("เปิดรายการต้องได้ 200 แต่ได้ %d", lw.Code)
	}

	// Then (1) A เป็น rejected
	if s := tcRequestStatus(t, tx, a); s != "rejected" {
		t.Errorf("(1) A ต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (2) B และ C ยังเป็น requested
	if s := tcRequestStatus(t, tx, b); s != "requested" {
		t.Errorf("(2) B ต้องยังเป็น requested แต่ได้ %q", s)
	}
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
	// Then (3) ในรายการ B มี isFirst: true
	found := false
	for _, it := range list.Items {
		if it.RequestID == b {
			found = true
			if !it.IsFirst {
				t.Error("(3) B ต้องมี isFirst: true หลัง A ถูกปฏิเสธ")
			}
		}
	}
	if !found {
		t.Error("(3) B ต้องอยู่ในรายการ")
	}
}

// TC-APV-02-7 (AC-APV-02 ขอบ)
func TestTC_APV_02_7_RejectNotFirst(t *testing.T) {
	// Given AP-TEST-0001 มี A (คำขอแรก) และ B เป็น requested
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 ปฏิเสธ B (ไม่ใช่คำขอแรก) พร้อมเหตุผล
	w := tcReject(srv, tokenStaffDept1, b, "อุปกรณ์ต้องใช้สอบ", time.Now())

	// Then (1) ได้ 200 B เป็น rejected (Q-APV-27)
	if w.Code != http.StatusOK {
		t.Errorf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, b); s != "rejected" {
		t.Errorf("(1) B ต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (2) A ยังเป็น requested และยังเป็นคำขอแรก
	if s := tcRequestStatus(t, tx, a); s != "requested" {
		t.Errorf("(2) A ต้องยังเป็น requested แต่ได้ %q", s)
	}
	_, list := tcList(t, srv, tokenStaffDept1)
	foundA := false
	for _, it := range list.Items {
		if it.RequestID == a {
			foundA = true
			if !it.IsFirst {
				t.Error("(2) A ต้องยังเป็นคำขอแรก (isFirst: true)")
			}
		}
	}
	if !foundA {
		t.Error("(2) A ต้องยังอยู่ในรายการ")
	}
	// Then (3) equipment_item.status ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0001"); s != "available" {
		t.Errorf("(3) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
}

// TC-APV-02-8 (AC-APV-02 ผิด)
func TestTC_APV_02_8_NoTokenNotRejected(t *testing.T) {
	// Given ไม่มี token คำขอเดียวกับ TC-APV-02-1
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServer(tx, nil)

	// When เรียก POST .../reject พร้อมเหตุผล
	w := tcReject(srv, "", id, "อุปกรณ์ต้องใช้สอบ", time.Now())

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

// TC-APV-02-9 (AC-APV-02 ผิด)
func TestTC_APV_02_9_CancelledNotRejected(t *testing.T) {
	// Given คำขอ cancelled ของ AP-TEST-0005
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0005", "student-002")
	srv := tcServer(tx, nil)

	// When staff-001 ปฏิเสธพร้อมเหตุผล
	w := tcReject(srv, tokenStaffDept1, id, "อุปกรณ์ต้องใช้สอบ", time.Now())

	// Then (1) ได้ 409 REQUEST_CANCELLED ข้อความ "คำขอนี้ถูกผู้ขอยกเลิกแล้ว"
	if w.Code != http.StatusConflict {
		t.Errorf("(1) ต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if e := tcError(t, w); e.Code != "REQUEST_CANCELLED" || e.Message != "คำขอนี้ถูกผู้ขอยกเลิกแล้ว" {
		t.Errorf("(1) ต้องได้ REQUEST_CANCELLED พร้อมข้อความตาม FR-APV-01 (ซ) แต่ได้ %+v", e)
	}
	// Then (2) คำขอยังเป็น cancelled และไม่มีแถว approval_decision ของคำขอนั้น
	if s := tcRequestStatus(t, tx, id); s != "cancelled" {
		t.Errorf("(2) คำขอต้องยังเป็น cancelled แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(2) ต้องไม่มีแถว approval_decision แต่มี %d แถว", n)
	}
}

// TC-APV-02-10 (AC-APV-02 ผิด)
func TestTC_APV_02_10_AlreadyApprovedNotRejected(t *testing.T) {
	// Given คำขอของ AP-TEST-0002 ที่ staff-001 อนุมัติไปแล้ว และหน้าจอโหลดรายการหลังการอนุมัตินั้น (loadedAt หลัง decided_at)
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, nil)
	if w := tcApprove(srv, tokenStaffDept1, id); w.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: อนุมัติต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var decidedAt time.Time
	if err := tx.QueryRow(context.Background(),
		`SELECT decided_at FROM approval_decision WHERE request_id = $1`, id).Scan(&decidedAt); err != nil {
		t.Fatal(err)
	}
	loadedAt := decidedAt.Add(time.Minute)

	// When staff-001 ปฏิเสธคำขอนั้นพร้อมเหตุผล
	w := tcReject(srv, tokenStaffDept1, id, "อุปกรณ์ต้องใช้สอบ", loadedAt)

	// Then (1) ได้ 409 REQUEST_ALREADY_DECIDED ข้อความ "คำขอนี้ได้รับการพิจารณาแล้ว"
	if w.Code != http.StatusConflict {
		t.Errorf("(1) ต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if e := tcError(t, w); e.Code != "REQUEST_ALREADY_DECIDED" || e.Message != "คำขอนี้ได้รับการพิจารณาแล้ว" {
		t.Errorf("(1) ต้องได้ REQUEST_ALREADY_DECIDED พร้อมข้อความตาม FR-APV-01 (ซ) แต่ได้ %+v", e)
	}
	// Then (2) คำขอยังเป็น approved และ approval_decision ยังเป็นผลอนุมัติเดิม
	if s := tcRequestStatus(t, tx, id); s != "approved" {
		t.Errorf("(2) คำขอต้องยังเป็น approved แต่ได้ %q", s)
	}
	var result string
	var n int
	if err := tx.QueryRow(context.Background(),
		`SELECT min(result), count(*) FROM approval_decision WHERE request_id = $1`, id).Scan(&result, &n); err != nil {
		t.Fatal(err)
	}
	if n != 1 || result != "approved" {
		t.Errorf("(2) approval_decision ต้องมีแถวเดียวผล approved แต่มี %d แถว ผล %q", n, result)
	}
}

// TC-APV-02-11 (AC-APV-02 ขอบ) แก้ตาม spec v3 CR-009 2569-11-13
func TestTC_APV_02_11_OvernightRejected(t *testing.T) {
	// Given คำขอค้างคืน C ของ AP-TEST-0001
	tx := devTx(t)
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	srv := tcServer(tx, nil)

	// When staff-001 ปฏิเสธ C พร้อมเหตุผล
	w := tcReject(srv, tokenStaffDept1, c, "อุปกรณ์ต้องใช้สอบ", time.Now())

	// Then (1) ได้ 200 status = rejected
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var res ResultDTO
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res.Status != "rejected" {
		t.Errorf("(1) status ต้องเป็น rejected แต่ได้ %q", res.Status)
	}
	// Then (2) C เป็น rejected
	if s := tcRequestStatus(t, tx, c); s != "rejected" {
		t.Errorf("(2) C ต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (3) A และ B ของ AP-TEST-0001 ยังเป็น requested
	for _, who := range []string{"student-001", "student-002"} {
		if s := tcRequestStatus(t, tx, requestIDOf(t, tx, "AP-TEST-0001", who)); s != "requested" {
			t.Errorf("(3) คำขอของ %s ต้องยังเป็น requested แต่ได้ %q", who, s)
		}
	}
}

// ---------- AC-APV-06 ----------

// tcExpectReasonRequired ตรวจ 422 VALIDATION_FAILED พร้อมข้อความ "กรุณาระบุเหตุผลในการปฏิเสธ"
func tcExpectReasonRequired(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("(1) ต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
		return
	}
	if e := tcError(t, w); e.Code != "VALIDATION_FAILED" || e.Message != "กรุณาระบุเหตุผลในการปฏิเสธ" {
		t.Errorf("(1) ต้องได้ VALIDATION_FAILED พร้อมข้อความ \"กรุณาระบุเหตุผลในการปฏิเสธ\" แต่ได้ %+v", e)
	}
}

// TC-APV-06-1 (AC-APV-06 ปกติ)
func TestTC_APV_06_1_EmptyReasonRejected(t *testing.T) {
	// Given staff-001 ยืนยันตัวตนแล้ว คำขอ requested ของ AP-TEST-0003
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServer(tx, nil)

	// When staff-001 เรียก POST .../reject ด้วย reason ""
	w := tcReject(srv, tokenStaffDept1, id, "", time.Now())

	// Then (1) ได้ 422 VALIDATION_FAILED ข้อความ "กรุณาระบุเหตุผลในการปฏิเสธ"
	tcExpectReasonRequired(t, w)
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	// Then (3) ไม่มีแถว approval_decision ของคำขอนั้น
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(3) ต้องไม่มีแถว approval_decision แต่มี %d แถว", n)
	}
}

// TC-APV-06-3 (AC-APV-06 ขอบ)
func TestTC_APV_06_3_WhitespaceReasonRejected(t *testing.T) {
	// Given คำขอ requested ของ AP-TEST-0003 เหตุผลที่มีแต่ช่องว่าง "   "
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServer(tx, nil)

	// When staff-001 เรียก POST .../reject ด้วย reason "   "
	w := tcReject(srv, tokenStaffDept1, id, "   ", time.Now())

	// Then (1) ได้ 422 VALIDATION_FAILED ข้อความ "กรุณาระบุเหตุผลในการปฏิเสธ"
	tcExpectReasonRequired(t, w)
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-06-4 (AC-APV-06 ขอบ)
func TestTC_APV_06_4_OneCharReasonAccepted(t *testing.T) {
	// Given คำขอ requested ของ AP-TEST-0003 เหตุผลสั้นที่สุดที่มีตัวอักษร 1 ตัว "x"
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServer(tx, nil)

	// When staff-001 ปฏิเสธด้วย reason "x"
	w := tcReject(srv, tokenStaffDept1, id, "x", time.Now())

	// Then (1) ได้ 200 คำขอเป็น rejected
	if w.Code != http.StatusOK {
		t.Errorf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if s := tcRequestStatus(t, tx, id); s != "rejected" {
		t.Errorf("(1) คำขอต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (2) approval_decision.reason เป็น "x"
	var reason string
	if err := tx.QueryRow(context.Background(),
		`SELECT reason FROM approval_decision WHERE request_id = $1`, id).Scan(&reason); err != nil {
		t.Fatalf("(2) ต้องมีแถว approval_decision: %v", err)
	}
	if reason != "x" {
		t.Errorf("(2) reason ต้องเป็น \"x\" แต่ได้ %q", reason)
	}
}

// TC-APV-06-6 (AC-APV-06 ขอบ)
func TestTC_APV_06_6_VeryLongReason(t *testing.T) {
	// Given คำขอ requested ของ AP-TEST-0003 (แต่ละกรณีใช้ฐานข้อมูลใหม่ คำขอจึงเริ่มเป็น requested ทุกครั้ง)
	// When  staff-001 ปฏิเสธด้วยเหตุผลยาว 500 ตัวอักษร, 501 ตัวอักษร และภาษาไทย 501 code point
	reject := func(t *testing.T, reason string) (pgx.Tx, int, *httptest.ResponseRecorder) {
		tx := devTx(t)
		id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
		return tx, id, tcReject(tcServer(tx, nil), tokenStaffDept1, id, reason, time.Now())
	}

	// Then (1) เหตุผลยาว 500 ตัวอักษรพอดี (นับตาม code point): ได้ 200 คำขอเป็น rejected
	t.Run("500", func(t *testing.T) {
		tx, id, w := reject(t, strings.Repeat("a", 500))
		if w.Code != http.StatusOK {
			t.Errorf("(1) 500 ตัวอักษรต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
		}
		if s := tcRequestStatus(t, tx, id); s != "rejected" {
			t.Errorf("(1) คำขอต้องเป็น rejected แต่ได้ %q", s)
		}
	})
	// Then (2) เหตุผลยาว 501 ตัวอักษร: ได้ 422 VALIDATION_FAILED ข้อความ "เหตุผลต้องไม่เกิน 500 ตัวอักษร" และสถานะคำขอยังเป็น requested
	t.Run("501", func(t *testing.T) {
		tx, id, w := reject(t, strings.Repeat("a", 501))
		tcExpectReasonTooLong(t, w, "(2)")
		if s := tcRequestStatus(t, tx, id); s != "requested" {
			t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
		}
	})
	// Then (3) ภาษาไทย "ที่" ซ้ำ 167 ครั้ง (เห็นบนจอ 167 ตัว แต่เป็น 501 code point): ได้ 422 VALIDATION_FAILED เช่นกัน
	t.Run("thai-501-code-points", func(t *testing.T) {
		reason := strings.Repeat("ที่", 167)
		if n := utf8.RuneCountInString(reason); n != 501 {
			t.Fatalf("เตรียมข้อมูล: ต้องเป็น 501 code point แต่ได้ %d", n)
		}
		tx, id, w := reject(t, reason)
		tcExpectReasonTooLong(t, w, "(3)")
		if s := tcRequestStatus(t, tx, id); s != "requested" {
			t.Errorf("(3) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
		}
	})
}

// tcExpectReasonTooLong ตรวจ 422 VALIDATION_FAILED พร้อมข้อความ "เหตุผลต้องไม่เกิน 500 ตัวอักษร" (Q-APV-28)
func tcExpectReasonTooLong(t *testing.T, w *httptest.ResponseRecorder, part string) {
	t.Helper()
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("%s ต้องได้ 422 แต่ได้ %d body=%s", part, w.Code, w.Body.String())
		return
	}
	if e := tcError(t, w); e.Code != "VALIDATION_FAILED" || e.Message != "เหตุผลต้องไม่เกิน 500 ตัวอักษร" {
		t.Errorf("%s ต้องได้ VALIDATION_FAILED พร้อมข้อความ \"เหตุผลต้องไม่เกิน 500 ตัวอักษร\" แต่ได้ %+v", part, e)
	}
}

// TC-APV-06-7 (AC-APV-06 ผิด)
func TestTC_APV_06_7_CancelledAndEmptyReason(t *testing.T) {
	// Given คำขอ cancelled ของ AP-TEST-0005 และ reason ""
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0005", "student-002")
	srv := tcServer(tx, nil)

	// When staff-001 เรียก POST .../reject
	w := tcReject(srv, tokenStaffDept1, id, "", time.Now())

	// Then (1) ได้ 422 VALIDATION_FAILED ข้อความ "กรุณาระบุเหตุผลในการปฏิเสธ" เพราะตรวจเหตุผลก่อนดูสถานะคำขอ (Q-APV-29)
	tcExpectReasonRequired(t, w)
	// Then (2) คำขอยังเป็น cancelled และไม่มีแถว approval_decision ของคำขอนั้น
	if s := tcRequestStatus(t, tx, id); s != "cancelled" {
		t.Errorf("(2) คำขอต้องยังเป็น cancelled แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(2) ต้องไม่มีแถว approval_decision แต่มี %d แถว", n)
	}
}

// TC-APV-06-8 (AC-APV-06 ผิด)
func TestTC_APV_06_8_StudentEmptyReasonForbidden(t *testing.T) {
	// Given student-001 (token mock-student-001) และ reason ""
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServer(tx, nil)

	// When เรียก POST .../reject ของคำขอ AP-TEST-0003
	w := tcReject(srv, tokenStudentDept1, id, "", time.Now())

	// Then (1) ได้ 403 FORBIDDEN_ROLE ข้อความ "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้"
	if w.Code != http.StatusForbidden {
		t.Errorf("(1) ต้องได้ 403 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if e := tcError(t, w); e.Code != "FORBIDDEN_ROLE" || e.Message != "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้" {
		t.Errorf("(1) ต้องได้ FORBIDDEN_ROLE พร้อมข้อความตาม FR-APV-01 (ซ) แต่ได้ %+v", e)
	}
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// ---------- AC-APV-11 ----------

// TC-APV-11-1 (AC-APV-11 ปกติ)
func TestTC_APV_11_1_ListOwnDepartmentPending(t *testing.T) {
	// Given ข้อมูล approve_dev.sql และเพิ่มคำขอ approved 1 คำขอของ AP-TEST-0005
	tx := devTx(t)
	if _, err := tx.Exec(context.Background(), `
		INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose,
		                            start_at, end_at, status, submitted_at)
		SELECT 'student-002', 'คำขอที่อนุมัติแล้ว', 'student-002', id, 'normal', 'study', now(), now(), 'approved',
		       now() - interval '50 minutes'
		FROM equipment_item WHERE asset_code = 'AP-TEST-0005'`); err != nil {
		t.Fatal(err)
	}
	srv := tcServer(tx, nil)

	// When staff-001 เรียก GET /v1/approvals/requests
	w, got := tcList(t, srv, tokenStaffDept1)

	// Then (1) ได้ 200 และ 6 รายการ (แก้ตาม spec v3 CR-009 2569-11-13)
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if len(got.Items) != 6 {
		t.Fatalf("(1) ต้องได้ 6 รายการ แต่ได้ %d: %+v", len(got.Items), got.Items)
	}
	// Then (2) ไม่มีคำขอ cancelled, คำขอ approved และคำขอของ AP-TEST-0101
	for _, it := range got.Items {
		switch {
		case it.AssetCode == "AP-TEST-0005":
			t.Errorf("(2) ต้องไม่มีคำขอ cancelled หรือ approved ของ AP-TEST-0005 แต่ได้ %+v", it)
		case it.AssetCode == "AP-TEST-0101":
			t.Error("(2) ต้องไม่มีคำขอของ AP-TEST-0101 (ภาควิชา 2)")
		}
	}
	// Then (3) เรียงตาม submittedAt จากเก่าไปใหม่: A, 0002, B, 0003, 0004, C
	want := []struct{ asset, name string }{
		{"AP-TEST-0001", "นักศึกษาตัวอย่าง ก"},
		{"AP-TEST-0002", "นักศึกษาตัวอย่าง ก"},
		{"AP-TEST-0001", "นักศึกษาตัวอย่าง ข"},
		{"AP-TEST-0003", "นักศึกษาตัวอย่าง ข"},
		{"AP-TEST-0004", "นักศึกษาตัวอย่าง ก"},
		{"AP-TEST-0001", "นักศึกษาตัวอย่าง ค"},
	}
	for i, w := range want {
		if got.Items[i].AssetCode != w.asset || got.Items[i].RequesterName != w.name {
			t.Errorf("(3) แถว %d ต้องเป็น %s ของ %s แต่ได้ %s ของ %s", i, w.asset, w.name, got.Items[i].AssetCode, got.Items[i].RequesterName)
		}
		if i > 0 && got.Items[i].SubmittedAt.Before(got.Items[i-1].SubmittedAt) {
			t.Errorf("(3) แถว %d ต้องเรียงตาม submittedAt จากเก่าไปใหม่", i)
		}
	}
	// Then (4) ทุกรายการมีชื่อผู้ขอ ชื่อรุ่น รหัสครุภัณฑ์ วัตถุประสงค์ startAt และ endAt
	for i, it := range got.Items {
		if it.RequesterName == "" || it.ModelName == "" || it.AssetCode == "" || it.Purpose == "" || it.StartAt.IsZero() || it.EndAt.IsZero() {
			t.Errorf("(4) แถว %d ข้อมูลไม่ครบ: %+v", i, it)
		}
	}
	// Then (5) มี loadedAt
	if got.LoadedAt.IsZero() {
		t.Error("(5) ต้องมี loadedAt")
	}
	// Then (6) C เป็น overnight isFirst=false มี attachmentUrl รายการอื่นเป็น normal ไม่มี attachmentUrl
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	for i, it := range got.Items {
		if it.RequestID == c {
			want := fmt.Sprintf("/v1/loan-requests/%d/attachment", c)
			if it.LoanType != "overnight" || it.IsFirst || it.AttachmentURL != want {
				t.Errorf("(6) C ต้องเป็น overnight isFirst=false attachmentUrl=%s แต่ได้ %+v", want, it)
			}
		} else if it.LoanType != "normal" || it.AttachmentURL != "" {
			t.Errorf("(6) แถว %d ต้องเป็น normal ไม่มี attachmentUrl แต่ได้ %+v", i, it)
		}
	}
}

// TC-APV-11-3 (AC-APV-11 ขอบ)
func TestTC_APV_11_3_DecidedRequestsDisappear(t *testing.T) {
	// Given staff-001 อนุมัติคำขอของ AP-TEST-0002 และปฏิเสธคำขอของ AP-TEST-0003 ก่อนเปิดรายการ
	tx := devTx(t)
	approved := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	rejected := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, approved))
	if w := tcApprove(srv, tokenStaffDept1, approved); w.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: อนุมัติต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if w := tcReject(srv, tokenStaffDept1, rejected, "อุปกรณ์ต้องใช้สอบ", time.Now()); w.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: ปฏิเสธต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}

	// When staff-001 เรียก GET /v1/approvals/requests
	w, got := tcList(t, srv, tokenStaffDept1)
	if w.Code != http.StatusOK {
		t.Fatalf("ต้องได้ 200 แต่ได้ %d", w.Code)
	}

	// Then (1) ได้ 4 รายการ (A, B, 0004, C) แก้ตาม spec v3 CR-009 2569-11-13
	if len(got.Items) != 4 {
		t.Errorf("(1) ต้องได้ 4 รายการ แต่ได้ %d: %+v", len(got.Items), got.Items)
	}
	// Then (2) ไม่มีคำขอของ AP-TEST-0002 และ 0003 แล้ว
	for _, it := range got.Items {
		if it.AssetCode == "AP-TEST-0002" || it.AssetCode == "AP-TEST-0003" {
			t.Errorf("(2) ต้องไม่มีคำขอที่พิจารณาแล้ว แต่ได้ %+v", it)
		}
	}
}

// TC-APV-11-4 (AC-APV-11 ขอบ)
func TestTC_APV_11_4_SameSubmittedAtOrder(t *testing.T) {
	// Given ภาควิชา 1 มีคำขอยืมปกติ requested 2 คำขอ ของอุปกรณ์ต่างชิ้นกัน (AP-TEST-0002 และ 0003) ที่ submitted_at เท่ากันพอดี
	tx := devTx(t)
	r2 := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	r3 := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	if _, err := tx.Exec(context.Background(),
		`UPDATE borrow_request SET submitted_at = date_trunc('second', now()) - interval '18 minutes' WHERE id IN ($1, $2)`, r2, r3); err != nil {
		t.Fatal(err)
	}
	low, high := r2, r3
	if high < low {
		low, high = high, low
	}
	srv := tcServer(tx, nil)

	// When staff-001 เรียก GET /v1/approvals/requests
	w, got := tcList(t, srv, tokenStaffDept1)
	if w.Code != http.StatusOK {
		t.Fatalf("ต้องได้ 200 แต่ได้ %d", w.Code)
	}

	// Then (1) ทั้ง 2 คำขออยู่ในรายการ
	pos := map[int]int{}
	for i, it := range got.Items {
		pos[it.RequestID] = i
	}
	pLow, okLow := pos[low]
	pHigh, okHigh := pos[high]
	if !okLow || !okHigh {
		t.Fatalf("(1) ทั้ง 2 คำขอต้องอยู่ในรายการ แต่ได้ %+v", got.Items)
	}
	if !got.Items[pLow].SubmittedAt.Equal(got.Items[pHigh].SubmittedAt) {
		t.Fatalf("เตรียมข้อมูล: submittedAt ต้องเท่ากัน แต่ได้ %v และ %v", got.Items[pLow].SubmittedAt, got.Items[pHigh].SubmittedAt)
	}
	// Then (2) คำขอที่รหัสคำขอน้อยกว่าอยู่ก่อน (Q-APV-30)
	if pLow > pHigh {
		t.Errorf("(2) คำขอรหัส %d ต้องอยู่ก่อนรหัส %d แต่อยู่ตำแหน่ง %d และ %d", low, high, pLow, pHigh)
	}
}

// TC-APV-11-6 (AC-APV-11 ผิด)
func TestTC_APV_11_6_OtherDepartmentSeesOwnOnly(t *testing.T) {
	// Given เจ้าหน้าที่ภาควิชา 2 (auth จำลองไม่มี token เจ้าหน้าที่ภาควิชา 2 จึงใส่ผู้ใช้ใน context ของ service โดยตรง)
	tx := devTx(t)

	// When เรียกรายการคำขอผ่าน service
	items, _, err := NewService(NewRepository(tx)).ListPending(staffCtx("staff-201", 2))
	if err != nil {
		t.Fatal(err)
	}

	// Then (1) ได้ 1 รายการ คือคำขอของ AP-TEST-0101
	if len(items) != 1 || items[0].AssetCode != "AP-TEST-0101" {
		t.Errorf("(1) ต้องได้ 1 รายการของ AP-TEST-0101 แต่ได้ %+v", items)
	}
	// Then (2) ไม่มีคำขอของ AP-TEST-0001 ถึง 0005
	for _, it := range items {
		switch it.AssetCode {
		case "AP-TEST-0001", "AP-TEST-0002", "AP-TEST-0003", "AP-TEST-0004", "AP-TEST-0005":
			t.Errorf("(2) ต้องไม่มีคำขอของภาควิชา 1 แต่ได้ %+v", it)
		}
	}
}

// TC-APV-11-7 (AC-APV-11 ผิด)
func TestTC_APV_11_7_StudentForbidden(t *testing.T) {
	// Given student-001
	srv := tcServer(devTx(t), nil)

	// When student-001 เรียก GET /v1/approvals/requests
	w, _ := tcList(t, srv, tokenStudentDept1)

	// Then (1) ได้ 403 FORBIDDEN_ROLE ข้อความ "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้"
	if w.Code != http.StatusForbidden {
		t.Errorf("(1) ต้องได้ 403 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if e := tcError(t, w); e.Code != "FORBIDDEN_ROLE" || e.Message != "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้" {
		t.Errorf("(1) ต้องได้ FORBIDDEN_ROLE พร้อมข้อความตาม FR-APV-01 (ซ) แต่ได้ %+v", e)
	}
	// Then (2) body ไม่มี items
	if tcHasItems(w) {
		t.Error("(2) body ต้องไม่มี items")
	}
}

// TC-APV-11-8 (AC-APV-11 ผิด)
func TestTC_APV_11_8_NoTokenUnauthenticated(t *testing.T) {
	// Given ไม่มี token
	srv := tcServer(devTx(t), nil)

	// When เรียก GET /v1/approvals/requests
	w, _ := tcList(t, srv, "")

	// Then (1) ได้ 401 UNAUTHENTICATED
	if w.Code != http.StatusUnauthorized {
		t.Errorf("(1) ต้องได้ 401 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if e := tcError(t, w); e.Code != "UNAUTHENTICATED" {
		t.Errorf("(1) code ต้องเป็น UNAUTHENTICATED แต่ได้ %+v", e)
	}
	// Then (2) body ไม่มี items
	if tcHasItems(w) {
		t.Error("(2) body ต้องไม่มี items")
	}
}
