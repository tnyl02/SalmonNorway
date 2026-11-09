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
