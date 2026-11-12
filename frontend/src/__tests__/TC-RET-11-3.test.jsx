// TC-RET-11-3 (AC-RET-11 ทางปกติ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen, fireEvent } from '@testing-library/react'
import ReturnPage from '../pages/Return/ReturnPage.jsx'

test('TC-RET-11-3: ค้นหาแล้วได้ FORBIDDEN_DEPARTMENT แจ้งว่าไม่มีสิทธิ์ และไม่แสดงฟอร์ม', async () => {
  // Given หน้าบันทึกการรับคืน หลังบ้านตอบ FORBIDDEN_DEPARTMENT ตอนค้นหา
  const client = {
    returns: {
      lookup: vi.fn().mockResolvedValue({
        ok: false,
        status: 403,
        data: { code: 'FORBIDDEN_DEPARTMENT', message: 'ไม่มีสิทธิ์: อุปกรณ์นี้เป็นของภาควิชาอื่น' },
      }),
      record: vi.fn(),
    },
  }
  render(<ReturnPage client={client} />)

  // When เจ้าหน้าที่กรอกรหัสครุภัณฑ์แล้วกดค้นหา
  fireEvent.change(screen.getByLabelText('รหัสครุภัณฑ์'), { target: { value: 'XX-NB-0002' } })
  fireEvent.click(screen.getByRole('button', { name: 'ค้นหา' }))

  // Then (1) แจ้งว่าไม่มีสิทธิ์
  expect((await screen.findByRole('alert')).textContent).toContain('ไม่มีสิทธิ์')
  // (2) ไม่แสดงตัวเลือกสภาพและปุ่มบันทึก
  expect(screen.queryByRole('radio')).toBeNull()
  expect(screen.queryByRole('button', { name: 'บันทึกการรับคืน' })).toBeNull()
})
