export interface FamilyInfoResponse {
  ID: number;

  student_profile_id: number;
  father_name: string;
  father_age: number;
  father_occupation: string;
  father_income: number;
  
  mother_name: string;
  mother_age: number;
  mother_occupation: string;
  mother_income: number;

  guardian_name: string;
  guardian_age: number;
  guardian_occupation: string;
  guardian_income: number;
  guardian_relation: string;
  
}
