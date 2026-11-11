// 006 T-15: หน้า OpenLoansPage ด้วย client จำลอง (AC-RET-10 ส่วนหน้าจอ)
import { render, screen, fireEvent, within } from '@testing-library/react'
import OpenLoansPage from '../OpenLoansPage.jsx'

const loans = [
  { loanId: 11, assetCode: 'CS-NB-0005', modelName: 'Notebook ตัวอย่างรุ่น A', borrowerId: 'student-002', loanType: 'normal', dueAt: '2026-11-10T09:00:00Z', isOverdue: true },
  { loanId: 12, assetCode: 'CS-NB-0004', modelName: 'Notebook ตัวอย่างรุ่น A', borrowerId: 'student-001', loanType: 'normal', dueAt: '2026-11-11T09:00:00Z', isOverdue: false },
  { loanId: 13, assetCode: 'CS-PJ-0003', modelName: 'Projector ตัวอย่างรุ่น B', borrowerId: 'student-001', loanType: 'normal', dueAt: '2026-11-12T09:00:00Z', isOverdue: false },
]

function clientWith(res) {
  return { returns: { openLoans: vi.fn().mockResolvedValue(res) } }
}

// AC-RET-10 (FR-RET-01 หน้ารายการ) บนหน้าจอ
//   Given ภาควิชาของเจ้าหน้าที่มีการยืมที่ยังไม่คืน 3 รายการ (หลังบ้านกรองภาควิชาแล้ว)
//   When  เจ้าหน้าที่เปิดหน้ารายการที่ยังไม่คืน
//   Then  เห็น 3 รายการ และเลือกรายการเพื่อบันทึกการรับคืนได้
test('AC-RET-10: แสดง 3 แถว และกดเลือกแล้วส่งรหัสครุภัณฑ์ไปหน้าบันทึกคืน', async () => {
  const onSelect = vi.fn()
  render(<OpenLoansPage client={clientWith({ ok: true, status: 200, data: loans })} onSelect={onSelect} />)

  const table = await screen.findByRole('table')
  expect(within(table).getAllByRole('row')).toHaveLength(1 + 3) // หัวตาราง + 3 แถว
  expect(screen.getByText('CS-PJ-0003')).toBeTruthy()

  fireEvent.click(screen.getByRole('button', { name: 'เลือก CS-NB-0004' }))
  expect(onSelect).toHaveBeenCalledWith('CS-NB-0004')
})

test('T-15: ไม่มีการยืมค้าง แสดงว่าไม่มี', async () => {
  render(<OpenLoansPage client={clientWith({ ok: true, status: 200, data: [] })} />)
  expect(await screen.findByText('ไม่มีการยืมที่ยังไม่คืน')).toBeTruthy()
})

test('T-15: โหลดไม่สำเร็จ แสดงข้อความพร้อม code', async () => {
  render(<OpenLoansPage client={clientWith({ ok: false, status: 0, data: { code: 'UPSTREAM_UNAVAILABLE' } })} />)
  expect((await screen.findByRole('alert')).textContent).toBe('โหลดรายการที่ยังไม่คืนไม่สำเร็จ (UPSTREAM_UNAVAILABLE)')
})

// AC-RET-13 (FR-RET-06) หน้ารายการที่ยังไม่คืน: แถวที่เลยกำหนดแสดงป้าย "เลยกำหนดคืน" เฉพาะแถวนั้น
test('AC-RET-13: แถวที่ isOverdue แสดง "เลยกำหนดคืน" แถวอื่นไม่แสดง', async () => {
  render(<OpenLoansPage client={clientWith({ ok: true, status: 200, data: loans })} />)
  const table = await screen.findByRole('table')
  const rows = within(table).getAllByRole('row').slice(1)
  expect(within(rows[0]).getByText('เลยกำหนดคืน')).toBeTruthy() // CS-NB-0005 isOverdue
  expect(within(rows[1]).queryByText('เลยกำหนดคืน')).toBeNull()
  expect(within(rows[2]).queryByText('เลยกำหนดคืน')).toBeNull()
})
