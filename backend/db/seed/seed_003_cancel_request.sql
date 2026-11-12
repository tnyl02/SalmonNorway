-- ข้อมูลตั้งต้นของฟีเจอร์ 003 ยกเลิกคำขอยืม (T-02) ใช้ระหว่างที่หน้าส่งคำขอของ 002 ยังไม่เสร็จ (ASM-01)
-- โหลดหลัง migration และ seed.sql ของ 001 แล้วเท่านั้น (ใช้อุปกรณ์และภาควิชาจาก seed.sql)
--   docker compose -f deploy/docker-compose.yml exec -T db psql -U salmon -d salmon < backend/db/seed/seed_003_cancel_request.sql
-- โหลดซ้ำได้ (ON CONFLICT DO NOTHING) ไม่แตะตาราง equipment_item (CON-CXL-01, CON-DATA-01)
--
-- ชื่อ รหัสนักศึกษา และ Email เป็นค่าสมมติที่ทีมเลือก (2569-11-11 ตอนทำ T-02) ไม่ได้มาจาก spec
-- วัตถุประสงค์ เวลาเริ่มและคืน เป็นค่าทดสอบตามรูปแบบตาราง 0003 ของ 002 ไม่ได้ใช้ในฟีเจอร์ 003
--
-- คำขอที่ยังยกเลิกได้ (3001-3004, 3008) ใช้เวลาเทียบกับ now() ไม่ใช้วันที่ตายตัว (แก้ 2569-11-12 หลัง merge 004)
-- เพราะงานเบื้องหลังของ 004 (approve.AutoRejecter, FR-APV-01 ฉ) ปฏิเสธคำขอแบบ normal ที่ยัง requested
-- และส่งก่อนเที่ยงคืนวันนี้ทันทีที่เปิดเซิร์ฟเวอร์ เวลาส่งห่างกันเป็นนาทีและเรียงเหมือนเดิม 3008 > 3001 > 3002 > 3003 > 3004
-- (FR-CXL-06 เรียงใหม่ไปเก่า) ถ้าโหลดช่วง 00:00-00:30 น. คำขอบางรายการอาจยังถูกปฏิเสธได้
-- คำขอที่จบแล้ว (3005-3007) คงวันที่ตายตัวไว้ เป็นประวัติ ไม่โดนงานเบื้องหลัง
--
-- ตาราง 0003 ห้ามผู้ขอคนเดิมมีคำขอที่ยังไม่จบบนอุปกรณ์ชิ้นเดียวกันเกิน 1 ใบ (uq_request_open_per_item)
-- คำขอที่ยังไม่จบของ student-001 จึงอยู่คนละชิ้น และใช้เฉพาะชิ้นที่ "ว่าง" (CON-DATA-01)

INSERT INTO borrow_request
    (id, requester_id, requester_name, student_code, item_id, loan_type, purpose,
     start_at, end_at, status, submitted_at, cancel_reason, cancelled_at)
VALUES
    -- student-001: คำขอที่ยังยกเลิกได้ 4 รายการ (CON-CXL-02) หน้า "คำขอที่รออนุมัติของฉัน" ต้องเห็นทั้ง 4 เรียงใหม่ไปเก่า
    (3001, 'student-001', 'ใจดี สวยมาก', '660710755', 1, 'normal', 'study',
        now() + interval '1 day', now() + interval '1 day 6 hours', 'requested', now() - interval '10 minutes', NULL, NULL), -- CS-NB-0001
    (3002, 'student-001', 'ใจดี สวยมาก', '660710755', 6, 'normal', 'activity',
        now() + interval '1 day', now() + interval '1 day 3 hours', 'requested', now() - interval '20 minutes', NULL, NULL), -- CS-PJ-0001
    (3003, 'student-001', 'ใจดี สวยมาก', '660710755', 2, 'overnight', 'study',
        now() + interval '2 days', now() + interval '3 days', 'approved_by_officer', now() - interval '25 minutes', NULL, NULL), -- CS-NB-0002 ค้างคืน ASM-04
    (3004, 'student-001', 'ใจดี สวยมาก', '660710755', 3, 'overnight', 'other',
        now() + interval '3 days', now() + interval '4 days', 'pending_advisor_approval', now() - interval '30 minutes', NULL, NULL), -- CS-NB-0003 ค้างคืน ASM-04
    -- student-001: คำขอที่ยกเลิกไม่ได้ หน้านี้ต้องไม่แสดง (FR-CXL-02)
    (3005, 'student-001', 'ใจดี สวยมาก', '660710755', 8, 'normal', 'study',
        '2026-11-12 09:00+07', '2026-11-12 16:00+07', 'approved', '2026-11-07 10:00+07', NULL, NULL),           -- XX-NB-0001
    (3006, 'student-001', 'ใจดี สวยมาก', '660710755', 1, 'normal', 'study',
        '2026-11-06 09:00+07', '2026-11-06 16:00+07', 'rejected', '2026-11-05 10:00+07', NULL, NULL),           -- CS-NB-0001
    (3007, 'student-001', 'ใจดี สวยมาก', '660710755', 1, 'normal', 'study',
        '2026-11-05 09:00+07', '2026-11-05 16:00+07', 'cancelled', '2026-11-04 10:00+07',
        'ไม่ต้องใช้อุปกรณ์แล้ว', '2026-11-04 12:00+07'),                                                        -- CS-NB-0001
    -- student-002: คำขอบนอุปกรณ์ชิ้นเดียวกับ 3001 ใช้ตรวจว่ายกเลิกของ student-001 ไม่กระทบคำขอของผู้อื่น (DOM-REQ-01, FR-CXL-03)
    (3008, 'student-002', 'มีสุข ใจเย็น', '660710756', 1, 'normal', 'study',
        now() + interval '1 day', now() + interval '1 day 6 hours', 'requested', now() - interval '5 minutes', NULL, NULL)  -- CS-NB-0001
    -- student-101 (ภาควิชา 2) ไม่มีคำขอ ใช้ลองกรณี "ไม่มีคำขอที่รออนุมัติ" (FR-CXL-06, AC-CXL-12)
ON CONFLICT (id) DO NOTHING;

-- Email เจ้าหน้าที่จำลองของภาควิชา 1 สำหรับ Notifier จำลอง (ASM-10, FR-CXL-09)
INSERT INTO staff_contact (department_id, email)
SELECT 1, 'staff-it@example.test'
WHERE NOT EXISTS (SELECT 1 FROM staff_contact WHERE department_id = 1 AND email = 'staff-it@example.test');

-- ให้ id ที่สร้างอัตโนมัติ (คำขอจากหน้าของ 002) ต่อจากข้อมูลตั้งต้น
SELECT setval('borrow_request_id_seq', GREATEST((SELECT max(id) FROM borrow_request), 1));
