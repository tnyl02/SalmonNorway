import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import ModelRow from '../../components/ModelRow.jsx'
import ItemTable from '../../components/ItemTable.jsx'
import SearchBar from '../../components/SearchBar.jsx'
import EmptyState from '../../components/EmptyState.jsx'
import EquipmentDetail from '../EquipmentDetail/EquipmentDetail.jsx'

const NO_FILTER = { categories: [], statuses: [] }

// หน้ารายการอุปกรณ์ (UC-01 main flow ขั้น 3-6 และ alternative 6a, FR-EQS-01, FR-EQS-02, FR-EQS-03)
// แสดงรายการรุ่น 1 แถว = 1 รุ่น กรองตามประเภทและสถานะ แล้วเมื่อเลือกรุ่นแสดงอุปกรณ์รายชิ้นตามรหัสครุภัณฑ์
// รายชิ้นใช้ตัวกรองเดียวกับรายการรุ่น (plan.md ข้อ 6.3)
// client ส่งเข้ามาได้ ตอน test ใช้ client จำลองแทนหลังบ้านจริง
export default function EquipmentList({ client = api }) {
  const [filter, setFilter] = useState(NO_FILTER)
  const [categories, setCategories] = useState([])
  const [models, setModels] = useState(null)
  const [selectedId, setSelectedId] = useState(null)
  const [items, setItems] = useState(null)
  const [failed, setFailed] = useState(false)
  const [openAsset, setOpenAsset] = useState(null) // รหัสครุภัณฑ์ที่เปิดดูรายละเอียดอยู่ (FR-EQS-04)

  useEffect(() => {
    let active = true
    setModels(null)
    client
      .listModels(filter)
      .then((res) => {
        if (!active) return
        setModels(res.items)
        // ตัวเลือกประเภทมาจากรายการตอนยังไม่กรอง ไม่มี endpoint ประเภทแยกใน plan.md ข้อ 5
        if (filter === NO_FILTER) setCategories([...new Set(res.items.map((m) => m.category))])
      })
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
  }, [client, filter])

  useEffect(() => {
    if (selectedId == null) return
    let active = true
    setItems(null)
    client
      .listItems(selectedId, filter)
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
  }, [client, selectedId, filter])

  if (failed) {
    return <p role="alert" className="mt-4 text-red-700">ดึงข้อมูลอุปกรณ์ไม่สำเร็จ</p>
  }

  // เปิดรายละเอียดแทนรายการ ไม่ล้างตัวกรองและรุ่นที่เลือก กดกลับแล้วเห็นรายการเดิม
  if (openAsset) {
    return <EquipmentDetail assetCode={openAsset} client={client} onBack={() => setOpenAsset(null)} />
  }

  const selected = models?.find((m) => m.modelId === selectedId)

  return (
    <section className="mt-4 grid gap-4">
      <h2 className="text-lg font-semibold text-slate-900">รายการอุปกรณ์</h2>
      <SearchBar categories={categories} filter={filter} onChange={setFilter} />

      {models == null ? (
        <p className="text-slate-600">กำลังโหลด...</p>
      ) : models.length === 0 ? (
        // FR-EQS-03: ไม่พบ แจ้งและให้ล้างตัวกรองกลับไปดูรายการทั้งหมด
        <EmptyState
          onClear={() => {
            setSelectedId(null)
            setFilter(NO_FILTER)
          }}
        />
      ) : (
        <ul className="grid gap-2">
          {models.map((m) => (
            <ModelRow key={m.modelId} model={m} selected={m.modelId === selectedId} onSelect={setSelectedId} />
          ))}
        </ul>
      )}

      {selected && (
        <div className="grid gap-2">
          <h3 className="font-medium text-slate-900">{selected.name}</h3>
          {items == null ? <p className="text-slate-600">กำลังโหลด...</p> : <ItemTable items={items} onOpen={setOpenAsset} />}
        </div>
      )}
    </section>
  )
}
