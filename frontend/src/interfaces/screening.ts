import type { StatusScreening } from "./status_screening";
import type { ScholarshipResponse } from "./scholarship";
import type { AdminProfileResponse, ApplicationResponse, FeatureScholarshipResponse, StudentProfileResponse } from ".";
import type { SemasterResponse } from "./semaster";

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
    application_scholarship?: ApplicationResponse;

    student_profile?: StudentProfileResponse;   
    admin_profile?: AdminProfileResponse; 

    feature_scholarships?: FeatureScholarshipResponse[];
    semaster?: SemasterResponse[];
}


export interface UpdateScreeningStatusPayload {
    status_screening_id: number;
    rejection_reason?: string | null; 
}

export interface CreateScreeningPayload {
    admin_profile_id: number;
    student_profile_id: number;
    application_id: number;
    status_screening_id: number;
    rejection_reason?: string | null; 
}