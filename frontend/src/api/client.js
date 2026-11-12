// จุดเดียวที่หน้าจอใช้เรียก API หลังบ้าน (ตามสัญญา API ใน plan.md ข้อ 5 ของแต่ละฟีเจอร์)
// ตอน test ให้ส่ง client จำลองเข้าไปในหน้าจอแทน ไม่ต้องรันหลังบ้านจริง
// เรียกผ่าน /api (ดู proxy ใน vite.config.js) หลังบ้านต้องรันอยู่ที่ port 8000
// ฟังก์ชันของแต่ละ endpoint ให้ task ของฟีเจอร์นั้นเพิ่มเอง (เช่น T-06 ของ 001)
import { MOCK_TOKEN } from '../config.js'

export const BASE = import.meta.env.VITE_API_BASE ?? '/api'

// ApiError บอกว่าหลังบ้านตอบไม่สำเร็จ status คือ HTTP status code ส่วน code มาจาก Error contract ใน plan.md
// serverMessage คือข้อความภาษาไทยจากหลังบ้าน (CON-TECH-03) ใช้ชื่อเดียวกับ api/handover.js ของ 005
// หน้าจอแสดง serverMessage ได้ แต่ตัดสินใจจาก code เท่านั้น
export class ApiError extends Error {
  constructor(status, code, serverMessage) {
    super(`API ${status} ${code ?? ''}`.trim())
    this.status = status
    this.code = code
    this.serverMessage = serverMessage
  }
}

// ---- session จำลอง (IF-IDP-01, AC-EQS-09, T-21 ของ 001) ทีมเลือกเอง ไม่ได้มาจาก spec รอ Q-SH-03 ----
// token เก็บใน sessionStorage ของ browser ลบ token ทิ้ง = session หมดอายุ (request ถัดไปได้ 401)
// เปิดระบบครั้งแรกใช้ MOCK_TOKEN จาก config.js (VITE_MOCK_TOKEN หรือค่าเริ่มต้น) เหมือนเดิม ฟีเจอร์อื่นไม่ต้องเปลี่ยนวิธีลอง
const TOKEN_KEY = 'salmon.mockToken'
let memoryToken = null // ใช้แทนเมื่อ browser ไม่ให้ใช้ sessionStorage

function tokenStore() {
  try {
    return window.sessionStorage
  } catch {
    return null
  }
}

// getToken คืน token ปัจจุบัน null = ไม่มี session
export function getToken() {
  const s = tokenStore()
  return s ? s.getItem(TOKEN_KEY) : memoryToken
}

// setToken ตั้ง token หลังล็อกอินจำลอง ส่ง null = ลบ session
export function setToken(token) {
  const s = tokenStore()
  if (!s) {
    memoryToken = token
  } else if (token) {
    s.setItem(TOKEN_KEY, token)
  } else {
    s.removeItem(TOKEN_KEY)
  }
}

if (getToken() == null) setToken(MOCK_TOKEN)

function authHeaders() {
  const token = getToken()
  return token ? { Authorization: `Bearer ${token}` } : {}
}

// onUnauthorized ลงทะเบียนให้ถูกเรียกทุกครั้งที่หลังบ้านตอบ 401 (App ใช้แสดงหน้าล็อกอินจำลอง) คืนฟังก์ชันยกเลิก
const unauthorizedListeners = new Set()
export function onUnauthorized(listener) {
  unauthorizedListeners.add(listener)
  return () => unauthorizedListeners.delete(listener)
}

// failed แปลง response ที่ไม่สำเร็จเป็น ApiError และแจ้งผู้ฟังเมื่อเป็น 401 (AC-EQS-09)
async function failed(res) {
  const err = await apiError(res)
  if (res.status === 401) unauthorizedListeners.forEach((listener) => listener(err))
  return err
}

// ทุก request ต้องยืนยันตัวตน (IF-IDP-01, NFR-SEC-01) ส่ง token จาก session จำลอง
// cache: 'no-store' ให้เบราว์เซอร์ถามหลังบ้านทุกครั้ง ไม่ใช้ข้อมูลเก่า (NFR-DATA-01, plan.md ข้อ 6.2)
// หลังบ้านส่ง Cache-Control: no-store มาด้วย (platform.NoStore)
async function getJSON(path) {
  const res = await fetch(`${BASE}${path}`, {
    headers: authHeaders(),
    cache: 'no-store',
  })
  if (!res.ok) throw await failed(res)
  return res.json()
}

// postJSON ส่ง body เป็น JSON แบบเดียวกับ getJSON (token, no-store, ApiError) เพิ่มโดย 003 T-07 ฟีเจอร์อื่นใช้ร่วมได้
async function postJSON(path, body) {
  const res = await fetch(`${BASE}${path}`, {
    method: 'POST',
    headers: { ...authHeaders(), 'Content-Type': 'application/json' },
    cache: 'no-store',
    body: JSON.stringify(body),
  })
  if (!res.ok) throw await failed(res)
  return res.json()
}

// apiError อ่าน code จาก Error contract ใน plan.md (ถ้าอ่านไม่ได้ code เป็น undefined)
async function apiError(res) {
  let body
  try {
    body = await res.json()
  } catch {
    body = undefined
  }
  return new ApiError(res.status, body?.code, body?.message)
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

  // 003 FR-CXL-06: คำขอที่ยังยกเลิกได้ของผู้ใช้ที่เข้าสู่ระบบ เรียงใหม่ไปเก่า (plan.md ข้อ 4 ของ 003)
  //   -> { items: [{ requestId, modelName, assetCode, submittedAt }] }
  listMyPendingRequests: () => getJSON('/v1/my/pending-requests'),

  // 003 FR-CXL-01, FR-CXL-08: ยกเลิกคำขอพร้อมเหตุผล (plan.md ข้อ 4 ของ 003)
  //   -> สำเร็จ หรือ throw ApiError(400 REASON_*, 404 REQUEST_NOT_FOUND, 409 ALREADY_APPROVED / NOT_CANCELLABLE, 503)
  cancelRequest: (requestId, reason) =>
    postJSON(`/v1/borrow-requests/${encodeURIComponent(requestId)}/cancel`, { reason }),
}
