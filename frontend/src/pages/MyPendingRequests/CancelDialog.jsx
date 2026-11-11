import { useState } from 'react'
import { CANCEL_REASON_MAX_LEN, LABEL_RETRY } from '../../config.js'

// countReason นับความยาวเหตุผลตาม code point แบบเดียวกับหลังบ้าน เช่น "ที่" = 3 (FR-CXL-08, ASM-13)
export function countReason(text) {
  return [...text].length
}

// reasonProblem บอกว่าเหตุผลใช้ไม่ได้เพราะอะไร: 'empty' ว่างหรือมีแต่ช่องว่าง, 'too_long' เกิน CANCEL_REASON_MAX_LEN, null ใช้ได้ (FR-CXL-08)
export function reasonProblem(text) {
  if (text.trim() === '') return 'empty'
  if (countReason(text) > CANCEL_REASON_MAX_LEN) return 'too_long'
  return null
}

// กล่องยืนยันการยกเลิกคำขอ พร้อมช่องเหตุผล (FR-CXL-07, FR-CXL-08)
// request คือแถวที่จะยกเลิก onConfirm(reason) เรียกเมื่อกดยืนยัน onCancel เรียกเมื่อกดไม่ยืนยัน busy = กำลังส่งไปหลังบ้าน
// ช่องเหตุผลเก็บใน state ของกล่อง ปิดแล้วเปิดใหม่จึงว่างเสมอ (FR-CXL-07, ASM-13)
// ปุ่มยืนยันกดไม่ได้เมื่อเหตุผลว่างหรือเกิน และแสดงตัวนับ "x/200" ใต้ช่องเหตุผล (FR-CXL-08, ASM-13)
// หน้าจอกันไว้เพื่อความสะดวก กฎจริงอยู่ที่หลังบ้าน (T-09 ตอบ 400)
// error = ข้อความเมื่อยกเลิกไม่สำเร็จเพราะระบบมีปัญหา แสดงในกล่อง เหตุผลที่พิมพ์ยังอยู่ ปุ่มยืนยันเปลี่ยนเป็น "ลองใหม่" (FR-CXL-01, ASM-06)
export default function CancelDialog({ request, onConfirm, onCancel, busy = false, error = null }) {
  const [reason, setReason] = useState('')
  const count = countReason(reason)
  const invalid = reasonProblem(reason) != null
  const counterId = `cancel-reason-count-${request.requestId}`
  const titleId = `cancel-dialog-title-${request.requestId}`

  return (
    <div className="fixed inset-0 z-10 grid place-items-center bg-slate-900/40 p-4">
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleId}
        className="grid w-full max-w-md gap-3 rounded-lg bg-white p-4 shadow-lg"
      >
        <h3 id={titleId} className="font-semibold text-slate-900">
          ยืนยันการยกเลิกคำขอยืม {request.modelName} ({request.assetCode})
        </h3>
        <label className="grid gap-1 text-sm text-slate-700">
          เหตุผลการยกเลิก
          <textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            rows={3}
            aria-describedby={counterId}
            className="rounded border border-slate-300 p-2 text-slate-900"
          />
        </label>
        <p
          id={counterId}
          className={count > CANCEL_REASON_MAX_LEN ? 'text-right text-xs text-red-700' : 'text-right text-xs text-slate-500'}
        >
          {`${count}/${CANCEL_REASON_MAX_LEN}`}
        </p>
        {error && (
          <p role="alert" className="rounded border border-red-300 bg-red-50 p-2 text-sm text-red-800">
            {error}
          </p>
        )}
        <div className="flex justify-end gap-2">
          {/* ไม่ยืนยัน: ปิดกล่องโดยไม่ยกเลิกและไม่เรียก API (FR-CXL-07) */}
          <button
            type="button"
            disabled={busy}
            onClick={onCancel}
            className="rounded border border-slate-300 px-3 py-1 text-sm text-slate-700 disabled:opacity-50"
          >
            ไม่ยืนยัน
          </button>
          <button
            type="button"
            disabled={busy || invalid}
            onClick={() => onConfirm(reason)}
            className="rounded bg-red-700 px-3 py-1 text-sm text-white disabled:opacity-50"
          >
            {error ? LABEL_RETRY : 'ยืนยัน'}
          </button>
        </div>
      </div>
    </div>
  )
}
