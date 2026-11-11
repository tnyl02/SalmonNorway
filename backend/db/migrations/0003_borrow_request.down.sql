-- ย้อน 0003: คืนสภาพให้เหมือนก่อนฟีเจอร์ 002 สร้างตาราง borrow_request

-- คืนรายการ action ของ DM-07 เป็นชุดเดิมที่ 0002 สร้างไว้ (ยังไม่มีค่า 'request' ของ 002)
ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel'));

ALTER TABLE audit_log DROP CONSTRAINT fk_audit_log_request;

DROP TABLE borrow_request;
