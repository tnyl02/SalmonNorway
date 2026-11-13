package platform

// HTTPS ของหลังบ้าน (CON-TECH-02 "HTTPS เท่านั้น", NFR-SEC-01 "TLS 1.2 ขึ้นไป", AC-APV-19)
// 004 รับทำเองตาม Q-APV-38 (ปิด 2569-11-13) แบบหลังบ้าน Go เปิด TLS เอง ไม่ใช้ reverse proxy (T-50 ของ 004)

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"

	"salmonnorway/backend/internal/config"
)

// TLSConfig คือการตั้งค่า TLS ของหลังบ้าน รับเฉพาะ TLS 1.2 ขึ้นไป (NFR-SEC-01)
func TLSConfig() *tls.Config {
	return &tls.Config{MinVersion: config.TLSMinVersion}
}

// ErrTLSFilesIncomplete เกิดเมื่อตั้งไฟล์ cert หรือ key อย่างใดอย่างหนึ่ง ไม่ครบคู่ ไม่เปิด server เพื่อไม่ให้หลุดไปเป็น HTTP โดยไม่รู้ตัว
var ErrTLSFilesIncomplete = errors.New("ต้องตั้งทั้ง TLS_CERT_FILE และ TLS_KEY_FILE หรือไม่ตั้งทั้งคู่")

// Serve เปิดรับ request บน listener ที่ให้มา
//   - มีทั้ง certFile และ keyFile: เปิดเฉพาะ HTTPS ด้วย TLSConfig() ไม่มีช่องทาง HTTP ธรรมดา (CON-TECH-02, AC-APV-19)
//   - ไม่มีทั้งคู่: HTTP แบบเดิม สำหรับเครื่องพัฒนาและ test (ทีมเลือกเอง ไม่ได้มาจาก spec)
//   - มีอย่างเดียว: ErrTLSFilesIncomplete
func Serve(l net.Listener, certFile, keyFile string, h http.Handler) error {
	srv := &http.Server{Handler: h, TLSConfig: TLSConfig()}
	switch {
	case certFile != "" && keyFile != "":
		log.Printf("serving HTTPS (TLS 1.2 ขึ้นไป) on %s", l.Addr())
		return srv.ServeTLS(l, certFile, keyFile)
	case certFile == "" && keyFile == "":
		log.Printf("serving HTTP on %s (ไม่ได้ตั้ง TLS_CERT_FILE / TLS_KEY_FILE ใช้กับเครื่องพัฒนาเท่านั้น)", l.Addr())
		return srv.Serve(l)
	default:
		return ErrTLSFilesIncomplete
	}
}

// ListenAndServe เปิด listener ที่ addr แล้วเรียก Serve
func ListenAndServe(addr, certFile, keyFile string, h http.Handler) error {
	if (certFile == "") != (keyFile == "") {
		return ErrTLSFilesIncomplete
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return Serve(l, certFile, keyFile, h)
}

// WriteSelfSignedCert สร้างใบรับรองที่ลงนามเองสำหรับเครื่องพัฒนาและ test (ไม่ใช่ของจริงสำหรับใช้งาน)
// hosts คือชื่อเครื่องหรือ IP ที่ใบรับรองใช้ได้ เช่น localhost, 127.0.0.1
// อายุใบรับรองและชนิดกุญแจทีมเลือกเอง ไม่ได้มาจาก spec
func WriteSelfSignedCert(certPath, keyPath string, hosts []string) error {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 62))
	if err != nil {
		return err
	}
	tmpl := x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{Organization: []string{"SalmonNorway dev"}},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		} else {
			tmpl.DNSNames = append(tmpl.DNSNames, h)
		}
	}
	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		return err
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}
	if err := os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o644); err != nil {
		return err
	}
	return os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), 0o600)
}
