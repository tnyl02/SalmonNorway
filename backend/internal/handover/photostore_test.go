package handover

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// T-07 เสร็จเมื่อ: เขียนไฟล์ลง <root>/<loanId>/ ใน temp dir และลบทิ้งได้เมื่อ rollback (FR-HND-06)
func TestPhotoStore_SaveAndRemove(t *testing.T) {
	root := t.TempDir()
	store := NewPhotoStore(root)
	const loanID = 42

	jpg := bytes.Repeat([]byte{0xFF}, 1024)
	path1, size1, err := store.Save(loanID, 1, "image/jpeg", bytes.NewReader(jpg))
	if err != nil {
		t.Fatalf("Save jpeg error = %v", err)
	}
	path2, size2, err := store.Save(loanID, 2, "image/png", bytes.NewReader([]byte("png")))
	if err != nil {
		t.Fatalf("Save png error = %v", err)
	}

	wantDir := filepath.Join(root, "42")
	for _, tc := range []struct {
		path string
		size int64
		want string
		n    int64
	}{
		{path1, size1, filepath.Join(wantDir, "1.jpg"), 1024},
		{path2, size2, filepath.Join(wantDir, "2.png"), 3},
	} {
		if tc.path != tc.want {
			t.Errorf("path = %q ต้องเป็น %q", tc.path, tc.want)
		}
		if tc.size != tc.n {
			t.Errorf("size = %d ต้องเป็น %d", tc.size, tc.n)
		}
		if _, err := os.Stat(tc.path); err != nil {
			t.Errorf("ไม่พบไฟล์ %q: %v", tc.path, err)
		}
	}

	t.Run("ชนิดที่ไม่รู้จักไม่เขียนไฟล์", func(t *testing.T) {
		_, _, err := store.Save(loanID, 3, "image/gif", bytes.NewReader([]byte("gif")))
		if !errors.Is(err, ErrPhotoType) {
			t.Fatalf("error = %v ต้องเป็น ErrPhotoType", err)
		}
		entries, _ := os.ReadDir(wantDir)
		if len(entries) != 2 {
			t.Errorf("มีไฟล์ %d ไฟล์ ต้องเป็น 2", len(entries))
		}
	})

	t.Run("ลำดับซ้ำไม่เขียนทับ", func(t *testing.T) {
		if _, _, err := store.Save(loanID, 1, "image/jpeg", bytes.NewReader([]byte("x"))); err == nil {
			t.Fatal("ต้อง error เมื่อมีไฟล์ลำดับเดิมอยู่แล้ว")
		}
	})

	t.Run("rollback ลบทั้งโฟลเดอร์ของการยืม", func(t *testing.T) {
		if err := store.RemoveLoan(loanID); err != nil {
			t.Fatalf("RemoveLoan error = %v", err)
		}
		if _, err := os.Stat(wantDir); !os.IsNotExist(err) {
			t.Errorf("โฟลเดอร์ยังอยู่: %v", err)
		}
		if err := store.RemoveLoan(loanID); err != nil {
			t.Errorf("ลบซ้ำต้องไม่ error: %v", err)
		}
	})
}
