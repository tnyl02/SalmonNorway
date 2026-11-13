// TC-APV-20-5 (AC-APV-20 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-20-5 อนุมัติแล้วได้ UPSTREAM_UNAVAILABLE แสดงข้อความ และทั้ง 2 แถวยังอยู่', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ และตอบ error รหัส UPSTREAM_UNAVAILABLE เมื่ออนุมัติ
  const client = {
    ...fakeClient(),
    approveRequest: async () => Promise.reject(Object.assign(new Error('503'), { code: 'UPSTREAM_UNAVAILABLE' })),
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "อนุมัติ" ที่คำขอแรก
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) แสดง "ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง"
  expect((await screen.findByRole('alert')).textContent).toBe('ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง')

  // Then (2) คำขอนั้นยังอยู่ในรายการ ทั้ง 2 แถวยังอยู่
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
  expect(screen.getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
})
