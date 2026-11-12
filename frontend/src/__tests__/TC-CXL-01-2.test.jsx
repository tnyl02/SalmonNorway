// TC-CXL-01-2 (AC-CXL-01 ปกติ) จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'
import { REQ_A1, fakeClient } from './helpers/cxlFakeClient.js'

test('TC-CXL-01-2 ยกเลิกพร้อมเหตุผล แสดงข้อความสำเร็จ และแถวหายจากหน้า', async () => {
  // Given อยู่ในหน้า "คำขอที่รออนุมัติของฉัน" มีคำขอ 1 รายการ client จำลองตอบ 200 เมื่อยกเลิก
  const client = fakeClient([REQ_A1])
  render(<MyPendingRequests client={client} />)
  const row = await screen.findByRole('listitem')

  // When กดยกเลิก พิมพ์เหตุผล แล้วกดยืนยัน
  fireEvent.click(within(row).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then (1) client ถูกเรียกยกเลิก 1 ครั้งด้วย requestId และเหตุผลที่พิมพ์
  await waitFor(() => expect(client.cancelRequest).toHaveBeenCalledTimes(1))
  expect(client.cancelRequest).toHaveBeenCalledWith(101, 'ไม่ต้องใช้อุปกรณ์แล้ว')

  // Then (2) แสดง "ยกเลิกคำขอยืม <ชื่อรุ่น> (<รหัสครุภัณฑ์>) เรียบร้อยแล้ว" โดยแทนค่าจริงของแถวนั้น
  expect(await screen.findByText('ยกเลิกคำขอยืม Notebook ตัวอย่างรุ่น A (CS-NB-0001) เรียบร้อยแล้ว')).toBeTruthy()

  // Then (3) แถวนั้นหายจากหน้า
  await waitFor(() => expect(screen.queryAllByRole('listitem')).toHaveLength(0))
})
