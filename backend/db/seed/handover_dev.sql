-- ข้อมูลตัวอย่างสำหรับลองฟีเจอร์ 005 ตามจุดตรวจใน specs/005-handover/tasks.md (ทีมเลือกเอง ไม่ได้มาจาก spec)
-- เป็นข้อมูลตัวอย่าง ไม่ใช่ข้อมูลจริง ใช้หลัง seed.sql ของ 001 โหลดแล้วเท่านั้น
-- โหลดเอง: docker compose -f deploy/docker-compose.yml exec -T db psql -U salmon -d salmon < backend/db/seed/handover_dev.sql
-- เพิ่มรุ่นและอุปกรณ์ชุดใหม่ของตัวเอง ไม่แตะอุปกรณ์ของ seed.sql ตัวเลขที่ใช้ลองของ 001 จึงไม่เปลี่ยน
--
-- AC-HND-14: คำขออนุมัติแล้วและยังไม่ส่งมอบ 2 รายการ (อุปกรณ์คนละชิ้น) ยังไม่อนุมัติ 1 รายการ ส่งมอบแล้ว 1 รายการ

INSERT INTO equipment_model (name, category_id, department_id) VALUES
    ('อุปกรณ์ทดลองส่งมอบ (005)', 1, 1);

INSERT INTO equipment_item (asset_code, model_id, status)
SELECT code, (SELECT id FROM equipment_model WHERE name = 'อุปกรณ์ทดลองส่งมอบ (005)'), st
FROM (VALUES
    ('HO-TEST-0001', 'available'),   -- คำขออนุมัติแล้ว รอส่งมอบ
    ('HO-TEST-0002', 'available'),   -- คำขออนุมัติแล้ว รอส่งมอบ
    ('HO-TEST-0003', 'available'),   -- คำขอยังไม่อนุมัติ
    ('HO-TEST-0004', 'borrowed')     -- ส่งมอบแล้ว (DOM-REQ-01: เปลี่ยนเป็นถูกยืมตอนส่งมอบ)
) AS v(code, st);

-- คำขอยืมแบบปกติ (ตาราง borrow_request ของ 002) ผู้ขอคือ student-001 ของ auth.MockProvider
INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type, purpose, start_at, end_at, status)
SELECT 'student-001', 'นักศึกษาตัวอย่าง', 'student-001', i.id, 'normal', 'study',
       (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') + interval '9 hours') AT TIME ZONE 'Asia/Bangkok',
       (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') + interval '16 hours') AT TIME ZONE 'Asia/Bangkok',
       v.st
FROM (VALUES
    ('HO-TEST-0001', 'approved'),
    ('HO-TEST-0002', 'approved'),
    ('HO-TEST-0003', 'requested'),
    ('HO-TEST-0004', 'handed_over')
) AS v(code, st)
JOIN equipment_item i ON i.asset_code = v.code;

-- คำขอที่ส่งมอบแล้ว ต้องมี Loan และแถว handover คู่กัน (FR-HND-07)
INSERT INTO loan (item_id, borrower_id, loan_type, due_at)
SELECT id, 'student-001', 'normal',
       (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') + interval '16 hours') AT TIME ZONE 'Asia/Bangkok'
FROM equipment_item WHERE asset_code = 'HO-TEST-0004';

INSERT INTO handover (loan_id, request_id, handed_over_at, handed_over_by,
                      borrower_student_id, borrower_first_name, borrower_last_name, identity_method, card_checked)
SELECT l.id, r.id, now(), 'staff-001', 'student-001', 'นักศึกษา', 'ตัวอย่าง', 'auth_api', true
FROM equipment_item i
JOIN loan l ON l.item_id = i.id AND l.returned_at IS NULL
JOIN borrow_request r ON r.item_id = i.id AND r.status = 'handed_over'
WHERE i.asset_code = 'HO-TEST-0004';
