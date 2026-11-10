package equipment

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	return mux
}

// listModels: GET /v1/equipment/models?category=&status= รายการรุ่นพร้อมจำนวนว่าง/ทั้งหมด (FR-EQS-01, FR-EQS-02, AC-EQS-01, AC-EQS-02)
func (h *Handler) listModels(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListModels(r.Context(), filterFrom(r))
	if err != nil {
		writeUpstreamError(w, err)
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
		writeUpstreamError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ListDTO[ItemRow]{Items: nonNil(items), Total: len(items)})
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

// writeUpstreamError ตอบเมื่ออ่านฐานข้อมูลไม่สำเร็จ ตาม Error contract ใน plan.md ข้อ 5 (FR-EQS-08)
// ไม่ส่งข้อมูลเก่าหรือรายละเอียดภายในกลับไป (test และหน้าจอของกรณีนี้อยู่ใน T-20)
func writeUpstreamError(w http.ResponseWriter, err error) {
	log.Printf("equipment: database error: %v", err)
	writeJSON(w, http.StatusServiceUnavailable, ErrorDTO{Code: "UPSTREAM_UNAVAILABLE", Retryable: true})
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
