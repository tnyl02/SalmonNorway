import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import { formatDueAt } from '../EquipmentDetail/EquipmentDetail.jsx'

// formatSubmittedAt แสดงเวลาที่ส่งคำขอตามเขตเวลา Asia/Bangkok (FR-CXL-06)
// ใช้รูปแบบเดียวกับ 001 "2569-11-10 10:30" (ทีมตัดสิน 2569-11-11 ตอนทำ T-05)
const formatSubmittedAt = formatDueAt

// หน้า "คำขอที่รออนุมัติของฉัน" (FR-CXL-06, UC-15)
// แสดงเฉพาะคำขอของนักศึกษาที่เข้าสู่ระบบที่ยังยกเลิกได้ หลังบ้านกรองและเรียงใหม่ไปเก่ามาแล้ว (FR-CXL-03, CON-CXL-02)
// T-05 ทำเฉพาะการแสดงรายการ กล่องยืนยันอยู่ใน T-07 ข้อความไม่มีคำขอและโหลดไม่สำเร็จอยู่ใน T-12
export default function MyPendingRequests({ client = api }) {
  const [items, setItems] = useState(null)

  useEffect(() => {
    let active = true
    client
      .listMyPendingRequests()
      .then((res) => active && setItems(res.items))
      .catch(() => {
        // ข้อความและปุ่มลองใหม่ของกรณีนี้อยู่ใน T-12 (FR-CXL-06) ตอนนี้ไม่แสดงรายการเก่า
        if (active) setItems(null)
      })
    return () => {
      active = false
    }
  }, [client])

  return (
    <section className="mt-4 grid gap-4">
      <h2 className="text-lg font-semibold text-slate-900">คำขอที่รออนุมัติของฉัน</h2>
      {items == null ? (
        <p className="text-slate-600">กำลังโหลด...</p>
      ) : (
        <ul className="grid gap-2">
          {items.map((r) => (
            <li
              key={r.requestId}
              className="flex flex-wrap items-center justify-between gap-2 rounded border border-slate-200 p-3"
            >
              <div className="grid gap-0.5">
                <span className="font-medium text-slate-900">{r.modelName}</span>
                <span className="font-mono text-sm text-slate-700">{r.assetCode}</span>
                <span className="text-sm text-slate-600">ส่งเมื่อ {formatSubmittedAt(r.submittedAt)}</span>
              </div>
              {/* ปุ่มยกเลิก (FR-CXL-06) กล่องยืนยันและการยกเลิกจริงอยู่ใน T-07 */}
              <button type="button" className="rounded border border-red-300 px-3 py-1 text-sm text-red-700">
                ยกเลิก
              </button>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
