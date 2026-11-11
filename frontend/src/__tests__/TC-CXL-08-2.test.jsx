// TC-CXL-08-2 (AC-CXL-08 ปกติ) จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { render, screen, within } from '@testing-library/react'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'
import { REQ_A1, REQ_A2, fakeClient } from './helpers/cxlFakeClient.js'

test('TC-CXL-08-2 แสดงคำขอที่รออนุมัติของตนเอง ใหม่ไปเก่า พร้อมชื่อรุ่น รหัสครุภัณฑ์ เวลา และปุ่มยกเลิก', async () => {
  // Given client จำลองคืนคำขอที่รออนุมัติของ A 2 รายการ เวลาส่งต่างกัน (REQ_A1 ส่งหลัง REQ_A2)
  const client = fakeClient([REQ_A1, REQ_A2])

  // When เปิดหน้า "คำขอที่รออนุมัติของฉัน"
  render(<MyPendingRequests client={client} />)

  // Then (1) แสดง 2 แถว
  const rows = await screen.findAllByRole('listitem')
  expect(rows).toHaveLength(2)

  // Then (2) แถวที่ส่งล่าสุดอยู่บนสุด
  expect(rows[0].textContent).toContain('CS-NB-0001')
  expect(rows[1].textContent).toContain('CS-PJ-0001')

  // Then (3) แต่ละแถวแสดงชื่อรุ่น รหัสครุภัณฑ์ และเวลาที่ส่งตามเขตเวลา Asia/Bangkok
  expect(rows[0].textContent).toContain('Notebook ตัวอย่างรุ่น A')
  expect(rows[0].textContent).toContain('10:30')
  expect(rows[1].textContent).toContain('Projector ตัวอย่างรุ่น B')
  expect(rows[1].textContent).toContain('09:00')

  // Then (4) แต่ละแถวมีปุ่มยกเลิก
  for (const row of rows) {
    expect(within(row).getByRole('button', { name: 'ยกเลิก' })).toBeTruthy()
  }
})
