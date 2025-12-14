<script setup lang="ts">
import { ref, computed, reactive, watch } from 'vue';

// --- 1. Interfaces & Types ---
interface InterviewRound {
  ID: number;
  name: string;
  description: string;
  scholarship_name: string;
  scholarship_type: 'merit' | 'financial' | 'activity';
  date: string;
  start_time: string;
  end_time: string;
  slot_duration: number;
  mode: 'Online' | 'Onsite';
  location_name?: string;
  meeting_link?: string;
  interviewers: string[];
  total_slots: number;
  booked_slots: number;
  status: 'Open' | 'Closed' | 'Draft' | 'Full';
}

interface SlotPreview {
  time: string;
  isBooked: boolean;
  studentName?: string;
}

// --- 2. State Management ---
const isModalOpen = ref(false);
const isFilterOpen = ref(false);
const modalMode = ref<'create' | 'view'>('create');
const activeTab = ref<'active' | 'history'>('active');

// Filter & Sort State
const searchQuery = ref('');
const filterStatus = ref('All');
const filterScholarship = ref('All');
const sortBy = ref('DateAsc');

// --- 3. Mock Data ---
const scholarships = [
  { id: 1, name: 'ทุนเรียนดี', type: 'merit' },
  { id: 2, name: 'ทุนขาดแคลนทุนทรัพย์', type: 'financial' },
  { id: 3, name: 'ทุนกิจกรรมดีเด่น', type: 'activity' }
];

const locations = [
  { id: 1, name: 'อาคารเรียนรวม 1 - ห้อง 101' },
  { id: 2, name: 'อาคารวิศวกรรม - ห้องประชุม 2' }
];

const interviewersList = [
  { id: 1, name: 'ดร.สมชาย ใจดี' },
  { id: 2, name: 'ผศ.ดร.วิจัย เก่งมาก' },
  { id: 3, name: 'อ.นิสิต รักเรียน' }
];

const rounds = ref<InterviewRound[]>([
  {
    ID: 1,
    name: 'สัมภาษณ์ทุนเรียนดี รุ่นที่ 1/2567',
    scholarship_name: 'ทุนเรียนดี',
    scholarship_type: 'merit',
    description: 'เตรียม Portfolio มาด้วย',
    date: '2024-12-25',
    start_time: '09:00',
    end_time: '12:00',
    slot_duration: 30,
    mode: 'Onsite',
    location_name: 'อาคารเรียนรวม 2 ห้อง 201',
    interviewers: ['ดร.สมชาย', 'อ.ใจดี'],
    total_slots: 6,
    booked_slots: 6,
    status: 'Full'
  },
  {
    ID: 2,
    name: 'สัมภาษณ์ทุนขาดแคลน (รอบเก็บตก)',
    scholarship_name: 'ทุนขาดแคลนทุนทรัพย์',
    scholarship_type: 'financial',
    description: 'สัมภาษณ์ผ่าน Zoom',
    date: '2024-12-28',
    start_time: '13:00',
    end_time: '16:00',
    slot_duration: 20,
    mode: 'Online',
    meeting_link: 'https://zoom.us/j/123456',
    interviewers: ['อ.วิจัย'],
    total_slots: 9,
    booked_slots: 2,
    status: 'Open'
  },
  {
    ID: 3,
    name: 'สัมภาษณ์ทุนกิจกรรม (เก่า)',
    scholarship_name: 'ทุนกิจกรรมดีเด่น',
    scholarship_type: 'activity',
    description: 'รอบที่แล้ว',
    date: '2024-10-15',
    start_time: '09:00',
    end_time: '12:00',
    slot_duration: 15,
    mode: 'Onsite',
    location_name: 'ห้องประชุมวิศวะ',
    interviewers: ['อ.นิสิต'],
    total_slots: 10,
    booked_slots: 10,
    status: 'Closed'
  }
]);

// Form Data
const formData = reactive({
  name: '', scholarship_id: '', description: '', date: '', 
  start_time: '09:00', end_time: '16:00', slot_duration: 30, 
  mode: 'Onsite' as 'Onsite' | 'Online', location_id: '', meeting_link: '', interviewer_ids: [] as string[]
});

// --- 4. Computed Logic ---

const stats = computed(() => {
  const activeRounds = rounds.value.filter(r => r.status !== 'Closed');
  const totalRounds = activeRounds.length;
  const totalCapacity = activeRounds.reduce((acc, r) => acc + r.total_slots, 0);
  const totalBooked = activeRounds.reduce((acc, r) => acc + r.booked_slots, 0);
  const bookedPercent = totalCapacity > 0 ? Math.round((totalBooked / totalCapacity) * 100) : 0;
  return { totalRounds, totalCapacity, totalBooked, bookedPercent };
});

const filteredRounds = computed(() => {
  let result = rounds.value;

  // 1. Tab Logic
  if (activeTab.value === 'active') {
    result = result.filter(r => r.status !== 'Closed');
  } else {
    result = result.filter(r => r.status === 'Closed');
  }

  // 2. Search
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(r => 
      r.name.toLowerCase().includes(q) || 
      r.scholarship_name.toLowerCase().includes(q)
    );
  }

  // 3. Filters
  if (filterStatus.value !== 'All') {
    result = result.filter(r => r.status === filterStatus.value);
  }

  if (filterScholarship.value !== 'All') {
    result = result.filter(r => r.scholarship_type === filterScholarship.value.toLowerCase());
  }

  // 4. Sort
  result = [...result].sort((a, b) => {
    if (sortBy.value === 'DateAsc') return new Date(a.date).getTime() - new Date(b.date).getTime();
    if (sortBy.value === 'DateDesc') return new Date(b.date).getTime() - new Date(a.date).getTime();
    if (sortBy.value === 'BookedDesc') return b.booked_slots - a.booked_slots;
    return 0;
  });

  return result;
});

const previewSlots = computed<SlotPreview[]>(() => {
  const slots: SlotPreview[] = [];
  if (!formData.start_time || !formData.end_time || !formData.slot_duration) return slots;

  const [startH = 0, startM = 0] = formData.start_time.split(':').map(Number);
  const [endH = 0, endM = 0] = formData.end_time.split(':').map(Number);
  
  let currentMin = startH * 60 + startM;
  const endMin = endH * 60 + endM;

  if (formData.slot_duration <= 0 || currentMin >= endMin) return slots;

  while (currentMin + formData.slot_duration <= endMin) {
    const h = Math.floor(currentMin / 60);
    const m = currentMin % 60;
    const timeStr = `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`;
    const isBookedMock = modalMode.value === 'view' ? Math.random() > 0.6 : false;
    
    slots.push({
      time: timeStr,
      isBooked: isBookedMock,
      studentName: isBookedMock ? 'นายทดสอบ ระบบ' : undefined
    });

    currentMin += formData.slot_duration;
  }
  return slots;
});

const getInitials = (name: string) => name.charAt(0);

// --- 5. Watchers ---
watch(activeTab, () => {
  searchQuery.value = '';
  filterStatus.value = 'All';
  filterScholarship.value = 'All';
  isFilterOpen.value = false;
  sortBy.value = activeTab.value === 'active' ? 'DateAsc' : 'DateDesc';
});

// --- 6. Methods ---
const openCreateModal = () => {
  modalMode.value = 'create';
  Object.assign(formData, {
    name: '', scholarship_id: '', description: '', date: '', 
    start_time: '09:00', end_time: '12:00', slot_duration: 30, 
    mode: 'Onsite', location_id: '', meeting_link: '', interviewer_ids: []
  });
  isModalOpen.value = true;
};

const openViewModal = (round: InterviewRound) => {
  modalMode.value = 'view';
  formData.name = round.name;
  formData.date = round.date;
  formData.start_time = round.start_time;
  formData.end_time = round.end_time;
  formData.slot_duration = round.slot_duration;
  formData.mode = round.mode;
  formData.description = round.description;
  isModalOpen.value = true;
};

const handleDelete = (id: number) => {
  if(confirm('ต้องการลบรอบสัมภาษณ์นี้ใช่หรือไม่?')) {
    rounds.value = rounds.value.filter(r => r.ID !== id);
  }
};

const saveRound = () => {
  alert('บันทึกข้อมูลเรียบร้อย (Mock)');
  isModalOpen.value = false;
};
</script>

<template>
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible" data-theme="light">
    
    <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4 mb-8">
      <div>
      </div>
      <button 
        @click="openCreateModal" 
        class="btn btn-sm bg-white border border-gray-300 text-gray-700 hover:bg-gray-100 hover:border-gray-400 flex items-center justify-center gap-2 rounded-full px-5 h-10 w-full md:w-auto shadow-sm transition-all duration-150"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd" />
        </svg>
        สร้างรอบสัมภาษณ์
      </button>
    </div>

    <div class="stats stats-vertical lg:stats-horizontal shadow bg-white border border-gray-100 w-full mb-8 relative z-0">
      <div class="stat">
        <div class="stat-figure text-[#1e3a8a] bg-blue-50 p-3 rounded-full">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
        </div>
        <div class="stat-title text-slate-500">รอบที่เปิดอยู่</div>
        <div class="stat-value text-[#1e3a8a]">{{ stats.totalRounds }}</div>
        <div class="stat-desc">กำลังดำเนินการ</div>
      </div>
      
      <div class="stat">
        <div class="stat-figure text-secondary bg-pink-50 p-3 rounded-full">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
        </div>
        <div class="stat-title text-slate-500">นศ. ที่จองแล้ว</div>
        <div class="stat-value text-secondary">{{ stats.totalBooked }}</div>
        <div class="stat-desc text-success font-medium">↗︎ {{ stats.bookedPercent }}% ของที่นั่งทั้งหมด</div>
      </div>

      <div class="stat">
        <div class="stat-figure text-orange-500 bg-orange-50 p-3 rounded-full">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block w-8 h-8 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
        </div>
        <div class="stat-title text-slate-500">ที่นั่งว่างเหลือ</div>
        <div class="stat-value text-orange-500">{{ stats.totalCapacity - stats.totalBooked }}</div>
        <div class="stat-desc">ที่นั่ง</div>
      </div>
    </div>

    <div class="relative z-50 flex flex-col xl:flex-row items-end xl:items-center justify-between gap-4 mb-6 border-b border-gray-200">
      
      <div class="flex gap-8 -mb-[1px]">
        <a 
          @click="activeTab = 'active'" 
          class="pb-3 px-1 text-base font-medium cursor-pointer transition-all duration-200 border-b-[3px]"
          :class="activeTab === 'active' 
            ? 'text-[#1e3a8a] border-[#1e3a8a]' 
            : 'text-slate-500 border-transparent hover:text-slate-700 hover:border-slate-300'"
        >
           รอบปัจจุบัน
        </a> 
        <a 
          @click="activeTab = 'history'" 
          class="pb-3 px-1 text-base font-medium cursor-pointer transition-all duration-200 border-b-[3px]"
          :class="activeTab === 'history' 
            ? 'text-[#1e3a8a] border-[#1e3a8a]' 
            : 'text-slate-500 border-transparent hover:text-slate-700 hover:border-slate-300'"
        >
           ประวัติย้อนหลัง
        </a>
      </div>

      <div class="flex flex-col md:flex-row items-center gap-2 w-full xl:w-auto pb-4 xl:pb-2">
         
        <div class="relative w-full md:w-auto">
            <button @click="isFilterOpen = !isFilterOpen" class="btn btn-sm btn-outline bg-white border-gray-300 text-gray-600 hover:bg-gray-50 hover:border-gray-400 gap-2 h-10 rounded-full font-normal px-5 w-full md:w-auto shadow-sm">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M3 3a1 1 0 011-1h12a1 1 0 011 1v3a1 1 0 01-.293.707L12 11.414V15a1 1 0 01-.293.707l-2 2A1 1 0 018 17v-5.586L3.293 6.707A1 1 0 013 6V3z" clip-rule="evenodd" />
                </svg>
                ตัวกรอง
            </button>

            <div v-if="isFilterOpen" class="fixed inset-0 z-50" @click="isFilterOpen = false"></div>
            <div v-if="isFilterOpen" class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-2xl z-[100] border p-4 animate-pop-in">
                
                <div v-if="activeTab === 'active'">
                  <p class="font-bold text-base mb-3 text-slate-700">กรองรอบปัจจุบัน</p>
                  <div class="space-y-3">
                      <div>
                          <label class="text-xs font-medium text-gray-500">สถานะ</label>
                          <select v-model="filterStatus" class="select select-bordered select-sm w-[98%] h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] shadow-sm font-medium text-gray-600">
                              <option value="All">ทั้งหมด</option>
                              <option value="Open">เปิดให้จอง</option>
                              <option value="Full">เต็มแล้ว</option>
                              <option value="Closed">ปิดแล้ว</option>
                          </select>
                      </div>
                      <div>
                          <label class="text-xs font-medium text-gray-500">ประเภททุน</label>
                          <select v-model="filterScholarship" class="select select-bordered select-sm w-[98%] h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] shadow-sm font-medium text-gray-600">
                              <option value="All">ทั้งหมด</option>
                              <option value="Merit">ทุนเรียนดี</option>
                              <option value="Financial">ทุนขาดแคลน</option>
                              <option value="Activity">ทุนกิจกรรม</option>
                          </select>
                      </div>
                  </div>
                </div>

                <div v-else>
                  <p class="font-bold text-base mb-3 text-slate-700">กรองประวัติย้อนหลัง</p>
                  <div class="space-y-3">
                      <div>
                          <label class="text-xs font-medium text-gray-500">ประเภททุน</label>
                          <select v-model="filterScholarship" class="select select-bordered select-sm w-[98%] h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] shadow-sm font-medium text-gray-600">
                              <option value="All">ทั้งหมด</option>
                              <option value="Merit">ทุนเรียนดี</option>
                              <option value="Financial">ทุนขาดแคลน</option>
                              <option value="Activity">ทุนกิจกรรม</option>
                          </select>
                      </div>
                  </div>
                </div>

            </div>
        </div>

         <select v-model="sortBy" class="select select-bordered select-sm rounded-full h-10 bg-white text-sm border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a] w-full md:w-40 shadow-sm px-6">
            <option value="DateAsc">วันที่ (เก่าสุดก่อน)</option>
            <option value="DateDesc">วันที่ (ใหม่สุดก่อน)</option>
            <option value="BookedDesc">ยอดจองมากที่สุด</option>
         </select>

         <div class="relative w-full md:w-64">
            <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
               <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
               </svg>
            </span>
            <input type="text" v-model="searchQuery" placeholder="ค้นหาชื่อรอบ, ทุน..." 
               class="input input-bordered input-sm h-10 w-full rounded-full pl-10 bg-white border-gray-300 shadow-sm focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] text-sm" />
         </div>

      </div>
    </div>

    <div class="bg-slate-50/30 rounded-3xl shadow-sm border border-slate-100 p-6 min-h-[600px] flex-1 relative z-0">
      
      <div v-if="activeTab === 'active'" class="animate-fade-in">
        <div v-if="filteredRounds.length > 0" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
          <div v-for="round in filteredRounds" :key="round.ID" 
            class="group card bg-white border border-gray-200 hover:border-blue-400 hover:shadow-xl transition-all duration-300 rounded-2xl overflow-visible">
            
            <div class="px-5 py-4 border-b border-gray-50 bg-slate-50/50 flex justify-between items-start rounded-t-2xl relative">
              <div>
                <div class="flex gap-2 mb-2">
                  <div class="badge badge-sm font-medium border-none text-white shadow-sm" 
                    :class="{
                      'badge-success': round.status === 'Open',
                      'badge-error': round.status === 'Full',
                      'badge-ghost text-gray-500': round.status === 'Closed'
                    }">
                    {{ round.status === 'Open' ? 'เปิดให้จอง' : (round.status === 'Full' ? 'เต็มแล้ว' : 'ปิด') }}
                  </div>
                   <div class="badge badge-sm badge-outline text-gray-500">{{ round.scholarship_name }}</div>
                </div>
                <h3 class="font-bold text-[#1e3a8a] text-lg leading-tight group-hover:text-blue-600 transition-colors">{{ round.name }}</h3>
              </div>

              <div class="dropdown dropdown-end">
                <label tabindex="0" class="btn btn-circle btn-ghost btn-sm text-gray-400 hover:bg-white">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="w-5 h-5 stroke-current"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path></svg>
                </label>
                <ul tabindex="0" class="dropdown-content menu p-2 shadow-lg bg-white rounded-box w-40 z-[10] border border-gray-100">
                  <li><a @click="openViewModal(round)" class="text-slate-600">รายละเอียด</a></li>
                  <li><a class="text-slate-600">แก้ไข</a></li>
                  <div class="divider my-0"></div>
                  <li><a @click="handleDelete(round.ID)" class="text-error hover:bg-red-50">ลบรายการ</a></li>
                </ul>
              </div>
            </div>

            <div class="p-5 space-y-4 cursor-pointer" @click="openViewModal(round)">
              <div class="flex items-center gap-3 text-sm text-gray-600">
                 <div class="w-8 h-8 rounded-full bg-blue-50 flex items-center justify-center text-blue-600">
                   <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                 </div>
                 <div>
                   <p class="font-semibold text-slate-700">วันที่ {{ new Date(round.date).toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: '2-digit'}) }}</p>
                   <p class="text-xs">{{ round.start_time }} - {{ round.end_time }}</p>
                 </div>
              </div>
              
              <div class="flex items-center gap-3 text-sm text-gray-600">
                 <div class="w-8 h-8 rounded-full flex items-center justify-center" 
                  :class="round.mode === 'Onsite' ? 'bg-orange-50 text-orange-600' : 'bg-purple-50 text-purple-600'">
                   <svg v-if="round.mode === 'Onsite'" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" /></svg>
                   <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor"><path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z" /></svg>
                 </div>
                 <div>
                   <p class="font-semibold text-slate-700">{{ round.mode }}</p>
                   <p class="text-xs truncate max-w-[180px]">{{ round.mode === 'Onsite' ? round.location_name : 'Online Meeting' }}</p>
                 </div>
              </div>

              <div class="divider my-1"></div>

              <div class="flex justify-between items-end">
                <div class="flex-1 mr-4">
                  <div class="flex justify-between text-xs mb-1">
                    <span class="text-gray-500">การจอง ({{ Math.round((round.booked_slots/round.total_slots)*100) }}%)</span>
                    <span class="font-bold text-slate-700">{{ round.booked_slots }}/{{ round.total_slots }}</span>
                  </div>
                  <progress class="progress w-full h-2" 
                    :class="round.booked_slots === round.total_slots ? 'progress-error' : 'progress-primary'" 
                    :value="round.booked_slots" :max="round.total_slots"></progress>
                </div>

                <div class="avatar-group -space-x-3">
                  <div v-for="(intr, idx) in round.interviewers.slice(0, 3)" :key="idx" class="avatar placeholder border-white">
                    <div class="bg-neutral-focus text-neutral-content w-8 h-8 text-[10px]">
                      <span>{{ getInitials(intr) }}</span>
                    </div>
                  </div>
                  <div v-if="round.interviewers.length > 3" class="avatar placeholder border-white">
                    <div class="bg-gray-200 text-gray-600 w-8 h-8 text-[10px]">
                      <span>+{{ round.interviewers.length - 3 }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="flex flex-col items-center justify-center py-20 text-center animate-fade-in">
          <div class="bg-gray-50 p-6 rounded-full mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
          </div>
          <h3 class="text-lg font-bold text-slate-600">ไม่มีรอบสัมภาษณ์ที่เปิดอยู่</h3>
          <p class="text-gray-500">สร้างรอบใหม่ได้ที่ปุ่มมุมขวาบน</p>
        </div>
      </div>

      <div v-else class="animate-fade-in">
        <div class="overflow-x-auto bg-white rounded-2xl border border-gray-100">
          <table class="table w-full">
            <thead class="bg-slate-50 text-slate-600 uppercase text-xs">
              <tr>
                <th class="py-4">วันที่</th>
                <th>ชื่อรอบสัมภาษณ์</th>
                <th>ประเภททุน</th>
                <th>สถิติการจอง</th>
                <th>สถานะ</th>
                <th class="text-right">จัดการ</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="round in filteredRounds" :key="round.ID" class="hover:bg-slate-50 transition-colors">
                <td class="font-medium text-slate-700">
                  <div class="flex flex-col">
                    <span>{{ new Date(round.date).toLocaleDateString('th-TH', { dateStyle: 'medium' }) }}</span>
                    <span class="text-xs text-gray-400 font-light">{{ round.start_time }} - {{ round.end_time }}</span>
                  </div>
                </td>
                <td>
                  <div class="font-bold text-[#1e3a8a]">{{ round.name }}</div>
                  <div class="text-xs text-gray-500 flex gap-2 items-center mt-1">
                     <span class="badge badge-xs badge-ghost">{{ round.mode }}</span>
                     {{ round.mode === 'Onsite' ? round.location_name : 'Online' }}
                  </div>
                </td>
                <td>
                  <div class="badge badge-outline text-xs">{{ round.scholarship_name }}</div>
                </td>
                <td>
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-semibold">{{ round.booked_slots }}/{{ round.total_slots }}</span>
                    <progress class="progress progress-primary w-16 h-1.5" :value="round.booked_slots" :max="round.total_slots"></progress>
                  </div>
                </td>
                <td>
                  <div class="badge bg-gray-100 text-gray-500 border-none">เสร็จสิ้น</div>
                </td>
                <td class="text-right">
                  <button @click="openViewModal(round)" class="btn btn-ghost btn-xs text-blue-600">ดูข้อมูล</button>
                </td>
              </tr>
              <tr v-if="filteredRounds.length === 0">
                <td colspan="6" class="text-center py-12 text-gray-400">
                  ไม่พบประวัติย้อนหลัง
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

    </div>

    <div v-if="isModalOpen" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="bg-white w-full max-w-5xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in">
        
        <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
          <div>
            <h2 class="text-xl font-bold text-[#1e3a8a]">
              {{ modalMode === 'create' ? 'สร้างรอบสัมภาษณ์ใหม่' : 'รายละเอียดรอบสัมภาษณ์' }}
            </h2>
            <p class="text-xs text-gray-500 mt-1">
              {{ modalMode === 'create' ? 'กรอกข้อมูลเพื่อสร้าง Slot อัตโนมัติ' : 'ตรวจสอบรายชื่อผู้จองและจัดการข้อมูล' }}
            </p>
          </div>
          <button @click="isModalOpen = false" class="btn btn-circle btn-ghost btn-sm text-gray-500">✕</button>
        </div>

        <div class="flex-1 overflow-y-auto p-6 bg-[#f8fafc]">
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div class="lg:col-span-1 space-y-6">
              <div class="card bg-white shadow-sm border border-gray-100 p-5 space-y-4">
                <h3 class="font-bold text-slate-700 border-b pb-2">ข้อมูลทั่วไป</h3>
                <div class="form-control w-full">
                  <label class="label"><span class="label-text font-medium">ชื่อรอบสัมภาษณ์</span></label>
                  <input v-model="formData.name" type="text" placeholder="เช่น รอบเรียนดี 1/67" class="input input-bordered input-sm w-full" :disabled="modalMode==='view'"/>
                </div>
                <div class="form-control w-full">
                  <label class="label"><span class="label-text font-medium">สำหรับทุน</span></label>
                  <select v-model="formData.scholarship_id" class="select select-bordered select-sm w-full" :disabled="modalMode==='view'">
                    <option disabled value="">เลือกทุน...</option>
                    <option v-for="s in scholarships" :key="s.id" :value="s.id">{{ s.name }}</option>
                  </select>
                </div>
                <div class="form-control w-full">
                  <label class="label"><span class="label-text font-medium">รายละเอียดเพิ่มเติม</span></label>
                  <textarea v-model="formData.description" class="textarea textarea-bordered text-sm h-20" placeholder="เช่น สิ่งที่ต้องเตรียมมา..." :disabled="modalMode==='view'"></textarea>
                </div>
              </div>
              <div class="card bg-white shadow-sm border border-gray-100 p-5 space-y-4">
                <h3 class="font-bold text-slate-700 border-b pb-2">เวลา & สถานที่</h3>
                <div class="form-control w-full">
                  <label class="label"><span class="label-text font-medium">วันที่</span></label>
                  <input v-model="formData.date" type="date" class="input input-bordered input-sm w-full" :disabled="modalMode==='view'"/>
                </div>
                <div class="grid grid-cols-2 gap-2">
                  <div class="form-control">
                    <label class="label"><span class="label-text font-medium">เริ่ม</span></label>
                    <input v-model="formData.start_time" type="time" class="input input-bordered input-sm w-full" :disabled="modalMode==='view'"/>
                  </div>
                  <div class="form-control">
                    <label class="label"><span class="label-text font-medium">สิ้นสุด</span></label>
                    <input v-model="formData.end_time" type="time" class="input input-bordered input-sm w-full" :disabled="modalMode==='view'"/>
                  </div>
                </div>
                <div class="form-control w-full">
                  <label class="label"><span class="label-text font-medium">ระยะเวลา/คน (นาที)</span></label>
                  <select v-model="formData.slot_duration" class="select select-bordered select-sm w-full" :disabled="modalMode==='view'">
                    <option :value="10">10 นาที</option>
                    <option :value="15">15 นาที</option>
                    <option :value="20">20 นาที</option>
                    <option :value="30">30 นาที</option>
                    <option :value="60">1 ชั่วโมง</option>
                  </select>
                </div>
                <div class="divider my-2"></div>
                <div class="form-control">
                  <label class="label cursor-pointer justify-start gap-4">
                    <span class="label-text font-medium">รูปแบบ:</span>
                    <label class="flex items-center gap-2 cursor-pointer">
                      <input type="radio" v-model="formData.mode" value="Onsite" class="radio radio-primary radio-sm" :disabled="modalMode==='view'"/>
                      <span class="text-sm">Onsite</span>
                    </label>
                    <label class="flex items-center gap-2 cursor-pointer">
                      <input type="radio" v-model="formData.mode" value="Online" class="radio radio-primary radio-sm" :disabled="modalMode==='view'"/>
                      <span class="text-sm">Online</span>
                    </label>
                  </label>
                </div>
                <div v-if="formData.mode === 'Onsite'" class="form-control w-full animate-fade-in">
                  <label class="label"><span class="label-text font-medium">สถานที่</span></label>
                  <select v-model="formData.location_id" class="select select-bordered select-sm w-full" :disabled="modalMode==='view'">
                    <option disabled value="">เลือกห้องสัมภาษณ์...</option>
                    <option v-for="loc in locations" :key="loc.id" :value="loc.id">{{ loc.name }}</option>
                  </select>
                </div>
                <div v-else class="form-control w-full animate-fade-in">
                  <label class="label"><span class="label-text font-medium">Meeting Link</span></label>
                  <input v-model="formData.meeting_link" type="text" placeholder="https://zoom.us/..." class="input input-bordered input-sm w-full" :disabled="modalMode==='view'"/>
                </div>
              </div>
            </div>
            <div class="lg:col-span-2 space-y-6">
              <div class="card bg-white shadow-sm border border-gray-100 p-5">
                <h3 class="font-bold text-slate-700 mb-3">กรรมการสัมภาษณ์</h3>
                <div class="flex flex-wrap gap-2">
                   <label v-for="interviewer in interviewersList" :key="interviewer.id" 
                    class="cursor-pointer border rounded-lg px-3 py-2 flex items-center gap-2 transition-all hover:bg-gray-50"
                    :class="formData.interviewer_ids.includes(String(interviewer.id)) ? 'border-blue-500 bg-blue-50 ring-1 ring-blue-500' : 'border-gray-200'">
                    <input type="checkbox" :value="String(interviewer.id)" v-model="formData.interviewer_ids" class="checkbox checkbox-primary checkbox-xs" :disabled="modalMode==='view'"/>
                    <span class="text-sm">{{ interviewer.name }}</span>
                   </label>
                </div>
              </div>
              <div class="card bg-white shadow-sm border border-gray-100 p-5 flex-1 min-h-[400px]">
                <div class="flex items-center justify-between mb-4">
                  <h3 class="font-bold text-slate-700 flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {{ modalMode === 'create' ? 'พรีวิวตารางเวลา' : 'ตารางการจอง' }}
                  </h3>
                  <div class="text-xs text-gray-500">
                    ทั้งหมด: {{ previewSlots.length }} สล็อต
                  </div>
                </div>
                <div v-if="previewSlots.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
                  <div v-for="(slot, idx) in previewSlots" :key="idx" 
                    class="relative p-3 rounded-xl border text-center transition-all duration-200 flex flex-col items-center justify-center gap-1 group"
                    :class="[
                      slot.isBooked 
                        ? 'bg-gray-100 border-gray-200 cursor-help' 
                        : 'bg-white border-green-200 hover:border-green-500 hover:shadow-md cursor-default'
                    ]">
                    <span class="font-bold text-lg" :class="slot.isBooked ? 'text-gray-400' : 'text-slate-700'">
                      {{ slot.time }}
                    </span>
                    <div v-if="slot.isBooked" class="badge badge-ghost badge-xs text-gray-400">จองแล้ว</div>
                    <div v-else class="badge badge-success badge-xs text-white">ว่าง</div>
                    <div v-if="slot.isBooked && modalMode === 'view'" 
                      class="absolute bottom-full mb-2 bg-slate-800 text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap z-10">
                      {{ slot.studentName }}
                    </div>
                  </div>
                </div>
                <div v-else class="flex flex-col items-center justify-center h-40 text-gray-400 border-2 border-dashed border-gray-200 rounded-xl">
                  <p>กรุณาระบุเวลาเริ่ม-จบ และระยะเวลา</p>
                  <p class="text-xs">ระบบจะคำนวณรอบให้อัตโนมัติ</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="p-4 border-t bg-white flex justify-end gap-2">
          <button @click="isModalOpen = false" class="btn btn-ghost text-gray-500">ยกเลิก</button>
          <button v-if="modalMode === 'create'" @click="saveRound" class="btn bg-[#1e3a8a] text-white hover:bg-[#152c6f]">
            ยืนยันการสร้าง
          </button>
        </div>

      </div>
    </div>

  </div>
</template>

<style scoped>
.animate-pop-in {
  animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
.animate-fade-in {
  animation: fade 0.3s ease forwards;
}
@keyframes pop-in {
  0% { opacity: 0; transform: scale(0.95); }
  100% { opacity: 1; transform: scale(1); }
}
@keyframes fade {
  from { opacity: 0; transform: translateY(-5px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Custom Scrollbar */
::-webkit-scrollbar { width: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }
::-webkit-scrollbar-thumb:hover { background: #94a3b8; }
</style>