// TC-RET-04-3 (AC-RET-04 ทางปกติ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen, fireEvent } from '@testing-library/react'
import ReturnPage from '../pages/Return/ReturnPage.jsx'

test('TC-RET-04-3: ตัวเลือกสภาพมีแค่ ปกติ และ ชำรุด ไม่มี สูญหาย', async () => {
  // Given หน้าบันทึกการรับคืน ค้นเจอการยืมแล้ว
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

  // When ดูตัวเลือกสภาพ
  const radios = screen.getAllByRole('radio')

  // Then (1) มีตัวเลือก 2 ตัวคือ "ปกติ" และ "ชำรุด"
  expect(radios).toHaveLength(2)
  expect(screen.getByLabelText('ปกติ')).toBeTruthy()
  expect(screen.getByLabelText('ชำรุด')).toBeTruthy()
  // (2) ไม่มีตัวเลือก "สูญหาย"
  expect(screen.queryByLabelText('สูญหาย')).toBeNull()
  expect(screen.queryByText('สูญหาย')).toBeNull()
})
