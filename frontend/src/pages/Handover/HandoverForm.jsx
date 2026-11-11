import { useState } from 'react'
import { handoverApi } from '../../api/client.js'
import { TIME_ZONE } from '../../config.js'

const fmt = (iso) => new Date(iso).toLocaleString('th-TH', { timeZone: TIME_ZONE })

// ฟอร์มบันทึกการส่งมอบของคำขอ 1 รายการ (UC-08) ไม่มี mockup ทีมออกแบบเอง
// - ชื่อและนามสกุลผู้ยืม เจ้าหน้าที่พิมพ์ตามบัตรนักศึกษาหรือบัตรประชาชน (FR-HND-03, ASM-HND-10)
// - ยืนยันว่าตรวจบัตรแล้วทุกครั้ง ไม่ติ๊ก = บันทึกไม่ได้ (FR-HND-03, ASM-HND-01)
// - แนบรูปสภาพอุปกรณ์ (FR-HND-06) หลังบ้านตรวจจำนวนและชนิดเสมอ
// - บันทึก แล้วแสดงเวลาที่ส่งมอบและกำหนดคืน (FR-HND-05, FR-HND-07)
export default function HandoverForm({ request, client = handoverApi, onBack }) {
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [cardChecked, setCardChecked] = useState(false)
  const [photos, setPhotos] = useState([])
  const [saving, setSaving] = useState(false)
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)

  async function submit(e) {
    e.preventDefault()
    if (!cardChecked) return // FR-HND-03
    const form = new FormData()
    form.append('borrowerFirstName', firstName)
    form.append('borrowerLastName', lastName)
    form.append('cardChecked', 'true')
    for (const p of photos) form.append('photos', p)
    setSaving(true)
    setError(null)
    try {
      setResult(await client.handover(request.requestId, form))
    } catch (err) {
      setError(err.serverMessage ?? 'บันทึกการส่งมอบไม่สำเร็จ')
    } finally {
      setSaving(false)
    }
  }

  if (result) {
    return (
      <section className="mt-4 grid gap-2" aria-label="ผลการบันทึกการส่งมอบ">
        <h2 className="text-lg font-semibold text-slate-900">บันทึกการส่งมอบแล้ว</h2>
        <p>รหัสครุภัณฑ์ {result.assetCode}</p>
        <p>เวลาที่ส่งมอบ {fmt(result.handedOverAt)}</p>
        <p>กำหนดคืน {fmt(result.dueAt)}</p>
        {onBack && (
          <button type="button" className="justify-self-start rounded border px-3 py-1" onClick={onBack}>
            กลับไปรายการคำขอ
          </button>
        )}
      </section>
    )
  }

  return (
    <form className="mt-4 grid gap-3" onSubmit={submit} aria-label="บันทึกการส่งมอบ">
      <h2 className="text-lg font-semibold text-slate-900">บันทึกการส่งมอบ {request.assetCode}</h2>
      <p className="text-sm text-slate-600">ผู้ส่งคำขอ {request.requesterStudentId}</p>

      <label className="grid gap-1">
        ชื่อผู้ยืม (ตามบัตร)
        <input className="rounded border px-2 py-1" value={firstName} onChange={(e) => setFirstName(e.target.value)} required />
      </label>
      <label className="grid gap-1">
        นามสกุลผู้ยืม (ตามบัตร)
        <input className="rounded border px-2 py-1" value={lastName} onChange={(e) => setLastName(e.target.value)} required />
      </label>
      <label className="flex items-center gap-2">
        <input type="checkbox" checked={cardChecked} onChange={(e) => setCardChecked(e.target.checked)} />
        ตรวจบัตรนักศึกษาหรือบัตรประชาชนของผู้ยืมแล้ว
      </label>
      <label className="grid gap-1">
        รูปสภาพอุปกรณ์
        <input type="file" multiple accept="image/jpeg,image/png" onChange={(e) => setPhotos([...e.target.files])} />
      </label>

      {error && (
        <p role="alert" className="text-red-700">
          {error}
        </p>
      )}
      <div className="flex gap-2">
        <button
          type="submit"
          className="rounded bg-teal-700 px-3 py-1 text-white disabled:opacity-50"
          disabled={!cardChecked || saving}
        >
          บันทึกการส่งมอบ
        </button>
        {onBack && (
          <button type="button" className="rounded border px-3 py-1" onClick={onBack}>
            กลับ
          </button>
        )}
      </div>
    </form>
  )
}
