import { Get, Post, Put, Delete } from './https';
import type { InterviewRound, Interviewer, InterviewBooking, InterviewMode } from '@/interfaces/interview';

// We can define more specific types for create/update if they differ from the response types
export interface InterviewRoundCreate {
  name: string;
  description: string;
  start_date_time: string;
  end_date_time: string;
  slot_duration: number;
  scholarship_id: number;
  admin_profile_id: number;
  interview_mode_id: number;
  location_id?: number | null;
  meeting_link?: string;
  interviewer_ids: number[];
  slots: { status: string }[];
}
export type InterviewRoundUpdate = Partial<Omit<InterviewRoundCreate, 'slots'>> & {
  slots?: { id: number, status: string }[];
};

export type InterviewerCreate = Omit<Interviewer, 'ID'>;

export type InterviewBookingCreate = Omit<InterviewBooking, 'ID'>;


export const InterviewAPI = {
  // Interview Rounds
  getAllRounds: (): Promise<InterviewRound[]> => Get("/interview-rounds"),
  getRoundById: (id: number): Promise<InterviewRound> => Get(`/interview-rounds/${id}`),
  createRound: (data: InterviewRoundCreate): Promise<InterviewRound> => Post("/interview-rounds", data),
  updateRound: (id: number, data: InterviewRoundUpdate): Promise<InterviewRound> => Put(`/interview-rounds/${id}`, data),
  deleteRound: (id: number): Promise<any> => Delete(`/interview-rounds/${id}`),

  // Interviewers
  getAllInterviewers: (): Promise<Interviewer[]> => Get("/interviewers"),
  createInterviewer: (data: InterviewerCreate): Promise<Interviewer> => Post("/interviewers", data),

  // Interview Bookings
  createBooking: (data: InterviewBookingCreate): Promise<InterviewBooking> => Post("/interview-bookings", data),
  getStudentBookings: (studentProfileId: number): Promise<InterviewBooking[]> => Get(`/students/${studentProfileId}/interview-bookings`),
  deleteBooking: (id: number): Promise<any> => Delete(`/interview-bookings/${id}`),

  // Interview Modes
  getAllModes: (): Promise<InterviewMode[]> => Get("/interview-modes"),
};
