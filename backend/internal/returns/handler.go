package returns

import (
	"encoding/json"
	"net/http"

	"salmonnorway/backend/internal/auth"
)

// Handler รวม endpoint ของการรับคืนตาม plan.md ข้อ 4 (FR-RET-01 ถึง FR-RET-07)
type Handler struct {
	repo *Repository
}

// NewHandler สร้าง Handler ของฟีเจอร์รับคืน
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Routes คืนฟังก์ชันที่ผูก endpoint ทั้งหมดของ 006 เข้ากับ router
// ทุก endpoint ต้องผ่าน auth.Middleware ก่อนเสมอ (IF-IDP-01, NFR-SEC-01)
func (h *Handler) Routes(provider auth.Provider) func(*http.ServeMux) {
	return func(mux *http.ServeMux) {
		for _, pattern := range []string{
			"GET /api/v1/returns/open-loans",                      // FR-RET-01(ข), FR-RET-06 ทำใน T-17
			"GET /api/v1/returns/lookup",                          // FR-RET-01(ก) ทำใน T-05
			"POST /api/v1/returns",                                // FR-RET-02 ถึง 07 ทำใน T-04
			"GET /api/v1/returns/loans/{loanId}",                  // FR-RET-06 ทำใน T-14
			"GET /api/v1/returns/loans/{loanId}/photos/{photoId}", // FR-RET-05 ทำใน T-14
		} {
			mux.Handle(pattern, auth.Middleware(provider, http.HandlerFunc(notImplemented)))
		}
	}
}

// notImplemented ตอบ 501 ให้ endpoint ที่ยังไม่ถึง task ของมัน (T-03 โครงเท่านั้น)
func notImplemented(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"code": "NOT_IMPLEMENTED"})
}

// writeJSON ตอบ JSON และตั้ง Cache-Control: no-store ทุก response ของ 006 (plan.md ข้อ 4, NFR-DATA-01)
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
