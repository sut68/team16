import type { AdminProfileResponse } from './user'
import type { ApplicationResponse, ApplicationDocumentResponse } from './application'
import type { ApprovalRequirementResponse } from './scholarship'

export interface ApprovalDecisionResponse {
  ID: number
  CreatedAt: string
  decision_at: string
  decision: string
  comment: string
  task_id: number
}

export interface ApprovalTaskResponse {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  status: string
  admin_id: number
  admin_profile?: AdminProfileResponse
  document_id: number
  application_document: ApplicationDocumentResponse
  application_id: number
  application: ApplicationResponse
  requirement_id: number
  approval_requirement: ApprovalRequirementResponse
  approval_decisions: ApprovalDecisionResponse[]
  round?: string
  submission_date?: string
}
