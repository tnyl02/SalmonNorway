-- ย้อน 0009: คืนรายการ action ของ DM-07 เป็นชุดหลัง 0008 ของ 004 (ไม่มี handover)
-- ถ้ามีแถว audit_log ที่ใช้ค่า handover แล้ว คำสั่งนี้จะล้ม เพราะ audit_log เป็น append-only
ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel', 'request', 'notify',
                      'approve', 'reject', 'auto_reject'));
