// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-10-6 (AC-BRQ-10, FR-BRQ-10 ข้อ ข)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ข้อความอ้างจาก config.js ที่เดียว (AGENTS.md กฎข้อ 9)
import { render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { LOAN_NORMAL_RETURN_HOUR, MSG_LOAN_NORMAL_CONDITIONS } from '../config.js'

test('TC-BRQ-10-6 ยืมปกติต้องเห็นเงื่อนไขครบ และไม่มีช่องแนบไฟล์ (AC-BRQ-10)', () => {
  // Given ผู้ใช้เลือกประเภทการยืมเป็น "ปกติ" ซึ่งเป็นค่าตั้งต้นของฟอร์ม
  const client = { createLoanRequest: vi.fn() }
  render(<BorrowRequestForm client={client} />)

  // When อยู่ในแบบฟอร์มก่อนกดส่ง
  const conditions = screen.getByText(MSG_LOAN_NORMAL_CONDITIONS)

  // Then (1) เห็นเงื่อนไขว่าต้องคืนภายใน 16:00 น. ของวันเดียวกัน
  expect(conditions.textContent).toContain(String(LOAN_NORMAL_RETURN_HOUR))
  expect(conditions.textContent).toContain('วันเดียวกัน')
  // Then (2) เห็นเงื่อนไขว่าใช้งานภายในพื้นที่มหาวิทยาลัย
  expect(conditions.textContent).toContain('พื้นที่มหาวิทยาลัย')
  // Then (3) ไม่มีช่องแนบไฟล์
  expect(screen.queryByLabelText('ใบบันทึกข้อความ')).toBeNull()
})
