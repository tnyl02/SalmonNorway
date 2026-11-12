package migrationtest

import (
	"context"
	"testing"

	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/platform"
)

const seed003File = "../../../db/seed/seed_003_cancel_request.sql"

// T-02 เสร็จเมื่อ: โหลด seed ของ 003 แล้วได้คำขอตามที่ tasks.md กำหนด และสถานะอุปกรณ์ไม่เปลี่ยน (ASM-01, ASM-10, CON-CXL-02, CON-DATA-01)
// โหลดใน test นี้ test เดียว ข้อมูลจึงคงอยู่ใน package (test อื่นใช้ transaction ที่ rollback จึงไม่ชนกัน)
func TestSeed003_RequestsAndStaffContact(t *testing.T) {
	ctx := context.Background()
	db := pool(t)
	url := dbURL(t)

	itemsBefore := itemStatuses(t)

	// โหลด 2 รอบ: รอบสองต้องไม่ error และไม่เพิ่มข้อมูลซ้ำ
	for i := 0; i < 2; i++ {
		if err := platform.Seed(ctx, url, seed003File); err != nil {
			t.Fatalf("โหลด seed ของ 003 รอบที่ %d ไม่สำเร็จ: %v", i+1, err)
		}
	}

	count := func(requester, status string) int {
		var n int
		if err := db.QueryRow(ctx,
			`SELECT count(*) FROM borrow_request WHERE requester_id = $1 AND status = $2`, requester, status).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	want := []struct {
		requester, status string
		n                 int
	}{
		{"student-001", "requested", 2},
		{"student-001", "approved_by_officer", 1},
		{"student-001", "pending_advisor_approval", 1},
		{"student-001", "approved", 1},
		{"student-001", "rejected", 1},
		{"student-001", "cancelled", 1},
		{"student-002", "requested", 1},
	}
	for _, w := range want {
		if got := count(w.requester, w.status); got != w.n {
			t.Errorf("%s สถานะ %s = %d ต้องเป็น %d", w.requester, w.status, got, w.n)
		}
	}

	// คำขอที่ยังยกเลิกได้ของ student-001 ต้องมี 4 รายการ (ตามจุดตรวจเฟส 1)
	var cancellable int
	if err := db.QueryRow(ctx,
		`SELECT count(*) FROM borrow_request WHERE requester_id = 'student-001' AND status = ANY($1)`,
		config.CancellableStatuses()).Scan(&cancellable); err != nil {
		t.Fatal(err)
	}
	if cancellable != 4 {
		t.Errorf("คำขอที่ยกเลิกได้ของ student-001 = %d ต้องเป็น 4", cancellable)
	}

	// requested รายการหนึ่งเป็น CS-NB-0001 และ student-002 ขอชิ้นเดียวกัน
	var shared int
	if err := db.QueryRow(ctx,
		`SELECT count(DISTINCT r.requester_id) FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id
		 WHERE i.asset_code = 'CS-NB-0001' AND r.status = 'requested'`).Scan(&shared); err != nil {
		t.Fatal(err)
	}
	if shared != 2 {
		t.Errorf("CS-NB-0001 ต้องมีคำขอ requested จาก 2 คน ได้ %d", shared)
	}

	// student-101 ไม่มีคำขอ
	var none int
	if err := db.QueryRow(ctx, `SELECT count(*) FROM borrow_request WHERE requester_id = 'student-101'`).Scan(&none); err != nil {
		t.Fatal(err)
	}
	if none != 0 {
		t.Errorf("student-101 ต้องไม่มีคำขอ ได้ %d", none)
	}

	// ภาควิชา 1 มี Email เจ้าหน้าที่จำลอง 1 รายการ
	var staff int
	if err := db.QueryRow(ctx, `SELECT count(*) FROM staff_contact WHERE department_id = 1`).Scan(&staff); err != nil {
		t.Fatal(err)
	}
	if staff != 1 {
		t.Errorf("staff_contact ภาควิชา 1 = %d ต้องเป็น 1", staff)
	}

	// CON-CXL-01, CON-DATA-01: seed ไม่เปลี่ยนสถานะอุปกรณ์
	itemsAfter := itemStatuses(t)
	for id, s := range itemsBefore {
		if itemsAfter[id] != s {
			t.Errorf("สถานะอุปกรณ์ id %d เปลี่ยนจาก %s เป็น %s", id, s, itemsAfter[id])
		}
	}
}

func itemStatuses(t *testing.T) map[int]string {
	t.Helper()
	rows, err := pool(t).Query(context.Background(), `SELECT id, status FROM equipment_item`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	m := map[int]string{}
	for rows.Next() {
		var id int
		var s string
		if err := rows.Scan(&id, &s); err != nil {
			t.Fatal(err)
		}
		m[id] = s
	}
	return m
}
