// TC-EQS-15-2 (AC-EQS-15 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-15-2 กด chip "ถูกยืม" แถวรุ่น A มี 1 แถว และแสดงจำนวน 2', async () => {
  // Given client จำลองที่นับตามตัวกรอง ข้อมูลแบบรุ่น A
  render(<EquipmentList client={eqsFakeClient()} />)
  await screen.findByText('ว่าง 3 / ทั้งหมด 5')

  // When กด chip "ถูกยืม"
  fireEvent.click(within(screen.getByRole('group', { name: 'สถานะ' })).getByRole('button', { name: 'ถูกยืม' }))
  await waitFor(() => expect(screen.queryByText('ว่าง 3 / ทั้งหมด 5')).toBeNull())

  // Then (1) แถวรุ่น A มี 1 แถว
  const rows = await screen.findAllByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })
  expect(rows).toHaveLength(1)
  // Then (2) แถวแสดงจำนวน 2
  expect(rows[0].textContent).toMatch(/(^|\D)2(\D|$)/)
  // Then (3) ข้อความกำกับตัวเลขเมื่อกรองสถานะที่ไม่ใช่ "ว่าง": ยังไม่ตรวจ เพราะรอ Q-21
})
