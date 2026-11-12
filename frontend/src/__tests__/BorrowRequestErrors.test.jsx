// T-20 ข้อความแจ้งผลและข้อผิดพลาดบนหน้าฟอร์มคำขอยืม (FR-BRQ-01, 03, 06, 09, 10, CON-TECH-03)
// ครอบส่วนข้อความที่ผู้ใช้เห็นของ AC-BRQ-01, AC-BRQ-03, AC-BRQ-06, AC-BRQ-09, AC-BRQ-13
// หน้าจอเลือกข้อความจาก "รหัส" ไม่ใช่จากข้อความที่หลังบ้านส่งมา (CON-TECH-03)
// กฎจริงทั้งหมดถูกตรวจที่หลังบ้านแล้ว ที่นี่ตรวจเฉพาะว่าผู้ใช้เห็นอะไร
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { ApiError } from '../api/client.js'
import { BORROW_REQUEST_MESSAGES, FIELD_LABELS, MSG_FIELDS_TO_FIX } from '../config.js'

// submitWithError ให้หลังบ้านจำลองตอบ error ตามรหัสที่ระบุ แล้วกดส่งฟอร์ม
async function submitWithError(code, body = {}) {
  const client = {
    createLoanRequest: vi.fn(async () => {
      throw new ApiError(400, code, 'ข้อความจากหลังบ้าน', { code, ...body })
    }),
  }
  render(<BorrowRequestForm client={client} />)
  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))
  return screen.findByRole('alert')
}

test('ทุกรหัส error ที่ฟีเจอร์นี้ใช้ มีข้อความของตัวเองบนหน้าจอ (CON-TECH-03)', async () => {
  for (const [code, message] of Object.entries(BORROW_REQUEST_MESSAGES)) {
    const alert = await submitWithError(code)
    expect(alert.textContent).toContain(message)
    screen.getByRole('button', { name: 'ส่งคำขอยืม' }) // ฟอร์มยังอยู่ ผู้ขอแก้แล้วส่งใหม่ได้ (FR-BRQ-09)
    cleanupRender()
  }
})

test('ข้อมูลไม่ครบ ต้องบอกว่าช่องไหนที่ต้องแก้ (AC-BRQ-03, FR-BRQ-03)', async () => {
  const alert = await submitWithError('VALIDATION_FAILED', { fields: ['studentCode', 'endAt'] })

  expect(alert.textContent).toContain(BORROW_REQUEST_MESSAGES.VALIDATION_FAILED)
  expect(alert.textContent).toContain(MSG_FIELDS_TO_FIX)
  expect(alert.textContent).toContain(FIELD_LABELS.studentCode)
  expect(alert.textContent).toContain(FIELD_LABELS.endAt)
  // ช่องที่ผิดถูกชี้ไว้ที่ตัวช่องเองด้วย ไม่ใช่แค่ในข้อความรวม (NFR-USE-01)
  expect(screen.getByLabelText(FIELD_LABELS.studentCode).getAttribute('aria-invalid')).toBe('true')
  expect(screen.getByLabelText(FIELD_LABELS.endAt).getAttribute('aria-invalid')).toBe('true')
  expect(screen.getByLabelText(FIELD_LABELS.requesterName).getAttribute('aria-invalid')).toBeNull()
})

test('เวลาคืนเกิน 16:00 น. ของวันเดียวกัน ต้องชี้ช่องวันและเวลาที่คืน (AC-BRQ-13, FR-BRQ-03 ข้อ ค)', async () => {
  const alert = await submitWithError('VALIDATION_FAILED', { fields: ['endAt'] })

  expect(alert.textContent).toContain(FIELD_LABELS.endAt)
  expect(screen.getByLabelText(FIELD_LABELS.endAt).getAttribute('aria-invalid')).toBe('true')
})

test('อุปกรณ์ถูกยืมไประหว่างกรอก ต้องแจ้งพร้อมสถานะล่าสุด (AC-BRQ-01, FR-BRQ-01)', async () => {
  const alert = await submitWithError('ITEM_NOT_AVAILABLE', { status: 'borrowed' })

  expect(alert.textContent).toContain(BORROW_REQUEST_MESSAGES.ITEM_NOT_AVAILABLE)
  expect(alert.textContent).toContain('ถูกยืม') // คำไทยชุดเดียวกับป้ายสถานะของ 001 (STATUS_LABEL)
})

test('มีของคืนเกินกำหนดค้างอยู่ ต้องแจ้งเหตุผล (AC-BRQ-06, FR-BRQ-06)', async () => {
  const alert = await submitWithError('HAS_OVERDUE_LOAN')

  expect(alert.textContent).toContain(BORROW_REQUEST_MESSAGES.HAS_OVERDUE_LOAN)
})

test('ระบบขัดข้องตอนบันทึก ต้องแจ้งให้ลองใหม่ (AC-BRQ-09, FR-BRQ-09)', async () => {
  const alert = await submitWithError('UPSTREAM_UNAVAILABLE')

  expect(alert.textContent).toContain(BORROW_REQUEST_MESSAGES.UPSTREAM_UNAVAILABLE)
})

test('รหัสที่หน้าจอยังไม่รู้จัก ต้องไม่กลืนข้อผิดพลาด ใช้ข้อความจากหลังบ้านแทน (CON-TECH-03)', async () => {
  const alert = await submitWithError('SOME_NEW_CODE')

  expect(alert.textContent).toContain('ข้อความจากหลังบ้าน')
})

// cleanupRender ล้างหน้าจอระหว่างวนรหัสใน test แรก ไม่ให้ alert ของรอบก่อนค้างอยู่
function cleanupRender() {
  document.body.innerHTML = ''
}
