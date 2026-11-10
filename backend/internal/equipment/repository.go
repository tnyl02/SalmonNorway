package equipment

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Querier คือส่วนของ pgxpool.Pool ที่ repository ใช้ (อ่านอย่างเดียว CON-DATA-02)
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Repository อ่านข้อมูลอุปกรณ์จาก view v_equipment_item_status (plan.md ข้อ 4)
// ทุก query ใช้ parameter ($1, $2, ...) ห้ามต่อสตริง SQL (plan.md ข้อ 6.3)
type Repository struct {
	db Querier
}

// NewRepository สร้าง Repository จาก connection pool
func NewRepository(db Querier) *Repository {
	return &Repository{db: db}
}

// ListModels คืนรายการรุ่น 1 แถว = 1 รุ่น พร้อมจำนวนว่างและจำนวนทั้งหมด (FR-EQS-01)
// จำนวนว่างนับจากสถานะ "ว่าง" อย่างเดียว คำขอยืมที่ค้างไม่ลดจำนวน (DOM-REQ-01)
// ยังไม่แบ่งหน้า (รอ Q-05 ดู T-23) และยังไม่จำกัดภาควิชา (T-18)
func (r *Repository) ListModels(ctx context.Context) ([]ModelRowDTO, error) {
	const q = `
SELECT model_id, model_name, category_name,
       count(DISTINCT item_id) FILTER (WHERE status = $1) AS available_count,
       count(DISTINCT item_id)                            AS total_count
FROM v_equipment_item_status
GROUP BY model_id, model_name, category_name
ORDER BY model_name, model_id`
	rows, err := r.db.Query(ctx, q, string(StatusAvailable))
	if err != nil {
		return nil, err
	}
	out, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (ModelRowDTO, error) {
		var m ModelRowDTO
		err := row.Scan(&m.ModelID, &m.Name, &m.Category, &m.AvailableCount, &m.TotalCount)
		return m, err
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ItemRow คืออุปกรณ์รายชิ้น 1 ชิ้นในรุ่นที่เลือก ระบุด้วยรหัสครุภัณฑ์ (FR-EQS-01, CON-DATA-03)
// กำหนดคืน ป้ายเลยกำหนดคืน และข้อมูลตามบทบาทจะเพิ่มใน T-10, T-12, T-17
type ItemRow struct {
	AssetCode string `json:"assetCode"`
	Status    Status `json:"status"`
}

// ListItems คืนอุปกรณ์รายชิ้นทั้งหมดของรุ่นที่เลือก พร้อมสถานะล่าสุดของแต่ละชิ้น เรียงตามรหัสครุภัณฑ์ (FR-EQS-01)
// รุ่นที่ไม่มีอยู่ได้รายการว่าง (plan.md ข้อ 5: ไม่พบ = items ว่าง ไม่ใช่ error)
func (r *Repository) ListItems(ctx context.Context, modelID int) ([]ItemRow, error) {
	const q = `
SELECT DISTINCT asset_code, status
FROM v_equipment_item_status
WHERE model_id = $1
ORDER BY asset_code`
	rows, err := r.db.Query(ctx, q, modelID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (ItemRow, error) {
		var it ItemRow
		err := row.Scan(&it.AssetCode, &it.Status)
		return it, err
	})
}
