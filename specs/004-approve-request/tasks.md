# Tasks: อนุมัติคำขอยืมปกติ (Approve Request)
Feature: 004-approve-request | Spec ID: SPEC-APV-004 (spec.md Draft v2.1)
อ้างอิง: plan.md Draft v1.1, specs/000-shared/{constraints.md v5 (CON-TECH-03 รหัสของ 004 ลงทะเบียน 2569-11-12), data-model.md v9, glossary.md}
Updated: 2569-11-12

สรุป: **20 เฟส 39 task** แบ่งเป็นเฟส 0 ฐานราก, เฟส 1 ถึง 18 (เฟสละ 1 AC ครบทั้ง 18 AC) และเฟส 19 รอคำตอบ
**3 task รอ Open Questions** (T-37 ถึง T-39 ใน Q-APV-04, Q-SH-03) ส่วน T-21, T-22, T-30, T-31 ทำได้เลยแต่ ID ของ AC ยังรอทีมอนุมัติ (Q-APV-13, Q-APV-22)
**เฟสแรกที่ลองใช้งานได้จริงคือเฟส 1** เจ้าหน้าที่เปิดหน้ารายการคำขอบนหน้าจอได้ และเฟส 2 เป็นเฟสแรกที่อนุมัติได้จริง

เงื่อนไขก่อนเริ่มเฟส 0 (ไม่ใช่ Open Question ของ spec แต่ต้องตกลงในกลุ่ม ตาม plan ข้อ 8):
- ~~จองเลข migration ของ 004~~ ได้เลข **0008** แล้ว (ทีมแจ้ง 2569-11-12)
- ~~เพิ่มค่า action ในตาราง DM-07 ของ data-model.md~~ ทีมตัดสิน 2569-11-12 ยังไม่ต้องเขียนใน data-model.md ค่า `approve`, `reject`, `auto_reject` และเลข 0008 กำหนดใน migration 0008 ที่ฐานข้อมูลเท่านั้น

ลำดับการทดสอบ (ทีมตัดสิน 2569-11-12): implement ทุก task ก่อน แล้วค่อยรัน test และลองตามจุดตรวจทีเดียวตอนท้าย
ระหว่าง implement ไม่รัน test (เครื่องยังไม่มี Go และ Docker) task ที่เขียนเสร็จใช้สถานะ "เสร็จ รอทีมตรวจ" พร้อมหมายเหตุว่ายังไม่ได้รัน test

วิธีเปิดระบบที่ใช้ทุกจุดตรวจ:
- หลังบ้าน: `docker compose -f deploy/docker-compose.yml up -d db` แล้ว `cd backend && go run ./cmd/server` (port 8000)
- ข้อมูลตัวอย่างของ 004 (โหลดเองหลังหลังบ้านรัน migration และ seed.sql แล้ว แบบเดียวกับ `handover_dev.sql` ของ 005):
  `docker compose -f deploy/docker-compose.yml exec -T db psql -U salmon -d salmon < backend/db/seed/approve_dev.sql`
- หน้าจอ (ในฐานะเจ้าหน้าที่ภาควิชา 1): `cd frontend && VITE_MOCK_TOKEN=mock-staff-001 npm run dev` แล้วเปิด `/approve.html`
- test: `cd backend && go test ./... -v` และ `cd frontend && npm test`

---

## เฟส 0: ฐานราก (ไม่มี AC เป็นงานพื้นฐานของทุกเฟส)

### T-01 สร้าง migration ของ 004 (ตาราง approval_decision, approval_email และค่า action ใหม่)
- รองรับ: FR-APV-01 (ฌ), FR-APV-02 (ง), CON-TECH-01, RQ-25
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-07 ถึง T-36
- ไฟล์ที่แตะ: `backend/db/migrations/0008_approval.up.sql`, `backend/db/migrations/0008_approval.down.sql`, `backend/internal/platform/migrate_approval_test.go` (ไฟล์ test ใหม่ ทีมตกลง 2569-11-12), `backend/internal/platform/migrate_test.go` (แก้ 1 บรรทัดใน test ของ 001 ที่ยืนยันว่า 'approve' ต้องถูกปฏิเสธ ทีมตกลง 2569-11-12)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: migration 0008 รันต่อจาก migration ที่มีอยู่ผ่าน (ตอนนี้ใน repo มีถึง 0003 เลข 0004 ถึง 0007 เป็นของฟีเจอร์อื่น) และ `audit_log_action_check` ที่ ADD ใหม่ต้องรวมค่า action ของทุก migration ก่อนหน้า และ test ยืนยันว่า (ก) `approval_decision.request_id` UNIQUE และมี FK ไป `borrow_request` (ข) `result` / `approval_email.kind` / `approval_email.status` รับเฉพาะค่าตาม plan ข้อ 3 (ค) `audit_log.action` รับ `approve`, `reject`, `auto_reject` ได้ (DROP แล้ว ADD `audit_log_action_check` ไม่แก้ migration เดิม)
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ทีมให้ implement ก่อนแล้วค่อยรัน test ยังไม่ได้รัน `TestMigration_ApprovalTables`, `TestMigration_ApprovalAuditActions`)

### T-02 เพิ่มค่าที่ตั้งได้ของ 004 ใน config
- รองรับ: FR-APV-01 (ค)(จ)(ฉ)(ซ), FR-APV-02, DOM-LOAN-01, CON-TECH-03
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-07 ถึง T-36
- ไฟล์ที่แตะ: `backend/internal/config/config.go`, `frontend/src/config.js`, `backend/internal/config/config_approval_test.go` (ไฟล์ test ใหม่ของ 004 แยกจาก config_test.go ของ 002 ตามแนวที่ทีมตกลงใน T-01)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: ค่าครบทุกแถวในตาราง "ค่าที่ตั้งได้" ของ plan ข้อ 5 (รวมข้อความ error 9 ข้อความ และ `ApprovalEmailRetryCount` = 3 นับเฉพาะครั้งที่ส่งซ้ำ) อยู่ที่เดียวพร้อมคอมเมนต์อ้าง ID และค้นในโค้ดแล้วไม่เจอค่าเหล่านี้ที่อื่น / ค่าที่ 002 ประกาศไว้แล้วใช้ร่วม ไม่ประกาศซ้ำ: เวลาตัด 16:00 น. ใช้ `LoanNormalReturnHour` (DOM-LOAN-01) และค่าสถานะของคำขอใช้ `RequestStatuses` / `StatusRequested`
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestConfig_ApprovalValues`, `TestConfig_ApprovalMessagesMatchFrontend` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-03 เพิ่มข้อมูลตัวอย่างสำหรับลองระบบของ 004
- รองรับ: FR-APV-01 (ก)(ข)(จ)
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นข้อมูลที่จุดตรวจของทุกเฟสใช้
- ไฟล์ที่แตะ: `backend/db/seed/approve_dev.sql` (ไฟล์ใหม่ โหลดเองด้วย psql ไม่แตะ `seed.sql` เพราะ `TestSeed_BorrowRequestScenarios` ของ 002 บังคับว่า seed.sql ต้องไม่มีคำขอยืม ทำแบบเดียวกับ `handover_dev.sql` ของ 005), `backend/internal/platform/migrate_approval_test.go` (ไฟล์ test ของ 004 จาก T-01 เพิ่ม `TestSeed_ApproveDevScenarios`)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: โหลด `approve_dev.sql` หลัง seed.sql ได้โดยไม่ error (โหลดร่วมกับ `handover_dev.sql` ได้) เพิ่มรุ่นและอุปกรณ์ชุดของตัวเอง (รหัสครุภัณฑ์ขึ้นต้น `AP-TEST-` ไม่แตะอุปกรณ์ของ seed.sql) และมีคำขอ `requested` ยืมปกติของภาควิชา 1 อย่างน้อย 3 ใบ (2 ใบอยู่บนอุปกรณ์ชิ้นเดียวกัน), คำขอค้างคืน 1 ใบบนชิ้นเดียวกันนั้น, คำขอของอุปกรณ์ภาควิชา 2 อย่างน้อย 1 ใบ, คำขอ 1 ใบบนอุปกรณ์ที่ `borrowed` และคำขอ `cancelled` 1 ใบ / seed.sql ไม่เปลี่ยน และ `TestSeed_BorrowRequestScenarios` ยังผ่าน
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestSeed_ApproveDevScenarios` และยังไม่ได้โหลดจริงด้วย psql ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 0
- test ที่ต้องผ่านทั้งหมด: `TestMigration_ApprovalTables`, `TestMigration_ApprovalAuditActions` และ test เดิมของ 001 / 002 / 005 ทั้งหมดต้องยังผ่าน (รวม `TestSeed_BorrowRequestScenarios`)
- เปิดระบบ: `docker compose -f deploy/docker-compose.yml up -d db` แล้ว `cd backend && go test ./... -v` และโหลด `approve_dev.sql` ตามหัวไฟล์
- คนลองเองตาม AC: ยังไม่มี AC ในเฟสนี้ ให้ลองว่า
  Given รัน migration ครบทุกไฟล์ seed.sql และ approve_dev.sql แล้ว
  When ใช้ psql `INSERT` แถว `approval_decision` ซ้ำ request_id เดิม 2 แถว
  Then ฐานข้อมูลต้องปฏิเสธแถวที่ 2 ด้วย UNIQUE constraint
- ยังขาด: ไม่มี

---

## เฟส 1: เจ้าหน้าที่เห็นรายการคำขอรอพิจารณาของภาควิชาตนเอง (AC-APV-11)

### T-04 สร้าง repository อ่านรายการคำขอรอพิจารณา
- รองรับ: FR-APV-01 (ก)(ข), DOM-REQ-01, ASM-06
- ตรวจด้วย: AC-APV-11
- ไฟล์ที่แตะ: `backend/internal/approve/model.go`, `backend/internal/approve/repository.go`, `backend/internal/approve/repository_test.go`, `backend/internal/config/config.go` (เพิ่ม `LoanTypeNormal` ตามกฎข้อ 9)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01, T-03
- เสร็จเมื่อ: query คืนเฉพาะคำขอ `requested` + `loan_type = 'normal'` ของภาควิชาที่ระบุ เรียงตาม `submitted_at` พร้อม `is_first` ถูกต้องต่ออุปกรณ์รายชิ้น และ `TestAC_APV_11` ระดับ repository ผ่าน
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_11`, `TestListPending_OvernightDoesNotBlockQueue` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-05 สร้าง GET /v1/approvals/requests พร้อมตัวตรวจบทบาทและผูก route
- รองรับ: FR-APV-01 (ก), IF-IDP-01, NFR-SEC-01, CON-TECH-02, CON-TECH-03
- ตรวจด้วย: AC-APV-11
- ไฟล์ที่แตะ: `backend/internal/approve/handler.go`, `backend/internal/approve/dto.go`, `backend/internal/approve/service.go`, `backend/internal/approve/handler_test.go`, `backend/cmd/server/main.go` (ของกลาง บันทึกใน prompt-log), `backend/internal/config/config.go` + `frontend/src/config.js` + `backend/internal/config/config_approval_test.go` (ข้อความ UPSTREAM_UNAVAILABLE ตาม Q-APV-24)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_APV_11` ผ่านทาง HTTP คืนรายการพร้อม `loaded_at` และทุก endpoint ใต้ `/v1/approvals/` ตอบ `FORBIDDEN_ROLE` เมื่อไม่ใช่เจ้าหน้าที่ (`{code, message}` ตาม CON-TECH-03)
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_11_HTTP`, `TestApprovals_RequiresStaff` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-06 สร้างหน้าจอ ApproveRequests แสดงรายการคำขอ
- รองรับ: FR-APV-01 (ก)
- ตรวจด้วย: AC-APV-11
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `getApprovalRequests`), `frontend/approve.html` + `frontend/src/approve-main.jsx` (หน้าแยกแบบ handover.html ของ 005 ไม่แตะ App.jsx ของ 001 ทีมตัดสิน 2569-11-12)
- แบบหน้าจอ: ไม่มี (ไม่มี mockup และไม่มี UI ID ใน spec ทีมออกแบบเอง ตาม plan ข้อ 4)
- ต้องทำหลัง: T-05
- เสร็จเมื่อ: vitest `AC-APV-11` ผ่านด้วย client จำลอง แสดงชื่อผู้ขอ อุปกรณ์ วัตถุประสงค์ ช่วงเวลาที่ต้องการยืม เรียงตามเวลาที่ส่ง
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน vitest `AC-APV-11` ตามลำดับการทดสอบที่ทีมตัดสิน และเครื่องยังไม่มี node_modules)

#### จุดตรวจเฟส 1
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_11` (repository + HTTP), vitest `AC-APV-11` และ test ของเฟส 0
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์ (`VITE_MOCK_TOKEN=mock-staff-001`)
- คนลองเองตาม AC-APV-11:
  Given โหลด approve_dev.sql จาก T-03 และเปิดหน้าจอในฐานะ `mock-staff-001` (ภาควิชา 1)
  When เปิดหน้ารายการคำขอ
  Then เห็นเฉพาะคำขอยืมปกติ `requested` ของภาควิชา 1 เรียงจากเก่าไปใหม่ ไม่เห็นคำขอของภาควิชา 2 คำขอค้างคืน และคำขอ `cancelled`
- ยังขาด: ไม่มี

---

## เฟส 2: อนุมัติคำขอแรกสำเร็จ (AC-APV-01)

### T-07 สร้าง service อนุมัติ (ล็อกแถว ตรวจสถานะ บันทึกผลและประวัติ)
- รองรับ: FR-APV-01 (ค)(ฌ), CON-DATA-01, DOM-REQ-01, RQ-25
- ตรวจด้วย: AC-APV-01
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/repository.go`, `backend/internal/approve/service_test.go`, `backend/internal/approve/model.go` (เพิ่ม requestRow, decision, Result)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-02, T-04
- เสร็จเมื่อ: `TestAC_APV_01` ผ่าน: ใน transaction เดียวล็อกคำขอทุกใบ (ทุกประเภทการยืม) ของอุปกรณ์ชิ้นนั้นด้วย `SELECT ... FOR UPDATE` คำขอเป็น `approved` มีแถว `approval_decision` และ `audit_log` action `approve` และ `equipment_item.status` ยังเป็น `available`
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_01`, `TestApprove_OnlyRequestedNormal` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-08 สร้าง POST /v1/approvals/requests/{id}/approve
- รองรับ: FR-APV-01, CON-TECH-03, IF-IDP-01
- ตรวจด้วย: AC-APV-01
- ไฟล์ที่แตะ: `backend/internal/approve/handler.go`, `backend/internal/approve/dto.go`, `backend/internal/approve/handler_test.go`, `backend/internal/config/config.go` + `frontend/src/config.js` + `backend/internal/config/config_approval_test.go` (REQUEST_NOT_FOUND ตาม Q-APV-25)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-05, T-07
- เสร็จเมื่อ: `TestAC_APV_01` ทาง HTTP ผ่าน รับ `loaded_at` และคืนผลพร้อมรายการคำขอที่ถูกปฏิเสธตาม (ว่างในเฟสนี้)
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_01_HTTP`, `TestApproveHTTP_NotFound` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-09 เพิ่มปุ่มอนุมัติบนหน้าจอ ApproveRequests
- รองรับ: FR-APV-01 (ข)(ค)
- ตรวจด้วย: AC-APV-01
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `postJSON`, `approveRequest`)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06, T-08
- เสร็จเมื่อ: vitest `AC-APV-01` ผ่าน กดอนุมัติแล้วคำขอหายจากรายการ ส่ง `loaded_at` ที่ได้จากการโหลดรายการ
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน vitest `AC-APV-01` ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 2
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_01` (service + HTTP), vitest `AC-APV-01` และ test ของเฟส 0 ถึง 1
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-01:
  Given approve_dev.sql มีคำขอแรกของอุปกรณ์ที่ "ว่าง" ในภาควิชา 1
  When กด "อนุมัติ" ที่คำขอนั้น
  Then คำขอหายจากรายการ เปิด psql เห็น `borrow_request.status = 'approved'` และ `equipment_item.status` ยังเป็น `available`
- ยังขาด: ไม่มี

---

## เฟส 3: ปฏิเสธคำขอพร้อมเหตุผล (AC-APV-02)

### T-10 สร้าง service และ POST /v1/approvals/requests/{id}/reject
- รองรับ: FR-APV-01 (ง)(ฌ), CON-DATA-01, RQ-25
- ตรวจด้วย: AC-APV-02
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/handler.go`, `backend/internal/approve/dto.go`, `backend/internal/approve/service_test.go`, `backend/internal/approve/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-08
- เสร็จเมื่อ: `TestAC_APV_02` ผ่าน: คำขอเป็น `rejected` บันทึกเหตุผลใน `approval_decision` มี `audit_log` action `reject` และสถานะอุปกรณ์ไม่เปลี่ยน
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_02`, `TestAC_APV_02_HTTP` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-11 เพิ่มปุ่มปฏิเสธพร้อมช่องเหตุผลบนหน้าจอ
- รองรับ: FR-APV-01 (ง)
- ตรวจด้วย: AC-APV-02
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `rejectRequest`)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-09, T-10
- เสร็จเมื่อ: vitest `AC-APV-02` ผ่าน พิมพ์เหตุผลแล้วกดปฏิเสธ คำขอหายจากรายการ
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน vitest `AC-APV-02` ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 3
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_02`, vitest `AC-APV-02` และ test ของเฟส 0 ถึง 2
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-02:
  Given approve_dev.sql มีคำขอ `requested` ของภาควิชา 1
  When กด "ปฏิเสธ" พิมพ์เหตุผล "อุปกรณ์ต้องใช้สอบ" แล้วยืนยัน
  Then คำขอหายจากรายการ psql เห็น `rejected` และเหตุผลตรงกับที่พิมพ์ สถานะอุปกรณ์ไม่เปลี่ยน
- ยังขาด: ไม่มี

---

## เฟส 4: ปฏิเสธโดยไม่ระบุเหตุผลไม่ได้ (AC-APV-06)

### T-12 ตรวจเหตุผลการปฏิเสธที่ว่างหรือมีแต่ช่องว่าง
- รองรับ: FR-APV-01 (ง)(ซ), CON-TECH-03
- ตรวจด้วย: AC-APV-06
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/handler_test.go`, `backend/internal/approve/handler.go` + `dto.go` (แปลง error เป็น 422 พร้อมรายชื่อช่อง)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `TestAC_APV_06` ผ่าน: ตอบ 422 `VALIDATION_FAILED` ระบุช่อง `reason` พร้อมข้อความ "กรุณาระบุเหตุผลในการปฏิเสธ" และสถานะคงเดิม
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_06` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-13 แสดงข้อความเมื่อไม่ระบุเหตุผลบนหน้าจอ
- รองรับ: FR-APV-01 (ซ)
- ตรวจด้วย: AC-APV-06
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-11, T-12
- เสร็จเมื่อ: vitest `AC-APV-06` ผ่าน แสดง "กรุณาระบุเหตุผลในการปฏิเสธ" จาก `APPROVAL_MESSAGES` และคำขอยังอยู่ในรายการ
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน vitest `AC-APV-06` ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 4
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_06`, vitest `AC-APV-06` และ test ของเฟส 0 ถึง 3
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-06:
  Given approve_dev.sql มีคำขอ `requested`
  When กด "ปฏิเสธ" เว้นช่องเหตุผลว่าง (หรือพิมพ์แต่ช่องว่าง) แล้วยืนยัน
  Then เห็น "กรุณาระบุเหตุผลในการปฏิเสธ" และคำขอยังอยู่ในรายการ
- ยังขาด: ไม่มี

---

## เฟส 5: ผู้ขอได้รับ Email ผลอนุมัติหรือปฏิเสธภายใน 60 วินาที (AC-APV-04)

### T-14 สร้าง Mailer, RecipientResolver (ตัวจำลอง) และตัวสร้าง Email 3 แบบ
- รองรับ: FR-APV-02 (ก)(ข)(ค), IF-NOT-01
- ตรวจด้วย: AC-APV-04
- ไฟล์ที่แตะ: `backend/internal/approve/mailer.go`, `backend/internal/approve/email.go`, `backend/internal/approve/email_test.go`, `backend/internal/config/config.go` (เพิ่ม `ApprovalPurposeLabels`)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-02
- เสร็จเมื่อ: test ยืนยันว่า Email แต่ละแบบมีหัวข้อและข้อมูลครบตาม FR-APV-02 (ก)(ข)(ค) จาก config และ Email ผลอนุมัติไม่มีจำนวนอุปกรณ์ / `RecipientResolver` ตัวจำลองคืนที่อยู่ตายตัวต่อผู้ขอ (ตัวจริงรอ Q-APV-04 ดู T-37)
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestBuildEmail_Contents`, `TestBuildEmail_PurposeOther`, `TestMockRecipientResolver` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-15 เขียน outbox ใน transaction เดียวกับผล และตัวส่งเบื้องหลัง
- รองรับ: FR-APV-02, IF-NOT-01
- ตรวจด้วย: AC-APV-04
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/sender.go`, `backend/internal/approve/sender_test.go`, `backend/internal/approve/repository.go` (insertEmail), `backend/internal/config/config.go` (สถานะการส่ง, ช่วงตรวจคิว, จำนวนต่อรอบ)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10, T-14
- เสร็จเมื่อ: `TestAC_APV_04` ผ่านด้วยนาฬิกาจำลอง: หลังอนุมัติหรือปฏิเสธมีแถว `approval_email` สถานะ `pending` แล้วตัวส่งจำลองได้รับ Email ภายใน `ApprovalEmailDeadline` และแถวเป็น `sent`
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_04`, `TestSender_PollWithinDeadline` ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 5
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_04`, test ของ `email_test.go` และ test ของเฟส 0 ถึง 4
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์ (ตัวส่งจำลองเขียน Email ลง log ของหลังบ้าน)
- คนลองเองตาม AC-APV-04:
  Given approve_dev.sql มีคำขอ `requested`
  When อนุมัติคำขอหนึ่ง และปฏิเสธอีกคำขอพร้อมเหตุผล
  Then ภายใน 60 วินาที log ของหลังบ้านมี Email 2 ฉบับ หัวข้อ "แจ้งผลการอนุมัติการยืมอุปกรณ์ภาคคอมพิวเตอร์" และ "แจ้งผลการปฏิเสธการยืมอุปกรณ์ภาคคอมพิวเตอร์" ข้อมูลครบตาม FR-APV-02
- ยังขาด: ส่งถึง Email จริงของผู้ขอ (รอ Q-APV-04 ดู T-37) และชื่อผู้อนุมัติใน Email (Principal ยังไม่มีชื่อ ดู T-38)

---

## เฟส 6: อนุมัติแล้วคำขออื่นของชิ้นเดียวกันถูกปฏิเสธตาม (AC-APV-03)

### T-16 ปฏิเสธคำขออื่นของอุปกรณ์ชิ้นเดียวกันทุกประเภทใน transaction ของการอนุมัติ
- รองรับ: FR-APV-01 (จ)(ฌ), FR-APV-02 (ข), DOM-REQ-01, Q-SH-06, Q-APV-14
- ตรวจด้วย: AC-APV-03
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/service_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07, T-15
- เสร็จเมื่อ: `TestAC_APV_03` ผ่าน: B เป็น `rejected` เหตุผล `ApprovalSupersededReason` ผู้ดำเนินการคือเจ้าหน้าที่ที่อนุมัติ A มี `audit_log` action `reject` และ Email ผลปฏิเสธของ B
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_03` ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 6
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_03` และ test ของเฟส 0 ถึง 5
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-03:
  Given approve_dev.sql มีคำขอ A (ส่งก่อน) และ B บนอุปกรณ์ชิ้นเดียวกัน
  When อนุมัติ A
  Then ทั้ง A และ B หายจากรายการ psql เห็น B เป็น `rejected` ด้วยเหตุผล "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว" และ log มี Email ผลปฏิเสธของ B
- ยังขาด: ไม่มี

---

## เฟส 7: อนุมัติข้ามลำดับไม่ได้ (AC-APV-07)

### T-17 ตรวจคำขอแรกก่อนอนุมัติ
- รองรับ: FR-APV-01 (ข)(ซ), DOM-REQ-01, ASM-04, CON-TECH-03
- ตรวจด้วย: AC-APV-07
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/service_test.go`, `backend/internal/approve/handler.go` (แปลงเป็น 409 NOT_FIRST_IN_QUEUE)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07
- เสร็จเมื่อ: `TestAC_APV_07` ผ่าน: อนุมัติ B ก่อน A ได้ 409 `NOT_FIRST_IN_QUEUE` ข้อความ "ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน" ทั้งสองคำขอคงเดิม (ปฏิเสธ B ได้ตามปกติ)
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_07` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-18 แสดงปุ่มอนุมัติเฉพาะคำขอแรก และแสดงข้อความเมื่อได้ NOT_FIRST_IN_QUEUE
- รองรับ: FR-APV-01 (ข)(ซ)
- ตรวจด้วย: AC-APV-07
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-09, T-17
- เสร็จเมื่อ: vitest `AC-APV-07` ผ่าน: แถวที่ `is_first = false` ไม่มีปุ่มอนุมัติ และเมื่อ client จำลองตอบ `NOT_FIRST_IN_QUEUE` หน้าจอแสดงข้อความตาม config
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน vitest `AC-APV-07` ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 7
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_07`, vitest `AC-APV-07` และ test ของเฟส 0 ถึง 6
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-07:
  Given approve_dev.sql มีคำขอ A และ B บนอุปกรณ์ชิ้นเดียวกัน (ถ้าเฟส 6 ใช้ไปแล้ว ให้ล้างฐานข้อมูลแล้วโหลด approve_dev.sql ใหม่)
  When ดูแถวของ B แล้วยิง `curl -X POST localhost:8000/v1/approvals/requests/<id ของ B>/approve -H "Authorization: Bearer mock-staff-001" -d '{"loaded_at":"..."}'`
  Then บนหน้าจอแถว B ไม่มีปุ่มอนุมัติ และ curl ได้ 409 `NOT_FIRST_IN_QUEUE`
- ยังขาด: ไม่มี

---

## เฟส 8: อนุมัติไม่ได้เมื่ออุปกรณ์ไม่ "ว่าง" (AC-APV-08)

### T-19 ตรวจสถานะอุปกรณ์ก่อนอนุมัติ
- รองรับ: FR-APV-01 (ค)(ซ), CON-DATA-01, CON-TECH-03
- ตรวจด้วย: AC-APV-08
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/service_test.go`, `backend/internal/approve/handler.go` (แปลงเป็น 409 ITEM_NOT_AVAILABLE), `backend/internal/config/config.go` (`ItemStatusAvailable`)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07
- เสร็จเมื่อ: `TestAC_APV_08` ผ่าน: อุปกรณ์ `borrowed` หรือ `unavailable` ได้ 409 `ITEM_NOT_AVAILABLE` ข้อความตาม config สถานะคงเดิม
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_08` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-20 แสดงข้อความ ITEM_NOT_AVAILABLE บนหน้าจอ
- รองรับ: FR-APV-01 (ซ)
- ตรวจด้วย: AC-APV-08
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-09, T-19
- เสร็จเมื่อ: vitest `AC-APV-08` ผ่าน แสดง "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้"
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน vitest `AC-APV-08` ตามลำดับการทดสอบที่ทีมตัดสิน / ApproveRequests.jsx ไม่ต้องแก้ เพราะแสดงข้อความตามรหัสอยู่แล้วตั้งแต่ T-09)

#### จุดตรวจเฟส 8
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_08`, vitest `AC-APV-08` และ test ของเฟส 0 ถึง 7
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-08:
  Given approve_dev.sql มีคำขอแรกของอุปกรณ์ที่ `borrowed`
  When กด "อนุมัติ" ที่คำขอนั้น
  Then เห็น "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้" และคำขอยังเป็น `requested`
- ยังขาด: ไม่มี

---

## เฟส 9: อนุมัติไม่ได้เมื่อเลย 16:00 น. ของวันที่เริ่มยืม (AC-APV-16)

### T-21 ตรวจเวลาตัด 16:00 น. ก่อนอนุมัติ
- รองรับ: FR-APV-01 (ค)(ซ), DOM-LOAN-01, Q-APV-11, Q-APV-16
- ตรวจด้วย: AC-APV-16
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/service_test.go`, `backend/internal/approve/handler.go` (แปลงเป็น 409 APPROVAL_CUTOFF_PASSED)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07
- เสร็จเมื่อ: `TestAC_APV_16` ผ่านด้วยนาฬิกาจำลอง: เลย `LoanNormalReturnHour` (16:00 น.) ของวันที่เริ่มยืม (Asia/Bangkok) ได้ 409 `APPROVAL_CUTOFF_PASSED` สถานะคงเดิม คำขอนั้นยังเป็นคำขอแรก (คำขอถัดไปยังได้ `NOT_FIRST_IN_QUEUE`) และปฏิเสธเองได้ (ID ของ AC-APV-16 รอ Q-APV-13 ถ้าทีมเปลี่ยน ID ให้เปลี่ยนชื่อ test ตาม)
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_16`, `TestApprove_BeforeCutoffAllowed` ตามลำดับการทดสอบที่ทีมตัดสิน / ID ของ AC-APV-16 ยังรอ Q-APV-13)

### T-22 แสดงข้อความ APPROVAL_CUTOFF_PASSED บนหน้าจอ
- รองรับ: FR-APV-01 (ซ)
- ตรวจด้วย: AC-APV-16
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-09, T-21
- เสร็จเมื่อ: vitest `AC-APV-16` ผ่าน แสดง "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้"
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน vitest `AC-APV-16` ตามลำดับการทดสอบที่ทีมตัดสิน / ApproveRequests.jsx ไม่ต้องแก้ เพราะแสดงข้อความตามรหัสอยู่แล้วตั้งแต่ T-09)

#### จุดตรวจเฟส 9
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_16`, vitest `AC-APV-16` และ test ของเฟส 0 ถึง 8
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-16:
  Given ใช้ psql แก้ `start_at` ของคำขอแรกเป็นเมื่อวาน (เลย 16:00 น. ของวันที่เริ่มยืมแน่นอน)
  When กด "อนุมัติ" ที่คำขอนั้น
  Then เห็น "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้" และคำขอยังเป็น `requested`
- ยังขาด: ทีมยังไม่อนุมัติ ID ของ AC-APV-16 (Q-APV-13)

---

## เฟส 10: พิจารณาคำขอที่ถูกยกเลิกหรือพิจารณาแล้วไม่ได้ (AC-APV-09)

### T-23 ตรวจสถานะคำขอก่อนพิจารณา และแยกพิจารณาแล้วด้วย loaded_at
- รองรับ: FR-APV-01 (ซ), CON-TECH-03
- ตรวจด้วย: AC-APV-09
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/service_test.go`, `backend/internal/approve/handler.go` (ส่ง loadedAt และแปลง 3 รหัส), `backend/internal/approve/repository.go` (decisionTime), `backend/internal/config/config.go` (StatusCancelled), `sender_test.go` (เพิ่มอาร์กิวเมนต์ loadedAt)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `TestAC_APV_09` ผ่าน: `cancelled` ได้ 409 `REQUEST_CANCELLED` / พิจารณาก่อน `loaded_at` ได้ 409 `REQUEST_ALREADY_DECIDED` / พิจารณาหลัง `loaded_at` ได้ `REQUEST_DECIDED_CONCURRENTLY` ทั้งอนุมัติและปฏิเสธ ผลเดิมคงไว้
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_09` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-24 แสดงข้อความคำขอถูกยกเลิก / พิจารณาแล้วบนหน้าจอ
- รองรับ: FR-APV-01 (ซ)
- ตรวจด้วย: AC-APV-09
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-11, T-23
- เสร็จเมื่อ: vitest `AC-APV-09` ผ่าน แสดง "คำขอนี้ถูกผู้ขอยกเลิกแล้ว" หรือ "คำขอนี้ได้รับการพิจารณาแล้ว" ตามรหัส
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน vitest `AC-APV-09` ตามลำดับการทดสอบที่ทีมตัดสิน / ApproveRequests.jsx ไม่ต้องแก้ เพราะแสดงข้อความตามรหัสอยู่แล้วตั้งแต่ T-09)

#### จุดตรวจเฟส 10
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_09`, vitest `AC-APV-09` และ test ของเฟส 0 ถึง 9
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-09:
  Given เปิดหน้ารายการค้างไว้ แล้วใช้ psql เปลี่ยนคำขอหนึ่งเป็น `cancelled`
  When กด "อนุมัติ" ที่คำขอนั้นโดยไม่รีเฟรช
  Then เห็น "คำขอนี้ถูกผู้ขอยกเลิกแล้ว" และสถานะยังเป็น `cancelled`
- ยังขาด: ไม่มี

---

## เฟส 11: พิจารณาคำขอของอุปกรณ์ภาควิชาอื่นไม่ได้ (AC-APV-12)

### T-25 ตรวจภาควิชาของเจ้าหน้าที่กับภาควิชาเจ้าของอุปกรณ์ก่อนพิจารณา
- รองรับ: FR-APV-01 (ก)(ซ), IF-IDP-01, ASM-06, CON-TECH-03
- ตรวจด้วย: AC-APV-12
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/handler_test.go`, `backend/internal/approve/handler.go` (แปลงเป็น 403 FORBIDDEN_DEPARTMENT)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `TestAC_APV_12` ผ่าน: `mock-staff-001` (ภาควิชา 1) อนุมัติหรือปฏิเสธคำขอของอุปกรณ์ภาควิชา 2 ได้ 403 `FORBIDDEN_DEPARTMENT` ข้อความ "คำขอนี้เป็นของอุปกรณ์ภาควิชาอื่น คุณไม่มีสิทธิ์พิจารณา" ไม่ส่งข้อมูลอุปกรณ์กลับ สถานะคงเดิม
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_12` ตามลำดับการทดสอบที่ทีมตัดสิน)

### T-26 แสดงข้อความ FORBIDDEN_DEPARTMENT บนหน้าจอ
- รองรับ: FR-APV-01 (ซ)
- ตรวจด้วย: AC-APV-12
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-11, T-25
- เสร็จเมื่อ: vitest `AC-APV-12` ผ่าน เมื่อ client จำลองตอบ `FORBIDDEN_DEPARTMENT` หน้าจอแสดงข้อความตาม config
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน vitest `AC-APV-12` ตามลำดับการทดสอบที่ทีมตัดสิน / ApproveRequests.jsx ไม่ต้องแก้ เพราะแสดงข้อความตามรหัสอยู่แล้วตั้งแต่ T-09)

#### จุดตรวจเฟส 11
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_12`, vitest `AC-APV-12` และ test ของเฟส 0 ถึง 10
- เปิดระบบ: หลังบ้านตามหัวไฟล์ (หน้าจอไม่แสดงคำขอภาควิชาอื่นอยู่แล้ว จึงลองด้วย curl)
- คนลองเองตาม AC-APV-12:
  Given approve_dev.sql มีคำขอของอุปกรณ์ภาควิชา 2
  When ยิง `curl -X POST localhost:8000/v1/approvals/requests/<id>/approve -H "Authorization: Bearer mock-staff-001" -d '{"loaded_at":"..."}'`
  Then ได้ 403 `{"code":"FORBIDDEN_DEPARTMENT","message":"คำขอนี้เป็นของอุปกรณ์ภาควิชาอื่น คุณไม่มีสิทธิ์พิจารณา"}` และคำขอยังเป็น `requested`
- ยังขาด: ไม่มี

---

## เฟส 12: ผู้ที่ไม่ได้ยืนยันตัวตนหรือนักศึกษาพิจารณาคำขอไม่ได้ (AC-APV-05)

### T-27 ทดสอบการตรวจตัวตนและบทบาทของ endpoint อนุมัติและปฏิเสธ
- รองรับ: NFR-SEC-01, IF-IDP-01, CON-TECH-03
- ตรวจด้วย: AC-APV-05
- ไฟล์ที่แตะ: `backend/internal/approve/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-08, T-10
- เสร็จเมื่อ: `TestAC_APV_05` ผ่าน: ไม่มี token ได้ 401 `UNAUTHENTICATED` / `mock-student-001` ได้ 403 `FORBIDDEN_ROLE` ทั้ง approve และ reject สถานะคำขอคงเดิม (401 มาจาก `auth.Middleware` ของกลางซึ่งยังไม่มี `message` ดูสิ่งที่ยังไม่ทำ)
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_05` ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 12
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_05` และ test ของเฟส 0 ถึง 11
- เปิดระบบ: หลังบ้านตามหัวไฟล์
- คนลองเองตาม AC-APV-05:
  Given approve_dev.sql มีคำขอ `requested`
  When ยิง curl approve ครั้งแรกไม่ใส่ `Authorization` และครั้งที่สองใส่ `Bearer mock-student-001`
  Then ได้ 401 `UNAUTHENTICATED` และ 403 `FORBIDDEN_ROLE` ตามลำดับ และคำขอยังเป็น `requested`
- ยังขาด: ไม่มี

---

## เฟส 13: ทุกการพิจารณามีประวัติ (AC-APV-15)

### T-28 ทดสอบประวัติของการอนุมัติ ปฏิเสธ และปฏิเสธตาม
- รองรับ: FR-APV-01 (ฌ), RQ-25, ASM-09
- ตรวจด้วย: AC-APV-15
- ไฟล์ที่แตะ: `backend/internal/approve/service_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-16
- เสร็จเมื่อ: `TestAC_APV_15` ผ่าน: อนุมัติ ปฏิเสธ และปฏิเสธตาม แต่ละกรณีมีแถว `approval_decision` (รหัสคำขอ ผล เหตุผล ผู้ดำเนินการ วันเวลา) และ `audit_log` action `approve` / `reject` ตรงกรณี (ชื่อเจ้าหน้าที่ยังว่างได้ ดู T-38) / กรณีปฏิเสธอัตโนมัติเพิ่มใน `TestAC_APV_15` ที่ T-29
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_15` ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 13
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_15` และ test ของเฟส 0 ถึง 12
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-15:
  Given ทำเฟส 2, 3 และ 6 บนหน้าจอแล้ว
  When เปิด psql `SELECT * FROM approval_decision` และ `SELECT * FROM audit_log WHERE action IN ('approve','reject','auto_reject')`
  Then ทุกการพิจารณามีแถวครบ ผู้ดำเนินการเป็น `staff-001`
- ยังขาด: ชื่อเจ้าหน้าที่ในประวัติ (Principal ยังไม่มีชื่อ ดู T-38) และกรณีปฏิเสธอัตโนมัติ (ทำในเฟส 14 T-29)

---

## เฟส 14: คำขอยืมปกติที่ค้างถึงเที่ยงคืนถูกปฏิเสธอัตโนมัติ (AC-APV-13)

### T-29 สร้างงานเที่ยงคืนปฏิเสธคำขอยืมปกติที่ค้าง
- รองรับ: FR-APV-01 (ฉ)(ฌ), FR-APV-02 (ค), ASM-05, Q-APV-14, Q-APV-15
- ตรวจด้วย: AC-APV-13, AC-APV-15 (กรณีปฏิเสธอัตโนมัติ)
- ไฟล์ที่แตะ: `backend/internal/approve/autoreject.go`, `backend/internal/approve/autoreject_test.go`, `backend/internal/approve/service_test.go`, `backend/internal/config/config.go` (`ApprovalSystemActorName`)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-15, T-28
- เสร็จเมื่อ: `TestAC_APV_13` ผ่านด้วยนาฬิกาจำลอง: ที่ 00:00 น. หลังวันที่ D คำขอยืมปกติ `requested` ที่ส่งวัน D เป็น `rejected` result `auto_rejected` ผู้ดำเนินการ `system` เหตุผล `ApprovalAutoRejectReason` และมี Email ปฏิเสธอัตโนมัติภายใน 60 วินาที / คำขอที่ส่งวัน D+1 และคำขอค้างคืนไม่ถูกแตะ / `TestAC_APV_15` เพิ่มกรณีปฏิเสธอัตโนมัติ (`audit_log` action `auto_reject`) แล้วผ่าน
- สถานะ: เสร็จ (ทีมตรวจแล้ว 2569-11-12 / ยังไม่ได้รัน `TestAC_APV_13`, `TestAC_APV_15` กรณีปฏิเสธอัตโนมัติ ตามลำดับการทดสอบที่ทีมตัดสิน)

#### จุดตรวจเฟส 14
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_13`, `TestAC_APV_15` (รวมกรณีปฏิเสธอัตโนมัติ) และ test ของเฟส 0 ถึง 13
- เปิดระบบ: หลังบ้านตามหัวไฟล์
- คนลองเองตาม AC-APV-13 (รอเที่ยงคืนจริงไม่สะดวก ให้ลองด้วยนาฬิกาจำลองใน test เป็นหลัก):
  Given ใช้ psql แก้ `submitted_at` ของคำขอยืมปกติหนึ่งใบเป็นเมื่อวาน แล้วเปิดหลังบ้านไว้ข้ามเที่ยงคืน (หรือใช้วิธีของ T-30)
  When ถึง 00:00 น.
  Then คำขอเป็น `rejected` ผู้ดำเนินการ `system` และ log มี Email หัวข้อ "แจ้งปฏิเสธคำขอยืมอุปกรณ์อัตโนมัติ"
- ยังขาด: ไม่มี

---

## เฟส 15: ระบบเริ่มทำงานหลังเที่ยงคืนแล้วไล่ปฏิเสธคำขอที่ค้าง (AC-APV-18 [เสนอ])

### T-30 ไล่ปฏิเสธคำขอที่เลยเวลาตอน service เริ่มทำงาน และเริ่ม goroutine ใน main.go
- รองรับ: FR-APV-01 (ฉ), FR-APV-02 (ค), Q-APV-21
- ตรวจด้วย: AC-APV-18
- ไฟล์ที่แตะ: `backend/internal/approve/autoreject.go`, `backend/internal/approve/autoreject_test.go`, `backend/cmd/server/main.go` (ของกลาง บันทึกใน prompt-log)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-29
- เสร็จเมื่อ: `TestAC_APV_18` ผ่าน: ตั้งนาฬิกาเลยเที่ยงคืนแล้วเริ่ม service คำขอยืมปกติที่ส่งวันก่อนเป็น `rejected` ทันทีพร้อม Email และ main.go เริ่มตัวส่ง Email กับงานเที่ยงคืน (ID ของ AC-APV-18 รอ Q-APV-22)
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 15
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_18` และ test ของเฟส 0 ถึง 14
- เปิดระบบ: หลังบ้านตามหัวไฟล์
- คนลองเองตาม AC-APV-18:
  Given ปิดหลังบ้าน แล้วใช้ psql แก้ `submitted_at` ของคำขอยืมปกติหนึ่งใบเป็นเมื่อวาน
  When เปิดหลังบ้านใหม่ `go run ./cmd/server`
  Then คำขอนั้นเป็น `rejected` ผู้ดำเนินการ `system` ทันที และ log มี Email ปฏิเสธอัตโนมัติ
- ยังขาด: ทีมยังไม่อนุมัติ ID ของ AC-APV-18 (Q-APV-22)

---

## เฟส 16: คำขอค้างคืนถูกปฏิเสธตาม แต่ไม่ถูกปฏิเสธตอนเที่ยงคืน (AC-APV-17 [เสนอ])

### T-31 ทดสอบคำขอค้างคืนในการปฏิเสธตามและงานเที่ยงคืน
- รองรับ: FR-APV-01 (จ)(ฉ), Q-APV-14
- ตรวจด้วย: AC-APV-17
- ไฟล์ที่แตะ: `backend/internal/approve/service_test.go`, `backend/internal/approve/autoreject_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-16, T-29
- เสร็จเมื่อ: `TestAC_APV_17` ผ่าน: อนุมัติคำขอปกติ A แล้วคำขอค้างคืน C ของชิ้นเดียวกันเป็น `rejected` ด้วย `ApprovalSupersededReason` และมี Email / ถ้าไม่มีการอนุมัติ ถึง 00:00 น. แล้ว C ยังเป็น `requested` (ID ของ AC-APV-17 รอ Q-APV-22)
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 16
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_17` และ test ของเฟส 0 ถึง 15
- เปิดระบบ: หลังบ้านและหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-17:
  Given approve_dev.sql มีคำขอยืมปกติ A และคำขอค้างคืน C บนอุปกรณ์ชิ้นเดียวกัน (C ไม่ขึ้นในหน้ารายการ)
  When อนุมัติ A บนหน้าจอ
  Then psql เห็น C เป็น `rejected` ด้วยเหตุผล "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว" และ log มี Email ผลปฏิเสธของ C
- ยังขาด: ทีมยังไม่อนุมัติ ID ของ AC-APV-17 (Q-APV-22)

---

## เฟส 17: Email ส่งไม่สำเร็จ ส่งซ้ำ และเจ้าหน้าที่กดส่งซ้ำได้ (AC-APV-14)

### T-32 ส่งซ้ำอัตโนมัติ 3 ครั้งหลังครั้งแรก และบันทึก Failed
- รองรับ: FR-APV-02 (ง), IF-NOT-01, Q-APV-17
- ตรวจด้วย: AC-APV-14
- ไฟล์ที่แตะ: `backend/internal/approve/sender.go`, `backend/internal/approve/sender_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-15
- เสร็จเมื่อ: test ด้วยนาฬิกาจำลองและตัวส่งที่ล้มทุกครั้ง: ตัวส่งถูกเรียกรวม 4 ครั้ง ห่างกัน `ApprovalEmailRetryInterval` แล้วแถวเป็น `failed` พร้อม `last_error` และ `last_attempt_at` ผลการพิจารณาไม่เปลี่ยน
- สถานะ: พร้อมทำ

### T-33 สร้าง GET /v1/approvals/failed-emails และ POST /v1/approvals/emails/{id}/resend
- รองรับ: FR-APV-02 (ง), ASM-08, ASM-06
- ตรวจด้วย: AC-APV-14
- ไฟล์ที่แตะ: `backend/internal/approve/handler.go`, `backend/internal/approve/dto.go`, `backend/internal/approve/repository.go`, `backend/internal/approve/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-05, T-32
- เสร็จเมื่อ: `TestAC_APV_14` ผ่าน: รายการคืนเฉพาะ Email `failed` ของภาควิชาตนเอง (รวมของ System) และกดส่งซ้ำแล้วเป็น `sent` โดย `approval_decision` ไม่เปลี่ยน
- สถานะ: พร้อมทำ

### T-34 สร้างหน้าจอ FailedEmails พร้อมปุ่มส่งซ้ำ
- รองรับ: FR-APV-02 (ง), ASM-08
- ตรวจด้วย: AC-APV-14
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/FailedEmails.jsx`, `frontend/src/pages/ApproveRequest/FailedEmails.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `getFailedEmails`, `resendEmail`), `frontend/src/approve-main.jsx` (แสดงต่อจาก ApproveRequests ในหน้า /approve.html)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-09, T-33
- เสร็จเมื่อ: vitest `AC-APV-14` ผ่าน แสดงรายการ "Email ที่ส่งไม่สำเร็จ" พร้อมสาเหตุและเวลา กดส่งซ้ำแล้วรายการหายไป
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 17
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_14`, test ของ `sender_test.go`, vitest `AC-APV-14` และ test ของเฟส 0 ถึง 16
- เปิดระบบ: หลังบ้านโดยตั้งตัวส่งจำลองให้ล้ม (ค่าใน config ของตัวส่งจำลอง ทีมเลือกเอง ไม่ได้มาจาก spec) และหน้าจอตามหัวไฟล์
- คนลองเองตาม AC-APV-14:
  Given ตัวส่งจำลองล้มทุกครั้ง
  When อนุมัติคำขอหนึ่ง รอประมาณ 3 นาที แล้วเปิดหน้า "Email ที่ส่งไม่สำเร็จ"
  Then เห็นคำขอนั้นพร้อมสาเหตุและเวลา จากนั้นตั้งตัวส่งให้กลับมาปกติ กดส่งซ้ำ แล้วรายการหายไปและ log มี Email ผลอนุมัติ คำขอยังเป็น `approved`
- ยังขาด: ไม่มี

---

## เฟส 18: เจ้าหน้าที่ 2 คนกดพร้อมกัน (AC-APV-10)

### T-35 ทดสอบการกดพร้อมกันด้วย goroutine บน PostgreSQL จริง
- รองรับ: FR-APV-01 (ช)(ซ), CON-TECH-03
- ตรวจด้วย: AC-APV-10
- ไฟล์ที่แตะ: `backend/internal/approve/service_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-16, T-23
- เสร็จเมื่อ: `TestAC_APV_10` ผ่านซ้ำ 20 รอบติด (`-count=20`): 2 goroutine อนุมัติ/ปฏิเสธคำขอเดียวกันพร้อมกัน สำเร็จ 1 อีกคนได้ `REQUEST_DECIDED_CONCURRENTLY` มี `approval_decision` 1 แถวและ `approval_email` 1 แถว
- สถานะ: พร้อมทำ

### T-36 แสดงข้อความกดพร้อมกันบนหน้าจอ
- รองรับ: FR-APV-01 (ซ)
- ตรวจด้วย: AC-APV-10
- ไฟล์ที่แตะ: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-24, T-35
- เสร็จเมื่อ: vitest `AC-APV-10` ผ่าน แสดง "คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น กรุณารีเฟรชหน้า"
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 18
- test ที่ต้องผ่านทั้งหมด: `TestAC_APV_10` (`-count=20`), vitest `AC-APV-10` และ test ของ**ทุกเฟส**
- เปิดระบบ: หลังบ้านตามหัวไฟล์ และเปิดหน้าจอ 2 แท็บด้วย `mock-staff-001`
- คนลองเองตาม AC-APV-10:
  Given เปิดหน้ารายการคำขอ 2 แท็บค้างไว้
  When แท็บแรกกดอนุมัติคำขอหนึ่ง แล้วแท็บที่สองกดอนุมัติหรือปฏิเสธคำขอเดียวกันโดยไม่รีเฟรช
  Then แท็บที่สองเห็น "คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น กรุณารีเฟรชหน้า" และ psql มี `approval_decision` ของคำขอนั้นแถวเดียว
- ยังขาด: การกดพร้อมกันจริงบนหน้าจอทำด้วยมือไม่ได้ (ลองได้แค่แบบ "หน้าจอเก่า") ส่วนพร้อมกันจริงตรวจด้วย test

---

## เฟส 19: รอคำตอบ

### T-37 เปลี่ยน RecipientResolver เป็นตัวจริงที่หา Email มหาวิทยาลัยของผู้ขอ
- รองรับ: FR-APV-02, IF-NOT-01
- ตรวจด้วย: AC-APV-04, AC-APV-13
- ไฟล์ที่แตะ: `backend/internal/approve/mailer.go` (และไฟล์ที่ได้จากคำตอบ เช่นคอลัมน์ Email ของ 002 หรือ Principal)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-14
- เสร็จเมื่อ: Email ส่งถึงที่อยู่ของผู้ขอตามแหล่งที่ทีมตัดสิน และ `TestAC_APV_04`, `TestAC_APV_13` ยังผ่าน
- สถานะ: รอ Q-APV-04

### T-38 ใส่ชื่อเจ้าหน้าที่ในประวัติและใน Email
- รองรับ: FR-APV-01 (ฌ), FR-APV-02 (ก)(ข), ASM-09
- ตรวจด้วย: AC-APV-04, AC-APV-15
- ไฟล์ที่แตะ: `backend/internal/approve/service.go`, `backend/internal/approve/email.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-28
- เสร็จเมื่อ: `actor_name` ไม่ว่าง และ Email แสดงชื่อผู้อนุมัติ / ผู้พิจารณา
- สถานะ: รอ Q-APV-04 (plan ข้อ 8 เสนอรวมชื่อเจ้าหน้าที่ใน CR เดียวกับ Email เพราะ Principal DM-08 ยังไม่มีชื่อ)

### T-39 เปลี่ยน auth จำลองเป็นระบบยืนยันตัวตนจริงของมหาวิทยาลัย
- รองรับ: IF-IDP-01, NFR-SEC-01
- ตรวจด้วย: AC-APV-05, AC-APV-12
- ไฟล์ที่แตะ: ไม่แตะไฟล์ของ 004 (ของกลาง `backend/internal/auth/` ที่ทุกฟีเจอร์ใช้) 004 ตรวจแค่ว่า test ของตัวเองยังผ่าน
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-27
- เสร็จเมื่อ: หลังเปลี่ยน auth `TestAC_APV_05` และ `TestAC_APV_12` ยังผ่าน
- สถานะ: รอ Q-SH-03

---

## ตารางตรวจความครบ

### 1. AC กับ task
| AC ID | task ที่ตรวจ AC นี้ |
|---|---|
| AC-APV-01 | T-07, T-08, T-09 |
| AC-APV-02 | T-10, T-11 |
| AC-APV-03 | T-16 |
| AC-APV-04 | T-14, T-15, T-37, T-38 |
| AC-APV-05 | T-27, T-39 |
| AC-APV-06 | T-12, T-13 |
| AC-APV-07 | T-17, T-18 |
| AC-APV-08 | T-19, T-20 |
| AC-APV-09 | T-23, T-24 |
| AC-APV-10 | T-35, T-36 |
| AC-APV-11 | T-04, T-05, T-06 |
| AC-APV-12 | T-25, T-26, T-39 |
| AC-APV-13 | T-29, T-37 |
| AC-APV-14 | T-32, T-33, T-34 |
| AC-APV-15 | T-28, T-29, T-38 |
| AC-APV-16 [เสนอ] | T-21, T-22 |
| AC-APV-17 [เสนอ] | T-31 |
| AC-APV-18 [เสนอ] | T-30 |

### 2. Constraint กับ task
| Constraint ID | task ที่ทำให้เป็นจริง |
|---|---|
| CON-TECH-01 | T-01 |
| CON-TECH-02 | T-05 (route อยู่ใน server เดิมหลัง proxy HTTPS ของกลาง) |
| CON-TECH-03 | T-02, T-05, T-12, T-17, T-19, T-23, T-25, T-27, T-35 |
| CON-DATA-01 | T-07, T-10, T-19 |
| DOM-REQ-01 | T-04, T-07, T-16, T-17 |
| DOM-LOAN-01 | T-02, T-21 |
| IF-IDP-01 | T-05, T-25, T-27, T-39 |
| IF-NOT-01 | T-14, T-15, T-32, T-37 |
| NFR-SEC-01 | T-05, T-27, T-39 |

### 3. FR กับ AC
| FR ID | AC ที่ตรวจ FR นี้ | เฟส |
|---|---|---|
| FR-APV-01 | AC-APV-01, 02, 03, 05 ถึง 12, 15, 16, 17, 18 | 1 ถึง 4, 6 ถึง 16, 18 |
| FR-APV-02 | AC-APV-04, 13, 14, 17, 18 | 5, 14 ถึง 17 |

---

## สิ่งที่ยังไม่ทำ
Open Questions จาก spec.md ที่ยังเปิดอยู่:
| คำถาม | task ที่รอ |
|---|---|
| Q-SH-03 auth service ใช้โปรโตคอลอะไร ส่งบทบาทกับภาควิชามาแบบไหน | T-39 (ระหว่างนี้ใช้ auth จำลอง ASM-02) |
| Q-SH-08 ระบบยืนยันตัวตนเชื่อมต่อไม่ได้ ต้องทำอย่างไร | ไม่มี task (ยังไม่มีคำตอบให้สร้าง) |
| Q-APV-01 รับคำขอยืมแบบค้างคืน | ไม่มี task (Out of scope ของ spec นี้ ทุก task กรอง `loan_type = 'normal'` ยกเว้นการปฏิเสธตามใน T-16 ตาม Q-APV-14) |
| Q-APV-04 Email ของผู้ขอมาจากไหน | T-37, T-38 |
| Q-APV-13 อนุมัติ ID ของ AC-APV-16 | T-21, T-22 ทำได้ ชื่อ test รอ |
| Q-APV-22 อนุมัติ ID ของ AC-APV-17, 18 | T-30, T-31 ทำได้ ชื่อ test รอ |

เรื่องที่ต้องคุยกับคนอื่น (จาก plan ข้อ 8 ไม่ใช่ task ของ 004):
- เลข migration ได้ 0008 แล้ว / ค่า action กำหนดใน migration 0008 ไม่เขียนใน data-model.md (ทีมตัดสิน 2569-11-12)
- ถ้า migration 0004 ถึง 0007 ของฟีเจอร์อื่นขึ้น main **หลัง** 0008 ฐานข้อมูลที่รันถึง 0008 ไปแล้วจะไม่รันเลขที่น้อยกว่าให้ (golang-migrate) -> ตกลงลำดับ merge ในกลุ่ม
- `auth.Middleware` ของกลางตอบ 401 ไม่มี `message` (เกี่ยวกับ T-27) -> เจ้าของ 001
- Q-APV-20 ทีมตัดสินไม่แจ้ง 003 เรื่องล็อกแถว ความเสี่ยงยกเลิกทับหลังอนุมัติยังอยู่
