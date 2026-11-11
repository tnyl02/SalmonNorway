// AC-CXL-01 (FR-CXL-01, FR-CXL-07) ยกเลิกคำขอพร้อมเหตุผลแล้วเห็นข้อความสำเร็จ ใช้ client จำลองแทนหลังบ้าน (T-07)
//   Given นักศึกษามีคำขอที่ยังยกเลิกได้ 2 รายการ
//   When  กดยกเลิกแถวแรก พิมพ์เหตุผล แล้วกดยืนยัน
//   Then  client ถูกเรียกยกเลิกด้วย requestId และเหตุผล เห็น "ยกเลิกคำขอยืม {ชื่อรุ่น} ({รหัสครุภัณฑ์}) เรียบร้อยแล้ว"
//         แถวที่ยกเลิกหายจากรายการ แถวอื่นยังอยู่ และกล่องยืนยันปิด
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { ApiError } from '../../../api/client.js'
import MyPendingRequests from '../MyPendingRequests.jsx'
import { REQ_A1, REQ_A2, fakeClient } from '../../../__tests__/helpers/cxlFakeClient.js'

test('AC_CXL_01 cancels with reason and shows success', async () => {
  // Given
  const client = fakeClient([REQ_A1, REQ_A2])
  render(<MyPendingRequests client={client} />)
  const rows = await screen.findAllByRole('listitem')

  // When
  fireEvent.click(within(rows[0]).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  expect(client.cancelRequest).not.toHaveBeenCalled() // ต้องยืนยันก่อน (FR-CXL-07)
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then
  await waitFor(() => expect(client.cancelRequest).toHaveBeenCalledWith(101, 'ไม่ต้องใช้อุปกรณ์แล้ว'))
  expect(await screen.findByRole('status')).toHaveProperty(
    'textContent',
    'ยกเลิกคำขอยืม Notebook ตัวอย่างรุ่น A (CS-NB-0001) เรียบร้อยแล้ว',
  )
  const left = screen.getAllByRole('listitem')
  expect(left).toHaveLength(1)
  expect(left[0].textContent).toContain('CS-PJ-0001')
  expect(screen.queryByRole('dialog')).toBeNull()
})

test('AC_CXL_01 keeps row when cancel fails', async () => {
  // Given client จำลองยกเลิกไม่สำเร็จ (ข้อความของแต่ละกรณีอยู่ใน T-16 ขึ้นไป)
  const client = fakeClient([REQ_A1])
  client.cancelRequest.mockRejectedValueOnce(new Error('fail'))
  render(<MyPendingRequests client={client} />)
  const row = await screen.findByRole('listitem')

  // When
  fireEvent.click(within(row).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then แถวยังอยู่ และไม่มีข้อความสำเร็จ
  await waitFor(() => expect(client.cancelRequest).toHaveBeenCalledTimes(1))
  expect(screen.getAllByRole('listitem')).toHaveLength(1)
  expect(screen.queryByRole('status')).toBeNull()
})

// AC-CXL-09 (FR-CXL-07) กดไม่ยืนยันแล้วไม่ยกเลิก เปิดใหม่ช่องเหตุผลว่าง กดยืนยันรอบสองแล้วยกเลิก (T-08)
//   Given นักศึกษามีคำขอที่ยังยกเลิกได้ 1 รายการ
//   When  กดยกเลิก พิมพ์เหตุผล แล้วกดไม่ยืนยัน
//   Then  กล่องปิด client ไม่ถูกเรียก แถวยังอยู่
//   When  กดยกเลิกอีกครั้ง
//   Then  ช่องเหตุผลว่าง พิมพ์ใหม่แล้วกดยืนยัน แถวหายไป
test('AC_CXL_09 confirm dialog cancel then confirm', async () => {
  // Given
  const client = fakeClient([REQ_A1])
  render(<MyPendingRequests client={client} />)
  const row = await screen.findByRole('listitem')

  // When ไม่ยืนยัน
  fireEvent.click(within(row).getByRole('button', { name: 'ยกเลิก' }))
  let dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'พิมพ์ไว้แล้วเปลี่ยนใจ' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ไม่ยืนยัน' }))

  // Then
  expect(screen.queryByRole('dialog')).toBeNull()
  expect(client.cancelRequest).not.toHaveBeenCalled()
  expect(screen.getAllByRole('listitem')).toHaveLength(1)

  // When เปิดใหม่
  fireEvent.click(within(screen.getByRole('listitem')).getByRole('button', { name: 'ยกเลิก' }))
  dialog = await screen.findByRole('dialog')

  // Then ช่องว่าง แล้วยืนยันรอบสอง
  expect(within(dialog).getByRole('textbox').value).toBe('')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))
  await waitFor(() => expect(client.cancelRequest).toHaveBeenCalledTimes(1))
  await waitFor(() => expect(screen.queryAllByRole('listitem')).toHaveLength(0))
})

// AC-CXL-14 (FR-CXL-08) ปุ่มยืนยันกดไม่ได้เมื่อเหตุผลว่าง มีแต่ช่องว่าง หรือเกิน 200 ตัวอักษร (T-10)
//   Given เปิดกล่องยืนยันของคำขอหนึ่ง
//   When  เหตุผลว่าง / ช่องว่างล้วน / 201 ตัวอักษร
//   Then  ปุ่มยืนยันกดไม่ได้ client ไม่ถูกเรียก ตัวนับแสดง "x/200" นับตาม code point
//   When  เหตุผล 200 ตัวอักษรพอดี
//   Then  กดยืนยันได้
test('AC_CXL_14 confirm disabled when reason empty or too long', async () => {
  // Given
  const client = fakeClient([REQ_A1])
  render(<MyPendingRequests client={client} />)
  fireEvent.click(within(await screen.findByRole('listitem')).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  const reason = within(dialog).getByRole('textbox')
  const confirm = within(dialog).getByRole('button', { name: 'ยืนยัน' })

  // Then ตอนเปิดกล่อง ช่องว่าง ตัวนับ 0/200 และกดไม่ได้
  expect(within(dialog).getByText('0/200')).toBeTruthy()
  expect(confirm.disabled).toBe(true)

  // When / Then กรณีกดไม่ได้
  for (const [value, counter] of [['   ', '3/200'], ['a'.repeat(201), '201/200'], ['ที่'.repeat(67), '201/200']]) {
    fireEvent.change(reason, { target: { value } })
    expect(within(dialog).getByText(counter)).toBeTruthy()
    expect(confirm.disabled).toBe(true)
  }
  expect(client.cancelRequest).not.toHaveBeenCalled()

  // When / Then 200 ตัวอักษรพอดี กดได้
  fireEvent.change(reason, { target: { value: 'a'.repeat(200) } })
  expect(within(dialog).getByText('200/200')).toBeTruthy()
  expect(confirm.disabled).toBe(false)
  fireEvent.click(confirm)
  await waitFor(() => expect(client.cancelRequest).toHaveBeenCalledWith(101, 'a'.repeat(200)))
})

// AC-CXL-04 (FR-CXL-04, FR-CXL-07) ฝั่งหน้าจอ T-22
//   Given นักศึกษาเปิดหน้ารายการ กดยกเลิกคำขอหนึ่ง กล่องยืนยันเปิดอยู่
//         แล้วเจ้าหน้าที่อนุมัติคำขอนั้นก่อนนักศึกษากดยืนยัน (client จำลองเปลี่ยนเป็นตอบ 409 ALREADY_APPROVED หลังเปิดกล่อง)
//   When  นักศึกษาพิมพ์เหตุผลแล้วกดยืนยัน
//   Then  ไม่มีข้อความยกเลิกสำเร็จ แสดง "ยกเลิกไม่ได้ เพราะคำขอได้รับการอนุมัติแล้ว"
//         กล่องปิด และแถวนั้นหายหลังโหลดรายการใหม่ (หลังบ้านไม่ส่งคำขอที่อนุมัติแล้ว)
test('AC_CXL_04 approved while dialog open', async () => {
  // Given หน้ารายการแสดง 2 แถว แล้วเปิดกล่องยืนยันของแถวแรก
  const client = fakeClient([REQ_A1, REQ_A2])
  render(<MyPendingRequests client={client} />)
  const rows = await screen.findAllByRole('listitem')
  fireEvent.click(within(rows[0]).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')

  // Given (ต่อ) เจ้าหน้าที่อนุมัติระหว่างที่กล่องเปิดอยู่: หลังบ้านตอบ 409 และรายการใหม่ไม่มี REQ_A1 แล้ว
  client.cancelRequest.mockRejectedValueOnce(new ApiError(409, 'ALREADY_APPROVED'))
  client.listMyPendingRequests.mockResolvedValueOnce({ items: [REQ_A2] })

  // When
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then
  expect((await screen.findByRole('alert')).textContent).toBe('ยกเลิกไม่ได้ เพราะคำขอได้รับการอนุมัติแล้ว')
  expect(client.cancelRequest).toHaveBeenCalledWith(101, 'ไม่ต้องใช้อุปกรณ์แล้ว')
  expect(screen.queryByRole('status')).toBeNull()
  expect(screen.queryByRole('dialog')).toBeNull()
  await waitFor(() => expect(screen.getAllByRole('listitem')).toHaveLength(1))
  expect(screen.getByRole('listitem').textContent).toContain('CS-PJ-0001')
  expect(screen.getByRole('listitem').textContent).not.toContain('CS-NB-0001')
})

