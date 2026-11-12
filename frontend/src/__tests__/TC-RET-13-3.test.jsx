// TC-RET-13-3 (AC-RET-13 ทางปกติ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen, within } from '@testing-library/react'
import OpenLoansPage from '../pages/Return/OpenLoansPage.jsx'

test('TC-RET-13-3: แถวที่ isOverdue แสดง "เลยกำหนดคืน" แถวอื่นไม่แสดง', async () => {
  // Given หน้ารายการที่ยังไม่คืน หลังบ้านตอบ 2 รายการ รายการที่สอง isOverdue = true
  const loans = [
    { loanId: 41, assetCode: 'TC13-A', modelName: 'รุ่น A', borrowerId: 'student-001', loanType: 'normal', dueAt: '2026-11-12T09:00:00Z', isOverdue: false },
    { loanId: 42, assetCode: 'TC13-B', modelName: 'รุ่น A', borrowerId: 'student-002', loanType: 'normal', dueAt: '2026-11-10T09:00:00Z', isOverdue: true },
  ]
  const client = { returns: { openLoans: vi.fn().mockResolvedValue({ ok: true, status: 200, data: loans }) } }

  // When เจ้าหน้าที่เปิดหน้ารายการ
  render(<OpenLoansPage client={client} onSelect={vi.fn()} />)
  const rows = within(await screen.findByRole('table')).getAllByRole('row').slice(1)

  // Then (1) แถวที่สองแสดง "เลยกำหนดคืน"
  expect(within(rows[1]).getByText('เลยกำหนดคืน')).toBeTruthy()
  // (2) แถวแรกไม่แสดง
  expect(within(rows[0]).queryByText('เลยกำหนดคืน')).toBeNull()
})
