// TC-APV-16-2 (AC-APV-16 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-16-2 หลังบ้านตอบ APPROVAL_CUTOFF_PASSED แสดงข้อความ คำขอยังอยู่และยังปฏิเสธได้', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ และตอบ error APPROVAL_CUTOFF_PASSED เมื่ออนุมัติ
  const client = {
    ...fakeClient(),
    approveRequest: async () => Promise.reject(Object.assign(new Error('409'), { code: 'APPROVAL_CUTOFF_PASSED' })),
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "อนุมัติ" ที่คำขอแรก
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) แสดง "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้"
  expect((await screen.findByRole('alert')).textContent).toBe('เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้')

  // Then (2) คำขอนั้นยังอยู่ในรายการ และยังมีปุ่ม "ปฏิเสธ"
  const stillFirst = within(table).getAllByRole('row')[1]
  expect(within(stillFirst).getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
  expect(within(stillFirst).getByRole('button', { name: 'ปฏิเสธ' })).toBeTruthy()
})
