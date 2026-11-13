package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-06 (ห้ามแก้ให้ผ่าน)

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// TC-EQS-06-1 (ปกติ)
func TestTC_EQS_06_1_UnavailableNoLoan(t *testing.T) {
	// Given CS-PJ-0002 สถานะ "ไม่พร้อมใช้งาน"
	srv := newTestServer(t)

	// When mock-student-001 เปิดรายละเอียด CS-PJ-0002
	rec := get(t, srv, PathPrefix+"items/CS-PJ-0002", tokenStudentDept1)

	// Then (1) 200 สถานะ "ไม่พร้อมใช้งาน"
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	raw := rec.Body.String()
	got := decode[StudentItemDTO](t, rec)
	if got.Status != StatusUnavailable {
		t.Errorf("สถานะ = %q ต้องเป็น unavailable", got.Status)
	}
	// Then (2) ส่งคำขอยืมไม่ได้ canRequestLoan false (FR-EQS-05)
	if !strings.Contains(raw, `"canRequestLoan":false`) {
		t.Errorf("ต้องมี canRequestLoan false: %s", raw)
	}
	// Then (3) ไม่มีกำหนดคืน (CR-002)
	if got.DueAt != nil || strings.Contains(raw, "dueAt") {
		t.Errorf("ต้องไม่มีกำหนดคืน: %s", raw)
	}
}

// TC-EQS-06-3 (ขอบ)
func TestTC_EQS_06_3_AllUnavailableZeroAvailable(t *testing.T) {
	// Given รุ่นทดสอบ 2 ชิ้น "ไม่พร้อมใช้งาน" ทั้ง 2 ชิ้น
	id := insertTempModel(t, "TC รุ่นไม่พร้อมใช้งานทั้งรุ่น", StatusUnavailable, StatusUnavailable)
	srv := newTestServer(t)

	// When เปิดหน้ารายการ
	rows := modelsByID(t, srv, url.Values{})

	// Then รุ่นนั้นยังแสดง 1 แถว ว่าง 0 / ทั้งหมด 2 (FR-EQS-01)
	m, ok := rows[id]
	if !ok {
		t.Fatalf("ไม่มีแถวของรุ่น %d", id)
	}
	if m.AvailableCount != 0 || m.TotalCount != 2 {
		t.Errorf("ได้ ว่าง %d / ทั้งหมด %d ต้องเป็น ว่าง 0 / ทั้งหมด 2", m.AvailableCount, m.TotalCount)
	}
}

// TC-EQS-06-5 (ทางผิด)
func TestTC_EQS_06_5_AvailableCanLoan(t *testing.T) {
	// Given CS-PJ-0001 สถานะ "ว่าง"
	srv := newTestServer(t)

	// When mock-student-001 เปิดรายละเอียด CS-PJ-0001
	rec := get(t, srv, PathPrefix+"items/CS-PJ-0001", tokenStudentDept1)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	got := decode[StudentItemDTO](t, rec)

	// Then (1) สถานะ "ว่าง"
	if got.Status != StatusAvailable {
		t.Errorf("สถานะ = %q ต้องเป็น available", got.Status)
	}
	// Then (2) ส่งคำขอยืมได้ canRequestLoan true (FR-EQS-06)
	if !got.CanRequestLoan {
		t.Error("canRequestLoan ต้องเป็น true")
	}
}
