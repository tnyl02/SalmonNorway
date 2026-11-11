// AC-CXL-01 (FR-CXL-01, FR-CXL-07) ยกเลิกคำขอพร้อมเหตุผลแล้วเห็นข้อความสำเร็จ ใช้ client จำลองแทนหลังบ้าน (T-07)
//   Given นักศึกษามีคำขอที่ยังยกเลิกได้ 2 รายการ
//   When  กดยกเลิกแถวแรก พิมพ์เหตุผล แล้วกดยืนยัน
//   Then  client ถูกเรียกยกเลิกด้วย requestId และเหตุผล เห็น "ยกเลิกคำขอยืม {ชื่อรุ่น} ({รหัสครุภัณฑ์}) เรียบร้อยแล้ว"
//         แถวที่ยกเลิกหายจากรายการ แถวอื่นยังอยู่ และกล่องยืนยันปิด
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
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
