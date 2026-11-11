// TC-CXL-14-2 (AC-CXL-14 ปกติ) จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'
import { REQ_A1, fakeClient } from './helpers/cxlFakeClient.js'

test('TC-CXL-14-2 เหตุผลว่าง ช่องว่างล้วน หรือ 201 ตัวอักษร กดยืนยันไม่ได้ เหตุผลถูกต้องกดได้', async () => {
  // Given เปิดกล่องยืนยันของคำขอหนึ่ง
  const client = fakeClient([REQ_A1])
  render(<MyPendingRequests client={client} />)
  const row = await screen.findByRole('listitem')
  fireEvent.click(within(row).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  const reason = within(dialog).getByRole('textbox')
  const confirm = within(dialog).getByRole('button', { name: 'ยืนยัน' })

  // When ลองช่องเหตุผลว่าง, ช่องว่างล้วน, 201 ตัวอักษร
  // Then (1) 3 กรณีแรกปุ่มยืนยันกดไม่ได้ และ client ไม่ถูกเรียก
  for (const value of ['', '   ', 'a'.repeat(201)]) {
    fireEvent.change(reason, { target: { value } })
    expect(confirm.disabled).toBe(true)
    fireEvent.click(confirm)
  }
  expect(client.cancelRequest).not.toHaveBeenCalled()

  // When แก้เป็น "ไม่ต้องใช้อุปกรณ์แล้ว"
  fireEvent.change(reason, { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })

  // Then (2) กรณีสุดท้ายกดยืนยันได้ และ client ถูกเรียกด้วยเหตุผลนั้น 1 ครั้ง
  expect(confirm.disabled).toBe(false)
  fireEvent.click(confirm)
  await waitFor(() => expect(client.cancelRequest).toHaveBeenCalledTimes(1))
  expect(client.cancelRequest).toHaveBeenCalledWith(101, 'ไม่ต้องใช้อุปกรณ์แล้ว')
})
