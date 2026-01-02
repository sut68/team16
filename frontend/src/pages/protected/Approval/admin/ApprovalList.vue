<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import DocumentDetailModal from './DocumentDetailModal.vue';
import { getApprovalTasks } from '@/services/api/approval';
import type { ApprovalTaskResponse, SemasterResponse } from '@/interfaces';

interface ApprovalTaskDisplay extends ApprovalTaskResponse {
  roundText: string;
  submission_date: string;
  semaster: SemasterResponse;
}

const allTasks = ref<ApprovalTaskDisplay[]>([]);
const isLoading = ref(true);
const error = ref<string | null>(null);

const activeTab = ref<'pending' | 'history'>('pending');
const searchQuery = ref('');
const sortOrder = ref('newest');
const filterStatus = ref('all');

const filterYear = ref('all');
const filterTerm = ref('all');
const filterRound = ref('all');

const availableYears = ref<{ label: string; value: string }[]>([]);
const availableTerms = ref<{ label: string; value: string }[]>([]);
const availableRounds = ref<{ label: string; value: string }[]>([]);

const isModalOpen = ref(false);
const isFilterOpen = ref(false);
const selectedDocument = ref<ApprovalTaskDisplay | null>(null);

const fetchTasks = async () => {
  isLoading.value = true;
  error.value = null;
  try {
    const tasks = await getApprovalTasks();
    const processedTasks = tasks.map((task: ApprovalTaskResponse) => {
      const semaster = task.application_document?.application_scholarship?.application?.semaster;
      const roundText = semaster
        ? `ปี: ${semaster.academic_year} เทอม: ${semaster.term} รอบ: ${semaster.round}`
        : 'N/A';
      return {
        ...task,
        roundText: roundText,
        semaster: semaster,
        submission_date: new Date(task.CreatedAt).toLocaleDateString('th-TH', {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
        }),
      };
    });
    allTasks.value = processedTasks;

    const uniqueYears = [...new Set(processedTasks.map(t => t.semaster?.academic_year).filter(Boolean))];
    availableYears.value = uniqueYears.map(y => ({ label: `ปีการศึกษา ${y}`, value: y }));

  } catch (err) {
    error.value = 'ไม่สามารถโหลดข้อมูลได้';
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};

watch(filterYear, (newYear) => {
  filterTerm.value = 'all';
  filterRound.value = 'all';
  availableTerms.value = [];
  availableRounds.value = [];

  if (newYear !== 'all') {
    const termsForYear = [...new Set(allTasks.value
      .filter(t => t.semaster?.academic_year === newYear)
      .map(t => t.semaster?.term)
      .filter(Boolean))];
    availableTerms.value = termsForYear.map(t => ({ label: `เทอม ${t}`, value: t }));
  }
});

watch(filterTerm, (newTerm) => {
  filterRound.value = 'all';
  availableRounds.value = [];

  if (filterYear.value !== 'all' && newTerm !== 'all') {
    const roundsForTerm = [...new Set(allTasks.value
      .filter(t => t.semaster?.academic_year === filterYear.value && t.semaster?.term === newTerm)
      .map(t => t.semaster?.round)
      .filter(Boolean))];
    availableRounds.value = roundsForTerm.map(r => ({ label: `รอบ ${r}`, value: r }));
  }
});


const checkIsResubmitted = (task: ApprovalTaskDisplay) => {
  if (task.status?.toLowerCase() !== 'pending') return false;
  const decisions = task.approval_decisions;
  if (!decisions || decisions.length === 0) return false;
  const sortedDecisions = [...decisions].sort((a: any, b: any) => (b.ID || 0) - (a.ID || 0));
  const latestDecision = sortedDecisions[0];
  return latestDecision?.decision === 'request-change';
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

const stats = computed(() => {
    if (activeTab.value === 'pending') {
        const purePending = pendingItems.value.filter(t => !checkIsResubmitted(t));
        const resubmitted = pendingItems.value.filter(t => checkIsResubmitted(t));

        return {
            title1: 'เอกสารรอตรวจทั้งหมด',
            value1: pendingItems.value.length,
            desc1: 'รายการ',
            
            title2: 'ยื่นครั้งแรก',
            value2: purePending.length,
            desc2: 'รายการ',
            
            title3: 'ส่งแก้ไข',
            value3: resubmitted.length,
            desc3: 'รายการ'
        };
    } else { // history tab
        const approvedCount = historyItems.value.filter(i => i.status === 'approved').length;
        const rejectedCount = historyItems.value.filter(i => i.status === 'rejected').length;
        const totalHistory = historyItems.value.length;

        return {
            title1: 'เอกสารทั้งหมด',
            value1: totalHistory,
            desc1: 'ในประวัติ',
            
            title2: 'อนุมัติทั้งหมด',
            value2: approvedCount,
            desc2: `จาก ${totalHistory} รายการ`,
            
            title3: 'ปฏิเสธทั้งหมด',
            value3: rejectedCount,
            desc3: `จาก ${totalHistory} รายการ`
        };
    }
});

const filteredItems = computed((): ApprovalTaskDisplay[] => {
  let result: ApprovalTaskDisplay[] = [];

  if (activeTab.value === 'pending') {
    result = [...pendingItems.value];
  } else {
    result = [...historyItems.value];
  }

  if (filterYear.value !== 'all') {
    result = result.filter(item => item.semaster?.academic_year === filterYear.value);
  }
  if (filterTerm.value !== 'all') {
    result = result.filter(item => item.semaster?.term === filterTerm.value);
  }
  if (filterRound.value !== 'all') {
    result = result.filter(item => item.semaster?.round === filterRound.value);
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
      const scholarshipName = item.application_document.application_scholarship?.scholarship?.scholarship_name || '';
      const firstName = item.application_document.application_scholarship?.application?.student_profile?.first_name_th || '';
      const lastName = item.application_document.application_scholarship?.application?.student_profile?.last_name_th || '';
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
  filterYear.value = 'all';
  filterTerm.value = 'all';
  filterRound.value = 'all';
});

const handleCardClick = (item: ApprovalTaskDisplay) => {
  selectedDocument.value = item;
  isModalOpen.value = true;
};

const handleActionCompleted = () => {
  isModalOpen.value = false;
  fetchTasks();
};
</script>

<template>
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible" data-theme="light" data-testid="approval-list-page">
    
    <h1 class="text-2xl font-bold text-slate-800 mb-10">อนุมัติเอกสาร</h1>
    
    <!-- Stats Section -->
    <div class="grid grid-cols-1 md:grid-cols-3 bg-white shadow rounded-2xl border border-gray-100 w-full mb-8 divide-y md:divide-y-0 md:divide-x divide-gray-100">
        
        <div class="p-4 flex flex-row items-center justify-between">
            <div>
                <div class="text-slate-500 text-sm mb-1">{{ stats.title1 }}</div>
                <div class="text-blue-600 text-3xl font-bold">{{ stats.value1 }}</div>
                <div class="text-xs text-gray-500 mt-1">{{ stats.desc1 }}</div>
            </div>
            <div class="text-blue-600 bg-blue-50 p-3 rounded-full">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>
            </div>
        </div>

        <div class="p-4 flex flex-row items-center justify-between">
            <div>
                <div class="text-slate-500 text-sm mb-1">{{ stats.title2 }}</div>
                <div class="text-emerald-700 text-3xl font-bold">{{ stats.value2 }}</div>
                <div class="text-xs text-gray-500 mt-1">{{ stats.desc2 }}</div>
            </div>
            <div class="text-emerald-700 bg-green-50 p-3 rounded-full">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            </div>
        </div>

        <div class="p-4 flex flex-row items-center justify-between">
            <div>
                <div class="text-slate-500 text-sm mb-1">{{ stats.title3 }}</div>
                <div class="text-3xl font-bold" :class="activeTab === 'pending' ? 'text-orange-500' : 'text-red-600'">{{ stats.value3 }}</div>
                <div class="text-xs text-gray-500 mt-1">{{ stats.desc3 }}</div>
            </div>
            <div class="p-3 rounded-full" :class="activeTab === 'pending' ? 'text-orange-500 bg-orange-50' : 'text-red-600 bg-red-50'">
                <svg v-if="activeTab === 'pending'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122"></path></svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path></svg>
            </div>
        </div>
        
    </div>

    <div class="flex flex-col xl:flex-row items-end xl:items-center justify-between gap-4 mb-6 border-b border-gray-200">
      
      <div class="flex gap-8 -mb-[1px] w-full xl:w-auto overflow-x-auto hide-scrollbar">
        <a 
          @click="activeTab = 'pending'" 
          class="pb-3 px-1 text-base font-medium cursor-pointer transition-all duration-200 border-b-[3px] flex items-center gap-2 whitespace-nowrap"
          :class="activeTab === 'pending' 
            ? 'text-[#1e3a8a] border-[#1e3a8a]' 
            : 'text-slate-500 border-transparent hover:text-slate-700 hover:border-slate-300'"
          data-testid="tab-pending"
        >
           เอกสารที่รอการอนุมัติ
           <span v-if="pendingItems.length > 0" 
             class="badge badge-error text-white border-none h-5 px-1.5 text-xs"
             :class="activeTab === 'pending' ? '' : 'opacity-70'">
             {{ pendingItems.length }}
           </span>
        </a> 
        <a 
          @click="activeTab = 'history'" 
          class="pb-3 px-1 text-base font-medium cursor-pointer transition-all duration-200 border-b-[3px] whitespace-nowrap"
          :class="activeTab === 'history' 
            ? 'text-[#1e3a8a] border-[#1e3a8a]' 
            : 'text-slate-500 border-transparent hover:text-slate-700 hover:border-slate-300'"
          data-testid="tab-history"
        >
           ประวัติการอนุมัติ
        </a>
      </div>

      <div class="flex flex-col md:flex-row items-center gap-2 w-full xl:w-auto pb-4 xl:pb-2">
        
        <select v-if="activeTab === 'pending'" v-model="sortOrder"
          class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-auto shadow-sm px-4"
          data-testid="sort-order-select"
        >
          <option value="newest">ใหม่ล่าสุด</option>
          <option value="oldest">ส่งมานานสุด</option>
        </select>

        <div v-if="activeTab === 'history'" class="relative w-full md:w-auto">
          <button @click="isFilterOpen = !isFilterOpen" class="btn btn-sm btn-outline bg-white border-gray-300 text-gray-600 hover:bg-gray-50 hover:border-gray-400 gap-2 h-10 rounded-full font-normal px-5 w-full md:w-auto shadow-sm" data-testid="filter-button">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M3 3a1 1 0 011-1h12a1 1 0 011 1v3a1 1 0 01-.293.707L12 11.414V15a1 1 0 01-.293.707l-2 2A1 1 0 018 17v-5.586L3.293 6.707A1 1 0 013 6V3z" clip-rule="evenodd" />
            </svg>
            ตัวกรอง
          </button>

          <div v-if="isFilterOpen" class="fixed inset-0 z-20" @click="isFilterOpen = false"></div>
          <div v-if="isFilterOpen" class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-2xl z-30 border p-4 animate-pop-in" data-testid="filter-popup">
              <p class="font-bold text-base mb-3 text-slate-700">ตัวกรองประวัติ</p>
              <div class="space-y-3">
                <div>
                    <label class="text-xs font-medium text-gray-500">ปีการศึกษา</label>
                    <select v-model="filterYear"
                      class="select select-bordered select-sm w-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] font-medium text-gray-600 shadow-sm"
                      data-testid="filter-year-select"
                    >
                      <option value="all">ทุกปีการศึกษา</option>
                      <option v-for="year in availableYears" :key="year.value" :value="year.value">{{ year.label }}</option>
                    </select>
                </div>
                <div>
                    <label class="text-xs font-medium text-gray-500">เทอม</label>
                    <select v-model="filterTerm" :disabled="filterYear === 'all'"
                      class="select select-bordered select-sm w-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] font-medium text-gray-600 shadow-sm disabled:bg-gray-100"
                      data-testid="filter-term-select"
                    >
                      <option value="all">ทุกเทอม</option>
                      <option v-for="term in availableTerms" :key="term.value" :value="term.value">{{ term.label }}</option>
                    </select>
                </div>
                <div>
                    <label class="text-xs font-medium text-gray-500">รอบ</label>
                    <select v-model="filterRound" :disabled="filterTerm === 'all'"
                      class="select select-bordered select-sm w-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] font-medium text-gray-600 shadow-sm disabled:bg-gray-100"
                      data-testid="filter-round-select"
                    >
                      <option value="all">ทุกรอบ</option>
                      <option v-for="round in availableRounds" :key="round.value" :value="round.value">{{ round.label }}</option>
                    </select>
                </div>
                <div>
                    <label class="text-xs font-medium text-gray-500">สถานะ</label>
                    <select v-model="filterStatus"
                      class="select select-bordered select-sm w-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] shadow-sm font-medium text-gray-600"
                      data-testid="filter-status-select"
                    >
                      <option value="all">สถานะทั้งหมด</option>
                      <option value="approved">อนุมัติแล้ว</option>
                      <option value="rejected">ปฏิเสธ</option>
                    </select>
                </div>
              </div>
          </div>
        </div>
        
        <div class="relative w-full md:w-64">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </span>
          <input type="text" v-model="searchQuery" placeholder="ค้นหาชื่อทุน, ผู้สมัคร..."
            class="input input-bordered input-sm h-10 w-full rounded-full pl-10 bg-white border-gray-300 shadow-sm focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] text-sm"
            data-testid="search-input"
          />
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="text-center py-20 text-gray-500" data-testid="loading-indicator">
      <span class="loading loading-spinner loading-lg"></span>
      <p>Loading data...</p>
    </div>
    <div v-if="error" class="text-center py-20 text-red-500" data-testid="error-message">
      <p>{{ error }}</p>
      <button @click="fetchTasks" class="btn btn-sm btn-outline mt-4">ลองใหม่อีกครั้ง</button>
    </div>
    
    <div v-if="!isLoading && !error" class="space-y-4 pb-10 overflow-y-auto pr-1 custom-scrollbar flex-1" data-testid="approval-list-container">
      <transition-group name="fade" tag="div" class="space-y-4">
        <div v-for="item in filteredItems" :key="item.ID" @click="handleCardClick(item)"
          class="card bg-white border border-gray-250 shadow-sm rounded-2xl cursor-pointer hover:border-blue-200 hover:shadow-md transition-all duration-300 transform hover:-translate-y-1"
          :data-testid="`approval-card-${item.ID}`"
        >
          <div class="card-body p-4 md:p-5 flex flex-row items-center justify-between min-h-[6rem]">
            <div class="flex flex-col overflow-hidden pr-2">
              <div class="flex flex-wrap items-center gap-2 mb-1">
                <h3 class="font-bold text-[#1e3a8a] text-base md:text-lg leading-tight truncate">
                  {{ item.application_document.application_scholarship?.scholarship?.scholarship_name || 'N/A' }}
                </h3>
                <template v-if="activeTab === 'pending'">
                  <span v-if="item.status?.toLowerCase() === 'request-change'"
                    class="badge bg-orange-500 text-white badge-xs py-2 px-2 font-normal animate-pulse shadow-sm">
                    รอผู้สมัครแก้ไข
                  </span>
                  <span v-else-if="checkIsResubmitted(item)"
                    class="badge badge-info text-white badge-xs py-2 px-2 font-normal animate-pulse shadow-sm">
                    มีการส่งแก้ไขใหม่
                  </span>
                  <span v-else-if="item.status?.toLowerCase() === 'pending'"
                    class="badge badge-ghost bg-blue-50 text-blue-700 badge-xs py-2 px-2 border-none font-normal">
                    ยื่นใหม่
                  </span>
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
                  {{ item.roundText }}
                </span>
                <div class="flex items-center gap-2">
                  <span class="font-medium text-gray-600 truncate">
                    {{ item.application_document.application_scholarship?.application?.student_profile?.first_name_th }}
                    {{ item.application_document.application_scholarship?.application?.student_profile?.last_name_th }}
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
          class="flex flex-col items-center justify-center py-16 text-gray-400"
          data-testid="empty-list-message"
        >
          <div class="bg-gray-50 p-4 rounded-full mb-3">
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
    <DocumentDetailModal v-show="isModalOpen" :isOpen="isModalOpen" :documentData="selectedDocument"
      @close="isModalOpen = false" @action-completed="handleActionCompleted" />
  </div>
</template>

<style scoped>
/* Style เดิม */
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

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 20px;
}
</style>