// TC-EQS-14-5 (AC-EQS-14 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

test('TC-EQS-14-5 เจ้าหน้าที่เห็นทั้งกำหนดคืน 2569-10-10 16:00 และป้าย "เลยกำหนดคืน"', async () => {
  // Given client จำลองคืนแบบที่เจ้าหน้าที่ได้: "ถูกยืม" + เลยกำหนดคืน + กำหนดคืน 2569-10-10 16:00 (= 2026-10-10 09:00 UTC)
  const client = {
    getItem: async () => ({
      assetCode: 'CS-NB-0005',
      status: 'borrowed',
      dueAt: '2026-10-10T09:00:00Z',
      isOverdue: true,
      canRequestLoan: false,
    }),
    listItemRequests: async () => [],
  }

  // When เจ้าหน้าที่เปิดหน้ารายละเอียด
  render(<EquipmentDetail assetCode="CS-NB-0005" client={client} />)

  // Then (1) แสดงกำหนดคืน "2569-10-10 16:00"
  expect(await screen.findByText('2569-10-10 16:00')).toBeTruthy()
  // Then (2) แสดงป้าย "เลยกำหนดคืน"
  expect(screen.getByText('เลยกำหนดคืน')).toBeTruthy()
})
