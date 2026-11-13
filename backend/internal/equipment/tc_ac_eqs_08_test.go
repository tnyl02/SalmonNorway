package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-08 (ห้ามแก้ให้ผ่าน)

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
	"salmonnorway/backend/internal/platform"
)

// TC-EQS-08-1 (ปกติ)
func TestTC_EQS_08_1_DBDown503(t *testing.T) {
	// Given ฐานข้อมูลปิดอยู่: pool ที่ชี้ฐานข้อมูลทดสอบแล้วถูกปิด
	closed, err := pgxpool.New(context.Background(), testPool(t).Config().ConnString())
	if err != nil {
		t.Fatal(err)
	}
	closed.Close()
	h := NewHandler(NewService(NewRepository(closed)))
	mux := http.NewServeMux()
	mux.Handle(PathPrefix, auth.Middleware(auth.NewMockProvider(), h.Routes()))
	srv := platform.NoStore(platform.MethodNotAllowedJSON(mux))

	// When เปิดหน้ารายการอุปกรณ์
	rec := get(t, srv, PathPrefix+"models", tokenStudentDept1)

	// Then (1) 503 code UPSTREAM_UNAVAILABLE พร้อม message ไม่ว่าง (CON-TECH-03)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d ต้องเป็น 503", rec.Code)
	}
	raw := rec.Body.String()
	if e := errorBody(t, raw); e["code"] != "UPSTREAM_UNAVAILABLE" || e["message"] == "" {
		t.Errorf("body ต้องเป็น {code: UPSTREAM_UNAVAILABLE, message} ได้ %s", raw)
	}
	// Then (2) ไม่มี items ในคำตอบ (FR-EQS-08 ไม่แสดงข้อมูลเก่า)
	if strings.Contains(raw, "items") {
		t.Errorf("คำตอบมี items: %s", raw)
	}
}

// TC-EQS-08-3 (ขอบ) ใช้เวลาจริงราว 5 วินาที
func TestTC_EQS_08_3_SlowDB5s(t *testing.T) {
	// Given ตาราง equipment_item ถูกล็อก ฐานข้อมูลค้าง
	pool := testPool(t)
	srv := newTestServer(t)
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "LOCK TABLE equipment_item IN ACCESS EXCLUSIVE MODE"); err != nil {
		t.Fatal(err)
	}

	// When เปิดหน้ารายการอุปกรณ์
	start := time.Now()
	rec := get(t, srv, PathPrefix+"models", tokenStudentDept1)
	elapsed := time.Since(start)

	// Then (1) 503 UPSTREAM_UNAVAILABLE
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d ต้องเป็น 503", rec.Code)
	}
	if e := errorBody(t, rec.Body.String()); e["code"] != "UPSTREAM_UNAVAILABLE" {
		t.Errorf("code = %q ต้องเป็น UPSTREAM_UNAVAILABLE", e["code"])
	}
	// Then (2) ตอบหลังรอครบ 5 วินาที (config.EquipmentDBTimeout, Q-08) และไม่เกิน 5 วินาทีบวกส่วนเผื่อ 2 วินาที (ทีมตัดสิน 2569-11-12)
	const margin = 2 * time.Second
	if elapsed < config.EquipmentDBTimeout {
		t.Errorf("ตอบหลัง %v เร็วกว่า %v", elapsed, config.EquipmentDBTimeout)
	}
	if elapsed > config.EquipmentDBTimeout+margin {
		t.Errorf("ตอบหลัง %v ค้างนานเกิน %v", elapsed, config.EquipmentDBTimeout+margin)
	}
}

// TC-EQS-08-6 (ทางผิด)
func TestTC_EQS_08_6_DBUpNormal(t *testing.T) {
	// Given ฐานข้อมูลตอบปกติ
	srv := newTestServer(t)

	// When เปิดหน้ารายการอุปกรณ์
	rec := get(t, srv, PathPrefix+"models", tokenStudentDept1)

	// Then 200 มีรายการรุ่นตามข้อมูลตั้งต้น ไม่มี code error (AC-EQS-01)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	raw := rec.Body.String()
	if strings.Contains(raw, `"code"`) {
		t.Errorf("คำตอบมี code error: %s", raw)
	}
	ids := map[int]bool{}
	for _, m := range decode[ListDTO[ModelRowDTO]](t, rec).Items {
		ids[m.ModelID] = true
	}
	if !ids[1] || !ids[2] {
		t.Errorf("ต้องมีรุ่น A (1) และรุ่น B (2) ตามข้อมูลตั้งต้น ได้ %v", ids)
	}
}
