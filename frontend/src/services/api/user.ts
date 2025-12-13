import { Get, Post, Put, Delete } from './https';
import type { 
  UserResponse, 
  RoleResponse, 
  MajorResponse, 
  CreateUserPayload, 
  MyProfileResponse 
} from '@/interfaces/user'; // ตรวจสอบ path import ให้ตรงกับโปรเจกต์คุณ

/**
 * ดึงข้อมูล Users ทั้งหมด (สำหรับ Admin)
 */
export const getUsers = async (): Promise<UserResponse[]> => {
  const res = await Get('/users');
  return res.data as UserResponse[];
};

/**
 * ดึงข้อมูล Roles ทั้งหมด (สำหรับ Dropdown)
 */
export const getRoles = async (): Promise<RoleResponse[]> => {
  const res = await Get('/roles');
  return res.data as RoleResponse[];
};

/**
 * ดึงข้อมูล Majors ทั้งหมด (สำหรับ Dropdown)
 */
export const getMajors = async (): Promise<MajorResponse[]> => {
  const res = await Get('/majors');
  return res.data as MajorResponse[];
};

/**
 * สร้าง User ใหม่ (Admin Only)
 */
export const createUser = async (data: CreateUserPayload) => {
  return await Post('/users', data);
};

/**
 * ลบ User (Admin Only)
 */
export const deleteUser = async (id: number) => {
  return await Delete(`/users/${id}`);
};

/**
 * ดึงข้อมูล Profile ส่วนตัว (Student/Admin)
 */
export const getMyProfile = async (): Promise<MyProfileResponse> => {
  const res = await Get('/profile/me');
  // Backend ส่งกลับมาเป็น { role: "...", data: { ... }, family: { ... } }
  return res as MyProfileResponse; 
};

/**
 * อัปเดตข้อมูล Profile ส่วนตัว
 * (ส่งไปทั้ง Object, Backend จะเลือก update เฉพาะ field ที่อนุญาตเอง)
 */
export const updateMyProfile = async (data: any) => {
  return await Put('/profile/me', data);
};

export const updateUser = async (id: number, data: any) => {
  return await Put(`/users/${id}`, data);
};