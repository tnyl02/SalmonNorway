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

## 2569-11-10 สร้าง spec ฉบับแรกของ 003-cancel-request

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (และตาราง "ใครอ้างใช้บ้าง" ใน 000-shared/constraints.md)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: อ่าน AGENTS.md, specs/000-shared/ ทั้ง 3 ไฟล์ และ SRS หัวข้อ 3.0 กับ 3.1.3
  แล้วสร้าง specs/003-cancel-request/spec.md ครบ 6 ส่วน + Assumptions & Open Questions, Constraint ใช้ร่วมอ้าง ID,
  ใส่ชื่อในช่อง Owner และเพิ่มชื่อฟีเจอร์ในตาราง "ใครอ้างใช้บ้าง"
- สิ่งที่ AI ทำ: spec Draft v1 (FR-CXL-01 ถึง 05, AC-CXL-01 ถึง 07), constraints.md ตาราง "ใครอ้างใช้บ้าง",
  specs/README.md สถานะ, CHANGELOG
- สิ่งที่ AI ตีความหรือร่างเอง ทีมต้องรู้:
  - FR-CXL-03 (ยกเลิกได้เฉพาะคำขอของตนเอง) ตีความจาก RQ-03 เทียบ RQ-10 เขียนเป็น ASM-03 รอยืนยัน
  - FR-CXL-04 (ตัดสินจากสถานะ ณ เวลาที่บันทึก) ร่างจากคำว่า "ก่อน" ใน RQ-03 สำหรับกรณีกดพร้อมกับเจ้าหน้าที่
  - FR-CXL-05 (บันทึกประวัติ) มาจาก RQ-25 ที่เก็บยังไม่รู้ (Q-07)
  - CON-CXL-01, CON-CXL-02 เป็น ID ใหม่เฉพาะฟีเจอร์นี้
  - ไม่ได้อ้าง DOM-AUDIT-01, DOM-PRIV-01, DOM-ACCESS-01 เพราะข้อความพูดถึงข้อมูลอุปกรณ์ ไม่ใช่คำขอ ทีมตรวจใน /clarify
- สิ่งที่ AI ไม่ได้ทำ (ตั้งใจ): ไม่ตั้งชื่อสถานะของคำขอ (รอ 002, Q-01), ไม่เพิ่มหน้ารายการคำขอ (UC-06 ไม่อยู่ใน prototype, Q-04),
  ยังไม่เขียน FR-CXL ใน SRS หัวข้อ 3.1.3

---

## 2569-11-10 ตอบ Open Questions ของ 003 (spec Draft v2)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ตอบ Q-01, Q-02, Q-04 ของ spec Draft v1

### คำถามที่ AI ถาม และคำตอบของทีม
1. Q-01 สถานะไหนนับว่า "ยังไม่อนุมัติ" -> "ถูกยืม, ไม่พร้อมใช้งาน"
   AI ทักกลับ: เป็นสถานะของอุปกรณ์ (CON-DATA-01) ไม่ใช่สถานะของคำขอ และขัดกับ DOM-REQ-01 จึงยังไม่ใส่ใน spec รอยืนยัน
2. Q-02 คำขอค้างคืนที่เจ้าหน้าที่ผ่านขั้นแรกแล้วแต่อาจารย์ยังไม่เซ็น ยกเลิกได้ไหม -> ยกเลิกได้
3. Q-04 เลือกคำขอจากหน้าไหน -> ฟีเจอร์ 003 มีหน้า "คำขอที่รออนุมัติของฉัน" ของตัวเอง (เฉพาะคำขอของตนที่รออนุมัติ,
   แสดงอุปกรณ์ เวลาที่ส่ง ปุ่มยกเลิก, ยืนยันก่อนยกเลิก, ยกเลิกแล้วหายจากหน้า, อนุมัติก่อนยืนยันต้องไม่ยกเลิก, ไม่ใช่ UC-06)

### สิ่งที่ AI ทำ
- spec Draft v2: เพิ่ม FR-CXL-06, FR-CXL-07, AC-CXL-08 ถึง 11, Q-11 แก้ FR-CXL-01, 02, 04 และ Clarifications, CHANGELOG, README
- สิ่งที่ AI ตีความเอง ทีมต้องรู้:
  - "อาจารย์ยังไม่เซ็น" = เจ้าหน้าที่ยังไม่บันทึกผลอนุมัติของอาจารย์ (UC-14) เพราะอาจารย์เซ็นนอกระบบ
  - FR-CXL-07 "ถ้าไม่ยืนยัน สถานะคงเดิม" เขียนเพิ่มจากคำว่า "ต้องยืนยันอีกครั้ง"
  - คำตอบ Q-04 ข้อ "ไม่ยกเลิกคำขอที่ปฏิเสธ/ยกเลิกไปแล้ว" นำไปใส่ใน FR-CXL-02 ด้วย

---

## 2569-11-10 ปิด Q-01 ของ 003 (สถานะของคำขอ)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code
- คำถามที่ AI ถาม: ขออธิบาย Q-01 ใหม่ AI แยกสถานะของอุปกรณ์กับสถานะของคำขอ แล้วเสนอตารางสถานะของคำขอ 8 ค่าจาก workflow บน Miro
  (data-model.md) พร้อมช่องยกเลิกได้/ไม่ได้ (ชื่อภาษาไทยในตาราง AI แปลเอง)
- คำตอบของทีม: ใช้ตามตารางได้เลย
- สิ่งที่ AI ทำ: CON-CXL-02 ใส่ตารางสถานะ, FR-CXL-01, 02, 06 อ้างชื่อสถานะ, ปิด Q-01, เพิ่ม Q-12, CHANGELOG (รวมใน Draft v2 เพราะยังไม่ commit)
- สิ่งที่ AI ตีความเอง: Approved by Officer ใช้กับค้างคืนเท่านั้น คำขอปกติที่เจ้าหน้าที่อนุมัติเป็น Approved ทันที (Q-12)
  ชื่อสถานะยังต้องตรวจกับ spec ของ 002

---

## 2569-11-10 ตอบคำถามจากหน้าสรุปสถานะของคำขอ (003)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ให้ AI ทำหน้าเว็บสรุปตารางสถานะของคำขอ (artifact บน claude.ai) เพื่อส่งให้เจ้าของ 002 และ 004 ตรวจ แล้วตอบคำถามในหน้านั้น

### คำถามที่ AI ถาม และคำตอบของทีม
1. ชื่อสถานะตรงกับ spec ของ 002 ไหม -> ตรงแล้ว
2. Q-12 คำขอปกติผ่าน Approved by Officer ไหม -> ไม่ผ่าน
3. ปฏิเสธได้ที่ขั้นไหนบ้าง -> ทั้ง 2 ขั้น ขั้นเจ้าหน้าที่ และหลังอาจารย์ไม่เซ็น

### สิ่งที่ AI ทำ
- spec: CON-CXL-02 เพิ่มเส้นทางของคำขอ ปิด Q-12 เพิ่ม Clarifications, CHANGELOG (รวมใน Draft v2 เพราะยังไม่ commit), อัปเดตหน้าเว็บสรุป
- หมายเหตุ: ยังไม่มี spec ของ 002 ใน repo (branch นี้) คำว่า "ตรงแล้ว" บันทึกตามที่เจ้าของ 003 ยืนยัน

---

## 2569-11-10 เขียน FR-CXL ลง SRS หัวข้อ 3.1.3

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (docs/srs/SRS.md หัวข้อ 3.1.3)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: เขียน FR-CXL ลง SRS หัวข้อ 3.1.3
- สิ่งที่ AI ทำ: แทน `[เติม]` ด้วย FR-CXL-01 ถึง 07 แบบสรุป (ID และความหมายตาม spec Draft v2 ไม่เพิ่มข้อใหม่) ตามรูปแบบหัวข้อ 3.1.1
  เพิ่ม RQ-25 ในบรรทัดที่มา และบรรทัด "สิ่งที่รู้แล้ว", CHANGELOG
- สิ่งที่ AI ไม่ได้ทำ (ตั้งใจ): ไม่เปลี่ยนเลขเวอร์ชันและวันที่หัวเอกสาร SRS เพราะเจ้าของฟีเจอร์อื่นแก้ SRS ใน branch ของตัวเองพร้อมกัน จะ conflict

---

## 2569-11-10 /clarify รอบที่ 1 ของ 003 (spec Draft v3)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /clarify ไม่ได้โหลดใน session เพราะเปิดจากโฟลเดอร์แม่ AI อ่าน .claude/commands/clarify.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /clarify specs/003-cancel-request/spec.md แล้วให้ทำหน้าเว็บใบตอบคำถาม (artifact บน claude.ai) ตอบผ่านหน้าเว็บแล้ววางคำตอบกลับมา

### คำถามที่ AI ถาม และคำตอบของทีม
1. ประวัติการยกเลิกเก็บที่ไหน บันทึกกรณีไหน -> audit_log (DM-07) เฉพาะผู้ใช้ที่ล็อกอินแล้ว รวมคำขอของผู้อื่น
2. ยกเลิกที่ Approved by Officer ขัด RQ-03 -> บันทึกเป็นการตัดสินของทีม (ASM) และทำใน prototype ด้วย seed
3. ตารางสถานะ / กฎใครบันทึกก่อน เป็นของกลางไหม -> เก็บไว้ใน 003 ให้ 004 เขียนกฎฝั่งตัวเอง
4. หน้า "คำขอที่รออนุมัติของฉัน" แสดงอะไร -> ชื่อรุ่น + รหัสครุภัณฑ์, ใหม่ไปเก่า, "ไม่มีคำขอที่รออนุมัติ"
5. ยกเลิกไม่ได้แจ้งว่าอะไร -> "คำขอนี้ยกเลิกไม่ได้" ทุกกรณี คำขอของผู้อื่น "ไม่พบคำขอ"
6. บันทึกไม่สำเร็จ -> แจ้งข้อผิดพลาด สถานะคงเดิม ลองใหม่ได้ (แบบ FR-EQS-08)
7. session หมดอายุ -> กลับหน้า "คำขอที่รออนุมัติของฉัน" ไม่ยกเลิกอัตโนมัติ
8. Q-05, Q-06, Q-09, Q-10 ตัดออกจาก prototype ได้ไหม -> ไม่ตัด ต้องได้คำตอบก่อน plan

### สิ่งที่ AI แก้
- spec Draft v3: FR-CXL-01, 02, 03, 05, 06, AC-CXL-03, 06, 07, 08, 11 แก้ตามคำตอบ, เพิ่ม AC-CXL-12, 13, ASM-04 ถึง 07, Q-13
  ปิด Q-03, Q-07, Q-08, Q-11, Clarifications, Traceability
- SRS 3.1.3, specs/README.md, CHANGELOG
- สิ่งที่ AI ตีความเอง ทีมต้องรู้:
  - ข้อ 5 "ทุกกรณี" ตีความว่าเฉพาะกรณีที่ยังไม่มีข้อความ กรณี Approved ยังแจ้ง "ได้รับการอนุมัติแล้ว" ตามคำตอบ Q-04 เดิม
  - เพิ่ม AC ใหม่ 2 ข้อ (AC-CXL-12, 13) เพื่อให้คำตอบข้อ 4 และ 6 ทดสอบได้ ไม่ได้เพิ่ม FR ใหม่
  - พบว่า DM-07 ไม่มีฟิลด์รหัสคำขอ เปิดเป็น Q-13 (กระทบของกลางและ 001)

---

## 2569-11-10 ตอบ Q-05, Q-06, Q-09, Q-10 ของ 003 (รวมใน spec Draft v3)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (และตาราง "ใครอ้างใช้บ้าง" ใน constraints.md, SRS 3.1.3 และ 3.3)
- เครื่องมือ: Claude Code

### คำถามที่ AI ถาม และคำตอบของทีม
1. Q-05 ต้องระบุเหตุผลไหม -> ต้องระบุ / AI ถามต่อ: แบบไหน -> พิมพ์เอง บังคับไม่ว่าง
2. Q-06 แจ้งเจ้าหน้าที่ไหม -> แจ้ง / AI ถามต่อ: แจ้งใคร ส่งไม่สำเร็จทำอย่างไร -> Email เจ้าหน้าที่ภาควิชาเจ้าของอุปกรณ์ ยกเลิกสำเร็จเสมอ
3. Q-09 กี่วินาที -> 15 วินาที / AI ถามต่อ: ใช้กับอะไร -> ทั้งยกเลิกและเปิดหน้า p95
4. Q-10 เจ้าหน้าที่ยกเลิกแทนได้ไหม -> ได้ / AI ถามต่อ: ทำแค่ไหน -> เขียนใน spec แต่ไม่ทำใน prototype

### สิ่งที่ AI แก้
- spec: เพิ่ม FR-CXL-08, 09, 10, NFR-PERF-CXL-01, AC-CXL-14 ถึง 17, ASM-08, ASM-09, Q-14 ถึง Q-16 แก้ FR-CXL-07, Scope, Constraints (IF-NOT-01)
- SRS 3.1.3 และ 3.3, constraints.md (IF-NOT-01), CHANGELOG
- สิ่งที่ AI ตีความหรือร่างเอง ทีมต้องรู้:
  - คำตอบทั้ง 4 ข้อ บันทึกเป็นการตัดสินของเจ้าของฟีเจอร์ (ASM-08) ยังไม่ได้ถามนักศึกษาหรือเจ้าหน้าที่จริง
  - เนื้อหาใน Email (ผู้ยกเลิก ชื่อรุ่น รหัสครุภัณฑ์ เวลา เหตุผล) AI ร่างเอง
  - ใช้ ID ใหม่ NFR-PERF-CXL-01 (ไม่ใช้ NFR-PERF-02 กันชนกับฟีเจอร์อื่น) และเพิ่มแถวใน SRS 3.3
  - เหตุผลการยกเลิกต้องมีช่องในตาราง borrow_request ของ 002 (ASM-09)

---

## 2569-11-10 /plan ของ 003 (plan v1, spec Draft v3.1)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (plan mode) คำสั่ง /plan ไม่ได้โหลดใน session เพราะเปิดจากโฟลเดอร์แม่ AI อ่าน .claude/commands/plan.md แล้วทำตามเอง
- สิ่งที่สั่ง: /plan specs/003-cancel-request/spec.md

### คำถามที่ AI ถาม และคำตอบของทีม
1. ตาราง borrow_request ยังไม่มี ใครสร้าง -> 003 ร่างเบื้องต้น ให้ 002 รับช่วง (ASM-11)
2. Q-14 เหตุผลยาวสูงสุด -> 200 ตัวอักษร
3. Q-15 prototype ส่ง Email จริงหรือจำลอง -> จำลอง แบบ auth จำลอง (ASM-10)

### สิ่งที่ AI ทำ
- spec Draft v3.1: FR-CXL-08, AC-CXL-14 (200 ตัวอักษร), ASM-10, ASM-11, ปิด Q-14, Q-15 ปิดบางส่วน, SRS 3.1.3, CHANGELOG, README
- สร้าง specs/003-cancel-request/plan.md Draft v1 (8 หัวข้อ, Constraints 8 ตัว, AC 17 ตัว, ลำดับงาน 9 ขั้น)
- อ้างโครงกลางจาก origin/001-equipment-status (auth.Middleware, Provider pattern, config.go, migration 0001, seed)
- สิ่งที่ AI เลือกเอง ทีมต้องตรวจ:
  - ชื่อ package `cancelrequest`, หน้า `MyPendingRequests`, path API, รหัส error (REQUEST_NOT_FOUND, ALREADY_APPROVED, NOT_CANCELLABLE ฯลฯ)
  - ค่าสถานะในฐานข้อมูลเป็นภาษาอังกฤษตัวเล็ก "ยกเลิก" = `cancelled` ต้องตรงกับ 002
  - ตาราง staff_contact (จำลอง) และ `CancelDBTimeout` (ตัวเลขทีมเลือกเอง)
  - ส่ง Email หลัง commit และไม่เก็บเหตุผลใน audit_log
- สิ่งที่อยากเดาแต่ไม่ได้เดา: รูปแบบช่องรหัสคำขอใน audit_log (Q-13)

---

## 2569-11-10 /tasks ของ 003 (tasks v1)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /tasks ไม่ได้โหลดใน session AI อ่าน .claude/commands/tasks.md แล้วทำตามเอง)
- สิ่งที่สั่ง: `/tasks specs/003-cancel-request/spec.md`
- ผลลัพธ์: specs/003-cancel-request/tasks.md 18 เฟส (ฐานราก + 16 เฟสตาม AC + รอคำตอบ) 32 task รอ Q-13 3 task (T-30 ถึง T-32)
  ตาราง AC 17/17, Constraint 8/8, FR 10 ตัว (FR-CXL-10 ไม่มี AC และไม่มี task เพราะไม่ทำใน prototype)
  แก้ specs/README.md สถานะเป็น tasks v1
- สิ่งที่ AI เลือกเอง ทีมต้องตรวจ:
  - ใช้ไฟล์ seed แยก `seed_003_cancel_request.sql` ไม่แก้ seed.sql ของ 001
  - เพิ่มค่า `MSG_ALREADY_APPROVED` ใน config.js (T-03) ซึ่ง plan หัวข้อ 5 ไม่ได้ระบุ แต่ FR-CXL-02 ต้องใช้ข้อความนี้
  - เรื่องที่รอ 002 (schema) และรอ 001 merge ไม่ใช่ Q ใน spec จึงเขียนไว้ในส่วนหัว ไม่ได้ใส่เป็นสถานะของ task
  - จุดตรวจเฟส 8-11 ให้ลองด้วย curl และ psql เพราะหน้าจอไปถึงกรณีเหล่านี้ไม่ได้ และ 004 ยังไม่มี
  - AC-CXL-17 ต้องมีผู้ใช้จำลอง 50 คน แต่ auth จำลองของกลางมี 4 คน ทีมต้องตัดสินวิธีทดสอบ

---

## 2569-11-11 merge main เข้า branch 003 และปิด Q-13 (spec Draft v3.2)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request และของกลาง (specs/000-shared/data-model.md แก้เพราะ 003)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ดึงงานของเพื่อนใน main เข้า branch 003 โดยไม่ให้งาน 003 เข้า main
- สิ่งที่ AI ทำ: `git merge origin/main` บน branch 003 (เจ้าของฟีเจอร์เริ่ม merge เอง AI แก้ conflict)
  conflict 2 ไฟล์ prompt-log.md และ specs/CHANGELOG.md แก้โดยเก็บทั้งสองฝั่ง ไม่ลบบันทึกใด
  (CHANGELOG: รายการ 003 อยู่บน, prompt-log: บันทึกของ 001 ก่อน แล้วต่อด้วยของ 003) merge commit 65b9588 ยังไม่ push

### คำถามที่ AI ถาม และคำตอบของทีม
1. main มี request_id และ action cancel แล้ว แต่ไม่มี cancel_notify_failed จะบันทึก Email ล้มเหลวอย่างไร
   -> เอาตามที่ AI แนะนำ: เพิ่ม action `notify` ด้วย migration ใหม่ (ใช้ร่วมกับ 004 ได้)
2. แต่ละกรณีใช้ outcome ค่าไหน -> เอาตามที่ AI แนะนำ: `ok` / `forbidden` + สาเหตุใน query, `error` เฉพาะระบบผิดพลาด
   คำขอที่ไม่มีอยู่ให้ request_id ว่าง และบันทึกกรณีเหตุผลไม่ผ่านด้วย

### สิ่งที่ AI แก้
- spec Draft v3.2: FR-CXL-05 (ตารางการบันทึก), FR-CXL-09, ASM-12, ปิด Q-13, Clarifications, Traceability
- data-model.md v4.1 (ของกลาง): เพิ่มแถว action `notify` ใน DM-07 ต้องแจ้งเจ้าของ 001
- plan.md v1.1, tasks.md v1.1 (T-30 ถึง T-32 พร้อมทำ, เฟส 17 มีจุดตรวจ), CHANGELOG
- สิ่งที่ AI พบ ทีมต้องรู้: data-model.md จองเลข migration `0003_borrow_request` ให้ 002 และ `0004` ขึ้นไปให้ 003
  ขัดกับ ASM-11 ที่ให้ 003 ร่าง borrow_request ต้องตกลงกับ 002 ว่าใครสร้างไฟล์ 0003

---

## 2569-11-11 /testcases โหมดร่าง ของ 003 (AC-CXL-08, 01, 09, 14, 12)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /testcases ไม่ได้โหลดใน session AI อ่าน .claude/commands/testcases.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /testcases เส้นทางหลัก 5 AC (เจ้าของฟีเจอร์เลือกจาก 3 ทางที่ AI เสนอ)
- โหมด: ร่าง สร้าง specs/003-cancel-request/test-cases.md ตามแบบในภาคผนวก ยังไม่เขียนโค้ด test
- TC ที่เสนอ 30 แถว สถานะ "ร่าง" ทั้งหมด:
  TC-CXL-08-1 ถึง 08-6, TC-CXL-01-1 ถึง 01-7, TC-CXL-09-1 ถึง 09-6, TC-CXL-14-1 ถึง 14-6, TC-CXL-12-1 ถึง 12-5
- ปรับจากไฟล์คำสั่งที่เขียนสำหรับ Python ให้ตรง AGENTS.md: ตรวจด้วย go test, ชื่อ test `TestTC_CXL_...`
- Then ที่ spec ไม่ได้บอก (เสนอเป็นคำถาม): ลำดับเมื่อเวลาส่งเท่ากัน (08-4), เจ้าหน้าที่เปิดหน้านี้ (08-6),
  ข้อความยกเลิกสำเร็จ (01-2), ช่องเหตุผลหลังกดไม่ยืนยัน (09-5), วิธีนับ "ตัวอักษร" ภาษาไทย (14-5), โหลดรายการไม่สำเร็จ (12-5)

---

## 2569-11-11 ตอบ 6 จุดที่ spec ไม่ได้บอก จาก /testcases (spec Draft v3.3)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code

### คำถามที่ AI ถาม และคำตอบของทีม
1. TC-CXL-08-4 ส่งเวลาเดียวกัน -> เรียงตาม requestId
2. TC-CXL-08-6 เจ้าหน้าที่เรียก -> 403 ถูกปฏิเสธ
3. TC-CXL-01-2 ข้อความสำเร็จ -> "ยกเลิกคำขอยืม {ชื่อรุ่น} ({รหัสครุภัณฑ์}) เรียบร้อยแล้ว"
4. TC-CXL-09-5 กดไม่ยืนยันแล้วเปิดใหม่ -> ล้างช่องเหตุผล
5. TC-CXL-14-5 นับตัวอักษร -> นับตาม code point และแสดงตัวนับ "x/200" ใต้ช่อง
6. TC-CXL-12-5 โหลดไม่สำเร็จ -> "โหลดรายการไม่สำเร็จ กรุณาลองใหม่อีกครั้ง" พร้อมปุ่ม "ลองใหม่"

### สิ่งที่ AI แก้
- spec Draft v3.3: FR-CXL-01, 03, 06, 07, 08 และ ASM-13, Clarifications (แก้ spec ก่อน แล้วค่อยแก้ test-cases ตามกติกาว่า spec เป็นแหล่งข้อมูลหลัก)
- plan v1.2, tasks (T-03, T-04, T-07, T-08, T-10, T-12), test-cases 6 แถว (ยังเป็นสถานะร่าง ทีมต้องเปลี่ยนเป็นใช้ได้เอง), CHANGELOG
- สิ่งที่ AI ตีความเอง ทีมต้องตรวจ:
  - "เรียงตาม requestId" ตีความเป็นมากไปน้อย ให้ตรงกับ "ใหม่ไปเก่า"
  - 403 ใช้กับการสั่งยกเลิกของเจ้าหน้าที่ด้วย ไม่ใช่แค่การเปิดรายการ ตั้งรหัส error `STUDENT_ONLY` เอง
  - TC-CXL-14-5 เพิ่มให้ตรวจทั้ง go test และ vitest เพื่อยืนยันว่านับแบบเดียวกัน

---

## 2569-11-11 /testcases โหมดเขียน test ของ 003 (AC-CXL-08, 01, 09, 14, 12)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (AI อ่าน .claude/commands/testcases.md แล้วทำตามเอง)
- สิ่งที่สั่ง: เขียนโค้ด test จากแถวที่ "ใช้ได้" (ทีมเปลี่ยนสถานะครบ 30 แถว และ commit 3be55cd แล้ว)
- โหมด: เขียน test ไม่ได้แก้โค้ดของระบบ
- ไฟล์ที่สร้าง (ปรับจากไฟล์คำสั่งที่เขียนสำหรับ Python ให้ตรง AGENTS.md):
  - หลังบ้าน `backend/internal/cancelrequest/`: tc_helpers_test.go (ตัวช่วย ไม่ใช่ test), tc_ac_cxl_08_test.go, tc_ac_cxl_01_test.go,
    tc_ac_cxl_09_test.go, tc_ac_cxl_14_test.go, tc_ac_cxl_12_test.go รวม 20 test ชื่อ `TestTC_CXL_...`
  - หน้าจอ `frontend/src/__tests__/`: TC-CXL-08-2, 01-2, 09-1, 09-3, 09-4, 09-5, 14-2, 14-5, 12-2, 12-5 (.test.jsx) รวม 10 ไฟล์
    และ helpers/cxlFakeClient.js (ตัวช่วย ไม่ใช่ test)
  - TC-CXL-09-2 เป็น "คน" ไม่ได้เขียนโค้ด
- ผล test:
  - `go test -short ./...`: auth, equipment, platform ผ่าน / cancelrequest build failed: undefined NewHandler, NewService, NewRepository
  - `npm test`: test ของ 001 ผ่าน 10 ไฟล์ 17 test / TC-CXL 10 ไฟล์ไม่ผ่าน: Failed to resolve import ../pages/MyPendingRequests/MyPendingRequests.jsx
  - ทั้งหมดเป็นกรณี "task ที่ตรวจด้วย AC นี้ยังไม่ได้ทำ" (T-01, T-04 ถึง T-12) ไม่ใช่บั๊ก
- สิ่งที่ test กำหนดไว้ ทีมต้องรู้ตอน implement:
  - หลังบ้าน: NewRepository(pool) / NewService(repo) / NewHandler(svc) / Routes() path `/v1/my/pending-requests` และ `/v1/borrow-requests/{id}/cancel` (ไม่มี /api เหมือน 001)
    package มี TestMain ใน tc_helpers_test.go แล้ว test อื่นของ 003 ห้ามประกาศ TestMain ซ้ำ ให้ใช้ tcServer / tcInsert
  - หน้าจอ: `<MyPendingRequests client>` client มี listMyPendingRequests() และ cancelRequest(id, reason) แถวเป็น listitem
    กล่องยืนยันเป็น dialog มี textbox ตัวนับ "x/200" ปุ่ม "ยืนยัน" / "ไม่ยืนยัน" ปุ่มในแถวชื่อ "ยกเลิก" ปุ่มโหลดใหม่ชื่อ "ลองใหม่"
  - ล้างข้อมูลระหว่าง test ด้วย DELETE FROM borrow_request ถ้า 002 เพิ่ม FK จาก audit_log ต้องแก้ตัวช่วย
- ต้องติดตั้ง dependency หน้าจอก่อนรัน: `cd frontend && npm ci`

---

## 2569-11-11 /implement T-03 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (แก้ของกลาง config.go และ config.js เพราะ 003)
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้โหลดใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-03 เพิ่มค่าที่ตั้งได้ของ 003 ใน config
- ไฟล์ที่แก้หรือสร้าง (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - backend/internal/config/config.go: CancellableStatuses() (CON-CXL-02), CancelReasonMaxLen = 200 (FR-CXL-08), CancelPerfP95 = 15s (NFR-PERF-CXL-01), CancelDBTimeout = 5s (ASM-06)
  - backend/internal/config/config_test.go (ใหม่): TestConfig_CancellableStatuses, TestConfig_CancelLimits
  - frontend/src/config.js: CANCEL_REASON_MAX_LEN, MSG_NO_PENDING, MSG_LOAD_FAILED, LABEL_RETRY, MSG_CANCEL_SUCCESS, MSG_ALREADY_APPROVED, MSG_NOT_CANCELLABLE, MSG_REQUEST_NOT_FOUND
  - ไม่เพิ่ม DisplayTimezone / DISPLAY_TIMEZONE เพราะ 001 มี TimeZone / TIME_ZONE = Asia/Bangkok แล้ว ใช้ตัวเดิม
- ผล test: `go test ./internal/config/ -v` ผ่าน 2/2, `go build ./...` ผ่าน, test ของ auth / equipment / platform ยังผ่าน,
  `npm test` test ของ 001 ยังผ่าน 17/17 (TC-CXL 10 ไฟล์ยังไม่ผ่านเพราะยังไม่มีหน้าจอ ตามคาด)
- สิ่งที่เกือบต้องเดาแต่ถามแทน:
  1. CancelDBTimeout ใช้กี่วินาที -> ทีม: 5 วินาที (ตามตัวอย่างใน plan ของ 001)
  2. ข้อความกรณีอนุมัติแล้ว -> ทีม: "ยกเลิกไม่ได้ เพราะคำขอได้รับการอนุมัติแล้ว" (ตามคำใน FR-CXL-02)
- สถานะใน tasks.md: T-03 เสร็จ

---

## 2569-11-11 บันทึกคำตอบของเจ้าของ 004 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: บันทึกคำตอบของณัฐสุดา (004) 3 ข้อ
  1. ตอนอนุมัติใช้เงื่อนไขสถานะแบบเดียวกัน ไม่เจอแถวแจ้งเจ้าหน้าที่ -> รับไปเขียนใน spec ของ 004
  2. คำขอปกติอนุมัติแล้วเป็น approved ทันที -> ยืนยัน
  3. หน้าอนุมัติไม่แสดงคำขอ cancelled -> ซ่อน
- ไฟล์ที่แก้: spec.md (Draft v3.4, ASM-05 และ Clarifications), plan.md หัวข้อ 8, tasks.md หมายเหตุหัวไฟล์, specs/CHANGELOG.md
- ไม่แก้ FR, AC, test หรือโค้ด เพราะคำตอบตรงกับที่ spec สมมติไว้แล้ว

---

## 2569-11-11 บันทึกการจองเลข migration (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: แจ้งว่าจองเลข migration 0004 และ 0005 ในกลุ่มได้แล้ว ขอสรุปงานที่เหลือ
- ไฟล์ที่แก้: tasks.md หมายเหตุหัวไฟล์, plan.md หัวข้อ 8

---

## 2569-11-11 /implement T-05 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (แก้ของกลาง App.jsx และ api/client.js เพราะ 003)
- เครื่องมือ: Claude Code (AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-05 หน้า "คำขอที่รออนุมัติของฉัน" แสดงรายการ
- ตรวจก่อนเริ่มไม่ผ่านข้อ 2 (ต้องทำหลัง T-04 ซึ่งยังไม่เสร็จ) AI หยุดถาม ทีมตัดสินให้แก้ T-05 เป็นต้องทำหลัง T-03 (บันทึกใน tasks.md และ CHANGELOG)
- ไฟล์ที่แก้หรือสร้าง (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx (ใหม่): รายการ ชื่อรุ่น รหัสครุภัณฑ์ เวลาที่ส่ง ปุ่มยกเลิก (ยังไม่มีการทำงาน รอ T-07)
  - frontend/src/pages/MyPendingRequests/__tests__/MyPendingRequests.test.jsx (ใหม่): AC_CXL_08 shows own pending requests newest first
  - frontend/src/api/client.js: listMyPendingRequests -> GET /v1/my/pending-requests
  - frontend/src/App.jsx: ปุ่มสลับหน้า "รายการอุปกรณ์" / "คำขอที่รออนุมัติของฉัน" (เริ่มที่หน้าของ 001 เหมือนเดิม)
- ผล test: `npx vitest run` ผ่าน 19 test ได้แก่ test ของ 001 17 ตัว, AC_CXL_08 และ TC-CXL-08-2
  ไม่ผ่าน 9 ไฟล์ TC-CXL ตามคาด (ไม่มี dialog รอ T-07/T-08/T-10, ข้อความว่างและโหลดไม่สำเร็จรอ T-12) `npm run build` ผ่าน
- สิ่งที่เกือบต้องเดาแต่ถามแทน:
  1. ลำดับ task (T-04 ยังไม่เสร็จ) -> ทีม: แก้ให้ทำหลัง T-03
  2. รูปแบบเวลาที่ส่ง -> ทีม: ใช้ formatDueAt ของ 001 "2569-11-10 10:30"
  3. ทางเข้าหน้าใน App.jsx -> ทีม: ปุ่มสลับ 2 ปุ่ม ไม่ใช้ router
- สถานะใน tasks.md: T-05 เสร็จ

---

## 2569-11-11 บันทึกคำตอบของเจ้าของ 002 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (อ่านคำตอบจากหน้าเว็บที่ธัญลักษณ์ทำ https://claude.ai/artifact/SrBS61rXxEcNSCGiZkK5uW)
- สิ่งที่สั่ง: บันทึกคำตอบลง spec, plan, tasks ยังไม่ทำ T-01
- คำตอบของ 002: (1) 002 สร้าง borrow_request เองที่ 0003 (2) ชื่อคอลัมน์และสถานะ 8 ค่าใช้ตามที่ 003 เสนอ
  (3) cancel_reason VARCHAR(200) NULL, cancelled_at TIMESTAMPTZ NULL สร้างให้ใน 0003 (4) FK จาก audit_log.request_id เพิ่มใน 0003
  แจ้งเพิ่ม: Q-SH-06 ปิด (ปฏิเสธอัตโนมัติเป็น rejected), 002 เพิ่ม action request, ยืมรายชิ้น
- ไฟล์ที่แก้: spec.md (Draft v3.5: ASM-09, ASM-11, Q-SH-06, Clarifications), plan.md (v1.3 หัวข้อ 2, 3, 7, 8), tasks.md (หัวไฟล์, T-01, T-30, T-31), CHANGELOG, README
- สิ่งที่ AI พบและไม่ได้ตัดสินเอง: index (requester_id, status, submitted_at DESC) ไม่อยู่ในโครงของ 002 ยังไม่รู้ว่าใครเพิ่ม (บันทึกใน plan หัวข้อ 8)
- 002 ขอคำตอบกลับ 2 ข้อ: ชนิดข้อมูลตรงไหม (ตรง) และต้องได้ตารางเมื่อไร (เร็วที่สุด)

---

## 2569-11-11 บันทึกคำตอบรอบที่ 2 ของเจ้าของ 002 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (อ่านจาก https://claude.ai/artifact/SrBS61rXxEcNSCGiZkK5uW ส่วน "รอบที่ 2")
- คำตอบ: (3) data-model v5-v8 อยู่ในเครื่อง 002 ยังไม่ push 003 นำแถว notify ไปวางบน v8 เป็น v9
  (4) 002 ใส่ index idx_request_owner_status ใน 0003 (5) constraint ชื่อ audit_log_action_check 0005 ใส่ค่าครบรวม request
  002 ขอ: ห้าม push 0005 เข้า main ก่อน 0003
- ไฟล์ที่แก้: spec.md (Draft v3.6 ASM-11, ASM-12, Clarifications), plan.md (v1.4 หัวข้อ 3, 8), tasks.md (หัวไฟล์, T-30, T-31), CHANGELOG, README
- สิ่งที่ AI พบและไม่ได้ตัดสินเอง: audit_log เป็น append-only (trigger ใน 0002) เมื่อมี FK จาก 0003 ตัวช่วย test ล้าง borrow_request ไม่ได้
  และแผนเดิมที่ให้ลบแถว audit_log ก่อนทำไม่ได้ บันทึกทางเลือก 3 ทางใน plan หัวข้อ 8 ให้ทีมตัดสินก่อน T-31

---

## 2569-11-11 บันทึกคำตอบของเจ้าของ 001 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code
- คำตอบของพรรณธร: รับ action notify (migration ใหม่ ไม่แก้ 0002, constraint audit_log_action_check), writer.go (T-19 ของ 001) จะรับ pool และ pgx.Tx
  ยังไม่ระบุวันเสร็จ รอไม่ไหวใช้แผนสำรองได้, ขอให้ยืนยันเลข 0003/0004 ในกลุ่ม, รับทราบ App.jsx และ formatDueAt
- ไฟล์ที่แก้: spec.md (Draft v3.7 ASM-12, Clarifications), plan.md (v1.5 หัวข้อ 3, 8), tasks.md (หัวไฟล์), CHANGELOG, README
- สิ่งที่ AI พบ: 001 เขียนว่า CHECK ต้องมี 5 ค่า (ไม่มี request) ขัดกับที่ตกลงกับ 002 ว่า 0005 ต้องมี 6 ค่า ใช้ 6 ค่าและต้องแจ้ง 001
  คำตอบเรื่องวันเสร็จของ T-19 ยังเป็นช่องว่าง "[ใส่วันที่]" ไม่ได้เดาวันที่

---

## 2569-11-11 แก้ลำดับงานหน้าจอ (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: แก้ลำดับงานหน้าจอ ให้ไม่ต้องรองานหลังบ้าน (ใช้เหตุผลเดียวกับ T-05)
- ไฟล์ที่แก้: tasks.md (ช่อง "ต้องทำหลัง" ของ T-07, 10, 12, 16, 18, 20, 22, 24, 26 และหมายเหตุหัวไฟล์), CHANGELOG
- T-08 ไม่ต้องแก้ (ต้องทำหลัง T-07 อยู่แล้ว) ไม่ได้เริ่ม task ไหน

---

## 2569-11-11 /implement T-07 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (แก้ของกลาง frontend/src/api/client.js เพราะ 003)
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-07 กล่องยืนยันการยกเลิกและแสดงผลสำเร็จ
- ไฟล์ที่แก้หรือสร้าง (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - frontend/src/pages/MyPendingRequests/CancelDialog.jsx (ใหม่): กล่อง role=dialog ช่องเหตุผล ปุ่มยืนยัน
  - frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx: ปุ่มยกเลิกเปิดกล่อง ยกเลิกสำเร็จแล้วเอาแถวออกและแสดงข้อความ role=status
  - frontend/src/pages/MyPendingRequests/__tests__/CancelDialog.test.jsx (ใหม่): AC_CXL_01 cancels with reason and shows success, AC_CXL_01 keeps row when cancel fails
  - frontend/src/api/client.js: เพิ่ม postJSON (ใช้ร่วม), apiError (แยกจาก getJSON พฤติกรรมเดิม), cancelRequest
- ผล test: `npx vitest run` ผ่าน 23 test (001 17 ตัว, AC_CXL_08, AC_CXL_01 x2, TC-CXL-08-2, 01-2, 09-3)
  ไม่ผ่าน 7 ตามคาด: TC-CXL-09-1, 09-5 (T-08), 14-2, 14-5 (T-10), 12-2, 12-5, 09-4 (T-12) `npm run build` ผ่าน
- สิ่งที่เกือบต้องเดาแต่ถามแทน:
  1. วิธีส่ง POST ใน client.js -> ทีม: เพิ่ม postJSON ใช้ร่วม
  2. ข้อความสำเร็จแสดงตรงไหน นานแค่ไหน -> ทีม: แถบเหนือรายการ ค้างจนเริ่มยกเลิกรายการถัดไป
- ที่ AI เลือกเองในขอบเขต task: ยกเลิกไม่สำเร็จคงแถวและกล่องไว้โดยยังไม่แสดงข้อความ (ข้อความอยู่ใน T-16, 18, 20, 24) ส่งเหตุผลตามที่พิมพ์ไม่ตัดช่องว่าง (หลังบ้านตรวจใน T-09)
- สถานะใน tasks.md: T-07 เสร็จ

---

## 2569-11-11 /implement T-08 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-08 ปุ่มไม่ยืนยันปิดกล่องโดยไม่ยกเลิก
- ไฟล์ที่แก้:
  - CancelDialog.jsx: ปุ่ม "ไม่ยืนยัน" เรียก onCancel ช่องเหตุผลอยู่ใน state ของกล่อง ปิดแล้วเปิดใหม่จึงว่าง
  - MyPendingRequests.jsx: ส่ง onCancel ให้กล่อง ปิดกล่องโดยไม่เรียก API (ไฟล์นอกช่อง "ไฟล์ที่แตะ" AI หยุดถาม ทีมอนุญาตและเพิ่มในช่องของ T-08 แล้ว)
  - __tests__/CancelDialog.test.jsx: AC_CXL_09 confirm dialog cancel then confirm
- ผล test: `npx vitest run` ผ่าน 25 (เพิ่ม AC_CXL_09 และ TC-CXL-09-1) ไม่ผ่าน 6 ตามคาด:
  TC-CXL-09-5 เหลือแค่ตัวนับ "0/200" (T-10), 14-2, 14-5 (T-10), 12-2, 12-5, 09-4 (T-12) `npm run build` ผ่าน
- สิ่งที่เกือบต้องเดาแต่ถามแทน: แตะ MyPendingRequests.jsx ซึ่งไม่อยู่ในช่องไฟล์ที่แตะ -> ทีม: เพิ่มไฟล์ใน T-08
- สถานะใน tasks.md: T-08 เสร็จ

---

## 2569-11-11 /implement T-10 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-10 ปิดปุ่มยืนยันเมื่อเหตุผลว่างหรือเกิน และแสดงตัวนับ x/200
- ไฟล์ที่แก้ (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - CancelDialog.jsx: countReason (นับ code point), reasonProblem (empty / too_long), ตัวนับใต้ช่อง (aria-describedby), ปุ่มยืนยัน disabled เมื่อเหตุผลใช้ไม่ได้ ใช้ CANCEL_REASON_MAX_LEN จาก config.js
  - __tests__/CancelDialog.test.jsx: AC_CXL_14 confirm disabled when reason empty or too long (รวมขอบ 200 พอดีกดได้)
- ผล test: `npx vitest run` ผ่าน 29 (เพิ่ม AC_CXL_14, TC-CXL-09-5, 14-2, 14-5) ไม่ผ่าน 3 ตามคาด: TC-CXL-12-2, 12-5, 09-4 (T-12) `npm run build` ผ่าน
- สิ่งที่เกือบต้องเดา: ไม่มี (FR-CXL-08, ASM-13 บอกครบ)
- ที่ AI เลือกเองในขอบเขต task: ตัวนับนับข้อความตามที่พิมพ์รวมช่องว่าง ("   " = 3/200) ส่วนการเช็กว่างใช้ trim ตัวนับเป็นสีแดงเมื่อเกิน
- สถานะใน tasks.md: T-10 เสร็จ

---

## 2569-11-11 แก้บั๊ก migrate.go บน Windows (ไฟล์ของกลางของ 001 แก้เพราะ 003)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code
- ปัญหา: ตอนตรวจ task ที่รอทีมตรวจบน Windows `go test ./internal/equipment/ ./internal/platform/` ไม่ผ่าน 21 test และ `go run ./cmd/server` เปิดไม่ขึ้น
  error: migrate init: failed to open source "file://D:\...\db\migrations": invalid port (URL อ่าน path แบบ Windows ไม่ได้)
- แก้: backend/internal/platform/migrate.go บรรทัดเดียว `"file://"+abs` -> `"file://"+filepath.ToSlash(abs)` บน Linux/Codespaces ผลเหมือนเดิม
- ผล test หลังแก้: go build ./... ผ่าน, go test auth / equipment / platform / config ผ่านทั้งหมด
- ต้องแจ้งเจ้าของ 001 (พรรณธร) ไม่อยู่ใน task ของ 003
- หมายเหตุเครื่องพัฒนา: พอร์ต 5432 ถูกใช้โดยโปรเจกต์อื่นและ PostgreSQL ของ Windows จึงรันฐานข้อมูลด้วย docker run ที่พอร์ต 5433 และตั้ง DATABASE_URL ใน terminal ไม่ได้แก้ไฟล์ใน repo

---

## 2569-11-11 ทีมตรวจ T-03, T-05, T-07, T-08, T-10 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจและเช็กสถานะเครื่อง)
- วิธีตรวจ:
  - A รัน test: `npm test` ผ่าน 29 ไม่ผ่าน 3 ตามคาด (TC-CXL-12-2, 12-5, 09-4 เป็นของ T-12), `go test ./internal/config/` ผ่าน
    test ของ 001 (auth, equipment, platform) ผ่านหลังแก้บั๊ก migrate.go บน Windows
  - B รันระบบจริง: ฐานข้อมูล (docker ที่พอร์ต 5433 เพราะ 5432 ถูกใช้ในเครื่อง), หลังบ้าน `go run ./cmd/server`, หน้าเว็บ `npm run dev`
    /healthz 200, /v1/equipment/models มี token ได้ JSON ไม่มี token 401, /v1/my/pending-requests 404 ตามคาด (T-04 ยังไม่ทำ)
    หน้าเว็บมีปุ่มสลับ 2 หน้า หน้า 003 แสดง "กำลังโหลด..." และส่ง GET /api/v1/my/pending-requests จาก client.js
  - C อ่านโค้ดตามเช็กลิสต์
- ข้อจำกัด: กล่องยืนยัน (T-07, T-08, T-10) ตรวจผ่าน test เท่านั้น ยังลองบนจอจริงไม่ได้จนกว่า T-04, T-06 เสร็จ (จุดตรวจเฟส 1, 2)
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-03, T-05, T-07, T-08, T-10 เป็น "เสร็จ"

---

## 2569-11-11 /implement T-12 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-12 ข้อความเมื่อไม่มีคำขอ และโหลดไม่สำเร็จพร้อมปุ่มลองใหม่
- ไฟล์ที่แก้ (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - MyPendingRequests.jsx: ไม่มีคำขอแสดง MSG_NO_PENDING, โหลดไม่สำเร็จแสดง MSG_LOAD_FAILED (role=alert) และปุ่ม LABEL_RETRY โหลดใหม่ ไม่แสดงรายการเก่า
  - __tests__/MyPendingRequests.test.jsx: AC_CXL_12 shows empty message, AC_CXL_12 load failed shows retry and recovers
- ผล test: `npx vitest run` ผ่านทั้งหมด 22 ไฟล์ 34 test รวม TC-CXL หน้าจอครบ 10/10 (TC-CXL-12-2, 12-5, 09-4 ผ่านเพิ่ม) `npm run build` ผ่าน
- สิ่งที่เกือบต้องเดา: ไม่มี (ข้อความอยู่ใน config.js ตาม ASM-13)
- สถานะใน tasks.md: T-12 เสร็จ

---

## 2569-11-11 merge origin/main เข้า branch 003 (ได้ 0003 ของ 002)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ตรวจการ pull จาก main และแก้ conflict
- ได้จาก main: backend/db/migrations/0003_borrow_request (ของ 002 มี cancel_reason, cancelled_at, idx_request_owner_status, FK audit_log.request_id, action request),
  spec/plan/tasks ของ 002, CR-003, CR-004, SRS, migrate_test.go ของ 001
- conflict 5 ไฟล์ แก้แบบรวมทั้งสองฝั่ง:
  - data-model.md: ใช้เนื้อหา v9 ของ 002 แล้วเพิ่มแถว notify ของ 003 เป็น v10 (ตกลงไว้ว่า v9 แต่ 002 ใช้เลข v9 แล้ว)
    แก้หมายเหตุ request_id ว่ามี FK แล้ว และใส่เลข 0004, 0005 ของ 003 ในตารางเลข migration (จองในกลุ่มแล้ว)
  - constraints.md: ตาราง "ใครอ้างใช้บ้าง" รวม 002 และ 003 ในแต่ละแถว
  - specs/README.md: แถว 002 จาก main แถว 003 ของ 003
  - CHANGELOG.md, prompt-log.md: เก็บทั้งสองฝั่ง
- ผล test หลัง merge: go build ผ่าน, go test auth / equipment / platform / config ผ่าน (platform รวม test 0003 ของ 002), npm test ผ่าน 34/34
- สิ่งที่พบ: tc_helpers_test.go (tcInsert) ของ 003 ใช้ purpose 'เพื่อการเรียน' และไม่มี requester_name, student_code, start_at, end_at
  ซึ่ง 0003 จริงบังคับ จดเพิ่มในช่องไฟล์ที่แตะของ T-01 ให้แก้ตอนทำ T-01

---

## 2569-11-11 merge origin/main เข้า branch 003 รอบที่ 2 (001 แก้ migrate.go เอง)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code
- ได้จาก main: PR #8 commit 011960b ของ 001 แก้ platform.Migrate บน Windows ด้วย filepath.ToSlash (บรรทัดเดียวกับที่ 003 แก้ไว้)
- conflict 2 ไฟล์: migrate.go ใช้ของ 001 ทั้งไฟล์ (โค้ดเหมือนกัน ต่างแค่คอมเมนต์ 001 เป็นเจ้าของไฟล์) การแก้ของ 003 จึงไม่เหลือในไฟล์ของกลางแล้ว
  prompt-log.md เก็บทั้งสองฝั่ง

---

## 2569-11-11 /implement T-01 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (เพิ่มไฟล์ใน backend/db/migrations ซึ่งเป็นของกลาง เพราะ 003)
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-01 สร้าง migration staff_contact และตรวจตาราง borrow_request ของ 002
- ไฟล์ที่แก้หรือสร้าง:
  - backend/db/migrations/0004_staff_contact.up.sql (ใหม่): staff_contact(department_id FK department, email) ตาม plan ข้อ 3 (ASM-10)
  - backend/internal/cancelrequest/migrationtest/migration_test.go (ใหม่ package แยก): TestMigration_BorrowRequestStatusCheck,
    TestMigration_CancelReasonMaxLen, TestMigration_StaffContactExists, TestMigration_NoEquipmentSideEffects
  - backend/internal/cancelrequest/tc_helpers_test.go (ตัวช่วย ไม่ใช่ TC): tcInsert ใส่คอลัมน์ตาม 0003 จริง
    (requester_name 'ผู้ทดสอบ', student_code 'TEST-000', purpose 'study', start_at = submitted_at + 1 วัน, end_at = start_at + 6 ชม. ค่าทดสอบ ไม่ได้มาจาก spec)
- ผล test: `go test ./internal/cancelrequest/migrationtest/ -v` ผ่าน 4/4, platform / equipment / config / auth ยังผ่าน
  package cancelrequest ยัง compile ไม่ผ่าน (undefined NewHandler) ตามคาด รอ T-04, T-06
- สิ่งที่เกือบต้องเดาแต่ถามแทน: ที่วาง migration_test.go (package cancelrequest compile ไม่ผ่าน) -> ทีม: package แยก migrationtest
  แก้ช่อง "ไฟล์ที่แตะ" ของ T-01 และ T-30 ให้ตรง
- ไม่ได้ทำ: 0004 ไม่มีไฟล์ .down.sql (ไม่อยู่ในช่องไฟล์ที่แตะ 0001, 0002 ของ 001 ก็ไม่มี)
  ห้าม push 0004 เข้า main ก่อน 0003 ตามที่ตกลง (0003 อยู่ใน main แล้ว จึงไม่ติด)
- สถานะใน tasks.md: T-01 เสร็จ

---

## 2569-11-11 ทีมตรวจ T-01, T-12 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ T-01:
  - `go test -count=1 -v ./internal/cancelrequest/migrationtest/` ผ่าน 4/4 และ test ของ 001/002 (auth, equipment, platform, config) ผ่าน
    (รอบแรก platform ขึ้น "rootless Docker is not supported on Windows" เป็นอาการติดต่อ Docker ไม่ได้ชั่วขณะ รันซ้ำแล้วผ่าน)
  - ฐานข้อมูลในเครื่อง (salmon-db พอร์ต 5433): schema_migrations version 4 dirty f, มีตาราง borrow_request และ staff_contact
    borrow_request มี cancel_reason varchar(200), cancelled_at, idx_request_owner_status, ใส่สถานะ 'ยกเลิก' ถูกปฏิเสธด้วย check constraint
- วิธีตรวจ T-12: หน้าเว็บ http://localhost:5173 หน้า "คำขอที่รออนุมัติของฉัน" แสดง "โหลดรายการไม่สำเร็จ กรุณาลองใหม่อีกครั้ง" และปุ่ม "ลองใหม่"
  (หลังบ้านยังตอบ 404 เพราะ T-04 ยังไม่ทำ) ข้อความ "ไม่มีคำขอที่รออนุมัติ" ตรวจผ่าน test (TC-CXL-12-2) ยังลองบนจอจริงไม่ได้จนกว่า T-04 เสร็จ
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-01, T-12 เป็น "เสร็จ"

---

## 2569-11-11 /implement T-02 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-02 ข้อมูลตั้งต้นคำขอทุกสถานะและ Email เจ้าหน้าที่จำลอง
- ค่าที่ทีมเลือก (ไม่ได้มาจาก spec): student-001 = "ใจดี สวยมาก" 660710755, student-002 = "มีสุข ใจเย็น" 660710756 (AI เสนอ ทีมรับ),
  Email เจ้าหน้าที่ภาควิชา 1 = staff-it@example.test
- ไฟล์ที่แก้หรือสร้าง:
  - backend/db/seed/seed_003_cancel_request.sql (ใหม่): คำขอ id 3001-3008 ตามช่อง "เสร็จเมื่อ" และ staff_contact ภาควิชา 1 โหลดซ้ำได้
  - backend/internal/cancelrequest/migrationtest/seed_test.go (ใหม่ ทีมอนุญาตเพิ่มในช่องไฟล์ที่แตะ): TestSeed003_RequestsAndStaffContact
  - backend/internal/cancelrequest/migrationtest/migration_test.go: เพิ่ม dbURL, ผู้ขอใน test ของ T-01 เปลี่ยนเป็น student-test ไม่ชนกับ seed
- ผล test: `go test -count=1 -v ./internal/cancelrequest/migrationtest/` ผ่าน 5/5
- สิ่งที่ AI เลือกเองในขอบเขต task:
  - ตาราง 0003 มี unique index ห้ามผู้ขอคนเดิมมีคำขอที่ยังไม่จบบนชิ้นเดียวกันเกิน 1 ใบ คำขอที่ยังไม่จบของ student-001 จึงอยู่คนละชิ้น
    (CS-NB-0001, CS-PJ-0001, CS-NB-0002, CS-NB-0003, XX-NB-0001) ใช้เฉพาะชิ้นที่ว่าง คำขอ approved อยู่บน XX-NB-0001 ของภาควิชา 2 เพราะชิ้นว่างของภาควิชา 1 ไม่พอ
  - approved_by_officer และ pending_advisor_approval เป็นยืมค้างคืน (ASM-04), เวลาที่ส่งเป็นวันที่ตายตัว 2569-11-04 ถึง 11-10 วัตถุประสงค์และเวลายืมเป็นค่าทดสอบ
- สถานะใน tasks.md: T-02 เสร็จ รอทีมตรวจ

---

## 2569-11-11 ทีมตรวจ T-02 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ และโหลด seed ใหม่ในฐานข้อมูลของเครื่องพัฒนา)
- วิธีตรวจ:
  - A: `go test -count=1 -v ./internal/cancelrequest/migrationtest/` ผ่าน 5/5 (รวม TestSeed003_RequestsAndStaffContact)
  - B: ฐานข้อมูลในเครื่อง (salmon-db) schema_migrations version 4, โหลด seed ได้ INSERT 0 8 และ INSERT 0 1, โหลดซ้ำได้ INSERT 0 0,
    จำนวนคำขอแต่ละสถานะตรงช่อง "เสร็จเมื่อ", ชื่อภาษาไทยและเหตุผลของ 3007 ถูกต้อง, staff_contact = staff-it@example.test, สถานะอุปกรณ์ไม่เปลี่ยน
  - C: อ่าน seed_003_cancel_request.sql แล้วค่าตรงกับที่ตกลง
- ปัญหาที่พบตอนตรวจ: โหลด seed ด้วย `Get-Content ... | docker exec -i ... psql` ใน PowerShell 5.1 ทำให้ภาษาไทยกลายเป็น ? ในฐานข้อมูล (ไฟล์ seed ไม่ผิด)
  แก้โดยลบแถว 3001-3008 และ staff_contact (ไม่มี audit_log อ้างถึง) แล้วโหลดใหม่ด้วย `docker cp` + `psql -f`
  ถ้าดูผลใน PowerShell ให้ตั้ง [Console]::OutputEncoding เป็น UTF8 ก่อน ไม่อย่างนั้นหน้าจอแสดงภาษาไทยเป็น ?
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-02 เป็น "เสร็จ"

---

## 2569-11-11 /implement T-04 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (แก้ของกลาง backend/cmd/server/main.go เพราะ 003)
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-04 API รายการคำขอที่รออนุมัติของฉัน
- ไฟล์ที่แก้หรือสร้าง:
  - backend/internal/cancelrequest/dto.go, repository.go, service.go, handler.go (ใหม่): GET /v1/my/pending-requests
    เฉพาะคำขอของ Principal.UserID สถานะใน config.CancellableStatuses() เรียง submitted_at DESC, id DESC เจ้าหน้าที่ 403 STUDENT_ONLY
  - backend/internal/cancelrequest/handler_test.go (ใหม่): TestAC_CXL_08_ListOwnPendingSorted, TestListPending_EmptyIsArray
  - backend/cmd/server/main.go: ลงทะเบียน /v1/my/ และ /v1/borrow-requests/ หลัง auth.Middleware
  - specs/003-cancel-request/test-cases.md แถว TC-CXL-08-1, 08-3 (ทีมตัดสิน) และ tc_ac_cxl_08_test.go, tc_helpers_test.go (tcItemAvailable2) ตาม
- ผล test: package cancelrequest compile ผ่านแล้ว AC_CXL_08 ผ่าน, TC-CXL-08-1, 08-3, 08-4, 08-5, 08-6, 12-1, 12-3, 12-4 ผ่าน
  TC ของการยกเลิก (01, 09-6, 14) ยังไม่ผ่านเพราะยังไม่มี route POST (T-06 ขึ้นไป) ตามคาด
  auth / equipment / platform / config / migrationtest ผ่าน, npm test 34/34, ลองหลังบ้านจริงกับ salmon-db (พอร์ต 8010 ชั่วคราว):
  student-001 ได้ 4 รายการ (3001, 3002, 3003, 3004) student-101 ได้ [] staff 403 ไม่มี token 401
- สิ่งที่เกือบต้องเดาแต่ถามแทน:
  1. TC-CXL-08-1, 08-3 ใส่ข้อมูลไม่ได้เพราะชน uq_request_open_per_item ของ 002 -> ทีม: แก้ Given ใน test-cases.md
  2. ลงทะเบียน route ที่ไหน (tasks เขียน platform/http.go แต่ 001 ใช้ main.go) -> ทีม: main.go แบบ 001
- ที่ AI เลือกเองในขอบเขต task: submittedAt ส่งเป็น ISO 8601 เวลา UTC (หน้าจอแปลงเป็น Asia/Bangkok)
- สถานะใน tasks.md: T-04 เสร็จ รอทีมตรวจ

---

## 2569-11-11 ทีมตรวจ T-04 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ:
  - A: `go test -run "AC_CXL_08|ListPending_|TC_CXL_08|TC_CXL_12" ./internal/cancelrequest/` (รันจากโฟลเดอร์ backend) และ test ของ 001/002/migrationtest
  - B: หลังบ้านจริงกับ salmon-db: student-001 ได้ 4 รายการ, student-002 ได้ของตัวเอง, student-101 ได้ [], staff 403, ไม่มี token 401
  - C: หน้าเว็บ "คำขอที่รออนุมัติของฉัน" แสดงคำขอจาก seed (จุดตรวจเฟส 1)
  - D: อ่าน main.go (route ของ 003 อยู่หลัง auth.Middleware), service.go (ผู้ขอมาจาก p.UserID),
    test-cases.md แถว TC-CXL-08-1, 08-3 (เปลี่ยนแค่ชิ้นอุปกรณ์ ช่อง Then เหมือนเดิม) ทีมยืนยันการแก้ Given
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-04 เป็น "เสร็จ"

---

## 2569-11-11 /implement T-06 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-06 API ยกเลิกคำขอ (กรณีสำเร็จ)
- ไฟล์ที่แก้ (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - dto.go: CancelRequestDTO, CancelResultDTO
  - repository.go: Querier เพิ่ม Begin, Cancel ด้วย UPDATE ... WHERE id AND requester_id AND status = ANY(cancellable) RETURNING id (FR-CXL-04)
  - service.go: Cancel เปิด transaction, UPDATE, commit (เผื่อ audit_log ใน transaction เดียวกันของ T-31) เจ้าหน้าที่ ErrStudentOnly
  - handler.go: POST /v1/borrow-requests/{id}/cancel ตอบ 200 {requestId, status:"cancelled"}
  - handler_test.go: TestAC_CXL_01_CancelPendingNormal, TestCancel_SecondTimeDoesNotOverwrite
- ผล test: AC_CXL_01 ผ่าน, TC-CXL-01-1, 01-3, 01-4, 01-5, 09-6, 14-1, 14-3 ผ่าน (รวม TC เดิมของ 08 และ 12)
  ไม่ผ่านตามคาด: TC-CXL-01-6 (T-15), 01-7 (T-19), 14-4, 14-5, 14-6 (T-09) auth / equipment / platform / config / migrationtest ผ่าน
- สิ่งที่ AI เลือกเองชั่วคราวในขอบเขต task (task ถัดไปจะแทนที่):
  - UPDATE ไม่เจอแถว ตอบ 409 NOT_CANCELLABLE ทุกกรณี จนกว่า T-15, T-17, T-19 จะแยกสาเหตุ
  - body อ่านไม่ได้ ตอบ 400 REASON_REQUIRED, id ไม่ใช่ตัวเลข ตอบ 404 REQUEST_NOT_FOUND
  - เก็บเหตุผลตามที่ส่งมา ยังไม่ตรวจว่างหรือยาวเกิน (T-09)
- สถานะใน tasks.md: T-06 เสร็จ รอทีมตรวจ

---

## 2569-11-11 /implement T-09 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-09 ตรวจเหตุผลที่หลังบ้าน
- หมายเหตุ: ก่อนเริ่ม ทีมเปลี่ยนสถานะ T-06 เป็น "เสร็จ" ใน tasks.md เอง (ยังไม่ได้ commit ตอนนั้น)
- ไฟล์ที่แก้ (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - service.go: ValidateReason (trim ว่าง -> ErrReasonRequired, นับ code point ของข้อความเดิมเกิน config.CancelReasonMaxLen -> ErrReasonTooLong) เรียกก่อนเปิด transaction
  - handler.go: 400 REASON_REQUIRED / REASON_TOO_LONG
  - handler_test.go: TestAC_CXL_14_ReasonValidation, TestValidateReason_Boundary
- ผล test: AC_CXL_14 ผ่าน, TC-CXL-14-1, 14-3, 14-4, 14-5, 14-6 ผ่านครบ ไม่ผ่านตามคาด: TC-CXL-01-6 (T-15), 01-7 (T-19)
- สิ่งที่เกือบต้องเดา: ไม่มี (FR-CXL-08, ASM-13)
- ที่ AI เลือกเองในขอบเขต task: นับความยาวจากข้อความเดิมไม่ตัดช่องว่าง ให้ตรงกับตัวนับบนหน้าจอ (T-10) เก็บเหตุผลตามที่ส่งมา
  เจ้าหน้าที่ตรวจ 403 ก่อนตรวจเหตุผล
- สถานะใน tasks.md: T-09 เสร็จ รอทีมตรวจ

---

## 2569-11-11 ทีมตรวจ T-09 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ: `go test -run "AC_CXL_14|ValidateReason|TC_CXL_14" ./internal/cancelrequest/` ผ่าน 7 ตัว
  และเรียก POST /v1/borrow-requests/3002/cancel ด้วยเหตุผล "   " ได้ 400 REASON_REQUIRED คำขอไม่ถูกยกเลิก
  ไม่ได้ทดสอบหน้าเว็บ เพราะ T-09 ไม่แตะ frontend และหน้าเว็บกันเหตุผลผิดไว้แล้วตั้งแต่ T-10
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-09 เป็น "เสร็จ" (T-06 ทีมเปลี่ยนเป็น "เสร็จ" เองก่อนหน้า)

---

## 2569-11-11 /implement T-15 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-15 แยกสาเหตุเมื่อยกเลิกไม่สำเร็จ กรณีอนุมัติแล้ว
- ไฟล์ที่แก้ (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - repository.go: OwnStatus (SELECT status WHERE id AND requester_id ใน transaction เดียวกัน)
  - service.go: UPDATE ไม่เจอแถว -> อ่านสถานะ ถ้า approved -> ErrAlreadyApproved นอกนั้นยังเป็น ErrNotCancelled
  - handler.go: 409 ALREADY_APPROVED
  - handler_test.go: TestAC_CXL_02_ApprovedNotCancellable
- ผล test: AC_CXL_02 ผ่าน, TC-CXL-01-6 ผ่าน, test อื่นในแพ็กเกจยังผ่าน ไม่ผ่านตามคาด: TC-CXL-01-7 (T-19)
- สิ่งที่เกือบต้องเดา: ไม่มี (plan.md ข้อ 4 ขั้น 2)
- สถานะใน tasks.md: T-15 เสร็จ รอทีมตรวจ

---

## 2569-11-11 ทีมตรวจ T-15 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ: `go test -run "AC_CXL_02|TC_CXL_01_6" ./internal/cancelrequest/` ผ่าน
  หลังบ้านจริง: POST cancel ของ 3005 (approved) ได้ 409 ALREADY_APPROVED สถานะยัง approved
  ไม่ได้ทดสอบหน้าเว็บ เพราะ T-15 ไม่แตะ frontend ข้อความบนหน้าเว็บอยู่ใน T-16
- ปัญหาที่พบตอนตรวจ (ไม่ใช่บั๊กของโค้ด):
  - รอบแรกได้ NOT_CANCELLABLE เพราะหลังบ้านที่รันอยู่เปิดก่อนเขียนโค้ด T-15 ต้องปิดแล้วเปิด go run ใหม่ทุกครั้งหลังแก้โค้ด
  - คำสั่ง curl ที่ AI ให้แบบ -d "{\"reason\":...}" ใช้ไม่ได้ใน PowerShell 5.1 (body เสีย ได้ REASON_REQUIRED)
    ข้อความที่มีช่องว่างก็ถูกตัด ใช้ -d '{\"reason\":\"test\"}' หรือเขียน body ลงไฟล์แล้ว --data-binary @body.json
    การตรวจ T-09 ด้วย curl รอบก่อนจึงยังไม่ได้พิสูจน์อะไร AI ตรวจซ้ำด้วย body ที่ถูกต้องแล้ว: เหตุผล "   " ได้ 400 REASON_REQUIRED
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-15 เป็น "เสร็จ"

---

## 2569-11-11 /implement T-19 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-19 ตอบไม่พบคำขอเมื่อไม่ใช่คำขอของตนเอง
- ไฟล์ที่แก้:
  - service.go: UPDATE ไม่เจอแถว และ OwnStatus ไม่พบ (ไม่มี หรือของผู้อื่น) -> ErrRequestNotFound
  - handler.go: 404 REQUEST_NOT_FOUND (ไฟล์นอกช่อง "ไฟล์ที่แตะ" AI หยุดถาม ทีมอนุญาตและเพิ่มในช่องของ T-19 แล้ว)
  - handler_test.go: TestAC_CXL_03_OthersRequestNotFound (คำขอของผู้อื่นและ id ไม่มี ตอบเหมือนกัน body มีแค่ code)
- ผล test: package cancelrequest ผ่านทั้งหมด 28 test รวม TC หลังบ้านทุกตัว (TC-CXL-01-7 ผ่านเพิ่ม), `go test -p 1 ./...` ผ่านทุก package
- สิ่งที่เกือบต้องเดาแต่ถามแทน: แตะ handler.go ซึ่งไม่อยู่ในช่องไฟล์ที่แตะ -> ทีม: เพิ่มใน T-19
- สถานะใน tasks.md: T-19 เสร็จ รอทีมตรวจ

---

## 2569-11-12 ทีมตรวจ T-19 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ: `go test -count=1 -p 1 ./...` ใน backend ผ่านทุก package (auth, cancelrequest, migrationtest, config, equipment, platform)
  รวม TestAC_CXL_03_OthersRequestNotFound และ TC หลังบ้านของ 003 ทุกตัว
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-19 เป็น "เสร็จ"

---

## 2569-11-12 /implement T-16 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-16 แสดงข้อความเมื่อคำขอได้รับการอนุมัติแล้ว
- ไฟล์ที่แก้ (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - MyPendingRequests.jsx: client ตอบ 409 ALREADY_APPROVED -> ปิดกล่อง แสดง MSG_ALREADY_APPROVED (role=alert แถบแดง) แล้วโหลดรายการใหม่
  - __tests__/MyPendingRequests.test.jsx: AC_CXL_02 shows already approved message and reloads
- ผล test: `npx vitest run` ผ่าน 35/35 (รวม TC-CXL หน้าจอ 10/10 และ 001 17 ตัว) `npm run build` ผ่าน
- สิ่งที่เกือบต้องเดา: ไม่มี (plan.md ข้อ 4: 409/404 แสดงข้อความแล้วโหลดรายการใหม่ ข้อความจาก config.js ตาม T-03)
- ที่ AI เลือกเองในขอบเขต task: ปิดกล่องยืนยันเมื่อได้ ALREADY_APPROVED, แถบแดงอยู่ตำแหน่งเดียวกับแถบเขียวของข้อความสำเร็จ
  ค้างไว้จนเริ่มยกเลิกรายการถัดไป (แบบเดียวกับข้อความสำเร็จที่ทีมตัดสินใน T-07)
- สถานะใน tasks.md: T-16 เสร็จ รอทีมตรวจ

---

## 2569-11-12 ทีมตรวจ T-16 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ: `npm test` ผ่าน 35/35 และลองบนจอจริงตาม AC-CXL-04 ด้วยมือ:
  เปิดหน้าเห็นคำขอ 3002 (CS-PJ-0001) -> UPDATE ในฐานข้อมูลเป็น approved แทน 004 -> กดยกเลิกโดยไม่รีเฟรช
  -> กล่องปิด แถบแดง "ยกเลิกไม่ได้ เพราะคำขอได้รับการอนุมัติแล้ว" แถว CS-PJ-0001 หายหลังโหลดรายการใหม่
- รอบแรกไม่เห็นข้อความเพราะหน้าโหลดรายการใหม่หลัง UPDATE แล้ว (ไม่ได้ผ่านขั้นกดยกเลิก) ไม่ใช่บั๊ก
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-16 เป็น "เสร็จ"

---

## 2569-11-12 /implement T-18 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-18 แสดงข้อความ "คำขอนี้ยกเลิกไม่ได้"
- ไฟล์ที่แก้ (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - MyPendingRequests.jsx: เปลี่ยนเงื่อนไขของ T-16 เป็นตาราง CANCEL_REJECTED (รหัส error -> ข้อความใน config.js)
    เพิ่ม NOT_CANCELLABLE -> MSG_NOT_CANCELLABLE ปิดกล่อง แสดงแถบแดง โหลดรายการใหม่ (พฤติกรรมของ ALREADY_APPROVED เหมือนเดิม)
  - __tests__/MyPendingRequests.test.jsx: AC_CXL_11 shows not cancellable message, cancel error clears when opening next dialog
- ผล test: `npx vitest run` ผ่าน 37/37 (TC-CXL หน้าจอ 10/10, AC_CXL_02 ของ T-16 ยังผ่าน) `npm run build` ผ่าน
- สิ่งที่เกือบต้องเดา: ไม่มี (plan.md ข้อ 4)
- สถานะใน tasks.md: T-18 เสร็จ รอทีมตรวจ

---

## 2569-11-12 /implement T-20 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-20 แสดงข้อความ "ไม่พบคำขอ"
- ไฟล์ที่แก้ (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - MyPendingRequests.jsx: เพิ่ม REQUEST_NOT_FOUND -> MSG_REQUEST_NOT_FOUND ในตาราง CANCEL_REJECTED (ปิดกล่อง แถบแดง โหลดรายการใหม่)
  - __tests__/MyPendingRequests.test.jsx: AC_CXL_03 shows request not found message
- ผล test: `npx vitest run` ผ่าน 38/38 `npm run build` ผ่าน
- สิ่งที่เกือบต้องเดา: ไม่มี (plan.md ข้อ 4)
- หมายเหตุ: ต้องทำหลัง T-18 ซึ่งยังเป็น "เสร็จ รอทีมตรวจ" (commit แล้ว)
- สถานะใน tasks.md: T-20 เสร็จ รอทีมตรวจ

---

## 2569-11-12 ทีมตรวจ T-18, T-20 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ: `npm test` ผ่าน 38/38 และลองบนจอจริงโดยเปลี่ยนข้อมูลคำขอ 3002 ระหว่างที่หน้าเปิดค้าง
  - T-18: status = rejected -> กดยกเลิก -> แถบแดง "คำขอนี้ยกเลิกไม่ได้" แถวหายหลังโหลดใหม่
  - T-20: requester_id = student-002 -> กดยกเลิก -> แถบแดง "ไม่พบคำขอ" แถวหายหลังโหลดใหม่
  คืนข้อมูล 3002 เป็น requested ของ student-001 หลังทดสอบ
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-18, T-20 เป็น "เสร็จ"

---

## 2569-11-12 /implement T-11, T-13, T-14, T-17, T-21, T-25 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: ทำกลุ่ม task test ยืนยันหลังบ้าน 6 ตัวในรอบเดียว (ทีมสั่งเป็นกลุ่ม ไม่ใช่ทีละ task ตาม implement.md)
- ไฟล์ที่แก้: backend/internal/cancelrequest/handler_test.go เท่านั้น ไม่ต้องแก้โค้ด เพราะพฤติกรรมถูกต้องแล้วจาก T-04, T-06, T-15, T-19
  - T-11: TestAC_CXL_12_EmptyList (มีแต่คำขอที่ยกเลิกไม่ได้ -> items [] ไม่ใช่ null)
  - T-13: TestAC_CXL_05_OtherRequestsOnSameItemUnchanged (ยกเลิกของ A คำขอของ B บนชิ้นเดียวกันและสถานะอุปกรณ์ไม่เปลี่ยน)
  - T-14: TestAC_CXL_10_CancelOvernightApprovedByOfficer (approved_by_officer และ pending_advisor_approval ยกเลิกได้)
  - T-17: TestAC_CXL_11_RejectedAndCancelledNotCancellable (rejected, cancelled, handed_over, returned -> 409 NOT_CANCELLABLE)
    ช่องไฟล์ที่แตะมี service.go ด้วย แต่ไม่ต้องแก้ เพราะ T-15/T-19 แยกสาเหตุไว้แล้ว
  - T-21: TestAC_CXL_04_ApprovedBeforeConfirm (GET เห็นคำขอ -> UPDATE เป็น approved แทน 004 -> POST -> 409 ALREADY_APPROVED)
  - T-25: TestAC_CXL_07_ExpiredSessionNoCancel (ไม่มี token และ token ที่ไม่รู้จัก -> 401 สถานะไม่เปลี่ยน ส่วน audit_log อยู่ใน T-31)
- ผล test: 6 test ใหม่ผ่านทั้งหมด package cancelrequest ผ่านทั้งหมด
- สิ่งที่เกือบต้องเดา: ไม่มี (AC ใน spec ครบ)
- สถานะใน tasks.md: T-11, T-13, T-14, T-17, T-21, T-25 เสร็จ รอทีมตรวจ

---

## 2569-11-12 ทีมตรวจ T-11, T-13, T-14, T-17, T-21, T-25 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ: รัน test ใน backend ผ่านทั้ง 6 ตัว (TestAC_CXL_12_EmptyList, 05, 10, 11, 04, 07) ไม่ได้ลองบนจอจริงเพราะเป็น test อย่างเดียว
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-11, T-13, T-14, T-17, T-21, T-25 เป็น "เสร็จ"

---

## 2569-11-12 /implement T-22 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-22 ทดสอบหน้าจอเมื่อกดยืนยันแล้วได้ผลว่าอนุมัติแล้ว
- ไฟล์ที่แก้ (ตรงกับช่อง "ไฟล์ที่แตะ"): frontend/src/pages/MyPendingRequests/__tests__/CancelDialog.test.jsx
  เพิ่ม AC_CXL_04 approved while dialog open (เปิดกล่องก่อน แล้ว client จำลองเปลี่ยนเป็นตอบ 409 ALREADY_APPROVED) ไม่ต้องแก้โค้ด เพราะ T-16 ทำไว้แล้ว
- ผล test: `npx vitest run` ผ่าน 39/39
- สิ่งที่เกือบต้องเดา: ไม่มี
- สถานะใน tasks.md: T-22 เสร็จ รอทีมตรวจ

---

## 2569-11-12 /implement T-23, T-24 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (แก้ของกลาง frontend/src/config.js เพราะ 003)
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: ทำ T-23 และ T-24 ในรอบเดียว
- T-23 (หลังบ้าน):
  - service.go: Service มี dbTimeout = config.CancelDBTimeout, Cancel ใช้ context.WithTimeout ฐานข้อมูลไม่ตอบ -> error -> 503 และ rollback
  - handler.go: คอมเมนต์ writeUpstreamError (503 UPSTREAM_UNAVAILABLE retryable มีอยู่แล้วตั้งแต่ T-04)
  - handler_test.go: TestAC_CXL_13_DBDownReturns503
  - ต่างจากช่อง "เสร็จเมื่อ": ไม่หยุด container เพราะ test ทั้งแพ็กเกจใช้ฐานข้อมูลร่วมกัน (หยุดแล้วเปิดใหม่พอร์ตเปลี่ยน)
    จำลองด้วย pool ที่ต่อไม่ได้ (127.0.0.1:1) และฐานข้อมูลที่ค้างเกินเวลา (dbTimeout 50ms) แล้วตรวจกับฐานข้อมูลจริงว่าสถานะยัง requested และยกเลิกใหม่ได้
- T-24 (หน้าเว็บ):
  - config.js: MSG_CANCEL_FAILED = "ยกเลิกไม่สำเร็จ กรุณาลองใหม่อีกครั้ง" (ไฟล์นอกช่องไฟล์ที่แตะ เพิ่มในช่องของ T-24 แล้ว)
  - CancelDialog.jsx: prop error แสดงข้อความในกล่อง (role=alert) ปุ่มยืนยันเปลี่ยนชื่อเป็น "ลองใหม่" เหตุผลที่พิมพ์ยังอยู่
  - MyPendingRequests.jsx: error ที่ไม่อยู่ใน CANCEL_REJECTED (503, ส่งไม่ถึง) -> dialogError ไม่ปิดกล่อง ไม่โหลดรายการใหม่ แถวยังอยู่
  - CancelDialog.test.jsx: AC_CXL_13 shows error with retry, dialog error clears when reopened
- สิ่งที่เกือบต้องเดาแต่ถามแทน:
  1. ข้อความเมื่อ 503 -> ทีม: "ยกเลิกไม่สำเร็จ กรุณาลองใหม่อีกครั้ง"
  2. ปุ่มลองใหม่อยู่ที่ไหน -> ทีม: ในกล่องยืนยัน ส่งเหตุผลเดิม
- ผล test: go test ./internal/cancelrequest/ ผ่านทั้งหมด, npx vitest run ผ่าน 41/41, npm run build ผ่าน
- สถานะใน tasks.md: T-23, T-24 เสร็จ รอทีมตรวจ

---

## 2569-11-12 ทีมตรวจ T-23, T-24 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ: go test ./... ผ่านทุก package, npm test ผ่าน 41/41 และลองบนจอจริงตามวิธีที่ AI ให้ (หยุด salmon-db ระหว่างกล่องยืนยันเปิด -> ข้อความแดงในกล่องและปุ่มลองใหม่ -> เปิดฐานข้อมูลกลับแล้วลองใหม่สำเร็จ)
- ทีมยอมรับวิธีทดสอบ T-23 ที่จำลองฐานข้อมูลล่มแทนการหยุด container ของ test
- ข้อจำกัดที่รู้: session หมดอายุ (401) ตอนนี้ขึ้นข้อความ "ยกเลิกไม่สำเร็จ" ด้วย จนกว่า T-26 จะพาไปหน้าเข้าสู่ระบบ
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-23, T-24 เป็น "เสร็จ"

---

## 2569-11-12 ทีมตรวจ T-22 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ทีมเปลี่ยนสถานะ T-22 ใน tasks.md เป็น "เสร็จ" เอง (บันทึกย้อนหลังโดย Claude Code ตามที่ทีมแจ้ง)
- วิธีตรวจ: npm test ผ่าน (รวม AC_CXL_04 approved while dialog open) ลำดับเหตุการณ์ AC-CXL-04 บนจอจริงทดสอบไปแล้วตอนตรวจ T-16

---

## 2569-11-12 /implement T-30 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request (เพิ่มไฟล์ใน backend/db/migrations ซึ่งเป็นของกลาง เพราะ 003)
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-30 เพิ่ม action notify ใน audit_log
- ก่อนเริ่มเช็ก origin/main: ยังไม่มี migration เลข 0005 และ CHECK ใน 0003 ของ 002 เป็น 5 ค่า (list, search, detail, cancel, request)
- ไฟล์ที่แก้หรือสร้าง (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - backend/db/migrations/0005_audit_log_notify_action.up.sql (ใหม่): DROP แล้ว ADD audit_log_action_check ชื่อเดิม 6 ค่า รวม notify ไม่แก้ 0002
  - backend/internal/cancelrequest/migrationtest/migration_test.go: TestMigration_AuditLogAcceptsNotify (ใส่ใน transaction ที่ rollback ทิ้ง)
- ผล test: `go test -p 1 ./...` ผ่านทุก package (รวม test audit_log ของ 001 และ 002 ใน platform)
- สิ่งที่เกือบต้องเดา: ไม่มี (ตกลงชื่อ constraint และค่าครบกับ 001, 002 แล้ว)
- ไม่ได้ทำ: ไม่มีไฟล์ .down.sql (ไม่อยู่ในช่องไฟล์ที่แตะ แบบเดียวกับ 0004)
  data-model.md v10 มีแถว notify อยู่แล้วตั้งแต่ merge main ไม่ต้องแก้
- สถานะใน tasks.md: T-30 เสร็จ รอทีมตรวจ

---

## 2569-11-12 ทีมตัดสินวิธีล้างข้อมูลใน test และ audit writer (003-cancel-request)

- ผู้ตัดสิน: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (อธิบายทางเลือกและผลกระทบ)
- เรื่องที่ 1 ล้างข้อมูลระหว่าง test เมื่อ audit_log (append-only) มี FK ไป borrow_request: เลือก (ข)
  ปิด trigger และ FK ชั่วคราวด้วย SET LOCAL session_replication_role = replica ใน transaction ของการล้างเท่านั้น
  ข้อระวังที่ AI แจ้ง: ห้ามใช้ SET ธรรมดาเพราะค้างใน connection ของ pool, ใช้ได้เพราะผู้ใช้ test เป็น superuser,
  ต้องเขียนคอมเมนต์ว่าใช้กับฐานข้อมูล test เท่านั้น กฎ append-only ยังตรวจด้วย test ของ 001
- เรื่องที่ 2 audit writer: รอ 001 ทำ writer.go (T-19 ของ 001) ก่อน ยังไม่ใช้แผนสำรอง T-31, T-32 จึงยังไม่เริ่ม
- ไฟล์ที่แก้: specs/003-cancel-request/plan.md หัวข้อ 8 (บันทึกการตัดสิน) ยังไม่ได้แก้โค้ด

---

## 2569-11-12 ทีมตรวจ T-30 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ: go test (TestMigration_AuditLogAcceptsNotify และทุก package) ผ่าน
  ฐานข้อมูลในเครื่อง: schema_migrations version 5, audit_log_action_check มี 6 ค่า, ใส่ notify ได้ (ใน transaction ที่ rollback),
  ใส่ cancel_notify_failed ถูกปฏิเสธด้วย check constraint ไม่มีข้อมูลทดสอบค้างใน audit_log
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-30 เป็น "เสร็จ"

---

## 2569-11-12 merge origin/main เข้า branch 003 รอบที่ 3 (ได้งาน 002 PR #9 และ 005 PR #10)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code
- ได้จาก main: backend/internal/borrowrequest (002), backend/internal/handover และ 0006_handover, handover_dev.sql, หน้า handover (005),
  spec/plan/tasks ของ 005, seed.sql (แค่คอมเมนต์ของ 002), migrate_test.go, SRS
- ตรวจก่อนแก้: 0006_handover ของ 005 ไม่แตะ audit_log จึงไม่ลบ action notify ของ 0005
  handover_dev.sql ใช้อุปกรณ์ชุดของตัวเอง (HO-TEST-*) และ id อัตโนมัติ ไม่ชนกับ seed_003 (3001-3008)
  ถ้าโหลดทั้งสองไฟล์ student-001 จะเห็น HO-TEST-0003 (requested) ในหน้า 003 เพิ่มอีก 1 รายการ
- conflict 7 ไฟล์ แก้แบบรวมทั้งสองฝั่ง:
  - backend/internal/config/config.go, frontend/src/config.js: ค่าของ 003 และ 002 อยู่คู่กัน
  - backend/internal/config/config_test.go: รวม import (slices, testing, time) และ test ของ 003 กับ 002
  - backend/cmd/server/main.go: import ทั้ง borrowrequest และ cancelrequest (เรียงตามตัวอักษร) route ของทุกฟีเจอร์อยู่ครบ
  - specs/000-shared/constraints.md: ตาราง "ใครอ้างใช้บ้าง" รวม 002, 003, 005 ในแต่ละแถว
  - specs/CHANGELOG.md, prompt-log.md: เก็บทั้งสองฝั่ง
- ผล test หลัง merge: go build, go vet ผ่าน, go test -p 1 ./... ผ่านทุก package (รวม borrowrequest, handover), npm test 47/47, npm run build ผ่าน

---

## 2569-11-12 /implement T-27 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-27 สร้าง Notifier จำลองและส่งหลังยกเลิกสำเร็จ
- ไฟล์ที่แก้หรือสร้าง:
  - notifier.go (ใหม่): CancelNotice, interface Notifier, LogNotifier (Email จำลองเขียนลง log ASM-10)
  - repository.go: CancelledNotice อ่านชื่อรุ่น รหัสครุภัณฑ์ ภาควิชาเจ้าของอุปกรณ์ cancelled_at และ Email จาก staff_contact ใน transaction
  - service.go: NewService ใช้ LogNotifier เป็นค่าเริ่มต้น, WithNotifier เปลี่ยนตัวส่ง, ส่งหลัง commit เท่านั้น ส่งไม่สำเร็จเขียน log ไม่ย้อนการยกเลิก
  - notifier_test.go (ใหม่): TestAC_CXL_15_NotifyDepartmentStaff, TestNotify_NotSentWhenNotCancelled, TestNotify_NoStaffEmailStillCancelled
  - ไม่แก้ backend/cmd/server/main.go (อยู่ในช่องไฟล์ที่แตะ) เพราะ NewService ใช้ LogNotifier เป็นค่าเริ่มต้นแล้ว ไม่ต้องแตะไฟล์กลาง
    ไม่แก้ tc_helpers_test.go เพราะวิธีสร้าง Service ไม่เปลี่ยน
- ผล test: package cancelrequest ผ่านทั้งหมด `go test -p 1 ./...` ผ่านทุก package
- ที่ AI เลือกเองในขอบเขต task (spec ไม่ได้บอก ต้องให้ทีมตรวจ):
  - ภาควิชาที่ไม่มี Email ใน staff_contact: ยกเลิกสำเร็จตามปกติ ไม่เรียก Notifier เขียน log ไว้
  - ผู้ยกเลิกในข้อความเป็นรหัสผู้ใช้ (Principal.UserID) เพราะ IF-IDP-01 ไม่ส่งชื่อมา
  - ส่ง Email แบบรอผล (ไม่แยก goroutine) หลัง commit ใช้ context.WithoutCancel ไม่ให้ dbTimeout ตัดการส่ง
- สถานะใน tasks.md: T-27 เสร็จ รอทีมตรวจ

---

## 2569-11-12 แก้รูปแบบเวลาใน Email จำลอง (T-27, 003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ตอนลองบนจอจริง เวลาในฐานข้อมูล (2026-11-11 18:25:39+00 UTC) กับใน Email จำลอง (2026-11-12T01:25:39+07:00) ดูไม่เหมือนกัน
  ขอให้แสดงเหมือนกันและไม่ต้องมีวินาที (เป็นเวลาเดียวกัน ต่างแค่เขตเวลา)
- ไฟล์ที่แก้: notifier.go เพิ่ม FormatNoticeTime เวลาไทย (config.TimeZone) ปี พ.ศ. ไม่มีวินาที "2569-11-12 01:25" แบบเดียวกับหน้าเว็บ
  notifier_test.go เพิ่ม TestFormatNoticeTime ค่าที่เก็บในฐานข้อมูลไม่เปลี่ยน
- ที่ AI เลือก: ใช้ปี พ.ศ. ให้ตรงกับเวลาบนหน้าเว็บ (formatDueAt ของ 001), ตัดวินาทีทิ้งไม่ปัด
- คำสั่งดูเวลาในฐานข้อมูลแบบเดียวกัน: to_char ใน psql (ให้ไว้ในแชต)

---

## 2569-11-12 ลอง T-27 บนจอจริงผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- วิธีตรวจ: เปิดหลังบ้านใหม่ ยกเลิกคำขอบนหน้าเว็บ Email จำลองใน log มีผู้รับ staff-it@example.test ผู้ยกเลิก ชื่อรุ่น รหัสครุภัณฑ์
  เวลา (รูปแบบใหม่ 2569-11-12 01:25 ตรงกับ cancelled_at ในฐานข้อมูลเมื่อแปลงเป็นเวลาไทย) และเหตุผล
- ยังค้าง: ทีมยังไม่ตอบ 3 เรื่องที่ spec ไม่ได้บอก (ภาควิชาที่ไม่มี Email, ผู้ยกเลิกเป็นรหัสหรือชื่อ, รอผลหรือส่งเบื้องหลัง) สถานะ T-27 ยังเป็น "เสร็จ รอทีมตรวจ"

---

## 2569-11-12 ทีมตัดสิน 3 เรื่องของ Email (T-27, 003-cancel-request)

- ผู้ตัดสิน: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code
- คำตอบ: (1) ภาควิชาที่ไม่มี Email เจ้าหน้าที่ นับเป็นส่งไม่สำเร็จ (2) ผู้ยกเลิกแสดงทั้งชื่อ/รหัสนักศึกษาจากคำขอและรหัสผู้ใช้ (3) ส่งแบบรอผล
- ไฟล์ที่แก้:
  - notifier.go: CancelNotice เพิ่ม RequesterName, StudentCode, FormatCanceller, ErrNoRecipients
  - repository.go: CancelledNotice อ่าน requester_name, student_code
  - service.go: notify คืน error (ไม่มีผู้รับ = ErrNoRecipients) ให้ T-32 ใช้บันทึก audit การยกเลิกยังสำเร็จ
  - notifier_test.go: ตรวจ FormatCanceller, เพิ่ม TestNotify_ReportsFailure
  - spec.md ASM-10 (บันทึก 4 ข้อของ Email), Q-15 (เพิ่มคำถามเวลาส่งของตัวส่งจริง), tasks.md T-32, CHANGELOG
- ผล test: `go test -p 1 ./...` ผ่านทุก package

---

## 2569-11-12 ทีมตรวจ T-27 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ:
  - A: test 5 ตัว (AC_CXL_15, Notify_NotSent, Notify_NoStaffEmail, Notify_ReportsFailure, FormatNoticeTime) และ go test -p 1 ./... ผ่าน
  - B: บนจอจริง กรณียกเลิกสำเร็จ Email จำลองมีผู้รับ ผู้ยกเลิก (ชื่อ รหัสนักศึกษา รหัสผู้ใช้) ชื่อรุ่น รหัสครุภัณฑ์ เวลา (ตรงกับฐานข้อมูล) เหตุผล
    กรณียกเลิกไม่สำเร็จไม่มี Email, กรณีภาควิชา 2 ไม่มี Email ขึ้น "ส่งไม่สำเร็จ: no staff email" และยกเลิกสำเร็จ
  - C: อ่านโค้ด ส่งหลัง commit และหาภาควิชาจากอุปกรณ์
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-27 เป็น "เสร็จ"

---

## 2569-11-12 /implement T-28 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-28 ไม่ย้อนการยกเลิกเมื่อ Notifier ล้มเหลว
- ไฟล์ที่แก้: backend/internal/cancelrequest/notifier_test.go เท่านั้น
  - TestAC_CXL_16_NotifyFailureStillCancelled (Notifier ปลอมคืน error -> 200 cancelled เหตุผลถูกเก็บ คำขอหายจากรายการ)
  - tcHandlerOnly ตัวช่วยตรวจรายการโดยไม่ล้างข้อมูล
  - ไม่ต้องแก้ service.go (อยู่ในช่องไฟล์ที่แตะ) เพราะ T-27 ส่ง Email หลัง commit และไม่ย้อนการยกเลิกไว้แล้ว
- ผล test: package cancelrequest ผ่านทั้งหมด (รอบแรกใช้ 31 วินาทีเพราะ Docker เริ่มฐานข้อมูล test ช้า รันซ้ำ 3.4 วินาที)
- ข้อค้นพบที่ต้องให้ทีมตัดสิน (ไม่ได้แก้): Notifier ที่ "ค้างไม่ตอบ" (ไม่ใช่คืน error) จะทำให้นักศึกษารอไม่สิ้นสุด
  เพราะส่งแบบรอผลและไม่มีเวลาจำกัด ตัวจำลองไม่เกิด แต่ตัวส่งจริงอาจเกิด เสนอเพิ่มเวลารอการส่ง Email ใน config (ต้องให้ทีมกำหนดตัวเลข)
- สถานะใน tasks.md: T-28 เสร็จ รอทีมตรวจ

---

## 2569-11-12 ทีมตรวจ T-28 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- วิธีตรวจ: `go test -run AC_CXL_16 ./internal/cancelrequest/` ผ่าน (ลองบนจอจริงไม่ได้ เพราะตัวจำลองส่งสำเร็จเสมอ)
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-28 เป็น "เสร็จ"
- ยังค้าง: ทีมยังไม่ตัดสินเรื่องเวลารอการส่ง Email (NotifyTimeout) กรณีตัวส่งค้างไม่ตอบ

---

## 2569-11-12 /implement T-29 (003-cancel-request)

- ผู้สั่ง: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- ฟีเจอร์: 003-cancel-request
- เครื่องมือ: Claude Code (คำสั่ง /implement ไม่ได้ติดตั้งใน session AI อ่าน .claude/commands/implement.md แล้วทำตามเอง)
- สิ่งที่สั่ง: /implement T-29 เขียน load test ด้วย k6 (ทีมเลือกรัน k6 ผ่าน Docker เพราะเครื่องไม่มี k6)
- ไฟล์ที่สร้าง (ตรงกับช่อง "ไฟล์ที่แตะ"):
  - backend/db/seed/seed_003_load.sql: อุปกรณ์ทดสอบ LD-0001 ถึง LD-1000 (รุ่น 5001 ภาควิชา 1) และคำขอ 50001-51000
    (ครึ่งแรก student-001 ครึ่งหลัง student-002) โหลดซ้ำ = รีเซ็ตเป็น requested ด้วย UPDATE (ไม่ใช้ DELETE เผื่อมี audit_log ใน T-31)
  - tests/load/cancel_request.js: 50 VU per-vu-iterations คนละ 20 รอบ แต่ละรอบ GET รายการ + POST ยกเลิกคำขอของตัวเอง (ไม่แย่งใบเดียวกัน)
    threshold p95 < 15000 ms แยกของ list และ cancel (ค่าจาก config.CancelPerfP95) และ checks rate == 1.0
- ที่ AI เลือกเอง (ไม่ได้มาจาก spec): จำนวนรอบ 20 ต่อผู้ใช้ (รวม 1,000), รันกับฐานข้อมูลแยก salmon_load ใน container เดิม
  เพื่อไม่ปนกับข้อมูลที่ใช้ลองหน้าเว็บ, เรียกหลังบ้านตรงไม่ผ่าน proxy ของ vite
- ผล: หลังบ้านชั่วคราวพอร์ต 8010 ต่อ salmon_load, k6 v2.3.0 ใน Docker
  list p95 = 106.63 ms, cancel p95 = 38.9 ms, checks 2000/2000, http_req_failed 0%, คำขอ 1,000 ใบเป็น cancelled
  Email จำลอง 1,000 บรรทัด ฐานข้อมูล salmon (ที่ใช้ลองหน้าเว็บ) ไม่มีข้อมูล load test
- ข้อจำกัด: วัดบนเครื่องพัฒนาเครื่องเดียว (หลังบ้าน ฐานข้อมูล และ k6 อยู่เครื่องเดียวกัน) ไม่ใช่สภาพแวดล้อมจริง และยังไม่มี audit_log (T-31)
- สถานะใน tasks.md: T-29 เสร็จ รอทีมตรวจ

---

## 2569-11-12 ทีมตรวจ T-29 ผ่าน (003-cancel-request)

- ผู้ตรวจ: 660710763 นางสาวนภัสนัน ดีเสมอ (@napatsanan-np)
- เครื่องมือ: Claude Code (ช่วยอธิบายวิธีตรวจ)
- วิธีตรวจ: หลังบ้านต่อฐานข้อมูลแยก salmon_load, โหลด seed_003_load.sql (รีเซ็ตคำขอ 1,000 ใบ),
  รัน k6 ผ่าน Docker (grafana/k6) กับ tests/load/cancel_request.js ผ่าน threshold ทั้ง 3 (list p95, cancel p95 < 15000 ms, checks 100%)
- ผล: เปลี่ยนสถานะใน tasks.md ของ T-29 เป็น "เสร็จ"


---

## 2569-11-12 CR-005 อนุมัติ: spec 001 v5, SRS, constraints, plan v1.7, tasks v2.2

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ของกลาง `specs/000-shared/constraints.md` ข้อ Q-SH-05 และ `docs/srs/SRS.md` 3.1.1 เพราะ CR-005 ของ 001)
- เครื่องมือ: Claude Code
- ที่มา: ทีมรับ CR-005 ทั้งหมด 2569-11-12 ผู้สั่งบันทึกผลใน CR-005 และแก้ spec เป็น v5 เองบางส่วน แล้วสั่ง AI ทำส่วนที่เหลือทั้งหมด
- การตัดสินของทีมตอนอนุมัติ (ผู้สั่งตอบตามที่ AI เสนอ):
  - ก. หน้าจอเรียก endpoint ลำดับคำขอของ 002 ทุกครั้ง ถ้าได้ 403 ไม่แสดงส่วนนั้น ไม่สร้าง endpoint บอกบทบาทเพิ่ม
    แก้ประโยคใน CR-005 จาก "หน้าจอไม่เรียกเมื่อเป็นนักศึกษา" เป็น "นักศึกษาไม่เห็นส่วนนี้"
  - ข. นักศึกษาได้ 403 รหัส `FORBIDDEN_ROLE` (มีในรายการกลาง CON-TECH-03 แล้ว)
- สิ่งที่ AI ตรวจและแก้ใน spec v5 ที่ผู้สั่งแก้มา:
  - บรรทัด "แสดงลำดับคำขอที่ค้างอยู่ให้เจ้าหน้าที่" อยู่ใต้ Out of scope ย้ายขึ้นไป In scope
  - แถว Traceability มี 6 คอลัมน์ แก้เป็น 3 คอลัมน์ตามตาราง
  - หมายเหตุ Updated ที่หัวไฟล์ระบุ CR-005 / Clarifications เพิ่ม Q/A ของข้อ ก และ ข
- สิ่งที่ AI แก้:
  - CR-005: ตาราง API (FORBIDDEN_ROLE, เรียกทุกครั้ง), ผลกระทบ (plan, DOM-PRIV-01), หัวข้อตัดสิน (ผู้ตัดสิน สิ่งที่แก้ตอนอนุมัติ สิ่งที่ทำแล้ว)
  - SRS 3.1.1: แถว FR-EQS-11
  - constraints.md Q-SH-05: 001 เพิ่ม FR-EQS-11 แล้ว พร้อม endpoint และรหัส 403
  - CHANGELOG: รายการ Draft v5 + plan v1.7 + tasks v2.2 และเอาบรรทัดว่างเกินท้ายไฟล์ออก
  - plan v1.7: ข้อ 5 วิธีเรียกและแสดงลำดับคำขอ (200 / 403 / error อื่น), ข้อ 7 test ของ AC-EQS-17, ข้อ 10 เรื่องที่ spec ยังไม่บอก
  - tasks v2.2: เฟส 15 ใหม่ T-38 (หน้าจอเท่านั้น หลังบ้านเป็นของ 002), เฟสรอคำตอบเลื่อนเป็น 16, ตาราง AC / Constraint / FR
- สิ่งที่ AI ไม่ได้ทำ: ไม่แก้โค้ด ไม่ commit ไม่ได้แก้บรรทัด " 003-cancel-request" ที่เหลือค้างใน CHANGELOG (มาจาก commit ของ 003)
- เรื่องที่ spec ยังไม่บอก (เสนอทีม): ไม่มีคำขอค้าง หรือ 002 ไม่ตอบ หน้าจอแสดงข้อความอะไร (ตอนนี้ไม่แสดงส่วนนี้)

---

## 2569-11-12 /implement T-37 ตอบ error ทุกตัวด้วย {code, message} ตาม CON-TECH-03

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ของกลาง `backend/internal/platform/` และ `backend/internal/auth/middleware.go` เพราะ T-37 ของ 001)
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement T-37`
- ก่อนเริ่ม: ดึง main ล่าสุดแบบ rebase (commit README ของผู้สั่งยังไม่ได้ push)
- ไฟล์ที่สร้าง: `backend/internal/platform/apierror.go` (รหัส 4 ตัว, ข้อความภาษาไทยตาม plan ข้อ 5, ErrorBody, ErrorMessage, WriteError),
  `backend/internal/platform/apierror_test.go`
- ไฟล์ที่แก้: `backend/internal/auth/middleware.go` (ใช้ platform.WriteError ลบ writeError เดิม),
  `backend/internal/auth/middleware_test.go` (ตรวจ code และ message ใน 401, 403), `backend/internal/equipment/handler.go` (503 ใช้ platform.WriteError เลิกใช้ retryable)
- ไม่ได้แตะ `backend/internal/equipment/handler_test.go` (อยู่ในช่องไฟล์ที่แตะ แต่ test 503 อยู่ใน apierror_test.go แล้ว)
- test: `TestCON_TECH_03_ErrorBodyHasCodeAndMessage` ผ่าน (401 x2, 403 PRINCIPAL_INCOMPLETE, 503 รายการรุ่น และ 503 รายละเอียด
  ทุกกรณี body มีแค่ code และ message และ message ไม่ว่าง), `TestErrorMessage_AllCodesHaveMessage` ผ่าน, test ของ auth เดิม 3 ตัวผ่าน
  test หลังบ้านทั้งระบบผ่าน (รวม borrowrequest, cancelrequest ที่มี TestTC_, handover) หน้าจอ 23 test ผ่าน
- ผลต่อฟีเจอร์อื่น: auth.Middleware ใช้ร่วมโดย 002, 003, 005 ตอนนี้ 401/403 มี message เพิ่ม test ของทุกฟีเจอร์ตรวจแค่ code จึงยังผ่าน
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - test 503 ใช้ฐานข้อมูลจำลองที่ตอบ error ทุก query ไม่ต้องหยุด Docker
  - รหัสที่ยังไม่มีข้อความ ErrorMessage คืนรหัสเองแทน เพื่อให้ message ไม่ว่าง และ log เตือน
  - ไม่ลบ `ErrorDTO` ใน `equipment/dto.go` ที่เลิกใช้แล้ว เพราะ dto.go ไม่อยู่ในไฟล์ที่แตะของ T-37 (ลบได้ใน task ถัดไปที่แตะ dto.go)
- ข้อค้นพบ: ข้อความของ UPSTREAM_UNAVAILABLE ของ 005 ("ระบบไม่สามารถเชื่อมต่อฐานข้อมูลได้ กรุณาลองใหม่") ต่างจากของ 001
  รหัสเดียวกันควรใช้ข้อความเดียวกัน 002, 005 เรียก platform.WriteError ได้ถ้าทีมตกลง / 404, 405 ยังตอบข้อความธรรมดา (รอรหัสในรายการกลาง)

---

## 2569-11-12 /implement-all เฟส 9: T-17 แยก DTO ให้นักศึกษาไม่มีข้อมูลผู้ยืม และ API สรุปจำนวนตามสถานะสำหรับเจ้าหน้าที่

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ของกลาง `backend/internal/platform/apierror.go` เพิ่มรหัส FORBIDDEN_ROLE เพราะ T-17)
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 9` (task เดียวของเฟส)
- คำถามที่ AI ถามก่อนทำ และคำตอบของทีม (Pannathorn J):
  1. นักศึกษาเรียกสรุปได้ 403 รหัสอะไร plan ไม่ได้บอก -> ทีมตอบ: FORBIDDEN_ROLE จากรายการกลาง CON-TECH-03
     ข้อความ "เฉพาะเจ้าหน้าที่เท่านั้นที่ดูข้อมูลนี้ได้" (ทีมเลือกเอง ไม่ได้มาจาก spec) เพิ่มใน platform/apierror.go
  2. query สรุปต้องอยู่ใน repository.go ซึ่งไม่อยู่ในไฟล์ที่แตะ -> ทีมตอบ: เพิ่ม repository.go และ service.go
  เพิ่ม 3 ไฟล์นี้ในช่องไฟล์ที่แตะของ T-17 แล้ว
- ไฟล์ที่แก้: `backend/internal/equipment/dto.go` (StatusSummaryDTO, ลบ ErrorDTO ที่เลิกใช้ตั้งแต่ T-37),
  `repository.go` (CountByStatus), `service.go` (Summary ตรวจบทบาทที่ server, ErrStaffOnly), `handler.go` (GET /v1/equipment/summary),
  `handler_test.go`, `backend/internal/platform/apierror.go` (CodeForbiddenRole)
- test: `TestAC_EQS_07_StudentNoBorrower` ผ่าน (student-002 เปิดรายละเอียด CS-NB-0004, รายชิ้น และรายการรุ่น ไม่มีคำว่า borrower หรือ student-001),
  `TestFR_EQS_07_StaffSummaryByStatus` ผ่าน (เจ้าหน้าที่ได้ 3 สถานะ + total ตรงกับฐานข้อมูล / นักศึกษาได้ 403 FORBIDDEN_ROLE มีแค่ code และ message)
  test หลังบ้านทั้งระบบผ่าน หน้าจอ 47 test ผ่าน
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - response สรุปเป็น `{available, borrowed, unavailable, total}` plan ไม่ได้กำหนดรูปแบบ
  - สรุปนับทุกภาควิชา ตาม plan ข้อ 6.1 (ขอบเขตของเจ้าหน้าที่รอ Q-P3 ดู T-28)
  - ส่วนนักศึกษาไม่เห็นผู้ยืม: DTO ของนักศึกษาไม่มีช่องผู้ยืมอยู่แล้วตั้งแต่ T-05, T-12 task นี้เพิ่ม test ยืนยันทั้ง 3 endpoint
  - ยังไม่มีหน้าจอของสรุป (T-17 เป็นงานหลังบ้าน และ FR-EQS-07 ส่วนสรุปยังไม่มี AC)

---

## 2569-11-12 /implement-all เฟส 10: T-18 จำกัดทุก query ของนักศึกษาด้วยภาควิชา และตอบ 403 เมื่อขอข้ามภาควิชา

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ `frontend/src/api/client.js` ที่ 003 ใช้ร่วม เพิ่ม serverMessage แบบเดียวกับ api/handover.js ของ 005)
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 10` (task เดียวของเฟส)
- คำถามที่ AI ถามก่อนทำ และคำตอบของทีม (Pannathorn J):
  1. ต้องแตะ service.go, client.js, EquipmentDetail.jsx เพิ่ม -> ทีมตอบ: เพิ่มทั้ง 3 ไฟล์ (เพิ่มในช่องไฟล์ที่แตะแล้ว)
  2. เจ้าหน้าที่ส่ง ?department= มาต้องทำอะไร -> ทีมตอบ: ยังไม่ทำ เจ้าหน้าที่เห็นทุกภาควิชาตาม plan 6.1 รอ Q-P3
- ไฟล์ที่แก้:
  - `backend/internal/equipment/repository.go`: Filter.DepartmentID, เงื่อนไขภาควิชาในทุก query รายการ ($4), GetItem คืน department_id, ModelDepartment
  - `backend/internal/equipment/service.go`: scopeFor (นโยบายเดียว นักศึกษา = ภาควิชาตน / เจ้าหน้าที่ = ทุกภาควิชา / ไม่มีผู้ใช้ = ปฏิเสธ),
    ErrForbiddenDepartment, ListModels รับ requestedDept, ListItems และ GetItem ตรวจภาควิชาก่อนสร้าง DTO
  - `backend/internal/equipment/handler.go`: requestedDepartment, writeServiceError (403 FORBIDDEN_DEPARTMENT / 503)
  - `backend/internal/equipment/handler_test.go`, `frontend/src/api/client.js` (ApiError.serverMessage), `frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx` (ใช้ ErrorState)
- ไฟล์ที่สร้าง: `frontend/src/components/ErrorState.jsx`, `frontend/src/__tests__/AC-EQS-13.test.jsx`
- test: `TestAC_EQS_13_OtherDepartment403` ผ่าน (นักศึกษาภาควิชา 2: รายการเห็นแค่รุ่น 3 / รายละเอียด CS-NB-0001, รายชิ้นรุ่น 1,
  ?department=1, ?department=abc ได้ 403 FORBIDDEN_DEPARTMENT ไม่มีข้อมูลอุปกรณ์ใน body / ของภาควิชาตนยังเปิดได้),
  `TestFR_EQS_10_ScopeByRole` ผ่าน (นักศึกษาภาควิชา 1 ไม่เห็นรุ่น 3 / เจ้าหน้าที่เห็นทุกภาควิชา),
  `AC-EQS-13.test.jsx` 3 test ผ่าน / test หลังบ้านทั้งระบบผ่าน หน้าจอ 26 ไฟล์ 50 test ผ่าน (รวม TC-CXL ของ 003) `vite build` ผ่าน
- ลองจริง: หลังบ้านที่ port 8010 mock-student-101 เปิด CS-NB-0001 ได้ 403 + ข้อความไทย / รายการเห็นแค่รุ่น C / XX-NB-0001 ได้ 200
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - รุ่นหรือรหัสครุภัณฑ์ที่ไม่มีอยู่จริงยังได้รายการว่าง / 404 แบบเดิม ตรวจภาควิชาเฉพาะของที่มีอยู่
    (ข้อแลกเปลี่ยนตาม plan 6.1: 403 บอกได้ว่ารหัสนั้นมีอยู่จริง แต่ spec สั่งให้แจ้งไม่มีสิทธิ์)
  - หน้าจอแสดง message จากหลังบ้านผ่าน ErrorState ถ้าไม่มี message ใช้ "ดึงข้อมูลอุปกรณ์ไม่สำเร็จ" เดิม
  - หน้ารายการ (EquipmentList.jsx) ยังใช้ข้อความ error เดิม ไม่อยู่ในไฟล์ของ T-18 นักศึกษาจะไม่ได้ 403 จากรายการเว้นแต่ส่ง department เอง (T-20 แตะไฟล์นี้)

---

## 2569-11-12 /implement-all เฟส 11: T-19 บันทึก audit log เมื่อเปิดรายละเอียด

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (สร้างของกลาง `backend/internal/audit/` ที่ 003 รอใช้)
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 11` (task เดียวของเฟส)
- ไฟล์ที่สร้าง: `backend/internal/audit/writer.go` (Write, Entry, ค่าคงที่ action และ outcome, interface DB ที่ใช้ได้ทั้ง pool และ pgx.Tx),
  `backend/internal/audit/writer_test.go`
- ไฟล์ที่แก้: `backend/internal/equipment/handler.go` (getItem บันทึก audit ทุกการเปิดรายละเอียด, auditDetail)
- test: `TestAC_EQS_11_DetailWritesAudit` ผ่าน (มีแถว actor student-001, actor_role student, เวลาใกล้ปัจจุบัน, asset_code CS-NB-0002, outcome ok),
  `TestAudit_ForbiddenDetailIsRecorded` ผ่าน (403 ภาควิชาอื่นบันทึก outcome forbidden), `TestAudit_WriteFailsThenNoData` ผ่าน
  (เขียน audit ไม่สำเร็จได้ 503 ไม่มีข้อมูลอุปกรณ์), `TestWrite_WorksInTransaction` ผ่าน (ยกเลิก transaction แล้วไม่มีแถวค้าง / ช่องว่างเป็น NULL)
  test หลังบ้านทั้งระบบผ่าน หน้าจอ 50 test ผ่าน
- ลองจริง: หลังบ้านที่ port 8010 student-002 เปิด CS-NB-0004 ได้ 200 แล้วตาราง audit_log มีแถว detail ok ของ student-002
- ทำภายในไฟล์ที่กำหนดได้: test ของ AC-EQS-11 อยู่ใน audit/writer_test.go เรียกผ่าน handler จริง (แยก package audit_test)
  handler ได้ฐานข้อมูลสำหรับ audit จาก repository ใน NewHandler ไม่ต้องแก้ main.go
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - เขียน audit ไม่สำเร็จตอบ 503 ไม่ส่งข้อมูล (fail closed) ตามที่ plan ข้อ 6.4 เสนอ ยังรอยืนยัน Q-P4 (T-29 แก้ได้จุดเดียวใน getItem)
  - รหัสครุภัณฑ์ที่ไม่มีอยู่ (404) ไม่บันทึก เพราะไม่มีข้อมูลอุปกรณ์ถูกเข้าถึง plan ไม่ได้บอกกรณีนี้
  - Write ใช้ INSERT ... RETURNING ผ่าน QueryRow ทำให้ใช้ได้ทั้ง pool และ transaction โดยไม่ต้องเพิ่มเมธอด Exec ให้ repository
  - การดูรายการและการค้นหายังไม่บันทึก (รอ Q-11 ดู T-29) แต่มีค่าคงที่ ActionList, ActionSearch ไว้แล้ว

---

## 2569-11-12 /implement-all เฟส 12: T-20 ตอบ 503 เมื่อต่อฐานข้อมูลไม่ได้ และหน้าจอแสดงปุ่มลองใหม่โดยไม่แสดงข้อมูลเก่า

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 12` (T-37 ในเฟสนี้เสร็จไปก่อนแล้ว ทำแค่ T-20)
- ไฟล์ที่แก้: `backend/internal/equipment/handler_test.go`, `frontend/src/components/ErrorState.jsx` (ปุ่ม "ลองใหม่" เมื่อส่ง onRetry),
  `frontend/src/pages/EquipmentList/EquipmentList.jsx` (เก็บ error, ล้างรายการรุ่นและรายชิ้นเมื่อพลาด, retry โหลดใหม่ตามตัวกรองเดิม)
- ไฟล์ที่สร้าง: `frontend/src/__tests__/AC-EQS-08.test.jsx`
- ไม่ได้แตะ `backend/internal/platform/http.go` (อยู่ในช่องไฟล์ที่แตะ) เพราะ 503 ตามรูปแบบกลางทำไว้แล้วใน T-37 ไม่ต้องแก้เพิ่ม
- test: `TestAC_EQS_08_DBDown503` ผ่าน (pool ที่ชี้ไปฐานข้อมูลที่ไม่มีอยู่: รายการรุ่น รายชิ้น รายละเอียด ได้ 503 UPSTREAM_UNAVAILABLE
  มีแค่ code และ message ไม่มีข้อมูลอุปกรณ์ และมี Cache-Control: no-store), `AC-EQS-08.test.jsx` 2 test ผ่าน
  (เปิดหน้าตอนฐานข้อมูลล่มเห็นข้อความ + ปุ่มลองใหม่ ไม่ค้างที่กำลังโหลด / ล่มระหว่างใช้งานแล้วรายการเก่าหายไป กดลองใหม่กลับมาได้)
  test หลังบ้านทั้งระบบผ่าน หน้าจอ 27 ไฟล์ 52 test ผ่าน `vite build` ผ่าน
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - ป้ายปุ่ม "ลองใหม่" ใช้คำตาม AC-EQS-08 / ข้อความแสดง message จากหลังบ้าน ถ้าไม่มีใช้ "ดึงข้อมูลอุปกรณ์ไม่สำเร็จ" เดิม
  - หน้ารายละเอียดยังไม่มีปุ่มลองใหม่ (AC-EQS-08 พูดถึงหน้ารายการ) กดกลับไปรายการแล้วเปิดใหม่ได้
  - "ไม่แสดงหน้าจอค้าง" ตอนนี้ครอบคลุมกรณีหลังบ้านตอบ error ทันที กรณีหลังบ้านช้าจนไม่ตอบต้องมี timeout รอ Q-08 (T-26)

---

## 2569-11-12 /implement-all เฟส 13: T-21 จัดการ 401 บนหน้าจอ ไปหน้าล็อกอินจำลองพร้อมข้อความ แล้วกลับหน้าเดิม

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ `frontend/src/api/client.js` และ `frontend/src/App.jsx` ที่ 003 ใช้ร่วม)
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all ... เฟส 13` (task เดียวของเฟส)
- คำถามที่ AI ถามก่อนทำ และคำตอบของทีม (Pannathorn J):
  1. token มาจากไหน -> ทีมตอบ: เก็บใน browser (sessionStorage) เปิดครั้งแรกใช้ MOCK_TOKEN / VITE_MOCK_TOKEN เดิม ฟีเจอร์อื่นไม่ต้องเปลี่ยนวิธีลอง
  2. ไม่มี router จะกลับหน้าเดิมอย่างไร -> ทีมตอบ: จำหน้าไว้ในหน่วยความจำ returnTo รับเฉพาะชื่อหน้าในตาราง PAGES
  เพิ่ม App.jsx ในช่องไฟล์ที่แตะของ T-21 แล้ว (คำถามระบุไว้ว่าต้องแก้ App.jsx)
- ไฟล์ที่แก้: `frontend/src/api/client.js` (getToken, setToken, onUnauthorized, ส่ง token จาก session, แจ้งผู้ฟังเมื่อได้ 401 ทั้ง getJSON และ postJSON ของ 003),
  `frontend/src/App.jsx` (safePageKey, แสดง MockLogin เมื่อ 401 แล้วกลับหน้าเดิม), `backend/internal/auth/middleware_test.go`
- ไฟล์ที่สร้าง: `frontend/src/pages/Login/MockLogin.jsx` (เลือกผู้ใช้จำลอง 4 คนตาม mock_provider.go), `frontend/src/__tests__/AC-EQS-09.test.jsx`
- ไม่ได้แตะ `frontend/src/api/handover.js` ของ 005 ซึ่งยังใช้ MOCK_TOKEN ตายตัว (หน้าของ 005 ยังไม่อยู่ใน App) ควรแจ้ง 005 ให้เปลี่ยนไปใช้ getToken
- test: `AC-EQS-09.test.jsx` 3 test ผ่าน (ไม่มี token -> หน้าเข้าสู่ระบบพร้อมข้อความ "กรุณาเข้าสู่ระบบก่อนใช้งาน" -> เลือก student-001 -> เห็นรายการและส่ง token ใหม่ /
  หมดอายุตอนอยู่หน้า "คำขอที่รออนุมัติของฉัน" ล็อกอินแล้วกลับหน้านั้น / safePageKey ปฏิเสธ URL ภายนอก),
  `TestAC_EQS_09_ExpiredSessionIs401` ผ่าน / test หลังบ้านทั้งระบบผ่าน หน้าจอ 28 ไฟล์ 55 test ผ่าน `vite build` ผ่าน
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - ข้อความบนหน้าเข้าสู่ระบบใช้ message ของ 401 จากหลังบ้าน (CON-TECH-03) ถ้าไม่มีใช้ "กรุณาเข้าสู่ระบบก่อนใช้งาน" ข้อความเดียวกัน
  - หัวข้อ "เข้าสู่ระบบ (จำลอง)" และรายชื่อผู้ใช้จำลองเป็นของ prototype เท่านั้น ระบบจริงไปหน้าของมหาวิทยาลัย (Q-SH-03, T-30)
  - กลับหน้าเดิม = กลับแท็บเดิมใน App สถานะภายในหน้า (เช่น ตัวกรองที่เลือก) เริ่มใหม่ เพราะหน้าถูกสร้างใหม่หลังล็อกอิน

---

## 2569-11-12 ตรวจก่อนเปิด PR (T-37, T-17 ถึง T-21) และลงชื่อ 001 ในผู้ใช้ FORBIDDEN_ROLE

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ของกลาง `specs/000-shared/constraints.md` เพราะ T-17 ของ 001 ตอบ 403 `FORBIDDEN_ROLE`)
- เครื่องมือ: Claude Code
- คำสั่ง: "เปิด PR เข้า main ต้องใส่อะไรบ้าง"
- สิ่งที่ AI ตรวจ: PR ไม่ชนกับ main, `go vet` และ `go test ./...` ผ่าน, หน้าจอ 28 ไฟล์ 55 test ผ่าน
- สิ่งที่ AI พบ: ช่อง "ใครใช้แล้ว" ของ `FORBIDDEN_ROLE` มีแค่ 004 แต่ 001 ใช้แล้วตั้งแต่ T-17 (API สรุปตามสถานะ)
- ไฟล์ที่แก้ (คนแก้เอง): `specs/000-shared/constraints.md` แถว `FORBIDDEN_ROLE` (เพิ่ม 001) และแถว CON-TECH-03 ในตาราง "ใครอ้างใช้บ้าง"
- ไม่เปลี่ยนเนื้อหา Constraint จึงไม่ทำ CR และไม่ขึ้นเวอร์ชัน (แบบเดียวกับที่ 005 ลงชื่อใน v5.1)
- ต้องแจ้ง 003, 005: body ของ 401 / 403 มี `message` แล้ว (T-37) และหน้าจอไปหน้าล็อกอินจำลองเมื่อได้ 401 (T-21)

---

## 2569-11-12 /implement-all เฟส 15: T-38 แสดงคำขอยืมที่ค้างอยู่บนหน้ารายละเอียด (เฉพาะเจ้าหน้าที่)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ `frontend/src/api/client.js` ที่ 003 ใช้ร่วม เพิ่ม listItemRequests เท่านั้น ไม่เปลี่ยนฟังก์ชันเดิม)
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement-all specs/001-equipment-status/tasks.md เฟส 15` (task เดียวของเฟส)
- คำถามที่ AI ถามก่อนทำ และคำตอบของทีม (Pannathorn J):
  1. หน้ารายละเอียดเรียก endpoint ลำดับคำขอเพิ่มทุกครั้ง ทำให้ `AC-EQS-05.test.jsx` ที่นับ request ทั้งหมด (คาด 2) ไม่ผ่าน
     และไฟล์นี้ไม่อยู่ในไฟล์ที่แตะของ T-38 -> ทีมตอบ: เพิ่มไฟล์นี้ในไฟล์ที่แตะ และให้นับเฉพาะ request ของรายละเอียด
     (ยังตรวจว่ารีเฟรชแล้วถามหลังบ้านใหม่ 2 ครั้ง และทุก request เป็น no-store เหมือนเดิม)
- ไฟล์ที่สร้าง: `frontend/src/components/RequestQueue.jsx` (ตาราง ลำดับที่ / เวลาที่ส่ง / ชื่อผู้ขอ เรียงตาม queuePosition
  พร้อมหมายเหตุว่าชื่อเป็นชื่อที่ผู้ขอพิมพ์เอง), `frontend/src/__tests__/AC-EQS-17.test.jsx`
- ไฟล์ที่แก้: `frontend/src/api/client.js` (listItemRequests เรียก `GET /v1/equipment/items/{assetCode}/requests` ของ 002 ตาม CR-005),
  `frontend/src/pages/EquipmentDetail/EquipmentDetail.jsx` (โหลดคำขอแยกจาก getItem ทุกครั้งที่เปิดหน้า),
  `frontend/src/__tests__/AC-EQS-05.test.jsx` (ตามคำตอบข้อ 1), `specs/001-equipment-status/tasks.md` (ไฟล์ที่แตะและสถานะ T-38)
- test: `AC-EQS-17.test.jsx` 4 test ผ่าน (เจ้าหน้าที่เห็น 2 คำขอ ลำดับ 1 = A 2569-11-10 10:30, 2 = B 2569-11-10 11:15 พร้อมชื่อ /
  นักศึกษาได้ 403 FORBIDDEN_ROLE ไม่เห็นส่วนนี้ ไม่มีชื่อ ไม่มีข้อความ error / 002 ตอบ 503 สถานะและกำหนดคืนยังแสดง /
  ไม่มีคำขอค้างไม่แสดงส่วนนี้) หน้าจอทั้งระบบ 29 ไฟล์ 59 test ผ่าน `vite build` ผ่าน หลังบ้าน `go test ./...` ผ่าน (ไม่ได้แก้หลังบ้าน)
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - หัวข้อ "คำขอยืมที่ค้างอยู่" หัวคอลัมน์ "ลำดับที่", "เวลาที่ส่ง", "ชื่อผู้ขอ" และหมายเหตุ "ชื่อผู้ขอคือชื่อที่ผู้ขอพิมพ์เองในฟอร์มคำขอ ระบบไม่ได้ตรวจ"
    ใช้คำจาก FR-EQS-11 ไม่มี mockup
  - แสดงส่วนนี้เฉพาะเมื่อรายละเอียดโหลดสำเร็จ ถ้ารายละเอียดเป็น error (เช่น 503 หรือต่างภาควิชา) ไม่แสดงคำขอด้วย
  - หน้าจอเรียงตาม queuePosition ซ้ำ แม้ 002 ส่งมาเรียงแล้ว
  - client ที่ไม่มี listItemRequests (client จำลองของ test เดิม) นับเป็นกรณี 002 ไม่ตอบ จึงไม่ต้องแก้ test เดิมตัวอื่น
- ยังลองกับระบบจริงไม่ได้: endpoint ของ 002 (T-18 ของ 002) ยังไม่มีใน main ตอนนี้หน้าจะได้ 404 และไม่แสดงส่วนนี้

---

## 2569-11-12 ปิด Q-08 (timeout ของ FR-EQS-08 / AC-EQS-08)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status
- เครื่องมือ: Claude Code
- คำสั่ง: "003 และ 005: PR #14 ทำให้ 401 / 403 มี message ภาษาไทย และหน้าจอไปหน้าล็อกอินจำลองเมื่อได้ 401 (มีผลกับแท็บของ 003 ด้วย) test ของทุกฟีเจอร์ยังผ่าน มีผลยังไง"
  -> AI วิเคราะห์ผลกระทบ ไม่แก้ไฟล์ สรุป: 003 ได้พฤติกรรมของ AC-CXL-07 จาก App แล้ว แต่ T-26 ของ 003 ยังไม่มี vitest
  และ plan 003 ข้อ 129 เขียน returnTo เป็น path ส่วน 005 แสดง message ใหม่ในฟอร์ม แต่ไม่ไปหน้าล็อกอิน (handover.html แยกจาก App)
  test ผ่านเพราะ test ของ 003 / 005 ใช้ client จำลอง และตรวจแค่ code
- คำสั่ง: "ทำ T26 31 32 ไปแล้วใช่มั้ย ของ001" -> AI ตอบ: T-31, T-32 เสร็จ T-26 รอ Q-08
- คำสั่ง: "คำถาม Q-08 ใน spec.md:166 ควรตอบว่าอะไร" -> AI เสนอ 5 วินาที (เท่ากับ CancelDBTimeout ของ 003, มากกว่า p95 3 วินาที
  ของ NFR-PERF-01) และหน้าจอรอนานกว่าหลังบ้านเล็กน้อย (6 วินาที) ตัวเลือกอื่น 3 วินาที / 10 วินาที
- คำตอบของทีม (Pannathorn J): "งั้นเอาที่แนะนำเลย" -> หลังบ้าน 5 วินาที หน้าจอ 6 วินาที
- ไฟล์ที่แก้: `specs/001-equipment-status/spec.md` (v5.1 ปิด Q-08 + Clarifications), `plan.md` (v1.8 ข้อ 6.2),
  `tasks.md` (v2.3 T-26 พร้อมทำ), `specs/CHANGELOG.md`, `specs/README.md` (เวอร์ชัน), `docs/srs/SRS.md` (ตาราง 4.3 แถว FR-EQS-08)
- ไม่ได้แก้โค้ด ตัวเลขจะเข้า `config.go` / `config.js` ตอนทำ T-26
- Q-08 เดิมให้ถามนักคอมพิวเตอร์ ทีมตัดสินเองแทน

---

## 2569-11-12 /implement T-26 ตั้ง timeout ของ query และ request (001)

- ผู้สั่ง: Pannathorn J
- ฟีเจอร์: 001-equipment-status (แก้ `frontend/src/api/client.js` ที่ 002, 003 ใช้ร่วม: เพิ่มตัวเลือก timeout ใน getJSON
  ใช้เฉพาะฟังก์ชันของ 001 request ของ 002, 003 ทำงานเหมือนเดิม)
- เครื่องมือ: Claude Code
- คำสั่ง: `/implement T-26`
- คำถามที่ AI ถามก่อนทำ และคำตอบของทีม (Pannathorn J):
  1. ช่อง "ไฟล์ที่แตะ" มีแค่ config.go, config.js, client.js ไม่พอทำ timeout ของฐานข้อมูลและไม่มีที่วาง test
     -> เพิ่ม `backend/internal/equipment/handler.go`, `handler_test.go`, `frontend/src/__tests__/AC-EQS-08.test.jsx`
  2. client.js ใช้ร่วมกับ 003 (ยกเลิกคำขอรอส่ง Email ยอมช้าถึง 15 วินาที NFR-PERF-CXL-01) timeout 6 วินาทีใช้กับ request ไหน
     -> เฉพาะ request ของ 001 (listModels, listItems, getItem)
  3. test ฐานข้อมูลช้าทำแบบไหน -> ล็อกตาราง equipment_item จริงใน testcontainers ใช้ 5 วินาทีจริงจาก config
- ไฟล์ที่แก้:
  - `backend/internal/config/config.go`: `EquipmentDBTimeout = 5 * time.Second` (FR-EQS-08, Q-08)
  - `backend/internal/equipment/handler.go`: `withDBTimeout` ห่อ router ของ 001 ทุก endpoint ด้วย context.WithTimeout
    เกินเวลา query คืน error แล้วตอบ 503 `UPSTREAM_UNAVAILABLE` ผ่านทางเดิม
  - `backend/internal/equipment/handler_test.go`: `TestAC_EQS_08_SlowDBTimeout`
  - `frontend/src/config.js`: `EQUIPMENT_REQUEST_TIMEOUT_MS = 6000` (FR-EQS-08, Q-08)
  - `frontend/src/api/client.js`: getJSON รับ `{ timeoutMs }` ใช้ AbortController ส่งเฉพาะจาก listModels, listItems, getItem
  - `frontend/src/__tests__/AC-EQS-08.test.jsx`: 2 test ใหม่ (หน้าจอเลิกรอเมื่อครบ 6 วินาที / request ของ 003 ไม่ถูกตัด)
  - `specs/001-equipment-status/tasks.md`: ไฟล์ที่แตะและสถานะของ T-26
- test: `TestAC_EQS_08_SlowDBTimeout` ผ่าน (รอครบ 5 วินาทีแล้วได้ 503 พร้อม message ไม่ส่งรายการ) `go vet ./...` ผ่าน
  `go test ./...` ผ่านทุก package / vitest 29 ไฟล์ 61 test ผ่าน `vite build` ผ่าน
- สิ่งที่ AI เลือกเอง (ทีมต้องรู้):
  - หน้าจอเลิกรอแล้วแสดงข้อความเดิมของหน้าจอ "ดึงข้อมูลอุปกรณ์ไม่สำเร็จ" พร้อมปุ่มลองใหม่ (แบบเดียวกับต่อหลังบ้านไม่ได้)
    ไม่สร้างข้อความหรือรหัสใหม่
  - listItemRequests (endpoint ของ 002 ที่หน้ารายละเอียดของ 001 เรียก) ไม่ได้ใส่ timeout ตามคำตอบข้อ 2
  - การเขียน audit ของหน้ารายละเอียดอยู่ใน context เดียวกัน จึงถูกจำกัด 5 วินาทีด้วย
- คำเตือน act ใน `setup.test.jsx` มีอยู่ก่อน T-26 ไม่ได้แก้

---

## 2569-11-11 คำสั่ง: เริ่ม 004-approve-request ใหม่จากศูนย์ (spec Draft v1)

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code
- คำสั่ง: "อยากลบทำใหม่อะ มันงง" แล้วเลือก "เริ่ม 004 จากศูนย์"
- สิ่งที่ทำ:
  - เปลี่ยนชื่อ branch เดิมเป็น `004-approve-request-old` (เก็บไว้เป็นสำรอง มี spec v2.1 / plan v1.1 / tasks v1.1 และบันทึก prompt-log ของ 004 ก่อนหน้า ไม่ได้ลบ branch บน GitHub)
  - สร้าง branch `004-approve-request` ใหม่จาก origin/main (743e36c)
  - เขียน specs/004-approve-request/spec.md Draft v1 จาก SRS 3.1.4 (ยังเป็น `[เติม]`), RQ-14, RQ-21, RQ-25, BL-08 และไฟล์ร่วมล่าสุด (constraints.md v5, data-model.md v9, migration 0003 ของ 002)
  - เพิ่ม 004 ในตาราง "ใครอ้างใช้บ้าง" ของ constraints.md
- ผลลัพธ์: FR-APV-01, FR-APV-02, NFR-SEC-01, AC-APV-01 ถึง 05 (ID รอทีมอนุมัติ), ASM-01 ถึง 03, Q-APV-01 ถึง 10
- สิ่งที่ AI ไม่เดา: เนื้อหา Email, เวลา "ทันที", ที่มาของ Email ผู้ขอ, ลำดับการอนุมัติ, ข้อมูลบนแดชบอร์ด, ประวัติ, กรณีพิเศษ -> เขียนเป็น Open Questions ให้ /clarify ถาม
- สิ่งที่ AI ใช้จากไฟล์ร่วม (ไม่ใช่การเดา): ค่าสถานะ requested / approved / rejected (data-model.md v9), คำขออื่นถูกปฏิเสธอัตโนมัติ (Q-SH-06 ในconstraints.md), CON-TECH-03, คอลัมน์ของ borrow_request (0003)
- ยังไม่ได้แก้: SRS 3.1.4 (รอเจ้าของสั่งหลัง /clarify)
- ขั้นต่อไป: /clarify

---

## 2569-11-11 คำสั่ง: /clarify ของ 004-approve-request (spec Draft v1) รอบที่ 1

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/clarify.md เพราะ /clarify ไม่ได้ติดตั้งใน session)
- คำสั่ง: `/clarify`
- AI ทำหน้าเว็บให้ทีมตอบ: https://claude.ai/artifact/Da6xDoDVyThdMF571RRVUH

### คำถามที่ AI ถาม
1. 004 ต้องจัดการคำขอยืมแบบค้างคืนด้วยไหม หรือทำเฉพาะยืมปกติ
2. เงื่อนไขที่อนุมัติไม่ได้ และข้อความของแต่ละกรณี (อุปกรณ์ไม่ว่าง, ยกเลิกแล้ว, พิจารณาแล้ว, กดพร้อมกัน)
3. ต้องอนุมัติตามลำดับเวลาที่ส่งไหม ข้ามได้ไหม ต้องระบุเหตุผลไหม
4. ปฏิเสธต้องระบุเหตุผลไหม คำขอที่ถูกปฏิเสธตามต้องส่ง Email ไหม บันทึกใครเป็นผู้ดำเนินการ
5. คำขอที่ไม่มีใครพิจารณาค้างได้นานแค่ไหน และระบบต้องทำอะไร
6. "ทันที" กี่วินาที เนื้อหา Email และกรณีส่งไม่สำเร็จ
7. Email ของผู้ขอมาจากไหน
8. หน้ารายการคำขอแสดงคำขอใด เรียงอย่างไร แสดงข้อมูลอะไร
9. อนุมัติ ID ของ AC-APV-01 ถึง 05 ไหม และเรื่องอื่น

### คำตอบของทีม (คัดลอกจากหน้าเว็บ)
1. ต้องทำด้วยหรอ ไม่ทำเป็นอะไรมั้ย
2. ไม่ให้ดำเนินการทุกกรณีที่คำขอไม่ใช่ requested หรืออุปกรณ์ไม่ว่าง ใช้ ITEM_NOT_AVAILABLE กับรหัสที่ตั้งเอง และคนที่กดก่อนชนะ
3. ต้องอนุมัติตามลำดับเวลาที่ส่ง ไม่มีการข้ามคำขอ มีแต่ยกเลิกคำขอ
4. ระบุเหตุผล บันทึกว่าเจ้าหน้าที่ที่กดอนุมัติเป็นคนปฏิเสธ แล้วส่ง Email แจ้งด้วย
5. ค้างได้จนถึงเที่ยงคืนก่อนระบบจะรีเซ็ทคำขอหลังเที่ยงคืน
6. 60 วินาที / Email อนุมัติ: หัวข้อ "แจ้งผลการอนุมัติการยืมอุปกรณ์ภาคคอมพิวเตอร์", ชื่อผู้ยืม, รายการอุปกรณ์, จำนวนอุปกรณ์, วัตถุประสงค์, วันเวลาที่ต้องการยืม, วันเวลาที่กำหนดคืน, ชื่อผู้อนุมัติ, สถานที่รับ "คณะวิทยาศาสตร์ ตึก 1 ชั้น 2 ห้องอุปกรณ์ภาควิชาคอมพิวเตอร์" / Email ปฏิเสธ: หัวข้อ "แจ้งผลการปฏิเสธการยืมอุปกรณ์ภาคคอมพิวเตอร์", ชื่อผู้ยืม, รายการอุปกรณ์, วันที่ยื่นคำขอ, เหตุผลในการปฏิเสธ, ชื่อผู้พิจารณา / ส่งไม่สำเร็จ: บันทึก Failed ไม่เปลี่ยนผล, บันทึกสาเหตุและเวลา, ส่งซ้ำอัตโนมัติ 3 ครั้งเว้นระยะ, ยังไม่สำเร็จแจ้งผู้ดูแลระบบ, ผู้ดูแลระบบกดส่งซ้ำได้โดยไม่พิจารณาใหม่
7. Email ของมหาวิทยาลัยที่ login เข้าสู่ระบบ
8. แสดงทุกคำขอ requested ของภาควิชาตนเอง เรียงตามเวลาที่ส่ง พร้อมชื่อผู้ขอ อุปกรณ์ วัตถุประสงค์ และช่วงเวลาที่ต้องการยืม
9. ทีมอนุมัติ

### ข้อสังเกตของ AI (ยังไม่แก้ spec เพราะคำตอบหลายข้อขัดกับ spec หรือไฟล์ร่วม ต้องถามซ้ำตามกฎของ /clarify)
- ข้อ 1 เป็นคำถามกลับ ยังไม่ได้ตัดสิน
- AC-APV-01 ใน spec ถูกแก้ด้วยมือเป็น "สถานะของอุปกรณ์รายชิ้นเป็น ไม่ว่าง" ขัดกับ CON-DATA-01 (ไม่มีสถานะ "ไม่ว่าง"), FR-APV-01 (ข) และ data-model.md แถว UC-12
- ข้อ 2 ยังไม่มีข้อความของแต่ละกรณี และรหัสที่ตั้งเองต้องลงทะเบียนตาม CON-TECH-03
- ข้อ 3 "มีแต่ยกเลิกคำขอ" ไม่ชัดว่าเจ้าหน้าที่ทำอะไรกับคำขอแรกที่ไม่ผ่าน (ยกเลิก UC-15 เป็นของนักศึกษา 003)
- ข้อ 4 ไม่บอกเหตุผลที่ใช้กับคำขอที่ถูกปฏิเสธตาม
- ข้อ 5 "รีเซ็ท" ไม่บอกว่าเปลี่ยนเป็นสถานะใด เที่ยงคืนของวันไหน และต้องแจ้งผู้ขอไหม
- ข้อ 6 "จำนวนอุปกรณ์" ขัดกับ 1 คำขอ = 1 ชิ้น (DM-06, FR-BRQ-01), "กำหนดคืน" 005 เป็นผู้บันทึกตอนส่งมอบ (CR-003), "ผู้ดูแลระบบ" ไม่มีบทบาทนี้ในระบบ, ไม่บอกระยะห่างของการส่งซ้ำ, "ผู้ยืม" ตาม glossary คือคนที่มี Loan ควรใช้ "ผู้ขอ"
- ข้อ 7 IF-IDP-01 และ DM-08 ไม่มี Email และตอนอนุมัติคนที่ login คือเจ้าหน้าที่ ไม่ใช่ผู้ขอ ต้องให้ 002 เก็บ Email ตอนส่งคำขอ (กระทบไฟล์ร่วมและ 002)
- ข้อ 9 เรื่องรหัส error กรณีนักศึกษาและข้อมูลประวัติ RQ-25 ยังไม่ตอบ

### /clarify รอบที่ 2 (2569-11-11)
คำถามที่ AI ถามซ้ำ (หน้าเว็บเดิม): R1 ค้างคืน, R2 สถานะอุปกรณ์หลังอนุมัติ, R3 ข้อความแต่ละกรณี, R4 คำขอแรกไม่ผ่าน, R5 เหตุผลของคำขอที่ถูกปฏิเสธตาม, R6 รีเซ็ทหลังเที่ยงคืน, R7 ข้อมูลใน Email ที่ขัดกับไฟล์ร่วม, R8 เก็บ Email ผู้ขอ, R9 รหัส error นักศึกษาและประวัติ

คำตอบของทีม:
- R1: ข. รับคำขอค้างคืนด้วย
- R2: approve แล้วต้องเป็นถูกยืม
- R3: (ทีมตอบเป็นรายการ Email แจ้งยกเลิกคำขออัตโนมัติ 10 ข้อ ซ้ำกับ R6 ไม่ได้ตอบข้อความ error)
- R4: ก. เจ้าหน้าที่ปฏิเสธคำขอแรกพร้อมเหตุผล แล้วคำขอถัดไปกลายเป็นคำขอแรก
- R5: ก. ข้อความตายตัว "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว"
- R6: เปลี่ยนเป็น rejected / เที่ยงคืนของวันที่ส่งคำขอ / Email: หัวข้อ "แจ้งยกเลิกคำขอยืมอุปกรณ์อัตโนมัติ", ผู้รับ Email ผู้ยื่นคำขอ, Request ID, ชื่อผู้ยื่น, ชื่อและจำนวนอุปกรณ์, วันที่ยื่น, สถานะ "ยกเลิกอัตโนมัติ (Cancelled)", เหตุผล "คำขอยังคงค้างอยู่และไม่ได้รับการพิจารณาภายในเวลาที่กำหนด", เวลา 00:00 น., แนวทาง "หากยังต้องการยืมอุปกรณ์ สามารถยื่นคำขอใหม่ได้" / ผู้ดำเนินการ system / อนุมัติคำขอที่เลยเวลาเริ่มยืมหรือเลย 16:00: ไม่ได้
- R7: จำนวนอุปกรณ์: ยืมปกติไม่ใส่ ค้างคืนใส่ได้ / กำหนดคืนใช้เวลาสิ้นสุดที่ผู้ขอระบุ: ใช่ / ส่งซ้ำห่างกัน 1 นาที / "ผู้ดูแลระบบ" = เจ้าหน้าที่ภาควิชาเจ้าของอุปกรณ์ แจ้งทาง Email
- R8: รอ 002 ตอบ (AI ร่างข้อความถึง 002: https://claude.ai/artifact/Y26U9FzkA3uDEAsdJztWzG)
- R9: ขอให้ขยายความ

ข้อสังเกตของ AI (ยังไม่แก้ spec ต้องถามรอบที่ 3):
- R1 ขัดกับขอบเขต prototype (specs/README.md และ SRS 3.1.7 ให้ UC-11, UC-14 ไม่ทำใน prototype, SRS 3.1.4 และ BL-08 เป็นยืมปกติ) ต้องให้ทีมตัดสินและทำ CR และ /clarify เพิ่ม FR เองไม่ได้ และขัดกับ R6 (คำขอค้างคืนรออาจารย์เซ็นหลายวัน แต่จะถูกปฏิเสธตอนเที่ยงคืน)
- R2 ขัดกับ DOM-REQ-01 และ data-model.md (UC-12 ว่างไม่เปลี่ยน, UC-08 ว่าง -> ถูกยืม, "ถูกยืม" = มี Loan) ต้องทำ CR และกระทบ 001 002 005
- R3 ยังไม่ได้คำตอบเรื่องข้อความ error
- R6 สถานะ rejected แต่ Email เขียนว่า "ยกเลิกอัตโนมัติ (Cancelled)" ซึ่ง cancelled เป็นของ 003 / "จำนวน" ขัดกับ 1 คำขอ = 1 ชิ้น / "ไม่ได้" ยังไม่ชัดว่าเลยเวลาเริ่มยืม หรือเลย 16:00
- R7 จำนวนอุปกรณ์ของค้างคืนขัดกับ FR-BRQ-01 (รายชิ้นทั้ง 2 แบบ) / แจ้งเจ้าหน้าที่ทาง Email ต้องมี Email เจ้าหน้าที่ ซึ่งไม่มีในระบบเหมือน R8

### /clarify รอบที่ 3 (2569-11-11)
คำถามที่ AI ถาม (หน้าเว็บเดิม): S1 ค้างคืน, S2 สถานะ "ถูกยืม", S3 ข้อความ error (AI ร่าง 7 ข้อความ), S4 คำใน Email เที่ยงคืนและเวลาที่อนุมัติไม่ได้, S5 จำนวนอุปกรณ์, S6 แจ้งเจ้าหน้าที่เมื่อ Email ล้มเหลว, S7 รหัส error นักศึกษาและประวัติ

คำตอบของทีม:
- S1: ผู้ขอต้องไปขอลายเซ็นจากอาจารย์ก่อน ค่อยนำใบมาแนบตอนส่งคำขอ ระบบไม่ต้องรอหลายวัน
- S2: ก. ยังเป็น "ว่าง" จนกว่า 005 ส่งมอบ ไม่ต้องทำ CR
- S3: ตามนี้ได้เลย
- S4: ก. ใช้ "ปฏิเสธอัตโนมัติ" ทั้งหัวข้อและสถานะ (ไม่ได้ตอบเรื่องเวลา)
- S5: ข. ค้างคืนต้องขอหลายชิ้นได้ (ต้องคุยกับ 002)
- S6: ก. รายการ "Email ที่ส่งไม่สำเร็จ" บนหน้าของ 004 พร้อมปุ่มส่งซ้ำ
- S7: ลงทะเบียน FORBIDDEN_ROLE (403) "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้" / ประวัติ: รหัสคำขอ, ผล, เหตุผล, ผู้ดำเนินการ, วันเวลา

### สิ่งที่แก้หลัง /clarify 3 รอบ
- spec.md Draft v1 -> Draft v2: FR-APV-01 (ก) ถึง (ฌ) และ FR-APV-02 (ก) ถึง (ง) ตามคำตอบ (ID เดิม ไม่เพิ่ม FR), แก้ AC-APV-01 กลับเป็น "ยังเป็น ว่าง" ตาม S2, แก้ AC-APV-03 ถึง 05, เสนอ AC-APV-06 ถึง 15 (ID รอทีมอนุมัติ Q-APV-12), เพิ่ม ASM-04 ถึง 09, ปิด Q-APV-02, 03, 05 ถึง 10, เปิด Q-APV-01 (ค้างคืน), Q-APV-04 (Email ผู้ขอ), Q-APV-11 (เวลาที่อนุมัติไม่ได้), Q-APV-12, หัวข้อ Clarifications
- constraints.md: ลงทะเบียน `FORBIDDEN_ROLE` (403) ในตารางรหัส CON-TECH-03 ตาม S7
- สิ่งที่ AI ทำเอง (ทีมตรวจ): ใช้คำว่า "ผู้ขอ" แทน "ผู้ยืม" ตาม glossary (ผู้ยืม = คนที่มี Loan); ไม่ใส่ค้างคืนใน spec แม้ทีมตอบว่ารับ เพราะขัดกับไฟล์ร่วมหลายจุดและ /clarify เพิ่ม FR เองไม่ได้ (เขียนเป็น Q-APV-01); "เที่ยงคืนของวันที่ส่ง" เขียนเป็น 00:00 น. หลังวันที่ส่ง; ASM-06 สิทธิ์พิจารณาเฉพาะภาควิชาตนเอง อนุมานจากคำตอบหน้ารายการ (Q8) และ AC-APV-12 ใช้ FORBIDDEN_DEPARTMENT

### ตอบเพิ่มในแชต (2569-11-11)
- Q-APV-11 อนุมัติไม่ได้หลังเลยเวลาไหน -> ทีม: "เลย 16:00"
- Q-APV-12 อนุมัติ ID ของ AC-APV-06 ถึง 15 -> ทีม: "อนุมัติ"
- สิ่งที่แก้: spec.md FR-APV-01 (ค) อนุมัติไม่ได้เมื่อเลย 16:00 น. ของวันที่เริ่มยืม, เอาป้าย [เสนอ] ออกจาก AC-APV-06 ถึง 15, ปิด Q-APV-11 และ 12, เสนอ AC-APV-16 (กรณีเลย 16:00) และเปิด Q-APV-13
- สิ่งที่ AI ตีความ (ทีมตรวจ): "16:00 ของวันที่เริ่มยืม" (ทีมไม่ได้บอกว่าวันไหน) / ข้อความ "เลยเวลาเริ่มยืมแล้ว อนุมัติไม่ได้" ที่ทีมรับใน S3 ไม่ตรงกับเงื่อนไขใหม่ AI ไม่แก้เอง ถามใน Q-APV-13
- ตอบเพิ่ม: "1-2 ยืนยัน" -> (1) 16:00 น. ของวันที่เริ่มยืม ถูกต้อง (2) ใช้ข้อความ "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้" แก้ใน FR-APV-01 (ซ) และ Q-APV-13 / ID ของ AC-APV-16 ยังไม่ได้อนุมัติ

---

## 2569-11-11 คำสั่ง: /plan ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/plan.md เพราะ /plan ไม่ได้ติดตั้งใน session) ผ่าน plan mode
- คำสั่ง: "พร้อม" (ให้ทำ /plan ต่อจาก spec Draft v2)
- ผลลัพธ์: specs/004-approve-request/plan.md (Draft v1) 8 หัวข้อ, ตารางตรวจ Constraints ครบ 8 ข้อ (CON-TECH-01/02/03, CON-DATA-01, DOM-REQ-01, DOM-LOAN-01, IF-IDP-01, IF-NOT-01), ค่าที่ตั้งได้ 13 รายการ, แผนทดสอบ TestAC_APV_01 ถึง 16, ลำดับงาน 8 ขั้น; บันทึก CHANGELOG
- Constraint ที่ยังไม่ได้ใช้: ไม่มี (IF-IDP-01 ใช้ auth จำลอง, IF-NOT-01 ใช้ตัวส่งจำลองและตัวหา Email จำลอง)
- AC ที่ทดสอบยาก: AC-APV-10 (กดพร้อมกัน ต้องใช้ goroutine กับ PostgreSQL จริง), AC-APV-04, 13, 14, 16 (60 วินาที, 00:00, ส่งซ้ำห่าง 1 นาที, 16:00 ต้องใช้นาฬิกาจำลอง ทดสอบกับระบบ Email จริงไม่ได้)
- สิ่งที่ AI อยากเดาแต่ไม่เดา: ที่มาของ Email ผู้ขอ (Q-APV-04), รับค้างคืน (Q-APV-01), ชื่อเจ้าหน้าที่ (Principal ไม่มีชื่อ เจอตอน plan), เลข migration, แถว action ใน DM-07
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ตาราง approval_decision และ approval_email, path `/v1/approvals/...`, รหัส error ใหม่ 5 ตัว (REQUEST_CANCELLED, REQUEST_ALREADY_DECIDED, REQUEST_DECIDED_CONCURRENTLY, NOT_FIRST_IN_QUEUE, APPROVAL_CUTOFF_PASSED) ต้องลงทะเบียนใน CON-TECH-03 ตอน implement, ใช้ VALIDATION_FAILED กับเหตุผลว่าง, แยกพิจารณาแล้ว/กดพร้อมกันด้วย loaded_at, ส่ง Email ผ่านตาราง outbox และ goroutine
- เรื่องที่พบตอน plan: auth.Middleware ของกลางตอบ error ไม่มี message (CON-TECH-03) ต้องแจ้ง 001

---

## 2569-11-12 คำสั่ง: อ่านไฟล์และตอบคำถามค้างของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004) / คำตอบในหน้าเว็บบันทึกจากบัญชีชื่อ "Gojo"
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code
- คำสั่ง (ตามลำดับ):
  1. "อ่าน `AGENTS.md`, `specs/000-shared/` (ทั้ง 3 ไฟล์) และ SRS หัวข้อ 3.0 กับหัวข้อของ004-approve-request"
  2. "อ่าน specs/004-approve-request/spec.md กับ plan.md"
  3. "ขอ html สำหรับตอบ"
  4. "ตอบแล้ว"
- ผลลัพธ์ของข้อ 1-2: AI สรุปและพบประเด็นที่ spec ไม่ได้บอกหรือ plan ยังไม่ได้ตัดสิน 8 ข้อ (ไม่แก้ไฟล์)
- ผลลัพธ์ของข้อ 3: หน้าเว็บให้ทีมตอบ https://claude.ai/artifact/LdTfyp1kd4ZEEpYnUhU6pc

คำถามที่ AI ถาม และคำตอบของทีม:
- Q-APV-14 ปฏิเสธอัตโนมัติใช้กับคำขอค้างคืนไหม -> เที่ยงคืน (ฉ) เฉพาะยืมปกติ / ปฏิเสธตาม (จ) รวมคำขอค้างคืนของชิ้นเดียวกันด้วย
- Q-APV-15 เที่ยงคืนนับจากวันที่ส่งหรือวันที่เริ่มยืม -> วันที่ส่ง (ตามเดิม)
- Q-APV-16 คำขอแรกที่เลย 16:00 น. ขวางคิว -> ยอมรับ เจ้าหน้าที่ปฏิเสธเองพร้อมเหตุผล
- Q-APV-17 "ส่งซ้ำ 3 ครั้ง" รวมครั้งแรกไหม -> ครั้งแรก + ซ้ำ 3 = 4 ครั้ง
- Q-APV-18 ข้อความ FORBIDDEN_DEPARTMENT -> "คำขอนี้เป็นของอุปกรณ์ภาควิชาอื่น คุณไม่มีสิทธิ์พิจารณา" (AI ร่าง)
- Q-APV-19 ลงทะเบียนรหัส error 5 ตัว -> ตามชื่อใน plan และเพิ่ม 004 ในรหัสเดิม 4 ตัว
- Q-APV-20 003 ยกเลิกพร้อมกับ 004 อนุมัติ -> ไม่ต้องทำอะไรเพิ่ม
- Q-APV-21 server ไม่ทำงานตอนเที่ยงคืน -> ไล่ปฏิเสธคำขอที่เลยเวลาตอนเริ่มทำงาน

สิ่งที่แก้หลังได้คำตอบ:
- spec.md Draft v2 -> v2.1, plan.md Draft v1 -> v1.1, constraints.md CON-TECH-03 (รหัสใหม่ 5 ตัว + 004 ในรหัสเดิม 4 ตัว), CHANGELOG
- สิ่งที่ AI ทำเอง (ทีมตรวจ): ตั้งเลข Q-APV-14 ถึง 22; เสนอ AC-APV-17 (ค้างคืนถูกปฏิเสธตาม / ไม่ถูกแตะตอนเที่ยงคืน) และ AC-APV-18 (ไล่ปฏิเสธตอนเริ่มทำงาน) เป็น [เสนอ] รอทีมอนุมัติ ID (Q-APV-22); ไม่ทำ CR โดยอ้างแบบเดียวกับตอนเพิ่ม FORBIDDEN_ROLE
- ความเสี่ยงที่ทีมรับทราบ (Q-APV-20): ถ้า 003 ไม่ล็อกแถวเดียวกัน คำขอที่เพิ่งอนุมัติอาจถูกยกเลิกทับได้

---

## 2569-11-12 คำสั่ง: /tasks ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/tasks.md)
- คำสั่ง: "ทำ /tasks"
- ผลลัพธ์: specs/004-approve-request/tasks.md 20 เฟส 39 task (เฟส 0 ฐานราก 3 task, เฟส 1 ถึง 18 เฟสละ 1 AC ครบ 18 AC, เฟส 19 รอคำตอบ 3 task: T-37, T-38 รอ Q-APV-04, T-39 รอ Q-SH-03), ตารางตรวจความครบ 3 ตาราง (AC, Constraint, FR) ไม่มีช่องว่าง
- task ที่ยากที่สุด: T-07 / T-16 / T-35 (ล็อกแถวทุกใบของชิ้นเดียวกันรวมค้างคืน และพิสูจน์ว่ากดพร้อมกันได้ผลเดียว) และ T-29 / T-30 (งานเที่ยงคืนกับการไล่ปฏิเสธตอนเริ่มทำงานไม่ให้ทำซ้ำ)
- AC ที่ทดสอบยาก: AC-APV-10 (พร้อมกันจริง ใช้ goroutine + -count=20), AC-APV-13 / 18 (เที่ยงคืน ใช้นาฬิกาจำลอง), AC-APV-14 (รอ 4 ครั้ง ใช้นาฬิกาจำลอง)
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): เรียงเฟสตามที่เจ้าหน้าที่ใช้งาน (ดูรายการ -> อนุมัติ -> ปฏิเสธ -> Email -> กรณีพิเศษ -> งานเบื้องหลัง -> พร้อมกัน); seed ของ 004 ต่อท้าย backend/db/seed/seed.sql; ลองหน้าจอด้วย VITE_MOCK_TOKEN=mock-staff-001; ใช้ curl ลอง AC-APV-07 / 12 / 05; ตัวส่ง Email จำลองตั้งให้ล้มได้ (ทีมเลือกเอง ไม่ได้มาจาก spec)
- เงื่อนไขก่อนเริ่ม T-01: จองเลข migration และเพิ่มค่า action ใน DM-07 (ไม่ใช่ Open Question ของ spec)

---

## 2569-11-12 คำสั่ง: แจ้งเลข migration ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code
- คำสั่ง: "migration 0008"
- สิ่งที่แก้: tasks.md (หัวไฟล์, T-01 ใช้ `0008_approval.up.sql` / `.down.sql`, สิ่งที่ยังไม่ทำ), plan.md (ข้อ 3, 7.1, 8)
- สิ่งที่ AI ไม่แก้เอง: ตาราง "เลข migration" ใน data-model.md (นอกขอบเขตที่ 004 แก้ได้) -> แจ้งทีมเพิ่มแถว 0008 ของ 004
- ข้อสังเกตของ AI: ใน repo มี migration ถึง 0003 ถ้า 0004 ถึง 0007 ขึ้น main หลัง 0008 golang-migrate จะไม่รันให้ฐานข้อมูลที่อยู่ที่ 0008 แล้ว ต้องตกลงลำดับ merge / `audit_log_action_check` ที่ 0008 ADD ใหม่ต้องรวมค่า action ของ migration ก่อนหน้าทั้งหมด

---

## 2569-11-12 คำสั่ง: ค่า action และเลข migration ไม่ต้องเขียนใน data-model

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code
- คำสั่ง: "ยังไม่ต้องกำหนดใน datamodel 0008 อยูในdb"
- สิ่งที่ AI เข้าใจ (ทีมตรวจ): ยังไม่ต้องเพิ่มแถว action (`approve`, `reject`, `auto_reject`) ใน DM-07 และแถวเลข migration 0008 ใน data-model.md ค่าเหล่านี้กำหนดใน migration 0008 ที่ฐานข้อมูลเท่านั้น
- สิ่งที่แก้: tasks.md (หัวไฟล์, T-01 ไม่ต้องรอ DM-07, สิ่งที่ยังไม่ทำ), plan.md (ข้อ 7.1, 8)
- ผล: T-01 ไม่ติดเงื่อนไขแล้ว เหลือเรื่องลำดับ merge ของ 0004 ถึง 0007 ที่ต้องตกลงในกลุ่ม

---

## 2569-11-12 คำสั่ง: /implement T-01 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่ง: "/implement  T-01"
- ไฟล์ที่สร้าง: `backend/db/migrations/0008_approval.up.sql`, `backend/db/migrations/0008_approval.down.sql`, `backend/internal/platform/migrate_approval_test.go`
- ไฟล์ที่แก้: `backend/internal/platform/migrate_test.go` (ของกลาง แก้เพราะ 004: `TestMigration_AuditLogAction` ของ 001 ใช้ 'approve' เป็นตัวอย่าง action ที่ต้องถูกปฏิเสธ เปลี่ยนเป็น 'delete' เจตนาของ test เดิมไม่เปลี่ยน -> ต้องแจ้งเจ้าของ 001), `specs/004-approve-request/tasks.md` (ช่องไฟล์ที่แตะและสถานะของ T-01)
- ผล test: **ยังไม่ได้รัน** เครื่องไม่มี Go (`go` ไม่อยู่ใน PATH) และ Docker daemon ไม่ได้เปิด (testcontainers ต้องใช้)
- คำถามที่ AI ถามก่อนแตะไฟล์นอกช่อง "ไฟล์ที่แตะ": (1) วาง test ที่ไหน -> ทีม: ไฟล์ใหม่ `migrate_approval_test.go` ใน platform (2) test ของ 001 ที่จะไม่ผ่าน -> ทีม: แก้ 1 บรรทัดแล้วแจ้ง 001
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): `approval_decision` เป็น append-only ด้วย trigger แบบเดียวกับ audit_log (plan ข้อ 3 เขียนว่าไม่มี UPDATE / DELETE); CHECK ให้ result ที่เป็นการปฏิเสธต้องมีเหตุผลไม่ว่าง (FR-APV-01 ง); `approval_email.request_id` UNIQUE และ FK ไป `approval_decision.request_id` (1 ผล = 1 Email, AC-APV-10); `actor_department_id` อ้าง `department(id)`; index `(status, next_attempt_at)` ให้ตัวส่ง Email
- ข้อสังเกต: ค่า action ใน `audit_log_action_check` ของ 0008 รวมค่าถึง 0003 เท่านั้น ถ้า 0004 ถึง 0007 เพิ่มค่า ต้องแก้ 0008 ให้รวมค่าเหล่านั้น / down ของ 0008 จะล้มถ้ามีแถว audit_log ที่ใช้ค่าใหม่แล้ว (append-only)

---

## 2569-11-12 คำสั่ง: เลื่อนการรัน test ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ต้อง implement ก่อนค่อย test"
- สิ่งที่ AI เข้าใจ (ทีมตรวจ): implement task ต่อไปก่อน แล้วรัน test ทีหลังเมื่อเครื่องมี Go และ Docker
- สิ่งที่แก้: tasks.md สถานะ T-01 เป็น "เสร็จ รอทีมตรวจ" พร้อมหมายเหตุว่ายังไม่ได้รัน test
- ความเสี่ยงที่ AI แจ้ง: ข้อผิดพลาดของ migration จะไปเจอตอนรัน test ทีเดียวหลายเฟส และ test ทุกตัวของ 004 ใช้ฐานข้อมูลจริงผ่าน testcontainers

---

## 2569-11-12 คำสั่ง: เปลี่ยนสถานะ T-01 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "เปลี่ยนสถานะเป็นเสร็จ"
- สิ่งที่แก้: tasks.md T-01 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (ทีมตรวจแล้ว) พร้อมหมายเหตุว่า test ยังไม่ได้รัน

---

## 2569-11-12 คำสั่ง: แก้ conflict ตอน merge main เข้า 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request (แก้ไฟล์ร่วม constraints.md)
- คำสั่ง: "อ่านไฟล์ที่พึ่ง pull ลงมา" แล้ว "แก้ conflict ตามนี้แล้ว commit merge"
- สิ่งที่แก้:
  - `specs/000-shared/constraints.md`: รวมแถวของ 004 (รหัส 5 ตัว) และ 005 (REQUEST_NOT_APPROVED, BORROWER_MISMATCH) ในตาราง CON-TECH-03, `FORBIDDEN_ROLE` ใครใช้แล้ว = 004, 005, ตาราง "ใครอ้างใช้บ้าง" ใส่ทั้ง 004 และ 005, IF-NOT-01 คงเป็น 004 (ฝั่ง 005 เขียน `-`), เพิ่มหมายเหตุในบรรทัด Status
  - `specs/CHANGELOG.md`: เก็บทั้งสองฝั่ง ของ 004 (2569-11-12) อยู่บน ของ main (2569-11-11) อยู่ล่าง
  - `prompt-log.md`: เก็บทั้งสองฝั่ง ไม่ลบบันทึกใด ของ main อยู่ก่อนของ 004
- สิ่งที่พบจากไฟล์ที่ merge มา (ยังไม่แก้ รอสั่ง):
  - `TestSeed_BorrowRequestScenarios` ของ 002 บังคับให้ seed.sql ไม่มีคำขอยืม -> T-03 ของ 004 ต้องย้ายไปไฟล์ seed แยกแบบ `handover_dev.sql` ของ 005
  - config.go มี `LoanNormalReturnHour` (DOM-LOAN-01) และ `RequestStatuses` ของ 002 แล้ว -> T-02 ของ 004 ควรใช้ร่วม ไม่ประกาศซ้ำ

---

## 2569-11-12 คำสั่ง: ปรับ T-02 และ T-03 ของ 004 ตามโค้ดที่ merge จาก main

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "แก้เลย"
- สิ่งที่แก้:
  - tasks.md T-03: ไฟล์ที่แตะเปลี่ยนจาก `seed.sql` เป็นไฟล์ใหม่ `backend/db/seed/approve_dev.sql` โหลดเองด้วย psql (แบบ `handover_dev.sql` ของ 005) เพราะ `TestSeed_BorrowRequestScenarios` ของ 002 บังคับว่า seed.sql ไม่มีคำขอยืม; อุปกรณ์ของ 004 ใช้รหัส `AP-TEST-`; เพิ่มคำสั่งโหลดในหัวไฟล์; จุดตรวจทุกเฟสอ้าง approve_dev.sql แทน seed
  - tasks.md T-02, T-21 และ plan.md ข้อ 5: เวลาตัด 16:00 น. ใช้ `LoanNormalReturnHour` ที่ 002 ประกาศไว้ ไม่สร้าง `ApprovalApproveCutoff` ใหม่ (AGENTS.md กฎข้อ 9) และใช้ `RequestStatuses` / `StatusRequested` ของ 002
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ชื่อไฟล์ `approve_dev.sql` และคำนำหน้ารหัสครุภัณฑ์ `AP-TEST-`

---

## 2569-11-12 คำสั่ง: /implement T-02 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่ง: "/implement T-02"
- ไฟล์ที่แก้: `backend/internal/config/config.go` (บล็อกค่าของ 004 คั่นระหว่างค่าของ 002 กับ type Config), `frontend/src/config.js` (`APPROVAL_MESSAGES` วางต่อจากบล็อกของ 005 ต้นไฟล์), `specs/004-approve-request/tasks.md`
- ไฟล์ที่สร้าง: `backend/internal/config/config_approval_test.go` (`TestConfig_ApprovalValues`, `TestConfig_ApprovalMessagesMatchFrontend` ตรวจว่าข้อความใน config.js ตรงกับ config.go)
- ผล test: Go test ยังไม่ได้รัน (เครื่องไม่มี Go ทีมให้ implement ก่อนแล้วค่อย test) / ตรวจ syntax ของ config.js ด้วย `node --check` ผ่าน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ใช้ `LoanNormalReturnHour`, `RequestStatuses` ของ 002 ร่วม; เพิ่ม `StatusApproved`, `StatusRejected`, `ApprovalSystemActorID = "system"`, `ApprovalKind*` (ค่าเดียวกับ CHECK ใน 0008) ซึ่งไม่อยู่ในตาราง plan ข้อ 5 แต่เป็นค่าที่ T-07 ถึง T-29 ต้องใช้; `ApprovalEmailTemplates` เก็บเป็นรายการชื่อข้อมูลตามคำใน spec ส่วนรูปแบบการจัดวางเนื้อ Email ยกไป T-14; ข้อความของ `VALIDATION_FAILED` ใช้ "กรุณาระบุเหตุผลในการปฏิเสธ" เพราะช่องเดียวที่ 004 ตรวจคือเหตุผล
- สิ่งที่เกือบต้องเดา: รูปแบบเนื้อ Email (spec ให้แค่หัวข้อและรายการข้อมูล) -> ไม่เดาใน T-02 เก็บเฉพาะสิ่งที่ spec เขียน

---

## 2569-11-12 คำสั่ง: ลำดับการทดสอบของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ต้อง test หลัง implement task ทั้งหมด"
- สิ่งที่แก้: tasks.md หัวไฟล์ เพิ่ม "ลำดับการทดสอบ" ว่า implement ทุก task ก่อน แล้วค่อยรัน test และลองตามจุดตรวจตอนท้าย (ต่างจากขั้นที่ 2 ของ .claude/commands/implement.md ที่ให้รัน test ทุก task ทีมตัดสินเอง)
- ความเสี่ยงที่ AI แจ้ง: ข้อผิดพลาดจะไปเจอพร้อมกันตอนท้าย และจุดตรวจรายเฟสจะยังไม่ได้ลอง

---

## 2569-11-12 คำสั่ง: /implement T-03 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่ง: "/implement T-03"
- ไฟล์ที่สร้าง: `backend/db/seed/approve_dev.sql` (รุ่นทดลอง 2 รุ่น ภาควิชา 1 และ 2, อุปกรณ์ AP-TEST-0001 ถึง 0005 และ 0101, Loan ของ AP-TEST-0004, คำขอ 8 ใบ: A B C บน 0001 (C ค้างคืน), 0002, 0003, 0004 (อุปกรณ์ถูกยืม), 0005 cancelled, 0101 ภาควิชา 2)
- ไฟล์ที่แก้: `backend/internal/platform/migrate_approval_test.go` (ไฟล์ test ของ 004 เพิ่ม `TestSeed_ApproveDevScenarios` โหลด handover_dev.sql + approve_dev.sql ใน transaction แล้วนับกรณี), `specs/004-approve-request/tasks.md` (ไฟล์ที่แตะและสถานะ T-03)
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน (รันหลัง implement ทุก task)
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ผู้ขอ `student-003` ไม่มีใน MockProvider ใช้เป็นข้อมูลอย่างเดียว (unique index ของ 002 ห้ามผู้ขอคนเดิมมีคำขอค้างบนชิ้นเดิมเกิน 1 ใบ); วันที่เริ่มยืม = วันนี้ถ้าโหลดก่อน 16:00 ไม่อย่างนั้นพรุ่งนี้ (กันชนเวลาตัด FR-APV-01 ค); เวลาที่ส่ง = now() ลบไม่กี่นาที; สร้าง Loan ให้อุปกรณ์ที่ถูกยืม; attachment_path ของคำขอค้างคืนเป็น path ตัวอย่างไม่มีไฟล์จริง
- ข้อสังเกต: วันที่เริ่มยืมเป็นพรุ่งนี้ถ้าโหลดหลัง 16:00 ซึ่งเป็นการยืมปกติล่วงหน้า (ยังรอ 002 ยืนยันว่าทำได้ไหม จาก Q-APV-15)

---

## 2569-11-12 คำสั่ง: /implement T-04 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ทีมตรวจแล้วผ่าน" -> เปลี่ยนสถานะ T-03 เป็น "เสร็จ" (T-04 ต้องรอ T-03)
- คำสั่ง: "/implement T-04"
- คำถามที่ AI ถามระหว่างทำ และคำตอบ: Q-APV-23 คำขอแรก (FR-APV-01 ข) นับคำขอค้างคืนด้วยไหม -> "นับเฉพาะยืมปกติ" (ไม่อย่างนั้นคำขอปกติหลังคำขอค้างคืนจะอนุมัติไม่ได้ตลอด) แก้ spec.md FR-APV-01 (ข) และปิด Q-APV-23
- ไฟล์ที่สร้าง: `backend/internal/approve/model.go` (`PendingRequest`), `backend/internal/approve/repository.go` (`Querier`, `Repository.ListPending` ใช้ row_number() หา is_first ต่ออุปกรณ์รายชิ้น), `backend/internal/approve/repository_test.go` (TestMain/testPool แบบ 002, `devTx` โหลด approve_dev.sql ใน transaction, `TestAC_APV_11`, `TestListPending_OvernightDoesNotBlockQueue`)
- ไฟล์ที่แก้: `backend/internal/config/config.go` (เพิ่ม `LoanTypeNormal`), `specs/004-approve-request/spec.md`, `tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): เรียงซ้ำด้วย id เมื่อเวลาที่ส่งเท่ากัน; รายการว่างคืน [] ไม่ใช่ null; Querier รับได้ทั้ง pool และ transaction เพื่อใช้ต่อใน T-07

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-04 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ทีมตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-04 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-05 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่ง: "/implement T-05"
- คำถามที่ AI ถามระหว่างทำ และคำตอบ: Q-APV-24 ข้อความเมื่อฐานข้อมูลไม่ตอบ (UPSTREAM_UNAVAILABLE) -> "ใช้ข้อความเดียวกับ 002" ("ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง") เพิ่มแถวในตาราง FR-APV-01 (ซ) และปิด Q-APV-24
- ไฟล์ที่สร้าง: `backend/internal/approve/service.go` (`Service`, `staff()`, `ListPending` คืน loadedAt, นาฬิกา inject ได้), `dto.go` (`PendingRequestDTO`, `PendingListDTO`, `ErrorDTO` camelCase แบบ 002), `handler.go` (`PathPrefix = "/v1/approvals/"`, `requireStaff` ห่อทุก endpoint, `GET /v1/approvals/requests`), `handler_test.go` (`TestAC_APV_11_HTTP`, `TestApprovals_RequiresStaff`)
- ไฟล์ที่แก้: `backend/cmd/server/main.go` (ของกลาง แก้เพราะ 004: ผูก route `/v1/approvals/` หลัง auth.Middleware), `config.go` + `config.js` + `config_approval_test.go` (ข้อความ UPSTREAM_UNAVAILABLE), `specs/000-shared/constraints.md` (UPSTREAM_UNAVAILABLE ใครใช้แล้ว เพิ่ม 004), `spec.md`, `tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ชื่อช่อง JSON เป็น camelCase แบบ 002; ตรวจบทบาทครั้งเดียวที่ `requireStaff` ครอบทุก endpoint; กรณีไม่มี principal (เกิดเมื่อผูก route ผิดเท่านั้น) ตอบ UNAUTHENTICATED โดย message ว่าง เพราะ spec ของ 004 ไม่มีข้อความนี้ (ข้อความ 401 ปกติมาจาก auth.Middleware ของกลาง ซึ่งยังไม่มี message ตาม plan ข้อ 8)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-05 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ทีมตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-05 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-06 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่ง: "งั้นทำ /implement T-06 ต่อเลย"
- คำถามที่ AI ถามระหว่างทำ และคำตอบ: tasks.md ให้เพิ่ม route ใน App.jsx แต่ App.jsx ของ 001 ไม่มีระบบ route -> ทีม: "หน้าแยกแบบ 005" (approve.html + approve-main.jsx) แก้ช่องไฟล์ที่แตะของ T-06 และ T-34 และวิธีเปิดหน้าจอในหัวไฟล์ tasks.md
- ไฟล์ที่สร้าง: `frontend/approve.html`, `frontend/src/approve-main.jsx`, `frontend/src/pages/ApproveRequest/ApproveRequests.jsx` (ตาราง ผู้ขอ / อุปกรณ์ / วัตถุประสงค์ / ช่วงเวลาที่ต้องการยืม / เวลาที่ส่ง, ข้อความ error ตามรหัสจาก APPROVAL_MESSAGES), `ApproveRequests.test.jsx` (2 test ขึ้นต้นด้วย AC-APV-11)
- ไฟล์ที่แก้: `frontend/src/api/client.js` (เพิ่ม `getApprovalRequests`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน (และ frontend ยังไม่มี node_modules ต้อง `npm install` ก่อน)
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ไม่มี mockup จึงใช้รูปแบบตารางและข้อความ "กำลังโหลด..." แบบหน้าของ 005; หัวข้อหน้า "คำขอยืมที่รอพิจารณา" และหัวคอลัมน์ทีมออกแบบเอง; วันเวลาแสดงแบบ th-TH เขตเวลา Asia/Bangkok; ไม่มีข้อความเมื่อรายการว่าง (spec ไม่มีข้อความนี้ ตารางแสดงแค่หัวคอลัมน์)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-06 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ทีมตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-06 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-02 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "t-02 ตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-02 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-07 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ทีมตรวจแล้วผ่าน" (T-06) และ "t-02 ตรวจแล้วผ่าน" -> T-02, T-06 เป็น "เสร็จ"
- คำสั่ง: "/implement T-07"
- ไฟล์ที่แก้: `backend/internal/approve/repository.go` (Querier เพิ่ม Exec, Begin; `inTx`, `lockItemRequests` ใช้ `SELECT ... FOR UPDATE OF r` ล็อกคำขอทุกใบของชิ้นนั้น, `setStatus`, `insertDecision`, `insertAudit`), `service.go` (`Approve`, `ErrRequestNotFound`, `ErrNotRequested`), `model.go` (`requestRow`, `decision`, `Result`), `specs/004-approve-request/tasks.md`
- ไฟล์ที่สร้าง: `backend/internal/approve/service_test.go` (`TestAC_APV_01`, `TestApprove_OnlyRequestedNormal`)
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ใส่การตรวจ "คำขอต้องเป็น requested" ใน T-07 เพราะเป็นเงื่อนไขของ FR-APV-01 (ค) ที่ T-07 รองรับ (การแยกข้อความยกเลิก / พิจารณาแล้ว / กดพร้อมกัน ยังเป็นงานของ T-23); คำขอค้างคืนหรือรหัสที่ไม่มีอยู่ถือเป็น "ไม่พบ" (ErrRequestNotFound) เพราะอยู่นอกขอบเขตของ 004; `actor_department_id` ของเจ้าหน้าที่ใช้ภาควิชาของเจ้าหน้าที่ (เท่ากับภาควิชาเจ้าของอุปกรณ์เมื่อ T-25 ตรวจแล้ว); audit_log เขียนเองใน repository ของ 004 เพราะยังไม่มี package audit กลาง
- เรื่องที่ต้องถามตอน T-08: HTTP ของ ErrRequestNotFound ใช้รหัสอะไร CON-TECH-03 ยังไม่มีรหัส "ไม่พบคำขอ"

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-07 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-07 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-08 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจแล้วผ่าน" (T-07 -> เสร็จ), "ขอคำสั่ง commit"
- คำสั่ง: "/implement T-08"
- คำถามที่ AI ถามระหว่างทำ และคำตอบ: Q-APV-25 พิจารณาคำขอที่ไม่มีอยู่หรือเป็นคำขอค้างคืน ตอบรหัสอะไร -> "ลงทะเบียน REQUEST_NOT_FOUND" (404, "ไม่พบคำขอนี้") เพิ่มใน constraints.md CON-TECH-03, spec FR-APV-01 (ซ), plan ข้อ 4, config.go, config.js
- ไฟล์ที่แก้: `backend/internal/approve/handler.go` (`POST /v1/approvals/requests/{id}/approve`, map ErrRequestNotFound / ErrNotRequested), `dto.go` (`DecideInputDTO`, `ResultDTO`), `handler_test.go` (`TestAC_APV_01_HTTP`, `TestApproveHTTP_NotFound`), config 3 ไฟล์, `specs/000-shared/constraints.md`, `spec.md`, `plan.md`, `tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): รับ `loadedAt` แต่ยังไม่ใช้และไม่ปฏิเสธ body ที่อ่านไม่ออก จนถึง T-23; ErrNotRequested ตอบ `REQUEST_ALREADY_DECIDED` ชั่วคราวจนถึง T-23 ที่แยกเป็น 3 รหัส; คำตอบสำเร็จคือ `{requestId, status, decidedAt, rejectedRequestIds}` (รายการว่างจนถึง T-16)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-08 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-08 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-09 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจแล้วผ่าน" (T-08 -> เสร็จ)
- คำสั่ง: "/implement T-09"
- ไฟล์ที่แก้: `frontend/src/api/client.js` (เพิ่ม `postJSON` และ `approveRequest`), `frontend/src/pages/ApproveRequest/ApproveRequests.jsx` (ปุ่ม "อนุมัติ" ทุกแถว, ส่ง loadedAt, เอาคำขอที่อนุมัติและ rejectedRequestIds ออกจากรายการ, ข้อความ error ตามรหัสเหนือรายการ, ปิดปุ่มระหว่างส่ง), `ApproveRequests.test.jsx` (test `AC-APV-01`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ปุ่มอนุมัติแสดงทุกแถวไปก่อน การซ่อนปุ่มแถวที่ไม่ใช่คำขอแรกเป็นงานของ T-18; ข้อความ error ของการอนุมัติแสดงเหนือตารางโดยรายการยังอยู่; ไม่มีหน้าต่างยืนยันก่อนอนุมัติ (spec ไม่ได้กำหนด)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-09 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-09 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-10 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจแล้วผ่าน" (T-09 -> เสร็จ)
- คำสั่ง: "/implement T-10 ต่อได้เลย"
- ไฟล์ที่แก้: `backend/internal/approve/service.go` (`Reject` ใน transaction เดียว: ล็อก, ตรวจ requested + ยืมปกติ, ตั้ง rejected, approval_decision พร้อมเหตุผล, audit_log reject), `handler.go` (`POST /v1/approvals/requests/{id}/reject`), `dto.go` (`DecideInputDTO` เพิ่มช่อง reason), `service_test.go` (`TestAC_APV_02`), `handler_test.go` (`TestAC_APV_02_HTTP`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ใช้ body เดียวกับการอนุมัติ (`DecideInputDTO` เพิ่ม reason); ยังไม่ตรวจเหตุผลว่างที่ service (T-12) ระหว่างนี้ถ้าเหตุผลว่าง CHECK `approval_decision_reason_required` ใน 0008 จะกันไว้และตอบ 503

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-10 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ทีมตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-10 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-11 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ทีมตรวจแล้วผ่าน" (T-10 -> เสร็จ)
- คำสั่ง: "/implement T-11"
- ไฟล์ที่แก้: `frontend/src/api/client.js` (เพิ่ม `rejectRequest`), `frontend/src/pages/ApproveRequest/ApproveRequests.jsx` (ปุ่ม "ปฏิเสธ" เปิดช่อง "เหตุผลในการปฏิเสธ" พร้อมปุ่ม "ยืนยันการปฏิเสธ" และ "ยกเลิก" ในแถวนั้น, รวมการส่งผลอนุมัติ/ปฏิเสธเป็นฟังก์ชัน `decide` เดียว), `ApproveRequests.test.jsx` (test `AC-APV-02`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน / `client.js` ผ่าน `node --check`
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ช่องเหตุผลเปิดในแถวเดียวกัน (ไม่ใช่หน้าต่างลอย); ข้อความ "เหตุผลในการปฏิเสธ" ใช้คำจาก FR-APV-02 (ข), ปุ่ม "ยืนยันการปฏิเสธ" / "ยกเลิก" ทีมออกแบบเอง; ยังไม่ตรวจเหตุผลว่างที่หน้าจอ (T-13)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-11 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ทีมตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-11 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-12 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ทีมตรวจแล้วผ่าน" (T-11 -> เสร็จ)
- คำสั่ง: "/implement T-12"
- ไฟล์ที่แก้: `backend/internal/approve/service.go` (`ErrReasonRequired` ตรวจ strings.TrimSpace ก่อนเปิด transaction), `handler.go` (แปลงเป็น 422 VALIDATION_FAILED พร้อม fields ["reason"]), `dto.go` (`ValidationErrorDTO`), `handler_test.go` (`TestAC_APV_06` 3 กรณี: ว่าง, มีแต่ช่องว่าง, ไม่ส่งช่อง), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- ไฟล์นอกช่อง "ไฟล์ที่แตะ": handler.go และ dto.go เป็นไฟล์ของ 004 เอง ต้องแก้เพื่อแปลง error เป็น HTTP 422 เพิ่มในช่องแล้ว
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ชื่อช่อง `fields` ใน body ของ VALIDATION_FAILED ตามแบบ 002; body ที่อ่านไม่ออกถือว่าไม่ได้ระบุเหตุผล

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-12 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ทีมตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-12 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: รับเรื่องจาก 003 และแก้ conflict ตอน merge main (PR #12 ของ 003) เข้า 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request (แก้ไฟล์ร่วม constraints.md, main.go, client.js เพราะ merge)
- คำสั่ง: ข้อความจากเจ้าของ 003 (ภาพ) "ข้อ 1 ให้เพิ่ม ข้อ 2 บอกเฉยๆ เผื่อ"
  1. migration 0005 ของ 003 เพิ่ม action `notify` -> 0008 ของ 004 ที่เขียน CHECK ใหม่ต้องใส่ `notify` ด้วย
  2. 003 เพิ่ม `postJSON` ใน api/client.js แล้ว ตอน merge ใช้ของที่มีอยู่ ไม่ต้องเพิ่มซ้ำ
- สิ่งที่แก้:
  - `backend/db/migrations/0008_approval.up.sql`: เพิ่ม `notify` ใน `audit_log_action_check` / `0008_approval.down.sql`: คืนเป็นชุดของ 0005 (มี notify)
  - `backend/internal/platform/migrate_approval_test.go`: `TestMigration_ApprovalAuditActions` ตรวจว่า notify ยังใช้ได้
  - แก้ conflict: `main.go` เก็บ route ทั้ง 004 และ 003 / `client.js` ใช้ `postJSON` + `apiError` ของ 003 ตัดของ 004 ออก เก็บฟังก์ชัน API ทั้งของ 004 และ 003 / `constraints.md` ตาราง "ใครอ้างใช้บ้าง" ใส่ทั้ง 003 และ 004, IF-NOT-01 ใส่ทั้ง 003 และ 004 / `CHANGELOG.md`, `prompt-log.md` เก็บทั้งสองฝั่ง
  - `CHANGELOG.md` ของ main มีเศษ marker ค้างจากการ merge ครั้งก่อน (บรรทัด ` 003-cancel-request`, `=======`, ` main`) ลบเฉพาะเศษ marker เนื้อหาครบ
- ข้อสังเกตของ AI (ไม่แก้ ไม่ใช่ของ 004): 003 ใช้รหัส `REQUEST_NOT_FOUND` และ `ALREADY_APPROVED`, `NOT_CANCELLABLE`, `REASON_*` ใน client.js แต่ CON-TECH-03 มีแค่ `REQUEST_NOT_FOUND` (ที่ 004 ลงทะเบียน) -> แจ้ง 003 ให้เพิ่ม 003 ในคอลัมน์ "ใครใช้แล้ว" และลงทะเบียนรหัสที่เหลือ

---

## 2569-11-12 คำสั่ง: /implement T-13 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: ทีม commit merge `eb59817` (main ที่มี 003 PR #12)
- คำสั่ง: "/implement T-13"
- ไฟล์ที่แก้: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx` (`reject` ตรวจเหตุผลว่างหรือมีแต่ช่องว่างก่อนส่ง แสดงข้อความ VALIDATION_FAILED โดยไม่เรียกหลังบ้าน), `ApproveRequests.test.jsx` (2 test `AC-APV-06`: ตรวจที่หน้าจอ และหลังบ้านตอบ VALIDATION_FAILED), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ตรวจเหตุผลว่างที่หน้าจอก่อนส่งด้วย (ไม่ต้องรอหลังบ้าน) แต่หลังบ้านยังตรวจซ้ำเสมอ (T-12); ข้อความแสดงที่จุดเดียวกับ error อื่นเหนือตาราง ช่องเหตุผลยังเปิดอยู่ให้พิมพ์ต่อ

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-13 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบผ่าน"
- สิ่งที่แก้: tasks.md T-13 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-14 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบผ่าน" (T-13 -> เสร็จ)
- คำสั่ง: "/implement T-14"
- คำถามที่ AI ถามระหว่างทำ และคำตอบ: รูปแบบเนื้อ Email (spec กำหนดแค่หัวข้อและรายการข้อมูล) -> "หัวข้อ: ค่า บรรทัดละช่อง" เรียงตาม spec ไม่มีคำขึ้นต้น/ลงท้าย เวลาเป็น พ.ศ. แบบ 003
- ไฟล์ที่สร้าง: `backend/internal/approve/mailer.go` (`Email`, `Mailer`, `LogMailer` เขียนลง log, `RecipientResolver`, `MockRecipientResolver` คืน `<รหัสผู้ใช้>@example.com`), `email.go` (`EmailData`, `BuildEmail` 3 แบบ, formatDateTime / formatDate แบบ พ.ศ.), `email_test.go` (`TestBuildEmail_Contents` เทียบทุกบรรทัด, `TestBuildEmail_PurposeOther`, `TestMockRecipientResolver`)
- ไฟล์ที่แก้: `backend/internal/config/config.go` (เพิ่ม `ApprovalPurposeLabels` ข้อความวัตถุประสงค์ภาษาไทยฝั่งหลังบ้าน), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): "รายการอุปกรณ์ที่ขอยืม" = "ชื่อรุ่น (รหัสครุภัณฑ์)"; "วันที่ยื่นคำขอ" แสดงแค่วันที่ ส่วนวันเวลาอื่นแสดงถึงนาที; "ชื่อผู้อนุมัติ/ผู้พิจารณา" ใช้รหัสเจ้าหน้าที่ไปก่อนเพราะ Principal ยังไม่มีชื่อ (T-38); เขียน formatter เวลาเองแบบเดียวกับ FormatNoticeTime ของ 003 ไม่ import package ของ 003; ที่อยู่ Email จำลองใช้โดเมน example.com ซึ่งสงวนไว้สำหรับตัวอย่าง

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-14 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-14 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-15 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจแล้วผ่าน" (T-14 -> เสร็จ)
- คำสั่ง: "/implement T-15"
- ไฟล์ที่สร้าง: `backend/internal/approve/sender.go` (`Sender` + `Run` ตรวจคิวทุก ApprovalSenderPollInterval, `SendDue` หยิบ Email pending ที่ถึงเวลาด้วย FOR UPDATE SKIP LOCKED แล้วส่งผ่าน Mailer, `markAttempt` บันทึก sent หรือ attempts/last_error/next_attempt_at), `sender_test.go` (`TestAC_APV_04` ด้วยนาฬิกาจำลองและ Mailer จำลอง, `TestSender_PollWithinDeadline`)
- ไฟล์ที่แก้: `service.go` (Approve / Reject ใส่ Email ลงคิวใน transaction เดียวกับผล), `repository.go` (`insertEmail`), `config.go` (`ApprovalEmailPending/Sent/Failed`, `ApprovalSenderPollInterval` = 5 วินาที, `ApprovalSenderBatchSize` = 50), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ช่วงตรวจคิว 5 วินาทีและจำนวนต่อรอบ 50 ฉบับ (ทีมเลือกเอง ไม่ได้มาจาก spec อยู่ใน config); ส่ง Email ภายใน transaction ที่ล็อกแถวไว้ (ง่ายและพอสำหรับ prototype); ส่งไม่สำเร็จตอนนี้ยังเป็น pending ลองใหม่ทุก 1 นาทีไม่จำกัดครั้ง การหยุดหลังครบ 3 ครั้งแล้วเป็น failed เป็นงานของ T-32; ยังไม่เริ่ม goroutine ใน main.go (T-30) ระหว่างนี้ระบบจริงยังไม่ส่ง Email

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-15 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-15 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-16 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจแล้วผ่าน" (T-15 -> เสร็จ)
- คำสั่ง: "/implement T-16"
- ไฟล์ที่แก้: `backend/internal/approve/service.go` (`Approve` เรียก `rejectSuperseded` ใน transaction เดียวกัน: คำขอ requested อื่นของชิ้นเดียวกันทุกประเภทเป็น rejected พร้อม approval_decision เหตุผล ApprovalSupersededReason, audit_log reject, Email ผลปฏิเสธ คืนรหัสใน RejectedBy), `service_test.go` (`TestAC_APV_03` ตรวจ A approved, B และค้างคืน C rejected), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ใช้แถวที่ล็อกไว้แล้วจาก lockItemRequests ไม่ query ซ้ำ; test ตรวจคำขอค้างคืน C ด้วย (ส่วนของ AC-APV-17 ที่ยัง [เสนอ] T-31 จะตรวจแยกอีกครั้ง)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-16 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจแล้วผ่าน"
- สิ่งที่แก้: tasks.md T-16 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-17 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจแล้วผ่าน" (T-16 -> เสร็จ)
- คำสั่ง: "/implement T-17"
- ไฟล์ที่แก้: `backend/internal/approve/service.go` (`ErrNotFirst`, `firstInQueue` หาคำขอยืมปกติ requested ที่ส่งก่อนสุดจากแถวที่ล็อกไว้, `Approve` ตรวจหลังตรวจสถานะ requested), `handler.go` (แปลงเป็น 409 NOT_FIRST_IN_QUEUE), `service_test.go` (`TestAC_APV_07`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- ไฟล์นอกช่อง "ไฟล์ที่แตะ": handler.go ของ 004 เอง ต้องแก้เพื่อแปลง error เป็น 409 เพิ่มในช่องแล้ว
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ลำดับใช้กฎเดียวกับ is_first ในหน้ารายการ (submitted_at แล้ว id, นับเฉพาะยืมปกติ Q-APV-23); บังคับลำดับเฉพาะการอนุมัติ ปฏิเสธคำขอที่ไม่ใช่คำขอแรกได้ (ASM-04 "คำขอแรกไม่ผ่านให้ปฏิเสธ" ไม่ได้ห้ามปฏิเสธคำขออื่น)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-17 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล้วผ่าน ขอคำสั่ง commit ที่ไม่ต้องมีวงเล็บรอtest"
- สิ่งที่แก้: tasks.md T-17 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-18 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล้วผ่าน" (T-17 -> เสร็จ) และขอให้ข้อความ commit ไม่มีวงเล็บรอ test ต่อท้าย
- คำสั่ง: "/implement T-18"
- ไฟล์ที่แก้: `frontend/src/pages/ApproveRequest/ApproveRequests.jsx` (แสดงปุ่ม "อนุมัติ" เฉพาะแถว isFirst แถวอื่นมีแค่ปุ่มปฏิเสธ), `ApproveRequests.test.jsx` (2 test `AC-APV-07`: ปุ่มเฉพาะคำขอแรก และหลังบ้านตอบ NOT_FIRST_IN_QUEUE), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ซ่อนปุ่มอนุมัติไปเลย (ไม่แสดงแบบกดไม่ได้); ข้อความ NOT_FIRST_IN_QUEUE ใช้ทางเดียวกับ error อื่นที่ทำไว้ใน T-09 จึงไม่ต้องแก้โค้ดส่วนแสดงข้อความ

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-18 และ /implement T-19 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่ง: "ตรวจสอบผ่าน /implement T-19"
- T-18: tasks.md "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)
- T-19 ไฟล์ที่แก้: `backend/internal/approve/service.go` (`ErrItemNotAvailable` ตรวจ target.ItemStatus หลังตรวจคำขอแรก), `handler.go` (แปลงเป็น 409 ITEM_NOT_AVAILABLE), `config.go` (`ItemStatusAvailable` ค่าเดียวกับ equipment.StatusAvailable ของ 001 ไม่ import package ของ 001), `service_test.go` (`TestAC_APV_08` ทั้งอุปกรณ์ถูกยืมและไม่พร้อมใช้งาน), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ลำดับการตรวจ: สถานะคำขอ -> คำขอแรก -> อุปกรณ์ว่าง (คำขอที่ไม่ใช่คำขอแรกได้ NOT_FIRST_IN_QUEUE ก่อน ITEM_NOT_AVAILABLE)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-19 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบผ่าน"
- สิ่งที่แก้: tasks.md T-19 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-20 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบผ่าน" (T-19 -> เสร็จ)
- คำสั่ง: "/implement T-20"
- ไฟล์ที่แก้: `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx` (test `AC-APV-08`), `specs/004-approve-request/tasks.md`
- ไม่แก้ `ApproveRequests.jsx` แม้อยู่ในช่องไฟล์ที่แตะ เพราะหน้าจอแสดงข้อความตามรหัส error จาก APPROVAL_MESSAGES อยู่แล้วตั้งแต่ T-09 (decide)
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-20 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบผ่าน"
- สิ่งที่แก้: tasks.md T-20 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-21 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบผ่าน" (T-20 -> เสร็จ)
- คำสั่ง: "/implement T-21"
- ไฟล์ที่แก้: `backend/internal/approve/service.go` (`ErrApprovalCutoff`, `approvalCutoff` = LoanNormalReturnHour ของวันที่เริ่มยืมตามเวลาไทย, ตรวจด้วยนาฬิกาที่ inject ได้หลังตรวจอุปกรณ์ว่าง), `handler.go` (แปลงเป็น 409 APPROVAL_CUTOFF_PASSED), `service_test.go` (`TestAC_APV_16`, `TestApprove_BeforeCutoffAllowed`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- ไฟล์นอกช่อง "ไฟล์ที่แตะ": handler.go ของ 004 เอง เพิ่มในช่องแล้ว
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): "เลย 16:00 น." = หลัง 16:00:00 พอดี (16:00:00 ยังอนุมัติได้); ลำดับการตรวจ: สถานะคำขอ -> คำขอแรก -> อุปกรณ์ว่าง -> เวลาตัด; ID ของ AC-APV-16 ยังรอ Q-APV-13 ชื่อ test ตั้งรอไว้

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-21 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-21 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-22 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-21 -> เสร็จ)
- คำสั่ง: "/implement T-22"
- ไฟล์ที่แก้: `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx` (test `AC-APV-16` ตรวจข้อความ คำขอยังอยู่ และยังมีปุ่มปฏิเสธ), `specs/004-approve-request/tasks.md`
- ไม่แก้ `ApproveRequests.jsx` เพราะแสดงข้อความตามรหัสจาก APPROVAL_MESSAGES อยู่แล้วตั้งแต่ T-09
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน / ID ของ AC-APV-16 ยังรอ Q-APV-13

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-22 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-22 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-23 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-22 -> เสร็จ), "/implement T-22" ซ้ำ (AI ไม่ทำซ้ำ ถามกลับ), "ตรวจสอบแล่วผ่าน" (ไม่มี task รอตรวจ)
- คำสั่ง: "/implement T-23"
- ไฟล์ที่แก้: `backend/internal/approve/service.go` (ลบ ErrNotRequested แทนด้วย `ErrRequestCancelled`, `ErrAlreadyDecided`, `ErrDecidedConcurrently` และ `notRequestedError`; Approve / Reject รับ loadedAt), `repository.go` (`decisionTime`), `handler.go` (ส่ง loadedAt จาก body และแปลงเป็น 409 REQUEST_CANCELLED / REQUEST_ALREADY_DECIDED / REQUEST_DECIDED_CONCURRENTLY), `config.go` (`StatusCancelled`), `service_test.go` + `sender_test.go` (เพิ่มอาร์กิวเมนต์ loadedAt ทุกที่ที่เรียก, test อนุมัติซ้ำคาด ErrAlreadyDecided, เพิ่ม `TestAC_APV_09` 3 กรณีทั้งอนุมัติและปฏิเสธ), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): เทียบ approval_decision.decided_at (เวลาของฐานข้อมูล) กับ loadedAt (เวลาของหลังบ้านตอนโหลดรายการ) ตาม plan ข้อ 4; ไม่มี loadedAt หรือ body อ่านไม่ออก ถือเป็น "พิจารณาแล้ว"; คำขอที่สถานะถูกเปลี่ยนโดยฟีเจอร์อื่นโดยไม่มี approval_decision (เช่น handed_over) ได้ REQUEST_ALREADY_DECIDED

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-23 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-23 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-24 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-23 -> เสร็จ)
- คำสั่ง: "/implement T-24"
- ไฟล์ที่แก้: `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx` (4 test `AC-APV-09`: REQUEST_CANCELLED และ REQUEST_ALREADY_DECIDED ทั้งตอนอนุมัติและปฏิเสธ), `specs/004-approve-request/tasks.md`
- ไม่แก้ `ApproveRequests.jsx` เพราะแสดงข้อความตามรหัสจาก APPROVAL_MESSAGES อยู่แล้วตั้งแต่ T-09
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-24 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-24 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-25 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-24 -> เสร็จ), ทีมวางคำสั่ง commit ในแชต AI รันแล้วถูกหยุดกลางคัน ทีม commit `a9bd482` เอง -> AI จดไว้ว่าต่อไปให้แค่คำสั่ง git ไม่รันเอง
- คำสั่ง: "/implement T-25"
- ไฟล์ที่แก้: `backend/internal/approve/service.go` (`ErrForbiddenDepartment` ตรวจภาควิชาเจ้าของอุปกรณ์กับภาควิชาของเจ้าหน้าที่ทั้งใน Approve และ Reject เป็นข้อแรกหลังพบคำขอ), `handler.go` (แปลงเป็น 403 FORBIDDEN_DEPARTMENT), `handler_test.go` (`TestAC_APV_12` ทั้งอนุมัติและปฏิเสธ ตรวจว่า body มีแค่ code/message), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- ไฟล์นอกช่อง "ไฟล์ที่แตะ": handler.go ของ 004 เอง เพิ่มในช่องแล้ว
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ตรวจภาควิชาก่อนสถานะคำขอและเรื่องอื่นทั้งหมด เพื่อไม่ให้เจ้าหน้าที่ภาควิชาอื่นรู้สถานะของคำขอผ่านรหัส error ที่ต่างกัน; รหัสคำขอที่ไม่มีอยู่ยังได้ REQUEST_NOT_FOUND ก่อน (ตรวจก่อนภาควิชา)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-25 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-25 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-26 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-25 -> เสร็จ)
- คำสั่ง: "/implement T-26"
- ไฟล์ที่แก้: `frontend/src/pages/ApproveRequest/ApproveRequests.test.jsx` (2 test `AC-APV-12` ตอนอนุมัติและตอนปฏิเสธ ใช้ `failWith` ที่เพิ่มไว้ใน T-24), `specs/004-approve-request/tasks.md`
- ไม่แก้ `ApproveRequests.jsx` เพราะแสดงข้อความตามรหัสจาก APPROVAL_MESSAGES อยู่แล้วตั้งแต่ T-09
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-26 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-26 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-27 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-26 -> เสร็จ)
- คำสั่ง: "/implement T-27"
- ไฟล์ที่แก้: `backend/internal/approve/handler_test.go` (`TestAC_APV_05` ไม่มี token -> 401 UNAUTHENTICATED / นักศึกษา -> 403 FORBIDDEN_ROLE พร้อมข้อความ ทั้ง approve และ reject แล้วตรวจว่าสถานะ ผลการพิจารณา และ audit ไม่เปลี่ยน), `specs/004-approve-request/tasks.md`
- ไม่แก้โค้ดหลัก เพราะ requireStaff (T-05) และ auth.Middleware ครอบ endpoint ทั้งสองอยู่แล้ว
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- ข้อสังเกต: 401 ตรวจแค่ code ไม่ตรวจ message เพราะ auth.Middleware ของกลางยังไม่ส่ง message (plan ข้อ 8 แจ้ง 001)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-27 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-27 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-28 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-27 -> เสร็จ)
- คำสั่ง: "/implement T-28"
- ไฟล์ที่แก้: `backend/internal/approve/service_test.go` (`assertHistory` ตรวจ approval_decision + audit_log ของคำขอ 1 ใบ, `TestAC_APV_15` อนุมัติ / ปฏิเสธ / ปฏิเสธตาม), `specs/004-approve-request/tasks.md`
- ไม่แก้โค้ดหลัก เพราะประวัติถูกบันทึกตั้งแต่ T-07, T-10, T-16
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): แยก `assertHistory` เป็นตัวช่วยให้ T-29 ใช้ต่อกับกรณีปฏิเสธอัตโนมัติ

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-28 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-28 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-29 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-28 -> เสร็จ)
- คำสั่ง: "/implement T-29"
- ไฟล์ที่สร้าง: `backend/internal/approve/autoreject.go` (`AutoRejecter`, `autoRejectCutoff` = 00:00 น. ของวันนี้ตามเวลาไทย, `Run` ตื่นทุกเที่ยงคืน, `RejectStale` ปฏิเสธคำขอยืมปกติ requested ที่ส่งก่อน cutoff: rejected + approval_decision auto_rejected ผู้ดำเนินการ system/System + audit_log auto_reject + Email ปฏิเสธอัตโนมัติ, `lockStaleRequests` ใช้ FOR UPDATE แบบรอ), `autoreject_test.go` (`TestAC_APV_13`)
- ไฟล์ที่แก้: `service_test.go` (`TestAC_APV_15` เพิ่มกรณีปฏิเสธอัตโนมัติ), `config.go` (`ApprovalSystemActorName = "System"`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ใช้ FOR UPDATE แบบรอ (ไม่ใช่ SKIP LOCKED) เพราะงานรันวันละครั้ง ถ้าข้ามคำขอที่กำลังถูกพิจารณาจะค้างไปอีกวัน; actor_role ของ audit_log ใช้ "system"; ภาควิชาของผู้ดำเนินการ System = ภาควิชาเจ้าของอุปกรณ์ (plan ข้อ 3); ยังไม่เริ่ม goroutine ใน main.go และยังไม่ไล่ปฏิเสธตอนเริ่มทำงาน (T-30)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-29 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-29 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-30 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-29 -> เสร็จ)
- คำสั่ง: "/implement T-30"
- ไฟล์ที่แก้: `backend/internal/approve/autoreject.go` (`CatchUp` ปฏิเสธคำขอที่เลยเวลาทันที, `Run` เรียก CatchUp ก่อนเข้าวงรอเที่ยงคืน), `autoreject_test.go` (`TestAC_APV_18` server เริ่มใหม่ 09:30 น. หลังเที่ยงคืนที่พลาดไป), `backend/cmd/server/main.go` (ของกลาง แก้เพราะ 004: เริ่ม goroutine `approve.NewSender(db, LogMailer, MockRecipientResolver).Run` และ `approve.NewAutoRejecter(db).Run`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน / ID ของ AC-APV-18 ยังรอ Q-APV-22
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ไล่ปฏิเสธตอนเริ่มทำงานใช้ RejectStale ตัวเดียวกับรอบเที่ยงคืน (เหตุผล ผู้ดำเนินการ Email เหมือนกัน); goroutine ใช้ ctx ของ main (context.Background) หยุดเมื่อโปรแกรมจบ; ใช้ LogMailer และ MockRecipientResolver ในระบบจริงไปก่อน (รอ Q-APV-04, T-37)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-30 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-30 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-31 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-30 -> เสร็จ)
- คำสั่ง: "/implement T-31"
- ไฟล์ที่แก้: `backend/internal/approve/service_test.go` (`TestAC_APV_17` 2 subtest: อนุมัติ A แล้วคำขอค้างคืน C ถูกปฏิเสธตามพร้อมเหตุผลและ Email / เที่ยงคืนแล้ว C ยังเป็น requested), `specs/004-approve-request/tasks.md`
- ไม่แก้ `autoreject_test.go` แม้อยู่ในช่องไฟล์ที่แตะ เพราะรวมทั้งสองกรณีไว้ใน test เดียวตามชื่อ AC; ไม่แก้โค้ดหลัก เพราะพฤติกรรมทำไว้แล้วใน T-16 และ T-29
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน / ID ของ AC-APV-17 ยังรอ Q-APV-22

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-31 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-31 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-32 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-31 -> เสร็จ)
- คำสั่ง: "/implement T-32"
- ไฟล์ที่แก้: `backend/internal/approve/sender.go` (`markAttempt` ส่งไม่สำเร็จครบ 1 + ApprovalEmailRetryCount = 4 ครั้งแล้วเป็น failed ใน UPDATE เดียว), `sender_test.go` (fakeMailer จดเวลาทุกครั้งที่ถูกเรียก, `TestSender_RetryThenFailed`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): นับครั้งใน SQL (attempts + 1 >= 4) ไม่อ่านค่าก่อน; next_attempt_at ยังถูกตั้งตอนเป็น failed แต่ไม่มีผลเพราะตัวส่งหยิบเฉพาะ pending

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-32 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-32 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-33 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-32 -> เสร็จ)
- คำสั่ง: "/implement T-33"
- ไฟล์ที่แก้: `backend/internal/approve/handler.go` (`GET /v1/approvals/failed-emails`, `POST /v1/approvals/emails/{id}/resend`), `dto.go` (`FailedEmailDTO`, `FailedEmailListDTO`, `ResendResultDTO`), `repository.go` (`ListFailedEmails` กรองตามภาควิชาเจ้าของอุปกรณ์, `lockEmail`, `requeueEmail`), `service.go` (`ListFailedEmails`, `Resend`, `ErrEmailNotFound`), `model.go` (`FailedEmail`), `handler_test.go` (`TestAC_APV_14`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน
- ไฟล์นอกช่อง "ไฟล์ที่แตะ": service.go และ model.go ของ 004 เอง เพิ่มในช่องแล้ว
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): กดส่งซ้ำ = ส่ง Email กลับเข้าคิว (pending, นับครั้งใหม่) ให้ตัวส่งเบื้องหลังส่งในรอบถัดไปภายใน 5 วินาที ไม่ได้ส่งทันทีใน request; Email ที่ไม่ได้ failed (ส่งแล้ว/อยู่ในคิว) กดซ้ำไม่ทำอะไร คืนสถานะเดิม จึงไม่เกิด Email ซ้ำและไม่ต้องลงทะเบียนรหัส error ใหม่; Email ที่ไม่มีอยู่ตอบ REQUEST_NOT_FOUND (รหัสที่ลงทะเบียนแล้ว); รายการกรองตามภาควิชาเจ้าของอุปกรณ์ (รวมของ System)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-33 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-33 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: /implement T-34 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- เครื่องมือ: Claude Code (ทำตาม .claude/commands/implement.md เพราะ /implement ไม่ได้ติดตั้งใน session)
- คำสั่งก่อนหน้า: "ตรวจสอบแล่วผ่าน" (T-33 -> เสร็จ)
- คำสั่ง: "/implement T-34"
- ไฟล์ที่สร้าง: `frontend/src/pages/ApproveRequest/FailedEmails.jsx` (ตาราง ผู้ขอ / อุปกรณ์ / Email / สาเหตุ / เวลา พร้อมปุ่ม "ส่งซ้ำ" ส่งแล้วแถวหาย error แสดงตามรหัส), `FailedEmails.test.jsx` (2 test `AC-APV-14`)
- ไฟล์ที่แก้: `frontend/src/api/client.js` (`getFailedEmails`, `resendEmail`), `frontend/src/approve-main.jsx` (แสดง FailedEmails ต่อจาก ApproveRequests ในหน้า /approve.html), `frontend/src/config.js` (`APPROVAL_EMAIL_KIND_LABELS`), `specs/004-approve-request/tasks.md`
- ผล test: ยังไม่ได้รัน ตามลำดับการทดสอบที่ทีมตัดสิน / `client.js`, `config.js` ผ่าน `node --check`
- สิ่งที่ AI เลือกเอง (ทีมตรวจ): ไม่มี mockup ทีมออกแบบเอง; หัวข้อ "Email ที่ส่งไม่สำเร็จ" ใช้คำจาก FR-APV-02 (ง), หัวคอลัมน์และข้อความแบบ Email (ผลอนุมัติ / ผลปฏิเสธ / ปฏิเสธอัตโนมัติ) ทีมออกแบบเอง; รายการว่างแสดงแค่หัวตาราง (spec ไม่มีข้อความ); ไม่โหลดรายการใหม่เองหลังมี Email ล้มเพิ่ม (ต้องรีเฟรชหน้า)

---

## 2569-11-12 คำสั่ง: ทีมตรวจ T-34 ของ 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request
- คำสั่ง: "ตรวจสอบแล่วผ่าน"
- สิ่งที่แก้: tasks.md T-34 "เสร็จ รอทีมตรวจ" -> "เสร็จ" (test ยังไม่ได้รัน)

---

## 2569-11-12 คำสั่ง: อ่านไฟล์ที่ pull และแก้ conflict ตอน merge main (PR #16 ของ 001) เข้า 004-approve-request

- ผู้สั่ง: 660710758 นางสาวณัฐสุดา กัญญาบุญ (เจ้าของฟีเจอร์ 004)
- ฟีเจอร์: 004-approve-request (แก้ไฟล์ร่วม constraints.md เพราะ merge)
- คำสั่ง: "อ่านที่ pull ลงมาให้หน่อย" แล้ว "แก้ตามนี้"
- คำสั่งก่อนหน้า (ยังรอทีม): 002 แจ้งว่า `ALREADY_APPROVED` ของ 003 ซ้ำกับ `REQUEST_ALREADY_DECIDED` ของ 004 AI เสนอ 3 ทาง (003 ใช้รหัสของ 004 ร่วม / แยกกันโดย 003 ลงทะเบียนเอง / 004 เปลี่ยนตาม ไม่แนะนำ) ทีมยังไม่ตัดสิน
- สิ่งที่เข้ามาจาก 001: `platform/apierror.go` (WriteError {code, message}), `auth.Middleware` ตอบ message แล้ว, `audit/writer.go`, RequestQueue (FR-EQS-11, CR-005), `EquipmentDBTimeout` ใน config.go
- สิ่งที่แก้:
  - `specs/000-shared/constraints.md`: `FORBIDDEN_ROLE` ใครใช้แล้ว = 001, 004, 005 / ตาราง "ใครอ้างใช้บ้าง" CON-TECH-01, 02, 03, CON-DATA-01 ใส่ครบ 001 ถึง 005 และเก็บหมายเหตุ "FORBIDDEN_ROLE ใน T-17" ของ 001
  - `specs/CHANGELOG.md`: เก็บทั้งสองฝั่ง ของ 004 อยู่บน ลบเศษ marker ` 003-cancel-request` ที่ค้างใน main
  - `prompt-log.md`: เก็บทั้งสองฝั่ง ของ main ก่อนแล้วตามด้วยของ 004
  - `specs/004-approve-request/plan.md` ข้อ 8: ปิดเรื่อง auth.Middleware ไม่มี message (001 แก้แล้ว)
- ข้อสังเกตของ AI: FORBIDDEN_ROLE ข้อความของ 001 กับ 004 ต่างกัน (หน้าจอเลือกจากรหัส 004 ใช้ข้อความตาม FR-APV-01 ซ) ไม่ต้องแก้; TestAC_APV_05 (T-27) ตรวจข้อความของ 401 เพิ่มได้แล้ว ยังไม่แก้ รอสั่ง
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

---

## 2569-11-12 /implement T-13 (002-borrow-request) — ปิดเฟส 8

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 8 AC-BRQ-12 กฎปลดล็อก 00:00 ของวันถัดไป)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-13`

### คำถามที่ AI ถามก่อนเริ่ม และคำตอบของทีม
T-13 มีแค่ `service.go` ในช่องไฟล์ที่แตะ แต่กฎข้อ ข ต้องรู้ `returned_at` ของใบที่คืนแล้ว
ส่วน query เดียวที่มีคือ `EarliestOpenLoanDue` จาก T-12 ซึ่งกรอง `returned_at IS NULL` ทิ้งไป
ใบที่คืนแล้วจึงไม่มีทางถึงมือ service ทำด้วยไฟล์เดียวไม่ได้จริง
→ **ทีมตอบ: ทางเลือก ก** เพิ่ม `repository.go` ในช่องไฟล์ที่แตะของ T-13 พร้อมเหตุผลลงวันที่
(เป็นครั้งที่ 4 ที่ช่องไฟล์ที่แตะไม่พอ ต่อจาก T-09 T-11 T-12 ทีมน่าทบทวนวิธีเขียนช่องนี้ตอน retrospective)

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/service.go` (ตรงกับช่องไฟล์ที่แตะ)
  ขยาย `hasOverdueLoan` ให้ครอบข้อ ข, เพิ่ม `returnedLateToday` และ `startOfDay`
- `backend/internal/borrowrequest/repository.go` (เพิ่มในช่องไฟล์ที่แตะวันนี้ ตามคำตอบของทีม)
  เพิ่ม `ReturnedLoan`, `returnedSinceSQL`, `ReturnedSince` และเพิ่ม `Query` ใน interface `Querier`
- `backend/internal/borrowrequest/service_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** ต่อ test ท้ายไฟล์เดิม

### การตัดสินใจ
- กฎ "00:00 น. ของวันถัดไป" เขียนเป็น "การคืนช้านั้นเกิดขึ้นในวันเดียวกับวันนี้หรือไม่"
  โดย service ส่งเวลา 00:00 ของวันนี้ (`startOfDay`) ลงไปเป็นเงื่อนไขของ query
  กฎจึงอยู่ใน Go ส่วน SQL เป็นแค่ตัวกรองช่วงเวลา และได้ผลพลอยได้คือ query คืนแถวไม่กี่แถวเสมอ
- "คืนช้า" ใช้ `overdueDeadline` ตัวเดียวกับข้อ ก (16:00 ของกำหนดคืน) ไม่มีนิยามที่สอง
- ใช้ `config.OverdueBlockUntilNextDayStart` จริงในโค้ด ไม่ใช่ค่าที่ตั้งไว้เฉย ๆ
  ถ้าทีมเปลี่ยนเป็น false กฎข้อ ข จะปิดทันที (ตรงกับตาราง "ค่าที่ตั้งได้" ใน plan ข้อ 5)
- เพิ่ม `Query` ใน `Querier` เพราะกฎนี้ต้องอ่านหลายแถว ของเดิมมีแต่ `QueryRow`
  `*pgxpool.Pool` รองรับทั้งสองอยู่แล้ว ไม่กระทบ main.go และ test

### สิ่งที่เกือบต้องเดาแต่ไม่ได้เดา
- **"คืนอุปกรณ์แล้ว" หมายถึงคืนใบไหน** spec ข้อ ข เขียนลอย ๆ ว่า "เมื่อคืนอุปกรณ์แล้ว"
  แต่ AC-BRQ-12 เขียนชัดว่า "คืนอุปกรณ์**ที่เกินกำหนด**แล้ว" และ FR-BRQ-06 ทั้งข้อพูดถึงผู้ที่คืนเกินกำหนด
  จึงอ่านว่าใบที่คืนตรงเวลาไม่เข้ากฎนี้ ไม่ได้เดาเอง และเขียน test ยืนยันไว้
  (`TestOverdue_OnTimeReturnDoesNotBlock`) ถ้าทีมอ่านต่างจากนี้ให้บอก

### ผล test
- `go test ./internal/borrowrequest/ -run "TestAC_BRQ_12|TestOverdue" -v`: PASS ทั้ง 4 ตัว
  (`TestAC_BRQ_12`, `TestOverdue_BoundaryAtReturnHour`, `TestOverdue_OtherDepartmentNotCounted`,
  `TestOverdue_OnTimeReturnDoesNotBlock`)
- `go test ./... -count=1`: ทุกแพ็กเกจ ok
- `go vet ./internal/borrowrequest/`: ไม่มีข้อความเตือน

### จุดตรวจเฟส 8 ครึ่งหลัง (ทีมลองเองด้วยมือ 2569-11-12)

- ผู้ลอง: ธัญลักษณ์ ด. (@tnyl02) จำลองการคืนด้วย psql แล้วยิง curl ด้วย `mock-student-002`
- ตั้ง `returned_at` = 14:00 **ของวันนี้** (คืนช้าวันนี้): **409** `{"code":"HAS_OVERDUE_LOAN"}`
- ตั้ง `returned_at` = 14:00 **ของเมื่อวาน** (ข้ามวันแล้ว): **201** requestId 10
- ยืนยันกฎ FR-BRQ-06 ข้อ ข บนของจริง ตรงกับ `TestAC_BRQ_12` ที่ใช้นาฬิกาตรึงเวลา 23:59 กับ 00:01
- หมายเหตุ: ทดสอบด้วยการขยับ `returned_at` แทนการรอข้ามวันจริง ซึ่งเป็นวิธีเดียวที่ทำได้ในเวลาจำกัด
  ส่วน test อัตโนมัติไม่ต้องใช้วิธีนี้เพราะตั้งนาฬิกาเองได้
- เฟส 8 ผ่านครบทั้ง AC-BRQ-06 และ AC-BRQ-12

---

## 2569-11-12 /implement T-14 (002-borrow-request) — เฟส 9 ครึ่งแรก

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 9 AC-BRQ-10 ไฟล์แนบใบบันทึกข้อความของการยืมค้างคืน)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-14`

### คำถามที่ AI ถามก่อนเริ่ม และคำตอบของทีม
ช่อง "เสร็จเมื่อ" บังคับให้ "บันทึก path ในตาราง" แต่ช่องไฟล์ที่แตะมีแค่ attachment.go, handler.go, validate.go
การเขียนคอลัมน์ `attachment_path` ต้องผ่าน repository และต้องมีทางเรียกจาก service
→ **ทีมตอบ: ทางเลือก ก** เพิ่ม `service.go` กับ `repository.go` ในช่องไฟล์ที่แตะ (ครั้งที่ 5 ต่อจาก T-09 T-11 T-12 T-13)

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/attachment.go` **ไฟล์ใหม่** — `Upload`, `parseCreateForm`, `AttachmentStore`, `safeFileName`
- `backend/internal/borrowrequest/validate.go` — `attachmentProblems`, `mediaType`, ช่อง `fieldAttachment`
- `backend/internal/borrowrequest/handler.go` — รับ multipart, จำกัดขนาด request, รวมกฎไฟล์แนบเข้ากับรายการช่องที่ผิด
- `backend/internal/borrowrequest/service.go` (เพิ่มในช่องวันนี้) — `CreateWithAttachment` และฟิลด์ `attachments`
- `backend/internal/borrowrequest/repository.go` (เพิ่มในช่องวันนี้) — `SetAttachmentPath`, `Delete`, เพิ่ม `Exec` ใน `Querier`
- `backend/internal/borrowrequest/handler_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** — test ใหม่ และแก้ test เดิม 1 ตัว (ดูด้านล่าง)

### การตัดสินใจ
- ลำดับการทำงานเป็น สร้างคำขอ -> ได้ id -> เขียนไฟล์ `<id>-<ชื่อที่ปลอดภัย>` -> อัปเดต `attachment_path`
  เพราะ plan.md ข้อ 3 กำหนดให้ชื่อไฟล์มีรหัสคำขออยู่ด้วย จึงต้องรู้ id ก่อน
  ถ้าเขียนไฟล์หรืออัปเดตไม่สำเร็จ จะลบคำขอใบนั้นทิ้ง ไม่ทิ้งคำขอค้างที่ไม่มีไฟล์แนบ (FR-BRQ-09)
- ตรวจชนิดไฟล์จาก Content-Type ที่ผู้ส่งแจ้งมา เทียบกับ `config.AttachmentAllowedTypes`
  **ยังไม่ได้ตรวจเนื้อไฟล์จริง (content sniffing)** spec ไม่ได้สั่ง จึงไม่เพิ่มเอง แต่ทีมควรรู้ว่าเป็นช่องโหว่
- `safeFileName` ตัดชื่อไฟล์จากผู้ขอให้เหลือเฉพาะอักษรปลอดภัย ชื่ออย่าง `../../etc/passwd` จึงพาไฟล์ออกนอกโฟลเดอร์ไม่ได้
- การยืมปกติไม่บังคับให้แนบไฟล์ (FR-BRQ-10 ข้อ ก พูดถึงเฉพาะค้างคืน) แต่ถ้าแนบมาก็ตรวจชนิดและขนาดด้วยกฎเดียวกัน
- โฟลเดอร์ไฟล์แนบอ่านจาก `config.Load().AttachmentDir` ที่เดียว test ใช้ `t.Setenv("ATTACHMENT_DIR", t.TempDir())`

### ข้อผิดพลาดระหว่างทำ และการแก้
1. `TestSafeFileName` จับได้ว่าชื่อไฟล์ภาษาไทยถูกแปลงเป็น `_` ทุกตัวแล้วถูกตัดจนเหลือ `pdf`
   คือ **นามสกุลไฟล์หาย** แก้ที่โค้ด (ยุบตัวคั่นซ้ำ และไม่ตัดจุดนำหน้านามสกุลทิ้ง) ไม่ได้แก้ test ให้ผ่าน
2. `TestLoanWindow_NormalBoundaryAndOvernight` (ของ T-08) ไม่ผ่านหลังเพิ่มกฎ
   เพราะมันยิงยืมค้างคืนโดยไม่แนบไฟล์ ซึ่งกฎใหม่ปฏิเสธถูกต้องแล้ว
   แก้ที่ test ให้ส่งแบบ multipart พร้อมไฟล์ PDF เจตนาเดิมของ test (กฎเวลาคืน) ยังอยู่ครบ

### สิ่งที่ยังไม่ได้ทำใน task นี้
- AC-BRQ-10 ท่อนที่ว่า "ยืมปกติต้องแสดงเงื่อนไขคืนภายใน 16:00 และใช้ในพื้นที่มหาวิทยาลัย" เป็นงานหน้าจอ รอ Q-BRQ-14 (T-19, T-20)
- endpoint เปิดดูไฟล์แนบเป็นของ T-15

### ผล test
- `go test ./internal/borrowrequest/ -run "TestAC_BRQ_10|TestSafeFileName"`: PASS
  (`TestAC_BRQ_10` มี 4 กรณีย่อย: ไม่แนบไฟล์ / ชนิดที่ไม่รับ / ใหญ่เกิน 10MB / แนบ PDF ถูกต้อง)
- `go test ./... -count=1`: ทุกแพ็กเกจ ok
- `go vet ./internal/borrowrequest/`: ไม่มีข้อความเตือน

### จุดตรวจเฟส 9 ส่วนของ T-14 (ทีมลองเองด้วยมือ 2569-11-12)

- ผู้ลอง: ธัญลักษณ์ ด. (@tnyl02) ด้วย `curl -F` ผ่าน multipart จริง
- ยืมค้างคืนไม่แนบไฟล์: **422** `{"code":"VALIDATION_FAILED","fields":["attachment"]}`
- ยิงใหม่พร้อมแนบ `/tmp/memo.pdf` (79 ไบต์, `type=application/pdf`): **201** requestId 11
- `ls backend/var/attachments` เห็น `11-memo.pdf` ขนาด 79 ไบต์ เท่าไฟล์ต้นฉบับ
  และในตารางมี `attachment_path` = `var/attachments/11-memo.pdf` ชี้ไฟล์เดียวกัน
  ชื่อไฟล์ขึ้นต้นด้วยรหัสคำขอตาม plan.md ข้อ 3
- ข้อสังเกต: สิทธิ์ไฟล์ขึ้นเป็น `-rwxrwxrwx` ทั้งที่โค้ดเขียนด้วย 0600
  เพราะโฟลเดอร์อยู่บนไดรฟ์ Windows ที่ WSL mount มา ซึ่งไม่เก็บ mode bits ของ Linux
  บนเซิร์ฟเวอร์จริงหรือใน container จะได้ 0600 ตามที่เขียน ไม่ใช่ข้อบกพร่องของโค้ด
- ล้างข้อมูลที่ลองไว้: ลบคำขอของ student-001 และลบโฟลเดอร์ `backend/var/attachments`

---

## 2569-11-12 /implement T-15 (002-borrow-request) — ปิดเฟส 9

- ผู้สั่ง: ธัญลักษณ์ ด. (@tnyl02)
- ฟีเจอร์: 002-borrow-request (เฟส 9 endpoint เปิดดูไฟล์แนบ จำกัดสิทธิ์)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: `/implement T-15`

### คำถามที่ AI ถามก่อนเริ่ม และคำตอบของทีม (3 ข้อ)
1. ขอบเขตไฟล์: ต้องอ่าน `requester_id` กับ `attachment_path` จากฐานข้อมูล แต่ช่องไฟล์ที่แตะมีแค่ handler.go กับ main.go
   → **ทีมตอบ ก** เพิ่ม `service.go` กับ `repository.go` (ครั้งที่ 6 ต่อจาก T-09 T-11 T-12 T-13 T-14)
2. รหัส 403 สำหรับ "ไม่ใช่เจ้าของคำขอ" ซึ่งตารางกลางไม่มี
   → **ทีมตอบ ก** ใช้ `FORBIDDEN_ROLE` ที่ 004 ลงทะเบียนไว้แล้ว ไม่ตั้งรหัสใหม่
3. คำขอที่ไม่มีอยู่จริงหรือไม่มีไฟล์แนบ ซึ่งตารางกลางไม่มีรหัส 404 เลย
   → **ทีมตอบ ก** ลงทะเบียน `REQUEST_NOT_FOUND` (404) เป็นรหัสกลาง ใช้ร่วมกับ 003 และ 004

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/borrowrequest/handler.go` (ตรงกับช่อง) — route `GET /v1/loan-requests/{id}/attachment` และการแปลง error
- `backend/cmd/server/main.go` (ตรงกับช่อง) — ผูก `PathPrefix+"/"` เพิ่ม
- `backend/internal/borrowrequest/service.go` (เพิ่มในช่องวันนี้) — `OpenAttachment`, `ErrAttachmentNotFound`, `ErrNotRequestOwner`
- `backend/internal/borrowrequest/repository.go` (เพิ่มในช่องวันนี้) — `AttachmentOf`, `RequestAttachment`
- `backend/internal/borrowrequest/handler_test.go` **ไม่อยู่ในช่องไฟล์ที่แตะ** — test ใหม่
- **ไฟล์ของกลาง (แก้เพราะ T-15 ของฟีเจอร์ 002 ตามที่ AGENTS.md สั่งให้บันทึกเหตุผล)**
  - `specs/000-shared/constraints.md` v5.1 -> v5.2 เพิ่มแถว `REQUEST_NOT_FOUND` และเติม 002 ในแถว `FORBIDDEN_ROLE`
  - `specs/changes/CR-007.md` ไฟล์ใหม่ (ข้าม CR-006 ไว้ให้ 003 ซึ่งกำลังเสนอ CR เรื่องรหัส error รอบเดียวกัน)
  - `specs/CHANGELOG.md` บันทึก shared constraints v5.2

### การตัดสินใจ
- `main.go` ต้องผูก `PathPrefix+"/"` เพิ่มอีกบรรทัด เพราะ mux ชั้นนอกจับ path แบบเป๊ะ
  ถ้าไม่ผูก path ลูกอย่าง `/v1/loan-requests/11/attachment` จะไปไม่ถึง router ของฟีเจอร์ (เจอตอนออกแบบ ไม่ได้เจอตอน test ล้ม)
- `AttachmentOf` รวม "ไม่มีคำขอ" กับ "ไม่มีไฟล์แนบ" เป็นคำตอบเดียว ผู้เรียกแยกไม่ได้ (DOM-PRIV-01)
- เจ้าหน้าที่เปิดได้ทุกใบ **ไม่ได้จำกัดว่าต้องเป็นภาควิชาเดียวกัน** เพราะ FR-BRQ-10 (ก) เขียนแค่ "เฉพาะเจ้าหน้าที่กับเจ้าของคำขอ"
  ไม่ได้เขียนเรื่องภาควิชา และ DOM-ACCESS-01 พูดถึงนักศึกษา ถ้าทีมต้องการจำกัดต้องเพิ่มใน spec ก่อน **ไม่ได้เดาเอง**
- ใช้ `http.ServeFile` ส่งไฟล์ ชนิดไฟล์ให้เดาจากนามสกุลซึ่ง `safeFileName` รักษาไว้ตั้งแต่ T-14

### ข้อจำกัดที่บันทึกไว้ใน CR-007
คำขอที่มีอยู่จริงและมีไฟล์แนบแต่เป็นของคนอื่น จะได้ 403 ไม่ใช่ 404
นักศึกษาจึงยังเดาได้ว่าเลขคำขอนั้นมีไฟล์แนบอยู่ แม้ไม่เห็นตัวไฟล์
เลือกแบบนี้เพราะช่อง "เสร็จเมื่อ" ของ T-15 กำหนดให้นักศึกษาคนอื่นได้ 403

### สภาพแวดล้อม
รอบแรกรัน test ไม่ได้เพราะ Docker Desktop ถูกปิดอยู่ (`docker` หายไปจาก WSL)
ทีมเปิดให้แล้วจึงรันต่อได้ ไม่ได้แก้โค้ดใด ๆ ระหว่างนั้น

### ผล test
- `go test ./internal/borrowrequest/ -run TestAC_BRQ_10 -v`: PASS ทั้งหมด
  `TestAC_BRQ_10_AttachmentAccess` มี 5 กรณีย่อย: เจ้าหน้าที่ 200 / เจ้าของคำขอ 200 /
  นักศึกษาคนอื่น 403 `FORBIDDEN_ROLE` และต้องไม่มีเนื้อไฟล์ติดไปใน body / ไม่มีคำขอ 404 `REQUEST_NOT_FOUND` / ไม่มี token 401
- `go test ./... -count=1`: ทุกแพ็กเกจ ok
- `go vet ./...`: ไม่มีข้อความเตือน
