package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-09-1 ถึง TC-BRQ-09-3 (AC-BRQ-09, FR-BRQ-09)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// failAttachmentUpdateDB ทำให้คำสั่งบันทึก path ของไฟล์แนบล้มเหลว
// ใช้จำลองกรณีที่ไฟล์ถูกเขียนลงดิสก์ไปแล้วแต่ฐานข้อมูลล้ม ซึ่งเป็นจุดที่ไฟล์อาจค้างได้จริง
type failAttachmentUpdateDB struct {
	Querier
}

func (d failAttachmentUpdateDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if strings.Contains(sql, "attachment_path") {
		return pgconn.CommandTag{}, errors.New("จำลองฐานข้อมูลขัดข้องตอนบันทึก path ของไฟล์แนบ")
	}
	return d.Querier.Exec(ctx, sql, args...)
}

// countFilesIn นับไฟล์ในโฟลเดอร์ ใช้ยืนยันว่าไม่มีไฟล์แนบค้างหลังบันทึกไม่สำเร็จ
func countFilesIn(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	return len(entries)
}

// TC-BRQ-09-1 (ทางปกติ) ล้มเหลวครั้งแรก แล้วสำเร็จในครั้งที่สอง
func TestTC_BRQ_09_1_RetryAfterFailure(t *testing.T) {
	pool := testPool(t)
	const assetCode = "TC-BRQ-09-1"
	insertTempItem(t, assetCode)

	// Given ฐานข้อมูลจะล้มเหลวเฉพาะการเขียนคำขอครั้งแรก
	svc := NewServiceWithClock(NewRepository(&failOnceDB{Querier: pool}), testClock())
	srv := auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes())

	// When ผู้ขอกดส่งคำขอครั้งแรก
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 503 UPSTREAM_UNAVAILABLE
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("(1) บันทึกไม่สำเร็จต้องได้ 503 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if code := errorCodeOf(t, w); code != "UPSTREAM_UNAVAILABLE" {
		t.Errorf("(1) รหัสต้องเป็น UPSTREAM_UNAVAILABLE ตาม CON-TECH-03 แต่ได้ %q", code)
	}
	// Then (2) หลังครั้งแรกตารางไม่มีแถวค้าง
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
		t.Fatalf("(2) บันทึกไม่สำเร็จต้องไม่ทิ้งคำขอค้างไว้ แต่พบ %d ใบ", n)
	}

	// When กดส่งใหม่อีกครั้ง
	w = postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (3) ครั้งที่สองได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(3) ส่งใหม่ต้องสำเร็จ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	// Then (4) สุดท้ายมีคำขอเพียง 1 รายการ
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(4) ต้องมีคำขอเพียง 1 รายการ แต่พบ %d ใบ", n)
	}
}

// TC-BRQ-09-2 (ขอบ) บันทึกล้มเหลวในคำขอยืมค้างคืนที่มีไฟล์แนบ
// spec ไม่ได้เขียนเรื่องไฟล์ ทีมตัดสิน 2569-11-13 ว่าต้องไม่เหลือไฟล์ค้าง
// ตรวจจุดที่ล้มได้ 2 จุด คือล้มตอนเขียนคำขอ (ไฟล์ยังไม่ถูกเขียน) และล้มตอนบันทึก path (ไฟล์ถูกเขียนไปแล้ว)
func TestTC_BRQ_09_2_FailureWithAttachment(t *testing.T) {
	pool := testPool(t)

	memo := &testFile{name: "memo.pdf", contentType: "application/pdf", content: []byte("%PDF-1.4 ใบบันทึกข้อความ")}

	t.Run("ล้มตอนเขียนคำขอ", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("ATTACHMENT_DIR", dir)
		const assetCode = "TC-BRQ-09-2a"
		insertTempItem(t, assetCode)

		// Given การบันทึกคำขอล้มเหลว และผู้ขอแนบไฟล์มาด้วย
		svc := NewServiceWithClock(NewRepository(&failOnceDB{Querier: pool}), testClock())
		srv := auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes())
		ct, body := multipartForm(t, assetCode, config.LoanTypeOvernight, memo)

		// When ส่งคำขอ
		w := postMultipart(t, srv, tokenStudentDept1, ct, body)

		// Then (1) ได้ 503
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("(1) ต้องได้ 503 แต่ได้ %d body=%s", w.Code, w.Body.String())
		}
		// Then (2) ไม่มีแถวค้าง
		if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
			t.Errorf("(2) ต้องไม่มีคำขอค้าง แต่พบ %d ใบ", n)
		}
		// Then (3) ไม่มีไฟล์แนบค้างในโฟลเดอร์
		if n := countFilesIn(t, dir); n != 0 {
			t.Errorf("(3) ต้องไม่มีไฟล์แนบค้าง แต่พบ %d ไฟล์", n)
		}
	})

	t.Run("ล้มตอนบันทึก path ของไฟล์แนบ", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("ATTACHMENT_DIR", dir)
		const assetCode = "TC-BRQ-09-2b"
		insertTempItem(t, assetCode)

		// Given ไฟล์ถูกเขียนลงดิสก์แล้ว แต่ฐานข้อมูลล้มตอนบันทึก path
		svc := NewServiceWithClock(NewRepository(failAttachmentUpdateDB{Querier: pool}), testClock())
		srv := auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes())
		ct, body := multipartForm(t, assetCode, config.LoanTypeOvernight, memo)

		// When ส่งคำขอ
		w := postMultipart(t, srv, tokenStudentDept1, ct, body)

		// Then (1) ได้ 503
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("(1) ต้องได้ 503 แต่ได้ %d body=%s", w.Code, w.Body.String())
		}
		// Then (2) ไม่มีแถวค้าง เพราะ service ย้อนลบคำขอให้
		if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 0 {
			t.Errorf("(2) ต้องไม่มีคำขอค้าง แต่พบ %d ใบ", n)
		}
		// Then (3) ไฟล์ที่เขียนไปแล้วต้องถูกลบ ไม่เหลือค้าง
		if n := countFilesIn(t, dir); n != 0 {
			t.Errorf("(3) ต้องไม่มีไฟล์แนบค้าง แต่พบ %d ไฟล์", n)
		}
	})
}

// TC-BRQ-09-3 (ทางผิด) ทำให้เงื่อนไข "บันทึกไม่สำเร็จ" ไม่เป็นจริง
func TestTC_BRQ_09_3_NoFailure(t *testing.T) {
	pool := testPool(t)
	srv := newTestServer(t)
	const assetCode = "TC-BRQ-09-3"

	// Given ฐานข้อมูลทำงานปกติ การบันทึกสำเร็จตั้งแต่ครั้งแรก
	insertTempItem(t, assetCode)

	// When ส่งคำขอ 1 ครั้ง
	w := postRequest(t, srv, tokenStudentDept1, bodyFor(assetCode))

	// Then (1) ได้ 201
	if w.Code != http.StatusCreated {
		t.Fatalf("(1) ต้องได้ 201 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	decodeCreated(t, w)
	// Then (2) มีคำขอ 1 รายการ
	if n := countRequestsOnItem(t, pool, "student-001", assetCode); n != 1 {
		t.Errorf("(2) ต้องมีคำขอ 1 รายการ แต่พบ %d ใบ", n)
	}
	// Then (3) ไม่มีข้อความแจ้งข้อผิดพลาดติดมาในคำตอบ
	if strings.Contains(w.Body.String(), "UPSTREAM_UNAVAILABLE") || strings.Contains(w.Body.String(), "code") {
		t.Errorf("(3) คำตอบที่สำเร็จต้องไม่มีรหัสข้อผิดพลาด แต่ได้ %s", w.Body.String())
	}
}
