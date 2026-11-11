package cancelrequest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"salmonnorway/backend/internal/auth"
)

// path ที่หลังบ้านรับ plan.md ข้อ 4 เขียนเป็น /api/v1/... ส่วน /api ถูกตัดโดย proxy ใน frontend/vite.config.js (แบบเดียวกับ 001)
const (
	PathPrefixMine     = "/v1/my/"              // GET /v1/my/pending-requests
	PathPrefixRequests = "/v1/borrow-requests/" // POST /v1/borrow-requests/{id}/cancel (T-06)
	pathPendingList    = PathPrefixMine + "pending-requests"
)

// Handler รับ HTTP request ของฟีเจอร์ 003 ทุก endpoint ต้องอยู่หลัง auth.Middleware (IF-IDP-01, NFR-SEC-01)
type Handler struct {
	svc *Service
}

// NewHandler สร้าง Handler
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes คืน router ของฟีเจอร์นี้ ให้ main.go ห่อด้วย auth.Middleware แล้วผูกกับ PathPrefixMine และ PathPrefixRequests
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+pathPendingList, h.listPending)
	mux.HandleFunc("POST "+PathPrefixRequests+"{id}/cancel", h.cancel)
	return mux
}

// listPending: GET /v1/my/pending-requests คำขอที่ยังยกเลิกได้ของนักศึกษาที่เข้าสู่ระบบ (FR-CXL-06, FR-CXL-03, AC-CXL-08)
func (h *Handler) listPending(w http.ResponseWriter, r *http.Request) {
	p, ok := auth.FromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorDTO{Code: "UNAUTHENTICATED"}) // ไม่ควรเกิดถ้าอยู่หลัง middleware
		return
	}
	items, err := h.svc.ListMyPending(r.Context(), p)
	switch {
	case errors.Is(err, ErrStudentOnly):
		writeJSON(w, http.StatusForbidden, ErrorDTO{Code: "STUDENT_ONLY"})
	case err != nil:
		writeUpstreamError(w, err)
	default:
		writeJSON(w, http.StatusOK, PendingListDTO{Items: items})
	}
}

// cancel: POST /v1/borrow-requests/{id}/cancel ยกเลิกคำขอพร้อมเหตุผล (FR-CXL-01, FR-CXL-04, AC-CXL-01)
func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	p, ok := auth.FromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, ErrorDTO{Code: "UNAUTHENTICATED"}) // ไม่ควรเกิดถ้าอยู่หลัง middleware
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusNotFound, ErrorDTO{Code: "REQUEST_NOT_FOUND"}) // รหัสที่ไม่ใช่ตัวเลข ถือว่าไม่มีคำขอนี้ (FR-CXL-03)
		return
	}
	var in CancelRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorDTO{Code: "REASON_REQUIRED"}) // body อ่านไม่ได้ = ไม่มีเหตุผล (FR-CXL-08)
		return
	}
	res, err := h.svc.Cancel(r.Context(), p, id, in.Reason)
	switch {
	case errors.Is(err, ErrStudentOnly):
		writeJSON(w, http.StatusForbidden, ErrorDTO{Code: "STUDENT_ONLY"})
	case errors.Is(err, ErrReasonRequired):
		writeJSON(w, http.StatusBadRequest, ErrorDTO{Code: "REASON_REQUIRED"}) // FR-CXL-08
	case errors.Is(err, ErrReasonTooLong):
		writeJSON(w, http.StatusBadRequest, ErrorDTO{Code: "REASON_TOO_LONG"}) // FR-CXL-08
	case errors.Is(err, ErrAlreadyApproved):
		writeJSON(w, http.StatusConflict, ErrorDTO{Code: "ALREADY_APPROVED"}) // FR-CXL-02, FR-CXL-04
	case errors.Is(err, ErrNotCancelled):
		writeJSON(w, http.StatusConflict, ErrorDTO{Code: "NOT_CANCELLABLE"}) // สาเหตุอื่น แยกต่อใน T-17, T-19
	case err != nil:
		writeUpstreamError(w, err)
	default:
		writeJSON(w, http.StatusOK, res)
	}
}

// writeUpstreamError ตอบเมื่อฐานข้อมูลไม่ตอบสนอง ไม่ส่งรายละเอียดภายในกลับไป (plan.md ข้อ 4, ASM-06)
func writeUpstreamError(w http.ResponseWriter, err error) {
	log.Printf("cancelrequest: database error: %v", err)
	writeJSON(w, http.StatusServiceUnavailable, ErrorDTO{Code: "UPSTREAM_UNAVAILABLE", Retryable: true})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
