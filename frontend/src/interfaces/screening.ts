import type { StatusScreening } from "./status_screening";
//import type { ApplicationResponse } from "./application";
//import type { AdminProfileResponse } from "./user";
import type { ScholarshipResponse } from "./scholarship";
export interface Screening {
    id: number;
    admin_profile_id: number;
    student_profile_id: number;
    application_id: number;
    status_screening_id: number;
    rejection_reason?: string | null; 
    created_at: string;
    updated_at: string;
}

export interface ScreeningResponse extends Screening {
    status_screening?: StatusScreening;
    scholarship?: ScholarshipResponse;

    application?: any;       
    student_profile?: any;   
    admin_profile?: any; 
}

export interface CreateScreeningPayload {
    admin_profile_id: number;
    scholarship_id: number;
    student_profile_id: number;
    application_id: number;
    status_screening_id: number;
    rejection_reason?: string | null;
}

export interface UpdateScreeningStatusPayload {
    status_screening_id: number;
    rejection_reason?: string | null; 
}