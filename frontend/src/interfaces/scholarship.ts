
export interface ScholarshipStatusResponse {
  ID: number
  status_name: string
}

export interface ScholarshipTypeResponse {
  ID: number
  type_name: string
}

export interface ScholarshipResponse {
  ID: number
  scholarship_name: string
  description: string
  open_date: string
  close_date: string
  statusscholarship_id: number
  statusscholarship: ScholarshipStatusResponse
  typescholarship_id: number
  typescholarship: ScholarshipTypeResponse
  approval_requirements: ApprovalRequirementResponse[]
}

export interface ApprovalRequirementResponse {
  ID: number
  name: string
  description: string
  scholarship_id: number
  scholarship: ScholarshipResponse
}