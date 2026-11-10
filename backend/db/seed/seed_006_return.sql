-- ข้อมูลตัวอย่างเพิ่มของ 006-return สำหรับลองตามจุดตรวจใน specs/006-return/tasks.md
-- รันหลัง backend/db/seed/seed.sql (ของ 001) และหลัง migration 0006 เป็นข้อมูลตัวอย่าง ไม่ใช่ข้อมูลจริง
-- เพิ่มจาก seed.sql:
--   อุปกรณ์ "ถูกยืม" ของภาควิชา 2 พร้อมการยืมที่ยังไม่คืน 1 รายการ (AC-RET-10 ไม่ต้องเห็น, AC-RET-11 บันทึกคืนไม่ได้)
--   การยืมที่คืนแล้ว 2 รายการของภาควิชา 1 (AC-RET-10 ไม่ต้องเห็นในรายการที่ยังไม่คืน)
-- การยืมที่ยังไม่คืนของภาควิชา 1 ใช้ของ seed.sql (CS-NB-0004, CS-NB-0005)

INSERT INTO equipment_item (id, asset_code, model_id, status) VALUES
    (9, 'XX-NB-0002', 3, 'borrowed');

-- กำหนดคืน 16:00 วันนี้ตามเขตเวลา Asia/Bangkok (glossary, DOM-LOAN-01) แบบเดียวกับ seed.sql
INSERT INTO loan (item_id, borrower_id, loan_type, due_at) VALUES
    (9, 'student-101', 'normal',
        (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') + interval '16 hours') AT TIME ZONE 'Asia/Bangkok');

-- การยืมที่คืนแล้วของ CS-NB-0001 และ CS-NB-0002 (ตอนนี้สถานะ "ว่าง") เมื่อวาน สภาพปกติ คืนทันกำหนด
INSERT INTO loan (item_id, borrower_id, loan_type, due_at, returned_at, return_condition, returned_late) VALUES
    (1, 'student-001', 'normal',
        (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') - interval '1 day' + interval '16 hours') AT TIME ZONE 'Asia/Bangkok',
        (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') - interval '1 day' + interval '15 hours') AT TIME ZONE 'Asia/Bangkok',
        'normal', false),
    (2, 'student-002', 'normal',
        (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') - interval '1 day' + interval '16 hours') AT TIME ZONE 'Asia/Bangkok',
        (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') - interval '1 day' + interval '14 hours') AT TIME ZONE 'Asia/Bangkok',
        'normal', false);

SELECT setval('equipment_item_id_seq', (SELECT max(id) FROM equipment_item));
