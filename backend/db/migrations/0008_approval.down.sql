-- ย้อน 0008: คืนสภาพให้เหมือนก่อนฟีเจอร์ 004 สร้างตาราง

-- คืนรายการ action ของ DM-07 เป็นชุดก่อน 0008 (ตาม 0005 ของ 003 ถ้า 0006 / 0007 เพิ่มค่า ต้องแก้รายการนี้ให้ตรง)
-- ถ้ามีแถว audit_log ที่ใช้ค่า approve / reject / auto_reject แล้ว คำสั่งนี้จะล้ม เพราะ audit_log เป็น append-only
ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel', 'request', 'notify'));

DROP TABLE approval_email;
DROP TABLE approval_decision;
DROP FUNCTION approval_decision_append_only();
