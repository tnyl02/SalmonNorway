// test จากตาราง specs/004-approve-request/test-cases.md แถว TC-APV-04-x, 10-x, 13-x, 14-x, 15-x, 18-x สถานะ "ใช้ได้" ห้ามแก้ให้ผ่าน
// ถ้าไม่ผ่าน ต้องหาว่าโค้ดผิด หรือ AC ผิด (ถ้า AC ผิด แก้ spec ก่อน)
// ใช้ helper จาก tc_ac_apv_*_test.go ไฟล์อื่น และ fakeMailer, fakeClock, concurrencyDB, setSubmittedAt ของ test เดิม
package approve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"salmonnorway/backend/internal/config"
)

const tcMsgConcurrent = "คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น กรุณารีเฟรชหน้า"
const tcAutoReason = "คำขอยังคงค้างอยู่และไม่ได้รับการพิจารณาภายในเวลาที่กำหนด"

// tcEmailEnv คือ service กับตัวส่ง Email ที่ใช้นาฬิกาจำลองเดียวกัน
type tcEmailEnv struct {
	tx     pgx.Tx
	clock  *fakeClock
	svc    *Service
	mailer *fakeMailer
	sender *Sender
	srv    http.Handler
}

func tcNewEmailEnv(t *testing.T, tx pgx.Tx, at time.Time) *tcEmailEnv {
	t.Helper()
	clock := &fakeClock{t: at}
	svc := NewService(NewRepository(tx))
	svc.now = clock.Now
	mailer := &fakeMailer{clock: clock}
	return &tcEmailEnv{
		tx: tx, clock: clock, svc: svc, mailer: mailer,
		sender: &Sender{db: tx, mailer: mailer, resolver: MockRecipientResolver{}, now: clock.Now},
		srv:    tcServer(tx, clock),
	}
}

func (e *tcEmailEnv) sendDue(t *testing.T) {
	t.Helper()
	if err := e.sender.SendDue(context.Background()); err != nil {
		t.Fatal(err)
	}
}

// tcThaiDateTime / tcThaiDate เขียนวันเวลาแบบ พ.ศ. ตามเวลาไทยเอง (ไม่ใช้ฟังก์ชันของโค้ดที่ทดสอบ)
func tcThaiDateTime(tm time.Time) string {
	lt := localTime(tm)
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d", lt.Year()+543, lt.Month(), lt.Day(), lt.Hour(), lt.Minute())
}

func tcThaiDate(tm time.Time) string {
	lt := localTime(tm)
	return fmt.Sprintf("%d-%02d-%02d", lt.Year()+543, lt.Month(), lt.Day())
}

func tcEndAt(t *testing.T, tx pgx.Tx, id int) time.Time {
	t.Helper()
	var at time.Time
	if err := tx.QueryRow(context.Background(), `SELECT end_at FROM borrow_request WHERE id = $1`, id).Scan(&at); err != nil {
		t.Fatal(err)
	}
	return at
}

func tcEmailRow(t *testing.T, tx pgx.Tx, id int) (status string, attempts int, lastError string, lastAttempt time.Time) {
	t.Helper()
	var le *string
	var la *time.Time
	if err := tx.QueryRow(context.Background(),
		`SELECT status, attempts, last_error, last_attempt_at FROM approval_email WHERE request_id = $1`, id).Scan(&status, &attempts, &le, &la); err != nil {
		t.Fatal(err)
	}
	if le != nil {
		lastError = *le
	}
	if la != nil {
		lastAttempt = *la
	}
	return
}

func tcEmailID(t *testing.T, tx pgx.Tx, id int) int {
	t.Helper()
	var eid int
	if err := tx.QueryRow(context.Background(), `SELECT id FROM approval_email WHERE request_id = $1`, id).Scan(&eid); err != nil {
		t.Fatal(err)
	}
	return eid
}

func tcAuditCount(t *testing.T, tx pgx.Tx, id int, action string) int {
	t.Helper()
	var n int
	if err := tx.QueryRow(context.Background(),
		`SELECT count(*) FROM audit_log WHERE request_id = $1 AND action = $2`, id, action).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func tcContainsAll(t *testing.T, part, body string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(body, w) {
			t.Errorf("%s เนื้อ Email ต้องมี %q\n%s", part, w, body)
		}
	}
}

// tcMessageOf คืนข้อความที่ handler ตอบสำหรับ error ของ service (FR-APV-01 ซ)
func tcMessageOf(err error) (int, ErrorDTO) {
	w := httptest.NewRecorder()
	writeServiceError(w, err)
	var e ErrorDTO
	_ = json.Unmarshal(w.Body.Bytes(), &e)
	return w.Code, e
}

// tcMidnight คือ 00:00 น. หลังวันที่ D (D = วันนี้ตามเวลาไทย)
func tcMidnight() time.Time { return autoRejectCutoff(time.Now()).AddDate(0, 0, 1) }

// ---------- AC-APV-04 ----------

// TC-APV-04-1 (AC-APV-04 ปกติ)
func TestTC_APV_04_1_ApprovedEmailContent(t *testing.T) {
	// Given staff-001 อนุมัติคำขอของ student-001 บน AP-TEST-0002 สำเร็จที่เวลา T ตัวส่งจำลองรับได้ทุกฉบับ
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	T := tcBeforeCutoff(t, tx, id).Now()
	env := tcNewEmailEnv(t, tx, T)
	if _, err := env.svc.Approve(staffCtx("staff-001", 1), id, time.Time{}); err != nil {
		t.Fatal(err)
	}

	// When ตัวส่ง Email ทำงานที่ T + ApprovalSenderPollInterval
	env.clock.Advance(config.ApprovalSenderPollInterval)
	env.sendDue(t)

	// Then (1) ได้ Email 1 ฉบับ ภายใน 60 วินาทีหลัง T
	if len(env.mailer.sent) != 1 {
		t.Fatalf("(1) ต้องได้ Email 1 ฉบับ แต่ได้ %d", len(env.mailer.sent))
	}
	if env.clock.Now().Sub(T) > 60*time.Second {
		t.Errorf("(1) ต้องส่งภายใน 60 วินาที แต่ส่งที่ %v หลังบันทึกผล", env.clock.Now().Sub(T))
	}
	m := env.mailer.sent[0]
	// Then (2) หัวข้อ
	if m.Subject != "แจ้งผลการอนุมัติการยืมอุปกรณ์ภาคคอมพิวเตอร์" {
		t.Errorf("(2) หัวข้อได้ %q", m.Subject)
	}
	// Then (3) เนื้อหาครบตาม FR-APV-02 (ก)
	tcContainsAll(t, "(3)", m.Body,
		"นักศึกษาตัวอย่าง ก", "อุปกรณ์ทดลองอนุมัติ (004) (AP-TEST-0002)", "เพื่อการเรียน",
		tcThaiDateTime(tcStartAt(t, tx, id)), tcThaiDateTime(tcEndAt(t, tx, id)),
		"คณะวิทยาศาสตร์ ตึก 1 ชั้น 2 ห้องอุปกรณ์ภาควิชาคอมพิวเตอร์")
	// Then (4) ไม่มีจำนวนอุปกรณ์
	if strings.Contains(m.Body, "จำนวน") {
		t.Errorf("(4) Email ผลอนุมัติยืมปกติต้องไม่มีจำนวนอุปกรณ์\n%s", m.Body)
	}
	// Then (5) มีหัวข้อ "ชื่อผู้อนุมัติ" ค่าจะเป็นชื่อจริงหรือไม่ยังไม่ตรวจ (รอ Q-APV-04)
	if !strings.Contains(m.Body, "ชื่อผู้อนุมัติ") {
		t.Errorf("(5) ต้องมีหัวข้อ ชื่อผู้อนุมัติ\n%s", m.Body)
	}
	// Then (6) ที่อยู่ผู้รับ: ยังไม่ตรวจ (รอ Q-APV-04 ตอนนี้ใช้ที่อยู่จำลอง ASM-10)
	// Then (7) approval_email เป็น sent
	if st, _, _, _ := tcEmailRow(t, tx, id); st != "sent" {
		t.Errorf("(7) approval_email ต้องเป็น sent แต่ได้ %q", st)
	}
}

// TC-APV-04-2 (AC-APV-04 ปกติ)
func TestTC_APV_04_2_RejectedEmailContent(t *testing.T) {
	// Given staff-001 ปฏิเสธคำขอของ student-002 บน AP-TEST-0003 ด้วยเหตุผล "อุปกรณ์ต้องใช้สอบ" ที่เวลา T
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	T := time.Now()
	env := tcNewEmailEnv(t, tx, T)
	if _, err := env.svc.Reject(staffCtx("staff-001", 1), id, "อุปกรณ์ต้องใช้สอบ", time.Time{}); err != nil {
		t.Fatal(err)
	}

	// When ตัวส่ง Email ทำงานที่ T + 5 วินาที
	env.clock.Advance(5 * time.Second)
	env.sendDue(t)

	// Then (1) ได้ Email 1 ฉบับภายใน 60 วินาทีหลัง T
	if len(env.mailer.sent) != 1 {
		t.Fatalf("(1) ต้องได้ Email 1 ฉบับ แต่ได้ %d", len(env.mailer.sent))
	}
	m := env.mailer.sent[0]
	// Then (2) หัวข้อ
	if m.Subject != "แจ้งผลการปฏิเสธการยืมอุปกรณ์ภาคคอมพิวเตอร์" {
		t.Errorf("(2) หัวข้อได้ %q", m.Subject)
	}
	// Then (3) เนื้อหาครบตาม FR-APV-02 (ข)
	var submitted time.Time
	if err := tx.QueryRow(context.Background(), `SELECT submitted_at FROM borrow_request WHERE id = $1`, id).Scan(&submitted); err != nil {
		t.Fatal(err)
	}
	tcContainsAll(t, "(3)", m.Body, "นักศึกษาตัวอย่าง ข", "อุปกรณ์ทดลองอนุมัติ (004) (AP-TEST-0003)",
		tcThaiDate(submitted), "อุปกรณ์ต้องใช้สอบ")
	// Then (4) มีหัวข้อ "ชื่อผู้พิจารณา" ค่ายังไม่ตรวจ (รอ Q-APV-04)
	if !strings.Contains(m.Body, "ชื่อผู้พิจารณา") {
		t.Errorf("(4) ต้องมีหัวข้อ ชื่อผู้พิจารณา\n%s", m.Body)
	}
}

// TC-APV-04-3 (AC-APV-04 ขอบ)
func TestTC_APV_04_3_SendWithinDeadline(t *testing.T) {
	// Given staff-001 อนุมัติคำขอของ AP-TEST-0002 ที่เวลา T
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	env := tcNewEmailEnv(t, tx, tcBeforeCutoff(t, tx, id).Now())
	if _, err := env.svc.Approve(staffCtx("staff-001", 1), id, time.Time{}); err != nil {
		t.Fatal(err)
	}

	// When ตัวส่ง Email ทำงานที่ T พอดี (0 วินาที)
	env.sendDue(t)

	// Then (1) Email ส่งได้ทันทีในรอบนั้น
	if len(env.mailer.sent) != 1 {
		t.Errorf("(1) ต้องส่งได้ทันที แต่ได้ %d ฉบับ", len(env.mailer.sent))
	}
	// Then (2) ApprovalSenderPollInterval น้อยกว่า ApprovalEmailDeadline (60 วินาที)
	if !(config.ApprovalSenderPollInterval < config.ApprovalEmailDeadline) || config.ApprovalEmailDeadline != 60*time.Second {
		t.Errorf("(2) รอบตรวจ %v ต้องน้อยกว่า %v (60 วินาที)", config.ApprovalSenderPollInterval, config.ApprovalEmailDeadline)
	}
}

// TC-APV-04-4 (AC-APV-04 ขอบ)
func TestTC_APV_04_4_PurposeOtherText(t *testing.T) {
	// Given คำขอของ AP-TEST-0002 วัตถุประสงค์ "อื่น ๆ" พิมพ์ว่า "ถ่ายวิดีโอโครงงาน" staff-001 อนุมัติ
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	if _, err := tx.Exec(context.Background(),
		`UPDATE borrow_request SET purpose = $2, purpose_other = 'ถ่ายวิดีโอโครงงาน' WHERE id = $1`, id, config.PurposeOther); err != nil {
		t.Fatal(err)
	}
	env := tcNewEmailEnv(t, tx, tcBeforeCutoff(t, tx, id).Now())
	if _, err := env.svc.Approve(staffCtx("staff-001", 1), id, time.Time{}); err != nil {
		t.Fatal(err)
	}

	// When ตัวส่ง Email ทำงาน
	env.sendDue(t)

	// Then (1) วัตถุประสงค์ใน Email เป็น "ถ่ายวิดีโอโครงงาน" ไม่ใช่ "อื่น ๆ"
	if len(env.mailer.sent) != 1 {
		t.Fatalf("ต้องได้ Email 1 ฉบับ แต่ได้ %d", len(env.mailer.sent))
	}
	body := env.mailer.sent[0].Body
	if !strings.Contains(body, "ถ่ายวิดีโอโครงงาน") || strings.Contains(body, "อื่น ๆ") {
		t.Errorf("(1) วัตถุประสงค์ต้องเป็นข้อความที่พิมพ์\n%s", body)
	}
}

// TC-APV-04-5 (AC-APV-04 ขอบ)
func TestTC_APV_04_5_SupersededEmail(t *testing.T) {
	// Given AP-TEST-0001 มี A และ B staff-001 อนุมัติ A (B ถูกปฏิเสธตาม)
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	env := tcNewEmailEnv(t, tx, tcBeforeCutoff(t, tx, a).Now())
	if _, err := env.svc.Approve(staffCtx("staff-001", 1), a, time.Time{}); err != nil {
		t.Fatal(err)
	}

	// When ตัวส่ง Email ทำงาน
	env.sendDue(t)

	// Then (1) ได้ Email ผลอนุมัติถึงผู้ขอ A และผลปฏิเสธถึงผู้ขอ B
	var approved, rejectedB *Email
	for i := range env.mailer.sent {
		m := &env.mailer.sent[i]
		switch {
		case m.Subject == "แจ้งผลการอนุมัติการยืมอุปกรณ์ภาคคอมพิวเตอร์" && strings.Contains(m.Body, "นักศึกษาตัวอย่าง ก"):
			approved = m
		case m.Subject == "แจ้งผลการปฏิเสธการยืมอุปกรณ์ภาคคอมพิวเตอร์" && strings.Contains(m.Body, "นักศึกษาตัวอย่าง ข"):
			rejectedB = m
		}
	}
	if approved == nil || rejectedB == nil {
		t.Fatalf("(1) ต้องได้ Email ผลอนุมัติของ A และผลปฏิเสธของ B แต่ได้ %d ฉบับ", len(env.mailer.sent))
	}
	// Then (2) Email ของ B มีเหตุผล "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว"
	if !strings.Contains(rejectedB.Body, tcMsgSuperseded) {
		t.Errorf("(2) Email ของ B ต้องมีเหตุผล %q\n%s", tcMsgSuperseded, rejectedB.Body)
	}
}

// TC-APV-04-6 (AC-APV-04 ผิด)
func TestTC_APV_04_6_NoEmailWhenNotSaved(t *testing.T) {
	// Given staff-001 อนุมัติ B ของ AP-TEST-0001 (ไม่ใช่คำขอแรก ได้ 409)
	tx := devTx(t)
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	env := tcNewEmailEnv(t, tx, tcBeforeCutoff(t, tx, b).Now())
	if w := tcApprove(env.srv, tokenStaffDept1, b); w.Code != http.StatusConflict {
		t.Fatalf("เตรียมข้อมูล: ต้องได้ 409 แต่ได้ %d", w.Code)
	}

	// When ตัวส่ง Email ทำงาน
	env.sendDue(t)

	// Then (1) ไม่มีแถว approval_email ของ B
	if n := tcEmailCount(t, tx, b); n != 0 {
		t.Errorf("(1) ต้องไม่มี approval_email ของ B แต่มี %d", n)
	}
	// Then (2) ตัวส่งจำลองไม่ได้ Email ใด
	if len(env.mailer.sent) != 0 {
		t.Errorf("(2) ต้องไม่มี Email แต่ได้ %d ฉบับ", len(env.mailer.sent))
	}
}

// ---------- AC-APV-10 ----------

// tcConcurrentRound สร้างคำขอใหม่ (commit จริง) บนอุปกรณ์ของคำขอ like แล้วให้ 2 action ทำพร้อมกัน
func tcConcurrentRound(t *testing.T, pool *pgxpool.Pool, like int, requester string, acts [2]func(id int) error) (int, [2]error) {
	t.Helper()
	ctx := context.Background()
	// คำขอ requested เดิมของชิ้นนี้ยกเลิกก่อน คำขอใหม่จึงเป็นคำขอแรก (เตรียมข้อมูล ไม่ใช่สิ่งที่ตรวจ)
	if _, err := pool.Exec(ctx, `
		UPDATE borrow_request SET status = 'cancelled', cancel_reason = 'เตรียมรอบใหม่', cancelled_at = now()
		WHERE status = 'requested' AND item_id = (SELECT item_id FROM borrow_request WHERE id = $1)`, like); err != nil {
		t.Fatal(err)
	}
	var id int
	if err := pool.QueryRow(ctx, `
		INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose,
		                            start_at, end_at, status, submitted_at)
		SELECT $2, 'ผู้ขอ '||$2, $2, r.item_id, 'normal', 'study', r.start_at, r.end_at, 'requested', now() - interval '1 hour'
		FROM borrow_request r WHERE r.id = $1 RETURNING id`, like, requester).Scan(&id); err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	var errs [2]error
	var wg sync.WaitGroup
	for i := range acts {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			errs[i] = acts[i](id)
		}(i)
	}
	close(start)
	wg.Wait()
	return id, errs
}

func tcPoolIDOf(t *testing.T, pool *pgxpool.Pool, asset, requester string) int {
	t.Helper()
	var id int
	if err := pool.QueryRow(context.Background(), `
		SELECT r.id FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id
		WHERE i.asset_code = $1 AND r.requester_id = $2`, asset, requester).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

// tcCheckOneWinner ตรวจว่าสำเร็จ 1 อีกคนได้ ErrDecidedConcurrently และผลในฐานข้อมูลตรงกัน
func tcCheckOneWinner(t *testing.T, pool *pgxpool.Pool, round, id int, errs [2]error) {
	t.Helper()
	ok, conc := 0, 0
	for _, err := range errs {
		switch {
		case err == nil:
			ok++
		case errors.Is(err, ErrDecidedConcurrently):
			conc++
		default:
			t.Errorf("รอบ %d ได้ error ที่ไม่คาด: %v", round, err)
		}
	}
	// Then (1) สำเร็จ 1 คน อีกคนได้ REQUEST_DECIDED_CONCURRENTLY
	if ok != 1 || conc != 1 {
		t.Errorf("(1) รอบ %d ต้องสำเร็จ 1 และ DECIDED_CONCURRENTLY 1 แต่ได้ %d / %d (%v)", round, ok, conc, errs)
	}
	// Then (2) approval_decision แถวเดียว และสถานะตรงกับผล
	var status, result string
	var n int
	if err := pool.QueryRow(context.Background(), `
		SELECT r.status, min(d.result), count(d.*) FROM borrow_request r LEFT JOIN approval_decision d ON d.request_id = r.id
		WHERE r.id = $1 GROUP BY r.status`, id).Scan(&status, &result, &n); err != nil {
		t.Fatal(err)
	}
	if n != 1 || status != result {
		t.Errorf("(2) รอบ %d ต้องมี approval_decision แถวเดียวและตรงกับสถานะ แต่ได้ %d แถว สถานะ %q ผล %q", round, n, status, result)
	}
}

// TC-APV-10-1 (AC-APV-10 ปกติ)
func TestTC_APV_10_1_ConcurrentApprove(t *testing.T) {
	// Given staff-001 และ staff-002 (ภาควิชา 1) โหลดรายการแล้ว คำขอเป็น requested (commit จริง)
	pool := concurrencyDB(t)
	svc := NewService(NewRepository(pool))
	loadedAt := time.Now().Add(-time.Hour)
	like := tcPoolIDOf(t, pool, "AP-TEST-0002", "student-001")
	approveAs := func(staffID string) func(int) error {
		return func(id int) error { _, err := svc.Approve(staffCtx(staffID, 1), id, loadedAt); return err }
	}

	// When ทั้งสองกดอนุมัติพร้อมกัน ทำซ้ำ 20 รอบ
	for round := 1; round <= 20; round++ {
		id, errs := tcConcurrentRound(t, pool, like, fmt.Sprintf("student-c1-%d", round),
			[2]func(int) error{approveAs("staff-001"), approveAs("staff-002")})
		// Then (1)(2)
		tcCheckOneWinner(t, pool, round, id, errs)
		if s := func() string {
			var s string
			_ = pool.QueryRow(context.Background(), `SELECT status FROM borrow_request WHERE id = $1`, id).Scan(&s)
			return s
		}(); s != "approved" {
			t.Errorf("(2) รอบ %d สถานะต้องเป็น approved แต่ได้ %q", round, s)
		}
	}

	// Then (3) ข้อความของคนที่ช้ากว่า
	code, e := tcMessageOf(ErrDecidedConcurrently)
	if code != http.StatusConflict || e.Code != "REQUEST_DECIDED_CONCURRENTLY" || e.Message != tcMsgConcurrent {
		t.Errorf("(3) ต้องได้ 409 REQUEST_DECIDED_CONCURRENTLY %q แต่ได้ %d %+v", tcMsgConcurrent, code, e)
	}
}

// TC-APV-10-2 (AC-APV-10 ปกติ)
func TestTC_APV_10_2_ConcurrentApproveReject(t *testing.T) {
	// Given เหมือน TC-APV-10-1
	pool := concurrencyDB(t)
	svc := NewService(NewRepository(pool))
	loadedAt := time.Now().Add(-time.Hour)
	like := tcPoolIDOf(t, pool, "AP-TEST-0003", "student-002")
	approve := func(id int) error { _, err := svc.Approve(staffCtx("staff-001", 1), id, loadedAt); return err }
	reject := func(id int) error {
		_, err := svc.Reject(staffCtx("staff-002", 1), id, "กดพร้อมกัน", loadedAt)
		return err
	}

	// When คนหนึ่งกดอนุมัติ อีกคนกดปฏิเสธพร้อมเหตุผล พร้อมกัน ทำซ้ำ 20 รอบ
	for round := 1; round <= 20; round++ {
		id, errs := tcConcurrentRound(t, pool, like, fmt.Sprintf("student-c2-%d", round), [2]func(int) error{approve, reject})
		// Then (1)(2)
		tcCheckOneWinner(t, pool, round, id, errs)
	}
}

// TC-APV-10-4 (AC-APV-10 ขอบ)
func TestTC_APV_10_4_DecidedAfterLoad(t *testing.T) {
	// Given เจ้าหน้าที่คนที่ 2 โหลดรายการที่ T0 คนที่ 1 อนุมัติที่ T1 (หลัง T0)
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	svc := NewService(NewRepository(tx))
	t0 := time.Now().Add(-time.Hour)
	if _, err := svc.Approve(staffCtx("staff-001", 1), id, time.Time{}); err != nil {
		t.Fatal(err)
	}

	// When คนที่ 2 กดปฏิเสธด้วย loadedAt = T0
	_, err := svc.Reject(staffCtx("staff-002", 1), id, "อุปกรณ์ต้องใช้สอบ", t0)

	// Then (1) ได้ 409 REQUEST_DECIDED_CONCURRENTLY ไม่ใช่ REQUEST_ALREADY_DECIDED
	code, e := tcMessageOf(err)
	if code != http.StatusConflict || e.Code != "REQUEST_DECIDED_CONCURRENTLY" {
		t.Errorf("(1) ต้องได้ 409 REQUEST_DECIDED_CONCURRENTLY แต่ได้ %d %+v (err %v)", code, e, err)
	}
	// Then (2) คำขอยังเป็น approved
	if s := tcRequestStatus(t, tx, id); s != "approved" {
		t.Errorf("(2) คำขอต้องยังเป็น approved แต่ได้ %q", s)
	}
}

// TC-APV-10-5 (AC-APV-10 ผิด)
func TestTC_APV_10_5_ConcurrentOtherDepartment(t *testing.T) {
	// Given staff-001 (ภาควิชา 1) กับเจ้าหน้าที่ภาควิชา 2 และคำขอของ AP-TEST-0002 (commit จริง)
	pool := concurrencyDB(t)
	svc := NewService(NewRepository(pool))
	id := tcPoolIDOf(t, pool, "AP-TEST-0002", "student-001")

	// When ทั้งสองกดอนุมัติพร้อมกัน
	start := make(chan struct{})
	var errDept1, errDept2 error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); <-start; _, errDept1 = svc.Approve(staffCtx("staff-001", 1), id, time.Time{}) }()
	go func() { defer wg.Done(); <-start; _, errDept2 = svc.Approve(staffCtx("staff-201", 2), id, time.Time{}) }()
	close(start)
	wg.Wait()

	// Then (1) staff-001 สำเร็จ
	if errDept1 != nil {
		t.Errorf("(1) staff-001 ต้องสำเร็จ แต่ได้ %v", errDept1)
	}
	// Then (2) เจ้าหน้าที่ภาควิชา 2 ได้ FORBIDDEN_DEPARTMENT ไม่ใช่ REQUEST_DECIDED_CONCURRENTLY
	if !errors.Is(errDept2, ErrForbiddenDepartment) {
		t.Errorf("(2) ภาควิชา 2 ต้องได้ ErrForbiddenDepartment แต่ได้ %v", errDept2)
	}
	// Then (3) มี approval_decision แถวเดียว
	var n int
	if err := pool.QueryRow(context.Background(), `SELECT count(*) FROM approval_decision WHERE request_id = $1`, id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("(3) ต้องมี approval_decision แถวเดียว แต่มี %d", n)
	}
}

// ---------- AC-APV-13 ----------

// tcRunMidnight รันงานปฏิเสธอัตโนมัติด้วยนาฬิกาที่ระบุ
func tcRunMidnight(t *testing.T, tx pgx.Tx, at time.Time) {
	t.Helper()
	job := &AutoRejecter{db: tx, now: func() time.Time { return at }}
	if _, err := job.RejectStale(context.Background()); err != nil {
		t.Fatal(err)
	}
}

// TC-APV-13-1 (AC-APV-13 ปกติ)
func TestTC_APV_13_1_MidnightAutoReject(t *testing.T) {
	// Given คำขอยืมปกติของ AP-TEST-0003 ส่งวันที่ D เวลา 21:00 น. ยังเป็น requested
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight.Add(-3*time.Hour))

	// When งานปฏิเสธอัตโนมัติทำงานที่ 00:00 น. หลังวันที่ D
	tcRunMidnight(t, tx, midnight)

	// Then (1) คำขอเป็น rejected
	if s := tcRequestStatus(t, tx, id); s != "rejected" {
		t.Errorf("(1) คำขอต้องเป็น rejected แต่ได้ %q", s)
	}
	// Then (2) approval_decision ผล auto_rejected ผู้ดำเนินการ system ชื่อ System (3) เหตุผลตามที่กำหนด
	var result, reason, actorID, actorName string
	if err := tx.QueryRow(context.Background(),
		`SELECT result, reason, actor_id, actor_name FROM approval_decision WHERE request_id = $1`, id).Scan(&result, &reason, &actorID, &actorName); err != nil {
		t.Fatalf("(2) ต้องมี approval_decision: %v", err)
	}
	if result != "auto_rejected" || actorID != "system" || actorName != "System" {
		t.Errorf("(2) ต้องเป็น auto_rejected โดย system/System แต่ได้ %q โดย %q/%q", result, actorID, actorName)
	}
	if reason != tcAutoReason {
		t.Errorf("(3) เหตุผลต้องเป็น %q แต่ได้ %q", tcAutoReason, reason)
	}
	// Then (4) มีแถว approval_email แบบ auto_rejected
	if k := tcEmailKind(t, tx, id); k != "auto_rejected" {
		t.Errorf("(4) Email ต้องเป็นแบบ auto_rejected แต่ได้ %q", k)
	}
}

// TC-APV-13-2 (AC-APV-13 ปกติ)
func TestTC_APV_13_2_AutoRejectEmailContent(t *testing.T) {
	// Given ต่อจาก TC-APV-13-1
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	submitted := midnight.Add(-3 * time.Hour)
	setSubmittedAt(t, tx, id, submitted)
	tcRunMidnight(t, tx, midnight)

	// When ตัวส่ง Email ทำงานที่ 00:00:05 น.
	env := tcNewEmailEnv(t, tx, midnight.Add(5*time.Second))
	env.sendDue(t)

	// Then (1) ได้ Email 1 ฉบับภายใน 60 วินาทีหลัง 00:00 น.
	var mine []Email
	for _, m := range env.mailer.sent {
		if strings.Contains(m.Body, fmt.Sprint(id)) && strings.Contains(m.Body, "AP-TEST-0003") {
			mine = append(mine, m)
		}
	}
	if len(mine) != 1 {
		t.Fatalf("(1) ต้องได้ Email ของคำขอนี้ 1 ฉบับ แต่ได้ %d", len(mine))
	}
	// Then (2) หัวข้อ
	if mine[0].Subject != "แจ้งปฏิเสธคำขอยืมอุปกรณ์อัตโนมัติ" {
		t.Errorf("(2) หัวข้อได้ %q", mine[0].Subject)
	}
	// Then (3) เนื้อหาครบตาม FR-APV-02 (ค)
	tcContainsAll(t, "(3)", mine[0].Body, fmt.Sprint(id), "นักศึกษาตัวอย่าง ข", "อุปกรณ์ทดลองอนุมัติ (004) (AP-TEST-0003)",
		tcThaiDate(submitted), "ปฏิเสธอัตโนมัติ", tcAutoReason, "00:00 น.", "หากยังต้องการยืมอุปกรณ์ สามารถยื่นคำขอใหม่ได้")
}

// TC-APV-13-3 (AC-APV-13 ขอบ)
func TestTC_APV_13_3_SubmittedJustBeforeMidnight(t *testing.T) {
	// Given คำขอส่งวันที่ D เวลา 23:59:59 น.
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight.Add(-time.Second))

	// When งานปฏิเสธอัตโนมัติทำงานที่ 00:00 น. หลังวันที่ D
	tcRunMidnight(t, tx, midnight)

	// Then (1) คำขอเป็น rejected ผล auto_rejected
	if res, _, _, _, _ := tcDecision(t, tx, id); res != "auto_rejected" || tcRequestStatus(t, tx, id) != "rejected" {
		t.Errorf("(1) ต้องถูกปฏิเสธอัตโนมัติ แต่ได้ผล %q สถานะ %q", res, tcRequestStatus(t, tx, id))
	}
}

// TC-APV-13-4 (AC-APV-13 ขอบ)
func TestTC_APV_13_4_SubmittedAtMidnight(t *testing.T) {
	// Given คำขอส่งวันที่ D+1 เวลา 00:00:00 น. พอดี
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight)

	// When งานปฏิเสธอัตโนมัติทำงานที่ 00:00 น. หลังวันที่ D
	tcRunMidnight(t, tx, midnight)

	// Then (1) คำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(1) คำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}

// TC-APV-13-5 (AC-APV-13 ขอบ)
func TestTC_APV_13_5_CountFromSubmitDate(t *testing.T) {
	// Given คำขอส่งวันที่ D แต่วันที่เริ่มยืมเป็น D+2
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight.Add(-3*time.Hour))
	if _, err := tx.Exec(context.Background(), `
		UPDATE borrow_request SET start_at = $2, end_at = $3 WHERE id = $1`,
		id, midnight.Add(24*time.Hour+9*time.Hour), midnight.Add(24*time.Hour+16*time.Hour)); err != nil {
		t.Fatal(err)
	}

	// When งานปฏิเสธอัตโนมัติทำงานที่ 00:00 น. หลังวันที่ D
	tcRunMidnight(t, tx, midnight)

	// Then (1) คำขอเป็น rejected ผล auto_rejected
	if res, _, _, _, _ := tcDecision(t, tx, id); res != "auto_rejected" {
		t.Errorf("(1) ต้องถูกปฏิเสธอัตโนมัติ แต่ได้ผล %q", res)
	}
}

// TC-APV-13-6 (AC-APV-13 ผิด)
func TestTC_APV_13_6_OnlyRequestedRejected(t *testing.T) {
	// Given คำขอที่ส่งวันที่ D แต่ไม่ใช่ requested แล้ว: อนุมัติแล้ว (0002), ปฏิเสธแล้ว (0003), ยกเลิกแล้ว (0005)
	tx := devTx(t)
	approved := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	rejected := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	cancelled := requestIDOf(t, tx, "AP-TEST-0005", "student-002")
	svc := NewService(NewRepository(tx))
	svc.now = tcBeforeCutoff(t, tx, approved).Now
	if _, err := svc.Approve(staffCtx("staff-001", 1), approved, time.Time{}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Reject(staffCtx("staff-001", 1), rejected, "อุปกรณ์ต้องใช้สอบ", time.Time{}); err != nil {
		t.Fatal(err)
	}
	midnight := tcMidnight()
	want := map[int]string{approved: "approved", rejected: "rejected", cancelled: "cancelled"}
	for id := range want {
		setSubmittedAt(t, tx, id, midnight.Add(-3*time.Hour))
	}

	// When งานปฏิเสธอัตโนมัติทำงานที่ 00:00 น. หลังวันที่ D
	tcRunMidnight(t, tx, midnight)

	for id, st := range want {
		// Then (1) ทั้ง 3 คำขอคงสถานะเดิม
		if s := tcRequestStatus(t, tx, id); s != st {
			t.Errorf("(1) คำขอ %d ต้องยังเป็น %s แต่ได้ %q", id, st, s)
		}
		// Then (2) ไม่มีผล auto_rejected ของคำขอเหล่านี้
		if res, _, _, _, _ := tcDecision(t, tx, id); res == "auto_rejected" {
			t.Errorf("(2) คำขอ %d ต้องไม่มีผล auto_rejected", id)
		}
	}
}

// ---------- AC-APV-14 ----------

// tcFailAll อนุมัติคำขอของ AP-TEST-0002 ที่เวลา T แล้วให้ตัวส่งล้มที่ T, T+1, ... (rounds ครั้ง)
func tcFailRounds(t *testing.T, env *tcEmailEnv, rounds int) {
	t.Helper()
	for i := 0; i < rounds; i++ {
		if i > 0 {
			env.clock.Advance(config.ApprovalEmailRetryInterval)
		}
		env.sendDue(t)
	}
}

func tcApprovedWithFailingMailer(t *testing.T) (*tcEmailEnv, int, time.Time) {
	t.Helper()
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	T := tcBeforeCutoff(t, tx, id).Now().Add(-10 * time.Minute)
	env := tcNewEmailEnv(t, tx, T)
	if _, err := env.svc.Approve(staffCtx("staff-001", 1), id, time.Time{}); err != nil {
		t.Fatal(err)
	}
	env.mailer.fail = errors.New("smtp down")
	return env, id, T
}

// TC-APV-14-1 (AC-APV-14 ปกติ)
func TestTC_APV_14_1_RetryThenFailed(t *testing.T) {
	// Given staff-001 อนุมัติคำขอของ AP-TEST-0002 ที่เวลา T ตัวส่งล้มทุกครั้งด้วยสาเหตุ "smtp down"
	env, id, T := tcApprovedWithFailingMailer(t)

	// When ตัวส่งทำงานที่ T, T+1 นาที, T+2 นาที, T+3 นาที
	tcFailRounds(t, env, 4)

	// Then (1) ตัวส่งจำลองถูกเรียก 4 ครั้ง ห่างกันครั้งละ 1 นาที
	if len(env.mailer.calls) != 4 {
		t.Fatalf("(1) ต้องถูกเรียก 4 ครั้ง แต่ได้ %d", len(env.mailer.calls))
	}
	for i := 1; i < 4; i++ {
		if d := env.mailer.calls[i].Sub(env.mailer.calls[i-1]); d != time.Minute {
			t.Errorf("(1) ครั้งที่ %d ห่างจากครั้งก่อน %v ต้องเป็น 1 นาที", i+1, d)
		}
	}
	// Then (2) failed attempts 4 last_error "smtp down" last_attempt_at = T+3 นาที
	st, attempts, lastErr, lastAt := tcEmailRow(t, env.tx, id)
	if st != "failed" || attempts != 4 || lastErr != "smtp down" || !lastAt.Equal(T.Add(3*time.Minute)) {
		t.Errorf("(2) ได้ %q attempts %d %q %v", st, attempts, lastErr, lastAt)
	}
	// Then (3) คำขอยังเป็น approved และ approval_decision ไม่เปลี่ยน
	if s := tcRequestStatus(t, env.tx, id); s != "approved" {
		t.Errorf("(3) คำขอต้องยังเป็น approved แต่ได้ %q", s)
	}
	if res, _, _, _, ok := tcDecision(t, env.tx, id); !ok || res != "approved" || tcDecisionCount(t, env.tx, id) != 1 {
		t.Errorf("(3) approval_decision ต้องเป็น approved แถวเดียว แต่ได้ %q", res)
	}
}

// TC-APV-14-2 (AC-APV-14 ปกติ)
func TestTC_APV_14_2_FailedEmailListed(t *testing.T) {
	// Given ต่อจาก TC-APV-14-1
	env, id, _ := tcApprovedWithFailingMailer(t)
	tcFailRounds(t, env, 4)

	// When staff-001 เรียก GET /v1/approvals/failed-emails
	w := getJSON(env.srv, tokenStaffDept1, PathPrefix+"failed-emails")

	// Then (1) ได้ 200 และมีรายการของคำขอนั้น
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d", w.Code)
	}
	var got FailedEmailListDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	var item *FailedEmailDTO
	for i := range got.Items {
		if got.Items[i].RequestID == id {
			item = &got.Items[i]
		}
	}
	if item == nil {
		t.Fatalf("(1) ต้องมีรายการของคำขอ %d แต่ได้ %+v", id, got.Items)
	}
	// Then (2) emailId, requestId, แบบ approved, ชื่อผู้ขอ, ชื่อรุ่น, รหัสครุภัณฑ์, สาเหตุ และเวลา
	if item.EmailID == 0 || item.Kind != "approved" || item.RequesterName != "นักศึกษาตัวอย่าง ก" ||
		item.ModelName != "อุปกรณ์ทดลองอนุมัติ (004)" || item.AssetCode != "AP-TEST-0002" ||
		item.LastError != "smtp down" || item.LastAttemptAt.IsZero() {
		t.Errorf("(2) ข้อมูลไม่ครบ: %+v", item)
	}
}

// TC-APV-14-3 (AC-APV-14 ปกติ)
func TestTC_APV_14_3_ResendSucceeds(t *testing.T) {
	// Given ต่อจาก TC-APV-14-1 ตัวส่งกลับมาปกติ
	env, id, _ := tcApprovedWithFailingMailer(t)
	tcFailRounds(t, env, 4)
	env.mailer.fail = nil
	emailID := tcEmailID(t, env.tx, id)

	// When staff-001 เรียก POST /v1/approvals/emails/{id}/resend แล้วตัวส่งทำงาน
	w := postJSON(env.srv, tokenStaffDept1, fmt.Sprintf("%semails/%d/resend", PathPrefix, emailID), `{}`)
	env.sendDue(t)

	// Then (1) resend ได้ 200
	if w.Code != http.StatusOK {
		t.Errorf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	// Then (2) ตัวส่งได้ Email 1 ฉบับ และ approval_email เป็น sent
	if len(env.mailer.sent) != 1 {
		t.Errorf("(2) ต้องได้ Email 1 ฉบับ แต่ได้ %d", len(env.mailer.sent))
	}
	if st, _, _, _ := tcEmailRow(t, env.tx, id); st != "sent" {
		t.Errorf("(2) approval_email ต้องเป็น sent แต่ได้ %q", st)
	}
	// Then (3) คำขอยังเป็น approved และ approval_decision ไม่เปลี่ยน
	if s := tcRequestStatus(t, env.tx, id); s != "approved" || tcDecisionCount(t, env.tx, id) != 1 {
		t.Errorf("(3) คำขอต้องยังเป็น approved และมีผลแถวเดียว แต่ได้ %q", s)
	}
	// Then (4) Email นั้นไม่อยู่ใน failed-emails แล้ว
	lw := getJSON(env.srv, tokenStaffDept1, PathPrefix+"failed-emails")
	var got FailedEmailListDTO
	_ = json.Unmarshal(lw.Body.Bytes(), &got)
	for _, it := range got.Items {
		if it.EmailID == emailID {
			t.Error("(4) Email นี้ต้องไม่อยู่ใน failed-emails แล้ว")
		}
	}
}

// TC-APV-14-5 (AC-APV-14 ขอบ)
func TestTC_APV_14_5_NotFailedBeforeLastRetry(t *testing.T) {
	// Given ตัวส่งล้มทุกครั้ง ทำงานแค่ T, T+1, T+2 นาที
	env, id, _ := tcApprovedWithFailingMailer(t)
	tcFailRounds(t, env, 3)

	// When staff-001 เรียก failed-emails
	w := getJSON(env.srv, tokenStaffDept1, PathPrefix+"failed-emails")

	// Then (1) approval_email ยังไม่เป็น failed
	if st, _, _, _ := tcEmailRow(t, env.tx, id); st == "failed" {
		t.Errorf("(1) ยังต้องไม่เป็น failed แต่ได้ %q", st)
	}
	// Then (2) ยังไม่อยู่ในรายการ failed-emails
	var got FailedEmailListDTO
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	for _, it := range got.Items {
		if it.RequestID == id {
			t.Error("(2) ต้องยังไม่อยู่ในรายการ")
		}
	}
}

// TC-APV-14-6 (AC-APV-14 ขอบ)
func TestTC_APV_14_6_SucceedOnRetry(t *testing.T) {
	// Given ตัวส่งล้ม 2 ครั้งแรก แล้วสำเร็จครั้งที่ 3 (T+2 นาที)
	env, id, _ := tcApprovedWithFailingMailer(t)

	// When ตัวส่งทำงานที่ T, T+1, T+2 นาที
	tcFailRounds(t, env, 2)
	env.mailer.fail = nil
	env.clock.Advance(config.ApprovalEmailRetryInterval)
	env.sendDue(t)

	// Then (1) approval_email เป็น sent
	if st, _, _, _ := tcEmailRow(t, env.tx, id); st != "sent" {
		t.Errorf("(1) ต้องเป็น sent แต่ได้ %q", st)
	}
	// Then (2) ไม่อยู่ในรายการ failed-emails
	w := getJSON(env.srv, tokenStaffDept1, PathPrefix+"failed-emails")
	var got FailedEmailListDTO
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	for _, it := range got.Items {
		if it.RequestID == id {
			t.Error("(2) ต้องไม่อยู่ในรายการ")
		}
	}
	// Then (3) ตัวส่งไม่ถูกเรียกอีกที่ T+3 นาที
	before := len(env.mailer.calls)
	env.clock.Advance(config.ApprovalEmailRetryInterval)
	env.sendDue(t)
	if len(env.mailer.calls) != before {
		t.Errorf("(3) ต้องไม่ถูกเรียกอีก แต่ถูกเรียกเพิ่ม %d ครั้ง", len(env.mailer.calls)-before)
	}
}

// TC-APV-14-7 (AC-APV-14 ขอบ)
func TestTC_APV_14_7_SystemEmailListed(t *testing.T) {
	// Given Email ปฏิเสธอัตโนมัติ (System) ของคำขอภาควิชา 1 ส่งไม่สำเร็จครบ 4 ครั้ง
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight.Add(-3*time.Hour))
	tcRunMidnight(t, tx, midnight)
	env := tcNewEmailEnv(t, tx, midnight)
	env.mailer.fail = errors.New("smtp down")
	tcFailRounds(t, env, 4)

	// When staff-001 เรียก failed-emails
	w := getJSON(env.srv, tokenStaffDept1, PathPrefix+"failed-emails")

	// Then (1) มีรายการนั้น แบบ auto_rejected
	var got FailedEmailListDTO
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	found := false
	for _, it := range got.Items {
		if it.RequestID == id && it.Kind == "auto_rejected" {
			found = true
		}
	}
	if !found {
		t.Errorf("(1) ต้องมี Email ปฏิเสธอัตโนมัติของคำขอ %d ในรายการ แต่ได้ %+v", id, got.Items)
	}
}

// TC-APV-14-8 (AC-APV-14 ผิด)
func TestTC_APV_14_8_OtherDepartmentCannotResend(t *testing.T) {
	// Given เจ้าหน้าที่ภาควิชา 2 กับ Email ที่ส่งไม่สำเร็จของ AP-TEST-0002 (ภาควิชา 1)
	env, id, _ := tcApprovedWithFailingMailer(t)
	tcFailRounds(t, env, 4)
	emailID := tcEmailID(t, env.tx, id)
	ctx2 := staffCtx("staff-201", 2)

	// When เรียกรายการ failed-emails และส่งซ้ำ Email นั้นผ่าน service
	items, err := env.svc.ListFailedEmails(ctx2)
	if err != nil {
		t.Fatal(err)
	}
	_, resendErr := env.svc.Resend(ctx2, emailID)

	// Then (1) รายการไม่มี Email นั้น
	for _, it := range items {
		if it.ID == emailID {
			t.Error("(1) ภาควิชา 2 ต้องไม่เห็น Email นี้")
		}
	}
	// Then (2) ส่งซ้ำได้ error ภาควิชาอื่น
	if !errors.Is(resendErr, ErrForbiddenDepartment) {
		t.Errorf("(2) ต้องได้ ErrForbiddenDepartment แต่ได้ %v", resendErr)
	}
	// Then (3) approval_email ยังเป็น failed
	if st, _, _, _ := tcEmailRow(t, env.tx, id); st != "failed" {
		t.Errorf("(3) ต้องยังเป็น failed แต่ได้ %q", st)
	}
}

// TC-APV-14-9 (AC-APV-14 ผิด)
func TestTC_APV_14_9_ResendSentEmail(t *testing.T) {
	// Given Email ที่ส่งสำเร็จแล้ว (sent) ของ AP-TEST-0002
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	env := tcNewEmailEnv(t, tx, tcBeforeCutoff(t, tx, id).Now())
	if _, err := env.svc.Approve(staffCtx("staff-001", 1), id, time.Time{}); err != nil {
		t.Fatal(err)
	}
	env.sendDue(t)
	_, attemptsBefore, _, _ := tcEmailRow(t, tx, id)
	sentBefore := len(env.mailer.sent)

	// When staff-001 เรียก POST /v1/approvals/emails/{id}/resend
	w := postJSON(env.srv, tokenStaffDept1, fmt.Sprintf("%semails/%d/resend", PathPrefix, tcEmailID(t, tx, id)), `{}`)

	// Then (1) ได้ 200 พร้อมสถานะ sent เดิม
	if w.Code != http.StatusOK {
		t.Fatalf("(1) ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var got ResendResultDTO
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got.Status != "sent" {
		t.Errorf("(1) สถานะต้องเป็น sent แต่ได้ %q", got.Status)
	}
	// Then (2) ไม่ส่งใหม่: ตัวส่งไม่ได้ Email เพิ่ม และ attempts ไม่เพิ่ม
	env.clock.Advance(time.Minute)
	env.sendDue(t)
	if len(env.mailer.sent) != sentBefore {
		t.Errorf("(2) ต้องไม่ส่งใหม่ แต่ได้ Email เพิ่ม %d ฉบับ", len(env.mailer.sent)-sentBefore)
	}
	if _, attempts, _, _ := tcEmailRow(t, tx, id); attempts != attemptsBefore {
		t.Errorf("(2) attempts ต้องเท่าเดิม %d แต่ได้ %d", attemptsBefore, attempts)
	}
	// Then (3) approval_decision ไม่เปลี่ยน
	if res, _, _, _, _ := tcDecision(t, tx, id); res != "approved" || tcDecisionCount(t, tx, id) != 1 {
		t.Errorf("(3) approval_decision ต้องเป็น approved แถวเดียว แต่ได้ %q", res)
	}
}

// ---------- AC-APV-15 ----------

// tcDBClockSlack เผื่อเวลาของฐานข้อมูลใน container ต่างจากเครื่องที่รัน test เล็กน้อย (decided_at มาจาก now() ของฐานข้อมูล)
const tcDBClockSlack = 5 * time.Second

// TC-APV-15-1 (AC-APV-15 ปกติ)
func TestTC_APV_15_1_ApproveHistory(t *testing.T) {
	// Given staff-001 อนุมัติคำขอของ AP-TEST-0002 ระหว่างเวลา T1 ถึง T2
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	srv := tcServer(tx, tcBeforeCutoff(t, tx, id))
	t1 := time.Now()
	if w := tcApprove(srv, tokenStaffDept1, id); w.Code != http.StatusOK {
		t.Fatalf("อนุมัติต้องได้ 200 แต่ได้ %d", w.Code)
	}
	t2 := time.Now()

	// When อ่านประวัติของคำขอนั้น
	res, reason, actor, decidedAt, ok := tcDecision(t, tx, id)

	// Then (1) แถวใหม่ 1 แถว ผล approved ไม่มีเหตุผล ผู้ดำเนินการ staff-001 decided_at ระหว่าง T1 ถึง T2
	if !ok || tcDecisionCount(t, tx, id) != 1 || res != "approved" || reason != "" || actor != "staff-001" {
		t.Errorf("(1) ได้ %q / %q / %q (มีแถว %v)", res, reason, actor, ok)
	}
	if decidedAt.Before(t1.Add(-tcDBClockSlack)) || decidedAt.After(t2.Add(tcDBClockSlack)) {
		t.Errorf("(1) decided_at %v ต้องอยู่ระหว่าง %v ถึง %v", decidedAt, t1, t2)
	}
	// Then (2) ชื่อเจ้าหน้าที่: ยังไม่ตรวจ (รอ Q-APV-04)
	// Then (3) มีแถว audit_log action approve
	if n := tcAuditCount(t, tx, id, "approve"); n != 1 {
		t.Errorf("(3) ต้องมี audit_log approve 1 แถว แต่มี %d", n)
	}
}

// TC-APV-15-2 (AC-APV-15 ปกติ)
func TestTC_APV_15_2_RejectHistory(t *testing.T) {
	// Given staff-001 ปฏิเสธคำขอของ AP-TEST-0003 ด้วยเหตุผล "อุปกรณ์ต้องใช้สอบ"
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	if w := tcReject(tcServer(tx, nil), tokenStaffDept1, id, "อุปกรณ์ต้องใช้สอบ", time.Now()); w.Code != http.StatusOK {
		t.Fatalf("ปฏิเสธต้องได้ 200 แต่ได้ %d", w.Code)
	}

	// When อ่านประวัติของคำขอนั้น
	res, reason, actor, decidedAt, ok := tcDecision(t, tx, id)

	// Then (1) แถวใหม่ผล rejected เหตุผลตรงตัว ผู้ดำเนินการ staff-001 มีวันเวลา
	if !ok || res != "rejected" || reason != "อุปกรณ์ต้องใช้สอบ" || actor != "staff-001" || decidedAt.IsZero() {
		t.Errorf("(1) ได้ %q / %q / %q / %v", res, reason, actor, decidedAt)
	}
	// Then (2) มีแถว audit_log action reject
	if n := tcAuditCount(t, tx, id, "reject"); n != 1 {
		t.Errorf("(2) ต้องมี audit_log reject 1 แถว แต่มี %d", n)
	}
}

// TC-APV-15-3 (AC-APV-15 ปกติ)
func TestTC_APV_15_3_AutoRejectHistory(t *testing.T) {
	// Given คำขอที่ส่งเมื่อวาน ถูกปฏิเสธอัตโนมัติตอนเที่ยงคืน
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight.Add(-3*time.Hour))
	tcRunMidnight(t, tx, midnight)

	// When อ่านประวัติของคำขอนั้น
	var result, reason, actorID, actorName string
	var decidedAt time.Time
	if err := tx.QueryRow(context.Background(),
		`SELECT result, reason, actor_id, actor_name, decided_at FROM approval_decision WHERE request_id = $1`, id).
		Scan(&result, &reason, &actorID, &actorName, &decidedAt); err != nil {
		t.Fatal(err)
	}

	// Then (1) ผล auto_rejected เหตุผลตามที่กำหนด ผู้ดำเนินการ system ชื่อ System มีวันเวลา
	if result != "auto_rejected" || reason != tcAutoReason || actorID != "system" || actorName != "System" || decidedAt.IsZero() {
		t.Errorf("(1) ได้ %q / %q / %q / %q / %v", result, reason, actorID, actorName, decidedAt)
	}
	// Then (2) มีแถว audit_log action auto_reject
	if n := tcAuditCount(t, tx, id, "auto_reject"); n != 1 {
		t.Errorf("(2) ต้องมี audit_log auto_reject 1 แถว แต่มี %d", n)
	}
}

// TC-APV-15-4 (AC-APV-15 ขอบ)
func TestTC_APV_15_4_SupersededHistory(t *testing.T) {
	// Given AP-TEST-0001 มี A, B และ C staff-001 อนุมัติ A
	tx := devTx(t)
	a := requestIDOf(t, tx, "AP-TEST-0001", "student-001")
	b := requestIDOf(t, tx, "AP-TEST-0001", "student-002")
	c := requestIDOf(t, tx, "AP-TEST-0001", "student-003")
	if w := tcApprove(tcServer(tx, tcBeforeCutoff(t, tx, a)), tokenStaffDept1, a); w.Code != http.StatusOK {
		t.Fatalf("อนุมัติ A ต้องได้ 200 แต่ได้ %d", w.Code)
	}

	// When อ่านประวัติของ A, B, C
	// Then (1) มีประวัติ 1 แถวต่อคำขอ รวม 3 แถว
	for name, id := range map[string]int{"A": a, "B": b, "C": c} {
		if n := tcDecisionCount(t, tx, id); n != 1 {
			t.Errorf("(1) %s ต้องมีประวัติ 1 แถว แต่มี %d", name, n)
		}
	}
	// Then (2) B และ C ผล rejected เหตุผลตามที่กำหนด ผู้ดำเนินการ staff-001
	for name, id := range map[string]int{"B": b, "C": c} {
		if res, reason, actor, _, _ := tcDecision(t, tx, id); res != "rejected" || reason != tcMsgSuperseded || actor != "staff-001" {
			t.Errorf("(2) %s ได้ %q / %q / %q", name, res, reason, actor)
		}
	}
}

// TC-APV-15-5 (AC-APV-15 ขอบ)
func TestTC_APV_15_5_HistoryAppendOnly(t *testing.T) {
	// Given มีประวัติการอนุมัติของ AP-TEST-0002 แล้ว
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	if w := tcApprove(tcServer(tx, tcBeforeCutoff(t, tx, id)), tokenStaffDept1, id); w.Code != http.StatusOK {
		t.Fatalf("อนุมัติต้องได้ 200 แต่ได้ %d", w.Code)
	}
	before, _, _, beforeAt, _ := tcDecision(t, tx, id)

	// When สั่ง UPDATE และ DELETE แถวนั้นตรง ๆ (ใน savepoint เพื่อให้ transaction ของ test ใช้ต่อได้)
	try := func(sql string) error {
		ctx := context.Background()
		sp, err := tx.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		_, execErr := sp.Exec(ctx, sql, id)
		_ = sp.Rollback(ctx)
		return execErr
	}
	updErr := try(`UPDATE approval_decision SET result = 'rejected', reason = 'แก้' WHERE request_id = $1`)
	delErr := try(`DELETE FROM approval_decision WHERE request_id = $1`)

	// Then (1) ฐานข้อมูลไม่ยอมทั้ง 2 คำสั่ง
	if updErr == nil || delErr == nil {
		t.Errorf("(1) ต้องแก้และลบไม่ได้ แต่ได้ update=%v delete=%v", updErr, delErr)
	}
	// Then (2) ประวัติยังเหมือนเดิม
	after, _, _, afterAt, ok := tcDecision(t, tx, id)
	if !ok || after != before || !afterAt.Equal(beforeAt) {
		t.Errorf("(2) ประวัติต้องเหมือนเดิม %q %v แต่ได้ %q %v (มีแถว %v)", before, beforeAt, after, afterAt, ok)
	}
}

// TC-APV-15-6 (AC-APV-15 ผิด)
func TestTC_APV_15_6_NoHistoryWhenNotSaved(t *testing.T) {
	// Given staff-001 อนุมัติคำขอของ AP-TEST-0004 (ถูกยืม ได้ 409)
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0004", "student-001")
	if w := tcApprove(tcServer(tx, tcBeforeCutoff(t, tx, id)), tokenStaffDept1, id); w.Code != http.StatusConflict {
		t.Fatalf("เตรียมข้อมูล: ต้องได้ 409 แต่ได้ %d", w.Code)
	}

	// When อ่านประวัติของคำขอนั้น
	// Then (1) ไม่มีแถว approval_decision
	if n := tcDecisionCount(t, tx, id); n != 0 {
		t.Errorf("(1) ต้องไม่มี approval_decision แต่มี %d", n)
	}
	// Then (2) ไม่มีแถว audit_log action approve ของคำขอนั้น
	if n := tcAuditCount(t, tx, id, "approve"); n != 0 {
		t.Errorf("(2) ต้องไม่มี audit_log approve แต่มี %d", n)
	}
}

// ---------- AC-APV-18 ----------

func tcCatchUp(tx pgx.Tx, at time.Time) {
	job := &AutoRejecter{db: tx, now: func() time.Time { return at }}
	job.CatchUp(context.Background())
}

// TC-APV-18-1 (AC-APV-18 ปกติ)
func TestTC_APV_18_1_CatchUpOnStart(t *testing.T) {
	// Given คำขอยืมปกติของ AP-TEST-0003 ส่งวันที่ D ระบบหยุดก่อนเที่ยงคืน เริ่มใหม่ 09:30 น. ของ D+1
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight.Add(-3*time.Hour))

	// When CatchUp
	tcCatchUp(tx, midnight.Add(9*time.Hour+30*time.Minute))

	// Then (1) rejected ผล auto_rejected ผู้ดำเนินการ System
	var result, actorName string
	if err := tx.QueryRow(context.Background(),
		`SELECT result, actor_name FROM approval_decision WHERE request_id = $1`, id).Scan(&result, &actorName); err != nil {
		t.Fatalf("(1) ต้องมีผลการพิจารณา: %v", err)
	}
	if tcRequestStatus(t, tx, id) != "rejected" || result != "auto_rejected" || actorName != "System" {
		t.Errorf("(1) ได้ %q / %q / %q", tcRequestStatus(t, tx, id), result, actorName)
	}
	// Then (2) มีแถว approval_email แบบ auto_rejected
	if k := tcEmailKind(t, tx, id); k != "auto_rejected" {
		t.Errorf("(2) Email ต้องเป็นแบบ auto_rejected แต่ได้ %q", k)
	}
}

// TC-APV-18-2 (AC-APV-18 ปกติ)
func TestTC_APV_18_2_CatchUpEmail(t *testing.T) {
	// Given ต่อจาก TC-APV-18-1
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight.Add(-3*time.Hour))
	started := midnight.Add(9*time.Hour + 30*time.Minute)
	tcCatchUp(tx, started)

	// When ตัวส่ง Email ทำงานที่ 09:30:05 น.
	env := tcNewEmailEnv(t, tx, started.Add(5*time.Second))
	env.sendDue(t)

	// Then (1) ได้ Email ปฏิเสธอัตโนมัติภายใน 60 วินาทีหลังเริ่มทำงาน
	var mine *Email
	for i := range env.mailer.sent {
		if strings.Contains(env.mailer.sent[i].Body, "AP-TEST-0003") {
			mine = &env.mailer.sent[i]
		}
	}
	if mine == nil || mine.Subject != "แจ้งปฏิเสธคำขอยืมอุปกรณ์อัตโนมัติ" {
		t.Fatalf("(1) ต้องได้ Email ปฏิเสธอัตโนมัติของ AP-TEST-0003 แต่ได้ %+v", env.mailer.sent)
	}
	// Then (2) เวลาที่ดำเนินการยังเป็น "00:00 น."
	if !strings.Contains(mine.Body, "00:00 น.") || strings.Contains(mine.Body, "09:30") {
		t.Errorf("(2) เวลาที่ดำเนินการต้องเป็น 00:00 น.\n%s", mine.Body)
	}
}

// TC-APV-18-3 (AC-APV-18 ขอบ)
func TestTC_APV_18_3_CatchUpManyDays(t *testing.T) {
	// Given คำขอยืมปกติส่งวันที่ D, D+1 และ D+2 อย่างละ 1 คำขอ เริ่มทำงานใหม่วันที่ D+3
	tx := devTx(t)
	midnight := tcMidnight() // 00:00 น. หลังวันที่ D
	ids := []int{
		requestIDOf(t, tx, "AP-TEST-0002", "student-001"),
		requestIDOf(t, tx, "AP-TEST-0003", "student-002"),
		requestIDOf(t, tx, "AP-TEST-0004", "student-001"),
	}
	for day, id := range ids {
		setSubmittedAt(t, tx, id, midnight.AddDate(0, 0, day).Add(-3*time.Hour))
	}

	// When CatchUp วันที่ D+3
	tcCatchUp(tx, midnight.AddDate(0, 0, 3).Add(9*time.Hour))

	// Then (1) ทั้ง 3 คำขอเป็น rejected ผล auto_rejected ในรอบเดียว
	for _, id := range ids {
		if res, _, _, _, _ := tcDecision(t, tx, id); res != "auto_rejected" || tcRequestStatus(t, tx, id) != "rejected" {
			t.Errorf("(1) คำขอ %d ต้องถูกปฏิเสธอัตโนมัติ แต่ได้ผล %q", id, res)
		}
	}
}

// TC-APV-18-4 (AC-APV-18 ขอบ)
func TestTC_APV_18_4_CatchUpTwice(t *testing.T) {
	// Given เหมือน TC-APV-18-1 แล้วเรียก CatchUp ซ้ำอีกครั้ง
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight.Add(-3*time.Hour))
	started := midnight.Add(9*time.Hour + 30*time.Minute)
	tcCatchUp(tx, started)

	// When CatchUp ครั้งที่ 2
	tcCatchUp(tx, started.Add(time.Minute))

	// Then (1) ประวัติแถวเดียว และ approval_email แถวเดียว
	if n := tcDecisionCount(t, tx, id); n != 1 {
		t.Errorf("(1) ต้องมีประวัติแถวเดียว แต่มี %d", n)
	}
	if n := tcEmailCount(t, tx, id); n != 1 {
		t.Errorf("(1) ต้องมี approval_email แถวเดียว แต่มี %d", n)
	}
}

// TC-APV-18-5 (AC-APV-18 ผิด)
func TestTC_APV_18_5_StartBeforeMidnight(t *testing.T) {
	// Given เริ่มทำงานก่อนเที่ยงคืนของวันที่ D (23:00 น.) คำขอส่งวันที่ D
	tx := devTx(t)
	id := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	midnight := tcMidnight()
	setSubmittedAt(t, tx, id, midnight.Add(-3*time.Hour))

	// When CatchUp ที่ 23:00 น. ของวันที่ D
	tcCatchUp(tx, midnight.Add(-time.Hour))

	// Then (1) คำขอยังเป็น requested
	if s := tcRequestStatus(t, tx, id); s != "requested" {
		t.Errorf("(1) คำขอต้องยังเป็น requested แต่ได้ %q", s)
	}
}
