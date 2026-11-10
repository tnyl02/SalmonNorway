package equipment

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"

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

// newTestServer ประกอบ handler แบบเดียวกับ cmd/server/main.go (อยู่หลัง auth.Middleware)
func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	h := NewHandler(NewService(NewRepository(testPool(t))))
	mux := http.NewServeMux()
	mux.Handle(PathPrefix, auth.Middleware(auth.NewMockProvider(), h.Routes()))
	return mux
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
