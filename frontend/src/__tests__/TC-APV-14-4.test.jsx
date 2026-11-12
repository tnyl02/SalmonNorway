// TC-APV-14-4 (AC-APV-14 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import FailedEmails from '../pages/ApproveRequest/FailedEmails.jsx'

test('TC-APV-14-4 แสดงสาเหตุและเวลา กดส่งซ้ำแล้วเรียก client ด้วย emailId และแถวนั้นหายจากรายการ', async () => {
  // Given client จำลองคืน Email ที่ส่งไม่สำเร็จ 2 รายการ และตอบ 200 เมื่อส่งซ้ำ
  // lastAttemptAt 03:03Z = 10:03 น. เวลาไทย
  const calls = []
  const client = {
    getFailedEmails: async () => ({
      items: [
        { emailId: 71, requestId: 1, kind: 'approved', requesterName: 'ผู้ขอแถวแรก', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)',
          assetCode: 'AP-TEST-0002', attempts: 4, lastError: 'smtp down', lastAttemptAt: '2026-11-12T03:03:00Z' },
        { emailId: 72, requestId: 2, kind: 'auto_rejected', requesterName: 'ผู้ขอแถวสอง', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)',
          assetCode: 'AP-TEST-0003', attempts: 4, lastError: 'smtp timeout', lastAttemptAt: '2026-11-12T04:30:00Z' },
      ],
    }),
    resendEmail: async (emailId) => {
      calls.push(emailId)
      return { emailId, status: 'pending' }
    },
  }

  // When เปิดหน้า "Email ที่ส่งไม่สำเร็จ"
  render(<FailedEmails client={client} />)
  const table = await screen.findByRole('table')
  const [row1, row2] = within(table).getAllByRole('row').slice(1)

  // Then (1) แต่ละแถวแสดงสาเหตุและเวลาที่ลองส่งครั้งล่าสุด
  expect(within(row1).getByText('smtp down')).toBeTruthy()
  expect(row1.textContent).toMatch(/10:03/)
  expect(within(row2).getByText('smtp timeout')).toBeTruthy()
  expect(row2.textContent).toMatch(/11:30/)

  // When กด "ส่งซ้ำ" ที่แถวแรก
  fireEvent.click(within(row1).getByRole('button', { name: 'ส่งซ้ำ' }))

  // Then (2) client ถูกเรียกส่งซ้ำ 1 ครั้งด้วย emailId ของแถวนั้น
  await waitFor(() => expect(calls).toEqual([71]))

  // Then (3) แถวนั้นหายจากรายการ อีกแถวยังอยู่
  await waitFor(() => expect(screen.queryByText('ผู้ขอแถวแรก')).toBeNull())
  expect(screen.getByText('ผู้ขอแถวสอง')).toBeTruthy()
})
