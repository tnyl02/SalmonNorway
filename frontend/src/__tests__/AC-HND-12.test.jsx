// AC-HND-12 (FR-HND-06) ตรวจรูปที่หน้าจอก่อนส่ง ใช้ client จำลองแทนหลังบ้าน
//   Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว
//   When  เจ้าหน้าที่บันทึกการส่งมอบโดยแนบรูป (ก) 1 รูป หรือ (ข) ไฟล์ GIF
//   Then  ระบบไม่บันทึกการส่งมอบในทุกกรณี และไม่มี Loan ถูกสร้าง
// หน้าจอไม่เรียกหลังบ้านเลยในกรณีนี้ (หลังบ้านตรวจซ้ำใน TestAC_HND_12_InvalidPhotos) กรณี (ค) 11 MB รอ Q-HND-10 (T-28)
import { fireEvent, render, screen } from '@testing-library/react'
import HandoverForm from '../pages/Handover/HandoverForm.jsx'

const request = { requestId: 11, assetCode: 'HO-TEST-0001', requesterStudentId: 'student-001' }

function fakeClient() {
  const calls = []
  return {
    calls,
    verifyBorrower: async () => ({ result: 'verified', userId: 'student-001' }),
    handover: async (requestId, form) => {
      calls.push({ requestId, form })
      return { loanId: 1, assetCode: 'HO-TEST-0001', handedOverAt: '2026-11-12T03:00:00Z', dueAt: '2026-11-12T09:00:00Z' }
    },
  }
}

const jpg = (name) => new File(['x'], name, { type: 'image/jpeg' })
const gif = (name) => new File(['x'], name, { type: 'image/gif' })

test.each([
  ['ก 1 รูป', [jpg('a.jpg')]],
  ['ข ไฟล์ GIF', [jpg('a.jpg'), gif('b.gif')]],
])('AC-HND-12 %s ไม่ส่งไปหลังบ้าน และแจ้งเงื่อนไขรูป', async (_, photos) => {
  const client = fakeClient()
  render(<HandoverForm request={request} client={client} />)
  fireEvent.change(screen.getByLabelText('บัญชีผู้ยืม'), { target: { value: 'mock-student-001' } })
  fireEvent.click(screen.getByRole('button', { name: 'ยืนยันตัวตนผู้ยืม' }))
  await screen.findByText(/ยืนยันตัวตนผู้ยืมแล้ว/)
  fireEvent.change(screen.getByLabelText('ชื่อผู้ยืม (ตามบัตร)'), { target: { value: 'สมชาย' } })
  fireEvent.change(screen.getByLabelText('นามสกุลผู้ยืม (ตามบัตร)'), { target: { value: 'ใจดี' } })
  fireEvent.change(screen.getByLabelText('รูปสภาพอุปกรณ์'), { target: { files: photos } })
  fireEvent.click(screen.getByLabelText('ตรวจบัตรนักศึกษาหรือบัตรประชาชนของผู้ยืมแล้ว'))
  fireEvent.click(screen.getByRole('button', { name: 'บันทึกการส่งมอบ' }))

  expect(await screen.findByRole('alert')).toBeTruthy()
  expect(client.calls).toHaveLength(0)
})
