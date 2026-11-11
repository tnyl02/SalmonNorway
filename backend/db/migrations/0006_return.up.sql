-- 0006 การรับคืนอุปกรณ์ของ 006-return (UC-09) ตาม specs/006-return/plan.md ข้อ 3
-- เลข 0006 เจ้าของ 006 จองไว้ใน specs/000-shared/data-model.md หัวข้อ "เลข migration"
-- ห้ามแก้ไฟล์ 0001, 0002 จึงเพิ่มคอลัมน์และค่า CHECK ด้วยไฟล์นี้

-- DM-05 การยืม: ข้อมูลตอนรับคืน ว่างได้ทั้งหมด เพื่อไม่ให้กระทบ 001 (อ่าน) และ 005 (สร้างการยืม)
ALTER TABLE loan
    -- FR-RET-02 สภาพตอนรับคืน 2 ค่า ปกติ / ชำรุด (สูญหายไม่ใช่สภาพ)
    ADD COLUMN return_condition TEXT NULL
        CHECK (return_condition IN ('normal', 'damaged')),
    -- FR-RET-06 คืนเกินกำหนด = returned_at > due_at บันทึกตอนรับคืน ไม่คำนวณใหม่ภายหลัง
    ADD COLUMN returned_late BOOLEAN NULL,
    -- FR-RET-02 คืนแล้วต้องมีเวลาคืน สภาพ และเครื่องหมายคืนเกินกำหนดครบ / ยังไม่คืนต้องว่างทั้ง 3 ค่า
    ADD CONSTRAINT loan_return_fields_complete CHECK (
        (returned_at IS NULL) = (return_condition IS NULL)
        AND (returned_at IS NULL) = (returned_late IS NULL)
    );

-- ไม่มีตารางรูปตอนรับคืน: FR-RET-05 ยกเลิก 2569-11-11 (RQ-04 ฝั่งรับคืนอยู่ใน BL-16 Could ไม่อยู่ใน prototype)

-- DM-07 audit_log: เพิ่มค่า action ของ 006 (DOM-AUDIT-01, FR-RET-07) ตามตาราง action ใน data-model.md
--   return_list   เปิดหน้ารายการที่ยังไม่คืน
--   return_lookup ค้นการยืมจากรหัสครุภัณฑ์
--   return_detail เปิดรายละเอียดการยืม
--   return        บันทึกการรับคืน
ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel',
                      'return_list', 'return_lookup', 'return_detail', 'return'));
