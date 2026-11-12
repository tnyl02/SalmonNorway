package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-02 (ห้ามแก้ให้ผ่าน)

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// TC-EQS-02-1 (ปกติ)
func TestTC_EQS_02_1_FilterAvailable(t *testing.T) {
	// Given อยู่ในหน้ารายการอุปกรณ์ (ข้อมูลตั้งต้น)
	srv := newTestServer(t)

	// When กรองสถานะ "ว่าง" (ส่วนค้นหาคำว่า "Notebook" รอ Q-06)
	models := modelsByID(t, srv, url.Values{"status": {"available"}})
	items := decode[ListDTO[ItemRow]](t, get(t, srv, PathPrefix+"models/1/items?status=available", tokenStudentDept1)).Items

	// Then (1) ทุกแถวรุ่นนับเฉพาะชิ้นที่ว่าง รุ่น A ว่าง 3 / ทั้งหมด 3
	for id, m := range models {
		if m.AvailableCount != m.TotalCount {
			t.Errorf("(1) รุ่น %d นับชิ้นที่ไม่ว่างด้วย %d / %d", id, m.AvailableCount, m.TotalCount)
		}
	}
	if a := models[1]; a.AvailableCount != 3 || a.TotalCount != 3 {
		t.Errorf("(1) รุ่น A ได้ %d / %d ต้องเป็น 3 / 3", a.AvailableCount, a.TotalCount)
	}
	// Then (2) รายชิ้นของรุ่น A มี 3 ชิ้น ทุกชิ้นสถานะว่าง
	if len(items) != 3 {
		t.Errorf("(2) ต้องได้ 3 ชิ้น ได้ %d", len(items))
	}
	for _, it := range items {
		if it.Status != StatusAvailable {
			t.Errorf("(2) %s สถานะ %s", it.AssetCode, it.Status)
		}
	}
	// Then (3) ชื่อหรือรหัสตรงกับคำค้น: ยังไม่ตรวจ เพราะรอ Q-06 (วิธีจับคู่คำค้นยังไม่มีใครตัดสิน)
}

// TC-EQS-02-3 (ขอบ)
func TestTC_EQS_02_3_OnlyModelsWithMatch(t *testing.T) {
	// Given รุ่นทดสอบที่มีชิ้นว่าง 0 ชิ้น (ถูกยืม 2) และรุ่น B ที่มีชิ้นว่างพอดี 1 ชิ้น
	srv := newTestServer(t)
	id := insertTempModel(t, "รุ่นทดสอบ TC-02-3", StatusBorrowed, StatusBorrowed)

	// When กรองสถานะ "ว่าง"
	models := modelsByID(t, srv, url.Values{"status": {"available"}})

	// Then (1) รุ่นที่ไม่มีชิ้นว่างไม่แสดง
	if _, ok := models[id]; ok {
		t.Error("(1) รุ่นที่ไม่มีชิ้นว่างต้องไม่แสดง")
	}
	// Then (2) รุ่น B แสดงด้วยจำนวน 1
	if b, ok := models[2]; !ok || b.TotalCount != 1 {
		t.Errorf("(2) รุ่น B ต้องแสดงด้วยจำนวน 1 ได้ %+v (มี %v)", b, ok)
	}
}

// TC-EQS-02-4 (ทางผิด)
func TestTC_EQS_02_4_FilterUnauthenticated(t *testing.T) {
	// Given ไม่ได้อยู่ในระบบ (ไม่มี token)
	srv := newTestServer(t)

	// When กรองสถานะ "ว่าง"
	rec := get(t, srv, PathPrefix+"models?status=available", "")
	raw := rec.Body.String()

	// Then 401 UNAUTHENTICATED พร้อม message ไม่มีรายการใดถูกส่งกลับ
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d ต้องเป็น 401", rec.Code)
	}
	if e := errorBody(t, raw); e["code"] != "UNAUTHENTICATED" || e["message"] == "" {
		t.Errorf("body = %s", raw)
	}
	if strings.Contains(raw, "items") {
		t.Errorf("ส่งรายการกลับไป: %s", raw)
	}
}

// TC-EQS-02-5 (ทางผิด)
func TestTC_EQS_02_5_UnknownStatus(t *testing.T) {
	// Given อยู่ในหน้ารายการอุปกรณ์
	srv := newTestServer(t)

	// When กรองด้วยค่าสถานะที่ไม่มีใน CON-DATA-01 (เช่น "pending")
	rec := get(t, srv, PathPrefix+"models?status=pending", tokenStudentDept1)

	// Then 200 รายการว่าง {items: [], total: 0} ค่าที่ไม่รู้จักไม่ตรงกับชิ้นใด (ทีมตัดสิน 2569-11-12)
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d ต้องเป็น 200", rec.Code)
	}
	if raw := rec.Body.String(); !strings.Contains(raw, `"items":[]`) || !strings.Contains(raw, `"total":0`) {
		t.Errorf("body = %s ต้องเป็น {items: [], total: 0}", raw)
	}
}
