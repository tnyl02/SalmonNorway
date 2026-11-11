// AC-EQS-05 (FR-EQS-06, NFR-DATA-01) ส่วนสถานะเปลี่ยนเป็น "ว่าง"
//   Given เจ้าหน้าที่บันทึกการรับคืนใน UC-09 สภาพปกติ
//   When  นักศึกษารีเฟรชหน้ารายการภายใน 30 วินาที
//   Then  สถานะเปลี่ยนเป็น "ว่าง"
// ส่วน "ปุ่มส่งคำขอยืมกลับมาใช้งานได้" รอ Q-16 (T-34) / หน้าจอรีเฟรชเองไหม รอ Q-18 (T-35)
import { render, screen } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import { api } from '../api/client.js'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

afterEach(() => vi.unstubAllGlobals())

// fetch จำลองที่ตอบตามสถานะปัจจุบันของ "ฐานข้อมูล" และจดว่าเรียกด้วย option อะไร
function stubFetch(db) {
  const calls = []
  vi.stubGlobal('fetch', async (url, opts) => {
    calls.push({ url, opts })
    return new Response(JSON.stringify(db.item), { status: 200, headers: { 'Content-Type': 'application/json' } })
  })
  return calls
}

test('AC-EQS-05 รับคืนแล้ว กดรีเฟรชเห็นสถานะ "ว่าง" ไม่ใช่ข้อมูลเก่า', async () => {
  // ก่อนรับคืน: ชิ้นนี้ถูกยืมอยู่
  const db = { item: { assetCode: 'CS-NB-0004', status: 'borrowed', dueAt: '2026-11-10T09:00:00Z' } }
  const calls = stubFetch(db)
  const first = render(<EquipmentDetail assetCode="CS-NB-0004" client={api} />)
  expect(await screen.findByText('ถูกยืม')).toBeTruthy()

  // Given: UC-09 บันทึกการรับคืนสภาพปกติ สถานะในฐานข้อมูลเป็นว่าง
  db.item = { assetCode: 'CS-NB-0004', status: 'available' }

  // When: นักศึกษารีเฟรชหน้า (เปิดหน้าใหม่)
  first.unmount()
  render(<EquipmentDetail assetCode="CS-NB-0004" client={api} />)

  // Then: สถานะเป็น "ว่าง" และไม่มีกำหนดคืนเก่าค้าง
  expect(await screen.findByText('ว่าง')).toBeTruthy()
  expect(screen.queryByText('ถูกยืม')).toBeNull()
  expect(screen.queryByText('กำหนดคืน')).toBeNull()

  // Then: ถามหลังบ้านใหม่ทุกครั้ง และบอกเบราว์เซอร์ไม่ให้ใช้ cache (NFR-DATA-01)
  expect(calls).toHaveLength(2)
  for (const c of calls) expect(c.opts.cache).toBe('no-store')
})
