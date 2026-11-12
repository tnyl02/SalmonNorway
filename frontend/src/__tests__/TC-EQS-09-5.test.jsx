// TC-EQS-09-5 (AC-EQS-09 ทางผิด) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { safePageKey } from '../App.jsx'

test('TC-EQS-09-5 returnTo เป็นที่อยู่นอกระบบ หลังเข้าสู่ระบบไม่ไปที่นั่น กลับหน้าในระบบแทน', () => {
  // Given returnTo เป็นที่อยู่นอกระบบ
  const returnTo = 'https://example.com'

  // When เข้าสู่ระบบสำเร็จ (App ใช้ safePageKey ตัดสินหน้าที่จะกลับไป)
  const page = safePageKey(returnTo)

  // Then ไม่ไปที่อยู่นอกระบบ กลับหน้าในระบบแทน (plan ข้อ 6.5)
  expect(page).not.toBe(returnTo)
  expect(['equipment', 'my-pending']).toContain(page)
})
