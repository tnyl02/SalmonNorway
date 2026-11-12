import { useState } from 'react'
import { handoverApi } from '../../api/handover.js'
import { HANDOVER_MIN_PHOTOS, HANDOVER_PHOTO_MAX_BYTES, HANDOVER_PHOTO_TYPES, TIME_ZONE } from '../../config.js'

// ตรวจรูปก่อนส่งเพื่อความสะดวก (FR-HND-06, ASM-HND-07) หลังบ้านตรวจซ้ำจากเนื้อไฟล์เสมอ
// ขนาดไม่เกิน HANDOVER_PHOTO_MAX_BYTES ต่อรูป ไม่จำกัดจำนวนรูป (Q-HND-10, T-28)
function photosValid(photos) {
  return (
    photos.length >= HANDOVER_MIN_PHOTOS &&
    photos.every((p) => HANDOVER_PHOTO_TYPES.includes(p.type) && p.size <= HANDOVER_PHOTO_MAX_BYTES)
  )
}

const MAX_MB = HANDOVER_PHOTO_MAX_BYTES / (1024 * 1024)

const fmt = (iso) => new Date(iso).toLocaleString('th-TH', { timeZone: TIME_ZONE })

// ฟอร์มบันทึกการส่งมอบของคำขอ 1 รายการ (UC-08) ไม่มี mockup ทีมออกแบบเอง
// - ชื่อและนามสกุลผู้ยืม เจ้าหน้าที่พิมพ์ตามบัตรนักศึกษาหรือบัตรประชาชน (FR-HND-03, ASM-HND-10)
// - ยืนยันว่าตรวจบัตรแล้วทุกครั้ง ไม่ติ๊ก = บันทึกไม่ได้ (FR-HND-03, ASM-HND-01)
// - แนบรูปสภาพอุปกรณ์ (FR-HND-06) หลังบ้านตรวจจำนวนและชนิดเสมอ
// - บันทึก แล้วแสดงเวลาที่ส่งมอบและกำหนดคืน (FR-HND-05, FR-HND-07)
// ผลการยืนยันผู้ยืม (plan.md ข้อ 4) verified / unavailable บันทึกต่อได้ mismatch บันทึกไม่ได้ (FR-HND-02, ASM-HND-04)
const VERIFY_TEXT = {
  verified: 'ยืนยันตัวตนผู้ยืมแล้ว',
  mismatch: 'ผู้ยืมไม่ใช่ผู้ส่งคำขอนี้',
  unavailable: 'เชื่อมต่อระบบยืนยันตัวตนไม่สำเร็จ ใช้ผลการตรวจบัตรแทน',
}

export default function HandoverForm({ request, client = handoverApi, onBack }) {
  // FR-HND-02, ASM-HND-09: ผู้ยืมล็อกอินบัญชีมหาวิทยาลัยบนเครื่องเจ้าหน้าที่ prototype ใช้ token จำลอง (Q-SH-03)
  const [credential, setCredential] = useState('')
  const [verify, setVerify] = useState(null) // { result, userId }
  const [verifyError, setVerifyError] = useState(null)
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [cardChecked, setCardChecked] = useState(false)
  const [photos, setPhotos] = useState([])
  const [saving, setSaving] = useState(false)
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)

  async function verifyBorrower() {
    setVerify(null)
    setVerifyError(null)
    try {
      setVerify(await client.verifyBorrower(request.requestId, credential))
    } catch (err) {
      setVerifyError(err.serverMessage ?? 'ยืนยันตัวตนผู้ยืมไม่สำเร็จ')
    }
  }

  // บันทึกได้เมื่อยืนยันผู้ยืมแล้ว (หรือระบบยืนยันตัวตนไม่ตอบ) และตรวจบัตรแล้ว หลังบ้านตรวจซ้ำเสมอ
  const canSave = cardChecked && (verify?.result === 'verified' || verify?.result === 'unavailable')

  async function submit(e) {
    e.preventDefault()
    if (!canSave) return // FR-HND-02, FR-HND-03
    if (!photosValid(photos)) {
      setError(`ต้องแนบรูปสภาพอุปกรณ์อย่างน้อย ${HANDOVER_MIN_PHOTOS} รูป เป็นไฟล์ JPG หรือ PNG ขนาดไม่เกิน ${MAX_MB} MB ต่อรูป`)
      return
    }
    const form = new FormData()
    form.append('borrowerCredential', credential)
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

      <fieldset className="grid gap-1 rounded border p-2">
        <legend className="px-1">ผู้ยืมล็อกอินบัญชีมหาวิทยาลัย</legend>
        <label className="grid gap-1">
          บัญชีผู้ยืม
          <input
            type="password"
            className="rounded border px-2 py-1"
            value={credential}
            onChange={(e) => {
              setCredential(e.target.value)
              setVerify(null)
            }}
          />
        </label>
        <button type="button" className="justify-self-start rounded border px-3 py-1" onClick={verifyBorrower} disabled={!credential}>
          ยืนยันตัวตนผู้ยืม
        </button>
        {verify && (
          <p role="status" className={verify.result === 'mismatch' ? 'text-red-700' : 'text-slate-700'}>
            {VERIFY_TEXT[verify.result]} {verify.userId && `(${verify.userId})`}
          </p>
        )}
        {verifyError && (
          <p role="alert" className="text-red-700">
            {verifyError}
          </p>
        )}
      </fieldset>

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
          disabled={!canSave || saving}
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
