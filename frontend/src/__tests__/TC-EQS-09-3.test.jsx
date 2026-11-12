// TC-EQS-09-3 (AC-EQS-09 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import App from '../App.jsx'
import { setToken } from '../api/client.js'

const MESSAGE = 'กรุณาเข้าสู่ระบบก่อนใช้งาน'

function json(body, status = 200) {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })
}

afterEach(() => {
  vi.unstubAllGlobals()
  setToken('mock-student-001')
})

test('TC-EQS-09-3 session หมดอายุที่หน้ารายการ ไปหน้าเข้าสู่ระบบพร้อมข้อความ ล็อกอินแล้วกลับหน้ารายการ', async () => {
  // Given อยู่หน้ารายการอุปกรณ์ client จำลองตอบ 401 (ไม่มี token = session หมดอายุ)
  vi.stubGlobal('fetch', async (url, opts = {}) => {
    if (!opts.headers?.Authorization) return json({ code: 'UNAUTHENTICATED', message: MESSAGE }, 401)
    if (String(url).includes('/v1/equipment/models')) {
      return json({ items: [{ modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 3, totalCount: 5 }], total: 1 })
    }
    return json({ items: [], total: 0 })
  })
  setToken(null)

  // When เรียกรายการ
  render(<App />)

  // Then (1) แสดงหน้าเข้าสู่ระบบ
  expect(await screen.findByText('เข้าสู่ระบบ (จำลอง)')).toBeTruthy()
  // Then (2) มีข้อความแจ้ง (ถ้อยคำทีมเลือก plan ข้อ 5)
  expect(screen.getByRole('alert').textContent).toBe(MESSAGE)

  // When เข้าสู่ระบบในหน้าเข้าสู่ระบบ (จำลอง)
  fireEvent.click(screen.getByRole('button', { name: /student-001/ }))

  // Then (3) กลับมาหน้ารายการอุปกรณ์ และแสดงรายการ
  expect(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeTruthy()
  expect(screen.getByRole('button', { name: 'รายการอุปกรณ์' }).getAttribute('aria-current')).toBe('page')
})
