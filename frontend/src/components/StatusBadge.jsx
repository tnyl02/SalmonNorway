// ป้ายสถานะของอุปกรณ์รายชิ้น มี 3 สถานะเท่านั้น (CON-DATA-01, CR-002)
// ค่าที่ได้จากหลังบ้านเป็นภาษาอังกฤษ (available / borrowed / unavailable) หน้าจอแสดงเป็นคำใน glossary
// ป้าย "เลยกำหนดคืน" (FR-EQS-04) จะเพิ่มใน T-13
export const STATUS_LABEL = {
  available: 'ว่าง',
  borrowed: 'ถูกยืม',
  unavailable: 'ไม่พร้อมใช้งาน',
}

export default function StatusBadge({ status }) {
  return (
    <span className="inline-block rounded border border-slate-300 px-2 py-0.5 text-sm text-slate-700">
      {STATUS_LABEL[status] ?? status}
    </span>
  )
}
