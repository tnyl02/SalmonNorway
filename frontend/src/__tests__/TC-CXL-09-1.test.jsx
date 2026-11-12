// TC-CXL-09-1 (AC-CXL-09 ปกติ) จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'
import { REQ_A1, fakeClient } from './helpers/cxlFakeClient.js'

test('TC-CXL-09-1 กดไม่ยืนยันแล้วไม่ยกเลิก กดยืนยันรอบสองแล้วยกเลิก', async () => {
  // Given อยู่ในหน้า "คำขอที่รออนุมัติของฉัน" มีคำขอ 1 รายการ client จำลองตอบ 200 เมื่อยกเลิก
  const client = fakeClient([REQ_A1])
  render(<MyPendingRequests client={client} />)
  const row = await screen.findByRole('listitem')

  // When กดยกเลิก กดไม่ยืนยัน
  fireEvent.click(within(row).getByRole('button', { name: 'ยกเลิก' }))
  let dialog = await screen.findByRole('dialog')
  fireEvent.click(within(dialog).getByRole('button', { name: 'ไม่ยืนยัน' }))

  // Then (1) หลังกดไม่ยืนยัน client ยังไม่ถูกเรียกยกเลิก และแถวยังอยู่
  await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull())
  expect(client.cancelRequest).not.toHaveBeenCalled()
  expect(screen.getAllByRole('listitem')).toHaveLength(1)

  // When กดยกเลิกอีกครั้ง พิมพ์เหตุผล แล้วกดยืนยัน
  fireEvent.click(within(screen.getByRole('listitem')).getByRole('button', { name: 'ยกเลิก' }))
  dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then (2) หลังกดยืนยันรอบสอง client ถูกเรียก 1 ครั้ง
  await waitFor(() => expect(client.cancelRequest).toHaveBeenCalledTimes(1))

  // Then (3) แถวหายจากหน้า
  await waitFor(() => expect(screen.queryAllByRole('listitem')).toHaveLength(0))
})
