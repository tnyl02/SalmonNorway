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
//   When  เจ้าหน้าที่กรอกรหัส เห็นการยืม เลือก "ปกติ" กดบันทึก
//   Then  เห็น "บันทึกการรับคืนแล้ว สถานะ: ว่าง"
test('AC-RET-01: กรอกรหัส เลือกสภาพปกติ บันทึก แล้วเห็นสถานะว่าง', async () => {
  const client = fakeClient()
  render(<ReturnPage client={client} />)

  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0004' } })
  fireEvent.click(screen.getByRole('button', { name: 'ค้นหา' }))

  await screen.findByText('Notebook ตัวอย่างรุ่น A')
  expect(client.returns.lookup).toHaveBeenCalledWith('CS-NB-0004')
  expect(screen.getByText('student-001')).toBeTruthy()

  fireEvent.click(screen.getByLabelText('ปกติ'))
  fireEvent.click(screen.getByRole('button', { name: 'บันทึกการรับคืน' }))

  await waitFor(() => expect(screen.getByRole('status').textContent).toBe('บันทึกการรับคืนแล้ว สถานะ: ว่าง'))
  expect(client.returns.record).toHaveBeenCalledWith({ assetCode: 'CS-NB-0004', condition: 'normal' })
})

// FR-RET-02: ตัวเลือกสภาพมาจาก config.js หน้าจอไม่มีช่องกรอกเวลาคืน (ASM-04) และไม่มีช่องแนบรูป (FR-RET-05 ยกเลิก)
test('ตัวเลือกสภาพมีแค่ ปกติ และ ชำรุด ไม่มีช่องเวลาคืน และไม่มีช่องแนบรูป', async () => {
  render(<ReturnPage client={fakeClient()} />)
  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0004' } })
  fireEvent.click(screen.getByRole('button', { name: 'ค้นหา' }))
  await screen.findByText('Notebook ตัวอย่างรุ่น A')

  expect(screen.getAllByRole('radio').map((r) => r.value)).toEqual(['normal', 'damaged'])
  expect(document.querySelector('input[type="datetime-local"], input[type="time"], input[type="date"]')).toBeNull()
  expect(document.querySelector('input[type="file"]')).toBeNull()
})

// เปิดหน้าแล้วหาการยืมให้เจอ (ขั้นก่อนกรอกสภาพ)
async function openLoanForm(client) {
  render(<ReturnPage client={client} />)
  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0004' } })
  fireEvent.click(screen.getByRole('button', { name: 'ค้นหา' }))
  await screen.findByText('Notebook ตัวอย่างรุ่น A')
}

// AC-RET-04 (FR-RET-02) บนหน้าจอ
//   Given เจ้าหน้าที่กำลังบันทึกการรับคืน
//   When  ไม่ได้ระบุสภาพ แล้วกดบันทึก
//   Then  ระบบไม่บันทึก แจ้งให้เลือกสภาพ ปกติ หรือ ชำรุด และไม่มีตัวเลือก "สูญหาย"
test('AC-RET-04: กดบันทึกโดยไม่เลือกสภาพ แจ้งให้เลือก ปกติ หรือ ชำรุด และไม่ส่งไปหลังบ้าน', async () => {
  const client = fakeClient()
  await openLoanForm(client)

  expect(screen.getAllByRole('radio').map((r) => r.value)).toEqual(['normal', 'damaged'])
  expect(screen.queryByLabelText('สูญหาย')).toBeNull()

  fireEvent.click(screen.getByRole('button', { name: 'บันทึกการรับคืน' }))

  expect(screen.getByRole('alert').textContent).toBe('กรุณาเลือกสภาพ ปกติ หรือ ชำรุด')
  expect(client.returns.record).not.toHaveBeenCalled()
  expect(screen.queryByRole('status')).toBeNull()
})

// หลังบ้านตอบ CONDITION_REQUIRED (เช่นส่งค่าที่ไม่รู้จัก) แสดงข้อความเดียวกัน (FR-RET-02)
test('AC-RET-04: หลังบ้านตอบ CONDITION_REQUIRED แสดงข้อความให้เลือกสภาพ', async () => {
  const client = fakeClient()
  client.returns.record.mockResolvedValue({ ok: false, status: 422, data: { code: 'CONDITION_REQUIRED' } })
  await openLoanForm(client)

  fireEvent.click(screen.getByLabelText('ชำรุด'))
  fireEvent.click(screen.getByRole('button', { name: 'บันทึกการรับคืน' }))

  await waitFor(() => expect(screen.getByRole('alert').textContent).toBe('กรุณาเลือกสภาพ ปกติ หรือ ชำรุด'))
})

// T-12: หลังบันทึกสำเร็จ มีลิงก์ไปรายละเอียดการยืม (FR-RET-06)
test('T-12: บันทึกสำเร็จแล้วกด "ดูรายละเอียดการยืม" เห็นรายละเอียดของการยืมนั้น', async () => {
  const client = fakeClient()
  client.returns.loanDetail = vi.fn().mockResolvedValue({
    ok: true,
    status: 200,
    data: { loanId: 1, assetCode: 'CS-NB-0004', borrowerId: 'student-001', dueAt: '2026-11-11T09:00:00Z', returnedAt: '2026-11-11T08:00:00Z', condition: 'normal', returnedLate: false },
  })
  await openLoanForm(client)
  fireEvent.click(screen.getByLabelText('ปกติ'))
  fireEvent.click(screen.getByRole('button', { name: 'บันทึกการรับคืน' }))
  await screen.findByRole('status')

  fireEvent.click(screen.getByRole('button', { name: 'ดูรายละเอียดการยืม' }))

  await screen.findByRole('region', { name: 'รายละเอียดการยืม' })
  expect(client.returns.loanDetail).toHaveBeenCalledWith(1)
})
