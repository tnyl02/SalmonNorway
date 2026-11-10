package returns

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
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
// ทุก endpoint ต้องผ่าน auth.Middleware ก่อนเสมอ (IF-IDP-01, NFR-SEC-01)
func (h *Handler) Routes(provider auth.Provider) func(*http.ServeMux) {
	return func(mux *http.ServeMux) {
		for pattern, handle := range map[string]http.HandlerFunc{
			"GET /api/v1/returns/open-loans":                      notImplemented, // FR-RET-01(ข), FR-RET-06 ทำใน T-17
			"GET /api/v1/returns/lookup":                          notImplemented, // FR-RET-01(ก) ทำใน T-05
			"POST /api/v1/returns":                                h.recordReturn, // FR-RET-02 ถึง 07
			"GET /api/v1/returns/loans/{loanId}":                  notImplemented, // FR-RET-06 ทำใน T-14
			"GET /api/v1/returns/loans/{loanId}/photos/{photoId}": notImplemented, // FR-RET-05 ทำใน T-14
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

// recordReturn รับ multipart: assetCode, condition, photos[] แล้วบันทึกการรับคืน (FR-RET-02, FR-RET-03, FR-RET-05)
// ไม่อ่านฟิลด์เวลาใด ๆ จาก client เวลาคืนมาจากฐานข้อมูลเท่านั้น (FR-RET-02, ASM-04)
func (h *Handler) recordReturn(w http.ResponseWriter, r *http.Request) {
	// เก็บในหน่วยความจำไม่เกินขนาดรูป 1 รูป ส่วนเกินพักลงไฟล์ชั่วคราว (ค่าจาก config.go)
	// การจำกัดขนาด request ทั้งก้อนทำใน T-12
	if err := r.ParseMultipartForm(config.ReturnPhotoMaxBytes); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"code": "BAD_REQUEST"})
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll() //nolint:errcheck
	}

	photos, err := readPhotos(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"code": "BAD_REQUEST"})
		return
	}

	rec, err := h.svc.RecordReturn(r.Context(), r.FormValue("assetCode"), r.FormValue("condition"), photos)
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

// readPhotos อ่านไฟล์ทุกไฟล์ในฟิลด์ photos (FR-RET-05)
// ชนิดไฟล์ดูจาก magic bytes ไม่เชื่อนามสกุลหรือ header ของ client (plan.md ข้อ 4)
// การตรวจจำนวน ชนิด และขนาดให้ตรง FR-RET-05 ทำใน T-11
func readPhotos(r *http.Request) ([]Photo, error) {
	if r.MultipartForm == nil {
		return nil, nil
	}
	var photos []Photo
	for _, fh := range r.MultipartForm.File["photos"] {
		f, err := fh.Open()
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			return nil, err
		}
		photos = append(photos, Photo{ContentType: http.DetectContentType(data), Data: data})
	}
	return photos, nil
}

// notImplemented ตอบ 501 ให้ endpoint ที่ยังไม่ถึง task ของมัน
func notImplemented(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{"code": "NOT_IMPLEMENTED"})
}

// writeJSON ตอบ JSON และตั้ง Cache-Control: no-store ทุก response ของ 006 (plan.md ข้อ 4, NFR-DATA-01)
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
