package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-08-1 ถึง TC-BRQ-08-3 (AC-BRQ-08, FR-BRQ-08)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้" (TC-BRQ-08-3 ปลดล็อกเมื่อปิด Q-BRQ-17 วันเดียวกัน)
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
//
// ตาราง audit_log เป็น append-only (trigger ของ migration 0002) ลบแถวทิ้งไม่ได้
// test จึงใช้ผู้ขอเฉพาะของตัวเอง หรือนับจำนวนก่อนและหลัง

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// failAuditDB ทำให้เฉพาะการเขียน audit_log ล้มเหลว ส่วนคำสั่งอื่นวิ่งไปฐานข้อมูลจริงตามปกติ
// ใช้ทดสอบพฤติกรรม best-effort ที่ทีมเลือกตอนปิด Q-BRQ-17 เมื่อ 2569-11-13
type failAuditDB struct {
	Querier
}

func (d failAuditDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if strings.Contains(sql, "INSERT INTO audit_log") {
		return errRow{err: errors.New("จำลองเขียน audit log ไม่สำเร็จ")}
	}
	return d.Querier.QueryRow(ctx, sql, args...)
}

// auditRowsOf คืนข้อมูลของแถว audit ทุกแถวของผู้ใช้คนนั้น เรียงตามเวลา
type auditRow struct {
	actorID    string
	occurredAt time.Time
	assetCode  *string
	outcome    string
}

func auditRowsOf(t *testing.T, actorID string) []auditRow {
	t.Helper()
	rows, err := testPool(t).Query(context.Background(), `
		SELECT actor_id, occurred_at, asset_code, outcome FROM audit_log
		WHERE actor_id = $1 AND action = $2 ORDER BY id`, actorID, config.AuditActionRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var out []auditRow
	for rows.Next() {
		var r auditRow
		if err := rows.Scan(&r.actorID, &r.occurredAt, &r.assetCode, &r.outcome); err != nil {
			t.Fatal(err)
		}
		out = append(out, r)
	}
	return out
}

// TC-BRQ-08-1 (ทางปกติ) ส่งสำเร็จ 1 ครั้ง และถูกปฏิเสธ 1 ครั้ง
func TestTC_BRQ_08_1_TwoAuditRows(t *testing.T) {
	pool := testPool(t)
	const requester = "student-tc-08-1" // ผู้ขอของ test นี้เท่านั้น เพราะ audit_log ลบไม่ได้
	const assetCode = "TC-BRQ-08-1"
	insertTempItem(t, assetCode)
	ctx := studentCtx(requester)
	svc := NewServiceWithClock(NewRepository(pool), testClock())
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM borrow_request WHERE requester_id = $1`, requester)
	})

	// Given ส่งคำขอสำเร็จ 1 ครั้ง
	if _, err := svc.CreateWithAttachment(ctx, newTestRequest(assetCode), nil); err != nil {
		t.Fatalf("คำขอใบแรกต้องสำเร็จ แต่ได้ error: %v", err)
	}
	// Given ถูกปฏิเสธ 1 ครั้ง (ส่งซ้ำของชิ้นเดิม FR-BRQ-05 ข)
	if _, err := svc.CreateWithAttachment(ctx, newTestRequest(assetCode), nil); !errors.Is(err, ErrDuplicateRequest) {
		t.Fatalf("คำขอใบที่สองต้องถูกปฏิเสธด้วย ErrDuplicateRequest แต่ได้ %v", err)
	}

	// When ตรวจตาราง audit_log
	rows := auditRowsOf(t, requester)

	// Then (1) มี 2 แถวที่ action = 'request'
	if len(rows) != 2 {
		t.Fatalf("(1) ต้องมี audit 2 แถว แต่พบ %d แถว", len(rows))
	}
	for i, r := range rows {
		// Then (2) ทุกแถวมีผู้ขอ
		if r.actorID != requester {
			t.Errorf("(2) แถวที่ %d ต้องมีผู้ขอ %q แต่ได้ %q", i+1, requester, r.actorID)
		}
		// Then (3) ทุกแถวมีเวลา
		if r.occurredAt.IsZero() {
			t.Errorf("(3) แถวที่ %d ต้องมีเวลาที่บันทึก", i+1)
		}
		// Then (4) ทุกแถวมีรหัสครุภัณฑ์
		if r.assetCode == nil || *r.assetCode != assetCode {
			t.Errorf("(4) แถวที่ %d ต้องมีรหัสครุภัณฑ์ %q แต่ได้ %v", i+1, assetCode, r.assetCode)
		}
	}
	// Then (5) ผลลัพธ์ของสองแถวต่างกัน
	if rows[0].outcome == rows[1].outcome {
		t.Errorf("(5) ผลลัพธ์ของครั้งที่สำเร็จกับครั้งที่ถูกปฏิเสธต้องต่างกัน แต่เป็น %q ทั้งคู่", rows[0].outcome)
	}
}

// TC-BRQ-08-2 (ขอบ) ถูกปฏิเสธตั้งแต่การตรวจข้อมูล (422) ซึ่งยังไม่มีคำขอในระบบเลย
func TestTC_BRQ_08_2_AuditOnValidationReject(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const actor = "student-001"
	const assetCode = "TC-BRQ-08-2"
	insertTempItem(t, assetCode)

	// Given นับ audit ก่อน เพราะตารางนี้ลบแถวไม่ได้
	before := countAudit(t, pool, actor)

	// When ส่งคำขอที่ข้อมูลไม่ครบ จนได้ 422
	w := postRequest(t, srv, tokenStudentDept1, bodyWithout(assetCode, "endAt"))
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("ต้องถูกปฏิเสธด้วย 422 ก่อน จึงจะทดสอบแถวนี้ได้ แต่ได้ %d", w.Code)
	}

	// Then (1) ต้องมีแถวบันทึกการปฏิเสธครั้งนั้นเพิ่มขึ้น 1 แถว (FR-BRQ-08 "ทุกครั้ง")
	if after := countAudit(t, pool, actor); after != before+1 {
		t.Fatalf("(1) audit ต้องเพิ่มขึ้น 1 แถว แต่เปลี่ยนจาก %d เป็น %d", before, after)
	}
	// Then (2) แถวมีรหัสครุภัณฑ์ที่ผู้ขอพยายามยืม
	_, lastAsset := lastAuditOutcome(t, pool, actor)
	if lastAsset == nil || *lastAsset != assetCode {
		t.Errorf("(2) แถวล่าสุดต้องมีรหัสครุภัณฑ์ %q แต่ได้ %v", assetCode, lastAsset)
	}
}

// TC-BRQ-08-3 (ทางผิด) การเขียน audit_log ล้มเหลว ส่วนการสร้างคำขอสำเร็จแล้ว
// Q-BRQ-17 ปิดแล้ว 2569-11-13 ทีมเลือก best-effort กับการเขียนที่สำเร็จไปแล้ว เพราะการย้อนลบขัด FR-BRQ-09
func TestTC_BRQ_08_3_AuditWriteFails(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ATTACHMENT_DIR", dir)
	pool := testPool(t)
	const assetCode = "TC-BRQ-08-3"
	insertTempItem(t, assetCode)

	// Given การเขียน audit_log ล้มเหลว ส่วนคำสั่งอื่นทำงานปกติ และผู้ขอแนบไฟล์มาด้วย
	svc := NewServiceWithClock(NewRepository(failAuditDB{Querier: pool}), testClock())
	srv := auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes())
	ct, body := multipartForm(t, assetCode, config.LoanTypeOvernight,
		&testFile{name: "memo.pdf", contentType: "application/pdf", content: []byte("%PDF-1.4 ใบบันทึกข้อความ")})

	// When ส่งคำขอ
	w := postMultipart(t, srv, tokenStudentDept1, ct, body)

	// Then (1) ผู้ขอยังได้ 201 เหมือนเดิม
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) เขียน audit ไม่สำเร็จต้องไม่กระทบคำตอบ ต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	// Then (2) คำขอที่สร้างแล้วยังอยู่ในตาราง ไม่ถูกย้อนลบ
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(2) คำขอที่สร้างสำเร็จแล้วต้องไม่ถูกย้อนลบ แต่พบ %d ใบ", n)
	}
	// Then (3) ไฟล์แนบยังอยู่
	if n := countFilesIn(t, dir); n != 1 {
		t.Errorf("(3) ไฟล์แนบต้องยังอยู่ 1 ไฟล์ แต่พบ %d ไฟล์", n)
	}
}
