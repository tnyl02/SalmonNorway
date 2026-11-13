// TC-RET-13-2 (AC-RET-13 ทางปกติ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import LoanDetailPage from '../pages/Return/LoanDetailPage.jsx'

test('TC-RET-13-2: หน้ารายละเอียดของการยืมที่ returnedLate แสดง "คืนเกินกำหนด"', async () => {
  // Given หน้ารายละเอียดการยืม หลังบ้านตอบการยืมที่ returnedLate = true
  const client = {
    returns: {
      loanDetail: vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        data: { loanId: 31, assetCode: 'TC-RET-13-2', borrowerId: 'student-001', dueAt: '2026-10-10T09:00:00Z', returnedAt: '2026-10-10T10:00:00Z', condition: 'normal', returnedLate: true },
      }),
    },
  }

  // When เจ้าหน้าที่เปิดหน้ารายละเอียด
  render(<LoanDetailPage loanId={31} client={client} />)

  // Then แสดงว่า "คืนเกินกำหนด"
  expect(await screen.findByText('คืนเกินกำหนด')).toBeTruthy()
})
