import type { AdminProfileResponse } from './user'
import type { ApplicationDocumentResponse } from './application'

export interface ApprovalDecisionResponse {
  ID: number
  CreatedAt: string
  decision_at: string
  decision: string
  comment: string
  task_id: number
  admin_id: number
  admin_profile?: AdminProfileResponse
}

export interface ApprovalTaskResponse {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  status: string
  document_id: number
  application_document: Omit<ApplicationDocumentResponse, 'approval_tasks'>
  approval_decisions: ApprovalDecisionResponse[]
}
