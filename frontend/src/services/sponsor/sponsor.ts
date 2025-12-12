import { httpsTest } from "../api/httpsTest";
import type { ContactPayload, SponsorPayload, SponsorResponse } from "../../interfaces/sponsor";

export const SponsorService = {
  async getAll(): Promise<SponsorResponse[]> {
    const res = await httpsTest.get<SponsorResponse[]>("/sponsors");
    return res.data;
  },

  async getById(id: number): Promise<SponsorResponse> {
    const res = await httpsTest.get<SponsorResponse>(`/sponsors/${id}`);
    return res.data;
  },

  async create(payload: SponsorPayload): Promise<SponsorResponse> {
    const body = stripUndefined(payload);
    const res = await httpsTest.post<SponsorResponse>("/sponsors", body);
    return res.data;
  },

  async update(id: number, payload: SponsorUpdatePayload): Promise<SponsorResponse> {
    const body = stripUndefined(payload);
    const res = await httpsTest.patch<SponsorResponse>(`/sponsors/${id}`, body);
    return res.data;
  },

  async updateContacts(
    id: number,
    payload: BatchContactsPayload
  ): Promise<{ contacts: ContactPayload[] }> {
    const body = stripUndefined(payload);
    const res = await httpsTest.patch<{ contacts: ContactPayload[] }>(`/sponsors/${id}/contacts`, body);
    return res.data;
  },

  async delete(id: number): Promise<void> {
    await httpsTest.delete(`/sponsors/${id}`);
  },
};

export default SponsorService;

/* helper */
function stripUndefined<T extends Record<string, any>>(obj: T): Partial<T> {
  const out: Partial<T> = {};
  for (const key in obj) {
    if (obj[key] !== undefined) out[key] = obj[key];
  }
  return out;
}

/* types */
export interface SponsorUpdatePayload extends Partial<Omit<SponsorPayload, "ID" | "contacts">> {}

export interface BatchContactsPayload {
  upsert?: ContactPayload[];
  delete_ids?: number[];
}
