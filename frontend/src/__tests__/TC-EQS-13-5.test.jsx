// TC-EQS-13-5 (AC-EQS-13 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import { ApiError } from '../api/client.js'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

const MESSAGE = 'ไม่มีสิทธิ์เข้าถึงอุปกรณ์ของภาควิชาอื่น'

test('TC-EQS-13-5 เปิดรายละเอียดของภาควิชาอื่น แจ้งไม่มีสิทธิ์ ไม่แสดงข้อมูลอุปกรณ์', async () => {
  // Given client จำลองตอบ 403 FORBIDDEN_DEPARTMENT พร้อม message
  const client = {
    getItem: async () => {
      throw new ApiError(403, 'FORBIDDEN_DEPARTMENT', MESSAGE)
    },
    listItemRequests: async () => [],
  }

  // When นักศึกษาเปิดรายละเอียด CS-NB-0001
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)

  // Then (1) แสดงข้อความแจ้งว่าไม่มีสิทธิ์เข้าถึง (ถ้อยคำจาก message ของหลังบ้าน)
  expect((await screen.findByRole('alert')).textContent).toBe(MESSAGE)
  // Then (2) ไม่แสดงสถานะหรือข้อมูลอุปกรณ์
  for (const label of ['สถานะ', 'ว่าง', 'ถูกยืม', 'ไม่พร้อมใช้งาน', 'กำหนดคืน', 'ส่งคำขอยืม']) {
    expect(screen.queryByText(label)).toBeNull()
  }
})
