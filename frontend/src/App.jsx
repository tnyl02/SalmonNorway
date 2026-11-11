// หน้าหลักของหน้าจอ ตอนนี้แสดงหน้าของ 006 หน้าของฟีเจอร์อื่นจะถูกเพิ่มตาม task ของฟีเจอร์นั้น
import { useState } from 'react'
import OpenLoansPage from './pages/Return/OpenLoansPage.jsx'
import ReturnPage from './pages/Return/ReturnPage.jsx'

export default function App() {
  // รหัสครุภัณฑ์ที่เลือกจากรายการที่ยังไม่คืน: เปลี่ยน key เพื่อเปิด ReturnPage ใหม่พร้อมรหัสนั้น (006 FR-RET-01(ข))
  const [selected, setSelected] = useState({ assetCode: '', key: 0 })

  return (
    <main className="mx-auto max-w-2xl p-6">
      <h1 className="text-2xl font-bold text-teal-800">ระบบยืมอุปกรณ์ภาควิชาคอมพิวเตอร์</h1>
      {/* 006-return บันทึกการรับคืน (เจ้าหน้าที่) */}
      <ReturnPage key={selected.key} initialAssetCode={selected.assetCode} />
      <OpenLoansPage onSelect={(assetCode) => setSelected((s) => ({ assetCode, key: s.key + 1 }))} />
    </main>
  )
}
