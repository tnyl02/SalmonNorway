package handover

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
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

// testPhotoRoot คือโฟลเดอร์เก็บรูปของ test (ไม่เขียนลง /data/photos ของเครื่อง)
var testPhotoRoot = filepath.Join(os.TempDir(), "salmon-handover-test-"+strconv.Itoa(os.Getpid()))

// routesOnly คือ router ที่ผ่าน auth.Middleware เหมือนใน main.go เก็บรูปที่ testPhotoRoot
func routesOnly(db DB) http.Handler {
	return routesAt(db, time.Time{})
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
			`DELETE FROM loan WHERE item_id = $1`,
		} {
			if _, err := pool.Exec(ctx, q, id); err != nil {
				t.Errorf("ลบข้อมูลทดสอบไม่สำเร็จ: %v", err)
			}
		}
		// audit_log เป็น append-only (trigger ใน 0002) ลบแถวไม่ได้ และอ้าง borrow_request ด้วย foreign key
		// คำขอที่มี audit log แล้วจึงลบไม่ได้ ตั้งเป็น cancelled แทน ไม่ให้ค้างในรายการรอส่งมอบของ test อื่น (AC-HND-14)
		var audited bool
		if err := pool.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM audit_log WHERE request_id IN (SELECT id FROM borrow_request WHERE item_id = $1))`,
			id).Scan(&audited); err != nil {
			t.Errorf("ตรวจ audit_log ไม่สำเร็จ: %v", err)
			return
		}
		if audited {
			if _, err := pool.Exec(ctx, `UPDATE borrow_request SET status = 'cancelled' WHERE item_id = $1`, id); err != nil {
				t.Errorf("ปิดคำขอทดสอบไม่สำเร็จ: %v", err)
			}
			return
		}
		for _, q := range []string{
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

// bangkok คือเวลา Asia/Bangkok ของวันที่และเวลาที่กำหนด (2569 = ค.ศ. 2026)
func bangkok(t *testing.T, year int, month time.Month, day, hour, minute int) time.Time {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		t.Fatal(err)
	}
	return time.Date(year, month, day, hour, minute, 0, 0, loc)
}

// routesAt คือ router ที่นาฬิกาหยุดอยู่ที่ at (FR-HND-05) at ว่าง = ใช้เวลาจริง
func routesAt(db DB, at time.Time) http.Handler {
	provider := auth.NewMockProvider()
	now := time.Now
	if !at.IsZero() {
		now = func() time.Time { return at }
	}
	h := newHandler(db, NewMockBorrowerVerifier(provider), NewPhotoStore(testPhotoRoot), now)
	return auth.Middleware(provider, h.Routes())
}

// AC-HND-05 (FR-HND-05, DOM-LOAN-01)
// Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว
// When  เจ้าหน้าที่บันทึกการส่งมอบสำเร็จเมื่อ 2569-11-12 เวลา 10:00 น. (Asia/Bangkok)
// Then  การยืมนั้นมีเวลาที่ส่งมอบ 2569-11-12 10:00 น. และกำหนดคืน 2569-11-12 16:00 น.
func TestAC_HND_05_DueAt1600SameDay(t *testing.T) {
	pool := testPool(t)
	itemID := createItem(t, pool, "T05-HO-0001", "available")
	requestID := createRequest(t, pool, "student-001", itemID, "approved")
	at := bangkok(t, 2026, time.November, 12, 10, 0)
	wantDue := bangkok(t, 2026, time.November, 12, 16, 0)

	rec := postHandover(t, routesAt(pool, at), tokenStaff, requestID, validForm(tokenStudentA), twoJPEGs())
	if rec.Code != http.StatusCreated {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
	}
	var handedOverAt, dueAt time.Time
	if err := pool.QueryRow(context.Background(),
		`SELECT h.handed_over_at, l.due_at FROM handover h JOIN loan l ON l.id = h.loan_id WHERE h.request_id = $1`,
		requestID).Scan(&handedOverAt, &dueAt); err != nil {
		t.Fatal(err)
	}
	if !handedOverAt.Equal(at) {
		t.Errorf("เวลาที่ส่งมอบ = %v ต้องเป็น %v", handedOverAt.In(at.Location()), at)
	}
	if !dueAt.Equal(wantDue) {
		t.Errorf("กำหนดคืน = %v ต้องเป็น %v", dueAt.In(at.Location()), wantDue)
	}
}

// AC-HND-01 (FR-HND-01)
// Given คำขอยืมแบบปกติที่ยังไม่ได้รับอนุมัติ และอุปกรณ์รายชิ้นในคำขอมีสถานะ "ว่าง"
// When  เจ้าหน้าที่บันทึกการส่งมอบของคำขอนั้น
// Then  ระบบไม่บันทึกการส่งมอบ ไม่มี Loan ถูกสร้าง และอุปกรณ์ยังมีสถานะ "ว่าง"
// เพิ่มกรณีคำขอที่ส่งมอบไปแล้ว (plan.md ข้อ 4 รหัสเดียวกัน)
func TestAC_HND_01_RequestNotApproved(t *testing.T) {
	pool := testPool(t)
	h := routesOnly(pool)
	for i, status := range []string{"requested", "handed_over"} {
		t.Run(status, func(t *testing.T) {
			itemID := createItem(t, pool, "T01-NA-000"+strconv.Itoa(i+1), "available")
			requestID := createRequest(t, pool, "student-001", itemID, status)

			rec := postHandover(t, h, tokenStaff, requestID, validForm(tokenStudentA), twoJPEGs())
			var e ErrorDTO
			_ = json.Unmarshal(rec.Body.Bytes(), &e)
			if rec.Code != http.StatusConflict || e.Code != "REQUEST_NOT_APPROVED" {
				t.Fatalf("HTTP %d %+v ต้องเป็น 409 REQUEST_NOT_APPROVED", rec.Code, e)
			}
			if n := countLoans(t, pool, itemID); n != 0 {
				t.Errorf("มี Loan %d รายการ ต้องไม่มี", n)
			}
			if s := statusOf(t, pool, "equipment_item", itemID); s != "available" {
				t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น available", s)
			}
			if s := statusOf(t, pool, "borrow_request", requestID); s != status {
				t.Errorf("สถานะคำขอ = %q ต้องเป็น %q", s, status)
			}
		})
	}
}

// AC-HND-04 (FR-HND-04)
// Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว แต่อุปกรณ์รายชิ้นในคำขอมีสถานะ "ไม่พร้อมใช้งาน" หรือ "ถูกยืม"
// When  เจ้าหน้าที่บันทึกการส่งมอบ
// Then  ระบบไม่สร้างการยืม และสถานะอุปกรณ์ไม่เปลี่ยน
func TestAC_HND_04_ItemNotAvailable(t *testing.T) {
	pool := testPool(t)
	h := routesOnly(pool)
	for i, status := range []string{"unavailable", "borrowed"} {
		t.Run(status, func(t *testing.T) {
			itemID := createItem(t, pool, "T04-NA-000"+strconv.Itoa(i+1), status)
			requestID := createRequest(t, pool, "student-001", itemID, "approved")

			rec := postHandover(t, h, tokenStaff, requestID, validForm(tokenStudentA), twoJPEGs())
			var e ErrorDTO
			_ = json.Unmarshal(rec.Body.Bytes(), &e)
			if rec.Code != http.StatusConflict || e.Code != "ITEM_NOT_AVAILABLE" {
				t.Fatalf("HTTP %d %+v ต้องเป็น 409 ITEM_NOT_AVAILABLE", rec.Code, e)
			}
			if n := countLoans(t, pool, itemID); n != 0 {
				t.Errorf("มี Loan %d รายการ ต้องไม่มี", n)
			}
			if s := statusOf(t, pool, "equipment_item", itemID); s != status {
				t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น %q", s, status)
			}
		})
	}
}

// AC-HND-10 (FR-HND-02)
// Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้วของนักศึกษา A
// When  ผู้ที่ยืนยันตัวตนตอนรับอุปกรณ์คือนักศึกษา B
// Then  ระบบปฏิเสธการส่งมอบ ไม่มี Loan ถูกสร้าง และอุปกรณ์ยังมีสถานะ "ว่าง"
func TestAC_HND_10_BorrowerMismatch(t *testing.T) {
	pool := testPool(t)
	itemID := createItem(t, pool, "T10-MM-0001", "available")
	requestID := createRequest(t, pool, "student-001", itemID, "approved") // นักศึกษา A

	rec := postHandover(t, routesOnly(pool), tokenStaff, requestID, validForm(tokenStudentB), twoJPEGs()) // นักศึกษา B ล็อกอิน
	var e ErrorDTO
	_ = json.Unmarshal(rec.Body.Bytes(), &e)
	if rec.Code != http.StatusConflict || e.Code != "BORROWER_MISMATCH" {
		t.Fatalf("HTTP %d %+v ต้องเป็น 409 BORROWER_MISMATCH", rec.Code, e)
	}
	if n := countLoans(t, pool, itemID); n != 0 {
		t.Errorf("มี Loan %d รายการ ต้องไม่มี", n)
	}
	if s := statusOf(t, pool, "equipment_item", itemID); s != "available" {
		t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น available", s)
	}
	if s := statusOf(t, pool, "borrow_request", requestID); s != "approved" {
		t.Errorf("สถานะคำขอ = %q ต้องเป็น approved", s)
	}
}

// AC-HND-06 (FR-HND-06)
// Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว
// When  เจ้าหน้าที่แนบรูปสภาพอุปกรณ์ JPG 2 รูป ขนาดรูปละ 2 MB แล้วบันทึกการส่งมอบ
// Then  ระบบบันทึกการส่งมอบ และรูปทั้ง 2 รูปถูกบันทึกผูกกับการยืมนั้น
func TestAC_HND_06_TwoJpgPhotosSaved(t *testing.T) {
	pool := testPool(t)
	itemID := createItem(t, pool, "T06-PH-0001", "available")
	requestID := createRequest(t, pool, "student-001", itemID, "approved")
	photos := []photoFile{
		{"a.jpg", "image/jpeg", testJPEG(2 << 20)},
		{"b.jpg", "image/jpeg", testJPEG(2 << 20)},
	}

	rec := postHandover(t, routesOnly(pool), tokenStaff, requestID, validForm(tokenStudentA), photos)
	if rec.Code != http.StatusCreated {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
	}
	var out HandoverOutput
	_ = json.Unmarshal(rec.Body.Bytes(), &out)

	rows, err := pool.Query(context.Background(),
		`SELECT file_path, content_type, size_bytes FROM handover_photo WHERE loan_id = $1 ORDER BY id`, out.LoanID)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	n := 0
	for rows.Next() {
		var path, ct string
		var size int64
		if err := rows.Scan(&path, &ct, &size); err != nil {
			t.Fatal(err)
		}
		n++
		if ct != "image/jpeg" || size != 2<<20 {
			t.Errorf("รูปที่ %d: ชนิด %q ขนาด %d ต้องเป็น image/jpeg 2 MB", n, ct, size)
		}
		if filepath.Dir(path) != filepath.Join(testPhotoRoot, strconv.Itoa(out.LoanID)) {
			t.Errorf("รูปที่ %d อยู่ที่ %q ไม่ใช่โฟลเดอร์ของการยืมนี้", n, path)
		}
		if st, err := os.Stat(path); err != nil || st.Size() != size {
			t.Errorf("ไฟล์รูปที่ %d บนดิสก์: %v", n, err)
		}
	}
	if n != 2 {
		t.Errorf("มีรูป %d รูปผูกกับการยืม ต้องเป็น 2", n)
	}
}

// dirEntries นับรายการในโฟลเดอร์รูปของ test (ไม่มีโฟลเดอร์ = 0)
func dirEntries(t *testing.T) int {
	t.Helper()
	entries, err := os.ReadDir(testPhotoRoot)
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return len(entries)
}

// AC-HND-12 (FR-HND-06) กรณี (ก) 1 รูป (ข) GIF (T-21) และ (ค) 11 MB (T-28, 10 MB ต่อรูป ตาม Q-HND-10)
// Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว
// When  เจ้าหน้าที่บันทึกการส่งมอบโดยแนบรูป (ก) 1 รูป หรือ (ข) ไฟล์ GIF หรือ (ค) รูปขนาด 11 MB
// Then  ระบบไม่บันทึกการส่งมอบในทุกกรณี และไม่มี Loan ถูกสร้าง (และไม่มีไฟล์ค้างบนดิสก์ T-21)
func TestAC_HND_12_InvalidPhotos(t *testing.T) {
	pool := testPool(t)
	h := routesOnly(pool)
	gif := []byte("GIF89a\x01\x00\x01\x00\x00\x00\x00")

	for i, tc := range []struct {
		name   string
		photos []photoFile
	}{
		{"ก_1_รูป", []photoFile{{"1.jpg", "image/jpeg", testJPEG(1024)}}},
		{"ข_GIF", []photoFile{{"1.jpg", "image/jpeg", testJPEG(1024)}, {"2.gif", "image/gif", gif}}},
		{"ข_GIF_ตั้งชื่อเป็น_jpg", []photoFile{{"1.jpg", "image/jpeg", testJPEG(1024)}, {"2.jpg", "image/jpeg", gif}}},
		{"ค_11_MB", []photoFile{{"1.jpg", "image/jpeg", testJPEG(1024)}, {"2.jpg", "image/jpeg", testJPEG(11 << 20)}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			itemID := createItem(t, pool, "T12-PH-000"+strconv.Itoa(i+1), "available")
			requestID := createRequest(t, pool, "student-001", itemID, "approved")
			before := dirEntries(t)

			rec := postHandover(t, h, tokenStaff, requestID, validForm(tokenStudentA), tc.photos)
			var e ErrorDTO
			_ = json.Unmarshal(rec.Body.Bytes(), &e)
			if rec.Code != http.StatusUnprocessableEntity || e.Code != "VALIDATION_FAILED" || len(e.Fields) != 1 || e.Fields[0] != "photos" {
				t.Fatalf("HTTP %d %+v ต้องเป็น 422 VALIDATION_FAILED ช่อง photos", rec.Code, e)
			}
			if n := countLoans(t, pool, itemID); n != 0 {
				t.Errorf("มี Loan %d รายการ ต้องไม่มี", n)
			}
			if s := statusOf(t, pool, "equipment_item", itemID); s != "available" {
				t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น available", s)
			}
			if after := dirEntries(t); after != before {
				t.Errorf("โฟลเดอร์รูปมี %d รายการ (ก่อน %d) มีไฟล์ค้าง", after, before)
			}
		})
	}
}

// AC-HND-13 (FR-HND-07, DOM-REQ-01)
// Given อุปกรณ์รายชิ้นหนึ่งมีคำขอยืม 2 รายการ คำขอ X ได้รับอนุมัติแล้ว และคำขอ Y ยังไม่ได้ตัดสิน
// When  เจ้าหน้าที่บันทึกการส่งมอบของคำขอ X สำเร็จ
// Then  สถานะของคำขอ Y ไม่เปลี่ยน
//
//	และเมื่อเจ้าหน้าที่บันทึกการส่งมอบของคำขอ Y ระบบไม่สร้างการยืม เพราะอุปกรณ์ "ถูกยืม" (FR-HND-04)
//
// หมายเหตุ: Given นี้ไม่เกิดในระบบจริงหลัง Q-SH-06 ปิดบางส่วน (004 ปฏิเสธคำขออื่นอัตโนมัติ) ตั้งข้อมูลเองตาม AC จนกว่าทีมจะแก้ AC
// ขั้น "ส่งมอบ Y" ตั้ง Y เป็น approved ในฐานข้อมูลเอง เพราะ 005 ส่งมอบได้เฉพาะคำขอที่อนุมัติแล้ว (FR-HND-01)
func TestAC_HND_13_OtherRequestUntouched(t *testing.T) {
	pool := testPool(t)
	h := routesOnly(pool)
	itemID := createItem(t, pool, "T13-OR-0001", "available")
	requestX := createRequest(t, pool, "student-001", itemID, "approved")
	requestY := createRequest(t, pool, "student-002", itemID, "requested")

	if rec := postHandover(t, h, tokenStaff, requestX, validForm(tokenStudentA), twoJPEGs()); rec.Code != http.StatusCreated {
		t.Fatalf("ส่งมอบ X: HTTP %d: %s", rec.Code, rec.Body.String())
	}
	if s := statusOf(t, pool, "borrow_request", requestY); s != "requested" {
		t.Fatalf("สถานะคำขอ Y = %q ต้องไม่เปลี่ยน (requested)", s)
	}

	if _, err := pool.Exec(context.Background(), `UPDATE borrow_request SET status = 'approved' WHERE id = $1`, requestY); err != nil {
		t.Fatal(err)
	}
	rec := postHandover(t, h, tokenStaff, requestY, validForm("mock-student-002"), twoJPEGs())
	var e ErrorDTO
	_ = json.Unmarshal(rec.Body.Bytes(), &e)
	if rec.Code != http.StatusConflict || e.Code != "ITEM_NOT_AVAILABLE" {
		t.Fatalf("ส่งมอบ Y: HTTP %d %+v ต้องเป็น 409 ITEM_NOT_AVAILABLE", rec.Code, e)
	}
	if n := countLoans(t, pool, itemID); n != 1 {
		t.Errorf("มี Loan %d รายการ ต้องเป็น 1 (ของ X เท่านั้น)", n)
	}
}

// AC-HND-09 (CON-HND-01)
// Given ผู้ใช้บทบาทนักศึกษายืนยันตัวตนแล้ว
// When  เรียกบันทึกการส่งมอบของคำขอที่ได้รับอนุมัติแล้วโดยตรง
// Then  ระบบปฏิเสธ ไม่มี Loan ถูกสร้าง และสถานะอุปกรณ์ไม่เปลี่ยน
func TestAC_HND_09_StudentForbidden(t *testing.T) {
	pool := testPool(t)
	itemID := createItem(t, pool, "T09-ST-0001", "available")
	requestID := createRequest(t, pool, "student-001", itemID, "approved")

	rec := postHandover(t, routesOnly(pool), tokenStudentA, requestID, validForm(tokenStudentA), twoJPEGs())
	var e ErrorDTO
	_ = json.Unmarshal(rec.Body.Bytes(), &e)
	if rec.Code != http.StatusForbidden || e.Code != "FORBIDDEN_ROLE" {
		t.Fatalf("HTTP %d %+v ต้องเป็น 403 FORBIDDEN_ROLE", rec.Code, e)
	}
	if n := countLoans(t, pool, itemID); n != 0 {
		t.Errorf("มี Loan %d รายการ ต้องไม่มี", n)
	}
	if s := statusOf(t, pool, "equipment_item", itemID); s != "available" {
		t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น available", s)
	}
	if s := statusOf(t, pool, "borrow_request", requestID); s != "approved" {
		t.Errorf("สถานะคำขอ = %q ต้องเป็น approved", s)
	}
}

// AC-HND-11 (FR-HND-05, DOM-LOAN-01)
// Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว
// When  เจ้าหน้าที่บันทึกการส่งมอบเมื่อ 2569-11-12 เวลา 16:30 น. (Asia/Bangkok)
// Then  ระบบบันทึกการส่งมอบ การยืมนั้นมีเวลาที่ส่งมอบ 2569-11-12 16:30 น. และอุปกรณ์มีสถานะ "ถูกยืม"
//
//	และกำหนดคืน 2569-11-13 16:00 น. (Q-HND-03 ปิดแล้ว, T-25)
func TestAC_HND_11_HandoverAfter1600(t *testing.T) {
	pool := testPool(t)
	itemID := createItem(t, pool, "T11-AF-0001", "available")
	requestID := createRequest(t, pool, "student-001", itemID, "approved")
	at := bangkok(t, 2026, time.November, 12, 16, 30)

	rec := postHandover(t, routesAt(pool, at), tokenStaff, requestID, validForm(tokenStudentA), twoJPEGs())
	if rec.Code != http.StatusCreated {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
	}
	var handedOverAt, dueAt time.Time
	if err := pool.QueryRow(context.Background(),
		`SELECT h.handed_over_at, l.due_at FROM handover h JOIN loan l ON l.id = h.loan_id WHERE h.request_id = $1`,
		requestID).Scan(&handedOverAt, &dueAt); err != nil {
		t.Fatal(err)
	}
	if !handedOverAt.Equal(at) {
		t.Errorf("เวลาที่ส่งมอบ = %v ต้องเป็น %v", handedOverAt.In(at.Location()), at)
	}
	if want := bangkok(t, 2026, time.November, 13, 16, 0); !dueAt.Equal(want) {
		t.Errorf("กำหนดคืน = %v ต้องเป็น %v", dueAt.In(at.Location()), want)
	}
	if s := statusOf(t, pool, "equipment_item", itemID); s != "borrowed" {
		t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น borrowed", s)
	}
}

// auditRow คือแถว audit_log ของการส่งมอบ 1 แถว (DM-07)
type auditRow struct {
	actorID, actorRole, assetCode, outcome string
	occurredAt                             time.Time
}

// handoverAudits คืนแถว audit_log action handover ของคำขอ requestID เรียงตามเวลา (FR-HND-08)
func handoverAudits(t *testing.T, pool *pgxpool.Pool, requestID int) []auditRow {
	t.Helper()
	rows, err := pool.Query(context.Background(),
		`SELECT actor_id, actor_role, COALESCE(asset_code, ''), outcome, occurred_at
		 FROM audit_log WHERE action = 'handover' AND request_id = $1 ORDER BY id`, requestID)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []auditRow
	for rows.Next() {
		var r auditRow
		if err := rows.Scan(&r.actorID, &r.actorRole, &r.assetCode, &r.outcome, &r.occurredAt); err != nil {
			t.Fatal(err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

// AC-HND-08 (FR-HND-08, DOM-AUDIT-01)
// Given เจ้าหน้าที่บันทึกการส่งมอบ (ทั้งสำเร็จและไม่สำเร็จ Q-HND-14 ปิดแล้ว)
// When  การบันทึกเสร็จสิ้น
// Then  มี audit log ที่ระบุผู้ทำ เวลา รหัสครุภัณฑ์ รหัสคำขอยืม และผลลัพธ์
// subtest สำเร็จ (T-26) กรณีไม่สำเร็จอยู่ใน T-27
func TestAC_HND_08_AuditLogged(t *testing.T) {
	pool := testPool(t)
	h := routesOnly(pool)

	t.Run("สำเร็จ", func(t *testing.T) {
		itemID := createItem(t, pool, "T08-OK-0001", "available")
		requestID := createRequest(t, pool, "student-001", itemID, "approved")
		before := time.Now().Add(-time.Second)

		rec := postHandover(t, h, tokenStaff, requestID, validForm(tokenStudentA), twoJPEGs())
		if rec.Code != http.StatusCreated {
			t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
		}
		got := handoverAudits(t, pool, requestID)
		if len(got) != 1 {
			t.Fatalf("audit log ของการส่งมอบ = %d แถว ต้องเป็น 1: %+v", len(got), got)
		}
		r := got[0]
		if r.actorID != "staff-001" || r.actorRole != "staff" || r.assetCode != "T08-OK-0001" || r.outcome != "ok" {
			t.Errorf("audit log = %+v ต้องเป็นของ staff-001 อุปกรณ์ T08-OK-0001 ผล ok", r)
		}
		if r.occurredAt.Before(before) || r.occurredAt.After(time.Now().Add(time.Second)) {
			t.Errorf("เวลาใน audit log = %v ต้องเป็นเวลาที่บันทึกจริง", r.occurredAt)
		}
	})
}

// hangingVerifier จำลอง Auth API ที่ไม่ตอบ: ค้างจนกว่า ctx จะหมดเวลา (AC-HND-03)
type hangingVerifier struct{}

func (hangingVerifier) VerifyBorrower(ctx context.Context, _ string) (string, error) {
	<-ctx.Done()
	return "", ctx.Err()
}

// testAuthTimeout คือ BorrowerAuthTimeout แบบสั้นใน test แทน 10 วินาทีจริง (plan.md ข้อ 6)
const testAuthTimeout = 50 * time.Millisecond

// routesWithVerifier คือ router ที่ใช้ verifier ที่กำหนด ห่อด้วย WithTimeout เหมือน NewHandler
func routesWithVerifier(db DB, v BorrowerVerifier) http.Handler {
	h := newHandler(db, WithTimeout(v, testAuthTimeout), NewPhotoStore(testPhotoRoot), time.Now)
	return auth.Middleware(auth.NewMockProvider(), h.Routes())
}

// AC-HND-03 (FR-HND-02, FR-HND-03)
// Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว และ Auth API ไม่ตอบภายใน 10 วินาที
// When  เจ้าหน้าที่ยืนยันว่าตรวจบัตรแล้ว และบันทึกการส่งมอบ
// Then  ระบบบันทึกการส่งมอบได้โดยใช้ผลการตรวจบัตร
//
//	ถ้าเจ้าหน้าที่ไม่ยืนยันการตรวจบัตร ระบบไม่บันทึกการส่งมอบ
//
// รหัสผู้ยืมของกรณี card_only = requester_id ของคำขอ (ทีมตัดสิน 2569-11-12 ตอน /implement T-13)
func TestAC_HND_03_AuthApiTimeoutCardOnly(t *testing.T) {
	pool := testPool(t)
	h := routesWithVerifier(pool, hangingVerifier{})

	t.Run("verify_borrower ได้ unavailable", func(t *testing.T) {
		requestID := createRequest(t, pool, "student-001", createItem(t, pool, "T03-VB-0001", "available"), "approved")
		start := time.Now()
		rec := doJSON(h, http.MethodPost, verifyPath(requestID), tokenStaff, `{"borrowerCredential":"`+tokenStudentA+`"}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
		}
		var got VerifyBorrowerResult
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got.Result != "unavailable" {
			t.Errorf("ผล = %+v ต้องเป็น unavailable", got)
		}
		if waited := time.Since(start); waited > time.Second {
			t.Errorf("รอ %v ต้องหยุดรอที่ timeout", waited)
		}
	})

	t.Run("ตรวจบัตรแล้ว บันทึกด้วย card_only", func(t *testing.T) {
		itemID := createItem(t, pool, "T03-HO-0001", "available")
		requestID := createRequest(t, pool, "student-001", itemID, "approved")

		rec := postHandover(t, h, tokenStaff, requestID, validForm(tokenStudentA), twoJPEGs())
		if rec.Code != http.StatusCreated {
			t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
		}
		var method, studentID, borrowerID string
		if err := pool.QueryRow(context.Background(),
			`SELECT h.identity_method, h.borrower_student_id, l.borrower_id
			 FROM handover h JOIN loan l ON l.id = h.loan_id WHERE h.request_id = $1`, requestID).Scan(&method, &studentID, &borrowerID); err != nil {
			t.Fatalf("ไม่พบการส่งมอบ: %v", err)
		}
		if method != "card_only" || studentID != "student-001" || borrowerID != "student-001" {
			t.Errorf("identity_method = %q รหัสนักศึกษา = %q ผู้ยืม = %q ต้องเป็น card_only ของ student-001", method, studentID, borrowerID)
		}
		if s := statusOf(t, pool, "equipment_item", itemID); s != "borrowed" {
			t.Errorf("สถานะอุปกรณ์ = %q ต้องเป็น borrowed", s)
		}
	})

	t.Run("ไม่ยืนยันการตรวจบัตร ไม่บันทึก", func(t *testing.T) {
		itemID := createItem(t, pool, "T03-HO-0002", "available")
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

	t.Run("Auth API ตอบทันเวลา ยังใช้ auth_api", func(t *testing.T) {
		itemID := createItem(t, pool, "T03-HO-0003", "available")
		requestID := createRequest(t, pool, "student-001", itemID, "approved")
		fast := routesWithVerifier(pool, NewMockBorrowerVerifier(auth.NewMockProvider()))

		rec := postHandover(t, fast, tokenStaff, requestID, validForm(tokenStudentA), twoJPEGs())
		if rec.Code != http.StatusCreated {
			t.Fatalf("HTTP %d: %s", rec.Code, rec.Body.String())
		}
		var method string
		if err := pool.QueryRow(context.Background(),
			`SELECT identity_method FROM handover WHERE request_id = $1`, requestID).Scan(&method); err != nil {
			t.Fatal(err)
		}
		if method != "auth_api" {
			t.Errorf("identity_method = %q ต้องเป็น auth_api", method)
		}
	})
}
