package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-04 (ห้ามแก้ให้ผ่าน)

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TC-EQS-04-1 (ปกติ)
func TestTC_EQS_04_1_BorrowedShowsDueAt(t *testing.T) {
	// Given จำลอง UC-08 ส่งมอบชิ้นใหม่ของรุ่น A ยืมปกติ กำหนดคืน 16:00 ของพรุ่งนี้ (ยังไม่เลยกำหนด)
	const code = "TC04-1-0001"
	now := time.Now().In(bangkok(t, 2026, time.January, 1, 0, 0).Location())
	due := time.Date(now.Year(), now.Month(), now.Day()+1, 16, 0, 0, 0, now.Location())
	insertLoanItem(t, code, "normal", due)
	srv := newTestServer(t)

	// When mock-student-002 เปิดรายละเอียดชิ้นนั้น
	rec := get(t, srv, PathPrefix+"items/"+code, "mock-student-002")

	// Then (1) 200 สถานะ "ถูกยืม"
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	raw := rec.Body.String()
	got := decode[StudentItemDTO](t, rec)
	if got.Status != StatusBorrowed {
		t.Errorf("สถานะ = %q ต้องเป็น borrowed", got.Status)
	}
	// Then (2) มีกำหนดคืนเป็นวันที่และเวลาตรงกับที่บันทึก
	if got.DueAt == nil || !got.DueAt.Equal(due) {
		t.Errorf("กำหนดคืน = %v ต้องเป็น %s", got.DueAt, due)
	}
	// Then (3) ส่งคำขอยืมไม่ได้ canRequestLoan false (FR-EQS-04)
	if !strings.Contains(raw, `"canRequestLoan":false`) {
		t.Errorf("ต้องมี canRequestLoan false: %s", raw)
	}
	// Then (4) ไม่มีป้ายเลยกำหนดคืน
	if got.IsOverdue {
		t.Error("ยังไม่เลยกำหนด ต้องไม่มีป้ายเลยกำหนดคืน")
	}
}

// TC-EQS-04-3 (ขอบ)
func TestTC_EQS_04_3_OneMinuteBeforeDue(t *testing.T) {
	// Given ยืมปกติ กำหนดคืน 2569-10-10 16:00 เวลาปัจจุบัน 2569-10-10 15:59
	const code = "TC04-3-0001"
	due := bangkok(t, 2026, time.October, 10, 16, 0)
	insertLoanItem(t, code, "normal", due)
	svc := serviceAt(t, bangkok(t, 2026, time.October, 10, 15, 59))

	// When นักศึกษาเปิดรายละเอียด
	got, found, err := svc.GetItem(studentCtx, code)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	student := got.(StudentItemDTO)

	// Then (1) เห็นกำหนดคืน 2569-10-10 16:00
	if student.DueAt == nil || !student.DueAt.Equal(due) {
		t.Errorf("กำหนดคืน = %v ต้องเป็น %s", student.DueAt, due)
	}
	// Then (2) ไม่มีป้ายเลยกำหนดคืน (FR-EQS-04 ก)
	if student.IsOverdue {
		t.Error("15:59 ยังไม่เลยกำหนด ต้องไม่มีป้าย")
	}
}

// TC-EQS-04-5 (ทางผิด)
func TestTC_EQS_04_5_NotHandedOverStillAvailable(t *testing.T) {
	// Given CS-NB-0001 ยังไม่ส่งมอบ มีคำขอยืมค้าง 1 รายการ
	insertPendingRequest(t, "CS-NB-0001")
	srv := newTestServer(t)

	// When mock-student-001 เปิดรายละเอียด CS-NB-0001
	rec := get(t, srv, PathPrefix+"items/CS-NB-0001", tokenStudentDept1)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	raw := rec.Body.String()
	got := decode[StudentItemDTO](t, rec)

	// Then (1) สถานะ "ว่าง" (DOM-REQ-01)
	if got.Status != StatusAvailable {
		t.Errorf("สถานะ = %q ต้องเป็น available", got.Status)
	}
	// Then (2) ไม่มีกำหนดคืน (CR-002)
	if got.DueAt != nil || strings.Contains(raw, "dueAt") {
		t.Errorf("ต้องไม่มีกำหนดคืน: %s", raw)
	}
	// Then (3) ส่งคำขอยืมได้ canRequestLoan true (FR-EQS-06)
	if !got.CanRequestLoan {
		t.Error("canRequestLoan ต้องเป็น true")
	}
}

// TC-EQS-04-6 (ทางผิด)
func TestTC_EQS_04_6_ReadOnly405(t *testing.T) {
	// Given CS-NB-0004 สถานะ "ถูกยืม"
	pool := testPool(t)
	srv := newTestServer(t)

	// When mock-student-001 เรียก POST, PUT, DELETE ที่รายละเอียด CS-NB-0004
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, PathPrefix+"items/CS-NB-0004", strings.NewReader(`{"status":"available"}`))
		req.Header.Set("Authorization", "Bearer "+tokenStudentDept1)
		rec := httptest.NewRecorder()
		srv.ServeHTTP(rec, req)

		// Then (1) ทุก method ตอบ 405 METHOD_NOT_ALLOWED (CON-DATA-02, T-40)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s status = %d ต้องเป็น 405", method, rec.Code)
			continue
		}
		if e := errorBody(t, rec.Body.String()); e["code"] != "METHOD_NOT_ALLOWED" {
			t.Errorf("%s code = %q ต้องเป็น METHOD_NOT_ALLOWED", method, e["code"])
		}
	}

	// Then (2) สถานะในฐานข้อมูลยังเป็น "ถูกยืม"
	var status string
	if err := pool.QueryRow(context.Background(), `SELECT status FROM equipment_item WHERE asset_code = 'CS-NB-0004'`).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(StatusBorrowed) {
		t.Errorf("สถานะ = %q ต้องยังเป็น borrowed", status)
	}
}
