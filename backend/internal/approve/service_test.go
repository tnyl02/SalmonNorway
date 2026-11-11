package approve

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
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
	res, err := NewService(NewRepository(tx)).Approve(ctx, id)
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
	if _, err := svc.Approve(ctx, id); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Approve(ctx, id); err != ErrNotRequested {
		t.Errorf("อนุมัติซ้ำต้องได้ ErrNotRequested แต่ได้ %v", err)
	}

	overnight := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	if _, err := svc.Approve(ctx, overnight); err != ErrRequestNotFound {
		t.Errorf("คำขอค้างคืนต้องได้ ErrRequestNotFound (Q-APV-01) แต่ได้ %v", err)
	}
	if _, err := svc.Approve(ctx, 999999); err != ErrRequestNotFound {
		t.Errorf("คำขอที่ไม่มีอยู่ต้องได้ ErrRequestNotFound แต่ได้ %v", err)
	}

	student := auth.WithPrincipal(context.Background(), auth.Principal{UserID: "student-001", Role: auth.RoleStudent, DepartmentID: 1})
	if _, err := svc.Approve(student, requestIDOf(t, tx, "AP-TEST-0003", "student-002")); err != ErrForbiddenRole {
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
	res, err := NewService(NewRepository(tx)).Reject(ctx, id, reason)
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

	res, err := NewService(NewRepository(tx)).Approve(ctx, a)
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
	if _, err := NewService(NewRepository(tx)).Reject(staffCtx("staff-001", 1), b, "ขอยืมซ้ำ"); err != nil {
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

// AC-APV-16 (FR-APV-01 ค, DOM-LOAN-01) [ID รอทีมอนุมัติ Q-APV-13]
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
	if _, err := svc.Approve(staffCtx("staff-001", 1), b); err != ErrNotFirst {
		t.Errorf("B ต้องได้ ErrNotFirst เพราะ A ยังเป็นคำขอแรก แต่ได้ %v", err)
	}
	// เจ้าหน้าที่ปฏิเสธ A เองได้ (Q-APV-16)
	if _, err := svc.Reject(staffCtx("staff-001", 1), a, "เลยเวลาที่อนุมัติได้"); err != nil {
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
	if _, err := svc.Approve(staffCtx("staff-001", 1), id); err != nil {
		t.Errorf("ก่อน 16:00 น. ต้องอนุมัติได้ แต่ได้ %v", err)
	}
}
