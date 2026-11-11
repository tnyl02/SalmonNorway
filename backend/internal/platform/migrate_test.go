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
	for _, a := range []string{"cancel_notify_failed", "approve", ""} {
		stmt := `INSERT INTO audit_log (actor_id, actor_role, action, outcome) VALUES ('student-001', 'student', '` + a + `', 'ok')`
		if err := execInTx(t, conn, stmt); err == nil {
			t.Errorf("action %q ไม่อยู่ใน DM-07 ต้องถูกปฏิเสธ", a)
		}
	}
	// การยกเลิกคำขอของ 003: มี request_id ไม่มีรหัสครุภัณฑ์ / การเปิดรายละเอียดของ 001: กลับกัน
	if err := execInTx(t, conn,
		`INSERT INTO audit_log (actor_id, actor_role, action, request_id, outcome) VALUES ('student-001', 'student', 'cancel', 42, 'ok')`,
		insertAudit,
	); err != nil {
		t.Errorf("request_id และ asset_code ต้องว่างได้ทั้งคู่: %v", err)
	}
}

// 006 T-01 เสร็จเมื่อ: audit_log รับ action ใหม่ของ 006 ทั้ง 4 ค่า (DOM-AUDIT-01, FR-RET-07)
func TestMigration006_AuditLogReturnActions(t *testing.T) {
	conn := testDB(t)

	for _, a := range []string{"return_list", "return_lookup", "return_detail", "return"} {
		stmt := `INSERT INTO audit_log (actor_id, actor_role, action, asset_code, outcome) VALUES ('staff-001', 'staff', '` + a + `', 'CS-NB-0004', 'ok')`
		if err := execInTx(t, conn, stmt); err != nil {
			t.Errorf("action %q ของ 006 ต้องบันทึกได้ แต่ได้ error: %v", a, err)
		}
	}
}

// 006 T-01 เสร็จเมื่อ: loan มี return_condition (ปกติ/ชำรุด เท่านั้น) และ returned_late
// และคืนแล้วต้องมีครบทั้งเวลาคืน สภาพ เครื่องหมายคืนเกินกำหนด (FR-RET-02, FR-RET-06)
func TestMigration006_LoanReturnColumns(t *testing.T) {
	conn := testDB(t)

	// การยืมของ CS-NB-0004 (item_id 4) ในข้อมูลตั้งต้น ยังไม่คืน
	const loanOf4 = `(SELECT id FROM loan WHERE item_id = 4 AND returned_at IS NULL)`

	ok := map[string]string{
		"ยังไม่คืน ว่างทั้ง 3 ค่า (ข้อมูลเดิมของ 001)": `SELECT 1`,
		"คืนสภาพปกติ":  `UPDATE loan SET returned_at = now(), return_condition = 'normal', returned_late = false WHERE id = ` + loanOf4,
		"คืนสภาพชำรุด": `UPDATE loan SET returned_at = now(), return_condition = 'damaged', returned_late = true WHERE id = ` + loanOf4,
	}
	for name, stmt := range ok {
		if err := execInTx(t, conn, stmt); err != nil {
			t.Errorf("%s ต้องบันทึกได้: %v", name, err)
		}
	}

	bad := map[string]string{
		"สภาพสูญหาย (ไม่ใช่สภาพ FR-RET-02)":   `UPDATE loan SET returned_at = now(), return_condition = 'lost', returned_late = false WHERE id = ` + loanOf4,
		"คืนแล้วแต่ไม่มีสภาพ":                 `UPDATE loan SET returned_at = now(), returned_late = false WHERE id = ` + loanOf4,
		"คืนแล้วแต่ไม่มีเครื่องหมายเกินกำหนด": `UPDATE loan SET returned_at = now(), return_condition = 'normal' WHERE id = ` + loanOf4,
		"มีสภาพแต่ไม่มีเวลาคืน":               `UPDATE loan SET return_condition = 'normal', returned_late = false WHERE id = ` + loanOf4,
	}
	for name, stmt := range bad {
		if err := execInTx(t, conn, stmt); err == nil {
			t.Errorf("%s ต้องถูกปฏิเสธ", name)
		}
	}
}
