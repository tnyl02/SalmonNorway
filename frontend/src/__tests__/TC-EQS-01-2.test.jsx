// TC-EQS-01-2 (AC-EQS-01 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-01-2 แถวรุ่น A แสดงประเภทและ "ว่าง 3 / ทั้งหมด 5" แล้วเห็น 5 ชิ้นพร้อมป้ายสถานะ', async () => {
  // Given client จำลองคืนรุ่น A ว่าง 3 / ทั้งหมด 5 และรายชิ้น 5 ชิ้น (ว่าง 3 ถูกยืม 2)
  render(<EquipmentList client={eqsFakeClient()} />)

  // When เปิดหน้ารายการอุปกรณ์ แล้วกดแถวรุ่น A
  const rows = await screen.findAllByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })

  // Then (1) มีแถวรุ่น A 1 แถว
  expect(rows).toHaveLength(1)
  // Then (2) แถวแสดง "คอมพิวเตอร์พกพา"
  expect(within(rows[0]).getByText('คอมพิวเตอร์พกพา')).toBeTruthy()
  // Then (3) แถวแสดง "ว่าง 3 / ทั้งหมด 5"
  expect(within(rows[0]).getByText('ว่าง 3 / ทั้งหมด 5')).toBeTruthy()

  fireEvent.click(rows[0])
  const table = await screen.findByRole('table')
  // Then (4) ตารางรายชิ้น 5 แถว แต่ละแถวมีรหัสครุภัณฑ์และป้าย "ว่าง" หรือ "ถูกยืม"
  const bodyRows = within(table).getAllByRole('row').slice(1)
  expect(bodyRows).toHaveLength(5)
  bodyRows.forEach((r) => {
    expect(r.textContent).toMatch(/CS-NB-000[1-5]/)
    expect(within(r).queryByText('ว่าง') ?? within(r).queryByText('ถูกยืม')).toBeTruthy()
  })
})
