import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import {
  LABEL_RETRY,
  MSG_ALREADY_APPROVED,
  MSG_CANCEL_SUCCESS,
  MSG_LOAD_FAILED,
  MSG_NO_PENDING,
  MSG_NOT_CANCELLABLE,
} from '../../config.js'
import { formatDueAt } from '../EquipmentDetail/EquipmentDetail.jsx'
import CancelDialog from './CancelDialog.jsx'

// formatSubmittedAt แสดงเวลาที่ส่งคำขอตามเขตเวลา Asia/Bangkok (FR-CXL-06)
// ใช้รูปแบบเดียวกับ 001 "2569-11-10 10:30" (ทีมตัดสิน 2569-11-11 ตอนทำ T-05)
const formatSubmittedAt = formatDueAt

// หน้า "คำขอที่รออนุมัติของฉัน" (FR-CXL-06, UC-15)
// แสดงเฉพาะคำขอของนักศึกษาที่เข้าสู่ระบบที่ยังยกเลิกได้ หลังบ้านกรองและเรียงใหม่ไปเก่ามาแล้ว (FR-CXL-03, CON-CXL-02)
// กดยกเลิกแล้วต้องยืนยันพร้อมเหตุผลก่อน ยกเลิกสำเร็จแล้วแถวหายจากหน้า (FR-CXL-07, FR-CXL-01)
// ไม่มีคำขอแสดง MSG_NO_PENDING โหลดไม่สำเร็จแสดง MSG_LOAD_FAILED พร้อมปุ่มลองใหม่ ไม่แสดงรายการเก่า (FR-CXL-06, ASM-13)
// ยกเลิกไม่ได้ตามตาราง error ใน plan.md ข้อ 4: แสดงข้อความแล้วโหลดรายการใหม่ (FR-CXL-02) ดู CANCEL_REJECTED
// ข้อความของกรณีอื่นอยู่ใน T-20, T-24

// CANCEL_REJECTED คือรหัส error จากหลังบ้านที่แปลว่ายกเลิกไม่ได้ จับคู่กับข้อความใน config.js (plan.md ข้อ 4)
// ได้รหัสเหล่านี้แล้ว ปิดกล่อง แสดงข้อความ และโหลดรายการใหม่ให้ตรงกับฐานข้อมูล
const CANCEL_REJECTED = {
  ALREADY_APPROVED: MSG_ALREADY_APPROVED, // 409 FR-CXL-02, FR-CXL-04 (T-16)
  NOT_CANCELLABLE: MSG_NOT_CANCELLABLE, // 409 FR-CXL-02 rejected, handed_over, returned, cancelled (T-18)
}

export default function MyPendingRequests({ client = api }) {
  const [items, setItems] = useState(null)
  const [loadFailed, setLoadFailed] = useState(false)
  const [reloadKey, setReloadKey] = useState(0) // เพิ่มค่าเมื่อกดลองใหม่ ให้โหลดรายการอีกครั้ง
  const [target, setTarget] = useState(null) // คำขอที่เปิดกล่องยืนยันอยู่
  const [busy, setBusy] = useState(false)
  const [success, setSuccess] = useState(null) // ข้อความยกเลิกสำเร็จ ค้างจนเริ่มยกเลิกรายการถัดไป (ทีมตัดสิน 2569-11-11 T-07)
  const [cancelError, setCancelError] = useState(null) // ข้อความเมื่อยกเลิกไม่ได้ ค้างจนเริ่มยกเลิกรายการถัดไป แบบเดียวกับ success

  useEffect(() => {
    let active = true
    setItems(null)
    setLoadFailed(false)
    client
      .listMyPendingRequests()
      .then((res) => active && setItems(res.items))
      .catch(() => {
        // โหลดไม่สำเร็จ: ไม่แสดงรายการเก่าและไม่แสดงว่าไม่มีคำขอ (FR-CXL-06)
        if (active) setLoadFailed(true)
      })
    return () => {
      active = false
    }
  }, [client, reloadKey])

  // เปิดกล่องยืนยัน (FR-CXL-07)
  function openDialog(request) {
    setSuccess(null)
    setCancelError(null)
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
    } catch (err) {
      const message = CANCEL_REJECTED[err?.code]
      if (message) {
        // FR-CXL-02, FR-CXL-04: สถานะเปลี่ยนไปแล้ว ปิดกล่อง แจ้งเหตุผล แล้วโหลดรายการใหม่
        setTarget(null)
        setCancelError(message)
        setReloadKey((k) => k + 1)
      }
      // กรณีอื่น: คงแถวและกล่องไว้ ข้อความอยู่ใน T-20, T-24
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
      {cancelError && (
        <p role="alert" className="rounded border border-red-300 bg-red-50 p-3 text-sm text-red-800">
          {cancelError}
        </p>
      )}
      {loadFailed ? (
        <div className="grid justify-items-start gap-2">
          <p role="alert" className="text-red-700">
            {MSG_LOAD_FAILED}
          </p>
          <button
            type="button"
            onClick={() => setReloadKey((k) => k + 1)}
            className="rounded border border-teal-700 px-3 py-1 text-sm text-teal-700"
          >
            {LABEL_RETRY}
          </button>
        </div>
      ) : items == null ? (
        <p className="text-slate-600">กำลังโหลด...</p>
      ) : items.length === 0 ? (
        <p className="text-slate-600">{MSG_NO_PENDING}</p>
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
      {target && (
        <CancelDialog
          key={target.requestId}
          request={target}
          onConfirm={confirmCancel}
          onCancel={() => setTarget(null)} // ไม่ยืนยัน: ปิดกล่อง แถวยังอยู่ (FR-CXL-07, T-08)
          busy={busy}
        />
      )}
    </section>
  )
}
