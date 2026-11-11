package borrowrequest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	t.Helper()
	svc := NewService(NewRepository(testPool(t)))
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
