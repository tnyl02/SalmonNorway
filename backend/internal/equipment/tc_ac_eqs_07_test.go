package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-07 (ห้ามแก้ให้ผ่าน)

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// TC-EQS-07-1 (ปกติ)
func TestTC_EQS_07_1_StudentNoBorrower(t *testing.T) {
	// Given CS-NB-0004 ถูกยืมโดย student-001 (ข้อมูลตั้งต้น)
	srv := newTestServer(t)

	// When mock-student-002 เปิดรายละเอียด CS-NB-0004
	rec := get(t, srv, PathPrefix+"items/CS-NB-0004", "mock-student-002")

	// Then (1) 200 สถานะ "ถูกยืม"
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	raw := rec.Body.String()
	if got := decode[StudentItemDTO](t, rec); got.Status != StatusBorrowed {
		t.Errorf("สถานะ = %q ต้องเป็น borrowed", got.Status)
	}
	// Then (2) ไม่มีช่องผู้ยืม และไม่มี "student-001" ที่ใดในคำตอบ (FR-EQS-07, DOM-PRIV-01)
	if strings.Contains(strings.ToLower(raw), "borrower") {
		t.Errorf("คำตอบมีช่องผู้ยืม: %s", raw)
	}
	if strings.Contains(raw, "student-001") {
		t.Errorf("คำตอบมีรหัสผู้ยืม student-001: %s", raw)
	}
}

// TC-EQS-07-3 (ปกติ)
func TestTC_EQS_07_3_StaffSeesBorrower(t *testing.T) {
	// Given CS-NB-0004 ถูกยืมโดย student-001
	srv := newTestServer(t)

	// When mock-staff-001 เปิดรายละเอียด CS-NB-0004
	rec := get(t, srv, PathPrefix+"items/CS-NB-0004", "mock-staff-001")

	// Then (1) 200 สถานะ "ถูกยืม"
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	if got := decode[StaffItemDTO](t, rec); got.Status != StatusBorrowed {
		t.Errorf("สถานะ = %q ต้องเป็น borrowed", got.Status)
	}
	// Then (2) เห็นชื่อผู้ยืมปัจจุบัน: ยังไม่ตรวจ รอ Q-20 (ที่มาของชื่อผู้ยืม, T-36)
}

// TC-EQS-07-4 (ขอบ)
func TestTC_EQS_07_4_BorrowerSelfNoBorrowerField(t *testing.T) {
	// Given CS-NB-0004 ถูกยืมโดย student-001
	srv := newTestServer(t)

	// When student-001 (ผู้ยืมเอง) เปิดรายละเอียด CS-NB-0004
	rec := get(t, srv, PathPrefix+"items/CS-NB-0004", tokenStudentDept1)

	// Then ไม่มีช่องผู้ยืมในคำตอบ เหมือนนักศึกษาคนอื่น (FR-EQS-07 "นักศึกษาไม่เห็น")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	if raw := rec.Body.String(); strings.Contains(strings.ToLower(raw), "borrower") {
		t.Errorf("ผู้ยืมเองได้ช่องผู้ยืม: %s", raw)
	}
}

// TC-EQS-07-5 (ทางผิด)
func TestTC_EQS_07_5_AvailableNoBorrower(t *testing.T) {
	// Given CS-NB-0001 สถานะ "ว่าง" ไม่มีผู้ยืม
	srv := newTestServer(t)

	// When mock-staff-001 เปิดรายละเอียด CS-NB-0001
	rec := get(t, srv, PathPrefix+"items/CS-NB-0001", "mock-staff-001")

	// Then (1) 200 สถานะ "ว่าง" ไม่มีกำหนดคืน (CR-002 กำหนดคืนมีเฉพาะ "ถูกยืม")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	raw := rec.Body.String()
	got := decode[StaffItemDTO](t, rec)
	if got.Status != StatusAvailable {
		t.Errorf("สถานะ = %q ต้องเป็น available", got.Status)
	}
	if got.DueAt != nil || strings.Contains(raw, "dueAt") {
		t.Errorf("ชิ้นที่ว่างต้องไม่มีกำหนดคืน: %s", raw)
	}
	// Then (2) ชื่อผู้ยืมไม่แสดง: ยังไม่ตรวจ รอ Q-20
}

// TC-EQS-07-6 (ปกติ)
func TestTC_EQS_07_6_StaffSummary(t *testing.T) {
	// Given ข้อมูลตั้งต้น
	srv := newTestServer(t)

	// When mock-staff-001 เรียกสรุปจำนวนแยกตามสถานะ
	rec := get(t, srv, PathPrefix+"summary", "mock-staff-001")

	// Then (1) 200 มีช่อง ว่าง ถูกยืม ไม่พร้อมใช้งาน และรวม (FR-EQS-07)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200 (body: %s)", rec.Code, rec.Body.String())
	}
	var got map[string]int
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v", err)
	}
	for _, k := range []string{"available", "borrowed", "unavailable", "total"} {
		if _, ok := got[k]; !ok {
			t.Errorf("สรุปไม่มีช่อง %q: %v", k, got)
		}
	}
	// Then (2) รวม = ผลบวก 3 สถานะ
	if sum := got["available"] + got["borrowed"] + got["unavailable"]; sum != got["total"] {
		t.Errorf("ผลบวก 3 สถานะ %d ไม่เท่ากับรวม %d", sum, got["total"])
	}
	// Then (3) จำนวนตรงข้อมูลตั้งต้น ว่าง 5 ถูกยืม 2 ไม่พร้อมใช้งาน 1 รวม 8: ยังไม่ตรวจ รอ Q-P3 (นับทุกภาควิชาหรือไม่)
}

// TC-EQS-07-7 (ทางผิด)
func TestTC_EQS_07_7_StudentSummary403(t *testing.T) {
	// Given ข้อมูลตั้งต้น
	srv := newTestServer(t)

	// When mock-student-001 เรียกสรุปจำนวนแยกตามสถานะ
	rec := get(t, srv, PathPrefix+"summary", tokenStudentDept1)

	// Then (1) 403 code FORBIDDEN_ROLE พร้อม message (FR-EQS-07, CON-TECH-03)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d ต้องเป็น 403", rec.Code)
	}
	raw := rec.Body.String()
	if e := errorBody(t, raw); e["code"] != "FORBIDDEN_ROLE" || e["message"] == "" {
		t.Errorf("body ต้องเป็น {code: FORBIDDEN_ROLE, message} ได้ %s", raw)
	}
	// Then (2) ไม่มีตัวเลขสรุปในคำตอบ
	for _, k := range []string{"available", "borrowed", "unavailable", "total"} {
		if strings.Contains(raw, `"`+k+`"`) {
			t.Errorf("คำตอบมีช่องสรุป %q: %s", k, raw)
		}
	}
}
