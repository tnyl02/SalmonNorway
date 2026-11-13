// TC-EQS-04-4 (AC-EQS-04 ขอบ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-04-4 ตารางรายชิ้นของรุ่น A ชิ้นที่ถูกยืมแสดงป้าย "ถูกยืม" ไม่มีปุ่มส่งคำขอยืม', async () => {
  // Given รุ่น A มีชิ้นที่ถูกยืม (CS-NB-0004, CS-NB-0005)
  render(<EquipmentList client={eqsFakeClient()} />)

  // When นักศึกษาเปิดตารางรายชิ้นของรุ่น A
  fireEvent.click(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ }))

  for (const code of ['CS-NB-0004', 'CS-NB-0005']) {
    const row = (await screen.findByRole('button', { name: code })).closest('tr')
    // Then (1) ชิ้นที่ถูกยืมแสดงป้าย "ถูกยืม"
    expect(within(row).getByText('ถูกยืม')).toBeTruthy()
    // Then (2) ไม่มีปุ่ม "ส่งคำขอยืม" ในแถวนั้น (Q-16)
    expect(within(row).queryByRole('button', { name: 'ส่งคำขอยืม' })).toBeNull()
  }
  // Then (3) ตารางรายชิ้นไม่แสดงกำหนดคืน กำหนดคืนดูได้ที่หน้ารายละเอียด (ทีมตัดสิน 2569-11-12)
  const table = screen.getByRole('table')
  expect(within(table).queryByText('กำหนดคืน')).toBeNull()
  expect(within(table).queryByText(/\d{4}-\d{2}-\d{2} \d{2}:\d{2}/)).toBeNull()
})
