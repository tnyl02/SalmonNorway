# Tasks: ตรวจสอบสถานะอุปกรณ์ (Equipment Status)
Spec: SPEC-EQS-001 (spec.md Draft v4.3) | Plan: plan.md Draft v1.5 | tasks v2: 2569-11-10 (v1 สร้าง 2569-11-02) | Owner: 660710767 นางสาวพรรณธร จันทร์สายทอง (@Pannathorn05)

สรุป
- 14 เฟสที่ใช้งาน (ฐานราก + 12 เฟสตาม AC + รอคำตอบ) รวม 36 task: เสร็จ 4, พร้อมทำ 19, รอ Q-xx 11, ยกเลิก 2 (T-15, T-25)
- รอ Q-xx 11 task: แบ่งหน้า (Q-05), วิธีค้นหา (Q-06), timeout (Q-08), จอมือถือ (Q-09), ขอบเขตเจ้าหน้าที่ (Q-P3), audit ของรายการ (Q-11, Q-P4),
  auth จริง (Q-SH-03), ปุ่มส่งคำขอยืม (Q-16), รีเฟรชเอง (Q-18), ชื่อผู้ยืม (Q-20), load test (Q-17)
- เฟสแรกที่ลองใช้งานได้จริง: เฟส 1 (เห็นรายการรุ่นและรายชิ้นบนหน้าจอ)

เปลี่ยนจาก v1: เลข task เดิมคงไว้ทั้งหมด (มี commit อ้างอยู่) task ใหม่เริ่มที่ T-32
- Q-15 ปิด: เพิ่ม T-32, T-33 (ตัวเลขตามตัวกรอง และเลือกหลายค่า AC-EQS-15, 16) / T-30 เหลือแค่ auth จริง
- Q-16 ยังเปิด: แยกเรื่องปุ่มส่งคำขอยืมออกจาก T-10, T-11, T-14, T-16 ไปไว้ที่ T-34 (รอ Q-16)
- Q-18, Q-20 ใหม่: T-35 รีเฟรชเอง (รอ Q-18), T-36 ชื่อผู้ยืมของเจ้าหน้าที่ แยกจาก T-17 (รอ Q-20)
- Q-17 ใหม่: T-22 load test ย้ายไปเฟสรอคำตอบ
- T-04 แตะ `cmd/server/main.go` ด้วย (รัน migration และผูก auth.Middleware ตาม plan ข้อ 10)

เทคโนโลยีตาม plan.md ข้อ 2: Go (net/http + pgx), PostgreSQL 16, React + Vite (ทีมเลือกเอง ไม่ได้มาจาก spec ยกเว้น PostgreSQL / HTTPS / Docker)
คำสั่งเปิดระบบ (ใช้ในทุกจุดตรวจ):
- ฐานข้อมูล: `docker compose -f deploy/docker-compose.yml up -d db`
- หลังบ้าน: `cd backend && go run ./cmd/server` (port 8000)
- หน้าจอ: `cd frontend && npm run dev` (port 5173)
- ผู้ใช้จำลอง: header `Authorization: Bearer mock-student-001` (ภาควิชา 1), `mock-student-101` (ภาควิชา 2), `mock-staff-001` (เจ้าหน้าที่)
- test: `cd backend && go test ./... -v` และ `cd frontend && npm test`

---

## เฟส 0: ฐานราก
หมายเหตุ: มี 4 task เกินกติกา 3 task เพราะ T-31 แยกออกมาทีหลังให้ 003 ใช้ตาราง audit_log ก่อน (2569-11-10)

### T-01 สร้างโครงหลังบ้าน Go หน้าจอ React และฐานข้อมูลใน Docker
- รองรับ: CON-TECH-01, CON-TECH-02
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของทุก task
- ไฟล์ที่แตะ: `backend/go.mod`, `backend/cmd/server/main.go`, `backend/internal/config/config.go`, `backend/internal/platform/db.go`, `backend/internal/platform/http.go`, `deploy/docker-compose.yml`, `frontend/` (โครงจาก template รายวิชา), `.gitignore`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: `GET /healthz` ตอบ 200 เมื่อต่อฐานข้อมูลได้ และ `npm test` ของโครงหน้าจอผ่าน
- สถานะ: เสร็จ

### T-02 สร้างตาราง read model และข้อมูลตั้งต้น
- รองรับ: CON-DATA-01, CON-DATA-03, ASM-01
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-04 ขึ้นไป
- ไฟล์ที่แตะ: `backend/db/migrations/0001_equipment_read_model.up.sql` (golang-migrate บังคับชื่อ .up.sql), `backend/db/seed/seed.sql`, `backend/internal/platform/migrate.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: `TestMigration_StatusCheck` ผ่าน (insert สถานะนอก 3 ค่าถูกปฏิเสธ) และ `TestMigration_AssetCodeUnique` ผ่าน (รหัสครุภัณฑ์ซ้ำถูกปฏิเสธ)
- สถานะ: เสร็จ 

### T-03 สร้าง middleware ยืนยันตัวตนแบบจำลอง
- รองรับ: IF-IDP-01, NFR-SEC-01 (ส่วน "ต้องยืนยันตัวตนก่อนทุกครั้ง")
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐาน (AC-EQS-09 ตรวจในเฟส 12)
- ไฟล์ที่แตะ: `backend/internal/auth/principal.go`, `backend/internal/auth/middleware.go`, `backend/internal/auth/mock_provider.go`, `backend/internal/auth/middleware_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: `TestAuth_NoToken_401` และ `TestAuth_MissingDepartment_Rejected` ผ่าน (ไม่ได้ภาควิชา = ปฏิเสธ ตาม plan 6.1)
- สถานะ: เสร็จ 

### T-31 สร้างตาราง audit_log ใช้ร่วม (แยกจาก T-19)
- รองรับ: DOM-AUDIT-01
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-19 และของ 003 (AC-EQS-11 ตรวจใน T-19)
- ไฟล์ที่แตะ: `backend/db/migrations/0002_audit_log.up.sql`, `backend/internal/platform/migrate_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-02
- เสร็จเมื่อ: `TestMigration_AuditLogAppendOnly` ผ่าน (INSERT ได้ / UPDATE และ DELETE ถูกปฏิเสธ) และ `TestMigration_AuditLogAction` ผ่าน (action นอกตาราง DM-07 ถูกปฏิเสธ, request_id ว่างได้)
- สถานะ: เสร็จ


#### จุดตรวจเฟส 0
- test ที่ต้องผ่านทั้งหมด: TestHealthz_OK_WhenDBReachable, TestHealthz_503_WhenDBDown, TestMigration_StatusCheck, TestMigration_AssetCodeUnique, TestSeed_Model1HasFiveItemsThreeAvailable, TestMigration_AuditLogAppendOnly, TestMigration_AuditLogAction, TestAuth_NoToken_401, TestAuth_MissingDepartment_Rejected, TestAuth_ValidToken_PassesPrincipal
- เปิดระบบ: ฐานข้อมูลและหลังบ้านตามคำสั่งด้านบน
- คนลองเอง (ไม่มี AC เป็นงานพื้นฐาน):
  Given เปิดฐานข้อมูลและหลังบ้านแล้ว
  When เรียก `GET http://localhost:8000/healthz`
  Then ได้ 200 `{"status":"ok"}`
- ยังขาด: ไม่มี (ผ่านแล้วทั้งเฟส)

---

## เฟส 1: เห็นรายการรุ่นและรายชิ้น (AC-EQS-01)

### T-04 สร้าง API รายการรุ่นพร้อมจำนวนว่าง/ทั้งหมด และผูกระบบใน main.go
- รองรับ: FR-EQS-01, DOM-REQ-01 (นับตามสถานะ "ว่าง" อย่างเดียว คำขอที่ค้างไม่ลดจำนวน), IF-IDP-01 (endpoint อยู่หลัง auth.Middleware)
- ตรวจด้วย: AC-EQS-01
- ไฟล์ที่แตะ: `backend/internal/equipment/repository.go`, `backend/internal/equipment/service.go`, `backend/internal/equipment/handler.go`, `backend/internal/equipment/dto.go`, `backend/internal/equipment/handler_test.go`, `backend/cmd/server/main.go` (เรียก platform.Migrate, platform.Seed และผูก auth.Middleware ตาม plan ข้อ 10)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-02, T-03
- เสร็จเมื่อ: `TestAC_EQS_01_ModelRowCounts` ผ่าน (รุ่น 5 ชิ้น ว่าง 3 -> 1 แถว availableCount=3 totalCount=5) และเปิดหลังบ้านแล้วเรียก `GET /api/v1/equipment/models` ด้วย token จำลองได้ข้อมูลตั้งต้น
- สถานะ: เสร็จ

### T-05 สร้าง API รายชิ้นของรุ่น
- รองรับ: FR-EQS-01, CON-DATA-03
- ตรวจด้วย: AC-EQS-01
- ไฟล์ที่แตะ: `backend/internal/equipment/repository.go`, `backend/internal/equipment/handler.go`, `backend/internal/equipment/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_EQS_01_ItemsOfModel` ผ่าน (ได้ 5 ชิ้นตามรหัสครุภัณฑ์ พร้อมสถานะแต่ละชิ้น)
- สถานะ: เสร็จ 

### T-06 สร้างหน้ารายการอุปกรณ์และรายชิ้น
- รองรับ: FR-EQS-01
- ตรวจด้วย: AC-EQS-01
- ไฟล์ที่แตะ: `frontend/src/pages/EquipmentList/EquipmentList.jsx`, `frontend/src/components/ModelRow.jsx`, `frontend/src/components/ItemTable.jsx`, `frontend/src/components/StatusBadge.jsx`, `frontend/src/api/client.js`, `frontend/src/__tests__/AC-EQS-01.test.jsx`, `frontend/src/App.jsx` (แสดงหน้ารายการ), `frontend/src/config.js` (token จำลอง ก่อนมีหน้าล็อกอินใน T-21) เพิ่ม 2 ไฟล์นี้ตามที่ทีมตัดสิน 2569-11-10
- แบบหน้าจอ: ไม่มี (ยังไม่มี mockup)
- ต้องทำหลัง: T-05
- เสร็จเมื่อ: `AC-EQS-01.test.jsx` ผ่าน (เห็น "ว่าง 3 / ทั้งหมด 5" และกดแล้วเห็น 5 ชิ้น) ยังไม่มีปุ่มส่งคำขอยืม (รอ Q-16 ดู T-34)
- สถานะ: เสร็จ 

#### จุดตรวจเฟส 1
- test ที่ต้องผ่านทั้งหมด: TestMigration_*, TestAuth_*, TestHealthz_*, TestAC_EQS_01_*, AC-EQS-01.test.jsx
- เปิดระบบ: ดูคำสั่งด้านบน ใช้ token `mock-student-001`
- คนลองเองตาม AC-EQS-01:
  Given ข้อมูลตั้งต้นมีรุ่น "Notebook ตัวอย่างรุ่น A" 5 ชิ้น ว่าง 3 ชิ้น
  When เปิดหน้ารายการอุปกรณ์ แล้วเลือกรุ่นนั้น
  Then เห็นรุ่นนั้น 1 แถว มีประเภท และ "ว่าง 3 / ทั้งหมด 5" และเห็น 5 ชิ้นตามรหัสครุภัณฑ์พร้อมสถานะ
- ยังขาด: แบ่งหน้า (รอ Q-05) / ปุ่มส่งคำขอยืม (รอ Q-16) / ลำดับคำขอแสดงที่หน้าไหน (รอ Q-SH-05)

---

## เฟส 2: ค้นหาและกรอง (AC-EQS-02, AC-EQS-15, AC-EQS-16)

### T-07 เพิ่มตัวกรองประเภทและสถานะใน API
- รองรับ: FR-EQS-02
- ตรวจด้วย: AC-EQS-02 (ส่วนกรองสถานะ)
- ไฟล์ที่แตะ: `backend/internal/equipment/repository.go`, `backend/internal/equipment/handler.go`, `backend/internal/equipment/handler_test.go`, `backend/internal/equipment/service.go` (ส่งตัวกรองผ่าน service เพิ่มตามที่ทีมตัดสิน 2569-11-10)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-05
- เสร็จเมื่อ: `TestAC_EQS_02_FilterStatus` ผ่าน (กรอง "ว่าง" ได้เฉพาะชิ้นที่ว่าง และแสดงเฉพาะรุ่นที่มีอย่างน้อย 1 ชิ้นตรงเงื่อนไข) parameterized query ทุกจุด
- สถานะ: เสร็จ

### T-32 ตัวกรองเลือกหลายค่า และตัวเลขของรุ่นนับตามตัวกรองใน API
- รองรับ: FR-EQS-01, FR-EQS-02 (คำตอบ Q-15: กรองที่รายชิ้น ตัวเลขเปลี่ยนตามตัวกรอง เลือกได้หลายค่า)
- ตรวจด้วย: AC-EQS-15, AC-EQS-16
- ไฟล์ที่แตะ: `backend/internal/equipment/repository.go`, `backend/internal/equipment/handler.go`, `backend/internal/equipment/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07
- เสร็จเมื่อ: `TestAC_EQS_15_CountsFollowFilter` ผ่าน (รุ่น 5 ชิ้น ว่าง 3 ถูกยืม 2 กรอง "ถูกยืม" -> 1 แถว ตัวเลขนับ 2) และ `TestAC_EQS_16_MultiStatusFilter` ผ่าน (กรอง "ว่าง" และ "ถูกยืม" พร้อมกัน ไม่มีชิ้นที่ไม่พร้อมใช้งาน)
- สถานะ: เสร็จ 

### T-08 เพิ่มช่องค้นหาและตัวกรองบนหน้าจอ
- รองรับ: FR-EQS-02
- ตรวจด้วย: AC-EQS-02 (ส่วนกรองสถานะ)
- ไฟล์ที่แตะ: `frontend/src/pages/EquipmentList/EquipmentList.jsx`, `frontend/src/components/SearchBar.jsx`, `frontend/src/__tests__/AC-EQS-02.test.jsx`, `frontend/src/api/client.js` (ส่งตัวกรองไปหลังบ้าน เพิ่มตามที่ทีมตัดสิน 2569-11-10)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07
- เสร็จเมื่อ: `AC-EQS-02.test.jsx` ส่วนกรองสถานะผ่าน
- สถานะ: เสร็จ 

### T-33 หน้าจอเลือกตัวกรองหลายค่า และแสดงตัวเลขตามตัวกรอง
- รองรับ: FR-EQS-01, FR-EQS-02
- ตรวจด้วย: AC-EQS-15, AC-EQS-16
- ไฟล์ที่แตะ: `frontend/src/components/SearchBar.jsx`, `frontend/src/components/ModelRow.jsx`, `frontend/src/__tests__/AC-EQS-15.test.jsx`, `frontend/src/__tests__/AC-EQS-16.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-32, T-08
- เสร็จเมื่อ: `AC-EQS-15.test.jsx` (กรอง "ถูกยืม" แล้วแถวรุ่นแสดงตัวเลข 2) และ `AC-EQS-16.test.jsx` (เลือก 2 สถานะพร้อมกัน) ผ่าน
- สถานะ: เสร็จ 

#### จุดตรวจเฟส 2
- test ที่ต้องผ่านทั้งหมด: test ของเฟส 1 + TestAC_EQS_02_FilterStatus + TestAC_EQS_15_* + TestAC_EQS_16_* + AC-EQS-02, 15, 16.test.jsx
- คนลองเองตาม AC-EQS-15 และ AC-EQS-16:
  Given รุ่น "Notebook ตัวอย่างรุ่น A" มี 5 ชิ้น ว่าง 3 ถูกยืม 2 (ข้อมูลตั้งต้น)
  When กรองสถานะ "ถูกยืม" แล้วเพิ่ม "ว่าง" เป็นตัวกรองที่ 2
  Then ตอนกรอง "ถูกยืม" รุ่นนั้นแสดงตัวเลข 2 / ตอนกรอง 2 สถานะ ไม่เห็น Projector ที่ไม่พร้อมใช้งาน
- ยังขาด: การจับคู่คำค้น "Notebook" (รอ Q-06 ดู T-24) / ข้อความกำกับตัวเลขเมื่อกรอง (รอ Q-21)

---

## เฟส 3: ค้นหาไม่พบ (AC-EQS-03)

### T-09 แสดงข้อความไม่พบและปุ่มล้างตัวกรอง
- รองรับ: FR-EQS-03
- ตรวจด้วย: AC-EQS-03
- ไฟล์ที่แตะ: `frontend/src/components/EmptyState.jsx`, `frontend/src/pages/EquipmentList/EquipmentList.jsx`, `frontend/src/__tests__/AC-EQS-03.test.jsx`, `backend/internal/equipment/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-08
- เสร็จเมื่อ: `TestAC_EQS_03_EmptyIs200` ผ่าน (ไม่พบ = 200 items ว่าง ไม่ใช่ error) และ `AC-EQS-03.test.jsx` ผ่าน
- สถานะ: เสร็จ 

#### จุดตรวจเฟส 3
- คนลองเองตาม AC-EQS-03:
  Given ค้นหาด้วยคำที่ไม่มีในระบบ
  When กดค้นหา
  Then เห็นข้อความไม่พบอุปกรณ์ และกดล้างตัวกรองแล้วกลับไปเห็นรายการทั้งหมด
- ยังขาด: ไม่มี

---

## เฟส 4: รายละเอียดอุปกรณ์ที่ถูกยืม (AC-EQS-04)

### T-10 สร้าง API รายละเอียดรายชิ้น พร้อมกำหนดคืน (วันที่และเวลา)
- รองรับ: FR-EQS-04 (ส่วนกำหนดคืน), CON-DATA-02
- ตรวจด้วย: AC-EQS-04 (ส่วนสถานะและกำหนดคืน)
- ไฟล์ที่แตะ: `backend/internal/equipment/repository.go`, `backend/internal/equipment/service.go`, `backend/internal/equipment/dto.go`, `backend/internal/equipment/handler_test.go`, `backend/internal/equipment/handler.go` (endpoint รายละเอียด เพิ่มตามที่ทีมตัดสิน 2569-11-10)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-05
- เสร็จเมื่อ: `TestAC_EQS_04_BorrowedShowsDueDate` ผ่าน (ถูกยืม ยังไม่เลยกำหนด -> มี dueAt วันที่และเวลา) และไม่มี endpoint ใดเขียนสถานะ (CON-DATA-02)
- สถานะ: เสร็จ 

### T-11 สร้างหน้ารายละเอียดอุปกรณ์
- รองรับ: FR-EQS-04 (ส่วนกำหนดคืน)
- ตรวจด้วย: AC-EQS-04 (ส่วนสถานะและกำหนดคืน)
- ไฟล์ที่แตะ: `frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx`, `frontend/src/__tests__/AC-EQS-04.test.jsx`, `frontend/src/components/ItemTable.jsx`, `frontend/src/pages/EquipmentList/EquipmentList.jsx` (กดรหัสครุภัณฑ์เพื่อเปิดรายละเอียด), `frontend/src/api/client.js` (getItem), `frontend/src/config.js` (เขตเวลา Asia/Bangkok) เพิ่มตามที่ทีมตัดสิน 2569-11-10
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `AC-EQS-04.test.jsx` ผ่าน (เห็น "ถูกยืม" และกำหนดคืนเป็นวันที่และเวลา)
- สถานะ: เสร็จ 

#### จุดตรวจเฟส 4
- คนลองเองตาม AC-EQS-04:
  Given ข้อมูลตั้งต้น CS-NB-0004 สถานะ "ถูกยืม" กำหนดคืน 16:00 วันนี้ (ลองก่อน 16:00)
  When นักศึกษาเปิดดูอุปกรณ์ชิ้นนั้น
  Then เห็นสถานะ "ถูกยืม" และกำหนดคืน (วันที่และเวลา)
- ยังขาด: "ไม่มีปุ่มส่งคำขอยืม" ของ AC-EQS-04 (รอ Q-16 ดู T-34)

---

## เฟส 5: เลยกำหนดคืน (AC-EQS-14, CR-001, CR-002)

### T-12 คำนวณ isOverdue ด้วยวันที่และเวลา Asia/Bangkok และแยก DTO นักศึกษา/เจ้าหน้าที่
- รองรับ: FR-EQS-04, DOM-PRIV-01, DOM-LOAN-01, CON-DATA-01 (ป้ายไม่ใช่สถานะที่ 4)
- ตรวจด้วย: AC-EQS-14
- ไฟล์ที่แตะ: `backend/internal/config/config.go`, `backend/internal/equipment/service.go`, `backend/internal/equipment/dto.go`, `backend/internal/equipment/service_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `TestAC_EQS_14_OverdueStudentVsStaff` ผ่าน (กำหนดคืน 2569-10-10 16:00 ตั้งนาฬิกา 17:00 วันเดียวกัน: นักศึกษาได้ isOverdue=true ไม่มี dueAt, เจ้าหน้าที่ได้ทั้งคู่ สถานะยังเป็น borrowed / ตั้งนาฬิกา 15:59 -> isOverdue=false)
- สถานะ: เสร็จ 

### T-13 แสดงป้าย "เลยกำหนดคืน"
- รองรับ: FR-EQS-04
- ตรวจด้วย: AC-EQS-14
- ไฟล์ที่แตะ: `frontend/src/components/StatusBadge.jsx`, `frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx`, `frontend/src/__tests__/AC-EQS-14.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-12, T-11
- เสร็จเมื่อ: `AC-EQS-14.test.jsx` ผ่าน (นักศึกษาเห็นป้ายไม่เห็นกำหนดคืน / เจ้าหน้าที่เห็นทั้งคู่)
- สถานะ: เสร็จ 

#### จุดตรวจเฟส 5
- คนลองเองตาม AC-EQS-14:
  Given ข้อมูลตั้งต้น CS-NB-0005 ยืมปกติ กำหนดคืน 16:00 เมื่อวาน
  When นักศึกษา (mock-student-001) เปิดดู แล้วเจ้าหน้าที่ (mock-staff-001) เปิดดูชิ้นเดียวกัน
  Then นักศึกษาเห็น "ถูกยืม" + ป้าย "เลยกำหนดคืน" ไม่เห็นกำหนดคืน / เจ้าหน้าที่เห็นกำหนดคืนและป้าย
- ยังขาด: "ไม่มีปุ่มส่งคำขอยืม" ของ AC-EQS-14 (รอ Q-16 ดู T-34) / กำหนดคืนของยืมค้างคืน (Q-SH-04) ไม่กระทบโค้ดเพราะอ่าน due_at จาก loan

---

## เฟส 6: อุปกรณ์ไม่พร้อมใช้งาน (AC-EQS-06)

### T-14 ป้าย "ไม่พร้อมใช้งาน"
- รองรับ: FR-EQS-05 (ส่วนป้าย), CON-DATA-01
- ตรวจด้วย: AC-EQS-06 (ส่วนป้าย)
- ไฟล์ที่แตะ: `frontend/src/components/StatusBadge.jsx`, `frontend/src/__tests__/AC-EQS-06.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-11
- เสร็จเมื่อ: `AC-EQS-06.test.jsx` ผ่าน (สถานะ unavailable แสดงป้าย "ไม่พร้อมใช้งาน")
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 6
- คนลองเองตาม AC-EQS-06:
  Given ข้อมูลตั้งต้น CS-PJ-0002 สถานะ "ไม่พร้อมใช้งาน"
  When นักศึกษาเปิดดูรายละเอียด
  Then เห็นป้าย "ไม่พร้อมใช้งาน"
- ยังขาด: "ไม่มีปุ่มส่งคำขอยืม" ของ AC-EQS-06 (รอ Q-16 ดู T-34)

---

## เฟส 7: [ยกเลิกโดย CR-002] อุปกรณ์รออนุมัติ (AC-EQS-12)

### T-15 [ยกเลิกโดย CR-002] ไม่นับ "รออนุมัติ" ในจำนวนพร้อมให้ยืม
- เหตุผล: FR-EQS-09 และ AC-EQS-12 ถูกยกเลิก ไม่มีสถานะรออนุมัติใน CON-DATA-01 แล้ว (RQ-09)
- สถานะ: ยกเลิก

---

## เฟส 8: คืนแล้วกลับมาว่าง (AC-EQS-05)

### T-16 ไม่ cache ข้อมูลสถานะ ให้รีเฟรชแล้วเห็นสถานะล่าสุด
- รองรับ: NFR-DATA-01, FR-EQS-06 (ส่วนสถานะ "ว่าง")
- ตรวจด้วย: AC-EQS-05 (ส่วนสถานะเปลี่ยนเป็น "ว่าง")
- ไฟล์ที่แตะ: `backend/internal/platform/http.go` (Cache-Control: no-store), `backend/internal/equipment/handler_test.go`, `frontend/src/api/client.js`, `frontend/src/__tests__/AC-EQS-05.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-11
- เสร็จเมื่อ: `TestAC_EQS_05_ReturnedBecomesAvailable` ผ่าน (ปิด loan จำลอง UC-09 แล้วเรียกใหม่ได้ status available และมี header no-store) และ `AC-EQS-05.test.jsx` ผ่าน (โหลดใหม่ไม่ใช้ข้อมูลเก่า)
- สถานะ: พร้อมทำ (ปุ่มส่งคำขอยืมกลับมาใช้ได้ ย้ายไป T-34 รอ Q-16 / หน้าจอรีเฟรชเองไหม รอ Q-18 ดู T-35)

#### จุดตรวจเฟส 8
- คนลองเองตาม AC-EQS-05:
  Given ปิด loan ของ CS-NB-0004 ในฐานข้อมูล และตั้งสถานะเป็น available (จำลองการรับคืนสภาพปกติของ UC-09)
  When นักศึกษากดรีเฟรชหน้ารายการภายใน 30 วินาที
  Then สถานะของชิ้นนั้นเป็น "ว่าง" และตัวเลขว่างของรุ่นเพิ่มขึ้น 1
- ยังขาด: ปุ่มส่งคำขอยืมกลับมาใช้ได้ (รอ Q-16) / รีเฟรชเองหรือกดเอง (รอ Q-18)

---

## เฟส 9: เห็นข้อมูลตามสิทธิ์ (AC-EQS-07)

### T-17 แยก DTO ให้นักศึกษาไม่มีข้อมูลผู้ยืม และ API สรุปจำนวนตามสถานะสำหรับเจ้าหน้าที่
- รองรับ: FR-EQS-07 (ส่วนนักศึกษาไม่เห็นผู้ยืม และสรุปตามสถานะ), DOM-PRIV-01
- ตรวจด้วย: AC-EQS-07 (ส่วนนักศึกษา)
- ไฟล์ที่แตะ: `backend/internal/equipment/dto.go`, `backend/internal/equipment/handler.go`, `backend/internal/equipment/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-12
- เสร็จเมื่อ: `TestAC_EQS_07_StudentNoBorrower` ผ่าน (response ของนักศึกษาไม่มีคำว่า borrower เลย) และ `TestFR_EQS_07_StaffSummaryByStatus` ผ่าน (เจ้าหน้าที่ได้จำนวนแยก 3 สถานะ นักศึกษาเรียกได้ 403)
- สถานะ: พร้อมทำ (ชื่อผู้ยืมฝั่งเจ้าหน้าที่ย้ายไป T-36 รอ Q-20)

#### จุดตรวจเฟส 9
- คนลองเองตาม AC-EQS-07:
  Given CS-NB-0004 ถูกยืมโดย student-001
  When student-002 (mock-student-002) เปิดดูรายละเอียด และเจ้าหน้าที่เรียกสรุปตามสถานะ
  Then นักศึกษาไม่เห็นข้อมูลผู้ยืมใน response / เจ้าหน้าที่เห็นจำนวนแยกตามสถานะ
- ยังขาด: เจ้าหน้าที่เห็นชื่อผู้ยืม (รอ Q-20 ดู T-36) / สรุปตามสถานะยังไม่มี AC ดู "สิ่งที่ยังไม่ทำ"

---

## เฟส 10: เฉพาะภาควิชาของตนเอง (AC-EQS-13)

### T-18 จำกัดทุก query ของนักศึกษาด้วยภาควิชา และตอบ 403 เมื่อขอข้ามภาควิชา
- รองรับ: FR-EQS-10, DOM-ACCESS-01, IF-IDP-01
- ตรวจด้วย: AC-EQS-13
- ไฟล์ที่แตะ: `backend/internal/equipment/repository.go`, `backend/internal/equipment/handler.go`, `backend/internal/equipment/handler_test.go`, `frontend/src/components/ErrorState.jsx`, `frontend/src/__tests__/AC-EQS-13.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: `TestAC_EQS_13_OtherDepartment403` ผ่าน ทั้งจาก list และเรียกรายละเอียดตรง body ไม่มีข้อมูลอุปกรณ์
- สถานะ: พร้อมทำ (ขอบเขตของเจ้าหน้าที่รอ Q-P3 ดู T-28 / ภาควิชาเดียวหรือหลายภาควิชา Q-19 ไม่กระทบวิธีตรวจ เพราะเทียบภาควิชาของผู้ใช้กับของอุปกรณ์เหมือนกัน)

#### จุดตรวจเฟส 10
- คนลองเองตาม AC-EQS-13:
  Given ล็อกอินจำลองเป็นนักศึกษาภาควิชา 2 (mock-student-101)
  When เปิดรายละเอียดของ CS-NB-0001 (ภาควิชา 1) โดยตรง
  Then เห็นข้อความไม่มีสิทธิ์ และไม่มีข้อมูลอุปกรณ์ใน response
- ยังขาด: ไม่มี

---

## เฟส 11: audit log ของการเปิดรายละเอียด (AC-EQS-11)

### T-19 บันทึก audit log เมื่อเปิดรายละเอียด
- รองรับ: DOM-AUDIT-01
- ตรวจด้วย: AC-EQS-11
- ไฟล์ที่แตะ: `backend/internal/audit/writer.go`, `backend/internal/audit/writer_test.go`, `backend/internal/equipment/handler.go` (ตารางสร้างแล้วใน T-31)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10, T-31
- เสร็จเมื่อ: `TestAC_EQS_11_DetailWritesAudit` ผ่าน (มีแถว action=detail ที่มีผู้เข้าถึง เวลา รหัสครุภัณฑ์)
- สถานะ: พร้อมทำ (การดูรายการ/ค้นหารอ Q-11 ดู T-29, ถ้าเขียน audit ไม่สำเร็จทำอย่างไรรอ Q-P4 ดู T-29)

#### จุดตรวจเฟส 11
- คนลองเองตาม AC-EQS-11:
  Given เปิดดูรายละเอียดอุปกรณ์ชิ้นหนึ่ง
  When การเข้าถึงเสร็จ
  Then ในตาราง audit_log มีแถวของผู้ใช้นั้น เวลา และรหัสครุภัณฑ์
- ยังขาด: audit ของรายการและการค้นหา (Q-11)

---

## เฟส 12: ฐานข้อมูลล่ม (AC-EQS-08)

### T-20 ตอบ 503 เมื่อต่อฐานข้อมูลไม่ได้ และหน้าจอแสดงปุ่มลองใหม่โดยไม่แสดงข้อมูลเก่า
- รองรับ: FR-EQS-08
- ตรวจด้วย: AC-EQS-08
- ไฟล์ที่แตะ: `backend/internal/platform/http.go`, `backend/internal/equipment/handler_test.go`, `frontend/src/components/ErrorState.jsx`, `frontend/src/pages/EquipmentList/EquipmentList.jsx`, `frontend/src/__tests__/AC-EQS-08.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestAC_EQS_08_DBDown503` ผ่าน (ปิดการเชื่อมต่อ -> 503 retryable) และ `AC-EQS-08.test.jsx` ผ่าน (รายการเดิมถูกล้างก่อนแสดง error)
- สถานะ: พร้อมทำ (กรณีต่อได้แต่ช้า ต้องรอกี่วินาทีรอ Q-08 ดู T-26)

#### จุดตรวจเฟส 12
- คนลองเองตาม AC-EQS-08:
  Given หยุด container ฐานข้อมูล (`docker compose -f deploy/docker-compose.yml stop db`)
  When เปิดหน้ารายการอุปกรณ์
  Then เห็นข้อความผิดพลาดพร้อมปุ่มลองใหม่ ไม่มีรายการเก่าค้าง
- ยังขาด: timeout (Q-08)

---

## เฟส 13: session หมดอายุ (AC-EQS-09)

### T-21 จัดการ 401 บนหน้าจอ: ไปหน้าล็อกอินจำลองพร้อมข้อความ แล้วกลับหน้าเดิม
- รองรับ: IF-IDP-01, NFR-SEC-01
- ตรวจด้วย: AC-EQS-09
- ไฟล์ที่แตะ: `frontend/src/api/client.js`, `frontend/src/pages/Login/MockLogin.jsx`, `frontend/src/__tests__/AC-EQS-09.test.jsx`, `backend/internal/auth/middleware_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-03, T-06
- เสร็จเมื่อ: `AC-EQS-09.test.jsx` ผ่าน (401 -> ไปหน้าล็อกอินพร้อมข้อความ -> กลับ path เดิม) และ returnTo ภายนอกระบบถูกปฏิเสธ
- สถานะ: พร้อมทำ (หน้าล็อกอินเป็นแบบจำลอง เพราะการยืนยันตัวตนเป็นระบบภายนอก IF-IDP-01 ไม่ใช่ use case ของระบบนี้)

#### จุดตรวจเฟส 13
- คนลองเองตาม AC-EQS-09:
  Given ลบ token ใน browser (จำลอง session หมดอายุ)
  When เรียกดูรายการอุปกรณ์
  Then ถูกส่งไปหน้าล็อกอินพร้อมข้อความ และหลังล็อกอินกลับมาหน้าเดิม
- ยังขาด: auth จริง (Q-SH-03) / auth ล่มตอนใช้งาน (Q-SH-08)

---

## เฟส 14: [ย้ายไปเฟสรอคำตอบ] ประสิทธิภาพ (AC-EQS-10)
T-22 ย้ายไปเฟส 15 เพราะยังไม่รู้ว่าต้องวัดอะไร (Q-17: 1,000 รายการนับชิ้นหรือรุ่น 50 คนวัดอย่างไร)

---

## เฟส 15: รอคำตอบ

### T-22 ทดสอบโหลด 1,000 รายการ ผู้ใช้พร้อมกัน 50 คน
- รองรับ: NFR-PERF-01
- ตรวจด้วย: AC-EQS-10
- ไฟล์ที่แตะ: `backend/db/seed/seed_1000.sql`, `tests/load/search.js` (k6 ผ่าน docker image `grafana/k6`), migration ใหม่สำหรับ index ค้นหา (เลขตามตารางเลข migration ใน data-model.md)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07
- เสร็จเมื่อ: รายงาน k6 แสดง p95 <= 3 วินาที บันทึกผลลง `tests/load/RESULT.md` (วัดบนเครื่อง Codespace ไม่ใช่ server จริง ต้องเขียนข้อจำกัดนี้ในผล)
- สถานะ: รอ Q-17

### T-23 แบ่งหน้ารายการรุ่น
- รองรับ: FR-EQS-01
- ตรวจด้วย: AC-EQS-01
- ไฟล์ที่แตะ: `backend/internal/config/config.go`, `backend/internal/equipment/repository.go`, `frontend/src/components/Pagination.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestAC_EQS_01_Pagination` ผ่าน ด้วยจำนวนต่อหน้าที่ได้จากคำตอบ
- สถานะ: รอ Q-05

### T-24 จับคู่คำค้นกับชื่อรุ่นและรหัสครุภัณฑ์
- รองรับ: FR-EQS-02 (ค้นจากชื่อหรือรหัสครุภัณฑ์ ตาม Q-14)
- ตรวจด้วย: AC-EQS-02
- ไฟล์ที่แตะ: `backend/internal/equipment/repository.go`, `backend/internal/equipment/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07
- เสร็จเมื่อ: `TestAC_EQS_02_SearchNotebook` ผ่าน
- สถานะ: รอ Q-06

### T-25 [ยกเลิกโดย CR-002] กำหนดสีของป้ายสถานะ
- เหตุผล: Q-07 ปิดแล้ว เหลือป้ายเดียว "ไม่พร้อมใช้งาน"
- สถานะ: ยกเลิก

### T-26 ตั้ง timeout ของ query และ request
- รองรับ: FR-EQS-08
- ตรวจด้วย: AC-EQS-08
- ไฟล์ที่แตะ: `backend/internal/config/config.go`, `frontend/src/config.js`, `frontend/src/api/client.js`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-20
- เสร็จเมื่อ: `TestAC_EQS_08_SlowDBTimeout` ผ่าน ด้วยค่าที่ได้จากคำตอบ
- สถานะ: รอ Q-08

### T-27 ตรวจหน้าจอที่ความกว้างมือถือ
- รองรับ: NFR-USE-01
- ตรวจด้วย: ไม่มี AC ควรเสนอทีมเพิ่ม
- ไฟล์ที่แตะ: `frontend/src/index.css`, หน้าจอทุกหน้าของฟีเจอร์นี้
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-13
- เสร็จเมื่อ: หน้าจอไม่ล้นแนวนอนที่ความกว้างที่ได้จากคำตอบ
- สถานะ: รอ Q-09

### T-28 ขอบเขตภาควิชาของเจ้าหน้าที่
- รองรับ: FR-EQS-07, DOM-ACCESS-01
- ตรวจด้วย: AC-EQS-07
- ไฟล์ที่แตะ: `backend/internal/auth/principal.go` (`scopeFor`), `backend/internal/equipment/dto.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-17, T-18
- เสร็จเมื่อ: test ตามคำตอบผ่าน
- สถานะ: รอ Q-P3

### T-29 audit ของการดูรายการและการค้นหา และกรณีเขียน audit ไม่สำเร็จ
- รองรับ: DOM-AUDIT-01
- ตรวจด้วย: AC-EQS-11
- ไฟล์ที่แตะ: `backend/internal/audit/writer.go`, `backend/internal/equipment/handler.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-19
- เสร็จเมื่อ: test ตามคำตอบผ่าน
- สถานะ: รอ Q-11, Q-P4

### T-30 เชื่อม auth service จริง
- รองรับ: IF-IDP-01
- ตรวจด้วย: AC-EQS-09
- ไฟล์ที่แตะ: `backend/internal/auth/` (Provider ตัวจริงแทน MockProvider)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-21
- เสร็จเมื่อ: test ตามคำตอบผ่าน
- สถานะ: รอ Q-SH-03 (เรื่องตัวเลขเมื่อกรองที่เคยอยู่ในนี้ ปิดด้วย Q-15 ย้ายไป T-32)

### T-34 ปุ่มส่งคำขอยืม: แสดงเมื่อว่าง ซ่อนเมื่อถูกยืมหรือไม่พร้อมใช้งาน และส่งต่อไป UC-02
- รองรับ: FR-EQS-04, FR-EQS-05, FR-EQS-06 (ส่วนปุ่ม)
- ตรวจด้วย: AC-EQS-04, AC-EQS-05, AC-EQS-06, AC-EQS-14 (ส่วน "มี/ไม่มีปุ่มส่งคำขอยืม")
- ไฟล์ที่แตะ: `backend/internal/equipment/service.go`, `backend/internal/equipment/dto.go` (canRequestLoan), `backend/internal/equipment/handler_test.go`, `frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx`, `frontend/src/components/ModelRow.jsx`, `frontend/src/__tests__/RequestButton.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-16, T-14, T-13, T-06
- เสร็จเมื่อ: test ตามคำตอบ Q-16 ผ่าน (ปุ่มอยู่หน้าไหน ส่งรายชิ้นหรือรุ่นไป UC-02)
- สถานะ: รอ Q-16 (ถ้าคำตอบเปลี่ยน FR-EQS-06 เป็นขอยืมระดับรุ่น ต้องทำ CR ก่อน)

### T-35 หน้าจอรีเฟรชสถานะเอง
- รองรับ: NFR-DATA-01
- ตรวจด้วย: AC-EQS-05
- ไฟล์ที่แตะ: `frontend/src/pages/EquipmentList/EquipmentList.jsx`, `frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx`, `frontend/src/config.js`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-16
- เสร็จเมื่อ: test ตามคำตอบผ่าน (ถ้าคำตอบคือผู้ใช้กดเอง task นี้ยกเลิก เพราะ T-16 ครอบคลุมแล้ว)
- สถานะ: รอ Q-18

### T-36 เจ้าหน้าที่เห็นชื่อผู้ยืมปัจจุบัน
- รองรับ: FR-EQS-07 (ส่วนชื่อผู้ยืม), DOM-PRIV-01
- ตรวจด้วย: AC-EQS-07 (ส่วนเจ้าหน้าที่)
- ไฟล์ที่แตะ: `backend/internal/equipment/dto.go`, `backend/internal/equipment/handler_test.go`, `frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx`, `frontend/src/__tests__/AC-EQS-07.test.jsx`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-17
- เสร็จเมื่อ: `TestAC_EQS_07_StaffSeesBorrower` ผ่าน ด้วยชื่อจากแหล่งตามคำตอบ
- สถานะ: รอ Q-20 (ระบบเก็บแค่ borrower_id ยังไม่รู้ว่าชื่อมาจากไหน)

---

## ตารางตรวจความครบ

### 1. AC -> task
| AC ID | task ที่ตรวจ AC นี้ |
|---|---|
| AC-EQS-01 | T-04, T-05, T-06, T-23 |
| AC-EQS-02 | T-07, T-08, T-24 |
| AC-EQS-03 | T-09 |
| AC-EQS-04 | T-10, T-11, T-34 |
| AC-EQS-05 | T-16, T-34, T-35 |
| AC-EQS-06 | T-14, T-34 |
| AC-EQS-07 | T-17, T-36, T-28 |
| AC-EQS-08 | T-20, T-26 |
| AC-EQS-09 | T-21, T-30 |
| AC-EQS-10 | T-22 |
| AC-EQS-11 | T-19, T-29 |
| AC-EQS-12 | (ยกเลิกโดย CR-002) |
| AC-EQS-13 | T-18 |
| AC-EQS-14 | T-12, T-13, T-34 |
| AC-EQS-15 | T-32, T-33 |
| AC-EQS-16 | T-32, T-33 |

### 2. Constraint -> task (รวมที่ใช้ร่วมจาก specs/000-shared/constraints.md)
| Constraint ID | task ที่ทำให้เป็นจริง |
|---|---|
| CON-TECH-01 | T-01, T-02 |
| CON-TECH-02 | T-01 (Docker) ส่วน HTTPS/TLS ที่ reverse proxy: ยังไม่มี task ดูหมายเหตุ |
| CON-DATA-01 | T-02, T-12, T-14 |
| CON-DATA-02 | T-10 |
| CON-DATA-03 | T-02, T-05 |
| DOM-PRIV-01 | T-12, T-17, T-36 |
| DOM-AUDIT-01 | T-31, T-19, T-29 |
| DOM-ACCESS-01 | T-18, T-28 |
| DOM-LOAN-01 | T-12 |
| DOM-REQ-01 | T-04 |
| IF-IDP-01 | T-03, T-04, T-18, T-21, T-30 |

หมายเหตุ CON-TECH-02 และ NFR-SEC-01 (TLS 1.2): การตั้ง reverse proxy ด้วย HTTPS เป็นงานของทั้งระบบ ไม่ใช่ของฟีเจอร์นี้คนเดียว
เสนอให้ทีมตัดสินว่าใครทำ (ของกลาง) ก่อนส่งงาน

### 3. FR -> AC -> เฟส
| FR ID | AC ที่ตรวจ FR นี้ | เฟส |
|---|---|---|
| FR-EQS-01 | AC-EQS-01, AC-EQS-15 | 1, 2, 15 |
| FR-EQS-02 | AC-EQS-02, AC-EQS-15, AC-EQS-16 | 2, 15 |
| FR-EQS-03 | AC-EQS-03 | 3 |
| FR-EQS-04 | AC-EQS-04, AC-EQS-14 | 4, 5, 15 |
| FR-EQS-05 | AC-EQS-06 | 6, 15 |
| FR-EQS-06 | AC-EQS-05 | 8, 15 |
| FR-EQS-07 | AC-EQS-07 (เฉพาะชื่อผู้ยืม) | 9, 15 |
| FR-EQS-08 | AC-EQS-08 | 12, 15 |
| FR-EQS-09 | (ยกเลิกโดย CR-002) | - |
| FR-EQS-10 | AC-EQS-13 | 10 |

FR-EQS-07 ส่วน "สรุปจำนวนอุปกรณ์แยกตามสถานะ" ไม่มี AC ควรเสนอทีมเพิ่ม (T-17 ใช้ test ชื่อ TestFR_EQS_07_StaffSummaryByStatus ไปก่อน)
NFR ที่ไม่มี AC: NFR-SEC-01 (บางส่วนตรวจผ่าน AC-EQS-09), NFR-AVAIL-01, NFR-USE-01 -> ไม่มี AC ควรเสนอทีมเพิ่ม
NFR-USE-01 ต้องทดสอบกับผู้ใช้จริง 10 คน เป็นงานของคน ไม่ใช่ task เขียนโค้ด

## สิ่งที่ยังไม่ทำ (Open Questions จาก spec และของกลาง)
| คำถาม | task ที่รอ |
|---|---|
| Q-04 ประเภทอุปกรณ์ตายตัวหรือเพิ่มได้ | ไม่มี (ใช้ข้อมูลตั้งต้น ตาม backlog) |
| Q-05 แบ่งหน้ากี่รายการ | T-23 |
| Q-06 วิธีจับคู่คำค้น | T-24 |
| Q-08 timeout | T-26 |
| Q-09 ความกว้างจอมือถือ | T-27 |
| Q-11 audit ของรายการ/ค้นหา | T-29 |
| Q-16 ปุ่มส่งคำขอยืมอยู่หน้าไหน นักศึกษาเลือกชิ้นเองได้ไหม | T-34 |
| Q-17 1,000 รายการ / 50 คน วัดอย่างไร | T-22 |
| Q-18 หน้าจอรีเฟรชเองไหม | T-35 |
| Q-19 ภาควิชาเดียวหรือหลายภาควิชา | ไม่มี task รอ (T-18 ทำได้ทั้งสองแบบ) |
| Q-20 ชื่อผู้ยืมมาจากไหน | T-36 |
| Q-21 ข้อความกำกับตัวเลขเมื่อกรองสถานะที่ไม่ใช่ "ว่าง" | T-33 (เฉพาะข้อความกำกับ ส่วนนับตัวเลขทำได้แล้ว) |
| Q-P3 ขอบเขตภาควิชาของเจ้าหน้าที่ | T-28 |
| Q-P4 เขียน audit ไม่สำเร็จ | T-29 |
| Q-SH-03 โปรโตคอล auth | T-30 |
| Q-SH-04 กำหนดคืนของยืมค้างคืน | ไม่มี (กระทบฟีเจอร์ 005) |
| Q-SH-05 ลำดับคำขอแสดงให้ใครและที่หน้าไหน | ไม่มี task (ถ้าคำตอบคือหน้ารายละเอียดของ 001 ต้องแก้ spec ผ่าน CR) |
| Q-SH-08 auth ล่มตอนใช้งานทั่วไป | ไม่มี task (ถ้าได้คำตอบ ต้องเพิ่ม AC และ task ต่อจาก T-21) |
| Q-SH-09 คืนแล้วชำรุด ใครเปลี่ยนสถานะ | ไม่มี (001 อ่านอย่างเดียว กระทบ 006) |
| ปิดแล้ว: Q-01, 02, 03, 07, 10, 12, 13, 14, 15, Q-SH-01, Q-SH-02, Q-P2 (ตอบใน Q-15), Q-P5 | - |

เรื่องที่ spec ยังไม่บอก และยังไม่มี Q (เสนอทีมเพิ่ม ห้ามเดา):
- FR-EQS-07: สรุปจำนวนแยกตามสถานะยังไม่มี AC
