// คำขอยืมที่ค้างอยู่ของอุปกรณ์ชิ้นหนึ่ง แสดงบนหน้ารายละเอียดให้เจ้าหน้าที่เห็น (FR-EQS-11, DOM-REQ-01, CR-005)
// แสดงลำดับที่ เวลาที่ส่ง และชื่อผู้ขอ เรียงตามลำดับที่ส่ง (queuePosition 1 = ส่งก่อน)
// ชื่อผู้ขอเป็นชื่อที่ผู้ขอพิมพ์เองในฟอร์มคำขอของ 002 (FR-BRQ-03 ก) ระบบไม่ได้ตรวจ จึงบอกไว้ใต้ตาราง
// component นี้แสดงอย่างเดียว ไม่เรียก API เอง หน้ารายละเอียดเป็นคนเรียกและตัดสินว่าจะแสดงหรือไม่
// requests ว่าง หรือไม่ใช่ array = ไม่แสดงอะไรเลย (spec ไม่บอกข้อความกรณีนี้ plan ข้อ 10 ของ 001)
// formatTime ส่งมาจากหน้ารายละเอียด ใช้รูปแบบเดียวกับกำหนดคืน "2569-10-10 16:00" (plan ข้อ 5)
export default function RequestQueue({ requests, formatTime }) {
  if (!Array.isArray(requests) || requests.length === 0) return null

  const sorted = [...requests].sort((a, b) => a.queuePosition - b.queuePosition)

  return (
    <section aria-labelledby="request-queue-heading" className="grid gap-2">
      <h3 id="request-queue-heading" className="text-base font-semibold text-slate-900">
        คำขอยืมที่ค้างอยู่
      </h3>
      <div className="overflow-x-auto">
        <table className="min-w-full text-left text-sm">
          <thead className="text-slate-600">
            <tr>
              <th scope="col" className="py-1 pr-4 font-medium">ลำดับที่</th>
              <th scope="col" className="py-1 pr-4 font-medium">เวลาที่ส่ง</th>
              <th scope="col" className="py-1 font-medium">ชื่อผู้ขอ</th>
            </tr>
          </thead>
          <tbody className="text-slate-900">
            {sorted.map((r) => (
              <tr key={r.requestId} className="border-t border-slate-200">
                <td className="py-1 pr-4 tabular-nums">{r.queuePosition}</td>
                <td className="py-1 pr-4 tabular-nums">{formatTime(r.submittedAt)}</td>
                <td className="py-1">{r.requesterName}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <p className="text-xs text-slate-600">ชื่อผู้ขอคือชื่อที่ผู้ขอพิมพ์เองในฟอร์มคำขอ ระบบไม่ได้ตรวจ</p>
    </section>
  )
}
