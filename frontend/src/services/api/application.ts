import { Get, Delete, Post } from './https';
import type { ApplicationScholarshipResponse } from '@/interfaces';

/**
 * Fetches all scholarship applications for a given student.
 * @param studentProfileId The ID of the student.
 * @returns A promise that resolves to an array of scholarship applications.
 */
export const getStudentApplications = async (studentProfileId: number): Promise<ApplicationScholarshipResponse[]> => {
  const response = await Get(`/students/${studentProfileId}/applications`);
  return response as ApplicationScholarshipResponse[];
};

/**
 * Fetches all application scholarships, optionally filtered by status.
 * @param status Optional status filter (e.g., 'qualified', 'pending')
 * @returns A promise that resolves to an array of application scholarships.
 */
export const getAllApplicationScholarships = async (status?: string): Promise<ApplicationScholarshipResponse[]> => {
  let url = '/application-scholarships';
  if (status) {
    url += `?status=${status}`;
  }
  const response = await Get(url);
  return response as ApplicationScholarshipResponse[];
};

/**
 * Uploads a document for a specific scholarship application.
 * @param applicationScholarshipId The ID of the scholarship application link.
 * @param file The file to upload.
 * @param studentProfileId The ID of the uploading student.
 * @returns A promise that resolves to the API response.
 */
export const uploadDocument = async (applicationScholarshipId: number, file: File, studentProfileId: number) => {
  const formData = new FormData();
  formData.append('document', file);
  formData.append('application_scholarship_id', String(applicationScholarshipId));
  formData.append('uploaded_by', String(studentProfileId));

  // Use Post wrapper which handles auth headers and FormData correctly
  const response = await Post('/application-documents', formData);

  if (response?.error) {
    throw new Error(response.error || 'Upload failed');
  }

  return response;
};

/**
 * ยกเลิกการสมัครทุน
 * @param applicationScholarshipId ID ของการสมัครทุน
 * @returns Promise ที่ resolve เป็น response จาก API
 */
export const cancelApplicationScholarship = async (applicationScholarshipId: number): Promise<{ message: string }> => {
  const response = await Delete(`/application-scholarships/${applicationScholarshipId}/cancel`);
  return response as { message: string };
};
