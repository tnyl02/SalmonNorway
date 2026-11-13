# RTM: บันทึกการส่งมอบอุปกรณ์ (Handover)
อ้างอิง: spec.md SPEC-HND-005 Draft v2.10 | plan.md v1.5 | tasks.md v1.5 | test-cases.md (63 แถว สถานะ "ร่าง" ทั้งหมด ยังไม่มี test `TestTC_`)
สร้างด้วย /verify เมื่อ 2569-11-13 22:15 (Asia/Bangkok) ที่ commit e550ddd บน main (รอบที่ 2 / รอบ 1 2569-11-12 22:29 ที่ e57cd3f) | test: 32 ผ่าน 0 ไม่ผ่าน (ของ 005: หลังบ้าน 25 test function ใน `backend/internal/handover/` และหน้าจอ 7 test ใน `AC-HND-02/12/14.test.jsx`)
ทั้งระบบ: `go test ./... -v -count=1` 539 ผ่าน 0 ไม่ผ่าน SKIP 2 (ของ 001) / `npx vitest run` 213 ผ่าน 0 ไม่ผ่าน SKIP 1 (ของ 001)
รอบนี้รันโดย 001 (พรรณธร) ตามคำสั่งในกลุ่มก่อนส่ง โค้ดใน `backend/internal/handover/` และหน้าจอของ 005 ไม่เปลี่ยนตั้งแต่รอบ 1
สิ่งที่เปลี่ยนรอบนี้: spec v2.9 (Q-SH-08 ปิดสำหรับ prototype), HTTPS ของ 004 ใน deploy/docker-compose.yml (กระทบ F-04, F-05), Q-SH-06 ปิดในไฟล์กลาง (กระทบ AC-HND-13 ดู F-09)
branch `005-handover` (commit 9d9212f เพิ่ม volume `photos`) ยังไม่ได้ merge เข้า main (PR #38 ถูกปิดเพราะ conflict) รอบนี้ตรวจตาม main
**อัปเดตรอบ 3 (2569-11-13 23:00 หลัง merge `005-handover` ที่ a0c1b5a):** docker-compose มี `PHOTO_DIR` และ volume `photos` แล้ว F-05 ย้ายไป "แก้แล้ว" / spec v2.10 แก้สถานะ Q-SH-06 / ไม่มีโค้ด Go หรือหน้าจอเปลี่ยน ไม่ได้รัน test ซ้ำ
**หมายเหตุ spec v2.11 (2569-11-13):** ผล test ในไฟล์นี้อ้าง spec Draft v2.10 ส่วน v2.11 เพิ่มแค่หัวข้อ Traceability ไม่เปลี่ยน FR หรือ AC จึงไม่ได้รัน test ซ้ำ

## 1. ตามรอยไปข้างหน้า (requirement ไป โค้ด ไป test)
| ID | AC | task | โค้ด (ไฟล์: ฟังก์ชัน) | test (ผล) | สถานะ |
|---|---|---|---|---|---|
| FR-HND-01 | AC-HND-01, AC-HND-14 | T-04, T-05, T-16 | repository.go: `ListReadyRequests`, `CreateHandover` (ตรวจ `readyStatuses`); handler.go: `listRequests`; config.go: `HandoverReadyStatuses`; HandoverList.jsx | TestAC_HND_01_RequestNotApproved (ผ่าน), TestAC_HND_14_ListApprovedOnly (ผ่าน), AC-HND-14.test.jsx (ผ่าน) | ช่องโหว่ (F-03 ส่วนค้างคืนจาก UC-14 ไม่มี AC) |
| FR-HND-02 | AC-HND-02, AC-HND-03, AC-HND-10 | T-08, T-09, T-10, T-12, T-13, T-20 | verifier.go: `MockBorrowerVerifier`, `WithTimeout`; service.go: `identifyBorrower`; repository.go: `CreateHandover` (เทียบ `requester_id`); handler.go: `verifyBorrower`; HandoverForm.jsx | TestAC_HND_02_AuthApiVerifiedAndCardChecked, TestAC_HND_03_AuthApiTimeoutCardOnly, TestAC_HND_10_BorrowerMismatch, TestVerifyBorrower_OtherResults, TestVerifier_Mock, AC-HND-02.test.jsx (ผ่านทั้งหมด) | ช่องโหว่ (F-06, F-08) |
| FR-HND-03 | AC-HND-02, AC-HND-03 | T-10, T-11 | service.go: `handover` (`CardChecked`, ชื่อ นามสกุล); repository.go: INSERT `handover`; 0006_handover.up.sql (`card_checked CHECK`) | TestAC_HND_02_HandoverRecordsBorrower, TestAC_HND_03_AuthApiTimeoutCardOnly, TestMigration_HandoverIdentityMethod, AC-HND-02.test.jsx (ผ่าน) | ครบ |
| FR-HND-04 | AC-HND-04, AC-HND-13 | T-06, T-17, T-18 | repository.go: `CreateHandover` (`SELECT ... FOR UPDATE` ตรวจ `available`) | TestAC_HND_04_ItemNotAvailable, TestHandover_ConcurrentSameItem, TestRepository_CreateHandover, TestAC_HND_13_OtherRequestUntouched (ผ่าน) | ครบ |
| FR-HND-05 | AC-HND-05, AC-HND-11 | T-14, T-25, T-30 | service.go: `dueAt`, `normalDueAt`, `overnightDueAt`; config.go: `LoanNormalReturnHour`, `LoanOvernightReturnHour`, `TimeZone` | TestAC_HND_05_DueAt1600SameDay, TestAC_HND_11_HandoverAfter1600, TestHandover_OvernightDueAt (ผ่าน) | ช่องโหว่ (F-03 ส่วนค้างคืนไม่มี AC) |
| FR-HND-06 | AC-HND-06, AC-HND-12 | T-07, T-15, T-21, T-22, T-28 | photos.go: `validatePhotos`, `inspectPhoto`, `photoSaver`; photostore.go: `PhotoStore`; config.go: `HandoverMinPhotos`, `HandoverPhotoMaxBytes`, `HandoverPhotoTypes`; HandoverForm.jsx: `photosValid` | TestAC_HND_06_TwoJpgPhotosSaved, TestAC_HND_12_InvalidPhotos, TestPhotoStore_SaveAndRemove, TestMigration_HandoverPhotoContentType, AC-HND-12.test.jsx (ผ่าน) | ช่องโหว่ (F-03 ส่วนค้างคืน, F-07) |
| FR-HND-07 | AC-HND-07, AC-HND-13 | T-06, T-10, T-19 | repository.go: `CreateHandover` (INSERT `loan`, UPDATE `equipment_item`, UPDATE เฉพาะคำขอนี้) | TestAC_HND_07_HandoverCreatesLoan, TestAC_HND_13_OtherRequestUntouched, TestRepository_CreateHandover (ผ่าน) | ครบ (ดู F-09: Q-SH-06 ปิดแล้ว Given ของ AC-HND-13 เกิดในระบบจริงไม่ได้) |
| FR-HND-08 | AC-HND-08 | T-26, T-27 | repository.go: `audit.Write` ใน transaction (ok); service.go: `Handover`, `AuditFailure` (error); handler.go: `requireStaff` (forbidden); 0009_audit_log_handover_action.up.sql | TestAC_HND_08_AuditLogged 5 กรณี (ผ่าน) | ครบ |
| NFR-SEC-01 | ไม่มี AC | ไม่มี task | ทุก route อยู่หลัง `auth.Middleware` (main.go) / TLS 1.2 ขึ้นไปของทั้งระบบ: platform/tls.go (T-50 ของ 004) | TestHandover_NonStaffForbidden (401 เมื่อไม่มี token ผ่าน) / TLS: tls_test.go, TestTC_APV_19_* ของ 004 (ผ่าน) | ช่องโหว่ (F-04 ไม่มี AC ของ 005) |
| NFR-AVAIL-01 | ไม่มี AC | ไม่มี task | ไม่มี | ไม่มี | ช่องโหว่ (F-04) |
| CON-TECH-01 | ไม่มี AC ตรง | T-02, T-26 | 0006_handover.up.sql, 0009_audit_log_handover_action.up.sql | TestMigration_* 3 ตัว (ผ่าน) | ครบ |
| CON-TECH-02 | ไม่มี AC ตรง | T-03, T-07 | photostore.go (`config.PhotoDir()` อ่าน env `PHOTO_DIR`) / docker-compose service backend แบบ HTTPS (004) ตั้ง `PHOTO_DIR: /data/photos` และ volume `photos:/data/photos` (T-07) | HTTPS: TestTC_APV_19_* ของ 004 (ผ่าน) / volume: ตรวจด้วยการอ่าน deploy/docker-compose.yml และ `docker compose config` ผ่าน ไม่มี test อัตโนมัติ | ครบ (F-05 แก้แล้ว) |
| CON-TECH-03 | AC-HND-01, 04, 09, 10, 12 | T-03, T-10, T-16, T-20, T-21 | errors.go: `writeError`, `errorMessages`; handler.go: `handover`, `verifyBorrower` | TestHandover_NonStaffForbidden, TestAC_HND_01/04/09/10/12 (ตรวจ code ผ่าน) | ช่องโหว่ (F-01) |
| CON-DATA-01 | AC-HND-04, AC-HND-07 | T-06, T-17 | repository.go: ค่าคงที่ `itemStatusAvailable`, `itemStatusBorrowed` | TestAC_HND_04_ItemNotAvailable, TestAC_HND_07_HandoverCreatesLoan (ผ่าน) | ครบ |
| CON-DATA-03 | AC-HND-07, AC-HND-14 | T-04, T-26 | repository.go: `ListReadyRequests` (`assetCode`), `CreateHandover` (`AssetCode`), audit `asset_code` | TestAC_HND_14_ListApprovedOnly, TestAC_HND_07_HandoverCreatesLoan, TestAC_HND_08_AuditLogged (ผ่าน) | ครบ |
| DOM-AUDIT-01 | AC-HND-08 | T-26, T-27 | ดู FR-HND-08 / `listRequests` และ `verifyBorrower` ไม่เขียน audit | TestAC_HND_08_AuditLogged (ผ่าน) | ช่องโหว่ (F-02) |
| DOM-LOAN-01 | AC-HND-05, AC-HND-11 | T-01, T-14, T-25 (T-31 ยกเลิก Q-HND-04) | service.go: `normalDueAt` | TestAC_HND_05_DueAt1600SameDay, TestAC_HND_11_HandoverAfter1600 (ผ่าน) | ครบ |
| DOM-REQ-01 | AC-HND-07, AC-HND-13 | T-06, T-18, T-19 | repository.go: `CreateHandover` (เปลี่ยนเป็น `borrowed` ที่นี่ที่เดียว ไม่แตะคำขออื่น) | TestHandover_ConcurrentSameItem, TestAC_HND_13_OtherRequestUntouched (ผ่าน) | ครบ |
| IF-IDP-01 | AC-HND-02, AC-HND-03, AC-HND-09 | T-03, T-08, T-09, T-13 (T-29 ยกเลิกเทอมนี้) | main.go: `auth.Middleware`; verifier.go: `MockBorrowerVerifier` (ตัวจำลอง ทีมตัดสินเทอมนี้ไม่เชื่อม Auth API จริง) / Q-SH-08 ปิดสำหรับ prototype (ตอบ 401) | TestHandover_NonStaffForbidden, TestAC_HND_02/03 (ผ่าน) | ครบ |
| CON-HND-01 | AC-HND-09 | T-03, T-23 | handler.go: `requireStaff` | TestHandover_NonStaffForbidden, TestAC_HND_09_StudentForbidden (ผ่าน) | ครบ |

สรุป 20 แถว: ครบ 12 / ยังไม่ถึง 0 / รอ 0 / ช่องโหว่ 8 (FR-HND-01, 02, 05, 06, NFR-SEC-01, NFR-AVAIL-01, CON-TECH-03, DOM-AUDIT-01) (รอบ 1: ครบ 11 / ช่องโหว่ 9)
ช่องทีมตัดสินของข้อค้นพบยังว่างทั้งหมด (F-01 ถึง F-04, F-06 ถึง F-09) เจ้าของ 005 ต้องเขียนเอง

## 2. ตามรอยย้อนกลับ (โค้ด ไป requirement)
| โค้ด (ไฟล์: ฟังก์ชัน หรือ endpoint) | อ้าง ID | ตรงกับข้อความใน spec ไหม | หมายเหตุ |
|---|---|---|---|
| `GET /v1/handover/requests` handler.go: `listRequests` | FR-HND-01 | ตรง | ไม่เขียน audit (ดู F-02) |
| `POST /v1/handover/requests/{id}/verify-borrower` handler.go: `verifyBorrower` | FR-HND-02, IF-IDP-01 | ตรง (plan ข้อ 4) | credential ใช้ไม่ได้ตอบ 422 ไม่มีใน plan (F-06) ไม่มีคำขอตอบ 404 ไม่มี `{code, message}` (F-01) |
| `POST /v1/handover/requests/{id}/handover` handler.go: `handover` | FR-HND-02 ถึง 08 | ตรง | ไม่มีคำขอตอบ 404 ไม่มี `{code, message}` (F-01) |
| handler.go: `requireStaff` | CON-HND-01, FR-HND-08 | ตรง | เขียน audit forbidden เฉพาะ route บันทึกการส่งมอบ (plan ข้อ 8 Q-HND-14) |
| handler.go: `requestIDFromPath`, `photoInputs`, `multipartMaxMemory` (8 MB) | - | ไม่มี FR (รายละเอียดทางเทคนิค) | 8 MB คือหน่วยความจำตอนแยกฟอร์ม ทีมเลือกเอง ไม่เปลี่ยนสิ่งที่ผู้ใช้ได้ ไม่นับเป็นข้อค้นพบ |
| service.go: `Handover`, `handover` | FR-HND-02, 03, 05, 07 | ตรง | ตรวจชื่อ นามสกุล credential ว่างแล้วตอบ 422 ช่องที่ขาด plan ไม่ได้ระบุ (F-06) |
| service.go: `AuditFailure` | FR-HND-08, DOM-AUDIT-01, Q-HND-14 | ตรง | คำขอที่ไม่มีจริง บันทึกโดยไม่มี request_id และ asset_code (spec ไม่ได้บอก ดู test-cases TC-HND-08-5) |
| service.go: `identifyBorrower` | FR-HND-02, FR-HND-03, IF-IDP-01 | ตรง | card_only ใช้ `requester_id` ตาม spec Clarifications 2569-11-12 |
| service.go: `dueAt`, `normalDueAt`, `overnightDueAt` | FR-HND-05, DOM-LOAN-01, Q-SH-04 | ตรง | ค้างคืนหลังวันสิ้นสุด ใช้ 16:00 ของวันสิ้นสุดเสมอ ตามที่ทีมสั่งระหว่างรอคำถามที่ยังไม่มี ID ใน spec (ไม่นับเป็นการเดา) |
| verifier.go: `MockBorrowerVerifier`, `WithTimeout`, `timeoutVerifier` | FR-HND-02, IF-IDP-01, ASM-HND-02 | ตรง | timeout อ่านจาก `config.BorrowerAuthTimeout` แต่ไม่มี test ยืนยันค่า 10 วินาที (F-08) |
| photos.go: `validatePhotos`, `inspectPhoto`, `photoSaver` | FR-HND-06 | ตรง | ตรวจชนิดจากเนื้อไฟล์ และขนาดต่อรูปจาก config |
| photostore.go: `PhotoStore`, `photoExt` | FR-HND-06, CON-TECH-02 | ตรง | รายการชนิด JPG/PNG ประกาศซ้ำใน `photoExt` นอกจาก `config.HandoverPhotoTypes` (F-07) |
| repository.go: `CreateHandover` | FR-HND-01, 02, 04, 07, 08, CON-DATA-01, DOM-REQ-01 | ตรง | audit ของกรณีสำเร็จเขียนใน transaction เดียวกัน |
| repository.go: `ListReadyRequests`, `GetRequester` | FR-HND-01, FR-HND-02 | ตรง | ไม่จำกัดภาควิชาของเจ้าหน้าที่ (spec 005 ไม่ได้สั่ง ดูหมายเหตุท้ายตาราง) |
| errors.go: `errorMessages` | CON-TECH-03 | ตรง | ข้อความภาษาไทยทีมเลือกเอง รหัสอยู่ในรายการกลางทุกตัว |
| config.go บล็อก 005 | FR-HND-02, 05, 06 | ตรง | `LoanOvernightReturnHour` อ้าง `LoanNormalReturnHour` ไม่ประกาศเลขซ้ำ |
| 0006_handover.up.sql | FR-HND-02, 03, 05, 06, 07, CON-TECH-01 | ตรง | ไม่มีคอลัมน์เลขบัตร (FR-HND-03) ตรวจด้วย TestMigration_HandoverNoCardNumberColumn |
| 0009_audit_log_handover_action.up.sql | FR-HND-08, Q-HND-06 | ตรง | CHECK มีค่าเดิมของ 001 ถึง 004 ครบ |
| frontend: HandoverList.jsx | FR-HND-01 | ตรง | ข้อความ "ไม่มีคำขอที่รอส่งมอบ" และ "ดึงรายการคำขอไม่สำเร็จ" ทีมออกแบบเอง (ไม่มี mockup ทีมดูหน้าจอแล้ว) |
| frontend: HandoverForm.jsx | FR-HND-02, 03, 05, 06, 07 | ตรง | ข้อความบนจอทีมออกแบบเอง (ทีมดูหน้าจอแล้ว T-11, T-12, T-22, T-28) |
| frontend: HandoverPage.jsx, HandoverApp.jsx, handover-main.jsx, handover.html, api/handover.js | FR-HND-01 ถึง 07 | ตรง | หน้าแยกที่ `/handover.html` ไม่มีลิงก์จากหน้าอื่น และไม่อยู่ใน `vite build` (เรื่องของไฟล์กลาง `vite.config.js` ไม่นับเป็นข้อค้นพบของ 005) |
| backend/db/seed/handover_dev.sql | - | ข้อมูลลองระบบ ทีมเลือกเอง (tasks T-02) | ไม่ใช่โค้ดที่ผู้ใช้ใช้ |

เอกสารที่ยังไม่ครบ (ไม่ใช่ข้อค้นพบของโค้ด แจ้งไว้ให้เจ้าของแก้): spec.md ไม่มีหัวข้อ `## Traceability` (ฟีเจอร์อื่นมีครบ) / test-cases.md 63 แถวยังเป็น "ร่าง" จึงยังไม่นับเป็น "AC ไม่มี test" / spec บรรทัด 8 ยังไม่ได้เทียบกับ US-008 และ UC-08 บน Miro
หมายเหตุ (ไม่ใช่ข้อค้นพบของโค้ด 005): 001, 002, 004 จำกัดเจ้าหน้าที่หรือนักศึกษาตามภาควิชา (`FORBIDDEN_DEPARTMENT`) แต่ spec 005 ไม่ได้สั่ง ยังเปิดอยู่ใน plan ข้อ 8 Q-P2 ทีมควรตัดสินว่าต้องเพิ่มความต้องการไหม

## 3. ข้อค้นพบ
ชนิด: AC ไม่มี test / test อ่อน / โค้ดไม่มี FR / FR ไม่มี AC / เดา Q-xx / ละเมิด Constraint / ตัวเลขไม่ตรง spec / อ้าง ID ผิดเรื่อง
ทีมตัดสิน: แก้โค้ด / แก้ spec / เพิ่ม Q-xx / ไม่ใช่ปัญหา (พร้อมเหตุผล 1 บรรทัด)

| F-ID | ชนิด | อยู่ที่ | ขัดกับ | รายละเอียด | ทีมตัดสิน |
|---|---|---|---|---|---|
| F-01 | ละเมิด Constraint | handler.go: `handover` (`ErrRequestNotFound` -> `http.NotFound`), `verifyBorrower` (ไม่มีคำขอ), `requestIDFrom` ไม่ใช่ตัวเลข | CON-TECH-03 | ไม่มีคำขอรหัสที่เรียก ตอบ 404 เป็นข้อความธรรมดาของ Go ไม่ใช่ `{code, message}` ทั้งที่ CON-TECH-03 ให้ทุก endpoint ตอบ error แบบนี้ และรายการกลางมี `REQUEST_NOT_FOUND` (404) อยู่แล้ว (plan ข้อ 4 ไม่มีแถวของกรณีนี้) ทางแก้ที่น่าจะเป็น: ใช้ `writeError(w, 404, "REQUEST_NOT_FOUND")` และลงชื่อ 005 ในคอลัมน์ "ใครใช้แล้ว" ของรายการกลาง ต้องเพิ่มแถวใน plan ข้อ 4 ก่อน | ไม่ใช่ปัญหา: ผู้ใช้ไม่เห็น 404 นี้ เพราะหน้าจอส่งเฉพาะ requestId จากรายการรอส่งมอบ แนวเดียวกับ 001 F-07 (2569-11-13) |
| F-02 | ละเมิด Constraint | handler.go: `listRequests` (ส่งรหัสครุภัณฑ์และรุ่นของอุปกรณ์), `verifyBorrower` | DOM-AUDIT-01 | DOM-AUDIT-01 ให้บันทึก audit log ทุกครั้งที่เข้าถึงข้อมูลอุปกรณ์ แต่รายการรอส่งมอบแสดงข้อมูลอุปกรณ์โดยไม่บันทึก audit ทีมตัดสินตอน review T-27 (2569-11-12) ว่า "เรียกดูรายการหรือยืนยันตัวตนผู้ยืมไม่ต้องเก็บ" ซึ่งเป็นการตัดสินใน prompt-log / plan ข้อ 8 แต่ขัดกับ Constraint ใช้ร่วม ต้องเลือก: แก้โค้ดให้บันทึก (ค่า action ต้องมีใน CHECK) หรือทำ CR ยกเว้นใน constraints.md |เพิ่ม Q-xx: ใช้ Q-SH-10 ของกลาง (CR-010) ร่วมกับ 002 ว่าการอ่านข้อมูลอุปกรณ์ต้องมี audit แยกไหม ระหว่างรอคงพฤติกรรมเดิม (2569-11-13) |
| F-03 | FR ไม่มี AC | spec.md FR-HND-01 (ค้างคืนจาก UC-14), FR-HND-05 (ค้างคืน = 16:00 ของวันสิ้นสุด), FR-HND-06 (ทั้งยืมปกติและค้างคืน) | spec AC | ส่วนยืมค้างคืนของ 3 FR ไม่มี AC เลย (T-30 ทำแล้ว ตรวจด้วย `TestHandover_OvernightDueAt` ซึ่งไม่ได้ผูกกับ AC) และกรณีส่งมอบค้างคืนหลังวันสิ้นสุดการยืมยังเป็นคำถามที่ยังไม่มี ID ใน spec ควรเพิ่ม AC ของค้างคืนผ่าน CR และตั้ง ID ให้คำถาม |ไม่ใช่ปัญหาใน prototype: กำหนดคืนของยืมค้างคืนตรวจด้วย TestHandover_OvernightDueAt แล้ว AC ของค้างคืนจะเสนอผ่าน CR เมื่อทำต่อ (2569-11-13) |
| F-04 | FR ไม่มี AC | spec.md NFR-SEC-01, NFR-AVAIL-01 | spec AC | NFR 2 ข้อไม่มี AC ไม่มี task และไม่มี test (TLS 1.2 ขึ้นไป / พร้อมใช้งาน 99% ช่วง 08.30-16.30 น.) ส่วน "ต้องผ่านการยืนยันตัวตนก่อนทุกครั้ง" มี test ผ่าน TestHandover_NonStaffForbidden แล้ว |ไม่ใช่ปัญหาใน prototype: TLS ทำโดย 004 ตรวจด้วย AC-APV-19 ส่วน NFR-AVAIL-01 เป็น NFR ของระบบรวม ตรวจที่ deploy แนวเดียวกับ 001 F-08 และ 006 F-04 (2569-11-13) |
| F-06 | โค้ดไม่มี FR | service.go: `handover` (`MissingFieldsError`), handler.go: `verifyBorrower` / `handover` (`ErrBorrowerNotVerified` -> 422 `VALIDATION_FAILED` ช่อง `borrowerCredential`) | plan.md ข้อ 4 ตาราง error | โค้ดตัดสินเองว่า ชื่อ นามสกุล หรือ credential ว่าง และ credential ใช้ไม่ได้ ตอบ 422 `VALIDATION_FAILED` พร้อมช่อง plan ข้อ 4 มีแค่แถว `cardChecked` และ `photos` (ดู test-cases TC-HND-02-4, 02-5, 07-3 ที่ Then เขียนว่า spec ไม่ได้บอก) ควรเพิ่มแถวใน plan ข้อ 4 ถ้าทีมรับพฤติกรรมนี้ |ไม่ใช่ปัญหา: FR-HND-03 บังคับบันทึกชื่อ นามสกุล ช่องว่างจึงต้องปฏิเสธ ใช้ VALIDATION_FAILED จากรายการกลาง CON-TECH-03 (2569-11-13) |
| F-07 | ตัวเลขไม่ตรง spec | photostore.go: `photoExt` | AGENTS.md กฎข้อ 9, plan ข้อ 5 (`HandoverPhotoTypes`) | ชนิดรูป JPG/PNG ของ FR-HND-06 อยู่ใน config.go แล้ว แต่ `photoExt` ประกาศรายการชนิดซ้ำอีกที่ ถ้าเพิ่มชนิดใน config ตามคำตอบในอนาคต `Save` จะปฏิเสธด้วย `ErrPhotoType` หลังจากผ่านการตรวจแล้ว |ไม่ใช่ปัญหาใน prototype: photoExt กับ config.HandoverPhotoTypes ตรงกัน (JPG, PNG) ถ้าเพิ่มชนิดภายหลังต้องแก้ 2 ที่ แนวเดียวกับ 001 F-10 (2569-11-13) |
| F-08 | test อ่อน | config.go: `BorrowerAuthTimeout`, handler.go: `NewHandler` | AC-HND-03 ("ไม่ตอบภายใน 10 วินาที") | TestAC_HND_03 ใช้ timeout 50 ms ผ่าน `WithTimeout` ไม่มี test ใดตรวจว่าระบบจริงใช้ 10 วินาที ถ้าแก้ `BorrowerAuthTimeout` เป็นค่าอื่น หรือ `NewHandler` ไม่ห่อ `WithTimeout` test ทุกตัวยังผ่าน (test-cases TC-HND-03-4 เสนอไว้แล้ว) | ไม่ใช่ปัญหา: NewHandler ใช้ config.BorrowerAuthTimeout (10 วินาที) ที่เดียว test ใช้ 50 ms เพื่อไม่ต้องรอจริง ตรวจด้วยการอ่าน handler.go บรรทัด 44 (2569-11-13)|
| F-09 | ช่องว่างของ spec | spec.md AC-HND-13 / handler_test.go: TestAC_HND_13_OtherRequestUntouched | Q-SH-06 (ปิดแล้วใน constraints.md v5.5: คำขออื่นของชิ้นเดียวกันถูก 004 ปฏิเสธทันทีเมื่อมีคำขอได้รับอนุมัติ) | Given ของ AC-HND-13 ("X อนุมัติแล้ว Y ยังไม่ได้ตัดสิน" ของชิ้นเดียวกัน) เกิดในระบบจริงไม่ได้แล้ว เพราะ 004 เปลี่ยน Y เป็น rejected ตอนอนุมัติ X test ตั้งสถานะ Y เองจึงยังผ่าน แต่ตรวจสถานการณ์ที่ไม่เกิดจริง spec 005 บรรทัด 195 เขียนไว้เองว่า "ทีมต้องตัดสินว่าแก้ Given" (สถานะ Q-SH-06 ใน spec แก้ให้ตรงไฟล์กลางแล้วใน v2.10 แต่ Given ของ AC-HND-13 ยังไม่ได้แก้) แนวทาง: แก้ Given ของ AC-HND-13 ให้ Y เป็นคำขอที่ 004 ปฏิเสธแล้ว (ส่งมอบ Y ไม่ได้) หรือบันทึกว่าเป็น test ป้องกันไว้ และแก้สถานะ Q-SH-06 ใน spec | ไม่ใช่ปัญหา: AC-HND-13 เป็น test ป้องกันไว้ ระบบจริง 004 ปฏิเสธ Y ตอนอนุมัติ X ตาม Q-SH-06 แต่ถ้ามีคำขออื่นค้างอยู่ 005 ก็ไม่แตะ (2569-11-13)|

ข้อค้นพบเดิมที่ยังเจออยู่ (F-01 ถึง F-04, F-06 ถึง F-08) คง F-ID ช่องทีมตัดสินยังว่าง (เจ้าของยังไม่ได้เขียน) โค้ดยังเหมือนรอบ 1 / F-05 แก้แล้ว (รอบ 3)
ข้อค้นพบใหม่รอบนี้: F-09 (ช่องทีมตัดสินว่าง ทีมเขียนเอง)

## 4. แก้แล้ว
| F-ID | แก้อย่างไร | รู้ได้อย่างไร |
|---|---|---|
| F-05 | HTTPS: 004 ทำ service backend :8443 พร้อม TLS 1.2 (T-50 ของ 004) / volume: บรรลุเพิ่ม `PHOTO_DIR: /data/photos` และ volume `photos` ใน service backend (7d52293, T-07) merge เข้า main โดยพรรณธรที่ a0c1b5a แทน PR #38 ที่ปิดเพราะ conflict (ช่องทีมตัดสินเดิมว่าง) | อ่าน deploy/docker-compose.yml บรรทัด 36-48 และ `docker compose config -q` ผ่าน / config.PhotoDir() อ่าน PHOTO_DIR ตัวเดียวกัน |
