# RTM: ยกเลิกคำขอยืม (Cancel Borrow Request) 003
อ้างอิง: spec.md Draft v3.8 | tasks.md (เสร็จ 33/33) | test-cases.md (30 แถว ใช้ได้)
สร้างด้วย /verify เมื่อ 2569-11-12 (รอบที่ 2 หลัง T-26, T-31, T-32, T-33) | test: 186 ผ่าน 0 ไม่ผ่าน
(หลังบ้าน `go test -count=1 -p 1 ./...` 129 ผ่าน ทั้งโปรเจกต์ ของ 003 อยู่ใน cancelrequest และ migrationtest 50 ตัว / หน้าจอ `npm test` 57 ผ่าน ของ 003 26 ตัว
/ load test k6 AC-CXL-17 รันด้วยมือ 2569-11-12 ผ่าน list p95 106.63 ms, cancel p95 38.9 ms ไม่อยู่ใน go test หรือ npm test ไม่ได้รันซ้ำรอบนี้)

## 1. ตามรอยไปข้างหน้า (requirement ไป โค้ด ไป test)
| ID | AC | task | โค้ด (ไฟล์: ฟังก์ชัน) | test (ผล) | สถานะ |
|---|---|---|---|---|---|
| FR-CXL-01 | AC-CXL-01, 10, 13 | T-06, T-07, T-14, T-23, T-24 | repository.go: Cancel (UPDATE มีเงื่อนไข), service.go: Cancel (transaction, dbTimeout), handler.go: cancel, MyPendingRequests.jsx: confirmCancel + MSG_CANCEL_SUCCESS, CancelDialog.jsx | TestAC_CXL_01, TestTC_CXL_01_1/3/4, TestAC_CXL_10, TestAC_CXL_13, TC-CXL-01-2, AC_CXL_01, AC_CXL_13 (ผ่าน) | ครบ |
| FR-CXL-02 | AC-CXL-02, 04, 11 | T-15, T-16, T-17, T-18, T-21, T-22, T-33 | service.go: Cancel (RequestState แยก ALREADY_APPROVED / NOT_CANCELLABLE แนบ status ผ่าน StatusError), handler.go: cancel, MyPendingRequests.jsx: CANCEL_REJECTED | TestAC_CXL_02, TestAC_CXL_04, TestAC_CXL_11 (ตรวจ status ใน 409), TestTC_CXL_01_6, TestTC_CXL_09_6, AC_CXL_02, AC_CXL_04, AC_CXL_11 (ผ่าน) | ครบ |
| FR-CXL-03 | AC-CXL-03, 08 | T-04, T-19, T-20, T-26 | service.go: ListMyPending / Cancel (Principal.UserID, ErrStudentOnly, owner != p.UserID -> ErrRequestNotFound), repository.go: RequestState, MyPendingRequests.jsx: FORBIDDEN_ROLE ไม่มีปุ่มลองใหม่ | TestAC_CXL_03, TestTC_CXL_01_7, TestTC_CXL_08_6, TestAC_CXL_08, TestFR_CXL_03_StaffCancelForbidden, AC_CXL_03, FR_CXL_03 staff sees forbidden message without retry (ผ่าน) | ครบ |
| FR-CXL-04 | AC-CXL-04 | T-06, T-21, T-22 | repository.go: cancelSQL (`UPDATE ... WHERE status = ANY(cancellable)` คำสั่งเดียว) | TestAC_CXL_04_ApprovedBeforeConfirm, AC_CXL_04 approved while dialog open (ผ่าน) | ครบ |
| FR-CXL-05 | AC-CXL-06, 07 (ส่วน audit) | T-30, T-31 | service.go: Cancel (audit.Write ok + reason_len), auditRejected (forbidden + cause, commit), repository.go: RequestState (แยก not_found / not_owner) ใช้ audit/writer.go ของ 001 ใน transaction เดียวกัน | TestAC_CXL_06_AuditLogAttempts, TestFR_CXL_05_AuditLogOtherCauses, TestAC_CXL_07 (ไม่มีแถวใหม่), TestMigration_AuditLogAcceptsNotify (ผ่าน) | ช่องโหว่ (F-09) |
| FR-CXL-06 | AC-CXL-08, 12 | T-04, T-05, T-11, T-12 | repository.go: ListPending (ORDER BY submitted_at DESC, id DESC), MyPendingRequests.jsx (รายการ, MSG_NO_PENDING, MSG_LOAD_FAILED + ลองใหม่) | TestAC_CXL_08, TestTC_CXL_08_1/3/4, TestAC_CXL_12, TestTC_CXL_12_1/3/4, TC-CXL-08-2, 12-2, 12-5, AC_CXL_08, AC_CXL_12 x2 (ผ่าน) | ครบ |
| FR-CXL-07 | AC-CXL-04, 09, 14 | T-07, T-08 | CancelDialog.jsx (ยืนยัน / ไม่ยืนยัน, เหตุผลใน state ของกล่อง), MyPendingRequests.jsx: openDialog / onCancel | AC_CXL_09, TC-CXL-09-1, 09-3, 09-4, 09-5 (ผ่าน) | ครบ |
| FR-CXL-08 | AC-CXL-14 | T-09, T-10, T-33 | service.go: ValidateReason (TrimSpace, RuneCountInString, config.CancelReasonMaxLen), errors.go: msgValidationFailed (ตัวเลขจาก config), CancelDialog.jsx: countReason / reasonProblem, 0003 cancel_reason VARCHAR(200) | TestAC_CXL_14 (422 message fields), TestValidateReason_Boundary, TestTC_CXL_14_1/3/4/5/6, TestMigration_CancelReasonMaxLen, AC_CXL_14, TC-CXL-14-2, 14-5 (ผ่าน) | ครบ |
| FR-CXL-09 | AC-CXL-15, 16 | T-27, T-28, T-32 | notifier.go: Notifier / LogNotifier / FormatCanceller / FormatNoticeTime, repository.go: CancelledNotice, service.go: notify (หลัง commit, ErrNoRecipients), auditNotifyFailed (notify error {"for":"cancel"}) | TestAC_CXL_15, TestAC_CXL_16 (รวมแถว notify error), TestNotify_NoStaffEmailStillCancelled (รวมแถว notify error), TestNotify_* (ผ่าน) | ช่องโหว่ (F-04) |
| FR-CXL-10 | ไม่มี AC | ไม่มี task (ไม่ทำใน prototype) | ไม่มี | ไม่มี | ช่องโหว่ (F-07 FR ไม่มี AC ตามที่ spec ตั้งใจ) |
| NFR-SEC-01 | AC-CXL-07 (ส่วนยืนยันตัวตน) | T-25, T-26 | auth.Middleware ครอบทุก route ของ 003 (main.go) ส่วน TLS ไม่มีในระบบ | TestAC_CXL_07, TestTC_CXL_01_5, TestTC_CXL_08_5, AC_CXL_07 redirects to login and back (ผ่าน) | ช่องโหว่ (F-03, F-08 ส่วน TLS) |
| NFR-AVAIL-01 | ไม่มี AC | ไม่มี task | ไม่มี | ไม่มี | ช่องโหว่ (F-08) |
| NFR-PERF-CXL-01 | AC-CXL-17 | T-29 | tests/load/cancel_request.js, backend/db/seed/seed_003_load.sql | k6 50 VU: list p95 106.63 ms, cancel p95 38.9 ms, checks 100% (ผ่าน รันด้วยมือก่อน T-31, T-32 ไม่ได้รันซ้ำหลังเพิ่ม audit) | ช่องโหว่ (F-04, F-05) |
| CON-TECH-01 | ไม่มี AC ตรง | T-01, T-30 | 0004_staff_contact, 0005_audit_log_notify_action, pgx ทั้งหมด | migrationtest ทั้ง 6 ตัว (ผ่าน) | ครบ |
| CON-TECH-02 | ไม่มี AC ตรง | ไม่มี task ของ 003 (plan ข้อ 5: งานค้างร่วมกับ 001) | deploy/docker-compose.yml มีแค่ฐานข้อมูล ไม่มี HTTPS / reverse proxy | ไม่มี | ช่องโหว่ (F-03) |
| CON-TECH-03 | ไม่มี AC ตรง | T-33 (CR-006) | errors.go (รหัสและข้อความ), dto.go: ErrorDTO {code, message, fields, status}, handler.go ใช้รหัสจากรายการกลาง v5.2 ไม่มี retryable | TestAC_CXL_03 (body ตรงทุกตัวอักษร), TestAC_CXL_13 (503 body), TestAC_CXL_14 (422), TestAC_CXL_02/11 (409 status), TestTC_CXL_08_6, TestTC_CXL_14_4/5/6 (ผ่าน) | ช่องโหว่ (F-10) |
| CON-DATA-01 | AC-CXL-01, 05 | T-01, T-06, T-13 | ไม่มีโค้ดของ 003 เขียน equipment_item (repository.go แก้แค่ borrow_request, audit_log) | TestAC_CXL_01, TestAC_CXL_05, TestTC_CXL_01_3, TestMigration_NoEquipmentSideEffects, TestSeed003 (ผ่าน) | ครบ |
| DOM-REQ-01 | AC-CXL-05 | T-13 | cancelSQL แก้เฉพาะแถว id + requester_id | TestAC_CXL_05_OtherRequestsOnSameItemUnchanged (ผ่าน) | ครบ |
| IF-IDP-01 | AC-CXL-07 | T-04, T-25, T-26 | ผู้ขอมาจาก auth.FromContext เท่านั้น (handler.go), หน้าจอใช้กลไก 401 ของ 001 T-21 (client.js onUnauthorized, App.jsx, MockLogin.jsx, safePageKey) | TestAC_CXL_07, AC_CXL_07 redirects to login and back (render App + client จริง) (ผ่าน) | ครบ |
| IF-NOT-01 | AC-CXL-15, 16 | T-27, T-28 | notifier.go: Notifier interface + LogNotifier จำลอง (ASM-10) ตัวส่งจริงรอ Q-15 | TestAC_CXL_15, TestAC_CXL_16 (ผ่าน) | ครบ (prototype ใช้ตัวจำลองตาม ASM-10) |
| CON-CXL-01 | AC-CXL-01, 05 | T-01, T-06, T-13 | เหมือน CON-DATA-01 | TestMigration_NoEquipmentSideEffects, TestAC_CXL_01, TestTC_CXL_01_3 (ผ่าน) | ครบ |
| CON-CXL-02 | AC-CXL-01, 02, 08, 10, 11 | T-01, T-03 | config.go: CancellableStatuses(), 0003 CHECK 8 ค่า | TestConfig_CancellableStatuses, TestMigration_BorrowRequestStatusCheck, TestAC_CXL_10, TestAC_CXL_11 (ผ่าน) | ครบ |

## 2. ตามรอยย้อนกลับ (โค้ด ไป requirement)
| โค้ด (ไฟล์: ฟังก์ชัน หรือ endpoint) | อ้าง ID | ตรงกับข้อความใน spec ไหม | หมายเหตุ |
|---|---|---|---|
| `GET /v1/my/pending-requests` (handler.go: listPending, service.go: ListMyPending, repository.go: ListPending) | FR-CXL-06, FR-CXL-03 | ตรง | ไม่ส่งรหัสผู้ขอกลับ เจ้าหน้าที่ 403 FORBIDDEN_ROLE |
| `POST /v1/borrow-requests/{id}/cancel` (handler.go: cancel, service.go: Cancel) | FR-CXL-01..05, 08 | ตรง | ยิงหลังบ้านจริง 2569-11-12: body ไม่ใช่ JSON -> 422 VALIDATION_FAILED fields [reason], id "abc" -> 404 ไม่พบคำขอ, เจ้าหน้าที่ -> 403 FORBIDDEN_ROLE ตาม plan ข้อ 4 / 2 กรณีแรกไม่บันทึก audit (F-09) |
| repository.go: cancelSQL | FR-CXL-04, DOM-REQ-01, CON-CXL-01 | ตรง | UPDATE คำสั่งเดียว ตัดสินจากสถานะ ณ เวลาบันทึก |
| repository.go: RequestState | FR-CXL-02, FR-CXL-03, FR-CXL-05 | ตรง | อ่านเจ้าของเพื่อแยก cause ใน audit เท่านั้น หน้าจอได้ 404 เหมือนกันทั้งไม่มีคำขอและของผู้อื่น |
| service.go: auditRejected, Cancel (audit ok) | FR-CXL-05, ASM-07, ASM-12 | ตรงกับตาราง FR-CXL-05 | เหตุผลผิดกับคำขอที่ไม่มีอยู่ -> request_id ว่าง เก็บ id ใน query (AI เลือก เจ้าของรับตอนตรวจ T-31) บันทึกไม่ได้ -> 503 |
| service.go: auditNotifyFailed | FR-CXL-09, FR-CXL-05 | ตรง | เขียนหลัง commit นอก transaction รอไม่เกิน dbTimeout บันทึกไม่ได้แค่ log (AI เลือก) |
| repository.go: CancelledNotice, service.go: notify | FR-CXL-09, ASM-10 | ตรง | ส่งหลัง commit แบบรอผล ไม่มีเวลาจำกัด (F-04) |
| notifier.go: LogNotifier, FormatCanceller, FormatNoticeTime, ErrNoRecipients | FR-CXL-09, IF-NOT-01, ASM-10 | ตรง | log มีชื่อ รหัสนักศึกษา และเหตุผลของผู้ยกเลิก (ไม่มี Constraint ห้าม) |
| service.go: ValidateReason | FR-CXL-08 | ตรง | ใช้ config.CancelReasonMaxLen |
| errors.go: รหัส ข้อความ StatusError, dto.go: ErrorDTO | CON-TECH-03, CR-006 | ตรงกับ constraints.md v5.2 | msgUnauthenticated "กรุณาเข้าสู่ระบบ" ไม่ถูกใช้จริง (middleware ตอบ 401 ก่อนด้วย "กรุณาเข้าสู่ระบบก่อนใช้งาน") ไม่กระทบผู้ใช้ |
| service.go: dbTimeout = config.CancelDBTimeout | ASM-06 | ตรง | 5 วินาทีทีมเลือกเอง (T-03) |
| config.go: CancellableStatuses, CancelReasonMaxLen, CancelPerfP95, CancelDBTimeout | CON-CXL-02, FR-CXL-08, NFR-PERF-CXL-01, ASM-06 | ตรง | |
| cmd/server/main.go: route /v1/my/, /v1/borrow-requests/ หลัง auth.Middleware | IF-IDP-01, NFR-SEC-01 | ตรง | ไฟล์ของกลาง |
| db/migrations/0004_staff_contact | ASM-10 | ตรง | ตารางจำลองรอ Q-15 FK ไป department |
| db/migrations/0005_audit_log_notify_action | FR-CXL-09, ASM-12 | ตรง | ใช้แล้วใน auditNotifyFailed (T-32) |
| db/seed/seed_003_cancel_request.sql | ASM-01, ASM-04, ASM-10 | ตรง | ON CONFLICT DO NOTHING โหลดซ้ำไม่คืนสถานะคำขอที่ยกเลิกไปแล้ว |
| db/seed/seed_003_load.sql, tests/load/cancel_request.js | NFR-PERF-CXL-01, AC-CXL-17 | ตรง | ตัวเลข 15000 และ 50 ฝังใน script (F-05) |
| frontend MyPendingRequests.jsx | FR-CXL-06, 07, 01, 02, 03 | ตรง | 403 FORBIDDEN_ROLE แสดง serverMessage ไม่มีปุ่มลองใหม่ (plan v1.6 ข้อ 4 ทีมตัดสิน T-26) |
| frontend CancelDialog.jsx | FR-CXL-07, 08, ASM-06 | ตรง | ปุ่ม "ลองใหม่" ในกล่องเมื่อ 503 (ทีมตัดสิน T-24) |
| frontend api/client.js: listMyPendingRequests, cancelRequest, postJSON | FR-CXL-06, FR-CXL-01 | ตรง ยกเว้นคอมเมนต์ | คอมเมนต์ของ cancelRequest ยังเขียน 400 REASON_* (F-10) ไฟล์ของกลาง |
| frontend App.jsx, MockLogin.jsx (ของ 001 T-21) | IF-IDP-01 (AC-CXL-07) | ตรง | 003 ใช้ร่วม ไม่ได้แก้ |
| frontend config.js: ข้อความของ 003 9 ตัว | FR-CXL-01, 02, 03, 06, 08, ASM-06, ASM-13 | ตรง | MSG_CANCEL_FAILED ทีมเลือกคำ (T-24) |

หมายเหตุ (ไม่ใช่ข้อค้นพบของโค้ด): หัวไฟล์ tasks.md ยังอ้าง spec Draft v3.2 / plan Draft v1.1 และ test-cases.md อ้าง spec Draft v3.3 (ปัจจุบัน spec v3.8, plan v1.6)
หมายเหตุ Q-SH-05 และบรรทัด Status ของ spec ที่ /verify รอบแรกบันทึกไว้ แก้แล้วใน spec v3.8 / plan v1.6

## 3. ข้อค้นพบ
ชนิด: AC ไม่มี test / test อ่อน / โค้ดไม่มี FR / FR ไม่มี AC / เดา Q-xx / ละเมิด Constraint / ตัวเลขไม่ตรง spec / อ้าง ID ผิดเรื่อง
ทีมตัดสิน: แก้โค้ด / แก้ spec / เพิ่ม Q-xx / ไม่ใช่ปัญหา (พร้อมเหตุผล 1 บรรทัด)

| F-ID | ชนิด | อยู่ที่ | ขัดกับ | รายละเอียด | ทีมตัดสิน |
|---|---|---|---|---|---|
| F-03 | ละเมิด Constraint | deploy/docker-compose.yml (ระบบทั้งหมด) | CON-TECH-02, NFR-SEC-01 (TLS 1.2 ขึ้นไป) | ไม่มี HTTPS หรือ reverse proxy ในระบบ หลังบ้านฟัง HTTP ที่ :8000 plan ข้อ 5 ของ 003 บันทึกไว้แล้วว่าเป็นงานค้างร่วมกับ 001 ไม่มี task ของใครรับ | |
| F-04 | ตัวเลขไม่ตรง spec | backend/internal/cancelrequest/service.go: Cancel / notify (`context.WithoutCancel`) | NFR-PERF-CXL-01 ("ไม่นับเวลาส่ง Email"), AC-CXL-16 ("ระบบแจ้งเตือน Email ไม่ตอบสนอง") | Email ส่งแบบรอผลก่อนตอบนักศึกษา (ทีมตัดสิน ASM-10) และไม่มีเวลาจำกัดการส่ง ถ้าตัวส่งค้าง นักศึกษาจะรอไม่สิ้นสุด และเวลาตอบสนองของการยกเลิกรวมเวลาส่ง Email ซึ่ง spec บอกว่าไม่นับ ตัวจำลองตอนนี้ไม่ค้างจึงยังไม่เกิด ค้างจาก T-28 ว่าจะเพิ่ม NotifyTimeout ใน config หรือจดเป็นคำถามคู่ Q-15 (รอบ 2: ยังเหมือนเดิม auditNotifyFailed ที่เพิ่มใน T-32 มีเวลาจำกัด dbTimeout แต่ตัวส่ง Email ยังไม่มี) |  แก้โค้ด: เพิ่ม NotifyTimeout 5 วินาที (ทีมเลือกเอง) ใน config.go ส่งไม่ทันนับเป็นส่งไม่สำเร็จ บันทึก notify error ทำใน T-34 spec v3.9 ASM-10 (ทีมตัดสิน 2569-11-12) |
| F-05 | ตัวเลขไม่ตรง spec (ตัวเลขฝังในโค้ด) | tests/load/cancel_request.js: P95_LIMIT_MS = 15000, VUS = 50 | NFR-PERF-CXL-01, AGENTS.md กฎข้อ 9 | ค่าตรงกับ spec แต่ฝังใน script แทนที่จะอ่านจาก config.CancelPerfP95 (script ของ k6 อ่าน config ของ Go ไม่ได้ มีคอมเมนต์บอกให้แก้คู่กัน และเปลี่ยนผ่าน env P95_LIMIT_MS ได้) | ไม่ใช่ปัญหา: k6 อ่าน config.go ไม่ได้ มีคอมเมนต์ผูกกับ config.CancelPerfP95 และเปลี่ยนผ่าน env P95_LIMIT_MS ได้ (ทีมตัดสิน 2569-11-12) |
| F-07 | FR ไม่มี AC | spec.md FR-CXL-10 | | FR-CXL-10 (เจ้าหน้าที่ยกเลิกแทน) ไม่มี AC และไม่มีโค้ด spec เขียนไว้ชัดว่าไม่ทำใน prototype และยังไม่มี AC จนกว่าจะนำเข้า scope (รอ Q-16) | ไม่ใช่ปัญหา: spec ตั้งใจไม่ทำใน prototype จะเขียน AC เมื่อนำเข้า scope (Q-16) (ทีมตัดสิน 2569-11-12) |
| F-08 | FR ไม่มี AC | spec.md NFR-AVAIL-01, NFR-SEC-01 (ส่วน TLS) | | NFR-AVAIL-01 (99% ช่วง 08.30-16.30 น.) และส่วนเข้ารหัส TLS ของ NFR-SEC-01 ไม่มี AC ใน spec 003 วัดผลไม่ได้ใน prototype | |
| F-09 | ตัวเลขไม่ตรง spec (เงื่อนไขไม่ตรง) | backend/internal/cancelrequest/handler.go: cancel (ตอบก่อนเรียก Service) / service.go: auditRejected | FR-CXL-05 ("ทุกครั้งที่นักศึกษาที่ยืนยันตัวตนแล้วพยายามยกเลิกคำขอ" และตาราง: เหตุผลผิด request_id = รหัสคำขอ) | 3 กรณีที่โค้ดต่างจากตาราง: (1) id ที่ไม่ใช่ตัวเลข ตอบ 404 ใน handler ไม่บันทึก audit (2) body ที่ไม่ใช่ JSON ตอบ 422 ใน handler ไม่บันทึก audit (ยิงหลังบ้านจริงแล้ว ได้ 422 ไม่ผ่าน Service) (3) เหตุผลผิดกับรหัสคำขอที่ไม่มีอยู่ บันทึก reason_invalid แบบ request_id ว่างและเก็บ id ใน query เพราะ FK อ้างคำขอที่ไม่มีไม่ได้ ทั้ง 3 ข้อ AI เลือกเองตอน T-31 และเจ้าของรับไว้ตอนตรวจ T-31 (prompt-log) แต่ spec ยังไม่ได้เขียนกรณีเหล่านี้ ทางที่เป็นไปได้: เพิ่มแถวในตาราง FR-CXL-05 (แก้ spec) หรือให้ handler ส่ง 2 กรณีแรกผ่าน Service ให้บันทึก (แก้โค้ด) |  แก้ spec: เพิ่ม 2 แถวในตาราง FR-CXL-05 (เหตุผลผิดกับคำขอที่ไม่มี, id / body อ่านไม่ได้ไม่บันทึก) ผู้ใช้หน้าเว็บส่งแบบนี้ไม่ได้ spec v3.9 ไม่แก้โค้ด (ทีมตัดสิน 2569-11-12) |
| F-10 | อ้าง ID ผิดเรื่อง (คอมเมนต์ไม่ตรงสัญญา API) | frontend/src/api/client.js: cancelRequest (บรรทัด 136) | plan.md ข้อ 4 (CR-006), CON-TECH-03 | คอมเมนต์ยังเขียนว่า throw ApiError(400 REASON_*, ...) แต่หลัง CR-006 หลังบ้านตอบ 422 VALIDATION_FAILED และมี 403 FORBIDDEN_ROLE ด้วย คนที่อ่านคอมเมนต์จะเข้าใจสัญญาผิด ไม่กระทบการทำงาน (หน้าจอตัดสินจาก code) ไฟล์เป็นของกลาง แต่บรรทัดนี้ 003 เขียน (T-07) | แก้โค้ด: แก้คอมเมนต์ client.js บรรทัด 136 เป็น 403 FORBIDDEN_ROLE, 404, 409, 422 VALIDATION_FAILED, 503 และแก้ msgUnauthenticated ใน errors.go ให้ตรงของกลาง (เจ้าของแก้เอง ทีมตัดสิน 2569-11-12) |

## 4. แก้แล้ว
| F-ID | แก้อย่างไร | รู้ได้อย่างไร |
|---|---|---|
| F-01 | CR-006 + T-33: errors.go (รหัสและข้อความ), ErrorDTO {code, message, fields, status}, handler.go ใช้รหัสกลาง, service.go แนบสถานะใน 409 (StatusError) แก้ constraints v5.2, spec v3.8, plan ข้อ 4, test-cases 4 แถว / ทีมตัดสินเดิม: "แก้ spec และโค้ด: CR-006 ลงทะเบียน REQUEST_NOT_FOUND (ร่วม 004), ALREADY_APPROVED, NOT_CANCELLABLE (แนบ status) ใช้ FORBIDDEN_ROLE, VALIDATION_FAILED 422 ตัด retryable ทำใน T-33 (ทีมตัดสิน 2569-11-12)" | /verify รอบ 2: ค้นทั้ง backend/internal และ frontend/src ไม่เหลือ STUDENT_ONLY, REASON_REQUIRED, REASON_TOO_LONG หรือ retryable ในโค้ด (เหลือคอมเมนต์ F-10) test ตรวจ message, fields และ status ผ่าน ยิงหลังบ้านจริงได้ {code, message} ทุกกรณีที่ลอง |
| F-02 | เพิ่ม TestFR_CXL_03_StaffCancelForbidden ใน handler_test.go (เจ้าหน้าที่ POST cancel ได้ 403 สถานะไม่เปลี่ยน หลัง CR-006 ตรวจ FORBIDDEN_ROLE) ไม่แก้โค้ด / ทีมตัดสินเดิม: "แก้โค้ด: เพิ่ม test TestFR_CXL_03_StaffCancelForbidden ไม่แก้โค้ด (ทีมตัดสิน 2569-11-12)" | /verify รอบ 2: test ผ่าน ลองปิดบรรทัดตรวจบทบาทแล้ว test ไม่ผ่าน (2569-11-12) และยิงหลังบ้านจริงด้วย mock-staff-001 ได้ 403 FORBIDDEN_ROLE |
| F-06 | plan ข้อ 4 กำหนด 2 กรณีแล้ว (CR-006) body อ่านไม่ได้ -> 422 VALIDATION_FAILED, id ไม่ใช่ตัวเลข -> 404 REQUEST_NOT_FOUND / ทีมตัดสินเดิม: "แก้ spec: รวมกับ F-01 ใน CR-006 plan ข้อ 4 กำหนดแล้ว body อ่านไม่ได้ -> 422 VALIDATION_FAILED, id ไม่ใช่ตัวเลข -> 404 REQUEST_NOT_FOUND (ทีมตัดสิน 2569-11-12)" | /verify รอบ 2 ยิงหลังบ้านจริง: body "not json" ได้ 422 VALIDATION_FAILED fields [reason], id "abc" ได้ 404 {"code":"REQUEST_NOT_FOUND","message":"ไม่พบคำขอ"} ตรง plan ยังไม่มี test อัตโนมัติของ 2 กรณีนี้ (เรื่อง audit ของ 2 กรณีนี้ดู F-09) |
