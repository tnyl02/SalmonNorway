-- ข้อมูลตั้งต้นสำหรับลองระบบในเครื่องพัฒนา (ASM-01: ระบบจริงให้เจ้าหน้าที่นำเข้าเอง)
-- เป็นข้อมูลตัวอย่าง ไม่ใช่ข้อมูลอุปกรณ์จริงของภาควิชา
-- ออกแบบให้ใช้ลองตามจุดตรวจใน tasks.md ของ 001 ได้:
--   รุ่น 1 มี 5 ชิ้น ว่าง 3 ชิ้น (AC-EQS-01), มีคำว่า Notebook (AC-EQS-02)
--   ถูกยืมแบบยังไม่เลยกำหนด และเลยกำหนดแล้ว (AC-EQS-04, AC-EQS-14), ไม่พร้อมใช้งาน (AC-EQS-06)
--   อุปกรณ์ของภาควิชาอื่น (AC-EQS-13)

INSERT INTO department (id, name) VALUES
    (1, 'ภาควิชาคอมพิวเตอร์'),
    (2, 'ภาควิชาตัวอย่างอื่น');

INSERT INTO category (id, name) VALUES
    (1, 'คอมพิวเตอร์พกพา'),
    (2, 'อุปกรณ์นำเสนอ');

INSERT INTO equipment_model (id, name, category_id, department_id) VALUES
    (1, 'Notebook ตัวอย่างรุ่น A', 1, 1),
    (2, 'Projector ตัวอย่างรุ่น B', 2, 1),
    (3, 'Notebook ตัวอย่างรุ่น C', 1, 2);

INSERT INTO equipment_item (id, asset_code, model_id, status) VALUES
    (1, 'CS-NB-0001', 1, 'available'),
    (2, 'CS-NB-0002', 1, 'available'),
    (3, 'CS-NB-0003', 1, 'available'),
    (4, 'CS-NB-0004', 1, 'borrowed'),
    (5, 'CS-NB-0005', 1, 'borrowed'),
    (6, 'CS-PJ-0001', 2, 'available'),
    (7, 'CS-PJ-0002', 2, 'unavailable'),
    (8, 'XX-NB-0001', 3, 'available');

-- กำหนดคืนคิดตามเขตเวลา Asia/Bangkok (glossary) และยืมปกติคืน 16:00 (DOM-LOAN-01)
-- ชิ้นที่ 4: ยืมปกติ กำหนดคืน 16:00 วันนี้ (ก่อน 16:00 ยังไม่เลยกำหนด)
-- ชิ้นที่ 5: ยืมปกติ กำหนดคืน 16:00 เมื่อวาน (เลยกำหนดแล้ว)
INSERT INTO loan (item_id, borrower_id, loan_type, due_at) VALUES
    (4, 'student-001', 'normal',
        (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') + interval '16 hours') AT TIME ZONE 'Asia/Bangkok'),
    (5, 'student-002', 'normal',
        (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') - interval '1 day' + interval '16 hours') AT TIME ZONE 'Asia/Bangkok');

-- ให้ id ที่สร้างอัตโนมัติต่อจากข้อมูลตั้งต้น
SELECT setval('department_id_seq', (SELECT max(id) FROM department));
SELECT setval('category_id_seq', (SELECT max(id) FROM category));
SELECT setval('equipment_model_id_seq', (SELECT max(id) FROM equipment_model));
SELECT setval('equipment_item_id_seq', (SELECT max(id) FROM equipment_item));
