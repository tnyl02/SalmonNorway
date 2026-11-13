package approve

import (
	"strconv"
	"time"

	"salmonnorway/backend/internal/borrowrequest"
	"salmonnorway/backend/internal/config"
)

// PendingRequestDTO คือ 1 แถวในหน้ารายการคำขอของเจ้าหน้าที่ (FR-APV-01 ก ข, plan.md ข้อ 4)
// ชื่อช่องเป็น camelCase แบบเดียวกับ API ของ 002
// ไม่ส่งรหัสนักศึกษา เพราะ FR-APV-01 (ก) ไม่ได้สั่งให้แสดง (rtm.md F-05 ทีมตัดสินแก้โค้ด 2569-11-13)
type PendingRequestDTO struct {
	RequestID     int       `json:"requestId"`
	RequesterName string    `json:"requesterName"`
	ModelName     string    `json:"modelName"`
	AssetCode     string    `json:"assetCode"`
	Purpose       string    `json:"purpose"`
	PurposeOther  string    `json:"purposeOther,omitempty"`
	StartAt       time.Time `json:"startAt"`
	EndAt         time.Time `json:"endAt"`
	SubmittedAt   time.Time `json:"submittedAt"`
	LoanType      string    `json:"loanType"` // normal / overnight (FR-APV-03 ก, plan.md v2 ข้อ 4)
	IsFirst       bool      `json:"isFirst"` // หน้าจอแสดงปุ่มอนุมัติเฉพาะแถวนี้ (FR-APV-01 ข, T-18) server ตรวจซ้ำเสมอ (T-17) คำขอค้างคืนเป็น false เสมอ เพราะไม่มีคิว
	// AttachmentURL คือลิงก์เปิดไฟล์ใบบันทึกข้อความผ่าน endpoint ของ 002 (FR-APV-03 ก, FR-BRQ-10 ก, Q-APV-41) มีเฉพาะคำขอค้างคืนที่มีไฟล์
	AttachmentURL string `json:"attachmentUrl,omitempty"`
}

// PendingListDTO คือคำตอบของ GET /v1/approvals/requests
// loadedAt ให้หน้าจอส่งกลับมาตอนอนุมัติหรือปฏิเสธ (plan.md ข้อ 4)
type PendingListDTO struct {
	Items    []PendingRequestDTO `json:"items"`
	LoadedAt time.Time           `json:"loadedAt"`
}

// ErrorDTO คือรูปแบบ error กลาง {code, message} (CON-TECH-03)
type ErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func toPendingDTO(p PendingRequest) PendingRequestDTO {
	return PendingRequestDTO{
		RequestID:     p.ID,
		RequesterName: p.RequesterName,
		ModelName:     p.ModelName,
		AssetCode:     p.AssetCode,
		Purpose:       p.Purpose,
		PurposeOther:  p.PurposeOther,
		StartAt:       p.StartAt,
		EndAt:         p.EndAt,
		SubmittedAt:   p.SubmittedAt,
		LoanType:      p.LoanType,
		IsFirst:       p.IsFirst,
		AttachmentURL: attachmentURL(p),
	}
}

// attachmentURL ใช้ path ของ 002 (GET /v1/loan-requests/{id}/attachment) ไม่เขียน path ซ้ำในฟีเจอร์นี้
func attachmentURL(p PendingRequest) string {
	if p.LoanType != config.LoanTypeOvernight || !p.HasAttachment {
		return ""
	}
	return borrowrequest.PathPrefix + "/" + strconv.Itoa(p.ID) + "/attachment"
}

// DecideInputDTO คือ body ของการอนุมัติและการปฏิเสธ (plan.md ข้อ 4)
type DecideInputDTO struct {
	LoadedAt time.Time `json:"loadedAt"`         // เวลาที่หน้าจอโหลดรายการ ใช้ใน T-23
	Reason   string    `json:"reason,omitempty"` // เหตุผลการปฏิเสธ (FR-APV-01 ง) การอนุมัติไม่ใช้ช่องนี้
}

// ResultDTO คือคำตอบเมื่ออนุมัติหรือปฏิเสธสำเร็จ (plan.md ข้อ 4)
type ResultDTO struct {
	RequestID          int       `json:"requestId"`
	Status             string    `json:"status"`             // approved / rejected
	DecidedAt          time.Time `json:"decidedAt"`          // FR-APV-01 ฌ
	RejectedRequestIDs []int     `json:"rejectedRequestIds"` // คำขออื่นที่ถูกปฏิเสธตาม (FR-APV-01 จ, T-16) ว่างได้
}

func toResultDTO(r Result) ResultDTO {
	ids := r.RejectedBy
	if ids == nil {
		ids = []int{}
	}
	return ResultDTO{RequestID: r.RequestID, Status: r.Status, DecidedAt: r.DecidedAt, RejectedRequestIDs: ids}
}

// ValidationErrorDTO คือ {code, message} ที่แนบรายชื่อช่องที่ผิด CON-TECH-03 กำหนดให้ VALIDATION_FAILED แนบช่องมาด้วย
// ชื่อช่อง fields แบบเดียวกับ 002
type ValidationErrorDTO struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Fields  []string `json:"fields"`
}

// FailedEmailDTO คือ 1 แถวในรายการ "Email ที่ส่งไม่สำเร็จ" (FR-APV-02 ง, plan.md ข้อ 4)
type FailedEmailDTO struct {
	EmailID       int       `json:"emailId"`
	RequestID     int       `json:"requestId"`
	Kind          string    `json:"kind"`
	RequesterName string    `json:"requesterName"`
	ModelName     string    `json:"modelName"`
	AssetCode     string    `json:"assetCode"`
	Attempts      int       `json:"attempts"`
	LastError     string    `json:"lastError"`     // สาเหตุ
	LastAttemptAt time.Time `json:"lastAttemptAt"` // เวลา
}

// FailedEmailListDTO คือคำตอบของ GET /v1/approvals/failed-emails
type FailedEmailListDTO struct {
	Items []FailedEmailDTO `json:"items"`
}

// ResendResultDTO คือคำตอบของ POST /v1/approvals/emails/{id}/resend สถานะการส่งหลังกด (pending = เข้าคิวแล้ว)
type ResendResultDTO struct {
	EmailID int    `json:"emailId"`
	Status  string `json:"status"`
}
