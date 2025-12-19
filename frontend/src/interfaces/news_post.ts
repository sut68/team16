import type { AdminProfileResponse } from "./user";
import type { ScholarshipResponse } from "./scholarship";
import type { FeatureScholarshipResponse } from "./featurescholarship"; 

// 1. Interface สำหรับข้อมูลที่ได้จาก Backend (NewsPost object)
export interface NewsPost {
    ID: number; 
    title: string;
    post_detail: string;
    file_path: string | null;

    admin_id: number;
    admin_profile?: AdminProfileResponse; // Object ที่ Preload มา

    scholarship_id: number | null;
    scholarship?: ScholarshipResponse; // Object ที่ Preload มา

    status_news_id: number | null;
    StatusNews?: StatusNewsPost; // Object ที่ Preload มา

    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt?: string | null;
}

// 2. Interface สำหรับ Response Body ที่ Go Backend ส่งจาก GetNewsPostByID
// 🔥 NEW INTERFACE 🔥
export interface NewsPostDetailResponse {
    news_post: NewsPost;
    features_scholarship: FeatureScholarshipResponse[];
    features?: any;
}


// 3. Interface สำหรับ Payload สำหรับการอัปเดต (ส่งไป Backend)
export interface UpdateNewsPost {
    title?: string;
    post_detail?: string;
    file_path?: null | File ; 
    status_news_id?: number;
    admin_id?: number;
    admin_profile?: AdminProfileResponse;
    scholarship_id?: number | null;
}


// 4. Interface สำหรับ Payload สำหรับการสร้าง
export interface CreateNewsPostPayload {
    title: string;
    post_detail: string;

    file_path?: null | File; 
    admin_id: number;
    admin_profile?: AdminProfileResponse;
    scholarship_id: number | null;
    status_news_id: number;
}

// 5. Interface สำหรับ Status
export interface StatusNewsPost {
    ID: number;
    status_news: string;
}

export interface StudentFavoriteNewsPost {
    ID: number;
    student_profile_id: number;
    news_post_id: number;
    
    news_post?: NewsPost;
    CreatedAt: string;
}
