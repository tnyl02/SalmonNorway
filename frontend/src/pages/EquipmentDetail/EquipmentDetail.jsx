import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import { TIME_ZONE } from '../../config.js'
import StatusBadge, { OverdueBadge } from '../../components/StatusBadge.jsx'
import ErrorState from '../../components/ErrorState.jsx'
import RequestQueue from '../../components/RequestQueue.jsx'
import { RequestLoanButton } from '../../components/ItemTable.jsx'

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
// ปุ่มส่งคำขอยืมแสดงเมื่อหลังบ้านบอก canRequestLoan (สถานะ "ว่าง") กดแล้วส่ง assetCode ทาง onRequestLoan (FR-EQS-04, 05, 06, T-34)
// คำขอยืมที่ค้างอยู่ของชิ้นนี้ (FR-EQS-11, T-38) โหลดแยกจาก getItem ดู useEffect ตัวที่สอง
export default function EquipmentDetail({ assetCode, client = api, onBack, onRequestLoan }) {
  const [item, setItem] = useState(null)
  const [error, setError] = useState(null)
  const [requests, setRequests] = useState(null)

  // FR-EQS-11, CR-005: เรียก endpoint ลำดับคำขอของ 002 ทุกครั้งที่เปิดหน้า หน้าจอไม่ต้องรู้บทบาทของผู้ใช้ (plan ข้อ 5)
  // - 200: ส่งรายการให้ RequestQueue แสดง
  // - 403 FORBIDDEN_ROLE (นักศึกษา): ไม่แสดงส่วนนี้ และไม่แสดงข้อความ error (DOM-PRIV-01)
  // - error อื่น หรือ 002 ไม่ตอบ: ไม่แสดงส่วนนี้ สถานะและกำหนดคืนยังแสดงตามปกติ (CR-005 "สิ่งที่ต้องแลก")
  // เรียกผ่าน Promise.resolve().then เพื่อให้ client ที่ไม่มี listItemRequests (เช่น client จำลองของ test เดิม)
  // กลายเป็นกรณี "ไม่ตอบ" แทนที่จะทำให้หน้าพัง
  useEffect(() => {
    let active = true
    setRequests(null)
    Promise.resolve()
      .then(() => client.listItemRequests(assetCode))
      .then((res) => active && setRequests(res))
      .catch(() => active && setRequests(null))
    return () => {
      active = false
    }
  }, [client, assetCode])

  useEffect(() => {
    let active = true
    setItem(null)
    setError(null)
    client
      .getItem(assetCode)
      .then((res) => active && setItem(res))
      .catch((err) => {
        // FR-EQS-08: ไม่แสดงข้อมูลเก่า ปุ่มลองใหม่และ test ของกรณีนี้อยู่ใน T-20
        // FR-EQS-10: 403 ภาควิชาอื่น แสดงข้อความไม่มีสิทธิ์ ไม่มีข้อมูลอุปกรณ์
        if (active) {
          setItem(null)
          setError(err)
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

      {error ? (
        <ErrorState error={error} />
      ) : item == null ? (
        <p className="text-slate-600">กำลังโหลด...</p>
      ) : (
        <>
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
        {/* FR-EQS-04, 05, 06 (T-34): ปุ่มตาม canRequestLoan จากหลังบ้าน ถูกยืม เลยกำหนด หรือไม่พร้อมใช้งาน = ไม่มีปุ่ม */}
        {item.canRequestLoan === true && (
          <div>
            <RequestLoanButton assetCode={item.assetCode ?? assetCode} onRequestLoan={onRequestLoan} />
          </div>
        )}
        {/* FR-EQS-11: แสดงคู่กับรายละเอียดที่โหลดสำเร็จเท่านั้น ถ้ารายละเอียดเป็น error ไม่แสดงส่วนนี้ */}
        <RequestQueue requests={requests} formatTime={formatDueAt} />
        </>
      )}
    </section>
  )
}
