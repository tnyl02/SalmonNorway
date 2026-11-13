package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-07-1 ถึง TC-BRQ-07-3 (AC-BRQ-07, FR-BRQ-07)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

// insertTempItemDept2 ใส่อุปกรณ์ของภาควิชา 2 (model 3 ตาม seed.sql) ไว้ใช้เฉพาะ test
// ต่างจาก insertTempItem ที่ใส่เป็น model 1 ซึ่งเป็นของภาควิชา 1
func insertTempItemDept2(t *testing.T, assetCode string) {
	t.Helper()
	pool := testPool(t)
	if _, err := pool.Exec(context.Background(),
		`INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 3, 'available')`, assetCode); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(),
			`DELETE FROM borrow_request WHERE item_id IN (SELECT id FROM equipment_item WHERE asset_code = $1)`, assetCode)
		_, _ = pool.Exec(context.Background(), `DELETE FROM equipment_item WHERE asset_code = $1`, assetCode)
	})
}

// expectNoItemInfo ยืนยันว่า body ไม่มีข้อมูลอุปกรณ์หลุดไปให้คนต่างภาควิชา (DOM-ACCESS-01, DOM-PRIV-01)
func expectNoItemInfo(t *testing.T, body, assetCode string) {
	t.Helper()
	for _, leak := range []string{"available", "borrowed", "unavailable", "Notebook", "Projector", assetCode} {
		if strings.Contains(body, leak) {
			t.Errorf("body ต้องไม่มีข้อมูลอุปกรณ์ แต่พบคำว่า %q", leak)
		}
	}
}

// TC-BRQ-07-1 (ทางปกติ) ผู้ขอภาควิชา 2 ขออุปกรณ์ของภาควิชา 1
func TestTC_BRQ_07_1_CrossDepartment(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-07-1" // อุปกรณ์ของภาควิชา 1 (model 1)

	// Given ผู้ขอสังกัดภาควิชา 2 อุปกรณ์เป้าหมายเป็นของภาควิชา 1
	insertTempItem(t, assetCode)

	// When ส่งคำขอยืมอุปกรณ์ชิ้นนั้น
	w := postRequest(t, srv, tokenStudentDept2, bodyFor(assetCode))

	// Then (1) ได้ 403 FORBIDDEN_DEPARTMENT
	if w.Code != http.StatusForbidden {
		t.Fatalf("(1) คำขอข้ามภาควิชาต้องได้ 403 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if code := errorCodeOf(t, w); code != "FORBIDDEN_DEPARTMENT" {
		t.Errorf("(1) รหัสต้องเป็น FORBIDDEN_DEPARTMENT ตาม CON-TECH-03 แต่ได้ %q", code)
	}
	// Then (2) ไม่มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-101", assetCode); n != 0 {
		t.Errorf("(2) คำขอข้ามภาควิชาต้องไม่ถูกสร้าง แต่พบ %d ใบ", n)
	}
	// Then (3) body ไม่มีข้อมูลอุปกรณ์ เช่น ชื่อรุ่นหรือสถานะ
	expectNoItemInfo(t, w.Body.String(), assetCode)
}

// TC-BRQ-07-2 (ทางผิด) ทำให้เงื่อนไข "ต่างภาควิชา" ไม่เป็นจริง
func TestTC_BRQ_07_2_SameDepartment(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-07-2" // อุปกรณ์ของภาควิชา 2 (model 3)

	// Given ผู้ขอและอุปกรณ์อยู่ภาควิชาเดียวกัน (ทั้งคู่เป็นภาควิชา 2)
	insertTempItemDept2(t, assetCode)

	// When ส่งคำขอยืม
	w := postRequest(t, srv, tokenStudentDept2, bodyFor(assetCode))

	// Then (1) ได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ขอยืมของภาควิชาตนเองต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	// Then (2) มีแถวใหม่
	if n := countRequestsOnItem(t, pool, "student-101", assetCode); n != 1 {
		t.Errorf("(2) ต้องมีคำขอ 1 ใบ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-07-3 (ขอบ) เข้าได้ทั้ง 403 และ 409 พร้อมกัน ต้องได้ 403 ก่อน
// ทีมปิด Q-BRQ-15 เมื่อ 2569-11-12 ว่าตรวจสิทธิ์ก่อนตรวจสถานะอุปกรณ์
func TestTC_BRQ_07_3_ForbiddenBeforeNotAvailable(t *testing.T) {
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-07-3"

	// Given ผู้ขอภาควิชา 2 ขออุปกรณ์ของภาควิชา 1 ที่สถานะเป็น borrowed ด้วย
	insertTempItem(t, assetCode)
	setItemStatus(t, assetCode, "borrowed")

	// When ส่งคำขอยืม
	w := postRequest(t, srv, tokenStudentDept2, bodyFor(assetCode))

	// Then (1) ได้ 403 FORBIDDEN_DEPARTMENT ไม่ใช่ 409
	if w.Code != http.StatusForbidden {
		t.Fatalf("(1) ต้องได้ 403 ตามลำดับการตรวจที่ทีมปิด Q-BRQ-15 ไว้ แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if code := errorCodeOf(t, w); code != "FORBIDDEN_DEPARTMENT" {
		t.Errorf("(1) รหัสต้องเป็น FORBIDDEN_DEPARTMENT ไม่ใช่ ITEM_NOT_AVAILABLE แต่ได้ %q", code)
	}
	// Then (2) body ไม่มีสถานะอุปกรณ์
	expectNoItemInfo(t, w.Body.String(), assetCode)
}
