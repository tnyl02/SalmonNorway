-- 0010 การรับคืนอุปกรณ์ของ 006-return (UC-09) ตาม specs/006-return/plan.md ข้อ 3
-- เลข 0010 ต่อจาก 0009 ของ 005 (เดิม 0006 แล้ว 0009 ย้ายเมื่อ 2569-11-12 เพราะ 0009 ชนกับ 005 ตอนรวม main)
-- ห้ามแก้ไฟล์ 0001, 0002 จึงเพิ่มคอลัมน์และค่า CHECK ด้วยไฟล์นี้

-- DM-05 การยืม: ข้อมูลตอนรับคืน ว่างได้ทั้งหมด เพื่อไม่ให้กระทบ 001 (อ่าน) และ 005 (สร้างการยืม)
ALTER TABLE loan
    -- FR-RET-02 สภาพตอนรับคืน 2 ค่า ปกติ / ชำรุด (สูญหายไม่ใช่สภาพ)
    ADD COLUMN return_condition TEXT NULL
        CHECK (return_condition IN ('normal', 'damaged')),
    -- FR-RET-06 คืนเกินกำหนด = returned_at > due_at บันทึกตอนรับคืน ไม่คำนวณใหม่ภายหลัง
    ADD COLUMN returned_late BOOLEAN NULL,
    -- FR-RET-02, FR-RET-06 สภาพและเครื่องหมายคืนเกินกำหนดมีได้เฉพาะการยืมที่คืนแล้ว (มี returned_at)
    -- ไม่บังคับกลับทาง (คืนแล้วต้องมีครบ) เพราะ test ของ 001, 002 จำลองการคืนด้วย returned_at อย่างเดียว
    -- การคืนผ่าน 006 ใส่ครบทั้ง 3 ค่าเสมอ (repository.go: RecordReturn) ผ่อนเมื่อรวม main 2569-11-12
    ADD CONSTRAINT loan_return_fields_complete CHECK (
        (return_condition IS NULL OR returned_at IS NOT NULL)
        AND (returned_late IS NULL OR returned_at IS NOT NULL)
    );

-- ไม่มีตารางรูปตอนรับคืน: FR-RET-05 ยกเลิก 2569-11-11 (RQ-04 ฝั่งรับคืนอยู่ใน BL-16 Could ไม่อยู่ใน prototype)

-- DM-07 audit_log: เพิ่มค่า action ของ 006 (DOM-AUDIT-01, FR-RET-07) ตามตาราง action ใน data-model.md
--   return_list   เปิดหน้ารายการที่ยังไม่คืน
--   return_lookup ค้นการยืมจากรหัสครุภัณฑ์
--   return_detail เปิดรายละเอียดการยืม
--   return        บันทึกการรับคืน
-- CHECK ใหม่มีค่าเดิมเฉพาะที่อยู่ใน DM-07 ของ branch นี้ตอนนี้ (0002 ของ 001: list, search, detail, cancel)
-- ยังไม่ใส่ค่าของ 0005 (003), 0007 (005), 0008 (004) เพราะยังไม่ merge และ test ของ 001 ปฏิเสธค่าที่ไม่อยู่ใน DM-07
-- เมื่อฟีเจอร์เหล่านั้น merge แล้ว 0009 รันหลังสุดจะทับ CHECK ของเขา -> ต้องรวมค่าของเขาเพิ่มที่นี่ (ตกลงในกลุ่มตอน merge)
-- ต้องคงค่าเดิมของทุกฟีเจอร์ (ล่าสุดคือ 0009 ของ 005) แล้วเพิ่ม 4 ค่าของ 006 ห้ามตัดค่าของคนอื่น
ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel', 'request', 'notify',
                      'approve', 'reject', 'auto_reject', 'handover',
                      'return_list', 'return_lookup', 'return_detail', 'return'));
