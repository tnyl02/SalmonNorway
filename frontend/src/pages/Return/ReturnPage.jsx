// หน้าบันทึกการรับคืนอุปกรณ์ของเจ้าหน้าที่ (006-return, UC-09)
// ไม่มี mockup ทีมออกแบบเอง (spec ไม่มี UI ID) ตาม plan.md ข้อ 4 หน้าจอ ReturnPage
// หลังบ้านเป็นผู้ตรวจทุกกฎ หน้านี้แค่ส่งข้อมูลและแสดงผล ข้อความ error แยกตามกรณีอยู่ใน T-10, T-13, T-22, T-24
import { useState } from 'react'
import { api } from '../../api/client.js'
import { RETURN_CONDITIONS, RETURN_PHOTO_TYPES, STATUS_LABELS, TIMEZONE } from '../../config.js'

// formatDateTime แสดงวันที่และเวลาตามเขตเวลา Asia/Bangkok (glossary)
function formatDateTime(iso) {
  return new Date(iso).toLocaleString('th-TH', { timeZone: TIMEZONE, dateStyle: 'medium', timeStyle: 'short' })
}

// ReturnPage: กรอกรหัสครุภัณฑ์ (FR-RET-01(ก)) เลือกสภาพ (FR-RET-02) แนบรูป (FR-RET-05) แล้วบันทึก (FR-RET-03)
export default function ReturnPage({ client = api }) {
  const [assetCode, setAssetCode] = useState('')
  const [loan, setLoan] = useState(null) // การยืมที่ยังไม่คืนที่หาเจอ
  const [condition, setCondition] = useState('')
  const [photos, setPhotos] = useState([])
  const [result, setResult] = useState(null) // ผลการบันทึกที่สำเร็จ
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)

  // หาการยืมที่ยังไม่คืนจากรหัสครุภัณฑ์ (FR-RET-01(ก))
  async function handleLookup(e) {
    e.preventDefault()
    setBusy(true)
    setError('')
    setLoan(null)
    setResult(null)
    const res = await client.returns.lookup(assetCode.trim())
    setBusy(false)
    if (res.ok) {
      setLoan(res.data)
      setCondition('')
      setPhotos([])
    } else {
      setError(`ค้นหาไม่สำเร็จ (${res.data?.code ?? res.status})`)
    }
  }

  // บันทึกการรับคืน ไม่ส่งเวลา เวลาคืนมาจากระบบ (FR-RET-02, ASM-04)
  async function handleSave(e) {
    e.preventDefault()
    setBusy(true)
    setError('')
    const res = await client.returns.record({ assetCode: loan.assetCode, condition, photos })
    setBusy(false)
    if (res.ok) {
      setResult(res.data)
      setLoan(null)
    } else {
      setError(`บันทึกไม่สำเร็จ (${res.data?.code ?? res.status})`)
    }
  }

  return (
    <section className="mt-6 space-y-4">
      <h2 className="text-xl font-semibold">บันทึกการรับคืนอุปกรณ์</h2>

      <form onSubmit={handleLookup} className="flex flex-wrap items-end gap-2">
        <label className="flex flex-col">
          <span>รหัสครุภัณฑ์</span>
          <input
            className="rounded border px-2 py-1"
            value={assetCode}
            onChange={(e) => setAssetCode(e.target.value)}
            required
          />
        </label>
        <button type="submit" disabled={busy} className="rounded bg-teal-700 px-3 py-1 text-white">
          ค้นหา
        </button>
      </form>

      {loan && (
        <form onSubmit={handleSave} className="space-y-3 rounded border p-3" aria-label="ฟอร์มบันทึกการรับคืน">
          <dl className="grid grid-cols-[auto_1fr] gap-x-3">
            <dt>รหัสครุภัณฑ์</dt>
            <dd>{loan.assetCode}</dd>
            <dt>รุ่น</dt>
            <dd>{loan.modelName}</dd>
            <dt>ผู้ยืม</dt>
            <dd>{loan.borrowerId}</dd>
            <dt>กำหนดคืน</dt>
            <dd>{formatDateTime(loan.dueAt)}</dd>
          </dl>

          <fieldset>
            <legend>สภาพ</legend>
            {RETURN_CONDITIONS.map((c) => (
              <label key={c.value} className="mr-4">
                <input
                  type="radio"
                  name="condition"
                  value={c.value}
                  checked={condition === c.value}
                  onChange={() => setCondition(c.value)}
                />{' '}
                {c.label}
              </label>
            ))}
          </fieldset>

          <label className="flex flex-col">
            <span>รูปสภาพอุปกรณ์</span>
            <input
              type="file"
              multiple
              accept={RETURN_PHOTO_TYPES.join(',')}
              onChange={(e) => setPhotos(Array.from(e.target.files ?? []))}
            />
          </label>

          <button type="submit" disabled={busy} className="rounded bg-teal-700 px-3 py-1 text-white">
            บันทึกการรับคืน
          </button>
        </form>
      )}

      {result && (
        <p role="status" className="rounded bg-teal-50 p-3">
          บันทึกการรับคืนแล้ว สถานะ: {STATUS_LABELS[result.itemStatus] ?? result.itemStatus}
        </p>
      )}

      {error && (
        <p role="alert" className="rounded bg-red-50 p-3 text-red-800">
          {error}
        </p>
      )}
    </section>
  )
}
