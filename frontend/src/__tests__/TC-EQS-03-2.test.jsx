// TC-EQS-03-2 (AC-EQS-03 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-03-2 ไม่พบอุปกรณ์ แจ้งและล้างตัวกรองกลับไปเห็นรายการทั้งหมด', async () => {
  // Given client จำลองคืนรายการว่างเมื่อกรอง "อุปกรณ์นำเสนอ" + "ถูกยืม"
  render(<EquipmentList client={eqsFakeClient()} />)
  await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })

  // When กด chip ทั้งสอง แล้วกด "ล้างตัวกรอง"
  fireEvent.click(within(screen.getByRole('group', { name: 'ประเภท' })).getByRole('button', { name: 'อุปกรณ์นำเสนอ' }))
  fireEvent.click(within(screen.getByRole('group', { name: 'สถานะ' })).getByRole('button', { name: 'ถูกยืม' }))

  // Then (1) แจ้ง "ไม่พบอุปกรณ์"
  expect(await screen.findByText('ไม่พบอุปกรณ์')).toBeTruthy()
  // Then (2) มีทางเลือก "ล้างตัวกรอง"
  const clear = screen.getByRole('button', { name: 'ล้างตัวกรอง' })
  expect(clear).toBeTruthy()

  fireEvent.click(clear)
  // Then (3) กดแล้วกลับไปเห็นรายการทั้งหมด (รุ่น A และ B) และไม่มี chip ใดถูกเลือกค้าง
  expect(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeTruthy()
  expect(screen.getByRole('button', { name: /Projector ตัวอย่างรุ่น B/ })).toBeTruthy()
  await waitFor(() => expect(screen.queryByText('ไม่พบอุปกรณ์')).toBeNull())
  const pressed = screen.queryAllByRole('button', { pressed: true }).filter((b) => b.closest('[role="group"]'))
  expect(pressed).toHaveLength(0)
})
