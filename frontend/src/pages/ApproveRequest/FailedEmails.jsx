import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import { APPROVAL_EMAIL_KIND_LABELS, APPROVAL_MESSAGES, TIME_ZONE } from '../../config.js'

// แสดงวันเวลาตามเขตเวลา Asia/Bangkok (glossary)
const dateTime = new Intl.DateTimeFormat('th-TH', { timeZone: TIME_ZONE, dateStyle: 'medium', timeStyle: 'short' })

// รายการ "Email ที่ส่งไม่สำเร็จ" ของภาควิชาเจ้าหน้าที่ พร้อมปุ่มส่งซ้ำ (FR-APV-02 ง, ASM-08, AC-APV-14, T-34)
// แสดงสาเหตุและเวลาที่ลองส่งครั้งล่าสุด กดส่งซ้ำแล้ว Email กลับเข้าคิว แถวนั้นหายจากรายการ ไม่มีการพิจารณาใหม่
// หน้าจอไม่มี mockup ทีมออกแบบเอง (plan.md ข้อ 4) client ส่งเข้ามาได้ ตอน test ใช้ client จำลอง
export default function FailedEmails({ client = api }) {
  const [items, setItems] = useState(null)
  const [errorCode, setErrorCode] = useState(null) // โหลดรายการไม่สำเร็จ
  const [actionError, setActionError] = useState(null) // ส่งซ้ำไม่สำเร็จ รายการยังอยู่
  const [busyId, setBusyId] = useState(null)

  useEffect(() => {
    let active = true
    client
      .getFailedEmails()
      .then((res) => active && setItems(res.items))
      .catch((err) => active && setErrorCode(err?.code ?? 'UPSTREAM_UNAVAILABLE'))
    return () => {
      active = false
    }
  }, [client])

  async function resend(emailId) {
    setBusyId(emailId)
    setActionError(null)
    try {
      await client.resendEmail(emailId)
      setItems((cur) => cur.filter((e) => e.emailId !== emailId))
    } catch (err) {
      setActionError(err?.code ?? 'UPSTREAM_UNAVAILABLE')
    } finally {
      setBusyId(null)
    }
  }

  const message = (code) => APPROVAL_MESSAGES[code] ?? APPROVAL_MESSAGES.UPSTREAM_UNAVAILABLE

  if (errorCode) {
    return (
      <p role="alert" className="mt-4 text-red-700">
        {message(errorCode)}
      </p>
    )
  }
  if (items == null) return <p className="mt-4 text-slate-600">กำลังโหลด...</p>

  return (
    <section className="mt-8 grid gap-2">
      <h2 className="text-lg font-semibold text-slate-900">Email ที่ส่งไม่สำเร็จ</h2>
      {actionError && (
        <p role="alert" className="text-red-700">
          {message(actionError)}
        </p>
      )}
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="border-b">
            <th className="py-1">ผู้ขอ</th>
            <th className="py-1">อุปกรณ์</th>
            <th className="py-1">Email</th>
            <th className="py-1">สาเหตุ</th>
            <th className="py-1">เวลา</th>
            <th className="py-1"></th>
          </tr>
        </thead>
        <tbody>
          {items.map((e) => (
            <tr key={e.emailId} className="border-b">
              <td className="py-1">{e.requesterName}</td>
              <td className="py-1">
                {e.modelName} <span className="text-slate-500">({e.assetCode})</span>
              </td>
              <td className="py-1">{APPROVAL_EMAIL_KIND_LABELS[e.kind] ?? e.kind}</td>
              <td className="py-1">{e.lastError}</td>
              <td className="py-1">{dateTime.format(new Date(e.lastAttemptAt))}</td>
              <td className="py-1">
                <button
                  type="button"
                  onClick={() => resend(e.emailId)}
                  disabled={busyId != null}
                  className="rounded border border-teal-700 px-3 py-1 text-teal-700 disabled:opacity-50"
                >
                  ส่งซ้ำ
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}
