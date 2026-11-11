package cancelrequest

// ตัวช่วยของ test จาก specs/003-cancel-request/test-cases.md (ไม่ใช่ test เอง แก้ได้ถ้าโครงโค้ดเปลี่ยน เช่น NewService รับ Notifier ใน T-27)
// ใช้รูปแบบเดียวกับ backend/internal/equipment/handler_test.go ของ 001: PostgreSQL 16 ใน testcontainers + migration + seed ของ 001
//
// สิ่งที่ test เหล่านี้คาดจากโค้ดของ 003 (ตาม plan.md ข้อ 4 และ tasks.md T-04, T-06):
//   - NewRepository(*pgxpool.Pool), NewService(*Repository), NewHandler(*Service), (*Handler).Routes() http.Handler
//   - Routes() ผูก "GET /v1/my/pending-requests" และ "POST /v1/borrow-requests/{id}/cancel"
//     (หลังบ้านไม่มี /api นำหน้า เหมือน equipment.PathPrefix = "/v1/equipment/" เพราะ vite proxy ตัด /api ออก)
//   - ตาราง borrow_request ตาม plan.md ข้อ 3 (สร้างใน T-01)
// ก่อน T-01, T-04, T-06 เสร็จ test ในโฟลเดอร์นี้จะ compile ไม่ผ่าน เป็นเรื่องปกติ ไม่ใช่บั๊ก

import (
	"context"
	"encoding/json"
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
	tcMigrationsDir = "../../db/migrations"
	tcSeedFile      = "../../db/seed/seed.sql"

	// ผู้ใช้จำลองจาก auth.NewMockProvider (A = student-001, B = student-002 ตามหัวตารางใน test-cases.md)
	tcTokenA     = "mock-student-001"
	tcTokenStaff = "mock-staff-001"
	tcUserA      = "student-001"
	tcUserB      = "student-002"

	tcPathMine = "/v1/my/pending-requests"

	// อุปกรณ์จาก seed.sql ของ 001
	tcItemAvailable = 1 // CS-NB-0001 Notebook ตัวอย่างรุ่น A สถานะ available
	tcItemBorrowed  = 4 // CS-NB-0004 Notebook ตัวอย่างรุ่น A สถานะ borrowed
	tcItemOther     = 6 // CS-PJ-0001 Projector ตัวอย่างรุ่น B สถานะ available
)

var (
	tcDBOnce      sync.Once
	tcDBPool      *pgxpool.Pool
	tcDBErr       error
	tcDBContainer *tcpostgres.PostgresContainer
)

// TestMain ปิดฐานข้อมูลหลัง test ทั้ง package จบ (package นี้มี TestMain ได้ตัวเดียว test อื่นของ 003 ใช้ตัวนี้ร่วม)
func TestMain(m *testing.M) {
	code := m.Run()
	if tcDBPool != nil {
		tcDBPool.Close()
	}
	if tcDBContainer != nil {
		_ = testcontainers.TerminateContainer(tcDBContainer)
	}
	os.Exit(code)
}

// tcPool เปิด PostgreSQL ครั้งเดียวต่อ package แล้วรัน migration และ seed ของ 001
func tcPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testing.Short() {
		t.Skip("ต้องใช้ Docker: ข้ามเมื่อรันด้วย -short")
	}
	tcDBOnce.Do(func() {
		ctx := context.Background()
		tcDBContainer, tcDBErr = tcpostgres.Run(ctx, "postgres:16",
			tcpostgres.WithDatabase("salmon"),
			tcpostgres.WithUsername("salmon"),
			tcpostgres.WithPassword("salmon"),
			tcpostgres.BasicWaitStrategies(),
		)
		if tcDBErr != nil {
			return
		}
		var url string
		if url, tcDBErr = tcDBContainer.ConnectionString(ctx, "sslmode=disable"); tcDBErr != nil {
			return
		}
		if tcDBErr = platform.Migrate(url, tcMigrationsDir); tcDBErr != nil {
			return
		}
		if tcDBErr = platform.Seed(ctx, url, tcSeedFile); tcDBErr != nil {
			return
		}
		tcDBPool, tcDBErr = pgxpool.New(ctx, url)
	})
	if tcDBErr != nil {
		t.Fatalf("เตรียมฐานข้อมูลไม่สำเร็จ: %v", tcDBErr)
	}
	return tcDBPool
}

// tcServer ล้างคำขอทั้งหมด แล้วประกอบ handler แบบเดียวกับ cmd/server/main.go (หลัง auth.Middleware ห่อด้วย platform.NoStore)
// ล้างด้วย DELETE เพราะ audit_log เป็น append-only (TRUNCATE ... CASCADE จะโดน trigger)
// ถ้า 002 เพิ่ม foreign key จาก audit_log.request_id มาที่ borrow_request ต้องกลับมาแก้วิธีล้างตรงนี้
func tcServer(t *testing.T) http.Handler {
	t.Helper()
	pool := tcPool(t)
	if _, err := pool.Exec(context.Background(), `DELETE FROM borrow_request`); err != nil {
		t.Fatalf("ล้างตาราง borrow_request ไม่สำเร็จ (สร้างใน T-01): %v", err)
	}
	h := NewHandler(NewService(NewRepository(pool)))
	routes := auth.Middleware(auth.NewMockProvider(), h.Routes())
	mux := http.NewServeMux()
	mux.Handle("/v1/my/", routes)
	mux.Handle("/v1/borrow-requests/", routes)
	return platform.NoStore(mux)
}

// tcInsert ใส่คำขอยืม 1 คำขอ คืนรหัสคำขอ (คอลัมน์ตาม plan.md ข้อ 3)
func tcInsert(t *testing.T, requester string, itemID int, loanType, status string, submittedAt time.Time) int {
	t.Helper()
	var id int
	err := tcPool(t).QueryRow(context.Background(),
		`INSERT INTO borrow_request (requester_id, item_id, loan_type, purpose, status, submitted_at)
		 VALUES ($1, $2, $3, 'เพื่อการเรียน', $4, $5) RETURNING id`,
		requester, itemID, loanType, status, submittedAt).Scan(&id)
	if err != nil {
		t.Fatalf("ใส่คำขอไม่สำเร็จ: %v", err)
	}
	return id
}

// tcRow อ่านสถานะ เหตุผล และเวลาที่ยกเลิกของคำขอจากฐานข้อมูล
type tcRowState struct {
	Status       string
	CancelReason *string
	CancelledAt  *time.Time
}

func tcRow(t *testing.T, id int) tcRowState {
	t.Helper()
	var r tcRowState
	err := tcPool(t).QueryRow(context.Background(),
		`SELECT status, cancel_reason, cancelled_at FROM borrow_request WHERE id = $1`, id).
		Scan(&r.Status, &r.CancelReason, &r.CancelledAt)
	if err != nil {
		t.Fatalf("อ่านคำขอ %d ไม่สำเร็จ: %v", id, err)
	}
	return r
}

// tcItemStatus อ่านสถานะอุปกรณ์รายชิ้น (CON-DATA-01)
func tcItemStatus(t *testing.T, itemID int) string {
	t.Helper()
	var s string
	if err := tcPool(t).QueryRow(context.Background(),
		`SELECT status FROM equipment_item WHERE id = $1`, itemID).Scan(&s); err != nil {
		t.Fatalf("อ่านสถานะอุปกรณ์ %d ไม่สำเร็จ: %v", itemID, err)
	}
	return s
}

// tcGet เรียก GET ด้วย token (ว่าง = ไม่ส่ง token)
func tcGet(t *testing.T, srv http.Handler, path, token string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

// tcCancel เรียก POST /v1/borrow-requests/{id}/cancel พร้อมเหตุผล
func tcCancel(t *testing.T, srv http.Handler, id any, reason, token string) *httptest.ResponseRecorder {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"reason": reason})
	req := httptest.NewRequest(http.MethodPost, "/v1/borrow-requests/"+toString(id)+"/cancel", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func toString(v any) string {
	b, _ := json.Marshal(v)
	return strings.Trim(string(b), `"`)
}

// รูปแบบ JSON ตาม plan.md ข้อ 4 (ตั้งชื่อขึ้นต้น tc เพื่อไม่ชนกับ DTO ของโค้ดจริง)
type tcListItem struct {
	RequestID   int    `json:"requestId"`
	ModelName   string `json:"modelName"`
	AssetCode   string `json:"assetCode"`
	SubmittedAt string `json:"submittedAt"`
}

type tcListBody struct {
	Items []tcListItem `json:"items"`
}

type tcCancelBody struct {
	RequestID int    `json:"requestId"`
	Status    string `json:"status"`
}

type tcErrorBody struct {
	Code string `json:"code"`
}

func tcDecode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v (body: %s)", err, rec.Body.String())
	}
	return v
}

func tcWantStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("HTTP status = %d ต้องเป็น %d (body: %s)", rec.Code, want, rec.Body.String())
	}
}
