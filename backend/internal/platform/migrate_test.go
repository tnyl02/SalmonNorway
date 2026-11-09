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
