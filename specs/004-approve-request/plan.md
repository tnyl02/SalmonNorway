# Implementation Plan: อนุมัติคำขอยืมปกติ (Approve Request)
Spec: SPEC-APV-004 (specs/004-approve-request/spec.md, Draft v2) | Plan status: Draft v1 | Updated: 2569-11-11
Owner: 660710758 นางสาวณัฐสุดา กัญญาบุญ (@natsuda05)
อ้างอิง: specs/000-shared/{constraints.md v5, data-model.md v9, glossary.md}, migration 0003 ของ 002

## 1. สรุปแนวทาง
- เจ้าหน้าที่ภาควิชาเดียวกับอุปกรณ์เปิดหน้ารายการ เห็นคำขอยืมแบบปกติ `requested` เรียงตามเวลาที่ส่ง แล้วอนุมัติหรือปฏิเสธ (FR-APV-01)
- หลังบ้านเป็น package `backend/internal/approve/` หน้าจอเป็น `frontend/src/pages/ApproveRequest/` ตรวจบทบาท ภาควิชา ลำดับ สถานะ และเวลาที่ฝั่ง server ทั้งหมด
- บันทึกผลใน transaction เดียว ล็อกแถวคำขอทุกใบของอุปกรณ์ชิ้นนั้น ปฏิเสธคำขออื่นของชิ้นเดียวกันในครั้งเดียว และบันทึกประวัติ (FR-APV-01 จ ฌ)
- Email เขียนลงตาราง outbox ใน transaction เดียวกับผล แล้วตัวส่งเบื้องหลังส่งภายใน 60 วินาที ส่งซ้ำ 3 ครั้งห่าง 1 นาที (FR-APV-02)
- งานตั้งเวลาใน process ปฏิเสธคำขอที่ค้างเมื่อถึง 00:00 น. หลังวันที่ส่ง (FR-APV-01 ฉ)

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
| `borrow_request` (DM-06) | id, requester_id, requester_name, student_code, item_id, loan_type, purpose, purpose_other, start_at, end_at, status, submitted_at | อ่าน และเปลี่ยน status `requested` -> `approved` / `rejected` เท่านั้น: FR-APV-01 | 002 (migration 0003) |
| `equipment_item`, `equipment_model` (DM-04, DM-03) | status, asset_code, ชื่อรุ่น, department_id | อ่านอย่างเดียว: อุปกรณ์ "ว่าง" (FR-APV-01 ค) และภาควิชา (FR-APV-01 ก) | 001 (migration 0001) |
| `approval_decision` (ใหม่ เฉพาะ 004) | id, request_id (UNIQUE, FK), result (`approved` / `rejected` / `auto_rejected`), reason, actor_id, actor_name, actor_department_id, decided_at | ประวัติ RQ-25: FR-APV-01 ฌ, ข้อมูลใน Email: FR-APV-02 | 004 |
| `approval_email` (ใหม่ เฉพาะ 004) | id, request_id, kind (`approved` / `rejected` / `auto_rejected`), status (`pending` / `sent` / `failed`), attempts, last_error, last_attempt_at, next_attempt_at, created_at | ส่ง Email และส่งซ้ำ: FR-APV-02 ง | 004 |
| `audit_log` (DM-07) | actor_id, actor_role, action, request_id, asset_code, outcome | บันทึกทุกการพิจารณา (RQ-25) เพิ่มค่า action `approve`, `reject`, `auto_reject` | 001 (migration 0002) |

- ไม่เปลี่ยน `equipment_item.status` (CON-DATA-01, data-model.md แถว UC-12)
- หน้ารายการและการอนุมัติกรองเฉพาะ `loan_type = 'normal'` (ค้างคืนรอ Q-APV-01)
- `approval_decision` ไม่มี UPDATE / DELETE และ request_id UNIQUE จึงพิจารณาซ้ำหรือเปลี่ยนผลไม่ได้ (Out of scope)
- ผลแต่ละแบบ:
  | กรณี | `borrow_request.status` | `approval_decision.result` | reason | actor |
  |---|---|---|---|---|
  | อนุมัติ | `approved` | `approved` | ว่าง | เจ้าหน้าที่ |
  | ปฏิเสธ | `rejected` | `rejected` | ที่เจ้าหน้าที่พิมพ์ | เจ้าหน้าที่ |
  | ถูกปฏิเสธตามเพราะมีคำขออื่นได้อนุมัติ (FR-APV-01 จ) | `rejected` | `rejected` | `ApprovalSupersededReason` | เจ้าหน้าที่ที่กดอนุมัติ |
  | เที่ยงคืน (FR-APV-01 ฉ) | `rejected` | `auto_rejected` | `ApprovalAutoRejectReason` | `system` / `System` + ภาควิชาเจ้าของอุปกรณ์ |
- ไม่มีคอลัมน์จำนวนอุปกรณ์ เพราะ 1 คำขอ = อุปกรณ์รายชิ้น 1 ชิ้น (FR-BRQ-01, ASM-BRQ-02) และไม่มีคอลัมน์ Email (รอ Q-APV-04)
- `actor_name` ว่างได้ เพราะ Principal (DM-08) ยังไม่มีชื่อ ดูข้อ 8
- migration ของ 004 ใช้เลขที่จองในกลุ่ม (0004 ขึ้นไป) และเพิ่มค่า action ด้วยวิธีเดียวกับ 0003 (DROP แล้ว ADD `audit_log_action_check`) ห้ามแก้ migration เดิม

## 4. API / หน้าจอ
ทุก endpoint ขึ้นต้นด้วย `/v1/approvals/` ผ่าน `auth.Middleware` (IF-IDP-01) แล้วตรวจ `Role == staff` (AC-APV-05)
และภาควิชาของเจ้าหน้าที่ = ภาควิชาเจ้าของอุปกรณ์ (ASM-06) ผูกใน `backend/cmd/server/main.go` แบบเดียวกับ `equipment.PathPrefix`

| method path | input / output หลัก | รองรับ |
|---|---|---|
| GET /v1/approvals/requests | out: คำขอ `requested` ยืมปกติของภาควิชา เรียงตาม submitted_at (รหัสคำขอ, ชื่อผู้ขอ, ชื่อรุ่น, รหัสครุภัณฑ์, วัตถุประสงค์, start_at, end_at, submitted_at, `is_first`) และ `loaded_at` | FR-APV-01 ก ข, AC-APV-11 |
| POST /v1/approvals/requests/{id}/approve | in: `loaded_at` / out: ผล และรายการคำขอที่ถูกปฏิเสธตาม | FR-APV-01 ข ค จ ช, AC-APV-01, 03, 07, 08, 09, 10, 16 |
| POST /v1/approvals/requests/{id}/reject | in: `reason`, `loaded_at` / out: ผล | FR-APV-01 ง, AC-APV-02, 06 |
| GET /v1/approvals/failed-emails | out: Email `failed` ที่ส่งซ้ำครบแล้วของภาควิชา (รวมของ System) | FR-APV-02 ง, AC-APV-14 |
| POST /v1/approvals/emails/{id}/resend | out: สถานะการส่งใหม่ ไม่แตะ `approval_decision` | FR-APV-02 ง, AC-APV-14 |

error (CON-TECH-03) ข้อความตาม FR-APV-01 (ซ) อยู่ใน config:
| code | HTTP | กรณี | สถานะในตารางกลาง | AC |
|---|---|---|---|---|
| `UNAUTHENTICATED` | 401 | ไม่ยืนยันตัวตน | มีแล้ว | AC-APV-05 |
| `FORBIDDEN_ROLE` | 403 | ไม่ใช่เจ้าหน้าที่ | มีแล้ว (004 ลงทะเบียน) | AC-APV-05 |
| `FORBIDDEN_DEPARTMENT` | 403 | อุปกรณ์ภาควิชาอื่น | มีแล้ว | AC-APV-12 |
| `ITEM_NOT_AVAILABLE` | 409 | อุปกรณ์ไม่ "ว่าง" | มีแล้ว | AC-APV-08 |
| `VALIDATION_FAILED` | 422 | ไม่ระบุเหตุผลการปฏิเสธ (field `reason`) | มีแล้ว | AC-APV-06 |
| `REQUEST_CANCELLED` | 409 | คำขอ `cancelled` | ต้องลงทะเบียน | AC-APV-09 |
| `REQUEST_ALREADY_DECIDED` | 409 | พิจารณาไปแล้วก่อนหน้าจอโหลด | ต้องลงทะเบียน | AC-APV-09 |
| `REQUEST_DECIDED_CONCURRENTLY` | 409 | ถูกพิจารณาหลังหน้าจอโหลด (กดพร้อมกัน) | ต้องลงทะเบียน | AC-APV-10 |
| `NOT_FIRST_IN_QUEUE` | 409 | ไม่ใช่คำขอแรกของชิ้นนั้น | ต้องลงทะเบียน | AC-APV-07 |
| `APPROVAL_CUTOFF_PASSED` | 409 | เลย 16:00 น. ของวันที่เริ่มยืม | ต้องลงทะเบียน | AC-APV-16 |

- แยก `REQUEST_ALREADY_DECIDED` กับ `REQUEST_DECIDED_CONCURRENTLY` ด้วยการเทียบ `approval_decision.decided_at` กับ `loaded_at` ที่หน้าจอส่งมา (ทีมเลือกเอง ไม่ได้มาจาก spec)
- auth.Middleware ของกลางตอบแค่ `{code}` ไม่มี `message` ดูข้อ 8 / handler ของ 004 ตอบครบทั้ง `code` และ `message`
- หน้าจอ (ไม่มี mockup ทีมออกแบบเอง ไม่มี UI ID ใน spec):
  - `ApproveRequests` รายการคำขอ ปุ่มอนุมัติ (เฉพาะแถว `is_first`) ปุ่มปฏิเสธพร้อมช่องเหตุผล และแสดงข้อความ error ตามรหัส
  - `FailedEmails` รายการ Email ที่ส่งไม่สำเร็จ ปุ่มส่งซ้ำ
  - เรียกหลังบ้านผ่าน `frontend/src/api/client.js` (เพิ่ม `postJSON` และฟังก์ชันของ 004) test ใช้ client จำลอง

## 5. ตารางตรวจ Constraints
| Constraint ID | ถูกนำไปใช้ที่ไหนใน plan | สถานะ |
|---|---|---|
| CON-TECH-01 | ตาราง `approval_decision`, `approval_email` และ migration ของ 004 ใน PostgreSQL (ข้อ 3) | ใช้แล้ว |
| CON-TECH-02 | route อยู่ใน server เดิมหลัง proxy HTTPS และ Docker ของกลาง (ข้อ 2, 4) | ใช้แล้ว |
| CON-TECH-03 | ตาราง error ข้อ 4 รหัสใหม่ 5 ตัวลงทะเบียนก่อนใช้ | ใช้แล้ว (middleware ของกลางยังไม่มี message ดูข้อ 8) |
| CON-DATA-01 | อนุมัติได้เฉพาะอุปกรณ์ `available` และไม่เปลี่ยนสถานะอุปกรณ์ (ข้อ 3, `ITEM_NOT_AVAILABLE`) | ใช้แล้ว |
| DOM-REQ-01 | เรียงตาม submitted_at, `is_first`, `NOT_FIRST_IN_QUEUE`, ปฏิเสธคำขออื่นเมื่อมีคำขอได้อนุมัติ (ข้อ 3, 4) | ใช้แล้ว |
| DOM-LOAN-01 | เวลาตัด 16:00 น. ของวันที่เริ่มยืม (`APPROVAL_CUTOFF_PASSED`) / กำหนดคืนใน Email = end_at ที่ 002 ตรวจแล้ว | ใช้แล้ว |
| IF-IDP-01 | `auth.Middleware` + บทบาทและภาควิชาจาก Principal (ข้อ 4) | ใช้แล้ว (auth จำลอง รอ Q-SH-03) |
| IF-NOT-01 | `Mailer`, `RecipientResolver`, ตาราง `approval_email`, ส่งซ้ำ (ข้อ 2, 3, 4) | ใช้แล้ว (ตัวส่งจำลอง ที่อยู่ Email รอ Q-APV-04) |

ค่าที่ตั้งได้ (อยู่ที่ `backend/internal/config/config.go` ที่เดียว พร้อมคอมเมนต์อ้าง ID / ข้อความที่หน้าจอแสดงเองอยู่ที่ `frontend/src/config.js`):
| ค่า | มาจาก ID | ชื่อตัวแปรและไฟล์ |
|---|---|---|
| ส่ง Email ภายใน 60 วินาทีหลังบันทึกผล | FR-APV-02 | `ApprovalEmailDeadline` (config.go) |
| ส่งซ้ำอัตโนมัติ 3 ครั้ง | FR-APV-02 ง | `ApprovalEmailRetryCount` (config.go) |
| ห่างกันครั้งละ 1 นาที | FR-APV-02 ง | `ApprovalEmailRetryInterval` (config.go) |
| อนุมัติไม่ได้หลัง 16:00 น. ของวันที่เริ่มยืม | FR-APV-01 ค, Q-APV-11, DOM-LOAN-01 | `ApprovalApproveCutoff` (config.go) |
| ปฏิเสธอัตโนมัติ 00:00 น. หลังวันที่ส่ง | FR-APV-01 ฉ, ASM-05 | `ApprovalAutoRejectAt` (config.go) |
| เขตเวลา Asia/Bangkok | glossary | `TimeZone` (config.go มีอยู่แล้ว) |
| สถานที่รับ "คณะวิทยาศาสตร์ ตึก 1 ชั้น 2 ห้องอุปกรณ์ภาควิชาคอมพิวเตอร์" | FR-APV-02 ก | `ApprovalPickupLocation` (config.go) |
| หัวข้อและเนื้อหา Email 3 แบบ | FR-APV-02 ก ข ค | `ApprovalEmailSubjects`, `ApprovalEmailTemplates` (config.go) |
| เหตุผล "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว" | FR-APV-01 จ | `ApprovalSupersededReason` (config.go) |
| เหตุผล "คำขอยังคงค้างอยู่และไม่ได้รับการพิจารณาภายในเวลาที่กำหนด" | FR-APV-01 ฉ | `ApprovalAutoRejectReason` (config.go) |
| แนวทาง "หากยังต้องการยืมอุปกรณ์ สามารถยื่นคำขอใหม่ได้" | FR-APV-02 ค | ใน `ApprovalEmailTemplates` (config.go) |
| ข้อความ error 8 ข้อความ | FR-APV-01 ซ, ASM-07 | `ApprovalMessages` (config.go) และ `APPROVAL_MESSAGES` (config.js) |
| ค่า action `approve`, `reject`, `auto_reject` | RQ-25, DM-07 | `AuditActionApprove`, `AuditActionReject`, `AuditActionAutoReject` (config.go) |

## 6. แผนทดสอบจาก Acceptance Criteria
test หลังบ้านอยู่ใน `backend/internal/approve/*_test.go` ใช้ PostgreSQL จริงผ่าน testcontainers (ข้ามด้วย `-short`)
นาฬิกาจำลองใช้ทดสอบ 16:00, 00:00, 60 วินาที และ 1 นาที / ตัวส่ง Email จำลองบันทึกเวลาและสั่งให้ล้มได้

| AC ID | ชื่อ test | ทดสอบอย่างไร |
|---|---|---|
| AC-APV-01 | TestAC_APV_01 | อนุมัติคำขอแรกของอุปกรณ์ `available` -> `approved`, `equipment_item.status` ยังเป็น `available`, มีแถว `approval_decision` |
| AC-APV-02 | TestAC_APV_02 | ปฏิเสธพร้อมเหตุผล -> `rejected` เหตุผลถูกบันทึก สถานะอุปกรณ์ไม่เปลี่ยน |
| AC-APV-03 | TestAC_APV_03 | คำขอ A, B ของชิ้นเดียวกัน อนุมัติ A -> A `approved`, B `rejected` reason = `ApprovalSupersededReason` actor = เจ้าหน้าที่คนเดียวกัน และมี Email ผลปฏิเสธของ B |
| AC-APV-04 | TestAC_APV_04 | หลังอนุมัติ/ปฏิเสธ ตัวส่งจำลองได้ Email ภายใน 60 วินาที หัวข้อและข้อมูลครบตาม FR-APV-02 ก ข (ไม่มีจำนวนอุปกรณ์) |
| AC-APV-05 | TestAC_APV_05 | ไม่มี token -> 401 `UNAUTHENTICATED` / `mock-student-001` -> 403 `FORBIDDEN_ROLE` สถานะคำขอคงเดิม |
| AC-APV-06 | TestAC_APV_06 (+ vitest) | reason ว่างหรือมีแต่ช่องว่าง -> 422 `VALIDATION_FAILED` สถานะคงเดิม / หน้าจอแสดง "กรุณาระบุเหตุผลในการปฏิเสธ" |
| AC-APV-07 | TestAC_APV_07 (+ vitest) | อนุมัติ B ก่อน A -> 409 `NOT_FIRST_IN_QUEUE` ทั้งคู่คงเดิม / หน้าจอแสดงข้อความและไม่มีปุ่มอนุมัติที่แถว B |
| AC-APV-08 | TestAC_APV_08 (+ vitest) | อุปกรณ์ `borrowed` / `unavailable` -> 409 `ITEM_NOT_AVAILABLE` / หน้าจอแสดงข้อความ |
| AC-APV-09 | TestAC_APV_09 (+ vitest) | คำขอ `cancelled` -> `REQUEST_CANCELLED` / พิจารณาแล้ว -> `REQUEST_ALREADY_DECIDED` ผลเดิมคงไว้ / หน้าจอแสดงข้อความตามกรณี |
| AC-APV-10 | TestAC_APV_10 (+ vitest) | 2 goroutine อนุมัติ/ปฏิเสธคำขอเดียวกันพร้อมกัน -> สำเร็จ 1 อีกคน `REQUEST_DECIDED_CONCURRENTLY` มี `approval_decision` 1 แถว Email 1 ฉบับ |
| AC-APV-11 | TestAC_APV_11 (+ vitest) | ข้อมูลผสม (requested ภาควิชา 1 หลายใบ, พิจารณาแล้ว, ภาควิชา 2, ค้างคืน) -> เห็นเฉพาะ requested ยืมปกติของภาควิชาตนเอง เรียงตาม submitted_at พร้อมข้อมูลครบ |
| AC-APV-12 | TestAC_APV_12 | เจ้าหน้าที่ภาควิชา 2 อนุมัติ/ปฏิเสธคำขอของอุปกรณ์ภาควิชา 1 -> 403 `FORBIDDEN_DEPARTMENT` สถานะคงเดิม |
| AC-APV-13 | TestAC_APV_13 | ตั้งนาฬิกาเป็น 00:00 น. หลังวันที่ D -> คำขอ requested ที่ส่งวัน D เป็น `rejected` result `auto_rejected` actor `system` และ Email ปฏิเสธอัตโนมัติภายใน 60 วินาที / คำขอที่ส่งวัน D+1 ไม่ถูกแตะ |
| AC-APV-14 | TestAC_APV_14 (+ vitest) | ตัวส่งจำลองล้มทุกครั้ง -> ส่งซ้ำ 3 ครั้งห่าง 1 นาที (นาฬิกาจำลอง) แล้ว `failed` พร้อม last_error และเวลา อยู่ใน failed-emails / กดส่งซ้ำหลังตัวส่งกลับมาปกติ -> `sent` และ `approval_decision` ไม่เปลี่ยน |
| AC-APV-15 | TestAC_APV_15 | หลังอนุมัติ ปฏิเสธ ปฏิเสธตาม และปฏิเสธอัตโนมัติ -> แต่ละครั้งมีแถว `approval_decision` (รหัสคำขอ ผล เหตุผล ผู้ดำเนินการ วันเวลา) และ `audit_log` action ตรงกับกรณี |
| AC-APV-16 | TestAC_APV_16 | ตั้งนาฬิกาเลย 16:00 น. ของวันที่เริ่มยืม -> 409 `APPROVAL_CUTOFF_PASSED` สถานะคงเดิม (ID ของ AC นี้รอทีมอนุมัติ Q-APV-13) |

vitest อยู่ที่ `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx` และ `FailedEmails.test.jsx` ตั้งชื่อ test ขึ้นต้นด้วย AC ID

## 7. ลำดับงาน
1. migration ของ 004: ตาราง `approval_decision`, `approval_email` และเพิ่มค่า action ใน `audit_log_action_check` (FR-APV-01 ฌ, FR-APV-02 ง) รอจองเลข migration และแถวใน DM-07
2. ลงทะเบียนรหัส error 5 ตัวใน CON-TECH-03 และเพิ่มค่าในตารางข้อ 5 ลง `config.go` / `config.js` (FR-APV-01 ซ, FR-APV-02)
3. repository: อ่านรายการตามภาควิชาและลำดับ, ล็อกคำขอทุกใบของอุปกรณ์ชิ้นนั้นด้วย `SELECT ... FOR UPDATE` (AC-APV-10, 11)
4. service ตรวจก่อนพิจารณา (บทบาท, ภาควิชา, สถานะคำขอ, คำขอแรก, สถานะอุปกรณ์, 16:00, เหตุผล) แล้ว approve / reject + ปฏิเสธคำขออื่น + `approval_decision` + `audit_log` + outbox (AC-APV-01, 02, 03, 05 ถึง 10, 12, 15, 16)
5. `Mailer`, `RecipientResolver`, ตัวส่งเบื้องหลังจาก `approval_email` ส่งซ้ำ 3 ครั้ง (AC-APV-04, 14)
6. งานเที่ยงคืนปฏิเสธคำขอค้าง (AC-APV-13)
7. handler และ endpoint failed-emails / resend, ผูก route ใน `backend/cmd/server/main.go` และเริ่ม goroutine (แก้ของกลาง บันทึกใน prompt-log ว่าแก้เพราะ 004) (AC-APV-14)
8. หน้าจอ `ApproveRequests` และ `FailedEmails` + ฟังก์ชันใน `frontend/src/api/client.js` (AC-APV-06 ถึง 11, 14)

## 8. สิ่งที่ยังไม่ทำ
คำถามที่เปิดอยู่ใน spec (ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ):
- Q-SH-03 auth service ใช้โปรโตคอลอะไร และส่งบทบาทกับภาควิชามาในรูปแบบใด -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ใช้ auth จำลองไปก่อน ASM-02)
- Q-SH-08 ถ้าระบบยืนยันตัวตนเชื่อมต่อไม่ได้ ระบบต้องทำอย่างไร -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ
- Q-APV-01 รับคำขอยืมแบบค้างคืน -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (plan กรองเฉพาะ `loan_type = 'normal'`)
- Q-APV-04 Email ของผู้ขอมาจากไหน -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (`RecipientResolver` ตัวจริง)
- Q-APV-13 ทีมอนุมัติ ID ของ AC-APV-16 -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (TestAC_APV_16 ตั้งชื่อรอไว้)

เรื่องที่พบตอนทำ plan (004 ไม่ตัดสินเอง):
- **ชื่อเจ้าหน้าที่ไม่มีใน Principal (DM-08, IF-IDP-01)** แต่ spec ต้องการชื่อผู้อนุมัติ / ผู้พิจารณาใน Email (FR-APV-02 ก ข) และชื่อเจ้าหน้าที่ในประวัติ (FR-APV-01 ฌ)
  -> เสนอรวมใน CR เดียวกับ Email (Q-APV-04) ให้ระบบยืนยันตัวตนส่งชื่อและ Email มาด้วย ระหว่างรอ `actor_name` ว่างได้
- **`auth.Middleware` ของกลางตอบแค่ `{code}` ไม่มี `message`** ไม่ครบ CON-TECH-03 -> แจ้งเจ้าของ 001 (ไม่แก้ใน 004)
- **ค่า action `approve`, `reject`, `auto_reject`** ต้องเพิ่มแถวในตาราง DM-07 ของ data-model.md -> คุยกับ 001 / ทีม
- **เลข migration ของ 004** ต้องจองในกลุ่ม (0004 ขึ้นไป) -> คุยในกลุ่ม
