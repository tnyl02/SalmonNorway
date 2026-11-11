// Package migrationtest ตรวจตารางที่ 003 ใช้หลังรัน migration ทั้งหมด (T-01)
//
// แยกเป็น package ของตัวเอง เพราะ package cancelrequest ยัง compile ไม่ผ่านจนกว่า T-04, T-06 จะเสร็จ
// (ไฟล์ TC เรียก NewHandler ล่วงหน้า) ทีมตัดสิน 2569-11-11 ตอนทำ T-01
package migrationtest

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/platform"
)

const (
	migrationsDir = "../../../db/migrations"
	seedFile      = "../../../db/seed/seed.sql"
	itemAvailable = 1 // CS-NB-0001 จาก seed.sql ของ 001
)

var (
	dbOnce      sync.Once
	dbPool      *pgxpool.Pool
	dbErr       error
	dbContainer *tcpostgres.PostgresContainer
)

func TestMain(m *testing.M) {
	code := m.Run()
	if dbPool != nil {
		dbPool.Close()
	}
	if dbContainer != nil {
		_ = testcontainers.TerminateContainer(dbContainer)
	}
	os.Exit(code)
}

// pool เปิด PostgreSQL 16 (CON-TECH-01) ครั้งเดียว แล้วรัน migration ทั้งหมดและ seed ของ 001
func pool(t *testing.T) *pgxpool.Pool {
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
		var url string
		if url, dbErr = dbContainer.ConnectionString(ctx, "sslmode=disable"); dbErr != nil {
			return
		}
		if dbErr = platform.Migrate(url, migrationsDir); dbErr != nil {
			return
		}
		if dbErr = platform.Seed(ctx, url, seedFile); dbErr != nil {
			return
		}
		dbPool, dbErr = pgxpool.New(ctx, url)
	})
	if dbErr != nil {
		t.Fatalf("เตรียมฐานข้อมูลไม่สำเร็จ: %v", dbErr)
	}
	return dbPool
}

// insertRequest ใส่คำขอ 1 แถวใน transaction ที่ rollback ทิ้ง เพื่อไม่ทิ้งข้อมูลไว้ (audit_log เป็น append-only)
// คอลัมน์ที่ 003 ไม่ใช้ใส่ค่าทดสอบตายตัว ไม่ได้มาจาก spec ของ 003
func insertRequest(t *testing.T, status string, cancelReason *string) error {
	t.Helper()
	ctx := context.Background()
	tx, err := pool(t).Begin(ctx)
	if err != nil {
		t.Fatalf("เปิด transaction ไม่สำเร็จ: %v", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose,
		                             start_at, end_at, status, cancel_reason)
		 VALUES ('student-001', 'ผู้ทดสอบ', 'TEST-000', $1, 'normal', 'study',
		         now() + interval '1 day', now() + interval '1 day 6 hours', $2, $3)`,
		itemAvailable, status, cancelReason)
	return err
}

// CON-CXL-02: ตารางของ 002 รับค่าสถานะครบ 8 ค่าที่ 003 ใช้ และปฏิเสธค่าอื่น
func TestMigration_BorrowRequestStatusCheck(t *testing.T) {
	statuses := []string{"requested", "approved_by_officer", "pending_advisor_approval",
		"approved", "rejected", "handed_over", "returned", "cancelled"}
	for _, s := range statuses {
		if err := insertRequest(t, s, nil); err != nil {
			t.Errorf("สถานะ %q ต้องใส่ได้: %v", s, err)
		}
	}
	// สถานะที่ยกเลิกได้ใน config ต้องเป็นค่าที่ตารางรับ
	for _, s := range config.CancellableStatuses() {
		if err := insertRequest(t, s, nil); err != nil {
			t.Errorf("สถานะที่ยกเลิกได้ %q ต้องใส่ได้: %v", s, err)
		}
	}
	for _, s := range []string{"ยกเลิก", "Cancelled", "pending"} {
		if err := insertRequest(t, s, nil); err == nil {
			t.Errorf("สถานะ %q นอก 8 ค่าต้องถูกปฏิเสธ", s)
		}
	}
}

// FR-CXL-08, ASM-09: cancel_reason เก็บได้ไม่เกิน CancelReasonMaxLen ตัวอักษร นับตาม code point เหมือนหน้าจอ
func TestMigration_CancelReasonMaxLen(t *testing.T) {
	ok := strings.Repeat("ที่", config.CancelReasonMaxLen/3) + "ab" // 66*3 + 2 = 200 code point
	if err := insertRequest(t, "cancelled", &ok); err != nil {
		t.Errorf("เหตุผล 200 ตัวอักษรต้องเก็บได้: %v", err)
	}
	tooLong := strings.Repeat("ที่", 67) // 201 code point
	if err := insertRequest(t, "cancelled", &tooLong); err == nil {
		t.Errorf("เหตุผล 201 ตัวอักษรต้องถูกปฏิเสธ")
	}
}

// ASM-10: ตาราง staff_contact มีอยู่ และอ้างภาควิชาได้
func TestMigration_StaffContactExists(t *testing.T) {
	ctx := context.Background()
	tx, err := pool(t).Begin(ctx)
	if err != nil {
		t.Fatalf("เปิด transaction ไม่สำเร็จ: %v", err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `INSERT INTO staff_contact (department_id, email) VALUES (1, 'staff@example.test')`); err != nil {
		t.Fatalf("ใส่ staff_contact ไม่สำเร็จ: %v", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO staff_contact (department_id, email) VALUES (9999, 'x@example.test')`); err == nil {
		t.Errorf("ภาควิชาที่ไม่มีอยู่ต้องถูกปฏิเสธ")
	}
}

// CON-CXL-01, CON-DATA-01: ไม่มี trigger บน borrow_request และไม่มี FK จาก equipment_item ไป borrow_request
// การยกเลิกจึงไม่เปลี่ยนสถานะอุปกรณ์ทางอ้อม
func TestMigration_NoEquipmentSideEffects(t *testing.T) {
	ctx := context.Background()
	var triggers int
	if err := pool(t).QueryRow(ctx,
		`SELECT count(*) FROM pg_trigger WHERE tgrelid = 'borrow_request'::regclass AND NOT tgisinternal`).Scan(&triggers); err != nil {
		t.Fatal(err)
	}
	if triggers != 0 {
		t.Errorf("borrow_request ต้องไม่มี trigger แต่มี %d", triggers)
	}
	var fks int
	if err := pool(t).QueryRow(ctx,
		`SELECT count(*) FROM pg_constraint
		 WHERE contype = 'f' AND conrelid = 'equipment_item'::regclass AND confrelid = 'borrow_request'::regclass`).Scan(&fks); err != nil {
		t.Fatal(err)
	}
	if fks != 0 {
		t.Errorf("equipment_item ต้องไม่มี FK ไป borrow_request แต่มี %d", fks)
	}
}
