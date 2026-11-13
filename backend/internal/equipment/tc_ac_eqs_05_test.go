package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-05 (ห้ามแก้ให้ผ่าน)

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// TC-EQS-05-1 (ปกติ)
func TestTC_EQS_05_1_ReturnedIsAvailable(t *testing.T) {
	// Given ชิ้นใหม่ของรุ่น A ถูกยืมอยู่ แล้วจำลอง UC-09 รับคืนสภาพปกติ
	const code = "TC05-1-0001"
	insertLoanItem(t, code, "normal", time.Now().Add(24*time.Hour))
	returnItem(t, code, StatusAvailable)
	srv := newTestServer(t)

	// When mock-student-001 เปิดรายละเอียดชิ้นนั้นทันที (ภายใน 30 วินาที)
	rec := get(t, srv, PathPrefix+"items/"+code, tokenStudentDept1)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	raw := rec.Body.String()
	got := decode[StudentItemDTO](t, rec)

	// Then (1) สถานะ "ว่าง"
	if got.Status != StatusAvailable {
		t.Errorf("สถานะ = %q ต้องเป็น available", got.Status)
	}
	// Then (2) ส่งคำขอยืมได้ canRequestLoan true (FR-EQS-06)
	if !got.CanRequestLoan {
		t.Error("canRequestLoan ต้องเป็น true")
	}
	// Then (3) ไม่มีกำหนดคืนค้าง (CR-002)
	if got.DueAt != nil || strings.Contains(raw, "dueAt") {
		t.Errorf("ต้องไม่มีกำหนดคืน: %s", raw)
	}
}

// TC-EQS-05-3 (ขอบ)
func TestTC_EQS_05_3_AvailableCountPlusOne(t *testing.T) {
	// Given ชิ้นใหม่ของรุ่น A ถูกยืม จดจำนวนว่างของรุ่น A ไว้
	const code = "TC05-3-0001"
	insertLoanItem(t, code, "normal", time.Now().Add(24*time.Hour))
	srv := newTestServer(t)
	before := modelsByID(t, srv, url.Values{})[1].AvailableCount

	// When จำลอง UC-09 รับคืน แล้วเปิดหน้ารายการ
	returnItem(t, code, StatusAvailable)
	after := modelsByID(t, srv, url.Values{})[1].AvailableCount

	// Then จำนวนว่างของรุ่น A เพิ่มขึ้น 1 (FR-EQS-01)
	if after != before+1 {
		t.Errorf("จำนวนว่าง %d -> %d ต้องเพิ่ม 1", before, after)
	}
}

// TC-EQS-05-4 (ขอบ)
func TestTC_EQS_05_4_NoStoreHeaders(t *testing.T) {
	// Given ข้อมูลตั้งต้น
	srv := newTestServer(t)

	// When เรียกหน้ารายการ รายชิ้น และรายละเอียด
	for _, p := range []string{PathPrefix + "models", PathPrefix + "models/1/items", PathPrefix + "items/CS-NB-0001"} {
		rec := get(t, srv, p, tokenStudentDept1)

		// Then ทุก endpoint ตอบ Cache-Control: no-store (NFR-DATA-01, plan ข้อ 6.2)
		if rec.Code != http.StatusOK {
			t.Errorf("%s status = %d", p, rec.Code)
		}
		if cc := rec.Header().Get("Cache-Control"); cc != "no-store" {
			t.Errorf("%s Cache-Control = %q ต้องเป็น no-store", p, cc)
		}
	}
}

// TC-EQS-05-6 (ทางผิด)
func TestTC_EQS_05_6_ReturnedDamagedUnavailable(t *testing.T) {
	// Given ชิ้นใหม่ของรุ่น A ถูกยืม แล้วรับคืนแบบชำรุด สถานะถูกเปลี่ยนเป็น "ไม่พร้อมใช้งาน" (ใครเปลี่ยน รอ Q-SH-09)
	const code = "TC05-6-0001"
	insertLoanItem(t, code, "normal", time.Now().Add(24*time.Hour))
	returnItem(t, code, StatusUnavailable)
	srv := newTestServer(t)

	// When mock-student-001 เปิดรายละเอียดชิ้นนั้น
	rec := get(t, srv, PathPrefix+"items/"+code, tokenStudentDept1)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	raw := rec.Body.String()
	got := decode[StudentItemDTO](t, rec)

	// Then (1) สถานะ "ไม่พร้อมใช้งาน"
	if got.Status != StatusUnavailable {
		t.Errorf("สถานะ = %q ต้องเป็น unavailable", got.Status)
	}
	// Then (2) ส่งคำขอยืมไม่ได้ canRequestLoan false (FR-EQS-05)
	if !strings.Contains(raw, `"canRequestLoan":false`) {
		t.Errorf("ต้องมี canRequestLoan false: %s", raw)
	}
}
