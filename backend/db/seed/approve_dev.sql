-- ข้อมูลตัวอย่างสำหรับลองฟีเจอร์ 004 อนุมัติคำขอยืมปกติ ตามจุดตรวจใน specs/004-approve-request/tasks.md (T-03)
-- ทีมเลือกเอง ไม่ได้มาจาก spec เป็นข้อมูลตัวอย่าง ไม่ใช่ข้อมูลจริง
-- ใช้หลัง migration (ถึง 0008) และ seed.sql ของ 001 โหลดแล้วเท่านั้น โหลดร่วมกับ handover_dev.sql ของ 005 ได้
-- โหลดเอง: docker compose -f deploy/docker-compose.yml exec -T db psql -U salmon -d salmon < backend/db/seed/approve_dev.sql
--
-- ไม่แตะ seed.sql เพราะ TestSeed_BorrowRequestScenarios ของ 002 บังคับว่า seed.sql ต้องไม่มีคำขอยืม
-- เพิ่มรุ่นและอุปกรณ์ชุดของตัวเอง (รหัสขึ้นต้น AP-TEST-) ตัวเลขที่ใช้ลองของ 001 / 002 / 005 จึงไม่เปลี่ยน
--
-- กรณีที่จุดตรวจของ 004 ใช้:
--   AP-TEST-0001 ว่าง: คำขอปกติ A (ส่งก่อน), B (ส่งทีหลัง) และคำขอค้างคืน C -> คำขอแรก / ปฏิเสธตาม / ข้ามลำดับไม่ได้ (AC-APV-03, 07, 17)
--   AP-TEST-0002 ว่าง: คำขอปกติใบเดียว -> อนุมัติสำเร็จ (AC-APV-01)
--   AP-TEST-0003 ว่าง: คำขอปกติใบเดียว -> ปฏิเสธพร้อมเหตุผล / ไม่ระบุเหตุผล (AC-APV-02, 06)
--   AP-TEST-0004 ถูกยืม: คำขอปกติเป็นคำขอแรก -> อนุมัติไม่ได้ ITEM_NOT_AVAILABLE (AC-APV-08)
--   AP-TEST-0005 ว่าง: คำขอที่ถูกยกเลิกแล้ว -> ไม่ขึ้นในรายการ (AC-APV-09, 11)
--   AP-TEST-0101 ว่าง ภาควิชา 2: คำขอปกติของภาควิชาอื่น -> ไม่ขึ้นในรายการ / FORBIDDEN_DEPARTMENT (AC-APV-11, 12)
--
-- ผู้ขอ: student-001, student-002 (ภาควิชา 1) และ student-101 (ภาควิชา 2) ของ auth.MockProvider
-- student-003 ไม่มีใน MockProvider ใช้เป็นข้อมูลผู้ขออีกคนเท่านั้น เพราะผู้ขอคนเดิมมีคำขอที่ยังไม่จบบนชิ้นเดิมได้ใบเดียว (uq_request_open_per_item ของ 002)
--
-- เวลา (Asia/Bangkok, glossary):
--   วันที่เริ่มยืม = วันนี้ถ้าโหลดก่อน 16:00 น. ไม่อย่างนั้นเป็นพรุ่งนี้ เพื่อให้อนุมัติได้ (FR-APV-01 ค เลย 16:00 น. ของวันที่เริ่มยืมอนุมัติไม่ได้)
--   เวลาที่ส่ง = ย้อนจากเวลาที่โหลดไม่กี่นาที เรียง A < B < C ถ้าโหลดช่วงหลังเที่ยงคืนไม่กี่นาที คำขออาจนับเป็นของเมื่อวานและถูกปฏิเสธอัตโนมัติ (FR-APV-01 ฉ)

INSERT INTO equipment_model (name, category_id, department_id) VALUES
    ('อุปกรณ์ทดลองอนุมัติ (004)', 1, 1),
    ('อุปกรณ์ทดลองอนุมัติ ภาควิชาอื่น (004)', 1, 2);

INSERT INTO equipment_item (asset_code, model_id, status)
SELECT v.code, m.id, v.st
FROM (VALUES
    ('AP-TEST-0001', 'available',   'อุปกรณ์ทดลองอนุมัติ (004)'),
    ('AP-TEST-0002', 'available',   'อุปกรณ์ทดลองอนุมัติ (004)'),
    ('AP-TEST-0003', 'available',   'อุปกรณ์ทดลองอนุมัติ (004)'),
    ('AP-TEST-0004', 'borrowed',    'อุปกรณ์ทดลองอนุมัติ (004)'),
    ('AP-TEST-0005', 'available',   'อุปกรณ์ทดลองอนุมัติ (004)'),
    ('AP-TEST-0101', 'available',   'อุปกรณ์ทดลองอนุมัติ ภาควิชาอื่น (004)')
) AS v(code, st, model)
JOIN equipment_model m ON m.name = v.model;

-- "ถูกยืม" = มี Loan ที่ยังไม่คืน (data-model.md) จึงสร้าง Loan ให้ AP-TEST-0004 ด้วย กำหนดคืน 16:00 น. วันนี้ (DOM-LOAN-01)
INSERT INTO loan (item_id, borrower_id, loan_type, due_at)
SELECT id, 'student-003', 'normal',
       (date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') + interval '16 hours') AT TIME ZONE 'Asia/Bangkok'
FROM equipment_item WHERE asset_code = 'AP-TEST-0004';

-- คำขอยืม (ตาราง borrow_request ของ 002)
WITH d AS (
    SELECT CASE WHEN extract(hour FROM now() AT TIME ZONE 'Asia/Bangkok') < 16
                THEN date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok')
                ELSE date_trunc('day', now() AT TIME ZONE 'Asia/Bangkok') + interval '1 day'
           END AS start_day
)
INSERT INTO borrow_request (requester_id, requester_name, student_code, item_id, loan_type,
                            purpose, purpose_other, start_at, end_at, status, submitted_at,
                            attachment_path, cancel_reason, cancelled_at)
SELECT v.requester, v.name, v.requester, i.id, v.loan_type,
       v.purpose, NULL,
       (d.start_day + interval '9 hours') AT TIME ZONE 'Asia/Bangkok',
       -- ยืมปกติคืน 16:00 น. วันเดียวกัน / ค้างคืนคืนวันถัดไป (FR-BRQ-10)
       (d.start_day + CASE WHEN v.loan_type = 'overnight' THEN interval '1 day 16 hours' ELSE interval '16 hours' END)
           AT TIME ZONE 'Asia/Bangkok',
       v.st,
       now() - v.ago,
       CASE WHEN v.loan_type = 'overnight' THEN 'dev/approve-overnight-sample.pdf' END,  -- path ตัวอย่าง ไม่มีไฟล์จริง
       CASE WHEN v.st = 'cancelled' THEN 'ข้อมูลตัวอย่าง: ผู้ขอยกเลิก' END,
       CASE WHEN v.st = 'cancelled' THEN now() - interval '1 minute' END
FROM d,
(VALUES
    -- code,          ผู้ขอ,          ชื่อ,                  ประเภท,      วัตถุประสงค์, สถานะ,      ส่งเมื่อ
    ('AP-TEST-0001', 'student-001', 'นักศึกษาตัวอย่าง ก', 'normal',    'study',    'requested', interval '30 minutes'),  -- A คำขอแรก
    ('AP-TEST-0001', 'student-002', 'นักศึกษาตัวอย่าง ข', 'normal',    'activity', 'requested', interval '20 minutes'),  -- B
    ('AP-TEST-0001', 'student-003', 'นักศึกษาตัวอย่าง ค', 'overnight', 'activity', 'requested', interval '10 minutes'),  -- C ค้างคืน
    ('AP-TEST-0002', 'student-001', 'นักศึกษาตัวอย่าง ก', 'normal',    'study',    'requested', interval '25 minutes'),
    ('AP-TEST-0003', 'student-002', 'นักศึกษาตัวอย่าง ข', 'normal',    'study',    'requested', interval '15 minutes'),
    ('AP-TEST-0004', 'student-001', 'นักศึกษาตัวอย่าง ก', 'normal',    'study',    'requested', interval '12 minutes'),
    ('AP-TEST-0005', 'student-002', 'นักศึกษาตัวอย่าง ข', 'normal',    'study',    'cancelled', interval '40 minutes'),
    ('AP-TEST-0101', 'student-101', 'นักศึกษาตัวอย่าง ภาควิชาอื่น', 'normal', 'study', 'requested', interval '35 minutes')
) AS v(code, requester, name, loan_type, purpose, st, ago)
JOIN equipment_item i ON i.asset_code = v.code;
