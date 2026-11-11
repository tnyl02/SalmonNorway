package approve

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"salmonnorway/backend/internal/config"
)

// PathPrefix คือ path ของทุก endpoint ของฟีเจอร์ 004 (plan.md ข้อ 4)
// หน้าจอเรียก /api/v1/approvals/... ส่วน /api ถูกตัดโดย proxy ใน frontend/vite.config.js (แบบเดียวกับ 002)
const PathPrefix = "/v1/approvals/"

// Handler รับ HTTP request ของฟีเจอร์ 004
// ต้องอยู่หลัง auth.Middleware (IF-IDP-01, NFR-SEC-01) และให้บริการผ่าน HTTPS เท่านั้น (CON-TECH-02)
type Handler struct {
	svc *Service
}

// NewHandler สร้าง Handler
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes คืน router ของฟีเจอร์นี้ ให้ main.go ห่อด้วย auth.Middleware แล้วผูกกับ PathPrefix
// ทุก endpoint ใต้ PathPrefix ตรวจว่าเป็นเจ้าหน้าที่ก่อน (ASM-01, AC-APV-05 FORBIDDEN_ROLE)
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+PathPrefix+"requests", h.list)
	return requireStaff(mux)
}

// requireStaff ปฏิเสธผู้ที่ไม่ใช่เจ้าหน้าที่ก่อนถึง endpoint ใด ๆ ของ 004 (AC-APV-05, CON-TECH-03)
func requireStaff(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := staff(r.Context()); err != nil {
			writeServiceError(w, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// list: GET /v1/approvals/requests รายการคำขอยืมปกติที่รอพิจารณาของภาควิชาตนเอง (FR-APV-01 ก, AC-APV-11)
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	items, loadedAt, err := h.svc.ListPending(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	out := PendingListDTO{Items: make([]PendingRequestDTO, 0, len(items)), LoadedAt: loadedAt}
	for _, it := range items {
		out.Items = append(out.Items, toPendingDTO(it))
	}
	writeJSON(w, http.StatusOK, out)
}

// writeServiceError แปลง error ของ service เป็น {code, message} ตาม CON-TECH-03 ข้อความจาก config (FR-APV-01 ซ)
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNoPrincipal):
		// กันพลาดเมื่อผูก router ผิด ข้อความใช้แบบ auth.Middleware ไม่มี message ใน spec ของ 004 (ดู plan.md ข้อ 8)
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED")
	case errors.Is(err, ErrForbiddenRole):
		writeError(w, http.StatusForbidden, "FORBIDDEN_ROLE")
	default:
		log.Printf("approve: %v", err)
		writeError(w, http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE")
	}
}

// writeError ส่ง error ด้วยรหัสจากรายการกลางและข้อความจาก config.ApprovalMessages (CON-TECH-03)
func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, ErrorDTO{Code: code, Message: config.ApprovalMessages[code]})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
