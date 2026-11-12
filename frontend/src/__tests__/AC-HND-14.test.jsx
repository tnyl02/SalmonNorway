// AC-HND-14 (FR-HND-01) หน้ารายการคำขอที่รอส่งมอบ ใช้ client จำลองแทนหลังบ้าน
//   Given มีคำขอที่ได้รับอนุมัติแล้วและยังไม่ส่งมอบ 2 รายการ (อุปกรณ์คนละชิ้น) คำขอที่ยังไม่อนุมัติ 1 รายการ และคำขอที่ส่งมอบแล้ว 1 รายการ
//   When  เจ้าหน้าที่เปิดรายการคำขอเพื่อส่งมอบ
//   Then  รายการแสดงเฉพาะคำขอที่ได้รับอนุมัติแล้วและยังไม่ส่งมอบทั้ง 2 รายการ
// การคัดกรองทำที่หลังบ้าน (TestAC_HND_14_ListApprovedOnly) client จำลองคืนเฉพาะ 2 รายการแบบเดียวกับหลังบ้าน
import { render, screen, within } from '@testing-library/react'
import HandoverList from '../pages/Handover/HandoverList.jsx'

function fakeClient() {
  return {
    listRequests: async () => [
      { requestId: 11, assetCode: 'HO-TEST-0001', modelName: 'อุปกรณ์ทดลองส่งมอบ (005)', requesterStudentId: 'student-001', submittedAt: '2026-11-12T02:00:00Z' },
      { requestId: 12, assetCode: 'HO-TEST-0002', modelName: 'อุปกรณ์ทดลองส่งมอบ (005)', requesterStudentId: 'student-002', submittedAt: '2026-11-12T03:00:00Z' },
    ],
  }
}

test('AC-HND-14 แสดงเฉพาะคำขอที่อนุมัติแล้วและยังไม่ส่งมอบ 2 รายการ', async () => {
  render(<HandoverList client={fakeClient()} onSelect={() => {}} />)

  const table = await screen.findByRole('table')
  const rows = within(table).getAllByRole('row').slice(1) // ตัดแถวหัวตาราง
  expect(rows).toHaveLength(2)
  expect(within(rows[0]).getByText('HO-TEST-0001')).toBeTruthy()
  expect(within(rows[1]).getByText('HO-TEST-0002')).toBeTruthy()
})
