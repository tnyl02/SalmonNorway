package equipment

// test จาก specs/001-equipment-status/test-cases.md ของ AC-EQS-11 (ห้ามแก้ให้ผ่าน)

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/platform"
)

// auditRow คือ 1 แถวของ audit_log ที่ test อ่าน
type auditRow struct {
	ActorID, ActorRole, Action, AssetCode, Outcome string
	OccurredAt                                     time.Time
}

// countDetailAudit นับบันทึกการเปิดรายละเอียดของผู้ใช้และรหัสครุภัณฑ์นั้น
func countDetailAudit(t *testing.T, actorID, assetCode string) int {
	t.Helper()
	var n int
	err := testPool(t).QueryRow(context.Background(),
		`SELECT count(*) FROM audit_log WHERE actor_id = $1 AND asset_code = $2`, actorID, assetCode).Scan(&n)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// lastDetailAudit อ่านบันทึกล่าสุดของผู้ใช้และรหัสครุภัณฑ์นั้น
func lastDetailAudit(t *testing.T, actorID, assetCode string) auditRow {
	t.Helper()
	var r auditRow
	err := testPool(t).QueryRow(context.Background(), `
SELECT actor_id, actor_role, action, asset_code, outcome, occurred_at FROM audit_log
WHERE actor_id = $1 AND asset_code = $2 ORDER BY id DESC LIMIT 1`, actorID, assetCode).
		Scan(&r.ActorID, &r.ActorRole, &r.Action, &r.AssetCode, &r.Outcome, &r.OccurredAt)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// TC-EQS-11-1 (ปกติ)
func TestTC_EQS_11_1_DetailAuditRow(t *testing.T) {
	// Given mock-student-001 จดเวลาก่อนเรียก
	srv := newTestServer(t)
	n0 := countDetailAudit(t, "student-001", "CS-NB-0001")
	before := time.Now()

	// When เปิดรายละเอียด CS-NB-0001
	rec := get(t, srv, PathPrefix+"items/CS-NB-0001", tokenStudentDept1)
	after := time.Now()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}

	// Then มี audit log ใหม่ 1 แถว
	if n := countDetailAudit(t, "student-001", "CS-NB-0001"); n != n0+1 {
		t.Fatalf("audit log เพิ่ม %d แถว ต้องเพิ่ม 1", n-n0)
	}
	r := lastDetailAudit(t, "student-001", "CS-NB-0001")
	// (1) ผู้เข้าถึง student-001 บทบาท student
	if r.ActorID != "student-001" || r.ActorRole != "student" {
		t.Errorf("ผู้เข้าถึง = %s/%s ต้องเป็น student-001/student", r.ActorID, r.ActorRole)
	}
	// (2) เวลาอยู่ระหว่างก่อนและหลังเรียก (เผื่อนาฬิกาของ container ต่างกันได้ 1 วินาที)
	if r.OccurredAt.Before(before.Add(-time.Second)) || r.OccurredAt.After(after.Add(time.Second)) {
		t.Errorf("เวลา %s ไม่อยู่ระหว่าง %s ถึง %s", r.OccurredAt, before, after)
	}
	// (3) รหัสครุภัณฑ์ CS-NB-0001
	if r.AssetCode != "CS-NB-0001" {
		t.Errorf("รหัสครุภัณฑ์ = %q", r.AssetCode)
	}
	// (4) ประเภท detail ผล ok (plan ข้อ 6.4)
	if r.Action != "detail" || r.Outcome != "ok" {
		t.Errorf("ประเภท/ผล = %s/%s ต้องเป็น detail/ok", r.Action, r.Outcome)
	}
}

// TC-EQS-11-2 (ปกติ)
func TestTC_EQS_11_2_StaffDetailAudit(t *testing.T) {
	// Given mock-staff-001
	srv := newTestServer(t)
	n0 := countDetailAudit(t, "staff-001", "CS-NB-0004")

	// When เปิดรายละเอียด CS-NB-0004
	if rec := get(t, srv, PathPrefix+"items/CS-NB-0004", "mock-staff-001"); rec.Code != http.StatusOK {
		t.Fatalf("status = %d ต้องเป็น 200", rec.Code)
	}

	// Then มี audit log ใหม่ ผู้เข้าถึง staff-001 บทบาท staff รหัสครุภัณฑ์ CS-NB-0004
	if n := countDetailAudit(t, "staff-001", "CS-NB-0004"); n != n0+1 {
		t.Fatalf("audit log เพิ่ม %d แถว ต้องเพิ่ม 1", n-n0)
	}
	r := lastDetailAudit(t, "staff-001", "CS-NB-0004")
	if r.ActorID != "staff-001" || r.ActorRole != "staff" || r.AssetCode != "CS-NB-0004" {
		t.Errorf("บันทึก = %+v ต้องเป็น staff-001 / staff / CS-NB-0004", r)
	}
}

// TC-EQS-11-3 (ขอบ)
func TestTC_EQS_11_3_TwiceTwoRows(t *testing.T) {
	// Given mock-student-001
	srv := newTestServer(t)
	n0 := countDetailAudit(t, "student-001", "CS-NB-0001")

	// When เปิดรายละเอียด CS-NB-0001 2 ครั้งติดกัน
	for i := 0; i < 2; i++ {
		if rec := get(t, srv, PathPrefix+"items/CS-NB-0001", tokenStudentDept1); rec.Code != http.StatusOK {
			t.Fatalf("ครั้งที่ %d status = %d", i+1, rec.Code)
		}
	}

	// Then มี audit log เพิ่ม 2 แถว (DOM-AUDIT-01 "ทุกครั้งที่เข้าถึง")
	if n := countDetailAudit(t, "student-001", "CS-NB-0001"); n != n0+2 {
		t.Errorf("audit log เพิ่ม %d แถว ต้องเพิ่ม 2", n-n0)
	}
}

// TC-EQS-11-4 (ทางผิด)
func TestTC_EQS_11_4_ForbiddenAudited(t *testing.T) {
	// Given mock-student-101 (ภาควิชา 2)
	srv := newTestServer(t)
	n0 := countDetailAudit(t, "student-101", "CS-NB-0001")

	// When เปิดรายละเอียด CS-NB-0001 (ภาควิชา 1)
	rec := get(t, srv, PathPrefix+"items/CS-NB-0001", "mock-student-101")

	// Then (1) 403 FORBIDDEN_DEPARTMENT
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d ต้องเป็น 403", rec.Code)
	}
	if e := errorBody(t, rec.Body.String()); e["code"] != "FORBIDDEN_DEPARTMENT" {
		t.Errorf("code = %q ต้องเป็น FORBIDDEN_DEPARTMENT", e["code"])
	}
	// Then (2) มี audit log ผู้เข้าถึง student-101 รหัส CS-NB-0001 ผล forbidden (plan ข้อ 6.4)
	if n := countDetailAudit(t, "student-101", "CS-NB-0001"); n != n0+1 {
		t.Fatalf("audit log เพิ่ม %d แถว ต้องเพิ่ม 1", n-n0)
	}
	if r := lastDetailAudit(t, "student-101", "CS-NB-0001"); r.Outcome != "forbidden" {
		t.Errorf("ผล = %q ต้องเป็น forbidden", r.Outcome)
	}
}

// TC-EQS-11-5 (ทางผิด)
func TestTC_EQS_11_5_NotFoundAudit(t *testing.T) {
	// Given mock-student-001
	srv := newTestServer(t)
	n0 := countDetailAudit(t, "student-001", "CS-NB-9999")

	// When เปิดรายละเอียดรหัสที่ไม่มีอยู่ (CS-NB-9999)
	rec := get(t, srv, PathPrefix+"items/CS-NB-9999", tokenStudentDept1)

	// Then (1) 404 ITEM_NOT_FOUND
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d ต้องเป็น 404", rec.Code)
	}
	if e := errorBody(t, rec.Body.String()); e["code"] != "ITEM_NOT_FOUND" {
		t.Errorf("code = %q ต้องเป็น ITEM_NOT_FOUND", e["code"])
	}
	// Then (2) ไม่มี audit log ใหม่ (ทีมตัดสิน 2569-11-12 รับพฤติกรรมของโค้ด)
	if n := countDetailAudit(t, "student-001", "CS-NB-9999"); n != n0 {
		t.Errorf("audit log เพิ่ม %d แถว ต้องไม่เพิ่ม", n-n0)
	}
}

// auditFailDB ส่งทุก query ไปฐานข้อมูลจริง ยกเว้นการเขียน audit_log ที่ล้มเหลว (Given ของ TC-EQS-11-6)
type auditFailDB struct{ *pgxpool.Pool }

type failedRow struct{}

func (failedRow) Scan(...any) error {
	return errors.New("audit_log เขียนไม่ได้ (จำลอง)")
}

func (d auditFailDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if strings.Contains(sql, "INSERT INTO audit_log") {
		return failedRow{}
	}
	return d.Pool.QueryRow(ctx, sql, args...)
}

// TC-EQS-11-6 (ทางผิด)
func TestTC_EQS_11_6_AuditFailClosed(t *testing.T) {
	// Given เขียน audit log ไม่สำเร็จ
	h := NewHandler(NewService(NewRepository(auditFailDB{testPool(t)})))
	mux := http.NewServeMux()
	mux.Handle(PathPrefix, auth.Middleware(auth.NewMockProvider(), h.Routes()))
	srv := platform.NoStore(platform.MethodNotAllowedJSON(mux))

	// When mock-student-001 เปิดรายละเอียด CS-NB-0001
	rec := get(t, srv, PathPrefix+"items/CS-NB-0001", tokenStudentDept1)

	// Then (1) 503 UPSTREAM_UNAVAILABLE
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d ต้องเป็น 503", rec.Code)
	}
	if e := errorBody(t, rec.Body.String()); e["code"] != "UPSTREAM_UNAVAILABLE" {
		t.Errorf("code = %q ต้องเป็น UPSTREAM_UNAVAILABLE", e["code"])
	}
	// Then (2) ไม่มีข้อมูลอุปกรณ์ในคำตอบ: ยังไม่ตรวจ รอ Q-P4 (fail closed ของ plan ข้อ 6.4 ยังรอยืนยัน)
}

// TC-EQS-11-7 (ขอบ)
func TestTC_EQS_11_7_ListSearchAudit(t *testing.T) {
	// Given mock-student-001
	srv := newTestServer(t)

	// When เปิดหน้ารายการ และค้นหาด้วยตัวกรอง
	if rec := get(t, srv, PathPrefix+"models", tokenStudentDept1); rec.Code != http.StatusOK {
		t.Fatalf("รายการ status = %d", rec.Code)
	}
	if rec := get(t, srv, PathPrefix+"models?status=available", tokenStudentDept1); rec.Code != http.StatusOK {
		t.Fatalf("ค้นหา status = %d", rec.Code)
	}

	// Then บันทึก audit log ของการดูรายการและการค้นหาด้วยหรือไม่: ยังไม่ตรวจ รอ Q-11
	t.Skip("ยังไม่ตรวจ Then รอ Q-11")
}
