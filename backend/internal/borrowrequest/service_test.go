package borrowrequest

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/platform"
)

const (
	migrationsDir = "../../db/migrations"
	seedFile      = "../../db/seed/seed.sql"
)

var (
	dbOnce      sync.Once
	dbPool      *pgxpool.Pool
	dbErr       error
	dbContainer *tcpostgres.PostgresContainer
)

// TestMain ปิดฐานข้อมูลหลัง test ทั้ง package จบ
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

// testPool เปิด PostgreSQL 16 ใน container ครั้งเดียว แล้วรัน migration และข้อมูลตั้งต้น (plan.md ข้อ 2)
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

// studentCtx จำลองผู้ใช้ที่ยืนยันตัวตนแล้ว แบบเดียวกับที่ auth.Middleware ใส่ใน context (IF-IDP-01)
func studentCtx(userID string) context.Context {
	return auth.WithPrincipal(context.Background(),
		auth.Principal{UserID: userID, Role: auth.RoleStudent, DepartmentID: 1})
}

// newTestRequest คืนข้อมูลฟอร์มที่กรอกครบ สำหรับอุปกรณ์รหัสที่ระบุ
func newTestRequest(assetCode string) NewRequest {
	start := time.Now()
	return NewRequest{
		AssetCode:     assetCode,
		RequesterName: "ธัญลักษณ์ ทดสอบ",
		StudentCode:   "660710715",
		LoanType:      "normal",
		Purpose:       "study",
		StartAt:       start,
		EndAt:         start.Add(2 * time.Hour),
	}
}

// AC-BRQ-04 (FR-BRQ-04)
// Given นักศึกษาไม่มีรายการคืนเกินกำหนดค้างอยู่ และอุปกรณ์ชิ้นนั้นสถานะ "ว่าง"
// When  กรอกข้อมูลครบถ้วนและกดส่งคำขอ
// Then  ระบบต้องสร้างคำขอด้วยสถานะ "รออนุมัติ" พร้อมผู้ขอ อุปกรณ์รายชิ้น ประเภทการยืม
//
//	และเวลาที่ส่งตามเวลา Asia/Bangkok โดยสถานะของอุปกรณ์ยังเป็น "ว่าง"
func TestAC_BRQ_04(t *testing.T) {
	pool := testPool(t)
	ctx := studentCtx("student-001")
	svc := NewService(NewRepository(pool))

	// Given: CS-NB-0001 เป็นอุปกรณ์ว่างของภาควิชา 1 ในข้อมูลตั้งต้น (T-03)
	const assetCode = "CS-NB-0001"
	var before string
	if err := pool.QueryRow(ctx, `SELECT status FROM equipment_item WHERE asset_code = $1`, assetCode).Scan(&before); err != nil {
		t.Fatal(err)
	}
	if before != "available" {
		t.Fatalf("ข้อมูลตั้งต้นต้องให้ %s เป็น available แต่เป็น %q", assetCode, before)
	}

	// When
	req, err := svc.Create(ctx, newTestRequest(assetCode))
	if err != nil {
		t.Fatalf("ส่งคำขอต้องสำเร็จ แต่ได้ error: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM borrow_request WHERE id = $1`, req.ID)
	})

	// Then: สถานะเริ่มต้นเป็น requested เท่านั้น (FR-BRQ-04, data-model.md v9)
	if req.Status != config.StatusRequested {
		t.Errorf("สถานะเริ่มต้นต้องเป็น %q แต่ได้ %q", config.StatusRequested, req.Status)
	}
	if req.ID == 0 {
		t.Error("คำขอที่สร้างต้องมีรหัสคำขอ")
	}

	// Then: ผู้ขอที่บันทึกคือรหัสจากการยืนยันตัวตน ไม่ใช่รหัสนักศึกษาที่พิมพ์ในฟอร์ม (ASM-BRQ-06)
	var requesterID, studentCode, loanType, purpose string
	var purposeOther *string
	if err := pool.QueryRow(ctx, `
		SELECT requester_id, student_code, loan_type, purpose, purpose_other
		FROM borrow_request WHERE id = $1`, req.ID,
	).Scan(&requesterID, &studentCode, &loanType, &purpose, &purposeOther); err != nil {
		t.Fatal(err)
	}
	if requesterID != "student-001" {
		t.Errorf("ผู้ขอต้องเป็น student-001 จากการยืนยันตัวตน แต่ได้ %q", requesterID)
	}
	if studentCode != "660710715" {
		t.Errorf("รหัสนักศึกษาที่พิมพ์ในฟอร์มต้องถูกเก็บไว้ แต่ได้ %q", studentCode)
	}
	if loanType != "normal" || purpose != "study" {
		t.Errorf("ประเภทการยืมและวัตถุประสงค์ต้องถูกเก็บ แต่ได้ %q, %q", loanType, purpose)
	}
	if purposeOther != nil {
		t.Errorf("ไม่ได้เลือก \"อื่น ๆ\" ช่องพิมพ์ระบุต้องว่าง แต่ได้ %q", *purposeOther)
	}

	// Then: เวลาที่ส่งเทียบด้วยเขตเวลา Asia/Bangkok และใกล้เคียงเวลาปัจจุบัน (FR-BRQ-04, glossary)
	if req.SubmittedAt.Location().String() != config.TimeZone {
		t.Errorf("เวลาที่ส่งต้องอยู่ในเขตเวลา %s แต่ได้ %s", config.TimeZone, req.SubmittedAt.Location())
	}
	if d := time.Since(req.SubmittedAt); d < -time.Minute || d > time.Minute {
		t.Errorf("เวลาที่ส่งต้องใกล้เคียงเวลาปัจจุบัน แต่ต่างกัน %v", d)
	}

	// Then: ฟีเจอร์นี้ต้องไม่เปลี่ยนสถานะอุปกรณ์ (DOM-REQ-01, CON-DATA-01)
	var after string
	if err := pool.QueryRow(ctx, `SELECT status FROM equipment_item WHERE asset_code = $1`, assetCode).Scan(&after); err != nil {
		t.Fatal(err)
	}
	if after != "available" {
		t.Errorf("อุปกรณ์ที่มีคำขอค้างต้องยังเป็น available (DOM-REQ-01) แต่เป็น %q", after)
	}
}

// T-04 เสร็จเมื่อ (ส่วนเสริม): ส่งคำขอของอุปกรณ์ที่ไม่ใช่สถานะ "ว่าง" ต้องไม่เกิดคำขอ (CON-DATA-01)
// การแปลงเป็นรหัส error และการแจ้งสถานะล่าสุดให้ผู้ขอเป็นงานของ T-09
func TestCreate_RejectsItemNotAvailable(t *testing.T) {
	pool := testPool(t)
	ctx := studentCtx("student-001")
	svc := NewService(NewRepository(pool))

	// CS-NB-0004 ถูกยืมอยู่ในข้อมูลตั้งต้น
	if _, err := svc.Create(ctx, newTestRequest("CS-NB-0004")); err != ErrItemNotAvailable {
		t.Errorf("อุปกรณ์ที่ถูกยืมอยู่ต้องได้ ErrItemNotAvailable แต่ได้ %v", err)
	}
	// รหัสครุภัณฑ์ที่ไม่มีอยู่จริง
	if _, err := svc.Create(ctx, newTestRequest("NO-SUCH-0001")); err != ErrItemNotAvailable {
		t.Errorf("รหัสครุภัณฑ์ที่ไม่มีอยู่ต้องได้ ErrItemNotAvailable แต่ได้ %v", err)
	}

	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM borrow_request`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("ต้องไม่มีคำขอเกิดขึ้นเลย แต่พบ %d ใบ", count)
	}
}

// T-04 เสร็จเมื่อ (ส่วนเสริม): เรียก service โดยไม่มีผู้ใช้ที่ยืนยันตัวตนต้องถูกปฏิเสธ ไม่เดาว่าเป็นใคร (IF-IDP-01)
func TestCreate_RequiresPrincipal(t *testing.T) {
	pool := testPool(t)
	svc := NewService(NewRepository(pool))

	if _, err := svc.Create(context.Background(), newTestRequest("CS-NB-0001")); err != ErrNoPrincipal {
		t.Errorf("ไม่มีผู้ใช้ใน context ต้องได้ ErrNoPrincipal แต่ได้ %v", err)
	}
}
