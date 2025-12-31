import { Get } from "./https";

export interface DashboardStats {
    total_applications: number;
    active_scholarships: number;
    pending_reviews: number;
    total_students: number;
    applications_status: StatusDistribution[];
    popular_scholarships: PopularScholarship[];
    recent_applications: RecentApplication[];
}

export interface StatusDistribution {
    status: string;
    count: number;
}

export interface PopularScholarship {
    id: number;
    name: string;
    applicant_count: number;
    status: string;
}

export interface RecentApplication {
    id: number;
    student_name: string;
    student_code: string;
    scholarship_name: string;
    applied_at: string;
    status: string;
}

export const getDashboardStats = async (): Promise<DashboardStats | null> => {
    const result = await Get("/dashboard", true);
    // Check if result is error or data
    // The Get function handles some errors but returns object.
    // Standard check:
    if (result && result.total_applications !== undefined) {
        return result as DashboardStats;
    }
    return null;
};
