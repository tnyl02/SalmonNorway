package borrowrequest

import "strings"

// ชื่อช่องที่ส่งกลับไปบอกหน้าจอ ใช้ชื่อเดียวกับใน JSON ที่หน้าจอส่งมา (CreateRequestDTO)
// เพื่อให้หน้าจอชี้จุดที่ต้องแก้ได้ถูกช่อง (FR-BRQ-03, การแสดงผลเป็นงานของ T-20)
const (
	fieldRequesterName = "requesterName"
	fieldStudentCode   = "studentCode"
	fieldStartAt       = "startAt"
	fieldEndAt         = "endAt"
	fieldPurpose       = "purpose"
)

// missingFields คืนรายชื่อช่องที่ยังไม่ได้กรอก ตาม 5 ช่องที่ทีมตัดสินเมื่อ 2569-11-11 (FR-BRQ-03, RQ-22)
// ชื่อ-นามสกุล, รหัสนักศึกษา, วันและเวลาที่เริ่มยืม, วันและเวลาที่คืน, วัตถุประสงค์การยืม
//
// ตรวจแค่ว่า "กรอกมาหรือยัง" เท่านั้น
// ค่าวัตถุประสงค์ต้องเป็น 1 ใน 3 ค่าหรือไม่เป็นงานของ T-07 (FR-BRQ-02)
// กฎเวลาคืนของการยืมปกติเป็นงานของ T-08 (FR-BRQ-03 ข้อ ค, DOM-LOAN-01)
func missingFields(in NewRequest) []string {
	var missing []string
	if strings.TrimSpace(in.RequesterName) == "" {
		missing = append(missing, fieldRequesterName)
	}
	if strings.TrimSpace(in.StudentCode) == "" {
		missing = append(missing, fieldStudentCode)
	}
	if in.StartAt.IsZero() {
		missing = append(missing, fieldStartAt)
	}
	if in.EndAt.IsZero() {
		missing = append(missing, fieldEndAt)
	}
	if strings.TrimSpace(in.Purpose) == "" {
		missing = append(missing, fieldPurpose)
	}
	return missing
}
