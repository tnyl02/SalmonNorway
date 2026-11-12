// Package audit เขียนบันทึกการเข้าถึงลงตาราง audit_log ที่ใช้ร่วมทุกฟีเจอร์ (DOM-AUDIT-01, DM-07, T-19 ของ 001)
// ตารางเป็น append-only บังคับด้วย trigger ใน migration 0002 (T-31) แพ็กเกจนี้จึงมีแค่การเพิ่มแถว
package audit

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

// ค่าของ action ที่ตกลงแล้วใน DM-07 (specs/000-shared/data-model.md) และ CHECK ของ audit_log
// ฟีเจอร์อื่นใช้ค่าของตัวเองได้ (cancel, request, notify) ค่านอก CHECK จะถูกฐานข้อมูลปฏิเสธ
const (
	ActionList   = "list"   // 001 ดูรายการ (รอ Q-11 ดู T-29)
	ActionSearch = "search" // 001 ค้นหา (รอ Q-11 ดู T-29)
	ActionDetail = "detail" // 001 เปิดรายละเอียดอุปกรณ์ (AC-EQS-11)
)

// ค่าของ outcome ตาม CHECK ของ audit_log (migration 0002)
const (
	OutcomeOK        = "ok"
	OutcomeForbidden = "forbidden"
	OutcomeError     = "error"
)

// DB คือสิ่งที่ Write ใช้เขียน ทั้ง *pgxpool.Pool และ pgx.Tx ใช้ได้
// ส่ง transaction มาได้ เพื่อเขียน audit ใน transaction เดียวกับงานหลัก (ตกลงกับ 003 FR-CXL-05)
type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Entry คือบันทึก 1 แถว เวลา (occurred_at) ฐานข้อมูลใส่ให้เอง (DOM-AUDIT-01 ผู้เข้าถึง เวลา และสิ่งที่เข้าถึง)
type Entry struct {
	ActorID   string         // รหัสผู้ใช้จากระบบยืนยันตัวตน (DM-08)
	ActorRole string         // student / staff
	Action    string         // ค่าจาก DM-07
	Query     map[string]any // เงื่อนไขที่ค้น เก็บเฉพาะพารามิเตอร์ค้นหา ไม่เก็บ token (plan 001 ข้อ 6.4) nil = ไม่มี
	AssetCode string         // รหัสครุภัณฑ์ (CON-DATA-03) ว่าง = ไม่มี
	RequestID *int           // รหัสคำขอยืม (ของ 003 ต้องมีคำขอจริงตาม FK ของ 0003) nil = ไม่มี
	Outcome   string         // ok / forbidden / error
}

// Write เพิ่มบันทึก 1 แถวลง audit_log และคืนรหัสแถว (DOM-AUDIT-01)
// ค่าว่างของช่องที่ไม่บังคับบันทึกเป็น NULL
func Write(ctx context.Context, db DB, e Entry) (int64, error) {
	var query []byte
	if e.Query != nil {
		b, err := json.Marshal(e.Query)
		if err != nil {
			return 0, err
		}
		query = b
	}
	var assetCode *string
	if e.AssetCode != "" {
		assetCode = &e.AssetCode
	}
	const q = `
INSERT INTO audit_log (actor_id, actor_role, action, query, asset_code, request_id, outcome)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id`
	var id int64
	err := db.QueryRow(ctx, q, e.ActorID, e.ActorRole, e.Action, query, assetCode, e.RequestID, e.Outcome).Scan(&id)
	return id, err
}
