# RTM: ตรวจสอบสถานะอุปกรณ์ (Equipment Status)
อ้างอิง: spec.md Draft v5.1 | plan.md Draft v1.9 | tasks.md v2.4 | test-cases.md (ชุด ก ข ค ง)
สร้างด้วย /verify เมื่อ 2569-11-12 04:14 (Asia/Bangkok) ที่ commit d95407a | test: 551 ผ่าน 0 ไม่ผ่าน (SKIP 3)

ผล test ที่รัน (ดูอย่างเดียว ไม่ได้แก้อะไร)
- หลังบ้าน `cd backend && go test ./... -v`: 383 ผ่าน 0 ไม่ผ่าน SKIP 2 (TestTC_EQS_11_7 รอ Q-11, TestTC_EQS_13_7 รอ Q-P3)
  เฉพาะ package ของ 001 และของกลาง: equipment 83 ผ่าน + SKIP 2, audit 4, auth 4, platform 18, config 6
- หน้าจอ `cd frontend && npx vitest run`: 168 ผ่าน 0 ไม่ผ่าน SKIP 1 (TC-EQS-05-5 รอ Q-18)
- คนลอง: TC-EQS-17-6 ผ่าน (Pannathorn J 2569-11-12 บันทึกใน tasks.md จุดตรวจเฟส 15) / TC-EQS-09-7, 10-1 ถึง 10-4 ยังไม่ได้ลอง (รอ Q-SH-08, Q-17)

## 1. ตามรอยไปข้างหน้า (requirement ไป โค้ด ไป test)
| ID | AC | task | โค้ด (ไฟล์: ฟังก์ชัน) | test (ผล) | สถานะ |
|---|---|---|---|---|---|
| FR-EQS-01 | AC-EQS-01, AC-EQS-15 | T-04, T-05, T-06, T-32, T-33, T-23 (รอ Q-05) | equipment/repository.go: ListModels, ListItems / service.go: ListModels, ListItems / frontend EquipmentList.jsx, ModelRow.jsx, ItemTable.jsx | TestTC_EQS_01_*, TestTC_EQS_15_*, TC-EQS-01-2, TC-EQS-15-2, TestAC_EQS_01_* (ผ่าน) | รอ Q-05 (แบ่งหน้ายังไม่ทำ ส่วนอื่นครบ) |
| FR-EQS-02 | AC-EQS-02, AC-EQS-15, AC-EQS-16 | T-07, T-08, T-32, T-33, T-24 (รอ Q-06) | repository.go: filterWhere / handler.go: filterFrom / SearchBar.jsx | TestTC_EQS_02_*, TestTC_EQS_16_*, TC-EQS-02-2, TC-EQS-16-2 (ผ่าน) | รอ Q-06 (ค้นด้วยชื่อหรือรหัสยังไม่มี กรองประเภทและสถานะครบ) |
| FR-EQS-03 | AC-EQS-03 | T-09 | handler.go: nonNil (200 รายการว่าง) / EmptyState.jsx / EquipmentList.jsx | TestTC_EQS_03_1, TC-EQS-03-2, 03-3, 03-4 (ผ่าน) | ครบ |
| FR-EQS-04 | AC-EQS-04, AC-EQS-14 | T-10, T-11, T-12, T-13, T-34 | service.go: GetItem (isOverdue, ตัด dueAt ของนักศึกษา) / dto.go: StudentItemDTO, StaffItemDTO, canRequestLoan / EquipmentDetail.jsx: formatDueAt / StatusBadge.jsx: OverdueBadge | TestTC_EQS_04_*, TestTC_EQS_14_*, TC-EQS-04-2, 04-4, 14-4, 14-5, TestAC_EQS_14_* (ผ่าน) | ครบ |
| FR-EQS-05 | AC-EQS-06 | T-14, T-34 | dto.go: canRequestLoan / StatusBadge.jsx / ItemTable.jsx | TestTC_EQS_06_*, TC-EQS-06-2, 06-4 (ผ่าน) | ครบ |
| FR-EQS-06 | AC-EQS-05 | T-16, T-34 | dto.go: canRequestLoan / ItemTable.jsx: RequestLoanButton / EquipmentList.jsx, EquipmentDetail.jsx (onRequestLoan) / App.jsx ไม่ได้ส่ง onRequestLoan | TestTC_EQS_05_*, TC-EQS-05-2, RequestButton.test.jsx (ผ่าน) | ช่องโหว่ (F-01) |
| FR-EQS-07 | AC-EQS-07 (ตรวจแค่ชื่อผู้ยืม) | T-17, T-28 (รอ Q-P3), T-36 (รอ Q-20) | dto.go: StudentItemDTO ไม่มีช่องผู้ยืม / handler.go: summary / service.go: Summary / repository.go: CountByStatus / หน้าจอสรุปไม่มี | TestTC_EQS_07_*, TestFR_EQS_07_StaffSummaryByStatus, TC-EQS-07-2 (ผ่าน) ชื่อผู้ยืมไม่มี assert (รอ Q-20) | ช่องโหว่ (F-03, F-04) |
| FR-EQS-08 | AC-EQS-08 | T-20, T-26 | handler.go: withDBTimeout, writeUpstreamError / config.go: EquipmentDBTimeout / client.js: getJSON timeoutMs / EquipmentList.jsx: retry / ErrorState.jsx / EquipmentDetail.jsx (ไม่มีปุ่มลองใหม่) | TestTC_EQS_08_*, TestAC_EQS_08_*, TC-EQS-08-2, 08-4, 08-5 (ผ่าน) | ช่องโหว่ (F-02) |
| FR-EQS-10 | AC-EQS-13 | T-18 | service.go: scopeFor, ListModels, ListItems, GetItem / repository.go: filterWhere ($4 ภาควิชา) | TestTC_EQS_13_1 ถึง 13_6, TestAC_EQS_13_OtherDepartment403, TC-EQS-13-5 (ผ่าน) | ครบ |
| FR-EQS-11 | AC-EQS-17 | T-38 | client.js: listItemRequests / EquipmentDetail.jsx (useEffect ตัวที่สอง) / RequestQueue.jsx / endpoint ของ 002 ใน main.go | TC-EQS-17-1 ถึง 17-5, AC-EQS-17.test.jsx (ผ่าน) / TC-EQS-17-6 คนลองผ่าน | ครบ |
| NFR-PERF-01 | AC-EQS-10 | T-22 (รอ Q-17) | ไม่มีโค้ดเฉพาะ (ไม่มี cache ตาม plan ข้อ 1) | TestTC_EQS_10_5 (ผ่าน p95 ราว 21 ms คนเดียว) / k6 ยังไม่มี | รอ Q-17 |
| NFR-DATA-01 | AC-EQS-05 | T-16, T-35 (รอ Q-18) | platform/http.go: NoStore / client.js: cache 'no-store' | TestTC_EQS_05_4, TestAC_EQS_05_ReturnedBecomesAvailable (ผ่าน) / TC-EQS-05-5 SKIP | รอ Q-18 |
| NFR-SEC-01 | ไม่มี AC (AC-EQS-09 ตรวจแค่ส่วนยืนยันตัวตน) | T-03, T-21 / TLS ไม่มี task | auth/middleware.go: Middleware / main.go ห่อทุกฟีเจอร์ / TLS ไม่มี | TestAuth_*, TestTC_EQS_09_* (ผ่าน) | ช่องโหว่ (F-06, F-08) |
| NFR-AVAIL-01 | ไม่มี AC | ไม่มี task | platform/http.go: healthz / deploy/docker-compose.yml มีแค่ db | TestHealthz_* (ผ่าน ตรวจแค่ healthz) | ช่องโหว่ (F-08) |
| NFR-USE-01 | ไม่มี AC | T-27 (รอ Q-09) | plan ข้อ 6.7 mobile-first (class ของ Tailwind) | ไม่มี | ช่องโหว่ (F-08) |
| CON-TECH-01 | ไม่มี AC (Constraint) | T-01, T-02 | platform/db.go: OpenDB / migrate.go / db/migrations | TestMigration_*, test ทุกตัวที่ใช้ testcontainers postgres:16 (ผ่าน) | ครบ |
| CON-TECH-02 | ไม่มี AC (Constraint) | T-01 (Docker) / HTTPS ไม่มี task | deploy/docker-compose.yml (db อย่างเดียว ไม่มี reverse proxy) | ไม่มี | ช่องโหว่ (F-06) |
| CON-TECH-03 | ไม่มี AC (Constraint) | T-03, T-37, T-39, T-40 | platform/apierror.go: WriteError / http.go: MethodNotAllowedJSON / ทุก error ของ equipment และ auth | TestCON_TECH_03_*, errorBody ใน TC ทุกตัว (ผ่าน) | ช่องโหว่ (F-07) |
| CON-DATA-01 | ไม่มี AC (Constraint) ใช้ใน AC-EQS-06, 14 | T-02, T-12, T-14 | migration 0001 CHECK status / dto.go: Status / StatusBadge.jsx: STATUS_LABEL | TestMigration_StatusCheck, TestTC_EQS_14_1 (สถานะยังเป็น borrowed) (ผ่าน) | ครบ |
| CON-DATA-02 | AC-EQS-04 (อ้าง) | T-10, T-40 | handler.go: Routes มีแต่ GET / repository.go: Querier ไม่มี Exec | TestCON_DATA_02_NoWriteEndpoints, TestTC_EQS_04_6 (ผ่าน) | ครบ |
| CON-DATA-03 | ไม่มี AC (Constraint) | T-02, T-05 | migration 0001 asset_code UNIQUE / repository.go: GetItem ตามรหัส | TestMigration_AssetCodeUnique (ผ่าน) | ครบ |
| DOM-PRIV-01 | AC-EQS-07, AC-EQS-14, AC-EQS-17 | T-12, T-17, T-38, T-36 (รอ Q-20) | dto.go: StudentItemDTO (ไม่มีช่องผู้ยืม) / service.go: GetItem / EquipmentDetail.jsx ซ่อนส่วนคำขอเมื่อ 403 | TestTC_EQS_07_1, 07_4, TestAC_EQS_07_StudentNoBorrower, TC-EQS-17-3 (ผ่าน) | ครบ |
| DOM-AUDIT-01 | AC-EQS-11 | T-31, T-19, T-29 (รอ Q-11, Q-P4) | audit/writer.go: Write / handler.go: auditDetail (เฉพาะรายละเอียด) / migration 0002 append-only | TestTC_EQS_11_1 ถึง 11_6, TestAC_EQS_11_DetailWritesAudit, TestMigration_AuditLogAppendOnly (ผ่าน) / 11_7 SKIP | รอ Q-11 (ดู F-05 ด้วย) |
| DOM-ACCESS-01 | AC-EQS-13 | T-18, T-28 (รอ Q-P3) | service.go: scopeFor / repository.go: filterWhere | TestTC_EQS_13_*, TestFR_EQS_10_ScopeByRole (ผ่าน) | ครบ (ส่วนนักศึกษาตามข้อความ DOM-ACCESS-01) |
| DOM-LOAN-01 | AC-EQS-14 | T-12 | service.go: GetItem อ่าน due_at ที่บันทึกไว้ ไม่คำนวณ 16:00 เอง / config.go: TimeZone | TestTC_EQS_14_6, 14_7, 14_8 (ค้างคืน Q-SH-04) (ผ่าน) | ครบ |
| DOM-REQ-01 | AC-EQS-01, AC-EQS-17 | T-04, T-38 | repository.go: ListModels นับว่างจากสถานะอย่างเดียว / RequestQueue.jsx | TestTC_EQS_01_* (คำขอค้างยังนับว่าง), TestTC_EQS_04_5, TC-EQS-17-* (ผ่าน) | ครบ |
| IF-IDP-01 | AC-EQS-09 | T-03, T-21, T-30 (รอ Q-SH-03) | auth/mock_provider.go (จำลอง) / principal.go: Complete / middleware.go / App.jsx, MockLogin.jsx | TestAuth_*, TestTC_EQS_09_*, TC-EQS-09-3, 09-4, 09-5 (ผ่าน) | รอ Q-SH-03 (ระบบจริง) |

FR-EQS-09 และ AC-EQS-12 ยกเลิกโดย CR-002 ไม่นับในตาราง (T-15, T-25 ยกเลิก)
สรุปสถานะ 27 แถว: ครบ 13 / ยังไม่ถึง 0 / รอ 6 / ช่องโหว่ 8

## 2. ตามรอยย้อนกลับ (โค้ด ไป requirement)
| โค้ด (ไฟล์: ฟังก์ชัน หรือ endpoint) | อ้าง ID | ตรงกับข้อความใน spec ไหม | หมายเหตุ |
|---|---|---|---|
| `GET /v1/equipment/models` handler.go: listModels | FR-EQS-01, 02, 10, AC-EQS-13 | ตรง | ไม่มี q (รอ Q-06) ไม่มี page (รอ Q-05) ตามที่ tasks บอก |
| `GET /v1/equipment/models/{modelId}/items` handler.go: listItems | FR-EQS-01, 02, 10 | ตรง | รหัสรุ่นไม่ใช่ตัวเลข = 404 ITEM_NOT_FOUND (T-39) |
| `GET /v1/equipment/items/{assetCode}` handler.go: getItem | FR-EQS-04 ถึง 07, 10, DOM-AUDIT-01 | ตรง | audit fail closed = การเลือกของ plan รอ Q-P4 (F-05) / 404 ไม่บันทึก audit (ทีมตัดสิน TC-EQS-11-5) |
| `GET /v1/equipment/summary` handler.go: summary | FR-EQS-07 | ตรงส่วนหลังบ้าน | ไม่มีหน้าจอเรียก (F-03) นับทุกภาควิชา (F-04) |
| handler.go: withDBTimeout | FR-EQS-08, Q-08 | ตรง | ใช้ config.EquipmentDBTimeout |
| handler.go: requestedDepartment | FR-EQS-10 | ตรง | ค่าที่ไม่ใช่ตัวเลข = -1 นักศึกษาได้ 403 |
| service.go: scopeFor | DOM-ACCESS-01, plan 6.1 | ตรงส่วนนักศึกษา | เจ้าหน้าที่ = ทุกภาควิชา (รอ Q-P3, F-04) |
| service.go: GetItem | FR-EQS-04 ก ข, CR-001, DOM-PRIV-01 | ตรง | เลยกำหนด = now > dueAt (ตรง 16:00 ยังไม่เลย ตรง FR-EQS-04 ข) |
| repository.go: ListModels คอมเมนต์ "ยังไม่จำกัดภาควิชา (T-18)" | T-18 | ไม่ตรงกับโค้ด | โค้ดจำกัดภาควิชาแล้วผ่าน filterWhere คอมเมนต์ค้างจากก่อนทำ T-18 (F-11) |
| dto.go: canRequestLoan | FR-EQS-04, 05, 06 | ตรง | ว่าง = true เท่านั้น |
| audit/writer.go: Write | DOM-AUDIT-01, DM-07 | ตรง | ของกลาง ใช้ร่วม 003 |
| auth/middleware.go: Middleware, mock_provider.go | IF-IDP-01, NFR-SEC-01 | ตรง (แบบจำลอง) | ผู้ใช้จำลอง 4 คน ทีมเลือกเอง รอ Q-SH-03 |
| platform/apierror.go: WriteError, errorMessages | CON-TECH-03 | ตรง | ข้อความภาษาไทยทีมเลือกเอง (plan ข้อ 5) |
| platform/http.go: MethodNotAllowedJSON | CON-TECH-03, CON-DATA-02 | ตรง | ครอบทุกฟีเจอร์ใน main.go (T-40) |
| platform/http.go: healthz | CON-TECH-01, 02 | ไม่ตรง CON-TECH-03 | 503 ตอบ `{"status":"db_unavailable"}` ไม่ใช่ {code, message} (F-07) |
| platform/http.go: NoStore | NFR-DATA-01 | ตรง | |
| cmd/server/main.go: ผูก equipment และ `GET /v1/equipment/items/{assetCode}/requests` ของ 002 | IF-IDP-01, FR-EQS-11 | ตรง | path ที่ไม่มีใต้ /v1/equipment/ ได้ 404 text/plain ของ ServeMux (F-07) |
| config.go: EquipmentDBTimeout, TimeZone | Q-08, FR-EQS-04 | ตรง | |
| frontend config.js: EQUIPMENT_REQUEST_TIMEOUT_MS | Q-08 | ตรง | 6000 มากกว่า 5 วินาทีของหลังบ้าน |
| client.js: listModels, listItems, getItem (timeout), listItemRequests (ไม่มี timeout) | FR-EQS-01 ถึง 04, 08, 11 | ตรง | listItemRequests ไม่ใช้ timeout ของ 001 เพราะเป็น endpoint ของ 002 |
| client.js: session จำลอง getToken/setToken/onUnauthorized | IF-IDP-01, AC-EQS-09 | ตรง (แบบจำลอง) | ทีมเลือกเอง รอ Q-SH-03 |
| App.jsx: PAGES, safePageKey, หน้าเข้าสู่ระบบเมื่อ 401 | AC-EQS-09, plan 6.5 | ตรง | `<Component />` ไม่ได้ส่ง onRequestLoan ให้ EquipmentList (F-01) |
| EquipmentList.jsx | FR-EQS-01, 02, 03, 08 | ตรง | ตัวเลือกประเภทมาจากรายการตอนไม่กรอง (ไม่มี endpoint ประเภทใน plan) |
| EquipmentDetail.jsx | FR-EQS-04, 06, 08, 10, 11 | ไม่ตรง FR-EQS-08 | `<ErrorState error={error} />` ไม่มี onRetry คอมเมนต์ใน catch บอกว่าปุ่มลองใหม่อยู่ใน T-20 แต่หน้านี้ไม่มี (F-02) |
| EquipmentDetail.jsx: formatDueAt, BUDDHIST_YEAR_OFFSET 543 | FR-EQS-04 | ตรง | 543 เป็นค่าปฏิทิน ไม่ใช่ตัวเลขจาก spec รูปแบบ "2569-10-10 16:00" ทีมตัดสิน 2569-11-10 |
| ItemTable.jsx: RequestLoanButton | FR-EQS-06, Q-16 | ตรงส่วนแสดงปุ่ม | กดแล้วไม่ไปไหนในระบบจริง (F-01) คอมเมนต์เขียนไว้ว่า "ยังไม่ได้ต่อ" |
| ModelRow.jsx | FR-EQS-01, AC-EQS-15 | ตรง | ข้อความเมื่อกรองไม่ใช่ "ว่าง" คงรูปแบบเดิม ทีมตัดสินชั่วคราว 2569-11-10 รอ Q-21 |
| SearchBar.jsx | FR-EQS-02, Q-15 | ตรง | ไม่มีช่องค้นด้วยคำ (รอ Q-06) |
| EmptyState.jsx, ErrorState.jsx | FR-EQS-03, 08 | ตรง | |
| StatusBadge.jsx: STATUS_LABEL, OverdueBadge | CON-DATA-01, FR-EQS-04, 05 | ตรง | ป้ายสถานะซ้ำกับ STATUS_LABELS ใน config.js (F-10) |
| RequestQueue.jsx | FR-EQS-11 | ตรง | ไม่มีคำขอ = ซ่อน (ทีมตัดสิน TC-EQS-17-4) |
| MockLogin.jsx | IF-IDP-01, AC-EQS-09 | ตรง (แบบจำลอง) | ไม่ใช่ use case ของระบบ ใช้แทนหน้าเข้าสู่ระบบของมหาวิทยาลัย |
| tc_ac_eqs_10_test.go: const perfP95 | NFR-PERF-01 | ตรงตัวเลข | ตัวเลขจาก spec อยู่ในไฟล์ test ไม่อยู่ใน config.go (F-09) |

ค้นข้อมูลที่ Constraint ปกป้อง (DOM-PRIV-01 ชื่อผู้ยืม / ผู้ขอ): ไม่พบ borrower ใน DTO ของนักศึกษา ตารางรายชิ้น log หรือข้อความ error
log มีที่เดียวคือ writeUpstreamError เขียนข้อความ error ของฐานข้อมูล ไม่มี token หรือข้อมูลผู้ใช้ / audit_log.query ของ 001 ยังไม่ได้ใช้
ไม่พบของแถมที่อยู่ใน Out of scope (ไม่มี endpoint ส่งคำขอ ส่งมอบ รับคืน หรือจัดการอุปกรณ์ในโค้ดของ 001)

## 3. ข้อค้นพบ
ชนิด: AC ไม่มี test / test อ่อน / โค้ดไม่มี FR / FR ไม่มี AC / เดา Q-xx / ละเมิด Constraint / ตัวเลขไม่ตรง spec / อ้าง ID ผิดเรื่อง
ทีมตัดสิน: แก้โค้ด / แก้ spec / เพิ่ม Q-xx / ไม่ใช่ปัญหา (พร้อมเหตุผล 1 บรรทัด)

| F-ID | ชนิด | อยู่ที่ | ขัดกับ | รายละเอียด | ทีมตัดสิน |
|---|---|---|---|---|---|
| F-01 | test อ่อน | frontend/src/App.jsx (`<Component />` บรรทัด 80), ItemTable.jsx: RequestLoanButton, RequestButton.test.jsx | FR-EQS-06 "ส่งข้อมูลอุปกรณ์ที่เลือกไปยัง UC-02", AC-EQS-05 "ปุ่มส่งคำขอยืมกลับมาใช้งานได้" | ในระบบจริง App ไม่ได้ส่ง onRequestLoan ให้ EquipmentList ปุ่ม "ส่งคำขอยืม" จึงกดแล้วไม่เกิดอะไร test ตรวจแค่ว่า callback ได้ assetCode ไม่มี test ไหนพังถ้าปุ่มไม่ต่อกับฟอร์ม ตอนนี้ฟอร์มของ 002 อยู่ใน main แล้ว (BorrowRequestForm รับ prop assetCode, PAGES มี 'borrow-request') แต่ไม่มี task ไหนต่อสองหน้านี้ (plan ข้อ 5 "ต้องตกลงกับ 002") แนวทาง: ตกลงกับ 002 แล้วเพิ่ม task ต่อปุ่มกับหน้า borrow-request พร้อม test ที่กดปุ่มใน App แล้วเห็นฟอร์มที่มีรหัสครุภัณฑ์ | แก้โค้ด ผู้ใช้เจอจริง ฟอร์มของ 002 พร้อมแล้ว ต้องบอก 002 ก่อน (2569-11-12) |
| F-02 | test อ่อน | frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx (`<ErrorState error={error} />`) | FR-EQS-08 "แจ้งข้อผิดพลาดพร้อมเปิดให้ลองใหม่" | หน้ารายละเอียดเมื่อฐานข้อมูลไม่ตอบ (503) หรือหน้าจอเลิกรอ แสดงข้อความแต่ไม่มีปุ่ม "ลองใหม่" มีแค่ปุ่มกลับไปรายการ AC-EQS-08 และ test ทุกตัวตรวจแค่หน้ารายการ จึงไม่มี test จับได้ คอมเมนต์ใน catch เขียนว่า "ปุ่มลองใหม่และ test ของกรณีนี้อยู่ใน T-20" แต่ T-20 ทำเฉพาะหน้ารายการ แนวทาง: ส่ง onRetry ให้ ErrorState ในหน้ารายละเอียด หรือทีมตัดสินว่าปุ่มกลับไปรายการถือเป็นการลองใหม่ | แก้โค้ด งานเล็ก ตรงข้อความ FR-EQS-08 (2569-11-12) |
| F-03 | FR ไม่มี AC | spec.md FR-EQS-07 / backend handler.go: summary (ไม่มีหน้าจอ) | FR-EQS-07 "เจ้าหน้าที่เห็น ... สรุปจำนวนอุปกรณ์แยกตามสถานะ" | AC-EQS-07 ตรวจแค่ชื่อผู้ยืม ส่วนสรุปจำนวนไม่มี AC (tasks.md บอกไว้แล้วใน "สิ่งที่ยังไม่ทำ") หลังบ้านมี `GET /v1/equipment/summary` แต่หน้าจอไม่มีที่ไหนเรียก client.js ไม่มีฟังก์ชันนี้ เจ้าหน้าที่จึงยังไม่เห็นสรุปบนหน้าจอ และไม่มี task หน้าจอ แนวทาง: เสนอ AC ของส่วนสรุปผ่าน CR แล้วเพิ่ม task หน้าจอ หรือทีมตัดสินว่าเทอมนี้มีแค่ API | ไม่ใช่ปัญหา: เทอมนี้มีแค่ API ยังไม่ทำหน้าจอสรุป (2569-11-12) |
| F-04 | เดา Q-P3 | backend service.go: scopeFor, repository.go: CountByStatus / TestFR_EQS_10_ScopeByRole, TestFR_EQS_07_StaffSummaryByStatus | Q-P3 (plan ข้อ 10) ยังไม่มีคำตอบ | โค้ดให้เจ้าหน้าที่เห็นทุกภาควิชา (plan ข้อ 6.1 เขียนว่าชั่วคราว) และ test 2 ตัวที่ไม่ได้มาจาก test-cases.md assert ไว้แล้วว่า "เจ้าหน้าที่ต้องเห็นทุกภาควิชา" และสรุปนับทั้งฐานข้อมูล ถ้า Q-P3 ตอบว่าเห็นเฉพาะภาควิชาตน test ทั้ง 2 จะพัง ส่วน TC-EQS-13-7 และ 07-6 ส่วน (3) เว้นไว้ถูกต้องแล้ว แนวทาง: ถาม Q-P3 หรือย้าย assert ส่วนเจ้าหน้าที่ใน 2 test นั้นออกเป็นคอมเมนต์รอ Q-P3 | เพิ่ม Q-xx (ถาม Q-P3) ตอนนี้เป็นการตัดสินชั่วคราวที่ plan เขียนไว้แล้ว (2569-11-12) |
| F-05 | เดา Q-P4 | backend handler.go: getItem (fail closed) / TestAudit_WriteFailsThenNoData, TestTC_EQS_11_6 ส่วน (1) | Q-P4 (plan ข้อ 10) ยังไม่มีคำตอบ | เขียน audit ไม่สำเร็จแล้วตอบ 503 ไม่ส่งข้อมูล เป็นข้อเสนอของ plan ข้อ 6.4 ที่ยังรอ Q-P4 แต่ TestAudit_WriteFailsThenNoData assert ว่าไม่มีข้อมูลอุปกรณ์ในคำตอบแล้ว ขณะที่ TC-EQS-11-6 ส่วน (2) เว้นไว้เพราะรอ Q-P4 สองที่จึงไม่ตรงกัน แนวทาง: ถาม Q-P4 หรือทีมรับ fail closed เป็นคำตอบแล้วปิด Q-P4 | ไม่ใช่ปัญหา spec เขียนว่า 'ทุกครั้ง (2569-11-12) |
| F-06 | ละเมิด Constraint | deploy/docker-compose.yml (มีแค่ db) / ไม่มี reverse proxy | CON-TECH-02 "HTTPS เท่านั้น", NFR-SEC-01 "TLS 1.2 ขึ้นไป" | ไม่มี HTTPS หรือ TLS ที่ใดในระบบ และไม่มี task ไหนรับงานนี้ (tasks.md ตาราง Constraint และหมายเหตุท้ายไฟล์บอกว่าเป็นงานของทั้งระบบ) ระบบตอนนี้รัน HTTP ล้วน แนวทาง: ทีมตกลงว่าใครทำ reverse proxy แล้วเพิ่ม task หรือบันทึกว่า prototype ไม่ทำ | ไม่ใช่ปัญหาสำหรับ prototype เป็นงานของทั้งระบบ ต้องคุยกับทีม (2569-11-12) |
| F-07 | ละเมิด Constraint | platform/http.go: healthz / ServeMux ของ equipment (path ที่ไม่มี) | CON-TECH-03 "ทุก endpoint ตอบ error ด้วย body {code, message}" | /healthz ตอบ 503 ด้วย `{"status":"db_unavailable"}` และ path ที่ไม่มีใต้ /v1/equipment/ (เช่น /v1/equipment/foo) ได้ 404 ข้อความธรรมดา "404 page not found" จาก ServeMux ไม่ผ่านตัวครอบ MethodNotAllowedJSON ที่ดักแค่ 405 หน้าจอไม่ได้เรียกสองกรณีนี้ ผู้ใช้จึงไม่เห็น แนวทาง: ตัดสินว่า healthz และ path ที่ไม่มีนับเป็น "endpoint" ไหม | ไม่ใช่ปัญหา ผู้ใช้ไม่เห็น (2569-11-12) |
| F-08 | FR ไม่มี AC | spec.md NFR-SEC-01, NFR-AVAIL-01, NFR-USE-01 | AC ทุกข้อใน spec | NFR 3 ข้อไม่มี AC ตรวจ (tasks.md บรรทัด "NFR ที่ไม่มี AC" บอกไว้แล้ว) NFR-AVAIL-01 ไม่มี task เลย NFR-USE-01 ต้องทดสอบกับคน 10 คน (T-27 รอ Q-09) NFR-SEC-01 มีแค่ส่วนยืนยันตัวตนผ่าน AC-EQS-09 แนวทาง: เสนอ AC ผ่าน CR หรือบันทึกว่าตรวจด้วยคนนอกระบบ test | ไม่ใช่ปัญหา ตรวจด้วยคนนอก test (tasks.md บันทึกไว้แล้ว) (2569-11-12) |
| F-09 | ตัวเลขไม่ตรง spec | backend/internal/equipment/tc_ac_eqs_10_test.go: const perfP95 | AGENTS.md กฎข้อ 9 (ตัวเลขจาก spec อยู่ที่ config.go ที่เดียว), NFR-PERF-01 | 3 วินาทีของ NFR-PERF-01 อยู่ในไฟล์ test เพราะ /testcases แก้ config.go ไม่ได้ (บอกไว้ใน prompt-log แล้ว) ตัวเลขตรงกับ spec แต่อยู่ผิดที่ 003 มี CancelPerfP95 ใน config.go เป็นตัวอย่าง แนวทาง: เพิ่ม EquipmentPerfP95 ใน config.go แล้วให้ test อ้างค่านั้น | แก้โค้ด ย้ายไป config.go งาน 2 บรรทัด (2569-11-12) |
| F-10 | ตัวเลขไม่ตรง spec | frontend/src/components/StatusBadge.jsx: STATUS_LABEL / frontend/src/config.js: STATUS_LABELS (ของ 006) | AGENTS.md กฎข้อ 9, CON-DATA-01 | คำของ 3 สถานะมี 2 ที่ คำตรงกันตอนนี้ แต่ถ้าแก้ที่หนึ่ง อีกที่ไม่เปลี่ยนตาม ป้ายบนหน้าจอ 001 กับ 006 จะไม่ตรงกัน แนวทาง: ให้ StatusBadge ใช้ STATUS_LABELS จาก config.js (ต้องคุยกับ 006) | ไม่ใช่ปัญหา คำตรงกันอยู่แล้ว (2569-11-12) |
| F-11 | อ้าง ID ผิดเรื่อง | backend/internal/equipment/repository.go: คอมเมนต์ของ ListModels "และยังไม่จำกัดภาควิชา (T-18)" | T-18 (เสร็จแล้ว) | คอมเมนต์ค้างจากก่อนทำ T-18 โค้ดจำกัดภาควิชาแล้วผ่าน filterWhere ($4) คนอ่านโค้ดอาจเข้าใจผิดว่ายังไม่ป้องกันข้ามภาควิชา แนวทาง: ลบคำว่า "และยังไม่จำกัดภาควิชา (T-18)" | แก้โค้ด ลบคอมเมนต์ 1 บรรทัด (2569-11-12) |

ไม่นับเป็นข้อค้นพบ (task ยังไม่ได้ทำเพราะรอคำตอบ): แบ่งหน้า (T-23 รอ Q-05), ค้นด้วยคำ (T-24 รอ Q-06), ชื่อผู้ยืมของเจ้าหน้าที่ (T-36 รอ Q-20),
audit ของรายการและการค้นหา (T-29 รอ Q-11), load test (T-22 รอ Q-17), รีเฟรชเอง (T-35 รอ Q-18), auth จริง (T-30 รอ Q-SH-03), จอมือถือ (T-27 รอ Q-09)
แถว "ใช้ได้" ที่ยังเขียนว่า "spec ไม่ได้บอก": TC-EQS-10-1 ส่วน (2) และ TC-EQS-10-3 (ยังไม่นับ เพราะ T-22 ยังไม่ได้ทำ แต่ต้องตัดสินก่อนรัน k6)

## 4. แก้แล้ว
| F-ID | แก้อย่างไร | รู้ได้อย่างไร |
|---|---|---|
| (ยังไม่มี rtm.md รอบก่อน) | | |
