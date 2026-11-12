// TC-APV-06-5 (AC-APV-06 ขอบ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-06-5 เหตุผลที่มีแต่ช่องว่างนับเป็นไม่ระบุ แสดงข้อความและไม่เรียก client', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ
  const client = fakeClient()
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "ปฏิเสธ" พิมพ์ "   " (ช่องว่างล้วน) แล้วกด "ยืนยันการปฏิเสธ"
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'ปฏิเสธ' }))
  fireEvent.change(within(first).getByLabelText('เหตุผลในการปฏิเสธ'), { target: { value: '   ' } })
  fireEvent.click(within(first).getByRole('button', { name: 'ยืนยันการปฏิเสธ' }))

  // Then (1) แสดง "กรุณาระบุเหตุผลในการปฏิเสธ"
  expect((await screen.findByRole('alert')).textContent).toBe('กรุณาระบุเหตุผลในการปฏิเสธ')

  // Then (2) client ไม่ถูกเรียกปฏิเสธ
  expect(client.calls.reject).toHaveLength(0)
})
