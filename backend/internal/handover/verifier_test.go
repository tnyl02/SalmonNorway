package handover

import (
	"context"
	"errors"
	"testing"

	"salmonnorway/backend/internal/auth"
)

// T-08 เสร็จเมื่อ: BorrowerVerifier คืนรหัสผู้ใช้ของผู้ที่ยืนยันตัวตน ตัวจำลองใช้ผู้ใช้จาก auth.MockProvider (FR-HND-02, IF-IDP-01)
func TestVerifier_Mock(t *testing.T) {
	var v BorrowerVerifier = NewMockBorrowerVerifier(auth.NewMockProvider())
	ctx := context.Background()

	t.Run("credential ถูกต้อง คืนรหัสผู้ใช้", func(t *testing.T) {
		for credential, want := range map[string]string{
			"mock-student-001": "student-001",
			"mock-student-101": "student-101",
		} {
			got, err := v.VerifyBorrower(ctx, credential)
			if err != nil {
				t.Fatalf("VerifyBorrower(%q) error = %v", credential, err)
			}
			if got != want {
				t.Errorf("VerifyBorrower(%q) = %q ต้องเป็น %q", credential, got, want)
			}
		}
	})

	t.Run("credential ใช้ไม่ได้ คืน ErrBorrowerNotVerified", func(t *testing.T) {
		for _, credential := range []string{"", "mock-unknown"} {
			got, err := v.VerifyBorrower(ctx, credential)
			if !errors.Is(err, ErrBorrowerNotVerified) {
				t.Errorf("VerifyBorrower(%q) error = %v ต้องเป็น ErrBorrowerNotVerified", credential, err)
			}
			if got != "" {
				t.Errorf("VerifyBorrower(%q) = %q ต้องว่าง", credential, got)
			}
		}
	})
}
