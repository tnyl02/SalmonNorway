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

