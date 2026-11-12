package cancelrequest

// test cases ของ AC-CXL-14 จาก specs/003-cancel-request/test-cases.md (สถานะใช้ได้) ห้ามแก้ให้ผ่าน
// TC-CXL-14-2 และส่วนหน้าจอของ TC-CXL-14-5 อยู่ที่ frontend/src/__tests__/

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// TC-CXL-14-1 (ปกติ)
func TestTC_CXL_14_1_ReasonStored(t *testing.T) {
	// Given A มีคำขอ requested
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	// When A ยกเลิกพร้อม reason "ไม่ต้องใช้อุปกรณ์แล้ว"
	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	// Then (1) ได้ 200
	tcWantStatus(t, rec, http.StatusOK)
	row := tcRow(t, id)
	// Then (2) สถานะเป็น cancelled
	if row.Status != "cancelled" {
		t.Errorf("สถานะ = %s ต้องเป็น cancelled", row.Status)
	}
	// Then (3) cancel_reason เท่ากับ "ไม่ต้องใช้อุปกรณ์แล้ว" ตรงตัว
	if row.CancelReason == nil || *row.CancelReason != "ไม่ต้องใช้อุปกรณ์แล้ว" {
		t.Errorf("cancel_reason = %v ต้องเป็น \"ไม่ต้องใช้อุปกรณ์แล้ว\" ตรงตัว", row.CancelReason)
	}
}

// TC-CXL-14-3 (ขอบ)
func TestTC_CXL_14_3_Reason200Accepted(t *testing.T) {
	// Given เหตุผลยาว 200 ตัวอักษรพอดี (ตัวอักษรภาษาอังกฤษ)
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())
	reason := strings.Repeat("a", 200)

	// When A ยกเลิกพร้อมเหตุผลนั้น
	rec := tcCancel(t, srv, id, reason, tcTokenA)

	// Then (1) ได้ 200
	tcWantStatus(t, rec, http.StatusOK)
	// Then (2) สถานะเป็น cancelled และ cancel_reason ยาว 200 ตัวอักษร
	row := tcRow(t, id)
	if row.Status != "cancelled" {
		t.Errorf("สถานะ = %s ต้องเป็น cancelled", row.Status)
	}
	if row.CancelReason == nil || utf8.RuneCountInString(*row.CancelReason) != 200 {
		t.Errorf("cancel_reason ต้องยาว 200 ตัวอักษร ได้ %v", row.CancelReason)
	}
}

// TC-CXL-14-4 (ขอบ)
func TestTC_CXL_14_4_Reason201Rejected(t *testing.T) {
	// Given เหตุผลยาว 201 ตัวอักษร (ตัวอักษรภาษาอังกฤษ)
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	// When A ยกเลิกพร้อมเหตุผลนั้น
	rec := tcCancel(t, srv, id, strings.Repeat("a", 201), tcTokenA)

	// Then (1) ได้ 422 VALIDATION_FAILED แนบ fields reason (แก้ตาม CR-006)
	tcWantStatus(t, rec, http.StatusUnprocessableEntity)
	if b := tcDecode[tcErrorBody](t, rec); b.Code != "VALIDATION_FAILED" || len(b.Fields) != 1 || b.Fields[0] != "reason" {
		t.Errorf("body = %+v ต้องเป็น VALIDATION_FAILED fields [reason]", b)
	}
	// Then (2) สถานะยังเป็น requested
	if s := tcRow(t, id).Status; s != "requested" {
		t.Errorf("สถานะ = %s ต้องยังเป็น requested", s)
	}
}

// TC-CXL-14-5 (ขอบ) ส่วนหลังบ้าน ส่วนหน้าจออยู่ที่ frontend/src/__tests__/TC-CXL-14-5.test.jsx
func TestTC_CXL_14_5_ThaiCharacterCounting(t *testing.T) {
	// Given เหตุผลภาษาไทย "ที่" ซ้ำ 67 ครั้ง (เห็นบนจอ 67 ตัว แต่เป็น 201 code point)
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())
	reason := strings.Repeat("ที่", 67)
	if n := utf8.RuneCountInString(reason); n != 201 {
		t.Fatalf("เตรียมข้อมูลผิด: ต้องเป็น 201 code point ได้ %d", n)
	}

	// When A ยกเลิกพร้อมเหตุผลนั้น
	rec := tcCancel(t, srv, id, reason, tcTokenA)

	// Then (1) ได้ 422 VALIDATION_FAILED เพราะนับตาม code point ได้ 201 (แก้ตาม CR-006)
	tcWantStatus(t, rec, http.StatusUnprocessableEntity)
	if code := tcDecode[tcErrorBody](t, rec).Code; code != "VALIDATION_FAILED" {
		t.Errorf("code = %q ต้องเป็น VALIDATION_FAILED", code)
	}
	// Then (2) สถานะยังเป็น requested
	if s := tcRow(t, id).Status; s != "requested" {
		t.Errorf("สถานะ = %s ต้องยังเป็น requested", s)
	}
	// Then (3) หน้าจอนับแบบเดียวกัน: ตรวจใน TC-CXL-14-5.test.jsx
}

// TC-CXL-14-6 (ผิด)
func TestTC_CXL_14_6_EmptyReasonRejected(t *testing.T) {
	// Given A มีคำขอ requested
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	for _, reason := range []string{"", "   "} {
		// When ยกเลิกด้วย reason "" แล้วยกเลิกด้วย reason "   " (ช่องว่างล้วน)
		rec := tcCancel(t, srv, id, reason, tcTokenA)

		// Then (1) ได้ 422 VALIDATION_FAILED ทั้ง 2 ครั้ง (แก้ตาม CR-006)
		tcWantStatus(t, rec, http.StatusUnprocessableEntity)
		if code := tcDecode[tcErrorBody](t, rec).Code; code != "VALIDATION_FAILED" {
			t.Errorf("reason %q: code = %q ต้องเป็น VALIDATION_FAILED", reason, code)
		}
	}
	// Then (2) สถานะยังเป็น requested และ cancel_reason ว่าง
	row := tcRow(t, id)
	if row.Status != "requested" {
		t.Errorf("สถานะ = %s ต้องยังเป็น requested", row.Status)
	}
	var reason *string
	_ = tcPool(t).QueryRow(context.Background(), `SELECT cancel_reason FROM borrow_request WHERE id = $1`, id).Scan(&reason)
	if reason != nil && *reason != "" {
		t.Errorf("cancel_reason = %q ต้องว่าง", *reason)
	}
}
