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
		{http.MethodGet, "/api/v1/returns/loans/1/photos/1"},
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

// ไฟล์รูปขนาดเล็กที่ http.DetectContentType รู้จักจาก magic bytes
var (
	tinyJPEG = []byte("\xFF\xD8\xFF\xE0\x00\x10JFIF\x00test-jpeg")
	tinyPNG  = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDRtest-png")
)

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
func postReturn(t *testing.T, router http.Handler, token string, fields map[string]string, photos ...[]byte) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	for k, v := range fields {
		if err := mw.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
	}
	for i, p := range photos {
		fw, err := mw.CreateFormFile("photos", fmt.Sprintf("photo-%d", i))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := fw.Write(p); err != nil {
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
//	When  เจ้าหน้าที่บันทึกการรับคืน สภาพ "ปกติ" พร้อมรูป 1 รูป
//	Then  การยืมนั้นมีเวลาคืน และสถานะอุปกรณ์เป็น "ว่าง"
func TestAC_RET_01_NormalReturnBecomesAvailable(t *testing.T) {
	pool := testPool(t)
	router := newTestRouter(pool)
	ctx := context.Background()

	// Given
	insertBorrowedItem(t, pool, "T04-AC01", 2*time.Hour)

	// When
	rec := postReturn(t, router, "mock-staff-001",
		map[string]string{"assetCode": "T04-AC01", "condition": "normal"}, tinyJPEG)

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
	var photoCount int
	if err := pool.QueryRow(ctx, `
		SELECT l.returned_at, i.status,
		       (SELECT count(*) FROM loan_return_photo p WHERE p.loan_id = l.id)
		FROM loan l JOIN equipment_item i ON i.id = l.item_id
		WHERE i.asset_code = 'T04-AC01'`).Scan(&returnedAt, &status, &photoCount); err != nil {
		t.Fatal(err)
	}
	if returnedAt == nil {
		t.Error("การยืมต้องมีเวลาคืน (FR-RET-02)")
	}
	if status != "available" {
		t.Errorf("สถานะ = %q, want available (FR-RET-03)", status)
	}
	if photoCount != 1 {
		t.Errorf("รูปที่เก็บ = %d, want 1 (FR-RET-05)", photoCount)
	}
}
