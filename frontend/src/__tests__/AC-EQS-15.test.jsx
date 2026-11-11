// AC-EQS-15 (FR-EQS-01, FR-EQS-02) ใช้ client จำลองแทนหลังบ้าน
//   Given รุ่นหนึ่งมี 5 ชิ้น ว่าง 3 ถูกยืม 2
//   When  กรองสถานะ "ถูกยืม"
//   Then  รุ่นนั้นแสดง 1 แถว และตัวเลขนับเฉพาะชิ้นที่ตรงตัวกรอง (2 ชิ้น)
// ข้อความกำกับตัวเลขคงรูปแบบ "ว่าง x / ทั้งหมด y" ไปก่อน ตามที่ทีมตัดสิน 2569-11-10 รอ Q-21
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

test('AC-EQS-15 กรอง "ถูกยืม" แล้วรุ่นนั้นแสดง 1 แถว และนับเฉพาะ 2 ชิ้นที่ตรง', async () => {
  // Given: รุ่น A 5 ชิ้น ว่าง 3 ถูกยืม 2
  render(<EquipmentList client={fakeClient()} />)
  await screen.findByText('ว่าง 3 / ทั้งหมด 5')

  // When: กรองสถานะ "ถูกยืม"
  const statusGroup = screen.getByRole('group', { name: 'สถานะ' })
  fireEvent.click(within(statusGroup).getByRole('button', { name: 'ถูกยืม' }))

  // Then: รุ่นนั้น 1 แถว และตัวเลขนับเฉพาะชิ้นที่ตรงตัวกรอง (2 ชิ้น)
  await waitFor(() => expect(screen.queryByText('ว่าง 3 / ทั้งหมด 5')).toBeNull())
  const rows = screen.getAllByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })
  expect(rows).toHaveLength(1)
  expect(within(rows[0]).getByText('ว่าง 0 / ทั้งหมด 2')).toBeTruthy()
})
