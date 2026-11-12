package platform

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	migrationsDir = "../../db/migrations"
	seedFile      = "../../db/seed/seed.sql"
)

var (
	dbOnce      sync.Once
	dbURL       string
	dbErr       error
	dbContainer *tcpostgres.PostgresContainer
)

// TestMain ปิด container หลัง test ทั้ง package จบ (ถ้ามีการเปิด) test ที่ไม่ใช้ฐานข้อมูลไม่ต้องมี Docker
func TestMain(m *testing.M) {
	code := m.Run()
	if dbContainer != nil {
		_ = testcontainers.TerminateContainer(dbContainer)
	}
	os.Exit(code)
}

// testDB เปิด PostgreSQL 16 ใน container ครั้งเดียวต่อการรัน test (plan.md ข้อ 2: testcontainers)
// แล้วรัน migration และข้อมูลตั้งต้น
func testDB(t *testing.T) *pgx.Conn {
	t.Helper()
	if testing.Short() {
		t.Skip("ต้องใช้ Docker: ข้ามเมื่อรันด้วย -short")
	}
	dbOnce.Do(func() {
		ctx := context.Background()
		dbContainer, dbErr = tcpostgres.Run(ctx, "postgres:16",
			tcpostgres.WithDatabase("salmon"),
			tcpostgres.WithUsername("salmon"),
			tcpostgres.WithPassword("salmon"),
			tcpostgres.BasicWaitStrategies(),
		)
		if dbErr != nil {
			return
		}
		if dbURL, dbErr = dbContainer.ConnectionString(ctx, "sslmode=disable"); dbErr != nil {
			return
		}
		if dbErr = Migrate(dbURL, migrationsDir); dbErr != nil {
			return
		}
		dbErr = Seed(ctx, dbURL, seedFile)
	})
	if dbErr != nil {
		t.Fatalf("เตรียมฐานข้อมูลไม่สำเร็จ: %v", dbErr)
	}
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close(context.Background()) })
	return conn
}

// insertItem พยายามเพิ่มอุปกรณ์ 1 ชิ้นใน transaction ที่ยกเลิกทิ้งเสมอ ข้อมูลจึงไม่ค้างข้าม test
func insertItem(t *testing.T, conn *pgx.Conn, assetCode, status string) error {
	t.Helper()
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, $2)`, assetCode, status)
	return err
}

// T-02 เสร็จเมื่อ: insert สถานะนอก 3 ค่าถูกปฏิเสธ (CON-DATA-01)
func TestMigration_StatusCheck(t *testing.T) {
	conn := testDB(t)

	for _, s := range []string{"available", "borrowed", "unavailable"} {
		if err := insertItem(t, conn, "TEST-OK-"+s, s); err != nil {
			t.Errorf("สถานะ %q ต้องบันทึกได้ แต่ได้ error: %v", s, err)
		}
	}
	for _, s := range []string{"pending_approval", "damaged", "maintenance", "lost", ""} {
		if err := insertItem(t, conn, "TEST-BAD-"+s, s); err == nil {
			t.Errorf("สถานะ %q ต้องถูกปฏิเสธ (CON-DATA-01 มี 3 ค่าเท่านั้น)", s)
		}
	}
}

// T-02 เสร็จเมื่อ: รหัสครุภัณฑ์ซ้ำถูกปฏิเสธ (CON-DATA-03)
func TestMigration_AssetCodeUnique(t *testing.T) {
	conn := testDB(t)

	// CS-NB-0001 มีอยู่แล้วในข้อมูลตั้งต้น
	if err := insertItem(t, conn, "CS-NB-0001", "available"); err == nil {
		t.Fatal("รหัสครุภัณฑ์ซ้ำต้องถูกปฏิเสธ (CON-DATA-03)")
	}
}

// ข้อมูลตั้งต้นใส่ได้ครบตามที่จุดตรวจของ 001 ใช้ (รุ่น 1: 5 ชิ้น ว่าง 3 ชิ้น)
func TestSeed_Model1HasFiveItemsThreeAvailable(t *testing.T) {
	conn := testDB(t)

	var total, available int
	err := conn.QueryRow(context.Background(),
		`SELECT count(*), count(*) FILTER (WHERE status = 'available') FROM equipment_item WHERE model_id = 1`,
	).Scan(&total, &available)
	if err != nil {
		t.Fatal(err)
	}
	if total != 5 || available != 3 {
		t.Fatalf("รุ่น 1 ได้ %d ชิ้น ว่าง %d ชิ้น ต้องเป็น 5 และ 3", total, available)
	}
}

// execInTx รันคำสั่ง SQL ใน transaction ที่ยกเลิกทิ้งเสมอ ข้อมูลจึงไม่ค้างข้าม test
func execInTx(t *testing.T, conn *pgx.Conn, sqls ...string) error {
	t.Helper()
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	for _, s := range sqls {
		if _, err := tx.Exec(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

const insertAudit = `INSERT INTO audit_log (actor_id, actor_role, action, asset_code, outcome)
	VALUES ('student-001', 'student', 'detail', 'CS-NB-0001', 'ok')`

// T-31 เสร็จเมื่อ: audit_log เป็น append-only คือ INSERT ได้ แต่ UPDATE และ DELETE ถูกปฏิเสธ (DOM-AUDIT-01, plan.md ข้อ 4)
func TestMigration_AuditLogAppendOnly(t *testing.T) {
	conn := testDB(t)

	if err := execInTx(t, conn, insertAudit); err != nil {
		t.Fatalf("INSERT ต้องทำได้: %v", err)
	}
	for name, stmt := range map[string]string{
		"UPDATE":   `UPDATE audit_log SET outcome = 'error'`,
		"DELETE":   `DELETE FROM audit_log`,
		"TRUNCATE": `TRUNCATE audit_log`,
	} {
		if err := execInTx(t, conn, insertAudit, stmt); err == nil {
			t.Errorf("%s ต้องถูกปฏิเสธ (append-only)", name)
		}
	}
}

// T-31 เสร็จเมื่อ: action นอกตาราง DM-07 ถูกปฏิเสธ และ request_id ว่างได้ (คำขอของ 003)
func TestMigration_AuditLogAction(t *testing.T) {
	conn := testDB(t)

	for _, a := range []string{"list", "search", "detail", "cancel"} {
		stmt := `INSERT INTO audit_log (actor_id, actor_role, action, outcome) VALUES ('student-001', 'student', '` + a + `', 'ok')`
		if err := execInTx(t, conn, stmt); err != nil {
			t.Errorf("action %q อยู่ใน DM-07 ต้องบันทึกได้ แต่ได้ error: %v", a, err)
		}
	}
	// 'approve' เคยอยู่ในรายการนี้ ย้ายเป็น 'delete' เพราะ migration 0008 ของ 004 เพิ่ม 'approve' ใน DM-07 แล้ว (ทีมตกลง 2569-11-12)
	for _, a := range []string{"cancel_notify_failed", "delete", ""} {
		stmt := `INSERT INTO audit_log (actor_id, actor_role, action, outcome) VALUES ('student-001', 'student', '` + a + `', 'ok')`
		if err := execInTx(t, conn, stmt); err == nil {
			t.Errorf("action %q ไม่อยู่ใน DM-07 ต้องถูกปฏิเสธ", a)
		}
	}
	// การยกเลิกคำขอของ 003: มี request_id ไม่มีรหัสครุภัณฑ์ / การเปิดรายละเอียดของ 001: กลับกัน
	// ตั้งแต่ migration 0003 ของ 002 คอลัมน์ request_id มี foreign key แล้ว (ASM-BRQ-03)
	// จึงต้องมีคำขอจริงก่อน ไม่ใช่เลขสมมติ เจตนาเดิมของ test ยังเหมือนเดิมคือทั้งสองคอลัมน์ว่างได้
	if err := execInTx(t, conn,
		`INSERT INTO borrow_request (id, requester_id, requester_name, student_code, item_id,
			loan_type, purpose, start_at, end_at, status)
		 VALUES (42, 'student-001', 'ผู้ขอทดสอบ', '660710715', 1, 'normal', 'study', now(), now(), 'requested')`,
		`INSERT INTO audit_log (actor_id, actor_role, action, request_id, outcome) VALUES ('student-001', 'student', 'cancel', 42, 'ok')`,
		insertAudit,
	); err != nil {
		t.Errorf("request_id และ asset_code ต้องว่างได้ทั้งคู่: %v", err)
	}
}

// insertRequest พยายามเพิ่มคำขอยืม 1 ใบใน transaction ที่ยกเลิกทิ้งเสมอ (T-01 ของ 002)
func insertRequest(t *testing.T, conn *pgx.Conn, requesterID string, itemID int, status string, extra ...string) error {
	t.Helper()
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	for _, s := range extra {
		if _, err := tx.Exec(ctx, s); err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, `INSERT INTO borrow_request
		(requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
		VALUES ($1, 'ธัญลักษณ์ ทดสอบ', '660710715', $2, 'normal', 'study', now(), now(), $3)`,
		requesterID, itemID, status)
	return err
}

// T-01 (002) เสร็จเมื่อ (ก): สถานะของคำขอนอก 8 ค่าถูกปฏิเสธ (data-model.md v9, FR-BRQ-04)
func TestMigration_BorrowRequestConstraints(t *testing.T) {
	conn := testDB(t)

	for _, s := range []string{"requested", "approved_by_officer", "pending_advisor_approval",
		"approved", "rejected", "handed_over", "returned", "cancelled"} {
		if err := insertRequest(t, conn, "student-001", 1, s); err != nil {
			t.Errorf("สถานะ %q อยู่ใน data-model v9 ต้องบันทึกได้ แต่ได้ error: %v", s, err)
		}
	}
	for _, s := range []string{"Requested", "pending", "approve", ""} {
		if err := insertRequest(t, conn, "student-001", 1, s); err == nil {
			t.Errorf("สถานะ %q ไม่อยู่ในรายการ 8 ค่า ต้องถูกปฏิเสธ", s)
		}
	}
	// อ้างอุปกรณ์ที่ไม่มีอยู่จริงต้องถูกปฏิเสธ (FR-BRQ-01 ขอยืมรายชิ้น)
	if err := insertRequest(t, conn, "student-001", 99999, "requested"); err == nil {
		t.Error("item_id ที่ไม่มีในตาราง equipment_item ต้องถูกปฏิเสธ")
	}
}

// T-01 (002) เสร็จเมื่อ (ข): unique index กันคำขอซ้ำของผู้ขอคนเดิมบนชิ้นเดิมใน 4 สถานะที่ยังไม่จบ
// (FR-BRQ-05 ข ตามขอบเขตที่ CR-004 ขยาย) แต่คำขอของคนอื่นบนชิ้นเดียวกันต้องส่งซ้อนได้ (DOM-REQ-01)
func TestMigration_DuplicateRequestBlocked(t *testing.T) {
	conn := testDB(t)

	const first = `INSERT INTO borrow_request
		(requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
		VALUES ('student-001', 'ก', '1', 1, 'normal', 'study', now(), now(), `

	for _, s := range []string{"requested", "approved_by_officer", "pending_advisor_approval", "approved"} {
		if err := insertRequest(t, conn, "student-001", 1, "requested", first+`'`+s+`')`); err == nil {
			t.Errorf("คำขอใบเดิมสถานะ %q ยังไม่จบ ต้องส่งซ้ำไม่ได้ (FR-BRQ-05 ข, CR-004)", s)
		}
	}
	for _, s := range []string{"rejected", "cancelled", "handed_over", "returned"} {
		if err := insertRequest(t, conn, "student-001", 1, "requested", first+`'`+s+`')`); err != nil {
			t.Errorf("คำขอใบเดิมสถานะ %q จบแล้ว ต้องส่งใหม่ได้ แต่ได้ error: %v", s, err)
		}
	}
	// คนละคนบนอุปกรณ์ชิ้นเดียวกัน ต้องส่งซ้อนได้ (DOM-REQ-01)
	if err := insertRequest(t, conn, "student-002", 1, "requested", first+`'requested')`); err != nil {
		t.Errorf("คำขอของผู้อื่นบนชิ้นเดียวกันต้องส่งซ้อนได้ (DOM-REQ-01): %v", err)
	}
}

// T-01 (002) เสร็จเมื่อ (ค)(ง): audit_log.request_id มี foreign key แล้ว และใช้ค่า action 'request' ได้ (Q-BRQ-07, ASM-BRQ-03)
func TestMigration_AuditLogRequestAction(t *testing.T) {
	conn := testDB(t)

	const insertRequestRow = `INSERT INTO borrow_request
		(id, requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
		VALUES (9001, 'student-001', 'ก', '1', 1, 'normal', 'study', now(), now(), 'requested')`

	if err := execInTx(t, conn, insertRequestRow,
		`INSERT INTO audit_log (actor_id, actor_role, action, asset_code, request_id, outcome)
		 VALUES ('student-001', 'student', 'request', 'CS-NB-0001', 9001, 'ok')`,
	); err != nil {
		t.Errorf("action 'request' พร้อม request_id ที่มีจริง ต้องบันทึกได้: %v", err)
	}
	if err := execInTx(t, conn,
		`INSERT INTO audit_log (actor_id, actor_role, action, request_id, outcome)
		 VALUES ('student-001', 'student', 'request', 9999, 'ok')`,
	); err == nil {
		t.Error("request_id ที่ไม่มีในตาราง borrow_request ต้องถูกปฏิเสธด้วย foreign key (ASM-BRQ-03)")
	}
}

// T-03 (002) เสร็จเมื่อ: ข้อมูลตั้งต้นมีครบ 3 กรณีที่จุดตรวจของฟีเจอร์ 002 ต้องใช้
// เป็นเกราะกันไม่ให้ใครแก้ seed แล้วทำให้จุดตรวจของ 002 ใช้ไม่ได้โดยไม่รู้ตัว
func TestSeed_BorrowRequestScenarios(t *testing.T) {
	conn := testDB(t)
	ctx := context.Background()

	// (1) อุปกรณ์ว่างของภาควิชา 1 สำหรับส่งคำขอสำเร็จ (FR-BRQ-01)
	var availableOwn int
	if err := conn.QueryRow(ctx, `
		SELECT count(*) FROM equipment_item i
		JOIN equipment_model m ON m.id = i.model_id
		WHERE i.status = 'available' AND m.department_id = 1`).Scan(&availableOwn); err != nil {
		t.Fatal(err)
	}
	if availableOwn < 1 {
		t.Error("ต้องมีอุปกรณ์ว่างของภาควิชา 1 อย่างน้อย 1 ชิ้น สำหรับลองส่งคำขอ")
	}

	// (2) อุปกรณ์ว่างของภาควิชาอื่น สำหรับลองคำขอข้ามภาควิชา (FR-BRQ-07, DOM-ACCESS-01)
	var availableOther int
	if err := conn.QueryRow(ctx, `
		SELECT count(*) FROM equipment_item i
		JOIN equipment_model m ON m.id = i.model_id
		WHERE i.status = 'available' AND m.department_id <> 1`).Scan(&availableOther); err != nil {
		t.Fatal(err)
	}
	if availableOther < 1 {
		t.Error("ต้องมีอุปกรณ์ว่างของภาควิชาอื่นอย่างน้อย 1 ชิ้น สำหรับลองคำขอข้ามภาควิชา")
	}

	// (3) ผู้ขอที่มีของคืนเกินกำหนดค้างอยู่ สำหรับลองกฎปฏิเสธอัตโนมัติ (FR-BRQ-06)
	var overdue int
	if err := conn.QueryRow(ctx, `
		SELECT count(*) FROM loan
		WHERE borrower_id = 'student-002' AND returned_at IS NULL AND due_at < now()`).Scan(&overdue); err != nil {
		t.Fatal(err)
	}
	if overdue < 1 {
		t.Error("ต้องมี loan ของ student-002 ที่เลยกำหนดคืนและยังไม่คืน สำหรับลองกฎของค้างเกินกำหนด")
	}

	// ตาราง borrow_request ต้องว่างเปล่า คำขอทุกใบต้องเกิดจาก API ของ 002 เท่านั้น
	var requests int
	if err := conn.QueryRow(ctx, `SELECT count(*) FROM borrow_request`).Scan(&requests); err != nil {
		t.Fatal(err)
	}
	if requests != 0 {
		t.Errorf("ข้อมูลตั้งต้นต้องไม่มีคำขอยืม แต่พบ %d ใบ", requests)
	}
}
