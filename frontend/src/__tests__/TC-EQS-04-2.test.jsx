// TC-EQS-04-2 (AC-EQS-04 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

test('TC-EQS-04-2 ชิ้นที่ถูกยืมยังไม่เลยกำหนด แสดง "ถูกยืม" กำหนดคืน และไม่มีปุ่มส่งคำขอยืม', async () => {
  // Given client จำลองคืนชิ้นสถานะ "ถูกยืม" กำหนดคืน 2569-11-13 16:00 (= 2026-11-13 09:00 UTC) ส่งคำขอยืมไม่ได้
  const client = {
    getItem: async () => ({ assetCode: 'CS-NB-0004', status: 'borrowed', dueAt: '2026-11-13T09:00:00Z', canRequestLoan: false }),
    listItemRequests: async () => [],
  }

  // When นักศึกษาเปิดหน้ารายละเอียด
  render(<EquipmentDetail assetCode="CS-NB-0004" client={client} />)

  // Then (1) แสดงป้าย "ถูกยืม"
  expect(await screen.findByText('ถูกยืม')).toBeTruthy()
  // Then (2) แสดงกำหนดคืน "2569-11-13 16:00"
  expect(screen.getByText('2569-11-13 16:00')).toBeTruthy()
  // Then (3) ไม่มีปุ่ม "ส่งคำขอยืม"
  expect(screen.queryByRole('button', { name: 'ส่งคำขอยืม' })).toBeNull()
  // Then (4) ไม่มีป้าย "เลยกำหนดคืน"
  expect(screen.queryByText('เลยกำหนดคืน')).toBeNull()
})
