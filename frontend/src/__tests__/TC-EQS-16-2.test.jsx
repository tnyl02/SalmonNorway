// TC-EQS-16-2 (AC-EQS-16 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-16-2 เลือก "ว่าง" และ "ถูกยืม" พร้อมกัน รายชิ้นรุ่น B ไม่มี CS-PJ-0002', async () => {
  // Given client จำลองที่กรองตามตัวกรองที่ได้รับ ข้อมูลแบบข้อมูลตั้งต้น
  const client = eqsFakeClient()
  render(<EquipmentList client={client} />)
  await screen.findByRole('button', { name: /Projector ตัวอย่างรุ่น B/ })

  // When กด chip "ว่าง" แล้วกด "ถูกยืม" แล้วกดแถวรุ่น B
  const group = screen.getByRole('group', { name: 'สถานะ' })
  fireEvent.click(within(group).getByRole('button', { name: 'ว่าง' }))
  fireEvent.click(within(group).getByRole('button', { name: 'ถูกยืม' }))

  // Then (1) chip ทั้งสองถูกเลือกอยู่พร้อมกัน
  expect(within(group).getByRole('button', { name: 'ว่าง' }).getAttribute('aria-pressed')).toBe('true')
  expect(within(group).getByRole('button', { name: 'ถูกยืม' }).getAttribute('aria-pressed')).toBe('true')
  // Then (2) หน้าจอส่งตัวกรองสถานะ 2 ค่าไปหลังบ้าน
  await waitFor(() => expect(client.calls).toContainEqual(['listModels', { categories: [], statuses: ['available', 'borrowed'] }]))

  fireEvent.click(await screen.findByRole('button', { name: /Projector ตัวอย่างรุ่น B/ }))
  const table = await screen.findByRole('table')
  // Then (3) ตารางรายชิ้นไม่มี CS-PJ-0002 และไม่มีป้าย "ไม่พร้อมใช้งาน"
  await waitFor(() => expect(within(table).queryByText('CS-PJ-0002')).toBeNull())
  expect(within(table).queryByText('ไม่พร้อมใช้งาน')).toBeNull()
})
