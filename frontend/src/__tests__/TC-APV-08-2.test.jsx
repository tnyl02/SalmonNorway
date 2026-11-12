// TC-APV-08-2 (AC-APV-08 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-08-2 หลังบ้านตอบ ITEM_NOT_AVAILABLE แสดงข้อความ และคำขอยังอยู่', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ และตอบ error ITEM_NOT_AVAILABLE เมื่ออนุมัติ
  const client = {
    ...fakeClient(),
    approveRequest: async () => Promise.reject(Object.assign(new Error('409'), { code: 'ITEM_NOT_AVAILABLE' })),
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "อนุมัติ" ที่คำขอแรก
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) แสดง "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้"
  expect((await screen.findByRole('alert')).textContent).toBe('อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้')

  // Then (2) คำขอนั้นยังอยู่ในรายการ
  expect(screen.getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
})
