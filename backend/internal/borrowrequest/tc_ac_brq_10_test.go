package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-10-1 ถึง TC-BRQ-10-5 (AC-BRQ-10, FR-BRQ-10)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
//
// ตัวเลขขอบคือ config.AttachmentMaxBytes และชนิดไฟล์คือ config.AttachmentAllowedTypes
// ไม่เขียนเลข 10MB หรือชื่อชนิดไฟล์ตรง ๆ ใน test (AGENTS.md กฎข้อ 9)

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

// multipartWithoutPurpose สร้าง body ของการยืมค้างคืนที่แนบไฟล์แล้ว แต่ไม่ได้เลือกวัตถุประสงค์
// multipartForm เดิมใส่ purpose ให้เสมอ จึงใช้ทดสอบแถวนี้ไม่ได้
func multipartWithoutPurpose(t *testing.T, assetCode string, file *testFile) (string, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for field, value := range map[string]string{
		"assetCode":     assetCode,
		"requesterName": "ธัญลักษณ์ ทดสอบ",
		"studentCode":   "660710715",
		"loanType":      config.LoanTypeOvernight,
		"startAt":       bangkokToday(9, 0).Format(time.RFC3339),
		"endAt":         bangkokToday(15, 0).AddDate(0, 0, 2).Format(time.RFC3339),
	} {
		if err := w.WriteField(field, value); err != nil {
			t.Fatal(err)
		}
	}
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
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return w.FormDataContentType(), &buf
}

// pdfOfSize สร้างไฟล์ PDF จำลองขนาดตามที่ระบุ
func pdfOfSize(n int) *testFile {
	content := make([]byte, n)
	copy(content, []byte("%PDF-1.4"))
	return &testFile{name: "memo.pdf", contentType: "application/pdf", content: content}
}

// TC-BRQ-10-1 (ทางปกติ) เลือกค้างคืนแต่ไม่แนบไฟล์
func TestTC_BRQ_10_1_OvernightNoAttachment(t *testing.T) {
	t.Setenv("ATTACHMENT_DIR", t.TempDir())
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-10-1"

	// Given เลือกประเภทการยืมเป็นค้างคืน และไม่แนบไฟล์ใบบันทึกข้อความ
	insertTempItem(t, assetCode)
	ct, body := multipartForm(t, assetCode, config.LoanTypeOvernight, nil)

	// When ส่ง POST
	w := postMultipart(t, srv, tokenStudentDept1, ct, body)

	// Then (1) ได้ 422 VALIDATION_FAILED
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ค้างคืนที่ไม่แนบไฟล์ต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	code, fields := validationFieldsOf(t, w)
	if code != "VALIDATION_FAILED" {
		t.Errorf("(1) รหัสต้องเป็น VALIDATION_FAILED ตาม CON-TECH-03 แต่ได้ %q", code)
	}
	// Then (2) รายชื่อช่องที่ผิดมี attachment
	if !slices.Contains(fields, "attachment") {
		t.Errorf("(2) ต้องชี้ช่อง attachment แต่ได้ %v", fields)
	}
	// Then (3) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-10-2 (ทางปกติ) เลือกค้างคืน แนบไฟล์แล้ว แต่ไม่เลือกวัตถุประสงค์
func TestTC_BRQ_10_2_OvernightNoPurpose(t *testing.T) {
	t.Setenv("ATTACHMENT_DIR", t.TempDir())
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-10-2"

	// Given เลือกค้างคืน แนบ PDF แล้ว แต่ไม่เลือกวัตถุประสงค์
	insertTempItem(t, assetCode)
	ct, body := multipartWithoutPurpose(t, assetCode, pdfOfSize(1024))

	// When ส่ง POST
	w := postMultipart(t, srv, tokenStudentDept1, ct, body)

	// Then (1) ได้ 422
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ไม่เลือกวัตถุประสงค์ต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	_, fields := validationFieldsOf(t, w)
	// Then (2) รายชื่อช่องที่ผิดมี purpose
	if !slices.Contains(fields, "purpose") {
		t.Errorf("(2) ต้องชี้ช่อง purpose แต่ได้ %v", fields)
	}
	// Then (3) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-10-3 (ขอบ) ไฟล์ขนาดเท่าเพดานพอดี
func TestTC_BRQ_10_3_AttachmentExactly10MB(t *testing.T) {
	t.Setenv("ATTACHMENT_DIR", t.TempDir())
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-10-3"

	// Given แนบไฟล์ PDF ขนาดเท่า config.AttachmentMaxBytes พอดี
	insertTempItem(t, assetCode)
	ct, body := multipartForm(t, assetCode, config.LoanTypeOvernight, pdfOfSize(config.AttachmentMaxBytes))

	// When ส่ง POST
	w := postMultipart(t, srv, tokenStudentDept1, ct, body)

	// Then (1) ได้ 201 เพราะ "ไม่เกิน" รวมขนาดเท่าเพดานพอดี
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ไฟล์ขนาดเท่าเพดานพอดีต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	// Then (2) มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(2) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-10-4 (ขอบ) ไฟล์ใหญ่เกินเพดาน 1 ไบต์
func TestTC_BRQ_10_4_AttachmentOver10MB(t *testing.T) {
	t.Setenv("ATTACHMENT_DIR", t.TempDir())
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-10-4"

	// Given แนบไฟล์ PDF ที่ใหญ่กว่าเพดาน 1 ไบต์
	insertTempItem(t, assetCode)
	ct, body := multipartForm(t, assetCode, config.LoanTypeOvernight, pdfOfSize(config.AttachmentMaxBytes+1))

	// When ส่ง POST
	w := postMultipart(t, srv, tokenStudentDept1, ct, body)

	// Then (1) ได้ 422
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ไฟล์ใหญ่เกินเพดานต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	_, fields := validationFieldsOf(t, w)
	if !slices.Contains(fields, "attachment") {
		t.Errorf("(1) ต้องชี้ช่อง attachment แต่ได้ %v", fields)
	}
	// Then (2) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(2) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
	// การบอกขนาดที่รับได้เป็นหน้าที่ของหน้าจอ ตรวจที่ TC-BRQ-10-7 (ทีมตัดสิน 2569-11-13 แก้แถวในตาราง)
}

// TC-BRQ-10-5 (ทางผิด) ไฟล์ชนิดที่ไม่อยู่ใน config.AttachmentAllowedTypes
func TestTC_BRQ_10_5_AttachmentWrongType(t *testing.T) {
	t.Setenv("ATTACHMENT_DIR", t.TempDir())
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-10-5"

	// Given แนบไฟล์ชนิดที่ไม่อยู่ในรายการที่รับได้
	insertTempItem(t, assetCode)
	exe := &testFile{name: "setup.exe", contentType: "application/octet-stream", content: []byte("MZ")}
	ct, body := multipartForm(t, assetCode, config.LoanTypeOvernight, exe)

	// When ส่ง POST
	w := postMultipart(t, srv, tokenStudentDept1, ct, body)

	// Then (1) ได้ 422
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("(1) ไฟล์ชนิดที่ไม่รับต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	_, fields := validationFieldsOf(t, w)
	if !slices.Contains(fields, "attachment") {
		t.Errorf("(1) ต้องชี้ช่อง attachment แต่ได้ %v", fields)
	}
	// Then (2) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(2) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
	// การบอกชนิดไฟล์ที่รับได้เป็นหน้าที่ของหน้าจอ ตรวจที่ TC-BRQ-10-7 (ทีมตัดสิน 2569-11-13 แก้แถวในตาราง)
}
