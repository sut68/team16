import type { StatusScreening } from "./status_screening";
import type { ScholarshipResponse } from "./scholarship";
import type { AdminProfileResponse, FeatureScholarshipResponse, StudentProfileResponse } from ".";
import type { SemasterResponse } from "./semaster";
import type { ApplicationScholarshipResponse } from "./application_scholarship";

export interface Screening {
    id: number;
    admin_profile_id: number;
    student_profile_id: number;
    application_scholarship_id: number;
    status_screening_id: number;
    rejection_reason?: string | null; 
    created_at: string;
    updated_at: string;
}

export interface ScreeningResponse extends Screening {
    status_screening?: StatusScreening;
    scholarship?: ScholarshipResponse;
    application_scholarship?: ApplicationScholarshipResponse;

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
    application_scholarship_id: number;
    status_screening_id: number;
    rejection_reason?: string | null; 
}