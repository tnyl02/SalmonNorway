package borrowrequest

import (
	"time"

	"salmonnorway/backend/internal/config"
)

// Clock บอกเวลาปัจจุบันให้กฎที่ต้องเทียบเวลา แยกเป็น interface เพื่อให้ test ตั้งเวลาเองได้
// กฎคืนเกินกำหนดเทียบกับ "ตอนนี้" (FR-BRQ-06) ถ้าผูกกับนาฬิกาจริง test จะได้ผลไม่เหมือนกันในแต่ละช่วงวัน
type Clock interface {
	Now() time.Time
}

// SystemClock ใช้นาฬิกาจริงของเครื่อง เป็นตัวที่ใช้ตอนรันของจริง
type SystemClock struct{}

// Now คืนเวลาปัจจุบันของเครื่อง
func (SystemClock) Now() time.Time { return time.Now() }

// FixedClock คือนาฬิกาที่หยุดอยู่กับที่ ใช้ใน test เพื่อตรวจกฎเวลาโดยไม่ต้องรอเวลาจริง (AC-BRQ-06)
type FixedClock struct {
	T time.Time
}

// Now คืนเวลาที่ตั้งไว้
func (c FixedClock) Now() time.Time { return c.T }

// overdueDeadline คืนเวลาที่เริ่มถือว่าเกินกำหนด คือ 16:00 น. ของวันกำหนดคืน (FR-BRQ-06 ข้อ ก, DOM-LOAN-01)
// ใช้ config.LoanNormalReturnHour ตัวเดียวกับกฎเวลาคืนใน validate.go ไม่ฝังเลข 16 ซ้ำ (AGENTS.md กฎข้อ 9)
// เทียบด้วยเขตเวลา Asia/Bangkok ตาม glossary ไม่ใช่เวลาของเครื่องที่รัน
func overdueDeadline(dueAt time.Time, loc *time.Location) time.Time {
	d := dueAt.In(loc)
	return time.Date(d.Year(), d.Month(), d.Day(), config.LoanNormalReturnHour, 0, 0, 0, loc)
}
