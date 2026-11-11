# RTM: บันทึกการรับคืนอุปกรณ์ (Return)
อ้างอิง: spec.md Draft v2.1 | tasks.md (T-01 ถึง T-23 เสร็จ รอทีมตรวจ, T-24 ถึง T-27 รอคำตอบ) | test-cases.md: ไม่มี
สร้างด้วย /verify เมื่อ 2569-11-11 | test: 50 ผ่าน 0 ไม่ผ่าน (backend 30: auth 3, platform 7, returns 20 / frontend 20)

## 1. ตามรอยไปข้างหน้า (requirement ไป โค้ด ไป test)
| ID | AC | task | โค้ด (ไฟล์: ฟังก์ชัน) | test (ผล) | สถานะ |
|---|---|---|---|---|---|
| FR-RET-01 | AC-RET-03, 08, 10, 11 | T-05, T-06, T-14, T-15, T-18 ถึง T-22 | `repository.go`: `FindOpenLoan`, `ListOpenLoans`, `checkItemDepartment`, `RecordReturn` (`status = 'borrowed'`, `returned_at IS NULL`) / `handler.go`: `lookup`, `openLoans`, `staffOnly` / `ReturnPage.jsx`, `OpenLoansPage.jsx`, `App.jsx` | `TestReturns_Lookup_FoundOpenLoan`, `TestAC_RET_03_*`, `TestAC_RET_08_*`, `TestAC_RET_10_*`, `TestAC_RET_11_*`, `TestReturns_ConcurrentReturn_OneWins` (ผ่าน) / vitest AC-RET-03, AC-RET-10, AC-RET-11, T-15 (ผ่าน) | ช่องโหว่ (F-03) |
| FR-RET-02 | AC-RET-04, AC-RET-09 | T-04, T-07, T-09, T-10 | `repository.go`: `RecordReturn` (`returned_at = now()`) / `service.go`: `RecordReturn` (สภาพจาก `config.ReturnConditionStatus`) / `handler.go`: `recordReturn` ไม่อ่านเวลาจาก form / `ReturnPage.jsx`: ตัวเลือกจาก `RETURN_CONDITIONS` | `TestAC_RET_04_ConditionRequired`, `TestAC_RET_09_ReturnTimeFromServer`, `TestMigration006_LoanReturnColumns`, `TestConfig006_ReturnConditionStatus` (ผ่าน) / vitest AC-RET-04 x2, ตัวเลือกสภาพ (ผ่าน) | ครบ |
| FR-RET-03 | AC-RET-01 | T-04, T-06 | `service.go`: `RecordReturn` -> `config.ReturnConditionStatus["normal"] = "available"` / `repository.go`: `RecordReturn` UPDATE `equipment_item` | `TestAC_RET_01_NormalReturnBecomesAvailable` (ผ่าน) / vitest AC-RET-01 (ผ่าน) | ครบ |
| FR-RET-04 | AC-RET-02 | T-08 | `config.ReturnConditionStatus["damaged"] = "unavailable"` / `RecordReturn` UPDATE สถานะครั้งเดียว | `TestAC_RET_02_DamagedReturnBecomesUnavailable` (ผ่าน ตรวจประวัติสถานะด้วย trigger) | ครบ |
| FR-RET-05 | AC-RET-05, AC-RET-12 (ยกเลิก) | - (task รูปลบออก) | ไม่มี (ตรวจแล้ว: ไม่มีตาราง/โค้ด/ช่องแนบรูป) | vitest "ไม่มีช่องแนบรูป" (ผ่าน) | ยกเลิก 2569-11-11 (ไม่นับ) |
| FR-RET-06 | AC-RET-06, AC-RET-13 | T-11 ถึง T-13, T-16, T-17 | `RecordReturn` (`returned_late = now() > due_at`), `GetLoan`, `ListOpenLoans` (`isOverdue`) / `LoanDetailPage.jsx` ป้าย "คืนเกินกำหนด", `OpenLoansPage.jsx` ป้าย "เลยกำหนดคืน" | `TestAC_RET_06_LateReturnFlagged`, `TestAC_RET_13_OverdueShownToStaff`, `TestReturns_LoanDetail` (ผ่าน) / vitest AC-RET-13 x2 (ผ่าน) | ครบ |
| FR-RET-07 | AC-RET-07 | T-01, T-23 | `audit.go`: `insertAudit`, `WriteAudit`, `outcomeFor` / `RecordReturn` เขียน audit ใน transaction / `handler.go`: `respond`, `staffOnly` | `TestAC_RET_07_ReturnAudited`, `TestMigration006_AuditLogReturnActions` (ผ่าน) | ช่องโหว่ (F-01, F-02) |
| NFR-DATA-01 | AC-EQS-05 (ของ 001) | ไม่มีใน 006 (001 T-16) | `RecordReturn` เขียนทันทีใน transaction / `writeJSON` ตั้ง `Cache-Control: no-store` | ไม่มี: test ของ AC-EQS-05 ยังไม่มีใน 001 | ยังไม่ถึง (รอ 001 T-16) |
| NFR-SEC-01 | ไม่มี AC ใน 006 (ส่วน auth: AC-RET-08) | T-03, T-22 | ทุก route ผ่าน `auth.Middleware` + `staffOnly` / TLS: ไม่มีใน `deploy/docker-compose.yml` | `TestReturns_Skeleton_RoutedThroughAuth` (ผ่าน) / TLS ไม่มี test | ยังไม่ถึง (TLS เป็น deploy ของกลาง ยังไม่ทำ) |
| NFR-AVAIL-01 | ไม่มี AC | ไม่มี | ไม่มี | ไม่มี | ช่องโหว่ (F-04) |
| CON-TECH-01 | ไม่มี AC ตรง | T-01 | migration `0006_return.up.sql`, `pgxpool` ผ่าน `platform.OpenDB` | `TestMigration006_*` บน PostgreSQL 16 จริง (ผ่าน) | ครบ |
| CON-TECH-02 | ไม่มี AC ตรง | T-03 (ใช้ deploy ของ 001) | ใช้ `deploy/docker-compose.yml` เดิม ยังไม่มี HTTPS / reverse proxy | ไม่มี | ยังไม่ถึง (deploy ของกลาง) |
| CON-DATA-01 | AC-RET-01, 02, 03 | T-01, T-04, T-08, T-18 | เขียนเฉพาะ `available` / `unavailable` จาก `borrowed` (CHECK เดิมของ 0001) | `TestAC_RET_01_*`, `TestAC_RET_02_*`, `TestAC_RET_03_*`, `TestMigration_StatusCheck` (001) (ผ่าน) | ครบ |
| CON-DATA-03 | AC-RET-01, 03 | T-04, T-05 | หาการยืมด้วย `asset_code` (UNIQUE) ใน `FindOpenLoan`, `RecordReturn`, `checkItemDepartment` | `TestReturns_Lookup_FoundOpenLoan`, `TestAC_RET_03_*` (ผ่าน) | ครบ |
| DOM-AUDIT-01 | AC-RET-07 | T-01, T-23 | ดู FR-RET-07: action `return_list`, `return_lookup` (เก็บ query), `return_detail`, `return` | `TestAC_RET_07_ReturnAudited` (ผ่าน) | ช่องโหว่ (F-01, F-02) |
| DOM-LOAN-01 | AC-RET-06 | T-13 | `returned_late` / `isOverdue` เทียบ `due_at` ที่ 005 บันทึก ไม่คำนวณ 16:00 เอง | `TestAC_RET_06_LateReturnFlagged` (ผ่าน) | ครบ |
| DOM-REQ-01 | ไม่มี AC | ไม่มี (plan ข้อ 5: ยังไม่ได้ใช้) | ไม่มี: 006 ไม่แตะ `borrow_request` | ไม่มี | รอ Q-08 |
| IF-IDP-01 | AC-RET-08 | T-03, T-14, T-20, T-22 (T-27 รอ Q-SH-08) | `auth.Middleware`, `auth.FromContext` ให้บทบาทและภาควิชา (`principal`, `staffOnly`, ภาควิชาใน repository) | `TestReturns_Skeleton_RoutedThroughAuth`, `TestAC_RET_08_*`, `TestAC_RET_10_*`, `TestAC_RET_11_*` (ผ่าน) | ครบ (กรณีระบบยืนยันตัวตนล่ม รอ Q-SH-08 T-27) |

สรุปสถานะ 18 แถว: ครบ 9 / ยังไม่ถึง 3 / รอ 1 / ช่องโหว่ 4 / ยกเลิก 1 (FR-RET-05)

## 2. ตามรอยย้อนกลับ (โค้ด ไป requirement)
| โค้ด (ไฟล์: ฟังก์ชัน หรือ endpoint) | อ้าง ID | ตรงกับข้อความใน spec ไหม | หมายเหตุ |
|---|---|---|---|
| `POST /api/v1/returns` (`handler.go: recordReturn`, `service.go: RecordReturn`, `repository.go: RecordReturn`) | FR-RET-02, 03, 04, 06, 07 | ตรง | ไม่รับเวลาจาก client, สภาพ 2 ค่า, สถานะตาม config, returned_late, audit ใน transaction |
| `GET /api/v1/returns/lookup` (`lookup`, `FindOpenLoan`) | FR-RET-01(ก) | ตรง | คืน `borrowerId` ให้เจ้าหน้าที่เท่านั้น (staffOnly) |
| `GET /api/v1/returns/open-loans` (`openLoans`, `ListOpenLoans`) | FR-RET-01(ข), FR-RET-06 | ตรง | ไม่แบ่งหน้า ตาม plan Q-P1 (open point ของ plan ไม่ใช่ข้อกำหนด) |
| `GET /api/v1/returns/loans/{loanId}` (`loanDetail`, `GetLoan`) | FR-RET-06 | ตรง | error code `LOAN_NOT_FOUND` (404) ไม่มีในตาราง error ของ plan ข้อ 4 (เป็น HTTP status ทั่วไป ไม่นับเป็นข้อค้นพบ ควรเพิ่มใน plan) |
| `handler.go: staffOnly` | FR-RET-01, AC-RET-08 | ตรง | FR-RET-01 "เจ้าหน้าที่" / RQ-11 |
| `repository.go: checkItemDepartment` | FR-RET-01, ASM-03 | ตรง | ตรวจภาควิชาก่อนการยืมค้าง ตามลำดับใน plan ข้อ 4 |
| `handler.go: respond`, `errorResponse` / `audit.go: outcomeFor` | DOM-AUDIT-01 | ตรงบางส่วน | การจับคู่ outcome และ fail closed เป็นการตัดสินใจของ AI (F-01, F-02) |
| `audit.go` (package returns) | DOM-AUDIT-01, FR-RET-07 | ตรง | 001 ยังไม่มี package `audit` กลาง (plan 001 ข้อ 3) ถ้ามีภายหลังควรย้ายไปใช้ |
| `backend/db/migrations/0006_return.up.sql` | FR-RET-02, 06, 07 | ตรง | ไม่มีตารางรูป (FR-RET-05 ยกเลิก) ตรวจแล้ว / แทนที่ CHECK action ของ 0002 โดยคงค่าเดิมครบ |
| `backend/db/seed/seed_006_return.sql` | (ข้อมูลตัวอย่าง จุดตรวจ AC-RET-10, 11) | ตรง | ข้อมูลตัวอย่าง ไม่ใช่ข้อมูลจริง |
| `backend/internal/config/config.go`: `Timezone`, `ReturnConditionStatus` | glossary, FR-RET-02 ถึง 04 | ตรง | ไม่พบตัวเลขจาก spec ฝังในโค้ด 006 |
| `frontend/src/config.js`: `TIMEZONE`, `STATUS_LABELS`, `RETURN_CONDITIONS` | glossary, CON-DATA-01, FR-RET-02 | ตรง | |
| `frontend/src/api/client.js`: `returns.lookup/record/openLoans/loanDetail`, `authHeaders` | FR-RET-01, 02, 06, IF-IDP-01 | ตรง | token จาก localStorage ทีมเลือก (T-06) |
| `ReturnPage.jsx` | FR-RET-01(ก), 02, 03, 04, AC-RET-03, 04, 11 | ตรง | ข้อความ error ของ `FORBIDDEN_ROLE` แสดงเป็น "ค้นหาไม่สำเร็จ (FORBIDDEN_ROLE)" (หน้าจอเป็นของเจ้าหน้าที่ ไม่นับเป็นข้อค้นพบ) |
| `OpenLoansPage.jsx` | FR-RET-01(ข), FR-RET-06 | ตรง | |
| `LoanDetailPage.jsx` | FR-RET-06 | ตรง | |
| `App.jsx` (แสดง ReturnPage + OpenLoansPage, เลือกแถวแล้วเปิดฟอร์มใหม่ด้วย key) | FR-RET-01(ข) | ตรง | ไม่มี test (F-03) |
| `cmd/server/main.go` ผูก route ของ 006 กับ `platform.NewRouter` ของ 001 | IF-IDP-01 | ตรง | `platform/` เหมือนของ 001 ทุกตัวอักษร |
| `frontend/vite.config.js` เอา rewrite `/api` ออก | (ทีมสั่ง T-06) | - | ไฟล์กลาง กระทบ 001 ควรแจ้งเจ้าของ 001 |

ตรวจ Out of scope ทีละบรรทัด: ไม่พบโค้ดของ นักศึกษาแจ้งคืนเอง (UC-04) / แนบรูป / ปิดการยืมแบบสูญหาย / ระงับคำขอค้าง / ส่งมอบ / ซ่อมแล้วกลับเป็นว่าง / บทลงโทษ / รายงานประวัติรายผู้ยืม / การแจ้งเตือน
ตรวจ Open Questions: ไม่พบโค้ดของ Q-01 (สูญหาย), Q-03 (หน้าผู้ยืม), Q-04 (รายละเอียดความเสียหาย), Q-08 (ระงับคำขอ), Q-09 (สถานะก่อน/หลังใน audit), Q-SH-08
ตรวจข้อมูลที่ต้องปกป้อง: ไม่มีคอลัมน์หรือฟิลด์เลขบัตร/ชื่อผู้ยืม log มีแค่ข้อความ error ไม่มีข้อมูลผู้ใช้ `borrowerId` ส่งเฉพาะเจ้าหน้าที่ภาควิชาเดียวกัน

## 3. ข้อค้นพบ
ชนิด: AC ไม่มี test / test อ่อน / โค้ดไม่มี FR / FR ไม่มี AC / เดา Q-xx / ละเมิด Constraint / ตัวเลขไม่ตรง spec / อ้าง ID ผิดเรื่อง
ทีมตัดสิน: แก้โค้ด / แก้ spec / เพิ่ม Q-xx / ไม่ใช่ปัญหา (พร้อมเหตุผล 1 บรรทัด)

| F-ID | ชนิด | อยู่ที่ | ขัดกับ | รายละเอียด | ทีมตัดสิน |
|---|---|---|---|---|---|
| F-01 | เดา Q-xx (ไม่มี Q รองรับ) | `backend/internal/returns/audit.go: outcomeFor` | DOM-AUDIT-01, DM-07 (outcome มี 3 ค่า) | AI จับคู่ outcome เอง: 2xx = ok, 403 = forbidden, อื่น ๆ (409 `NO_OPEN_LOAN`, 422 `CONDITION_REQUIRED`, 404, 503) = error ทำให้ "บันทึกไม่สำเร็จเพราะไม่มีการยืมค้าง" ปนกับ "ระบบขัดข้อง" ในประวัติ spec/plan/DM-07 ไม่ได้บอกความหมายของแต่ละค่า เสนอ: เพิ่ม Q ให้ทีมกำหนดความหมายของ outcome ใน DM-07 (อาจต้องเพิ่มค่า เช่น rejected ผ่านของกลาง) | |
| F-02 | เดา Q-xx (Q-P4 ของ 001) | `backend/internal/returns/handler.go: respond` | DOM-AUDIT-01 | เขียน audit ไม่สำเร็จแล้วตอบ 503 และไม่ส่งข้อมูล (fail closed) ใช้แนวจาก plan 6.4 ของ 001 ซึ่งยังเป็น Q-P4 ที่ไม่มีคำตอบ ไม่มี test กรณี audit เขียนไม่ได้ (การอ่าน) ส่วนการบันทึกคืน audit อยู่ใน transaction จึง rollback ทั้งก้อนอยู่แล้ว | |
| F-03 | test อ่อน | `frontend/src/App.jsx` (onSelect -> ReturnPage key/initialAssetCode) | AC-RET-10 Then "เลือกรายการเพื่อบันทึกการรับคืนได้" | test ตรวจแยกส่วน: OpenLoansPage เรียก `onSelect` และ ReturnPage รับ `initialAssetCode` แต่ไม่มี test ที่ตรวจการต่อกันใน App ถ้าลบบรรทัด `setSelected` ใน App.jsx ไม่มี test ตัวไหนไม่ผ่าน เสนอ: เพิ่ม test ของ App ด้วย client จำลอง กดเลือกแถวแล้วเห็นฟอร์มของรหัสนั้น | |
| F-04 | FR ไม่มี AC | `specs/006-return/spec.md` Quality Requirements | NFR-AVAIL-01 | NFR-AVAIL-01 (99% ช่วง 08.30-16.30 จ-ศ) ไม่มี AC และไม่มี task ใน 006 เป็นช่องโหว่ของ spec ไม่ใช่ของโค้ด เสนอ: ระบุว่าเป็น NFR ของระบบรวม (ตรวจที่ 001 หรือ deploy) หรือเพิ่ม AC | |

## 4. แก้แล้ว
| F-ID | แก้อย่างไร | รู้ได้อย่างไร |
|---|---|---|
| - | (rtm.md รอบแรก ยังไม่มีข้อค้นพบเดิม) | - |
