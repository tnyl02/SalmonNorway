package approve

// test จาก specs/004-approve-request/test-cases.md แถว TC-APV-22-* ถึง TC-APV-26-* (คำขอยืมค้างคืน FR-APV-03, CR-009)
// ห้ามแก้ test ในไฟล์นี้ให้ผ่าน (AGENTS.md) ถ้าไม่ผ่านให้หาว่าโค้ดผิด หรือแถวในตารางผิด

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

	"salmonnorway/backend/internal/config"
)

const (
	tcMsgAdvisorRequired = "กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ"
	tcMsgAdvisorTooLong  = "ชื่ออาจารย์ต้องไม่เกิน 100 ตัวอักษร"
	tcAdvisor            = "อ.สมชาย ใจดี"
)

// tcApproveAdvisor เรียก POST .../approve พร้อม loadedAt และ advisorName (advisor = nil คือไม่ส่งช่องนี้)
func tcApproveAdvisor(srv http.Handler, token string, id int, advisor *string, loadedAt time.Time) *httptest.ResponseRecorder {
	body := map[string]any{"loadedAt": loadedAt.UTC().Format(time.RFC3339)}
	if advisor != nil {
		body["advisorName"] = *advisor
	}
	b, _ := json.Marshal(body)
	return postJSON(srv, token, fmt.Sprintf("%srequests/%d/approve", PathPrefix, id), string(b))
}

func tcStr(s string) *string { return &s }

// tcOvernightC คืนรหัสคำขอค้างคืน C ของ AP-TEST-0001 ใน approve_dev.sql
func tcOvernightC(t *testing.T, tx pgx.Tx) int {
	t.Helper()
	return requestIDOf(t, tx, "AP-TEST-0001", "student-003")
}

// tcAddOvernight เพิ่มคำขอค้างคืน requested ของอุปกรณ์ asset ส่งเมื่อ ago ก่อน (มีไฟล์แนบตัวอย่าง)
func tcAddOvernight(t *testing.T, tx pgx.Tx, asset, requester string, ago time.Duration) int {
	t.Helper()
	var id int
	if err := tx.QueryRow(context.Background(), `
		INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose,
		                            start_at, end_at, status, submitted_at, attachment_path)
		SELECT $2, 'นักศึกษาค้างคืน '||$2, $2, id, 'overnight', 'activity', now(), now() + interval '1 day',
		       'requested', $3, 'dev/x.pdf'
		FROM equipment_item WHERE asset_code = $1 RETURNING id`, asset, requester, time.Now().Add(-ago)).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

// tcAdvisorName คืน approval_decision.advisor_name (nil = NULL)
func tcAdvisorName(t *testing.T, tx pgx.Tx, id int) *string {
	t.Helper()
	var name *string
	if err := tx.QueryRow(context.Background(), `SELECT advisor_name FROM approval_decision WHERE request_id = $1`, id).Scan(&name); err != nil {
		t.Fatal(err)
	}
	return name
}

// tcSteps คืน status ของ approval_step เรียงตามลำดับที่บันทึก
func tcSteps(t *testing.T, tx pgx.Tx, id int) []string {
	t.Helper()
	rows, err := tx.Query(context.Background(), `SELECT status FROM approval_step WHERE request_id = $1 ORDER BY id`, id)
	if err != nil {
		t.Fatal(err)
	}
	steps, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		t.Fatal(err)
	}
	return steps
}

// tcExpectAdvisorError ตรวจ 422 VALIDATION_FAILED fields ["advisorName"] พร้อมข้อความ
func tcExpectAdvisorError(t *testing.T, part string, w *httptest.ResponseRecorder, msg string) {
	t.Helper()
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("%s ต้องได้ 422 แต่ได้ %d body=%s", part, w.Code, w.Body.String())
		return
	}
	var e ValidationErrorDTO
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if e.Code != "VALIDATION_FAILED" || e.Message != msg || len(e.Fields) != 1 || e.Fields[0] != "advisorName" {
		t.Errorf("%s ต้องได้ VALIDATION_FAILED ช่อง advisorName ข้อความ %q แต่ได้ %+v", part, msg, e)
	}
}

// tcMailTo คืน Email ที่ส่งถึงผู้รับ to
func tcMailTo(m *fakeMailer, to string) []Email {
	var out []Email
	for _, e := range m.sent {
		if e.To == to {
			out = append(out, e)
		}
	}
	return out
}

// ---------- AC-APV-22 ----------

// TC-APV-22-1 (AC-APV-22 ปกติ)
func TestTC_APV_22_1_ApproveOvernightWithAdvisor(t *testing.T) {
	// Given คำขอค้างคืน C (student-003) ของ AP-TEST-0001 "ว่าง" ภาควิชา 1 สถานะ requested
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	srv := tcServer(tx, nil)

	// When staff-001 อนุมัติ C พร้อม advisorName = "อ.สมชาย ใจดี"
	w := tcApproveAdvisor(srv, tokenStaffDept1, c, tcStr(tcAdvisor), time.Now())

	// Then (1) ได้ 200 status = approved
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if r := tcApproveResult(t, w); r.Status != "approved" {
		t.Errorf("(1) status ต้องเป็น approved แต่ได้ %q", r.Status)
	}
	// Then (2) C เป็น approved ทันที
	if s := tcRequestStatus(t, tx, c); s != "approved" {
		t.Errorf("(2) C ต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (3) AP-TEST-0001 ยังเป็น available
	if s := tcItemStatus(t, tx, "AP-TEST-0001"); s != "available" {
		t.Errorf("(3) AP-TEST-0001 ต้องยังเป็น available แต่ได้ %q", s)
	}
	// Then (4) approval_decision มี advisor_name = "อ.สมชาย ใจดี" ผู้ดำเนินการ staff-001
	if n := tcAdvisorName(t, tx, c); n == nil || *n != tcAdvisor {
		t.Errorf("(4) advisor_name ต้องเป็น %q แต่ได้ %v", tcAdvisor, n)
	}
	if _, _, actor, _, ok := tcDecision(t, tx, c); !ok || actor != "staff-001" {
		t.Errorf("(4) ผู้ดำเนินการต้องเป็น staff-001 แต่ได้ %q (มีผล %v)", actor, ok)
	}
	// Then (5) approval_step 3 แถวตามลำดับ
	want := []string{"approved_by_officer", "pending_advisor_approval", "approved"}
	if got := tcSteps(t, tx, c); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("(5) approval_step ต้องเป็น %v แต่ได้ %v", want, got)
	}
	// Then (6) audit_log มี action approve ของ C
	if n := tcAuditCount(t, tx, c, "approve"); n != 1 {
		t.Errorf("(6) audit_log action approve ของ C ต้องมี 1 แถว แต่มี %d", n)
	}
}

// TC-APV-22-2 (AC-APV-22 ขอบ)
func TestTC_APV_22_2_AdvisorName100(t *testing.T) {
	// Given ชื่ออาจารย์ยาว 100 code point พอดี
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	name := strings.Repeat("ก", 100)

	// When staff-001 อนุมัติ C พร้อมชื่อนั้น
	w := tcApproveAdvisor(tcServer(tx, nil), tokenStaffDept1, c, &name, time.Now())

	// Then (1) ได้ 200 approved
	if w.Code != http.StatusOK || tcApproveResult(t, w).Status != "approved" {
		t.Fatalf("(1) ต้องได้ 200 approved แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	// Then (2) advisor_name เก็บครบ 100 code point
	if n := tcAdvisorName(t, tx, c); n == nil || utf8.RuneCountInString(*n) != 100 || *n != name {
		t.Errorf("(2) advisor_name ต้องเก็บครบ 100 code point แต่ได้ %v", n)
	}
}

// TC-APV-22-3 (AC-APV-22 ขอบ)
func TestTC_APV_22_3_AdvisorName101(t *testing.T) {
	// Given ชื่ออาจารย์ยาว 101 code point
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	name := strings.Repeat("ก", 101)

	// When staff-001 อนุมัติ C พร้อมชื่อนั้น
	w := tcApproveAdvisor(tcServer(tx, nil), tokenStaffDept1, c, &name, time.Now())

	// Then (1) ได้ 422 VALIDATION_FAILED fields ["advisorName"] ข้อความ "ชื่ออาจารย์ต้องไม่เกิน 100 ตัวอักษร"
	tcExpectAdvisorError(t, "(1)", w, tcMsgAdvisorTooLong)
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
	// Then (3) ไม่มี approval_step ของ C
	if got := tcSteps(t, tx, c); len(got) != 0 {
		t.Errorf("(3) ต้องไม่มี approval_step แต่ได้ %v", got)
	}
}

// TC-APV-22-4 (AC-APV-22 ขอบ)
func TestTC_APV_22_4_AdvisorNameSpaces(t *testing.T) {
	// Given ชื่ออาจารย์มีช่องว่างหน้าและหลัง "  อ.สมชาย  "
	tx := devTx(t)
	c := tcOvernightC(t, tx)

	// When staff-001 อนุมัติ C พร้อมชื่อนั้น
	w := tcApproveAdvisor(tcServer(tx, nil), tokenStaffDept1, c, tcStr("  อ.สมชาย  "), time.Now())

	// Then (1) ได้ 200 approved (ไม่ใช่ช่องว่างล้วน)
	if w.Code != http.StatusOK || tcApproveResult(t, w).Status != "approved" {
		t.Errorf("(1) ต้องได้ 200 approved แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	// Then (2) ค่าใน advisor_name เก็บตามที่พิมพ์หรือตัดช่องว่าง: ยังไม่ตรวจ เพราะ spec ไม่ได้บอก (คำถามในรายงาน /testcases 2569-11-13)
	// Then (3) ช่องว่างหน้าหลังนับรวมใน 100 ตัวอักษรไหม: ยังไม่ตรวจ เพราะ spec ไม่ได้บอก
}

// TC-APV-22-5 (AC-APV-22 ผิด)
func TestTC_APV_22_5_AdvisorNameMissing(t *testing.T) {
	// Given ไม่ส่ง advisorName และส่ง "   " (มีแต่ช่องว่าง)
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	srv := tcServer(tx, nil)

	for _, tc := range []struct {
		name    string
		advisor *string
	}{{"ไม่ส่ง advisorName", nil}, {"มีแต่ช่องว่าง", tcStr("   ")}} {
		// When staff-001 อนุมัติ C
		w := tcApproveAdvisor(srv, tokenStaffDept1, c, tc.advisor, time.Now())
		// Then (1) ได้ 422 VALIDATION_FAILED fields ["advisorName"] ข้อความ "กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ"
		tcExpectAdvisorError(t, "(1) "+tc.name, w, tcMsgAdvisorRequired)
	}
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-22-6 (AC-APV-22 ผิด)
func TestTC_APV_22_6_ItemNotAvailable(t *testing.T) {
	// Given คำขอค้างคืน requested ของ AP-TEST-0004 ที่ถูกยืมอยู่
	tx := devTx(t)
	id := tcAddOvernight(t, tx, "AP-TEST-0004", "student-104", time.Minute)

	// When staff-001 อนุมัติคำขอนั้นพร้อม advisorName
	w := tcApproveAdvisor(tcServer(tx, nil), tokenStaffDept1, id, tcStr(tcAdvisor), time.Now())

	// Then (1) ได้ 409 ITEM_NOT_AVAILABLE ข้อความ "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้"
	tcExpectError(t, "(1)", w, http.StatusConflict, "ITEM_NOT_AVAILABLE", tcMsgItemNotAvailable)
	// Then (2) คำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(2) คำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-22-7 (AC-APV-22 ผิด)
func TestTC_APV_22_7_OtherDepartmentFirst(t *testing.T) {
	// Given เจ้าหน้าที่ภาควิชา 2 (staff-201 ใส่ใน context ของ service โดยตรง) และไม่ส่งชื่ออาจารย์
	tx := devTx(t)
	c := tcOvernightC(t, tx)

	// When staff-201 อนุมัติ C
	_, err := NewService(NewRepository(tx)).ApproveWithAdvisor(staffCtx("staff-201", 2), c, time.Now(), "")

	// Then (1) ได้ FORBIDDEN_DEPARTMENT ไม่ใช่ VALIDATION_FAILED
	if _, e := tcMessageOf(err); e.Code != "FORBIDDEN_DEPARTMENT" {
		t.Errorf("(1) ต้องได้ FORBIDDEN_DEPARTMENT แต่ได้ %+v (err %v)", e, err)
	}
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-22-8 (AC-APV-22 ผิด)
func TestTC_APV_22_8_CancelledOvernight(t *testing.T) {
	// Given C ถูกยกเลิกแล้ว (cancelled)
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	tcSetStatus(t, tx, c, "cancelled")
	srv := tcServer(tx, nil)

	// When (ก) ไม่ส่งชื่ออาจารย์
	w := tcApproveAdvisor(srv, tokenStaffDept1, c, nil, time.Now())
	// Then (1) ได้ 422 VALIDATION_FAILED ช่อง advisorName
	tcExpectAdvisorError(t, "(1)", w, tcMsgAdvisorRequired)

	// When (ข) ส่งชื่ออาจารย์
	w = tcApproveAdvisor(srv, tokenStaffDept1, c, tcStr(tcAdvisor), time.Now())
	// Then (2) ได้ 409 REQUEST_CANCELLED ข้อความ "คำขอนี้ถูกผู้ขอยกเลิกแล้ว"
	tcExpectError(t, "(2)", w, http.StatusConflict, "REQUEST_CANCELLED", "คำขอนี้ถูกผู้ขอยกเลิกแล้ว")

	// Then (3) C ยังเป็น cancelled
	if s := tcRequestStatus(t, tx, c); s != "cancelled" {
		t.Errorf("(3) C ต้องยังเป็น cancelled แต่ได้ %q", s)
	}
}

// TC-APV-22-9 (AC-APV-22 ผิด)
func TestTC_APV_22_9_StudentForbidden(t *testing.T) {
	// Given student-001
	tx := devTx(t)
	c := tcOvernightC(t, tx)

	// When student-001 อนุมัติ C พร้อม advisorName
	w := tcApproveAdvisor(tcServer(tx, nil), tokenStudentDept1, c, tcStr(tcAdvisor), time.Now())

	// Then (1) ได้ 403 FORBIDDEN_ROLE ข้อความ "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้"
	tcExpectError(t, "(1)", w, http.StatusForbidden, "FORBIDDEN_ROLE", "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้")
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-22-12 (AC-APV-22 ปกติ)
func TestTC_APV_22_12_ApprovedEmail(t *testing.T) {
	// Given ตัวส่ง Email จำลองและนาฬิกาจำลอง หลังอนุมัติ C แบบ 22-1
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	env := tcNewEmailEnv(t, tx, time.Now())
	if _, err := env.svc.ApproveWithAdvisor(staffCtx("staff-001", 1), c, env.clock.Now(), tcAdvisor); err != nil {
		t.Fatal(err)
	}
	decidedAt := env.clock.Now()

	// When ตัวส่งตรวจคิวรอบถัดไป
	env.clock.Advance(config.ApprovalSenderPollInterval)
	env.sendDue(t)

	// Then (1) ส่ง Email ถึง student-003 ภายใน 60 วินาที หัวข้อ "แจ้งผลการอนุมัติการยืมอุปกรณ์ภาคคอมพิวเตอร์"
	mails := tcMailTo(env.mailer, "student-003@example.com")
	if len(mails) == 0 {
		t.Fatalf("(1) ต้องมี Email ถึง student-003 แต่ไม่มี (ส่งทั้งหมด %d ฉบับ)", len(env.mailer.sent))
	}
	if mails[0].Subject != "แจ้งผลการอนุมัติการยืมอุปกรณ์ภาคคอมพิวเตอร์" {
		t.Errorf("(1) หัวข้อไม่ถูกต้อง: %q", mails[0].Subject)
	}
	if d := env.clock.Now().Sub(decidedAt); d > 60*time.Second {
		t.Errorf("(1) ต้องส่งภายใน 60 วินาที แต่ใช้ %v", d)
	}
	// Then (2) มีสถานที่รับ
	tcContainsAll(t, "(2)", mails[0].Body, "คณะวิทยาศาสตร์ ตึก 1 ชั้น 2 ห้องอุปกรณ์ภาควิชาคอมพิวเตอร์")
	// Then (3) Email ของ C ส่ง 1 ฉบับ ไม่มีฉบับของขั้นระหว่างทาง
	if len(mails) != 1 || tcEmailCount(t, tx, c) != 1 {
		t.Errorf("(3) Email ของ C ต้องมี 1 ฉบับ แต่ส่ง %d ฉบับ และในคิว %d ฉบับ", len(mails), tcEmailCount(t, tx, c))
	}
	// Then (4) Email ของค้างคืนแสดงจำนวนอุปกรณ์ไหม: ยังไม่ตรวจ เพราะ spec ไม่ได้บอก (คำถามในรายงาน /testcases 2569-11-13)
	// Then (5) Email แสดงชื่ออาจารย์ไหม: ยังไม่ตรวจ เพราะ spec ไม่ได้บอก
}

// ---------- AC-APV-23 ----------

// tcXY เพิ่มคำขอค้างคืน X (ส่ง 3 ชั่วโมงก่อน) และ Y (1 นาทีก่อน) ของ AP-TEST-0003 คืน X, Y และคำขอปกติเดิมของ student-002
func tcXY(t *testing.T, tx pgx.Tx) (x, y, normal int) {
	t.Helper()
	x = tcAddOvernight(t, tx, "AP-TEST-0003", "student-101", 3*time.Hour)
	y = tcAddOvernight(t, tx, "AP-TEST-0003", "student-102", time.Minute)
	return x, y, requestIDOf(t, tx, "AP-TEST-0003", "student-002")
}

// tcExpectSuperseded ตรวจว่าคำขอถูกปฏิเสธตาม เหตุผลของ FR-APV-01 (จ) ผู้ดำเนินการ staff-001
func tcExpectSuperseded(t *testing.T, part string, tx pgx.Tx, id int) {
	t.Helper()
	result, reason, actor, _, ok := tcDecision(t, tx, id)
	if s := tcRequestStatus(t, tx, id); s != "rejected" || !ok || result != "rejected" || reason != tcMsgSuperseded || actor != "staff-001" {
		t.Errorf("%s คำขอ %d ต้องถูกปฏิเสธตามโดย staff-001 แต่ได้ สถานะ %q ผล %q เหตุผล %q ผู้ดำเนินการ %q", part, id, s, result, reason, actor)
	}
}

// TC-APV-23-1 (AC-APV-23 ปกติ)
func TestTC_APV_23_1_ApproveLaterOvernight(t *testing.T) {
	// Given AP-TEST-0003 มีค้างคืน X (ส่งก่อน) และ Y และคำขอปกติของ student-002
	tx := devTx(t)
	x, y, normal := tcXY(t, tx)

	// When staff-001 อนุมัติ Y พร้อม advisorName
	w := tcApproveAdvisor(tcServer(tx, nil), tokenStaffDept1, y, tcStr(tcAdvisor), time.Now())

	// Then (1) ได้ 200 approved ไม่ได้ NOT_FIRST_IN_QUEUE
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	res := tcApproveResult(t, w)
	if res.Status != "approved" {
		t.Errorf("(1) status ต้องเป็น approved แต่ได้ %q", res.Status)
	}
	// Then (2) X เป็น rejected เหตุผลปฏิเสธตาม ผู้ดำเนินการ staff-001
	tcExpectSuperseded(t, "(2)", tx, x)
	// Then (3) คำขอปกติของ AP-TEST-0003 เป็น rejected เหตุผลเดียวกัน
	tcExpectSuperseded(t, "(3)", tx, normal)
	// Then (4) rejectedRequestIds มี X และคำขอปกติ
	if !tcContains(res.RejectedRequestIDs, x) || !tcContains(res.RejectedRequestIDs, normal) {
		t.Errorf("(4) rejectedRequestIds ต้องมี %d และ %d แต่ได้ %v", x, normal, res.RejectedRequestIDs)
	}
	// Then (5) X มี Email ผลปฏิเสธในคิว
	if k := tcEmailKind(t, tx, x); k != "rejected" {
		t.Errorf("(5) Email ของ X ต้องเป็นแบบ rejected แต่ได้ %q", k)
	}
}

// TC-APV-23-2 (AC-APV-23 ขอบ)
func TestTC_APV_23_2_ApproveEarlierOvernight(t *testing.T) {
	// Given เหมือน 23-1
	tx := devTx(t)
	x, y, normal := tcXY(t, tx)

	// When staff-001 อนุมัติ X (ส่งก่อน) แทน
	w := tcApproveAdvisor(tcServer(tx, nil), tokenStaffDept1, x, tcStr(tcAdvisor), time.Now())

	// Then (1) ได้ 200 approved
	if w.Code != http.StatusOK || tcApproveResult(t, w).Status != "approved" {
		t.Fatalf("(1) ต้องได้ 200 approved แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	// Then (2) Y และคำขอปกติของ AP-TEST-0003 เป็น rejected เหตุผลปฏิเสธตาม
	tcExpectSuperseded(t, "(2) Y", tx, y)
	tcExpectSuperseded(t, "(2) คำขอปกติ", tx, normal)
}

// TC-APV-23-3 (AC-APV-23 ผิด)
func TestTC_APV_23_3_SupersededCannotApprove(t *testing.T) {
	// Given อนุมัติ Y แบบ 23-1 ไปแล้ว (X ถูกปฏิเสธตาม)
	tx := devTx(t)
	x, y, _ := tcXY(t, tx)
	srv := tcServer(tx, nil)
	if w := tcApproveAdvisor(srv, tokenStaffDept1, y, tcStr(tcAdvisor), time.Now()); w.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: อนุมัติ Y ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}

	// When staff-001 อนุมัติ X ด้วย loadedAt หลังเวลาที่ Y ได้อนุมัติ
	w := tcApproveAdvisor(srv, tokenStaffDept1, x, tcStr(tcAdvisor), time.Now().Add(time.Minute))

	// Then (1) ได้ 409 REQUEST_ALREADY_DECIDED ข้อความ "คำขอนี้ได้รับการพิจารณาแล้ว"
	tcExpectError(t, "(1)", w, http.StatusConflict, "REQUEST_ALREADY_DECIDED", tcMsgAlreadyDecided)
	// Then (2) X ยังเป็น rejected Y ยังเป็น approved
	if sx, sy := tcRequestStatus(t, tx, x), tcRequestStatus(t, tx, y); sx != "rejected" || sy != "approved" {
		t.Errorf("(2) X ต้องเป็น rejected และ Y เป็น approved แต่ได้ %q / %q", sx, sy)
	}
}

// ---------- AC-APV-24 ----------

// TC-APV-24-1 (AC-APV-24 ปกติ)
func TestTC_APV_24_1_OvernightSupersedesNormal(t *testing.T) {
	// Given AP-TEST-0001 มีคำขอปกติ A, B และค้างคืน C ตัวส่ง Email จำลองและนาฬิกาจำลอง
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	env := tcNewEmailEnv(t, tx, time.Now())

	// When staff-001 อนุมัติ C พร้อม advisorName แล้วตัวส่งตรวจคิวรอบถัดไป
	w := tcApproveAdvisor(env.srv, tokenStaffDept1, c, tcStr(tcAdvisor), env.clock.Now())
	if w.Code != http.StatusOK {
		t.Fatalf("อนุมัติ C ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	res := tcApproveResult(t, w)
	decidedAt := env.clock.Now()
	env.clock.Advance(config.ApprovalSenderPollInterval)
	env.sendDue(t)

	// Then (1) C เป็น approved
	if s := tcRequestStatus(t, tx, c); s != "approved" {
		t.Errorf("(1) C ต้องเป็น approved แต่ได้ %q", s)
	}
	// Then (2) A และ B เป็น rejected เหตุผลปฏิเสธตาม ผู้ดำเนินการ staff-001
	tcExpectSuperseded(t, "(2) A", tx, a)
	tcExpectSuperseded(t, "(2) B", tx, b)
	// Then (3) rejectedRequestIds มี A และ B
	if !tcContains(res.RejectedRequestIDs, a) || !tcContains(res.RejectedRequestIDs, b) {
		t.Errorf("(3) rejectedRequestIds ต้องมี %d และ %d แต่ได้ %v", a, b, res.RejectedRequestIDs)
	}
	// Then (4) student-001 ได้ Email ภายใน 60 วินาที หัวข้อผลปฏิเสธ มีเหตุผลปฏิเสธตาม
	mails := tcMailTo(env.mailer, "student-001@example.com")
	if len(mails) != 1 {
		t.Fatalf("(4) student-001 ต้องได้ Email 1 ฉบับ แต่ได้ %d", len(mails))
	}
	if mails[0].Subject != "แจ้งผลการปฏิเสธการยืมอุปกรณ์ภาคคอมพิวเตอร์" {
		t.Errorf("(4) หัวข้อไม่ถูกต้อง: %q", mails[0].Subject)
	}
	tcContainsAll(t, "(4)", mails[0].Body, tcMsgSuperseded)
	if d := env.clock.Now().Sub(decidedAt); d > 60*time.Second {
		t.Errorf("(4) ต้องส่งภายใน 60 วินาที แต่ใช้ %v", d)
	}
}

// TC-APV-24-2 (AC-APV-24 ขอบ)
func TestTC_APV_24_2_NormalSupersedesOvernight(t *testing.T) {
	// Given อนุมัติ A (คำขอแรกของยืมปกติ) ก่อน 16:00 น. ของวันที่เริ่มยืม (นาฬิกาจำลอง)
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, a))

	// When staff-001 อนุมัติ A
	w := tcApprove(srv, tokenStaffDept1, a)

	// Then (1) A เป็น approved
	if w.Code != http.StatusOK || tcRequestStatus(t, tx, a) != "approved" {
		t.Fatalf("(1) A ต้องเป็น approved แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	// Then (2) C และ B เป็น rejected เหตุผลปฏิเสธตาม
	tcExpectSuperseded(t, "(2) C", tx, c)
	tcExpectSuperseded(t, "(2) B", tx, b)
	// Then (3) C มี Email ผลปฏิเสธในคิว
	if k := tcEmailKind(t, tx, c); k != "rejected" {
		t.Errorf("(3) Email ของ C ต้องเป็นแบบ rejected แต่ได้ %q", k)
	}
}

// TC-APV-24-3 (AC-APV-24 ผิด)
func TestTC_APV_24_3_SupersededNormalCannotApprove(t *testing.T) {
	// Given อนุมัติ C แบบ 24-1 ไปแล้ว
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	srv := tcServer(tx, nil)
	if w := tcApproveAdvisor(srv, tokenStaffDept1, c, tcStr(tcAdvisor), time.Now()); w.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: อนุมัติ C ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}

	// When staff-001 อนุมัติ A ด้วย loadedAt หลังเวลาที่ C ได้อนุมัติ
	w := tcApproveWithLoaded(srv, a, time.Now().Add(time.Minute))

	// Then (1) ได้ 409 REQUEST_ALREADY_DECIDED ข้อความ "คำขอนี้ได้รับการพิจารณาแล้ว"
	tcExpectError(t, "(1)", w, http.StatusConflict, "REQUEST_ALREADY_DECIDED", tcMsgAlreadyDecided)
	// Then (2) A ยังเป็น rejected
	if s := tcRequestStatus(t, tx, a); s != "rejected" {
		t.Errorf("(2) A ต้องยังเป็น rejected แต่ได้ %q", s)
	}
}

// ---------- AC-APV-25 ----------

// TC-APV-25-1 (AC-APV-25 ปกติ)
func TestTC_APV_25_1_RejectOvernight(t *testing.T) {
	const reason = "เอกสารไม่มีลายเซ็นอาจารย์"
	// Given C สถานะ requested ตัวส่ง Email จำลองและนาฬิกาจำลอง
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	env := tcNewEmailEnv(t, tx, time.Now())

	// When staff-001 ปฏิเสธ C พร้อมเหตุผล (ไม่ส่งชื่ออาจารย์) แล้วตัวส่งตรวจคิวรอบถัดไป
	w := tcReject(env.srv, tokenStaffDept1, c, reason, env.clock.Now())
	decidedAt := env.clock.Now()
	env.clock.Advance(config.ApprovalSenderPollInterval)
	env.sendDue(t)

	// Then (1) ได้ 200 status = rejected
	if w.Code != http.StatusOK || tcApproveResult(t, w).Status != "rejected" {
		t.Fatalf("(1) ต้องได้ 200 rejected แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	// Then (2) C เป็น rejected ประวัติมีเหตุผล ผู้ดำเนินการ staff-001 ไม่มี advisor_name
	result, gotReason, actor, _, ok := tcDecision(t, tx, c)
	if s := tcRequestStatus(t, tx, c); s != "rejected" || !ok || result != "rejected" || gotReason != reason || actor != "staff-001" {
		t.Errorf("(2) ได้ สถานะ %q ผล %q เหตุผล %q ผู้ดำเนินการ %q", s, result, gotReason, actor)
	}
	if n := tcAdvisorName(t, tx, c); n != nil {
		t.Errorf("(2) ต้องไม่มี advisor_name แต่ได้ %q", *n)
	}
	// Then (3) ไม่มี approval_step ของ C
	if got := tcSteps(t, tx, c); len(got) != 0 {
		t.Errorf("(3) ต้องไม่มี approval_step แต่ได้ %v", got)
	}
	// Then (4) A และ B ยังเป็น requested
	for _, who := range []string{"student-001", "student-002"} {
		if s := tcRequestStatus(t, tx, requestIDOf(t, tx, "AP-TEST-0001", who)); s != "requested" {
			t.Errorf("(4) คำขอของ %s ต้องยังเป็น requested แต่ได้ %q", who, s)
		}
	}
	// Then (5) student-003 ได้ Email ภายใน 60 วินาที หัวข้อผลปฏิเสธ มีเหตุผล
	mails := tcMailTo(env.mailer, "student-003@example.com")
	if len(mails) != 1 {
		t.Fatalf("(5) student-003 ต้องได้ Email 1 ฉบับ แต่ได้ %d", len(mails))
	}
	if mails[0].Subject != "แจ้งผลการปฏิเสธการยืมอุปกรณ์ภาคคอมพิวเตอร์" {
		t.Errorf("(5) หัวข้อไม่ถูกต้อง: %q", mails[0].Subject)
	}
	tcContainsAll(t, "(5)", mails[0].Body, "เหตุผลในการปฏิเสธ: "+reason)
	if d := env.clock.Now().Sub(decidedAt); d > 60*time.Second {
		t.Errorf("(5) ต้องส่งภายใน 60 วินาที แต่ใช้ %v", d)
	}
}

// TC-APV-25-2 (AC-APV-25 ขอบ)
func TestTC_APV_25_2_Reason500(t *testing.T) {
	// Given เหตุผลยาว 500 code point
	tx := devTx(t)
	c := tcOvernightC(t, tx)

	// When staff-001 ปฏิเสธ C
	w := tcReject(tcServer(tx, nil), tokenStaffDept1, c, strings.Repeat("ก", 500), time.Now())

	// Then (1) ได้ 200 rejected
	if w.Code != http.StatusOK || tcApproveResult(t, w).Status != "rejected" {
		t.Errorf("(1) ต้องได้ 200 rejected แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
}

// TC-APV-25-3 (AC-APV-25 ขอบ)
func TestTC_APV_25_3_Reason501(t *testing.T) {
	// Given เหตุผลยาว 501 code point
	tx := devTx(t)
	c := tcOvernightC(t, tx)

	// When staff-001 ปฏิเสธ C
	w := tcReject(tcServer(tx, nil), tokenStaffDept1, c, strings.Repeat("ก", 501), time.Now())

	// Then (1) ได้ 422 VALIDATION_FAILED fields ["reason"] ข้อความ "เหตุผลต้องไม่เกิน 500 ตัวอักษร"
	tcExpectReasonTooLong(t, w, "(1)")
	var e ValidationErrorDTO
	_ = json.Unmarshal(w.Body.Bytes(), &e)
	if len(e.Fields) != 1 || e.Fields[0] != "reason" {
		t.Errorf("(1) fields ต้องเป็น [reason] แต่ได้ %v", e.Fields)
	}
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-25-4 (AC-APV-25 ผิด)
func TestTC_APV_25_4_ReasonMissing(t *testing.T) {
	// Given เหตุผลว่าง และมีแต่ช่องว่าง
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	srv := tcServer(tx, nil)

	for _, reason := range []string{"", "   "} {
		// When staff-001 ปฏิเสธ C
		w := tcReject(srv, tokenStaffDept1, c, reason, time.Now())
		// Then (1) ได้ 422 VALIDATION_FAILED fields ["reason"] ข้อความ "กรุณาระบุเหตุผลในการปฏิเสธ"
		tcExpectReasonRequired(t, w)
		var e ValidationErrorDTO
		_ = json.Unmarshal(w.Body.Bytes(), &e)
		if len(e.Fields) != 1 || e.Fields[0] != "reason" {
			t.Errorf("(1) เหตุผล %q fields ต้องเป็น [reason] แต่ได้ %v", reason, e.Fields)
		}
	}
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-25-5 (AC-APV-25 ผิด)
func TestTC_APV_25_5_RejectTwice(t *testing.T) {
	// Given ปฏิเสธ C ไปแล้ว
	tx := devTx(t)
	c := tcOvernightC(t, tx)
	srv := tcServer(tx, nil)
	if w := tcReject(srv, tokenStaffDept1, c, "เอกสารไม่ครบ", time.Now()); w.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: ปฏิเสธครั้งแรกต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}

	// When staff-001 ปฏิเสธ C อีกครั้งด้วย loadedAt หลังเวลาที่ปฏิเสธครั้งแรก
	w := tcReject(srv, tokenStaffDept1, c, "เอกสารไม่ครบ", time.Now().Add(time.Minute))

	// Then (1) ได้ 409 REQUEST_ALREADY_DECIDED ข้อความ "คำขอนี้ได้รับการพิจารณาแล้ว"
	tcExpectError(t, "(1)", w, http.StatusConflict, "REQUEST_ALREADY_DECIDED", tcMsgAlreadyDecided)
	// Then (2) มี Email ผลปฏิเสธของ C ฉบับเดียว
	if n := tcEmailCount(t, tx, c); n != 1 {
		t.Errorf("(2) Email ของ C ต้องมี 1 ฉบับ แต่มี %d", n)
	}
}

// ---------- AC-APV-26 ----------

// tcLateOvernight ตั้ง C ให้ส่งเมื่อวาน 08:00 น. เริ่มยืมเมื่อวาน 09:00 น. เวลาไทย คืน C และวันนี้ 17:00 น.
func tcLateOvernight(t *testing.T, tx pgx.Tx) (int, time.Time) {
	t.Helper()
	c := tcOvernightC(t, tx)
	today := autoRejectCutoff(time.Now()) // 00:00 น. ของวันนี้ เวลาไทย
	yesterday := today.AddDate(0, 0, -1)
	if _, err := tx.Exec(context.Background(), `UPDATE borrow_request SET submitted_at = $2, start_at = $3, end_at = $4 WHERE id = $1`,
		c, yesterday.Add(8*time.Hour), yesterday.Add(9*time.Hour), today.Add(16*time.Hour)); err != nil {
		t.Fatal(err)
	}
	return c, today.Add(17 * time.Hour)
}

// TC-APV-26-1 (AC-APV-26 ปกติ)
func TestTC_APV_26_1_NoCutoffNoAutoReject(t *testing.T) {
	// Given C ส่งเมื่อวาน 08:00 น. เริ่มยืมเมื่อวาน 09:00 น. นาฬิกาจำลองคือวันนี้ 17:00 น.
	tx := devTx(t)
	c, now := tcLateOvernight(t, tx)

	// When งานเที่ยงคืน (RejectStale) และงานไล่ตาม (CatchUp) ทำงาน
	tcRunMidnight(t, tx, now)
	tcCatchUp(tx, now)

	// Then (1) C ยังเป็น requested ไม่มีผลการพิจารณา
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(1) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
	if n := tcDecisionCount(t, tx, c); n != 0 {
		t.Errorf("(1) ต้องไม่มีผลการพิจารณาของ C แต่มี %d", n)
	}

	// When staff-001 อนุมัติ C พร้อม advisorName
	w := tcApproveAdvisor(tcServer(tx, &fakeClock{t: now}), tokenStaffDept1, c, tcStr(tcAdvisor), now)

	// Then (2) ได้ 200 approved ไม่ได้ APPROVAL_CUTOFF_PASSED
	if w.Code != http.StatusOK || tcApproveResult(t, w).Status != "approved" {
		t.Errorf("(2) ต้องได้ 200 approved แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
}

// TC-APV-26-2 (AC-APV-26 ขอบ)
func TestTC_APV_26_2_At1600(t *testing.T) {
	// Given C และคำขอปกติ A ของชิ้นเดียวกันเริ่มยืมวันนี้ นาฬิกาจำลองคือ 16:00:00 น. ตรง (ข้อมูลชุดใหม่ทีละแบบ)

	// แยก subtest เพื่อให้ transaction ของแต่ละแบบปิดก่อนเปิดชุดใหม่ (devTx 2 ชุดพร้อมกันจะรอล็อกของข้อมูลตั้งต้นกันเอง)
	t.Run("A ยืมปกติ", func(t *testing.T) {
		// When staff-001 อนุมัติ A
		tx := devTx(t)
		a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
		w := tcApprove(tcServer(tx, &fakeClock{t: tcAt(t, tx, a, 16, 0, 0)}), tokenStaffDept1, a)
		// Then (1) A ได้ 409 APPROVAL_CUTOFF_PASSED ข้อความ "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้"
		tcExpectError(t, "(1)", w, http.StatusConflict, "APPROVAL_CUTOFF_PASSED", tcMsgCutoff)
	})
	t.Run("C ค้างคืน", func(t *testing.T) {
		// When staff-001 อนุมัติ C
		tx := devTx(t)
		c := tcOvernightC(t, tx)
		at := tcAt(t, tx, c, 16, 0, 0)
		w := tcApproveAdvisor(tcServer(tx, &fakeClock{t: at}), tokenStaffDept1, c, tcStr(tcAdvisor), at)
		// Then (2) C ได้ 200 approved
		if w.Code != http.StatusOK || tcApproveResult(t, w).Status != "approved" {
			t.Errorf("(2) C ต้องได้ 200 approved แต่ได้ %d body=%s", w.Code, w.Body.String())
		}
	})
}

// TC-APV-26-3 (AC-APV-26 ขอบ)
func TestTC_APV_26_3_MidnightSkipsOvernight(t *testing.T) {
	// Given A (ปกติ) และ C (ค้างคืน) ส่งวันที่ D นาฬิกาจำลองคือ 00:00:00 น. หลังวันที่ D
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	c := tcOvernightC(t, tx)
	midnight := tcMidnight()
	setSubmittedAt(t, tx, a, midnight.Add(-2*time.Hour))
	setSubmittedAt(t, tx, c, midnight.Add(-2*time.Hour))

	// When งานเที่ยงคืน RejectStale ทำงาน
	tcRunMidnight(t, tx, midnight)

	// Then (1) A เป็น rejected ผู้ดำเนินการ System เหตุผลปฏิเสธอัตโนมัติ
	result, reason, actor, _, ok := tcDecision(t, tx, a)
	if s := tcRequestStatus(t, tx, a); s != "rejected" || !ok || result != "auto_rejected" || reason != tcAutoReason || actor != config.ApprovalSystemActorID {
		t.Errorf("(1) A ได้ สถานะ %q ผล %q เหตุผล %q ผู้ดำเนินการ %q", s, result, reason, actor)
	}
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-26-4 (AC-APV-26 ผิด)
func TestTC_APV_26_4_LateWithoutAdvisor(t *testing.T) {
	// Given สภาพเวลาเดียวกับ 26-1 และไม่ส่งชื่ออาจารย์
	tx := devTx(t)
	c, now := tcLateOvernight(t, tx)

	// When staff-001 อนุมัติ C
	w := tcApproveAdvisor(tcServer(tx, &fakeClock{t: now}), tokenStaffDept1, c, nil, now)

	// Then (1) ได้ 422 VALIDATION_FAILED ช่อง advisorName ไม่ใช่ APPROVAL_CUTOFF_PASSED
	tcExpectAdvisorError(t, "(1)", w, tcMsgAdvisorRequired)
	// Then (2) C ยังเป็น requested
	if s := tcRequestStatus(t, tx, c); s != "requested" {
		t.Errorf("(2) C ต้องยังเป็น requested แต่ได้ %q", s)
	}
}
