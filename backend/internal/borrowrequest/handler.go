package borrowrequest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"salmonnorway/backend/internal/auth"
)

// PathPrefix คือ path ของ endpoint ส่งคำขอยืม
// plan.md ข้อ 4 เขียน path ที่หน้าจอเรียกเป็น /api/v1/loan-requests ส่วน /api ถูกตัดออกโดย proxy
// ใน frontend/vite.config.js (template รายวิชา) หลังบ้านจึงรับที่ /v1/loan-requests
const PathPrefix = "/v1/loan-requests"

// QueuePath คือ endpoint ลำดับคำขอที่ค้างอยู่ของอุปกรณ์ชิ้นหนึ่ง (FR-BRQ-05 ข้อ ก)
// path อยู่ใต้ /v1/equipment/ ของ 001 แต่เจ้าของคือ 002 ตามสัญญา API ใน specs/changes/CR-005.md
// main.go ผูก pattern นี้ที่ mux ชั้นนอก Go 1.22 ส่ง request ให้ pattern ที่เจาะจงกว่า จึงไม่ชนกับของ 001
const QueuePath = "/v1/equipment/items/{assetCode}/requests"

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
	// เปิดดูไฟล์ใบบันทึกข้อความของคำขอใบนั้น (FR-BRQ-10 ข้อ ก, plan.md ข้อ 4)
	mux.HandleFunc("GET "+PathPrefix+"/{id}/attachment", h.attachment)
	// ลำดับคำขอที่ค้างอยู่ของอุปกรณ์ชิ้นนั้น ให้ 001 เรียกไปแสดง (FR-BRQ-05 ข้อ ก, CR-005)
	mux.HandleFunc("GET "+QueuePath, h.itemRequests)
	// method อื่นบน path เดียวกันได้ 405 จาก router
	return mux
}

// attachment: GET /v1/loan-requests/{id}/attachment ส่งไฟล์ใบบันทึกข้อความกลับไป (FR-BRQ-10 ข้อ ก)
// เปิดดูได้เฉพาะเจ้าหน้าที่กับเจ้าของคำขอ ตามที่ทีมตัดสินใน Q-BRQ-13 (ASM-BRQ-07, DOM-PRIV-01)
func (h *Handler) attachment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeError(w, http.StatusNotFound, "REQUEST_NOT_FOUND", "ไม่พบคำขอที่ระบุ")
		return
	}

	path, err := h.svc.OpenAttachment(r.Context(), id)
	switch {
	case err == nil:
		// ชนิดไฟล์ให้ net/http เดาจากนามสกุล ซึ่งถูกรักษาไว้ตอนเก็บไฟล์ (safeFileName)
		http.ServeFile(w, r, path)
	case errors.Is(err, ErrAttachmentNotFound):
		// ไม่แยกว่า "ไม่มีคำขอ" หรือ "ไม่มีไฟล์แนบ" เพื่อไม่บอกว่าเลขคำขอใดมีอยู่จริง (DOM-PRIV-01)
		writeError(w, http.StatusNotFound, "REQUEST_NOT_FOUND", "ไม่พบไฟล์แนบของคำขอที่ระบุ")
	case errors.Is(err, ErrNotRequestOwner):
		writeError(w, http.StatusForbidden, "FORBIDDEN_ROLE",
			"เปิดดูไฟล์แนบได้เฉพาะเจ้าหน้าที่และเจ้าของคำขอ")
	case errors.Is(err, ErrNoPrincipal):
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "ต้องเข้าสู่ระบบก่อนเปิดดูไฟล์แนบ")
	default:
		log.Printf("borrowrequest attachment: %v", err)
		writeError(w, http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", "ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง")
	}
}

// create: POST /v1/loan-requests สร้างคำขอยืมใบใหม่ (FR-BRQ-04, AC-BRQ-04)
// ผู้ขอมาจาก context ที่ auth.Middleware ใส่ไว้ ไม่ได้มาจาก body (ASM-BRQ-06)
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	// จำกัดขนาด request ทั้งก้อนตามเพดานไฟล์แนบ กันไฟล์ใหญ่เกินถูกอ่านเข้าเครื่องทั้งหมด (FR-BRQ-10 ข้อ ก)
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)

	var in CreateRequestDTO
	var up *Upload
	if isMultipart(r) {
		// ส่งมาพร้อมไฟล์แนบ (plan.md ข้อ 4) ใช้กับการยืมค้างคืนเป็นหลัก
		var err error
		if in, up, err = parseCreateForm(r); err != nil {
			writeError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "รูปแบบข้อมูลที่ส่งมาไม่ถูกต้อง")
			return
		}
	} else if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		// body ที่อ่านไม่ออก ไม่ใช่ข้อมูลที่ส่งมาไม่ครบ (การตรวจความครบถ้วนรายช่องเป็นงานของ T-06)
		writeError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "รูปแบบข้อมูลที่ส่งมาไม่ถูกต้อง")
		return
	}

	// ตรวจความครบถ้วนของข้อมูลก่อนอนุญาตให้ส่ง ถ้าไม่ครบ ต้องไม่สร้างคำขอ (FR-BRQ-03, AC-BRQ-03)
	// รวมกฎไฟล์แนบของการยืมค้างคืนไว้ในรายการเดียวกัน ผู้ขอจะได้เห็นทุกช่องที่ต้องแก้พร้อมกัน (FR-BRQ-10, AC-BRQ-10)
	newReq := in.toNewRequest()
	problems := append(missingFields(newReq), attachmentProblems(newReq.LoanType, up)...)
	if len(problems) > 0 {
		writeValidationError(w, problems)
		return
	}

	req, err := h.svc.CreateWithAttachment(r.Context(), newReq, up)
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
	case errors.Is(err, ErrDuplicateRequest):
		// ผู้ขอคนเดิมมีคำขอของอุปกรณ์ชิ้นนี้ที่ยังไม่จบอยู่แล้ว (FR-BRQ-05 ข, CON-TECH-03, AC-BRQ-11)
		writeError(w, http.StatusConflict, "DUPLICATE_REQUEST",
			"คุณมีคำขอของอุปกรณ์ชิ้นนี้ค้างอยู่แล้ว")
	case errors.Is(err, ErrHasOverdueLoan):
		// แจ้งเหตุผลให้ผู้ขอรู้ว่าต้องคืนของก่อน (FR-BRQ-06, CON-TECH-03, AC-BRQ-06)
		writeError(w, http.StatusConflict, "HAS_OVERDUE_LOAN",
			"คุณมีอุปกรณ์ที่คืนเกินกำหนดค้างอยู่ กรุณาคืนอุปกรณ์ก่อนส่งคำขอใหม่")
	case errors.Is(err, ErrNoPrincipal):
		// ปกติ auth.Middleware กันไว้ก่อนถึงตรงนี้แล้ว เหลือเป็นกันพลาดเมื่อผูก router ผิด (IF-IDP-01)
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "ต้องเข้าสู่ระบบก่อนส่งคำขอ")
	default:
		log.Printf("borrowrequest create: %v", err)
		writeError(w, http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", "ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง")
	}
}

// queueItemDTO คือคำขอ 1 ใบในลำดับคิวของอุปกรณ์ชิ้นนั้น รูปแบบล็อกไว้ในสัญญา API ของ CR-005
// ไม่มีรหัสผู้ใช้ของผู้ขอ มีแต่ชื่อที่ผู้ขอพิมพ์เอง เพราะผู้เรียกคือเจ้าหน้าที่ที่เห็นชื่อได้ (DOM-PRIV-01)
// (ประกาศไว้ที่นี่เพราะ dto.go ไม่อยู่ในช่องไฟล์ที่แตะของ T-18 ย้ายไปรวมกับ DTO ตัวอื่นได้ภายหลัง)
type queueItemDTO struct {
	RequestID     int       `json:"requestId"`
	QueuePosition int       `json:"queuePosition"`
	SubmittedAt   time.Time `json:"submittedAt"`
	RequesterName string    `json:"requesterName"`
}

// itemRequests: GET /v1/equipment/items/{assetCode}/requests ลำดับคำขอที่ค้างอยู่ของอุปกรณ์ชิ้นนั้น
// เฉพาะเจ้าหน้าที่ นักศึกษาได้ 403 เพราะต้องไม่เห็นคำขอของผู้อื่น (FR-BRQ-05 ข้อ ก, DOM-PRIV-01, CR-005)
// 002 เป็นเจ้าของ endpoint นี้ ส่วนการแสดงผลเป็นของ 001 ตาม CR-003
func (h *Handler) itemRequests(w http.ResponseWriter, r *http.Request) {
	p, ok := auth.FromContext(r.Context())
	if !ok || !p.Complete() {
		// ปกติ auth.Middleware กันไว้ก่อนถึงตรงนี้แล้ว (IF-IDP-01)
		writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED", "ต้องเข้าสู่ระบบก่อน")
		return
	}
	if p.Role != auth.RoleStaff {
		writeError(w, http.StatusForbidden, "FORBIDDEN_ROLE", "ดูลำดับคำขอได้เฉพาะเจ้าหน้าที่")
		return
	}

	// เรียก repository ผ่าน service โดยตรง เพราะ service.go ไม่อยู่ในช่องไฟล์ที่แตะของ T-18
	// ย้ายเป็น method ของ service ได้ภายหลัง งานนี้ไม่มีกฎของฟีเจอร์นอกจากการเรียงลำดับ
	pending, err := h.svc.repo.PendingByItem(r.Context(), r.PathValue("assetCode"))
	if err != nil {
		log.Printf("borrowrequest itemRequests: %v", err)
		writeError(w, http.StatusServiceUnavailable, "UPSTREAM_UNAVAILABLE", "ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง")
		return
	}

	// ลำดับที่ 1 คือคำขอที่ส่งก่อนสุด นับจากการเรียงของฐานข้อมูล (DOM-REQ-01)
	out := make([]queueItemDTO, 0, len(pending))
	for i, q := range pending {
		out = append(out, queueItemDTO{
			RequestID:     q.RequestID,
			QueuePosition: i + 1,
			SubmittedAt:   q.SubmittedAt.In(appLoc), // เวลาตามเขตเวลา Asia/Bangkok (glossary)
			RequesterName: q.RequesterName,
		})
	}
	writeJSON(w, http.StatusOK, out)
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
