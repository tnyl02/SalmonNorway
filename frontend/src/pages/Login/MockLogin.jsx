import { setToken } from '../../api/client.js'

// ผู้ใช้จำลองตรงกับ backend/internal/auth/mock_provider.go (ทีมเลือกเอง ไม่ได้มาจาก spec รอ Q-SH-03)
// ระบบจริงจะส่งไปหน้าเข้าสู่ระบบของมหาวิทยาลัย (IF-IDP-01) ไม่ใช่หน้านี้ การยืนยันตัวตนไม่ใช่ use case ของระบบนี้
export const MOCK_USERS = [
  { token: 'mock-student-001', label: 'นักศึกษา student-001 (ภาควิชา 1)' },
  { token: 'mock-student-002', label: 'นักศึกษา student-002 (ภาควิชา 1)' },
  { token: 'mock-student-101', label: 'นักศึกษา student-101 (ภาควิชา 2)' },
  { token: 'mock-staff-001', label: 'เจ้าหน้าที่ staff-001 (ภาควิชา 1)' },
]

// หน้าเข้าสู่ระบบจำลอง (IF-IDP-01, AC-EQS-09) แสดงเมื่อ session หมดอายุ พร้อมข้อความแจ้ง
// เลือกผู้ใช้แล้วตั้ง token ใหม่ จากนั้นให้ App พากลับหน้าเดิม (onLoggedIn)
export default function MockLogin({ message, onLoggedIn }) {
  const login = (token) => {
    setToken(token)
    onLoggedIn()
  }
  return (
    <section className="mt-4 grid gap-4">
      <h2 className="text-lg font-semibold text-slate-900">เข้าสู่ระบบ (จำลอง)</h2>
      {message && (
        <p role="alert" className="text-red-700">
          {message}
        </p>
      )}
      <ul className="grid gap-2">
        {MOCK_USERS.map((u) => (
          <li key={u.token}>
            <button
              type="button"
              onClick={() => login(u.token)}
              className="w-full rounded border border-slate-300 bg-white px-4 py-2 text-left"
            >
              {u.label}
            </button>
          </li>
        ))}
      </ul>
    </section>
  )
}
