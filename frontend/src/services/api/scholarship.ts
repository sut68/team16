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
 * Submits an application for a specific scholarship.
 * @param scholarshipId The ID of the scholarship to apply for.
 * @param studentProfileId The ID of the student applying.
 * @returns A promise that resolves to the API response.
 */
export const applyForScholarship = async (scholarshipId: number, studentProfileId: number) => {
  const payload = {
    student_profile_id: studentProfileId,
  };
  const response = await Post(`/scholarships/${scholarshipId}/apply`, payload);
  return response;
};

export const ScholarshipAPI = {
  create: (data: ScholarshipCreate) => Post("/scholarship", data),
  getAll: () => Get("/scholarship"),
  getById: (id: number) => Get(`/scholarship/${id}`),
  update: (id: number, data: ScholarshipUpdate) => Put(`/scholarship/${id}`, data),
  delete: (id: number) => Delete(`/scholarship/${id}`),
};