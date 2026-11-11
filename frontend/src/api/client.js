// จุดเดียวที่หน้าจอใช้เรียก API หลังบ้าน (ตามสัญญา API ใน plan.md ของแต่ละฟีเจอร์)
// ตอน test ให้ส่ง client จำลองเข้าไปในหน้าจอแทน ไม่ต้องรันหลังบ้านจริง
// เรียกผ่าน /api (ดู proxy ใน vite.config.js) หลังบ้านต้องรันอยู่ที่ port 8000
// ฟังก์ชันของแต่ละ endpoint ให้ task ของฟีเจอร์นั้นเพิ่มเอง
export const BASE = import.meta.env.VITE_API_BASE ?? '/api'

// TOKEN_KEY คือ key ใน localStorage ที่เก็บ token ของ auth จำลอง (IF-IDP-01 รอ Q-SH-03)
// ทีมเลือกเอง 2569-11-11 (006 T-06): หน้าล็อกอินจำลองของ 001 (T-21) เป็นผู้เขียนค่านี้
// ระหว่างที่ยังไม่มีหน้านั้น ตั้งเองใน console ของเบราว์เซอร์ เช่น localStorage.setItem('token', 'mock-staff-001')
export const TOKEN_KEY = 'token'

// authHeaders ใส่ Authorization: Bearer <token> ทุก request (IF-IDP-01, NFR-SEC-01)
function authHeaders() {
  let token = null
  try {
    token = window.localStorage.getItem(TOKEN_KEY)
  } catch {
    token = null
  }
  return token ? { Authorization: `Bearer ${token}` } : {}
}

// request เรียกหลังบ้านแล้วคืน { ok, status, data } เสมอ หน้าจอไม่ต้องจัดการ exception เอง
// ต่อหลังบ้านไม่ได้ -> status 0, code UPSTREAM_UNAVAILABLE
async function request(path, options = {}) {
  try {
    const res = await fetch(`${BASE}${path}`, {
      ...options,
      headers: { ...authHeaders(), ...(options.headers ?? {}) },
    })
    const data = await res.json().catch(() => ({}))
    return { ok: res.ok, status: res.status, data }
  } catch {
    return { ok: false, status: 0, data: { code: 'UPSTREAM_UNAVAILABLE' } }
  }
}

export const api = {
  // ---- 006-return (specs/006-return/plan.md ข้อ 4) ----
  returns: {
    // lookup หาการยืมที่ยังไม่คืนจากรหัสครุภัณฑ์ (FR-RET-01(ก))
    lookup: (assetCode) =>
      request(`/v1/returns/lookup?assetCode=${encodeURIComponent(assetCode)}`),

    // record บันทึกการรับคืน: รหัสครุภัณฑ์ สภาพ และรูป (FR-RET-02, FR-RET-05) ไม่ส่งเวลา (เวลาคืนมาจากระบบ ASM-04)
    record: ({ assetCode, condition, photos }) => {
      const form = new FormData()
      form.append('assetCode', assetCode)
      form.append('condition', condition)
      for (const p of photos) form.append('photos', p)
      return request('/v1/returns', { method: 'POST', body: form })
    },
  },
}
