// หน้าหลักของ 006-return สำหรับเจ้าหน้าที่ รวมหน้าบันทึกการรับคืนกับรายการที่ยังไม่คืนไว้ด้วยกัน
// ย้ายมาจาก App.jsx (T-15) ตอนรวม main 2569-11-12 เพื่อให้ App มีแค่ 1 แถวใน PAGES เหมือนฟีเจอร์อื่น
import { useState } from 'react'
import OpenLoansPage from './OpenLoansPage.jsx'
import ReturnPage from './ReturnPage.jsx'

// ReturnsHome: เลือกแถวจากรายการที่ยังไม่คืน แล้วเปิด ReturnPage ใหม่พร้อมรหัสนั้น (FR-RET-01(ข))
// เปลี่ยน key เพื่อให้ ReturnPage เริ่มใหม่และค้นหารหัสที่เลือกทันที
export default function ReturnsHome({ client }) {
  const [selected, setSelected] = useState({ assetCode: '', key: 0 })
  const clientProps = client ? { client } : {}

  return (
    <>
      <ReturnPage key={selected.key} initialAssetCode={selected.assetCode} {...clientProps} />
      <OpenLoansPage onSelect={(assetCode) => setSelected((s) => ({ assetCode, key: s.key + 1 }))} {...clientProps} />
    </>
  )
}
