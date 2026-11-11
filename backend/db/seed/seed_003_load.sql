-- ข้อมูลสำหรับ load test ของฟีเจอร์ 003 (T-29, NFR-PERF-CXL-01, AC-CXL-17) ใช้กับ tests/load/cancel_request.js
-- โหลดหลัง migration และ seed.sql ของ 001 ควรใช้ฐานข้อมูลแยก (เช่น salmon_load) ไม่ปนกับข้อมูลที่ใช้ลองหน้าเว็บ
-- เป็นข้อมูลทดสอบ ไม่ใช่ข้อมูลจริง ทีมเลือกเอง ไม่ได้มาจาก spec
--
-- auth จำลองมีนักศึกษาภาควิชา 1 แค่ student-001 และ student-002 และตาราง 0003 ห้ามผู้ขอคนเดิมมีคำขอที่ยังไม่จบซ้ำชิ้น
-- (uq_request_open_per_item) จึงสร้างอุปกรณ์ทดสอบของตัวเอง 1,000 ชิ้น (LD-0001 ถึง LD-1000) ไม่แตะอุปกรณ์ของ seed.sql
-- คำขอ 50001-50500 เป็นของ student-001 และ 50501-51000 เป็นของ student-002 สถานะ requested ทั้งหมด
--
-- โหลดซ้ำได้ และโหลดซ้ำ = รีเซ็ตคำขอกลับเป็น requested ให้รัน load test รอบใหม่ได้
-- (ใช้ UPDATE ไม่ใช้ DELETE เพราะเมื่อมี audit_log อ้างถึงคำขอ (T-31) จะลบคำขอไม่ได้)

INSERT INTO equipment_model (id, name, category_id, department_id)
VALUES (5001, 'อุปกรณ์ทดสอบโหลด (003)', 1, 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO equipment_item (id, asset_code, model_id, status)
SELECT 50000 + n, 'LD-' || lpad(n::text, 4, '0'), 5001, 'available'
FROM generate_series(1, 1000) AS n
ON CONFLICT (id) DO NOTHING;

INSERT INTO borrow_request
    (id, requester_id, requester_name, student_code, item_id, loan_type, purpose,
     start_at, end_at, status, submitted_at)
SELECT 50000 + n,
       CASE WHEN n <= 500 THEN 'student-001' ELSE 'student-002' END,
       CASE WHEN n <= 500 THEN 'ผู้ทดสอบโหลด หนึ่ง' ELSE 'ผู้ทดสอบโหลด สอง' END,
       CASE WHEN n <= 500 THEN 'LOAD-001' ELSE 'LOAD-002' END,
       50000 + n, 'normal', 'study',
       now() + interval '1 day', now() + interval '1 day 6 hours',
       'requested', now() - (n || ' minutes')::interval
FROM generate_series(1, 1000) AS n
ON CONFLICT (id) DO UPDATE SET status = 'requested', cancel_reason = NULL, cancelled_at = NULL;

-- Email เจ้าหน้าที่ภาควิชา 1 ให้ Notifier จำลองมีผู้รับ (เหมือนระบบจริงที่ต้องส่ง Email ทุกครั้งที่ยกเลิก FR-CXL-09)
INSERT INTO staff_contact (department_id, email)
SELECT 1, 'staff-it@example.test'
WHERE NOT EXISTS (SELECT 1 FROM staff_contact WHERE department_id = 1);

SELECT setval('equipment_model_id_seq', GREATEST((SELECT max(id) FROM equipment_model), 1));
SELECT setval('equipment_item_id_seq', GREATEST((SELECT max(id) FROM equipment_item), 1));
SELECT setval('borrow_request_id_seq', GREATEST((SELECT max(id) FROM borrow_request), 1));
