package cancelrequest

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"salmonnorway/backend/internal/auth"
)

// fakeNotifier เก็บ Email ที่ถูกส่งไว้ตรวจใน test ตั้ง err เพื่อจำลองส่งไม่สำเร็จ (T-28)
type fakeNotifier struct {
	mu   sync.Mutex
	sent []CancelNotice
	err  error
}

func (f *fakeNotifier) NotifyCancelled(_ context.Context, n CancelNotice) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sent = append(f.sent, n)
	return f.err
}

// notifyServer ประกอบ handler แบบเดียวกับ tcServer แต่ใช้ Notifier ที่ให้มา และ Email เจ้าหน้าที่ภาควิชา 1 จาก staff_contact
func notifyServer(t *testing.T, n Notifier) http.Handler {
	t.Helper()
	tcServer(t) // ล้างคำขอเดิม
	pool := tcPool(t)
	ctx := context.Background()
	if _, err := pool.Exec(ctx, `DELETE FROM staff_contact`); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO staff_contact (department_id, email) VALUES (1, 'staff-it@example.test')`); err != nil {
		t.Fatal(err)
	}
	svc := NewService(NewRepository(pool)).WithNotifier(n)
	mux := http.NewServeMux()
	mux.Handle(PathPrefixRequests, auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes()))
	return mux
}

// AC-CXL-15 (FR-CXL-09, IF-NOT-01, ASM-10) T-27
//
//	Given คำขอของ A เป็นอุปกรณ์ของภาควิชา 1 (CS-NB-0001 Notebook ตัวอย่างรุ่น A) ภาควิชา 1 มี Email เจ้าหน้าที่
//	When  A ยกเลิกคำขอนั้นสำเร็จพร้อมเหตุผล "ไม่ต้องใช้อุปกรณ์แล้ว"
//	Then  Notifier ได้ 1 ข้อความถึง Email เจ้าหน้าที่ภาควิชา 1 มีผู้ยกเลิก ชื่อรุ่น รหัสครุภัณฑ์ เวลาที่ยกเลิก และเหตุผล
//	      และส่งหลัง commit (เวลาในข้อความตรงกับ cancelled_at ในฐานข้อมูล)
func TestAC_CXL_15_NotifyDepartmentStaff(t *testing.T) {
	fake := &fakeNotifier{}
	srv := notifyServer(t, fake)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())

	tcWantStatus(t, tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA), http.StatusOK)

	if len(fake.sent) != 1 {
		t.Fatalf("ต้องส่ง 1 ข้อความ ได้ %d", len(fake.sent))
	}
	n := fake.sent[0]
	row := tcRow(t, id)
	if len(n.To) != 1 || n.To[0] != "staff-it@example.test" {
		t.Errorf("ผู้รับ = %v ต้องเป็น staff-it@example.test", n.To)
	}
	if n.CancelledBy != tcUserA || n.ModelName != "Notebook ตัวอย่างรุ่น A" || n.AssetCode != "CS-NB-0001" || n.Reason != "ไม่ต้องใช้อุปกรณ์แล้ว" {
		t.Errorf("ข้อความ = %+v ต้องมีผู้ยกเลิก ชื่อรุ่น รหัสครุภัณฑ์ และเหตุผลครบ", n)
	}
	// ผู้ยกเลิกแสดงทั้งชื่อ รหัสนักศึกษา (จากคำขอ) และรหัสผู้ใช้ (ทีมตัดสิน 2569-11-12)
	if got := FormatCanceller(n); got != "ผู้ทดสอบ (รหัสนักศึกษา TEST-000, ผู้ใช้ student-001)" {
		t.Errorf("FormatCanceller = %q", got)
	}
	if row.CancelledAt == nil || !n.CancelledAt.Equal(*row.CancelledAt) {
		t.Errorf("เวลาในข้อความ %v ต้องตรงกับ cancelled_at %v", n.CancelledAt, row.CancelledAt)
	}
}

// FR-CXL-09: ยกเลิกไม่สำเร็จต้องไม่ส่ง Email
func TestNotify_NotSentWhenNotCancelled(t *testing.T) {
	fake := &fakeNotifier{}
	srv := notifyServer(t, fake)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "approved", time.Now())

	tcWantStatus(t, tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA), http.StatusConflict)
	tcWantStatus(t, tcCancel(t, srv, id, "   ", tcTokenA), http.StatusUnprocessableEntity)

	if len(fake.sent) != 0 {
		t.Errorf("ยกเลิกไม่สำเร็จต้องไม่ส่ง Email ได้ %d ข้อความ", len(fake.sent))
	}
}

// ASM-10: อุปกรณ์ของภาควิชาที่ไม่มี Email เจ้าหน้าที่ ยกเลิกสำเร็จตามปกติ ไม่เรียก Notifier
// และนับเป็นส่งไม่สำเร็จ บันทึก notify error เหมือน AC-CXL-16 (ErrNoRecipients ทีมตัดสิน 2569-11-12, T-32)
func TestNotify_NoStaffEmailStillCancelled(t *testing.T) {
	fake := &fakeNotifier{err: errors.New("ไม่ควรถูกเรียก")}
	srv := notifyServer(t, fake)
	id := tcInsert(t, tcUserA, 8, "normal", "requested", time.Now()) // XX-NB-0001 ภาควิชา 2
	mark := tcAuditMark(t)

	tcWantStatus(t, tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA), http.StatusOK)

	if len(fake.sent) != 0 {
		t.Errorf("ภาควิชาที่ไม่มี Email ต้องไม่เรียก Notifier ได้ %d", len(fake.sent))
	}
	if s := tcRow(t, id).Status; s != "cancelled" {
		t.Errorf("สถานะ = %s ต้องเป็น cancelled", s)
	}
	wantNotifyError(t, tcAuditSince(t, mark), id)
}

// wantNotifyError ตรวจว่า audit_log มีแถว cancel ok ตามด้วย notify error ของคำขอ id (FR-CXL-05 ตาราง, FR-CXL-09, T-32)
func wantNotifyError(t *testing.T, rows []tcAuditRow, id int) {
	t.Helper()
	if len(rows) != 2 {
		t.Fatalf("audit_log มี %d แถว ต้องมี 2 (cancel ok, notify error): %+v", len(rows), rows)
	}
	if r := rows[0]; r.Action != "cancel" || r.Outcome != "ok" {
		t.Errorf("แถวแรก = %+v ต้องเป็น cancel ok", r)
	}
	r := rows[1]
	if r.Action != "notify" || r.Outcome != "error" || r.RequestID == nil || *r.RequestID != id ||
		r.Query["for"] != "cancel" || r.ActorID != tcUserA || !r.HasTime {
		t.Errorf("แถวที่สอง = %+v ต้องเป็น notify error request_id %d query for=cancel ของ %s", r, id, tcUserA)
	}
}

// FR-CXL-09: เวลาใน Email เป็นเวลาไทย ปี พ.ศ. ไม่มีวินาที เหมือนหน้าเว็บ
func TestFormatNoticeTime(t *testing.T) {
	cases := map[string]string{
		"2026-11-11T18:25:39Z":      "2569-11-12 01:25", // ข้ามวันเมื่อเป็นเวลาไทย
		"2026-11-10T03:30:59Z":      "2569-11-10 10:30", // ตัดวินาที ไม่ปัด
		"2026-11-12T01:25:39+07:00": "2569-11-12 01:25", // เวลาเดียวกับกรณีแรก
	}
	for in, want := range cases {
		ts, err := time.Parse(time.RFC3339, in)
		if err != nil {
			t.Fatal(err)
		}
		if got := FormatNoticeTime(ts); got != want {
			t.Errorf("FormatNoticeTime(%s) = %q ต้องเป็น %q", in, got, want)
		}
	}
}

// ทีมตัดสิน 2569-11-12: ไม่มีผู้รับ = ส่งไม่สำเร็จ (ErrNoRecipients) และ Notifier error ส่งต่อออกมา
func TestNotify_ReportsFailure(t *testing.T) {
	failing := errors.New("smtp down")
	svc := NewService(nil).WithNotifier(&fakeNotifier{err: failing})
	if err := svc.notify(context.Background(), CancelNotice{RequestID: 1}); !errors.Is(err, ErrNoRecipients) {
		t.Errorf("ไม่มีผู้รับ ต้องได้ ErrNoRecipients ได้ %v", err)
	}
	if err := svc.notify(context.Background(), CancelNotice{RequestID: 1, To: []string{"a@example.test"}}); !errors.Is(err, failing) {
		t.Errorf("Notifier error ต้องส่งต่อออกมา ได้ %v", err)
	}
	ok := NewService(nil).WithNotifier(&fakeNotifier{})
	if err := ok.notify(context.Background(), CancelNotice{RequestID: 1, To: []string{"a@example.test"}}); err != nil {
		t.Errorf("ส่งสำเร็จต้องได้ nil ได้ %v", err)
	}
}

// AC-CXL-16 (FR-CXL-09) T-28, ส่วนบันทึก audit_log notify error T-32
//
//	Given ระบบแจ้งเตือน Email ไม่ตอบสนอง (Notifier ปลอมคืน error)
//	When  A ยกเลิกคำขอและยืนยัน
//	Then  ได้ 200 {requestId, status:"cancelled"} สถานะในฐานข้อมูลเป็น cancelled พร้อมเหตุผล
//	      Notifier ถูกเรียก 1 ครั้ง (ส่งหลัง commit) และคำขอหายจากรายการ
//	      และ audit_log มีแถว notify outcome error พร้อม request_id ของคำขอนั้น
func TestAC_CXL_16_NotifyFailureStillCancelled(t *testing.T) {
	fake := &fakeNotifier{err: errors.New("smtp: connection refused")}
	srv := notifyServer(t, fake)
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())
	mark := tcAuditMark(t)

	rec := tcCancel(t, srv, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)

	tcWantStatus(t, rec, http.StatusOK)
	if body := tcDecode[tcCancelBody](t, rec); body.RequestID != id || body.Status != "cancelled" {
		t.Errorf("body = %+v ต้องเป็น {requestId:%d status:cancelled}", body, id)
	}
	if len(fake.sent) != 1 {
		t.Errorf("Notifier ต้องถูกเรียก 1 ครั้ง ได้ %d", len(fake.sent))
	}
	r := tcRow(t, id)
	if r.Status != "cancelled" || r.CancelReason == nil || *r.CancelReason != "ไม่ต้องใช้อุปกรณ์แล้ว" || r.CancelledAt == nil {
		t.Errorf("แถว = %+v ต้องเป็น cancelled พร้อมเหตุผลและเวลา (ไม่ย้อนการยกเลิก)", r)
	}
	list := tcDecode[tcListBody](t, tcGet(t, tcHandlerOnly(t), tcPathMine, tcTokenA))
	if len(list.Items) != 0 {
		t.Errorf("คำขอที่ยกเลิกแล้วต้องไม่อยู่ในรายการ ได้ %+v", list.Items)
	}
	wantNotifyError(t, tcAuditSince(t, mark), id)
}

// hangingNotifier จำลองระบบ Email ที่ค้างไม่ตอบ และไม่สนใจ context (กรณีแย่ที่สุด) ปล่อยเมื่อ test จบ
type hangingNotifier struct {
	release chan struct{}
	called  chan struct{}
}

func (h *hangingNotifier) NotifyCancelled(_ context.Context, _ CancelNotice) error {
	close(h.called)
	<-h.release
	return nil
}

// AC-CXL-16 (FR-CXL-09, NFR-PERF-CXL-01, ASM-10) T-34 ระบบ Email ค้าง
//
//	Given ระบบแจ้งเตือน Email ไม่ตอบสนอง (ค้างไม่ตอบ)
//	When  A ยกเลิกคำขอและยืนยัน
//	Then  ได้ 200 ภายในเวลาที่จำกัด (ไม่รอ Email ไม่สิ้นสุด) สถานะเป็น cancelled
//	      และ audit_log มีแถว notify error ของคำขอนั้น (ส่งไม่ทันนับเป็นส่งไม่สำเร็จ)
func TestAC_CXL_16_NotifyHangsStillCancelled(t *testing.T) {
	hang := &hangingNotifier{release: make(chan struct{}), called: make(chan struct{})}
	t.Cleanup(func() { close(hang.release) })
	notifyServer(t, hang) // ล้างคำขอเดิมและใส่ Email เจ้าหน้าที่ภาควิชา 1
	svc := NewService(NewRepository(tcPool(t))).WithNotifier(hang)
	svc.notifyTimeout = 200 * time.Millisecond // แทน config.NotifyTimeout ให้ test เร็ว
	mux := http.NewServeMux()
	mux.Handle(PathPrefixRequests, auth.Middleware(auth.NewMockProvider(), NewHandler(svc).Routes()))
	id := tcInsert(t, tcUserA, tcItemAvailable, "normal", "requested", time.Now())
	mark := tcAuditMark(t)

	start := time.Now()
	rec := tcCancel(t, mux, id, "ไม่ต้องใช้อุปกรณ์แล้ว", tcTokenA)
	elapsed := time.Since(start)

	tcWantStatus(t, rec, http.StatusOK)
	select {
	case <-hang.called:
	default:
		t.Errorf("ต้องเรียกตัวส่ง Email")
	}
	if elapsed > 3*time.Second {
		t.Errorf("ใช้เวลา %v ต้องหยุดรอ Email เมื่อครบ notifyTimeout", elapsed)
	}
	if s := tcRow(t, id).Status; s != "cancelled" {
		t.Errorf("สถานะ = %s ต้องเป็น cancelled", s)
	}
	wantNotifyError(t, tcAuditSince(t, mark), id)
}

// tcHandlerOnly ประกอบ handler แบบเดียวกับ tcServer แต่ไม่ล้างข้อมูล ใช้ตรวจผลหลังทำรายการ
func tcHandlerOnly(t *testing.T) http.Handler {
	t.Helper()
	h := NewHandler(NewService(NewRepository(tcPool(t))))
	mux := http.NewServeMux()
	mux.Handle(PathPrefixMine, auth.Middleware(auth.NewMockProvider(), h.Routes()))
	return mux
}
