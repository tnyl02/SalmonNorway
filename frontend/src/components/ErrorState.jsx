// แจ้งข้อผิดพลาดจากหลังบ้าน (CON-TECH-03, FR-EQS-08, FR-EQS-10)
// แสดงข้อความภาษาไทยที่หลังบ้านส่งมา (serverMessage) เช่น 503 UPSTREAM_UNAVAILABLE "ดึงข้อมูลไม่สำเร็จ กรุณาลองใหม่อีกครั้ง"
// ถ้าไม่มีข้อความ (เช่น ต่อหลังบ้านไม่ได้เลย) ใช้ข้อความเดิมของหน้าจอ
// มีปุ่ม "ลองใหม่" เมื่อหน้าที่เรียกส่ง onRetry มา (FR-EQS-08 "เปิดให้ลองใหม่", AC-EQS-08 "ปุ่มลองใหม่")
// ไม่แสดงข้อมูลเก่าใด ๆ หน้าที่เรียกต้องล้างข้อมูลเดิมก่อนแสดงส่วนนี้
export const DEFAULT_ERROR_MESSAGE = 'ดึงข้อมูลอุปกรณ์ไม่สำเร็จ'

export default function ErrorState({ error, onRetry }) {
  return (
    <div className="grid justify-items-start gap-2">
      <p role="alert" className="text-red-700">
        {error?.serverMessage || DEFAULT_ERROR_MESSAGE}
      </p>
      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="rounded border border-teal-600 px-3 py-1 text-sm text-teal-700"
        >
          ลองใหม่
        </button>
      )}
    </div>
  )
}
