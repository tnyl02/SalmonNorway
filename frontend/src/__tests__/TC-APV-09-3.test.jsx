// TC-APV-09-3 (AC-APV-09 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { REQ_1, REQ_2, fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-09-3 แสดงข้อความยกเลิกแล้วและพิจารณาแล้วตามรหัส', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ ตอบ REQUEST_CANCELLED เมื่ออนุมัติแถวแรก และ REQUEST_ALREADY_DECIDED เมื่ออนุมัติแถวที่สอง
  const codes = { [REQ_1.requestId]: 'REQUEST_CANCELLED', [REQ_2.requestId]: 'REQUEST_ALREADY_DECIDED' }
  const client = {
    ...fakeClient(),
    approveRequest: async (requestId) => Promise.reject(Object.assign(new Error('409'), { code: codes[requestId] })),
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')
  const [row1, row2] = within(table).getAllByRole('row').slice(1)

  // When กด "อนุมัติ" ที่แถวแรก
  fireEvent.click(within(row1).getByRole('button', { name: 'อนุมัติ' }))
  // Then (1) ครั้งแรกแสดง "คำขอนี้ถูกผู้ขอยกเลิกแล้ว"
  await waitFor(() => expect(screen.getByRole('alert').textContent).toBe('คำขอนี้ถูกผู้ขอยกเลิกแล้ว'))

  // When กด "อนุมัติ" ที่แถวที่สอง
  fireEvent.click(within(row2).getByRole('button', { name: 'อนุมัติ' }))
  // Then (2) ครั้งที่สองแสดง "คำขอนี้ได้รับการพิจารณาแล้ว"
  await waitFor(() => expect(screen.getByRole('alert').textContent).toBe('คำขอนี้ได้รับการพิจารณาแล้ว'))
})
