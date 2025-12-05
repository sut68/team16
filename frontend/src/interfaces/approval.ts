import { type AdminProfileResponse } from './user'
import { type ApplicationDocumentResponse, type ApplicationResponse } from './application'
import { type ApprovalRequirementResponse } from './scholarship'

// Response from backend
// ----------------------

export interface ApprovalDecisionResponse {
  ID: number
  decision_at: string
  decision: string
  comment: string
  task_id: number
}

export interface ApprovalTaskResponse {
  ID: number
  status: string
  admin_id: number
  admin_profile: AdminProfileResponse
  document_id: number
  application_document: ApplicationDocumentResponse
  application_id: number
  application: ApplicationResponse
  requirement_id: number
  approval_requirement: ApprovalRequirementResponse
  approval_decisions: ApprovalDecisionResponse[]
}
