// TC-RET-03-3 (AC-RET-03 ทางปกติ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen, fireEvent } from '@testing-library/react'
import ReturnPage from '../pages/Return/ReturnPage.jsx'

test('TC-RET-03-3: หลังบ้านตอบ NO_OPEN_LOAN แจ้งข้อความ และไม่แสดงฟอร์มบันทึก', async () => {
  // Given หน้าบันทึกการรับคืน หลังบ้านตอบ NO_OPEN_LOAN
  const client = {
    returns: {
      lookup: vi.fn().mockResolvedValue({
        ok: false,
        status: 409,
        data: { code: 'NO_OPEN_LOAN', message: 'อุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน' },
      }),
      record: vi.fn(),
    },
  }
  render(<ReturnPage client={client} />)

  // When เจ้าหน้าที่กรอกรหัสครุภัณฑ์แล้วกดค้นหา
  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'CS-NB-0001' } })
  fireEvent.click(screen.getByRole('button', { name: 'ค้นหา' }))

  // Then (1) แจ้ง "อุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน"
  expect((await screen.findByRole('alert')).textContent).toBe('อุปกรณ์นี้ไม่มีการยืมที่ยังไม่คืน')
  // (2) ไม่แสดงตัวเลือกสภาพและปุ่มบันทึก (ระบบไม่บันทึก)
  expect(screen.queryByRole('radio')).toBeNull()
  expect(screen.queryByRole('button', { name: 'บันทึกการรับคืน' })).toBeNull()
  expect(client.returns.record).not.toHaveBeenCalled()
})
