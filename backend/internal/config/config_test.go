package config

import (
	"slices"
	"testing"
)

// T-02 เสร็จเมื่อ: ค่าที่มาจาก spec ของฟีเจอร์ 002 อยู่ในไฟล์เดียวและตรงกับตาราง "ค่าที่ตั้งได้" ใน plan.md ข้อ 5
func TestConfig_BorrowRequestValues(t *testing.T) {
	if LoanNormalReturnHour != 16 {
		t.Errorf("การยืมปกติต้องคืนภายใน 16:00 น. (DOM-LOAN-01) แต่ได้ %d", LoanNormalReturnHour)
	}
	if !OverdueBlockUntilNextDayStart {
		t.Error("ผู้ที่คืนของเกินกำหนดต้องส่งคำขอใหม่ได้เมื่อถึงต้นวันถัดไป (FR-BRQ-06 ข้อ ข)")
	}
	if AttachmentMaxBytes != 10*1024*1024 {
		t.Errorf("ไฟล์แนบสูงสุด 10MB (FR-BRQ-10 ข้อ ก) แต่ได้ %d ไบต์", AttachmentMaxBytes)
	}
	if TimeZone != "Asia/Bangkok" {
		t.Errorf("ต้องเทียบเวลาด้วย Asia/Bangkok (glossary) แต่ได้ %q", TimeZone)
	}
	if AuditActionRequest != "request" {
		t.Errorf("action ของ audit log ต้องเป็น request (Q-BRQ-07) แต่ได้ %q", AuditActionRequest)
	}

	wantTypes := []string{"application/pdf", "image/jpeg", "image/png"}
	if !slices.Equal(AttachmentAllowedTypes, wantTypes) {
		t.Errorf("ชนิดไฟล์แนบต้องเป็น %v (FR-BRQ-10 ข้อ ก) แต่ได้ %v", wantTypes, AttachmentAllowedTypes)
	}

	wantPurposes := []string{"study", "activity", "other"}
	if !slices.Equal(PurposeValues, wantPurposes) {
		t.Errorf("วัตถุประสงค์ต้องมี 3 ค่า %v (FR-BRQ-02) แต่ได้ %v", wantPurposes, PurposeValues)
	}
	if !slices.Contains(PurposeValues, PurposeOther) {
		t.Error("PurposeOther ต้องเป็นหนึ่งในค่าของ PurposeValues")
	}

	// data-model.md v9: สถานะของคำขอมี 8 ค่า และ 002 ตั้งได้เฉพาะ requested (FR-BRQ-04)
	wantStatuses := []string{"requested", "approved_by_officer", "pending_advisor_approval",
		"approved", "rejected", "handed_over", "returned", "cancelled"}
	if !slices.Equal(RequestStatuses, wantStatuses) {
		t.Errorf("สถานะของคำขอต้องมี 8 ค่าตาม data-model.md v9 แต่ได้ %v", RequestStatuses)
	}
	if RequestStatuses[0] != StatusRequested {
		t.Errorf("สถานะเริ่มต้นที่ 002 ตั้งต้องเป็น %q", StatusRequested)
	}
}

// T-02 เสร็จเมื่อ: โฟลเดอร์เก็บไฟล์แนบอ่านจาก environment ได้ (ASM-BRQ-07)
func TestConfig_AttachmentDirFromEnv(t *testing.T) {
	if got := Load().AttachmentDir; got == "" {
		t.Error("AttachmentDir ต้องมีค่าเริ่มต้นสำหรับเครื่องพัฒนา")
	}
	t.Setenv("ATTACHMENT_DIR", "/tmp/salmon-attachments")
	if got := Load().AttachmentDir; got != "/tmp/salmon-attachments" {
		t.Errorf("ATTACHMENT_DIR ต้องเปลี่ยนค่าได้ แต่ได้ %q", got)
	}
}
