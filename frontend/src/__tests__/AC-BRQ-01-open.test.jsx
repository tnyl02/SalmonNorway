// T-25 ตรวจสถานะอุปกรณ์ตอนเปิดฟอร์ม และล็อกรหัสครุภัณฑ์ที่มาจาก 001 (FR-BRQ-01, NFR-USE-01)
// ครอบส่วน Given ของ AC-BRQ-01 ที่ว่า "นักศึกษาเปิดแบบฟอร์มคำขอของอุปกรณ์ชิ้นที่สถานะว่าง"
// ซึ่งเดิมไม่มีใครตรวจ เพราะ test ของ AC-BRQ-01 ที่มีอยู่จำลองการเปลี่ยนสถานะแล้วยิง POST อย่างเดียว
// ข้อความและคำไทยของสถานะอ้างจาก config.js ที่เดียว (AGENTS.md กฎข้อ 9)
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { vi } from 'vitest'
import BorrowRequestForm from '../pages/BorrowRequest/BorrowRequestForm.jsx'
import { BORROW_REQUEST_MESSAGES } from '../config.js'
import { STATUS_LABEL } from '../components/StatusBadge.jsx'

const ASSET_CODE = 'CS-NB-0003'

// clientWithItem จำลอง client ที่มีทั้ง getItem ของ 001 และ createLoanRequest ของ 002
function clientWithItem(getItem) {
  return {
    getItem,
    createLoanRequest: vi.fn(async () => ({ requestId: 5, status: 'requested', submittedAt: '2026-11-13T09:00:00+07:00' })),
  }
}

test('AC-BRQ-01 เปิดฟอร์มด้วยรหัสจาก 001 แล้วอุปกรณ์ไม่ว่างแล้ว ต้องเตือนทันทีโดยไม่ต้องกดส่ง', async () => {
  // Given อุปกรณ์ชิ้นนั้นถูกเปลี่ยนเป็น "ถูกยืม" ก่อนที่ผู้ขอจะเปิดฟอร์ม
  const client = clientWithItem(vi.fn(async () => ({ assetCode: ASSET_CODE, status: 'borrowed' })))

  // When เปิดฟอร์มด้วยรหัสที่ส่งมาจากปุ่มของ 001
  render(<BorrowRequestForm client={client} assetCode={ASSET_CODE} />)

  // Then เห็นคำเตือนพร้อมคำไทยของสถานะล่าสุด ตั้งแต่ยังไม่กดส่ง
  const warning = await screen.findByText(new RegExp(BORROW_REQUEST_MESSAGES.ITEM_NOT_AVAILABLE))
  expect(warning.textContent).toContain(STATUS_LABEL.borrowed)
  expect(client.getItem).toHaveBeenCalledWith(ASSET_CODE)
  // Then ช่องรหัสครุภัณฑ์ที่มาจาก 001 แก้ไม่ได้
  expect(screen.getByLabelText('รหัสครุภัณฑ์').readOnly).toBe(true)
})

test('AC-BRQ-01 อุปกรณ์ยังว่างอยู่ตอนเปิดฟอร์ม ต้องไม่มีคำเตือน', async () => {
  // Given อุปกรณ์ชิ้นนั้นยังว่าง
  const client = clientWithItem(vi.fn(async () => ({ assetCode: ASSET_CODE, status: 'available' })))

  // When เปิดฟอร์มด้วยรหัสที่ส่งมาจาก 001
  render(<BorrowRequestForm client={client} assetCode={ASSET_CODE} />)
  await waitFor(() => expect(client.getItem).toHaveBeenCalled())

  // Then ไม่มีคำเตือนใด ๆ
  expect(screen.queryByText(new RegExp(BORROW_REQUEST_MESSAGES.ITEM_NOT_AVAILABLE))).toBeNull()
})

test('AC-BRQ-01 อ่านสถานะตอนเปิดฟอร์มไม่สำเร็จ ต้องไม่ปิดกั้นการส่งคำขอ (FR-BRQ-09)', async () => {
  // Given การอ่านสถานะล้มเหลว เช่นระบบของ 001 ขัดข้องชั่วคราว
  const client = clientWithItem(vi.fn(async () => { throw new Error('อ่านสถานะไม่สำเร็จ') }))

  // When เปิดฟอร์มแล้วกดส่งคำขอ
  render(<BorrowRequestForm client={client} assetCode={ASSET_CODE} />)
  await waitFor(() => expect(client.getItem).toHaveBeenCalled())
  fireEvent.click(screen.getByRole('button', { name: 'ส่งคำขอยืม' }))

  // Then ไม่มีคำเตือน และคำขอยังถูกส่งไปหลังบ้านตามปกติ เพราะการตรวจจริงอยู่ที่หลังบ้าน
  expect(screen.queryByText(new RegExp(BORROW_REQUEST_MESSAGES.ITEM_NOT_AVAILABLE))).toBeNull()
  await waitFor(() => expect(client.createLoanRequest).toHaveBeenCalledTimes(1))
})

test('AC-BRQ-01 เปิดหน้าตรงจากเมนู (ไม่มีรหัสจาก 001) ต้องไม่เรียก API และยังพิมพ์รหัสเองได้', () => {
  // Given ผู้ใช้เปิดหน้าจากเมนู ไม่ได้มาจากปุ่มของ 001
  const client = clientWithItem(vi.fn())

  // When หน้าจอถูกเปิดขึ้นมา
  render(<BorrowRequestForm client={client} />)

  // Then ไม่ยิง API ของ 001 เลย และช่องรหัสยังแก้ได้
  expect(client.getItem).not.toHaveBeenCalled()
  expect(screen.getByLabelText('รหัสครุภัณฑ์').readOnly).toBe(false)
})
