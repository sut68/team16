import type { UpdateScreeningStatusPayload } from '@/interfaces/screening';

interface ValidationErrors {
  rejection_reason?: string;
}

export const validateScreeningStatusForm = (payload: UpdateScreeningStatusPayload) => {
  const errors: ValidationErrors = {};

  // ตรวจสอบเฉพาะกรณีสถานะคือ 3 (ไม่ผ่าน/Rejected)
  if (payload.status_screening_id === 3) {
    // 1. ตรวจสอบว่าต้องใส่ (Required)
    if (!payload.rejection_reason || payload.rejection_reason.trim() === '') {
      errors.rejection_reason = "กรุณาระบุเหตุผลที่ไม่ผ่านการคัดกรอง";
    } 
    // 2. ตรวจสอบความยาวสูงสุด
    else if (payload.rejection_reason.length > 100) {
      errors.rejection_reason = "เหตุผลต้องไม่เกิน 100 ตัวอักษร";
    }
  }

  return {
    valid: Object.keys(errors).length === 0,
    errors
  };
};