// AC-HND-02 (FR-HND-02, FR-HND-03) ฟอร์มบันทึกการส่งมอบ ใช้ client จำลองแทนหลังบ้าน
//   Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว Auth API ยืนยันตัวตนผู้ยืมสำเร็จ และผู้ยืมเป็นผู้ส่งคำขอ
//   When  เจ้าหน้าที่ยืนยันว่าตรวจบัตรแล้ว และบันทึกการส่งมอบ
//   Then  ระบบบันทึกการส่งมอบ พร้อมชื่อ นามสกุล รหัสนักศึกษาของผู้ยืม และวันที่และเวลาที่ส่งมอบ
//         ถ้าเจ้าหน้าที่ไม่ยืนยันการตรวจบัตร ระบบไม่บันทึกการส่งมอบ
// หลังบ้านตรวจการบันทึกจริงใน TestAC_HND_02_* ที่นี่ตรวจว่าหน้าจอส่งข้อมูลครบและไม่ส่งเมื่อไม่ติ๊กตรวจบัตร
import { fireEvent, render, screen } from '@testing-library/react'
import HandoverForm from '../pages/Handover/HandoverForm.jsx'

const request = { requestId: 11, assetCode: 'HO-TEST-0001', requesterStudentId: 'student-001' }

function fakeClient() {
  const calls = []
  return {
    calls,
    handover: async (requestId, form) => {
      calls.push({ requestId, form })
      return { loanId: 1, assetCode: 'HO-TEST-0001', handedOverAt: '2026-11-12T03:00:00Z', dueAt: '2026-11-12T09:00:00Z' }
    },
  }
}

function fillBorrower() {
  fireEvent.change(screen.getByLabelText('ชื่อผู้ยืม (ตามบัตร)'), { target: { value: 'สมชาย' } })
  fireEvent.change(screen.getByLabelText('นามสกุลผู้ยืม (ตามบัตร)'), { target: { value: 'ใจดี' } })
  const photos = [new File(['a'], 'a.jpg', { type: 'image/jpeg' }), new File(['b'], 'b.jpg', { type: 'image/jpeg' })]
  fireEvent.change(screen.getByLabelText('รูปสภาพอุปกรณ์'), { target: { files: photos } })
}

test('AC-HND-02 ไม่ติ๊กว่าตรวจบัตรแล้ว ปุ่มบันทึกใช้ไม่ได้และไม่ส่งข้อมูล', () => {
  const client = fakeClient()
  render(<HandoverForm request={request} client={client} />)
  fillBorrower()

  const save = screen.getByRole('button', { name: 'บันทึกการส่งมอบ' })
  expect(save.disabled).toBe(true)
  fireEvent.submit(screen.getByRole('form', { name: 'บันทึกการส่งมอบ' }))
  expect(client.calls).toHaveLength(0)
})

test('AC-HND-02 ติ๊กตรวจบัตรแล้ว ส่งชื่อ นามสกุล ผลตรวจบัตร และรูป แล้วแสดงเวลาที่ส่งมอบ', async () => {
  const client = fakeClient()
  render(<HandoverForm request={request} client={client} />)
  fillBorrower()
  fireEvent.click(screen.getByLabelText('ตรวจบัตรนักศึกษาหรือบัตรประชาชนของผู้ยืมแล้ว'))
  fireEvent.click(screen.getByRole('button', { name: 'บันทึกการส่งมอบ' }))

  expect(await screen.findByText('บันทึกการส่งมอบแล้ว')).toBeTruthy()
  expect(client.calls).toHaveLength(1)
  const { requestId, form } = client.calls[0]
  expect(requestId).toBe(11)
  expect(form.get('borrowerFirstName')).toBe('สมชาย')
  expect(form.get('borrowerLastName')).toBe('ใจดี')
  expect(form.get('cardChecked')).toBe('true')
  expect(form.getAll('photos')).toHaveLength(2)
  expect(screen.getByText(/เวลาที่ส่งมอบ/)).toBeTruthy()
})
