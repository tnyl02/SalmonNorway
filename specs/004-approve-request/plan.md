# Implementation Plan: อนุมัติคำขอยืม ปกติและค้างคืน (Approve Request)
Spec: SPEC-APV-004 (specs/004-approve-request/spec.md, Draft v3) | Plan status: Draft v2 | Updated: 2569-11-13 (v2 ตาม CR-009 คำขอยืมค้างคืน FR-APV-03, AC-APV-22 ถึง 26 / v1.2 ตาม CR-008: Q-APV-28, 31 / v1.1 ตามคำตอบ Q-APV-14 ถึง 21)
Owner: 660710758 นางสาวณัฐสุดา กัญญาบุญ (@natsuda05)
อ้างอิง: specs/000-shared/{constraints.md v5.4, data-model.md v10, glossary.md}, migration 0003 ของ 002, endpoint ไฟล์แนบของ 002 (FR-BRQ-10 ก), specs/changes/CR-009.md

## 1. สรุปแนวทาง
- เจ้าหน้าที่ภาควิชาเดียวกับอุปกรณ์เปิดหน้ารายการ เห็นคำขอยืมแบบปกติ `requested` เรียงตามเวลาที่ส่ง แล้วอนุมัติหรือปฏิเสธ (FR-APV-01)
- หลังบ้านเป็น package `backend/internal/approve/` หน้าจอเป็น `frontend/src/pages/ApproveRequest/` ตรวจบทบาท ภาควิชา ลำดับ สถานะ และเวลาที่ฝั่ง server ทั้งหมด
- บันทึกผลใน transaction เดียว ล็อกแถวคำขอทุกใบของอุปกรณ์ชิ้นนั้น ปฏิเสธคำขออื่นของชิ้นเดียวกันในครั้งเดียว และบันทึกประวัติ (FR-APV-01 จ ฌ)
- Email เขียนลงตาราง outbox ใน transaction เดียวกับผล แล้วตัวส่งเบื้องหลังส่งภายใน 60 วินาที ส่งซ้ำ 3 ครั้งห่าง 1 นาที (FR-APV-02)
- คำขอยืมค้างคืน (FR-APV-03, CR-009): ขึ้นในรายการเดียวกันพร้อมประเภทการยืมและลิงก์เปิดใบบันทึกข้อความของ 002 เจ้าหน้าที่ตรวจเอกสารแล้วกดอนุมัติครั้งเดียวพร้อมชื่ออาจารย์
  ไม่มีคิว ไม่มีเวลาตัด ไม่ปฏิเสธอัตโนมัติ ใช้ endpoint เดิมแยกกฎตาม `loan_type` ฝั่ง server
- งานตั้งเวลาใน process ปฏิเสธคำขอยืมปกติที่ค้างเมื่อถึง 00:00 น. หลังวันที่ส่ง และไล่ปฏิเสธคำขอที่เลยเวลาแล้วทันทีตอน server เริ่มทำงาน (FR-APV-01 ฉ, Q-APV-21)

## 2. เทคโนโลยีที่ใช้
| สิ่งที่เลือก | มาจาก | หมายเหตุ |
|---|---|---|
| PostgreSQL | CON-TECH-01 | ตารางเดิมจาก migration 0001-0003 + migration ใหม่ของ 004 |
| HTTPS เท่านั้น, deploy ด้วย Docker | CON-TECH-02, NFR-SEC-01 | ใช้ server, proxy และ compose ของกลาง ไม่เปิดพอร์ตใหม่ |
| error `{code, message}` + รหัสจากตารางกลาง | CON-TECH-03 | รหัสใหม่ต้องลงทะเบียนใน constraints.md ก่อนใช้ (ข้อ 4) |
| Go (net/http + pgx) | ทีมเลือกเอง ไม่ได้มาจาก spec | ตาม AGENTS.md และ backend ที่มีอยู่ |
| React + Vite, test ด้วย vitest | ทีมเลือกเอง ไม่ได้มาจาก spec | ตาม frontend ที่มีอยู่ |
| golang-migrate, testcontainers (PostgreSQL 16) | ทีมเลือกเอง ไม่ได้มาจาก spec | ใช้ `platform.Migrate` และแบบเดียวกับ `platform/migrate_test.go` |
| auth จำลอง (`backend/internal/auth/mock_provider.go`) | ASM-02, IF-IDP-01 | โปรโตคอลจริงรอ Q-SH-03 |
| interface `Mailer` + ตัวส่งจำลอง | ทีมเลือกเอง ไม่ได้มาจาก spec | ใช้แทนระบบแจ้งเตือน Email (IF-NOT-01) ใน prototype เปลี่ยนเป็นตัวจริงได้โดยไม่แก้ส่วนอื่น |
| interface `RecipientResolver` + ตัวจำลอง | ทีมเลือกเอง ไม่ได้มาจาก spec | หา Email ผู้ขอ ระหว่างรอ Q-APV-04 |
| งานเบื้องหลังเป็น goroutine + นาฬิกาที่ inject ได้ | ทีมเลือกเอง ไม่ได้มาจาก spec | ตัวส่ง Email และงานเที่ยงคืน ไม่ใช้ cron หรือ queue ภายนอก |

## 3. โมเดลข้อมูล
| entity | ฟิลด์หลัก | รองรับ | เจ้าของ |
|---|---|---|---|
| `borrow_request` (DM-06) | id, requester_id, requester_name, student_code, item_id, loan_type, purpose, purpose_other, start_at, end_at, status, submitted_at, attachment | อ่าน และเปลี่ยน status: ยืมปกติ `requested` -> `approved` / `rejected` (FR-APV-01) / ค้างคืน `requested` -> `approved_by_officer` -> `pending_advisor_approval` -> `approved` ใน transaction เดียว หรือ `requested` -> `rejected` (FR-APV-03) | 002 (migration 0003) |
| `equipment_item`, `equipment_model` (DM-04, DM-03) | status, asset_code, ชื่อรุ่น, department_id | อ่านอย่างเดียว: อุปกรณ์ "ว่าง" (FR-APV-01 ค) และภาควิชา (FR-APV-01 ก) | 001 (migration 0001) |
| `approval_decision` (ใหม่ เฉพาะ 004) | id, request_id (UNIQUE, FK), result (`approved` / `rejected` / `auto_rejected`), reason, actor_id, actor_name, actor_department_id, decided_at, advisor_name (เพิ่มใน 0011 ว่างได้ ใช้กับค้างคืนที่อนุมัติ) | ประวัติ RQ-25: FR-APV-01 ฌ, ข้อมูลใน Email: FR-APV-02 | 004 |
| `approval_step` (ใหม่ใน 0011 เฉพาะ 004) | id, request_id (FK), status (`approved_by_officer` / `pending_advisor_approval` / `approved`), actor_id, created_at | ประวัติสถานะระหว่างทางของคำขอค้างคืนที่ผ่านในครั้งเดียว: FR-APV-03 ข (ทีมเลือกเอง ไม่ได้มาจาก spec) | 004 |
| `approval_email` (ใหม่ เฉพาะ 004) | id, request_id, kind (`approved` / `rejected` / `auto_rejected`), status (`pending` / `sent` / `failed`), attempts, last_error, last_attempt_at, next_attempt_at, created_at | ส่ง Email และส่งซ้ำ: FR-APV-02 ง | 004 |
| `audit_log` (DM-07) | actor_id, actor_role, action, request_id, asset_code, outcome | บันทึกทุกการพิจารณา (RQ-25) เพิ่มค่า action `approve`, `reject`, `auto_reject` | 001 (migration 0002) |

- ไม่เปลี่ยน `equipment_item.status` (CON-DATA-01, data-model.md แถว UC-12)
- งานเที่ยงคืนและเวลาตัด 16:00 น. ใช้เฉพาะ `loan_type = 'normal'` (FR-APV-01 ฉ ค, FR-APV-03 ง) / หน้ารายการ อนุมัติ และปฏิเสธ รับทั้ง 2 ประเภท แยกกฎตาม loan_type (v2)
- ลำดับคำขอแรก (`is_first`, `NOT_FIRST_IN_QUEUE`) ใช้เฉพาะยืมปกติ คำขอค้างคืนไม่มีคิวและไม่ขวางคิว (FR-APV-03 ค, Q-APV-23) / อุปกรณ์ต้อง "ว่าง" ทั้ง 2 ประเภท (Q-APV-42)
- การปฏิเสธตามเมื่อมีคำขออื่นได้อนุมัติ (FR-APV-01 จ) **ไม่กรอง** loan_type: คำขอ `requested` ทุกประเภทของชิ้นนั้นถูกปฏิเสธ (Q-APV-14)
  จึงต้องล็อกคำขอทุกประเภทของชิ้นนั้นด้วย
- `approval_decision` ไม่มี UPDATE / DELETE และ request_id UNIQUE จึงพิจารณาซ้ำหรือเปลี่ยนผลไม่ได้ (Out of scope)
- ผลแต่ละแบบ:
  | กรณี | `borrow_request.status` | `approval_decision.result` | reason | actor |
  |---|---|---|---|---|
  | อนุมัติ | `approved` | `approved` | ว่าง | เจ้าหน้าที่ |
  | ปฏิเสธ | `rejected` | `rejected` | ที่เจ้าหน้าที่พิมพ์ | เจ้าหน้าที่ |
  | ถูกปฏิเสธตามเพราะมีคำขออื่นได้อนุมัติ (FR-APV-01 จ) | `rejected` | `rejected` | `ApprovalSupersededReason` | เจ้าหน้าที่ที่กดอนุมัติ |
  | เที่ยงคืน (FR-APV-01 ฉ) | `rejected` | `auto_rejected` | `ApprovalAutoRejectReason` | `system` / `System` + ภาควิชาเจ้าของอุปกรณ์ |
  | อนุมัติค้างคืน (FR-APV-03 ข) | `approved` (ผ่าน 2 สถานะระหว่างทางใน transaction เดียว บันทึกใน `approval_step` 3 แถว) | `approved` + advisor_name | ว่าง | เจ้าหน้าที่ |
  | ปฏิเสธค้างคืน (FR-APV-03 จ) | `rejected` | `rejected` | ที่เจ้าหน้าที่พิมพ์ ไม่มี advisor_name (Q-APV-43) | เจ้าหน้าที่ |
- ไม่มีคอลัมน์จำนวนอุปกรณ์ เพราะ 1 คำขอ = อุปกรณ์รายชิ้น 1 ชิ้น (FR-BRQ-01, ASM-BRQ-02) และไม่มีคอลัมน์ Email (รอ Q-APV-04)
- `actor_name` ว่างได้ เพราะ Principal (DM-08) ยังไม่มีชื่อ ดูข้อ 8
- migration ใหม่ของ CR-009 คือ `0011_approval_overnight` (เลขถัดจาก 0010 ของ 006 ตรวจแล้วว่ายังไม่มีใน branch ใด 2569-11-13): ADD COLUMN advisor_name, CREATE TABLE approval_step (append-only แบบ approval_decision) ไม่แตะ CHECK ของ audit_log (ใช้ action `approve`, `reject` เดิม)
- migration ของ 004 คือ `0008_approval` (ทีมจองให้ 2569-11-12) และเพิ่มค่า action ด้วยวิธีเดียวกับ 0003 (DROP แล้ว ADD `audit_log_action_check`) ห้ามแก้ migration เดิม

## 4. API / หน้าจอ
ทุก endpoint ขึ้นต้นด้วย `/v1/approvals/` ผ่าน `auth.Middleware` (IF-IDP-01) แล้วตรวจ `Role == staff` (AC-APV-05)
และภาควิชาของเจ้าหน้าที่ = ภาควิชาเจ้าของอุปกรณ์ (ASM-06) ผูกใน `backend/cmd/server/main.go` แบบเดียวกับ `equipment.PathPrefix`

| method path | input / output หลัก | รองรับ |
|---|---|---|
| GET /v1/approvals/requests | out: คำขอ `requested` ทั้งยืมปกติและค้างคืนของภาควิชา เรียงตาม submitted_at (รหัสคำขอ, ชื่อผู้ขอ, ชื่อรุ่น, รหัสครุภัณฑ์, วัตถุประสงค์, start_at, end_at, submitted_at, `loan_type`, `is_first` เฉพาะยืมปกติ, `attachment_url` เฉพาะค้างคืน) และ `loaded_at` | FR-APV-01 ก ข, FR-APV-03 ก, AC-APV-11 |
| POST /v1/approvals/requests/{id}/approve | in: `loaded_at`, `advisorName` (บังคับเฉพาะค้างคืน ยืมปกติไม่ใช้) / out: ผล และรายการคำขอที่ถูกปฏิเสธตาม | FR-APV-01 ข ค จ ช, FR-APV-03 ข ค ง ฉ, AC-APV-01, 03, 07, 08, 09, 10, 16, 22, 23, 24, 26 |
| POST /v1/approvals/requests/{id}/reject | in: `reason`, `loaded_at` / out: ผล (ใช้กับทั้ง 2 ประเภท) | FR-APV-01 ง, FR-APV-03 จ, AC-APV-02, 06, 25 |
| GET /v1/loan-requests/{id}/attachment (ของ 002) | หน้าจอ 004 ลิงก์ไปเปิดไฟล์ใบบันทึกข้อความ ไม่สร้าง endpoint ใหม่ | FR-APV-03 ก, Q-APV-41, FR-BRQ-10 ก |
| GET /v1/approvals/failed-emails | out: Email `failed` ที่ส่งซ้ำครบแล้วของภาควิชา (รวมของ System) | FR-APV-02 ง, AC-APV-14 |
| POST /v1/approvals/emails/{id}/resend | out: สถานะการส่งใหม่ ไม่แตะ `approval_decision` | FR-APV-02 ง, AC-APV-14 |

error (CON-TECH-03) ข้อความตาม FR-APV-01 (ซ) อยู่ใน config:
| code | HTTP | กรณี | สถานะในตารางกลาง | AC |
|---|---|---|---|---|
| `UNAUTHENTICATED` | 401 | ไม่ยืนยันตัวตน | มีแล้ว | AC-APV-05 |
| `FORBIDDEN_ROLE` | 403 | ไม่ใช่เจ้าหน้าที่ | มีแล้ว (004 ลงทะเบียน) | AC-APV-05 |
| `FORBIDDEN_DEPARTMENT` | 403 | อุปกรณ์ภาควิชาอื่น (ข้อความ Q-APV-18) | มีแล้ว | AC-APV-12 |
| `ITEM_NOT_AVAILABLE` | 409 | อุปกรณ์ไม่ "ว่าง" | มีแล้ว | AC-APV-08 |
| `VALIDATION_FAILED` | 422 | ไม่ระบุเหตุผลการปฏิเสธ หรือเหตุผลยาวเกิน 500 ตัวอักษร (field `reason` Q-APV-28) / ไม่ระบุชื่ออาจารย์ หรือยาวเกิน 100 ตัวอักษรตอนอนุมัติค้างคืน (field `advisorName` Q-APV-44) ข้อความต่างกันตามกรณี | มีแล้ว | AC-APV-06, 22 |
| `REQUEST_CANCELLED` | 409 | คำขอ `cancelled` | ลงทะเบียนแล้ว 2569-11-12 (Q-APV-19) | AC-APV-09 |
| `REQUEST_ALREADY_DECIDED` | 409 | พิจารณาไปแล้วก่อนหน้าจอโหลด | ลงทะเบียนแล้ว 2569-11-12 (Q-APV-19) | AC-APV-09 |
| `REQUEST_DECIDED_CONCURRENTLY` | 409 | ถูกพิจารณาหลังหน้าจอโหลด (กดพร้อมกัน) | ลงทะเบียนแล้ว 2569-11-12 (Q-APV-19) | AC-APV-10 |
| `NOT_FIRST_IN_QUEUE` | 409 | ไม่ใช่คำขอแรกของชิ้นนั้น | ลงทะเบียนแล้ว 2569-11-12 (Q-APV-19) | AC-APV-07 |
| `APPROVAL_CUTOFF_PASSED` | 409 | เลย 16:00 น. ของวันที่เริ่มยืม | ลงทะเบียนแล้ว 2569-11-12 (Q-APV-19) | AC-APV-16 |
| `REQUEST_NOT_FOUND` | 404 | ไม่มีคำขอรหัสนี้ (v2: คำขอค้างคืนไม่ได้ 404 แล้ว) | ลงทะเบียนแล้ว 2569-11-12 (Q-APV-25) | AC-APV-21 |
| `UPSTREAM_UNAVAILABLE` | 503 | ฐานข้อมูลไม่ตอบ | มีแล้ว (Q-APV-24) | AC-APV-20 |

- แยก `REQUEST_ALREADY_DECIDED` กับ `REQUEST_DECIDED_CONCURRENTLY` ด้วยการเทียบ `approval_decision.decided_at` กับ `loaded_at` ที่หน้าจอส่งมา (ทีมเลือกเอง ไม่ได้มาจาก spec)
- auth.Middleware ของกลางตอบ `{code, message}` แล้ว (001 แก้ 2569-11-12 ดูข้อ 8) / handler ของ 004 ตอบครบทั้ง `code` และ `message`
- หน้าจอ (ไม่มี mockup ทีมออกแบบเอง ไม่มี UI ID ใน spec):
  - `ApproveRequests` รายการคำขอ ปุ่มอนุมัติ (ยืมปกติเฉพาะแถว `is_first`, ค้างคืนทุกแถว) ปุ่มปฏิเสธพร้อมช่องเหตุผล และแสดงข้อความ error ตามรหัส
    v2: คอลัมน์ประเภทการยืม / แถวค้างคืนมีลิงก์ "เปิดใบบันทึกข้อความ" ไปที่ endpoint ของ 002 และช่องชื่ออาจารย์ก่อนกดอนุมัติ (ตรวจว่าง และ 100 ตัวอักษรก่อนส่ง หลังบ้านตรวจซ้ำ)
  - `FailedEmails` รายการ Email ที่ส่งไม่สำเร็จ ปุ่มส่งซ้ำ
  - เรียกหลังบ้านผ่าน `frontend/src/api/client.js` (เพิ่ม `postJSON` และฟังก์ชันของ 004) test ใช้ client จำลอง

## 5. ตารางตรวจ Constraints
| Constraint ID | ถูกนำไปใช้ที่ไหนใน plan | สถานะ |
|---|---|---|
| CON-TECH-01 | ตาราง `approval_decision`, `approval_email` และ migration ของ 004 ใน PostgreSQL (ข้อ 3) | ใช้แล้ว |
| CON-TECH-02 | route อยู่ใน server เดิมหลัง proxy HTTPS และ Docker ของกลาง (ข้อ 2, 4) | ใช้แล้ว |
| CON-TECH-03 | ตาราง error ข้อ 4 รหัสใหม่ 5 ตัวลงทะเบียนแล้ว 2569-11-12 (Q-APV-19) | ใช้แล้ว |
| CON-DATA-01 | อนุมัติได้เฉพาะอุปกรณ์ `available` และไม่เปลี่ยนสถานะอุปกรณ์ (ข้อ 3, `ITEM_NOT_AVAILABLE`) | ใช้แล้ว |
| DOM-REQ-01 | เรียงตาม submitted_at, `is_first`, `NOT_FIRST_IN_QUEUE`, ปฏิเสธคำขออื่นเมื่อมีคำขอได้อนุมัติ (ข้อ 3, 4) / ค้างคืนไม่มีคิว ใครได้ approved ก่อนได้อุปกรณ์ (FR-APV-03 ค) | ใช้แล้ว |
| DOM-LOAN-01 | เวลาตัด 16:00 น. ของวันที่เริ่มยืม (`APPROVAL_CUTOFF_PASSED`) เฉพาะยืมปกติ / ค้างคืนไม่มีเวลาตัด (FR-APV-03 ง) / กำหนดคืนใน Email = end_at ที่ 002 ตรวจแล้ว | ใช้แล้ว |
| IF-IDP-01 | `auth.Middleware` + บทบาทและภาควิชาจาก Principal (ข้อ 4) | ใช้แล้ว (auth จำลอง รอ Q-SH-03) |
| IF-NOT-01 | `Mailer`, `RecipientResolver`, ตาราง `approval_email`, ส่งซ้ำ (ข้อ 2, 3, 4) | ใช้แล้ว (ตัวส่งจำลอง ที่อยู่ Email รอ Q-APV-04) |

ค่าที่ตั้งได้ (อยู่ที่ `backend/internal/config/config.go` ที่เดียว พร้อมคอมเมนต์อ้าง ID / ข้อความที่หน้าจอแสดงเองอยู่ที่ `frontend/src/config.js`):
| ค่า | มาจาก ID | ชื่อตัวแปรและไฟล์ |
|---|---|---|
| ส่ง Email ภายใน 60 วินาทีหลังบันทึกผล | FR-APV-02 | `ApprovalEmailDeadline` (config.go) |
| ส่งซ้ำอัตโนมัติอีก 3 ครั้งหลังครั้งแรก (รวม 4 ครั้ง) | FR-APV-02 ง, Q-APV-17 | `ApprovalEmailRetryCount` = 3 นับเฉพาะครั้งที่ส่งซ้ำ (config.go) |
| ห่างกันครั้งละ 1 นาที | FR-APV-02 ง | `ApprovalEmailRetryInterval` (config.go) |
| อนุมัติไม่ได้หลัง 16:00 น. ของวันที่เริ่มยืม | FR-APV-01 ค, Q-APV-11, DOM-LOAN-01 | ใช้ `LoanNormalReturnHour` ที่ 002 ประกาศไว้แล้ว (config.go) ไม่ประกาศซ้ำ (AGENTS.md กฎข้อ 9) |
| ปฏิเสธอัตโนมัติ 00:00 น. หลังวันที่ส่ง | FR-APV-01 ฉ, ASM-05 | `ApprovalAutoRejectAt` (config.go) |
| เขตเวลา Asia/Bangkok | glossary | `TimeZone` (config.go มีอยู่แล้ว) |
| สถานที่รับ "คณะวิทยาศาสตร์ ตึก 1 ชั้น 2 ห้องอุปกรณ์ภาควิชาคอมพิวเตอร์" | FR-APV-02 ก | `ApprovalPickupLocation` (config.go) |
| หัวข้อและเนื้อหา Email 3 แบบ | FR-APV-02 ก ข ค | `ApprovalEmailSubjects`, `ApprovalEmailTemplates` (config.go) |
| เหตุผล "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว" | FR-APV-01 จ | `ApprovalSupersededReason` (config.go) |
| เหตุผล "คำขอยังคงค้างอยู่และไม่ได้รับการพิจารณาภายในเวลาที่กำหนด" | FR-APV-01 ฉ | `ApprovalAutoRejectReason` (config.go) |
| แนวทาง "หากยังต้องการยืมอุปกรณ์ สามารถยื่นคำขอใหม่ได้" | FR-APV-02 ค | ใน `ApprovalEmailTemplates` (config.go) |
| ข้อความ error 9 ข้อความ | FR-APV-01 ซ, ASM-07, Q-APV-18 | `ApprovalMessages` (config.go) และ `APPROVAL_MESSAGES` (config.js) |
| เหตุผลการปฏิเสธไม่เกิน 500 ตัวอักษร นับตาม code point และข้อความ "เหตุผลต้องไม่เกิน 500 ตัวอักษร" | FR-APV-01 ง ซ, Q-APV-28 | `ApprovalRejectReasonMaxLength`, `ApprovalReasonTooLongMessage` (config.go) และ `APPROVAL_REASON_MAX_LENGTH`, `APPROVAL_REASON_TOO_LONG` (config.js) |
| ข้อความหน้ารายการว่าง "ไม่มีคำขอที่รออนุมัติ" | FR-APV-01 ก, Q-APV-31 | `APPROVAL_EMPTY_LIST` (config.js) |
| ชื่ออาจารย์บังคับ ไม่เกิน 100 ตัวอักษร นับตาม code point และข้อความ "กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ" / "ชื่ออาจารย์ต้องไม่เกิน 100 ตัวอักษร" | FR-APV-03 ช, FR-APV-01 ซ, Q-APV-44 | `ApprovalAdvisorNameMaxLength`, `ApprovalAdvisorNameRequiredMessage`, `ApprovalAdvisorNameTooLongMessage` (config.go) และค่าเดียวกันใน config.js |
| ค่า loan_type `overnight` และสถานะ `approved_by_officer`, `pending_advisor_approval` | FR-APV-03, data-model.md | ใช้ค่าที่ 002 ประกาศไว้แล้วใน config.go ไม่ประกาศซ้ำ (AGENTS.md กฎข้อ 9) |
| ค่า action `approve`, `reject`, `auto_reject` | RQ-25, DM-07 | `AuditActionApprove`, `AuditActionReject`, `AuditActionAutoReject` (config.go) |

## 6. แผนทดสอบจาก Acceptance Criteria
test หลังบ้านอยู่ใน `backend/internal/approve/*_test.go` ใช้ PostgreSQL จริงผ่าน testcontainers (ข้ามด้วย `-short`)
นาฬิกาจำลองใช้ทดสอบ 16:00, 00:00, 60 วินาที และ 1 นาที / ตัวส่ง Email จำลองบันทึกเวลาและสั่งให้ล้มได้

| AC ID | ชื่อ test | ทดสอบอย่างไร |
|---|---|---|
| AC-APV-01 | TestAC_APV_01 | อนุมัติคำขอแรกของอุปกรณ์ `available` -> `approved`, `equipment_item.status` ยังเป็น `available`, มีแถว `approval_decision` |
| AC-APV-02 | TestAC_APV_02 | ปฏิเสธพร้อมเหตุผล -> `rejected` เหตุผลถูกบันทึก สถานะอุปกรณ์ไม่เปลี่ยน |
| AC-APV-03 | TestAC_APV_03 | คำขอ A, B ของชิ้นเดียวกัน อนุมัติ A -> A `approved`, B `rejected` reason = `ApprovalSupersededReason` actor = เจ้าหน้าที่คนเดียวกัน และมี Email ผลปฏิเสธของ B |
| AC-APV-17 | TestAC_APV_17 | คำขอปกติ A และคำขอค้างคืน C ของชิ้นเดียวกัน อนุมัติ A -> C `rejected` reason = `ApprovalSupersededReason` มี Email ผลปฏิเสธของ C / อีกกรณี: ถึง 00:00 น. หลังวันที่ส่งของ C -> C ยังเป็น `requested` |
| AC-APV-04 | TestAC_APV_04 | หลังอนุมัติ/ปฏิเสธ ตัวส่งจำลองได้ Email ภายใน 60 วินาที หัวข้อและข้อมูลครบตาม FR-APV-02 ก ข (ไม่มีจำนวนอุปกรณ์) |
| AC-APV-05 | TestAC_APV_05 | ไม่มี token -> 401 `UNAUTHENTICATED` / `mock-student-001` -> 403 `FORBIDDEN_ROLE` สถานะคำขอคงเดิม |
| AC-APV-06 | TestAC_APV_06 (+ vitest) | reason ว่างหรือมีแต่ช่องว่าง -> 422 `VALIDATION_FAILED` สถานะคงเดิม / หน้าจอแสดง "กรุณาระบุเหตุผลในการปฏิเสธ" |
| AC-APV-07 | TestAC_APV_07 (+ vitest) | อนุมัติ B ก่อน A -> 409 `NOT_FIRST_IN_QUEUE` ทั้งคู่คงเดิม / หน้าจอแสดงข้อความและไม่มีปุ่มอนุมัติที่แถว B |
| AC-APV-08 | TestAC_APV_08 (+ vitest) | อุปกรณ์ `borrowed` / `unavailable` -> 409 `ITEM_NOT_AVAILABLE` / หน้าจอแสดงข้อความ |
| AC-APV-09 | TestAC_APV_09 (+ vitest) | คำขอ `cancelled` -> `REQUEST_CANCELLED` / พิจารณาแล้ว -> `REQUEST_ALREADY_DECIDED` ผลเดิมคงไว้ / หน้าจอแสดงข้อความตามกรณี |
| AC-APV-10 | TestAC_APV_10 (+ vitest) | 2 goroutine อนุมัติ/ปฏิเสธคำขอเดียวกันพร้อมกัน -> สำเร็จ 1 อีกคน `REQUEST_DECIDED_CONCURRENTLY` มี `approval_decision` 1 แถว Email 1 ฉบับ |
| AC-APV-11 | TestAC_APV_11 (+ vitest) | ข้อมูลผสม (requested ภาควิชา 1 หลายใบ, พิจารณาแล้ว, ภาควิชา 2, ค้างคืน) -> เห็นเฉพาะ requested ของภาควิชาตนเอง เรียงตาม submitted_at พร้อมข้อมูลครบ (v2: ค้างคืนขึ้นด้วยพร้อม loan_type ต้องแก้ test เดิมและแถว TC-APV-11-1) |
| AC-APV-12 | TestAC_APV_12 | เจ้าหน้าที่ภาควิชา 2 อนุมัติ/ปฏิเสธคำขอของอุปกรณ์ภาควิชา 1 -> 403 `FORBIDDEN_DEPARTMENT` message ตาม Q-APV-18 สถานะคงเดิม |
| AC-APV-13 | TestAC_APV_13 | ตั้งนาฬิกาเป็น 00:00 น. หลังวันที่ D -> คำขอ requested ที่ส่งวัน D เป็น `rejected` result `auto_rejected` actor `system` และ Email ปฏิเสธอัตโนมัติภายใน 60 วินาที / คำขอที่ส่งวัน D+1 และคำขอค้างคืนที่ส่งวัน D ไม่ถูกแตะ |
| AC-APV-18 | TestAC_APV_18 | ตั้งนาฬิกาเลย 00:00 น. หลังวันที่ D แล้วจึงเริ่ม service (จำลอง server ไม่ได้ทำงานตอนเที่ยงคืน) -> คำขอปกติที่ส่งวัน D เป็น `rejected` result `auto_rejected` และมี Email ปฏิเสธอัตโนมัติ |
| AC-APV-19 | (ยังไม่มี) | เรียกผ่าน HTTP ธรรมดาหรือ TLS ต่ำกว่า 1.2 แล้วไม่ได้บริการ ทดสอบได้เมื่อทำ HTTPS ตาม Q-APV-38 |
| AC-APV-20 | (ยังไม่มี) | ฐานข้อมูลไม่ตอบ (เช่น Querier ที่คืน error) -> 503 `UPSTREAM_UNAVAILABLE` ข้อความจาก config สถานะคงเดิม / หน้าจอแสดงข้อความ |
| AC-APV-21 | TestApproveHTTP_NotFound, TestTC_APV_21_* | รหัสที่ไม่มีอยู่ -> 404 `REQUEST_NOT_FOUND` "ไม่พบคำขอนี้" (v2: TC-APV-01-12, 02-11 ที่ตรวจว่าค้างคืนได้ 404 ต้องให้ทีมแก้แถวก่อน) |
| AC-APV-22 | TestAC_APV_22 (+ vitest) | อนุมัติค้างคืนพร้อมชื่ออาจารย์ -> `approved`, อุปกรณ์ยัง `available`, `approval_decision.advisor_name` ตรง, `approval_step` 3 แถว / ทางผิด: ชื่อว่าง ยาว 101 ตัว อุปกรณ์ไม่ว่าง |
| AC-APV-23 | TestAC_APV_23 | ค้างคืน X (ส่งก่อน) และ Y ชิ้นเดียวกัน อนุมัติ Y -> Y `approved` X `rejected` เหตุผลปฏิเสธตาม |
| AC-APV-24 | TestAC_APV_24 | ค้างคืน C และปกติ A ชิ้นเดียวกัน อนุมัติ C -> A `rejected` เหตุผลปฏิเสธตาม และมี Email ผลปฏิเสธของ A |
| AC-APV-25 | TestAC_APV_25 | ปฏิเสธค้างคืนพร้อมเหตุผล -> `rejected` ไม่มี advisor_name ตัวส่งจำลองได้ Email ผลปฏิเสธภายใน 60 วินาที |
| AC-APV-26 | TestAC_APV_26 | นาฬิกาจำลองเลย 16:00 น. ของวันที่เริ่มยืมและเลยเที่ยงคืนของวันที่ส่ง อนุมัติค้างคืน -> `approved` / งานเที่ยงคืนไม่แตะคำขอค้างคืน |
| AC-APV-14 | TestAC_APV_14 (+ vitest) | ตัวส่งจำลองล้มทุกครั้ง -> ส่งครั้งแรก + ส่งซ้ำ 3 ครั้งห่าง 1 นาที (ตัวส่งจำลองถูกเรียกรวม 4 ครั้ง นาฬิกาจำลอง) แล้ว `failed` พร้อม last_error และเวลา อยู่ใน failed-emails / กดส่งซ้ำหลังตัวส่งกลับมาปกติ -> `sent` และ `approval_decision` ไม่เปลี่ยน |
| AC-APV-15 | TestAC_APV_15 | หลังอนุมัติ ปฏิเสธ ปฏิเสธตาม และปฏิเสธอัตโนมัติ -> แต่ละครั้งมีแถว `approval_decision` (รหัสคำขอ ผล เหตุผล ผู้ดำเนินการ วันเวลา) และ `audit_log` action ตรงกับกรณี |
| AC-APV-16 | TestAC_APV_16 | ตั้งนาฬิกาเลย 16:00 น. ของวันที่เริ่มยืม -> 409 `APPROVAL_CUTOFF_PASSED` สถานะคงเดิม |

vitest อยู่ที่ `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx` และ `FailedEmails.test.jsx` ตั้งชื่อ test ขึ้นต้นด้วย AC ID

## 7. ลำดับงาน
1. migration ของ 004: ตาราง `approval_decision`, `approval_email` และเพิ่มค่า action ใน `audit_log_action_check` (FR-APV-01 ฌ, FR-APV-02 ง) ไฟล์ `0008_approval` (ค่า action กำหนดใน migration นี้ ไม่เขียนใน data-model.md ทีมตัดสิน 2569-11-12)
2. เพิ่มค่าในตารางข้อ 5 ลง `config.go` / `config.js` (FR-APV-01 ซ, FR-APV-02) (รหัส error 5 ตัวลงทะเบียนใน CON-TECH-03 แล้ว 2569-11-12)
3. repository: อ่านรายการตามภาควิชาและลำดับ, ล็อกคำขอทุกใบ (ทุกประเภทการยืม) ของอุปกรณ์ชิ้นนั้นด้วย `SELECT ... FOR UPDATE` (AC-APV-10, 11, 17)
4. service ตรวจก่อนพิจารณา (บทบาท, ภาควิชา, สถานะคำขอ, คำขอแรก, สถานะอุปกรณ์, 16:00, เหตุผล) แล้ว approve / reject + ปฏิเสธคำขออื่นทุกประเภท + `approval_decision` + `audit_log` + outbox (AC-APV-01, 02, 03, 05 ถึง 10, 12, 15, 16, 17)
5. `Mailer`, `RecipientResolver`, ตัวส่งเบื้องหลังจาก `approval_email` ส่งครั้งแรก + ส่งซ้ำ 3 ครั้ง (AC-APV-04, 14)
6. งานเที่ยงคืนปฏิเสธคำขอยืมปกติที่ค้าง และไล่ปฏิเสธคำขอที่เลยเวลาตอน service เริ่มทำงาน (AC-APV-13, 17, 18)
7. handler และ endpoint failed-emails / resend, ผูก route ใน `backend/cmd/server/main.go` และเริ่ม goroutine (แก้ของกลาง บันทึกใน prompt-log ว่าแก้เพราะ 004) (AC-APV-14)
8. หน้าจอ `ApproveRequests` และ `FailedEmails` + ฟังก์ชันใน `frontend/src/api/client.js` (AC-APV-06 ถึง 11, 14)

v2 (CR-009 คำขอยืมค้างคืน) ทำต่อจากของเดิม ไม่แก้ส่วนยืมปกติที่ผ่าน test แล้ว:
9. migration `0011_approval_overnight`: คอลัมน์ advisor_name และตาราง approval_step + ค่าใน config.go / config.js (FR-APV-03 ข ช)
10. repository: รายการรวมค้างคืนพร้อม loan_type และไฟล์แนบ / บันทึก approval_step (AC-APV-11, 22)
11. service: แยกกฎตาม loan_type ค้างคืนตรวจชื่ออาจารย์และอุปกรณ์ "ว่าง" ไม่ตรวจคิวและเวลาตัด ใช้การปฏิเสธตามเดิม (AC-APV-22 ถึง 26)
12. handler: รับ advisorName ตอบ VALIDATION_FAILED ช่อง advisorName (AC-APV-22)
13. หน้าจอ: คอลัมน์ประเภท ลิงก์ไฟล์แนบ ช่องชื่ออาจารย์ (FR-APV-03 ก ข)
14. แก้ test เดิมที่ยึดว่าค้างคืนไม่ขึ้นรายการหรือได้ 404 (TestAC_APV_11, TestListPending_*, TestApproveHTTP_NotFound ถ้ามีกรณีค้างคืน) / แถว TC ที่ขัด (TC-APV-01-12, 02-11, 11-1) ต้องให้ทีมแก้ใน test-cases.md ก่อน

## 8. สิ่งที่ยังไม่ทำ
คำถามที่เปิดอยู่ใน spec (ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ):
- Q-SH-03 auth service ใช้โปรโตคอลอะไร และส่งบทบาทกับภาควิชามาในรูปแบบใด -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ใช้ auth จำลองไปก่อน ASM-02)
- Q-SH-08 [ปิดสำหรับ prototype ทีมตัดสิน 2569-11-13] ถ้าระบบยืนยันตัวตนเชื่อมต่อไม่ได้ ระบบต้องทำอย่างไร -> คงตอบ 401 ของ auth.Middleware ไม่สร้างอะไรเพิ่ม
- Q-APV-01 [ปิดแล้ว 2569-11-13 CR-009 ทาง ก] รับคำขอค้างคืน UC-14 1 คำขอ = 1 ชิ้น -> ข้อ 3, 4, 7 (v2)
- Q-APV-38 prototype ต้องมี HTTPS และ deploy หลังบ้านด้วย Docker ไหม (ข้อเสนอของ 004: ทำ) -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (AC-APV-19)
- Q-APV-04 [ปิดแล้ว 2569-11-13] ใช้ Email แบบจำลองตลอด prototype (ASM-10) `RecipientResolver` ตัวจำลองใช้ต่อ
- Q-APV-13 [ปิดแล้ว 2569-11-12] ทีมอนุมัติ ID ของ AC-APV-16 แล้ว
- Q-APV-22 [ปิดแล้ว 2569-11-12] ทีมอนุมัติ ID ของ AC-APV-17, 18 แล้ว

ทีมตอบแล้ว 2569-11-12 (Q-APV-14 ถึง 21) และแก้ plan ตามแล้ว ข้อที่ไม่ต้องแก้ plan:
- Q-APV-15 เที่ยงคืนนับจากวันที่ส่ง (เดิม) / Q-APV-16 คำขอแรกที่เลย 16:00 น. ยังขวางคิว (เดิม ใช้ `APPROVAL_CUTOFF_PASSED` และเจ้าหน้าที่ปฏิเสธเอง)
- Q-APV-20 ไม่แจ้ง 003 เรื่องล็อกแถว ความเสี่ยงที่เหลือ: ถ้า 003 ไม่ล็อกแถวเดียวกัน คำขอที่เพิ่งอนุมัติอาจถูกยกเลิกทับได้ (ทีมรับทราบ)

เรื่องที่พบตอนทำ plan v2 (2569-11-13):
- **วิธีเก็บสถานะระหว่างทางของค้างคืน** spec ให้ตัดสินใน plan (FR-APV-03 ข) plan เลือกตาราง `approval_step` แบบ append-only (ทีมเลือกเอง ไม่ได้มาจาก spec) ทางอื่นคือเพิ่ม action ใน audit_log ซึ่งต้องเขียน CHECK กลางใหม่ ทีมตรวจได้
- **endpoint ไฟล์แนบของ 002 ให้เจ้าหน้าที่ทุกภาควิชาเปิดได้** (ไม่ตรวจภาควิชา) 004 ลิงก์เฉพาะคำขอที่แสดงในรายการของภาควิชาตนเอง ไม่แก้ของ 002

เรื่องที่พบตอนทำ plan (004 ไม่ตัดสินเอง):
- **ชื่อเจ้าหน้าที่ไม่มีใน Principal (DM-08, IF-IDP-01)** [ปิดแล้ว ASM-11 ใช้รหัสแทนชื่อ] แต่ spec ต้องการชื่อผู้อนุมัติ / ผู้พิจารณาใน Email (FR-APV-02 ก ข) และชื่อเจ้าหน้าที่ในประวัติ (FR-APV-01 ฌ)
  -> เสนอรวมใน CR เดียวกับ Email (Q-APV-04) ให้ระบบยืนยันตัวตนส่งชื่อและ Email มาด้วย ระหว่างรอ `actor_name` ว่างได้
- **`auth.Middleware` ของกลางตอบแค่ `{code}` ไม่มี `message`** [ปิดแล้ว 2569-11-12] 001 แก้ใน T-37 ของ 001 ตอบ `{code, message}` ผ่าน `platform.WriteError` แล้ว (merge จาก main PR #16)
- **ค่า action `approve`, `reject`, `auto_reject` และเลข migration 0008** [ปิดแล้ว 2569-11-12 ทีมตัดสิน] ยังไม่ต้องเขียนใน data-model.md กำหนดใน migration 0008 ที่ฐานข้อมูลเท่านั้น
- **ลำดับ migration**: ตกลงลำดับ merge ให้ 0004 ถึง 0007 ขึ้นก่อน 0008 เพราะ golang-migrate ไม่รันเลขที่น้อยกว่าเวอร์ชันที่ฐานข้อมูลอยู่แล้ว
