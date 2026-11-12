// TC-RET-13-5 (AC-RET-13 ขอบ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import LoanDetailPage from '../pages/Return/LoanDetailPage.jsx'

test('TC-RET-13-5: หน้ารายละเอียดของการยืมที่คืนทันกำหนด ไม่แสดง "คืนเกินกำหนด"', async () => {
  // Given หน้ารายละเอียดการยืม หลังบ้านตอบการยืมที่ returnedLate = false
  const client = {
    returns: {
      loanDetail: vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        data: { loanId: 51, assetCode: 'TC-RET-13-5', borrowerId: 'student-001', dueAt: '2026-10-10T09:00:00Z', returnedAt: '2026-10-10T08:00:00Z', condition: 'normal', returnedLate: false },
      }),
    },
  }

  // When เจ้าหน้าที่เปิดหน้ารายละเอียด
  render(<LoanDetailPage loanId={51} client={client} />)
  await screen.findByText('TC-RET-13-5')

  // Then ไม่แสดง "คืนเกินกำหนด"
  expect(screen.queryByText('คืนเกินกำหนด')).toBeNull()
})
