<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { getScholarships, applyForScholarship } from '@/services/api/scholarship';
import { getStudentApplications } from '@/services/api/application';
import { getMyProfile, updateMyProfile } from '@/services/api/user';
import type { ScholarshipResponse, MyProfileResponse, StudentProfileResponse } from '@/interfaces';
import type { ApplicationScholarshipResponse } from '@/interfaces/application_scholarship';
import type { FamilyInfoResponse } from '@/interfaces/family_info';
import Swal from 'sweetalert2';

const scholarships = ref<ScholarshipResponse[]>([]);
const studentApplications = ref<ApplicationScholarshipResponse[]>([]);
const studentProfile = ref<MyProfileResponse | null>(null);
const isLoading = ref(true);
const error = ref<string | null>(null);
const router = useRouter();

// Modal state
const isApplicationModalOpen = ref(false);
const selectedScholarship = ref<ScholarshipResponse | null>(null);
const isSubmitting = ref(false);
const isSavingProfile = ref(false);

// Editable fields
const editableData = ref({
  email: '',
  phone: '',
  father_occupation: '',
  father_income: 0,
  mother_occupation: '',
  mother_income: 0,
  guardian_occupation: '',
  guardian_income: 0,
  application_reason: '' // เหตุผลในการสมัคร/เล่าเรื่องราว
});

// Get student profile data
const studentData = computed(() => {
  if (studentProfile.value?.role === 'student') {
    return studentProfile.value.data as StudentProfileResponse;
  }
  return null;
});

const familyData = computed(() => {
  if (studentProfile.value?.family) {
    return studentProfile.value.family as FamilyInfoResponse;
  }
  return null;
});

// Get student profile ID from logged-in user
const studentProfileId = computed(() => {
  return studentData.value?.ID || null;
});

// Check if scholarship is already applied
const isAlreadyApplied = (scholarshipId: number) => {
  return studentApplications.value.some(app => app.scholarship_id === scholarshipId);
};

const fetchScholarships = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    // Fetch user profile first
    studentProfile.value = await getMyProfile();
    
    // Fetch scholarships
    scholarships.value = await getScholarships();
    
    // Fetch student's existing applications
    if (studentProfileId.value) {
      studentApplications.value = await getStudentApplications(studentProfileId.value);
    }
    
    console.log('Data received for scholarships:', scholarships.value);
  } catch (err) {
    error.value = 'ไม่สามารถโหลดข้อมูลทุนการศึกษาได้';
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

// Initialize editable data from profile
const initEditableData = () => {
  if (studentData.value) {
    editableData.value.email = studentData.value.email || '';
    editableData.value.phone = studentData.value.phone || '';
  }
  if (familyData.value) {
    editableData.value.father_occupation = familyData.value.father_occupation || '';
    editableData.value.father_income = familyData.value.father_income || 0;
    editableData.value.mother_occupation = familyData.value.mother_occupation || '';
    editableData.value.mother_income = familyData.value.mother_income || 0;
    editableData.value.guardian_occupation = familyData.value.guardian_occupation || '';
    editableData.value.guardian_income = familyData.value.guardian_income || 0;
  }
  editableData.value.application_reason = '';
};

// Open application modal
const openApplicationModal = (scholarship: ScholarshipResponse) => {
  selectedScholarship.value = scholarship;
  initEditableData();
  isApplicationModalOpen.value = true;
};

// Close modal
const closeModal = () => {
  if (isSubmitting.value || isSavingProfile.value) return;
  isApplicationModalOpen.value = false;
  selectedScholarship.value = null;
};

// Save profile changes
const saveProfileChanges = async () => {
  if (!studentData.value) return false;
  
  isSavingProfile.value = true;
  try {
    const updatePayload = {
      email: editableData.value.email,
      phone: editableData.value.phone,
      first_name_en: studentData.value.first_name_en,
      last_name_en: studentData.value.last_name_en,
      birth_date: studentData.value.birth_date,
      permanent_address: studentData.value.permanent_address,
      current_address: studentData.value.current_address,
      province: studentData.value.province,
      siblings_count: studentData.value.siblings_count,
      family_info: {
        father_name: familyData.value?.father_name || '',
        father_occupation: editableData.value.father_occupation,
        father_income: editableData.value.father_income,
        mother_name: familyData.value?.mother_name || '',
        mother_occupation: editableData.value.mother_occupation,
        mother_income: editableData.value.mother_income,
        guardian_name: familyData.value?.guardian_name || '',
        guardian_relation: familyData.value?.guardian_relation || '',
        guardian_occupation: editableData.value.guardian_occupation,
        guardian_income: editableData.value.guardian_income,
      }
    };
    
    await updateMyProfile(updatePayload);
    return true;
  } catch (err) {
    console.error('Error saving profile:', err);
    return false;
  } finally {
    isSavingProfile.value = false;
  }
};

// Submit application
const submitApplication = async () => {
  if (!studentProfileId.value || !selectedScholarship.value) return;
  
  isSubmitting.value = true;
  try {
    // Save profile changes first
    const saved = await saveProfileChanges();
    if (!saved) {
      Swal.fire('ข้อผิดพลาด', 'ไม่สามารถบันทึกข้อมูลโปรไฟล์ได้', 'error');
      return;
    }

    // Submit application
    const apiResult = await applyForScholarship(selectedScholarship.value.ID, studentProfileId.value);
    console.log('Application result:', apiResult);
    
    closeModal();
    
    await Swal.fire({
      title: 'สมัครสำเร็จ!',
      text: 'ระบบจะนำคุณไปยังหน้าติดตามสถานะ',
      icon: 'success',
      confirmButtonText: 'ตกลง',
      confirmButtonColor: '#1e3a8a',
      timer: 2000,
      timerProgressBar: true,
    });
    
    router.push('/dashboard/track-status');
  } catch (err: any) {
    const msg = err?.response?.data?.error || 'เกิดข้อผิดพลาดในการสมัครทุน หรือคุณอาจเคยสมัครทุนนี้ไปแล้ว';
    Swal.fire('ข้อผิดพลาด', msg, 'error');
    console.error(err);
  } finally {
    isSubmitting.value = false;
  }
};

// Calculate total family income
const totalFamilyIncome = computed(() => {
  return (editableData.value.father_income || 0) + 
         (editableData.value.mother_income || 0) + 
         (editableData.value.guardian_income || 0);
});

onMounted(fetchScholarships);
</script>

<template>
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible font-sans text-slate-800" data-theme="light">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-[#1e3a8a] mb-1 flex items-center gap-3">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
          stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
        </svg>
        ยื่นสมัครทุนการศึกษา
      </h1>
      <p class="text-gray-500 text-sm">เลือกทุนการศึกษาที่คุณต้องการสมัคร</p>
    </div>

      <div v-if="isLoading" class="text-center py-20 text-gray-500">
        <span class="loading loading-spinner loading-lg"></span>
        <p class="mt-2">กำลังโหลดรายการทุน...</p>
      </div>

      <div v-else-if="error" class="text-center py-20 text-red-500">
        <p>{{ error }}</p>
        <button @click="fetchScholarships" class="btn btn-sm btn-outline mt-4">ลองใหม่อีกครั้ง</button>
      </div>

      <div v-else class="grid grid-cols-1 gap-5">
        <div v-for="scholarship in scholarships" :key="scholarship.ID"
          class="card bg-white shadow-sm hover:shadow-lg transition-all border border-transparent rounded-2xl overflow-hidden group"
          :class="{'ring-2 ring-green-400 ring-offset-2': isAlreadyApplied(scholarship.ID)}">

          <div class="card-body p-5">
            <div class="flex flex-col md:flex-row gap-4">
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2 flex-wrap">
                  <span class="bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-1 rounded-full">{{
                    scholarship.typescholarship?.type_name || 'ทั่วไป' }}</span>
                  <span v-if="scholarship.semaster" class="bg-gray-100 text-gray-800 text-xs font-semibold px-2.5 py-1 rounded-full">
                    รอบ: {{ scholarship.semaster.term }}/{{ scholarship.semaster.academic_year }}
                  </span>
                  <span v-if="isAlreadyApplied(scholarship.ID)" class="bg-green-100 text-green-800 text-xs font-bold px-2.5 py-1 rounded-full flex items-center gap-1">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                    สมัครแล้ว
                  </span>
                  <span class="text-xs text-gray-400">ปิดรับสมัคร: {{ new
                    Date(scholarship.close_date).toLocaleDateString('th-TH') }}</span>
                </div>
                <h2 class="font-bold text-lg text-slate-800 group-hover:text-[#1e3a8a] transition-colors">
                  {{ scholarship.scholarship_name }}
                </h2>
                <p class="text-sm text-gray-600 mt-1 line-clamp-2">
                  {{ scholarship.description }}
                </p>
              </div>
              <div
                class="flex-shrink-0 flex flex-col items-center justify-center gap-2 border-l-0 md:border-l md:pl-4 mt-4 md:mt-0 pt-4 md:pt-0 border-t md:border-t-0">
                <button v-if="!isAlreadyApplied(scholarship.ID)" @click="openApplicationModal(scholarship)"
                  class="btn btn-primary btn-sm bg-[#1e3a8a] hover:bg-blue-800 border-none text-white rounded-full px-5">
                  สมัครทุนนี้
                </button>
                <router-link v-else to="/dashboard/track-status"
                  class="btn btn-sm bg-green-600 hover:bg-green-700 text-white border-none rounded-full px-5">
                  ดูสถานะ
                </router-link>
              </div>
            </div>
          </div>
        </div>
        <div v-if="scholarships.length === 0" class="text-center py-20 text-gray-500">
          <p>ไม่พบทุนการศึกษาที่เปิดรับสมัครในขณะนี้</p>
        </div>
      </div>
  </div>

  <!-- Application Modal -->
  <div v-if="isApplicationModalOpen && selectedScholarship" 
       class="fixed inset-0 z-[200] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
    <div class="bg-white w-full max-w-4xl max-h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
      
      <!-- Header -->
      <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50 shrink-0">
        <div>
          <h2 class="text-xl font-bold text-[#1e3a8a]">ยืนยันการสมัครทุน</h2>
          <p class="text-sm text-gray-500">{{ selectedScholarship.scholarship_name }}</p>
        </div>
        <button @click="closeModal" :disabled="isSubmitting" class="btn btn-circle btn-ghost btn-sm hover:bg-slate-200">✕</button>
      </div>

      <!-- Body -->
      <div class="flex-1 overflow-y-auto p-6 space-y-6">
        
        <!-- Scholarship Info -->
        <div class="bg-blue-50 border border-blue-100 rounded-xl p-4">
          <h3 class="font-bold text-blue-800 mb-2">ข้อมูลทุนการศึกษา</h3>
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div><span class="text-gray-500">ประเภท:</span> <span class="font-medium">{{ selectedScholarship.typescholarship?.type_name || 'ทั่วไป' }}</span></div>
            <div><span class="text-gray-500">รอบ:</span> <span class="font-medium">{{ selectedScholarship.semaster?.term }}/{{ selectedScholarship.semaster?.academic_year }}</span></div>
            <div class="col-span-2"><span class="text-gray-500">ปิดรับสมัคร:</span> <span class="font-medium">{{ new Date(selectedScholarship.close_date).toLocaleDateString('th-TH') }}</span></div>
          </div>
        </div>

        <!-- Student Profile (Read-only) -->
        <div class="bg-slate-50 border border-slate-200 rounded-xl p-4">
          <h3 class="font-bold text-slate-700 mb-3 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            ข้อมูลนักศึกษา (ไม่สามารถแก้ไขได้)
          </h3>
          <div class="grid grid-cols-2 md:grid-cols-3 gap-4 text-sm">
            <div><span class="text-gray-500">รหัสนักศึกษา:</span> <span class="font-medium">{{ studentData?.student_id }}</span></div>
            <div><span class="text-gray-500">ชื่อ-สกุล:</span> <span class="font-medium">{{ studentData?.first_name_th }} {{ studentData?.last_name_th }}</span></div>
            <div><span class="text-gray-500">สาขา:</span> <span class="font-medium">{{ studentData?.major?.major_name || '-' }}</span></div>
            <div><span class="text-gray-500">ชั้นปี:</span> <span class="font-medium">{{ studentData?.current_year }}</span></div>
            <div><span class="text-gray-500">GPAX:</span> <span class="font-bold text-blue-600">{{ studentData?.gpax?.toFixed(2) }}</span></div>
            <div><span class="text-gray-500">จำนวนพี่น้อง:</span> <span class="font-medium">{{ studentData?.siblings_count }} คน</span></div>
          </div>
        </div>

        <!-- Editable Contact Info -->
        <div class="bg-white border border-gray-200 rounded-xl p-4">
          <h3 class="font-bold text-slate-700 mb-3 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            ข้อมูลติดต่อ (แก้ไขได้)
          </h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="form-control">
              <label class="label py-1"><span class="label-text font-medium">อีเมล</span></label>
              <input v-model="editableData.email" type="email" class="input input-bordered input-sm bg-white" placeholder="example@email.com" />
            </div>
            <div class="form-control">
              <label class="label py-1"><span class="label-text font-medium">เบอร์โทรศัพท์</span></label>
              <input v-model="editableData.phone" type="tel" class="input input-bordered input-sm bg-white" placeholder="0812345678" />
            </div>
          </div>
        </div>

        <!-- Editable Family Income -->
        <div class="bg-white border border-gray-200 rounded-xl p-4">
          <h3 class="font-bold text-slate-700 mb-3 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            ข้อมูลครอบครัว (แก้ไขได้)
          </h3>
          
          <!-- Father -->
          <div class="mb-4 p-3 bg-slate-50 rounded-lg">
            <p class="font-medium text-sm text-slate-600 mb-2">บิดา: {{ familyData?.father_name || '-' }}</p>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div class="form-control">
                <label class="label py-0"><span class="label-text text-xs">อาชีพ</span></label>
                <input v-model="editableData.father_occupation" type="text" class="input input-bordered input-sm bg-white" />
              </div>
              <div class="form-control">
                <label class="label py-0"><span class="label-text text-xs">รายได้ (บาท/เดือน)</span></label>
                <input v-model.number="editableData.father_income" type="number" min="0" class="input input-bordered input-sm bg-white" />
              </div>
            </div>
          </div>

          <!-- Mother -->
          <div class="mb-4 p-3 bg-slate-50 rounded-lg">
            <p class="font-medium text-sm text-slate-600 mb-2">มารดา: {{ familyData?.mother_name || '-' }}</p>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div class="form-control">
                <label class="label py-0"><span class="label-text text-xs">อาชีพ</span></label>
                <input v-model="editableData.mother_occupation" type="text" class="input input-bordered input-sm bg-white" />
              </div>
              <div class="form-control">
                <label class="label py-0"><span class="label-text text-xs">รายได้ (บาท/เดือน)</span></label>
                <input v-model.number="editableData.mother_income" type="number" min="0" class="input input-bordered input-sm bg-white" />
              </div>
            </div>
          </div>

          <!-- Guardian -->
          <div class="p-3 bg-slate-50 rounded-lg">
            <p class="font-medium text-sm text-slate-600 mb-2">ผู้ปกครอง: {{ familyData?.guardian_name || '-' }} ({{ familyData?.guardian_relation || '-' }})</p>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div class="form-control">
                <label class="label py-0"><span class="label-text text-xs">อาชีพ</span></label>
                <input v-model="editableData.guardian_occupation" type="text" class="input input-bordered input-sm bg-white" />
              </div>
              <div class="form-control">
                <label class="label py-0"><span class="label-text text-xs">รายได้ (บาท/เดือน)</span></label>
                <input v-model.number="editableData.guardian_income" type="number" min="0" class="input input-bordered input-sm bg-white" />
              </div>
            </div>
          </div>

          <!-- Total Income Summary -->
          <div class="mt-4 p-3 bg-blue-50 border border-blue-100 rounded-lg flex justify-between items-center">
            <span class="font-medium text-blue-800">รายได้รวมครอบครัว:</span>
            <span class="font-bold text-lg text-blue-900">{{ totalFamilyIncome.toLocaleString('th-TH') }} บาท/เดือน</span>
          </div>
        </div>

        <!-- Application Reason -->
        <div class="bg-white border border-gray-200 rounded-xl p-4">
          <h3 class="font-bold text-slate-700 mb-3 flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            เหตุผลในการสมัครทุน / เล่าเรื่องราวของคุณ
          </h3>
          <textarea 
            v-model="editableData.application_reason"
            class="textarea textarea-bordered w-full h-32 bg-white"
            placeholder="กรุณาเล่าถึงความจำเป็นที่ต้องการทุนการศึกษา ความมุ่งมั่น หรือเรื่องราวที่คุณอยากให้คณะกรรมการทราบ..."
          ></textarea>
        </div>

      </div>

      <!-- Footer -->
      <div class="px-6 py-4 border-t bg-slate-50 flex justify-end gap-3 shrink-0">
        <button @click="closeModal" :disabled="isSubmitting" class="btn btn-ghost">
          ยกเลิก
        </button>
        <button @click="submitApplication" :disabled="isSubmitting" 
                class="btn bg-[#1e3a8a] hover:bg-[#152c6f] text-white border-none">
          <span v-if="isSubmitting" class="loading loading-spinner loading-xs"></span>
          <span v-else>ยืนยันสมัครทุน</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes pop-in {
  0% { opacity: 0; transform: scale(0.95) translateY(10px); }
  100% { opacity: 1; transform: scale(1) translateY(0); }
}
.animate-pop-in { animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards; }
</style>