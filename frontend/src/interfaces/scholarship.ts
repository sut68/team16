
export interface ScholarshipStatusResponse {
  ID: number
  status_name: string
}

export interface ScholarshipTypeResponse {
  ID: number
  type_name: string
}

export interface RequirementResponse {
  ID: number;
  name: string;
}

export interface ApprovalRequirementResponse {
  ID: number
  scholarship_id: number
  requirement_id: number
  requirement: RequirementResponse
}

export interface ScholarshipResponse {
  ID: number
  scholarship_name: string
  description: string
  open_date: string
  close_date: string
  statusscholarship: ScholarshipStatusResponse
  typescholarship: ScholarshipTypeResponse
}