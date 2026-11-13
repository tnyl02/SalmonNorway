// TC-APV-20-7 (AC-APV-20 ขอบ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-20-7 รหัสที่ไม่มีในตาราง FR-APV-01 (ซ) แสดงข้อความเดียวกับ UPSTREAM_UNAVAILABLE', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ และตอบ error รหัส UNAUTHENTICATED (ไม่มีในตาราง) เมื่ออนุมัติ
  const client = {
    ...fakeClient(),
    approveRequest: async () => Promise.reject(Object.assign(new Error('401'), { code: 'UNAUTHENTICATED' })),
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "อนุมัติ" ที่คำขอแรก
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) แสดงข้อความเดียวกับ UPSTREAM_UNAVAILABLE
  expect((await screen.findByRole('alert')).textContent).toBe('ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง')

  // Then (2) คำขอนั้นยังอยู่ในรายการ
  expect(screen.getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
})
