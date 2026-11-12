package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-15 (ห้ามแก้ให้ผ่าน)

import (
	"net/url"
	"testing"
)

// TC-EQS-15-1 (ปกติ)
func TestTC_EQS_15_1_CountsFollowFilter(t *testing.T) {
	// Given รุ่น A 5 ชิ้น ว่าง 3 ถูกยืม 2
	srv := newTestServer(t)

	// When กรองสถานะ "ถูกยืม"
	models := modelsByID(t, srv, url.Values{"status": {"borrowed"}})

	// Then (1) รุ่น A แสดง 1 แถว
	a, ok := models[1]
	if !ok {
		t.Fatal("(1) รุ่น A ต้องแสดง")
	}
	// Then (2) ตัวเลขนับเฉพาะชิ้นที่ตรงตัวกรอง totalCount 2 (availableCount 0)
	if a.TotalCount != 2 || a.AvailableCount != 0 {
		t.Errorf("(2) ได้ %d / %d ต้องเป็น 0 / 2", a.AvailableCount, a.TotalCount)
	}
}

// TC-EQS-15-3 (ขอบ)
func TestTC_EQS_15_3_ZeroMatchHidden(t *testing.T) {
	// Given รุ่น A 5 ชิ้น ว่าง 3 ถูกยืม 2 ไม่พร้อมใช้งาน 0
	srv := newTestServer(t)

	// When กรองสถานะ "ไม่พร้อมใช้งาน"
	models := modelsByID(t, srv, url.Values{"status": {"unavailable"}})

	// Then รุ่น A ไม่แสดง เพราะมีชิ้นที่ตรง 0 ชิ้น
	if _, ok := models[1]; ok {
		t.Error("รุ่น A ต้องไม่แสดงเมื่อไม่มีชิ้นที่ตรง")
	}
}

// TC-EQS-15-4 (ทางผิด)
func TestTC_EQS_15_4_NoStatusFilterWholeModel(t *testing.T) {
	// Given รุ่น A 5 ชิ้น ว่าง 3 ถูกยืม 2
	srv := newTestServer(t)

	// When ไม่กรองสถานะ (กรองแค่ประเภท "คอมพิวเตอร์พกพา")
	models := modelsByID(t, srv, url.Values{"category": {"คอมพิวเตอร์พกพา"}})

	// Then ตัวเลขกลับเป็นทั้งรุ่น availableCount 3 totalCount 5
	if a := models[1]; a.AvailableCount != 3 || a.TotalCount != 5 {
		t.Errorf("ได้ %d / %d ต้องเป็น 3 / 5", a.AvailableCount, a.TotalCount)
	}
}
