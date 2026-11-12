package equipment

import (
	"context"
	"errors"
	"time"

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

// Filter คือตัวกรองของหน้ารายการ (FR-EQS-02) ค่าว่าง = ไม่กรองด้วยเงื่อนไขนั้น
// การกรองใช้กับสถานะของอุปกรณ์รายชิ้น (Q-15) ค่าที่ไม่รู้จักจะไม่ตรงกับชิ้นใดเลย ไม่ถือเป็น error
type Filter struct {
	Categories []string // ชื่อประเภท (glossary: "ประเภท" = กลุ่มใหญ่ของรุ่น, Q-13)
	Statuses   []string // สถานะรายชิ้น ใช้ค่าของ CON-DATA-01 (available / borrowed / unavailable)

	// DepartmentID คือขอบเขตภาควิชาที่ service ตั้งจากผู้ใช้ (scopeFor) ไม่ได้มาจาก query string
	// 0 = ไม่จำกัด (เจ้าหน้าที่ รอ Q-P3) นักศึกษา = ภาควิชาของตน (DOM-ACCESS-01, FR-EQS-10)
	DepartmentID int
}

// args คืนตัวกรองเป็น array สำหรับ parameter ของ SQL ไม่ส่ง nil เพื่อให้ cardinality() ได้ 0 แทน NULL
func (f Filter) args() (categories, statuses []string) {
	categories, statuses = f.Categories, f.Statuses
	if categories == nil {
		categories = []string{}
	}
	if statuses == nil {
		statuses = []string{}
	}
	return categories, statuses
}

// เงื่อนไขกรองที่ใช้ร่วมกันทั้งรายการรุ่นและรายชิ้น (plan.md ข้อ 6.3: รายชิ้นใช้ตัวกรองเดียวกัน)
// $2 = ประเภท, $3 = สถานะ, $4 = ภาควิชา (0 = ไม่จำกัด) เป็น parameter ทั้งหมด ไม่ต่อสตริงจากค่าที่ผู้ใช้ส่งมา
// เงื่อนไขภาควิชาอยู่ในทุก query ของรายการ ที่ชั้น repository (plan.md ข้อ 6.1, DOM-ACCESS-01)
const filterWhere = `
  (cardinality($2::text[]) = 0 OR category_name = ANY($2::text[]))
  AND (cardinality($3::text[]) = 0 OR status = ANY($3::text[]))
  AND ($4::int = 0 OR department_id = $4::int)`

// ListModels คืนรายการรุ่น 1 แถว = 1 รุ่น พร้อมจำนวนว่างและจำนวนทั้งหมด (FR-EQS-01)
// จำนวนว่างนับจากสถานะ "ว่าง" อย่างเดียว คำขอยืมที่ค้างไม่ลดจำนวน (DOM-REQ-01)
// กรองที่รายชิ้นก่อนนับ จึงแสดงเฉพาะรุ่นที่มีอย่างน้อย 1 ชิ้นตรงเงื่อนไข และตัวเลขนับเฉพาะชิ้นที่ตรง (FR-EQS-02, Q-15)
// ยังไม่แบ่งหน้า (รอ Q-05 ดู T-23) และยังไม่จำกัดภาควิชา (T-18)
func (r *Repository) ListModels(ctx context.Context, f Filter) ([]ModelRowDTO, error) {
	const q = `
SELECT model_id, model_name, category_name,
       count(DISTINCT item_id) FILTER (WHERE status = $1) AS available_count,
       count(DISTINCT item_id)                            AS total_count
FROM v_equipment_item_status
WHERE` + filterWhere + `
GROUP BY model_id, model_name, category_name
ORDER BY model_name, model_id`
	categories, statuses := f.args()
	rows, err := r.db.Query(ctx, q, string(StatusAvailable), categories, statuses, f.DepartmentID)
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

// ListItems คืนอุปกรณ์รายชิ้นของรุ่นที่เลือก พร้อมสถานะล่าสุดของแต่ละชิ้น เรียงตามรหัสครุภัณฑ์ (FR-EQS-01)
// ใช้ตัวกรองเดียวกับรายการรุ่น (FR-EQS-02, plan.md ข้อ 6.3)
// รุ่นที่ไม่มีอยู่ได้รายการว่าง (plan.md ข้อ 5: ไม่พบ = items ว่าง ไม่ใช่ error)
func (r *Repository) ListItems(ctx context.Context, modelID int, f Filter) ([]ItemRow, error) {
	const q = `
SELECT DISTINCT asset_code, status
FROM v_equipment_item_status
WHERE model_id = $1 AND` + filterWhere + `
ORDER BY asset_code`
	categories, statuses := f.args()
	rows, err := r.db.Query(ctx, q, modelID, categories, statuses, f.DepartmentID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (ItemRow, error) {
		var it ItemRow
		err := row.Scan(&it.AssetCode, &it.Status)
		return it, err
	})
}

// ItemRecord คือข้อมูลอุปกรณ์รายชิ้น 1 ชิ้นตามที่อยู่ในฐานข้อมูล ยังไม่ผ่านกฎของ service
// dueAt มาจากการยืมที่ยังไม่คืน (loan.returned_at IS NULL) ว่าง = ไม่มีการยืมค้าง
type ItemRecord struct {
	AssetCode    string
	Status       Status
	DueAt        *time.Time
	DepartmentID int // ภาควิชาเจ้าของอุปกรณ์ ใช้ตรวจสิทธิ์ใน service ไม่ส่งออกไปใน DTO (DOM-ACCESS-01)
}

// GetItem อ่านอุปกรณ์รายชิ้นจากรหัสครุภัณฑ์ (FR-EQS-04, CON-DATA-03) อ่านอย่างเดียว (CON-DATA-02)
// found = false เมื่อไม่มีรหัสนี้
func (r *Repository) GetItem(ctx context.Context, assetCode string) (rec ItemRecord, found bool, err error) {
	const q = `
SELECT asset_code, status, due_at, department_id
FROM v_equipment_item_status
WHERE asset_code = $1
LIMIT 1`
	err = r.db.QueryRow(ctx, q, assetCode).Scan(&rec.AssetCode, &rec.Status, &rec.DueAt, &rec.DepartmentID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ItemRecord{}, false, nil
	}
	if err != nil {
		return ItemRecord{}, false, err
	}
	return rec, true, nil
}

// CountByStatus นับอุปกรณ์รายชิ้นแยกตาม 3 สถานะ (FR-EQS-07, CON-DATA-01) อ่านอย่างเดียว (CON-DATA-02)
// ยังไม่จำกัดภาควิชาของเจ้าหน้าที่ (รอ Q-P3 ดู T-28 / plan.md ข้อ 6.1 ชั่วคราวให้เห็นทุกภาควิชา)
func (r *Repository) CountByStatus(ctx context.Context) (StatusSummaryDTO, error) {
	const q = `
SELECT count(DISTINCT item_id) FILTER (WHERE status = $1),
       count(DISTINCT item_id) FILTER (WHERE status = $2),
       count(DISTINCT item_id) FILTER (WHERE status = $3),
       count(DISTINCT item_id)
FROM v_equipment_item_status`
	var out StatusSummaryDTO
	err := r.db.QueryRow(ctx, q, string(StatusAvailable), string(StatusBorrowed), string(StatusUnavailable)).
		Scan(&out.Available, &out.Borrowed, &out.Unavailable, &out.Total)
	return out, err
}

// ModelDepartment คืนภาควิชาเจ้าของรุ่น ใช้ตรวจว่านักศึกษาขอรายชิ้นของรุ่นภาควิชาอื่นหรือไม่ (FR-EQS-10, DOM-ACCESS-01)
// found = false เมื่อไม่มีรุ่นนี้
func (r *Repository) ModelDepartment(ctx context.Context, modelID int) (departmentID int, found bool, err error) {
	err = r.db.QueryRow(ctx, `SELECT department_id FROM equipment_model WHERE id = $1`, modelID).Scan(&departmentID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return departmentID, true, nil
}
