import { BORROW_REQUEST_MESSAGES, FIELD_LABELS, MSG_FIELDS_TO_FIX, MSG_SUBMIT_FAILED } from '../config.js'
import { STATUS_LABEL } from './StatusBadge.jsx'

// สรุปข้อผิดพลาดของการส่งคำขอยืมให้ผู้ขอเห็นว่าเกิดอะไรและต้องแก้ตรงไหน (FR-BRQ-03, FR-BRQ-09, CON-TECH-03)
// เลือกข้อความจาก "รหัส" ที่หลังบ้านส่งมา ไม่ใช่จากข้อความ (CON-TECH-03) รหัสที่ยังไม่รู้จักใช้ข้อความจากหลังบ้านแทน
// รหัส ITEM_NOT_AVAILABLE แสดงสถานะล่าสุดของอุปกรณ์ต่อท้าย ตามที่ FR-BRQ-01 สั่ง โดยใช้คำในภาษาไทยชุดเดียวกับ 001
// รหัส VALIDATION_FAILED แสดงรายชื่อช่องที่ต้องแก้จาก fields ที่หลังบ้านส่งมา (FR-BRQ-03, AC-BRQ-03, AC-BRQ-10)
export default function RequestErrorSummary({ error }) {
  if (!error) return null

  return (
    <div role="alert" className="grid justify-items-start gap-1 text-red-700">
      <p>{summaryText(error)}</p>
      {error.fields?.length > 0 && (
        <p className="text-sm">
          {MSG_FIELDS_TO_FIX} {error.fields.map((field) => FIELD_LABELS[field] ?? field).join(', ')}
        </p>
      )}
    </div>
  )
}

// summaryText เลือกข้อความหลักจากรหัส แล้วเติมสถานะล่าสุดเมื่ออุปกรณ์ไม่ว่าง (FR-BRQ-01)
export function summaryText({ code, message, status }) {
  const text = BORROW_REQUEST_MESSAGES[code] ?? message ?? MSG_SUBMIT_FAILED
  if (code === 'ITEM_NOT_AVAILABLE' && status) {
    return `${text} (สถานะล่าสุด: ${STATUS_LABEL[status] ?? status})`
  }
  return text
}
