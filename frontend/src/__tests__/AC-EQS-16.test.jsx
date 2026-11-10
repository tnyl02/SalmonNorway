// AC-EQS-16 (FR-EQS-02) ใช้ client จำลองแทนหลังบ้าน
//   Given มีอุปกรณ์สถานะ ว่าง ถูกยืม และไม่พร้อมใช้งาน
//   When  เลือกตัวกรองสถานะ "ว่าง" และ "ถูกยืม" พร้อมกัน
//   Then  แสดงเฉพาะชิ้นที่ว่างหรือถูกยืม ไม่แสดงชิ้นที่ไม่พร้อมใช้งาน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'

const ITEMS = [
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0001', status: 'available' },
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0002', status: 'available' },
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0003', status: 'available' },
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0004', status: 'borrowed' },
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', assetCode: 'CS-NB-0005', status: 'borrowed' },
  { modelId: 2, name: 'Projector ตัวอย่างรุ่น B', category: 'อุปกรณ์นำเสนอ', assetCode: 'CS-PJ-0001', status: 'available' },
  { modelId: 2, name: 'Projector ตัวอย่างรุ่น B', category: 'อุปกรณ์นำเสนอ', assetCode: 'CS-PJ-0002', status: 'unavailable' },
]

function matches(it, filter = {}) {
  const cats = filter.categories ?? []
  const sts = filter.statuses ?? []
  return (cats.length === 0 || cats.includes(it.category)) && (sts.length === 0 || sts.includes(it.status))
}

// client จำลองกรองที่รายชิ้นแล้วนับตามตัวกรอง แบบเดียวกับหลังบ้าน (T-32)
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

test('AC-EQS-16 เลือก "ว่าง" และ "ถูกยืม" พร้อมกัน แล้วไม่เห็นชิ้นที่ไม่พร้อมใช้งาน', async () => {
  // Given: Projector รุ่น B มีชิ้นว่าง 1 และไม่พร้อมใช้งาน 1
  const client = fakeClient()
  render(<EquipmentList client={client} />)
  await screen.findByText('ว่าง 1 / ทั้งหมด 2')

  // When: เลือก "ว่าง" และ "ถูกยืม" พร้อมกัน
  const statusGroup = screen.getByRole('group', { name: 'สถานะ' })
  fireEvent.click(within(statusGroup).getByRole('button', { name: 'ว่าง' }))
  fireEvent.click(within(statusGroup).getByRole('button', { name: 'ถูกยืม' }))

  // Then: ทั้งสองตัวกรองถูกเลือกอยู่ และส่งไปหลังบ้านพร้อมกัน
  expect(within(statusGroup).getByRole('button', { name: 'ว่าง' }).getAttribute('aria-pressed')).toBe('true')
  expect(within(statusGroup).getByRole('button', { name: 'ถูกยืม' }).getAttribute('aria-pressed')).toBe('true')
  await waitFor(() =>
    expect(client.calls).toContainEqual(['listModels', { categories: [], statuses: ['available', 'borrowed'] }]),
  )

  // Then: Projector นับแค่ชิ้นที่ว่าง และรายชิ้นไม่มีชิ้นที่ไม่พร้อมใช้งาน
  const projector = await screen.findByRole('button', { name: /Projector ตัวอย่างรุ่น B/ })
  await waitFor(() => expect(within(projector).getByText('ว่าง 1 / ทั้งหมด 1')).toBeTruthy())
  fireEvent.click(projector)
  const table = await screen.findByRole('table')
  expect(within(table).getByText('CS-PJ-0001')).toBeTruthy()
  expect(within(table).queryByText('CS-PJ-0002')).toBeNull()
  expect(within(table).queryByText('ไม่พร้อมใช้งาน')).toBeNull()
})
