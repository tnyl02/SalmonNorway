package auth

import (
	"net/http"
	"strings"

	"salmonnorway/backend/internal/platform"
)

// Middleware บังคับให้ทุก request ยืนยันตัวตนก่อน (IF-IDP-01, NFR-SEC-01 "ต้องผ่านการยืนยันตัวตนก่อนทุกครั้ง")
//   - ไม่มี token หรือ token ใช้ไม่ได้ -> 401 UNAUTHENTICATED (plan ข้อ 5, หน้าจอใช้ต่อใน AC-EQS-09)
//   - body ของ error เป็น {code, message} ตาม CON-TECH-03 (platform.WriteError, T-37)
//   - ยืนยันตัวตนได้ แต่ไม่มีภาควิชาหรือบทบาทไม่รู้จัก -> 403 PRINCIPAL_INCOMPLETE (fail closed ตาม plan 6.1
//     ทีมเลือก 403 เมื่อ 2569-11-09 เพื่อไม่ให้หน้าจอส่งกลับไปล็อกอินซ้ำไม่จบ)
//
// token อ่านจาก header "Authorization: Bearer <token>" (รูปแบบของ auth จำลอง ทีมเลือกเอง รอ Q-SH-03)
func Middleware(p Provider, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r)
		if !ok {
			platform.WriteError(w, http.StatusUnauthorized, platform.CodeUnauthenticated)
			return
		}
		principal, err := p.Authenticate(r.Context(), token)
		if err != nil {
			platform.WriteError(w, http.StatusUnauthorized, platform.CodeUnauthenticated)
			return
		}
		if !principal.Complete() {
			platform.WriteError(w, http.StatusForbidden, platform.CodePrincipalIncomplete)
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
