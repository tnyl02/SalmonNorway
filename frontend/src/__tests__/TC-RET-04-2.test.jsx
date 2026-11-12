// TC-RET-04-2 (AC-RET-04 ทางปกติ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen, fireEvent } from '@testing-library/react'
import ReturnPage from '../pages/Return/ReturnPage.jsx'

test('TC-RET-04-2: กดบันทึกโดยยังไม่เลือกสภาพ แจ้งให้เลือก ปกติ หรือ ชำรุด และไม่ส่งไปหลังบ้าน', async () => {
  // Given หน้าบันทึกการรับคืน ค้นเจอการยืมแล้ว ยังไม่ได้เลือกสภาพ
  const client = {
    returns: {
      lookup: vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        data: { loanId: 1, assetCode: 'CS-NB-0004', modelName: 'Notebook ตัวอย่างรุ่น A', borrowerId: 'student-001', loanType: 'normal', dueAt: '2026-11-11T09:00:00Z', isOverdue: false },
      }),
      record: vi.fn(),
    },
  }
  render(<ReturnPage client={client} />)
  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0004' } })
  fireEvent.click(screen.getByRole('button', { name: 'ค้นหา' }))
  await screen.findByText('Notebook ตัวอย่างรุ่น A')

  // When กดบันทึก
  fireEvent.click(screen.getByRole('button', { name: 'บันทึกการรับคืน' }))

  // Then (1) แจ้งให้เลือกสภาพ ข้อความมีคำว่า "ปกติ" และ "ชำรุด"
  const alert = screen.getByRole('alert').textContent
  expect(alert).toContain('ปกติ')
  expect(alert).toContain('ชำรุด')
  // (2) ไม่ส่งคำขอบันทึกไปหลังบ้าน
  expect(client.returns.record).not.toHaveBeenCalled()
})
