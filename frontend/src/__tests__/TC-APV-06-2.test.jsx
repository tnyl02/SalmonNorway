// TC-APV-06-2 (AC-APV-06 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-06-2 ปฏิเสธโดยไม่พิมพ์เหตุผลแสดงข้อความ ไม่เรียก client และแถวยังอยู่', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ
  const client = fakeClient()
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "ปฏิเสธ" ที่แถวหนึ่ง ไม่พิมพ์อะไร แล้วกด "ยืนยันการปฏิเสธ"
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'ปฏิเสธ' }))
  fireEvent.click(within(first).getByRole('button', { name: 'ยืนยันการปฏิเสธ' }))

  // Then (1) แสดง "กรุณาระบุเหตุผลในการปฏิเสธ"
  expect((await screen.findByRole('alert')).textContent).toBe('กรุณาระบุเหตุผลในการปฏิเสธ')

  // Then (2) client ไม่ถูกเรียกปฏิเสธ
  expect(client.calls.reject).toHaveLength(0)

  // Then (3) แถวนั้นยังอยู่ในรายการ
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
  expect(screen.getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
})
