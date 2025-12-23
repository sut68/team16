// src/validators/approval_validators.ts

export interface ApprovalDecisionPayload {
  task_id: number;
  decision: 'approve' | 'reject' | 'request-change';
  comment: string;
  admin_id: number;
}

export interface ValidationResult {
  valid: boolean;
  errors: Record<string, any>;
}

export function validateApprovalDecision(form: Partial<ApprovalDecisionPayload>): ValidationResult {
  const errors: Record<string, any> = {};

  if (!form.decision || !['approve', 'reject', 'request-change'].includes(form.decision)) {
    errors.decision = 'โปรดเลือกการตัดสินใจ';
  }

  if ((form.decision === 'reject' || form.decision === 'request-change') && (!form.comment || form.comment.trim().length === 0)) {
    errors.comment = 'กรุณาระบุเหตุผลประกอบการพิจารณา';
  }

  if (!form.task_id || form.task_id <= 0) {
    errors.task_id = 'Task ID ไม่ถูกต้อง';
  }

  if (!form.admin_id || form.admin_id <= 0) {
    errors.admin_id = 'Admin ID ไม่ถูกต้อง';
  }

  return { valid: Object.keys(errors).length === 0, errors };
}
