// TC-APV-23-4 (AC-APV-23 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'

test('TC-APV-23-4 คำขอค้างคืน 2 แถวของอุปกรณ์ชิ้นเดียวกันมีปุ่มอนุมัติและช่องชื่ออาจารย์ทั้งคู่', async () => {
  // Given หน้ารายการได้คำขอค้างคืน 2 แถวของอุปกรณ์ชิ้นเดียวกัน isFirst = false ทั้งคู่ จาก client จำลอง
  const overnight = (requestId, requesterName, submittedAt) => ({
    requestId, requesterName, modelName: 'อุปกรณ์ทดลองอนุมัติ (004)', assetCode: 'AP-TEST-0003', purpose: 'activity',
    startAt: '2026-11-12T02:00:00Z', endAt: '2026-11-13T09:00:00Z', submittedAt, loanType: 'overnight', isFirst: false,
  })
  const client = {
    getApprovalRequests: async () => ({
      loadedAt: '2026-11-12T03:00:00Z',
      items: [overnight(41, 'ค้างคืน X', '2026-11-12T00:00:00Z'), overnight(42, 'ค้างคืน Y', '2026-11-12T02:59:00Z')],
    }),
    approveRequest: async () => ({}),
    openAttachment: async () => {},
  }

  // When เจ้าหน้าที่เปิดหน้ารายการ
  render(<ApproveRequests client={client} />)
  const rows = within(await screen.findByRole('table')).getAllByRole('row').slice(1)

  // Then (1) ทั้ง 2 แถวมีปุ่ม "อนุมัติ" และช่องชื่ออาจารย์
  expect(rows).toHaveLength(2)
  for (const row of rows) {
    expect(within(row).getByRole('button', { name: 'อนุมัติ' })).toBeTruthy()
    expect(within(row).getByLabelText('ชื่ออาจารย์ที่เซ็นอนุมัติ')).toBeTruthy()
  }
})
