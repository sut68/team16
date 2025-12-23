
// src/validators/interview_validators.ts
import type { InterviewMode } from '@/interfaces/interview';

export interface InterviewRoundForm {
    name: string;
    scholarship_id: number | null;
    description: string;
    date: string; // YYYY-MM-DD
    end_date: string; // YYYY-MM-DD
    start_time: string; // HH:mm
    end_time: string; // HH:mm
    slot_duration: number;
    interviewer_ids: number[];
    interview_mode_id: number | null;
    location_id: number | null;
    meeting_link: string;
}

export interface ValidationResult {
  valid: boolean;
  errors: Record<string, any>;
}

function isValidURL(url = '') {
  try {
    new URL(url);
    return /^https?:\/\//i.test(String(url).trim());
  } catch (_) {
    return false;
  }
}

export function validateInterviewRoundForm(form: Partial<InterviewRoundForm>, modes: InterviewMode[] = []): ValidationResult {
  const errors: Record<string, any> = {};

  if (!form.name || form.name.trim().length < 3) {
    errors.name = 'โปรดระบุชื่อรอบสัมภาษณ์ (อย่างน้อย 3 ตัวอักษร)';
  }

  if (!form.scholarship_id) {
    errors.scholarship_id = 'โปรดเลือกทุนการศึกษา';
  }

  if (!form.date) {
    errors.date = 'โปรดระบุวันที่เริ่ม';
  }
  
  if (!form.end_date) {
    errors.end_date = 'โปรดระบุวันที่สิ้นสุด';
  }

  if (form.date && form.end_date && new Date(form.end_date) < new Date(form.date)) {
    errors.end_date = 'วันที่สิ้นสุดต้องไม่ก่อนวันที่เริ่ม';
  }

  if (!form.start_time) {
    errors.start_time = 'โปรดระบุเวลาเริ่ม';
  }

  if (!form.end_time) {
    errors.end_time = 'โปรดระบุเวลาสิ้นสุด';
  }

  if (form.date && form.end_date && form.start_time && form.end_time && form.date === form.end_date && form.end_time <= form.start_time) {
    errors.end_time = 'เวลาสิ้นสุดต้องหลังจากเวลาเริ่ม';
  }
  
  if (!form.slot_duration || form.slot_duration <= 0) {
    errors.slot_duration = 'โปรดกำหนดระยะเวลาที่ถูกต้อง';
  }
  
  if (!form.interviewer_ids || form.interviewer_ids.length === 0) {
    errors.interviewer_ids = 'โปรดเลือกกรรมการสัมภาษณ์อย่างน้อย 1 คน';
  }

  if (!form.interview_mode_id) {
    errors.interview_mode_id = 'โปรดเลือกรูปแบบการสัมภาษณ์';
  } else {
    const selectedMode = modes.find(m => m.ID === form.interview_mode_id);
    if (selectedMode) {
      if (selectedMode.name.toLowerCase() === 'onsite' && !form.location_id) {
        errors.location_id = 'โปรดเลือกสถานที่สำหรับ Onsite';
      }
      if (selectedMode.name.toLowerCase() === 'online' && (!form.meeting_link || !isValidURL(form.meeting_link))) {
        errors.meeting_link = 'โปรดระบุ URL สำหรับ Online Meeting ให้ถูกต้อง (ต้องมี http/https)';
      }
    }
  }

  return { valid: Object.keys(errors).length === 0, errors };
}
