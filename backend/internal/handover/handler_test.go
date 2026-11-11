package handover

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
