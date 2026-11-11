// 006 T-02 เสร็จเมื่อ: ค่าใน config.js ตรงกับ backend/internal/config/config.go (FR-RET-02, FR-RET-05)
import {
  RETURN_PHOTO_MIN,
  RETURN_PHOTO_MAX,
  RETURN_PHOTO_MAX_BYTES,
  RETURN_PHOTO_TYPES,
  RETURN_CONDITIONS,
} from '../config.js'

test('ค่ารูปของการรับคืนตรงกับ spec (FR-RET-05)', () => {
  expect(RETURN_PHOTO_MIN).toBe(1)
  expect(RETURN_PHOTO_MAX).toBe(5)
  expect(RETURN_PHOTO_MAX_BYTES).toBe(5 * 1024 * 1024)
  expect(RETURN_PHOTO_TYPES).toEqual(['image/jpeg', 'image/png'])
})

test('สภาพตอนรับคืนมี 2 ค่า ไม่มีสูญหาย (FR-RET-02)', () => {
  expect(RETURN_CONDITIONS.map((c) => c.value)).toEqual(['normal', 'damaged'])
  expect(RETURN_CONDITIONS.map((c) => c.label)).toEqual(['ปกติ', 'ชำรุด'])
})
