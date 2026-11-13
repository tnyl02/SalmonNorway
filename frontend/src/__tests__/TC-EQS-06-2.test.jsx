// TC-EQS-06-2 (AC-EQS-06 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

test('TC-EQS-06-2 ชิ้นที่ไม่พร้อมใช้งาน แสดงป้าย "ไม่พร้อมใช้งาน" ไม่มีปุ่มส่งคำขอยืม', async () => {
  // Given client จำลองคืน CS-PJ-0002 สถานะ "ไม่พร้อมใช้งาน"
  const client = {
    getItem: async () => ({ assetCode: 'CS-PJ-0002', status: 'unavailable', canRequestLoan: false }),
    listItemRequests: async () => [],
  }

  // When นักศึกษาเปิดหน้ารายละเอียด
  render(<EquipmentDetail assetCode="CS-PJ-0002" client={client} />)

  // Then (1) แสดงป้าย "ไม่พร้อมใช้งาน"
  expect(await screen.findByText('ไม่พร้อมใช้งาน')).toBeTruthy()
  // Then (2) ไม่มีปุ่ม "ส่งคำขอยืม"
  expect(screen.queryByRole('button', { name: 'ส่งคำขอยืม' })).toBeNull()
})
