// TC-EQS-05-2 (AC-EQS-05 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

test('TC-EQS-05-2 รับคืนแล้วรีเฟรช แสดง "ว่าง" และมีปุ่มส่งคำขอยืม', async () => {
  // Given client จำลองคืน CS-NB-0004 "ถูกยืม" ครั้งแรก แล้วคืน "ว่าง" หลังรับคืน
  let calls = 0
  const client = {
    getItem: async () =>
      ++calls === 1
        ? { assetCode: 'CS-NB-0004', status: 'borrowed', dueAt: '2026-11-13T09:00:00Z', canRequestLoan: false }
        : { assetCode: 'CS-NB-0004', status: 'available', canRequestLoan: true },
    listItemRequests: async () => [],
  }

  // When นักศึกษาเปิดรายละเอียด แล้วรีเฟรช (เปิดหน้าใหม่)
  const first = render(<EquipmentDetail assetCode="CS-NB-0004" client={client} />)
  expect(await screen.findByText('ถูกยืม')).toBeTruthy()
  first.unmount()
  render(<EquipmentDetail assetCode="CS-NB-0004" client={client} />)

  // Then (1) หลังรีเฟรชแสดงป้าย "ว่าง" ไม่มีป้าย "ถูกยืม"
  expect(await screen.findByText('ว่าง')).toBeTruthy()
  expect(screen.queryByText('ถูกยืม')).toBeNull()
  // Then (2) มีปุ่ม "ส่งคำขอยืม"
  expect(screen.getByRole('button', { name: 'ส่งคำขอยืม' })).toBeTruthy()
  // Then (3) หน้าจอรีเฟรชเองหรือผู้ใช้กดเอง: ยังไม่ตรวจ รอ Q-18
})
