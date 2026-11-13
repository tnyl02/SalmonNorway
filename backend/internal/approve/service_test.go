package approve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/platform"
)

// staffCtx จำลองเจ้าหน้าที่ที่ยืนยันตัวตนแล้ว แบบเดียวกับที่ auth.Middleware ใส่ใน context (IF-IDP-01)
func staffCtx(userID string, dept int) context.Context {
	return auth.WithPrincipal(context.Background(),
		auth.Principal{UserID: userID, Role: auth.RoleStaff, DepartmentID: dept})
}

// requestIDOf คืนรหัสคำขอของอุปกรณ์และผู้ขอที่ระบุ จากข้อมูลใน approve_dev.sql
func requestIDOf(t *testing.T, tx pgx.Tx, assetCode, requesterID string) int {
	t.Helper()
	var id int
	err := tx.QueryRow(context.Background(), `
		SELECT r.id FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id
		WHERE i.asset_code = $1 AND r.requester_id = $2`, assetCode, requesterID).Scan(&id)
	if err != nil {
		t.Fatalf("หาคำขอของ %s / %s ไม่พบ: %v", assetCode, requesterID, err)
	}
	return id
}

// AC-APV-01 (FR-APV-01, CON-DATA-01)
// Given เจ้าหน้าที่ผ่านการยืนยันตัวตนแล้ว และมีคำขอยืมแบบปกติสถานะ requested
//
//	ที่เป็นคำขอแรกของอุปกรณ์รายชิ้นที่สถานะ "ว่าง" (AP-TEST-0002 ใน approve_dev.sql)
//
// When  เจ้าหน้าที่อนุมัติคำขอนั้น
// Then  สถานะของคำขอเป็น approved และสถานะของอุปกรณ์รายชิ้นยังเป็น "ว่าง"
//
//	(T-07 เสร็จเมื่อ: มีแถว approval_decision และ audit_log action approve)
func TestAC_APV_01(t *testing.T) {
	tx := devTx(t)
	ctx := staffCtx("staff-001", 1)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")

	// When
	res, err := NewService(NewRepository(tx)).Approve(ctx, id, time.Time{})
	if err != nil {
		t.Fatalf("อนุมัติต้องสำเร็จ แต่ได้ error: %v", err)
	}
	if res.Status != config.StatusApproved || res.DecidedAt.IsZero() {
		t.Errorf("ผลต้องเป็น approved พร้อมเวลาที่บันทึก แต่ได้ %+v", res)
	}

	// Then: สถานะของคำขอเป็น approved และสถานะอุปกรณ์ยังเป็น "ว่าง" (CON-DATA-01)
	var reqStatus, itemStatus string
	if err := tx.QueryRow(context.Background(), `
		SELECT r.status, i.status FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id
		WHERE r.id = $1`, id).Scan(&reqStatus, &itemStatus); err != nil {
		t.Fatal(err)
	}
	if reqStatus != "approved" {
		t.Errorf("สถานะของคำขอต้องเป็น approved แต่ได้ %q", reqStatus)
	}
	if itemStatus != "available" {
		t.Errorf("สถานะอุปกรณ์ต้องยังเป็น available (ว่าง) แต่ได้ %q", itemStatus)
	}

	// ประวัติ (FR-APV-01 ฌ, RQ-25)
	var result, actor string
	var reason *string
	if err := tx.QueryRow(context.Background(),
		`SELECT result, reason, actor_id FROM approval_decision WHERE request_id = $1`, id).Scan(&result, &reason, &actor); err != nil {
		t.Fatalf("ต้องมีแถว approval_decision: %v", err)
	}
	if result != "approved" || reason != nil || actor != "staff-001" {
		t.Errorf("approval_decision ไม่ถูกต้อง: result=%q reason=%v actor=%q", result, reason, actor)
	}
	var audits int
	if err := tx.QueryRow(context.Background(),
		`SELECT count(*) FROM audit_log WHERE request_id = $1 AND action = 'approve' AND actor_id = 'staff-001'`, id).Scan(&audits); err != nil {
		t.Fatal(err)
	}
	if audits != 1 {
		t.Errorf("ต้องมี audit_log action approve 1 แถว แต่ได้ %d", audits)
	}
}

// T-07: อนุมัติได้เฉพาะคำขอ requested อนุมัติซ้ำไม่ได้ และคำขอค้างคืนไม่อยู่ในขอบเขต (FR-APV-01 ค, Q-APV-01)
func TestApprove_OnlyRequestedNormal(t *testing.T) {
	tx := devTx(t)
	ctx := staffCtx("staff-001", 1)
	svc := NewService(NewRepository(tx))

	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	if _, err := svc.Approve(ctx, id, time.Time{}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Approve(ctx, id, time.Time{}); err != ErrAlreadyDecided {
		t.Errorf("อนุมัติซ้ำต้องได้ ErrAlreadyDecided แต่ได้ %v", err)
	}

	// v3 (T-44): คำขอค้างคืนอยู่ในขอบเขตแล้ว อนุมัติโดยไม่มีชื่ออาจารย์ได้ ErrAdvisorNameRequired (FR-APV-03 ช)
	overnight := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	if _, err := svc.Approve(ctx, overnight, time.Time{}); err != ErrAdvisorNameRequired {
		t.Errorf("คำขอค้างคืนที่ไม่มีชื่ออาจารย์ต้องได้ ErrAdvisorNameRequired แต่ได้ %v", err)
	}
	if _, err := svc.Approve(ctx, 999999, time.Time{}); err != ErrRequestNotFound {
		t.Errorf("คำขอที่ไม่มีอยู่ต้องได้ ErrRequestNotFound แต่ได้ %v", err)
	}

	student := auth.WithPrincipal(context.Background(), auth.Principal{UserID: "student-001", Role: auth.RoleStudent, DepartmentID: 1})
	if _, err := svc.Approve(student, requestIDOf(t, tx, "AP-TEST-0003", "student-002"), time.Time{}); err != ErrForbiddenRole {
		t.Errorf("นักศึกษาต้องได้ ErrForbiddenRole แต่ได้ %v", err)
	}
}

// AC-APV-02 (FR-APV-01)
// Given เจ้าหน้าที่ผ่านการยืนยันตัวตนแล้ว และมีคำขอยืมแบบปกติสถานะ requested (AP-TEST-0003 ใน approve_dev.sql)
// When  เจ้าหน้าที่ปฏิเสธคำขอนั้นพร้อมระบุเหตุผล
// Then  สถานะของคำขอเป็น rejected และสถานะของอุปกรณ์รายชิ้นไม่เปลี่ยน
//
//	(T-10 เสร็จเมื่อ: เหตุผลถูกบันทึกใน approval_decision และมี audit_log action reject)
func TestAC_APV_02(t *testing.T) {
	tx := devTx(t)
	ctx := staffCtx("staff-001", 1)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	const reason = "อุปกรณ์ต้องใช้สอบ"

	var itemBefore string
	if err := tx.QueryRow(context.Background(), `
		SELECT i.status FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id WHERE r.id = $1`, id).Scan(&itemBefore); err != nil {
		t.Fatal(err)
	}

	// When
	res, err := NewService(NewRepository(tx)).Reject(ctx, id, reason, time.Time{})
	if err != nil {
		t.Fatalf("ปฏิเสธต้องสำเร็จ แต่ได้ error: %v", err)
	}
	if res.Status != config.StatusRejected {
		t.Errorf("ผลต้องเป็น rejected แต่ได้ %+v", res)
	}

	// Then
	var reqStatus, itemAfter string
	if err := tx.QueryRow(context.Background(), `
		SELECT r.status, i.status FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id WHERE r.id = $1`, id).Scan(&reqStatus, &itemAfter); err != nil {
		t.Fatal(err)
	}
	if reqStatus != "rejected" {
		t.Errorf("สถานะของคำขอต้องเป็น rejected แต่ได้ %q", reqStatus)
	}
	if itemAfter != itemBefore {
		t.Errorf("สถานะอุปกรณ์ต้องไม่เปลี่ยน (%q) แต่ได้ %q (CON-DATA-01)", itemBefore, itemAfter)
	}
	var result, gotReason string
	if err := tx.QueryRow(context.Background(),
		`SELECT result, reason FROM approval_decision WHERE request_id = $1`, id).Scan(&result, &gotReason); err != nil {
		t.Fatalf("ต้องมีแถว approval_decision: %v", err)
	}
	if result != "rejected" || gotReason != reason {
		t.Errorf("approval_decision ต้องเป็น rejected เหตุผล %q แต่ได้ %q / %q", reason, result, gotReason)
	}
	var audits int
	if err := tx.QueryRow(context.Background(),
		`SELECT count(*) FROM audit_log WHERE request_id = $1 AND action = 'reject'`, id).Scan(&audits); err != nil {
		t.Fatal(err)
	}
	if audits != 1 {
		t.Errorf("ต้องมี audit_log action reject 1 แถว แต่ได้ %d", audits)
	}
}

// AC-APV-03 (FR-APV-01 จ, DOM-REQ-01)
// Given อุปกรณ์รายชิ้นเดียวกันมีคำขอ A (ส่งก่อน) และ B สถานะ requested (AP-TEST-0001 ใน approve_dev.sql)
// When  เจ้าหน้าที่อนุมัติคำขอ A
// Then  คำขอ A เป็น approved คำขอ B เป็น rejected ด้วยเหตุผล "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว"
//
//	ผู้ดำเนินการคือเจ้าหน้าที่คนนั้น (T-16 เสร็จเมื่อ: มี audit_log action reject และ Email ผลปฏิเสธของ B)
func TestAC_APV_03(t *testing.T) {
	tx := devTx(t)
	ctx := staffCtx("staff-001", 1)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003") // คำขอค้างคืนของชิ้นเดียวกัน (Q-APV-14)

	res, err := NewService(NewRepository(tx)).Approve(ctx, a, time.Time{})
	if err != nil {
		t.Fatalf("อนุมัติ A ต้องสำเร็จ: %v", err)
	}
	if len(res.RejectedBy) != 2 {
		t.Errorf("ต้องปฏิเสธตาม 2 คำขอ (B และค้างคืน C) แต่ได้ %v", res.RejectedBy)
	}

	status := func(id int) string {
		t.Helper()
		var s string
		if err := tx.QueryRow(context.Background(), `SELECT status FROM borrow_request WHERE id = $1`, id).Scan(&s); err != nil {
			t.Fatal(err)
		}
		return s
	}
	if status(a) != "approved" || status(b) != "rejected" || status(c) != "rejected" {
		t.Errorf("A ต้อง approved, B และ C ต้อง rejected แต่ได้ %s / %s / %s", status(a), status(b), status(c))
	}

	for _, id := range []int{b, c} {
		var result, reason, actor string
		if err := tx.QueryRow(context.Background(),
			`SELECT result, reason, actor_id FROM approval_decision WHERE request_id = $1`, id).Scan(&result, &reason, &actor); err != nil {
			t.Fatalf("คำขอ %d ต้องมี approval_decision: %v", id, err)
		}
		if result != "rejected" || reason != config.ApprovalSupersededReason || actor != "staff-001" {
			t.Errorf("คำขอ %d ได้ %q / %q / %q", id, result, reason, actor)
		}
		var audits, emails int
		if err := tx.QueryRow(context.Background(), `
			SELECT (SELECT count(*) FROM audit_log WHERE request_id = $1 AND action = 'reject'),
			       (SELECT count(*) FROM approval_email WHERE request_id = $1 AND kind = 'rejected')`, id).Scan(&audits, &emails); err != nil {
			t.Fatal(err)
		}
		if audits != 1 || emails != 1 {
			t.Errorf("คำขอ %d ต้องมี audit_log reject 1 แถว และ Email ผลปฏิเสธ 1 ฉบับ แต่ได้ %d / %d", id, audits, emails)
		}
	}
}

// AC-APV-07 (FR-APV-01 ข, DOM-REQ-01)
// Given อุปกรณ์รายชิ้นเดียวกันมีคำขอ A (ส่งก่อน) และ B สถานะ requested (AP-TEST-0001)
// When  เจ้าหน้าที่พยายามอนุมัติคำขอ B
// Then  ระบบไม่อนุมัติ แสดง "ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน" และทั้งสองคำขอคงเดิม
//
//	(T-17 เสร็จเมื่อ: ได้ 409 NOT_FIRST_IN_QUEUE และปฏิเสธ B ได้ตามปกติ)
func TestAC_APV_07(t *testing.T) {
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")

	// When: อนุมัติ B ผ่าน HTTP
	w := postJSON(newTestServer(tx), tokenStaffDept1, fmt.Sprintf("%srequests/%d/approve", PathPrefix, b), `{}`)

	// Then: 409 NOT_FIRST_IN_QUEUE พร้อมข้อความ
	if w.Code != http.StatusConflict {
		t.Fatalf("ต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var e ErrorDTO
	_ = json.Unmarshal(w.Body.Bytes(), &e)
	if e.Code != "NOT_FIRST_IN_QUEUE" || e.Message != "ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน" {
		t.Errorf("ได้ %+v", e)
	}
	// ทั้งสองคำขอคงเดิม ไม่มีผลการพิจารณา
	var requested, decisions int
	if err := tx.QueryRow(context.Background(), `
		SELECT count(*) FILTER (WHERE status = 'requested'),
		       (SELECT count(*) FROM approval_decision WHERE request_id = ANY($1))
		FROM borrow_request WHERE id = ANY($1)`, []int{a, b}).Scan(&requested, &decisions); err != nil {
		t.Fatal(err)
	}
	if requested != 2 || decisions != 0 {
		t.Errorf("A และ B ต้องยังเป็น requested และไม่มี approval_decision แต่ได้ %d / %d", requested, decisions)
	}

	// ปฏิเสธ B ได้ตามปกติ (ลำดับบังคับเฉพาะการอนุมัติ FR-APV-01 ข)
	if _, err := NewService(NewRepository(tx)).Reject(staffCtx("staff-001", 1), b, "ขอยืมซ้ำ", time.Time{}); err != nil {
		t.Errorf("ปฏิเสธ B ต้องทำได้ แต่ได้ %v", err)
	}
}

// AC-APV-08 (FR-APV-01 ค, CON-DATA-01)
// Given คำขอแรกสถานะ requested ของอุปกรณ์รายชิ้นที่สถานะไม่ใช่ "ว่าง"
//
//	(AP-TEST-0004 ถูกยืม ใน approve_dev.sql และ AP-TEST-0002 ที่ test เปลี่ยนเป็นไม่พร้อมใช้งาน)
//
// When  เจ้าหน้าที่พยายามอนุมัติ
// Then  ระบบไม่อนุมัติ ตอบ ITEM_NOT_AVAILABLE แสดง "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้" และสถานะคงเดิม
func TestAC_APV_08(t *testing.T) {
	tx := devTx(t)
	if _, err := tx.Exec(context.Background(),
		`UPDATE equipment_item SET status = 'unavailable' WHERE asset_code = 'AP-TEST-0002'`); err != nil {
		t.Fatal(err)
	}
	srv := newTestServer(tx)

	for _, c := range []struct{ asset, requester, itemStatus string }{
		{"AP-TEST-0004", "student-001", "borrowed"},
		{"AP-TEST-0002", "student-001", "unavailable"},
	} {
		id := requestIDOf(t, tx, c.asset, c.requester)
		w := postJSON(srv, tokenStaffDept1, fmt.Sprintf("%srequests/%d/approve", PathPrefix, id), `{}`)
		if w.Code != http.StatusConflict {
			t.Errorf("%s (%s) ต้องได้ 409 แต่ได้ %d", c.asset, c.itemStatus, w.Code)
			continue
		}
		var e ErrorDTO
		_ = json.Unmarshal(w.Body.Bytes(), &e)
		if e.Code != "ITEM_NOT_AVAILABLE" || e.Message != "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้" {
			t.Errorf("%s ได้ %+v", c.asset, e)
		}
		var reqStatus, itemStatus string
		if err := tx.QueryRow(context.Background(), `
			SELECT r.status, i.status FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id WHERE r.id = $1`, id).Scan(&reqStatus, &itemStatus); err != nil {
			t.Fatal(err)
		}
		if reqStatus != "requested" || itemStatus != c.itemStatus {
			t.Errorf("%s สถานะต้องคงเดิม (requested / %s) แต่ได้ %s / %s", c.asset, c.itemStatus, reqStatus, itemStatus)
		}
	}
}

// AC-APV-16 (FR-APV-01 ค, DOM-LOAN-01)
// Given คำขอแรกสถานะ requested ของอุปกรณ์ที่ "ว่าง" และเวลาปัจจุบันเลย 16:00 น. ของวันที่เริ่มยืมแล้ว
// When  เจ้าหน้าที่พยายามอนุมัติ
// Then  ระบบไม่อนุมัติ แสดงข้อความกรณีเลย 16:00 น. และสถานะคงเดิม
//
//	(T-21 เสร็จเมื่อ: คำขอนั้นยังเป็นคำขอแรก คำขอถัดไปยังได้ NOT_FIRST_IN_QUEUE และปฏิเสธเองได้)
func TestAC_APV_16(t *testing.T) {
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	var startAt time.Time
	if err := tx.QueryRow(context.Background(), `SELECT start_at FROM borrow_request WHERE id = $1`, a).Scan(&startAt); err != nil {
		t.Fatal(err)
	}

	// นาฬิกาจำลอง: เลยเวลาตัดของวันที่เริ่มยืม 1 นาที / ก่อนเวลาตัด 1 นาทียังอนุมัติได้ (ตรวจด้วย approvalCutoff ตรง ๆ)
	cutoff := approvalCutoff(startAt)
	if h := localTime(cutoff).Hour(); h != config.LoanNormalReturnHour {
		t.Fatalf("เวลาตัดต้องเป็น %d:00 น. เวลาไทย แต่ได้ %v", config.LoanNormalReturnHour, localTime(cutoff))
	}
	clock := &fakeClock{t: cutoff.Add(time.Minute)}
	svc := NewService(NewRepository(tx))
	svc.now = clock.Now
	srv := auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes())

	// When: อนุมัติ A
	w := postJSON(srv, tokenStaffDept1, fmt.Sprintf("%srequests/%d/approve", PathPrefix, a), `{}`)

	// Then: 409 APPROVAL_CUTOFF_PASSED พร้อมข้อความ และสถานะคงเดิม
	if w.Code != http.StatusConflict {
		t.Fatalf("ต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var e ErrorDTO
	_ = json.Unmarshal(w.Body.Bytes(), &e)
	if e.Code != "APPROVAL_CUTOFF_PASSED" || e.Message != "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้" {
		t.Errorf("ได้ %+v", e)
	}
	var status string
	if err := tx.QueryRow(context.Background(), `SELECT status FROM borrow_request WHERE id = $1`, a).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "requested" {
		t.Errorf("สถานะของ A ต้องคงเป็น requested แต่ได้ %q", status)
	}

	// A ยังเป็นคำขอแรก: อนุมัติ B ได้ NOT_FIRST_IN_QUEUE (Q-APV-16)
	if _, err := svc.Approve(staffCtx("staff-001", 1), b, time.Time{}); err != ErrNotFirst {
		t.Errorf("B ต้องได้ ErrNotFirst เพราะ A ยังเป็นคำขอแรก แต่ได้ %v", err)
	}
	// เจ้าหน้าที่ปฏิเสธ A เองได้ (Q-APV-16)
	if _, err := svc.Reject(staffCtx("staff-001", 1), a, "เลยเวลาที่อนุมัติได้", time.Time{}); err != nil {
		t.Errorf("ปฏิเสธ A ต้องทำได้ แต่ได้ %v", err)
	}

	// เวลาตัดเป็นวันเดียวกับวันที่เริ่มยืม (ไม่ใช่วันที่ส่งหรือวันนี้) กันกรณีคิดผิดวัน
	if localTime(cutoff).Format("2006-01-02") != localTime(startAt).Format("2006-01-02") {
		t.Errorf("เวลาตัดต้องอยู่วันเดียวกับวันที่เริ่มยืม: start %v cutoff %v", localTime(startAt), localTime(cutoff))
	}
}

// T-21: ก่อนเวลาตัด 1 นาทีอนุมัติได้ (ขอบเขตของ "เลย 16:00 น.")
func TestApprove_BeforeCutoffAllowed(t *testing.T) {
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	var startAt time.Time
	if err := tx.QueryRow(context.Background(), `SELECT start_at FROM borrow_request WHERE id = $1`, id).Scan(&startAt); err != nil {
		t.Fatal(err)
	}
	clock := &fakeClock{t: approvalCutoff(startAt).Add(-time.Minute)}
	svc := NewService(NewRepository(tx))
	svc.now = clock.Now
	if _, err := svc.Approve(staffCtx("staff-001", 1), id, time.Time{}); err != nil {
		t.Errorf("ก่อน 16:00 น. ต้องอนุมัติได้ แต่ได้ %v", err)
	}
}

// AC-APV-09 (FR-APV-01 ซ)
// Given คำขอสถานะ cancelled หรือ approved / rejected
// When  เจ้าหน้าที่พยายามอนุมัติหรือปฏิเสธ
// Then  ระบบไม่ดำเนินการ แสดง "คำขอนี้ถูกผู้ขอยกเลิกแล้ว" หรือ "คำขอนี้ได้รับการพิจารณาแล้ว" ตามกรณี และสถานะคงเดิม
//
//	(T-23 เสร็จเมื่อ: พิจารณาหลัง loadedAt ได้ REQUEST_DECIDED_CONCURRENTLY ทั้งอนุมัติและปฏิเสธ ผลเดิมคงไว้)
func TestAC_APV_09(t *testing.T) {
	tx := devTx(t)
	srv := newTestServer(tx)
	ctx := context.Background()

	call := func(action string, id int, loadedAt time.Time) ErrorDTO {
		t.Helper()
		body := fmt.Sprintf(`{"loadedAt":%q,"reason":"ทดสอบ"}`, loadedAt.UTC().Format(time.RFC3339Nano))
		w := postJSON(srv, tokenStaffDept1, fmt.Sprintf("%srequests/%d/%s", PathPrefix, id, action), body)
		if w.Code != http.StatusConflict {
			t.Errorf("%s คำขอ %d ต้องได้ 409 แต่ได้ %d body=%s", action, id, w.Code, w.Body.String())
		}
		var e ErrorDTO
		_ = json.Unmarshal(w.Body.Bytes(), &e)
		return e
	}
	statusOf := func(id int) string {
		t.Helper()
		var s string
		if err := tx.QueryRow(ctx, `SELECT status FROM borrow_request WHERE id = $1`, id).Scan(&s); err != nil {
			t.Fatal(err)
		}
		return s
	}

	// (1) คำขอที่ผู้ขอยกเลิกแล้ว (AP-TEST-0005 ใน approve_dev.sql)
	cancelled := requestIDOf(t, tx, "AP-TEST-0005", "student-002")
	for _, action := range []string{"approve", "reject"} {
		if e := call(action, cancelled, time.Now()); e.Code != "REQUEST_CANCELLED" || e.Message != "คำขอนี้ถูกผู้ขอยกเลิกแล้ว" {
			t.Errorf("%s คำขอที่ยกเลิกแล้วได้ %+v", action, e)
		}
	}
	if s := statusOf(cancelled); s != "cancelled" {
		t.Errorf("สถานะต้องคงเป็น cancelled แต่ได้ %q", s)
	}

	// เตรียมคำขอที่พิจารณาไปแล้ว: อนุมัติ AP-TEST-0002
	decided := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	if _, err := NewService(NewRepository(tx)).Approve(staffCtx("staff-001", 1), decided, time.Time{}); err != nil {
		t.Fatal(err)
	}
	var decidedAt time.Time
	if err := tx.QueryRow(ctx, `SELECT decided_at FROM approval_decision WHERE request_id = $1`, decided).Scan(&decidedAt); err != nil {
		t.Fatal(err)
	}

	// (2) พิจารณาไปแล้วก่อนหน้าจอโหลด (loadedAt หลังเวลาที่บันทึกผล)
	for _, action := range []string{"approve", "reject"} {
		if e := call(action, decided, decidedAt.Add(time.Minute)); e.Code != "REQUEST_ALREADY_DECIDED" || e.Message != "คำขอนี้ได้รับการพิจารณาแล้ว" {
			t.Errorf("%s คำขอที่พิจารณาแล้วได้ %+v", action, e)
		}
	}

	// (3) ถูกพิจารณาหลังหน้าจอโหลด (กดพร้อมกัน loadedAt ก่อนเวลาที่บันทึกผล)
	for _, action := range []string{"approve", "reject"} {
		if e := call(action, decided, decidedAt.Add(-time.Minute)); e.Code != "REQUEST_DECIDED_CONCURRENTLY" ||
			e.Message != "คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น กรุณารีเฟรชหน้า" {
			t.Errorf("%s คำขอที่เพิ่งถูกพิจารณาได้ %+v", action, e)
		}
	}

	// ผลเดิมคงไว้: ยัง approved และมี approval_decision แถวเดียว
	var decisions int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM approval_decision WHERE request_id = $1`, decided).Scan(&decisions); err != nil {
		t.Fatal(err)
	}
	if s := statusOf(decided); s != "approved" || decisions != 1 {
		t.Errorf("ผลเดิมต้องคงไว้ (approved, 1 แถว) แต่ได้ %q, %d แถว", s, decisions)
	}
}

// assertHistory ตรวจประวัติของคำขอ 1 ใบ (FR-APV-01 ฌ, RQ-25): แถว approval_decision ที่มีรหัสคำขอ ผล เหตุผล
// ผู้ดำเนินการ และวันเวลา และ audit_log action ตรงกรณี 1 แถว ใช้ใน TestAC_APV_15 และกรณีปฏิเสธอัตโนมัติของ T-29
func assertHistory(t *testing.T, tx pgx.Tx, requestID int, result, reason, actorID, action string) {
	t.Helper()
	var gotResult, gotActor string
	var gotReason *string
	var decidedAt time.Time
	if err := tx.QueryRow(context.Background(), `
		SELECT result, reason, actor_id, decided_at FROM approval_decision WHERE request_id = $1`, requestID).
		Scan(&gotResult, &gotReason, &gotActor, &decidedAt); err != nil {
		t.Errorf("คำขอ %d ต้องมีแถว approval_decision: %v", requestID, err)
		return
	}
	r := ""
	if gotReason != nil {
		r = *gotReason
	}
	if gotResult != result || r != reason || gotActor != actorID || decidedAt.IsZero() {
		t.Errorf("คำขอ %d ประวัติไม่ตรง: ได้ (%q, %q, %q, %v) ต้องการ (%q, %q, %q, มีวันเวลา)",
			requestID, gotResult, r, gotActor, decidedAt, result, reason, actorID)
	}
	var audits int
	if err := tx.QueryRow(context.Background(),
		`SELECT count(*) FROM audit_log WHERE request_id = $1 AND action = $2 AND actor_id = $3`, requestID, action, actorID).Scan(&audits); err != nil {
		t.Fatal(err)
	}
	if audits != 1 {
		t.Errorf("คำขอ %d ต้องมี audit_log action %q 1 แถว แต่ได้ %d", requestID, action, audits)
	}
}

// AC-APV-15 (FR-APV-01 ฌ, RQ-25)
// Given การอนุมัติ ปฏิเสธ หรือปฏิเสธอัตโนมัติใด ๆ
// When  ระบบบันทึกผลสำเร็จ
// Then  มีประวัติใหม่ที่มีรหัสคำขอ ผล เหตุผล (ถ้ามี) ผู้ดำเนินการ และวันเวลา
//
//	T-28 ตรวจอนุมัติ ปฏิเสธ และปฏิเสธตาม / กรณีปฏิเสธอัตโนมัติเพิ่มใน T-29
//	ชื่อเจ้าหน้าที่ (actor_name) ยังว่างได้เพราะ Principal ยังไม่มีชื่อ (T-38)
func TestAC_APV_15(t *testing.T) {
	tx := devTx(t)
	svc := NewService(NewRepository(tx))
	ctx := staffCtx("staff-001", 1)

	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	rejected := requestIDOf(t, tx, "AP-TEST-0003", "student-002")

	// อนุมัติ A (B ถูกปฏิเสธตาม) และปฏิเสธคำขอของ AP-TEST-0003 พร้อมเหตุผล
	if _, err := svc.Approve(ctx, a, time.Time{}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Reject(ctx, rejected, "อุปกรณ์ต้องใช้สอบ", time.Time{}); err != nil {
		t.Fatal(err)
	}

	assertHistory(t, tx, a, config.ApprovalKindApproved, "", "staff-001", config.AuditActionApprove)
	assertHistory(t, tx, rejected, config.ApprovalKindRejected, "อุปกรณ์ต้องใช้สอบ", "staff-001", config.AuditActionReject)
	assertHistory(t, tx, b, config.ApprovalKindRejected, config.ApprovalSupersededReason, "staff-001", config.AuditActionReject)

	// T-29: ปฏิเสธอัตโนมัติ คำขอของ AP-TEST-0004 ส่งเมื่อวาน แล้วงานเที่ยงคืนทำงาน
	stale := requestIDOf(t, tx, "AP-TEST-0004", "student-001")
	midnight := autoRejectCutoff(time.Now())
	if _, err := tx.Exec(context.Background(), `UPDATE borrow_request SET submitted_at = $2 WHERE id = $1`,
		stale, midnight.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	job := &AutoRejecter{db: tx, now: func() time.Time { return midnight.Add(time.Minute) }}
	if _, err := job.RejectStale(context.Background()); err != nil {
		t.Fatal(err)
	}
	assertHistory(t, tx, stale, config.ApprovalKindAutoRejected, config.ApprovalAutoRejectReason,
		config.ApprovalSystemActorID, config.AuditActionAutoReject)
}

// AC-APV-17 (FR-APV-01 จ ฉ, Q-APV-14)
// Given อุปกรณ์รายชิ้นเดียวกันมีคำขอยืมปกติ A (ส่งก่อน) และคำขอค้างคืน C สถานะ requested (AP-TEST-0001)
// When  เจ้าหน้าที่อนุมัติคำขอ A / หรือเมื่อถึง 00:00 น. หลังวันที่ส่งของ C
// Then  กรณีอนุมัติ A: C เป็น rejected ด้วยเหตุผล "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว" และผู้ขอได้รับ Email ผลปฏิเสธ
//
//	/ กรณีเที่ยงคืน: C ยังเป็น requested
func TestAC_APV_17(t *testing.T) {
	t.Run("อนุมัติ A แล้ว C ค้างคืนถูกปฏิเสธตาม", func(t *testing.T) {
		tx := devTx(t)
		a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
		c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")

		res, err := NewService(NewRepository(tx)).Approve(staffCtx("staff-001", 1), a, time.Time{})
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, id := range res.RejectedBy {
			found = found || id == c
		}
		if !found {
			t.Errorf("C ต้องอยู่ในรายการคำขอที่ถูกปฏิเสธตาม แต่ได้ %v", res.RejectedBy)
		}

		var status, reason, kind string
		if err := tx.QueryRow(context.Background(), `
			SELECT r.status, d.reason, e.kind FROM borrow_request r
			JOIN approval_decision d ON d.request_id = r.id
			JOIN approval_email    e ON e.request_id = r.id
			WHERE r.id = $1`, c).Scan(&status, &reason, &kind); err != nil {
			t.Fatalf("C ต้องมีผลการพิจารณาและ Email: %v", err)
		}
		if status != "rejected" || reason != config.ApprovalSupersededReason || kind != config.ApprovalKindRejected {
			t.Errorf("C ได้ %q / %q / Email %q", status, reason, kind)
		}
	})

	t.Run("เที่ยงคืนหลังวันที่ส่งของ C แล้ว C ยังเป็น requested", func(t *testing.T) {
		tx := devTx(t)
		c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
		midnight := autoRejectCutoff(time.Now()).AddDate(0, 0, 1)
		if _, err := tx.Exec(context.Background(), `UPDATE borrow_request SET submitted_at = $2 WHERE id = $1`,
			c, midnight.Add(-2*time.Hour)); err != nil {
			t.Fatal(err)
		}

		job := &AutoRejecter{db: tx, now: func() time.Time { return midnight.Add(time.Minute) }}
		if _, err := job.RejectStale(context.Background()); err != nil {
			t.Fatal(err)
		}

		var status string
		var decisions int
		if err := tx.QueryRow(context.Background(), `
			SELECT r.status, (SELECT count(*) FROM approval_decision d WHERE d.request_id = r.id)
			FROM borrow_request r WHERE r.id = $1`, c).Scan(&status, &decisions); err != nil {
			t.Fatal(err)
		}
		if status != "requested" || decisions != 0 {
			t.Errorf("C ค้างคืนต้องยังเป็น requested ไม่มีผลการพิจารณา แต่ได้ %q / %d", status, decisions)
		}
	})
}

// concurrencyDB สร้างฐานข้อมูลแยกใน container เดียวกันสำหรับ test ที่ต้องใช้หลาย connection พร้อมกัน (T-35)
// ใช้ devTx ไม่ได้เพราะทุกอย่างใน transaction เดียวไม่มีการแย่งกันจริง และข้อมูลที่ commit แล้วลบไม่ได้
// (approval_decision และ audit_log แก้หรือลบไม่ได้) จึงแยกฐานข้อมูลแล้วลบทั้งฐานตอนจบ
func concurrencyDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	base := testPool(t)
	name := fmt.Sprintf("approve_conc_%d", time.Now().UnixNano())
	if _, err := base.Exec(ctx, "CREATE DATABASE "+name); err != nil {
		t.Fatal(err)
	}
	baseURL, err := dbContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	url := strings.Replace(baseURL, "/salmon?", "/"+name+"?", 1)
	if err := platform.Migrate(url, migrationsDir); err != nil {
		t.Fatal(err)
	}
	if err := platform.Seed(ctx, url, seedFile); err != nil {
		t.Fatal(err)
	}
	if err := platform.Seed(ctx, url, approveDevSQL); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		pool.Close()
		_, _ = base.Exec(ctx, "DROP DATABASE IF EXISTS "+name+" WITH (FORCE)")
	})
	return pool
}

// AC-APV-10 (FR-APV-01 ช)
// Given เจ้าหน้าที่ 2 คนเปิดคำขอเดียวกันสถานะ requested
// When  ทั้งสองกดอนุมัติหรือปฏิเสธพร้อมกัน
// Then  ระบบบันทึกเฉพาะผลของคนที่สำเร็จก่อน อีกคนเห็น "คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น กรุณารีเฟรชหน้า"
//
//	(T-35 เสร็จเมื่อ: ผ่านซ้ำ -count=20 มี approval_decision 1 แถวและ approval_email 1 แถว)
func TestAC_APV_10(t *testing.T) {
	pool := concurrencyDB(t)
	svc := NewService(NewRepository(pool))
	// หน้าจอของทั้งสองคนโหลดก่อนกด เผื่อเวลาของ container กับเครื่องต่างกันเล็กน้อย จึงย้อนไป 1 ชั่วโมง
	loadedAt := time.Now().Add(-time.Hour)

	type action func(ctx context.Context, id int) error
	approve := func(ctx context.Context, id int) error { _, err := svc.Approve(ctx, id, loadedAt); return err }
	reject := func(ctx context.Context, id int) error {
		_, err := svc.Reject(ctx, id, "กดพร้อมกัน", loadedAt)
		return err
	}

	cases := []struct {
		name, asset, requester string
		dept                   int
		first, second          action
	}{
		{"อนุมัติ+อนุมัติ", "AP-TEST-0002", "student-001", 1, approve, approve},
		{"อนุมัติ+ปฏิเสธ", "AP-TEST-0003", "student-002", 1, approve, reject},
		{"ปฏิเสธ+ปฏิเสธ", "AP-TEST-0101", "student-101", 2, reject, reject},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			var id int
			if err := pool.QueryRow(ctx, `
				SELECT r.id FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id
				WHERE i.asset_code = $1 AND r.requester_id = $2`, c.asset, c.requester).Scan(&id); err != nil {
				t.Fatal(err)
			}

			// When: 2 goroutine เริ่มพร้อมกันหลังสัญญาณ start
			start := make(chan struct{})
			errs := make([]error, 2)
			var wg sync.WaitGroup
			for i, act := range []action{c.first, c.second} {
				wg.Add(1)
				go func(i int, act action, staffID string) {
					defer wg.Done()
					<-start
					errs[i] = act(staffCtx(staffID, c.dept), id)
				}(i, act, fmt.Sprintf("staff-%d0%d", c.dept, i+1))
			}
			close(start)
			wg.Wait()

			// Then: สำเร็จ 1 อีกคนได้ ErrDecidedConcurrently
			ok, concurrent := 0, 0
			for _, err := range errs {
				switch {
				case err == nil:
					ok++
				case errors.Is(err, ErrDecidedConcurrently):
					concurrent++
				default:
					t.Errorf("ได้ error ที่ไม่คาด: %v", err)
				}
			}
			if ok != 1 || concurrent != 1 {
				t.Errorf("ต้องสำเร็จ 1 และได้ ErrDecidedConcurrently 1 แต่ได้ %d / %d (%v)", ok, concurrent, errs)
			}
			var decisions, emails int
			if err := pool.QueryRow(ctx, `
				SELECT (SELECT count(*) FROM approval_decision WHERE request_id = $1),
				       (SELECT count(*) FROM approval_email WHERE request_id = $1)`, id).Scan(&decisions, &emails); err != nil {
				t.Fatal(err)
			}
			if decisions != 1 || emails != 1 {
				t.Errorf("ต้องมี approval_decision 1 แถวและ Email 1 ฉบับ แต่ได้ %d / %d", decisions, emails)
			}
		})
	}
}

// AC-APV-22 (FR-APV-03 ข, CON-DATA-01)
// Given คำขอยืมแบบค้างคืนสถานะ requested ของอุปกรณ์ภาควิชาเดียวกับเจ้าหน้าที่ (C ของ AP-TEST-0001 ใน approve_dev.sql)
// When  เจ้าหน้าที่อนุมัติพร้อมระบุชื่ออาจารย์ที่เซ็น
// Then  คำขอเป็น approved สถานะอุปกรณ์ยังเป็น "ว่าง" และประวัติมีชื่ออาจารย์
//
//	(T-44 เสร็จเมื่อ: approval_decision.advisor_name ตรง, approval_step 3 แถวตามลำดับ, ไม่ติดคิวของ A และ B)
func TestAC_APV_22(t *testing.T) {
	tx := devTx(t)
	ctx := staffCtx("staff-001", 1)
	bg := context.Background()
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	const advisor = "อ.สมชาย ใจดี"

	// When: C ส่งหลัง A และ B แต่ไม่มีคิว (FR-APV-03 ค) และไม่มีเวลาตัด (FR-APV-03 ง) ใช้เวลาจริงได้
	res, err := NewService(NewRepository(tx)).ApproveWithAdvisor(ctx, c, time.Now(), advisor)
	if err != nil {
		t.Fatalf("อนุมัติคำขอค้างคืนต้องสำเร็จ แต่ได้ %v", err)
	}

	// Then: คำขอเป็น approved
	if res.Status != config.StatusApproved {
		t.Errorf("ผลต้องเป็น approved แต่ได้ %q", res.Status)
	}
	var status, itemStatus string
	if err := tx.QueryRow(bg, `
		SELECT r.status, i.status FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id WHERE r.id = $1`, c).
		Scan(&status, &itemStatus); err != nil {
		t.Fatal(err)
	}
	if status != config.StatusApproved {
		t.Errorf("สถานะคำขอต้องเป็น approved แต่ได้ %q", status)
	}
	// Then: สถานะอุปกรณ์ยังเป็น "ว่าง" (CON-DATA-01, data-model.md แถว UC-14)
	if itemStatus != config.ItemStatusAvailable {
		t.Errorf("สถานะอุปกรณ์ต้องยังเป็น available แต่ได้ %q", itemStatus)
	}
	// Then: ประวัติมีชื่ออาจารย์ (FR-APV-03 ข ช)
	var got *string
	if err := tx.QueryRow(bg, `SELECT advisor_name FROM approval_decision WHERE request_id = $1`, c).Scan(&got); err != nil {
		t.Fatal(err)
	}
	if got == nil || *got != advisor {
		t.Errorf("approval_decision.advisor_name ต้องเป็น %q แต่ได้ %v", advisor, got)
	}
	// FR-APV-03 (ข): approval_step 3 แถวตามลำดับ ผู้ดำเนินการคือเจ้าหน้าที่
	rows, err := tx.Query(bg, `SELECT status, actor_id FROM approval_step WHERE request_id = $1 ORDER BY id`, c)
	if err != nil {
		t.Fatal(err)
	}
	var steps []string
	for rows.Next() {
		var s, actor string
		if err := rows.Scan(&s, &actor); err != nil {
			t.Fatal(err)
		}
		if actor != "staff-001" {
			t.Errorf("approval_step.actor_id ต้องเป็น staff-001 แต่ได้ %q", actor)
		}
		steps = append(steps, s)
	}
	want := []string{config.StatusApprovedByOfficer, config.StatusPendingAdvisorApproval, config.StatusApproved}
	if strings.Join(steps, ",") != strings.Join(want, ",") {
		t.Errorf("approval_step ต้องเป็น %v แต่ได้ %v", want, steps)
	}
	// ประวัติ approval_decision และ audit_log action approve แบบเดียวกับยืมปกติ (FR-APV-01 ฌ, FR-APV-03 ช)
	assertHistory(t, tx, c, config.ApprovalKindApproved, "", "staff-001", config.AuditActionApprove)
}

// AC-APV-22 ทางผิด (FR-APV-03 ข ช, Q-APV-42, Q-APV-44)
//   - ชื่ออาจารย์ว่าง มีแต่ช่องว่าง หรือยาวเกิน 100 code point -> error และสถานะคงเดิม ไม่มี approval_step
//   - ชื่อยาว 100 code point พอดีผ่าน
//   - อุปกรณ์ไม่ว่าง -> ErrItemNotAvailable
//   - ลำดับ: ภาควิชาก่อนชื่ออาจารย์ และชื่ออาจารย์ก่อนสถานะคำขอ (ทีมตอบ 2569-11-13 ตอน T-44)
func TestApproveOvernight_Rules(t *testing.T) {
	ctx := staffCtx("staff-001", 1)
	bg := context.Background()
	long := strings.Repeat("ก", config.ApprovalAdvisorNameMaxLength+1)
	exact := strings.Repeat("ก", config.ApprovalAdvisorNameMaxLength)

	for _, tc := range []struct {
		name, advisor string
		want          error
	}{
		{"ชื่อว่าง", "", ErrAdvisorNameRequired},
		{"มีแต่ช่องว่าง", "   ", ErrAdvisorNameRequired},
		{"ยาว 101 code point", long, ErrAdvisorNameTooLong},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tx := devTx(t)
			c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
			if _, err := NewService(NewRepository(tx)).ApproveWithAdvisor(ctx, c, time.Now(), tc.advisor); err != tc.want {
				t.Fatalf("ต้องได้ %v แต่ได้ %v", tc.want, err)
			}
			if s := tcRequestStatus(t, tx, c); s != config.StatusRequested {
				t.Errorf("สถานะต้องคงเดิม requested แต่ได้ %q", s)
			}
			var n int
			if err := tx.QueryRow(bg, `SELECT count(*) FROM approval_step WHERE request_id = $1`, c).Scan(&n); err != nil {
				t.Fatal(err)
			}
			if n != 0 {
				t.Errorf("ต้องไม่มี approval_step แต่มี %d แถว", n)
			}
		})
	}

	t.Run("ยาว 100 code point พอดีผ่าน", func(t *testing.T) {
		tx := devTx(t)
		c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
		if _, err := NewService(NewRepository(tx)).ApproveWithAdvisor(ctx, c, time.Now(), exact); err != nil {
			t.Fatalf("ชื่อ 100 code point ต้องอนุมัติได้ แต่ได้ %v", err)
		}
	})

	t.Run("อุปกรณ์ไม่ว่าง", func(t *testing.T) {
		tx := devTx(t)
		// AP-TEST-0004 ถูกยืมอยู่ใน approve_dev.sql เพิ่มคำขอค้างคืนของชิ้นนั้น
		var id int
		if err := tx.QueryRow(bg, `
			INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose,
			                            start_at, end_at, status, attachment_path)
			SELECT 'student-003', 'ค้างคืนอุปกรณ์ไม่ว่าง', 'student-003', id, 'overnight', 'activity',
			       now(), now() + interval '1 day', 'requested', 'dev/x.pdf'
			FROM equipment_item WHERE asset_code = 'AP-TEST-0004' RETURNING id`).Scan(&id); err != nil {
			t.Fatal(err)
		}
		if _, err := NewService(NewRepository(tx)).ApproveWithAdvisor(ctx, id, time.Now(), "อ.สมชาย"); err != ErrItemNotAvailable {
			t.Fatalf("ต้องได้ ErrItemNotAvailable (Q-APV-42) แต่ได้ %v", err)
		}
		if s := tcRequestStatus(t, tx, id); s != config.StatusRequested {
			t.Errorf("สถานะต้องคงเดิม requested แต่ได้ %q", s)
		}
	})

	t.Run("ภาควิชาอื่นได้ ErrForbiddenDepartment ก่อนชื่ออาจารย์", func(t *testing.T) {
		tx := devTx(t)
		c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
		if _, err := NewService(NewRepository(tx)).ApproveWithAdvisor(staffCtx("staff-201", 2), c, time.Now(), ""); err != ErrForbiddenDepartment {
			t.Fatalf("ต้องได้ ErrForbiddenDepartment แต่ได้ %v", err)
		}
	})

	t.Run("คำขอยกเลิกแล้วและไม่มีชื่ออาจารย์ได้ ErrAdvisorNameRequired", func(t *testing.T) {
		tx := devTx(t)
		c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
		if _, err := tx.Exec(bg, `UPDATE borrow_request SET status = 'cancelled', cancel_reason = 'ทดสอบ', cancelled_at = now() WHERE id = $1`, c); err != nil {
			t.Fatal(err)
		}
		svc := NewService(NewRepository(tx))
		if _, err := svc.ApproveWithAdvisor(ctx, c, time.Now(), ""); err != ErrAdvisorNameRequired {
			t.Fatalf("ชื่ออาจารย์ต้องตรวจก่อนสถานะ ได้ ErrAdvisorNameRequired แต่ได้ %v", err)
		}
		if _, err := svc.ApproveWithAdvisor(ctx, c, time.Now(), "อ.สมชาย"); err != ErrRequestCancelled {
			t.Fatalf("มีชื่ออาจารย์แล้วต้องได้ ErrRequestCancelled แต่ได้ %v", err)
		}
	})

	t.Run("ยืมปกติส่งชื่ออาจารย์มาก็ไม่บันทึก", func(t *testing.T) {
		tx := devTx(t)
		a := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
		svc := NewService(NewRepository(tx))
		svc.now = func() time.Time { return approvalCutoff(tcStartAt(t, tx, a)).Add(-time.Hour) }
		if _, err := svc.ApproveWithAdvisor(ctx, a, time.Now(), "อ.สมชาย"); err != nil {
			t.Fatal(err)
		}
		var got *string
		if err := tx.QueryRow(bg, `SELECT advisor_name FROM approval_decision WHERE request_id = $1`, a).Scan(&got); err != nil {
			t.Fatal(err)
		}
		if got != nil {
			t.Errorf("ยืมปกติต้องไม่มี advisor_name แต่ได้ %q", *got)
		}
	})
}
