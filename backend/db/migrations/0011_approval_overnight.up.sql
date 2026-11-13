-- 0011 คำขอยืมค้างคืนของฟีเจอร์ 004 (specs/004-approve-request/plan.md Draft v2 ข้อ 3, T-40, CR-009)
-- เลข 0011 ยังไม่มีฟีเจอร์ใดใช้ (ตรวจทุก branch 2569-11-13)
-- ไม่แตะ CHECK ของ audit_log: อนุมัติ / ปฏิเสธคำขอค้างคืนใช้ action approve / reject เดิมจาก 0008 (plan ข้อ 3)
-- ไม่แตะตาราง equipment_item: อนุมัติค้างคืนแล้วสถานะอุปกรณ์ยังเป็น "ว่าง" (FR-APV-03 ข, data-model.md แถว UC-14)

-- FR-APV-03 (ข)(ช), Q-APV-44: ชื่ออาจารย์ที่เซ็นอนุมัติ ว่างได้ เพราะคำขอปกติไม่มีชื่ออาจารย์
-- กฎบังคับกรอกและความยาวตรวจใน service (config.ApprovalAdvisorNameMaxLength) ไม่ฝังตัวเลขในฐานข้อมูล
ALTER TABLE approval_decision ADD COLUMN advisor_name TEXT NULL;

-- FR-APV-03 (ข), FR-APV-01 (ฌ): สถานะกลางของคำขอค้างคืนที่อนุมัติ เก็บเป็นประวัติ
-- borrow_request.status ไปที่ approved ทีเดียว ส่วนตารางนี้เก็บ approved_by_officer, pending_advisor_approval, approved ตามลำดับ
CREATE TABLE approval_step (
    id          SERIAL      PRIMARY KEY,
    request_id  INT         NOT NULL REFERENCES borrow_request (id),
    status      TEXT        NOT NULL
        CHECK (status IN ('approved_by_officer', 'pending_advisor_approval', 'approved')),
    actor_id    TEXT        NOT NULL,                       -- FR-APV-01 (ฌ): รหัสเจ้าหน้าที่ (DM-08)
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()          -- FR-APV-01 (ฌ): วันเวลา
);

-- RQ-25, FR-APV-01 (ฌ): ประวัติแก้หรือลบไม่ได้ เหมือน approval_decision ใน 0008
CREATE FUNCTION approval_step_append_only() RETURNS trigger AS $$
BEGIN
    RAISE EXCEPTION 'approval_step is append-only: % is not allowed', TG_OP;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER approval_step_no_update_delete
    BEFORE UPDATE OR DELETE ON approval_step
    FOR EACH ROW EXECUTE FUNCTION approval_step_append_only();

CREATE TRIGGER approval_step_no_truncate
    BEFORE TRUNCATE ON approval_step
    FOR EACH STATEMENT EXECUTE FUNCTION approval_step_append_only();
