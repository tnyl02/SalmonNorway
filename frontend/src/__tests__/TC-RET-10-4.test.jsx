// TC-RET-10-4 (AC-RET-10 ทางปกติ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen, fireEvent, within } from '@testing-library/react'
import OpenLoansPage from '../pages/Return/OpenLoansPage.jsx'

const loans = [
  { loanId: 21, assetCode: 'TC10-A', modelName: 'รุ่น A', borrowerId: 'student-001', loanType: 'normal', dueAt: '2026-11-10T09:00:00Z', isOverdue: false },
  { loanId: 22, assetCode: 'TC10-B', modelName: 'รุ่น A', borrowerId: 'student-002', loanType: 'normal', dueAt: '2026-11-11T09:00:00Z', isOverdue: false },
  { loanId: 23, assetCode: 'TC10-C', modelName: 'รุ่น B', borrowerId: 'student-001', loanType: 'normal', dueAt: '2026-11-12T09:00:00Z', isOverdue: false },
]

test('TC-RET-10-4: แสดง 3 แถว และกดเลือกรายการที่สองแล้วส่งรหัสครุภัณฑ์ไปบันทึกคืน', async () => {
  // Given หน้ารายการที่ยังไม่คืน หลังบ้านตอบ 3 รายการ
  const onSelect = vi.fn()
  const client = { returns: { openLoans: vi.fn().mockResolvedValue({ ok: true, status: 200, data: loans }) } }

  // When เจ้าหน้าที่ดูรายการ
  render(<OpenLoansPage client={client} onSelect={onSelect} />)
  const table = await screen.findByRole('table')

  // Then (1) แสดง 3 แถว (ไม่นับหัวตาราง)
  expect(within(table).getAllByRole('row')).toHaveLength(1 + 3)

  // When กดเลือกรายการที่สอง
  fireEvent.click(screen.getByRole('button', { name: 'เลือก TC10-B' }))

  // Then (2) ส่งรหัสครุภัณฑ์ของรายการนั้นไปบันทึกการรับคืน
  expect(onSelect).toHaveBeenCalledWith('TC10-B')
})
