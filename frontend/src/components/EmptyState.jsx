// แจ้งเมื่อค้นหาหรือกรองแล้วไม่พบอุปกรณ์ พร้อมทางเลือกล้างตัวกรองกลับไปดูรายการทั้งหมด (FR-EQS-03, UC-01 alternative 6a)
// ข้อความ "ไม่พบอุปกรณ์" และ "ล้างตัวกรอง" ใช้คำตาม FR-EQS-03 และ AC-EQS-03
export default function EmptyState({ onClear }) {
  return (
    <div role="status" className="grid justify-items-start gap-2 rounded border border-slate-200 bg-white px-4 py-6">
      <p className="text-slate-800">ไม่พบอุปกรณ์</p>
      <button
        type="button"
        onClick={onClear}
        className="rounded border border-teal-600 px-3 py-1 text-sm text-teal-700"
      >
        ล้างตัวกรอง
      </button>
    </div>
  )
}
