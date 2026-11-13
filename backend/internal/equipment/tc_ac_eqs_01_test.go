package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-01 (ห้ามแก้ให้ผ่าน)

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// TC-EQS-01-1 (ปกติ)
func TestTC_EQS_01_1_ModelRowAndItems(t *testing.T) {
	// Given นักศึกษายืนยันตัวตนแล้ว ข้อมูลตั้งต้น รุ่น A 5 ชิ้น ว่าง 3
	srv := newTestServer(t)

	// When เรียกรายการรุ่น แล้วเรียกรายชิ้นของรุ่น A
	rec := get(t, srv, PathPrefix+"models", tokenStudentDept1)
	items := decode[ListDTO[ModelRowDTO]](t, rec).Items
	itemsRec := get(t, srv, PathPrefix+"models/1/items", tokenStudentDept1)

	// Then (1) รุ่น A แสดง 1 แถว
	var rows []ModelRowDTO
	for _, m := range items {
		if m.Name == "Notebook ตัวอย่างรุ่น A" {
			rows = append(rows, m)
		}
	}
	if len(rows) != 1 {
		t.Fatalf("(1) รุ่น A ต้องแสดง 1 แถว ได้ %d", len(rows))
	}
	// Then (2) ประเภท "คอมพิวเตอร์พกพา"
	if rows[0].Category != "คอมพิวเตอร์พกพา" {
		t.Errorf("(2) ประเภท = %q", rows[0].Category)
	}
	// Then (3) availableCount 3 totalCount 5
	if rows[0].AvailableCount != 3 || rows[0].TotalCount != 5 {
		t.Errorf("(3) ได้ %d / %d ต้องเป็น 3 / 5", rows[0].AvailableCount, rows[0].TotalCount)
	}
	// Then (4) รายชิ้น 5 ชิ้นตามรหัสครุภัณฑ์ CS-NB-0001 ถึง 0005 พร้อมสถานะของแต่ละชิ้น
	list := decode[ListDTO[ItemRow]](t, itemsRec).Items
	want := map[string]Status{"CS-NB-0001": StatusAvailable, "CS-NB-0002": StatusAvailable, "CS-NB-0003": StatusAvailable, "CS-NB-0004": StatusBorrowed, "CS-NB-0005": StatusBorrowed}
	if len(list) != 5 {
		t.Fatalf("(4) ต้องได้ 5 ชิ้น ได้ %d", len(list))
	}
	for _, it := range list {
		if want[it.AssetCode] != it.Status {
			t.Errorf("(4) %s สถานะ %q ต้องเป็น %q", it.AssetCode, it.Status, want[it.AssetCode])
		}
	}
	// Then (5) แบ่งหน้า: ยังไม่ตรวจ เพราะรอ Q-05 (จำนวนต่อหน้ายังไม่มีใครตัดสิน)
}

// TC-EQS-01-3 (ขอบ)
func TestTC_EQS_01_3_PendingRequestStillAvailable(t *testing.T) {
	// Given รุ่น A 5 ชิ้น ว่าง 3 และชิ้นที่ว่าง 1 ชิ้นมีคำขอยืมค้างอยู่ (สถานะคำขอ requested)
	srv := newTestServer(t)
	insertPendingRequest(t, "CS-NB-0001")

	// When เรียกรายการรุ่น
	m := modelsByID(t, srv, url.Values{})[1]

	// Then รุ่น A ยังเป็น availableCount 3 totalCount 5 คำขอที่ค้างไม่ลดจำนวนว่าง (DOM-REQ-01)
	if m.AvailableCount != 3 || m.TotalCount != 5 {
		t.Errorf("ได้ %d / %d ต้องเป็น 3 / 5 (DOM-REQ-01)", m.AvailableCount, m.TotalCount)
	}
}

// TC-EQS-01-4 (ขอบ)
func TestTC_EQS_01_4_ZeroAvailableStillListed(t *testing.T) {
	// Given รุ่นทดสอบ 5 ชิ้น ว่าง 0 (ถูกยืม 5)
	srv := newTestServer(t)
	id := insertTempModel(t, "รุ่นทดสอบ TC-01-4", StatusBorrowed, StatusBorrowed, StatusBorrowed, StatusBorrowed, StatusBorrowed)

	// When เรียกรายการรุ่น
	models := modelsByID(t, srv, url.Values{})

	// Then รุ่นนั้นยังแสดง 1 แถว availableCount 0 totalCount 5
	m, ok := models[id]
	if !ok {
		t.Fatal("รุ่นที่ว่าง 0 ชิ้นต้องยังแสดงเมื่อไม่มีตัวกรอง")
	}
	if m.AvailableCount != 0 || m.TotalCount != 5 {
		t.Errorf("ได้ %d / %d ต้องเป็น 0 / 5", m.AvailableCount, m.TotalCount)
	}
}

// TC-EQS-01-5 (ทางผิด)
func TestTC_EQS_01_5_Unauthenticated401(t *testing.T) {
	// Given ยังไม่ยืนยันตัวตน (ไม่มี token)
	srv := newTestServer(t)

	// When เรียกรายการรุ่น และรายชิ้นของรุ่น A
	for _, path := range []string{PathPrefix + "models", PathPrefix + "models/1/items"} {
		rec := get(t, srv, path, "")
		raw := rec.Body.String()

		// Then (1) 401 รหัส UNAUTHENTICATED พร้อม message
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("(1) %s status = %d ต้องเป็น 401", path, rec.Code)
		}
		e := errorBody(t, raw)
		if e["code"] != "UNAUTHENTICATED" || e["message"] == "" {
			t.Errorf("(1) %s body = %s", path, raw)
		}
		// Then (2) body ไม่มีข้อมูลรุ่นหรือรายชิ้น
		if strings.Contains(raw, "items") || strings.Contains(raw, "CS-NB") || strings.Contains(raw, "Notebook") {
			t.Errorf("(2) %s ส่งข้อมูลอุปกรณ์กลับไป: %s", path, raw)
		}
	}
}

// TC-EQS-01-6 (ทางผิด)
func TestTC_EQS_01_6_UnknownModel(t *testing.T) {
	// Given นักศึกษายืนยันตัวตนแล้ว
	srv := newTestServer(t)

	// When เลือกรุ่นที่ไม่มีอยู่ (รหัส 9999) และรหัสรุ่นที่ไม่ใช่ตัวเลข ("abc")
	missing := get(t, srv, PathPrefix+"models/9999/items", tokenStudentDept1)
	bad := get(t, srv, PathPrefix+"models/abc/items", tokenStudentDept1)

	// Then (1) รุ่นที่ไม่มีอยู่ได้ 200 รายการว่าง {items: [], total: 0}
	if missing.Code != http.StatusOK {
		t.Errorf("(1) status = %d ต้องเป็น 200", missing.Code)
	}
	if raw := missing.Body.String(); !strings.Contains(raw, `"items":[]`) || !strings.Contains(raw, `"total":0`) {
		t.Errorf("(1) body = %s ต้องเป็น {items: [], total: 0}", raw)
	}
	// Then (2) รหัสไม่ใช่ตัวเลขได้ 404 รหัส ITEM_NOT_FOUND พร้อม message
	if bad.Code != http.StatusNotFound {
		t.Errorf("(2) status = %d ต้องเป็น 404", bad.Code)
	}
	if e := errorBody(t, bad.Body.String()); e["code"] != "ITEM_NOT_FOUND" || e["message"] == "" {
		t.Errorf("(2) body = %v", e)
	}
}
