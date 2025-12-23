import type { StudentProfileResponse, FamilyInfo } from '@/interfaces/user';

/**
 * ฟิลด์บังคับสำหรับโปรไฟล์นักศึกษา
 */
const REQUIRED_STUDENT_FIELDS = [
  'student_id',
  'first_name_th',
  'last_name_th',
  'national_id',
  'email',
  'phone',
  'current_address',
  'province',
  'major_id',
  'current_year',
  'gpax',
] as const;

/**
 * ฟิลด์บังคับสำหรับข้อมูลครอบครัว
 */
const REQUIRED_FAMILY_FIELDS = [
  'father_name',
  'father_occupation',
  'mother_name',
  'mother_occupation',
] as const;

/**
 * ตรวจสอบว่าโปรไฟล์นักศึกษาครบถ้วนหรือไม่
 * @param student ข้อมูลนักศึกษา
 * @param family ข้อมูลครอบครัว
 * @returns object ประกอบด้วย isComplete, missingFields และ completionPercentage
 */
export const checkProfileCompleteness = (
  student: StudentProfileResponse | null | undefined,
  family: FamilyInfo | null | undefined
): {
  isComplete: boolean;
  missingFields: string[];
  completionPercentage: number;
} => {
  const missingFields: string[] = [];
  const totalFields = REQUIRED_STUDENT_FIELDS.length + REQUIRED_FAMILY_FIELDS.length;
  let completedFields = 0;

  if (!student) {
    return {
      isComplete: false,
      missingFields: ['ไม่พบข้อมูลโปรไฟล์'],
      completionPercentage: 0,
    };
  }

  // ตรวจสอบฟิลด์นักศึกษา
  for (const field of REQUIRED_STUDENT_FIELDS) {
    const value = (student as any)[field];
    if (value === null || value === undefined || value === '' || value === 0) {
      missingFields.push(getFieldLabel(field));
    } else {
      completedFields++;
    }
  }

  // ตรวจสอบฟิลด์ครอบครัว
  if (!family) {
    REQUIRED_FAMILY_FIELDS.forEach(field => {
      missingFields.push(getFieldLabel(field));
    });
  } else {
    for (const field of REQUIRED_FAMILY_FIELDS) {
      const value = (family as any)[field];
      if (value === null || value === undefined || value === '') {
        missingFields.push(getFieldLabel(field));
      } else {
        completedFields++;
      }
    }
  }

  return {
    isComplete: missingFields.length === 0,
    missingFields,
    completionPercentage: Math.round((completedFields / totalFields) * 100),
  };
};

/**
 * แปลงชื่อ field เป็นภาษาไทย
 */
const fieldLabels: Record<string, string> = {
  student_id: 'รหัสนักศึกษา',
  first_name_th: 'ชื่อ (ภาษาไทย)',
  last_name_th: 'นามสกุล (ภาษาไทย)',
  national_id: 'เลขบัตรประชาชน',
  email: 'อีเมล',
  phone: 'เบอร์โทรศัพท์',
  current_address: 'ที่อยู่ปัจจุบัน',
  province: 'จังหวัด',
  major_id: 'สาขาวิชา',
  current_year: 'ชั้นปีที่กำลังศึกษา',
  gpax: 'เกรดเฉลี่ยสะสม (GPAX)',
  father_name: 'ชื่อบิดา',
  father_occupation: 'อาชีพบิดา',
  mother_name: 'ชื่อมารดา',
  mother_occupation: 'อาชีพมารดา',
};

const getFieldLabel = (field: string): string => {
  return fieldLabels[field] || field;
};

/**
 * Routes ที่ต้องตรวจสอบโปรไฟล์ก่อนเข้า
 */
export const PROFILE_REQUIRED_ROUTES = [
  '/dashboard/apply-scholarship',
  '/dashboard/track-status',
  '/dashboard/schedule',
];
