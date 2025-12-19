import { Get, Post } from "./https"; // ไม่ต้องใช้ Delete แล้ว
import type { StudentFavoriteNewsPost } from "@/interfaces/news_post";

// 1. ดึงรายการที่ชอบ (GET)
// Backend Route: r.GET("/student_favs/my_favs/:student_id", ...)
export const getMyFavoriteNews = async (studentId: number): Promise<StudentFavoriteNewsPost[]> => {
    // แก้ URL ให้ตรงกับ Backend
    const response: any = await Get(`/student_favs/my_favs/${studentId}`);

    if (Array.isArray(response)) return response;
    if (response && Array.isArray(response.data)) return response.data;
    
    return [];
}

// 2. กดปุ่มหัวใจ (Toggle)
// Backend Route: r.POST("/student_favs/toggle", ...)
// รับแค่ studentId กับ newsId ก็พอ (ไม่ต้องใช้ ID ของตาราง Fav)
export const toggleStudentFavoriteNews = async (studentId: number, newsPostId: number): Promise<any> => {
    const payload = {
        // 🔥 จุดสำคัญ: ชื่อ field ต้องตรงกับ struct ใน Go (json:"student_profile_id")
        student_profile_id: studentId, 
        news_post_id: newsPostId
    };

    // ยิงไปที่ Path toggle
    const response: any = await Post(`/student_favs/toggle`, payload);
    
    // คืนค่า response กลับไป (Backend อาจจะส่ง message ว่า "Favorited" หรือ "Unfavorited")
    return response;
}
