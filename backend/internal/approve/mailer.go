package approve

import (
	"context"
	"log"
)

// Email คือ Email 1 ฉบับที่ส่งให้ผู้ขอ (FR-APV-02, IF-NOT-01)
type Email struct {
	To      string
	Subject string
	Body    string
}

// Mailer คือช่องทางส่ง Email ผ่านระบบแจ้งเตือน (IF-NOT-01, plan.md ข้อ 2)
// เป็น interface เปลี่ยนเป็นตัวส่งจริงได้โดยไม่แก้ส่วนอื่น (ทีมเลือกเอง ไม่ได้มาจาก spec)
// ส่งไม่สำเร็จให้คืน error เพื่อให้ตัวส่งเบื้องหลังส่งซ้ำ (FR-APV-02 ง, T-32)
type Mailer interface {
	Send(ctx context.Context, e Email) error
}

// LogMailer คือ Mailer จำลองของ prototype เขียน Email ลง log ของหลังบ้านแทนการส่งจริง
// แบบเดียวกับ LogNotifier ของ 003 (ทีมเลือกเอง ไม่ได้มาจาก spec)
type LogMailer struct{}

// Send เขียน Email ลง log สำเร็จเสมอ
func (LogMailer) Send(_ context.Context, e Email) error {
	log.Printf("approve: [Email จำลอง] ถึง %s\nหัวข้อ: %s\n%s", e.To, e.Subject, e.Body)
	return nil
}

// RecipientResolver หา Email มหาวิทยาลัยของผู้ขอจากรหัสผู้ใช้ (FR-APV-02 "Email มหาวิทยาลัยของผู้ขอที่ใช้ login")
// ที่มาจริงยังไม่รู้ (Q-APV-04) ตัวจริงเป็นงานของ T-37
type RecipientResolver interface {
	Resolve(ctx context.Context, requesterID string) (string, error)
}

// MockRecipientResolver คืนที่อยู่จำลองตายตัวต่อผู้ขอ "<รหัสผู้ใช้>@example.com" ระหว่างรอ Q-APV-04
// example.com เป็นโดเมนที่สงวนไว้สำหรับตัวอย่าง ส่งไม่ถึงใครจริง (ทีมเลือกเอง ไม่ได้มาจาก spec)
type MockRecipientResolver struct{}

// Resolve คืนที่อยู่จำลองของผู้ขอ
func (MockRecipientResolver) Resolve(_ context.Context, requesterID string) (string, error) {
	return requesterID + "@example.com", nil
}
