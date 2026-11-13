// client จำลองของ 001 สำหรับ test ที่มาจาก specs/001-equipment-status/test-cases.md (TC-EQS-*.test.jsx)
// ข้อมูลแบบเดียวกับข้อมูลตั้งต้นของหลังบ้าน (backend/db/seed/seed.sql ภาควิชา 1)
// กรองที่รายชิ้นแล้วนับตามตัวกรอง แบบเดียวกับหลังบ้าน (FR-EQS-02, T-32) และจดว่าหน้าจอเรียกด้วยตัวกรองอะไร
import { ApiError } from '../../api/client.js'

const A = { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา' }
const B = { modelId: 2, name: 'Projector ตัวอย่างรุ่น B', category: 'อุปกรณ์นำเสนอ' }

export const SEED_ITEMS = [
  { ...A, assetCode: 'CS-NB-0001', status: 'available' },
  { ...A, assetCode: 'CS-NB-0002', status: 'available' },
  { ...A, assetCode: 'CS-NB-0003', status: 'available' },
  { ...A, assetCode: 'CS-NB-0004', status: 'borrowed' },
  { ...A, assetCode: 'CS-NB-0005', status: 'borrowed' },
  { ...B, assetCode: 'CS-PJ-0001', status: 'available' },
  { ...B, assetCode: 'CS-PJ-0002', status: 'unavailable' },
]

function matches(item, filter = {}) {
  const cats = filter.categories ?? []
  const sts = filter.statuses ?? []
  return (cats.length === 0 || cats.includes(item.category)) && (sts.length === 0 || sts.includes(item.status))
}

// eqsFakeClient({ down }) down = true ทำให้ทุก request ตอบ 503 UPSTREAM_UNAVAILABLE (CON-TECH-03)
export function eqsFakeClient({ down = false } = {}) {
  const calls = []
  const fail = () => new ApiError(503, 'UPSTREAM_UNAVAILABLE', 'ดึงข้อมูลไม่สำเร็จ กรุณาลองใหม่อีกครั้ง')
  return {
    calls,
    listModels: async (filter) => {
      calls.push(['listModels', filter])
      if (down) throw fail()
      const byModel = new Map()
      for (const it of SEED_ITEMS.filter((i) => matches(i, filter))) {
        const m = byModel.get(it.modelId) ?? { modelId: it.modelId, name: it.name, category: it.category, availableCount: 0, totalCount: 0 }
        m.totalCount++
        if (it.status === 'available') m.availableCount++
        byModel.set(it.modelId, m)
      }
      const items = [...byModel.values()]
      return { items, total: items.length }
    },
    listItems: async (modelId, filter) => {
      calls.push(['listItems', modelId, filter])
      if (down) throw fail()
      const items = SEED_ITEMS.filter((i) => i.modelId === modelId && matches(i, filter)).map(({ assetCode, status }) => ({
        assetCode,
        status,
        canRequestLoan: status === 'available',
      }))
      return { items, total: items.length }
    },
    getItem: async (assetCode) => {
      if (down) throw fail()
      const it = SEED_ITEMS.find((i) => i.assetCode === assetCode)
      return { assetCode, status: it.status, canRequestLoan: it.status === 'available' }
    },
    listItemRequests: async () => [],
  }
}
