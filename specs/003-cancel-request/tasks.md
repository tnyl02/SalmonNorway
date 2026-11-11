# Tasks: ยกเลิกคำขอยืม (Cancel Borrow Request)
Spec: SPEC-CXL-003 (spec.md Draft v3.2) | Plan: plan.md Draft v1.1 | สร้าง: 2569-11-10 | แก้: 2569-11-11 (ปิด Q-13) | Owner: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)

สรุป
- 18 เฟส (ฐานราก + 17 เฟสตาม AC) รวม 32 task (เดิมเฟส 17 เป็นเฟสรอคำตอบ ปิด Q-13 แล้วจึงเป็นเฟสของ AC-CXL-06)
- รอ Q-xx 0 task
- เฟสแรกที่ลองใช้งานได้จริง: เฟส 1 (นักศึกษาเห็นรายการคำขอที่รออนุมัติของตนเองบนหน้าจอ)

เทคโนโลยีตาม plan.md ข้อ 2: Go (net/http + pgx), PostgreSQL 16, React + Vite, go test + testcontainers, vitest, k6
(ทีมเลือกเอง ไม่ได้มาจาก spec ยกเว้น PostgreSQL / HTTPS / Docker) ใช้โครงกลางของ 001 ไม่สร้างโครงแอปใหม่

คำสั่งเปิดระบบ (ใช้ในทุกจุดตรวจ)
- ฐานข้อมูล: `docker compose -f deploy/docker-compose.yml up -d db`
- ข้อมูลตั้งต้นของ 003: `docker compose -f deploy/docker-compose.yml exec -T db psql -U salmon -d salmon < backend/db/seed/seed_003_cancel_request.sql`
- หลังบ้าน: `cd backend && go run ./cmd/server`
- หน้าจอ: `cd frontend && npm run dev`
- test: `cd backend && go test ./... -v` และ `cd frontend && npm test`
- ผู้ใช้จำลอง (ของกลาง `backend/internal/auth/mock_provider.go`): `mock-student-001`, `mock-student-002` (ภาควิชา 1), `mock-student-101` (ภาควิชา 2)

ต้องเคลียร์ก่อนเริ่ม /implement (ไม่ใช่ Q ใน spec จึงไม่ได้ใส่เป็นสถานะของ task)
- ~~โครงกลางยังไม่อยู่ใน main~~ merge เข้า main และเข้า branch 003 แล้ว (2569-11-11)
- เลข migration ต้องจองในกลุ่มก่อน: data-model.md จอง `0003_borrow_request` ให้ 002 และ `0004` ขึ้นไปให้ 003 (T-01, T-30)
  จองในกลุ่มแล้ว 2569-11-11: 003 ได้ `0004` และ `0005`
- ~~ชื่อคอลัมน์และค่าสถานะของ `borrow_request` ต้องตกลงกับเจ้าของ 002 ก่อนทำ T-01~~ ตกลงแล้ว 2569-11-11: 002 สร้างเองใน 0003 พร้อม cancel_reason, cancelled_at และ FK ของ audit_log (ASM-09, ASM-11)
  ~~T-01 ยังต้องรอ 002 push `0003_borrow_request` เข้า main~~ merge main เข้า branch 003 แล้ว 2569-11-11 T-01 เริ่มได้
  002 ใส่ index `idx_request_owner_status` ใน 0003 ให้แล้ว (ตกลงรอบที่ 2)
- งานหน้าจอ (T-05, 07, 08, 10, 12, 16, 18, 20, 22, 24, 26) ไม่รองานหลังบ้าน (ทีมตัดสิน 2569-11-11) ทำต่อกันตามลำดับหน้าจอได้เลย
  แต่จุดตรวจของแต่ละเฟสยังต้องรอทั้งหน้าจอและหลังบ้านของเฟสนั้นเสร็จ
- ห้าม push 0004 และ 0005 ของ 003 เข้า main ก่อน 0003 ของ 002 (ตกลงกับ 002 2569-11-11)
  ถ้าเผลอ merge ไม่ตรงลำดับ ทุกคนต้องรัน `docker compose -f deploy/docker-compose.yml down -v` (001 แจ้ง)
- 001 รับ action notify แล้ว (2569-11-11) และ writer.go (T-19 ของ 001) จะรับ pgx.Tx ได้ ยังไม่ได้วันเสร็จ (T-31)
- data-model.md: 002 มี v5 ถึง v8 ยังไม่ได้ push ตอน merge main ให้นำแถว notify ของ 003 ไปวางบน v8 เป็น v9 ไม่ใช้ v4.1
- 004 ต้องใช้เงื่อนไขสถานะตอนอนุมัติ (ASM-05) ไม่กระทบ task ของ 003 แต่กระทบผลของ AC-CXL-04 ตอนลองทั้งระบบ
  ณัฐสุดายืนยันแล้ว 2569-11-11: รับกฎไปเขียนใน spec ของ 004 และหน้าอนุมัติซ่อนคำขอ cancelled
- task ที่แตะไฟล์ของกลาง (`backend/internal/platform/http.go`, `backend/cmd/server/main.go`, `frontend/src/App.jsx`, `frontend/src/api/client.js`, `backend/internal/config/config.go`, `backend/db/migrations/`) ต้องจดใน prompt-log ว่าแก้เพราะ 003 (AGENTS.md)

---

## เฟส 0: ฐานราก

### T-01 สร้าง migration staff_contact และตรวจตาราง borrow_request ของ 002
- รองรับ: CON-TECH-01, CON-TECH-02, CON-CXL-02, CON-CXL-01, CON-DATA-01, FR-CXL-08 (ช่อง cancel_reason), ASM-09, ASM-10, ASM-11
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-04 ขึ้นไป
- ไฟล์ที่แตะ: `backend/db/migrations/0004_staff_contact.up.sql`, `backend/internal/cancelrequest/migrationtest/migration_test.go` (package แยก เพราะ package cancelrequest compile ไม่ผ่านจนกว่า T-04 ทีมตัดสิน 2569-11-11) (ไม่สร้าง borrow_request เพราะ 002 สร้างใน 0003 แก้ 2569-11-11), `backend/internal/cancelrequest/tc_helpers_test.go` (ตัวช่วย ไม่ใช่ TC: tcInsert ต้องใส่ requester_name, student_code, start_at, end_at และ purpose เป็น study / activity / other ตาม 0003 จริงของ 002 เพิ่ม 2569-11-11 หลัง merge main)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี task ของ 003 แต่ต้องมี `0003_borrow_request` ของ 002 ใน branch ก่อน (ดูส่วนหัว)
- เสร็จเมื่อ: `TestMigration_BorrowRequestStatusCheck` (สถานะนอก 8 ค่าถูกปฏิเสธ) และ `TestMigration_CancelReasonMaxLen` (201 ตัวอักษรถูกปฏิเสธ) ผ่านกับตารางของ 002, ตาราง staff_contact มีอยู่ และไม่มี trigger หรือ FK ที่แก้ `equipment_item.status`
- สถานะ: เสร็จ

### T-02 เพิ่มข้อมูลตั้งต้นคำขอทุกสถานะและ Email เจ้าหน้าที่จำลอง
- รองรับ: ASM-01, ASM-04 (คำขอค้างคืนจาก seed), ASM-10, CON-CXL-02
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของจุดตรวจทุกเฟส
- ไฟล์ที่แตะ: `backend/db/seed/seed_003_cancel_request.sql` (ไฟล์แยก ไม่แก้ seed.sql ของ 001), `backend/internal/cancelrequest/migrationtest/seed_test.go` (เพิ่ม 2569-11-11 ทีมอนุญาต: test ยืนยันช่องเสร็จเมื่อ), `backend/internal/cancelrequest/migrationtest/migration_test.go` (เพิ่ม dbURL และเปลี่ยนผู้ขอใน test เป็น student-test ไม่ชนกับ seed)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: โหลด seed แล้ว student-001 มีคำขอ `requested` 2 รายการ (คนละเวลา หนึ่งในนั้นเป็นอุปกรณ์ CS-NB-0001), `approved_by_officer` 1, `pending_advisor_approval` 1, `approved` 1, `rejected` 1, `cancelled` 1, student-002 มีคำขอ `requested` บน CS-NB-0001 1 รายการ, student-101 ไม่มีคำขอ และภาควิชา 1 มี Email เจ้าหน้าที่จำลอง 1 รายการ
- สถานะ: เสร็จ

### T-03 เพิ่มค่าที่ตั้งได้ของ 003 ใน config
- รองรับ: CON-CXL-02, FR-CXL-02, FR-CXL-03, FR-CXL-06, FR-CXL-08, NFR-PERF-CXL-01, ASM-06
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-04 ขึ้นไป
- ไฟล์ที่แตะ: `backend/internal/config/config.go` (`CancellableStatuses`, `CancelReasonMaxLen`, `CancelPerfP95`, `CancelDBTimeout`, `DisplayTimezone` ถ้า 001 ยังไม่เพิ่ม), `frontend/src/config.js` (`CANCEL_REASON_MAX_LEN`, `DISPLAY_TIMEZONE`, `MSG_NO_PENDING`, `MSG_NOT_CANCELLABLE`, `MSG_REQUEST_NOT_FOUND`, `MSG_ALREADY_APPROVED`, `MSG_CANCEL_SUCCESS`, `MSG_LOAD_FAILED`, `LABEL_RETRY` ตาม ASM-13), `backend/internal/config/config_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: `TestConfig_CancellableStatuses` ผ่าน (มี 3 ค่า requested, approved_by_officer, pending_advisor_approval) และ `CancelDBTimeout` น้อยกว่า 15 วินาที (ตัวเลขทีมเลือกเอง จดใน prompt-log)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 0
- test ที่ต้องผ่านทั้งหมด: `TestMigration_BorrowRequestStatusCheck`, `TestMigration_CancelReasonMaxLen`, `TestConfig_CancellableStatuses` และ test ของ 001 ต้องยังผ่าน
- เปิดระบบ: คำสั่งฐานข้อมูลและ seed ของ 003 ในส่วนหัว
- คนลองเอง (ไม่มี AC):
  Given เปิดฐานข้อมูลและโหลด seed ของ 003
  When รัน `SELECT requester_id, status, count(*) FROM borrow_request GROUP BY 1, 2` ใน psql
  Then เห็นคำขอของ student-001 และ student-002 ตามที่ T-02 กำหนด และ `equipment_item` ยังมีสถานะเหมือนก่อนโหลด seed
- ยังขาด: ไม่มี

---

## เฟส 1: เห็นรายการคำขอที่รออนุมัติของตนเอง (AC-CXL-08)

### T-04 สร้าง API รายการคำขอที่รออนุมัติของฉัน
- รองรับ: FR-CXL-06, FR-CXL-03, CON-CXL-02, IF-IDP-01
- ตรวจด้วย: AC-CXL-08
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/repository.go`, `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/handler.go`, `backend/internal/cancelrequest/dto.go`, `backend/internal/cancelrequest/handler_test.go`, `backend/cmd/server/main.go` (ของกลาง ลงทะเบียน route หลัง `auth.Middleware` แบบเดียวกับ 001 เดิมเขียน platform/http.go ทีมแก้ 2569-11-11), `backend/internal/cancelrequest/tc_helpers_test.go` และ `tc_ac_cxl_08_test.go` (แก้ Given ของ TC-CXL-08-1, 08-3 ตาม test-cases.md ที่ทีมแก้ 2569-11-11)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01, T-03
- เสร็จเมื่อ: `TestAC_CXL_08_ListOwnPendingSorted` ผ่าน (`GET /api/v1/my/pending-requests` คืนเฉพาะคำขอของผู้เรียกที่ยกเลิกได้ เรียงใหม่ไปเก่า เวลาเท่ากันเรียง requestId มากไปน้อย มี requestId, modelName, assetCode, submittedAt เจ้าหน้าที่ได้ 403 `STUDENT_ONLY`)
- สถานะ: เสร็จ

### T-05 สร้างหน้า "คำขอที่รออนุมัติของฉัน" แสดงรายการ
- รองรับ: FR-CXL-06
- ตรวจด้วย: AC-CXL-08
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/MyPendingRequests.test.jsx`, `frontend/src/api/client.js` (ของกลาง เพิ่ม `listMyPendingRequests`), `frontend/src/App.jsx` (ของกลาง เพิ่มทางเข้าหน้า)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-03 (เดิม T-04 ทีมแก้ 2569-11-11: สัญญา API อยู่ใน plan ข้อ 4 แล้ว และ test หน้าจอใช้ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจเฟส 1 หลัง T-04 เสร็จ)
- เสร็จเมื่อ: vitest `AC_CXL_08 shows own pending requests newest first` ผ่าน (แต่ละแถวมีชื่อรุ่น รหัสครุภัณฑ์ เวลาที่ส่งแบบ Asia/Bangkok และปุ่มยกเลิก)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 1
- test ที่ต้องผ่านทั้งหมด: `TestMigration_BorrowRequestStatusCheck`, `TestMigration_CancelReasonMaxLen`, `TestConfig_CancellableStatuses`, `TestAC_CXL_08_ListOwnPendingSorted`, vitest `AC_CXL_08 ...` และ test ของ 001 ต้องยังผ่าน
- เปิดระบบ: ตามคำสั่งในส่วนหัว (ฐานข้อมูล, seed ของ 003, หลังบ้าน, หน้าจอ)
- คนลองเองตาม AC-CXL-08:
  Given ใช้ข้อมูลตั้งต้นของ T-02 และเข้าระบบเป็น student-001
  When เปิดหน้า "คำขอที่รออนุมัติของฉัน"
  Then เห็น 4 รายการ (requested 2, approved_by_officer 1, pending_advisor_approval 1) เรียงจากส่งล่าสุดอยู่บนสุด แต่ละแถวมีชื่อรุ่น รหัสครุภัณฑ์ เวลาที่ส่ง และปุ่มยกเลิก ไม่เห็นคำขอที่ approved, rejected, cancelled และไม่เห็นคำขอของ student-002
- ยังขาด: ไม่มี

---

## เฟส 2: ยกเลิกคำขอสำเร็จ (AC-CXL-01)

### T-06 สร้าง API ยกเลิกคำขอ (กรณีสำเร็จ)
- รองรับ: FR-CXL-01, FR-CXL-04, FR-CXL-08 (เก็บ cancel_reason), CON-CXL-01, CON-DATA-01, IF-IDP-01
- ตรวจด้วย: AC-CXL-01
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/repository.go`, `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/handler.go`, `backend/internal/cancelrequest/dto.go`, `backend/internal/cancelrequest/handler_test.go` (route ลงทะเบียนใน main.go แล้วตั้งแต่ T-04 ไม่ต้องแตะไฟล์กลาง)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_CXL_01_CancelPendingNormal` ผ่าน (`POST /api/v1/borrow-requests/{id}/cancel` ด้วย UPDATE แบบมีเงื่อนไขคำสั่งเดียว ตอบ 200 สถานะเป็น `cancelled` เก็บ cancel_reason และ cancelled_at และ `equipment_item.status` ยังเป็น `available`)
- สถานะ: เสร็จ

### T-07 สร้างกล่องยืนยันการยกเลิกและแสดงผลสำเร็จ
- รองรับ: FR-CXL-01, FR-CXL-07, FR-CXL-08
- ตรวจด้วย: AC-CXL-01
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/CancelDialog.jsx`, `frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/CancelDialog.test.jsx`, `frontend/src/api/client.js` (ของกลาง เพิ่ม `cancelRequest`)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-05 (เดิม T-05, T-06 ทีมแก้ 2569-11-11: งานหน้าจอไม่รองานหลังบ้าน ใช้สัญญา API ใน plan ข้อ 4 และ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจของเฟส)
- เสร็จเมื่อ: vitest `AC_CXL_01 cancels with reason and shows success` ผ่าน (พิมพ์เหตุผล กดยืนยัน เห็น "ยกเลิกคำขอยืม {ชื่อรุ่น} ({รหัสครุภัณฑ์}) เรียบร้อยแล้ว" และแถวหายจากรายการ)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 2
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_01_CancelPendingNormal`, vitest `AC_CXL_01 ...` และ test ของเฟส 0-1
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-01:
  Given เข้าระบบเป็น student-001 และเปิดหน้า "คำขอที่รออนุมัติของฉัน" (โหลด seed ใหม่ก่อนถ้าเคยลองแล้ว)
  When กดยกเลิกคำขอของ CS-NB-0001 พิมพ์เหตุผล "ไม่ต้องใช้อุปกรณ์แล้ว" แล้วกดยืนยัน
  Then เห็นข้อความยกเลิกสำเร็จ แถวนั้นหายจากหน้า และเมื่อเปิดหน้ารายการอุปกรณ์ของ 001 CS-NB-0001 ยังเป็น "ว่าง"
- ยังขาด: ยังไม่มีบันทึก audit_log (เฟส 17) และยังไม่ส่ง Email (เฟส 14)

---

## เฟส 3: ยืนยันหรือไม่ยืนยันในกล่อง (AC-CXL-09)

### T-08 ทำปุ่มไม่ยืนยันให้ปิดกล่องโดยไม่ยกเลิก
- รองรับ: FR-CXL-07
- ตรวจด้วย: AC-CXL-09
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/CancelDialog.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/CancelDialog.test.jsx`, `frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx` (เพิ่ม 2569-11-11 ทีมอนุญาต: หน้าเป็นผู้ปิดกล่อง)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-07
- เสร็จเมื่อ: vitest `AC_CXL_09 confirm dialog cancel then confirm` ผ่าน (ไม่ยืนยันแล้ว client จำลองไม่ถูกเรียก แถวยังอยู่ เปิดกล่องใหม่ช่องเหตุผลว่าง ยืนยันรอบสองแล้วแถวหาย)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 3
- test ที่ต้องผ่านทั้งหมด: vitest `AC_CXL_09 ...` และ test ของเฟส 0-2
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-09:
  Given เข้าระบบเป็น student-001 มีคำขอที่รออนุมัติอย่างน้อย 1 รายการ
  When กดยกเลิก แล้วกดไม่ยืนยัน จากนั้นกดยกเลิกอีกครั้ง พิมพ์เหตุผล แล้วกดยืนยัน
  Then ครั้งแรกแถวยังอยู่ (รีเฟรชหน้าแล้วยังอยู่) ครั้งที่สองแถวหายจากหน้า
- ยังขาด: ไม่มี

---

## เฟส 4: ต้องพิมพ์เหตุผลไม่เกิน 200 ตัวอักษร (AC-CXL-14)

### T-09 ตรวจเหตุผลที่หลังบ้าน
- รองรับ: FR-CXL-08
- ตรวจด้วย: AC-CXL-14
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/handler.go`, `backend/internal/cancelrequest/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestAC_CXL_14_ReasonValidation` ผ่าน (ว่าง และช่องว่างล้วน -> 400 `REASON_REQUIRED`, 201 ตัวอักษร -> 400 `REASON_TOO_LONG` สถานะไม่เปลี่ยน, "ไม่ต้องใช้อุปกรณ์แล้ว" -> 200 และ cancel_reason ตรงกัน)
- สถานะ: เสร็จ

### T-10 ปิดปุ่มยืนยันเมื่อเหตุผลว่างหรือเกิน และแสดงจำนวนตัวอักษร
- รองรับ: FR-CXL-08, FR-CXL-07
- ตรวจด้วย: AC-CXL-14
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/CancelDialog.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/CancelDialog.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-08 (เดิม T-08, T-09 ทีมแก้ 2569-11-11: งานหน้าจอไม่รองานหลังบ้าน ใช้สัญญา API ใน plan ข้อ 4 และ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจของเฟส)
- เสร็จเมื่อ: vitest `AC_CXL_14 confirm disabled when reason empty or too long` ผ่าน (ตัวนับ "x/200" นับตาม code point เท่ากับหลังบ้าน)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 4
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_14_ReasonValidation`, vitest `AC_CXL_14 ...` และ test ของเฟส 0-3
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-14:
  Given เข้าระบบเป็น student-001 และกดยกเลิกคำขอหนึ่ง
  When ไม่พิมพ์เหตุผล (หรือพิมพ์แต่ช่องว่าง), วางข้อความยาว 201 ตัวอักษร, แล้วแก้เป็น "ไม่ต้องใช้อุปกรณ์แล้ว" แล้วกดยืนยัน
  Then สองกรณีแรกกดยืนยันไม่ได้ กรณีที่สามยกเลิกสำเร็จ แถวหายจากหน้า
- ยังขาด: ไม่มี

---

## เฟส 5: ไม่มีคำขอที่รออนุมัติ (AC-CXL-12)

### T-11 ทดสอบ API กรณีไม่มีคำขอ
- รองรับ: FR-CXL-06
- ตรวจด้วย: AC-CXL-12
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/handler_test.go` (แก้ `handler.go` เฉพาะถ้าคืน `null` แทน `[]`)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_CXL_12_EmptyList` ผ่าน (คืน `items: []`)
- สถานะ: เสร็จ

### T-12 แสดงข้อความเมื่อไม่มีคำขอ
- รองรับ: FR-CXL-06
- ตรวจด้วย: AC-CXL-12
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/MyPendingRequests.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-05 (เดิม T-05, T-11 ทีมแก้ 2569-11-11: งานหน้าจอไม่รองานหลังบ้าน ใช้สัญญา API ใน plan ข้อ 4 และ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจของเฟส)
- เสร็จเมื่อ: vitest `AC_CXL_12 shows empty message` ผ่าน (แสดง `MSG_NO_PENDING` "ไม่มีคำขอที่รออนุมัติ" และไม่มีแถว ส่วนโหลดไม่สำเร็จแสดง `MSG_LOAD_FAILED` พร้อมปุ่ม "ลองใหม่" แทน)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 5
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_12_EmptyList`, vitest `AC_CXL_12 ...` และ test ของเฟส 0-4
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-12:
  Given เข้าระบบเป็น student-101 (ไม่มีคำขอใน seed)
  When เปิดหน้า "คำขอที่รออนุมัติของฉัน"
  Then เห็นข้อความ "ไม่มีคำขอที่รออนุมัติ" และไม่มีรายการใด
- ยังขาด: ไม่มี

---

## เฟส 6: คำขออื่นของอุปกรณ์ชิ้นเดียวกันไม่เปลี่ยน (AC-CXL-05)

### T-13 ทดสอบว่ายกเลิกแล้วไม่กระทบคำขออื่นและสถานะอุปกรณ์
- รองรับ: DOM-REQ-01, CON-CXL-01, CON-DATA-01, FR-CXL-01
- ตรวจด้วย: AC-CXL-05
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestAC_CXL_05_OtherRequestsOnSameItemUnchanged` ผ่าน
- สถานะ: เสร็จ

#### จุดตรวจเฟส 6
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_05_OtherRequestsOnSameItemUnchanged` และ test ของเฟส 0-5
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-05:
  Given โหลด seed ใหม่ (student-001 และ student-002 ต่างมีคำขอ requested บน CS-NB-0001)
  When เข้าระบบเป็น student-001 แล้วยกเลิกคำขอของ CS-NB-0001 จากนั้นเข้าระบบเป็น student-002 แล้วเปิดหน้า "คำขอที่รออนุมัติของฉัน"
  Then student-002 ยังเห็นคำขอของ CS-NB-0001 อยู่ และหน้ารายการอุปกรณ์ของ 001 CS-NB-0001 ยังเป็น "ว่าง"
- ยังขาด: ไม่มี

---

## เฟส 7: ยกเลิกคำขอค้างคืนที่ยังไม่มีผลอนุมัติของอาจารย์ (AC-CXL-10)

### T-14 ทดสอบการยกเลิกคำขอค้างคืน 2 สถานะ
- รองรับ: FR-CXL-01, CON-CXL-02, ASM-04
- ตรวจด้วย: AC-CXL-10
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestAC_CXL_10_CancelOvernightApprovedByOfficer` ผ่าน (`approved_by_officer` และ `pending_advisor_approval` ยกเลิกได้ทั้งคู่)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 7
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_10_CancelOvernightApprovedByOfficer` และ test ของเฟส 0-6
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-10:
  Given เข้าระบบเป็น student-001 (seed มีคำขอค้างคืน approved_by_officer และ pending_advisor_approval)
  When ยกเลิกคำขอค้างคืนทั้ง 2 รายการพร้อมเหตุผล
  Then ยกเลิกสำเร็จทั้งคู่ และทั้ง 2 แถวหายจากหน้า
- ยังขาด: ไม่มี

---

## เฟส 8: คำขอที่อนุมัติแล้วยกเลิกไม่ได้ (AC-CXL-02)

### T-15 แยกสาเหตุเมื่อยกเลิกไม่สำเร็จ กรณีอนุมัติแล้ว
- รองรับ: FR-CXL-02, FR-CXL-04, CON-CXL-02
- ตรวจด้วย: AC-CXL-02
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/repository.go`, `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/handler.go`, `backend/internal/cancelrequest/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestAC_CXL_02_ApprovedNotCancellable` ผ่าน (UPDATE ไม่เจอแถว แล้ว SELECT แยกสาเหตุ คำขอ `approved` -> 409 `ALREADY_APPROVED` สถานะยัง `approved`)
- สถานะ: เสร็จ

### T-16 แสดงข้อความเมื่อคำขอได้รับการอนุมัติแล้ว
- รองรับ: FR-CXL-02
- ตรวจด้วย: AC-CXL-02
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/MyPendingRequests.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-07 (เดิม T-07, T-15 ทีมแก้ 2569-11-11: งานหน้าจอไม่รองานหลังบ้าน ใช้สัญญา API ใน plan ข้อ 4 และ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจของเฟส)
- เสร็จเมื่อ: vitest `AC_CXL_02 shows already approved message and reloads` ผ่าน (client จำลองตอบ 409 `ALREADY_APPROVED` -> แสดง `MSG_ALREADY_APPROVED` แล้วโหลดรายการใหม่)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 8
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_02_ApprovedNotCancellable`, vitest `AC_CXL_02 ...` และ test ของเฟส 0-7
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-02:
  Given เข้าระบบเป็น student-001 ใช้ id ของคำขอ approved ใน seed
  When ส่งคำสั่งยกเลิกโดยตรง (เช่น `curl -X POST -H "Authorization: Bearer mock-student-001" -d '{"reason":"ทดสอบ"}' http://localhost:8000/api/v1/borrow-requests/<id>/cancel`)
  Then ได้ 409 `ALREADY_APPROVED` และคำขอยังเป็น approved
- ยังขาด: ไม่มี

---

## เฟส 9: คำขอที่ถูกปฏิเสธหรือยกเลิกแล้วยกเลิกไม่ได้ (AC-CXL-11)

### T-17 แยกสาเหตุกรณีสถานะอื่นที่ยกเลิกไม่ได้
- รองรับ: FR-CXL-02, CON-CXL-02
- ตรวจด้วย: AC-CXL-11
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-15
- เสร็จเมื่อ: `TestAC_CXL_11_RejectedAndCancelledNotCancellable` ผ่าน (`rejected` และ `cancelled` -> 409 `NOT_CANCELLABLE` สถานะคงเดิม)
- สถานะ: เสร็จ

### T-18 แสดงข้อความ "คำขอนี้ยกเลิกไม่ได้"
- รองรับ: FR-CXL-02
- ตรวจด้วย: AC-CXL-11
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/MyPendingRequests.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-16 (เดิม T-16, T-17 ทีมแก้ 2569-11-11: งานหน้าจอไม่รองานหลังบ้าน ใช้สัญญา API ใน plan ข้อ 4 และ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจของเฟส)
- เสร็จเมื่อ: vitest `AC_CXL_11 shows not cancellable message` ผ่าน (409 `NOT_CANCELLABLE` -> แสดง `MSG_NOT_CANCELLABLE` แล้วโหลดรายการใหม่)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 9
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_11_RejectedAndCancelledNotCancellable`, vitest `AC_CXL_11 ...` และ test ของเฟส 0-8
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-11:
  Given เข้าระบบเป็น student-001 ใช้ id ของคำขอ rejected และ cancelled ใน seed
  When ส่งคำสั่งยกเลิกแต่ละคำขอโดยตรงด้วย curl แบบเฟส 8
  Then ได้ 409 `NOT_CANCELLABLE` ทั้ง 2 ครั้ง และสถานะคงเดิม
- ยังขาด: ไม่มี

---

## เฟส 10: ยกเลิกคำขอของผู้อื่นไม่ได้ (AC-CXL-03)

### T-19 ตอบไม่พบคำขอเมื่อไม่ใช่คำขอของตนเอง
- รองรับ: FR-CXL-03, IF-IDP-01
- ตรวจด้วย: AC-CXL-03
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/handler_test.go`, `backend/internal/cancelrequest/handler.go` (เพิ่ม 2569-11-11 ทีมอนุญาต: แปลง error เป็น 404)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-15
- เสร็จเมื่อ: `TestAC_CXL_03_OthersRequestNotFound` ผ่าน (คำขอของผู้อื่นและ id ที่ไม่มี -> 404 `REQUEST_NOT_FOUND` สถานะไม่เปลี่ยน body ไม่มีข้อมูลคำขอ)
- สถานะ: เสร็จ

### T-20 แสดงข้อความ "ไม่พบคำขอ"
- รองรับ: FR-CXL-03
- ตรวจด้วย: AC-CXL-03
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/MyPendingRequests.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-18 (เดิม T-18, T-19 ทีมแก้ 2569-11-11: งานหน้าจอไม่รองานหลังบ้าน ใช้สัญญา API ใน plan ข้อ 4 และ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจของเฟส)
- เสร็จเมื่อ: vitest `AC_CXL_03 shows request not found message` ผ่าน (404 -> แสดง `MSG_REQUEST_NOT_FOUND` แล้วโหลดรายการใหม่)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 10
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_03_OthersRequestNotFound`, vitest `AC_CXL_03 ...` และ test ของเฟส 0-9
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-03:
  Given ใช้ id ของคำขอ requested ของ student-002 ใน seed
  When ส่งคำสั่งยกเลิกด้วย token `mock-student-001`
  Then ได้ 404 `REQUEST_NOT_FOUND` และเมื่อเข้าระบบเป็น student-002 คำขอนั้นยังอยู่ในหน้า
- ยังขาด: ไม่มี

---

## เฟส 11: เจ้าหน้าที่อนุมัติก่อนนักศึกษากดยืนยัน (AC-CXL-04)

### T-21 ทดสอบกรณีสถานะเปลี่ยนระหว่างเปิดหน้ากับกดยืนยัน
- รองรับ: FR-CXL-04, FR-CXL-02
- ตรวจด้วย: AC-CXL-04
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-15
- เสร็จเมื่อ: `TestAC_CXL_04_ApprovedBeforeConfirm` ผ่าน (GET เห็นคำขอ -> UPDATE ในฐานข้อมูลเป็น approved แทน 004 -> POST cancel -> 409 `ALREADY_APPROVED` สถานะยัง approved)
- สถานะ: เสร็จ

### T-22 ทดสอบหน้าจอเมื่อกดยืนยันแล้วได้ผลว่าอนุมัติแล้ว
- รองรับ: FR-CXL-04, FR-CXL-07
- ตรวจด้วย: AC-CXL-04
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/__tests__/CancelDialog.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-16 (เดิม T-16, T-21 ทีมแก้ 2569-11-11: งานหน้าจอไม่รองานหลังบ้าน ใช้สัญญา API ใน plan ข้อ 4 และ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจของเฟส)
- เสร็จเมื่อ: vitest `AC_CXL_04 approved while dialog open` ผ่าน (เปิดกล่องแล้ว client จำลองเปลี่ยนเป็นตอบ 409 -> กดยืนยัน -> แสดง `MSG_ALREADY_APPROVED` แถวนั้นหายหลังโหลดใหม่)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 11
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_04_ApprovedBeforeConfirm`, vitest `AC_CXL_04 ...` และ test ของเฟส 0-10
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-04:
  Given เข้าระบบเป็น student-001 เปิดหน้า "คำขอที่รออนุมัติของฉัน" กดยกเลิกคำขอหนึ่ง และพิมพ์เหตุผลไว้ (ยังไม่กดยืนยัน)
  When อีกหน้าต่างหนึ่งรัน `UPDATE borrow_request SET status='approved' WHERE id=<id>` ใน psql (แทน 004 ที่ยังไม่มี) แล้วกลับมากดยืนยัน
  Then เห็นข้อความว่าคำขอได้รับการอนุมัติแล้ว และคำขอยังเป็น approved
- ยังขาด: ลองกับหน้าอนุมัติของ 004 จริงเมื่อ 004 เสร็จ (ASM-05)

---

## เฟส 12: ฐานข้อมูลไม่ตอบสนองตอนยกเลิก (AC-CXL-13)

### T-23 ตอบ 503 เมื่อฐานข้อมูลล่มหรือเกินเวลา
- รองรับ: FR-CXL-01, ASM-06
- ตรวจด้วย: AC-CXL-13
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/handler.go`, `backend/internal/cancelrequest/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06, T-03
- เสร็จเมื่อ: `TestAC_CXL_13_DBDownReturns503` ผ่าน (หยุด container ก่อน POST -> 503 `UPSTREAM_UNAVAILABLE` retryable, ใช้ `CancelDBTimeout`, เปิดฐานข้อมูลกลับมาแล้วสถานะยัง requested)
- สถานะ: เสร็จ

### T-24 แสดงข้อผิดพลาดพร้อมปุ่มลองใหม่
- รองรับ: FR-CXL-01, ASM-06
- ตรวจด้วย: AC-CXL-13
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/CancelDialog.jsx`, `frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/CancelDialog.test.jsx`, `frontend/src/config.js` (เพิ่ม 2569-11-12: MSG_CANCEL_FAILED ทีมเลือกคำ)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-07 (เดิม T-07, T-23 ทีมแก้ 2569-11-11: งานหน้าจอไม่รองานหลังบ้าน ใช้สัญญา API ใน plan ข้อ 4 และ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจของเฟส)
- เสร็จเมื่อ: vitest `AC_CXL_13 shows error with retry` ผ่าน (503 -> ข้อความผิดพลาด ปุ่มลองใหม่ แถวยังอยู่ กดลองใหม่แล้วเรียก API ซ้ำ)
- สถานะ: เสร็จ

#### จุดตรวจเฟส 12
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_13_DBDownReturns503`, vitest `AC_CXL_13 ...` และ test ของเฟส 0-11
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-13:
  Given เข้าระบบเป็น student-001 เปิดหน้า "คำขอที่รออนุมัติของฉัน" แล้วหยุดฐานข้อมูล `docker compose -f deploy/docker-compose.yml stop db`
  When กดยกเลิกคำขอ พิมพ์เหตุผล แล้วกดยืนยัน
  Then เห็นข้อความผิดพลาดพร้อมปุ่มลองใหม่ แถวยังอยู่ เปิดฐานข้อมูลกลับ (`start db`) แล้วคำขอยังเป็น requested
- ยังขาด: ไม่มี

---

## เฟส 13: session หมดอายุตอนกดยกเลิก (AC-CXL-07)

### T-25 ทดสอบว่า 401 ไม่ยกเลิกคำขอ
- รองรับ: IF-IDP-01, NFR-SEC-01
- ตรวจด้วย: AC-CXL-07
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestAC_CXL_07_ExpiredSessionNoCancel` ผ่าน (ไม่มี token หรือ token ไม่รู้จัก -> 401 `UNAUTHENTICATED` สถานะไม่เปลี่ยน) ส่วนตรวจ audit_log อยู่ที่ T-31
- สถานะ: เสร็จ

### T-26 พาไปหน้าเข้าสู่ระบบและกลับมาหน้าเดิม
- รองรับ: IF-IDP-01
- ตรวจด้วย: AC-CXL-07
- ไฟล์ที่แตะ: `frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx`, `frontend/src/pages/MyPendingRequests/__tests__/MyPendingRequests.test.jsx` (ถ้า 001 ทำกลไก redirect พร้อม `returnTo` ไว้แล้วใน AC-EQS-09 ให้ใช้ของเดิม ไม่เขียนใหม่)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (plan ข้อ 4)
- ต้องทำหลัง: T-07 (เดิม T-07, T-25 ทีมแก้ 2569-11-11: งานหน้าจอไม่รองานหลังบ้าน ใช้สัญญา API ใน plan ข้อ 4 และ client จำลอง ลองกับหลังบ้านจริงที่จุดตรวจของเฟส)
- เสร็จเมื่อ: vitest `AC_CXL_07 redirects to login with returnTo` ผ่าน (401 -> ไปหน้าเข้าสู่ระบบพร้อม `returnTo` ของหน้านี้ ไม่เรียกยกเลิกซ้ำอัตโนมัติ `returnTo` ภายนอกระบบถูกปฏิเสธ)
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 13
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_07_ExpiredSessionNoCancel`, vitest `AC_CXL_07 ...` และ test ของเฟส 0-12
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-07:
  Given เข้าระบบเป็น student-001 เปิดหน้า "คำขอที่รออนุมัติของฉัน" แล้วลบ token ออก (จำลอง session หมดอายุ)
  When กดยกเลิกคำขอ พิมพ์เหตุผล แล้วกดยืนยัน
  Then ถูกพาไปหน้าเข้าสู่ระบบพร้อมข้อความแจ้ง เข้าสู่ระบบใหม่แล้วกลับมาหน้า "คำขอที่รออนุมัติของฉัน" และคำขอยังอยู่ (ต้องกดยกเลิกใหม่)
- ยังขาด: ตรวจว่าไม่มีบันทึกใน audit_log (T-31 เฟส 17)

---

## เฟส 14: แจ้งเจ้าหน้าที่ทาง Email (AC-CXL-15)

### T-27 สร้าง Notifier จำลองและส่งหลังยกเลิกสำเร็จ
- รองรับ: FR-CXL-09, IF-NOT-01, ASM-10
- ตรวจด้วย: AC-CXL-15
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/notifier.go`, `backend/internal/cancelrequest/notifier_test.go`, `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/repository.go` (อ่าน staff_contact), `backend/cmd/server/main.go` (ของกลาง ส่ง Notifier จำลองเข้า handler)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-02, T-06
- เสร็จเมื่อ: `TestAC_CXL_15_NotifyDepartmentStaff` ผ่าน (Notifier ปลอมได้ 1 ข้อความถึง Email ของภาควิชาเจ้าของอุปกรณ์ มีผู้ยกเลิก ชื่อรุ่น รหัสครุภัณฑ์ เวลา และเหตุผล ส่งหลัง commit)
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 14
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_15_NotifyDepartmentStaff` และ test ของเฟส 0-13
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-15:
  Given เข้าระบบเป็น student-001 คำขอเป็นอุปกรณ์ของภาควิชา 1
  When ยกเลิกคำขอพร้อมเหตุผล "ไม่ต้องใช้อุปกรณ์แล้ว"
  Then log ของหลังบ้านมีข้อความ Email จำลองถึง Email เจ้าหน้าที่ภาควิชา 1 ระบุผู้ยกเลิก ชื่อรุ่น รหัสครุภัณฑ์ เวลา และเหตุผลนั้น
- ยังขาด: Email จริง (Q-15 ส่วนระบบจริง ไม่ทำใน prototype ตาม ASM-10)

---

## เฟส 15: ส่ง Email ไม่สำเร็จแต่ยกเลิกสำเร็จ (AC-CXL-16)

### T-28 ไม่ย้อนการยกเลิกเมื่อ Notifier ล้มเหลว
- รองรับ: FR-CXL-09
- ตรวจด้วย: AC-CXL-16
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/notifier_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-27
- เสร็จเมื่อ: `TestAC_CXL_16_NotifyFailureStillCancelled` ผ่านส่วนที่ไม่เกี่ยวกับ audit (Notifier ปลอมคืน error -> 200 สถานะ `cancelled`) ส่วนบันทึกใน audit_log อยู่ที่ T-32
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 15
- test ที่ต้องผ่านทั้งหมด: `TestAC_CXL_16_NotifyFailureStillCancelled` (ส่วนที่ไม่ใช่ audit) และ test ของเฟส 0-14
- เปิดระบบ: ตามคำสั่งในส่วนหัว โดยตั้งให้ Notifier จำลองคืน error (วิธีตั้งเลือกตอน implement)
- คนลองเองตาม AC-CXL-16:
  Given เข้าระบบเป็น student-001 และ Notifier จำลองถูกตั้งให้ล้มเหลว
  When ยกเลิกคำขอพร้อมเหตุผลแล้วยืนยัน
  Then เห็นข้อความยกเลิกสำเร็จ แถวหายจากหน้า และคำขอเป็น cancelled
- ยังขาด: บันทึกว่าส่งแจ้งเตือนไม่สำเร็จใน audit_log (T-32 เฟส 17)

---

## เฟส 16: ผู้ใช้พร้อมกัน 50 คน (AC-CXL-17)

### T-29 เขียน load test ด้วย k6
- รองรับ: NFR-PERF-CXL-01
- ตรวจด้วย: AC-CXL-17
- ไฟล์ที่แตะ: `tests/load/cancel_request.js`, `backend/db/seed/seed_003_load.sql` (คำขอจำนวนพอสำหรับ 50 ผู้ใช้)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06, T-27
- เสร็จเมื่อ: `k6 run tests/load/cancel_request.js` ผ่าน threshold `http_req_duration p(95) < 15000` (ค่ามาจาก `CancelPerfP95`) ทั้ง GET รายการและ POST ยกเลิก ที่ 50 VU
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 16
- test ที่ต้องผ่านทั้งหมด: k6 `AC_CXL_17` และ test ของเฟส 0-15
- เปิดระบบ: ตามคำสั่งในส่วนหัว แล้วโหลด `seed_003_load.sql`
- คนลองเองตาม AC-CXL-17:
  Given ผู้ใช้จำลอง 50 คน แต่ละคนมีคำขอที่รออนุมัติ
  When รัน `k6 run tests/load/cancel_request.js`
  Then รายงานของ k6 แสดง p95 ของทั้ง 2 อย่างไม่เกิน 15 วินาที
- ยังขาด: auth จำลองมีผู้ใช้แค่ 4 คน ต้องตัดสินว่าจะเพิ่มผู้ใช้จำลองอย่างไร (แก้ของกลาง) ดูรายงานในแชต

---

## เฟส 17: บันทึกประวัติการยกเลิก (AC-CXL-06)

### T-30 เพิ่ม action notify ใน audit_log
- รองรับ: FR-CXL-05, FR-CXL-09, ASM-12, CON-TECH-01
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-31, T-32
- ไฟล์ที่แตะ: `backend/db/migrations/0005_audit_log_notify_action.up.sql` (ของกลาง ไฟล์ใหม่ เลขจองในกลุ่ม ห้ามแก้ 0002: DROP แล้ว ADD `audit_log_action_check` ชื่อเดิม ค่า list, search, detail, cancel, request, notify), `backend/internal/cancelrequest/migrationtest/migration_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: `TestMigration_AuditLogAcceptsNotify` ผ่าน (insert action `notify` ได้ ค่า action อื่นนอกตาราง DM-07 ยังถูกปฏิเสธ และ action เดิมของ 001 กับ `request` ของ 002 ยังใช้ได้)
- สถานะ: เสร็จ รอทีมตรวจ

### T-31 บันทึก audit_log ทุกครั้งที่พยายามยกเลิก
- รองรับ: FR-CXL-05, ASM-07
- ตรวจด้วย: AC-CXL-06, AC-CXL-07 (ส่วน audit)
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/repository.go`, `backend/internal/cancelrequest/handler_test.go` (ถ้า `backend/internal/audit/writer.go` ของ 001 รับ transaction ได้ ให้เรียกใช้แทนเขียนเอง), `backend/internal/cancelrequest/tc_helpers_test.go` (ตัวช่วย ไม่ใช่ TC: แก้วิธีล้างข้อมูล เพราะ FK จาก 0003 และ audit_log ลบไม่ได้ ต้องตัดสินทางเลือกใน plan หัวข้อ 8 ก่อน)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-30, T-09, T-19, T-25
- เสร็จเมื่อ: `TestAC_CXL_06_AuditLogAttempts` ผ่าน (3 แถว action `cancel`: `ok` 1, `forbidden` 2 มี request_id และ cause ใน query ตามตาราง FR-CXL-05 เขียนใน transaction เดียวกับการยกเลิก) และ `TestAC_CXL_07_ExpiredSessionNoCancel` ตรวจเพิ่มว่าไม่มีแถวใหม่
- สถานะ: พร้อมทำ

### T-32 บันทึกเมื่อส่ง Email ไม่สำเร็จ
- รองรับ: FR-CXL-09, FR-CXL-05
- ตรวจด้วย: AC-CXL-16 (ส่วน audit)
- ไฟล์ที่แตะ: `backend/internal/cancelrequest/service.go`, `backend/internal/cancelrequest/notifier_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-28, T-31
- เสร็จเมื่อ: `TestAC_CXL_16_NotifyFailureStillCancelled` ตรวจเพิ่มว่ามีแถว action `notify` outcome `error` และ request_id ของคำขอนั้นใน audit_log
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 17
- test ที่ต้องผ่านทั้งหมด: `TestMigration_AuditLogAcceptsNotify`, `TestAC_CXL_06_AuditLogAttempts`, `TestAC_CXL_07_ExpiredSessionNoCancel` (ส่วน audit), `TestAC_CXL_16_NotifyFailureStillCancelled` (ส่วน audit) และ test ของเฟส 0-16
- เปิดระบบ: ตามคำสั่งในส่วนหัว
- คนลองเองตาม AC-CXL-06:
  Given เข้าระบบเป็น student-001 (โหลด seed ใหม่)
  When ยกเลิกคำขอ requested ของตนเอง 1 ครั้ง, ส่งคำสั่งยกเลิกคำขอ approved ด้วย curl 1 ครั้ง และคำขอของ student-002 ด้วย curl 1 ครั้ง แล้วรัน `SELECT actor_id, action, outcome, request_id, query FROM audit_log WHERE action='cancel' ORDER BY id DESC LIMIT 3` ใน psql
  Then เห็น 3 แถวของ student-001 เป็น `ok` 1 แถว `forbidden` 2 แถว (cause `already_approved` และ `not_owner`)
- ยังขาด: ไม่มี

---

## ตารางตรวจความครบ

### 1. AC -> task
| AC ID | task ที่ตรวจ AC นี้ |
|---|---|
| AC-CXL-01 | T-06, T-07 |
| AC-CXL-02 | T-15, T-16 |
| AC-CXL-03 | T-19, T-20 |
| AC-CXL-04 | T-21, T-22 |
| AC-CXL-05 | T-13 |
| AC-CXL-06 | T-31 |
| AC-CXL-07 | T-25, T-26, T-31 (ส่วน audit) |
| AC-CXL-08 | T-04, T-05 |
| AC-CXL-09 | T-08 |
| AC-CXL-10 | T-14 |
| AC-CXL-11 | T-17, T-18 |
| AC-CXL-12 | T-11, T-12 |
| AC-CXL-13 | T-23, T-24 |
| AC-CXL-14 | T-09, T-10 |
| AC-CXL-15 | T-27 |
| AC-CXL-16 | T-28, T-32 (ส่วน audit) |
| AC-CXL-17 | T-29 |

### 2. Constraint -> task
| Constraint ID | task ที่ทำให้เป็นจริง |
|---|---|
| CON-TECH-01 | T-01 (ตารางใน PostgreSQL), T-30 (migration audit_log) |
| CON-TECH-02 | T-01 (รันใน PostgreSQL บน Docker ของกลาง) reverse proxy / TLS เป็นงานค้างของ 001 |
| CON-DATA-01 | T-01 (ไม่มี trigger แก้สถานะอุปกรณ์), T-06, T-13 |
| DOM-REQ-01 | T-13 |
| IF-IDP-01 | T-04, T-06 (อยู่หลัง `auth.Middleware`), T-19, T-25, T-26 |
| IF-NOT-01 | T-27, T-28 |
| CON-CXL-01 | T-01, T-06, T-13 |
| CON-CXL-02 | T-01, T-02, T-03, T-14, T-15, T-17 |

### 3. FR -> AC -> เฟส
| FR ID | AC ที่ตรวจ FR นี้ | เฟส |
|---|---|---|
| FR-CXL-01 | AC-CXL-01, AC-CXL-05, AC-CXL-10, AC-CXL-13 | 2, 6, 7, 12 |
| FR-CXL-02 | AC-CXL-02, AC-CXL-04, AC-CXL-11 | 8, 11, 9 |
| FR-CXL-03 | AC-CXL-03, AC-CXL-08 | 10, 1 |
| FR-CXL-04 | AC-CXL-04 | 11 |
| FR-CXL-05 | AC-CXL-06, AC-CXL-07 (ส่วน audit) | 17, 13 |
| FR-CXL-06 | AC-CXL-08, AC-CXL-12 | 1, 5 |
| FR-CXL-07 | AC-CXL-01, AC-CXL-04, AC-CXL-09, AC-CXL-14 | 2, 11, 3, 4 |
| FR-CXL-08 | AC-CXL-14 | 4 |
| FR-CXL-09 | AC-CXL-15, AC-CXL-16 | 14, 15, 17 |
| FR-CXL-10 | ไม่มี AC ควรเสนอทีมเพิ่ม (ไม่ทำใน prototype ตาม spec จึงไม่มี task) | - |

NFR-PERF-CXL-01 ตรวจด้วย AC-CXL-17 (เฟส 16)

## สิ่งที่ยังไม่ทำ
คัดลอกจาก Open Questions ใน spec

- ~~Q-13~~ ปิดแล้ว 2569-11-11 (ASM-12) T-30 ถึง T-32 เปลี่ยนเป็นพร้อมทำ ข้อความเดิม: DM-07 AuditLog มีฟิลด์ "รหัสครุภัณฑ์" แต่ไม่มีฟิลด์สำหรับรหัสคำขอ FR-CXL-05 ต้องบันทึกคำขอที่ยกเลิก
  ต้องเพิ่มฟิลด์ใน data-model.md (ของกลาง ใช้ร่วมกับ 001) หรือเก็บรหัสคำขอในช่องอื่น? -> ทีม (แจ้งเจ้าของ 001)
  task ที่รอ: T-30, T-31, T-32
- **Q-15** (ส่วนระบบจริง) รายชื่อ Email ของเจ้าหน้าที่แต่ละภาควิชามาจากไหน? IF-IDP-01 ส่งข้อมูลเฉพาะผู้ใช้ที่ล็อกอิน -> ถามนักเทคโนโลยีสารสนเทศ
  task ที่รอ: ไม่มีใน prototype (T-27 ใช้ Notifier จำลองตาม ASM-10) Notifier ตัวจริงยังไม่มี task
- **Q-16** (จาก FR-CXL-10 ตอบเมื่อนำเข้า scope) เจ้าหน้าที่ยกเลิกจากหน้าไหน ต้องระบุเหตุผลไหม ต้องแจ้งนักศึกษาไหม และยกเลิกได้เฉพาะคำขอของภาควิชาตนเองใช่ไหม? -> ถามเจ้าหน้าที่
  task ที่รอ: ไม่มี (FR-CXL-10 ไม่ทำใน prototype)
- คำถามของกลางที่เกี่ยวข้อง: Q-SH-03 (auth จริง ใช้ auth จำลองไปก่อน ไม่มี task), Q-SH-05 (ไม่แสดงลำดับคิวในหน้านี้ ไม่มี task), Q-SH-06 (ปิดแล้ว ปฏิเสธอัตโนมัติเป็น rejected 003 ไม่เปลี่ยนคำขออื่นเอง ไม่มี task), Q-SH-08 (ใช้ 401 ของ middleware กลาง ไม่มี task)
