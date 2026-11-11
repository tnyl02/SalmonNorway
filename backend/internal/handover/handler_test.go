package handover

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/auth"
)

// ผู้ใช้จำลองจาก auth.NewMockProvider
const (
	tokenStaff    = "mock-staff-001"
	tokenStudentA = "mock-student-001"
	tokenStudentB = "mock-student-101"
)

// routesOnly คือ router ที่ผ่าน auth.Middleware เหมือนใน main.go (ใช้กับ test ที่ไม่แตะฐานข้อมูล)
func routesOnly(db DB) http.Handler {
	provider := auth.NewMockProvider()
	return auth.Middleware(provider, NewHandler(db, provider).Routes())
}

func doRequest(h http.Handler, method, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

var allRoutes = []struct{ method, path string }{
	{http.MethodGet, PathPrefix + "requests"},
	{http.MethodPost, PathPrefix + "requests/1/verify-borrower"},
	{http.MethodPost, PathPrefix + "requests/1/handover"},
}

// T-03 เสร็จเมื่อ: ทั้ง 3 route ตอบ 403 {code, message} FORBIDDEN_ROLE เมื่อเป็นนักศึกษา และ 401 เมื่อไม่มี token (CON-HND-01, IF-IDP-01, CON-TECH-03)
func TestHandover_NonStaffForbidden(t *testing.T) {
	h := routesOnly(nil)
	for _, rt := range allRoutes {
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			rec := doRequest(h, rt.method, rt.path, tokenStudentA)
			if rec.Code != http.StatusForbidden {
				t.Fatalf("นักศึกษา: HTTP %d ต้องเป็น 403", rec.Code)
			}
			var body ErrorDTO
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body ไม่ใช่ JSON: %s", rec.Body.String())
			}
			if body.Code != "FORBIDDEN_ROLE" || body.Message == "" {
				t.Errorf("body = %+v ต้องเป็น FORBIDDEN_ROLE พร้อม message", body)
			}

			if rec := doRequest(h, rt.method, rt.path, ""); rec.Code != http.StatusUnauthorized {
				t.Errorf("ไม่มี token: HTTP %d ต้องเป็น 401", rec.Code)
			}
		})
	}
}

// testPool คืน connection pool ไปยังฐานข้อมูลของ test (ใช้กับ test ผ่าน HTTP ที่ต้อง commit จริง)
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	testDB(t)
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close) // ลงทะเบียนก่อน จึงปิดหลังการลบข้อมูลทดสอบ (Cleanup ทำงานย้อนลำดับ)
	return pool
}

// createItem เพิ่มอุปกรณ์รายชิ้นของรุ่น 1 แบบ commit จริง และลบทุกอย่างที่ผูกกับชิ้นนี้ตอน test จบ
func createItem(t *testing.T, pool *pgxpool.Pool, assetCode, status string) int {
	t.Helper()
	ctx := context.Background()
	var id int
	if err := pool.QueryRow(ctx,
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, $2) RETURNING id`,
		assetCode, status).Scan(&id); err != nil {
		t.Fatalf("เพิ่มอุปกรณ์ไม่สำเร็จ: %v", err)
	}
	t.Cleanup(func() {
		for _, q := range []string{
			`DELETE FROM handover_photo WHERE loan_id IN (SELECT id FROM loan WHERE item_id = $1)`,
			`DELETE FROM handover WHERE loan_id IN (SELECT id FROM loan WHERE item_id = $1)`,
			`DELETE FROM audit_log WHERE request_id IN (SELECT id FROM borrow_request WHERE item_id = $1)`,
			`DELETE FROM loan WHERE item_id = $1`,
			`DELETE FROM borrow_request WHERE item_id = $1`,
			`DELETE FROM equipment_item WHERE id = $1`,
		} {
			if _, err := pool.Exec(ctx, q, id); err != nil {
				t.Errorf("ลบข้อมูลทดสอบไม่สำเร็จ: %v", err)
			}
		}
	})
	return id
}

// createRequest เพิ่มคำขอยืมแบบปกติแบบ commit จริง คืน id (ข้อมูลลบโดย cleanup ของ createItem)
func createRequest(t *testing.T, pool *pgxpool.Pool, requesterID string, itemID int, status string) int {
	t.Helper()
	var id int
	if err := pool.QueryRow(context.Background(),
		`INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
		 VALUES ($1, 'ผู้ขอทดสอบ', $1, $2, 'normal', 'study', now(), now(), $3) RETURNING id`,
		requesterID, itemID, status).Scan(&id); err != nil {
		t.Fatalf("เพิ่มคำขอยืมไม่สำเร็จ: %v", err)
	}
	return id
}

// AC-HND-14 (FR-HND-01)
// Given มีคำขอที่ได้รับอนุมัติแล้วและยังไม่ส่งมอบ 2 รายการ (อุปกรณ์คนละชิ้น) คำขอที่ยังไม่อนุมัติ 1 รายการ และคำขอที่ส่งมอบแล้ว 1 รายการ
// When  เจ้าหน้าที่เปิดรายการคำขอเพื่อส่งมอบ
// Then  รายการแสดงเฉพาะคำขอที่ได้รับอนุมัติแล้วและยังไม่ส่งมอบทั้ง 2 รายการ
func TestAC_HND_14_ListApprovedOnly(t *testing.T) {
	pool := testPool(t)
	approved1 := createRequest(t, pool, "student-001", createItem(t, pool, "T14-AP-0001", "available"), "approved")
	approved2 := createRequest(t, pool, "student-002", createItem(t, pool, "T14-AP-0002", "available"), "approved")
	createRequest(t, pool, "student-001", createItem(t, pool, "T14-RQ-0001", "available"), "requested")
	createRequest(t, pool, "student-001", createItem(t, pool, "T14-HO-0001", "borrowed"), "handed_over")

	rec := doRequest(routesOnly(pool), http.MethodGet, PathPrefix+"requests", tokenStaff)
	if rec.Code != http.StatusOK {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
	}
	var got []ReadyRequestDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("body ไม่ใช่ JSON array: %s", rec.Body.String())
	}
	want := map[int]string{approved1: "T14-AP-0001", approved2: "T14-AP-0002"}
	if len(got) != len(want) {
		t.Fatalf("ได้ %d รายการ ต้องเป็น 2: %+v", len(got), got)
	}
	for _, d := range got {
		if want[d.RequestID] != d.AssetCode {
			t.Errorf("รายการที่ไม่คาดคิด: %+v", d)
		}
		if d.ModelName == "" || d.RequesterStudentID == "" || d.SubmittedAt.IsZero() {
			t.Errorf("ข้อมูลไม่ครบ: %+v", d)
		}
	}
}

func doJSON(h http.Handler, method, path, token, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func verifyPath(requestID int) string {
	return PathPrefix + "requests/" + strconv.Itoa(requestID) + "/verify-borrower"
}

// AC-HND-02 (FR-HND-02, FR-HND-03)
// Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว Auth API ยืนยันตัวตนผู้ยืมสำเร็จ และผู้ยืมเป็นผู้ส่งคำขอ
// When  เจ้าหน้าที่ยืนยันว่าตรวจบัตรแล้ว และบันทึกการส่งมอบ
// Then  ระบบบันทึกการส่งมอบ พร้อมชื่อ นามสกุล รหัสนักศึกษาของผู้ยืม และวันที่และเวลาที่ส่งมอบ
//
//	ถ้าเจ้าหน้าที่ไม่ยืนยันการตรวจบัตร ระบบไม่บันทึกการส่งมอบ
//
// subtest verify_borrower (T-09) ส่วนบันทึกการส่งมอบอยู่ใน T-10
func TestAC_HND_02_AuthApiVerifiedAndCardChecked(t *testing.T) {
	pool := testPool(t)
	h := routesOnly(pool)
	requestID := createRequest(t, pool, "student-001", createItem(t, pool, "T02-VB-0001", "available"), "approved")

	t.Run("verify_borrower", func(t *testing.T) {
		rec := doJSON(h, http.MethodPost, verifyPath(requestID), tokenStaff, `{"borrowerCredential":"`+tokenStudentA+`"}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
		}
		var got VerifyBorrowerResult
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got.Result != "verified" || got.UserID != "student-001" {
			t.Errorf("ผล = %+v ต้องเป็น verified ของ student-001", got)
		}
	})
}

// T-09: ผลอื่นของ verify-borrower ผู้ยืมไม่ใช่ผู้ส่งคำขอ = mismatch, ล็อกอินไม่สำเร็จ = 422, ไม่มีคำขอ = 404 (FR-HND-02, ASM-HND-04)
func TestVerifyBorrower_OtherResults(t *testing.T) {
	pool := testPool(t)
	h := routesOnly(pool)
	requestID := createRequest(t, pool, "student-001", createItem(t, pool, "T09-VB-0001", "available"), "approved")

	rec := doJSON(h, http.MethodPost, verifyPath(requestID), tokenStaff, `{"borrowerCredential":"`+tokenStudentB+`"}`)
	var got VerifyBorrowerResult
	_ = json.Unmarshal(rec.Body.Bytes(), &got)
	if rec.Code != http.StatusOK || got.Result != "mismatch" {
		t.Errorf("ผู้ยืมคนอื่น: HTTP %d ผล %+v ต้องเป็น 200 mismatch", rec.Code, got)
	}

	for _, body := range []string{`{"borrowerCredential":"mock-unknown"}`, `{}`, `not json`} {
		rec := doJSON(h, http.MethodPost, verifyPath(requestID), tokenStaff, body)
		var e ErrorDTO
		_ = json.Unmarshal(rec.Body.Bytes(), &e)
		if rec.Code != http.StatusUnprocessableEntity || e.Code != "VALIDATION_FAILED" {
			t.Errorf("body %s: HTTP %d %+v ต้องเป็น 422 VALIDATION_FAILED", body, rec.Code, e)
		}
	}

	if rec := doJSON(h, http.MethodPost, verifyPath(999999), tokenStaff, `{"borrowerCredential":"`+tokenStudentA+`"}`); rec.Code != http.StatusNotFound {
		t.Errorf("ไม่มีคำขอ: HTTP %d ต้องเป็น 404", rec.Code)
	}
}

// testJPEG คือข้อมูลที่ขึ้นต้นด้วย signature ของ JPEG ยาว size byte (ใช้แทนรูปจริงใน test)
func testJPEG(size int) []byte {
	b := make([]byte, size)
	copy(b, []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00})
	return b
}

// photoFile คือรูป 1 ไฟล์ที่แนบในฟอร์ม
type photoFile struct {
	name        string
	contentType string
	data        []byte
}

// twoJPEGs คือรูปขั้นต่ำที่ถูกต้อง JPG 2 รูป (FR-HND-06)
func twoJPEGs() []photoFile {
	return []photoFile{
		{"1.jpg", "image/jpeg", testJPEG(1024)},
		{"2.jpg", "image/jpeg", testJPEG(1024)},
	}
}

// validForm คือช่องในฟอร์มบันทึกการส่งมอบที่ถูกต้องของผู้ยืม credential
func validForm(credential string) map[string]string {
	return map[string]string{
		"borrowerCredential": credential,
		"borrowerFirstName":  "สมชาย",
		"borrowerLastName":   "ใจดี",
		"cardChecked":        "true",
	}
}

// postHandover ส่งฟอร์ม multipart ไปที่ POST .../handover
func postHandover(t *testing.T, h http.Handler, token string, requestID int, fields map[string]string, photos []photoFile) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for k, v := range fields {
		_ = mw.WriteField(k, v)
	}
	for _, p := range photos {
		hdr := textproto.MIMEHeader{}
		hdr.Set("Content-Disposition", `form-data; name="photos"; filename="`+p.name+`"`)
		hdr.Set("Content-Type", p.contentType)
		part, err := mw.CreatePart(hdr)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = part.Write(p.data)
	}
	_ = mw.Close()
	req := httptest.NewRequest(http.MethodPost, PathPrefix+"requests/"+strconv.Itoa(requestID)+"/handover", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func countLoans(t *testing.T, pool *pgxpool.Pool, itemID int) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM loan WHERE item_id = $1`, itemID).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func statusOf(t *testing.T, pool *pgxpool.Pool, table string, id int) string {
	t.Helper()
	var s string
	if err := pool.QueryRow(context.Background(), `SELECT status FROM `+table+` WHERE id = $1`, id).Scan(&s); err != nil {
		t.Fatal(err)
	}
	return s
}

// AC-HND-02 ส่วนบันทึกการส่งมอบ (T-10)
func TestAC_HND_02_HandoverRecordsBorrower(t *testing.T) {
	pool := testPool(t)
	h := routesOnly(pool)

	t.Run("ตรวจบัตรแล้ว บันทึกชื่อ นามสกุล รหัสนักศึกษา และเวลาที่ส่งมอบ", func(t *testing.T) {
		itemID := createItem(t, pool, "T02-HO-0001", "available")
		requestID := createRequest(t, pool, "student-001", itemID, "approved")
		before := time.Now().Add(-time.Second)

		rec := postHandover(t, h, tokenStaff, requestID, validForm(tokenStudentA), twoJPEGs())
		if rec.Code != http.StatusCreated {
			t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
		}
		var first, last, studentID, method, by string
		var at time.Time
		if err := pool.QueryRow(context.Background(),
			`SELECT borrower_first_name, borrower_last_name, borrower_student_id, handed_over_at, identity_method, handed_over_by
			 FROM handover WHERE request_id = $1`, requestID).Scan(&first, &last, &studentID, &at, &method, &by); err != nil {
			t.Fatalf("ไม่พบแถว handover: %v", err)
		}
		if first != "สมชาย" || last != "ใจดี" || studentID != "student-001" {
			t.Errorf("ผู้ยืม = %q %q %q", first, last, studentID)
		}
		if at.Before(before) || at.After(time.Now().Add(time.Second)) {
			t.Errorf("handed_over_at = %v ต้องเป็นเวลาที่ส่งมอบจริง", at)
		}
		if method != "auth_api" || by != "staff-001" {
			t.Errorf("identity_method = %q handed_over_by = %q", method, by)
		}
	})

	t.Run("ไม่ยืนยันการตรวจบัตร ไม่บันทึก", func(t *testing.T) {
		itemID := createItem(t, pool, "T02-HO-0002", "available")
		requestID := createRequest(t, pool, "student-001", itemID, "approved")
		form := validForm(tokenStudentA)
		form["cardChecked"] = "false"

		rec := postHandover(t, h, tokenStaff, requestID, form, twoJPEGs())
		var e ErrorDTO
		_ = json.Unmarshal(rec.Body.Bytes(), &e)
		if rec.Code != http.StatusUnprocessableEntity || e.Code != "VALIDATION_FAILED" || len(e.Fields) != 1 || e.Fields[0] != "cardChecked" {
			t.Fatalf("HTTP %d %+v ต้องเป็น 422 VALIDATION_FAILED ช่อง cardChecked", rec.Code, e)
		}
		if n := countLoans(t, pool, itemID); n != 0 {
			t.Errorf("มี Loan %d รายการ ต้องไม่มี", n)
		}
		if s := statusOf(t, pool, "equipment_item", itemID); s != "available" {
			t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น available", s)
		}
	})
}

// AC-HND-07 (FR-HND-07, DOM-REQ-01)
// Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว อุปกรณ์รายชิ้นมีสถานะ "ว่าง" และยืนยันตัวตนผู้ยืมแล้ว
// When  เจ้าหน้าที่บันทึกการส่งมอบ
// Then  มี Loan ใหม่ 1 รายการของผู้ยืมและอุปกรณ์ชิ้นนั้น ที่ยังไม่มีเวลาคืน
//
//	สถานะอุปกรณ์เป็น "ถูกยืม" และสถานะของคำขอเป็น "ส่งมอบแล้ว"
func TestAC_HND_07_HandoverCreatesLoan(t *testing.T) {
	pool := testPool(t)
	itemID := createItem(t, pool, "T07-HO-0001", "available")
	requestID := createRequest(t, pool, "student-001", itemID, "approved")

	rec := postHandover(t, routesOnly(pool), tokenStaff, requestID, validForm(tokenStudentA), twoJPEGs())
	if rec.Code != http.StatusCreated {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
	}
	var out HandoverOutput
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.LoanID == 0 || out.AssetCode != "T07-HO-0001" || out.HandedOverAt.IsZero() || out.DueAt.IsZero() {
		t.Errorf("response = %+v", out)
	}

	var n int
	var borrower string
	if err := pool.QueryRow(context.Background(),
		`SELECT count(*), max(borrower_id) FROM loan WHERE item_id = $1 AND returned_at IS NULL`, itemID).Scan(&n, &borrower); err != nil {
		t.Fatal(err)
	}
	if n != 1 || borrower != "student-001" {
		t.Errorf("Loan ที่ยังไม่คืน = %d รายการ ผู้ยืม %q ต้องเป็น 1 รายการของ student-001", n, borrower)
	}
	if s := statusOf(t, pool, "equipment_item", itemID); s != "borrowed" {
		t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น borrowed", s)
	}
	if s := statusOf(t, pool, "borrow_request", requestID); s != "handed_over" {
		t.Errorf("สถานะคำขอ = %q ต้องเป็น handed_over", s)
	}
}
