import type { CreateNewsPostPayload } from '@/interfaces/news_post';

// กำหนดโครงสร้าง Error ให้สอดคล้องกับฟิลด์ใน Payload
interface ValidationErrors {
  title?: string;
  post_detail?: string;
  file_path?: string;
  scholarship_id?: string;
  status_news_id?: string;
}

export const validateNewsPostForm = (payload: CreateNewsPostPayload) => {
  const errors: ValidationErrors = {};

  // 1. ตรวจสอบหัวข้อข่าว (Title) - ตรงกับ Backend: stringlength(5|200)
  if (!payload.title || payload.title.trim() === '') {
    errors.title = "กรุณากรอกหัวข้อข่าว";
  } else if (payload.title.length < 5) {
    errors.title = "หัวข้อข่าวต้องมีอย่างน้อย 5 ตัวอักษร";
  } else if (payload.title.length > 200) {
    errors.title = "หัวข้อข่าวต้องไม่เกิน 200 ตัวอักษร";
  }

  // 2. ตรวจสอบเนื้อหาข่าว (Post Detail) - ตรงกับ Backend: stringlength(10|500)
  if (!payload.post_detail || payload.post_detail.trim() === '') {
    errors.post_detail = "กรุณากรอกเนื้อหาข่าว";
  } else if (payload.post_detail.length < 10) {
    errors.post_detail = "เนื้อหาข่าวต้องมีอย่างน้อย 10 ตัวอักษร";
  } else if (payload.post_detail.length > 500) {
    errors.post_detail = "เนื้อหาข่าวต้องไม่เกิน 500 ตัวอักษร";
  }

  // 3. ตรวจสอบรูปภาพ (File Path)
  // ใน CreateNewsPostPayload ของคุณ file_path เป็น File | null
  if (!payload.file_path) {
    errors.file_path = "กรุณาเลือกรูปภาพหน้าปกข่าว";
  }

  // 4. ตรวจสอบทุนการศึกษา (Scholarship ID)
  if (payload.scholarship_id === null || payload.scholarship_id === undefined) {
    errors.scholarship_id = "กรุณาเลือกทุนการศึกษาที่เกี่ยวข้อง";
  }

  // 5. ตรวจสอบสถานะข่าว (Status News ID) - ตรงกับ Backend: in(1|2|3|4|5)
  const validStatuses = [1, 2, 3, 4, 5];
  if (!payload.status_news_id || !validStatuses.includes(payload.status_news_id)) {
    errors.status_news_id = "กรุณาระบุสถานะข่าวให้ถูกต้อง";
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors
  };
};