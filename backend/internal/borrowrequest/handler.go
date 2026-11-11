package borrowrequest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// PathPrefix คือ path ของ endpoint ส่งคำขอยืม
// plan.md ข้อ 4 เขียน path ที่หน้าจอเรียกเป็น /api/v1/loan-requests ส่วน /api ถูกตัดออกโดย proxy
// ใน frontend/vite.config.js (template รายวิชา) หลังบ้านจึงรับที่ /v1/loan-requests
const PathPrefix = "/v1/loan-requests"

// Handler รับ HTTP request ของฟีเจอร์ 002
// ทุก endpoint ต้องอยู่หลัง auth.Middleware (IF-IDP-01, NFR-SEC-01) และให้บริการผ่าน HTTPS เท่านั้น (CON-TECH-02)
type Handler struct {
	svc *Service
}

// NewHandler สร้าง Handler
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes คืน router ของฟีเจอร์นี้ ให้ main.go ห่อด้วย auth.Middleware แล้วผูกกับ PathPrefix
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+PathPrefix, h.create)
	// method อื่นบน path เดียวกันได้ 405 จาก router ฟีเจอร์นี้มีแต่การสร้างคำขอ
	return mux
}

// create: POST /v1/loan-requests สร้างคำขอยืมใบใหม่ (FR-BRQ-04, AC-BRQ-04)
// ผู้ขอมาจาก context ที่ auth.Middleware ใส่ไว้ ไม่ได้มาจาก body (ASM-BRQ-06)
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in CreateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		// body ที่อ่านไม่ออก ไม่ใช่ข้อมูลที่ส่งมาไม่ครบ (การตรวจความครบถ้วนรายช่องเป็นงานของ T-06)
		writeError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "รูปแบบข้อมูลที่ส่งมาไม่ถูกต้อง")
		return
	}

	// ตรวจความครบถ้วนของข้อมูลก่อนอนุญาตให้ส่ง ถ้าไม่ครบ ต้องไม่สร้างคำขอ (FR-BRQ-03, AC-BRQ-03)
	newReq := in.toNewRequest()
	if missing := missingFields(newReq); len(missing) > 0 {
		writeValidationError(w, missing)
		return
	}

	req, err := h.svc.Create(r.Context(), newReq)
	switch {
	case err == nil:
		writeJSON(w, http.StatusCreated, CreatedDTO{
			RequestID:   req.ID,
			Status:      req.Status,
			SubmittedAt: req.SubmittedAt,
		})
	case errors.Is(err, ErrForbiddenDepartment):
		// ไม่ส่งข้อมูลอุปกรณ์ใด ๆ กลับไป แม้แต่สถานะ (FR-BRQ-07, DOM-ACCESS-01, AC-BRQ-07)
		writeError(w, http.StatusForbidden, "FORBIDDEN_DEPARTMENT",
			"ส่งคำขอยืมได้เฉพาะอุปกรณ์ของภาควิชาตนเอง")
	case errors.Is(err, ErrItemNotAvailable):
		writeItemNotAvailable(w, err)
	case errors.Is(err, ErrNoPrincipal):
		// ปกติ auth.Middleware กันไว้ก่อนถึงตรงนี้แล้ว เหลือเป็นกันพลาดเมื่อผูก router ผิด (IF-IDP-01)
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "ต้องเข้าสู่ระบบก่อนส่งคำขอ")
	default:
		log.Printf("borrowrequest create: %v", err)
		writeError(w, http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", "ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง")
	}
}

// writeJSON ส่งคำตอบเป็น JSON
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// writeError ส่ง error ตามรูปแบบกลาง {code, message} ด้วยรหัสจากรายการของ CON-TECH-03 (CR-004)
func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorDTO{Code: code, Message: message})
}

// validationErrorDTO คือรูปแบบกลาง {code, message} ที่แนบรายชื่อช่องที่ยังขาดมาด้วย
// CON-TECH-03 กำหนดให้ VALIDATION_FAILED "แนบรายชื่อช่องที่ผิด" (CR-004)
type validationErrorDTO struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Fields  []string `json:"fields"`
}

// itemNotAvailableDTO คือรูปแบบกลาง {code, message} ที่แนบสถานะล่าสุดของอุปกรณ์มาด้วย (CON-TECH-03)
// ค่า status เป็นค่าของ CON-DATA-01 (available / borrowed / unavailable)
// หน้าจอแปลงเป็นคำไทยเอง แบบเดียวกับ StatusBadge ของ 001
type itemNotAvailableDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
}

// writeItemNotAvailable แจ้งว่าส่งคำขอไม่ได้ พร้อมสถานะล่าสุดของอุปกรณ์ (FR-BRQ-01, AC-BRQ-01)
// ถ้าไม่มีรหัสครุภัณฑ์นั้นอยู่จริง จะไม่มีสถานะให้แจ้ง ตอบเหมือนกรณีไม่ว่างเพื่อไม่เปิดเผยว่ามีรหัสนั้นไหม
func writeItemNotAvailable(w http.ResponseWriter, err error) {
	body := itemNotAvailableDTO{
		Code:    "ITEM_NOT_AVAILABLE",
		Message: "อุปกรณ์ชิ้นนี้ไม่อยู่ในสถานะว่างแล้ว",
	}
	var notAvail *ItemNotAvailableError
	if errors.As(err, &notAvail) {
		body.Status = notAvail.Status
	}
	writeJSON(w, http.StatusConflict, body)
}

// writeValidationError บอกผู้ขอว่าช่องไหนยังไม่ได้กรอก และไม่สร้างคำขอ (FR-BRQ-03, AC-BRQ-03)
func writeValidationError(w http.ResponseWriter, fields []string) {
	writeJSON(w, http.StatusUnprocessableEntity, validationErrorDTO{
		Code:    "VALIDATION_FAILED",
		Message: "กรุณากรอกข้อมูลให้ครบทุกช่องก่อนส่งคำขอ",
		Fields:  fields,
	})
}
