import axios from "axios";
import type { SponsorPayload, SponsorResponse } from "../../interfaces/sponsor";

const api = axios.create({
  baseURL: "http://localhost:8080/api",
  timeout: 10000,
})

export type SponsorUpdatePayload = Partial<Omit<SponsorPayload, 'ID'>>;

function stripUndefined<T extends Record<string, any>>(obj: T): Partial<T> {
  return Object.entries(obj).reduce((acc, [k, v]) => {
    if (v !== undefined) (acc as any)[k] = v;
    return acc;
  }, {} as Partial<T>);
}

export const SponsorService = {
  async getAll(): Promise<SponsorResponse[]> {
    const res = await api.get<SponsorResponse[]>('/sponsors');
    return res.data;
  },

  async getById(id: number): Promise<SponsorResponse> {
    const res = await api.get<SponsorResponse>(`/sponsors/${id}`);
    return res.data;
  },

  async create(data: SponsorPayload): Promise<SponsorResponse> {
    const res = await api.post<SponsorResponse>('/sponsors', data);
    return res.data;
  },

  async update(id: number, data: SponsorUpdatePayload): Promise<SponsorResponse> {
    const payload = stripUndefined(data);
    const res = await api.patch<SponsorResponse>(`/sponsors/${id}`, payload);
    return res.data;
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/sponsors/${id}`);
  },
};