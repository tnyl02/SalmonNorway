// AC-EQS-08 (FR-EQS-08) ใช้ client จำลองแทนหลังบ้าน
//   Given ฐานข้อมูลไม่ตอบสนอง
//   When  เปิดหน้ารายการอุปกรณ์
//   Then  แสดงข้อความผิดพลาดพร้อมปุ่มลองใหม่ ไม่แสดงหน้าจอค้างและไม่แสดงข้อมูลเก่า
// กรณีต่อได้แต่ช้า: หน้าจอรอไม่เกิน EQUIPMENT_REQUEST_TIMEOUT_MS (Q-08, T-26) ใช้ client จริงกับ fetch จำลองที่ไม่ตอบ
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import { ApiError, api } from '../api/client.js'
import { DEFAULT_ERROR_MESSAGE } from '../components/ErrorState.jsx'
import { EQUIPMENT_REQUEST_TIMEOUT_MS } from '../config.js'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'

afterEach(() => {
  vi.useRealTimers()
  vi.unstubAllGlobals()
})

// fetch จำลองที่ไม่ตอบเลย จนกว่า request ถูกยกเลิกด้วย signal (หลังบ้านค้าง)
function hangingFetch() {
  return vi.fn(
    (url, init = {}) =>
      new Promise((_, reject) => {
        init.signal?.addEventListener('abort', () => reject(new DOMException('aborted', 'AbortError')))
      }),
  )
}

const DOWN_MESSAGE = 'ดึงข้อมูลไม่สำเร็จ กรุณาลองใหม่อีกครั้ง'
const MODELS = [
  { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 3, totalCount: 5 },
]

// client จำลองที่สลับได้ว่าฐานข้อมูลล่มหรือไม่ แบบเดียวกับหลังบ้านที่ตอบ 503 UPSTREAM_UNAVAILABLE (CON-TECH-03)
function switchableClient() {
  const state = { down: false, calls: 0 }
  return {
    state,
    listModels: async () => {
      state.calls++
      if (state.down) throw new ApiError(503, 'UPSTREAM_UNAVAILABLE', DOWN_MESSAGE)
      return { items: MODELS, total: MODELS.length }
    },
    listItems: async () => ({ items: [], total: 0 }),
  }
}

test('AC-EQS-08 ฐานข้อมูลไม่ตอบ แสดงข้อความผิดพลาดพร้อมปุ่มลองใหม่ ไม่ค้างที่ "กำลังโหลด"', async () => {
  // Given: ฐานข้อมูลไม่ตอบสนอง
  const client = switchableClient()
  client.state.down = true

  // When: เปิดหน้ารายการอุปกรณ์
  render(<EquipmentList client={client} />)

  // Then: ข้อความผิดพลาดจากหลังบ้าน + ปุ่มลองใหม่ และไม่ค้างที่ "กำลังโหลด..."
  expect((await screen.findByRole('alert')).textContent).toBe(DOWN_MESSAGE)
  expect(screen.getByRole('button', { name: 'ลองใหม่' })).toBeTruthy()
  expect(screen.queryByText('กำลังโหลด...')).toBeNull()
})

test('AC-EQS-08 เคยโหลดสำเร็จแล้วฐานข้อมูลล่ม ต้องไม่แสดงรายการเก่า และกดลองใหม่แล้วกลับมาได้', async () => {
  const client = switchableClient()
  render(<EquipmentList client={client} />)
  await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })

  // ฐานข้อมูลล่มระหว่างใช้งาน แล้วผู้ใช้เปลี่ยนตัวกรอง (หน้าจอต้องโหลดใหม่)
  client.state.down = true
  fireEvent.click(within(screen.getByRole('group', { name: 'สถานะ' })).getByRole('button', { name: 'ว่าง' }))

  // ไม่แสดงข้อมูลเก่าที่อาจไม่ถูกต้อง: รายการเดิมหายไปเหลือแค่ข้อผิดพลาด
  await screen.findByRole('alert')
  expect(screen.queryByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeNull()
  expect(screen.queryByText('ว่าง 3 / ทั้งหมด 5')).toBeNull()

  // ฐานข้อมูลกลับมา กดลองใหม่แล้วเห็นรายการอีกครั้ง
  client.state.down = false
  const callsBefore = client.state.calls
  fireEvent.click(screen.getByRole('button', { name: 'ลองใหม่' }))
  expect(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeTruthy()
  await waitFor(() => expect(client.state.calls).toBeGreaterThan(callsBefore))
  expect(screen.queryByRole('alert')).toBeNull()
})

// AC-EQS-08 (FR-EQS-08) กรณีหลังบ้านค้าง (T-26, Q-08)
//   Given หลังบ้านไม่ตอบเลย
//   When  เปิดหน้ารายการอุปกรณ์
//   Then  รอไม่เกิน EQUIPMENT_REQUEST_TIMEOUT_MS แล้วแสดงข้อความผิดพลาดพร้อมปุ่มลองใหม่ ไม่ค้างที่ "กำลังโหลด"
test('AC-EQS-08 หลังบ้านไม่ตอบ หน้าจอเลิกรอเมื่อครบ EQUIPMENT_REQUEST_TIMEOUT_MS แล้วแสดงข้อผิดพลาดพร้อมปุ่มลองใหม่', async () => {
  // Given: หลังบ้านไม่ตอบเลย
  vi.useFakeTimers()
  vi.stubGlobal('fetch', hangingFetch())

  // When: เปิดหน้ารายการอุปกรณ์ด้วย client จริง
  render(<EquipmentList client={api} />)

  // ยังไม่ครบเวลา: ยังรออยู่ ไม่แสดงข้อผิดพลาดเร็วเกิน
  await act(() => vi.advanceTimersByTimeAsync(EQUIPMENT_REQUEST_TIMEOUT_MS - 1))
  expect(screen.queryByRole('alert')).toBeNull()
  expect(screen.getByText('กำลังโหลด...')).toBeTruthy()

  // Then: ครบเวลาแล้วแสดงข้อความผิดพลาดพร้อมปุ่มลองใหม่ และไม่ค้างที่ "กำลังโหลด..."
  await act(() => vi.advanceTimersByTimeAsync(1))
  vi.useRealTimers()
  expect((await screen.findByRole('alert')).textContent).toBe(DEFAULT_ERROR_MESSAGE)
  expect(screen.getByRole('button', { name: 'ลองใหม่' })).toBeTruthy()
  expect(screen.queryByText('กำลังโหลด...')).toBeNull()
})

// T-26 timeout ใช้เฉพาะ request ของ 001 (ทีมตัดสิน 2569-11-12) request ของ 003 ต้องรอหลังบ้านเหมือนเดิม
// เพราะ 003 ยอมให้ช้าได้ถึง 15 วินาที (NFR-PERF-CXL-01)
test('T-26 request ของ 003 ไม่ถูกตัดที่ EQUIPMENT_REQUEST_TIMEOUT_MS', async () => {
  vi.useFakeTimers()
  const fetchMock = hangingFetch()
  vi.stubGlobal('fetch', fetchMock)

  let settled = false
  const done = () => {
    settled = true
  }
  api.listMyPendingRequests().then(done, done)
  api.cancelRequest(101, 'ไม่ต้องใช้อุปกรณ์แล้ว').then(done, done)
  await vi.advanceTimersByTimeAsync(EQUIPMENT_REQUEST_TIMEOUT_MS * 2)

  expect(settled).toBe(false)
  for (const [, init] of fetchMock.mock.calls) expect(init.signal).toBeUndefined()
})
