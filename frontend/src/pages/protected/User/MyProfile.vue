<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import * as UserAPI from '@/services/api/user';
import StudentProfile from '@/components/ui/ProfileStudent.vue';
import AdminProfile from '@/components/ui/ProfileAdmin.vue';
import Swal from 'sweetalert2';
import { checkProfileCompleteness } from '@/utils/profileValidator';
import { clearProfileCache } from '@/routes/index';

const route = useRoute();
const role = ref('');
const loading = ref(true);
const profileData = ref<any>(null);
const familyData = ref<any>(null);

// Check if redirected with required flag
const isProfileRequired = computed(() => route.query.required === 'true');

// Profile completeness
const profileStatus = computed(() => {
  if (!profileData.value) return { isComplete: true, missingFields: [], completionPercentage: 0 };
  return checkProfileCompleteness(profileData.value, familyData.value);
});

const loadProfile = async () => {
  loading.value = true;
  try {
    const res = await UserAPI.getMyProfile();
    if (res) {
      role.value = res.role;
      familyData.value = res.family || (res.data as any).family_info || null;
      // รวม data และ family_info ให้ Component ลูก
      profileData.value = { 
        ...res.data, 
        family_info: familyData.value || {} 
      };
    }
  } catch (err) {
    console.error(err);
    Swal.fire('ข้อผิดพลาด', 'ไม่สามารถโหลดข้อมูลโปรไฟล์ได้', 'error');
  } finally {
    loading.value = false;
  }
};

// Refresh and clear cache
const handleRefresh = async () => {
  clearProfileCache();
  await loadProfile();
};

onMounted(loadProfile);
</script>

<template>
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-y-auto" data-theme="light">
    
    <!-- Profile Required Alert -->
    <div v-if="isProfileRequired && !profileStatus.isComplete" 
         class="mb-6 bg-gradient-to-r from-amber-50 to-orange-50 border-l-4 border-amber-500 p-5 rounded-r-xl shadow-sm">
      <div class="flex items-start gap-4">
        <div class="shrink-0 bg-amber-100 p-3 rounded-full">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <div class="flex-1">
          <h3 class="font-bold text-lg text-amber-800 mb-2">กรุณากรอกข้อมูลให้ครบถ้วน</h3>
          <p class="text-sm text-amber-700 mb-3">
            คุณต้องกรอกข้อมูลโปรไฟล์ให้ครบถ้วนก่อนจึงจะสามารถสมัครทุนการศึกษาได้
          </p>
          
          <!-- Progress Bar -->
          <div class="mb-3">
            <div class="flex items-center justify-between text-xs mb-1">
              <span class="text-amber-600 font-medium">ความสมบูรณ์ของข้อมูล</span>
              <span class="text-amber-800 font-bold">{{ profileStatus.completionPercentage }}%</span>
            </div>
            <div class="w-full bg-amber-200 rounded-full h-2.5">
              <div class="bg-amber-600 h-2.5 rounded-full transition-all duration-500" 
                   :style="{ width: profileStatus.completionPercentage + '%' }"></div>
            </div>
          </div>

          <!-- Missing Fields -->
          <div class="bg-white/60 rounded-lg p-3 border border-amber-200">
            <p class="text-xs text-amber-700 font-semibold mb-2">ข้อมูลที่ยังไม่ได้กรอก:</p>
            <div class="flex flex-wrap gap-2">
              <span v-for="(field, idx) in profileStatus.missingFields" :key="idx"
                    class="inline-flex items-center px-2.5 py-1 bg-red-100 text-red-700 text-xs rounded-full">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
                {{ field }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Profile Completion Success -->
    <div v-if="isProfileRequired && profileStatus.isComplete && !loading" 
         class="mb-6 bg-gradient-to-r from-green-50 to-emerald-50 border-l-4 border-green-500 p-5 rounded-r-xl shadow-sm">
      <div class="flex items-center gap-4">
        <div class="shrink-0 bg-green-100 p-3 rounded-full">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <div class="flex-1">
          <h3 class="font-bold text-lg text-green-800">ข้อมูลโปรไฟล์ครบถ้วนแล้ว</h3>
          <p class="text-sm text-green-700">คุณสามารถไปสมัครทุนการศึกษาได้แล้ว</p>
        </div>
        <router-link to="/dashboard/apply-scholarship" 
                     class="btn bg-green-600 hover:bg-green-700 text-white border-none">
          ไปสมัครทุน
        </router-link>
      </div>
    </div>
    
    <div v-if="loading" class="flex justify-center items-center h-64">
      <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
    </div>

    <div v-else>
      <StudentProfile 
        v-if="role === 'student'" 
        :initialData="profileData" 
        @refresh="handleRefresh" 
      />

      <AdminProfile 
        v-if="role === 'admin'" 
        :initialData="profileData" 
        @refresh="handleRefresh" 
      />
    </div>
  </div>
</template>