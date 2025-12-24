// src/validators/sponsorValidator.ts
import type { SponsorPayload, ContactPayload, ContactResponse } from '@/interfaces/sponsor';

export interface ValidationResult {
  valid: boolean;
  errors: Record<string, any>;
}

function isValidURL(url = '') {
  return /^https?:\/\//i.test(String(url).trim());
}
function isValidEmail(email = '') {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(email).trim());
}

export function validateSponsorForm(form: Partial<SponsorPayload>): ValidationResult {
  const errors: Record<string, any> = {};

  if (!form.company_name || String(form.company_name).trim().length < 2) {
    errors.company_name = 'โปรดระบุชื่อบริษัท (อย่างน้อย 2 ตัวอักษร)';
  }

  if (form.website && !isValidURL(String(form.website))) {
    errors.website = 'ใส่ URL ให้มี http:// หรือ https://';
  }

  if (!form.status) {
    errors.status = 'โปรดระบุสถานะ';
  }

  // contacts validation
  const contacts = (form.contacts ?? []) as ContactPayload[];
  const contactErrors: Record<number, any> = {};
  contacts.forEach((c, idx) => {
    const ce: Record<string, string> = {};
    if (!c?.name || String(c.name).trim().length < 2) {
      ce.name = 'โปรดระบุชื่อ (อย่างน้อย 2 ตัวอักษร)';
    }
    if (c?.email && !isValidEmail(String(c.email))) {
      ce.email = 'อีเมลไม่ถูกต้อง';
    }
    if (Object.keys(ce).length) contactErrors[idx] = ce;
  });

  if (Object.keys(contactErrors).length) errors.contacts = contactErrors;

  return { valid: Object.keys(errors).length === 0, errors };
}

export function validateContacts(contacts: ContactPayload[]) {
  const errors: Record<number, Record<string, string>> = {}
  contacts.forEach((c, idx) => {
    const e: Record<string,string> = {}
    if (!c.name || String(c.name).trim().length < 2) e.name = 'โปรดระบุชื่อ (อย่างน้อย 2 ตัวอักษร)'
    if (!c.email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(c.email))) e.email = 'ใส่อีเมลที่ถูกต้อง'
    if (Object.keys(e).length) errors[idx] = e
  })
  return { valid: Object.keys(errors).length === 0, errors }
}

export function buildContactsBatch(original: ContactResponse[] = [], current: ContactPayload[] = []) {
  const origIds = new Set<number>(original.filter(c => c?.ID != null).map(c => c.ID));
  const currIds = new Set<number>(current.filter(c => (c as any).ID != null).map(c => (c as any).ID));
  const delete_ids = Array.from(origIds).filter(id => !currIds.has(id));

  const upsert = (current || []).map(c => {
    const out: any = {
      name: String(c.name || '').trim(),
      email: String(c.email || '').trim(),
      phone: String(c.phone || '').trim(),
      position: c.position == null ? null : String(c.position).trim(),
    };
    if ((c as any).ID != null) out.ID = (c as any).ID;
    return out;
  }).filter(c => !!(c.name || c.email)); // optional: remove empty rows

  return { upsert: upsert.length ? upsert : undefined, delete_ids: delete_ids.length ? delete_ids : undefined };
}