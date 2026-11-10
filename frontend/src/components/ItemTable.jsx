import StatusBadge from './StatusBadge.jsx'

// อุปกรณ์รายชิ้นของรุ่นที่เลือก ตามรหัสครุภัณฑ์ พร้อมสถานะล่าสุดของแต่ละชิ้น (FR-EQS-01, CON-DATA-03)
export default function ItemTable({ items }) {
  return (
    <table className="w-full text-left text-sm">
      <thead>
        <tr className="border-b border-slate-200 text-slate-600">
          <th className="py-2 pr-4 font-medium">รหัสครุภัณฑ์</th>
          <th className="py-2 font-medium">สถานะ</th>
        </tr>
      </thead>
      <tbody>
        {items.map((item) => (
          <tr key={item.assetCode} className="border-b border-slate-100">
            <td className="py-2 pr-4 font-mono">{item.assetCode}</td>
            <td className="py-2">
              <StatusBadge status={item.status} />
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}
