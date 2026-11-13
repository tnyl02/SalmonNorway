-- ย้อน 0010 ของ 006-return: คืนรายการ action ของ DM-07 เป็นชุดหลัง 0009 ของ 005 (ไม่มี action ของ 006)
-- และถอดคอลัมน์ตอนรับคืนออกจาก loan
-- ถ้ามีแถว audit_log ที่ใช้ action ของ 006 แล้ว คำสั่งนี้จะล้ม เพราะ audit_log เป็น append-only
ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel', 'request', 'notify',
                      'approve', 'reject', 'auto_reject', 'handover'));

ALTER TABLE loan
    DROP CONSTRAINT loan_return_fields_complete,
    DROP COLUMN returned_late,
    DROP COLUMN return_condition;
