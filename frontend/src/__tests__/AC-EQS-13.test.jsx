// AC-EQS-13 (FR-EQS-10, DOM-ACCESS-01) ใช้ client จำลองแทนหลังบ้าน
//   Given นักศึกษาภาควิชา A ยืนยันตัวตนแล้ว และอุปกรณ์เป็นของภาควิชา B
//   When  เรียกรายละเอียดอุปกรณ์ของภาควิชา B โดยตรง
//   Then  แจ้งว่าไม่มีสิทธิ์เข้าถึง และไม่มีข้อมูลอุปกรณ์ของภาควิชา B ถูกส่งกลับ
// การกรองรายการให้เห็นเฉพาะภาควิชาตนเอง ตรวจที่หลังบ้าน (TestAC_EQS_13_OtherDepartment403)
import { render, screen } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import { api, ApiError } from '../api/client.js'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

afterEach(() => vi.unstubAllGlobals())

const FORBIDDEN_MESSAGE = 'ไม่มีสิทธิ์เข้าถึงอุปกรณ์ของภาควิชาอื่น'

test('AC-EQS-13 เปิดรายละเอียดอุปกรณ์ภาควิชาอื่น แจ้งไม่มีสิทธิ์ และไม่แสดงข้อมูลอุปกรณ์', async () => {
  // Given: หลังบ้านตอบ 403 FORBIDDEN_DEPARTMENT ตาม CON-TECH-03
  const client = {
    getItem: async () => {
      throw new ApiError(403, 'FORBIDDEN_DEPARTMENT', FORBIDDEN_MESSAGE)
    },
  }

  // When: เรียกรายละเอียดของภาควิชา B โดยตรง
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)

  // Then: แจ้งไม่มีสิทธิ์ และไม่มีสถานะหรือกำหนดคืนของอุปกรณ์
  const alert = await screen.findByRole('alert')
  expect(alert.textContent).toBe(FORBIDDEN_MESSAGE)
  for (const label of ['ว่าง', 'ถูกยืม', 'ไม่พร้อมใช้งาน', 'กำหนดคืน', 'สถานะ']) {
    expect(screen.queryByText(label)).toBeNull()
  }
})

test('client.js เก็บ message จากหลังบ้านไว้ใน serverMessage (CON-TECH-03)', async () => {
  vi.stubGlobal('fetch', async () =>
    new Response(JSON.stringify({ code: 'FORBIDDEN_DEPARTMENT', message: FORBIDDEN_MESSAGE }), {
      status: 403,
      headers: { 'Content-Type': 'application/json' },
    }),
  )
  const err = await api.getItem('CS-NB-0001').catch((e) => e)
  expect(err).toBeInstanceOf(ApiError)
  expect(err.status).toBe(403)
  expect(err.code).toBe('FORBIDDEN_DEPARTMENT')
  expect(err.serverMessage).toBe(FORBIDDEN_MESSAGE)
})

test('error ที่ไม่มีข้อความจากหลังบ้าน ใช้ข้อความเดิม "ดึงข้อมูลอุปกรณ์ไม่สำเร็จ"', async () => {
  const client = {
    getItem: async () => {
      throw new TypeError('Failed to fetch')
    },
  }
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)
  expect((await screen.findByRole('alert')).textContent).toBe('ดึงข้อมูลอุปกรณ์ไม่สำเร็จ')
})
