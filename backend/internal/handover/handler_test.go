package handover

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/auth"
)

// ผู้ใช้จำลองจาก auth.NewMockProvider
const (
	tokenStaff    = "mock-staff-001"
	tokenStudentA = "mock-student-001"
	tokenStudentB = "mock-student-101"
)

// routesOnly คือ router ที่ผ่าน auth.Middleware เหมือนใน main.go (ใช้กับ test ที่ไม่แตะฐานข้อมูล)
func routesOnly(db DB) http.Handler {
	provider := auth.NewMockProvider()
	return auth.Middleware(provider, NewHandler(db, provider).Routes())
}

func doRequest(h http.Handler, method, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

var allRoutes = []struct{ method, path string }{
	{http.MethodGet, PathPrefix + "requests"},
	{http.MethodPost, PathPrefix + "requests/1/verify-borrower"},
	{http.MethodPost, PathPrefix + "requests/1/handover"},
}

// T-03 เสร็จเมื่อ: ทั้ง 3 route ตอบ 403 {code, message} FORBIDDEN_ROLE เมื่อเป็นนักศึกษา และ 401 เมื่อไม่มี token (CON-HND-01, IF-IDP-01, CON-TECH-03)
func TestHandover_NonStaffForbidden(t *testing.T) {
	h := routesOnly(nil)
	for _, rt := range allRoutes {
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			rec := doRequest(h, rt.method, rt.path, tokenStudentA)
			if rec.Code != http.StatusForbidden {
				t.Fatalf("นักศึกษา: HTTP %d ต้องเป็น 403", rec.Code)
			}
			var body ErrorDTO
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body ไม่ใช่ JSON: %s", rec.Body.String())
			}
			if body.Code != "FORBIDDEN_ROLE" || body.Message == "" {
				t.Errorf("body = %+v ต้องเป็น FORBIDDEN_ROLE พร้อม message", body)
			}

			if rec := doRequest(h, rt.method, rt.path, ""); rec.Code != http.StatusUnauthorized {
				t.Errorf("ไม่มี token: HTTP %d ต้องเป็น 401", rec.Code)
			}
		})
	}
}

// testPool คืน connection pool ไปยังฐานข้อมูลของ test (ใช้กับ test ผ่าน HTTP ที่ต้อง commit จริง)
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	testDB(t)
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close) // ลงทะเบียนก่อน จึงปิดหลังการลบข้อมูลทดสอบ (Cleanup ทำงานย้อนลำดับ)
	return pool
}

// createItem เพิ่มอุปกรณ์รายชิ้นของรุ่น 1 แบบ commit จริง และลบทุกอย่างที่ผูกกับชิ้นนี้ตอน test จบ
func createItem(t *testing.T, pool *pgxpool.Pool, assetCode, status string) int {
	t.Helper()
	ctx := context.Background()
	var id int
	if err := pool.QueryRow(ctx,
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, $2) RETURNING id`,
		assetCode, status).Scan(&id); err != nil {
		t.Fatalf("เพิ่มอุปกรณ์ไม่สำเร็จ: %v", err)
	}
	t.Cleanup(func() {
		for _, q := range []string{
			`DELETE FROM handover_photo WHERE loan_id IN (SELECT id FROM loan WHERE item_id = $1)`,
			`DELETE FROM handover WHERE loan_id IN (SELECT id FROM loan WHERE item_id = $1)`,
			`DELETE FROM audit_log WHERE request_id IN (SELECT id FROM borrow_request WHERE item_id = $1)`,
			`DELETE FROM loan WHERE item_id = $1`,
			`DELETE FROM borrow_request WHERE item_id = $1`,
			`DELETE FROM equipment_item WHERE id = $1`,
		} {
			if _, err := pool.Exec(ctx, q, id); err != nil {
				t.Errorf("ลบข้อมูลทดสอบไม่สำเร็จ: %v", err)
			}
		}
	})
	return id
}

// createRequest เพิ่มคำขอยืมแบบปกติแบบ commit จริง คืน id (ข้อมูลลบโดย cleanup ของ createItem)
func createRequest(t *testing.T, pool *pgxpool.Pool, requesterID string, itemID int, status string) int {
	t.Helper()
	var id int
	if err := pool.QueryRow(context.Background(),
		`INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
		 VALUES ($1, 'ผู้ขอทดสอบ', $1, $2, 'normal', 'study', now(), now(), $3) RETURNING id`,
		requesterID, itemID, status).Scan(&id); err != nil {
		t.Fatalf("เพิ่มคำขอยืมไม่สำเร็จ: %v", err)
	}
	return id
}

// AC-HND-14 (FR-HND-01)
// Given มีคำขอที่ได้รับอนุมัติแล้วและยังไม่ส่งมอบ 2 รายการ (อุปกรณ์คนละชิ้น) คำขอที่ยังไม่อนุมัติ 1 รายการ และคำขอที่ส่งมอบแล้ว 1 รายการ
// When  เจ้าหน้าที่เปิดรายการคำขอเพื่อส่งมอบ
// Then  รายการแสดงเฉพาะคำขอที่ได้รับอนุมัติแล้วและยังไม่ส่งมอบทั้ง 2 รายการ
func TestAC_HND_14_ListApprovedOnly(t *testing.T) {
	pool := testPool(t)
	approved1 := createRequest(t, pool, "student-001", createItem(t, pool, "T14-AP-0001", "available"), "approved")
	approved2 := createRequest(t, pool, "student-002", createItem(t, pool, "T14-AP-0002", "available"), "approved")
	createRequest(t, pool, "student-001", createItem(t, pool, "T14-RQ-0001", "available"), "requested")
	createRequest(t, pool, "student-001", createItem(t, pool, "T14-HO-0001", "borrowed"), "handed_over")

	rec := doRequest(routesOnly(pool), http.MethodGet, PathPrefix+"requests", tokenStaff)
	if rec.Code != http.StatusOK {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
	}
	var got []ReadyRequestDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("body ไม่ใช่ JSON array: %s", rec.Body.String())
	}
	want := map[int]string{approved1: "T14-AP-0001", approved2: "T14-AP-0002"}
	if len(got) != len(want) {
		t.Fatalf("ได้ %d รายการ ต้องเป็น 2: %+v", len(got), got)
	}
	for _, d := range got {
		if want[d.RequestID] != d.AssetCode {
			t.Errorf("รายการที่ไม่คาดคิด: %+v", d)
		}
		if d.ModelName == "" || d.RequesterStudentID == "" || d.SubmittedAt.IsZero() {
			t.Errorf("ข้อมูลไม่ครบ: %+v", d)
		}
	}
}
