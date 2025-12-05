import type { StudentProfileResponse } from './user'
import type { ApprovalRequirementResponse } from './scholarship'
import type { ApprovalTaskResponse } from './approval'

// Response from backend
// ----------------------

export interface ApplicationDocumentResponse {
  ID: number
  file_name: string
  uploaded_by: string
  application_id: number
  requirement_id: number
  approval_requirement: ApprovalRequirementResponse
  approval_tasks: ApprovalTaskResponse[]
}

export interface ApplicationResponse {
  ID: number
  student_profile_id: number
  student_profile: StudentProfileResponse
  application_documents: ApplicationDocumentResponse[]
}
