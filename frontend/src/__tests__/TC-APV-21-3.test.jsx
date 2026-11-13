// TC-APV-21-3 (AC-APV-21 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-21-3 อนุมัติแล้วได้ REQUEST_NOT_FOUND แสดง "ไม่พบคำขอนี้"', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ และตอบ error รหัส REQUEST_NOT_FOUND เมื่ออนุมัติ
  const client = {
    ...fakeClient(),
    approveRequest: async () => Promise.reject(Object.assign(new Error('404'), { code: 'REQUEST_NOT_FOUND' })),
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "อนุมัติ" ที่คำขอแรก
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) แสดง "ไม่พบคำขอนี้"
  expect((await screen.findByRole('alert')).textContent).toBe('ไม่พบคำขอนี้')
})
