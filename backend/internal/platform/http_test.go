package platform

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
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
