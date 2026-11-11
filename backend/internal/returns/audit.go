package returns

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"

	"salmonnorway/backend/internal/auth"
)

// ค่า action ของ 006 ในตาราง audit_log (DM-07 ใน specs/000-shared/data-model.md, migration 0006)
const (
	actionReturnList   = "return_list"   // เปิดหน้ารายการที่ยังไม่คืน
	actionReturnLookup = "return_lookup" // ค้นการยืมจากรหัสครุภัณฑ์
	actionReturnDetail = "return_detail" // เปิดรายละเอียดการยืม
	actionReturn       = "return"        // บันทึกการรับคืน
)

// AuditEntry คือ 1 แถวของ audit_log ใช้เฉพาะคอลัมน์ที่ DM-07 มี (DOM-AUDIT-01, FR-RET-07) ไม่เพิ่มฟิลด์เอง (Q-09)
type AuditEntry struct {
	ActorID   string
	ActorRole string
	Action    string
	Query     any    // เงื่อนไขที่ค้น เก็บเป็น JSONB (nil = NULL)
	AssetCode string // ว่าง = NULL
	Outcome   string // ok / forbidden / error ตาม CHECK ของ audit_log
}

// newAudit สร้างแถว audit ของผู้ใช้ที่ยืนยันตัวตนแล้ว (IF-IDP-01)
func newAudit(p auth.Principal, action string) AuditEntry {
	return AuditEntry{ActorID: p.UserID, ActorRole: string(p.Role), Action: action}
}

// outcomeFor แปลง HTTP status เป็นผลลัพธ์ของ audit_log: 2xx = ok, 403 = forbidden, อื่น ๆ = error
// (DM-07 มีแค่ 3 ค่า การจับคู่นี้ AI เลือกเอง ทีมยืนยัน ดู prompt-log T-23)
func outcomeFor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "ok"
	case status == http.StatusForbidden:
		return "forbidden"
	default:
		return "error"
	}
}

// execer คือสิ่งที่รัน Exec ได้ ใช้ได้ทั้ง pool และ transaction
type execer interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// insertAudit เขียน 1 แถวลง audit_log (append-only ตาม 0002) (DOM-AUDIT-01, RQ-25)
func insertAudit(ctx context.Context, db execer, e AuditEntry) error {
	var query []byte
	if e.Query != nil {
		b, err := json.Marshal(e.Query)
		if err != nil {
			return err
		}
		query = b
	}
	var asset *string
	if e.AssetCode != "" {
		asset = &e.AssetCode
	}
	_, err := db.Exec(ctx, `
		INSERT INTO audit_log (actor_id, actor_role, action, query, asset_code, outcome)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		e.ActorID, e.ActorRole, e.Action, query, asset, e.Outcome)
	return err
}

// WriteAudit เขียน audit นอก transaction ใช้กับการอ่าน และกรณีที่ถูกปฏิเสธหรือผิดพลาด (DOM-AUDIT-01)
func (r *Repository) WriteAudit(ctx context.Context, e AuditEntry) error {
	return insertAudit(ctx, r.pool, e)
}
