
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
  sponsor_id?: number
}

export interface SponsorResponse {
  ID: number
  company_name: string
  industry_id?: number | null
  industry?: IndustryResponse | null
  website?: string | null
  status: string
  description?: string | null
  contacts?: ContactResponse[]
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
