-- ข้อมูลตัวอย่างเพิ่มของ 006-return สำหรับลองตามจุดตรวจใน specs/006-return/tasks.md (ทีมเลือกเอง ไม่ได้มาจาก spec)
-- เป็นข้อมูลตัวอย่าง ไม่ใช่ข้อมูลจริง ใช้หลัง seed.sql ของ 001 และ migration 0010 แล้วเท่านั้น
-- โหลดเอง: docker compose -f deploy/docker-compose.yml exec -T db psql -U salmon -d salmon < backend/db/seed/seed_006_return.sql
-- ไม่กำหนด id ตายตัว อ้างอุปกรณ์ด้วยรหัสครุภัณฑ์ จึงโหลดร่วมกับ seed ของฟีเจอร์อื่นได้ไม่ชนกัน
-- เพิ่มจาก seed.sql:
--   อุปกรณ์ "ถูกยืม" ของภาควิชา 2 พร้อมการยืมที่ยังไม่คืน 1 รายการ (AC-RET-10 ไม่ต้องเห็น, AC-RET-11 บันทึกคืนไม่ได้)
--   การยืมที่คืนแล้ว 2 รายการของภาควิชา 1 (AC-RET-10 ไม่ต้องเห็นในรายการที่ยังไม่คืน)
-- การยืมที่ยังไม่คืนของภาควิชา 1 ใช้ของ seed.sql (CS-NB-0004, CS-NB-0005)

INSERT INTO equipment_item (asset_code, model_id, status) VALUES
    ('XX-NB-0002', 3, 'borrowed');

-- กำหนดคืน 16:00 วันนี้ตามเขตเวลา Asia/Bangkok (glossary, DOM-LOAN-01) แบบเดียวกับ seed.sql
INSERT INTO loan (item_id, borrower_id, loan_type, due_at)
SELECT id, 'student-101', 'normal',
       (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') + interval '16 hours') AT TIME ZONE 'Asia/Bangkok'
FROM equipment_item WHERE asset_code = 'XX-NB-0002';

-- การยืมที่คืนแล้วของ CS-NB-0001 และ CS-NB-0002 (ตอนนี้สถานะ "ว่าง") เมื่อวาน สภาพปกติ คืนทันกำหนด
INSERT INTO loan (item_id, borrower_id, loan_type, due_at, returned_at, return_condition, returned_late)
SELECT i.id, v.borrower_id, 'normal',
       (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') - interval '1 day' + interval '16 hours') AT TIME ZONE 'Asia/Bangkok',
       (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') - interval '1 day' + v.returned_hour) AT TIME ZONE 'Asia/Bangkok',
       'normal', false
FROM (VALUES ('CS-NB-0001', 'student-001', interval '15 hours'),
             ('CS-NB-0002', 'student-002', interval '14 hours')) AS v(asset_code, borrower_id, returned_hour)
JOIN equipment_item i ON i.asset_code = v.asset_code;
