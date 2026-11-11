// AC-EQS-06 (FR-EQS-05) ส่วนป้าย ใช้ client จำลองแทนหลังบ้าน
//   Given อุปกรณ์มีสถานะ "ไม่พร้อมใช้งาน"
//   When  นักศึกษาเปิดดูรายละเอียด
//   Then  แสดงป้าย "ไม่พร้อมใช้งาน"
// ส่วน "ไม่มีปุ่มส่งคำขอยืม" รอ Q-16 (T-34)
import { render, screen, within } from '@testing-library/react'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'
import ItemTable from '../components/ItemTable.jsx'

test('AC-EQS-06 อุปกรณ์ไม่พร้อมใช้งาน แสดงป้าย "ไม่พร้อมใช้งาน"', async () => {
  // Given: หลังบ้านตอบสถานะ unavailable (เช่น CS-PJ-0002 ในข้อมูลตั้งต้น)
  const client = { getItem: async () => ({ assetCode: 'CS-PJ-0002', status: 'unavailable' }) }

  // When: นักศึกษาเปิดดูรายละเอียด
  render(<EquipmentDetail assetCode="CS-PJ-0002" client={client} />)

  // Then: แสดงป้าย "ไม่พร้อมใช้งาน" ไม่มีกำหนดคืนหรือป้ายเลยกำหนดคืน และไม่ได้แสดงคำภาษาอังกฤษจากหลังบ้าน
  expect(await screen.findByText('ไม่พร้อมใช้งาน')).toBeTruthy()
  expect(screen.queryByText('กำหนดคืน')).toBeNull()
  expect(screen.queryByText('เลยกำหนดคืน')).toBeNull()
  expect(screen.queryByText('unavailable')).toBeNull()
})

test('FR-EQS-05 ตารางรายชิ้นแสดงป้าย "ไม่พร้อมใช้งาน" ด้วย', () => {
  render(<ItemTable items={[{ assetCode: 'CS-PJ-0002', status: 'unavailable' }]} />)
  const row = screen.getAllByRole('row')[1]
  expect(within(row).getByText('ไม่พร้อมใช้งาน')).toBeTruthy()
})
