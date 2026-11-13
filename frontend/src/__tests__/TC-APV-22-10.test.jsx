// TC-APV-22-10 (AC-APV-22 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'

const LOADED_AT = '2026-11-12T03:00:00Z'

test('TC-APV-22-10 พิมพ์ชื่ออาจารย์แล้วกดอนุมัติคำขอค้างคืน ส่ง advisorName กับ loadedAt และแถวหายจากรายการ', async () => {
  // Given หน้ารายการได้คำขอค้างคืน C 1 แถว (loanType = overnight) จาก client จำลอง
  const calls = []
  const client = {
    getApprovalRequests: async () => ({
      loadedAt: LOADED_AT,
      items: [
        { requestId: 31, requesterName: 'นักศึกษาตัวอย่าง ค', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)', assetCode: 'AP-TEST-0001',
          purpose: 'activity', startAt: '2026-11-12T02:00:00Z', endAt: '2026-11-13T09:00:00Z',
          submittedAt: '2026-11-12T01:50:00Z', loanType: 'overnight', isFirst: false,
          attachmentUrl: '/v1/loan-requests/31/attachment' },
      ],
    }),
    approveRequest: async (requestId, loadedAt, advisorName) => {
      calls.push({ requestId, loadedAt, advisorName })
      return { requestId, status: 'approved', decidedAt: '2026-11-12T03:05:00Z', rejectedRequestIds: [] }
    },
    openAttachment: async () => {},
  }
  render(<ApproveRequests client={client} />)
  const row = within(await screen.findByRole('table')).getAllByRole('row')[1]

  // Then (1) แถว C มีช่องชื่ออาจารย์และปุ่มอนุมัติ
  const input = within(row).getByLabelText('ชื่ออาจารย์ที่เซ็นอนุมัติ')
  const button = within(row).getByRole('button', { name: 'อนุมัติ' })

  // When พิมพ์ "อ.สมชาย ใจดี" แล้วกด "อนุมัติ"
  fireEvent.change(input, { target: { value: 'อ.สมชาย ใจดี' } })
  fireEvent.click(button)

  // Then (2) ส่งรหัสคำขอ loadedAt และ advisorName ไปหลังบ้าน
  await waitFor(() => expect(calls).toHaveLength(1))
  expect(calls[0]).toEqual({ requestId: 31, loadedAt: LOADED_AT, advisorName: 'อ.สมชาย ใจดี' })

  // Then (3) แถว C หายจากรายการ
  await waitFor(() => expect(screen.queryByText('นักศึกษาตัวอย่าง ค')).toBeNull())
})
