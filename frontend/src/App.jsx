// หน้าหลักของหน้าจอ ตอนนี้แสดงหน้าของ 006 (ReturnPage) หน้าของฟีเจอร์อื่นจะถูกเพิ่มตาม task ของฟีเจอร์นั้น
import ReturnPage from './pages/Return/ReturnPage.jsx'

export default function App() {
  return (
    <main className="mx-auto max-w-2xl p-6">
      <h1 className="text-2xl font-bold text-teal-800">ระบบยืมอุปกรณ์ภาควิชาคอมพิวเตอร์</h1>
      {/* 006-return บันทึกการรับคืน (เจ้าหน้าที่) */}
      <ReturnPage />
    </main>
  )
}
