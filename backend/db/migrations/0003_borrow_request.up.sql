-- 0003 ตาราง borrow_request ของฟีเจอร์ 002 (DM-06 ใน specs/000-shared/data-model.md)
-- เลขไฟล์จองไว้ในตาราง "เลข migration" ของ data-model.md (ASM-BRQ-03)
-- ฟีเจอร์นี้ไม่เปลี่ยนสถานะอุปกรณ์ (DOM-REQ-01) จึงไม่มีคำสั่งใดแตะตาราง equipment_item

CREATE TABLE borrow_request (
    id              SERIAL       PRIMARY KEY,                   -- ASM-BRQ-03: INT ให้ชนิดตรงกับ audit_log.request_id
    requester_id    TEXT         NOT NULL,                      -- FR-BRQ-04: รหัสผู้ใช้จาก IF-IDP-01 เท่านั้น (ASM-BRQ-06)
    requester_name  TEXT         NOT NULL,                      -- FR-BRQ-03 (ก): ผู้ขอพิมพ์เอง เป็นข้อมูลประกอบ
    student_code    TEXT         NOT NULL,                      -- FR-BRQ-03 (ก): ผู้ขอพิมพ์เอง ไม่ใช้ตรวจสิทธิ์
    item_id         INT          NOT NULL REFERENCES equipment_item (id),  -- FR-BRQ-01: ขอยืมรายชิ้น (ASM-BRQ-02)
    loan_type       TEXT         NOT NULL
        CHECK (loan_type IN ('normal', 'overnight')),           -- FR-BRQ-10, DM-05 (ปกติ / ค้างคืน)
    purpose         TEXT         NOT NULL
        CHECK (purpose IN ('study', 'activity', 'other')),      -- FR-BRQ-02: 3 ค่าตาม data-model.md v9
    purpose_other   TEXT         NULL,                          -- FR-BRQ-02 (ก): บังคับเมื่อ purpose='other' ตรวจที่ service (T-07)
    start_at        TIMESTAMPTZ  NOT NULL,                      -- FR-BRQ-03 (ข): วันและเวลาที่เริ่มยืม
    end_at          TIMESTAMPTZ  NOT NULL,                      -- FR-BRQ-03 (ข)(ค): วันและเวลาที่คืน กฎ 16:00 ตรวจที่ service (T-08)
    status          TEXT         NOT NULL
        CHECK (status IN ('requested', 'approved_by_officer', 'pending_advisor_approval',
                          'approved', 'rejected', 'handed_over', 'returned', 'cancelled')),  -- data-model.md v9
    submitted_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),        -- FR-BRQ-04, FR-BRQ-05: ใช้เรียงลำดับคำขอ (DOM-REQ-01)
    attachment_path TEXT         NULL,                          -- FR-BRQ-10 (ก): path ของใบบันทึกข้อความ ไฟล์อยู่นอกฐานข้อมูล (ASM-BRQ-07)
    cancel_reason   VARCHAR(200) NULL,                          -- ของ 003 (UC-15) 002 สร้างให้แต่ไม่เขียนค่า
    cancelled_at    TIMESTAMPTZ  NULL                           -- ของ 003 (UC-15)
);

-- ไม่มีคอลัมน์กำหนดคืน เพราะ 005 เป็นผู้บันทึกตอนส่งมอบ (FR-BRQ-10 ค, DM-05)
-- ไม่มีคอลัมน์อีเมล เพราะยังไม่มีที่มาในระบบ (DM-08 มีแค่รหัสผู้ใช้ บทบาท ภาควิชา)

-- FR-BRQ-05 (ข) ตามขอบเขตที่ CR-004 ขยาย: ผู้ขอคนเดิมมีคำขอที่ยังไม่จบของอุปกรณ์ชิ้นเดิมได้ใบเดียว
-- บังคับที่ฐานข้อมูลด้วย เพื่อกันการกดส่งซ้ำพร้อมกัน (FR-BRQ-09)
CREATE UNIQUE INDEX uq_request_open_per_item
    ON borrow_request (requester_id, item_id)
    WHERE status IN ('requested', 'approved_by_officer', 'pending_advisor_approval', 'approved');

-- หน้า "คำขอของฉัน" ของ 003 (ขอโดยเจ้าของ 003 2569-11-11)
CREATE INDEX idx_request_owner_status
    ON borrow_request (requester_id, status, submitted_at DESC);

-- FR-BRQ-05 (ก): ลำดับคำขอของอุปกรณ์ชิ้นนั้น ให้ 001 เรียกไปแสดงแก่เจ้าหน้าที่
CREATE INDEX idx_request_item_queue
    ON borrow_request (item_id, submitted_at);

-- ASM-BRQ-03: 0002 ของ 001 เขียนคอมเมนต์ฝากให้ 002 เพิ่ม foreign key เมื่อมีตารางนี้
ALTER TABLE audit_log
    ADD CONSTRAINT fk_audit_log_request
    FOREIGN KEY (request_id) REFERENCES borrow_request (id);

-- Q-BRQ-07: เพิ่มค่า 'request' ของ 002 ในรายการ action ของ DM-07
-- ใช้ชื่อ constraint เดิมที่ PostgreSQL ตั้งให้ใน 0002 เพื่อให้ฟีเจอร์ถัดไปเพิ่มค่าด้วยวิธีเดียวกัน (ตกลงกับ 003 2569-11-11)
-- ห้ามแก้ไฟล์ 0002 ตามกติกาที่เขียนไว้ในไฟล์นั้น
ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;
ALTER TABLE audit_log ADD CONSTRAINT audit_log_action_check
    CHECK (action IN ('list', 'search', 'detail', 'cancel', 'request'));
