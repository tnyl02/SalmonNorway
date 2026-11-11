import { useState } from 'react'
import HandoverList from './HandoverList.jsx'
import HandoverForm from './HandoverForm.jsx'

// หน้าบันทึกการส่งมอบของเจ้าหน้าที่ (UC-08): เลือกคำขอจากรายการ (FR-HND-01) แล้วบันทึกในฟอร์ม (FR-HND-02 ถึง 07)
// ไม่มี mockup ทีมออกแบบเอง
export default function HandoverPage() {
  const [selected, setSelected] = useState(null)
  if (selected) return <HandoverForm request={selected} onBack={() => setSelected(null)} />
  return <HandoverList onSelect={setSelected} />
}
