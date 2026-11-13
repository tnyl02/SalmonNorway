// test นี้มาจากตาราง test-cases.md แถว TC-BRQ-04-4 (AC-BRQ-04, FR-BRQ-04)
// ทีมตรวจแถวแล้วเมื่อ 2569-11-13 สถานะ "ใช้ได้"
// **ห้ามแก้ test นี้ให้ผ่าน** ถ้าไม่ผ่านต้องหาว่าโค้ดผิดหรือ AC ผิด (AGENTS.md)
// ข้อความไทยอ้างจาก config.js ที่เดียว ไม่เขียนสตริงซ้ำใน test (AGENTS.md กฎข้อ 9)
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { MSG_SUBMIT_SUCCESS } from '../config.js'

const REQUEST_ID = 77

test('TC-BRQ-04-4 ส่งสำเร็จแล้วต้องเห็นหมายเลขคำขอและสถานะรออนุมัติ (AC-BRQ-04)', async () => {
  // Given หลังบ้านตอบ 201 พร้อมหมายเลขคำขอ
  const client = {
    createLoanRequest: vi.fn(async () => ({
      requestId: REQUEST_ID,
      status: 'requested',
      submittedAt: '2026-11-13T09:00:00+07:00',
    })),
  }
  render(<BorrowRequestForm client={client} />)
  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0003' } })
  fireEvent.change(screen.getByLabelText('ชื่อ-นามสกุล'), { target: { value: 'ธัญลักษณ์ ทดสอบ' } })
  fireEvent.change(screen.getByLabelText('รหัสนักศึกษา'), { target: { value: '660710715' } })
  fireEvent.change(screen.getByLabelText('วันและเวลาที่เริ่มยืม'), { target: { value: '2026-11-13T09:00' } })
  fireEvent.change(screen.getByLabelText('วันและเวลาที่คืน'), { target: { value: '2026-11-13T15:00' } })
  fireEvent.change(screen.getByLabelText('วัตถุประสงค์การยืม'), { target: { value: 'study' } })

  // When ผู้ใช้กดส่งคำขอบนหน้าจอ
  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  const status = await screen.findByRole('status')
  // Then (1) เห็นหมายเลขคำขอที่หลังบ้านส่งมา
  expect(status.textContent).toContain(String(REQUEST_ID))
  // Then (2) เห็นคำว่าสถานะ "รออนุมัติ" ซึ่งอยู่ในข้อความผลสำเร็จของ config.js
  expect(status.textContent).toContain(MSG_SUBMIT_SUCCESS(REQUEST_ID))
})
