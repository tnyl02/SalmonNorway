// TC-APV-07-2 (AC-APV-07 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { REQ_1, fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-07-2 หลังบ้านตอบ NOT_FIRST_IN_QUEUE แสดงข้อความ และทั้ง 2 แถวยังอยู่', async () => {
  // Given client จำลองคืน 2 คำขอที่ทั้งคู่มี isFirst: true และตอบ error NOT_FIRST_IN_QUEUE เมื่ออนุมัติแถวที่สอง
  const A = { ...REQ_1, requestId: 51, requesterName: 'ผู้ขอ A', assetCode: 'AP-TEST-0001', isFirst: true }
  const B = { ...REQ_1, requestId: 52, requesterName: 'ผู้ขอ B', assetCode: 'AP-TEST-0001', isFirst: true }
  const client = {
    ...fakeClient([A, B]),
    approveRequest: async () => Promise.reject(Object.assign(new Error('409'), { code: 'NOT_FIRST_IN_QUEUE' })),
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "อนุมัติ" ที่แถวที่สอง
  const rowB = within(table).getAllByRole('row')[2]
  fireEvent.click(within(rowB).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) แสดง "ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน"
  expect((await screen.findByRole('alert')).textContent).toBe('ต้องพิจารณาคำขอที่ส่งก่อนหน้าของอุปกรณ์ชิ้นนี้ก่อน')

  // Then (2) ทั้ง 2 แถวยังอยู่ในรายการ
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(2)
})
