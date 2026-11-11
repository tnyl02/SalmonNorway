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

// postJSON ส่ง body เป็น JSON แบบเดียวกับ getJSON (เพิ่มโดย 004 T-09 plan.md ข้อ 4)
// ไม่สำเร็จ throw ApiError ที่มี code จากหลังบ้าน (CON-TECH-03) ให้หน้าจอเลือกข้อความจากรหัส
async function postJSON(path, body) {
  const res = await fetch(`${BASE}${path}`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${MOCK_TOKEN}`, 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
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

  // 004 FR-APV-01 (ก)(ข): คำขอยืมปกติที่รอพิจารณาของภาควิชาเจ้าหน้าที่ เรียงตามเวลาที่ส่ง
  //   -> { items: [{ requestId, requesterName, studentCode, modelName, assetCode, purpose, purposeOther?,
  //                  startAt, endAt, submittedAt, isFirst }], loadedAt }
  getApprovalRequests: () => getJSON('/v1/approvals/requests'),

  // 004 FR-APV-01 (ค): อนุมัติคำขอ ส่ง loadedAt ที่ได้จาก getApprovalRequests
  //   -> { requestId, status, decidedAt, rejectedRequestIds }
  approveRequest: (requestId, loadedAt) =>
    postJSON(`/v1/approvals/requests/${encodeURIComponent(requestId)}/approve`, { loadedAt }),
}
