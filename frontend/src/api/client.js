// จุดเดียวที่หน้าจอใช้เรียก API หลังบ้าน (ตามสัญญา API ใน plan.md ข้อ 5 ของแต่ละฟีเจอร์)
// ตอน test ให้ส่ง client จำลองเข้าไปในหน้าจอแทน ไม่ต้องรันหลังบ้านจริง
// เรียกผ่าน /api (ดู proxy ใน vite.config.js) หลังบ้านต้องรันอยู่ที่ port 8000
// ฟังก์ชันของแต่ละ endpoint ให้ task ของฟีเจอร์นั้นเพิ่มเอง (เช่น T-06 ของ 001)
import { EQUIPMENT_REQUEST_TIMEOUT_MS, MOCK_TOKEN } from '../config.js'

export const BASE = import.meta.env.VITE_API_BASE ?? '/api'

// ApiError บอกว่าหลังบ้านตอบไม่สำเร็จ status คือ HTTP status code ส่วน code มาจาก Error contract ใน plan.md
// serverMessage คือข้อความภาษาไทยจากหลังบ้าน (CON-TECH-03) ใช้ชื่อเดียวกับ api/handover.js ของ 005
// หน้าจอแสดง serverMessage ได้ แต่ตัดสินใจจาก code เท่านั้น
//
// body คือ body ทั้งก้อนที่หลังบ้านส่งมา เพิ่มโดย 002 T-20 (แก้ไฟล์กลางเพราะ 002)
// CON-TECH-03 อนุญาตให้บางรหัสแนบช่องเพิ่มได้ เช่น VALIDATION_FAILED แนบ fields
// และ ITEM_NOT_AVAILABLE แนบ status ของอุปกรณ์ ซึ่ง 002 ต้องใช้ตาม FR-BRQ-03 และ FR-BRQ-01
// ฟีเจอร์อื่นไม่ต้องเปลี่ยนอะไร code กับ serverMessage ยังทำงานเหมือนเดิมทุกประการ
export class ApiError extends Error {
  constructor(status, code, serverMessage, body) {
    super(`API ${status} ${code ?? ''}`.trim())
    this.status = status
    this.code = code
    this.serverMessage = serverMessage
    this.body = body
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
// timeoutMs (ถ้าส่งมา) คือเวลารอหลังบ้านทั้งหมดรวมการอ่าน body เกินแล้วยกเลิก request (FR-EQS-08, Q-08, T-26)
// fetch ถูกยกเลิกจะ throw AbortError ซึ่งไม่ใช่ ApiError หน้าจอจึงแสดงข้อความเดิมพร้อมปุ่มลองใหม่ เหมือนต่อหลังบ้านไม่ได้
// ไม่ส่ง timeoutMs = รอจนกว่าหลังบ้านตอบ (พฤติกรรมเดิม ใช้กับ request ของ 002, 003)
async function getJSON(path, { timeoutMs } = {}) {
  const controller = timeoutMs ? new AbortController() : null
  const timer = controller ? setTimeout(() => controller.abort(), timeoutMs) : null
  try {
    const res = await fetch(`${BASE}${path}`, {
      headers: authHeaders(),
      cache: 'no-store',
      ...(controller ? { signal: controller.signal } : {}),
    })
    if (!res.ok) throw await failed(res)
    return await res.json()
  } finally {
    if (timer) clearTimeout(timer)
  }
}

// EQUIPMENT_TIMEOUT ใช้กับ request ของ 001 เท่านั้น (FR-EQS-08, Q-08, ทีมตัดสิน 2569-11-12 ตอนทำ T-26)
const EQUIPMENT_TIMEOUT = { timeoutMs: EQUIPMENT_REQUEST_TIMEOUT_MS }

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

// postForm ส่ง multipart/form-data พร้อม token เพิ่มโดย 002 T-19 สำหรับไฟล์แนบของการยืมค้างคืน (FR-BRQ-10 ข้อ ก)
// ห้ามตั้ง Content-Type เอง ต้องให้เบราว์เซอร์ใส่ boundary ให้ ไม่อย่างนั้นหลังบ้านอ่าน multipart ไม่ออก
async function postForm(path, formData) {
  const res = await fetch(`${BASE}${path}`, {
    method: 'POST',
    headers: authHeaders(),
    cache: 'no-store',
    body: formData,
  })
  if (!res.ok) throw await failed(res)
  return res.json()
}

// loanRequestForm แปลงข้อมูลฟอร์มเป็น FormData ตามชื่อช่องที่หลังบ้านรับ (plan.md ข้อ 4 ของ 002)
function loanRequestForm(input, attachment) {
  const form = new FormData()
  for (const [key, value] of Object.entries(input)) {
    if (value !== undefined && value !== null && value !== '') form.append(key, value)
  }
  form.append('attachment', attachment)
  return form
}

// apiError อ่าน code จาก Error contract ใน plan.md (ถ้าอ่านไม่ได้ code เป็น undefined)
async function apiError(res) {
  let body
  try {
    body = await res.json()
  } catch {
    body = undefined
  }
  return new ApiError(res.status, body?.code, body?.message, body)
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

// ---- 006-return (specs/006-return/plan.md ข้อ 4) ----
// returnsRequest คืนค่าเป็น { ok, status, data } เสมอ หน้าจอของ 006 ตัดสินใจจาก data.code (T-06)
// ใช้ token จาก session จำลองของ 001 (authHeaders) และแจ้ง 401 ให้ App แสดงหน้าล็อกอินจำลองเหมือนฟีเจอร์อื่น (IF-IDP-01)
// ต่อหลังบ้านไม่ได้ -> status 0, code UPSTREAM_UNAVAILABLE
async function returnsRequest(path, options = {}) {
  try {
    const res = await fetch(`${BASE}${path}`, {
      ...options,
      cache: 'no-store',
      headers: { ...authHeaders(), ...(options.headers ?? {}) },
    })
    const data = await res.json().catch(() => ({}))
    if (res.status === 401) {
      const err = new ApiError(401, data?.code, data?.message)
      unauthorizedListeners.forEach((listener) => listener(err))
    }
    return { ok: res.ok, status: res.status, data }
  } catch {
    return { ok: false, status: 0, data: { code: 'UPSTREAM_UNAVAILABLE' } }
  }
}
// ---- จบ 006-return ----

export const api = {
  // 001 FR-EQS-01, FR-EQS-02: รายการรุ่น 1 แถว = 1 รุ่น ตามตัวกรอง
  //   -> { items: [{ modelId, name, category, availableCount, totalCount }], total }
  listModels: (filter) => getJSON(`/v1/equipment/models${filterQuery(filter)}`, EQUIPMENT_TIMEOUT),

  // 001 FR-EQS-01, FR-EQS-02: อุปกรณ์รายชิ้นของรุ่น ตามตัวกรองเดียวกัน -> { items: [{ assetCode, status }], total }
  listItems: (modelId, filter) =>
    getJSON(`/v1/equipment/models/${encodeURIComponent(modelId)}/items${filterQuery(filter)}`, EQUIPMENT_TIMEOUT),

  // 001 FR-EQS-04: รายละเอียดอุปกรณ์รายชิ้น -> { assetCode, status, dueAt? } (dueAt มีเฉพาะถูกยืม)
  getItem: (assetCode) => getJSON(`/v1/equipment/items/${encodeURIComponent(assetCode)}`, EQUIPMENT_TIMEOUT),

  // 001 FR-EQS-11 (T-38): คำขอยืมที่ค้างอยู่ของอุปกรณ์ชิ้นนี้ เรียก endpoint ของ 002 (FR-BRQ-05 ก, สัญญาใน CR-005)
  //   -> [{ requestId, queuePosition, submittedAt, requesterName }] หรือ throw ApiError(403 FORBIDDEN_ROLE) เมื่อไม่ใช่เจ้าหน้าที่
  listItemRequests: (assetCode) => getJSON(`/v1/equipment/items/${encodeURIComponent(assetCode)}/requests`),

  // 004 FR-APV-01 (ก)(ข), FR-APV-03 (ก): คำขอยืมปกติและค้างคืนที่รอพิจารณาของภาควิชาเจ้าหน้าที่ เรียงตามเวลาที่ส่ง
  //   -> { items: [{ requestId, requesterName, modelName, assetCode, purpose, purposeOther?,
  //                  startAt, endAt, submittedAt, loanType, isFirst, attachmentUrl? }], loadedAt }
  getApprovalRequests: () => getJSON('/v1/approvals/requests'),

  // 004 FR-APV-01 (ค): อนุมัติคำขอ ส่ง loadedAt ที่ได้จาก getApprovalRequests
  //   -> { requestId, status, decidedAt, rejectedRequestIds }
  approveRequest: (requestId, loadedAt) =>
    postJSON(`/v1/approvals/requests/${encodeURIComponent(requestId)}/approve`, { loadedAt }),

  // 004 FR-APV-01 (ง): ปฏิเสธคำขอพร้อมเหตุผล -> { requestId, status, decidedAt, rejectedRequestIds }
  rejectRequest: (requestId, reason, loadedAt) =>
    postJSON(`/v1/approvals/requests/${encodeURIComponent(requestId)}/reject`, { reason, loadedAt }),

  // 004 FR-APV-03 (ก), Q-APV-41: เปิดไฟล์ใบบันทึกข้อความผ่าน endpoint ของ 002 (attachmentUrl จาก getApprovalRequests)
  //   ลิงก์ <a href> ธรรมดาส่ง token ไม่ได้ (endpoint อยู่หลัง auth.Middleware) จึงดึงไฟล์พร้อม token แล้วเปิดในแท็บใหม่
  //   เปิดแท็บก่อน await เพื่อไม่ให้เบราว์เซอร์บล็อก popup / ไม่สำเร็จปิดแท็บแล้ว throw ApiError ให้หน้าจอแสดงข้อความ
  //   (ทีมเลือกเอง ไม่ได้มาจาก spec เจ้าของ 004 อนุมัติให้แก้ไฟล์นี้ 2569-11-13)
  openAttachment: async (attachmentUrl) => {
    const win = window.open('', '_blank')
    try {
      const res = await fetch(`${BASE}${attachmentUrl}`, { headers: authHeaders(), cache: 'no-store' })
      if (!res.ok) throw await failed(res)
      const url = URL.createObjectURL(await res.blob())
      if (win) win.location.href = url
      else window.open(url, '_blank')
    } catch (err) {
      win?.close()
      throw err
    }
  },

  // 004 FR-APV-02 (ง): Email ที่ส่งไม่สำเร็จของภาควิชาเจ้าหน้าที่
  //   -> { items: [{ emailId, requestId, kind, requesterName, modelName, assetCode, attempts, lastError, lastAttemptAt }] }
  getFailedEmails: () => getJSON('/v1/approvals/failed-emails'),

  // 004 FR-APV-02 (ง): ส่ง Email ซ้ำโดยไม่พิจารณาใหม่ -> { emailId, status } (pending = เข้าคิวแล้ว)
  resendEmail: (emailId) => postJSON(`/v1/approvals/emails/${encodeURIComponent(emailId)}/resend`, {}),

  // 002 FR-BRQ-04: ส่งคำขอยืม 1 ใบ (plan.md ข้อ 4 ของ 002)
  //   in: { assetCode, requesterName, studentCode, loanType, purpose, purposeOther?, startAt, endAt }
  //   attachment: ไฟล์ใบบันทึกข้อความของการยืมค้างคืน ส่งเป็น multipart เมื่อมีไฟล์ (FR-BRQ-10 ข้อ ก)
  //   -> { requestId, status, submittedAt } หรือ throw ApiError(403 FORBIDDEN_DEPARTMENT,
  //      409 ITEM_NOT_AVAILABLE / DUPLICATE_REQUEST / HAS_OVERDUE_LOAN, 422 VALIDATION_FAILED, 503)
  createLoanRequest: (input, attachment) =>
    attachment ? postForm('/v1/loan-requests', loanRequestForm(input, attachment)) : postJSON('/v1/loan-requests', input),

  // 003 FR-CXL-06: คำขอที่ยังยกเลิกได้ของผู้ใช้ที่เข้าสู่ระบบ เรียงใหม่ไปเก่า (plan.md ข้อ 4 ของ 003)
  //   -> { items: [{ requestId, modelName, assetCode, submittedAt }] }
  listMyPendingRequests: () => getJSON('/v1/my/pending-requests'),

  // 003 FR-CXL-01, FR-CXL-08: ยกเลิกคำขอพร้อมเหตุผล (plan.md ข้อ 4 ของ 003)
  //   -> สำเร็จ หรือ throw ApiError(403 FORBIDDEN_ROLE, 404 REQUEST_NOT_FOUND, 409 ALREADY_APPROVED / NOT_CANCELLABLE, 422 VALIDATION_FAILED, 503)
  cancelRequest: (requestId, reason) =>
    postJSON(`/v1/borrow-requests/${encodeURIComponent(requestId)}/cancel`, { reason }),

  // ---- 006-return: บันทึกการรับคืน (เฉพาะเจ้าหน้าที่) คืนค่า { ok, status, data } ----
  returns: {
    // FR-RET-01(ข), FR-RET-06: การยืมที่ยังไม่คืนของภาควิชาเจ้าหน้าที่
    openLoans: () => returnsRequest('/v1/returns/open-loans'),

    // FR-RET-01(ก): หาการยืมที่ยังไม่คืนจากรหัสครุภัณฑ์
    lookup: (assetCode) => returnsRequest(`/v1/returns/lookup?assetCode=${encodeURIComponent(assetCode)}`),

    // FR-RET-06: รายละเอียดการยืม เวลาคืน สภาพ คืนเกินกำหนด
    loanDetail: (loanId) => returnsRequest(`/v1/returns/loans/${encodeURIComponent(loanId)}`),

    // FR-RET-02: บันทึกการรับคืน รหัสครุภัณฑ์และสภาพ ไม่ส่งเวลา (เวลาคืนมาจากระบบ ASM-04)
    record: ({ assetCode, condition }) => {
      const form = new FormData()
      form.append('assetCode', assetCode)
      form.append('condition', condition)
      return returnsRequest('/v1/returns', { method: 'POST', body: form })
    },
  },
}
