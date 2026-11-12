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
	svc := NewServiceWithClock(NewRepository(pool), testClock())

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
	svc := NewServiceWithClock(NewRepository(pool), testClock())

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
	svc := NewServiceWithClock(NewRepository(pool), testClock())

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
	svc := NewServiceWithClock(NewRepository(pool), testClock())
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
	svc := NewServiceWithClock(NewRepository(pool), testClock())
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

// testClock ตั้งเวลาไว้ที่ 09:00 น. ของวันนี้ ตามเขตเวลา Asia/Bangkok
// test ที่ไม่ได้ทดสอบกฎคืนเกินกำหนดต้องไม่ขึ้นกับเวลาจริงของเครื่อง
// เพราะข้อมูลตั้งต้นให้ student-001 ยืมของโดยมีกำหนดคืน 16:00 ของวันนี้ (T-03)
// ถ้าปล่อยให้ใช้นาฬิกาจริงแล้วรัน test ตอนเย็น ผู้ขอคนนั้นจะถูกกฎ FR-BRQ-06 ปฏิเสธ
func testClock() Clock { return FixedClock{T: bangkokToday(9, 0)} }

// AC-BRQ-06 (FR-BRQ-06)
// Given นักศึกษามีอุปกรณ์ที่ยังไม่คืน และเลยเวลา 16:00 น. ของกำหนดคืนไปแล้ว
// When  พยายามส่งคำขอยืมใหม่
// Then  ระบบต้องปฏิเสธคำขอ แจ้งสาเหตุ และไม่สร้างคำขอ
func TestAC_BRQ_06(t *testing.T) {
	pool := testPool(t)
	// Given: ข้อมูลตั้งต้นให้ student-002 ยืม CS-NB-0005 กำหนดคืน 16:00 ของเมื่อวาน และยังไม่คืน (T-03)
	const requester = "student-002"
	ctx := studentCtx(requester)
	svc := NewServiceWithClock(NewRepository(pool), testClock())

	var dueAt time.Time
	var returnedAt *time.Time
	if err := pool.QueryRow(ctx,
		`SELECT due_at, returned_at FROM loan WHERE borrower_id = $1`, requester).Scan(&dueAt, &returnedAt); err != nil {
		t.Fatalf("ข้อมูลตั้งต้นต้องมีการยืมของ %s: %v", requester, err)
	}
	if returnedAt != nil {
		t.Fatalf("ข้อมูลตั้งต้นต้องเป็นการยืมที่ยังไม่คืน แต่คืนแล้วเมื่อ %v", *returnedAt)
	}
	if !testClock().Now().After(overdueDeadline(dueAt, appLoc)) {
		t.Fatalf("ข้อมูลตั้งต้นต้องเลย 16:00 ของกำหนดคืนแล้ว (กำหนดคืน %v เวลาที่ใช้ทดสอบ %v)", dueAt, testClock().Now())
	}

	// When
	_, err := svc.Create(ctx, newTestRequest("CS-NB-0003"))

	// Then: ปฏิเสธพร้อมสาเหตุที่แยกจาก error ตัวอื่นได้ (ข้อความถึงผู้ขออยู่ที่ handler)
	if !errors.Is(err, ErrHasOverdueLoan) {
		t.Errorf("ผู้ขอที่มีของเกินกำหนดต้องได้ ErrHasOverdueLoan แต่ได้ %v", err)
	}
	// Then: ไม่สร้างคำขอ
	if n := countRequests(t, pool, requester); n != 0 {
		t.Errorf("ต้องไม่มีคำขอเกิดขึ้นเลย แต่พบ %d ใบ", n)
	}
}

// T-12 เสร็จเมื่อ (ส่วนเสริม): เส้นแบ่ง "เกินกำหนด" อยู่ที่ 16:00 น. ของวันกำหนดคืนพอดี (FR-BRQ-06 ข้อ ก)
// ใช้ข้อมูลตั้งต้นที่ student-001 ยืม CS-NB-0004 กำหนดคืน 16:00 ของวันนี้ และยังไม่คืน (T-03)
func TestOverdue_BoundaryAtReturnHour(t *testing.T) {
	pool := testPool(t)
	const requester = "student-001"
	ctx := studentCtx(requester)
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM borrow_request WHERE requester_id = $1`, requester)
	})

	// ก่อน 16:00 ของวันกำหนดคืน ยังไม่ถือว่าเกินกำหนด ส่งคำขอได้ตามปกติ
	before := NewServiceWithClock(NewRepository(pool), FixedClock{T: bangkokToday(15, 59)})
	if _, err := before.Create(ctx, newTestRequest("CS-NB-0001")); err != nil {
		t.Errorf("เวลา 15:59 ของวันกำหนดคืน ยังต้องส่งคำขอได้ แต่ได้ error: %v", err)
	}

	// เลย 16:00 แล้ว ถือว่าเกินกำหนดทันที
	after := NewServiceWithClock(NewRepository(pool), FixedClock{T: bangkokToday(16, 1)})
	if _, err := after.Create(ctx, newTestRequest("CS-NB-0003")); !errors.Is(err, ErrHasOverdueLoan) {
		t.Errorf("เวลา 16:01 ของวันกำหนดคืน ต้องได้ ErrHasOverdueLoan แต่ได้ %v", err)
	}
}

// T-12 เสร็จเมื่อ (ส่วนเสริม): ของค้างเกินกำหนดในภาควิชาอื่นไม่นับ (FR-BRQ-06 ข้อ ค, DOM-ACCESS-01)
func TestOverdue_OtherDepartmentNotCounted(t *testing.T) {
	pool := testPool(t)
	const requester = "student-013" // ผู้ขอของ test นี้เท่านั้น ไม่ให้ชนกับข้อมูลตั้งต้น
	ctx := studentCtx(requester)    // สังกัดภาควิชา 1
	svc := NewServiceWithClock(NewRepository(pool), testClock())

	// ตรวจข้อมูลตั้งต้นก่อนว่า XX-NB-0001 เป็นของภาควิชาอื่นจริง (ไม่ใช่ 1 ที่ผู้ขอสังกัด)
	const otherDeptItem = "XX-NB-0001"
	var itemDept int
	if err := pool.QueryRow(ctx, `
		SELECT m.department_id FROM equipment_item i
		JOIN equipment_model m ON m.id = i.model_id WHERE i.asset_code = $1`, otherDeptItem).Scan(&itemDept); err != nil {
		t.Fatal(err)
	}
	if itemDept == 1 {
		t.Fatalf("%s ต้องเป็นอุปกรณ์ของภาควิชาอื่น แต่อยู่ภาควิชาเดียวกับผู้ขอ", otherDeptItem)
	}

	// Given: มีของค้างเกินกำหนดของภาควิชาอื่น กำหนดคืน 16:00 ของเมื่อวาน
	if _, err := pool.Exec(ctx, `
		INSERT INTO loan (item_id, borrower_id, loan_type, due_at)
		SELECT id, $1, 'normal', $2 FROM equipment_item WHERE asset_code = $3`,
		requester, bangkokToday(16, 0).AddDate(0, 0, -1), otherDeptItem); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM loan WHERE borrower_id = $1`, requester)
		_, _ = pool.Exec(context.Background(), `DELETE FROM borrow_request WHERE requester_id = $1`, requester)
	})

	// When/Then: ส่งคำขออุปกรณ์ของภาควิชาตนเองได้ตามปกติ เพราะของค้างอยู่คนละภาควิชา
	if _, err := svc.Create(ctx, newTestRequest("CS-NB-0003")); err != nil {
		t.Errorf("ของค้างในภาควิชาอื่นต้องไม่นับ ควรส่งคำขอได้ แต่ได้ error: %v", err)
	}
}

// insertLoan เพิ่มการยืม 1 ใบในข้อมูลทดสอบ จำลองงานของ 005 (ส่งมอบ) และ 006 (รับคืน)
// 002 ไม่เคยเขียนตาราง loan เอง ที่นี่เขียนเพื่อสร้างสถานการณ์ของ AC เท่านั้น
func insertLoan(t *testing.T, pool *pgxpool.Pool, borrowerID, assetCode string, dueAt time.Time, returnedAt *time.Time) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), `
		INSERT INTO loan (item_id, borrower_id, loan_type, due_at, returned_at)
		SELECT id, $1, 'normal', $2, $3 FROM equipment_item WHERE asset_code = $4`,
		borrowerID, dueAt, returnedAt, assetCode); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM loan WHERE borrower_id = $1`, borrowerID)
		_, _ = pool.Exec(context.Background(), `DELETE FROM borrow_request WHERE requester_id = $1`, borrowerID)
	})
}

// AC-BRQ-12 (FR-BRQ-06 ข้อ ข)
// Given นักศึกษาคืนอุปกรณ์ที่เกินกำหนดแล้วเมื่อเวลา 14:00 น. ของวันนี้
// When  ส่งคำขอยืมใหม่เวลา 23:59 น. ของวันเดียวกัน แล้วส่งอีกครั้งเวลา 00:01 น. ของวันถัดไป
//
//	(เทียบด้วยเขตเวลา Asia/Bangkok)
//
// Then  คำขอเวลา 23:59 น. ต้องถูกปฏิเสธ ส่วนคำขอเวลา 00:01 น. ต้องส่งได้
func TestAC_BRQ_12(t *testing.T) {
	pool := testPool(t)
	const requester = "student-014" // ผู้ขอของ test นี้เท่านั้น
	ctx := studentCtx(requester)

	// Given: ยืมโดยมีกำหนดคืน 16:00 ของเมื่อวาน แล้วเพิ่งมาคืนเอา 14:00 ของวันนี้ (คืนช้า)
	returnedAt := bangkokToday(14, 0)
	insertLoan(t, pool, requester, "CS-NB-0001", bangkokToday(16, 0).AddDate(0, 0, -1), &returnedAt)

	// When/Then: 23:59 ของวันที่คืน ยังส่งคำขอไม่ได้
	sameDay := NewServiceWithClock(NewRepository(pool), FixedClock{T: bangkokToday(23, 59)})
	if _, err := sameDay.Create(ctx, newTestRequest("CS-NB-0003")); !errors.Is(err, ErrHasOverdueLoan) {
		t.Errorf("เวลา 23:59 ของวันที่คืนช้า ต้องได้ ErrHasOverdueLoan แต่ได้ %v", err)
	}
	if n := countRequests(t, pool, requester); n != 0 {
		t.Errorf("คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}

	// When/Then: 00:01 ของวันถัดไป ปลดล็อกแล้ว ส่งได้ตามปกติ
	nextDay := NewServiceWithClock(NewRepository(pool), FixedClock{T: bangkokToday(0, 1).AddDate(0, 0, 1)})
	if _, err := nextDay.Create(ctx, newTestRequest("CS-NB-0003")); err != nil {
		t.Fatalf("เวลา 00:01 ของวันถัดไป ต้องส่งคำขอได้ แต่ได้ error: %v", err)
	}
	if n := countRequests(t, pool, requester); n != 1 {
		t.Errorf("ต้องมีคำขอ 1 ใบหลังปลดล็อก แต่พบ %d ใบ", n)
	}
}

// T-13 เสร็จเมื่อ (ส่วนเสริม): คืนตรงเวลาไม่เข้ากฎนี้ (FR-BRQ-06 พูดถึงผู้ที่คืนเกินกำหนดเท่านั้น)
func TestOverdue_OnTimeReturnDoesNotBlock(t *testing.T) {
	pool := testPool(t)
	const requester = "student-015"
	ctx := studentCtx(requester)

	// Given: กำหนดคืน 16:00 ของวันนี้ และคืนไปแล้วตั้งแต่ 10:00 ของวันนี้ คือคืนตรงเวลา
	returnedAt := bangkokToday(10, 0)
	insertLoan(t, pool, requester, "CS-NB-0001", bangkokToday(16, 0), &returnedAt)

	// When/Then: ส่งคำขอเวลา 11:00 ของวันเดียวกันได้เลย ไม่ต้องรอข้ามวัน
	svc := NewServiceWithClock(NewRepository(pool), FixedClock{T: bangkokToday(11, 0)})
	if _, err := svc.Create(ctx, newTestRequest("CS-NB-0003")); err != nil {
		t.Errorf("คืนตรงเวลาต้องส่งคำขอได้ทันที แต่ได้ error: %v", err)
	}
}
