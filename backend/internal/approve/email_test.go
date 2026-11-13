package approve

import (
	"context"
	"strings"
	"testing"
	"time"

	"salmonnorway/backend/internal/config"
)

// ข้อมูลตัวอย่าง: ยื่นคำขอ 08:30 น. ยืม 09:00 ถึง 16:00 น. ของวันที่ 12 พ.ย. 2026 (เวลาไทย)
func sampleEmailData() EmailData {
	bkk := time.FixedZone("ICT", 7*60*60)
	return EmailData{
		RequestID:     42,
		RequesterName: "นักศึกษาตัวอย่าง ก",
		ModelName:     "อุปกรณ์ทดลองอนุมัติ (004)",
		AssetCode:     "AP-TEST-0002",
		Purpose:       "study",
		StartAt:       time.Date(2026, 11, 12, 9, 0, 0, 0, bkk),
		EndAt:         time.Date(2026, 11, 12, 16, 0, 0, 0, bkk),
		SubmittedAt:   time.Date(2026, 11, 12, 8, 30, 0, 0, bkk),
		Reason:        "อุปกรณ์ต้องใช้สอบ",
		ActorName:     "staff-001",
	}
}

// T-14 เสร็จเมื่อ: Email แต่ละแบบมีหัวข้อและข้อมูลครบตาม FR-APV-02 (ก)(ข)(ค) จาก config และ Email ผลอนุมัติไม่มีจำนวนอุปกรณ์
func TestBuildEmail_Contents(t *testing.T) {
	d := sampleEmailData()
	cases := map[string][]string{
		config.ApprovalKindApproved: {
			"ชื่อผู้ขอ: นักศึกษาตัวอย่าง ก",
			"รายการอุปกรณ์ที่ขอยืม: อุปกรณ์ทดลองอนุมัติ (004) (AP-TEST-0002)",
			"วัตถุประสงค์: เพื่อการเรียน",
			"วันเวลาที่ต้องการยืม: 2569-11-12 09:00",
			"วันเวลากำหนดคืน: 2569-11-12 16:00",
			"ชื่อผู้อนุมัติ: staff-001",
			"สถานที่รับ: " + config.ApprovalPickupLocation,
		},
		config.ApprovalKindRejected: {
			"ชื่อผู้ขอ: นักศึกษาตัวอย่าง ก",
			"รายการอุปกรณ์ที่ขอยืม: อุปกรณ์ทดลองอนุมัติ (004) (AP-TEST-0002)",
			"วันที่ยื่นคำขอ: 2569-11-12",
			"เหตุผลในการปฏิเสธ: อุปกรณ์ต้องใช้สอบ",
			"ชื่อผู้พิจารณา: staff-001",
		},
		config.ApprovalKindAutoRejected: {
			"เลขที่คำขอ: 42",
			"ชื่อผู้ขอ: นักศึกษาตัวอย่าง ก",
			"รายการอุปกรณ์ที่ขอยืม: อุปกรณ์ทดลองอนุมัติ (004) (AP-TEST-0002)",
			"วันที่ยื่นคำขอ: 2569-11-12",
			"สถานะ: ปฏิเสธอัตโนมัติ",
			"เหตุผล: คำขอยังคงค้างอยู่และไม่ได้รับการพิจารณาภายในเวลาที่กำหนด",
			"เวลาที่ดำเนินการ: 00:00 น.",
			"แนวทาง: หากยังต้องการยืมอุปกรณ์ สามารถยื่นคำขอใหม่ได้",
		},
	}
	for kind, wantLines := range cases {
		subject, body, err := BuildEmail(kind, d)
		if err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
		if subject != config.ApprovalEmailSubjects[kind] {
			t.Errorf("%s: หัวข้อต้องเป็น %q แต่ได้ %q", kind, config.ApprovalEmailSubjects[kind], subject)
		}
		got := strings.Split(strings.TrimSpace(body), "\n")
		if strings.Join(got, "\n") != strings.Join(wantLines, "\n") {
			t.Errorf("%s: เนื้อ Email ไม่ตรง\nได้:\n%s\nต้องการ:\n%s", kind, strings.Join(got, "\n"), strings.Join(wantLines, "\n"))
		}
	}

	// FR-APV-02 (ก): ยืมปกติไม่แสดงจำนวนอุปกรณ์
	if _, body, _ := BuildEmail(config.ApprovalKindApproved, d); strings.Contains(body, "จำนวน") {
		t.Error("Email ผลอนุมัติต้องไม่มีจำนวนอุปกรณ์ (FR-APV-02 ก)")
	}
}

// FR-BRQ-02 (ก): วัตถุประสงค์ "อื่น ๆ" แสดงข้อความที่ผู้ขอพิมพ์
// T-52: Email ผลอนุมัติของคำขอค้างคืนมีชื่ออาจารย์ต่อจากชื่อผู้อนุมัติ ยืมปกติไม่มี และไม่มีจำนวนอุปกรณ์ (FR-APV-02 ก v3.1, rtm F-08)
func TestBuildEmail_OvernightAdvisor(t *testing.T) {
	d := EmailData{RequesterName: "ค", ModelName: "รุ่น", AssetCode: "AP-1", Purpose: "study",
		StartAt: time.Now(), EndAt: time.Now(), ActorName: "staff-001", AdvisorName: "อ.สมชาย ใจดี"}
	_, body, err := BuildEmail(config.ApprovalKindApproved, d)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, "ชื่อผู้อนุมัติ: staff-001\nชื่ออาจารย์ที่เซ็นอนุมัติ: อ.สมชาย ใจดี\n") {
		t.Errorf("ต้องมีชื่ออาจารย์ต่อจากชื่อผู้อนุมัติ:\n%s", body)
	}
	if strings.Contains(body, "จำนวน") {
		t.Errorf("ต้องไม่มีจำนวนอุปกรณ์:\n%s", body)
	}
	d.AdvisorName = ""
	if _, body, _ = BuildEmail(config.ApprovalKindApproved, d); strings.Contains(body, "ชื่ออาจารย์") {
		t.Errorf("ยืมปกติต้องไม่มีชื่ออาจารย์:\n%s", body)
	}
}

func TestBuildEmail_PurposeOther(t *testing.T) {
	d := sampleEmailData()
	d.Purpose, d.PurposeOther = config.PurposeOther, "ถ่ายวิดีโอโครงงาน"
	_, body, err := BuildEmail(config.ApprovalKindApproved, d)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, "วัตถุประสงค์: ถ่ายวิดีโอโครงงาน\n") {
		t.Errorf("ต้องแสดงข้อความที่ผู้ขอพิมพ์ แต่ได้:\n%s", body)
	}
}

// T-14 เสร็จเมื่อ: RecipientResolver ตัวจำลองคืนที่อยู่ตายตัวต่อผู้ขอ (ตัวจริงรอ Q-APV-04, T-37)
func TestMockRecipientResolver(t *testing.T) {
	got, err := MockRecipientResolver{}.Resolve(context.Background(), "student-001")
	if err != nil || got != "student-001@example.com" {
		t.Errorf("ต้องได้ student-001@example.com แต่ได้ %q (%v)", got, err)
	}
}
