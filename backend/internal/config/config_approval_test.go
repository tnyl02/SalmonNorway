package config

// test ของค่าที่ตั้งได้ของฟีเจอร์ 004 (T-02, T-41 ใน specs/004-approve-request/tasks.md)
// แยกไฟล์จาก config_test.go ของ 002 เพื่อไม่แก้ไฟล์ของฟีเจอร์อื่น

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"
)

// T-02 (004) เสร็จเมื่อ: ค่าครบทุกแถวในตาราง "ค่าที่ตั้งได้" ของ plan ข้อ 5 และตรงกับ spec
func TestConfig_ApprovalValues(t *testing.T) {
	if ApprovalEmailDeadline != 60*time.Second {
		t.Errorf("ต้องส่ง Email ภายใน 60 วินาที (FR-APV-02) แต่ได้ %v", ApprovalEmailDeadline)
	}
	if ApprovalEmailRetryCount != 3 {
		t.Errorf("ส่งซ้ำ 3 ครั้งหลังครั้งแรก (FR-APV-02 ง, Q-APV-17) แต่ได้ %d", ApprovalEmailRetryCount)
	}
	if ApprovalEmailRetryInterval != time.Minute {
		t.Errorf("ส่งซ้ำห่างกันครั้งละ 1 นาที (FR-APV-02 ง) แต่ได้ %v", ApprovalEmailRetryInterval)
	}
	if ApprovalAutoRejectHour != 0 {
		t.Errorf("ปฏิเสธอัตโนมัติเวลา 00:00 น. (FR-APV-01 ฉ) แต่ได้ %d", ApprovalAutoRejectHour)
	}
	// FR-APV-01 (ค), Q-APV-11: เวลาตัดใช้ค่าของ 002 ร่วม
	if LoanNormalReturnHour != 16 {
		t.Errorf("อนุมัติไม่ได้หลัง 16:00 น. (FR-APV-01 ค) แต่ LoanNormalReturnHour = %d", LoanNormalReturnHour)
	}
	if ApprovalSupersededReason != "อุปกรณ์ชิ้นนี้ได้รับการอนุมัติให้คำขออื่นแล้ว" {
		t.Errorf("เหตุผลปฏิเสธตามไม่ตรง FR-APV-01 (จ): %q", ApprovalSupersededReason)
	}
	if ApprovalAutoRejectReason != "คำขอยังคงค้างอยู่และไม่ได้รับการพิจารณาภายในเวลาที่กำหนด" {
		t.Errorf("เหตุผลปฏิเสธอัตโนมัติไม่ตรง FR-APV-01 (ฉ): %q", ApprovalAutoRejectReason)
	}
	if ApprovalPickupLocation != "คณะวิทยาศาสตร์ ตึก 1 ชั้น 2 ห้องอุปกรณ์ภาควิชาคอมพิวเตอร์" {
		t.Errorf("สถานที่รับไม่ตรง FR-APV-02 (ก): %q", ApprovalPickupLocation)
	}
	if ApprovalAutoRejectEmail.Guidance != "หากยังต้องการยืมอุปกรณ์ สามารถยื่นคำขอใหม่ได้" ||
		ApprovalAutoRejectEmail.Status != "ปฏิเสธอัตโนมัติ" || ApprovalAutoRejectEmail.DoneAt != "00:00 น." {
		t.Errorf("ค่าตายตัวใน Email ปฏิเสธอัตโนมัติไม่ตรง FR-APV-02 (ค): %+v", ApprovalAutoRejectEmail)
	}

	// สถานะที่ 004 ตั้งต้องอยู่ในรายการ 8 ค่าของ 002 (data-model.md v9)
	for _, s := range []string{StatusApproved, StatusRejected} {
		if !slices.Contains(RequestStatuses, s) {
			t.Errorf("สถานะ %q ต้องอยู่ใน RequestStatuses", s)
		}
	}

	// ค่า action ต้องตรงกับ CHECK ใน migration 0008
	if AuditActionApprove != "approve" || AuditActionReject != "reject" || AuditActionAutoReject != "auto_reject" {
		t.Error("ค่า action ของ audit log ต้องเป็น approve / reject / auto_reject (migration 0008)")
	}

	wantSubjects := map[string]string{
		"approved":      "แจ้งผลการอนุมัติการยืมอุปกรณ์ภาคคอมพิวเตอร์",
		"rejected":      "แจ้งผลการปฏิเสธการยืมอุปกรณ์ภาคคอมพิวเตอร์",
		"auto_rejected": "แจ้งปฏิเสธคำขอยืมอุปกรณ์อัตโนมัติ",
	}
	for kind, want := range wantSubjects {
		if ApprovalEmailSubjects[kind] != want {
			t.Errorf("หัวข้อ Email แบบ %q ต้องเป็น %q แต่ได้ %q", kind, want, ApprovalEmailSubjects[kind])
		}
		if len(ApprovalEmailTemplates[kind]) == 0 {
			t.Errorf("Email แบบ %q ต้องมีรายการข้อมูล (FR-APV-02)", kind)
		}
	}
	// FR-APV-02 (ก): ยืมปกติไม่แสดงจำนวนอุปกรณ์
	for _, f := range ApprovalEmailTemplates[ApprovalKindApproved] {
		if strings.Contains(f, "จำนวน") {
			t.Errorf("Email ผลอนุมัติต้องไม่มีจำนวนอุปกรณ์ (FR-APV-02 ก) แต่พบ %q", f)
		}
	}
}

// T-02 (004) เสร็จเมื่อ: ข้อความ error ครบ 9 รหัส และหน้าจอ (config.js) ใช้ข้อความเดียวกัน (FR-APV-01 ซ, CON-TECH-03)
func TestConfig_ApprovalMessagesMatchFrontend(t *testing.T) {
	codes := []string{"ITEM_NOT_AVAILABLE", "REQUEST_CANCELLED", "REQUEST_ALREADY_DECIDED",
		"REQUEST_DECIDED_CONCURRENTLY", "NOT_FIRST_IN_QUEUE", "APPROVAL_CUTOFF_PASSED",
		"VALIDATION_FAILED", "FORBIDDEN_ROLE", "FORBIDDEN_DEPARTMENT", "UPSTREAM_UNAVAILABLE", "REQUEST_NOT_FOUND"}
	if len(ApprovalMessages) != len(codes) {
		t.Errorf("ต้องมีข้อความ %d รหัส แต่ได้ %d", len(codes), len(ApprovalMessages))
	}

	src, err := os.ReadFile(filepath.Join("..", "..", "..", "frontend", "src", "config.js"))
	if err != nil {
		t.Fatalf("อ่าน frontend/src/config.js ไม่ได้: %v", err)
	}
	for _, code := range codes {
		msg, ok := ApprovalMessages[code]
		if !ok || msg == "" {
			t.Errorf("ไม่มีข้อความของรหัส %s", code)
			continue
		}
		re := regexp.MustCompile(code + `:\s*'([^']*)'`)
		m := re.FindSubmatch(src)
		if m == nil {
			t.Errorf("APPROVAL_MESSAGES ใน config.js ไม่มีรหัส %s", code)
			continue
		}
		if string(m[1]) != msg {
			t.Errorf("ข้อความของ %s ไม่ตรงกัน: backend %q / frontend %q", code, msg, m[1])
		}
	}
}

// T-41 (004) เสร็จเมื่อ: ค่าของคำขอค้างคืนตรงกับ spec (FR-APV-03 ช, Q-APV-44) และหน้าจอใช้ค่าเดียวกัน
// ใช้ค่า loan_type และชื่อสถานะของ 002 ที่มีอยู่แล้ว ไม่ประกาศซ้ำ
func TestConfig_ApprovalOvernightValues(t *testing.T) {
	if ApprovalAdvisorNameMaxLength != 100 {
		t.Errorf("ชื่ออาจารย์ยาวไม่เกิน 100 ตัวอักษร (Q-APV-44) แต่ได้ %d", ApprovalAdvisorNameMaxLength)
	}
	if ApprovalAdvisorNameRequiredMessage != "กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ" {
		t.Errorf("ข้อความไม่ระบุชื่ออาจารย์ไม่ตรง FR-APV-01 (ซ): %q", ApprovalAdvisorNameRequiredMessage)
	}
	if ApprovalAdvisorNameTooLongMessage != "ชื่ออาจารย์ต้องไม่เกิน 100 ตัวอักษร" {
		t.Errorf("ข้อความชื่ออาจารย์ยาวเกินไม่ตรง FR-APV-01 (ซ): %q", ApprovalAdvisorNameTooLongMessage)
	}
	if LoanTypeOvernight != "overnight" || StatusApprovedByOfficer != "approved_by_officer" ||
		StatusPendingAdvisorApproval != "pending_advisor_approval" {
		t.Error("ค่า loan_type และสถานะของ 002 ต้องตรงกับ CHECK ของ approval_step ใน migration 0011")
	}

	src, err := os.ReadFile(filepath.Join("..", "..", "..", "frontend", "src", "config.js"))
	if err != nil {
		t.Fatalf("อ่าน frontend/src/config.js ไม่ได้: %v", err)
	}
	if m := regexp.MustCompile(`APPROVAL_ADVISOR_NAME_MAX_LENGTH\s*=\s*(\d+)`).FindSubmatch(src); m == nil || string(m[1]) != "100" {
		t.Errorf("APPROVAL_ADVISOR_NAME_MAX_LENGTH ใน config.js ต้องเป็น 100 แต่ได้ %q", m)
	}
	for name, want := range map[string]string{
		"APPROVAL_ADVISOR_NAME_REQUIRED": ApprovalAdvisorNameRequiredMessage,
		"APPROVAL_ADVISOR_NAME_TOO_LONG": ApprovalAdvisorNameTooLongMessage,
	} {
		m := regexp.MustCompile(name + `\s*=\s*'([^']*)'`).FindSubmatch(src)
		if m == nil || string(m[1]) != want {
			t.Errorf("%s ใน config.js ต้องเป็น %q แต่ได้ %q", name, want, m)
		}
	}
}
