// load test ของฟีเจอร์ 003 ยกเลิกคำขอยืม (T-29, NFR-PERF-CXL-01, AC-CXL-17)
//   Given ผู้ใช้พร้อมกัน 50 คน
//   When  เปิดหน้า "คำขอที่รออนุมัติของฉัน" และยกเลิกคำขอ
//   Then  p95 ของเวลาตอบสนองของทั้ง 2 อย่าง <= 15 วินาที
//
// ต้องโหลด backend/db/seed/seed_003_load.sql ก่อน (คำขอ 50001-51000) และเปิดหลังบ้านไว้
// รันด้วย k6 ในเครื่อง:   k6 run tests/load/cancel_request.js
// รันด้วย k6 ใน Docker:   Get-Content tests\load\cancel_request.js -Raw | docker run --rm -i -e BASE_URL=http://host.docker.internal:8000 grafana/k6 run -
//
// k6 เป็นเครื่องมือที่ทีมเลือกเอง (plan.md ข้อ 2) ไม่ได้มาจาก spec
import http from 'k6/http'
import { check } from 'k6'

// ที่อยู่หลังบ้าน (ไม่มี /api เพราะเรียกหลังบ้านตรง ไม่ผ่าน proxy ของ vite)
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8000'

// ค่าจาก backend/internal/config/config.go: CancelPerfP95 = 15 วินาที (NFR-PERF-CXL-01) เปลี่ยนที่ config.go แล้วต้องแก้ตรงนี้ให้ตรงกัน
const P95_LIMIT_MS = Number(__ENV.P95_LIMIT_MS || 15000)

// ผู้ใช้พร้อมกัน 50 คน (AC-CXL-17) แต่ละคนยกเลิกคำขอของตัวเองคนละ 20 ใบ ไม่แย่งใบเดียวกัน รวม 1,000 ใบ = คำขอใน seed_003_load.sql
const VUS = 50
const PER_VU = 20
const FIRST_ID = 50001
const SPLIT_ID = 50500 // 50001-50500 ของ student-001, 50501-51000 ของ student-002

export const options = {
  scenarios: {
    cancel: { executor: 'per-vu-iterations', vus: VUS, iterations: PER_VU, maxDuration: '5m' },
  },
  thresholds: {
    'http_req_duration{name:list}': [`p(95)<${P95_LIMIT_MS}`], // เปิดหน้า "คำขอที่รออนุมัติของฉัน" (FR-CXL-06)
    'http_req_duration{name:cancel}': [`p(95)<${P95_LIMIT_MS}`], // ยกเลิกคำขอ (FR-CXL-01)
    checks: ['rate==1.0'], // ทุก request ต้องได้ผลถูก ไม่ใช่เร็วเพราะ error
  },
}

export default function () {
  // VU ที่ n (เริ่ม 1) รอบที่ i (เริ่ม 0) ยกเลิกคำขอ FIRST_ID + (n-1)*PER_VU + i
  const id = FIRST_ID + (__VU - 1) * PER_VU + __ITER
  const token = id <= SPLIT_ID ? 'mock-student-001' : 'mock-student-002'
  const auth = { headers: { Authorization: `Bearer ${token}` } }

  const list = http.get(`${BASE_URL}/v1/my/pending-requests`, { ...auth, tags: { name: 'list' } })
  check(list, { 'list 200': (r) => r.status === 200 })

  const res = http.post(
    `${BASE_URL}/v1/borrow-requests/${id}/cancel`,
    JSON.stringify({ reason: 'load test' }),
    { headers: { ...auth.headers, 'Content-Type': 'application/json' }, tags: { name: 'cancel' } },
  )
  check(res, { 'cancel 200': (r) => r.status === 200 })
}
