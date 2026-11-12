package equipment

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	// ผู้ใช้จำลองจาก auth.NewMockProvider
	tokenStudentDept1 = "mock-student-001"
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

// testPool เปิด PostgreSQL 16 ใน container ครั้งเดียว แล้วรัน migration และข้อมูลตั้งต้น (plan.md ข้อ 2: testcontainers)
// ข้อมูลตั้งต้นคือ backend/db/seed/seed.sql (รุ่น 1 "Notebook ตัวอย่างรุ่น A" 5 ชิ้น ว่าง 3)
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

// newTestServer ประกอบ handler แบบเดียวกับ cmd/server/main.go (อยู่หลัง auth.Middleware และห่อด้วย platform.NoStore)
func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	h := NewHandler(NewService(NewRepository(testPool(t))))
	mux := http.NewServeMux()
	mux.Handle(PathPrefix, auth.Middleware(auth.NewMockProvider(), h.Routes()))
	return platform.NoStore(mux)
}

// get เรียก endpoint ด้วย token ที่กำหนด (ว่าง = ไม่ส่ง token)
func get(t *testing.T, srv http.Handler, path, token string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.NewDecoder(rec.Body).Decode(&v); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
	}
	return v
}

// AC-EQS-01 (FR-EQS-01) ส่วนรายการรุ่น
//
//	Given นักศึกษายืนยันตัวตนแล้ว และมีอุปกรณ์รุ่นหนึ่ง 5 ชิ้น ว่าง 3 ชิ้น
//	When  เปิดหน้ารายการอุปกรณ์
//	Then  รายการแสดงรุ่นนั้น 1 แถว ระบุประเภทและ "ว่าง 3 / ทั้งหมด 5"
func TestAC_EQS_01_ModelRowCounts(t *testing.T) {
	// Given: ข้อมูลตั้งต้นมีรุ่น 1 "Notebook ตัวอย่างรุ่น A" 5 ชิ้น ว่าง 3 และนักศึกษาภาควิชา 1 ยืนยันตัวตนแล้ว
	srv := newTestServer(t)

	// When: เปิดหน้ารายการอุปกรณ์
	rec := get(t, srv, PathPrefix+"models", tokenStudentDept1)

	// Then
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, ต้องเป็น 200 (body: %s)", rec.Code, rec.Body.String())
	}
	body := decode[ListDTO[ModelRowDTO]](t, rec)

	var rows []ModelRowDTO
	for _, m := range body.Items {
		if m.ModelID == 1 {
			rows = append(rows, m)
		}
	}
	if len(rows) != 1 {
		t.Fatalf("รุ่นนั้นต้องแสดง 1 แถว แต่ได้ %d แถว", len(rows))
	}
	row := rows[0]
	if row.Name != "Notebook ตัวอย่างรุ่น A" {
		t.Errorf("ชื่อรุ่น = %q", row.Name)
	}
	if row.Category != "คอมพิวเตอร์พกพา" {
		t.Errorf("ประเภท = %q, ต้องเป็น \"คอมพิวเตอร์พกพา\"", row.Category)
	}
	if row.AvailableCount != 3 || row.TotalCount != 5 {
		t.Errorf("ได้ ว่าง %d / ทั้งหมด %d, ต้องเป็น ว่าง 3 / ทั้งหมด 5", row.AvailableCount, row.TotalCount)
	}
	if body.Total != len(body.Items) {
		t.Errorf("total = %d แต่ items มี %d แถว", body.Total, len(body.Items))
	}
}

// ทุก endpoint ต้องผ่านการยืนยันตัวตนก่อน (IF-IDP-01, NFR-SEC-01) ตรวจว่า endpoint ใหม่อยู่หลัง middleware จริง
func TestModels_RequiresAuthentication(t *testing.T) {
	srv := newTestServer(t)

	rec := get(t, srv, PathPrefix+"models", "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("ไม่มี token ต้องได้ 401 แต่ได้ %d", rec.Code)
	}
}

// AC-EQS-01 (FR-EQS-01, CON-DATA-03) ส่วนรายชิ้น
//
//	Given นักศึกษายืนยันตัวตนแล้ว และมีอุปกรณ์รุ่นหนึ่ง 5 ชิ้น ว่าง 3 ชิ้น
//	When  เลือกรุ่นนั้น
//	Then  เห็นอุปกรณ์ 5 ชิ้นตามรหัสครุภัณฑ์พร้อมสถานะของแต่ละชิ้น
func TestAC_EQS_01_ItemsOfModel(t *testing.T) {
	// Given: รุ่น 1 ในข้อมูลตั้งต้น CS-NB-0001 ถึง 0003 ว่าง, 0004 และ 0005 ถูกยืม
	srv := newTestServer(t)

	// When: เลือกรุ่นนั้น
	rec := get(t, srv, PathPrefix+"models/1/items", tokenStudentDept1)

	// Then
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, ต้องเป็น 200 (body: %s)", rec.Code, rec.Body.String())
	}
	body := decode[ListDTO[ItemRow]](t, rec)
	want := []ItemRow{
		{"CS-NB-0001", StatusAvailable},
		{"CS-NB-0002", StatusAvailable},
		{"CS-NB-0003", StatusAvailable},
		{"CS-NB-0004", StatusBorrowed},
		{"CS-NB-0005", StatusBorrowed},
	}
	if len(body.Items) != len(want) || body.Total != len(want) {
		t.Fatalf("ต้องเห็น 5 ชิ้น แต่ได้ %d ชิ้น (total %d): %+v", len(body.Items), body.Total, body.Items)
	}
	for i, w := range want {
		if body.Items[i] != w {
			t.Errorf("ชิ้นที่ %d = %+v, ต้องเป็น %+v", i+1, body.Items[i], w)
		}
	}
}

// รุ่นที่ไม่มีอยู่ได้รายการว่าง ไม่ใช่ error (plan.md ข้อ 5) และรหัสรุ่นที่ไม่ใช่ตัวเลขได้ 404
func TestItems_UnknownModel(t *testing.T) {
	srv := newTestServer(t)

	rec := get(t, srv, PathPrefix+"models/9999/items", tokenStudentDept1)
	if rec.Code != http.StatusOK {
		t.Fatalf("รุ่นที่ไม่มีอยู่ status = %d, ต้องเป็น 200", rec.Code)
	}
	if body := decode[ListDTO[ItemRow]](t, rec); len(body.Items) != 0 || body.Items == nil {
		t.Errorf("ต้องได้ items เป็น [] แต่ได้ %+v", body.Items)
	}

	if rec := get(t, srv, PathPrefix+"models/abc/items", tokenStudentDept1); rec.Code != http.StatusNotFound {
		t.Errorf("รหัสรุ่นไม่ใช่ตัวเลข status = %d, ต้องเป็น 404", rec.Code)
	}
}

// modelsByID เรียกรายการรุ่นด้วย query string ที่กำหนด แล้วคืนเป็น map ตามรหัสรุ่น
func modelsByID(t *testing.T, srv http.Handler, query url.Values) map[int]ModelRowDTO {
	t.Helper()
	rec := get(t, srv, PathPrefix+"models?"+query.Encode(), tokenStudentDept1)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, ต้องเป็น 200 (body: %s)", rec.Code, rec.Body.String())
	}
	out := map[int]ModelRowDTO{}
	for _, m := range decode[ListDTO[ModelRowDTO]](t, rec).Items {
		out[m.ModelID] = m
	}
	return out
}

// AC-EQS-02 (FR-EQS-02) ส่วนกรองสถานะ (ส่วนค้นหาด้วยคำว่า "Notebook" รอ Q-06 ดู T-24)
//
//	Given อยู่ในหน้ารายการอุปกรณ์
//	When  กรองสถานะ "ว่าง"
//	Then  แสดงเฉพาะอุปกรณ์ที่มีสถานะว่างเท่านั้น
func TestAC_EQS_02_FilterStatus(t *testing.T) {
	// Given: ข้อมูลตั้งต้น รุ่น 1 ว่าง 3 ถูกยืม 2 / รุ่น 2 ว่าง 1 ไม่พร้อมใช้งาน 1 / รุ่น 3 ว่าง 1
	srv := newTestServer(t)

	// When: กรองสถานะ "ว่าง"
	models := modelsByID(t, srv, url.Values{"status": {"available"}})

	// Then: ทุกแถวนับเฉพาะชิ้นที่ว่าง
	for id, m := range models {
		if m.TotalCount != m.AvailableCount {
			t.Errorf("รุ่น %d นับชิ้นที่ไม่ว่างด้วย: ว่าง %d / ทั้งหมด %d", id, m.AvailableCount, m.TotalCount)
		}
	}
	if m, ok := models[1]; !ok || m.TotalCount != 3 {
		t.Errorf("รุ่น 1 ต้องเหลือ 3 ชิ้นที่ว่าง ได้ %+v", m)
	}

	// Then: รายชิ้นของรุ่นนั้นใช้ตัวกรองเดียวกัน เห็นเฉพาะชิ้นที่ว่าง
	rec := get(t, srv, PathPrefix+"models/1/items?status=available", tokenStudentDept1)
	items := decode[ListDTO[ItemRow]](t, rec).Items
	if len(items) != 3 {
		t.Fatalf("รายชิ้นต้องเหลือ 3 ชิ้น ได้ %d: %+v", len(items), items)
	}
	for _, it := range items {
		if it.Status != StatusAvailable {
			t.Errorf("%s สถานะ %s ไม่ใช่ว่าง", it.AssetCode, it.Status)
		}
	}
}

// FR-EQS-02 แสดงเฉพาะรุ่นที่มีอย่างน้อย 1 ชิ้นตรงเงื่อนไข ทั้งตัวกรองสถานะและประเภท
func TestFR_EQS_02_OnlyModelsWithMatchingItems(t *testing.T) {
	srv := newTestServer(t)

	// มีชิ้นที่ไม่พร้อมใช้งานแค่ CS-PJ-0002 ของรุ่น 2
	models := modelsByID(t, srv, url.Values{"status": {"unavailable"}})
	if len(models) != 1 || models[2].TotalCount != 1 {
		t.Errorf("กรอง \"ไม่พร้อมใช้งาน\" ต้องเหลือแค่รุ่น 2 จำนวน 1 ชิ้น ได้ %+v", models)
	}

	// ประเภท "อุปกรณ์นำเสนอ" มีแค่รุ่น 2
	models = modelsByID(t, srv, url.Values{"category": {"อุปกรณ์นำเสนอ"}})
	if len(models) != 1 || models[2].TotalCount != 2 {
		t.Errorf("กรองประเภท \"อุปกรณ์นำเสนอ\" ต้องเหลือแค่รุ่น 2 ได้ %+v", models)
	}

	// ค่าที่ไม่รู้จัก (รวมข้อความที่หน้าตาเป็น SQL) ไม่ตรงกับชิ้นใด และไม่ทำให้ query พัง (parameterized query)
	models = modelsByID(t, srv, url.Values{"category": {"x' OR '1'='1"}, "status": {"pending"}})
	if len(models) != 0 {
		t.Errorf("ค่าที่ไม่รู้จักต้องไม่ตรงกับรุ่นใด ได้ %+v", models)
	}
}

// AC-EQS-15 (FR-EQS-01, FR-EQS-02)
//
//	Given รุ่นหนึ่งมี 5 ชิ้น ว่าง 3 ถูกยืม 2
//	When  กรองสถานะ "ถูกยืม"
//	Then  รุ่นนั้นแสดง 1 แถว และตัวเลขนับเฉพาะชิ้นที่ตรงตัวกรอง (2 ชิ้น)
func TestAC_EQS_15_CountsFollowFilter(t *testing.T) {
	// Given: รุ่น 1 ในข้อมูลตั้งต้น 5 ชิ้น ว่าง 3 ถูกยืม 2
	srv := newTestServer(t)

	// When: กรองสถานะ "ถูกยืม"
	rec := get(t, srv, PathPrefix+"models?status=borrowed", tokenStudentDept1)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, ต้องเป็น 200", rec.Code)
	}
	var rows []ModelRowDTO
	for _, m := range decode[ListDTO[ModelRowDTO]](t, rec).Items {
		if m.ModelID == 1 {
			rows = append(rows, m)
		}
	}

	// Then: รุ่นนั้น 1 แถว และนับเฉพาะชิ้นที่ตรงตัวกรอง (2 ชิ้น ซึ่งไม่มีชิ้นไหนว่าง)
	if len(rows) != 1 {
		t.Fatalf("รุ่นนั้นต้องแสดง 1 แถว ได้ %d แถว", len(rows))
	}
	if rows[0].TotalCount != 2 {
		t.Errorf("ตัวเลขต้องนับเฉพาะชิ้นที่ตรงตัวกรอง (2 ชิ้น) ได้ totalCount = %d", rows[0].TotalCount)
	}
	if rows[0].AvailableCount != 0 {
		t.Errorf("ชิ้นที่ตรงตัวกรอง \"ถูกยืม\" ไม่มีชิ้นที่ว่าง ได้ availableCount = %d", rows[0].AvailableCount)
	}
}

// AC-EQS-16 (FR-EQS-02)
//
//	Given มีอุปกรณ์สถานะ ว่าง ถูกยืม และไม่พร้อมใช้งาน
//	When  เลือกตัวกรองสถานะ "ว่าง" และ "ถูกยืม" พร้อมกัน
//	Then  แสดงเฉพาะชิ้นที่ว่างหรือถูกยืม ไม่แสดงชิ้นที่ไม่พร้อมใช้งาน
func TestAC_EQS_16_MultiStatusFilter(t *testing.T) {
	// Given: ข้อมูลตั้งต้นมีครบ 3 สถานะ (CS-PJ-0002 ของรุ่น 2 ไม่พร้อมใช้งาน)
	srv := newTestServer(t)
	filter := url.Values{"status": {"available", "borrowed"}}

	// When: เลือก "ว่าง" และ "ถูกยืม" พร้อมกัน
	models := modelsByID(t, srv, filter)

	// Then: รุ่น 1 ได้ครบ 5 ชิ้น (ว่าง 3 + ถูกยืม 2) และรุ่น 2 ไม่นับชิ้นที่ไม่พร้อมใช้งาน
	if models[1].TotalCount != 5 {
		t.Errorf("รุ่น 1 ต้องนับ 5 ชิ้น ได้ %+v", models[1])
	}
	if models[2].TotalCount != 1 {
		t.Errorf("รุ่น 2 ต้องนับแค่ชิ้นที่ว่าง 1 ชิ้น ได้ %+v", models[2])
	}

	// Then: รายชิ้นของรุ่น 2 ไม่มีชิ้นที่ไม่พร้อมใช้งาน
	rec := get(t, srv, PathPrefix+"models/2/items?"+filter.Encode(), tokenStudentDept1)
	items := decode[ListDTO[ItemRow]](t, rec).Items
	if len(items) != 1 || items[0].AssetCode != "CS-PJ-0001" {
		t.Fatalf("รุ่น 2 ต้องเหลือแค่ CS-PJ-0001 ได้ %+v", items)
	}
	for _, it := range items {
		if it.Status == StatusUnavailable {
			t.Errorf("%s ไม่พร้อมใช้งาน ต้องไม่แสดง", it.AssetCode)
		}
	}
}

// AC-EQS-03 (FR-EQS-03) ส่วนหลังบ้าน: ไม่พบ = 200 กับรายการว่าง ไม่ใช่ error (plan.md ข้อ 5)
// ยังไม่มีช่องค้นหาด้วยคำ (รอ Q-06 ดู T-24) จึงใช้ตัวกรองที่ไม่ตรงกับชิ้นใด ตามที่ทีมตัดสิน 2569-11-10
//
//	Given ค้นหา/กรองด้วยเงื่อนไขที่ไม่มีในระบบ
//	When  กดค้นหา
//	Then  ได้รายการว่าง (หน้าจอแจ้งว่าไม่พบอุปกรณ์)
func TestAC_EQS_03_EmptyIs200(t *testing.T) {
	// Given: ประเภท "อุปกรณ์นำเสนอ" ไม่มีชิ้นที่ถูกยืมในข้อมูลตั้งต้น
	srv := newTestServer(t)
	query := url.Values{"category": {"อุปกรณ์นำเสนอ"}, "status": {"borrowed"}}

	// When
	rec := get(t, srv, PathPrefix+"models?"+query.Encode(), tokenStudentDept1)

	// Then: 200 และ items เป็น [] ไม่ใช่ null หรือ error
	if rec.Code != http.StatusOK {
		t.Fatalf("ไม่พบต้องได้ 200 แต่ได้ %d (body: %s)", rec.Code, rec.Body.String())
	}
	raw := rec.Body.String()
	body := decode[ListDTO[ModelRowDTO]](t, rec)
	if body.Items == nil || len(body.Items) != 0 || body.Total != 0 {
		t.Errorf("ต้องได้ {items: [], total: 0} แต่ได้ %s", raw)
	}
}

// AC-EQS-04 (FR-EQS-04, CON-DATA-02) ส่วนสถานะและกำหนดคืน (ส่วน "ไม่มีปุ่มส่งคำขอยืม" รอ Q-16 ดู T-34)
//
//	Given เจ้าหน้าที่บันทึกการส่งมอบอุปกรณ์ชิ้นหนึ่งใน UC-08 แล้ว และยังไม่เลยกำหนดคืน
//	When  นักศึกษาเปิดดูอุปกรณ์ชิ้นนั้น
//	Then  สถานะแสดงเป็น "ถูกยืม" พร้อมกำหนดคืน (วันที่และเวลา)
func TestAC_EQS_04_BorrowedShowsDueDate(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()

	// Given: จำลอง UC-08 ส่งมอบชิ้นใหม่ของรุ่น 1 กำหนดคืนพรุ่งนี้ (ยังไม่เลยกำหนด) แล้วลบทิ้งหลัง test
	const code = "TEST-AC04-0001"
	var want time.Time
	err := pool.QueryRow(ctx, `
WITH item AS (
  INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, 'borrowed') RETURNING id
)
INSERT INTO loan (item_id, borrower_id, loan_type, due_at)
SELECT id, 'student-002', 'normal', date_trunc('minute', now()) + interval '1 day' FROM item
RETURNING due_at`, code).Scan(&want)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM loan WHERE item_id = (SELECT id FROM equipment_item WHERE asset_code = $1)`, code)
		_, _ = pool.Exec(ctx, `DELETE FROM equipment_item WHERE asset_code = $1`, code)
	})
	srv := newTestServer(t)

	// When: นักศึกษาเปิดดูอุปกรณ์ชิ้นนั้น
	rec := get(t, srv, PathPrefix+"items/"+code, tokenStudentDept1)

	// Then: สถานะ "ถูกยืม" พร้อมกำหนดคืน (วันที่และเวลา)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, ต้องเป็น 200 (body: %s)", rec.Code, rec.Body.String())
	}
	item := decode[ItemDetailDTO](t, rec)
	if item.Status != StatusBorrowed {
		t.Errorf("สถานะ = %s, ต้องเป็น borrowed", item.Status)
	}
	if item.DueAt == nil {
		t.Fatal("ต้องมีกำหนดคืน (dueAt)")
	}
	if !item.DueAt.Equal(want) {
		t.Errorf("กำหนดคืน = %s, ต้องเป็น %s (วันที่และเวลาตามที่ UC-08 บันทึก)", item.DueAt, want)
	}
}

// FR-EQS-04 กำหนดคืนมีเฉพาะชิ้นที่ถูกยืม และรหัสที่ไม่มีได้ 404
func TestItemDetail_AvailableHasNoDueAtAndUnknownIs404(t *testing.T) {
	srv := newTestServer(t)

	rec := get(t, srv, PathPrefix+"items/CS-NB-0001", tokenStudentDept1)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if item := decode[ItemDetailDTO](t, rec); item.Status != StatusAvailable || item.DueAt != nil {
		t.Errorf("ชิ้นที่ว่างต้องไม่มีกำหนดคืน ได้ %+v", item)
	}

	if rec := get(t, srv, PathPrefix+"items/NO-SUCH-CODE", tokenStudentDept1); rec.Code != http.StatusNotFound {
		t.Errorf("รหัสที่ไม่มีต้องได้ 404 ได้ %d", rec.Code)
	}
}

// CON-DATA-02 ฟีเจอร์นี้อ่านอย่างเดียว: ไม่มี endpoint ใดเปลี่ยนสถานะอุปกรณ์ได้
func TestCON_DATA_02_NoWriteEndpoints(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	paths := []string{PathPrefix + "models", PathPrefix + "models/1/items", PathPrefix + "items/CS-NB-0001"}

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		for _, p := range paths {
			req := httptest.NewRequest(method, p, strings.NewReader(`{"status":"unavailable"}`))
			req.Header.Set("Authorization", "Bearer "+tokenStudentDept1)
			rec := httptest.NewRecorder()
			srv.ServeHTTP(rec, req)
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("%s %s ได้ %d ต้องเป็น 405 (ไม่มี endpoint เขียนข้อมูล)", method, p, rec.Code)
			}
		}
	}

	var status string
	if err := pool.QueryRow(context.Background(), `SELECT status FROM equipment_item WHERE asset_code = 'CS-NB-0001'`).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(StatusAvailable) {
		t.Errorf("สถานะของ CS-NB-0001 เปลี่ยนเป็น %s", status)
	}
}

// AC-EQS-05 (FR-EQS-06, NFR-DATA-01) ส่วนสถานะเปลี่ยนเป็น "ว่าง"
// (ส่วน "ปุ่มส่งคำขอยืมกลับมาใช้งานได้" รอ Q-16 ดู T-34 / หน้าจอรีเฟรชเองไหม รอ Q-18 ดู T-35)
//
//	Given เจ้าหน้าที่บันทึกการรับคืนใน UC-09 สภาพปกติ
//	When  นักศึกษารีเฟรชหน้ารายการภายใน 30 วินาที
//	Then  สถานะเปลี่ยนเป็น "ว่าง"
func TestAC_EQS_05_ReturnedBecomesAvailable(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	srv := newTestServer(t)

	// เตรียม: ชิ้นใหม่ของรุ่น 1 ถูกยืมอยู่ (จำลอง UC-08) ลบทิ้งหลัง test
	const code = "TEST-AC05-0001"
	_, err := pool.Exec(ctx, `
WITH item AS (
  INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, 'borrowed') RETURNING id
)
INSERT INTO loan (item_id, borrower_id, loan_type, due_at)
SELECT id, 'student-001', 'normal', now() + interval '1 day' FROM item`, code)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM loan WHERE item_id = (SELECT id FROM equipment_item WHERE asset_code = $1)`, code)
		_, _ = pool.Exec(ctx, `DELETE FROM equipment_item WHERE asset_code = $1`, code)
	})
	before := modelsByID(t, srv, url.Values{})[1].AvailableCount
	if got := decode[ItemDetailDTO](t, get(t, srv, PathPrefix+"items/"+code, tokenStudentDept1)); got.Status != StatusBorrowed {
		t.Fatalf("ก่อนรับคืนต้องเป็น borrowed ได้ %s", got.Status)
	}

	// Given: จำลอง UC-09 รับคืนสภาพปกติ ปิดการยืมและตั้งสถานะเป็นว่าง
	for _, q := range []string{
		`UPDATE loan SET returned_at = now() WHERE item_id = (SELECT id FROM equipment_item WHERE asset_code = $1) AND returned_at IS NULL`,
		`UPDATE equipment_item SET status = 'available', updated_at = now() WHERE asset_code = $1`,
	} {
		if _, err := pool.Exec(ctx, q, code); err != nil {
			t.Fatal(err)
		}
	}

	// When: นักศึกษารีเฟรช (เรียกใหม่ทันที อยู่ในช่วง 30 วินาทีแน่นอน)
	rec := get(t, srv, PathPrefix+"items/"+code, tokenStudentDept1)

	// Then: สถานะเป็น "ว่าง" ไม่มีกำหนดคืนค้าง และตัวเลขว่างของรุ่นเพิ่มขึ้น 1
	item := decode[ItemDetailDTO](t, rec)
	if item.Status != StatusAvailable || item.DueAt != nil {
		t.Errorf("หลังรับคืนต้องเป็น available ไม่มี dueAt ได้ %+v", item)
	}
	if after := modelsByID(t, srv, url.Values{})[1].AvailableCount; after != before+1 {
		t.Errorf("ตัวเลขว่างของรุ่น 1 ต้องเพิ่มจาก %d เป็น %d ได้ %d", before, before+1, after)
	}

	// Then: response บอกห้าม cache ทุก endpoint กดรีเฟรชจึงได้ข้อมูลจากฐานข้อมูลใหม่เสมอ (plan.md ข้อ 6.2)
	for _, p := range []string{PathPrefix + "items/" + code, PathPrefix + "models", PathPrefix + "models/1/items"} {
		if cc := get(t, srv, p, tokenStudentDept1).Header().Get("Cache-Control"); cc != "no-store" {
			t.Errorf("%s Cache-Control = %q ต้องเป็น no-store", p, cc)
		}
	}
}

// AC-EQS-07 (FR-EQS-07, DOM-PRIV-01) ส่วนนักศึกษา (ส่วนเจ้าหน้าที่เห็นชื่อผู้ยืมรอ Q-20 ดู T-36)
//
//	Given อุปกรณ์กำลังถูกยืมโดยนักศึกษาคนอื่น
//	When  นักศึกษาทั่วไปเปิดดูรายละเอียด
//	Then  นักศึกษาไม่เห็นชื่อผู้ยืม
func TestAC_EQS_07_StudentNoBorrower(t *testing.T) {
	// Given: ข้อมูลตั้งต้น CS-NB-0004 ถูกยืมโดย student-001
	srv := newTestServer(t)

	// When: student-002 (นักศึกษาอีกคน) เปิดดูรายละเอียด และเปิดรายชิ้นของรุ่นนั้น
	for _, path := range []string{PathPrefix + "items/CS-NB-0004", PathPrefix + "models/1/items", PathPrefix + "models"} {
		rec := get(t, srv, path, "mock-student-002")
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d", path, rec.Code)
		}
		body := strings.ToLower(rec.Body.String())

		// Then: response ไม่มีข้อมูลผู้ยืมเลย ทั้งชื่อช่องและรหัสผู้ยืม (DOM-PRIV-01, plan.md ข้อ 6.1)
		if strings.Contains(body, "borrower") || strings.Contains(body, "student-001") {
			t.Errorf("%s ส่งข้อมูลผู้ยืมให้นักศึกษา: %s", path, rec.Body.String())
		}
	}
}

// FR-EQS-07 "สรุปจำนวนอุปกรณ์แยกตามสถานะ" ให้เจ้าหน้าที่ นักศึกษาไม่เห็น (ส่วนนี้ยังไม่มี AC ดู tasks.md "สิ่งที่ยังไม่ทำ")
func TestFR_EQS_07_StaffSummaryByStatus(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)

	// เจ้าหน้าที่ได้จำนวนแยก 3 สถานะ ตรงกับที่นับจากฐานข้อมูลโดยตรง
	rec := get(t, srv, PathPrefix+"summary", "mock-staff-001")
	if rec.Code != http.StatusOK {
		t.Fatalf("เจ้าหน้าที่ status = %d (body: %s)", rec.Code, rec.Body.String())
	}
	got := decode[StatusSummaryDTO](t, rec)
	var want StatusSummaryDTO
	err := pool.QueryRow(context.Background(), `
SELECT count(*) FILTER (WHERE status = 'available'), count(*) FILTER (WHERE status = 'borrowed'),
       count(*) FILTER (WHERE status = 'unavailable'), count(*) FROM equipment_item`).
		Scan(&want.Available, &want.Borrowed, &want.Unavailable, &want.Total)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("สรุป = %+v ต้องเป็น %+v", got, want)
	}
	if got.Available+got.Borrowed+got.Unavailable != got.Total {
		t.Errorf("ผลรวม 3 สถานะ (%d) ไม่เท่ากับ total (%d)", got.Available+got.Borrowed+got.Unavailable, got.Total)
	}

	// นักศึกษาได้ 403 FORBIDDEN_ROLE ตาม CON-TECH-03 และไม่ได้ตัวเลขใด ๆ
	rec = get(t, srv, PathPrefix+"summary", tokenStudentDept1)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("นักศึกษา status = %d ต้องเป็น 403", rec.Code)
	}
	e := decode[map[string]string](t, rec)
	if e["code"] != "FORBIDDEN_ROLE" || e["message"] == "" || e["message"] == "FORBIDDEN_ROLE" || len(e) != 2 {
		t.Errorf("body ต้องเป็น {code: FORBIDDEN_ROLE, message: ข้อความไทย} ได้ %v", e)
	}
}

// AC-EQS-13 (FR-EQS-10, DOM-ACCESS-01)
//
//	Given นักศึกษาภาควิชา A ยืนยันตัวตนแล้ว และอุปกรณ์เป็นของภาควิชา B
//	When  เปิดหน้ารายการอุปกรณ์ หรือเรียกรายละเอียดอุปกรณ์ของภาควิชา B โดยตรง
//	Then  แจ้งว่าไม่มีสิทธิ์เข้าถึง และไม่มีข้อมูลอุปกรณ์ของภาควิชา B ถูกส่งกลับ
func TestAC_EQS_13_OtherDepartment403(t *testing.T) {
	// Given: นักศึกษาภาควิชา 2 (mock-student-101) / อุปกรณ์ภาควิชา 1 คือรุ่น 1, 2 และ CS-NB-*, CS-PJ-*
	srv := newTestServer(t)
	const studentDept2 = "mock-student-101"

	// When: เปิดหน้ารายการอุปกรณ์
	rec := get(t, srv, PathPrefix+"models", studentDept2)

	// Then: เห็นเฉพาะรุ่นของภาควิชาตน ไม่มีข้อมูลรุ่นของภาควิชา 1 ปนมา
	if rec.Code != http.StatusOK {
		t.Fatalf("รายการรุ่น status = %d", rec.Code)
	}
	for _, m := range decode[ListDTO[ModelRowDTO]](t, rec).Items {
		if m.ModelID != 3 {
			t.Errorf("นักศึกษาภาควิชา 2 เห็นรุ่น %d (%s) ของภาควิชาอื่น", m.ModelID, m.Name)
		}
	}

	// When: ขอข้อมูลของภาควิชา 1 โดยตรง 3 ทาง
	// Then: 403 FORBIDDEN_DEPARTMENT body มีแค่ code และ message ไม่มีข้อมูลอุปกรณ์
	for _, path := range []string{
		PathPrefix + "items/CS-NB-0001",      // รายละเอียดรายชิ้น
		PathPrefix + "models/1/items",        // รายชิ้นของรุ่นภาควิชาอื่น
		PathPrefix + "models?department=1",   // ระบุภาควิชาอื่นในรายการ
		PathPrefix + "models?department=abc", // ค่าภาควิชาที่ไม่ถูกต้อง
	} {
		rec := get(t, srv, path, studentDept2)
		if rec.Code != http.StatusForbidden {
			t.Errorf("%s status = %d ต้องเป็น 403", path, rec.Code)
			continue
		}
		raw := rec.Body.String()
		e := decode[map[string]string](t, rec)
		if e["code"] != "FORBIDDEN_DEPARTMENT" || e["message"] == "" || len(e) != 2 {
			t.Errorf("%s body ต้องเป็น {code: FORBIDDEN_DEPARTMENT, message} ได้ %s", path, raw)
		}
		if strings.Contains(raw, "CS-") || strings.Contains(raw, "Notebook") || strings.Contains(raw, "status") {
			t.Errorf("%s ส่งข้อมูลอุปกรณ์ของภาควิชาอื่นกลับไป: %s", path, raw)
		}
	}

	// ภาควิชาของตนยังเปิดได้ปกติ และระบุ department ของตนได้
	if rec := get(t, srv, PathPrefix+"items/XX-NB-0001", studentDept2); rec.Code != http.StatusOK {
		t.Errorf("อุปกรณ์ภาควิชาตนเอง status = %d ต้องเป็น 200", rec.Code)
	}
	if rec := get(t, srv, PathPrefix+"models?department=2", studentDept2); rec.Code != http.StatusOK {
		t.Errorf("ระบุภาควิชาตนเอง status = %d ต้องเป็น 200", rec.Code)
	}
}

// FR-EQS-10 กลับด้าน: นักศึกษาภาควิชา 1 ไม่เห็นรุ่นของภาควิชา 2 / เจ้าหน้าที่เห็นทุกภาควิชา (plan.md ข้อ 6.1 รอ Q-P3)
func TestFR_EQS_10_ScopeByRole(t *testing.T) {
	srv := newTestServer(t)

	student := modelsByID(t, srv, url.Values{})
	if _, ok := student[3]; ok {
		t.Error("นักศึกษาภาควิชา 1 เห็นรุ่น 3 ของภาควิชา 2")
	}
	if rec := get(t, srv, PathPrefix+"items/XX-NB-0001", tokenStudentDept1); rec.Code != http.StatusForbidden {
		t.Errorf("นักศึกษาภาควิชา 1 เปิด XX-NB-0001 status = %d ต้องเป็น 403", rec.Code)
	}

	rec := get(t, srv, PathPrefix+"models", "mock-staff-001")
	ids := map[int]bool{}
	for _, m := range decode[ListDTO[ModelRowDTO]](t, rec).Items {
		ids[m.ModelID] = true
	}
	if !ids[1] || !ids[2] || !ids[3] {
		t.Errorf("เจ้าหน้าที่ต้องเห็นทุกภาควิชา (รุ่น 1, 2, 3) ได้ %v", ids)
	}
	if rec := get(t, srv, PathPrefix+"items/XX-NB-0001", "mock-staff-001"); rec.Code != http.StatusOK {
		t.Errorf("เจ้าหน้าที่เปิด XX-NB-0001 status = %d ต้องเป็น 200", rec.Code)
	}
}

// AC-EQS-08 (FR-EQS-08) ส่วนหลังบ้าน
//
//	Given ฐานข้อมูลไม่ตอบสนอง
//	When  เปิดหน้ารายการอุปกรณ์
//	Then  แจ้งข้อผิดพลาด (503 UPSTREAM_UNAVAILABLE พร้อม message) และไม่ส่งข้อมูลเก่าหรือรายการใด ๆ
func TestAC_EQS_08_DBDown503(t *testing.T) {
	// Given: connection pool ที่ชี้ไปฐานข้อมูลที่ไม่มีอยู่ (port 1 ไม่มีใครฟัง) แทนการหยุด container จริง
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	down, err := pgxpool.New(ctx, "postgres://salmon:salmon@127.0.0.1:1/salmon?sslmode=disable&connect_timeout=1")
	if err != nil {
		t.Fatal(err)
	}
	defer down.Close()
	h := NewHandler(NewService(NewRepository(down)))
	mux := http.NewServeMux()
	mux.Handle(PathPrefix, auth.Middleware(auth.NewMockProvider(), h.Routes()))
	srv := platform.NoStore(mux)

	// When: เปิดหน้ารายการอุปกรณ์ (และหน้าอื่นที่อ่านฐานข้อมูล)
	for _, path := range []string{PathPrefix + "models", PathPrefix + "models/1/items", PathPrefix + "items/CS-NB-0001"} {
		rec := get(t, srv, path, tokenStudentDept1)

		// Then: 503 UPSTREAM_UNAVAILABLE ตาม CON-TECH-03 มีแค่ code และ message ไม่มี items หรือข้อมูลอุปกรณ์
		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("%s status = %d ต้องเป็น 503", path, rec.Code)
			continue
		}
		raw := rec.Body.String()
		e := decode[map[string]string](t, rec)
		if e["code"] != "UPSTREAM_UNAVAILABLE" || e["message"] == "" || len(e) != 2 {
			t.Errorf("%s body ต้องเป็น {code: UPSTREAM_UNAVAILABLE, message} ได้ %s", path, raw)
		}
		if strings.Contains(raw, "items") || strings.Contains(raw, "CS-NB") {
			t.Errorf("%s ส่งข้อมูลเก่ากลับไป: %s", path, raw)
		}
		if cc := rec.Header().Get("Cache-Control"); cc != "no-store" {
			t.Errorf("%s error ต้องไม่ถูก cache (NFR-DATA-01) ได้ Cache-Control %q", path, cc)
		}
	}
}
