package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-14-1 ถึง TC-BRQ-14-4 (AC-BRQ-14, FR-BRQ-02 ข้อ ก)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)

import (
	"net/http"
	"slices"
	"testing"

	"salmonnorway/backend/internal/config"
)

// TC-BRQ-14-1 (ทางปกติ) เลือก "อื่น ๆ" และพิมพ์ข้อความระบุ
func TestTC_BRQ_14_1_OtherWithText(t *testing.T) {
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-14-1"
	const text = "ซ้อมดนตรี"

	// Given เลือกวัตถุประสงค์ other และพิมพ์ข้อความระบุ
	insertTempItem(t, assetCode)

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose(assetCode, config.PurposeOther, text))

	// Then (1) ได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) เลือกอื่น ๆ พร้อมข้อความต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	got := decodeCreated(t, w)
	purpose, other := purposeOfRequest(t, got.RequestID)
	// Then (2) แถวมี purpose = 'other'
	if purpose != config.PurposeOther {
		t.Errorf("(2) วัตถุประสงค์ต้องเป็น %q แต่ได้ %q", config.PurposeOther, purpose)
	}
	// Then (3) แถวมีข้อความที่พิมพ์
	if other != text {
		t.Errorf("(3) ข้อความระบุต้องเป็น %q แต่ได้ %q", text, other)
	}
}

// TC-BRQ-14-2 (ทางผิด) เลือก "อื่น ๆ" แต่ไม่พิมพ์ข้อความระบุ
func TestTC_BRQ_14_2_OtherWithoutText(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-14-2"

	// Given เลือก other แต่ไม่พิมพ์ข้อความระบุ
	insertTempItem(t, assetCode)

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose(assetCode, config.PurposeOther, ""))

	// Then (1) ได้ 422 VALIDATION_FAILED
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) เลือกอื่น ๆ โดยไม่พิมพ์ระบุต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	code, fields := validationFieldsOf(t, w)
	if code != "VALIDATION_FAILED" {
		t.Errorf("(1) รหัสต้องเป็น VALIDATION_FAILED แต่ได้ %q", code)
	}
	// Then (2) รายชื่อช่องที่ผิดมี purposeOther
	if !slices.Contains(fields, "purposeOther") {
		t.Errorf("(2) ต้องชี้ช่อง purposeOther แต่ได้ %v", fields)
	}
	// Then (3) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-14-3 (ขอบ) เลือก "อื่น ๆ" แล้วพิมพ์เฉพาะช่องว่าง
// spec ไม่ได้เขียนกรณีนี้ ทีมตัดสิน 2569-11-13 ว่าช่องว่างล้วนยังไม่นับว่าระบุแล้ว
func TestTC_BRQ_14_3_OtherBlankText(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-14-3"

	// Given เลือก other และพิมพ์เฉพาะช่องว่าง
	insertTempItem(t, assetCode)

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose(assetCode, config.PurposeOther, "   "))

	// Then (1) ได้ 422 VALIDATION_FAILED
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ช่องว่างล้วนต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	code, fields := validationFieldsOf(t, w)
	if code != "VALIDATION_FAILED" {
		t.Errorf("(1) รหัสต้องเป็น VALIDATION_FAILED แต่ได้ %q", code)
	}
	// Then (2) รายชื่อช่องที่ผิดมี purposeOther
	if !slices.Contains(fields, "purposeOther") {
		t.Errorf("(2) ต้องชี้ช่อง purposeOther แต่ได้ %v", fields)
	}
	// Then (3) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-14-4 (ทางผิด) ทำให้เงื่อนไข "เลือกอื่น ๆ" ไม่เป็นจริง แต่ยังส่งข้อความระบุมาด้วย
func TestTC_BRQ_14_4_OtherTextIgnored(t *testing.T) {
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-14-4"

	// Given เลือก study แต่ส่งข้อความระบุมาด้วย
	insertTempItem(t, assetCode)

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose(assetCode, "study", "ข้อความที่ไม่ควรถูกเก็บ"))

	// Then (1) ได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	got := decodeCreated(t, w)
	purpose, other := purposeOfRequest(t, got.RequestID)
	// Then (2) แถวมี purpose = 'study'
	if purpose != "study" {
		t.Errorf("(2) วัตถุประสงค์ต้องเป็น study แต่ได้ %q", purpose)
	}
	// Then (3) ข้อความระบุต้องไม่ถูกเก็บ (FR-BRQ-02 ข้อ ข)
	if other != "" {
		t.Errorf("(3) ไม่ได้เลือกอื่น ๆ จึงต้องไม่เก็บข้อความระบุ แต่เก็บ %q", other)
	}
}
