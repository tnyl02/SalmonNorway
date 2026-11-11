import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import { APPROVAL_MESSAGES, PURPOSE_OPTIONS, PURPOSE_OTHER, TIME_ZONE } from '../../config.js'

// แสดงวันเวลาตามเขตเวลา Asia/Bangkok (glossary)
const dateTime = new Intl.DateTimeFormat('th-TH', { timeZone: TIME_ZONE, dateStyle: 'medium', timeStyle: 'short' })
const fmt = (iso) => dateTime.format(new Date(iso))

// วัตถุประสงค์: ข้อความของ 3 ค่าจาก config ถ้าเป็น "อื่น ๆ" แสดงข้อความที่ผู้ขอพิมพ์ (FR-BRQ-02, FR-APV-01 ก)
function purposeLabel(r) {
  if (r.purpose === PURPOSE_OTHER && r.purposeOther) return r.purposeOther
  return PURPOSE_OPTIONS.find((o) => o.value === r.purpose)?.label ?? r.purpose
}

// หน้ารายการคำขอยืมแบบปกติที่รอพิจารณาของภาควิชาตนเอง (FR-APV-01 ก, AC-APV-11)
// เรียงตามเวลาที่ส่งจากเก่าไปใหม่ตามที่หลังบ้านส่งมา ไม่เรียงใหม่ที่หน้าจอ (DOM-REQ-01)
// ปุ่มอนุมัติ / ปฏิเสธเพิ่มใน T-09, T-11 / หน้าจอไม่มี mockup ทีมออกแบบเอง (plan.md ข้อ 4)
// client ส่งเข้ามาได้ ตอน test ใช้ client จำลองแทนหลังบ้านจริง
export default function ApproveRequests({ client = api }) {
  const [list, setList] = useState(null) // { items, loadedAt } จาก GET /v1/approvals/requests
  const [errorCode, setErrorCode] = useState(null)

  useEffect(() => {
    let active = true
    client
      .getApprovalRequests()
      .then((res) => active && setList(res))
      .catch((err) => active && setErrorCode(err?.code ?? 'UPSTREAM_UNAVAILABLE'))
    return () => {
      active = false
    }
  }, [client])

  if (errorCode) {
    // ข้อความเลือกตามรหัส ไม่ใช่ตามข้อความจากหลังบ้าน (CON-TECH-03)
    const message = APPROVAL_MESSAGES[errorCode] ?? APPROVAL_MESSAGES.UPSTREAM_UNAVAILABLE
    return (
      <p role="alert" className="mt-4 text-red-700">
        {message}
      </p>
    )
  }
  if (list == null) return <p className="mt-4 text-slate-600">กำลังโหลด...</p>

  return (
    <section className="mt-4 grid gap-2">
      <h2 className="text-lg font-semibold text-slate-900">คำขอยืมที่รอพิจารณา</h2>
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="border-b">
            <th className="py-1">ผู้ขอ</th>
            <th className="py-1">อุปกรณ์</th>
            <th className="py-1">วัตถุประสงค์</th>
            <th className="py-1">ช่วงเวลาที่ต้องการยืม</th>
            <th className="py-1">เวลาที่ส่ง</th>
          </tr>
        </thead>
        <tbody>
          {list.items.map((r) => (
            <tr key={r.requestId} className="border-b">
              <td className="py-1">{r.requesterName}</td>
              <td className="py-1">
                {r.modelName} <span className="text-slate-500">({r.assetCode})</span>
              </td>
              <td className="py-1">{purposeLabel(r)}</td>
              <td className="py-1">
                {fmt(r.startAt)} - {fmt(r.endAt)}
              </td>
              <td className="py-1">{fmt(r.submittedAt)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}
