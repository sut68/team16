<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import DocumentDetailModal from './DetailScreening.vue';
import { getAllScreenings, getScreeningById } from '@/services/api/screening'; 
import type { ScreeningResponse } from '@/interfaces/screening';

// คง Interface เดิมไว้
interface DocumentItem {
  id: number;
  title: string;
  applicant: string;
  date: string;
  round: string;
  status: 'pending' | 'approved' | 'rejected' | 'request-change';
  rejection_reason?: string;
  raw_data?: ScreeningResponse;
  created_at_date?: Date; // เพิ่ม field นี้เพื่อช่วยเรื่อง sort (internal use)
}

const isLoading = ref(false);
const errorMsg = ref('');

const allItems = ref<DocumentItem[]>([]);

const activeTab = ref<'pending' | 'history'>('pending');
const searchQuery = ref('');
const sortOrder = ref('newest');
const filterStatus = ref('all');
const filterRound = ref('all');

const isModalOpen = ref(false);
const selectedDocument = ref<any>(null);

const availableRounds = [
  { label: '1/2568', value: '1/2568' },
  { label: '2/2567', value: '2/2567' },
  { label: '3/2567', value: '3/2567' }
];

const mapStatusIdToString = (id: number): DocumentItem['status'] => {
  switch (id) {
    case 1: return 'pending';
    case 2: return 'approved';
    case 3: return 'rejected';
    case 4: return 'request-change';
    default: return 'pending';
  }
};

const formatDate = (dateString: string) => {
  if (!dateString || dateString === '0001-01-01T00:00:00Z') return '-';
  const date = new Date(dateString);
  return date.toLocaleDateString('th-TH', {
    day: '2-digit', month: 'short', year: 'numeric'
  });
};

const fetchData = async () => {
  isLoading.value = true;
  errorMsg.value = '';
  try {
    const response = await getAllScreenings();
    
    // [ปรับ Flow 2] Map ข้อมูลใส่ allItems ครั้งเดียว
    allItems.value = response.map((item: any) => {
      const app = item.Application || item.application;
      const scholarship = item.Scholarship || item.scholarship || (app?.ApplicationScholarships && app.ApplicationScholarships[0]?.Scholarship) || null;
      const title = scholarship?.ScholarshipName || scholarship?.scholarship_name || 'ไม่ระบุชื่อทุน';

      const student = app?.StudentProfile || app?.student_profile || app?.Student || app?.student || null;
      const firstName = student?.first_name_th || student?.FirstNameTH || student?.first_name_en || student?.FirstNameEN || '';
      const lastName = student?.last_name_th || student?.LastNameTH || student?.last_name_en || student?.LastNameEN || '';
      const applicantName = `${firstName} ${lastName}`.trim() || 'ไม่ระบุชื่อ';

      const statusId = item.StatusScreeningID || item.status_screening_id;
      const createdAt = item.CreatedAt || item.created_at;

      return {
        id: item.ID || item.id,
        title,
        applicant: applicantName,
        date: formatDate(createdAt),
        round: '1/2568', // Mock หรือดึงจริงถ้ามี
        status: mapStatusIdToString(statusId),
        rejection_reason: item.RejectionReason || item.rejection_reason,
        raw_data: item,
        created_at_date: new Date(createdAt)
      };
    });

  } catch (err: any) {
    console.error('Error fetching screenings:', err);
    errorMsg.value = 'ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ภายหลัง';
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => { fetchData(); });

// [ปรับ Flow 3] ใช้ Computed แยก Pending/History (เหมือน ApprovalList)
const pendingItems = computed(() => 
  allItems.value.filter(item => item.status === 'pending' || item.status === 'request-change')
);

const historyItems = computed(() => 
  allItems.value.filter(item => item.status === 'approved' || item.status === 'rejected')
);

// [ปรับ Flow 4] Logic การ Filter/Sort แบบเดียวกับ Approval
const filteredItems = computed(() => {
  let result: DocumentItem[] = [];

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
      result.sort((a, b) => (a.created_at_date?.getTime() || 0) - (b.created_at_date?.getTime() || 0));
    } else {
      result.sort((a, b) => (b.created_at_date?.getTime() || 0) - (a.created_at_date?.getTime() || 0));
    }
  } else {
    if (filterStatus.value !== 'all') {
      result = result.filter(item => item.status === filterStatus.value);
    }
    // History เอาใหม่สุดขึ้นก่อนเสมอ
    result.sort((a, b) => (b.created_at_date?.getTime() || 0) - (a.created_at_date?.getTime() || 0));
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

const handleCardClick = async (item: DocumentItem | undefined) => {
  if (!item) return;
  try {
    const fullData = await getScreeningById(item.id);
    selectedDocument.value = fullData; // ส่งข้อมูลเต็มเข้า Modal
    isModalOpen.value = true;
  } catch (error) {
    console.error("Error fetching detail:", error);
    // Fallback ใช้ข้อมูลเดิมที่มีถ้าโหลดไม่ผ่าน
    selectedDocument.value = { data: item.raw_data }; 
    isModalOpen.value = true;
  }
};

const handleExport = () => {
  alert(`กำลังดาวน์โหลดรายงานรายชื่อผู้ขอรับทุน (Round: ${filterRound.value === 'all' ? 'ทั้งหมด' : filterRound.value})`);
};

// [ปรับ Flow 5] Refresh ข้อมูลเมื่อ Modal ทำงานเสร็จ
const handleActionCompleted = () => {
  isModalOpen.value = false;
  fetchData(); 
};
</script>

<template>
<div class="min-h-screen bg-[#f0f2f5] p-6 font-sans text-slate-800">
  <div v-if="isLoading" class="text-center py-20 text-gray-500">
      <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
      <p class="mt-2">กำลังโหลดข้อมูล...</p>
  </div>

  <div v-else>
    <div class="flex flex-col xl:flex-row items-center justify-between gap-4 mb-6 px-1">
      <div class="flex items-center gap-2 overflow-x-auto pb-2 xl:pb-0 hide-scrollbar shrink-0">
        <button
          class="relative btn btn-sm h-11 px-6 rounded-full border-0 text-base font-bold transition-all whitespace-nowrap"
          :class="activeTab === 'pending' ? 'bg-white shadow text-[#1e3a8a]' : 'btn-ghost text-gray-500 hover:bg-white/50'"
          @click="activeTab = 'pending'">
          <div class="flex items-center gap-2">
            <span>รายการที่รอตรวจสอบ</span>
            <span v-if="pendingItems.length > 0" class="badge badge-error badge-sm text-white border-none h-5 px-1.5">{{ pendingItems.length }}</span>
          </div>
          <div v-if="activeTab === 'pending'" class="absolute bottom-1 left-6 right-6 h-[3px] bg-[#1e3a8a] rounded-full"></div>
        </button>

        <button
          class="relative btn btn-sm h-11 px-6 rounded-full border-0 text-base font-bold transition-all whitespace-nowrap"
          :class="activeTab === 'history' ? 'bg-white shadow text-[#1e3a8a]' : 'btn-ghost text-gray-500 hover:bg-white/50'"
          @click="activeTab = 'history'">
          <span>ประวัติการคัดกรอง</span>
          <div v-if="activeTab === 'history'" class="absolute bottom-1 left-6 right-6 h-[3px] bg-[#1e3a8a] rounded-full"></div>
        </button>
      </div>

      <div class="flex flex-col md:flex-row items-center gap-2 w-full xl:w-auto">
        <button v-if="activeTab === 'pending'" @click="handleExport" class="btn btn-sm btn-outline bg-white border-gray-300 text-gray-600 hover:bg-gray-50 hover:border-gray-400 gap-2 h-10 rounded-full font-normal px-5 w-full md:w-auto shadow-sm">
           Export
        </button>

        <select v-model="filterRound" class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-auto font-medium text-gray-600 shadow-sm px-4">
          <option value="all">ทุกรอบการรับสมัคร</option>
          <option v-for="round in availableRounds" :key="round.value" :value="round.value">{{ round.label }}</option>
        </select>

        <select v-if="activeTab === 'pending'" v-model="sortOrder" class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-auto shadow-sm px-4">
          <option value="newest">ใหม่ล่าสุด</option>
          <option value="oldest">ส่งมานานสุด</option>
        </select>

        <select v-if="activeTab === 'history'" v-model="filterStatus" class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-auto shadow-sm px-4">
          <option value="all">สถานะทั้งหมด</option>
          <option value="approved">อนุมัติแล้ว</option>
          <option value="rejected">ปฏิเสธ</option>
        </select>

        <div class="relative w-full md:w-64">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
          </span>
          <input type="text" v-model="searchQuery" placeholder="ค้นหาชื่อทุน, ผู้สมัคร..."
            class="input input-bordered input-sm h-10 w-full rounded-full pl-10 bg-white border-gray-300 shadow-sm focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] text-sm" />
        </div>
      </div>
    </div>

    <div class="space-y-4 pb-10">
      <div v-if="errorMsg" class="alert alert-error text-white text-sm py-2 rounded-lg">{{ errorMsg }}</div>

      <transition-group name="fade" tag="div" class="space-y-4">
        <div v-for="item in filteredItems" :key="item.id" @click="handleCardClick(item)"
          class="card bg-white shadow-sm rounded-2xl cursor-pointer border border-transparent hover:border-blue-200 hover:shadow-md transition-all duration-300 transform hover:-translate-y-1">
          <div class="card-body p-4 md:p-5 flex flex-row items-center justify-between min-h-[6rem]">
            <div class="flex flex-col overflow-hidden pr-2">
              <div class="flex flex-wrap items-center gap-2 mb-1">
                <h3 class="font-bold text-[#1e3a8a] text-base md:text-lg leading-tight truncate">{{ item.title }}</h3>
              </div>
              <div class="flex flex-col md:flex-row md:items-center gap-1 md:gap-3 mt-1 text-sm text-gray-500">
                <span class="inline-flex items-center gap-1 bg-slate-50 px-2 py-0.5 rounded border border-gray-200 text-xs font-semibold text-gray-600 w-fit">
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
                <span class="hidden md:inline">{{ item.status === 'approved' ? 'อนุมัติแล้ว' : item.status === 'rejected' ? 'ปฏิเสธ' : 'แก้ไข' }}</span>
              </div>
              <div class="w-8 h-8 md:w-10 md:h-10 rounded-full flex items-center justify-center text-[#1e3a8a] bg-slate-50 shadow-sm group-hover:bg-blue-50 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-5 h-5 md:w-6 md:h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" /></svg>
              </div>
            </div>
          </div>
        </div>

        <div v-if="!isLoading && filteredItems.length === 0" :key="'empty'" class="flex flex-col items-center justify-center py-16 text-gray-400">
          <div class="bg-white/50 p-4 rounded-full mb-3">
             <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          </div>
          <p v-if="searchQuery">ไม่พบข้อมูลที่ตรงกับ "{{ searchQuery }}"</p>
          <p v-else>ไม่พบรายการในรอบที่เลือก</p>
        </div>
      </transition-group>
    </div>
  </div>

  <DocumentDetailModal
    :isOpen="isModalOpen"
    :documentData="selectedDocument"
    @close="isModalOpen = false"
    @action-completed="handleActionCompleted"
  />
</div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active { transition: opacity 0.2s ease, transform 0.2s ease; }
.fade-enter-from,
.fade-leave-to { opacity: 0; transform: translateY(10px); }
.hide-scrollbar::-webkit-scrollbar { display: none; }
.hide-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
</style>