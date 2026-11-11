package handover

import (
	"errors"
	"io"
	"net/http"
	"slices"

	"salmonnorway/backend/internal/config"
)

// ErrPhotoInvalid หมายถึงรูปไม่ครบตามจำนวนขั้นต่ำ หรือชนิดไฟล์ไม่ใช่ JPG/PNG (FR-HND-06, ASM-HND-07)
// ขนาดไฟล์ตรวจใน T-28 (รอ Q-HND-10)
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

// validatePhotos ตรวจจำนวนขั้นต่ำ (config.HandoverMinPhotos) และชนิดไฟล์จากเนื้อไฟล์ ไม่ใช่จากนามสกุลหรือ header ที่ส่งมา
// (config.HandoverPhotoTypes) (FR-HND-06, ASM-HND-07) ไม่ผ่านข้อใด = ErrPhotoInvalid
func validatePhotos(photos []PhotoInput) ([]validPhoto, error) {
	if len(photos) < config.HandoverMinPhotos {
		return nil, ErrPhotoInvalid
	}
	allowed := config.HandoverPhotoTypes()
	out := make([]validPhoto, 0, len(photos))
	for _, p := range photos {
		ct, err := sniffContentType(p.Open)
		if err != nil {
			return nil, err
		}
		if !slices.Contains(allowed, ct) {
			return nil, ErrPhotoInvalid
		}
		out = append(out, validPhoto{open: p.Open, contentType: ct})
	}
	return out, nil
}

// sniffContentType อ่าน 512 byte แรกของไฟล์แล้วหาชนิดด้วย http.DetectContentType
func sniffContentType(open func() (io.ReadCloser, error)) (string, error) {
	f, err := open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, err := io.ReadFull(f, buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", err
	}
	return http.DetectContentType(buf[:n]), nil
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
