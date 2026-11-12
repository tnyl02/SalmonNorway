-- 0007 เพิ่มค่า action 'handover' ใน audit_log ของฟีเจอร์ 005 (DM-07, FR-HND-08, Q-HND-06 ปิด 2569-11-12) เลขจองในกลุ่ม 2569-11-11
-- บันทึกการส่งมอบทั้งสำเร็จ (outcome 'ok') และไม่สำเร็จ ('forbidden' / 'error') ตาม Q-HND-14
-- ห้ามแก้ 0002 ตามกติกาในไฟล์นั้น ใช้ชื่อ constraint เดิม audit_log_action_check แบบเดียวกับ 0003, 0005
-- ต้องใส่ค่าครบทุกค่าที่มีก่อนหน้า: list, search, detail (001), cancel (003 ใน 0002), request (002 ใน 0003), notify (003 ใน 0005)
-- 0008 ของ 004 สร้าง CHECK ใหม่ภายหลัง จึงเพิ่ม 'handover' ในรายการของ 0008 ด้วย (ตามคอมเมนต์ในไฟล์ 0008)

ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel', 'request', 'notify', 'handover'));
