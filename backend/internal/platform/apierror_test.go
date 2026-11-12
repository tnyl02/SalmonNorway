package platform_test

// test แยก package (platform_test) เพราะต้องเรียก auth และ equipment ที่ import platform อยู่แล้ว
import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/equipment"
	"salmonnorway/backend/internal/platform"
)

// incompleteProvider ยืนยันตัวตนได้ แต่ไม่มีภาควิชา (กรณี PRINCIPAL_INCOMPLETE)
type incompleteProvider struct{}

func (incompleteProvider) Authenticate(context.Context, string) (auth.Principal, error) {
	return auth.Principal{UserID: "student-x", Role: auth.RoleStudent}, nil
}

// downDB จำลองฐานข้อมูลที่ไม่ตอบ ทุก query คืน error (กรณี UPSTREAM_UNAVAILABLE)
type downDB struct{}

var errDown = errors.New("database is down")

func (downDB) Query(context.Context, string, ...any) (pgx.Rows, error) { return nil, errDown }
func (downDB) QueryRow(context.Context, string, ...any) pgx.Row        { return errRow{} }

type errRow struct{}

func (errRow) Scan(...any) error { return errDown }

// T-37 เสร็จเมื่อ: error ทุกตัวที่ 001 และ auth ส่ง มีทั้ง code และ message และไม่มีช่องอื่น (CON-TECH-03, CR-004)
func TestCON_TECH_03_ErrorBodyHasCodeAndMessage(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	eq := equipment.NewHandler(equipment.NewService(equipment.NewRepository(downDB{})))

	cases := []struct {
		name       string
		handler    http.Handler
		path       string
		token      string
		wantStatus int
		wantCode   string
	}{
		{"401 ไม่มี token", auth.Middleware(auth.NewMockProvider(), ok), "/v1/equipment/models", "", http.StatusUnauthorized, "UNAUTHENTICATED"},
		{"401 token ไม่รู้จัก", auth.Middleware(auth.NewMockProvider(), ok), "/v1/equipment/models", "mock-unknown", http.StatusUnauthorized, "UNAUTHENTICATED"},
		{"403 ไม่มีภาควิชา", auth.Middleware(incompleteProvider{}, ok), "/v1/equipment/models", "any", http.StatusForbidden, "PRINCIPAL_INCOMPLETE"},
		{"503 ฐานข้อมูลไม่ตอบ (รายการรุ่น)", auth.Middleware(auth.NewMockProvider(), eq.Routes()), "/v1/equipment/models", "mock-student-001", http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE"},
		{"503 ฐานข้อมูลไม่ตอบ (รายละเอียด)", auth.Middleware(auth.NewMockProvider(), eq.Routes()), "/v1/equipment/items/CS-NB-0001", "mock-student-001", http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, c.path, nil)
			if c.token != "" {
				req.Header.Set("Authorization", "Bearer "+c.token)
			}
			rec := httptest.NewRecorder()
			c.handler.ServeHTTP(rec, req)

			if rec.Code != c.wantStatus {
				t.Fatalf("status = %d, ต้องเป็น %d", rec.Code, c.wantStatus)
			}
			var body map[string]any
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("body ไม่ใช่ JSON: %v", err)
			}
			if len(body) != 2 {
				t.Errorf("body ต้องมีแค่ code และ message ได้ %v", body)
			}
			if body["code"] != c.wantCode {
				t.Errorf("code = %v, ต้องเป็น %s", body["code"], c.wantCode)
			}
			if msg, _ := body["message"].(string); msg == "" || msg == c.wantCode {
				t.Errorf("message ต้องเป็นข้อความภาษาไทยที่ไม่ว่าง ได้ %q", msg)
			}
		})
	}
}

// ทุกรหัสที่ของกลางและ 001 ใช้ มีข้อความของตัวเอง ไม่ตกไปใช้รหัสแทนข้อความ (CON-TECH-03)
func TestErrorMessage_AllCodesHaveMessage(t *testing.T) {
	for _, code := range []string{
		platform.CodeUnauthenticated, platform.CodePrincipalIncomplete,
		platform.CodeForbiddenDepartment, platform.CodeUpstreamUnavailable,
	} {
		if m := platform.ErrorMessage(code); m == "" || m == code {
			t.Errorf("%s ไม่มีข้อความ", code)
		}
	}
	if m := platform.ErrorMessage("NOT_REGISTERED"); m != "NOT_REGISTERED" {
		t.Errorf("รหัสที่ไม่มีข้อความต้องคืนรหัสเอง (message ไม่ว่าง) ได้ %q", m)
	}
}
