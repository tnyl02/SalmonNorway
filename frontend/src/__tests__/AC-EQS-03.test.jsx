// AC-EQS-03 (FR-EQS-03) ใช้ client จำลองแทนหลังบ้าน
//   Given ค้นหาด้วยคำที่ไม่มีในระบบ
//   When  กดค้นหา
//   Then  แจ้งว่าไม่พบอุปกรณ์ และมีทางเลือกล้างตัวกรองกลับไปดูรายการทั้งหมด
// ยังไม่มีช่องค้นหาด้วยคำ (รอ Q-06 ดู T-24) จึงใช้ตัวกรองที่ไม่ตรงกับชิ้นใดแทน ตามที่ทีมตัดสิน 2569-11-10
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'

const MODELS = [
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 3, totalCount: 5 },
  { modelId: 2, name: 'Projector ตัวอย่างรุ่น B', category: 'อุปกรณ์นำเสนอ', availableCount: 1, totalCount: 2 },
]

// client จำลอง: ไม่กรอง = ได้ทุกรุ่น / ประเภท "อุปกรณ์นำเสนอ" + สถานะ "ถูกยืม" = ไม่พบ
function fakeClient() {
  return {
    listModels: async (filter = {}) => {
      const none = (filter.categories ?? []).includes('อุปกรณ์นำเสนอ') && (filter.statuses ?? []).includes('borrowed')
      const items = none ? [] : MODELS
      return { items, total: items.length }
    },
    listItems: async () => ({ items: [], total: 0 }),
  }
}

test('AC-EQS-03 ไม่พบอุปกรณ์ แจ้งและล้างตัวกรองกลับไปเห็นรายการทั้งหมด', async () => {
  render(<EquipmentList client={fakeClient()} />)
  await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })

  // Given / When: กรองด้วยเงื่อนไขที่ไม่มีในระบบ
  fireEvent.click(within(screen.getByRole('group', { name: 'ประเภท' })).getByRole('button', { name: 'อุปกรณ์นำเสนอ' }))
  fireEvent.click(within(screen.getByRole('group', { name: 'สถานะ' })).getByRole('button', { name: 'ถูกยืม' }))

  // Then: แจ้งว่าไม่พบอุปกรณ์ ไม่ใช่ข้อความผิดพลาด
  const empty = await screen.findByRole('status')
  expect(within(empty).getByText('ไม่พบอุปกรณ์')).toBeTruthy()
  expect(screen.queryByRole('alert')).toBeNull()

  // Then: กดล้างตัวกรอง กลับไปเห็นรายการทั้งหมด และไม่มีตัวกรองค้าง
  fireEvent.click(within(empty).getByRole('button', { name: 'ล้างตัวกรอง' }))
  expect(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeTruthy()
  expect(screen.getByRole('button', { name: /Projector ตัวอย่างรุ่น B/ })).toBeTruthy()
  await waitFor(() => expect(screen.queryByText('ไม่พบอุปกรณ์')).toBeNull())
  const stillPressed = screen
    .queryAllByRole('button', { pressed: true })
    .filter((b) => b.closest('[role="group"]')) // เฉพาะ chip ตัวกรอง ไม่นับแถวรุ่น
  expect(stillPressed.map((b) => b.textContent)).toEqual([])
})
