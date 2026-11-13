package platform

// test จาก specs/004-approve-request/test-cases.md แถว TC-APV-19-* (AC-APV-19, NFR-SEC-01, CON-TECH-02)
// อยู่ใน package platform เพราะ HTTPS เป็นของกลาง (T-50 ของ 004) ห้ามแก้ให้ผ่าน

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// tcGetAPI เรียก API แล้วคืน body ถ้าได้คำตอบ (err != nil = ไม่ได้คำตอบเลย)
func tcGetAPI(c *http.Client, url string) (int, string, error) {
	resp, err := c.Get(url)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b), nil
}

// TC-APV-19-1 (AC-APV-19 ปกติ)
func TestTC_APV_19_1_TLS12Served(t *testing.T) {
	// Given หลังบ้านเปิดด้วยไฟล์ cert และ key ที่ test สร้าง client เชื่อใบรับรองนั้น
	addr, roots := tlsTestServer(t)
	// When client ใช้ TLS 1.2 เรียก API
	code, body, err := tcGetAPI(tlsClient(roots, tls.VersionTLS12, tls.VersionTLS12), "https://"+addr+"/")
	// Then (1) จับมือสำเร็จ
	if err != nil {
		t.Fatalf("(1) TLS 1.2 ต้องจับมือสำเร็จ: %v", err)
	}
	// Then (2) ได้คำตอบของ API ตามปกติ
	if code != http.StatusOK || body != "api-ok" {
		t.Errorf("(2) ต้องได้คำตอบของ API แต่ได้ %d %q", code, body)
	}
}

// TC-APV-19-2 (AC-APV-19 ขอบ)
func TestTC_APV_19_2_TLS13Served(t *testing.T) {
	// Given เหมือน 19-1
	addr, roots := tlsTestServer(t)
	// When client ใช้ TLS 1.3 เรียก API
	code, body, err := tcGetAPI(tlsClient(roots, tls.VersionTLS13, tls.VersionTLS13), "https://"+addr+"/")
	// Then (1) จับมือสำเร็จ
	if err != nil {
		t.Fatalf("(1) TLS 1.3 ต้องจับมือสำเร็จ: %v", err)
	}
	// Then (2) ได้คำตอบของ API
	if code != http.StatusOK || body != "api-ok" {
		t.Errorf("(2) ต้องได้คำตอบของ API แต่ได้ %d %q", code, body)
	}
}

// TC-APV-19-3 (AC-APV-19 ขอบ)
func TestTC_APV_19_3_TLS11Refused(t *testing.T) {
	// Given client จำกัดรุ่นสูงสุดที่ TLS 1.1
	addr, roots := tlsTestServer(t)
	// When client เรียก API
	_, body, err := tcGetAPI(tlsClient(roots, tls.VersionTLS10, tls.VersionTLS11), "https://"+addr+"/")
	// Then (1) จับมือไม่สำเร็จ ไม่ได้คำตอบของ API
	if err == nil {
		t.Errorf("(1) TLS 1.1 ต้องจับมือไม่สำเร็จ แต่ได้คำตอบ %q", body)
	}
}

// TC-APV-19-4 (AC-APV-19 ผิด)
func TestTC_APV_19_4_TLS10Refused(t *testing.T) {
	// Given client จำกัดรุ่นที่ TLS 1.0
	addr, roots := tlsTestServer(t)
	// When client เรียก API
	_, body, err := tcGetAPI(tlsClient(roots, tls.VersionTLS10, tls.VersionTLS10), "https://"+addr+"/")
	// Then (1) จับมือไม่สำเร็จ ไม่ได้คำตอบของ API
	if err == nil {
		t.Errorf("(1) TLS 1.0 ต้องจับมือไม่สำเร็จ แต่ได้คำตอบ %q", body)
	}
}

// TC-APV-19-5 (AC-APV-19 ผิด)
func TestTC_APV_19_5_PlainHTTPRefused(t *testing.T) {
	// Given client ส่ง HTTP ธรรมดาไปที่ port HTTPS ของหลังบ้าน
	addr, _ := tlsTestServer(t)
	// When client เรียก API
	code, body, err := tcGetAPI(&http.Client{Timeout: 5 * time.Second}, "http://"+addr+"/")
	// Then (1) ไม่ได้คำตอบของ API (ได้ 400 จาก server หรือถูกตัดการเชื่อมต่อ)
	if err == nil && (code == http.StatusOK || strings.Contains(body, "api-ok")) {
		t.Errorf("(1) HTTP ธรรมดาต้องไม่ได้คำตอบของ API แต่ได้ %d %q", code, body)
	}
	// Then (2) ไม่มี port HTTP ธรรมดาอีกช่องทาง: ใน deploy/docker-compose.yml หลังบ้านเปิด port เดียว และตั้ง cert กับ key ไว้
	compose, err := os.ReadFile(filepath.Join("..", "..", "..", "deploy", "docker-compose.yml"))
	if err != nil {
		t.Fatalf("(2) อ่าน deploy/docker-compose.yml ไม่ได้: %v", err)
	}
	backend := string(compose)
	i := strings.Index(backend, "\n  backend:")
	if i < 0 {
		t.Fatal("(2) docker-compose.yml ต้องมี service backend")
	}
	backend = backend[i:]
	if j := strings.Index(backend, "\nvolumes:"); j > 0 {
		backend = backend[:j]
	}
	for _, want := range []string{"TLS_CERT_FILE:", "TLS_KEY_FILE:", `HTTP_ADDR: ":8443"`, `- "8443:8443"`} {
		if !strings.Contains(backend, want) {
			t.Errorf("(2) service backend ต้องมี %s", want)
		}
	}
	if n := strings.Count(backend, `- "`); n != 1 {
		t.Errorf("(2) service backend ต้องเปิด port เดียว (HTTPS) แต่พบ %d รายการใน ports", n)
	}
}

// TC-APV-19-6 (AC-APV-19 ผิด)
func TestTC_APV_19_6_IncompleteTLSFiles(t *testing.T) {
	// Given ตั้งไฟล์ cert อย่างเดียว ไม่ตั้ง key (และกลับกัน)
	for _, files := range [][2]string{{"cert.pem", ""}, {"", "key.pem"}} {
		// When เริ่มหลังบ้าน
		err := ListenAndServe("127.0.0.1:0", files[0], files[1], http.NotFoundHandler())
		// Then (1) ไม่เปิดรับ request ได้ ErrTLSFilesIncomplete
		if !errors.Is(err, ErrTLSFilesIncomplete) {
			t.Errorf("(1) cert=%q key=%q ต้องได้ ErrTLSFilesIncomplete แต่ได้ %v", files[0], files[1], err)
		}
	}
	// Serve กับ listener ที่เปิดไว้แล้วก็ต้องไม่ให้บริการเช่นกัน
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if err := Serve(l, "cert.pem", "", http.NotFoundHandler()); !errors.Is(err, ErrTLSFilesIncomplete) {
		t.Errorf("(1) Serve ต้องได้ ErrTLSFilesIncomplete แต่ได้ %v", err)
	}
}

// TC-APV-19-7 (AC-APV-19 ผิด)
func TestTC_APV_19_7_NoTLSFilesDevMode(t *testing.T) {
	// Given หลังบ้านเริ่มโดยไม่ตั้งไฟล์ cert และ key (แบบเครื่องพัฒนาและ test ของฟีเจอร์อื่น)
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "api-ok") })
	go func() { _ = Serve(l, "", "", h) }()
	t.Cleanup(func() { _ = l.Close() })

	// When client เรียก API ผ่าน HTTP ธรรมดา
	code, body, err := tcGetAPI(&http.Client{Timeout: 5 * time.Second}, "http://"+l.Addr().String()+"/")

	// Then (1) ได้คำตอบของ API ผ่าน HTTP ตามปกติ (ASM-13 แถวแก้ตามคำตอบทีม 2569-11-13 ทาง ก)
	if err != nil {
		t.Fatalf("(1) ไม่ตั้ง cert ต้องตอบผ่าน HTTP ได้ แต่ได้ error %v", err)
	}
	if code != http.StatusOK || body != "api-ok" {
		t.Errorf("(1) ต้องได้คำตอบของ API แต่ได้ %d %q", code, body)
	}
}
