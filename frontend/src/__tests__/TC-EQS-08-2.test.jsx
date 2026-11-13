// TC-EQS-08-2 (AC-EQS-08 ปกติ) จาก specs/001-equipment-status/test-cases.md ห้ามแก้ให้ผ่าน
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { ApiError } from '../api/client.js'
import EquipmentList from '../pages/EquipmentList/EquipmentList.jsx'
import { eqsFakeClient } from './helpers/eqsFakeClient.js'

test('TC-EQS-08-2 ฐานข้อมูลไม่ตอบ แสดงข้อผิดพลาดพร้อมปุ่มลองใหม่ กดลองใหม่หลังฐานข้อมูลกลับมาแล้วเห็นรายการ', async () => {
  // Given client จำลองที่ฐานข้อมูลไม่ตอบ (ตอบ 503) สลับกลับเป็นปกติได้
  const state = { down: true, listCalls: 0 }
  const normal = eqsFakeClient()
  const client = {
    ...normal,
    listModels: async (filter) => {
      state.listCalls++
      if (state.down) throw new ApiError(503, 'UPSTREAM_UNAVAILABLE', 'ดึงข้อมูลไม่สำเร็จ กรุณาลองใหม่อีกครั้ง')
      return normal.listModels(filter)
    },
  }

  // When เปิดหน้ารายการ
  render(<EquipmentList client={client} />)

  // Then (1) แสดงข้อความผิดพลาด
  expect(await screen.findByRole('alert')).toBeTruthy()
  // Then (2) มีปุ่ม "ลองใหม่"
  const retry = screen.getByRole('button', { name: 'ลองใหม่' })
  // Then (3) ไม่มีแถวรุ่นใด
  expect(screen.queryByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeNull()
  expect(screen.queryByRole('button', { name: /Projector ตัวอย่างรุ่น B/ })).toBeNull()

  // When ทำให้ client กลับมาปกติ แล้วกด "ลองใหม่"
  state.down = false
  const before = state.listCalls
  fireEvent.click(retry)

  // Then (4) เรียกหลังบ้านอีกครั้งและแสดงรายการรุ่น A
  expect(await screen.findByRole('button', { name: /Notebook ตัวอย่างรุ่น A/ })).toBeTruthy()
  await waitFor(() => expect(state.listCalls).toBeGreaterThan(before))
})
