import type { CreateNewsPostPayload } from '@/interfaces/news_post';

interface ValidationErrors {
  title?: string;
  file_path?: string;
  post_detail?: string;
  scholarship_id?: string;
  status_news_id?: string;
}

export const validateNewsPostForm = (payload: CreateNewsPostPayload) => {
  const errors: ValidationErrors = {};

  // ===== หัวข้อข่าว =====
  if (!payload.title || payload.title.trim() === '') {
    errors.title = 'กรุณากรอกหัวข้อข่าว';
  }

  // ===== เนื้อหาข่าว =====
  if (!payload.post_detail || payload.post_detail.trim() === '') {
    errors.post_detail = 'กรุณากรอกเนื้อหาข่าว';
  }

  // ===== รูปภาพหน้าปก =====
  if (!payload.file_path) {
    errors.file_path = 'กรุณาอัปโหลดรูปภาพหน้าปกข่าว';
  }

  // ===== ทุนการศึกษา ===== (Optional - ไม่บังคับ)
  // ข่าวสารอาจไม่จำเป็นต้องผูกกับทุนการศึกษา

  // ===== สถานะข่าว =====
  const validStatuses = [1, 2, 3, 4, 5];
  if (payload.status_news_id === undefined) {
    errors.status_news_id = 'กรุณาเลือกสถานะข่าว';
  } else if (!validStatuses.includes(payload.status_news_id)) {
    errors.status_news_id = 'สถานะข่าวไม่ถูกต้อง';
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors
  };
};
