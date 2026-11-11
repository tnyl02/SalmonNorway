import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import { MSG_CANCEL_SUCCESS } from '../../config.js'
import { formatDueAt } from '../EquipmentDetail/EquipmentDetail.jsx'
import CancelDialog from './CancelDialog.jsx'

// formatSubmittedAt แสดงเวลาที่ส่งคำขอตามเขตเวลา Asia/Bangkok (FR-CXL-06)
// ใช้รูปแบบเดียวกับ 001 "2569-11-10 10:30" (ทีมตัดสิน 2569-11-11 ตอนทำ T-05)
const formatSubmittedAt = formatDueAt

// หน้า "คำขอที่รออนุมัติของฉัน" (FR-CXL-06, UC-15)
// แสดงเฉพาะคำขอของนักศึกษาที่เข้าสู่ระบบที่ยังยกเลิกได้ หลังบ้านกรองและเรียงใหม่ไปเก่ามาแล้ว (FR-CXL-03, CON-CXL-02)
// กดยกเลิกแล้วต้องยืนยันพร้อมเหตุผลก่อน ยกเลิกสำเร็จแล้วแถวหายจากหน้า (FR-CXL-07, FR-CXL-01)
// ข้อความไม่มีคำขอและโหลดไม่สำเร็จอยู่ใน T-12 ข้อความเมื่อยกเลิกไม่ได้อยู่ใน T-16, T-18, T-20, T-24
export default function MyPendingRequests({ client = api }) {
  const [items, setItems] = useState(null)
  const [target, setTarget] = useState(null) // คำขอที่เปิดกล่องยืนยันอยู่
  const [busy, setBusy] = useState(false)
  const [success, setSuccess] = useState(null) // ข้อความยกเลิกสำเร็จ ค้างจนเริ่มยกเลิกรายการถัดไป (ทีมตัดสิน 2569-11-11 T-07)

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

  // เปิดกล่องยืนยัน (FR-CXL-07)
  function openDialog(request) {
    setSuccess(null)
    setTarget(request)
  }

  // ยืนยันการยกเลิก (FR-CXL-01, FR-CXL-07) หลังบ้านเป็นผู้ตัดสินว่ายกเลิกได้ไหม (FR-CXL-04)
  async function confirmCancel(reason) {
    const request = target
    setBusy(true)
    try {
      await client.cancelRequest(request.requestId, reason)
      setItems((current) => current.filter((r) => r.requestId !== request.requestId))
      setSuccess(MSG_CANCEL_SUCCESS(request.modelName, request.assetCode))
      setTarget(null)
    } catch {
      // ยกเลิกไม่สำเร็จ: คงแถวและกล่องไว้ ข้อความของแต่ละกรณีอยู่ใน T-16, T-18, T-20, T-24
    } finally {
      setBusy(false)
    }
  }

  return (
    <section className="mt-4 grid gap-4">
      <h2 className="text-lg font-semibold text-slate-900">คำขอที่รออนุมัติของฉัน</h2>
      {success && (
        <p role="status" className="rounded border border-green-300 bg-green-50 p-3 text-sm text-green-800">
          {success}
        </p>
      )}
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
              <button
                type="button"
                onClick={() => openDialog(r)}
                className="rounded border border-red-300 px-3 py-1 text-sm text-red-700"
              >
                ยกเลิก
              </button>
            </li>
          ))}
        </ul>
      )}
      {target && <CancelDialog key={target.requestId} request={target} onConfirm={confirmCancel} busy={busy} />}
    </section>
  )
}
