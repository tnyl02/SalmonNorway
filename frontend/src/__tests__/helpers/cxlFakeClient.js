// ตัวช่วยของ test หน้าจอ 003 จาก specs/003-cancel-request/test-cases.md (ไม่ใช่ test เอง แก้ได้ถ้าโครงโค้ดเปลี่ยน)
//
// สิ่งที่ test เหล่านี้คาดจากหน้าจอของ 003 (plan.md ข้อ 4, tasks.md T-05, T-07):
//   - หน้า frontend/src/pages/MyPendingRequests/MyPendingRequests.jsx รับ prop `client` แบบเดียวกับหน้าของ 001
//   - client.listMyPendingRequests() -> { items: [{ requestId, modelName, assetCode, submittedAt }] } (เรียงมาจากหลังบ้านแล้ว)
//   - client.cancelRequest(requestId, reason) -> สำเร็จ หรือ throw ApiError(status, code) จาก api/client.js
//   - แต่ละคำขอเป็น 1 รายการ (role listitem) มีปุ่ม "ยกเลิก"
//   - กล่องยืนยันเป็น role dialog มีช่องเหตุผล (textbox) ตัวนับ "x/200" ปุ่ม "ยืนยัน" และ "ไม่ยืนยัน" (คำจาก FR-CXL-07, AC-CXL-09)
// ก่อน T-05 และ T-07 เสร็จ test เหล่านี้จะไม่ผ่านเพราะยังไม่มีไฟล์หน้า เป็นเรื่องปกติ ไม่ใช่บั๊ก
import { vi } from 'vitest'

export const REQ_A1 = {
  requestId: 101,
  modelName: 'Notebook ตัวอย่างรุ่น A',
  assetCode: 'CS-NB-0001',
  submittedAt: '2026-11-10T03:30:00Z', // 10:30 เวลา Asia/Bangkok
}

export const REQ_A2 = {
  requestId: 102,
  modelName: 'Projector ตัวอย่างรุ่น B',
  assetCode: 'CS-PJ-0001',
  submittedAt: '2026-11-10T02:00:00Z', // 09:00 เวลา Asia/Bangkok
}

// client จำลองที่จำสถานะ: ยกเลิกสำเร็จแล้วคำขอหายจากรายการ ทั้งกรณีหน้าจอโหลดรายการใหม่หรือเอาแถวออกเอง
export function fakeClient(items) {
  let current = [...items]
  return {
    listMyPendingRequests: vi.fn(async () => ({ items: [...current] })),
    cancelRequest: vi.fn(async (requestId) => {
      current = current.filter((r) => r.requestId !== requestId)
      return { requestId, status: 'cancelled' }
    }),
  }
}
