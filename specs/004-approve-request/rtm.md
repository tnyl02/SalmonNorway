# RTM: อนุมัติคำขอยืมปกติ (Approve Request)
อ้างอิง: spec.md SPEC-APV-004 Draft v2.2 (CR-008 รวมแก้ไขหลังอนุมัติ 2569-11-13) | plan.md v1.2 | tasks.md (T-01 ถึง T-36 เสร็จ, T-37 ถึง T-39 รอ Q) | test-cases.md (124 แถว ใช้ได้ทั้งหมด)
สร้างด้วย /verify เมื่อ 2569-11-12 / /verify รอบ 2 เมื่อ 2569-11-13 (หลังทีมตัดสิน F-01 ถึง F-06 และ pull main PR #23 ถึง #29) | test: 228 ผ่าน 0 ไม่ผ่าน
(รอบ 2: go test `./internal/approve/ ./internal/config/ ./internal/platform/ ./internal/auth/` -count=1: 191 ผ่าน นับรวม subtest / vitest ของหน้าจอ 004: 37 ผ่าน ใน 19 ไฟล์ / แถว "คน" 2 แถว ผ่าน / vitest ทั้งโปรเจกต์ 184 ผ่าน 1 skip ของ 001)

## 1. ตามรอยไปข้างหน้า (requirement ไป โค้ด ไป test)
FR-APV-01 และ FR-APV-02 แยกแถวตามข้อย่อย เพราะแต่ละข้อย่อยมี AC, โค้ด และ test ต่างกัน

| ID | AC | task | โค้ด (ไฟล์: ฟังก์ชัน) | test (ผล) | สถานะ |
|---|---|---|---|---|---|
| FR-APV-01 (ก) หน้ารายการ ภาควิชาตนเอง เรียงตามเวลาที่ส่ง | AC-APV-11, 12 | T-04, T-05, T-06, T-25 | repository.go: ListPending (กรอง requested, normal, ภาควิชา, ORDER BY submitted_at, id) / service.go: ListPending / handler.go: list / ApproveRequests.jsx | TestAC_APV_11, TestAC_APV_11_HTTP, TestTC_APV_11_1, 3, 4, 6, 7, 8, vitest AC-APV-11, TC-APV-11-2, 11-5 (ผ่าน) | ครบ (F-05 แก้แล้ว รอบ 2) |
| FR-APV-01 (ข) อนุมัติเฉพาะคำขอแรก ค้างคืนไม่ขวางคิว ปฏิเสธได้ทุกคำขอ | AC-APV-07, 16, 17 | T-17, T-18 | service.go: Approve (คำขอแรกจาก lockItemRequests) / repository.go: is_first / ApproveRequests.jsx ปุ่มอนุมัติเฉพาะ isFirst | TestAC_APV_07, TestListPending_OvernightDoesNotBlockQueue, TestTC_APV_01_6, 01_7, 01_9, 02_6, 02_7, 07_3 ถึง 07_6, vitest TC-APV-07-1, 07-2 (ผ่าน) | ครบ |
| FR-APV-01 (ค) อนุมัติเมื่อ requested + "ว่าง" ก่อน 16:00 น. ของวันที่เริ่มยืม | AC-APV-01, 08, 16 | T-07, T-08, T-09, T-19, T-20, T-21, T-22 | service.go: Approve, approvalCutoff (`!Before`, LoanNormalReturnHour) | TestAC_APV_01, 08, 16, TestTC_APV_01_1, 01_4, 01_5, 01_10, 08_1, 08_3, 16_1, 16_3, 16_4 (ผ่าน) | ครบ |
| FR-APV-01 (ง) ปฏิเสธต้องมีเหตุผล ไม่เกิน 500 ตัวอักษร | AC-APV-02, 06 | T-10, T-11, T-12, T-13 | service.go: Reject (TrimSpace, utf8.RuneCountInString) / handler.go: VALIDATION_FAILED / ApproveRequests.jsx: reject | TestAC_APV_02, 06, TestTC_APV_02_1, 06_1, 06_3, 06_4, 06_6, 06_7, vitest TC-APV-02-2, 06-2, 06-5 (ผ่าน) | ครบ |
| FR-APV-01 (จ) ปฏิเสธตามเมื่ออนุมัติ ทั้งปกติและค้างคืน | AC-APV-03, 17 | T-16, T-31 | service.go: Approve (rejectSuperseded) | TestAC_APV_03, 17, TestTC_APV_03_1 ถึง 03_7, 17_1, 17_5, 04_5, vitest TC-APV-03-2 (ผ่าน) | ครบ |
| FR-APV-01 (ฉ) ปฏิเสธอัตโนมัติ 00:00 น. ยืมปกติ นับจากวันที่ส่ง ไล่ตอนเริ่มทำงาน | AC-APV-13, 17, 18 | T-29, T-30, T-31 | autoreject.go: Run, CatchUp, RejectStale, autoRejectCutoff / main.go: go NewAutoRejecter(db).Run | TestAC_APV_13, 18, TestTC_APV_13_1 ถึง 13_6, 17_2 ถึง 17_4, 18_1 ถึง 18_5 (ผ่าน) | ครบ |
| FR-APV-01 (ช) กดพร้อมกัน บันทึกเฉพาะคนแรก | AC-APV-10 | T-23, T-35, T-36 | repository.go: lockItemRequests (FOR UPDATE) / service.go: notRequestedError (loadedAt) | TestAC_APV_10 (-count=20), TestTC_APV_10_1, 10_2, 10_4, 10_5, vitest TC-APV-10-3 (ผ่าน) | ครบ |
| FR-APV-01 (ซ) ไม่เปลี่ยนสถานะ + ข้อความตามกรณี + ลำดับการตรวจ | AC-APV-05, 06, 07, 08, 09, 10, 12, 16 (ไม่มี AC ของ UPSTREAM_UNAVAILABLE และ REQUEST_NOT_FOUND) | T-12, T-17, T-19, T-21, T-23, T-25, T-27 | handler.go: writeServiceError / config.go: ApprovalMessages, ApprovalReasonTooLongMessage / config.js: APPROVAL_MESSAGES | TestAC_APV_09, 12, 05, TestTC_APV_05_*, 09_*, 12_*, vitest TC-APV-05-3, 05-8, 08-2, 09-3, 12-3, 16-2 (ผ่าน) ไม่มี test ของ 503 | ช่องโหว่ (F-04 รอ Q-APV-40) |
| FR-APV-01 (ฌ) ประวัติ: รหัสคำขอ ผล เหตุผล ผู้ดำเนินการ (รหัสและชื่อ) วันเวลา | AC-APV-15 | T-28, T-29 | repository.go: insertDecision, insertAudit / migration 0008: approval_decision (append-only trigger) | TestAC_APV_15, TestTC_APV_15_1 ถึง 15_6 (ผ่าน) ชื่อเจ้าหน้าที่ไม่ได้ assert ตาม ASM-11 | รอ Q-APV-04 (ASM-11 ใช้รหัสแทนชื่อ) |
| FR-APV-02 (ก)(ข) Email ผลอนุมัติ / ปฏิเสธ ภายใน 60 วินาที | AC-APV-04 | T-14, T-15, T-37, T-38 | email.go: BuildEmail, emailValues / sender.go: Sender.Run, SendDue, loadDueEmails / config.go: ApprovalEmailSubjects, ApprovalEmailTemplates | TestAC_APV_04, TestBuildEmail_*, TestTC_APV_04_1 ถึง 04_6 (ผ่าน) ชื่อผู้อนุมัติ/ผู้พิจารณาและผู้รับไม่ได้ assert (ASM-10, ASM-11) | รอ Q-APV-04 (T-37, T-38) |
| FR-APV-02 (ค) Email ปฏิเสธอัตโนมัติ | AC-APV-13, 18 | T-29, T-30 | email.go: emailValues (auto_rejected) / config.go: ApprovalAutoRejectEmail | TestAC_APV_13, 18, TestTC_APV_13_2, 18_2 (ผ่าน) | ครบ |
| FR-APV-02 (ง) ส่งไม่สำเร็จ: Failed + สาเหตุ + เวลา, ซ้ำ 3 ครั้งห่าง 1 นาที, รายการให้กดส่งซ้ำ | AC-APV-14 | T-32, T-33, T-34 | sender.go: markAttempt / service.go: ListFailedEmails, Resend / handler.go: failedEmails, resend / FailedEmails.jsx | TestAC_APV_14, TestSender_RetryThenFailed, TestTC_APV_14_1 ถึง 14_9, vitest AC-APV-14, TC-APV-14-4 (ผ่าน) | ครบ |
| NFR-SEC-01 ยืนยันตัวตนก่อนทุกครั้ง + TLS 1.2 ขึ้นไป | AC-APV-05 (ส่วนยืนยันตัวตนเท่านั้น ไม่มี AC ของ TLS) | T-05, T-27, T-39 | auth.Middleware (ของกลาง) / handler.go: requireStaff / main.go: ผูก route หลัง Middleware | TestAC_APV_05, TestApprovals_RequiresStaff, TestTC_APV_05_1 ถึง 05_7 (ผ่าน) | ช่องโหว่ (F-01, F-02 รอ Q-APV-38, 39) |
| ค่าที่มาจาก spec อยู่ใน config.go (หมวด Quality Requirements, AGENTS.md กฎข้อ 9) | ไม่มี AC | ทุก task | config.go บล็อก 004 / config.js | TestConfig_Approval* ใน config_approval_test.go (ผ่าน) | ครบ (F-06 แก้แล้ว รอบ 2) |
| CON-TECH-01 PostgreSQL | ไม่มี AC | T-01 | migration 0008_approval.up.sql / down.sql | TestMigration_ApprovalTables, TestMigration_ApprovalAuditActions, TestSeed_ApproveDevScenarios (ผ่าน) | ครบ |
| CON-TECH-02 HTTPS เท่านั้น deploy ด้วย Docker | ไม่มี AC | T-05 | ไม่มี: deploy/docker-compose.yml มีแค่ service `db` หลังบ้านรัน HTTP ที่ :8000 | ไม่มี test | ช่องโหว่ (F-01 รอ Q-APV-38) |
| CON-TECH-03 error `{code, message}` จากรายการกลาง | AC-APV-05, 06, 07, 08, 09, 10, 12, 16 | T-02, T-05, T-12, T-17, T-19, T-23, T-25, T-27, T-35 | handler.go: writeError, writeServiceError / platform: WriteError, MethodNotAllowedJSON | TestTC_APV_05_1 (message ไม่ว่าง), 12_1 (body มีแค่ code, message) และทุก test ที่ตรวจรหัส (ผ่าน) | ครบ |
| CON-DATA-01 สถานะอุปกรณ์ 3 ค่า อนุมัติแล้วยัง "ว่าง" | AC-APV-01, 02, 08 | T-07, T-10, T-19 | service.go: Approve, Reject (ไม่แตะ equipment_item) | TestAC_APV_01, 02, TestTC_APV_01_1, 02_1, 02_4, 08_1 (ผ่าน) | ครบ |
| DOM-REQ-01 คำขอซ้อนได้ คำขอที่อนุมัติก่อนได้อุปกรณ์ | AC-APV-03, 07 | T-04, T-07, T-16, T-17 | service.go: Approve (คำขอแรก + rejectSuperseded) | TestAC_APV_03, 07, TestTC_APV_03_*, 07_* (ผ่าน) | ครบ |
| DOM-LOAN-01 ยืมปกติคืนภายใน 16:00 น. | AC-APV-16 | T-02, T-21 | service.go: approvalCutoff ใช้ config.LoanNormalReturnHour ของ 002 | TestAC_APV_16, TestTC_APV_16_*, 01_4, 01_5 (ผ่าน) | ครบ |
| IF-IDP-01 ยืนยันตัวตนผ่านระบบมหาวิทยาลัย ได้บทบาทและภาควิชา | AC-APV-05, 12 | T-05, T-25, T-27, T-39 | auth.NewMockProvider (จำลอง ASM-02) / handler.go: requireStaff / service.go: staff() | TestAC_APV_05, 12, TestTC_APV_05_*, 12_* (ผ่าน กับ auth จำลอง) | รอ Q-SH-03 (T-39) |
| IF-NOT-01 แจ้งผ่านระบบ Email | AC-APV-04, 13, 14 | T-14, T-15, T-32, T-37 | mailer.go: LogMailer, MockRecipientResolver (Email แบบจำลอง ASM-10) | TestAC_APV_04, TestMockRecipientResolver, TestTC_APV_04_*, 14_* (ผ่าน กับตัวส่งจำลอง) | รอ Q-APV-04 (T-37, T-38) |

## 2. ตามรอยย้อนกลับ (โค้ด ไป requirement)
| โค้ด (ไฟล์: ฟังก์ชัน หรือ endpoint) | อ้าง ID | ตรงกับข้อความใน spec ไหม | หมายเหตุ |
|---|---|---|---|
| GET /v1/approvals/requests (handler.go: list) | FR-APV-01 ก ข, AC-APV-11 | ตรง | รอบ 2: ไม่มี `studentCode` แล้ว (F-05 แก้แล้ว) / `isFirst`, `loadedAt`, `submittedAt` มาจาก plan ข้อ 4 |
| POST /v1/approvals/requests/{id}/approve (handler.go: approve) | FR-APV-01 ข ค จ ช, AC-APV-01 | ตรง | ลำดับการตรวจตรงกับ FR-APV-01 (ซ) Q-APV-35, 36 |
| POST /v1/approvals/requests/{id}/reject (handler.go: reject) | FR-APV-01 ง, AC-APV-02 | ตรง | ลำดับการตรวจตรงกับ Q-APV-29, 34 |
| GET /v1/approvals/failed-emails (handler.go: failedEmails) | FR-APV-02 ง, AC-APV-14 | ตรง | รวม Email ของ System ตาม plan ข้อ 4 |
| POST /v1/approvals/emails/{id}/resend (handler.go: resend) | FR-APV-02 ง, AC-APV-14 | ตรง | ส่งซ้ำ Email ที่ sent แล้วตอบ sent เดิม (Q-APV-37) |
| handler.go: requireStaff | AC-APV-05, NFR-SEC-01 | ตรง | ครอบทุก route ของ 004 |
| handler.go: writeServiceError (default -> 503 UPSTREAM_UNAVAILABLE) | FR-APV-01 ซ, Q-APV-24 | ตรง | ไม่มี AC และไม่มี test ของกรณีนี้ (F-04) / log ข้อความ error ลง log ของ server |
| service.go: Approve, Reject, ListPending, ListFailedEmails, Resend | FR-APV-01, FR-APV-02 ง | ตรง | รอบ 2: ข้อความใน `ErrApprovalCutoff` ไม่มีตัวเลขแล้ว (F-06 แก้แล้ว) |
| repository.go: lockItemRequests, insertDecision, insertAudit, insertEmail, requeueEmail | FR-APV-01 ช ฌ, RQ-25 | ตรง | insertDecision ของเจ้าหน้าที่ไม่มี actor_name ตาม ASM-11 |
| sender.go: Sender.Run, SendDue, markAttempt | FR-APV-02, FR-APV-02 ง | ตรงกับ ASM-11 | loadDueEmails ใช้ `coalesce(actor_name, actor_id)` ชื่อใน Email จึงเป็นรหัส เช่น "staff-001" ตาม ASM-11 / ค่ารอบตรวจ 5 วินาที และ 50 ฉบับต่อรอบ ทีมเลือกเอง ระบุใน config.go แล้ว |
| autoreject.go: Run, CatchUp, RejectStale | FR-APV-01 ฉ, Q-APV-14, 15, 21 | ตรง | |
| email.go: BuildEmail, emailValues, formatDateTime, formatDate | FR-APV-02 ก ข ค | ตรง | รูปแบบวันที่ พ.ศ. ทีมตัดสิน 2569-11-12 (เหมือน 003) |
| mailer.go: LogMailer, MockRecipientResolver | IF-NOT-01, ASM-10, Q-APV-04 | ตรงกับ ASM-10 (Email แบบจำลอง) | LogMailer เขียนเนื้อ Email ทั้งฉบับ (ชื่อผู้ขอ อุปกรณ์ เหตุผล) ลง log ของ server ไม่ขัด Constraint ของ 004 แต่ทีมควรรู้ |
| migration 0008: approval_decision, approval_email, audit_log action | FR-APV-01 ฌ, FR-APV-02 ง, RQ-25 | ตรง | approval_decision แก้และลบไม่ได้ (TC-APV-15-5) |
| config.go บล็อก 004 / config.js บล็อก 004 | Quality Requirements, AGENTS.md กฎข้อ 9 | ตรง | ข้อความ Email และ error อยู่ใน config ครบ |
| main.go: ผูก approve.PathPrefix หลัง auth.Middleware, go Sender.Run, go AutoRejecter.Run | IF-IDP-01, FR-APV-01 ฉ, FR-APV-02 | ตรง | server เปิด HTTP ธรรมดา (F-01) |
| ApproveRequests.jsx | FR-APV-01 ก ข ค ง ซ | ตรง | รอบ 2: ไม่มีคอลัมน์ "เวลาที่ส่ง" แล้ว (F-05 แก้แล้ว) / ไม่มี mockup ทีมออกแบบเอง |
| FailedEmails.jsx | FR-APV-02 ง, ASM-08 | ตรง | ป้ายแบบ Email ทีมเลือกเอง ระบุใน config.js แล้ว |
| api/client.js: getApprovalRequests, approveRequest, rejectRequest, getFailedEmails, resendEmail | plan ข้อ 4 | ตรง | |
| approve.html, approve-main.jsx | plan ข้อ 4 | ตรง | หน้าแยกของ 004 ใช้ token จำลองจาก VITE_MOCK_TOKEN (รอ Q-SH-03) |
| backend/db/seed/approve_dev.sql | ASM-01 | ตรง | ข้อมูลทดลองของเครื่องพัฒนาเท่านั้น |

## 3. ข้อค้นพบ
ชนิด: AC ไม่มี test / test อ่อน / โค้ดไม่มี FR / FR ไม่มี AC / เดา Q-xx / ละเมิด Constraint / ตัวเลขไม่ตรง spec / อ้าง ID ผิดเรื่อง
ทีมตัดสิน: แก้โค้ด / แก้ spec / เพิ่ม Q-xx / ไม่ใช่ปัญหา (พร้อมเหตุผล 1 บรรทัด)

| F-ID | ชนิด | อยู่ที่ | ขัดกับ | รายละเอียด | ทีมตัดสิน |
|---|---|---|---|---|---|
| F-01 | ละเมิด Constraint | deploy/docker-compose.yml, backend/cmd/server/main.go (`http.ListenAndServe`), tasks.md บรรทัดตรวจ CON-TECH-02 | CON-TECH-02, NFR-SEC-01 (TLS 1.2) | compose มีแค่ service `db` หลังบ้านไม่ได้ deploy ด้วย Docker และเปิด HTTP ธรรมดาที่ :8000 ไม่มี HTTPS / tasks.md เขียนว่า "route อยู่ใน server เดิมหลัง proxy HTTPS ของกลาง" แต่ใน repo ไม่มี proxy นั้น เป็นเรื่องระดับทีม (003 พบเรื่องเดียวกัน) ทางเลือก: เพิ่ม reverse proxy HTTPS ใน compose หรือเขียน Assumption ว่า prototype ไม่ทำ HTTPS | เพิ่ม Q-xx: ให้ทีมตัดสินรวมเป็นคำถามระดับทีม (ตัดสิน 2569-11-13) |
| F-02 | FR ไม่มี AC | spec.md หมวด Quality Requirements และ Acceptance Criteria | NFR-SEC-01 | ส่วน "เข้ารหัสขณะรับส่งด้วย TLS 1.2 ขึ้นไป" ไม่มี AC ใดตรวจ AC-APV-05 ตรวจแค่การยืนยันตัวตน (003 พบเรื่องเดียวกัน) ถ้าจะเพิ่ม AC ต้องเสนอเป็นคำถามก่อน | เพิ่ม Q-xx: เสนอ AC ใหม่สำหรับ TLS ให้ทีมอนุมัติ (ตัดสิน 2569-11-13) |
| F-04 | FR ไม่มี AC | spec.md FR-APV-01 (ซ) แถว `UPSTREAM_UNAVAILABLE` และ `REQUEST_NOT_FOUND` / handler.go: writeServiceError | FR-APV-01 (ซ) | 2 กรณีในตารางข้อความไม่มี AC รองรับ `REQUEST_NOT_FOUND` มี test (TC-APV-01-12, 02-11, TestApproveHTTP_NotFound) แต่ `UPSTREAM_UNAVAILABLE` (ฐานข้อมูลไม่ตอบ 503) ไม่มีทั้ง AC และ test ถ้าลบบรรทัด default ใน writeServiceError ไม่มี test ไหนไม่ผ่าน | เพิ่ม Q-xx: เสนอ AC ใหม่ของ 503 UPSTREAM_UNAVAILABLE และ 404 REQUEST_NOT_FOUND ให้ทีมอนุมัติ แล้วทำ test-cases ตาม (ตัดสิน 2569-11-13) |

## 4. แก้แล้ว
| F-ID | แก้อย่างไร | รู้ได้อย่างไร |
|---|---|---|
| F-03 | ทีมตัดสินแก้ spec: เพิ่ม ASM-11 ใน spec Draft v2.2 (CR-008 แก้ไขหลังอนุมัติ 2569-11-13) ว่า prototype ใช้รหัสเจ้าหน้าที่แทนชื่อจนกว่าจะได้คำตอบ Q-APV-04 / Q-SH-03 | spec.md มี ASM-11 และ CR-008 มีข้อ 13 / โค้ดทำตาม ASM-11 อยู่แล้ว (sender.go coalesce) จึงไม่ใช่การเดาอีกต่อไป ส่วนชื่อจริงยังรอ Q-APV-04 (แถว FR-APV-02 ก ข และ FR-APV-01 ฌ เป็น "รอ Q-APV-04") |
| F-05 | ทีมตัดสินแก้โค้ด: ตัด `studentCode` ออกจาก PendingRequestDTO (dto.go) และตัดคอลัมน์ "เวลาที่ส่ง" ออกจาก ApproveRequests.jsx | อ่านโค้ด: ไม่มี `studentCode` ใน dto.go และไม่มีคอลัมน์ "เวลาที่ส่ง" ใน ApproveRequests.jsx / test ของ 004 ผ่านหลังแก้ (go test approve, vitest 37) |
| F-06 | ทีมตัดสินแก้โค้ด: ข้อความของ `ErrApprovalCutoff` ใน service.go เปลี่ยนเป็น "เลยเวลาตัดของวันที่เริ่มยืมแล้ว อนุมัติไม่ได้" ไม่มีตัวเลข | อ่านโค้ด: "16:00" ใน service.go เหลือแค่ในคอมเมนต์ / ข้อความที่ผู้ใช้เห็นยังมาจาก config (TestTC_APV_16_1 ผ่าน) |
