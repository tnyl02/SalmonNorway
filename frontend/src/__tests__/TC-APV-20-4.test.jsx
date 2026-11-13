// TC-APV-20-4 (AC-APV-20 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'

test('TC-APV-20-4 โหลดรายการแล้วได้ UPSTREAM_UNAVAILABLE แสดงข้อความและไม่มีตาราง', async () => {
  // Given client จำลองตอบ error รหัส UPSTREAM_UNAVAILABLE ตอนโหลดรายการ
  const client = {
    getApprovalRequests: async () => Promise.reject(Object.assign(new Error('503'), { code: 'UPSTREAM_UNAVAILABLE' })),
  }

  // When เปิดหน้ารายการคำขอ
  render(<ApproveRequests client={client} />)

  // Then (1) แสดง "ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง"
  expect((await screen.findByRole('alert')).textContent).toBe('ระบบขัดข้องชั่วคราว กรุณาลองใหม่อีกครั้ง')

  // Then (2) ไม่มีตาราง
  expect(screen.queryByRole('table')).toBeNull()
})
