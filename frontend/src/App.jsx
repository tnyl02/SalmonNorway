import EquipmentList from './pages/EquipmentList/EquipmentList.jsx'

// หน้าจอหลักของระบบ ตอนนี้มีหน้ารายการอุปกรณ์ของ 001 (T-06)
// หน้าจอของฟีเจอร์อื่นจะถูกสร้างใน src/pages/ ตาม task ของฟีเจอร์นั้น
export default function App() {
  return (
    <main className="mx-auto max-w-2xl p-4 sm:p-6">
      <h1 className="text-2xl font-bold text-teal-800">ระบบยืมอุปกรณ์ภาควิชาคอมพิวเตอร์</h1>
      <EquipmentList />
    </main>
  )
}
