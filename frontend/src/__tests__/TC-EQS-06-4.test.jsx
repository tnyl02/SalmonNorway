// TC-EQS-06-4 (AC-EQS-06 ขอบ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-06-4 ตารางรายชิ้นรุ่น B ชิ้นไม่พร้อมใช้งานไม่มีปุ่ม ชิ้นที่ว่างมีปุ่ม', async () => {
  // Given รุ่น B มี CS-PJ-0001 "ว่าง" และ CS-PJ-0002 "ไม่พร้อมใช้งาน" (client จำลอง)
  render(<EquipmentList client={eqsFakeClient()} />)

  // When นักศึกษาเปิดตารางรายชิ้นของรุ่น B
  fireEvent.click(await screen.findByRole('button', { name: /Projector ตัวอย่างรุ่น B/ }))
  const unavailable = (await screen.findByRole('button', { name: 'CS-PJ-0002' })).closest('tr')
  const available = screen.getByRole('button', { name: 'CS-PJ-0001' }).closest('tr')

  // Then (1) แถว CS-PJ-0002 แสดงป้าย "ไม่พร้อมใช้งาน" และไม่มีปุ่ม "ส่งคำขอยืม"
  expect(within(unavailable).getByText('ไม่พร้อมใช้งาน')).toBeTruthy()
  expect(within(unavailable).queryByRole('button', { name: 'ส่งคำขอยืม' })).toBeNull()
  // Then (2) แถว CS-PJ-0001 มีปุ่ม "ส่งคำขอยืม" (Q-16)
  expect(within(available).getByRole('button', { name: 'ส่งคำขอยืม' })).toBeTruthy()
})
