import type { IndustryResponse } from "../../interfaces/sponsor";
import { Get } from "../api";

export const IndustryService = {
  async getAll(requireAuth: boolean = true): Promise<IndustryResponse[]> {
    const data: IndustryResponse[] = await Get("/industries", requireAuth);
    return data;
  }
};