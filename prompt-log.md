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

## 2569-11-10 อ่านของกลางและ SRS ก่อนเริ่ม 006 แล้วตอบคำถามที่ยังไม่ชัด

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: ให้อ่าน AGENTS.md, specs/000-shared/ (3 ไฟล์) และ SRS หัวข้อ 3.0 กับ 3.1.6 แล้ว AI สรุปเรื่องที่ยังไม่ชัด 7 ข้อ

### คำถามที่ AI ถาม และคำตอบของเจ้าของฟีเจอร์
1. ถ้ารับคืนแล้วอุปกรณ์ชำรุด สถานะเป็นอะไร (ลำดับ UC-09 กับ UC-07) -> เป็น "ไม่พร้อมใช้งาน"
2. UC-07 ตรวจสอบอุปกรณ์ที่คืน อยู่ในขอบเขตของ 006 ไหม -> อยู่ในขอบเขต
3. ต้องแนบรูปสภาพอุปกรณ์ตอนรับคืนไหม (RQ-04) -> ทำ
4. สภาพอุปกรณ์ตอนรับคืนมีค่าอะไรบ้าง -> ปกติ, ชำรุด (สูญหายไม่ใช้เป็นสภาพ)
5. คืนเลยกำหนดต้องทำอะไรไหม -> บันทึกหรือทำเครื่องหมายไว้
6. UC-04 คืนอุปกรณ์ (นักศึกษา) ที่ไม่มี RQ รองรับ -> รวมกับ UC-09
7. กำหนดคืนของยืมค้างคืนยังไม่รู้ (Q-SH-04) กระทบ 006 ไหม -> "ไม่มีทางไม่รู้กำหนด" (รอ AI ยืนยันความหมาย)

### ยังไม่ได้แก้ไฟล์ใด คำตอบข้อ 1, 2, 6 ขัดกับ SRS / README / data-model.md ต้องทำ CR ก่อน

---

## 2569-11-10 เขียน FR ของ 006 ใน SRS หัวข้อ 3.1.6

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: "เขียน FR ของตัวเองใน SRS (ตอนนี้ยังเป็น [เติม])"
- สิ่งที่ AI แก้: `docs/srs/SRS.md` หัวข้อ 3.1.6 เพิ่ม FR-RET-01 ถึง FR-RET-07 จาก RQ-04, RQ-11, RQ-13, RQ-25, US-010
  และคำตอบของเจ้าของฟีเจอร์ในบันทึกก่อนหน้า เพิ่ม RQ-04, RQ-25 ในบรรทัดที่มา และหมายเหตุ "รอ CR"
- ยังค้าง `[ยืนยัน]` 3 จุด: กรณีสูญหายใครจัดการ (FR-RET-02), รูปบังคับหรือไม่/กี่รูป (FR-RET-05), ใครเห็นเครื่องหมายคืนเกินกำหนด (FR-RET-06)
- ไม่ได้แก้: หัวข้อ 3.1.7, 4.3 และ data-model.md (ต้องทำ CR ก่อน)

---

## 2569-11-10 สร้าง specs/006-return/spec.md (Draft v1)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code
- สิ่งที่สั่ง: สร้าง spec.md ครบ 6 ส่วน + Assumptions & Open Questions, Constraint ใช้ร่วมอ้าง ID ไม่ copy ข้อความ,
  ใส่ชื่อในช่อง Owner, เพิ่ม 006-return ในตาราง "ใครอ้างใช้บ้าง" ของ constraints.md
- สิ่งที่ AI สร้าง/แก้:
  - `specs/006-return/spec.md` FR-RET-01 ถึง 07 ตาม SRS 3.1.6, AC-RET-01 ถึง 08 (AI ร่าง), ASM-01, 02, Q-01 ถึง Q-08,
    Clarifications 2569-11-10, Traceability (มีหัวข้อ Clarifications และ Traceability เพิ่มตามแบบของ 001)
  - `specs/000-shared/constraints.md` ตาราง "ใครอ้างใช้บ้าง": เพิ่ม 006-return ใน CON-TECH-01, CON-TECH-02, CON-DATA-01,
    CON-DATA-03, DOM-AUDIT-01, DOM-LOAN-01, IF-IDP-01 (แก้ของกลางเพราะฟีเจอร์ 006 ตามกติกาท้ายไฟล์)
- สิ่งที่ AI เลือกเอง (ทีมต้องตรวจ): ไม่อ้าง DOM-PRIV-01, DOM-ACCESS-01, DOM-REQ-01 ในใช้ร่วม (ถามเป็น Q-03, Q-06, Q-08 แทน)
  และเลือก NFR-DATA-01, NFR-SEC-01, NFR-AVAIL-01 จาก SRS 3.3

---

## 2569-11-10 คำสั่ง: /clarify (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code (`.claude/commands/clarify.md`)
- ไฟล์: specs/006-return/spec.md (Draft v1)
- คำสั่ง: "/clarify -> /plan -> /tasks -> /implement T-xx ทีละ task แล้ว commit ทุกขั้น" (เริ่มที่ /clarify ซึ่งต้องหยุดรอคำตอบ)

### คำถามที่ AI ถาม (ทั้งหมด)
ส่วน A (จุดกำกวม 12 แถว): 1 เวลาคืนมาจากไหน, 2 วิธีเลือกการยืม, 3 FR-RET-02 มี 2 พฤติกรรม, 4 รูปบังคับ/กี่รูป/ชนิด/ขนาด,
5 ใครเห็นคืนเกินกำหนด, 6 FR-RET-04/05 ขัด SRS 3.1.7 และ data-model, 7 IF-IDP-01 ล่ม, 8 FR-RET-07 กับ DOM-AUDIT-01 (อ่าน vs เขียน),
9 ASM-01 การยืมค้างคืนไม่มีกำหนดคืน, 10 Goal ไม่มีตัวเลข, 11 คำว่า "สภาพ" ไม่มีใน glossary, 12 คืนมาชำรุดแต่มีคำขอค้าง (DOM-REQ-01)
ส่วน B: Q1 เวลาคืน, Q2 วิธีหาการยืม, Q3 เงื่อนไขรูป, Q4 ใครเห็นคืนเกินกำหนด, Q5 สูญหาย, Q6 ภาควิชาอื่น, Q7 ความหมาย "ไม่มีทางไม่รู้กำหนด", Q8 ทำ CR ก่อน plan ไหม
ถามซ้ำ 4 ข้อเพราะคำตอบขัดกัน: รูป (A4 กับ Q3), เวลาคืนแก้ได้ไหม (A1 กับ Q1), สูญหายเป็นสถานะ (ขัด CON-DATA-01), ชำรุดเป็น "ส่งซ่อม" + ระงับคำขอ (ขัด CON-DATA-01 และเป็น FR ใหม่)

### คำตอบของทีมและเหตุผล
- A1 ระบบบันทึก / A2 ดูรายการที่ยังไม่คืน / A3 แยก test ทีละเรื่อง / A4 แนบทุกครั้ง 1 รูป JPG, PNG
- A5 เจ้าหน้าที่เห็นในรายการยืม/คืน หรือหน้ารายงานประวัติรายผู้ยืม / A6 ปล่อยไว้ก่อน
- A7 (ตอบแถว DOM-AUDIT-01) ให้ DOM-AUDIT-01 ครอบทั้งอ่านและเขียน ไม่ให้ FR-RET-07 กำหนดเอง (แถว IF-IDP-01 ไม่ได้ตอบ คงเป็น Q-SH-08)
- A8 การยืมทุกรายการมีกำหนดคืนเสมอ / A9 Goal ไม่เป็นไร / A10 เพิ่ม "สภาพ (ตอนรับคืน)" ใน glossary
- A11 ชำรุดให้เป็น "ส่งซ่อม/ไม่พร้อมให้ยืม" ระงับคำขอค้าง และแจ้งผู้ขอ
- Q1 ใช้เวลาคืนเป็นค่าเริ่มต้น / Q2 กรอกรหัสครุภัณฑ์เป็นหลัก และมีหน้ารายการที่ยังไม่คืน
- Q3 ปกติไม่บังคับ ชำรุดบังคับ >= 1 รูป สูงสุด 5 รูป JPG/PNG <= 5 MB / Q4 เจ้าหน้าที่เห็นในหน้ารายละเอียดการยืมและรายการที่ยังไม่คืน ผู้ยืมเห็นเฉพาะของตนเอง
- Q5 ปิดการยืมแบบ "สูญหาย" ต้องมีหมายเหตุ สถานะเป็น "สูญหาย" / Q6 เฉพาะภาควิชาตนเอง / Q7 ใช่ / Q8 ทำ CR ก่อน /plan
- คำตอบที่ถามซ้ำ: รูป = บังคับทุกครั้ง + ขอบเขตจาก Q3 / เวลาคืนแก้ไม่ได้ / สูญหาย = "ไม่พร้อมใช้งาน" + เหตุสูญหาย /
  ชำรุด = "ไม่พร้อมใช้งาน" ส่วนระงับคำขอและแจ้งผู้ขอเสนอใน CR

### สิ่งที่แก้ใน spec.md (v1 เป็น v2)
- FR-RET-01 เลือกได้ 2 ทาง (รหัสครุภัณฑ์ / หน้ารายการที่ยังไม่คืน) เฉพาะภาควิชาตนเอง
- FR-RET-02 เวลาคืน = เวลาระบบ แก้ไม่ได้ / FR-RET-05 บังคับ 1-5 รูป JPG/PNG <= 5 MB
- FR-RET-06 ระบุที่แสดงของเจ้าหน้าที่ / FR-RET-07 อ้าง DOM-AUDIT-01 ไม่กำหนดฟิลด์เอง
- AC: แก้ AC-RET-01, 04, 05, 07, 08 และเพิ่ม AC-RET-09 ถึง 13 (AI ร่าง ทีมต้องตรวจ)
- Constraints ใช้ร่วม: เพิ่ม DOM-REQ-01 / Scope: เพิ่มหน้ารายการที่ยังไม่คืน และย้ายสูญหาย/ระงับคำขอไป Out of scope (รอ CR-003)
- ASM-01 ถึง 07, Q-01 ถึง Q-09 (ปิด Q-02, Q-05, Q-06 เพิ่ม Q-09), Clarifications, Traceability
- ของกลางที่แก้เพราะ 006: glossary.md เพิ่ม "สภาพ (ตอนรับคืน)" (v2.1), constraints.md ตารางใครอ้างใช้ เพิ่ม 006 ที่ DOM-REQ-01

---

## 2569-11-10 คำสั่ง: /plan (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code (`.claude/commands/plan.md`)
- คำสั่งของทีม: "ไม่ต้องทำ CR-003.md ทำ plan ไปเลย" (ขัดกับ Q-07 ใน spec v2 ที่เขียนว่าทำ CR ก่อน plan ทีมตัดสินใหม่)
- ก่อนเริ่ม: commit spec v2 (`0a7e66c 006: spec v2`)
- ผลลัพธ์: specs/006-return/plan.md (Draft v1)
- Constraint ที่ AI ยังไม่ได้ใช้: DOM-REQ-01 (การระงับคำขอค้างเป็น FR ใหม่ รอ Q-08)
- สิ่งที่ AI บอกว่าอยากเดาแต่ไม่ได้เดา: ช่องหมายเหตุความเสียหาย (Q-04), หน้าของผู้ยืม (Q-03), ปิดการยืมแบบสูญหาย (Q-01),
  ฟิลด์ audit สถานะก่อน/หลัง (Q-09), พฤติกรรมเมื่อ auth ล่ม (Q-SH-08)
- สิ่งที่ AI เลือกเอง (ทีมเลือกเอง ไม่ได้มาจาก spec): เก็บรูปเป็น BYTEA, 5 MB = 5 × 1024 × 1024, ไม่แบ่งหน้ารายการที่ยังไม่คืน,
  path `/api/v1/returns/...`, error code ใหม่ (FORBIDDEN_ROLE, NO_OPEN_LOAN, CONDITION_REQUIRED, PHOTO_*), action audit ใหม่ 4 ค่า
- ของกลางที่ plan จะแก้ตอน implement: คอลัมน์ใหม่ในตาราง `loan` (DM-05), ค่า action ใหม่ของ `audit_log` (DM-07), เลข migration ต้องจองในกลุ่ม

---

## 2569-11-10 คำสั่ง: /tasks (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code (`.claude/commands/tasks.md`)
- ก่อนเริ่ม: ทีมแจ้งว่า "ตรวจ plan เรียบร้อยแล้ว" (ไม่ได้ตอบข้อ 1-4 ที่ AI ถามหลัง /plan: ข้อความ "รอ CR-003" ใน spec, Q-P1 ถึง Q-P5, เลข migration, คอลัมน์ใหม่ใน loan)
- ผลลัพธ์: specs/006-return/tasks.md รอบแรก 15 เฟส (ฐานราก + 13 เฟสตาม AC + รอคำตอบ) 30 task
- รอ Q-xx 4 task: T-27 (Q-04), T-28 (Q-03), T-29 (Q-09), T-30 (Q-SH-08) และ T-01 ต้องจองเลข migration ก่อน (Q-P4)
- เฟสแรกที่ลองใช้งานได้จริง: เฟส 1 (บันทึกคืนสภาพปกติ อุปกรณ์กลับเป็น "ว่าง")
- ไม่มี task: Q-01 (สูญหาย) และ Q-08 (ระงับคำขอค้าง) เพราะไม่มี FR รองรับ / DOM-REQ-01 ไม่มี task ด้วยเหตุผลเดียวกับ plan
- test เพิ่มที่ไม่ได้มาจาก AC ตรงๆ: `TestReturns_Lookup_FoundOpenLoan`, `TestReturns_ConcurrentReturn_OneWins`

---

## 2569-11-10 คำสั่ง: /implement T-01 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 0 ฐานราก)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "/implement T-xx ทีละ task แล้ว commit 006-return ทุกขั้น" -> `/implement T-01 specs/006-return/tasks.md`
- ก่อนหน้า: ทีมตอบ "จองเองเลย" เรื่องเลข migration (ใช้แค่ใน branch นี้) -> AI ใช้ `0006` แก้ data-model.md ตารางเลข migration, plan Q-P4, tasks T-01
  (รอบแรกที่สั่ง implement ทีมหยุดและให้ถอน commit `89e9cb4` และบันทึกของรอบนั้นออก แล้วสั่งใหม่)

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- สร้าง `backend/db/migrations/0006_return.up.sql`: `loan.return_condition` (normal/damaged), `loan.returned_late`,
  CHECK `loan_return_fields_complete`, ตาราง `loan_return_photo` (JPG/PNG, BYTEA), action ใหม่ 4 ค่าใน `audit_log_action_check`
- แก้ `specs/000-shared/data-model.md` (ของกลาง แก้เพราะ 006): DM-05 ฟิลด์ตอนรับคืน, DM-07 แถว action ของ 006, ตารางเลข migration แถว 0006
- แก้ `backend/internal/platform/migrate_test.go`: เพิ่ม 3 test ต่อท้าย ไม่แก้ test เดิมของ 001

### ผล test (`cd backend && go test ./... -v`)
- TestMigration006_AuditLogReturnActions PASS
- TestMigration006_LoanReturnColumns PASS (ปกติ/ชำรุดได้, lost / ไม่ครบ 3 ค่า ถูกปฏิเสธ)
- TestMigration006_LoanReturnPhoto PASS (JPG/PNG ได้, GIF/PDF/ว่าง และ loan ที่ไม่มีอยู่ ถูกปฏิเสธ)
- test เดิมของ 001 (auth, healthz, migration, seed, audit_log) ยัง PASS ทั้งหมด

### สิ่งที่เกือบต้องเดา
- ชื่อ CHECK เดิมของ audit_log: ใช้ชื่อที่ PostgreSQL ตั้งให้อัตโนมัติ `audit_log_action_check` (ยืนยันด้วย test ที่ผ่าน)
- ไม่ได้ใส่เพดานขนาดรูปเป็น CHECK ใน DB เพราะ AGENTS.md กฎข้อ 9 ให้ตัวเลขอยู่ใน config.go (T-02) ที่เดียว

---

## 2569-11-10 คำสั่ง: /implement T-02 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 0 ฐานราก)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "เอาเลย" -> `/implement T-02 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- แก้ `backend/internal/config/config.go`: `Timezone`, `ReturnPhotoMin`, `ReturnPhotoMax`, `ReturnPhotoMaxBytes`,
  `ReturnRequestMaxBytes`, `ReturnPhotoTypes`, `ReturnConditionStatus` พร้อมคอมเมนต์อ้าง FR-RET-02 ถึง 05, ASM-05, ASM-06, glossary
- สร้าง `backend/internal/config/config_test.go` (ไม่อยู่ในช่องไฟล์ที่แตะ เพิ่มเพราะ AGENTS.md ให้ test อยู่ข้างโค้ด)
- ไม่ได้ทำ `frontend/src/config.js`: AI สร้างแล้วรัน vitest ไม่ได้ (ไม่มี node_modules) ถามว่าจะรัน `npm ci` ไหม
  ทีมตอบ "ไม่ต้องทำอะไรกับ frontend" -> AI ลบ `frontend/src/config.js` และ `frontend/src/__tests__/config.test.js` ที่สร้างไว้ออก

### ผล test (`cd backend && go test ./internal/config -v`)
- TestReturnPhotoLimits PASS / TestReturnConditionStatus PASS / TestTimezoneLoads PASS

### สิ่งที่เกือบต้องเดา
- ที่เผื่อ 1 MB ใน `ReturnRequestMaxBytes` และ 5 MB = 5 × 1024 × 1024 มาจาก plan.md (ทีมเลือกเอง) ไม่ได้ตั้งใหม่
- ยังไม่รู้ว่า "ไม่ต้องทำ frontend" ใช้กับ task หน้าจอที่เหลือ (T-06, T-10, T-13, T-15, T-18, T-20, T-22, T-24) ด้วยหรือไม่ ต้องถามทีม

---

## 2569-11-10 คำสั่ง: /implement T-03 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 0 ฐานราก)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "เอาเลยย" -> `/implement T-03 specs/006-return/tasks.md`
- คำถามค้างจาก T-02 ที่ทีมยังไม่ตอบ: "ไม่ต้องทำ frontend" ใช้กับ task หน้าจอที่เหลือด้วยไหม (T-03 เป็นหลังบ้าน ไม่กระทบ)

### ไฟล์ที่สร้างหรือแก้
- สร้าง `backend/internal/returns/handler.go` (5 route ตาม plan ข้อ 4 ผ่าน `auth.Middleware` ตอบ 501 จนกว่าจะถึง task ของแต่ละ endpoint, `Cache-Control: no-store`)
- สร้าง `backend/internal/returns/repository.go` (โครง Repository)
- แก้ `backend/internal/platform/http.go` (ของกลาง แก้เพราะ 006): เพิ่ม `type Routes` และ `NewRouter(db, routes...)` แบบ variadic
  ให้ฟีเจอร์ผูก route เองโดย platform ไม่ต้อง import ฟีเจอร์ / test เดิมของ 001 ที่เรียก `NewRouter(db)` ใช้ได้เหมือนเดิม
- แก้ `backend/cmd/server/main.go`: สร้าง `auth.NewMockProvider()` และผูก route ของ 006 (ครั้งแรกที่ middleware ถูกผูกเข้า router จริง)
- สร้าง `backend/db/seed/seed_006_return.sql`: `XX-NB-0002` ภาควิชา 2 ถูกยืม + การยืมที่ยังไม่คืน, การยืมที่คืนแล้ว 2 รายการของ CS-NB-0001/0002
- สร้าง `backend/internal/returns/handler_test.go` (ไม่อยู่ในช่องไฟล์ที่แตะ เพิ่มเพราะ AGENTS.md ให้ test อยู่ข้างโค้ด)

### ผล test (`cd backend && go test ./... -v`)
- TestReturns_Skeleton_RoutedThroughAuth PASS (staff -> 501, ไม่มี token -> 401 ทั้ง 5 endpoint)
- TestReturns_Seed006Loaded PASS (ภาควิชา 1 ยังไม่คืน 2 / คืนแล้ว 2, ภาควิชา 2 ยังไม่คืน 1)
- test เดิมทั้งหมด (auth, config, platform) ยัง PASS

### สิ่งที่เกือบต้องเดา
- วิธีผูก route: เลือกส่งฟังก์ชันเข้า `NewRouter` แทนให้ platform import returns เพื่อไม่ให้ของกลางผูกกับฟีเจอร์ใดฟีเจอร์หนึ่ง (002-005 ใช้แบบเดียวกันได้)
- ยังไม่ตรวจบทบาทเจ้าหน้าที่ (นักศึกษาก็ได้ 501) เพราะเป็นงานของ T-25 ตาม tasks.md

---

## 2569-11-10 คำสั่ง: /implement T-04 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 1)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "เอาเลย" -> `/implement T-04 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `backend/internal/returns/repository.go`: `RecordReturn` 1 transaction (`SELECT ... FOR UPDATE OF l, i` -> ปิดการยืมด้วย `now()` ->
  เปลี่ยนสถานะครั้งเดียว -> เก็บรูป -> commit), `ErrNoOpenLoan`
- `backend/internal/returns/service.go` (ใหม่): สถานะหลังคืนจาก `config.ReturnConditionStatus`, `ErrUnknownCondition`
- `backend/internal/returns/handler.go`: `POST /api/v1/returns` multipart (`assetCode`, `condition`, `photos`) ชนิดรูปจาก `http.DetectContentType`
- `backend/internal/returns/handler_test.go`: ตัวช่วย `insertBorrowedItem`, `postReturn` และ AC test / แก้ `TestReturns_Seed006Loaded` ให้นับเฉพาะข้อมูลตัวอย่าง

### ผล test (`cd backend && go test ./... -v`)
- TestAC_RET_01_NormalReturnBecomesAvailable PASS (200, returned_at มีค่า, สถานะ available, รูป 1 รูป, Cache-Control no-store)
- test เดิมทั้งหมดยัง PASS

### สิ่งที่ต้องทำเกินขอบ task เล็กน้อย (บอกทีม)
- `returned_late = (now() > due_at)` ต้องใส่ใน T-04 เพราะ CHECK `loan_return_fields_complete` ของ T-01 บังคับ (test ของค่านี้ยังเป็น T-16)
- ตอบ 422 `CONDITION_REQUIRED` เมื่อสภาพไม่รู้จัก และ 409 `NO_OPEN_LOAN` แบบขั้นต่ำ เพราะต้องมีทางออกของ error เหล่านี้
  (test และข้อความแจ้งอยู่ใน T-09 และ T-21)
- ยังไม่ทำ: ตรวจรูป (T-11, T-12) ชนิดผิดตอนนี้ DB ปฏิเสธแล้วตอบ 503, ตรวจภาควิชา (T-23), บทบาท (T-25), audit (T-26)

---

## 2569-11-11 แก้คำสั่งเรื่อง frontend และทำส่วนหน้าจอของ T-02 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code
- คำสั่ง: "แก้ไข มีหน้า frontend ด้วย" (เปลี่ยนจากคำตอบใน T-02 ที่ว่า "ไม่ต้องทำอะไรกับ frontend")
  -> ตอบคำถามค้างหลัง T-02 ถึง T-04 ว่า task หน้าจอ (T-06, T-10, T-13, T-15, T-18, T-20, T-22, T-24) ต้องทำ
- สิ่งที่ AI ทำ: สร้าง `frontend/src/config.js` (ค่าเดียวกับ config.go: RETURN_PHOTO_MIN/MAX/MAX_BYTES/TYPES, RETURN_CONDITIONS)
  และ `frontend/src/__tests__/config.test.js` ใหม่ (เป็นไฟล์เดิมที่ลบไปใน T-02), รัน `npm ci` ใน frontend/ (ติดตั้งตาม package-lock.json)
- ผล test (`cd frontend && npm test`): config.test.js 2 tests PASS, setup.test.jsx 1 test PASS
- หมายเหตุ: `npm ci` เตือนว่า esbuild มี install script ที่ยังไม่ได้อนุญาต (allowScripts) แต่ vitest รันได้ปกติ

---

## 2569-11-11 /tasks แก้รอบที่ 1 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code
- คำสั่ง: "ช่วยแยกให้หน่อยได้มั้ย task ไหนต้องทำ ควรทำ หรือไม่ต้องทำก็ได้" แล้ว "แยกเข้าไปใน task.md เลยได้มั้ย"
- แก้รอบที่ 1: เพิ่มหัวข้อ "ลำดับความสำคัญ" ใต้ส่วนหัว และบรรทัด `- ความสำคัญ:` ในทุก task (30 task) ไม่ได้แก้ task เดิมช่องอื่น ไม่ผิดกฎข้อใด เป็นการเพิ่มตามที่ทีมขอ
  - ต้องทำ 13 (T-01 ถึง T-06, T-08, T-09, T-11, T-21, T-23, T-25, T-26) / ควรทำ 13 / ไม่ต้องทำตอนนี้ 4 (T-27 ถึง T-30 รอ Q)
- เกณฑ์ที่ AI ใช้ (ทีมต้องตรวจ): "ต้องทำ" = ไม่ทำแล้วระบบผิดหรือไม่ปลอดภัย, "ควรทำ" = มี AC รองรับแต่ระบบหลักใช้ได้, "ไม่ต้องทำตอนนี้" = ติด Q
- AI เตือน: ข้าม task ไหน RTM จะแสดงว่า AC นั้นยังไม่ผ่าน / ลำดับที่แนะนำดึง T-25 ขึ้นก่อน แต่ช่อง "ต้องทำหลัง" ของ T-25 ยังชี้ T-17, T-14 ต้องให้ทีมปรับก่อน

---

## 2569-11-11 คำสั่ง: /implement T-05 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 1)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำ task ต่อไปได้เลย" -> `/implement T-05 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `backend/internal/returns/repository.go`: `OpenLoan`, `FindOpenLoan` (isOverdue = now() > due_at ที่ฐานข้อมูล)
- `backend/internal/returns/handler.go`: `GET /api/v1/returns/lookup?assetCode=` -> 200 / 409 `NO_OPEN_LOAN` / 503
- `backend/internal/returns/handler_test.go`: `TestReturns_Lookup_FoundOpenLoan`

### ผล test (`cd backend && go test ./... -v`)
- TestReturns_Lookup_FoundOpenLoan PASS (พบการยืม ข้อมูลครบ isOverdue ถูกทั้งยังไม่เลยและเลยกำหนด, คืนแล้วได้ 409)
- test เดิมทั้งหมดยัง PASS

### สิ่งที่เกือบต้องเดา
- ไม่ส่ง `assetCode` มาเลย: ตอบเหมือนหาไม่พบ (409 `NO_OPEN_LOAN`) ไม่สร้าง error code ใหม่ เพราะ plan ข้อ 4 ไม่มี code สำหรับกรณีนี้
- ยังไม่กรองภาควิชา (T-23) และไม่ลง audit `return_lookup` (T-26)

---

## 2569-11-11 คำสั่ง: /implement T-06 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 1)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำเลย" -> `/implement T-06 specs/006-return/tasks.md`

### คำถามที่ AI ถาม และคำตอบของทีม
1. หน้าจอได้ token จากไหน (หน้าล็อกอินจำลองเป็นงานของ 001 T-21 ที่ยังไม่ทำ) -> อ่านจาก localStorage key `token`
2. `vite.config.js` ตัด `/api` ออกก่อนส่งหลังบ้าน แต่หลังบ้านใช้ `/api/v1/...` -> เอา rewrite ออก

### ไฟล์ที่สร้างหรือแก้
- ตามช่องไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`,
  `frontend/src/api/client.js` (`api.returns.lookup`, `api.returns.record`, `authHeaders` จาก localStorage, `request` คืน `{ok,status,data}`),
  `frontend/src/App.jsx` (แสดง ReturnPage ใต้หัวเรื่องเดิม)
- นอกช่องไฟล์ที่แตะ (ทีมอนุญาต/ตามกติกา): `frontend/vite.config.js` (เอา rewrite ออก ตามคำตอบข้อ 2 เป็นไฟล์กลาง กระทบ 001),
  `frontend/src/config.js` (เพิ่ม `TIMEZONE`, `STATUS_LABELS` ตามกฎข้อ 9)
- `specs/006-return/tasks.md`: สถานะ T-06 และวิธีตั้ง token ในจุดตรวจเฟส 1

### ผล test (`cd frontend && npm test`)
- ReturnPage.test.jsx 2 tests PASS (AC-RET-01 บนหน้าจอ / ตัวเลือกสภาพมีแค่ปกติ-ชำรุด และไม่มีช่องเวลา)
- config.test.js, setup.test.jsx ยัง PASS / `vite build` สำเร็จ

### สิ่งที่เกือบต้องเดา
- ข้อความ error ตอนนี้เป็นแบบรวม "ค้นหาไม่สำเร็จ (code)" / "บันทึกไม่สำเร็จ (code)" ข้อความแยกตามกรณีอยู่ใน T-10, T-13, T-22, T-24
- ยังไม่มี router หน้าเดียวแสดงใน App.jsx (plan ไม่ได้กำหนด router)

---

## 2569-11-11 คำสั่ง: /implement T-07 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 2)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำเลยย" (ไม่ระบุ task) -> AI เลือก T-07 ตามลำดับเฟส เพราะ T-25 ที่แนะนำยังติดช่อง "ต้องทำหลัง" (T-17, T-14) ที่ทีมยังไม่สั่งแก้

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/returns/handler_test.go`: `TestAC_RET_09_ReturnTimeFromServer`
- `backend/internal/returns/handler.go`: ไม่ต้องแก้ (T-04 ไม่อ่านฟิลด์เวลาจาก client อยู่แล้ว)

### ผล test (`cd backend && go test ./internal/returns -v`)
- TestAC_RET_09_ReturnTimeFromServer PASS (ส่ง `returnedAt` / `returned_at` ปลอม -> returned_at อยู่ระหว่าง now() ก่อนและหลังบันทึก)
- test เดิมของ returns ยัง PASS

### สิ่งที่เกือบต้องเดา
- AC ใช้เวลา 14:30 ตายตัว แต่เวลาคืนมาจาก now() ของฐานข้อมูล จึงตรวจแบบช่วงเวลา (ก่อน/หลังบันทึก) แทน ตามที่ plan.md ข้อ 6 เสนอ

---

## 2569-11-11 คำสั่ง: /implement T-08 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 3)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำตามลำดับ" -> `/implement T-08 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/returns/handler_test.go`: `TestAC_RET_02_DamagedReturnBecomesUnavailable`
  (สร้าง trigger + ตาราง `test_status_history` เฉพาะใน test แล้วลบทิ้งใน Cleanup)
- `backend/internal/returns/service.go`: ไม่ต้องแก้ (T-04 ใช้ `config.ReturnConditionStatus` ซึ่งมี damaged -> unavailable แล้ว)

### ผล test (`cd backend && go test ./internal/returns -v`)
- TestAC_RET_02_DamagedReturnBecomesUnavailable PASS (200, itemStatus unavailable, returned_at มีค่า, ประวัติสถานะ = [unavailable] ครั้งเดียว)
- test เดิมของ returns ยัง PASS

### สิ่งที่เกือบต้องเดา
- ไม่มี: วิธีตรวจ "ไม่เคยเป็นว่างระหว่างทาง" ด้วย trigger มาจาก plan.md ข้อ 6

---

## 2569-11-11 คำสั่ง: /implement T-09 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 4)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำ task ต่อไปเลย" -> `/implement T-09 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/returns/handler_test.go`: `TestAC_RET_04_ConditionRequired`
- `backend/internal/returns/handler.go`: ไม่ต้องแก้ (T-04 ตอบ 422 `CONDITION_REQUIRED` ก่อนแตะฐานข้อมูลอยู่แล้ว ผ่าน `ErrUnknownCondition` ใน service)

### ผล test (`cd backend && go test ./internal/returns -v`)
- TestAC_RET_04_ConditionRequired PASS (ไม่ส่ง / ว่าง / lost / broken / "ปกติ" ภาษาไทย -> 422 CONDITION_REQUIRED, การยืม สถานะ รูป ไม่เปลี่ยน)
- test เดิมของ returns ยัง PASS

### สิ่งที่เกือบต้องเดา
- ไม่มี

---

## 2569-11-11 คำสั่ง: /implement T-10 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 4)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำเลยย" -> `/implement T-10 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `frontend/src/pages/Return/ReturnPage.jsx`: `conditionRequiredMessage` (สร้างจากชื่อสภาพใน config.js), ตรวจก่อนส่งถ้ายังไม่เลือกสภาพ,
  แสดงข้อความเดียวกันเมื่อหลังบ้านตอบ `CONDITION_REQUIRED`
- `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`: ตัวช่วย `openLoanForm` และ test AC-RET-04 2 ตัว

### ผล test (`cd frontend && npm test`)
- AC-RET-04: กดบันทึกโดยไม่เลือกสภาพ -> "กรุณาเลือกสภาพ ปกติ หรือ ชำรุด", ไม่เรียก record, ไม่มี "สูญหาย" PASS
- AC-RET-04: หลังบ้านตอบ CONDITION_REQUIRED -> ข้อความเดียวกัน PASS
- test เดิมทั้งหมด PASS (7 tests)

### สิ่งที่เกือบต้องเดา
- ข้อความ "กรุณาเลือกสภาพ ปกติ หรือ ชำรุด" ไม่มี mockup ทีมออกแบบเอง ใช้ถ้อยคำตาม Then ของ AC-RET-04

---

## 2569-11-11 คำสั่ง: /implement T-11 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 5)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำ task ต่อเลย" -> `/implement T-11 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `backend/internal/returns/repository.go`: `LoanDetail`, `GetLoan`, `ErrLoanNotFound`
- `backend/internal/returns/handler.go`: `GET /api/v1/returns/loans/{loanId}` -> 200 / 404 `LOAN_NOT_FOUND` / 503
- `backend/internal/returns/handler_test.go`: `TestReturns_LoanDetail`

### ผล test (`cd backend && go test ./internal/returns -v`)
- TestReturns_LoanDetail PASS (คืนแล้ว: ข้อมูลครบ condition damaged returnedLate false / ยังไม่คืน: returnedAt null / ไม่มีหรือไม่ใช่ตัวเลข: 404)
- test เดิมของ returns ยัง PASS (8 tests)

### สิ่งที่เกือบต้องเดา
- error code `LOAN_NOT_FOUND` (404) ไม่มีในตาราง error ของ plan ข้อ 4 AI ตั้งเอง ทีมต้องยืนยันและเพิ่มใน plan

---

## 2569-11-11 คำสั่ง: /implement T-12 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 5)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำเลยย" -> `/implement T-12 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- ตามช่องไฟล์ที่แตะ: สร้าง `frontend/src/pages/Return/LoanDetailPage.jsx` (เวลาคืน สภาพ "ยังไม่คืน" ข้อความ error ตาม code),
  `frontend/src/pages/Return/__tests__/LoanDetailPage.test.jsx`, แก้ `frontend/src/api/client.js` (`api.returns.loanDetail`)
- นอกช่องไฟล์ที่แตะ (จำเป็นตามช่อง "เสร็จเมื่อ" ที่ให้ ReturnPage มีลิงก์): `ReturnPage.jsx` ปุ่ม "ดูรายละเอียดการยืม" หลังบันทึกสำเร็จ,
  `ReturnPage.test.jsx` test ของปุ่มนี้

### ผล test (`cd frontend && npm test`)
- LoanDetailPage.test.jsx 3 tests PASS (คืนแล้ว / ยังไม่คืน / โหลดไม่สำเร็จ)
- ReturnPage.test.jsx เพิ่ม 1 test PASS (บันทึกแล้วกดดูรายละเอียด) / รวม 10 tests PASS / `vite build` สำเร็จ

### สิ่งที่เกือบต้องเดา
- แอปยังไม่มี router (plan ไม่ได้กำหนด) "ลิงก์" จึงเป็นปุ่มที่เปิด LoanDetailPage ในหน้าเดียวกัน ไม่เพิ่ม library ใหม่
- ป้าย "คืนเกินกำหนด" ยังไม่แสดง เป็นงานของ T-17

---

## 2569-11-11 คำสั่ง: /implement T-13 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 6)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำต่อเลย" -> `/implement T-13 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/returns/handler_test.go`: `TestAC_RET_06_LateReturnFlagged`
- `backend/internal/returns/repository.go`: ไม่ต้องแก้ (T-04 บันทึก `returned_late = (now() > due_at)` ไว้แล้ว เพราะ CHECK ของ T-01 บังคับ)

### ผล test (`cd backend && go test ./internal/returns -v`)
- TestAC_RET_06_LateReturnFlagged PASS (กำหนดคืนผ่านไปแล้ว -> true / ยังไม่ถึง -> false ทั้งใน response และฐานข้อมูล,
  แก้ due_at ภายหลังแล้ว returned_late ไม่เปลี่ยน)
- test เดิมของ returns ยัง PASS (9 tests)

### สิ่งที่เกือบต้องเดา
- AC ใช้เวลา 17:00 / 15:00 ตายตัว แต่เวลาคืนมาจาก now() ของฐานข้อมูล จึงใช้กำหนดคืนสัมพัทธ์ +/- 1 ชั่วโมงแทน ตาม plan.md ข้อ 6

---

## 2569-11-11 คำสั่ง: /implement T-14 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 7)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "เอาเลย" -> `/implement T-14 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `backend/internal/returns/repository.go`: `ListOpenLoans(departmentID)` กรองภาควิชาที่ชั้น repository เรียงตาม due_at คืน `[]` เมื่อไม่มี
- `backend/internal/returns/handler.go`: `GET /api/v1/returns/open-loans` ใช้ภาควิชาจาก `auth.FromContext` (IF-IDP-01) ไม่รับจาก query
  ลบ `notImplemented` เพราะทุก endpoint ทำแล้ว
- `backend/internal/returns/handler_test.go`: `TestAC_RET_10_OpenLoansOwnDepartment`, ตัวช่วย `fixedProvider`, `insertDepartment`, `insertLoanOfModel`
  แก้ `TestReturns_Skeleton_RoutedThroughAuth` (test ของ T-03 ที่ AI เขียนเอง ไม่ใช่ TestTC_) ให้เลิกคาด 501 เหลือเฉพาะการตรวจ auth

### ผล test (`cd backend && go test ./internal/returns -v`)
- TestAC_RET_10_OpenLoansOwnDepartment PASS (ได้ 3 รายการของภาควิชาตนเองพอดี เรียงตามกำหนดคืน ไม่มีที่คืนแล้วหรือของภาควิชาอื่น,
  isOverdue ถูก, ภาควิชาที่ไม่มีการยืมค้างได้ `[]`)
- test เดิมของ returns ยัง PASS (10 tests)

### สิ่งที่เกือบต้องเดา
- ฐานข้อมูลทดสอบใช้ร่วมกันทุก test นับ "3 รายการพอดี" ในภาควิชา 1 ไม่ได้ จึงสร้างภาควิชาใหม่และเจ้าหน้าที่จำลองเฉพาะ test (ผ่าน `auth.Provider` ไม่แก้ mock ของกลาง)
- ไม่แบ่งหน้า ตาม plan.md Q-P1 ที่ยังรอทีมยืนยัน

---

## 2569-11-11 คำสั่ง: /implement T-15 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 7)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำเลย" -> `/implement T-15 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- ตามช่องไฟล์ที่แตะ: สร้าง `frontend/src/pages/Return/OpenLoansPage.jsx` (ตาราง รหัสครุภัณฑ์ รุ่น ผู้ยืม กำหนดคืน ปุ่ม "เลือก", รายการว่าง, error ตาม code),
  `frontend/src/pages/Return/__tests__/OpenLoansPage.test.jsx`, แก้ `frontend/src/api/client.js` (`api.returns.openLoans`),
  `frontend/src/pages/Return/ReturnPage.jsx` (รับ `initialAssetCode` แล้วกรอกและค้นหาให้ทันที)
- นอกช่องไฟล์ที่แตะ: `frontend/src/App.jsx` (แสดง OpenLoansPage ใต้ ReturnPage และเปิด ReturnPage ใหม่ด้วย key เมื่อเลือกแถว),
  `ReturnPage.test.jsx` (test ของ initialAssetCode) เพราะแอปยังไม่มี router

### ผล test (`cd frontend && npm test`)
- OpenLoansPage.test.jsx 3 tests PASS (AC-RET-10: 3 แถว + กดเลือกส่งรหัส / รายการว่าง / โหลดไม่สำเร็จ)
- ReturnPage.test.jsx เพิ่ม 1 test PASS (initialAssetCode) / รวม 14 tests PASS / `vite build` สำเร็จ
- มีคำเตือน act(...) จาก `setup.test.jsx` (test โครงของรายวิชา) เพราะ App โหลดรายการตอนเปิด test ยังผ่าน ไม่ได้แก้ไฟล์นั้น

### สิ่งที่เกือบต้องเดา
- ไม่มี router "เปิด ReturnPage ที่กรอกรหัสไว้" จึงทำด้วย state ใน App.jsx แทนการเปลี่ยนหน้า
- รายการไม่รีเฟรชเองหลังบันทึกคืน (spec ไม่ได้กำหนด) ต้องโหลดหน้าใหม่ ป้าย "เลยกำหนดคืน" อยู่ใน T-17

---

## 2569-11-11 คำสั่ง: /implement T-16 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 8)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำเลย" -> `/implement T-16 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/returns/handler_test.go`: `TestAC_RET_13_OverdueShownToStaff`
- `backend/internal/returns/repository.go`: ไม่ต้องแก้ (`isOverdue` มีใน `ListOpenLoans` จาก T-14 และ `returnedLate` มีใน `GetLoan` จาก T-11)

### ผล test (`cd backend && go test ./internal/returns -v`)
- TestAC_RET_13_OverdueShownToStaff PASS (รายละเอียดการยืมที่คืนเกินกำหนดได้ returnedLate true, รายการที่ยังไม่คืนมีเฉพาะรายการที่เลยกำหนดพร้อม isOverdue true)
- test เดิมของ returns ยัง PASS (11 tests)

### สิ่งที่เกือบต้องเดา
- ไม่มี (ใช้ภาควิชาใหม่เฉพาะ test เหมือน T-14 เพื่อให้รายการที่ยังไม่คืนนับได้แน่นอน)

---

## 2569-11-11 คำสั่ง: /implement T-17 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 8)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำเล้ยยย" -> `/implement T-17 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `frontend/src/pages/Return/LoanDetailPage.jsx`: ป้าย "คืนเกินกำหนด" เมื่อ `returnedLate`
- `frontend/src/pages/Return/OpenLoansPage.jsx`: ป้าย "เลยกำหนดคืน" ในช่องกำหนดคืนเมื่อ `isOverdue`
- `LoanDetailPage.test.jsx`, `OpenLoansPage.test.jsx`: test AC-RET-13 ฝั่งหน้าจอ

### ผล test (`cd frontend && npm test`)
- AC-RET-13: หน้ารายละเอียดแสดง "คืนเกินกำหนด" เมื่อ returnedLate (ไม่แสดงเมื่อคืนทัน) PASS
- AC-RET-13: แถวที่ isOverdue แสดง "เลยกำหนดคืน" แถวอื่นไม่แสดง PASS
- รวม 16 tests PASS / `vite build` สำเร็จ

### สิ่งที่เกือบต้องเดา
- คำบนป้าย: ใช้ "คืนเกินกำหนด" ตาม FR-RET-06 และ "เลยกำหนดคืน" ตามป้ายใน glossary.md (คำเดียวกับที่ 001 ใช้) ไม่มี mockup

---

## 2569-11-11 คำสั่ง: /implement T-18 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 9)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำเลย" -> `/implement T-18 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- `backend/internal/returns/handler_test.go`: `TestAC_RET_03_NoOpenLoanRejected`, `TestReturns_ConcurrentReturn_OneWins`, ตัวช่วย `insertItemWithoutLoan`, `itemStatus`
- `service.go`, `repository.go`: ไม่ต้องแก้ (409 `NO_OPEN_LOAN` มีตั้งแต่ T-04 / T-05 และ `SELECT ... FOR UPDATE OF l, i` กันบันทึกซ้อนอยู่แล้ว)

### ผล test (`cd backend && go test ./internal/returns -v`)
- TestAC_RET_03_NoOpenLoanRejected PASS (ว่าง / ไม่พร้อมใช้งาน: lookup และ POST ได้ 409 NO_OPEN_LOAN สถานะไม่เปลี่ยน, รหัสที่ไม่มีได้ 409)
- TestReturns_ConcurrentReturn_OneWins PASS (5 ชิ้น ยิงพร้อมกันชิ้นละ 2 request ได้ 200 หนึ่ง 409 หนึ่ง ปิดการยืมครั้งเดียว)
- รันแยก 5 รอบด้วย `-race` ผ่านทุกรอบ ไม่พบ data race / test เดิมของ returns ยัง PASS (13 tests)

### สิ่งที่เกือบต้องเดา
- ไม่มี / หมายเหตุ: test ใน package นี้ใช้รหัสครุภัณฑ์ตายตัว รันด้วย `-count` มากกว่า 1 จะชนรหัสเดิม ให้รันทีละครั้ง

---

## 2569-11-11 คำสั่ง: /implement T-19 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 9)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำต่อได้เลยครับ" -> `/implement T-19 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `frontend/src/pages/Return/ReturnPage.jsx`: `noOpenLoanMessage` "อุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน" เมื่อ lookup หรือ record ได้ `NO_OPEN_LOAN`
  (lookup: ไม่แสดงฟอร์มสภาพ / record ที่ชนกับคนอื่น: ปิดฟอร์ม)
- `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`: test AC-RET-03 บนหน้าจอ และกรณีบันทึกซ้อน

### ผล test (`cd frontend && npm test`)
- AC-RET-03: ค้นหาแล้วได้ NO_OPEN_LOAN แจ้งข้อความ และไม่แสดงฟอร์มสภาพ PASS
- T-19: บันทึกแล้วได้ NO_OPEN_LOAN แจ้งข้อความ และปิดฟอร์ม PASS
- รวม 18 tests PASS / `vite build` สำเร็จ

### สิ่งที่เกือบต้องเดา
- ข้อความใช้ตาม Then ของ AC-RET-03 ไม่มี mockup / กรณีบันทึกซ้อนใช้ข้อความเดียวกัน (ไม่ได้สร้างข้อความใหม่)

---

## 2569-11-11 คำสั่ง: /implement T-20 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 10)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำต่อได้เลย" -> `/implement T-20 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- ตามช่องไฟล์ที่แตะ: `repository.go` (`ErrForbiddenDepartment`, `checkItemDepartment` ตรวจภาควิชาก่อนเรื่องการยืมค้าง,
  `RecordReturn` / `FindOpenLoan` / `GetLoan` รับ departmentID), `service.go` (`RecordReturn` รับ departmentID),
  `handler_test.go` (`TestAC_RET_11_OtherDepartmentForbidden`)
- นอกช่องไฟล์ที่แตะ: `handler.go` (จำเป็นเพื่อส่งภาควิชาจาก `auth.FromContext` และแปลง error เป็น 403 `FORBIDDEN_DEPARTMENT`, เพิ่ม `principal`, `writeForbiddenDepartment`)

### ผล test (`cd backend && go test ./... -v`)
- TestAC_RET_11_OtherDepartmentForbidden PASS (lookup / POST / loan detail ของอุปกรณ์ภาควิชาอื่น ทั้งที่ถูกยืมและที่ว่าง ได้ 403 body มีแค่ code,
  การยืมไม่ถูกปิด สถานะยัง borrowed)
- test เดิมทั้งหมด PASS (returns 14 tests, platform, config, auth)

### สิ่งที่เกือบต้องเดา
- ลำดับ 403 กับ 409: ตรวจภาควิชาก่อนเสมอตาม plan.md ข้อ 4 อุปกรณ์ภาควิชาอื่นที่ไม่มีการยืมค้างก็ได้ 403 (ไม่บอกใบ้ว่ามีการยืมหรือไม่)
- รหัสครุภัณฑ์ที่ไม่มีในระบบยังได้ 409 `NO_OPEN_LOAN` ตามเดิม (AC-RET-03)

---

## 2569-11-11 คำสั่ง: /implement T-21 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 10)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำต่อเลยคับ" -> `/implement T-21 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `frontend/src/pages/Return/ReturnPage.jsx`: `forbiddenDepartmentMessage` "ไม่มีสิทธิ์: อุปกรณ์นี้เป็นของภาควิชาอื่น" เมื่อ lookup หรือ record ได้ `FORBIDDEN_DEPARTMENT`
  (ไม่แสดงฟอร์ม / ปิดฟอร์ม)
- `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`: test AC-RET-11 บนหน้าจอ และกรณีตอนบันทึก

### ผล test (`cd frontend && npm test`)
- AC-RET-11: ค้นหาแล้วได้ FORBIDDEN_DEPARTMENT แจ้งว่าไม่มีสิทธิ์ ไม่แสดงฟอร์ม ไม่เรียก record PASS
- T-21: บันทึกแล้วได้ FORBIDDEN_DEPARTMENT แจ้งข้อความเดียวกัน ปิดฟอร์ม PASS
- รวม 20 tests PASS / `vite build` สำเร็จ

### สิ่งที่เกือบต้องเดา
- ถ้อยคำ: AC-RET-11 บอกแค่ "ปฏิเสธว่าไม่มีสิทธิ์" AI ใช้ "ไม่มีสิทธิ์: อุปกรณ์นี้เป็นของภาควิชาอื่น" ไม่มี mockup ทีมตรวจถ้อยคำ

---

## 2569-11-11 คำสั่ง: /implement T-22 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 11)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำต่อเลย" -> `/implement T-22 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้ (ตรงกับช่องไฟล์ที่แตะ)
- `backend/internal/returns/handler.go`: `staffOnly` ครอบทุก route ของ 006 ต่อจาก `auth.Middleware` บทบาทอื่นได้ 403 `FORBIDDEN_ROLE` ก่อนถึง handler
- `backend/internal/returns/handler_test.go`: `TestAC_RET_08_StudentForbidden`

### ผล test (`cd backend && go test ./... -v`)
- TestAC_RET_08_StudentForbidden PASS (token นักศึกษาเรียกทั้ง 4 endpoint ได้ 403 body มีแค่ code, การยืมและสถานะไม่เปลี่ยน, เจ้าหน้าที่ยังได้ 200)
- test ทั้งหมดของ backend ยัง PASS (returns 15 tests)

### สิ่งที่เกือบต้องเดา
- ไม่มี (error code `FORBIDDEN_ROLE` มาจากตาราง error ของ plan.md ข้อ 4)

---

## 2569-11-11 คำสั่ง: /implement T-23 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return (เฟส 12)
- เครื่องมือ: Claude Code (`.claude/commands/implement.md`)
- คำสั่ง: "ทำเลย" -> `/implement T-23 specs/006-return/tasks.md`

### ไฟล์ที่สร้างหรือแก้
- ตามช่องไฟล์ที่แตะ: `handler.go` (ทุกคำตอบผ่าน `respond` ที่เขียน audit ก่อนตอบ, `errorResponse` รวมการแปลง error,
  `staffOnly` เขียน audit forbidden ด้วย action ของ endpoint), `service.go` (ส่ง AuditEntry เข้า transaction),
  `handler_test.go` (`TestAC_RET_07_ReturnAudited`, ตัวช่วย `auditRowsOf`)
- นอกช่องไฟล์ที่แตะ: สร้าง `backend/internal/returns/audit.go` (ยังไม่มี package audit ของ 001 จึงเขียนใน returns ใช้คอลัมน์ DM-07 เดิม),
  แก้ `repository.go` (`RecordReturn` เขียน audit "return" ใน transaction เดียวกัน ตามช่อง "เสร็จเมื่อ")

### ผล test (`cd backend && go test ./... -v`)
- TestAC_RET_07_ReturnAudited PASS (lookup / return / detail / list ของเจ้าหน้าที่ outcome ok ครบตามลำดับ, บันทึกซ้ำ 409 = error,
  lookup เก็บเงื่อนไขที่ค้นใน query, audit ของการคืนสำเร็จมีเวลาเท่า returned_at = transaction เดียวกัน, นักศึกษาได้ outcome forbidden)
- test ทั้งหมดของ backend ยัง PASS (returns 16 tests)

### สิ่งที่เกือบต้องเดา (ทีมต้องยืนยัน)
- การจับคู่ outcome: 2xx = ok, 403 = forbidden, อื่น ๆ (409, 422, 404, 503) = error เพราะ DM-07 มีแค่ 3 ค่า
- เขียน audit ไม่สำเร็จตอบ 503 ไม่ส่งข้อมูล (fail closed ตามแนว plan 6.4 ของ 001)
- ไม่เก็บสถานะก่อน/หลังใน audit (รอ Q-09 T-26)

---

## 2569-11-11 ใช้โครงกลางจาก 001 และย้ายโค้ดเข้า package ของ 006 (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code
- คำสั่ง: "ใช้โครงกลางจาก 001 (Go, PostgreSQL, auth จำลอง) ห้ามสร้างโครงแอปใหม่เอง และเขียนโค้ดเฉพาะใน package ของตัวเอง"

### สิ่งที่ AI ตรวจ
- ไม่ได้สร้างโครงแอปใหม่: ใช้ `platform.NewRouter`, `platform.OpenDB/Migrate/Seed`, `auth.Middleware`, `auth.NewMockProvider`, PostgreSQL + testcontainers ของ 001 ตลอด
- ไฟล์นอก package ของ 006 ที่เคยแก้: `platform/http.go`, `platform/migrate_test.go`, `config/config_test.go`, `frontend/src/__tests__/config.test.js`

### สิ่งที่แก้
- คืน `backend/internal/platform/http.go` และ `migrate_test.go` เป็นของ 001 (เหมือน commit 5edd87d ทุกตัวอักษร)
- `backend/cmd/server/main.go`: ใช้ `platform.NewRouter(db)` เดิมแล้วผูก route ของ 006 ด้วย `returnsHandler.Routes(provider)(router)`
- ย้าย test ของ migration 0006 และ config ของ 006 ไป `backend/internal/returns/migration_test.go` ลบ `config/config_test.go` (ไฟล์ที่ 006 สร้างเอง)
- `handler_test.go`: ตัวช่วย `newRouterWith` ใช้ router กลางของ 001 แบบเดียวกับ main.go
- ย้าย `frontend/src/__tests__/config.test.js` ไป `frontend/src/pages/Return/__tests__/config.test.js`

### ที่ยังอยู่นอก package ของ 006 (ตั้งใจ)
- `config/config.go`, `frontend/src/config.js`: AGENTS.md กฎข้อ 9 ให้ค่าจาก spec อยู่ที่นี่ที่เดียว
- `cmd/server/main.go`, `frontend/src/api/client.js`, `frontend/src/App.jsx`, migration/seed: จุดที่โครงของ 001 ให้ฟีเจอร์เพิ่มเอง
- `frontend/vite.config.js`: ทีมสั่งให้เอา rewrite ออกตอน T-06

### ผล test
- `cd backend && go test ./...`: PASS ทั้งหมด (platform ของ 001 เท่าเดิม, returns 20 tests)
- `cd frontend && npm test`: 20 tests PASS

---

## 2569-11-11 คำสั่ง: /verify (006-return)

- ผู้สั่ง: 660710737 นายวรรัตน์ หลักฐาน
- ฟีเจอร์: 006-return
- เครื่องมือ: Claude Code (`.claude/commands/verify.md`)
- คำสั่ง: "ตอนท้ายรัน /verify ทำ rtm.md ของฟีเจอร์ตัวเอง" -> `/verify specs/006-return/spec.md`
- ผลลัพธ์: `specs/006-return/rtm.md` (รอบแรก) ไม่ได้แก้โค้ด test spec plan tasks

### ผล test (ดูอย่างเดียว)
- `cd backend && go test ./... -v`: 30 ผ่าน 0 ไม่ผ่าน (auth 3, platform 7, returns 20)
- `cd frontend && npm test`: 20 ผ่าน 0 ไม่ผ่าน

### ตารางไปข้างหน้า 18 แถว
- ครบ 9 / ยังไม่ถึง 3 (NFR-DATA-01 รอ 001 T-16, NFR-SEC-01 และ CON-TECH-02 รอ TLS ของ deploy กลาง) / รอ 1 (DOM-REQ-01 รอ Q-08) /
  ช่องโหว่ 4 (FR-RET-01, FR-RET-07, DOM-AUDIT-01, NFR-AVAIL-01) / ยกเลิก 1 (FR-RET-05)

### ข้อค้นพบใหม่
- F-01 เดา: การจับคู่ outcome ของ audit (409/422/404 = error) ไม่มีที่มาใน spec/DM-07
- F-02 เดา: fail closed 503 เมื่อเขียน audit ไม่ได้ ใช้แนว Q-P4 ของ 001 ที่ยังไม่มีคำตอบ
- F-03 test อ่อน: การต่อ OpenLoansPage -> ReturnPage ใน App.jsx ไม่มี test (AC-RET-10)
- F-04 FR ไม่มี AC: NFR-AVAIL-01
