import type { FamilyInfoResponse } from "./family_info";
export interface LoginUserRequest {
  username: string;
  password: string;
}

export interface CreateUserRequest {
  username: string;
  password: string;
}


export interface RoleResponse {
  ID: number;
  name: string;
}

export interface UserResponse {
  ID: number;
  CreatedAt?: string;
  username: string;
  role_id?: number | null;
  role?: RoleResponse | null;
  student_profile?: StudentProfileResponse[];
}

export interface AdminProfileResponse {
  ID: number;
  admin_firstname: string;
  admin_lastname: string;
  position: number;
  department: number;
  email: string;
  user_id: number;
  user: UserResponse;
}

export interface MajorResponse {
  ID: number;
  major_name: string;
}

export interface StudentProfileResponse {
  ID: number;
  student_id: string;
  first_name_th: string;
  last_name_th: string;
  first_name_en: string;
  last_name_en: string;
  national_id: string;
  birth_date: string;
  current_year: number;
  gpax: number;
  advisor_name: string;
  phone: string;
  email: string;
  permanent_address: string;
  current_address: string;
  province: string;
  siblings_count: number;
  user_id: number;
  user: UserResponse;
  major_id: number;
  major: MajorResponse | null;

  //เพิ่ม FamilyInfo เข้ามา
  family_info: FamilyInfoResponse[];
}


// --- เพิ่มต่อท้ายไฟล์เดิม ---

export interface FamilyInfo {
  ID?: number;
  father_name?: string;
  father_occupation?: string;
  father_income?: number;
  mother_name?: string;
  mother_occupation?: string;
  mother_income?: number;
  guardian_name?: string;
  guardian_occupation?: string;
  guardian_income?: number;
  guardian_relation?: string;
  profile_id?: number;
}

// ใช้สำหรับฟอร์มสร้าง User ใหม่ (Admin Dashboard)
export interface CreateUserPayload {
  username: string;
  password?: string;
  role_id: number;

  // Student Fields
  student_id?: string;
  first_name_th?: string;
  last_name_th?: string;
  first_name_en?: string;
  last_name_en?: string;
  national_id?: string;
  major_id?: number;
  gpax?: number;
  advisor_name?: string;

  // Admin Fields
  admin_firstname?: string;
  admin_lastname?: string;
  position?: string; // แก้เป็น string ให้ตรงกับ Backend

  // Common Fields
  email?: string;
  phone?: string;
}

// ใช้สำหรับ Response หน้า Profile (รวม Role + Data)
export interface MyProfileResponse {
  role: 'student' | 'admin';
  data: StudentProfileResponse | AdminProfileResponse;
  family?: FamilyInfo; // มีเฉพาะ Student
}