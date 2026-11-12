package platform

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

// T-01 เสร็จเมื่อ: GET /healthz ตอบ 200 เมื่อต่อฐานข้อมูลได้
func TestHealthz_OK_WhenDBReachable(t *testing.T) {
	rec := httptest.NewRecorder()
	NewRouter(fakePinger{}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

// ถ้าต่อฐานข้อมูลไม่ได้ ต้องไม่ตอบ 200
func TestHealthz_503_WhenDBDown(t *testing.T) {
	rec := httptest.NewRecorder()
	NewRouter(fakePinger{err: errors.New("down")}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want 503", rec.Code)
	}
}

// T-40 เสร็จเมื่อ: 405 ได้ body {code: METHOD_NOT_ALLOWED, message} ไม่มีช่องอื่น header Allow ยังอยู่
// และ 200, 401, error JSON อื่นไม่ถูกแตะ (CON-TECH-03, CON-DATA-02)
func TestCON_TECH_03_MethodNotAllowed(t *testing.T) {
	// router แบบเดียวกับของฟีเจอร์: ประกาศ method ใน pattern แล้วให้ http.ServeMux ตอบ 405 เอง
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/x", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /v1/unauth", func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusUnauthorized, CodeUnauthenticated)
	})
	mux.HandleFunc("GET /v1/down", func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, http.StatusServiceUnavailable, CodeUpstreamUnavailable)
	})
	mux.HandleFunc("GET /v1/plain", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello")) // Write โดยไม่เรียก WriteHeader = 200
	})
	srv := MethodNotAllowedJSON(mux)

	serve := func(method, path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, httptest.NewRequest(method, path, nil))
		return rec
	}
	decode := func(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
		t.Helper()
		var body map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("body ไม่ใช่ JSON: %v", err)
		}
		return body
	}

	t.Run("405 ได้ {code, message} และ header Allow", func(t *testing.T) {
		rec := serve(http.MethodPost, "/v1/x")
		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, ต้องเป็น 405", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
			t.Errorf("Content-Type = %q ต้องเป็น JSON", ct)
		}
		if allow := rec.Header().Get("Allow"); !strings.Contains(allow, http.MethodGet) {
			t.Errorf("header Allow = %q ต้องยังบอก GET", allow)
		}
		raw := rec.Body.String()
		if strings.Contains(raw, "Method Not Allowed") {
			t.Errorf("ข้อความเดิมของ router ยังหลุดมา: %q", raw)
		}
		body := decode(t, rec)
		if len(body) != 2 || body["code"] != CodeMethodNotAllowed || body["message"] != "ไม่รองรับการเรียกแบบนี้" {
			t.Errorf("body = %v ต้องเป็น {code: METHOD_NOT_ALLOWED, message: ไม่รองรับการเรียกแบบนี้}", body)
		}
	})

	t.Run("200 ไม่ถูกแตะ", func(t *testing.T) {
		rec := serve(http.MethodGet, "/v1/x")
		if rec.Code != http.StatusOK || decode(t, rec)["status"] != "ok" {
			t.Errorf("status = %d body = %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("200 แบบ Write อย่างเดียวไม่ถูกแตะ", func(t *testing.T) {
		rec := serve(http.MethodGet, "/v1/plain")
		if rec.Code != http.StatusOK || rec.Body.String() != "hello" {
			t.Errorf("status = %d body = %q", rec.Code, rec.Body.String())
		}
	})

	for _, c := range []struct {
		path   string
		status int
		code   string
	}{
		{"/v1/unauth", http.StatusUnauthorized, CodeUnauthenticated},
		{"/v1/down", http.StatusServiceUnavailable, CodeUpstreamUnavailable},
	} {
		t.Run(c.code+" ไม่ถูกแตะ", func(t *testing.T) {
			rec := serve(http.MethodGet, c.path)
			if rec.Code != c.status {
				t.Fatalf("status = %d, ต้องเป็น %d", rec.Code, c.status)
			}
			if body := decode(t, rec); body["code"] != c.code || body["message"] != ErrorMessage(c.code) {
				t.Errorf("body = %v", body)
			}
		})
	}
}
