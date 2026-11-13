package approve

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/config"
)

// Sender คือตัวส่ง Email เบื้องหลัง อ่านคิว approval_email แล้วส่งผ่าน Mailer (FR-APV-02, IF-NOT-01, plan.md ข้อ 1)
// เป็น goroutine ใน process เดียวกับหลังบ้าน ไม่ใช้ cron หรือ queue ภายนอก (plan.md ข้อ 2 ทีมเลือกเอง ไม่ได้มาจาก spec)
// ส่งไม่สำเร็จส่งซ้ำ 3 ครั้งหลังครั้งแรกแล้วเป็น failed (T-32) / main.go เริ่ม Run เป็น goroutine (T-30)
type Sender struct {
	db       Querier
	mailer   Mailer
	resolver RecipientResolver
	now      func() time.Time // นาฬิกาที่ inject ได้ ใช้ทดสอบ 60 วินาที (plan.md ข้อ 6)
}

// NewSender สร้างตัวส่งที่ใช้นาฬิกาจริง
func NewSender(db Querier, mailer Mailer, resolver RecipientResolver) *Sender {
	return &Sender{db: db, mailer: mailer, resolver: resolver, now: time.Now}
}

// Run ตรวจคิวทุก config.ApprovalSenderPollInterval จนกว่า ctx ถูกยกเลิก
// ช่วงตรวจสั้นกว่า config.ApprovalEmailDeadline มาก Email จึงออกภายใน 60 วินาทีหลังบันทึกผล (FR-APV-02)
func (s *Sender) Run(ctx context.Context) {
	t := time.NewTicker(config.ApprovalSenderPollInterval)
	defer t.Stop()
	for {
		if err := s.SendDue(ctx); err != nil {
			log.Printf("approve sender: %v", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// dueEmail คือ Email ในคิวที่ถึงเวลาส่ง พร้อมข้อมูลที่ใช้สร้างเนื้อ Email (FR-APV-02 ก ข ค)
type dueEmail struct {
	ID          int
	Kind        string
	RequesterID string
	Data        EmailData
}

// SendDue ส่ง Email ทุกฉบับที่สถานะ pending และถึงเวลาส่งแล้ว คืน error เฉพาะเมื่อฐานข้อมูลมีปัญหา
// แต่ละรอบทำใน transaction เดียว ล็อกแถวด้วย FOR UPDATE SKIP LOCKED กันส่งซ้ำถ้ามีตัวส่งมากกว่า 1 ตัว
func (s *Sender) SendDue(ctx context.Context) error {
	now := s.now()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	due, err := loadDueEmails(ctx, tx, now)
	if err != nil {
		return err
	}
	for _, e := range due {
		sendErr := s.sendOne(ctx, e)
		if err := markAttempt(ctx, tx, e.ID, now, sendErr); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// sendOne หาผู้รับ สร้างเนื้อ Email แล้วส่ง (FR-APV-02)
func (s *Sender) sendOne(ctx context.Context, e dueEmail) error {
	to, err := s.resolver.Resolve(ctx, e.RequesterID)
	if err != nil {
		return err
	}
	subject, body, err := BuildEmail(e.Kind, e.Data)
	if err != nil {
		return err
	}
	return s.mailer.Send(ctx, Email{To: to, Subject: subject, Body: body})
}

func loadDueEmails(ctx context.Context, tx pgx.Tx, now time.Time) ([]dueEmail, error) {
	const q = `
SELECT e.id, e.kind, r.requester_id,
       r.id, r.requester_name, m.name, i.asset_code, r.purpose, coalesce(r.purpose_other, ''),
       r.start_at, r.end_at, r.submitted_at, coalesce(d.reason, ''), coalesce(d.actor_name, d.actor_id),
       coalesce(d.advisor_name, '')
FROM approval_email e
JOIN approval_decision d ON d.request_id = e.request_id
JOIN borrow_request    r ON r.id = e.request_id
JOIN equipment_item    i ON i.id = r.item_id
JOIN equipment_model   m ON m.id = i.model_id
WHERE e.status = $1 AND e.next_attempt_at <= $2
ORDER BY e.next_attempt_at, e.id
LIMIT $3
FOR UPDATE OF e SKIP LOCKED`
	rows, err := tx.Query(ctx, q, config.ApprovalEmailPending, now, config.ApprovalSenderBatchSize)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (dueEmail, error) {
		var e dueEmail
		d := &e.Data
		err := row.Scan(&e.ID, &e.Kind, &e.RequesterID,
			&d.RequestID, &d.RequesterName, &d.ModelName, &d.AssetCode, &d.Purpose, &d.PurposeOther,
			&d.StartAt, &d.EndAt, &d.SubmittedAt, &d.Reason, &d.ActorName, &d.AdvisorName)
		return e, err
	})
}

// markAttempt บันทึกผลการส่ง 1 ครั้ง (FR-APV-02 ง: สถานะการส่ง สาเหตุ และเวลา) ผลการพิจารณาไม่ถูกแตะเลย
//   - สำเร็จ: เป็น sent
//   - ไม่สำเร็จ: นับครั้ง เก็บสาเหตุและเวลา แล้วลองใหม่หลัง config.ApprovalEmailRetryInterval (1 นาที)
//     ถ้านับรวมครั้งนี้ครบครั้งแรก + config.ApprovalEmailRetryCount (= 4 ครั้ง Q-APV-17) แล้ว เป็น failed หยุดส่งซ้ำ
//     แล้วขึ้นในรายการ "Email ที่ส่งไม่สำเร็จ" (T-33) (AC-APV-14, T-32)
func markAttempt(ctx context.Context, tx pgx.Tx, emailID int, now time.Time, sendErr error) error {
	if sendErr == nil {
		_, err := tx.Exec(ctx, `
UPDATE approval_email SET status = $2, attempts = attempts + 1, last_attempt_at = $3, last_error = NULL
WHERE id = $1`, emailID, config.ApprovalEmailSent, now)
		return err
	}
	// ใน SET ค่า attempts ทางขวาคือค่าก่อนแก้ attempts + 1 จึงเป็นจำนวนครั้งรวมครั้งนี้
	_, err := tx.Exec(ctx, `
UPDATE approval_email
SET attempts = attempts + 1, last_attempt_at = $2, last_error = $3, next_attempt_at = $4,
    status = CASE WHEN attempts + 1 >= $5 THEN $6 ELSE status END
WHERE id = $1`,
		emailID, now, sendErr.Error(), now.Add(config.ApprovalEmailRetryInterval),
		1+config.ApprovalEmailRetryCount, config.ApprovalEmailFailed)
	return err
}
