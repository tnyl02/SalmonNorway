import { useState } from 'react'
import { api } from '../../api/client.js'
import {
  ATTACHMENT_ALLOWED_TYPES,
  LABEL_SUBMIT_REQUEST,
  MSG_ATTACHMENT_HINT,
  MSG_LOAN_NORMAL_CONDITIONS,
  MSG_SUBMIT_FAILED,
  MSG_SUBMIT_SUCCESS,
  PURPOSE_OPTIONS,
  PURPOSE_OTHER,
} from '../../config.js'

// หน้าแบบฟอร์มส่งคำขอยืมอุปกรณ์ (UC-02 รวม UC-05 ยืมค้างคืน)
// ช่องที่ต้องกรอกตาม FR-BRQ-03 ข้อ ก: ชื่อ-นามสกุล, รหัสนักศึกษา, วันและเวลาที่เริ่มยืม, วันและเวลาที่คืน, วัตถุประสงค์
// วัตถุประสงค์เลือกจากรายการ 3 ค่า เลือก "อื่น ๆ" ต้องพิมพ์ระบุเพิ่ม (FR-BRQ-02, AC-BRQ-02, AC-BRQ-14)
// เลือก "ค้างคืน" ต้องแนบใบบันทึกข้อความตั้งแต่ตอนส่ง ส่วน "ปกติ" ต้องเห็นเงื่อนไขก่อนส่ง (FR-BRQ-10, AC-BRQ-10)
// ไม่มี mockup ทีมออกแบบเอง (Q-BRQ-14 ปิดแล้ว 2569-11-12) ข้อความที่ FR สั่งให้แสดงอยู่ใน config.js ที่เดียว
//
// หน้าจอไม่ตัดสินใจแทนหลังบ้าน: กดส่งได้เสมอ แล้วให้หลังบ้านตรวจและตอบรหัสตาม CON-TECH-03
// การแสดงข้อความของแต่ละรหัสและการชี้ช่องที่ผิดเป็นงานของ T-20 ที่นี่แสดง message ที่หลังบ้านส่งมาไปก่อน
export default function BorrowRequestForm({ client = api, assetCode: initialAssetCode = '' }) {
  const [assetCode, setAssetCode] = useState(initialAssetCode)
  const [requesterName, setRequesterName] = useState('')
  const [studentCode, setStudentCode] = useState('')
  const [loanType, setLoanType] = useState('normal')
  const [startAt, setStartAt] = useState('')
  const [endAt, setEndAt] = useState('')
  const [purpose, setPurpose] = useState('')
  const [purposeOther, setPurposeOther] = useState('')
  const [attachment, setAttachment] = useState(null)
  const [sending, setSending] = useState(false)
  const [created, setCreated] = useState(null)
  const [error, setError] = useState(null)

  const overnight = loanType === 'overnight'

  async function handleSubmit(event) {
    event.preventDefault()
    setSending(true)
    setError(null)
    setCreated(null)
    try {
      const result = await client.createLoanRequest(
        {
          assetCode,
          requesterName,
          studentCode,
          loanType,
          purpose,
          // ไม่ได้เลือก "อื่น ๆ" ไม่ต้องส่งข้อความที่พิมพ์ไว้ (FR-BRQ-02 ข้อ ข)
          purposeOther: purpose === PURPOSE_OTHER ? purposeOther : '',
          startAt: toBangkok(startAt),
          endAt: toBangkok(endAt),
        },
        overnight ? attachment : null,
      )
      setCreated(result)
    } catch (err) {
      setError(err?.serverMessage || MSG_SUBMIT_FAILED)
    } finally {
      setSending(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="mx-auto flex max-w-xl flex-col gap-4 p-4">
      <h2 className="text-lg font-semibold">{LABEL_SUBMIT_REQUEST}</h2>

      <Field id="assetCode" label="รหัสครุภัณฑ์">
        <input id="assetCode" value={assetCode} onChange={(e) => setAssetCode(e.target.value)} className={INPUT} />
      </Field>

      <Field id="requesterName" label="ชื่อ-นามสกุล">
        <input id="requesterName" value={requesterName} onChange={(e) => setRequesterName(e.target.value)} className={INPUT} />
      </Field>

      <Field id="studentCode" label="รหัสนักศึกษา">
        <input id="studentCode" value={studentCode} onChange={(e) => setStudentCode(e.target.value)} className={INPUT} />
      </Field>

      <fieldset className="flex flex-col gap-1">
        <legend className="text-sm">ประเภทการยืม</legend>
        <label className="flex items-center gap-2">
          <input type="radio" name="loanType" value="normal" checked={!overnight} onChange={() => setLoanType('normal')} />
          ปกติ
        </label>
        <label className="flex items-center gap-2">
          <input type="radio" name="loanType" value="overnight" checked={overnight} onChange={() => setLoanType('overnight')} />
          ค้างคืน
        </label>
      </fieldset>

      {/* ยืมปกติต้องแสดงเงื่อนไขก่อนส่ง (FR-BRQ-10 ข้อ ข, AC-BRQ-10) */}
      {!overnight && <p className="text-sm text-slate-600">{MSG_LOAN_NORMAL_CONDITIONS}</p>}

      <Field id="startAt" label="วันและเวลาที่เริ่มยืม">
        <input id="startAt" type="datetime-local" value={startAt} onChange={(e) => setStartAt(e.target.value)} className={INPUT} />
      </Field>

      <Field id="endAt" label="วันและเวลาที่คืน">
        <input id="endAt" type="datetime-local" value={endAt} onChange={(e) => setEndAt(e.target.value)} className={INPUT} />
      </Field>

      <Field id="purpose" label="วัตถุประสงค์การยืม">
        <select id="purpose" value={purpose} onChange={(e) => setPurpose(e.target.value)} className={INPUT}>
          <option value="">เลือกวัตถุประสงค์</option>
          {PURPOSE_OPTIONS.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </Field>

      {/* เลือก "อื่น ๆ" ต้องพิมพ์ระบุเพิ่ม จึงจะส่งได้ (FR-BRQ-02 ข้อ ก, AC-BRQ-14) */}
      {purpose === PURPOSE_OTHER && (
        <Field id="purposeOther" label="ระบุวัตถุประสงค์">
          <input id="purposeOther" value={purposeOther} onChange={(e) => setPurposeOther(e.target.value)} className={INPUT} />
        </Field>
      )}

      {/* ยืมค้างคืนต้องแนบใบบันทึกข้อความตั้งแต่ตอนส่งคำขอ (FR-BRQ-10 ข้อ ก, AC-BRQ-10) */}
      {overnight && (
        <Field id="attachment" label="ใบบันทึกข้อความ" hint={MSG_ATTACHMENT_HINT}>
          <input
            id="attachment"
            type="file"
            accept={ATTACHMENT_ALLOWED_TYPES.join(',')}
            onChange={(e) => setAttachment(e.target.files?.[0] ?? null)}
          />
        </Field>
      )}

      <button type="submit" disabled={sending} className="rounded bg-teal-600 px-4 py-2 text-white disabled:opacity-50">
        {LABEL_SUBMIT_REQUEST}
      </button>

      {created && <p role="status">{MSG_SUBMIT_SUCCESS(created.requestId)}</p>}
      {error && <p role="alert">{error}</p>}
    </form>
  )
}

const INPUT = 'rounded border border-slate-300 px-2 py-1'

// Field ครอบ label กับช่องกรอกให้ผูกกันด้วย id เพื่อให้ผู้ใช้คลิกป้ายแล้วโฟกัสช่องได้ (NFR-USE-01)
function Field({ id, label, hint, children }) {
  return (
    <div className="flex flex-col gap-1">
      <label htmlFor={id} className="text-sm">
        {label}
      </label>
      {children}
      {hint && <span className="text-xs text-slate-600">{hint}</span>}
    </div>
  )
}

// toBangkok แปลงค่าจากช่อง datetime-local เป็นเวลาแบบ RFC3339 โดยถือว่าเป็นเวลาไทย (glossary, config.TIME_ZONE)
// ช่อง datetime-local ไม่มีเขตเวลาติดมา ถ้าไม่ระบุ หลังบ้านจะตีความตามเขตเวลาของเครื่องที่รัน
// ซึ่งทำให้กฎ "คืนภายใน 16:00 น. ของวันเดียวกัน" เพี้ยนเมื่อผู้ใช้อยู่คนละเขตเวลา (DOM-LOAN-01)
export function toBangkok(value) {
  return value ? `${value}:00+07:00` : ''
}
