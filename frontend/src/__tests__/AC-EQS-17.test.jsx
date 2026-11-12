// AC-EQS-17 (FR-EQS-11, DOM-REQ-01, DOM-PRIV-01) ใช้ client จำลองแทนหลังบ้าน (endpoint เป็นของ 002 ตามสัญญาใน CR-005)
//   Given อุปกรณ์ชิ้นหนึ่งสถานะ "ว่าง" มีคำขอค้างอยู่ 2 รายการ ผู้ขอ A ส่งก่อน ผู้ขอ B ส่งตาม
//   When  เจ้าหน้าที่เปิดดูรายละเอียด และนักศึกษาเปิดดูรายละเอียดเดียวกัน
//   Then  เจ้าหน้าที่เห็นคำขอ 2 รายการ ลำดับที่ 1 คือ A ลำดับที่ 2 คือ B พร้อมเวลาที่ส่งและชื่อผู้ขอ
//         ส่วนนักศึกษาไม่เห็นรายการคำขอและชื่อผู้ขอใด ๆ
// หลังบ้านของ 002 เป็นคนตัดสินบทบาท: เจ้าหน้าที่ได้ 200 นักศึกษาได้ 403 FORBIDDEN_ROLE (CON-TECH-03)
import { act, render, screen, within } from '@testing-library/react'
import { vi } from 'vitest'
import { ApiError } from '../api/client.js'
import EquipmentDetail from '../pages/EquipmentDetail/EquipmentDetail.jsx'

const ITEM = { assetCode: 'CS-NB-0001', status: 'available' }

// Given: A ส่งก่อน (10:30 เวลาไทย) B ส่งตาม (11:15 เวลาไทย)
const REQ_A = { requestId: 201, queuePosition: 1, submittedAt: '2026-11-10T03:30:00Z', requesterName: 'ผู้ขอ A' }
const REQ_B = { requestId: 202, queuePosition: 2, submittedAt: '2026-11-10T04:15:00Z', requesterName: 'ผู้ขอ B' }

const FORBIDDEN_ROLE = new ApiError(403, 'FORBIDDEN_ROLE', 'เฉพาะเจ้าหน้าที่เท่านั้นที่ดูข้อมูลนี้ได้')

function clientFor({ item = ITEM, requests }) {
  return {
    getItem: vi.fn(async () => item),
    listItemRequests: vi.fn(async () => {
      if (requests instanceof Error) throw requests
      return requests
    }),
  }
}

// รอให้ promise ของการโหลดคำขอจบก่อนตรวจว่า "ไม่มี" ส่วนนั้น
async function settle(client) {
  await screen.findByText('ว่าง')
  await act(async () => {})
  expect(client.listItemRequests).toHaveBeenCalledWith('CS-NB-0001')
}

test('AC-EQS-17 เจ้าหน้าที่เห็นคำขอ 2 รายการ ลำดับ 1 = A ลำดับ 2 = B พร้อมเวลาที่ส่งและชื่อผู้ขอ', async () => {
  // When: เจ้าหน้าที่เปิดดูรายละเอียด (หลังบ้านคืนลำดับ แม้ส่งมาสลับกัน หน้าจอต้องเรียงตามลำดับที่)
  const client = clientFor({ requests: [REQ_B, REQ_A] })
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)

  // Then
  const queue = await screen.findByRole('region', { name: 'คำขอยืมที่ค้างอยู่' })
  const rows = within(queue).getAllByRole('row').slice(1) // แถวแรกคือหัวตาราง
  expect(rows).toHaveLength(2)
  expect(within(rows[0]).getByText('1')).toBeTruthy()
  expect(within(rows[0]).getByText('2569-11-10 10:30')).toBeTruthy()
  expect(within(rows[0]).getByText('ผู้ขอ A')).toBeTruthy()
  expect(within(rows[1]).getByText('2')).toBeTruthy()
  expect(within(rows[1]).getByText('2569-11-10 11:15')).toBeTruthy()
  expect(within(rows[1]).getByText('ผู้ขอ B')).toBeTruthy()

  // FR-EQS-11: บอกว่าชื่อผู้ขอเป็นชื่อที่พิมพ์เอง และสถานะยังเป็น "ว่าง" (DOM-REQ-01)
  expect(within(queue).getByText(/ชื่อที่ผู้ขอพิมพ์เอง/)).toBeTruthy()
  expect(screen.getByText('ว่าง')).toBeTruthy()
})

test('AC-EQS-17 นักศึกษา (403 FORBIDDEN_ROLE) ไม่เห็นรายการคำขอและชื่อผู้ขอใด ๆ', async () => {
  // When: นักศึกษาเปิดดูรายละเอียดเดียวกัน
  const client = clientFor({ requests: FORBIDDEN_ROLE })
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)
  await settle(client)

  // Then: ไม่มีส่วนคำขอ ไม่มีชื่อผู้ขอ และไม่มีข้อความ error (DOM-PRIV-01)
  expect(screen.queryByRole('region', { name: 'คำขอยืมที่ค้างอยู่' })).toBeNull()
  expect(screen.queryByText('ผู้ขอ A')).toBeNull()
  expect(screen.queryByText('ผู้ขอ B')).toBeNull()
  expect(screen.queryByRole('alert')).toBeNull()
})

test('FR-EQS-11 002 ไม่ตอบ: สถานะและกำหนดคืนยังแสดง ไม่มีส่วนคำขอ (CR-005)', async () => {
  const borrowed = { assetCode: 'CS-NB-0001', status: 'borrowed', dueAt: '2026-11-10T09:00:00Z' }
  const client = clientFor({ item: borrowed, requests: new ApiError(503, 'UPSTREAM_UNAVAILABLE') })
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)

  expect(await screen.findByText('ถูกยืม')).toBeTruthy()
  await act(async () => {})
  expect(screen.getByText('2569-11-10 16:00')).toBeTruthy()
  expect(screen.queryByRole('region', { name: 'คำขอยืมที่ค้างอยู่' })).toBeNull()
  expect(screen.queryByRole('alert')).toBeNull()
})

test('FR-EQS-11 ไม่มีคำขอค้าง: ไม่แสดงส่วนนี้ (spec ไม่บอกข้อความ plan ข้อ 10)', async () => {
  const client = clientFor({ requests: [] })
  render(<EquipmentDetail assetCode="CS-NB-0001" client={client} />)
  await settle(client)

  expect(screen.queryByRole('region', { name: 'คำขอยืมที่ค้างอยู่' })).toBeNull()
})
