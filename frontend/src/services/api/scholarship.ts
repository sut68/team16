import { Get, Post } from './https';
import type { ScholarshipResponse } from '@/interfaces';

/**
 * Fetches all scholarships from the backend.
 * @returns A promise that resolves to an array of scholarships.
 */
export const getScholarships = async (): Promise<ScholarshipResponse[]> => {
  // The Get function in https.ts already extracts the .data property on success
  const scholarships = await Get('/scholarships');
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