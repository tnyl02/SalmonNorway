import { useEffect, useState } from 'react'
import { api } from '../../api/client.js'
import {
  APPROVAL_ADVISOR_NAME_MAX_LENGTH,
  APPROVAL_ADVISOR_NAME_REQUIRED,
  APPROVAL_ADVISOR_NAME_TOO_LONG,
  APPROVAL_ATTACHMENT_LINK,
  APPROVAL_EMPTY_LIST,
  APPROVAL_LOAN_TYPE_LABELS,
  APPROVAL_MESSAGES,
  APPROVAL_REASON_MAX_LENGTH,
  APPROVAL_REASON_TOO_LONG,
  PURPOSE_OPTIONS,
  PURPOSE_OTHER,
  TIME_ZONE,
} from '../../config.js'

// รหัสภายในหน้าจอสำหรับเหตุผลยาวเกิน (Q-APV-28) หลังบ้านตอบ VALIDATION_FAILED เหมือนกรณีไม่ระบุเหตุผล แต่ข้อความต่างกัน
const REASON_TOO_LONG = 'REASON_TOO_LONG'
// รหัสภายในหน้าจอของชื่ออาจารย์ (FR-APV-03 ช, Q-APV-44) หลังบ้านตอบ VALIDATION_FAILED ช่อง advisorName ทั้ง 2 กรณี
const ADVISOR_REQUIRED = 'ADVISOR_REQUIRED'
const ADVISOR_TOO_LONG = 'ADVISOR_TOO_LONG'
// ค่า loanType ของคำขอค้างคืนที่หลังบ้านส่งมา (config.LoanTypeOvernight ของ 002) ใช้ค่าเดียวกับแบบฟอร์มของ 002
const LOAN_TYPE_OVERNIGHT = 'overnight'
// ข้อความของรหัสภายในหน้าจอ ข้อความมาจาก config.js (AGENTS.md กฎข้อ 9)
const LOCAL_MESSAGES = {
  [REASON_TOO_LONG]: APPROVAL_REASON_TOO_LONG,
  [ADVISOR_REQUIRED]: APPROVAL_ADVISOR_NAME_REQUIRED,
  [ADVISOR_TOO_LONG]: APPROVAL_ADVISOR_NAME_TOO_LONG,
}

// ตรวจชื่ออาจารย์แบบเดียวกับหลังบ้าน ว่างหรือมีแต่ช่องว่างไม่นับ ไม่เกิน 100 นับตาม code point (Q-APV-44) คืนรหัสภายในหรือ null
// รับชื่อที่ตัดช่องว่างหน้าหลังแล้ว (FR-APV-03 ช v3.1, rtm F-07, T-51)
function advisorError(name) {
  if (name.trim() === '') return ADVISOR_REQUIRED
  if ([...name].length > APPROVAL_ADVISOR_NAME_MAX_LENGTH) return ADVISOR_TOO_LONG
  return null
}

// แสดงวันเวลาตามเขตเวลา Asia/Bangkok (glossary)
const dateTime = new Intl.DateTimeFormat('th-TH', { timeZone: TIME_ZONE, dateStyle: 'medium', timeStyle: 'short' })
const fmt = (iso) => dateTime.format(new Date(iso))

// วัตถุประสงค์: ข้อความของ 3 ค่าจาก config ถ้าเป็น "อื่น ๆ" แสดงข้อความที่ผู้ขอพิมพ์ (FR-BRQ-02, FR-APV-01 ก)
function purposeLabel(r) {
  if (r.purpose === PURPOSE_OTHER && r.purposeOther) return r.purposeOther
  return PURPOSE_OPTIONS.find((o) => o.value === r.purpose)?.label ?? r.purpose
}

// หน้ารายการคำขอยืมที่รอพิจารณาของภาควิชาตนเอง ทั้งยืมปกติและค้างคืน (FR-APV-01 ก, FR-APV-03 ก, AC-APV-11)
// คอลัมน์ประเภทการยืม และลิงก์เปิดใบบันทึกข้อความเฉพาะคำขอค้างคืนที่มีไฟล์ (T-43, Q-APV-41)
// เรียงตามเวลาที่ส่งจากเก่าไปใหม่ตามที่หลังบ้านส่งมา ไม่เรียงใหม่ที่หน้าจอ (DOM-REQ-01)
// แสดงเฉพาะ 4 ข้อมูลตาม FR-APV-01 (ก) ไม่มีคอลัมน์เวลาที่ส่ง (rtm.md F-05 ทีมตัดสินแก้โค้ด 2569-11-13)
// ปุ่มอนุมัติ (T-09) / ปฏิเสธพร้อมช่องเหตุผล (T-11) / หน้าจอไม่มี mockup ทีมออกแบบเอง (plan.md ข้อ 4)
// คำขอค้างคืนทุกแถวมีช่องชื่ออาจารย์และปุ่มอนุมัติ เพราะไม่มีคิว (FR-APV-03 ข ค, T-45)
// client ส่งเข้ามาได้ ตอน test ใช้ client จำลองแทนหลังบ้านจริง
export default function ApproveRequests({ client = api }) {
  const [list, setList] = useState(null) // { items, loadedAt } จาก GET /v1/approvals/requests
  const [errorCode, setErrorCode] = useState(null) // โหลดรายการไม่สำเร็จ
  const [actionError, setActionError] = useState(null) // อนุมัติไม่สำเร็จ แสดงเหนือรายการ รายการยังอยู่
  const [busyId, setBusyId] = useState(null) // คำขอที่กำลังส่งอยู่ กันกดซ้ำ
  const [rejectingId, setRejectingId] = useState(null) // คำขอที่เปิดช่องเหตุผลการปฏิเสธอยู่ (T-11)
  const [reason, setReason] = useState('')
  const [advisorNames, setAdvisorNames] = useState({}) // ชื่ออาจารย์ที่พิมพ์ของแต่ละคำขอค้างคืน { [requestId]: string } (T-45)

  // ส่งผลการพิจารณา 1 คำขอ ส่ง loadedAt ของรายการที่เห็นอยู่ (plan.md ข้อ 4)
  // สำเร็จแล้วเอาคำขอนั้น และคำขอที่หลังบ้านปฏิเสธตาม (FR-APV-01 จ) ออกจากรายการ
  // ไม่สำเร็จแสดงข้อความตามรหัส (FR-APV-01 ซ, CON-TECH-03) และไม่เปลี่ยนรายการ
  async function decide(requestId, send) {
    setBusyId(requestId)
    setActionError(null)
    try {
      const res = await send()
      const gone = new Set([res.requestId, ...(res.rejectedRequestIds ?? [])])
      setList((cur) => ({ ...cur, items: cur.items.filter((r) => !gone.has(r.requestId)) }))
      setRejectingId(null)
      setReason('')
    } catch (err) {
      setActionError(err?.code ?? 'UPSTREAM_UNAVAILABLE')
    } finally {
      setBusyId(null)
    }
  }

  // เปิดใบบันทึกข้อความของคำขอค้างคืน (FR-APV-03 ก, Q-APV-41) ไม่สำเร็จแสดงข้อความตามรหัสเหนือรายการ
  const openAttachment = (r) => {
    setActionError(null)
    Promise.resolve(client.openAttachment(r.attachmentUrl)).catch((err) =>
      setActionError(err?.code ?? 'UPSTREAM_UNAVAILABLE'),
    )
  }

  // อนุมัติ (FR-APV-01 ค, AC-APV-01)
  const approve = (requestId) => decide(requestId, () => client.approveRequest(requestId, list.loadedAt))

  // อนุมัติคำขอค้างคืนพร้อมชื่ออาจารย์ (FR-APV-03 ข ช, AC-APV-22)
  // ชื่อว่างหรือยาวเกินไม่ส่งไปหลังบ้าน แสดงข้อความเลย (Q-APV-44) หลังบ้านตรวจซ้ำเสมอ (T-44)
  // หลังบ้านตอบ VALIDATION_FAILED ช่อง advisorName แปลงเป็นรหัสภายใน เพราะข้อความต่างจากของช่องเหตุผล
  const approveOvernight = (requestId) => {
    // ตัดช่องว่างหน้าหลังก่อนตรวจและก่อนส่ง นับความยาวหลังตัดแบบเดียวกับหลังบ้าน (FR-APV-03 ช v3.1, T-51)
    const name = (advisorNames[requestId] ?? '').trim()
    const invalid = advisorError(name)
    if (invalid) {
      setActionError(invalid)
      return
    }
    return decide(requestId, () =>
      client.approveRequest(requestId, list.loadedAt, name).catch((err) => {
        if (err?.code === 'VALIDATION_FAILED' && err?.body?.fields?.includes('advisorName')) {
          throw { code: [...name].length > APPROVAL_ADVISOR_NAME_MAX_LENGTH ? ADVISOR_TOO_LONG : ADVISOR_REQUIRED }
        }
        throw err
      }),
    )
  }

  // ปฏิเสธพร้อมเหตุผลที่เจ้าหน้าที่พิมพ์ (FR-APV-01 ง, AC-APV-02)
  // ไม่ระบุเหตุผล (ว่างหรือมีแต่ช่องว่าง) ไม่ส่งไปหลังบ้าน แสดงข้อความของ VALIDATION_FAILED เลย คำขอยังอยู่ (AC-APV-06, T-13)
  // หลังบ้านตรวจซ้ำเสมอ (T-12) ถ้าหลังบ้านตอบ VALIDATION_FAILED ก็แสดงข้อความเดียวกันผ่าน decide
  const reject = (requestId) => {
    if (reason.trim() === '') {
      setActionError('VALIDATION_FAILED')
      return
    }
    // ยาวเกินกำหนด นับตาม code point แบบเดียวกับหลังบ้าน (FR-APV-01 ง, Q-APV-28) ไม่ส่งไปหลังบ้าน
    if ([...reason].length > APPROVAL_REASON_MAX_LENGTH) {
      setActionError(REASON_TOO_LONG)
      return
    }
    return decide(requestId, () => client.rejectRequest(requestId, reason, list.loadedAt))
  }

  useEffect(() => {
    let active = true
    client
      .getApprovalRequests()
      .then((res) => active && setList(res))
      .catch((err) => active && setErrorCode(err?.code ?? 'UPSTREAM_UNAVAILABLE'))
    return () => {
      active = false
    }
  }, [client])

  if (errorCode) {
    // ข้อความเลือกตามรหัส ไม่ใช่ตามข้อความจากหลังบ้าน (CON-TECH-03)
    const message = APPROVAL_MESSAGES[errorCode] ?? APPROVAL_MESSAGES.UPSTREAM_UNAVAILABLE
    return (
      <p role="alert" className="mt-4 text-red-700">
        {message}
      </p>
    )
  }
  if (list == null) return <p className="mt-4 text-slate-600">กำลังโหลด...</p>

  return (
    <section className="mt-4 grid gap-2">
      <h2 className="text-lg font-semibold text-slate-900">คำขอยืมที่รอพิจารณา</h2>
      {actionError && (
        <p role="alert" className="text-red-700">
          {LOCAL_MESSAGES[actionError] ?? APPROVAL_MESSAGES[actionError] ?? APPROVAL_MESSAGES.UPSTREAM_UNAVAILABLE}
        </p>
      )}
      {/* ไม่มีคำขอที่รอพิจารณา แสดงข้อความแทนตาราง (FR-APV-01 ก, Q-APV-31) */}
      {list.items.length === 0 && <p className="text-slate-600">{APPROVAL_EMPTY_LIST}</p>}
      {list.items.length > 0 && (
      <table className="w-full text-left text-sm">
        <thead>
          <tr className="border-b">
            <th className="py-1">ผู้ขอ</th>
            <th className="py-1">อุปกรณ์</th>
            <th className="py-1">วัตถุประสงค์</th>
            <th className="py-1">ช่วงเวลาที่ต้องการยืม</th>
            <th className="py-1">ประเภทการยืม</th>
            <th className="py-1"></th>
          </tr>
        </thead>
        <tbody>
          {list.items.map((r) => (
            <tr key={r.requestId} className="border-b">
              <td className="py-1">{r.requesterName}</td>
              <td className="py-1">
                {r.modelName} <span className="text-slate-500">({r.assetCode})</span>
              </td>
              <td className="py-1">{purposeLabel(r)}</td>
              <td className="py-1">
                {fmt(r.startAt)} - {fmt(r.endAt)}
              </td>
              <td className="py-1">
                {APPROVAL_LOAN_TYPE_LABELS[r.loanType] ?? r.loanType}
                {r.attachmentUrl && (
                  <>
                    {' '}
                    <a
                      href={r.attachmentUrl}
                      onClick={(e) => {
                        e.preventDefault()
                        openAttachment(r)
                      }}
                      className="text-teal-700 underline"
                    >
                      {APPROVAL_ATTACHMENT_LINK}
                    </a>
                  </>
                )}
              </td>
              <td className="py-1">
                {/* ปุ่มอนุมัติเฉพาะคำขอแรกของชิ้นนั้น (FR-APV-01 ข, AC-APV-07, T-18) หลังบ้านตรวจซ้ำเสมอ (T-17)
                    แถวอื่นปฏิเสธได้อย่างเดียว */}
                {/* คำขอค้างคืน: ช่องชื่ออาจารย์และปุ่มอนุมัติทุกแถว ไม่มีคิว (FR-APV-03 ข ค ช, T-45) ทีมออกแบบเอง ไม่มี mockup */}
                {r.loanType === LOAN_TYPE_OVERNIGHT && (
                  <span className="mb-1 grid gap-1">
                    <label className="grid gap-1">
                      ชื่ออาจารย์ที่เซ็นอนุมัติ
                      <input
                        type="text"
                        value={advisorNames[r.requestId] ?? ''}
                        onChange={(e) => setAdvisorNames((cur) => ({ ...cur, [r.requestId]: e.target.value }))}
                        className="rounded border border-slate-300 p-1"
                      />
                    </label>
                    <span>
                      <button
                        type="button"
                        onClick={() => approveOvernight(r.requestId)}
                        disabled={busyId != null}
                        className="rounded bg-teal-700 px-3 py-1 text-white disabled:opacity-50"
                      >
                        อนุมัติ
                      </button>
                    </span>
                  </span>
                )}
                {r.loanType !== LOAN_TYPE_OVERNIGHT && r.isFirst && (
                  <>
                    <button
                      type="button"
                      onClick={() => approve(r.requestId)}
                      disabled={busyId != null}
                      className="rounded bg-teal-700 px-3 py-1 text-white disabled:opacity-50"
                    >
                      อนุมัติ
                    </button>{' '}
                  </>
                )}
                {rejectingId === r.requestId ? (
                  // ช่องเหตุผลการปฏิเสธ (FR-APV-01 ง) ทีมออกแบบเอง ไม่มี mockup
                  <span className="mt-1 grid gap-1">
                    <label className="grid gap-1">
                      เหตุผลในการปฏิเสธ
                      <textarea
                        value={reason}
                        onChange={(e) => setReason(e.target.value)}
                        rows={2}
                        className="rounded border border-slate-300 p-1"
                      />
                    </label>
                    <span className="flex gap-1">
                      <button
                        type="button"
                        onClick={() => reject(r.requestId)}
                        disabled={busyId != null}
                        className="rounded bg-red-700 px-3 py-1 text-white disabled:opacity-50"
                      >
                        ยืนยันการปฏิเสธ
                      </button>
                      <button
                        type="button"
                        onClick={() => {
                          setRejectingId(null)
                          setReason('')
                        }}
                        className="rounded border border-slate-400 px-3 py-1"
                      >
                        ยกเลิก
                      </button>
                    </span>
                  </span>
                ) : (
                  <button
                    type="button"
                    onClick={() => {
                      setRejectingId(r.requestId)
                      setReason('')
                    }}
                    disabled={busyId != null}
                    className="rounded border border-red-700 px-3 py-1 text-red-700 disabled:opacity-50"
                  >
                    ปฏิเสธ
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      )}
    </section>
  )
}
