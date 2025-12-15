import type { SemasterResponse } from "./semaster";
import type { SponsorResponse } from "./sponsor"

export interface ScholarshipStatusResponse {
  ID: number
  status_name: string
}

export interface ScholarshipStatusCreate {
  status_name: string
}

export interface ScholarshipStatusUpdate {
  ID: number
  status_name?: string
}

export interface ScholarshipTypeResponse {
  ID: number
  type_name: string
}

export interface ScholarshipTypeCreate {
  type_name: string
}

export interface ScholarshipTypeUpdate {
  ID: number
  type_name?: string
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
  
  statusscholarship_id: number;
  statusscholarship: ScholarshipStatusResponse

  typescholarship_id: number;
  typescholarship: ScholarshipTypeResponse

  sponsor_id: number
  sponsor: SponsorResponse

  semaster_id: number
  semaster: SemasterResponse
  

  approval_requirements: ApprovalRequirementResponse[]
}

export interface ScholarshipCreate {
  scholarship_name: string;
  description: string;
  open_date: string;
  close_date: string;

  statusscholarship_id: number;
  typescholarship_id: number;

  sponsor_id: number

  semaster_id: number
  
}

export interface ScholarshipUpdate {
  ID: number;
  scholarship_name?: string;
  description?: string;
  open_date?: string;
  close_date?: string;

  statusscholarship_id?: number;
  typescholarship_id?: number;

  sponsor_id?: number

  semaster_id?: number
}