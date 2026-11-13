// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-13-5 (AC-BRQ-13, FR-BRQ-03 ข้อ ค)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ข้อความเงื่อนไขอ้างจาก config.js ที่เดียว และเลข 16 มาจาก LOAN_NORMAL_RETURN_HOUR
import { render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { LOAN_NORMAL_RETURN_HOUR, MSG_LOAN_NORMAL_CONDITIONS } from '../config.js'

test('TC-BRQ-13-5 เลือกยืมปกติแล้วต้องเห็นเงื่อนไขเวลาคืนก่อนกดส่ง (AC-BRQ-13)', () => {
  // Given ผู้ใช้เลือกประเภทการยืมเป็น "ปกติ" ซึ่งเป็นค่าตั้งต้นของฟอร์ม
  const client = { createLoanRequest: vi.fn() }
  render(<BorrowRequestForm client={client} />)
  expect(screen.getByLabelText('ปกติ').checked).toBe(true)

  // When อยู่ในแบบฟอร์มก่อนกดส่ง
  const conditions = screen.getByText(MSG_LOAN_NORMAL_CONDITIONS)

  // Then (1) เห็นข้อความเงื่อนไขที่บอกเวลา 16:00 น. ของวันเดียวกัน
  expect(conditions).toBeTruthy()
  expect(conditions.textContent).toContain(String(LOAN_NORMAL_RETURN_HOUR))
})
