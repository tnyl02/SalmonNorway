package approve

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/config"
)

// setSubmittedAt เปลี่ยนเวลาที่ส่งของคำขอ ใช้จำลองคำขอที่ส่งเมื่อวาน / หลังเที่ยงคืน
func setSubmittedAt(t *testing.T, tx pgx.Tx, id int, at time.Time) {
	t.Helper()
	if _, err := tx.Exec(context.Background(), `UPDATE borrow_request SET submitted_at = $2 WHERE id = $1`, id, at); err != nil {
		t.Fatal(err)
	}
}

// AC-APV-13 (FR-APV-01 ฉ, FR-APV-02 ค)
// Given คำขอที่ส่งวันที่ D และยังเป็น requested
// When  ถึง 00:00 น. หลังวันที่ D
// Then  คำขอเป็น rejected ผู้ดำเนินการคือ System และผู้ขอได้รับ Email ปฏิเสธอัตโนมัติตาม FR-APV-02 (ค) ภายใน 60 วินาที
//
//	(T-29 เสร็จเมื่อ: คำขอที่ส่งวัน D+1 และคำขอค้างคืนที่ส่งวัน D ไม่ถูกแตะ)
func TestAC_APV_13(t *testing.T) {
	tx := devTx(t)
	// นาฬิกาจำลอง: 00:00 น. เวลาไทยของวันถัดไป (วัน D+1) บวก 1 นาที
	midnight := autoRejectCutoff(time.Now()).AddDate(0, 0, 1)
	clock := &fakeClock{t: midnight.Add(time.Minute)}

	stale := requestIDOf(t, tx, "AP-TEST-0003", "student-002")     // ยืมปกติ ส่งวัน D
	fresh := requestIDOf(t, tx, "AP-TEST-0002", "student-001")     // ยืมปกติ ส่งวัน D+1 หลังเที่ยงคืน
	overnight := requestIDOf(t, tx, "AP-TEST-0001", "student-003") // ค้างคืน ส่งวัน D (Q-APV-14)
	setSubmittedAt(t, tx, stale, midnight.Add(-14*time.Hour))
	setSubmittedAt(t, tx, overnight, midnight.Add(-14*time.Hour))
	setSubmittedAt(t, tx, fresh, midnight.Add(30*time.Second))

	// When: งานเที่ยงคืนทำงาน
	job := &AutoRejecter{db: tx, now: clock.Now}
	if _, err := job.RejectStale(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Then: คำขอที่ส่งวัน D เป็น rejected result auto_rejected ผู้ดำเนินการ System
	statusOf := func(id int) string {
		t.Helper()
		var s string
		if err := tx.QueryRow(context.Background(), `SELECT status FROM borrow_request WHERE id = $1`, id).Scan(&s); err != nil {
			t.Fatal(err)
		}
		return s
	}
	if s := statusOf(stale); s != "rejected" {
		t.Errorf("คำขอที่ส่งวัน D ต้องเป็น rejected แต่ได้ %q", s)
	}
	var actorName string
	if err := tx.QueryRow(context.Background(),
		`SELECT actor_name FROM approval_decision WHERE request_id = $1`, stale).Scan(&actorName); err != nil {
		t.Fatal(err)
	}
	if actorName != "System" {
		t.Errorf("ชื่อผู้ดำเนินการต้องเป็น System แต่ได้ %q", actorName)
	}
	// คำขอที่ส่งวัน D+1 และคำขอค้างคืนไม่ถูกแตะ
	if s := statusOf(fresh); s != "requested" {
		t.Errorf("คำขอที่ส่งวัน D+1 ต้องยังเป็น requested แต่ได้ %q", s)
	}
	if s := statusOf(overnight); s != "requested" {
		t.Errorf("คำขอค้างคืนต้องยังเป็น requested แต่ได้ %q (Q-APV-14)", s)
	}

	// Email ปฏิเสธอัตโนมัติภายใน 60 วินาที (ตัวส่งตรวจคิวรอบถัดไป)
	mailer := &fakeMailer{}
	clock.Advance(config.ApprovalSenderPollInterval)
	sender := &Sender{db: tx, mailer: mailer, resolver: MockRecipientResolver{}, now: clock.Now}
	if err := sender.SendDue(context.Background()); err != nil {
		t.Fatal(err)
	}
	// คำขอยืมปกติอื่นใน approve_dev.sql ก็ส่งก่อนเที่ยงคืนนี้ จึงถูกปฏิเสธพร้อมกัน หา Email ของ stale (ผู้ขอ student-002)
	var e Email
	for _, m := range mailer.sent {
		if m.To == "student-002@example.com" && strings.Contains(m.Body, "AP-TEST-0003") {
			e = m
		}
	}
	if e.Subject != config.ApprovalEmailSubjects[config.ApprovalKindAutoRejected] ||
		!strings.Contains(e.Body, "สถานะ: ปฏิเสธอัตโนมัติ") || !strings.Contains(e.Body, "เวลาที่ดำเนินการ: 00:00 น.") {
		t.Errorf("Email ปฏิเสธอัตโนมัติไม่ถูกต้อง: %+v", e)
	}
	if took := clock.Now().Sub(midnight.Add(time.Minute)); took > config.ApprovalEmailDeadline {
		t.Errorf("Email ต้องออกภายใน %v แต่ใช้ %v", config.ApprovalEmailDeadline, took)
	}

	// รันซ้ำไม่ปฏิเสธซ้ำ
	if n, err := job.RejectStale(context.Background()); err != nil || n != 0 {
		t.Errorf("รันซ้ำต้องไม่มีคำขอให้ปฏิเสธ แต่ได้ %d (%v)", n, err)
	}
}

// AC-APV-18 (FR-APV-01 ฉ, Q-APV-21)
// Given คำขอยืมปกติที่ส่งวันที่ D ยังเป็น requested และระบบไม่ได้ทำงานตอน 00:00 น. หลังวันที่ D
// When  ระบบเริ่มทำงานใหม่หลังเวลานั้น
// Then  คำขอเป็น rejected ผู้ดำเนินการคือ System และผู้ขอได้รับ Email ปฏิเสธอัตโนมัติตาม FR-APV-02 (ค)
func TestAC_APV_18(t *testing.T) {
	tx := devTx(t)
	// server ดับตั้งแต่ก่อนเที่ยงคืน แล้วเริ่มใหม่ตอน 09:30 น. ของวันถัดไป (ไม่มีรอบเที่ยงคืนเกิดขึ้น)
	midnight := autoRejectCutoff(time.Now()).AddDate(0, 0, 1)
	startedAt := midnight.Add(9*time.Hour + 30*time.Minute)
	stale := requestIDOf(t, tx, "AP-TEST-0003", "student-002")
	setSubmittedAt(t, tx, stale, midnight.Add(-3*time.Hour))

	// When: ระบบเริ่มทำงาน (สิ่งที่ Run ทำเป็นอย่างแรก)
	job := &AutoRejecter{db: tx, now: func() time.Time { return startedAt }}
	job.CatchUp(context.Background())

	// Then
	var status, result, actorName string
	if err := tx.QueryRow(context.Background(), `
		SELECT r.status, d.result, d.actor_name FROM borrow_request r JOIN approval_decision d ON d.request_id = r.id
		WHERE r.id = $1`, stale).Scan(&status, &result, &actorName); err != nil {
		t.Fatalf("คำขอต้องถูกปฏิเสธอัตโนมัติตอนเริ่มทำงาน: %v", err)
	}
	if status != "rejected" || result != "auto_rejected" || actorName != "System" {
		t.Errorf("ได้ %q / %q / %q", status, result, actorName)
	}
	var emails int
	if err := tx.QueryRow(context.Background(),
		`SELECT count(*) FROM approval_email WHERE request_id = $1 AND kind = 'auto_rejected'`, stale).Scan(&emails); err != nil {
		t.Fatal(err)
	}
	if emails != 1 {
		t.Errorf("ต้องมี Email ปฏิเสธอัตโนมัติ 1 ฉบับในคิว แต่ได้ %d", emails)
	}
}
