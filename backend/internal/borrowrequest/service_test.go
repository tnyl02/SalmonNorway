package borrowrequest

import (
	"context"
	"errors"
	"os"
	"slices"
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
	// ใช้ errors.Is เพราะตั้งแต่ T-09 error พกสถานะล่าสุดมาด้วย (ItemNotAvailableError)
	if _, err := svc.Create(ctx, newTestRequest("CS-NB-0004")); !errors.Is(err, ErrItemNotAvailable) {
		t.Errorf("อุปกรณ์ที่ถูกยืมอยู่ต้องได้ ErrItemNotAvailable แต่ได้ %v", err)
	}
	// รหัสครุภัณฑ์ที่ไม่มีอยู่จริง
	if _, err := svc.Create(ctx, newTestRequest("NO-SUCH-0001")); !errors.Is(err, ErrItemNotAvailable) {
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

	if _, err := svc.Create(context.Background(), newTestRequest("CS-NB-0001")); !errors.Is(err, ErrNoPrincipal) {
		t.Errorf("ไม่มีผู้ใช้ใน context ต้องได้ ErrNoPrincipal แต่ได้ %v", err)
	}
}

// สถานะที่ฟีเจอร์อื่น (003 ถึง 006) เป็นคนตั้ง แต่ test ของ AC-BRQ-11 ต้องจำลองขึ้นมาเอง
// ฟีเจอร์ 002 ตั้งได้เฉพาะ config.StatusRequested (FR-BRQ-04) ค่าที่เหลือมาจาก data-model.md v9
const (
	statusApproved  = "approved"
	statusCancelled = "cancelled"
	statusRejected  = "rejected"
)

// setStatus จำลองว่าฟีเจอร์อื่นเปลี่ยนสถานะคำขอใบนั้นแล้ว
// ยืนยันก่อนว่าค่าที่ใช้อยู่ในรายการสถานะที่ลงทะเบียนไว้ ไม่ใช่ค่าที่ test คิดขึ้นเอง (config.RequestStatuses)
func setStatus(t *testing.T, pool *pgxpool.Pool, id int, status string) {
	t.Helper()
	if !slices.Contains(config.RequestStatuses, status) {
		t.Fatalf("สถานะ %q ไม่อยู่ในรายการของ data-model.md (config.RequestStatuses)", status)
	}
	if _, err := pool.Exec(context.Background(),
		`UPDATE borrow_request SET status = $1 WHERE id = $2`, status, id); err != nil {
		t.Fatal(err)
	}
}

// countRequests นับคำขอของผู้ขอคนนั้น ใช้ยืนยันว่า "ไม่สร้างคำขอใหม่" เป็นจริง
func countRequests(t *testing.T, pool *pgxpool.Pool, requesterID string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*) FROM borrow_request WHERE requester_id = $1`, requesterID).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// AC-BRQ-11 (FR-BRQ-05 ข)
// Given นักศึกษาคนหนึ่งมีคำขอของอุปกรณ์ชิ้นหนึ่งอยู่แล้ว ไม่ว่าสถานะจะเป็น `requested` หรือ `approved`
// When  ส่งคำขอของอุปกรณ์ชิ้นเดิมซ้ำอีกครั้ง
// Then  ระบบต้องปฏิเสธ ไม่สร้างคำขอใหม่ และแจ้งว่ามีคำขอของอุปกรณ์ชิ้นนี้ค้างอยู่แล้ว
// และ
// Given คำขอเดิมถูกยกเลิกหรือถูกปฏิเสธไปแล้ว
// When  ส่งคำขอของอุปกรณ์ชิ้นเดิมใหม่
// Then  ระบบต้องรับคำขอได้ตามปกติ
func TestAC_BRQ_11(t *testing.T) {
	pool := testPool(t)
	// ใช้ผู้ขอของ test นี้เท่านั้น เพื่อไม่ให้คำขอของ test อื่นมาชนกฎกันซ้ำ
	const requester = "student-011"
	const assetCode = "CS-NB-0003" // อุปกรณ์ว่างของภาควิชา 1 ในข้อมูลตั้งต้น (T-03)
	ctx := studentCtx(requester)
	svc := NewService(NewRepository(pool))
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM borrow_request WHERE requester_id = $1`, requester)
	})

	// Given: มีคำขอใบแรกอยู่แล้ว สถานะ requested
	first, err := svc.Create(ctx, newTestRequest(assetCode))
	if err != nil {
		t.Fatalf("คำขอใบแรกต้องสำเร็จ แต่ได้ error: %v", err)
	}

	// When/Then: ส่งซ้ำขณะคำขอเดิมยังเป็น requested
	if _, err := svc.Create(ctx, newTestRequest(assetCode)); !errors.Is(err, ErrDuplicateRequest) {
		t.Errorf("ส่งซ้ำขณะสถานะ %s ต้องได้ ErrDuplicateRequest แต่ได้ %v", config.StatusRequested, err)
	}
	if n := countRequests(t, pool, requester); n != 1 {
		t.Errorf("ต้องไม่มีคำขอใหม่เกิดขึ้น ควรมี 1 ใบ แต่พบ %d ใบ", n)
	}

	// When/Then: ส่งซ้ำขณะคำขอเดิมถูกอนุมัติแล้ว (approved) ก็ยังต้องถูกปฏิเสธ
	setStatus(t, pool, first.ID, statusApproved)
	if _, err := svc.Create(ctx, newTestRequest(assetCode)); !errors.Is(err, ErrDuplicateRequest) {
		t.Errorf("ส่งซ้ำขณะสถานะ %s ต้องได้ ErrDuplicateRequest แต่ได้ %v", statusApproved, err)
	}
	if n := countRequests(t, pool, requester); n != 1 {
		t.Errorf("ต้องไม่มีคำขอใหม่เกิดขึ้น ควรมี 1 ใบ แต่พบ %d ใบ", n)
	}

	// Given: คำขอเดิมถูกยกเลิก -> When ส่งใหม่ -> Then รับได้ตามปกติ
	setStatus(t, pool, first.ID, statusCancelled)
	second, err := svc.Create(ctx, newTestRequest(assetCode))
	if err != nil {
		t.Fatalf("คำขอเดิมถูกยกเลิกแล้ว ต้องส่งใหม่ได้ แต่ได้ error: %v", err)
	}
	if n := countRequests(t, pool, requester); n != 2 {
		t.Errorf("ต้องมีคำขอ 2 ใบ (ใบที่ยกเลิกและใบใหม่) แต่พบ %d ใบ", n)
	}

	// Given: คำขอเดิมถูกปฏิเสธ -> When ส่งใหม่ -> Then รับได้ตามปกติ
	setStatus(t, pool, second.ID, statusRejected)
	if _, err := svc.Create(ctx, newTestRequest(assetCode)); err != nil {
		t.Fatalf("คำขอเดิมถูกปฏิเสธแล้ว ต้องส่งใหม่ได้ แต่ได้ error: %v", err)
	}
	if n := countRequests(t, pool, requester); n != 3 {
		t.Errorf("ต้องมีคำขอ 3 ใบ แต่พบ %d ใบ", n)
	}
}

// T-11 เสร็จเมื่อ (ส่วนเสริม): ยิงพร้อมกัน 2 ครั้งต้องเหลือคำขอแถวเดียว (FR-BRQ-05 ข)
// กันซ้ำที่ฐานข้อมูลด้วย unique index จาก T-01 ไม่ใช่การอ่านมาตรวจก่อนเขียน ซึ่งจะมีช่องให้ผ่านทั้งคู่
func TestDuplicateRequest_ConcurrentSubmitsKeepOneRow(t *testing.T) {
	pool := testPool(t)
	const requester = "student-012"
	const assetCode = "CS-NB-0002" // อุปกรณ์ว่างของภาควิชา 1 ในข้อมูลตั้งต้น (T-03)
	ctx := studentCtx(requester)
	svc := NewService(NewRepository(pool))
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM borrow_request WHERE requester_id = $1`, requester)
	})

	errs := make([]error, 2)
	var wg sync.WaitGroup
	wg.Add(len(errs))
	for i := range errs {
		go func(i int) {
			defer wg.Done()
			_, errs[i] = svc.Create(ctx, newTestRequest(assetCode))
		}(i)
	}
	wg.Wait()

	var created, duplicate int
	for _, err := range errs {
		switch {
		case err == nil:
			created++
		case errors.Is(err, ErrDuplicateRequest):
			duplicate++
		default:
			t.Fatalf("ได้ error ที่ไม่คาดไว้: %v", err)
		}
	}
	if created != 1 || duplicate != 1 {
		t.Errorf("ยิงพร้อมกัน 2 ครั้ง ต้องสำเร็จ 1 และซ้ำ 1 แต่ได้สำเร็จ %d ซ้ำ %d", created, duplicate)
	}
	if n := countRequests(t, pool, requester); n != 1 {
		t.Errorf("ต้องมีคำขอแถวเดียว แต่พบ %d แถว", n)
	}
}
