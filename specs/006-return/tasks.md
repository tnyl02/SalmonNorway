# Tasks: บันทึกการรับคืนอุปกรณ์ (Return)
Spec: SPEC-RET-006 (spec.md Draft v2.1) | Plan: plan.md Draft v1.1 | สร้าง: 2569-11-10 | แก้: 2569-11-11 (ลบ task รูป เลื่อนเลข) | Owner: 660710737 นายวรรัตน์ หลักฐาน (@660710737)

สรุป
- 14 เฟส (ฐานราก + 12 เฟส + รอคำตอบ) รวม 27 task
- รอ Q-xx 4 task (T-24 Q-04, T-25 Q-03, T-26 Q-09, T-27 Q-SH-08) (เลข migration 0006 จองแล้ว)
- เฟสแรกที่ลองใช้งานได้จริง: เฟส 1 (เจ้าหน้าที่กรอกรหัสครุภัณฑ์ บันทึกคืนสภาพปกติ แล้วอุปกรณ์กลับเป็น "ว่าง")

เลข task และเฟสถูกเลื่อน 2569-11-11 (ทีมสั่ง): ลบ task รูปเดิม T-11 ถึง T-13 และเฟส 5 เดิม (FR-RET-05 ยกเลิก)
prompt-log และ commit ก่อนวันที่นี้ใช้เลขเก่า เทียบได้ดังนี้: T-14 -> T-11, T-15 -> T-12, ... , T-30 -> T-27 (เลขใหม่ = เลขเก่า - 3)
เฟส 6 -> 5, ... , เฟส 14 -> 13 (เลขใหม่ = เลขเก่า - 1) / T-01 ถึง T-10 และเฟส 0 ถึง 4 ไม่เปลี่ยน

เทคโนโลยีตาม plan.md ข้อ 2: Go (net/http + pgx), PostgreSQL 16, golang-migrate, React + Vite, vitest, testcontainers
(ทีมเลือกเอง ไม่ได้มาจาก spec ยกเว้น PostgreSQL / HTTPS / Docker)
ทีมตัดสิน 2569-11-11: ไม่ทำ CR-003 ถอด FR-RET-05 (รูป) และ UC-04 ออก / FR-RET-04 อยู่ใน tasks นี้ตาม plan.md v1.1

คำสั่งเปิดระบบ (ใช้ในทุกจุดตรวจ):
- ฐานข้อมูล: `docker compose -f deploy/docker-compose.yml up -d db`
- ข้อมูลตัวอย่าง: `docker compose -f deploy/docker-compose.yml exec -T db psql -U salmon salmon < backend/db/seed/seed.sql`
  แล้วตามด้วย `... < backend/db/seed/seed_006_return.sql` (T-03)
- หลังบ้าน: `cd backend && go run ./cmd/server`
- หน้าจอ: `cd frontend && npm run dev`
- ล็อกอินจำลอง: เจ้าหน้าที่ `Authorization: Bearer mock-staff-001` (ภาควิชา 1), นักศึกษา `mock-student-001`
- test: `cd backend && go test ./... -v` และ `cd frontend && npm test`

## ลำดับความสำคัญ (ทีมขอ 2569-11-11 AI ร่าง ทีมตรวจ)
task ส่วนใหญ่ผูกกับ AC ใน spec ถ้าข้าม task ไหน RTM จาก /verify จะแสดงว่า AC นั้นยังไม่ผ่าน "ควรทำ" จึงไม่ได้แปลว่าตัดทิ้งได้

| ความสำคัญ | task | เกณฑ์ |
|---|---|---|
| ต้องทำ | T-01 ถึง T-06, T-08, T-09 (เสร็จ), T-18, T-20, T-22, T-23 | ไม่ทำแล้วระบบทำงานผิดหรือไม่ปลอดภัย: ทางหลักของการบันทึกคืน, บังคับสภาพ (FR-RET-02), กันบันทึกซ้อน, กันภาควิชาอื่น (ASM-03), กันนักศึกษา (AC-RET-08), audit (DOM-AUDIT-01) |
| ควรทำ | T-07, T-10 (เสร็จ), T-11, T-12, T-13, T-14, T-15, T-16, T-17, T-19, T-21 | มี AC รองรับ แต่ระบบหลักใช้ได้แม้ยังไม่ทำ: test ของโค้ดที่ทำแล้ว (T-07, T-13), ข้อความ error บนหน้าจอ, รายละเอียดการยืม, หน้ารายการที่ยังไม่คืน, ป้ายคืนเกินกำหนด |
| ไม่ต้องทำตอนนี้ | T-24, T-25, T-26, T-27 | ติด Open Question (Q-04, Q-03, Q-09, Q-SH-08) ทำไม่ได้จนกว่าจะได้คำตอบ |

ลำดับที่แนะนำ: T-05 -> T-06 (จบเฟส 1) -> T-22 -> T-20 -> T-18 -> T-09 -> T-08 -> T-23 -> T-07, T-13 -> ที่เหลือของ "ควรทำ"
หมายเหตุ: ช่อง "ต้องทำหลัง" ของ T-22 ระบุ T-14, T-11 (กลุ่ม "ควรทำ") และ T-20 ระบุ T-18 ถ้าจะทำ T-22 ก่อน ต้องให้ทีมปรับช่องนั้นก่อน

---

## เฟส 0: ฐานราก

### T-01 สร้าง migration ของการรับคืน
- รองรับ: FR-RET-02, FR-RET-06, FR-RET-07, CON-TECH-01, CON-DATA-01, DOM-AUDIT-01
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-04, T-08, T-11, T-13, T-23
- ไฟล์ที่แตะ: `backend/db/migrations/0006_return.up.sql`, `specs/000-shared/data-model.md` (เพิ่มแถว action `return_list`, `return_lookup`, `return_detail`, `return` ใน DM-07 และคอลัมน์ใหม่ของ DM-05), `backend/internal/platform/migrate_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: migration รันบน DB ว่างได้ ตาราง `loan` มี `return_condition`, `returned_late` พร้อม CHECK และ `audit_log` รับ action ใหม่ 4 ค่า (ตาราง `loan_return_photo` ถอดออก 2569-11-11 พร้อม FR-RET-05) (test ใน `migrate_test.go` ผ่าน)
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ต้องทำ

### T-02 เพิ่มค่าที่ตั้งได้ของการรับคืน
- รองรับ: FR-RET-02, FR-RET-03, FR-RET-04
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-04, T-08
- ไฟล์ที่แตะ: `backend/internal/config/config.go`, `frontend/src/config.js`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: มี `Timezone` และ `ReturnConditionStatus` ใน config.go และ `TIMEZONE`, `STATUS_LABELS`, `RETURN_CONDITIONS` ใน config.js ทุกตัวมีคอมเมนต์อ้าง ID (ค่ารูปถอดออก 2569-11-11 พร้อม FR-RET-05)
- สถานะ: เสร็จ รอทีมตรวจ (config.go 2569-11-10, config.js เพิ่มตามทีมสั่ง 2569-11-11)
- ความสำคัญ: ต้องทำ

### T-03 สร้างโครง package returns ผูก route และข้อมูลตัวอย่าง
- รองรับ: IF-IDP-01, CON-TECH-02
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของทุกเฟส
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/repository.go`, `backend/internal/platform/http.go` (ผูก route `/api/v1/returns/` ผ่าน auth middleware), `backend/cmd/server/main.go`, `backend/db/seed/seed_006_return.sql` (การยืมที่ยังไม่คืนของภาควิชา 2 และการยืมที่คืนแล้ว 2 รายการของภาควิชา 1 สำหรับ AC-RET-10, AC-RET-11)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: เรียก `/api/v1/returns/open-loans` ด้วย token เจ้าหน้าที่ได้ 501 (ยังไม่ทำ) และไม่มี token ได้ 401 ตาม middleware ของ 001
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ต้องทำ

---

## เฟส 1: เจ้าหน้าที่บันทึกคืนสภาพปกติ อุปกรณ์กลับเป็น "ว่าง" (AC-RET-01)

### T-04 ทำ `POST /api/v1/returns` สภาพปกติ
- รองรับ: FR-RET-02, FR-RET-03, CON-DATA-01, CON-DATA-03
- ตรวจด้วย: AC-RET-01
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/service.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01, T-02, T-03
- เสร็จเมื่อ: `TestAC_RET_01_NormalReturnBecomesAvailable` ผ่าน (1 transaction: `SELECT ... FOR UPDATE` ปิดการยืมด้วย `now()` เปลี่ยนสถานะเป็น `available`)
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ต้องทำ

### T-05 ทำ `GET /api/v1/returns/lookup?assetCode=`
- รองรับ: FR-RET-01, CON-DATA-03
- ตรวจด้วย: AC-RET-01 (ขั้นหาการยืมก่อนบันทึก)
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-03
- เสร็จเมื่อ: `TestReturns_Lookup_FoundOpenLoan` ผ่าน (คืน `{loanId, assetCode, modelName, borrowerId, loanType, dueAt, isOverdue}`)
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ต้องทำ

### T-06 สร้างหน้า `ReturnPage` กรอกรหัส เลือกสภาพ บันทึก
- รองรับ: FR-RET-01, FR-RET-02
- ตรวจด้วย: AC-RET-01
- ไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `api.returns.lookup`, `api.returns.record`), `frontend/src/App.jsx` (แสดงหน้า)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (spec ไม่มี UI ID)
- ต้องทำหลัง: T-04, T-05
- เสร็จเมื่อ: vitest ด้วย client จำลอง: กรอกรหัส เห็นการยืม เลือก "ปกติ" กดบันทึก แล้วเห็นผล "บันทึกการรับคืนแล้ว สถานะ: ว่าง"
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ต้องทำ

#### จุดตรวจเฟส 1
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_01_NormalReturnBecomesAvailable`, `TestReturns_Lookup_FoundOpenLoan`, `ReturnPage.test.jsx` และ test เดิมของ 001 ต้องยังผ่าน
- เปิดระบบ: ดูคำสั่งด้านบน ล็อกอินจำลองเป็น `staff-001` โดยพิมพ์ `localStorage.setItem('token', 'mock-staff-001')` ใน console ของเบราว์เซอร์ (จนกว่า 001 T-18 จะมีหน้าล็อกอิน)
- คนลองเองตาม AC-RET-01:
  Given ข้อมูลตัวอย่าง `CS-NB-0004` สถานะ "ถูกยืม" มีการยืมที่ยังไม่คืน
  When กรอก `CS-NB-0004` เลือกสภาพ "ปกติ" แล้วกดบันทึก
  Then เห็นข้อความบันทึกสำเร็จ และเปิดหน้ารายการอุปกรณ์ของ 001 แล้ว `CS-NB-0004` เป็น "ว่าง" (ภายใน 30 วินาที NFR-DATA-01)
- ยังขาด: ไม่มี

---

## เฟส 2: เวลาคืนเป็นเวลาของระบบ (AC-RET-09)

### T-07 บังคับให้เวลาคืนมาจากระบบ และทดสอบว่าไม่รับเวลาจาก client
- รองรับ: FR-RET-02
- ตรวจด้วย: AC-RET-09
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_RET_09_ReturnTimeFromServer` ผ่าน (ส่ง `returnedAt` มาแล้ว `returned_at` ยังเป็น `now()` ของ DB)
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ควรทำ

#### จุดตรวจเฟส 2
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_09_ReturnTimeFromServer` และ test ของเฟส 1
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-09:
  Given รีเซ็ตข้อมูลตัวอย่าง และดูเวลาบนเครื่อง
  When บันทึกคืน `CS-NB-0004` บนหน้าจอ (หน้าจอไม่มีช่องให้กรอกเวลา) แล้วดู `returned_at` ใน DB ด้วย psql
  Then `returned_at` ตรงกับเวลาที่กดบันทึก (Asia/Bangkok) และหน้าจอไม่มีช่องแก้เวลา
- ยังขาด: ไม่มี

---

## เฟส 3: คืนมาชำรุด เปลี่ยนเป็น "ไม่พร้อมใช้งาน" (AC-RET-02)

### T-08 รองรับสภาพ "ชำรุด" ไปเป็น `unavailable` โดยตรง
- รองรับ: FR-RET-04, CON-DATA-01
- ตรวจด้วย: AC-RET-02
- ไฟล์ที่แตะ: `backend/internal/returns/service.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี (`ReturnPage` มีตัวเลือก "ชำรุด" แล้วจาก T-06)
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_RET_02_DamagedReturnBecomesUnavailable` ผ่าน (trigger ใน test เก็บประวัติสถานะ แล้วพบว่าไม่เคยเป็น `available`)
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ต้องทำ

#### จุดตรวจเฟส 3
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_02_DamagedReturnBecomesUnavailable` และ test ของเฟส 1-2
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-02:
  Given `CS-NB-0005` สถานะ "ถูกยืม"
  When กรอก `CS-NB-0005` เลือก "ชำรุด" แล้วกดบันทึก
  Then เห็นผล "สถานะ: ไม่พร้อมใช้งาน" และหน้าของ 001 แสดงป้าย "ไม่พร้อมใช้งาน" ไม่มีปุ่มส่งคำขอยืม
- ยังขาด: ช่องรายละเอียดความเสียหาย (T-24 รอ Q-04)

---

## เฟส 4: ต้องระบุสภาพก่อนบันทึก (AC-RET-04)

### T-09 ตรวจสภาพที่ server (ไม่มี / ค่าไม่รู้จัก รวม `lost`)
- รองรับ: FR-RET-02
- ตรวจด้วย: AC-RET-04
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_RET_04_ConditionRequired` ผ่าน (422 `CONDITION_REQUIRED` การยืมไม่เปลี่ยน)
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ต้องทำ

### T-10 แจ้งให้เลือกสภาพบนหน้าจอ และมีแค่ 2 ตัวเลือก
- รองรับ: FR-RET-02
- ตรวจด้วย: AC-RET-04
- ไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-09
- เสร็จเมื่อ: vitest: ตัวเลือกมีแค่ "ปกติ" กับ "ชำรุด" (ไม่มี "สูญหาย") และกดบันทึกโดยไม่เลือก แล้วเห็นข้อความให้เลือกสภาพ
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ควรทำ

#### จุดตรวจเฟส 4
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_04_ConditionRequired`, `ReturnPage.test.jsx` และ test ของเฟส 1-3
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-04:
  Given เปิด `ReturnPage` กรอก `CS-NB-0004` แล้ว
  When ไม่เลือกสภาพ แล้วกดบันทึก
  Then เห็นข้อความให้เลือกสภาพ ปกติ หรือ ชำรุด ไม่มีตัวเลือก "สูญหาย" และ `CS-NB-0004` ยัง "ถูกยืม"
- ยังขาด: ไม่มี

---

## เฟส 5: ดูรายละเอียดการยืมที่คืนแล้ว (ไม่มี AC เป็นพื้นฐานของ AC-RET-13 เดิมคือ AC-RET-05 ที่ยกเลิก)

### T-11 ทำ `GET /returns/loans/{loanId}`
- รองรับ: FR-RET-06
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-16 (AC-RET-13)
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestReturns_LoanDetail` ผ่าน (การยืมที่คืนแล้วได้ `{loanId, assetCode, borrowerId, dueAt, returnedAt, condition, returnedLate}`)
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ควรทำ

### T-12 สร้างหน้า `LoanDetailPage` แสดงเวลาคืน และสภาพ
- รองรับ: FR-RET-06
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-17 (AC-RET-13)
- ไฟล์ที่แตะ: `frontend/src/pages/Return/LoanDetailPage.jsx`, `frontend/src/pages/Return/__tests__/LoanDetailPage.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `api.returns.loanDetail`)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-11, T-06
- เสร็จเมื่อ: vitest: แสดงเวลาคืน และสภาพจาก client จำลอง และ `ReturnPage` มีลิงก์ไปหน้านี้หลังบันทึกสำเร็จ
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ควรทำ

#### จุดตรวจเฟส 5
- test ที่ต้องผ่านทั้งหมด: `TestReturns_LoanDetail`, `LoanDetailPage.test.jsx` และ test ของเฟส 1-4
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเอง (ไม่มี AC):
  Given เปิด `ReturnPage` กรอก `CS-NB-0004`
  When เลือก "ปกติ" กดบันทึก แล้วกดลิงก์ไปดูรายละเอียดการยืม
  Then เห็นเวลาคืนและสภาพ "ปกติ"
- ยังขาด: ไม่มี

---

## เฟส 6: บันทึกว่าคืนเกินกำหนด (AC-RET-06)

### T-13 คำนวณและบันทึก `returned_late` ตอนรับคืน
- รองรับ: FR-RET-06, DOM-LOAN-01
- ตรวจด้วย: AC-RET-06
- ไฟล์ที่แตะ: `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี (หน้าจอแสดงผลในเฟส 8)
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_RET_06_LateReturnFlagged` ผ่าน (`due_at` ก่อน now 1 ชั่วโมงได้ `true` / หลัง now 1 ชั่วโมงได้ `false` อ่าน `due_at` จากการยืม ไม่คำนวณ 16:00 เอง)
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ควรทำ

#### จุดตรวจเฟส 6
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_06_LateReturnFlagged` และ test ของเฟส 1-5
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-06:
  Given ข้อมูลตัวอย่าง `CS-NB-0005` กำหนดคืน 16:00 เมื่อวาน และ `CS-NB-0004` กำหนดคืน 16:00 วันนี้ (ลองก่อน 16:00)
  When บันทึกคืนทั้งสองชิ้น แล้วดู `returned_late` ใน DB ด้วย psql
  Then `CS-NB-0005` เป็น `true` และ `CS-NB-0004` เป็น `false`
- ยังขาด: ไม่มี

---

## เฟส 7: ดูรายการที่ยังไม่คืนของภาควิชาตนเอง (AC-RET-10)

### T-14 ทำ `GET /api/v1/returns/open-loans` กรองตามภาควิชา
- รองรับ: FR-RET-01, IF-IDP-01
- ตรวจด้วย: AC-RET-10
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-03
- เสร็จเมื่อ: `TestAC_RET_10_OpenLoansOwnDepartment` ผ่าน (ภาควิชาตน 3 เปิด + 2 คืนแล้ว + ภาควิชาอื่น 1 เปิด ได้ 3 รายการพอดี เรียงตาม `dueAt`)
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ควรทำ

### T-15 สร้างหน้า `OpenLoansPage` เลือกแถวแล้วไปบันทึกคืน
- รองรับ: FR-RET-01
- ตรวจด้วย: AC-RET-10
- ไฟล์ที่แตะ: `frontend/src/pages/Return/OpenLoansPage.jsx`, `frontend/src/pages/Return/__tests__/OpenLoansPage.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `api.returns.openLoans`), `frontend/src/pages/Return/ReturnPage.jsx` (รับ `assetCode` ที่ส่งมา)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-14, T-06
- เสร็จเมื่อ: vitest: แสดง 3 แถวจาก client จำลอง กดแถวแล้วเปิด `ReturnPage` ที่กรอกรหัสครุภัณฑ์ไว้แล้ว
- สถานะ: เสร็จ รอทีมตรวจ
- ความสำคัญ: ควรทำ

#### จุดตรวจเฟส 7
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_10_OpenLoansOwnDepartment`, `OpenLoansPage.test.jsx` และ test ของเฟส 1-6
- เปิดระบบ: ดูคำสั่งด้านบน (โหลด `seed_006_return.sql` ด้วย)
- คนลองเองตาม AC-RET-10:
  Given รีเซ็ตข้อมูลตัวอย่าง ภาควิชา 1 มีการยืมที่ยังไม่คืน และภาควิชา 2 มี 1 รายการ
  When ล็อกอินเป็น `staff-001` เปิดหน้ารายการที่ยังไม่คืน แล้วกดแถว `CS-NB-0004`
  Then เห็นเฉพาะการยืมที่ยังไม่คืนของภาควิชา 1 (ไม่เห็นรายการที่คืนแล้วหรือของภาควิชา 2) และเปิดหน้าบันทึกคืนที่มีรหัส `CS-NB-0004` แล้ว
- ยังขาด: การแบ่งหน้า (Q-P1 ของ plan ยังไม่แบ่ง)

---

## เฟส 8: เจ้าหน้าที่เห็นเครื่องหมายคืนเกินกำหนด (AC-RET-13)

### T-16 ส่ง `isOverdue` ใน open-loans และ `returnedLate` ใน loan detail
- รองรับ: FR-RET-06
- ตรวจด้วย: AC-RET-13
- ไฟล์ที่แตะ: `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-11, T-13, T-14
- เสร็จเมื่อ: `TestAC_RET_13_OverdueShownToStaff` ผ่าน (การยืมที่คืนเกินกำหนดได้ `returnedLate=true` และการยืมที่ยังไม่คืนแต่เลยกำหนดได้ `isOverdue=true`)
- สถานะ: พร้อมทำ
- ความสำคัญ: ควรทำ

### T-17 แสดงป้าย "คืนเกินกำหนด" และ "เลยกำหนดคืน" บนหน้าจอ
- รองรับ: FR-RET-06
- ตรวจด้วย: AC-RET-13
- ไฟล์ที่แตะ: `frontend/src/pages/Return/LoanDetailPage.jsx`, `frontend/src/pages/Return/OpenLoansPage.jsx`, `frontend/src/pages/Return/__tests__/LoanDetailPage.test.jsx`, `frontend/src/pages/Return/__tests__/OpenLoansPage.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-16, T-12, T-15
- เสร็จเมื่อ: vitest: `LoanDetailPage` แสดง "คืนเกินกำหนด" เมื่อ `returnedLate`, `OpenLoansPage` แสดง "เลยกำหนดคืน" เมื่อ `isOverdue`
- สถานะ: พร้อมทำ
- ความสำคัญ: ควรทำ

#### จุดตรวจเฟส 8
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_13_OverdueShownToStaff`, `LoanDetailPage.test.jsx`, `OpenLoansPage.test.jsx` และ test ของเฟส 1-7
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-13:
  Given รีเซ็ตข้อมูลตัวอย่าง `CS-NB-0005` เลยกำหนดคืนแล้ว
  When เปิดหน้ารายการที่ยังไม่คืน จากนั้นบันทึกคืน `CS-NB-0005` แล้วเปิดรายละเอียดการยืมของชิ้นนั้น
  Then หน้ารายการแสดงป้าย "เลยกำหนดคืน" ที่ `CS-NB-0005` และหน้ารายละเอียดแสดง "คืนเกินกำหนด"
- ยังขาด: ผู้ยืมเห็นเครื่องหมายของตนเอง (T-25 รอ Q-03)

---

## เฟส 9: อุปกรณ์ที่ไม่มีการยืมค้าง บันทึกคืนไม่ได้ (AC-RET-03)

### T-18 ตอบ 409 `NO_OPEN_LOAN` และกันการบันทึกซ้อน
- รองรับ: FR-RET-01, CON-DATA-01
- ตรวจด้วย: AC-RET-03
- ไฟล์ที่แตะ: `backend/internal/returns/service.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04, T-05
- เสร็จเมื่อ: `TestAC_RET_03_NoOpenLoanRejected` ผ่าน (อุปกรณ์ `available` และ `unavailable` ได้ 409 ทั้ง lookup และ POST สถานะไม่เปลี่ยน) และ `TestReturns_ConcurrentReturn_OneWins` ผ่าน (2 request พร้อมกันได้ 200 หนึ่งครั้ง 409 หนึ่งครั้ง)
- สถานะ: พร้อมทำ
- ความสำคัญ: ต้องทำ

### T-19 แจ้ง "อุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน" บนหน้าจอ
- รองรับ: FR-RET-01
- ตรวจด้วย: AC-RET-03
- ไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-18
- เสร็จเมื่อ: vitest: client จำลองตอบ `NO_OPEN_LOAN` แล้วหน้าจอแสดงข้อความ และไม่แสดงฟอร์มสภาพ
- สถานะ: พร้อมทำ
- ความสำคัญ: ควรทำ

#### จุดตรวจเฟส 9
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_03_NoOpenLoanRejected`, `TestReturns_ConcurrentReturn_OneWins`, `ReturnPage.test.jsx` และ test ของเฟส 1-8
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-03:
  Given `CS-NB-0001` สถานะ "ว่าง" และ `CS-PJ-0002` สถานะ "ไม่พร้อมใช้งาน"
  When กรอกรหัสทั้งสองทีละรหัสใน `ReturnPage`
  Then เห็น "อุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน" และสถานะทั้งสองชิ้นไม่เปลี่ยน
- ยังขาด: ไม่มี

---

## เฟส 10: อุปกรณ์ภาควิชาอื่น บันทึกคืนไม่ได้ (AC-RET-11)

### T-20 ตอบ 403 `FORBIDDEN_DEPARTMENT` สำหรับอุปกรณ์ภาควิชาอื่น
- รองรับ: FR-RET-01, IF-IDP-01
- ตรวจด้วย: AC-RET-11
- ไฟล์ที่แตะ: `backend/internal/returns/service.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-18
- เสร็จเมื่อ: `TestAC_RET_11_OtherDepartmentForbidden` ผ่าน (lookup, POST และ loan detail ได้ 403 การยืมและสถานะไม่เปลี่ยน)
- สถานะ: พร้อมทำ
- ความสำคัญ: ต้องทำ

### T-21 แจ้ง "ไม่มีสิทธิ์" บนหน้าจอ
- รองรับ: FR-RET-01
- ตรวจด้วย: AC-RET-11
- ไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-20
- เสร็จเมื่อ: vitest: client จำลองตอบ `FORBIDDEN_DEPARTMENT` แล้วหน้าจอแสดงข้อความไม่มีสิทธิ์
- สถานะ: พร้อมทำ
- ความสำคัญ: ควรทำ

#### จุดตรวจเฟส 10
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_11_OtherDepartmentForbidden`, `ReturnPage.test.jsx` และ test ของเฟส 1-9
- เปิดระบบ: ดูคำสั่งด้านบน (โหลด `seed_006_return.sql` ด้วย)
- คนลองเองตาม AC-RET-11:
  Given อุปกรณ์ "ถูกยืม" ของภาควิชา 2 จาก `seed_006_return.sql`
  When ล็อกอินเป็น `staff-001` (ภาควิชา 1) กรอกรหัสครุภัณฑ์ของชิ้นนั้น
  Then เห็นข้อความไม่มีสิทธิ์ และการยืมของชิ้นนั้นยังไม่คืน
- ยังขาด: ไม่มี

---

## เฟส 11: นักศึกษาบันทึกคืนไม่ได้ (AC-RET-08)

### T-22 ตรวจบทบาทเจ้าหน้าที่ทุก endpoint ของ returns
- รองรับ: FR-RET-01, IF-IDP-01
- ตรวจด้วย: AC-RET-08
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี (AC-RET-08 เป็นการเรียกโดยตรง หน้าจอของ 006 เป็นของเจ้าหน้าที่เท่านั้น)
- ต้องทำหลัง: T-03, T-14, T-11
- เสร็จเมื่อ: `TestAC_RET_08_StudentForbidden` ผ่าน (token นักศึกษาเรียกทั้ง 4 endpoint ได้ 403 `FORBIDDEN_ROLE` body ไม่มีข้อมูลการยืม DB ไม่เปลี่ยน)
- สถานะ: พร้อมทำ
- ความสำคัญ: ต้องทำ

#### จุดตรวจเฟส 11
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_08_StudentForbidden` และ test ของเฟส 1-10
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-08:
  Given หลังบ้านรันอยู่
  When `curl -H "Authorization: Bearer mock-student-001" http://localhost:8000/api/v1/returns/open-loans` และ `curl -X POST ... /api/v1/returns -F assetCode=CS-NB-0004 -F condition=normal`
  Then ได้ 403 `FORBIDDEN_ROLE` ทั้งสองครั้ง และ `CS-NB-0004` ยัง "ถูกยืม"
- ยังขาด: ไม่มี

---

## เฟส 12: ทุกการรับคืนมี audit log (AC-RET-07)

### T-23 เขียน audit ของ returns ทุก endpoint
- รองรับ: FR-RET-07, DOM-AUDIT-01
- ตรวจด้วย: AC-RET-07
- ไฟล์ที่แตะ: `backend/internal/returns/service.go`, `backend/internal/returns/handler.go`, `backend/internal/returns/handler_test.go` (ใช้ตาราง `audit_log` ของ 001 T-31 ผ่าน package audit ถ้ามีแล้ว)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01, T-18, T-20, T-22
- เสร็จเมื่อ: `TestAC_RET_07_ReturnAudited` ผ่าน (บันทึกสำเร็จมีแถว `action='return'` `outcome='ok'` อยู่ใน transaction เดียวกัน, ถูกปฏิเสธมี `outcome='forbidden'`, การอ่านมี `return_list` / `return_lookup` / `return_detail`)
- สถานะ: พร้อมทำ
- ความสำคัญ: ต้องทำ

#### จุดตรวจเฟส 12
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_07_ReturnAudited` และ test ทั้งหมดของเฟส 1-11
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-07:
  Given รีเซ็ตข้อมูลตัวอย่าง
  When บันทึกคืน `CS-NB-0004` บนหน้าจอ แล้ว `SELECT * FROM audit_log ORDER BY id DESC LIMIT 5`
  Then มีแถว `return` ที่มี `actor_id = staff-001`, เวลา, `asset_code = CS-NB-0004`, `outcome = ok`
- ยังขาด: สถานะก่อน/หลังใน audit (T-26 รอ Q-09)

---

## เฟส 13: รอคำตอบ

### T-24 เพิ่มรายละเอียดความเสียหายเมื่อสภาพชำรุด
- รองรับ: FR-RET-02, FR-RET-04
- ตรวจด้วย: AC-RET-02 (ต้องเพิ่ม AC เมื่อได้คำตอบ)
- ไฟล์ที่แตะ: ตามคำตอบ Q-04
- แบบหน้าจอ: ไม่มี mockup
- ต้องทำหลัง: T-08
- เสร็จเมื่อ: กำหนดเมื่อได้คำตอบ Q-04
- สถานะ: รอ Q-04
- ความสำคัญ: ไม่ต้องทำตอนนี้ (รอคำตอบ)

### T-25 ให้ผู้ยืมเห็นเครื่องหมายคืนเกินกำหนดของตนเอง
- รองรับ: FR-RET-06
- ตรวจด้วย: AC-RET-13 (ส่วนของผู้ยืม ต้องเพิ่ม AC เมื่อได้คำตอบ)
- ไฟล์ที่แตะ: ตามคำตอบ Q-03 (อาจเป็นหน้าของฟีเจอร์อื่น)
- แบบหน้าจอ: ไม่มี mockup
- ต้องทำหลัง: T-16
- เสร็จเมื่อ: กำหนดเมื่อได้คำตอบ Q-03
- สถานะ: รอ Q-03
- ความสำคัญ: ไม่ต้องทำตอนนี้ (รอคำตอบ)

### T-26 บันทึกสถานะก่อน/หลังใน audit ของการเขียน
- รองรับ: FR-RET-07, DOM-AUDIT-01
- ตรวจด้วย: AC-RET-07
- ไฟล์ที่แตะ: ตามคำตอบ Q-09 (migration ใหม่ของ `audit_log` ถ้าขยาย DM-07)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-23
- เสร็จเมื่อ: กำหนดเมื่อได้คำตอบ Q-09
- สถานะ: รอ Q-09
- ความสำคัญ: ไม่ต้องทำตอนนี้ (รอคำตอบ)

### T-27 จัดการกรณีระบบยืนยันตัวตนเชื่อมต่อไม่ได้
- รองรับ: IF-IDP-01
- ตรวจด้วย: AC-RET-08 (ต้องเพิ่ม AC เมื่อได้คำตอบ)
- ไฟล์ที่แตะ: ตามคำตอบ Q-SH-08 (น่าจะอยู่ใน `backend/internal/auth/` ของกลาง)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-22
- เสร็จเมื่อ: กำหนดเมื่อได้คำตอบ Q-SH-08
- สถานะ: รอ Q-SH-08
- ความสำคัญ: ไม่ต้องทำตอนนี้ (รอคำตอบ)

---

## ตารางตรวจความครบ

### 1. AC กับ task
| AC ID | task ที่ตรวจ AC นี้ |
|---|---|
| AC-RET-01 | T-04, T-05, T-06 |
| AC-RET-02 | T-08 (T-24 รอ Q-04) |
| AC-RET-03 | T-18, T-19 |
| AC-RET-04 | T-09, T-10 |
| AC-RET-05 | ยกเลิก 2569-11-11 พร้อม FR-RET-05 |
| AC-RET-06 | T-13 |
| AC-RET-07 | T-23 (T-26 รอ Q-09) |
| AC-RET-08 | T-22 (T-27 รอ Q-SH-08) |
| AC-RET-09 | T-07 |
| AC-RET-10 | T-14, T-15 |
| AC-RET-11 | T-20, T-21 |
| AC-RET-12 | ยกเลิก 2569-11-11 พร้อม FR-RET-05 |
| AC-RET-13 | T-16, T-17 (T-25 รอ Q-03) |

### 2. Constraint กับ task
| Constraint ID | task ที่ทำให้เป็นจริง |
|---|---|
| CON-TECH-01 | T-01 |
| CON-TECH-02 | T-03 (ใช้ deploy ของ 001) |
| CON-DATA-01 | T-01, T-04, T-08, T-18 |
| CON-DATA-03 | T-04, T-05 |
| DOM-AUDIT-01 | T-01, T-23 (T-26 รอ Q-09) |
| DOM-LOAN-01 | T-13 |
| DOM-REQ-01 | ไม่มี task: plan.md ข้อ 5 ยังไม่ได้ใช้ เพราะการระงับคำขอค้างเป็น FR ใหม่รอ Q-08 (ไม่มี FR รองรับ จึงสร้าง task ไม่ได้) |
| IF-IDP-01 | T-03, T-14, T-20, T-22 (T-27 รอ Q-SH-08) |

### 3. FR กับ AC
| FR ID | AC ที่ตรวจ FR นี้ | เฟส |
|---|---|---|
| FR-RET-01 | AC-RET-03, AC-RET-08, AC-RET-10, AC-RET-11 | 1, 7, 9, 10, 11 |
| FR-RET-02 | AC-RET-04, AC-RET-09 | 1, 2, 4 |
| FR-RET-03 | AC-RET-01 | 1 |
| FR-RET-04 | AC-RET-02 | 3 |
| FR-RET-05 | ยกเลิก 2569-11-11 | - |
| FR-RET-06 | AC-RET-06, AC-RET-13 | 5, 6, 8 |
| FR-RET-07 | AC-RET-07 | 12 |

## สิ่งที่ยังไม่ทำ
Open Questions จาก spec.md (ทุกข้อ) และ task ที่รออยู่
- Q-01 อุปกรณ์สูญหาย (ปิดการยืมแบบ "สูญหาย") เป็น FR ใหม่ที่ต้องทำ CR ทีมตัดสินไม่ทำใน prototype -> ไม่มี task (ไม่มี FR รองรับ)
- Q-02 [ไม่ใช้แล้ว 2569-11-11 เพราะ FR-RET-05 ยกเลิก] -> ไม่มี task รอ
- Q-03 ผู้ยืมเห็นเครื่องหมายคืนเกินกำหนดของตนเองที่หน้าไหน -> T-25
- Q-04 สภาพ "ชำรุด" ต้องเขียนรายละเอียดความเสียหายไหม -> T-24
- Q-05 [ปิดแล้ว] ดู ASM-01 -> ไม่มี task รอ
- Q-06 [ปิดแล้ว] ดู ASM-03 -> ไม่มี task รอ
- Q-07 [ปิดแล้ว 2569-11-11] ไม่ทำ CR-003 ถอด FR-RET-05 (task รูปเดิม T-11 ถึง T-13 ลบออก, T-11 T-12 เหลือส่วนรายละเอียดการยืม) และ UC-04 ออก / T-08 (FR-RET-04) คงไว้
- Q-08 คืนมาชำรุดแต่มีคำขอค้าง (DOM-REQ-01) ระงับคำขอและแจ้งผู้ขอ เป็น FR ใหม่ -> ไม่มี task (ไม่มี FR รองรับ)
- Q-09 ขยาย DOM-AUDIT-01 ให้ครอบการเขียน -> T-26
- Q-SH-08 ระบบยืนยันตัวตนเชื่อมต่อไม่ได้ -> T-27
- Q-P1 ถึง Q-P5 (open points ของ plan.md): Q-P1 การแบ่งหน้าเขียนไว้ในจุดตรวจเฟส 7
