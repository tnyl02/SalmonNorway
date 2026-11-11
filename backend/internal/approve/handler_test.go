package approve

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"salmonnorway/backend/internal/auth"
	"salmonnorway/backend/internal/config"
)

// ผู้ใช้จำลองจาก auth.NewMockProvider (รอ Q-SH-03 สำหรับโปรโตคอลจริง)
const (
	tokenStaffDept1   = "mock-staff-001"
	tokenStudentDept1 = "mock-student-001"
)

// newTestServer ประกอบ handler เข้ากับ auth.Middleware แบบเดียวกับ main.go (IF-IDP-01)
// ใช้ Querier ที่ให้ (transaction ของ test) ข้อมูลที่ test ใส่ไว้จึงเห็นได้และถูกยกเลิกทิ้งตอนจบ
func newTestServer(q Querier) http.Handler {
	return auth.Middleware(auth.NewMockProvider(), NewHandler(NewService(NewRepository(q))).Routes())
}

// getJSON ยิง GET ด้วย token ที่ให้ ถ้า token ว่างคือไม่ส่ง header เลย
func getJSON(srv http.Handler, token, path string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)
	return w
}

// AC-APV-11 (FR-APV-01 ก) ยิงผ่าน HTTP จริง
// Given มีคำขอ requested ของภาควิชาตนเองหลายคำขอ และคำขอของภาควิชาอื่น (approve_dev.sql)
// When  เจ้าหน้าที่ภาควิชา 1 เรียก GET /v1/approvals/requests
// Then  ได้ 200 รายการเฉพาะภาควิชาตนเอง เรียงตามเวลาที่ส่ง พร้อมข้อมูลที่ต้องแสดง is_first และ loadedAt
func TestAC_APV_11_HTTP(t *testing.T) {
	srv := newTestServer(devTx(t))

	w := getJSON(srv, tokenStaffDept1, PathPrefix+"requests")
	if w.Code != http.StatusOK {
		t.Fatalf("ต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var got PendingListDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.LoadedAt.IsZero() {
		t.Error("ต้องมี loadedAt ให้หน้าจอส่งกลับมาตอนพิจารณา (plan.md ข้อ 4)")
	}
	wantAssets := []string{"AP-TEST-0001", "AP-TEST-0002", "AP-TEST-0001", "AP-TEST-0003", "AP-TEST-0004"}
	if len(got.Items) != len(wantAssets) {
		t.Fatalf("ต้องได้ %d คำขอ แต่ได้ %d: %+v", len(wantAssets), len(got.Items), got.Items)
	}
	for i, a := range wantAssets {
		it := got.Items[i]
		if it.AssetCode != a {
			t.Errorf("แถว %d ต้องเป็น %s แต่ได้ %s", i, a, it.AssetCode)
		}
		if it.RequesterName == "" || it.ModelName == "" || it.Purpose == "" || it.StartAt.IsZero() || it.EndAt.IsZero() {
			t.Errorf("แถว %d ข้อมูลไม่ครบตาม FR-APV-01 (ก): %+v", i, it)
		}
	}
	if got.Items[2].IsFirst {
		t.Error("คำขอ B ของ AP-TEST-0001 ต้องไม่ใช่คำขอแรก (FR-APV-01 ข)")
	}
}

// T-05 เสร็จเมื่อ: endpoint ใต้ /v1/approvals/ ตอบ FORBIDDEN_ROLE เมื่อไม่ใช่เจ้าหน้าที่ (ASM-01, CON-TECH-03)
// และไม่มี token ได้ 401 UNAUTHENTICATED จาก auth.Middleware (IF-IDP-01)
func TestApprovals_RequiresStaff(t *testing.T) {
	srv := newTestServer(devTx(t))

	w := getJSON(srv, tokenStudentDept1, PathPrefix+"requests")
	if w.Code != http.StatusForbidden {
		t.Fatalf("นักศึกษาต้องได้ 403 แต่ได้ %d", w.Code)
	}
	var e ErrorDTO
	if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil {
		t.Fatal(err)
	}
	if e.Code != "FORBIDDEN_ROLE" || e.Message != config.ApprovalMessages["FORBIDDEN_ROLE"] {
		t.Errorf("ต้องตอบ {FORBIDDEN_ROLE, %q} แต่ได้ %+v", config.ApprovalMessages["FORBIDDEN_ROLE"], e)
	}
	if len(w.Body.Bytes()) > 0 && json.Valid(w.Body.Bytes()) {
		var raw map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &raw)
		if _, leaked := raw["items"]; leaked {
			t.Error("นักศึกษาต้องไม่ได้รายการคำขอกลับไป")
		}
	}

	if w := getJSON(srv, "", PathPrefix+"requests"); w.Code != http.StatusUnauthorized {
		t.Errorf("ไม่มี token ต้องได้ 401 แต่ได้ %d", w.Code)
	}
}

// postJSON ยิง POST ด้วย token และ body ที่ให้
func postJSON(srv http.Handler, token, path, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, r)
	return w
}

// AC-APV-01 (FR-APV-01, CON-DATA-01) ยิงผ่าน HTTP จริง
// Given เจ้าหน้าที่ภาควิชา 1 และคำขอแรกของอุปกรณ์ที่ "ว่าง" (AP-TEST-0002)
// When  POST /v1/approvals/requests/{id}/approve พร้อม loadedAt
// Then  ได้ 200 สถานะ approved และสถานะอุปกรณ์ยังเป็น "ว่าง"
func TestAC_APV_01_HTTP(t *testing.T) {
	tx := devTx(t)
	srv := newTestServer(tx)
	id := requestIDOf(t, tx, "AP-TEST-0002", "student-001")

	w := postJSON(srv, tokenStaffDept1, fmt.Sprintf("%srequests/%d/approve", PathPrefix, id), `{"loadedAt":"2026-11-12T03:00:00Z"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("อนุมัติสำเร็จต้องได้ 200 แต่ได้ %d body=%s", w.Code, w.Body.String())
	}
	var got ResultDTO
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.RequestID != id || got.Status != "approved" || got.DecidedAt.IsZero() || got.RejectedRequestIDs == nil {
		t.Errorf("คำตอบไม่ถูกต้อง: %+v", got)
	}

	var itemStatus string
	if err := tx.QueryRow(context.Background(), `
		SELECT i.status FROM borrow_request r JOIN equipment_item i ON i.id = r.item_id WHERE r.id = $1`, id).Scan(&itemStatus); err != nil {
		t.Fatal(err)
	}
	if itemStatus != "available" {
		t.Errorf("สถานะอุปกรณ์ต้องยังเป็น available แต่ได้ %q (CON-DATA-01)", itemStatus)
	}
}

// Q-APV-25: รหัสคำขอที่ไม่มีอยู่หรือไม่ใช่ตัวเลข ได้ 404 REQUEST_NOT_FOUND พร้อมข้อความจาก config
func TestApproveHTTP_NotFound(t *testing.T) {
	srv := newTestServer(devTx(t))
	for _, path := range []string{PathPrefix + "requests/999999/approve", PathPrefix + "requests/abc/approve"} {
		w := postJSON(srv, tokenStaffDept1, path, `{}`)
		if w.Code != http.StatusNotFound {
			t.Errorf("%s ต้องได้ 404 แต่ได้ %d", path, w.Code)
			continue
		}
		var e ErrorDTO
		_ = json.Unmarshal(w.Body.Bytes(), &e)
		if e.Code != "REQUEST_NOT_FOUND" || e.Message != config.ApprovalMessages["REQUEST_NOT_FOUND"] {
			t.Errorf("%s ต้องได้ {REQUEST_NOT_FOUND, ไม่พบคำขอนี้} แต่ได้ %+v", path, e)
		}
	}
}
