import type { StudentProfileResponse } from './user'
import type { ApprovalTaskResponse } from './approval'
import type { ApplicationScholarshipResponse } from './application_scholarship'
import type { SemasterResponse } from './semaster'

export interface ApplicationDocumentResponse {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  file_name: string
  file_path: string
  file_type: string
  uploaded_by: string
  application_scholarship_id: number
  application_scholarship: ApplicationScholarshipResponse
  approval_tasks: ApprovalTaskResponse[]
}

export interface ApplicationResponse {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  student_profile_id: number
  semaster_id: number
  student_profile: StudentProfileResponse
  semaster: SemasterResponse
  application_scholarships: any[]
}