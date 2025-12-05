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
  username: string;
  role_id?: number | null;
  role?: RoleResponse | null;
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
}
