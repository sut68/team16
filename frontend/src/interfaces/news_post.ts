import type { AdminProfileResponse } from "./user";
import type { ScholarshipResponse } from "./scholarship";
// สมมติว่า FeaturescholarshipResponse ถูก import แล้ว
import type { FeatureScholarshipResponse } from "./featurescholarship"; 

// 1. Interface สำหรับข้อมูลที่ได้จาก Backend (NewsPost object)
export interface NewsPost {
    ID: number; 
    title: string;
    post_detail: string;
    file_path: string; // คาดหวัง string จาก Backend
    
    admin_id: number;
    Admin?: AdminProfileResponse; // Object ที่ Preload มา

    scholarship_id: number;
    scholarship?: ScholarshipResponse; // Object ที่ Preload มา

    status_news_id: number;
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
}


// 3. Interface สำหรับ Payload สำหรับการอัปเดต (ส่งไป Backend)
export interface UpdateNewsPost {
    title?: string;
    post_detail?: string;
    // อนุญาตให้เป็น File object สำหรับอัปโหลดรูป
    file_path?: string | File ; 
    status_news_id?: number;

    admin_id?: number;
    scholarship_id?: number;
}


// 4. Interface สำหรับ Payload สำหรับการสร้าง
export interface CreateNewsPostPayload {
    title: string;
    post_detail: string;

    file_path?: string | File; 
    admin_id: number;
    scholarship_id: number;
    status_news_id: number;
}

// 5. Interface สำหรับ Status
export interface StatusNewsPost {
    ID: number;
    status_news: string;
}