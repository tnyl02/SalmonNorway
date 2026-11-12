// TC-APV-01-2 (AC-APV-01 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'

const LOADED_AT = '2026-11-12T03:00:00Z'

test('TC-APV-01-2 กดอนุมัติคำขอแรกแล้วส่ง requestId กับ loadedAt และแถวนั้นหายจากรายการ', async () => {
  // Given client จำลองคืนรายการ 2 คำขอ คำขอแรก isFirst: true และ client ตอบ 200 เมื่ออนุมัติ
  const calls = []
  const client = {
    getApprovalRequests: async () => ({
      loadedAt: LOADED_AT,
      items: [
        { requestId: 11, requesterName: 'นักศึกษาตัวอย่าง ก', studentCode: 'student-001', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)',
          assetCode: 'AP-TEST-0002', purpose: 'study', startAt: '2026-11-12T02:00:00Z', endAt: '2026-11-12T09:00:00Z',
          submittedAt: '2026-11-12T01:30:00Z', isFirst: true },
        { requestId: 12, requesterName: 'นักศึกษาตัวอย่าง ข', studentCode: 'student-002', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)',
          assetCode: 'AP-TEST-0003', purpose: 'study', startAt: '2026-11-12T02:00:00Z', endAt: '2026-11-12T09:00:00Z',
          submittedAt: '2026-11-12T01:40:00Z', isFirst: true },
      ],
    }),
    approveRequest: async (requestId, loadedAt) => {
      calls.push({ requestId, loadedAt })
      return { requestId, status: 'approved', decidedAt: '2026-11-12T03:05:00Z', rejectedRequestIds: [] }
    },
  }
  render(<ApproveRequests client={client} />)
  const table = await screen.findByRole('table')

  // When เปิดหน้ารายการคำขอ แล้วกด "อนุมัติ" ที่คำขอแรก
  const first = within(table).getAllByRole('row')[1]
  fireEvent.click(within(first).getByRole('button', { name: 'อนุมัติ' }))

  // Then (1) client ถูกเรียกอนุมัติ 1 ครั้ง ด้วย requestId ของแถวนั้น และ loadedAt ที่ได้ตอนโหลดรายการ
  await waitFor(() => expect(calls).toHaveLength(1))
  expect(calls[0]).toEqual({ requestId: 11, loadedAt: LOADED_AT })

  // Then (2) แถวนั้นหายจากรายการ
  await waitFor(() => expect(screen.queryByText('นักศึกษาตัวอย่าง ก')).toBeNull())

  // Then (3) อีกแถวยังอยู่
  const rows = within(table).getAllByRole('row').slice(1)
  expect(rows).toHaveLength(1)
  expect(within(rows[0]).getByText('นักศึกษาตัวอย่าง ข')).toBeTruthy()
})
