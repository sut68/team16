import { Get, Post, Put, Delete } from './https';
import type { NewsPost, CreateNewsPostPayload, UpdateNewsPost } from '@/interfaces/news_post';

// Helper Function: ปรับปรุงให้ Safe สำหรับ Update
function convertToFormData(payload: Partial<UpdateNewsPost>): FormData {
    const formData = new FormData();

    if (payload.title !== undefined) formData.append("title", payload.title);
    if (payload.post_detail !== undefined) formData.append("post_detail", payload.post_detail);
    if (payload.admin_id !== undefined) formData.append("admin_id", payload.admin_id.toString());
    if (payload.scholarship_id !== undefined && payload.scholarship_id !== null) formData.append("scholarship_id", payload.scholarship_id.toString());
    if (payload.status_news_id !== undefined) formData.append("status_news_id", payload.status_news_id.toString());

    if (payload.file_path) {
    if (payload.file_path instanceof File) {
        formData.append("file_path", payload.file_path);
    } else if (typeof payload.file_path === "string") {
        formData.append("file_path", payload.file_path);
    }
}


    return formData;
}

export const getAllNewsPosts = async (): Promise<NewsPost[]> => {
    const response: any = await Get('/newsposts');
    if (Array.isArray(response)) return response;
    if (response && Array.isArray(response.data)) return response.data;
    return [];
}

interface NewsPostDetailResponse {
    news_post: NewsPost;
    features: any[]; // ใช้ any[] เพื่อความยืดหยุ่น ถ้าไม่ต้องการ import FeaturescholarshipResponse
}

/**
 * ดึงข้อมูลข่าวสารพร้อมคุณสมบัติทุนตาม ID
 * @param id ID ของข่าวสาร
 * @returns Promise<NewsPostDetailResponse> โครงสร้าง { news_post: NewsPost, features: [...] }
 */
export const getNewsPostById = async (id: number): Promise<NewsPostDetailResponse> => {
    // 1. เรียก API
    const response: any = await Get(`/newsposts/${id}`);
    
    // 2. ตรวจสอบโครงสร้าง Response ที่ Go Backend ส่งมา
    if (response && response.news_post) {
        // ถ้าโครงสร้างถูกต้อง (มี news_post อยู่) ให้ return ออกไป
        return response as NewsPostDetailResponse; 
    }
    
    // 3. ถ้าไม่พบ หรือเป็น Error ให้โยน Error
    throw new Error('Invalid or empty API response structure for NewsPost ID: ' + id);
}

export const createNewsPost = async (payload: CreateNewsPostPayload | FormData): Promise<NewsPost> => {
    const dataToSend = payload instanceof FormData ? payload : convertToFormData(payload);
    const response: any = await Post('/newsposts', dataToSend);
    return response.data;
}

export const updateNewsPost = async (id: number, data: UpdateNewsPost | FormData): Promise<any> => {
    const response: any = await Put(`/newsposts/${id}`, data); 
    return response;
};


export const deleteNewsPost = async (id: number): Promise<void> => {
    await Delete(`/newsposts/${id}`);
}