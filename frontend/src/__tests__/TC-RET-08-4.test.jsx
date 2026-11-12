// TC-RET-08-4 (AC-RET-08 ทางปกติ หน้าจอ) จาก specs/006-return/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import OpenLoansPage from '../pages/Return/OpenLoansPage.jsx'

test('TC-RET-08-4: นักศึกษาเปิดหน้ารายการ หลังบ้านตอบ FORBIDDEN_ROLE แสดง error และไม่แสดงแถวการยืม', async () => {
  // Given นักศึกษาเปิดหน้ารายการที่ยังไม่คืน หลังบ้านตอบ 403 FORBIDDEN_ROLE
  const client = {
    returns: {
      openLoans: vi.fn().mockResolvedValue({
        ok: false,
        status: 403,
        data: { code: 'FORBIDDEN_ROLE', message: 'เฉพาะเจ้าหน้าที่เท่านั้นที่บันทึกการรับคืนได้' },
      }),
    },
  }

  // When หน้ารายการโหลดข้อมูล
  render(<OpenLoansPage client={client} onSelect={vi.fn()} />)

  // Then (1) แสดงข้อความ error
  expect((await screen.findByRole('alert')).textContent).not.toBe('')
  // (2) ไม่แสดงแถวการยืมใด ๆ
  expect(screen.queryByRole('table')).toBeNull()
  expect(screen.queryByRole('button', { name: /^เลือก / })).toBeNull()
})
