package handover

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// PathPrefix คือส่วนต้นของ path ทุก endpoint ของ 005 ที่หลังบ้านรับ
// plan.md ข้อ 4 เขียน /api/v1/handover/... ส่วน /api ถูกตัดโดย proxy ใน frontend/vite.config.js (แบบเดียวกับ 001)
const PathPrefix = "/v1/handover/"

// DB คือส่วนของ pgxpool.Pool ที่ 005 ใช้ (CON-TECH-01)
type DB interface {
	TxBeginner
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// Handler รับ HTTP request ของ 005 ทุก endpoint ต้องอยู่หลัง auth.Middleware (IF-IDP-01) และเป็นเจ้าหน้าที่ (CON-HND-01)
type Handler struct {
	db       DB
	verifier BorrowerVerifier // FR-HND-02 ยืนยันผู้ยืมที่เคาน์เตอร์ (ASM-HND-09)
}

// NewHandler สร้าง Handler จาก connection pool และระบบยืนยันตัวตน
// ผู้ยืมล็อกอินบัญชีมหาวิทยาลัยบนเครื่องเจ้าหน้าที่ (ASM-HND-09) prototype ใช้ provider จำลองตัวเดียวกับ middleware (Q-SH-03)
func NewHandler(db DB, provider auth.Provider) *Handler {
	return &Handler{db: db, verifier: NewMockBorrowerVerifier(provider)}
}

// requestIDFrom อ่าน {requestId} จาก path ok = false ถ้าไม่ใช่ตัวเลขบวก
func requestIDFrom(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("requestId"))
	return id, err == nil && id > 0
}

// Routes คืน router ของ 005 ห่อด้วย requireStaff ให้ main.go ห่อด้วย auth.Middleware แล้วผูกกับ PathPrefix
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+PathPrefix+"requests", h.listRequests)                                // FR-HND-01 (T-04)
	mux.HandleFunc("POST "+PathPrefix+"requests/{requestId}/verify-borrower", h.verifyBorrower) // FR-HND-02 (T-09)
	mux.HandleFunc("POST "+PathPrefix+"requests/{requestId}/handover", h.handover)              // FR-HND-02 ถึง 08 (T-10)
	return requireStaff(mux)
}

// requireStaff ปฏิเสธผู้ใช้ที่ไม่ใช่เจ้าหน้าที่ก่อนถึง handler ใด ๆ (CON-HND-01, AC-HND-09) ตอบ 403 FORBIDDEN_ROLE
func requireStaff(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := auth.FromContext(r.Context())
		if !ok || p.Role != auth.RoleStaff {
			writeError(w, http.StatusForbidden, codeForbiddenRole)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// listRequests: GET /v1/handover/requests รายการคำขอที่อนุมัติแล้วและยังไม่ส่งมอบทั้งหมด (FR-HND-01, AC-HND-14)
// ตอบเป็น JSON array ตาม plan.md ข้อ 4 ไม่มีคำขอ = []
func (h *Handler) listRequests(w http.ResponseWriter, r *http.Request) {
	items, err := ListReadyRequests(r.Context(), h.db, config.HandoverReadyStatuses())
	if err != nil {
		writeUpstreamError(w, err)
		return
	}
	if items == nil {
		items = []ReadyRequestDTO{}
	}
	writeJSON(w, http.StatusOK, items)
}

// VerifyBorrowerInput คือ body ของ verify-borrower (plan.md ข้อ 4)
type VerifyBorrowerInput struct {
	BorrowerCredential string `json:"borrowerCredential"` // ผู้ยืมล็อกอินบนเครื่องเจ้าหน้าที่ (ASM-HND-09) prototype = token จำลอง
}

// VerifyBorrowerResult คือผลการยืนยันผู้ยืม result = verified | mismatch | unavailable (plan.md ข้อ 4)
type VerifyBorrowerResult struct {
	Result string `json:"result"`
	UserID string `json:"userId,omitempty"`
}

// ผลการยืนยันผู้ยืม (plan.md ข้อ 4)
const (
	verifyVerified    = "verified"    // ยืนยันได้และเป็นผู้ส่งคำขอ
	verifyMismatch    = "mismatch"    // ยืนยันได้แต่ไม่ใช่ผู้ส่งคำขอ (FR-HND-02, ASM-HND-04)
	verifyUnavailable = "unavailable" // Auth API ไม่ตอบภายในเวลา ใช้ผลตรวจบัตรแทน (T-13)
)

// verifyBorrower: POST /v1/handover/requests/{requestId}/verify-borrower (FR-HND-02, IF-IDP-01, ASM-HND-09)
// ใช้แสดงผลบนหน้าจอก่อนบันทึก หลังบ้านยืนยันซ้ำเองอีกครั้งตอนบันทึก (plan.md v1.4)
// ล็อกอินไม่สำเร็จ (credential ใช้ไม่ได้) ตอบ 422 VALIDATION_FAILED ระบุช่อง borrowerCredential
func (h *Handler) verifyBorrower(w http.ResponseWriter, r *http.Request) {
	requestID, ok := requestIDFrom(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	var in VerifyBorrowerInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.BorrowerCredential == "" {
		writeError(w, http.StatusUnprocessableEntity, codeValidationFailed, "borrowerCredential")
		return
	}
	requesterID, found, err := GetRequester(r.Context(), h.db, requestID)
	if err != nil {
		writeUpstreamError(w, err)
		return
	}
	if !found {
		http.NotFound(w, r) // ไม่มีคำขอรหัสนี้ (plan ไม่ได้กำหนด error code ของกรณีนี้ ใช้แบบเดียวกับ 001)
		return
	}
	userID, err := h.verifier.VerifyBorrower(r.Context(), in.BorrowerCredential)
	if errors.Is(err, ErrBorrowerNotVerified) {
		writeError(w, http.StatusUnprocessableEntity, codeValidationFailed, "borrowerCredential")
		return
	}
	if err != nil {
		writeUpstreamError(w, err)
		return
	}
	result := verifyVerified
	if userID != requesterID {
		result = verifyMismatch
	}
	writeJSON(w, http.StatusOK, VerifyBorrowerResult{Result: result, UserID: userID})
}

// handover: POST /v1/handover/requests/{requestId}/handover (FR-HND-02 ถึง 08) ทำใน T-10
func (h *Handler) handover(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
