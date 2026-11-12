package cancelrequest

// test cases ของ AC-CXL-08 จาก specs/003-cancel-request/test-cases.md (สถานะใช้ได้) ห้ามแก้ให้ผ่าน
// TC-CXL-08-2 เป็น test หน้าจอ อยู่ที่ frontend/src/__tests__/TC-CXL-08-2.test.jsx

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TC-CXL-08-1 (ปกติ)
func TestTC_CXL_08_1_ListOwnPendingSorted(t *testing.T) {
	// Given A ยืนยันตัวตนแล้ว มีคำขอ requested 2 คำขอ (ส่งต่างเวลากัน) approved 1 rejected 1 cancelled 1 และ B มีคำขอ requested 1
	srv := tcServer(t)
	base := time.Date(2026, 11, 10, 2, 0, 0, 0, time.UTC)
	older := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", base)
	newer := tcInsert(t, tcUserA, tcItemOther, "normal", "requested", base.Add(time.Hour))
	approved := tcInsert(t, tcUserA, tcItemAvailable2, "normal", "approved", base) // คนละชิ้นกับ requested (test-cases.md แก้ 2569-11-11)
	rejected := tcInsert(t, tcUserA, tcItemAvailable, "normal", "rejected", base)
	cancelled := tcInsert(t, tcUserA, tcItemAvailable, "normal", "cancelled", base)
	ofB := tcInsert(t, tcUserB, tcItemAvailable, "normal", "requested", base.Add(2*time.Hour))

	// When A เรียก GET /api/v1/my/pending-requests
	rec := tcGet(t, srv, tcPathMine, tcTokenA)

	// Then (1) ได้ 200 และ 2 รายการ
	tcWantStatus(t, rec, http.StatusOK)
	body := tcDecode[tcListBody](t, rec)
	if len(body.Items) != 2 {
		t.Fatalf("ได้ %d รายการ ต้องเป็น 2: %+v", len(body.Items), body.Items)
	}
	// Then (2) ทั้ง 2 รายการเป็นคำขอของ A
	got := map[int]bool{body.Items[0].RequestID: true, body.Items[1].RequestID: true}
	if !got[older] || !got[newer] {
		t.Errorf("ต้องได้คำขอ %d และ %d ของ A ได้ %+v", older, newer, body.Items)
	}
	// Then (3) เรียงจาก submittedAt ล่าสุดก่อน
	if body.Items[0].RequestID != newer || body.Items[1].RequestID != older {
		t.Errorf("ลำดับต้องเป็น [%d, %d] (ใหม่ไปเก่า) ได้ [%d, %d]", newer, older, body.Items[0].RequestID, body.Items[1].RequestID)
	}
	// Then (4) แต่ละรายการมี requestId, modelName, assetCode, submittedAt
	for _, it := range body.Items {
		if it.RequestID == 0 || it.ModelName == "" || it.AssetCode == "" || it.SubmittedAt == "" {
			t.Errorf("รายการต้องมี requestId, modelName, assetCode, submittedAt ครบ: %+v", it)
		}
	}
	// Then (5) ไม่มีคำขอ approved, rejected, cancelled ของ A
	for _, id := range []int{approved, rejected, cancelled} {
		if got[id] {
			t.Errorf("ไม่ควรเห็นคำขอ %d ที่ยกเลิกไม่ได้", id)
		}
	}
	// Then (6) ไม่มีคำขอของ B
	if got[ofB] {
		t.Errorf("ไม่ควรเห็นคำขอ %d ของ B", ofB)
	}
}

// TC-CXL-08-3 (ขอบ)
func TestTC_CXL_08_3_OvernightStatusesListed(t *testing.T) {
	// Given A มี requested 1, approved_by_officer 1, pending_advisor_approval 1 (ค้างคืน)
	srv := tcServer(t)
	base := time.Date(2026, 11, 10, 2, 0, 0, 0, time.UTC)
	r1 := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", base)
	r2 := tcInsert(t, tcUserA, tcItemOther, "overnight", "approved_by_officer", base.Add(time.Hour))
	r3 := tcInsert(t, tcUserA, tcItemAvailable2, "overnight", "pending_advisor_approval", base.Add(2*time.Hour)) // คนละชิ้น (test-cases.md แก้ 2569-11-11)

	// When A เรียก GET /api/v1/my/pending-requests
	rec := tcGet(t, srv, tcPathMine, tcTokenA)

	// Then (1) ได้ 3 รายการ เพราะทั้ง 3 สถานะนับเป็นรออนุมัติ
	tcWantStatus(t, rec, http.StatusOK)
	body := tcDecode[tcListBody](t, rec)
	if len(body.Items) != 3 {
		t.Fatalf("ได้ %d รายการ ต้องเป็น 3: %+v", len(body.Items), body.Items)
	}
	// Then (2) เรียงจาก submittedAt ล่าสุดก่อน
	want := []int{r3, r2, r1}
	for i, id := range want {
		if body.Items[i].RequestID != id {
			t.Errorf("ลำดับที่ %d ต้องเป็น %d ได้ %d", i, id, body.Items[i].RequestID)
		}
	}
}

// TC-CXL-08-4 (ขอบ)
func TestTC_CXL_08_4_SameSubmittedAtOrder(t *testing.T) {
	// Given A มีคำขอที่รออนุมัติ 2 คำขอ ที่ submitted_at เท่ากันพอดี
	srv := tcServer(t)
	same := time.Date(2026, 11, 10, 2, 0, 0, 0, time.UTC)
	first := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", same)
	second := tcInsert(t, tcUserA, tcItemOther, "normal", "requested", same)

	// When A เรียก GET /api/v1/my/pending-requests
	rec := tcGet(t, srv, tcPathMine, tcTokenA)

	// Then (1) ได้ 2 รายการ
	tcWantStatus(t, rec, http.StatusOK)
	body := tcDecode[tcListBody](t, rec)
	if len(body.Items) != 2 {
		t.Fatalf("ได้ %d รายการ ต้องเป็น 2", len(body.Items))
	}
	// Then (2) รายการที่ requestId มากกว่าอยู่ก่อน
	hi, lo := second, first
	if first > second {
		hi, lo = first, second
	}
	if body.Items[0].RequestID != hi || body.Items[1].RequestID != lo {
		t.Errorf("ลำดับต้องเป็น [%d, %d] ได้ [%d, %d]", hi, lo, body.Items[0].RequestID, body.Items[1].RequestID)
	}
}

// TC-CXL-08-5 (ผิด)
func TestTC_CXL_08_5_NoTokenUnauthenticated(t *testing.T) {
	// Given ไม่มี token (ยังไม่ยืนยันตัวตน) และมีคำขอของ A อยู่
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	// When เรียก GET /api/v1/my/pending-requests
	rec := tcGet(t, srv, tcPathMine, "")

	// Then (1) ได้ 401 UNAUTHENTICATED
	tcWantStatus(t, rec, http.StatusUnauthorized)
	if code := tcDecode[tcErrorBody](t, rec).Code; code != "UNAUTHENTICATED" {
		t.Errorf("code = %q ต้องเป็น UNAUTHENTICATED", code)
	}
	// Then (2) body ไม่มีข้อมูลคำขอใด
	if strings.Contains(rec.Body.String(), "items") || strings.Contains(rec.Body.String(), "requestId") {
		t.Errorf("body ต้องไม่มีข้อมูลคำขอ: %s", rec.Body.String())
	}
	_ = id
}

// TC-CXL-08-6 (ผิด)
func TestTC_CXL_08_6_StaffDoesNotSeeStudentRequests(t *testing.T) {
	// Given เจ้าหน้าที่ (staff-001) ยืนยันตัวตนแล้ว และมีคำขอที่รออนุมัติของ A อยู่
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	// When staff-001 เรียก GET /api/v1/my/pending-requests
	rec := tcGet(t, srv, tcPathMine, tcTokenStaff)

	// Then (1) ได้ 403 FORBIDDEN_ROLE (แก้ตาม CR-006)
	tcWantStatus(t, rec, http.StatusForbidden)
	if code := tcDecode[tcErrorBody](t, rec).Code; code != "FORBIDDEN_ROLE" {
		t.Errorf("code = %q ต้องเป็น FORBIDDEN_ROLE", code)
	}
	// Then (2) body ไม่มีข้อมูลคำขอของ A
	if strings.Contains(rec.Body.String(), "requestId") || strings.Contains(rec.Body.String(), strconv.Itoa(id)) {
		t.Errorf("body ต้องไม่มีข้อมูลคำขอของ A: %s", rec.Body.String())
	}
}
