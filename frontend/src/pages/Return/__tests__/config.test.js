// 006 T-02 เสร็จเมื่อ: ค่าใน config.js ตรงกับ backend/internal/config/config.go (FR-RET-02) อยู่ในโฟลเดอร์ของ 006
// ค่ารูป (FR-RET-05) ถอดออก 2569-11-11 เพราะ FR-RET-05 ยกเลิก
import { RETURN_CONDITIONS } from '../../../config.js'

test('สภาพตอนรับคืนมี 2 ค่า ไม่มีสูญหาย (FR-RET-02)', () => {
  expect(RETURN_CONDITIONS.map((c) => c.value)).toEqual(['normal', 'damaged'])
  expect(RETURN_CONDITIONS.map((c) => c.label)).toEqual(['ปกติ', 'ชำรุด'])
})
