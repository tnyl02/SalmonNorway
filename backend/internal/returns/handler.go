package returns

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
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
// ทุก endpoint ต้องผ่าน auth.Middleware (IF-IDP-01, NFR-SEC-01) แล้วตามด้วย staffOnly (FR-RET-01, AC-RET-08)
func (h *Handler) Routes(provider auth.Provider) func(*http.ServeMux) {
	return func(mux *http.ServeMux) {
		for _, rt := range []struct {
			pattern, action string
			handle          http.HandlerFunc
		}{
			{"GET /v1/returns/open-loans", actionReturnList, h.openLoans},        // FR-RET-01(ข), FR-RET-06
			{"GET /v1/returns/lookup", actionReturnLookup, h.lookup},             // FR-RET-01(ก)
			{"POST /v1/returns", actionReturn, h.recordReturn},                   // FR-RET-02, 03, 04, 06, 07
			{"GET /v1/returns/loans/{loanId}", actionReturnDetail, h.loanDetail}, // FR-RET-06
		} {
			mux.Handle(rt.pattern, auth.Middleware(provider, h.staffOnly(rt.action, rt.handle)))
		}
	}
}

// returnResponse คือ output ของ POST /v1/returns ตาม plan.md ข้อ 4
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
// สำเร็จ: audit "return" อยู่ใน transaction แล้ว / ไม่สำเร็จ: เขียน audit นอก transaction (FR-RET-07)
func (h *Handler) recordReturn(w http.ResponseWriter, r *http.Request) {
	p, ok := principal(w, r)
	if !ok {
		return
	}
	assetCode := r.FormValue("assetCode")
	e := newAudit(p, actionReturn)
	e.AssetCode = assetCode
	rec, err := h.svc.RecordReturn(r.Context(), assetCode, r.FormValue("condition"), p.DepartmentID, e)
	if err == nil {
		writeJSON(w, http.StatusOK, returnResponse(rec))
		return
	}
	status, body := errorResponse(err, "record return")
	h.respond(w, r, e, status, body)
}

// lookup หาการยืมที่ยังไม่คืนจากรหัสครุภัณฑ์ ก่อนบันทึกการรับคืน (FR-RET-01(ก)) audit "return_lookup" (DOM-AUDIT-01)
// ไม่พบ -> 409 NO_OPEN_LOAN (ข้อความบนหน้าจอและ test ของ AC-RET-03 อยู่ใน T-18, T-19)
func (h *Handler) lookup(w http.ResponseWriter, r *http.Request) {
	p, ok := principal(w, r)
	if !ok {
		return
	}
	assetCode := r.URL.Query().Get("assetCode")
	e := newAudit(p, actionReturnLookup)
	e.AssetCode, e.Query = assetCode, map[string]string{"assetCode": assetCode}
	ol, err := h.svc.repo.FindOpenLoan(r.Context(), assetCode, p.DepartmentID)
	if err != nil {
		status, body := errorResponse(err, "lookup")
		h.respond(w, r, e, status, body)
		return
	}
	h.respond(w, r, e, http.StatusOK, ol)
}

// loanDetail แสดงรายละเอียดการยืม: เวลาคืน สภาพ และคืนเกินกำหนด (FR-RET-06) audit "return_detail" (DOM-AUDIT-01)
// รหัสไม่ใช่ตัวเลข หรือไม่มีการยืมนี้ -> 404 LOAN_NOT_FOUND
func (h *Handler) loanDetail(w http.ResponseWriter, r *http.Request) {
	p, ok := principal(w, r)
	if !ok {
		return
	}
	e := newAudit(p, actionReturnDetail)
	e.Query = map[string]string{"loanId": r.PathValue("loanId")}
	loanID, err := strconv.Atoi(r.PathValue("loanId"))
	if err != nil {
		h.respond(w, r, e, http.StatusNotFound, errorBody("LOAN_NOT_FOUND"))
		return
	}
	d, err := h.svc.repo.GetLoan(r.Context(), loanID, p.DepartmentID)
	if err != nil {
		status, body := errorResponse(err, "loan detail")
		h.respond(w, r, e, status, body)
		return
	}
	e.AssetCode = d.AssetCode
	h.respond(w, r, e, http.StatusOK, d)
}

// openLoans แสดงรายการการยืมที่ยังไม่คืนของภาควิชาเดียวกับเจ้าหน้าที่ (FR-RET-01(ข), ASM-03) audit "return_list" (DOM-AUDIT-01)
// ภาควิชามาจากผู้ใช้ที่ยืนยันตัวตนแล้ว (IF-IDP-01) ไม่รับจาก query string
func (h *Handler) openLoans(w http.ResponseWriter, r *http.Request) {
	p, ok := principal(w, r)
	if !ok {
		return
	}
	e := newAudit(p, actionReturnList)
	loans, err := h.svc.repo.ListOpenLoans(r.Context(), p.DepartmentID)
	if err != nil {
		status, body := errorResponse(err, "open loans")
		h.respond(w, r, e, status, body)
		return
	}
	h.respond(w, r, e, http.StatusOK, loans)
}

// errorResponse แปลง error ของฟีเจอร์เป็น HTTP status และ body ตามตาราง error ใน plan.md ข้อ 4
func errorResponse(err error, op string) (int, any) {
	switch {
	case errors.Is(err, ErrForbiddenDepartment):
		return http.StatusForbidden, errorBody("FORBIDDEN_DEPARTMENT") // AC-RET-11
	case errors.Is(err, ErrUnknownCondition):
		return http.StatusUnprocessableEntity, errorBody("CONDITION_REQUIRED") // AC-RET-04
	case errors.Is(err, ErrNoOpenLoan):
		return http.StatusConflict, errorBody("NO_OPEN_LOAN") // AC-RET-03
	case errors.Is(err, ErrLoanNotFound):
		return http.StatusNotFound, errorBody("LOAN_NOT_FOUND")
	default:
		log.Printf("returns: %s: %v", op, err)
		return http.StatusServiceUnavailable, unavailableBody()
	}
}

// errorBody คือ body ของ error ตาม CON-TECH-03: {"code", "message"} ข้อความจาก config.ReturnMessages
func errorBody(code string) map[string]any {
	return map[string]any{"code": code, "message": config.ReturnMessages[code]}
}

// unavailableBody คือ 503 UPSTREAM_UNAVAILABLE ตาม CON-TECH-03 มีแค่ code และ message
// ไม่แนบ retryable เพราะตารางกลางไม่ได้กำหนด (rtm F-05 ทีมตัดสิน 2569-11-12 แบบเดียวกับ 003 CR-006)
func unavailableBody() map[string]any {
	return errorBody("UPSTREAM_UNAVAILABLE")
}

// respond เขียน audit ของ request นี้ก่อน แล้วจึงตอบ (DOM-AUDIT-01 "ทุกครั้งที่เข้าถึง")
// เขียน audit ไม่สำเร็จ -> 503 และไม่ส่งข้อมูลกลับ (fail closed แบบเดียวกับ plan 6.4 ของ 001)
func (h *Handler) respond(w http.ResponseWriter, r *http.Request, e AuditEntry, status int, body any) {
	e.Outcome = outcomeFor(status)
	if err := h.svc.repo.WriteAudit(r.Context(), e); err != nil {
		log.Printf("returns: write audit: %v", err)
		writeJSON(w, http.StatusServiceUnavailable, unavailableBody())
		return
	}
	writeJSON(w, status, body)
}

// staffOnly ให้ผ่านเฉพาะเจ้าหน้าที่ บทบาทอื่น (นักศึกษา) ได้ 403 FORBIDDEN_ROLE ก่อนถึง handler
// จึงไม่มีการอ่านหรือเขียนข้อมูลการยืมเลย (FR-RET-01 "เจ้าหน้าที่", RQ-11, AC-RET-08)
// การถูกปฏิเสธบันทึกลง audit ด้วย action ของ endpoint นั้น outcome forbidden (DOM-AUDIT-01)
func (h *Handler) staffOnly(action string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, ok := principal(w, r)
		if !ok {
			return
		}
		if p.Role != auth.RoleStaff {
			h.respond(w, r, newAudit(p, action), http.StatusForbidden, errorBody("FORBIDDEN_ROLE"))
			return
		}
		next(w, r)
	}
}

// principal อ่านผู้ใช้ที่ middleware ยืนยันตัวตนแล้ว (IF-IDP-01) ไม่มี -> ตอบ 401 แล้วคืน false
func principal(w http.ResponseWriter, r *http.Request) (auth.Principal, bool) {
	p, ok := auth.FromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errorBody("UNAUTHENTICATED"))
	}
	return p, ok
}

// writeJSON ตอบ JSON และตั้ง Cache-Control: no-store ทุก response ของ 006 (plan.md ข้อ 4, NFR-DATA-01)
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
