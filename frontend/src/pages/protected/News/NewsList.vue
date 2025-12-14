<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getAllNewsPosts, deleteNewsPost } from '@/services/api/news_post';
import type { NewsPost } from '@/interfaces/news_post';

// Import Child Components
import AddNewsPostScreen from './CreateNews.vue'; 
import EditNewsPostScreen from './EditNews.vue'; 
import EditPrivacyModal from './EditPrivacy.vue'; 
import NewsPreviewModal from './PreviewPage.vue'; 

// --- State Definitions ---
const isLoading = ref(false);
const newsList = ref<NewsPost[]>([]);

const viewMode = ref<'list' | 'create' | 'edit'>('list'); 
const editId = ref<number>(0); 

const isPrivacyOpen = ref(false);
const privacyTargetId = ref(0);
const isPreviewOpen = ref(false);
const previewId = ref(0);

// --- Custom Popup System Definition ---
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

/**
 * ฟังก์ชันสำหรับเรียกแสดงหน้าต่างแจ้งเตือน (Modal Dialog)
 * รองรับการทำงานแบบ Await เพื่อรอผลลัพธ์จากผู้ใช้
 */
const showPopup = (type: 'success' | 'error' | 'confirm', title: string, message: string) => {
  return new Promise<boolean>((resolve) => {
    popup.value = {
      isOpen: true,
      type,
      title,
      message,
      resolve
    };
  });
};

/**
 * ฟังก์ชันสำหรับปิดหน้าต่างแจ้งเตือนและส่งคืนค่าผลลัพธ์
 */
const closePopup = (result: boolean) => {
  if (popup.value.resolve) {
    popup.value.resolve(result);
  }
  popup.value.isOpen = false;
};

// --- Helper Functions: Data Formatting ---
const getStatusLabel = (id: number) => {
  switch (id) {
    case 1: return 'เผยแพร่สาธารณะ (Public)';
    case 2: return 'ฉบับร่าง (Draft)';
    case 3: return 'จัดเก็บ (Archived)';
    case 5: return 'เฉพาะสมาชิก (Members Only)';
    case 4: return 'ลบ (Deleted)';
    default: return 'ไม่ระบุสถานะ';
  }
};

const getStatusColor = (id: number) => {
  switch (id) {
    case 1: return 'bg-green-500';
    case 2: return 'bg-orange-400';
    case 3: return 'bg-slate-400';
    case 5: return 'bg-indigo-500';
    case 4: return 'bg-red-500';
    default: return 'bg-gray-200';
  }
};

// --- API Integration ---
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
    console.error("Error fetching news:", error);
    showPopup('error', 'เกิดข้อผิดพลาด', 'ไม่สามารถโหลดข้อมูลข่าวสารได้ กรุณาลองใหม่อีกครั้ง');
  } finally {
    isLoading.value = false;
  }
};

// --- View State Handlers ---
const handleCreateNews = () => { viewMode.value = 'create'; };
const handleCloseCreate = () => { viewMode.value = 'list'; };

const handleCreateSuccess = () => {
  viewMode.value = 'list';
  fetchNews();
  showPopup('success', 'ดำเนินการสำเร็จ', 'เพิ่มข่าวสารใหม่เข้าสู่ระบบเรียบร้อยแล้ว');
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
  showPopup('success', 'ดำเนินการสำเร็จ', 'ปรับปรุงข้อมูลข่าวสารเรียบร้อยแล้ว');
};

// --- Delete Operation Handler ---
const handleDelete = async (id: number) => {
  // เรียกใช้งาน Confirmation Dialog
  const isConfirmed = await showPopup(
    'confirm', 
    'ยืนยันการลบข้อมูล', 
    'ท่านต้องการลบข่าวสารนี้ใช่หรือไม่? การกระทำนี้ไม่สามารถเรียกคืนได้'
  );

  if (isConfirmed) {
    try {
      await deleteNewsPost(id);
      // อัปเดตรายการหน้าบ้านทันที (Optimistic UI)
      newsList.value = newsList.value.filter(item => item.ID !== id);
      
      showPopup('success', 'ดำเนินการสำเร็จ', 'ข้อมูลถูกลบออกจากระบบเรียบร้อยแล้ว');

    } catch (error) {
      console.error("Delete operation failed:", error);
      showPopup('error', 'เกิดข้อผิดพลาด', 'ไม่สามารถลบข้อมูลได้ กรุณาติดต่อผู้ดูแลระบบ');
    }
  }
};

// --- Modal Trigger Handlers ---
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
  showPopup('success', 'ดำเนินการสำเร็จ', 'สถานะการแสดงผลถูกบันทึกเรียบร้อยแล้ว');
};

// --- Lifecycle Hook ---
onMounted(() => {
  fetchNews();
});
</script>

<template>
  <div class="min-h-screen bg-[#f0f2f5] relative">
    
    <div v-if="viewMode === 'create'" class="w-full h-full">
      <AddNewsPostScreen @close="handleCloseCreate" @success="handleCreateSuccess" />
    </div>

    <div v-else-if="viewMode === 'edit'" class="w-full h-full">
      <EditNewsPostScreen 
        :id="editId" 
        @close="handleCloseEdit" 
        @success="handleEditSuccess" 
      />
    </div>

    <div v-else class="p-6 font-sans text-slate-800">
      
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-[#1e3a8a]">จัดการข่าวสารประชาสัมพันธ์</h2>
        <button 
            @click="handleCreateNews"
            class="btn bg-white border-none shadow-sm text-[#1e3a8a] hover:bg-blue-50 hover:shadow-md rounded-full px-6 gap-2 h-11 font-bold transition-all"
        >
          <span class="text-xl leading-none font-normal">+</span> เพิ่มข่าวทุนใหม่
        </button>
      </div>

      <div v-if="isLoading" class="text-center py-20 text-gray-500">
        <span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span>
        <p class="mt-2 text-sm">กำลังประมวลผล...</p>
      </div>

      <div v-else class="space-y-4 pb-10">
        <transition-group name="fade" tag="div" class="space-y-4">
          <div 
            v-for="item in newsList" 
            :key="item.ID" 
            @click="handleCardClick(item.ID)"
            class="card bg-white shadow-sm rounded-2xl cursor-pointer border border-transparent hover:border-blue-200 hover:shadow-md transition-all duration-300 transform hover:-translate-y-1 relative overflow-visible z-0 hover:z-50 focus-within:z-50"
          >
            <div class="card-body p-5 flex flex-row items-center justify-between min-h-[5rem]">
              
              <div class="flex flex-col pr-4 gap-1">
                <div class="flex items-center gap-2">
                    <span 
                        class="w-2.5 h-2.5 rounded-full"
                        :class="getStatusColor(item.status_news_id)"
                    ></span>

                    <h3 class="font-bold text-[#1e3a8a] text-lg leading-tight line-clamp-1">
                        {{ item.title }}
                    </h3>
                </div>
                <div class="flex items-center gap-3 pl-4 pt-1">
                    <span class="text-xs text-gray-500 font-medium">
                        {{ getStatusLabel(item.status_news_id) }}
                    </span>
                </div>
              </div>
              
              <div class="dropdown dropdown-end" @click.stop>
                <div tabindex="0" role="button" class="btn btn-ghost btn-circle btn-sm text-gray-400 hover:bg-blue-50 hover:text-[#1e3a8a]">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM12.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM18.75 12a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0Z" />
                  </svg>
                </div>
                <ul tabindex="0" class="dropdown-content z-[999] menu p-2 shadow-xl bg-white rounded-xl w-48 border border-gray-100 mt-2">
                  <li><a @click="handlePreview(item.ID)" class="hover:bg-blue-50 hover:text-blue-700">ดูตัวอย่าง</a></li>
                  <li><a @click="handleEdit(item.ID)" class="hover:bg-blue-50 hover:text-blue-700">แก้ไขเนื้อหา</a></li>
                  <li><a @click="handleEditPrivacy(item.ID)" class="hover:bg-blue-50 hover:text-blue-700">แก้ไขความเป็นส่วนตัว</a></li>
                  <div class="divider my-1"></div>
                  <li><a @click="handleDelete(item.ID)" class="text-red-500 hover:bg-red-50 hover:text-red-700">ลบข้อมูล</a></li>
                </ul>
              </div>

            </div>
          </div>
        </transition-group>

        <div v-if="newsList.length === 0" class="flex flex-col items-center justify-center py-20 text-gray-400">
            <p>ไม่พบข้อมูลข่าวสารในระบบ</p>
        </div>

        <EditPrivacyModal 
          :isOpen="isPrivacyOpen" 
          :newsId="privacyTargetId" 
          @close="isPrivacyOpen = false" 
          @save="handlePrivacySaved" 
        />

        <NewsPreviewModal 
            :isOpen="isPreviewOpen" 
            :id="previewId" 
            @close="isPreviewOpen = false" 
        />
      </div>
    </div>

    <div v-if="popup.isOpen" class="fixed inset-0 z-[2000] flex items-center justify-center bg-black/50 backdrop-blur-sm transition-all p-4">
        <div class="bg-white w-full max-w-sm rounded-2xl shadow-2xl p-6 text-center transform scale-100 animate-bounce-in">
            
            <div class="mx-auto mb-4 w-16 h-16 flex items-center justify-center rounded-full"
                :class="{
                    'bg-green-100 text-green-600': popup.type === 'success',
                    'bg-red-100 text-red-600': popup.type === 'error',
                    'bg-orange-100 text-orange-500': popup.type === 'confirm'
                }"
            >
                <svg v-if="popup.type === 'success'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-8 h-8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                </svg>
                <svg v-else-if="popup.type === 'error'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-8 h-8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-8 h-8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z" />
                </svg>
            </div>

            <h3 class="text-xl font-bold text-slate-800 mb-2">{{ popup.title }}</h3>
            <p class="text-gray-500 mb-6 text-sm">{{ popup.message }}</p>

            <div class="flex gap-3 justify-center">
                <button 
                    v-if="popup.type === 'confirm'"
                    @click="closePopup(false)" 
                    class="btn btn-md bg-white border border-gray-300 text-gray-700 hover:bg-gray-50 flex-1"
                >
                    ยกเลิก
                </button>
                
                <button 
                    @click="closePopup(true)" 
                    class="btn btn-md border-none text-white flex-1"
                    :class="{
                        'bg-green-600 hover:bg-green-700': popup.type === 'success',
                        'bg-red-600 hover:bg-red-700': popup.type === 'error' || popup.type === 'confirm'
                    }"
                >
                    {{ popup.type === 'confirm' ? 'ยืนยัน' : 'ตกลง' }}
                </button>
            </div>
        </div>
    </div>

  </div>
</template>

<style scoped>
@keyframes bounceIn {
    0% { opacity: 0; transform: scale(0.9); }
    100% { opacity: 1; transform: scale(1); }
}
.animate-bounce-in {
    animation: bounceIn 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.fade-enter-active, .fade-leave-active { 
    transition: opacity 0.2s ease, transform 0.2s ease; 
}
.fade-enter-from, .fade-leave-to { 
    opacity: 0; 
    transform: translateY(10px); 
}
</style>