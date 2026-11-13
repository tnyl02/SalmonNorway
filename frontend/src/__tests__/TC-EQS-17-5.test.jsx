// TC-EQS-17-5 (AC-EQS-17 ทางผิด) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
// Then (3) ทีมตัดสิน 2569-11-12 รับพฤติกรรมของโค้ด: ซ่อนเงียบ
import { act, render, screen } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'
import { queueClient } from './helpers/eqsQueueClient.js'

test('TC-EQS-17-5 endpoint ของ 002 ไม่ตอบ สถานะและกำหนดคืนยังแสดง ส่วนคำขอซ่อนเงียบ', async () => {
  // Given endpoint ลำดับคำขอของ 002 ไม่ตอบ (เชื่อมต่อไม่ได้)
  const borrowed = { assetCode: 'CS-NB-0004', status: 'borrowed', dueAt: '2026-11-10T09:00:00Z', canRequestLoan: false }
  const client = queueClient({ item: borrowed, requests: new TypeError('Failed to fetch') })

  // When เจ้าหน้าที่เปิดรายละเอียด CS-NB-0004 (ถูกยืม)
  render(<EquipmentDetail assetCode="CS-NB-0004" client={client} />)
  expect(await screen.findByText('ถูกยืม')).toBeTruthy()
  await act(async () => {})
  expect(client.listItemRequests).toHaveBeenCalledWith('CS-NB-0004')

  // Then (1) สถานะ "ถูกยืม" และกำหนดคืนยังแสดงตามปกติ
  expect(screen.getByText('ถูกยืม')).toBeTruthy()
  expect(screen.getByText('2569-11-10 16:00')).toBeTruthy()
  // Then (2) ส่วนรายการคำขอแสดงไม่ได้ ไม่มีชื่อผู้ขอ
  expect(screen.queryByText(/ผู้ขอ A|ผู้ขอ B/)).toBeNull()
  // Then (3) ไม่มีส่วนรายการคำขอและไม่มีข้อความ error
  expect(screen.queryByRole('region', { name: 'คำขอยืมที่ค้างอยู่' })).toBeNull()
  expect(screen.queryByRole('alert')).toBeNull()
})
