// TC-APV-22-11 (AC-APV-22 ผิด) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'

test('TC-APV-22-11 ไม่พิมพ์ชื่ออาจารย์ หรือพิมพ์ยาว 101 ตัวอักษร แสดงข้อความและไม่เรียกหลังบ้าน', async () => {
  // Given เหมือน 22-10 หน้ารายการได้คำขอค้างคืน C 1 แถว
  const calls = []
  const client = {
    getApprovalRequests: async () => ({
      loadedAt: '2026-11-12T03:00:00Z',
      items: [
        { requestId: 31, requesterName: 'นักศึกษาตัวอย่าง ค', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)', assetCode: 'AP-TEST-0001',
          purpose: 'activity', startAt: '2026-11-12T02:00:00Z', endAt: '2026-11-13T09:00:00Z',
          submittedAt: '2026-11-12T01:50:00Z', loanType: 'overnight', isFirst: false },
      ],
    }),
    approveRequest: async (...args) => {
      calls.push(args)
      return { requestId: 31, status: 'approved', decidedAt: '2026-11-12T03:05:00Z', rejectedRequestIds: [] }
    },
    openAttachment: async () => {},
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')
  const row = within(table).getAllByRole('row')[1]
  const button = within(row).getByRole('button', { name: 'อนุมัติ' })

  // When (ก) กด "อนุมัติ" โดยไม่พิมพ์ชื่อ
  fireEvent.click(button)
  // Then (1) แสดง "กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ"
  expect(screen.getByRole('alert').textContent).toBe('กรุณาระบุชื่ออาจารย์ที่เซ็นอนุมัติ')

  // When (ข) พิมพ์ชื่อยาว 101 ตัวอักษรแล้วกด "อนุมัติ"
  fireEvent.change(within(row).getByLabelText('ชื่ออาจารย์ที่เซ็นอนุมัติ'), { target: { value: 'ก'.repeat(101) } })
  fireEvent.click(button)
  // Then (2) แสดง "ชื่ออาจารย์ต้องไม่เกิน 100 ตัวอักษร"
  expect(screen.getByRole('alert').textContent).toBe('ชื่ออาจารย์ต้องไม่เกิน 100 ตัวอักษร')

  // Then (3) ไม่เรียกหลังบ้านทั้ง 2 แบบ และแถว C ยังอยู่
  expect(calls).toHaveLength(0)
  expect(within(table).getByText('นักศึกษาตัวอย่าง ค')).toBeTruthy()
})
