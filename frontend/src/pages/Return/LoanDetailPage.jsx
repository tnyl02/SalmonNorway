// หน้ารายละเอียดการยืมของเจ้าหน้าที่ (006-return) ไม่มี mockup ทีมออกแบบเอง ตาม plan.md ข้อ 4 หน้าจอ LoanDetailPage
// แสดงเวลาคืนและสภาพ (FR-RET-06) ป้าย "คืนเกินกำหนด" อยู่ใน T-17
import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import { RETURN_CONDITIONS, TIMEZONE } from '../../config.js'

// formatDateTime แสดงวันที่และเวลาตามเขตเวลา Asia/Bangkok (glossary)
function formatDateTime(iso) {
  return new Date(iso).toLocaleString('th-TH', { timeZone: TIMEZONE, dateStyle: 'medium', timeStyle: 'short' })
}

// conditionLabel แปลงค่าสภาพเป็นคำที่แสดง จาก config.js (FR-RET-02)
function conditionLabel(value) {
  return RETURN_CONDITIONS.find((c) => c.value === value)?.label ?? value
}

// LoanDetailPage: รายละเอียดการยืม 1 รายการจาก GET /api/v1/returns/loans/{loanId} (FR-RET-06)
export default function LoanDetailPage({ loanId, client = api }) {
  const [detail, setDetail] = useState(null)
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    setDetail(null)
    setError('')
    client.returns.loanDetail(loanId).then((res) => {
      if (cancelled) return
      if (res.ok) setDetail(res.data)
      else setError(`โหลดรายละเอียดการยืมไม่สำเร็จ (${res.data?.code ?? res.status})`)
    })
    return () => {
      cancelled = true
    }
  }, [loanId, client])

  if (error) {
    return (
      <p role="alert" className="rounded bg-red-50 p-3 text-red-800">
        {error}
      </p>
    )
  }
  if (!detail) return <p>กำลังโหลด...</p>

  return (
    <section className="mt-6 space-y-2 rounded border p-3" aria-label="รายละเอียดการยืม">
      <h2 className="text-xl font-semibold">รายละเอียดการยืม</h2>
      <dl className="grid grid-cols-[auto_1fr] gap-x-3">
        <dt>รหัสครุภัณฑ์</dt>
        <dd>{detail.assetCode}</dd>
        <dt>ผู้ยืม</dt>
        <dd>{detail.borrowerId}</dd>
        <dt>กำหนดคืน</dt>
        <dd>{formatDateTime(detail.dueAt)}</dd>
        <dt>เวลาคืน</dt>
        <dd>{detail.returnedAt ? formatDateTime(detail.returnedAt) : 'ยังไม่คืน'}</dd>
        <dt>สภาพ</dt>
        <dd>{detail.condition ? conditionLabel(detail.condition) : '-'}</dd>
      </dl>
    </section>
  )
}
