import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import ApproveRequests from './pages/ApproveRequest/ApproveRequests.jsx'

// จุดเริ่มของหน้าอนุมัติคำขอยืมปกติ (004-approve-request, UC-12) เปิดจาก frontend/approve.html
// เฉพาะเจ้าหน้าที่ (ASM-01) หลังบ้านตรวจบทบาทและภาควิชาทุก request (AC-APV-05, ASM-06)
// ลองในเครื่อง: VITE_MOCK_TOKEN=mock-staff-001 npm run dev แล้วเปิด /approve.html
createRoot(document.getElementById('root')).render(
  <StrictMode>
    <main className="mx-auto max-w-4xl p-4 sm:p-6">
      <ApproveRequests />
    </main>
  </StrictMode>,
)
