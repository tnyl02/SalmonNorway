// TC-RET-09-4 (AC-RET-09 ทางปกติ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import ReturnPage from '../pages/Return/ReturnPage.jsx'

test('TC-RET-09-4: ฟอร์มไม่มีช่องเวลาคืน และส่งไปหลังบ้านแค่รหัสครุภัณฑ์กับสภาพ', async () => {
  // Given หน้าบันทึกการรับคืน ค้นเจอการยืมแล้ว
  const client = {
    returns: {
      lookup: vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        data: { loanId: 1, assetCode: 'CS-NB-0004', modelName: 'Notebook ตัวอย่างรุ่น A', borrowerId: 'student-001', loanType: 'normal', dueAt: '2026-11-11T09:00:00Z', isOverdue: false },
      }),
      record: vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        data: { loanId: 1, assetCode: 'CS-NB-0004', condition: 'normal', returnedLate: false, itemStatus: 'available' },
      }),
    },
  }
  render(<ReturnPage client={client} />)
  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0004' } })
  fireEvent.click(screen.getByRole('button', { name: 'ค้นหา' }))
  await screen.findByText('Notebook ตัวอย่างรุ่น A')

  // Then (1) ไม่มีช่องกรอกวันที่หรือเวลาคืน
  expect(document.querySelector('input[type="date"], input[type="time"], input[type="datetime-local"]')).toBeNull()

  // When เลือก "ปกติ" กดบันทึก
  fireEvent.click(screen.getByLabelText('ปกติ'))
  fireEvent.click(screen.getByRole('button', { name: 'บันทึกการรับคืน' }))

  // Then (2) ข้อมูลที่ส่งไปหลังบ้านมีแค่รหัสครุภัณฑ์และสภาพ ไม่มีเวลา
  await waitFor(() => expect(client.returns.record).toHaveBeenCalledTimes(1))
  const sent = client.returns.record.mock.calls[0][0]
  expect(Object.keys(sent).sort()).toEqual(['assetCode', 'condition'])
})
