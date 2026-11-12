// TC-CXL-09-5 (AC-CXL-09 ผิด) จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'
import { REQ_A1, fakeClient } from './helpers/cxlFakeClient.js'

test('TC-CXL-09-5 พิมพ์เหตุผลแล้วกดไม่ยืนยัน เปิดใหม่ช่องเหตุผลว่างและตัวนับเป็น 0/200', async () => {
  // Given พิมพ์เหตุผลไว้แล้วกดไม่ยืนยัน
  const client = fakeClient([REQ_A1])
  render(<MyPendingRequests client={client} />)
  const row = await screen.findByRole('listitem')
  fireEvent.click(within(row).getByRole('button', { name: 'ยกเลิก' }))
  let dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'พิมพ์ไว้แล้วเปลี่ยนใจ' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ไม่ยืนยัน' }))
  await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull())

  // When กดยกเลิกแถวเดิมอีกครั้ง
  fireEvent.click(within(screen.getByRole('listitem')).getByRole('button', { name: 'ยกเลิก' }))
  dialog = await screen.findByRole('dialog')

  // Then (1) แถวยังอยู่และ client ไม่ถูกเรียก
  expect(screen.getAllByRole('listitem')).toHaveLength(1)
  expect(client.cancelRequest).not.toHaveBeenCalled()

  // Then (2) ช่องเหตุผลว่าง และตัวนับแสดง "0/200"
  expect(within(dialog).getByRole('textbox').value).toBe('')
  expect(within(dialog).getByText('0/200')).toBeTruthy()
})
