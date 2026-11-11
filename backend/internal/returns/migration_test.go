package returns

// test ของ migration 0009 และค่าใน config ที่ 006 ใช้ อยู่ใน package ของ 006 เอง ไม่แก้ test ของ 001
// (ย้ายมาจาก platform/migrate_test.go และ config/config_test.go 2569-11-11)

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/config"
)

// execRolledBack รันคำสั่ง SQL ใน transaction ที่ยกเลิกทิ้งเสมอ ข้อมูลจึงไม่ค้างข้าม test
func execRolledBack(t *testing.T, pool *pgxpool.Pool, sqls ...string) error {
	t.Helper()
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	for _, s := range sqls {
		if _, err := tx.Exec(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

// T-01 เสร็จเมื่อ: audit_log รับ action ใหม่ของ 006 ทั้ง 4 ค่า (DOM-AUDIT-01, FR-RET-07)
func TestMigration006_AuditLogReturnActions(t *testing.T) {
	pool := testPool(t)
	for _, a := range []string{"return_list", "return_lookup", "return_detail", "return"} {
		stmt := `INSERT INTO audit_log (actor_id, actor_role, action, asset_code, outcome) VALUES ('staff-001', 'staff', '` + a + `', 'CS-NB-0004', 'ok')`
		if err := execRolledBack(t, pool, stmt); err != nil {
			t.Errorf("action %q ของ 006 ต้องบันทึกได้ แต่ได้ error: %v", a, err)
		}
	}
	// ค่าเดิมของ 001 และ 003 ต้องยังใช้ได้หลัง 0009 แทนที่ CHECK
	for _, a := range []string{"list", "search", "detail", "cancel"} {
		stmt := `INSERT INTO audit_log (actor_id, actor_role, action, outcome) VALUES ('student-001', 'student', '` + a + `', 'ok')`
		if err := execRolledBack(t, pool, stmt); err != nil {
			t.Errorf("action %q เดิมต้องยังบันทึกได้: %v", a, err)
		}
	}
}

// T-01 เสร็จเมื่อ: loan มี return_condition (ปกติ/ชำรุด เท่านั้น) และ returned_late
// และคืนแล้วต้องมีครบทั้งเวลาคืน สภาพ เครื่องหมายคืนเกินกำหนด (FR-RET-02, FR-RET-06)
func TestMigration006_LoanReturnColumns(t *testing.T) {
	pool := testPool(t)
	insertBorrowedItem(t, pool, "T01-MIGRATE", 2*time.Hour)
	loan := `(SELECT l.id FROM loan l JOIN equipment_item i ON i.id = l.item_id WHERE i.asset_code = 'T01-MIGRATE')`

	for name, stmt := range map[string]string{
		"คืนสภาพปกติ":  `UPDATE loan SET returned_at = now(), return_condition = 'normal', returned_late = false WHERE id = ` + loan,
		"คืนสภาพชำรุด": `UPDATE loan SET returned_at = now(), return_condition = 'damaged', returned_late = true WHERE id = ` + loan,
	} {
		if err := execRolledBack(t, pool, stmt); err != nil {
			t.Errorf("%s ต้องบันทึกได้: %v", name, err)
		}
	}
	for name, stmt := range map[string]string{
		"สภาพสูญหาย (ไม่ใช่สภาพ FR-RET-02)":   `UPDATE loan SET returned_at = now(), return_condition = 'lost', returned_late = false WHERE id = ` + loan,
		"คืนแล้วแต่ไม่มีสภาพ":                 `UPDATE loan SET returned_at = now(), returned_late = false WHERE id = ` + loan,
		"คืนแล้วแต่ไม่มีเครื่องหมายเกินกำหนด": `UPDATE loan SET returned_at = now(), return_condition = 'normal' WHERE id = ` + loan,
		"มีสภาพแต่ไม่มีเวลาคืน":               `UPDATE loan SET return_condition = 'normal', returned_late = false WHERE id = ` + loan,
	} {
		if err := execRolledBack(t, pool, stmt); err == nil {
			t.Errorf("%s ต้องถูกปฏิเสธ", name)
		}
	}
}

// T-02 เสร็จเมื่อ: สภาพมี 2 ค่า และสถานะหลังคืนเป็นค่าใน CON-DATA-01 (FR-RET-02, FR-RET-03, FR-RET-04)
func TestConfig006_ReturnConditionStatus(t *testing.T) {
	want := map[string]string{"normal": "available", "damaged": "unavailable"}
	if fmt.Sprint(config.ReturnConditionStatus) != fmt.Sprint(want) {
		t.Fatalf("ReturnConditionStatus = %v, want %v", config.ReturnConditionStatus, want)
	}
	if _, ok := config.ReturnConditionStatus["lost"]; ok {
		t.Error("สูญหายไม่ใช่สภาพ (FR-RET-02)")
	}
}

// T-02: เขตเวลาต้องโหลดได้จริง (glossary.md)
func TestConfig006_TimezoneLoads(t *testing.T) {
	if _, err := time.LoadLocation(config.Timezone); err != nil {
		t.Fatalf("โหลดเขตเวลา %q ไม่ได้: %v", config.Timezone, err)
	}
}
