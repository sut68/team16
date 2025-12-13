import type { AdminProfileResponse } from "./user";
import type { ScholarshipResponse } from "./scholarship";

// 1. Interface สำหรับข้อมูลที่ได้จาก Backend (Response)
export interface NewsPost {
    ID: number;           // Go struct มี ID
    title: string;
    post_detail: string;  // เพิ่ม field นี้ตาม Go Backend
    file_path: string | File ;    // รับมาเป็น URL หรือ Path string

    admin_id: number;     // ค่า FK
    Admin?: AdminProfileResponse; // Object ที่ Preload มา

    scholarship_id: number;
    Scholarship?: ScholarshipResponse;

    status_news_id: number;
    StatusNews?: StatusNewsPost;

    CreatedAt: string;    
    UpdatedAt: string;    
    DeletedAt?: string | null;
}

export interface UpdateNewsPost {
    title?: string;
    post_detail?: string;
    file_path?: string | File ;
    status_news_id?: number;

    admin_id?: number;
    scholarship_id?: number;
}


export interface CreateNewsPostPayload {
    title: string;
    post_detail: string;

    file_path?: string; 
    admin_id: number;
    scholarship_id: number;
    status_news_id: number;
}

export interface StatusNewsPost {
    ID: number;
    status_news: string;
}