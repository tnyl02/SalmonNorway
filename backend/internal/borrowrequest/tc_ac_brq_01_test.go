package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-01-1 ถึง TC-BRQ-01-3 (AC-BRQ-01, FR-BRQ-01)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)

import (
	"context"
	"net/http"
	"testing"
)

// setItemStatus เปลี่ยนสถานะอุปกรณ์ จำลองว่ามีฟีเจอร์อื่น (005 ส่งมอบ) เปลี่ยนระหว่างผู้ขอกรอกฟอร์ม
func setItemStatus(t *testing.T, assetCode, status string) {
	t.Helper()
	if _, err := testPool(t).Exec(context.Background(),
		`UPDATE equipment_item SET status = $1 WHERE asset_code = $2`, status, assetCode); err != nil {
		t.Fatal(err)
	}
}

// itemNotAvailableBody คือรูปแบบคำตอบของกรณีอุปกรณ์ไม่ว่าง ซึ่งแนบสถานะล่าสุดมาด้วย (CON-TECH-03)
type itemNotAvailableBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// TC-BRQ-01-1 (ทางปกติ) อุปกรณ์ถูกเปลี่ยนเป็น borrowed ระหว่างกรอก
func TestTC_BRQ_01_1_ItemBecameBorrowed(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-01-1"

	// Given อุปกรณ์สถานะ available ตอนเปิดฟอร์ม แล้วถูกเปลี่ยนเป็น borrowed
	insertTempItem(t, assetCode)
	setItemStatus(t, assetCode, "borrowed")

	// When ส่ง POST ของชิ้นนั้น
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 409 ITEM_NOT_AVAILABLE
	if w.Code != http.StatusConflict {
		t.Fatalf("(1) อุปกรณ์ถูกยืมไปแล้วต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var body itemNotAvailableBody
	decodeJSON(t, w, &body)
	if body.Code != "ITEM_NOT_AVAILABLE" {
		t.Errorf("(1) รหัสต้องเป็น ITEM_NOT_AVAILABLE ตาม CON-TECH-03 แต่ได้ %q", body.Code)
	}
	// Then (2) body มีสถานะล่าสุด borrowed (FR-BRQ-01 ให้แจ้งสถานะล่าสุด)
	if body.Status != "borrowed" {
		t.Errorf("(2) ต้องแจ้งสถานะล่าสุดว่า borrowed แต่ได้ %q", body.Status)
	}
	// Then (3) ไม่มีแถวของชิ้นนี้เพิ่ม
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) อุปกรณ์ไม่ว่างต้องไม่สร้างคำขอ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-01-2 (ขอบ) สถานะอีกค่าหนึ่งที่ไม่ใช่ available ก็ต้องถูกปฏิเสธเหมือนกัน
func TestTC_BRQ_01_2_ItemUnavailable(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-01-2"

	// Given อุปกรณ์ถูกเปลี่ยนเป็น unavailable ระหว่างกรอก
	insertTempItem(t, assetCode)
	setItemStatus(t, assetCode, "unavailable")

	// When ส่ง POST ของชิ้นนั้น
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 409 ITEM_NOT_AVAILABLE เหมือนกัน เพราะ CON-DATA-01 ให้ยืมได้เฉพาะ available
	if w.Code != http.StatusConflict {
		t.Fatalf("(1) อุปกรณ์ไม่พร้อมใช้งานต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var body itemNotAvailableBody
	decodeJSON(t, w, &body)
	if body.Code != "ITEM_NOT_AVAILABLE" {
		t.Errorf("(1) รหัสต้องเป็น ITEM_NOT_AVAILABLE แต่ได้ %q", body.Code)
	}
	// Then (2) body มีสถานะล่าสุด unavailable
	if body.Status != "unavailable" {
		t.Errorf("(2) ต้องแจ้งสถานะล่าสุดว่า unavailable แต่ได้ %q", body.Status)
	}
	// Then (3) ไม่มีแถวเพิ่ม
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Errorf("(3) อุปกรณ์ไม่ว่างต้องไม่สร้างคำขอ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-01-3 (ทางผิด) ทำให้เงื่อนไขใน Given ไม่เป็นจริง คือสถานะยัง available ตอนกดส่ง
func TestTC_BRQ_01_3_StillAvailable(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-01-3"

	// Given สถานะยังเป็น available ตอนกดส่ง
	insertTempItem(t, assetCode)

	// When ส่ง POST ของชิ้นนั้น
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) อุปกรณ์ยังว่างต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	// Then (2) มีแถวใหม่ 1 แถว
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(2) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
	// Then (3) สถานะอุปกรณ์ยังเป็น available เพราะ 002 ไม่เปลี่ยนสถานะอุปกรณ์ (CON-DATA-01, DOM-REQ-01)
	var status string
	if err := pool.QueryRow(context.Background(),
		`SELECT status FROM equipment_item WHERE asset_code = $1`, assetCode).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != "available" {
		t.Errorf("(3) อุปกรณ์ที่มีคำขอค้างต้องยังเป็น available แต่เป็น %q", status)
	}
}
