# Implementation Plan: บันทึกการรับคืนอุปกรณ์ (Return)
Spec: SPEC-RET-006 (specs/006-return/spec.md, Draft v2.2) | Plan status: Draft v1.2 | Updated: 2569-11-12 (รวม main: path /v1/returns, migration 0010, CON-TECH-03)
ทีมตัดสิน 2569-11-11: **ไม่ทำ CR-003** และถอดส่วนที่ต้องทำ CR ออก: FR-RET-05 (รูปตอนรับคืน) ยกเลิก, UC-04 ไม่ทำ
FR-RET-04 อยู่ใน plan (data-model.md ให้ 006 ยืนยันลำดับกับ UC-07 เอง) / ปิดการยืมแบบสูญหาย Q-01, ระงับคำขอค้าง Q-08, ขยาย DOM-AUDIT-01 Q-09 **ไม่อยู่ใน plan** (ดูข้อ 8)

## 1. สรุปแนวทาง
- เจ้าหน้าที่ (เท่านั้น) บันทึกการรับคืนอุปกรณ์ของภาควิชาตนเอง: หาการยืมจากรหัสครุภัณฑ์หรือหน้ารายการที่ยังไม่คืน (FR-RET-01)
- บันทึก 1 ครั้ง = 1 transaction: ปิดการยืม (เวลาคืนจากระบบ สภาพ เครื่องหมายคืนเกินกำหนด) + เปลี่ยนสถานะอุปกรณ์ + audit (FR-RET-02, 03, 04, 06, 07)
- สภาพกำหนดสถานะหลังคืน: ปกติเป็น `available`, ชำรุดเป็น `unavailable` ตรงไปเลยไม่ผ่าน `available` (FR-RET-03, FR-RET-04)
- ทุกกฎตรวจที่ฝั่ง server (บทบาท ภาควิชา สภาพ เวลา) หน้าจอตรวจซ้ำเพื่อความสะดวกเท่านั้น
- เขียน DB ทันทีไม่มี cache/queue เพื่อให้ 001 เห็นผลภายใน 30 วินาที (NFR-DATA-01)

## 2. เทคโนโลยีที่ใช้
| สิ่งที่เลือก | มาจาก | หมายเหตุ |
|---|---|---|
| PostgreSQL 16 | CON-TECH-01 | ใช้ DB และตารางเดียวกับ 001 (`loan`, `equipment_item`, `audit_log`) |
| HTTPS ผ่าน reverse proxy, deploy ด้วย Docker Compose | CON-TECH-02, NFR-SEC-01 | ใช้ `deploy/docker-compose.yml` ที่มีอยู่ ไม่ตั้งค่าใหม่ |
| Go (net/http + pgx v5) | ทีมเลือกเอง ไม่ได้มาจาก spec (AGENTS.md, plan.md ข้อ 2 ของ 001) | package `backend/internal/returns/` |
| golang-migrate (SQL ล้วน) | ทีมเลือกเอง ไม่ได้มาจาก spec | migration ใหม่ ห้ามแก้ 0001, 0002 |
| React + Vite, vitest | ทีมเลือกเอง ไม่ได้มาจาก spec (template รายวิชา) | `frontend/src/pages/Return/` |
| go test + testcontainers | ทีมเลือกเอง ไม่ได้มาจาก spec | มีใน go.mod แล้ว (ใช้กับ 001) |
| Auth จำลอง `Authorization: Bearer mock-<รหัสผู้ใช้>` | ASM-02 (รอ Q-SH-03) | ใช้ `auth.Principal{UserID, Role, DepartmentID}` ของ 001 T-03 |

## 3. โมเดลข้อมูล
ใช้ตารางร่วมที่มีแล้ว (`backend/db/migrations/0001`, `0002`) และเพิ่มด้วย migration ใหม่ 1 ไฟล์ (`0010_return.up.sql` ต่อจาก 0009 ของ 005 บันทึกใน data-model.md)

**`loan` (DM-05 ใช้ร่วม) เพิ่มคอลัมน์ ว่างได้ทั้งหมด เพื่อไม่ให้กระทบ 001 และ 005:**
| คอลัมน์ | ชนิด | รองรับ |
|---|---|---|
| `returned_at` (มีแล้ว) | `TIMESTAMPTZ NULL` | FR-RET-02 เวลาคืน = `now()` ของ DB ตอนบันทึก (ASM-04) |
| `return_condition` | `TEXT NULL CHECK (return_condition IN ('normal','damaged'))` | FR-RET-02 สภาพ 2 ค่า ไม่มี `lost` |
| `returned_late` | `BOOLEAN NULL` | FR-RET-06 = `returned_at > due_at` คำนวณและบันทึกตอนคืน ไม่คำนวณใหม่ภายหลัง |
| CHECK เพิ่ม | `(returned_at IS NULL) = (return_condition IS NULL) AND (returned_at IS NULL) = (returned_late IS NULL)` | คืนแล้วต้องมีครบทั้ง 3 ค่า |

**`equipment_item` (DM-04 ใช้ร่วม):** ไม่เพิ่มคอลัมน์ แก้เฉพาะ `status` และ `updated_at` (FR-RET-03, FR-RET-04, CON-DATA-01)
ไม่มีคอลัมน์ "เหตุ" (ส่งซ่อม/สูญหาย) เพราะ spec ไม่มี FR รองรับ (Q-01, Q-08)

**`audit_log` (DM-07 ใช้ร่วม):** เพิ่มค่าใน CHECK ของ `action` ด้วย migration เดียวกัน และเพิ่มแถวในตาราง action ของ data-model.md
| action | ใช้ตอน |
|---|---|
| `return_list` | เปิดหน้ารายการที่ยังไม่คืน |
| `return_lookup` | ค้นการยืมจากรหัสครุภัณฑ์ (`query = {"assetCode": ...}`) |
| `return_detail` | เปิดรายละเอียดการยืม |
| `return` | บันทึกการรับคืน (`asset_code`, `outcome`) |

ใช้แค่คอลัมน์ที่ DM-07 มีอยู่ ไม่กำหนดฟิลด์ audit เพิ่มเอง (FR-RET-07 ตามคำตอบทีม การขยายไปการเขียนรอ Q-09)

ไม่มีฟิลด์ข้อมูลส่วนตัวใหม่: ไม่เก็บเลขบัตร ชื่อ หรือข้อมูลผู้ยืมเพิ่มจาก `borrower_id` ที่มีอยู่

## 4. API / หน้าจอ
ทุก endpoint: HTTPS, ต้องผ่าน auth middleware, `Role == staff` เท่านั้น (นักศึกษา: 403 `FORBIDDEN_ROLE` AC-RET-08)
และกรองด้วย `equipment_model.department_id = principal.DepartmentID` ที่ชั้น repository ทุก query (ASM-03)
ทุก response ตั้ง `Cache-Control: no-store`

| Method / path | รองรับ | input | output หลัก |
|---|---|---|---|
| `GET /v1/returns/open-loans` | FR-RET-01(ข), FR-RET-06 | - | `[{loanId, assetCode, modelName, borrowerId, loanType, dueAt, isOverdue}]` เรียงตาม `dueAt` |
| `GET /v1/returns/lookup?assetCode=` | FR-RET-01(ก) | รหัสครุภัณฑ์ | การยืมที่ยังไม่คืน 1 รายการ (รูปแบบเดียวกับแถวด้านบน) |
| `POST /v1/returns` | FR-RET-02, 03, 04, 06, 07 | form (urlencoded หรือ multipart): `assetCode`, `condition` (`normal`/`damaged`) | `{loanId, assetCode, returnedAt, condition, returnedLate, itemStatus}` |
| `GET /v1/returns/loans/{loanId}` | FR-RET-06 | - | `{loanId, assetCode, borrowerId, dueAt, returnedAt, condition, returnedLate}` |

`isOverdue` ของรายการที่ยังไม่คืน = `now() > due_at` คำนวณที่ server (FR-RET-06 ส่วน "เห็นการยืมที่เลยกำหนดคืนแล้ว")
ไม่รับเวลาจาก client: ถ้า form มีฟิลด์เวลาใดๆ ให้ไม่สนใจ (AC-RET-09)

**ลำดับใน `POST /v1/returns` (1 transaction):**
1. ตรวจ `condition` ก่อนแตะ DB
2. `SELECT ... FROM loan JOIN equipment_item ... WHERE asset_code = $1 AND returned_at IS NULL FOR UPDATE`
   ไม่พบ: 409 `NO_OPEN_LOAN` (AC-RET-03) / ต่างภาควิชา: 403 `FORBIDDEN_DEPARTMENT` (AC-RET-11) / สถานะไม่ใช่ `borrowed`: 409 `NO_OPEN_LOAN`
3. `UPDATE loan SET returned_at = now(), return_condition = $c, returned_late = (now() > due_at)`
4. `UPDATE equipment_item SET status = <ตาม config>, updated_at = now()` ครั้งเดียว (ไม่ผ่าน `available` FR-RET-04)
5. `INSERT audit_log (action='return', outcome='ok')`
6. commit ถ้าขั้นไหนพลาดให้ rollback ทั้งหมด และเขียน audit `outcome='error'`/`'forbidden'` นอก transaction
บันทึกพร้อมกัน 2 คน: `FOR UPDATE` ทำให้คนที่สองเห็น `returned_at` แล้ว ได้ 409 `NO_OPEN_LOAN`

**Error contract** (ทุก error ตอบ `{"code", "message"}` ตาม CON-TECH-03 ข้อความอยู่ใน `config.ReturnMessages`):
| กรณี | HTTP | code | AC |
|---|---|---|---|
| ไม่มี/หมดอายุ token | 401 | `UNAUTHENTICATED` (ของ 001) | - |
| นักศึกษาเรียก | 403 | `FORBIDDEN_ROLE` | AC-RET-08 |
| อุปกรณ์ภาควิชาอื่น | 403 | `FORBIDDEN_DEPARTMENT` | AC-RET-11 |
| ไม่มีการยืมที่ยังไม่คืน | 409 | `NO_OPEN_LOAN` | AC-RET-03 |
| ไม่ระบุสภาพ / ค่าไม่รู้จัก (รวม `lost`) | 422 | `CONDITION_REQUIRED` | AC-RET-04 |
| ไม่พบการยืม (รายละเอียด) | 404 | `LOAN_NOT_FOUND` (เพิ่มใน T-11 ลงทะเบียนแล้ว Q-10) | - |
| DB ล่ม | 503 | `UPSTREAM_UNAVAILABLE` (ของ 001) | - |

**หน้าจอ** (spec ไม่มีหมวด UI และไม่มี UI ID ไม่มี mockup ทีมออกแบบเอง):
| หน้าจอ | รองรับ | หมายเหตุ |
|---|---|---|
| `ReturnPage` (`frontend/src/pages/Return/ReturnPage.jsx`) ไม่มี mockup ทีมออกแบบเอง | FR-RET-01(ก), 02 | ช่องรหัสครุภัณฑ์ แสดงการยืมที่พบ เลือกสภาพ ปกติ/ชำรุด (ไม่มีตัวเลือกสูญหาย) ปุ่มบันทึก แสดงข้อความ error ตาม code |
| `OpenLoansPage` (`frontend/src/pages/Return/OpenLoansPage.jsx`) ไม่มี mockup ทีมออกแบบเอง | FR-RET-01(ข), 06 | ตารางการยืมที่ยังไม่คืน ป้าย "เลยกำหนดคืน" เลือกแถวแล้วไป `ReturnPage` พร้อม `assetCode` |
| `LoanDetailPage` (`frontend/src/pages/Return/LoanDetailPage.jsx`) ไม่มี mockup ทีมออกแบบเอง | FR-RET-06 | แสดงเวลาคืน สภาพ และเครื่องหมาย "คืนเกินกำหนด" |

## 5. ตารางตรวจ Constraints
| Constraint ID | ถูกนำไปใช้ที่ไหนใน plan | สถานะ |
|---|---|---|
| CON-TECH-01 | ข้อ 2, 3: ใช้ PostgreSQL ตารางร่วม + migration ใหม่ | ใช้แล้ว |
| CON-TECH-02 | ข้อ 2: ผ่าน reverse proxy HTTPS และ Docker Compose ที่มีอยู่ | ใช้แล้ว |
| CON-TECH-03 | ข้อ 4: error ทุกตัว `{code, message}` (`errorBody` ใน handler.go) | ใช้แล้ว (ลงทะเบียน NO_OPEN_LOAN, CONDITION_REQUIRED, LOAN_NOT_FOUND ใน constraints.md v5.4 ปิด Q-10) |
| CON-DATA-01 | ข้อ 3, 4: เขียนได้เฉพาะ `available` / `unavailable` จาก `borrowed` CHECK เดิมของ 0001 กันค่าอื่น | ใช้แล้ว |
| CON-DATA-03 | ข้อ 4: หาการยืมด้วย `asset_code` (UNIQUE) ใน lookup และ POST | ใช้แล้ว |
| DOM-AUDIT-01 | ข้อ 3: action `return_list`, `return_lookup`, `return_detail` (อ่าน) และ `return` (เขียน) ด้วยคอลัมน์ DM-07 เดิม | ใช้แล้ว (ส่วนขยายไปการเขียนตาม Q-09 ไม่มี CR ใช้เท่าที่ DM-07 มี) |
| DOM-LOAN-01 | ข้อ 3, 4: `returned_late` / `isOverdue` เทียบกับ `due_at` ที่ 005 บันทึก (16:00 สำหรับยืมปกติ) ไม่คำนวณ 16:00 เอง (ASM-01) | ใช้แล้ว |
| DOM-REQ-01 | ไม่มี 006 ไม่อ่านหรือแก้ `borrow_request` | ยังไม่ได้ใช้ เพราะการระงับคำขอค้างตอนคืนชำรุดเป็น FR ใหม่ที่รอ Q-08 |
| IF-IDP-01 | ข้อ 4: `auth.Principal` ให้ `Role` (ตรวจ staff) และ `DepartmentID` (กรองภาควิชา) | ใช้แล้ว (กรณีระบบยืนยันตัวตนล่ม รอ Q-SH-08) |

**ค่าที่ตั้งได้** (อยู่ใน `backend/internal/config/config.go` และ `frontend/src/config.js` เท่านั้น พร้อมคอมเมนต์อ้าง ID):
| ค่า | มาจาก ID | ชื่อตัวแปรและไฟล์ |
|---|---|---|
| สภาพ 2 ค่า `normal`, `damaged` และสถานะหลังคืน `normal` เป็น `available`, `damaged` เป็น `unavailable` | FR-RET-02, FR-RET-03, FR-RET-04 | `ReturnConditionStatus` (map) ใน config.go / `RETURN_CONDITIONS` ใน config.js |
| เขตเวลาแสดงผล Asia/Bangkok | glossary.md (วันที่/เวลาปัจจุบัน) | `Timezone` ใน config.go (ถ้า 001 เพิ่มก่อนให้ใช้ตัวเดียวกัน) |

## 6. แผนทดสอบจาก Acceptance Criteria
ไฟล์: `backend/internal/returns/handler_test.go` (integration ด้วย testcontainers PostgreSQL + migration จริง) และ `frontend/src/pages/Return/__tests__/` (vitest)
เวลาในการทดสอบ: ใช้ `now()` ของ DB ผ่านฟังก์ชันเดียวใน repository ที่ test override ได้ (เช่นตั้ง `due_at` สัมพัทธ์กับเวลาจริง)

| AC ID | ชื่อ test | ทดสอบอย่างไร |
|---|---|---|
| AC-RET-01 | `TestAC_RET_01_NormalReturnBecomesAvailable` | seed loan เปิด + item `borrowed` -> POST `normal` + 1 JPG -> 200, `returned_at` ไม่ว่าง, item `available` |
| AC-RET-02 | `TestAC_RET_02_DamagedReturnBecomesUnavailable` | POST `damaged` -> item `unavailable`; ตรวจว่ามี UPDATE สถานะครั้งเดียวด้วย trigger นับ/เก็บประวัติสถานะใน test schema ว่าไม่เคยเป็น `available` |
| AC-RET-03 | `TestAC_RET_03_NoOpenLoanRejected` | item `available` และ `unavailable` -> lookup/POST ได้ 409 `NO_OPEN_LOAN`, สถานะไม่เปลี่ยน |
| AC-RET-04 | `TestAC_RET_04_ConditionRequired` | ไม่ส่ง `condition` และส่ง `lost` -> 422 `CONDITION_REQUIRED`, loan ไม่เปลี่ยน; frontend test: ตัวเลือกมีแค่ ปกติ/ชำรุด |
| AC-RET-05 | - | ยกเลิก 2569-11-11 พร้อม FR-RET-05 |
| AC-RET-06 | `TestAC_RET_06_LateReturnFlagged` | `due_at` = 1 ชั่วโมงก่อน now -> `returned_late = true`; `due_at` = 1 ชั่วโมงหลัง now -> `false` (แทนเวลา 17:00/15:00 ของ AC ด้วยเวลาสัมพัทธ์) |
| AC-RET-07 | `TestAC_RET_07_ReturnAudited` | หลัง POST สำเร็จมีแถว `audit_log` `action='return'`, `actor_id`, `actor_role='staff'`, `asset_code`, `outcome='ok'` |
| AC-RET-08 | `TestAC_RET_08_StudentForbidden` | token นักศึกษาเรียกทุก endpoint -> 403 `FORBIDDEN_ROLE`, body ไม่มีข้อมูลการยืม, DB ไม่เปลี่ยน |
| AC-RET-09 | `TestAC_RET_09_ReturnTimeFromServer` | ส่งฟิลด์ `returnedAt=2569-10-10T09:00` มาด้วย -> `returned_at` ใกล้ `now()` ของ DB (ต่างไม่เกินไม่กี่วินาที) ไม่ใช่ค่าที่ส่ง |
| AC-RET-10 | `TestAC_RET_10_OpenLoansOwnDepartment` | seed ภาควิชาตน 3 เปิด + 2 คืนแล้ว, ภาควิชาอื่น 1 เปิด -> GET open-loans ได้ 3 รายการพอดี |
| AC-RET-11 | `TestAC_RET_11_OtherDepartmentForbidden` | lookup/POST อุปกรณ์ภาควิชาอื่น -> 403 `FORBIDDEN_DEPARTMENT`, DB ไม่เปลี่ยน |
| AC-RET-12 | - | ยกเลิก 2569-11-11 พร้อม FR-RET-05 |
| AC-RET-13 | `TestAC_RET_13_OverdueShownToStaff` | loan คืนเกินกำหนด -> GET detail `returnedLate=true`; loan เปิดเลยกำหนด -> open-loans `isOverdue=true`; frontend test แสดงป้าย |

ทดสอบเพิ่ม: บันทึกพร้อมกัน 2 request ได้ 200 หนึ่งครั้งและ 409 หนึ่งครั้ง
NFR-DATA-01 ตรวจโดย AC-EQS-05 ของ 001 (ไม่ทำ test ซ้ำใน 006)

## 7. ลำดับงาน
1. Migration `0010_return.up.sql` (คอลัมน์ `loan`, ค่า action ใหม่ของ `audit_log`) + แถวใหม่ใน data-model.md (FR-RET-02, 06, 07)
2. ค่าใน `config.go` / `config.js` ตามตารางข้อ 5 (FR-RET-02 ถึง 04)
3. Repository + สิทธิ์ (staff + ภาควิชา) และ `GET /returns/lookup` (FR-RET-01(ก), AC-RET-03, 08, 11)
4. `POST /returns` สภาพปกติ: transaction ปิดการยืม + `available` + audit (FR-RET-02, 03, 07, AC-RET-01, 04, 07, 09)
5. สภาพชำรุด เป็น `unavailable` (FR-RET-04, AC-RET-02)
7. คืนเกินกำหนด + `GET /returns/open-loans` + `GET /returns/loans/{id}` (FR-RET-01(ข), 06, AC-RET-06, 10, 13)
8. หน้าจอ `ReturnPage`, `OpenLoansPage`, `LoanDetailPage` + vitest (FR-RET-01, 02, 06)
9. ผูก route เข้ากับ router ใน `main.go` / `platform` และทดสอบรวมกับ AC-EQS-05 ของ 001 (NFR-DATA-01)

## 8. สิ่งที่ยังไม่ทำ
Open Questions จาก spec (ทุกข้อ) **ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ**
- Q-01 อุปกรณ์สูญหาย: ทีมเสนอให้เจ้าหน้าที่ปิดการยืมแบบ "สูญหาย" ต้องระบุหมายเหตุ ระบบบันทึกเวลาปิด และสถานะเป็น "ไม่พร้อมใช้งาน" (เหตุ: สูญหาย)
  เป็น FR ใหม่ ต้องทำ CR -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ทีมตัดสินไม่ทำใน prototype)
- Q-02 [ไม่ใช้แล้ว 2569-11-11 เพราะ FR-RET-05 ยกเลิก]
- Q-03 ผู้ยืมเห็นเครื่องหมายคืนเกินกำหนดของตนเองที่หน้าไหน -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (plan มีแค่หน้าของเจ้าหน้าที่)
- Q-04 สภาพ "ชำรุด" ต้องเขียนรายละเอียดความเสียหายไหม -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ไม่มีช่องหมายเหตุ)
- Q-05 [ปิดแล้ว] ดู ASM-01
- Q-06 [ปิดแล้ว] ดู ASM-03
- Q-07 [ปิดแล้ว 2569-11-11] ไม่ทำ CR-003 ถอด FR-RET-05 และ UC-04 ออก plan ตรงกับ SRS 3.1.7, 4.3 และ specs/README.md แล้ว
- Q-08 คืนมาชำรุดแต่มีคำขอค้าง (DOM-REQ-01): ระงับคำขอและแจ้งผู้ขอ -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (006 ไม่แตะ `borrow_request`)
- Q-09 ขยาย DOM-AUDIT-01 ให้ครอบการเขียน -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ
  (plan ใช้คอลัมน์ DM-07 เดิม ไม่เพิ่มฟิลด์ เช่น สถานะก่อน/หลัง)
- Q-SH-08 ระบบยืนยันตัวตนเชื่อมต่อไม่ได้ -> ส่วนที่เกี่ยวข้องกับข้อนี้จะยังไม่สร้างจนกว่าจะได้คำตอบ (ใช้พฤติกรรมของ middleware 001 ตามเดิม)

Open points ของ plan (ไม่ใช่ข้อกำหนด ทีมยืนยัน):
- Q-P1 หน้ารายการที่ยังไม่คืนต้องแบ่งหน้าไหม spec ไม่บอก plan คืนทั้งหมดของภาควิชา (ขนาดข้อมูล prototype เล็ก)
- Q-P2 [ไม่ใช้แล้ว 2569-11-11 เพราะ FR-RET-05 ยกเลิก]
- Q-P3 [ไม่ใช้แล้ว 2569-11-11 เพราะ FR-RET-05 ยกเลิก]
- Q-P4 [ปิดแล้ว 2569-11-10, เปลี่ยนเลข 2569-11-11] เดิมใช้ `0006_return.up.sql` แต่ 0006 และ 0007 เป็นของ 005 ส่วน 0004 และ 0005 เป็นของ 003, 0008 เป็นของ 004 จึงย้ายไป `0009_return.up.sql` / 2569-11-12 ตอนรวม main พบว่า 0009 เป็นของ 005 แล้ว ย้ายเป็น `0010_return.up.sql`
- Q-P5 หน้าจอแสดง `borrowerId` แทนชื่อ เพราะยังไม่มีที่มาของชื่อผู้ยืม (ประเด็นเดียวกับ 001 จาก T-02)
