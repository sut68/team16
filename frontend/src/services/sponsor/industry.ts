import axios from "axios";
import type { IndustryResponse } from "../../interfaces/sponsor";

const api = axios.create({
  baseURL: "http://localhost:8080/api",
  timeout: 10000,
})

export const IndustryService = {
  async getAll(): Promise<IndustryResponse[]> {
    const res = await api.get<IndustryResponse[]>('/industries');
    return res.data;
  },
};