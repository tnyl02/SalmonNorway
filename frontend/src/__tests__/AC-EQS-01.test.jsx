// AC-EQS-01 (FR-EQS-01) หน้าจอรายการอุปกรณ์ ใช้ client จำลองแทนหลังบ้าน
//   Given นักศึกษายืนยันตัวตนแล้ว และมีอุปกรณ์รุ่นหนึ่ง 5 ชิ้น ว่าง 3 ชิ้น
//   When  เปิดหน้ารายการอุปกรณ์ แล้วเลือกรุ่นนั้น
//   Then  รายการแสดงรุ่นนั้น 1 แถว ระบุประเภทและ "ว่าง 3 / ทั้งหมด 5"
//         และเมื่อเลือกแล้วเห็นอุปกรณ์ 5 ชิ้นตามรหัสครุภัณฑ์พร้อมสถานะของแต่ละชิ้น
import { fireEvent, render, screen, within } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'

// Given: client จำลองคืนข้อมูลแบบเดียวกับที่หลังบ้านตอบ (plan.md ข้อ 5)
function fakeClient() {
  return {
    listModels: async () => ({
      items: [
        { modelId: 1, name: 'Notebook ตัวอย่างรุ่น A', category: 'คอมพิวเตอร์พกพา', availableCount: 3, totalCount: 5 },
        { modelId: 2, name: 'Projector ตัวอย่างรุ่น B', category: 'อุปกรณ์นำเสนอ', availableCount: 1, totalCount: 2 },
      ],
      total: 2,
    }),
    listItems: async (modelId) => {
      if (modelId !== 1) return { items: [], total: 0 }
      return {
        items: [
          { assetCode: 'CS-NB-0001', status: 'available' },
          { assetCode: 'CS-NB-0002', status: 'available' },
          { assetCode: 'CS-NB-0003', status: 'available' },
          { assetCode: 'CS-NB-0004', status: 'borrowed' },
          { assetCode: 'CS-NB-0005', status: 'borrowed' },
        ],
        total: 5,
      }
    },
  }
}

test('AC-EQS-01 รุ่น 1 แถว แสดงประเภทและ "ว่าง 3 / ทั้งหมด 5" แล้วเลือกรุ่นเห็น 5 ชิ้นพร้อมสถานะ', async () => {
  // When: เปิดหน้ารายการอุปกรณ์
  render(<EquipmentList client={fakeClient()} />)

  // Then: รุ่นนั้น 1 แถว ระบุประเภทและ "ว่าง 3 / ทั้งหมด 5"
  const rows = await screen.findAllByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })
  expect(rows).toHaveLength(1)
  const row = rows[0]
  expect(within(row).getByText('คอมพิวเตอร์พกพา')).toBeTruthy()
  expect(within(row).getByText('ว่าง 3 / ทั้งหมด 5')).toBeTruthy()

  // When: เลือกรุ่นนั้น
  fireEvent.click(row)

  // Then: เห็น 5 ชิ้นตามรหัสครุภัณฑ์พร้อมสถานะของแต่ละชิ้น
  const table = await screen.findByRole('table')
  const bodyRows = within(table).getAllByRole('row').slice(1) // ตัดแถวหัวตาราง
  expect(bodyRows).toHaveLength(5)
  const expected = [
    ['CS-NB-0001', 'ว่าง'],
    ['CS-NB-0002', 'ว่าง'],
    ['CS-NB-0003', 'ว่าง'],
    ['CS-NB-0004', 'ถูกยืม'],
    ['CS-NB-0005', 'ถูกยืม'],
  ]
  expected.forEach(([code, label], i) => {
    expect(within(bodyRows[i]).getByText(code)).toBeTruthy()
    expect(within(bodyRows[i]).getByText(label)).toBeTruthy()
  })
})
