package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-13-1 ถึง TC-BRQ-13-4 (AC-BRQ-13, FR-BRQ-03 ข้อ ค)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
//
// ตัวเลขขอบของชุดนี้คือ config.LoanNormalReturnHour (16:00 น.) ไม่เขียนเลข 16 ตรง ๆ ใน test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"slices"
	"testing"
	"time"

	"salmonnorway/backend/internal/config"
)

// multipartWithEnd สร้าง body แบบ multipart ที่กำหนดเวลาคืนเองได้
// ต่างจาก multipartForm ตรงที่อันนั้นคิดเวลาคืนให้เอง แต่ TC-BRQ-13-4 ต้องระบุ 18:00 น. ของวันนี้
func multipartWithEnd(t *testing.T, assetCode, loanType string, end time.Time, file *testFile) (string, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for field, value := range map[string]string{
		"assetCode":     assetCode,
		"requesterName": "ธัญลักษณ์ ทดสอบ",
		"studentCode":   "660710715",
		"loanType":      loanType,
		"purpose":       "study",
		"startAt":       bangkokToday(9, 0).Format(time.RFC3339),
		"endAt":         end.Format(time.RFC3339),
	} {
		if err := w.WriteField(field, value); err != nil {
			t.Fatal(err)
		}
	}
	if file != nil {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="attachment"; filename=%q`, file.name))
		h.Set("Content-Type", file.contentType)
		part, err := w.CreatePart(h)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write(file.content); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return w.FormDataContentType(), &buf
}

// TC-BRQ-13-1 (ทางปกติ) ยืมปกติ เริ่ม 09:00 คืน 18:00 ของวันเดียวกัน
func TestTC_BRQ_13_1_NormalReturn1800(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-13-1"

	// Given ยืมปกติ เริ่มยืมวันนี้ 09:00 น.
	insertTempItem(t, assetCode)

	// When ระบุเวลาคืนวันนี้ 18:00 น. แล้วส่ง POST
	body := bodyWithWindow(assetCode, config.LoanTypeNormal, bangkokToday(9, 0), bangkokToday(18, 0))
	w := postRequest(t, srv, tokenStudentDept1, body)

	// Then (1) ได้ 422 VALIDATION_FAILED
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ยืมปกติคืน 18:00 ต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	code, fields := validationFieldsOf(t, w)
	if code != "VALIDATION_FAILED" {
		t.Errorf("(1) รหัสต้องเป็น VALIDATION_FAILED แต่ได้ %q", code)
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

// TC-BRQ-13-2 (ขอบ) คืนเวลา 16:00 น. พอดี ซึ่งเป็นตัวเลขขอบจาก config.LoanNormalReturnHour
func TestTC_BRQ_13_2_NormalReturnExactly1600(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-13-2"

	// Given ยืมปกติ เริ่มยืมวันนี้ 09:00 น.
	insertTempItem(t, assetCode)

	// When ระบุเวลาคืนวันนี้ 16:00 น. พอดี แล้วส่ง POST
	body := bodyWithWindow(assetCode, config.LoanTypeNormal,
		bangkokToday(9, 0), bangkokToday(config.LoanNormalReturnHour, 0))
	w := postRequest(t, srv, tokenStudentDept1, body)

	// Then (1) ได้ 201 เพราะ "ภายใน 16:00 น." รวมเวลา 16:00 น. พอดี
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) คืนเวลา %d:00 พอดีต้องได้ 201 แต่ได้ %d body=%s",
			config.LoanNormalReturnHour, w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	// Then (2) มีแถวใหม่ 1 แถว
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(2) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-13-3 (ขอบ) ขยับไปอีก 1 นาที ผลต้องเปลี่ยน
func TestTC_BRQ_13_3_NormalReturn1601(t *testing.T) {
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-13-3"

	// Given ยืมปกติ เริ่มยืมวันนี้ 09:00 น.
	insertTempItem(t, assetCode)

	// When ระบุเวลาคืนวันนี้ 16:01 น. แล้วส่ง POST
	body := bodyWithWindow(assetCode, config.LoanTypeNormal,
		bangkokToday(9, 0), bangkokToday(config.LoanNormalReturnHour, 1))
	w := postRequest(t, srv, tokenStudentDept1, body)

	// Then (1) ได้ 422
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) คืนเวลา %d:01 ต้องได้ 422 แต่ได้ %d body=%s",
			config.LoanNormalReturnHour, w.Code, w.Body.String())
	}
	// Then (2) รายชื่อช่องที่ผิดมี endAt
	if _, fields := validationFieldsOf(t, w); !slices.Contains(fields, "endAt") {
		t.Errorf("(2) ต้องชี้ช่อง endAt แต่ได้ %v", fields)
	}
}

// TC-BRQ-13-4 (ทางผิด) ทำให้เงื่อนไข "ปกติ" ไม่เป็นจริง คือเลือกค้างคืนและแนบไฟล์
func TestTC_BRQ_13_4_OvernightNotLimitedBy1600(t *testing.T) {
	t.Setenv("ATTACHMENT_DIR", t.TempDir()) // ต้องตั้งก่อนสร้าง service เพราะอ่าน config ตอนสร้าง
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-13-4"

	// Given เลือกประเภทการยืมเป็นค้างคืน และแนบ PDF ตามที่ FR-BRQ-10 ข้อ ก บังคับ
	insertTempItem(t, assetCode)
	ct, body := multipartWithEnd(t, assetCode, config.LoanTypeOvernight, bangkokToday(18, 0),
		&testFile{name: "memo.pdf", contentType: "application/pdf", content: []byte("%PDF-1.4 ใบบันทึกข้อความ")})

	// When ระบุเวลาคืนวันนี้ 18:00 น. แล้วส่ง POST
	w := postMultipart(t, srv, tokenStudentDept1, ct, body)

	// Then (1) ได้ 201 เพราะกฎ 16:00 น. ใช้กับการยืมปกติเท่านั้น (DOM-LOAN-01)
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ยืมค้างคืนต้องไม่ติดกฎ %d:00 ต้องได้ 201 แต่ได้ %d body=%s",
			config.LoanNormalReturnHour, w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(1) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
}
