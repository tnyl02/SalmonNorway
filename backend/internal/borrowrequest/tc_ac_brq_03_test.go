package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-03-1 ถึง TC-BRQ-03-4 (AC-BRQ-03, FR-BRQ-03)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)

import (
	"net/http"
	"slices"
	"testing"
	"time"
)

// ช่องที่ AC-BRQ-03 ระบุว่าต้องกรอกครบ 5 ช่อง (ทีมตัดสิน 2569-11-11)
var fiveRequiredFields = []string{"requesterName", "studentCode", "startAt", "endAt", "purpose"}

// TC-BRQ-03-1 (ทางปกติ) กรอกครบทุกช่อง ยกเว้นวันและเวลาที่คืน
func TestTC_BRQ_03_1_MissingEndAt(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-03-1"

	// Given กรอกครบทุกช่อง ยกเว้นไม่ระบุวันและเวลาที่คืน
	insertTempItem(t, assetCode)

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyWithout(assetCode, "endAt"))

	// Then (1) ได้ 422 VALIDATION_FAILED
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ไม่ระบุเวลาคืนต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	code, fields := validationFieldsOf(t, w)
	if code != "VALIDATION_FAILED" {
		t.Errorf("(1) รหัสต้องเป็น VALIDATION_FAILED ตาม CON-TECH-03 แต่ได้ %q", code)
	}
	// Then (2) รายชื่อช่องที่ผิดมี endAt
	if !slices.Contains(fields, "endAt") {
		t.Errorf("(2) ต้องชี้ช่อง endAt แต่ได้ %v", fields)
	}
	// Then (3) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-03-2 (ขอบ) ไม่กรอกเลยทั้ง 5 ช่องที่ AC ระบุ
func TestTC_BRQ_03_2_AllFiveMissing(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-03-2"

	// Given ส่งมาแต่รหัสครุภัณฑ์กับประเภทการยืม ไม่มีทั้ง 5 ช่องที่บังคับ
	insertTempItem(t, assetCode)
	body := `{"assetCode":"` + assetCode + `","loanType":"normal"}`

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, body)

	// Then (1) ได้ 422
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ไม่กรอกเลยต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	_, fields := validationFieldsOf(t, w)
	// Then (2) รายชื่อช่องที่ผิดมีครบทั้ง 5 ช่อง
	for _, field := range fiveRequiredFields {
		if !slices.Contains(fields, field) {
			t.Errorf("(2) ต้องชี้ช่อง %s ด้วย แต่ได้ %v", field, fields)
		}
	}
	// Then (3) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-03-3 (ขอบ) ฝั่งตรงข้ามของขอบ คือกรอกครบทั้ง 5 ช่องพอดี
func TestTC_BRQ_03_3_AllFiveFilled(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-03-3"

	// Given กรอกครบทั้ง 5 ช่อง
	insertTempItem(t, assetCode)

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) กรอกครบต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	// Then (2) มีแถวใหม่ 1 แถว
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(2) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-03-4 (ทางผิด) มีช่องชื่อ-นามสกุลแต่ไม่มีค่า
func TestTC_BRQ_03_4_EmptyRequesterName(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-03-4"

	// Given ส่งชื่อ-นามสกุลเป็นสตริงว่าง ช่องอื่นกรอกครบ
	insertTempItem(t, assetCode)
	body := `{"assetCode":"` + assetCode + `","requesterName":"","studentCode":"660710715",` +
		`"loanType":"normal","purpose":"study",` +
		`"startAt":"` + bangkokToday(9, 0).Format(time.RFC3339) + `",` +
		`"endAt":"` + bangkokToday(15, 0).Format(time.RFC3339) + `"}`

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, body)

	// Then (1) ได้ 422
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ชื่อว่างต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	_, fields := validationFieldsOf(t, w)
	// Then (2) รายชื่อช่องที่ผิดมี requesterName
	if !slices.Contains(fields, "requesterName") {
		t.Errorf("(2) ต้องชี้ช่อง requesterName แต่ได้ %v", fields)
	}
	// Then (3) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}
