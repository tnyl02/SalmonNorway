import { STATUS_LABEL } from './StatusBadge.jsx'

// ตัวกรองของหน้ารายการอุปกรณ์: ประเภท และสถานะ (FR-EQS-02, plan.md ข้อ 6.7 ตัวกรองเป็น chip)
// เลือกได้หลายค่าในแต่ละตัวกรอง กดซ้ำเพื่อเอาค่านั้นออก (FR-EQS-02, Q-15, T-33)
// ยังไม่มีช่องค้นหาด้วยคำ เพราะวิธีจับคู่คำค้นรอ Q-06 (T-24) ตามที่ทีมตัดสิน 2569-11-10
//
// filter = { categories: [...], statuses: [...] }  ค่าว่าง = ไม่กรอง
export default function SearchBar({ categories, filter, onChange }) {
  const pick = (key, value) => {
    const current = filter[key] ?? []
    const next = current.includes(value) ? current.filter((v) => v !== value) : [...current, value]
    onChange({ ...filter, [key]: next })
  }

  return (
    <div className="grid gap-3">
      <ChipGroup
        label="ประเภท"
        options={categories.map((c) => ({ value: c, label: c }))}
        selected={filter.categories ?? []}
        onPick={(v) => pick('categories', v)}
      />
      <ChipGroup
        label="สถานะ"
        options={Object.entries(STATUS_LABEL).map(([value, label]) => ({ value, label }))}
        selected={filter.statuses ?? []}
        onPick={(v) => pick('statuses', v)}
      />
    </div>
  )
}

function ChipGroup({ label, options, selected, onPick }) {
  return (
    <div role="group" aria-label={label} className="flex flex-wrap items-center gap-2">
      <span className="text-sm text-slate-600">{label}</span>
      {options.map((o) => {
        const on = selected.includes(o.value)
        return (
          <button
            key={o.value}
            type="button"
            aria-pressed={on}
            onClick={() => onPick(o.value)}
            className={`rounded-full border px-3 py-1 text-sm ${
              on ? 'border-teal-600 bg-teal-600 text-white' : 'border-slate-300 bg-white text-slate-700'
            }`}
          >
            {o.label}
          </button>
        )
      })}
    </div>
  )
}
