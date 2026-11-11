// AC-EQS-14 (FR-EQS-04, DOM-PRIV-01, DOM-LOAN-01, CR-001, CR-002) ใช้ client จำลองแทนหลังบ้าน
//   Given อุปกรณ์ชิ้นหนึ่งสถานะ "ถูกยืม" แบบปกติ กำหนดคืน 2569-10-10 16:00 และเวลาปัจจุบันคือ 2569-10-10 17:00
//   When  นักศึกษาเปิดดูรายละเอียด และเจ้าหน้าที่เปิดดูรายละเอียดเดียวกัน
//   Then  นักศึกษาเห็นสถานะ "ถูกยืม" กับป้าย "เลยกำหนดคืน" ไม่เห็นกำหนดคืน 2569-10-10 16:00
//         ส่วนเจ้าหน้าที่เห็นทั้งกำหนดคืน 2569-10-10 16:00 และป้าย "เลยกำหนดคืน"
// ส่วน "ไม่มีปุ่มส่งคำขอยืม" รอ Q-16 (T-34)
// หลังบ้านเป็นคนตัดสินว่าเลยกำหนดและเลือกข้อมูลตามบทบาท (T-12) client จำลองคืนผลแบบเดียวกับหลังบ้าน
import { render, screen } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

// 2026-10-10T09:00:00Z = 2569-10-10 16:00 เวลาไทย
const clientFor = (response) => ({ getItem: async () => response })
const STUDENT_VIEW = { assetCode: 'CS-NB-0005', status: 'borrowed', isOverdue: true }
const STAFF_VIEW = { assetCode: 'CS-NB-0005', status: 'borrowed', dueAt: '2026-10-10T09:00:00Z', isOverdue: true }

test('AC-EQS-14 นักศึกษาเห็น "ถูกยืม" กับป้าย "เลยกำหนดคืน" ไม่เห็นกำหนดคืน', async () => {
  // When: นักศึกษาเปิดดูรายละเอียด
  render(<EquipmentDetail assetCode="CS-NB-0005" client={clientFor(STUDENT_VIEW)} />)

  // Then
  expect(await screen.findByText('ถูกยืม')).toBeTruthy()
  expect(screen.getByText('เลยกำหนดคืน')).toBeTruthy()
  expect(screen.queryByText('กำหนดคืน')).toBeNull()
  expect(screen.queryByText('2569-10-10 16:00')).toBeNull()
})

test('AC-EQS-14 เจ้าหน้าที่เห็นทั้งกำหนดคืน 2569-10-10 16:00 และป้าย "เลยกำหนดคืน"', async () => {
  // When: เจ้าหน้าที่เปิดดูรายละเอียดเดียวกัน
  render(<EquipmentDetail assetCode="CS-NB-0005" client={clientFor(STAFF_VIEW)} />)

  // Then
  expect(await screen.findByText('ถูกยืม')).toBeTruthy()
  expect(screen.getByText('เลยกำหนดคืน')).toBeTruthy()
  expect(screen.getByText('2569-10-10 16:00')).toBeTruthy()
})

test('FR-EQS-04 ข้อ ก ยังไม่เลยกำหนด ไม่มีป้าย "เลยกำหนดคืน"', async () => {
  render(
    <EquipmentDetail
      assetCode="CS-NB-0004"
      client={clientFor({ assetCode: 'CS-NB-0004', status: 'borrowed', dueAt: '2026-10-10T09:00:00Z' })}
    />,
  )
  expect(await screen.findByText('2569-10-10 16:00')).toBeTruthy()
  expect(screen.queryByText('เลยกำหนดคืน')).toBeNull()
})
