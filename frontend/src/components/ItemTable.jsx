import StatusBadge from './StatusBadge.jsx'

// อุปกรณ์รายชิ้นของรุ่นที่เลือก ตามรหัสครุภัณฑ์ พร้อมสถานะล่าสุดของแต่ละชิ้น (FR-EQS-01, CON-DATA-03)
// กดรหัสครุภัณฑ์เพื่อเปิดรายละเอียด (FR-EQS-04, ทีมตัดสิน 2569-11-10 ตอนทำ T-11)
// ปุ่ม "ส่งคำขอยืม" แสดงเฉพาะชิ้นที่หลังบ้านบอก canRequestLoan (สถานะ "ว่าง") ซ่อนเมื่อถูกยืมหรือไม่พร้อมใช้งาน
// (FR-EQS-04, 05, 06, Q-16 ปุ่มอยู่ทั้งตารางรายชิ้นและหน้ารายละเอียด, T-34)
// กดแล้วส่ง assetCode ของชิ้นนั้นออกไปทาง onRequestLoan ฟอร์มของ 002 เป็นผู้รับ (FR-BRQ-01) ยังไม่ได้ต่อ ดู tasks.md T-34
export default function ItemTable({ items, onOpen, onRequestLoan }) {
  return (
    <table className="w-full text-left text-sm">
      <thead>
        <tr className="border-b border-slate-200 text-slate-600">
          <th className="py-2 pr-4 font-medium">รหัสครุภัณฑ์</th>
          <th className="py-2 pr-4 font-medium">สถานะ</th>
          <th className="py-2 font-medium" aria-hidden="true"></th>
        </tr>
      </thead>
      <tbody>
        {items.map((item) => (
          <tr key={item.assetCode} className="border-b border-slate-100">
            <td className="py-2 pr-4 font-mono">
              {onOpen ? (
                <button type="button" onClick={() => onOpen(item.assetCode)} className="text-teal-700 underline">
                  {item.assetCode}
                </button>
              ) : (
                item.assetCode
              )}
            </td>
            <td className="py-2 pr-4">
              <StatusBadge status={item.status} />
            </td>
            <td className="py-2">
              {item.canRequestLoan === true && (
                <RequestLoanButton assetCode={item.assetCode} onRequestLoan={onRequestLoan} />
              )}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

// ปุ่มส่งคำขอยืมของอุปกรณ์ 1 ชิ้น ใช้ทั้งตารางรายชิ้นและหน้ารายละเอียด (FR-EQS-06 ส่งข้อมูลอุปกรณ์ที่เลือกไปยัง UC-02, T-34)
// ส่งแค่ assetCode (plan.md ข้อ 5) ไม่ได้ส่ง onRequestLoan มา = กดแล้วยังไม่ไปไหน จนกว่าจะต่อกับฟอร์มของ 002
export function RequestLoanButton({ assetCode, onRequestLoan }) {
  return (
    <button
      type="button"
      onClick={() => onRequestLoan?.(assetCode)}
      className="rounded border border-teal-600 px-3 py-1 text-sm text-teal-700"
    >
      ส่งคำขอยืม
    </button>
  )
}
