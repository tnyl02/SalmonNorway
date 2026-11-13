// TC-EQS-14-4 (AC-EQS-14 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

test('TC-EQS-14-4 นักศึกษาเห็น "ถูกยืม" กับป้าย "เลยกำหนดคืน" ไม่เห็นกำหนดคืน ไม่มีปุ่ม', async () => {
  // Given client จำลองคืนแบบที่นักศึกษาได้: "ถูกยืม" + เลยกำหนดคืน ไม่มีกำหนดคืน
  const client = {
    getItem: async () => ({ assetCode: 'CS-NB-0005', status: 'borrowed', isOverdue: true, canRequestLoan: false }),
    listItemRequests: async () => [],
  }

  // When นักศึกษาเปิดหน้ารายละเอียด
  render(<EquipmentDetail assetCode="CS-NB-0005" client={client} />)

  // Then (1) แสดงป้าย "ถูกยืม"
  expect(await screen.findByText('ถูกยืม')).toBeTruthy()
  // Then (2) แสดงป้าย "เลยกำหนดคืน"
  expect(screen.getByText('เลยกำหนดคืน')).toBeTruthy()
  // Then (3) ไม่มี "กำหนดคืน" และไม่มี "2569-10-10 16:00"
  expect(screen.queryByText('กำหนดคืน')).toBeNull()
  expect(screen.queryByText('2569-10-10 16:00')).toBeNull()
  // Then (4) ไม่มีปุ่ม "ส่งคำขอยืม"
  expect(screen.queryByRole('button', { name: 'ส่งคำขอยืม' })).toBeNull()
})
