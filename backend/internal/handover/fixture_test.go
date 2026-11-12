package handover

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"salmonnorway/backend/internal/platform"
)

// ข้อมูลและฐานข้อมูลสำหรับ test ของ 005 (แบบเดียวกับ platform/migrate_test.go ของ 001)
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

// TestMain ปิด container และลบโฟลเดอร์รูปของ test หลัง test ทั้ง package จบ test ที่ไม่ใช้ฐานข้อมูลไม่ต้องมี Docker
func TestMain(m *testing.M) {
	code := m.Run()
	_ = os.RemoveAll(testPhotoRoot)
	if dbContainer != nil {
		_ = testcontainers.TerminateContainer(dbContainer)
	}
	os.Exit(code)
}

// testDB เปิด PostgreSQL 16 ใน container ครั้งเดียวต่อการรัน test แล้วรัน migration ทั้งหมดและข้อมูลตั้งต้นของ 001
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
		if dbErr = platform.Migrate(dbURL, migrationsDir); dbErr != nil {
			return
		}
		dbErr = platform.Seed(ctx, dbURL, seedFile)
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

// rollbackTx เปิด transaction ที่ยกเลิกทิ้งเสมอเมื่อ test จบ ข้อมูลจึงไม่ค้างข้าม test
func rollbackTx(t *testing.T, conn *pgx.Conn) pgx.Tx {
	t.Helper()
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })
	return tx
}

// insertRequest เพิ่มคำขอยืมแบบปกติของอุปกรณ์ itemID ในสถานะ status (ตาราง borrow_request ของ 002) คืน id
func insertRequest(t *testing.T, tx pgx.Tx, requesterID string, itemID int, status string) int {
	t.Helper()
	var id int
	err := tx.QueryRow(context.Background(),
		`INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
		 VALUES ($1, 'ผู้ขอทดสอบ', $1, $2, 'normal', 'study', now(), now(), $3) RETURNING id`,
		requesterID, itemID, status).Scan(&id)
	if err != nil {
		t.Fatalf("เพิ่มคำขอยืมไม่สำเร็จ: %v", err)
	}
	return id
}

// insertLoan เพิ่มการยืมแบบปกติของอุปกรณ์ itemID ด้วยคอลัมน์เดิมของ 001 (DM-05) คืน id
func insertLoan(t *testing.T, tx pgx.Tx, borrowerID string, itemID int) int {
	t.Helper()
	var id int
	err := tx.QueryRow(context.Background(),
		`INSERT INTO loan (item_id, borrower_id, loan_type, due_at) VALUES ($1, $2, 'normal', now()) RETURNING id`,
		itemID, borrowerID).Scan(&id)
	if err != nil {
		t.Fatalf("เพิ่มการยืมไม่สำเร็จ: %v", err)
	}
	return id
}
