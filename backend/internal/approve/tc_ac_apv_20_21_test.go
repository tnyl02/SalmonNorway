// test จากตาราง specs/004-approve-request/test-cases.md แถว TC-APV-20-x, TC-APV-21-x สถานะ "ใช้ได้" ห้ามแก้ให้ผ่าน
// ถ้าไม่ผ่าน ต้องหาว่าโค้ดผิด หรือ AC ผิด (ถ้า AC ผิด แก้ spec ก่อน)
// ใช้ helper จาก tc_ac_apv_*_test.go ไฟล์อื่น
package approve

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"salmonnorway/backend/internal/auth"
)

const tcMsgUpstream = "ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง"
const tcMsgNotFound = "ไม่พบคำขอนี้"

var errTCDatabaseDown = errors.New("ฐานข้อมูลจำลองไม่ตอบ")

// tcErrRow คือแถวที่ Scan แล้วได้ error เสมอ
type tcErrRow struct{ err error }

func (r tcErrRow) Scan(...any) error { return r.err }

// tcDownDB คือตัวเชื่อมฐานข้อมูลจำลองที่ตอบ error ทุกคำสั่ง (TC-APV-20-1, 20-8)
type tcDownDB struct{}

func (tcDownDB) Query(context.Context, string, ...any) (pgx.Rows, error) { return nil, errTCDatabaseDown }
func (tcDownDB) QueryRow(context.Context, string, ...any) pgx.Row       { return tcErrRow{errTCDatabaseDown} }
func (tcDownDB) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, errTCDatabaseDown
}
func (tcDownDB) Begin(context.Context) (pgx.Tx, error) { return nil, errTCDatabaseDown }

// tcBeginFailDB ครอบฐานข้อมูลจริง อ่านได้ แต่เปิด transaction ไม่ได้ (TC-APV-20-2, 20-3)
type tcBeginFailDB struct{ pgx.Tx }

func (tcBeginFailDB) Begin(context.Context) (pgx.Tx, error) { return nil, errTCDatabaseDown }

// tcFailInsertDB ครอบฐานข้อมูลจริง transaction เปิดได้ แต่คำสั่งที่มีข้อความ fail ล้ม (TC-APV-20-6)
type tcFailInsertDB struct {
	pgx.Tx
	fail string
}

func (d tcFailInsertDB) Begin(ctx context.Context) (pgx.Tx, error) {
	sp, err := d.Tx.Begin(ctx) // savepoint ใน transaction ของ test
	if err != nil {
		return nil, err
	}
	return tcFailInsertTx{Tx: sp, fail: d.fail}, nil
}

type tcFailInsertTx struct {
	pgx.Tx
	fail string
}

func (t tcFailInsertTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if strings.Contains(sql, t.fail) {
		return tcErrRow{errTCDatabaseDown}
	}
	return t.Tx.QueryRow(ctx, sql, args...)
}

func (t tcFailInsertTx) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if strings.Contains(sql, t.fail) {
		return pgconn.CommandTag{}, errTCDatabaseDown
	}
	return t.Tx.Exec(ctx, sql, args...)
}

// tcServerOn ประกอบ handler กับ auth.Middleware บนตัวเชื่อมฐานข้อมูลที่ให้
func tcServerOn(q Querier) http.Handler {
	return auth.Middleware(auth.NewMockProvider(), NewHandler(NewService(NewRepository(q))).Routes())
}

// tcCountAll นับแถวทั้งตาราง
func tcCountAll(t *testing.T, tx pgx.Tx, table string) int {
	t.Helper()
	var n int
	if err := tx.QueryRow(context.Background(), "SELECT count(*) FROM "+table).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// ---------- AC-APV-20 ----------

// TC-APV-20-1 (AC-APV-20 ปกติ)
func TestTC_APV_20_1_ListDatabaseDown(t *testing.T) {
	// Given staff-001 ยืนยันตัวตนแล้ว ฐานข้อมูลไม่ตอบ (ตัวเชื่อมจำลองตอบ error ทุกคำสั่ง)
	srv := tcServerOn(tcDownDB{})

	// When staff-001 เรียก GET /v1/approvals/requests
	w := getJSON(srv, tokenStaffDept1, PathPrefix+"requests")

	// Then (1) ได้ 503 UPSTREAM_UNAVAILABLE ข้อความตามที่กำหนด
	tcExpectError(t, "(1)", w, http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", tcMsgUpstream)
	// Then (2) body มีแค่ code กับ message ไม่มี items
	var raw map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if _, ok := raw["items"]; ok || len(raw) != 2 {
		t.Errorf("(2) body ต้องมีแค่ code กับ message แต่ได้ %v", raw)
	}
}

// TC-APV-20-2 (AC-APV-20 ปกติ)
func TestTC_APV_20_2_ApproveDatabaseDown(t *testing.T) {
	// Given คำขอ requested ของ AP-TEST-0002 อยู่ในฐานข้อมูล แต่การเปิด transaction ล้ม
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServerOn(tcBeginFailDB{tx})

	// When staff-001 เรียก POST .../approve
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 503 UPSTREAM_UNAVAILABLE ข้อความตามที่กำหนด
	tcExpectError(t, "(1)", w, http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", tcMsgUpstream)
	// Then (2) สถานะคำขอยังเป็น requested และไม่มีแถว approval_decision
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision แต่มี %d", n)
	}
}

// TC-APV-20-3 (AC-APV-20 ปกติ)
func TestTC_APV_20_3_RejectDatabaseDown(t *testing.T) {
	// Given เหมือน TC-APV-20-2 กับคำขอของ AP-TEST-0003
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServerOn(tcBeginFailDB{tx})

	// When staff-001 เรียก POST .../reject พร้อมเหตุผล
	w := tcReject(srv, tokenStaffDept1, id, "อุปกรณ์ต้องใช้สอบ", time.Now())

	// Then (1) ได้ 503 UPSTREAM_UNAVAILABLE ข้อความเดียวกัน
	tcExpectError(t, "(1)", w, http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", tcMsgUpstream)
	// Then (2) สถานะคำขอยังเป็น requested และไม่มีแถว approval_decision
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision แต่มี %d", n)
	}
}

// TC-APV-20-6 (AC-APV-20 ขอบ)
func TestTC_APV_20_6_FailMidTransactionRollsBack(t *testing.T) {
	// Given transaction เปิดได้และเปลี่ยนสถานะแล้ว แต่การบันทึก approval_decision ล้ม คำขอของ AP-TEST-0002
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	svc := NewService(NewRepository(tcFailInsertDB{Tx: tx, fail: "INSERT INTO approval_decision"}))
	svc.now = tcBeforeCutoff(t, tx, id).Now
	srv := auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes())

	// When staff-001 เรียก POST .../approve
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 503 UPSTREAM_UNAVAILABLE
	tcExpectError(t, "(1)", w, http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", tcMsgUpstream)
	// Then (2) สถานะคำขอยังเป็น requested เพราะยกเลิกทั้ง transaction
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	// Then (3) ไม่มีแถว approval_decision, approval_email และ audit_log action approve ของคำขอนั้น
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_decision แต่มี %d", n)
	}
	if n := tcEmailCount(t, tx, id); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_email แต่มี %d", n)
	}
	if n := tcAuditCount(t, tx, id, "approve"); n != 0 {
		t.Errorf("(3) ต้องไม่มี audit_log approve แต่มี %d", n)
	}
}

// TC-APV-20-8 (AC-APV-20 ผิด)
func TestTC_APV_20_8_NoTokenBeforeDatabase(t *testing.T) {
	// Given ไม่มี token ฐานข้อมูลไม่ตอบเหมือน TC-APV-20-1
	srv := tcServerOn(tcDownDB{})

	// When เรียก GET /v1/approvals/requests
	w := getJSON(srv, "", PathPrefix+"requests")

	// Then (1) ได้ 401 UNAUTHENTICATED ไม่ใช่ 503
	tcExpectError(t, "(1)", w, http.StatusUnauthorized, "UNAUTHENTICATED", "")
}

// ---------- AC-APV-21 ----------

// TC-APV-21-1 (AC-APV-21 ปกติ)
func TestTC_APV_21_1_ApproveUnknownRequest(t *testing.T) {
	// Given staff-001 ยืนยันตัวตนแล้ว ไม่มีคำขอรหัส 999999
	tx := devTx(t)
	srv := tcServer(tx, nil)

	// When staff-001 เรียก POST /v1/approvals/requests/999999/approve
	w := tcApprove(srv, tokenStaffDept1, 999999)

	// Then (1) ได้ 404 REQUEST_NOT_FOUND ข้อความ "ไม่พบคำขอนี้"
	tcExpectError(t, "(1)", w, http.StatusNotFound, "REQUEST_NOT_FOUND", tcMsgNotFound)
	// Then (2) body มีแค่ code กับ message
	var raw map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if len(raw) != 2 {
		t.Errorf("(2) body ต้องมีแค่ code กับ message แต่ได้ %v", raw)
	}
	// Then (3) ไม่มีแถว approval_decision และ approval_email เกิดขึ้น
	if n := tcCountAll(t, tx, "approval_decision"); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_decision แต่มี %d", n)
	}
	if n := tcCountAll(t, tx, "approval_email"); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_email แต่มี %d", n)
	}
}

// TC-APV-21-2 (AC-APV-21 ปกติ)
func TestTC_APV_21_2_RejectUnknownRequest(t *testing.T) {
	// Given เหมือน TC-APV-21-1
	tx := devTx(t)
	srv := tcServer(tx, nil)

	// When staff-001 เรียก POST .../999999/reject พร้อมเหตุผล
	w := tcReject(srv, tokenStaffDept1, 999999, "อุปกรณ์ต้องใช้สอบ", time.Now())

	// Then (1) ได้ 404 REQUEST_NOT_FOUND ข้อความ "ไม่พบคำขอนี้"
	tcExpectError(t, "(1)", w, http.StatusNotFound, "REQUEST_NOT_FOUND", tcMsgNotFound)
	// Then (2) ไม่มีแถว approval_decision เกิดขึ้น
	if n := tcCountAll(t, tx, "approval_decision"); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision แต่มี %d", n)
	}
}

// TC-APV-21-4 (AC-APV-21 ขอบ)
func TestTC_APV_21_4_NonNumericId(t *testing.T) {
	// Given รหัสคำขอที่ไม่ใช่ตัวเลข "abc"
	srv := tcServer(devTx(t), nil)

	// When staff-001 เรียก .../abc/approve และ .../abc/reject
	for _, action := range []string{"approve", "reject"} {
		w := postJSON(srv, tokenStaffDept1, PathPrefix+"requests/abc/"+action, `{"reason":"อุปกรณ์ต้องใช้สอบ"}`)

		// Then (1) ได้ 404 REQUEST_NOT_FOUND ข้อความ "ไม่พบคำขอนี้" ทั้ง 2 ครั้ง
		tcExpectError(t, "(1) "+action, w, http.StatusNotFound, "REQUEST_NOT_FOUND", tcMsgNotFound)
	}
}

// TC-APV-21-5 (AC-APV-21 ขอบ)
func TestTC_APV_21_5_ZeroAndNegativeId(t *testing.T) {
	// Given รหัสคำขอ 0 และ -1
	srv := tcServer(devTx(t), nil)

	// When staff-001 เรียก .../0/approve และ .../-1/approve
	for _, id := range []string{"0", "-1"} {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, PathPrefix+"requests/"+id+"/approve", strings.NewReader(`{}`))
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer "+tokenStaffDept1)
		srv.ServeHTTP(w, r)

		// Then (1) ได้ 404 REQUEST_NOT_FOUND ทั้ง 2 ครั้ง
		tcExpectError(t, "(1) id "+id, w, http.StatusNotFound, "REQUEST_NOT_FOUND", tcMsgNotFound)
	}
}

// TC-APV-21-6 (AC-APV-21 ผิด)
func TestTC_APV_21_6_UnknownAndEmptyReason(t *testing.T) {
	// Given รหัส 999999 และเหตุผลว่าง
	srv := tcServer(devTx(t), nil)

	// When staff-001 เรียก POST .../999999/reject ด้วย reason ""
	w := tcReject(srv, tokenStaffDept1, 999999, "", time.Now())

	// Then (1) ได้ 404 REQUEST_NOT_FOUND ไม่ใช่ 422 VALIDATION_FAILED
	tcExpectError(t, "(1)", w, http.StatusNotFound, "REQUEST_NOT_FOUND", tcMsgNotFound)
}

// TC-APV-21-7 (AC-APV-21 ผิด)
func TestTC_APV_21_7_StudentUnknownRequest(t *testing.T) {
	// Given student-001 กับรหัส 999999
	srv := tcServer(devTx(t), nil)

	// When student-001 เรียก POST .../999999/approve
	w := tcApprove(srv, tokenStudentDept1, 999999)

	// Then (1) ได้ 403 FORBIDDEN_ROLE ไม่ใช่ 404
	tcExpectError(t, "(1)", w, http.StatusForbidden, "FORBIDDEN_ROLE", tcMsgForbiddenRole)
}
