import { Get, Post, Put, Delete } from './https';
import type { ScholarshipCreate, ScholarshipUpdate, ScholarshipResponse } from '@/interfaces';

/**
 * Fetches all scholarships from the backend.
 * @returns A promise that resolves to an array of scholarships.
 */
export const getScholarships = async (): Promise<ScholarshipResponse[]> => {
  // The Get function in https.ts already extracts the .data property on success
  const scholarships = await Get('/scholarship');
  return scholarships as ScholarshipResponse[];
};

/**
 * Submits an application for a specific scholarship with full data.
 * @param scholarshipId The ID of the scholarship to apply for.
 * @param payload The application data including student profile info and family info.
 * @returns A promise that resolves to the API response.
 */
export interface ApplyScholarshipPayload {
  student_profile_id: number;
  application_reason?: string;
  email?: string;
  phone?: string;
  father_occupation?: string;
  father_income?: number;
  mother_occupation?: string;
  mother_income?: number;
  guardian_occupation?: string;
  guardian_income?: number;
}

export const applyForScholarship = async (scholarshipId: number, payload: ApplyScholarshipPayload) => {
  const response = await Post(`/scholarship/${scholarshipId}/apply`, payload);
  return response;
};

export const ScholarshipAPI = {
  create: (data: ScholarshipCreate) => Post("/scholarship", data),
  getAll: () => Get("/scholarship"),
  getById: (id: number) => Get(`/scholarship/${id}`),
  update: (id: number, data: ScholarshipUpdate) => Put(`/scholarship/${id}`, data),
  delete: (id: number) => Delete(`/scholarship/${id}`),
};