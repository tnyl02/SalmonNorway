// T-34 ปุ่มส่งคำขอยืม (FR-EQS-04, 05, 06) ส่วน "มี/ไม่มีปุ่มส่งคำขอยืม" ของ AC-EQS-04, AC-EQS-05, AC-EQS-06, AC-EQS-14
// ใช้ client จำลองแทนหลังบ้าน หลังบ้านเป็นคนคำนวณ canRequestLoan จากสถานะ (TestAC_EQS_04_06_14_CanRequestLoan)
//   ว่าง = มีปุ่ม "ส่งคำขอยืม" ทั้งตารางรายชิ้นและหน้ารายละเอียด (FR-EQS-06, AC-EQS-05 ปุ่มกลับมาใช้งานได้)
//   ถูกยืม (AC-EQS-04) เลยกำหนดคืน (AC-EQS-14) ไม่พร้อมใช้งาน (AC-EQS-06) = ไม่มีปุ่ม
//   กดปุ่มแล้วส่ง assetCode ของชิ้นนั้นออกไป (FR-EQS-06 ส่งข้อมูลอุปกรณ์ที่เลือกไปยัง UC-02)
import { fireEvent, render, screen, within } from '@testing-library/react'
import { vi } from 'vitest'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'

const BUTTON = 'ส่งคำขอยืม'

const ITEMS = {
  available: { assetCode: 'CS-NB-0001', status: 'available', canRequestLoan: true },
  borrowed: { assetCode: 'CS-NB-0004', status: 'borrowed', dueAt: '2026-11-10T09:00:00Z', canRequestLoan: false },
  overdue: { assetCode: 'CS-NB-0005', status: 'borrowed', isOverdue: true, canRequestLoan: false },
  unavailable: { assetCode: 'CS-PJ-0002', status: 'unavailable', canRequestLoan: false },
}

function listClient() {
  return {
    listModels: async () => ({
      items: [{ modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 1, totalCount: 4 }],
      total: 1,
    }),
    listItems: async () => ({
      items: Object.values(ITEMS).map(({ assetCode, status, canRequestLoan }) => ({ assetCode, status, canRequestLoan })),
      total: 4,
    }),
  }
}

const detailClient = (item) => ({ getItem: async () => item })

test('ตารางรายชิ้น: มีปุ่มเฉพาะชิ้นที่ว่าง กดแล้วส่ง assetCode ของชิ้นนั้น (AC-EQS-04, 05, 06, 14)', async () => {
  const onRequestLoan = vi.fn()
  render(<EquipmentList client={listClient()} onRequestLoan={onRequestLoan} />)
  fireEvent.click(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ }))

  const table = await screen.findByRole('table')
  const rowOf = (code) => within(table).getByText(code).closest('tr')

  // ว่าง = มีปุ่ม
  expect(within(rowOf('CS-NB-0001')).getByRole('button', { name: BUTTON })).toBeTruthy()
  // ถูกยืม เลยกำหนดคืน ไม่พร้อมใช้งาน = ไม่มีปุ่ม
  for (const code of ['CS-NB-0004', 'CS-NB-0005', 'CS-PJ-0002']) {
    expect(within(rowOf(code)).queryByRole('button', { name: BUTTON })).toBeNull()
  }
  expect(within(table).getAllByRole('button', { name: BUTTON })).toHaveLength(1)

  // กดแล้วส่ง assetCode ของชิ้นนั้นออกไป และไม่เปิดหน้ารายละเอียดแทน
  fireEvent.click(within(rowOf('CS-NB-0001')).getByRole('button', { name: BUTTON }))
  expect(onRequestLoan).toHaveBeenCalledTimes(1)
  expect(onRequestLoan).toHaveBeenCalledWith('CS-NB-0001')
})

test('หน้ารายละเอียด: ว่าง = มีปุ่ม กดแล้วส่ง assetCode (FR-EQS-06, AC-EQS-05)', async () => {
  const onRequestLoan = vi.fn()
  render(<EquipmentDetail assetCode="CS-NB-0001" client={detailClient(ITEMS.available)} onRequestLoan={onRequestLoan} />)

  fireEvent.click(await screen.findByRole('button', { name: BUTTON }))
  expect(onRequestLoan).toHaveBeenCalledWith('CS-NB-0001')
})

test.each([
  ['ถูกยืม (AC-EQS-04)', ITEMS.borrowed, 'ถูกยืม'],
  ['เลยกำหนดคืน (AC-EQS-14)', ITEMS.overdue, 'เลยกำหนดคืน'],
  ['ไม่พร้อมใช้งาน (AC-EQS-06)', ITEMS.unavailable, 'ไม่พร้อมใช้งาน'],
])('หน้ารายละเอียด: %s = ไม่มีปุ่มส่งคำขอยืม', async (_, item, label) => {
  render(<EquipmentDetail assetCode={item.assetCode} client={detailClient(item)} onRequestLoan={vi.fn()} />)

  expect(await screen.findByText(label)).toBeTruthy()
  expect(screen.queryByRole('button', { name: BUTTON })).toBeNull()
})

test('ไม่มี canRequestLoan จากหลังบ้าน = ไม่มีปุ่ม (หน้าจอไม่ตัดสินเองจากสถานะ plan.md ข้อ 6.7)', async () => {
  render(<EquipmentDetail assetCode="CS-NB-0001" client={detailClient({ assetCode: 'CS-NB-0001', status: 'available' })} />)

  expect(await screen.findByText('ว่าง')).toBeTruthy()
  expect(screen.queryByRole('button', { name: BUTTON })).toBeNull()
})
