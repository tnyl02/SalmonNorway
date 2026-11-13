package platform

// test ของ migration 0008 และ 0011 ฟีเจอร์ 004 (T-01, T-40 ใน specs/004-approve-request/tasks.md)
// ใช้ testDB และ execInTx จาก migrate_test.go ทุกคำสั่งรันใน transaction ที่ยกเลิกทิ้ง ข้อมูลจึงไม่ค้างข้าม test

import (
	"context"
	"fmt"
	"os"
	"testing"
)

// คำขอยืม 1 ใบสำหรับให้ผลการพิจารณาอ้างถึง (อุปกรณ์ id 1 ของข้อมูลตั้งต้น ภาควิชา 1)
func approvalTestRequest(id int) string {
	return fmt.Sprintf(`INSERT INTO borrow_request
		(id, requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
		VALUES (%d, 'student-001', 'ผู้ขอทดสอบ', '660710758', 1, 'normal', 'study', now(), now(), 'requested')`, id)
}

func approvalTestDecision(requestID int, result, reason string) string {
	r := "NULL"
	if reason != "" {
		r = "'" + reason + "'"
	}
	return fmt.Sprintf(`INSERT INTO approval_decision (request_id, result, reason, actor_id, actor_department_id)
		VALUES (%d, '%s', %s, 'staff-001', 1)`, requestID, result, r)
}

func approvalTestEmail(requestID int, kind, status string) string {
	return fmt.Sprintf(`INSERT INTO approval_email (request_id, kind, status) VALUES (%d, '%s', '%s')`, requestID, kind, status)
}

// T-01 (004) เสร็จเมื่อ (ก)(ข): approval_decision.request_id UNIQUE และมี FK ไป borrow_request
// result / kind / status รับเฉพาะค่าตาม plan ข้อ 3, ปฏิเสธต้องมีเหตุผล (FR-APV-01 ง) และประวัติแก้หรือลบไม่ได้ (RQ-25)
func TestMigration_ApprovalTables(t *testing.T) {
	conn := testDB(t)
	const req = 9101

	// ค่า result ที่ plan ข้อ 3 กำหนด บันทึกได้ทั้ง 3 ค่า
	for _, c := range []struct{ result, reason string }{
		{"approved", ""},
		{"rejected", "อุปกรณ์ต้องใช้สอบ"},
		{"auto_rejected", "คำขอยังคงค้างอยู่และไม่ได้รับการพิจารณาภายในเวลาที่กำหนด"},
	} {
		if err := execInTx(t, conn, approvalTestRequest(req), approvalTestDecision(req, c.result, c.reason)); err != nil {
			t.Errorf("result %q ต้องบันทึกได้: %v", c.result, err)
		}
	}
	for _, bad := range []string{"approve", "cancelled", "Approved", ""} {
		if err := execInTx(t, conn, approvalTestRequest(req), approvalTestDecision(req, bad, "เหตุผล")); err == nil {
			t.Errorf("result %q ไม่อยู่ใน plan ข้อ 3 ต้องถูกปฏิเสธ", bad)
		}
	}

	// FR-APV-01 (ง): ปฏิเสธต้องมีเหตุผล เหตุผลว่างหรือมีแต่ช่องว่างถูกปฏิเสธ
	for _, result := range []string{"rejected", "auto_rejected"} {
		for _, reason := range []string{"", "   "} {
			if err := execInTx(t, conn, approvalTestRequest(req), approvalTestDecision(req, result, reason)); err == nil {
				t.Errorf("result %q ที่เหตุผล %q ต้องถูกปฏิเสธ (FR-APV-01 ง)", result, reason)
			}
		}
	}

	// (ก) UNIQUE: 1 คำขอมีผลเดียว พิจารณาซ้ำไม่ได้
	if err := execInTx(t, conn, approvalTestRequest(req),
		approvalTestDecision(req, "approved", ""),
		approvalTestDecision(req, "rejected", "ซ้ำ"),
	); err == nil {
		t.Error("ผลการพิจารณาซ้ำของคำขอเดิมต้องถูกปฏิเสธ (request_id UNIQUE)")
	}
	// (ก) FK: อ้างคำขอที่ไม่มีอยู่จริงไม่ได้
	if err := execInTx(t, conn, approvalTestDecision(99999, "approved", "")); err == nil {
		t.Error("request_id ที่ไม่มีในตาราง borrow_request ต้องถูกปฏิเสธ")
	}

	// RQ-25, plan ข้อ 3: ประวัติแก้หรือลบไม่ได้
	for name, stmt := range map[string]string{
		"UPDATE":   `UPDATE approval_decision SET result = 'rejected', reason = 'แก้ผล'`,
		"DELETE":   `DELETE FROM approval_decision`,
		"TRUNCATE": `TRUNCATE approval_decision CASCADE`,
	} {
		if err := execInTx(t, conn, approvalTestRequest(req), approvalTestDecision(req, "approved", ""), stmt); err == nil {
			t.Errorf("%s approval_decision ต้องถูกปฏิเสธ (append-only)", name)
		}
	}

	// approval_email: kind และ status รับเฉพาะค่าตาม plan ข้อ 3 และ 1 คำขอมี Email เดียว
	withDecision := []string{approvalTestRequest(req), approvalTestDecision(req, "approved", "")}
	for _, kind := range []string{"approved", "rejected", "auto_rejected"} {
		for _, status := range []string{"pending", "sent", "failed"} {
			if err := execInTx(t, conn, append(withDecision, approvalTestEmail(req, kind, status))...); err != nil {
				t.Errorf("Email kind %q status %q ต้องบันทึกได้: %v", kind, status, err)
			}
		}
	}
	for _, c := range []struct{ kind, status string }{{"cancelled", "pending"}, {"approved", "Failed"}, {"approved", "retrying"}} {
		if err := execInTx(t, conn, append(withDecision, approvalTestEmail(req, c.kind, c.status))...); err == nil {
			t.Errorf("Email kind %q status %q ไม่อยู่ใน plan ข้อ 3 ต้องถูกปฏิเสธ", c.kind, c.status)
		}
	}
	if err := execInTx(t, conn, append(withDecision,
		approvalTestEmail(req, "approved", "pending"),
		approvalTestEmail(req, "approved", "pending"))...); err == nil {
		t.Error("Email ฉบับที่ 2 ของคำขอเดิมต้องถูกปฏิเสธ (AC-APV-10: Email 1 ฉบับ)")
	}
	if err := execInTx(t, conn, approvalTestRequest(req), approvalTestEmail(req, "approved", "pending")); err == nil {
		t.Error("Email ของคำขอที่ยังไม่มีผลการพิจารณาต้องถูกปฏิเสธ (FK ไป approval_decision)")
	}
}

// T-01 (004) เสร็จเมื่อ (ค): audit_log.action รับ approve, reject, auto_reject ได้ (RQ-25)
// และค่าเดิมของ 001 / 002 / 003 ยังใช้ได้ ค่านอกรายการยังถูกปฏิเสธ
func TestMigration_ApprovalAuditActions(t *testing.T) {
	conn := testDB(t)
	const req = 9102

	// notify ของ 003 (0005) ต้องยังใช้ได้หลัง 0008 เขียน CHECK ใหม่
	for _, a := range []string{"list", "search", "detail", "cancel", "request", "notify", "approve", "reject", "auto_reject"} {
		stmt := fmt.Sprintf(`INSERT INTO audit_log (actor_id, actor_role, action, request_id, outcome)
			VALUES ('staff-001', 'staff', '%s', %d, 'ok')`, a, req)
		if err := execInTx(t, conn, approvalTestRequest(req), stmt); err != nil {
			t.Errorf("action %q ต้องบันทึกได้หลัง migration 0008: %v", a, err)
		}
	}
	for _, a := range []string{"approved", "auto-reject", "delete", ""} {
		stmt := fmt.Sprintf(`INSERT INTO audit_log (actor_id, actor_role, action, outcome)
			VALUES ('staff-001', 'staff', '%s', 'ok')`, a)
		if err := execInTx(t, conn, stmt); err == nil {
			t.Errorf("action %q ไม่อยู่ในรายการ ต้องถูกปฏิเสธ", a)
		}
	}
}

// T-03 (004) เสร็จเมื่อ: approve_dev.sql โหลดหลัง seed.sql ได้ (ร่วมกับ handover_dev.sql ของ 005)
// และมีกรณีครบตามที่จุดตรวจของ 004 ใช้ ทำใน transaction ที่ยกเลิกทิ้ง seed ของ test อื่นจึงไม่เปลี่ยน
func TestSeed_ApproveDevScenarios(t *testing.T) {
	conn := testDB(t)
	ctx := context.Background()

	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)

	for _, f := range []string{"../../db/seed/handover_dev.sql", "../../db/seed/approve_dev.sql"} {
		sql, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			t.Fatalf("โหลด %s ไม่สำเร็จ: %v", f, err)
		}
	}

	count := func(q string) int {
		t.Helper()
		var n int
		if err := tx.QueryRow(ctx, q).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	const base = `SELECT count(*) FROM borrow_request r
		JOIN equipment_item i ON i.id = r.item_id
		JOIN equipment_model m ON m.id = i.model_id
		WHERE i.asset_code LIKE 'AP-TEST-%' `

	if n := count(base + `AND r.status = 'requested' AND r.loan_type = 'normal' AND m.department_id = 1`); n < 3 {
		t.Errorf("ต้องมีคำขอยืมปกติ requested ของภาควิชา 1 อย่างน้อย 3 ใบ แต่ได้ %d", n)
	}
	if n := count(base + `AND i.asset_code = 'AP-TEST-0001' AND r.status = 'requested' AND r.loan_type = 'normal'`); n < 2 {
		t.Errorf("AP-TEST-0001 ต้องมีคำขอยืมปกติอย่างน้อย 2 ใบ (AC-APV-03, 07) แต่ได้ %d", n)
	}
	if n := count(base + `AND i.asset_code = 'AP-TEST-0001' AND r.loan_type = 'overnight' AND r.status = 'requested'`); n != 1 {
		t.Errorf("AP-TEST-0001 ต้องมีคำขอค้างคืน 1 ใบ (AC-APV-17) แต่ได้ %d", n)
	}
	if n := count(base + `AND m.department_id = 2 AND r.status = 'requested'`); n < 1 {
		t.Errorf("ต้องมีคำขอของอุปกรณ์ภาควิชา 2 อย่างน้อย 1 ใบ (AC-APV-11, 12) แต่ได้ %d", n)
	}
	if n := count(base + `AND i.status = 'borrowed' AND r.status = 'requested'`); n < 1 {
		t.Errorf("ต้องมีคำขอบนอุปกรณ์ที่ถูกยืม (AC-APV-08) แต่ได้ %d", n)
	}
	if n := count(base + `AND r.status = 'cancelled'`); n != 1 {
		t.Errorf("ต้องมีคำขอที่ถูกยกเลิก 1 ใบ (AC-APV-09) แต่ได้ %d", n)
	}
	// คำขอแรกของ AP-TEST-0001 ต้องเป็นยืมปกติ (ลำดับตาม submitted_at, DOM-REQ-01)
	var firstType string
	if err := tx.QueryRow(ctx, `SELECT r.loan_type FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id
		WHERE i.asset_code = 'AP-TEST-0001' ORDER BY r.submitted_at LIMIT 1`).Scan(&firstType); err != nil {
		t.Fatal(err)
	}
	if firstType != "normal" {
		t.Errorf("คำขอแรกของ AP-TEST-0001 ต้องเป็นยืมปกติ แต่ได้ %q", firstType)
	}
}

// ---- 0011 คำขอยืมค้างคืน (T-40, CR-009) ----

func approvalTestStep(requestID int, status string) string {
	return fmt.Sprintf(`INSERT INTO approval_step (request_id, status, actor_id) VALUES (%d, '%s', 'staff-001')`, requestID, status)
}

// T-40 เสร็จเมื่อ: approval_decision.advisor_name ว่างได้ (FR-APV-03 ข ช, Q-APV-44)
// approval_step มี FK ไป borrow_request รับ status 3 ค่า และแก้หรือลบไม่ได้ (FR-APV-03 ข, FR-APV-01 ฌ)
func TestMigration_ApprovalOvernight(t *testing.T) {
	conn := testDB(t)
	const req = 99011

	// advisor_name ว่างได้ (คำขอปกติ) และเก็บชื่อได้ (คำขอค้างคืน)
	if err := execInTx(t, conn, approvalTestRequest(req), approvalTestDecision(req, "approved", "")); err != nil {
		t.Errorf("ผลการพิจารณาที่ไม่มีชื่ออาจารย์ต้องบันทึกได้: %v", err)
	}
	if err := execInTx(t, conn, approvalTestRequest(req),
		fmt.Sprintf(`INSERT INTO approval_decision (request_id, result, actor_id, actor_department_id, advisor_name)
			VALUES (%d, 'approved', 'staff-001', 1, 'อ.ทดสอบ ใจดี')`, req)); err != nil {
		t.Errorf("ผลการพิจารณาที่มีชื่ออาจารย์ต้องบันทึกได้: %v", err)
	}

	// status รับเฉพาะ 3 ค่า
	for _, s := range []string{"approved_by_officer", "pending_advisor_approval", "approved"} {
		if err := execInTx(t, conn, approvalTestRequest(req), approvalTestStep(req, s)); err != nil {
			t.Errorf("approval_step.status = %q ต้องบันทึกได้: %v", s, err)
		}
	}
	for _, s := range []string{"requested", "rejected", "auto_rejected", ""} {
		if err := execInTx(t, conn, approvalTestRequest(req), approvalTestStep(req, s)); err == nil {
			t.Errorf("approval_step.status = %q ต้องถูกปฏิเสธ", s)
		}
	}

	// FK: อ้างคำขอที่ไม่มีอยู่จริงไม่ได้
	if err := execInTx(t, conn, approvalTestStep(99999, "approved")); err == nil {
		t.Error("approval_step.request_id ที่ไม่มีในตาราง borrow_request ต้องถูกปฏิเสธ")
	}

	// append-only: UPDATE / DELETE / TRUNCATE ไม่ได้
	for name, s := range map[string]string{
		"UPDATE":   `UPDATE approval_step SET status = 'approved' WHERE request_id = 99011`,
		"DELETE":   `DELETE FROM approval_step WHERE request_id = 99011`,
		"TRUNCATE": `TRUNCATE approval_step`,
	} {
		if err := execInTx(t, conn, approvalTestRequest(req), approvalTestStep(req, "approved_by_officer"), s); err == nil {
			t.Errorf("%s บน approval_step ต้องถูกปฏิเสธ (append-only)", name)
		}
	}
}

// T-40 เสร็จเมื่อ: down ย้อนได้ รันไฟล์ down ใน transaction ที่ยกเลิกทิ้ง แล้วตรวจว่าตารางและคอลัมน์หายไป
func TestMigration_ApprovalOvernightDown(t *testing.T) {
	conn := testDB(t)
	down, err := os.ReadFile(migrationsDir + "/0011_approval_overnight.down.sql")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, string(down)); err != nil {
		t.Fatalf("down ของ 0011 ต้องรันได้: %v", err)
	}
	var n int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM information_schema.tables WHERE table_name = 'approval_step'`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Error("หลัง down ตาราง approval_step ต้องหายไป")
	}
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM information_schema.columns
		WHERE table_name = 'approval_decision' AND column_name = 'advisor_name'`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Error("หลัง down คอลัมน์ approval_decision.advisor_name ต้องหายไป")
	}
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM pg_proc WHERE proname = 'approval_step_append_only'`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Error("หลัง down ฟังก์ชัน approval_step_append_only ต้องหายไป")
	}
}
