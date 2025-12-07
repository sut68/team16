<script setup lang="ts">
import { ref, onMounted, type Ref } from 'vue';
import TrackingModal from './TrackingModal.vue';
import { getStudentApplications, uploadDocument } from '@/services/api/application';
import type { ApplicationScholarshipResponse } from '@/interfaces';

const studentApplications: Ref<ApplicationScholarshipResponse[]> = ref([]);
const isLoading = ref(true);
const error = ref<string | null>(null);

const fetchStudentApplications = async () => {
  isLoading.value = true;
  error.value = null;
  // In a real app, studentProfileId would come from a global auth store.
  // We are hardcoding it to 1 for demonstration purposes.
  const studentProfileId = 1; 
  try {
    const data = await getStudentApplications(studentProfileId);
    studentApplications.value = data;
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
  if (!selectedApp.value) return;

  // In a real app, studentProfileId would come from a global auth store.
  const studentProfileId = 1;

  try {
    await uploadDocument(selectedApp.value.ID, file, studentProfileId);
    alert('อัปโหลดไฟล์สำเร็จ! สถานะจะถูกอัปเดตในการโหลดครั้งถัดไป');
    
    // Refetch data to show changes
    await fetchStudentApplications();
    isModalOpen.value = false;

  } catch (err) {
    alert('เกิดข้อผิดพลาดระหว่างการอัปโหลด');
    console.error(err);
  }
};
</script>

<template>
  <div class="min-h-screen bg-[#f0f2f5] p-6 font-sans text-slate-800">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-2xl font-bold text-[#1e3a8a] mb-6 flex items-center gap-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" /></svg>
        ติดตามสถานะการสมัครทุน
      </h1>

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
          @click="openTracking(app)"
          class="card bg-white shadow-sm hover:shadow-md transition-all cursor-pointer border border-transparent hover:border-blue-200 group rounded-xl overflow-hidden">
          
          <div class="card-body p-5 flex flex-col md:flex-row items-center justify-between gap-4">
            <div v-if="app.scholarship" class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <span class="bg-blue-50 text-blue-700 text-xs px-2 py-0.5 rounded border border-blue-100 font-semibold">APP-{{ app.ID }}</span>
                <span class="text-xs text-gray-400">ยื่นเมื่อ: {{ new Date(app.CreatedAt).toLocaleDateString('th-TH') }}</span>
              </div>
              <h3 class="font-bold text-lg text-slate-800 group-hover:text-[#1e3a8a] transition-colors">
                {{ app.scholarship?.scholarship_name }}
              </h3>
              <p class="text-sm text-gray-500">
                {{ app.scholarship?.typescholarship?.type_name }}
              </p>
            </div>
            <div v-else class="flex-1">
              <p class="text-gray-500">ไม่สามารถโหลดข้อมูลทุนได้</p>
            </div>

            <div class="flex items-center gap-4 w-full md:w-auto justify-between md:justify-end">
              <div class="text-right">
                <p class="text-xs text-gray-400 mb-1">สถานะปัจจุบัน</p>
                <span v-if="app.status === 'approved'" class="badge badge-success text-white gap-1">
                  อนุมัติแล้ว
                </span>
                <span v-else-if="app.status === 'rejected'" class="badge badge-error text-white gap-1">
                  ไม่ผ่านการพิจารณา
                </span>
                <span v-else class="badge badge-info bg-blue-100 text-blue-800 border-none gap-1">
                  {{ app.status }}
                </span>
              </div>
              
              <div class="bg-slate-50 p-2 rounded-full group-hover:bg-blue-50 text-gray-400 group-hover:text-[#1e3a8a]">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
              </div>
            </div>
          </div>
          
          <!-- The progress bar logic is removed as it's too complex with the new structure without more data -->
        </div>
      </div>
       <div v-else class="text-center py-20 text-gray-500">
          <p>คุณยังไม่มีรายการสมัครทุน</p>
          <router-link to="/dashboard/apply-scholarship" class="btn btn-primary btn-sm mt-4 bg-[#1e3a8a]">ไปที่หน้าสมัครทุน</router-link>
        </div>
    </div>

    <TrackingModal 
      v-if="isModalOpen" 
      :isOpen="isModalOpen" 
      :applicationData="selectedApp"
      @close="isModalOpen = false" 
      @upload-file="handleFileUpload"
    />
  </div>
</template>
