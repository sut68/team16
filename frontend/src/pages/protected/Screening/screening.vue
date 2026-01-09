<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import ScreeningDetailModal from './DetailScreening.vue';
import { getAllScreenings, getScreeningById } from '@/services/api/screening';
import type { ScreeningResponse } from '@/interfaces/screening';
import type { SemasterResponse } from '@/interfaces';

interface DocumentItem {
  id: number;
  title: string;
  applicant: string;
  submission_date: string;
  round_label: string; // เก็บไว้ใช้กับ logic เดิมถ้าจำเป็น
  roundText: string;   // <--- [เพิ่ม] สำหรับโชว์ข้อความเต็ม (ปี/เทอม/รอบ)
  status: 'pending' | 'approved' | 'rejected' | 'request-change';
  rejection_reason?: string;
  raw_data?: ScreeningResponse;
  created_at_date?: Date;
  semaster: Partial<SemasterResponse>;
}

const isLoading = ref(false);
const errorMsg = ref('');
const allItems = ref<DocumentItem[]>([]);
const activeTab = ref<'pending' | 'history'>('pending');
// Filter States
const searchQuery = ref('');
const sortOrder = ref('newest');
const filterStatus = ref('all');
const filterYear = ref('all');
const filterTerm = ref('all');
const filterRound = ref('all');

// UI States
const isModalOpen = ref(false);
const isFilterOpen = ref(false);
const selectedDocument = ref<any | null>(null);

const availableYears = ref<{ label: string; value: string }[]>([]);
const availableTerms = ref<{ label: string; value: string }[]>([]);
const availableRounds = ref<{ label: string; value: string }[]>([]);

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
  return new Date(dateString).toLocaleDateString('th-TH', {
    day: '2-digit', month: 'short', year: 'numeric'
  });
};

const fetchData = async (background = false) => {
  if (!background) isLoading.value = true;
  errorMsg.value = '';
  try {
    const response = await getAllScreenings();
    
    allItems.value = response.map((item: any) => {
      const appSch = item.application_scholarship || {};
      const scholarship = appSch.scholarship || {};
      const application = appSch.application || {};
      const student = application.student_profile || {};
      const sem = scholarship.semaster || {};

      const termVal = sem.term || sem.Term || '';
      const yearVal = sem.academic_year || sem.AcademicYear || '';
      const roundVal = sem.round || '1';

      const semData = {
        term: termVal,
        academic_year: yearVal,
        semaster_name: sem.semester_name || sem.SemasterName || '',
        round: roundVal 
      };

      const applicantName = student.first_name_th
        ? `${student.first_name_th} ${student.last_name_th}`
        : 'ไม่ระบุชื่อ';

      // --- [แก้ไข] สร้างข้อความแสดงผล ปี/เทอม/รอบ ---
      const roundTextDisplay = (yearVal && termVal) 
        ? `ปี: ${yearVal} เทอม: ${termVal} รอบ: ${roundVal}`
        : `รอบ: ${roundVal}`;

      return {
        id: item.ID,
        title: scholarship.scholarship_name ?? 'ไม่ระบุชื่อทุน',
        applicant: applicantName,
        submission_date: formatDate(item.CreatedAt),
        round_label: roundVal ? `${roundVal}` : 'ไม่ระบุ',
        roundText: roundTextDisplay, // <--- ใส่ค่าตรงนี้
        semaster: semData,
        status: mapStatusIdToString(item.status_screening_id),
        rejection_reason: item.rejection_reason,
        raw_data: item,
        created_at_date: new Date(item.CreatedAt),
      };
    });

    const uniqueYears = [...new Set(allItems.value
      .map(item => item.semaster?.academic_year)
      .filter(y => y && y !== 'ไม่ระบุ')
    )];
    availableYears.value = uniqueYears.sort().reverse().map(y => ({ label: `ปีการศึกษา ${y}`, value: y! }));

  } catch (err) {
    console.error('Error fetching screenings:', err);
    errorMsg.value = 'ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ภายหลัง';
  } finally {
    isLoading.value = false;
  }
};

// ... (ส่วน Watch Filter เหมือนเดิม ไม่ต้องแก้) ...
watch(filterYear, (newYear) => {
  filterTerm.value = 'all';
  filterRound.value = 'all';
  availableTerms.value = [];
  availableRounds.value = [];
  if (newYear !== 'all') {
    const terms = [...new Set(allItems.value
      .filter(i => i.semaster?.academic_year === newYear)
      .map(i => i.semaster?.term)
      .filter(t => t))];
    availableTerms.value = terms.sort().map(t => ({ label: `เทอม ${t}`, value: t! }));
  }
});

watch(filterTerm, (newTerm) => {
  filterRound.value = 'all';
  availableRounds.value = [];
  if (filterYear.value !== 'all' && newTerm !== 'all') {
    const rounds = [...new Set(allItems.value
      .filter(i => i.semaster?.academic_year === filterYear.value && i.semaster?.term === newTerm)
      .map(i => i.semaster?.round || '1'))];
    availableRounds.value = rounds.sort().map(r => ({ label: `รอบ ${r}`, value: String(r) }));
  }
});

watch(activeTab, () => {
  searchQuery.value = '';
  sortOrder.value = 'newest';
  filterStatus.value = 'all';
  filterYear.value = 'all';
  isFilterOpen.value = false; 
});

onMounted(() => { 
  fetchData(); 
});

// ... (Computed Filtered Items เหมือนเดิม) ...
const pendingItems = computed(() =>
  allItems.value.filter(item => item.status === 'pending' || item.status === 'request-change')
);
const historyItems = computed(() =>
  allItems.value.filter(item => item.status === 'approved' || item.status === 'rejected')
);
const filteredItems = computed(() => {
  let result: DocumentItem[] = activeTab.value === 'pending' ? [...pendingItems.value] : [...historyItems.value];
  if (filterYear.value !== 'all') result = result.filter(item => item.semaster?.academic_year === filterYear.value);
  if (filterTerm.value !== 'all') result = result.filter(item => item.semaster?.term === filterTerm.value);
  if (filterRound.value !== 'all') result = result.filter(item => String(item.semaster?.round) === filterRound.value); // cast string เพื่อความชัวร์

  if (activeTab.value === 'history' && filterStatus.value !== 'all') {
    result = result.filter(item => item.status === filterStatus.value);
  }
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(item => 
      item.title.toLowerCase().includes(query) || item.applicant.toLowerCase().includes(query)
    );
  }
  if (sortOrder.value === 'newest') {
    result.sort((a, b) => (b.created_at_date?.getTime() || 0) - (a.created_at_date?.getTime() || 0));
  } else {
    result.sort((a, b) => (a.created_at_date?.getTime() || 0) - (b.created_at_date?.getTime() || 0));
  }
  return result;
});

const stats = computed(() => {
    // Calculate stats from ALL items regardless of tab
    //const pendingCount = allItems.value.filter(i => i.status === 'pending').length;
    const approvedCount = allItems.value.filter(i => i.status === 'approved').length;
    const rejectedCount = allItems.value.filter(i => i.status === 'rejected').length;
    
    // Optionally include 'request-change' in pending count or show separately?
    // Based on user request: "รอจรวจสอบ", "ผ่าน", "ไม่ผ่าน" -> seems to map to pending, approved, rejected.
    // If 'request-change' is considered pending action, maybe group it or just stick to requested 3 columns.
    // Let's stick to the explicit request: "รอการคัดกรอง (pending)", "ผ่าน (approved)", "ไม่ผ่าน (rejected)"
    
    // If 'request-change' is important, we can add it to pending or keep it separate. 
    // Usually 'request-change' is still in 'pending' tab workflow. 
    // For now, let's map: 
    // Col 1: รอการคัดกรอง (Pending + Request Change?) Or just Pending? 
    // User said: "รอการคัดกรอง ผ่านการคัดกรอง ไม่ผ่านการคัดกรอง" (3 col)
    
    const totalPending = allItems.value.filter(i => i.status === 'pending' || i.status === 'request-change').length; // Grouping both as actionable

    return {
        title1: 'รอการคัดกรอง',
        value1: totalPending,
        desc1: 'รายการ',
        
        title2: 'ผ่านการคัดกรอง',
        value2: approvedCount,
        desc2: 'รายการ',
        
        title3: 'ไม่ผ่านการคัดกรอง',
        value3: rejectedCount,
        desc3: 'รายการ'
    };
});



const handleCardClick = async (item: DocumentItem) => {
  if (!item.raw_data && !item.id) return;
  const fullData = await getScreeningById(Number(item.id));
  if (!fullData) {
    selectedDocument.value = item.raw_data ?? null;
    isModalOpen.value = !!selectedDocument.value;
    return;
  }
  const appSch = (fullData?.application_scholarship || item.raw_data?.application_scholarship || {}) as any;
  const application = appSch?.application || {};
  const scholarship = appSch?.scholarship || {};
  const features = scholarship.feature_scholarships || scholarship.FeatureScholarships || [];
  const studentProfile = application?.student_profile || {};
  const familyInfo = studentProfile?.family_info || {};

  selectedDocument.value = {
    ...fullData,
    scholarship: { ...scholarship, feature_scholarships: features },
    application: { ...application, student_profile: { ...studentProfile, family_info: familyInfo } },
    feature_scholarships: features,
    student_profile: studentProfile,
    family_info: familyInfo
  };
  isModalOpen.value = !!selectedDocument.value;
};

const handleActionCompleted = () => {
  isModalOpen.value = false;
  fetchData(); // หรือ fetchData() เช็คชื่อฟังก์ชันให้ตรงกัน
};
</script>

<template>
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible font-sans text-slate-800">
    
    <h1 class="text-2xl font-bold text-slate-800 mb-10">คัดกรองผู้สมัคร</h1>
    
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
            <div class="text-3xl font-bold text-red-600">{{ stats.value3 }}</div>
                <div class="text-xs text-gray-500 mt-1">{{ stats.desc3 }}</div>
            </div>
            <div class="p-3 rounded-full text-red-600 bg-red-50">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"></path></svg>
            </div>
        </div>
        
    </div>

    <div class="flex flex-col xl:flex-row items-end xl:items-center justify-between gap-4 mb-6 border-b border-gray-200">
      <div class="flex gap-8 -mb-[1px] w-full xl:w-auto overflow-x-auto hide-scrollbar">
        <a @click="activeTab = 'pending'" 
           class="pb-3 px-1 text-base font-bold cursor-pointer transition-all duration-200 border-b-[3px] flex items-center gap-2 whitespace-nowrap"
           :class="activeTab === 'pending' ? 'text-[#1e3a8a] border-[#1e3a8a]' : 'text-slate-500 border-transparent hover:text-slate-700 hover:border-slate-300'">
            รายการที่รอตรวจสอบ
            <span v-if="pendingItems.length > 0" class="badge badge-error text-white border-none h-5 px-1.5 text-xs" :class="activeTab === 'pending' ? '' : 'opacity-70'">
              {{ pendingItems.length }}
            </span>
        </a> 
        <a @click="activeTab = 'history'" 
           class="pb-3 px-1 text-base font-bold cursor-pointer transition-all duration-200 border-b-[3px] whitespace-nowrap"
           :class="activeTab === 'history' ? 'text-[#1e3a8a] border-[#1e3a8a]' : 'text-slate-500 border-transparent hover:text-slate-700 hover:border-slate-300'">
            ประวัติการคัดกรอง
        </a>
      </div>

      <div class="flex flex-col md:flex-row items-center gap-2 w-full xl:w-auto pb-4 xl:pb-2">
        <button @click="fetchData()" class="btn btn-sm btn-ghost bg-white border border-gray-300 text-gray-600 hover:bg-gray-50 h-10 w-10 p-0 rounded-full shadow-sm" title="รีเฟรชข้อมูล">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
          </svg>
        </button>

        <select v-if="activeTab === 'pending'" v-model="sortOrder"
          class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-auto shadow-sm px-4">
          <option value="newest">ใหม่ล่าสุด</option>
          <option value="oldest">ส่งมานานสุด</option>
        </select>

        <div v-if="activeTab === 'history'" class="relative w-full md:w-auto">
          <button @click="isFilterOpen = !isFilterOpen" class="btn btn-sm btn-outline bg-white border-gray-300 text-gray-600 hover:bg-gray-50 hover:border-gray-400 gap-2 h-10 rounded-full font-normal px-5 w-full md:w-auto shadow-sm">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M3 3a1 1 0 011-1h12a1 1 0 011 1v3a1 1 0 01-.293.707L12 11.414V15a1 1 0 01-.293.707l-2 2A1 1 0 018 17v-5.586L3.293 6.707A1 1 0 013 6V3z" clip-rule="evenodd" /></svg>
            ตัวกรอง
          </button>
          <div v-if="isFilterOpen" class="fixed inset-0 z-20" @click="isFilterOpen = false"></div>
          <div v-if="isFilterOpen" class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-2xl z-30 border p-4 animate-pop-in">
              <p class="font-bold text-base mb-3 text-slate-700">ตัวกรองประวัติ</p>
              <div class="space-y-3">
                <div>
                    <label class="text-xs font-medium text-gray-500">ปีการศึกษา</label>
                    <select v-model="filterYear" class="select select-bordered select-sm w-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] font-medium text-gray-600 shadow-sm">
                      <option value="all">ทุกปีการศึกษา</option>
                      <option v-for="year in availableYears" :key="year.value" :value="year.value">{{ year.label }}</option>
                    </select>
                </div>
                <div>
                    <label class="text-xs font-medium text-gray-500">เทอม</label>
                    <select v-model="filterTerm" :disabled="filterYear === 'all'" class="select select-bordered select-sm w-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] font-medium text-gray-600 shadow-sm disabled:bg-gray-100">
                      <option value="all">ทุกเทอม</option>
                      <option v-for="term in availableTerms" :key="term.value" :value="term.value">{{ term.label }}</option>
                    </select>
                </div>
                <div>
                    <label class="text-xs font-medium text-gray-500">รอบ</label>
                    <select v-model="filterRound" :disabled="filterTerm === 'all'" class="select select-bordered select-sm w-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] font-medium text-gray-600 shadow-sm disabled:bg-gray-100">
                      <option value="all">ทุกรอบ</option>
                      <option v-for="round in availableRounds" :key="round.value" :value="round.value">{{ round.label }}</option>
                    </select>
                </div>
                <div>
                    <label class="text-xs font-medium text-gray-500">สถานะ</label>
                    <select v-model="filterStatus" class="select select-bordered select-sm w-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] shadow-sm font-medium text-gray-600">
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
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
          </span>
          <input type="text" v-model="searchQuery" placeholder="ค้นหาชื่อทุน, ผู้สมัคร..."
            class="input input-bordered input-sm h-10 w-full rounded-full pl-10 bg-white border-gray-300 shadow-sm focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] text-sm" />
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="text-center py-20 text-gray-500">
      <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
      <p class="mt-2">กำลังโหลดข้อมูล...</p>
    </div>
    <div v-else-if="errorMsg" class="alert alert-error text-white text-sm py-2 rounded-lg mb-4">
      {{ errorMsg }}
    </div>
    <div v-else class="space-y-4 pb-10 overflow-y-auto pr-1 custom-scrollbar flex-1">
      <transition-group name="fade" tag="div" class="space-y-4">
        <div v-for="item in filteredItems" :key="item.id" @click="handleCardClick(item)"
          class="card bg-white border border-gray-250 shadow-sm rounded-2xl cursor-pointer hover:border-blue-200 hover:shadow-md transition-all duration-300 transform hover:-translate-y-1">
          <div class="card-body p-4 md:p-5 flex flex-row items-center justify-between min-h-[6rem]">
            
            <div class="flex flex-col overflow-hidden pr-2">
              <div class="flex flex-wrap items-center gap-2 mb-1">
                <h3 class="font-bold text-[#1e3a8a] text-base md:text-lg leading-tight truncate">
                  {{ item.title }}
                </h3>
                <span v-if="item.status === 'request-change'" class="badge bg-orange-500 text-white badge-xs py-2 px-2 font-normal animate-pulse shadow-sm">
                  รอแก้ไข
                </span>
                <span v-else-if="item.status === 'pending'" class="badge badge-ghost bg-blue-50 text-blue-700 badge-xs py-2 px-2 border-none font-normal">
                  รอตรวจสอบ
                </span>
              </div>
              
              <div class="flex flex-col md:flex-row md:items-center gap-1 md:gap-3 mt-1 text-sm text-gray-500">
                
                <span class="inline-flex items-center gap-1 bg-slate-50 px-2 py-0.5 rounded border border-gray-200 text-xs font-semibold text-gray-600 w-fit">
                   <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3 text-blue-500">
                     <path fill-rule="evenodd" d="M5.75 2a.75.75 0 01.75.75V4h7V2.75a.75.75 0 011.5 0V4h.25A2.75 2.75 0 0118 6.75v8.5A2.75 2.75 0 0115.25 18H4.75A2.75 2.75 0 012 15.25v-8.5A2.75 2.75 0 014.75 4H5V2.75A.75.75 0 015.75 2zm-1 5.5c-.69 0-1.25.56-1.25 1.25v6.5c0 .69.56 1.25 1.25 1.25h10.5c.69 0 1.25-.56 1.25-1.25v-6.5c0-.69-.56-1.25-1.25-1.25H4.75z" clip-rule="evenodd" />
                   </svg>
                   {{ item.roundText }}
                </span>

                <div class="flex items-center gap-2">
                  <span class="font-medium text-gray-600 truncate">{{ item.applicant }}</span>
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
                  'badge-warning text-white': item.status === 'request-change'
                }">
                {{ item.status === 'approved' ? 'อนุมัติแล้ว' : item.status === 'rejected' ? 'ปฏิเสธ' : 'แก้ไข' }}
              </div>
              <div class="w-8 h-8 md:w-10 md:h-10 rounded-full flex items-center justify-center text-[#1e3a8a] bg-slate-50 shadow-sm group-hover:bg-blue-50 transition-colors">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-5 h-5 md:w-6 md:h-6"><path stroke-linecap="round" stroke-linejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" /></svg>
              </div>
            </div>
          </div>
        </div>

        <div v-if="filteredItems.length === 0" :key="'empty'" class="flex flex-col items-center justify-center py-16 text-gray-400">
          <div class="bg-gray-50 p-4 rounded-full mb-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          </div>
          <p v-if="searchQuery">ไม่พบข้อมูลที่ตรงกับ "{{ searchQuery }}"</p>
          <p v-else>ไม่พบรายการ</p>
        </div>
      </transition-group>
    </div>

    <screening-detail-modal
      v-if="isModalOpen && selectedDocument"
      :isOpen="isModalOpen"
      :document-data="selectedDocument"
      @close="isModalOpen = false"
      @action-completed="handleActionCompleted"
    />
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease, transform 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; transform: translateY(10px); }
.hide-scrollbar::-webkit-scrollbar { display: none; }
.hide-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background-color: #cbd5e1; border-radius: 20px; }
@keyframes pop-in { 0% { opacity: 0; transform: scale(0.95); } 100% { opacity: 1; transform: scale(1); } }
.animate-pop-in { animation: pop-in 0.1s ease-out; }
</style>