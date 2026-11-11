// 006 T-12: หน้า LoanDetailPage ด้วย client จำลอง (FR-RET-06) ไม่มี AC ตรง เป็นพื้นฐานของ T-17 (AC-RET-13)
import { render, screen } from '@testing-library/react'
import LoanDetailPage from '../LoanDetailPage.jsx'

function clientWith(res) {
  return { returns: { loanDetail: vi.fn().mockResolvedValue(res) } }
}

test('T-12: แสดงเวลาคืน และสภาพของการยืมที่คืนแล้ว', async () => {
  const client = clientWith({
    ok: true,
    status: 200,
    data: {
      loanId: 7,
      assetCode: 'CS-NB-0005',
      borrowerId: 'student-002',
      dueAt: '2026-11-10T09:00:00Z',
      returnedAt: '2026-11-11T03:30:00Z',
      condition: 'damaged',
      returnedLate: true,
    },
  })
  render(<LoanDetailPage loanId={7} client={client} />)

  await screen.findByText('CS-NB-0005')
  expect(client.returns.loanDetail).toHaveBeenCalledWith(7)
  expect(screen.getByText('student-002')).toBeTruthy()
  expect(screen.getByText('ชำรุด')).toBeTruthy()
  expect(screen.queryByText('ยังไม่คืน')).toBeNull()
})

test('T-12: การยืมที่ยังไม่คืน แสดงว่ายังไม่คืน', async () => {
  const client = clientWith({
    ok: true,
    status: 200,
    data: { loanId: 8, assetCode: 'CS-NB-0004', borrowerId: 'student-001', dueAt: '2026-11-11T09:00:00Z', returnedAt: null, condition: null, returnedLate: null },
  })
  render(<LoanDetailPage loanId={8} client={client} />)
  expect(await screen.findByText('ยังไม่คืน')).toBeTruthy()
})

test('T-12: โหลดไม่สำเร็จ แสดงข้อความพร้อม code', async () => {
  render(<LoanDetailPage loanId={9} client={clientWith({ ok: false, status: 404, data: { code: 'LOAN_NOT_FOUND' } })} />)
  expect((await screen.findByRole('alert')).textContent).toBe('โหลดรายละเอียดการยืมไม่สำเร็จ (LOAN_NOT_FOUND)')
})

// AC-RET-13 (FR-RET-06) หน้ารายละเอียด: การยืมที่คืนเกินกำหนดแสดงป้าย "คืนเกินกำหนด" / คืนทันไม่แสดง
test('AC-RET-13: หน้ารายละเอียดแสดง "คืนเกินกำหนด" เมื่อ returnedLate', async () => {
  const base = { assetCode: 'CS-NB-0005', borrowerId: 'student-002', dueAt: '2026-11-10T09:00:00Z', returnedAt: '2026-11-11T03:30:00Z', condition: 'normal' }
  const { unmount } = render(<LoanDetailPage loanId={21} client={clientWith({ ok: true, status: 200, data: { ...base, loanId: 21, returnedLate: true } })} />)
  expect(await screen.findByText('คืนเกินกำหนด')).toBeTruthy()
  unmount()

  render(<LoanDetailPage loanId={22} client={clientWith({ ok: true, status: 200, data: { ...base, loanId: 22, returnedLate: false } })} />)
  await screen.findByText('CS-NB-0005')
  expect(screen.queryByText('คืนเกินกำหนด')).toBeNull()
})
