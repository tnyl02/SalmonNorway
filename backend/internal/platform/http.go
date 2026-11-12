package platform

import (
	"encoding/json"
	"net/http"
)

// NewRouter สร้าง HTTP router ของหลังบ้าน
// ตอนนี้มีแค่ /healthz (T-01) endpoint ของฟีเจอร์ต่าง ๆ จะถูกเพิ่มใน task ของฟีเจอร์นั้น
func NewRouter(db Pinger) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz(db))
	return mux
}

// healthz ตอบ 200 เมื่อต่อฐานข้อมูลได้ และ 503 เมื่อต่อไม่ได้ (CON-TECH-01, CON-TECH-02: ใช้ตรวจ container ใน Docker)
func healthz(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(r.Context()); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "db_unavailable"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

// NoStore บอกเบราว์เซอร์และ proxy ไม่ให้เก็บ response ไว้ใช้ซ้ำ (NFR-DATA-01, plan.md ข้อ 6.2)
// สถานะอุปกรณ์ต้องคลาดเคลื่อนจากฐานข้อมูลไม่เกิน 30 วินาที หลังบ้านไม่ cache และไม่ให้ใคร cache ต่อ
// กดรีเฟรชแล้วจึงได้ข้อมูลล่าสุดจากฐานข้อมูลเสมอ ใช้ห่อ router ทั้งหมดใน cmd/server/main.go
func NoStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// MethodNotAllowedJSON เปลี่ยนคำตอบ 405 เป็น body {code: METHOD_NOT_ALLOWED, message} ตาม CON-TECH-03 (T-40)
// 405 มาจาก router ของ Go (http.ServeMux) ซึ่งเขียนข้อความธรรมดา "Method Not Allowed" เอง handler ของฟีเจอร์แก้ไม่ได้
// จึงครอบ router ทั้งหมดใน cmd/server/main.go ทุกฟีเจอร์ได้รูปแบบเดียวกัน (plan.md ข้อ 5 และข้อ 10 ของ 001)
// ดักเฉพาะการตอบ 405 ทิ้ง body เดิมแล้วเขียน {code, message} แทน header Allow ที่ router ตั้งไว้ยังอยู่
// คำตอบอื่น (200, 401, error JSON อื่น) ผ่านไปตามเดิมไม่ถูกแตะ
func MethodNotAllowedJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&methodNotAllowedWriter{ResponseWriter: w}, r)
	})
}

// methodNotAllowedWriter ส่งต่อทุกอย่างให้ ResponseWriter จริง ยกเว้นเมื่อ status เป็น 405 (CON-TECH-03, T-40)
type methodNotAllowedWriter struct {
	http.ResponseWriter
	wroteHeader bool // เขียน status ไปแล้วหรือยัง
	replaced    bool // ตอบ 405 แบบ {code, message} แล้ว body ที่ตามมาจาก router ให้ทิ้ง
}

func (m *methodNotAllowedWriter) WriteHeader(status int) {
	if m.wroteHeader || status != http.StatusMethodNotAllowed {
		m.wroteHeader = true
		m.ResponseWriter.WriteHeader(status)
		return
	}
	m.wroteHeader = true
	m.replaced = true
	m.ResponseWriter.Header().Del("Content-Length") // ความยาวของข้อความเดิมใช้ไม่ได้แล้ว
	WriteError(m.ResponseWriter, status, CodeMethodNotAllowed)
}

func (m *methodNotAllowedWriter) Write(p []byte) (int, error) {
	if !m.wroteHeader {
		m.WriteHeader(http.StatusOK) // Write โดยไม่เรียก WriteHeader ก่อน = 200 ตามกติกาของ net/http
	}
	if m.replaced {
		return len(p), nil // ทิ้งข้อความ "Method Not Allowed" ของ router
	}
	return m.ResponseWriter.Write(p)
}

// Unwrap ให้ http.ResponseController เข้าถึง ResponseWriter จริงได้ (เช่น Flush)
func (m *methodNotAllowedWriter) Unwrap() http.ResponseWriter { return m.ResponseWriter }

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
