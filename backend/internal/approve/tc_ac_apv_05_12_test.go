// test จากตาราง specs/004-approve-request/test-cases.md แถว TC-APV-05-x, TC-APV-12-x สถานะ "ใช้ได้" ห้ามแก้ให้ผ่าน
// ถ้าไม่ผ่าน ต้องหาว่าโค้ดผิด หรือ AC ผิด (ถ้า AC ผิด แก้ spec ก่อน)
// ใช้ helper จาก tc_ac_apv_01_test.go และ tc_ac_apv_02_06_11_test.go
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
)

// tcEmailCount นับแถว approval_email ของคำขอ
func tcEmailCount(t *testing.T, tx pgx.Tx, id int) int {
	t.Helper()
	var n int
	if err := tx.QueryRow(context.Background(), `SELECT count(*) FROM approval_email WHERE request_id = $1`, id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// tcExpectError ตรวจ HTTP status รหัส และข้อความ (ข้อความว่าง = ไม่ตรวจข้อความ แต่ต้องไม่ว่าง)
func tcExpectError(t *testing.T, part string, w *httptest.ResponseRecorder, status int, code, message string) {
	t.Helper()
	if w.Code != status {
		t.Errorf("%s ต้องได้ %d แต่ได้ %d body=%s", part, status, w.Code, w.Body.String())
		return
	}
	e := tcError(t, w)
	if e.Code != code {
		t.Errorf("%s code ต้องเป็น %s แต่ได้ %+v", part, code, e)
	}
	if message != "" && e.Message != message {
		t.Errorf("%s message ต้องเป็น %q แต่ได้ %q", part, message, e.Message)
	}
	if e.Message == "" {
		t.Errorf("%s message ต้องไม่ว่าง (CON-TECH-03)", part)
	}
}

const tcMsgForbiddenRole = "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้"
const tcMsgForbiddenDepartment = "คำขอนี้เป็นของอุปกรณ์ภาควิชาอื่น คุณไม่มีสิทธิ์พิจารณา"

// ---------- AC-APV-05 ----------

// TC-APV-05-1 (AC-APV-05 ปกติ)
func TestTC_APV_05_1_NoTokenApproveReject(t *testing.T) {
	// Given ไม่มี token คำขอ requested ของ AP-TEST-0002
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, nil)

	// When เรียก POST .../approve แล้วเรียก POST .../reject พร้อมเหตุผล
	wa := tcApprove(srv, "", id)
	wr := tcReject(srv, "", id, "อุปกรณ์ต้องใช้สอบ", time.Now())

	// Then (1) ได้ 401 UNAUTHENTICATED ทั้ง 2 ครั้ง (2) body เป็น {code, message} และ message ไม่ว่าง
	tcExpectError(t, "(1)(2) approve", wa, http.StatusUnauthorized, "UNAUTHENTICATED", "")
	tcExpectError(t, "(1)(2) reject", wr, http.StatusUnauthorized, "UNAUTHENTICATED", "")
	// Then (3) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(3) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	// Then (4) ไม่มีแถว approval_decision และไม่มีแถว approval_email ของคำขอนั้น
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(4) ต้องไม่มี approval_decision แต่มี %d", n)
	}
	if n := tcEmailCount(t, tx, id); n != 0 {
		t.Errorf("(4) ต้องไม่มี approval_email แต่มี %d", n)
	}
}

// TC-APV-05-2 (AC-APV-05 ปกติ)
func TestTC_APV_05_2_StudentForbidden(t *testing.T) {
	// Given student-001 ยืนยันตัวตนแล้ว คำขอ requested ของ AP-TEST-0003 (ของ student-002)
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	srv := tcServer(tx, nil)

	// When student-001 เรียก POST .../approve แล้วเรียก POST .../reject พร้อมเหตุผล
	wa := tcApprove(srv, tokenStudentDept1, id)
	wr := tcReject(srv, tokenStudentDept1, id, "อุปกรณ์ต้องใช้สอบ", time.Now())

	// Then (1) ได้ 403 FORBIDDEN_ROLE ข้อความ "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้" ทั้ง 2 ครั้ง
	tcExpectError(t, "(1) approve", wa, http.StatusForbidden, "FORBIDDEN_ROLE", tcMsgForbiddenRole)
	tcExpectError(t, "(1) reject", wr, http.StatusForbidden, "FORBIDDEN_ROLE", tcMsgForbiddenRole)
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	// Then (3) ไม่มีแถว approval_decision และไม่มีแถว approval_email ของคำขอนั้น
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_decision แต่มี %d", n)
	}
	if n := tcEmailCount(t, tx, id); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_email แต่มี %d", n)
	}
}

// TC-APV-05-4 (AC-APV-05 ขอบ)
func TestTC_APV_05_4_InvalidTokenUnauthenticated(t *testing.T) {
	// Given มี token แต่ใช้ไม่ได้ (mock-nobody) คำขอ requested ของ AP-TEST-0002
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, nil)

	// When เรียก POST .../approve
	w := tcApprove(srv, "mock-nobody", id)

	// Then (1) ได้ 401 UNAUTHENTICATED
	tcExpectError(t, "(1)", w, http.StatusUnauthorized, "UNAUTHENTICATED", "")
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-05-5 (AC-APV-05 ขอบ)
func TestTC_APV_05_5_RequesterCannotApproveOwn(t *testing.T) {
	// Given student-001 กับคำขอ requested ของตัวเองบน AP-TEST-0002
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, nil)

	// When student-001 เรียก POST .../approve
	w := tcApprove(srv, tokenStudentDept1, id)

	// Then (1) ได้ 403 FORBIDDEN_ROLE
	tcExpectError(t, "(1)", w, http.StatusForbidden, "FORBIDDEN_ROLE", tcMsgForbiddenRole)
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-05-6 (AC-APV-05 ขอบ)
func TestTC_APV_05_6_StudentOtherDepartmentRoleFirst(t *testing.T) {
	// Given student-101 (ภาควิชา 2) กับคำขอ requested ของภาควิชา 1 บน AP-TEST-0002
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, nil)

	// When student-101 เรียก POST .../approve
	w := tcApprove(srv, "mock-student-101", id)

	// Then (1) ได้ 403 FORBIDDEN_ROLE ไม่ใช่ FORBIDDEN_DEPARTMENT
	tcExpectError(t, "(1)", w, http.StatusForbidden, "FORBIDDEN_ROLE", tcMsgForbiddenRole)
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-05-7 (AC-APV-05 ขอบ)
func TestTC_APV_05_7_FailedEmailsRequireStaff(t *testing.T) {
	// Given endpoint อื่นของ 004: รายการ "Email ที่ส่งไม่สำเร็จ" และการส่งซ้ำ
	srv := tcServer(devTx(t), nil)
	call := func(token, method string) *httptest.ResponseRecorder {
		if method == http.MethodGet {
			return getJSON(srv, token, PathPrefix+"failed-emails")
		}
		return postJSON(srv, token, PathPrefix+"emails/1/resend", `{}`)
	}

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		// When student-001 เรียก แล้วเรียกโดยไม่มี token
		ws := call(tokenStudentDept1, method)
		wn := call("", method)

		// Then (1) student-001 ได้ 403 FORBIDDEN_ROLE
		tcExpectError(t, "(1) "+method, ws, http.StatusForbidden, "FORBIDDEN_ROLE", tcMsgForbiddenRole)
		// Then (2) ไม่มี token ได้ 401 UNAUTHENTICATED
		tcExpectError(t, "(2) "+method, wn, http.StatusUnauthorized, "UNAUTHENTICATED", "")
		// Then (3) body ไม่มีรายการ Email
		if tcHasItems(ws) || tcHasItems(wn) {
			t.Errorf("(3) %s body ต้องไม่มี items", method)
		}
	}
}

// ---------- AC-APV-12 ----------

// tcOnlyCodeMessage ตรวจว่า body มีแค่ code กับ message และไม่มีรหัสครุภัณฑ์
func tcOnlyCodeMessage(t *testing.T, part string, w *httptest.ResponseRecorder) {
	t.Helper()
	var raw map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	_, hasCode := raw["code"]
	_, hasMsg := raw["message"]
	if len(raw) != 2 || !hasCode || !hasMsg {
		t.Errorf("%s body ต้องมีแค่ code กับ message แต่ได้ %v", part, raw)
	}
	if strings.Contains(w.Body.String(), "AP-TEST-0101") {
		t.Errorf("%s body ต้องไม่มีรหัสครุภัณฑ์", part)
	}
}

// TC-APV-12-1 (AC-APV-12 ปกติ)
func TestTC_APV_12_1_ApproveOtherDepartment(t *testing.T) {
	// Given staff-001 (ภาควิชา 1) คำขอ requested ของ student-101 บน AP-TEST-0101 (ภาควิชา 2)
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0101", "student-101")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, id))

	// When staff-001 เรียก POST .../approve
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 403 FORBIDDEN_DEPARTMENT ข้อความ "คำขอนี้เป็นของอุปกรณ์ภาควิชาอื่น คุณไม่มีสิทธิ์พิจารณา"
	tcExpectError(t, "(1)", w, http.StatusForbidden, "FORBIDDEN_DEPARTMENT", tcMsgForbiddenDepartment)
	// Then (2) body มีแค่ code กับ message ไม่มีรหัสครุภัณฑ์หรือข้อมูลคำขอ
	tcOnlyCodeMessage(t, "(2)", w)
	// Then (3) สถานะคำขอยังเป็น requested และไม่มีแถว approval_decision
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(3) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_decision แต่มี %d", n)
	}
	// Then (4) equipment_item.status ของ AP-TEST-0101 ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0101"); s != "available" {
		t.Errorf("(4) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
}

// TC-APV-12-2 (AC-APV-12 ปกติ)
func TestTC_APV_12_2_RejectOtherDepartment(t *testing.T) {
	// Given เหมือน TC-APV-12-1
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0101", "student-101")
	srv := tcServer(tx, nil)

	// When staff-001 เรียก POST .../reject พร้อมเหตุผล "อุปกรณ์ต้องใช้สอบ"
	w := tcReject(srv, tokenStaffDept1, id, "อุปกรณ์ต้องใช้สอบ", time.Now())

	// Then (1) ได้ 403 FORBIDDEN_DEPARTMENT ข้อความเดียวกับ TC-APV-12-1
	tcExpectError(t, "(1)", w, http.StatusForbidden, "FORBIDDEN_DEPARTMENT", tcMsgForbiddenDepartment)
	// Then (2) body มีแค่ code กับ message
	tcOnlyCodeMessage(t, "(2)", w)
	// Then (3) สถานะคำขอยังเป็น requested และไม่มีแถว approval_decision
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(3) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(3) ต้องไม่มี approval_decision แต่มี %d", n)
	}
}

// TC-APV-12-4 (AC-APV-12 ขอบ)
func TestTC_APV_12_4_Department2CannotDecideDepartment1(t *testing.T) {
	// Given เจ้าหน้าที่ภาควิชา 2 (ใส่ผู้ใช้ใน context ของ service) กับคำขอ requested ของ AP-TEST-0002 (ภาควิชา 1)
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	svc := NewService(NewRepository(tx))
	svc.now = tcBeforeCutoff(t, tx, id).Now

	// When เรียกอนุมัติผ่าน service
	_, err := svc.Approve(staffCtx("staff-201", 2), id, time.Time{})

	// Then (1) ได้ error ภาควิชาอื่น (FORBIDDEN_DEPARTMENT)
	if !errors.Is(err, ErrForbiddenDepartment) {
		t.Errorf("(1) ต้องได้ ErrForbiddenDepartment แต่ได้ %v", err)
	}
	// Then (2) สถานะคำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) สถานะคำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-12-5 (AC-APV-12 ผิด)
func TestTC_APV_12_5_SameDepartmentAllowed(t *testing.T) {
	// Given เจ้าหน้าที่ภาควิชา 2 (ใส่ผู้ใช้ใน context ของ service) กับคำขอ requested ของ AP-TEST-0101 (ภาควิชา 2 คำขอแรก อุปกรณ์ "ว่าง")
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0101", "student-101")
	svc := NewService(NewRepository(tx))
	svc.now = tcBeforeCutoff(t, tx, id).Now

	// When เรียกอนุมัติผ่าน service
	_, err := svc.Approve(staffCtx("staff-201", 2), id, time.Time{})

	// Then (1) อนุมัติสำเร็จ คำขอเป็น approved
	if err != nil {
		t.Errorf("(1) ต้องอนุมัติได้ แต่ได้ %v", err)
	}
	if s := tcRequestStatus(t, tx, id); s != "approved" {
		t.Errorf("(1) คำขอต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (2) equipment_item.status ของ AP-TEST-0101 ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0101"); s != "available" {
		t.Errorf("(2) equipment_item.status ต้องยังเป็น available แต่ได้ %q", s)
	}
}

// TC-APV-12-6 (AC-APV-12 ผิด)
func TestTC_APV_12_6_OtherDepartmentAndCancelled(t *testing.T) {
	// Given คำขอของ AP-TEST-0101 (ภาควิชา 2) ถูกผู้ขอยกเลิกแล้ว staff-001 ภาควิชา 1
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0101", "student-101")
	if _, err := tx.Exec(context.Background(),
		`UPDATE borrow_request SET status = 'cancelled', cancel_reason = 'ไม่ต้องใช้แล้ว', cancelled_at = now() WHERE id = $1`, id); err != nil {
		t.Fatal(err)
	}
	srv := tcServer(tx, nil)

	// When staff-001 เรียก POST .../approve
	w := tcApprove(srv, tokenStaffDept1, id)

	// Then (1) ได้ 403 FORBIDDEN_DEPARTMENT เพราะตรวจภาควิชาก่อนสถานะ (Q-APV-33)
	tcExpectError(t, "(1)", w, http.StatusForbidden, "FORBIDDEN_DEPARTMENT", tcMsgForbiddenDepartment)
	// Then (2) คำขอยังเป็น cancelled และไม่มีแถว approval_decision
	if s := tcRequestStatus(t, tx, id); s != "cancelled" {
		t.Errorf("(2) คำขอต้องยังเป็น cancelled แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision แต่มี %d", n)
	}
}

// TC-APV-12-7 (AC-APV-12 ผิด)
func TestTC_APV_12_7_OtherDepartmentAndEmptyReason(t *testing.T) {
	// Given staff-001 ภาควิชา 1 กับคำขอ requested ของ AP-TEST-0101 (ภาควิชา 2) และ reason ""
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0101", "student-101")
	srv := tcServer(tx, nil)

	// When staff-001 เรียก POST .../reject
	w := tcReject(srv, tokenStaffDept1, id, "", time.Now())

	// Then (1) ได้ 403 FORBIDDEN_DEPARTMENT เพราะตรวจภาควิชาก่อนเหตุผล (Q-APV-34)
	tcExpectError(t, "(1)", w, http.StatusForbidden, "FORBIDDEN_DEPARTMENT", tcMsgForbiddenDepartment)
	// Then (2) คำขอยังเป็น requested และไม่มีแถว approval_decision
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) คำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(2) ต้องไม่มี approval_decision แต่มี %d", n)
	}
}
