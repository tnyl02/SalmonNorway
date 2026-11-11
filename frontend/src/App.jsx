import { useState } from 'react'
import EquipmentList from './pages/EquipmentList/EquipmentList.jsx'
import MyPendingRequests from './pages/MyPendingRequests/MyPendingRequests.jsx'

// หน้าที่สลับได้จากปุ่มบนหัวหน้า ยังไม่ใช้ router (ทีมตัดสิน 2569-11-11 ตอนทำ T-05 ของ 003)
// หน้าของฟีเจอร์อื่นเพิ่มเป็นแถวใหม่ในตารางนี้
const PAGES = [
  { key: 'equipment', label: 'รายการอุปกรณ์', Component: EquipmentList }, // 001 (T-06)
  { key: 'my-pending', label: 'คำขอที่รออนุมัติของฉัน', Component: MyPendingRequests }, // 003 FR-CXL-06 (T-05)
]

// หน้าจอหลักของระบบ เริ่มที่หน้ารายการอุปกรณ์ของ 001 เหมือนเดิม
export default function App() {
  const [current, setCurrent] = useState(PAGES[0].key)
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
