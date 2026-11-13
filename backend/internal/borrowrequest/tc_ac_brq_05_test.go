package borrowrequest

// test ชุดนี้มาจากตาราง test-cases.md แถว TC-BRQ-05-1 ถึง TC-BRQ-05-3 (AC-BRQ-05, FR-BRQ-05 ข้อ ก)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test เหล่านี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"salmonnorway/backend/internal/auth"
)

// queueFixture เตรียม service กับ server ที่ใช้ร่วมกันในชุดนี้ และลบคำขอของผู้ขอที่ระบุเมื่อ test จบ
func queueFixture(t *testing.T, requesters ...string) (*Service, http.Handler) {
	t.Helper()
	svc := NewServiceWithClock(NewRepository(testPool(t)), testClock())
	t.Cleanup(func() {
		_, _ = testPool(t).Exec(context.Background(),
			`DELETE FROM borrow_request WHERE requester_id = ANY($1)`, requesters)
	})
	return svc, auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes())
}

// decodeQueue อ่านลำดับคำขอของอุปกรณ์ชิ้นหนึ่งจากคำตอบของ endpoint คิว
func decodeQueue(t *testing.T, w *httptest.ResponseRecorder) []queueItemDTO {
	t.Helper()
	var got []queueItemDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("อ่าน JSON ไม่ได้: %v body=%s", err, w.Body.String())
	}
	return got
}

// TC-BRQ-05-1 (ทางปกติ) มีคำขอค้าง 1 ใบ แล้วอีกคนส่งซ้อน
func TestTC_BRQ_05_1_QueueOrder(t *testing.T) {
	const assetCode = "TC-BRQ-05-1"
	const requesterA = "student-tc-05-1a"
	const requesterB = "student-tc-05-1b"
	insertTempItem(t, assetCode)
	svc, srv := queueFixture(t, requesterA, requesterB)

	// Given อุปกรณ์ชิ้นนี้มีคำขอค้างของนักศึกษา A อยู่ 1 คำขอ
	first, err := svc.CreateWithAttachment(studentCtx(requesterA), newTestRequest(assetCode), nil)
	if err != nil {
		t.Fatalf("คำขอใบแรกต้องสำเร็จ แต่ได้ error: %v", err)
	}

	// When นักศึกษา B ส่งคำขอของชิ้นเดียวกัน แล้วเจ้าหน้าที่เรียกคิว
	second, err := svc.CreateWithAttachment(studentCtx(requesterB), newTestRequest(assetCode), nil)
	// Then (1) คำขอของ B สำเร็จ เพราะส่งซ้อนกันได้ (DOM-REQ-01)
	if err != nil {
		t.Fatalf("(1) คำขอซ้อนของผู้อื่นต้องสำเร็จ แต่ได้ error: %v", err)
	}
	w := getQueue(t, srv, tokenStaff, assetCode)
	if w.Code != http.StatusOK {
		t.Fatalf("เจ้าหน้าที่ต้องเรียกคิวได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	queue := decodeQueue(t, w)

	// Then (2) ได้ 2 แถว
	if len(queue) != 2 {
		t.Fatalf("(2) ต้องได้คำขอ 2 ใบ แต่ได้ %d ใบ", len(queue))
	}
	// Then (3) คำขอของ A เป็นลำดับที่ 1
	if queue[0].RequestID != first.ID || queue[0].QueuePosition != 1 {
		t.Errorf("(3) ลำดับที่ 1 ต้องเป็นคำขอ %d แต่ได้ %d (ลำดับ %d)", first.ID, queue[0].RequestID, queue[0].QueuePosition)
	}
	// Then (4) คำขอของ B เป็นลำดับที่ 2
	if queue[1].RequestID != second.ID || queue[1].QueuePosition != 2 {
		t.Errorf("(4) ลำดับที่ 2 ต้องเป็นคำขอ %d แต่ได้ %d (ลำดับ %d)", second.ID, queue[1].RequestID, queue[1].QueuePosition)
	}
	// Then (5) เรียงตาม submitted_at
	if queue[0].SubmittedAt.After(queue[1].SubmittedAt) {
		t.Error("(5) ต้องเรียงตามเวลาที่ส่ง ใบแรกต้องมาก่อนใบที่สอง")
	}
}

// TC-BRQ-05-2 (ขอบ) ขยับจาก 1 คำขอใน Given เป็น 3 คำขอจาก 3 คน
func TestTC_BRQ_05_2_QueueThree(t *testing.T) {
	const assetCode = "TC-BRQ-05-2"
	requesters := []string{"student-tc-05-2a", "student-tc-05-2b", "student-tc-05-2c"}
	insertTempItem(t, assetCode)
	svc, srv := queueFixture(t, requesters...)

	// Given มีคำขอค้าง 3 คำขอจาก 3 คน ส่งเรียงกันไป
	var ids []int
	for _, requester := range requesters {
		req, err := svc.CreateWithAttachment(studentCtx(requester), newTestRequest(assetCode), nil)
		if err != nil {
			t.Fatalf("คำขอของ %s ต้องสำเร็จ แต่ได้ error: %v", requester, err)
		}
		ids = append(ids, req.ID)
	}

	// When เจ้าหน้าที่เรียกคิวของชิ้นนั้น
	w := getQueue(t, srv, tokenStaff, assetCode)
	if w.Code != http.StatusOK {
		t.Fatalf("ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	queue := decodeQueue(t, w)

	// Then (1) ได้ 3 แถว
	if len(queue) != 3 {
		t.Fatalf("(1) ต้องได้คำขอ 3 ใบ แต่ได้ %d ใบ", len(queue))
	}
	// Then (2) ลำดับเป็น 1, 2, 3 ตามเวลาที่ส่ง
	for i, item := range queue {
		if item.QueuePosition != i+1 {
			t.Errorf("(2) แถวที่ %d ต้องมีลำดับ %d แต่ได้ %d", i+1, i+1, item.QueuePosition)
		}
		if item.RequestID != ids[i] {
			t.Errorf("(2) ลำดับที่ %d ต้องเป็นคำขอ %d แต่ได้ %d", i+1, ids[i], item.RequestID)
		}
	}
}

// TC-BRQ-05-3 (ทางผิด) ผู้เรียกเป็นนักศึกษา ไม่ใช่เจ้าหน้าที่
// spec ไม่ได้ระบุรหัส ทีมตัดสิน 2569-11-13 ให้ใช้ FORBIDDEN_ROLE ที่มีอยู่แล้วใน CON-TECH-03
func TestTC_BRQ_05_3_QueueStudentForbidden(t *testing.T) {
	const assetCode = "TC-BRQ-05-3"
	const owner = "student-tc-05-3"
	insertTempItem(t, assetCode)
	svc, srv := queueFixture(t, owner)

	// Given อุปกรณ์ชิ้นนี้มีคำขอค้างของนักศึกษาคนอื่นอยู่
	other := newTestRequest(assetCode)
	other.RequesterName = "ผู้ขอรายอื่น ทดสอบ"
	if _, err := svc.CreateWithAttachment(studentCtx(owner), other, nil); err != nil {
		t.Fatalf("คำขอตั้งต้นต้องสำเร็จ แต่ได้ error: %v", err)
	}

	// When นักศึกษา (ไม่ใช่เจ้าหน้าที่) เรียกคิวของอุปกรณ์ชิ้นนั้น
	w := getQueue(t, srv, tokenStudentDept1, assetCode)

	// Then (1) ได้ 403 FORBIDDEN_ROLE
	if w.Code != http.StatusForbidden {
		t.Fatalf("(1) นักศึกษาต้องได้ 403 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	if code := errorCodeOf(t, w); code != "FORBIDDEN_ROLE" {
		t.Errorf("(1) รหัสต้องเป็น FORBIDDEN_ROLE แต่ได้ %q", code)
	}
	// Then (2) body ไม่มีข้อมูลคำขอของผู้อื่น (DOM-PRIV-01)
	if strings.Contains(w.Body.String(), other.RequesterName) {
		t.Error("(2) คำตอบต้องไม่มีชื่อผู้ขอรายอื่นติดไปด้วย")
	}
	if strings.Contains(w.Body.String(), "queuePosition") {
		t.Error("(2) คำตอบต้องไม่มีลำดับคิวติดไปด้วย")
	}
}
