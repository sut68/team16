import type {
  ApplicationResponse,
  ScholarshipResponse,
  ApplicationDocumentResponse,
} from './';
import type { ScreeningResponse } from './screening';
import type { InterviewBooking } from './interview';
import type { EvaluationResponse } from './evaluation';

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

  // Interview bookings
  interviewe_bookings?: InterviewBooking[];

  // Evaluation results
  evaluations?: EvaluationResponse[];
  final_total_score?: number;
  final_decision?: 'pending' | 'approved' | 'rejected' | 'waitlist';
}
