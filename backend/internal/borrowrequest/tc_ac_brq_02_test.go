package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-02-1 ถึง TC-BRQ-02-4 (AC-BRQ-02, FR-BRQ-02)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"salmonnorway/backend/internal/config"
)

// validationFieldsOf อ่านรายชื่อช่องที่ผิดจากคำตอบ 422 (CON-TECH-03 แนบ fields มาด้วย)
func validationFieldsOf(t *testing.T, w *httptest.ResponseRecorder) (string, []string) {
	t.Helper()
	var got struct {
		Code   string   `json:"code"`
		Fields []string `json:"fields"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v body=%s", err, w.Body.String())
	}
	return got.Code, got.Fields
}

// purposeOfRequest อ่านวัตถุประสงค์และข้อความที่พิมพ์เพิ่มของคำขอใบหนึ่ง
func purposeOfRequest(t *testing.T, id int) (string, string) {
	t.Helper()
	var purpose string
	var other *string
	if err := testPool(t).QueryRow(context.Background(),
		`SELECT purpose, purpose_other FROM borrow_request WHERE id = $1`, id).Scan(&purpose, &other); err != nil {
		t.Fatal(err)
	}
	if other == nil {
		return purpose, ""
	}
	return purpose, *other
}

// TC-BRQ-02-1 (ทางปกติ) เลือก "เพื่อการเรียน" แล้วส่งสำเร็จ
func TestTC_BRQ_02_1_PurposeStudySaved(t *testing.T) {
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-02-1"

	// Given ผู้ขอกรอกครบ เลือกวัตถุประสงค์ study
	insertTempItem(t, assetCode)

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose(assetCode, "study", ""))

	// Then (1) ได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	got := decodeCreated(t, w)
	// Then (2) แถวที่บันทึกมี purpose = 'study'
	if purpose, _ := purposeOfRequest(t, got.RequestID); purpose != "study" {
		t.Errorf("(2) วัตถุประสงค์ที่บันทึกต้องเป็น study แต่ได้ %q", purpose)
	}
}

// TC-BRQ-02-2 (ขอบ) ค่าที่เหลือในรายการ 3 ค่าของ config.PurposeValues
func TestTC_BRQ_02_2_AllPurposeValues(t *testing.T) {
	srv := newTestServer(t)

	for _, tc := range []struct {
		purpose   string
		other     string
		assetCode string
	}{
		{"activity", "", "TC-BRQ-02-2a"},
		{config.PurposeOther, "ซ้อมดนตรี", "TC-BRQ-02-2b"}, // other ต้องพิมพ์ระบุด้วย (FR-BRQ-02 ข้อ ก)
	} {
		t.Run(tc.purpose, func(t *testing.T) {
			// Given อุปกรณ์ว่างคนละชิ้นกัน เพื่อไม่ให้ติดกฎคำขอซ้ำ (FR-BRQ-05 ข)
			insertTempItem(t, tc.assetCode)

			// When ส่ง POST ด้วยค่าวัตถุประสงค์นี้
			w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose(tc.assetCode, tc.purpose, tc.other))

			// Then (1) ได้ 201
			if w.Code != http.StatusCreated {
				t.Fatalf("(1) ค่า %q ต้องส่งได้ 201 แต่ได้ %d body=%s", tc.purpose, w.Code, w.Body.String())
			}
			got := decodeCreated(t, w)
			// Then (2) purpose ในแถวตรงกับที่ส่ง
			if purpose, _ := purposeOfRequest(t, got.RequestID); purpose != tc.purpose {
				t.Errorf("(2) วัตถุประสงค์ที่บันทึกต้องเป็น %q แต่ได้ %q", tc.purpose, purpose)
			}
		})
	}
}

// TC-BRQ-02-3 (ทางผิด) ไม่เลือกวัตถุประสงค์
func TestTC_BRQ_02_3_PurposeMissing(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-02-3"

	// Given ผู้ขอกรอกช่องอื่นครบ แต่ไม่เลือกวัตถุประสงค์
	insertTempItem(t, assetCode)

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose(assetCode, "", ""))

	// Then (1) ได้ 422 VALIDATION_FAILED
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ไม่เลือกวัตถุประสงค์ต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	code, fields := validationFieldsOf(t, w)
	if code != "VALIDATION_FAILED" {
		t.Errorf("(1) รหัสต้องเป็น VALIDATION_FAILED ตาม CON-TECH-03 แต่ได้ %q", code)
	}
	// Then (2) รายชื่อช่องที่ผิดมี purpose
	if !slices.Contains(fields, "purpose") {
		t.Errorf("(2) ต้องชี้ช่อง purpose แต่ได้ %v", fields)
	}
	// Then (3) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-02-4 (ทางผิด) ค่าที่ไม่อยู่ในรายการ 3 ค่า
// spec ไม่ได้เขียนกรณีนี้ ทีมตัดสิน 2569-11-13 ให้ถือว่าไม่ผ่านการตรวจ เหมือนไม่ได้เลือก
func TestTC_BRQ_02_4_PurposeNotInList(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-02-4"

	// Given ส่งค่าวัตถุประสงค์ที่ไม่อยู่ในรายการ 3 ค่าของ config.PurposeValues
	insertTempItem(t, assetCode)

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose(assetCode, "travel", ""))

	// Then (1) ได้ 422 VALIDATION_FAILED
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ค่านอกรายการต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	code, fields := validationFieldsOf(t, w)
	if code != "VALIDATION_FAILED" {
		t.Errorf("(1) รหัสต้องเป็น VALIDATION_FAILED แต่ได้ %q", code)
	}
	// Then (2) รายชื่อช่องที่ผิดมี purpose
	if !slices.Contains(fields, "purpose") {
		t.Errorf("(2) ต้องชี้ช่อง purpose แต่ได้ %v", fields)
	}
	// Then (3) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}
