package auth

import (
	"context"
	"errors"
)

// ErrInvalidToken หมายถึงไม่มี token หรือ token ใช้ไม่ได้ (หมดอายุหรือไม่รู้จัก)
var ErrInvalidToken = errors.New("invalid or expired token")

// Provider คือจุดเดียวที่ต่อกับระบบยืนยันตัวตน (IF-IDP-01) เปลี่ยนเป็นตัวจริงได้โดยไม่ต้องแก้ middleware (T-30, Q-SH-03)
type Provider interface {
	Authenticate(ctx context.Context, token string) (Principal, error)
}

// MockProvider คือระบบยืนยันตัวตนจำลองสำหรับ prototype (ทีมเลือกเอง ไม่ได้มาจาก spec รอ Q-SH-03)
// ผู้ใช้ตัวอย่างตรงกับข้อมูลตั้งต้นใน backend/db/seed/seed.sql (ภาควิชา 1 และ 2, ผู้ยืม student-001 และ student-002)
type MockProvider struct {
	users map[string]Principal // token -> ผู้ใช้
}

// NewMockProvider สร้างผู้ใช้จำลอง ใช้ token แบบ "mock-<รหัสผู้ใช้>"
func NewMockProvider() *MockProvider {
	return &MockProvider{users: map[string]Principal{
		"mock-student-001": {UserID: "student-001", Role: RoleStudent, DepartmentID: 1},
		"mock-student-002": {UserID: "student-002", Role: RoleStudent, DepartmentID: 1},
		"mock-student-101": {UserID: "student-101", Role: RoleStudent, DepartmentID: 2},
		"mock-staff-001":   {UserID: "staff-001", Role: RoleStaff, DepartmentID: 1},
	}}
}

// Authenticate คืนผู้ใช้ของ token หรือ ErrInvalidToken (IF-IDP-01)
func (m *MockProvider) Authenticate(_ context.Context, token string) (Principal, error) {
	p, ok := m.users[token]
	if !ok {
		return Principal{}, ErrInvalidToken
	}
	return p, nil
}
