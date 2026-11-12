// การเรียก API ของ 005-handover (plan.md ข้อ 4 ของ 005) ใช้ BASE และ ApiError ตัวเดียวกับ client.js
// แยกไฟล์เพื่อไม่แก้ client.js ที่หลายฟีเจอร์แก้อยู่ (ไม่ให้ชนตอน merge) ตอน test ส่ง client จำลองเข้าไปในหน้าจอแทน
// error ตอบเป็น { code, message, fields } ตาม CON-TECH-03 หน้าจอตัดสินใจจาก code และ fields
import { ApiError, BASE } from './client.js'
import { MOCK_TOKEN } from '../config.js'

async function handoverRequest(path, init = {}) {
  const res = await fetch(`${BASE}${path}`, {
    ...init,
    headers: { Authorization: `Bearer ${MOCK_TOKEN}`, ...(init.headers ?? {}) },
    cache: 'no-store',
  })
  const body = await res.json().catch(() => undefined)
  if (!res.ok) {
    const err = new ApiError(res.status, body?.code)
    err.fields = body?.fields ?? []
    err.serverMessage = body?.message // ข้อความภาษาไทยจากหลังบ้าน (CON-TECH-03)
    throw err
  }
  return body
}

export const handoverApi = {
  // FR-HND-01: คำขอที่อนุมัติแล้วและยังไม่ส่งมอบ -> [{ requestId, assetCode, modelName, requesterStudentId, submittedAt }]
  listRequests: () => handoverRequest('/v1/handover/requests'),

  // FR-HND-02: ยืนยันผู้ยืมที่เคาน์เตอร์ (ASM-HND-09) -> { result: 'verified' | 'mismatch' | 'unavailable', userId }
  verifyBorrower: (requestId, borrowerCredential) =>
    handoverRequest(`/v1/handover/requests/${encodeURIComponent(requestId)}/verify-borrower`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ borrowerCredential }),
    }),

  // FR-HND-02 ถึง 07: บันทึกการส่งมอบ formData มี borrowerCredential, borrowerFirstName, borrowerLastName, cardChecked, photos
  //   -> { loanId, assetCode, handedOverAt, dueAt }
  handover: (requestId, formData) =>
    handoverRequest(`/v1/handover/requests/${encodeURIComponent(requestId)}/handover`, { method: 'POST', body: formData }),
}
