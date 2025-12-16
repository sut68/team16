<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { getAllNewsPosts, deleteNewsPost } from '@/services/api/news_post';
import type { NewsPost } from '@/interfaces/news_post';

// Components
import AddNewsPostScreen from './CreateNews.vue'; 
import EditNewsPostScreen from './EditNews.vue'; 
import EditPrivacyModal from './EditPrivacy.vue'; 
import NewsPreviewModal from './PreviewPage.vue'; 

// --- State ---
const isLoading = ref(false);
const newsList = ref<NewsPost[]>([]);

const viewMode = ref<'list' | 'create' | 'edit'>('list'); 
const editId = ref<number>(0); 

const isPrivacyOpen = ref(false);
const privacyTargetId = ref(0);
const isPreviewOpen = ref(false);
const previewId = ref(0);

const activeStatus = ref<number | 'all'>('all');
const searchQuery = ref('');

// Popup
interface PopupState {
  isOpen: boolean;
  type: 'success' | 'error' | 'confirm';
  title: string;
  message: string;
  resolve: ((value: boolean) => void) | null;
}
const popup = ref<PopupState>({
  isOpen: false,
  type: 'success',
  title: '',
  message: '',
  resolve: null
});
const showPopup = (type: 'success' | 'error' | 'confirm', title: string, message: string) => {
  return new Promise<boolean>((resolve) => {
    popup.value = { isOpen: true, type, title, message, resolve };
  });
};
const closePopup = (result: boolean) => {
  if (popup.value.resolve) popup.value.resolve(result);
  popup.value.isOpen = false;
};

// --- Helper ---
const getStatusLabel = (id: number) => {
  switch (id) {
    case 1: return 'สาธารณะ';
    case 2: return 'ฉบับร่าง';
    case 3: return 'จัดเก็บ';
    case 4: return 'เฉพาะสมาชิก';
    default: return 'ไม่ระบุ';
  }
};
const getStatusBadgeClass = (id: number) => {
  switch (id) {
    case 1: return 'badge-success text-white';        
    case 2: return 'badge-warning text-white';        
    case 3: return 'badge-ghost text-slate-500';      
    case 5: return 'badge-info text-white';           
    case 4: return 'badge-error text-white';          
    default: return 'badge-ghost';
  }
};
const formatDate = (dateString: string) => {
  if (!dateString) return '-';
  return new Date(dateString).toLocaleDateString('th-TH', {
    day: '2-digit', month: 'short', year: 'numeric'
  });
};

// --- API ---
const fetchNews = async () => {
  isLoading.value = true;
  try {
    const data = await getAllNewsPosts();
    if (Array.isArray(data)) {
      newsList.value = data.sort((a, b) => b.ID - a.ID); 
    } else {
      newsList.value = [];
    }
  } catch (error) {
    console.error("Error:", error);
    showPopup('error', 'เกิดข้อผิดพลาด', 'โหลดข้อมูลไม่สำเร็จ');
  } finally {
    isLoading.value = false;
  }
};

// --- Computed ---
const filteredNews = computed(() => {
  return newsList.value.filter(item => {
    const matchStatus = activeStatus.value === 'all' || item.status_news_id === activeStatus.value;
    const matchSearch = item.title.toLowerCase().includes(searchQuery.value.toLowerCase());
    return matchStatus && matchSearch;
  });
});

// --- Handlers ---
const handleCreateNews = () => { viewMode.value = 'create'; };
const handleCloseCreate = () => { viewMode.value = 'list'; };
const handleCreateSuccess = () => {
  viewMode.value = 'list';
  fetchNews();
  showPopup('success', 'สำเร็จ', 'เพิ่มข่าวสารเรียบร้อย');
};

const handleEdit = (id: number) => {
  editId.value = id;     
  viewMode.value = 'edit'; 
};
const handleCloseEdit = () => {
  viewMode.value = 'list';
  editId.value = 0;
};
const handleEditSuccess = () => {
  viewMode.value = 'list';
  editId.value = 0;
  fetchNews(); 
  showPopup('success', 'สำเร็จ', 'แก้ไขข้อมูลเรียบร้อย');
};

const handleDelete = async (id: number) => {
  const isConfirmed = await showPopup('confirm', 'ยืนยันการลบ', 'ต้องการลบข่าวสารนี้ใช่หรือไม่?');
  if (isConfirmed) {
    try {
      await deleteNewsPost(id);
      newsList.value = newsList.value.filter(item => item.ID !== id);
      showPopup('success', 'สำเร็จ', 'ลบข้อมูลเรียบร้อย');
    } catch (error) {
      showPopup('error', 'ผิดพลาด', 'ลบข้อมูลไม่สำเร็จ');
    }
  }
};

const handleCardClick = (id: number) => { handlePreview(id); };
const handlePreview = (id: number) => {
  previewId.value = id;
  isPreviewOpen.value = true;
};
const handleEditPrivacy = (id: number) => {
  privacyTargetId.value = id; 
  isPrivacyOpen.value = true; 
};
const handlePrivacySaved = () => {
  isPrivacyOpen.value = false;
  fetchNews();
  showPopup('success', 'สำเร็จ', 'บันทึกการตั้งค่าเรียบร้อย');
};

onMounted(() => { fetchNews(); });
</script>

<template>
<div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow overflow-visible font-sans text-slate-800">

  <div v-if="viewMode !== 'list'" class="w-full h-full animate-fade-in-up">
    <AddNewsPostScreen v-if="viewMode === 'create'" @close="handleCloseCreate" @success="handleCreateSuccess" />
    <EditNewsPostScreen v-else :id="editId" @close="handleCloseEdit" @success="handleEditSuccess" />
  </div>

  <div v-else class="flex flex-col h-full">
    
    <!-- Header -->
    <div class="flex flex-col xl:flex-row items-end xl:items-center justify-between gap-4 mb-6 border-b border-gray-200 pb-4">
      <div>
          <h2 class="text-xl font-bold text-[#1e3a8a] flex items-center gap-2">
              จัดการข่าวสารประชาสัมพันธ์
              <span class="badge badge-neutral text-white text-xs font-normal h-5">{{ newsList.length }} รายการ</span>
          </h2>
          <p class="text-sm text-gray-500 mt-1">บริหารจัดการข้อมูลข่าวทุนและกิจกรรมต่างๆ</p>
      </div>
      
      <button @click="handleCreateNews" class="btn btn-sm h-10 px-6 rounded-full bg-[#1e3a8a] hover:bg-[#152c6f] text-white border-none shadow-sm gap-2 font-medium">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-5 h-5">
          <path d="M10.75 4.75a.75.75 0 00-1.5 0v4.5h-4.5a.75.75 0 000 1.5h4.5v4.5a.75.75 0 001.5 0v-4.5h4.5a.75.75 0 000-1.5h-4.5v-4.5z" />
        </svg>
        เพิ่มข่าวใหม่
      </button>
    </div>

    <!-- Filter + Search -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-3 mb-4">
      <div class="flex items-center gap-2">
        <label class="text-sm text-gray-500">สถานะ:</label>
        <select v-model.number="activeStatus" class="select select-bordered select-sm rounded-full bg-white border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a]">
            <option :value="'all'">ทั้งหมด</option>
            <option :value="1">สาธารณะ</option>
            <option :value="2">ฉบับร่าง</option>
            <option :value="3">จัดเก็บ</option>
            <option :value="5">เฉพาะสมาชิก</option>
        </select>

      </div>
      <input type="text" v-model="searchQuery" placeholder="ค้นหาข่าว..." 
             class="input input-bordered input-sm rounded-full w-full md:w-64 border-gray-300 focus:border-[#1e3a8a] focus:ring-[#1e3a8a]" />
    </div>

    <!-- News Grid -->
    <div v-if="isLoading" class="flex-1 flex flex-col items-center justify-center text-gray-500">
      <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
      <p class="mt-2 text-sm">กำลังโหลดข้อมูล...</p>
    </div>

    <div v-else class="flex-1 overflow-y-auto pr-1 custom-scrollbar">
      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6 pb-10">
        <div v-for="item in filteredNews" :key="item.ID" @click="handleCardClick(item.ID)"
             class="card bg-white border border-gray-200 shadow-sm rounded-2xl hover:border-blue-200 hover:shadow-md transition-all duration-300 group !overflow-visible hover:z-50">
          <div class="card-body p-5 flex flex-col justify-between h-full">
            <div>
              <div class="flex justify-between items-start mb-3">
                <span class="badge badge-sm border-none px-2 py-3 font-medium" :class="getStatusBadgeClass(item.status_news_id)">
                  {{ getStatusLabel(item.status_news_id) }}
                </span>
                <div class="dropdown dropdown-end" @click.stop>
                  <div tabindex="0" role="button" class="btn btn-square btn-ghost btn-sm text-slate-400 hover:bg-slate-100 hover:text-[#1e3a8a]">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" /></svg>
                  </div>
                  <ul tabindex="0" class="dropdown-content z-[999] menu p-2 shadow-xl bg-white rounded-xl w-48 border border-gray-100 mt-1 text-sm">
                    <li><a @click="handlePreview(item.ID)">ดูตัวอย่าง</a></li>
                    <li><a @click="handleEdit(item.ID)">แก้ไข</a></li>
                    <li><a @click="handleEditPrivacy(item.ID)">ความเป็นส่วนตัว</a></li>
                    <div class="divider my-1"></div>
                    <li><a @click="handleDelete(item.ID)" class="text-error">ลบ</a></li>
                  </ul>
                </div>
              </div>
              <h3 class="font-bold text-[#1e3a8a] text-lg leading-snug line-clamp-2 mb-2 group-hover:underline">
                {{ item.title }}
              </h3>
              <p class="text-xs text-gray-400 flex items-center gap-1">
                สร้างเมื่อ: {{ formatDate(item.CreatedAt) }}
              </p>
            </div>
          </div>
        </div>

        <div v-if="filteredNews.length === 0" class="col-span-full flex flex-col items-center justify-center py-16 text-gray-400">
          <p>ไม่พบข่าวสารในระบบ</p>
        </div>
      </div>
    </div>
  </div>

  <EditPrivacyModal :isOpen="isPrivacyOpen" :newsId="privacyTargetId" @close="isPrivacyOpen = false" @save="handlePrivacySaved" />
  <NewsPreviewModal :isOpen="isPreviewOpen" :id="previewId" @close="isPreviewOpen = false" />

  <div v-if="popup.isOpen" class="fixed inset-0 z-[2000] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
    <div class="bg-white w-full max-w-sm rounded-2xl shadow-2xl p-6 text-center animate-bounce-in">
      <h3 class="text-lg font-bold text-slate-800 mb-2">{{ popup.title }}</h3>
      <p class="text-gray-500 mb-6 text-sm">{{ popup.message }}</p>
      <div class="flex gap-2 justify-center">
        <button v-if="popup.type === 'confirm'" @click="closePopup(false)" class="btn btn-sm btn-ghost text-gray-500">ยกเลิก</button>
        <button @click="closePopup(true)" class="btn btn-sm text-white border-none" 
                :class="popup.type === 'success' ? 'btn-success' : (popup.type === 'confirm' ? 'bg-[#1e3a8a]' : 'btn-error')">
          ตกลง
        </button>
      </div>
    </div>
  </div>

</div>
</template>
