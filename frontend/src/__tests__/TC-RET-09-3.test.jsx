// TC-RET-09-3 (AC-RET-09 ขอบ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import LoanDetailPage from '../pages/Return/LoanDetailPage.jsx'

test('TC-RET-09-3: เวลาคืน 07:30Z แสดงเป็น 14:30 ตามเขตเวลา Asia/Bangkok', async () => {
  // Given การยืมที่คืนแล้ว เวลาคืน 2026-10-10T07:30:00Z
  const client = {
    returns: {
      loanDetail: vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        data: {
          loanId: 9,
          assetCode: 'TC-RET-09-3',
          borrowerId: 'student-001',
          dueAt: '2026-10-10T09:00:00Z',
          returnedAt: '2026-10-10T07:30:00Z',
          condition: 'normal',
          returnedLate: false,
        },
      }),
    },
  }

  // When เจ้าหน้าที่เปิดหน้ารายละเอียดการยืม
  render(<LoanDetailPage loanId={9} client={client} />)
  await screen.findByText('TC-RET-09-3')

  // Then แสดงเวลาคืนเป็น 14:30 (Asia/Bangkok)
  expect(document.body.textContent).toContain('14:30')
  expect(document.body.textContent).not.toContain('07:30')
})
