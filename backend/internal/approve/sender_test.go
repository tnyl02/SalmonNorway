package approve

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"salmonnorway/backend/internal/config"
)

// fakeMailer จดทุก Email ที่ได้รับ และสั่งให้ส่งไม่สำเร็จได้ (plan.md ข้อ 6: ตัวส่ง Email จำลองบันทึกเวลาและสั่งให้ล้มได้)
type fakeMailer struct {
	mu    sync.Mutex
	sent  []Email
	fail  error
	calls []time.Time // เวลาของนาฬิกาจำลองทุกครั้งที่ถูกเรียก (ทั้งสำเร็จและล้ม) ใช้กับ clock
	clock *fakeClock
}

func (m *fakeMailer) Send(_ context.Context, e Email) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.clock != nil {
		m.calls = append(m.calls, m.clock.Now())
	}
	if m.fail != nil {
		return m.fail
	}
	m.sent = append(m.sent, e)
	return nil
}

// fakeClock คือนาฬิกาจำลองที่เลื่อนเวลาได้
type fakeClock struct{ t time.Time }

func (c *fakeClock) Now() time.Time          { return c.t }
func (c *fakeClock) Advance(d time.Duration) { c.t = c.t.Add(d) }

// AC-APV-04 (FR-APV-02, IF-NOT-01)
// Given เจ้าหน้าที่เพิ่งอนุมัติหรือปฏิเสธคำขอ
// When  ระบบบันทึกผลสำเร็จ
// Then  ระบบส่ง Email ถึงผู้ขอภายใน 60 วินาที หัวข้อและข้อมูลครบตาม FR-APV-02 (ก) หรือ (ข)
func TestAC_APV_04(t *testing.T) {
	tx := devTx(t)
	clock := &fakeClock{t: time.Now()}
	svc := NewService(NewRepository(tx))
	svc.now = clock.Now
	mailer := &fakeMailer{}
	sender := &Sender{db: tx, mailer: mailer, resolver: MockRecipientResolver{}, now: clock.Now}

	ctx := staffCtx("staff-001", 1)
	approveID := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	rejectID := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	if _, err := svc.Approve(ctx, approveID, time.Time{}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Reject(ctx, rejectID, "อุปกรณ์ต้องใช้สอบ", time.Time{}); err != nil {
		t.Fatal(err)
	}

	// บันทึกผลแล้ว Email อยู่ในคิวสถานะ pending (outbox ใน transaction เดียวกับผล)
	var pending int
	if err := tx.QueryRow(context.Background(),
		`SELECT count(*) FROM approval_email WHERE request_id = ANY($1) AND status = 'pending'`, []int{approveID, rejectID}).Scan(&pending); err != nil {
		t.Fatal(err)
	}
	if pending != 2 {
		t.Fatalf("ต้องมี Email รอส่ง 2 ฉบับ แต่ได้ %d", pending)
	}

	// When: ตัวส่งตรวจคิวรอบถัดไป (ช้าสุดหลังบันทึกผล 1 ช่วง poll)
	clock.Advance(config.ApprovalSenderPollInterval)
	if err := sender.SendDue(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Then: ส่งครบ 2 ฉบับถึงผู้ขอ หัวข้อตรงแบบ
	if len(mailer.sent) != 2 {
		t.Fatalf("ต้องส่ง Email 2 ฉบับ แต่ได้ %d", len(mailer.sent))
	}
	byTo := map[string]Email{}
	for _, e := range mailer.sent {
		byTo[e.To] = e
	}
	approved, rejected := byTo["student-001@example.com"], byTo["student-002@example.com"]
	if approved.Subject != config.ApprovalEmailSubjects[config.ApprovalKindApproved] ||
		!strings.Contains(approved.Body, "สถานที่รับ: "+config.ApprovalPickupLocation) ||
		!strings.Contains(approved.Body, "AP-TEST-0002") {
		t.Errorf("Email ผลอนุมัติไม่ถูกต้อง: %+v", approved)
	}
	if rejected.Subject != config.ApprovalEmailSubjects[config.ApprovalKindRejected] ||
		!strings.Contains(rejected.Body, "เหตุผลในการปฏิเสธ: อุปกรณ์ต้องใช้สอบ") {
		t.Errorf("Email ผลปฏิเสธไม่ถูกต้อง: %+v", rejected)
	}

	// ภายใน 60 วินาทีหลังบันทึกผล และสถานะการส่งเป็น sent
	rows, err := tx.Query(context.Background(),
		`SELECT status, extract(epoch FROM last_attempt_at - created_at)::float8 FROM approval_email WHERE request_id = ANY($1)`, []int{approveID, rejectID})
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var seconds float64
		if err := rows.Scan(&status, &seconds); err != nil {
			t.Fatal(err)
		}
		took := time.Duration(seconds * float64(time.Second))
		if status != config.ApprovalEmailSent || took > config.ApprovalEmailDeadline {
			t.Errorf("Email ต้องเป็น sent ภายใน %v แต่ได้ %q ใช้เวลา %v", config.ApprovalEmailDeadline, status, took)
		}
	}
}

// T-15: ตัวส่งตรวจคิวถี่พอให้ส่งทันภายใน 60 วินาที (FR-APV-02)
func TestSender_PollWithinDeadline(t *testing.T) {
	if config.ApprovalSenderPollInterval <= 0 || config.ApprovalSenderPollInterval >= config.ApprovalEmailDeadline {
		t.Errorf("ช่วงตรวจคิว %v ต้องมากกว่า 0 และน้อยกว่า %v", config.ApprovalSenderPollInterval, config.ApprovalEmailDeadline)
	}
}

// T-32 (FR-APV-02 ง, Q-APV-17) ส่วนส่งซ้ำของ AC-APV-14
// Given บันทึกผลแล้ว แต่ระบบแจ้งเตือน Email ส่งไม่สำเร็จทุกครั้ง
// When  ส่งครั้งแรกแล้วส่งซ้ำอัตโนมัติอีกครบ 3 ครั้ง ห่างกันครั้งละ 1 นาที
// Then  ตัวส่งถูกเรียกรวม 4 ครั้ง สถานะการส่งเป็น failed พร้อมสาเหตุและเวลา และผลการพิจารณาไม่เปลี่ยน
func TestSender_RetryThenFailed(t *testing.T) {
	tx := devTx(t)
	clock := &fakeClock{t: time.Now()}
	svc := NewService(NewRepository(tx))
	svc.now = clock.Now
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")
	if _, err := svc.Approve(staffCtx("staff-001", 1), id, time.Time{}); err != nil {
		t.Fatal(err)
	}

	mailer := &fakeMailer{fail: errors.New("smtp: connection refused"), clock: clock}
	sender := &Sender{db: tx, mailer: mailer, resolver: MockRecipientResolver{}, now: clock.Now}
	run := func() {
		t.Helper()
		if err := sender.SendDue(context.Background()); err != nil {
			t.Fatal(err)
		}
	}

	run() // ครั้งแรก
	// ยังไม่ถึง 1 นาที ต้องไม่ลองซ้ำ
	clock.Advance(config.ApprovalEmailRetryInterval - time.Second)
	run()
	if len(mailer.calls) != 1 {
		t.Fatalf("ก่อนครบ 1 นาทีต้องไม่ส่งซ้ำ แต่ถูกเรียก %d ครั้ง", len(mailer.calls))
	}
	// ส่งซ้ำ 3 ครั้ง ห่างกันครั้งละ 1 นาที
	clock.Advance(time.Second)
	for i := 0; i < config.ApprovalEmailRetryCount; i++ {
		run()
		clock.Advance(config.ApprovalEmailRetryInterval)
	}
	run() // หลังเป็น failed แล้วต้องไม่ส่งอีก

	if want := 1 + config.ApprovalEmailRetryCount; len(mailer.calls) != want {
		t.Fatalf("ตัวส่งต้องถูกเรียกรวม %d ครั้ง แต่ได้ %d", want, len(mailer.calls))
	}
	for i := 1; i < len(mailer.calls); i++ {
		if gap := mailer.calls[i].Sub(mailer.calls[i-1]); gap < config.ApprovalEmailRetryInterval {
			t.Errorf("ครั้งที่ %d ห่างจากครั้งก่อน %v ต้องไม่น้อยกว่า %v", i+1, gap, config.ApprovalEmailRetryInterval)
		}
	}

	var status, lastError, decision, reqStatus string
	var attempts int
	var lastAttempt *time.Time
	if err := tx.QueryRow(context.Background(), `
		SELECT e.status, e.attempts, e.last_error, e.last_attempt_at, d.result, r.status
		FROM approval_email e JOIN approval_decision d ON d.request_id = e.request_id
		JOIN borrow_request r ON r.id = e.request_id WHERE e.request_id = $1`, id).
		Scan(&status, &attempts, &lastError, &lastAttempt, &decision, &reqStatus); err != nil {
		t.Fatal(err)
	}
	if status != config.ApprovalEmailFailed || attempts != 4 || !strings.Contains(lastError, "connection refused") || lastAttempt == nil {
		t.Errorf("Email ต้องเป็น failed 4 ครั้ง พร้อมสาเหตุและเวลา แต่ได้ %q %d %q %v", status, attempts, lastError, lastAttempt)
	}
	if decision != "approved" || reqStatus != "approved" {
		t.Errorf("ผลการพิจารณาต้องไม่เปลี่ยน แต่ได้ %q / %q", decision, reqStatus)
	}
}
