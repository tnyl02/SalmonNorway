package borrowrequest

import "time"

// CreateRequestDTO คือข้อมูลที่หน้าจอส่งมาตอนกดส่งคำขอยืม (FR-BRQ-04)
// ไม่มีช่องรหัสผู้ใช้ เพราะระบบใช้รหัสจากการยืนยันตัวตนเสมอ (IF-IDP-01, ASM-BRQ-06)
// การตรวจความครบถ้วนของแต่ละช่องเป็นงานของ T-06 (FR-BRQ-03)
type CreateRequestDTO struct {
	AssetCode     string    `json:"assetCode"`
	RequesterName string    `json:"requesterName"`
	StudentCode   string    `json:"studentCode"`
	LoanType      string    `json:"loanType"`
	Purpose       string    `json:"purpose"`
	PurposeOther  string    `json:"purposeOther,omitempty"`
	StartAt       time.Time `json:"startAt"`
	EndAt         time.Time `json:"endAt"`
}

// toNewRequest แปลงข้อมูลจากหน้าจอเป็นโครงที่ service ใช้
func (d CreateRequestDTO) toNewRequest() NewRequest {
	return NewRequest{
		AssetCode:     d.AssetCode,
		RequesterName: d.RequesterName,
		StudentCode:   d.StudentCode,
		LoanType:      d.LoanType,
		Purpose:       d.Purpose,
		PurposeOther:  d.PurposeOther,
		StartAt:       d.StartAt,
		EndAt:         d.EndAt,
	}
}

// CreatedDTO คือคำตอบเมื่อสร้างคำขอสำเร็จ (FR-BRQ-04, plan.md ข้อ 4)
type CreatedDTO struct {
	RequestID   int       `json:"requestId"`
	Status      string    `json:"status"`      // ค่าเดียวที่ฟีเจอร์นี้ตั้งได้คือ requested
	SubmittedAt time.Time `json:"submittedAt"` // เวลาที่ส่งตามเขตเวลา Asia/Bangkok (glossary)
}

// ErrorDTO คือรูปแบบคำตอบเมื่อเกิดข้อผิดพลาด ใช้ร่วมทุกฟีเจอร์ (CON-TECH-03, CR-004)
// หน้าจอตัดสินใจจากช่อง code ส่วน message เป็นข้อความภาษาไทยสำหรับผู้ใช้
type ErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
