// Package approve คือฟีเจอร์ 004 อนุมัติคำขอยืมปกติ (UC-12, specs/004-approve-request/)
//
// เจ้าหน้าที่ตรวจและอนุมัติหรือปฏิเสธคำขอยืมแบบปกติของภาควิชาตนเอง (FR-APV-01) และแจ้งผลผู้ขอทาง Email (FR-APV-02)
// ฟีเจอร์นี้เปลี่ยนแค่สถานะของคำขอ **ไม่เปลี่ยนสถานะอุปกรณ์** (CON-DATA-01, data-model.md แถว UC-12)
// กฎทุกข้อบังคับที่ฝั่ง server ไม่พึ่งการซ่อนปุ่มที่หน้าจอ (plan.md ข้อ 1)
package approve

import "time"

// PendingRequest คือคำขอยืมแบบปกติที่รอพิจารณา 1 แถวในหน้ารายการของเจ้าหน้าที่ (FR-APV-01 ก ข)
type PendingRequest struct {
	ID            int       // รหัสคำขอ (DM-06)
	RequesterName string    // ชื่อผู้ขอ (FR-APV-01 ก) ผู้ขอพิมพ์เองตอนส่งคำขอ (FR-BRQ-03 ก)
	StudentCode   string    // รหัสนักศึกษาที่ผู้ขอพิมพ์ (ASM-03)
	ModelName     string    // อุปกรณ์: ชื่อรุ่น (FR-APV-01 ก)
	AssetCode     string    // อุปกรณ์: รหัสครุภัณฑ์ (CON-DATA-03)
	Purpose       string    // วัตถุประสงค์ 1 ใน 3 ค่า (FR-APV-01 ก, FR-BRQ-02)
	PurposeOther  string    // ข้อความที่ผู้ขอพิมพ์เมื่อเลือก "อื่น ๆ" ว่างได้
	StartAt       time.Time // ช่วงเวลาที่ต้องการยืม: เริ่ม (FR-APV-01 ก)
	EndAt         time.Time // ช่วงเวลาที่ต้องการยืม: คืน (FR-APV-01 ก)
	SubmittedAt   time.Time // เวลาที่ส่ง ใช้เรียงลำดับ (FR-APV-01 ก, DOM-REQ-01)
	IsFirst       bool      // เป็นคำขอแรกของอุปกรณ์รายชิ้นนั้นหรือไม่ อนุมัติได้เฉพาะคำขอแรก (FR-APV-01 ข, Q-APV-23)
}
