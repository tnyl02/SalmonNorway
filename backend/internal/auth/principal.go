// Package auth ตรวจว่าผู้ใช้ยืนยันตัวตนแล้วก่อนเข้าถึงข้อมูลทุกครั้ง (IF-IDP-01, NFR-SEC-01)
//
// ตอนนี้ใช้ผู้ให้บริการจำลอง (mock_provider.go) เพราะยังไม่รู้โปรโตคอลของระบบยืนยันตัวตนของมหาวิทยาลัย (Q-SH-03)
// เป็นการตัดสินใจของทีม ไม่ได้มาจาก spec เมื่อได้คำตอบ Q-SH-03 ให้เขียน Provider ตัวจริงมาแทน (T-30)
package auth

import "context"

// Role คือบทบาทของผู้ใช้ระบบ มี 2 บทบาทตาม Use Case Diagram (specs/000-shared/constraints.md, CR-002)
type Role string

const (
	RoleStudent Role = "student" // นักศึกษา
	RoleStaff   Role = "staff"   // เจ้าหน้าที่ (นักเทคโนโลยีสารสนเทศ / นักคอมพิวเตอร์)
)

// Valid บอกว่าเป็นบทบาทที่ระบบรู้จักหรือไม่
func (r Role) Valid() bool {
	return r == RoleStudent || r == RoleStaff
}

// Principal คือผู้ใช้ที่ยืนยันตัวตนแล้ว (DM-08) ได้จากระบบยืนยันตัวตน (IF-IDP-01)
type Principal struct {
	UserID       string
	Role         Role
	DepartmentID int // ภาควิชาของผู้ใช้ ใช้ตรวจ DOM-ACCESS-01 ค่า 0 = ไม่ได้รับมา
}

// Complete บอกว่าได้ข้อมูลครบตาม IF-IDP-01 (มีรหัสผู้ใช้ บทบาทที่รู้จัก และภาควิชา) ไม่ครบ = ต้องปฏิเสธ ไม่เดา (plan 6.1)
func (p Principal) Complete() bool {
	return p.UserID != "" && p.Role.Valid() && p.DepartmentID > 0
}

type ctxKey struct{}

// WithPrincipal ใส่ผู้ใช้ลงใน context ของ request (ใช้โดย middleware และใน test ของฟีเจอร์อื่น)
func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, ctxKey{}, p)
}

// FromContext อ่านผู้ใช้จาก context ที่ middleware ใส่ไว้ ok = false ถ้าไม่มี
func FromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(ctxKey{}).(Principal)
	return p, ok
}
