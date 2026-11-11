import { useEffect, useState } from 'react'
import { handoverApi } from '../../api/handover.js'
import { TIME_ZONE } from '../../config.js'

// รายการคำขอที่ได้รับอนุมัติแล้วและยังไม่ส่งมอบทั้งหมด ให้เจ้าหน้าที่เลือกคำขอที่จะบันทึกการส่งมอบ (FR-HND-01, AC-HND-14, ASM-HND-05)
// ไม่มี mockup ทีมออกแบบเอง client ส่งเข้ามาได้ ตอน test ใช้ client จำลอง
export default function HandoverList({ client = handoverApi, onSelect }) {
  const [requests, setRequests] = useState(null)
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    let active = true
    client
      .listRequests()
      .then((rows) => active && setRequests(rows))
      .catch(() => active && setFailed(true))
    return () => {
      active = false
    }
  }, [client])

  if (failed) return <p role="alert" className="mt-4 text-red-700">ดึงรายการคำขอไม่สำเร็จ</p>
  if (requests == null) return <p className="mt-4 text-slate-600">กำลังโหลด...</p>

  return (
    <section className="mt-4 grid gap-2">
      <h2 className="text-lg font-semibold text-slate-900">คำขอที่รอส่งมอบ</h2>
      {requests.length === 0 ? (
        <p className="text-slate-600">ไม่มีคำขอที่รอส่งมอบ</p>
      ) : (
        <table className="w-full text-left text-sm">
          <thead>
            <tr className="border-b">
              <th className="py-1">รหัสครุภัณฑ์</th>
              <th className="py-1">รุ่น</th>
              <th className="py-1">ผู้ส่งคำขอ</th>
              <th className="py-1">เวลาที่ส่งคำขอ</th>
              <th className="py-1"></th>
            </tr>
          </thead>
          <tbody>
            {requests.map((r) => (
              <tr key={r.requestId} className="border-b">
                <td className="py-1">{r.assetCode}</td>
                <td className="py-1">{r.modelName}</td>
                <td className="py-1">{r.requesterStudentId}</td>
                <td className="py-1">{new Date(r.submittedAt).toLocaleString('th-TH', { timeZone: TIME_ZONE })}</td>
                <td className="py-1">
                  {onSelect && (
                    <button type="button" className="rounded bg-teal-700 px-2 py-1 text-white" onClick={() => onSelect(r)}>
                      เลือก
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
