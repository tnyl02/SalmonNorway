import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import ModelRow from '../../components/ModelRow.jsx'
import ItemTable from '../../components/ItemTable.jsx'

// หน้ารายการอุปกรณ์ (UC-01 main flow ขั้น 3-4, FR-EQS-01)
// แสดงรายการรุ่น 1 แถว = 1 รุ่น แล้วเมื่อเลือกรุ่นแสดงอุปกรณ์รายชิ้นตามรหัสครุภัณฑ์
// client ส่งเข้ามาได้ ตอน test ใช้ client จำลองแทนหลังบ้านจริง
export default function EquipmentList({ client = api }) {
  const [models, setModels] = useState(null)
  const [selectedId, setSelectedId] = useState(null)
  const [items, setItems] = useState(null)
  const [failed, setFailed] = useState(false)

  useEffect(() => {
    let active = true
    client
      .listModels()
      .then((res) => active && setModels(res.items))
      .catch(() => {
        // FR-EQS-08: ไม่แสดงข้อมูลเก่า ปุ่มลองใหม่และ test ของกรณีนี้อยู่ใน T-20
        if (active) {
          setModels(null)
          setFailed(true)
        }
      })
    return () => {
      active = false
    }
  }, [client])

  useEffect(() => {
    if (selectedId == null) return
    let active = true
    setItems(null)
    client
      .listItems(selectedId)
      .then((res) => active && setItems(res.items))
      .catch(() => {
        if (active) {
          setItems(null)
          setFailed(true)
        }
      })
    return () => {
      active = false
    }
  }, [client, selectedId])

  if (failed) {
    return <p role="alert" className="mt-4 text-red-700">ดึงข้อมูลอุปกรณ์ไม่สำเร็จ</p>
  }
  if (models == null) {
    return <p className="mt-4 text-slate-600">กำลังโหลด...</p>
  }

  const selected = models.find((m) => m.modelId === selectedId)

  return (
    <section className="mt-4 grid gap-4">
      <h2 className="text-lg font-semibold text-slate-900">รายการอุปกรณ์</h2>
      <ul className="grid gap-2">
        {models.map((m) => (
          <ModelRow key={m.modelId} model={m} selected={m.modelId === selectedId} onSelect={setSelectedId} />
        ))}
      </ul>

      {selected && (
        <div className="grid gap-2">
          <h3 className="font-medium text-slate-900">{selected.name}</h3>
          {items == null ? <p className="text-slate-600">กำลังโหลด...</p> : <ItemTable items={items} />}
        </div>
      )}
    </section>
  )
}
