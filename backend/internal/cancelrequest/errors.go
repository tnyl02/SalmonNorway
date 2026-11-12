package cancelrequest

import (
	"fmt"
	"net/http"

	"salmonnorway/backend/internal/config"
)

// รหัส error ของ 003 ตามรายการกลาง CON-TECH-03 (specs/000-shared/constraints.md v5.2, CR-006)
const (
	codeUnauthenticated     = "UNAUTHENTICATED"
	codeForbiddenRole       = "FORBIDDEN_ROLE"
	codeValidationFailed    = "VALIDATION_FAILED"
	codeRequestNotFound     = "REQUEST_NOT_FOUND"
	codeAlreadyApproved     = "ALREADY_APPROVED"
	codeNotCancellable      = "NOT_CANCELLABLE"
	codeUpstreamUnavailable = "UPSTREAM_UNAVAILABLE"
)

// fieldReason คือชื่อช่องที่แนบใน VALIDATION_FAILED (CON-TECH-03 แนบรายชื่อช่องที่ผิด)
const fieldReason = "reason"

// ข้อความของแต่ละรหัส (CR-006) "ไม่พบคำขอ" มาจาก FR-CXL-03 ข้อความอื่นทีมเลือกเอง ไม่ได้มาจาก spec
var (
	msgUnauthenticated     = "กรุณาเข้าสู่ระบบก่อนใช้งาน"
	msgForbiddenRole       = "เฉพาะนักศึกษาเท่านั้นที่ยกเลิกคำขอของตนเองได้"
	msgRequestNotFound     = "ไม่พบคำขอ"
	msgAlreadyApproved     = "ยกเลิกไม่ได้ เพราะคำขอได้รับการอนุมัติแล้ว"
	msgNotCancellable      = "คำขอนี้ยกเลิกไม่ได้"
	msgUpstreamUnavailable = "ระบบไม่ตอบสนอง กรุณาลองใหม่อีกครั้ง"
)

// msgValidationFailed คือข้อความเดียวของเหตุผลว่างหรือยาวเกิน ตัวเลขมาจาก config.CancelReasonMaxLen (FR-CXL-08, CR-006)
func msgValidationFailed() string {
	return fmt.Sprintf("กรุณาระบุเหตุผลการยกเลิกไม่เกิน %d ตัวอักษร", config.CancelReasonMaxLen)
}

// StatusError คือคำขอของตนที่ยกเลิกไม่ได้ พร้อมสถานะปัจจุบันของคำขอ ใช้แนบ status ใน 409 (CR-006)
// Err เป็น ErrAlreadyApproved หรือ ErrNotCancelled ใช้ errors.Is ได้
type StatusError struct {
	Err    error
	Status string
}

func (e *StatusError) Error() string { return e.Err.Error() + ": " + e.Status }
func (e *StatusError) Unwrap() error { return e.Err }

func writeError(w http.ResponseWriter, httpStatus int, body ErrorDTO) {
	writeJSON(w, httpStatus, body)
}

func errUnauthenticated() ErrorDTO {
	return ErrorDTO{Code: codeUnauthenticated, Message: msgUnauthenticated}
}

func errForbiddenRole() ErrorDTO {
	return ErrorDTO{Code: codeForbiddenRole, Message: msgForbiddenRole}
}

func errValidationFailed() ErrorDTO {
	return ErrorDTO{Code: codeValidationFailed, Message: msgValidationFailed(), Fields: []string{fieldReason}}
}

func errRequestNotFound() ErrorDTO {
	return ErrorDTO{Code: codeRequestNotFound, Message: msgRequestNotFound}
}

func errUpstreamUnavailable() ErrorDTO {
	return ErrorDTO{Code: codeUpstreamUnavailable, Message: msgUpstreamUnavailable}
}
