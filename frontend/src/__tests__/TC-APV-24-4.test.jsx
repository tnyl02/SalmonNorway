// TC-APV-24-4 (AC-APV-24 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'

test('TC-APV-24-4 อนุมัติคำขอค้างคืน C แล้วแถว A B และ C หายจากรายการทั้งหมด', async () => {
  // Given หน้ารายการได้ A, B (ปกติ) และ C (ค้างคืน) ของอุปกรณ์ชิ้นเดียวกัน
  //       client จำลองตอบการอนุมัติ C ว่า rejectedRequestIds = [A, B]
  const base = {
    modelName: 'อุปกรณ์ทดลองอนุมัติ (004)', assetCode: 'AP-TEST-0001', purpose: 'study',
    startAt: '2026-11-12T02:00:00Z', endAt: '2026-11-12T09:00:00Z',
  }
  const client = {
    getApprovalRequests: async () => ({
      loadedAt: '2026-11-12T03:00:00Z',
      items: [
        { ...base, requestId: 51, requesterName: 'นักศึกษาตัวอย่าง ก', submittedAt: '2026-11-12T01:30:00Z', loanType: 'normal', isFirst: true },
        { ...base, requestId: 52, requesterName: 'นักศึกษาตัวอย่าง ข', submittedAt: '2026-11-12T01:40:00Z', loanType: 'normal', isFirst: false },
        { ...base, requestId: 53, requesterName: 'นักศึกษาตัวอย่าง ค', submittedAt: '2026-11-12T01:50:00Z', loanType: 'overnight',
          isFirst: false, endAt: '2026-11-13T09:00:00Z' },
      ],
    }),
    approveRequest: async (requestId) => ({
      requestId, status: 'approved', decidedAt: '2026-11-12T03:05:00Z', rejectedRequestIds: [51, 52],
    }),
    openAttachment: async () => {},
  }
  render(<ApproveRequests client={client} />)
  const rowC = within(await screen.findByRole('table')).getAllByRole('row')[3]

  // When พิมพ์ชื่ออาจารย์ที่แถว C แล้วกด "อนุมัติ"
  fireEvent.change(within(rowC).getByLabelText('ชื่ออาจารย์ที่เซ็นอนุมัติ'), { target: { value: 'อ.สมชาย ใจดี' } })
  fireEvent.click(within(rowC).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) แถว A, B และ C หายจากรายการทั้งหมด
  await waitFor(() => expect(screen.queryByText('นักศึกษาตัวอย่าง ค')).toBeNull())
  expect(screen.queryByText('นักศึกษาตัวอย่าง ก')).toBeNull()
  expect(screen.queryByText('นักศึกษาตัวอย่าง ข')).toBeNull()
})
