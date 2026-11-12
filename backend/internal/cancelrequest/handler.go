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
		writeError(w, http.StatusUnauthorized, errUnauthenticated()) // ไม่ควรเกิดถ้าอยู่หลัง middleware
		return
	}
	items, err := h.svc.ListMyPending(r.Context(), p)
	switch {
	case errors.Is(err, ErrStudentOnly):
		writeError(w, http.StatusForbidden, errForbiddenRole()) // ASM-13, CR-006
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
		writeError(w, http.StatusUnauthorized, errUnauthenticated()) // ไม่ควรเกิดถ้าอยู่หลัง middleware
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeError(w, http.StatusNotFound, errRequestNotFound()) // รหัสที่ไม่ใช่ตัวเลข ถือว่าไม่มีคำขอนี้ (FR-CXL-03)
		return
	}
	var in CancelRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusUnprocessableEntity, errValidationFailed()) // body อ่านไม่ได้ = เหตุผลผิด (FR-CXL-08, CR-006 ปิด F-06)
		return
	}
	res, err := h.svc.Cancel(r.Context(), p, id, in.Reason)
	switch {
	case errors.Is(err, ErrStudentOnly):
		writeError(w, http.StatusForbidden, errForbiddenRole()) // ASM-13, CR-006
	case errors.Is(err, ErrReasonRequired), errors.Is(err, ErrReasonTooLong):
		writeError(w, http.StatusUnprocessableEntity, errValidationFailed()) // FR-CXL-08, CR-006
	case errors.Is(err, ErrRequestNotFound):
		writeError(w, http.StatusNotFound, errRequestNotFound()) // FR-CXL-03
	case errors.Is(err, ErrAlreadyApproved):
		writeError(w, http.StatusConflict, ErrorDTO{Code: codeAlreadyApproved, Message: msgAlreadyApproved, Status: statusOf(err)}) // FR-CXL-02, FR-CXL-04
	case errors.Is(err, ErrNotCancelled):
		writeError(w, http.StatusConflict, ErrorDTO{Code: codeNotCancellable, Message: msgNotCancellable, Status: statusOf(err)}) // FR-CXL-02 สถานะอื่นที่ยกเลิกไม่ได้
	case err != nil:
		writeUpstreamError(w, err)
	default:
		writeJSON(w, http.StatusOK, res)
	}
}

// statusOf คืนสถานะของคำขอที่แนบมากับ StatusError (CR-006) ไม่มีคืน ""
func statusOf(err error) string {
	var se *StatusError
	if errors.As(err, &se) {
		return se.Status
	}
	return ""
}

// writeUpstreamError ตอบเมื่อฐานข้อมูลไม่ตอบสนองหรือเกิน CancelDBTimeout ไม่ส่งรายละเอียดภายในกลับไป (plan.md ข้อ 4, ASM-06, AC-CXL-13)
// 503 ลองใหม่ได้เสมอ ไม่แนบ retryable (CR-006)
func writeUpstreamError(w http.ResponseWriter, err error) {
	log.Printf("cancelrequest: database error: %v", err)
	writeError(w, http.StatusServiceUnavailable, errUpstreamUnavailable())
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
