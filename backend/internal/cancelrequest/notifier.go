package cancelrequest

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"salmonnorway/backend/internal/config"
)

// CancelNotice คือเนื้อหา Email แจ้งเจ้าหน้าที่เมื่อมีการยกเลิกคำขอ (FR-CXL-09, AC-CXL-15)
// ระบุผู้ยกเลิก ชื่อรุ่น รหัสครุภัณฑ์ เวลาที่ยกเลิก และเหตุผล (ASM-08)
type CancelNotice struct {
	To            []string  // Email เจ้าหน้าที่ของภาควิชาเจ้าของอุปกรณ์ (staff_contact, ASM-10)
	RequestID     int       // ใช้อ้างอิงเท่านั้น
	CancelledBy   string    // รหัสผู้ใช้ของผู้ยกเลิก (DM-08 Principal.UserID) มาจากระบบยืนยันตัวตน
	RequesterName string    // ชื่อผู้ขอที่พิมพ์ไว้ในคำขอ (requester_name ของ 002) ทีมตัดสินให้แสดงคู่กับรหัสผู้ใช้ 2569-11-12
	StudentCode   string    // รหัสนักศึกษาที่พิมพ์ไว้ในคำขอ (student_code ของ 002)
	ModelName     string    // ชื่อรุ่น
	AssetCode     string    // รหัสครุภัณฑ์
	CancelledAt   time.Time // เวลาที่ยกเลิก (cancelled_at ในฐานข้อมูล)
	Reason        string    // เหตุผลการยกเลิก (FR-CXL-08)
}

// Notifier คือช่องทางส่ง Email ผ่าน IF-NOT-01 ออกแบบเป็น interface แบบเดียวกับ auth.Provider (plan.md ข้อ 2)
// เปลี่ยนเป็นตัวส่งจริงได้โดยไม่แก้ service เมื่อได้คำตอบ Q-15
type Notifier interface {
	NotifyCancelled(ctx context.Context, n CancelNotice) error
}

// LogNotifier คือ Notifier จำลองของ prototype เขียนเนื้อหา Email ลง log แทนการส่งจริง (ASM-10)
// ทีมเลือกเอง ไม่ได้มาจาก spec
type LogNotifier struct{}

// NotifyCancelled เขียน Email ลง log ของหลังบ้าน
func (LogNotifier) NotifyCancelled(_ context.Context, n CancelNotice) error {
	log.Printf("cancelrequest: [Email จำลอง] ถึง %v: %s ยกเลิกคำขอยืม %s (%s) เมื่อ %s เหตุผล: %s (คำขอ %d)",
		n.To, FormatCanceller(n), n.ModelName, n.AssetCode, FormatNoticeTime(n.CancelledAt), n.Reason, n.RequestID)
	return nil
}

// FormatCanceller แสดงผู้ยกเลิกใน Email เป็นชื่อและรหัสนักศึกษาจากคำขอ คู่กับรหัสผู้ใช้จากระบบยืนยันตัวตน (FR-CXL-09)
// เช่น "ใจดี สวยมาก (รหัสนักศึกษา 660710755, ผู้ใช้ student-001)" ทีมตัดสิน 2569-11-12
// ชื่อและรหัสนักศึกษาเป็นข้อมูลที่ผู้ขอพิมพ์เอง (FR-BRQ-03 ของ 002) รหัสผู้ใช้จึงต้องอยู่ด้วยเพื่อยืนยันตัวตน
func FormatCanceller(n CancelNotice) string {
	return fmt.Sprintf("%s (รหัสนักศึกษา %s, ผู้ใช้ %s)", n.RequesterName, n.StudentCode, n.CancelledBy)
}

// ErrNoRecipients คือภาควิชาเจ้าของอุปกรณ์ไม่มี Email เจ้าหน้าที่ใน staff_contact
// ทีมตัดสิน 2569-11-12 ให้นับเป็นส่ง Email ไม่สำเร็จ (บันทึก audit notify error ใน T-32 เหมือนส่งพลาด)
var ErrNoRecipients = errors.New("no staff email for department")

// buddhistYearOffset คือ ปี พ.ศ. = ปี ค.ศ. + 543 (ค่าปฏิทิน ไม่ใช่ค่าจาก spec)
const buddhistYearOffset = 543

// FormatNoticeTime แสดงเวลาใน Email เป็นเวลาไทย (config.TimeZone, glossary) ปี พ.ศ. ไม่มีวินาที เช่น "2569-11-12 01:25"
// รูปแบบเดียวกับเวลาบนหน้าเว็บ (formatDueAt ของ 001 ทีมตัดสิน 2569-11-10) ทีมขอให้ตรงกันและไม่แสดงวินาที 2569-11-12
func FormatNoticeTime(t time.Time) string {
	loc, err := time.LoadLocation(config.TimeZone)
	if err != nil {
		loc = time.UTC // ไม่ควรเกิด เพราะ config ฝัง tzdata ไว้แล้ว
	}
	lt := t.In(loc)
	return fmt.Sprintf("%d-%s", lt.Year()+buddhistYearOffset, lt.Format("01-02 15:04"))
}
