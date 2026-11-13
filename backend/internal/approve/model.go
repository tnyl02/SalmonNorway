// Package approve คือฟีเจอร์ 004 อนุมัติคำขอยืมปกติ (UC-12, specs/004-approve-request/)
//
// เจ้าหน้าที่ตรวจและอนุมัติหรือปฏิเสธคำขอยืมแบบปกติของภาควิชาตนเอง (FR-APV-01) และแจ้งผลผู้ขอทาง Email (FR-APV-02)
// ฟีเจอร์นี้เปลี่ยนแค่สถานะของคำขอ **ไม่เปลี่ยนสถานะอุปกรณ์** (CON-DATA-01, data-model.md แถว UC-12)
// กฎทุกข้อบังคับที่ฝั่ง server ไม่พึ่งการซ่อนปุ่มที่หน้าจอ (plan.md ข้อ 1)
package approve

import "time"

// PendingRequest คือคำขอยืมที่รอพิจารณา 1 แถวในหน้ารายการของเจ้าหน้าที่ ทั้งยืมปกติและค้างคืน (FR-APV-01 ก ข, FR-APV-03 ก)
type PendingRequest struct {
	ID            int       // รหัสคำขอ (DM-06)
	RequesterName string    // ชื่อผู้ขอ (FR-APV-01 ก) ผู้ขอพิมพ์เองตอนส่งคำขอ (FR-BRQ-03 ก)
	StudentCode   string    // รหัสนักศึกษาที่ผู้ขอพิมพ์ (ASM-03)
	ModelName     string    // อุปกรณ์: ชื่อรุ่น (FR-APV-01 ก)
	AssetCode     string    // อุปกรณ์: รหัสครุภัณฑ์ (CON-DATA-03)
	Purpose       string    // วัตถุประสงค์ 1 ใน 3 ค่า (FR-APV-01 ก, FR-BRQ-02)
	PurposeOther  string    // ข้อความที่ผู้ขอพิมพ์เมื่อเลือก "อื่น ๆ" ว่างได้
	StartAt       time.Time // ช่วงเวลาที่ต้องการยืม: เริ่ม (FR-APV-01 ก)
	EndAt         time.Time // ช่วงเวลาที่ต้องการยืม: คืน (FR-APV-01 ก)
	SubmittedAt   time.Time // เวลาที่ส่ง ใช้เรียงลำดับ (FR-APV-01 ก, DOM-REQ-01)
	LoanType      string    // normal / overnight (config.LoanTypeNormal, config.LoanTypeOvernight) ใช้แยกกฎ (FR-APV-03, plan.md v2)
	IsFirst       bool      // เป็นคำขอแรกของอุปกรณ์รายชิ้นนั้นหรือไม่ อนุมัติได้เฉพาะคำขอแรก คิดเฉพาะยืมปกติ (FR-APV-01 ข, Q-APV-23)
	HasAttachment bool      // มีไฟล์ใบบันทึกข้อความที่ 002 เก็บไว้ (FR-APV-03 ก, Q-APV-41)
}

// requestRow คือคำขอ 1 ใบที่ล็อกไว้ระหว่างพิจารณา พร้อมข้อมูลอุปกรณ์ที่ต้องใช้ตรวจ (FR-APV-01 ข ค, ASM-06)
type requestRow struct {
	ID           int
	ItemID       int
	LoanType     string
	Status       string    // สถานะของคำขอ (data-model.md v9)
	StartAt      time.Time // ใช้ตรวจเวลาตัด 16:00 น. ของวันที่เริ่มยืม (T-21)
	SubmittedAt  time.Time
	AssetCode    string
	ItemStatus   string // สถานะอุปกรณ์ (CON-DATA-01) ใช้ตรวจ "ว่าง" (T-19)
	DepartmentID int    // ภาควิชาเจ้าของอุปกรณ์ (DM-03, ASM-06)
}

// decision คือผลการพิจารณาที่บันทึกลง approval_decision (FR-APV-01 ฌ, plan.md ข้อ 3)
type decision struct {
	RequestID         int
	Result            string // config.ApprovalKind*
	Reason            string // ว่างได้เฉพาะการอนุมัติ
	ActorID           string // รหัสเจ้าหน้าที่ หรือ config.ApprovalSystemActorID
	ActorName         string // ว่างได้ Principal ยังไม่มีชื่อ (plan.md ข้อ 8, T-38)
	ActorDepartmentID int
	AdvisorName       string // ชื่ออาจารย์ที่เซ็น เฉพาะอนุมัติคำขอค้างคืน ว่าง = NULL (FR-APV-03 ข ช, migration 0011)
}

// Result คือผลของการอนุมัติหรือปฏิเสธที่ส่งกลับให้ handler (plan.md ข้อ 4)
type Result struct {
	RequestID  int
	Status     string    // approved / rejected
	DecidedAt  time.Time // เวลาที่บันทึกผล (FR-APV-01 ฌ)
	RejectedBy []int     // รหัสคำขออื่นของชิ้นเดียวกันที่ถูกปฏิเสธตาม (FR-APV-01 จ, T-16) ว่างได้
}

// FailedEmail คือ Email ที่ส่งไม่สำเร็จครบทุกครั้งแล้ว 1 แถวในรายการ "Email ที่ส่งไม่สำเร็จ" (FR-APV-02 ง, ASM-08)
type FailedEmail struct {
	ID            int
	RequestID     int
	Kind          string // approved / rejected / auto_rejected
	RequesterName string
	ModelName     string
	AssetCode     string
	Attempts      int
	LastError     string    // สาเหตุ (FR-APV-02 ง)
	LastAttemptAt time.Time // เวลา (FR-APV-02 ง)
}
