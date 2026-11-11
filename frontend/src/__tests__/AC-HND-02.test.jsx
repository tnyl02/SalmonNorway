// AC-HND-02 (FR-HND-02, FR-HND-03) ฟอร์มบันทึกการส่งมอบ ใช้ client จำลองแทนหลังบ้าน
//   Given คำขอยืมแบบปกติที่ได้รับอนุมัติแล้ว Auth API ยืนยันตัวตนผู้ยืมสำเร็จ และผู้ยืมเป็นผู้ส่งคำขอ
//   When  เจ้าหน้าที่ยืนยันว่าตรวจบัตรแล้ว และบันทึกการส่งมอบ
//   Then  ระบบบันทึกการส่งมอบ พร้อมชื่อ นามสกุล รหัสนักศึกษาของผู้ยืม และวันที่และเวลาที่ส่งมอบ
//         ถ้าเจ้าหน้าที่ไม่ยืนยันการตรวจบัตร ระบบไม่บันทึกการส่งมอบ
// หลังบ้านตรวจการบันทึกจริงใน TestAC_HND_02_* ที่นี่ตรวจว่าหน้าจอยืนยันผู้ยืมก่อน (T-12) ส่งข้อมูลครบ และไม่ส่งเมื่อไม่ติ๊กตรวจบัตร (T-11)
import { fireEvent, render, screen } from '@testing-library/react'
import HandoverForm from '../pages/Handover/HandoverForm.jsx'

const request = { requestId: 11, assetCode: 'HO-TEST-0001', requesterStudentId: 'student-001' }

function fakeClient(verifyResult = 'verified') {
  const calls = []
  return {
    calls,
    verifyBorrower: async (requestId, credential) => {
      calls.push({ verify: true, requestId, credential })
      return { result: verifyResult, userId: verifyResult === 'unavailable' ? undefined : 'student-001' }
    },
    handover: async (requestId, form) => {
      calls.push({ requestId, form })
      return { loanId: 1, assetCode: 'HO-TEST-0001', handedOverAt: '2026-11-12T03:00:00Z', dueAt: '2026-11-12T09:00:00Z' }
    },
  }
}

async function verifyBorrower(expectedText) {
  fireEvent.change(screen.getByLabelText('บัญชีผู้ยืม'), { target: { value: 'mock-student-001' } })
  fireEvent.click(screen.getByRole('button', { name: 'ยืนยันตัวตนผู้ยืม' }))
  expect(await screen.findByText(new RegExp(expectedText))).toBeTruthy()
}

function fillBorrower() {
  fireEvent.change(screen.getByLabelText('ชื่อผู้ยืม (ตามบัตร)'), { target: { value: 'สมชาย' } })
  fireEvent.change(screen.getByLabelText('นามสกุลผู้ยืม (ตามบัตร)'), { target: { value: 'ใจดี' } })
  const photos = [new File(['a'], 'a.jpg', { type: 'image/jpeg' }), new File(['b'], 'b.jpg', { type: 'image/jpeg' })]
  fireEvent.change(screen.getByLabelText('รูปสภาพอุปกรณ์'), { target: { files: photos } })
}

const saveButton = () => screen.getByRole('button', { name: 'บันทึกการส่งมอบ' })
const handoverCalls = (client) => client.calls.filter((c) => !c.verify)

test('AC-HND-02 ยืนยันผู้ยืมแล้วแต่ไม่ติ๊กว่าตรวจบัตร ปุ่มบันทึกใช้ไม่ได้และไม่ส่งข้อมูล', async () => {
  const client = fakeClient()
  render(<HandoverForm request={request} client={client} />)
  await verifyBorrower('ยืนยันตัวตนผู้ยืมแล้ว')
  fillBorrower()

  expect(saveButton().disabled).toBe(true)
  fireEvent.submit(screen.getByRole('form', { name: 'บันทึกการส่งมอบ' }))
  expect(handoverCalls(client)).toHaveLength(0)
})

test('AC-HND-02 ยืนยันผู้ยืมและตรวจบัตรแล้ว ส่งบัญชีผู้ยืม ชื่อ นามสกุล ผลตรวจบัตร และรูป แล้วแสดงเวลาที่ส่งมอบ', async () => {
  const client = fakeClient()
  render(<HandoverForm request={request} client={client} />)
  await verifyBorrower('ยืนยันตัวตนผู้ยืมแล้ว')
  fillBorrower()
  fireEvent.click(screen.getByLabelText('ตรวจบัตรนักศึกษาหรือบัตรประชาชนของผู้ยืมแล้ว'))
  fireEvent.click(saveButton())

  expect(await screen.findByText('บันทึกการส่งมอบแล้ว')).toBeTruthy()
  const calls = handoverCalls(client)
  expect(calls).toHaveLength(1)
  const { requestId, form } = calls[0]
  expect(requestId).toBe(11)
  expect(form.get('borrowerCredential')).toBe('mock-student-001')
  expect(form.get('borrowerFirstName')).toBe('สมชาย')
  expect(form.get('borrowerLastName')).toBe('ใจดี')
  expect(form.get('cardChecked')).toBe('true')
  expect(form.getAll('photos')).toHaveLength(2)
  expect(screen.getByText(/เวลาที่ส่งมอบ/)).toBeTruthy()
})

test('T-12 ยังไม่ยืนยันผู้ยืม หรือผู้ยืมไม่ใช่ผู้ส่งคำขอ บันทึกไม่ได้ (FR-HND-02)', async () => {
  const client = fakeClient('mismatch')
  render(<HandoverForm request={request} client={client} />)
  fillBorrower()
  fireEvent.click(screen.getByLabelText('ตรวจบัตรนักศึกษาหรือบัตรประชาชนของผู้ยืมแล้ว'))
  expect(saveButton().disabled).toBe(true) // ยังไม่ยืนยันผู้ยืม

  await verifyBorrower('ผู้ยืมไม่ใช่ผู้ส่งคำขอนี้')
  expect(saveButton().disabled).toBe(true)
  expect(handoverCalls(client)).toHaveLength(0)
})
