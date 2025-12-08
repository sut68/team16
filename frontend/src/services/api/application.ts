import { Get } from './https';
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

  // We cannot use the JSON Post wrapper because this is multipart/form-data.
  // We need a custom fetch or a different axios setup.
  // For now, let's use a raw fetch.
  // NOTE: This will not have the auth headers from the wrapper.
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";
  const response = await fetch(`${API_URL}/application-documents`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || 'Upload failed');
  }

  return response.json();
};
