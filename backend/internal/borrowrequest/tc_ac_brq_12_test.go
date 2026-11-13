package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-12-1 ถึง TC-BRQ-12-4 (AC-BRQ-12, FR-BRQ-06 ข้อ ข)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
//
// ชุดนี้เรียกชั้น service ตรง ๆ ไม่ผ่าน HTTP ต่างจากชุดของ AC-BRQ-06 เพราะต้องใช้ผู้ขอที่ไม่มี
// ข้อมูลยืมอื่นปนเลย แต่ผู้ใช้จำลองของ auth.MockProvider ในภาควิชา 1 มีแค่ student-001 กับ
// student-002 ซึ่งมีรายการยืมค้างอยู่ตาม seed.sql ทั้งคู่ จึงแยกเคสไม่ได้ถ้ายิงผ่าน HTTP
// (วิธีเดียวกับ TestAC_BRQ_12 เดิม) ErrHasOverdueLoan คือ error ที่ handler แปลงเป็น
// 409 HAS_OVERDUE_LOAN ตาม CON-TECH-03 จึงเป็นกฎเดียวกันกับที่แถวในตารางเขียนไว้

import (
	"errors"
	"testing"

	"salmonnorway/backend/internal/config"
)

// TC-BRQ-12-1 (ทางปกติ) คืนของที่เกินกำหนดเมื่อ 14:00 น. ของวันนี้ แล้วส่งคำขอ 23:59 น. วันเดียวกัน
func TestTC_BRQ_12_1_SameDayStillBlocked(t *testing.T) {
	pool := testPool(t)
	const requester = "student-tc-12-1" // ผู้ขอของ test นี้เท่านั้น
	ctx := studentCtx(requester)

	// Given กำหนดคืนเมื่อวาน 16:00 น. แต่เพิ่งคืนเอา 14:00 น. ของวันนี้ (คืนช้า)
	returnedAt := bangkokToday(14, 0)
	insertLoan(t, pool, requester, "CS-NB-0001",
		bangkokToday(config.LoanNormalReturnHour, 0).AddDate(0, 0, -1), &returnedAt)

	// When ส่งคำขอยืมใหม่เวลา 23:59 น. ของวันเดียวกัน
	svc := NewServiceWithClock(NewRepository(pool), FixedClock{T: bangkokToday(23, 59)})
	_, err := svc.Create(ctx, newTestRequest("CS-NB-0003"))

	// Then (1) ถูกปฏิเสธด้วยกฎของค้างเกินกำหนด (handler แปลงเป็น 409 HAS_OVERDUE_LOAN)
	if !errors.Is(err, ErrHasOverdueLoan) {
		t.Fatalf("(1) เวลา 23:59 ของวันที่คืนช้า ต้องได้ ErrHasOverdueLoan แต่ได้ %v", err)
	}
	// Then (2) ไม่มีแถวใหม่
	if n := countRequests(t, pool, requester); n != 0 {
		t.Errorf("(2) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-12-2 (ทางปกติ) เงื่อนไขเดียวกัน แต่ส่งเวลา 00:01 น. ของวันถัดไป
func TestTC_BRQ_12_2_NextDayAllowed(t *testing.T) {
	pool := testPool(t)
	const requester = "student-tc-12-2"
	ctx := studentCtx(requester)

	// Given คืนของที่เกินกำหนดแล้วเมื่อ 14:00 น. ของวันนี้
	returnedAt := bangkokToday(14, 0)
	insertLoan(t, pool, requester, "CS-NB-0001",
		bangkokToday(config.LoanNormalReturnHour, 0).AddDate(0, 0, -1), &returnedAt)

	// When ส่งคำขอยืมใหม่เวลา 00:01 น. ของวันถัดไป (เทียบเขตเวลา Asia/Bangkok)
	svc := NewServiceWithClock(NewRepository(pool), FixedClock{T: bangkokToday(0, 1).AddDate(0, 0, 1)})
	_, err := svc.Create(ctx, newTestRequest("CS-NB-0003"))

	// Then (1) ส่งได้ ไม่มี error
	if err != nil {
		t.Fatalf("(1) เวลา 00:01 ของวันถัดไปต้องส่งคำขอได้ แต่ได้ error: %v", err)
	}
	// Then (2) มีแถวใหม่ 1 แถว
	if n := countRequests(t, pool, requester); n != 1 {
		t.Errorf("(2) ต้องมีคำขอ 1 ใบหลังปลดล็อก แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-12-3 (ขอบ) จุดปลดล็อกตาม FR-BRQ-06 ข คือ 00:00 น. ของวันถัดไป ลองที่เวลานั้นพอดี
func TestTC_BRQ_12_3_ExactMidnight(t *testing.T) {
	pool := testPool(t)
	const requester = "student-tc-12-3"
	ctx := studentCtx(requester)

	// Given คืนของที่เกินกำหนดแล้วเมื่อ 14:00 น. ของวันนี้
	returnedAt := bangkokToday(14, 0)
	insertLoan(t, pool, requester, "CS-NB-0001",
		bangkokToday(config.LoanNormalReturnHour, 0).AddDate(0, 0, -1), &returnedAt)

	// When ส่งคำขอยืมใหม่เวลา 00:00:00 น. ของวันถัดไปพอดี
	svc := NewServiceWithClock(NewRepository(pool), FixedClock{T: bangkokToday(0, 0).AddDate(0, 0, 1)})
	_, err := svc.Create(ctx, newTestRequest("CS-NB-0003"))

	// Then (1) ส่งได้ เพราะถึงจุดปลดล็อกแล้ว
	if err != nil {
		t.Fatalf("(1) เวลา 00:00:00 ของวันถัดไปต้องส่งคำขอได้ แต่ได้ error: %v", err)
	}
	if n := countRequests(t, pool, requester); n != 1 {
		t.Errorf("(1) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-12-4 (ทางผิด) ทำให้เงื่อนไข "คืนเกินกำหนด" ไม่เป็นจริง คือคืนตรงเวลา
func TestTC_BRQ_12_4_ReturnedOnTime(t *testing.T) {
	pool := testPool(t)
	const requester = "student-tc-12-4"
	ctx := studentCtx(requester)

	// Given กำหนดคืนวันนี้ 16:00 น. และคืนไปแล้วตั้งแต่ 14:00 น. ของวันนี้ คือคืนตรงเวลา
	returnedAt := bangkokToday(14, 0)
	insertLoan(t, pool, requester, "CS-NB-0001", bangkokToday(config.LoanNormalReturnHour, 0), &returnedAt)

	// When ส่งคำขอยืมใหม่เวลา 23:59 น. ของวันเดียวกัน
	svc := NewServiceWithClock(NewRepository(pool), FixedClock{T: bangkokToday(23, 59)})
	_, err := svc.Create(ctx, newTestRequest("CS-NB-0003"))

	// Then (1) ส่งได้ เพราะการปิดกั้นใช้กับผู้ที่คืนช้าเท่านั้น
	if err != nil {
		t.Fatalf("(1) คืนตรงเวลาต้องส่งคำขอได้ทันที แต่ได้ error: %v", err)
	}
	if n := countRequests(t, pool, requester); n != 1 {
		t.Errorf("(1) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
}
