-- 0006 ตารางการส่งมอบของฟีเจอร์ 005 (UC-08, specs/005-handover/plan.md ข้อ 3)
-- เลขไฟล์ทีม 005 กำหนด 2569-11-11 (data-model.md ตาราง "เลข migration")
-- ไม่ ALTER ตาราง loan ของ 001 (plan v1.1) ข้อมูลการส่งมอบอยู่ในตาราง handover ที่ผูก 1:1 กับ loan
-- ไม่มีคอลัมน์เลขบัตรประชาชนหรือเลขบัตรนักศึกษา FR-HND-03 บันทึกแค่ชื่อ นามสกุล รหัสนักศึกษา เวลา และผลการตรวจบัตร

-- การส่งมอบ 1 ครั้ง = Loan 1 รายการ = คำขอยืม 1 รายการ (FR-HND-02, 03, 05, 07)
CREATE TABLE handover (
    loan_id             INT          PRIMARY KEY REFERENCES loan (id),       -- FR-HND-07: การยืมที่สร้างตอนส่งมอบ
    request_id          INT          NOT NULL UNIQUE
                                     REFERENCES borrow_request (id),         -- FR-HND-01, 07: คำขอที่ส่งมอบ ส่งมอบได้ครั้งเดียว
    handed_over_at      TIMESTAMPTZ  NOT NULL,                               -- FR-HND-03, 05: วันที่และเวลาที่ส่งมอบจริง
    handed_over_by      TEXT         NOT NULL,                               -- FR-HND-08: รหัสเจ้าหน้าที่ผู้ส่งมอบ (CON-HND-01)
    borrower_student_id TEXT         NOT NULL,                               -- FR-HND-03: รหัสนักศึกษาของผู้ยืม
    borrower_first_name TEXT         NOT NULL,                               -- FR-HND-03: ชื่อ (ที่มารอ Q-HND-13)
    borrower_last_name  TEXT         NOT NULL,                               -- FR-HND-03: นามสกุล (ที่มารอ Q-HND-13)
    identity_method     TEXT         NOT NULL
        CHECK (identity_method IN ('auth_api', 'card_only')),                -- FR-HND-02: Auth API สำเร็จ / ใช้ผลตรวจบัตรแทน
    card_checked        BOOLEAN      NOT NULL CHECK (card_checked)           -- FR-HND-03: ต้องตรวจบัตรทุกครั้ง (ASM-HND-01)
);

-- รูปสภาพอุปกรณ์ตอนส่งมอบ ไฟล์อยู่บน Docker volume เก็บแค่ path (FR-HND-06, plan.md Q-P1)
-- จำนวนขั้นต่ำและขนาดสูงสุดตรวจที่ service จากค่าใน config.go (AGENTS.md กฎข้อ 9) ไม่ฝังตัวเลขที่นี่
CREATE TABLE handover_photo (
    id           SERIAL       PRIMARY KEY,
    loan_id      INT          NOT NULL REFERENCES handover (loan_id),        -- FR-HND-06: รูปผูกกับการยืมนั้น
    file_path    TEXT         NOT NULL,
    content_type TEXT         NOT NULL
        CHECK (content_type IN ('image/jpeg', 'image/png')),                 -- FR-HND-06: JPG หรือ PNG (ASM-HND-07)
    size_bytes   BIGINT       NOT NULL CHECK (size_bytes > 0),
    uploaded_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_handover_photo_loan ON handover_photo (loan_id);
