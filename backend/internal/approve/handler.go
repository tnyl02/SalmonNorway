package approve

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

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
	mux.HandleFunc("POST "+PathPrefix+"requests/{id}/approve", h.approve)
	mux.HandleFunc("POST "+PathPrefix+"requests/{id}/reject", h.reject)
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

// approve: POST /v1/approvals/requests/{id}/approve อนุมัติคำขอยืมปกติ (FR-APV-01 ค, AC-APV-01)
// body { loadedAt } มาจาก GET รายการ ใช้แยก "พิจารณาแล้ว" กับ "กดพร้อมกัน" ใน T-23 (plan.md ข้อ 4)
func (h *Handler) approve(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "REQUEST_NOT_FOUND") // Q-APV-25
		return
	}
	var in DecideInputDTO
	// body ที่อ่านไม่ออกหรือไม่มี loadedAt ไม่ปฏิเสธ แค่ถือว่าไม่รู้เวลาโหลด คำขอที่พิจารณาไปแล้วจึงได้ REQUEST_ALREADY_DECIDED (T-23)
	_ = json.NewDecoder(r.Body).Decode(&in)

	res, err := h.svc.Approve(r.Context(), id, in.LoadedAt)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResultDTO(res))
}

// reject: POST /v1/approvals/requests/{id}/reject ปฏิเสธคำขอยืมปกติพร้อมเหตุผล (FR-APV-01 ง, AC-APV-02)
// body { loadedAt, reason } เหตุผลว่างได้ 422 VALIDATION_FAILED (T-12, AC-APV-06)
func (h *Handler) reject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "REQUEST_NOT_FOUND") // Q-APV-25
		return
	}
	var in DecideInputDTO
	_ = json.NewDecoder(r.Body).Decode(&in) // body ที่อ่านไม่ออกได้เหตุผลว่าง จึงได้ VALIDATION_FAILED ช่อง reason

	res, err := h.svc.Reject(r.Context(), id, in.Reason, in.LoadedAt)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toResultDTO(res))
}

// writeServiceError แปลง error ของ service เป็น {code, message} ตาม CON-TECH-03 ข้อความจาก config (FR-APV-01 ซ)
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrReasonRequired):
		// CON-TECH-03: VALIDATION_FAILED แนบรายชื่อช่องที่ผิด ช่องเดียวที่ 004 ตรวจคือ reason (AC-APV-06)
		writeJSON(w, http.StatusUnprocessableEntity, ValidationErrorDTO{
			Code: "VALIDATION_FAILED", Message: config.ApprovalMessages["VALIDATION_FAILED"], Fields: []string{"reason"},
		})
	case errors.Is(err, ErrForbiddenDepartment):
		// ไม่ส่งข้อมูลอุปกรณ์กลับ ตอบแค่ {code, message} (CON-TECH-03, AC-APV-12)
		writeError(w, http.StatusForbidden, "FORBIDDEN_DEPARTMENT")
	case errors.Is(err, ErrApprovalCutoff):
		writeError(w, http.StatusConflict, "APPROVAL_CUTOFF_PASSED") // FR-APV-01 (ค)(ซ), AC-APV-16
	case errors.Is(err, ErrItemNotAvailable):
		writeError(w, http.StatusConflict, "ITEM_NOT_AVAILABLE") // FR-APV-01 (ค)(ซ), AC-APV-08
	case errors.Is(err, ErrNotFirst):
		writeError(w, http.StatusConflict, "NOT_FIRST_IN_QUEUE") // FR-APV-01 (ข)(ซ), AC-APV-07
	case errors.Is(err, ErrRequestNotFound):
		writeError(w, http.StatusNotFound, "REQUEST_NOT_FOUND") // Q-APV-25
	case errors.Is(err, ErrRequestCancelled):
		writeError(w, http.StatusConflict, "REQUEST_CANCELLED") // FR-APV-01 (ซ), AC-APV-09
	case errors.Is(err, ErrAlreadyDecided):
		writeError(w, http.StatusConflict, "REQUEST_ALREADY_DECIDED") // FR-APV-01 (ซ), AC-APV-09
	case errors.Is(err, ErrDecidedConcurrently):
		writeError(w, http.StatusConflict, "REQUEST_DECIDED_CONCURRENTLY") // FR-APV-01 (ช)(ซ), AC-APV-10
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
