package approve

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"

	"salmonnorway/backend/internal/config"
)

// AutoRejecter คืองานตั้งเวลาที่ปฏิเสธคำขอยืมแบบปกติที่ยังเป็น requested เมื่อถึงเที่ยงคืนของวันที่ส่ง
// (FR-APV-01 ฉ, FR-APV-02 ค, ASM-05, Q-APV-14, Q-APV-15)
// เป็น goroutine ใน process เดียวกับหลังบ้าน (plan.md ข้อ 2 ทีมเลือกเอง ไม่ได้มาจาก spec)
// การไล่ปฏิเสธตอน server เริ่มทำงานและการเริ่ม goroutine ใน main.go เป็นงานของ T-30
type AutoRejecter struct {
	db  Querier
	now func() time.Time // นาฬิกาที่ inject ได้ ใช้ทดสอบ 00:00 น. (plan.md ข้อ 6)
}

// NewAutoRejecter สร้างงานที่ใช้นาฬิกาจริง
func NewAutoRejecter(db Querier) *AutoRejecter {
	return &AutoRejecter{db: db, now: time.Now}
}

// autoRejectCutoff คือเวลา config.ApprovalAutoRejectHour (00:00 น.) ของวันนี้ตามเวลาไทย (glossary)
// คำขอที่ส่งก่อนเวลานี้ = ส่งก่อนวันนี้ = เลยเที่ยงคืนของวันที่ส่งแล้ว (ASM-05 นับจากวันที่ส่ง Q-APV-15)
func autoRejectCutoff(now time.Time) time.Time {
	lt := localTime(now)
	return time.Date(lt.Year(), lt.Month(), lt.Day(), config.ApprovalAutoRejectHour, 0, 0, 0, lt.Location())
}

// Run ตื่นทุก 00:00 น. เวลาไทยแล้วปฏิเสธคำขอที่ค้าง จนกว่า ctx ถูกยกเลิก
func (a *AutoRejecter) Run(ctx context.Context) {
	for {
		next := autoRejectCutoff(a.now()).AddDate(0, 0, 1)
		t := time.NewTimer(next.Sub(a.now()))
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
		}
		if n, err := a.RejectStale(ctx); err != nil {
			log.Printf("approve auto-reject: %v", err)
		} else if n > 0 {
			log.Printf("approve auto-reject: ปฏิเสธอัตโนมัติ %d คำขอ", n)
		}
	}
}

// staleRequest คือคำขอที่เลยเที่ยงคืนของวันที่ส่งแล้วและยังไม่มีใครพิจารณา
type staleRequest struct {
	ID           int
	AssetCode    string
	DepartmentID int // ภาควิชาเจ้าของอุปกรณ์ บันทึกเป็นภาควิชาของผู้ดำเนินการ System (plan.md ข้อ 3)
}

// RejectStale ปฏิเสธคำขอยืมแบบปกติทุกใบที่ยังเป็น requested และส่งก่อน 00:00 น. ของวันนี้ ใน transaction เดียว คืนจำนวนที่ปฏิเสธ
//   - คำขอค้างคืนไม่ถูกแตะ (Q-APV-14)
//   - result auto_rejected ผู้ดำเนินการ System เหตุผล config.ApprovalAutoRejectReason (FR-APV-01 ฉ)
//   - audit_log action auto_reject (RQ-25) และ Email ปฏิเสธอัตโนมัติลงคิว (FR-APV-02 ค)
//
// ล็อกแถวด้วย FOR UPDATE (รอ ไม่ข้าม) คำขอที่เจ้าหน้าที่กำลังพิจารณาอยู่จะถูกตรวจสถานะใหม่หลังรอ
// คำขอที่ถูกพิจารณาไปแล้วระหว่างนั้นจึงไม่ถูกปฏิเสธซ้ำ
func (a *AutoRejecter) RejectStale(ctx context.Context) (int, error) {
	now := a.now()
	cutoff := autoRejectCutoff(now)
	repo := NewRepository(a.db)
	count := 0
	err := repo.inTx(ctx, func(tx *Repository) error {
		stale, err := tx.lockStaleRequests(ctx, cutoff)
		if err != nil {
			return err
		}
		for _, r := range stale {
			if err := tx.setStatus(ctx, r.ID, config.StatusRejected); err != nil {
				return err
			}
			if _, err := tx.insertDecision(ctx, decision{
				RequestID:         r.ID,
				Result:            config.ApprovalKindAutoRejected,
				Reason:            config.ApprovalAutoRejectReason,
				ActorID:           config.ApprovalSystemActorID,
				ActorName:         config.ApprovalSystemActorName, // FR-APV-01 (ฉ) "ผู้ดำเนินการคือ System"
				ActorDepartmentID: r.DepartmentID,
			}); err != nil {
				return err
			}
			if err := tx.insertAudit(ctx, config.ApprovalSystemActorID, config.ApprovalSystemActorID,
				config.AuditActionAutoReject, r.AssetCode, r.ID); err != nil {
				return err
			}
			if err := tx.insertEmail(ctx, r.ID, config.ApprovalKindAutoRejected, now); err != nil {
				return err
			}
			count++
		}
		return nil
	})
	return count, err
}

// lockStaleRequests ล็อกและคืนคำขอยืมแบบปกติที่ยังเป็น requested และส่งก่อน cutoff (FR-APV-01 ฉ)
func (r *Repository) lockStaleRequests(ctx context.Context, cutoff time.Time) ([]staleRequest, error) {
	const q = `
SELECT r.id, i.asset_code, m.department_id
FROM borrow_request r
JOIN equipment_item  i ON i.id = r.item_id
JOIN equipment_model m ON m.id = i.model_id
WHERE r.status = $1 AND r.loan_type = $2 AND r.submitted_at < $3
ORDER BY r.submitted_at, r.id
FOR UPDATE OF r`
	rows, err := r.db.Query(ctx, q, config.StatusRequested, config.LoanTypeNormal, cutoff)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (staleRequest, error) {
		var s staleRequest
		err := row.Scan(&s.ID, &s.AssetCode, &s.DepartmentID)
		return s, err
	})
}
