package handover

import (
	"encoding/json"
	"log"
	"net/http"
)

// รหัส error ที่ 005 ใช้ ทุกตัวอยู่ในรายการกลางของ CON-TECH-03 (specs/000-shared/constraints.md, CR-004)
const (
	codeForbiddenRole       = "FORBIDDEN_ROLE"       // 403 CON-HND-01, AC-HND-09 (ลงทะเบียนโดย 005, Q-P4)
	codeRequestNotApproved  = "REQUEST_NOT_APPROVED" // 409 FR-HND-01, AC-HND-01 (ลงทะเบียนโดย 005, Q-P4)
	codeBorrowerMismatch    = "BORROWER_MISMATCH"    // 409 FR-HND-02, AC-HND-10 (ลงทะเบียนโดย 005, Q-P4)
	codeItemNotAvailable    = "ITEM_NOT_AVAILABLE"   // 409 FR-HND-04, AC-HND-04, 13
	codeValidationFailed    = "VALIDATION_FAILED"    // 422 FR-HND-03, 06 (ตรวจบัตร, รูป)
	codeUpstreamUnavailable = "UPSTREAM_UNAVAILABLE" // 503 ฐานข้อมูลไม่ตอบ
)

// ข้อความภาษาไทยของแต่ละรหัส (CON-TECH-03 ต้องมี message) ถ้อยคำทีมเลือกเอง ไม่ได้มาจาก spec
// หน้าจอตัดสินใจจาก code ไม่ใช่จากข้อความ
var errorMessages = map[string]string{
	codeForbiddenRole:       "เฉพาะเจ้าหน้าที่เท่านั้นที่บันทึกการส่งมอบได้",
	codeRequestNotApproved:  "คำขอนี้ยังไม่ได้รับอนุมัติ หรือส่งมอบไปแล้ว",
	codeBorrowerMismatch:    "ผู้ที่ยืนยันตัวตนไม่ใช่ผู้ส่งคำขอนี้",
	codeItemNotAvailable:    "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่าง",
	codeValidationFailed:    "ข้อมูลไม่ครบหรือไม่ถูกต้อง",
	codeUpstreamUnavailable: "ระบบไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่",
}

// ErrorDTO คือ body ของ error ทุกตัว `{code, message}` (CON-TECH-03) fields ใช้กับ VALIDATION_FAILED บอกช่องที่ผิด
type ErrorDTO struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Fields  []string `json:"fields,omitempty"`
}

// writeError ตอบ error ตาม CON-TECH-03
func writeError(w http.ResponseWriter, status int, code string, fields ...string) {
	writeJSON(w, status, ErrorDTO{Code: code, Message: errorMessages[code], Fields: fields})
}

// writeUpstreamError ตอบ 503 เมื่อฐานข้อมูลไม่ตอบ ไม่ส่งรายละเอียดภายในกลับไป
func writeUpstreamError(w http.ResponseWriter, err error) {
	log.Printf("handover: database error: %v", err)
	writeError(w, http.StatusServiceUnavailable, codeUpstreamUnavailable)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
