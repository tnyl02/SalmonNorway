// TC-APV-11-5 (AC-APV-11 ขอบ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-11-5 ไม่มีคำขอเลยแสดง "ไม่มีคำขอที่รออนุมัติ" และไม่มีตาราง', async () => {
  // Given ไม่มีคำขอ requested ของภาควิชาตนเองเลย (client จำลองคืน items: [])
  // When  เปิดหน้ารายการคำขอ
  render(<ApproveRequests client={fakeClient([])} />)

  // Then (1) แสดง "ไม่มีคำขอที่รออนุมัติ" (Q-APV-31)
  expect(await screen.findByText('ไม่มีคำขอที่รออนุมัติ')).toBeTruthy()

  // Then (2) ไม่มีตารางรายการคำขอ
  expect(screen.queryByRole('table')).toBeNull()
})
