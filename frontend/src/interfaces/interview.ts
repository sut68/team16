import type { ScholarshipResponse } from "./scholarship";
import type { User } from "./user";


export interface InterviewRound {
    ID: number;
    name: string;
    description: string;
    start_date_time: string;
    end_date_time: string;
    slot_duration: number;
    scholarship_id: number;
    scholarship: ScholarshipResponse;
    admin_profile_id: number;
    admin_profile: User;
    slots: Slot[];
}

export interface Interviewer {
    ID: number;
    interviewer_firstname: string;
    interviewer_lastname: string;
    email: string;
}

export interface Slot {
    ID: number;
    start_time: string;
    end_time: string;
    is_booked: boolean;
    interview_round_id: number;
    interviewer_slots: InterviewerSlot[];
    interview_booking: InterviewBooking | null;
}

export interface InterviewerSlot {
    ID: number;
    interviewer_id: number;
    interviewer: Interviewer;
    slot_id: number;
}

export interface InterviewBooking {
    ID: number;
    booking_date: string;
    booking_time: string;
    status: string;
    slot_id: number;
    application_scholarship_id: number;
}
