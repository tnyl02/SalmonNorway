import { useState } from 'react'

// กล่องยืนยันการยกเลิกคำขอ พร้อมช่องเหตุผล (FR-CXL-07, FR-CXL-08)
// request คือแถวที่จะยกเลิก onConfirm(reason) เรียกเมื่อกดยืนยัน onCancel เรียกเมื่อกดไม่ยืนยัน busy = กำลังส่งไปหลังบ้าน
// ช่องเหตุผลเก็บใน state ของกล่อง ปิดแล้วเปิดใหม่จึงว่างเสมอ (FR-CXL-07, ASM-13)
// การตรวจเหตุผลว่างหรือเกินและตัวนับ x/200 อยู่ใน T-10
export default function CancelDialog({ request, onConfirm, onCancel, busy = false }) {
  const [reason, setReason] = useState('')
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
            className="rounded border border-slate-300 p-2 text-slate-900"
          />
        </label>
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
            disabled={busy}
            onClick={() => onConfirm(reason)}
            className="rounded bg-red-700 px-3 py-1 text-sm text-white disabled:opacity-50"
          >
            ยืนยัน
          </button>
        </div>
      </div>
    </div>
  )
}
