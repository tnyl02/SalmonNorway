// TC-APV-10-3 (AC-APV-10 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-10-3 หลังบ้านตอบ REQUEST_DECIDED_CONCURRENTLY แสดงข้อความให้รีเฟรชหน้า', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ และตอบ error รหัส REQUEST_DECIDED_CONCURRENTLY เมื่ออนุมัติ
  const client = {
    ...fakeClient(),
    approveRequest: async () => Promise.reject(Object.assign(new Error('409'), { code: 'REQUEST_DECIDED_CONCURRENTLY' })),
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "อนุมัติ" ที่คำขอแรก
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) แสดง "คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น กรุณารีเฟรชหน้า"
  expect((await screen.findByRole('alert')).textContent).toBe('คำขอนี้เพิ่งถูกพิจารณาโดยเจ้าหน้าที่คนอื่น กรุณารีเฟรชหน้า')
})
