// Package cancelrequest คือฟีเจอร์ 003 ยกเลิกคำขอยืม (UC-15, specs/003-cancel-request/spec.md)
package cancelrequest

import "time"

// PendingRequestDTO คือคำขอ 1 รายการในหน้า "คำขอที่รออนุมัติของฉัน" (FR-CXL-06, plan.md ข้อ 4)
// ส่งเฉพาะข้อมูลที่หน้าจอใช้ ไม่ส่งรหัสผู้ขอหรือข้อมูลของผู้อื่น (FR-CXL-03)
type PendingRequestDTO struct {
	RequestID   int       `json:"requestId"`
	ModelName   string    `json:"modelName"`
	AssetCode   string    `json:"assetCode"`
	SubmittedAt time.Time `json:"submittedAt"` // เวลา UTC แบบ ISO 8601 หน้าจอแปลงเป็น Asia/Bangkok เอง
}

// PendingListDTO คือ body ของ GET /v1/my/pending-requests ไม่มีรายการ = items เป็น [] (FR-CXL-06, AC-CXL-12)
type PendingListDTO struct {
	Items []PendingRequestDTO `json:"items"`
}

// CancelRequestDTO คือ body ของ POST /v1/borrow-requests/{id}/cancel (FR-CXL-08)
type CancelRequestDTO struct {
	Reason string `json:"reason"`
}

// CancelResultDTO คือ body เมื่อยกเลิกสำเร็จ (FR-CXL-01, plan.md ข้อ 4)
type CancelResultDTO struct {
	RequestID int    `json:"requestId"`
	Status    string `json:"status"`
}

// ErrorDTO คือ body เมื่อตอบไม่สำเร็จ ตามตาราง error ใน plan.md ข้อ 4
type ErrorDTO struct {
	Code      string `json:"code"`
	Retryable bool   `json:"retryable,omitempty"`
}
