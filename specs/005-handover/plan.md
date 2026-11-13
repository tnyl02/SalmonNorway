# Implementation Plan: บันทึกการส่งมอบอุปกรณ์ (Handover)
Spec: SPEC-HND-005 (specs/005-handover/spec.md, Draft v2.8) | Plan status: Draft v1.5 | Updated: 2569-11-12 (v1.1 ไม่กระทบฟีเจอร์อื่น, v1.2 รับผล CR-003, v1.3 รับผล CR-004, migration 0003 และ 004 spec v1, v1.4 ปิด Q-HND-11, Q-HND-13, Q-P4, v1.5 รับ spec v2.7-v2.8 migration 0009 และตาราง test ข้อ 6 ครบสำหรับ RTM)
Owner: 660710764 นายบรรลุ ฉิมสิงห์ (@mmoonn101219)

## 1. สรุปแนวทาง
- เจ้าหน้าที่ (CON-HND-01) ดูรายการคำขอยืมที่อนุมัติแล้วและยังไม่ส่งมอบ เลือก 1 คำขอ แล้วบันทึกการส่งมอบ (FR-HND-01)
- ก่อนบันทึก ระบบยืนยันตัวตนผู้ยืมผ่าน Auth API ภายใน 10 วินาที และเจ้าหน้าที่ต้องยืนยันว่าตรวจบัตรแล้วทุกครั้ง (FR-HND-02, 03)
- บันทึกการส่งมอบเป็น **transaction เดียว**: ล็อกอุปกรณ์ ตรวจว่า "ว่าง" สร้าง Loan เปลี่ยนสถานะอุปกรณ์และคำขอ เก็บรูป และเขียน audit log (FR-HND-04 ถึง 08)
- ตรวจทุกกฎที่ฝั่ง server หน้าจอตรวจซ้ำเพื่อความสะดวกเท่านั้น
- โค้ดอยู่ใน package ของตัวเอง `backend/internal/handover/` และ `frontend/src/pages/Handover/` ใช้ของกลาง (`auth`, `audit`, `platform`, `config`) ที่ 001 สร้างไว้

## 2. เทคโนโลยีที่ใช้
| สิ่งที่เลือก | มาจาก | หมายเหตุ |
|---|---|---|
| PostgreSQL 16 | CON-TECH-01 | ใช้ฐานข้อมูลเดียวกับ 001 |
| HTTPS, Docker Compose (on-campus server / cloud VM) | CON-TECH-02 | ใช้ `deploy/` เดิมของทีม |
| Go (net/http + pgx) | ทีมเลือกเอง ไม่ได้มาจาก spec | เหมือน 001 (AGENTS.md ตารางเทียบ Go) |
| React + Vite, vitest | ทีมเลือกเอง ไม่ได้มาจาก spec | ค่าเริ่มต้นของรายวิชา |
| golang-migrate (SQL ล้วน) | ทีมเลือกเอง ไม่ได้มาจาก spec | migration ใหม่เท่านั้น ห้ามแก้ 0001, 0002 |
| เก็บไฟล์รูปบน Docker volume (โฟลเดอร์บนดิสก์) | ทีมเลือกเอง ไม่ได้มาจาก spec (Q-P1 ปิดแล้ว) | SRS 2.3 บอกว่ามี "ที่เก็บไฟล์" แต่ไม่บอกชนิด ต้องรวม volume ใน backup รายวัน |
| auth จำลอง (`auth.MockProvider`) สำหรับเจ้าหน้าที่และผู้ยืม | ทีมเลือกเอง ไม่ได้มาจาก spec (Q-SH-03) | ผู้ยืมผ่าน interface `BorrowerVerifier` ผู้ยืมล็อกอินบนเครื่องเจ้าหน้าที่ (ASM-HND-09) ทีมตัดสิน 2569-11-12 เทอมนี้ใช้ตัวจำลอง (T-29 ยกเลิกเทอมนี้) |

## 3. โมเดลข้อมูล
ตาราง `equipment_item`, `loan` (0001) และ `audit_log` (0002) มีแล้ว ตาราง `borrow_request` เป็นของ 002 (migration 0003)
ฟีเจอร์นี้เพิ่มด้วย migration ใหม่ 2 ไฟล์ เลข **0006 และ 0009** (0006 ทีมกำหนด 2569-11-11 / ไฟล์ action ย้ายจาก 0007 เป็น 0009 ทีมเลือก 2569-11-12 เพราะ 0008 ของ 004 สร้าง CHECK ใหม่) `0006_handover.up.sql` และ `0009_audit_log_handover_action.up.sql`

| Entity / การเปลี่ยน | ฟิลด์หลัก | รองรับ |
|---|---|---|
| `loan` (DM-05) **ไม่แก้โครงตาราง** | INSERT ด้วยคอลัมน์เดิม `item_id`, `borrower_id` (รหัสนักศึกษาที่ยืนยันแล้ว), `loan_type`, `due_at` | FR-HND-07 |
| `handover` (ใหม่ เฉพาะ 005 ผูก 1:1 กับ `loan`) | `loan_id` PK/FK -> `loan`, `request_id` FK -> `borrow_request` UNIQUE, `handed_over_at TIMESTAMPTZ`, `handed_over_by` (รหัสเจ้าหน้าที่), `borrower_student_id`, `borrower_first_name`, `borrower_last_name`, `identity_method` CHECK IN (`auth_api`, `card_only`), `card_checked BOOLEAN NOT NULL CHECK (card_checked)` | FR-HND-02, 03, 05, 07 (ชื่อ-นามสกุลเจ้าหน้าที่พิมพ์ตามบัตร ASM-HND-10, รหัสนักศึกษา = รหัสผู้ใช้จากการยืนยัน ASM-HND-09) |
| `loan.due_at` (มีแล้ว) | ยืมปกติ = 16:00 ของวันส่งมอบ / ค้างคืน = 16:00 ของวันที่สิ้นสุดการยืมในคำขอ (Q-SH-04, CR-003) คำนวณใน 005 เอง ไม่อ่านค่าจากตารางของ 004 | FR-HND-05 (ส่งมอบหลัง 16:00 = 16:00 ของวันถัดไป Q-HND-03 ปิด 2569-11-12, ค้างคืนทำแล้ว Q-HND-07 ปิด 2569-11-12) |
| `handover_photo` (ใหม่) | `id`, `loan_id` FK -> `handover`, `file_path`, `content_type` CHECK IN (`image/jpeg`, `image/png`), `size_bytes`, `uploaded_at` | FR-HND-06 |
| `equipment_item.status` (มีแล้ว) | `available` -> `borrowed` | FR-HND-04, 07, CON-DATA-01 |
| `borrow_request.status` (ของ 002) | `approved` -> `handed_over` (data-model v9 ยืนยันแล้ว) อ่านค่าที่ส่งมอบได้จาก config | FR-HND-01, 07 |
| `audit_log.action` เพิ่มค่าใน CHECK (migration แยก) | ค่า `handover` (Q-HND-06 ปิด 2569-11-12) ใช้ `asset_code`, `request_id`, `outcome` เดิม บันทึกทั้งสำเร็จและไม่สำเร็จ (Q-HND-14 ปิด 2569-11-12) | FR-HND-08 |

สิ่งที่ **ไม่มีฟิลด์** โดยตั้งใจ: เลขบัตรประชาชน และเลขบัตรนักศึกษา FR-HND-03 สั่งให้บันทึกแค่ชื่อ นามสกุล รหัสนักศึกษา และเวลา การตรวจบัตรเก็บเป็น `card_checked` อย่างเดียว

ไม่กระทบฟีเจอร์อื่น (แก้ใน v1.1):
- ไม่ `ALTER TABLE loan` ข้อมูลการส่งมอบอยู่ในตาราง `handover` ของ 005 เอง seed และ test ของ 001 ที่ INSERT `loan` 4 คอลัมน์จึงไม่พัง
  และไม่ต้องแก้ DM-05 ใน data-model.md (ตามกติกา "ฟิลด์ที่ใช้เฉพาะฟีเจอร์เดียวอยู่ใน plan.md ของฟีเจอร์นั้น")
- เลข migration: 0003 ของ 002, 0004 และ 0005 ของ 003 (จองแล้ว) 005 ใช้ 2 เลขถัดไปที่ว่าง ต้องจองในกลุ่มก่อนสร้างไฟล์
  ไฟล์แรก `0006_handover` (ตาราง `handover`, `handover_photo`) ไฟล์ที่สอง `0009_audit_log_handover_action` (Q-HND-06 ปิดแล้ว)
- migration ของ action ต้อง `DROP CONSTRAINT` แล้วสร้าง CHECK ใหม่ที่มี**ทุกค่าที่มีอยู่ ณ เวลานั้น** ใช้ชื่อเดิม `audit_log_action_check` ตามที่ 0003 ของ 002 ตกลงกับ 003 (ตอนนี้มี list, search, detail, cancel, request บวก notify ของ 003 และของ 004 ถ้าไฟล์มาก่อน) บวก `handover` เพื่อไม่ลบค่าของคนอื่น

## 4. API / หน้าจอ
ทุก endpoint ต้องผ่าน `auth` middleware เดิม (IF-IDP-01) และตรวจ `role == staff` (CON-HND-01) ถ้าไม่ใช่ตอบ 403

| Method / Path | รองรับ | Input / Output หลัก |
|---|---|---|
| `GET /api/v1/handover/requests` | FR-HND-01 | out: `[{requestId, assetCode, modelName, requesterStudentId, submittedAt}]` เฉพาะสถานะอนุมัติแล้วและยังไม่ส่งมอบ |
| `POST /api/v1/handover/requests/{requestId}/verify-borrower` | FR-HND-02 | in: `{borrowerCredential}` ผู้ยืมล็อกอินบัญชีมหาวิทยาลัยบนเครื่องเจ้าหน้าที่ (ASM-HND-09, prototype = token จำลอง) out: `{result: "verified"|"mismatch"|"unavailable", userId}` ถ้า Auth API ไม่ตอบภายใน 10 วินาที = `unavailable` ใช้แสดงผลบนหน้าจอก่อนบันทึก |
| `POST /api/v1/handover/requests/{requestId}/handover` (multipart) | FR-HND-02 ถึง 08 | in: `borrowerCredential`, `borrowerFirstName`, `borrowerLastName` (เจ้าหน้าที่พิมพ์ตามบัตร ASM-HND-10), `cardChecked=true`, `photos[]` out 201: `{loanId, assetCode, handedOverAt, dueAt}` **หลังบ้านยืนยันผู้ยืมซ้ำเองด้วย `borrowerCredential` และกำหนด `identity_method` เอง ไม่เชื่อค่าที่หน้าจอส่งมา** รหัสนักศึกษาที่บันทึก = รหัสผู้ใช้จากการยืนยัน (ASM-HND-09) |

Error contract ของ `POST ./handover` ตอบเป็น `{code, message}` ตาม CON-TECH-03 (CR-004) ใช้รหัสจากรายการกลางเท่านั้น:
| กรณี | HTTP | code | อ้างอิง |
|---|---|---|---|
| ไม่ใช่เจ้าหน้าที่ | 403 | `FORBIDDEN_ROLE` | CON-HND-01, AC-HND-09 |
| คำขอไม่ใช่ `approved` (ยังไม่อนุมัติ หรือส่งมอบไปแล้ว) | 409 | `REQUEST_NOT_APPROVED` | FR-HND-01, AC-HND-01 |
| ผู้ยืมไม่ใช่ผู้ส่งคำขอ | 409 | `BORROWER_MISMATCH` | FR-HND-02, AC-HND-10 |
| ไม่ได้ยืนยันการตรวจบัตร | 422 | `VALIDATION_FAILED` (ช่องที่ผิด: `cardChecked`) | FR-HND-03, AC-HND-02, 03 |
| อุปกรณ์ไม่ "ว่าง" | 409 | `ITEM_NOT_AVAILABLE` | FR-HND-04, AC-HND-04, 13 |
| รูปไม่ครบ ชนิดผิด หรือใหญ่เกิน | 422 | `VALIDATION_FAILED` (ช่องที่ผิด: `photos`) | FR-HND-06, AC-HND-12 |
| ฐานข้อมูลล่ม | 503 | `UPSTREAM_UNAVAILABLE` | ใช้ error mapper เดิมใน `platform` |

ลำดับใน transaction (FR-HND-04, 07): `SELECT ... FROM equipment_item WHERE id = $1 FOR UPDATE` -> ตรวจ `available` -> ตรวจคำขอยังอยู่ในสถานะที่ส่งมอบได้ -> `INSERT loan` -> `INSERT handover` -> `UPDATE equipment_item SET status='borrowed'` -> `UPDATE borrow_request SET status='handed_over' WHERE id=$1 AND status = ANY($ready)` -> `INSERT handover_photo` -> commit
`FOR UPDATE` กันเจ้าหน้าที่ 2 คนส่งมอบอุปกรณ์ชิ้นเดียวกันพร้อมกัน ไม่แตะคำขออื่นของชิ้นเดียวกัน (ASM-HND-08)
ไฟล์รูปเขียนลงดิสก์ก่อน commit ถ้า transaction ล้มเหลว ลบไฟล์ทิ้ง audit log เขียนนอก transaction เพื่อให้บันทึกได้ทั้งกรณีสำเร็จและไม่สำเร็จ (ถ้า Q-HND-14 ตอบว่าบันทึก)
audit log: ถ้า `backend/internal/audit/writer.go` ของ 001 (T-19 ของ 001) มีแล้วใช้ของเดิม ถ้ายังไม่มี เขียน INSERT ใน repository ของ 005 เอง (แผนสำรองแบบเดียวกับ 003) ไม่สร้าง package `audit` เอง

ไฟล์ของกลางที่ต้องแตะ (เพิ่มต่อท้ายเท่านั้น ไม่แก้บรรทัดของคนอื่น):
- `config.go`, `config.js`: เพิ่มเป็นบล็อก `// ---- 005-handover ----` แบบที่ 003 ทำ
- `cmd/server/main.go`, `frontend/src/App.jsx`, `frontend/src/api/client.js`: เพิ่ม route และฟังก์ชันของ 005 ต่อท้าย ทำหลัง merge 003 เพื่อไม่ให้ชนกัน
- `deploy/docker-compose.yml`: เพิ่ม named volume ใหม่ `photos` mount ที่ `/data/photos` 005 ใช้แค่ `/data/photos/handover/` (006 ใช้โฟลเดอร์ข้างกันได้โดยไม่ต้องแก้ volume)
- `platform/migrate_test.go` ของ 001: ไม่แตะ test ของ migration 005 อยู่ใน `backend/internal/handover/`

หน้าจอ (spec ไม่มี UI ID และไม่มี mockup):
| หน้าจอ | รองรับ | แบบ |
|---|---|---|
| `HandoverList` (`frontend/src/pages/Handover/HandoverList.jsx`) | FR-HND-01 | ไม่มี mockup ทีมออกแบบเอง |
| `HandoverForm` (`frontend/src/pages/Handover/HandoverForm.jsx`) ยืนยันผู้ยืม, ช่องยืนยันตรวจบัตร, แนบรูป, ปุ่มบันทึก | FR-HND-02, 03, 06 | ไม่มี mockup ทีมออกแบบเอง |

## 5. ตารางตรวจ Constraints
| Constraint ID | ถูกนำไปใช้ที่ไหนใน plan | สถานะ |
|---|---|---|
| CON-TECH-01 | ข้อ 2, 3: PostgreSQL ตารางเดิม + migration ใหม่ | ใช้แล้ว |
| CON-TECH-02 | ข้อ 2: ใช้ HTTPS/Docker ของ `deploy/` เดิม ไม่เปิด endpoint ผ่าน HTTP | ใช้แล้ว |
| CON-TECH-03 | ข้อ 4: error contract ใช้ `{code, message}` และรหัสกลาง รหัสใหม่ 3 ตัวลงทะเบียนแล้ว (Q-P4) | ใช้แล้ว |
| CON-DATA-01 | ข้อ 3, 4: เปลี่ยนเฉพาะ `available` -> `borrowed` ไม่มีค่าใหม่ | ใช้แล้ว |
| CON-DATA-03 | ข้อ 4: ระบุอุปกรณ์ด้วย `asset_code` ใน response และ audit log | ใช้แล้ว |
| DOM-AUDIT-01 | ข้อ 3, 4: เขียน `audit_log` ทุกครั้งที่บันทึกการส่งมอบ ผ่าน package `audit` เดิม | ใช้แล้ว (action `handover` migration 0009 บันทึกทั้งสำเร็จและไม่สำเร็จ) |
| DOM-LOAN-01 | ข้อ 3, 5: `due_at` = 16:00 Asia/Bangkok ของวันส่งมอบ | ใช้แล้ว ส่วน "ใช้ได้เฉพาะในพื้นที่มหาวิทยาลัย" เป็นกฎที่เจ้าหน้าที่แจ้ง ระบบไม่ทำอะไร (Q-HND-04 ปิด 2569-11-12) |
| DOM-REQ-01 | ข้อ 4: สถานะเปลี่ยนเป็น `borrowed` ที่ transaction ของฟีเจอร์นี้เท่านั้น ไม่แตะคำขออื่น | ใช้แล้ว |
| IF-IDP-01 | ข้อ 4: `auth` middleware สำหรับเจ้าหน้าที่ + `BorrowerVerifier` สำหรับผู้ยืม timeout 10 วินาที | ใช้แล้ว (ผู้ยืมล็อกอินบนเครื่องเจ้าหน้าที่ ASM-HND-09 เทอมนี้ใช้ตัวจำลอง ไม่เชื่อม Auth API จริง) |
| CON-HND-01 | ข้อ 4: ตรวจ `role == staff` ทุก endpoint | ใช้แล้ว |

### ค่าที่ตั้งได้
| ค่า | มาจาก ID | ชื่อตัวแปรและไฟล์ |
|---|---|---|
| 10 วินาที | FR-HND-02, ASM-HND-02 | `BorrowerAuthTimeout` ใน `backend/internal/config/config.go` |
| 16:00 น. | DOM-LOAN-01, FR-HND-05 | `LoanNormalReturnHour = 16` ใน `config.go` **ของ 002** (merge เข้า main แล้ว 2569-11-12) 005 ใช้ตัวเดียวกัน ไม่ประกาศซ้ำ (กฎข้อ 9) |
| สถานะคำขอที่ส่งมอบได้ | FR-HND-01, data-model v9 | `HandoverReadyStatuses()` ใน `config.go` = `approved` |
| Asia/Bangkok | glossary, FR-HND-05 | `TimeZone` ใน `config.go` (มีแล้ว), `TIME_ZONE` ใน `frontend/src/config.js` (มีแล้ว) |
| อย่างน้อย 2 รูป | FR-HND-06, ASM-HND-07 | `HandoverMinPhotos` ใน `config.go`, `HANDOVER_MIN_PHOTOS` ใน `frontend/src/config.js` |
| ไม่เกิน 10 MB | FR-HND-06, ASM-HND-07 | `HandoverPhotoMaxBytes` ใน `config.go`, `HANDOVER_PHOTO_MAX_BYTES` ใน `config.js` (ต่อรูป ไม่จำกัดจำนวน Q-HND-10 ปิด 2569-11-12) |
| JPG, PNG | FR-HND-06, ASM-HND-07 | `HandoverPhotoTypes` ใน `config.go`, `HANDOVER_PHOTO_TYPES` ใน `config.js` |
| ที่เก็บไฟล์รูป | ทีมเลือกเอง ไม่ได้มาจาก spec (Q-P1) | `PhotoDir` ใน `config.go` อ่านจาก env |

## 6. แผนทดสอบจาก Acceptance Criteria
test หลังบ้านอยู่ที่ `backend/internal/handover/*_test.go` ใช้ PostgreSQL 16 จริงใน container (testcontainers) แบบเดียวกับ 001
เวลาและ `BorrowerVerifier` inject ได้ เพื่อทดสอบ 10:00, 16:30 และ timeout โดยไม่ต้องรอจริง
test หน้าจออยู่ที่ `frontend/src/__tests__/AC-HND-*.test.jsx` ใช้ client จำลอง
คำสั่งรัน: `cd backend && go test ./internal/handover/ -v` (ต้องมี Docker) และ `cd frontend && npx vitest run src/__tests__/AC-HND`
(v1.5 2569-11-12 ตารางนี้ตรงกับชื่อ test ที่มีจริงในโค้ดทุกแถว ใช้ทำ RTM ได้)

| AC ID | ชื่อ test | ทดสอบอย่างไร |
|---|---|---|
| AC-HND-01 | `TestAC_HND_01_RequestNotApproved` | คำขอสถานะยังไม่อนุมัติ เรียก POST handover ได้ 409 ไม่มี Loan อุปกรณ์ยัง `available` |
| AC-HND-02 | `TestAC_HND_02_AuthApiVerifiedAndCardChecked` (verify-borrower) + `TestAC_HND_02_HandoverRecordsBorrower` (บันทึก) + `frontend/src/__tests__/AC-HND-02.test.jsx` | verifier คืน verified + `cardChecked=true` ได้ 201 และ Loan มีชื่อ นามสกุล รหัสนักศึกษา เวลา / `cardChecked=false` ได้ 422 / หน้าจอไม่ติ๊กตรวจบัตรส่งไม่ได้ |
| AC-HND-03 | `TestAC_HND_03_AuthApiTimeoutCardOnly` | verifier ค้างจน context หมดเวลา (timeout สั้นใน test) ได้ `unavailable` แล้วบันทึกด้วย `card_only` ได้ 201 / ไม่ยืนยันบัตรได้ 422 |
| AC-HND-04 | `TestAC_HND_04_ItemNotAvailable` | อุปกรณ์ `unavailable` และ `borrowed` ได้ 409 สถานะไม่เปลี่ยน |
| AC-HND-05 | `TestAC_HND_05_DueAt1600SameDay` | clock = 2569-11-12 10:00 Asia/Bangkok ได้ `handed_over_at` 10:00 และ `due_at` 16:00 วันเดียวกัน |
| AC-HND-06 | `TestAC_HND_06_TwoJpgPhotosSaved` | แนบ JPG 2 รูป รูปละ 2 MB ได้ 201 มี `handover_photo` 2 แถวผูก Loan และมีไฟล์บนดิสก์ |
| AC-HND-07 | `TestAC_HND_07_HandoverCreatesLoan` | ได้ Loan 1 แถว `returned_at` ว่าง อุปกรณ์ `borrowed` คำขอเป็นสถานะส่งมอบแล้ว |
| AC-HND-08 | `TestAC_HND_08_AuditLogged` | หลังบันทึกสำเร็จและไม่สำเร็จ มี audit log ระบุผู้ทำ เวลา รหัสครุภัณฑ์ รหัสคำขอ ผลลัพธ์ (Q-HND-14 ปิด 2569-11-12) |
| AC-HND-09 | `TestAC_HND_09_StudentForbidden` | token นักศึกษา ได้ 403 ไม่มี Loan |
| AC-HND-10 | `TestAC_HND_10_BorrowerMismatch` | verifier คืนนักศึกษา B สำหรับคำขอของ A ได้ 409 อุปกรณ์ยัง `available` |
| AC-HND-11 | `TestAC_HND_11_HandoverAfter1600` | clock = 16:30 ได้ 201 `handed_over_at` 16:30 อุปกรณ์ `borrowed` และ `due_at` 16:00 ของวันถัดไป (Q-HND-03 ปิด 2569-11-12) |
| AC-HND-12 | `TestAC_HND_12_InvalidPhotos` + `frontend/src/__tests__/AC-HND-12.test.jsx` | (ก) 1 รูป (ข) GIF (ค) 11 MB ทุกกรณีได้ 422 ไม่มี Loan และไม่มีไฟล์ค้างบนดิสก์ / หน้าจอไม่ส่งไปหลังบ้าน |
| AC-HND-13 | `TestAC_HND_13_OtherRequestUntouched` | คำขอ X ส่งมอบสำเร็จ คำขอ Y สถานะไม่เปลี่ยน แล้วส่งมอบ Y ได้ 409 `ITEM_NOT_AVAILABLE`  (Given ไม่เกิดในระบบจริงหลัง Q-SH-06 ปิดบางส่วน รอทีมตัดสิน) |
| AC-HND-14 | `TestAC_HND_14_ListApprovedOnly` + `frontend/src/__tests__/AC-HND-14.test.jsx` | API คืนเฉพาะ 2 คำขอที่อนุมัติแล้วและยังไม่ส่งมอบ / หน้าจอแสดง 2 แถว |

test อื่นที่ไม่มี AC แยก (ตรวจช่อง "เสร็จเมื่อ" ของ task หรือ Constraint):
| ชื่อ test | ไฟล์ | รองรับ | task |
|---|---|---|---|
| `TestHandover_NonStaffForbidden` | handler_test.go | CON-HND-01, IF-IDP-01, CON-TECH-03 | T-03 |
| `TestVerifyBorrower_OtherResults` | handler_test.go | FR-HND-02 (mismatch, 422, 404) | T-09 |
| `TestHandover_ConcurrentSameItem` | repository_test.go | FR-HND-04, DOM-REQ-01 (`FOR UPDATE`) | T-18 |
| `TestHandover_OvernightDueAt` | handler_test.go | FR-HND-05 ยืมค้างคืน (Q-SH-04, Q-HND-07) spec ยังไม่มี AC | T-30 |
| `TestRepository_CreateHandover` | repository_test.go | FR-HND-04, FR-HND-07, CON-DATA-01 | T-06 |
| `TestMigration_HandoverIdentityMethod`, `TestMigration_HandoverPhotoContentType`, `TestMigration_HandoverNoCardNumberColumn` | migration_test.go | CON-TECH-01, FR-HND-02, 03, 06 (ไม่มีคอลัมน์เลขบัตร) | T-02 |
| `TestPhotoStore_SaveAndRemove` | photostore_test.go | FR-HND-06 | T-07 |
| `TestVerifier_Mock` | verifier_test.go | FR-HND-02, IF-IDP-01 | T-08 |
| `AC-HND-02.test.jsx` กรณี "T-12 ..." | frontend | FR-HND-02 (ยังไม่ยืนยันหรือ mismatch บันทึกไม่ได้) | T-12 |

## 7. ลำดับงาน
1. migration ใช้เลข 0006 และ 0009 (0006 ทีมกำหนด 2569-11-11, 0009 แทน 0007 ทีมเลือก 2569-11-12) `borrow_request` (0003) ของ 002 merge แล้ว (FR-HND-01, 07)
2. เพิ่มค่าที่ตั้งได้ใน `config.go` และ `frontend/src/config.js` (FR-HND-02, 05, 06)
3. migration `0006_handover`: ตาราง `handover` (ไม่แก้ `loan`) และ `handover_photo` แล้ว migration แยกสำหรับค่า action ของ `audit_log` (FR-HND-03, 06, 08)
4. repository + transaction ส่งมอบ ล็อกอุปกรณ์ (FR-HND-04, 07, AC-HND-04, 07, 13)
5. service: ตรวจบทบาท สถานะคำขอ เวลา/กำหนดคืน (CON-HND-01, FR-HND-01, 05, AC-HND-01, 05, 09, 11)
6. `BorrowerVerifier` + timeout 10 วินาที + การตรวจบัตร (FR-HND-02, 03, AC-HND-02, 03, 10)
7. รับและตรวจรูป เก็บไฟล์ ลบเมื่อ rollback (FR-HND-06, AC-HND-06, 12)
8. เขียน audit log ผ่าน writer.go ของ 001 ถ้ามีแล้ว ถ้ายังไม่มีใช้แผนสำรองใน repository ของ 005 (FR-HND-08, AC-HND-08)
9. `GET /handover/requests` และผูก route ใน `main.go` (FR-HND-01, AC-HND-14)
10. หน้าจอ `HandoverList` และ `HandoverForm` (FR-HND-01, 02, 03, 06)

## 8. สิ่งที่ยังไม่ทำ
ข้อที่ยังเปิด ส่วนที่เกี่ยวข้องจะยังไม่สร้างจนกว่าจะได้คำตอบ ข้อที่ปิดแล้วเก็บไว้เพื่อตามรอย

- Q-HND-03 [ปิดแล้ว 2569-11-12] ส่งมอบหลัง 16:00 กำหนดคืน = 16:00 ของวันถัดไป ส่งมอบตรง 16:00:00 นับเป็นวันเดียวกัน (FR-HND-05, AC-HND-11) ทำแล้วใน T-25
- Q-HND-04 [ปิดแล้ว 2569-11-12] "ใช้งานได้เฉพาะภายในพื้นที่มหาวิทยาลัย" เป็นกฎที่เจ้าหน้าที่แจ้งผู้ยืม ระบบไม่ทำอะไร (T-31 ยกเลิก)
- Q-HND-05 [ปิดแล้ว 2569-11-11 data-model v9 + migration 0003] ส่งมอบได้เมื่อ `approved` ตั้งเป็น `handed_over` ใช้ `requester_id`, `item_id`, `end_at` ของ `borrow_request`
  ไม่ต้องรอแล้ว
- Q-HND-06 [ปิดแล้ว 2569-11-12] ค่า `action` = `handover` ใน migration 0009 (หลัง 0008 ของ 004 ไม่แก้ไฟล์ของ 004) ทำแล้วใน T-26
- Q-HND-07 [ปิดแล้ว 2569-11-12] ทำการส่งมอบยืมค้างคืนเทอมนี้ กำหนดคืน 16:00 ของวันสิ้นสุดการยืมตามคำขอ ทำแล้วใน T-30 (spec ยังไม่มี AC ของค้างคืน)
  ยังเปิด (ยังไม่มี ID): ส่งมอบค้างคืนหลังวันสิ้นสุดการยืม กำหนดคืนเป็นอย่างไร -> ทีมตกลงและระบุใน spec ก่อนทำ
- Q-HND-09 การบันทึกการส่งมอบต้องตอบสนองภายในกี่วินาที -> ถามเจ้าหน้าที่
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ไม่มี test ด้านเวลา)
- Q-HND-10 [ปิดแล้ว 2569-11-12] 10 MB ต่อรูป ไม่จำกัดจำนวนรูป ทำแล้วใน T-28
- Q-HND-11 [ปิดแล้ว 2569-11-11] ผู้ยืมล็อกอินบัญชีมหาวิทยาลัยบนเครื่องเจ้าหน้าที่ (ASM-HND-09)
  ไม่ต้องรอแล้ว prototype ใช้ verifier จำลอง ทีมตัดสิน 2569-11-12 เทอมนี้ไม่เชื่อม Auth API จริง (T-29 ยกเลิกเทอมนี้)
- Q-HND-13 [ปิดแล้ว 2569-11-11] เจ้าหน้าที่พิมพ์ชื่อและนามสกุลตามบัตรนักศึกษาหรือบัตรประชาชน (ASM-HND-10)
  ไม่ต้องรอแล้ว
- Q-HND-14 [ปิดแล้ว 2569-11-12] บันทึกทั้งสำเร็จและไม่สำเร็จ outcome: นักศึกษาเรียกบันทึกการส่งมอบ = forbidden กรณีอื่น = error ดูรายการและยืนยันผู้ยืมไม่เก็บ ทำแล้วใน T-27
- Q-HND-15 ย้าย CON-HND-01 ไป constraints.md ไหม -> ทีม
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ไม่กระทบโค้ด กระทบแค่ที่อยู่ของข้อความ)
- Q-SH-04 [ปิดแล้ว 2569-11-11] กำหนดคืนยืมค้างคืน = 16:00 ของวันที่สิ้นสุดการยืมที่ผู้ขอระบุ (005 แปลงตอนส่งมอบ ตาม CR-003)
  ไม่ต้องรอแล้ว แปลงจาก `borrow_request.end_at` (ใช้แล้วใน T-30)
- Q-SH-06 [ปิดบางส่วน 2569-11-11] 004 ปฏิเสธคำขออื่นของชิ้นเดียวกันอัตโนมัติเมื่ออนุมัติ (005 ไม่แตะตาม ASM-HND-08 ยังสอดคล้อง)
  ผลต่อ plan: Given ของ AC-HND-13 ไม่เกิดในระบบจริง test ของ AC-HND-13 ยังเขียนได้ด้วยข้อมูลที่ตั้งเอง แต่ทีมต้องตัดสินว่าแก้หรือยกเลิก AC นี้ (ผ่าน CR)
- Q-SH-08 [ปิดสำหรับ prototype ทีมตัดสิน 2569-11-13] ระบบยืนยันตัวตนเชื่อมต่อไม่ได้ตอนเข้าใช้ทั่วไป
  คงตอบ 401 ของ middleware เดิมของ 001 ถ้าต่อระบบจริงต้องถามนักเทคโนโลยีสารสนเทศ

คำถามจากการวางแผน (ไม่มีใน spec ห้ามเดา):
- Q-P1 [ปิดแล้ว 2569-11-11 ทีมเลือกเอง ไม่ได้มาจาก spec] เก็บไฟล์รูปในโฟลเดอร์บน Docker volume (`/data/photos/handover/<loanId>/`)
  เก็บ path ในตาราง `handover_photo` volume นี้ต้องอยู่ใน backup รายวัน (SRS 2.3)
  ยังเหลือ: เก็บรูปนานเท่าไร -> ถามนักเทคโนโลยีสารสนเทศ และควรตกลงกับ 006 ให้ใช้ที่เก็บแบบเดียวกัน (RQ-04 รูปตอนรับคืน)
- Q-P2 [ปิดบางส่วน 2569-11-11] 004 ให้อุปกรณ์ 1 ชิ้นมีคำขออนุมัติแล้วได้ทีละ 1 รายการ (เดิมอ้าง AC-APV-08 ซึ่งไม่มีแล้วใน 004 spec v1 ตอนนี้คือ FR-APV-01 (ค))
  แก้ AC-HND-13 และ AC-HND-14 ใน spec Draft v2.3 แล้ว อ้าง FR-APV-01 (ค) ตั้งแต่ Draft v2.5
  ยังเปิด: ส่งมอบหลัง 16:00 กับข้อความใน Email ของ 004 (เนื้อหา Email รอ Q-APV-03 ของ 004, กำหนดคืนตอบแล้วใน Q-HND-03)
  และจะจำกัดให้เจ้าหน้าที่ส่งมอบได้เฉพาะอุปกรณ์ภาควิชาตนเองไหม (004 spec v1 ยังเปิดที่ Q-APV-08) -> ทีมตัดสิน
- Q-P3 [หมดไปแล้ว 2569-11-11] 004 spec v1 ไม่มีตารางที่เก็บกำหนดคืนตอนอนุมัติ ช่วงรับก่อน 16:00 และงานปฏิเสธอัตโนมัติ 16:00 แล้ว 005 เป็นที่เดียวที่บันทึกกำหนดคืน (ตรงกับ DM-05)
- Q-P4 [ปิดแล้ว 2569-11-11 ทีมเลือกทางที่ 1] ลงทะเบียน `FORBIDDEN_ROLE` (403), `REQUEST_NOT_APPROVED` (409), `BORROWER_MISMATCH` (409) ในตาราง CON-TECH-03 ของ `constraints.md`
  ตามกติกาของ CR-004 ("เพิ่มแถวในตารางก่อนใช้") ไม่ต้องทำ CR แต่ต้องแจ้งกลุ่ม
  ไม่ต้องรอแล้ว (T-03, T-16, T-20)
