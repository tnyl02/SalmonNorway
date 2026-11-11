package returns

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"salmonnorway/backend/internal/auth"
)

// Handler รวม endpoint ของการรับคืนตาม plan.md ข้อ 4 (FR-RET-01 ถึง FR-RET-07)
type Handler struct {
	svc *Service
}

// NewHandler สร้าง Handler ของฟีเจอร์รับคืน
func NewHandler(repo *Repository) *Handler {
	return &Handler{svc: NewService(repo)}
}

// Routes คืนฟังก์ชันที่ผูก endpoint ทั้งหมดของ 006 เข้ากับ router
// ทุก endpoint ต้องผ่าน auth.Middleware ก่อนเสมอ (IF-IDP-01, NFR-SEC-01)
func (h *Handler) Routes(provider auth.Provider) func(*http.ServeMux) {
	return func(mux *http.ServeMux) {
		for pattern, handle := range map[string]http.HandlerFunc{
			"GET /api/v1/returns/open-loans":     h.openLoans,    // FR-RET-01(ข), FR-RET-06
			"GET /api/v1/returns/lookup":         h.lookup,       // FR-RET-01(ก)
			"POST /api/v1/returns":               h.recordReturn, // FR-RET-02, 03, 04, 06, 07
			"GET /api/v1/returns/loans/{loanId}": h.loanDetail,   // FR-RET-06
		} {
			mux.Handle(pattern, auth.Middleware(provider, handle))
		}
	}
}

// returnResponse คือ output ของ POST /api/v1/returns ตาม plan.md ข้อ 4
type returnResponse struct {
	LoanID       int       `json:"loanId"`
	AssetCode    string    `json:"assetCode"`
	ReturnedAt   time.Time `json:"returnedAt"`
	Condition    string    `json:"condition"`
	ReturnedLate bool      `json:"returnedLate"`
	ItemStatus   string    `json:"itemStatus"`
}

// recordReturn รับ assetCode, condition (form แบบ urlencoded หรือ multipart) แล้วบันทึกการรับคืน (FR-RET-02, FR-RET-03, FR-RET-04)
// ไม่อ่านฟิลด์เวลาใด ๆ จาก client เวลาคืนมาจากฐานข้อมูลเท่านั้น (FR-RET-02, ASM-04)
func (h *Handler) recordReturn(w http.ResponseWriter, r *http.Request) {
	rec, err := h.svc.RecordReturn(r.Context(), r.FormValue("assetCode"), r.FormValue("condition"))
	switch {
	case errors.Is(err, ErrUnknownCondition):
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"code": "CONDITION_REQUIRED"})
	case errors.Is(err, ErrNoOpenLoan):
		writeJSON(w, http.StatusConflict, map[string]string{"code": "NO_OPEN_LOAN"})
	case err != nil:
		log.Printf("returns: record return: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"code": "UPSTREAM_UNAVAILABLE", "retryable": true})
	default:
		writeJSON(w, http.StatusOK, returnResponse(rec))
	}
}

// lookup หาการยืมที่ยังไม่คืนจากรหัสครุภัณฑ์ ก่อนบันทึกการรับคืน (FR-RET-01(ก))
// ไม่พบ -> 409 NO_OPEN_LOAN (ข้อความบนหน้าจอและ test ของ AC-RET-03 อยู่ใน T-18, T-19)
func (h *Handler) lookup(w http.ResponseWriter, r *http.Request) {
	ol, err := h.svc.repo.FindOpenLoan(r.Context(), r.URL.Query().Get("assetCode"))
	switch {
	case errors.Is(err, ErrNoOpenLoan):
		writeJSON(w, http.StatusConflict, map[string]string{"code": "NO_OPEN_LOAN"})
	case err != nil:
		log.Printf("returns: lookup: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"code": "UPSTREAM_UNAVAILABLE", "retryable": true})
	default:
		writeJSON(w, http.StatusOK, ol)
	}
}

// loanDetail แสดงรายละเอียดการยืม: เวลาคืน สภาพ และคืนเกินกำหนด (FR-RET-06)
// รหัสไม่ใช่ตัวเลข หรือไม่มีการยืมนี้ -> 404 LOAN_NOT_FOUND
func (h *Handler) loanDetail(w http.ResponseWriter, r *http.Request) {
	loanID, err := strconv.Atoi(r.PathValue("loanId"))
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"code": "LOAN_NOT_FOUND"})
		return
	}
	d, err := h.svc.repo.GetLoan(r.Context(), loanID)
	switch {
	case errors.Is(err, ErrLoanNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"code": "LOAN_NOT_FOUND"})
	case err != nil:
		log.Printf("returns: loan detail: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"code": "UPSTREAM_UNAVAILABLE", "retryable": true})
	default:
		writeJSON(w, http.StatusOK, d)
	}
}

// openLoans แสดงรายการการยืมที่ยังไม่คืนของภาควิชาเดียวกับเจ้าหน้าที่ (FR-RET-01(ข), ASM-03)
// ภาควิชามาจากผู้ใช้ที่ยืนยันตัวตนแล้ว (IF-IDP-01) ไม่รับจาก query string
func (h *Handler) openLoans(w http.ResponseWriter, r *http.Request) {
	p, ok := auth.FromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"code": "UNAUTHENTICATED"})
		return
	}
	loans, err := h.svc.repo.ListOpenLoans(r.Context(), p.DepartmentID)
	if err != nil {
		log.Printf("returns: open loans: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"code": "UPSTREAM_UNAVAILABLE", "retryable": true})
		return
	}
	writeJSON(w, http.StatusOK, loans)
}

// writeJSON ตอบ JSON และตั้ง Cache-Control: no-store ทุก response ของ 006 (plan.md ข้อ 4, NFR-DATA-01)
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
