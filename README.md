# SalmomNorway: ระบบยืม-คืนอุปกรณ์ภาควิชา

Term Project รายวิชา 520461-165 Software Requirement Specification and Management (1/2569)
ภาควิชาคอมพิวเตอร์ คณะวิทยาศาสตร์ มหาวิทยาลัยศิลปากร

## ทีม

- ชื่อทีม: SalmonNorway
- เครื่องมือ AI ที่ใช้: Claude Code

| สมาชิก | GitHub | ฟีเจอร์ที่รับผิดชอบ | ส่วนอื่นที่ทำ |
|---|---|---|---|
| 660710767 นางสาวพรรณธร จันทร์สายทอง | https://github.com/Pannathorn05 | 001-equipment-status | วางของกลาง (000-shared, AGENTS.md), ร่าง SRS, โครง Go + ฐานข้อมูล + auth จำลอง, audit_log และแพ็กเกจ audit ใช้ร่วม, รหัส error กลาง (ITEM_NOT_FOUND, METHOD_NOT_ALLOWED และตัวครอบ 405), หน้าล็อกอินจำลอง, แก้ Migrate ให้ใช้บน Windows ได้ |
| 660710715 นางสาวธัญลักษณ์ เดชนครชัย  | https://github.com/tnyl02 | 002-borrow-request | สร้าง repo ของทีม, CR-003 ตรวจคำกำกวมที่ปิด Q-SH-04 (กำหนดคืนยืมค้างคืน 16:00) และ Q-SH-05 (ลำดับคำขอแสดงที่ 001) ให้ทุกฟีเจอร์, CR-004 ยกรูปแบบ error เป็นมาตรฐานกลาง CON-TECH-03 ที่ทุกฟีเจอร์ใช้, CR-007 ใช้ REQUEST_NOT_FOUND และ FORBIDDEN_ROLE ร่วม, endpoint ลำดับคำขอ `GET /api/v1/equipment/items/{assetCode}/requests` ให้หน้ารายละเอียดของ 001 (CR-005), ปรับข้อมูลตั้งต้น seed.sql ให้ลองจุดตรวจของ 001 และ 002 ได้, `postForm` (ส่งไฟล์แนบ) และ `ApiError.body` ใน api/client.js |
| 660710763 นางสาวนภัสนัน ดีเสมอ | https://github.com/napatsanan-np | 003-cancel-request | ฟังก์ชัน postJSON ใน api/client.js ที่ 002 และ 004 ใช้ส่งข้อมูลด้วย, migration 0004 (staff_contact) และ 0005 (เพิ่ม action notify ใน audit_log ที่ 004 ใช้ร่วม), CR-006 ลงทะเบียนรหัส error ALREADY_APPROVED, NOT_CANCELLABLE และใช้ REQUEST_NOT_FOUND ร่วมกับ 004 ในตาราง CON-TECH-03, load test k6 (AC-CXL-17) ผ่าน Docker |
| 660710758 นางสาวณัฐสุดา กัญญาบุญ | https://github.com/natsuda05 | 004-approve-request | HTTPS ของหลังบ้าน (TLS 1.2 ขึ้นไป, CON-TECH-02, NFR-SEC-01) และ deploy หลังบ้านด้วย Docker (Dockerfile, service backend ใน docker-compose, โปรแกรมสร้าง cert สำหรับเครื่องพัฒนา cmd/devcert) |
| 660710764 นายบรรลุ ฉิมสิงห์ | https://github.com/mmoonn101219 | 005-handover | ลงทะเบียนรหัส error กลาง FORBIDDEN_ROLE (001, 002, 006 ใช้ต่อ), REQUEST_NOT_APPROVED และ BORROWER_MISMATCH ใน CON-TECH-03, migration 0009 เพิ่ม action `handover` ใน audit_log |
| 660710737 นายวรรัตน์ หลักฐาน | https://github.com/660710737 | 006-return | ลงทะเบียนรหัส error NO_OPEN_LOAN, CONDITION_REQUIRED, LOAN_NOT_FOUND ใน CON-TECH-03, เพิ่มฟิลด์รับคืนของ loan และ action รับคืนของ audit_log ใน data-model.md (migration 0010_return), ตัดสิน Q-SH-09 ของกลาง (คืนแล้วชำรุดเปลี่ยนเป็น "ไม่พร้อมใช้งาน" ที่ UC-09 ตาม FR-RET-04) |

ดูว่าใครทำอะไรได้จาก 3 ที่: ตารางนี้, ช่อง `Owner:` ใน spec.md ของแต่ละฟีเจอร์, และ `git log --author=<ชื่อ>`

## สิ่งที่ส่ง

| งาน | อยู่ที่ |
|---|---|
| SRS | `docs/srs/SRS.md` (diagram อยู่ใน Miro ลิงก์ในภาคผนวกของ SRS) |
| spec.md | `specs/00N-<feature>/spec.md` ทั้ง 6 ฟีเจอร์ |
| Prototype | `backend/` (Go) และ `frontend/` (React) |
| Traceability Matrix | `specs/00N-<feature>/rtm.md` (ทำหลังทุกฟีเจอร์เสร็จ) |

## โครงของ repo

```
README.md                 ไฟล์นี้
AGENTS.md                 กติกาที่ AI ต้องทำตาม (รวมตารางเทียบ Go กับคำสั่งของรายวิชา)
CLAUDE.md                 ชี้ไป AGENTS.md (สำหรับ Claude Code)
prompt-log.md             บันทึกการสั่ง AI ทุกครั้ง พร้อมชื่อคนสั่ง
docs/srs/SRS.md           SRS ฉบับเต็ม
specs/README.md           ดัชนีฟีเจอร์และผู้รับผิดชอบ
specs/000-shared/         Constraint ข้อมูล และคำศัพท์ที่ทุกฟีเจอร์ใช้ร่วมกัน
specs/00N-<feature>/      spec.md plan.md tasks.md test-cases.md rtm.md ของแต่ละฟีเจอร์
specs/backlog.md          backlog จัดลำดับแล้ว และบันทึกการเจรจา N-01
specs/changes/            change request (CR-001 ถึง CR-009)
specs/CHANGELOG.md        บันทึกว่า spec แต่ละเวอร์ชันเปลี่ยนอะไร
backend/                  หลังบ้าน Go + PostgreSQL
frontend/                 หน้าจอ React + Vite
.claude/commands/         คำสั่ง /clarify /plan /tasks /implement /testcases /verify (Claude Code)
.github/prompts/          คำสั่งชุดเดียวกัน (Copilot)
.cursor/commands/         คำสั่งชุดเดียวกัน (Cursor)
```
