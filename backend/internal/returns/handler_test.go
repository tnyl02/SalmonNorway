package returns

import (
	"context"
	"net/http"
	"net/http/httptest"
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
