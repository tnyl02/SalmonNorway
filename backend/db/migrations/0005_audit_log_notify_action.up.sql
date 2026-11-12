-- 0005 เพิ่มค่า action 'notify' ใน audit_log ของฟีเจอร์ 003 (DM-07, ASM-12) เลขจองในกลุ่ม 2569-11-11
-- ใช้บันทึกเมื่อส่ง Email แจ้งเจ้าหน้าที่ไม่สำเร็จ (FR-CXL-09, outcome 'error') ฟีเจอร์อื่นใช้ร่วมได้ เช่น 004 แจ้งผลอนุมัติ (RQ-21)
-- ห้ามแก้ 0002 ตามกติกาในไฟล์นั้น ใช้ชื่อ constraint เดิมที่ PostgreSQL ตั้งให้ใน 0002 (ตกลงกับ 001 และ 002 2569-11-11)
-- ต้องใส่ค่าครบทุกค่าที่มีก่อนหน้า: list, search, detail (001), cancel (003 ใน 0002), request (002 ใน 0003)

ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel', 'request', 'notify'));
