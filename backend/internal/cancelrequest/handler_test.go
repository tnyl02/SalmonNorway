package cancelrequest

// test ของ task (ไม่ใช่ TC) ใช้ตัวช่วย tcServer / tcInsert / tcGet จาก tc_helpers_test.go

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"salmonnorway/backend/internal/config"
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

// AC-CXL-01 (FR-CXL-01, FR-CXL-04, FR-CXL-08, CON-CXL-01) T-06
//
//	Given A มีคำขอแบบปกติ requested ของอุปกรณ์ X สถานะ available และ B มีคำขอ requested บน X เช่นกัน
//	When  A ยกเลิกคำขอของตนพร้อมเหตุผล
//	Then  ได้ 200 {requestId, status:"cancelled"} ฐานข้อมูลเป็น cancelled เก็บเหตุผลและเวลาที่ยกเลิก
//	      อุปกรณ์ X ยัง available คำขอของ B ไม่เปลี่ยน และคำขอนั้นหายจากรายการของ A
func TestAC_CXL_01_CancelPendingNormal(t *testing.T) {
	srv := tcServer(t)
	now := time.Now()
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", now)
	ofB := tcInsert(t, tcUserB, tcItemAvailable, "normal", "requested", now)

	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	tcWantStatus(t, rec, http.StatusOK)
	if body := tcDecode[tcCancelBody](t, rec); body.RequestID != id || body.Status != "cancelled" {
		t.Errorf("body = %+v ต้องเป็น {requestId:%d status:cancelled}", body, id)
	}
	row := tcRow(t, id)
	if row.Status != "cancelled" || row.CancelReason == nil || *row.CancelReason != "ไม่ต้องใช้อุปกรณ์แล้ว" || row.CancelledAt == nil {
		t.Errorf("แถวในฐานข้อมูล = %+v ต้องเป็น cancelled พร้อมเหตุผลและเวลา", row)
	}
	if s := tcItemStatus(t, tcItemAvailable); s != "available" {
		t.Errorf("สถานะอุปกรณ์ = %s ต้องยังเป็น available (CON-CXL-01)", s)
	}
	if s := tcRow(t, ofB).Status; s != "requested" {
		t.Errorf("คำขอของ B = %s ต้องยังเป็น requested (DOM-REQ-01)", s)
	}
	list := tcDecode[tcListBody](t, tcGet(t, srv, tcPathMine, tcTokenA))
	if len(list.Items) != 0 {
		t.Errorf("หลังยกเลิกรายการของ A ต้องว่าง ได้ %+v", list.Items)
	}
}

// FR-CXL-04: ยกเลิกซ้ำไม่เขียนทับเหตุผลเดิม เพราะ UPDATE มีเงื่อนไขสถานะ (ตอบ 409 จนกว่า T-17 จะแยกรหัส)
func TestCancel_SecondTimeDoesNotOverwrite(t *testing.T) {
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())
	tcWantStatus(t, tcCancel(t, srv, id, "ครั้งแรก", tcTokenA), http.StatusOK)

	rec := tcCancel(t, srv, id, "ครั้งที่สอง", tcTokenA)

	tcWantStatus(t, rec, http.StatusConflict)
	if r := tcRow(t, id); r.CancelReason == nil || *r.CancelReason != "ครั้งแรก" {
		t.Errorf("cancel_reason = %v ต้องยังเป็น ครั้งแรก", r.CancelReason)
	}
}

// AC-CXL-14 (FR-CXL-08) T-09
//
//	Given A มีคำขอ requested
//	When  ยกเลิกด้วยเหตุผลว่าง / ช่องว่างล้วน / 201 ตัวอักษร / ภาษาไทย 201 code point
//	Then  ได้ 400 REASON_REQUIRED หรือ REASON_TOO_LONG สถานะยัง requested
//	When  ยกเลิกด้วย "ไม่ต้องใช้อุปกรณ์แล้ว"
//	Then  ได้ 200 และ cancel_reason ตรงกัน
func TestAC_CXL_14_ReasonValidation(t *testing.T) {
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	cases := []struct {
		reason, code string
	}{
		{"", "REASON_REQUIRED"},
		{"   ", "REASON_REQUIRED"},
		{"\t\n ", "REASON_REQUIRED"},
		{strings.Repeat("a", config.CancelReasonMaxLen+1), "REASON_TOO_LONG"},
		{strings.Repeat("ที่", 67), "REASON_TOO_LONG"}, // 201 code point
	}
	for _, c := range cases {
		rec := tcCancel(t, srv, id, c.reason, tcTokenA)
		tcWantStatus(t, rec, http.StatusBadRequest)
		if code := tcDecode[tcErrorBody](t, rec).Code; code != c.code {
			t.Errorf("เหตุผลยาว %d: code = %q ต้องเป็น %s", len(c.reason), code, c.code)
		}
	}
	if s := tcRow(t, id).Status; s != "requested" {
		t.Fatalf("สถานะ = %s ต้องยังเป็น requested", s)
	}

	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)
	tcWantStatus(t, rec, http.StatusOK)
	if r := tcRow(t, id); r.CancelReason == nil || *r.CancelReason != "ไม่ต้องใช้อุปกรณ์แล้ว" {
		t.Errorf("cancel_reason = %v ต้องเป็น ไม่ต้องใช้อุปกรณ์แล้ว", r.CancelReason)
	}
}

// FR-CXL-08: ขอบ 200 code point พอดียกเลิกได้ (นับเหมือนหน้าจอ)
func TestValidateReason_Boundary(t *testing.T) {
	ok := strings.Repeat("ที่", 66) + "ab" // 200 code point
	if err := ValidateReason(ok); err != nil {
		t.Errorf("200 code point ต้องผ่าน ได้ %v", err)
	}
	if err := ValidateReason(ok + "c"); err != ErrReasonTooLong {
		t.Errorf("201 code point ต้องได้ ErrReasonTooLong ได้ %v", err)
	}
}

