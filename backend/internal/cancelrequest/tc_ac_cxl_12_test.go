package cancelrequest

// test cases ของ AC-CXL-12 จาก specs/003-cancel-request/test-cases.md (สถานะใช้ได้) ห้ามแก้ให้ผ่าน
// TC-CXL-12-2 และ TC-CXL-12-5 เป็น test หน้าจอ อยู่ที่ frontend/src/__tests__/

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// tcItemsRaw อ่านค่า items แบบดิบ เพื่อแยก [] ออกจาก null
func tcItemsRaw(t *testing.T, body []byte) json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v (body: %s)", err, body)
	}
	raw, ok := m["items"]
	if !ok {
		t.Fatalf("body ต้องมี items: %s", body)
	}
	return raw
}

// TC-CXL-12-1 (ปกติ)
func TestTC_CXL_12_1_EmptyListIsArray(t *testing.T) {
	// Given A ยืนยันตัวตนแล้ว ไม่มีคำขอใดเลย
	srv := tcServer(t)

	// When A เรียก GET /api/v1/my/pending-requests
	rec := tcGet(t, srv, tcPathMine, tcTokenA)

	// Then (1) ได้ 200
	tcWantStatus(t, rec, http.StatusOK)
	// Then (2) items เป็นอาร์เรย์ว่าง [] ไม่ใช่ null
	if raw := string(tcItemsRaw(t, rec.Body.Bytes())); raw != "[]" {
		t.Errorf("items = %s ต้องเป็น []", raw)
	}
}

// TC-CXL-12-3 (ขอบ)
func TestTC_CXL_12_3_OnlyNonCancellableIsEmpty(t *testing.T) {
	// Given A มีคำขอแต่ทุกคำขอยกเลิกไม่ได้ (approved, rejected, handed_over, returned, cancelled อย่างละ 1)
	srv := tcServer(t)
	for _, s := range []string{"approved", "rejected", "handed_over", "returned", "cancelled"} {
		tcInsert(t, tcUserA, tcItemAvailable, "normal", s, time.Now())
	}

	// When A เรียก GET /api/v1/my/pending-requests
	rec := tcGet(t, srv, tcPathMine, tcTokenA)

	// Then (1) ได้ items: []
	tcWantStatus(t, rec, http.StatusOK)
	if raw := string(tcItemsRaw(t, rec.Body.Bytes())); raw != "[]" {
		t.Errorf("items = %s ต้องเป็น []", raw)
	}
}

// TC-CXL-12-4 (ขอบ)
func TestTC_CXL_12_4_OthersRequestsNotShown(t *testing.T) {
	// Given A ไม่มีคำขอ แต่ B มีคำขอ requested 2 คำขอ
	srv := tcServer(t)
	tcInsert(t, tcUserB, tcItemAvailable, "normal", "requested", time.Now())
	tcInsert(t, tcUserB, tcItemOther, "normal", "requested", time.Now())

	// When A เรียก GET /api/v1/my/pending-requests
	rec := tcGet(t, srv, tcPathMine, tcTokenA)

	// Then (1) ได้ items: [] ไม่เห็นคำขอของ B
	tcWantStatus(t, rec, http.StatusOK)
	if raw := string(tcItemsRaw(t, rec.Body.Bytes())); raw != "[]" {
		t.Errorf("items = %s ต้องเป็น [] (ไม่เห็นคำขอของ B)", raw)
	}
}
