// test หน้าจอ FailedEmails ของ 004 ใช้ client จำลองแทนหลังบ้าน (plan.md ข้อ 6)
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import FailedEmails from './FailedEmails.jsx'
import { APPROVAL_MESSAGES } from '../../config.js'

// คำตอบแบบเดียวกับ GET /v1/approvals/failed-emails (หลังบ้านกรองภาควิชาแล้ว TestAC_APV_14)
function fakeClient(overrides = {}) {
  return {
    getFailedEmails: async () => ({
      items: [
        { emailId: 7, requestId: 1, kind: 'approved', requesterName: 'นักศึกษาตัวอย่าง ก',
          modelName: 'อุปกรณ์ทดลองอนุมัติ (004)', assetCode: 'AP-TEST-0002', attempts: 4,
          lastError: 'smtp: connection refused', lastAttemptAt: '2026-11-12T03:03:00Z' },
        { emailId: 8, requestId: 2, kind: 'auto_rejected', requesterName: 'นักศึกษาตัวอย่าง ข',
          modelName: 'อุปกรณ์ทดลองอนุมัติ (004)', assetCode: 'AP-TEST-0003', attempts: 4,
          lastError: 'smtp: timeout', lastAttemptAt: '2026-11-12T17:03:00Z' },
      ],
    }),
    ...overrides,
  }
}

// AC-APV-14 (FR-APV-02 ง) ฝั่งหน้าจอ
//   Given บันทึกผลแล้ว แต่ Email ส่งไม่สำเร็จทุกครั้ง (หลังบ้านเปลี่ยนเป็น Failed แล้ว)
//   Then  คำขออยู่ในรายการ "Email ที่ส่งไม่สำเร็จ" พร้อมสาเหตุและเวลา
//   When  เจ้าหน้าที่กดส่งซ้ำ
//   Then  ส่งรหัส Email ไปหลังบ้าน แล้วแถวนั้นหายจากรายการ
test('AC-APV-14 แสดง Email ที่ส่งไม่สำเร็จพร้อมสาเหตุและเวลา กดส่งซ้ำแล้วหายจากรายการ', async () => {
  const calls = []
  render(<FailedEmails client={fakeClient({
    resendEmail: async (emailId) => {
      calls.push(emailId)
      return { emailId, status: 'pending' }
    },
  })} />)

  expect(await screen.findByText('Email ที่ส่งไม่สำเร็จ')).toBeTruthy()
  const table = screen.getByRole('table')
  const rows = within(table).getAllByRole('row').slice(1)
  expect(rows).toHaveLength(2)
  expect(within(rows[0]).getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
  expect(within(rows[0]).getByText('smtp: connection refused')).toBeTruthy() // สาเหตุ
  expect(rows[0].cells[4].textContent).not.toBe('') // เวลา
  expect(within(rows[1]).getByText('ปฏิเสธอัตโนมัติ')).toBeTruthy() // รวมของ System

  fireEvent.click(within(rows[0]).getByRole('button', { name: 'ส่งซ้ำ' }))

  await waitFor(() => expect(within(table).getAllByRole('row').slice(1)).toHaveLength(1))
  expect(calls).toEqual([7])
  expect(screen.queryByText('นักศึกษาตัวอย่าง ก')).toBeNull()
})

// AC-APV-14: ส่งซ้ำไม่สำเร็จ (เช่น Email ของภาควิชาอื่น) แสดงข้อความตามรหัสและรายการไม่เปลี่ยน (CON-TECH-03)
test('AC-APV-14 ส่งซ้ำแล้วหลังบ้านตอบ error แสดงข้อความตามรหัส', async () => {
  render(<FailedEmails client={fakeClient({
    resendEmail: async () => Promise.reject(Object.assign(new Error('403'), { code: 'FORBIDDEN_DEPARTMENT' })),
  })} />)

  const table = await screen.findByRole('table')
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'ส่งซ้ำ' }))

  expect((await screen.findByRole('alert')).textContent).toBe(APPROVAL_MESSAGES.FORBIDDEN_DEPARTMENT)
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
})
