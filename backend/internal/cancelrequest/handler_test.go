package cancelrequest

// test ของ task (ไม่ใช่ TC) ใช้ตัวช่วย tcServer / tcInsert / tcGet จาก tc_helpers_test.go

import (
	"context"
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

// AC-CXL-02 (FR-CXL-02, FR-CXL-04, CON-CXL-02) T-15
//
//	Given A มีคำขอที่เจ้าหน้าที่อนุมัติแล้ว (approved)
//	When  A ยกเลิกคำขอนั้นพร้อมเหตุผล
//	Then  ได้ 409 ALREADY_APPROVED สถานะยัง approved ไม่มีเหตุผลหรือเวลายกเลิกถูกบันทึก
func TestAC_CXL_02_ApprovedNotCancellable(t *testing.T) {
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "approved", time.Now())

	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	tcWantStatus(t, rec, http.StatusConflict)
	if code := tcDecode[tcErrorBody](t, rec).Code; code != "ALREADY_APPROVED" {
		t.Errorf("code = %q ต้องเป็น ALREADY_APPROVED", code)
	}
	r := tcRow(t, id)
	if r.Status != "approved" || r.CancelReason != nil || r.CancelledAt != nil {
		t.Errorf("แถว = %+v ต้องยังเป็น approved ไม่มีเหตุผลและเวลายกเลิก", r)
	}
}

// AC-CXL-03 (FR-CXL-03, IF-IDP-01) T-19
//
//	Given B มีคำขอ requested
//	When  A สั่งยกเลิกคำขอของ B และสั่งยกเลิกรหัสคำขอที่ไม่มีในระบบ
//	Then  ได้ 404 REQUEST_NOT_FOUND เหมือนกันทั้ง 2 กรณี คำขอของ B ไม่เปลี่ยน และ body ไม่มีข้อมูลคำขอ
func TestAC_CXL_03_OthersRequestNotFound(t *testing.T) {
	srv := tcServer(t)
	ofB := tcInsert(t, tcUserB, tcItemAvailable, "normal", "requested", time.Now())

	for _, id := range []int{ofB, 999999} {
		rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)
		tcWantStatus(t, rec, http.StatusNotFound)
		if got := rec.Body.String(); got != "{\"code\":\"REQUEST_NOT_FOUND\"}\n" {
			t.Errorf("id %d: body = %q ต้องมีแค่ code REQUEST_NOT_FOUND", id, got)
		}
	}
	r := tcRow(t, ofB)
	if r.Status != "requested" || r.CancelReason != nil {
		t.Errorf("คำขอของ B = %+v ต้องยังเป็น requested ไม่มีเหตุผล", r)
	}
}

// AC-CXL-12 (FR-CXL-06) T-11
//
//	Given A ไม่มีคำขอที่ยังยกเลิกได้ (มีแต่ approved, rejected, cancelled) และ B มีคำขอ requested
//	When  A เรียก GET /v1/my/pending-requests
//	Then  ได้ 200 และ items เป็น [] (ไม่ใช่ null) ไม่มีคำขอของ B
func TestAC_CXL_12_EmptyList(t *testing.T) {
	srv := tcServer(t)
	now := time.Now()
	tcInsert(t, tcUserA, tcItemAvailable, "normal", "approved", now)
	tcInsert(t, tcUserA, tcItemAvailable, "normal", "rejected", now)
	tcInsert(t, tcUserA, tcItemAvailable, "normal", "cancelled", now)
	tcInsert(t, tcUserB, tcItemOther, "normal", "requested", now)

	rec := tcGet(t, srv, tcPathMine, tcTokenA)

	tcWantStatus(t, rec, http.StatusOK)
	if got := rec.Body.String(); got != "{\"items\":[]}\n" {
		t.Errorf("body = %q ต้องเป็น {\"items\":[]}", got)
	}
}

// AC-CXL-05 (DOM-REQ-01, CON-CXL-01, CON-DATA-01, FR-CXL-01) T-13
//
//	Given อุปกรณ์ X (available) มีคำขอที่ยังไม่อนุมัติ 2 คำขอ จาก A และ B
//	When  A ยกเลิกคำขอของตนเอง
//	Then  คำขอของ A เป็น cancelled คำขอของ B ยัง requested และ X ยัง available
//	      B ยังเห็นคำขอของตนในหน้ารายการ
func TestAC_CXL_05_OtherRequestsOnSameItemUnchanged(t *testing.T) {
	srv := tcServer(t)
	now := time.Now()
	ofA := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", now)
	ofB := tcInsert(t, tcUserB, tcItemAvailable, "normal", "requested", now.Add(time.Minute))

	tcWantStatus(t, tcCancel(t, srv, ofA, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA), http.StatusOK)

	if s := tcRow(t, ofA).Status; s != "cancelled" {
		t.Errorf("คำขอของ A = %s ต้องเป็น cancelled", s)
	}
	rb := tcRow(t, ofB)
	if rb.Status != "requested" || rb.CancelReason != nil || rb.CancelledAt != nil {
		t.Errorf("คำขอของ B = %+v ต้องไม่เปลี่ยน (DOM-REQ-01)", rb)
	}
	if s := tcItemStatus(t, tcItemAvailable); s != "available" {
		t.Errorf("สถานะอุปกรณ์ = %s ต้องยังเป็น available (CON-CXL-01)", s)
	}
	list := tcDecode[tcListBody](t, tcGet(t, srv, tcPathMine, "mock-student-002"))
	if len(list.Items) != 1 || list.Items[0].RequestID != ofB {
		t.Errorf("B ต้องยังเห็นคำขอ %d ได้ %+v", ofB, list.Items)
	}
}

// AC-CXL-10 (FR-CXL-01, CON-CXL-02, ASM-04) T-14
//
//	Given A มีคำขอค้างคืนสถานะ approved_by_officer และ pending_advisor_approval (คนละชิ้น)
//	When  A ยกเลิกแต่ละคำขอพร้อมเหตุผล
//	Then  ได้ 200 และทั้งคู่เป็น cancelled
func TestAC_CXL_10_CancelOvernightApprovedByOfficer(t *testing.T) {
	srv := tcServer(t)
	now := time.Now()
	officer := tcInsert(t, tcUserA, tcItemAvailable, "overnight", "approved_by_officer", now)
	advisor := tcInsert(t, tcUserA, tcItemAvailable2, "overnight", "pending_advisor_approval", now)

	for _, id := range []int{officer, advisor} {
		tcWantStatus(t, tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA), http.StatusOK)
		if s := tcRow(t, id).Status; s != "cancelled" {
			t.Errorf("คำขอ %d = %s ต้องเป็น cancelled", id, s)
		}
	}
}

// AC-CXL-11 (FR-CXL-02, CON-CXL-02) T-17
//
//	Given A มีคำขอ rejected, cancelled และ handed_over, returned (สถานะที่ยกเลิกไม่ได้ตาม CON-CXL-02)
//	When  A ส่งคำสั่งยกเลิกแต่ละคำขอโดยตรง
//	Then  ได้ 409 NOT_CANCELLABLE ทุกครั้ง และสถานะคงเดิม
func TestAC_CXL_11_RejectedAndCancelledNotCancellable(t *testing.T) {
	srv := tcServer(t)
	now := time.Now()
	ids := map[int]string{}
	for _, status := range []string{"rejected", "cancelled", "handed_over", "returned"} {
		ids[tcInsert(t, tcUserA, tcItemAvailable, "normal", status, now)] = status
	}

	for id, status := range ids {
		rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)
		tcWantStatus(t, rec, http.StatusConflict)
		if code := tcDecode[tcErrorBody](t, rec).Code; code != "NOT_CANCELLABLE" {
			t.Errorf("%s: code = %q ต้องเป็น NOT_CANCELLABLE", status, code)
		}
		if got := tcRow(t, id).Status; got != status {
			t.Errorf("%s: สถานะเปลี่ยนเป็น %s", status, got)
		}
	}
}

// AC-CXL-04 (FR-CXL-04, FR-CXL-02) T-21
//
//	Given A เปิดหน้ารายการเห็นคำขอ requested
//	      แล้วเจ้าหน้าที่อนุมัติคำขอนั้นก่อน A กดยืนยัน (จำลอง 004 ด้วย UPDATE ในฐานข้อมูล)
//	When  A กดยืนยันยกเลิก
//	Then  ได้ 409 ALREADY_APPROVED และสถานะยังเป็น approved
func TestAC_CXL_04_ApprovedBeforeConfirm(t *testing.T) {
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())
	list := tcDecode[tcListBody](t, tcGet(t, srv, tcPathMine, tcTokenA))
	if len(list.Items) != 1 || list.Items[0].RequestID != id {
		t.Fatalf("ก่อนอนุมัติ A ต้องเห็นคำขอ %d ได้ %+v", id, list.Items)
	}
	if _, err := tcPool(t).Exec(context.Background(), `UPDATE borrow_request SET status = 'approved' WHERE id = $1`, id); err != nil {
		t.Fatalf("จำลองการอนุมัติไม่สำเร็จ: %v", err)
	}

	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	tcWantStatus(t, rec, http.StatusConflict)
	if code := tcDecode[tcErrorBody](t, rec).Code; code != "ALREADY_APPROVED" {
		t.Errorf("code = %q ต้องเป็น ALREADY_APPROVED", code)
	}
	r := tcRow(t, id)
	if r.Status != "approved" || r.CancelReason != nil {
		t.Errorf("แถว = %+v ต้องยังเป็น approved ไม่มีเหตุผลยกเลิก", r)
	}
}

// AC-CXL-07 (IF-IDP-01, NFR-SEC-01) T-25 ส่วนหลังบ้าน (ส่วน audit_log อยู่ใน T-31 ส่วนหน้าจออยู่ใน T-26)
//
//	Given A มีคำขอ requested และ session หมดอายุ (ไม่มี token หรือ token ที่ระบบไม่รู้จัก)
//	When  ส่งคำสั่งยกเลิก
//	Then  ได้ 401 UNAUTHENTICATED และสถานะไม่เปลี่ยน
func TestAC_CXL_07_ExpiredSessionNoCancel(t *testing.T) {
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	for _, token := range []string{"", "mock-expired-token"} {
		rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", token)
		tcWantStatus(t, rec, http.StatusUnauthorized)
		if code := tcDecode[tcErrorBody](t, rec).Code; code != "UNAUTHENTICATED" {
			t.Errorf("token %q: code = %q ต้องเป็น UNAUTHENTICATED", token, code)
		}
	}
	r := tcRow(t, id)
	if r.Status != "requested" || r.CancelReason != nil {
		t.Errorf("แถว = %+v ต้องยังเป็น requested", r)
	}
}

