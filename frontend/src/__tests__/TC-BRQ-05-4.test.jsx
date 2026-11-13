// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-05-4 (AC-BRQ-05, FR-BRQ-05 ข้อ ก, DOM-PRIV-01)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { MSG_STACKED_NOTICE, MSG_SUBMIT_SUCCESS } from '../config.js'

test('TC-BRQ-05-4 หน้าจอของผู้ขอต้องไม่แสดงข้อมูลคำขอของผู้อื่น (AC-BRQ-05)', async () => {
  // Given หลังบ้านตอบสำเร็จ และแนบข้อมูลคิวของผู้อื่นมาด้วย (เผื่อวันหนึ่งสัญญา API เปลี่ยน)
  // หน้าจอต้องไม่หยิบค่าพวกนี้มาแสดง ไม่ว่าจะได้รับมาหรือไม่
  const client = {
    createLoanRequest: vi.fn(async () => ({
      requestId: 91,
      status: 'requested',
      submittedAt: '2026-11-13T09:00:00+07:00',
      queuePosition: 3,
      requesterName: 'ผู้ขอรายอื่น ทดสอบ',
    })),
  }
  render(<BorrowRequestForm client={client} />)

  // When ผู้ขอส่งคำขอสำเร็จ แล้วดูผลลัพธ์บนหน้าจอ
  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))
  const status = await screen.findByRole('status')

  // ต้องเห็นผลของคำขอตัวเองตามปกติ
  expect(status.textContent).toContain(MSG_SUBMIT_SUCCESS(91))
  expect(status.textContent).toContain(MSG_STACKED_NOTICE)

  const page = document.body.textContent
  // Then (1) ไม่มีจำนวนคำขอของผู้อื่น (ไม่มีเลขอื่นนอกจากหมายเลขคำขอของตัวเองในกล่องผลลัพธ์)
  expect(status.textContent).not.toContain('3')
  // Then (2) ไม่มีลำดับคิวของผู้ขอรายนี้
  // (ข้อความ MSG_STACKED_NOTICE มีคำว่า "ตามลำดับการส่งคำขอ" ซึ่งเป็นการอธิบายกติกา ไม่ใช่การบอกลำดับที่)
  expect(page).not.toMatch(/ลำดับที่/)
  // Then (3) ไม่มีชื่อผู้ขออื่น
  expect(page).not.toContain('ผู้ขอรายอื่น ทดสอบ')
})
