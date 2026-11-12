package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-13 (ห้ามแก้ให้ผ่าน)

import (
	"net/http"
	"strings"
	"testing"
)

const tokenStudentDept2 = "mock-student-101"

// TC-EQS-13-1 (ปกติ)
func TestTC_EQS_13_1_ListOwnDeptOnly(t *testing.T) {
	// Given mock-student-101 (ภาควิชา 2)
	srv := newTestServer(t)

	// When เปิดหน้ารายการอุปกรณ์
	rec := get(t, srv, PathPrefix+"models", tokenStudentDept2)

	// Then (1) 200 มีเฉพาะรุ่น C
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	ids := map[int]bool{}
	for _, m := range decode[ListDTO[ModelRowDTO]](t, rec).Items {
		ids[m.ModelID] = true
	}
	if !ids[3] || len(ids) != 1 {
		t.Errorf("ต้องมีเฉพาะรุ่น C (3) ได้ %v", ids)
	}
	// Then (2) ไม่มีรุ่น A และรุ่น B ของภาควิชา 1 (FR-EQS-10, DOM-ACCESS-01)
	if ids[1] || ids[2] {
		t.Errorf("เห็นรุ่นของภาควิชา 1: %v", ids)
	}
}

// TC-EQS-13-2 (ปกติ)
func TestTC_EQS_13_2_DetailOtherDept403(t *testing.T) {
	// Given mock-student-101 (ภาควิชา 2)
	srv := newTestServer(t)

	// When เรียกรายละเอียด CS-NB-0001 (ภาควิชา 1) โดยตรง
	rec := get(t, srv, PathPrefix+"items/CS-NB-0001", tokenStudentDept2)

	// Then (1) 403 code FORBIDDEN_DEPARTMENT พร้อม message (CON-TECH-03)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d ต้องเป็น 403", rec.Code)
	}
	raw := rec.Body.String()
	if e := errorBody(t, raw); e["code"] != "FORBIDDEN_DEPARTMENT" || e["message"] == "" {
		t.Errorf("body ต้องเป็น {code: FORBIDDEN_DEPARTMENT, message} ได้ %s", raw)
	}
	// Then (2) คำตอบไม่มีรหัสครุภัณฑ์ สถานะ หรือข้อมูลอื่นของ CS-NB-0001
	for _, s := range []string{"CS-NB-0001", "assetCode", "status", "available", "dueAt"} {
		if strings.Contains(raw, s) {
			t.Errorf("คำตอบมี %q: %s", s, raw)
		}
	}
}

// TC-EQS-13-3 (ปกติ)
func TestTC_EQS_13_3_ItemsOtherDept403(t *testing.T) {
	// Given mock-student-101 (ภาควิชา 2)
	srv := newTestServer(t)

	// When เรียกรายชิ้นของรุ่น A (ภาควิชา 1) โดยตรง
	rec := get(t, srv, PathPrefix+"models/1/items", tokenStudentDept2)

	// Then (1) 403 FORBIDDEN_DEPARTMENT
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d ต้องเป็น 403", rec.Code)
	}
	raw := rec.Body.String()
	if e := errorBody(t, raw); e["code"] != "FORBIDDEN_DEPARTMENT" {
		t.Errorf("code = %q ต้องเป็น FORBIDDEN_DEPARTMENT", e["code"])
	}
	// Then (2) ไม่มีรายการอุปกรณ์ในคำตอบ
	if strings.Contains(raw, "items") || strings.Contains(raw, "CS-NB") {
		t.Errorf("คำตอบมีรายการอุปกรณ์: %s", raw)
	}
}

// TC-EQS-13-4 (ขอบ)
func TestTC_EQS_13_4_DeptParamOtherDept403(t *testing.T) {
	// Given mock-student-101 (ภาควิชา 2)
	srv := newTestServer(t)

	// When เรียกหน้ารายการโดยระบุภาควิชา 1 ใน query
	rec := get(t, srv, PathPrefix+"models?department=1", tokenStudentDept2)

	// Then (1) 403 FORBIDDEN_DEPARTMENT
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d ต้องเป็น 403", rec.Code)
	}
	raw := rec.Body.String()
	if e := errorBody(t, raw); e["code"] != "FORBIDDEN_DEPARTMENT" {
		t.Errorf("code = %q ต้องเป็น FORBIDDEN_DEPARTMENT", e["code"])
	}
	// Then (2) ไม่มีรุ่นของภาควิชา 1 ในคำตอบ
	if strings.Contains(raw, "Notebook") || strings.Contains(raw, "Projector") || strings.Contains(raw, "items") {
		t.Errorf("คำตอบมีรุ่นของภาควิชา 1: %s", raw)
	}
}

// TC-EQS-13-6 (ทางผิด)
func TestTC_EQS_13_6_OwnDeptDetail200(t *testing.T) {
	// Given mock-student-101 (ภาควิชา 2)
	srv := newTestServer(t)

	// When เปิดรายละเอียด XX-NB-0001 (ภาควิชาของตัวเอง)
	rec := get(t, srv, PathPrefix+"items/XX-NB-0001", tokenStudentDept2)

	// Then 200 เห็นสถานะ "ว่าง" ตามปกติ
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}
	if got := decode[StudentItemDTO](t, rec); got.Status != StatusAvailable {
		t.Errorf("สถานะ = %q ต้องเป็น available", got.Status)
	}
}

// TC-EQS-13-7 (ทางผิด)
func TestTC_EQS_13_7_StaffOtherDept(t *testing.T) {
	// Given mock-staff-001 (เจ้าหน้าที่ ภาควิชา 1)
	srv := newTestServer(t)

	// When เปิดรายละเอียด XX-NB-0001 (ภาควิชา 2)
	_ = get(t, srv, PathPrefix+"items/XX-NB-0001", "mock-staff-001")

	// Then เจ้าหน้าที่เห็นภาควิชาอื่นได้หรือไม่: ยังไม่ตรวจ รอ Q-P3 (plan ข้อ 6.1 ชั่วคราวให้เห็นทุกภาควิชา)
	t.Skip("ยังไม่ตรวจ Then รอ Q-P3")
}
