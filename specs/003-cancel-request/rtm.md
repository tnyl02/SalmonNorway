# RTM: ยกเลิกคำขอยืม (Cancel Borrow Request) 003
อ้างอิง: spec.md Draft v3.7 (เนื้อหา ASM-10 แก้ 2569-11-12 แต่บรรทัด Status ยังไม่ได้ขึ้นเวอร์ชัน) | tasks.md (เสร็จ 29/32) | test-cases.md (30 แถว ใช้ได้)
สร้างด้วย /verify เมื่อ 2569-11-12 | test: 161 ผ่าน 0 ไม่ผ่าน
(หลังบ้าน `go test -p 1 ./...` 114 ผ่าน ทั้งโปรเจกต์ ของ 003 อยู่ใน cancelrequest และ migrationtest 47 ตัว / หน้าจอ `npm test` 47 ผ่าน ของ 003 23 ตัว
/ load test k6 AC-CXL-17 รันด้วยมือ 2569-11-12 ผ่าน list p95 106.63 ms, cancel p95 38.9 ms ไม่อยู่ใน go test หรือ npm test)

## 1. ตามรอยไปข้างหน้า (requirement ไป โค้ด ไป test)
| ID | AC | task | โค้ด (ไฟล์: ฟังก์ชัน) | test (ผล) | สถานะ |
|---|---|---|---|---|---|
| FR-CXL-01 | AC-CXL-01, 10, 13 | T-06, T-07, T-14, T-23, T-24 | repository.go: Cancel (UPDATE มีเงื่อนไข), service.go: Cancel (transaction, dbTimeout), handler.go: cancel, MyPendingRequests.jsx: confirmCancel + MSG_CANCEL_SUCCESS, CancelDialog.jsx | TestAC_CXL_01, TestTC_CXL_01_1/3/4, TestAC_CXL_10, TestAC_CXL_13, TC-CXL-01-2, AC_CXL_01, AC_CXL_13 (ผ่าน) | ครบ |
| FR-CXL-02 | AC-CXL-02, 04, 11 | T-15, T-16, T-17, T-18, T-21, T-22 | service.go: Cancel (OwnStatus แยก ALREADY_APPROVED / NOT_CANCELLABLE), handler.go: cancel, MyPendingRequests.jsx: CANCEL_REJECTED | TestAC_CXL_02, TestAC_CXL_04, TestAC_CXL_11, TestTC_CXL_01_6, TestTC_CXL_09_6, AC_CXL_02, AC_CXL_04, AC_CXL_11 (ผ่าน) | ครบ |
| FR-CXL-03 | AC-CXL-03, 08 | T-04, T-19, T-20 | service.go: ListMyPending / Cancel (ใช้ Principal.UserID, ErrStudentOnly, ErrRequestNotFound), repository.go: OwnStatus (WHERE requester_id) | TestAC_CXL_03, TestTC_CXL_01_7, TestTC_CXL_08_6 (GET 403), TestAC_CXL_08, AC_CXL_03 (ผ่าน) | ช่องโหว่ (F-02) |
| FR-CXL-04 | AC-CXL-04 | T-06, T-21, T-22 | repository.go: cancelSQL (`UPDATE ... WHERE status = ANY(cancellable)` คำสั่งเดียว) | TestAC_CXL_04_ApprovedBeforeConfirm, AC_CXL_04 approved while dialog open (ผ่าน) | ครบ |
| FR-CXL-05 | AC-CXL-06, 07 (ส่วน audit) | T-30 (เสร็จ), T-31 | ยังไม่มีโค้ดเขียน audit_log ของ 003 (0005 เพิ่ม action notify แล้ว) | TestMigration_AuditLogAcceptsNotify (ผ่าน) ไม่มี TestAC_CXL_06 | ยังไม่ถึง (T-31 รอ writer.go ของ 001) |
| FR-CXL-06 | AC-CXL-08, 12 | T-04, T-05, T-11, T-12 | repository.go: ListPending (ORDER BY submitted_at DESC, id DESC), MyPendingRequests.jsx (รายการ, MSG_NO_PENDING, MSG_LOAD_FAILED + ลองใหม่) | TestAC_CXL_08, TestTC_CXL_08_1/3/4, TestAC_CXL_12, TestTC_CXL_12_1/3/4, TC-CXL-08-2, 12-2, 12-5, AC_CXL_08, AC_CXL_12 x2 (ผ่าน) | ครบ |
| FR-CXL-07 | AC-CXL-04, 09, 14 | T-07, T-08 | CancelDialog.jsx (ยืนยัน / ไม่ยืนยัน, เหตุผลใน state ของกล่อง), MyPendingRequests.jsx: openDialog / onCancel | AC_CXL_09, TC-CXL-09-1, 09-3, 09-4, 09-5 (ผ่าน) | ครบ |
| FR-CXL-08 | AC-CXL-14 | T-09, T-10 | service.go: ValidateReason (TrimSpace, RuneCountInString, config.CancelReasonMaxLen), CancelDialog.jsx: countReason / reasonProblem, 0003 cancel_reason VARCHAR(200) | TestAC_CXL_14, TestValidateReason_Boundary, TestTC_CXL_14_1/3/4/5/6, TestMigration_CancelReasonMaxLen, AC_CXL_14, TC-CXL-14-2, 14-5 (ผ่าน) | ครบ |
| FR-CXL-09 | AC-CXL-15, 16 | T-27, T-28, T-32 | notifier.go: Notifier / LogNotifier / FormatCanceller / FormatNoticeTime, repository.go: CancelledNotice (ภาควิชาจาก equipment_model), service.go: notify (หลัง commit, ErrNoRecipients) | TestAC_CXL_15, TestAC_CXL_16 (ส่วนไม่เกี่ยว audit), TestNotify_* (ผ่าน) ส่วน audit notify error ยังไม่มี | ยังไม่ถึง (T-32 ส่วน audit) และ ช่องโหว่ (F-04) |
| FR-CXL-10 | ไม่มี AC | ไม่มี task (ไม่ทำใน prototype) | ไม่มี | ไม่มี | ช่องโหว่ (F-07 FR ไม่มี AC ตามที่ spec ตั้งใจ) |
| NFR-SEC-01 | AC-CXL-07 (ส่วนยืนยันตัวตน) | T-25, T-26 | auth.Middleware ครอบทุก route ของ 003 (main.go) ส่วน TLS ไม่มีในระบบ | TestAC_CXL_07 (401 สถานะไม่เปลี่ยน), TestTC_CXL_01_5, TestTC_CXL_08_5 (ผ่าน) | ช่องโหว่ (F-03 ส่วน TLS) |
| NFR-AVAIL-01 | ไม่มี AC | ไม่มี task | ไม่มี | ไม่มี | ช่องโหว่ (F-08) |
| NFR-PERF-CXL-01 | AC-CXL-17 | T-29 | tests/load/cancel_request.js, backend/db/seed/seed_003_load.sql | k6 50 VU: list p95 106.63 ms, cancel p95 38.9 ms, checks 100% (ผ่าน รันด้วยมือ) | ช่องโหว่ (F-04, F-05) |
| CON-TECH-01 | ไม่มี AC ตรง | T-01, T-30 | 0004_staff_contact, 0005_audit_log_notify_action, pgx ทั้งหมด | migrationtest ทั้ง 6 ตัว (ผ่าน) | ครบ |
| CON-TECH-02 | ไม่มี AC ตรง | ไม่มี task ของ 003 (plan ข้อ 5: งานค้างร่วมกับ 001) | deploy/docker-compose.yml มีแค่ฐานข้อมูล ไม่มี HTTPS / reverse proxy | ไม่มี | ช่องโหว่ (F-03) |
| CON-TECH-03 (ใช้ร่วม ทุก endpoint ไม่อยู่ในรายการ "ใช้ร่วม" ของ spec 003) | ไม่มี AC ตรง | ไม่มี task | handler.go: ErrorDTO {code, retryable} ไม่มี message รหัส STUDENT_ONLY, REASON_REQUIRED, REASON_TOO_LONG, REQUEST_NOT_FOUND, ALREADY_APPROVED, NOT_CANCELLABLE ไม่ได้ลงทะเบียน | ไม่มี | ช่องโหว่ (F-01) |
| CON-DATA-01 | AC-CXL-01, 05 | T-01, T-06, T-13 | ไม่มีโค้ดของ 003 เขียน equipment_item (repository.go แก้แค่ borrow_request) | TestAC_CXL_01, TestAC_CXL_05, TestTC_CXL_01_3, TestMigration_NoEquipmentSideEffects, TestSeed003 (ผ่าน) | ครบ |
| DOM-REQ-01 | AC-CXL-05 | T-13 | cancelSQL แก้เฉพาะแถว id + requester_id | TestAC_CXL_05_OtherRequestsOnSameItemUnchanged (ผ่าน) | ครบ |
| IF-IDP-01 | AC-CXL-07 | T-04, T-25, T-26 | ผู้ขอมาจาก auth.FromContext เท่านั้น (handler.go), หน้าจอยังไม่พาไปหน้าเข้าสู่ระบบ | TestAC_CXL_07 (ผ่าน) ส่วนหน้าจอไม่มี | ยังไม่ถึง (T-26 รอหน้าล็อกอินของ 001) |
| IF-NOT-01 | AC-CXL-15, 16 | T-27, T-28 | notifier.go: Notifier interface + LogNotifier จำลอง (ASM-10) ตัวส่งจริงรอ Q-15 | TestAC_CXL_15, TestAC_CXL_16 (ผ่าน) | ครบ (prototype ใช้ตัวจำลองตาม ASM-10) |
| CON-CXL-01 | AC-CXL-01, 05 | T-01, T-06, T-13 | เหมือน CON-DATA-01 | TestMigration_NoEquipmentSideEffects, TestAC_CXL_01, TestTC_CXL_01_3 (ผ่าน) | ครบ |
| CON-CXL-02 | AC-CXL-01, 02, 08, 10, 11 | T-01, T-03 | config.go: CancellableStatuses(), 0003 CHECK 8 ค่า | TestConfig_CancellableStatuses, TestMigration_BorrowRequestStatusCheck, TestAC_CXL_10, TestAC_CXL_11 (ผ่าน) | ครบ |

## 2. ตามรอยย้อนกลับ (โค้ด ไป requirement)
| โค้ด (ไฟล์: ฟังก์ชัน หรือ endpoint) | อ้าง ID | ตรงกับข้อความใน spec ไหม | หมายเหตุ |
|---|---|---|---|
| `GET /v1/my/pending-requests` (handler.go: listPending, service.go: ListMyPending, repository.go: ListPending) | FR-CXL-06, FR-CXL-03 | ตรง | ไม่ส่งรหัสผู้ขอกลับ เจ้าหน้าที่ 403 |
| `POST /v1/borrow-requests/{id}/cancel` (handler.go: cancel, service.go: Cancel) | FR-CXL-01..04, 08 | ตรง ยกเว้นรูปแบบ error | รูปแบบ error ขัด CON-TECH-03 (F-01) id ไม่ใช่ตัวเลข / body อ่านไม่ได้ เป็นการตัดสินของ AI (F-06) |
| repository.go: cancelSQL | FR-CXL-04, DOM-REQ-01, CON-CXL-01 | ตรง | UPDATE คำสั่งเดียว ตัดสินจากสถานะ ณ เวลาบันทึก |
| repository.go: OwnStatus | FR-CXL-02, FR-CXL-03 | ตรง | อ่านเฉพาะของผู้ขอคนนี้ ไม่เปิดเผยคำขอของผู้อื่น |
| repository.go: CancelledNotice, service.go: notify | FR-CXL-09, ASM-10 | ตรง | ส่งหลัง commit แบบรอผล ไม่มีเวลาจำกัด (F-04) |
| notifier.go: LogNotifier, FormatCanceller, FormatNoticeTime, ErrNoRecipients | FR-CXL-09, IF-NOT-01, ASM-10 | ตรง | รูปแบบเวลาและผู้ยกเลิกเป็นการตัดสินของทีม (ASM-10 2569-11-12) log มีชื่อ รหัสนักศึกษา และเหตุผลของผู้ยกเลิก (ไม่มี Constraint ห้าม) |
| service.go: ValidateReason | FR-CXL-08 | ตรง | ใช้ config.CancelReasonMaxLen |
| service.go: dbTimeout = config.CancelDBTimeout | ASM-06 | ตรง | 5 วินาทีทีมเลือกเอง (T-03) |
| config.go: CancellableStatuses, CancelReasonMaxLen, CancelPerfP95, CancelDBTimeout | CON-CXL-02, FR-CXL-08, NFR-PERF-CXL-01, ASM-06 | ตรง | |
| cmd/server/main.go: route /v1/my/, /v1/borrow-requests/ หลัง auth.Middleware | IF-IDP-01, NFR-SEC-01 | ตรง | ไฟล์ของกลาง |
| db/migrations/0004_staff_contact | ASM-10 | ตรง | ตารางจำลองรอ Q-15 |
| db/migrations/0005_audit_log_notify_action | FR-CXL-09, ASM-12 | ตรง | ยังไม่มีโค้ดใช้ action notify (T-32) |
| db/seed/seed_003_cancel_request.sql | ASM-01, ASM-04, ASM-10 | ตรง | ข้อมูลตัวอย่าง |
| db/seed/seed_003_load.sql, tests/load/cancel_request.js | NFR-PERF-CXL-01, AC-CXL-17 | ตรง | ตัวเลข 15000 และ 50 ฝังใน script (F-05) |
| frontend MyPendingRequests.jsx | FR-CXL-06, 07, 01, 02, 03 | ตรง | ข้อความจาก config.js ทั้งหมด ป้าย "ส่งเมื่อ" เป็นคำที่ AI เลือก (T-05) |
| frontend CancelDialog.jsx | FR-CXL-07, 08, ASM-06 | ตรง | ปุ่ม "ลองใหม่" ในกล่องเมื่อ 503 (ทีมตัดสิน T-24) |
| frontend api/client.js: listMyPendingRequests, cancelRequest, postJSON | FR-CXL-06, FR-CXL-01 | ตรง | ไฟล์ของกลาง |
| frontend App.jsx: PAGES (ปุ่มสลับหน้า) | FR-CXL-06 | ตรง | ปุ่มนำทาง ไม่ใช่ของแถม |
| frontend config.js: ข้อความของ 003 9 ตัว | FR-CXL-01, 02, 03, 06, 08, ASM-06, ASM-13 | ตรง | MSG_CANCEL_FAILED ทีมเลือกคำ (T-24) |

หมายเหตุ (ไม่ใช่ข้อค้นพบของโค้ด): spec.md และ plan.md ยังเขียน Q-SH-05 เป็น "รอคำตอบ" ทั้งที่ 001 ปิดแล้วใน CR-005 (2569-11-12)
และบรรทัด Status ของ spec ยังเป็น Draft v3.7 แม้ ASM-10 และ Q-15 ถูกแก้ 2569-11-12 / test-cases.md อ้าง spec Draft v3.3

## 3. ข้อค้นพบ
ชนิด: AC ไม่มี test / test อ่อน / โค้ดไม่มี FR / FR ไม่มี AC / เดา Q-xx / ละเมิด Constraint / ตัวเลขไม่ตรง spec / อ้าง ID ผิดเรื่อง
ทีมตัดสิน: แก้โค้ด / แก้ spec / เพิ่ม Q-xx / ไม่ใช่ปัญหา (พร้อมเหตุผล 1 บรรทัด)

| F-ID | ชนิด | อยู่ที่ | ขัดกับ | รายละเอียด | ทีมตัดสิน |
|---|---|---|---|---|---|
| F-01 | ละเมิด Constraint | backend/internal/cancelrequest/handler.go: listPending, cancel, writeUpstreamError / dto.go: ErrorDTO | CON-TECH-03 (CR-004 2569-11-11) | CON-TECH-03 กำหนดให้ทุก endpoint ตอบ error เป็น `{"code","message"}` ด้วยรหัสจากรายการกลาง แต่ 003 ตอบ `{"code","retryable"}` ไม่มี `message` และใช้รหัสที่ไม่ได้ลงทะเบียน 6 ตัว (STUDENT_ONLY ซึ่งรายการกลางมี FORBIDDEN_ROLE ความหมายเดียวกัน, REASON_REQUIRED / REASON_TOO_LONG ซึ่งรายการกลางมี VALIDATION_FAILED 422, REQUEST_NOT_FOUND, ALREADY_APPROVED, NOT_CANCELLABLE) spec 003 เขียนก่อน CR-004 และไม่ได้ใส่ CON-TECH-03 ในรายการ "ใช้ร่วม" หน้าจอของ 003 ตัดสินจาก code อยู่แล้ว ทางแก้ที่เป็นไปได้: ลงทะเบียนรหัสของ 003 ในตารางกลาง หรือเปลี่ยนไปใช้รหัสกลาง แล้วเพิ่ม message (ต้องแก้ plan ข้อ 4, test และ TC ที่ตรวจรหัส) || แก้ spec และโค้ด: CR-006 ลงทะเบียน REQUEST_NOT_FOUND (ร่วม 004), ALREADY_APPROVED, NOT_CANCELLABLE (แนบ status) ใช้ FORBIDDEN_ROLE, VALIDATION_FAILED 422 ตัด retryable ทำใน T-33 (ทีมตัดสิน 2569-11-12) |
| F-02 | test อ่อน | backend/internal/cancelrequest/service.go: Cancel (`p.Role != auth.RoleStudent`) | FR-CXL-03 (ASM-13: เจ้าหน้าที่ที่สั่งยกเลิกต้องได้ 403) | มี test ว่าเจ้าหน้าที่เรียกรายการได้ 403 (TestTC_CXL_08_6, TestAC_CXL_08) แต่ไม่มี test ที่เจ้าหน้าที่เรียก `POST .../cancel` ถ้าลบบรรทัดตรวจบทบาทใน Cancel ออก เจ้าหน้าที่จะได้ 404 REQUEST_NOT_FOUND แทน 403 และไม่มี test ตัวไหนพัง | แก้โค้ด: เพิ่ม test TestFR_CXL_03_StaffCancelForbidden ไม่แก้โค้ด (ทีมตัดสิน 2569-11-12) |
| F-03 | ละเมิด Constraint | deploy/docker-compose.yml (ระบบทั้งหมด) | CON-TECH-02, NFR-SEC-01 (TLS 1.2 ขึ้นไป) | ไม่มี HTTPS หรือ reverse proxy ในระบบ หลังบ้านฟัง HTTP ที่ :8000 plan ข้อ 5 ของ 003 บันทึกไว้แล้วว่าเป็นงานค้างร่วมกับ 001 ไม่มี task ของใครรับ | |
| F-04 | ตัวเลขไม่ตรง spec | backend/internal/cancelrequest/service.go: Cancel / notify (`context.WithoutCancel`) | NFR-PERF-CXL-01 ("ไม่นับเวลาส่ง Email"), AC-CXL-16 ("ระบบแจ้งเตือน Email ไม่ตอบสนอง") | Email ส่งแบบรอผลก่อนตอบนักศึกษา (ทีมตัดสิน ASM-10) และไม่มีเวลาจำกัดการส่ง ถ้าตัวส่งค้าง นักศึกษาจะรอไม่สิ้นสุด และเวลาตอบสนองของการยกเลิกรวมเวลาส่ง Email ซึ่ง spec บอกว่าไม่นับ ตัวจำลองตอนนี้ไม่ค้างจึงยังไม่เกิด ค้างจาก T-28 ว่าจะเพิ่ม NotifyTimeout ใน config หรือจดเป็นคำถามคู่ Q-15 | |
| F-05 | ตัวเลขไม่ตรง spec (ตัวเลขฝังในโค้ด) | tests/load/cancel_request.js: P95_LIMIT_MS = 15000, VUS = 50 | NFR-PERF-CXL-01, AGENTS.md กฎข้อ 9 | ค่าตรงกับ spec แต่ฝังใน script แทนที่จะอ่านจาก config.CancelPerfP95 (script ของ k6 อ่าน config ของ Go ไม่ได้ มีคอมเมนต์บอกให้แก้คู่กัน และเปลี่ยนผ่าน env P95_LIMIT_MS ได้) | |
| F-06 | โค้ดไม่มี FR | backend/internal/cancelrequest/handler.go: cancel | plan.md ข้อ 4 ตาราง error | 2 กรณีที่ plan ไม่ได้กำหนด AI เลือกเอง: body อ่านไม่ได้ (ไม่ใช่ JSON) ตอบ 400 REASON_REQUIRED และ id ที่ไม่ใช่ตัวเลขตอบ 404 REQUEST_NOT_FOUND ผู้ใช้ผ่านหน้าเว็บไม่เจอ เกิดเฉพาะเรียก API ตรง || แก้ spec: รวมกับ F-01 ใน CR-006 plan ข้อ 4 กำหนดแล้ว body อ่านไม่ได้ -> 422 VALIDATION_FAILED, id ไม่ใช่ตัวเลข -> 404 REQUEST_NOT_FOUND (ทีมตัดสิน 2569-11-12) |
| F-07 | FR ไม่มี AC | spec.md FR-CXL-10 | | FR-CXL-10 (เจ้าหน้าที่ยกเลิกแทน) ไม่มี AC และไม่มีโค้ด spec เขียนไว้ชัดว่าไม่ทำใน prototype และยังไม่มี AC จนกว่าจะนำเข้า scope (รอ Q-16) | |
| F-08 | FR ไม่มี AC | spec.md NFR-AVAIL-01, NFR-SEC-01 (ส่วน TLS) | | NFR-AVAIL-01 (99% ช่วง 08.30-16.30 น.) และส่วนเข้ารหัส TLS ของ NFR-SEC-01 ไม่มี AC ใน spec 003 วัดผลไม่ได้ใน prototype | |

## 4. แก้แล้ว
| F-ID | แก้อย่างไร | รู้ได้อย่างไร |
|---|---|---|
| F-02 | เพิ่ม TestFR_CXL_03_StaffCancelForbidden ใน handler_test.go (เจ้าหน้าที่ POST cancel ได้ 403 STUDENT_ONLY สถานะไม่เปลี่ยน) ไม่แก้โค้ด | test ผ่าน และลองปิดบรรทัดตรวจบทบาทใน service.go: Cancel ชั่วคราว test ไม่ผ่าน (ได้ 404 แทน 403) แล้วคืนไฟล์ 2569-11-12 ต้องรัน /verify ซ้ำเพื่อยืนยันตามขั้นตอน |
| F-01 | CR-006 + T-33: errors.go (รหัสและข้อความ), ErrorDTO {code, message, fields, status}, handler.go ใช้รหัสกลาง, service.go แนบสถานะใน 409 (StatusError) แก้ constraints v5.2, spec v3.8, plan ข้อ 4, test-cases 4 แถว | go test ผ่านทั้งหมด (ของ 003 48) test ตรวจ body ทั้ง message, fields [reason] และ status ของ 409 |
| F-06 | plan ข้อ 4 กำหนด 2 กรณีแล้ว (CR-006) body อ่านไม่ได้เปลี่ยนเป็น 422 VALIDATION_FAILED | ตรวจด้วยมือหลังเปิดหลังบ้าน (ยังไม่มี test เฉพาะ) |
