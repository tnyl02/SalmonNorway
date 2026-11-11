// AC-EQS-04 (FR-EQS-04, CON-DATA-02) ส่วนสถานะและกำหนดคืน ใช้ client จำลองแทนหลังบ้าน
//   Given เจ้าหน้าที่บันทึกการส่งมอบอุปกรณ์ชิ้นหนึ่งใน UC-08 แล้ว และยังไม่เลยกำหนดคืน
//   When  นักศึกษาเปิดดูอุปกรณ์ชิ้นนั้น
//   Then  สถานะแสดงเป็น "ถูกยืม" พร้อมกำหนดคืน (วันที่และเวลา)
// ส่วน "ไม่มีปุ่มส่งคำขอยืม" รอ Q-16 (T-34)
import { fireEvent, render, screen, within } from '@testing-library/react'
import EquipmentDetail, { formatDueAt } from '../pages/EquipmentDetail/EquipmentDetail.jsx'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'

// 2026-11-10T09:00:00Z = 16:00 วันที่ 10 พ.ย. 2569 ตามเวลา Asia/Bangkok (UTC+7)
const BORROWED = { assetCode: 'CS-NB-0004', status: 'borrowed', dueAt: '2026-11-10T09:00:00Z' }

function fakeClient() {
  return {
    listModels: async () => ({
      items: [{ modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 3, totalCount: 5 }],
      total: 1,
    }),
    listItems: async () => ({
      items: [
        { assetCode: 'CS-NB-0001', status: 'available' },
        { assetCode: 'CS-NB-0004', status: 'borrowed' },
      ],
      total: 2,
    }),
    getItem: async (code) => (code === 'CS-NB-0004' ? BORROWED : { assetCode: code, status: 'available' }),
  }
}

test('AC-EQS-04 อุปกรณ์ที่ถูกยืมแสดง "ถูกยืม" พร้อมกำหนดคืนเป็นวันที่และเวลา', async () => {
  // Given: ชิ้นที่ถูกยืม กำหนดคืน 16:00 (ยังไม่เลยกำหนด)
  // When: นักศึกษาเปิดดูอุปกรณ์ชิ้นนั้น
  render(<EquipmentDetail assetCode="CS-NB-0004" client={fakeClient()} />)

  // Then: สถานะ "ถูกยืม" พร้อมกำหนดคืน (วันที่และเวลา)
  expect(await screen.findByText('ถูกยืม')).toBeTruthy()
  expect(screen.getByText('กำหนดคืน')).toBeTruthy()
  expect(screen.getByText('2569-11-10 16:00')).toBeTruthy()
})

test('FR-EQS-04 ชิ้นที่ว่างไม่มีกำหนดคืน', async () => {
  render(<EquipmentDetail assetCode="CS-NB-0001" client={fakeClient()} />)
  expect(await screen.findByText('ว่าง')).toBeTruthy()
  expect(screen.queryByText('กำหนดคืน')).toBeNull()
})

test('formatDueAt แสดงตามเวลา Asia/Bangkok ปี พ.ศ. แม้เวลา UTC จะเป็นคนละวัน', () => {
  // 2026-10-09T17:30:00Z = 00:30 วันที่ 10 ต.ค. 2569 เวลาไทย
  expect(formatDueAt('2026-10-09T17:30:00Z')).toBe('2569-10-10 00:30')
  expect(formatDueAt('2026-10-10T09:00:00Z')).toBe('2569-10-10 16:00')
})

test('กดรหัสครุภัณฑ์ในตารางรายชิ้นแล้วเปิดรายละเอียด และกดกลับไปรายการเดิม', async () => {
  render(<EquipmentList client={fakeClient()} />)
  fireEvent.click(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ }))
  const table = await screen.findByRole('table')

  fireEvent.click(within(table).getByRole('button', { name: 'CS-NB-0004' }))
  expect(await screen.findByText('2569-11-10 16:00')).toBeTruthy()

  fireEvent.click(screen.getByRole('button', { name: 'กลับไปรายการอุปกรณ์' }))
  expect(await screen.findByRole('table')).toBeTruthy()
})
