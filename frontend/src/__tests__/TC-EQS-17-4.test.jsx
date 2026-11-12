// TC-EQS-17-4 (AC-EQS-17 ขอบ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
// Then ทีมตัดสิน 2569-11-12 รับพฤติกรรมของโค้ด: ไม่มีคำขอค้าง = ซ่อนส่วนนี้
import { act, render, screen } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'
import { queueClient } from './helpers/eqsQueueClient.js'

test('TC-EQS-17-4 ไม่มีคำขอค้าง ไม่มีส่วนรายการคำขอ สถานะยังแสดง', async () => {
  // Given client จำลองคืนคำขอ 0 รายการ
  const client = queueClient({ requests: [] })

  // When เจ้าหน้าที่เปิดรายละเอียด CS-NB-0001
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)
  await screen.findByText('ว่าง')
  await act(async () => {})
  expect(client.listItemRequests).toHaveBeenCalledWith('CS-NB-0001')

  // Then (1) ไม่มีส่วนรายการคำขอ
  expect(screen.queryByRole('region', { name: 'คำขอยืมที่ค้างอยู่' })).toBeNull()
  // Then (2) สถานะ "ว่าง" ยังแสดง
  expect(screen.getByText('ว่าง')).toBeTruthy()
})
