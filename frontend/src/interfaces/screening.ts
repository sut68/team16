import type { StatusScreening } from "./status_screening";
import type { ApplicationResponse } from "./application";
import type { AdminProfileResponse } from "./user";
import type { ScholarshipResponse } from "./scholarship";

export interface Screening {
    ID: number;           
    CreatedAt: string;    
    UpdatedAt: string;    
    DeletedAt?: string | null;

    admin_profile_id: number;
    application_id: number;
    status_screening_id: number;

    scholarship_id: number;
    
    rejection_reason?: string | null; 
}

export interface ScreeningResponse extends Screening {
    status_screening?: StatusScreening;
    application?: ApplicationResponse;      
    admin_profile?: AdminProfileResponse;
    scholarship?: ScholarshipResponse;
}

export interface CreateScreeningPayload {
    admin_profile_id: number;
    scholarship_id: number;
    application_id: number;
    status_screening_id: number;
    rejection_reason?: string | null;
}

export interface UpdateScreeningStatusPayload {
    status_screening_id: number;
    rejection_reason?: string | null; 
}