// TC-CXL-09-3 (AC-CXL-09 ขอบ) จาก specs/003-cancel-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import MyPendingRequests from '../pages/MyPendingRequests/MyPendingRequests.jsx'
import { REQ_A1, REQ_A2, fakeClient } from './helpers/cxlFakeClient.js'

test('TC-CXL-09-3 มี 2 รายการ ยกเลิกแถวแรกแล้วแถวที่สองยังอยู่', async () => {
  // Given มีคำขอ 2 รายการ
  const client = fakeClient([REQ_A1, REQ_A2])
  render(<MyPendingRequests client={client} />)
  const rows = await screen.findAllByRole('listitem')

  // When กดยกเลิกแถวแรก พิมพ์เหตุผล แล้วกดยืนยัน
  fireEvent.click(within(rows[0]).getByRole('button', { name: 'ยกเลิก' }))
  const dialog = await screen.findByRole('dialog')
  fireEvent.change(within(dialog).getByRole('textbox'), { target: { value: 'ไม่ต้องใช้อุปกรณ์แล้ว' } })
  fireEvent.click(within(dialog).getByRole('button', { name: 'ยืนยัน' }))

  // Then (1) แถวแรกหายจากหน้า
  await waitFor(() => expect(screen.queryAllByRole('listitem')).toHaveLength(1))
  expect(screen.getAllByRole('listitem').some((li) => li.textContent.includes('CS-NB-0001'))).toBe(false)

  // Then (2) แถวที่สองยังอยู่
  expect(screen.getByRole('listitem').textContent).toContain('CS-PJ-0001')
})
