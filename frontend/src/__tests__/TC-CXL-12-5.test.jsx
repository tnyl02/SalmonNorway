// TC-CXL-12-5 (AC-CXL-12 ผิด) จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { vi } from 'vitest'
import { ApiError } from '../api/client.js'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'

test('TC-CXL-12-5 โหลดรายการไม่สำเร็จ แสดงข้อความและปุ่มลองใหม่ ไม่แสดงว่าไม่มีคำขอ', async () => {
  // Given client จำลองตอบ 503 ตอนโหลดรายการ
  const client = {
    listMyPendingRequests: vi.fn(async () => {
      throw new ApiError(503, 'UPSTREAM_UNAVAILABLE')
    }),
    cancelRequest: vi.fn(),
  }

  // When เปิดหน้า "คำขอที่รออนุมัติของฉัน"
  render(<MyPendingRequests client={client} />)

  // Then (2) แสดง "โหลดรายการไม่สำเร็จ กรุณาลองใหม่อีกครั้ง"
  expect(await screen.findByText('โหลดรายการไม่สำเร็จ กรุณาลองใหม่อีกครั้ง')).toBeTruthy()

  // Then (1) ไม่แสดงข้อความ "ไม่มีคำขอที่รออนุมัติ"
  expect(screen.queryByText('ไม่มีคำขอที่รออนุมัติ')).toBeNull()

  // Then (3) มีปุ่ม "ลองใหม่" กดแล้ว client ถูกเรียกโหลดรายการอีกครั้ง
  const calls = client.listMyPendingRequests.mock.calls.length
  fireEvent.click(screen.getByRole('button', { name: 'ลองใหม่' }))
  await waitFor(() => expect(client.listMyPendingRequests.mock.calls.length).toBe(calls + 1))
})
