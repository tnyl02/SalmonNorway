// TC-APV-05-8 (AC-APV-05 ผิด) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-05-8 ยืนยันตัวตนหมดอายุตอนกดอนุมัติ แสดงข้อความสำรองและคำขอยังอยู่ในรายการ', async () => {
  // Given หน้ารายการโหลดแล้ว client จำลองตอบ error รหัส UNAUTHENTICATED เมื่อกด "อนุมัติ"
  let called = 0
  const client = {
    ...fakeClient(),
    approveRequest: async () => {
      called++
      return Promise.reject(Object.assign(new Error('401'), { code: 'UNAUTHENTICATED' }))
    },
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "อนุมัติ" ที่คำขอแรก
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))
  await waitFor(() => expect(called).toBe(1))

  // Then (1) แสดง "ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง" (Q-APV-32)
  expect((await screen.findByRole('alert')).textContent).toBe('ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง')

  // Then (2) คำขอยังอยู่ในรายการ ทั้ง 2 แถวยังอยู่
  expect(screen.getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
})
