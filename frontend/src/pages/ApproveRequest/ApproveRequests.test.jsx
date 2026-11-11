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
