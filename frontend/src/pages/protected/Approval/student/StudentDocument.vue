<script setup lang="ts">
import { ref, onMounted, computed, type Ref } from 'vue';
import TrackingModal from './TrackingModal.vue';
import { getStudentApplications, uploadDocument, cancelApplicationScholarship } from '@/services/api/application';
import { getMyProfile } from '@/services/api/user';
import type { ApplicationScholarshipResponse, MyProfileResponse } from '@/interfaces';
import Swal from 'sweetalert2';

const studentApplications: Ref<ApplicationScholarshipResponse[]> = ref([]);
const studentProfile = ref<MyProfileResponse | null>(null);
const isLoading = ref(true);
const error = ref<string | null>(null);

// Get student profile ID from logged-in user
const studentProfileId = computed(() => {
  if (studentProfile.value?.role === 'student') {
    return (studentProfile.value.data as any).ID;
  }
  return null;
});

const fetchStudentApplications = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    // Fetch user profile first
    studentProfile.value = await getMyProfile();
    
    if (studentProfileId.value) {
      const data = await getStudentApplications(studentProfileId.value);
      studentApplications.value = data;
    }
  } catch (err) {
    error.value = 'ไม่สามารถโหลดข้อมูลการสมัครทุนได้';
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchStudentApplications);


const isModalOpen = ref(false);
const selectedApp = ref<ApplicationScholarshipResponse | null>(null);

const openTracking = (app: ApplicationScholarshipResponse) => {
  selectedApp.value = app;
  isModalOpen.value = true;
};

const handleFileUpload = async (file: File) => {
  if (!selectedApp.value || !studentProfileId.value) return;

  try {
    await uploadDocument(selectedApp.value.ID, file, studentProfileId.value);
    await fetchStudentApplications();
    isModalOpen.value = false;

  } catch (err) {
    console.error(err);
  }
};

// Check if cancellation is allowed for an application
const canCancel = (app: ApplicationScholarshipResponse): boolean => {
  const allowedStatuses = ['new', 'pending', ''];
  return allowedStatuses.includes(app.status);
};

// Handle cancel application
const handleCancelApplication = async (app: ApplicationScholarshipResponse, event: Event) => {
  event.stopPropagation(); // Prevent opening the modal
  
  const result = await Swal.fire({
    title: 'ยืนยันการยกเลิก',
    html: `
      <div class="text-left">
        <p class="mb-2">คุณต้องการยกเลิกการสมัครทุน</p>
        <p class="font-bold text-blue-800">"${app.scholarship?.scholarship_name}"</p>
        <p class="text-sm text-red-500 mt-3">⚠️ การดำเนินการนี้ไม่สามารถย้อนกลับได้</p>
      </div>
    `,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'ยืนยันยกเลิก',
    cancelButtonText: 'ไม่ยกเลิก',
    confirmButtonColor: '#dc2626',
    cancelButtonColor: '#6b7280',
    reverseButtons: true,
  });

  if (!result.isConfirmed) return;

  try {
    await cancelApplicationScholarship(app.ID);
    
    await Swal.fire({
      title: 'ยกเลิกสำเร็จ',
      text: 'การสมัครทุนถูกยกเลิกเรียบร้อยแล้ว',
      icon: 'success',
      timer: 2000,
      showConfirmButton: false,
    });

    // Refresh the list
    await fetchStudentApplications();
  } catch (err: any) {
    const msg = err?.response?.data?.error || err?.message || 'เกิดข้อผิดพลาดในการยกเลิก';
    Swal.fire('ข้อผิดพลาด', msg, 'error');
    console.error(err);
  }
};

// Get status display text
const getStatusDisplay = (status: string) => {
  const statusMap: Record<string, { text: string; class: string }> = {
    'new': { text: 'รอดำเนินการ', class: 'badge-info bg-blue-100 text-blue-800 border-none' },
    'pending': { text: 'รอตรวจสอบ', class: 'badge-info bg-yellow-100 text-yellow-800 border-none' },
    'screening': { text: 'กำลังคัดกรอง', class: 'badge-warning bg-orange-100 text-orange-800 border-none' },
    'approved': { text: 'อนุมัติแล้ว', class: 'badge-success text-white' },
    'rejected': { text: 'ไม่ผ่านการพิจารณา', class: 'badge-error text-white' },
    'interview': { text: 'นัดสัมภาษณ์', class: 'badge-info bg-purple-100 text-purple-800 border-none' },
  };
  return statusMap[status] || { text: status, class: 'badge-info bg-gray-100 text-gray-800 border-none' };
};
</script>

<template>
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible font-sans text-slate-800" data-theme="light">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-[#1e3a8a] mb-1 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
        </svg>
        ติดตามสถานะการสมัครทุน
      </h1>
      <p class="text-gray-500 text-sm">ตรวจสอบความคืบหน้าของใบสมัครทุนการศึกษา</p>
    </div>

      <div v-if="isLoading" class="text-center py-20 text-gray-500">
        <span class="loading loading-spinner loading-lg"></span>
        <p class="mt-2">กำลังโหลดข้อมูล...</p>
      </div>

      <div v-else-if="error" class="text-center py-20 text-red-500">
        <p>{{ error }}</p>
        <button @click="fetchStudentApplications" class="btn btn-sm btn-outline mt-4">ลองใหม่อีกครั้ง</button>
      </div>

      <div v-else-if="studentApplications.length > 0" class="grid grid-cols-1 gap-4">
        <div v-for="app in studentApplications" :key="app.ID"
          class="card bg-white shadow-sm hover:shadow-md transition-all cursor-pointer border border-transparent hover:border-blue-200 group rounded-xl overflow-hidden">

          <div class="card-body p-5 flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
            <div v-if="app.scholarship" class="flex-1" @click="openTracking(app)">
              <div class="flex items-center gap-2 mb-1">
                <span
                  class="bg-blue-50 text-blue-700 text-xs px-2 py-0.5 rounded border border-blue-100 font-semibold">APP-{{
                  app.ID }}</span>
                <span class="text-xs text-gray-400">ยื่นเมื่อ: {{ new Date(app.CreatedAt).toLocaleDateString('th-TH')
                  }}</span>
              </div>
              <h3 class="font-bold text-lg text-slate-800 group-hover:text-[#1e3a8a] transition-colors">
                {{ app.scholarship?.scholarship_name }}
              </h3>
              <p class="text-sm text-gray-500">
                {{ app.scholarship?.typescholarship?.type_name }}
              </p>
            </div>
            <div v-else class="flex-1" @click="openTracking(app)">
              <p class="text-gray-500">ไม่สามารถโหลดข้อมูลทุนได้</p>
            </div>

            <div class="flex items-center gap-3 w-full md:w-auto justify-between md:justify-end">
              <!-- Status Badge -->
              <div class="text-right">
                <p class="text-xs text-gray-400 mb-1">สถานะปัจจุบัน</p>
                <span :class="['badge gap-1', getStatusDisplay(app.status).class]">
                  {{ getStatusDisplay(app.status).text }}
                </span>
              </div>

              <!-- Cancel Button (only show for cancellable statuses) -->
              <button 
                v-if="canCancel(app)" 
                @click="handleCancelApplication(app, $event)"
                class="btn btn-sm btn-outline btn-error hover:bg-red-600 hover:border-red-600 hover:text-white"
                title="ยกเลิกการสมัคร"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
                ยกเลิก
              </button>

              <!-- Arrow Icon -->
              <div @click="openTracking(app)" class="bg-slate-50 p-2 rounded-full group-hover:bg-blue-50 text-gray-400 group-hover:text-[#1e3a8a]">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24"
                  stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="text-center py-20 text-gray-500">
        <p>คุณยังไม่มีรายการสมัครทุน</p>
        <router-link to="/dashboard/apply-scholarship"
          class="btn btn-primary btn-sm mt-4 bg-[#1e3a8a]">ไปที่หน้าสมัครทุน</router-link>
      </div>

    <TrackingModal v-if="isModalOpen" :isOpen="isModalOpen" :applicationData="selectedApp" @close="isModalOpen = false"
      @upload-file="handleFileUpload" />
  </div>
</template>