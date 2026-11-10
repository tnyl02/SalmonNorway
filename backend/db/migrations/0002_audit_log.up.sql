-- 0002 ตาราง audit_log ใช้ร่วมทุกฟีเจอร์ (DM-07 ใน specs/000-shared/data-model.md, DOM-AUDIT-01)
-- สร้างโดย 001 ใน T-31 ฟีเจอร์ที่ต้องการค่า action ใหม่ ให้เพิ่มใน CHECK ด้วย migration ใหม่ ห้ามแก้ไฟล์นี้

CREATE TABLE audit_log (
    id          BIGSERIAL PRIMARY KEY,
    actor_id    TEXT        NOT NULL,                 -- ผู้เข้าถึง (DM-08 รหัสผู้ใช้)
    actor_role  TEXT        NOT NULL,                 -- บทบาทของผู้เข้าถึง
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now(),   -- เวลา
    action      TEXT        NOT NULL
        CHECK (action IN ('list', 'search', 'detail', 'cancel')), -- ค่าตามตาราง action ใน DM-07
    query       JSONB       NULL,                     -- เงื่อนไขที่ค้น (DOM-AUDIT-01)
    asset_code  TEXT        NULL,                     -- รหัสครุภัณฑ์ (ว่างได้)
    request_id  INT         NULL,                     -- รหัสคำขอยืม (ว่างได้ ขอโดย 003) FK ให้ 002 เพิ่มเมื่อมีตาราง borrow_request
    outcome     TEXT        NOT NULL
        CHECK (outcome IN ('ok', 'forbidden', 'error'))
);

-- append-only: บันทึกแล้วแก้หรือลบไม่ได้ (plan.md ข้อ 4)
-- ใช้ trigger เพราะบังคับได้ทุกผู้ใช้ฐานข้อมูล รวมถึงเครื่องพัฒนาที่ใช้ผู้ใช้คนเดียว
CREATE FUNCTION audit_log_append_only() RETURNS trigger AS $$
BEGIN
    RAISE EXCEPTION 'audit_log is append-only: % is not allowed', TG_OP;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_log_no_update_delete
    BEFORE UPDATE OR DELETE ON audit_log
    FOR EACH ROW EXECUTE FUNCTION audit_log_append_only();

CREATE TRIGGER audit_log_no_truncate
    BEFORE TRUNCATE ON audit_log
    FOR EACH STATEMENT EXECUTE FUNCTION audit_log_append_only();
