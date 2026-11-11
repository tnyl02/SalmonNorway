package borrowrequest

import (
	"slices"
	"strings"

	"salmonnorway/backend/internal/config"
)

// ชื่อช่องที่ส่งกลับไปบอกหน้าจอ ใช้ชื่อเดียวกับใน JSON ที่หน้าจอส่งมา (CreateRequestDTO)
// เพื่อให้หน้าจอชี้จุดที่ต้องแก้ได้ถูกช่อง (FR-BRQ-03, การแสดงผลเป็นงานของ T-20)
const (
	fieldRequesterName = "requesterName"
	fieldStudentCode   = "studentCode"
	fieldStartAt       = "startAt"
	fieldEndAt         = "endAt"
	fieldPurpose       = "purpose"
	fieldPurposeOther  = "purposeOther"
)

// missingFields คืนรายชื่อช่องที่ยังกรอกไม่ถูกต้อง ให้หน้าจอชี้จุดที่ต้องแก้ (FR-BRQ-03, RQ-22)
// ครอบคลุม 5 ช่องที่ทีมตัดสินเมื่อ 2569-11-11
// (ชื่อ-นามสกุล, รหัสนักศึกษา, วันและเวลาที่เริ่มยืม, วันและเวลาที่คืน, วัตถุประสงค์การยืม)
// บวกกฎของช่องวัตถุประสงค์ตาม FR-BRQ-02 ที่เพิ่มใน T-07
//
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
	missing = append(missing, purposeProblems(in)...)
	return missing
}

// purposeProblems ตรวจช่องวัตถุประสงค์การยืมตาม FR-BRQ-02 (ทีมตัดสินค่าทั้ง 3 เมื่อ 2569-11-11)
//   - ต้องเลือกมา 1 ค่า และต้องเป็นค่าที่ระบบรู้จัก (config.PurposeValues)
//   - ถ้าเลือก "อื่น ๆ" ต้องพิมพ์ระบุเพิ่มด้วย จึงจะส่งคำขอได้ (FR-BRQ-02 ข้อ ก)
func purposeProblems(in NewRequest) []string {
	purpose := strings.TrimSpace(in.Purpose)
	if purpose == "" || !slices.Contains(config.PurposeValues, purpose) {
		// ไม่ได้เลือก หรือเลือกค่าที่ไม่มีในรายการ ถือเป็นช่องเดียวกันที่ต้องแก้
		return []string{fieldPurpose}
	}
	if purpose == config.PurposeOther && strings.TrimSpace(in.PurposeOther) == "" {
		return []string{fieldPurposeOther}
	}
	return nil
}
