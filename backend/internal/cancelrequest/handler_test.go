package cancelrequest

// test ของ task (ไม่ใช่ TC) ใช้ตัวช่วย tcServer / tcInsert / tcGet จาก tc_helpers_test.go

import (
	"net/http"
	"testing"
	"time"
)

// อุปกรณ์ว่างเพิ่มเติมจาก seed.sql ของ 001 (ตาราง 0003 ห้ามคำขอที่ยังไม่จบของผู้ขอคนเดิมซ้ำชิ้น จึงใช้คนละชิ้น)
const (
	itemNB2 = 2 // CS-NB-0002
	itemNB3 = 3 // CS-NB-0003
)

// AC-CXL-08 (FR-CXL-06, FR-CXL-03, CON-CXL-02) T-04
//
//	Given A มีคำขอ requested 2 รายการ (ส่งต่างเวลา), approved_by_officer 1, pending_advisor_approval 1 (คนละชิ้น)
//	      approved, rejected, cancelled อย่างละ 1 และ B มีคำขอ requested 1
//	When  A เปิดหน้า "คำขอที่รออนุมัติของฉัน" (GET /v1/my/pending-requests)
//	Then  ได้เฉพาะคำขอของ A ที่ยังยกเลิกได้ 4 รายการ เรียงใหม่ไปเก่า เวลาเท่ากันเรียง requestId มากไปน้อย
//	      แต่ละรายการมี requestId, modelName, assetCode, submittedAt และเจ้าหน้าที่ได้ 403 STUDENT_ONLY
func TestAC_CXL_08_ListOwnPendingSorted(t *testing.T) {
	srv := tcServer(t)
	base := time.Date(2026, 11, 10, 2, 0, 0, 0, time.UTC)

	// Given
	r1 := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", base)
	r2 := tcInsert(t, tcUserA, tcItemOther, "normal", "requested", base.Add(2*time.Hour))
	r3 := tcInsert(t, tcUserA, itemNB2, "overnight", "approved_by_officer", base.Add(time.Hour))
	r4 := tcInsert(t, tcUserA, itemNB3, "overnight", "pending_advisor_approval", base.Add(time.Hour)) // เวลาเท่ากับ r3
	tcInsert(t, tcUserA, 8, "normal", "approved", base.Add(3*time.Hour))
	tcInsert(t, tcUserA, tcItemAvailable, "normal", "rejected", base.Add(3*time.Hour))
	tcInsert(t, tcUserA, tcItemAvailable, "normal", "cancelled", base.Add(3*time.Hour))
	tcInsert(t, tcUserB, tcItemAvailable, "normal", "requested", base.Add(4*time.Hour))

	// When
	rec := tcGet(t, srv, tcPathMine, tcTokenA)

	// Then
	tcWantStatus(t, rec, http.StatusOK)
	body := tcDecode[tcListBody](t, rec)
	hi, lo := r3, r4
	if r4 > r3 {
		hi, lo = r4, r3
	}
	want := []int{r2, hi, lo, r1}
	if len(body.Items) != len(want) {
		t.Fatalf("ได้ %d รายการ ต้องเป็น %d: %+v", len(body.Items), len(want), body.Items)
	}
	for i, id := range want {
		if body.Items[i].RequestID != id {
			t.Errorf("ลำดับที่ %d ต้องเป็น %d ได้ %d", i, id, body.Items[i].RequestID)
		}
	}
	first := body.Items[3] // r1 บน CS-NB-0001
	if first.ModelName != "Notebook ตัวอย่างรุ่น A" || first.AssetCode != "CS-NB-0001" {
		t.Errorf("r1 ต้องเป็น Notebook ตัวอย่างรุ่น A (CS-NB-0001) ได้ %+v", first)
	}
	if ts, err := time.Parse(time.RFC3339, first.SubmittedAt); err != nil || !ts.Equal(base) {
		t.Errorf("submittedAt = %q ต้องเป็นเวลา %v แบบ RFC3339", first.SubmittedAt, base)
	}

	// เจ้าหน้าที่ได้ 403 STUDENT_ONLY
	staff := tcGet(t, srv, tcPathMine, tcTokenStaff)
	tcWantStatus(t, staff, http.StatusForbidden)
	if code := tcDecode[tcErrorBody](t, staff).Code; code != "STUDENT_ONLY" {
		t.Errorf("code = %q ต้องเป็น STUDENT_ONLY", code)
	}
}

// FR-CXL-06: ไม่มีคำขอได้ items เป็น [] ไม่ใช่ null
func TestListPending_EmptyIsArray(t *testing.T) {
	srv := tcServer(t)
	rec := tcGet(t, srv, tcPathMine, tcTokenA)
	tcWantStatus(t, rec, http.StatusOK)
	if got := rec.Body.String(); got != "{\"items\":[]}\n" {
		t.Errorf("body = %q ต้องเป็น {\"items\":[]}", got)
	}
}
