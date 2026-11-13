// TC-EQS-03-3 (AC-EQS-03 ขอบ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-03-3 ตัวกรองตรงพอดี 1 รุ่น แสดงรุ่น B 1 แถว ไม่แสดง "ไม่พบอุปกรณ์"', async () => {
  // Given client จำลองคืนรุ่นที่ตรงตัวกรองพอดี 1 รุ่น (สถานะ "ไม่พร้อมใช้งาน" มีแค่รุ่น B 1 ชิ้น)
  render(<EquipmentList client={eqsFakeClient()} />)
  await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })

  // When กด chip "ไม่พร้อมใช้งาน"
  fireEvent.click(within(screen.getByRole('group', { name: 'สถานะ' })).getByRole('button', { name: 'ไม่พร้อมใช้งาน' }))

  // Then (1) แสดงแถวรุ่น B 1 แถว
  await waitFor(() => expect(screen.queryByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeNull())
  expect(screen.getAllByRole('button', { name: /Projector ตัวอย่างรุ่น B/ })).toHaveLength(1)
  // Then (2) ไม่แสดง "ไม่พบอุปกรณ์"
  expect(screen.queryByText('ไม่พบอุปกรณ์')).toBeNull()
})
