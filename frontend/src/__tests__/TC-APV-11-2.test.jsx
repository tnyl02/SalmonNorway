// TC-APV-11-2 (AC-APV-11 ปกติ) จาก specs/004-approve-request/test-cases.md สถานะใช้ได้ ห้ามแก้ให้ผ่าน
import { render, screen, within } from '@testing-library/react'
import ApproveRequests from '../pages/ApproveRequest/ApproveRequests.jsx'
import { fakeClient } from './helpers/apvFakeClient.js'

test('TC-APV-11-2 แสดงคำขอตามลำดับ พร้อมชื่อผู้ขอ อุปกรณ์ วัตถุประสงค์ และช่วงเวลาตามเวลาไทย', async () => {
  // Given client จำลองคืน 2 คำขอเรียงแล้ว คำขอแรกวัตถุประสงค์ study คำขอที่สองวัตถุประสงค์ "อื่น ๆ" พิมพ์ว่า "ถ่ายวิดีโอโครงงาน"
  // When  เปิดหน้ารายการคำขอ
  render(<ApproveRequests client={fakeClient()} />)
  const table = await screen.findByRole('table')
  const rows = within(table).getAllByRole('row').slice(1)

  // Then (1) แสดง 2 แถวตามลำดับที่ได้มา
  expect(rows).toHaveLength(2)
  expect(within(rows[0]).getByText('นักศึกษาตัวอย่าง ก')).toBeTruthy()
  expect(within(rows[1]).getByText('นักศึกษาตัวอย่าง ข')).toBeTruthy()

  // Then (2) แต่ละแถวแสดงชื่อผู้ขอ ชื่อรุ่นและรหัสครุภัณฑ์ วัตถุประสงค์ และช่วงเวลาเริ่ม ถึง คืน
  rows.forEach((row, i) => {
    expect(within(row).getByText(/อุปกรณ์ทดลองอนุมัติ \(004\)/)).toBeTruthy()
    expect(within(row).getByText(i === 0 ? /AP-TEST-0002/ : /AP-TEST-0003/)).toBeTruthy()
  })

  // Then (3) วัตถุประสงค์แสดง "เพื่อการเรียน" และ "ถ่ายวิดีโอโครงงาน"
  expect(within(rows[0]).getByText('เพื่อการเรียน')).toBeTruthy()
  expect(within(rows[1]).getByText('ถ่ายวิดีโอโครงงาน')).toBeTruthy()

  // Then (4) เวลาแสดงตามเขตเวลา Asia/Bangkok: 02:00Z = 09:00 น. และ 09:00Z = 16:00 น. (ถ้าแสดงตาม UTC จะเห็น 02:00)
  for (const row of rows) {
    expect(row.textContent).toMatch(/09:00/)
    expect(row.textContent).toMatch(/16:00/)
    expect(row.textContent).not.toMatch(/02:00/)
  }
})
