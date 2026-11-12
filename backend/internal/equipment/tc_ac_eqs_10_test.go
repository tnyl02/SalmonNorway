package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-10 (ห้ามแก้ให้ผ่าน)
// TC-EQS-10-1 ถึง 10-4 ตรวจด้วยคน (k6, T-22 รอ Q-17) ไม่มีโค้ดในไฟล์นี้

import (
	"context"
	"net/http"
	"net/url"
	"sort"
	"testing"
	"time"
)

// perfP95 คือเวลาตอบสนอง p95 สูงสุดของการค้นหาและแสดงรายการ (NFR-PERF-01, AC-EQS-10)
// ควรย้ายไป backend/internal/config/config.go แต่ /testcases แก้โค้ดของระบบไม่ได้ จึงตั้งไว้ที่นี่ที่เดียว (เสนอทีมแล้ว)
const perfP95 = 3 * time.Second

// TC-EQS-10-5 (ทางผิด)
func TestTC_EQS_10_5_SingleUser1000Items(t *testing.T) {
	// Given ข้อมูลอุปกรณ์ 1,000 รายการ ผู้ใช้คนเดียว
	// Q-17 ยังไม่ตอบว่านับชิ้นหรือรุ่น จึงเตรียมให้ครบทั้งสองแบบ: 1,000 รุ่น รุ่นละ 1 ชิ้น = 1,000 รุ่นและ 1,000 ชิ้น
	// (ภาควิชา 1 ประเภท "คอมพิวเตอร์พกพา" สถานะสลับ ว่าง / ไม่พร้อมใช้งาน) ลบทิ้งหลัง test
	pool := testPool(t)
	ctx := context.Background()
	if _, err := pool.Exec(ctx, `
INSERT INTO equipment_model (name, category_id, department_id)
SELECT 'TC10-5 รุ่น ' || g, 1, 1 FROM generate_series(1, 1000) g`); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `
INSERT INTO equipment_item (asset_code, model_id, status)
SELECT 'TC10-5-' || lpad(m.id::text, 6, '0'), m.id, CASE WHEN m.id % 2 = 0 THEN 'available' ELSE 'unavailable' END
FROM equipment_model m WHERE m.name LIKE 'TC10-5 รุ่น %'`); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM equipment_item WHERE asset_code LIKE 'TC10-5-%'`)
		_, _ = pool.Exec(ctx, `DELETE FROM equipment_model WHERE name LIKE 'TC10-5 รุ่น %'`)
	})
	srv := newTestServer(t)
	path := PathPrefix + "models?" + url.Values{"category": {"คอมพิวเตอร์พกพา"}, "status": {"available"}}.Encode()

	// When ค้นหาด้วยตัวกรองประเภทและสถานะ 20 ครั้งติดกัน
	const runs = 20
	durations := make([]time.Duration, 0, runs)
	for i := 0; i < runs; i++ {
		start := time.Now()
		rec := get(t, srv, path, tokenStudentDept1)
		durations = append(durations, time.Since(start))
		if rec.Code != http.StatusOK {
			t.Fatalf("ครั้งที่ %d status = %d ต้องเป็น 200", i+1, rec.Code)
		}
	}

	// Then p95 ของเวลาตอบสนอง <= 3 วินาที (NFR-PERF-01)
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	p95 := durations[(runs*95+99)/100-1] // ลำดับที่ ceil(0.95 x 20) = 19
	t.Logf("p95 = %v (ช้าสุด %v) วัดบน Codespace ไม่ใช่ server จริง", p95, durations[runs-1])
	if p95 > perfP95 {
		t.Errorf("p95 = %v เกิน %v", p95, perfP95)
	}
}
