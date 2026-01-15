import type { ApprovalTaskResponse } from '@/interfaces/approval' // ตรวจสอบ path interface ให้ถูกต้อง
import { https } from '@/services/api'

export const getApprovalTasks = async (): Promise<ApprovalTaskResponse[]> => {
  // axiosInstance จะ return AxiosResponse เราต้อง .data เพื่อเอาเนื้อหา
  // Add a cache-busting query parameter
  const cacheBust = new Date().getTime();
  const response = await https.get<ApprovalTaskResponse[]>(`/approval-tasks?t=${cacheBust}`);
  return response.data;
}

export const getApprovalTaskById = async (id: number): Promise<ApprovalTaskResponse> => {
  const response = await https.get<ApprovalTaskResponse>(`/approval-tasks/${id}`)
  return response.data
}

interface DecisionPayload {
  task_id: number
  decision: 'approve' | 'reject' | 'request-change'
  comment: string
  admin_id: number
}

export const makeApprovalDecision = async (payload: DecisionPayload) => {
  // 1. Create the decision record
  const decisionResponse = await https.post('/approval-decisions', {
    task_id: payload.task_id,
    decision: payload.decision,
    comment: payload.comment,
    admin_id: payload.admin_id,
  })

  if (decisionResponse.status !== 201) {
    throw new Error('Failed to create approval decision')
  }

  return decisionResponse.data
}