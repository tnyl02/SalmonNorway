# Implementation Plan: บันทึกการส่งมอบอุปกรณ์ (Handover)
Spec: SPEC-HND-005 (specs/005-handover/spec.md, Draft v2.3) | Plan status: Draft v1 | Updated: 2569-11-11
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
| auth จำลอง (`auth.MockProvider`) สำหรับเจ้าหน้าที่และผู้ยืม | ทีมเลือกเอง ไม่ได้มาจาก spec (Q-SH-03) | ผู้ยืมผ่าน interface `BorrowerVerifier` เปลี่ยนเป็นตัวจริงได้เมื่อตอบ Q-HND-11 |

## 3. โมเดลข้อมูล
ตาราง `equipment_item`, `loan` (0001) และ `audit_log` (0002) มีแล้ว ตาราง `borrow_request` เป็นของ 002 (migration 0003)
ฟีเจอร์นี้เพิ่มด้วย migration ใหม่ 1 ไฟล์ **เลขรอจองในกลุ่ม** (data-model.md "0004 ขึ้นไป จองเลขในกลุ่มก่อนสร้างไฟล์") ชื่อชั่วคราว `00NN_handover.up.sql`

| Entity / การเปลี่ยน | ฟิลด์หลัก | รองรับ |
|---|---|---|
| `loan` (DM-05) เพิ่มคอลัมน์ | `request_id` FK -> `borrow_request`, `handed_over_at TIMESTAMPTZ`, `handed_over_by` (รหัสเจ้าหน้าที่), `borrower_student_id`, `borrower_first_name`, `borrower_last_name`, `identity_method` CHECK IN (`auth_api`, `card_only`), `card_checked BOOLEAN NOT NULL` | FR-HND-02, 03, 05, 07 (ชื่อ-นามสกุลรอ Q-HND-13) |
| `loan.due_at` (มีแล้ว) | ยืมปกติ = 16:00 ของวันส่งมอบ | FR-HND-05 (ส่งมอบหลัง 16:00 รอ Q-HND-03) |
| `handover_photo` (ใหม่) | `id`, `loan_id` FK, `file_path`, `content_type` CHECK IN (`image/jpeg`, `image/png`), `size_bytes`, `uploaded_at` | FR-HND-06 |
| `equipment_item.status` (มีแล้ว) | `available` -> `borrowed` | FR-HND-04, 07, CON-DATA-01 |
| `borrow_request.status` (ของ 002) | สถานะอนุมัติแล้ว -> "ส่งมอบแล้ว" ชื่อค่ารอ Q-HND-05 | FR-HND-01, 07 |
| `audit_log.action` เพิ่มค่าใน CHECK | ค่าสำหรับการส่งมอบ ชื่อรอ Q-HND-06 (ชื่อชั่วคราว `handover`) ใช้ `asset_code`, `request_id`, `outcome` เดิม | FR-HND-08 |

สิ่งที่ **ไม่มีฟิลด์** โดยตั้งใจ: เลขบัตรประชาชน และเลขบัตรนักศึกษา FR-HND-03 สั่งให้บันทึกแค่ชื่อ นามสกุล รหัสนักศึกษา และเวลา การตรวจบัตรเก็บเป็น `card_checked` อย่างเดียว

การแก้ `loan` เป็นการเปลี่ยน DM-05 ใน data-model.md ของกลาง ต้องแจ้งทีมก่อน (Q-HND-13) และใช้ `ALTER TABLE` ใน migration ใหม่ ไม่แก้ 0001

## 4. API / หน้าจอ
ทุก endpoint ต้องผ่าน `auth` middleware เดิม (IF-IDP-01) และตรวจ `role == staff` (CON-HND-01) ถ้าไม่ใช่ตอบ 403

| Method / Path | รองรับ | Input / Output หลัก |
|---|---|---|
| `GET /api/v1/handover/requests` | FR-HND-01 | out: `[{requestId, assetCode, modelName, requesterStudentId, submittedAt}]` เฉพาะสถานะอนุมัติแล้วและยังไม่ส่งมอบ |
| `POST /api/v1/handover/requests/{requestId}/verify-borrower` | FR-HND-02 | in: ข้อมูลยืนยันตัวตนผู้ยืม (รูปแบบรอ Q-HND-11) out: `{result: "verified"\|"mismatch"\|"unavailable"}` ถ้า Auth API ไม่ตอบภายใน 10 วินาที = `unavailable` |
| `POST /api/v1/handover/requests/{requestId}/handover` (multipart) | FR-HND-02 ถึง 08 | in: `cardChecked=true`, `identityMethod`, `photos[]` out 201: `{loanId, assetCode, handedOverAt, dueAt}` |

Error contract ของ `POST .../handover`:
| กรณี | HTTP | code | อ้างอิง |
|---|---|---|---|
| ไม่ใช่เจ้าหน้าที่ | 403 | `FORBIDDEN_ROLE` | CON-HND-01, AC-HND-09 |
| คำขอยังไม่อนุมัติ หรือส่งมอบไปแล้ว | 409 | `REQUEST_NOT_APPROVED` | FR-HND-01, AC-HND-01 |
| ผู้ยืมไม่ใช่ผู้ส่งคำขอ | 409 | `BORROWER_MISMATCH` | FR-HND-02, AC-HND-10 |
| ไม่ได้ยืนยันการตรวจบัตร | 422 | `CARD_NOT_CHECKED` | FR-HND-03, AC-HND-02, 03 |
| อุปกรณ์ไม่ "ว่าง" | 409 | `ITEM_NOT_AVAILABLE` | FR-HND-04, AC-HND-04, 13 |
| รูปไม่ครบ ชนิดผิด หรือใหญ่เกิน | 422 | `PHOTO_INVALID` | FR-HND-06, AC-HND-12 |
| ฐานข้อมูลล่ม | 503 | `UPSTREAM_UNAVAILABLE` | ใช้ error mapper เดิมใน `platform` |

ลำดับใน transaction (FR-HND-04, 07): `SELECT ... FROM equipment_item WHERE id = $1 FOR UPDATE` -> ตรวจ `available` -> ตรวจคำขอยังอนุมัติอยู่ -> `INSERT loan` -> `UPDATE equipment_item SET status='borrowed'` -> `UPDATE borrow_request` -> `INSERT handover_photo` -> commit
`FOR UPDATE` กันเจ้าหน้าที่ 2 คนส่งมอบอุปกรณ์ชิ้นเดียวกันพร้อมกัน ไม่แตะคำขออื่นของชิ้นเดียวกัน (ASM-HND-08)
ไฟล์รูปเขียนลงดิสก์ก่อน commit ถ้า transaction ล้มเหลว ลบไฟล์ทิ้ง audit log เขียนนอก transaction เพื่อให้บันทึกได้ทั้งกรณีสำเร็จและไม่สำเร็จ (ถ้า Q-HND-14 ตอบว่าบันทึก)

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
| CON-DATA-01 | ข้อ 3, 4: เปลี่ยนเฉพาะ `available` -> `borrowed` ไม่มีค่าใหม่ | ใช้แล้ว |
| CON-DATA-03 | ข้อ 4: ระบุอุปกรณ์ด้วย `asset_code` ใน response และ audit log | ใช้แล้ว |
| DOM-AUDIT-01 | ข้อ 3, 4: เขียน `audit_log` ทุกครั้งที่บันทึกการส่งมอบ ผ่าน package `audit` เดิม | ใช้แล้ว (ค่า action รอ Q-HND-06) |
| DOM-LOAN-01 | ข้อ 3, 5: `due_at` = 16:00 Asia/Bangkok ของวันส่งมอบ | ใช้แล้วบางส่วน ส่วน "ใช้ได้เฉพาะในพื้นที่มหาวิทยาลัย" ยังไม่ได้ใช้ เพราะรอ Q-HND-04 |
| DOM-REQ-01 | ข้อ 4: สถานะเปลี่ยนเป็น `borrowed` ที่ transaction ของฟีเจอร์นี้เท่านั้น ไม่แตะคำขออื่น | ใช้แล้ว |
| IF-IDP-01 | ข้อ 4: `auth` middleware สำหรับเจ้าหน้าที่ + `BorrowerVerifier` สำหรับผู้ยืม timeout 10 วินาที | ใช้แล้ว (วิธียืนยันผู้ยืมรอ Q-HND-11) |
| CON-HND-01 | ข้อ 4: ตรวจ `role == staff` ทุก endpoint | ใช้แล้ว |

### ค่าที่ตั้งได้
| ค่า | มาจาก ID | ชื่อตัวแปรและไฟล์ |
|---|---|---|
| 10 วินาที | FR-HND-02, ASM-HND-02 | `BorrowerAuthTimeout` ใน `backend/internal/config/config.go` |
| 16:00 น. | DOM-LOAN-01, FR-HND-05 | `NormalLoanDueHour = 16`, `NormalLoanDueMinute = 0` ใน `config.go` |
| Asia/Bangkok | glossary, FR-HND-05 | `TimeZone` ใน `config.go` (มีแล้ว), `TIME_ZONE` ใน `frontend/src/config.js` (มีแล้ว) |
| อย่างน้อย 2 รูป | FR-HND-06, ASM-HND-07 | `HandoverMinPhotos` ใน `config.go`, `HANDOVER_MIN_PHOTOS` ใน `frontend/src/config.js` |
| ไม่เกิน 10 MB | FR-HND-06, ASM-HND-07 | `HandoverPhotoMaxBytes` ใน `config.go`, `HANDOVER_PHOTO_MAX_BYTES` ใน `config.js` (ต่อรูปหรือรวม รอ Q-HND-10) |
| JPG, PNG | FR-HND-06, ASM-HND-07 | `HandoverPhotoTypes` ใน `config.go`, `HANDOVER_PHOTO_TYPES` ใน `config.js` |
| ที่เก็บไฟล์รูป | ทีมเลือกเอง ไม่ได้มาจาก spec (Q-P1) | `PhotoDir` ใน `config.go` อ่านจาก env |

## 6. แผนทดสอบจาก Acceptance Criteria
test หลังบ้านอยู่ที่ `backend/internal/handover/handler_test.go` ใช้ PostgreSQL จริงแบบเดียวกับ 001
เวลาและ `BorrowerVerifier` inject ได้ เพื่อทดสอบ 10:00, 16:30 และ timeout โดยไม่ต้องรอจริง

| AC ID | ชื่อ test | ทดสอบอย่างไร |
|---|---|---|
| AC-HND-01 | `TestAC_HND_01_RequestNotApproved` | คำขอสถานะยังไม่อนุมัติ เรียก POST handover ได้ 409 ไม่มี Loan อุปกรณ์ยัง `available` |
| AC-HND-02 | `TestAC_HND_02_AuthApiVerifiedAndCardChecked` | verifier คืน verified + `cardChecked=true` ได้ 201 และ Loan มีชื่อ นามสกุล รหัสนักศึกษา เวลา / `cardChecked=false` ได้ 422 |
| AC-HND-03 | `TestAC_HND_03_AuthApiTimeoutCardOnly` | verifier ค้างจน context หมดเวลา (timeout สั้นใน test) ได้ `unavailable` แล้วบันทึกด้วย `card_only` ได้ 201 / ไม่ยืนยันบัตรได้ 422 |
| AC-HND-04 | `TestAC_HND_04_ItemNotAvailable` | อุปกรณ์ `unavailable` และ `borrowed` ได้ 409 สถานะไม่เปลี่ยน |
| AC-HND-05 | `TestAC_HND_05_DueAt1600SameDay` | clock = 2569-11-12 10:00 Asia/Bangkok ได้ `handed_over_at` 10:00 และ `due_at` 16:00 วันเดียวกัน |
| AC-HND-06 | `TestAC_HND_06_TwoJpgPhotosSaved` | แนบ JPG 2 รูป รูปละ 2 MB ได้ 201 มี `handover_photo` 2 แถวผูก Loan และมีไฟล์บนดิสก์ |
| AC-HND-07 | `TestAC_HND_07_HandoverCreatesLoan` | ได้ Loan 1 แถว `returned_at` ว่าง อุปกรณ์ `borrowed` คำขอเป็นสถานะส่งมอบแล้ว |
| AC-HND-08 | `TestAC_HND_08_AuditLogged` | หลังบันทึกสำเร็จ มี audit log ระบุผู้ทำ เวลา รหัสครุภัณฑ์ รหัสคำขอ ผลลัพธ์ (กรณีไม่สำเร็จรอ Q-HND-14) |
| AC-HND-09 | `TestAC_HND_09_StudentForbidden` | token นักศึกษา ได้ 403 ไม่มี Loan |
| AC-HND-10 | `TestAC_HND_10_BorrowerMismatch` | verifier คืนนักศึกษา B สำหรับคำขอของ A ได้ 409 อุปกรณ์ยัง `available` |
| AC-HND-11 | `TestAC_HND_11_HandoverAfter1600` | clock = 16:30 ได้ 201 `handed_over_at` 16:30 อุปกรณ์ `borrowed` **ไม่ตรวจ `due_at`** จนกว่าจะตอบ Q-HND-03 |
| AC-HND-12 | `TestAC_HND_12_InvalidPhotos` | (ก) 1 รูป (ข) GIF (ค) 11 MB ทุกกรณีได้ 422 ไม่มี Loan และไม่มีไฟล์ค้างบนดิสก์ |
| AC-HND-13 | `TestAC_HND_13_OtherRequestUntouched` | คำขอ X ส่งมอบสำเร็จ คำขอ Y สถานะไม่เปลี่ยน แล้วส่งมอบ Y ได้ 409 `ITEM_NOT_AVAILABLE` |
| AC-HND-14 | `TestAC_HND_14_ListApprovedOnly` + `frontend/src/__tests__/AC-HND-14.test.jsx` | API คืนเฉพาะ 2 คำขอที่อนุมัติแล้วและยังไม่ส่งมอบ / หน้าจอแสดง 2 แถว |

เพิ่มเติม (ไม่มี AC แยก): test ส่งมอบพร้อมกัน 2 goroutine ชิ้นเดียวกัน สำเร็จได้แค่ 1 (ยืนยัน `FOR UPDATE` ของ FR-HND-04)

## 7. ลำดับงาน
1. จองเลข migration ในกลุ่ม รอ 002 สร้าง `borrow_request` (0003) และ 004 ใช้สถานะอนุมัติได้ (FR-HND-01, 07)
2. เพิ่มค่าที่ตั้งได้ใน `config.go` และ `frontend/src/config.js` (FR-HND-02, 05, 06)
3. migration `00NN_handover`: คอลัมน์ใหม่ของ `loan`, ตาราง `handover_photo`, ค่า action ใหม่ของ `audit_log` (FR-HND-03, 06, 08)
4. repository + transaction ส่งมอบ ล็อกอุปกรณ์ (FR-HND-04, 07, AC-HND-04, 07, 13)
5. service: ตรวจบทบาท สถานะคำขอ เวลา/กำหนดคืน (CON-HND-01, FR-HND-01, 05, AC-HND-01, 05, 09, 11)
6. `BorrowerVerifier` + timeout 10 วินาที + การตรวจบัตร (FR-HND-02, 03, AC-HND-02, 03, 10)
7. รับและตรวจรูป เก็บไฟล์ ลบเมื่อ rollback (FR-HND-06, AC-HND-06, 12)
8. เขียน audit log (FR-HND-08, AC-HND-08)
9. `GET /handover/requests` และผูก route ใน `main.go` (FR-HND-01, AC-HND-14)
10. หน้าจอ `HandoverList` และ `HandoverForm` (FR-HND-01, 02, 03, 06)

## 8. สิ่งที่ยังไม่ทำ
ส่วนที่เกี่ยวข้องกับแต่ละข้อจะยังไม่สร้างจนกว่าจะได้คำตอบ

- Q-HND-03 ส่งมอบยืมปกติหลัง 16:00 น. แล้วกำหนดคืนคือเวลาไหน (FR-HND-05, AC-HND-11) -> ถามนักคอมพิวเตอร์
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (การคำนวณ `due_at` เมื่อส่งมอบหลัง 16:00)
- Q-HND-04 "ใช้งานได้เฉพาะภายในพื้นที่มหาวิทยาลัย" ระบบต้องทำอะไร -> ถามนักคอมพิวเตอร์
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ
- Q-HND-05 ชื่อสถานะ "ส่งมอบแล้ว" และฟิลด์รหัสผู้ส่งคำขอ/รหัสครุภัณฑ์ในคำขอยืม -> เจ้าของ 002
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (การ UPDATE `borrow_request` และ FK)
- Q-HND-06 ค่า `action` ของ audit log สำหรับการส่งมอบ (จองในกลุ่ม)
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ใช้ชื่อชั่วคราว `handover` ใน plan เท่านั้น)
- Q-HND-07 เทอมนี้ต้องทดสอบการส่งมอบยืมแบบค้างคืนไหม -> ทีมตัดสิน
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (plan รองรับเฉพาะ `loan_type = normal`)
- Q-HND-09 การบันทึกการส่งมอบต้องตอบสนองภายในกี่วินาที -> ถามเจ้าหน้าที่
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ไม่มี test ด้านเวลา)
- Q-HND-10 "ไม่เกิน 10 MB" ต่อรูปหรือรวม และแนบได้มากสุดกี่รูป -> ถามเจ้าหน้าที่
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (การตรวจขนาดรูป)
- Q-HND-11 ผู้ยืมยืนยันตัวตนผ่าน Auth API ที่เคาน์เตอร์อย่างไร -> ถามนักเทคโนโลยีสารสนเทศ
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ใช้ verifier จำลองเท่านั้น)
- Q-HND-13 ชื่อและนามสกุลผู้ยืมได้มาจากไหน และการเพิ่มคอลัมน์ใน DM-05 -> ทีม / นักเทคโนโลยีสารสนเทศ
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ
- Q-HND-14 การส่งมอบที่ไม่สำเร็จต้องบันทึก audit log ด้วยไหม -> ถามเจ้าหน้าที่
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ
- Q-HND-15 ย้าย CON-HND-01 ไป constraints.md ไหม -> ทีม
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ไม่กระทบโค้ด กระทบแค่ที่อยู่ของข้อความ)
- Q-SH-04 กำหนดคืนของการยืมแบบค้างคืน -> ถามนักคอมพิวเตอร์
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ
- Q-SH-06 คำขออื่นของอุปกรณ์ชิ้นเดียวกันหลังอนุมัติ/ส่งมอบ (005 ไม่แตะตาม ASM-HND-08) -> นักคอมพิวเตอร์
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ
- Q-SH-08 ระบบยืนยันตัวตนเชื่อมต่อไม่ได้ตอนเข้าใช้ทั่วไป -> ถามนักเทคโนโลยีสารสนเทศ
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ใช้ middleware เดิมของ 001)

คำถามจากการวางแผน (ไม่มีใน spec ห้ามเดา):
- Q-P1 [ปิดแล้ว 2569-11-11 ทีมเลือกเอง ไม่ได้มาจาก spec] เก็บไฟล์รูปในโฟลเดอร์บน Docker volume (`/data/photos/handover/<loanId>/`)
  เก็บ path ในตาราง `handover_photo` volume นี้ต้องอยู่ใน backup รายวัน (SRS 2.3)
  ยังเหลือ: เก็บรูปนานเท่าไร -> ถามนักเทคโนโลยีสารสนเทศ และควรตกลงกับ 006 ให้ใช้ที่เก็บแบบเดียวกัน (RQ-04 รูปตอนรับคืน)
- Q-P2 [ปิดบางส่วน 2569-11-11] AC-APV-08 ของ 004 ให้อุปกรณ์ 1 ชิ้นมีคำขออนุมัติแล้วได้ทีละ 1 รายการ
  แก้ AC-HND-13 (คำขอ Y ยังไม่ได้ตัดสิน) และ AC-HND-14 (อุปกรณ์คนละชิ้น) ใน spec Draft v2.3 แล้ว
  ยังเปิด: ส่งมอบหลัง 16:00 ขัดกับ Email อนุมัติของ 004 ที่บอกคืน 16:00 วันเดียวกัน (Q-HND-03)
  และจะจำกัดให้เจ้าหน้าที่ส่งมอบได้เฉพาะอุปกรณ์ภาควิชาตนเองแบบ AC-APV-11 ของ 004 ไหม -> ทีมตัดสิน
