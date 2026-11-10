// ป้ายสถานะของอุปกรณ์รายชิ้น มี 3 สถานะเท่านั้น (CON-DATA-01, CR-002)
// ค่าที่ได้จากหลังบ้านเป็นภาษาอังกฤษ (available / borrowed / unavailable) หน้าจอแสดงเป็นคำใน glossary
export const STATUS_LABEL = {
  available: 'ว่าง',
  borrowed: 'ถูกยืม',
  unavailable: 'ไม่พร้อมใช้งาน',
}

// ป้าย "เลยกำหนดคืน" แสดงคู่กับสถานะ "ถูกยืม" ไม่ใช่สถานะที่ 4 (FR-EQS-04 ข้อ ข, CON-DATA-01, CR-001)
// หลังบ้านเป็นคนตัดสินว่าเลยกำหนดหรือยัง (isOverdue) หน้าจอไม่คำนวณเอง
export function OverdueBadge() {
  return (
    <span className="inline-block rounded border border-red-300 bg-red-50 px-2 py-0.5 text-sm text-red-700">
      เลยกำหนดคืน
    </span>
  )
}

export default function StatusBadge({ status }) {
  return (
    <span className="inline-block rounded border border-slate-300 px-2 py-0.5 text-sm text-slate-700">
      {STATUS_LABEL[status] ?? status}
    </span>
  )
}
