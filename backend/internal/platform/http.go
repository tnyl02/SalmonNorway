package platform

import (
	"encoding/json"
	"net/http"
)

// Routes คือฟังก์ชันที่ฟีเจอร์ใช้ผูก endpoint ของตัวเองเข้ากับ router
// ฟีเจอร์เป็นคนห่อ auth.Middleware เอง platform จึงไม่ต้อง import package ของฟีเจอร์ (เพิ่มโดย 006 T-03)
type Routes func(mux *http.ServeMux)

// NewRouter สร้าง HTTP router ของหลังบ้าน: /healthz (001 T-01) และ endpoint ของฟีเจอร์ที่ส่งมาใน routes
func NewRouter(db Pinger, routes ...Routes) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz(db))
	for _, register := range routes {
		register(mux)
	}
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

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
