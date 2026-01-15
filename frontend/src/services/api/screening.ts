import { Get, Post, Put, Delete } from './https';
import type {
  ScreeningResponse,
  CreateScreeningPayload,
  UpdateScreeningStatusPayload
} from './../../interfaces/screening';

/**
 * Service สำหรับจัดการข้อมูลการคัดกรอง (Screening)
 */

/**
 * ดึงข้อมูลการคัดกรองทั้งหมด
 * GET /screenings
 */
export const getAllScreenings = async (): Promise<ScreeningResponse[]> => {
  try {
    const res: any = await Get('/screening');

    // console.log('Debug API Response:', res);

    if (!res) return [];
    if (Array.isArray(res)) return res;
    if (Array.isArray(res.data)) return res.data;
    return [];

  } catch (error) {
    console.error("Error in getAllScreenings:", error);
    return [];
  }
};

/**
 * ดึงข้อมูลการคัดกรองตาม ID
 * GET /screenings/:id
 */
export const getScreeningById = async (id: number): Promise<ScreeningResponse> => {

  const response: any = await Get(`/screening/${id}`);
  return response.data;
};

/**
 * สร้างการคัดกรองใหม่
 * POST /screenings
 */
export const createScreening = async (payload: CreateScreeningPayload): Promise<ScreeningResponse> => {
  // 🔥 แก้ path: เติม s
  const response: any = await Post('/screening', payload);
  return response.data;
};

/**
 * อัปเดตสถานะการคัดกรอง
 * PUT /screenings/:id
 */
export const updateScreeningStatus = async (
  id: number,
  payload: UpdateScreeningStatusPayload
): Promise<ScreeningResponse> => {

  const response: any = await Put(`/screening/${id}`, payload);
  return response.data;
};


/**
 * ลบการคัดกรอง
 * DELETE /screenings/:id
 */
export const deleteScreening = async (screeningId: number): Promise<void> => {
  // 🔥 แก้ path: เติม s
  await Delete(`/screening/${screeningId}`);
};

// ==================== AUTO SCREENING ====================

export interface AutoScreenResult {
  screening_id: number;
  application_id: number;
  scholarship_id: number;
  scholarship_name: string;
  student_name: string;
  auto_approved: boolean;
  passed_criteria: number;
  total_criteria: number;
  failed_criteria?: string[];
  processed_by_admin?: string;
}

export interface AutoScreenBatchResponse {
  data: AutoScreenResult[];
  summary: {
    total: number;
    auto_approved: number;
    needs_review: number;
    scholarship: string;
  };
  message: string;
}

/**
 * คัดกรองอัตโนมัติ 1 คน
 * POST /screening/:id/auto-screen
 */
export const autoScreenSingle = async (screeningId: number): Promise<{
  data: AutoScreenResult;
  message: string;
}> => {
  const response: any = await Post(`/screening/${screeningId}/auto-screen`, {});
  return response;
};

/**
 * คัดกรองอัตโนมัติทั้งทุน (Batch)
 * POST /scholarship/:id/auto-screen-batch
 */
export const autoScreenBatch = async (scholarshipId: number): Promise<AutoScreenBatchResponse> => {
  const response: any = await Post(`/scholarship/${scholarshipId}/auto-screen-batch`, {});
  return response;
};