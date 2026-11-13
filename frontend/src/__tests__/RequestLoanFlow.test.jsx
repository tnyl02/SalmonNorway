// T-41 ปุ่มส่งคำขอยืมของ 001 เปิดฟอร์มของ 002 พร้อมรหัสครุภัณฑ์ (FR-EQS-06 "ส่งข้อมูลอุปกรณ์ที่เลือกไปยัง UC-02", F-01 ใน rtm.md)
// ตรวจส่วน "ปุ่มส่งคำขอยืมกลับมาใช้งานได้" ของ AC-EQS-05 ในระบบจริงทั้งหน้า (App) ไม่ใช่แค่ callback
// ใช้ fetch จำลองแทนหลังบ้าน รูปแบบตาม plan.md ข้อ 5 ของ 001
import { fireEvent, render, screen, within } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import App from '../App.jsx'
import { setToken } from '../api/client.js'

function json(body, status = 200) {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })
}

function stubBackend() {
  vi.stubGlobal('fetch', async (url) => {
    const u = String(url)
    if (u.endsWith('/v1/equipment/models')) {
      return json({ items: [{ modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 1, totalCount: 2 }], total: 1 })
    }
    if (u.includes('/v1/equipment/models/1/items')) {
      return json({
        items: [
          { assetCode: 'CS-NB-0001', status: 'available', canRequestLoan: true },
          { assetCode: 'CS-NB-0004', status: 'borrowed', canRequestLoan: false },
        ],
        total: 2,
      })
    }
    if (u.endsWith('/v1/equipment/items/CS-NB-0001')) {
      return json({ assetCode: 'CS-NB-0001', status: 'available', canRequestLoan: true })
    }
    return json({ code: 'ITEM_NOT_FOUND', message: 'ไม่พบอุปกรณ์รหัสนี้' }, 404)
  })
}

afterEach(() => {
  vi.unstubAllGlobals()
  setToken('mock-student-001')
})

async function openModelA() {
  setToken('mock-student-001')
  stubBackend()
  render(<App />)
  fireEvent.click(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ }))
  return screen.findByRole('table')
}

test('AC-EQS-05 กดปุ่มในตารางรายชิ้น เปิดฟอร์มของ 002 พร้อมรหัสครุภัณฑ์ของชิ้นนั้น', async () => {
  const table = await openModelA()
  const row = within(table).getByRole('button', { name: 'CS-NB-0001' }).closest('tr')

  fireEvent.click(within(row).getByRole('button', { name: 'ส่งคำขอยืม' }))

  // หน้าปัจจุบันเป็นหน้า "ส่งคำขอยืม" ของ 002 และช่องรหัสครุภัณฑ์มีรหัสของชิ้นที่กด
  expect(screen.getByRole('button', { name: 'ส่งคำขอยืม', current: 'page' })).toBeTruthy()
  expect(screen.getByLabelText('รหัสครุภัณฑ์').value).toBe('CS-NB-0001')
})

test('FR-EQS-06 กดปุ่มในหน้ารายละเอียด เปิดฟอร์มพร้อมรหัสครุภัณฑ์เดียวกัน', async () => {
  const table = await openModelA()
  fireEvent.click(within(table).getByRole('button', { name: 'CS-NB-0001' }))

  // รอหน้ารายละเอียดโหลดเสร็จ แล้วกดปุ่มในหน้า (ไม่ใช่ปุ่มบนแถบด้านบนที่ชื่อเดียวกัน)
  await screen.findByText('ว่าง')
  const inDetail = screen.getAllByRole('button', { name: 'ส่งคำขอยืม' }).find((b) => !b.closest('nav'))
  fireEvent.click(inDetail)

  expect(screen.getByLabelText('รหัสครุภัณฑ์').value).toBe('CS-NB-0001')
})

test('เปิดฟอร์มจากแถบด้านบนเองหลังเคยกดปุ่ม ช่องรหัสครุภัณฑ์ว่าง (ไม่มีรหัสค้างจากครั้งก่อน)', async () => {
  const table = await openModelA()
  const row = within(table).getByRole('button', { name: 'CS-NB-0001' }).closest('tr')
  fireEvent.click(within(row).getByRole('button', { name: 'ส่งคำขอยืม' }))

  const nav = screen.getByRole('navigation')
  fireEvent.click(within(nav).getByRole('button', { name: 'รายการอุปกรณ์' }))
  fireEvent.click(within(nav).getByRole('button', { name: 'ส่งคำขอยืม' }))

  expect(screen.getByLabelText('รหัสครุภัณฑ์').value).toBe('')
})
