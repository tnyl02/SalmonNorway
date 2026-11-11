package approve

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"salmonnorway/backend/internal/platform"
)

const (
	migrationsDir = "../../db/migrations"
	seedFile      = "../../db/seed/seed.sql"
	approveDevSQL = "../../db/seed/approve_dev.sql" // ข้อมูลตัวอย่างของ 004 (T-03)
)

var (
	dbOnce      sync.Once
	dbPool      *pgxpool.Pool
	dbErr       error
	dbContainer *tcpostgres.PostgresContainer
)

// TestMain ปิดฐานข้อมูลหลัง test ทั้ง package จบ (แบบเดียวกับ test ของ 002)
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

// testPool เปิด PostgreSQL 16 ใน container ครั้งเดียว แล้วรัน migration และ seed.sql (plan.md ข้อ 2, ข้อ 6)
// ไม่โหลด approve_dev.sql ที่นี่ test แต่ละตัวโหลดเองใน transaction ที่ยกเลิกทิ้ง ข้อมูลจึงไม่ค้างข้าม test
func testPool(t *testing.T) *pgxpool.Pool {
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

// devTx เปิด transaction ที่โหลด approve_dev.sql แล้ว และยกเลิกทิ้งเมื่อ test จบ
func devTx(t *testing.T) pgx.Tx {
	t.Helper()
	ctx := context.Background()
	tx, err := testPool(t).Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })
	sql, err := os.ReadFile(approveDevSQL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tx.Exec(ctx, string(sql)); err != nil {
		t.Fatalf("โหลด approve_dev.sql ไม่สำเร็จ: %v", err)
	}
	return tx
}

// AC-APV-11 (FR-APV-01 ก) ระดับ repository
// Given มีคำขอ requested ของภาควิชาตนเองหลายคำขอ คำขอที่พิจารณาแล้ว และคำขอของภาควิชาอื่น
// When  เจ้าหน้าที่เปิดหน้ารายการ
// Then  เห็นเฉพาะคำขอ requested ของภาควิชาตนเอง เรียงตามเวลาที่ส่งจากเก่าไปใหม่
//
//	พร้อมชื่อผู้ขอ อุปกรณ์ วัตถุประสงค์ ช่วงเวลาที่ต้องการยืม (และคำขอค้างคืนไม่ขึ้น ตาม Q-APV-01)
func TestAC_APV_11(t *testing.T) {
	tx := devTx(t)
	ctx := context.Background()

	// Given: คำขอที่พิจารณาแล้วของภาควิชา 1 เพิ่มจาก approve_dev.sql (ที่มี requested 5, ค้างคืน 1, cancelled 1, ภาควิชา 2 อีก 1)
	if _, err := tx.Exec(ctx, `
		INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose,
		                            start_at, end_at, status, submitted_at)
		SELECT 'student-002', 'พิจารณาแล้ว', 'student-002', id, 'normal', 'study', now(), now(), 'approved',
		       now() - interval '50 minutes'
		FROM equipment_item WHERE asset_code = 'AP-TEST-0005'`); err != nil {
		t.Fatal(err)
	}

	// When
	got, err := NewRepository(tx).ListPending(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Then: เฉพาะยืมปกติ requested ของภาควิชา 1 เรียงตามเวลาที่ส่ง (ลำดับตาม approve_dev.sql)
	want := []struct {
		asset   string
		isFirst bool
	}{
		{"AP-TEST-0001", true},  // A ส่งเมื่อ 30 นาทีก่อน
		{"AP-TEST-0002", true},  // 25 นาที
		{"AP-TEST-0001", false}, // B 20 นาที ไม่ใช่คำขอแรกของชิ้นนี้
		{"AP-TEST-0003", true},  // 15 นาที
		{"AP-TEST-0004", true},  // 12 นาที (อุปกรณ์ถูกยืม ยังขึ้นในรายการ แต่อนุมัติไม่ได้ใน T-19)
	}
	if len(got) != len(want) {
		t.Fatalf("ต้องได้ %d คำขอ แต่ได้ %d: %+v", len(want), len(got), got)
	}
	for i, w := range want {
		g := got[i]
		if g.AssetCode != w.asset || g.IsFirst != w.isFirst {
			t.Errorf("แถว %d ต้องเป็น %s is_first=%v แต่ได้ %s is_first=%v", i, w.asset, w.isFirst, g.AssetCode, g.IsFirst)
		}
		if i > 0 && g.SubmittedAt.Before(got[i-1].SubmittedAt) {
			t.Errorf("แถว %d ต้องเรียงตามเวลาที่ส่งจากเก่าไปใหม่", i)
		}
		// ข้อมูลที่หน้ารายการต้องแสดง (FR-APV-01 ก)
		if g.RequesterName == "" || g.ModelName == "" || g.Purpose == "" || g.StartAt.IsZero() || g.EndAt.IsZero() {
			t.Errorf("แถว %d ข้อมูลไม่ครบ (ชื่อผู้ขอ อุปกรณ์ วัตถุประสงค์ ช่วงเวลา): %+v", i, g)
		}
	}

	// ภาควิชา 2 เห็นเฉพาะคำขอของอุปกรณ์ภาควิชา 2 (ASM-06)
	other, err := NewRepository(tx).ListPending(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(other) != 1 || other[0].AssetCode != "AP-TEST-0101" {
		t.Errorf("ภาควิชา 2 ต้องเห็นเฉพาะ AP-TEST-0101 แต่ได้ %+v", other)
	}
}

// Q-APV-23: คำขอค้างคืนที่ส่งก่อนไม่ทำให้คำขอยืมปกติไม่เป็นคำขอแรก
func TestListPending_OvernightDoesNotBlockQueue(t *testing.T) {
	tx := devTx(t)
	ctx := context.Background()

	// ค้างคืนของ AP-TEST-0002 ส่งก่อนคำขอปกติที่มีอยู่ (25 นาที)
	if _, err := tx.Exec(ctx, `
		INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose,
		                            start_at, end_at, status, submitted_at)
		SELECT 'student-003', 'ค้างคืนส่งก่อน', 'student-003', id, 'overnight', 'activity',
		       now(), now() + interval '1 day', 'requested', now() - interval '2 hours'
		FROM equipment_item WHERE asset_code = 'AP-TEST-0002'`); err != nil {
		t.Fatal(err)
	}

	got, err := NewRepository(tx).ListPending(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, g := range got {
		if g.AssetCode == "AP-TEST-0002" && !g.IsFirst {
			t.Error("คำขอยืมปกติของ AP-TEST-0002 ต้องยังเป็นคำขอแรก เพราะคำขอค้างคืนไม่นับในลำดับ (Q-APV-23)")
		}
	}
}
