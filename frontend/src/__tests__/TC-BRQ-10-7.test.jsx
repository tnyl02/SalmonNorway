// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-10-7 (AC-BRQ-10, FR-BRQ-10 ข้อ ก)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ข้อความและค่ากำกับอ้างจาก config.js ที่เดียว (AGENTS.md กฎข้อ 9)
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { ATTACHMENT_MAX_BYTES, MSG_ATTACHMENT_HINT } from '../config.js'

test('TC-BRQ-10-7 เลือกค้างคืนต้องมีช่องแนบไฟล์พร้อมบอกชนิดและขนาดที่รับได้ (AC-BRQ-10)', () => {
  // Given ผู้ใช้อยู่ที่หน้าแบบฟอร์ม
  const client = { createLoanRequest: vi.fn() }
  render(<BorrowRequestForm client={client} />)

  // When เลือกประเภทการยืมเป็น "ค้างคืน"
  fireEvent.click(screen.getByLabelText('ค้างคืน'))

  // Then (1) มีช่องแนบใบบันทึกข้อความ
  expect(screen.getByLabelText('ใบบันทึกข้อความ')).toBeTruthy()
  // Then (2) เห็นข้อความบอกชนิดและขนาดไฟล์ที่รับได้
  const hint = screen.getByText(MSG_ATTACHMENT_HINT)
  expect(hint.textContent).toContain('PDF')
  expect(hint.textContent).toContain(String(ATTACHMENT_MAX_BYTES / (1024 * 1024)))
})
