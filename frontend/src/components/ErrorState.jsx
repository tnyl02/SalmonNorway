// แจ้งข้อผิดพลาดจากหลังบ้าน (CON-TECH-03, FR-EQS-08, FR-EQS-10)
// แสดงข้อความภาษาไทยที่หลังบ้านส่งมา (serverMessage) เช่น 403 FORBIDDEN_DEPARTMENT "ไม่มีสิทธิ์เข้าถึงอุปกรณ์ของภาควิชาอื่น"
// ถ้าไม่มีข้อความ (เช่น ต่อหลังบ้านไม่ได้เลย) ใช้ข้อความเดิมของหน้าจอ
// ไม่แสดงข้อมูลเก่าใด ๆ ปุ่มลองใหม่ของ FR-EQS-08 เพิ่มใน T-20
export const DEFAULT_ERROR_MESSAGE = 'ดึงข้อมูลอุปกรณ์ไม่สำเร็จ'

export default function ErrorState({ error }) {
  return (
    <p role="alert" className="text-red-700">
      {error?.serverMessage || DEFAULT_ERROR_MESSAGE}
    </p>
  )
}
