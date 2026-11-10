# Tasks: บันทึกการรับคืนอุปกรณ์ (Return)
Spec: SPEC-RET-006 (spec.md Draft v2) | Plan: plan.md Draft v1 | สร้าง: 2569-11-10 | Owner: 660710737 นายวรรัตน์ หลักฐาน (@660710737)

สรุป
- 15 เฟส (ฐานราก + 13 เฟสตาม AC + รอคำตอบ) รวม 30 task
- รอ Q-xx 4 task (T-27 Q-04, T-28 Q-03, T-29 Q-09, T-30 Q-SH-08) (เลข migration 0006 จองแล้ว)
- เฟสแรกที่ลองใช้งานได้จริง: เฟส 1 (เจ้าหน้าที่กรอกรหัสครุภัณฑ์ บันทึกคืนสภาพปกติ แล้วอุปกรณ์กลับเป็น "ว่าง")

เทคโนโลยีตาม plan.md ข้อ 2: Go (net/http + pgx), PostgreSQL 16, golang-migrate, React + Vite, vitest, testcontainers
(ทีมเลือกเอง ไม่ได้มาจาก spec ยกเว้น PostgreSQL / HTTPS / Docker)
ทีมตัดสิน 2569-11-10: ไม่ทำ CR-003 FR-RET-04, FR-RET-05 อยู่ใน tasks นี้ตาม plan.md

คำสั่งเปิดระบบ (ใช้ในทุกจุดตรวจ):
- ฐานข้อมูล: `docker compose -f deploy/docker-compose.yml up -d db`
- ข้อมูลตัวอย่าง: `docker compose -f deploy/docker-compose.yml exec -T db psql -U salmon salmon < backend/db/seed/seed.sql`
  แล้วตามด้วย `... < backend/db/seed/seed_006_return.sql` (T-03)
- หลังบ้าน: `cd backend && go run ./cmd/server`
- หน้าจอ: `cd frontend && npm run dev`
- ล็อกอินจำลอง: เจ้าหน้าที่ `Authorization: Bearer mock-staff-001` (ภาควิชา 1), นักศึกษา `mock-student-001`
- test: `cd backend && go test ./... -v` และ `cd frontend && npm test`

---

## เฟส 0: ฐานราก

### T-01 สร้าง migration ของการรับคืน
- รองรับ: FR-RET-02, FR-RET-05, FR-RET-06, FR-RET-07, CON-TECH-01, CON-DATA-01, DOM-AUDIT-01
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-04, T-08, T-12, T-14, T-16, T-26
- ไฟล์ที่แตะ: `backend/db/migrations/0006_return.up.sql`, `specs/000-shared/data-model.md` (เพิ่มแถว action `return_list`, `return_lookup`, `return_detail`, `return` ใน DM-07 และคอลัมน์ใหม่ของ DM-05), `backend/internal/platform/migrate_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: migration รันบน DB ว่างได้ ตาราง `loan` มี `return_condition`, `returned_late` พร้อม CHECK, มีตาราง `loan_return_photo` และ `audit_log` รับ action ใหม่ 4 ค่า (test ใน `migrate_test.go` ผ่าน)
- สถานะ: เสร็จ รอทีมตรวจ

### T-02 เพิ่มค่าที่ตั้งได้ของการรับคืน
- รองรับ: FR-RET-02, FR-RET-03, FR-RET-04, FR-RET-05
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของ T-04, T-08, T-12, T-13
- ไฟล์ที่แตะ: `backend/internal/config/config.go`, `frontend/src/config.js`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: ไม่มี
- เสร็จเมื่อ: มี `ReturnPhotoMin`, `ReturnPhotoMax`, `ReturnPhotoMaxBytes`, `ReturnPhotoTypes`, `ReturnRequestMaxBytes`, `ReturnConditionStatus`, `Timezone` ใน config.go และค่าคู่กันใน config.js ทุกตัวมีคอมเมนต์อ้าง ID ตาม plan.md ข้อ 5
- สถานะ: เสร็จ รอทีมตรวจ (เฉพาะ config.go ทีมสั่ง 2569-11-10 ไม่ทำ frontend/src/config.js)

### T-03 สร้างโครง package returns ผูก route และข้อมูลตัวอย่าง
- รองรับ: IF-IDP-01, CON-TECH-02
- ตรวจด้วย: ไม่มี AC ตรง ๆ เป็นงานพื้นฐานของทุกเฟส
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/repository.go`, `backend/internal/platform/http.go` (ผูก route `/api/v1/returns/` ผ่าน auth middleware), `backend/cmd/server/main.go`, `backend/db/seed/seed_006_return.sql` (การยืมที่ยังไม่คืนของภาควิชา 2 และการยืมที่คืนแล้ว 2 รายการของภาควิชา 1 สำหรับ AC-RET-10, AC-RET-11)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01
- เสร็จเมื่อ: เรียก `/api/v1/returns/open-loans` ด้วย token เจ้าหน้าที่ได้ 501 (ยังไม่ทำ) และไม่มี token ได้ 401 ตาม middleware ของ 001
- สถานะ: เสร็จ รอทีมตรวจ

---

## เฟส 1: เจ้าหน้าที่บันทึกคืนสภาพปกติ อุปกรณ์กลับเป็น "ว่าง" (AC-RET-01)

### T-04 ทำ `POST /api/v1/returns` สภาพปกติ พร้อมเก็บรูป
- รองรับ: FR-RET-02, FR-RET-03, FR-RET-05, CON-DATA-01, CON-DATA-03
- ตรวจด้วย: AC-RET-01
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/service.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01, T-02, T-03
- เสร็จเมื่อ: `TestAC_RET_01_NormalReturnBecomesAvailable` ผ่าน (1 transaction: `SELECT ... FOR UPDATE` ปิดการยืมด้วย `now()` เก็บรูป เปลี่ยนสถานะเป็น `available`)
- สถานะ: พร้อมทำ

### T-05 ทำ `GET /api/v1/returns/lookup?assetCode=`
- รองรับ: FR-RET-01, CON-DATA-03
- ตรวจด้วย: AC-RET-01 (ขั้นหาการยืมก่อนบันทึก)
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-03
- เสร็จเมื่อ: `TestReturns_Lookup_FoundOpenLoan` ผ่าน (คืน `{loanId, assetCode, modelName, borrowerId, loanType, dueAt, isOverdue}`)
- สถานะ: พร้อมทำ

### T-06 สร้างหน้า `ReturnPage` กรอกรหัส เลือกสภาพ แนบรูป บันทึก
- รองรับ: FR-RET-01, FR-RET-02, FR-RET-05
- ตรวจด้วย: AC-RET-01
- ไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `api.returns.lookup`, `api.returns.record`), `frontend/src/App.jsx` (แสดงหน้า)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง (spec ไม่มี UI ID)
- ต้องทำหลัง: T-04, T-05
- เสร็จเมื่อ: vitest ด้วย client จำลอง: กรอกรหัส เห็นการยืม เลือก "ปกติ" แนบ 1 รูป กดบันทึก แล้วเห็นผล "บันทึกการรับคืนแล้ว สถานะ: ว่าง"
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 1
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_01_NormalReturnBecomesAvailable`, `TestReturns_Lookup_FoundOpenLoan`, `ReturnPage.test.jsx` และ test เดิมของ 001 ต้องยังผ่าน
- เปิดระบบ: ดูคำสั่งด้านบน ล็อกอินจำลองเป็น `staff-001`
- คนลองเองตาม AC-RET-01:
  Given ข้อมูลตัวอย่าง `CS-NB-0004` สถานะ "ถูกยืม" มีการยืมที่ยังไม่คืน
  When กรอก `CS-NB-0004` เลือกสภาพ "ปกติ" แนบรูป JPG 1 รูป แล้วกดบันทึก
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
- สถานะ: พร้อมทำ

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
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 3
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_02_DamagedReturnBecomesUnavailable` และ test ของเฟส 1-2
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-02:
  Given `CS-NB-0005` สถานะ "ถูกยืม"
  When กรอก `CS-NB-0005` เลือก "ชำรุด" แนบรูป 1 รูป แล้วกดบันทึก
  Then เห็นผล "สถานะ: ไม่พร้อมใช้งาน" และหน้าของ 001 แสดงป้าย "ไม่พร้อมใช้งาน" ไม่มีปุ่มส่งคำขอยืม
- ยังขาด: ช่องรายละเอียดความเสียหาย (T-27 รอ Q-04)

---

## เฟส 4: ต้องระบุสภาพก่อนบันทึก (AC-RET-04)

### T-09 ตรวจสภาพที่ server (ไม่มี / ค่าไม่รู้จัก รวม `lost`)
- รองรับ: FR-RET-02
- ตรวจด้วย: AC-RET-04
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_RET_04_ConditionRequired` ผ่าน (422 `CONDITION_REQUIRED` การยืมไม่เปลี่ยน)
- สถานะ: พร้อมทำ

### T-10 แจ้งให้เลือกสภาพบนหน้าจอ และมีแค่ 2 ตัวเลือก
- รองรับ: FR-RET-02
- ตรวจด้วย: AC-RET-04
- ไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-09
- เสร็จเมื่อ: vitest: ตัวเลือกมีแค่ "ปกติ" กับ "ชำรุด" (ไม่มี "สูญหาย") และกดบันทึกโดยไม่เลือก แล้วเห็นข้อความให้เลือกสภาพ
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 4
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_04_ConditionRequired`, `ReturnPage.test.jsx` และ test ของเฟส 1-3
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-04:
  Given เปิด `ReturnPage` กรอก `CS-NB-0004` แล้ว
  When ไม่เลือกสภาพ แนบรูป แล้วกดบันทึก
  Then เห็นข้อความให้เลือกสภาพ ปกติ หรือ ชำรุด ไม่มีตัวเลือก "สูญหาย" และ `CS-NB-0004` ยัง "ถูกยืม"
- ยังขาด: ไม่มี

---

## เฟส 5: รูปต้องครบเงื่อนไข (AC-RET-12)

### T-11 ตรวจรูปที่ server: จำนวน ชนิดจาก magic bytes ขนาด
- รองรับ: FR-RET-05
- ตรวจด้วย: AC-RET-12
- ไฟล์ที่แตะ: `backend/internal/returns/photos.go`, `backend/internal/returns/photos_test.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_RET_12_PhotoRulesEnforced` ผ่านทั้ง 4 กรณีย่อย (`PHOTO_REQUIRED`, `PHOTO_TOO_MANY`, `PHOTO_TYPE`, `PHOTO_TOO_LARGE`) และ DB ไม่เปลี่ยน ใช้ค่าจาก config ไม่ฝังตัวเลข
- สถานะ: พร้อมทำ

### T-12 จำกัดขนาด request รวมของ `POST /returns`
- รองรับ: FR-RET-05
- ตรวจด้วย: AC-RET-12 (กรณี (ง) ไฟล์ใหญ่)
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-11
- เสร็จเมื่อ: request ใหญ่กว่า `ReturnRequestMaxBytes` ได้ 422 `PHOTO_TOO_LARGE` โดยไม่อ่านเข้าหน่วยความจำทั้งก้อน (`http.MaxBytesReader`) test ผ่าน
- สถานะ: พร้อมทำ

### T-13 แจ้งเหตุผลรูปไม่ผ่านบนหน้าจอ
- รองรับ: FR-RET-05
- ตรวจด้วย: AC-RET-12
- ไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-11
- เสร็จเมื่อ: vitest: ตรวจก่อนส่งด้วยค่าจาก `config.js` และแสดงข้อความตาม error code ทั้ง 4 แบบจาก client จำลอง
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 5
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_12_PhotoRulesEnforced`, test ของ T-12, `ReturnPage.test.jsx` และ test ของเฟส 1-4
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-12:
  Given เปิด `ReturnPage` กรอก `CS-NB-0004` เลือก "ปกติ"
  When ลองทีละกรณี: (ก) ไม่แนบรูป (ข) แนบ 6 รูป (ค) แนบ GIF (ง) แนบ JPG ขนาดเกิน 5 MB แล้วกดบันทึก
  Then ทุกกรณีเห็นเหตุผลที่ไม่ผ่าน และ `CS-NB-0004` ยัง "ถูกยืม"
- ยังขาด: ไม่มี

---

## เฟส 6: เปิดดูรูปที่แนบไว้ภายหลัง (AC-RET-05)

### T-14 ทำ `GET /returns/loans/{loanId}` และ `GET /returns/loans/{loanId}/photos/{photoId}`
- รองรับ: FR-RET-05, FR-RET-06
- ตรวจด้วย: AC-RET-05
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04, T-11
- เสร็จเมื่อ: `TestAC_RET_05_PhotosStoredAndViewable` ผ่าน (1 JPG + 4 PNG บันทึกได้ มี 5 รูปใน detail และ GET รูปได้ bytes กับ Content-Type เดิม)
- สถานะ: พร้อมทำ

### T-15 สร้างหน้า `LoanDetailPage` แสดงเวลาคืน สภาพ และรูป
- รองรับ: FR-RET-05
- ตรวจด้วย: AC-RET-05
- ไฟล์ที่แตะ: `frontend/src/pages/Return/LoanDetailPage.jsx`, `frontend/src/pages/Return/__tests__/LoanDetailPage.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `api.returns.loanDetail`)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-14, T-06
- เสร็จเมื่อ: vitest: แสดงเวลาคืน สภาพ และรูป 5 รูปจาก client จำลอง และ `ReturnPage` มีลิงก์ไปหน้านี้หลังบันทึกสำเร็จ
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 6
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_05_PhotosStoredAndViewable`, `LoanDetailPage.test.jsx` และ test ของเฟส 1-5
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-05:
  Given เปิด `ReturnPage` กรอก `CS-NB-0004`
  When เลือก "ปกติ" แนบ JPG 1 รูป และ PNG 4 รูป (แต่ละรูปไม่เกิน 5 MB) กดบันทึก แล้วกดลิงก์ไปดูรายละเอียดการยืม
  Then เห็นรูปครบ 5 รูป และกดเปิดแต่ละรูปได้
- ยังขาด: ไม่มี

---

## เฟส 7: บันทึกว่าคืนเกินกำหนด (AC-RET-06)

### T-16 คำนวณและบันทึก `returned_late` ตอนรับคืน
- รองรับ: FR-RET-06, DOM-LOAN-01
- ตรวจด้วย: AC-RET-06
- ไฟล์ที่แตะ: `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี (หน้าจอแสดงผลในเฟส 9)
- ต้องทำหลัง: T-04
- เสร็จเมื่อ: `TestAC_RET_06_LateReturnFlagged` ผ่าน (`due_at` ก่อน now 1 ชั่วโมงได้ `true` / หลัง now 1 ชั่วโมงได้ `false` อ่าน `due_at` จากการยืม ไม่คำนวณ 16:00 เอง)
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 7
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_06_LateReturnFlagged` และ test ของเฟส 1-6
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-06:
  Given ข้อมูลตัวอย่าง `CS-NB-0005` กำหนดคืน 16:00 เมื่อวาน และ `CS-NB-0004` กำหนดคืน 16:00 วันนี้ (ลองก่อน 16:00)
  When บันทึกคืนทั้งสองชิ้น แล้วดู `returned_late` ใน DB ด้วย psql
  Then `CS-NB-0005` เป็น `true` และ `CS-NB-0004` เป็น `false`
- ยังขาด: ไม่มี

---

## เฟส 8: ดูรายการที่ยังไม่คืนของภาควิชาตนเอง (AC-RET-10)

### T-17 ทำ `GET /api/v1/returns/open-loans` กรองตามภาควิชา
- รองรับ: FR-RET-01, IF-IDP-01
- ตรวจด้วย: AC-RET-10
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-03
- เสร็จเมื่อ: `TestAC_RET_10_OpenLoansOwnDepartment` ผ่าน (ภาควิชาตน 3 เปิด + 2 คืนแล้ว + ภาควิชาอื่น 1 เปิด ได้ 3 รายการพอดี เรียงตาม `dueAt`)
- สถานะ: พร้อมทำ

### T-18 สร้างหน้า `OpenLoansPage` เลือกแถวแล้วไปบันทึกคืน
- รองรับ: FR-RET-01
- ตรวจด้วย: AC-RET-10
- ไฟล์ที่แตะ: `frontend/src/pages/Return/OpenLoansPage.jsx`, `frontend/src/pages/Return/__tests__/OpenLoansPage.test.jsx`, `frontend/src/api/client.js` (เพิ่ม `api.returns.openLoans`), `frontend/src/pages/Return/ReturnPage.jsx` (รับ `assetCode` ที่ส่งมา)
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-17, T-06
- เสร็จเมื่อ: vitest: แสดง 3 แถวจาก client จำลอง กดแถวแล้วเปิด `ReturnPage` ที่กรอกรหัสครุภัณฑ์ไว้แล้ว
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 8
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_10_OpenLoansOwnDepartment`, `OpenLoansPage.test.jsx` และ test ของเฟส 1-7
- เปิดระบบ: ดูคำสั่งด้านบน (โหลด `seed_006_return.sql` ด้วย)
- คนลองเองตาม AC-RET-10:
  Given รีเซ็ตข้อมูลตัวอย่าง ภาควิชา 1 มีการยืมที่ยังไม่คืน และภาควิชา 2 มี 1 รายการ
  When ล็อกอินเป็น `staff-001` เปิดหน้ารายการที่ยังไม่คืน แล้วกดแถว `CS-NB-0004`
  Then เห็นเฉพาะการยืมที่ยังไม่คืนของภาควิชา 1 (ไม่เห็นรายการที่คืนแล้วหรือของภาควิชา 2) และเปิดหน้าบันทึกคืนที่มีรหัส `CS-NB-0004` แล้ว
- ยังขาด: การแบ่งหน้า (Q-P1 ของ plan ยังไม่แบ่ง)

---

## เฟส 9: เจ้าหน้าที่เห็นเครื่องหมายคืนเกินกำหนด (AC-RET-13)

### T-19 ส่ง `isOverdue` ใน open-loans และ `returnedLate` ใน loan detail
- รองรับ: FR-RET-06
- ตรวจด้วย: AC-RET-13
- ไฟล์ที่แตะ: `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-14, T-16, T-17
- เสร็จเมื่อ: `TestAC_RET_13_OverdueShownToStaff` ผ่าน (การยืมที่คืนเกินกำหนดได้ `returnedLate=true` และการยืมที่ยังไม่คืนแต่เลยกำหนดได้ `isOverdue=true`)
- สถานะ: พร้อมทำ

### T-20 แสดงป้าย "คืนเกินกำหนด" และ "เลยกำหนดคืน" บนหน้าจอ
- รองรับ: FR-RET-06
- ตรวจด้วย: AC-RET-13
- ไฟล์ที่แตะ: `frontend/src/pages/Return/LoanDetailPage.jsx`, `frontend/src/pages/Return/OpenLoansPage.jsx`, `frontend/src/pages/Return/__tests__/LoanDetailPage.test.jsx`, `frontend/src/pages/Return/__tests__/OpenLoansPage.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-19, T-15, T-18
- เสร็จเมื่อ: vitest: `LoanDetailPage` แสดง "คืนเกินกำหนด" เมื่อ `returnedLate`, `OpenLoansPage` แสดง "เลยกำหนดคืน" เมื่อ `isOverdue`
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 9
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_13_OverdueShownToStaff`, `LoanDetailPage.test.jsx`, `OpenLoansPage.test.jsx` และ test ของเฟส 1-8
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-13:
  Given รีเซ็ตข้อมูลตัวอย่าง `CS-NB-0005` เลยกำหนดคืนแล้ว
  When เปิดหน้ารายการที่ยังไม่คืน จากนั้นบันทึกคืน `CS-NB-0005` แล้วเปิดรายละเอียดการยืมของชิ้นนั้น
  Then หน้ารายการแสดงป้าย "เลยกำหนดคืน" ที่ `CS-NB-0005` และหน้ารายละเอียดแสดง "คืนเกินกำหนด"
- ยังขาด: ผู้ยืมเห็นเครื่องหมายของตนเอง (T-28 รอ Q-03)

---

## เฟส 10: อุปกรณ์ที่ไม่มีการยืมค้าง บันทึกคืนไม่ได้ (AC-RET-03)

### T-21 ตอบ 409 `NO_OPEN_LOAN` และกันการบันทึกซ้อน
- รองรับ: FR-RET-01, CON-DATA-01
- ตรวจด้วย: AC-RET-03
- ไฟล์ที่แตะ: `backend/internal/returns/service.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-04, T-05
- เสร็จเมื่อ: `TestAC_RET_03_NoOpenLoanRejected` ผ่าน (อุปกรณ์ `available` และ `unavailable` ได้ 409 ทั้ง lookup และ POST สถานะไม่เปลี่ยน) และ `TestReturns_ConcurrentReturn_OneWins` ผ่าน (2 request พร้อมกันได้ 200 หนึ่งครั้ง 409 หนึ่งครั้ง)
- สถานะ: พร้อมทำ

### T-22 แจ้ง "อุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน" บนหน้าจอ
- รองรับ: FR-RET-01
- ตรวจด้วย: AC-RET-03
- ไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-21
- เสร็จเมื่อ: vitest: client จำลองตอบ `NO_OPEN_LOAN` แล้วหน้าจอแสดงข้อความ และไม่แสดงฟอร์มสภาพ/รูป
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 10
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_03_NoOpenLoanRejected`, `TestReturns_ConcurrentReturn_OneWins`, `ReturnPage.test.jsx` และ test ของเฟส 1-9
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-03:
  Given `CS-NB-0001` สถานะ "ว่าง" และ `CS-PJ-0002` สถานะ "ไม่พร้อมใช้งาน"
  When กรอกรหัสทั้งสองทีละรหัสใน `ReturnPage`
  Then เห็น "อุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน" และสถานะทั้งสองชิ้นไม่เปลี่ยน
- ยังขาด: ไม่มี

---

## เฟส 11: อุปกรณ์ภาควิชาอื่น บันทึกคืนไม่ได้ (AC-RET-11)

### T-23 ตอบ 403 `FORBIDDEN_DEPARTMENT` สำหรับอุปกรณ์ภาควิชาอื่น
- รองรับ: FR-RET-01, IF-IDP-01
- ตรวจด้วย: AC-RET-11
- ไฟล์ที่แตะ: `backend/internal/returns/service.go`, `backend/internal/returns/repository.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-21
- เสร็จเมื่อ: `TestAC_RET_11_OtherDepartmentForbidden` ผ่าน (lookup, POST, loan detail และรูป ได้ 403 การยืมและสถานะไม่เปลี่ยน)
- สถานะ: พร้อมทำ

### T-24 แจ้ง "ไม่มีสิทธิ์" บนหน้าจอ
- รองรับ: FR-RET-01
- ตรวจด้วย: AC-RET-11
- ไฟล์ที่แตะ: `frontend/src/pages/Return/ReturnPage.jsx`, `frontend/src/pages/Return/__tests__/ReturnPage.test.jsx`
- แบบหน้าจอ: ไม่มี mockup ทีมออกแบบเอง
- ต้องทำหลัง: T-23
- เสร็จเมื่อ: vitest: client จำลองตอบ `FORBIDDEN_DEPARTMENT` แล้วหน้าจอแสดงข้อความไม่มีสิทธิ์
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 11
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_11_OtherDepartmentForbidden`, `ReturnPage.test.jsx` และ test ของเฟส 1-10
- เปิดระบบ: ดูคำสั่งด้านบน (โหลด `seed_006_return.sql` ด้วย)
- คนลองเองตาม AC-RET-11:
  Given อุปกรณ์ "ถูกยืม" ของภาควิชา 2 จาก `seed_006_return.sql`
  When ล็อกอินเป็น `staff-001` (ภาควิชา 1) กรอกรหัสครุภัณฑ์ของชิ้นนั้น
  Then เห็นข้อความไม่มีสิทธิ์ และการยืมของชิ้นนั้นยังไม่คืน
- ยังขาด: ไม่มี

---

## เฟส 12: นักศึกษาบันทึกคืนไม่ได้ (AC-RET-08)

### T-25 ตรวจบทบาทเจ้าหน้าที่ทุก endpoint ของ returns
- รองรับ: FR-RET-01, IF-IDP-01
- ตรวจด้วย: AC-RET-08
- ไฟล์ที่แตะ: `backend/internal/returns/handler.go`, `backend/internal/returns/handler_test.go`
- แบบหน้าจอ: ไม่มี (AC-RET-08 เป็นการเรียกโดยตรง หน้าจอของ 006 เป็นของเจ้าหน้าที่เท่านั้น)
- ต้องทำหลัง: T-03, T-17, T-14
- เสร็จเมื่อ: `TestAC_RET_08_StudentForbidden` ผ่าน (token นักศึกษาเรียกทั้ง 5 endpoint ได้ 403 `FORBIDDEN_ROLE` body ไม่มีข้อมูลการยืม DB ไม่เปลี่ยน)
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 12
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_08_StudentForbidden` และ test ของเฟส 1-11
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-08:
  Given หลังบ้านรันอยู่
  When `curl -H "Authorization: Bearer mock-student-001" http://localhost:8000/api/v1/returns/open-loans` และ `curl -X POST ... /api/v1/returns -F assetCode=CS-NB-0004 -F condition=normal`
  Then ได้ 403 `FORBIDDEN_ROLE` ทั้งสองครั้ง และ `CS-NB-0004` ยัง "ถูกยืม"
- ยังขาด: ไม่มี

---

## เฟส 13: ทุกการรับคืนมี audit log (AC-RET-07)

### T-26 เขียน audit ของ returns ทุก endpoint
- รองรับ: FR-RET-07, DOM-AUDIT-01
- ตรวจด้วย: AC-RET-07
- ไฟล์ที่แตะ: `backend/internal/returns/service.go`, `backend/internal/returns/handler.go`, `backend/internal/returns/handler_test.go` (ใช้ตาราง `audit_log` ของ 001 T-31 ผ่าน package audit ถ้ามีแล้ว)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-01, T-21, T-23, T-25
- เสร็จเมื่อ: `TestAC_RET_07_ReturnAudited` ผ่าน (บันทึกสำเร็จมีแถว `action='return'` `outcome='ok'` อยู่ใน transaction เดียวกัน, ถูกปฏิเสธมี `outcome='forbidden'`, การอ่านมี `return_list` / `return_lookup` / `return_detail`)
- สถานะ: พร้อมทำ

#### จุดตรวจเฟส 13
- test ที่ต้องผ่านทั้งหมด: `TestAC_RET_07_ReturnAudited` และ test ทั้งหมดของเฟส 1-12
- เปิดระบบ: ดูคำสั่งด้านบน
- คนลองเองตาม AC-RET-07:
  Given รีเซ็ตข้อมูลตัวอย่าง
  When บันทึกคืน `CS-NB-0004` บนหน้าจอ แล้ว `SELECT * FROM audit_log ORDER BY id DESC LIMIT 5`
  Then มีแถว `return` ที่มี `actor_id = staff-001`, เวลา, `asset_code = CS-NB-0004`, `outcome = ok`
- ยังขาด: สถานะก่อน/หลังใน audit (T-29 รอ Q-09)

---

## เฟส 14: รอคำตอบ

### T-27 เพิ่มรายละเอียดความเสียหายเมื่อสภาพชำรุด
- รองรับ: FR-RET-02, FR-RET-04
- ตรวจด้วย: AC-RET-02 (ต้องเพิ่ม AC เมื่อได้คำตอบ)
- ไฟล์ที่แตะ: ตามคำตอบ Q-04
- แบบหน้าจอ: ไม่มี mockup
- ต้องทำหลัง: T-08
- เสร็จเมื่อ: กำหนดเมื่อได้คำตอบ Q-04
- สถานะ: รอ Q-04

### T-28 ให้ผู้ยืมเห็นเครื่องหมายคืนเกินกำหนดของตนเอง
- รองรับ: FR-RET-06
- ตรวจด้วย: AC-RET-13 (ส่วนของผู้ยืม ต้องเพิ่ม AC เมื่อได้คำตอบ)
- ไฟล์ที่แตะ: ตามคำตอบ Q-03 (อาจเป็นหน้าของฟีเจอร์อื่น)
- แบบหน้าจอ: ไม่มี mockup
- ต้องทำหลัง: T-19
- เสร็จเมื่อ: กำหนดเมื่อได้คำตอบ Q-03
- สถานะ: รอ Q-03

### T-29 บันทึกสถานะก่อน/หลังใน audit ของการเขียน
- รองรับ: FR-RET-07, DOM-AUDIT-01
- ตรวจด้วย: AC-RET-07
- ไฟล์ที่แตะ: ตามคำตอบ Q-09 (migration ใหม่ของ `audit_log` ถ้าขยาย DM-07)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-26
- เสร็จเมื่อ: กำหนดเมื่อได้คำตอบ Q-09
- สถานะ: รอ Q-09

### T-30 จัดการกรณีระบบยืนยันตัวตนเชื่อมต่อไม่ได้
- รองรับ: IF-IDP-01
- ตรวจด้วย: AC-RET-08 (ต้องเพิ่ม AC เมื่อได้คำตอบ)
- ไฟล์ที่แตะ: ตามคำตอบ Q-SH-08 (น่าจะอยู่ใน `backend/internal/auth/` ของกลาง)
- แบบหน้าจอ: ไม่มี
- ต้องทำหลัง: T-25
- เสร็จเมื่อ: กำหนดเมื่อได้คำตอบ Q-SH-08
- สถานะ: รอ Q-SH-08

---

## ตารางตรวจความครบ

### 1. AC กับ task
| AC ID | task ที่ตรวจ AC นี้ |
|---|---|
| AC-RET-01 | T-04, T-05, T-06 |
| AC-RET-02 | T-08 (T-27 รอ Q-04) |
| AC-RET-03 | T-21, T-22 |
| AC-RET-04 | T-09, T-10 |
| AC-RET-05 | T-14, T-15 |
| AC-RET-06 | T-16 |
| AC-RET-07 | T-26 (T-29 รอ Q-09) |
| AC-RET-08 | T-25 (T-30 รอ Q-SH-08) |
| AC-RET-09 | T-07 |
| AC-RET-10 | T-17, T-18 |
| AC-RET-11 | T-23, T-24 |
| AC-RET-12 | T-11, T-12, T-13 |
| AC-RET-13 | T-19, T-20 (T-28 รอ Q-03) |

### 2. Constraint กับ task
| Constraint ID | task ที่ทำให้เป็นจริง |
|---|---|
| CON-TECH-01 | T-01 |
| CON-TECH-02 | T-03 (ใช้ deploy ของ 001) |
| CON-DATA-01 | T-01, T-04, T-08, T-21 |
| CON-DATA-03 | T-04, T-05 |
| DOM-AUDIT-01 | T-01, T-26 (T-29 รอ Q-09) |
| DOM-LOAN-01 | T-16 |
| DOM-REQ-01 | ไม่มี task: plan.md ข้อ 5 ยังไม่ได้ใช้ เพราะการระงับคำขอค้างเป็น FR ใหม่รอ Q-08 (ไม่มี FR รองรับ จึงสร้าง task ไม่ได้) |
| IF-IDP-01 | T-03, T-17, T-23, T-25 (T-30 รอ Q-SH-08) |

### 3. FR กับ AC
| FR ID | AC ที่ตรวจ FR นี้ | เฟส |
|---|---|---|
| FR-RET-01 | AC-RET-03, AC-RET-08, AC-RET-10, AC-RET-11 | 1, 8, 10, 11, 12 |
| FR-RET-02 | AC-RET-04, AC-RET-09 | 1, 2, 4 |
| FR-RET-03 | AC-RET-01 | 1 |
| FR-RET-04 | AC-RET-02 | 3 |
| FR-RET-05 | AC-RET-05, AC-RET-12 | 1, 5, 6 |
| FR-RET-06 | AC-RET-06, AC-RET-13 | 7, 9 |
| FR-RET-07 | AC-RET-07 | 13 |

## สิ่งที่ยังไม่ทำ
Open Questions จาก spec.md (ทุกข้อ) และ task ที่รออยู่
- Q-01 อุปกรณ์สูญหาย (ปิดการยืมแบบ "สูญหาย") เป็น FR ใหม่ที่ต้องผ่าน CR-003 ทีมตัดสินไม่ทำ CR-003 -> ไม่มี task (ไม่มี FR รองรับ)
- Q-02 [ปิดแล้ว] ดู ASM-06 -> ไม่มี task รอ
- Q-03 ผู้ยืมเห็นเครื่องหมายคืนเกินกำหนดของตนเองที่หน้าไหน -> T-28
- Q-04 สภาพ "ชำรุด" ต้องเขียนรายละเอียดความเสียหายไหม -> T-27
- Q-05 [ปิดแล้ว] ดู ASM-01 -> ไม่มี task รอ
- Q-06 [ปิดแล้ว] ดู ASM-03 -> ไม่มี task รอ
- Q-07 spec เขียนว่าต้องทำ CR-003 ก่อน /plan แต่ทีมตัดสินไม่ทำ -> T-08 (FR-RET-04) และ T-11 ถึง T-15 (FR-RET-05) ทำตามคำสั่งทีม
  ทีมยังต้องตัดสินว่าจะแก้ SRS 3.1.7, 4.3 และ specs/README.md อย่างไร
- Q-08 คืนมาชำรุดแต่มีคำขอค้าง (DOM-REQ-01) ระงับคำขอและแจ้งผู้ขอ เป็น FR ใหม่ -> ไม่มี task (ไม่มี FR รองรับ)
- Q-09 ขยาย DOM-AUDIT-01 ให้ครอบการเขียน -> T-29
- Q-SH-08 ระบบยืนยันตัวตนเชื่อมต่อไม่ได้ -> T-30
- Q-P1 ถึง Q-P5 (open points ของ plan.md): Q-P1 การแบ่งหน้าเขียนไว้ในจุดตรวจเฟส 8
