// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-03-5 (AC-BRQ-03, FR-BRQ-03)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ชื่อช่องภาษาไทยอ้างจาก config.js ที่เดียว ไม่เขียนสตริงซ้ำใน test (AGENTS.md กฎข้อ 9)
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { ApiError } from '../api/client.js'
import { FIELD_LABELS } from '../config.js'

test('TC-BRQ-03-5 หลังบ้านบอกช่องที่ผิด หน้าจอต้องชี้ให้ครบและไม่ชี้ช่องที่ถูก (AC-BRQ-03)', async () => {
  // Given หลังบ้านตอบ 422 พร้อมรายชื่อช่องที่ผิดคือ studentCode และ endAt
  const client = {
    createLoanRequest: vi.fn(async () => {
      throw new ApiError(422, 'VALIDATION_FAILED', 'ข้อความจากหลังบ้าน', {
        code: 'VALIDATION_FAILED',
        fields: ['studentCode', 'endAt'],
      })
    }),
  }
  render(<BorrowRequestForm client={client} />)

  // When ผู้ใช้กดส่งคำขอบนหน้าจอ
  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  const alert = await screen.findByRole('alert')
  // Then (1) เห็นชื่อช่องทั้งสองเป็นภาษาไทย
  expect(alert.textContent).toContain(FIELD_LABELS.studentCode)
  expect(alert.textContent).toContain(FIELD_LABELS.endAt)
  // Then (2) ทั้งสองช่องมี aria-invalid="true"
  expect(screen.getByLabelText(FIELD_LABELS.studentCode).getAttribute('aria-invalid')).toBe('true')
  expect(screen.getByLabelText(FIELD_LABELS.endAt).getAttribute('aria-invalid')).toBe('true')
  // Then (3) ช่องที่ไม่ผิดต้องไม่มี aria-invalid
  expect(screen.getByLabelText(FIELD_LABELS.requesterName).getAttribute('aria-invalid')).toBeNull()
  expect(screen.getByLabelText(FIELD_LABELS.startAt).getAttribute('aria-invalid')).toBeNull()
})
