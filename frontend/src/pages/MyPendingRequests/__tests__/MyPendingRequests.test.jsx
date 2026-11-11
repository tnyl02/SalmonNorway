// AC-CXL-08 (FR-CXL-06) หน้า "คำขอที่รออนุมัติของฉัน" ใช้ client จำลองแทนหลังบ้าน (T-05)
//   Given นักศึกษามีคำขอที่ยังยกเลิกได้ 2 รายการ ส่งคนละเวลา (หลังบ้านเรียงใหม่ไปเก่ามาแล้ว)
//   When  เปิดหน้า "คำขอที่รออนุมัติของฉัน"
//   Then  เห็น 2 แถว ส่งล่าสุดอยู่บนสุด แต่ละแถวมีชื่อรุ่น รหัสครุภัณฑ์ เวลาที่ส่งแบบ Asia/Bangkok และปุ่มยกเลิก
import { render, screen, within } from '@testing-library/react'
import MyPendingRequests from '../MyPendingRequests.jsx'
import { REQ_A1, REQ_A2, fakeClient } from '../../../__tests__/helpers/cxlFakeClient.js'

test('AC_CXL_08 shows own pending requests newest first', async () => {
  // Given
  const client = fakeClient([REQ_A1, REQ_A2])

  // When
  render(<MyPendingRequests client={client} />)

  // Then
  const rows = await screen.findAllByRole('listitem')
  expect(rows).toHaveLength(2)
  expect(client.listMyPendingRequests).toHaveBeenCalledTimes(1)

  expect(rows[0].textContent).toContain('Notebook ตัวอย่างรุ่น A')
  expect(rows[0].textContent).toContain('CS-NB-0001')
  expect(rows[0].textContent).toContain('2569-11-10 10:30') // 03:30Z = 10:30 Asia/Bangkok

  expect(rows[1].textContent).toContain('Projector ตัวอย่างรุ่น B')
  expect(rows[1].textContent).toContain('CS-PJ-0001')
  expect(rows[1].textContent).toContain('2569-11-10 09:00')

  for (const row of rows) {
    expect(within(row).getByRole('button', { name: 'ยกเลิก' })).toBeTruthy()
  }
})
