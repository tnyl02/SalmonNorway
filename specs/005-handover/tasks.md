# Tasks: บันทึกการส่งมอบอุปกรณ์ (Handover)
Spec: SPEC-HND-005 (spec.md Draft v2.6) | Plan: plan.md Draft v1.4 | tasks v1.4: 2569-11-11 (v1.1 ไม่กระทบฟีเจอร์อื่น, v1.2 รับผล CR-003, v1.3 รับผล CR-004, migration 0003 และ 004 spec v1, v1.4 ปิด Q-HND-11, Q-HND-13, Q-P4) | Owner: 660710764 นายบรรลุ ฉิมสิงห์ (@mmoonn101219)

สรุป
- 14 เฟส (ฐานราก + 12 เฟสตาม AC + รอคำตอบ) รวม 31 task
- รอ Q-xx 7 task (T-25 ถึง T-31 ในเฟสรอคำตอบ) เส้นทางหลักไม่ติดคำถามแล้ว (Q-HND-11, Q-HND-13, Q-P4 ปิด 2569-11-11)
  เสร็จแล้ว: T-01, T-02, T-06, T-07, T-08, T-18 ถัดไป: T-03 -> T-04 -> T-05 และ T-09 -> T-10 -> ที่เหลือ
- เฟสแรกที่ลองใช้งานได้จริง: เฟส 1 (เห็นรายการคำขอที่รอส่งมอบ) migration 0003 ของ 002 merge แล้ว

เทคโนโลยีตาม plan.md ข้อ 2: Go (net/http + pgx), PostgreSQL 16, React + Vite, เก็บรูปบน Docker volume
(ทีมเลือกเอง ไม่ได้มาจาก spec ยกเว้น PostgreSQL / HTTPS / Docker ตาม CON-TECH-01, 02)
คำสั่งเปิดระบบ (ใช้ในทุกจุดตรวจ):
- ฐานข้อมูล: `docker compose -f deploy/docker-compose.yml up -d db`
- หลังบ้าน: `cd backend && go run ./cmd/server` (port 8000)
- หน้าจอ: `cd frontend && npm run dev` (port 5173)
- ผู้ใช้จำลอง: `Authorization: Bearer mock-staff-001` (เจ้าหน้าที่), `mock-student-001` (นักศึกษา A), `mock-student-101` (นักศึกษา B)
- test: `cd backend && go test ./... -v` และ `cd frontend && npm test`

กติกาไม่กระทบฟีเจอร์อื่น (plan.md v1.1 ข้อ 3, 4)
- ไม่แก้โครงตาราง `loan` ของ 001 ข้อมูลการส่งมอบอยู่ในตาราง `handover` ของ 005 (seed และ test ของ 001 ไม่พัง ไม่ต้องแก้ DM-05)
- migration ของ 005 มี 2 ไฟล์ เลขถัดจาก 0005 ของ 003 ต้องจองในกลุ่มก่อน ไฟล์ action ต้องมีทุกค่าเดิมใน CHECK (list, search, detail, cancel, request, notify และของ 004)
- 16:00 ใช้ `LoanNormalReturnHour` ของ 002 ใน config.go (merge เข้า main แล้ว) ไม่ประกาศซ้ำ
- error ตอบเป็น `{code, message}` ตาม CON-TECH-03 ใช้รหัสกลางเท่านั้น รหัสใหม่ 3 ตัวของ 005 ลงทะเบียนแล้ว (Q-P4)
- สถานะคำขอที่ส่งมอบได้อ่านจาก `HandoverReadyStatuses()` ใน config ไม่ต้องแก้ของ 002 ไม่อ่านตารางของ 004
- audit log ใช้ writer.go ของ 001 ถ้ามีแล้ว ถ้ายังไม่มีเขียนใน repository ของ 005 (แผนสำรองแบบ 003) ไม่สร้าง package `audit` เอง
- ไฟล์ของกลาง (`config.go`, `config.js`, `main.go`, `App.jsx`, `client.js`, `docker-compose.yml`) เพิ่มต่อท้ายเป็นบล็อกของ 005 เท่านั้น ทำหลัง merge 003

---

## เฟส 0: ฐานราก

### T-01 เพิ่มค่าที่ตั้งได้ของการส่งมอบ
- รองรับ: FR-HND-02, FR-HND-05, FR-HND-06, DOM-LOAN-01
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-07, T-10, T-13, T-14, T-15, T-22
- ไฟล์ที่แตะ: `backend/internal/config/config.go`, `frontend/src/config.js`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: บล็อก `// ---- 005-handover ----` ใน `config.go` มี `BorrowerAuthTimeout`, `HandoverMinPhotos`, `HandoverPhotoMaxBytes`, `HandoverPhotoTypes`, `PhotoDir` (อ่านจาก env) (ค่า 16:00 ใช้ `LoanNormalReturnHour` ของ 002) บล็อกเดียวกันใน `config.js` มี `HANDOVER_MIN_PHOTOS`, `HANDOVER_PHOTO_MAX_BYTES`, `HANDOVER_PHOTO_TYPES` ทุกตัวมีคอมเมนต์อ้าง ID ไม่แก้บรรทัดของฟีเจอร์อื่น และ `go test ./...`, `npm test` เดิมยังผ่าน
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

### T-02 สร้าง migration ของการส่งมอบและข้อมูลตัวอย่าง
- รองรับ: CON-TECH-01, FR-HND-03, FR-HND-05, FR-HND-06, FR-HND-07
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-04 ขึ้นไป
- ไฟล์ที่แตะ: `backend/db/migrations/0006_handover.up.sql` (ทีมกำหนดเลข 0006 2569-11-11 ถัดจาก 0005 ของ 003) สร้างตาราง `handover` และ `handover_photo` **ไม่ ALTER `loan`**, `backend/internal/handover/migration_test.go` (ไม่แตะ `platform/migrate_test.go` ของ 001), `backend/internal/handover/fixture_test.go`, `backend/db/seed/handover_dev.sql` (ข้อมูลลองระบบ ทีมเลือกเอง ไม่ได้มาจาก spec)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01 (migration 0003 ของ 002 merge แล้ว)
- เสร็จเมื่อ: `TestMigration_HandoverPhotoContentType` ผ่าน (`content_type` นอก jpeg/png ถูกปฏิเสธ) `TestMigration_HandoverIdentityMethod` ผ่าน (`identity_method` นอก `auth_api`/`card_only` ถูกปฏิเสธ) ไม่มีคอลัมน์เลขบัตรใด ๆ และ test ของ 001 ทั้งหมดยังผ่าน
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

### T-03 สร้าง package handover ผูก route และกันผู้ที่ไม่ใช่เจ้าหน้าที่
- รองรับ: CON-HND-01, IF-IDP-01, CON-TECH-02
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐาน (AC-HND-09 ตรวจครบใน T-23)
- ไฟล์ที่แตะ: `backend/internal/handover/handler.go`, `backend/internal/handover/errors.go`, `backend/internal/handover/handler_test.go`, `backend/cmd/server/main.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: `TestHandover_NonStaffForbidden` ผ่าน (ทั้ง 3 route ตอบ 403 `{code, message}` `FORBIDDEN_ROLE` เมื่อใช้ `mock-student-001` และตอบ 401 เมื่อไม่มี token ผ่าน middleware เดิม)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

#### จุดตรวจเฟส 0
- test ที่ต้องผ่านทั้งหมด: TestMigration_HandoverPhotoContentType, TestMigration_HandoverIdentityMethod, TestHandover_NonStaffForbidden และ test เดิมของ 001 ทั้งหมด
- เปิดระบบ: ฐานข้อมูลและหลังบ้านตามคำสั่งด้านบน
- คนลองเอง (ไม่มี AC เป็นงานพื้นฐาน):
  Given เปิดฐานข้อมูลและหลังบ้านแล้ว
  When เรียก `GET http://localhost:8000/api/v1/handover/requests` ด้วย `Bearer mock-student-001`
  Then ได้ 403 และ body `{code, message}`
- ยังขาด: ไม่มี

---

## เฟส 1: เห็นรายการคำขอที่รอส่งมอบ (AC-HND-14)

### T-04 สร้าง API รายการคำขอที่อนุมัติแล้วและยังไม่ส่งมอบ
- รองรับ: FR-HND-01, CON-DATA-03
- ตรวจด้วย: AC-HND-14
- ไฟล์ที่แตะ: `backend/internal/config/config.go` (เพิ่ม `HandoverReadyStatuses()` ในบล็อกของ 005), `backend/internal/handover/repository.go`, `backend/internal/handover/handler.go`, `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-02, T-03
- เสร็จเมื่อ: `TestAC_HND_14_ListApprovedOnly` ผ่าน (มีข้อมูล 2 อนุมัติ, 1 ยังไม่อนุมัติ, 1 ส่งมอบแล้ว ได้กลับมา 2 รายการพร้อม `assetCode`)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

### T-05 สร้างหน้า HandoverList
- รองรับ: FR-HND-01
- ตรวจด้วย: AC-HND-14
- ไฟล์ที่แตะ: `frontend/src/pages/Handover/HandoverList.jsx`, `frontend/src/api/client.js` (เพิ่มฟังก์ชันของ 005 ต่อท้าย), `frontend/src/App.jsx` (เพิ่ม route ต่อท้าย ทำหลัง merge 003 แจ้งเจ้าของ 001), `frontend/src/__tests__/AC-HND-14.test.jsx`
- แบบหน้าจอ: ไม่มี (spec ไม่มี UI ID และไม่มี mockup ทีมออกแบบเอง)
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `AC-HND-14.test.jsx` ผ่าน (client จำลองคืน 2 คำขอ หน้าจอแสดง 2 แถว)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ผ่าน ทีมดูหน้าจอแล้วผ่าน)

#### จุดตรวจเฟส 1
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_14_ListApprovedOnly, AC-HND-14.test.jsx และ test ของเฟส 0
- เปิดระบบ: ฐานข้อมูล (โหลด `handover_dev.sql`) หลังบ้าน และหน้าจอ
- คนลองเองตาม AC-HND-14:
  Given ข้อมูลตัวอย่างมีคำขออนุมัติแล้ว 2 รายการ (อุปกรณ์คนละชิ้น) ยังไม่อนุมัติ 1 และส่งมอบแล้ว 1
  When เข้าหน้ารายการคำขอเพื่อส่งมอบในฐานะ `mock-staff-001`
  Then เห็นเฉพาะคำขอที่อนุมัติแล้วและยังไม่ส่งมอบ 2 รายการ
- ยังขาด: ไม่มี

---

## เฟส 2: ส่งมอบสำเร็จหลังยืนยันตัวตนและตรวจบัตร (AC-HND-02, AC-HND-07)
หมายเหตุ: รวม 2 AC ไว้เฟสเดียว เพราะการส่งมอบสำเร็จครั้งแรก (AC-HND-07) ต้องยืนยันผู้ยืมและตรวจบัตรก่อนเสมอ (FR-HND-02, 03)

### T-06 สร้าง transaction ส่งมอบในชั้น repository
- รองรับ: FR-HND-04, FR-HND-07, CON-DATA-01, DOM-REQ-01
- ตรวจด้วย: AC-HND-07
- ไฟล์ที่แตะ: `backend/internal/handover/repository.go`, `backend/internal/handover/repository_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-02
- เสร็จเมื่อ: `TestRepository_CreateHandover` ผ่าน (`SELECT ... FOR UPDATE` อุปกรณ์ สร้าง Loan ด้วยคอลัมน์เดิมของ 001 ที่ `returned_at` ว่าง และแถว `handover` คู่กัน อุปกรณ์เป็น `borrowed` คำขอเป็น `handed_over` ทั้งหมดใน transaction เดียว ไม่แตะคำขออื่น)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

### T-07 สร้างที่เก็บไฟล์รูปบน Docker volume
- รองรับ: FR-HND-06, CON-TECH-02
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-10, T-15 (AC-HND-06 ตรวจใน T-15)
- ไฟล์ที่แตะ: `backend/internal/handover/photostore.go`, `backend/internal/handover/photostore_test.go`, `deploy/docker-compose.yml` (เพิ่ม named volume ใหม่ `photos` ที่ `/data/photos` ไม่แก้ service เดิม 005 ใช้แค่ `/data/photos/handover/` แจ้งทีมและ 006)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: `TestPhotoStore_SaveAndRemove` ผ่าน (เขียนไฟล์ลง `PhotoDir/<loanId>/` ใน temp dir และลบทิ้งได้เมื่อ rollback)
- สถานะ: เสร็จ รอทีมตรวจ (ยังไม่แก้ docker-compose.yml ดู prompt-log)

### T-08 สร้าง BorrowerVerifier และตัวจำลอง
- รองรับ: FR-HND-02, IF-IDP-01
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-09, T-13, T-20
- ไฟล์ที่แตะ: `backend/internal/handover/verifier.go`, `backend/internal/handover/verifier_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: `TestVerifier_Mock` ผ่าน (interface คืนรหัสนักศึกษาของผู้ที่ยืนยันตัวตน ตัวจำลองใช้ผู้ใช้จาก `auth.MockProvider`)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

### T-09 สร้าง API ยืนยันตัวตนผู้ยืม
- รองรับ: FR-HND-02, IF-IDP-01
- ตรวจด้วย: AC-HND-02
- ไฟล์ที่แตะ: `backend/internal/handover/handler.go`, `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-03, T-08
- เสร็จเมื่อ: `TestAC_HND_02_AuthApiVerifiedAndCardChecked/verify_borrower` ผ่าน (`POST .../verify-borrower` คืน `verified`)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

### T-10 สร้าง API บันทึกการส่งมอบ
- รองรับ: FR-HND-02, FR-HND-03, FR-HND-05, FR-HND-07
- ตรวจด้วย: AC-HND-02, AC-HND-07
- ไฟล์ที่แตะ: `backend/internal/handover/service.go`, `backend/internal/handover/handler.go`, `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06, T-07, T-09
- เสร็จเมื่อ: `TestAC_HND_07_HandoverCreatesLoan` และ `TestAC_HND_02_AuthApiVerifiedAndCardChecked` ผ่าน (201 พร้อมชื่อ นามสกุล รหัสนักศึกษา เวลาที่ส่งมอบใน Loan / `cardChecked=false` ได้ 422 `VALIDATION_FAILED` ระบุช่อง `cardChecked` ไม่มี Loan)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

### T-11 สร้างหน้า HandoverForm ส่วนตรวจบัตร แนบรูป และบันทึก
- รองรับ: FR-HND-03, FR-HND-06, FR-HND-07
- ตรวจด้วย: AC-HND-02, AC-HND-07
- ไฟล์ที่แตะ: `frontend/src/pages/Handover/HandoverForm.jsx`, `frontend/src/api/client.js`, `frontend/src/__tests__/AC-HND-02.test.jsx`
- แบบหน้าจอ: ไม่มี (ไม่มี mockup ทีมออกแบบเอง)
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `AC-HND-02.test.jsx` ผ่าน (ไม่ติ๊กยืนยันตรวจบัตร ปุ่มบันทึกส่งไม่ได้ / ติ๊กแล้วส่ง multipart ไป client จำลอง และแสดงผลที่ได้กลับ)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ผ่าน ทีมดูหน้าจอแล้วผ่าน)

### T-12 เพิ่มส่วนยืนยันตัวตนผู้ยืมในหน้า HandoverForm
- รองรับ: FR-HND-02
- ตรวจด้วย: AC-HND-02
- ไฟล์ที่แตะ: `frontend/src/pages/Handover/HandoverForm.jsx`, `frontend/src/__tests__/AC-HND-02.test.jsx`
- แบบหน้าจอ: ไม่มี (ไม่มี mockup ทีมออกแบบเอง)
- ต้องทำหลัง: T-11
- เสร็จเมื่อ: `AC-HND-02.test.jsx` ผ่านกรณียืนยันตัวตนผู้ยืมก่อนบันทึก
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ผ่าน ทีมดูหน้าจอแล้วผ่าน)

#### จุดตรวจเฟส 2
- test ที่ต้องผ่านทั้งหมด: TestRepository_CreateHandover, TestPhotoStore_SaveAndRemove, TestVerifier_Mock, TestAC_HND_02_AuthApiVerifiedAndCardChecked, TestAC_HND_07_HandoverCreatesLoan, AC-HND-02.test.jsx และ test ของเฟส 0-1
- เปิดระบบ: ฐานข้อมูล หลังบ้าน และหน้าจอ
- คนลองเองตาม AC-HND-02, AC-HND-07:
  Given คำขออนุมัติแล้วของ `mock-student-001` อุปกรณ์ "ว่าง" และมีรูป JPG 2 รูปในเครื่อง
  When เจ้าหน้าที่เลือกคำขอ ยืนยันตัวตนผู้ยืม ติ๊กว่าตรวจบัตรแล้ว แนบรูป แล้วกดบันทึก
  Then บันทึกสำเร็จ หน้ารายละเอียดอุปกรณ์ของ 001 แสดง "ถูกยืม" และคำขอหายจากรายการรอส่งมอบ
       ลองอีกครั้งโดยไม่ติ๊กตรวจบัตร ต้องบันทึกไม่ได้
- ยังขาด: ไม่มี

---

## เฟส 3: Auth API ไม่ตอบ ใช้ผลตรวจบัตร (AC-HND-03)

### T-13 ใส่ timeout 10 วินาทีและโหมดตรวจบัตรอย่างเดียว
- รองรับ: FR-HND-02, FR-HND-03, IF-IDP-01
- ตรวจด้วย: AC-HND-03
- ไฟล์ที่แตะ: `backend/internal/handover/verifier.go`, `backend/internal/handover/service.go`, `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี (หน้าจอใช้ผล `unavailable` จาก T-12)
- ต้องทำหลัง: T-08, T-10
- เสร็จเมื่อ: `TestAC_HND_03_AuthApiTimeoutCardOnly` ผ่าน (verifier ค้างเกิน `BorrowerAuthTimeout` ที่ตั้งสั้นใน test ได้ `unavailable` บันทึกด้วย `card_only` ได้ 201 / ไม่ยืนยันบัตรได้ 422)
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 3
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_03_AuthApiTimeoutCardOnly และ test ของเฟส 0-2
- เปิดระบบ: ตั้ง env ให้ใช้ verifier จำลองแบบค้าง (ทีมเลือกเอง) แล้วเปิดหลังบ้านและหน้าจอ
- คนลองเองตาม AC-HND-03:
  Given คำขออนุมัติแล้ว และ verifier จำลองไม่ตอบ
  When เจ้าหน้าที่ยืนยันตัวตนผู้ยืม รอเกิน 10 วินาที ติ๊กตรวจบัตร แล้วบันทึก
  Then บันทึกสำเร็จ Loan มี `identity_method = card_only`
- ยังขาด: ไม่มี

---

## เฟส 4: เวลาส่งมอบและกำหนดคืน 16:00 (AC-HND-05)

### T-14 ทำให้นาฬิกา inject ได้และคำนวณกำหนดคืน 16:00 วันเดียวกัน
- รองรับ: FR-HND-05, DOM-LOAN-01
- ตรวจด้วย: AC-HND-05
- ไฟล์ที่แตะ: `backend/internal/handover/service.go`, `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `TestAC_HND_05_DueAt1600SameDay` ผ่าน (clock 2569-11-12 10:00 Asia/Bangkok ได้ `handed_over_at` 10:00 และ `due_at` 16:00 วันเดียวกัน ค่า 16:00 อ่านจาก `LoanNormalReturnHour`)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

#### จุดตรวจเฟส 4
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_05_DueAt1600SameDay และ test ของเฟส 0-3
- เปิดระบบ: หลังบ้านและหน้าจอ
- คนลองเองตาม AC-HND-05:
  Given คำขออนุมัติแล้ว ลองก่อน 16:00
  When บันทึกการส่งมอบ
  Then ผลที่หน้าจอแสดงเวลาส่งมอบตามเวลาจริง และกำหนดคืน 16:00 วันเดียวกัน
- ยังขาด: กรณีหลัง 16:00 (T-25 รอ Q-HND-03)

---

## เฟส 5: รูปสภาพอุปกรณ์ผูกกับการยืม (AC-HND-06)

### T-15 ตรวจจำนวนและชนิดรูป แล้วบันทึกแถว handover_photo
- รองรับ: FR-HND-06
- ตรวจด้วย: AC-HND-06
- ไฟล์ที่แตะ: `backend/internal/handover/photos.go`, `backend/internal/handover/service.go`, `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07, T-10
- เสร็จเมื่อ: `TestAC_HND_06_TwoJpgPhotosSaved` ผ่าน (JPG 2 รูป รูปละ 2 MB ได้ 201 มี `handover_photo` 2 แถวผูก Loan และไฟล์อยู่บนดิสก์ ชนิดไฟล์ตรวจจากเนื้อไฟล์ ไม่ใช่นามสกุล)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

#### จุดตรวจเฟส 5
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_06_TwoJpgPhotosSaved และ test ของเฟส 0-4
- เปิดระบบ: ฐานข้อมูล หลังบ้าน และหน้าจอ
- คนลองเองตาม AC-HND-06:
  Given คำขออนุมัติแล้ว และรูป JPG 2 รูป รูปละประมาณ 2 MB
  When แนบรูปทั้ง 2 แล้วบันทึกการส่งมอบ
  Then บันทึกสำเร็จ มีไฟล์ 2 ไฟล์ใน volume ใต้โฟลเดอร์ของ Loan นั้น
- ยังขาด: การตรวจขนาด 10 MB (T-28 รอ Q-HND-10)

---

## เฟส 6: คำขอยังไม่อนุมัติ (AC-HND-01)

### T-16 ปฏิเสธคำขอที่ยังไม่อนุมัติหรือส่งมอบไปแล้ว
- รองรับ: FR-HND-01
- ตรวจด้วย: AC-HND-01
- ไฟล์ที่แตะ: `backend/internal/handover/service.go`, `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `TestAC_HND_01_RequestNotApproved` ผ่าน (409 `REQUEST_NOT_APPROVED` ไม่มี Loan อุปกรณ์ยัง `available`)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

#### จุดตรวจเฟส 6
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_01_RequestNotApproved และ test ของเฟส 0-5
- เปิดระบบ: ฐานข้อมูลและหลังบ้าน
- คนลองเองตาม AC-HND-01:
  Given คำขอที่ยังไม่อนุมัติในข้อมูลตัวอย่าง อุปกรณ์ "ว่าง"
  When เรียก `POST /api/v1/handover/requests/{requestId}/handover` ด้วย `mock-staff-001`
  Then ได้ 409 และอุปกรณ์ยัง "ว่าง"
- ยังขาด: ไม่มี

---

## เฟส 7: อุปกรณ์ไม่ว่าง (AC-HND-04)

### T-17 ทดสอบการปฏิเสธเมื่ออุปกรณ์ไม่ "ว่าง"
- รองรับ: FR-HND-04, CON-DATA-01
- ตรวจด้วย: AC-HND-04
- ไฟล์ที่แตะ: `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `TestAC_HND_04_ItemNotAvailable` ผ่าน (อุปกรณ์ `unavailable` และ `borrowed` ได้ 409 `ITEM_NOT_AVAILABLE` สถานะไม่เปลี่ยน)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

### T-18 ทดสอบส่งมอบอุปกรณ์ชิ้นเดียวกันพร้อมกัน
- รองรับ: FR-HND-04, DOM-REQ-01
- ตรวจด้วย: AC-HND-04 (ไม่มี AC แยก plan ข้อ 6 สั่งเพิ่ม)
- ไฟล์ที่แตะ: `backend/internal/handover/repository_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestHandover_ConcurrentSameItem` ผ่าน (2 goroutine ส่งมอบชิ้นเดียวกัน สำเร็จ 1 อีกอันได้ `ITEM_NOT_AVAILABLE` มี Loan 1 แถว)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

#### จุดตรวจเฟส 7
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_04_ItemNotAvailable, TestHandover_ConcurrentSameItem และ test ของเฟส 0-6
- เปิดระบบ: ฐานข้อมูล หลังบ้าน และหน้าจอ
- คนลองเองตาม AC-HND-04:
  Given ตั้งอุปกรณ์ของคำขออนุมัติแล้วเป็น "ไม่พร้อมใช้งาน" ในฐานข้อมูล
  When เจ้าหน้าที่บันทึกการส่งมอบ
  Then บันทึกไม่ได้ หน้าจอแสดงข้อผิดพลาด และสถานะอุปกรณ์ยังเป็น "ไม่พร้อมใช้งาน"
- ยังขาด: ไม่มี

---

## เฟส 8: ไม่แตะคำขออื่นของอุปกรณ์ชิ้นเดียวกัน (AC-HND-13)

### T-19 ทดสอบว่าคำขอ Y ไม่เปลี่ยนและส่งมอบ Y ไม่ได้
- รองรับ: FR-HND-07, FR-HND-04, DOM-REQ-01
- ตรวจด้วย: AC-HND-13
- ไฟล์ที่แตะ: `backend/internal/handover/handler_test.go`, `backend/internal/handover/fixture_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10, T-17
- เสร็จเมื่อ: `TestAC_HND_13_OtherRequestUntouched` ผ่าน (ส่งมอบ X สำเร็จ สถานะของ Y เท่าเดิม แล้วส่งมอบ Y ได้ 409 `ITEM_NOT_AVAILABLE`)
- สถานะ: เสร็จ รอทีมตรวจ

#### จุดตรวจเฟส 8
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_13_OtherRequestUntouched และ test ของเฟส 0-7
- เปิดระบบ: ฐานข้อมูลและหลังบ้าน
- คนลองเองตาม AC-HND-13:
  Given อุปกรณ์ชิ้นเดียวมีคำขอ X อนุมัติแล้ว และคำขอ Y ยังไม่ได้ตัดสิน
  When ส่งมอบ X สำเร็จ แล้วดูสถานะ Y ในฐานข้อมูล
  Then สถานะ Y เท่าเดิม
- ยังขาด: ไม่มี (การส่งมอบ Y ต้องรอ 004 อนุมัติ Y ก่อน ทดสอบได้แค่ใน test ที่ตั้งสถานะเอง)

---

## เฟส 9: ผู้รับไม่ใช่ผู้ส่งคำขอ (AC-HND-10)

### T-20 ปฏิเสธเมื่อผู้ยืนยันตัวตนไม่ใช่ผู้ส่งคำขอ
- รองรับ: FR-HND-02
- ตรวจด้วย: AC-HND-10
- ไฟล์ที่แตะ: `backend/internal/handover/service.go`, `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-08, T-10
- เสร็จเมื่อ: `TestAC_HND_10_BorrowerMismatch` ผ่าน (verifier คืนนักศึกษา B สำหรับคำขอของ A ได้ 409 `BORROWER_MISMATCH` อุปกรณ์ยัง `available`)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

#### จุดตรวจเฟส 9
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_10_BorrowerMismatch และ test ของเฟส 0-8
- เปิดระบบ: ฐานข้อมูล หลังบ้าน และหน้าจอ
- คนลองเองตาม AC-HND-10:
  Given คำขออนุมัติแล้วของ `mock-student-001`
  When ผู้ที่ยืนยันตัวตนที่เคาน์เตอร์คือ `mock-student-101` แล้วเจ้าหน้าที่กดบันทึก
  Then บันทึกไม่ได้ และอุปกรณ์ยัง "ว่าง"
- ยังขาด: ไม่มี

---

## เฟส 10: รูปไม่ถูกต้อง (AC-HND-12)

### T-21 ทดสอบรูปไม่ครบและชนิดผิด และไม่มีไฟล์ค้าง
- รองรับ: FR-HND-06
- ตรวจด้วย: AC-HND-12
- ไฟล์ที่แตะ: `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-15
- เสร็จเมื่อ: `TestAC_HND_12_InvalidPhotos` กรณี (ก) 1 รูป และ (ข) GIF ผ่าน (422 `VALIDATION_FAILED` ระบุช่อง `photos` ไม่มี Loan ไม่มีไฟล์ค้างบนดิสก์)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

### T-22 ตรวจจำนวนและชนิดรูปที่หน้าจอก่อนส่ง
- รองรับ: FR-HND-06
- ตรวจด้วย: AC-HND-12
- ไฟล์ที่แตะ: `frontend/src/pages/Handover/HandoverForm.jsx`, `frontend/src/__tests__/AC-HND-12.test.jsx`
- แบบหน้าจอ: ไม่มี (ไม่มี mockup ทีมออกแบบเอง)
- ต้องทำหลัง: T-11, T-21
- เสร็จเมื่อ: `AC-HND-12.test.jsx` ผ่าน (1 รูป หรือ GIF ไม่เรียก client จำลอง ใช้ค่าจาก `config.js`)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ผ่าน ทีมดูหน้าจอแล้วผ่าน)

#### จุดตรวจเฟส 10
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_12_InvalidPhotos (ก, ข), AC-HND-12.test.jsx และ test ของเฟส 0-9
- เปิดระบบ: ฐานข้อมูล หลังบ้าน และหน้าจอ
- คนลองเองตาม AC-HND-12:
  Given คำขออนุมัติแล้ว
  When แนบรูป 1 รูป แล้วลองใหม่ด้วยไฟล์ GIF
  Then บันทึกไม่ได้ทั้ง 2 ครั้ง และคำขอยังอยู่ในรายการรอส่งมอบ
- ยังขาด: กรณี (ค) รูป 11 MB (T-28 รอ Q-HND-10)

---

## เฟส 11: นักศึกษาเรียกบันทึกการส่งมอบเอง (AC-HND-09)

### T-23 ทดสอบนักศึกษาถูกปฏิเสธและไม่มีอะไรเปลี่ยน
- รองรับ: CON-HND-01
- ตรวจด้วย: AC-HND-09
- ไฟล์ที่แตะ: `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-03, T-10
- เสร็จเมื่อ: `TestAC_HND_09_StudentForbidden` ผ่าน (`mock-student-001` เรียก POST handover ของคำขออนุมัติแล้ว ได้ 403 ไม่มี Loan สถานะอุปกรณ์ไม่เปลี่ยน)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

#### จุดตรวจเฟส 11
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_09_StudentForbidden และ test ของเฟส 0-10
- เปิดระบบ: ฐานข้อมูลและหลังบ้าน
- คนลองเองตาม AC-HND-09:
  Given คำขออนุมัติแล้ว
  When เรียก `POST /api/v1/handover/requests/{requestId}/handover` ด้วย `mock-student-001`
  Then ได้ 403 และอุปกรณ์ยัง "ว่าง"
- ยังขาด: ไม่มี

---

## เฟส 12: ส่งมอบหลัง 16:00 (AC-HND-11 บางส่วน)

### T-24 ทดสอบว่าส่งมอบเวลา 16:30 บันทึกได้
- รองรับ: FR-HND-05
- ตรวจด้วย: AC-HND-11
- ไฟล์ที่แตะ: `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-14
- เสร็จเมื่อ: `TestAC_HND_11_HandoverAfter1600` ผ่าน (clock 16:30 ได้ 201 `handed_over_at` 16:30 อุปกรณ์ `borrowed` **ไม่ตรวจ `due_at`**)
- สถานะ: เสร็จ (ตรวจ 2569-11-12 test ในช่องเสร็จเมื่อผ่าน)

#### จุดตรวจเฟส 12
- test ที่ต้องผ่านทั้งหมด: TestAC_HND_11_HandoverAfter1600 และ test ของเฟส 0-11
- เปิดระบบ: ฐานข้อมูล หลังบ้าน และหน้าจอ
- คนลองเองตาม AC-HND-11:
  Given คำขออนุมัติแล้ว ลองหลัง 16:00
  When บันทึกการส่งมอบ
  Then บันทึกสำเร็จ เวลาส่งมอบตามเวลาจริง อุปกรณ์ "ถูกยืม" (ไม่ต้องดูกำหนดคืน)
- ยังขาด: กำหนดคืนของกรณีนี้ (T-25 รอ Q-HND-03)

---

## เฟส AC-HND-08 (audit log): ย้ายไปเฟสรอคำตอบทั้งเฟส
ค่า `action` ยังไม่มีใน DM-07 (Q-HND-06) และกรณีไม่สำเร็จรอ Q-HND-14 จึงไม่มีส่วนที่ทำได้ก่อน ดู T-26, T-27

---

## เฟส 13: รอคำตอบ

### T-25 คำนวณกำหนดคืนเมื่อส่งมอบหลัง 16:00
- รองรับ: FR-HND-05, DOM-LOAN-01
- ตรวจด้วย: AC-HND-11
- ไฟล์ที่แตะ: `backend/internal/handover/service.go`, `backend/internal/config/config.go` (ถ้าคำตอบมีค่าใหม่), `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-24
- เสร็จเมื่อ: `TestAC_HND_11_HandoverAfter1600` ตรวจ `due_at` ตามคำตอบ Q-HND-03 และผ่าน
- สถานะ: รอ Q-HND-03

### T-26 เขียน audit log เมื่อส่งมอบสำเร็จ
- รองรับ: FR-HND-08, DOM-AUDIT-01, CON-DATA-03
- ตรวจด้วย: AC-HND-08
- ไฟล์ที่แตะ: `backend/db/migrations/0007_audit_log_handover_action.up.sql` (ทีมกำหนดเลข 0007 2569-11-11 CHECK ใหม่ต้องมีทุกค่าเดิมของทุกฟีเจอร์ ณ เวลานั้น บวกค่าของ 005), `backend/internal/handover/repository.go` (ใช้ `backend/internal/audit/writer.go` ของ 001 ถ้ามีแล้ว ถ้ายังไม่มี INSERT เองแบบแผนสำรองของ 003), `backend/internal/handover/handler_test.go`, `specs/000-shared/data-model.md` (เพิ่มแถว action ของ 005 ในตาราง DM-07 แจ้งเจ้าของ 001)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `TestAC_HND_08_AuditLogged` ผ่านกรณีสำเร็จ (มีแถวระบุผู้ทำ เวลา รหัสครุภัณฑ์ รหัสคำขอ ผลลัพธ์) และ `TestMigration_AuditLogAction` ของ 001 กับ test ของ 003 ที่ใช้ action `cancel`, `notify` ยังผ่าน
- สถานะ: รอ Q-HND-06

### T-27 เขียน audit log เมื่อส่งมอบไม่สำเร็จ
- รองรับ: FR-HND-08, DOM-AUDIT-01
- ตรวจด้วย: AC-HND-08
- ไฟล์ที่แตะ: `backend/internal/handover/service.go`, `backend/internal/handover/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-26
- เสร็จเมื่อ: `TestAC_HND_08_AuditLogged` กรณีไม่สำเร็จผ่าน ตามคำตอบ Q-HND-14
- สถานะ: รอ Q-HND-14

### T-28 ตรวจขนาดรูปไม่เกิน 10 MB
- รองรับ: FR-HND-06
- ตรวจด้วย: AC-HND-12, AC-HND-06
- ไฟล์ที่แตะ: `backend/internal/handover/photos.go`, `backend/internal/handover/handler_test.go`, `frontend/src/pages/Handover/HandoverForm.jsx`, `frontend/src/__tests__/AC-HND-12.test.jsx`
- แบบหน้าจอ: ไม่มี (ไม่มี mockup ทีมออกแบบเอง)
- ต้องทำหลัง: T-21, T-22
- เสร็จเมื่อ: `TestAC_HND_12_InvalidPhotos` กรณี (ค) 11 MB ผ่าน และ `AC-HND-12.test.jsx` กรณีขนาดผ่าน (ต่อรูปหรือรวมตามคำตอบ)
- สถานะ: รอ Q-HND-10

### T-29 เชื่อม Auth API จริงแทนตัวจำลอง
- รองรับ: FR-HND-02, IF-IDP-01
- ตรวจด้วย: AC-HND-02, AC-HND-03
- ไฟล์ที่แตะ: `backend/internal/handover/verifier.go` (ตัวจริง), `backend/internal/config/config.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-13
- เสร็จเมื่อ: ยืนยันตัวตนผู้ยืมกับ Auth API จริงได้ และ `TestAC_HND_02`, `TestAC_HND_03` ยังผ่าน
- สถานะ: รอ Q-SH-03 (วิธีล็อกอินที่เคาน์เตอร์ตอบแล้วใน ASM-HND-09 เหลือโปรโตคอลของ Auth API ตัวจริง)

### T-30 รองรับการส่งมอบยืมแบบค้างคืน
- รองรับ: FR-HND-01, FR-HND-05, FR-HND-06
- ตรวจด้วย: ไม่มี AC (spec ไม่มี AC ของค้างคืน ควรเสนอทีมเพิ่มเมื่อได้คำตอบ)
- ไฟล์ที่แตะ: `backend/internal/handover/service.go`, `backend/internal/config/config.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-14
- เสร็จเมื่อ: ส่งมอบคำขอค้างคืนได้ `loan_type = overnight` และ `due_at` ตามคำตอบ Q-SH-04
- สถานะ: รอ Q-HND-07 (004 spec v1 ตัด UC-14 ออก) กำหนดคืน = 16:00 ของ `borrow_request.end_at` (Q-SH-04, Q-HND-05 ปิดแล้ว)

### T-31 ทำกฎ "ใช้ได้เฉพาะในพื้นที่มหาวิทยาลัย"
- รองรับ: DOM-LOAN-01
- ตรวจด้วย: ไม่มี AC (ขึ้นกับคำตอบว่าระบบต้องทำอะไร)
- ไฟล์ที่แตะ: ยังบอกไม่ได้จนกว่าจะได้คำตอบ
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: ตามคำตอบ Q-HND-04 (ถ้าเป็นแค่กฎที่เจ้าหน้าที่แจ้งผู้ยืม task นี้ยกเลิก)
- สถานะ: รอ Q-HND-04

---

## ตารางตรวจความครบ

### 1. AC -> task
| AC ID | task ที่ตรวจ AC นี้ |
|---|---|
| AC-HND-01 | T-16 |
| AC-HND-02 | T-09, T-10, T-11, T-12, T-29 |
| AC-HND-03 | T-13, T-29 |
| AC-HND-04 | T-17, T-18 |
| AC-HND-05 | T-14 |
| AC-HND-06 | T-15, T-28 |
| AC-HND-07 | T-06, T-10, T-11 |
| AC-HND-08 | T-26, T-27 |
| AC-HND-09 | T-23 (และ T-03 ส่วน 403) |
| AC-HND-10 | T-20 |
| AC-HND-11 | T-24, T-25 |
| AC-HND-12 | T-21, T-22, T-28 |
| AC-HND-13 | T-19 |
| AC-HND-14 | T-04, T-05 |

### 2. Constraint -> task
| Constraint ID | task ที่ทำให้เป็นจริง |
|---|---|
| CON-TECH-01 | T-02 |
| CON-TECH-02 | T-03 (ใช้ HTTPS/Docker ของ `deploy/` เดิม), T-07 (volume รูปใน docker-compose) |
| CON-TECH-03 | T-03, T-10, T-16, T-20, T-21 (error `{code, message}` และรหัสกลาง) |
| CON-DATA-01 | T-06, T-17 |
| CON-DATA-03 | T-04, T-26 |
| DOM-AUDIT-01 | T-26, T-27 |
| DOM-LOAN-01 | T-01, T-14, T-25, T-31 |
| DOM-REQ-01 | T-06, T-18, T-19 |
| IF-IDP-01 | T-03, T-08, T-09, T-13, T-29 |
| CON-HND-01 | T-03, T-23 |

### 3. FR -> AC -> เฟส
| FR ID | AC ที่ตรวจ FR นี้ | เฟส |
|---|---|---|
| FR-HND-01 | AC-HND-01, AC-HND-14 | 1, 6 |
| FR-HND-02 | AC-HND-02, AC-HND-03, AC-HND-10 | 2, 3, 9 |
| FR-HND-03 | AC-HND-02, AC-HND-03 | 2, 3 |
| FR-HND-04 | AC-HND-04, AC-HND-13 | 7, 8 |
| FR-HND-05 | AC-HND-05, AC-HND-11 | 4, 12, 13 |
| FR-HND-06 | AC-HND-06, AC-HND-12 | 5, 10, 13 |
| FR-HND-07 | AC-HND-07, AC-HND-13 | 2, 8 |
| FR-HND-08 | AC-HND-08 | 13 |

FR ทุกข้อมี AC แต่ส่วนยืมแบบค้างคืนของ FR-HND-01, 05, 06 ไม่มี AC ควรเสนอทีมเพิ่ม (รอ Q-HND-07)

---

## สิ่งที่ยังไม่ทำ (Open Questions จาก spec และของกลาง)
| คำถาม | task ที่รอ |
|---|---|
| Q-HND-03 ส่งมอบยืมปกติหลัง 16:00 กำหนดคืนคือเวลาไหน (ถ้า 16:00 วันเดียวกัน ผู้ยืมเกินกำหนดทันทีและ FR-BRQ-06 ของ 002 ปฏิเสธคำขอใหม่) -> นักคอมพิวเตอร์ | T-25 |
| Q-HND-04 "ใช้งานได้เฉพาะภายในพื้นที่มหาวิทยาลัย" ระบบต้องทำอะไร -> นักคอมพิวเตอร์ | T-31 |
| Q-HND-05 [ปิดแล้ว] `approved` -> `handed_over`, `requester_id`, `item_id`, `end_at` มีใน migration 0003 | ไม่มี task รอ |
| Q-HND-06 ค่า `action` ของ audit log สำหรับการส่งมอบ -> จองในกลุ่ม | T-26 |
| Q-HND-07 ต้องทดสอบยืมแบบค้างคืนเทอมนี้ไหม (ขึ้นกับว่า 004 ทำ UC-14 ไหม) -> ทีม + เจ้าของ 004 | T-30 |
| Q-HND-09 เวลาตอบสนองของการบันทึกการส่งมอบ -> เจ้าหน้าที่ | ไม่มี task (ไม่มี NFR ID รองรับ จะเสนอเป็น NFR เมื่อได้คำตอบ) |
| Q-HND-10 10 MB ต่อรูปหรือรวม และมากสุดกี่รูป -> เจ้าหน้าที่ | T-28 |
| Q-HND-11 [ปิดแล้ว] ผู้ยืมล็อกอินบนเครื่องเจ้าหน้าที่ (ASM-HND-09) | ไม่มี task รอ |
| Q-HND-13 [ปิดแล้ว] เจ้าหน้าที่พิมพ์ชื่อตามบัตร (ASM-HND-10) | ไม่มี task รอ |
| Q-HND-14 ส่งมอบไม่สำเร็จต้องบันทึก audit log ไหม -> เจ้าหน้าที่ | T-27 |
| Q-HND-15 ย้าย CON-HND-01 ไป constraints.md ไหม -> ทีม | ไม่มี task (ไม่กระทบโค้ด) |
| Q-SH-03 auth จริง (ของกลาง) | T-29 |
| Q-SH-04 [ปิดแล้ว] กำหนดคืนยืมค้างคืน = 16:00 ของวันสิ้นสุดการยืม | ไม่มี task รอ (T-30 ใช้คำตอบนี้) |
| Q-SH-06 [ปิดบางส่วน] 004 ปฏิเสธคำขออื่นของชิ้นเดียวกันอัตโนมัติเมื่ออนุมัติ ทำให้ Given ของ AC-HND-13 ไม่เกิดในระบบจริง -> ทีมตัดสินแก้หรือยกเลิก AC-HND-13 (ผ่าน CR) | T-19 (ทำได้ด้วยข้อมูลที่ตั้งเอง แต่อาจถูกยกเลิก) |
| Q-SH-08 ระบบยืนยันตัวตนเชื่อมต่อไม่ได้ตอนเข้าใช้ทั่วไป | ไม่มี task (ใช้ middleware เดิมของ 001) |
| Q-P2 (plan) ส่งมอบหลัง 16:00 ขัดกับ Email ของ 004 / จำกัดเจ้าหน้าที่ตามภาควิชาไหม -> ทีม | T-25 (ส่วนเวลา) ส่วนภาควิชาไม่มี task เพราะไม่มีใน spec |
| Q-P3 [หมดไปแล้ว] 004 spec v1 ไม่มีที่เก็บกำหนดคืนตอนอนุมัติ 005 เป็นที่เดียวที่บันทึกกำหนดคืน | ไม่มี task รอ |
| Q-P4 [ปิดแล้ว] ลงทะเบียน `FORBIDDEN_ROLE`, `REQUEST_NOT_APPROVED`, `BORROWER_MISMATCH` ใน CON-TECH-03 แล้ว | ไม่มี task รอ |

นอกจากนี้ยังต้องทำก่อนเริ่ม (ไม่ใช่ Open Question): แจ้งกลุ่มว่า 005 ใช้ migration 0006, 0007 (T-02, T-26), merge 003 ก่อนแตะ `App.jsx`/`client.js`, แจ้ง 006 เรื่อง volume `photos`
