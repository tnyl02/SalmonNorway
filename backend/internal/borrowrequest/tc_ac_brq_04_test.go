package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-04-1 ถึง TC-BRQ-04-3 (AC-BRQ-04, FR-BRQ-04)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/config"
)

// countRequestsOnItem นับคำขอของผู้ขอคนหนึ่งบนอุปกรณ์ชิ้นหนึ่ง ใช้ตรวจว่าไม่มีคำขอเกิดซ้ำ
func countRequestsOnItem(t *testing.T, pool *pgxpool.Pool, requesterID, assetCode string) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), `
		SELECT count(*) FROM borrow_request r
		JOIN equipment_item i ON i.id = r.item_id
		WHERE r.requester_id = $1 AND i.asset_code = $2`, requesterID, assetCode).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// TC-BRQ-04-1 (ทางปกติ) ตรวจครบทั้ง 7 ส่วนของ Then
func TestTC_BRQ_04_1_CreateRequested(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-04-1"

	// Given ผู้ขอไม่มีของคืนเกินกำหนด (student-001 ตามข้อมูลตั้งต้น) และอุปกรณ์ชิ้นนั้นสถานะ available
	insertTempItem(t, assetCode)

	// When กรอกครบแล้วส่ง POST
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ส่งคำขอครบถ้วนต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	got := decodeCreated(t, w)

	var requesterID, loanType, status, itemAsset string
	var submittedAt time.Time
	if err := pool.QueryRow(context.Background(), `
		SELECT r.requester_id, r.loan_type, r.status, r.submitted_at, i.asset_code
		FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id
		WHERE r.id = $1`, got.RequestID).
		Scan(&requesterID, &loanType, &status, &submittedAt, &itemAsset); err != nil {
		t.Fatal(err)
	}

	// Then (2) สถานะคำขอเป็น requested
	if status != config.StatusRequested {
		t.Errorf("(2) สถานะที่บันทึกต้องเป็น %q แต่ได้ %q", config.StatusRequested, status)
	}
	// Then (3) requester_id เป็นรหัสผู้ใช้จาก IF-IDP-01 ไม่ใช่รหัสนักศึกษาที่ผู้ขอพิมพ์ (ASM-BRQ-06)
	if requesterID != "student-001" {
		t.Errorf("(3) ผู้ขอต้องเป็นรหัสผู้ใช้ student-001 แต่ได้ %q", requesterID)
	}
	// Then (4) อ้างอุปกรณ์รายชิ้นตรงกับ assetCode ที่ส่ง (ASM-BRQ-02)
	if itemAsset != assetCode {
		t.Errorf("(4) คำขอต้องผูกกับอุปกรณ์ %q แต่ได้ %q", assetCode, itemAsset)
	}
	// Then (5) loan_type ตรงกับที่เลือก
	if loanType != config.LoanTypeNormal {
		t.Errorf("(5) ประเภทการยืมต้องเป็น %q แต่ได้ %q", config.LoanTypeNormal, loanType)
	}
	// Then (6) submitted_at เทียบเป็นเวลา Asia/Bangkok แล้วตรงกับเวลาที่ส่ง
	if d := time.Since(submittedAt); d < -time.Minute || d > time.Minute {
		t.Errorf("(6) เวลาที่ส่งต้องใกล้เคียงเวลาปัจจุบัน แต่ต่างกัน %v", d)
	}
	if _, offset := submittedAt.In(appLoc).Zone(); offset != 7*60*60 {
		t.Errorf("(6) เวลาที่ส่งต้องเทียบได้กับเขตเวลา %s แต่ offset เป็น %d วินาที", config.TimeZone, offset)
	}
	// Then (7) สถานะอุปกรณ์ยังเป็น available เพราะ 002 ไม่เปลี่ยนสถานะอุปกรณ์ (DOM-REQ-01, CON-DATA-01)
	var itemStatus string
	if err := pool.QueryRow(context.Background(),
		`SELECT status FROM equipment_item WHERE asset_code = $1`, assetCode).Scan(&itemStatus); err != nil {
		t.Fatal(err)
	}
	if itemStatus != "available" {
		t.Errorf("(7) อุปกรณ์ที่มีคำขอค้างต้องยังเป็น available แต่เป็น %q", itemStatus)
	}
}

// TC-BRQ-04-2 (ขอบ) ประเภทการยืมอีกค่าหนึ่งก็ต้องได้สถานะ requested เหมือนกัน
func TestTC_BRQ_04_2_OvernightAlsoRequested(t *testing.T) {
	t.Setenv("ATTACHMENT_DIR", t.TempDir()) // ต้องตั้งก่อนสร้าง service เพราะอ่าน config ตอนสร้าง
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-04-2"

	// Given เลือกยืม "ค้างคืน" พร้อมแนบ PDF (FR-BRQ-10 ข้อ ก บังคับไฟล์แนบ)
	insertTempItem(t, assetCode)
	ct, body := multipartForm(t, assetCode, config.LoanTypeOvernight,
		&testFile{name: "memo.pdf", contentType: "application/pdf", content: []byte("%PDF-1.4 ใบบันทึกข้อความ")})

	// When ส่ง POST
	w := postMultipart(t, srv, tokenStudentDept1, ct, body)

	// Then (1) ได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ยืมค้างคืนที่แนบไฟล์ถูกกติกาต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	got := decodeCreated(t, w)

	var status, loanType string
	if err := pool.QueryRow(context.Background(),
		`SELECT status, loan_type FROM borrow_request WHERE id = $1`, got.RequestID).
		Scan(&status, &loanType); err != nil {
		t.Fatal(err)
	}
	// Then (2) สถานะยังเป็น requested
	if status != config.StatusRequested {
		t.Errorf("(2) สถานะต้องเป็น %q แต่ได้ %q", config.StatusRequested, status)
	}
	// Then (3) loan_type เป็น overnight
	if loanType != config.LoanTypeOvernight {
		t.Errorf("(3) ประเภทการยืมต้องเป็น %q แต่ได้ %q", config.LoanTypeOvernight, loanType)
	}
}

// TC-BRQ-04-3 (ทางผิด) ทำให้เงื่อนไขแรกใน Given ไม่เป็นจริง คือผู้ขอมีของคืนเกินกำหนดค้างอยู่
func TestTC_BRQ_04_3_OverdueBlocksCreate(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-04-3"

	// Given student-002 มีอุปกรณ์เลยกำหนดคืนค้างอยู่ตามข้อมูลตั้งต้น (T-03)
	insertTempItem(t, assetCode)
	before := countRequests(t, pool, "student-002")

	// When ส่ง POST
	w := postRequest(t, srv, tokenStudentOverdue, bodyFor(assetCode))

	// Then (1) ได้ 409 HAS_OVERDUE_LOAN
	if w.Code != http.StatusConflict {
		t.Fatalf("(1) ผู้ขอที่มีของเกินกำหนดต้องได้ 409 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if code := errorCodeOf(t, w); code != "HAS_OVERDUE_LOAN" {
		t.Errorf("(1) รหัสต้องเป็น HAS_OVERDUE_LOAN ตาม CON-TECH-03 แต่ได้ %q", code)
	}
	// Then (2) ไม่มีแถวใหม่
	if after := countRequests(t, pool, "student-002"); after != before {
		t.Errorf("(2) คำขอที่ถูกปฏิเสธต้องไม่ถูกสร้าง แต่จำนวนเปลี่ยนจาก %d เป็น %d", before, after)
	}
}
