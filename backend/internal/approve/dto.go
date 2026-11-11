package approve

import "time"

// PendingRequestDTO คือ 1 แถวในหน้ารายการคำขอของเจ้าหน้าที่ (FR-APV-01 ก ข, plan.md ข้อ 4)
// ชื่อช่องเป็น camelCase แบบเดียวกับ API ของ 002
type PendingRequestDTO struct {
	RequestID     int       `json:"requestId"`
	RequesterName string    `json:"requesterName"`
	StudentCode   string    `json:"studentCode"`
	ModelName     string    `json:"modelName"`
	AssetCode     string    `json:"assetCode"`
	Purpose       string    `json:"purpose"`
	PurposeOther  string    `json:"purposeOther,omitempty"`
	StartAt       time.Time `json:"startAt"`
	EndAt         time.Time `json:"endAt"`
	SubmittedAt   time.Time `json:"submittedAt"`
	IsFirst       bool      `json:"isFirst"` // หน้าจอแสดงปุ่มอนุมัติเฉพาะแถวนี้ (FR-APV-01 ข, T-18) server ตรวจซ้ำเสมอ (T-17)
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
		StudentCode:   p.StudentCode,
		ModelName:     p.ModelName,
		AssetCode:     p.AssetCode,
		Purpose:       p.Purpose,
		PurposeOther:  p.PurposeOther,
		StartAt:       p.StartAt,
		EndAt:         p.EndAt,
		SubmittedAt:   p.SubmittedAt,
		IsFirst:       p.IsFirst,
	}
}
