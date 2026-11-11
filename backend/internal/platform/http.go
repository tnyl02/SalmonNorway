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

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
