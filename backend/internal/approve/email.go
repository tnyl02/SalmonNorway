package approve

import (
	"fmt"
	"strings"
	"time"

	"salmonnorway/backend/internal/config"
)

// EmailData คือข้อมูลของคำขอและผลการพิจารณาที่ใช้สร้าง Email (FR-APV-02 ก ข ค)
type EmailData struct {
	RequestID     int       // เลขที่คำขอ (ค)
	RequesterName string    // ชื่อผู้ขอ (ก ข ค)
	ModelName     string    // รายการอุปกรณ์ที่ขอยืม: ชื่อรุ่น (ก ข ค)
	AssetCode     string    // รายการอุปกรณ์ที่ขอยืม: รหัสครุภัณฑ์ (ก ข ค)
	Purpose       string    // วัตถุประสงค์ 1 ใน 3 ค่า (ก)
	PurposeOther  string    // ข้อความที่ผู้ขอพิมพ์เมื่อเลือก "อื่น ๆ" (ก)
	StartAt       time.Time // วันเวลาที่ต้องการยืม (ก)
	EndAt         time.Time // วันเวลากำหนดคืน = เวลาสิ้นสุดที่ผู้ขอระบุ (ก)
	SubmittedAt   time.Time // วันที่ยื่นคำขอ (ข ค)
	Reason        string    // เหตุผลในการปฏิเสธ (ข)
	ActorName     string    // ชื่อผู้อนุมัติ / ผู้พิจารณา (ก ข) Principal ยังไม่มีชื่อ ใช้รหัสเจ้าหน้าที่ไปก่อน (T-38)
}

// BuildEmail สร้างหัวข้อและเนื้อ Email ตามแบบ (config.ApprovalKind*) (FR-APV-02 ก ข ค, AC-APV-04)
// หัวข้อจาก config.ApprovalEmailSubjects และเนื้อเป็น "หัวข้อข้อมูล: ค่า" บรรทัดละ 1 ช่อง
// เรียงตาม config.ApprovalEmailTemplates (ทีมตัดสินรูปแบบ 2569-11-12 ตอน T-14) ไม่มีคำขึ้นต้นหรือลงท้ายที่ spec ไม่ได้เขียน
func BuildEmail(kind string, d EmailData) (subject, body string, err error) {
	subject, ok := config.ApprovalEmailSubjects[kind]
	if !ok {
		return "", "", fmt.Errorf("approve: ไม่รู้จัก Email แบบ %q", kind)
	}
	values := emailValues(kind, d)
	var b strings.Builder
	for _, label := range config.ApprovalEmailTemplates[kind] {
		v, ok := values[label]
		if !ok {
			return "", "", fmt.Errorf("approve: Email แบบ %q ไม่มีค่าของ %q", kind, label)
		}
		fmt.Fprintf(&b, "%s: %s\n", label, v)
	}
	return subject, b.String(), nil
}

// emailValues คืนค่าของแต่ละหัวข้อข้อมูล ชื่อหัวข้อตรงกับ config.ApprovalEmailTemplates
func emailValues(kind string, d EmailData) map[string]string {
	item := fmt.Sprintf("%s (%s)", d.ModelName, d.AssetCode)
	switch kind {
	case config.ApprovalKindApproved: // FR-APV-02 (ก) ยืมปกติไม่แสดงจำนวนอุปกรณ์
		return map[string]string{
			"ชื่อผู้ขอ":             d.RequesterName,
			"รายการอุปกรณ์ที่ขอยืม": item,
			"วัตถุประสงค์":          purposeText(d),
			"วันเวลาที่ต้องการยืม":  formatDateTime(d.StartAt),
			"วันเวลากำหนดคืน":       formatDateTime(d.EndAt),
			"ชื่อผู้อนุมัติ":        d.ActorName,
			"สถานที่รับ":            config.ApprovalPickupLocation,
		}
	case config.ApprovalKindRejected: // FR-APV-02 (ข) รวมคำขอที่ถูกปฏิเสธตาม FR-APV-01 (จ)
		return map[string]string{
			"ชื่อผู้ขอ":             d.RequesterName,
			"รายการอุปกรณ์ที่ขอยืม": item,
			"วันที่ยื่นคำขอ":        formatDate(d.SubmittedAt),
			"เหตุผลในการปฏิเสธ":     d.Reason,
			"ชื่อผู้พิจารณา":        d.ActorName,
		}
	default: // config.ApprovalKindAutoRejected FR-APV-02 (ค) ข้อความตายตัวจาก config
		return map[string]string{
			"เลขที่คำขอ":            fmt.Sprint(d.RequestID),
			"ชื่อผู้ขอ":             d.RequesterName,
			"รายการอุปกรณ์ที่ขอยืม": item,
			"วันที่ยื่นคำขอ":        formatDate(d.SubmittedAt),
			"สถานะ":                 config.ApprovalAutoRejectEmail.Status,
			"เหตุผล":                config.ApprovalAutoRejectReason,
			"เวลาที่ดำเนินการ":      config.ApprovalAutoRejectEmail.DoneAt,
			"แนวทาง":                config.ApprovalAutoRejectEmail.Guidance,
		}
	}
}

// purposeText คือวัตถุประสงค์ภาษาไทย ถ้าเลือก "อื่น ๆ" แสดงข้อความที่ผู้ขอพิมพ์ (FR-BRQ-02 ก)
func purposeText(d EmailData) string {
	if d.Purpose == config.PurposeOther && d.PurposeOther != "" {
		return d.PurposeOther
	}
	if s, ok := config.ApprovalPurposeLabels[d.Purpose]; ok {
		return s
	}
	return d.Purpose
}

// buddhistYearOffset คือ ปี พ.ศ. = ปี ค.ศ. + 543 (ค่าปฏิทิน ไม่ใช่ค่าจาก spec)
const buddhistYearOffset = 543

// localTime แปลงเป็นเวลาไทย (config.TimeZone, glossary)
func localTime(t time.Time) time.Time {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		loc = time.UTC // ไม่ควรเกิด เพราะ config ฝัง tzdata ไว้แล้ว
	}
	return t.In(loc)
}

// formatDateTime แสดงวันเวลาเป็น พ.ศ. ไม่มีวินาที เช่น "2569-11-12 09:00"
// รูปแบบเดียวกับ FormatNoticeTime ของ 003 และเวลาบนหน้าเว็บของ 001 (ทีมตัดสิน 2569-11-12)
func formatDateTime(t time.Time) string {
	lt := localTime(t)
	return fmt.Sprintf("%d-%s", lt.Year()+buddhistYearOffset, lt.Format("01-02 15:04"))
}

// formatDate แสดงวันที่เป็น พ.ศ. เช่น "2569-11-12" ใช้กับ "วันที่ยื่นคำขอ" (FR-APV-02 ข ค)
func formatDate(t time.Time) string {
	lt := localTime(t)
	return fmt.Sprintf("%d-%s", lt.Year()+buddhistYearOffset, lt.Format("01-02"))
}
