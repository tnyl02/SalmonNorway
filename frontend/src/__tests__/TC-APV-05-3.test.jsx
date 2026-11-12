// TC-APV-05-3 (AC-APV-05 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'

test('TC-APV-05-3 นักศึกษาเปิดหน้าอนุมัติ แสดงข้อความ FORBIDDEN_ROLE ไม่มีตารางและปุ่ม', async () => {
  // Given client จำลองตอบ error รหัส FORBIDDEN_ROLE ตอนโหลดรายการ (นักศึกษาเปิดหน้าอนุมัติ)
  const client = {
    getApprovalRequests: async () => Promise.reject(Object.assign(new Error('403'), { code: 'FORBIDDEN_ROLE' })),
  }

  // When เปิดหน้ารายการคำขอ
  render(<ApproveRequests client={client} />)

  // Then (1) แสดง "บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้"
  expect((await screen.findByRole('alert')).textContent).toBe('บทบาทของคุณไม่มีสิทธิ์ทำรายการนี้')

  // Then (2) ไม่มีตาราง และไม่มีปุ่ม "อนุมัติ" หรือ "ปฏิเสธ"
  expect(screen.queryByRole('table')).toBeNull()
  expect(screen.queryByRole('button', { name: 'อนุมัติ' })).toBeNull()
  expect(screen.queryByRole('button', { name: 'ปฏิเสธ' })).toBeNull()
})
