// Response from backend
// ----------------------
export interface IndustryResponse {
  ID: number
  name: string
}

export interface ContactResponse {
  ID: number
  name: string
  email: string
  phone: string
  position?: string | null
}

export interface SponsorResponse {
  ID: number
  company_name: string
  industry_id?: number | null
  website?: string | null
  status: string
  description?: string | null
  contacts?: ContactResponse[]
  industry?: IndustryResponse | null
}
// ----------------------


// Payload send to Backend
// ----------------------
export interface ContactPayload {
  ID?: number
  name: string
  email: string
  phone: string
  position?: string | null
}

export interface SponsorPayload {
  ID?: number
  company_name: string
  industry_id?: number | null
  website?: string | null
  status: string
  description?: string | null
  contacts?: ContactPayload[]
}
// ----------------------


// For UI
// ----------------------
export interface SponsorView {
  ID: number
  company_name: string
  website?: string | null
  industry_name?: string | null
  industry_id?: number | null
  status: string
  contacts_count: number
}

// Sponsor Scholarships Response
// ----------------------
export interface SponsorScholarshipStatusResponse {
  ID: number
  status_name: string
}

export interface SponsorScholarshipTypeResponse {
  ID: number
  type_name: string
}

export interface SponsorSemasterResponse {
  ID: number
  year: number
  term: number
  is_active?: boolean
}

export interface SponsorScholarshipResponse {
  ID: number
  scholarship_name: string
  description: string
  open_date: string
  close_date: string
  statusscholarship_id: number
  statusscholarship: SponsorScholarshipStatusResponse
  typescholarship_id: number
  typescholarship: SponsorScholarshipTypeResponse
  semaster_id: number
  semaster: SponsorSemasterResponse
}