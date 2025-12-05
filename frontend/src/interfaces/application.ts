import type { StudentProfileResponse } from './user'
import type { ApprovalTaskResponse } from './approval'
import type { ApprovalRequirementResponse } from './scholarship'


export interface ApplicationDocumentResponse {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  file_name: string
  uploaded_by: string
  application_id: number
  requirement_id: number
  approval_requirement: ApprovalRequirementResponse
  approval_tasks: ApprovalTaskResponse[]
}

export interface ApplicationResponse {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  student_profile_id: number
  student_profile: StudentProfileResponse
  application_documents: ApplicationDocumentResponse[]
}
