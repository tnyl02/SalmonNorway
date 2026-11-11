import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import HandoverApp from './pages/Handover/HandoverApp.jsx'

// จุดเริ่มของหน้าบันทึกการส่งมอบ (005-handover, UC-08) เปิดจาก frontend/handover.html
// เฉพาะเจ้าหน้าที่ (CON-HND-01) หลังบ้านตรวจสิทธิ์ทุก request
createRoot(document.getElementById('root')).render(
  <StrictMode>
    <HandoverApp />
  </StrictMode>,
)
