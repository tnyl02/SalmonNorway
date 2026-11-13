// client จำลองของหน้าอนุมัติ 004 ใช้ร่วมใน test จาก specs/004-approve-request/test-cases.md
// รูปแบบเดียวกับคำตอบ GET /v1/approvals/requests (plan.md ข้อ 4)
export const LOADED_AT = '2026-11-12T03:00:00Z'

// 2 คำขอเรียงตามเวลาที่ส่งแล้ว (หลังบ้านกรองและเรียงให้ TC-APV-11-1)
// startAt 02:00Z = 09:00 น. และ endAt 09:00Z = 16:00 น. เวลาไทย
export const REQ_1 = {
  requestId: 21, requesterName: 'นักศึกษาตัวอย่าง ก', studentCode: 'student-001', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)',
  assetCode: 'AP-TEST-0002', purpose: 'study', startAt: '2026-11-12T02:00:00Z', endAt: '2026-11-12T09:00:00Z',
  submittedAt: '2026-11-12T01:30:00Z', isFirst: true,
}
export const REQ_2 = {
  requestId: 22, requesterName: 'นักศึกษาตัวอย่าง ข', studentCode: 'student-002', modelName: 'อุปกรณ์ทดลองอนุมัติ (004)',
  assetCode: 'AP-TEST-0003', purpose: 'other', purposeOther: 'ถ่ายวิดีโอโครงงาน', startAt: '2026-11-12T02:00:00Z',
  endAt: '2026-11-12T09:00:00Z', submittedAt: '2026-11-12T01:40:00Z', isFirst: true,
}

// fakeClient คืนรายการตาม items ที่ให้ และนับการเรียกอนุมัติและปฏิเสธไว้ใน calls
export function fakeClient(items = [REQ_1, REQ_2]) {
  const calls = { approve: [], reject: [] }
  return {
    calls,
    getApprovalRequests: async () => ({ loadedAt: LOADED_AT, items }),
    approveRequest: async (requestId, loadedAt) => {
      calls.approve.push({ requestId, loadedAt })
      return { requestId, status: 'approved', decidedAt: '2026-11-12T03:05:00Z', rejectedRequestIds: [] }
    },
    rejectRequest: async (requestId, reason, loadedAt) => {
      calls.reject.push({ requestId, reason, loadedAt })
      return { requestId, status: 'rejected', decidedAt: '2026-11-12T03:05:00Z', rejectedRequestIds: [] }
    },
  }
}
