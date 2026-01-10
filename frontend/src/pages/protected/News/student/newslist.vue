<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { Heart, Search, Bell, X } from 'lucide-vue-next'; 

// Import Component หน้า Preview (Modal)
import NewsPreview from './preview.vue';

// Import API Services
import { getAllNewsPosts } from '@/services/api/news_post'; 
import { getMyFavoriteNews, toggleStudentFavoriteNews } from '@/services/api/student_fav_news'; 

// อ่าน Base URL จาก ENV (ตัด /api ออกสำหรับ static files)
const API_BASE_URL = (import.meta.env.VITE_API_URL || 'http://localhost:8080/api').replace('/api', ''); 
const CURRENT_STUDENT_ID = 1; 

// --- State Variables ---
const newsItems = ref<any[]>([]);
const loading = ref(true);
const searchTerm = ref('');
const currentSlide = ref(0);
let slideInterval: any = null;
let pollingInterval: any = null;
const selectedNewsId = ref<number | null>(null);

// Notification State
const latestNewsIdFromAPI = ref(0); 
const lastSeenNewsId = ref(0); 
const showNotifications = ref(false);

// Filter State
const showFavoritesOnly = ref(false);

// --- Computed Properties ---

const highlightNews = computed(() => {
  return [...newsItems.value]
    .sort((a, b) => b.id - a.id)
    .slice(0, 3)
    .map((item, index) => ({
      ...item,
      bgGradient: [
        'bg-gradient-to-r from-[#1e3a8a] to-blue-400',
        'bg-gradient-to-r from-[#8B0025] to-red-400',
        'bg-gradient-to-r from-orange-500 to-amber-400'
      ][index % 3]
    }));
});

const filteredNewsList = computed(() => {
    let items = newsItems.value;
    if (showFavoritesOnly.value) {
        items = items.filter(item => item.isLiked);
    }
    if (searchTerm.value) {
        const term = searchTerm.value.toLowerCase();
        items = items.filter((item: any) => 
            item.title.toLowerCase().includes(term) || 
            item.caption.toLowerCase().includes(term)
        );
    }
    return items.sort((a, b) => b.id - a.id);
});

const recentNotifications = computed(() => {
    return [...newsItems.value].sort((a, b) => b.id - a.id).slice(0, 5);
});

const hasNewNotifications = computed(() => {
  return latestNewsIdFromAPI.value > lastSeenNewsId.value;
});

// --- Slider Logic ---
const nextSlide = () => {
  if (highlightNews.value.length > 0) {
    currentSlide.value = (currentSlide.value + 1) % highlightNews.value.length;
  }
};
const startAutoPlay = () => {
    clearInterval(slideInterval);
    slideInterval = setInterval(nextSlide, 4000);
};

// --- Actions ---
const toggleNotifications = () => {
  showNotifications.value = !showNotifications.value;
  if (showNotifications.value && hasNewNotifications.value) {
    lastSeenNewsId.value = latestNewsIdFromAPI.value;
    localStorage.setItem('lastSeenNewsId', latestNewsIdFromAPI.value.toString());
  }
};

const goToNews = (id: number) => {
    showNotifications.value = false;
    selectedNewsId.value = id;
};

const backToNewsList = () => { selectedNewsId.value = null; };

const handleToggleLike = async (newsId: number) => {
  const item = newsItems.value.find(n => n.id === newsId);
  if (item) item.isLiked = !item.isLiked;
  try { 
      await toggleStudentFavoriteNews(CURRENT_STUDENT_ID, newsId); 
  } catch (error) { 
      if (item) item.isLiked = !item.isLiked; 
  }
};

// --- Fetch Data ---
const fetchData = async (background = false) => {
  if (!background) loading.value = true;
  try {
    const storedId = localStorage.getItem('lastSeenNewsId');
    lastSeenNewsId.value = storedId ? Number(storedId) : 0;

    const [posts, favorites] = await Promise.all([
      getAllNewsPosts(),
      getMyFavoriteNews(CURRENT_STUDENT_ID).catch(() => [])
    ]);
    
    const favSet = new Set(favorites.map((f: any) => f.news_post_id));
    
    if (posts.length > 0) {
        const maxId = Math.max(...posts.map((p:any) => Number(p.ID || p.id)));
        latestNewsIdFromAPI.value = maxId;
    }

    // Helper to strip HTML tags
    const stripHtml = (html: string) => {
       if (!html) return "";
       return html.replace(/<[^>]*>?/gm, '');
    };

    // Filter: Show Public (1) and Members Only (4)
    const validPosts = posts.filter((post: any) => post.status_news_id === 1 || post.status_news_id === 4);

    newsItems.value = validPosts.map((post: any) => ({
      id: Number(post.ID || post.id), 
      title: post.title,
      desc: stripHtml(post.post_detail),
      caption: post.scholarship?.scholarship_name || 'ข่าวทั่วไป', // ชื่อทุน
      postDetail: stripHtml(post.post_detail),
      imagePath: post.file_path ? `${API_BASE_URL}/${post.file_path}` : null,
      isLiked: favSet.has(post.ID || post.id),
      createdAt: post.CreatedAt,
      isNewBadge: Number(post.ID || post.id) > lastSeenNewsId.value, 
      fallbackColor: ['bg-blue-100', 'bg-emerald-100', 'bg-orange-100', 'bg-purple-100'][(post.ID || post.id) % 4]
    }));

  } catch (error) {
    console.error("Failed to fetch news:", error);
  } finally {
    if (!background) loading.value = false;
  }
};

onMounted(() => { 
  fetchData(); 
  startAutoPlay();
  pollingInterval = setInterval(() => {
    fetchData(true);
  }, 30000); // 30 seconds
});
onUnmounted(() => { 
  clearInterval(slideInterval);
  if (pollingInterval) clearInterval(pollingInterval);
});
</script>

<template>
  <div class="flex h-screen bg-[#F5F7FA] font-prompt overflow-hidden rounded-tl-[50px] rounded-bl-[50px] shadow-2xl ml-0 relative z-10">
    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
      
      <NewsPreview v-if="selectedNewsId" :id="selectedNewsId" @back="backToNewsList" class="absolute inset-0 z-50" />

      <div class="flex flex-col h-full transition-all duration-300" :class="{ 'scale-95 opacity-50 blur-[2px]': selectedNewsId }">
          
          <header class="px-10 py-6 flex flex-col gap-6 bg-white/50 backdrop-blur-sm sticky top-0 z-30">
            <div class="flex justify-between items-center">
                <div class="flex items-center gap-4">
                    <h1 class="text-2xl font-bold text-[#1e3a8a]">ข่าวทุนการศึกษา</h1>
                    <div class="flex bg-gray-200/50 p-1 rounded-xl">
                        <button @click="showFavoritesOnly = false" class="px-4 py-1.5 rounded-lg text-sm font-medium transition-all duration-300" :class="!showFavoritesOnly ? 'bg-white text-[#1e3a8a] shadow-sm' : 'text-gray-500 hover:text-gray-700'">ทั้งหมด</button>
                        <button @click="showFavoritesOnly = true" class="px-4 py-1.5 rounded-lg text-sm font-medium transition-all duration-300 flex items-center gap-1.5" :class="showFavoritesOnly ? 'bg-white text-[#8B0025] shadow-sm' : 'text-gray-500 hover:text-gray-700'"><Heart :size="14" :class="showFavoritesOnly ? 'fill-[#8B0025]' : ''" /> ที่ถูกใจ</button>
                    </div>
                </div>
                <div class="flex items-center gap-4">
                  <div class="relative group">
                    <input v-model="searchTerm" type="text" placeholder="ค้นหาทุน..." class="w-64 pl-11 pr-4 py-2.5 bg-white border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-[#8B0025]/20 focus:border-[#8B0025] transition-all shadow-sm group-hover:shadow-md"/>
                    <Search class="absolute left-3.5 top-3 text-gray-400 group-hover:text-[#8B0025] transition-colors" :size="18" />
                  </div>
                  <div class="relative">
                      <button @click.stop="toggleNotifications" class="w-10 h-10 bg-white rounded-full border border-gray-200 flex items-center justify-center text-gray-500 hover:text-[#8B0025] hover:shadow-md transition-all relative group z-20 overflow-hidden">
                        <Bell :size="20" :class="hasNewNotifications ? 'animate-bell-ring text-[#8B0025]' : ''" />
                        <span v-if="hasNewNotifications" class="absolute top-2 right-2.5 flex h-2.5 w-2.5"><span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span><span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-red-500 border-2 border-white"></span></span>
                      </button>
                      <div v-if="showNotifications" class="absolute right-0 mt-3 w-80 bg-white rounded-2xl shadow-2xl border border-gray-100 z-50 overflow-hidden animate-fade-in-up origin-top-right">
                          <div class="px-4 py-3 border-b border-gray-50 bg-gray-50/50 flex justify-between items-center">
                              <span class="font-bold text-[#1e3a8a] text-sm">การแจ้งเตือนล่าสุด</span>
                              <button @click="showNotifications = false" class="text-gray-400 hover:text-red-500"><X :size="16" /></button>
                          </div>
                          <div class="max-h-80 overflow-y-auto custom-scrollbar">
                              <div v-if="recentNotifications.length === 0" class="p-6 text-center text-gray-400 text-sm">ไม่มีการแจ้งเตือน</div>
                              <div v-for="item in recentNotifications" :key="item.id" @click="goToNews(item.id)" class="p-3 border-b border-gray-50 hover:bg-blue-50/50 cursor-pointer flex gap-3 transition-colors relative">
                                  <div class="w-12 h-12 rounded-lg bg-gray-200 overflow-hidden flex-shrink-0">
                                      <img v-if="item.imagePath" :src="item.imagePath" class="w-full h-full object-cover" />
                                      <div v-else :class="`w-full h-full ${item.fallbackColor}`"></div>
                                  </div>
                                  <div class="flex-1 min-w-0">
                                      <h4 class="text-sm font-bold text-gray-700 truncate pr-4">{{ item.title }}</h4>
                                      <p class="text-xs text-gray-400 mt-0.5 truncate">{{ item.caption }}</p>
                                      <p class="text-[10px] text-gray-300 mt-1">{{ new Date(item.createdAt).toLocaleDateString('th-TH') }}</p>
                                  </div>
                                  <div v-if="item.isNewBadge" class="absolute top-3 right-3 w-2 h-2 bg-red-500 rounded-full animate-pulse"></div>
                              </div>
                          </div>
                      </div>
                      <div v-if="showNotifications" @click="showNotifications = false" class="fixed inset-0 z-10 cursor-default"></div>
                  </div>
                </div>
            </div>
          </header>

          <div class="flex-1 overflow-y-auto px-10 pb-10 custom-scrollbar">
            <div v-if="loading" class="flex justify-center items-center h-40"><span class="loading loading-spinner loading-lg text-[#1e3a8a]"></span></div>

            <div v-else>
                <div v-if="highlightNews.length > 0 && !showFavoritesOnly && !searchTerm" class="mb-8 mt-2"> 
                    <div class="relative w-full h-64 md:h-80 rounded-2xl overflow-hidden shadow-lg shadow-blue-900/10 group cursor-grab active:cursor-grabbing select-none">
                        <div class="flex h-full transition-transform duration-700 ease-in-out" :style="{ transform: `translateX(-${currentSlide * 100}%)` }">
                            <div v-for="item in highlightNews" :key="item.id" class="w-full flex-shrink-0 relative h-full cursor-pointer" @click="goToNews(item.id)">
                                <div v-if="item.imagePath" class="absolute inset-0">
                                    <img :src="item.imagePath" class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105 pointer-events-none" draggable="false" />
                                    <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent"></div>
                                </div>
                                <div v-else :class="`absolute inset-0 ${item.bgGradient}`"></div>
                                <div class="relative z-10 flex flex-col justify-end h-full px-8 pb-10 md:px-12 md:pb-12 text-white">
                                    <div class="transform transition-transform duration-500 translate-y-4 group-hover:translate-y-0">
                                        <span class="bg-white/20 w-fit px-3 py-1 rounded-full text-xs font-bold mb-3 backdrop-blur-md border border-white/20 shadow-sm inline-block">Latest News</span>
                                        <h2 class="text-2xl md:text-4xl font-bold truncate drop-shadow-md mb-2">{{ item.title }}</h2>
                                        <p class="text-sm md:text-base opacity-90 line-clamp-2 drop-shadow-sm max-w-2xl text-gray-200">{{ item.desc }}</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div v-if="filteredNewsList.length === 0" class="flex flex-col items-center justify-center py-20 text-gray-400">
                    <div class="w-20 h-20 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                        <Heart v-if="showFavoritesOnly" :size="40" class="text-gray-300" />
                        <Search v-else :size="40" class="text-gray-300" />
                    </div>
                    <p class="text-lg font-medium">{{ showFavoritesOnly ? 'ยังไม่มีข่าวที่ถูกใจ' : 'ไม่พบข้อมูลทุน' }}</p>
                </div>

                <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 xl:grid-cols-3 gap-8">
                    <div 
                        v-for="news in filteredNewsList" :key="news.id" @click="goToNews(news.id)"
                        class="group bg-white rounded-2xl p-4 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-300 border border-gray-100 flex flex-col cursor-pointer relative overflow-hidden"
                    >
                        <div :class="`relative w-full aspect-[4/3] rounded-xl overflow-hidden mb-4 flex items-center justify-center group-hover:opacity-95 transition bg-gray-50`">
                            <img v-if="news.imagePath" :src="news.imagePath" class="w-full h-full object-cover" @error="news.imagePath = null" />
                            <div v-else :class="`w-full h-full flex flex-col items-center justify-center text-center p-4 ${news.fallbackColor}`">
                                <h3 class="text-xl font-bold text-slate-700/80 line-clamp-2">{{ news.title }}</h3>
                                <p class="text-sm text-slate-500 mt-1 line-clamp-2">{{ news.desc }}</p>
                            </div>
                            <div v-if="news.isNewBadge" class="absolute top-3 left-3 bg-red-500 text-white px-3 py-1 rounded-lg text-xs font-bold shadow-md animate-pulse">NEW!</div>
                        </div>
                        
                        <div class="flex justify-between items-center mt-auto px-1 pt-1">
                            <div class="flex-1 min-w-0 pr-3">
                                <p class="text-sm font-bold text-gray-700 truncate">
                                    {{ news.title }}
                                </p>
                                <p class="text-xs font-medium text-gray-400 truncate flex items-center gap-1">
                                    {{ news.postDetail }} - {{ news.caption }}
                                </p>
                            </div>
                            
                            <button @click.stop="handleToggleLike(news.id)" class="w-10 h-10 rounded-full flex items-center justify-center transition-all duration-200 active:scale-90 bg-gray-50 hover:bg-red-50 flex-shrink-0">
                                <Heart :class="news.isLiked ? 'fill-[#8B0025] text-[#8B0025]' : 'text-gray-400'" :size="22" :stroke-width="news.isLiked ? 0 : 2" />
                            </button>
                        </div>

                    </div>
                </div>
            </div>
          </div>
      </div>
      
    </main>
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Prompt:wght@300;400;500;600;700&display=swap');
.font-prompt { font-family: 'Prompt', sans-serif; }

.custom-scrollbar::-webkit-scrollbar { width: 6px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background-color: #cbd5e1; border-radius: 20px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background-color: #94a3b8; }

@keyframes bell-ring {
  0%, 100% { transform: rotate(0deg); }
  20% { transform: rotate(-25deg); }
  40% { transform: rotate(25deg); }
  60% { transform: rotate(-15deg); }
  80% { transform: rotate(15deg); }
}
.animate-bell-ring {
  animation: bell-ring 1s infinite ease-in-out;
  transform-origin: 50% 0%; 
}
.animate-fade-in-up { animation: fade-in-up 0.2s ease-out forwards; }
@keyframes fade-in-up { from { opacity: 0; transform: translateY(-10px) scale(0.95); } to { opacity: 1; transform: translateY(0) scale(1); } }
</style>