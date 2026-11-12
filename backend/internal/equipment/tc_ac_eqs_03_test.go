package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-03 (ห้ามแก้ให้ผ่าน)

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// TC-EQS-03-1 (ปกติ)
func TestTC_EQS_03_1_EmptyIs200(t *testing.T) {
	// Given กรองด้วยเงื่อนไขที่ไม่มีในระบบ: ประเภท "อุปกรณ์นำเสนอ" + สถานะ "ถูกยืม" (แทนคำค้นที่ไม่มี รอ Q-06)
	srv := newTestServer(t)
	q := url.Values{"category": {"อุปกรณ์นำเสนอ"}, "status": {"borrowed"}}

	// When เรียกรายการรุ่น
	rec := get(t, srv, PathPrefix+"models?"+q.Encode(), tokenStudentDept1)

	// Then 200 {items: [], total: 0} ไม่ใช่ error
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d ต้องเป็น 200", rec.Code)
	}
	if raw := rec.Body.String(); !strings.Contains(raw, `"items":[]`) || !strings.Contains(raw, `"total":0`) || strings.Contains(raw, `"code"`) {
		t.Errorf("body = %s ต้องเป็น {items: [], total: 0} ไม่ใช่ error", raw)
	}
}
