// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-01-4 (AC-BRQ-01, FR-BRQ-01)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ข้อความและคำไทยของสถานะอ้างจาก config.js ที่เดียว (AGENTS.md กฎข้อ 9)
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { ApiError } from '../api/client.js'
import { BORROW_REQUEST_MESSAGES, STATUS_LABELS } from '../config.js'

test('TC-BRQ-01-4 อุปกรณ์ถูกยืมระหว่างกรอก ต้องแจ้งพร้อมสถานะล่าสุด (AC-BRQ-01)', async () => {
  // Given หลังบ้านตอบ 409 ITEM_NOT_AVAILABLE พร้อมสถานะล่าสุด borrowed (CON-TECH-03)
  const client = {
    createLoanRequest: vi.fn(async () => {
      throw new ApiError(409, 'ITEM_NOT_AVAILABLE', 'ข้อความจากหลังบ้าน', {
        code: 'ITEM_NOT_AVAILABLE',
        status: 'borrowed',
      })
    }),
  }
  render(<BorrowRequestForm client={client} />)

  // When ผู้ใช้กดส่งคำขอบนหน้าจอ
  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  const alert = await screen.findByRole('alert')
  // Then (1) ไม่ขึ้นข้อความสำเร็จ
  expect(screen.queryByRole('status')).toBeNull()
  // Then (2) ขึ้นข้อความแจ้งพร้อมคำไทยของสถานะล่าสุด
  expect(alert.textContent).toContain(BORROW_REQUEST_MESSAGES.ITEM_NOT_AVAILABLE)
  expect(alert.textContent).toContain(STATUS_LABELS.borrowed)
})
