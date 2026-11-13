// TC-EQS-17-2 (AC-EQS-17 ขอบ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen, within } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'
import { queueClient, REQ_A, REQ_B } from './helpers/eqsQueueClient.js'

test('TC-EQS-17-2 client คืน B ก่อน A หน้าจอยังเรียง A ลำดับ 1 ก่อน B ลำดับ 2', async () => {
  // Given client จำลองคืนคำขอ 2 รายการ แต่ B (ลำดับ 2) มาก่อน A (ลำดับ 1)
  const client = queueClient({ requests: [REQ_B, REQ_A] })

  // When เจ้าหน้าที่เปิดรายละเอียด CS-NB-0001
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)

  // Then หน้าจอเรียง A ลำดับที่ 1 ก่อน B ลำดับที่ 2 (FR-EQS-11)
  const queue = await screen.findByRole('region', { name: 'คำขอยืมที่ค้างอยู่' })
  const rows = within(queue).getAllByRole('row').slice(1)
  expect(within(rows[0]).getByText('ผู้ขอ A')).toBeTruthy()
  expect(within(rows[0]).getByText('1')).toBeTruthy()
  expect(within(rows[1]).getByText('ผู้ขอ B')).toBeTruthy()
  expect(within(rows[1]).getByText('2')).toBeTruthy()
})
