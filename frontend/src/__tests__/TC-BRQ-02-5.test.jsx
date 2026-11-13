// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-02-5 (AC-BRQ-02, FR-BRQ-02)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ข้อความไทยอ้างจาก config.js ที่เดียว ไม่เขียนสตริงซ้ำใน test (AGENTS.md กฎข้อ 9)
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { ApiError } from '../api/client.js'
import { FIELD_LABELS, MSG_FIELDS_TO_FIX } from '../config.js'

test('TC-BRQ-02-5 ไม่เลือกวัตถุประสงค์แล้วต้องบอกว่าต้องแก้ช่องไหน (AC-BRQ-02)', async () => {
  // Given หลังบ้านตอบ 422 พร้อมรายชื่อช่องที่ผิดมี purpose (CON-TECH-03)
  const client = {
    createLoanRequest: vi.fn(async () => {
      throw new ApiError(422, 'VALIDATION_FAILED', 'ข้อความจากหลังบ้าน', {
        code: 'VALIDATION_FAILED',
        fields: ['purpose'],
      })
    }),
  }
  render(<BorrowRequestForm client={client} />)

  // When ผู้ใช้กดส่งคำขอบนหน้าจอ
  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  const alert = await screen.findByRole('alert')
  // Then (1) ขึ้นข้อความว่าต้องแก้ช่องใด
  expect(alert.textContent).toContain(MSG_FIELDS_TO_FIX)
  expect(alert.textContent).toContain(FIELD_LABELS.purpose)
  // Then (2) ช่องวัตถุประสงค์มี aria-invalid="true" (NFR-USE-01 ไม่สื่อด้วยสีอย่างเดียว)
  expect(screen.getByLabelText(FIELD_LABELS.purpose).getAttribute('aria-invalid')).toBe('true')
})
