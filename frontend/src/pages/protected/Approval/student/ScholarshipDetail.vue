<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { getScholarships, applyForScholarship } from '@/services/api/scholarship';
import type { ScholarshipResponse } from '@/interfaces';

const scholarships = ref<ScholarshipResponse[]>([]);
const isLoading = ref(true);
const error = ref<string | null>(null);
const router = useRouter();

const fetchScholarships = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    scholarships.value = await getScholarships();
    console.log('Data received for scholarships:', scholarships.value);
  } catch (err) {
    error.value = 'ไม่สามารถโหลดข้อมูลทุนการศึกษาได้';
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

const handleApply = async (scholarshipId: number) => {
  const studentProfileId = 1;

  try {
    const result = await applyForScholarship(scholarshipId, studentProfileId);
    console.log('Application result:', result);
    router.push('/dashboard/track-status');

  } catch (err) {
    alert('เกิดข้อผิดพลาดในการสมัครทุน หรือคุณอาจเคยสมัครทุนนี้ไปแล้ว');
    console.error(err);
  }
};

onMounted(fetchScholarships);
</script>

<template>
  <div class="min-h-screen bg-[#f0f2f5] p-6 font-sans text-slate-800">
    <div class="max-w-4xl mx-auto">
      <h1 class="text-2xl font-bold text-[#1e3a8a] mb-6 flex items-center gap-3">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
          stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
        </svg>
        ยื่นสมัครทุนการศึกษา
      </h1>

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
          class="card bg-white shadow-sm hover:shadow-lg transition-all border border-transparent rounded-2xl overflow-hidden group">

          <div class="card-body p-5">
            <div class="flex flex-col md:flex-row gap-4">
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  <span class="bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-1 rounded-full">{{
                    scholarship.typescholarship?.type_name || 'ทั่วไป' }}</span>
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
                <button @click="handleApply(scholarship.ID)"
                  class="btn btn-primary btn-sm bg-[#1e3a8a] hover:bg-blue-800 border-none text-white rounded-full px-5">
                  สมัครทุนนี้
                </button>
              </div>
            </div>
          </div>
        </div>
        <div v-if="scholarships.length === 0" class="text-center py-20 text-gray-500">
          <p>ไม่พบทุนการศึกษาที่เปิดรับสมัครในขณะนี้</p>
        </div>
      </div>
    </div>
  </div>
</template>
