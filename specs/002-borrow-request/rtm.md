# RTM: ส่งคำขอยืมอุปกรณ์ (Borrow Request)
อ้างอิง: spec.md Draft v3.4 | plan.md | tasks.md (24 task เสร็จครบ) | test-cases.md (63 แถว "ใช้ได้" ครบ)
สร้างด้วย /verify เมื่อ 2569-11-13 | test: ผ่านทั้งหมด ไม่ผ่าน 0

ผลการรันที่ใช้อ้างอิง (รันบน WSL Ubuntu หลัง merge main ล่าสุด d09a02f)
- `go test ./...` ok ทุก package รวม `borrowrequest` (87 ฟังก์ชัน test ในนั้นเป็น `TestTC_` 51 ตัว)
- `npm test -- --run` 113 ไฟล์ 208 test ผ่าน 1 skipped (ตัวที่ skip คือ `TC-EQS-05-5` ของ 001 ที่รอ Q-18 ไม่ใช่ของ 002)
- test หน้าจอของ 002 มี 13 ไฟล์

**หมายเหตุสำคัญ:** AI รัน `go test` จากฝั่ง Windows ไม่ได้ (Smart App Control บล็อกไฟล์ที่ go สร้าง)
ตัวเลขข้างบนมาจากผลที่ทีมรันใน WSL และส่งให้ดู ไม่ใช่ AI รันเอง

## 1. ตามรอยไปข้างหน้า (requirement ไป โค้ด ไป test)

| ID | AC | task | โค้ด (ไฟล์: ฟังก์ชัน) | test (ผล) | สถานะ |
|---|---|---|---|---|---|
| FR-BRQ-01 | AC-BRQ-01 | T-09, T-19 | `repository.go: Create` (ตรวจ `status='available'` ในคำสั่งเดียวกับ INSERT), `repository.go: ItemStatus`, `handler.go: writeItemNotAvailable` | `TestAC_BRQ_01`, `TestTC_BRQ_01_1/2/3`, `TC-BRQ-01-4.test.jsx` (ผ่าน) | **ช่องโหว่** F-01 |
| FR-BRQ-02 | AC-BRQ-02, AC-BRQ-14 | T-07, T-19 | `validate.go: purposeProblems`, `service.go: normalizePurpose`, `config.go: PurposeValues` | `TestAC_BRQ_02`, `TestAC_BRQ_14`, `TestTC_BRQ_02_1..4`, `TestTC_BRQ_14_1..4`, `TC-BRQ-02-5`, `TC-BRQ-14-5` (ผ่าน) | ครบ |
| FR-BRQ-03 | AC-BRQ-03, AC-BRQ-13 | T-06, T-08, T-19, T-20 | `validate.go: missingFields`, `validate.go: loanWindowProblems`, `handler.go: writeValidationError` | `TestAC_BRQ_03`, `TestAC_BRQ_13`, `TestTC_BRQ_03_1..4`, `TestTC_BRQ_13_1..4`, `TC-BRQ-03-5`, `TC-BRQ-13-5` (ผ่าน) | ครบ |
| FR-BRQ-04 | AC-BRQ-04 | T-04, T-05, T-19 | `service.go: Create`, `repository.go: Create` (ตั้ง `config.StatusRequested`, `submitted_at` จาก DB), `dto.go: CreatedDTO` | `TestAC_BRQ_04`, `TestAC_BRQ_04_HTTP`, `TestTC_BRQ_04_1..3`, `TC-BRQ-04-4` (ผ่าน) | ครบ |
| FR-BRQ-05 | AC-BRQ-05, AC-BRQ-11 | T-11, T-17, T-22 | `repository.go: PendingByItem`, `handler.go: itemRequests`, unique index `uq_request_open_per_item` (migration 0003), `repository.go: isDuplicateRequest` | `TestAC_BRQ_05`, `TestAC_BRQ_11_HTTP`, `TestTC_BRQ_05_1..3`, `TestTC_BRQ_11_1..4`, `TC-BRQ-05-4`, `TC-BRQ-11-5` (ผ่าน) | ครบ |
| FR-BRQ-06 | AC-BRQ-06, AC-BRQ-12 | T-12, T-13 | `service.go: hasOverdueLoan`, `service.go: returnedLateToday`, `clock.go: overdueDeadline`, `repository.go: EarliestOpenLoanDue`, `ReturnedSince` | `TestAC_BRQ_06_HTTP`, `TestAC_BRQ_12`, `TestTC_BRQ_06_1..4`, `TestTC_BRQ_12_1..4` (ผ่าน) | ครบ |
| FR-BRQ-07 | AC-BRQ-07 | T-10, T-21 | `repository.go: checkDepartment` (เรียกก่อนตรวจสถานะ ตาม Q-BRQ-15), `handler.go` แปลงเป็น 403 `FORBIDDEN_DEPARTMENT` | `TestAC_BRQ_07`, `TestCheckOrder_DepartmentBeforeItemStatus`, `TestTC_BRQ_07_1..3` (ผ่าน) | ครบ |
| FR-BRQ-08 | AC-BRQ-08 | T-16, T-24 | `service.go: auditRequest`, `AuditValidationRejected`, `auditOutcome`, `audit.Write` | `TestAC_BRQ_08`, `TestAuditOnValidationRejected`, `TestTC_BRQ_08_1..3` (ผ่าน) | ครบ |
| FR-BRQ-09 | AC-BRQ-09 | T-15, T-20 | `handler.go: create` (503 `UPSTREAM_UNAVAILABLE`), `service.go: createRequest` (ย้อนลบคำขอและไฟล์เมื่อบันทึกไม่สำเร็จ) | `TestAC_BRQ_09`, `TestTC_BRQ_09_1..3`, `TC-BRQ-09-4` (ผ่าน) | ครบ |
| FR-BRQ-10 | AC-BRQ-10 | T-14, T-19 | `attachment.go: parseCreateForm`, `AttachmentStore.Save/Remove`, `validate.go: attachmentProblems`, `handler.go: attachment`, `service.go: OpenAttachment` | `TestAC_BRQ_10`, `TestAC_BRQ_10_AttachmentAccess`, `TestTC_BRQ_10_1..5`, `TC-BRQ-10-6`, `TC-BRQ-10-7`, แถวคน TC-BRQ-10-8 (ผ่าน) | ครบ |
| NFR-SEC-01 | ไม่มี AC | ไม่มี task | `auth.Middleware` (ยืนยันตัวตนทุก endpoint) ส่วน TLS ไม่มีในโค้ด | `TestCreate_WithoutToken_401`, `TestIdentityProviderUnavailable` (ผ่าน) ครอบเฉพาะส่วนยืนยันตัวตน | **ช่องโหว่** F-03 |
| NFR-AVAIL-01 | ไม่มี AC | ไม่มี task | ไม่มีโค้ดที่รับผิดชอบข้อนี้ | ไม่มี | **ช่องโหว่** F-04 |
| NFR-USE-01 | ไม่มี AC | T-19, T-20 | `BorrowRequestForm.jsx` (label ผูก id, `aria-invalid`, `max-w-xl` คอลัมน์เดียว) | `BorrowRequestErrors.test.jsx`, `TC-BRQ-03-5` ครอบเฉพาะ `aria-invalid` (ผ่าน) | **ช่องโหว่** F-05 |
| NFR-PERF-01 | ไม่ใช้กับฟีเจอร์นี้ | ไม่มี task | - | - | ไม่ใช้ (ทีมตัดสิน 2569-11-11 ปิด Q-BRQ-05) |
| CON-TECH-01 | ไม่มี AC | T-01, T-02 | migration `0003_borrow_request.up.sql`, `platform/db.go` (pgx + PostgreSQL 16) | `platform` package ok, ทุก test ใช้ PostgreSQL 16 จริงผ่าน testcontainers | ครบ |
| CON-TECH-02 | ไม่มี AC | ไม่มี task | `deploy/docker-compose.yml` มีเฉพาะฐานข้อมูล ไม่มี container ของหลังบ้าน และไม่มี HTTPS | ไม่มี | **ช่องโหว่** F-02 |
| CON-TECH-03 | ครอบผ่าน AC หลายข้อ | T-05, T-20 | `handler.go: writeError`, `writeValidationError`, `writeItemNotAvailable` ใช้รหัสจากรายการกลางทั้งหมด | `BorrowRequestErrors.test.jsx` (ครบทุกรหัสที่ฟีเจอร์ใช้), `TestTC_BRQ_*` ทุกตัวที่ตรวจ code (ผ่าน) | ครบ |
| CON-DATA-01 | AC-BRQ-01 | T-09 | `repository.go: Create` รับเฉพาะ `status='available'` และไม่มีคำสั่ง UPDATE สถานะอุปกรณ์ที่ใดเลยในฟีเจอร์นี้ | `TestTC_BRQ_01_1/2/3`, `TestTC_BRQ_04_1` ส่วน (7) (ผ่าน) | ครบ |
| CON-DATA-03 | ไม่มี AC | T-01 | `asset_code TEXT NOT NULL UNIQUE` (migration 0001 ของ 001), 002 ใช้ assetCode เป็นตัวระบุและแปลงเป็น `item_id` | `TestCreate_UnknownAssetCodeLooksTheSame` (ผ่าน) | ครบ |
| DOM-PRIV-01 | AC-BRQ-05 | T-11, T-22 | `handler.go: itemRequests` เปิดเฉพาะ staff, `BorrowRequestForm.jsx` แสดงข้อความคงที่ไม่มีข้อมูลผู้อื่น | `TestTC_BRQ_05_3`, `TC-BRQ-05-4`, `TestAC_BRQ_10_AttachmentAccess` (ผ่าน) | ครบ |
| DOM-AUDIT-01 | AC-BRQ-08 (ทางอ้อม) | T-16 | `service.go: auditRequest` เขียนแถว `action='request'` พร้อม `asset_code` ทุกครั้งที่ส่งคำขอ | `TestTC_BRQ_08_1/2` (ผ่าน) | **ช่องโหว่** F-06 |
| DOM-ACCESS-01 | AC-BRQ-07 | T-10 | `repository.go: checkDepartment` (WHERE department_id ของผู้ขอ), `service.go: hasOverdueLoan` ส่ง departmentID ไปกรอง | `TestTC_BRQ_07_1..3`, `TestOverdue_OtherDepartmentNotCounted` (ผ่าน) | ครบ |
| DOM-LOAN-01 | AC-BRQ-13 | T-08 | `validate.go: loanWindowProblems` ใช้ `config.LoanNormalReturnHour` และ `appLoc` | `TestTC_BRQ_13_1..4`, `TestLoanWindow_NormalBoundaryAndOvernight` (ผ่าน) | ครบ |
| DOM-REQ-01 | AC-BRQ-05, AC-BRQ-11 | T-11, T-17 | ไม่มีคอลัมน์สถานะอุปกรณ์ใน `borrow_request`, `PendingByItem` ให้ลำดับ, unique index จำกัดเฉพาะคำขอของตนเอง | `TestTC_BRQ_05_1/2`, `TestTC_BRQ_11_4`, `TestTC_BRQ_01_3` ส่วน (3) (ผ่าน) | ครบ |
| IF-IDP-01 | AC-BRQ-07 (ทางอ้อม) | T-03, T-23 | `auth.Middleware` + `auth.MockProvider` (ของกลาง), `service.go` อ่าน principal จาก context เสมอ | `TestCreate_WithoutToken_401`, `TestIdentityProviderUnavailable` (ผ่าน) | รอ Q-SH-03 (ยังเป็น IdP จำลอง) |

สรุป: ครบ 17 / ช่องโหว่ 6 / รอ 1 / ไม่ใช้ 1 / ยังไม่ถึง 0 (รวม 25 แถว)

## 2. ตามรอยย้อนกลับ (โค้ด ไป requirement)

| โค้ด (ไฟล์: ฟังก์ชัน หรือ endpoint) | อ้าง ID | ตรงกับข้อความใน spec ไหม | หมายเหตุ |
|---|---|---|---|
| `POST /v1/loan-requests` (`handler.go: create`) | FR-BRQ-01 ถึง 10 | ตรง | อยู่ใน plan ข้อ 4 รับทั้ง JSON และ multipart ตามที่ plan ระบุ |
| `GET /v1/loan-requests/{id}/attachment` (`handler.go: attachment`) | FR-BRQ-10 ก | ตรง | plan ข้อ 4 ระบุไว้ เปิดได้เฉพาะเจ้าหน้าที่และเจ้าของคำขอ |
| `GET /v1/equipment/items/{assetCode}/requests` (`handler.go: itemRequests`) | FR-BRQ-05 ก | ตรง | plan ข้อ 4 ระบุไว้ เฉพาะ staff ตาม DOM-PRIV-01 |
| `repository.go: Create` | FR-BRQ-01, 04, 05 ข, 07 | ตรง | ตรวจสถานะ ภาควิชา และคำขอซ้ำในคำสั่งเดียว กันการกดส่งพร้อมกัน (FR-BRQ-09) |
| `service.go: hasOverdueLoan` / `returnedLateToday` | FR-BRQ-06 ก ข ค | ตรง | ใช้ `Clock` ที่ฉีดได้ ไม่เรียก `time.Now()` ในชั้น service ตามหมายเหตุท้าย plan ข้อ 6 |
| `clock.go: overdueDeadline` | FR-BRQ-06 ก, DOM-LOAN-01 | ตรง | คิด 16:00 จาก `config.LoanNormalReturnHour` ไม่ฝังเลข |
| `service.go: auditRequest` | FR-BRQ-08 | ตรง | ไม่ผูก `request_id` โดยตั้งใจ เพราะ `audit_log` append-only แต่คำขอถูกลบได้ (เขียนเหตุผลไว้ในโค้ด) |
| `service.go: normalizePurpose` | FR-BRQ-02 ข | ตรง | ล้างข้อความระบุเมื่อไม่ได้เลือก "อื่น ๆ" ตรงกับข้อความของ FR |
| `attachment.go: AttachmentStore` | FR-BRQ-10 ก, ASM-BRQ-07 | ตรง | เก็บไฟล์บนเครื่องที่รันระบบตามที่ ASM-BRQ-07 เขียนไว้ |
| `attachment.go: safeFileName` | ไม่อ้าง ID | ไม่มี FR รองรับโดยตรง | เป็นการป้องกัน path traversal ซึ่งเป็นงานพื้นฐานด้านความปลอดภัย ไม่ใช่ฟีเจอร์ใหม่ ไม่นับเป็นของแถม |
| `service.go: ItemStatus` | FR-BRQ-01 | อ้าง ID ถูกเรื่อง แต่**ไม่มีผู้เรียกในโค้ดจริง** | ดู F-08 |
| ช่อง "รหัสครุภัณฑ์" ใน `BorrowRequestForm.jsx` | ไม่อ้าง ID | **ไม่ตรง** FR-BRQ-01 เขียนว่ารับรหัสครุภัณฑ์จาก FR-EQS-06 | ดู F-07 |
| `BorrowRequestForm.jsx: toBangkok` | glossary, FR-BRQ-04 | ตรง | เติมเขตเวลาไทยก่อนส่ง ไม่ปล่อยให้หลังบ้านเดาจากเครื่องที่รัน |
| `RequestErrorSummary.jsx` | CON-TECH-03, FR-BRQ-03, FR-BRQ-09 | ตรง | เลือกข้อความจาก "รหัส" ไม่ใช่จากข้อความที่หลังบ้านส่งมา |
| ค่าคงที่ทั้งหมดของฟีเจอร์ | AGENTS.md กฎข้อ 9 | ตรง | ตรวจแล้วไม่มีตัวเลขจาก spec ฝังในโค้ด เลข 16 และ 10MB ปรากฏเฉพาะในคอมเมนต์กับ `config.go` / `config.js` |
| ตาราง `borrow_request` (migration 0003) | DM-06 | ตรง | ไม่มีคอลัมน์เลขบัตรประชาชน ไม่มีอีเมล ไม่มีคอลัมน์สถานะอุปกรณ์ ไม่มีคอลัมน์กำหนดคืน (FR-BRQ-10 ค) |

ตรวจ Out of scope ทีละบรรทัดแล้ว **ไม่พบโค้ดที่ทำสิ่งที่อยู่ใน Out of scope**
โดยเฉพาะ ไม่มีการเปลี่ยนสถานะอุปกรณ์ ไม่มีการคำนวณกำหนดคืน ไม่มีการอนุมัติหรือยกเลิกคำขอ
และหน้าจอของ 002 ไม่แสดงลำดับคำขอของผู้อื่น (เป็นงานของ 001)

## 3. ข้อค้นพบ

ชนิด: AC ไม่มี test / test อ่อน / โค้ดไม่มี FR / FR ไม่มี AC / เดา Q-xx / ละเมิด Constraint / ตัวเลขไม่ตรง spec / อ้าง ID ผิดเรื่อง / ช่องว่างของ spec
ทีมตัดสิน: แก้โค้ด / แก้ spec / เพิ่ม Q-xx / ไม่ใช่ปัญหา (พร้อมเหตุผล 1 บรรทัด)

| F-ID | ชนิด | อยู่ที่ | ขัดกับ | รายละเอียด | ทีมตัดสิน |
|---|---|---|---|---|---|
| F-01 | test อ่อน | `frontend/src/pages/BorrowRequest/BorrowRequestForm.jsx` | FR-BRQ-01, plan.md ข้อ 4 | FR-BRQ-01 สั่งให้ตรวจสถานะอุปกรณ์ **ทั้งตอนเปิดฟอร์มและตอนกดส่ง** และ plan ข้อ 4 เขียนว่า "ตอนเปิดฟอร์ม หน้าจอเรียก `GET /api/v1/equipment/items/{assetCode}` ของ 001" แต่ฟอร์มจริงไม่มี `useEffect` หรือการเรียก API ใด ๆ ตอนเปิด ตรวจเฉพาะตอนกดส่งอย่างเดียว ผลคือผู้ขอกรอกฟอร์มจนจบแล้วค่อยรู้ว่าอุปกรณ์ไม่ว่าง ส่วน test ที่มีอยู่จำลองการเปลี่ยนสถานะที่ฐานข้อมูลแล้วยิง POST จึงไม่จับช่องว่างนี้ **วิธีแก้ที่เสนอ:** เพิ่มการเรียกสถานะตอนเปิดฟอร์มเป็น task ใหม่ หรือแก้ FR ให้ตรงกับที่ทำจริง | |
| F-02 | ละเมิด Constraint | `deploy/docker-compose.yml`, ทั้งระบบ | CON-TECH-02 | CON-TECH-02 บังคับ HTTPS เท่านั้นและ deploy ด้วย Docker แต่ compose มีเฉพาะฐานข้อมูล ไม่มี container ของหลังบ้านหรือหน้าจอ และไม่มี TLS ที่ใด เป็นงานระดับระบบที่ยังไม่มีเจ้าของ **004 บันทึกเรื่องเดียวกันไว้เป็น F-01 ของเขา** ควรยกเป็นงานของทีมรวม ไม่ใช่ของฟีเจอร์ใดฟีเจอร์หนึ่ง | |
| F-03 | FR ไม่มี AC | spec.md หมวด Quality Requirements | NFR-SEC-01 | NFR-SEC-01 มี 2 ส่วนคือ TLS 1.2 ขึ้นไป กับยืนยันตัวตนก่อนทุกครั้ง ส่วนที่สองมี test ครอบ (`TestCreate_WithoutToken_401`) แต่ส่วน TLS ไม่มี AC ไม่มี task และไม่มีใครตรวจ เกี่ยวพันกับ F-02 | |
| F-04 | FR ไม่มี AC | spec.md หมวด Quality Requirements | NFR-AVAIL-01 | "พร้อมใช้งาน 99% ช่วง 08.30-16.30 จันทร์-ศุกร์" ไม่มี AC ไม่มี task และทดสอบไม่ได้ในสภาพแวดล้อมของนักศึกษา ควรเขียนให้ชัดว่าเป็นข้อกำหนดของการ deploy จริง ไม่ใช่ของ prototype | |
| F-05 | FR ไม่มี AC | spec.md, `BorrowRequestForm.jsx` | NFR-USE-01 | NFR-USE-01 ว่า "ใช้งานบนหน้าจอมือถือได้" แต่ไม่มี AC และไม่มี test ที่ตรวจเรื่องขนาดหน้าจอเลย ที่มีคือ `aria-invalid` กับการผูก label ซึ่งเป็นเรื่องการเข้าถึง ไม่ใช่เรื่องมือถือ ฟอร์มใช้ `max-w-xl` คอลัมน์เดียวซึ่งน่าจะใช้บนมือถือได้ แต่ไม่มีใครยืนยัน **เสนอ:** เพิ่มแถว `ตรวจด้วย = คน` ใน test-cases ให้ลองบนมือถือจริง | |
| F-06 | ช่องว่างของ spec | `service.go: auditRequest` | DOM-AUDIT-01 | DOM-AUDIT-01 สั่งให้เขียน audit **ทุกครั้งที่เข้าถึงข้อมูลอุปกรณ์** และ plan ข้อ 5 เคลมว่า "ใช้แล้ว" ที่ขั้นที่ 8 แต่โค้ดจริงเขียนแถวเดียวต่อการส่งคำขอ 1 ครั้ง ด้วย `action='request'` ไม่ได้เขียนแถวแยกสำหรับการอ่านสถานะอุปกรณ์ และ `GET .../requests` ของเจ้าหน้าที่ก็ไม่เขียน audit เลย คำถามคือแถว `request` ถือว่าครอบคลุมการอ่านในขั้นตอนเดียวกันหรือไม่ และ endpoint คิวต้องมี audit ไหม | |
| F-07 | โค้ดไม่มี FR | `BorrowRequestForm.jsx` ช่อง "รหัสครุภัณฑ์" | FR-BRQ-01 | FR-BRQ-01 เขียนว่า "รับรหัสครุภัณฑ์จาก FR-EQS-06" คือมาจากหน้ารายการอุปกรณ์ของ 001 แต่ฟอร์มให้ผู้ขอ**พิมพ์รหัสเอง** ซึ่งไม่มี FR รองรับ และเปิดช่องให้พิมพ์รหัสของอุปกรณ์ที่ตนไม่เคยเห็น (หลังบ้านกันไว้แล้วด้วย 403 และข้อความกลางของ `TestCreate_UnknownAssetCodeLooksTheSame` จึงไม่รั่ว) เป็นทางออกชั่วคราวเพราะปุ่มของ 001 (T-34) ยังไม่ต่อกับฟอร์มนี้ **เสนอ:** ทำให้ช่องนี้เป็นอ่านอย่างเดียวเมื่อเปิดมาจาก 001 หรือบันทึกเป็น Assumption | |
| F-08 | โค้ดไม่มี FR | `service.go: ItemStatus` (บรรทัด 288) | - | เป็น wrapper ที่เรียก `repo.ItemStatus` ต่อ แต่**ไม่มีผู้เรียกในโค้ดจริงเลย** มีแต่ใน test การตรวจสถานะจริงเกิดในคำสั่งเดียวกับ INSERT ที่ `repository.go: Create` อยู่แล้ว น่าจะเป็นของเหลือจากตอนออกแบบ T-09 ไม่กระทบผู้ใช้ แต่ทำให้คนอ่านเข้าใจผิดว่ามีการตรวจ 2 ที่ | |
| F-09 | ช่องว่างของ spec | ข้ามฟีเจอร์ 001 กับ 002 | FR-BRQ-08, Q-BRQ-17 | Q-BRQ-17 ปิดแล้วเมื่อ 2569-11-13 ว่าการเขียนที่สำเร็จไปแล้วใช้ best-effort ซึ่ง 002 ทำตามและมี `TestTC_BRQ_08_3` ยืนยัน แต่ **001 ยังใช้ fail closed** (ตอบ 503 เมื่อเขียน audit ไม่ได้) ระบบจึงยังมีสองพฤติกรรมในเรื่องเดียวกัน เจ้าของ 001 ต้องพิจารณา Q-P4 ของตนเองอีกครั้ง 002 ไม่แก้โค้ดของ 001 | |
| F-10 | ช่องว่างของ spec | `attachment.go: AttachmentStore`, 003 | FR-BRQ-10 ก, ASM-BRQ-07 | ไฟล์ใบบันทึกข้อความถูกเก็บถาวรบนเครื่อง และถูกลบเฉพาะตอนที่บันทึกคำขอไม่สำเร็จเท่านั้น **ไม่มีใครลบเมื่อคำขอถูกยกเลิกโดย 003 หรือถูกปฏิเสธโดย 004** spec ไม่ได้เขียนเรื่องอายุของไฟล์ไว้เลย ในการใช้งานจริงไฟล์จะสะสมไปเรื่อย ๆ **เสนอ:** เปิดเป็น Q ใหม่ว่าไฟล์แนบต้องถูกลบเมื่อไร ใครเป็นคนลบ | |

## 4. แก้แล้ว

| F-ID | แก้อย่างไร | รู้ได้อย่างไร |
|---|---|---|
| - | เป็น rtm.md ฉบับแรกของฟีเจอร์นี้ ยังไม่มีข้อค้นพบเดิมให้ย้าย | - |
