// โปรแกรมสร้างใบรับรอง HTTPS ที่ลงนามเองสำหรับเครื่องพัฒนา (T-50 ของ 004, CON-TECH-02, NFR-SEC-01)
// ใช้: cd backend && go run ./cmd/devcert          -> เขียน ../deploy/certs/cert.pem และ key.pem
// ไม่ใช่ใบรับรองจริงสำหรับใช้งาน เบราว์เซอร์จะเตือนว่าไม่น่าเชื่อถือ (ทีมเลือกเอง ไม่ได้มาจาก spec)
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"salmonnorway/backend/internal/platform"
)

func main() {
	out := flag.String("out", "../deploy/certs", "โฟลเดอร์ที่เขียน cert.pem และ key.pem")
	hosts := flag.String("hosts", "localhost,127.0.0.1,backend", "ชื่อเครื่องหรือ IP ที่ใบรับรองใช้ได้ คั่นด้วย ,")
	flag.Parse()
	if err := os.MkdirAll(*out, 0o755); err != nil {
		log.Fatal(err)
	}
	cert, key := filepath.Join(*out, "cert.pem"), filepath.Join(*out, "key.pem")
	if err := platform.WriteSelfSignedCert(cert, key, strings.Split(*hosts, ",")); err != nil {
		log.Fatal(err)
	}
	log.Printf("เขียน %s และ %s แล้ว", cert, key)
}
