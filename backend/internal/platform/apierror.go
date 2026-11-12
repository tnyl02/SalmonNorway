package platform

import (
	"log"
	"net/http"
)

// รหัส error จากรายการกลางของ CON-TECH-03 (specs/000-shared/constraints.md, CR-004) ที่ของกลางและ 001 ใช้
// ห้ามตั้งรหัสใหม่ในโค้ดโดยไม่ลงทะเบียนในรายการกลางก่อน
const (
	CodeUnauthenticated     = "UNAUTHENTICATED"      // 401 ไม่มี session หรือหมดอายุ (IF-IDP-01)
	CodePrincipalIncomplete = "PRINCIPAL_INCOMPLETE" // 403 ยืนยันตัวตนได้ แต่ไม่มีภาควิชาหรือบทบาทไม่รู้จัก (IF-IDP-01)
	CodeForbiddenDepartment = "FORBIDDEN_DEPARTMENT" // 403 อุปกรณ์ของภาควิชาอื่น (FR-EQS-10, DOM-ACCESS-01) ใช้ใน T-18
	CodeForbiddenRole       = "FORBIDDEN_ROLE"       // 403 บทบาทนี้ไม่มีสิทธิ์ เช่น นักศึกษาเรียกสรุปของเจ้าหน้าที่ (FR-EQS-07, T-17)
	CodeUpstreamUnavailable = "UPSTREAM_UNAVAILABLE" // 503 ฐานข้อมูลไม่ตอบ (FR-EQS-08)
	CodeItemNotFound        = "ITEM_NOT_FOUND"       // 404 ไม่มีรหัสครุภัณฑ์นี้ หรือรหัสรุ่นไม่ถูกต้อง (FR-EQS-01, FR-EQS-04, 001 ลงทะเบียน 2569-11-12, T-39)
	CodeMethodNotAllowed    = "METHOD_NOT_ALLOWED"   // 405 method ที่ไม่รองรับ ทุกฟีเจอร์ผ่าน MethodNotAllowedJSON (CON-DATA-02, 001 ลงทะเบียน 2569-11-12, T-40)
)

// errorMessages คือข้อความภาษาไทยของแต่ละรหัส (CON-TECH-03 ต้องมี message)
// ถ้อยคำทีมเลือกเอง ไม่ได้มาจาก spec เสนอใน plan.md ข้อ 5 ของ 001 แก้คำได้ที่นี่จุดเดียว
// หน้าจอตัดสินใจจาก code เท่านั้น ห้ามเอา message ไปเทียบ
var errorMessages = map[string]string{
	CodeUnauthenticated:     "กรุณาเข้าสู่ระบบก่อนใช้งาน",
	CodePrincipalIncomplete: "บัญชีนี้ไม่มีข้อมูลภาควิชาหรือบทบาท กรุณาติดต่อเจ้าหน้าที่",
	CodeForbiddenDepartment: "ไม่มีสิทธิ์เข้าถึงอุปกรณ์ของภาควิชาอื่น",
	CodeForbiddenRole:       "เฉพาะเจ้าหน้าที่เท่านั้นที่ดูข้อมูลนี้ได้",
	CodeUpstreamUnavailable: "ดึงข้อมูลไม่สำเร็จ กรุณาลองใหม่อีกครั้ง",
	CodeItemNotFound:        "ไม่พบอุปกรณ์รหัสนี้",
	CodeMethodNotAllowed:    "ไม่รองรับการเรียกแบบนี้",
}

// ErrorBody คือ body ของ error ตามรูปแบบกลาง {"code", "message"} ไม่มีช่องอื่น (CON-TECH-03)
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorMessage คืนข้อความของรหัส ถ้ารหัสนี้ยังไม่มีข้อความ คืนรหัสเอง เพื่อให้ message ไม่ว่าง (CON-TECH-03)
func ErrorMessage(code string) string {
	if m, ok := errorMessages[code]; ok {
		return m
	}
	log.Printf("apierror: no message for code %q (เพิ่มข้อความใน platform/apierror.go)", code)
	return code
}

// WriteError ตอบ error ด้วย body {code, message} ตาม CON-TECH-03 ใช้ร่วมทุกฟีเจอร์ที่ใช้รหัสข้างบน
func WriteError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, ErrorBody{Code: code, Message: ErrorMessage(code)})
}
