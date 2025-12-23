// src/validators/student_interview_validator.ts

export interface BookingPayload {
  slot_id: number | null | undefined;
  application_scholarship_id: number | null | undefined;
}

export interface ValidationResult {
  valid: boolean;
  errors: Record<string, string>;
}

export function validateBookingPayload(payload: Partial<BookingPayload>): ValidationResult {
  const errors: Record<string, string> = {};

  if (!payload.slot_id) {
    errors.slot_id = 'กรุณาเลือกช่วงเวลาที่ต้องการจอง';
  }
  
  if (!payload.application_scholarship_id) {
    errors.application_scholarship_id = 'ไม่พบใบสมัครสำหรับทุนการศึกษานี้';
  }

  return { valid: Object.keys(errors).length === 0, errors };
}
