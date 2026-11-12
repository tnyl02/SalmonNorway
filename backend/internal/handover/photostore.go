package handover

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"salmonnorway/backend/internal/config"
)

// ErrPhotoType หมายถึงชนิดไฟล์ที่ที่เก็บไม่รู้จัก (การตรวจชนิดจริงจากเนื้อไฟล์อยู่ใน T-15)
var ErrPhotoType = errors.New("unsupported photo type")

// นามสกุลไฟล์ตามชนิดที่ FR-HND-06 รับ (ASM-HND-07) ตั้งชื่อไฟล์เองเสมอ ไม่ใช้ชื่อไฟล์ที่ผู้ใช้ส่งมา
var photoExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

// PhotoStore เก็บไฟล์รูปสภาพอุปกรณ์ตอนส่งมอบบนดิสก์ (FR-HND-06, plan.md Q-P1: Docker volume ทีมเลือกเอง ไม่ได้มาจาก spec)
// โครงโฟลเดอร์: <root>/<loanID>/<ลำดับ>.<นามสกุล> ฐานข้อมูลเก็บแค่ path ในตาราง handover_photo
type PhotoStore struct {
	root string
}

// NewPhotoStore สร้างที่เก็บรูปที่โฟลเดอร์ root
func NewPhotoStore(root string) *PhotoStore {
	return &PhotoStore{root: root}
}

// DefaultPhotoRoot คือโฟลเดอร์ของ 005 ใต้โฟลเดอร์รากของรูป (config.PhotoDir) ฟีเจอร์อื่นใช้โฟลเดอร์ข้างกันได้
func DefaultPhotoRoot() string {
	return filepath.Join(config.PhotoDir(), "handover")
}

// Save เขียนรูปลำดับที่ index ของการยืม loanID คืน path ของไฟล์และจำนวน byte ที่เขียน (FR-HND-06)
// ถ้าเขียนไม่สำเร็จ ลบไฟล์ที่เขียนค้างทิ้ง
func (s *PhotoStore) Save(loanID, index int, contentType string, r io.Reader) (path string, size int64, err error) {
	ext, ok := photoExt[contentType]
	if !ok {
		return "", 0, fmt.Errorf("%w: %s", ErrPhotoType, contentType)
	}
	dir := s.loanDir(loanID)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", 0, err
	}
	path = filepath.Join(dir, strconv.Itoa(index)+ext)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o640)
	if err != nil {
		return "", 0, err
	}
	size, err = io.Copy(f, r)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(path)
		return "", 0, err
	}
	return path, size, nil
}

// RemoveLoan ลบรูปทั้งหมดของการยืม loanID ใช้เมื่อ transaction การส่งมอบล้มเหลว
// ไม่ให้มีไฟล์ค้างบนดิสก์ (FR-HND-06, AC-HND-12) ถ้าไม่มีโฟลเดอร์ถือว่าสำเร็จ
func (s *PhotoStore) RemoveLoan(loanID int) error {
	return os.RemoveAll(s.loanDir(loanID))
}

func (s *PhotoStore) loanDir(loanID int) string {
	return filepath.Join(s.root, strconv.Itoa(loanID))
}
