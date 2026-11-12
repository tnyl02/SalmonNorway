package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-09 (ห้ามแก้ให้ผ่าน)
// TC-EQS-09-7 ตรวจด้วยคน (รอ Q-SH-08) ไม่มีโค้ด

import (
	"net/http"
	"strings"
	"testing"
)

// TC-EQS-09-1 (ปกติ)
func TestTC_EQS_09_1_UnknownToken401(t *testing.T) {
	// Given token ที่ระบบไม่รู้จัก (session หมดอายุแบบจำลอง)
	srv := newTestServer(t)

	// When เรียกรายการอุปกรณ์
	rec := get(t, srv, PathPrefix+"models", "mock-expired-999")

	// Then (1) 401 code UNAUTHENTICATED พร้อม message (CON-TECH-03)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d ต้องเป็น 401", rec.Code)
	}
	raw := rec.Body.String()
	if e := errorBody(t, raw); e["code"] != "UNAUTHENTICATED" || e["message"] == "" {
		t.Errorf("body ต้องเป็น {code: UNAUTHENTICATED, message} ได้ %s", raw)
	}
	// Then (2) ไม่มี items ในคำตอบ (IF-IDP-01)
	if strings.Contains(raw, "items") {
		t.Errorf("คำตอบมี items: %s", raw)
	}
}

// TC-EQS-09-2 (ขอบ)
func TestTC_EQS_09_2_NoToken401(t *testing.T) {
	// Given ไม่ส่ง token เลย
	srv := newTestServer(t)

	// When เรียกรายการอุปกรณ์ และรายละเอียด CS-NB-0001
	for _, path := range []string{PathPrefix + "models", PathPrefix + "items/CS-NB-0001"} {
		rec := get(t, srv, path, "")

		// Then ทั้ง 2 endpoint ตอบ 401 UNAUTHENTICATED ไม่มีข้อมูลอุปกรณ์
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("%s status = %d ต้องเป็น 401", path, rec.Code)
			continue
		}
		raw := rec.Body.String()
		if e := errorBody(t, raw); e["code"] != "UNAUTHENTICATED" {
			t.Errorf("%s code = %q ต้องเป็น UNAUTHENTICATED", path, e["code"])
		}
		if strings.Contains(raw, "items") || strings.Contains(raw, "CS-NB") || strings.Contains(raw, "status") {
			t.Errorf("%s คำตอบมีข้อมูลอุปกรณ์: %s", path, raw)
		}
	}
}

// TC-EQS-09-6 (ทางผิด)
func TestTC_EQS_09_6_ValidToken200(t *testing.T) {
	// Given token ที่ยังใช้ได้ (mock-student-001)
	srv := newTestServer(t)

	// When เรียกรายการอุปกรณ์
	rec := get(t, srv, PathPrefix+"models", tokenStudentDept1)

	// Then 200 มีรายการ ไม่ถูกส่งไปหน้าเข้าสู่ระบบ (AC-EQS-01)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	if got := decode[ListDTO[ModelRowDTO]](t, rec); len(got.Items) == 0 {
		t.Error("ต้องมีรายการรุ่น")
	}
}
