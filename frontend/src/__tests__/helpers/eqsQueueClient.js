// client จำลองของหน้ารายละเอียดสำหรับ test ของ AC-EQS-17 (TC-EQS-17-*.test.jsx)
// endpoint ลำดับคำขอเป็นของ 002 (CR-005) requests เป็นรายการ หรือ Error ที่จะ throw
import { vi } from 'vitest'

// ผู้ขอ A ส่งก่อน (10:30 เวลาไทย) ผู้ขอ B ส่งตาม (11:15 เวลาไทย)
export const REQ_A = { requestId: 201, queuePosition: 1, submittedAt: '2026-11-10T03:30:00Z', requesterName: 'ผู้ขอ A' }
export const REQ_B = { requestId: 202, queuePosition: 2, submittedAt: '2026-11-10T04:15:00Z', requesterName: 'ผู้ขอ B' }

export const AVAILABLE_ITEM = { assetCode: 'CS-NB-0001', status: 'available', canRequestLoan: true }

export function queueClient({ item = AVAILABLE_ITEM, requests }) {
  return {
    getItem: vi.fn(async () => item),
    listItemRequests: vi.fn(async () => {
      if (requests instanceof Error) throw requests
      return requests
    }),
  }
}
