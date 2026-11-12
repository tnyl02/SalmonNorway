// AC-EQS-08 (FR-EQS-08) ใช้ client จำลองแทนหลังบ้าน
//   Given ฐานข้อมูลไม่ตอบสนอง
//   When  เปิดหน้ารายการอุปกรณ์
//   Then  แสดงข้อความผิดพลาดพร้อมปุ่มลองใหม่ ไม่แสดงหน้าจอค้างและไม่แสดงข้อมูลเก่า
// กรณีต่อได้แต่ช้า (ต้องรอกี่วินาที) รอ Q-08 ดู T-26
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { ApiError } from '../api/client.js'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'

const DOWN_MESSAGE = 'ดึงข้อมูลไม่สำเร็จ กรุณาลองใหม่อีกครั้ง'
const MODELS = [
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 3, totalCount: 5 },
]

// client จำลองที่สลับได้ว่าฐานข้อมูลล่มหรือไม่ แบบเดียวกับหลังบ้านที่ตอบ 503 UPSTREAM_UNAVAILABLE (CON-TECH-03)
function switchableClient() {
  const state = { down: false, calls: 0 }
  return {
    state,
    listModels: async () => {
      state.calls++
      if (state.down) throw new ApiError(503, 'UPSTREAM_UNAVAILABLE', DOWN_MESSAGE)
      return { items: MODELS, total: MODELS.length }
    },
    listItems: async () => ({ items: [], total: 0 }),
  }
}

test('AC-EQS-08 ฐานข้อมูลไม่ตอบ แสดงข้อความผิดพลาดพร้อมปุ่มลองใหม่ ไม่ค้างที่ "กำลังโหลด"', async () => {
  // Given: ฐานข้อมูลไม่ตอบสนอง
  const client = switchableClient()
  client.state.down = true

  // When: เปิดหน้ารายการอุปกรณ์
  render(<EquipmentList client={client} />)

  // Then: ข้อความผิดพลาดจากหลังบ้าน + ปุ่มลองใหม่ และไม่ค้างที่ "กำลังโหลด..."
  expect((await screen.findByRole('alert')).textContent).toBe(DOWN_MESSAGE)
  expect(screen.getByRole('button', { name: 'ลองใหม่' })).toBeTruthy()
  expect(screen.queryByText('กำลังโหลด...')).toBeNull()
})

test('AC-EQS-08 เคยโหลดสำเร็จแล้วฐานข้อมูลล่ม ต้องไม่แสดงรายการเก่า และกดลองใหม่แล้วกลับมาได้', async () => {
  const client = switchableClient()
  render(<EquipmentList client={client} />)
  await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })

  // ฐานข้อมูลล่มระหว่างใช้งาน แล้วผู้ใช้เปลี่ยนตัวกรอง (หน้าจอต้องโหลดใหม่)
  client.state.down = true
  fireEvent.click(within(screen.getByRole('group', { name: 'สถานะ' })).getByRole('button', { name: 'ว่าง' }))

  // ไม่แสดงข้อมูลเก่าที่อาจไม่ถูกต้อง: รายการเดิมหายไปเหลือแค่ข้อผิดพลาด
  await screen.findByRole('alert')
  expect(screen.queryByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeNull()
  expect(screen.queryByText('ว่าง 3 / ทั้งหมด 5')).toBeNull()

  // ฐานข้อมูลกลับมา กดลองใหม่แล้วเห็นรายการอีกครั้ง
  client.state.down = false
  const callsBefore = client.state.calls
  fireEvent.click(screen.getByRole('button', { name: 'ลองใหม่' }))
  expect(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeTruthy()
  await waitFor(() => expect(client.state.calls).toBeGreaterThan(callsBefore))
  expect(screen.queryByRole('alert')).toBeNull()
})
