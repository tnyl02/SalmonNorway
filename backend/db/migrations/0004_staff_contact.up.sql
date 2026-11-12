-- 0004 ตาราง staff_contact ของฟีเจอร์ 003 (จำลอง ASM-10) เลขจองในกลุ่ม 2569-11-11
-- Email เจ้าหน้าที่ของแต่ละภาควิชา ใช้เฉพาะ Notifier จำลองหาผู้รับเมื่อมีการยกเลิกคำขอ (FR-CXL-09, IF-NOT-01)
-- ระบบจริงยังไม่รู้ว่ารายชื่อมาจากไหน (Q-15) ข้อมูลในตารางนี้มาจาก seed ของ 003 (T-02) ทีมเลือกเอง ไม่ได้มาจาก spec
-- ไม่แตะตาราง equipment_item และ borrow_request (CON-CXL-01, CON-DATA-01)

CREATE TABLE staff_contact (
    department_id INT  REFERENCES department (id),   -- ภาควิชาเจ้าของอุปกรณ์ (DM-01)
    email         TEXT NOT NULL                      -- ผู้รับ Email แจ้งการยกเลิก (FR-CXL-09)
);
