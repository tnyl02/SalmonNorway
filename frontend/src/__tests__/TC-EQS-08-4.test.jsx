// TC-EQS-08-4 (AC-EQS-08 ขอบ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
// "client จำลองที่ไม่ตอบเลย": ใช้ client จริงกับ fetch จำลองที่ไม่ตอบ เพราะการเลิกรอ 6 วินาทีอยู่ใน api/client.js (T-26)
import { act, render, screen } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import { api } from '../api/client.js'
import { EQUIPMENT_REQUEST_TIMEOUT_MS } from '../config.js'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'

afterEach(() => {
  vi.useRealTimers()
  vi.unstubAllGlobals()
})

test('TC-EQS-08-4 หลังบ้านไม่ตอบ ครบ 6 วินาทีแล้วแสดงข้อผิดพลาดพร้อมปุ่มลองใหม่', async () => {
  // Given client ที่ไม่ตอบเลย (fake timers)
  vi.useFakeTimers()
  vi.stubGlobal(
    'fetch',
    vi.fn(
      (url, init = {}) =>
        new Promise((_, reject) => {
          init.signal?.addEventListener('abort', () => reject(new DOMException('aborted', 'AbortError')))
        }),
    ),
  )

  // When เปิดหน้ารายการ แล้วเดินเวลาไป 6 วินาที (EQUIPMENT_REQUEST_TIMEOUT_MS, Q-08)
  render(<EquipmentList client={api} />)
  await act(() => vi.advanceTimersByTimeAsync(EQUIPMENT_REQUEST_TIMEOUT_MS - 1))

  // Then (1) ก่อนครบ 6 วินาทียังไม่แสดงข้อผิดพลาด
  expect(screen.queryByRole('alert')).toBeNull()

  // Then (2) ครบ 6 วินาทีแล้วแสดงข้อความผิดพลาดพร้อมปุ่ม "ลองใหม่" ไม่ค้างที่กำลังโหลด
  await act(() => vi.advanceTimersByTimeAsync(1))
  vi.useRealTimers()
  expect(await screen.findByRole('alert')).toBeTruthy()
  expect(screen.getByRole('button', { name: 'ลองใหม่' })).toBeTruthy()
  expect(screen.queryByText('กำลังโหลด...')).toBeNull()
})
