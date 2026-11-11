package returns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/platform"
)

const (
	migrationsDir = "../../db/migrations"
	seedFile      = "../../db/seed/seed.sql"
	seed006File   = "../../db/seed/seed_006_return.sql"
)

var (
	dbOnce      sync.Once
	dbURL       string
	dbErr       error
	dbContainer *tcpostgres.PostgresContainer
)

// TestMain ปิด container หลัง test ทั้ง package จบ (ถ้ามีการเปิด)
func TestMain(m *testing.M) {
	code := m.Run()
	if dbContainer != nil {
		_ = testcontainers.TerminateContainer(dbContainer)
	}
	os.Exit(code)
}

// testPool เปิด PostgreSQL 16 ใน container ครั้งเดียว รัน migration (รวม 0009) และข้อมูลตัวอย่างของ 001 และ 006
// (plan.md ข้อ 6: integration ด้วย testcontainers + migration จริง)
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
		if dbURL, dbErr = dbContainer.ConnectionString(ctx, "sslmode=disable"); dbErr != nil {
			return
		}
		if dbErr = platform.Migrate(dbURL, migrationsDir); dbErr != nil {
			return
		}
		if dbErr = platform.Seed(ctx, dbURL, seedFile); dbErr != nil {
			return
		}
		dbErr = platform.Seed(ctx, dbURL, seed006File)
	})
	if dbErr != nil {
		t.Fatalf("เตรียมฐานข้อมูลไม่สำเร็จ: %v", dbErr)
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	return pool
}

// newTestRouter สร้าง router เหมือน main.go: /healthz + endpoint ของ 006 ผ่าน auth จำลองของ 001
func newTestRouter(pool *pgxpool.Pool) http.Handler {
	return newRouterWith(pool, auth.NewMockProvider())
}

// newRouterWith ใช้ router กลางของ 001 (platform.NewRouter) แล้วผูก route ของ 006 เข้าไป แบบเดียวกับ main.go
func newRouterWith(pool *pgxpool.Pool, provider auth.Provider) *http.ServeMux {
	mux := platform.NewRouter(pool)
	NewHandler(NewRepository(pool)).Routes(provider)(mux)
	return mux
}

func call(router http.Handler, method, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// T-03 เสร็จเมื่อ: ทุก endpoint ของ 006 ผ่าน auth ก่อน ไม่มี token ได้ 401 (IF-IDP-01)
// (เดิมตรวจว่า open-loans ตอบ 501 ตอนยังไม่ทำ ตอนนี้ T-14 ทำแล้ว จึงเหลือเฉพาะส่วน auth)
func TestReturns_Skeleton_RoutedThroughAuth(t *testing.T) {
	router := newTestRouter(nil) // ไม่ต้องใช้ฐานข้อมูล: middleware ปฏิเสธก่อนถึง handler

	if rec := call(router, http.MethodGet, "/api/v1/returns/open-loans", ""); rec.Code != http.StatusUnauthorized {
		t.Errorf("ไม่มี token: status = %d, want 401 (IF-IDP-01)", rec.Code)
	}
	// ทุก endpoint ของ 006 ต้องผ่าน auth
	for _, ep := range []struct{ method, path string }{
		{http.MethodGet, "/api/v1/returns/lookup?assetCode=CS-NB-0004"},
		{http.MethodPost, "/api/v1/returns"},
		{http.MethodGet, "/api/v1/returns/loans/1"},
	} {
		if rec := call(router, ep.method, ep.path, ""); rec.Code != http.StatusUnauthorized {
			t.Errorf("%s %s ไม่มี token: status = %d, want 401", ep.method, ep.path, rec.Code)
		}
	}
}

// T-03 เสร็จเมื่อ (ข้อมูลตัวอย่าง): seed_006_return.sql ใส่ได้หลัง seed.sql และ migration 0009
// ภาควิชา 1 มีการยืมที่ยังไม่คืน 2 รายการ และคืนแล้ว 2 รายการ / ภาควิชา 2 มีการยืมที่ยังไม่คืน 1 รายการ
func TestReturns_Seed006Loaded(t *testing.T) {
	pool := testPool(t)

	rows, err := pool.Query(context.Background(), `
		SELECT m.department_id,
		       count(*) FILTER (WHERE l.returned_at IS NULL),
		       count(*) FILTER (WHERE l.returned_at IS NOT NULL)
		FROM loan l
		JOIN equipment_item i ON i.id = l.item_id
		JOIN equipment_model m ON m.id = i.model_id
		WHERE i.asset_code LIKE 'CS-%' OR i.asset_code LIKE 'XX-%' -- เฉพาะข้อมูลตัวอย่าง ไม่นับที่ test อื่นสร้าง
		GROUP BY m.department_id ORDER BY m.department_id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	got := map[int][2]int{}
	for rows.Next() {
		var dept, open, returned int
		if err := rows.Scan(&dept, &open, &returned); err != nil {
			t.Fatal(err)
		}
		got[dept] = [2]int{open, returned}
	}
	want := map[int][2]int{1: {2, 2}, 2: {1, 0}}
	for dept, w := range want {
		if got[dept] != w {
			t.Errorf("ภาควิชา %d: ยังไม่คืน/คืนแล้ว = %v, want %v", dept, got[dept], w)
		}
	}
}

// ---- ตัวช่วยสำหรับ test ของการบันทึกรับคืน ----

// insertBorrowedItem สร้างอุปกรณ์ "ถูกยืม" ของภาควิชา 1 (รุ่น 1) พร้อมการยืมที่ยังไม่คืน กำหนดคืนตาม dueIn นับจากตอนนี้
// ใช้รหัสครุภัณฑ์ใหม่ทุก test ข้อมูลจึงไม่ชนกันข้าม test (CON-DATA-03)
func insertBorrowedItem(t *testing.T, pool *pgxpool.Pool, assetCode string, dueIn time.Duration) {
	t.Helper()
	ctx := context.Background()
	var itemID int
	if err := pool.QueryRow(ctx,
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, 'borrowed') RETURNING id`,
		assetCode).Scan(&itemID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx,
		`INSERT INTO loan (item_id, borrower_id, loan_type, due_at) VALUES ($1, 'student-001', 'normal', now() + $2::interval)`,
		itemID, dueIn.String()); err != nil {
		t.Fatal(err)
	}
}

// postReturn ส่ง multipart ไปที่ POST /api/v1/returns ด้วย token ที่ระบุ
func postReturn(t *testing.T, router http.Handler, token string, fields map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	for k, v := range fields {
		if err := mw.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
	}
	if err := mw.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/returns", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// AC-RET-01 (FR-RET-03)
//
//	Given อุปกรณ์ชิ้นหนึ่งสถานะ "ถูกยืม" มีการยืมที่ยังไม่คืน
//	When  เจ้าหน้าที่บันทึกการรับคืน สภาพ "ปกติ"
//	Then  การยืมนั้นมีเวลาคืน และสถานะอุปกรณ์เป็น "ว่าง"
func TestAC_RET_01_NormalReturnBecomesAvailable(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)
	ctx := context.Background()

	// Given
	insertBorrowedItem(t, pool, "T04-AC01", 2*time.Hour)

	// When
	rec := postReturn(t, router, "mock-staff-001",
		map[string]string{"assetCode": "T04-AC01", "condition": "normal"})

	// Then
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 body=%s", rec.Code, rec.Body.String())
	}
	var resp struct {
		LoanID     int        `json:"loanId"`
		ReturnedAt *time.Time `json:"returnedAt"`
		Condition  string     `json:"condition"`
		ItemStatus string     `json:"itemStatus"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ReturnedAt == nil || resp.Condition != "normal" || resp.ItemStatus != "available" {
		t.Errorf("response = %+v, want returnedAt มีค่า condition=normal itemStatus=available", resp)
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store (NFR-DATA-01)", got)
	}

	var returnedAt *time.Time
	var status string
	if err := pool.QueryRow(ctx, `
		SELECT l.returned_at, i.status
		FROM loan l JOIN equipment_item i ON i.id = l.item_id
		WHERE i.asset_code = 'T04-AC01'`).Scan(&returnedAt, &status); err != nil {
		t.Fatal(err)
	}
	if returnedAt == nil {
		t.Error("การยืมต้องมีเวลาคืน (FR-RET-02)")
	}
	if status != "available" {
		t.Errorf("สถานะ = %q, want available (FR-RET-03)", status)
	}
}

// T-05 เสร็จเมื่อ: lookup คืนการยืมที่ยังไม่คืนของรหัสครุภัณฑ์ (FR-RET-01(ก)) ใช้ในขั้นหาการยืมของ AC-RET-01
func TestReturns_Lookup_FoundOpenLoan(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)

	// Given อุปกรณ์ "ถูกยืม" 2 ชิ้น: ชิ้นแรกยังไม่เลยกำหนด ชิ้นที่สองเลยกำหนดแล้ว
	insertBorrowedItem(t, pool, "T05-LOOKUP", 2*time.Hour)
	insertBorrowedItem(t, pool, "T05-LATE", -1*time.Hour)

	for _, tc := range []struct {
		assetCode string
		overdue   bool
	}{{"T05-LOOKUP", false}, {"T05-LATE", true}} {
		// When เจ้าหน้าที่กรอกรหัสครุภัณฑ์
		rec := call(router, http.MethodGet, "/api/v1/returns/lookup?assetCode="+tc.assetCode, "mock-staff-001")

		// Then ได้การยืมที่ยังไม่คืนพร้อมข้อมูลตาม plan.md ข้อ 4
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200 body=%s", tc.assetCode, rec.Code, rec.Body.String())
		}
		var got OpenLoan
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got.LoanID == 0 || got.AssetCode != tc.assetCode || got.ModelName == "" ||
			got.BorrowerID != "student-001" || got.LoanType != "normal" || got.DueAt.IsZero() {
			t.Errorf("%s: response = %+v ข้อมูลไม่ครบ", tc.assetCode, got)
		}
		if got.IsOverdue != tc.overdue {
			t.Errorf("%s: isOverdue = %v, want %v", tc.assetCode, got.IsOverdue, tc.overdue)
		}
		if rec.Header().Get("Cache-Control") != "no-store" {
			t.Errorf("%s: ต้องมี Cache-Control: no-store", tc.assetCode)
		}
	}

	// การยืมที่คืนแล้วต้องหาไม่พบ (CS-NB-0001 คืนแล้วใน seed_006_return.sql) กรณีอื่นของ AC-RET-03 อยู่ใน T-18
	if rec := call(router, http.MethodGet, "/api/v1/returns/lookup?assetCode=CS-NB-0001", "mock-staff-001"); rec.Code != http.StatusConflict {
		t.Errorf("อุปกรณ์ที่ไม่มีการยืมค้าง: status = %d, want 409", rec.Code)
	}
}

// AC-RET-09 (FR-RET-02 เวลาคืน)
//
//	Given เวลาของระบบคือเวลาปัจจุบันของฐานข้อมูล (Asia/Bangkok ใช้ตอนแสดงผล เวลาเก็บเป็น timestamptz)
//	When  เจ้าหน้าที่บันทึกการรับคืนสำเร็จ โดยส่งเวลาอื่นมาด้วย (2569-10-10 09:00 = 2026-10-10 09:00 +07:00)
//	Then  เวลาคืนในการยืมเป็นเวลาของระบบ ไม่ใช่เวลาที่ส่งมา
func TestAC_RET_09_ReturnTimeFromServer(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)
	ctx := context.Background()

	// Given
	insertBorrowedItem(t, pool, "T07-AC09", 2*time.Hour)
	var before time.Time
	if err := pool.QueryRow(ctx, `SELECT now()`).Scan(&before); err != nil {
		t.Fatal(err)
	}

	// When: ส่งเวลาคืนปลอมมาหลายชื่อฟิลด์ ระบบต้องไม่อ่านเลย
	const fake = "2026-10-10T09:00:00+07:00"
	rec := postReturn(t, router, "mock-staff-001", map[string]string{
		"assetCode":   "T07-AC09",
		"condition":   "normal",
		"returnedAt":  fake,
		"returned_at": fake,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 body=%s", rec.Code, rec.Body.String())
	}
	var after time.Time
	if err := pool.QueryRow(ctx, `SELECT now()`).Scan(&after); err != nil {
		t.Fatal(err)
	}

	// Then
	var returnedAt time.Time
	if err := pool.QueryRow(ctx, `
		SELECT l.returned_at FROM loan l JOIN equipment_item i ON i.id = l.item_id
		WHERE i.asset_code = 'T07-AC09'`).Scan(&returnedAt); err != nil {
		t.Fatal(err)
	}
	fakeTime, _ := time.Parse(time.RFC3339, fake)
	if returnedAt.Equal(fakeTime) {
		t.Fatal("เวลาคืนเป็นเวลาที่ client ส่งมา ต้องเป็นเวลาของระบบ (FR-RET-02, ASM-04)")
	}
	if returnedAt.Before(before) || returnedAt.After(after) {
		t.Errorf("returned_at = %v ต้องอยู่ระหว่างเวลาระบบก่อนและหลังบันทึก [%v, %v]", returnedAt, before, after)
	}
}

// AC-RET-02 (FR-RET-04)
//
//	Given อุปกรณ์ชิ้นหนึ่งสถานะ "ถูกยืม" มีการยืมที่ยังไม่คืน
//	When  เจ้าหน้าที่บันทึกการรับคืน สภาพ "ชำรุด"
//	Then  การยืมนั้นมีเวลาคืน และสถานะอุปกรณ์เป็น "ไม่พร้อมใช้งาน" ไม่เคยเป็น "ว่าง" ระหว่างทาง
//
// "ไม่เคยเป็นว่างระหว่างทาง" ตรวจด้วย trigger ที่สร้างเฉพาะใน test เก็บทุกค่าสถานะที่ถูกเขียน (plan.md ข้อ 6)
func TestAC_RET_02_DamagedReturnBecomesUnavailable(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)
	ctx := context.Background()

	// เตรียมตัวเก็บประวัติสถานะ (เฉพาะ test) แล้วลบทิ้งเมื่อจบ
	if _, err := pool.Exec(ctx, `
		CREATE TABLE test_status_history (asset_code TEXT, status TEXT);
		CREATE FUNCTION test_record_status() RETURNS trigger AS $$
		BEGIN
			INSERT INTO test_status_history VALUES (NEW.asset_code, NEW.status);
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
		CREATE TRIGGER test_status_history AFTER UPDATE OF status ON equipment_item
			FOR EACH ROW EXECUTE FUNCTION test_record_status();`); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `
			DROP TRIGGER test_status_history ON equipment_item;
			DROP FUNCTION test_record_status();
			DROP TABLE test_status_history;`)
	})

	// Given
	insertBorrowedItem(t, pool, "T08-AC02", 2*time.Hour)

	// When
	rec := postReturn(t, router, "mock-staff-001",
		map[string]string{"assetCode": "T08-AC02", "condition": "damaged"})

	// Then
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 body=%s", rec.Code, rec.Body.String())
	}
	var resp struct {
		ItemStatus string `json:"itemStatus"`
		Condition  string `json:"condition"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ItemStatus != "unavailable" || resp.Condition != "damaged" {
		t.Errorf("response = %+v, want condition=damaged itemStatus=unavailable", resp)
	}

	var returnedAt *time.Time
	var status string
	if err := pool.QueryRow(ctx, `
		SELECT l.returned_at, i.status FROM loan l JOIN equipment_item i ON i.id = l.item_id
		WHERE i.asset_code = 'T08-AC02'`).Scan(&returnedAt, &status); err != nil {
		t.Fatal(err)
	}
	if returnedAt == nil {
		t.Error("การยืมต้องมีเวลาคืน")
	}
	if status != "unavailable" {
		t.Errorf("สถานะ = %q, want unavailable (FR-RET-04)", status)
	}

	rows, err := pool.Query(ctx, `SELECT status FROM test_status_history WHERE asset_code = 'T08-AC02'`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var history []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			t.Fatal(err)
		}
		history = append(history, s)
	}
	if len(history) != 1 || history[0] != "unavailable" {
		t.Errorf("ประวัติสถานะ = %v, want [unavailable] (เปลี่ยนครั้งเดียว ไม่ผ่าน available)", history)
	}
}

// AC-RET-04 (FR-RET-02 สภาพ)
//
//	Given เจ้าหน้าที่กำลังบันทึกการรับคืน
//	When  ไม่ได้ระบุสภาพ แล้วกดบันทึก
//	Then  ระบบไม่บันทึก แจ้งให้เลือกสภาพ ปกติ หรือ ชำรุด และไม่มีตัวเลือก "สูญหาย"
//
// ฝั่งหลังบ้าน: ไม่มีสภาพ / ค่าว่าง / "lost" / ค่าอื่น -> 422 CONDITION_REQUIRED และการยืมกับสถานะไม่เปลี่ยน
// ฝั่งหน้าจอ (ตัวเลือกมีแค่ 2 ค่า และข้อความแจ้ง) อยู่ใน T-10
func TestAC_RET_04_ConditionRequired(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)
	ctx := context.Background()

	// Given
	insertBorrowedItem(t, pool, "T09-AC04", 2*time.Hour)

	cases := map[string]map[string]string{
		"ไม่ส่งฟิลด์สภาพ":               {"assetCode": "T09-AC04"},
		"สภาพว่าง":                      {"assetCode": "T09-AC04", "condition": ""},
		"สูญหาย (ไม่ใช่สภาพ FR-RET-02)": {"assetCode": "T09-AC04", "condition": "lost"},
		"ค่าที่ไม่รู้จัก":               {"assetCode": "T09-AC04", "condition": "broken"},
		"ภาษาไทยแทนค่าในระบบ":           {"assetCode": "T09-AC04", "condition": "ปกติ"},
	}
	for name, fields := range cases {
		// When
		rec := postReturn(t, router, "mock-staff-001", fields)

		// Then
		if rec.Code != http.StatusUnprocessableEntity {
			t.Errorf("%s: status = %d, want 422", name, rec.Code)
			continue
		}
		var body struct {
			Code string `json:"code"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Code != "CONDITION_REQUIRED" {
			t.Errorf("%s: code = %q, want CONDITION_REQUIRED", name, body.Code)
		}
	}

	var returnedAt *time.Time
	var status string
	if err := pool.QueryRow(ctx, `
		SELECT l.returned_at, i.status
		FROM loan l JOIN equipment_item i ON i.id = l.item_id
		WHERE i.asset_code = 'T09-AC04'`).Scan(&returnedAt, &status); err != nil {
		t.Fatal(err)
	}
	if returnedAt != nil || status != "borrowed" {
		t.Errorf("ต้องไม่บันทึกอะไรเลย: returned_at=%v status=%q", returnedAt, status)
	}
}

// T-11 เสร็จเมื่อ: การยืมที่คืนแล้วได้ {loanId, assetCode, borrowerId, dueAt, returnedAt, condition, returnedLate} (FR-RET-06)
// เป็นงานพื้นฐานของ T-16 (AC-RET-13)
func TestReturns_LoanDetail(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)

	// Given การยืมที่คืนแล้วสภาพชำรุด (บันทึกผ่าน POST จริง)
	insertBorrowedItem(t, pool, "T11-DETAIL", 2*time.Hour)
	rec := postReturn(t, router, "mock-staff-001", map[string]string{"assetCode": "T11-DETAIL", "condition": "damaged"})
	if rec.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: status = %d body=%s", rec.Code, rec.Body.String())
	}
	var created struct {
		LoanID int `json:"loanId"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	// When เปิดรายละเอียดการยืม
	rec = call(router, http.MethodGet, fmt.Sprintf("/api/v1/returns/loans/%d", created.LoanID), "mock-staff-001")

	// Then ได้ข้อมูลครบตาม plan.md ข้อ 4
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 body=%s", rec.Code, rec.Body.String())
	}
	var got LoanDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.LoanID != created.LoanID || got.AssetCode != "T11-DETAIL" || got.BorrowerID != "student-001" || got.DueAt.IsZero() {
		t.Errorf("response = %+v ข้อมูลการยืมไม่ครบ", got)
	}
	if got.ReturnedAt == nil || got.Condition == nil || *got.Condition != "damaged" || got.ReturnedLate == nil || *got.ReturnedLate {
		t.Errorf("returnedAt=%v condition=%v returnedLate=%v, want มีเวลาคืน, damaged, false", got.ReturnedAt, got.Condition, got.ReturnedLate)
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Error("ต้องมี Cache-Control: no-store")
	}

	// การยืมที่ยังไม่คืน: ข้อมูลตอนคืนเป็น null
	insertBorrowedItem(t, pool, "T11-OPEN", 2*time.Hour)
	var openID int
	if err := pool.QueryRow(context.Background(),
		`SELECT l.id FROM loan l JOIN equipment_item i ON i.id = l.item_id WHERE i.asset_code = 'T11-OPEN'`).Scan(&openID); err != nil {
		t.Fatal(err)
	}
	rec = call(router, http.MethodGet, fmt.Sprintf("/api/v1/returns/loans/%d", openID), "mock-staff-001")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `"returnedAt":null`) {
		t.Errorf("การยืมที่ยังไม่คืน: status = %d body=%s, want 200 และ returnedAt null", rec.Code, rec.Body.String())
	}

	// ไม่มีการยืมนี้ / รหัสไม่ใช่ตัวเลข -> 404
	for _, path := range []string{"/api/v1/returns/loans/999999", "/api/v1/returns/loans/abc"} {
		if rec := call(router, http.MethodGet, path, "mock-staff-001"); rec.Code != http.StatusNotFound {
			t.Errorf("%s: status = %d, want 404", path, rec.Code)
		}
	}
}

// AC-RET-06 (FR-RET-06, DOM-LOAN-01)
//
//	Given การยืมแบบปกติ กำหนดคืน 16:00
//	When  ระบบบันทึกการรับคืนเวลา 17:00 และอีกกรณีบันทึกเวลา 15:00
//	Then  กรณี 17:00 การยืมถูกบันทึกว่าคืนเกินกำหนด กรณี 15:00 ไม่ถูกบันทึก
//
// เวลาคืนมาจาก now() ของฐานข้อมูล (ASM-04) จึงแทน 17:00 / 15:00 ด้วยกำหนดคืนที่ผ่านไปแล้ว 1 ชั่วโมง / อีก 1 ชั่วโมง (plan.md ข้อ 6)
// กำหนดคืนอ่านจาก due_at ที่ 005 บันทึก ไม่คำนวณ 16:00 เอง (ASM-01)
func TestAC_RET_06_LateReturnFlagged(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)
	ctx := context.Background()

	cases := []struct {
		assetCode string
		dueIn     time.Duration
		wantLate  bool
	}{
		{"T13-LATE", -1 * time.Hour, true},   // คืน "17:00" หลังกำหนด 16:00
		{"T13-ONTIME", 1 * time.Hour, false}, // คืน "15:00" ก่อนกำหนด 16:00
	}
	for _, tc := range cases {
		// Given
		insertBorrowedItem(t, pool, tc.assetCode, tc.dueIn)

		// When
		rec := postReturn(t, router, "mock-staff-001", map[string]string{"assetCode": tc.assetCode, "condition": "normal"})
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d body=%s", tc.assetCode, rec.Code, rec.Body.String())
		}

		// Then: ทั้งใน response และในฐานข้อมูล
		var resp struct {
			ReturnedLate bool `json:"returnedLate"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		var stored *bool
		if err := pool.QueryRow(ctx, `
			SELECT l.returned_late FROM loan l JOIN equipment_item i ON i.id = l.item_id
			WHERE i.asset_code = $1`, tc.assetCode).Scan(&stored); err != nil {
			t.Fatal(err)
		}
		if resp.ReturnedLate != tc.wantLate || stored == nil || *stored != tc.wantLate {
			t.Errorf("%s: returnedLate response=%v stored=%v, want %v", tc.assetCode, resp.ReturnedLate, stored, tc.wantLate)
		}
	}

	// ค่าที่บันทึกแล้วไม่เปลี่ยนตามกำหนดคืนภายหลัง (FR-RET-06 "บันทึก" ไม่ใช่คำนวณใหม่)
	if _, err := pool.Exec(ctx, `
		UPDATE loan SET due_at = now() + interval '1 day'
		WHERE item_id = (SELECT id FROM equipment_item WHERE asset_code = 'T13-LATE')`); err != nil {
		t.Fatal(err)
	}
	var stillLate bool
	if err := pool.QueryRow(ctx, `
		SELECT l.returned_late FROM loan l JOIN equipment_item i ON i.id = l.item_id
		WHERE i.asset_code = 'T13-LATE'`).Scan(&stillLate); err != nil {
		t.Fatal(err)
	}
	if !stillLate {
		t.Error("returned_late ต้องคงค่าที่บันทึกตอนรับคืน")
	}
}

// fixedProvider คือระบบยืนยันตัวตนจำลองที่ test กำหนดผู้ใช้เอง (ใช้ auth.Provider ของ IF-IDP-01 ไม่แก้ mock ของกลาง)
type fixedProvider map[string]auth.Principal

func (f fixedProvider) Authenticate(_ context.Context, token string) (auth.Principal, error) {
	p, ok := f[token]
	if !ok {
		return auth.Principal{}, auth.ErrInvalidToken
	}
	return p, nil
}

// insertDepartment สร้างภาควิชาใหม่ 1 ภาควิชา พร้อมรุ่น 1 รุ่น คืนรหัสภาควิชาและรหัสรุ่น
func insertDepartment(t *testing.T, pool *pgxpool.Pool, name string) (deptID, modelID int) {
	t.Helper()
	ctx := context.Background()
	if err := pool.QueryRow(ctx, `INSERT INTO department (name) VALUES ($1) RETURNING id`, name).Scan(&deptID); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx,
		`INSERT INTO equipment_model (name, category_id, department_id) VALUES ($1, 1, $2) RETURNING id`,
		name+" รุ่นทดสอบ", deptID).Scan(&modelID); err != nil {
		t.Fatal(err)
	}
	return deptID, modelID
}

// insertLoanOfModel สร้างอุปกรณ์ของรุ่น modelID พร้อมการยืม ถ้า returned = true การยืมนั้นคืนแล้วและอุปกรณ์ "ว่าง"
func insertLoanOfModel(t *testing.T, pool *pgxpool.Pool, modelID int, assetCode string, dueIn time.Duration, returned bool) {
	t.Helper()
	ctx := context.Background()
	status := "borrowed"
	if returned {
		status = "available"
	}
	var itemID int
	if err := pool.QueryRow(ctx,
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, $2, $3) RETURNING id`,
		assetCode, modelID, status).Scan(&itemID); err != nil {
		t.Fatal(err)
	}
	if returned {
		_, err := pool.Exec(ctx, `INSERT INTO loan (item_id, borrower_id, loan_type, due_at, returned_at, return_condition, returned_late)
			VALUES ($1, 'student-201', 'normal', now() + $2::interval, now(), 'normal', false)`, itemID, dueIn.String())
		if err != nil {
			t.Fatal(err)
		}
		return
	}
	if _, err := pool.Exec(ctx, `INSERT INTO loan (item_id, borrower_id, loan_type, due_at)
		VALUES ($1, 'student-201', 'normal', now() + $2::interval)`, itemID, dueIn.String()); err != nil {
		t.Fatal(err)
	}
}

// AC-RET-10 (FR-RET-01 หน้ารายการ)
//
//	Given ภาควิชาของเจ้าหน้าที่มีการยืมที่ยังไม่คืน 3 รายการ คืนแล้ว 2 รายการ และภาควิชาอื่นมีการยืมที่ยังไม่คืน 1 รายการ
//	When  เจ้าหน้าที่เปิดหน้ารายการที่ยังไม่คืน
//	Then  เห็นเฉพาะ 3 รายการของภาควิชาตนเอง (ส่วนเลือกรายการเพื่อบันทึกคืนอยู่ในหน้าจอ T-15)
//
// ฐานข้อมูลทดสอบใช้ร่วมกันทุก test จึงสร้างภาควิชาใหม่ 2 ภาควิชาเฉพาะ test นี้ ให้นับได้ตรงตาม AC
func TestAC_RET_10_OpenLoansOwnDepartment(t *testing.T) {
	pool := testPool(t)

	// Given
	ownDept, ownModel := insertDepartment(t, pool, "T14 ภาควิชาของเจ้าหน้าที่")
	_, otherModel := insertDepartment(t, pool, "T14 ภาควิชาอื่น")
	insertLoanOfModel(t, pool, ownModel, "T14-OPEN-3", 3*time.Hour, false)
	insertLoanOfModel(t, pool, ownModel, "T14-OPEN-1", -1*time.Hour, false) // เลยกำหนดแล้ว
	insertLoanOfModel(t, pool, ownModel, "T14-OPEN-2", 1*time.Hour, false)
	insertLoanOfModel(t, pool, ownModel, "T14-DONE-1", -2*time.Hour, true)
	insertLoanOfModel(t, pool, ownModel, "T14-DONE-2", -3*time.Hour, true)
	insertLoanOfModel(t, pool, otherModel, "T14-OTHER", 1*time.Hour, false)

	provider := fixedProvider{"staff-t14": {UserID: "staff-t14", Role: auth.RoleStaff, DepartmentID: ownDept}}
	router := newRouterWith(pool, provider)

	// When
	rec := call(router, http.MethodGet, "/api/v1/returns/open-loans", "staff-t14")

	// Then
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 body=%s", rec.Code, rec.Body.String())
	}
	var got []OpenLoan
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	var codes []string
	for _, l := range got {
		codes = append(codes, l.AssetCode)
	}
	want := []string{"T14-OPEN-1", "T14-OPEN-2", "T14-OPEN-3"} // เรียงตามกำหนดคืน
	if strings.Join(codes, ",") != strings.Join(want, ",") {
		t.Errorf("รายการ = %v, want %v (เฉพาะที่ยังไม่คืนของภาควิชาตนเอง เรียงตามกำหนดคืน)", codes, want)
	}
	if len(got) == 3 && (!got[0].IsOverdue || got[1].IsOverdue) {
		t.Errorf("isOverdue = %v, %v want true, false", got[0].IsOverdue, got[1].IsOverdue)
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Error("ต้องมี Cache-Control: no-store")
	}

	// ภาควิชาที่ไม่มีการยืมค้างได้รายการว่าง [] ไม่ใช่ null
	emptyDept, _ := insertDepartment(t, pool, "T14 ภาควิชาว่าง")
	provider["staff-empty"] = auth.Principal{UserID: "staff-empty", Role: auth.RoleStaff, DepartmentID: emptyDept}
	if rec := call(router, http.MethodGet, "/api/v1/returns/open-loans", "staff-empty"); strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Errorf("ภาควิชาที่ไม่มีการยืมค้าง: body = %s, want []", rec.Body.String())
	}
}

// AC-RET-13 (FR-RET-06 การแสดงผล) ฝั่งหลังบ้าน (ป้ายบนหน้าจออยู่ใน T-17)
//
//	Given การยืมรายการหนึ่งคืนเกินกำหนดแล้ว และอีกรายการยังไม่คืนแต่เลยกำหนดคืนแล้ว
//	When  เจ้าหน้าที่เปิดหน้ารายละเอียดการยืมรายการแรก และเปิดหน้ารายการที่ยังไม่คืน
//	Then  หน้ารายละเอียดแสดงว่าคืนเกินกำหนด และหน้ารายการที่ยังไม่คืนแสดงว่ารายการที่สองเลยกำหนดคืน
//
// ใช้ภาควิชาใหม่เฉพาะ test นี้ เพื่อให้รายการที่ยังไม่คืนมีเฉพาะข้อมูลของ test
func TestAC_RET_13_OverdueShownToStaff(t *testing.T) {
	pool := testPool(t)
	dept, model := insertDepartment(t, pool, "T16 ภาควิชา")
	provider := fixedProvider{"staff-t16": {UserID: "staff-t16", Role: auth.RoleStaff, DepartmentID: dept}}
	router := newRouterWith(pool, provider)

	// Given รายการแรก: เลยกำหนดแล้วบันทึกคืน -> คืนเกินกำหนด / รายการที่สอง: เลยกำหนดแล้วยังไม่คืน
	insertLoanOfModel(t, pool, model, "T16-RETURNED-LATE", -2*time.Hour, false)
	insertLoanOfModel(t, pool, model, "T16-STILL-OUT", -1*time.Hour, false)
	rec := postReturn(t, router, "staff-t16", map[string]string{"assetCode": "T16-RETURNED-LATE", "condition": "normal"})
	if rec.Code != http.StatusOK {
		t.Fatalf("เตรียมข้อมูล: status = %d body=%s", rec.Code, rec.Body.String())
	}
	var created struct {
		LoanID int `json:"loanId"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}

	// When + Then: หน้ารายละเอียดของรายการแรก
	rec = call(router, http.MethodGet, fmt.Sprintf("/api/v1/returns/loans/%d", created.LoanID), "staff-t16")
	var detail LoanDetail
	if err := json.Unmarshal(rec.Body.Bytes(), &detail); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK || detail.ReturnedLate == nil || !*detail.ReturnedLate {
		t.Errorf("รายละเอียดรายการแรก: status=%d returnedLate=%v, want 200 และ true", rec.Code, detail.ReturnedLate)
	}

	// When + Then: หน้ารายการที่ยังไม่คืน มีเฉพาะรายการที่สอง และเลยกำหนดคืน
	rec = call(router, http.MethodGet, "/api/v1/returns/open-loans", "staff-t16")
	var open []OpenLoan
	if err := json.Unmarshal(rec.Body.Bytes(), &open); err != nil {
		t.Fatal(err)
	}
	if len(open) != 1 || open[0].AssetCode != "T16-STILL-OUT" || !open[0].IsOverdue {
		t.Errorf("รายการที่ยังไม่คืน = %+v, want [T16-STILL-OUT isOverdue=true]", open)
	}
}

// insertItemWithoutLoan สร้างอุปกรณ์ของภาควิชา 1 (รุ่น 1) ที่ไม่มีการยืมค้าง สถานะตามที่กำหนด
func insertItemWithoutLoan(t *testing.T, pool *pgxpool.Pool, assetCode, status string) {
	t.Helper()
	if _, err := pool.Exec(context.Background(),
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, $2)`, assetCode, status); err != nil {
		t.Fatal(err)
	}
}

func itemStatus(t *testing.T, pool *pgxpool.Pool, assetCode string) string {
	t.Helper()
	var s string
	if err := pool.QueryRow(context.Background(),
		`SELECT status FROM equipment_item WHERE asset_code = $1`, assetCode).Scan(&s); err != nil {
		t.Fatal(err)
	}
	return s
}

// AC-RET-03 (FR-RET-01)
//
//	Given อุปกรณ์ชิ้นหนึ่งสถานะ "ว่าง" หรือ "ไม่พร้อมใช้งาน" (ไม่มีการยืมที่ยังไม่คืน)
//	When  เจ้าหน้าที่กรอกรหัสครุภัณฑ์ของชิ้นนั้นเพื่อบันทึกการรับคืน
//	Then  ระบบไม่บันทึก แจ้งว่าอุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน และสถานะอุปกรณ์ไม่เปลี่ยน
func TestAC_RET_03_NoOpenLoanRejected(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)

	for _, tc := range []struct{ assetCode, status string }{
		{"T18-AVAILABLE", "available"},
		{"T18-UNAVAILABLE", "unavailable"},
	} {
		// Given
		insertItemWithoutLoan(t, pool, tc.assetCode, tc.status)

		// When: ทั้งขั้นค้นหา (lookup) และขั้นบันทึก (POST)
		lookup := call(router, http.MethodGet, "/api/v1/returns/lookup?assetCode="+tc.assetCode, "mock-staff-001")
		post := postReturn(t, router, "mock-staff-001", map[string]string{"assetCode": tc.assetCode, "condition": "normal"})

		// Then
		for name, rec := range map[string]*httptest.ResponseRecorder{"lookup": lookup, "POST": post} {
			if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), `"NO_OPEN_LOAN"`) {
				t.Errorf("%s %s: status = %d body=%s, want 409 NO_OPEN_LOAN", tc.status, name, rec.Code, rec.Body.String())
			}
		}
		if got := itemStatus(t, pool, tc.assetCode); got != tc.status {
			t.Errorf("%s: สถานะเปลี่ยนเป็น %q ต้องไม่เปลี่ยน", tc.assetCode, got)
		}
	}

	// รหัสครุภัณฑ์ที่ไม่มีในระบบ ก็ไม่มีการยืมค้างเช่นกัน
	if rec := postReturn(t, router, "mock-staff-001", map[string]string{"assetCode": "T18-NOT-EXIST", "condition": "normal"}); rec.Code != http.StatusConflict {
		t.Errorf("รหัสที่ไม่มี: status = %d, want 409", rec.Code)
	}
}

// T-18 เสร็จเมื่อ: บันทึกคืนชิ้นเดียวกันพร้อมกัน 2 request ได้ 200 หนึ่งครั้ง และ 409 หนึ่งครั้ง (plan.md ข้อ 4 FOR UPDATE)
// ทำซ้ำหลายชิ้นเพื่อให้มีโอกาสชนกันจริง
func TestReturns_ConcurrentReturn_OneWins(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		code := fmt.Sprintf("T18-RACE-%d", i)
		insertBorrowedItem(t, pool, code, 2*time.Hour)

		var wg sync.WaitGroup
		codes := make([]int, 2)
		start := make(chan struct{})
		for j := range codes {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				<-start
				codes[j] = postReturn(t, router, "mock-staff-001", map[string]string{"assetCode": code, "condition": "normal"}).Code
			}(j)
		}
		close(start)
		wg.Wait()

		ok, conflict := 0, 0
		for _, c := range codes {
			switch c {
			case http.StatusOK:
				ok++
			case http.StatusConflict:
				conflict++
			}
		}
		if ok != 1 || conflict != 1 {
			t.Errorf("%s: status = %v, want 200 หนึ่งครั้ง และ 409 หนึ่งครั้ง", code, codes)
		}

		var returns int
		if err := pool.QueryRow(ctx, `
			SELECT count(*) FROM loan l JOIN equipment_item i ON i.id = l.item_id
			WHERE i.asset_code = $1 AND l.returned_at IS NOT NULL`, code).Scan(&returns); err != nil {
			t.Fatal(err)
		}
		if returns != 1 {
			t.Errorf("%s: การยืมที่ปิดแล้ว = %d, want 1", code, returns)
		}
	}
}

// AC-RET-11 (FR-RET-01 ภาควิชา)
//
//	Given อุปกรณ์ "ถูกยืม" เป็นของภาควิชาอื่น
//	When  เจ้าหน้าที่กรอกรหัสครุภัณฑ์ของชิ้นนั้นเพื่อบันทึกการรับคืน
//	Then  ระบบปฏิเสธว่าไม่มีสิทธิ์ และการยืมกับสถานะอุปกรณ์ไม่เปลี่ยน
//
// ตรวจทุกทางที่เข้าถึงการยืมของภาควิชาอื่นได้: lookup, POST และรายละเอียดการยืม
func TestAC_RET_11_OtherDepartmentForbidden(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool) // mock-staff-001 อยู่ภาควิชา 1
	ctx := context.Background()

	// Given อุปกรณ์ "ถูกยืม" ของภาควิชาอื่นพร้อมการยืมค้าง และอีกชิ้นของภาควิชาอื่นที่ "ว่าง"
	_, otherModel := insertDepartment(t, pool, "T20 ภาควิชาอื่น")
	insertLoanOfModel(t, pool, otherModel, "T20-OTHER", 2*time.Hour, false)
	if _, err := pool.Exec(ctx,
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ('T20-OTHER-FREE', $1, 'available')`, otherModel); err != nil {
		t.Fatal(err)
	}
	var loanID int
	if err := pool.QueryRow(ctx,
		`SELECT l.id FROM loan l JOIN equipment_item i ON i.id = l.item_id WHERE i.asset_code = 'T20-OTHER'`).Scan(&loanID); err != nil {
		t.Fatal(err)
	}

	// When
	recs := map[string]*httptest.ResponseRecorder{
		"lookup":      call(router, http.MethodGet, "/api/v1/returns/lookup?assetCode=T20-OTHER", "mock-staff-001"),
		"POST":        postReturn(t, router, "mock-staff-001", map[string]string{"assetCode": "T20-OTHER", "condition": "normal"}),
		"loan detail": call(router, http.MethodGet, fmt.Sprintf("/api/v1/returns/loans/%d", loanID), "mock-staff-001"),
		"lookup ว่าง": call(router, http.MethodGet, "/api/v1/returns/lookup?assetCode=T20-OTHER-FREE", "mock-staff-001"),
		"POST ว่าง":   postReturn(t, router, "mock-staff-001", map[string]string{"assetCode": "T20-OTHER-FREE", "condition": "normal"}),
	}

	// Then: 403 ทุกทาง (ไม่บอกใบ้ว่ามีการยืมค้างหรือไม่) และไม่มีข้อมูลการยืมใน body
	for name, rec := range recs {
		if rec.Code != http.StatusForbidden || strings.TrimSpace(rec.Body.String()) != `{"code":"FORBIDDEN_DEPARTMENT"}` {
			t.Errorf("%s: status = %d body=%s, want 403 และ body มีแค่ code FORBIDDEN_DEPARTMENT", name, rec.Code, rec.Body.String())
		}
	}
	var returnedAt *time.Time
	if err := pool.QueryRow(ctx, `SELECT returned_at FROM loan WHERE id = $1`, loanID).Scan(&returnedAt); err != nil {
		t.Fatal(err)
	}
	if returnedAt != nil {
		t.Error("การยืมของภาควิชาอื่นต้องไม่ถูกปิด")
	}
	if got := itemStatus(t, pool, "T20-OTHER"); got != "borrowed" {
		t.Errorf("สถานะ = %q, want borrowed (ไม่เปลี่ยน)", got)
	}
}

// AC-RET-08 (FR-RET-01, IF-IDP-01)
//
//	Given นักศึกษายืนยันตัวตนแล้ว
//	When  เรียกบันทึกการรับคืน หรือเปิดหน้ารายการที่ยังไม่คืนโดยตรง
//	Then  ระบบปฏิเสธว่าไม่มีสิทธิ์ และไม่มีข้อมูลการยืมถูกส่งกลับหรือเปลี่ยน
//
// ตรวจทั้ง 4 endpoint ของ 006 ด้วย token นักศึกษาภาควิชาเดียวกับอุปกรณ์ (ภาควิชา 1)
func TestAC_RET_08_StudentForbidden(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)
	ctx := context.Background()

	// Given อุปกรณ์ภาควิชา 1 ที่ถูกยืมอยู่ และนักศึกษาภาควิชา 1
	insertBorrowedItem(t, pool, "T22-STUDENT", 2*time.Hour)
	var loanID int
	if err := pool.QueryRow(ctx,
		`SELECT l.id FROM loan l JOIN equipment_item i ON i.id = l.item_id WHERE i.asset_code = 'T22-STUDENT'`).Scan(&loanID); err != nil {
		t.Fatal(err)
	}
	const student = "mock-student-001"

	// When
	recs := map[string]*httptest.ResponseRecorder{
		"POST /returns":   postReturn(t, router, student, map[string]string{"assetCode": "T22-STUDENT", "condition": "normal"}),
		"GET /open-loans": call(router, http.MethodGet, "/api/v1/returns/open-loans", student),
		"GET /lookup":     call(router, http.MethodGet, "/api/v1/returns/lookup?assetCode=T22-STUDENT", student),
		"GET /loans/{id}": call(router, http.MethodGet, fmt.Sprintf("/api/v1/returns/loans/%d", loanID), student),
	}

	// Then: 403 FORBIDDEN_ROLE ทุก endpoint body มีแค่ code
	for name, rec := range recs {
		if rec.Code != http.StatusForbidden || strings.TrimSpace(rec.Body.String()) != `{"code":"FORBIDDEN_ROLE"}` {
			t.Errorf("%s: status = %d body=%s, want 403 และ body มีแค่ code FORBIDDEN_ROLE", name, rec.Code, rec.Body.String())
		}
	}
	var returnedAt *time.Time
	if err := pool.QueryRow(ctx, `SELECT returned_at FROM loan WHERE id = $1`, loanID).Scan(&returnedAt); err != nil {
		t.Fatal(err)
	}
	if returnedAt != nil || itemStatus(t, pool, "T22-STUDENT") != "borrowed" {
		t.Error("นักศึกษาเรียกแล้วการยืมหรือสถานะต้องไม่เปลี่ยน")
	}

	// เจ้าหน้าที่ยังใช้ได้ตามปกติ
	if rec := call(router, http.MethodGet, "/api/v1/returns/lookup?assetCode=T22-STUDENT", "mock-staff-001"); rec.Code != http.StatusOK {
		t.Errorf("เจ้าหน้าที่: status = %d, want 200", rec.Code)
	}
}

type auditRow struct {
	Role, Action, AssetCode, Outcome string
	Query                            *string
}

// auditRowsOf อ่าน audit_log ของผู้ใช้ actorID เรียงตามลำดับที่เขียน
func auditRowsOf(t *testing.T, pool *pgxpool.Pool, actorID string) []auditRow {
	t.Helper()
	rows, err := pool.Query(context.Background(), `
		SELECT actor_role, action, coalesce(asset_code, ''), outcome, query::text
		FROM audit_log WHERE actor_id = $1 ORDER BY id`, actorID)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []auditRow
	for rows.Next() {
		var a auditRow
		if err := rows.Scan(&a.Role, &a.Action, &a.AssetCode, &a.Outcome, &a.Query); err != nil {
			t.Fatal(err)
		}
		out = append(out, a)
	}
	return out
}

// AC-RET-07 (FR-RET-07, DOM-AUDIT-01)
//
//	Given เจ้าหน้าที่บันทึกการรับคืนสำเร็จ
//	When  การบันทึกเสร็จสิ้น
//	Then  มี audit log ของการรับคืนครั้งนั้นตามที่ DOM-AUDIT-01 กำหนด (ผู้เข้าถึง เวลา เงื่อนไขที่ค้น)
//
// ตรวจเพิ่มตาม T-23: การอ่านทุก endpoint มี audit และการถูกปฏิเสธมี outcome forbidden
func TestAC_RET_07_ReturnAudited(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	provider := fixedProvider{
		"staff-t23":   {UserID: "staff-t23", Role: auth.RoleStaff, DepartmentID: 1},
		"student-t23": {UserID: "student-t23", Role: auth.RoleStudent, DepartmentID: 1},
	}
	router := newRouterWith(pool, provider)

	// Given
	insertBorrowedItem(t, pool, "T23-AUDIT", 2*time.Hour)

	// When: เจ้าหน้าที่ค้นหา บันทึกคืนสำเร็จ เปิดรายละเอียด เปิดรายการ และบันทึกซ้ำ (ไม่มีการยืมค้างแล้ว)
	call(router, http.MethodGet, "/api/v1/returns/lookup?assetCode=T23-AUDIT", "staff-t23")
	rec := postReturn(t, router, "staff-t23", map[string]string{"assetCode": "T23-AUDIT", "condition": "normal"})
	if rec.Code != http.StatusOK {
		t.Fatalf("บันทึกคืน: status = %d body=%s", rec.Code, rec.Body.String())
	}
	var created struct {
		LoanID int `json:"loanId"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	call(router, http.MethodGet, fmt.Sprintf("/api/v1/returns/loans/%d", created.LoanID), "staff-t23")
	call(router, http.MethodGet, "/api/v1/returns/open-loans", "staff-t23")
	postReturn(t, router, "staff-t23", map[string]string{"assetCode": "T23-AUDIT", "condition": "normal"})

	// Then: แถว audit ของเจ้าหน้าที่ครบตามลำดับ
	got := auditRowsOf(t, pool, "staff-t23")
	want := []struct{ action, asset, outcome string }{
		{"return_lookup", "T23-AUDIT", "ok"},
		{"return", "T23-AUDIT", "ok"},
		{"return_detail", "T23-AUDIT", "ok"},
		{"return_list", "", "ok"},
		{"return", "T23-AUDIT", "error"}, // 409 NO_OPEN_LOAN
	}
	if len(got) != len(want) {
		t.Fatalf("audit ของเจ้าหน้าที่ = %+v, want %d แถว", got, len(want))
	}
	for i, w := range want {
		g := got[i]
		if g.Role != "staff" || g.Action != w.action || g.AssetCode != w.asset || g.Outcome != w.outcome {
			t.Errorf("แถว %d = %+v, want action=%s asset=%q outcome=%s role=staff", i, g, w.action, w.asset, w.outcome)
		}
	}
	if got[0].Query == nil || !strings.Contains(*got[0].Query, "T23-AUDIT") {
		t.Errorf("return_lookup ต้องเก็บเงื่อนไขที่ค้น: query = %v", got[0].Query)
	}

	// แถว "return" ที่สำเร็จมีเวลา และอยู่ใน transaction เดียวกับการปิดการยืม: เวลาเท่ากับ returned_at
	var sameTx bool
	if err := pool.QueryRow(ctx, `
		SELECT a.occurred_at = l.returned_at
		FROM audit_log a, loan l
		WHERE a.actor_id = 'staff-t23' AND a.action = 'return' AND a.outcome = 'ok' AND l.id = $1`, created.LoanID).Scan(&sameTx); err != nil {
		t.Fatal(err)
	}
	if !sameTx {
		t.Error("audit ของการรับคืนสำเร็จต้องเขียนใน transaction เดียวกับการปิดการยืม (now() เท่ากัน)")
	}

	// นักศึกษาถูกปฏิเสธ: มี audit outcome forbidden
	postReturn(t, router, "student-t23", map[string]string{"assetCode": "T23-AUDIT", "condition": "normal"})
	call(router, http.MethodGet, "/api/v1/returns/open-loans", "student-t23")
	st := auditRowsOf(t, pool, "student-t23")
	if len(st) != 2 || st[0].Action != "return" || st[1].Action != "return_list" ||
		st[0].Outcome != "forbidden" || st[1].Outcome != "forbidden" || st[0].Role != "student" {
		t.Errorf("audit ของนักศึกษา = %+v, want return และ return_list outcome forbidden", st)
	}
}
