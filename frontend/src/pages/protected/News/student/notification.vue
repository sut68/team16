<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ArrowLeft, Bell } from 'lucide-vue-next';
import { getAllNewsPosts } from '@/services/api/news_post';

const router = useRouter();
const API_BASE_URL = 'http://localhost:8080';

const loading = ref(true);
const notifications = ref<any[]>([]);
const lastSeenId = ref(0);

// --- Helper: Format วันที่ ---
const formatTime = (dateStr: string) => {
  if (!dateStr) return '';
  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays === 0) return 'วันนี้';
  if (diffDays === 1) return 'เมื่อวาน';
  return date.toLocaleDateString('th-TH', { day: 'numeric', month: 'short' });
};

// --- Fetch Data ---
const fetchData = async () => {
  loading.value = true;
  try {
    // 1. อ่านค่า ID ล่าสุดที่เคยดูจาก LocalStorage
    const storedId = localStorage.getItem('lastSeenNewsId');
    lastSeenId.value = storedId ? parseInt(storedId) : 0;

    // 2. ดึงข่าวทั้งหมด
    const posts = await getAllNewsPosts();

    // 3. เรียงลำดับจาก ID มากไปน้อย (ใหม่สุดอยู่บน)
    const sortedPosts = posts.sort((a: any, b: any) => b.ID - a.ID);

    notifications.value = sortedPosts.map((post: any) => ({
      id: post.ID,
      title: post.title,
      desc: post.post_detail,
      // เช็คว่าใหม่กว่าที่เคยดูหรือไม่
      isNew: (post.ID || post.id) > lastSeenId.value, 
      createdAt: post.CreatedAt,
      // 🔥 แก้ไข: ใช้ undefined แทน null เพื่อแก้ Type Error ในบาง IDE
      image: post.file_path ? `${API_BASE_URL}/${post.file_path}` : undefined
    }));

  } catch (error) {
    console.error("Failed to load notifications:", error);
  } finally {
    loading.value = false;
  }
};

// --- Action ---
const markAllAsRead = () => {
  if (notifications.value.length > 0) {
    // หา ID สูงสุด
    const maxId = Math.max(...notifications.value.map(n => n.id));
    // บันทึกว่าอ่านถึงอันล่าสุดแล้ว
    localStorage.setItem('lastSeenNewsId', maxId.toString());
  }
  router.back();
};

const goToNews = (id: number) => {
  // 1. อัปเดตว่าอ่านแล้ว
  if (notifications.value.length > 0) {
     const maxId = Math.max(...notifications.value.map(n => n.id));
     localStorage.setItem('lastSeenNewsId', maxId.toString());
  }

  // 2. 🔥 แก้ตรงนี้: แทนที่จะไปหน้า NewsPreview ให้กลับไปหน้าหลัก + ส่งรหัสข่าวไปด้วย
  // สมมติว่า path หน้าหลักคือ '/news' (เช็คใน router.ts ของพี่นะว่า path อะไร)
  router.push({ 
      path: '/news',      // <--- เปลี่ยนเป็น Path ของหน้า NewsList
      query: { openId: id } // <--- ฝากรหัสข่าวไปบอกให้เปิด Modal
  });
};

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="flex h-screen bg-[#F5F7FA] font-prompt overflow-hidden rounded-tl-[50px] rounded-bl-[50px] shadow-2xl ml-0 relative z-10">
    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
      
      <header class="px-6 py-4 bg-white/80 backdrop-blur-md sticky top-0 z-20 border-b border-gray-100 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <button 
            @click="markAllAsRead" 
            class="w-10 h-10 rounded-full bg-white border border-gray-200 flex items-center justify-center text-gray-600 hover:bg-[#8B0025] hover:text-white transition-all shadow-sm"
          >
            <ArrowLeft :size="20" />
          </button>
          <h1 class="text-xl font-bold text-[#1e3a8a]">การแจ้งเตือน</h1>
        </div>
        
        <button @click="markAllAsRead" class="text-xs text-gray-400 hover:text-[#8B0025] underline">
           ทำเครื่องหมายว่าอ่านแล้ว
        </button>
      </header>

      <div class="flex-1 overflow-y-auto custom-scrollbar p-6">
        
        <div v-if="loading" class="flex justify-center items-center h-40">
           <span class="loading loading-spinner text-[#1e3a8a]"></span>
        </div>

        <div v-else-if="notifications.length === 0" class="flex flex-col items-center justify-center h-64 text-gray-400">
           <Bell :size="48" class="mb-2 opacity-20" />
           <p>ไม่มีการแจ้งเตือน</p>
        </div>

        <div v-else class="space-y-3">
           <div 
             v-for="item in notifications" 
             :key="item.id"
             @click="goToNews(item.id)"
             :class="`
                relative p-4 rounded-2xl border cursor-pointer transition-all hover:shadow-md flex gap-4 items-start
                ${item.isNew ? 'bg-white border-red-100 shadow-sm' : 'bg-gray-50 border-transparent opacity-70'}
             `"
           >
              <div class="w-12 h-12 rounded-xl bg-gray-200 flex-shrink-0 overflow-hidden">
                 <img v-if="item.image" :src="item.image" class="w-full h-full object-cover" />
                 <div v-else class="w-full h-full flex items-center justify-center bg-blue-50 text-blue-300">
                    <Bell :size="20" />
                 </div>
              </div>

              <div class="flex-1 min-w-0">
                 <div class="flex justify-between items-start">
                    <h3 class="font-bold text-gray-800 text-sm truncate pr-2">{{ item.title }}</h3>
                    <span class="text-[10px] text-gray-400 whitespace-nowrap">{{ formatTime(item.createdAt) }}</span>
                 </div>
                 <p class="text-xs text-gray-500 line-clamp-2 mt-1">{{ item.desc }}</p>
              </div>

              <div v-if="item.isNew" class="absolute top-4 right-4 w-2 h-2 bg-red-500 rounded-full animate-pulse shadow-red-200 shadow-lg"></div>
           </div>
        </div>

      </div>
    </main>
  </div>
</template>