package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-06-1 ถึง TC-BRQ-06-4 (AC-BRQ-06, FR-BRQ-06)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
//
// ชุดนี้ยิงผ่าน HTTP จริงและตรึงนาฬิกาด้วย FixedClock เพราะทุกแถวขึ้นกับ "ตอนนี้กี่โมง"
// ใช้ข้อมูลตั้งต้นของ seed.sql ที่ตรงกับ Given อยู่แล้ว
//   - student-001 ยืมของที่กำหนดคืนวันนี้ 16:00 น. และยังไม่คืน
//   - student-002 ยืมของที่กำหนดคืนเมื่อวาน 16:00 น. และยังไม่คืน (เลยกำหนดแล้ว)

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"salmonnorway/backend/internal/config"
)

// decodeJSON อ่าน body ของคำตอบลงโครงสร้างที่ส่งมา
func decodeJSON(t *testing.T, w *httptest.ResponseRecorder, into any) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), into); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v body=%s", err, w.Body.String())
	}
}

// overdueCount นับรายการยืมที่ยังไม่คืนของผู้ยืมคนหนึ่ง ใช้ยืนยันว่าข้อมูลตั้งต้นตรงกับ Given
func overdueCount(t *testing.T, borrowerID string) int {
	t.Helper()
	var n int
	if err := testPool(t).QueryRow(context.Background(),
		`SELECT count(*) FROM loan WHERE borrower_id = $1 AND returned_at IS NULL AND due_at < now()`,
		borrowerID).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// TC-BRQ-06-1 (ทางปกติ) ผู้ขอมีของที่กำหนดคืนเมื่อวาน 16:00 น. และยังไม่คืน
func TestTC_BRQ_06_1_OverdueYesterday(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-06-1"

	// Given student-002 มีอุปกรณ์ที่ยังไม่คืน กำหนดคืนเมื่อวาน 16:00 น. (seed.sql ชิ้นที่ 5)
	insertTempItem(t, assetCode)
	before := countRequests(t, pool, "student-002")

	// When ส่งคำขอยืมใหม่
	w := postRequest(t, srv, tokenStudentOverdue, bodyFor(assetCode))

	// Then (1) ได้ 409 HAS_OVERDUE_LOAN
	if w.Code != http.StatusConflict {
		t.Fatalf("(1) ผู้ขอที่มีของเกินกำหนดต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if code := errorCodeOf(t, w); code != "HAS_OVERDUE_LOAN" {
		t.Errorf("(1) รหัสต้องเป็น HAS_OVERDUE_LOAN ตาม CON-TECH-03 แต่ได้ %q", code)
	}
	// Then (2) body มีเหตุผล
	var body ErrorDTO
	decodeJSON(t, w, &body)
	if body.Message == "" {
		t.Error("(2) ต้องมีข้อความแจ้งเหตุผลว่ามีของค้างเกินกำหนด (FR-BRQ-06)")
	}
	// Then (3) ไม่มีแถวใหม่
	if after := countRequests(t, pool, "student-002"); after != before {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่จำนวนเปลี่ยนจาก %d เป็น %d", before, after)
	}
}

// TC-BRQ-06-2 (ขอบ) กำหนดคืนวันนี้ 16:00 น. และตอนนี้คือ 15:59 น. ยังไม่เลยกำหนด
func TestTC_BRQ_06_2_BeforeDeadline(t *testing.T) {
	pool := testPool(t)
	const assetCode = "TC-BRQ-06-2"

	// Given student-001 มีของยังไม่คืน กำหนดคืนวันนี้ 16:00 น. (seed.sql ชิ้นที่ 4)
	insertTempItem(t, assetCode)
	// ตรึงนาฬิกาไว้ 1 นาทีก่อนเวลาคืนของการยืมปกติ
	srv := newTestServerAt(t, FixedClock{T: bangkokToday(config.LoanNormalReturnHour-1, 59)})

	// When ส่งคำขอยืมใหม่
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 201 เพราะยังไม่เลยเวลา 16:00 น. ตาม FR-BRQ-06
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ก่อนถึงเวลาคืนต้องส่งคำขอได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(1) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-06-3 (ขอบ) ขยับไปอีก 2 นาที เลยเวลาคืนแล้ว ผลต้องเปลี่ยน
func TestTC_BRQ_06_3_AfterDeadline(t *testing.T) {
	pool := testPool(t)
	const assetCode = "TC-BRQ-06-3"

	// Given เงื่อนไขเดิม แต่ตอนนี้คือ 16:01 น. ของวันกำหนดคืน
	insertTempItem(t, assetCode)
	srv := newTestServerAt(t, FixedClock{T: bangkokToday(config.LoanNormalReturnHour, 1)})

	// When ส่งคำขอยืมใหม่
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 409 HAS_OVERDUE_LOAN
	if w.Code != http.StatusConflict {
		t.Fatalf("(1) เลยเวลาคืนแล้วต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if code := errorCodeOf(t, w); code != "HAS_OVERDUE_LOAN" {
		t.Errorf("(1) รหัสต้องเป็น HAS_OVERDUE_LOAN แต่ได้ %q", code)
	}
	// Then (2) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(2) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-06-4 (ทางผิด) ของที่เกินกำหนดเป็นของนักศึกษาคนอื่น ส่วนผู้ขอเองยังไม่เลยกำหนด
func TestTC_BRQ_06_4_OtherStudentOverdue(t *testing.T) {
	pool := testPool(t)
	const assetCode = "TC-BRQ-06-4"

	// Given student-002 มีของเกินกำหนดค้างอยู่ ส่วน student-001 ยังไม่ถึงกำหนดคืน (นาฬิกา 09:00 น.)
	insertTempItem(t, assetCode)
	srv := newTestServer(t)
	if n := overdueCount(t, "student-002"); n == 0 {
		t.Fatal("ข้อมูลตั้งต้นต้องมีของเกินกำหนดของ student-002 อยู่ จึงจะทดสอบแถวนี้ได้")
	}

	// When student-001 ส่งคำขอยืมใหม่
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 201 เพราะกฎนับเฉพาะของค้างของผู้ขอเอง
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ของเกินกำหนดของคนอื่นต้องไม่ปิดกั้น ต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(1) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
}
