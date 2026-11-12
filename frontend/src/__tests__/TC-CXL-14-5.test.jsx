// TC-CXL-14-5 (AC-CXL-14 ขอบ) ส่วนหน้าจอ จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
// ส่วนหลังบ้านอยู่ที่ backend/internal/cancelrequest/tc_ac_cxl_14_test.go (TestTC_CXL_14_5_ThaiCharacterCounting)
import { fireEvent, render, screen, within } from '@testing-library/react'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'
import { REQ_A1, fakeClient } from './helpers/cxlFakeClient.js'

test('TC-CXL-14-5 "ที่" 67 ตัว นับตาม code point ได้ 201 ตัวนับแสดง 201/200 และกดยืนยันไม่ได้', async () => {
  // Given เหตุผลภาษาไทย "ที่" ซ้ำ 67 ครั้ง (เห็นบนจอ 67 ตัว แต่เป็น 201 code point)
  const reasonText = 'ที่'.repeat(67)
  expect([...reasonText]).toHaveLength(201)
  const client = fakeClient([REQ_A1])
  render(<MyPendingRequests client={client} />)
  const row = await screen.findByRole('listitem')
  fireEvent.click(within(row).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')

  // When พิมพ์เหตุผลนั้น
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: reasonText } })

  // Then (3) หน้าจอนับแบบเดียวกับหลังบ้าน: ตัวนับแสดง "201/200" และปุ่มยืนยันกดไม่ได้
  expect(within(dialog).getByText('201/200')).toBeTruthy()
  const confirm = within(dialog).getByRole('button', { name: 'ยืนยัน' })
  expect(confirm.disabled).toBe(true)
  fireEvent.click(confirm)
  expect(client.cancelRequest).not.toHaveBeenCalled()
})
