package approve

import (
	"context"
	"testing"

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
