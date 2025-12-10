import type { StatusScreening } from "./status_screening"; 
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
    application?: any;       
    student_profile?: any;   
    admin_profile?: any;     
}

export interface CreateScreeningPayload {
    admin_profile_id: number;
    student_profile_id: number;
    application_id: number;
    status_screening_id: number;
    rejection_reason?: string | null;
}

export interface UpdateScreeningStatusPayload {
    status_screening_id: number;
    rejection_reason?: string | null; 
}