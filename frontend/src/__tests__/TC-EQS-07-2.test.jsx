// TC-EQS-07-2 (AC-EQS-07 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { act, render, screen } from '@testing-library/react'
import { ApiError } from '../api/client.js'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

test('TC-EQS-07-2 นักศึกษาเปิดรายละเอียดชิ้นที่ถูกยืม เห็นป้าย "ถูกยืม" ไม่เห็นผู้ยืม', async () => {
  // Given client จำลองคืน CS-NB-0004 สถานะ "ถูกยืม" แบบที่นักศึกษาได้ (ไม่มีช่องผู้ยืม)
  const client = {
    getItem: async () => ({ assetCode: 'CS-NB-0004', status: 'borrowed', canRequestLoan: false }),
    listItemRequests: async () => {
      throw new ApiError(403, 'FORBIDDEN_ROLE', 'เฉพาะเจ้าหน้าที่เท่านั้นที่ดูข้อมูลนี้ได้')
    },
  }

  // When นักศึกษาเปิดหน้ารายละเอียด CS-NB-0004
  render(<EquipmentDetail assetCode="CS-NB-0004" client={client} />)

  // Then (1) แสดงป้าย "ถูกยืม"
  expect(await screen.findByText('ถูกยืม')).toBeTruthy()
  await act(async () => {})
  // Then (2) หน้าจอไม่มีหัวข้อหรือชื่อผู้ยืม
  expect(screen.queryByText(/ผู้ยืม/)).toBeNull()
  expect(screen.queryByText(/student-001/)).toBeNull()
})
