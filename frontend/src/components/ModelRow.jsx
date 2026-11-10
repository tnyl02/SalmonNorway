// 1 แถวในหน้ารายการอุปกรณ์ = 1 รุ่น แสดงชื่อรุ่น ประเภท และจำนวนที่ว่าง/จำนวนทั้งหมด (FR-EQS-01, Q-13)
// รูปแบบ "ว่าง 3 / ทั้งหมด 5" ตาม AC-EQS-01 กดแล้วเลือกรุ่นนั้นเพื่อดูรายชิ้น
// เมื่อมีตัวกรอง ตัวเลขนับเฉพาะชิ้นที่ตรงตัวกรอง (FR-EQS-02, AC-EQS-15) ข้อความกำกับยังคงรูปแบบเดิมไปก่อน
// ตามที่ทีมตัดสิน 2569-11-10 จนกว่าจะได้คำตอบ Q-21 (แก้ที่บรรทัดข้อความด้านล่างจุดเดียว)
export default function ModelRow({ model, selected, onSelect }) {
  return (
    <li>
      <button
        type="button"
        aria-pressed={selected}
        onClick={() => onSelect(model.modelId)}
        className={`w-full rounded border px-4 py-3 text-left ${
          selected ? 'border-teal-600 bg-teal-50' : 'border-slate-200 bg-white'
        }`}
      >
        <span className="block font-medium text-slate-900">{model.name}</span>
        <span className="block text-sm text-slate-600">{model.category}</span>
        <span className="block text-sm text-slate-800">
          ว่าง {model.availableCount} / ทั้งหมด {model.totalCount}
        </span>
      </button>
    </li>
  )
}
