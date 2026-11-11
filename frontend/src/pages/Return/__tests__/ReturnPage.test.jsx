// 006 T-06: หน้า ReturnPage ด้วย client จำลอง (ไม่ต่อหลังบ้านจริง)
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import ReturnPage from '../ReturnPage.jsx'

const openLoan = {
  loanId: 1,
  assetCode: 'CS-NB-0004',
  modelName: 'Notebook ตัวอย่างรุ่น A',
  borrowerId: 'student-001',
  loanType: 'normal',
  dueAt: '2026-11-11T09:00:00Z',
  isOverdue: false,
}

function fakeClient() {
  return {
    returns: {
      lookup: vi.fn().mockResolvedValue({ ok: true, status: 200, data: openLoan }),
      record: vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        data: { loanId: 1, assetCode: 'CS-NB-0004', condition: 'normal', returnedLate: false, itemStatus: 'available' },
      }),
    },
  }
}

// AC-RET-01 (FR-RET-03) บนหน้าจอ
//   Given อุปกรณ์ "ถูกยืม" มีการยืมที่ยังไม่คืน
//   When  เจ้าหน้าที่กรอกรหัส เห็นการยืม เลือก "ปกติ" แนบ 1 รูป กดบันทึก
//   Then  เห็น "บันทึกการรับคืนแล้ว สถานะ: ว่าง"
test('AC-RET-01: กรอกรหัส เลือกสภาพปกติ แนบรูป บันทึก แล้วเห็นสถานะว่าง', async () => {
  const client = fakeClient()
  render(<ReturnPage client={client} />)

  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0004' } })
  fireEvent.click(screen.getByRole('button', { name: 'ค้นหา' }))

  await screen.findByText('Notebook ตัวอย่างรุ่น A')
  expect(client.returns.lookup).toHaveBeenCalledWith('CS-NB-0004')
  expect(screen.getByText('student-001')).toBeTruthy()

  fireEvent.click(screen.getByLabelText('ปกติ'))
  const photo = new File(['jpeg'], 'photo.jpg', { type: 'image/jpeg' })
  fireEvent.change(screen.getByLabelText('รูปสภาพอุปกรณ์'), { target: { files: [photo] } })
  fireEvent.click(screen.getByRole('button', { name: 'บันทึกการรับคืน' }))

  await waitFor(() => expect(screen.getByRole('status').textContent).toBe('บันทึกการรับคืนแล้ว สถานะ: ว่าง'))
  expect(client.returns.record).toHaveBeenCalledWith({ assetCode: 'CS-NB-0004', condition: 'normal', photos: [photo] })
})

// FR-RET-02: ตัวเลือกสภาพมาจาก config.js และหน้าจอไม่มีช่องกรอกเวลาคืน (ASM-04)
test('ตัวเลือกสภาพมีแค่ ปกติ และ ชำรุด และไม่มีช่องเวลาคืน', async () => {
  render(<ReturnPage client={fakeClient()} />)
  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0004' } })
  fireEvent.click(screen.getByRole('button', { name: 'ค้นหา' }))
  await screen.findByText('Notebook ตัวอย่างรุ่น A')

  expect(screen.getAllByRole('radio').map((r) => r.value)).toEqual(['normal', 'damaged'])
  expect(document.querySelector('input[type="datetime-local"], input[type="time"], input[type="date"]')).toBeNull()
})
