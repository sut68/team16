// Response from backend
// ----------------------
export interface ScholarshipStatusResponse {
  ID: number
  // Add other fields from Statusscholarship if necessary
}

export interface ScholarshipTypeResponse {
  ID: number
  // Add other fields from Typescholarship if necessary
}

export interface ApprovalRequirementResponse {
  ID: number
  name: string
  description: string
  scholarship_id: number
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
