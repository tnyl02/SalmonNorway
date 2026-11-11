// Package handover บันทึกการส่งมอบอุปกรณ์ให้ผู้ยืม (UC-08, specs/005-handover/spec.md)
package handover

import (
	"context"
	"errors"
	"fmt"

	"salmonnorway/backend/internal/auth"
)

// ErrBorrowerNotVerified หมายถึงยืนยันตัวตนผู้ยืมไม่สำเร็จ (credential ใช้ไม่ได้) (FR-HND-02)
var ErrBorrowerNotVerified = errors.New("borrower not verified")

// BorrowerVerifier คือจุดเดียวที่ยืนยันตัวตนผู้ยืมที่เคาน์เตอร์ผ่าน Auth API ของมหาวิทยาลัย (FR-HND-02, IF-IDP-01)
// คืนรหัสผู้ใช้จากระบบยืนยันตัวตน ซึ่งใช้เทียบกับ borrow_request.requester_id (ASM-HND-04)
// รูปแบบ credential ที่เคาน์เตอร์ยังรอ Q-HND-11 ตัวจริงจะมาแทนใน T-29
// timeout 10 วินาทีและการใช้ผลตรวจบัตรแทนอยู่ใน T-13 ไม่ใช่ที่นี่
type BorrowerVerifier interface {
	VerifyBorrower(ctx context.Context, credential string) (userID string, err error)
}

// MockBorrowerVerifier คือตัวยืนยันผู้ยืมจำลองสำหรับ prototype ใช้ผู้ใช้จาก auth.Provider
// (ทีมเลือกเอง ไม่ได้มาจาก spec รอ Q-HND-11 และ Q-SH-03) credential คือ token แบบ "mock-<รหัสผู้ใช้>"
type MockBorrowerVerifier struct {
	provider auth.Provider
}

// NewMockBorrowerVerifier สร้างตัวยืนยันผู้ยืมจำลองจาก provider ที่ให้มา เช่น auth.NewMockProvider() (FR-HND-02)
func NewMockBorrowerVerifier(provider auth.Provider) *MockBorrowerVerifier {
	return &MockBorrowerVerifier{provider: provider}
}

// VerifyBorrower คืนรหัสผู้ใช้ของ credential หรือ ErrBorrowerNotVerified (FR-HND-02, IF-IDP-01)
func (v *MockBorrowerVerifier) VerifyBorrower(ctx context.Context, credential string) (string, error) {
	p, err := v.provider.Authenticate(ctx, credential)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrBorrowerNotVerified, err)
	}
	return p.UserID, nil
}
