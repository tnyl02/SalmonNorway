// TC-APV-03-2 (AC-APV-03 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { REQ_1, fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-03-2 อนุมัติ A แล้ว A และ B ที่ถูกปฏิเสธตามหายจากรายการ คำขอของชิ้นอื่นยังอยู่', async () => {
  // Given client จำลองคืน 3 คำขอ: A (AP-TEST-0001 คำขอแรก), B (AP-TEST-0001 ไม่ใช่คำขอแรก) และคำขอของ AP-TEST-0002
  //       เมื่ออนุมัติ A ตอบ rejectedRequestIds เป็นรหัสของ B
  const A = { ...REQ_1, requestId: 31, requesterName: 'ผู้ขอ A', assetCode: 'AP-TEST-0001', isFirst: true }
  const B = { ...REQ_1, requestId: 32, requesterName: 'ผู้ขอ B', assetCode: 'AP-TEST-0001', isFirst: false }
  const other = { ...REQ_1, requestId: 33, requesterName: 'ผู้ขอชิ้นอื่น', assetCode: 'AP-TEST-0002', isFirst: true }
  const client = {
    ...fakeClient([A, B, other]),
    approveRequest: async (requestId) => ({ requestId, status: 'approved', decidedAt: '2026-11-12T03:05:00Z', rejectedRequestIds: [32] }),
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When กด "อนุมัติ" ที่ A
  const rowA = within(table).getAllByRole('row')[1]
  fireEvent.click(within(rowA).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) A และ B หายจากรายการ
  await waitFor(() => expect(screen.queryByText('ผู้ขอ A')).toBeNull())
  expect(screen.queryByText('ผู้ขอ B')).toBeNull()

  // Then (2) คำขอของ AP-TEST-0002 ยังอยู่
  expect(screen.getByText('ผู้ขอชิ้นอื่น')).toBeTruthy()
  expect(within(table).getAllByRole('row').slice(1)).toHaveLength(1)
})
