<script setup lang="ts">
import { ref, onMounted } from 'vue';
import * as UserAPI from '@/services/api/user';
import StudentProfile from '@/components/ui/ProfileStudent.vue';
import AdminProfile from '@/components/ui/ProfileAdmin.vue';
import Swal from 'sweetalert2';

const role = ref('');
const loading = ref(true);
const profileData = ref<any>(null);

const loadProfile = async () => {
  loading.value = true;
  try {
    const res = await UserAPI.getMyProfile();
    if (res) {
      role.value = res.role;
      // รวม data และ family_info ให้ Component ลูก
      profileData.value = { 
        ...res.data, 
        family_info: res.family || (res.data as any).family_info || {} 
      };
    }
  } catch (err) {
    console.error(err);
    Swal.fire('ข้อผิดพลาด', 'ไม่สามารถโหลดข้อมูลโปรไฟล์ได้', 'error');
  } finally {
    loading.value = false;
  }
};

onMounted(loadProfile);
</script>

<template>
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-y-auto" data-theme="light">
    
    <div v-if="loading" class="flex justify-center items-center h-64">
      <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
    </div>

    <div v-else>
      <StudentProfile 
        v-if="role === 'student'" 
        :initialData="profileData" 
        @refresh="loadProfile" 
      />

      <AdminProfile 
        v-if="role === 'admin'" 
        :initialData="profileData" 
        @refresh="loadProfile" 
      />
    </div>
  </div>
</template>