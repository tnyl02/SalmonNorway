-- 0008 ตารางของฟีเจอร์ 004 อนุมัติคำขอยืมปกติ (specs/004-approve-request/plan.md ข้อ 3, T-01)
-- เลข 0008 ทีมจองให้ 2569-11-12 ค่า action ใหม่กำหนดในไฟล์นี้ ไม่เขียนใน data-model.md (ทีมตัดสิน 2569-11-12)
-- ฟีเจอร์นี้ไม่เปลี่ยนสถานะอุปกรณ์ (CON-DATA-01, data-model.md แถว UC-12) จึงไม่มีคำสั่งใดแตะตาราง equipment_item

-- ผลการพิจารณา 1 แถวต่อคำขอ = ประวัติตาม RQ-25 (FR-APV-01 ฌ) และข้อมูลใน Email (FR-APV-02)
CREATE TABLE approval_decision (
    id                  SERIAL      PRIMARY KEY,
    request_id          INT         NOT NULL UNIQUE                 -- UNIQUE: พิจารณาซ้ำหรือเปลี่ยนผลไม่ได้ (Out of scope "เปลี่ยนผลการพิจารณา", FR-APV-01 ช)
        REFERENCES borrow_request (id),
    result              TEXT        NOT NULL
        CHECK (result IN ('approved', 'rejected', 'auto_rejected')),  -- อนุมัติ / ปฏิเสธ (รวมปฏิเสธตาม จ) / ปฏิเสธอัตโนมัติ (ฉ)
    reason              TEXT        NULL,                           -- FR-APV-01 (ง)(จ)(ฉ): เหตุผล ว่างได้เฉพาะการอนุมัติ
    actor_id            TEXT        NOT NULL,                       -- FR-APV-01 (ฌ): รหัสเจ้าหน้าที่ (DM-08) หรือ 'system'
    actor_name          TEXT        NULL,                           -- FR-APV-01 (ฌ): ว่างได้ เพราะ Principal ยังไม่มีชื่อ (plan ข้อ 8, T-38)
    actor_department_id INT         NOT NULL REFERENCES department (id),  -- ASM-06: ภาควิชาเจ้าของอุปกรณ์ ใช้กรองรายการ Email ที่ส่งไม่สำเร็จ
    decided_at          TIMESTAMPTZ NOT NULL DEFAULT now(),         -- FR-APV-01 (ฌ): วันเวลา
    -- FR-APV-01 (ง): ปฏิเสธต้องมีเหตุผล (ไม่นับช่องว่างล้วน)
    CONSTRAINT approval_decision_reason_required
        CHECK (result = 'approved' OR (reason IS NOT NULL AND btrim(reason) <> ''))
);

-- ประวัติแก้หรือลบไม่ได้ (plan ข้อ 3: ไม่มี UPDATE / DELETE) ใช้ trigger แบบเดียวกับ audit_log ใน 0002
CREATE FUNCTION approval_decision_append_only() RETURNS trigger AS $$
BEGIN
    RAISE EXCEPTION 'approval_decision is append-only: % is not allowed', TG_OP;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER approval_decision_no_update_delete
    BEFORE UPDATE OR DELETE ON approval_decision
    FOR EACH ROW EXECUTE FUNCTION approval_decision_append_only();

CREATE TRIGGER approval_decision_no_truncate
    BEFORE TRUNCATE ON approval_decision
    FOR EACH STATEMENT EXECUTE FUNCTION approval_decision_append_only();

-- Email แจ้งผล 1 ฉบับต่อผลการพิจารณา เขียนใน transaction เดียวกับผล (outbox) แล้วตัวส่งเบื้องหลังส่ง (FR-APV-02, IF-NOT-01)
CREATE TABLE approval_email (
    id              SERIAL      PRIMARY KEY,
    request_id      INT         NOT NULL UNIQUE                     -- 1 คำขอมีผลเดียว จึงมี Email เดียว (AC-APV-10: Email 1 ฉบับ)
        REFERENCES approval_decision (request_id),
    kind            TEXT        NOT NULL
        CHECK (kind IN ('approved', 'rejected', 'auto_rejected')),  -- FR-APV-02 (ก)(ข)(ค)
    status          TEXT        NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'sent', 'failed')),            -- FR-APV-02 (ง): Failed
    attempts        INT         NOT NULL DEFAULT 0 CHECK (attempts >= 0),  -- FR-APV-02 (ง): ครั้งแรก + ส่งซ้ำ 3 ครั้ง (Q-APV-17)
    last_error      TEXT        NULL,                               -- FR-APV-02 (ง): สาเหตุ
    last_attempt_at TIMESTAMPTZ NULL,                               -- FR-APV-02 (ง): เวลา
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),             -- FR-APV-02 (ง): ห่างกันครั้งละ 1 นาที
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()              -- FR-APV-02: นับ 60 วินาทีหลังบันทึกผล
);

-- ตัวส่งเบื้องหลังหยิบ Email ที่ถึงเวลาส่ง (T-15, T-32)
CREATE INDEX idx_approval_email_due
    ON approval_email (status, next_attempt_at);

-- RQ-25: เพิ่มค่า action ของ 004 ใน audit_log (DM-07) ด้วยวิธีเดียวกับ 0003 ห้ามแก้ไฟล์ 0002 / 0003
-- รายการต้องรวมค่าของทุก migration ก่อนหน้า: list, search, detail (001), cancel (003 ใน 0002), request (002 ใน 0003),
-- notify (003 ใน 0005 แจ้งจาก 003 ก่อนเปิด PR 2569-11-12) ถ้า 0006 / 0007 เพิ่มค่าภายหลัง ต้องเพิ่มที่นี่ด้วย
ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel', 'request', 'notify',
                      'approve', 'reject', 'auto_reject'));
