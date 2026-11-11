# Tasks: ส่งคำขอยืมอุปกรณ์ (Borrow Request)
Feature: 002-borrow-request | Spec ID: SPEC-BRQ-002 (spec.md Draft v3.1)
อ้างอิง: plan.md Draft v1.1, specs/changes/{CR-003.md, CR-004.md}, specs/000-shared/{constraints.md v5, data-model.md v9}
Updated: 2569-11-11

สรุป: **14 เฟส 23 task** โดยเฟส 0 ถึง 12 เป็นงานที่ทำได้เลย ส่วนเฟสสุดท้ายคือเฟสรอคำตอบ
**5 task รอ Open Questions** (T-19 ถึง T-23) ซึ่งเป็นงานหน้าจอทั้งหมดและกรณีพิเศษ 3 เรื่อง
**เฟสแรกที่ลองใช้งานได้จริงคือเฟส 1** (ส่งคำขอสำเร็จ) แต่ลองผ่าน API ด้วย `curl` เพราะหน้าจอติด Q-BRQ-14
ทุกเฟสจึงลองด้วย `curl` จนกว่าจะได้คำตอบ Q-BRQ-14

---

## เฟส 0: ฐานราก (ไม่มี AC — เป็นงานพื้นฐานของทุกเฟส)

### T-01 สร้าง migration 0003 ตาราง borrow_request
- รองรับ: FR-BRQ-04, FR-BRQ-05, CON-TECH-01, CON-DATA-03, DOM-REQ-01
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-04 ถึง T-18
- ไฟล์ที่แตะ: `backend/db/migrations/0003_borrow_request.up.sql`, `backend/db/migrations/0003_borrow_request.down.sql`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: migration รันผ่าน และ test ยืนยันได้ว่า (ก) สถานะนอก 8 ค่าถูกปฏิเสธ (ข) unique index กันคำขอซ้ำของผู้ขอคนเดิมบนชิ้นเดิมใน 4 สถานะที่ยังไม่จบ (ค) `audit_log.request_id` มี foreign key แล้ว (ง) ค่า `request` ใช้ใน `audit_log.action` ได้
- สถานะ: เสร็จ

> **หมายเหตุ**: 003 ติดงาน 17 task และ 004 รอตารางนี้อยู่ ให้ส่ง T-01 ขึ้น main ทันทีที่ผ่าน โดยไม่ต้องรอ task อื่น

### T-02 เพิ่มค่าที่ตั้งได้ของฟีเจอร์นี้ใน config
- รองรับ: DOM-LOAN-01, FR-BRQ-02, FR-BRQ-06, FR-BRQ-10, CON-TECH-03
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-06 ถึง T-14
- ไฟล์ที่แตะ: `backend/internal/config/config.go`, `frontend/src/config.js`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: ค่าทั้ง 10 ตัวในตาราง "ค่าที่ตั้งได้" ของ plan ข้อ 5 อยู่ในไฟล์เดียว พร้อมคอมเมนต์อ้าง ID และไม่มีตัวเลขเหล่านี้ปรากฏที่อื่นในโค้ด
- สถานะ: เสร็จ

### T-03 เพิ่มข้อมูลตัวอย่างสำหรับลองระบบ
- รองรับ: FR-BRQ-01, FR-BRQ-06, FR-BRQ-07
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นข้อมูลที่จุดตรวจของทุกเฟสใช้
- ไฟล์ที่แตะ: `backend/db/seed/seed.sql`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: seed มีครบ 3 กรณีสำหรับลองมือ — อุปกรณ์ว่างของภาควิชา 1, อุปกรณ์ว่างของภาควิชา 2 (ไว้ลองข้ามภาควิชา) และ loan ที่เลยกำหนดคืนของ `student-002`
- สถานะ: เสร็จ

#### จุดตรวจเฟส 0
- test ที่ต้องผ่านทั้งหมด: `TestMigration_BorrowRequestConstraints`, `TestMigration_DuplicateRequestBlocked` และ test เดิมของ 001 ทั้งหมดต้องยังผ่าน
- เปิดระบบ: `cd backend && go test ./... -v` แล้ว `docker compose -f deploy/docker-compose.yml up -d db`
- คนลองเองตาม AC: ยังไม่มี AC ในเฟสนี้ ให้ลองว่า
  Given รัน migration ครบทั้ง 3 ไฟล์
  When เปิดตาราง `borrow_request` ด้วย psql แล้ว `INSERT` แถวสถานะ `pending` (ค่าที่ไม่มีในรายการ)
  Then ฐานข้อมูลต้องปฏิเสธด้วย CHECK constraint
- ยังขาด: ไม่มี

---

## เฟส 1: ส่งคำขอยืมสำเร็จ (AC-BRQ-04)

### T-04 สร้าง repository และ service ของการสร้างคำขอ
- รองรับ: FR-BRQ-04, DOM-REQ-01, CON-DATA-01
- ตรวจด้วย: AC-BRQ-04
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/repository.go`, `backend/internal/borrowrequest/service.go`, `backend/internal/borrowrequest/model.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01, T-02
- เสร็จเมื่อ: `TestAC_BRQ_04` ผ่าน — สร้างคำขอสถานะ `requested` พร้อมเวลาที่ส่งตาม Asia/Bangkok และ `equipment_item.status` ยังเป็น `available`
- สถานะ: เสร็จ

### T-05 สร้าง endpoint POST /api/v1/loan-requests และผูกเข้า router
- รองรับ: FR-BRQ-04, CON-TECH-02, CON-TECH-03, IF-IDP-01
- ตรวจด้วย: AC-BRQ-04
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/handler.go`, `backend/internal/borrowrequest/dto.go`, `backend/cmd/server/main.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_BRQ_04` ที่ยิงผ่าน HTTP ผ่าน และเรียกโดยไม่มี token ได้ 401 `UNAUTHENTICATED` ตาม CON-TECH-03
- สถานะ: เสร็จ รอทีมตรวจ

#### จุดตรวจเฟส 1
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_04` และ test ของเฟส 0 กับของ 001 ต้องยังผ่าน
- เปิดระบบ: `cd backend && go run ./cmd/server` (หน้าจอยังไม่มีในเฟสนี้ เพราะติด Q-BRQ-14)
- คนลองเองตาม AC-BRQ-04:
  Given ใช้ข้อมูลตัวอย่างจาก T-03 อุปกรณ์ `CS-NB-0001` ของภาควิชา 1 สถานะว่าง และล็อกอินเป็น `student-001`
  When ยิง `curl -X POST localhost:8000/api/v1/loan-requests -H "Authorization: Bearer mock-student-001" -d '{...ข้อมูลครบ...}'`
  Then ต้องได้ 201 พร้อม `{"requestId":1,"status":"requested","submittedAt":"..."}` และเมื่อเปิดดูอุปกรณ์ชิ้นนั้นในหน้าของ 001 สถานะต้องยังเป็น "ว่าง"
- ยังขาด: หน้าจอฟอร์มยังไม่มี (รอ Q-BRQ-14 ดู T-19)

---

## เฟส 2: กรอกไม่ครบแล้วส่งไม่ได้ (AC-BRQ-03)

### T-06 ตรวจความครบถ้วนของ 5 ช่องในฟอร์ม
- รองรับ: FR-BRQ-03, CON-TECH-03
- ตรวจด้วย: AC-BRQ-03
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/validate.go`, `backend/internal/borrowrequest/handler.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-05
- เสร็จเมื่อ: `TestAC_BRQ_03` ผ่าน — ขาดช่องใดช่องหนึ่งใน 5 ช่องต้องได้ 422 `VALIDATION_FAILED` พร้อมรายชื่อช่องที่ขาด และไม่มีแถวใหม่ในตาราง
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 2
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_03` และ test ของเฟส 0-1
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-03:
  Given ล็อกอินเป็น `student-001`
  When ยิง POST โดยไม่ส่ง `endAt`
  Then ต้องได้ 422 `VALIDATION_FAILED` ที่บอกว่าช่อง `endAt` ยังขาด และนับแถวในตารางแล้วต้องไม่เพิ่ม
- ยังขาด: การไฮไลต์ช่องที่ขาดบนหน้าจอ (รอ Q-BRQ-14 ดู T-20)

---

## เฟส 3: เลือกวัตถุประสงค์การยืม (AC-BRQ-02, AC-BRQ-14)

### T-07 รับค่าวัตถุประสงค์ 3 ค่า และบังคับพิมพ์ระบุเมื่อเลือก "อื่น ๆ"
- รองรับ: FR-BRQ-02
- ตรวจด้วย: AC-BRQ-02, AC-BRQ-14
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/validate.go`, `backend/internal/borrowrequest/service.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestAC_BRQ_02` และ `TestAC_BRQ_14` ผ่าน — ค่านอก 3 ค่าถูกปฏิเสธ, เลือก `other` โดยไม่พิมพ์ระบุได้ 422, และเมื่อพิมพ์ระบุต้องเก็บทั้งสองค่า
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 3
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_02`, `TestAC_BRQ_14` และ test ของเฟส 0-2
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-02 และ AC-BRQ-14:
  Given ล็อกอินเป็น `student-001`
  When ยิง POST ด้วย `purpose: "other"` แต่ไม่ส่ง `purposeOther` แล้วยิงอีกครั้งโดยส่งข้อความมาด้วย
  Then ครั้งแรกต้องได้ 422 ครั้งที่สองต้องได้ 201 และเปิดดูแถวในตารางต้องเห็นทั้ง `purpose='other'` และข้อความที่พิมพ์
- ยังขาด: dropdown บนหน้าจอ (รอ Q-BRQ-14 ดู T-19)

---

## เฟส 4: ยืมปกติห้ามคืนเกิน 16:00 น. (AC-BRQ-13)

### T-08 บังคับกฎเวลาคืนของการยืมแบบปกติ
- รองรับ: FR-BRQ-03, DOM-LOAN-01
- ตรวจด้วย: AC-BRQ-13
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/validate.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-06
- เสร็จเมื่อ: `TestAC_BRQ_13` ผ่าน — `loanType=normal` ที่ `endAt` เกิน 16:00 น. ของวันเดียวกับ `startAt` ต้องได้ 422 โดยค่า 16 อ่านจาก config ไม่ฝังในโค้ด
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 4
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_13` และ test ของเฟส 0-3
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-13:
  Given ล็อกอินเป็น `student-001` และเลือกการยืมแบบปกติ
  When ยิง POST ด้วย `startAt` วันนี้ 09:00 และ `endAt` วันนี้ 18:00
  Then ต้องได้ 422 พร้อมข้อความว่าการยืมแบบปกติต้องคืนภายใน 16:00 น. ของวันเดียวกัน
- ยังขาด: ข้อความเตือนบนหน้าจอ (รอ Q-BRQ-14 ดู T-20)

---

## เฟส 5: อุปกรณ์ถูกยืมไประหว่างกรอก (AC-BRQ-01)

### T-09 ตรวจสถานะอุปกรณ์ซ้ำในทรานแซกชันเดียวกับการสร้างคำขอ
- รองรับ: FR-BRQ-01, CON-DATA-01, CON-DATA-03, CON-TECH-03
- ตรวจด้วย: AC-BRQ-01
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/repository.go`, `backend/internal/borrowrequest/service.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-05
- เสร็จเมื่อ: `TestAC_BRQ_01` ผ่าน — เปลี่ยนสถานะเป็น `borrowed` ก่อนกดส่ง ต้องได้ 409 `ITEM_NOT_AVAILABLE` พร้อมสถานะล่าสุด และไม่มีแถวใหม่
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 5
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_01` และ test ของเฟส 0-4
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-01:
  Given อุปกรณ์ `CS-NB-0001` สถานะว่าง
  When แก้สถานะเป็น `borrowed` ด้วย psql แล้วยิง POST คำขอของชิ้นนั้น
  Then ต้องได้ 409 `ITEM_NOT_AVAILABLE` และ body บอกสถานะล่าสุดว่า "ถูกยืม"
- ยังขาด: การแจ้งสถานะล่าสุดบนหน้าจอ (รอ Q-BRQ-14 ดู T-20)

---

## เฟส 6: นักศึกษาต่างภาควิชาส่งคำขอไม่ได้ (AC-BRQ-07)

### T-10 บังคับขอบเขตภาควิชาที่ชั้น repository
- รองรับ: FR-BRQ-07, DOM-ACCESS-01, IF-IDP-01, CON-TECH-03
- ตรวจด้วย: AC-BRQ-07
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/repository.go`, `backend/internal/borrowrequest/handler.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-09
- เสร็จเมื่อ: `TestAC_BRQ_07` ผ่าน — `student-101` (ภาควิชา 2) ส่งคำขออุปกรณ์ภาควิชา 1 ต้องได้ 403 `FORBIDDEN_DEPARTMENT` และ body ต้องไม่มีชื่อรุ่น สถานะ หรือข้อมูลอุปกรณ์ใด ๆ
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 6
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_07` และ test ของเฟส 0-5
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-07:
  Given ล็อกอินเป็น `student-101` ซึ่งอยู่ภาควิชา 2
  When ยิง POST คำขออุปกรณ์ `CS-NB-0001` ของภาควิชา 1
  Then ต้องได้ 403 `FORBIDDEN_DEPARTMENT` และอ่าน body แล้วต้องไม่รู้เลยว่าอุปกรณ์ชิ้นนั้นเป็นรุ่นอะไรหรือสถานะอะไร
- ยังขาด: ลำดับการตรวจเมื่อข้ามภาควิชาและอุปกรณ์ถูกยืมพร้อมกัน (รอ Q-BRQ-15 ดู T-21)

---

## เฟส 7: ส่งคำขอซ้ำของตัวเองไม่ได้ (AC-BRQ-11)

### T-11 ปฏิเสธคำขอซ้ำของผู้ขอคนเดิมที่คำขอเก่ายังไม่จบ
- รองรับ: FR-BRQ-05, DOM-REQ-01, CON-TECH-03
- ตรวจด้วย: AC-BRQ-11
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/service.go`, `backend/internal/borrowrequest/repository.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-09
- เสร็จเมื่อ: `TestAC_BRQ_11` ผ่าน — ส่งซ้ำขณะสถานะ `requested` หรือ `approved` ได้ 409 `DUPLICATE_REQUEST` / หลังคำขอเดิมเป็น `cancelled` หรือ `rejected` ส่งใหม่ได้ / และเมื่อยิงพร้อมกัน 2 ครั้ง ต้องมีแถวเดียวเพราะ unique index จาก T-01 กันไว้
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 7
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_11` และ test ของเฟส 0-6
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-11:
  Given `student-001` ส่งคำขอของ `CS-NB-0002` ไปแล้ว 1 ใบ สถานะ `requested`
  When ยิง POST ของชิ้นเดิมซ้ำ แล้วแก้สถานะใบเดิมเป็น `cancelled` ด้วย psql แล้วยิงอีกครั้ง
  Then ครั้งแรกต้องได้ 409 `DUPLICATE_REQUEST` ครั้งหลังต้องได้ 201
- ยังขาด: ไม่มี

---

## เฟส 8: คนที่คืนของช้าส่งคำขอใหม่ไม่ได้ (AC-BRQ-06, AC-BRQ-12)

### T-12 เพิ่มนาฬิกาที่ฉีดได้ และกฎปฏิเสธเมื่อมีของคืนเกินกำหนด
- รองรับ: FR-BRQ-06, DOM-ACCESS-01, CON-TECH-03
- ตรวจด้วย: AC-BRQ-06
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/clock.go`, `backend/internal/borrowrequest/service.go`, `backend/internal/borrowrequest/repository.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-11
- เสร็จเมื่อ: `TestAC_BRQ_06` ผ่าน — ผู้ขอที่มี loan ซึ่ง `returned_at` ว่างและเลย 16:00 น. ของ `due_at` แล้ว ต้องได้ 409 `HAS_OVERDUE_LOAN` และต้องนับเฉพาะอุปกรณ์ของภาควิชาตนเอง
- สถานะ: พร้อมทำ

### T-13 กฎปลดล็อกเมื่อถึง 00:00 น. ของวันถัดไปหลังคืน
- รองรับ: FR-BRQ-06
- ตรวจด้วย: AC-BRQ-12
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/service.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-12
- เสร็จเมื่อ: `TestAC_BRQ_12` ผ่าน — ใช้นาฬิกาที่ตั้งได้ ยิงเวลา 23:59 ของวันที่คืนต้องถูกปฏิเสธ และ 00:01 ของวันถัดไปต้องผ่าน โดยเทียบตามเขตเวลา Asia/Bangkok
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 8
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_06`, `TestAC_BRQ_12` และ test ของเฟส 0-7
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-06:
  Given ใช้ข้อมูลตัวอย่างจาก T-03 ที่ `student-002` มีอุปกรณ์เลยกำหนดคืนค้างอยู่
  When ล็อกอินเป็น `student-002` แล้วยิง POST คำขอใหม่
  Then ต้องได้ 409 `HAS_OVERDUE_LOAN` พร้อมเหตุผลว่ามีของค้างเกินกำหนด
- ยังขาด: ไม่มี

---

## เฟส 9: ยืมค้างคืนพร้อมแนบใบบันทึกข้อความ (AC-BRQ-10)

### T-14 รับไฟล์แนบของการยืมค้างคืน พร้อมตรวจชนิดและขนาด
- รองรับ: FR-BRQ-10, CON-TECH-03
- ตรวจด้วย: AC-BRQ-10
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/attachment.go`, `backend/internal/borrowrequest/handler.go`, `backend/internal/borrowrequest/validate.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-07
- เสร็จเมื่อ: `TestAC_BRQ_10` ผ่าน — `loanType=overnight` ที่ไม่แนบไฟล์ได้ 422, ไฟล์เกิน 10MB หรือชนิดนอก 3 ชนิดได้ 422, และไฟล์ที่ผ่านถูกเก็บลงโฟลเดอร์ตาม config พร้อมบันทึก path ในตาราง
- สถานะ: พร้อมทำ

### T-15 สร้าง endpoint เปิดดูไฟล์แนบ จำกัดสิทธิ์
- รองรับ: FR-BRQ-10, DOM-PRIV-01, IF-IDP-01
- ตรวจด้วย: AC-BRQ-10
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/handler.go`, `backend/cmd/server/main.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-14
- เสร็จเมื่อ: เจ้าหน้าที่และเจ้าของคำขอเปิดไฟล์ได้ 200 ส่วนนักศึกษาคนอื่นได้ 403 โดยมี test ครอบทั้ง 3 กรณี
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 9
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_10` และ test ของเฟส 0-8
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-10:
  Given ล็อกอินเป็น `student-001` และเลือกการยืมแบบค้างคืน
  When ยิง POST แบบ multipart โดยไม่แนบไฟล์ แล้วยิงใหม่พร้อมแนบ PDF 1MB
  Then ครั้งแรกต้องได้ 422 ครั้งหลังต้องได้ 201 และเปิดโฟลเดอร์ไฟล์แนบแล้วต้องเห็นไฟล์นั้น
- ยังขาด: ช่องอัปโหลดบนหน้าจอ (รอ Q-BRQ-14 ดู T-19)

---

## เฟส 10: ทุกคำขอถูกบันทึกไว้ตรวจย้อนหลังได้ (AC-BRQ-08)

### T-16 เขียน audit log ของการส่งคำขอทุกผลลัพธ์
- รองรับ: FR-BRQ-08, DOM-AUDIT-01
- ตรวจด้วย: AC-BRQ-08
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/service.go`, `backend/internal/audit/` (ของกลาง)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-12
- เสร็จเมื่อ: `TestAC_BRQ_08` ผ่าน — ส่งสำเร็จ 1 ครั้งและถูกปฏิเสธ 1 ครั้ง ต้องได้ 2 แถวใน `audit_log` ที่ `action='request'` และ `outcome` ต่างกัน พร้อมผู้ขอ เวลา และรหัสครุภัณฑ์
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 10
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_08` และ test ของเฟส 0-9
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-08:
  Given ล้างตาราง `audit_log` ด้วย psql ไม่ได้ (append-only) จึงจดจำนวนแถวก่อนเริ่ม
  When ยิง POST ที่สำเร็จ 1 ครั้ง และ POST ที่ถูกปฏิเสธเพราะข้ามภาควิชา 1 ครั้ง
  Then แถวใน `audit_log` ต้องเพิ่ม 2 แถว โดยมี `outcome` เป็น `ok` และ `forbidden` อย่างละ 1
- ยังขาด: ไม่มี

---

## เฟส 11: ระบบขัดข้องแล้วไม่เกิดคำขอค้างหรือซ้ำ (AC-BRQ-09)

### T-17 จัดการกรณีบันทึกไม่สำเร็จและการส่งซ้ำ
- รองรับ: FR-BRQ-09, CON-TECH-03
- ตรวจด้วย: AC-BRQ-09
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/service.go`, `backend/internal/borrowrequest/handler.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-11
- เสร็จเมื่อ: `TestAC_BRQ_09` ผ่าน — ใช้ repository ปลอมที่ล้มเหลวครั้งแรก ครั้งแรกได้ 503 `UPSTREAM_UNAVAILABLE` และไม่มีแถวค้าง ครั้งที่สองสำเร็จและมีแถวเดียว
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 11
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_09` และ test ของเฟส 0-10
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-09:
  Given หยุด container ฐานข้อมูลด้วย `docker compose stop db`
  When ยิง POST คำขอ แล้วเปิด db กลับมาแล้วยิงซ้ำ
  Then ครั้งแรกต้องได้ 503 `UPSTREAM_UNAVAILABLE` ครั้งหลังต้องได้ 201 และนับแถวแล้วต้องมีใบเดียว
- ยังขาด: ปุ่ม "ลองใหม่" บนหน้าจอ (รอ Q-BRQ-14 ดู T-20)

---

## เฟส 12: เจ้าหน้าที่เห็นลำดับคำขอของอุปกรณ์ชิ้นนั้น (AC-BRQ-05)

### T-18 สร้าง endpoint ลำดับคำขอให้ 001 เรียกใช้
- รองรับ: FR-BRQ-05, DOM-REQ-01, DOM-PRIV-01, IF-IDP-01
- ตรวจด้วย: AC-BRQ-05
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/handler.go`, `backend/internal/borrowrequest/repository.go`, `backend/cmd/server/main.go`
- แบบหน้าจอ: ไม่มี (การแสดงผลเป็นงานของ 001 ตาม CR-003)
- ต้องทำหลัง: T-11
- เสร็จเมื่อ: `TestAC_BRQ_05` ผ่าน — เจ้าหน้าที่เรียกแล้วได้ลำดับเรียงตาม `submitted_at` พร้อมชื่อผู้ขอ ส่วนนักศึกษาเรียกต้องได้ 403
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 12
- test ที่ต้องผ่านทั้งหมด: `TestAC_BRQ_05` และ test ของเฟส 0-11
- เปิดระบบ: `cd backend && go run ./cmd/server`
- คนลองเองตาม AC-BRQ-05:
  Given `student-001` และ `student-002` ส่งคำขอของ `CS-NB-0003` ชิ้นเดียวกัน ตามลำดับ
  When ล็อกอินเป็น `staff-001` แล้วยิง `GET /api/v1/equipment/items/CS-NB-0003/requests`
  Then ต้องได้ 2 รายการ ลำดับ 1 คือ `student-001` และเมื่อยิงด้วย token ของนักศึกษาต้องได้ 403
- ยังขาด: การแสดงลำดับบนหน้ารายละเอียดอุปกรณ์เป็นงานของ 001 ซึ่งต้องทำ CR ก่อน

---

## เฟสรอคำตอบ

### T-19 สร้างหน้าจอแบบฟอร์มคำขอยืม
- รองรับ: FR-BRQ-01, FR-BRQ-02, FR-BRQ-03, FR-BRQ-10, NFR-USE-01
- ตรวจด้วย: AC-BRQ-02, AC-BRQ-03, AC-BRQ-10, AC-BRQ-14 (ส่วนที่ผู้ใช้เห็นบนหน้าจอ)
- ไฟล์ที่แตะ: `frontend/src/pages/BorrowRequest/BorrowRequestForm.jsx`, `frontend/src/api/client.js`, `frontend/src/config.js`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (spec ยังไม่มีหมวดหน้าจอ)
- ต้องทำหลัง: T-14
- เสร็จเมื่อ: กรอกฟอร์มบนหน้าจอแล้วส่งคำขอสำเร็จ และ test หน้าจอด้วย client จำลองผ่าน
- สถานะ: รอ Q-BRQ-14

### T-20 แสดงข้อความแจ้งผลและข้อผิดพลาดบนหน้าจอ
- รองรับ: FR-BRQ-01, FR-BRQ-03, FR-BRQ-06, FR-BRQ-09, FR-BRQ-10, CON-TECH-03
- ตรวจด้วย: AC-BRQ-01, AC-BRQ-03, AC-BRQ-06, AC-BRQ-09, AC-BRQ-13 (ส่วนข้อความที่ผู้ใช้เห็น)
- ไฟล์ที่แตะ: `frontend/src/pages/BorrowRequest/BorrowRequestForm.jsx`, `frontend/src/components/`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-19
- เสร็จเมื่อ: ทุกรหัส error ของ CON-TECH-03 ที่ฟีเจอร์นี้ใช้ มีข้อความแสดงบนหน้าจอที่ทีมตกลงไว้
- สถานะ: รอ Q-BRQ-14

### T-21 กำหนดลำดับการตรวจระหว่างข้ามภาควิชากับอุปกรณ์ไม่ว่าง
- รองรับ: FR-BRQ-01, FR-BRQ-07
- ตรวจด้วย: AC-BRQ-01, AC-BRQ-07 (กรณีที่เงื่อนไขซ้อนกัน)
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/service.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-10
- เสร็จเมื่อ: มี test ที่ยืนยันลำดับการตรวจตามที่ทีมตัดสิน (403 หรือ 409 มาก่อน)
- สถานะ: รอ Q-BRQ-15

### T-22 แสดงผลของคำขอที่ส่งซ้อนให้ผู้ขอทราบ
- รองรับ: FR-BRQ-05
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นส่วนขยายของ AC-BRQ-05
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/dto.go`, `frontend/src/pages/BorrowRequest/`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-18
- เสร็จเมื่อ: ผู้ขอเห็นข้อความที่ทีมตกลงว่าจะบอกอะไรเมื่อคำขอถูกส่งซ้อน
- สถานะ: รอ Q-SH-06

### T-23 กำหนดพฤติกรรมเมื่อระบบยืนยันตัวตนเชื่อมต่อไม่ได้
- รองรับ: IF-IDP-01, CON-TECH-03
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-05
- ไฟล์ที่แตะ: `backend/internal/borrowrequest/handler.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-05
- เสร็จเมื่อ: มี test ที่ยืนยันพฤติกรรมตามที่ทีมตัดสิน
- สถานะ: รอ Q-SH-08

---

## ตารางตรวจความครบ

### 1. AC ทุกข้อมี task ตรวจ

| AC ID | task ที่ตรวจ AC นี้ |
|---|---|
| AC-BRQ-01 | T-09 (เพิ่มเติมบนหน้าจอ T-20, กรณีซ้อน T-21) |
| AC-BRQ-02 | T-07 (หน้าจอ T-19) |
| AC-BRQ-03 | T-06 (หน้าจอ T-19, T-20) |
| AC-BRQ-04 | T-04, T-05 |
| AC-BRQ-05 | T-18 |
| AC-BRQ-06 | T-12 (หน้าจอ T-20) |
| AC-BRQ-07 | T-10 (กรณีซ้อน T-21) |
| AC-BRQ-08 | T-16 |
| AC-BRQ-09 | T-17 (หน้าจอ T-20) |
| AC-BRQ-10 | T-14, T-15 (หน้าจอ T-19) |
| AC-BRQ-11 | T-11 (unique index จาก T-01) |
| AC-BRQ-12 | T-13 |
| AC-BRQ-13 | T-08 (หน้าจอ T-20) |
| AC-BRQ-14 | T-07 (หน้าจอ T-19) |

### 2. Constraint ทุกข้อมี task ทำให้เป็นจริง

| Constraint ID | task ที่ทำให้เป็นจริง |
|---|---|
| CON-TECH-01 | T-01, T-03 |
| CON-TECH-02 | T-05 (ทุก endpoint ผ่าน HTTPS ที่ proxy ระดับระบบ ยังไม่มีเจ้าของเหมือน 001) |
| CON-TECH-03 | T-05, T-06, T-09, T-10, T-11, T-12, T-14, T-17 (ทุก error ใช้รหัสกลาง) |
| CON-DATA-01 | T-04, T-09 |
| CON-DATA-03 | T-01, T-09 |
| DOM-PRIV-01 | T-15, T-18 |
| DOM-AUDIT-01 | T-16 |
| DOM-ACCESS-01 | T-10, T-12 |
| DOM-LOAN-01 | T-02, T-08 |
| DOM-REQ-01 | T-01, T-04, T-11, T-18 |
| IF-IDP-01 | T-05, T-10, T-15, T-18, T-23 |

> IF-NOT-01 ไม่อยู่ในรายการ Constraint ของ spec ฟีเจอร์นี้ (การแจ้งเตือนเป็นของ 004 ตาม RQ-21)

### 3. FR ทุกข้อมี AC ตรวจ

| FR ID | AC ที่ตรวจ FR นี้ | เฟส |
|---|---|---|
| FR-BRQ-01 | AC-BRQ-01 | 5 |
| FR-BRQ-02 | AC-BRQ-02, AC-BRQ-14 | 3 |
| FR-BRQ-03 | AC-BRQ-03, AC-BRQ-13 | 2, 4 |
| FR-BRQ-04 | AC-BRQ-04 | 1 |
| FR-BRQ-05 | AC-BRQ-05, AC-BRQ-11 | 7, 12 |
| FR-BRQ-06 | AC-BRQ-06, AC-BRQ-12 | 8 |
| FR-BRQ-07 | AC-BRQ-07 | 6 |
| FR-BRQ-08 | AC-BRQ-08 | 10 |
| FR-BRQ-09 | AC-BRQ-09 | 11 |
| FR-BRQ-10 | AC-BRQ-10 | 9 |

---

## สิ่งที่ยังไม่ทำ

คัดลอก Open Questions ที่ยังเปิดอยู่จาก spec.md พร้อม task ที่รออยู่

- **Q-BRQ-14** spec ยังไม่มีหมวด "หน้าจอ (UI)" และไม่มี mockup ข้อความบนหน้าจอเอามาจากไหน -> ทีม
  *task ที่รอ: T-19, T-20*
- **Q-BRQ-15** คำขอข้ามภาควิชาของอุปกรณ์ที่ถูกยืมอยู่ จะแจ้งสถานะหรือปฏิเสธเงียบ -> ทีม / นักคอมพิวเตอร์
  *task ที่รอ: T-21*
- **Q-SH-06 (ส่วนที่เหลือ)** ต้องแจ้งผู้ขอที่ถูกปฏิเสธอัตโนมัติตาม RQ-21 ไหม และอนุมัติตามลำดับหรือข้ามได้ -> ทีม / เจ้าของ 004
  *task ที่รอ: T-22*
- **Q-SH-08** ถ้าระบบยืนยันตัวตนของมหาวิทยาลัยเชื่อมต่อไม่ได้ ระบบต้องทำอย่างไร -> นักเทคโนโลยีสารสนเทศ
  *task ที่รอ: T-23*
- **Q-SH-03** โปรโตคอลของระบบยืนยันตัวตน -> ระหว่างรอใช้ auth จำลองของ 001 ไม่มี task ของฟีเจอร์นี้รออยู่

### ค้างจาก CR-003 และ CR-004 (ไม่ใช่ Open Question)
- backlog BL-23 ยังเป็น Could "เลื่อนไว้ก่อน" ทั้งที่ส่วน RQ-20 ทำใน 002 แล้ว
- Use Case Spec UC-02 บน Miro ยังเขียนว่า NFR-PERF-01 ครอบคลุมแบบฟอร์มคำขอ
- 001 ต้องเพิ่ม FR แสดงลำดับคำขอผ่าน CR ของตนเอง (ผลของ Q-SH-05) ไม่งั้น T-18 จะไม่มีใครเรียกใช้
