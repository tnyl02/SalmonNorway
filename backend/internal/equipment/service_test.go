package equipment

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"salmonnorway/backend/internal/auth"
)

// ผู้ใช้จำลองสำหรับเรียก service ตรง ๆ (แบบเดียวกับที่ auth.Middleware ใส่ใน context)
var (
	studentCtx = auth.WithPrincipal(context.Background(), auth.Principal{UserID: "student-002", Role: auth.RoleStudent, DepartmentID: 1})
	staffCtx   = auth.WithPrincipal(context.Background(), auth.Principal{UserID: "staff-001", Role: auth.RoleStaff, DepartmentID: 1})
)

// insertBorrowed จำลอง UC-08 ส่งมอบอุปกรณ์ 1 ชิ้นแบบยืมปกติ กำหนดคืนตามที่ให้ แล้วลบทิ้งหลัง test
func insertBorrowed(t *testing.T, code string, dueAt time.Time) {
	t.Helper()
	pool := testPool(t)
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
WITH item AS (
  INSERT INTO equipment_item (asset_code, model_id, status) VALUES ($1, 1, 'borrowed') RETURNING id
)
INSERT INTO loan (item_id, borrower_id, loan_type, due_at)
SELECT id, 'student-001', 'normal', $2 FROM item`, code, dueAt)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM loan WHERE item_id = (SELECT id FROM equipment_item WHERE asset_code = $1)`, code)
		_, _ = pool.Exec(ctx, `DELETE FROM equipment_item WHERE asset_code = $1`, code)
	})
}

// serviceAt สร้าง Service ที่นาฬิกาหยุดอยู่ที่เวลาที่กำหนด (fix clock ตาม plan.md ข้อ 7)
func serviceAt(t *testing.T, now time.Time) *Service {
	t.Helper()
	s := NewService(NewRepository(testPool(t)))
	s.now = func() time.Time { return now }
	return s
}

// bangkok คืนเวลาตาม Asia/Bangkok (UTC+7) ปี ค.ศ.
func bangkok(t *testing.T, year int, month time.Month, day, hour, min int) time.Time {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		t.Fatal(err)
	}
	return time.Date(year, month, day, hour, min, 0, 0, loc)
}

func toJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// AC-EQS-14 (FR-EQS-04, DOM-PRIV-01, DOM-LOAN-01, CR-001, CR-002)
//
//	Given อุปกรณ์ชิ้นหนึ่งสถานะ "ถูกยืม" แบบปกติ กำหนดคืน 2569-10-10 16:00 และเวลาปัจจุบันคือ 2569-10-10 17:00
//	When  นักศึกษาเปิดดูรายละเอียด และเจ้าหน้าที่เปิดดูรายละเอียดเดียวกัน
//	Then  นักศึกษาเห็นสถานะ "ถูกยืม" กับป้าย "เลยกำหนดคืน" ไม่เห็นกำหนดคืน 2569-10-10 16:00
//	      ส่วนเจ้าหน้าที่เห็นทั้งกำหนดคืน 2569-10-10 16:00 และป้าย "เลยกำหนดคืน"
func TestAC_EQS_14_OverdueStudentVsStaff(t *testing.T) {
	// Given: 2569-10-10 = ค.ศ. 2026-10-10 กำหนดคืน 16:00 เวลาปัจจุบัน 17:00 (Asia/Bangkok)
	const code = "TEST-AC14-0001"
	due := bangkok(t, 2026, time.October, 10, 16, 0)
	insertBorrowed(t, code, due)
	svc := serviceAt(t, bangkok(t, 2026, time.October, 10, 17, 0))

	// When: นักศึกษาเปิดดูรายละเอียด
	got, found, err := svc.GetItem(studentCtx, code)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}

	// Then: นักศึกษาเห็น "ถูกยืม" + ป้ายเลยกำหนดคืน ไม่เห็นกำหนดคืน
	student, ok := got.(StudentItemDTO)
	if !ok {
		t.Fatalf("นักศึกษาต้องได้ StudentItemDTO ได้ %T", got)
	}
	if student.Status != StatusBorrowed {
		t.Errorf("สถานะต้องยังเป็น borrowed (ป้ายไม่ใช่สถานะใหม่ CON-DATA-01) ได้ %s", student.Status)
	}
	if !student.IsOverdue {
		t.Error("นักศึกษาต้องเห็นป้ายเลยกำหนดคืน (isOverdue=true)")
	}
	if student.DueAt != nil || strings.Contains(toJSON(t, student), "dueAt") {
		t.Errorf("นักศึกษาต้องไม่เห็นกำหนดคืนเมื่อเลยกำหนดแล้ว ได้ %s", toJSON(t, student))
	}

	// When: เจ้าหน้าที่เปิดดูรายละเอียดเดียวกัน
	got, _, err = svc.GetItem(staffCtx, code)
	if err != nil {
		t.Fatal(err)
	}

	// Then: เจ้าหน้าที่เห็นทั้งกำหนดคืน 16:00 และป้าย
	staff, ok := got.(StaffItemDTO)
	if !ok {
		t.Fatalf("เจ้าหน้าที่ต้องได้ StaffItemDTO ได้ %T", got)
	}
	if !staff.IsOverdue || staff.DueAt == nil || !staff.DueAt.Equal(due) || staff.Status != StatusBorrowed {
		t.Errorf("เจ้าหน้าที่ต้องเห็น borrowed + กำหนดคืน %s + ป้าย ได้ %s", due, toJSON(t, staff))
	}
}

// FR-EQS-04 ข้อ ก: ยังไม่เลยกำหนด (15:59 และตรง 16:00 พอดี) นักศึกษาเห็นกำหนดคืน ไม่มีป้าย
// เลยกำหนด = เวลาปัจจุบัน "มากกว่า" กำหนดคืน ตรงเวลาพอดียังไม่เลย
func TestFR_EQS_04_NotOverdueBeforeDueTime(t *testing.T) {
	const code = "TEST-AC14-0002"
	due := bangkok(t, 2026, time.October, 10, 16, 0)
	insertBorrowed(t, code, due)

	for _, now := range []time.Time{bangkok(t, 2026, time.October, 10, 15, 59), due} {
		got, _, err := serviceAt(t, now).GetItem(studentCtx, code)
		if err != nil {
			t.Fatal(err)
		}
		student := got.(StudentItemDTO)
		if student.IsOverdue {
			t.Errorf("เวลา %s ยังไม่เลยกำหนด แต่ได้ isOverdue=true", now.Format("15:04"))
		}
		if student.DueAt == nil || !student.DueAt.Equal(due) {
			t.Errorf("เวลา %s นักศึกษาต้องเห็นกำหนดคืน ได้ %s", now.Format("15:04"), toJSON(t, student))
		}
	}
}

// เทียบด้วยวันที่และเวลา ไม่ใช่แค่วันที่ (CR-002): วันถัดไป 00:01 ตามเวลาไทยเลยกำหนดแล้ว
// แม้ตามเวลา UTC ยังเป็นวันเดียวกับกำหนดคืน (2026-10-10T17:01Z)
func TestFR_EQS_04_OverdueUsesBangkokDateTime(t *testing.T) {
	const code = "TEST-AC14-0003"
	insertBorrowed(t, code, bangkok(t, 2026, time.October, 10, 16, 0))

	got, _, err := serviceAt(t, bangkok(t, 2026, time.October, 11, 0, 1)).GetItem(studentCtx, code)
	if err != nil {
		t.Fatal(err)
	}
	if !got.(StudentItemDTO).IsOverdue {
		t.Error("00:01 วันถัดไปตามเวลาไทยต้องเลยกำหนดคืนแล้ว")
	}
}
