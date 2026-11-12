package cancelrequest

// test cases ของ AC-CXL-01 จาก specs/003-cancel-request/test-cases.md (สถานะใช้ได้) ห้ามแก้ให้ผ่าน
// TC-CXL-01-2 เป็น test หน้าจอ อยู่ที่ frontend/src/__tests__/TC-CXL-01-2.test.jsx

import (
	"context"
	"net/http"
	"testing"
	"time"
)

// TC-CXL-01-1 (ปกติ)
func TestTC_CXL_01_1_CancelPendingNormal(t *testing.T) {
	// Given A ยืนยันตัวตนแล้ว มีคำขอแบบปกติ requested ของอุปกรณ์รายชิ้น X ซึ่งมีสถานะ available
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())
	if s := tcItemStatus(t, tcItemAvailable); s != "available" {
		t.Fatalf("เตรียมข้อมูลผิด: อุปกรณ์ต้องเป็น available ได้ %s", s)
	}

	// When A เรียก POST /api/v1/borrow-requests/{id}/cancel พร้อม reason "ไม่ต้องใช้อุปกรณ์แล้ว"
	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	// Then (1) ได้ 200 {requestId, status:"cancelled"}
	tcWantStatus(t, rec, http.StatusOK)
	body := tcDecode[tcCancelBody](t, rec)
	if body.RequestID != id || body.Status != "cancelled" {
		t.Errorf("body = %+v ต้องเป็น {requestId:%d status:cancelled}", body, id)
	}
	// Then (2) ในฐานข้อมูลสถานะของคำขอเป็น cancelled
	row := tcRow(t, id)
	if row.Status != "cancelled" {
		t.Errorf("สถานะของคำขอ = %s ต้องเป็น cancelled", row.Status)
	}
	// Then (3) cancel_reason ตรงกับที่ส่ง และ cancelled_at มีค่า
	if row.CancelReason == nil || *row.CancelReason != "ไม่ต้องใช้อุปกรณ์แล้ว" {
		t.Errorf("cancel_reason = %v ต้องเป็น \"ไม่ต้องใช้อุปกรณ์แล้ว\"", row.CancelReason)
	}
	if row.CancelledAt == nil {
		t.Errorf("cancelled_at ต้องมีค่า")
	}
	// Then (4) equipment_item.status ของ X ยังเป็น available
	if s := tcItemStatus(t, tcItemAvailable); s != "available" {
		t.Errorf("สถานะอุปกรณ์ = %s ต้องยังเป็น available (CON-CXL-01)", s)
	}
}

// TC-CXL-01-3 (ขอบ)
func TestTC_CXL_01_3_BorrowedItemStatusUnchanged(t *testing.T) {
	// Given อุปกรณ์ X สถานะ borrowed (มีคำขอ requested ของ A ค้างอยู่)
	srv := tcServer(t)
	if s := tcItemStatus(t, tcItemBorrowed); s != "borrowed" {
		t.Fatalf("เตรียมข้อมูลผิด: อุปกรณ์ต้องเป็น borrowed ได้ %s", s)
	}
	id := tcInsert(t, tcUserA, tcItemBorrowed, "normal", "requested", time.Now())

	// When A ยกเลิกคำขอนั้นพร้อมเหตุผล
	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	// Then (1) ได้ 200 สถานะของคำขอเป็น cancelled
	tcWantStatus(t, rec, http.StatusOK)
	if s := tcRow(t, id).Status; s != "cancelled" {
		t.Errorf("สถานะของคำขอ = %s ต้องเป็น cancelled", s)
	}
	// Then (2) equipment_item.status ของ X ยังเป็น borrowed ไม่เปลี่ยน
	if s := tcItemStatus(t, tcItemBorrowed); s != "borrowed" {
		t.Errorf("สถานะอุปกรณ์ = %s ต้องยังเป็น borrowed (CON-CXL-01)", s)
	}
}

// TC-CXL-01-4 (ขอบ)
func TestTC_CXL_01_4_CancelOvernightRequested(t *testing.T) {
	// Given คำขอแบบค้างคืนสถานะ requested ของ A
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "overnight", "requested", time.Now())

	// When A ยกเลิกคำขอนั้นพร้อมเหตุผล
	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	// Then (1) ได้ 200
	tcWantStatus(t, rec, http.StatusOK)
	// Then (2) สถานะของคำขอเป็น cancelled
	if s := tcRow(t, id).Status; s != "cancelled" {
		t.Errorf("สถานะของคำขอ = %s ต้องเป็น cancelled", s)
	}
}

// TC-CXL-01-5 (ผิด)
func TestTC_CXL_01_5_NoTokenNotCancelled(t *testing.T) {
	// Given ไม่มี token มีคำขอ requested ของ A
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	// When เรียก POST .../cancel พร้อมเหตุผล
	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", "")

	// Then (1) ได้ 401 UNAUTHENTICATED
	tcWantStatus(t, rec, http.StatusUnauthorized)
	if code := tcDecode[tcErrorBody](t, rec).Code; code != "UNAUTHENTICATED" {
		t.Errorf("code = %q ต้องเป็น UNAUTHENTICATED", code)
	}
	// Then (2) สถานะของคำขอยังเป็น requested
	if s := tcRow(t, id).Status; s != "requested" {
		t.Errorf("สถานะของคำขอ = %s ต้องยังเป็น requested", s)
	}
}

// TC-CXL-01-6 (ผิด)
func TestTC_CXL_01_6_ApprovedNotCancelled(t *testing.T) {
	// Given คำขอของ A สถานะ approved
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "approved", time.Now())

	// When A เรียก POST .../cancel พร้อมเหตุผล
	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	// Then (1) ได้ 409 ALREADY_APPROVED
	tcWantStatus(t, rec, http.StatusConflict)
	if code := tcDecode[tcErrorBody](t, rec).Code; code != "ALREADY_APPROVED" {
		t.Errorf("code = %q ต้องเป็น ALREADY_APPROVED", code)
	}
	// Then (2) สถานะของคำขอยังเป็น approved
	if s := tcRow(t, id).Status; s != "approved" {
		t.Errorf("สถานะของคำขอ = %s ต้องยังเป็น approved", s)
	}
}

// TC-CXL-01-7 (ผิด)
func TestTC_CXL_01_7_UnknownRequestNotFound(t *testing.T) {
	// Given รหัสคำขอที่ไม่มีในระบบ (มีคำขอ requested ของ A อยู่ 1 คำขอ ไว้ตรวจว่าไม่มีอะไรเปลี่ยน)
	srv := tcServer(t)
	existing := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	// When A เรียก POST /api/v1/borrow-requests/999999/cancel พร้อมเหตุผล
	rec := tcCancel(t, srv, 999999, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	// Then (1) ได้ 404 REQUEST_NOT_FOUND
	tcWantStatus(t, rec, http.StatusNotFound)
	if code := tcDecode[tcErrorBody](t, rec).Code; code != "REQUEST_NOT_FOUND" {
		t.Errorf("code = %q ต้องเป็น REQUEST_NOT_FOUND", code)
	}
	// Then (2) ไม่มีคำขอใดในฐานข้อมูลเปลี่ยนสถานะ
	var n int
	if err := tcPool(t).QueryRow(context.Background(),
		`SELECT count(*) FROM borrow_request WHERE status = 'cancelled'`).Scan(&n); err != nil {
		t.Fatalf("นับคำขอไม่สำเร็จ: %v", err)
	}
	if n != 0 {
		t.Errorf("มีคำขอถูกยกเลิก %d รายการ ต้องเป็น 0", n)
	}
	if s := tcRow(t, existing).Status; s != "requested" {
		t.Errorf("คำขอ %d สถานะ = %s ต้องยังเป็น requested", existing, s)
	}
}
