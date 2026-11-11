import HandoverPage from './HandoverPage.jsx'

// โครงหน้าของ 005 ใช้หัวเรื่องและระยะขอบแบบเดียวกับ App.jsx ของ 001 (เปิดจาก src/handover-main.jsx ที่ /handover.html)
export default function HandoverApp() {
  return (
    <main className="mx-auto max-w-2xl p-4 sm:p-6">
      <h1 className="text-2xl font-bold text-teal-800">ระบบยืมอุปกรณ์ภาควิชาคอมพิวเตอร์</h1>
      <HandoverPage />
    </main>
  )
}
