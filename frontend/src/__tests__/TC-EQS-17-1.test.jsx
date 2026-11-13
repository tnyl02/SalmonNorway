// TC-EQS-17-1 (AC-EQS-17 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen, within } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'
import { queueClient, REQ_A, REQ_B } from './helpers/eqsQueueClient.js'

test('TC-EQS-17-1 เจ้าหน้าที่เห็นคำขอ 2 รายการ ลำดับ 1 = A ลำดับ 2 = B พร้อมเวลาที่ส่งและชื่อผู้ขอ', async () => {
  // Given CS-NB-0001 สถานะ "ว่าง" มีคำขอค้าง 2 รายการ ผู้ขอ A ส่งก่อน ผู้ขอ B ส่งตาม
  const client = queueClient({ requests: [REQ_A, REQ_B] })

  // When เจ้าหน้าที่เปิดรายละเอียด CS-NB-0001
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)

  // Then (1) แสดงคำขอ 2 รายการ
  const queue = await screen.findByRole('region', { name: 'คำขอยืมที่ค้างอยู่' })
  const rows = within(queue).getAllByRole('row').slice(1) // แถวแรกคือหัวตาราง
  expect(rows).toHaveLength(2)
  // Then (2) ลำดับที่ 1 คือ A ลำดับที่ 2 คือ B
  expect(within(rows[0]).getByText('1')).toBeTruthy()
  expect(within(rows[1]).getByText('2')).toBeTruthy()
  // Then (3) เวลาที่ส่งของแต่ละรายการเป็นวันที่และเวลา
  expect(within(rows[0]).getByText('2569-11-10 10:30')).toBeTruthy()
  expect(within(rows[1]).getByText('2569-11-10 11:15')).toBeTruthy()
  // Then (4) ชื่อผู้ขอ A และ B
  expect(within(rows[0]).getByText('ผู้ขอ A')).toBeTruthy()
  expect(within(rows[1]).getByText('ผู้ขอ B')).toBeTruthy()
})
