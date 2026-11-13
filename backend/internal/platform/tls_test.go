package platform

// test ของ HTTPS หลังบ้าน (T-50 ของ 004, AC-APV-19, NFR-SEC-01, CON-TECH-02)

import (
	"crypto/tls"
	"crypto/x509"
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

// tlsTestServer เปิดหลังบ้านจำลองแบบ HTTPS ด้วยใบรับรองที่ลงนามเอง คืนที่อยู่และ pool ของใบรับรองให้ client เชื่อ
func tlsTestServer(t *testing.T) (addr string, roots *x509.CertPool) {
	t.Helper()
	dir := t.TempDir()
	cert, key := filepath.Join(dir, "cert.pem"), filepath.Join(dir, "key.pem")
	if err := WriteSelfSignedCert(cert, key, []string{"127.0.0.1", "localhost"}); err != nil {
		t.Fatal(err)
	}
	pemBytes, err := os.ReadFile(cert)
	if err != nil {
		t.Fatal(err)
	}
	roots = x509.NewCertPool()
	roots.AppendCertsFromPEM(pemBytes)

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "api-ok") })
	go func() { _ = Serve(l, cert, key, h) }()
	t.Cleanup(func() { _ = l.Close() })
	return l.Addr().String(), roots
}

func tlsClient(roots *x509.CertPool, min, max uint16) *http.Client {
	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: roots, MinVersion: min, MaxVersion: max}},
	}
}

// AC-APV-19 (NFR-SEC-01, CON-TECH-02)
// Given เรียก API ผ่าน HTTP ที่ไม่เข้ารหัส หรือ TLS ต่ำกว่า 1.2
// When  ส่งคำขอใด ๆ
// Then  ระบบไม่ให้บริการผ่านช่องทางนั้น (และ TLS 1.2 ขึ้นไปใช้ได้ตามปกติ)
func TestAC_APV_19(t *testing.T) {
	addr, roots := tlsTestServer(t)

	t.Run("TLS 1.2 ใช้ได้", func(t *testing.T) {
		resp, err := tlsClient(roots, tls.VersionTLS12, tls.VersionTLS12).Get("https://" + addr + "/")
		if err != nil {
			t.Fatalf("TLS 1.2 ต้องเชื่อมต่อได้: %v", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK || string(body) != "api-ok" {
			t.Errorf("TLS 1.2 ต้องได้คำตอบของ API แต่ได้ %d %q", resp.StatusCode, body)
		}
	})

	t.Run("TLS 1.1 ไม่ให้บริการ", func(t *testing.T) {
		resp, err := tlsClient(roots, tls.VersionTLS10, tls.VersionTLS11).Get("https://" + addr + "/")
		if err == nil {
			resp.Body.Close()
			t.Fatal("TLS ต่ำกว่า 1.2 ต้องจับมือไม่สำเร็จ แต่เชื่อมต่อได้")
		}
	})

	t.Run("HTTP ธรรมดาไม่ให้บริการ", func(t *testing.T) {
		resp, err := (&http.Client{Timeout: 5 * time.Second}).Get("http://" + addr + "/")
		if err != nil {
			return // ปิดการเชื่อมต่อทิ้งก็ถือว่าไม่ให้บริการ
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK || strings.Contains(string(body), "api-ok") {
			t.Errorf("HTTP ธรรมดาต้องไม่ได้คำตอบของ API แต่ได้ %d %q", resp.StatusCode, body)
		}
	})
}

// T-50: ตั้งไฟล์ cert หรือ key อย่างเดียว ต้องไม่เปิด server (ไม่หลุดไปเป็น HTTP โดยไม่รู้ตัว)
func TestServe_TLSFilesIncomplete(t *testing.T) {
	if err := ListenAndServe("127.0.0.1:0", "cert.pem", "", http.NotFoundHandler()); !errors.Is(err, ErrTLSFilesIncomplete) {
		t.Errorf("ต้องได้ ErrTLSFilesIncomplete แต่ได้ %v", err)
	}
	if err := ListenAndServe("127.0.0.1:0", "", "key.pem", http.NotFoundHandler()); !errors.Is(err, ErrTLSFilesIncomplete) {
		t.Errorf("ต้องได้ ErrTLSFilesIncomplete แต่ได้ %v", err)
	}
}

// T-50: TLSConfig รับต่ำสุด TLS 1.2 (NFR-SEC-01)
func TestTLSConfig_MinVersion(t *testing.T) {
	if v := TLSConfig().MinVersion; v != tls.VersionTLS12 {
		t.Errorf("MinVersion ต้องเป็น TLS 1.2 แต่ได้ %x", v)
	}
}
