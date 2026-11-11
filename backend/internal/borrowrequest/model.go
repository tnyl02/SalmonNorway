// Package borrowrequest คือฟีเจอร์ 002 ส่งคำขอยืมอุปกรณ์ (UC-02 รวม UC-05 ยืมค้างคืน)
//
// ฟีเจอร์นี้เขียนเฉพาะตาราง borrow_request และ **ไม่เปลี่ยนสถานะอุปกรณ์เลย** (DOM-REQ-01)
// กฎทุกข้อบังคับที่ฝั่ง server ไม่พึ่งการซ่อนปุ่มที่หน้าจอ (plan.md ข้อ 1)
package borrowrequest

import (
	"errors"
	"time"
)

// ErrItemNotAvailable เกิดเมื่ออุปกรณ์ที่ขอยืมไม่มีอยู่จริง หรือสถานะไม่ใช่ "ว่าง" แล้ว
// ส่งคำขอได้เฉพาะชิ้นที่สถานะ "ว่าง" (CON-DATA-01, FR-BRQ-01)
var ErrItemNotAvailable = errors.New("อุปกรณ์ชิ้นนี้ไม่อยู่ในสถานะว่าง")

// ErrNoPrincipal เกิดเมื่อเรียก service โดยไม่มีผู้ใช้ที่ยืนยันตัวตนแล้วใน context (IF-IDP-01)
// ไม่เดาว่าผู้ขอเป็นใคร ปฏิเสธไปเลย (plan.md ข้อ 6.1 ของ 001)
var ErrNoPrincipal = errors.New("ไม่มีผู้ใช้ที่ยืนยันตัวตนแล้วใน context")

// NewRequest คือข้อมูลที่ผู้ขอกรอกในแบบฟอร์มแล้วส่งมา (FR-BRQ-03, FR-BRQ-04)
// ไม่มีช่องรหัสผู้ใช้ เพราะระบบใช้รหัสจากการยืนยันตัวตนเสมอ (ASM-BRQ-06)
type NewRequest struct {
	AssetCode     string // รหัสครุภัณฑ์ที่รับมาจาก FR-EQS-06 (CON-DATA-03)
	RequesterName string // ชื่อ-นามสกุลที่ผู้ขอพิมพ์เอง เป็นข้อมูลประกอบ (FR-BRQ-03 ข้อ ก)
	StudentCode   string // รหัสนักศึกษาที่ผู้ขอพิมพ์เอง ไม่ใช้ตรวจสิทธิ์ (FR-BRQ-03 ข้อ ก, ASM-BRQ-06)
	LoanType      string // ปกติ / ค้างคืน (FR-BRQ-10, DM-05)
	Purpose       string // วัตถุประสงค์ 1 ใน 3 ค่า (FR-BRQ-02)
	PurposeOther  string // ข้อความที่พิมพ์เพิ่มเมื่อเลือก "อื่น ๆ" (FR-BRQ-02 ข้อ ก)
	StartAt       time.Time
	EndAt         time.Time
}

// Request คือคำขอยืม 1 ใบที่บันทึกแล้ว (DM-06)
type Request struct {
	ID          int
	RequesterID string    // รหัสผู้ใช้จากระบบยืนยันตัวตน (IF-IDP-01, ASM-BRQ-06)
	AssetCode   string    // รหัสครุภัณฑ์ของอุปกรณ์ที่ขอยืม (CON-DATA-03)
	LoanType    string
	Purpose     string
	Status      string    // ฟีเจอร์นี้ตั้งได้ค่าเดียวคือ config.StatusRequested (FR-BRQ-04)
	SubmittedAt time.Time // เวลาที่ส่ง ใช้เรียงลำดับคำขอที่มาก่อน (FR-BRQ-04, DOM-REQ-01)
}
