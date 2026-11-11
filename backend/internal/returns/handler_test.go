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

// testPool เปิด PostgreSQL 16 ใน container ครั้งเดียว รัน migration (รวม 0006) และข้อมูลตัวอย่างของ 001 และ 006
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

// newTestRouter สร้าง router เหมือน main.go: /healthz + endpoint ของ 006 ผ่าน auth จำลอง
func newTestRouter(pool *pgxpool.Pool) http.Handler {
	h := NewHandler(NewRepository(pool))
	return platform.NewRouter(pool, platform.Routes(h.Routes(auth.NewMockProvider())))
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

// T-03 เสร็จเมื่อ: เรียก open-loans ด้วย token เจ้าหน้าที่ได้ 501 (ยังไม่ทำ) และไม่มี token ได้ 401 (IF-IDP-01)
func TestReturns_Skeleton_RoutedThroughAuth(t *testing.T) {
	router := newTestRouter(nil) // T-03 ยังไม่แตะฐานข้อมูล

	if rec := call(router, http.MethodGet, "/api/v1/returns/open-loans", "mock-staff-001"); rec.Code != http.StatusNotImplemented {
		t.Errorf("token เจ้าหน้าที่: status = %d, want 501", rec.Code)
	}
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

// T-03 เสร็จเมื่อ (ข้อมูลตัวอย่าง): seed_006_return.sql ใส่ได้หลัง seed.sql และ migration 0006
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
