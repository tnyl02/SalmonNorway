// AC-EQS-02 (FR-EQS-02) ส่วนกรองสถานะ ใช้ client จำลองแทนหลังบ้าน
//   Given อยู่ในหน้ารายการอุปกรณ์
//   When  กรองสถานะ "ว่าง"
//   Then  แสดงเฉพาะอุปกรณ์ที่มีสถานะว่างเท่านั้น
// ส่วน "ค้นหาด้วยคำว่า Notebook" รอ Q-06 (T-24) ทีมตัดสิน 2569-11-10 ว่ายังไม่ทำช่องค้นหา
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'

// ข้อมูลรายชิ้นจำลองแบบเดียวกับข้อมูลตั้งต้นของหลังบ้าน client จำลองกรองเองตามตัวกรองที่ได้รับ
const ITEMS = [
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0001', status: 'available' },
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0002', status: 'available' },
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0003', status: 'available' },
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0004', status: 'borrowed' },
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0005', status: 'borrowed' },
  { modelId: 2, name: 'Projector ตัวอย่างรุ่น B', category: 'อุปกรณ์นำเสนอ', assetCode: 'CS-PJ-0002', status: 'unavailable' },
]

function matches(it, filter = {}) {
  const cats = filter.categories ?? []
  const sts = filter.statuses ?? []
  return (cats.length === 0 || cats.includes(it.category)) && (sts.length === 0 || sts.includes(it.status))
}

function fakeClient() {
  const calls = []
  return {
    calls,
    listModels: async (filter) => {
      calls.push(['listModels', filter])
      const byModel = new Map()
      for (const it of ITEMS.filter((i) => matches(i, filter))) {
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
      const items = ITEMS.filter((i) => i.modelId === modelId && matches(i, filter)).map(({ assetCode, status }) => ({ assetCode, status }))
      return { items, total: items.length }
    },
  }
}

test('AC-EQS-02 กรองสถานะ "ว่าง" แล้วเห็นเฉพาะอุปกรณ์ที่ว่าง', async () => {
  // Given: อยู่ในหน้ารายการอุปกรณ์
  const client = fakeClient()
  render(<EquipmentList client={client} />)
  await screen.findByRole('button', { name: /Projector ตัวอย่างรุ่น B/ })

  // When: กรองสถานะ "ว่าง"
  const statusGroup = screen.getByRole('group', { name: 'สถานะ' })
  fireEvent.click(within(statusGroup).getByRole('button', { name: 'ว่าง' }))

  // Then: หน้าจอส่งตัวกรองไปหลังบ้าน และรุ่นที่ไม่มีชิ้นว่างหายไป
  await waitFor(() => expect(screen.queryByRole('button', { name: /Projector ตัวอย่างรุ่น B/ })).toBeNull())
  expect(client.calls).toContainEqual(['listModels', { categories: [], statuses: ['available'] }])
  const row = screen.getByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })
  expect(within(row).getByText('ว่าง 3 / ทั้งหมด 3')).toBeTruthy()

  // Then: รายชิ้นของรุ่นนั้นเห็นเฉพาะชิ้นที่ว่าง
  fireEvent.click(row)
  const table = await screen.findByRole('table')
  const bodyRows = within(table).getAllByRole('row').slice(1)
  expect(bodyRows).toHaveLength(3)
  bodyRows.forEach((r) => expect(within(r).getByText('ว่าง')).toBeTruthy())
  expect(within(table).queryByText('ถูกยืม')).toBeNull()
  expect(client.calls).toContainEqual(['listItems', 1, { categories: [], statuses: ['available'] }])
})

test('FR-EQS-02 กรองตามประเภท และกดซ้ำเพื่อยกเลิกตัวกรอง', async () => {
  const client = fakeClient()
  render(<EquipmentList client={client} />)
  await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })

  // ตัวเลือกประเภทมาจากรายการตอนยังไม่กรอง
  const catGroup = screen.getByRole('group', { name: 'ประเภท' })
  const chip = within(catGroup).getByRole('button', { name: 'อุปกรณ์นำเสนอ' })
  fireEvent.click(chip)
  await waitFor(() => expect(screen.queryByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeNull())
  expect(screen.getByRole('button', { name: /Projector ตัวอย่างรุ่น B/ })).toBeTruthy()

  // กดซ้ำ = ยกเลิก กลับมาเห็นทุกรุ่น
  fireEvent.click(chip)
  expect(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeTruthy()
})
