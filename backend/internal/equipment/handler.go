package equipment

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"salmonnorway/backend/internal/platform"
)

// PathPrefix คือส่วนต้นของ path ทุก endpoint ของฟีเจอร์นี้ที่หลังบ้านรับ
// plan.md ข้อ 5 เขียน path ที่หน้าจอเรียกเป็น /api/v1/equipment/... ส่วน /api ถูกตัดออกโดย proxy
// ใน frontend/vite.config.js (template รายวิชา) หลังบ้านจึงรับที่ /v1/equipment/...
const PathPrefix = "/v1/equipment/"

// Handler รับ HTTP request ของฟีเจอร์ 001 ทุก endpoint ต้องอยู่หลัง auth.Middleware (IF-IDP-01, NFR-SEC-01)
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
	mux.HandleFunc("GET "+PathPrefix+"models", h.listModels)
	mux.HandleFunc("GET "+PathPrefix+"models/{modelId}/items", h.listItems)
	mux.HandleFunc("GET "+PathPrefix+"items/{assetCode}", h.getItem)
	mux.HandleFunc("GET "+PathPrefix+"summary", h.summary)
	// ไม่มี endpoint ที่เขียนข้อมูล ฟีเจอร์นี้อ่านอย่างเดียว (CON-DATA-02) method อื่นได้ 405 จาก router
	return mux
}

// listModels: GET /v1/equipment/models?category=&status=&department= รายการรุ่นพร้อมจำนวนว่าง/ทั้งหมด (FR-EQS-01, FR-EQS-02, AC-EQS-01, AC-EQS-02)
// นักศึกษาเห็นเฉพาะภาควิชาของตน และขอ department อื่นได้ 403 (FR-EQS-10, AC-EQS-13)
func (h *Handler) listModels(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListModels(r.Context(), filterFrom(r), requestedDepartment(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ListDTO[ModelRowDTO]{Items: nonNil(rows), Total: len(rows)})
}

// listItems: GET /v1/equipment/models/{modelId}/items อุปกรณ์รายชิ้นของรุ่นที่เลือก (FR-EQS-01, AC-EQS-01)
// ใช้ตัวกรองเดียวกับรายการรุ่น (FR-EQS-02)
func (h *Handler) listItems(w http.ResponseWriter, r *http.Request) {
	modelID, err := strconv.Atoi(r.PathValue("modelId"))
	if err != nil || modelID <= 0 {
		http.NotFound(w, r) // รหัสรุ่นที่ไม่ใช่ตัวเลข ถือว่าไม่มี path นี้
		return
	}
	items, err := h.svc.ListItems(r.Context(), modelID, filterFrom(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ListDTO[ItemRow]{Items: nonNil(items), Total: len(items)})
}

// getItem: GET /v1/equipment/items/{assetCode} รายละเอียดอุปกรณ์รายชิ้น พร้อมกำหนดคืนเมื่อถูกยืม (FR-EQS-04, AC-EQS-04)
func (h *Handler) getItem(w http.ResponseWriter, r *http.Request) {
	item, found, err := h.svc.GetItem(r.Context(), r.PathValue("assetCode"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if !found {
		http.NotFound(w, r) // ไม่มีรหัสครุภัณฑ์นี้ (plan ไม่ได้กำหนด error code ของกรณีนี้)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// summary: GET /v1/equipment/summary จำนวนอุปกรณ์แยกตามสถานะ เฉพาะเจ้าหน้าที่ (FR-EQS-07)
// นักศึกษาได้ 403 FORBIDDEN_ROLE ตาม CON-TECH-03 และไม่ได้ตัวเลขใด ๆ
func (h *Handler) summary(w http.ResponseWriter, r *http.Request) {
	out, err := h.svc.Summary(r.Context())
	if errors.Is(err, ErrStaffOnly) {
		platform.WriteError(w, http.StatusForbidden, platform.CodeForbiddenRole)
		return
	}
	if err != nil {
		writeUpstreamError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

// requestedDepartment อ่าน ?department= (plan.md ข้อ 5) ไม่ได้ส่ง = 0 / ค่าที่ไม่ใช่รหัสภาควิชา = -1 (ไม่ตรงกับภาควิชาใด)
// ตอนนี้ใช้ตรวจสิทธิ์ของนักศึกษาอย่างเดียว เจ้าหน้าที่ไม่ได้กรองด้วยค่านี้ (ทีมตัดสิน 2569-11-12 รอ Q-P3)
func requestedDepartment(r *http.Request) int {
	v := r.URL.Query().Get("department")
	if v == "" {
		return 0
	}
	id, err := strconv.Atoi(v)
	if err != nil || id <= 0 {
		return -1
	}
	return id
}

// writeServiceError แปลง error จาก service เป็นคำตอบตาม CON-TECH-03
//   - ขอข้อมูลภาควิชาอื่น -> 403 FORBIDDEN_DEPARTMENT ไม่มีข้อมูลอุปกรณ์ใน body (FR-EQS-10, AC-EQS-13)
//   - อื่น ๆ -> 503 UPSTREAM_UNAVAILABLE (FR-EQS-08)
func writeServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrForbiddenDepartment) {
		platform.WriteError(w, http.StatusForbidden, platform.CodeForbiddenDepartment)
		return
	}
	writeUpstreamError(w, err)
}

// filterFrom อ่านตัวกรองจาก query string (FR-EQS-02) เลือกได้หลายค่าในแต่ละตัวกรอง (Q-15)
// ส่งค่าซ้ำชื่อเดิมได้ เช่น ?status=available&status=borrowed&category=<ชื่อประเภท>
func filterFrom(r *http.Request) Filter {
	q := r.URL.Query()
	return Filter{
		Categories: nonEmpty(q["category"]),
		Statuses:   nonEmpty(q["status"]),
	}
}

// nonEmpty ตัดค่าว่างออก (?status= ไม่นับเป็นตัวกรอง)
func nonEmpty(values []string) []string {
	var out []string
	for _, v := range values {
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

// writeUpstreamError ตอบ 503 UPSTREAM_UNAVAILABLE เมื่ออ่านฐานข้อมูลไม่สำเร็จ (FR-EQS-08) body {code, message} ตาม CON-TECH-03
// ไม่ส่งข้อมูลเก่าหรือรายละเอียดภายในกลับไป (test และหน้าจอของกรณีนี้อยู่ใน T-20)
func writeUpstreamError(w http.ResponseWriter, err error) {
	log.Printf("equipment: database error: %v", err)
	platform.WriteError(w, http.StatusServiceUnavailable, platform.CodeUpstreamUnavailable)
}

// nonNil ให้ JSON เป็น [] แทน null เมื่อไม่มีข้อมูล (plan.md ข้อ 5: ไม่พบ = items ว่าง)
func nonNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
