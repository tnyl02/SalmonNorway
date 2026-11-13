// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-14-5 (AC-BRQ-14, FR-BRQ-02 ข้อ ก)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ค่าวัตถุประสงค์อ้างจาก config.js ที่เดียว ไม่เขียนค่าซ้ำใน test (AGENTS.md กฎข้อ 9)
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { PURPOSE_OTHER } from '../config.js'

test('TC-BRQ-14-5 เลือก "อื่น ๆ" แล้วต้องมีช่องพิมพ์ระบุเพิ่ม (AC-BRQ-14)', () => {
  // Given ผู้ใช้อยู่ที่หน้าแบบฟอร์ม
  const client = { createLoanRequest: vi.fn() }
  render(<BorrowRequestForm client={client} />)
  // Then (2) ก่อนเลือกต้องไม่มีช่องนี้
  expect(screen.queryByLabelText('ระบุวัตถุประสงค์')).toBeNull()

  // When เลือกวัตถุประสงค์เป็น "อื่น ๆ"
  fireEvent.change(screen.getByLabelText('วัตถุประสงค์การยืม'), { target: { value: PURPOSE_OTHER } })

  // Then (1) ช่อง "ระบุวัตถุประสงค์" ปรากฏขึ้น
  expect(screen.getByLabelText('ระบุวัตถุประสงค์')).toBeTruthy()
})
