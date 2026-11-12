import { useEffect, useRef, useState } from 'react'
import { onUnauthorized } from './api/client.js'
import MockLogin from './pages/Login/MockLogin.jsx'
import EquipmentList from './pages/EquipmentList/EquipmentList.jsx'
import MyPendingRequests from './pages/MyPendingRequests/MyPendingRequests.jsx'
import BorrowRequestForm from './pages/BorrowRequest/BorrowRequestForm.jsx'
import ReturnsHome from './pages/Return/ReturnsHome.jsx'

// หน้าที่สลับได้จากปุ่มบนหัวหน้า ยังไม่ใช้ router (ทีมตัดสิน 2569-11-11 ตอนทำ T-05 ของ 003)
// หน้าของฟีเจอร์อื่นเพิ่มเป็นแถวใหม่ในตารางนี้
const PAGES = [
  { key: 'equipment', label: 'รายการอุปกรณ์', Component: EquipmentList }, // 001 (T-06)
  { key: 'my-pending', label: 'คำขอที่รออนุมัติของฉัน', Component: MyPendingRequests }, // 003 FR-CXL-06 (T-05)
  { key: 'borrow-request', label: 'ส่งคำขอยืม', Component: BorrowRequestForm }, // 002 FR-BRQ-03 (T-19)
  { key: 'returns', label: 'รับคืนอุปกรณ์', Component: ReturnsHome }, // 006 FR-RET-01 เจ้าหน้าที่ (T-06, T-15)
]

// safePageKey รับเฉพาะชื่อหน้าที่อยู่ในตาราง PAGES ค่าอื่น (เช่น URL ภายนอก) กลับไปหน้าแรก กัน open redirect (plan 001 ข้อ 6.5)
export function safePageKey(key) {
  return PAGES.some((p) => p.key === key) ? key : PAGES[0].key
}

// ข้อความเมื่อหลังบ้านไม่ได้ส่ง message มา (ปกติใช้ message ของ 401 ตาม CON-TECH-03)
const SESSION_EXPIRED_MESSAGE = 'กรุณาเข้าสู่ระบบก่อนใช้งาน'

// หน้าจอหลักของระบบ เริ่มที่หน้ารายการอุปกรณ์ของ 001 เหมือนเดิม
// session หมดอายุ (401) แสดงหน้าล็อกอินจำลองพร้อมข้อความ แล้วกลับหน้าเดิมหลังล็อกอิน (IF-IDP-01, AC-EQS-09, T-21 ของ 001)
// ยังไม่ใช้ router จึงจำหน้าเดิมไว้ในหน่วยความจำ (ทีมตัดสิน 2569-11-12)
export default function App() {
  const [current, setCurrent] = useState(PAGES[0].key)
  const [login, setLogin] = useState(null) // { message, returnTo } เมื่อ session หมดอายุ
  const currentRef = useRef(current)
  currentRef.current = current

  useEffect(
    () =>
      onUnauthorized((err) =>
        setLogin((prev) => prev ?? { message: err.serverMessage || SESSION_EXPIRED_MESSAGE, returnTo: currentRef.current }),
      ),
    [],
  )

  if (login) {
    return (
      <main className="mx-auto max-w-2xl p-4 sm:p-6">
        <h1 className="text-2xl font-bold text-teal-800">ระบบยืมอุปกรณ์ภาควิชาคอมพิวเตอร์</h1>
        <MockLogin
          message={login.message}
          onLoggedIn={() => {
            setCurrent(safePageKey(login.returnTo))
            setLogin(null)
          }}
        />
      </main>
    )
  }

  const { Component } = PAGES.find((p) => p.key === current)

  return (
    <main className="mx-auto max-w-2xl p-4 sm:p-6">
      <h1 className="text-2xl font-bold text-teal-800">ระบบยืมอุปกรณ์ภาควิชาคอมพิวเตอร์</h1>
      <nav className="mt-3 flex flex-wrap gap-2">
        {PAGES.map((p) => (
          <button
            key={p.key}
            type="button"
            aria-current={p.key === current ? 'page' : undefined}
            onClick={() => setCurrent(p.key)}
            className={
              p.key === current
                ? 'rounded bg-teal-700 px-3 py-1 text-sm text-white'
                : 'rounded border border-teal-700 px-3 py-1 text-sm text-teal-700'
            }
          >
            {p.label}
          </button>
        ))}
      </nav>
      <Component />
    </main>
  )
}
