import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

// test หน้าจอรันใน jsdom (เบราว์เซอร์จำลอง) ไม่ต้องเปิดเบราว์เซอร์จริง
export default defineConfig({
  plugins: [react(), tailwindcss()],
  // หน้าจอเรียก /api/... แล้ว Vite ส่งต่อไปหลังบ้านที่ port 8000 ในเครื่องเดียวกัน
  // ทำให้ไม่ติดปัญหา CORS และใช้ได้ทั้ง Codespace บนเบราว์เซอร์และบน VS Code
  server: {
    host: true,
    port: 5173,
    // หลังบ้านใน Docker เปิด HTTPS ที่ :8443 (T-50 ของ 004): ตั้ง VITE_API_TARGET=https://localhost:8443
    // secure: false ยอมรับใบรับรองที่ลงนามเองของเครื่องพัฒนา ไม่ตั้ง = หลังบ้าน HTTP :8000 แบบเดิม
    proxy: {
      '/api': {
        target: process.env.VITE_API_TARGET ?? 'http://localhost:8000',
        secure: false,
        rewrite: (p) => p.replace(/^\/api/, ''),
      },
    },
  },
  test: { environment: 'jsdom', globals: true },
})
