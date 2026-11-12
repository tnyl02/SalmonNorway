# RTM: บันทึกการรับคืนอุปกรณ์ (Return)
อ้างอิง: spec.md Draft v2.2 | tasks.md (T-01 ถึง T-23, T-29 เสร็จ / T-24 ถึง T-27 รอคำตอบ) | test-cases.md (60 แถว: ใช้ได้ 59, รอ Q-SH-08 1)
สร้างด้วย /verify เมื่อ 2569-11-12 | test: 584 ผ่าน 0 ไม่ผ่าน (backend `go test ./...` 447 รวม subtest ทุกฟีเจอร์ ของ 006 `returns` 67 / frontend `npm test` 137 ของ 006 30)

## 1. ตามรอยไปข้างหน้า (requirement ไป โค้ด ไป test)
| ID | AC | task | โค้ด (ไฟล์: ฟังก์ชัน) | test (ผล) | สถานะ |
|---|---|---|---|---|---|
| FR-RET-01 | AC-RET-03, 08, 10, 11 | T-05, T-06, T-14, T-15, T-18 ถึง T-22 | `repository.go`: `FindOpenLoan`, `ListOpenLoans`, `checkItemDepartment`, `RecordReturn` (`status = 'borrowed'`, `returned_at IS NULL`) / `handler.go`: `lookup`, `openLoans`, `staffOnly` / `ReturnPage.jsx`, `OpenLoansPage.jsx`, `ReturnsHome.jsx`, `App.jsx` แถว `returns` | `TestAC_RET_03/08/10/11_*`, `TestTC_RET_03_*`, `TestTC_RET_08_1/2/3/5`, `TestTC_RET_10_1/2/3/6`, `TestTC_RET_11_1/2/4`, `TestReturns_ConcurrentReturn_OneWins` (ผ่าน) / vitest AC-RET-03, 10, 11, TC-RET-03-3, 08-4, 10-4, 11-3 (ผ่าน) / TC-RET-10-5 ให้คนลอง | ช่องโหว่ (F-03) |
| FR-RET-02 | AC-RET-04, AC-RET-09 | T-04, T-07, T-09, T-10 | `repository.go`: `RecordReturn` (`returned_at = now()`) / `service.go`: `RecordReturn` (สภาพจาก `config.ReturnConditionStatus`) / `handler.go`: `recordReturn` ไม่อ่านเวลาจาก form / `ReturnPage.jsx`: ตัวเลือกจาก `RETURN_CONDITIONS` / `LoanDetailPage.jsx` แสดงเวลาตาม `TIME_ZONE` | `TestAC_RET_04_*`, `TestAC_RET_09_*`, `TestTC_RET_04_1/4/5/6/7`, `TestTC_RET_09_1/2/5`, `TestMigration006_LoanReturnColumns`, `TestConfig006_*` (ผ่าน) / vitest AC-RET-04 x2, TC-RET-04-2, 04-3, 09-3, 09-4 (ผ่าน) | ครบ |
| FR-RET-03 | AC-RET-01 | T-04, T-06 | `service.go`: `RecordReturn` -> `config.ReturnConditionStatus["normal"] = "available"` / `repository.go`: `RecordReturn` UPDATE `equipment_item` | `TestAC_RET_01_*`, `TestTC_RET_01_1/2/3/4` (ผ่าน) / vitest AC-RET-01 (ผ่าน) / TC-RET-01-5 ให้คนลอง | ครบ |
| FR-RET-04 | AC-RET-02 | T-08 | `config.ReturnConditionStatus["damaged"] = "unavailable"` / `RecordReturn` UPDATE สถานะครั้งเดียว | `TestAC_RET_02_*`, `TestTC_RET_02_1/2/3/4` (ผ่าน ตรวจประวัติสถานะด้วย trigger) | ครบ |
| FR-RET-05 | AC-RET-05, AC-RET-12 (ยกเลิก) | - | ไม่มี (ตรวจแล้ว: ไม่มีตาราง โค้ด หรือช่องแนบรูป) | vitest "ไม่มีช่องแนบรูป" (ผ่าน) | ยกเลิก 2569-11-11 (ไม่นับ) |
| FR-RET-06 | AC-RET-06, AC-RET-13 | T-11 ถึง T-13, T-16, T-17 | `RecordReturn` (`returned_late = now() > due_at`), `GetLoan`, `ListOpenLoans` (`isOverdue`) / `LoanDetailPage.jsx` ป้าย "คืนเกินกำหนด", `OpenLoansPage.jsx` ป้าย "เลยกำหนดคืน" | `TestAC_RET_06_*`, `TestAC_RET_13_*`, `TestReturns_LoanDetail`, `TestTC_RET_06_1..6`, `TestTC_RET_13_1/4/6` (ผ่าน) / vitest AC-RET-13 x2, TC-RET-13-2, 13-3, 13-5 (ผ่าน) / ผู้ยืมเห็นที่ไหน รอ Q-03 (T-25) | ครบ (ส่วนของผู้ยืม รอ Q-03) |
| FR-RET-07 | AC-RET-07 | T-01, T-23 | `audit.go`: `insertAudit`, `WriteAudit`, `outcomeFor` / `RecordReturn` เขียน audit ใน transaction / `handler.go`: `respond`, `staffOnly` | `TestAC_RET_07_*`, `TestTC_RET_07_1..5`, `TestMigration006_AuditLogReturnActions` (ผ่าน) | ช่องโหว่ (F-01, F-02) |
| NFR-DATA-01 | AC-EQS-05 (ของ 001) | ไม่มีใน 006 (ของ 001) | `RecordReturn` เขียนทันทีใน transaction / `writeJSON` ตั้ง `Cache-Control: no-store` | `TestAC_EQS_05_ReturnedBecomesAvailable` (001, ผ่าน) / vitest `AC-EQS-05.test.jsx` (ผ่าน) / TC-RET-01-5 ให้คนลอง | ครบ |
| NFR-SEC-01 | ไม่มี AC ใน 006 (ส่วน auth: AC-RET-08) | T-03, T-22 | ทุก route ผ่าน `auth.Middleware` + `staffOnly` / TLS: `deploy/docker-compose.yml` ไม่มี HTTPS | `TestReturns_Skeleton_RoutedThroughAuth`, `TestTC_RET_01_4`, `TestTC_RET_08_5` (ผ่าน) / TLS ไม่มี test | ยังไม่ถึง (TLS เป็น deploy ของกลาง ยังไม่ทำ) |
| NFR-AVAIL-01 | ไม่มี AC | ไม่มี | ไม่มี | ไม่มี | ช่องโหว่ (F-04) |
| CON-TECH-01 | ไม่มี AC ตรง | T-01 | migration `0010_return.up.sql` / `.down.sql`, `pgxpool` ผ่าน `platform.OpenDB` | `TestMigration006_*` บน PostgreSQL 16 จริง (ผ่าน) | ครบ |
| CON-TECH-02 | ไม่มี AC ตรง | T-03 (ใช้ deploy ของ 001) | ใช้ `deploy/docker-compose.yml` เดิม ยังไม่มี HTTPS / reverse proxy | ไม่มี | ยังไม่ถึง (deploy ของกลาง) |
| CON-TECH-03 | ไม่มี AC ตรง (Constraint ใช้ร่วม) | T-29 | `handler.go`: `errorBody`, `unavailableBody`, `errorResponse` / `config.ReturnMessages` / รหัสลงทะเบียนใน constraints.md v5.4 (Q-10 ปิด) | `TestReturns_ErrorBodyFollowsConTech03`, `isErrorOnly` ใน AC-RET-08, 11, TC-RET-02-4, 03-6, 08-1, 08-2, 11-1, 11-2 (ผ่าน) | ช่องโหว่ (F-05) |
| CON-DATA-01 | AC-RET-01, 02, 03 | T-01, T-04, T-08, T-18 | เขียนเฉพาะ `available` / `unavailable` จาก `borrowed` (CHECK เดิมของ 0001) | `TestAC_RET_01/02/03_*`, `TestTC_RET_01_*`, `TestTC_RET_02_*`, `TestTC_RET_03_*` (ผ่าน) | ครบ |
| CON-DATA-03 | AC-RET-01, 03 | T-04, T-05 | หาการยืมด้วย `asset_code` (UNIQUE) ใน `FindOpenLoan`, `RecordReturn`, `checkItemDepartment` | `TestReturns_Lookup_FoundOpenLoan`, `TestAC_RET_03_*`, `TestTC_RET_03_5` (ผ่าน) | ครบ |
| DOM-AUDIT-01 | AC-RET-07 | T-01, T-23 | ดู FR-RET-07: action `return_list`, `return_lookup` (เก็บ query), `return_detail`, `return` | `TestAC_RET_07_*`, `TestTC_RET_07_1..5` (ผ่าน) | ช่องโหว่ (F-01, F-02) |
| DOM-LOAN-01 | AC-RET-06 | T-13 | `returned_late` / `isOverdue` เทียบ `due_at` ที่ 005 บันทึก ไม่คำนวณ 16:00 เอง | `TestAC_RET_06_*`, `TestTC_RET_06_1..6` (ผ่าน) | ครบ |
| DOM-REQ-01 | ไม่มี AC | ไม่มี (plan ข้อ 5: ยังไม่ได้ใช้) | ไม่มี: 006 ไม่แตะ `borrow_request` (FR-RET-08 ย้อนออก 2569-11-12) | ไม่มี | รอ Q-08 |
| IF-IDP-01 | AC-RET-08 | T-03, T-14, T-20, T-22 (T-27 รอ Q-SH-08) | `auth.Middleware`, `auth.FromContext` ให้บทบาทและภาควิชา (`principal`, `staffOnly`, ภาควิชาใน repository) | `TestReturns_Skeleton_RoutedThroughAuth`, `TestAC_RET_08/10/11_*`, `TestTC_RET_08_1/2/3/5`, `TestTC_RET_10_*`, `TestTC_RET_11_*` (ผ่าน) / TC-RET-08-6 ยังไม่เขียน (รอ Q-SH-08) | ครบ (กรณีระบบยืนยันตัวตนล่ม รอ Q-SH-08) |

สรุปสถานะ 19 แถว: ครบ 10 / ยังไม่ถึง 2 / รอ 1 / ช่องโหว่ 5 / ยกเลิก 1 (FR-RET-05)

## 2. ตามรอยย้อนกลับ (โค้ด ไป requirement)
| โค้ด (ไฟล์: ฟังก์ชัน หรือ endpoint) | อ้าง ID | ตรงกับข้อความใน spec ไหม | หมายเหตุ |
|---|---|---|---|
| `POST /v1/returns` (`handler.go: recordReturn`, `service.go: RecordReturn`, `repository.go: RecordReturn`) | FR-RET-02, 03, 04, 06, 07 | ตรง | ไม่รับเวลาจาก client (TC-RET-09-1, 09-2), สภาพ 2 ค่าตรงตัว (TC-RET-04-6), สถานะตาม config, returned_late, audit ใน transaction |
| `GET /v1/returns/lookup` (`lookup`, `FindOpenLoan`) | FR-RET-01(ก) | ตรง | คืน `borrowerId` ให้เจ้าหน้าที่ภาควิชาเดียวกันเท่านั้น |
| `GET /v1/returns/open-loans` (`openLoans`, `ListOpenLoans`) | FR-RET-01(ข), FR-RET-06 | ตรง | ไม่แบ่งหน้า ตาม plan Q-P1 / เรียงตาม `due_at` ตาม plan ข้อ 4 (TC-RET-10-1) |
| `GET /v1/returns/loans/{loanId}` (`loanDetail`, `GetLoan`) | FR-RET-06 | ตรง | `LOAN_NOT_FOUND` อยู่ในตาราง error ของ plan ข้อ 4 และลงทะเบียนใน CON-TECH-03 แล้ว |
| `handler.go: staffOnly` | FR-RET-01, AC-RET-08 | ตรง | ผู้ยืมเองก็คืนไม่ได้ (TC-RET-08-3, UC-04 อยู่ใน Out of scope) |
| `repository.go: checkItemDepartment` | FR-RET-01, ASM-03 | ตรง | ตรวจภาควิชาก่อนการยืมค้าง อุปกรณ์ว่างของภาควิชาอื่นได้ 403 ตามที่ทีมตัดสินใน TC-RET-03-6 |
| `handler.go: errorBody`, `unavailableBody` | CON-TECH-03 | ตรงบางส่วน | 503 แนบ `retryable: true` ซึ่ง CON-TECH-03 ไม่ได้กำหนด (F-05) |
| `handler.go: respond` / `audit.go: outcomeFor` | DOM-AUDIT-01 | ตรงบางส่วน | การจับคู่ outcome และ fail closed เป็นการตัดสินใจของ AI (F-01, F-02) |
| `handler.go: principal` (ตอบ 401 เมื่อไม่มีผู้ใช้ใน context) | IF-IDP-01 | ตรง | เป็นทางสำรอง ปกติ `auth.Middleware` ตอบ 401 ก่อนถึงตรงนี้ |
| `audit.go` (package returns: `AuditEntry`, `insertAudit`) | DOM-AUDIT-01, FR-RET-07 | ตรง | ไม่ใช้ package `audit` กลางของ 001 (ทีมตัดสิน 2569-11-12 ไม่ทำ ข้อ ค.) ใช้คอลัมน์ DM-07 เดิม ไม่ใส่ `request_id` |
| `backend/db/migrations/0010_return.up.sql` / `.down.sql` | FR-RET-02, 06, 07 | ตรง | CHECK ของ `audit_log` คงค่าของทุกฟีเจอร์ / CHECK ของ `loan` บังคับทางเดียว ไม่กระทบ 001, 002 |
| `backend/db/seed/seed_006_return.sql` | (ข้อมูลตัวอย่าง จุดตรวจ AC-RET-10, 11) | ตรง | ไม่ใช้ id ตายตัว |
| `backend/internal/config/config.go` (บล็อก 006: `Timezone`, `ReturnConditionStatus`, `ReturnMessages`) | glossary, FR-RET-02 ถึง 04, CON-TECH-03 | ตรง | ไม่พบตัวเลขจาก spec ฝังในโค้ด 006 |
| `frontend/src/config.js` (บล็อก 006: `STATUS_LABELS`, `RETURN_CONDITIONS`) + `TIME_ZONE` กลาง | glossary, CON-DATA-01, FR-RET-02 | ตรง | |
| `frontend/src/api/client.js`: `returnsRequest`, `api.returns.*` | FR-RET-01, 02, 06, IF-IDP-01 | ตรง | ใช้ `authHeaders` และ session ของ 001 / `record` ส่งแค่ `assetCode`, `condition` (TC-RET-09-4) |
| `ReturnPage.jsx` | FR-RET-01(ก), 02, 03, 04, AC-RET-03, 04, 11 | ตรง | error อื่นแสดงเป็น "ค้นหาไม่สำเร็จ (code)" ไม่ใช้ `message` จากหลังบ้าน (ทีมตัดสินไม่ทำข้อ ง. และ TC-RET-08-4) |
| `OpenLoansPage.jsx` | FR-RET-01(ข), FR-RET-06 | ตรง | |
| `LoanDetailPage.jsx` | FR-RET-06 | ตรง | |
| `ReturnsHome.jsx` + `App.jsx` แถว `returns` (เลือกแถวแล้วเปิด ReturnPage ใหม่ด้วย key) | FR-RET-01(ข) | ตรง | ไม่มี test อัตโนมัติ มีแค่ TC-RET-10-5 ให้คนลอง (F-03) |
| `cmd/server/main.go` ผูก route ของ 006 กับ `platform.NewRouter` ของ 001 | IF-IDP-01 | ตรง | |
| `specs/006-return/plan.md` ข้อ 6 แถว AC-RET-01 | AC-RET-01 | ไม่ตรง | ยังเขียน "+ 1 JPG" ของ FR-RET-05 ที่ยกเลิกแล้ว (F-06) |

ตรวจ Out of scope ทีละบรรทัด: ไม่พบโค้ดของ นักศึกษาแจ้งคืนเอง (UC-04) / แนบรูป / ปิดการยืมแบบสูญหาย / ระงับคำขอค้าง / ส่งมอบ / ซ่อมแล้วกลับเป็นว่าง / บทลงโทษ / รายงานประวัติรายผู้ยืม / การแจ้งเตือน
ตรวจ Open Questions: ไม่พบโค้ดของ Q-01 (สูญหาย), Q-03 (หน้าผู้ยืม), Q-04 (รายละเอียดความเสียหาย), Q-08 (ระงับคำขอ), Q-09 (สถานะก่อน/หลังใน audit), Q-SH-08 / ไม่มีโค้ดเปลี่ยน `borrow_request.status` (FR-RET-08 ย้อนออกแล้ว)
ตรวจข้อมูลที่ต้องปกป้อง: ไม่มีคอลัมน์หรือฟิลด์เลขบัตรหรือชื่อผู้ยืม log มีแค่ข้อความ error ไม่มีข้อมูลผู้ใช้ body ของ error มีแค่ code และ message (TC-RET-02-4, 03-6, 11-1, 11-2)

## 3. ข้อค้นพบ
ชนิด: AC ไม่มี test / test อ่อน / โค้ดไม่มี FR / FR ไม่มี AC / เดา Q-xx / ละเมิด Constraint / ตัวเลขไม่ตรง spec / อ้าง ID ผิดเรื่อง
ทีมตัดสิน: แก้โค้ด / แก้ spec / เพิ่ม Q-xx / ไม่ใช่ปัญหา (พร้อมเหตุผล 1 บรรทัด)

| F-ID | ชนิด | อยู่ที่ | ขัดกับ | รายละเอียด | ทีมตัดสิน |
|---|---|---|---|---|---|
| F-01 | เดา Q-xx (ไม่มี Q รองรับ) | `backend/internal/returns/audit.go: outcomeFor` | DOM-AUDIT-01, DM-07 (outcome มี 3 ค่า) | AI จับคู่ outcome เอง: 2xx = ok, 403 = forbidden, อื่น ๆ (409 `NO_OPEN_LOAN`, 422 `CONDITION_REQUIRED`, 404, 503) = error ทำให้ "บันทึกไม่สำเร็จเพราะไม่มีการยืมค้าง" ปนกับ "ระบบขัดข้อง" ในประวัติ spec/plan/DM-07 ไม่ได้บอกความหมายของแต่ละค่า เสนอ: เพิ่ม Q ให้ทีมกำหนดความหมายของ outcome ใน DM-07 (อาจต้องเพิ่มค่า เช่น rejected ผ่านของกลาง) / รอบ 2569-11-12: TC-RET-07-4 ที่ทีมยืนยันแล้วกำหนดให้ 409 = `error` ตรงกับโค้ด แต่คอมเมนต์ใน `outcomeFor` เขียนว่า "ทีมยืนยัน ดู prompt-log T-23" ขณะที่ prompt-log T-23 เขียนว่า "ทีมต้องยืนยัน" | |
| F-02 | เดา Q-xx (Q-P4 ของ 001) | `backend/internal/returns/handler.go: respond` | DOM-AUDIT-01 | เขียน audit ไม่สำเร็จแล้วตอบ 503 และไม่ส่งข้อมูล (fail closed) ใช้แนวจาก plan 6.4 ของ 001 ซึ่งยังเป็น Q-P4 ที่ไม่มีคำตอบ ไม่มี test กรณี audit เขียนไม่ได้ (การอ่าน) ส่วนการบันทึกคืน audit อยู่ใน transaction จึง rollback ทั้งก้อนอยู่แล้ว | |
| F-03 | test อ่อน | `frontend/src/pages/Return/ReturnsHome.jsx` (onSelect -> ReturnPage key/initialAssetCode) (เดิมอยู่ใน `App.jsx` ย้ายมาตอนรวม main) | AC-RET-10 Then "เลือกรายการเพื่อบันทึกการรับคืนได้" | test ตรวจแยกส่วน: OpenLoansPage เรียก `onSelect` (TC-RET-10-4) และ ReturnPage รับ `initialAssetCode` แต่ไม่มี test อัตโนมัติที่ตรวจการต่อกันใน ReturnsHome ถ้าลบบรรทัด `setSelected` ไม่มี test ตัวไหนไม่ผ่าน (ReturnsHome.test.jsx ถูกย้อนออก 2569-11-12) ตอนนี้มีแค่ TC-RET-10-5 ที่ให้คนลอง เสนอ: เพิ่ม vitest ของ ReturnsHome หรือยอมรับ TC-RET-10-5 | |
| F-04 | FR ไม่มี AC | `specs/006-return/spec.md` Quality Requirements | NFR-AVAIL-01 | NFR-AVAIL-01 (99% ช่วง 08.30-16.30 จ-ศ) ไม่มี AC และไม่มี task ใน 006 เป็นช่องโหว่ของ spec ไม่ใช่ของโค้ด เสนอ: ระบุว่าเป็น NFR ของระบบรวม (ตรวจที่ 001 หรือ deploy) หรือเพิ่ม AC | |
| F-05 | ละเมิด Constraint | `backend/internal/returns/handler.go: unavailableBody` | CON-TECH-03 (body `{"code", "message"}` ตารางกลางไม่ได้กำหนดฟิลด์แนบของ `UPSTREAM_UNAVAILABLE`) | 503 ตอบ `{"code":"UPSTREAM_UNAVAILABLE","message":...,"retryable":true}` มีฟิลด์ `retryable` ที่ CON-TECH-03 ไม่ได้กำหนด ขณะที่ 003 เอา `retryable` ออกแล้วตาม CR-006 (`cancelrequest/handler.go` "503 ลองใหม่ได้เสมอ ไม่แนบ retryable") หน้าจอของ 006 ไม่ได้อ่านฟิลด์นี้ ไม่มี test ตรวจ body ของ 503 เสนอ: เอา `retryable` ออก หรือลงทะเบียนฟิลด์แนบในตารางกลาง | |
| F-06 | อ้าง ID ผิดเรื่อง | `specs/006-return/plan.md` ข้อ 6 แถว AC-RET-01 | FR-RET-05 (ยกเลิก 2569-11-11), AC-RET-01 | แผนทดสอบของ AC-RET-01 ยังเขียน "POST `normal` + 1 JPG" ซึ่งเป็นส่วนของการแนบรูปที่ยกเลิกแล้ว test จริง (`TestAC_RET_01_*`, `TestTC_RET_01_1`) ไม่แนบรูป เป็นเอกสารไม่ตรงกับ spec ไม่ใช่โค้ด เสนอ: ลบ "+ 1 JPG" ใน plan.md | |

## 4. แก้แล้ว
| F-ID | แก้อย่างไร | รู้ได้อย่างไร |
|---|---|---|
| - | ข้อค้นพบเดิม F-01 ถึง F-04 ยังเจออยู่ทั้งหมด (F-03 ย้ายที่อยู่จาก App.jsx ไป ReturnsHome.jsx) | - |
