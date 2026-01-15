// src/validators/sponsor_validator.ts
import type { SponsorPayload, ContactPayload, ContactResponse } from '@/interfaces/sponsor'

// TYPES & INTERFACES

export interface ValidationResult {
  valid: boolean
  errors: Record<string, any>
}

export interface SponsorFormErrors {
  company_name?: string
  website?: string
  status?: string
  contacts?: Record<number, ContactErrors>
}

export interface ContactErrors {
  name?: string
  email?: string
  phone?: string
}

// VALIDATION RULES
export const SPONSOR_VALIDATION_RULES = {
  company_name: {
    minLength: 2,
    maxLength: 255,
    required: true,
  },
  website: {
    required: false,
    // รองรับทั้ง domain ธรรมดา (example.com) และ full URL (https://example.com)
    pattern: /^(https?:\/\/)?[a-zA-Z0-9][a-zA-Z0-9-]*(\.[a-zA-Z0-9-]+)+/i,
  },
  status: {
    required: true,
  },
} as const

export const CONTACT_VALIDATION_RULES = {
  name: {
    minLength: 2,
    maxLength: 100,
    required: true,
  },
  email: {
    required: true,
    pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  },
  phone: {
    maxLength: 20,
    required: false,
  },
} as const

// HELPER FUNCTIONS
function isValidURL(url = ''): boolean {
  const trimmed = String(url).trim()
  if (!trimmed) return true // Empty is valid (optional field)
  // รองรับทั้ง "example.com" และ "https://example.com"
  return SPONSOR_VALIDATION_RULES.website.pattern.test(trimmed)
}

function isValidEmail(email = ''): boolean {
  return CONTACT_VALIDATION_RULES.email.pattern.test(String(email).trim())
}

// VALIDATION FUNCTIONS - SPONSOR

export function validateSponsorForm(form: Partial<SponsorPayload>): ValidationResult {
  const errors: Record<string, any> = {}

  // Company Name validation
  if (!form.company_name || String(form.company_name).trim().length < SPONSOR_VALIDATION_RULES.company_name.minLength) {
    errors.company_name = `โปรดระบุชื่อบริษัท (อย่างน้อย ${SPONSOR_VALIDATION_RULES.company_name.minLength} ตัวอักษร)`
  }

  // Website validation - รองรับทั้ง domain และ full URL
  if (form.website && !isValidURL(String(form.website))) {
    errors.website = 'URL ไม่ถูกต้อง (เช่น example.com หรือ https://example.com)'
  }

  // Status validation
  if (!form.status) {
    errors.status = 'โปรดระบุสถานะ'
  }

  // Contacts validation
  const contacts = (form.contacts ?? []) as ContactPayload[]
  const contactErrors: Record<number, any> = {}
  contacts.forEach((c, idx) => {
    const ce: Record<string, string> = {}
    if (!c?.name || String(c.name).trim().length < CONTACT_VALIDATION_RULES.name.minLength) {
      ce.name = `โปรดระบุชื่อ (อย่างน้อย ${CONTACT_VALIDATION_RULES.name.minLength} ตัวอักษร)`
    }
    if (c?.email && !isValidEmail(String(c.email))) {
      ce.email = 'อีเมลไม่ถูกต้อง'
    }
    if (Object.keys(ce).length) contactErrors[idx] = ce
  })

  if (Object.keys(contactErrors).length) errors.contacts = contactErrors

  return { valid: Object.keys(errors).length === 0, errors }
}

// Individual field validators
export function validateCompanyName(name: string | undefined): string | undefined {
  const trimmed = String(name || '').trim()

  if (!trimmed) {
    return 'กรุณากรอกชื่อบริษัท'
  }
  if (trimmed.length < SPONSOR_VALIDATION_RULES.company_name.minLength) {
    return `ชื่อบริษัทต้องมีอย่างน้อย ${SPONSOR_VALIDATION_RULES.company_name.minLength} ตัวอักษร`
  }

  return undefined
}

export function validateWebsite(url: string | undefined): string | undefined {
  if (!url || !url.trim()) {
    return undefined // Optional field
  }
  if (!isValidURL(url)) {
    return 'URL ไม่ถูกต้อง (เช่น example.com หรือ https://example.com)'
  }

  return undefined
}

export function validateStatus(status: string | undefined): string | undefined {
  if (!status) {
    return 'กรุณาเลือกสถานะ'
  }
  return undefined
}

// VALIDATION FUNCTIONS - CONTACTS

export function validateContacts(contacts: ContactPayload[]) {
  const errors: Record<number, Record<string, string>> = {}
  contacts.forEach((c, idx) => {
    const e: Record<string, string> = {}
    if (!c.name || String(c.name).trim().length < CONTACT_VALIDATION_RULES.name.minLength) {
      e.name = `โปรดระบุชื่อ (อย่างน้อย ${CONTACT_VALIDATION_RULES.name.minLength} ตัวอักษร)`
    }
    if (!c.email || !isValidEmail(String(c.email))) {
      e.email = 'ใส่อีเมลที่ถูกต้อง'
    }
    if (Object.keys(e).length) errors[idx] = e
  })
  return { valid: Object.keys(errors).length === 0, errors }
}

export function validateContactName(name: string | undefined): string | undefined {
  const trimmed = String(name || '').trim()

  if (!trimmed) {
    return 'กรุณากรอกชื่อผู้ติดต่อ'
  }
  if (trimmed.length < CONTACT_VALIDATION_RULES.name.minLength) {
    return `ชื่อต้องมีอย่างน้อย ${CONTACT_VALIDATION_RULES.name.minLength} ตัวอักษร`
  }

  return undefined
}

export function validateContactEmail(email: string | undefined): string | undefined {
  const trimmed = String(email || '').trim()

  if (!trimmed) {
    return 'กรุณากรอกอีเมล'
  }
  if (!isValidEmail(trimmed)) {
    return 'รูปแบบอีเมลไม่ถูกต้อง'
  }

  return undefined
}

export function validateContactPhone(phone: string | undefined): string | undefined {
  if (!phone) {
    return undefined // Optional field
  }
  if (phone.length > CONTACT_VALIDATION_RULES.phone.maxLength) {
    return `เบอร์โทรต้องไม่เกิน ${CONTACT_VALIDATION_RULES.phone.maxLength} ตัวอักษร`
  }

  return undefined
}

// UTILITY FUNCTIONS

export function buildContactsBatch(original: ContactResponse[] = [], current: ContactPayload[] = []) {
  const origIds = new Set<number>(original.filter(c => c?.ID != null).map(c => c.ID))
  const currIds = new Set<number>(current.filter(c => (c as any).ID != null).map(c => (c as any).ID))
  const delete_ids = Array.from(origIds).filter(id => !currIds.has(id))

  const upsert = (current || []).map(c => {
    const out: any = {
      name: String(c.name || '').trim(),
      email: String(c.email || '').trim(),
      phone: String(c.phone || '').trim(),
      position: c.position == null ? null : String(c.position).trim(),
    }
    if ((c as any).ID != null) out.ID = (c as any).ID
    return out
  }).filter(c => !!(c.name || c.email))

  return { upsert: upsert.length ? upsert : undefined, delete_ids: delete_ids.length ? delete_ids : undefined }
}

// GROUPED EXPORTS
export const SponsorValidators = {
  sponsor: {
    rules: SPONSOR_VALIDATION_RULES,
    validateForm: validateSponsorForm,
    validateCompanyName,
    validateWebsite,
    validateStatus,
  },
  contact: {
    rules: CONTACT_VALIDATION_RULES,
    validateContacts,
    validateName: validateContactName,
    validateEmail: validateContactEmail,
    validatePhone: validateContactPhone,
    buildBatch: buildContactsBatch,
  },
} as const