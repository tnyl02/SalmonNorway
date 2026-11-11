import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import { TIME_ZONE } from '../../config.js'
import StatusBadge, { OverdueBadge } from '../../components/StatusBadge.jsx'

// ปี พ.ศ. = ปี ค.ศ. + 543 ใช้แสดงวันที่แบบตัวอย่างใน spec ("2569-10-10 16:00" ใน AC-EQS-14)
const BUDDHIST_YEAR_OFFSET = 543

// formatDueAt แสดงกำหนดคืนเป็นวันที่และเวลาตามเขตเวลา Asia/Bangkok รูปแบบ "2569-10-10 16:00"
// (FR-EQS-04, CR-002 กำหนดคืนมีเวลา, รูปแบบตามที่ทีมตัดสิน 2569-11-10)
export function formatDueAt(iso) {
  const parts = Object.fromEntries(
    new Intl.DateTimeFormat('en-GB', {
      timeZone: TIME_ZONE,
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hourCycle: 'h23',
    })
      .formatToParts(new Date(iso))
      .map((p) => [p.type, p.value]),
  )
  return `${Number(parts.year) + BUDDHIST_YEAR_OFFSET}-${parts.month}-${parts.day} ${parts.hour}:${parts.minute}`
}

// หน้ารายละเอียดอุปกรณ์รายชิ้น (UC-01 ขั้น 7-8, FR-EQS-04)
// แสดงสถานะ กำหนดคืน (วันที่และเวลา) เมื่อถูกยืม และป้าย "เลยกำหนดคืน" เมื่อเลยกำหนด (FR-EQS-04 ข้อ ก, ข)
// นักศึกษาที่เลยกำหนดแล้วไม่ได้รับกำหนดคืนจากหลังบ้าน จึงเห็นแค่ป้าย ส่วนเจ้าหน้าที่เห็นทั้งคู่ (CR-001)
// ปุ่มส่งคำขอยืม (T-34 รอ Q-16) ยังไม่มีในตอนนี้
export default function EquipmentDetail({ assetCode, client = api, onBack }) {
  const [item, setItem] = useState(null)
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    let active = true
    setItem(null)
    setFailed(false)
    client
      .getItem(assetCode)
      .then((res) => active && setItem(res))
      .catch(() => {
        // FR-EQS-08: ไม่แสดงข้อมูลเก่า ปุ่มลองใหม่และ test ของกรณีนี้อยู่ใน T-20
        if (active) {
          setItem(null)
          setFailed(true)
        }
      })
    return () => {
      active = false
    }
  }, [client, assetCode])

  return (
    <section className="mt-4 grid gap-4">
      {onBack && (
        <button type="button" onClick={onBack} className="justify-self-start text-sm text-teal-700 underline">
          กลับไปรายการอุปกรณ์
        </button>
      )}
      <h2 className="font-mono text-lg font-semibold text-slate-900">{assetCode}</h2>

      {failed ? (
        <p role="alert" className="text-red-700">ดึงข้อมูลอุปกรณ์ไม่สำเร็จ</p>
      ) : item == null ? (
        <p className="text-slate-600">กำลังโหลด...</p>
      ) : (
        <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm">
          <dt className="text-slate-600">สถานะ</dt>
          <dd className="flex flex-wrap gap-2">
            <StatusBadge status={item.status} />
            {item.isOverdue && <OverdueBadge />}
          </dd>
          {item.dueAt && (
            <>
              <dt className="text-slate-600">กำหนดคืน</dt>
              <dd className="text-slate-900">{formatDueAt(item.dueAt)}</dd>
            </>
          )}
        </dl>
      )}
    </section>
  )
}
