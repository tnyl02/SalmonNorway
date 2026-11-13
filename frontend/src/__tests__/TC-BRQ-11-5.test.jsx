// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-11-5 (AC-BRQ-11, FR-BRQ-05 ข)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ข้อความไทยอ้างจาก config.js ที่เดียว ไม่เขียนสตริงซ้ำใน test (AGENTS.md กฎข้อ 9)
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { ApiError } from '../api/client.js'
import { BORROW_REQUEST_MESSAGES } from '../config.js'

test('TC-BRQ-11-5 ส่งซ้ำแล้วต้องเห็นข้อความว่ามีคำขอค้างอยู่ และไม่มีข้อความสำเร็จ (AC-BRQ-11)', async () => {
  // Given หลังบ้านตอบ 409 DUPLICATE_REQUEST (CON-TECH-03)
  const client = {
    createLoanRequest: vi.fn(async () => {
      throw new ApiError(409, 'DUPLICATE_REQUEST', 'ข้อความจากหลังบ้าน', { code: 'DUPLICATE_REQUEST' })
    }),
  }
  render(<BorrowRequestForm client={client} />)

  // When ผู้ใช้กดส่งคำขอบนหน้าจอ
  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  const alert = await screen.findByRole('alert')
  // Then (1) ขึ้นข้อความว่ามีคำขอของอุปกรณ์ชิ้นนี้ค้างอยู่แล้ว
  expect(alert.textContent).toContain(BORROW_REQUEST_MESSAGES.DUPLICATE_REQUEST)
  // Then (2) ไม่ขึ้นข้อความสำเร็จ
  expect(screen.queryByRole('status')).toBeNull()
})
