import type { ScholarshipResponse } from "./scholarship";
//import type { UserResponse } from "./user";
import type { ApplicationScholarshipResponse } from "./application_scholarship";
import type { AdminProfileResponse } from "./user";


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
    //admin_profile: UserResponse;

    admin_profile: AdminProfileResponse;
    interview_mode_id?: number;
    interview_mode?: InterviewMode;
    location_id?: number | null;
    location?: Location | null;
    meeting_link?: string;
    slots: Slot[];
}

export interface InterviewMode {
    ID: number;
    name: string;
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
    status: string;
    interview_round_id: number;
    interviewer_slots: InterviewerSlot[];
    interviewe_bookings: InterviewBooking[];
}

export interface InterviewerSlot {
    ID: number;
    interviewer_id: number;
    interviewer: Interviewer;
    slot_id: number;
}

export interface InterviewBooking {
    ID: number;
    status: string;
    slot_id: number;
    slot?: Slot;
    application_scholarship_id: number;
    application_scholarship?: ApplicationScholarshipResponse;
}

export interface Location {
    ID: number;
    name: string;
    building: string;
    room: string;
    floor: number;
}
