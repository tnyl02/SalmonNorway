// test ตรวจว่าโครงหน้าจอพร้อมใช้ (ไม่ใช่ test ของ AC ใด)
import { render, screen } from '@testing-library/react'
import App from '../App.jsx'

test('โครงหน้าจอเปิดได้', () => {
  render(<App />)
  expect(screen.getByText('ระบบยืมอุปกรณ์ภาควิชาคอมพิวเตอร์')).toBeTruthy()
})
