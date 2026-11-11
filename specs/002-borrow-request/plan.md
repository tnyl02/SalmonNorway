# Implementation Plan: ส่งคำขอยืมอุปกรณ์ (Borrow Request)
Spec: SPEC-BRQ-002 (specs/002-borrow-request/spec.md, Draft v3.1) | Plan status: Draft v1.1 | Updated: 2569-11-11 (CR-004)
อ้างอิง: specs/changes/{CR-003.md, CR-004.md}, specs/000-shared/{constraints.md v5, data-model.md v9, glossary.md v2}

## 1. สรุปแนวทาง
ฟีเจอร์นี้ให้นักศึกษาที่ยืนยันตัวตนแล้ว ส่งคำขอยืมอุปกรณ์รายชิ้นที่สถานะ "ว่าง" จากหน้าอุปกรณ์ของ 001
โดยระบบตรวจความพร้อมของอุปกรณ์ สิทธิ์ภาควิชา และประวัติการคืนเกินกำหนดให้ก่อนสร้างคำขอ
แนวทางหลัก: ฟีเจอร์นี้เป็น **write-only ต่อตาราง `borrow_request`** และ **ไม่แตะสถานะอุปกรณ์เลย** (DOM-REQ-01)
กฎทุกข้อบังคับที่ฝั่ง server ไม่พึ่งการซ่อนปุ่มที่หน้าจอ และกฎที่ซ้ำซ้อนได้ (คำขอซ้ำของคนเดิม)
บังคับทั้งที่ชั้นโปรแกรมและที่ฐานข้อมูลด้วย unique index เพื่อกัน race condition ตอนกดส่งซ้ำ

## 2. เทคโนโลยีที่ใช้
| สิ่งที่เลือก | มาจาก | หมายเหตุ |
|---|---|---|
| PostgreSQL 16 | CON-TECH-01 | ตามที่ spec สั่ง |
| HTTPS เท่านั้น + Docker Compose | CON-TECH-02, NFR-SEC-01 | TLS ที่ reverse proxy เป็นงานระดับระบบ ยังไม่มีเจ้าของ |
| Go (net/http + pgx) หลังบ้าน | **ทีมเลือกเอง ไม่ได้มาจาก spec** | ตาม AGENTS.md และ plan ข้อ 2 ของ 001 ใช้ของเดิมทั้งหมด |
| React + Vite หน้าบ้าน | NFR-USE-01 (มือถือ) + **ทีมเลือกเอง** | ใช้โครงเดียวกับ 001 |
| golang-migrate (SQL ล้วน) | **ทีมเลือกเอง** | ไฟล์ของฟีเจอร์นี้คือ `0003_borrow_request.up.sql` (data-model.md v9) |
| go test + testcontainers | **ทีมเลือกเอง** | ใช้ชุดเดิมที่ 001 วางไว้ใน T-01, T-02 |
| auth จำลอง `Bearer mock-<userID>` | IF-IDP-01 (โปรโตคอลจริงรอ Q-SH-03) | ใช้ middleware ของ 001 (`backend/internal/auth`) ไม่เขียนใหม่ |
| เก็บไฟล์แนบเป็นไฟล์บนเครื่องที่รันระบบ | ASM-BRQ-07 (**ทีมตัดสิน ไม่มี Constraint รองรับ**) | ไม่เก็บลงฐานข้อมูล ตำแหน่งโฟลเดอร์อ่านจาก config |

ขนาดงาน: คำขอหลักร้อยรายการต่อเทอม ผู้ใช้พร้อมกัน 50 คน — ไม่ต้องมี cache, queue หรือ replica

## 3. โมเดลข้อมูล
ตารางใหม่ของฟีเจอร์นี้คือ `borrow_request` (DM-06) สร้างด้วย migration `0003`
ส่วน `equipment_item`, `loan` (อ่านอย่างเดียว) และ `audit_log` มีอยู่แล้วจาก 001

```
borrow_request(
  id             SERIAL PK,                       -- ASM-BRQ-03 ต้องเป็น INT ให้ตรงกับ audit_log.request_id
  requester_id   TEXT NOT NULL,                   -- FR-BRQ-04, ASM-BRQ-06 (รหัสผู้ใช้จาก IF-IDP-01 เท่านั้น)
  requester_name TEXT NOT NULL,                   -- FR-BRQ-03 (ก) ผู้ขอพิมพ์เอง เป็นข้อมูลประกอบ
  student_code   TEXT NOT NULL,                   -- FR-BRQ-03 (ก) ผู้ขอพิมพ์เอง ไม่ใช้ตรวจสิทธิ์
  item_id        INT NOT NULL REFERENCES equipment_item(id),  -- FR-BRQ-01, ASM-BRQ-02 (รายชิ้น)
  loan_type      TEXT NOT NULL CHECK (loan_type IN ('normal','overnight')),  -- FR-BRQ-10, DM-05
  purpose        TEXT NOT NULL CHECK (purpose IN ('study','activity','other')),  -- FR-BRQ-02 (3 ค่า)
  purpose_other  TEXT NULL,                       -- FR-BRQ-02 (ก) บังคับเมื่อ purpose='other'
  start_at       TIMESTAMPTZ NOT NULL,            -- FR-BRQ-03 (ข)
  end_at         TIMESTAMPTZ NOT NULL,            -- FR-BRQ-03 (ข)(ค) ยืมปกติห้ามเกิน 16:00 ของวันเดียวกัน
  status         TEXT NOT NULL CHECK (status IN
                 ('requested','approved_by_officer','pending_advisor_approval',
                  'approved','rejected','handed_over','returned','cancelled')),  -- data-model.md v9
  submitted_at   TIMESTAMPTZ NOT NULL DEFAULT now(),  -- FR-BRQ-04, FR-BRQ-05 (เรียงลำดับ DOM-REQ-01)
  attachment_path TEXT NULL,                      -- FR-BRQ-10 (ก) บังคับเมื่อ loan_type='overnight'
  cancel_reason  VARCHAR(200) NULL,               -- ของ 003 (002 สร้างให้ ไม่เขียนค่า)
  cancelled_at   TIMESTAMPTZ NULL                 -- ของ 003
)
```

Index และกฎที่บังคับในฐานข้อมูล
- `UNIQUE (requester_id, item_id) WHERE status IN ('requested','approved_by_officer','pending_advisor_approval','approved')`
  — บังคับ FR-BRQ-05 (ข) ตามขอบเขตที่ CR-004 ขยาย และกันคำขอซ้ำจากการกดส่งซ้ำ (FR-BRQ-09)
  เป็น partial index ครอบเฉพาะคำขอที่ยังไม่จบ เมื่อคำขอกลายเป็น `rejected`, `cancelled`, `handed_over`
  หรือ `returned` แถวจะหลุดจากเงื่อนไข ผู้ขอจึงส่งคำขอชิ้นเดิมใหม่ได้
  ผลต่อฟีเจอร์อื่น: 003-006 เพียง UPDATE สถานะของแถวเดิม ไม่เคย INSERT แถวสถานะเหล่านี้
  จึงไม่มีทางชน index นี้ (ตรวจแล้ว 2569-11-11)
- `INDEX (requester_id, status, submitted_at DESC)` — ตามคำขอของ 003 สำหรับหน้า "คำขอของฉัน"
- `INDEX (item_id, submitted_at)` — ใช้คำนวณลำดับคำขอของอุปกรณ์ชิ้นนั้น (FR-BRQ-05 ก)
- `ALTER TABLE audit_log ADD CONSTRAINT fk_audit_log_request FOREIGN KEY (request_id) REFERENCES borrow_request(id)` (ASM-BRQ-03)
- `ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;` แล้วสร้างใหม่ด้วยชื่อเดิม โดยเพิ่มค่า `request` (Q-BRQ-07)

สิ่งที่ **ไม่มี** ในตารางนี้ โดยตั้งใจ
- **ไม่มีคอลัมน์กำหนดคืน** — FR-BRQ-10 (ค) ระบุว่าเป็นงานของ 005 ตอนส่งมอบ (DM-05)
- **ไม่มีอีเมลผู้ขอ** — ไม่มีที่มาในระบบ (DM-08 มีแค่รหัสผู้ใช้ บทบาท ภาควิชา) 004 ถามมาแล้วและยังไม่มีคำตอบ
- **ไม่มีคอลัมน์สถานะอุปกรณ์** — ฟีเจอร์นี้ไม่เปลี่ยนสถานะอุปกรณ์ (DOM-REQ-01)

ไฟล์แนบ: เก็บไฟล์ไว้ที่ `<ATTACHMENT_DIR>/<requestID>-<ชื่อไฟล์ที่ทำให้ปลอดภัยแล้ว>` และเก็บเฉพาะ path ในฐานข้อมูล (ASM-BRQ-07)

## 4. API และหน้าจอ

| Endpoint | รองรับ FR | input / output หลัก |
|---|---|---|
| `POST /api/v1/loan-requests` (multipart เมื่อมีไฟล์แนบ) | FR-BRQ-01, 02, 03, 04, 05ข, 06, 07, 09, 10 | in: assetCode, loanType, purpose, purposeOther?, requesterName, studentCode, startAt, endAt, attachment? / out: `{requestId, status:"requested", submittedAt}` |
| `GET /api/v1/equipment/items/{assetCode}/requests` | FR-BRQ-05 (ก) | เฉพาะ staff / out: `[{requestId, queuePosition, submittedAt, requesterName}]` ให้ 001 เรียกไปแสดง |
| `GET /api/v1/loan-requests/{id}/attachment` | FR-BRQ-10 (ก) | เฉพาะเจ้าหน้าที่และเจ้าของคำขอ / out: ไฟล์ |

ตอนเปิดฟอร์ม หน้าจอเรียก `GET /api/v1/equipment/items/{assetCode}` ของ 001 เพื่อตรวจสถานะรอบแรก (FR-BRQ-01)
แล้ว `POST` จะตรวจสถานะซ้ำอีกครั้งในทรานแซกชันเดียวกับการสร้างคำขอ

Error contract — รูปแบบและรหัสมาจาก **CON-TECH-03** (ของกลาง เพิ่มโดย CR-004) ไม่ได้นิยามเองในฟีเจอร์นี้
| กรณี | HTTP | code | อ้างอิง |
|---|---|---|---|
| ไม่มี/หมดอายุ session | 401 | `UNAUTHENTICATED` | IF-IDP-01 |
| ไม่มีภาควิชาหรือบทบาทไม่รู้จัก | 403 | `PRINCIPAL_INCOMPLETE` | IF-IDP-01 (ใช้ของ 001 T-03) |
| อุปกรณ์ต่างภาควิชา | 403 | `FORBIDDEN_DEPARTMENT` (ไม่ส่งข้อมูลอุปกรณ์กลับ) | FR-BRQ-07, AC-BRQ-07 |
| สถานะไม่ใช่ "ว่าง" แล้ว | 409 | `ITEM_NOT_AVAILABLE` + สถานะล่าสุด | FR-BRQ-01, AC-BRQ-01 |
| มีคำขอที่ยังไม่จบของตนเองบนชิ้นเดิม | 409 | `DUPLICATE_REQUEST` | FR-BRQ-05ข, AC-BRQ-11, CR-004 |
| มีของคืนเกินกำหนดค้างอยู่ | 409 | `HAS_OVERDUE_LOAN` + เหตุผล | FR-BRQ-06, AC-BRQ-06, AC-BRQ-12 |
| ข้อมูลไม่ครบ / เลือก "อื่น ๆ" แต่ไม่พิมพ์ / เวลาคืนเกิน 16:00 / ไฟล์ผิดชนิดหรือใหญ่เกิน | 422 | `VALIDATION_FAILED` + รายชื่อช่องที่ผิด | FR-BRQ-02, 03, 10, AC-BRQ-02, 03, 10, 13, 14 |
| บันทึกไม่สำเร็จ | 503 | `UPSTREAM_UNAVAILABLE` (ไม่สร้างคำขอค้าง) | FR-BRQ-09, AC-BRQ-09 |

หน้าจอ
- `BorrowRequestForm` — **ไม่มี mockup ทีมออกแบบเอง** และ spec ยังไม่มีหมวดหน้าจอ (Q-BRQ-14)
  ข้อความบนหน้าจอทุกจุด (แจ้งสถานะล่าสุด, ช่องที่ขาด, กฎ 16:00 น., ชนิดและขนาดไฟล์) **ยังไม่สร้างจนกว่าจะได้คำตอบ Q-BRQ-14**
- ปุ่มเข้าหน้านี้มาจาก 001 (T-34) ซึ่งรอ CR ของ 001 อยู่ ฟีเจอร์นี้จึงเปิดหน้าตรงด้วย URL ได้โดยไม่ต้องรอ

## 5. ตารางตรวจ Constraints

| Constraint ID | ถูกนำไปใช้ที่ไหนใน plan | สถานะ |
|---|---|---|
| CON-TECH-01 | ข้อ 2 (PostgreSQL 16), ข้อ 3 ทั้งหมด | ใช้แล้ว |
| CON-TECH-02 | ข้อ 2 (HTTPS + Docker) | ใช้แล้ว (TLS ที่ proxy ยังเป็นงานระดับระบบ ไม่มีเจ้าของ เหมือน 001) |
| CON-TECH-03 | ข้อ 4 ตาราง Error contract ใช้รหัสจากรายการกลางทั้ง 8 ตัว | ใช้แล้ว |
| CON-DATA-01 | ข้อ 4 `POST` ตรวจ `status='available'` ก่อนสร้างคำขอ และไม่เขียนสถานะกลับ | ใช้แล้ว |
| CON-DATA-03 | ข้อ 4 รับ `assetCode` เป็นตัวระบุอุปกรณ์ แล้วแปลงเป็น `item_id` | ใช้แล้ว |
| DOM-PRIV-01 | ข้อ 4 `GET .../requests` เปิดเฉพาะ staff, response ของนักศึกษาไม่มีข้อมูลคำขอผู้อื่น | ใช้แล้ว |
| DOM-AUDIT-01 | ข้อ 7 ขั้นที่ 8 — เขียน audit log ตอนอ่านข้อมูลอุปกรณ์ในขั้นตอนส่งคำขอ | ใช้แล้ว |
| DOM-ACCESS-01 | ข้อ 4 `FORBIDDEN_DEPARTMENT` และ repository เติม `WHERE department_id = principal.departmentID` | ใช้แล้ว |
| DOM-LOAN-01 | ข้อ 3 (`end_at`), ข้อ 4 (422 เมื่อยืมปกติเกิน 16:00), config `LoanNormalReturnHour` | ใช้แล้ว |
| DOM-REQ-01 | ข้อ 3 (ไม่มีคอลัมน์สถานะอุปกรณ์), ข้อ 4 (`GET .../requests` ให้ลำดับ), unique index เฉพาะคำขอของตนเอง | ใช้แล้ว |
| IF-IDP-01 | ข้อ 2 (auth middleware ของ 001), ข้อ 3 (`requester_id`), ข้อ 4 (401/403) | ใช้แล้ว |
| IF-NOT-01 | - | **ยังไม่ได้ใช้ เพราะ** การแจ้งผลอนุมัติเป็นของ 004 (RQ-21) ฟีเจอร์นี้ไม่ส่งการแจ้งเตือน และ spec ไม่ได้อ้าง ID นี้ |

### ค่าที่ตั้งได้ (ห้ามฝังตัวเลขกระจายในโค้ด — AGENTS.md กฎข้อ 9)

| ค่า | มาจาก ID | ชื่อตัวแปรและไฟล์ |
|---|---|---|
| เวลาคืนของการยืมปกติ = 16:00 น. | DOM-LOAN-01, FR-BRQ-03 (ค) | `LoanNormalReturnHour = 16` — `backend/internal/config/config.go` |
| เขตเวลาที่ใช้เทียบทุกกรณี = Asia/Bangkok | glossary, FR-BRQ-04, FR-BRQ-06 (ข) | `AppTimeZone = "Asia/Bangkok"` — `config.go` (มีอยู่แล้วจาก 001 ใช้ซ้ำ) |
| จุดปลดล็อกคนคืนช้า = 00:00 น. ของวันถัดไป | FR-BRQ-06 (ข) | `OverdueBlockUntilNextDayStart = true` — `config.go` |
| ขนาดไฟล์แนบสูงสุด 10MB | FR-BRQ-10 (ก) | `AttachmentMaxBytes = 10 << 20` — `config.go` + `frontend/src/config.js` |
| ชนิดไฟล์แนบที่รับ (PDF, รูปถ่าย) | FR-BRQ-10 (ก) | `AttachmentAllowedTypes = ["application/pdf","image/jpeg","image/png"]` — `config.go` + `config.js` (ทีมยืนยันชนิดรูป 2569-11-11 spec เขียนแค่ "รูปถ่าย") |
| โฟลเดอร์เก็บไฟล์แนบ | ASM-BRQ-07 | `AttachmentDir` — `config.go` (อ่านจาก env ได้) |
| ค่าวัตถุประสงค์ 3 ค่า | FR-BRQ-02, data-model.md v9 | `PurposeValues = ["study","activity","other"]` — `config.go` / ข้อความไทยที่แสดงบนหน้าจอใน `config.js` |
| ค่าสถานะของคำขอ 8 ค่า | data-model.md v9, FR-BRQ-04 | `RequestStatuses` และ `StatusRequested = "requested"` — `config.go` |
| ค่า action ของ audit log = `request` | Q-BRQ-07, DM-07 | `AuditActionRequest = "request"` — `config.go` |
| เลข migration ของฟีเจอร์นี้ = 0003 | ASM-BRQ-03 | ชื่อไฟล์ `backend/db/migrations/0003_borrow_request.up.sql` |

## 6. แผนทดสอบจาก Acceptance Criteria

| AC ID | ชื่อ test | ทดสอบอย่างไร |
|---|---|---|
| AC-BRQ-01 | `TestAC_BRQ_01` | สร้าง item สถานะ available แล้วเปลี่ยนเป็น borrowed ก่อนเรียก POST — ต้องได้ 409 `ITEM_NOT_AVAILABLE` และตาราง `borrow_request` ว่าง |
| AC-BRQ-02 | `TestAC_BRQ_02` | POST ด้วย purpose=`study` ต้องสำเร็จและค่าในแถวตรงกัน / POST โดยไม่ส่ง purpose ต้องได้ 422 |
| AC-BRQ-03 | `TestAC_BRQ_03` | POST โดยไม่ส่ง `endAt` ต้องได้ 422 พร้อมรายชื่อช่องที่ขาด ครบ 5 ช่องตามที่ spec ระบุ |
| AC-BRQ-04 | `TestAC_BRQ_04` | POST ครบถ้วน ต้องได้ 201 สถานะ `requested`, `submitted_at` เป็นเวลา Asia/Bangkok และ `equipment_item.status` ยังเป็น `available` |
| AC-BRQ-05 | `TestAC_BRQ_05` | ผู้ใช้ A ส่งคำขอก่อน ผู้ใช้ B ส่งตาม แล้วเรียก `GET .../requests` ด้วย staff ต้องได้ลำดับ 1 และ 2 เรียงตาม `submitted_at` |
| AC-BRQ-06 | `TestAC_BRQ_06` | seed loan ที่ `due_at` เป็นเมื่อวาน 16:00 และ `returned_at` ว่าง แล้ว POST ต้องได้ 409 `HAS_OVERDUE_LOAN` |
| AC-BRQ-07 | `TestAC_BRQ_07` | ผู้ใช้ภาควิชา 2 ส่งคำขออุปกรณ์ภาควิชา 1 ต้องได้ 403 และ body ต้องไม่มีข้อมูลอุปกรณ์ |
| AC-BRQ-08 | `TestAC_BRQ_08` | ส่งคำขอสำเร็จ 1 ครั้งและถูกปฏิเสธ 1 ครั้ง แล้วนับแถวใน `audit_log` ต้องได้ 2 แถว มี `action='request'` และ `outcome` ต่างกัน |
| AC-BRQ-09 | `TestAC_BRQ_09` | ใช้ repository ปลอมที่ล้มเหลวครั้งแรก — ครั้งแรกต้อง 503 และไม่มีแถวค้าง ครั้งที่สองสำเร็จและมีแถวเดียว |
| AC-BRQ-10 | `TestAC_BRQ_10` | loanType=overnight โดยไม่แนบไฟล์ ต้อง 422 / แนบไฟล์ 11MB หรือไฟล์ .exe ต้อง 422 / แนบ PDF 1MB ต้องสำเร็จ |
| AC-BRQ-11 | `TestAC_BRQ_11` | ผู้ใช้เดิมส่งคำขอชิ้นเดิมซ้ำขณะสถานะ `requested` ต้องได้ 409 `DUPLICATE_REQUEST` และมีแถวเดียว |
| AC-BRQ-12 | `TestAC_BRQ_12` | ใช้นาฬิกาที่ตั้งได้: คืนของ 14:00 แล้วส่งคำขอ 23:59 ต้อง 409 / ส่ง 00:01 วันถัดไปต้องสำเร็จ |
| AC-BRQ-13 | `TestAC_BRQ_13` | loanType=normal, startAt วันนี้ 09:00, endAt วันนี้ 18:00 ต้องได้ 422 พร้อมข้อความกฎ 16:00 น. |
| AC-BRQ-14 | `TestAC_BRQ_14` | purpose=`other` โดยไม่ส่ง `purposeOther` ต้อง 422 / ส่งข้อความมาด้วยต้องสำเร็จและเก็บทั้งสองค่า |

หมายเหตุ: AC-BRQ-12 ต้องฉีดนาฬิกาเข้า service (แบบเดียวกับ T-12 ของ 001) ห้ามใช้ `time.Now()` ตรง ๆ ในชั้น service

## 7. ลำดับงาน

| ขั้น | งาน | เกี่ยวกับ |
|---|---|---|
| 1 | migration `0003_borrow_request.up.sql` — ตาราง, CHECK 8 ค่า, index 3 ตัว, FK ของ audit_log, เพิ่มค่า `request` ใน CHECK ของ action | ASM-BRQ-03, Q-BRQ-07, FR-BRQ-04 (**003 และ 004 รออยู่ ทำก่อนทุกขั้น**) |
| 2 | เพิ่มค่าที่ตั้งได้ทั้งหมดใน `config.go` และ `config.js` | ตารางข้อ 5 |
| 3 | repository + service สร้างคำขอ และทรานแซกชันตรวจสถานะซ้ำ | FR-BRQ-01, FR-BRQ-04, AC-BRQ-01, AC-BRQ-04 |
| 4 | กฎสิทธิ์ภาควิชาและคำขอซ้ำของตนเอง | FR-BRQ-05ข, FR-BRQ-07, AC-BRQ-07, AC-BRQ-11 |
| 5 | กฎคืนเกินกำหนด พร้อมนาฬิกาที่ฉีดได้ | FR-BRQ-06, AC-BRQ-06, AC-BRQ-12 |
| 6 | ตรวจความครบถ้วนของฟอร์ม รวมกฎ 16:00 น. และกฎ "อื่น ๆ ต้องพิมพ์ระบุ" | FR-BRQ-02, FR-BRQ-03, AC-BRQ-02, AC-BRQ-03, AC-BRQ-13, AC-BRQ-14 |
| 7 | ไฟล์แนบของการยืมค้างคืน (อัปโหลด ตรวจชนิดและขนาด และ endpoint เปิดดู) | FR-BRQ-10, AC-BRQ-10 |
| 8 | audit log ของการส่งคำขอทุกผลลัพธ์ | FR-BRQ-08, AC-BRQ-08, DOM-AUDIT-01 |
| 9 | endpoint ลำดับคำขอให้ 001 เรียกใช้ | FR-BRQ-05ก, AC-BRQ-05 |
| 10 | หน้าจอ `BorrowRequestForm` | FR-BRQ-01 ถึง 10 (**รอ Q-BRQ-14 ก่อน**) |

ขั้นที่ 1 ส่งขึ้น main ทันทีที่เสร็จ ไม่ต้องรอขั้นอื่น เพราะ 003 ติดอยู่ 17 task และ 004 ก็รอตารางเดียวกัน

## 8. สิ่งที่ยังไม่ทำ
คัดลอกจาก Open Questions ของ spec.md ทุกข้อที่ยังเปิดอยู่ **ส่วนที่เกี่ยวข้องกับแต่ละข้อจะยังไม่สร้างจนกว่าจะได้คำตอบ**

- **Q-BRQ-14** spec ยังไม่มีหมวด "หน้าจอ (UI)" และไม่มี mockup แต่ FR-BRQ-01, 02, 03, 10 สั่งให้แสดงข้อความหลายจุด
  ข้อความและหน้าตาเอามาจากไหน -> ทีม
  *ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ขั้นที่ 10 ทั้งขั้น)*
- **Q-BRQ-15** FR-BRQ-07 บอกว่าคำขอข้ามภาควิชาต้องไม่ได้รับข้อมูลอุปกรณ์ แต่ FR-BRQ-01 บอกว่าถ้าสถานะไม่ใช่ "ว่าง"
  ต้องแจ้งสถานะล่าสุด กรณีซ้อนกันจะแจ้งสถานะหรือปฏิเสธเงียบ -> ทีม / นักคอมพิวเตอร์
  *ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ลำดับการตรวจใน `POST` ระหว่าง 403 กับ 409)*
- **Q-SH-06 (ส่วนที่เหลือ)** ต้องแจ้งผู้ขอที่ถูกปฏิเสธอัตโนมัติตาม RQ-21 ไหม และเจ้าหน้าที่ต้องอนุมัติตามลำดับหรือข้ามได้
  -> ทีม / เจ้าของ 004
  *ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ข้อความอธิบายผลของคำขอที่ส่งซ้อนบนหน้าจอ)*
- **Q-SH-08** ถ้าระบบยืนยันตัวตนของมหาวิทยาลัยเชื่อมต่อไม่ได้ ระบบต้องทำอย่างไร -> นักเทคโนโลยีสารสนเทศ
  *ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (พฤติกรรมเมื่อ auth ล่ม ใช้ของ 001 ไปก่อน)*
- **Q-SH-03** โปรโตคอลของระบบยืนยันตัวตน -> ระหว่างรอใช้ auth จำลองของ 001 (ไม่ใช่ของฟีเจอร์นี้)

### ค้างจาก CR-003 (ไม่ใช่ Open Question แต่ยังไม่ได้ทำ)
- backlog BL-23 ยังเป็น Could "เลื่อนไว้ก่อน" ทั้งที่ส่วน RQ-20 ทำใน 002 แล้ว
- Use Case Spec UC-02 บน Miro ยังเขียนว่า NFR-PERF-01 ครอบคลุมแบบฟอร์มคำขอ ซึ่งทีมตัดสินว่าไม่ขยาย
- ที่เก็บไฟล์แนบยังไม่มี Constraint ใน 000-shared รองรับ (ASM-BRQ-07)
