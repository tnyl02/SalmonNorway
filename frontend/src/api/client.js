// จุดเดียวที่หน้าจอใช้เรียก API หลังบ้าน (ตามสัญญา API ใน plan.md ข้อ 5 ของแต่ละฟีเจอร์)
// ตอน test ให้ส่ง client จำลองเข้าไปในหน้าจอแทน ไม่ต้องรันหลังบ้านจริง
// เรียกผ่าน /api (ดู proxy ใน vite.config.js) หลังบ้านต้องรันอยู่ที่ port 8000
// ฟังก์ชันของแต่ละ endpoint ให้ task ของฟีเจอร์นั้นเพิ่มเอง (เช่น T-06 ของ 001)
import { MOCK_TOKEN } from '../config.js'

export const BASE = import.meta.env.VITE_API_BASE ?? '/api'

// ApiError บอกว่าหลังบ้านตอบไม่สำเร็จ status คือ HTTP status code ส่วน code มาจาก Error contract ใน plan.md
export class ApiError extends Error {
  constructor(status, code) {
    super(`API ${status} ${code ?? ''}`.trim())
    this.status = status
    this.code = code
  }
}

// ทุก request ต้องยืนยันตัวตน (IF-IDP-01, NFR-SEC-01) ส่ง token จำลองจาก config.js ไปก่อน (T-21 จะเปลี่ยนเป็นหน้าล็อกอินจำลอง)
// cache: 'no-store' ให้เบราว์เซอร์ถามหลังบ้านทุกครั้ง ไม่ใช้ข้อมูลเก่า (NFR-DATA-01, plan.md ข้อ 6.2)
// หลังบ้านส่ง Cache-Control: no-store มาด้วย (platform.NoStore)
async function getJSON(path) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { Authorization: `Bearer ${MOCK_TOKEN}` },
    cache: 'no-store',
  })
  if (!res.ok) {
    let code
    try {
      code = (await res.json()).code
    } catch {
      code = undefined
    }
    throw new ApiError(res.status, code)
  }
  return res.json()
}

// 001 FR-EQS-02: แปลงตัวกรอง { categories: [...], statuses: [...] } เป็น query string
// แต่ละค่าส่งเป็นชื่อซ้ำ เช่น ?status=available&status=borrowed (รูปแบบที่หลังบ้านรับใน T-32)
function filterQuery(filter = {}) {
  const params = new URLSearchParams()
  for (const c of filter.categories ?? []) params.append('category', c)
  for (const s of filter.statuses ?? []) params.append('status', s)
  const qs = params.toString()
  return qs ? `?${qs}` : ''
}

export const api = {
  // 001 FR-EQS-01, FR-EQS-02: รายการรุ่น 1 แถว = 1 รุ่น ตามตัวกรอง
  //   -> { items: [{ modelId, name, category, availableCount, totalCount }], total }
  listModels: (filter) => getJSON(`/v1/equipment/models${filterQuery(filter)}`),

  // 001 FR-EQS-01, FR-EQS-02: อุปกรณ์รายชิ้นของรุ่น ตามตัวกรองเดียวกัน -> { items: [{ assetCode, status }], total }
  listItems: (modelId, filter) =>
    getJSON(`/v1/equipment/models/${encodeURIComponent(modelId)}/items${filterQuery(filter)}`),

  // 001 FR-EQS-04: รายละเอียดอุปกรณ์รายชิ้น -> { assetCode, status, dueAt? } (dueAt มีเฉพาะถูกยืม)
  getItem: (assetCode) => getJSON(`/v1/equipment/items/${encodeURIComponent(assetCode)}`),
}

// ---- 005-handover (plan.md ข้อ 4 ของ 005) ----
// error ตอบเป็น { code, message, fields } ตาม CON-TECH-03 หน้าจอตัดสินใจจาก code และ fields

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
    throw err
  }
  return body
}

export const handoverApi = {
  // 005 FR-HND-01: คำขอที่อนุมัติแล้วและยังไม่ส่งมอบ -> [{ requestId, assetCode, modelName, requesterStudentId, submittedAt }]
  listRequests: () => handoverRequest('/v1/handover/requests'),

  // 005 FR-HND-02: ยืนยันผู้ยืมที่เคาน์เตอร์ (ASM-HND-09) -> { result: 'verified' | 'mismatch' | 'unavailable', userId }
  verifyBorrower: (requestId, borrowerCredential) =>
    handoverRequest(`/v1/handover/requests/${encodeURIComponent(requestId)}/verify-borrower`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ borrowerCredential }),
    }),

  // 005 FR-HND-02 ถึง 07: บันทึกการส่งมอบ formData มี borrowerCredential, borrowerFirstName, borrowerLastName, cardChecked, photos
  //   -> { loanId, assetCode, handedOverAt, dueAt }
  handover: (requestId, formData) =>
    handoverRequest(`/v1/handover/requests/${encodeURIComponent(requestId)}/handover`, { method: 'POST', body: formData }),
}
