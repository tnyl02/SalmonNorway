package borrowrequest

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"salmonnorway/backend/internal/config"
)

// Upload คือไฟล์ใบบันทึกข้อความที่ผู้ขอแนบมากับการยืมค้างคืน (FR-BRQ-10 ข้อ ก)
// เก็บเนื้อไฟล์ไว้ในหน่วยความจำได้ เพราะขนาดถูกจำกัดไว้ที่ config.AttachmentMaxBytes
type Upload struct {
	Filename    string
	ContentType string
	Size        int64
	Content     []byte
}

// attachmentField คือชื่อช่องไฟล์ใน multipart form ที่หน้าจอส่งมา (plan.md ข้อ 4)
const attachmentField = "attachment"

// parseMemoryBytes คือขนาดที่ยอมให้ multipart พักไว้ในหน่วยความจำก่อนเขียนลงไฟล์ชั่วคราว
// เป็นค่าทางเทคนิคของ net/http ไม่ได้มาจาก spec
const parseMemoryBytes = 1 << 20

// maxRequestBytes คือขนาดสูงสุดของทั้ง request เผื่อช่องข้อความอื่นจากขนาดไฟล์ที่ spec กำหนด
// กันไม่ให้ไฟล์ใหญ่เกินถูกอ่านเข้าเครื่องทั้งก้อนก่อนจะถูกปฏิเสธ (FR-BRQ-10 ข้อ ก)
const maxRequestBytes = config.AttachmentMaxBytes + 1<<20

// maxFileNameLen คือความยาวสูงสุดของชื่อไฟล์ที่เก็บลงดิสก์ เป็นค่าทางเทคนิค ไม่ได้มาจาก spec
const maxFileNameLen = 100

// isMultipart บอกว่า request นี้ส่งมาพร้อมไฟล์แนบหรือไม่ (plan.md ข้อ 4 "multipart เมื่อมีไฟล์แนบ")
func isMultipart(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data")
}

// parseCreateForm อ่านข้อมูลฟอร์มคำขอที่ส่งมาแบบ multipart/form-data พร้อมไฟล์แนบ (FR-BRQ-10 ข้อ ก)
// ช่องวันเวลาที่อ่านไม่ออกจะถูกปล่อยเป็นค่าว่าง แล้วให้การตรวจความครบถ้วนรายงานช่องนั้นเอง (FR-BRQ-03)
func parseCreateForm(r *http.Request) (CreateRequestDTO, *Upload, error) {
	if err := r.ParseMultipartForm(parseMemoryBytes); err != nil {
		return CreateRequestDTO{}, nil, err
	}

	in := CreateRequestDTO{
		AssetCode:     r.FormValue("assetCode"),
		RequesterName: r.FormValue("requesterName"),
		StudentCode:   r.FormValue("studentCode"),
		LoanType:      r.FormValue("loanType"),
		Purpose:       r.FormValue("purpose"),
		PurposeOther:  r.FormValue("purposeOther"),
		StartAt:       parseFormTime(r.FormValue("startAt")),
		EndAt:         parseFormTime(r.FormValue("endAt")),
	}

	file, header, err := r.FormFile(attachmentField)
	if errors.Is(err, http.ErrMissingFile) {
		return in, nil, nil // ไม่ได้แนบไฟล์มา กฎว่าต้องแนบหรือไม่อยู่ที่ attachmentProblems
	}
	if err != nil {
		return in, nil, err
	}
	defer file.Close()

	// อ่านเกินเพดานไป 1 ไบต์เพื่อให้รู้ว่าไฟล์ใหญ่เกิน โดยไม่ต้องอ่านทั้งก้อน (FR-BRQ-10 ข้อ ก)
	content, err := io.ReadAll(io.LimitReader(file, config.AttachmentMaxBytes+1))
	if err != nil {
		return in, nil, err
	}
	return in, &Upload{
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        int64(len(content)),
		Content:     content,
	}, nil
}

// parseFormTime แปลงวันเวลาแบบ RFC3339 ที่หน้าจอส่งมา ถ้าแปลงไม่ได้คืนค่าว่าง
func parseFormTime(v string) time.Time {
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(v))
	if err != nil {
		return time.Time{}
	}
	return t
}

// AttachmentStore เก็บไฟล์ใบบันทึกข้อความเป็นไฟล์บนเครื่องที่รันระบบ ไม่เก็บลงฐานข้อมูล (ASM-BRQ-07)
type AttachmentStore struct {
	dir string
}

// NewAttachmentStore สร้างที่เก็บไฟล์จากโฟลเดอร์ที่ config กำหนด (ASM-BRQ-07, plan.md ข้อ 5)
func NewAttachmentStore(dir string) *AttachmentStore {
	return &AttachmentStore{dir: dir}
}

// Save เขียนไฟล์เป็น <dir>/<requestID>-<ชื่อไฟล์ที่ทำให้ปลอดภัยแล้ว> แล้วคืน path ที่เก็บ (plan.md ข้อ 3)
func (s *AttachmentStore) Save(requestID int, up Upload) (string, error) {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(s.dir, fmt.Sprintf("%d-%s", requestID, safeFileName(up.Filename)))
	if err := os.WriteFile(path, up.Content, 0o600); err != nil {
		return "", err
	}
	return path, nil
}

// Remove ลบไฟล์ที่เพิ่งเก็บ ใช้ตอนย้อนกลับเมื่อบันทึกลงฐานข้อมูลไม่สำเร็จ
func (s *AttachmentStore) Remove(path string) error {
	return os.Remove(path)
}

// safeFileName ตัดส่วนที่เป็น path ออกและเหลือเฉพาะอักษรที่ปลอดภัย
// ชื่อไฟล์มาจากผู้ขอ จึงเชื่อไม่ได้ ชื่ออย่าง ../../etc/passwd ต้องไม่พาไฟล์ออกนอกโฟลเดอร์ที่กำหนด
func safeFileName(name string) string {
	name = filepath.Base(strings.ReplaceAll(name, `\`, "/"))

	var b strings.Builder
	lastWasFiller := false
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9',
			r == '.', r == '-', r == '_':
			b.WriteRune(r)
			lastWasFiller = false
		default:
			// อักษรอื่น เช่น ภาษาไทยหรือช่องว่าง แทนด้วย _ ตัวเดียวติดกัน ไม่ยืดยาวตามชื่อเดิม
			if !lastWasFiller {
				b.WriteRune('_')
				lastWasFiller = true
			}
		}
	}

	// ตัดจุดนำหน้าออกกันไฟล์ซ่อน และตัดตัวคั่นท้ายออก แต่ไม่ตัดจนนามสกุลไฟล์หาย
	out := strings.TrimRight(strings.TrimLeft(b.String(), "."), "._")
	if out == "" {
		out = "attachment"
	}
	if len(out) > maxFileNameLen {
		out = out[len(out)-maxFileNameLen:] // เก็บท้ายไว้เพื่อให้นามสกุลไฟล์ยังอยู่
	}
	return out
}
