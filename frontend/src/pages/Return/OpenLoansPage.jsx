// หน้ารายการการยืมที่ยังไม่คืนของภาควิชาเจ้าหน้าที่ (006-return) ไม่มี mockup ทีมออกแบบเอง ตาม plan.md ข้อ 4 หน้าจอ OpenLoansPage
// ป้าย "เลยกำหนดคืน" อยู่ใน T-17
import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import { TIMEZONE } from '../../config.js'

// formatDateTime แสดงวันที่และเวลาตามเขตเวลา Asia/Bangkok (glossary)
function formatDateTime(iso) {
  return new Date(iso).toLocaleString('th-TH', { timeZone: TIMEZONE, dateStyle: 'medium', timeStyle: 'short' })
}

// OpenLoansPage: แสดงการยืมที่ยังไม่คืน เลือกแถวแล้วส่งรหัสครุภัณฑ์ให้หน้าบันทึกคืน (FR-RET-01(ข))
export default function OpenLoansPage({ client = api, onSelect }) {
  const [loans, setLoans] = useState(null)
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    client.returns.openLoans().then((res) => {
      if (cancelled) return
      if (res.ok) setLoans(res.data)
      else setError(`โหลดรายการที่ยังไม่คืนไม่สำเร็จ (${res.data?.code ?? res.status})`)
    })
    return () => {
      cancelled = true
    }
  }, [client])

  return (
    <section className="mt-6 space-y-2" aria-label="รายการที่ยังไม่คืน">
      <h2 className="text-xl font-semibold">รายการที่ยังไม่คืน</h2>
      {error && (
        <p role="alert" className="rounded bg-red-50 p-3 text-red-800">
          {error}
        </p>
      )}
      {!error && !loans && <p>กำลังโหลด...</p>}
      {loans && loans.length === 0 && <p>ไม่มีการยืมที่ยังไม่คืน</p>}
      {loans && loans.length > 0 && (
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr>
                <th>รหัสครุภัณฑ์</th>
                <th>รุ่น</th>
                <th>ผู้ยืม</th>
                <th>กำหนดคืน</th>
                <th />
              </tr>
            </thead>
            <tbody>
              {loans.map((l) => (
                <tr key={l.loanId}>
                  <td>{l.assetCode}</td>
                  <td>{l.modelName}</td>
                  <td>{l.borrowerId}</td>
                  <td>{formatDateTime(l.dueAt)}</td>
                  <td>
                    <button
                      type="button"
                      onClick={() => onSelect?.(l.assetCode)}
                      aria-label={`เลือก ${l.assetCode}`}
                      className="text-teal-800 underline"
                    >
                      เลือก
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  )
}
