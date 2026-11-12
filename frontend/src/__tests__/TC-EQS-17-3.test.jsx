// TC-EQS-17-3 (AC-EQS-17 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { act, render, screen } from '@testing-library/react'
import { ApiError } from '../api/client.js'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'
import { queueClient } from './helpers/eqsQueueClient.js'

test('TC-EQS-17-3 นักศึกษา (403 FORBIDDEN_ROLE) ไม่เห็นรายการคำขอ สถานะและปุ่มยังแสดง', async () => {
  // Given client จำลองตอบ 403 FORBIDDEN_ROLE เมื่อขอรายการคำขอ (นักศึกษา)
  const client = queueClient({ requests: new ApiError(403, 'FORBIDDEN_ROLE', 'เฉพาะเจ้าหน้าที่เท่านั้นที่ดูข้อมูลนี้ได้') })

  // When นักศึกษาเปิดรายละเอียด CS-NB-0001
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)
  await screen.findByText('ว่าง')
  await act(async () => {})
  expect(client.listItemRequests).toHaveBeenCalledWith('CS-NB-0001')

  // Then (1) ไม่มีส่วนรายการคำขอ
  expect(screen.queryByRole('region', { name: 'คำขอยืมที่ค้างอยู่' })).toBeNull()
  // Then (2) ไม่มีชื่อผู้ขอ A หรือ B
  expect(screen.queryByText('ผู้ขอ A')).toBeNull()
  expect(screen.queryByText('ผู้ขอ B')).toBeNull()
  // Then (3) ไม่มีข้อความ error
  expect(screen.queryByRole('alert')).toBeNull()
  // Then (4) สถานะ "ว่าง" และปุ่มส่งคำขอยืมยังแสดง
  expect(screen.getByText('ว่าง')).toBeTruthy()
  expect(screen.getByRole('button', { name: 'ส่งคำขอยืม' })).toBeTruthy()
})
