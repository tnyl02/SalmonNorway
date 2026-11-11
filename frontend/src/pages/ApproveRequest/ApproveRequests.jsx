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
// ปุ่มอนุมัติ (T-09) / ปฏิเสธพร้อมช่องเหตุผล (T-11) / หน้าจอไม่มี mockup ทีมออกแบบเอง (plan.md ข้อ 4)
// client ส่งเข้ามาได้ ตอน test ใช้ client จำลองแทนหลังบ้านจริง
export default function ApproveRequests({ client = api }) {
  const [list, setList] = useState(null) // { items, loadedAt } จาก GET /v1/approvals/requests
  const [errorCode, setErrorCode] = useState(null) // โหลดรายการไม่สำเร็จ
  const [actionError, setActionError] = useState(null) // อนุมัติไม่สำเร็จ แสดงเหนือรายการ รายการยังอยู่
  const [busyId, setBusyId] = useState(null) // คำขอที่กำลังส่งอยู่ กันกดซ้ำ
  const [rejectingId, setRejectingId] = useState(null) // คำขอที่เปิดช่องเหตุผลการปฏิเสธอยู่ (T-11)
  const [reason, setReason] = useState('')

  // ส่งผลการพิจารณา 1 คำขอ ส่ง loadedAt ของรายการที่เห็นอยู่ (plan.md ข้อ 4)
  // สำเร็จแล้วเอาคำขอนั้น และคำขอที่หลังบ้านปฏิเสธตาม (FR-APV-01 จ) ออกจากรายการ
  // ไม่สำเร็จแสดงข้อความตามรหัส (FR-APV-01 ซ, CON-TECH-03) และไม่เปลี่ยนรายการ
  async function decide(requestId, send) {
    setBusyId(requestId)
    setActionError(null)
    try {
      const res = await send()
      const gone = new Set([res.requestId, ...(res.rejectedRequestIds ?? [])])
      setList((cur) => ({ ...cur, items: cur.items.filter((r) => !gone.has(r.requestId)) }))
      setRejectingId(null)
      setReason('')
    } catch (err) {
      setActionError(err?.code ?? 'UPSTREAM_UNAVAILABLE')
    } finally {
      setBusyId(null)
    }
  }

  // อนุมัติ (FR-APV-01 ค, AC-APV-01)
  const approve = (requestId) => decide(requestId, () => client.approveRequest(requestId, list.loadedAt))

  // ปฏิเสธพร้อมเหตุผลที่เจ้าหน้าที่พิมพ์ (FR-APV-01 ง, AC-APV-02) ข้อความเมื่อไม่ระบุเหตุผลเป็นงานของ T-13
  const reject = (requestId) =>
    decide(requestId, () => client.rejectRequest(requestId, reason, list.loadedAt))

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
      {actionError && (
        <p role="alert" className="text-red-700">
          {APPROVAL_MESSAGES[actionError] ?? APPROVAL_MESSAGES.UPSTREAM_UNAVAILABLE}
        </p>
      )}
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="border-b">
            <th className="py-1">ผู้ขอ</th>
            <th className="py-1">อุปกรณ์</th>
            <th className="py-1">วัตถุประสงค์</th>
            <th className="py-1">ช่วงเวลาที่ต้องการยืม</th>
            <th className="py-1">เวลาที่ส่ง</th>
            <th className="py-1"></th>
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
              <td className="py-1">
                {/* ปุ่มอนุมัติ (T-09) การซ่อนปุ่มแถวที่ไม่ใช่คำขอแรกเป็นงานของ T-18 หลังบ้านตรวจซ้ำเสมอ (T-17) */}
                <button
                  type="button"
                  onClick={() => approve(r.requestId)}
                  disabled={busyId != null}
                  className="rounded bg-teal-700 px-3 py-1 text-white disabled:opacity-50"
                >
                  อนุมัติ
                </button>{' '}
                {rejectingId === r.requestId ? (
                  // ช่องเหตุผลการปฏิเสธ (FR-APV-01 ง) ทีมออกแบบเอง ไม่มี mockup
                  <span className="mt-1 grid gap-1">
                    <label className="grid gap-1">
                      เหตุผลในการปฏิเสธ
                      <textarea
                        value={reason}
                        onChange={(e) => setReason(e.target.value)}
                        rows={2}
                        className="rounded border border-slate-300 p-1"
                      />
                    </label>
                    <span className="flex gap-1">
                      <button
                        type="button"
                        onClick={() => reject(r.requestId)}
                        disabled={busyId != null}
                        className="rounded bg-red-700 px-3 py-1 text-white disabled:opacity-50"
                      >
                        ยืนยันการปฏิเสธ
                      </button>
                      <button
                        type="button"
                        onClick={() => {
                          setRejectingId(null)
                          setReason('')
                        }}
                        className="rounded border border-slate-400 px-3 py-1"
                      >
                        ยกเลิก
                      </button>
                    </span>
                  </span>
                ) : (
                  <button
                    type="button"
                    onClick={() => {
                      setRejectingId(r.requestId)
                      setReason('')
                    }}
                    disabled={busyId != null}
                    className="rounded border border-red-700 px-3 py-1 text-red-700 disabled:opacity-50"
                  >
                    ปฏิเสธ
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}
