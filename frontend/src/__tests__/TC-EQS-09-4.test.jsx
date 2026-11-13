// TC-EQS-09-4 (AC-EQS-09 ขอบ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
// Then ทีมตัดสิน 2569-11-12 รับพฤติกรรมของโค้ด: "หน้าเดิม" = แท็บเดิม
import { fireEvent, render, screen } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import App from '../App.jsx'
import { setToken } from '../api/client.js'

function json(body, status = 200) {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })
}

afterEach(() => {
  vi.unstubAllGlobals()
  setToken('mock-student-001')
})

test('TC-EQS-09-4 session หมดอายุตอนเปิดรายละเอียด ล็อกอินแล้วกลับแท็บรายการอุปกรณ์ที่หน้ารายการ', async () => {
  // Given หลังบ้านจำลอง: ไม่มี token ตอบ 401 / มี token ตอบรายการรุ่น A และรายชิ้น CS-NB-0001
  vi.stubGlobal('fetch', async (url, opts = {}) => {
    const u = String(url)
    if (!opts.headers?.Authorization) {
      return json({ code: 'UNAUTHENTICATED', message: 'กรุณาเข้าสู่ระบบก่อนใช้งาน' }, 401)
    }
    if (u.includes('/models/1/items')) {
      return json({ items: [{ assetCode: 'CS-NB-0001', status: 'available', canRequestLoan: true }], total: 1 })
    }
    if (u.includes('/v1/equipment/models')) {
      return json({ items: [{ modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 3, totalCount: 5 }], total: 1 })
    }
    if (u.includes('/items/CS-NB-0001')) return json({ assetCode: 'CS-NB-0001', status: 'available', canRequestLoan: true })
    return json({ items: [], total: 0 })
  })
  setToken('mock-student-001')
  render(<App />)
  fireEvent.click(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ }))
  const openDetail = await screen.findByRole('button', { name: 'CS-NB-0001' })

  // When session หมดอายุ แล้วเปิดรายละเอียด CS-NB-0001 (client ตอบ 401)
  setToken(null)
  fireEvent.click(openDetail)

  // Then (1) แสดงหน้าเข้าสู่ระบบ
  expect(await screen.findByText('เข้าสู่ระบบ (จำลอง)')).toBeTruthy()

  // When เข้าสู่ระบบ
  fireEvent.click(screen.getByRole('button', { name: /student-001/ }))

  // Then (2) กลับมาแท็บ "รายการอุปกรณ์" ที่หน้ารายการ ไม่ได้กลับไปรายละเอียด CS-NB-0001
  expect(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeTruthy()
  expect(screen.getByRole('button', { name: 'รายการอุปกรณ์' }).getAttribute('aria-current')).toBe('page')
  expect(screen.queryByRole('heading', { name: 'CS-NB-0001' })).toBeNull()
})
