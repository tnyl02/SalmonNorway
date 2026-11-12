package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-14 (ห้ามแก้ให้ผ่าน)

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TC-EQS-14-1 (ปกติ)
func TestTC_EQS_14_1_StudentOverdue(t *testing.T) {
	// Given ยืมปกติ กำหนดคืน 2569-10-10 16:00 เวลาปัจจุบัน 2569-10-10 17:00
	const code = "TC14-1-0001"
	insertLoanItem(t, code, "normal", bangkok(t, 2026, time.October, 10, 16, 0))
	svc := serviceAt(t, bangkok(t, 2026, time.October, 10, 17, 0))

	// When นักศึกษาเปิดรายละเอียด
	got, found, err := svc.GetItem(studentCtx, code)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	student := got.(StudentItemDTO)

	// Then (1) สถานะ "ถูกยืม" ไม่ใช่สถานะใหม่ (CON-DATA-01)
	if student.Status != StatusBorrowed {
		t.Errorf("สถานะ = %q ต้องเป็น borrowed", student.Status)
	}
	// Then (2) มีป้ายเลยกำหนดคืน
	if !student.IsOverdue {
		t.Error("ต้องมีป้ายเลยกำหนดคืน")
	}
	// Then (3) ไม่มีกำหนดคืนในคำตอบ (CR-001)
	if student.DueAt != nil || strings.Contains(toJSON(t, student), "dueAt") {
		t.Errorf("นักศึกษาต้องไม่ได้กำหนดคืน: %s", toJSON(t, student))
	}
	// Then (4) ส่งคำขอยืมไม่ได้
	if student.CanRequestLoan {
		t.Error("canRequestLoan ต้องเป็น false")
	}
}

// TC-EQS-14-2 (ปกติ)
func TestTC_EQS_14_2_StaffOverdue(t *testing.T) {
	// Given เหมือน TC-EQS-14-1
	const code = "TC14-2-0001"
	due := bangkok(t, 2026, time.October, 10, 16, 0)
	insertLoanItem(t, code, "normal", due)
	svc := serviceAt(t, bangkok(t, 2026, time.October, 10, 17, 0))

	// When เจ้าหน้าที่เปิดรายละเอียดเดียวกัน
	got, found, err := svc.GetItem(staffCtx, code)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	staff := got.(StaffItemDTO)

	// Then (1) สถานะ "ถูกยืม"
	if staff.Status != StatusBorrowed {
		t.Errorf("สถานะ = %q ต้องเป็น borrowed", staff.Status)
	}
	// Then (2) มีกำหนดคืน 2569-10-10 16:00
	if staff.DueAt == nil || !staff.DueAt.Equal(due) {
		t.Errorf("กำหนดคืน = %v ต้องเป็น %s", staff.DueAt, due)
	}
	// Then (3) มีป้ายเลยกำหนดคืน
	if !staff.IsOverdue {
		t.Error("ต้องมีป้ายเลยกำหนดคืน")
	}
}

// TC-EQS-14-3 (ปกติ)
func TestTC_EQS_14_3_HttpStudentVsStaff(t *testing.T) {
	// Given CS-NB-0005 ยืมปกติ กำหนดคืน 16:00 ของเมื่อวาน (ข้อมูลตั้งต้น เลยกำหนดเสมอ)
	srv := newTestServer(t)

	// When mock-student-002 และ mock-staff-001 เรียกรายละเอียด CS-NB-0005 ผ่าน HTTP
	st := get(t, srv, PathPrefix+"items/CS-NB-0005", "mock-student-002")
	sf := get(t, srv, PathPrefix+"items/CS-NB-0005", "mock-staff-001")
	if st.Code != http.StatusOK || sf.Code != http.StatusOK {
		t.Fatalf("status นักศึกษา %d เจ้าหน้าที่ %d ต้องเป็น 200", st.Code, sf.Code)
	}

	// Then (1) JSON ของนักศึกษามี isOverdue true ไม่มีช่อง dueAt
	if raw := st.Body.String(); !strings.Contains(raw, `"isOverdue":true`) || strings.Contains(raw, "dueAt") {
		t.Errorf("นักศึกษาต้องได้ isOverdue true ไม่มี dueAt: %s", raw)
	}
	// Then (2) JSON ของเจ้าหน้าที่มีทั้ง dueAt และ isOverdue true
	if raw := sf.Body.String(); !strings.Contains(raw, `"isOverdue":true`) || !strings.Contains(raw, `"dueAt"`) {
		t.Errorf("เจ้าหน้าที่ต้องได้ทั้ง dueAt และ isOverdue true: %s", raw)
	}
}

// TC-EQS-14-6 (ขอบ)
func TestTC_EQS_14_6_ExactlyDueNotOverdue(t *testing.T) {
	// Given กำหนดคืน 2569-10-10 16:00 เวลาปัจจุบัน 16:00 พอดี
	const code = "TC14-6-0001"
	due := bangkok(t, 2026, time.October, 10, 16, 0)
	insertLoanItem(t, code, "normal", due)
	svc := serviceAt(t, due)

	// When นักศึกษาเปิดรายละเอียด
	got, _, err := svc.GetItem(studentCtx, code)
	if err != nil {
		t.Fatal(err)
	}
	student := got.(StudentItemDTO)

	// Then ยังไม่เลยกำหนด (เวลาปัจจุบัน > กำหนดคืน เท่านั้นจึงเลย)
	// (1) เห็นกำหนดคืน
	if student.DueAt == nil || !student.DueAt.Equal(due) {
		t.Errorf("กำหนดคืน = %v ต้องเป็น %s", student.DueAt, due)
	}
	// (2) ไม่มีป้ายเลยกำหนดคืน
	if student.IsOverdue {
		t.Error("16:00 พอดียังไม่เลยกำหนด ต้องไม่มีป้าย")
	}
}

// TC-EQS-14-7 (ขอบ)
func TestTC_EQS_14_7_OneMinuteAfterDue(t *testing.T) {
	// Given กำหนดคืน 2569-10-10 16:00 เวลาปัจจุบัน 16:01
	const code = "TC14-7-0001"
	insertLoanItem(t, code, "normal", bangkok(t, 2026, time.October, 10, 16, 0))
	svc := serviceAt(t, bangkok(t, 2026, time.October, 10, 16, 1))

	// When นักศึกษาเปิดรายละเอียด
	got, _, err := svc.GetItem(studentCtx, code)
	if err != nil {
		t.Fatal(err)
	}
	student := got.(StudentItemDTO)

	// Then เลยกำหนดแล้ว (1) มีป้ายเลยกำหนดคืน
	if !student.IsOverdue {
		t.Error("16:01 เลยกำหนดแล้ว ต้องมีป้าย")
	}
	// (2) ไม่มีกำหนดคืน
	if student.DueAt != nil {
		t.Errorf("นักศึกษาต้องไม่ได้กำหนดคืน ได้ %v", student.DueAt)
	}
}

// TC-EQS-14-8 (ทางผิด)
func TestTC_EQS_14_8_OvernightOverdue(t *testing.T) {
	// Given ยืมแบบค้างคืน วันสิ้นสุด 2569-10-11 กำหนดคืน 2569-10-11 16:00 (Q-SH-04) เวลาปัจจุบัน 2569-10-11 17:00
	const code = "TC14-8-0001"
	due := bangkok(t, 2026, time.October, 11, 16, 0)
	insertLoanItem(t, code, "overnight", due)
	svc := serviceAt(t, bangkok(t, 2026, time.October, 11, 17, 0))

	// When นักศึกษาและเจ้าหน้าที่เปิดรายละเอียด
	gotSt, _, err := svc.GetItem(studentCtx, code)
	if err != nil {
		t.Fatal(err)
	}
	gotSf, _, err := svc.GetItem(staffCtx, code)
	if err != nil {
		t.Fatal(err)
	}

	// Then ใช้กติกาเดียวกับยืมปกติ
	// (1) นักศึกษาเห็นป้ายเลยกำหนดคืน ไม่มีกำหนดคืน
	if st := gotSt.(StudentItemDTO); !st.IsOverdue || st.DueAt != nil {
		t.Errorf("นักศึกษาต้องได้ป้าย ไม่มีกำหนดคืน ได้ %s", toJSON(t, st))
	}
	// (2) เจ้าหน้าที่เห็นทั้งกำหนดคืน 2569-10-11 16:00 และป้าย
	if sf := gotSf.(StaffItemDTO); !sf.IsOverdue || sf.DueAt == nil || !sf.DueAt.Equal(due) {
		t.Errorf("เจ้าหน้าที่ต้องได้กำหนดคืน %s และป้าย ได้ %s", due, toJSON(t, sf))
	}
}

// TC-EQS-14-9 (ทางผิด)
func TestTC_EQS_14_9_ReturnedNoOverdue(t *testing.T) {
	// Given ชิ้นที่เคยถูกยืมเลยกำหนด แต่รับคืนแล้ว (สถานะ "ว่าง" การยืมปิดแล้ว)
	const code = "TC14-9-0001"
	insertLoanItem(t, code, "normal", bangkok(t, 2026, time.October, 10, 16, 0))
	returnItem(t, code, StatusAvailable)
	svc := serviceAt(t, bangkok(t, 2026, time.October, 10, 17, 0))

	// When นักศึกษาและเจ้าหน้าที่เปิดรายละเอียด
	for _, ctx := range []context.Context{studentCtx, staffCtx} {
		got, _, err := svc.GetItem(ctx, code)
		if err != nil {
			t.Fatal(err)
		}
		raw := toJSON(t, got)

		// Then (1) สถานะ "ว่าง"
		if !strings.Contains(raw, `"status":"available"`) {
			t.Errorf("สถานะต้องเป็น available: %s", raw)
		}
		// Then (2) ไม่มีป้ายเลยกำหนดคืน (FR-EQS-04 ป้ายมีเฉพาะ "ถูกยืม")
		if strings.Contains(raw, "isOverdue") {
			t.Errorf("ต้องไม่มีป้ายเลยกำหนดคืน: %s", raw)
		}
		// Then (3) ไม่มีกำหนดคืน (CR-002)
		if strings.Contains(raw, "dueAt") {
			t.Errorf("ต้องไม่มีกำหนดคืน: %s", raw)
		}
	}
}
