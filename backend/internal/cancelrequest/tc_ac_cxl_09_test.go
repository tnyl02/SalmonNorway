package cancelrequest

// test cases ของ AC-CXL-09 จาก specs/003-cancel-request/test-cases.md (สถานะใช้ได้) ห้ามแก้ให้ผ่าน
// TC-CXL-09-1, 09-3, 09-4, 09-5 เป็น test หน้าจอ อยู่ที่ frontend/src/__tests__/
// TC-CXL-09-2 ให้คนลองเอง (ดูรายงานใน prompt-log.md)

import (
	"net/http"
	"testing"
	"time"
)

// TC-CXL-09-6 (ผิด)
func TestTC_CXL_09_6_SecondCancelNotCancellable(t *testing.T) {
	// Given คำขอของ A ถูกยกเลิกไปแล้วครั้งหนึ่ง (สถานะ cancelled)
	srv := tcServer(t)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())
	tcWantStatus(t, tcCancel(t, srv, id, "เหตุผลครั้งแรก", tcTokenA), http.StatusOK)

	// When ส่ง POST .../cancel ซ้ำอีกครั้ง
	rec := tcCancel(t, srv, id, "เหตุผลครั้งที่สอง", tcTokenA)

	// Then (1) ได้ 409 NOT_CANCELLABLE
	tcWantStatus(t, rec, http.StatusConflict)
	if code := tcDecode[tcErrorBody](t, rec).Code; code != "NOT_CANCELLABLE" {
		t.Errorf("code = %q ต้องเป็น NOT_CANCELLABLE", code)
	}
	// Then (2) สถานะยังเป็น cancelled และ cancel_reason ไม่ถูกเขียนทับ
	row := tcRow(t, id)
	if row.Status != "cancelled" {
		t.Errorf("สถานะ = %s ต้องยังเป็น cancelled", row.Status)
	}
	if row.CancelReason == nil || *row.CancelReason != "เหตุผลครั้งแรก" {
		t.Errorf("cancel_reason = %v ต้องยังเป็น \"เหตุผลครั้งแรก\"", row.CancelReason)
	}
}
