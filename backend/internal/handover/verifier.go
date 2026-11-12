// Package handover บันทึกการส่งมอบอุปกรณ์ให้ผู้ยืม (UC-08, specs/005-handover/spec.md)
package handover

import (
	"context"
	"errors"
	"fmt"
	"time"

	"salmonnorway/backend/internal/auth"
)

var (
	// ErrBorrowerNotVerified หมายถึงยืนยันตัวตนผู้ยืมไม่สำเร็จ (credential ใช้ไม่ได้) (FR-HND-02)
	ErrBorrowerNotVerified = errors.New("borrower not verified")
	// ErrBorrowerAuthUnavailable หมายถึง Auth API ไม่ได้ผลภายใน BorrowerAuthTimeout ถือว่าเชื่อมต่อไม่สำเร็จ
	// ให้ใช้ผลการตรวจบัตรแทน (FR-HND-02, ASM-HND-02, AC-HND-03)
	ErrBorrowerAuthUnavailable = errors.New("borrower auth api unavailable")
)

// BorrowerVerifier คือจุดเดียวที่ยืนยันตัวตนผู้ยืมที่เคาน์เตอร์ผ่าน Auth API ของมหาวิทยาลัย (FR-HND-02, IF-IDP-01)
// คืนรหัสผู้ใช้จากระบบยืนยันตัวตน ซึ่งใช้เทียบกับ borrow_request.requester_id (ASM-HND-04)
// ผู้ยืมล็อกอินบนเครื่องเจ้าหน้าที่ (ASM-HND-09) ตัวจริงจะมาแทนตัวจำลองใน T-29
// การรอไม่เกิน 10 วินาทีอยู่ที่ WithTimeout (T-13)
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

// timeoutVerifier ห่อ BorrowerVerifier ให้รอผลไม่เกิน timeout (FR-HND-02, ASM-HND-02, IF-IDP-01)
type timeoutVerifier struct {
	inner   BorrowerVerifier
	timeout time.Duration
}

// WithTimeout คืน verifier ที่ถ้า inner ไม่ได้ผลภายใน timeout จะคืน ErrBorrowerAuthUnavailable (FR-HND-02, AC-HND-03)
// ใช้งานจริงส่ง config.BorrowerAuthTimeout (10 วินาที) test ส่งค่าสั้นเพื่อไม่ต้องรอจริง (plan.md ข้อ 6)
func WithTimeout(inner BorrowerVerifier, timeout time.Duration) BorrowerVerifier {
	return &timeoutVerifier{inner: inner, timeout: timeout}
}

// VerifyBorrower เรียก inner ใน goroutine แยก จึงหยุดรอได้ตรงเวลาแม้ inner ไม่สนใจ ctx (FR-HND-02)
// ถ้า ctx ของผู้เรียกถูกยกเลิกเอง (เช่น หน้าจอปิดการเชื่อมต่อ) คืน error ของ ctx ไม่นับเป็น Auth API ไม่ตอบ
func (v *timeoutVerifier) VerifyBorrower(ctx context.Context, credential string) (string, error) {
	tctx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	type result struct {
		userID string
		err    error
	}
	done := make(chan result, 1) // มีที่ว่าง 1 ช่อง goroutine จึงจบได้แม้ไม่มีใครรอผลแล้ว
	go func() {
		id, err := v.inner.VerifyBorrower(tctx, credential)
		done <- result{id, err}
	}()

	select {
	case r := <-done:
		if r.err != nil && ctx.Err() == nil && errors.Is(tctx.Err(), context.DeadlineExceeded) {
			return "", ErrBorrowerAuthUnavailable // inner คืน error เพราะหมดเวลาพอดี
		}
		return r.userID, r.err
	case <-tctx.Done():
		if err := ctx.Err(); err != nil {
			return "", err
		}
		return "", ErrBorrowerAuthUnavailable
	}
}
