// TC-APV-07-1 (AC-APV-07 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { REQ_1, fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-07-1 ปุ่มอนุมัติมีเฉพาะคำขอแรก คำขอที่ไม่ใช่คำขอแรกยังปฏิเสธได้', async () => {
  // Given client จำลองคืน A (isFirst: true) และ B (isFirst: false) ของอุปกรณ์ชิ้นเดียวกัน
  const A = { ...REQ_1, requestId: 41, requesterName: 'ผู้ขอ A', assetCode: 'AP-TEST-0001', isFirst: true }
  const B = { ...REQ_1, requestId: 42, requesterName: 'ผู้ขอ B', assetCode: 'AP-TEST-0001', isFirst: false }

  // When เปิดหน้ารายการคำขอ
  render(<ApproveRequests client={fakeClient([A, B])} />)
  const table = await screen.findByRole('table')
  const [rowA, rowB] = within(table).getAllByRole('row').slice(1)

  // Then (1) แถว A มีปุ่ม "อนุมัติ"
  expect(within(rowA).getByRole('button', { name: 'อนุมัติ' })).toBeTruthy()
  // Then (2) แถว B ไม่มีปุ่ม "อนุมัติ"
  expect(within(rowB).queryByRole('button', { name: 'อนุมัติ' })).toBeNull()
  // Then (3) แถว B ยังมีปุ่ม "ปฏิเสธ"
  expect(within(rowB).getByRole('button', { name: 'ปฏิเสธ' })).toBeTruthy()
})
