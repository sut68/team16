import axios from "axios";
import type { SponsorPayload, SponsorResponse } from "../../interfaces/sponsor";

const api = axios.create({
  baseURL: "http://localhost:8080/api",
  timeout: 10000,
})

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

  async update(id: number, data: SponsorPayload): Promise<SponsorResponse> {
    const res = await api.patch<SponsorResponse>(`/sponsors/${id}`, data);
    return res.data;
  },

  async delete(id: number): Promise<void> {
    const res = await api.delete<void>(`/sponsors/${id}`);
    return res.data;
  },
};