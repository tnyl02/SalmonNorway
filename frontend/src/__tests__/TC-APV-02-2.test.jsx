// TC-APV-02-2 (AC-APV-02 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { LOADED_AT, fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-02-2 ปฏิเสธพร้อมเหตุผลแล้วส่ง requestId เหตุผล loadedAt และแถวนั้นหายจากรายการ', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ และตอบ 200 เมื่อปฏิเสธ
  const client = fakeClient()
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "ปฏิเสธ" ที่แถวที่สอง พิมพ์ "อุปกรณ์ต้องใช้สอบ" ในช่อง "เหตุผลในการปฏิเสธ" แล้วกด "ยืนยันการปฏิเสธ"
  const second = within(table).getAllByRole('row')[2]
  fireEvent.click(within(second).getByRole('button', { name: 'ปฏิเสธ' }))
  fireEvent.change(within(second).getByLabelText('เหตุผลในการปฏิเสธ'), { target: { value: 'อุปกรณ์ต้องใช้สอบ' } })
  fireEvent.click(within(second).getByRole('button', { name: 'ยืนยันการปฏิเสธ' }))

  // Then (1) client ถูกเรียกปฏิเสธ 1 ครั้ง ด้วย requestId ของแถวนั้น เหตุผลที่พิมพ์ และ loadedAt ที่ได้ตอนโหลดรายการ
  await waitFor(() => expect(client.calls.reject).toHaveLength(1))
  expect(client.calls.reject[0]).toEqual({ requestId: 22, reason: 'อุปกรณ์ต้องใช้สอบ', loadedAt: LOADED_AT })

  // Then (2) แถวนั้นหายจากรายการ
  await waitFor(() => expect(screen.queryByText('นักศึกษาตัวอย่าง ข')).toBeNull())

  // Then (3) อีกแถวยังอยู่
  const rows = within(table).getAllByRole('row').slice(1)
  expect(rows).toHaveLength(1)
  expect(within(rows[0]).getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
})
