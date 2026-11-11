// TC-CXL-12-2 (AC-CXL-12 ปกติ) จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'
import { fakeClient } from './helpers/cxlFakeClient.js'

test('TC-CXL-12-2 ไม่มีคำขอ แสดง "ไม่มีคำขอที่รออนุมัติ" และไม่มีแถวหรือปุ่มยกเลิก', async () => {
  // Given client จำลองคืน items: []
  const client = fakeClient([])

  // When เปิดหน้า "คำขอที่รออนุมัติของฉัน"
  render(<MyPendingRequests client={client} />)

  // Then (1) แสดงข้อความ "ไม่มีคำขอที่รออนุมัติ" ตรงตัว
  expect(await screen.findByText('ไม่มีคำขอที่รออนุมัติ')).toBeTruthy()

  // Then (2) ไม่มีแถวและไม่มีปุ่มยกเลิก
  expect(screen.queryAllByRole('listitem')).toHaveLength(0)
  expect(screen.queryAllByRole('button', { name: 'ยกเลิก' })).toHaveLength(0)
})
