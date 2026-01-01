import type {
  ContactPayload,
  SponsorPayload,
  SponsorResponse,
  ContactResponse,
  SponsorScholarshipResponse
} from "../../interfaces/sponsor";
import { Get, Patch, Post, Delete } from "../api";

export const SponsorService = {
  async getAll(): Promise<SponsorResponse[]> {
    const data: SponsorResponse[] = await Get("/sponsors");
    return data;
  },

  async getById(id: number): Promise<SponsorResponse> {
    const data: SponsorResponse = await Get(`/sponsors/${id}`);
    return data;
  },

  async create(payload: SponsorPayload): Promise<SponsorResponse> {
    const body = stripUndefined(payload)
    const data: SponsorResponse = await Post("/sponsors", body);
    return data;
  },

  async update(id: number, payload: SponsorUpdatePayload): Promise<SponsorResponse> {
    const body = stripUndefined(payload);
    const data: SponsorResponse = await Patch(`/sponsors/${id}`, body);
    return data;
  },

  async updateContacts(
    id: number,
    payload: BatchContactsPayload
  ): Promise<{ contacts: ContactResponse[] }> {
    const body = stripUndefined(payload);
    const data: { contacts: ContactResponse[] } = await Patch(`/sponsors/${id}/contacts`, body);
    return data;
  },

  async delete(id: number): Promise<void> {
    await Delete(`/sponsors/${id}`);
  },

  async getScholarships(id: number): Promise<SponsorScholarshipResponse[]> {
    const data: SponsorScholarshipResponse[] = await Get(`/sponsors/${id}/scholarships`);
    return data ?? [];
  },
};

export default SponsorService

export type { SponsorScholarshipResponse } from "../../interfaces/sponsor";

// helper
function stripUndefined<T extends Record<string, any>>(obj: T): Partial<T> {
  const out: Partial<T> = {}
  for (const key in obj) {
    if (obj[key] !== undefined) out[key] = obj[key]
  }
  return out
}

export interface SponsorUpdatePayload
  extends Partial<Omit<SponsorPayload, "ID" | "contacts">> {}

export interface BatchContactsPayload {
  upsert?: ContactPayload[]
  delete_ids?: number[]
}