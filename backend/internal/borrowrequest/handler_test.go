package borrowrequest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// ผู้ใช้จำลองจาก auth.NewMockProvider (รอ Q-SH-03 สำหรับโปรโตคอลจริง)
const tokenStudentDept1 = "mock-student-001"

// newTestServer ประกอบ handler เข้ากับ auth.Middleware แบบเดียวกับ main.go (IF-IDP-01)
func newTestServer(t *testing.T) http.Handler {
	return newTestServerAt(t, testClock())
}

// newTestServerAt ประกอบ handler ด้วยนาฬิกาที่กำหนดเอง ใช้กับ test ที่ต้องตรึงเวลา (AC-BRQ-06)
func newTestServerAt(t *testing.T, clock Clock) http.Handler {
	t.Helper()
	svc := NewServiceWithClock(NewRepository(testPool(t)), clock)
	return auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes())
}

// postRequest ยิง POST /v1/loan-requests ด้วย token ที่ให้ ถ้า token ว่างคือไม่ส่ง header เลย
func postRequest(t *testing.T, srv http.Handler, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, PathPrefix, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)
	return w
}

// bodyFor สร้าง JSON ของฟอร์มที่กรอกครบ สำหรับอุปกรณ์รหัสที่ระบุ
func bodyFor(assetCode string) string {
	start := time.Now().Format(time.RFC3339)
	end := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
	return `{"assetCode":"` + assetCode + `","requesterName":"ธัญลักษณ์ ทดสอบ","studentCode":"660710715",` +
		`"loanType":"normal","purpose":"study","startAt":"` + start + `","endAt":"` + end + `"}`
}

// AC-BRQ-04 (FR-BRQ-04) ยิงผ่าน HTTP จริง
// Given นักศึกษาที่ยืนยันตัวตนแล้ว และอุปกรณ์ชิ้นนั้นสถานะ "ว่าง"
// When  POST /v1/loan-requests ด้วยข้อมูลครบถ้วน
// Then  ได้ 201 พร้อมรหัสคำขอ สถานะ requested และเวลาที่ส่ง โดยสถานะอุปกรณ์ยังเป็น "ว่าง"
func TestAC_BRQ_04_HTTP(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "CS-NB-0002" // อุปกรณ์ว่างของภาควิชา 1 ในข้อมูลตั้งต้น (T-03)

	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))
	if w.Code != http.StatusCreated {
		t.Fatalf("ส่งคำขอสำเร็จต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}

	var got CreatedDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM borrow_request WHERE id = $1`, got.RequestID)
	})

	if got.RequestID == 0 {
		t.Error("คำตอบต้องมีรหัสคำขอ")
	}
	if got.Status != config.StatusRequested {
		t.Errorf("สถานะต้องเป็น %q แต่ได้ %q", config.StatusRequested, got.Status)
	}
	if d := time.Since(got.SubmittedAt); d < -time.Minute || d > time.Minute {
		t.Errorf("เวลาที่ส่งต้องใกล้เคียงเวลาปัจจุบัน แต่ต่างกัน %v", d)
	}

	// สถานะอุปกรณ์ต้องไม่เปลี่ยน (DOM-REQ-01)
	var status string
	if err := pool.QueryRow(context.Background(),
		`SELECT status FROM equipment_item WHERE asset_code = $1`, assetCode).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "available" {
		t.Errorf("อุปกรณ์ที่มีคำขอค้างต้องยังเป็น available แต่เป็น %q", status)
	}
}

// T-05 เสร็จเมื่อ: เรียกโดยไม่มี token ต้องได้ 401 UNAUTHENTICATED ตามรูปแบบกลาง (IF-IDP-01, CON-TECH-03)
func TestCreate_WithoutToken_401(t *testing.T) {
	srv := newTestServer(t)

	w := postRequest(t, srv, "", bodyFor("CS-NB-0003"))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("ไม่มี token ต้องได้ 401 แต่ได้ %d", w.Code)
	}
	var e ErrorDTO
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
	}
	if e.Code != "UNAUTHENTICATED" {
		t.Errorf("รหัสต้องเป็น UNAUTHENTICATED ตาม CON-TECH-03 แต่ได้ %q", e.Code)
	}

	// ต้องไม่มีคำขอเกิดขึ้นจาก request ที่ถูกปฏิเสธ
	var count int
	if err := testPool(t).QueryRow(context.Background(), `SELECT count(*) FROM borrow_request`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("request ที่ไม่ผ่านการยืนยันตัวตนต้องไม่สร้างคำขอ แต่พบ %d ใบ", count)
	}
}

// T-05 เสร็จเมื่อ: body ที่อ่านไม่ออกต้องได้ VALIDATION_FAILED ตามรูปแบบกลาง (CON-TECH-03)
func TestCreate_BadJSON_422(t *testing.T) {
	srv := newTestServer(t)

	w := postRequest(t, srv, tokenStudentDept1, `{"assetCode":`)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("body ที่อ่านไม่ออกต้องได้ 422 แต่ได้ %d", w.Code)
	}
	var e ErrorDTO
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
	}
	if e.Code != "VALIDATION_FAILED" || e.Message == "" {
		t.Errorf("ต้องได้ {code: VALIDATION_FAILED, message: ...} แต่ได้ %+v", e)
	}
}

// T-05 เสร็จเมื่อ: ฟีเจอร์นี้มีแต่การสร้างคำขอ method อื่นบน path เดียวกันต้องไม่ผ่าน
func TestRoutes_OnlyPost(t *testing.T) {
	srv := newTestServer(t)

	r := httptest.NewRequest(http.MethodGet, PathPrefix, nil)
	r.Header.Set("Authorization", "Bearer "+tokenStudentDept1)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET บน path นี้ต้องได้ 405 แต่ได้ %d", w.Code)
	}
}

// bodyWithout สร้าง JSON ของฟอร์มที่ครบทุกช่อง ยกเว้นช่องที่ระบุให้ตัดออก (ใช้กับ AC-BRQ-03)
func bodyWithout(assetCode, omit string) string {
	fields := map[string]string{
		"requesterName": `"ธัญลักษณ์ ทดสอบ"`,
		"studentCode":   `"660710715"`,
		"startAt":       `"` + time.Now().Format(time.RFC3339) + `"`,
		"endAt":         `"` + time.Now().Add(2*time.Hour).Format(time.RFC3339) + `"`,
		"purpose":       `"study"`,
	}
	delete(fields, omit)
	parts := []string{`"assetCode":"` + assetCode + `"`, `"loanType":"normal"`}
	for k, v := range fields {
		parts = append(parts, `"`+k+`":`+v)
	}
	return "{" + strings.Join(parts, ",") + "}"
}

// AC-BRQ-03 (FR-BRQ-03)
// Given ข้อมูลคำขอไม่ครบถ้วน (เช่น ไม่ได้ระบุวันและเวลาที่คืน)
// When  นักศึกษากดส่งคำขอ
// Then  ระบบต้องไม่สร้างคำขอ และระบุช่องข้อมูลที่ยังต้องกรอก
//
//	(ช่องที่ตรวจ 5 ช่อง: ชื่อ-นามสกุล, รหัสนักศึกษา, วันและเวลาที่เริ่มยืม, วันและเวลาที่คืน, วัตถุประสงค์การยืม)
func TestAC_BRQ_03(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)

	for _, field := range []string{"requesterName", "studentCode", "startAt", "endAt", "purpose"} {
		w := postRequest(t, srv, tokenStudentDept1, bodyWithout("CS-NB-0001", field))
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("ขาดช่อง %s ต้องได้ 422 แต่ได้ %d body=%s", field, w.Code, w.Body.String())
			continue
		}

		var e struct {
			Code    string   `json:"code"`
			Message string   `json:"message"`
			Fields  []string `json:"fields"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
			t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
		}
		if e.Code != "VALIDATION_FAILED" || e.Message == "" {
			t.Errorf("ขาดช่อง %s ต้องได้ VALIDATION_FAILED พร้อมข้อความ แต่ได้ %+v", field, e)
		}
		if len(e.Fields) != 1 || e.Fields[0] != field {
			t.Errorf("ต้องระบุว่าขาดช่อง %s เท่านั้น แต่ได้ %v", field, e.Fields)
		}
	}

	// Then: ต้องไม่มีคำขอเกิดขึ้นเลยจากคำขอที่ข้อมูลไม่ครบ
	var count int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM borrow_request`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("ข้อมูลไม่ครบต้องไม่สร้างคำขอ แต่พบ %d ใบ", count)
	}
}

// T-06 เสร็จเมื่อ: ขาดหลายช่องพร้อมกัน ต้องบอกครบทุกช่องในครั้งเดียว (FR-BRQ-03)
func TestValidation_ReportsAllMissingFields(t *testing.T) {
	srv := newTestServer(t)

	w := postRequest(t, srv, tokenStudentDept1, `{"assetCode":"CS-NB-0001","loanType":"normal"}`)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("ต้องได้ 422 แต่ได้ %d", w.Code)
	}
	var e struct {
		Fields []string `json:"fields"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if len(e.Fields) != 5 {
		t.Errorf("ไม่กรอกอะไรเลยต้องบอกครบทั้ง 5 ช่อง แต่ได้ %v", e.Fields)
	}
}

// bodyWithPurpose สร้าง JSON ของฟอร์มที่ครบทุกช่อง โดยกำหนดค่าวัตถุประสงค์และข้อความระบุเพิ่มเอง
func bodyWithPurpose(assetCode, purpose, purposeOther string) string {
	start := time.Now().Format(time.RFC3339)
	end := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
	return `{"assetCode":"` + assetCode + `","requesterName":"ธัญลักษณ์ ทดสอบ","studentCode":"660710715",` +
		`"loanType":"normal","purpose":"` + purpose + `","purposeOther":"` + purposeOther + `",` +
		`"startAt":"` + start + `","endAt":"` + end + `"}`
}

// decodeCreated อ่านคำตอบตอนสร้างสำเร็จ แล้วสั่งลบแถวทิ้งหลัง test
func decodeCreated(t *testing.T, w *httptest.ResponseRecorder) CreatedDTO {
	t.Helper()
	var got CreatedDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
	}
	t.Cleanup(func() {
		_, _ = testPool(t).Exec(context.Background(), `DELETE FROM borrow_request WHERE id = $1`, got.RequestID)
	})
	return got
}

// AC-BRQ-02 (FR-BRQ-02)
// Given นักศึกษาเปิดแบบฟอร์มคำขอของอุปกรณ์ที่ว่าง
// When  เลือกวัตถุประสงค์ "เพื่อการเรียน" กรอกข้อมูลครบ แล้วส่งคำขอสำเร็จ
// Then  คำขอที่บันทึกต้องมีวัตถุประสงค์ตรงกับค่าที่เลือก
// และ   เมื่อยังไม่ได้เลือกวัตถุประสงค์ ต้องไม่สร้างคำขอ และระบุว่าต้องเลือก
func TestAC_BRQ_02(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)

	// When: เลือก study (เพื่อการเรียน) แล้วส่ง
	w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose("CS-NB-0001", "study", ""))
	if w.Code != http.StatusCreated {
		t.Fatalf("เลือกวัตถุประสงค์ที่ถูกต้องต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	got := decodeCreated(t, w)

	// Then: ค่าที่บันทึกต้องตรงกับที่เลือก และช่องพิมพ์ระบุต้องว่าง (FR-BRQ-02 ข้อ ข)
	var purpose string
	var purposeOther *string
	if err := pool.QueryRow(context.Background(),
		`SELECT purpose, purpose_other FROM borrow_request WHERE id = $1`, got.RequestID,
	).Scan(&purpose, &purposeOther); err != nil {
		t.Fatal(err)
	}
	if purpose != "study" {
		t.Errorf("วัตถุประสงค์ที่บันทึกต้องเป็น study แต่ได้ %q", purpose)
	}
	if purposeOther != nil {
		t.Errorf("ไม่ได้เลือก \"อื่น ๆ\" ช่องพิมพ์ระบุต้องว่าง แต่ได้ %q", *purposeOther)
	}

	// และ: ไม่เลือกวัตถุประสงค์ หรือเลือกค่าที่ไม่มีในรายการ ต้องถูกปฏิเสธทั้งคู่
	for _, bad := range []string{"", "exam", "เพื่อการเรียน"} {
		w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose("CS-NB-0002", bad, ""))
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("วัตถุประสงค์ %q ต้องได้ 422 แต่ได้ %d", bad, w.Code)
			continue
		}
		var e struct {
			Fields []string `json:"fields"`
		}
		_ = json.Unmarshal(w.Body.Bytes(), &e)
		if len(e.Fields) != 1 || e.Fields[0] != "purpose" {
			t.Errorf("วัตถุประสงค์ %q ต้องชี้ช่อง purpose แต่ได้ %v", bad, e.Fields)
		}
	}
}

// AC-BRQ-14 (FR-BRQ-02 ก)
// Given นักศึกษาเลือกวัตถุประสงค์ "อื่น ๆ"
// When  ไม่พิมพ์ข้อความระบุ แล้วกดส่งคำขอ
// Then  ระบบต้องไม่สร้างคำขอ และระบุว่าต้องพิมพ์ระบุวัตถุประสงค์
// และ   เมื่อพิมพ์ข้อความระบุแล้วส่งสำเร็จ ต้องเก็บทั้งค่า "อื่น ๆ" และข้อความที่พิมพ์
func TestAC_BRQ_14(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)

	// When: เลือก other แต่ไม่พิมพ์ระบุ
	w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose("CS-NB-0001", "other", ""))
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("เลือก \"อื่น ๆ\" โดยไม่พิมพ์ระบุต้องได้ 422 แต่ได้ %d", w.Code)
	}
	var e struct {
		Fields []string `json:"fields"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &e)
	if len(e.Fields) != 1 || e.Fields[0] != "purposeOther" {
		t.Errorf("ต้องชี้ช่อง purposeOther แต่ได้ %v", e.Fields)
	}
	// ช่องว่างล้วนก็ถือว่ายังไม่ได้พิมพ์
	if w := postRequest(t, srv, tokenStudentDept1, bodyWithPurpose("CS-NB-0001", "other", "   ")); w.Code != http.StatusUnprocessableEntity {
		t.Errorf("พิมพ์แต่ช่องว่างต้องได้ 422 แต่ได้ %d", w.Code)
	}

	// และ: พิมพ์ระบุแล้วต้องส่งได้ และเก็บทั้งสองค่า
	w = postRequest(t, srv, tokenStudentDept1, bodyWithPurpose("CS-NB-0001", "other", "ใช้ถ่ายทำวิดีโอส่งอาจารย์"))
	if w.Code != http.StatusCreated {
		t.Fatalf("เลือก \"อื่น ๆ\" พร้อมพิมพ์ระบุต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	got := decodeCreated(t, w)

	var purpose string
	var purposeOther *string
	if err := pool.QueryRow(context.Background(),
		`SELECT purpose, purpose_other FROM borrow_request WHERE id = $1`, got.RequestID,
	).Scan(&purpose, &purposeOther); err != nil {
		t.Fatal(err)
	}
	if purpose != "other" {
		t.Errorf("วัตถุประสงค์ต้องเป็น other แต่ได้ %q", purpose)
	}
	if purposeOther == nil || *purposeOther != "ใช้ถ่ายทำวิดีโอส่งอาจารย์" {
		t.Errorf("ข้อความที่พิมพ์ระบุต้องถูกเก็บไว้ แต่ได้ %v", purposeOther)
	}
}

// bodyWithWindow สร้าง JSON ของฟอร์มที่ครบทุกช่อง โดยกำหนดประเภทการยืมและช่วงเวลาเอง
func bodyWithWindow(assetCode, loanType string, start, end time.Time) string {
	return `{"assetCode":"` + assetCode + `","requesterName":"ธัญลักษณ์ ทดสอบ","studentCode":"660710715",` +
		`"loanType":"` + loanType + `","purpose":"study",` +
		`"startAt":"` + start.Format(time.RFC3339) + `","endAt":"` + end.Format(time.RFC3339) + `"}`
}

// bangkokToday คืนเวลาของวันนี้ตามเขตเวลา Asia/Bangkok ที่ชั่วโมงและนาทีที่ระบุ (glossary)
func bangkokToday(hour, min int) time.Time {
	now := time.Now().In(appLoc)
	return time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, appLoc)
}

// AC-BRQ-13 (FR-BRQ-03 ข้อ ค)
// Given นักศึกษาเลือกประเภทการยืมเป็น "ปกติ" และระบุเวลาเริ่มยืมวันนี้ 09:00 น.
// When  ระบุเวลาคืนเป็น 18:00 น. ของวันเดียวกัน แล้วกดส่งคำขอ
// Then  ระบบต้องปฏิเสธไม่ให้ส่ง ไม่สร้างคำขอ
//
//	(ข้อความที่บอกผู้ใช้ว่าต้องคืนภายใน 16:00 น. เป็นงานของ T-20 ตามที่ tasks.md แบ่งไว้)
func TestAC_BRQ_13(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)

	// When: ยืมปกติ 09:00 ถึง 18:00 ของวันเดียวกัน
	w := postRequest(t, srv, tokenStudentDept1,
		bodyWithWindow("CS-NB-0001", "normal", bangkokToday(9, 0), bangkokToday(18, 0)))
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("ยืมปกติที่คืนเกิน 16:00 ต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var e struct {
		Code   string   `json:"code"`
		Fields []string `json:"fields"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if e.Code != "VALIDATION_FAILED" || len(e.Fields) != 1 || e.Fields[0] != "endAt" {
		t.Errorf("ต้องชี้ช่อง endAt ด้วย VALIDATION_FAILED แต่ได้ %+v", e)
	}

	// Then: ต้องไม่มีคำขอเกิดขึ้น
	var count int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM borrow_request`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("คำขอที่ผิดกฎเวลาต้องไม่ถูกสร้าง แต่พบ %d ใบ", count)
	}
}

// T-08 เสร็จเมื่อ: กฎนี้ใช้กับการยืมปกติเท่านั้น และคืนพอดี 16:00 ต้องผ่าน (DOM-LOAN-01)
func TestLoanWindow_NormalBoundaryAndOvernight(t *testing.T) {
	t.Setenv("ATTACHMENT_DIR", t.TempDir()) // ยืมค้างคืนต้องแนบไฟล์ตั้งแต่ T-14 (FR-BRQ-10 ข้อ ก)
	srv := newTestServer(t)

	// คืนพอดี 16:00 ของวันเดียวกัน ต้องส่งได้ (ข้อกำหนดคือ "ภายใน 16:00 น.")
	w := postRequest(t, srv, tokenStudentDept1,
		bodyWithWindow("CS-NB-0001", "normal", bangkokToday(9, 0), bangkokToday(16, 0)))
	if w.Code != http.StatusCreated {
		t.Errorf("คืนพอดี 16:00 ต้องส่งได้ แต่ได้ %d body=%s", w.Code, w.Body.String())
	} else {
		decodeCreated(t, w)
	}

	// เกินไปนาทีเดียวต้องไม่ผ่าน
	if w := postRequest(t, srv, tokenStudentDept1,
		bodyWithWindow("CS-NB-0002", "normal", bangkokToday(9, 0), bangkokToday(16, 1))); w.Code != http.StatusUnprocessableEntity {
		t.Errorf("คืน 16:01 ต้องได้ 422 แต่ได้ %d", w.Code)
	}

	// ยืมค้างคืนข้ามวันต้องไม่ติดกฎนี้ (DOM-LOAN-01 ใช้กับยืมปกติเท่านั้น)
	// ตั้งแต่ T-14 การยืมค้างคืนต้องแนบใบบันทึกข้อความมาด้วย จึงส่งแบบ multipart (FR-BRQ-10 ข้อ ก)
	ct, body := multipartForm(t, "CS-NB-0003", config.LoanTypeOvernight,
		&testFile{name: "memo.pdf", contentType: "application/pdf", content: []byte("%PDF-1.4 ตัวอย่าง")})
	w = postMultipart(t, srv, tokenStudentDept1, ct, body)
	if w.Code != http.StatusCreated {
		t.Errorf("ยืมค้างคืนข้ามวันต้องส่งได้ แต่ได้ %d body=%s", w.Code, w.Body.String())
	} else {
		decodeCreated(t, w)
	}
}

// insertTempItem เพิ่มอุปกรณ์ชั่วคราวสถานะ "ว่าง" แล้วลบทิ้งหลัง test (ไม่แตะข้อมูลตั้งต้นของ 001)
func insertTempItem(t *testing.T, assetCode string) {
	t.Helper()
	pool := testPool(t)
	ctx := context.Background()
	if _, err := pool.Exec(ctx,
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, 'available')`, assetCode); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(),
			`DELETE FROM borrow_request WHERE item_id IN (SELECT id FROM equipment_item WHERE asset_code = $1)`, assetCode)
		_, _ = pool.Exec(context.Background(), `DELETE FROM equipment_item WHERE asset_code = $1`, assetCode)
	})
}

// AC-BRQ-01 (FR-BRQ-01)
// Given นักศึกษาเปิดแบบฟอร์มคำขอของอุปกรณ์ชิ้นที่สถานะ "ว่าง"
// When  อุปกรณ์ชิ้นนั้นถูกเปลี่ยนเป็น "ถูกยืม" ระหว่างกรอก แล้วนักศึกษากดส่งคำขอ
// Then  ระบบต้องไม่สร้างคำขอ และแจ้งสถานะล่าสุดของอุปกรณ์
func TestAC_BRQ_01(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	ctx := context.Background()
	const assetCode = "TMP-NB-0901"

	// Given: อุปกรณ์ว่าง ผู้ขอเปิดฟอร์มตอนนี้
	insertTempItem(t, assetCode)
	var atFormOpen string
	if err := pool.QueryRow(ctx,
		`SELECT status FROM equipment_item WHERE asset_code = $1`, assetCode).Scan(&atFormOpen); err != nil {
		t.Fatal(err)
	}
	if atFormOpen != "available" {
		t.Fatalf("ตอนเปิดฟอร์มอุปกรณ์ต้องว่าง แต่เป็น %q", atFormOpen)
	}

	// When: ระหว่างกรอก อุปกรณ์ถูกยืมไป (จำลอง UC-08 ของ 005)
	if _, err := pool.Exec(ctx,
		`UPDATE equipment_item SET status = 'borrowed' WHERE asset_code = $1`, assetCode); err != nil {
		t.Fatal(err)
	}

	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then: ต้องได้ 409 ITEM_NOT_AVAILABLE ตามรหัสกลางของ CON-TECH-03
	if w.Code != http.StatusConflict {
		t.Fatalf("อุปกรณ์ไม่ว่างแล้วต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var e struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if e.Code != "ITEM_NOT_AVAILABLE" {
		t.Errorf("รหัสต้องเป็น ITEM_NOT_AVAILABLE แต่ได้ %q", e.Code)
	}
	// Then: ต้องแจ้งสถานะล่าสุดของอุปกรณ์กลับไปด้วย (FR-BRQ-01)
	if e.Status != "borrowed" {
		t.Errorf("ต้องแจ้งสถานะล่าสุดว่า borrowed แต่ได้ %q", e.Status)
	}

	// Then: ต้องไม่มีคำขอเกิดขึ้น
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM borrow_request`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("อุปกรณ์ไม่ว่างต้องไม่สร้างคำขอ แต่พบ %d ใบ", count)
	}

	// Then: error ที่ service คืนต้องพกสถานะล่าสุดมาด้วย เพื่อให้แจ้งผู้ขอได้ (FR-BRQ-01)
	svc := NewServiceWithClock(NewRepository(pool), testClock())
	_, err := svc.Create(studentCtx("student-001"), newTestRequest(assetCode))
	var notAvail *ItemNotAvailableError
	if !errors.As(err, &notAvail) {
		t.Fatalf("ต้องได้ ItemNotAvailableError ที่พกสถานะมาด้วย แต่ได้ %v", err)
	}
	if notAvail.Status != "borrowed" {
		t.Errorf("สถานะล่าสุดต้องเป็น borrowed แต่ได้ %q", notAvail.Status)
	}
	if !errors.Is(err, ErrItemNotAvailable) {
		t.Error("ต้องยังเทียบกับ ErrItemNotAvailable ได้ เพื่อให้โค้ดเดิมใช้ได้เหมือนเดิม")
	}
}

// T-09 เสร็จเมื่อ: รหัสครุภัณฑ์ที่ไม่มีอยู่จริง ต้องตอบเหมือนกรณีไม่ว่าง ไม่เปิดเผยว่ามีรหัสนั้นไหม
func TestCreate_UnknownAssetCodeLooksTheSame(t *testing.T) {
	srv := newTestServer(t)

	w := postRequest(t, srv, tokenStudentDept1, bodyFor("NO-SUCH-9999"))
	if w.Code != http.StatusConflict {
		t.Fatalf("รหัสที่ไม่มีอยู่ต้องได้ 409 เหมือนกรณีไม่ว่าง แต่ได้ %d", w.Code)
	}
	var e struct {
		Code   string `json:"code"`
		Status string `json:"status"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &e)
	if e.Code != "ITEM_NOT_AVAILABLE" {
		t.Errorf("รหัสต้องเป็น ITEM_NOT_AVAILABLE แต่ได้ %q", e.Code)
	}
	// ไม่มีรหัสนี้อยู่จริง จึงไม่มีสถานะให้แจ้ง ต้องไม่เปิดเผยว่ามีรหัสนั้นในระบบไหม
	if e.Status != "" {
		t.Errorf("รหัสที่ไม่มีอยู่ต้องไม่มีสถานะใน response แต่ได้ %q", e.Status)
	}
}

// ผู้ใช้จำลองของภาควิชาที่ 2 (ดู backend/internal/auth/mock_provider.go)
const tokenStudentDept2 = "mock-student-101"

// AC-BRQ-07 (FR-BRQ-07)
// Given นักศึกษาสังกัดภาควิชาที่ 2 และอุปกรณ์ชิ้นเป้าหมายเป็นของภาควิชาที่ 1
// When  ส่งคำขอยืมอุปกรณ์ชิ้นนั้น
// Then  ระบบต้องปฏิเสธ ไม่สร้างคำขอ และไม่ส่งข้อมูลอุปกรณ์กลับไป
func TestAC_BRQ_07(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)

	// When: student-101 (ภาควิชา 2) ขอยืม CS-NB-0001 ซึ่งเป็นของภาควิชา 1
	w := postRequest(t, srv, tokenStudentDept2, bodyFor("CS-NB-0001"))

	// Then: 403 FORBIDDEN_DEPARTMENT ตามรหัสกลางของ CON-TECH-03
	if w.Code != http.StatusForbidden {
		t.Fatalf("คำขอข้ามภาควิชาต้องได้ 403 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var e struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if e.Code != "FORBIDDEN_DEPARTMENT" {
		t.Errorf("รหัสต้องเป็น FORBIDDEN_DEPARTMENT แต่ได้ %q", e.Code)
	}

	// Then: body ต้องไม่มีข้อมูลอุปกรณ์เลย ทั้งสถานะ ชื่อรุ่น และรหัสครุภัณฑ์
	if e.Status != "" {
		t.Errorf("ต้องไม่บอกสถานะอุปกรณ์ให้คนต่างภาควิชารู้ แต่ได้ %q", e.Status)
	}
	for _, leak := range []string{"available", "borrowed", "unavailable", "Notebook", "CS-NB-0001"} {
		if strings.Contains(w.Body.String(), leak) {
			t.Errorf("body ต้องไม่มีข้อมูลอุปกรณ์ แต่พบคำว่า %q", leak)
		}
	}

	// Then: ต้องไม่มีคำขอเกิดขึ้น
	var count int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM borrow_request`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("คำขอข้ามภาควิชาต้องไม่ถูกสร้าง แต่พบ %d ใบ", count)
	}
}

// T-10 เสร็จเมื่อ: กฎนี้กันเฉพาะข้ามภาควิชา ผู้ขอยังส่งคำขอของภาควิชาตนเองได้ตามปกติ (DOM-ACCESS-01)
func TestDepartmentScope_OwnDepartmentStillWorks(t *testing.T) {
	srv := newTestServer(t)

	// student-101 อยู่ภาควิชา 2 ขอยืม XX-NB-0001 ซึ่งเป็นของภาควิชา 2 ต้องผ่าน
	w := postRequest(t, srv, tokenStudentDept2, bodyFor("XX-NB-0001"))
	if w.Code != http.StatusCreated {
		t.Fatalf("ขอยืมอุปกรณ์ของภาควิชาตนเองต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)

	// และคนภาควิชา 1 ขอยืมของภาควิชา 2 ก็ต้องถูกปฏิเสธเช่นกัน (กันทั้งสองทาง)
	if w := postRequest(t, srv, tokenStudentDept1, bodyFor("XX-NB-0001")); w.Code != http.StatusForbidden {
		t.Errorf("คนภาควิชา 1 ขอยืมของภาควิชา 2 ต้องได้ 403 แต่ได้ %d", w.Code)
	}
}

// AC-BRQ-11 (FR-BRQ-05 ข) ยิงผ่าน HTTP จริง
// Given นักศึกษาคนหนึ่งมีคำขอของอุปกรณ์ชิ้นนั้นอยู่แล้ว (ใบแรกได้ 201 สถานะ requested)
// When  POST คำขอของอุปกรณ์ชิ้นเดิมซ้ำด้วยผู้ใช้คนเดิม
// Then  ต้องได้ 409 รหัส DUPLICATE_REQUEST ตามรายการของ CON-TECH-03 พร้อมข้อความแจ้งผู้ขอ
func TestAC_BRQ_11_HTTP(t *testing.T) {
	srv := newTestServer(t)
	const assetCode = "CS-NB-0003" // อุปกรณ์ว่างของภาควิชา 1 ในข้อมูลตั้งต้น (T-03)

	// Given: decodeCreated ลบคำขอใบนี้ให้ตอน test จบ
	decodeCreated(t, postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode)))

	// When
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then
	if w.Code != http.StatusConflict {
		t.Fatalf("ส่งซ้ำต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorDTO
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
	}
	if got.Code != "DUPLICATE_REQUEST" {
		t.Errorf("รหัสต้องเป็น DUPLICATE_REQUEST ตาม CON-TECH-03 แต่ได้ %q", got.Code)
	}
	if got.Message == "" {
		t.Error("ต้องมีข้อความแจ้งว่ามีคำขอของอุปกรณ์ชิ้นนี้ค้างอยู่แล้ว")
	}
}

// ผู้ใช้จำลองที่ข้อมูลตั้งต้นให้มีอุปกรณ์คืนเกินกำหนดค้างอยู่ (T-03)
const tokenStudentOverdue = "mock-student-002"

// AC-BRQ-06 (FR-BRQ-06) ยิงผ่าน HTTP จริง ตรงกับจุดตรวจเฟส 8
// Given ล็อกอินเป็น student-002 ซึ่งมีอุปกรณ์เลยกำหนดคืนค้างอยู่
// When  POST /v1/loan-requests
// Then  ได้ 409 รหัส HAS_OVERDUE_LOAN ตามรายการของ CON-TECH-03 พร้อมเหตุผล
func TestAC_BRQ_06_HTTP(t *testing.T) {
	srv := newTestServer(t)

	w := postRequest(t, srv, tokenStudentOverdue, bodyFor("CS-NB-0003"))

	if w.Code != http.StatusConflict {
		t.Fatalf("ผู้ขอที่มีของเกินกำหนดต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var got ErrorDTO
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
	}
	if got.Code != "HAS_OVERDUE_LOAN" {
		t.Errorf("รหัสต้องเป็น HAS_OVERDUE_LOAN ตาม CON-TECH-03 แต่ได้ %q", got.Code)
	}
	if got.Message == "" {
		t.Error("ต้องมีข้อความแจ้งเหตุผลว่ามีของค้างเกินกำหนด (FR-BRQ-06)")
	}
}

// testFile คือไฟล์ที่ใช้ทดสอบการแนบใบบันทึกข้อความ
type testFile struct {
	name        string
	contentType string
	content     []byte
}

// multipartForm สร้าง body แบบ multipart/form-data ของฟอร์มคำขอ (plan.md ข้อ 4)
// file เป็น nil แปลว่าไม่ได้แนบไฟล์มา
func multipartForm(t *testing.T, assetCode, loanType string, file *testFile) (string, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	end := bangkokToday(15, 0)
	if loanType == config.LoanTypeOvernight {
		end = end.AddDate(0, 0, 2) // ค้างคืนไม่ติดกฎคืนภายใน 16:00 วันเดียวกัน (DOM-LOAN-01)
	}
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
		// ใช้ CreatePart เพื่อกำหนด Content-Type ของไฟล์เอง แทน CreateFormFile ที่ตั้งเป็น octet-stream เสมอ
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

// postMultipart ยิง POST /v1/loan-requests แบบ multipart
func postMultipart(t *testing.T, srv http.Handler, token, contentType string, body *bytes.Buffer) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, PathPrefix, body)
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)
	return w
}

// expectAttachmentRejected ยืนยันว่าได้ 422 พร้อมชี้ช่อง attachment และไม่มีคำขอเกิดขึ้น
func expectAttachmentRejected(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("ต้องได้ 422 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var got validationErrorDTO
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
	}
	if got.Code != "VALIDATION_FAILED" {
		t.Errorf("รหัสต้องเป็น VALIDATION_FAILED ตาม CON-TECH-03 แต่ได้ %q", got.Code)
	}
	if !slices.Contains(got.Fields, "attachment") {
		t.Errorf("ต้องระบุว่าช่อง attachment คือสิ่งที่ยังขาดหรือผิด แต่ได้ %v", got.Fields)
	}
}

// AC-BRQ-10 (FR-BRQ-10)
// Given นักศึกษาเลือกประเภทการยืมเป็น "ค้างคืน"
// When  ยังไม่แนบไฟล์ใบบันทึกข้อความ แล้วกดส่งคำขอ
// Then  ระบบต้องไม่สร้างคำขอ และระบุสิ่งที่ยังขาด
// และ
// Given แนบไฟล์ที่ไม่ใช่ PDF หรือรูปถ่าย หรือไฟล์ขนาดเกิน 10MB
// When  กดส่งคำขอ
// Then  ระบบต้องไม่สร้างคำขอ และแจ้งชนิดกับขนาดไฟล์ที่รับได้
//
// ส่วน "ยังไม่เลือกวัตถุประสงค์" ของ AC ข้อนี้ถูกบังคับไว้แล้วตั้งแต่ T-07 (TestAC_BRQ_02)
// และส่วนที่ให้ "แสดงเงื่อนไขการยืมปกติ" เป็นงานหน้าจอ ซึ่งรอ Q-BRQ-14 (T-19, T-20)
func TestAC_BRQ_10(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ATTACHMENT_DIR", dir) // ต้องตั้งก่อนสร้าง service เพราะอ่าน config ตอนสร้าง
	pool := testPool(t)
	srv := newTestServer(t)

	t.Run("ค้างคืนแต่ไม่แนบไฟล์", func(t *testing.T) {
		ct, body := multipartForm(t, "CS-NB-0003", config.LoanTypeOvernight, nil)
		expectAttachmentRejected(t, postMultipart(t, srv, tokenStudentDept1, ct, body))
	})

	t.Run("ไฟล์ชนิดที่ไม่รับ", func(t *testing.T) {
		ct, body := multipartForm(t, "CS-NB-0003", config.LoanTypeOvernight,
			&testFile{name: "memo.txt", contentType: "text/plain", content: []byte("ไม่ใช่ PDF หรือรูปถ่าย")})
		expectAttachmentRejected(t, postMultipart(t, srv, tokenStudentDept1, ct, body))
	})

	t.Run("ไฟล์ใหญ่เกินเพดาน", func(t *testing.T) {
		big := make([]byte, config.AttachmentMaxBytes+1)
		copy(big, []byte("%PDF-1.4"))
		ct, body := multipartForm(t, "CS-NB-0003", config.LoanTypeOvernight,
			&testFile{name: "big.pdf", contentType: "application/pdf", content: big})
		expectAttachmentRejected(t, postMultipart(t, srv, tokenStudentDept1, ct, body))
	})

	// ถึงตรงนี้ยังต้องไม่มีคำขอของผู้ขอคนนี้เกิดขึ้นเลย (AC-BRQ-10 "ต้องไม่สร้างคำขอ")
	if n := countRequests(t, pool, "student-001"); n != 0 {
		t.Fatalf("คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}

	t.Run("แนบ PDF ที่ถูกต้อง", func(t *testing.T) {
		content := []byte("%PDF-1.4 ใบบันทึกข้อความตัวอย่าง")
		ct, body := multipartForm(t, "CS-NB-0003", config.LoanTypeOvernight,
			&testFile{name: "ใบบันทึก ข้อความ.pdf", contentType: "application/pdf", content: content})

		w := postMultipart(t, srv, tokenStudentDept1, ct, body)
		if w.Code != http.StatusCreated {
			t.Fatalf("แนบไฟล์ถูกต้องต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
		}
		got := decodeCreated(t, w)

		// Then: บันทึก path ไว้ในตาราง (FR-BRQ-10 ข้อ ก)
		var path *string
		if err := pool.QueryRow(context.Background(),
			`SELECT attachment_path FROM borrow_request WHERE id = $1`, got.RequestID).Scan(&path); err != nil {
			t.Fatal(err)
		}
		if path == nil {
			t.Fatal("คำขอค้างคืนที่แนบไฟล์แล้ว ต้องมี attachment_path ในตาราง")
		}

		// Then: ไฟล์อยู่ในโฟลเดอร์ตาม config และชื่อขึ้นต้นด้วยรหัสคำขอ (plan.md ข้อ 3, ASM-BRQ-07)
		if gotDir := filepath.Dir(*path); gotDir != dir {
			t.Errorf("ไฟล์ต้องอยู่ในโฟลเดอร์จาก config (%s) แต่อยู่ที่ %s", dir, gotDir)
		}
		if prefix := fmt.Sprintf("%d-", got.RequestID); !strings.HasPrefix(filepath.Base(*path), prefix) {
			t.Errorf("ชื่อไฟล์ต้องขึ้นต้นด้วย %q แต่ได้ %q", prefix, filepath.Base(*path))
		}
		saved, err := os.ReadFile(*path)
		if err != nil {
			t.Fatalf("ต้องเก็บไฟล์ไว้จริงบนเครื่อง: %v", err)
		}
		if !bytes.Equal(saved, content) {
			t.Error("เนื้อไฟล์ที่เก็บต้องตรงกับไฟล์ที่ผู้ขอแนบมา")
		}
	})
}

// ผู้ใช้จำลองที่เป็นเจ้าหน้าที่ (ดู backend/internal/auth/mock_provider.go)
const tokenStaff = "mock-staff-001"

// getAttachment ยิง GET /v1/loan-requests/{id}/attachment ด้วย token ที่ให้
func getAttachment(t *testing.T, srv http.Handler, token string, id int) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d/attachment", PathPrefix, id), nil)
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)
	return w
}

// AC-BRQ-10 (FR-BRQ-10 ข้อ ก) ส่วนของการเปิดดูไฟล์แนบ
// Given คำขอยืมค้างคืนของ student-001 ที่แนบใบบันทึกข้อความไว้แล้ว
// When  เจ้าหน้าที่ เจ้าของคำขอ และนักศึกษาคนอื่น เปิดไฟล์แนบของคำขอใบนั้น
// Then  สองคนแรกต้องได้ไฟล์ (200) ส่วนนักศึกษาคนอื่นต้องถูกปฏิเสธ (403) และต้องไม่ได้เนื้อไฟล์
//
// (ทีมตัดสิน 2569-11-12: ใช้รหัส FORBIDDEN_ROLE ที่ลงทะเบียนไว้แล้ว และ REQUEST_NOT_FOUND สำหรับ 404)
func TestAC_BRQ_10_AttachmentAccess(t *testing.T) {
	t.Setenv("ATTACHMENT_DIR", t.TempDir())
	srv := newTestServer(t)

	// Given: คำขอค้างคืนที่แนบไฟล์แล้ว 1 ใบ เจ้าของคือ student-001
	content := []byte("%PDF-1.4 ใบบันทึกข้อความของ student-001")
	ct, body := multipartForm(t, "CS-NB-0003", config.LoanTypeOvernight,
		&testFile{name: "memo.pdf", contentType: "application/pdf", content: content})
	w := postMultipart(t, srv, tokenStudentDept1, ct, body)
	if w.Code != http.StatusCreated {
		t.Fatalf("เตรียมคำขอไม่สำเร็จ ได้ %d body=%s", w.Code, w.Body.String())
	}
	created := decodeCreated(t, w)

	t.Run("เจ้าหน้าที่เปิดได้", func(t *testing.T) {
		got := getAttachment(t, srv, tokenStaff, created.RequestID)
		if got.Code != http.StatusOK {
			t.Fatalf("เจ้าหน้าที่ต้องได้ 200 แต่ได้ %d body=%s", got.Code, got.Body.String())
		}
		if !bytes.Equal(got.Body.Bytes(), content) {
			t.Error("ไฟล์ที่ส่งกลับต้องเป็นไฟล์เดียวกับที่ผู้ขอแนบมา")
		}
	})

	t.Run("เจ้าของคำขอเปิดได้", func(t *testing.T) {
		got := getAttachment(t, srv, tokenStudentDept1, created.RequestID)
		if got.Code != http.StatusOK {
			t.Fatalf("เจ้าของคำขอต้องได้ 200 แต่ได้ %d body=%s", got.Code, got.Body.String())
		}
		if !bytes.Equal(got.Body.Bytes(), content) {
			t.Error("ไฟล์ที่ส่งกลับต้องเป็นไฟล์เดียวกับที่ผู้ขอแนบมา")
		}
	})

	t.Run("นักศึกษาคนอื่นเปิดไม่ได้", func(t *testing.T) {
		got := getAttachment(t, srv, tokenStudentOverdue, created.RequestID) // student-002 ภาควิชาเดียวกัน
		if got.Code != http.StatusForbidden {
			t.Fatalf("นักศึกษาคนอื่นต้องได้ 403 แต่ได้ %d body=%s", got.Code, got.Body.String())
		}
		var e ErrorDTO
		if err := json.Unmarshal(got.Body.Bytes(), &e); err != nil {
			t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
		}
		if e.Code != "FORBIDDEN_ROLE" {
			t.Errorf("รหัสต้องเป็น FORBIDDEN_ROLE ตาม CON-TECH-03 แต่ได้ %q", e.Code)
		}
		if bytes.Contains(got.Body.Bytes(), content) {
			t.Error("คำตอบต้องไม่มีเนื้อไฟล์ของคนอื่นติดไปด้วย (DOM-PRIV-01)")
		}
	})

	t.Run("ไม่มีคำขอใบนั้น", func(t *testing.T) {
		got := getAttachment(t, srv, tokenStaff, created.RequestID+100000)
		if got.Code != http.StatusNotFound {
			t.Fatalf("คำขอที่ไม่มีอยู่ต้องได้ 404 แต่ได้ %d body=%s", got.Code, got.Body.String())
		}
		var e ErrorDTO
		if err := json.Unmarshal(got.Body.Bytes(), &e); err != nil {
			t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
		}
		if e.Code != "REQUEST_NOT_FOUND" {
			t.Errorf("รหัสต้องเป็น REQUEST_NOT_FOUND แต่ได้ %q", e.Code)
		}
	})

	t.Run("ไม่มี token", func(t *testing.T) {
		if got := getAttachment(t, srv, "", created.RequestID); got.Code != http.StatusUnauthorized {
			t.Errorf("ไม่มี token ต้องได้ 401 แต่ได้ %d", got.Code)
		}
	})
}

// T-14 เสร็จเมื่อ (ส่วนเสริม): ชื่อไฟล์จากผู้ขอต้องพาไฟล์ออกนอกโฟลเดอร์ที่กำหนดไม่ได้ (ASM-BRQ-07)
func TestSafeFileName(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"memo.pdf", "memo.pdf"},
		{"../../etc/passwd", "passwd"},
		{`C:\Users\USER\memo.pdf`, "memo.pdf"},
		{"ใบบันทึก ข้อความ.pdf", "_.pdf"}, // อักษรไทยกลายเป็นตัวคั่น แต่ต้องไม่ทำให้นามสกุลหาย
		{"...", "attachment"},
	} {
		if got := safeFileName(c.in); got != c.want {
			t.Errorf("safeFileName(%q) = %q ต้องเป็น %q", c.in, got, c.want)
		}
	}
}
