-- ย้อน 0011: คืนสภาพให้เหมือนก่อนรับคำขอยืมค้างคืน (CR-009, T-40)
-- คอลัมน์ advisor_name ถูกลบพร้อมข้อมูล ชื่ออาจารย์ที่บันทึกไว้แล้วจะหาย

DROP TABLE approval_step;
DROP FUNCTION approval_step_append_only();
ALTER TABLE approval_decision DROP COLUMN advisor_name;
