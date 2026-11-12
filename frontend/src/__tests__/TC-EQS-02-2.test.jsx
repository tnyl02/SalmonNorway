// TC-EQS-02-2 (AC-EQS-02 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-02-2 กด chip "ว่าง" ส่งตัวกรองไปหลังบ้าน แถวรุ่น A "ว่าง 3 / ทั้งหมด 3" และรายชิ้นมีแต่ "ว่าง"', async () => {
  // Given client จำลองที่กรองตามตัวกรองที่ได้รับ ข้อมูลแบบข้อมูลตั้งต้น
  const client = eqsFakeClient()
  render(<EquipmentList client={client} />)
  await screen.findByText('ว่าง 3 / ทั้งหมด 5')

  // When กด chip สถานะ "ว่าง" แล้วกดแถวรุ่น A
  fireEvent.click(within(screen.getByRole('group', { name: 'สถานะ' })).getByRole('button', { name: 'ว่าง' }))

  // Then (1) หน้าจอส่งตัวกรองสถานะ available ไปหลังบ้าน
  await waitFor(() => expect(client.calls).toContainEqual(['listModels', { categories: [], statuses: ['available'] }]))
  // Then (2) แถวรุ่น A แสดง "ว่าง 3 / ทั้งหมด 3"
  const row = await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })
  await waitFor(() => expect(within(row).getByText('ว่าง 3 / ทั้งหมด 3')).toBeTruthy())

  fireEvent.click(row)
  const table = await screen.findByRole('table')
  // Then (3) ตารางรายชิ้นมีแต่ป้าย "ว่าง" ไม่มี "ถูกยืม"
  expect(within(table).getAllByText('ว่าง').length).toBeGreaterThan(0)
  expect(within(table).queryByText('ถูกยืม')).toBeNull()
})
