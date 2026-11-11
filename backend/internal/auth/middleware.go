package auth

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Middleware บังคับให้ทุก request ยืนยันตัวตนก่อน (IF-IDP-01, NFR-SEC-01 "ต้องผ่านการยืนยันตัวตนก่อนทุกครั้ง")
//   - ไม่มี token หรือ token ใช้ไม่ได้ -> 401 UNAUTHENTICATED (plan ข้อ 5, หน้าจอใช้ต่อใน AC-EQS-09)
//   - ยืนยันตัวตนได้ แต่ไม่มีภาควิชาหรือบทบาทไม่รู้จัก -> 403 PRINCIPAL_INCOMPLETE (fail closed ตาม plan 6.1
//     ทีมเลือก 403 เมื่อ 2569-11-09 เพื่อไม่ให้หน้าจอส่งกลับไปล็อกอินซ้ำไม่จบ)
//
// token อ่านจาก header "Authorization: Bearer <token>" (รูปแบบของ auth จำลอง ทีมเลือกเอง รอ Q-SH-03)
func Middleware(p Provider, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED")
			return
		}
		principal, err := p.Authenticate(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "UNAUTHENTICATED")
			return
		}
		if !principal.Complete() {
			writeError(w, http.StatusForbidden, "PRINCIPAL_INCOMPLETE")
			return
		}
		next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), principal)))
	})
}

func bearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get("Authorization")
	token, found := strings.CutPrefix(h, "Bearer ")
	token = strings.TrimSpace(token)
	return token, found && token != ""
}

func writeError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"code": code})
}
