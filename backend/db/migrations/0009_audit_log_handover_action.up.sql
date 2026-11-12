-- 0009 เพิ่มค่า action 'handover' ใน audit_log ของฟีเจอร์ 005 (DM-07, FR-HND-08, Q-HND-06 ปิด 2569-11-12)
-- บันทึกการส่งมอบทั้งสำเร็จ (outcome 'ok') และไม่สำเร็จ ('forbidden' / 'error') ตาม Q-HND-14
-- ใช้เลขหลัง 0008 ของ 004 เพราะ 0008 สร้าง CHECK ใหม่ ถ้าเพิ่มค่าก่อน 0008 ค่าจะหาย (ทีมเลือก 2569-11-12 ไม่แก้ไฟล์ของ 004)
-- ห้ามแก้ migration เดิม ใช้ชื่อ constraint เดิม audit_log_action_check แบบเดียวกับ 0003, 0005, 0008
-- ต้องใส่ค่าครบทุกค่าที่มีก่อนหน้า: list, search, detail (001), cancel (003 ใน 0002), request (002 ใน 0003),
-- notify (003 ใน 0005), approve, reject, auto_reject (004 ใน 0008)

ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel', 'request', 'notify',
                      'approve', 'reject', 'auto_reject', 'handover'));
