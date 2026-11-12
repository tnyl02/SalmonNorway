import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import {
  LABEL_RETRY,
  MSG_ALREADY_APPROVED,
  MSG_CANCEL_FAILED,
  MSG_CANCEL_SUCCESS,
  MSG_LOAD_FAILED,
  MSG_NO_PENDING,
  MSG_NOT_CANCELLABLE,
  MSG_REQUEST_NOT_FOUND,
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
// ระบบไม่ตอบสนอง (503) หรือส่งไม่ถึง: กล่องยังเปิด แสดง MSG_CANCEL_FAILED ในกล่อง ลองใหม่ด้วยเหตุผลเดิม แถวยังอยู่ (FR-CXL-01, ASM-06)

// CANCEL_REJECTED คือรหัส error จากหลังบ้านที่แปลว่ายกเลิกไม่ได้ จับคู่กับข้อความใน config.js (plan.md ข้อ 4)
// ได้รหัสเหล่านี้แล้ว ปิดกล่อง แสดงข้อความ และโหลดรายการใหม่ให้ตรงกับฐานข้อมูล
const CANCEL_REJECTED = {
  ALREADY_APPROVED: MSG_ALREADY_APPROVED, // 409 FR-CXL-02, FR-CXL-04 (T-16)
  NOT_CANCELLABLE: MSG_NOT_CANCELLABLE, // 409 FR-CXL-02 rejected, handed_over, returned, cancelled (T-18)
  REQUEST_NOT_FOUND: MSG_REQUEST_NOT_FOUND, // 404 FR-CXL-03 ไม่มีคำขอนี้ หรือเป็นของผู้อื่น (T-20)
}

// CODE_FORBIDDEN_ROLE คือรหัสที่หลังบ้านตอบเมื่อเจ้าหน้าที่เปิดหน้านี้ (403, ASM-13, CON-TECH-03 CR-006)
// ลองใหม่ก็ไม่สำเร็จ จึงแสดง message จากหลังบ้านโดยไม่มีปุ่มลองใหม่ (plan.md ข้อ 4, ทีมตัดสิน 2569-11-12 ทาง ก, T-26)
const CODE_FORBIDDEN_ROLE = 'FORBIDDEN_ROLE'

// 401 ไม่ต้องจัดการที่หน้านี้: api/client.js แจ้ง App ให้แสดงหน้าเข้าสู่ระบบจำลอง แล้วกลับมาแท็บนี้ (IF-IDP-01, AC-CXL-07, T-21 ของ 001)

export default function MyPendingRequests({ client = api }) {
  const [items, setItems] = useState(null)
  const [loadFailed, setLoadFailed] = useState(false)
  const [forbidden, setForbidden] = useState(null) // message ของ 403 FORBIDDEN_ROLE (T-26)
  const [reloadKey, setReloadKey] = useState(0) // เพิ่มค่าเมื่อกดลองใหม่ ให้โหลดรายการอีกครั้ง
  const [target, setTarget] = useState(null) // คำขอที่เปิดกล่องยืนยันอยู่
  const [busy, setBusy] = useState(false)
  const [success, setSuccess] = useState(null) // ข้อความยกเลิกสำเร็จ ค้างจนเริ่มยกเลิกรายการถัดไป (ทีมตัดสิน 2569-11-11 T-07)
  const [cancelError, setCancelError] = useState(null) // ข้อความเมื่อยกเลิกไม่ได้ ค้างจนเริ่มยกเลิกรายการถัดไป แบบเดียวกับ success
  const [dialogError, setDialogError] = useState(null) // ข้อความในกล่องเมื่อระบบมีปัญหา (T-24)

  useEffect(() => {
    let active = true
    setItems(null)
    setLoadFailed(false)
    setForbidden(null)
    client
      .listMyPendingRequests()
      .then((res) => active && setItems(res.items))
      .catch((err) => {
        if (!active) return
        if (err?.code === CODE_FORBIDDEN_ROLE) {
          // เจ้าหน้าที่: แสดง message จากหลังบ้าน ไม่มีปุ่มลองใหม่ (ASM-13, T-26) ไม่มี message ใช้ข้อความโหลดไม่สำเร็จแทน
          setForbidden(err.serverMessage || MSG_LOAD_FAILED)
          return
        }
        // โหลดไม่สำเร็จ: ไม่แสดงรายการเก่าและไม่แสดงว่าไม่มีคำขอ (FR-CXL-06)
        setLoadFailed(true)
      })
    return () => {
      active = false
    }
  }, [client, reloadKey])

  // เปิดกล่องยืนยัน (FR-CXL-07)
  function openDialog(request) {
    setSuccess(null)
    setCancelError(null)
    setDialogError(null)
    setTarget(request)
  }

  // ยืนยันการยกเลิก (FR-CXL-01, FR-CXL-07) หลังบ้านเป็นผู้ตัดสินว่ายกเลิกได้ไหม (FR-CXL-04)
  async function confirmCancel(reason) {
    const request = target
    setBusy(true)
    setDialogError(null)
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
      } else {
        // FR-CXL-01, ASM-06: ระบบไม่ตอบสนอง สถานะคงเดิม คงแถวและกล่องไว้ให้ลองใหม่
        setDialogError(MSG_CANCEL_FAILED)
      }
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
      {forbidden ? (
        <p role="alert" className="text-red-700">
          {forbidden}
        </p>
      ) : loadFailed ? (
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
          error={dialogError}
        />
      )}
    </section>
  )
}
