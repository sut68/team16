<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import DocumentDetailModal from './DocumentDetailModal.vue';
import { getApprovalTasks } from '../../../services/api/approval';
import type { ApprovalTaskResponse } from '@/interfaces';

const allTasks = ref<ApprovalTaskResponse[]>([]);
const isLoading = ref(true);
const error = ref<string | null>(null);

const activeTab = ref<'pending' | 'history'>('pending');
const searchQuery = ref('');
const sortOrder = ref('newest');
const filterStatus = ref('all');
const filterRound = ref('all');

const isModalOpen = ref(false);
const selectedDocument = ref<ApprovalTaskResponse | null>(null);

const availableRounds = ref([
  { label: '1/2568', value: '1/2568' },
  { label: '2/2568', value: '2/2568' },
  { label: '3/2568', value: '3/2568' }
]);

const fetchTasks = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    const tasks = await getApprovalTasks();
    allTasks.value = tasks.map(task => ({
      ...task,
      round: '1/2568', // Mock ค่า
      submission_date: new Date(task.CreatedAt).toLocaleDateString('th-TH', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      }),
    }));
  } catch (err) {
    error.value = 'ไม่สามารถโหลดข้อมูลได้';
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchTasks);
const pendingItems = computed(() => {
  return allTasks.value.filter(item =>
    item.status?.toLowerCase() === 'pending' || item.status?.toLowerCase() === 'request-change'
  );
});

const historyItems = computed(() => {
  return allTasks.value.filter(item =>
    item.status === 'approved' || item.status === 'rejected'
  );
});

// Filter Logic
const filteredItems = computed(() => {
  let result: ApprovalTaskResponse[] = [];

  if (activeTab.value === 'pending') {
    result = [...pendingItems.value];
  } else {
    result = [...historyItems.value];
  }

  if (filterRound.value !== 'all') {
    result = result.filter(item => item.round === filterRound.value);
  }

  if (activeTab.value === 'pending') {
    if (sortOrder.value === 'oldest') {
      result.sort((a, b) => new Date(a.CreatedAt).getTime() - new Date(b.CreatedAt).getTime());
    } else {
      result.sort((a, b) => new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime());
    }
  } else {
    if (filterStatus.value !== 'all') {
      result = result.filter(item => item.status?.toLowerCase() === filterStatus.value.toLowerCase());
    }
  }

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(item => {
      const scholarshipName = item.approval_requirement?.scholarship?.scholarship_name || '';
      const firstName = item.application?.student_profile?.first_name_th || '';
      const lastName = item.application?.student_profile?.last_name_th || '';
      const applicantName = `${firstName} ${lastName}`;

      return scholarshipName.toLowerCase().includes(query) || applicantName.toLowerCase().includes(query);
    });
  }
  return result;
});

watch(activeTab, () => {
  searchQuery.value = '';
  sortOrder.value = 'newest';
  filterStatus.value = 'all';
  filterRound.value = 'all';
});
const handleCardClick = (item: ApprovalTaskResponse) => {
  selectedDocument.value = item;
  isModalOpen.value = true;
};

const handleExport = () => {
  alert(`กำลังดาวน์โหลดรายงานรายชื่อผู้ขอรับทุน (Round: ${filterRound.value === 'all' ? 'ทั้งหมด' : filterRound.value})`);
};

const handleActionCompleted = () => {
  isModalOpen.value = false;
  fetchTasks();
};
</script>

<template>
  <div class="min-h-screen bg-[#f0f2f5] p-6 font-sans text-slate-800">
    <div class="flex flex-col xl:flex-row items-center justify-between gap-4 mb-6 px-1">
      <div class="flex items-center gap-2 overflow-x-auto pb-2 xl:pb-0 hide-scrollbar shrink-0">
        <button
          class="relative btn btn-sm h-11 px-6 rounded-full border-0 text-base font-bold transition-all whitespace-nowrap"
          :class="activeTab === 'pending' ? 'bg-white shadow text-[#1e3a8a]' : 'btn-ghost text-gray-500 hover:bg-white/50'"
          @click="activeTab = 'pending'">
          <div class="flex items-center gap-2">
            <span>เอกสารที่รอการอนุมัติ</span>
            <span v-if="pendingItems.length > 0" class="badge badge-error badge-sm text-white border-none h-5 px-1.5">{{
              pendingItems.length }}</span>
          </div>
          <div v-if="activeTab === 'pending'"
            class="absolute bottom-1 left-6 right-6 h-[3px] bg-[#1e3a8a] rounded-full"></div>
        </button>
        <button
          class="relative btn btn-sm h-11 px-6 rounded-full border-0 text-base font-bold transition-all whitespace-nowrap"
          :class="activeTab === 'history' ? 'bg-white shadow text-[#1e3a8a]' : 'btn-ghost text-gray-500 hover:bg-white/50'"
          @click="activeTab = 'history'">
          <span>ประวัติการอนุมัติ</span>
          <div v-if="activeTab === 'history'"
            class="absolute bottom-1 left-6 right-6 h-[3px] bg-[#1e3a8a] rounded-full"></div>
        </button>
      </div>
      <div class="flex flex-col md:flex-row items-center gap-2 w-full xl:w-auto">
        <button v-if="activeTab === 'pending'" @click="handleExport"
          class="btn btn-sm btn-outline bg-white border-gray-300 text-gray-600 hover:bg-gray-50 hover:border-gray-400 gap-2 h-10 rounded-full font-normal px-5 w-full md:w-auto shadow-sm">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
            stroke="currentColor" class="w-4 h-4">
            <path stroke-linecap="round" stroke-linejoin="round"
              d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5M16.5 12L12 16.5m0 0L7.5 12m4.5 4.5V3" />
          </svg>
          Export
        </button>
        <select v-model="filterRound"
          class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-auto font-medium text-gray-600 shadow-sm px-4">
          <option value="all">ทุกรอบการรับสมัคร</option>
          <option v-for="round in availableRounds" :key="round.value" :value="round.value">{{ round.label }}</option>
        </select>
        <select v-if="activeTab === 'pending'" v-model="sortOrder"
          class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-auto shadow-sm px-4">
          <option value="newest">ใหม่ล่าสุด</option>
          <option value="oldest">ส่งมานานสุด</option>
        </select>
        <select v-if="activeTab === 'history'" v-model="filterStatus"
          class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-auto shadow-sm px-4">
          <option value="all">สถานะทั้งหมด</option>
          <option value="approved">อนุมัติแล้ว</option>
          <option value="rejected">ปฏิเสธ</option>
        </select>
        <div class="relative w-full md:w-64">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </span>
          <input type="text" v-model="searchQuery" placeholder="ค้นหาชื่อทุน, ผู้สมัคร..."
            class="input input-bordered input-sm h-10 w-full rounded-full pl-10 bg-white border-gray-300 shadow-sm focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] text-sm" />
        </div>
      </div>
    </div>
    <div v-if="isLoading" class="text-center py-20 text-gray-500">
      <span class="loading loading-spinner loading-lg"></span>
      <p>Loading data...</p>
    </div>
    <div v-if="error" class="text-center py-20 text-red-500">
      <p>{{ error }}</p>
      <button @click="fetchTasks" class="btn btn-sm btn-outline mt-4">ลองใหม่อีกครั้ง</button>
    </div>
    <div v-if="!isLoading && !error" class="space-y-4 pb-10">
      <transition-group name="fade" tag="div" class="space-y-4">
        <div v-for="item in filteredItems" :key="item.ID" @click="handleCardClick(item)"
          class="card bg-white shadow-sm rounded-2xl cursor-pointer border border-transparent hover:border-blue-200 hover:shadow-md transition-all duration-300 transform hover:-translate-y-1">
          <div class="card-body p-4 md:p-5 flex flex-row items-center justify-between min-h-[6rem]">
            <div class="flex flex-col overflow-hidden pr-2">
              <div class="flex flex-wrap items-center gap-2 mb-1">
                <h3 class="font-bold text-[#1e3a8a] text-base md:text-lg leading-tight truncate">
                  {{ item.approval_requirement?.scholarship?.scholarship_name || 'N/A' }}
                </h3>
                <template v-if="activeTab === 'pending'">
                  <span v-if="item.status?.toLowerCase() === 'request-change'"
                    class="badge bg-orange-500 text-white badge-xs py-2 px-2 font-normal animate-pulse shadow-sm">รอผู้สมัครแก้ไข</span>
                  <span v-if="item.status?.toLowerCase() === 'pending'"
                    class="badge badge-ghost bg-blue-50 text-blue-700 badge-xs py-2 px-2 border-none font-normal">ยื่นใหม่</span>
                </template>
              </div>
              <div class="flex flex-col md:flex-row md:items-center gap-1 md:gap-3 mt-1 text-sm text-gray-500">
                <span
                  class="inline-flex items-center gap-1 bg-slate-50 px-2 py-0.5 rounded border border-gray-200 text-xs font-semibold text-gray-600 w-fit">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor"
                    class="w-3 h-3 text-blue-500">
                    <path fill-rule="evenodd"
                      d="M5.75 2a.75.75 0 01.75.75V4h7V2.75a.75.75 0 011.5 0V4h.25A2.75 2.75 0 0118 6.75v8.5A2.75 2.75 0 0115.25 18H4.75A2.75 2.75 0 012 15.25v-8.5A2.75 2.75 0 014.75 4H5V2.75A.75.75 0 015.75 2zm-1 5.5c-.69 0-1.25.56-1.25 1.25v6.5c0 .69.56 1.25 1.25 1.25h10.5c.69 0 1.25-.56 1.25-1.25v-6.5c0-.69-.56-1.25-1.25-1.25H4.75z"
                      clip-rule="evenodd" />
                  </svg>
                  รอบ: {{ item.round }}
                </span>
                <div class="flex items-center gap-2">
                  <span class="font-medium text-gray-600 truncate">
                    {{ item.application?.student_profile?.first_name_th }} {{
                      item.application?.student_profile?.last_name_th }}
                  </span>
                  <span class="hidden md:inline text-gray-300">|</span>
                  <span class="text-xs md:text-sm text-gray-400">ส่งเมื่อ {{ item.submission_date }}</span>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2 md:gap-3 shrink-0">
              <div v-if="activeTab === 'history'" class="badge badge-sm font-medium h-6 px-2 md:px-3 whitespace-nowrap"
                :class="{
                  'badge-success text-white': item.status === 'approved',
                  'badge-error text-white': item.status === 'rejected',
                }">
                {{ item.status === 'approved' ? 'อนุมัติแล้ว' : 'ปฏิเสธ' }}
              </div>
              <div
                class="w-8 h-8 md:w-10 md:h-10 rounded-full flex items-center justify-center text-[#1e3a8a] bg-slate-50 shadow-sm group-hover:bg-blue-50 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5"
                  stroke="currentColor" class="w-5 h-5 md:w-6 md:h-6">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
                </svg>
              </div>
            </div>
          </div>
        </div>
        <div v-if="filteredItems.length === 0" :key="'empty'"
          class="flex flex-col items-center justify-center py-16 text-gray-400">
          <div class="bg-white/50 p-4 rounded-full mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 opacity-50" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p v-if="searchQuery">ไม่พบข้อมูลที่ตรงกับ "{{ searchQuery }}"</p>
          <p v-else>ไม่พบรายการ</p>
        </div>
      </transition-group>
    </div>
    <DocumentDetailModal v-if="isModalOpen" :isOpen="isModalOpen" :documentData="selectedDocument"
      @close="isModalOpen = false" @action-completed="handleActionCompleted" />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

.hide-scrollbar::-webkit-scrollbar {
  display: none;
}

.hide-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>