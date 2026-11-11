package platform

// test ของ migration 0008 ฟีเจอร์ 004 (T-01 ใน specs/004-approve-request/tasks.md)
// ใช้ testDB และ execInTx จาก migrate_test.go ทุกคำสั่งรันใน transaction ที่ยกเลิกทิ้ง ข้อมูลจึงไม่ค้างข้าม test

import (
	"fmt"
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

	for _, a := range []string{"list", "search", "detail", "cancel", "request", "approve", "reject", "auto_reject"} {
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
