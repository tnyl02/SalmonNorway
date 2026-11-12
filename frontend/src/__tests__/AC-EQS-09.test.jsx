// AC-EQS-09 (IF-IDP-01) ใช้ fetch จำลองแทนหลังบ้าน (ตอบ 401 เมื่อไม่มี token แบบเดียวกับ auth.Middleware)
//   Given session หมดอายุ
//   When  เรียกดูรายการอุปกรณ์
//   Then  ส่งกลับหน้าเข้าสู่ระบบพร้อมข้อความแจ้ง และกลับมาหน้าเดิมได้หลังล็อกอินสำเร็จ
// หน้าเข้าสู่ระบบเป็นแบบจำลอง (ระบบจริงรอ Q-SH-03) และจำหน้าเดิมไว้ในหน่วยความจำ (ทีมตัดสิน 2569-11-12)
import { fireEvent, render, screen } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import App, { safePageKey } from '../App.jsx'
import { getToken, setToken } from '../api/client.js'

const EXPIRED_MESSAGE = 'กรุณาเข้าสู่ระบบก่อนใช้งาน'
const VALID = new Set(['mock-student-001', 'mock-staff-001'])

function json(body, status = 200) {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })
}

// หลังบ้านจำลอง: token ที่ไม่รู้จักหรือไม่มี token ได้ 401 UNAUTHENTICATED (CON-TECH-03)
function stubBackend() {
  const seen = []
  vi.stubGlobal('fetch', async (url, opts = {}) => {
    const auth = opts.headers?.Authorization ?? ''
    seen.push({ url: String(url), auth })
    if (!VALID.has(auth.replace('Bearer ', ''))) return json({ code: 'UNAUTHENTICATED', message: EXPIRED_MESSAGE }, 401)
    if (String(url).includes('/v1/equipment/models')) {
      return json({ items: [{ modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 3, totalCount: 5 }], total: 1 })
    }
    return json({ items: [] })
  })
  return seen
}

afterEach(() => {
  vi.unstubAllGlobals()
  setToken('mock-student-001')
})

test('AC-EQS-09 session หมดอายุ ไปหน้าเข้าสู่ระบบพร้อมข้อความ แล้วล็อกอินกลับมาดูรายการได้', async () => {
  // Given: session หมดอายุ (ไม่มี token ใน browser)
  const seen = stubBackend()
  setToken(null)

  // When: เรียกดูรายการอุปกรณ์
  render(<App />)

  // Then: หน้าเข้าสู่ระบบพร้อมข้อความแจ้ง ไม่มีรายการอุปกรณ์
  expect(await screen.findByText('เข้าสู่ระบบ (จำลอง)')).toBeTruthy()
  expect(screen.getByRole('alert').textContent).toBe(EXPIRED_MESSAGE)
  expect(screen.queryByText(/Notebook/)).toBeNull()

  // Then: ล็อกอินสำเร็จแล้วกลับมาหน้ารายการ และ request ใหม่ส่ง token ที่ได้
  fireEvent.click(screen.getByRole('button', { name: /student-001/ }))
  expect(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeTruthy()
  expect(getToken()).toBe('mock-student-001')
  expect(seen.at(-1).auth).toBe('Bearer mock-student-001')
})

test('AC-EQS-09 กลับมาหน้าเดิม: หมดอายุตอนอยู่หน้าอื่น ล็อกอินแล้วกลับหน้านั้น ไม่กลับหน้าแรก', async () => {
  stubBackend()
  setToken('mock-student-001')
  render(<App />)
  await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })

  // session หมดอายุระหว่างใช้งาน แล้วผู้ใช้ไปหน้า "คำขอที่รออนุมัติของฉัน"
  setToken(null)
  fireEvent.click(screen.getByRole('button', { name: 'คำขอที่รออนุมัติของฉัน' }))
  expect(await screen.findByText('เข้าสู่ระบบ (จำลอง)')).toBeTruthy()

  // ล็อกอินแล้วกลับมาหน้าที่ค้างไว้
  fireEvent.click(screen.getByRole('button', { name: /staff-001/ }))
  const back = await screen.findByRole('button', { name: 'คำขอที่รออนุมัติของฉัน' })
  expect(back.getAttribute('aria-current')).toBe('page')
})

test('returnTo ภายนอกระบบถูกปฏิเสธ กลับหน้าแรกแทน (plan ข้อ 6.5 กัน open redirect)', () => {
  expect(safePageKey('https://evil.example/login')).toBe('equipment')
  expect(safePageKey('//evil.example')).toBe('equipment')
  expect(safePageKey(undefined)).toBe('equipment')
  expect(safePageKey('my-pending')).toBe('my-pending')
})
