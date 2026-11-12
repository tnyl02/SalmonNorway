# Prompt log

บันทึกทุกครั้งที่ใช้ AI กับ repo นี้ เขียนต่อท้ายเรื่อย ๆ ห้ามแก้หรือลบบันทึกเดิม
ทุกบันทึกต้องมี: วันที่ ผู้สั่ง (สมาชิกคนไหน) ฟีเจอร์ เครื่องมือ สิ่งที่สั่ง คำถามที่ AI ถาม คำตอบของทีม และสิ่งที่ AI ทำผิดหรือเดา

> หมายเหตุ: งานก่อนวันที่ 2569-11-02 (spec v1-v3, plan v1.1, backlog, CR-001) ไม่ได้บันทึกไว้ในไฟล์นี้
> ไฟล์นี้เริ่มบันทึกตั้งแต่ 2569-11-02 เป็นต้นไป

---

## 2569-11-02 วางฐานกลางของ repo และตรวจคำกำกวม

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: ของกลาง (000-shared, docs/srs) และ 001-equipment-status
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: สรุปงาน Term Project จากหน้าเรียน W1-W10 และข้อกำหนดที่อาจารย์สั่งในห้อง
  (ส่ง SRS, spec.md, prototype 6 ฟีเจอร์คนละฟีเจอร์, RTM ทำท้ายสุด, ต้องรู้ว่าใครทำอะไร, ตรวจคำกำกวม)
  แล้วให้วางฐานกลางและทำฟีเจอร์ 001

### คำถามที่ AI ถาม และคำตอบของทีม
1. 6 ฟีเจอร์ใช้ use case UC-01, 02, 07, 08, 09, 10 ได้ไหม -> ทีม: แยกให้เล็กลง
   (AI เสนอร่าง: เพิ่ม "ยกเลิกคำขอยืม" แทน UC-10 รอทีมยืนยันใน specs/README.md)
2. SRS อยู่ที่ไหน -> ทีม: เขียนใหม่ใน repo (docs/srs/SRS.md) diagram อยู่ใน Miro ไม่ต้องใส่ repo
3. ให้ AI ช่วยแค่ไหน -> ทีม: วางฐานกลาง + ฟีเจอร์ 001 สมาชิกคนอื่นทำฟีเจอร์ตัวเองจากบัญชีตัวเอง
4. เทคโนโลยี -> ทีม: Go + React ตาม plan.md

### สิ่งที่ AI ทำ
- ติดตั้ง agent-pack จาก template รายวิชา (ppsajja/swreqspec-template): คำสั่ง 7 ตัว + docs/*-guide.md
- เขียน AGENTS.md ใหม่ เพิ่มตารางเทียบ Python -> Go เพราะคำสั่งของรายวิชาเขียนไว้สำหรับ Python
- สร้าง specs/000-shared/ (constraints, data-model, glossary) ย้าย Constraint ใช้ร่วมจาก 001 ข้อความและ ID เดิม
- spec 001 -> Draft v3.1, plan 001 -> v1.2 (โครงโฟลเดอร์ backend/ frontend/)
- แก้ลิงก์ backlog ที่เสีย 3 ไฟล์, รหัส BL-11/BL-12 ที่ปนกันใน backlog.md, ตัวเลขกฎ 60% (27 -> 24, 81.5% -> 91.7%)
- ร่าง docs/srs/SRS.md

### ข้อค้นพบจากการตรวจคำกำกวม (เพิ่มเป็น Q ใน spec ไม่ได้แก้ FR เอง)
- spec 001: Q-05 แบ่งหน้า, Q-06 วิธีค้นหา, Q-07 "ชัดเจน", Q-08 timeout, Q-09 "หน้าจอมือถือ",
  Q-10 สิทธิ์อาจารย์, Q-11 ขอบเขต audit, Q-12 ที่มาของกำหนดคืน
- ของกลาง: Q-SH-01 สถานะหลังอนุมัติก่อนส่งมอบ (6 สถานะใน CON-DATA-01 ไม่มีค่าที่ตรง), Q-SH-02 ใครอนุมัติ, Q-SH-03 โปรโตคอล auth
- plan 001 v1.1 เคยเดาค่า page size 20 และ timeout 5/8 วินาที โดยไม่มีที่มาใน spec -> กลายเป็น Q-05, Q-08

### สิ่งที่ AI เดา หรือทีมต้องตรวจ
- ร่าง SRS หัวข้อ 3.1.2-3.1.6 เป็นโครงว่าง ให้เจ้าของฟีเจอร์เขียนเอง AI ไม่ได้แต่ง FR ของฟีเจอร์อื่น
- ฟีเจอร์ 003 (ยกเลิกคำขอยืม) ยังไม่มีที่มาใน backlog ต้องยืนยันว่ามีใน user stories บน Miro

---

## 2569-11-02 /tasks ของ 001-equipment-status

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code (ทำตามกติกาใน .claude/commands/tasks.md)
- ผลลัพธ์: specs/001-equipment-status/tasks.md 16 เฟส 30 task รอ Q-xx 8 task
- task ที่ยากที่สุด: T-12 (isOverdue ตาม Asia/Bangkok + แยก DTO ตามบทบาท) เพราะต้องทดสอบด้วยนาฬิกาที่ตั้งได้ และพลาดแล้วชื่อผู้ยืมหลุดถึงนักศึกษา
- AC ที่ทดสอบยากในสภาพแวดล้อมนักศึกษา: AC-EQS-10 (50 ผู้ใช้พร้อมกัน) เสนอวัดด้วย k6 บน Codespace แล้วเขียนข้อจำกัดในผล ทีมต้องตัดสินว่ายอมรับไหม
- สิ่งที่ AI ไม่ได้ทำ: HTTPS ที่ reverse proxy (CON-TECH-02) ไม่ได้ใส่เป็น task ของ 001 เพราะเป็นงานทั้งระบบ เสนอให้ทีมเลือกคนรับผิดชอบ
- ยังไม่เริ่มทำ task ใด ๆ

---

## 2569-11-02 ตรวจเทียบ repo กับ Miro แล้วแก้ให้ตรง (CR-002)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: ของกลาง, 001-equipment-status, backlog, SRS
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ส่งภาพจาก Miro (Journey Map, Stakeholder Register, System Context, Onion, User Story, Infrastructure,
  Requirement table, Use Case Spec, Use Case Diagram, Architecture, Backlog) ให้ตรวจว่าตรงไหนไม่เข้ากัน
  แล้วสั่ง: "อะไรที่ไม่ตรงกับ Miro เอาให้ตรงก่อน" (สถานะ, ยืมค้างคืน, กำหนดคืน 16:00, ผู้อนุมัติค้างคืน)

### สิ่งที่ AI ทำ
- เขียน CR-002 แล้วแก้ spec 001 (v4), constraints/data-model/glossary (v2), plan (v1.3), tasks, backlog (v4), SRS (v0.2), specs/README
- ยกเลิก FR-EQS-09 / AC-EQS-12 / T-15 / T-25 (ไม่ลบ ID ทิ้ง เขียนว่ายกเลิกโดย CR-002)

### สิ่งที่ AI ตัดสินเองเมื่อ Miro ขัดกันเอง (ทีมต้องตรวจ)
- สถานะ: ยึด Requirement ข้อ 9 (3 ค่า) ไม่ใช้ US-001 / Onion
- ผู้อนุมัติค้างคืน: ยึด Requirement ข้อ 15 + UC-14 + Architecture (อาจารย์เซ็นนอกระบบ) ไม่ใช้ US-012 / Infrastructure
- BL-04 ใช้ข้อความตาม Miro ("เจ้าหน้าที่/อาจารย์") แม้จะขัดกับ Use Case Diagram ที่ไม่มีอาจารย์เป็น actor -> แจ้งทีม
- NFR-PERF-01 เปลี่ยนเป็น 3 วินาทีตามข้อ 24 แต่คงเงื่อนไข 50 ผู้ใช้ไว้เป็นการตีความ "สภาวะปกติ" -> ต้องยืนยัน
- BL-20 (ยกเลิกคำขอ) ใส่ Must เพราะทีมเลือกเป็นฟีเจอร์ prototype ส่วน BL-12 ถึง 19, 21 ถึง 23 ไม่ได้จัดลำดับให้
- ภาพ Use Case Specification ความละเอียดต่ำ อ่านไม่ได้ ยังไม่ได้ตรวจ

---

## 2569-11-09 ทีมตอบ Q-SH-01

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: ของกลาง (กระทบ 001, 002, 004, 005)
- เครื่องมือ: Claude Code
- คำตอบของทีม: (ก) ส่งคำขอซ้อนได้ ใครได้รับอนุมัติก่อนได้ไป แต่ต้องแสดงเวลาและลำดับของคนที่ขอมาก่อน
- สิ่งที่ AI ทำ: ปิด Q-SH-01 เพิ่ม DOM-REQ-01 ใน constraints.md (v3) แก้ data-model, plan/tasks ของ 001, SRS, CHANGELOG
- สิ่งที่ AI ไม่ได้เดา เพิ่มเป็นคำถาม:
  - Q-SH-05 ลำดับแสดงให้ใคร (นักศึกษา/เจ้าหน้าที่) ถ้านักศึกษาเห็น จะขัด DOM-PRIV-01 ไหมถ้ามีชื่อ และแสดงที่หน้าของ 001 หรือ 002
  - Q-SH-06 หลังมีคำขอได้รับอนุมัติ คำขอที่เหลือถูกปฏิเสธอัตโนมัติไหม แจ้งผู้ขอไหม และต้องอนุมัติตามลำดับหรือข้ามได้
- หมายเหตุ: คำตอบนี้ทีมตัดสินเอง ยังไม่ได้ยืนยันกับนักคอมพิวเตอร์ ตามที่ constraints.md เคยระบุ

---

## 2569-11-09 จัดลำดับ backlog BL-12 ถึง BL-23 (backlog v5)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: backlog (ไม่กระทบ spec / plan / tasks ของ 001)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ให้ AI ร่าง MoSCoW, Kano, Value/Effort ของเรื่องใหม่จาก CR-002 แล้วคำนวณกฎ 60% ใหม่
  ทีมนำร่างไปแก้ใน Miro ส่งภาพกลับมาให้ AI ตรวจจนตรงกัน จากนั้นสั่งให้อัปเดต repo ตาม Miro
- สิ่งที่ AI ทำ: เขียน backlog.md หัวข้อ 1-4 ใหม่ให้ตรงกับ 4 ตารางบน Miro, กฎ 60% = 24/53 = 45.3% ผ่าน, SRS 3.1.7 ใส่ MoSCoW, CHANGELOG
- สิ่งที่ AI ร่างเอง ทีมต้องรู้:
  - Kano ของ BL-12 ถึง BL-23 เป็นการประเมิน ไม่ได้ถามผู้ใช้ (เขียนหมายเหตุไว้ใน backlog.md แล้ว)
  - BL-13 เป็น Won't เพราะ UC-04 ไม่มี RQ รองรับ ถ้าทีมเก็บ UC-04 ไว้ใน Miro ต้องทบทวน
  - ผ่านกฎ 60% เพราะตัวหารโตขึ้น ไม่ใช่เพราะ Must ลดลง เขียนข้อแลกเปลี่ยนไว้ใน backlog.md

---

## 2569-11-09 ทีมแก้จุดที่ Miro ขัดกันเอง

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: ของกลาง (ตาม CR-002)
- เครื่องมือ: Claude Code
- ทีมแจ้งว่าแก้ใน Miro แล้ว: (1) สถานะอุปกรณ์ 3 ค่า ว่าง / ถูกยืม / ไม่พร้อมใช้งาน (2) ค้างคืน อาจารย์เซ็นอนุมัตินอกระบบ เจ้าหน้าที่บันทึกผล
- ทีมยืนยันการเลือกที่ AI ทำไว้ใน CR-002 (ยึดตาราง Requirement ข้อ 9 และข้อ 15)
- สิ่งที่ AI ทำ: บันทึกผลใน SRS 4.3 และเพิ่มบรรทัดติดตามผลท้ายรายการใน CR-002 ไม่ได้เปลี่ยนข้อกำหนดใด ๆ

---

## 2569-11-09 เติมชื่อเจ้าของใน SRS และเปลี่ยนช่อง [เติม] ที่ไม่มีคำตอบเป็น Q

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: SRS, ของกลาง, 001
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ทีมพบว่า SRS ยังมี [เติม] 13 จุด (FR ของ 002-006, ช่องเจ้าของ, บทลงโทษ RQ-27)
- สิ่งที่ AI ทำ: ใส่ชื่อเจ้าของ 6 ฟีเจอร์ตาม README.md และ Owner ของ spec 001,
  เปลี่ยนเรื่อง RQ-27 และ auth ล่มเป็น Q-SH-07, Q-SH-08 ใน constraints.md
- สิ่งที่ AI ไม่ได้ทำ (ตั้งใจ): ไม่เขียน FR ของ 002-006 เพราะเป็นงานของเจ้าของแต่ละฟีเจอร์ (อาจารย์ต้องเห็นว่าใครทำอะไร)
  และไม่เดารูปแบบบทลงโทษ เพราะต้องถามผู้ใช้

---

## 2569-11-09 /implement T-01 (001-equipment-status)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (เฟส 0 ฐานราก ใช้ร่วมกับทุกฟีเจอร์)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-01 specs/001-equipment-status/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- `backend/go.mod`, `backend/go.sum` (go.sum สร้างอัตโนมัติจาก go mod tidy ไม่อยู่ในช่องไฟล์ที่แตะ)
- `backend/cmd/server/main.go`, `backend/internal/config/config.go`
- `backend/internal/platform/db.go`, `backend/internal/platform/http.go`
- `backend/internal/platform/http_test.go` (test ของ "เสร็จเมื่อ" ไม่อยู่ในช่องไฟล์ที่แตะ)
- `deploy/docker-compose.yml`
- `frontend/` คัดลอกจาก template รายวิชา (ppsajja/swreqspec-template) แก้ชื่อระบบใน index.html, App.jsx, setup.test.jsx, package.json
  และลบฟังก์ชัน getSlots / createBooking ของตัวอย่างจองคิวออกจาก src/api/client.js (ไม่ใช่ระบบของเรา)
- `.gitignore` ไม่ต้องแก้ (มี node_modules/ และ dist/ แล้ว)

### ผล test
- `cd backend && go test ./... -v`: TestHealthz_OK_WhenDBReachable PASS, TestHealthz_503_WhenDBDown PASS
- `cd frontend && npm test`: setup.test.jsx 1 passed
- ลองกับของจริง: docker compose db ขึ้น -> `curl localhost:8000/healthz` = 200 {"status":"ok"} / หยุด db -> 503 {"status":"db_unavailable"}

### สิ่งที่ AI เกือบต้องเดา (เลือกแบบเปลี่ยนน้อยที่สุด ทีมต้องยืนยัน)
- port หลังบ้าน 8000: ใช้ตาม proxy ใน vite.config.js ของ template แทนการเลือกเลขเอง
- รหัสผ่านฐานข้อมูล dev (salmon/salmon): ค่าของเครื่องพัฒนา ไม่ใช่ requirement อ่านจาก env DATABASE_URL ได้
- go directive 1.25.0 (pgx v5.12 ต้องการขั้นต่ำ 1.25) เพื่อให้เพื่อนที่ใช้ Go เก่ากว่าเครื่องนี้ build ได้
- ไม่ได้ใส่ HTTPS/reverse proxy (CON-TECH-02 ส่วน TLS) เพราะ tasks.md ระบุว่าเป็นงานของทั้งระบบที่ยังไม่มีเจ้าของ

---

## 2569-11-09 /implement T-02 (001-equipment-status)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (เฟส 0 ฐานราก ตารางใช้ร่วมกับ 002-006)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-02 specs/001-equipment-status/tasks.md`
- รอบแรกถูกปฏิเสธตามกติกา เพราะ T-01 ยังเป็น "เสร็จ รอทีมตรวจ" ทีมตรวจแล้วเปลี่ยนเป็น "เสร็จ" (commit 04f0fd7) จึงเริ่มรอบนี้

### ไฟล์ที่สร้างหรือแก้
- `backend/db/migrations/0001_equipment_read_model.up.sql` (**ชื่อต่างจาก tasks.md** ที่เขียน `0001_equipment_read_model.sql`
  เพราะ golang-migrate ที่ plan เลือกบังคับให้ลงท้าย `.up.sql`)
- `backend/db/seed/seed.sql`, `backend/internal/platform/migrate.go`
- `backend/internal/platform/migrate_test.go` (test ของ "เสร็จเมื่อ")
- `backend/go.mod`, `backend/go.sum` (เพิ่ม golang-migrate v4.20.1, testcontainers-go v0.44.0 ตาม plan ข้อ 2)

### ผล test (`cd backend && go test ./... -v`)
- TestMigration_StatusCheck PASS (3 ค่าบันทึกได้ / pending_approval, damaged, maintenance, lost, ค่าว่าง ถูกปฏิเสธ)
- TestMigration_AssetCodeUnique PASS
- TestSeed_Model1HasFiveItemsThreeAvailable PASS (test เพิ่ม ยืนยันข้อมูลตั้งต้นใช้กับ AC-EQS-01 ได้)
- TestHealthz_* ของ T-01 ยัง PASS, `go test -short ./...` ข้าม test ที่ต้องใช้ Docker ได้
- รอบแรก test ที่ 2-3 ล้มเพราะ container ถูกปิดหลัง test แรก แก้โดยย้ายการปิดไป TestMain

### สิ่งที่ AI เกือบต้องเดา / ทีมต้องรู้
- ชื่อไฟล์ migration (ดูด้านบน)
- go directive: testcontainers v0.45 และ golang.org/x/* ล่าสุดต้องการ Go 1.26 จึงล็อกเวอร์ชันที่ยังใช้ Go 1.25 ได้ (go 1.25.11)
- ข้อมูลตั้งต้นเป็นข้อมูลตัวอย่าง (ชื่อรุ่น ภาควิชาที่ 2 รหัสครุภัณฑ์ CS-NB-xxxx รหัสผู้ใช้ student-001) ไม่ใช่ข้อมูลจริง
- **ยังไม่มี task ไหนสั่งให้เรียก Migrate/Seed ตอนเปิดหลังบ้าน** (main.go ไม่อยู่ในไฟล์ที่แตะของ T-02) ไม่ได้แก้เอง
- **plan ไม่มีที่มาของชื่อผู้ยืม** (loan เก็บแค่ borrower_id แต่ FR-EQS-07 ให้เจ้าหน้าที่เห็นชื่อ) กระทบ T-17 ไม่ได้เพิ่มคอลัมน์เอง

---

## 2569-11-09 ตรวจ 001 เทียบกับของกลางล่าสุด แล้วแก้ 5 จุด (spec v4.1)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ถามว่าไฟล์ใน 001 เป็นฉบับล่าสุดไหมก่อนทำงานต่อ AI ตรวจแล้วเจอ 5 จุด ทีมสั่งให้แก้
- สิ่งที่ AI แก้: spec เพิ่ม DOM-REQ-01 ในใช้ร่วม และอ้าง Q-SH-05, Q-SH-08 / tasks: T-04 รองรับ DOM-REQ-01,
  ตาราง Constraint, ชื่อไฟล์ T-02 (.up.sql), T-21 เลิกอ้าง UC-13, ตาราง Q เพิ่ม Q-SH-08 / plan อ้าง v4.1 / CHANGELOG
- ยังไม่ได้ทำ (อยู่ในงาน spec v4.2 ถัดไป): กลิ่นจากตารางดมกลิ่น และที่มาของชื่อผู้ยืม (จาก T-02)

---

## 2569-11-09 /implement T-03 (001-equipment-status)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (เฟส 0 ฐานราก auth ใช้ร่วมกับ 002-006)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-03 specs/001-equipment-status/tasks.md`

### คำถามที่ AI ถาม และคำตอบของทีม
1. ผู้ใช้ยืนยันตัวตนได้ แต่ไม่มีภาควิชา (หรือบทบาทไม่รู้จัก) ตอบ HTTP code อะไร (plan 6.1 บอกแค่ "ปฏิเสธ")
   -> ทีมเลือก **403** (ตัวเลือกที่ AI แนะนำ) เพราะ 401 จะทำให้หน้าจอส่งกลับหน้าล็อกอินซ้ำไม่จบ (AC-EQS-09)

### ไฟล์ที่สร้าง (ตรงกับช่องไฟล์ที่แตะทั้งหมด)
- `backend/internal/auth/principal.go`, `middleware.go`, `mock_provider.go`, `middleware_test.go`

### ผล test (`cd backend && go test ./... -v`)
- TestAuth_NoToken_401 PASS (ไม่มี header / ไม่ใช่ Bearer / token ว่าง / token ไม่รู้จัก)
- TestAuth_MissingDepartment_Rejected PASS (ไม่มีภาควิชา / บทบาทไม่รู้จัก / ไม่มีรหัสผู้ใช้ -> 403)
- TestAuth_ValidToken_PassesPrincipal PASS (test เพิ่ม: ผู้ใช้ครบผ่านถึง handler พร้อมข้อมูลใน context)
- test ของ T-01, T-02 ยัง PASS

### สิ่งที่ AI เลือกเอง (ทีมต้องรู้)
- รูปแบบ token ของ auth จำลอง: header `Authorization: Bearer mock-<รหัสผู้ใช้>` (ไม่ใช่โปรโตคอลจริง รอ Q-SH-03)
- ผู้ใช้จำลอง 4 คน: student-001, student-002 (ภาควิชา 1), student-101 (ภาควิชา 2), staff-001 (ภาควิชา 1) ตรงกับ seed.sql
- error code `PRINCIPAL_INCOMPLETE` สำหรับ 403 กรณีข้อมูลผู้ใช้ไม่ครบ (plan ข้อ 5 ยังไม่มี code นี้ ควรเพิ่มในตาราง Error contract ของ plan)
- ยังไม่ได้ผูก middleware เข้ากับ router ใน main.go (main.go ไม่อยู่ในไฟล์ที่แตะ) จะผูกตอนมี endpoint แรกใน T-04

---

## 2569-11-10 ตอบคำขอของ 003 เรื่อง audit_log (แก้ของกลาง เพราะ 003)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: ของกลาง (DM-07) และ 001 (plan, tasks) แก้เพราะ 003-cancel-request
- เครื่องมือ: Claude Code
- ที่มา: นภัสนัน (003) ส่งหน้าสรุป "คุยกับใครก่อนทำ 003" ถาม 001 3 เรื่อง: (1) DM-07 ไม่มีช่องรหัสคำขอ (Q-13 ของ 003)
  (2) ขอเพิ่ม action cancel และ cancel_notify_failed (3) โครงกลางและ migration ของ audit_log จะ merge เมื่อไร
- ทีมตอบ (Pannathorn J ตัดสินตามที่ AI เสนอ): (1) เพิ่มคอลัมน์ request_id ว่างได้ (2) เพิ่ม cancel ได้ cancel_notify_failed รอทีม
  (3) โครงกลาง merge แล้ว (PR #2) ส่วนตาราง audit_log แยกเป็น T-31 ทำก่อน
- สิ่งที่ AI แก้: data-model.md v4 (DM-07 + ตารางเลข migration ที่เสนอ), plan 001 v1.4 ข้อ 4, tasks 001 (เพิ่ม T-31, T-19 ต้องทำหลัง T-31), CHANGELOG
- branch: 001-audit-log-table (แตกจาก main หลัง PR #2)

---

## 2569-11-10 /implement T-31 (001-equipment-status, ตาราง audit_log ใช้ร่วม)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (ของกลาง ทำก่อนเพราะ 003-cancel-request ต้องใช้)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-31 specs/001-equipment-status/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `backend/db/migrations/0002_audit_log.up.sql` (ใหม่)
- `backend/internal/platform/migrate_test.go` (เพิ่ม 2 test และ helper execInTx)

### ผล test (`cd backend && go test ./... -v`)
- TestMigration_AuditLogAppendOnly PASS (INSERT ได้ / UPDATE, DELETE, TRUNCATE ถูกปฏิเสธ)
- TestMigration_AuditLogAction PASS (list, search, detail, cancel ได้ / cancel_notify_failed, approve, ค่าว่าง ถูกปฏิเสธ / request_id และ asset_code ว่างได้)
- test เดิมของ T-01, T-02 ยัง PASS
- ลองกับฐานข้อมูลจริง: UPDATE และ DELETE ได้ error "audit_log is append-only" จาก trigger จริง

### สิ่งที่ AI เลือกเอง (ทีมต้องรู้)
- append-only ใช้ trigger แทนการจำกัดสิทธิ์ผู้ใช้ฐานข้อมูล (plan เขียนว่า "ไม่ให้ app user ทำ UPDATE/DELETE")
  เพราะเครื่องพัฒนาใช้ผู้ใช้คนเดียวที่มีสิทธิ์ทุกอย่าง จำกัดสิทธิ์แล้วทดสอบไม่ได้ trigger บังคับได้ทุกผู้ใช้
- request_id เป็น INT ให้ตรงกับ id แบบ SERIAL ของตารางอื่น ถ้า 002 ใช้ชนิดอื่นใน borrow_request ต้องแก้ด้วย migration ใหม่
- ไม่ใส่ CHECK ที่ actor_role และไม่ใส่ index เพิ่ม เพราะ plan ไม่ได้สั่ง

---

## 2569-11-10 spec 001 v4.2 (ตรวจคำกำกวมรอบ 2) และ plan v1.5

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (+ Q-SH-09 ในของกลาง กระทบ 006)
- เครื่องมือ: Claude Code
- ที่มา: ตารางดมกลิ่นบน Miro (ทีมตรวจแล้ว) + สิ่งที่ AI รายงานตอนทำ T-02, T-03, T-31
- สิ่งที่ AI ทำ: เพิ่ม Q-13 ถึง Q-20 ใน spec, Q-SH-09 ใน constraints.md, แก้ ASM-03 และ CON-DATA-02 ไม่ให้ตอบคำถามที่ยังเปิดไปก่อน,
  plan v1.5 (Error contract 403, append-only ด้วย trigger, main.go ใน T-04, ชื่อผู้ยืมรอ Q-20), CHANGELOG
- สิ่งที่ AI ไม่ได้ทำ: ไม่แก้ FR หรือ AC ข้อไหน ไม่เดาคำตอบของ Q ใหม่ ไม่แก้ tasks.md (รอทีมรัน /tasks เอง)
  ยังไม่ได้แก้ Use Case Spec ของ UC-01 บน Miro ให้มีกฎ CR-001 (ทีมต้องแก้เอง)

---

## 2569-11-10 Q-16 ย้อนกลับเป็นคำถามเปิด (001)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- เหตุผล: คำตอบ Q-16 (นักศึกษาเลือกชิ้นเองไม่ได้ เจ้าหน้าที่จัดให้) ขัดกับ FR-EQS-06 ที่ตกลงไว้ ถ้ารับต้องทำ CR และกระทบ 002, 005
  AI เสนอ 3 ทาง (ก. CR รับ / ข. CR ไม่รับ / ค. ย้อนเป็นคำถามเปิด) ทีมเลือก ค.
- สิ่งที่ AI แก้: Q-16 ใน spec เปลี่ยนเป็น [ยังเปิด] เก็บคำตอบเบื้องต้นไว้พร้อมหมายเหตุว่ายังไม่ยืนยัน และระหว่างนี้ยึด FR-EQS-06 เดิม (รายชิ้น)
- ผลต่อ tasks: ตอนรัน /tasks ใหม่ task ที่เกี่ยวกับปุ่มส่งคำขอยืม (T-06, T-11, T-16) ควรเป็น "รอ Q-16"

---

## 2569-11-10 /tasks specs/001-equipment-status/spec.md (tasks v2)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/tasks specs/001-equipment-status/spec.md`
- อ่านก่อนทำ: spec.md (Status เขียน Draft v4.2 แต่ commit 337031b ตั้งชื่อว่า v4.3 มีคำตอบ Q-13 ถึง Q-15 แล้ว), plan.md v1.5, AGENTS.md, specs/000-shared/*.md, tasks.md v1
- ผลลัพธ์: เขียน tasks.md ใหม่เป็น v2
  - 14 เฟสที่ใช้งาน (เฟส 7 ยกเลิกโดย CR-002 / เฟส 14 ย้ายไปเฟสรอคำตอบ) รวม 36 task: เสร็จ 4 (T-01, 02, 03, 31), พร้อมทำ 19, รอ Q-xx 11, ยกเลิก 2 (T-15, T-25)
  - คงเลข task เดิมทั้งหมด เพราะมี commit อ้างอยู่ task ใหม่เริ่มที่ T-32
  - เพิ่ม T-32, T-33 (ตัวเลขตามตัวกรองและเลือกได้หลายค่า AC-EQS-15, 16 จากคำตอบ Q-15)
  - แยกเรื่องปุ่มส่งคำขอยืมออกจาก T-10, T-11, T-14, T-16 ไปเป็น T-34 (รอ Q-16)
  - เพิ่ม T-35 รีเฟรชเอง (รอ Q-18), T-36 ชื่อผู้ยืมของเจ้าหน้าที่ (รอ Q-20)
  - ย้าย T-22 load test ไปเฟสรอคำตอบ (รอ Q-17)
  - T-30 เหลือแค่ auth จริง (รอ Q-SH-03) เพราะ Q-P2 ปิดด้วย Q-15 แล้ว
  - T-04 เพิ่ม `cmd/server/main.go` ในไฟล์ที่แตะ (ตาม plan ข้อ 10)
- สิ่งที่ AI ไม่ได้ทำ: ไม่เริ่มทำ task ใด ๆ ไม่แก้ spec ไม่ commit
- ข้อค้นพบที่ทีมต้องตัดสิน:
  - บรรทัด Status ของ spec ยังเขียน Draft v4.2 ทั้งที่ commit ตั้งชื่อว่า v4.3 / ตาราง Traceability ยังไม่มี AC-EQS-15, 16
  - AC-EQS-15: ตอนกรองสถานะอื่นที่ไม่ใช่ "ว่าง" ข้อความกำกับตัวเลขบนแถวรุ่นเป็นแบบไหน spec ไม่บอก (เสนอให้เพิ่มเป็น Q ใหม่)
  - FR-EQS-07 ส่วนสรุปจำนวนแยกตามสถานะไม่มี AC / NFR-USE-01, NFR-AVAIL-01 ไม่มี AC
  - CON-TECH-02 ส่วน HTTPS ไม่มี task ในฟีเจอร์นี้ เป็นงานของทั้งระบบ ต้องตกลงว่าใครทำ

---

## 2569-11-10 เก็บ spec 001 v4.3 (Clarifications, Traceability, Q-21) และแก้ tasks v2

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: ให้ AI อธิบายสิ่งที่ต้องทำในข้อ 1 ของหน้าสถานะงาน แล้วตรวจสิ่งที่ทีมแก้ แล้วให้แก้ prompt-log และ tasks.md
- สิ่งที่ AI ร่าง แล้วทีมนำไปแก้เองใน spec.md และ CHANGELOG:
  - Clarifications Session 2569-11-10 (v4.3): Q/A ของ Q-13, Q-14, Q-15 คัดจากคำตอบที่บันทึกใน Open Questions ไม่เพิ่มรายละเอียด
    และหมายเหตุว่า Q-16 ยังเปิด
  - Traceability: เพิ่ม AC-EQS-15 ในแถว FR-EQS-01 และ FR-EQS-02, เพิ่ม AC-EQS-16 ในแถว FR-EQS-02, เพิ่มที่มา Q-13, Q-14, Q-15
  - Q-21 (ข้อความกำกับตัวเลขเมื่อกรองสถานะที่ไม่ใช่ "ว่าง") AI เสนอเลข ทีมตกลงใช้เลขนี้
  - CHANGELOG รายการ Draft v4.3 + tasks v2
- สิ่งที่ AI ตรวจพบ (ทีมตัดสินเอง):
  - CHANGELOG ยังมีคำว่า "(ถ้าทีมตกลงเพิ่ม)" ที่ Q-21
  - AC-EQS-15 ตรวจว่าตัวเลขนับเฉพาะชิ้นที่ตรง แต่ FR-EQS-01, 02 ยังไม่เขียนเรื่องนี้ เสนอเติมท้าย FR-EQS-02 จากคำตอบ Q-15
  - หมายเหตุ Updated ที่หัว spec ยังเป็นของรอบ v4.2
- สิ่งที่ AI แก้เอง:
  - tasks.md: หัวไฟล์อ้าง spec Draft v4.3, T-33 และจุดตรวจเฟส 2 อ้าง Q-21 แทน "spec ไม่บอก",
    ตาราง "สิ่งที่ยังไม่ทำ" เพิ่มแถว Q-21 -> T-33 และเอาเรื่องนี้ออกจากรายการ "ยังไม่มี Q"
  - T-33 ยังเป็น "พร้อมทำ" เพราะส่วนนับตัวเลขทำได้ เหลือแค่ข้อความกำกับที่รอ Q-21 (ไม่ได้เลือกข้อความแทนทีม)
- สิ่งที่ AI ไม่ได้ทำ: ไม่แก้ spec.md หรือ CHANGELOG ไม่ commit

---

## 2569-11-10 /implement-all เฟส 1: T-04 สร้าง API รายการรุ่น และผูกระบบใน main.go

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 1`
- ไฟล์ที่สร้าง: `backend/internal/equipment/dto.go`, `repository.go`, `service.go`, `handler.go`, `handler_test.go`
- ไฟล์ที่แก้: `backend/cmd/server/main.go` (รัน Migrate, ใส่ข้อมูลตั้งต้นเมื่อฐานข้อมูลว่าง, ผูก auth.Middleware)
- test: `TestAC_EQS_01_ModelRowCounts` ผ่าน, `TestModels_RequiresAuthentication` ผ่าน, test เดิมทั้งหมดยังผ่าน
- ลองจริง: เปิดหลังบ้านกับ docker compose แล้วเรียก `/v1/equipment/models` ด้วย `mock-student-001` ได้ รุ่น A ว่าง 3 / ทั้งหมด 5
  ไม่มี token ได้ 401 / เปิดหลังบ้านครั้งที่ 2 ไม่ใส่ข้อมูลตั้งต้นซ้ำ
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - path ที่หลังบ้านรับคือ `/v1/equipment/...` ไม่ใช่ `/api/v1/...` เพราะ proxy ใน `frontend/vite.config.js` ตัด `/api` ออก
    หน้าจอเรียก `/api/v1/equipment/models` ตาม plan ข้อ 5 ได้ตรงตัวผ่าน Vite (port 5173) ส่วน curl ไปที่ port 8000 ตรง ๆ ต้องไม่มี `/api`
  - ใส่ข้อมูลตั้งต้นเฉพาะตอนตาราง department ว่าง (seed ใช้ id ตายตัว ใส่ซ้ำไม่ได้)
    ถ้าฐานข้อมูลเก็บไว้หลายวัน กำหนดคืนของ CS-NB-0004 (16:00 ของวันที่ seed) จะเลยกำหนดไปด้วย ล้างด้วย `docker compose -f deploy/docker-compose.yml down -v`
  - อ่านฐานข้อมูลไม่สำเร็จตอบ 503 `UPSTREAM_UNAVAILABLE` ตาม Error contract ใน plan (test ของกรณีนี้อยู่ใน T-20)
- สิ่งที่ยังไม่ทำ (เป็นของ task อื่น): ยังไม่จำกัดภาควิชา (T-18) ตอนนี้นักศึกษาเห็นรุ่น C ของภาควิชา 2 ด้วย, ยังไม่แบ่งหน้า (T-23 รอ Q-05)

---

## 2569-11-10 /implement-all เฟส 1: T-05 สร้าง API รายชิ้นของรุ่น

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 1` (task ที่ 2 ของรอบนี้)
- ไฟล์ที่แก้: `backend/internal/equipment/repository.go` (ListItems, ItemRow), `handler.go` (GET /v1/equipment/models/{modelId}/items), `handler_test.go`
- test: `TestAC_EQS_01_ItemsOfModel` ผ่าน (5 ชิ้น CS-NB-0001 ถึง 0005 พร้อมสถานะ), `TestItems_UnknownModel` ผ่าน, test เดิมทั้งหมดยังผ่าน
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - T-05 ไม่มี service.go และ dto.go ในช่อง "ไฟล์ที่แตะ" จึงวาง ItemRow ไว้ใน repository.go และ handler เรียก repository ตรง
    task ที่แตะ dto.go และ service.go (T-10, T-12) ย้ายไปผ่าน service ได้ตอนเพิ่มกำหนดคืน
  - รุ่นที่ไม่มีอยู่ได้ 200 กับรายการว่าง (ตาม plan "ไม่พบ = items ว่าง") รหัสรุ่นที่ไม่ใช่ตัวเลขได้ 404 แบบ path ที่ไม่มี
    plan ไม่ได้กำหนดสองกรณีนี้ ไม่ได้สร้าง error code ใหม่

---

## 2569-11-10 /implement-all เฟส 1: T-06 สร้างหน้ารายการอุปกรณ์และรายชิ้น

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 1` (task ที่ 3 ของรอบนี้)
- คำถามที่ AI ถามก่อนทำ และคำตอบของทีม (Pannathorn J):
  1. T-06 ไม่มี `frontend/src/App.jsx` ในช่องไฟล์ที่แตะ แต่ไม่แก้แล้วหน้ารายการจะไม่ขึ้น -> ทีมตอบ: แก้ App.jsx ได้
  2. ยังไม่มีหน้าล็อกอินจำลอง (T-21) หน้าจอจะเอา token จากไหน -> ทีมตอบ: ตั้งค่าใน `frontend/src/config.js`
     ค่าเริ่มต้น `mock-student-001` เปลี่ยนได้ด้วย `VITE_MOCK_TOKEN` (ทีมเลือกเอง ไม่ได้มาจาก spec)
  เพิ่ม 2 ไฟล์นี้ในช่องไฟล์ที่แตะของ T-06 ใน tasks.md แล้ว
- ไฟล์ที่สร้าง: `frontend/src/pages/EquipmentList/EquipmentList.jsx`, `frontend/src/components/ModelRow.jsx`, `ItemTable.jsx`,
  `StatusBadge.jsx`, `frontend/src/config.js`, `frontend/src/__tests__/AC-EQS-01.test.jsx`
- ไฟล์ที่แก้: `frontend/src/api/client.js` (listModels, listItems, ส่ง token), `frontend/src/App.jsx` (แสดง EquipmentList)
- test: `AC-EQS-01.test.jsx` ผ่าน, `setup.test.jsx` ยังผ่าน (มีคำเตือน act จาก React เพราะ App เรียก API จริงใน test นั้น ไม่ทำให้ test ล้ม)
  test หลังบ้านทั้งหมดยังผ่าน, `vite build` ผ่าน
- ลองจริง: เปิดหลังบ้านและ Vite แล้วเรียก `/api/v1/equipment/models/1/items` ผ่าน proxy ที่ port 5173 ได้ 5 ชิ้น
  ยังไม่ได้เปิดดูหน้าจอในเบราว์เซอร์ (เครื่องไม่มีเบราว์เซอร์) ทีมต้องลองเองตามจุดตรวจเฟส 1
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - ถ้าโหลดข้อมูลไม่สำเร็จ แสดงข้อความ "ดึงข้อมูลอุปกรณ์ไม่สำเร็จ" และไม่แสดงรายการเก่า (FR-EQS-08) ปุ่มลองใหม่และ test อยู่ใน T-20
  - ป้ายสถานะใช้คำใน glossary (ว่าง / ถูกยืม / ไม่พร้อมใช้งาน) ไม่กำหนดสีแยก เพราะ T-25 เรื่องสียกเลิกแล้ว
  - ยังไม่มีปุ่มส่งคำขอยืม (T-34 รอ Q-16)

---

## 2569-11-10 /implement-all เฟส 2: T-07 เพิ่มตัวกรองประเภทและสถานะใน API

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 2` (task ที่ 1 ของรอบนี้)
- คำถามที่ AI ถามก่อนเริ่มเฟส 2 และคำตอบของทีม (Pannathorn J):
  1. T-07, T-32 ไม่มี service.go ในไฟล์ที่แตะ แต่ตัวกรองต้องผ่าน service -> ทีมตอบ: เพิ่ม service.go ได้ (เพิ่มในช่องไฟล์ที่แตะของ T-07 แล้ว)
  2. T-08 ช่องค้นหา แต่วิธีจับคู่คำค้นรอ Q-06 -> ทีมตอบ: ยังไม่ทำช่องค้นหา T-08 ทำแค่ตัวกรอง ช่องค้นหาไปทำกับ T-24
  3. T-33 ข้อความกำกับตัวเลขเมื่อกรอง รอ Q-21 -> ทีมตอบ: คงรูปแบบเดิม "ว่าง x / ทั้งหมด y" นับเฉพาะชิ้นที่ตรงตัวกรอง (ชั่วคราว จนกว่าจะได้คำตอบ Q-21)
- ไฟล์ที่แก้: `backend/internal/equipment/repository.go` (Filter, เงื่อนไขกรองที่ใช้ร่วม), `service.go` (ส่งตัวกรอง, ListItems), `handler.go` (อ่าน ?category=&status=), `handler_test.go`
- test: `TestAC_EQS_02_FilterStatus` ผ่าน, `TestFR_EQS_02_OnlyModelsWithMatchingItems` ผ่าน (รวมค่าที่หน้าตาเป็น SQL ไม่ทำให้ query พัง), test เดิมทั้งหมดยังผ่าน
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - ตัวกรองประเภทใช้ "ชื่อประเภท" ไม่ใช่รหัส เพราะ DTO ของรายการรุ่นมีแค่ชื่อประเภท (plan ข้อ 5) หน้าจอจึงสร้างตัวเลือกจากรายการได้เลย
  - ค่าตัวกรองที่ไม่รู้จักไม่ตรงกับชิ้นใด ได้รายการว่าง 200 ไม่ได้สร้าง error code ใหม่
  - กรองที่รายชิ้นก่อนนับ ตัวเลขของรุ่นจึงนับตามตัวกรองตั้งแต่ task นี้ (FR-EQS-02 v4.3) T-32 เพิ่มการเลือกหลายค่าและ test ของ AC-EQS-15, 16

---

## 2569-11-10 /implement-all เฟส 2: T-32 ตัวกรองเลือกหลายค่า และตัวเลขของรุ่นนับตามตัวกรองใน API

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 2` (task ที่ 2 ของรอบนี้)
- ไฟล์ที่แก้: `backend/internal/equipment/handler.go` (อ่านค่าซ้ำของ category และ status), `handler_test.go`
  repository.go ไม่ต้องแก้ เพราะ T-07 ใช้ array ใน SQL รองรับหลายค่าไว้แล้ว
- test: `TestAC_EQS_15_CountsFollowFilter` ผ่าน (กรอง "ถูกยืม" รุ่น 1 แสดง 1 แถว นับ 2 ชิ้น ว่าง 0),
  `TestAC_EQS_16_MultiStatusFilter` ผ่าน (ว่าง + ถูกยืม ไม่มี CS-PJ-0002), test เดิมทั้งหมดยังผ่าน
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้): ส่งหลายค่าด้วยชื่อซ้ำ `?status=available&status=borrowed` (plan ไม่ได้กำหนดรูปแบบ)

---

## 2569-11-10 /implement-all เฟส 2: T-08 เพิ่มตัวกรองบนหน้าจอ

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 2` (task ที่ 3 ของรอบนี้)
- คำถามที่ AI ถามก่อนทำ และคำตอบของทีม (Pannathorn J):
  T-08 ไม่มี `frontend/src/api/client.js` ในไฟล์ที่แตะ แต่ไม่แก้แล้วหน้าจอจริงจะไม่ส่งตัวกรอง -> ทีมตอบ: เพิ่ม client.js ได้ (เพิ่มในช่องไฟล์ที่แตะของ T-08 แล้ว)
- ไฟล์ที่สร้าง: `frontend/src/components/SearchBar.jsx`, `frontend/src/__tests__/AC-EQS-02.test.jsx`
- ไฟล์ที่แก้: `frontend/src/pages/EquipmentList/EquipmentList.jsx` (ตัวกรอง ส่งตัวกรองทั้งรายการรุ่นและรายชิ้น), `frontend/src/api/client.js` (แปลงตัวกรองเป็น query string)
- test: `AC-EQS-02.test.jsx` 2 test ผ่าน (กรองสถานะ "ว่าง", กรองประเภทแล้วกดซ้ำเพื่อยกเลิก), test หน้าจอเดิมและ test หลังบ้านทั้งหมดยังผ่าน, `vite build` ผ่าน
- ตามที่ทีมตัดสินก่อนเริ่มเฟส: ยังไม่มีช่องค้นหาด้วยคำ (รอ Q-06 ทำกับ T-24) ชื่อไฟล์ SearchBar.jsx คงตาม tasks.md
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - ตัวกรองเป็น chip ตาม plan ข้อ 6.7 ตอนนี้เลือกได้ 1 ค่าต่อกลุ่ม กดซ้ำเพื่อยกเลิก (เลือกหลายค่าคือ T-33)
  - ตัวเลือกประเภทสร้างจากรายการตอนยังไม่กรอง เพราะ plan ข้อ 5 ไม่มี endpoint รายการประเภท
  - client.js ส่งหลายค่าได้แล้ว (ชื่อซ้ำ) T-33 จึงไม่ต้องแตะ client.js

---

## 2569-11-10 /implement-all เฟส 2: T-33 หน้าจอเลือกตัวกรองหลายค่า และแสดงตัวเลขตามตัวกรอง

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 2` (task ที่ 4 ของรอบนี้ ครบเฟส)
- ไฟล์ที่สร้าง: `frontend/src/__tests__/AC-EQS-15.test.jsx`, `frontend/src/__tests__/AC-EQS-16.test.jsx`
- ไฟล์ที่แก้: `frontend/src/components/SearchBar.jsx` (กดเพิ่ม/เอาออกทีละค่า), `frontend/src/components/ModelRow.jsx` (คอมเมนต์เรื่องข้อความกำกับตัวเลขรอ Q-21)
- test: `AC-EQS-15.test.jsx` ผ่าน (กรอง "ถูกยืม" รุ่น A 1 แถว "ว่าง 0 / ทั้งหมด 2"), `AC-EQS-16.test.jsx` ผ่าน (ว่าง + ถูกยืม ไม่เห็น CS-PJ-0002)
  test หน้าจอทั้งหมด 5 ไฟล์ 6 test ผ่าน, test หลังบ้านทั้งหมดผ่าน, `vite build` ผ่าน
- ลองจริง: หลังบ้านตัวใหม่ที่ port 8010 (port 8000 มีหลังบ้านที่ทีมเปิดค้างไว้ AI ไม่ได้ปิด)
  `?status=borrowed` ได้รุ่น A แถวเดียว ว่าง 0 / ทั้งหมด 2 / รายชิ้นรุ่น 2 `?status=available&status=borrowed` ได้แค่ CS-PJ-0001
- ตามที่ทีมตัดสินก่อนเริ่มเฟส: ข้อความกำกับตัวเลขคงรูปแบบ "ว่าง x / ทั้งหมด y" ที่นับเฉพาะชิ้นที่ตรงตัวกรอง ชั่วคราวจนกว่าจะได้คำตอบ Q-21

---

## 2569-11-10 /implement-all เฟส 3: T-09 แสดงข้อความไม่พบและปุ่มล้างตัวกรอง

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 3` (task เดียวของเฟส)
- คำถามที่ AI ถามก่อนทำ และคำตอบของทีม (Pannathorn J):
  AC-EQS-03 พูดถึง "ค้นหาด้วยคำที่ไม่มีในระบบ" แต่ยังไม่มีช่องค้นหา (รอ Q-06 ดู T-24) -> ทีมตอบ: ทำด้วยตัวกรองที่ไม่เจอ
  ส่วน "ค้นหาด้วยคำ" ยังขาด รอ T-24
- ไฟล์ที่สร้าง: `frontend/src/components/EmptyState.jsx`, `frontend/src/__tests__/AC-EQS-03.test.jsx`
- ไฟล์ที่แก้: `frontend/src/pages/EquipmentList/EquipmentList.jsx` (แสดง EmptyState เมื่อไม่พบ ล้างตัวกรองแล้วโหลดรายการทั้งหมด), `backend/internal/equipment/handler_test.go`
- test: `TestAC_EQS_03_EmptyIs200` ผ่าน (ไม่พบ = 200 `{items: [], total: 0}`), `AC-EQS-03.test.jsx` ผ่าน
  (แก้ test 1 รอบ: assertion สุดท้ายใช้ getAllByRole ซึ่ง error เมื่อไม่เจอ เปลี่ยนเป็น queryAllByRole ไม่ได้แก้โค้ดของระบบ)
  test หน้าจอทั้งหมด 6 ไฟล์ 7 test และ test หลังบ้านทั้งหมดผ่าน, `vite build` ผ่าน
- ข้อความบนหน้าจอใช้คำตาม FR-EQS-03 / AC-EQS-03: "ไม่พบอุปกรณ์" และปุ่ม "ล้างตัวกรอง"

---

## 2569-11-10 /implement-all เฟส 4: T-10 สร้าง API รายละเอียดรายชิ้น พร้อมกำหนดคืน

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 4` (task ที่ 1 ของรอบนี้)
- คำถามที่ AI ถามก่อนเริ่มเฟส 4 และคำตอบของทีม (Pannathorn J):
  1. T-10 ไม่มี handler.go ในไฟล์ที่แตะ แต่ต้องเพิ่ม endpoint -> ทีมตอบ: เพิ่ม handler.go ได้
  2. ยังไม่มี router จะเปิดหน้ารายละเอียดอย่างไร -> ทีมตอบ: กดรหัสครุภัณฑ์ในตารางรายชิ้น (ไม่เพิ่ม library)
     T-11 จึงแตะ ItemTable.jsx, EquipmentList.jsx, client.js, config.js เพิ่ม
  3. รูปแบบกำหนดคืน -> ทีมตอบ: "2569-10-10 16:00" ตามตัวอย่างใน AC-EQS-14 (ปี พ.ศ. เวลา Asia/Bangkok)
  เพิ่มไฟล์ในช่องไฟล์ที่แตะของ T-10 และ T-11 ใน tasks.md แล้ว
- ไฟล์ที่แก้: `backend/internal/equipment/dto.go` (ItemDetailDTO), `repository.go` (GetItem), `service.go` (กำหนดคืนเฉพาะถูกยืม),
  `handler.go` (GET /v1/equipment/items/{assetCode}), `handler_test.go`
- test: `TestAC_EQS_04_BorrowedShowsDueDate` ผ่าน (จำลอง UC-08 ส่งมอบชิ้นใหม่ กำหนดคืนพรุ่งนี้ ได้ status borrowed และ dueAt ตรงกับที่บันทึก),
  `TestItemDetail_AvailableHasNoDueAtAndUnknownIs404` ผ่าน, `TestCON_DATA_02_NoWriteEndpoints` ผ่าน (POST/PUT/PATCH/DELETE ได้ 405 สถานะไม่เปลี่ยน)
  test เดิมทั้งหมดยังผ่าน
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - รหัสครุภัณฑ์ที่ไม่มีได้ 404 แบบไม่มี error code (plan ไม่ได้กำหนดกรณีนี้)
  - ยังไม่มี isOverdue และยังไม่ซ่อนกำหนดคืนของนักศึกษาเมื่อเลยกำหนด (เป็นงานของ T-12) ตอนนี้ชิ้นที่ถูกยืมได้ dueAt เสมอ
  - test AC-EQS-04 ใส่ข้อมูลชิ้นใหม่เองแล้วลบทิ้งหลัง test ไม่พึ่งเวลา 16:00 ของข้อมูลตั้งต้น เพื่อไม่ให้ผลขึ้นกับเวลาที่รัน

---

## 2569-11-10 /implement-all เฟส 4: T-11 สร้างหน้ารายละเอียดอุปกรณ์

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 4` (task ที่ 2 ของรอบนี้ ครบเฟส)
- ไฟล์ที่สร้าง: `frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx` (หน้ารายละเอียด + formatDueAt), `frontend/src/__tests__/AC-EQS-04.test.jsx`
- ไฟล์ที่แก้: `frontend/src/components/ItemTable.jsx` (รหัสครุภัณฑ์กดได้), `frontend/src/pages/EquipmentList/EquipmentList.jsx` (เปิด/ปิดรายละเอียด),
  `frontend/src/api/client.js` (getItem), `frontend/src/config.js` (TIME_ZONE = Asia/Bangkok ตาม glossary)
- test: `AC-EQS-04.test.jsx` 4 test ผ่าน (ถูกยืมแสดง "ถูกยืม" + กำหนดคืน "2569-11-10 16:00", ชิ้นที่ว่างไม่มีกำหนดคืน,
  formatDueAt ใช้เวลาไทยแม้ UTC เป็นคนละวัน, กดรหัสในตารางเปิดรายละเอียดแล้วกดกลับ)
  test หน้าจอทั้งหมด 7 ไฟล์ 11 test และ test หลังบ้านทั้งหมดผ่าน, `vite build` ผ่าน
- ลองจริง: หลังบ้านที่ port 8010 `GET /v1/equipment/items/CS-NB-0004` ได้ borrowed + dueAt, CS-NB-0001 ไม่มี dueAt, POST ได้ 405
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - ปุ่ม "กลับไปรายการอุปกรณ์" เป็นการนำทางกลับ ไม่ใช่ฟีเจอร์ใหม่ (มาจากการตัดสินใจเปิดรายละเอียดในหน้าเดียวกัน) กลับแล้วตัวกรองและรุ่นที่เลือกยังอยู่
  - รายละเอียดแสดงรหัสครุภัณฑ์ สถานะ และกำหนดคืน ตามฟิลด์ใน plan ข้อ 5 ไม่ได้เพิ่มชื่อรุ่น
  - ยังไม่มีปุ่มส่งคำขอยืม (T-34 รอ Q-16) และป้ายเลยกำหนดคืน (T-13)

---

## 2569-11-10 /implement-all เฟส 5: T-12 คำนวณ isOverdue ด้วยวันที่และเวลา Asia/Bangkok และแยก DTO นักศึกษา/เจ้าหน้าที่

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 5` (task ที่ 1 ของรอบนี้)
- ไฟล์ที่แก้: `backend/internal/config/config.go` (TimeZone = Asia/Bangkok, ฝัง tzdata), `backend/internal/equipment/dto.go` (StudentItemDTO, StaffItemDTO),
  `backend/internal/equipment/service.go` (นาฬิกา, เขตเวลา, กฎเลยกำหนดคืน, เลือก DTO ตามบทบาท)
- ไฟล์ที่สร้าง: `backend/internal/equipment/service_test.go`
- test: `TestAC_EQS_14_OverdueStudentVsStaff` ผ่าน (กำหนดคืน 2569-10-10 16:00 นาฬิกา 17:00: นักศึกษาได้ isOverdue ไม่มี dueAt ใน JSON,
  เจ้าหน้าที่ได้ทั้งคู่ สถานะยังเป็น borrowed), `TestFR_EQS_04_NotOverdueBeforeDueTime` ผ่าน (15:59 และ 16:00 พอดียังไม่เลย),
  `TestFR_EQS_04_OverdueUsesBangkokDateTime` ผ่าน, test เดิมทั้งหมดยังผ่าน
- ทำภายในไฟล์ที่กำหนดได้โดยไม่แก้ handler.go และ main.go:
  - service อ่านบทบาทจาก context ที่ auth.Middleware ใส่ไว้ (auth.FromContext) ไม่ต้องส่งผ่าน handler
  - นาฬิกาเป็นฟิลด์ของ Service (ค่าเริ่มต้น time.Now) test ใน package เดียวกันตั้งเวลาเองได้ NewService ไม่เปลี่ยน signature
  - GetItem คืน `any` (StudentItemDTO หรือ StaffItemDTO) handler ส่งเป็น JSON ได้เหมือนเดิม
  - `ItemDetailDTO` คงไว้เป็นชื่อเดิม (alias ของ StaffItemDTO) ให้ test ของ T-10 ที่อ่าน JSON ยัง compile ได้
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - "เลยกำหนด" = เวลาปัจจุบันมากกว่ากำหนดคืน ตรง 16:00 พอดียังไม่เลย (FR-EQS-04 ใช้เครื่องหมาย >)
  - ไม่มีผู้ใช้ใน context ได้ข้อมูลแบบนักศึกษา (สิทธิ์น้อยที่สุด) ปกติไม่เกิดเพราะทุก endpoint อยู่หลัง middleware
  - isOverdue มีเฉพาะหน้ารายละเอียด รายการรุ่นและตารางรายชิ้นยังไม่แสดงป้าย (spec พูดถึง "เปิดดูรายละเอียด" เท่านั้น)

---

## 2569-11-10 /implement-all เฟส 5: T-13 แสดงป้าย "เลยกำหนดคืน"

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 5` (task ที่ 2 ของรอบนี้ ครบเฟส)
- ไฟล์ที่แก้: `frontend/src/components/StatusBadge.jsx` (OverdueBadge), `frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx` (แสดงป้ายเมื่อ isOverdue)
- ไฟล์ที่สร้าง: `frontend/src/__tests__/AC-EQS-14.test.jsx`
- test: `AC-EQS-14.test.jsx` 3 test ผ่าน (นักศึกษาเห็น "ถูกยืม" + "เลยกำหนดคืน" ไม่เห็นกำหนดคืน / เจ้าหน้าที่เห็น "2569-10-10 16:00" + ป้าย /
  ยังไม่เลยกำหนดไม่มีป้าย) test หน้าจอทั้งหมด 8 ไฟล์ 14 test และ test หลังบ้านทั้งหมดผ่าน, `vite build` ผ่าน
- ลองจริง: หลังบ้านที่ port 8010 CS-NB-0005 นักศึกษาได้ `{status: borrowed, isOverdue: true}` ไม่มี dueAt / เจ้าหน้าที่ได้ dueAt + isOverdue
  CS-NB-0004 (ยังไม่เลยกำหนด) ทั้งสองบทบาทได้ dueAt ไม่มีป้าย
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - ป้าย "เลยกำหนดคืน" ใช้กรอบสีแดงให้แยกจากป้ายสถานะ เป็นเรื่องหน้าตา (ยืดหยุ่นได้) ไม่ได้มาจาก spec
  - หน้าจอไม่คำนวณเลยกำหนดเอง ใช้ isOverdue จากหลังบ้านอย่างเดียว (plan ข้อ 5: คำนวณที่ server)

---

## 2569-11-10 /implement-all เฟส 6: T-14 ป้าย "ไม่พร้อมใช้งาน"

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 6` (task เดียวของเฟส)
- ไฟล์ที่แก้: `frontend/src/components/StatusBadge.jsx` (คอมเมนต์อ้าง FR-EQS-05 เท่านั้น)
- ไฟล์ที่สร้าง: `frontend/src/__tests__/AC-EQS-06.test.jsx`
- test: `AC-EQS-06.test.jsx` 2 test ผ่าน (หน้ารายละเอียดของชิ้นที่ unavailable แสดง "ไม่พร้อมใช้งาน" ไม่มีกำหนดคืนหรือป้ายเลยกำหนด /
  ตารางรายชิ้นแสดงป้ายเดียวกัน) test หน้าจอทั้งหมด 9 ไฟล์ 16 test และ test หลังบ้านทั้งหมดผ่าน, `vite build` ผ่าน
- ข้อสังเกต: StatusBadge แปลง unavailable เป็น "ไม่พร้อมใช้งาน" ไว้แล้วตั้งแต่ T-06 (ป้าย 3 สถานะตาม CON-DATA-01)
  task นี้จึงไม่ต้องเขียนโค้ดหน้าจอเพิ่ม เหลือแค่ test ที่ยืนยัน AC-EQS-06 ส่วน "ไม่มีปุ่มส่งคำขอยืม" รอ Q-16 (T-34)

---

## 2569-11-10 /implement-all เฟส 8: T-16 ไม่ cache ข้อมูลสถานะ ให้รีเฟรชแล้วเห็นสถานะล่าสุด

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ของกลาง `backend/internal/platform/http.go` เพราะฟีเจอร์ 001 ตาม AGENTS.md)
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 8` (task เดียวของเฟส)
- คำถามที่ AI ถามก่อนทำ และคำตอบของทีม (Pannathorn J):
  T-16 ใส่ no-store ใน platform/http.go แต่ router ผูกใน main.go ซึ่งไม่อยู่ในไฟล์ที่แตะ -> ทีมตอบ: เพิ่ม main.go ได้ (เพิ่มในช่องไฟล์ที่แตะแล้ว)
- ไฟล์ที่แก้:
  - `backend/internal/platform/http.go` (ของกลาง): เพิ่ม `NoStore` ใส่ `Cache-Control: no-store` ทุก response
  - `backend/cmd/server/main.go`: ห่อ router ทั้งหมดด้วย `platform.NoStore` (มีผลกับ endpoint ของฟีเจอร์อื่นที่จะมาต่อด้วย)
  - `backend/internal/equipment/handler_test.go`: newTestServer ห่อด้วย NoStore แบบเดียวกับ main.go และเพิ่ม test ของ AC-EQS-05
  - `frontend/src/api/client.js`: fetch ด้วย `cache: 'no-store'`
- ไฟล์ที่สร้าง: `frontend/src/__tests__/AC-EQS-05.test.jsx`
- test: `TestAC_EQS_05_ReturnedBecomesAvailable` ผ่าน (จำลอง UC-08 แล้ว UC-09 รับคืน: เรียกใหม่ได้ available ไม่มี dueAt, ตัวเลขว่างรุ่น 1 เพิ่ม 1,
  header no-store ทุก endpoint), `AC-EQS-05.test.jsx` ผ่าน (รับคืนแล้วเปิดหน้าใหม่เห็น "ว่าง" fetch ใหม่ทุกครั้งด้วย cache no-store)
  test หลังบ้านทั้งหมดและ test หน้าจอ 10 ไฟล์ 17 test ผ่าน, `vite build` ผ่าน
- แก้ระหว่างทาง 1 รอบ: test หลังบ้านส่ง SQL 2 คำสั่งพร้อม parameter ในครั้งเดียว pgx ไม่รับ แยกเป็น 2 คำสั่ง (แก้ใน test ไม่ได้แก้โค้ดของระบบ)
  และ AI เผลอเพิ่ม test ใน `platform/http_test.go` ซึ่งไม่อยู่ในไฟล์ที่แตะ จึงย้อนไฟล์นั้นกลับ (test ของ AC-EQS-05 ตรวจ header ครบแล้ว)
- ลองจริง: หลังบ้านที่ port 8010 ทั้ง 200 และ 401 มี `Cache-Control: no-store`

---

## 2569-11-10 เติม FR ของ 002-borrow-request ใน SRS (หัวข้อ 3.1.2)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: อ่าน AGENTS.md, specs/000-shared/ ทั้ง 3 ไฟล์ และ SRS หัวข้อ 3.0 กับ 002-borrow-request
  แล้วเขียน FR ของ 002 ใน SRS (เดิมเป็น `FR-BRQ-01 | [เติม]`)

### คำถามที่ AI ถาม และคำตอบของทีม
1. บันทึก prompt-log ช่อง "ผู้สั่ง" เป็นใคร -> ธัญลักษณ์ (@tnyl02) เจ้าของฟีเจอร์ 002

### สิ่งที่ AI ทำ
- `docs/srs/SRS.md` 3.1.2: เขียน FR-BRQ-01 ถึง FR-BRQ-08 พร้อมช่องที่มาทุกข้อ
  (US-002, RQ-18, RQ-20, RQ-22, RQ-25, FR-EQS-06, DOM-REQ-01, DOM-ACCESS-01, DOM-AUDIT-01, DM-06, DM-07)
- คงไว้เป็น `[เติม]` 2 ข้อที่ยังไม่มีข้อมูลใน Miro และ repo: FR-BRQ-09 (กรณียกเว้นของ UC-02),
  FR-BRQ-10 (ยืมค้างคืน RQ-16 / UC-05 รวมใน 002 หรือแยก)
- เพิ่มบรรทัด "นอกขอบเขตของ 002" และรายการคำกำกวมที่ต้องถามใต้ตาราง
- `docs/srs/SRS.md` บรรทัดหัว: Version 0.3 -> 0.4, Updated 2569-11-10
- `specs/CHANGELOG.md`: เพิ่มรายการ [SRS v0.4 FR-BRQ ของ 002] ไม่ทำ CR เพราะเป็นการเติมช่อง `[เติม]`
  ไม่ได้แก้ FR / AC / Constraint เดิมของฟีเจอร์อื่น

### สิ่งที่ AI ตีความเอง (ทีมต้องตรวจ)
- FR-BRQ-01 ข้อ "ตรวจสถานะอุปกรณ์ซ้ำตอนเปิดฟอร์มและตอนกดส่ง" มาจาก `specs/001-equipment-status/plan.md` ข้อ 5
  ที่ระบุว่าเป็น dependency ของ UC-02 ไม่ใช่ข้อความจาก RQ ใด ๆ โดยตรง
- FR-BRQ-08 รวม audit log ของการส่งคำขอ อ้าง DOM-AUDIT-01 (ซึ่งเขียนว่า "เข้าถึงข้อมูลอุปกรณ์") และ RQ-25
  ถ้าทีมเห็นว่าการส่งคำขอไม่เข้าข่าย DOM-AUDIT-01 ให้ตัดเหลือ RQ-25 อย่างเดียว
- FR-BRQ-06 (RQ-20) เขียนไว้เพราะ SRS 3.1.2 ระบุ RQ-20 เป็นที่มาของฟีเจอร์นี้ แต่ BL-23 เป็น Could
  ต้องยืนยันว่าอยู่ใน prototype หรือไม่
- ยังไม่ได้สร้าง `specs/002-borrow-request/` และยังไม่ได้เพิ่ม 002 ในตาราง "ใครอ้างใช้บ้าง" ของ constraints.md

---

## 2569-11-11 ปิด [เติม] ของ FR-BRQ-09 และ FR-BRQ-10 (SRS 3.1.2)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: อ่าน prompt-log ดูว่างานถึงไหนแล้ว จากนั้นส่งภาพ Use Case Specification UC-02 จาก Miro
  มาให้ใช้เติม FR-BRQ-09 และตอบว่า "ยืมค้างคืนรวมใน 002"

### คำตอบของทีมในรอบนี้
1. FR-BRQ-09 -> ใช้ Exception Flow ของ UC-02 ตามภาพที่ทีมส่งมา
2. FR-BRQ-10 -> การยืมแบบค้างคืน (RQ-16 / UC-05) **รวมอยู่ในฟีเจอร์ 002** ไม่แยกเป็นฟีเจอร์ใหม่

### สิ่งที่ AI ทำ
- `docs/srs/SRS.md` 3.1.2: เติม FR-BRQ-09 (บันทึกไม่สำเร็จ 10a + กลับไปแก้ฟอร์มก่อนส่ง 8a)
  และ FR-BRQ-10 (เลือกประเภทปกติ/ค้างคืน, ค้างคืนต้องมีเหตุผล + ใบบันทึกข้อความ, ปกติต้องแสดงเงื่อนไข DOM-LOAN-01)
- FR-BRQ-04: ระบุสถานะเริ่มต้นของคำขอเป็น "รออนุมัติ" ตาม UC-02 Postcondition (ปิดคำกำกวมเดิม 1 ข้อ)
- ปรับบรรทัดที่มา / สิ่งที่รู้แล้ว / นอกขอบเขต / คำกำกวม ของ 3.1.2, หัวเอกสาร v0.4 -> v0.5, `specs/CHANGELOG.md`
- ไม่ทำ CR เพราะเป็นการเติมช่อง `[เติม]` ไม่ได้แก้ข้อกำหนดเดิมของฟีเจอร์อื่น

### สิ่งที่ AI ไม่ได้เขียนเอง (ตามกฎข้อ 3 และ 5) ทีมต้องตัดสิน
- UC-02 Main Flow 6 "นักศึกษาระบุช่วงเวลาที่ต้องการยืม" ยังไม่มี FR-BRQ ข้อไหนรองรับ
  ถ้าต้องมี ต้องเป็น ID ใหม่ (เสนอ FR-BRQ-11) -> AI ไม่สร้างเอง เขียนไว้ใต้ตารางว่า "ยังไม่มี FR รองรับ"
- UC-02 เขียนว่า NFR-PERF-01 ครอบคลุม**แบบฟอร์มคำขอ** ด้วย แต่ NFR-PERF-01 ปัจจุบัน (RQ-24) พูดถึงเฉพาะหน้ารายการอุปกรณ์
  การขยายต้องทำ CR กับของกลาง -> ยังไม่แก้
- UC-02 ไม่ระบุว่าผู้ขอส่งคำขอของตนเองซ้ำบนอุปกรณ์ชิ้นเดิมได้หรือไม่ -> เพิ่มเป็นคำกำกวมใต้ตาราง ไม่เดา
- ชนิดและขนาดไฟล์เอกสารประกอบของ RQ-16 ไม่มีที่ไหนระบุ -> เพิ่มเป็นคำกำกวม

### สิ่งที่ AI ตีความเอง (ทีมต้องตรวจ)
- รวม UC-02 Main Flow 7 (แสดงเงื่อนไขยืมปกติ คืน 16:00 ใช้ในพื้นที่มหาวิทยาลัย) ไว้ใน FR-BRQ-10
  เพราะอยู่ในขั้นตอนเลือกประเภทการยืมเดียวกัน ถ้าทีมอยากแยกเป็น FR ของตัวเอง ต้องให้ ID ใหม่

### ยังไม่ได้ทำ
- ยังไม่สร้าง `specs/002-borrow-request/spec.md` และยังไม่เพิ่ม 002 ในตาราง "ใครอ้างใช้บ้าง" ของ constraints.md
- งาน SRS ของ 002 (รอบ 2569-11-10 และรอบนี้) ยังไม่ได้ commit

### คำตอบเพิ่มเติมของทีมในรอบเดียวกัน (2569-11-11)
- ถาม: UC-02 Main Flow 6 "ระบุช่วงเวลาที่ต้องการยืม" ยังไม่มี FR รองรับ จะเพิ่ม FR-BRQ-11 ไหม
  -> ทีมตอบ: **ถือว่ารวมอยู่ใน FR-BRQ-03 แล้ว** (ไม่สร้าง ID ใหม่)
- AI แก้ FR-BRQ-03 ให้ระบุว่าข้อมูลที่ต้องครบรวมช่วงเวลาที่ต้องการยืม และตัดข้อนี้ออกจากบรรทัด "ยังไม่มี FR รองรับ"
  เหลือค้างเฉพาะเรื่องขยาย NFR-PERF-01 ให้ครอบคลุมแบบฟอร์มคำขอ (ต้องทำ CR กับของกลาง)
- ทีมสั่งให้**หยุดไว้ก่อนเพื่อตรวจเอง** ยังไม่ commit และยังไม่สร้าง `specs/002-borrow-request/`

---

## 2569-11-11 สร้าง specs/002-borrow-request/spec.md (Draft v1)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (แตะของกลาง constraints.md ตารางใครอ้างใช้บ้าง เพราะเปิดฟีเจอร์ 002)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: สร้าง spec.md ให้ครบ 6 ส่วน + Assumptions & Open Questions, Constraint ใช้ร่วมให้อ้าง ID
  ห้าม copy ข้อความ, ใส่ชื่อเจ้าของ, เพิ่มชื่อฟีเจอร์ในตาราง "ใครอ้างใช้บ้าง" ของ constraints.md
- ก่อนหน้านี้ในวันเดียวกัน: commit `SRS v0.5 FR-BRQ ครบ 10 ข้อ` แล้ว `git pull origin main`
  (conflict 2 ไฟล์ prompt-log.md และ specs/CHANGELOG.md แก้แบบเก็บทั้งสองฝั่ง ไม่ลบบันทึกของใคร)

### ข้อความจากเจ้าของ 001 ที่ทีมส่งต่อมา และสิ่งที่ตรวจพบ
- ข้อความมา 2 บรรทัดที่**ขัดกันเอง**: บรรทัดแรกบอกว่า id ของ borrow_request ต้องเป็น TEXT
  บรรทัดที่สองบอกว่าขอเป็น SERIAL (INT)
- AI เปิดไฟล์จริง `backend/db/migrations/0002_audit_log.up.sql` พบว่า `request_id INT NULL`
  จึงยึด **SERIAL (INT)** ตามบรรทัดที่สอง เพราะตรงกับของจริงในฐานข้อมูล บันทึกไว้ที่ ASM-BRQ-03
- เลข migration ของ 002 = `0003_borrow_request.up.sql` ตามตารางใน data-model.md และ 002 เป็นผู้เพิ่ม FK

### ไฟล์ที่สร้างหรือแก้
- สร้าง `specs/002-borrow-request/spec.md` (Draft v1): FR-BRQ-01 ถึง 10 ใช้ ID เดิมจาก SRS 3.1.2,
  AC-BRQ-01 ถึง 10, ASM-BRQ-01 ถึง 04, Q-BRQ-01 ถึง 08, Traceability
- `specs/000-shared/constraints.md`: เพิ่ม 002-borrow-request ในตาราง "ใครอ้างใช้บ้าง" 10 ID (ไม่แก้ข้อความ constraint)
- `specs/README.md`: สถานะ 002 เป็น spec Draft v1
- `specs/CHANGELOG.md`: รายการ [SPEC-BRQ-002 Draft v1]

### สิ่งที่ AI ไม่ได้ทำ (ตามกฎ)
- ไม่สร้าง FR / NFR / CON ID ใหม่ เรื่องที่ไม่มี ID รองรับเขียนเป็น Q-BRQ แทน
  (ค่า action ของ audit_log = Q-BRQ-07, NFR-PERF-01 ครอบคลุมฟอร์มไหม = Q-BRQ-05)
- ไม่ตัดสิน Q-16 ของ 001 เขียนเป็น ASM-BRQ-02 พร้อมรายการ FR ที่จะกระทบถ้าคำตอบพลิก
- ไม่เขียน plan.md / tasks.md / test-cases.md

### สิ่งที่ AI ตีความเอง (ทีมต้องตรวจ)
- ASM-BRQ-04 ตีความว่า "แนบเอกสารประกอบ" (RQ-16) คือการอัปโหลดไฟล์เข้าระบบ ไม่ใช่เอกสารกระดาษนอกระบบ
  ตั้งเป็น Q-BRQ-06 ให้ยืนยัน
- AC-BRQ-05 เขียนว่า "ไม่เปิดเผยชื่อผู้ขอรายอื่น" ตาม DOM-PRIV-01 แต่รูปแบบจริงยังรอ Q-SH-05

---

## 2569-11-11 ทีมตอบ Q-BRQ-08 ค่าสถานะของคำขอยืม 8 ค่า

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (แตะของกลาง data-model.md เพราะเป็นตารางสถานะที่ทุกฟีเจอร์ใช้ร่วม)
- เครื่องมือ: Claude Code
- คำตอบของทีม: สถานะของคำขอมี 8 ค่า Requested / Approved by Officer / Pending Advisor Approval /
  Approved / Rejected / Handed Over / Returned / ยกเลิก
  โดย Approved by Officer และ Pending Advisor Approval ใช้กับการยืมค้างคืนเท่านั้น

### สิ่งที่ AI ทำ
- `specs/000-shared/data-model.md` v4 -> v5: เปลี่ยนหัวข้อ "สถานะของคำขอยืม" จากร่างเป็นตาราง 8 ค่า
  เพิ่มคอลัมน์ "ตั้งค่าโดย" เพื่อให้เห็นว่า 002 ตั้งได้ค่าเดียวคือ Requested ที่เหลือเป็นของ 003-006
- spec 002 -> Draft v1.1: FR-BRQ-04 ระบุว่า "รออนุมัติ" คือค่า `Requested` และ 002 ตั้งได้ค่าเดียว
- ปิด Q-BRQ-08, เพิ่ม Q-BRQ-09, อัปเดต CHANGELOG

### สิ่งที่ AI ไม่ได้เดา
- ค่าที่ 8 ทีมเขียนมาเป็นภาษาไทย ("ยกเลิก") ต่างจากอีก 7 ค่าที่เป็นอังกฤษ
  AI ไม่ตั้งชื่อค่าในฐานข้อมูลให้เอง เปิดเป็น Q-BRQ-09 ให้ทีมและเจ้าของ 003 ยืนยัน

### ต้องแจ้งเจ้าของฟีเจอร์อื่น
- ตารางสถานะนี้กระทบ 003 (ยกเลิก), 004 (อนุมัติ/ปฏิเสธ), 005 (ส่งมอบ), 006 (รับคืน)

---

## 2569-11-11 ทีมตกลงค่า action = request (ปิด Q-BRQ-07)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (แตะของกลาง data-model.md เพราะตาราง action ของ DM-07 ใช้ร่วม)
- เครื่องมือ: Claude Code
- คำตอบของทีม: ใช้ `request` ค่าเดียว ไม่แยกยืมปกติกับค้างคืน

### สิ่งที่ AI ทำ
- `specs/000-shared/data-model.md` v5 -> v6: เพิ่มแถว `request` (002) ในตาราง action ของ DM-07
  ระบุว่าใช้ค่าเดียวทั้ง 2 ประเภท และกรณีถูกปฏิเสธแยกด้วยคอลัมน์ outcome ที่มีอยู่แล้ว
- spec 002 -> Draft v1.2: ปิด Q-BRQ-07, ASM-BRQ-03 เพิ่มว่า migration 0003 ต้อง ALTER CHECK ของ audit_log.action
  โดยไม่แก้ไฟล์ 0002 ตามกติกาที่เขียนไว้ในไฟล์นั้น
- CHANGELOG

### ยังเปิดอยู่
- Q-BRQ-09 ชื่อค่าในฐานข้อมูลของสถานะ "ยกเลิก" (กระทบ 003)
- Q-BRQ-01 ถึง 06 และ Q-16 ของ 001

---

## 2569-11-11 ตรวจ BL-23 แล้วพบว่าเอกสารขัดกัน (ยังไม่ตัดสิน)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ถามว่า BL-23 อยู่ในฟีเจอร์ไหน และตัดออกจาก prototype ได้ไหม

### ข้อค้นพบ
- BL-23 ไม่ได้จับคู่กับฟีเจอร์ 001-006 ใด ๆ ใน SRS 3.1 และไม่อยู่ในตาราง 3.1.7 เพราะไม่มี use case รองรับ
- `specs/000-shared/constraints.md` (Q-SH-07) เขียนว่า BL-23 "ไม่อยู่ใน prototype"
  แต่ Use Case Spec UC-02 บน Miro เขียนกฎ RQ-20 ไว้ 4 จุด และ FR-BRQ-06 ของเราก็ทำข้อนี้อยู่
- BL-23 รวม RQ-20 (ปฏิเสธอัตโนมัติ ทำได้เลย) กับ RQ-27 (บทลงโทษ ยังไม่มีนิยาม) ไว้ในรายการเดียว
- ตาราง `loan` ของ 001 มี borrower_id, due_at, returned_at และ index idx_loan_open ครบแล้ว
  งานที่เหลือของ FR-BRQ-06 คือ query เดียว + error path + test (ประเมินคร่าว ๆ ครึ่งวัน)

### การตัดสินของทีม
- **ยังไม่ตัดสิน รอถามเจ้าหน้าที่ก่อน** AI ไม่แก้ spec และไม่ทำ CR
- บันทึกข้อขัดแย้งไว้ใน Q-BRQ-03 ของ spec 002 ให้ครบเพื่อถามทีเดียวจบ

---

## 2569-11-11 ทีมตอบคำถามรอบใหญ่ 11 ข้อ (spec 002 -> Draft v2)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (แตะของกลาง constraints.md และ data-model.md เพราะเป็นการปิด Q-SH)
- เครื่องมือ: Claude Code
- คำตอบที่ทีมส่งมา: Q-16 รายชิ้น / Q-BRQ-09 Cancelled / Q-BRQ-01 วัตถุประสงค์พิมพ์เอง /
  Q-BRQ-02 ช่องที่ต้องครบ 5 ช่อง + ช่วงเวลาเป็นวันที่-เวลาเริ่มถึงสิ้นสุด /
  Q-BRQ-03 เลย 16:00 ทันที, คืนแล้วยืมได้วันถัดไป, ทำใน prototype /
  Q-BRQ-06 PDF หรือรูปถ่าย 5GB แนบตอนส่งคำขอ / Q-BRQ-04 ปฏิเสธ / Q-BRQ-05 ไม่ขยาย /
  Q-SH-04 16:00 ของวันสิ้นสุด / Q-SH-05 เฉพาะเจ้าหน้าที่ ที่หน้า 001 เห็นชื่อผู้ขอได้ /
  Q-SH-06 ปฏิเสธอัตโนมัติ

### สิ่งที่ AI เขียนลง spec แล้ว
- FR-BRQ-03, 05, 06, 10, NFR-PERF-01, AC-BRQ-03, 05, 06, 10 และเพิ่ม AC-BRQ-11, AC-BRQ-12
- constraints.md v4 (ปิด Q-SH-04, Q-SH-05, Q-SH-06 บางส่วน, แก้ขอบเขต BL-23 ใน Q-SH-07)
- data-model.md v7 (Cancelled, กำหนดคืนค้างคืน), SRS v0.6, CHANGELOG

### คำตอบที่ AI ไม่นำไปเขียนทับ เพราะขัดกับข้อกำหนดเดิม (เปิดเป็นคำถามใหม่)
- Q-BRQ-10: "วัตถุประสงค์เป็นช่องพิมพ์เอง" ขัด RQ-18 ที่สั่งให้แยกประเภท และซ้ำกับช่อง "เหตุผลที่ยืม"
  ใน Q-BRQ-02 -> FR-BRQ-02 ยังไม่แก้
- Q-BRQ-11: ชื่อ-นามสกุลและรหัสนักศึกษา ระบบมีอยู่แล้วจาก IF-IDP-01 ให้กรอกเองหรือเติมให้
  (และ Q-20 ของ 001 ยังเปิดว่าชื่อมาจากไหน)
- Q-BRQ-12: ยืมปกติให้กรอกเวลาสิ้นสุด แต่ DOM-LOAN-01 บังคับ 16:00 วันเดียวกัน ถ้ากรอกเกินทำอย่างไร
- Q-BRQ-13: ขนาดไฟล์ 5GB น่าจะพิมพ์ผิด (เสนอ 10MB)

### ต้องแจ้งเจ้าของฟีเจอร์อื่น
- 001 ต้องทำ CR เพิ่ม FR แสดงลำดับคำขอให้เจ้าหน้าที่ (ผลของ Q-SH-05) และปิด Q-16 ใน spec ของตน
- 004 (ปฏิเสธอัตโนมัติ Q-SH-06), 005 (กำหนดคืนค้างคืน Q-SH-04), 003 (ค่า Cancelled)

---

## 2569-11-11 ตอบคำถามของเจ้าของ 003 เรื่องตาราง borrow_request

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (แตะของกลาง data-model.md เพราะเป็นตารางสถานะที่ใช้ร่วม)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ตอบคำถาม 4 ข้อที่เจ้าของ 003 ถามมา แล้วทำเป็นหน้าเว็บให้อ่านง่าย

### คำตอบที่ให้ไป (อ้างอิงของกลางและไฟล์จริง ไม่ได้เดา)
1. 002 สร้าง migration เอง ที่ `0003_borrow_request.up.sql` ตามตารางเลข migration ในของกลาง
   (ASM-11 ของ 003 ที่ว่าให้ 003 ร่างไว้ก่อน ไม่ตรงกับของกลาง) 003 ใช้ 0004, 0005 ตามเดิม ไม่ต้องขอเลขเพิ่ม
2. ชื่อคอลัมน์ที่ 003 เขียนไว้ใช้ได้ทั้งหมด ส่วนค่าสถานะ **ทีมตัดสินรับข้อเสนอของ 003 เป็นตัวพิมพ์เล็ก**
3. `cancel_reason` และ `cancelled_at` 002 สร้างให้ใน 0003 ขอให้ 003 ยืนยันชนิดข้อมูล
4. FK ของ `audit_log.request_id` มี และ 002 เป็นผู้เพิ่มใน 0003 (ตามคอมเมนต์ใน 0002 ของ 001)

### สิ่งที่ AI ตรวจพบและแก้
- ค่าสถานะที่ทีมเคาะเมื่อวาน (`Requested`, `Approved by Officer`) มีช่องว่างและตัวพิมพ์ใหญ่
  ไม่เข้ากับค่าอื่นในฐานข้อมูล (`available`, `normal`, `list`) จึงเสนอให้แยกชื่อใน spec
  ออกจากค่าในฐานข้อมูลแบบเดียวกับ CON-DATA-01 -> ทีมตัดสินใช้ตัวพิมพ์เล็ก data-model v8

### สิ่งที่ AI ไม่เติมให้
- กำหนดส่ง migration 0003 เพราะยังไม่มี plan.md และ tasks.md ของ 002 เขียนเป็นคำถามกลับไปหา 003 แทน

### ผลลัพธ์อื่น
- หน้าเว็บคำตอบสำหรับ 003 (artifact) และหน้าเว็บคำถามค้างของ 002 ที่อัปเดตเป็นเวอร์ชัน 3

---

## 2569-11-11 ตอบคำถามของ 004 และได้คำยืนยันจาก 003

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request
- เครื่องมือ: Claude Code

### 003 ตอบกลับ
- ยืนยันชนิดข้อมูล `cancel_reason VARCHAR(200) NULL` และ `cancelled_at TIMESTAMPTZ NULL`
- **003 ติดอยู่ที่ตาราง borrow_request ทั้งหมด 17 task** ขอให้ส่ง migration 0003 เร็วที่สุด
  และให้บอกเมื่อ push เข้า main แล้วเพื่อ merge เข้า branch ของตน

### 004 ถามมา และ AI ตอบให้ (ทำเป็นหน้าเว็บ artifact)
- คอลัมน์และชนิดข้อมูลทั้งตาราง, ค่าสถานะ 8 ค่าแบบ snake_case, FK ของ audit_log.request_id
- รหัสคำขอเป็น SERIAL ไม่ใช่ REQ-xxxx เพราะ audit_log.request_id เป็น INT

### 2 เรื่องที่ AI ไม่ตอบ เพราะไม่มีที่มาในเอกสาร (ส่งกลับเป็นคำถาม)
- **อีเมลผู้ขอ**: IF-NOT-01 บอกให้แจ้งเตือนผ่าน Email แต่ DM-08 มีแค่รหัสผู้ใช้ บทบาท ภาควิชา
  ไม่มีที่ไหนในระบบเก็บอีเมล -> กระทบ 004, 005, 006 ควรเปิดเป็นคำถามของกลาง (ยังไม่ได้เพิ่มใน constraints.md)
- **ASM-07 ของ 004** ที่ว่า 002 ต้องปฏิเสธคำขอที่ส่งหลัง 16:00 น.: ค้นทั้ง repo แล้วไม่มีข้อกำหนดไหนรองรับ
  16:00 น. ที่มีอยู่คือกำหนดคืนตาม DOM-LOAN-01 และจุดตัดเลยกำหนดคืนใน FR-BRQ-06 ซึ่งคนละเรื่องกัน
  -> ถ้าทีมต้องการต้องเพิ่มเป็น FR ใหม่ของ 002 และแก้ UC-02 บน Miro

---

## 2569-11-11 /clarify specs/002-borrow-request/spec.md (รอบที่ 1)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request
- เครื่องมือ: Claude Code (`.claude/commands/clarify.md`)
- คำสั่ง: `/clarify`

### สิ่งที่ AI ตรวจ
อ่าน spec.md (Draft v2.1), AGENTS.md และ specs/000-shared ทั้ง 3 ไฟล์ ไม่มีโฟลเดอร์ mockups/
พบจุดกำกวม ขัดแย้ง หรือทดสอบไม่ได้ 16 จุด แล้วถามกลับ 8 ข้อ

### คำถามที่ AI ถาม และคำตอบของทีม
1. ฟอร์มมีช่องเหตุผลกี่ช่อง เป็นตัวเลือกหรือพิมพ์ -> **1 ช่อง เป็นช่องพิมพ์**
2. ชื่อ-นามสกุลและรหัสนักศึกษา เติมให้หรือพิมพ์เอง -> **ผู้ขอพิมพ์เอง**
3. ยืมปกติที่เวลาคืนเกิน 16:00 น. ทำอย่างไร -> **ปฏิเสธไม่ให้ส่ง พร้อมข้อความอธิบาย**
4. "วันถัดไป" เริ่มเมื่อไร และนับข้ามภาควิชาไหม -> **00:00 น. ของวันถัดไป (Asia/Bangkok) นับเฉพาะภาควิชาตนเอง**
5. ไฟล์แนบขนาดเท่าไร เก็บที่ไหน ใครดูได้ -> **ไม่เกิน 10MB เก็บบนเครื่องที่รันระบบ เจ้าหน้าที่กับเจ้าของคำขอดูได้**
6. In scope ขัดกับ FR-BRQ-05 เรื่องการแสดงลำดับ -> **ตัดการแสดงออกจาก In scope ให้ 002 เปิดเผยผ่าน API**
7. กำหนดคืนของค้างคืนอยู่ที่ 002 หรือ 005 -> **002 เก็บแค่วันที่สิ้นสุด การแปลงเป็นกำหนดคืนเป็นงานของ 005**
8. audit log จะขยาย DOM-AUDIT-01 หรือใช้ RQ-25 -> **ใช้ RQ-25 อย่างเดียว ตัด DOM-AUDIT-01 ออก**

### สิ่งที่ AI แก้ (ทำ CR-003 เพราะ spec ผ่าน Draft v2 แล้ว)
- `specs/changes/CR-003.md` ใหม่
- spec 002 Draft v2.1 -> v3: Scope, Constraints, FR-BRQ-02, 03, 04, 05, 06, 08, 10,
  AC-BRQ-02, 03, 05, 10, 12 และเพิ่ม AC-BRQ-13
- ปิด Q-BRQ-01, 03, 10, 11, 12, 13 และ ASM-BRQ-04 / เพิ่ม ASM-BRQ-05, 06, 07 และ Q-BRQ-14, 15
- `docs/srs/SRS.md` 3.1.2 (v0.6 -> v0.7), `constraints.md` ตารางใครอ้างใช้, CHANGELOG

### สิ่งที่ AI ชี้ว่าคำตอบของทีมขัดกับข้ออื่น (ไม่ได้เลือกเอง)
- คำตอบข้อ 1 (ช่องพิมพ์อิสระ) **ขัดกับ RQ-18** ที่สั่งให้ "แยกประเภทวัตถุประสงค์"
  ทีมยืนยันซ้ำเป็นครั้งที่ 2 AI จึงเขียนตามคำตอบ แล้วบันทึกเป็น ASM-BRQ-05
  พร้อมระบุใน CR-003 ว่าต้องไปแก้ RQ-18 บน Miro และ SRS 3.1 ให้ตรงกัน
- คำตอบข้อ 2 (พิมพ์ชื่อและรหัสเอง) ทำให้ไม่ชัดว่า "ผู้ขอ" ที่ระบบใช้ตรวจสิทธิ์คือใคร
  AI เขียน ASM-BRQ-06 ว่ายังใช้รหัสผู้ใช้จาก IF-IDP-01 **และขอให้ทีมยืนยันอีกครั้ง**
- ที่เก็บไฟล์แนบยังไม่มี CON หรือ IF ข้อไหนรองรับ บันทึกเป็น ASM-BRQ-07

### ยังไม่ได้ทำ
- ยังไม่ commit / ยังไม่มี plan.md และ tasks.md
- ของที่ต้องตามต่ออยู่ใน CR-003 หัวข้อ "ยังไม่ได้ทำ ต้องตามต่อ"

### ทบทวนรอบสอง 2569-11-11 (วันเดียวกัน หลัง AI อธิบายผลข้างเคียง)
- ทีมอ่านคำอธิบาย 3 ข้อที่ AI ชี้ แล้วเปลี่ยนคำตอบ 2 ข้อ
  1. วัตถุประสงค์การยืม: จากช่องพิมพ์อิสระ -> **กลับไปใช้ตัวเลือกประเภท (dropdown)** ตรงกับ RQ-18 เดิม
     ผล: ไม่ต้องแก้ RQ-18 บน Miro อีกแล้ว ตัดงานค้างออก 1 รายการ และ ASM-BRQ-05 ถูกยกเลิก (คง ID ไว้)
  2. เจ้าของคำขอ: **ยืนยัน ASM-BRQ-06** ว่ายึดรหัสจากการล็อกอินเสมอ
     ส่วนชื่อและรหัสที่ผู้ขอพิมพ์เป็นข้อมูลประกอบในใบคำขอเท่านั้น
  3. ที่เก็บไฟล์แนบ: คงเป็น ASM-BRQ-07 ไม่ยกระดับเป็น CON ในรอบนี้
- AI แก้ FR-BRQ-02, FR-BRQ-03, FR-BRQ-10, AC-BRQ-02, AC-BRQ-03, AC-BRQ-10, SRS 3.1.2, CR-003 และ CHANGELOG ตาม
- คำถามใหม่ที่เกิดจากการกลับไปใช้ตัวเลือก: **Q-BRQ-16** รายการตัวเลือกมีค่าอะไรบ้าง
  มีตัวเลือก "อื่น ๆ ระบุ" ไหม และการยืมค้างคืนตาม RQ-16 ที่เขียนว่า "ระบุเหตุผล" ต้องมีช่องพิมพ์เพิ่มไหม

### ทีมตอบ Q-BRQ-16 (2569-11-11)
- วัตถุประสงค์การยืมมี 3 ค่า: เพื่อการเรียน / เพื่อทำกิจกรรม / อื่น ๆ (เลือกแล้วต้องพิมพ์ระบุ)
- การยืมค้างคืนใช้ช่องเดียวกันนี้ ไม่มีช่องเหตุผลแยกตาม RQ-16
- AI แก้ FR-BRQ-02, FR-BRQ-10 (ก), AC-BRQ-02, เพิ่ม AC-BRQ-14, SRS 3.1.2,
  data-model v9 (DM-06 ระบุ 3 ค่าให้ 003 และ 004 อ้างอิงได้) และปิด Q-BRQ-16
- ผล: CR-003 ไม่เหลือรายการ "ยังไม่ได้ทำ" เรื่องตัวเลือกวัตถุประสงค์แล้ว

---

## 2569-11-11 /plan specs/002-borrow-request/plan.md (Draft v1)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request
- เครื่องมือ: Claude Code (`.claude/commands/plan.md`)
- คำสั่ง: `/plan`

### สิ่งที่ AI ทำ
- สร้าง `specs/002-borrow-request/plan.md` ครบ 8 หัวข้อตามคำสั่ง
  (สรุปแนวทาง, เทคโนโลยี, โมเดลข้อมูล, API/หน้าจอ, ตารางตรวจ Constraints + ค่าที่ตั้งได้,
   แผนทดสอบจาก AC 14 ข้อ, ลำดับงาน 10 ขั้น, สิ่งที่ยังไม่ทำ)
- ใช้ Go ตาม AGENTS.md ไม่ใช่ Python FastAPI ที่เขียนไว้ในไฟล์คำสั่งของรายวิชา
- จัด migration 0003 เป็นขั้นที่ 1 และระบุว่าส่งขึ้น main ได้ทันทีโดยไม่ต้องรอขั้นอื่น
  เพราะ 003 ติดอยู่ 17 task และ 004 รอตารางเดียวกัน

### Constraint ที่ยังไม่ได้ใช้
- IF-NOT-01 (แจ้งเตือน Email) ฟีเจอร์นี้ไม่ส่งการแจ้งเตือน เป็นงานของ 004 ตาม RQ-21
  และ spec ของ 002 ไม่ได้อ้าง ID นี้อยู่แล้ว
- CON-TECH-02 ส่วน TLS ที่ reverse proxy ยังเป็นงานระดับระบบที่ไม่มีเจ้าของ (เหมือนที่ 001 บันทึกไว้)

### AC ที่ทดสอบยาก
- AC-BRQ-12 (คืนแล้วยืมใหม่ได้ 00:00 น. วันถัดไป) ต้องฉีดนาฬิกาเข้า service แบบเดียวกับ T-12 ของ 001
  ห้ามใช้ time.Now() ตรง ๆ ไม่งั้น test จะไม่เสถียรตามเวลาที่รัน
- AC-BRQ-09 (บันทึกไม่สำเร็จแล้วส่งใหม่) ต้องใช้ repository ปลอมที่ล้มเหลวครั้งแรก
- AC-BRQ-10 ส่วนไฟล์ 11MB ต้องสร้างไฟล์ชั่วคราวในเครื่องทดสอบ

### สิ่งที่ AI อยากเดาแต่ไม่ได้เดา (ติด Open Questions)
- ข้อความบนหน้าจอทุกจุดและหน้าตาฟอร์ม (Q-BRQ-14) -> ขั้นที่ 10 ยังไม่เริ่ม
- ลำดับการตรวจใน POST ว่าถ้าเป็นคนต่างภาควิชาและอุปกรณ์ถูกยืมพร้อมกัน จะตอบ 403 หรือ 409 ก่อน (Q-BRQ-15)
- ข้อความอธิบายผลของคำขอที่ส่งซ้อนให้ผู้ขอเห็น (Q-SH-06 ส่วนที่เหลือ)
- พฤติกรรมเมื่อระบบยืนยันตัวตนล่ม (Q-SH-08) ใช้ของ 001 ไปก่อน

### สิ่งที่ AI เสนอเอง ทีมต้องตรวจ
- ชนิดไฟล์รูปที่รับคือ image/jpeg และ image/png (spec เขียนแค่ "รูปถ่าย")
- unique index `(requester_id, item_id) WHERE status='requested'` ใช้บังคับ FR-BRQ-05 (ข)
  ที่ระดับฐานข้อมูลด้วย เพื่อกันการกดส่งซ้ำพร้อมกัน
- รหัส error ทั้ง 8 ตัวในตาราง Error contract (plan ข้อ 4) ยังไม่มีใน spec

### ทีมตอบระหว่างตรวจ plan (2569-11-11)
- ชนิดไฟล์รูปที่รับ: **image/jpeg และ image/png** (ยืนยันข้อเสนอของ AI) -> แก้ plan ข้อ 5 จาก "ข้อเสนอ" เป็น "ทีมยืนยัน"
- ถาม: unique index ที่เสนอกระทบใครไหม -> AI ไล่ดูทุกฟีเจอร์แล้วตอบว่า **ไม่กระทบ**
  เพราะเป็น partial index มีผลเฉพาะตอน INSERT แถวสถานะ requested ซึ่งมีแต่ 002 ที่ทำ
  ส่วน 003-006 เพียง UPDATE สถานะของแถวเดิม แถวจะหลุดจากเงื่อนไขทันที
- **ข้อค้นพบใหม่จากการไล่ดู**: FR-BRQ-05 (ข) เขียนครอบคลุมเฉพาะสถานะ `requested`
  แปลว่าหลังอนุมัติแล้วแต่ยังไม่ส่งมอบ ผู้ขอคนเดิมส่งคำขอชิ้นเดิมได้อีกใบ
  AI ไม่ตัดสินให้ เขียนเป็นคำถามไว้ใน plan ข้อ 8 รอทีมตัดสินว่าจะเพิ่มใน spec ไหม

### ทีมตัดสิน 2 เรื่องระหว่างตรวจ plan (2569-11-11) -> CR-004
1. **รหัส error เป็นมาตรฐานกลาง** -> เพิ่ม CON-TECH-03 ใน constraints.md (v4 -> v5)
   พร้อมรายการรหัส 8 ตัว โดย 4 ตัวแรกเป็นของเดิมที่ 001 ใช้อยู่แล้ว ไม่ต้องแก้โค้ดของ 001
2. **unique index กันถึงสถานะที่ยังไม่จบ** -> แก้ FR-BRQ-05 (ข) จากเฉพาะ `requested`
   เป็น 4 สถานะ (requested, approved_by_officer, pending_advisor_approval, approved)

### สิ่งที่ AI ทำ
- `specs/changes/CR-004.md` ใหม่ พร้อมตารางผลกระทบและสิ่งที่ต้องแลก
- constraints.md v5 (CON-TECH-03 + ตารางใครอ้างใช้), spec 002 Draft v3.1 (FR-BRQ-05 ข, AC-BRQ-11, Constraints)
- plan 002 Draft v1.1 (unique index 4 สถานะ, error contract อ้าง CON-TECH-03, ตารางตรวจเพิ่ม CON-TECH-03)
- SRS v0.8, CHANGELOG
- ลบคำถามเรื่องขอบเขต index ออกจาก plan ข้อ 8 เพราะทีมตัดสินแล้ว

### ต้องแจ้งในกลุ่ม
- CON-TECH-03 เป็นของกลางใหม่ ทุกฟีเจอร์ต้องใช้รูปแบบ {code, message} และลงทะเบียนรหัสใหม่ก่อนใช้
- 004 ต้องรู้ว่าหลังอนุมัติแล้วผู้ขอคนเดิมส่งคำขอชิ้นเดิมซ้ำไม่ได้จนกว่าคำขอจะจบ

---

## 2569-11-11 /tasks specs/002-borrow-request (tasks.md v1)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request
- เครื่องมือ: Claude Code (`.claude/commands/tasks.md`)
- คำสั่ง: `/tasks`

### ผลลัพธ์
- `specs/002-borrow-request/tasks.md`: **14 เฟส 23 task** แบ่งแบบชิ้นแนวตั้ง 1 เฟสต่อ 1 AC
  (ยกเว้นเฟส 3 ที่รวม AC-BRQ-02 กับ AC-BRQ-14 เพราะเป็น FR-BRQ-02 ข้อเดียวกัน
   และเฟส 8 ที่รวม AC-BRQ-06 กับ AC-BRQ-12 เพราะเป็น FR-BRQ-06 ข้อเดียวกัน)
- **5 task รอ Open Questions**: T-19, T-20 (รอ Q-BRQ-14), T-21 (รอ Q-BRQ-15),
  T-22 (รอ Q-SH-06), T-23 (รอ Q-SH-08)
- เฟส 0 มี 3 task ตามเพดานที่คำสั่งกำหนด โดย T-01 (migration 0003) ระบุให้ส่งขึ้น main ทันทีที่ผ่าน
  เพราะ 003 ติด 17 task และ 004 รอตารางเดียวกัน
- ตารางตรวจความครบ 3 ตาราง: AC 14 ข้อมี task ครบ / Constraint 11 ข้อมี task ครบ / FR 10 ข้อมี AC ครบทุกข้อ

### สิ่งที่ AI ตัดสินเอง ทีมต้องตรวจ
- **จุดตรวจทุกเฟสลองด้วย `curl` ไม่ใช่หน้าจอ** เพราะงานหน้าจอทั้งหมดติด Q-BRQ-14
  ซึ่งต่างจากรูปแบบปกติของคำสั่งที่ให้ลองบนหน้าจอ ถ้าทีมไม่ยอมรับต้องตอบ Q-BRQ-14 ก่อนเริ่มเฟส 1
- รวม AC 2 ข้อไว้ในเฟสเดียว 2 ครั้ง (เฟส 3 และเฟส 8) เพราะเป็น FR เดียวกันและผู้ใช้เห็นเป็นเรื่องเดียว
- แยก T-14 (รับไฟล์) กับ T-15 (เปิดดูไฟล์) ออกจากกัน เพราะรวมกันแล้วเกิน 2 ชั่วโมง

### task ที่ยากที่สุด
- T-13 (ปลดล็อก 00:00 น. ของวันถัดไป) ต้องฉีดนาฬิกาและเทียบเขตเวลา Asia/Bangkok ให้ถูก
  ถ้าใช้ time.Now() ตรง ๆ test จะไม่เสถียร
- T-14 (ไฟล์แนบ) แตะทั้ง multipart, ระบบไฟล์ และ validation 3 แบบในงานเดียว

### ยังไม่เริ่มทำ task ใด ๆ

---

## 2569-11-11 /implement T-01 (002-borrow-request) — เขียนเสร็จ แต่ยังรัน test ไม่ได้

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 0 ฐานราก ตาราง borrow_request ที่ 003 และ 004 รออยู่)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-01`

### ไฟล์ที่สร้างหรือแก้
- `backend/db/migrations/0003_borrow_request.up.sql` (ตรงกับช่องไฟล์ที่แตะ)
- `backend/db/migrations/0003_borrow_request.down.sql` (ตรงกับช่องไฟล์ที่แตะ)
- `backend/internal/platform/migrate_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** เพิ่ม test 3 ตัวต่อท้าย
  เหตุผล: ช่อง "เสร็จเมื่อ" ของ T-01 บังคับให้มี test แต่ไม่มีไฟล์ test ในรายการ
  และชุดเปิด container (testDB, execInTx) เป็น unexported ในแพ็กเกจ platform
  ถ้าสร้างไฟล์ใหม่ในแพ็กเกจของ 002 จะต้องคัดลอก harness มาราว 60 บรรทัด
  ทำแบบเดียวกับที่ T-02 และ T-31 ของ 001 เคยทำ (แก้ของกลางเพราะฟีเจอร์ 002)

### ผล test
- **ยังรันไม่ได้** เครื่องนี้ไม่มี Go (`go` ไม่อยู่ใน PATH และหาไม่เจอในตำแหน่งที่ติดตั้งปกติ GOROOT ว่าง)
  และ Docker Desktop ติดตั้งแล้วแต่ยังไม่ได้เปิด (testcontainers ต้องใช้)
- **ยังไม่เปลี่ยนสถานะ T-01 ใน tasks.md** เพราะช่อง "เสร็จเมื่อ" ยังพิสูจน์ไม่ได้

### สิ่งที่ AI เกือบต้องเดา แต่เขียนให้ test พิสูจน์แทน
- ชื่อ constraint ของ CHECK ใน `audit_log.action` ที่ 0002 เขียนแบบไม่ตั้งชื่อ
  AI ใช้ชื่อที่ PostgreSQL ตั้งให้อัตโนมัติคือ `audit_log_action_check`
  ถ้าชื่อไม่ตรง migration 0003 จะล้มทันทีตอนรัน ซึ่งจะเห็นผลในรอบ test
- `purpose_other` บังคับเมื่อ `purpose='other'` ตรวจที่ service (T-07) ไม่ได้ทำเป็น CHECK ในฐานข้อมูล
  เพราะ tasks.md แยกเรื่องนี้เป็น T-07 ไม่ใช่งานของ T-01

### สิ่งที่ AI ไม่ได้ทำ (ตามกติกา)
- ไม่แตะ `equipment_item` และไม่เขียนสถานะอุปกรณ์ใด ๆ (DOM-REQ-01)
- ไม่ใส่คอลัมน์กำหนดคืน (เป็นของ 005) และไม่ใส่คอลัมน์อีเมล (ยังไม่มีที่มาในระบบ)
- ไม่แก้ไฟล์ `0002_audit_log.up.sql` ตามกติกาที่เขียนไว้ในไฟล์นั้น

### ผล test ของ T-01 (รันเมื่อ 2569-11-11 หลังตั้งสภาพแวดล้อมเสร็จ)
- สภาพแวดล้อม: เครื่องนี้ติดตั้ง Go บน Windows ไม่ได้ เพราะ **Smart App Control บล็อก go.exe**
  (ติดตั้งสำเร็จที่ C:\Program Files\Go แต่รันไม่ได้) จึงย้ายไปรันใน **WSL2 Ubuntu** แทน
  ติดตั้ง Go 1.26.8 ด้วย apt และใช้ Docker Desktop ผ่าน WSL Integration
  **ต่อจากนี้งานรัน test ของฟีเจอร์ 002 ทำในฝั่ง Ubuntu ทั้งหมด**
- `go test ./internal/platform/ -run "TestMigration_BorrowRequest|TestMigration_Duplicate|TestMigration_AuditLogRequest" -v`
  - TestMigration_BorrowRequestConstraints PASS (34.63s)
  - TestMigration_DuplicateRequestBlocked PASS (0.03s)
  - TestMigration_AuditLogRequestAction PASS (0.01s)
  - ok salmonnorway/backend/internal/platform 35.182s
- **ผลพลอยได้**: migration 0003 รันผ่านทั้งไฟล์ แปลว่าชื่อ constraint `audit_log_action_check`
  ที่ AI เดาไว้ถูกต้อง (ถ้าผิดจะล้มที่บรรทัด DROP CONSTRAINT)
- เปลี่ยนสถานะ T-01 ใน tasks.md เป็น "เสร็จ รอทีมตรวจ"

### จุดตรวจเฟส 0: รันชุดเต็มแล้วพบว่า 0003 ทำ test ของ 001 พัง 1 ตัว (แก้แล้ว)
- `go test ./... -v` รอบแรก: test ของ 002 ผ่าน 3 ตัว แต่ **TestMigration_AuditLogAction ของ 001 FAIL**
  `insert or update on table "audit_log" violates foreign key constraint "fk_audit_log_request"`
- สาเหตุ: test ของ 001 (เขียนใน T-31) ใส่ `request_id = 42` ซึ่งเป็นเลขสมมติ
  เมื่อก่อนผ่านได้เพราะยังไม่มี foreign key พอ 0003 ของ 002 ใส่ FK ตามที่คอมเมนต์ใน 0002 สั่งไว้
  เลขที่ไม่มีอยู่จริงจึงถูกปฏิเสธ — **ไม่ใช่ migration ผิด แต่ข้อมูลใน test เก่าไม่สอดคล้องกับ FK ใหม่**
- **ทีมอนุญาตให้แก้ test ของ 001** (ธัญลักษณ์ 2569-11-11) AI จึงแก้ `backend/internal/platform/migrate_test.go`
  โดยเพิ่มการสร้างคำขอ id 42 ขึ้นก่อนในทรานแซกชันเดียวกัน พร้อมคอมเมนต์อธิบายว่า FK มาจาก 0003 (ASM-BRQ-03)
  เจตนาเดิมของ test ไม่เปลี่ยน (ยังพิสูจน์ว่า request_id และ asset_code ว่างได้ทั้งคู่)
  ไม่ได้แก้ test ให้ผ่านด้วยการลดเงื่อนไข และไม่ได้ถอด FK ออก
- ผลรอบที่สอง `go test ./... -v`: **ผ่านครบ 7 ตัว** (001 สี่ตัว + 002 สามตัว)
  ok salmonnorway/backend/internal/platform 4.625s
- **ต้องแจ้งเจ้าของ 001**: ตั้งแต่ migration 0003 เป็นต้นไป `audit_log.request_id` มี FK แล้ว
  การเขียน test หรือโค้ดที่ใส่ request_id ต้องอ้างคำขอที่มีอยู่จริง

---

## 2569-11-11 แก้ของกลาง platform/migrate.go ให้ใช้บน Windows ได้

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ของกลาง `backend/internal/platform/migrate.go` ที่ 001 สร้างใน T-02 ตาม AGENTS.md)
- เครื่องมือ: Claude Code
- ที่มา: เพื่อนในทีมที่ใช้ Windows เปิดระบบไม่ได้ ถาม AI ว่าทำไม
- สาเหตุที่ AI พบ: `Migrate` ต่อ path เป็น `"file://"+abs` บน Windows abs เป็น `C:\...\db\migrations`
  ได้ URL `file://C:\...` ซึ่ง Go แปลงไม่ได้ (invalid port) หลังบ้านจึงหยุดตั้งแต่ขั้นรัน migration
  AI ยืนยันด้วยการลองแปลง URL ทั้งสองแบบใน Go: แบบ `\` error / แบบ `/` ใช้ได้
- สิ่งที่แก้: `migrate.New("file://"+filepath.ToSlash(abs), ...)` บรรทัดเดียว บน Linux/Mac ไม่มีผล
- ก่อนแก้ดึง main ล่าสุด (fast-forward 10 commit รวม migration 0003 ของ 002)
- test: หลังบ้านผ่านทุก package (รวม migration 0001 ถึง 0003), หน้าจอ 10 ไฟล์ 17 test ผ่าน, `GOOS=windows go build` ผ่าน
- ยังไม่ได้ลองบนเครื่อง Windows จริง ให้เพื่อนที่ใช้ Windows ลองหลัง merge
- สาเหตุอื่นบน Windows ที่แจ้งทีม (ไม่ได้แก้โค้ด): ตั้งค่า `VITE_MOCK_TOKEN=...` แบบ Linux ใช้ใน PowerShell ไม่ได้,
  PostgreSQL ในเครื่องชน port 5432, ต้องเปิด Docker Desktop, ต้อง `npm install`

---

## 2569-11-11 spec 001 v4.4, plan v1.6, tasks v2.1 และร่าง CR-005 (หลัง PR ของ 002)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- ที่มา: สรุปผลกระทบต่อ 001 หลัง PR #6, #7 ของ 002 (CR-003, CR-004) ทีมส่งมาเป็นรูป AI อธิบายและตรวจกับ main
- คำถามที่ AI ถาม และคำตอบของทีม (Pannathorn J):
  1. ปุ่มส่งคำขอยืมอยู่หน้าไหน -> ทั้งสองหน้า (ตารางรายชิ้นและหน้ารายละเอียด)
  2. message ใน error body บังคับกับ 001 ไหม -> บังคับ
  3. เลข CR และ FR/AC ใหม่ -> CR-005 (จองในกลุ่มแล้ว merge ไฟล์ CR ก่อน), FR-EQS-11, AC-EQS-17 ต่อท้าย ทีมอนุมัติใน CR ก่อนแก้ spec
  4. API ของ 002 -> ใช้ `GET /api/v1/equipment/items/{assetCode}/requests` ตาม plan 002 ข้อ 4 002 ผูก route เองใน main.go
     ไม่ต้องแก้ handler.go ของ 001 / requesterName = ชื่อที่ผู้ขอพิมพ์เอง (FR-BRQ-03 ก) ไม่รวม Q-20
- สิ่งที่ AI แก้:
  - spec v4.4: Q-16 ปิด, Q-SH-05 ปิดชี้ไป CR-005, Q-12 บันทึก Q-SH-04 ปิด, เพิ่ม CON-TECH-03 ในรายการใช้ร่วม, Clarifications Session 2569-11-11
  - plan v1.6: Error contract ตาม CON-TECH-03 พร้อมข้อความที่เสนอ (ทีมเลือกเอง ไม่ได้มาจาก spec รอทีมตรวจ), canRequestLoan ทั้ง 2 หน้า,
    การส่งต่อไป UC-02 และลำดับคำขอ (CR-005), ข้อ 10 เพิ่ม Q-16, Q-SH-05 และเรื่อง 404/405
  - tasks v2.1: T-34 พร้อมทำ (เขียนไฟล์ที่แตะและเสร็จเมื่อใหม่), เพิ่ม T-37, T-20 ทำหลัง T-37, ตาราง Constraint และสิ่งที่ยังไม่ทำ
  - สร้าง specs/changes/CR-005.md สถานะ "เสนอ รอทีมอนุมัติ" และบันทึก CHANGELOG
- สิ่งที่ AI ไม่ได้ทำ: ไม่เพิ่ม FR-EQS-11 / AC-EQS-17 ใน spec (รออนุมัติ CR-005), ไม่แก้โค้ด, ไม่ commit
- ข้อค้นพบที่ทีมต้องรู้:
  - CR-004 เขียนว่า 001 "ไม่ต้องแก้โค้ด" แต่ body error ของ 001 ยังไม่มี message จึงเพิ่ม T-37
  - 404 / 405 ตอบเป็นข้อความธรรมดาของ Go ไม่ตรง CON-TECH-03 และรายการกลางยังไม่มีรหัส ต้องให้ทีมเพิ่ม
  - migration 0003 ของ 002 เพิ่ม action `request` ใน CHECK ของ audit_log แล้ว migration `notify` ของ 003 ต้องใส่ `request` ด้วย
## 2569-11-11 /implement T-02 (002-borrow-request)

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 0 ฐานราก ค่าที่ตั้งได้จาก spec)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-02`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/config/config.go` (ตรงกับช่องไฟล์ที่แตะ) เพิ่ม 9 ค่า + ฟิลด์ AttachmentDir
- `frontend/src/config.js` (ตรงกับช่องไฟล์ที่แตะ) เพิ่ม 5 ค่า
- `backend/internal/config/config_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** เป็น test ของ "เสร็จเมื่อ"
  (เหตุผลเดียวกับ T-01 และทำแบบเดียวกับที่ 001 เคยทำ)

### ค่าที่เพิ่ม ทุกตัวมีคอมเมนต์อ้าง ID
- config.go: LoanNormalReturnHour (DOM-LOAN-01), OverdueBlockUntilNextDayStart (FR-BRQ-06 ข),
  AttachmentMaxBytes และ AttachmentAllowedTypes (FR-BRQ-10 ก), PurposeValues และ PurposeOther (FR-BRQ-02),
  RequestStatuses และ StatusRequested (data-model v9, FR-BRQ-04), AuditActionRequest (Q-BRQ-07),
  AttachmentDir อ่านจาก env (ASM-BRQ-07)
- config.js: LOAN_NORMAL_RETURN_HOUR, ATTACHMENT_MAX_BYTES, ATTACHMENT_ALLOWED_TYPES,
  PURPOSE_OPTIONS (พร้อมข้อความไทยจากคำตอบ Q-BRQ-16), PURPOSE_OTHER
- **ไม่ได้เพิ่ม TimeZone ซ้ำ** เพราะ 001 ประกาศไว้แล้วใน config.go ใช้ตัวเดิม

### ผล test
- `go test ./internal/config/ -v`: TestConfig_BorrowRequestValues PASS, TestConfig_AttachmentDirFromEnv PASS
- `go test ./...`: auth, config, equipment, platform ok ทั้งหมด
- `npm test` (รันใน Ubuntu หลัง `npm install` ครั้งแรกของเครื่องนี้): **10 ไฟล์ 17 test ผ่านหมด**
- ตรวจด้วย grep แล้วไม่มีค่าเหล่านี้ฝังอยู่นอกไฟล์ config (ที่เจอเป็นคอมเมนต์และข้อมูลทดสอบของ 001)

### สิ่งที่ AI เกือบต้องเดา
- ค่าเริ่มต้นของ AttachmentDir (`./var/attachments`) เป็นค่าของเครื่องพัฒนาแบบเดียวกับรหัสผ่าน DB ของ T-01
  ไม่ใช่ requirement เปลี่ยนด้วย ATTACHMENT_DIR ได้
- OverdueBlockUntilNextDayStart เป็น bool ที่เป็น true เสมอ ทำตามชื่อที่ plan ข้อ 5 กำหนดไว้
  เพื่อให้กฎจาก spec อยู่ใน config ไม่ฝังในเงื่อนไข ถ้าทีมเห็นว่าไม่จำเป็นให้บอก จะถอดออกและแก้ plan

### หมายเหตุสภาพแวดล้อม
- เครื่องนี้ไม่เคยติดตั้ง node_modules ของ frontend มาก่อน (อยู่ใน .gitignore)
  ติดตั้ง Node ใน Ubuntu แล้ว `npm install` ครั้งแรก
  **ต้องติดตั้งและรันฝั่งเดียวกันเสมอ** ถ้าลงจาก Windows แล้วมารันใน WSL จะพังเพราะ shim คนละแบบ

---

## 2569-11-11 /implement T-03 (002-borrow-request) — ปิดเฟส 0

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 0 ฐานราก ข้อมูลตัวอย่าง)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-03`

### ข้อค้นพบ: ไม่ต้องเพิ่มข้อมูลใหม่เลย
ตรวจ seed.sql แล้วพบว่า 3 กรณีที่ "เสร็จเมื่อ" ของ T-03 ต้องการ มีอยู่แล้วครบจาก T-02 ของ 001
- CS-NB-0001 ถึง 0003 ว่างและเป็นของภาควิชา 1
- XX-NB-0001 ว่างแต่เป็นของภาควิชา 2 (รุ่น 3 -> department_id 2)
- ชิ้นที่ 5 ถูกยืมโดย student-002 กำหนดคืน 16:00 เมื่อวาน ยังไม่คืน

AI จึง**ไม่เพิ่มแถวใหม่** เพราะจะเป็นข้อมูลที่ไม่มีใครใช้ และเสี่ยงทำให้ test ของ 001
ที่นับจำนวน (รุ่น 1 มี 5 ชิ้น ว่าง 3 ชิ้น) พัง

### ไฟล์ที่สร้างหรือแก้
- `backend/db/seed/seed.sql` (ตรงกับช่องไฟล์ที่แตะ) **แก้เฉพาะคอมเมนต์หัวไฟล์**
  ระบุว่าแถวไหนรองรับจุดตรวจข้อไหนของ 002 และระบุว่าตาราง borrow_request ต้องไม่มีข้อมูลตั้งต้น
- `backend/internal/platform/migrate_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** เพิ่ม TestSeed_BorrowRequestScenarios
  (เหตุผลเดียวกับ T-01 และ T-02)

### ผล test
- `go test ./internal/platform/ -run TestSeed -v`
  - TestSeed_Model1HasFiveItemsThreeAvailable PASS (5.22s) ของ 001 ยังผ่าน
  - TestSeed_BorrowRequestScenarios PASS (0.01s)
  - ok salmonnorway/backend/internal/platform 5.804s

### สิ่งที่ AI เพิ่มเองโดยไม่ได้สั่ง (ทีมตรวจได้)
- test ข้อสุดท้ายบังคับว่า **ตาราง borrow_request ต้องว่างเปล่าในข้อมูลตั้งต้น**
  เพื่อกันไม่ให้ใครใส่คำขอไว้ล่วงหน้าแล้วทำให้จุดตรวจของเฟสถัดไปผ่านแบบไม่จริง
  ไม่ได้อยู่ใน "เสร็จเมื่อ" แต่เป็นเกราะที่สอดคล้องกับ FR-BRQ-04 ที่ว่าคำขอเกิดจากการส่งเท่านั้น

### หมายเหตุ
- เฟส 0 ครบทั้ง 3 task แล้ว เหลือรันจุดตรวจเฟส 0 ให้ครบ (go test ./... และลองมือด้วย psql)

---

## 2569-11-11 /implement T-04 (002-borrow-request) — เริ่มเฟส 1

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 1 ส่งคำขอยืมสำเร็จ AC-BRQ-04)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-04`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/model.go` (ตรงกับช่องไฟล์ที่แตะ) โครงข้อมูลและ error ของฟีเจอร์
- `backend/internal/borrowrequest/repository.go` (ตรงกับช่องไฟล์ที่แตะ) INSERT ... SELECT และ ItemStatus
- `backend/internal/borrowrequest/service.go` (ตรงกับช่องไฟล์ที่แตะ) กฎของฟีเจอร์ อ่านผู้ขอจาก auth
- `backend/internal/borrowrequest/service_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** เป็น test ของ "เสร็จเมื่อ"
  (เหตุผลเดียวกับ T-01 ถึง T-03) ใช้ชุดเปิด container แบบเดียวกับ handler_test.go ของ 001

### การตัดสินใจด้านโครงสร้าง
- ใช้ `INSERT ... SELECT ... WHERE asset_code = $10 AND status = 'available' RETURNING` คำสั่งเดียว
  การอ่านสถานะกับการเขียนคำขอจึงอยู่ในคำสั่งเดียวกัน ไม่มีช่องให้สถานะเปลี่ยนระหว่างกลาง
  (plan.md ข้อ 3 สั่งให้อยู่ในทรานแซกชันเดียวกัน คำสั่งเดียวคือทรานแซกชันเดียวโดยปริยาย)
- ไม่มีคำสั่งใดใน repository แตะตาราง equipment_item เลย (DOM-REQ-01) และ test ตรวจซ้ำหลังส่งคำขอ
- `requester_id` มาจาก auth.FromContext เท่านั้น ไม่ใช่ student_code ที่ผู้ขอพิมพ์ (ASM-BRQ-06)

### ผล test
- `go test ./internal/borrowrequest/ -v`: ok 4.680s (TestAC_BRQ_04,
  TestCreate_RejectsItemNotAvailable, TestCreate_RequiresPrincipal)
- `go test ./...`: auth, borrowrequest, config, equipment, platform ok ทั้งหมด

### สิ่งที่ AI ทำล่วงหน้าหรือคาบเกี่ยว task อื่น (ทีมตรวจได้)
- ใส่ `WHERE status = 'available'` ไว้แล้วใน T-04 เพราะ CON-DATA-01 อยู่ในช่อง "รองรับ" ของ T-04
  เหลือให้ T-09 ทำเรื่องการแจ้งสถานะล่าสุดและแปลงเป็น HTTP 409
- เพิ่มเมท็อด `ItemStatus()` ใน repository ซึ่งยังไม่มีใครเรียกในรอบนี้ เตรียมไว้ให้ T-09
- `Service` มีฟิลด์ `now` ที่ยังไม่ได้ใช้ เตรียมไว้ให้ T-12 และ T-13 ที่ต้องฉีดนาฬิกา

---

## 2569-11-11 /implement T-05 (002-borrow-request) — ปิดเฟส 1

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 1 ส่งคำขอยืมสำเร็จ AC-BRQ-04)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-05`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/handler.go` (ตรงกับช่องไฟล์ที่แตะ) POST /v1/loan-requests
- `backend/internal/borrowrequest/dto.go` (ตรงกับช่องไฟล์ที่แตะ) CreateRequestDTO, CreatedDTO, ErrorDTO
- `backend/cmd/server/main.go` (ตรงกับช่องไฟล์ที่แตะ) ผูก handler เข้า router หลัง auth.Middleware
- `backend/internal/borrowrequest/handler_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** เป็น test ของ "เสร็จเมื่อ"

### การตัดสินใจ
- `ErrorDTO` ของฟีเจอร์นี้เป็น `{code, message}` ตาม CON-TECH-03 (CR-004)
  ต่างจาก ErrorDTO ของ 001 ที่มีแค่ `{code}` เพราะเขียนก่อนมี constraint นี้
  CR-004 ระบุไว้แล้วว่า 001 ไม่ต้องแก้โค้ด
- path เป็น `/v1/loan-requests` ไม่ใช่ `/api/v1/...` เพราะ proxy ใน vite.config.js ตัด /api ออก
  (แบบเดียวกับ equipment.PathPrefix ของ 001)
- ตั้งชื่อ test ว่า `TestAC_BRQ_04_HTTP` เพราะชื่อ `TestAC_BRQ_04` ถูกใช้ใน T-04 แล้ว
  ชื่อซ้ำกันใน package เดียวไม่ได้ แต่ยังมี AC ID อยู่ในชื่อตามกติกา

### ผล test
- `go test ./internal/borrowrequest/ -v`: ok 4.947s
  (TestAC_BRQ_04, TestAC_BRQ_04_HTTP, TestCreate_WithoutToken_401, TestCreate_BadJSON_422,
   TestRoutes_OnlyPost, TestCreate_RejectsItemNotAvailable, TestCreate_RequiresPrincipal)
- `go build ./...` ผ่าน ยืนยันว่า main.go ที่แก้คอมไพล์ได้
- `go test ./...`: ทุกแพ็กเกจ ok

### สิ่งที่ AI ทำคาบเกี่ยว task อื่น (ทีมตรวจได้)
- แปลง ErrItemNotAvailable เป็น 409 ITEM_NOT_AVAILABLE แล้วใน handler
  แต่ **ยังไม่แนบสถานะล่าสุดกลับไปใน body** ซึ่งเป็นส่วนที่ T-09 ต้องเติม (เขียนคอมเมนต์กำกับไว้)
- body ที่อ่านไม่ออกตอบ 422 VALIDATION_FAILED ถือเป็นเรื่องของ handler
  ไม่ใช่การตรวจความครบถ้วนรายช่องซึ่งเป็นงานของ T-06

---

## 2569-11-11 /implement T-06 (002-borrow-request) — เฟส 2 ตรวจความครบถ้วนของฟอร์ม

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 2 AC-BRQ-03)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-06`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/validate.go` (ตรงกับช่องไฟล์ที่แตะ) ฟังก์ชัน missingFields ตรวจ 5 ช่อง
- `backend/internal/borrowrequest/handler.go` (ตรงกับช่องไฟล์ที่แตะ) เรียก missingFields ก่อนเรียก service
  และเพิ่ม writeValidationError
- `backend/internal/borrowrequest/handler_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** ต่อ test ท้ายไฟล์เดิม

### การตัดสินใจ
- CR-004 ระบุว่า VALIDATION_FAILED ต้อง "แนบรายชื่อช่องที่ผิด" แต่ ErrorDTO ใน dto.go มีแค่ {code, message}
  และ dto.go ไม่อยู่ในช่องไฟล์ที่แตะของ T-06 จึงประกาศ validationErrorDTO ไว้ใน handler.go แทน
  ยังมี code กับ message ครบตาม CON-TECH-03 เพิ่มแค่ fields
  **ถ้าทีมอยากให้รวมไว้ที่ dto.go ให้ทำตอน task หลังที่แตะ dto.go อยู่แล้ว**
- ใช้ strings.TrimSpace ไม่ให้ช่องว่างล้วนผ่าน และ time.IsZero สำหรับวันเวลาที่ไม่ได้ส่งมา
- เขียนคอมเมนต์กันขอบเขตไว้ใน validate.go ว่าค่าวัตถุประสงค์ 3 ค่าเป็นของ T-07
  และกฎเวลาคืน 16:00 เป็นของ T-08

### ผล test
- `go test ./internal/borrowrequest/ -run "TestAC_BRQ_03|TestValidation" -v`: PASS
- `go test ./...`: ทุกแพ็กเกจ ok (borrowrequest 4.922s)

### สิ่งที่ AI ไม่ได้ทำ (ตามขอบเขต)
- ไม่ตรวจว่า assetCode ว่างหรือไม่ เพราะไม่ได้อยู่ใน 5 ช่องที่ทีมตัดสิน
  รหัสครุภัณฑ์มาจาก FR-EQS-06 ไม่ใช่ช่องที่ผู้ขอพิมพ์ ถ้าไม่มีจะถูกปฏิเสธที่ชั้นฐานข้อมูลเป็น ITEM_NOT_AVAILABLE
- ไม่ตรวจว่า purpose เป็น 1 ใน 3 ค่า (T-07) และไม่ตรวจกฎ 16:00 (T-08)

---

## 2569-11-12 /implement T-07 (002-borrow-request) — ปิดเฟส 3

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 3 AC-BRQ-02 และ AC-BRQ-14)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-07`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/validate.go` (ตรงกับช่องไฟล์ที่แตะ) เพิ่ม purposeProblems
- `backend/internal/borrowrequest/service.go` (ตรงกับช่องไฟล์ที่แตะ) เพิ่ม normalizePurpose
- `backend/internal/borrowrequest/handler_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** ต่อ test ท้ายไฟล์เดิม

### ข้อจำกัดที่ต้องออกแบบรอบ
`handler.go` **ไม่อยู่ในช่องไฟล์ที่แตะของ T-07** จึงเพิ่ม branch ใหม่ใน handler ไม่ได้
ทางออกคือให้กฎใหม่ไหลผ่าน missingFields() ที่ handler เรียกอยู่แล้ว (เรียก purposeProblems ต่อท้าย)
ผลคือไม่ต้องแตะ handler เลย แต่ชื่อ missingFields ตอนนี้ครอบคลุมทั้ง "ไม่ได้กรอก" และ "กรอกค่าที่ไม่รู้จัก"
แก้คอมเมนต์หัวฟังก์ชันให้ตรงความหมายใหม่แล้ว **ถ้าทีมอยากเปลี่ยนชื่อเป็น invalidFields
ต้องทำตอน task ที่มี handler.go อยู่ในรายการ**

### การตัดสินใจอื่น
- ค่าที่ไม่รู้จัก (เช่น "exam" หรือส่งข้อความไทย "เพื่อการเรียน" มาแทน value) ถือเป็นปัญหาของช่อง purpose
  เหมือนกับกรณีไม่ได้เลือก ชี้ช่องเดียวกัน ไม่แยกรหัส error ใหม่
- normalizePurpose ล้างข้อความพิมพ์ระบุทิ้งเมื่อไม่ได้เลือก "อื่น ๆ" (FR-BRQ-02 ข้อ ข)
  กันข้อความค้างจากการที่ผู้ขอพิมพ์ไว้แล้วเปลี่ยนใจ
- ใช้ config.PurposeValues และ config.PurposeOther จาก T-02 ไม่มีการเขียน "study" หรือ "other" ฝังในตรรกะ

### ผล test
- `go test ./internal/borrowrequest/ -run "TestAC_BRQ_02|TestAC_BRQ_14" -v`: PASS
- `go test ./...`: ทุกแพ็กเกจ ok (borrowrequest 4.936s)

---

## 2569-11-11 เขียน FR ของ UC-08 ใน SRS หัวข้อ 3.1.5

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- คำสั่ง: "อ่าน AGENTS.md, specs/000-shared/ (ทั้ง 3 ไฟล์) และ SRS หัวข้อ 3.0 กับหัวข้อของฟีเจอร์ตัวเอง (005-handover)"
  แล้ว "เขียน FR ของตัวเองใน SRS (ตอนนี้ยังเป็น [เติม])"
- ไฟล์ที่แก้: `docs/srs/SRS.md` หัวข้อ 3.1.5 แทน `[เติม]` ด้วย FR-HND-01 ถึง FR-HND-08 และคำถาม Q-HND-01 ถึง Q-HND-06
- ที่มา: RQ-04, RQ-08, RQ-14, RQ-15, RQ-17, RQ-19, RQ-23, RQ-25 และ `specs/000-shared/` เท่านั้น
  AI ไม่มีข้อความ US-008 และ Use Case Specification ของ UC-08 (อยู่บน Miro) เจ้าของฟีเจอร์ต้องเทียบเอง
- ข้อที่ AI ไม่ตัดสินแทน (เขียนเป็นคำถาม): RQ-19 กับ RQ-23 และ US-003 เรื่องตรวจบัตร, รายละเอียดรูป RQ-04,
  ส่งมอบหลัง 16:00, ข้อจำกัดพื้นที่มหาวิทยาลัย, ชื่อสถานะคำขอ "ส่งมอบแล้ว", ค่า action ของ audit log

---

## 2569-11-11 สร้าง spec.md ของ 005-handover (Draft v1)

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (แก้ของกลาง `specs/000-shared/constraints.md` เพราะฟีเจอร์ 005)
- เครื่องมือ: Claude Code
- คำสั่ง: สร้าง `specs/00N-<feature>/spec.md` ครบ 6 ส่วน + Assumptions & Open Questions, Constraint ร่วมอ้าง ID แบบ "ใช้ร่วม",
  ใส่ชื่อตัวเองในช่อง Owner, เพิ่มชื่อฟีเจอร์ในตาราง "ใครอ้างใช้บ้าง" ท้าย constraints.md
- ไฟล์ที่สร้าง: `specs/005-handover/spec.md` (FR-HND-01 ถึง 08 จาก SRS 3.1.5, AC-HND-01 ถึง 09, Q-HND-01 ถึง 09)
- ไฟล์ที่แก้: `specs/000-shared/constraints.md` (ตาราง "ใครอ้างใช้บ้าง" เพิ่ม 005-handover 8 แถว),
  `specs/README.md` (สถานะ spec Draft v1), `specs/CHANGELOG.md`, `docs/srs/SRS.md` หัวข้อ 3.1.5 (อ้าง Draft v1 และ Q-HND-07 ถึง 09)
- ID ใหม่ที่ AI เสนอ (ทีมตรวจ): CON-HND-01 บทบาทเจ้าหน้าที่เท่านั้น (จาก actor ของ UC-08), Q-HND-07 ถึง Q-HND-09
- ข้อที่ AI ไม่ตัดสินแทน: ส่งมอบยืมค้างคืนต้องทดสอบเทอมนี้ไหม (Q-HND-07), เลือกรายชิ้นตอนส่งมอบถ้า Q-16 ของ 001 เปลี่ยน (Q-HND-08),
  เวลาตอบสนอง (Q-HND-09) และคำถามเดิม Q-HND-01 ถึง 06
- ยังไม่ได้เทียบกับ US-008 และ Use Case Specification ของ UC-08 บน Miro

---

## 2569-11-11 คำสั่ง: /clarify

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์: specs/005-handover/spec.md (v1)

### คำถามที่ AI ถาม (ทั้งหมด)

1. ตอนส่งมอบ ต้องตรวจบัตรทุกครั้ง (RQ-19) หรือเฉพาะเมื่อ Auth API ใช้ไม่ได้ (RQ-23) และระบบต้องบันทึกอะไรจากการตรวจบัตร? (FR-HND-02, FR-HND-03)
2. ณ เคาน์เตอร์ ผู้ยืมยืนยันตัวตนผ่าน Auth API อย่างไร และรอนานกี่วินาทีจึงถือว่าเชื่อมต่อไม่สำเร็จ? (FR-HND-02, IF-IDP-01)
3. คำขอยืมระบุอุปกรณ์รายชิ้นมาแล้ว หรือเจ้าหน้าที่ต้องเลือกรหัสครุภัณฑ์ตอนส่งมอบ? (FR-HND-04, FR-HND-07, Q-16 ของ 001)
4. ถ้าคนที่ยืนยันตัวตนตอนรับของไม่ใช่ผู้ส่งคำขอ ระบบต้องปฏิเสธการส่งมอบไหม? (FR-HND-01, FR-HND-02)
5. เจ้าหน้าที่หาคำขอที่อนุมัติแล้วเพื่อส่งมอบอย่างไร? (FR-HND-01)
6. หลังส่งมอบแล้ว คำขออื่นที่ค้างอยู่ของอุปกรณ์ชิ้นเดียวกันเป็นอย่างไร? (FR-HND-07, DOM-REQ-01, Q-SH-06)
7. ถ้าส่งมอบยืมปกติหลัง 16:00 ระบบต้องห้ามบันทึก หรือให้กำหนดคืนเป็นวันถัดไป? (FR-HND-05)
8. รูปสภาพอุปกรณ์ตอนส่งมอบบังคับแนบไหม แนบได้กี่รูป และไฟล์ชนิดใด ขนาดไม่เกินเท่าไร? (FR-HND-06)

### คำตอบของทีมและเหตุผล

ทีมตอบในแชต (Banlu C) ไม่ได้ระบุว่าถามผู้ใช้จริง จึงบันทึกเป็น ASM (ทีมตัดสิน)
1. ทุกครั้ง บันทึกชื่อ นามสกุล รหัสนักศึกษา วันเดือนปีเวลาที่ส่งมอบ -> ASM-HND-01
2. เกิน 10 วินาที -> ASM-HND-02 ส่วนวิธียืนยันตัวตนที่เคาน์เตอร์ไม่ได้ตอบ ย้ายไป Q-HND-11 (ถามนักเทคโนโลยีสารสนเทศ)
3. ระบุรายชิ้น -> ASM-HND-03 (ขัดกับคำตอบเบื้องต้นของ Q-16 ใน 001 ถ้า Q-16 ยืนยันเป็นระดับรุ่นต้องทำ CR)
4. ปฏิเสธ -> ASM-HND-04
5. ดูรายการทั้งหมด -> ASM-HND-05 (AI ตีความเป็นรายการคำขอที่ได้รับอนุมัติแล้วและยังไม่ส่งมอบทั้งหมด)
6. ปฏิเสธ และตอบกลับผู้ใช้งาน -> ASM-HND-08 ช่องทางแจ้งย้ายไป Q-HND-12 (Q-SH-06 เป็นคำถามร่วม ต้องตกลงกับ 002-004)
7. บันทึกตามเวลาที่ส่งจริงได้เลย แต่ไม่มีให้ยืมหลัง 16:00 -> ASM-HND-06 (AI ตีความว่าส่งมอบหลัง 16:00 ระบบไม่บันทึก)
8. บังคับอย่างน้อย 2 รูป jpg/png ไม่เกิน 10 MB -> ASM-HND-07 "10 MB ต่อรูปหรือรวม" ย้ายไป Q-HND-10

### สิ่งที่แก้ใน spec.md (v1 เป็น v2)

- FR-HND-01, 02, 03, 05, 06, 07 ปรับข้อความตามคำตอบ (ID เดิม ไม่เพิ่ม FR) FR-HND-04 เปลี่ยนคำเป็น "การยืม"
- AC-HND-02, 03, 06, 07 ปรับ, เพิ่ม AC-HND-10 ถึง 14
- เพิ่ม ASM-HND-01 ถึง 08, Q-HND-10 ถึง 15 ปิด Q-HND-01, 02, 03, 08 เพิ่มหัวข้อ Clarifications
- Status เป็น Draft v2
- ไฟล์อื่น: `docs/srs/SRS.md` หัวข้อ 3.1.5 และตาราง 4.3 (RQ-19 vs US-003), `specs/CHANGELOG.md`, `specs/README.md`

---

## 2569-11-11 คำสั่ง: /clarify (ยืนยันการตีความ)

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์: specs/005-handover/spec.md (v2 เป็น v2.1)

### คำถามที่ AI ถาม

1. Q5 "ดูรายการทั้งหมด" คือคำขอที่อนุมัติแล้วและยังไม่ส่งมอบทั้งหมดใช่ไหม
2. Q7 "ไม่มีให้ยืมหลัง 16:00" หมายความว่าระบบปฏิเสธการส่งมอบหลัง 16:00 ใช่ไหม
3. Q6 ปฏิเสธคำขออื่นหลังส่งมอบ เป็นคำถามร่วม Q-SH-06 ยังไม่ได้แก้ constraints.md

### คำตอบของทีมและเหตุผล

1. ใช่ -> ASM-HND-05 ยืนยันแล้ว
2. ไม่ใช่ ระบบบันทึกการส่งมอบหลัง 16:00 ได้ -> แก้ FR-HND-05, ASM-HND-06, AC-HND-11
   AI ชี้ว่ายังไม่รู้กำหนดคืนของกรณีนี้ (16:00 วันเดียวกันจะเลยกำหนดทันที) -> เปิด Q-HND-03 อีกครั้ง (ถามนักคอมพิวเตอร์)
3. ทีมถามกลับว่ากระทบใครบ้าง -> AI อธิบายในแชต ยังไม่แก้ไฟล์ของกลาง

### สิ่งที่แก้ใน spec.md (v2 เป็น v2.1)

- FR-HND-05, AC-HND-11, ASM-HND-05, ASM-HND-06, Q-HND-03 (เปิดอีกครั้ง), Clarifications เพิ่ม session ยืนยันการตีความ
- ไฟล์อื่น: `docs/srs/SRS.md` 3.1.5, `specs/CHANGELOG.md`, `specs/README.md`

---

## 2569-11-11 คำสั่ง: /clarify (Q6 หาแนวทางที่ไม่กระทบฟีเจอร์อื่น)

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์: specs/005-handover/spec.md (v2.1 เป็น v2.2)

### คำถามที่ AI ถาม / เสนอ

1. ทีมถาม: Q6 มีแนวทางที่ไม่กระทบคนอื่นไหม AI เสนอ 2 แบบ
   - แบบ A: 005 ไม่แตะคำขออื่น เจ้าหน้าที่ปฏิเสธเองผ่าน UC-12 (RQ-14) และ 004 แจ้งผล (RQ-21) FR-HND-04 กันส่งมอบซ้ำอยู่แล้ว
   - แบบ B: แบบ A + เตือนเจ้าหน้าที่ว่ามีคำขอค้าง (AI ไม่แนะนำ เพราะเป็นความต้องการใหม่ที่ไม่มี RQ)
2. ทีมถาม: 002 กระทบอย่างไร ต้องทำอะไรเพิ่ม

### คำตอบของทีมและเหตุผล

1. เลือกแบบ A -> แก้ FR-HND-07, ASM-HND-08, AC-HND-13, Scope ยกเลิก Q-HND-12
2. AI ตอบในแชต: 002 ไม่ต้องเพิ่มพฤติกรรม เหลือยืนยันชื่อสถานะ "ส่งมอบแล้ว" และฟิลด์ในคำขอ -> รวมไว้ใน Q-HND-05

### สิ่งที่แก้ใน spec.md (v2.1 เป็น v2.2)

- Scope, FR-HND-07, AC-HND-13, ASM-HND-08, Q-HND-05, Q-HND-12 (ยกเลิก), Q-SH-06, Clarifications
- ไฟล์อื่น: `docs/srs/SRS.md` 3.1.5, `specs/CHANGELOG.md`, `specs/README.md`
- ไม่ได้แก้ไฟล์ของกลาง (`specs/000-shared/`) Q-SH-06 ยังเปิด

---

## 2569-11-11 คำสั่ง: /plan

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ก่อนทำ: ทีมสั่งให้ทิ้ง commit `c90274b 005: spec v2.1` บน GitHub (spec ฉบับเก่าจาก session ก่อน) ทั้ง commit
  แล้วใช้ spec v2.2 ของ session นี้แทน AI แก้ชื่อ commit เป็น `67b2d07 005: spec v2.1` และ force push (with lease) ทับ
- ผลลัพธ์: specs/005-handover/plan.md (Draft v1) อ้าง spec Draft v2.2
  - โค้ดอยู่ที่ `backend/internal/handover/` และ `frontend/src/pages/Handover/` migration ใหม่ `00NN_handover` (เลขรอจองในกลุ่ม)
  - บันทึกการส่งมอบเป็น transaction เดียว ใช้ `SELECT ... FOR UPDATE` กันส่งมอบชิ้นเดียวกันพร้อมกัน
  - ค่าที่ตั้งได้: 10 วินาที, 16:00, อย่างน้อย 2 รูป, 10 MB, JPG/PNG อยู่ใน `config.go` และ `frontend/src/config.js`
- Constraint ที่ AI ยังไม่ได้ใช้: DOM-LOAN-01 ส่วน "ใช้ได้เฉพาะในพื้นที่มหาวิทยาลัย" (รอ Q-HND-04)
- สิ่งที่ AI บอกว่าอยากเดาแต่ไม่ได้เดา: กำหนดคืนเมื่อส่งมอบหลัง 16:00 (Q-HND-03), วิธียืนยันตัวตนผู้ยืมที่เคาน์เตอร์ (Q-HND-11),
  ชื่อสถานะ "ส่งมอบแล้ว" และค่า action (Q-HND-05, 06), ที่มาของชื่อ-นามสกุลผู้ยืม (Q-HND-13), 10 MB ต่อรูปหรือรวม (Q-HND-10)
- คำถามใหม่จากการวางแผน: Q-P1 ที่เก็บไฟล์รูป, Q-P2 กฎของ 004 ที่ให้อนุมัติได้ทีละ 1 คำขอต่อชิ้น (กระทบ AC-HND-13)

---

## 2569-11-11 คำสั่ง: /plan (ตอบคำถามจากการวางแผน)

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- คำถามที่ AI ถาม: Q-P1 เก็บไฟล์รูปที่ไหน (AI เสนอ A: Docker volume, B: bytea ในฐานข้อมูล, C: MinIO)
  Q-P2 spec ของ 004 กระทบ 005 อย่างไร (AI อ่าน spec ของ 004 แล้วพบ 5 จุด และเสนอทางที่ไม่กระทบฟีเจอร์อื่น)
- คำตอบของทีม: Q-P1 เลือก A / Q-P2 เลือกทำจุด 1 (ปรับ AC-HND-13, 14 ให้ตรงกับ AC-APV-08 ของ 004)
- สิ่งที่แก้: plan.md (Q-P1 ปิด, Q-P2 ปิดบางส่วน, อ้าง spec v2.3), spec.md Draft v2.3 (AC-HND-13, 14),
  `specs/CHANGELOG.md`, `specs/README.md`, `docs/srs/SRS.md` (อ้าง Draft v2.3)
- ยังรอทีม: จุด 3 ส่งมอบหลัง 16:00 (3A กำหนดคืน 16:00 เสมอ / 3B ไม่ให้ส่งมอบตั้งแต่ 16:00) และจุด 4 จำกัดภาควิชาของเจ้าหน้าที่ไหม

---

## 2569-11-11 คำสั่ง: /tasks

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- อินพุต: specs/005-handover/plan.md (Draft v1), spec.md (Draft v2.3)
- ผลลัพธ์: specs/005-handover/tasks.md (v1) 14 เฟส (ฐานราก + 12 เฟสตาม AC + รอคำตอบ) รวม 31 task รอ Q-xx 16 task
  - เฟส 2 รวม AC-HND-02 กับ AC-HND-07 เพราะการส่งมอบสำเร็จครั้งแรกต้องยืนยันผู้ยืมและตรวจบัตรก่อนเสมอ
  - AC-HND-08 ย้ายไปเฟสรอคำตอบทั้งเฟส (Q-HND-06, Q-HND-14)
  - task ที่ติดมากที่สุดคือ Q-HND-05 (7 task) และ Q-HND-13 (2 task) เฟสแรกที่ลองได้จริง (เฟส 1) ต้องรอ Q-HND-05 และ migration 0003 ของ 002
  - task ที่เริ่มได้ทันที: T-01 (config), T-03 (package + กันบทบาท), T-07 (ที่เก็บรูป), T-08 (BorrowerVerifier)
- ข้อที่ AI ทำต่างจาก plan และต้องให้ทีมยืนยัน: แยกการเพิ่มค่า `action` ของ audit_log ออกจาก migration ของ T-02 ไปไว้ T-26
  (plan ข้อ 3 บอก 1 ไฟล์) เพื่อไม่ให้ทุกเฟสรอ Q-HND-06
- สิ่งที่ AI พบ: plan อ้าง package `audit` ของกลาง แต่ใน repo ยังไม่มี
- AC ที่ทดสอบยาก: AC-HND-03 (Auth API ไม่ตอบ 10 วินาที) ใช้ verifier จำลองที่ค้าง + timeout สั้นใน test,
  T-18 (ส่งมอบพร้อมกัน) ใช้ 2 goroutine กับ PostgreSQL จริง

---

## 2569-11-11 คำสั่ง: /tasks แก้รอบที่ 1 (ไม่กระทบฟีเจอร์อื่น)

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- คำสั่ง: "เช็คว่าที่ทำจะมีอะไรไปกระทบของคนอื่นมั้ย" แล้ว "ลองทำแบบที่ไม่กระทบมาให้หน่อย"
- AI เช็ค branch ของ 001, 002, 003, 004, 006 บน GitHub พบ 8 จุด แก้ใน plan.md (v1.1) และ tasks.md (v1.1) ของ 005 เท่านั้น ไม่แตะไฟล์ของฟีเจอร์อื่นหรือไฟล์ของกลาง
- แก้อะไร เพราะอะไร:
  1. เลิก ALTER `loan` ย้ายข้อมูลการส่งมอบไปตาราง `handover` ของ 005 (T-02, T-06) เพราะ seed และ test ของ 001 INSERT `loan` 4 คอลัมน์จะพัง และไม่ต้องแก้ DM-05
  2. migration 2 ไฟล์ เลขถัดจาก 0005 ของ 003 ไฟล์ action ต้องมีทุกค่าเดิมใน CHECK (T-02, T-26) เพราะ 003 และ 004 ก็แก้ CHECK เดียวกัน
  3. ใช้ `NormalLoanReturnTime` ชื่อเดียวกับ plan ของ 004 แทน `NormalLoanDueHour/Minute` (T-01, T-14) เพราะกฎข้อ 9 ห้ามมีค่าเดียวกัน 2 ที่
  4. audit log ใช้ writer.go ของ 001 ถ้ามี ไม่มีใช้แผนสำรองแบบ 003 (T-26) เพราะ package `audit` ยังไม่มี
  5. สถานะคำขอที่ส่งมอบได้อ่านจาก `HandoverReadyStatuses()` ใน config (T-04) และ 005 ไม่อ่านตารางของ 004 (T-14)
  6. ไฟล์ของกลางเพิ่มต่อท้ายเป็นบล็อกของ 005 ทำหลัง merge 003, test migration อยู่ใน package handover ไม่แตะ `platform/migrate_test.go`, volume รูปเป็น named volume ใหม่
- ผลข้างเคียง: T-02 ไม่ติด Q แล้ว (เหลือรอ 002 merge 0003 และจองเลข) task ที่รอ Q ลดจาก 16 เป็น 15
- ยังต้องให้ทีมตัดสิน (AI ไม่ได้ตัดสินแทน): Q-P3 ใหม่ กำหนดคืน 2 ที่ (`approval_decision.return_due_at` ของ 004 กับ `loan.due_at` ของ 005) ตอบพร้อม Q-HND-03,
  Q-HND-05 เหลือคำถามว่าสถานะที่ส่งมอบได้คือ `approved_by_officer` หรือ `approved`

---

## 2569-11-11 pull main (002 spec v3 + CR-003) แล้วตรวจผลกระทบต่อ 005

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- คำสั่ง: "ไป pull main มาหน่อย เพื่อนเพิ่ม cr003 เข้ามา แล้วก็ดูงานของกูว่ามีอะไรกระทบมั้ย แล้วรายงานมา" แล้ว "อัปเดตตามนี้"
- merge: conflict 3 ไฟล์ (prompt-log.md, CHANGELOG.md, constraints.md) AI เก็บทั้งสองฝั่ง ไม่ลบบันทึกของใคร
  ตาราง "ใครอ้างใช้บ้าง" รวมแถวให้มีทั้ง 002 และ 005 commit `867412c 005: merge main (002 spec v3 + CR-003)`
- ผลกระทบที่ AI พบ 7 จุด แก้ใน spec v2.4, plan v1.2, tasks v1.2:
  1. Q-HND-05 ปิดบางส่วนจาก data-model v9: ส่งมอบได้เมื่อ `approved` แล้วตั้งเป็น `handed_over`
  2. Q-SH-04 ปิด + CR-003 ให้ 005 แปลงวันสิ้นสุดการยืมของค้างคืนเป็นกำหนดคืน -> แก้ FR-HND-05 และ SRS 3.1.5
  3. 002 รับคำขอค้างคืนแล้ว -> แก้ Scope และ Q-HND-07 (ยังขึ้นกับ UC-14 ของ 004)
  4. FR-BRQ-06 ของ 002 อยู่ใน prototype -> แก้ข้อความ Q-HND-03 (ถ้ากำหนดคืน 16:00 วันเดียวกัน ผู้ยืมจะถูกปฏิเสธคำขอใหม่ทันที)
  5. Q-SH-06 ปิดบางส่วน (004 ปฏิเสธคำขออื่นอัตโนมัติ) -> Given ของ AC-HND-13 ไม่เกิดจริง ใส่หมายเหตุ ไม่แก้ตัว AC
  6. 002 เพิ่ม action `request` -> migration action ของ 005 ต้องมีค่านี้
  7. (แจ้งเฉย ๆ) ชื่อในใบคำขอของ 002 เป็นอีกทางเลือกของ Q-HND-13, ที่เก็บไฟล์แนบของ 002 ควรตกลงร่วมกับ volume รูป
- tasks: T-04, 05, 06, 16, 20 ไม่รอ Q แล้ว task ที่รอ Q ลดจาก 15 เป็น 10 เฟส 1 ไม่ติด Q (รอแค่ 002 merge 0003)
- ยังต้องให้ทีมตัดสิน (AI ไม่ได้ตัดสินแทน): แก้หรือยกเลิก AC-HND-13 (ต้องทำ CR), Q-HND-07 ต้องถาม 004 เรื่อง UC-14

---

## 2569-11-11 pull main รอบ 2 (002 plan/tasks/migration 0003 + CR-004) แล้วแก้ 005 ตามผลกระทบ

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (แก้ของกลาง `specs/000-shared/constraints.md` แถว CON-TECH-03 ในตาราง "ใครอ้างใช้บ้าง" เพราะ 005)
- เครื่องมือ: Claude Code
- คำสั่ง: "ตอนนี้เพื่อน push main ไปใหม่ละ pull ลงมาแล้ว recheck สรุปอีกครั้งมาหน่อย" / "ว่ามันกระทบของเราตรงไหนบ้าง" แล้ว "แก้ตามนี้"
- merge: conflict 3 ไฟล์ (prompt-log.md, CHANGELOG.md, constraints.md) เก็บทั้งสองฝั่ง commit `32c0700`
- ผลกระทบที่ AI พบและแก้ (spec v2.5, plan v1.3, tasks v1.3):
  1. CON-TECH-03 (CR-004): เพิ่มในบรรทัด "ใช้ร่วม" error contract ใช้ `{code, message}` เปลี่ยน `CARD_NOT_CHECKED`, `PHOTO_INVALID` เป็น `VALIDATION_FAILED`
     รหัส `FORBIDDEN_ROLE`, `REQUEST_NOT_APPROVED`, `BORROWER_MISMATCH` ไม่อยู่ในรายการกลาง -> Q-P4 ใหม่ AI ไม่ได้เพิ่มแถวในตารางรหัสของกลางเอง
  2. migration 0003 merge แล้ว: Q-HND-05 ปิด (`requester_id`, `item_id`, `end_at`) T-02 ไม่รอ 002
  3. CHECK ของ action ใช้ชื่อ `audit_log_action_check` ตามที่ 0003 ตกลงกับ 003
  4. ค่า 16:00 ใช้ `LoanNormalReturnHour` ของ 002 แทน `NormalLoanReturnTime` ของ 004 (004 spec v1 ไม่มีแล้ว)
  5. 004 spec v1 (force-push): AC-APV-08 ไม่มีแล้ว เปลี่ยนไปอ้าง FR-APV-01 (ค), UC-14 ออกจาก scope -> แก้ Q-HND-07, Q-P3 หมดไป
  6. Q-HND-13: เพิ่มทางเลือกชื่อที่ผู้ขอพิมพ์ในคำขอของ 002
- ผลข้างเคียง: task ที่รอ Q เพิ่มจาก 10 เป็น 13 (T-03, T-16, T-20 รอ Q-P4) T-30 รอแค่ Q-HND-07
- ยังต้องให้ทีมตัดสิน: Q-P4 ลงทะเบียนรหัส error 3 ตัว (แจ้งกลุ่มก่อนแก้ constraints.md), AC-HND-13 แก้หรือยกเลิก, Q-HND-07

---

## 2569-11-11 กำหนดเลข migration และคำสั่ง /implement T-08 (หยุดก่อนเริ่ม)

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (แก้ของกลาง `specs/000-shared/data-model.md` ตาราง "เลข migration" เพราะ 005 จองเลข)
- เครื่องมือ: Claude Code
- คำสั่ง: "migration เอาเป็น 6 กับ 7 แล้วก็ /implement task ที่จะไม่กระทบกับคนอื่นมาก่อนได้เลย"
- สิ่งที่แก้: plan.md และ tasks.md ใช้ `0006_handover.up.sql`, `0007_audit_log_handover_action.up.sql` แทน `00NN`
  data-model.md เพิ่ม 2 แถวในตาราง "เลข migration" (ระบุว่ารอกลุ่มยืนยัน) T-02 ไม่ต้องรอจองเลขแล้ว
- /implement: AI เลือก T-08 (แตะแค่ไฟล์ใน `backend/internal/handover/`) แต่ตรวจก่อนเริ่มไม่ผ่านข้อ 2 ของคำสั่ง
  เพราะช่อง "ต้องทำหลัง" ของ T-08 คือ T-01 ซึ่งยังไม่เสร็จ AI หยุดและถามทีม ไม่ได้สร้างโค้ด
- ทุก task ที่ไม่กระทบคนอื่นต้องทำหลัง T-01 (หรือ task ที่รอ Q) ส่วน T-01 แตะ `config.go`, `config.js` ของกลาง (เพิ่มบล็อกต่อท้าย)

---

## 2569-11-11 คำสั่ง: /implement T-01

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (แก้ของกลาง `config.go`, `config.js` เพราะ 005 เพิ่มบล็อกต่อท้ายเท่านั้น)
- เครื่องมือ: Claude Code
- คำสั่ง: "ทีละ task แล้ว commit `00N: ...` ทุกขั้น" (หลังเลือกทำ T-01 ก่อน T-08)
- ไฟล์ที่แก้: `backend/internal/config/config.go` (บล็อก 005: `BorrowerAuthTimeout`, `LoanNormalReturnHour`, `HandoverMinPhotos`, `HandoverPhotoMaxBytes`, `HandoverPhotoTypes()`, `PhotoDir()` และ import `time`)
  `frontend/src/config.js` (บล็อก 005: `HANDOVER_MIN_PHOTOS`, `HANDOVER_PHOTO_MAX_BYTES`, `HANDOVER_PHOTO_TYPES`)
- ผล test: `go test ./...` ผ่านทุก package (auth, equipment, platform ใช้ PostgreSQL ใน testcontainers) `npx vitest run` 10 ไฟล์ 17 test ผ่าน (ต้อง `npm ci` ก่อน)
- ไม่ได้สร้าง config_test.go เพราะไม่อยู่ในช่อง "ไฟล์ที่แตะ" และ branch ของ 003 มีไฟล์ชื่อนี้อยู่แล้ว (จะชนตอน merge)
- สิ่งที่ AI เกือบต้องเดา: ค่าเริ่มต้นของ `PHOTO_DIR` (plan บอกแค่ `/data/photos/handover/`) ใช้ `/data/photos` ตามที่ plan v1.1 เขียนว่า volume mount ที่ `/data/photos`
  และใช้ชื่อ `LoanNormalReturnHour` ตาม plan ของ 002 ที่ยังไม่ได้ใส่ในโค้ด ต้องแจ้งเจ้าของ 002 ให้ใช้ตัวนี้

---

## 2569-11-11 คำสั่ง: /implement T-08

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่สร้าง: `backend/internal/handover/verifier.go` (interface `BorrowerVerifier`, `MockBorrowerVerifier`, `ErrBorrowerNotVerified`), `backend/internal/handover/verifier_test.go`
- ผล test: `TestVerifier_Mock` ผ่าน (2 subtest: token ถูกต้องคืนรหัสผู้ใช้ / token ว่างหรือไม่รู้จักคืน ErrBorrowerNotVerified)
- สิ่งที่ AI เกือบต้องเดาแต่ไม่ทำ: ไม่ตรวจว่าผู้ยืมต้องเป็นบทบาทนักศึกษา (spec ไม่บอก), ไม่คืนชื่อ-นามสกุล (ที่มารอ Q-HND-13),
  credential เป็น token จำลองเท่านั้น (รูปแบบจริงรอ Q-HND-11) ไม่ใส่ timeout (เป็นงานของ T-13)

---

## 2569-11-11 คำสั่ง: /implement T-02

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (เพิ่มไฟล์ใหม่ใน `backend/db/migrations/` และ `backend/db/seed/` ซึ่งเป็นโฟลเดอร์ของกลาง ไม่ได้แก้ไฟล์ของใคร)
- เครื่องมือ: Claude Code
- ไฟล์ที่สร้าง: `backend/db/migrations/0006_handover.up.sql` (ตาราง `handover`, `handover_photo` ไม่ ALTER `loan`),
  `backend/internal/handover/fixture_test.go` (testcontainers แบบเดียวกับ 001 + helper เพิ่มคำขอและการยืม),
  `backend/internal/handover/migration_test.go`, `backend/db/seed/handover_dev.sql` (ข้อมูลลองระบบ ทีมเลือกเอง)
- ผล test: `TestMigration_HandoverIdentityMethod` (4 กรณี), `TestMigration_HandoverPhotoContentType` (4 กรณี), `TestMigration_HandoverNoCardNumberColumn` ผ่าน
  `go test ./...` ทุก package ผ่าน (test ของ 001 ใน equipment และ platform ยังผ่าน)
  `handover_dev.sql` ลองโหลดใน container ของ test ชั่วคราวแล้ว ได้คำขอ approved 2, handed_over 1, แถว handover 1 (ลบไฟล์ test ชั่วคราวแล้ว)
- สิ่งที่ AI ตัดสินในขอบเขตที่ plan ให้: `card_checked` ใส่ CHECK ให้เป็น true เสมอ (ASM-HND-01 ตรวจบัตรทุกครั้ง),
  ชื่อ-นามสกุลเป็น NOT NULL เพราะ FR-HND-03 บังคับบันทึก (ที่มายังรอ Q-HND-13), ไม่เขียน down migration เพราะไม่อยู่ในช่อง "ไฟล์ที่แตะ"
  ข้อมูลลองระบบเพิ่มรุ่นและอุปกรณ์ชุดใหม่ (HO-TEST-*) ไม่แตะอุปกรณ์ของ seed.sql ของ 001
- ต้องแจ้งทีม: เลข 0006 มาก่อน 0004, 0005 ของ 003 ที่ยังไม่ merge ฐานข้อมูลในเครื่องที่รัน 0006 ไปแล้ว golang-migrate จะไม่รัน 0004, 0005 ย้อนหลัง
  ต้องล้างฐานข้อมูลในเครื่อง (docker compose down -v) หลัง 003 merge

---

## 2569-11-11 คำสั่ง: /implement T-06

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (แก้ของกลาง `config.go` เพิ่ม `HandoverReadyStatuses()` ในบล็อกของ 005 ตามกติกา /implement ที่ให้เพิ่มค่าที่ตั้งได้ได้)
- เครื่องมือ: Claude Code
- คำสั่ง: "ทำ T-06 ต่อได้เลย และอันที่กระทบคนอื่นสำคัญ ๆ รายงานด้วย"
- ไฟล์ที่สร้าง/แก้: `backend/internal/handover/repository.go` (`CreateHandover` transaction เดียว), `backend/internal/handover/repository_test.go`,
  `backend/internal/config/config.go` (`HandoverReadyStatuses()` = approved ซึ่ง T-04 ระบุไว้ ย้ายมาเพิ่มตอนนี้เพราะ T-06 ใช้ก่อน)
- ผล test: `TestRepository_CreateHandover` ผ่าน 4 กรณี (สำเร็จ / อุปกรณ์ไม่ว่าง / คำขอไม่พร้อม / ไม่มีคำขอ) `go test ./...` ทุก package ผ่าน
- สิ่งที่ AI ตัดสินในขอบเขตที่ plan ให้: ลำดับล็อกตาม plan ข้อ 4 (อุปกรณ์ FOR UPDATE ก่อน แล้วคำขอ FOR UPDATE),
  `loan.loan_type` เอาจากคำขอ, `loan.borrower_id` = รหัสผู้ใช้ที่ยืนยันแล้ว, รูปยังไม่บันทึก (T-15)
  error ของ repository (`ErrItemNotAvailable`, `ErrRequestNotReady`) ยังไม่ผูกรหัส HTTP (T-10, T-16, T-17)
- ผลกระทบต่อฟีเจอร์อื่นที่ต้องแจ้ง (เปลี่ยนข้อมูล ไม่ได้แก้ไฟล์ของใคร):
  001 เห็นอุปกรณ์เป็น "ถูกยืม" และกำหนดคืนจาก Loan ที่ 005 สร้าง / 002 คำขอ `handed_over` หลุดจาก unique index ผู้ขอส่งคำขอชิ้นเดิมได้อีก
  003 คำขอ `handed_over` ยกเลิกไม่ได้ / 004 ต้องล็อกแถวอุปกรณ์หรือคำขอลำดับเดียวกันถ้ามี transaction ที่แตะ 2 ตารางนี้ (กัน deadlock)
  006 ปิด Loan ที่ 005 สร้าง ต้องรู้ว่ามีแถว `handover` ผูก (FK จาก handover ไป loan ทำให้ลบ Loan ไม่ได้)

---

## 2569-11-11 คำสั่ง: /implement T-18

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/repository_test.go` (เพิ่ม `TestHandover_ConcurrentSameItem`)
- ผล test: ผ่าน 5 รอบติด (`-count=5`) สำเร็จ 1 ได้ ErrItemNotAvailable 1 มี Loan 1 แถว `go test ./...` ทุก package ผ่าน
- รอบแรกไม่ผ่าน: ปิด pool ก่อน cleanup ลบข้อมูล ทำให้ข้อมูลค้างและรอบถัดไปชน asset_code ซ้ำ แก้ให้ปิด pool ใน t.Cleanup (ทำงานหลังการลบ)
- test นี้ใช้ pool และ commit จริง (transaction ที่ยกเลิกทิ้งทดสอบการแย่งกันข้าม connection ไม่ได้) แล้วลบข้อมูลที่สร้างตอนจบ
- ข้อมูลตั้งต้นใช้คำขออนุมัติแล้ว 2 รายการของชิ้นเดียวกัน ซึ่งในระบบจริงไม่เกิด (004 FR-APV-01 ค) ตั้งเองเพื่อบังคับให้แย่งล็อกกัน

---

## 2569-11-11 คำสั่ง: /implement T-07

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- คำสั่ง: "ทำไปยาว ๆ เลย"
- ไฟล์ที่สร้าง: `backend/internal/handover/photostore.go` (`PhotoStore.Save`, `RemoveLoan`, `DefaultPhotoRoot`), `backend/internal/handover/photostore_test.go`
- ผล test: `TestPhotoStore_SaveAndRemove` ผ่าน (เขียน jpg/png ลง `<root>/<loanId>/`, ชนิดอื่นไม่เขียน, ลำดับซ้ำไม่เขียนทับ, ลบทั้งโฟลเดอร์ได้และลบซ้ำไม่ error) `go test ./...` ผ่าน
- **ไม่ได้แก้ `deploy/docker-compose.yml`** แม้อยู่ในช่อง "ไฟล์ที่แตะ" เพราะไฟล์นี้ตอนนี้มีแค่ service `db` ไม่มี service ของหลังบ้าน
  volume `photos` จึงไม่มีที่ mount (หลังบ้านรันด้วย `go run` เก็บรูปที่ `PHOTO_DIR` ค่าเริ่มต้น `/data/photos`)
  การเพิ่ม volume ที่ไม่มีใครใช้คือการแตะไฟล์ร่วมโดยไม่จำเป็น ต้องเพิ่มเมื่อทีมเพิ่ม service หลังบ้านใน compose (ตามหมายเหตุ HTTPS ของ 001)
- ตั้งชื่อไฟล์เองเป็น `<ลำดับ>.jpg|.png` ไม่ใช้ชื่อไฟล์ที่ผู้ใช้ส่งมา (กัน path traversal) การตรวจชนิดจากเนื้อไฟล์เป็นงานของ T-15

---

## 2569-11-11 ตอบ Q-HND-11, Q-HND-13, Q-P4

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (แก้ของกลาง `specs/000-shared/constraints.md` เพิ่มรหัส error 3 ตัวใน CON-TECH-03 เพราะ 005)
- เครื่องมือ: Claude Code
- คำตอบรอบแรกของทีม: "q-hnd-11 ชื่อนามสกุลบัตรนักศึกหรือบัตรประชาชน / q-hnd-13 บัตรนักศึกษาหรือบัตรประชาชน / q-p4 อธิบายเพิ่มเติมหน่อยคืออะไร"
- AI ชี้ว่าคำตอบ Q-HND-11 เป็นเรื่องตรวจบัตร ไม่ได้ตอบวิธียืนยันผ่าน Auth API (FR-HND-02) เสนอ A ล็อกอินบนเครื่องเจ้าหน้าที่ / B เจ้าหน้าที่พิมพ์รหัสให้ระบบค้น / C ตรวจบัตรอย่างเดียว (ต้องทำ CR)
  และอธิบาย Q-P4 (CON-TECH-03 ของ CR-004) เสนอ 1 ลงทะเบียนรหัสใหม่ 3 ตัว / 2 ใช้รหัสเดิมที่ความหมายไม่ตรง
- คำตอบสุดท้าย: "Q-HND-11 a / Q-P4 เอาทางเลือกที่ 1" (Q-HND-13 ใช้คำตอบรอบแรก ซึ่งตรงกับแบบ A ที่ให้เจ้าหน้าที่พิมพ์ชื่อจากบัตร)
- สิ่งที่แก้: spec Draft v2.6 (ASM-HND-09, ASM-HND-10, FR-HND-02, 03, Clarifications), plan v1.4, tasks v1.4, constraints.md v5.1, SRS 3.1.5, CHANGELOG, README
- ระหว่างแก้ AI ลบบรรทัด Q-HND-12 [ยกเลิก] ผิดตัว แล้วใส่คืนจาก commit ก่อนหน้า (ID ที่ยกเลิกต้องคงไว้)
- สิ่งที่ AI ตัดสินในขอบเขต plan แต่ทีมควรตรวจ: รหัสนักศึกษาที่บันทึก = รหัสผู้ใช้จากการล็อกอิน (บัตรประชาชนไม่มีรหัสนักศึกษา)
  และหลังบ้านยืนยันผู้ยืมซ้ำเองตอนบันทึก ไม่เชื่อ `identityMethod` ที่หน้าจอส่ง (กันการปลอมผล)

---

## 2569-11-11 คำสั่ง: /implement T-03

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (แก้ของกลาง `backend/cmd/server/main.go` เพิ่ม 2 บรรทัดผูก route ของ 005 และ import)
- เครื่องมือ: Claude Code
- ไฟล์ที่สร้าง/แก้: `backend/internal/handover/handler.go` (3 route + `requireStaff`), `errors.go` (รหัสและข้อความตาม CON-TECH-03), `handler_test.go`, `backend/cmd/server/main.go`
- ผล test: `TestHandover_NonStaffForbidden` ผ่าน 3 route (นักศึกษาได้ 403 `{code: FORBIDDEN_ROLE, message}` / ไม่มี token ได้ 401) `go test ./...` ผ่าน
- สิ่งที่ AI ตัดสินในขอบเขตที่ plan ให้: path หลังบ้านเป็น `/v1/handover/` เพราะ proxy ตัด `/api` แบบเดียวกับ 001,
  ข้อความภาษาไทยของ error ทีมเลือกเอง ไม่ได้มาจาก spec, route ของ T-04, T-09, T-10 ตอบ 501 ชั่วคราวจนกว่า task นั้นเสร็จ
  `NewHandler(db, provider)` รับ provider ไว้ใช้ยืนยันผู้ยืม (T-09) เพื่อไม่ต้องแก้ main.go อีกใน task ถัดไป
- ข้อสังเกตถึง 001: `auth.Middleware` ของ 001 ตอบ error เป็น `{code}` ไม่มี message ยังไม่ตรง CON-TECH-03 ทั้งหมด (CR-004 บอกว่า 001 ควรแก้รอบถัดไป)

---

## 2569-11-11 คำสั่ง: /implement T-04

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/repository.go` (`ListReadyRequests`, `ReadyRequestDTO`), `handler.go` (`listRequests`), `handler_test.go`
  (`HandoverReadyStatuses()` ใน config.go เพิ่มไปแล้วตอน T-06)
- ผล test: `TestAC_HND_14_ListApprovedOnly` ผ่าน (อนุมัติแล้ว 2 รายการ คนละชิ้น ยังไม่อนุมัติ 1 ส่งมอบแล้ว 1 ได้กลับมา 2 รายการพร้อม assetCode) `go test ./...` ผ่าน
- สิ่งที่ AI ตัดสินในขอบเขตที่ plan ให้: ตอบเป็น JSON array ตาม plan ข้อ 4 (001 ใช้ `{items, total}` แต่ plan ของ 005 เขียนเป็น array),
  เรียงตามเวลาที่ส่งคำขอเก่าก่อน (spec ไม่บอก ทีมเลือกเอง), test ผ่าน HTTP ใช้ pool และ commit จริงแล้วลบข้อมูลทดสอบตอนจบ

---

## 2569-11-11 คำสั่ง: /implement T-09

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/handler.go` (`verifyBorrower`, Handler ถือ `BorrowerVerifier`), `repository.go` (`GetRequester`), `handler_test.go`
- ผล test: `TestAC_HND_02_AuthApiVerifiedAndCardChecked/verify_borrower` ผ่าน (verified ของ student-001)
  `TestVerifyBorrower_OtherResults` ผ่าน (ผู้ยืมคนอื่น = mismatch, credential ใช้ไม่ได้ / ว่าง / ไม่ใช่ JSON = 422 VALIDATION_FAILED, ไม่มีคำขอ = 404)
  `go test ./...` ผ่าน
- สิ่งที่ AI ตัดสินในขอบเขตที่ plan ให้: ล็อกอินไม่สำเร็จไม่อยู่ใน result 3 ค่าของ plan จึงตอบ 422 VALIDATION_FAILED ช่อง `borrowerCredential`
  ไม่มีคำขอ = 404 แบบเดียวกับ 001 (plan ไม่ได้กำหนดรหัส), mismatch ตอบ 200 พร้อม result (ส่วน 409 BORROWER_MISMATCH อยู่ที่ POST handover ใน T-20)

---

## 2569-11-11 คำสั่ง: /implement T-10

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่สร้าง/แก้: `backend/internal/handover/service.go` (`Service.Handover`, `normalDueAt`), `handler.go` (`handover` รับ multipart), `handler_test.go`
- ผล test: `TestAC_HND_02_HandoverRecordsBorrower` ผ่าน 2 กรณี (ตรวจบัตรแล้ว 201 บันทึกชื่อ นามสกุล รหัสนักศึกษา เวลา / ไม่ตรวจบัตร 422 ช่อง cardChecked ไม่มี Loan)
  `TestAC_HND_07_HandoverCreatesLoan` ผ่าน (Loan 1 รายการที่ยังไม่คืน อุปกรณ์ borrowed คำขอ handed_over) `go test ./...` ผ่าน
- ต่างจากช่อง "เสร็จเมื่อ" เล็กน้อย: test ส่วนบันทึกของ AC-HND-02 แยกเป็น `TestAC_HND_02_HandoverRecordsBorrower` (ชื่อยังมี AC ID) เพราะต้องใช้ข้อมูลคนละชุดกับ subtest verify_borrower ของ T-09
- สิ่งที่ AI ตัดสินในขอบเขตที่ plan ให้: หลังบ้านยืนยันผู้ยืมซ้ำด้วย `borrowerCredential` (plan v1.4) รหัสนักศึกษา = รหัสผู้ใช้ (ASM-HND-09)
  ช่องที่ขาดตอบ 422 VALIDATION_FAILED พร้อมชื่อช่อง, กำหนดคืนใช้ `LoanNormalReturnHour` กับ time.Now (นาฬิกาที่ inject ได้อยู่ใน T-14)
- ยังไม่ทำใน task นี้ (เป็นของ task อื่น): รูป (T-15) test ส่งรูป JPG 2 รูปไว้แล้วแต่ยังไม่ถูกเก็บ, เทียบผู้ยืมกับผู้ขอ (T-20),
  timeout 10 วินาที (T-13), REQUEST_NOT_APPROVED (T-16 ตอนนี้คำขอที่ไม่พร้อมได้ 503 ชั่วคราว)

---

## 2569-11-12 คำสั่ง: /implement T-14

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/service.go` (Service รับนาฬิกา `now`), `handler.go` (`newHandler` ใส่ verifier และนาฬิกาได้), `handler_test.go`
- ผล test: `TestAC_HND_05_DueAt1600SameDay` ผ่าน (นาฬิกา 2569-11-12 10:00 Asia/Bangkok ได้เวลาส่งมอบ 10:00 และกำหนดคืน 16:00 วันเดียวกัน ตรวจจากฐานข้อมูล) `go test ./...` ผ่าน
- ค่า 16:00 อ่านจาก `config.LoanNormalReturnHour` เขตเวลาจาก `config.TimeZone` ไม่อ่านตารางของ 004

---

## 2569-11-12 คำสั่ง: /implement T-16

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/handler.go` (แปลง `ErrRequestNotReady` เป็น 409 `REQUEST_NOT_APPROVED`), `handler_test.go`
  ช่อง "ไฟล์ที่แตะ" ระบุ service.go แต่การตรวจสถานะอยู่ใน transaction ของ T-06 แล้ว งานที่เหลือคือแปลง error ซึ่งอยู่ใน handler.go
- ผล test: `TestAC_HND_01_RequestNotApproved` ผ่าน 2 กรณี (requested / handed_over ได้ 409 ไม่มี Loan อุปกรณ์ยัง available สถานะคำขอไม่เปลี่ยน) `go test ./...` ผ่าน

---

## 2569-11-12 คำสั่ง: /implement T-17

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/handler_test.go`
- ผล test: `TestAC_HND_04_ItemNotAvailable` ผ่าน 2 กรณี (unavailable / borrowed ได้ 409 ITEM_NOT_AVAILABLE ไม่มี Loan ใหม่ สถานะอุปกรณ์ไม่เปลี่ยน)

---

## 2569-11-12 คำสั่ง: /implement T-20

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/repository.go` (เทียบผู้ยืมกับ `requester_id` ใน transaction เดียวกับการล็อกคำขอ `ErrBorrowerMismatch`),
  `handler.go` (409 `BORROWER_MISMATCH`), `handler_test.go`, `repository_test.go`
  ช่อง "ไฟล์ที่แตะ" ระบุ service.go แต่ AI ย้ายการเทียบไปไว้ใน transaction เพื่อใช้แถวคำขอที่ล็อกแล้ว ไม่ต้องอ่านคำขอซ้ำ
- ผล test: `TestAC_HND_10_BorrowerMismatch` ผ่าน (นักศึกษา B ล็อกอินรับของคำขอของ A ได้ 409 ไม่มี Loan อุปกรณ์ available คำขอยัง approved)
- test ที่ไม่ผ่านระหว่างทาง: `TestHandover_ConcurrentSameItem` (T-18) เพราะข้อมูลตั้งต้นใช้ผู้ยืม student-001 กับคำขอของ student-002 ด้วย
  แก้ข้อมูลตั้งต้นของ test ให้ผู้ยืมตรงกับผู้ขอของแต่ละคำขอ (ไม่ใช่ test ที่ขึ้นต้นด้วย TestTC_ และไม่ได้ลดสิ่งที่ test ตรวจ) ผ่าน 3 รอบติด
- `go test ./...` ผ่านทุก package

---

## 2569-11-12 คำสั่ง: /implement T-15

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่สร้าง/แก้: `backend/internal/handover/photos.go` (ตรวจจำนวนขั้นต่ำและชนิดจากเนื้อไฟล์ด้วย `http.DetectContentType`, `photoSaver`),
  `service.go` (ตรวจรูปก่อนยืนยันผู้ยืม), `repository.go` (`CreateHandover` รับ `PhotoSaver` เขียนรูปหลังสร้าง Loan บันทึก `handover_photo` ใน transaction เดียวกัน ลบรูปเมื่อขั้นต่อไปล้มเหลว),
  `handler.go` (อ่านช่อง `photos`, 422 VALIDATION_FAILED ช่อง photos), `handler_test.go`, `repository_test.go` (ส่ง nil เป็น PhotoSaver), `fixture_test.go` (ลบโฟลเดอร์รูปของ test ตอนจบ)
- ผล test: `TestAC_HND_06_TwoJpgPhotosSaved` ผ่าน (JPG 2 รูป รูปละ 2 MB ได้ 201 มี 2 แถวผูกกับ Loan ไฟล์อยู่ที่ `<root>/<loanId>/` ขนาดตรง) `go test ./...` ผ่าน
- test เก็บรูปที่โฟลเดอร์ชั่วคราวของระบบ ไม่เขียนลง `/data/photos` ของเครื่อง
- ขนาด 10 MB ยังไม่ตรวจ (T-28 รอ Q-HND-10)

---

## 2569-11-12 คำสั่ง: /implement T-21

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/handler_test.go`
- ผล test: `TestAC_HND_12_InvalidPhotos` ผ่าน 3 กรณี ((ก) 1 รูป, (ข) GIF, (ข) GIF ที่ตั้งชื่อและ Content-Type เป็น jpg) ได้ 422 ช่อง photos ไม่มี Loan อุปกรณ์ available ไม่มีไฟล์ค้าง
- กรณี (ค) 11 MB ยังไม่ทำ (T-28 รอ Q-HND-10)

---

## 2569-11-12 คำสั่ง: /implement T-19

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/handler_test.go` (ไม่ได้แก้ fixture_test.go เพราะ helper ใน handler_test.go พอแล้ว)
- ผล test: `TestAC_HND_13_OtherRequestUntouched` ผ่าน (ส่งมอบ X สำเร็จ สถานะ Y ยัง requested / ตั้ง Y เป็น approved แล้วส่งมอบ Y ได้ 409 ITEM_NOT_AVAILABLE มี Loan 1 รายการ)
- ข้อมูลตั้งต้นตาม AC ตรง ๆ ซึ่งในระบบจริงไม่เกิดแล้ว (Q-SH-06 ปิดบางส่วน) เขียนหมายเหตุใน test ถ้าทีมแก้หรือยกเลิก AC-HND-13 ต้องแก้ test นี้ตาม
- ขั้น "ส่งมอบ Y" ต้องตั้ง Y เป็น approved เองก่อน เพราะ 005 ส่งมอบได้เฉพาะคำขอที่อนุมัติแล้ว (ถ้าไม่ตั้ง จะได้ REQUEST_NOT_APPROVED แทน ITEM_NOT_AVAILABLE ตามที่ AC เขียน)
- gofmt จัดรูปแบบคอมเมนต์ของ test ใหม่ (ย่อหน้าใน doc comment)

---

## 2569-11-12 คำสั่ง: /implement T-23

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/handler_test.go`
- ผล test: `TestAC_HND_09_StudentForbidden` ผ่าน (`mock-student-001` เรียก POST handover ของคำขออนุมัติแล้ว ได้ 403 FORBIDDEN_ROLE ไม่มี Loan อุปกรณ์ available คำขอยัง approved)

---

## 2569-11-12 คำสั่ง: /implement T-24

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `backend/internal/handover/handler_test.go`
- ผล test: `TestAC_HND_11_HandoverAfter1600` ผ่าน (นาฬิกา 16:30 ได้ 201 เวลาส่งมอบ 16:30 อุปกรณ์ borrowed) ไม่ตรวจ `due_at` ตามที่ spec สั่ง (รอ Q-HND-03, T-25)
- ข้อสังเกต: ตอนนี้โค้ดตั้งกำหนดคืนเป็น 16:00 ของวันเดียวกันแม้ส่งมอบหลัง 16:00 (เลยกำหนดทันที) จนกว่าทีมตอบ Q-HND-03

---

## 2569-11-12 คำสั่ง: /implement T-05

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (แก้ของกลาง `frontend/src/App.jsx` และ `frontend/src/api/client.js` เพิ่มต่อท้ายเท่านั้น)
- เครื่องมือ: Claude Code
- ไฟล์ที่สร้าง/แก้: `frontend/src/pages/Handover/HandoverList.jsx`, `frontend/src/pages/Handover/HandoverPage.jsx` (ไม่อยู่ในช่อง "ไฟล์ที่แตะ" เพิ่มเป็นหน้าที่รวมรายการกับฟอร์มของ T-11 เพื่อให้ App.jsx ไม่ต้องแก้อีก),
  `frontend/src/api/client.js` (บล็อก 005 `handoverApi` และ helper ที่อ่าน `{code, fields}` ตาม CON-TECH-03 ไม่แก้ `api` เดิมของ 001),
  `frontend/src/App.jsx` (เปิดหน้าส่งมอบเมื่อ path ขึ้นต้นด้วย `/handover` ไม่เพิ่มปุ่มหรือเมนู), `frontend/src/__tests__/AC-HND-14.test.jsx`
- ผล test: `AC-HND-14.test.jsx` ผ่าน (client จำลองคืน 2 คำขอ ตารางแสดง 2 แถว) `npx vitest run` 11 ไฟล์ 18 test ผ่าน `vite build` ผ่าน
- ไม่มี mockup ทีมออกแบบเอง หัวคอลัมน์และข้อความ "เลือก" "ไม่มีคำขอที่รอส่งมอบ" ทีมเลือกเอง ไม่ได้มาจาก spec
- ไม่มี router ใน template จึงใช้ path ของเบราว์เซอร์ตรง ๆ แทนการเพิ่ม library (plan ไม่ได้ระบุ)
- หมายเหตุ: App.jsx ใช้ CRLF แก้ด้วย Edit

---

## 2569-11-12 คำสั่ง: /implement T-11

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่สร้าง/แก้: `frontend/src/pages/Handover/HandoverForm.jsx`, `HandoverPage.jsx` (รายการ -> ฟอร์ม), `frontend/src/api/client.js` (เก็บ message จากหลังบ้านไว้แสดง), `frontend/src/__tests__/AC-HND-02.test.jsx`
- ผล test: `AC-HND-02.test.jsx` ผ่าน 2 test (ไม่ติ๊กตรวจบัตร ปุ่มใช้ไม่ได้และไม่ส่ง / ติ๊กแล้วส่งชื่อ นามสกุล cardChecked=true รูป 2 ไฟล์ และแสดงเวลาที่ส่งมอบ) `npx vitest run` 12 ไฟล์ 20 test ผ่าน
- ไม่มี mockup ข้อความบนฟอร์มทีมเลือกเอง ไม่ได้มาจาก spec ข้อผิดพลาดแสดง message ที่หลังบ้านส่งมา (CON-TECH-03)
- ช่องล็อกอินของผู้ยืม (borrowerCredential) ยังไม่มี เป็นงานของ T-12 ระหว่างนี้หลังบ้านจะตอบ 422 ถ้ายังไม่ได้ทำ T-12

---

## 2569-11-12 คำสั่ง: /implement T-12

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่แก้: `frontend/src/pages/Handover/HandoverForm.jsx` (ส่วน "ผู้ยืมล็อกอินบัญชีมหาวิทยาลัย" ปุ่มยืนยันตัวตนผู้ยืม แสดงผล verified / mismatch / unavailable),
  `frontend/src/__tests__/AC-HND-02.test.jsx` (2 test เดิมเพิ่มขั้นยืนยันผู้ยืม + test ใหม่กรณียังไม่ยืนยันและ mismatch)
- ผล test: `AC-HND-02.test.jsx` 3 test ผ่าน `npx vitest run` 12 ไฟล์ 21 test ผ่าน `vite build` ผ่าน
- ปุ่มบันทึกใช้ได้เมื่อผลเป็น verified หรือ unavailable และติ๊กตรวจบัตรแล้ว หลังบ้านยืนยันผู้ยืมซ้ำเองเสมอ (plan v1.4)
- ช่องบัญชีผู้ยืมเป็นช่องรหัสผ่าน prototype กรอก token จำลอง (เช่น mock-student-001) จนกว่าจะรู้โปรโตคอลของ Auth API (Q-SH-03)
- ข้อความบนหน้าจอทีมเลือกเอง ไม่ได้มาจาก spec

---

## 2569-11-12 คำสั่ง: /implement T-22

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- ไฟล์ที่สร้าง/แก้: `frontend/src/pages/Handover/HandoverForm.jsx` (ตรวจจำนวนและชนิดรูปจาก `HANDOVER_MIN_PHOTOS`, `HANDOVER_PHOTO_TYPES` ก่อนส่ง), `frontend/src/__tests__/AC-HND-12.test.jsx`
- ผล test: `AC-HND-12.test.jsx` ผ่าน 2 กรณี ((ก) 1 รูป, (ข) GIF ไม่เรียก client และแสดงข้อความแจ้ง) `npx vitest run` 13 ไฟล์ผ่านทั้งหมด `vite build` ผ่าน
- ข้อความแจ้งเรียบเรียงจาก FR-HND-06 (อย่างน้อย 2 รูป JPG หรือ PNG) ตัวเลขอ่านจาก config.js ขนาดไฟล์ยังไม่ตรวจ (T-28 รอ Q-HND-10)

---

## 2569-11-12 แก้ให้ PR ของ 005 ไม่กระทบฟีเจอร์อื่นตอน merge

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover (คืนไฟล์ของกลาง `App.jsx`, `api/client.js` เป็นแบบเดียวกับ main แก้ `main.jsx` และ `config.js` ของกลางแทน)
- เครื่องมือ: Claude Code
- คำสั่ง: "ถาม ถ้าจะ pr ไป main ตอนนี้จะกระทบใคร และอะไรมั้ย มีอะไรบ้างที่ต้อง merge" แล้ว "แก้ไขมาให้ไม่กระทบคนอื่นก่อน"
- สิ่งที่ AI พบจากการลอง merge จำลองกับ branch ของทุกคน:
  1. 002 ประกาศ `const LoanNormalReturnHour` ใน config.go เหมือน 005 git merge ได้เองแต่ build พัง (ประกาศซ้ำ)
  2. 002 commit `App.jsx`, `client.js`, `config.js`, `migrate.go` เป็น CRLF ทั้งไฟล์ ใครแก้ไฟล์เหล่านี้จะชนกับ 002 ทั้งไฟล์ (ไม่ใช่ปัญหาของ 005)
  3. 003 แก้ `App.jsx` (ปุ่มเลือกหน้า) และต่อท้าย `config.js` ชนกับ 005
- สิ่งที่แก้:
  1. config.go เปลี่ยนชื่อค่าของ 005 เป็น `HandoverDueHour` ชั่วคราว (ค่าเดียวกับ `LoanNormalReturnHour` ของ 002 ต้องรวมเป็นตัวเดียวหลัง 002 merge)
  2. ย้าย `handoverApi` จาก `api/client.js` ไปไฟล์ใหม่ `frontend/src/api/handover.js` (ใช้ BASE และ ApiError ของ client.js) แล้วคืน client.js เป็นแบบ main
  3. คืน `App.jsx` เป็นแบบ main (รอบแรกย้ายไป `main.jsx` แต่ 002 เปลี่ยนท้ายบรรทัดของ main.jsx ด้วย จึงย้ายอีกครั้ง)
     เปิดหน้าส่งมอบจากไฟล์ใหม่ `frontend/handover.html` + `src/handover-main.jsx` (Vite เปิดได้ที่ /handover.html โดยไม่ต้องแก้ vite.config.js) ไม่แก้ไฟล์เดิมของใคร
  4. ย้ายบล็อก 005 ใน `frontend/src/config.js` ไปไว้ต้นไฟล์ ไม่ชนกับ 002, 003 ที่ต่อท้าย
- ผล test: `go test ./...` ผ่านทุก package `npx vitest run` 13 ไฟล์ 23 test ผ่าน `vite build` ผ่าน
- ชนที่เหลือและหลีกไม่ได้: เอกสารที่ต่อท้าย (prompt-log, CHANGELOG, constraints.md, data-model.md) และ config.js กับ 002 เพราะ 002 เปลี่ยนท้ายบรรทัดทั้งไฟล์
- รอบที่ 2 (เอกสารของกลาง) หลังลอง merge จำลองซ้ำ:
  - 004 ลงทะเบียน `FORBIDDEN_ROLE` ไว้ใน branch ของตัวเองด้วย (รหัสเดียวกัน) AI ใช้แถวที่ข้อความตรงกับของ 004 ทุกตัวอักษรในตำแหน่งเดียวกัน git จึง merge ได้โดยไม่ชนและไม่เกิดแถวซ้ำ
    แถว `REQUEST_NOT_APPROVED`, `BORROWER_MISMATCH` ย้ายไปไว้ใต้ `ITEM_NOT_AVAILABLE` ไม่ติดกับจุดที่ 004 เพิ่ม
  - คืน `data-model.md` เป็นแบบ main เพราะ 003 แก้ตาราง "เลข migration" แถวเดียวกัน (003 เขียนว่า "0006 ขึ้นไป ฟีเจอร์ถัดไป") เลข 0006, 0007 ของ 005 บันทึกใน plan.md และต้องแจ้งกลุ่ม ลงตารางหลัง 003 merge

---

## 2569-11-12 merge main (PR #9 ของ 002) เข้า 005 และใช้ค่า 16:00 ของ 002

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- คำสั่ง: "ตรวจสอบ commit เมื่อ 7 นาทีที่แล้วของ branch main อีกครั้งหน่อย เพื่อนเพิ่งอัพเข้ามา" แล้ว "แก้ไขที่เราก่อน"
- main ได้ฟีเจอร์ 002 (PR #9) มี `LoanNormalReturnHour = 16` ใน config.go แล้ว
- conflict 2 ไฟล์: `frontend/src/config.js` (ใช้ของ main + บล็อก 005 ต้นไฟล์ 10 บรรทัดแรกที่ 002 ยังเป็น CRLF ถูกแปลงกลับเป็น LF ตาม autocrlf ซึ่งตรงกับก่อน 002)
  และ `prompt-log.md` (เก็บทั้งสองฝั่ง ของ main ก่อน)
- ลบ `HandoverDueHour` ของ 005 ใช้ `config.LoanNormalReturnHour` ของ 002 แทน (กฎข้อ 9 ค่าเดียวที่เดียว) แก้ plan ข้อ 5 และ tasks (หัวไฟล์, T-01, T-14)
- ผล test: `go test ./...` ผ่านทุก package รวม borrowrequest และ config_test ของ 002 `npx vitest run` ผ่าน

---

## 2569-11-12 ตรวจก่อน push และเปิด PR ไป main

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- คำสั่ง: "recheck ให้ชัวว่า push + pr แล้ว merge ไปแล้วงานของคนอื่นไม่พัง หรืองานเราไม่ได้ไปทับของคนอื่น จัดการมาได้เลย"
- ตรวจ: branch มี main ล่าสุด (9e16c56) แล้ว ลอง merge จำลองกับ main ไม่ชน ไฟล์นอก 005 ที่แตะ 8 ไฟล์ ทุกบรรทัดที่เปลี่ยนในไฟล์ของคนอื่นเป็นการเติม
  ข้อความเดิมของ 001, 002 ในตาราง "ใครอ้างใช้บ้าง" ยังอยู่ครบ prompt-log ไม่มีบรรทัดถูกลบ
- แก้: ข้อความ Status ของ constraints.md v5.1 ให้ตรงว่าเพิ่ม 2 รหัสใหม่ และ FORBIDDEN_ROLE แถวเดียวกับของ 004
- ผล test: `go test -p 1 ./...` ผ่านทุก package (auth, borrowrequest, config, equipment, handover, platform) `npx vitest run` 23 test ผ่าน `vite build` ผ่าน
  รอบแรกที่รันพร้อมกันหลาย package test ของ equipment ล้มเพราะ testcontainers ต่อ Docker ไม่ได้ชั่วคราว รันซ้ำแล้วผ่าน ไม่ใช่ปัญหาของโค้ด

---

## 2569-11-12 ตรวจชื่อ commit ให้ตรงรูปแบบหลัก

- ผู้สั่ง: Banlu C (660710764 นายบรรลุ ฉิมสิงห์)
- ฟีเจอร์: 005-handover
- เครื่องมือ: Claude Code
- คำสั่ง: "ไป recheck ชื่อ commit ไป push ไปดูของคนอื่น แล้วทำในรูปแบบหลัก" แล้ว "อนุญาต force push"
- รูปแบบหลัก (specs/README.md ขั้น 2-5 และ commit ใน main): `00N: spec vX`, `00N: plan vX`, `00N: tasks vX`, `00N: T-xx done` รายละเอียดเพิ่มต่อด้วย ` - `
- แก้ชื่อ 2 commit (เนื้อหาไฟล์ไม่เปลี่ยน ตรวจด้วย git diff): `005: spec v2.6 ปิด ...` เป็น `005: spec v2.6 - ปิด ...`
  และ `005: spec v2.5, plan v1.3, tasks v1.3` เป็น `005: tasks v1.3 - รวม spec v2.5 และ plan v1.3 ...` (ไม่มี commit tasks v1 แยกตาม README)
- force push ด้วย `--force-with-lease=005-handover:3e24726` ตามที่ทีมอนุญาต (branch ของ 005 คนเดียว ยังไม่มี PR) remote เป็น 1876206
- ไม่ได้สร้าง commit `T-xx reviewed` เพราะต้องให้คนในทีมตรวจจริงก่อน

---

## 2569-11-12 /implement T-08 (002-borrow-request) — ปิดเฟส 4

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 4 AC-BRQ-13 กฎเวลาคืนของการยืมปกติ)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-08`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/validate.go` (ตรงกับช่องไฟล์ที่แตะ) เพิ่ม loanWindowProblems และตัวแปร appLoc
- `backend/internal/config/config.go` เพิ่ม LoanTypeNormal, LoanTypeOvernight
  (คำสั่ง /implement อนุญาตให้เพิ่ม config ได้โดยไม่ต้องถาม เมื่อเป็นค่าที่มาจาก spec)
- `backend/internal/borrowrequest/handler_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** ต่อ test ท้ายไฟล์เดิม

### การตัดสินใจ
- เทียบ "วันเดียวกัน" ด้วยเขตเวลา Asia/Bangkok (glossary) ไม่ใช่เวลาของเครื่องที่รัน
  ถ้าใช้ UTC คำขอช่วงเย็นของไทยจะกลายเป็นคนละวันทันที
- "ภายใน 16:00 น." ตีความว่าคืนพอดี 16:00 ผ่านได้ เกินไป 1 นาทีไม่ผ่าน มี test ครอบทั้งสองกรณี
- กฎนี้ไม่ใช้กับการยืมค้างคืน ตามที่ DOM-LOAN-01 เขียนว่าใช้กับยืมปกติเท่านั้น มี test ยืนยัน
- ไม่ได้ตรวจว่า endAt มาก่อน startAt หรือไม่ เพราะไม่มี FR ข้อไหนสั่งไว้ (ไม่เพิ่มเอง)

### ขอบเขตที่ tasks.md แบ่งไว้แล้ว ไม่ได้ตัดเอง
AC-BRQ-13 เขียนว่าต้อง "แจ้งว่าการยืมแบบปกติต้องคืนภายใน 16:00 น." แต่ข้อความอยู่ใน handler.go
ซึ่งไม่อยู่ในช่องไฟล์ที่แตะของ T-08 — ตรวจแล้วพบว่า **ช่อง "ตรวจด้วย" ของ T-20 ระบุ AC-BRQ-13 ไว้
พร้อมวงเล็บว่า "ส่วนข้อความที่ผู้ใช้เห็น"** จึงเป็นการแบ่งที่ตั้งใจไว้ตั้งแต่ตอนทำ /tasks
ตอนนี้ API ตอบ fields: ["endAt"] ส่วนข้อความยังเป็นข้อความรวม รอ T-20 (ซึ่งรอ Q-BRQ-14)

### ผล test
- `go test ./internal/borrowrequest/ -run "TestAC_BRQ_13|TestLoanWindow" -v`: PASS ทั้ง 2 ตัว
- `go test ./...`: ทุกแพ็กเกจ ok (config รันใหม่ 0.007s ยืนยันว่าค่าที่เพิ่มไม่ทำ test เดิมพัง)

---

## 2569-11-12 /implement T-09 (002-borrow-request) — ปิดเฟส 5

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 5 AC-BRQ-01 ตรวจสถานะอุปกรณ์ซ้ำ)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-09`

### คำถามที่ AI ถาม และคำตอบของทีม
1. ช่อง "เสร็จเมื่อ" ต้องการให้ 409 แนบสถานะล่าสุด แต่ response ประกอบที่ handler.go
   ซึ่งไม่อยู่ในช่องไฟล์ที่แตะของ T-09 จะทำอย่างไร
   -> **ทีมเลือกทางเลือก ก: เพิ่ม handler.go ในช่องไฟล์ที่แตะของ T-09** (ตัวเลือกที่ AI แนะนำ)
   เพราะเป็นงานเดียวกัน แยกไป task อื่นจะขาดตอน บันทึกเหตุผลและวันที่ไว้ใน tasks.md แล้ว

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/repository.go` (ตรงกับช่องไฟล์ที่แตะ) อ่านสถานะล่าสุดเมื่อ INSERT ไม่ได้แถว
- `backend/internal/borrowrequest/service.go` (ตรงกับช่องไฟล์ที่แตะ) เพิ่ม ItemNotAvailableError
- `backend/internal/borrowrequest/handler.go` (เพิ่มในช่องไฟล์ที่แตะตามคำตอบข้อ 1) แนบ status ลง response
- `backend/internal/borrowrequest/handler_test.go` และ `service_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ**
  (test ของ "เสร็จเมื่อ" และการแก้ test เดิมที่ล้ม ดูด้านล่าง)

### test เดิมล้มเพราะการเปลี่ยนแปลงนี้ และวิธีแก้
- `go test ./...` รอบแรก **FAIL ที่ TestCreate_RejectsItemNotAvailable** ซึ่งเขียนไว้ตอน T-04
  โดยใช้ `err != ErrItemNotAvailable` เทียบตรง ๆ พอ error กลายเป็นชนิดที่พกสถานะมาด้วย การเทียบจึงไม่ตรง
- แก้เป็น `errors.Is(err, ...)` ทั้ง 3 จุดในไฟล์นั้น **ไม่ได้ลดเงื่อนไขของ test ให้ผ่าน**
  แต่เปลี่ยนวิธีเทียบให้ถูกหลักของ Go ซึ่งรองรับ error ที่ถูกห่อ (test เดิมเขียนแบบเปราะไว้ตั้งแต่แรก)

### การตัดสินใจ
- `ItemNotAvailableError` มีเมท็อด `Is()` ทำให้ `errors.Is(err, ErrItemNotAvailable)` ยังเป็นจริง
  โค้ดเดิมที่เทียบกับ sentinel จึงไม่พัง
- response ส่ง **ค่าดิบ** `"status":"borrowed"` ไม่ใช่คำไทย เพราะ 001 ก็ส่งค่าดิบแล้วให้ StatusBadge แปลง
  ถ้าหลังบ้านส่งคำไทยจะกลายเป็นมีที่แปล 2 ที่
- รหัสครุภัณฑ์ที่ไม่มีอยู่จริง **ไม่แนบ status** (omitempty) เพื่อไม่เปิดเผยว่ามีรหัสนั้นในระบบไหม มี test คุม

### ผล test
- `go test ./internal/borrowrequest/ -run "TestAC_BRQ_01|TestCreate_Unknown" -v`: PASS
- `go test ./...` รอบแรก FAIL (ตามด้านบน) รอบที่สองหลังแก้: **ทุกแพ็กเกจ ok**

---

## 2569-11-12 /implement T-10 (002-borrow-request) — ปิดเฟส 6

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 6 AC-BRQ-07 ขอบเขตภาควิชา)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-10`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/repository.go` (ตรงกับช่องไฟล์ที่แตะ) เพิ่ม checkDepartment และ ErrForbiddenDepartment
- `backend/internal/borrowrequest/handler.go` (ตรงกับช่องไฟล์ที่แตะ) แปลงเป็น 403 FORBIDDEN_DEPARTMENT
- `backend/internal/borrowrequest/handler_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** ต่อ test ท้ายไฟล์เดิม

### การตัดสินใจ
- ภาควิชาของผู้ขออ่านจาก auth.FromContext เท่านั้น ไม่ใช่ค่าที่ส่งมาใน request
  ถ้าเชื่อค่าจาก request ผู้ขอจะปลอมภาควิชาได้ทันที (IF-IDP-01, DM-08)
- repository อ่าน principal จาก context เอง แทนการเพิ่มพารามิเตอร์ที่ service ส่งมา
  เพราะ service.go ไม่อยู่ในช่องไฟล์ที่แตะของ T-10
- รหัสครุภัณฑ์ที่ไม่มีอยู่จริงตอบเป็น ItemNotAvailable เหมือนเดิม ไม่เปิดเผยว่ามีรหัสนั้นไหม

### สิ่งที่ AI เลือกไปก่อน เพราะคำถามยังไม่มีคำตอบ (ทีมต้องตรวจ)
**Q-BRQ-15 ยังเปิดอยู่** ว่าถ้าคนต่างภาควิชาขอยืมอุปกรณ์ที่ไม่ว่าง จะตอบ 403 หรือ 409 ก่อน
AI วางให้ **ตรวจภาควิชาก่อนเสมอ** พร้อมคอมเมนต์กำกับในโค้ดว่ารอ Q-BRQ-15 (ดู T-21)
เหตุผล: ถ้าตรวจสถานะก่อน คนต่างภาควิชาจะได้ 409 พร้อม "status":"borrowed"
ซึ่งรั่วข้อมูลอุปกรณ์ที่เขาไม่มีสิทธิ์เห็น ขัด FR-BRQ-07 ตรง ๆ
ทางที่เลือกแย่ที่สุดคือบอกน้อยเกินไป ซึ่งแก้ทีหลังได้ถ้าทีมตอบเป็นอย่างอื่น

### งานค้างที่ยังทำไม่ได้
AI เคยบอกว่า T-10 จะได้เปลี่ยนชื่อ missingFields เป็น invalidFields **ซึ่งผิด**
การเปลี่ยนชื่อต้องแก้ทั้ง validate.go (ที่ประกาศ) และ handler.go (ที่เรียก)
แต่ T-10 มีแค่ handler.go ในรายการ ต้องรอ task ที่มีทั้งสองไฟล์

### ผล test
- `go test ./internal/borrowrequest/ -run "TestAC_BRQ_07|TestDepartmentScope" -v`: PASS ทั้ง 2 ตัว
  (TestAC_BRQ_07 ไล่หาคำว่า available / borrowed / unavailable / Notebook / CS-NB-0001 ใน body ต้องไม่เจอสักคำ)
- `go test ./...`: ทุกแพ็กเกจ ok

---

## 2569-11-12 /implement T-11 (002-borrow-request) — ปิดเฟส 7

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 7 AC-BRQ-11 ส่งคำขอซ้ำของตัวเองไม่ได้)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-11`

### คำถามที่ AI ถามก่อนเริ่ม และคำตอบของทีม
ช่อง "เสร็จเมื่อ" ของ T-11 บังคับให้ตอบ 409 `DUPLICATE_REQUEST` แต่การแปลง error เป็นรหัสและสถานะ HTTP
อยู่ใน `switch` ของ handler.go ซึ่ง **ไม่มีในช่องไฟล์ที่แตะ** (มีแค่ service.go, repository.go)
ถ้าแตะแค่ 2 ไฟล์นั้น error ตัวใหม่จะตกลง default กลายเป็น 503 ซึ่งผิดรหัสตาม CON-TECH-03
→ **ทีมตอบ: ทางเลือก ก** เพิ่ม handler.go ในช่องไฟล์ที่แตะของ T-11 พร้อมเหตุผลลงวันที่ แบบเดียวกับ T-09
(แก้ tasks.md แล้ว ไม่ได้แก้ task อื่น)

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/repository.go` (ตรงกับช่องไฟล์ที่แตะ)
  เพิ่ม `ErrDuplicateRequest`, `isDuplicateRequest`, `uniqueOpenRequestIndex`, `uniqueViolationCode`
- `backend/internal/borrowrequest/handler.go` (เพิ่มในช่องไฟล์ที่แตะวันนี้ ตามคำตอบของทีม)
  เพิ่ม case แปลงเป็น 409 `DUPLICATE_REQUEST`
- `backend/internal/borrowrequest/service.go` (ตรงกับช่องไฟล์ที่แตะ)
  **แก้คอมเมนต์เท่านั้น** เขียนกำกับว่า Create คืน ErrDuplicateRequest เมื่อไหร่ ไม่มีโค้ดเพิ่ม
- `backend/internal/borrowrequest/service_test.go`, `handler_test.go`
  **ไม่อยู่ในช่องไฟล์ที่แตะ** (ทุก task ที่ผ่านมาก็ไม่เคยใส่ไฟล์ test ในช่องนี้) ต่อ test ท้ายไฟล์เดิม

### การตัดสินใจ
- **ไม่** เพิ่ม SELECT ไปอ่านคำขอเก่ามาตรวจก่อนเขียน ใช้การดักจับ unique violation (SQLSTATE 23505)
  จาก partial unique index `uq_request_open_per_item` ที่ T-01 สร้างไว้แทน
  เพราะการอ่านมาตรวจก่อนเขียนมีช่องว่างระหว่างอ่านกับเขียน ถ้าสองคำขอมาถึงพร้อมกันจะผ่านไปได้ทั้งคู่
  ส่วนรายการสถานะที่ถือว่า "ยังไม่จบ" ตาม FR-BRQ-05 ข อยู่ในเงื่อนไข WHERE ของ index นั้นที่เดียว
- `isDuplicateRequest` ตรวจ **ชื่อ constraint** ด้วย ไม่ได้ดูแต่รหัส 23505
  เพื่อไม่เหมาเอาว่า unique violation ตัวอื่นในตารางนี้ในอนาคตคือคำขอซ้ำ
- `ErrDuplicateRequest` ประกาศไว้ใน repository.go เพราะ model.go ไม่อยู่ในช่องไฟล์ที่แตะ
  (เหมือน ErrForbiddenDepartment ของ T-10 ย้ายไปรวมกันได้ภายหลัง)
- รหัส `23505` เป็นค่าคงที่ของ PostgreSQL ไม่ใช่ตัวเลขที่มาจาก spec จึงไม่ได้ย้ายไป config.go
  (AGENTS.md กฎข้อ 9 คุมตัวเลขที่มาจาก spec) แต่ประกาศเป็น const พร้อมคอมเมนต์ ไม่ฝังกลางโค้ด
- test ใช้ผู้ขอคนละรหัสกับ test อื่น (`student-011`, `student-012`) เพราะกฎกันซ้ำผูกกับคู่ (ผู้ขอ, อุปกรณ์)
  ถ้าใช้ student-001 ร่วมกับ test อื่นจะชนกันเอง
- สถานะ approved / cancelled / rejected ที่ test ต้องจำลอง เป็นของฟีเจอร์ 003 ถึง 006
  test จึงตรวจก่อนว่าค่านั้นอยู่ใน `config.RequestStatuses` จริง ไม่ใช่ค่าที่ test คิดขึ้นเอง

### สิ่งที่เกือบต้องเดาแต่ถามแทน
- ขอบเขตไฟล์ของ handler.go (ดูหัวข้อคำถามด้านบน) ไม่ได้แตะเองก่อนถาม
- **ไม่ได้เดา** ข้อความที่หน้าจอจะแสดง นอกเหนือจากข้อความใน response ของหลังบ้าน
  การแสดงผลเป็นงานของ T-19/T-20 ที่ยังรอ Q-BRQ-14

### ผล test
- `go test ./internal/borrowrequest/ -run "TestAC_BRQ_11|TestDuplicateRequest_ConcurrentSubmitsKeepOneRow" -v`
  PASS ทั้ง 3 ตัว (`TestAC_BRQ_11`, `TestAC_BRQ_11_HTTP`, `TestDuplicateRequest_ConcurrentSubmitsKeepOneRow`)
- `go test ./...`: ทุกแพ็กเกจ ok
- `go vet ./internal/borrowrequest/`: ไม่มีข้อความเตือน
- `gofmt -l` ขึ้นชื่อ `backend/internal/borrowrequest/model.go` ซึ่ง **ไม่ใช่ไฟล์ที่ T-11 แตะ**
  เป็นของเดิมมาก่อน T-11 ยังไม่แก้ เพราะไม่อยู่ในขอบเขต task นี้ (ทีมควรเปิด task เล็กจัดรูปแบบ)

### จุดตรวจเฟส 7 (ทีมลองเองด้วยมือ 2569-11-12)

- ผู้ลอง: ธัญลักษณ์ ด. (@tnyl02) ด้วย `go run ./cmd/server` + curl + psql ใน container
- ใบแรก `CS-NB-0002` ด้วย `mock-student-001`: **201** requestId 3 สถานะ `requested`
- ยิงซ้ำขณะสถานะ `requested`: **409** `{"code":"DUPLICATE_REQUEST"}` และ psql ยืนยันว่ายังมีแถวเดียว
- ยิงซ้ำขณะสถานะ `approved`: **409** `DUPLICATE_REQUEST` (ครึ่งหลังของ AC-BRQ-11)
- แก้ใบเดิมเป็น `cancelled` แล้วยิงใหม่: **201** requestId 8
- ข้อสังเกตระหว่างลอง: ครั้งแรกได้ 503 `UPSTREAM_UNAVAILABLE` เพราะหน้าต่างที่รัน `go run` เปิดไว้ตั้งแต่ก่อนแก้โค้ด
  `go run` คอมไพล์ครั้งเดียวตอนเปิด ต้อง Ctrl+C แล้วเปิดใหม่ทุกครั้งหลังแก้โค้ดก่อนลองด้วยมือ
  (ตอนนั้นฐานข้อมูลกันแถวซ้ำได้ถูกต้องแล้ว เห็นจาก psql ที่มีแถวเดียว ผิดแค่การแปลงเป็นรหัส HTTP ของ build เก่า)
- ล้างข้อมูลที่ลองไว้ด้วย `DELETE FROM borrow_request WHERE requester_id = 'student-001';`

---

## 2569-11-12 /implement T-12 (002-borrow-request) — เฟส 8 ครึ่งแรก

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 8 AC-BRQ-06 ปฏิเสธเมื่อมีของคืนเกินกำหนด)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-12 เพิ่ม handler.go ได้` (ทีมอนุญาตขยายช่องไฟล์ที่แตะไว้ล่วงหน้า เหมือน T-09 และ T-11)
- ทำหลังจาก merge main ลงมาแล้ว (ได้ 005-handover และ 001 รอบล่าสุดติดมาด้วย)

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/clock.go` **ไฟล์ใหม่** (ตรงกับช่องไฟล์ที่แตะ)
  `Clock`, `SystemClock`, `FixedClock`, `overdueDeadline`
- `backend/internal/borrowrequest/service.go` (ตรงกับช่องไฟล์ที่แตะ)
  เปลี่ยนฟิลด์ `now func() time.Time` เป็น `clock Clock`, เพิ่ม `NewServiceWithClock`, `ErrHasOverdueLoan`, `hasOverdueLoan`
- `backend/internal/borrowrequest/repository.go` (ตรงกับช่องไฟล์ที่แตะ)
  เพิ่ม `earliestOpenLoanSQL` และ `EarliestOpenLoanDue` (อ่านตาราง `loan` ของ 001 อย่างเดียว)
- `backend/internal/borrowrequest/handler.go` (เพิ่มในช่องไฟล์ที่แตะวันนี้ ตามที่ทีมอนุญาต)
  เพิ่ม case แปลงเป็น 409 `HAS_OVERDUE_LOAN`
- `backend/internal/borrowrequest/service_test.go`, `handler_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ**
  นอกจาก test ใหม่แล้ว ยังเปลี่ยนวิธีประกอบ service ของ test เดิมทุกตัวให้ใช้นาฬิกาตรึงเวลา (ดูหัวข้อถัดไป)

### การตัดสินใจ
- **test เดิมทุกตัวเปลี่ยนมาใช้ `NewServiceWithClock(..., testClock())` ตรึงเวลาไว้ที่ 09:00 ของวันนี้**
  เหตุผล: ข้อมูลตั้งต้น (T-03) ให้ `student-001` ยืม CS-NB-0004 กำหนดคืน 16:00 **ของวันนี้** และยังไม่คืน
  พอ T-12 บังคับกฎ FR-BRQ-06 แล้ว ถ้า test ยังใช้นาฬิกาจริง test เดิมที่ใช้ student-001 (เช่น TestAC_BRQ_04)
  จะผ่านตอนเช้าแต่ล้มตอนเย็น ไม่ได้ลดความเข้มของการตรวจ เปลี่ยนแค่เวลาที่ใช้ทดสอบให้แน่นอน
- ตรวจกฎคืนเกินกำหนด **ก่อน** ตรวจอุปกรณ์ เพราะเป็นเรื่องของตัวผู้ขอเอง คำตอบ `HAS_OVERDUE_LOAN`
  ไม่เปิดเผยข้อมูลอุปกรณ์ชิ้นใดเลย จึงไม่ขัด FR-BRQ-07 และไม่กระทบ AC-BRQ-07 ที่ผู้ขอภาควิชา 2 ไม่มีของค้าง
- ใช้ `SELECT min(l.due_at)` แทนการดึงทุกใบมาวนเทียบใน Go ใบที่กำหนดคืนเร็วที่สุดเกินกำหนดก็พอปฏิเสธแล้ว
- เส้นแบ่ง 16:00 ใช้ `config.LoanNormalReturnHour` ตัวเดียวกับกฎเวลาคืนใน validate.go ไม่ฝังเลข 16 ซ้ำ
- **ไม่ได้ทำ** กฎปลดล็อก 00:00 ของวันถัดไป (FR-BRQ-06 ข) เพราะเป็น T-13 query จึงกรอง `returned_at IS NULL`

### ข้อผิดพลาดระหว่างทำ และการแก้
`TestOverdue_OtherDepartmentNotCounted` รอบแรกไม่ผ่าน เพราะ AI เข้าใจข้อมูลตั้งต้นผิด
คิดว่า `CS-PJ-0001` เป็นของภาควิชา 2 แต่จริง ๆ เป็นรุ่น 2 ซึ่งอยู่ภาควิชา 1 ส่วนชิ้นของภาควิชา 2 คือ `XX-NB-0001`
แก้ที่ test (ข้อมูลที่ใช้ทดสอบผิด ไม่ใช่โค้ดผิด) และเพิ่มการตรวจภาควิชาของอุปกรณ์ก่อนเริ่ม test กันพลาดซ้ำ

### งานค้างที่ทีมต้องตัดสิน (พบระหว่างทำ T-12)
**T-13 (AC-BRQ-12) มีแค่ `service.go` ในช่องไฟล์ที่แตะ แต่กฎปลดล็อกต้องรู้ `returned_at` ของใบที่คืนแล้ว**
ซึ่ง query ปัจจุบันกรองเฉพาะใบที่ยังไม่คืน ตอนทำ T-13 น่าจะต้องขอเพิ่ม `repository.go` ในช่องไฟล์ที่แตะ
AI ไม่ได้ทำเผื่อไว้ เพราะเป็นกฎของ task อื่น

### ผล test
- `go test ./internal/borrowrequest/ -run "TestAC_BRQ_06|TestOverdue" -v`: PASS ทั้ง 4 ตัว
  (`TestAC_BRQ_06`, `TestAC_BRQ_06_HTTP`, `TestOverdue_BoundaryAtReturnHour`, `TestOverdue_OtherDepartmentNotCounted`)
- `go test ./... -count=1`: ทุกแพ็กเกจ ok รวม `internal/handover` ของ 005 ที่เพิ่ง merge ลงมา
- `go vet ./internal/borrowrequest/`: ไม่มีข้อความเตือน

### จุดตรวจเฟส 8 ครึ่งแรก (ทีมลองเองด้วยมือ 2569-11-12)

- ผู้ลอง: ธัญลักษณ์ ด. (@tnyl02) ด้วย `go run ./cmd/server` + curl + psql ใน container
- `mock-student-002` (มีของเกินกำหนดค้าง): **409** `{"code":"HAS_OVERDUE_LOAN"}` และ psql ยืนยันว่าไม่มีคำขอเกิดขึ้น (0 แถว)
- `mock-student-001` รอบแรก: ได้ 409 เหมือนกัน **ซึ่งถูกต้อง** เพราะข้อมูลในฐานข้อมูลของเครื่องพัฒนาเก่า
  `due_at` ของ CS-NB-0004 เป็น 16:00 ของ 11 พ.ย. (วันที่ seed) ไม่ใช่ของวันนี้ จึงเลยกำหนดไปแล้วจริง
- หลัง `UPDATE loan SET due_at = 16:00 ของวันนี้ WHERE borrower_id = 'student-001'` แล้วยิงซ้ำเวลา 09:10 น.: **201** requestId 9
  ยืนยันเส้นแบ่ง 16:00 บนของจริง ตรงกับ `TestOverdue_BoundaryAtReturnHour` ที่ตรวจ 15:59 กับ 16:01 ด้วยนาฬิกาตรึงเวลา
- **ข้อควรระวังสำหรับคนที่ลองต่อจากนี้** (001 เคยจดไว้แล้ว เจออีกครั้งตอน T-12)
  ข้อมูลตั้งต้นคิดกำหนดคืนจาก `now()` ตอน seed ถ้าฐานข้อมูลของเครื่องพัฒนาถูกสร้างไว้หลายวันแล้ว
  ของทุกคนจะกลายเป็นเกินกำหนดหมด ล้างด้วย `docker compose -f deploy/docker-compose.yml down -v`
  หรือ `UPDATE loan SET due_at = ...` ก่อนลองกฎที่เทียบกับวันนี้ (test อัตโนมัติไม่เจอปัญหานี้
  เพราะ testcontainers สร้างฐานข้อมูลใหม่ทุกครั้ง และ T-12 ตรึงเวลาด้วย FixedClock แล้ว)
- ล้างข้อมูลที่ลองไว้ด้วย `DELETE FROM borrow_request WHERE requester_id IN ('student-001','student-002');`
