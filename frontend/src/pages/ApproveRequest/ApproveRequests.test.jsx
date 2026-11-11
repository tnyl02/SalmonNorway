// test หน้าจอ ApproveRequests ของ 004 ใช้ client จำลองแทนหลังบ้าน (plan.md ข้อ 6)
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import ApproveRequests from './ApproveRequests.jsx'
import { APPROVAL_MESSAGES } from '../../config.js'

// คำตอบแบบเดียวกับ GET /v1/approvals/requests (หลังบ้านกรองภาควิชาและเรียงลำดับแล้ว TestAC_APV_11_HTTP)
function fakeClient() {
  return {
    getApprovalRequests: async () => ({
      loadedAt: '2026-11-12T03:00:00Z',
      items: [
        { requestId: 1, requesterName: 'นักศึกษาตัวอย่าง ก', studentCode: 'student-001', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)',
          assetCode: 'AP-TEST-0001', purpose: 'study', startAt: '2026-11-12T02:00:00Z', endAt: '2026-11-12T09:00:00Z',
          submittedAt: '2026-11-12T01:30:00Z', isFirst: true },
        { requestId: 2, requesterName: 'นักศึกษาตัวอย่าง ข', studentCode: 'student-002', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)',
          assetCode: 'AP-TEST-0001', purpose: 'other', purposeOther: 'ถ่ายวิดีโอโครงงาน', startAt: '2026-11-12T02:00:00Z',
          endAt: '2026-11-12T09:00:00Z', submittedAt: '2026-11-12T01:40:00Z', isFirst: false },
      ],
    }),
  }
}

// AC-APV-11 (FR-APV-01 ก)
//   Given มีคำขอ requested ของภาควิชาตนเองหลายคำขอ (หลังบ้านคัดคำขอที่พิจารณาแล้วและของภาควิชาอื่นออกแล้ว)
//   When  เจ้าหน้าที่เปิดหน้ารายการ
//   Then  เห็นคำขอเรียงตามเวลาที่ส่งจากเก่าไปใหม่ พร้อมชื่อผู้ขอ อุปกรณ์ วัตถุประสงค์ ช่วงเวลาที่ต้องการยืม
test('AC-APV-11 แสดงคำขอที่รอพิจารณาตามลำดับ พร้อมชื่อผู้ขอ อุปกรณ์ วัตถุประสงค์ และช่วงเวลา', async () => {
  render(<ApproveRequests client={fakeClient()} />)

  const table = await screen.findByRole('table')
  const rows = within(table).getAllByRole('row').slice(1) // ตัดแถวหัวตาราง
  expect(rows).toHaveLength(2)

  // ลำดับตามที่หลังบ้านส่งมา (เวลาที่ส่งจากเก่าไปใหม่)
  expect(within(rows[0]).getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
  expect(within(rows[1]).getByText('นักศึกษาตัวอย่าง ข')).toBeTruthy()

  // อุปกรณ์ วัตถุประสงค์ (รวมข้อความที่พิมพ์เมื่อเลือก "อื่น ๆ") และช่วงเวลา
  expect(within(rows[0]).getByText(/อุปกรณ์ทดลองอนุมัติ \(004\)/)).toBeTruthy()
  expect(within(rows[0]).getByText(/AP-TEST-0001/)).toBeTruthy()
  expect(within(rows[0]).getByText('เพื่อการเรียน')).toBeTruthy()
  expect(within(rows[1]).getByText('ถ่ายวิดีโอโครงงาน')).toBeTruthy()
  expect(rows[0].cells[3].textContent).toMatch(/ - /) // ช่วงเวลา เริ่ม - คืน
})

// FR-APV-01 (ซ), CON-TECH-03: เลือกข้อความจากรหัส เช่น นักศึกษาเปิดหน้านี้ได้ FORBIDDEN_ROLE
test('AC-APV-11 หลังบ้านตอบ error แสดงข้อความตามรหัสและไม่แสดงรายการ', async () => {
  const client = { getApprovalRequests: async () => Promise.reject(Object.assign(new Error('403'), { code: 'FORBIDDEN_ROLE' })) }
  render(<ApproveRequests client={client} />)

  expect((await screen.findByRole('alert')).textContent).toBe(APPROVAL_MESSAGES.FORBIDDEN_ROLE)
  expect(screen.queryByRole('table')).toBeNull()
})

// AC-APV-01 (FR-APV-01, CON-DATA-01) ฝั่งหน้าจอ (สถานะคำขอและอุปกรณ์ตรวจที่หลังบ้านใน TestAC_APV_01)
//   Given เจ้าหน้าที่เห็นคำขอแรกของอุปกรณ์ที่ "ว่าง" ในรายการ
//   When  กด "อนุมัติ" ที่คำขอนั้น
//   Then  ส่งรหัสคำขอและ loadedAt ของรายการไปหลังบ้าน แล้วคำขอนั้นหายจากรายการ
test('AC-APV-01 กดอนุมัติแล้วส่ง loadedAt และคำขอหายจากรายการ', async () => {
  const calls = []
  const client = {
    ...fakeClient(),
    approveRequest: async (requestId, loadedAt) => {
      calls.push({ requestId, loadedAt })
      return { requestId, status: 'approved', decidedAt: '2026-11-12T03:05:00Z', rejectedRequestIds: [] }
    },
  }
  render(<ApproveRequests client={client} />)

  const table = await screen.findByRole('table')
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  await waitFor(() => expect(within(table).getAllByRole('row').slice(1)).toHaveLength(1))
  expect(calls).toEqual([{ requestId: 1, loadedAt: '2026-11-12T03:00:00Z' }])
  expect(screen.queryByText('นักศึกษาตัวอย่าง ก')).toBeNull()
  expect(screen.getByText('นักศึกษาตัวอย่าง ข')).toBeTruthy()
})

// AC-APV-02 (FR-APV-01) ฝั่งหน้าจอ (สถานะคำขอและอุปกรณ์ตรวจที่หลังบ้านใน TestAC_APV_02)
//   Given เจ้าหน้าที่เห็นคำขอสถานะ requested ในรายการ
//   When  กด "ปฏิเสธ" พิมพ์เหตุผล แล้วกด "ยืนยันการปฏิเสธ"
//   Then  ส่งรหัสคำขอ เหตุผล และ loadedAt ไปหลังบ้าน แล้วคำขอนั้นหายจากรายการ
test('AC-APV-02 ปฏิเสธพร้อมเหตุผลแล้วคำขอหายจากรายการ', async () => {
  const calls = []
  const client = {
    ...fakeClient(),
    rejectRequest: async (requestId, reason, loadedAt) => {
      calls.push({ requestId, reason, loadedAt })
      return { requestId, status: 'rejected', decidedAt: '2026-11-12T03:05:00Z', rejectedRequestIds: [] }
    },
  }
  render(<ApproveRequests client={client} />)

  const table = await screen.findByRole('table')
  const second = within(table).getAllByRole('row')[2]
  fireEvent.click(within(second).getByRole('button', { name: 'ปฏิเสธ' }))
  fireEvent.change(within(second).getByLabelText('เหตุผลในการปฏิเสธ'), { target: { value: 'อุปกรณ์ต้องใช้สอบ' } })
  fireEvent.click(within(second).getByRole('button', { name: 'ยืนยันการปฏิเสธ' }))

  await waitFor(() => expect(within(table).getAllByRole('row').slice(1)).toHaveLength(1))
  expect(calls).toEqual([{ requestId: 2, reason: 'อุปกรณ์ต้องใช้สอบ', loadedAt: '2026-11-12T03:00:00Z' }])
  expect(screen.queryByText('นักศึกษาตัวอย่าง ข')).toBeNull()
})

// AC-APV-06 (FR-APV-01 ง) ฝั่งหน้าจอ (หลังบ้านตรวจซ้ำใน TestAC_APV_06)
//   Given คำขอสถานะ requested ในรายการ
//   When  เจ้าหน้าที่กดปฏิเสธโดยไม่ระบุเหตุผล (ว่าง หรือมีแต่ช่องว่าง)
//   Then  ระบบไม่ปฏิเสธ แสดง "กรุณาระบุเหตุผลในการปฏิเสธ" และคำขอยังอยู่ในรายการ
test('AC-APV-06 ปฏิเสธโดยไม่ระบุเหตุผลแสดงข้อความและไม่ส่งไปหลังบ้าน', async () => {
  let called = 0
  const client = { ...fakeClient(), rejectRequest: async () => { called++ } }
  render(<ApproveRequests client={client} />)

  const table = await screen.findByRole('table')
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'ปฏิเสธ' }))
  for (const value of ['', '   ']) {
    fireEvent.change(within(first).getByLabelText('เหตุผลในการปฏิเสธ'), { target: { value } })
    fireEvent.click(within(first).getByRole('button', { name: 'ยืนยันการปฏิเสธ' }))
    expect((await screen.findByRole('alert')).textContent).toBe('กรุณาระบุเหตุผลในการปฏิเสธ')
  }
  expect(called).toBe(0)
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
})

// AC-APV-06: ถ้าหลังบ้านตอบ VALIDATION_FAILED (ตรวจซ้ำ T-12) หน้าจอแสดงข้อความเดียวกันและคำขอยังอยู่
test('AC-APV-06 หลังบ้านตอบ VALIDATION_FAILED แสดงข้อความตามรหัส', async () => {
  const client = {
    ...fakeClient(),
    rejectRequest: async () => Promise.reject(Object.assign(new Error('422'), { code: 'VALIDATION_FAILED' })),
  }
  render(<ApproveRequests client={client} />)

  const table = await screen.findByRole('table')
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'ปฏิเสธ' }))
  fireEvent.change(within(first).getByLabelText('เหตุผลในการปฏิเสธ'), { target: { value: 'มีเหตุผล' } })
  fireEvent.click(within(first).getByRole('button', { name: 'ยืนยันการปฏิเสธ' }))

  expect((await screen.findByRole('alert')).textContent).toBe(APPROVAL_MESSAGES.VALIDATION_FAILED)
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
})

// AC-APV-07 (FR-APV-01 ข, DOM-REQ-01) ฝั่งหน้าจอ (หลังบ้านตรวจซ้ำใน TestAC_APV_07)
//   Given อุปกรณ์รายชิ้นเดียวกันมีคำขอ A (ส่งก่อน) และ B สถานะ requested
//   When  เจ้าหน้าที่ดูรายการ
//   Then  แถว B ไม่มีปุ่มอนุมัติ (ปฏิเสธได้) ส่วนแถว A มีปุ่มอนุมัติ
test('AC-APV-07 แสดงปุ่มอนุมัติเฉพาะคำขอแรก', async () => {
  render(<ApproveRequests client={fakeClient()} />)

  const table = await screen.findByRole('table')
  const [a, b] = within(table).getAllByRole('row').slice(1)
  expect(within(a).getByRole('button', { name: 'อนุมัติ' })).toBeTruthy()
  expect(within(b).queryByRole('button', { name: 'อนุมัติ' })).toBeNull()
  expect(within(b).getByRole('button', { name: 'ปฏิเสธ' })).toBeTruthy()
})

// AC-APV-07: หน้าจอเก่า (คำขอก่อนหน้าเพิ่งเข้ามา) กดอนุมัติแล้วหลังบ้านตอบ NOT_FIRST_IN_QUEUE
//   Then แสดง "ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน" และคำขอยังอยู่ในรายการ
test('AC-APV-07 หลังบ้านตอบ NOT_FIRST_IN_QUEUE แสดงข้อความตามรหัส', async () => {
  const client = {
    ...fakeClient(),
    approveRequest: async () => Promise.reject(Object.assign(new Error('409'), { code: 'NOT_FIRST_IN_QUEUE' })),
  }
  render(<ApproveRequests client={client} />)

  const table = await screen.findByRole('table')
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  expect((await screen.findByRole('alert')).textContent).toBe('ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน')
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
})

// AC-APV-08 (FR-APV-01 ค, CON-DATA-01) ฝั่งหน้าจอ (หลังบ้านตรวจใน TestAC_APV_08)
//   Given คำขอแรกของอุปกรณ์รายชิ้นที่สถานะไม่ใช่ "ว่าง"
//   When  เจ้าหน้าที่กดอนุมัติ แล้วหลังบ้านตอบ ITEM_NOT_AVAILABLE
//   Then  แสดง "อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้" และคำขอยังอยู่ในรายการ
test('AC-APV-08 หลังบ้านตอบ ITEM_NOT_AVAILABLE แสดงข้อความตามรหัส', async () => {
  const client = {
    ...fakeClient(),
    approveRequest: async () => Promise.reject(Object.assign(new Error('409'), { code: 'ITEM_NOT_AVAILABLE' })),
  }
  render(<ApproveRequests client={client} />)

  const table = await screen.findByRole('table')
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  expect((await screen.findByRole('alert')).textContent).toBe('อุปกรณ์ชิ้นนี้ไม่ได้อยู่ในสถานะว่างแล้ว อนุมัติไม่ได้')
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
})

// AC-APV-16 (FR-APV-01 ค, DOM-LOAN-01) ฝั่งหน้าจอ (หลังบ้านตรวจใน TestAC_APV_16) [ID รอทีมอนุมัติ Q-APV-13]
//   Given คำขอแรกของอุปกรณ์ที่ "ว่าง" และเวลาปัจจุบันเลย 16:00 น. ของวันที่เริ่มยืมแล้ว
//   When  เจ้าหน้าที่กดอนุมัติ แล้วหลังบ้านตอบ APPROVAL_CUTOFF_PASSED
//   Then  แสดง "เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้" และคำขอยังอยู่ในรายการ (ปฏิเสธเองได้ Q-APV-16)
test('AC-APV-16 หลังบ้านตอบ APPROVAL_CUTOFF_PASSED แสดงข้อความตามรหัส', async () => {
  const client = {
    ...fakeClient(),
    approveRequest: async () => Promise.reject(Object.assign(new Error('409'), { code: 'APPROVAL_CUTOFF_PASSED' })),
  }
  render(<ApproveRequests client={client} />)

  const table = await screen.findByRole('table')
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  expect((await screen.findByRole('alert')).textContent).toBe('เลย 16:00 น. ของวันที่ยืมแล้ว อนุมัติไม่ได้')
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
  expect(within(first).getByRole('button', { name: 'ปฏิเสธ' })).toBeTruthy()
})
