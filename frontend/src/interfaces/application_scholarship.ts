import type {
  ApplicationResponse,
  ScholarshipResponse,
  ApplicationDocumentResponse,
} from './';
import type { ScreeningResponse } from './screening';

export interface ApplicationScholarshipResponse {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  status: string;

  application_id: number;
  application: ApplicationResponse;

  scholarship_id: number;
  scholarship: ScholarshipResponse;

  application_documents: ApplicationDocumentResponse[];
  screening?: ScreeningResponse;
}
