package audit_test

// test แยก package (audit_test) เพราะ AC-EQS-11 ต้องเรียกผ่าน handler ของ equipment ที่ import audit อยู่แล้ว
import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"salmonnorway/backend/internal/audit"
	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/equipment"
	"salmonnorway/backend/internal/platform"
)

var (
	dbOnce      sync.Once
	dbPool      *pgxpool.Pool
	dbErr       error
	dbContainer *tcpostgres.PostgresContainer
)

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
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testing.Short() {
		t.Skip("ต้องใช้ Docker: ข้ามเมื่อรันด้วย -short")
	}
	dbOnce.Do(func() {
		ctx := context.Background()
		dbContainer, dbErr = tcpostgres.Run(ctx, "postgres:16",
			tcpostgres.WithDatabase("salmon"), tcpostgres.WithUsername("salmon"), tcpostgres.WithPassword("salmon"),
			tcpostgres.BasicWaitStrategies())
		if dbErr != nil {
			return
		}
		var url string
		if url, dbErr = dbContainer.ConnectionString(ctx, "sslmode=disable"); dbErr != nil {
			return
		}
		if dbErr = platform.Migrate(url, "../../db/migrations"); dbErr != nil {
			return
		}
		if dbErr = platform.Seed(ctx, url, "../../db/seed/seed.sql"); dbErr != nil {
			return
		}
		dbPool, dbErr = pgxpool.New(ctx, url)
	})
	if dbErr != nil {
		t.Fatalf("เตรียมฐานข้อมูลไม่สำเร็จ: %v", dbErr)
	}
	return dbPool
}

// server ประกอบ handler ของ 001 แบบเดียวกับ cmd/server/main.go
func server(db equipment.Querier) http.Handler {
	h := equipment.NewHandler(equipment.NewService(equipment.NewRepository(db)))
	mux := http.NewServeMux()
	mux.Handle(equipment.PathPrefix, auth.Middleware(auth.NewMockProvider(), h.Routes()))
	return mux
}

func open(srv http.Handler, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

type row struct {
	ActorRole  string
	OccurredAt time.Time
	Outcome    string
}

// lastDetail อ่านบันทึกการเปิดรายละเอียดล่าสุดของผู้ใช้และรหัสครุภัณฑ์นั้น
func lastDetail(t *testing.T, pool *pgxpool.Pool, actorID, assetCode string) (row, bool) {
	t.Helper()
	var r row
	err := pool.QueryRow(context.Background(), `
SELECT actor_role, occurred_at, outcome FROM audit_log
WHERE actor_id = $1 AND action = 'detail' AND asset_code = $2
ORDER BY id DESC LIMIT 1`, actorID, assetCode).Scan(&r.ActorRole, &r.OccurredAt, &r.Outcome)
	if err == pgx.ErrNoRows {
		return row{}, false
	}
	if err != nil {
		t.Fatal(err)
	}
	return r, true
}

// AC-EQS-11 (DOM-AUDIT-01)
//
//	Given มีการเปิดดูรายละเอียดอุปกรณ์
//	When  การเข้าถึงเสร็จสิ้น
//	Then  มี audit log ที่ระบุผู้เข้าถึง เวลา และรหัสครุภัณฑ์ที่เปิดดู
func TestAC_EQS_11_DetailWritesAudit(t *testing.T) {
	pool := testPool(t)
	srv := server(pool)
	before := time.Now().Add(-time.Minute)

	// Given / When: นักศึกษา student-001 เปิดดูรายละเอียด CS-NB-0002 จนเสร็จ
	if rec := open(srv, "/v1/equipment/items/CS-NB-0002", "mock-student-001"); rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}

	// Then: มีแถวที่ระบุผู้เข้าถึง เวลา และรหัสครุภัณฑ์
	r, ok := lastDetail(t, pool, "student-001", "CS-NB-0002")
	if !ok {
		t.Fatal("ไม่มีบันทึกการเปิดรายละเอียดใน audit_log")
	}
	if r.ActorRole != "student" || r.Outcome != audit.OutcomeOK {
		t.Errorf("บันทึก = %+v ต้องเป็น actor_role student, outcome ok", r)
	}
	if r.OccurredAt.Before(before) || r.OccurredAt.After(time.Now().Add(time.Minute)) {
		t.Errorf("เวลาในบันทึก %s ไม่ใช่เวลาที่เปิดดู", r.OccurredAt)
	}
}

// plan.md ข้อ 6.4: บันทึกกรณีถูกปฏิเสธด้วย (outcome forbidden) เพื่อตรวจสอบย้อนหลัง
func TestAudit_ForbiddenDetailIsRecorded(t *testing.T) {
	pool := testPool(t)
	srv := server(pool)

	// นักศึกษาภาควิชา 2 เปิดอุปกรณ์ภาควิชา 1 (FR-EQS-10)
	if rec := open(srv, "/v1/equipment/items/CS-NB-0003", "mock-student-101"); rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d ต้องเป็น 403", rec.Code)
	}
	r, ok := lastDetail(t, pool, "student-101", "CS-NB-0003")
	if !ok || r.Outcome != audit.OutcomeForbidden {
		t.Errorf("ต้องมีบันทึก outcome forbidden ได้ %+v (มี %v)", r, ok)
	}
}

// failAuditDB ส่งทุก query ไปฐานข้อมูลจริง ยกเว้นการเขียน audit_log ที่จำลองว่าล้มเหลว
type failAuditDB struct{ *pgxpool.Pool }

type errRow struct{ err error }

func (e errRow) Scan(...any) error { return e.err }

func (f failAuditDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if strings.Contains(sql, "INSERT INTO audit_log") {
		return errRow{err: context.DeadlineExceeded}
	}
	return f.Pool.QueryRow(ctx, sql, args...)
}

// plan.md ข้อ 6.4: เขียน audit ไม่สำเร็จ ไม่ส่งข้อมูลอุปกรณ์ ตอบ 503 (fail closed รอยืนยัน Q-P4 ดู T-29)
func TestAudit_WriteFailsThenNoData(t *testing.T) {
	srv := server(failAuditDB{testPool(t)})

	rec := open(srv, "/v1/equipment/items/CS-NB-0001", "mock-student-001")
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d ต้องเป็น 503", rec.Code)
	}
	if body := rec.Body.String(); strings.Contains(body, "CS-NB-0001") || strings.Contains(body, "available") {
		t.Errorf("เขียน audit ไม่สำเร็จแต่ยังส่งข้อมูลอุปกรณ์: %s", body)
	}
}

// Write ใช้ได้ทั้ง pool และ transaction (ตกลงกับ 003): ใน transaction ที่ยกเลิก บันทึกต้องหายไปด้วย
func TestWrite_WorksInTransaction(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	id, err := audit.Write(ctx, tx, audit.Entry{
		ActorID: "student-tx", ActorRole: "student", Action: audit.ActionSearch,
		Query: map[string]any{"status": []string{"available"}}, Outcome: audit.OutcomeOK,
	})
	if err != nil || id == 0 {
		t.Fatalf("เขียนใน transaction ไม่สำเร็จ: id=%d err=%v", id, err)
	}
	_ = tx.Rollback(ctx)

	var n int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM audit_log WHERE actor_id = 'student-tx'`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("transaction ยกเลิกแล้วแต่ยังมีบันทึก %d แถว", n)
	}

	// เขียนด้วย pool ตรง ๆ: query เก็บเป็น JSON และช่องที่ว่างเป็น NULL
	if _, err := audit.Write(ctx, pool, audit.Entry{ActorID: "student-pool", ActorRole: "student", Action: audit.ActionList, Outcome: audit.OutcomeOK}); err != nil {
		t.Fatal(err)
	}
	var assetIsNull, queryIsNull bool
	err = pool.QueryRow(ctx, `SELECT asset_code IS NULL, query IS NULL FROM audit_log WHERE actor_id = 'student-pool'`).Scan(&assetIsNull, &queryIsNull)
	if err != nil || !assetIsNull || !queryIsNull {
		t.Errorf("ช่องที่ไม่ได้ส่งต้องเป็น NULL ได้ asset_code null=%v query null=%v err=%v", assetIsNull, queryIsNull, err)
	}
}
