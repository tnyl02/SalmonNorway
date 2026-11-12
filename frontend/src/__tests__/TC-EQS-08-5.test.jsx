// TC-EQS-08-5 (AC-EQS-08 ขอบ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, within } from '@testing-library/react'
import { ApiError } from '../api/client.js'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-08-5 โหลดรอบแรกสำเร็จ รอบถัดไปตอบ 503 ไม่แสดงรายการเก่าค้าง', async () => {
  // Given client จำลองโหลดรอบแรกสำเร็จ (เห็นรุ่น A) แล้วรอบถัดไปตอบ 503
  const normal = eqsFakeClient()
  let calls = 0
  const client = {
    ...normal,
    listModels: async (filter) => {
      calls++
      if (calls > 1) throw new ApiError(503, 'UPSTREAM_UNAVAILABLE', 'ดึงข้อมูลไม่สำเร็จ กรุณาลองใหม่อีกครั้ง')
      return normal.listModels(filter)
    },
  }
  render(<EquipmentList client={client} />)
  await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })

  // When เปลี่ยนตัวกรองให้หน้าโหลดใหม่
  fireEvent.click(within(screen.getByRole('group', { name: 'สถานะ' })).getByRole('button', { name: 'ว่าง' }))

  // Then (1) แสดงข้อความผิดพลาดพร้อมปุ่ม "ลองใหม่"
  expect(await screen.findByRole('alert')).toBeTruthy()
  expect(screen.getByRole('button', { name: 'ลองใหม่' })).toBeTruthy()
  // Then (2) แถวรุ่น A ของรอบแรกไม่ค้างอยู่บนหน้าจอ (FR-EQS-08)
  expect(screen.queryByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeNull()
})
