package handover

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

const insertHandover = `INSERT INTO handover (loan_id, request_id, handed_over_at, handed_over_by,
	borrower_student_id, borrower_first_name, borrower_last_name, identity_method, card_checked)
	VALUES ($1, $2, now(), 'staff-001', 'student-001', 'ชื่อ', 'นามสกุล', $3, $4)`

// newLoanAndRequest เตรียมคำขอที่อนุมัติแล้วและการยืมของอุปกรณ์ชิ้นที่ 1 ใน transaction ที่ยกเลิกทิ้ง
func newLoanAndRequest(t *testing.T, tx pgx.Tx) (loanID, requestID int) {
	t.Helper()
	requestID = insertRequest(t, tx, "student-001", 1, "approved")
	loanID = insertLoan(t, tx, "student-001", 1)
	return loanID, requestID
}

// T-02 เสร็จเมื่อ: identity_method นอก auth_api / card_only ถูกปฏิเสธ และ card_checked ต้องเป็น true (FR-HND-02, 03)
func TestMigration_HandoverIdentityMethod(t *testing.T) {
	conn := testDB(t)
	ctx := context.Background()

	for _, tc := range []struct {
		name    string
		method  string
		checked bool
		wantOK  bool
	}{
		{"auth_api และตรวจบัตรแล้ว", "auth_api", true, true},
		{"card_only และตรวจบัตรแล้ว", "card_only", true, true},
		{"วิธีนอกรายการ", "password", true, false},
		{"ไม่ได้ตรวจบัตร", "auth_api", false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tx := rollbackTx(t, conn)
			loanID, requestID := newLoanAndRequest(t, tx)
			_, err := tx.Exec(ctx, insertHandover, loanID, requestID, tc.method, tc.checked)
			if tc.wantOK && err != nil {
				t.Fatalf("ต้องบันทึกได้ แต่ได้ error: %v", err)
			}
			if !tc.wantOK && err == nil {
				t.Fatal("ต้องถูกปฏิเสธ แต่บันทึกได้")
			}
		})
	}
}

// T-02 เสร็จเมื่อ: content_type นอก image/jpeg, image/png ถูกปฏิเสธ (FR-HND-06, ASM-HND-07)
func TestMigration_HandoverPhotoContentType(t *testing.T) {
	conn := testDB(t)
	ctx := context.Background()

	for _, tc := range []struct {
		contentType string
		wantOK      bool
	}{
		{"image/jpeg", true},
		{"image/png", true},
		{"image/gif", false},
		{"application/pdf", false},
	} {
		t.Run(tc.contentType, func(t *testing.T) {
			tx := rollbackTx(t, conn)
			loanID, requestID := newLoanAndRequest(t, tx)
			if _, err := tx.Exec(ctx, insertHandover, loanID, requestID, "auth_api", true); err != nil {
				t.Fatal(err)
			}
			_, err := tx.Exec(ctx,
				`INSERT INTO handover_photo (loan_id, file_path, content_type, size_bytes) VALUES ($1, 'x', $2, 1)`,
				loanID, tc.contentType)
			if tc.wantOK && err != nil {
				t.Fatalf("ต้องบันทึกได้ แต่ได้ error: %v", err)
			}
			if !tc.wantOK && err == nil {
				t.Fatal("ต้องถูกปฏิเสธ แต่บันทึกได้")
			}
		})
	}
}

// T-02 เสร็จเมื่อ: ไม่มีคอลัมน์เลขบัตรใด ๆ ในตารางของ 005 (FR-HND-03 บันทึกแค่ชื่อ นามสกุล รหัสนักศึกษา เวลา)
func TestMigration_HandoverNoCardNumberColumn(t *testing.T) {
	conn := testDB(t)
	rows, err := conn.Query(context.Background(),
		`SELECT table_name, column_name FROM information_schema.columns
		 WHERE table_name IN ('handover', 'handover_photo')
		   AND (column_name ILIKE '%card_no%' OR column_name ILIKE '%card_number%'
		        OR column_name ILIKE '%national%' OR column_name ILIKE '%citizen%' OR column_name ILIKE '%id_card%')`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var table, column string
		if err := rows.Scan(&table, &column); err != nil {
			t.Fatal(err)
		}
		t.Errorf("พบคอลัมน์ที่อาจเป็นเลขบัตร: %s.%s", table, column)
	}
}
