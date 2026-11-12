package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeProvider คืนผู้ใช้ที่กำหนดไว้ ใช้ทดสอบกรณีที่ผู้ให้บริการส่งข้อมูลไม่ครบ
type fakeProvider struct{ p Principal }

func (f fakeProvider) Authenticate(context.Context, string) (Principal, error) { return f.p, nil }

// serve ส่ง request ผ่าน middleware แล้วบอกว่า handler ด้านในถูกเรียกไหม และได้ผู้ใช้คนไหน
func serve(t *testing.T, p Provider, authHeader string) (*httptest.ResponseRecorder, bool, Principal) {
	t.Helper()
	called := false
	var got Principal
	h := Middleware(p, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		got, _ = FromContext(r.Context())
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/equipment/models", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec, called, got
}

// errorCode อ่าน body ของ error แล้วตรวจว่ามี message ภาษาไทย (CON-TECH-03, T-37)
func errorCode(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body struct{ Code, Message string }
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("body ไม่ใช่ JSON: %v", err)
	}
	if body.Message == "" {
		t.Errorf("error %s ต้องมี message (CON-TECH-03)", body.Code)
	}
	return body.Code
}

// T-03 เสร็จเมื่อ: ไม่มี token -> 401 และไม่เข้าถึงข้อมูล (IF-IDP-01, NFR-SEC-01)
func TestAuth_NoToken_401(t *testing.T) {
	for name, header := range map[string]string{
		"ไม่มี header":       "",
		"ไม่ใช่ Bearer":      "Basic abc",
		"token ว่าง":         "Bearer ",
		"token ที่ไม่รู้จัก": "Bearer mock-unknown",
	} {
		t.Run(name, func(t *testing.T) {
			rec, called, _ := serve(t, NewMockProvider(), header)
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want 401", rec.Code)
			}
			if code := errorCode(t, rec); code != "UNAUTHENTICATED" {
				t.Errorf("code = %q, want UNAUTHENTICATED", code)
			}
			if called {
				t.Error("handler ด้านในต้องไม่ถูกเรียก")
			}
		})
	}
}

// T-03 เสร็จเมื่อ: ไม่ได้ภาควิชา = ปฏิเสธ ไม่เดา (plan 6.1) ทีมเลือก 403
func TestAuth_MissingDepartment_Rejected(t *testing.T) {
	cases := map[string]Principal{
		"ไม่มีภาควิชา":    {UserID: "student-x", Role: RoleStudent},
		"บทบาทไม่รู้จัก":  {UserID: "lecturer-x", Role: "lecturer", DepartmentID: 1},
		"ไม่มีรหัสผู้ใช้": {Role: RoleStaff, DepartmentID: 1},
	}
	for name, p := range cases {
		t.Run(name, func(t *testing.T) {
			rec, called, _ := serve(t, fakeProvider{p}, "Bearer any")
			if rec.Code != http.StatusForbidden {
				t.Errorf("status = %d, want 403", rec.Code)
			}
			if code := errorCode(t, rec); code != "PRINCIPAL_INCOMPLETE" {
				t.Errorf("code = %q, want PRINCIPAL_INCOMPLETE", code)
			}
			if called {
				t.Error("handler ด้านในต้องไม่ถูกเรียก")
			}
		})
	}
}

// ผู้ใช้ครบ -> ผ่านไปถึง handler พร้อมข้อมูลผู้ใช้ใน context
func TestAuth_ValidToken_PassesPrincipal(t *testing.T) {
	rec, called, got := serve(t, NewMockProvider(), "Bearer mock-student-101")
	if !called || rec.Code != http.StatusOK {
		t.Fatalf("status = %d, called = %v, want handler ถูกเรียก", rec.Code, called)
	}
	want := Principal{UserID: "student-101", Role: RoleStudent, DepartmentID: 2}
	if got != want {
		t.Errorf("principal = %+v, want %+v", got, want)
	}
}
