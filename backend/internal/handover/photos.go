package handover

import (
	"errors"
	"io"
	"net/http"
	"slices"

	"salmonnorway/backend/internal/config"
)

// ErrPhotoInvalid หมายถึงรูปไม่ครบตามจำนวนขั้นต่ำ ชนิดไฟล์ไม่ใช่ JPG/PNG หรือรูปใดใหญ่เกิน 10 MB (FR-HND-06, ASM-HND-07, Q-HND-10)
var ErrPhotoInvalid = errors.New("photos invalid")

// PhotoInput คือรูป 1 ไฟล์ที่เจ้าหน้าที่แนบ เปิดอ่านได้หลายครั้ง (ตรวจชนิดครั้งหนึ่ง เขียนลงดิสก์อีกครั้ง)
type PhotoInput struct {
	Open func() (io.ReadCloser, error)
}

// validPhoto คือรูปที่ผ่านการตรวจแล้ว พร้อมชนิดที่ตรวจจากเนื้อไฟล์
type validPhoto struct {
	open        func() (io.ReadCloser, error)
	contentType string
}

// StoredPhoto คือรูปที่เขียนลงดิสก์แล้ว รอบันทึกแถวในตาราง handover_photo
type StoredPhoto struct {
	Path        string
	ContentType string
	Size        int64
}

// validatePhotos ตรวจจำนวนขั้นต่ำ (config.HandoverMinPhotos) ชนิดไฟล์จากเนื้อไฟล์ ไม่ใช่จากนามสกุลหรือ header ที่ส่งมา
// (config.HandoverPhotoTypes) และขนาดต่อรูปไม่เกิน config.HandoverPhotoMaxBytes ไม่จำกัดจำนวนรูป
// (FR-HND-06, ASM-HND-07, Q-HND-10) ไม่ผ่านข้อใด = ErrPhotoInvalid
func validatePhotos(photos []PhotoInput) ([]validPhoto, error) {
	if len(photos) < config.HandoverMinPhotos {
		return nil, ErrPhotoInvalid
	}
	allowed := config.HandoverPhotoTypes()
	out := make([]validPhoto, 0, len(photos))
	for _, p := range photos {
		ct, size, err := inspectPhoto(p.Open)
		if err != nil {
			return nil, err
		}
		if !slices.Contains(allowed, ct) || size > config.HandoverPhotoMaxBytes {
			return nil, ErrPhotoInvalid
		}
		out = append(out, validPhoto{open: p.Open, contentType: ct})
	}
	return out, nil
}

// inspectPhoto อ่านไฟล์ 1 รอบ: หาชนิดจาก 512 byte แรกด้วย http.DetectContentType และนับขนาด (FR-HND-06)
// นับไม่เกิน HandoverPhotoMaxBytes+1 byte พอรู้ว่าเกินแล้วหยุด ไม่อ่านไฟล์ใหญ่จนจบ
func inspectPhoto(open func() (io.ReadCloser, error)) (contentType string, size int64, err error) {
	f, err := open()
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, err := io.ReadFull(f, buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", 0, err
	}
	rest, err := io.Copy(io.Discard, io.LimitReader(f, config.HandoverPhotoMaxBytes+1-int64(n)))
	if err != nil {
		return "", 0, err
	}
	return http.DetectContentType(buf[:n]), int64(n) + rest, nil
}

// photoSaver เขียนรูปที่ตรวจแล้วของการยืม 1 รายการลง PhotoStore และลบทิ้งเมื่อ transaction ล้มเหลว (FR-HND-06)
type photoSaver struct {
	store  *PhotoStore
	photos []validPhoto
}

// SavePhotos เขียนรูปทุกไฟล์ลง <root>/<loanID>/ ถ้าไฟล์ใดล้มเหลว ลบทั้งโฟลเดอร์ทิ้ง
func (s photoSaver) SavePhotos(loanID int) ([]StoredPhoto, error) {
	stored := make([]StoredPhoto, 0, len(s.photos))
	for i, p := range s.photos {
		f, err := p.open()
		if err != nil {
			s.RemovePhotos(loanID)
			return nil, err
		}
		path, size, err := s.store.Save(loanID, i+1, p.contentType, f)
		_ = f.Close()
		if err != nil {
			s.RemovePhotos(loanID)
			return nil, err
		}
		stored = append(stored, StoredPhoto{Path: path, ContentType: p.contentType, Size: size})
	}
	return stored, nil
}

// RemovePhotos ลบรูปของการยืม loanID ทั้งหมด ไม่ให้มีไฟล์ค้างเมื่อบันทึกไม่สำเร็จ (AC-HND-12)
func (s photoSaver) RemovePhotos(loanID int) {
	_ = s.store.RemoveLoan(loanID)
}
