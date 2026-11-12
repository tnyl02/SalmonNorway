package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-16 (ห้ามแก้ให้ผ่าน)

import (
	"net/url"
	"testing"
)

// TC-EQS-16-1 (ปกติ)
func TestTC_EQS_16_1_TwoStatuses(t *testing.T) {
	// Given ข้อมูลตั้งต้นมีครบ 3 สถานะ (CS-PJ-0002 ไม่พร้อมใช้งาน)
	srv := newTestServer(t)
	q := url.Values{"status": {"available", "borrowed"}}

	// When เลือกตัวกรองสถานะ "ว่าง" และ "ถูกยืม" พร้อมกัน
	models := modelsByID(t, srv, q)
	items := decode[ListDTO[ItemRow]](t, get(t, srv, PathPrefix+"models/2/items?"+q.Encode(), tokenStudentDept1)).Items

	// Then (1) รุ่น A นับ 5 ชิ้น
	if models[1].TotalCount != 5 {
		t.Errorf("(1) รุ่น A ได้ %d ต้องเป็น 5", models[1].TotalCount)
	}
	// Then (2) รุ่น B นับ 1 ชิ้น
	if models[2].TotalCount != 1 {
		t.Errorf("(2) รุ่น B ได้ %d ต้องเป็น 1", models[2].TotalCount)
	}
	// Then (3) รายชิ้นของรุ่น B มีแค่ CS-PJ-0001 ไม่มี CS-PJ-0002
	if len(items) != 1 || items[0].AssetCode != "CS-PJ-0001" {
		t.Errorf("(3) รายชิ้นรุ่น B = %+v", items)
	}
}

// TC-EQS-16-3 (ขอบ)
func TestTC_EQS_16_3_AllStatusesEqualsNoFilter(t *testing.T) {
	// Given ข้อมูลตั้งต้นมีครบ 3 สถานะ
	srv := newTestServer(t)

	// When เลือกครบทั้ง 3 สถานะพร้อมกัน
	models := modelsByID(t, srv, url.Values{"status": {"available", "borrowed", "unavailable"}})

	// Then ผลเท่ากับไม่กรอง: รุ่น A ว่าง 3 / ทั้งหมด 5 และรุ่น B ทั้งหมด 2 รวม CS-PJ-0002
	if a := models[1]; a.AvailableCount != 3 || a.TotalCount != 5 {
		t.Errorf("รุ่น A ได้ %d / %d ต้องเป็น 3 / 5", a.AvailableCount, a.TotalCount)
	}
	if models[2].TotalCount != 2 {
		t.Errorf("รุ่น B ได้ %d ต้องเป็น 2", models[2].TotalCount)
	}
}

// TC-EQS-16-4 (ปกติ)
func TestTC_EQS_16_4_TwoCategories(t *testing.T) {
	// Given ข้อมูลตั้งต้น ประเภท "คอมพิวเตอร์พกพา" และ "อุปกรณ์นำเสนอ"
	srv := newTestServer(t)

	// When เลือกตัวกรองประเภททั้ง 2 ค่าพร้อมกัน
	models := modelsByID(t, srv, url.Values{"category": {"คอมพิวเตอร์พกพา", "อุปกรณ์นำเสนอ"}})

	// Then แสดงทั้งรุ่น A และรุ่น B
	if _, ok := models[1]; !ok {
		t.Error("ต้องแสดงรุ่น A")
	}
	if _, ok := models[2]; !ok {
		t.Error("ต้องแสดงรุ่น B")
	}
}

// TC-EQS-16-5 (ทางผิด)
func TestTC_EQS_16_5_UnknownValueMixed(t *testing.T) {
	// Given ข้อมูลตั้งต้นมีครบ 3 สถานะ
	srv := newTestServer(t)

	// When เลือก "ว่าง" คู่กับค่าที่ไม่มีใน CON-DATA-01 (เช่น "pending")
	mixed := modelsByID(t, srv, url.Values{"status": {"available", "pending"}})
	only := modelsByID(t, srv, url.Values{"status": {"available"}})

	// Then ผลเท่ากับกรอง "ว่าง" อย่างเดียว ค่าที่ไม่รู้จักไม่ตรงกับชิ้นใด (ทีมตัดสิน 2569-11-12)
	if len(mixed) != len(only) {
		t.Fatalf("ได้ %d รุ่น ต้องเท่ากับกรองว่างอย่างเดียว %d รุ่น", len(mixed), len(only))
	}
	for id, m := range only {
		if mixed[id] != m {
			t.Errorf("รุ่น %d ได้ %+v ต้องเป็น %+v", id, mixed[id], m)
		}
	}
}
