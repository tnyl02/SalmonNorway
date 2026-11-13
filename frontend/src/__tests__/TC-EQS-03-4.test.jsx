// TC-EQS-03-4 (AC-EQS-03 ทางผิด) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { render, screen } from '@testing-library/react'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-03-4 ฐานข้อมูลไม่ตอบ แจ้งข้อผิดพลาดพร้อมปุ่มลองใหม่ ไม่ใช่ "ไม่พบอุปกรณ์"', async () => {
  // Given client จำลองตอบ 503 UPSTREAM_UNAVAILABLE (ฐานข้อมูลไม่ตอบ ไม่ใช่ไม่พบ)
  // When เปิดหน้ารายการอุปกรณ์
  render(<EquipmentList client={eqsFakeClient({ down: true })} />)

  // Then (1) แจ้งข้อผิดพลาดพร้อมปุ่ม "ลองใหม่"
  expect(await screen.findByRole('alert')).toBeTruthy()
  expect(screen.getByRole('button', { name: 'ลองใหม่' })).toBeTruthy()
  // Then (2) ไม่แสดง "ไม่พบอุปกรณ์"
  expect(screen.queryByText('ไม่พบอุปกรณ์')).toBeNull()
})
