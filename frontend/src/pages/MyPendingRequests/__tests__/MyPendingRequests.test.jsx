// AC-CXL-08 (FR-CXL-06) หน้า "คำขอที่รออนุมัติของฉัน" ใช้ client จำลองแทนหลังบ้าน (T-05)
//   Given นักศึกษามีคำขอที่ยังยกเลิกได้ 2 รายการ ส่งคนละเวลา (หลังบ้านเรียงใหม่ไปเก่ามาแล้ว)
//   When  เปิดหน้า "คำขอที่รออนุมัติของฉัน"
//   Then  เห็น 2 แถว ส่งล่าสุดอยู่บนสุด แต่ละแถวมีชื่อรุ่น รหัสครุภัณฑ์ เวลาที่ส่งแบบ Asia/Bangkok และปุ่มยกเลิก
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import App from '../../../App.jsx'
import { ApiError, setToken } from '../../../api/client.js'
import MyPendingRequests from '../MyPendingRequests.jsx'
import { REQ_A1, REQ_A2, fakeClient } from '../../../__tests__/helpers/cxlFakeClient.js'

test('AC_CXL_08 shows own pending requests newest first', async () => {
  // Given
  const client = fakeClient([REQ_A1, REQ_A2])

  // When
  render(<MyPendingRequests client={client} />)

  // Then
  const rows = await screen.findAllByRole('listitem')
  expect(rows).toHaveLength(2)
  expect(client.listMyPendingRequests).toHaveBeenCalledTimes(1)

  expect(rows[0].textContent).toContain('Notebook ตัวอย่างรุ่น A')
  expect(rows[0].textContent).toContain('CS-NB-0001')
  expect(rows[0].textContent).toContain('2569-11-10 10:30') // 03:30Z = 10:30 Asia/Bangkok

  expect(rows[1].textContent).toContain('Projector ตัวอย่างรุ่น B')
  expect(rows[1].textContent).toContain('CS-PJ-0001')
  expect(rows[1].textContent).toContain('2569-11-10 09:00')

  for (const row of rows) {
    expect(within(row).getByRole('button', { name: 'ยกเลิก' })).toBeTruthy()
  }
})

// AC-CXL-12 (FR-CXL-06) ไม่มีคำขอแสดง "ไม่มีคำขอที่รออนุมัติ" โหลดไม่สำเร็จแสดงข้อความและปุ่มลองใหม่ (T-12)
//   Given นักศึกษาไม่มีคำขอที่ยังยกเลิกได้
//   When  เปิดหน้า "คำขอที่รออนุมัติของฉัน"
//   Then  แสดง "ไม่มีคำขอที่รออนุมัติ" และไม่มีแถว
//   และกรณีโหลดไม่สำเร็จ แสดง "โหลดรายการไม่สำเร็จ กรุณาลองใหม่อีกครั้ง" ไม่แสดงว่าไม่มีคำขอ กดลองใหม่แล้วโหลดได้
test('AC_CXL_12 shows empty message', async () => {
  // Given / When
  render(<MyPendingRequests client={fakeClient([])} />)

  // Then
  expect(await screen.findByText('ไม่มีคำขอที่รออนุมัติ')).toBeTruthy()
  expect(screen.queryAllByRole('listitem')).toHaveLength(0)
  expect(screen.queryByRole('alert')).toBeNull()
})

test('AC_CXL_12 load failed shows retry and recovers', async () => {
  // Given รอบแรกหลังบ้านตอบ 503 รอบที่สองตอบได้
  const client = fakeClient([REQ_A1])
  client.listMyPendingRequests.mockRejectedValueOnce(new ApiError(503, 'UPSTREAM_UNAVAILABLE'))

  // When
  render(<MyPendingRequests client={client} />)

  // Then ข้อความโหลดไม่สำเร็จ ไม่มีแถว ไม่บอกว่าไม่มีคำขอ
  expect((await screen.findByRole('alert')).textContent).toBe('โหลดรายการไม่สำเร็จ กรุณาลองใหม่อีกครั้ง')
  expect(screen.queryByText('ไม่มีคำขอที่รออนุมัติ')).toBeNull()
  expect(screen.queryAllByRole('listitem')).toHaveLength(0)

  // When กดลองใหม่
  fireEvent.click(screen.getByRole('button', { name: 'ลองใหม่' }))

  // Then โหลดรอบสองสำเร็จ ข้อความหาย เห็นรายการ
  await waitFor(() => expect(screen.queryAllByRole('listitem')).toHaveLength(1))
  expect(client.listMyPendingRequests).toHaveBeenCalledTimes(2)
  expect(screen.queryByRole('alert')).toBeNull()
})

// AC-CXL-02 (FR-CXL-02) ฝั่งหน้าจอ T-16
//   Given นักศึกษาเปิดหน้ารายการ แล้วคำขอถูกอนุมัติก่อนกดยืนยัน (หลังบ้านตอบ 409 ALREADY_APPROVED)
//   When  กดยกเลิก พิมพ์เหตุผล แล้วกดยืนยัน
//   Then  กล่องปิด แสดง "ยกเลิกไม่ได้ เพราะคำขอได้รับการอนุมัติแล้ว" และโหลดรายการใหม่ (คำขอที่อนุมัติแล้วหายไป)
test('AC_CXL_02 shows already approved message and reloads', async () => {
  // Given รอบแรกเห็น REQ_A1 หลังอนุมัติ หลังบ้านไม่ส่ง REQ_A1 แล้ว
  const client = fakeClient([REQ_A1, REQ_A2])
  client.listMyPendingRequests
    .mockResolvedValueOnce({ items: [REQ_A1, REQ_A2] })
    .mockResolvedValueOnce({ items: [REQ_A2] })
  client.cancelRequest.mockRejectedValueOnce(new ApiError(409, 'ALREADY_APPROVED'))
  render(<MyPendingRequests client={client} />)
  const rows = await screen.findAllByRole('listitem')

  // When
  fireEvent.click(within(rows[0]).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then
  expect((await screen.findByRole('alert')).textContent).toBe('ยกเลิกไม่ได้ เพราะคำขอได้รับการอนุมัติแล้ว')
  expect(screen.queryByRole('dialog')).toBeNull()
  await waitFor(() => expect(client.listMyPendingRequests).toHaveBeenCalledTimes(2))
  await waitFor(() => expect(screen.getAllByRole('listitem')).toHaveLength(1))
  expect(screen.getByRole('listitem').textContent).toContain('CS-PJ-0001')
  expect(screen.queryByRole('status')).toBeNull() // ไม่มีข้อความยกเลิกสำเร็จ
})

// AC-CXL-11 (FR-CXL-02) ฝั่งหน้าจอ T-18
//   Given นักศึกษาเปิดหน้ารายการ แล้วคำขอถูกปฏิเสธก่อนกดยืนยัน (หลังบ้านตอบ 409 NOT_CANCELLABLE)
//   When  กดยกเลิก พิมพ์เหตุผล แล้วกดยืนยัน
//   Then  กล่องปิด แสดง "คำขอนี้ยกเลิกไม่ได้" และโหลดรายการใหม่
test('AC_CXL_11 shows not cancellable message', async () => {
  // Given
  const client = fakeClient([REQ_A1, REQ_A2])
  client.listMyPendingRequests
    .mockResolvedValueOnce({ items: [REQ_A1, REQ_A2] })
    .mockResolvedValueOnce({ items: [REQ_A2] })
  client.cancelRequest.mockRejectedValueOnce(new ApiError(409, 'NOT_CANCELLABLE'))
  render(<MyPendingRequests client={client} />)
  const rows = await screen.findAllByRole('listitem')

  // When
  fireEvent.click(within(rows[0]).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then
  expect((await screen.findByRole('alert')).textContent).toBe('คำขอนี้ยกเลิกไม่ได้')
  expect(screen.queryByRole('dialog')).toBeNull()
  await waitFor(() => expect(client.listMyPendingRequests).toHaveBeenCalledTimes(2))
  await waitFor(() => expect(screen.getAllByRole('listitem')).toHaveLength(1))
  expect(screen.queryByRole('status')).toBeNull()
})

// ข้อความยกเลิกไม่ได้หายเมื่อเริ่มยกเลิกรายการถัดไป (แบบเดียวกับข้อความสำเร็จ T-07)
test('cancel error clears when opening next dialog', async () => {
  const client = fakeClient([REQ_A1, REQ_A2])
  client.cancelRequest.mockRejectedValueOnce(new ApiError(409, 'NOT_CANCELLABLE'))
  render(<MyPendingRequests client={client} />)
  fireEvent.click(within((await screen.findAllByRole('listitem'))[0]).getByRole('button', { name: 'ยกเลิก' }))
  let dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'เหตุผล' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))
  await screen.findByRole('alert')

  fireEvent.click(within((await screen.findAllByRole('listitem'))[0]).getByRole('button', { name: 'ยกเลิก' }))
  dialog = await screen.findByRole('dialog')
  expect(screen.queryByRole('alert')).toBeNull()
})

// AC-CXL-03 (FR-CXL-03) ฝั่งหน้าจอ T-20
//   Given นักศึกษาเปิดหน้ารายการ แล้วหลังบ้านตอบว่าไม่พบคำขอ (404 REQUEST_NOT_FOUND)
//   When  กดยกเลิก พิมพ์เหตุผล แล้วกดยืนยัน
//   Then  กล่องปิด แสดง "ไม่พบคำขอ" และโหลดรายการใหม่
test('AC_CXL_03 shows request not found message', async () => {
  // Given
  const client = fakeClient([REQ_A1, REQ_A2])
  client.listMyPendingRequests
    .mockResolvedValueOnce({ items: [REQ_A1, REQ_A2] })
    .mockResolvedValueOnce({ items: [REQ_A2] })
  client.cancelRequest.mockRejectedValueOnce(new ApiError(404, 'REQUEST_NOT_FOUND'))
  render(<MyPendingRequests client={client} />)
  const rows = await screen.findAllByRole('listitem')

  // When
  fireEvent.click(within(rows[0]).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then
  expect((await screen.findByRole('alert')).textContent).toBe('ไม่พบคำขอ')
  expect(screen.queryByRole('dialog')).toBeNull()
  await waitFor(() => expect(client.listMyPendingRequests).toHaveBeenCalledTimes(2))
  await waitFor(() => expect(screen.getAllByRole('listitem')).toHaveLength(1))
  expect(screen.queryByRole('status')).toBeNull()
})

// AC-CXL-07 (IF-IDP-01, NFR-SEC-01) ใช้ App และ api/client.js ของจริง fetch จำลองแทนหลังบ้าน (T-26)
//   Given session ของนักศึกษา A หมดอายุ
//   When  นักศึกษา A ส่งคำสั่งยกเลิกคำขอ
//   Then  ระบบไม่ยกเลิก ส่งกลับหน้าเข้าสู่ระบบพร้อมข้อความแจ้ง
//         และหลังล็อกอินสำเร็จกลับมาหน้า "คำขอที่รออนุมัติของฉัน" โดยคำขอยังไม่ถูกยกเลิก (ต้องกดยกเลิกใหม่)
// กลไก 401 เป็นของ 001 T-21 (จำชื่อแท็บ + safePageKey แทน returnTo แบบ path, plan.md ข้อ 4)
// ส่วน "ไม่มีบันทึกใน audit_log" เป็นของหลังบ้าน ตรวจใน T-31
const EXPIRED_MESSAGE = 'กรุณาเข้าสู่ระบบก่อนใช้งาน'

function jsonResponse(body, status = 200) {
  return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } })
}

// stubBackend จำลองหลังบ้านของ 003: token ที่ไม่รู้จักได้ 401 แบบ auth.Middleware, staff ได้ 403 FORBIDDEN_ROLE
function stubBackend() {
  const calls = []
  vi.stubGlobal('fetch', async (url, opts = {}) => {
    const token = (opts.headers?.Authorization ?? '').replace('Bearer ', '')
    const path = String(url)
    calls.push({ path, method: opts.method ?? 'GET', token })
    if (!['mock-student-001', 'mock-staff-001'].includes(token)) {
      return jsonResponse({ code: 'UNAUTHENTICATED', message: EXPIRED_MESSAGE }, 401)
    }
    if (token === 'mock-staff-001') {
      return jsonResponse({ code: 'FORBIDDEN_ROLE', message: 'เฉพาะนักศึกษาเท่านั้นที่ยกเลิกคำขอของตนเองได้' }, 403)
    }
    if (path.includes('/v1/my/pending-requests')) return jsonResponse({ items: [REQ_A1] })
    if (path.includes('/cancel')) return jsonResponse({ requestId: REQ_A1.requestId, status: 'cancelled' })
    return jsonResponse({ items: [], total: 0 }) // หน้าแรกของ 001
  })
  return calls
}

afterEach(() => {
  vi.unstubAllGlobals()
  setToken('mock-student-001')
})

test('AC_CXL_07 redirects to login and back', async () => {
  const calls = stubBackend()
  setToken('mock-student-001')
  render(<App />)
  fireEvent.click(await screen.findByRole('button', { name: 'คำขอที่รออนุมัติของฉัน' }))
  const row = await screen.findByRole('listitem')

  // Given: session ของ A หมดอายุระหว่างเปิดหน้านี้
  setToken(null)

  // When: ส่งคำสั่งยกเลิก
  fireEvent.click(within(row).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then: ไปหน้าเข้าสู่ระบบพร้อมข้อความแจ้ง
  expect(await screen.findByText('เข้าสู่ระบบ (จำลอง)')).toBeTruthy()
  expect(screen.getByRole('alert').textContent).toBe(EXPIRED_MESSAGE)
  const cancels = () => calls.filter((c) => c.method === 'POST')
  expect(cancels()).toHaveLength(1)
  expect(cancels()[0].token).toBe('') // ส่งไปโดยไม่มี session หลังบ้านจึงไม่ยกเลิก

  // Then: ล็อกอินแล้วกลับมาหน้า "คำขอที่รออนุมัติของฉัน" คำขอยังอยู่ ไม่ยกเลิกซ้ำอัตโนมัติ
  fireEvent.click(screen.getByRole('button', { name: /student-001/ }))
  const tab = await screen.findByRole('button', { name: 'คำขอที่รออนุมัติของฉัน' })
  expect(tab.getAttribute('aria-current')).toBe('page')
  expect((await screen.findByRole('listitem')).textContent).toContain(REQ_A1.assetCode)
  expect(screen.queryByRole('dialog')).toBeNull() // ต้องกดยกเลิกใหม่เอง
  expect(cancels()).toHaveLength(1)
})

// FR-CXL-03 (ASM-13) เจ้าหน้าที่เปิดหน้านี้ได้ 403 FORBIDDEN_ROLE (plan.md ข้อ 4, ทีมตัดสิน 2569-11-12 ทาง ก, T-26)
//   Then แสดง message จากหลังบ้าน ไม่มีปุ่มลองใหม่ ไม่มีรายการ
test('FR_CXL_03 staff sees forbidden message without retry', async () => {
  stubBackend()
  setToken('mock-staff-001')
  render(<App />)
  fireEvent.click(await screen.findByRole('button', { name: 'คำขอที่รออนุมัติของฉัน' }))

  expect((await screen.findByRole('alert')).textContent).toBe('เฉพาะนักศึกษาเท่านั้นที่ยกเลิกคำขอของตนเองได้')
  expect(screen.queryByRole('button', { name: 'ลองใหม่' })).toBeNull()
  expect(screen.queryByRole('listitem')).toBeNull()
  expect(screen.queryByText('ไม่มีคำขอที่รออนุมัติ')).toBeNull()
})
