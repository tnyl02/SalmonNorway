// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-09-4 (AC-BRQ-09, FR-BRQ-09)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ข้อความไทยอ้างจาก config.js ที่เดียว (AGENTS.md กฎข้อ 9)
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { ApiError } from '../api/client.js'
import { BORROW_REQUEST_MESSAGES } from '../config.js'

test('TC-BRQ-09-4 ระบบขัดข้องตอนบันทึก ต้องแจ้งให้ลองใหม่และฟอร์มยังอยู่ (AC-BRQ-09)', async () => {
  // Given หลังบ้านตอบ 503 UPSTREAM_UNAVAILABLE (CON-TECH-03)
  const client = {
    createLoanRequest: vi.fn(async () => {
      throw new ApiError(503, 'UPSTREAM_UNAVAILABLE', 'ข้อความจากหลังบ้าน', { code: 'UPSTREAM_UNAVAILABLE' })
    }),
  }
  render(<BorrowRequestForm client={client} />)

  // When ผู้ใช้กดส่งคำขอบนหน้าจอ
  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  const alert = await screen.findByRole('alert')
  // Then (1) ขึ้นข้อความให้ลองใหม่
  expect(alert.textContent).toContain(BORROW_REQUEST_MESSAGES.UPSTREAM_UNAVAILABLE)
  // Then (2) ฟอร์มยังอยู่ ผู้ขอกดส่งซ้ำได้
  const submit = screen.getByRole('button', { name: 'ส่งคำขอยืม' })
  expect(submit.disabled).toBe(false)
  fireEvent.click(submit)
  expect(client.createLoanRequest).toHaveBeenCalledTimes(2)
})
