# Implementation Plan: ยกเลิกคำขอยืม (Cancel Borrow Request)
Spec: SPEC-CXL-003 (specs/003-cancel-request/spec.md, Draft v3.9) | Plan status: Draft v1.7 | Updated: 2569-11-12 (ข้อ 5 เพิ่ม NotifyTimeout ตาม rtm F-04)
Owner: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)

## 1. สรุปแนวทาง
นักศึกษาที่ยืนยันตัวตนแล้วเปิดหน้า "คำขอที่รออนุมัติของฉัน" เห็นเฉพาะคำขอของตนที่ยังยกเลิกได้ แล้วยกเลิกพร้อมเหตุผล (FR-CXL-01, 06, 07, 08)
ระบบเปลี่ยนเฉพาะสถานะของคำขอ ไม่แตะสถานะอุปกรณ์ (CON-CXL-01) บันทึก audit_log ทุกครั้งที่พยายามยกเลิก (FR-CXL-05) และแจ้งเจ้าหน้าที่ทาง Email (FR-CXL-09)
แนวทางหลัก: เปลี่ยนสถานะด้วย `UPDATE ... WHERE status IN (สถานะที่ยกเลิกได้)` คำสั่งเดียว ให้ฐานข้อมูลตัดสินจากสถานะ ณ เวลาที่บันทึก (FR-CXL-04)
บังคับสิทธิ์ที่ server ทั้งหมด (FR-CXL-03) Email ส่งหลังบันทึกสำเร็จ ส่งไม่ได้ไม่ทำให้การยกเลิกล้ม (FR-CXL-09)
ใช้โครงกลางที่ 001 สร้างไว้ (Go, PostgreSQL, auth จำลอง) ไม่สร้างโครงแอปใหม่ ใช้ได้จริงเมื่อ 001 merge เข้า main แล้ว

## 2. เทคโนโลยีที่ใช้
ไฟล์คำสั่งของรายวิชาใช้ FastAPI เป็นค่าเริ่มต้น แต่ AGENTS.md กำหนดให้ทีมใช้ Go จึงใช้ตามโครงกลางของ 001

| สิ่งที่เลือก | มาจาก | หมายเหตุ |
|---|---|---|
| PostgreSQL 16 | CON-TECH-01 | ฐานข้อมูลเดียวกับ 001 |
| HTTPS เท่านั้น, Docker Compose | CON-TECH-02, NFR-SEC-01 | ใช้ `deploy/docker-compose.yml` ของกลาง HTTPS ทำโดย 004 T-50 (TLS 1.2 ขึ้นไปในหลังบ้าน `platform.ListenAndServe` ไม่ใช้ reverse proxy backend :8443) 003 ไม่ต้องแก้ |
| Go (net/http + pgx v5) | ทีมเลือกเอง ไม่ได้มาจาก spec | ตาม plan ข้อ 2 ของ 001 และ AGENTS.md |
| golang-migrate (SQL ล้วน) | ทีมเลือกเอง ไม่ได้มาจาก spec | ใช้ `backend/internal/platform/migrate.go` ของกลาง |
| React + Vite | ทีมเลือกเอง ไม่ได้มาจาก spec | โครงหน้าจอของ template รายวิชา |
| auth จำลอง (`backend/internal/auth`) | IF-IDP-01, ทีมเลือกเอง รอ Q-SH-03 | ใช้ `auth.Middleware` และ `auth.FromContext` ของกลาง ไม่แก้ |
| Email จำลอง (`Notifier` เขียนลง log) | IF-NOT-01, ASM-10 | ออกแบบเป็น interface แบบเดียวกับ `auth.Provider` เปลี่ยนเป็นตัวจริงได้โดยไม่แก้ service |
| go test + testcontainers | ทีมเลือกเอง ไม่ได้มาจาก spec | integration test กับ PostgreSQL จริง เหมือน 001 |
| vitest + Testing Library | ทีมเลือกเอง ไม่ได้มาจาก spec | test หน้าจอ |
| k6 | ทีมเลือกเอง ไม่ได้มาจาก spec | load test ของ NFR-PERF-CXL-01 (AC-CXL-17) |

### โครงไฟล์ของ 003 (ทีมเลือกชื่อเอง)
```
backend/internal/cancelrequest/   repository.go, service.go, handler.go, dto.go, notifier.go, *_test.go
backend/db/migrations/            0004_staff_contact.up.sql (จำลอง ASM-10), 0005_audit_log_notify_action.up.sql (ASM-12)
                                  (0003_borrow_request.up.sql เป็นของ 002 ไม่ใช่ของ 003 ASM-11)
backend/db/seed/                  เพิ่มคำขอตัวอย่างทุกสถานะ และ Email เจ้าหน้าที่จำลอง
frontend/src/pages/MyPendingRequests/   MyPendingRequests.jsx, CancelDialog.jsx, __tests__/
frontend/src/api/client.js        เพิ่มฟังก์ชัน listMyPendingRequests, cancelRequest (ไฟล์กลาง)
frontend/src/config.js            ข้อความและค่าที่หน้าจอใช้ (สร้างถ้ายังไม่มี)
```
เลขไฟล์ migration จองในกลุ่มแล้ว 2569-11-11: 003 ได้ 0004 และ 0005

## 3. โมเดลข้อมูล

### borrow_request (DM-06) 002 สร้างใน 0003_borrow_request (ASM-11 ตกลงกับเจ้าของ 002 แล้ว 2569-11-11)
ตารางจริงของ 002 มีคอลัมน์มากกว่านี้ (วันเวลาเริ่ม-สิ้นสุด เหตุผลและไฟล์แนบของการยืมค้างคืน) ข้างล่างเฉพาะส่วนที่ 003 ใช้ คอลัมน์ที่ 003 ใช้มีเครื่องหมาย *
```
borrow_request(
  id            SERIAL PRIMARY KEY,                                   * รหัสคำขอ
  requester_id  TEXT NOT NULL,                                        * ผู้ขอ (DM-08 UserID) ใช้ตรวจ FR-CXL-03
  item_id       INT  NOT NULL REFERENCES equipment_item(id),          * อุปกรณ์ที่ขอ (ชื่อรุ่น รหัสครุภัณฑ์ FR-CXL-06)
  loan_type     TEXT NOT NULL CHECK (loan_type IN ('normal','overnight')),       DM-06
  purpose       TEXT NOT NULL,                                          DM-06 (RQ-18) 002 กำหนดรายละเอียด
  status        TEXT NOT NULL CHECK (status IN (
                  'requested','approved_by_officer','pending_advisor_approval',
                  'approved','rejected','handed_over','returned','cancelled')),  * CON-CXL-02
  submitted_at  TIMESTAMPTZ NOT NULL DEFAULT now(),                   * เวลาที่ส่ง (DM-06, DOM-REQ-01) ใช้เรียง FR-CXL-06
  cancel_reason VARCHAR(200) NULL,                                    * FR-CXL-08 (200 = Q-14, ASM-09) ของ 003 002 ไม่เขียนค่า
  cancelled_at  TIMESTAMPTZ NULL                                      * เวลาที่ยกเลิก (FR-CXL-05, FR-CXL-09) ของ 003
)
ALTER TABLE audit_log ADD CONSTRAINT fk_audit_log_request FOREIGN KEY (request_id) REFERENCES borrow_request (id)   -- 002 เพิ่มใน 0003
CREATE INDEX idx_request_owner_status ON borrow_request (requester_id, status, submitted_at DESC)   -- FR-CXL-06 002 ใส่ใน 0003
```
การแมปชื่อสถานะ (CON-CXL-02): Requested=`requested`, Approved by Officer=`approved_by_officer`,
Pending Advisor Approval=`pending_advisor_approval`, Approved=`approved`, Rejected=`rejected`,
Handed Over=`handed_over`, Returned=`returned`, ยกเลิก=`cancelled`
ไม่มี foreign key หรือ trigger ที่เปลี่ยน `equipment_item.status` (CON-CXL-01, CON-DATA-01, DOM-REQ-01)

### audit_log (DM-07) ใช้ร่วมกับ 001
001 สร้างแล้วใน `backend/db/migrations/0002_audit_log.up.sql` มี `request_id INT NULL` และ action `cancel` (ปิด Q-13)
003 เพิ่ม migration ใหม่ให้ CHECK ของ `action` รับ `notify` (ASM-12) ห้ามแก้ไฟล์ 0002 ตามกติกาในไฟล์นั้น
การแมปกรณี -> action / outcome / request_id / query ตามตารางใน FR-CXL-05 ของ spec (ASM-12)
ถ้า 001 ทำ `backend/internal/audit/writer.go` (T-19 ของ 001) ที่รับ transaction ได้แล้ว ให้ใช้ของเดิม ถ้ายังไม่มีให้เขียนใน repository ของ 003 ภายใน transaction เดียวกัน
(001 ยืนยัน 2569-11-11: writer.go จะรับได้ทั้ง pool และ pgx.Tx คือรับอะไรก็ได้ที่มี Exec และอนุญาตแผนสำรองถ้ารอไม่ไหว)
บันทึกเฉพาะผู้ใช้ที่ยืนยันตัวตนแล้ว (ASM-07) ไม่เก็บเหตุผลการยกเลิกใน audit_log (เก็บที่ borrow_request แล้ว)

### staff_contact (จำลอง ASM-10)
```
staff_contact(department_id INT REFERENCES department(id), email TEXT NOT NULL)
```
ใช้เฉพาะ Notifier จำลองหาผู้รับ Email ของภาควิชาเจ้าของอุปกรณ์ (FR-CXL-09) ระบบจริงรอ Q-15

### ข้อมูลที่ไม่เก็บ
ไม่เก็บชื่อผู้ขอรายอื่นในหน้าของนักศึกษา และ API ไม่ส่งข้อมูลคำขอของผู้อื่นกลับเลย (FR-CXL-03)

## 4. API / หน้าจอ
ทุก endpoint อยู่หลัง `auth.Middleware` (IF-IDP-01, NFR-SEC-01) ใช้ `Principal.UserID` เป็นผู้ขอเสมอ ไม่รับรหัสผู้ใช้จาก request

| Method / path | รองรับ | input / output หลัก |
|---|---|---|
| `GET /api/v1/my/pending-requests` | FR-CXL-06, FR-CXL-03 | out: `{items:[{requestId, modelName, assetCode, submittedAt}]}` เฉพาะ `requester_id = UserID` และสถานะที่ยกเลิกได้ เรียง `submitted_at DESC, id DESC` (ASM-13) ไม่มีรายการ = `items: []` เจ้าหน้าที่ -> 403 `FORBIDDEN_ROLE` (ASM-13, CR-006) |
| `POST /api/v1/borrow-requests/{id}/cancel` | FR-CXL-01..05, 08, 09 | in: `{reason}` / out 200: `{requestId, status:"cancelled"}` |

ตาราง error ของ `POST .../cancel` (รูปแบบ `{code, message}` ตาม CON-TECH-03 แก้ตาม CR-006 2569-11-12 message ของแต่ละรหัสอยู่ใน CR-006)

| กรณี | HTTP | body | อ้างอิง |
|---|---|---|---|
| เหตุผลว่าง มีแต่ช่องว่าง หรือเกิน 200 ตัวอักษร หรือ body อ่านไม่ได้ | 422 | `{code:"VALIDATION_FAILED", message, fields:["reason"]}` | FR-CXL-08, AC-CXL-14, CR-006 |
| ไม่มีคำขอนี้ หรือเป็นคำขอของผู้อื่น | 404 | `{code:"REQUEST_NOT_FOUND", message}` หน้าจอแสดง "ไม่พบคำขอ" | FR-CXL-03, AC-CXL-03 |
| คำขอเป็น `approved` | 409 | `{code:"ALREADY_APPROVED", message, status:"approved"}` หน้าจอแจ้งว่าได้รับการอนุมัติแล้ว | FR-CXL-02, FR-CXL-04, AC-CXL-02, 04 |
| `rejected`, `handed_over`, `returned`, `cancelled` | 409 | `{code:"NOT_CANCELLABLE", message, status}` หน้าจอแสดง "คำขอนี้ยกเลิกไม่ได้" | FR-CXL-02, AC-CXL-11 |
| session หมดอายุ / ไม่มี token | 401 | `{code:"UNAUTHENTICATED"}` (จาก middleware กลาง) | AC-CXL-07 |
| ผู้ใช้บทบาทเจ้าหน้าที่ | 403 | `{code:"FORBIDDEN_ROLE", message}` ใช้ทั้ง GET และ POST | FR-CXL-03, ASM-13 |
| ฐานข้อมูลไม่ตอบสนอง | 503 | `{code:"UPSTREAM_UNAVAILABLE", message}` (ตัด retryable ตาม CR-006 503 ลองใหม่ได้เสมอ) | FR-CXL-01, ASM-06, AC-CXL-13 |
| Email ส่งไม่สำเร็จ | 200 | ยกเลิกสำเร็จตามปกติ | FR-CXL-09, AC-CXL-16 |

### ลำดับการทำงานของการยกเลิก (FR-CXL-04, FR-CXL-05, FR-CXL-09)
1. ตรวจเหตุผล (trim แล้วไม่ว่าง และไม่เกิน `CancelReasonMaxLen`) ไม่ผ่าน -> 422 VALIDATION_FAILED (ยังไม่แตะฐานข้อมูล)
2. เปิด transaction แล้วรัน
   `UPDATE borrow_request SET status='cancelled', cancel_reason=$reason, cancelled_at=now() WHERE id=$id AND requester_id=$uid AND status = ANY($cancellable) RETURNING ...`
   - ได้ 1 แถว = ยกเลิกสำเร็จ (ฐานข้อมูลตัดสินจากสถานะ ณ เวลาที่บันทึก ถ้า 004 อนุมัติไปก่อน คำสั่งนี้จะไม่เจอแถว)
   - ได้ 0 แถว = `SELECT status, requester_id` เพื่อแยกสาเหตุ: ไม่มีหรือไม่ใช่ของตน -> 404, `approved` -> 409 ALREADY_APPROVED, อื่น ๆ -> 409 NOT_CANCELLABLE
3. เขียน audit_log ใน transaction เดียวกัน ทั้งสำเร็จและไม่สำเร็จ (FR-CXL-05) แล้ว commit
4. หลัง commit เรียก `Notifier.NotifyCancelled(ภาควิชาเจ้าของอุปกรณ์, ผู้ยกเลิก, ชื่อรุ่น, รหัสครุภัณฑ์, เวลา, เหตุผล)` (FR-CXL-09)
   ล้มเหลว -> เขียน audit_log action `notify` outcome `error` แล้วตอบ 200 ตามปกติ
5. ฐานข้อมูลล่มที่ขั้น 2 หรือ 3 -> rollback ตอบ 503 สถานะคงเดิม (FR-CXL-01, ASM-06)
   query timeout ใช้ค่าใน config (`CancelDBTimeout`)

ฝั่ง 004 ต้องใช้เงื่อนไขแบบเดียวกัน (`WHERE status IN (...)`) เพื่อไม่ให้อนุมัติคำขอที่ถูกยกเลิกแล้ว (ASM-05) เป็นงานของ 004

### หน้าจอ
| หน้าจอ | รองรับ | แบบ |
|---|---|---|
| `MyPendingRequests` หน้า "คำขอที่รออนุมัติของฉัน" แต่ละแถวแสดงชื่อรุ่น รหัสครุภัณฑ์ เวลาที่ส่ง (Asia/Bangkok) และปุ่มยกเลิก ไม่มีรายการแสดง "ไม่มีคำขอที่รออนุมัติ" | FR-CXL-06 | ไม่มี mockup ทีมออกแบบเอง |
| `CancelDialog` กล่องยืนยันพร้อมช่องเหตุผล ตัวนับ "x/200" ใต้ช่อง (นับตาม code point เช่น `[...text].length`) ปุ่มยืนยันกดไม่ได้เมื่อว่างหรือเกิน 200 ปุ่มไม่ยืนยันปิดกล่องโดยไม่เรียก API และล้างช่องเหตุผล | FR-CXL-07, FR-CXL-08 | ไม่มี mockup ทีมออกแบบเอง |

- ยกเลิกสำเร็จ: แสดง "ยกเลิกคำขอยืม {ชื่อรุ่น} ({รหัสครุภัณฑ์}) เรียบร้อยแล้ว" แล้วเอาแถวออกจากรายการ (FR-CXL-01, FR-CXL-07)
- โหลดรายการไม่สำเร็จ: แสดง "โหลดรายการไม่สำเร็จ กรุณาลองใหม่อีกครั้ง" พร้อมปุ่ม "ลองใหม่" ไม่แสดงข้อความว่าไม่มีคำขอ (FR-CXL-06)
- 403 `FORBIDDEN_ROLE` (เจ้าหน้าที่เปิดหน้านี้ ASM-13): แสดง `message` จากหลังบ้าน ไม่แสดงปุ่ม "ลองใหม่" เพราะลองใหม่ก็ไม่สำเร็จ หน้าจอตัดสินจาก `code` (ทีมตัดสิน 2569-11-12 ทำใน T-26 ทีมเลือกเอง ไม่ได้มาจาก spec)
- 409 / 404: แสดงข้อความตามตาราง error แล้วโหลดรายการใหม่ (FR-CXL-02, FR-CXL-03)
- 503: แสดงข้อผิดพลาดพร้อมปุ่มลองใหม่ แถวยังอยู่ (ASM-06, AC-CXL-13)
- 401: ใช้กลไกกลางของ 001 (T-21: `onUnauthorized` ใน `api/client.js`, `App.jsx`, `MockLogin.jsx`) แสดงหน้าเข้าสู่ระบบจำลองพร้อม `message` ของ 401
  แล้วกลับมาแท็บ `my-pending` กลับมาแล้วไม่ยกเลิกให้อัตโนมัติ เหตุผลที่พิมพ์ค้างไว้หายไปด้วย (AC-CXL-07 ให้กดยกเลิกใหม่)
  ยังไม่มี router จึงจำชื่อแท็บไว้ในหน่วยความจำแทน `returnTo` แบบ path และรับเฉพาะชื่อแท็บที่มีในระบบ (`safePageKey` กัน open redirect แบบ plan 6.5 ของ 001, ทีมตัดสิน 2569-11-12)
- ไม่มี endpoint และหน้าจอของเจ้าหน้าที่ (FR-CXL-10 ไม่ทำใน prototype)

## 5. ตารางตรวจ Constraints

| Constraint ID | ถูกนำไปใช้ที่ไหนใน plan | สถานะ |
|---|---|---|
| CON-TECH-01 | borrow_request, audit_log, staff_contact อยู่ใน PostgreSQL ฐานเดียวกับ 001 (ข้อ 2, 3) | ใช้แล้ว |
| CON-TECH-02 | รันใน Docker Compose ของกลาง ให้บริการผ่าน HTTPS (ข้อ 2) TLS ทำโดย 004 T-50 ตรวจด้วย AC-APV-19 ของ 004 (แก้ 2569-11-13 rtm F-12) | ใช้แล้ว |
| CON-DATA-01 | ไม่มีโค้ดใดของ 003 เขียน `equipment_item.status` ค่าสถานะอุปกรณ์ยังมี 3 ค่า (ข้อ 3, AC-CXL-01, 05) | ใช้แล้ว |
| DOM-REQ-01 | ยกเลิกเฉพาะแถวของ `id` ที่ระบุ คำขออื่นของอุปกรณ์ชิ้นเดียวกันไม่เปลี่ยน อุปกรณ์ยัง "ว่าง" (ข้อ 4 ขั้น 2, AC-CXL-05) | ใช้แล้ว |
| IF-IDP-01 | ทุก endpoint อยู่หลัง `auth.Middleware` ผู้ขอมาจาก `Principal.UserID` (ข้อ 4) prototype ใช้ auth จำลอง รอ Q-SH-03 | ใช้แล้ว |
| IF-NOT-01 | `Notifier` ส่ง Email ถึงเจ้าหน้าที่ภาควิชาเจ้าของอุปกรณ์หลัง commit (ข้อ 4 ขั้น 4) prototype ใช้ตัวจำลอง (ASM-10) | ใช้แล้ว |
| CON-CXL-01 | เปลี่ยนเฉพาะ `borrow_request.status` (ข้อ 3, ข้อ 4 ขั้น 2) | ใช้แล้ว |
| CON-CXL-02 | CHECK ของ `borrow_request.status` 8 ค่า และ `CancellableStatuses` ใน config (ข้อ 3, ข้อ 5) | ใช้แล้ว |

Quality Requirements: NFR-SEC-01 (middleware กลาง + HTTPS), NFR-AVAIL-01 (container `restart: unless-stopped` + `/healthz` ของกลาง), NFR-PERF-CXL-01 (index ข้อ 3 + load test AC-CXL-17)

### ค่าที่ตั้งได้
ทุกค่าอยู่ที่ `backend/internal/config/config.go` (ต่อจากของ 001) และ `frontend/src/config.js` เท่านั้น พร้อมคอมเมนต์อ้าง ID ห้ามฝังในโค้ดที่อื่น

| ค่า | มาจาก ID | ชื่อตัวแปรและไฟล์ |
|---|---|---|
| สถานะที่ยกเลิกได้ `requested`, `approved_by_officer`, `pending_advisor_approval` | CON-CXL-02, FR-CXL-01, ASM-04 | `CancellableStatuses` ใน config.go |
| เหตุผลยาวสูงสุด 200 ตัวอักษร | FR-CXL-08, Q-14 | `CancelReasonMaxLen` ใน config.go และ `CANCEL_REASON_MAX_LEN` ใน config.js |
| เวลาตอบสนอง p95 <= 15 วินาที | NFR-PERF-CXL-01 | `CancelPerfP95` ใน config.go (ใช้เป็น threshold ของ k6 และ test) |
| เขตเวลาแสดงผล Asia/Bangkok | glossary | `DisplayTimezone` ใน config.go และ `DISPLAY_TIMEZONE` ใน config.js (ใช้ร่วมกับ 001 ถ้า 001 เพิ่มไว้แล้วให้ใช้ตัวเดิม) |
| ข้อความ "ไม่มีคำขอที่รออนุมัติ" | FR-CXL-06 | `MSG_NO_PENDING` ใน config.js |
| ข้อความ "คำขอนี้ยกเลิกไม่ได้" | FR-CXL-02 | `MSG_NOT_CANCELLABLE` ใน config.js |
| ข้อความ "ไม่พบคำขอ" | FR-CXL-03 | `MSG_REQUEST_NOT_FOUND` ใน config.js |
| ข้อความ "ยกเลิกคำขอยืม {ชื่อรุ่น} ({รหัสครุภัณฑ์}) เรียบร้อยแล้ว" | FR-CXL-01, ASM-13 | `MSG_CANCEL_SUCCESS` (แม่แบบ) ใน config.js |
| ข้อความ "โหลดรายการไม่สำเร็จ กรุณาลองใหม่อีกครั้ง" และปุ่ม "ลองใหม่" | FR-CXL-06, ASM-13 | `MSG_LOAD_FAILED`, `LABEL_RETRY` ใน config.js |
| ข้อความว่าอนุมัติแล้ว | FR-CXL-02 | `MSG_ALREADY_APPROVED` ใน config.js |
| timeout ของคำสั่งฐานข้อมูลตอนยกเลิก | ASM-06 (ตัวเลขทีมเลือกเอง ไม่ได้มาจาก spec) | `CancelDBTimeout` ใน config.go ต้องน้อยกว่า 15 วินาที (NFR-PERF-CXL-01) |
| เวลารอส่ง Email หลังยกเลิก | ASM-10, FR-CXL-09, NFR-PERF-CXL-01 (5 วินาที ทีมเลือกเอง ไม่ได้มาจาก spec rtm F-04 2569-11-12) | `NotifyTimeout` ใน config.go ส่งไม่ทันนับเป็นส่งไม่สำเร็จ (บันทึก notify error) |

## 6. แผนทดสอบจาก Acceptance Criteria
backend: go test + testcontainers (PostgreSQL จริง) ในไฟล์ `*_test.go` ของ `backend/internal/cancelrequest/`
frontend: vitest ใน `frontend/src/pages/MyPendingRequests/__tests__/` ใช้ client จำลอง

| AC ID | ชื่อ test | ทดสอบอย่างไร |
|---|---|---|
| AC-CXL-01 | `TestAC_CXL_01_CancelPendingNormal` | seed คำขอ `requested` ของ student-001 อุปกรณ์ `available` -> POST cancel -> 200, สถานะ `cancelled`, `equipment_item.status` ยัง `available` |
| AC-CXL-02 | `TestAC_CXL_02_ApprovedNotCancellable` | คำขอ `approved` -> 409 `ALREADY_APPROVED` สถานะยัง `approved` |
| AC-CXL-03 | `TestAC_CXL_03_OthersRequestNotFound` | student-001 ยกเลิกคำขอของ student-002 -> 404 `REQUEST_NOT_FOUND` สถานะไม่เปลี่ยน |
| AC-CXL-04 | `TestAC_CXL_04_ApprovedBeforeConfirm` | ดึงรายการ (เห็นคำขอ) -> UPDATE ในฐานข้อมูลเป็น `approved` (จำลอง 004) -> POST cancel -> 409 `ALREADY_APPROVED` สถานะยัง `approved` |
| AC-CXL-05 | `TestAC_CXL_05_OtherRequestsOnSameItemUnchanged` | 2 คำขอของ student-001 และ student-002 บนอุปกรณ์ชิ้นเดียวกัน -> ยกเลิกของ 001 -> ของ 002 ยัง `requested` อุปกรณ์ยัง `available` |
| AC-CXL-06 | `TestAC_CXL_06_AuditLogAttempts` | ยกเลิกสำเร็จ 1, คำขอ approved 1, คำขอของผู้อื่น 1 -> audit_log 3 แถว actor student-001 outcome `ok` 1 แถว `forbidden` 2 แถว พร้อม request_id และ cause ใน query |
| AC-CXL-07 | `TestAC_CXL_07_ExpiredSessionNoCancel` (Go) + `AC_CXL_07 redirects to login and back` (vitest) | ไม่มี token -> 401 สถานะไม่เปลี่ยน audit_log ไม่มีแถวใหม่ / หน้าจอได้ 401 -> หน้าเข้าสู่ระบบจำลองพร้อมข้อความ ล็อกอินแล้วกลับแท็บเดิม ไม่ยกเลิกซ้ำ (กลไกของ 001 T-21 ดูข้อ 4) |
| AC-CXL-08 | `TestAC_CXL_08_ListOwnPendingSorted` | student-001 มี requested 2 (เวลาต่างกัน) approved 1 rejected 1 cancelled 1 และ student-002 มี requested 1 -> GET ได้ 2 รายการของ 001 เรียงใหม่ไปเก่า มี modelName assetCode submittedAt |
| AC-CXL-09 | `AC_CXL_09 confirm dialog cancel then confirm` (vitest) | กดยกเลิก -> กดไม่ยืนยัน (ไม่เรียก API แถวยังอยู่) -> กดยกเลิก -> พิมพ์เหตุผล -> ยืนยัน -> แถวหายไป |
| AC-CXL-10 | `TestAC_CXL_10_CancelOvernightApprovedByOfficer` | seed คำขอค้างคืน `approved_by_officer` และ `pending_advisor_approval` (ASM-04 ใช้ seed) -> ยกเลิกได้ทั้งคู่ |
| AC-CXL-11 | `TestAC_CXL_11_RejectedAndCancelledNotCancellable` | `rejected` และ `cancelled` -> 409 `NOT_CANCELLABLE` ทั้งคู่ สถานะคงเดิม |
| AC-CXL-12 | `TestAC_CXL_12_EmptyList` (Go) + `AC_CXL_12 shows empty message` (vitest) | ไม่มีคำขอที่ยกเลิกได้ -> `items: []` / หน้าจอแสดง "ไม่มีคำขอที่รออนุมัติ" |
| AC-CXL-13 | `TestAC_CXL_13_DBDownReturns503` (Go) + `AC_CXL_13 shows error with retry` (vitest) | หยุด container ฐานข้อมูลก่อน POST -> 503 `UPSTREAM_UNAVAILABLE` `{code, message}` (CR-006 ไม่มี retryable) / หน้าจอแสดงข้อผิดพลาดและปุ่มลองใหม่ แถวยังอยู่ เมื่อฐานข้อมูลกลับมา สถานะยังเป็น `requested` |
| AC-CXL-14 | `TestAC_CXL_14_ReasonValidation` (Go) + `AC_CXL_14 confirm disabled` (vitest) | เหตุผลว่าง, ช่องว่างล้วน, 201 ตัวอักษร -> 422 `VALIDATION_FAILED` (CR-006) สถานะไม่เปลี่ยน / "ไม่ต้องใช้อุปกรณ์แล้ว" -> 200 และ `cancel_reason` ตรงกัน / หน้าจอปุ่มยืนยันกดไม่ได้ 2 กรณีแรก |
| AC-CXL-15 | `TestAC_CXL_15_NotifyDepartmentStaff` | ใช้ Notifier ปลอมใน test -> ยกเลิกสำเร็จ -> ได้ 1 ข้อความถึง Email เจ้าหน้าที่ภาควิชาของอุปกรณ์ (จาก staff_contact) มีผู้ยกเลิก ชื่อรุ่น รหัสครุภัณฑ์ เวลา และเหตุผล |
| AC-CXL-16 | `TestAC_CXL_16_NotifyFailureStillCancelled` | Notifier ปลอมคืน error -> 200 สถานะ `cancelled` และ audit_log มีแถว action `notify` outcome `error` |
| AC-CXL-17 | `AC_CXL_17` (สคริปต์ k6 `tests/load/cancel_request.js`) | 50 VU เปิดรายการและยกเลิก (seed คำขอพอสำหรับทุก VU) threshold `http_req_duration p(95) < 15000` |

ทดสอบเพิ่ม (ไม่มี AC ตรง): API ไม่ส่งข้อมูลคำขอของผู้อื่นใน GET, ไม่มีโค้ดของ 003 เขียน `equipment_item`
(การปฏิเสธหน้าเดิมที่เป็น URL ภายนอกอยู่ใน `safePageKey` ของ App.jsx ทดสอบโดย 001 ใน AC-EQS-09.test.jsx ไม่ใช่ test ของ 003)

## 7. ลำดับงาน
1. migration `0004_staff_contact` และ test ว่าตาราง `borrow_request` ของ 002 (0003) ตรงกับที่ 003 ใช้ + seed คำขอทุกสถานะ (ทั้งปกติและค้างคืน) ของ student-001/002 และ Email เจ้าหน้าที่จำลอง ต้องรอ 002 push 0003 (ASM-01, ASM-10, ASM-11, CON-CXL-02)
2. เพิ่มค่าใน `config.go` และ `config.js` ตามตารางข้อ 5
3. API รายการคำขอที่รออนุมัติของฉัน (FR-CXL-06, FR-CXL-03 / AC-CXL-08, 12)
4. API ยกเลิก: UPDATE แบบมีเงื่อนไข แยกสาเหตุ 404/409 (FR-CXL-01..04 / AC-CXL-01..05, 10, 11)
5. ตรวจเหตุผลและเก็บ `cancel_reason` (FR-CXL-08 / AC-CXL-14) และกรณีฐานข้อมูลล่ม (ASM-06 / AC-CXL-13)
6. migration เพิ่ม action `notify` แล้วเขียน audit_log ใน transaction เดียวกัน (FR-CXL-05 / AC-CXL-06, 07)
7. `Notifier` จำลองและส่งหลัง commit (FR-CXL-09 / AC-CXL-15, 16)
8. หน้าจอ `MyPendingRequests` + `CancelDialog` (FR-CXL-06, 07, 08 / AC-CXL-07, 09, 12, 13, 14 ฝั่งหน้าจอ)
9. load test k6 (NFR-PERF-CXL-01 / AC-CXL-17)

## 8. สิ่งที่ยังไม่ทำ
คัดลอกจาก Open Questions ใน spec ส่วนที่เกี่ยวข้องกับแต่ละข้อจะยังไม่สร้างจนกว่าจะได้คำตอบ

- ~~Q-13~~ ปิดแล้ว 2569-11-11 (ASM-12) ข้อความเดิม: DM-07 AuditLog มีฟิลด์ "รหัสครุภัณฑ์" แต่ไม่มีฟิลด์สำหรับรหัสคำขอ FR-CXL-05 ต้องบันทึกคำขอที่ยกเลิก
  ต้องเพิ่มฟิลด์ใน data-model.md (ของกลาง ใช้ร่วมกับ 001) หรือเก็บรหัสคำขอในช่องอื่น? -> ทีม (แจ้งเจ้าของ 001)
  ได้คำตอบแล้ว: ขั้นที่ 6, AC-CXL-06, AC-CXL-16 (ส่วน audit) ทำได้ ใช้ action `cancel` และ `notify` ตาม ASM-12
- **Q-15** (ส่วนระบบจริง) รายชื่อ Email ของเจ้าหน้าที่แต่ละภาควิชามาจากไหน? IF-IDP-01 ส่งข้อมูลเฉพาะผู้ใช้ที่ล็อกอิน -> ถามนักเทคโนโลยีสารสนเทศ
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ: Notifier ตัวจริง (prototype ใช้ตัวจำลองตาม ASM-10)
- **Q-16** (จาก FR-CXL-10 ตอบเมื่อนำเข้า scope) เจ้าหน้าที่ยกเลิกจากหน้าไหน ต้องระบุเหตุผลไหม ต้องแจ้งนักศึกษาไหม และยกเลิกได้เฉพาะคำขอของภาควิชาตนเองใช่ไหม? -> ถามเจ้าหน้าที่
  ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ: FR-CXL-10 ทั้งหมด (ไม่ทำใน prototype)
- **Q-SH-03** auth service ใช้โปรโตคอลอะไร และส่ง role + ภาควิชามาในรูปแบบใด -> ใช้ auth จำลองของกลางไปก่อน ไม่เขียน Provider ตัวจริง
- ~~Q-SH-05~~ ปิดแล้วโดย CR-005 ของ 001 (เจ้าหน้าที่เห็นลำดับคำขอ นักศึกษาไม่เห็น) ลำดับและเวลาของคำขอที่มาก่อน แสดงให้ใครเห็น -> หน้า "คำขอที่รออนุมัติของฉัน" ไม่แสดงลำดับคิวจนกว่าจะได้คำตอบ
- ~~Q-SH-06~~ ปิดแล้ว: คำขออื่นของชิ้นเดียวกันถูก 004 ปฏิเสธอัตโนมัติเป็น `rejected` 003 ไม่เปลี่ยนคำขออื่นเอง และตอบ `NOT_CANCELLABLE` ตาม AC-CXL-11
- **Q-SH-08** ระบบยืนยันตัวตนเชื่อมต่อไม่ได้ ต้องทำอย่างไร -> ใช้พฤติกรรม 401 ของ middleware กลางไปก่อน

### Dependency ที่ต้องเคลียร์ก่อน /implement
- โครงกลาง (`backend/`, `frontend/`) ยังอยู่ใน `origin/001-equipment-status` ยังไม่ merge เข้า main
- ~~migration `audit_log` ยังไม่มี~~ มีแล้วใน main (0002) โครงกลางก็ merge เข้า main แล้ว (2569-11-10)
- เลข migration: data-model.md จอง `0003_borrow_request` ให้ 002 และ `0004` ขึ้นไปให้ 003 เป็นต้นไป จองในกลุ่มแล้ว 2569-11-11: 003 ได้ `0004` และ `0005`
- ~~ชื่อคอลัมน์และค่าสถานะของ `borrow_request` ต้องตกลงกับเจ้าของ 002~~ ตกลงแล้ว 2569-11-11 (ASM-09, ASM-11) **ยังรอ 002 push `0003_borrow_request` เข้า main**
- ~~index ใครเพิ่ม~~ 002 ใส่ `idx_request_owner_status` ใน 0003 แล้ว (ตกลง 2569-11-11)
- ลำดับ push: 0004 และ 0005 ของ 003 ห้ามเข้า main ก่อน 0003 ของ 002 (ถ้าฐานข้อมูลรันถึงเลข 5 แล้ว golang-migrate จะไม่ย้อนรัน 0003 ที่มาทีหลัง)
- 0005 เขียนแบบ `ALTER TABLE audit_log DROP CONSTRAINT audit_log_action_check;` แล้ว `ADD CONSTRAINT audit_log_action_check CHECK (action IN ('list','search','detail','cancel','request','notify'))` (ASM-12)
- data-model.md: รอ 002 push v5 ถึง v8 แล้วนำแถว notify ของ 003 ไปวางบน v8 เป็น v9 (ตอน merge main เข้า branch 003 ให้แก้ conflict ตามนี้)
- ~~001 รับ action notify ไหม~~ รับแล้ว 2569-11-11 (ASM-12)
- writer.go ของ 001 (T-19) ยังไม่มี ยังไม่ได้วันที่เสร็จ ถ้าไม่เสร็จก่อน 003 ถึง T-31 ใช้แผนสำรอง (เขียน audit ใน repository ของ 003)
- ถ้า merge migration ไม่ตรงลำดับเลข ทุกคนต้องล้างฐานข้อมูลในเครื่องด้วย `docker compose -f deploy/docker-compose.yml down -v` (001 แจ้ง)
- **ยังไม่ตัดสิน: การล้างข้อมูลระหว่าง test** ตัวช่วย test `tcServer` ใช้ `DELETE FROM borrow_request` แต่เมื่อ 0003 เพิ่ม FK จาก audit_log
  และ audit_log เป็น append-only (trigger ใน 0002 ห้าม DELETE, UPDATE, TRUNCATE) ลบคำขอที่มีแถว audit อ้างอยู่ไม่ได้ และลบแถว audit ก็ไม่ได้
  ทางเลือก: (ก) ไม่ลบข้อมูลระหว่าง test ให้แต่ละ test ใช้คำขอที่สร้างใหม่และตรวจด้วย requestId ของตัวเอง
  (ข) ตอนล้างใน test ตั้ง `SET session_replication_role = replica` ปิด trigger และ FK ชั่วคราว (ทำได้เพราะ testcontainers ใช้ superuser)
  (ค) สร้างฐานข้อมูลใหม่ต่อ test (ช้า)
  **ทีมเลือก (ข) 2569-11-12** ใช้ `SET LOCAL session_replication_role = replica` ใน transaction ของการล้างเท่านั้น
  กฎกลับมาเองเมื่อ transaction จบ (ห้ามใช้ SET ธรรมดาเพราะค้างใน connection ของ pool) ใช้กับฐานข้อมูล test ของ testcontainers เท่านั้น
- writer.go ของ 001: ทีมเลือกรอ 001 ทำก่อน (2569-11-12) T-31, T-32 ยังไม่เริ่ม
- 004 ต้องใช้เงื่อนไขสถานะตอนอนุมัติ เพื่อไม่ให้อนุมัติคำขอที่ถูกยกเลิกแล้ว (ASM-05) เจ้าของ 004 รับกฎแล้ว 2569-11-11 และหน้าอนุมัติซ่อนคำขอ cancelled
