// T-19 แบบฟอร์มส่งคำขอยืม (FR-BRQ-01, 02, 03, 10, NFR-USE-01)
// ครอบส่วนที่ผู้ใช้เห็นบนหน้าจอของ AC-BRQ-02, AC-BRQ-03, AC-BRQ-10, AC-BRQ-14
// ใช้ client จำลองแทนหลังบ้าน การตรวจกฎจริงอยู่ที่หลังบ้านแล้ว (TestAC_BRQ_02, 03, 10, 14)
// ข้อความบนหน้าจออ้างจาก config.js ที่เดียว ไม่เขียนสตริงซ้ำใน test
import { fireEvent, render, screen } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { ApiError } from '../api/client.js'
import { MSG_ATTACHMENT_HINT, MSG_LOAN_NORMAL_CONDITIONS, MSG_SUBMIT_SUCCESS, PURPOSE_OPTIONS } from '../config.js'

// okClient คืนผลสำเร็จแบบเดียวกับหลังบ้าน (FR-BRQ-04) และจำค่าที่ถูกส่งไว้ตรวจ
function okClient() {
  return { createLoanRequest: vi.fn(async () => ({ requestId: 42, status: 'requested', submittedAt: '2026-11-13T09:00:00+07:00' })) }
}

// fillRequiredFields กรอก 5 ช่องที่ FR-BRQ-03 ข้อ ก บังคับ บวกรหัสครุภัณฑ์
function fillRequiredFields() {
  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0003' } })
  fireEvent.change(screen.getByLabelText('ชื่อ-นามสกุล'), { target: { value: 'ธัญลักษณ์ ทดสอบ' } })
  fireEvent.change(screen.getByLabelText('รหัสนักศึกษา'), { target: { value: '660710715' } })
  fireEvent.change(screen.getByLabelText('วันและเวลาที่เริ่มยืม'), { target: { value: '2026-11-13T09:00' } })
  fireEvent.change(screen.getByLabelText('วันและเวลาที่คืน'), { target: { value: '2026-11-13T15:00' } })
  fireEvent.change(screen.getByLabelText('วัตถุประสงค์การยืม'), { target: { value: 'study' } })
}

test('ฟอร์มมีช่องที่ FR-BRQ-03 บังคับครบ และวัตถุประสงค์เป็นรายการ 3 ค่า (AC-BRQ-02, AC-BRQ-14)', () => {
  render(<BorrowRequestForm client={okClient()} />)

  for (const label of ['ชื่อ-นามสกุล', 'รหัสนักศึกษา', 'วันและเวลาที่เริ่มยืม', 'วันและเวลาที่คืน', 'วัตถุประสงค์การยืม']) {
    expect(screen.getByLabelText(label)).toBeTruthy()
  }
  for (const option of PURPOSE_OPTIONS) {
    expect(screen.getByRole('option', { name: option.label })).toBeTruthy()
  }
})

test('เลือก "อื่น ๆ" แล้วต้องมีช่องพิมพ์ระบุเพิ่ม (AC-BRQ-14, FR-BRQ-02 ข้อ ก)', () => {
  render(<BorrowRequestForm client={okClient()} />)
  expect(screen.queryByLabelText('ระบุวัตถุประสงค์')).toBeNull()

  fireEvent.change(screen.getByLabelText('วัตถุประสงค์การยืม'), { target: { value: 'other' } })

  expect(screen.getByLabelText('ระบุวัตถุประสงค์')).toBeTruthy()
})

test('ยืมปกติเห็นเงื่อนไขก่อนส่ง ส่วนค้างคืนมีช่องแนบใบบันทึกข้อความ (AC-BRQ-10, FR-BRQ-10)', () => {
  render(<BorrowRequestForm client={okClient()} />)

  // ปกติ: แสดงเงื่อนไขคืนภายใน 16:00 น. ของวันเดียวกัน และใช้ในพื้นที่มหาวิทยาลัย (FR-BRQ-10 ข้อ ข)
  expect(screen.getByText(MSG_LOAN_NORMAL_CONDITIONS)).toBeTruthy()
  expect(screen.queryByLabelText('ใบบันทึกข้อความ')).toBeNull()

  fireEvent.click(screen.getByLabelText('ค้างคืน'))

  // ค้างคืน: มีช่องแนบไฟล์พร้อมบอกชนิดและขนาดที่รับได้ (FR-BRQ-10 ข้อ ก)
  expect(screen.getByLabelText('ใบบันทึกข้อความ')).toBeTruthy()
  expect(screen.getByText(MSG_ATTACHMENT_HINT)).toBeTruthy()
  expect(screen.queryByText(MSG_LOAN_NORMAL_CONDITIONS)).toBeNull()
})

test('กรอกครบแล้วกดส่ง ต้องส่งค่าที่กรอกไปหลังบ้านและแจ้งผลสำเร็จ (AC-BRQ-03, FR-BRQ-04)', async () => {
  const client = okClient()
  render(<BorrowRequestForm client={client} />)
  fillRequiredFields()

  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  expect(await screen.findByRole('status')).toHaveProperty('textContent', MSG_SUBMIT_SUCCESS(42))
  const [input, attachment] = client.createLoanRequest.mock.calls[0]
  expect(input.assetCode).toBe('CS-NB-0003')
  expect(input.loanType).toBe('normal')
  expect(input.purpose).toBe('study')
  expect(input.purposeOther).toBe('') // ไม่ได้เลือก "อื่น ๆ" จึงไม่ส่งข้อความที่พิมพ์ไว้ (FR-BRQ-02 ข้อ ข)
  // ส่งเวลาพร้อมเขตเวลาไทยเสมอ ไม่ปล่อยให้หลังบ้านเดาจากเครื่องที่รัน (glossary, DOM-LOAN-01)
  expect(input.startAt).toBe('2026-11-13T09:00:00+07:00')
  expect(input.endAt).toBe('2026-11-13T15:00:00+07:00')
  expect(attachment).toBeNull()
})

test('ยืมค้างคืนส่งไฟล์แนบไปด้วย (FR-BRQ-10 ข้อ ก)', async () => {
  const client = okClient()
  render(<BorrowRequestForm client={client} />)
  fillRequiredFields()
  fireEvent.click(screen.getByLabelText('ค้างคืน'))
  const file = new File(['%PDF-1.4'], 'memo.pdf', { type: 'application/pdf' })
  fireEvent.change(screen.getByLabelText('ใบบันทึกข้อความ'), { target: { files: [file] } })

  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  await screen.findByRole('status')
  const [input, attachment] = client.createLoanRequest.mock.calls[0]
  expect(input.loanType).toBe('overnight')
  expect(attachment).toBe(file)
})

test('หลังบ้านปฏิเสธ ต้องไม่ขึ้นข้อความสำเร็จ และแสดงเหตุผลที่หลังบ้านส่งมา (FR-BRQ-09)', async () => {
  // การจับคู่รหัส error กับข้อความบนหน้าจอเป็นงานของ T-20 ที่นี่แค่ต้องไม่กลืนข้อผิดพลาด
  const client = { createLoanRequest: vi.fn(async () => { throw new ApiError(409, 'DUPLICATE_REQUEST', 'คุณมีคำขอของอุปกรณ์ชิ้นนี้ค้างอยู่แล้ว') }) }
  render(<BorrowRequestForm client={client} />)
  fillRequiredFields()

  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  expect(await screen.findByRole('alert')).toHaveProperty('textContent', 'คุณมีคำขอของอุปกรณ์ชิ้นนี้ค้างอยู่แล้ว')
  expect(screen.queryByRole('status')).toBeNull()
})
