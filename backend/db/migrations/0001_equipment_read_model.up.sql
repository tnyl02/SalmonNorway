-- 0001 ตารางข้อมูลอุปกรณ์และการยืม (read model ของ 001-equipment-status)
-- ตามโครงใน specs/001-equipment-status/plan.md ข้อ 4 และ specs/000-shared/data-model.md
-- เจ้าของข้อมูลจริงคือ UC-10 (อุปกรณ์) และ UC-08/09 (การยืม) ฟีเจอร์นั้นรับช่วงตารางเหล่านี้ต่อได้

-- DM-01 ภาควิชา
CREATE TABLE department (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- DM-02 ประเภทอุปกรณ์ (ตายตัวหรือเพิ่มได้ รอ Q-04)
CREATE TABLE category (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

-- DM-03 รุ่น/ชนิด
CREATE TABLE equipment_model (
    id            SERIAL PRIMARY KEY,
    name          TEXT NOT NULL,
    category_id   INT  NOT NULL REFERENCES category (id),
    department_id INT  NOT NULL REFERENCES department (id)
);

-- DM-04 อุปกรณ์รายชิ้น
CREATE TABLE equipment_item (
    id         SERIAL PRIMARY KEY,
    asset_code TEXT        NOT NULL UNIQUE,                    -- CON-DATA-03 รหัสครุภัณฑ์ไม่ซ้ำกัน
    model_id   INT         NOT NULL REFERENCES equipment_model (id),
    status     TEXT        NOT NULL
        CHECK (status IN ('available', 'borrowed', 'unavailable')), -- CON-DATA-01 ว่าง / ถูกยืม / ไม่พร้อมใช้งาน
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- DM-05 การยืม (สร้างโดย UC-08 ปิดโดย UC-09)
CREATE TABLE loan (
    id          SERIAL PRIMARY KEY,
    item_id     INT         NOT NULL REFERENCES equipment_item (id),
    borrower_id TEXT        NOT NULL,                          -- รหัสผู้ใช้จากระบบยืนยันตัวตน (DM-08)
    loan_type   TEXT        NOT NULL CHECK (loan_type IN ('normal', 'overnight')), -- DM-05 ยืมปกติ / ค้างคืน
    due_at      TIMESTAMPTZ NOT NULL,                          -- กำหนดคืน วันที่และเวลา (CR-002, DOM-LOAN-01)
    returned_at TIMESTAMPTZ NULL                               -- ว่าง = ยังไม่คืน
);

-- index ตาม plan.md ข้อ 4 (index สำหรับค้นหาด้วย pg_trgm อยู่ใน T-22)
CREATE INDEX idx_item_model_status ON equipment_item (model_id, status);
CREATE INDEX idx_model_dept_category ON equipment_model (department_id, category_id);
CREATE INDEX idx_loan_open ON loan (item_id) WHERE returned_at IS NULL;

-- view สำหรับอ่าน ซ่อน join ไว้ที่เดียว (plan.md ข้อ 4)
-- due_at และ borrower_id มีค่าเฉพาะอุปกรณ์ที่มีการยืมที่ยังไม่คืน
CREATE VIEW v_equipment_item_status AS
SELECT i.id            AS item_id,
       i.asset_code,
       i.status,
       i.updated_at,
       m.id            AS model_id,
       m.name          AS model_name,
       c.id            AS category_id,
       c.name          AS category_name,
       d.id            AS department_id,
       d.name          AS department_name,
       l.due_at,
       l.borrower_id
FROM equipment_item i
JOIN equipment_model m ON m.id = i.model_id
JOIN category c ON c.id = m.category_id
JOIN department d ON d.id = m.department_id
LEFT JOIN loan l ON l.item_id = i.id AND l.returned_at IS NULL;
