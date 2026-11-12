package equipment

// ตัวช่วยเตรียมข้อมูลของ test ที่มาจาก specs/001-equipment-status/test-cases.md (TestTC_EQS_*)
// ทุกข้อมูลที่เพิ่ม ลบทิ้งหลัง test เสมอ เพื่อไม่ให้ test อื่นที่ใช้ข้อมูลตั้งต้นเพี้ยน

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

// insertTempModel เพิ่มรุ่นทดสอบในภาควิชา 1 ประเภท "คอมพิวเตอร์พกพา" พร้อมชิ้นตามสถานะที่ให้ คืนรหัสรุ่น
func insertTempModel(t *testing.T, name string, statuses ...Status) int {
	t.Helper()
	pool := testPool(t)
	ctx := context.Background()
	var modelID int
	if err := pool.QueryRow(ctx,
		`INSERT INTO equipment_model (name, category_id, department_id) VALUES ($1, 1, 1) RETURNING id`, name).Scan(&modelID); err != nil {
		t.Fatal(err)
	}
	for i, s := range statuses {
		code := fmt.Sprintf("TC-%d-%02d", modelID, i+1)
		if _, err := pool.Exec(ctx, `INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, $2, $3)`, code, modelID, string(s)); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM equipment_item WHERE model_id = $1`, modelID)
		_, _ = pool.Exec(ctx, `DELETE FROM equipment_model WHERE id = $1`, modelID)
	})
	return modelID
}

// insertPendingRequest เพิ่มคำขอยืมสถานะ requested ของชิ้นที่ให้ (คำขอค้าง DOM-REQ-01) แล้วลบทิ้งหลัง test
func insertPendingRequest(t *testing.T, assetCode string) {
	t.Helper()
	pool := testPool(t)
	ctx := context.Background()
	var id int
	err := pool.QueryRow(ctx, `
INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
SELECT 'student-002', 'ผู้ขอทดสอบ', '660710000', id, 'normal', 'study', now(), now(), 'requested'
FROM equipment_item WHERE asset_code = $1
RETURNING id`, assetCode).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM borrow_request WHERE id = $1`, id) })
}

// errorBody อ่าน body ของ error ตาม CON-TECH-03 แล้วตรวจว่ามีแค่ code และ message ที่ไม่ว่าง
func errorBody(t *testing.T, raw string) map[string]string {
	t.Helper()
	var e map[string]string
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatalf("body ไม่ใช่ JSON ของ error: %v (%s)", err, raw)
	}
	return e
}
