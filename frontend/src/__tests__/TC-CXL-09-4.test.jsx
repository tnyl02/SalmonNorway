// TC-CXL-09-4 (AC-CXL-09 ขอบ) จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'
import { REQ_A1, fakeClient } from './helpers/cxlFakeClient.js'

test('TC-CXL-09-4 ยกเลิกรายการสุดท้ายแล้วแสดง "ไม่มีคำขอที่รออนุมัติ"', async () => {
  // Given มีคำขอ 1 รายการ (ตาม AC)
  const client = fakeClient([REQ_A1])
  render(<MyPendingRequests client={client} />)
  const row = await screen.findByRole('listitem')

  // When กดยกเลิก พิมพ์เหตุผล แล้วกดยืนยัน
  fireEvent.click(within(row).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then (1) ไม่มีแถวเหลือ
  await waitFor(() => expect(screen.queryAllByRole('listitem')).toHaveLength(0))

  // Then (2) แสดงข้อความ "ไม่มีคำขอที่รออนุมัติ"
  expect(await screen.findByText('ไม่มีคำขอที่รออนุมัติ')).toBeTruthy()
})
