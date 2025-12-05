<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import DocumentDetailModal from './DocumentDetailModal.vue';

// 1. Interface
interface DocumentItem {
  id: number;
  title: string;
  applicant: string;
  date: string;
  round: string;
  subStatus?: 'new' | 'resubmitted';
  status: 'pending' | 'approved' | 'rejected' | 'request-change';
}

// 2. Mock Data
const pendingItems = ref<DocumentItem[]>([
  {
    id: 4,
    title: 'ทุนมูลนิธิ SCG (เพื่อการศึกษา)',
    applicant: 'นายสมชาย ใจดี',
    date: '03 ธ.ค. 2568',
    round: '1/2568',
    subStatus: 'new',
    status: 'pending'
  },
  {
    id: 3,
    title: 'ทุนกู้ยืม กยศ. (รายเก่า)',
    applicant: 'น.ส.มานี มีตา',
    date: '02 ธ.ค. 2568',
    round: '1/2568',
    subStatus: 'resubmitted',
    status: 'pending'
  },
  {
    id: 2,
    title: 'ทุนบริษัท ปตท. จำกัด (มหาชน)',
    applicant: 'นายปิติ รักเรียน',
    date: '01 ธ.ค. 2568',
    round: 'Special',
    subStatus: 'new',
    status: 'pending'
  },
  {
    id: 1,
    title: 'ทุนศิษย์เก่าสัมพันธ์',
    applicant: 'น.ส.ชูใจ ใส่ใจ',
    date: '30 พ.ย. 2568',
    round: '2/2567',
    subStatus: 'new',
    status: 'pending'
  },
]);

const historyItems = ref<DocumentItem[]>([
  { id: 6, title: 'ทุนมูลนิธิ SCG (เพื่อการศึกษา)', applicant: 'นายกล้าหาญ ชาญชัย', date: '20 พ.ย. 2568', round: '1/2568', status: 'approved' },
  { id: 5, title: 'ทุนการศึกษา ธนาคารกรุงเทพ', applicant: 'น.ส.แก้วตา ขวัญใจ', date: '15 พ.ย. 2568', round: '1/2568', status: 'rejected' },
]);

const availableRounds = [
  { label: 'รอบที่ 1/2568', value: '1/2568' },
  { label: 'รอบที่ 2/2567 (ตกค้าง)', value: '2/2567' },
  { label: 'รอบพิเศษ', value: 'Special' }
];

// 3. State
type TabType = 'pending' | 'history';
const activeTab = ref<TabType>('pending');
const searchQuery = ref('');
const sortOrder = ref('newest');
const filterStatus = ref('all');
const filterRound = ref('all');

const isModalOpen = ref(false);
const selectedDocument = ref<DocumentItem | null>(null);

// 4. Computed Logic
const filteredItems = computed(() => {
  let result: DocumentItem[] = [];

  if (activeTab.value === 'pending') result = [...pendingItems.value];
  else result = [...historyItems.value];

  if (filterRound.value !== 'all') {
    result = result.filter(item => item.round === filterRound.value);
  }

  if (activeTab.value === 'pending') {
    if (sortOrder.value === 'oldest') result.sort((a, b) => a.id - b.id);
    else result.sort((a, b) => b.id - a.id);
  } else {
    if (filterStatus.value !== 'all') result = result.filter(item => item.status === filterStatus.value);
  }

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(item =>
      item.title.toLowerCase().includes(query) ||
      item.applicant.toLowerCase().includes(query)
    );
  }
  return result;
});

watch(activeTab, () => {
  searchQuery.value = '';
  sortOrder.value = 'newest';
  filterStatus.value = 'all';
  filterRound.value = 'all';
});

const handleCardClick = (item: DocumentItem | undefined) => {
  if (!item) return;
  selectedDocument.value = item;
  isModalOpen.value = true;
};

const handleExport = () => {
  alert(`กำลังดาวน์โหลดรายงานรายชื่อผู้ขอรับทุน (Round: ${filterRound.value === 'all' ? 'ทั้งหมด' : filterRound.value})`);
};

const handleAction = (payload: any) => {
  if (payload.action === 'approve' || payload.action === 'reject' || payload.action === 'request-change') {
    const index = pendingItems.value.findIndex(i => i.id === payload.id);
    if (index !== -1) {
      const item = pendingItems.value[index];
      if (!item) return;

      item.status = payload.action === 'approve' ? 'approved'
        : payload.action === 'reject' ? 'rejected'
          : 'request-change';

      delete item.subStatus;
      const removedItem = pendingItems.value.splice(index, 1)[0];
      if (removedItem) {
        historyItems.value.unshift(removedItem);
      }
    }
  }
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

    <div class="space-y-4 pb-10">
      <transition-group name="fade" tag="div" class="space-y-4">

        <div v-for="item in filteredItems" :key="item.id" @click="handleCardClick(item)"
          class="card bg-white shadow-sm rounded-2xl cursor-pointer border border-transparent hover:border-blue-200 hover:shadow-md transition-all duration-300 transform hover:-translate-y-1">
          <div class="card-body p-4 md:p-5 flex flex-row items-center justify-between min-h-[6rem]">

            <div class="flex flex-col overflow-hidden pr-2">
              <div class="flex flex-wrap items-center gap-2 mb-1">
                <h3 class="font-bold text-[#1e3a8a] text-base md:text-lg leading-tight truncate">{{ item.title }}</h3>

                <template v-if="activeTab === 'pending'">
                  <span v-if="item.subStatus === 'resubmitted'"
                    class="badge bg-orange-500 text-white badge-xs py-2 px-2 font-normal animate-pulse shadow-sm">ส่งเอกสารแก้</span>
                  <span v-if="item.subStatus === 'new'"
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
                  <span class="font-medium text-gray-600 truncate">{{ item.applicant }}</span>
                  <span class="hidden md:inline text-gray-300">|</span>
                  <span class="text-xs md:text-sm text-gray-400">ส่งเมื่อ {{ item.date }}</span>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-2 md:gap-3 shrink-0">
              <div v-if="activeTab === 'history'" class="badge badge-sm font-medium h-6 px-2 md:px-3 whitespace-nowrap"
                :class="{
                  'badge-success text-white': item.status === 'approved',
                  'badge-error text-white': item.status === 'rejected',
                  'badge-warning text-white': item.status === 'request-change'
                }">
                <span class="hidden md:inline">{{ item.status === 'approved' ? 'อนุมัติแล้ว' : item.status ===
                  'rejected' ? 'ปฏิเสธ'
                  : 'แก้ไข' }}</span>
                <span class="md:hidden">{{ item.status === 'approved' ? 'อนุมัติ' : item.status === 'rejected' ?
                  'ปฏิเสธ' : 'แก้ไข'
                }}</span>
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
          <p v-else>ไม่พบรายการในรอบที่เลือก</p>
        </div>

      </transition-group>
    </div>

    <DocumentDetailModal :isOpen="isModalOpen" :documentData="selectedDocument" @close="isModalOpen = false"
      @approve="(payload: any) => handleAction({ ...payload, action: 'approve' })"
      @reject="(payload: any) => handleAction({ ...payload, action: 'reject' })"
      @request-change="(payload: any) => handleAction({ ...payload, action: 'request-change' })" />

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