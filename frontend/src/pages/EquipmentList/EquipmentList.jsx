import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import ModelRow from '../../components/ModelRow.jsx'
import ItemTable from '../../components/ItemTable.jsx'
import SearchBar from '../../components/SearchBar.jsx'
import EmptyState from '../../components/EmptyState.jsx'
import ErrorState from '../../components/ErrorState.jsx'
import EquipmentDetail from '../EquipmentDetail/EquipmentDetail.jsx'

const NO_FILTER = { categories: [], statuses: [] }

// หน้ารายการอุปกรณ์ (UC-01 main flow ขั้น 3-6, alternative 6a และ exception E2, FR-EQS-01, FR-EQS-02, FR-EQS-03, FR-EQS-08)
// แสดงรายการรุ่น 1 แถว = 1 รุ่น กรองตามประเภทและสถานะ แล้วเมื่อเลือกรุ่นแสดงอุปกรณ์รายชิ้นตามรหัสครุภัณฑ์
// รายชิ้นใช้ตัวกรองเดียวกับรายการรุ่น (plan.md ข้อ 6.3)
// client ส่งเข้ามาได้ ตอน test ใช้ client จำลองแทนหลังบ้านจริง
// onRequestLoan(assetCode) ส่งต่อให้ปุ่มส่งคำขอยืมทั้งตารางรายชิ้นและหน้ารายละเอียด (FR-EQS-06, T-34)
export default function EquipmentList({ client = api, onRequestLoan }) {
  const [filter, setFilter] = useState(NO_FILTER)
  const [categories, setCategories] = useState([])
  const [models, setModels] = useState(null)
  const [selectedId, setSelectedId] = useState(null)
  const [items, setItems] = useState(null)
  const [error, setError] = useState(null) // ดึงข้อมูลไม่สำเร็จ (FR-EQS-08)
  const [reloadKey, setReloadKey] = useState(0) // เพิ่มทีละ 1 เมื่อกดลองใหม่ ให้โหลดใหม่ทั้งรายการรุ่นและรายชิ้น
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
      .catch((err) => {
        // FR-EQS-08: ไม่แสดงข้อมูลเก่าที่อาจไม่ถูกต้อง ล้างรายการเดิมทิ้งก่อนแสดงข้อผิดพลาด
        if (active) {
          setModels(null)
          setItems(null)
          setError(err)
        }
      })
    return () => {
      active = false
    }
  }, [client, filter, reloadKey])

  useEffect(() => {
    if (selectedId == null) return
    let active = true
    setItems(null)
    client
      .listItems(selectedId, filter)
      .then((res) => active && setItems(res.items))
      .catch((err) => {
        if (active) {
          setModels(null)
          setItems(null)
          setError(err)
        }
      })
    return () => {
      active = false
    }
  }, [client, selectedId, filter, reloadKey])

  // FR-EQS-08: กดลองใหม่ = ล้างข้อผิดพลาดแล้วโหลดใหม่ตามตัวกรองเดิม
  const retry = () => {
    setError(null)
    setReloadKey((k) => k + 1)
  }

  if (error) {
    return (
      <section className="mt-4 grid gap-4">
        <h2 className="text-lg font-semibold text-slate-900">รายการอุปกรณ์</h2>
        <ErrorState error={error} onRetry={retry} />
      </section>
    )
  }

  // เปิดรายละเอียดแทนรายการ ไม่ล้างตัวกรองและรุ่นที่เลือก กดกลับแล้วเห็นรายการเดิม
  if (openAsset) {
    return (
      <EquipmentDetail
        assetCode={openAsset}
        client={client}
        onBack={() => setOpenAsset(null)}
        onRequestLoan={onRequestLoan}
      />
    )
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
          {items == null ? <p className="text-slate-600">กำลังโหลด...</p> : (
            <ItemTable items={items} onOpen={setOpenAsset} onRequestLoan={onRequestLoan} />
          )}
        </div>
      )}
    </section>
  )
}
