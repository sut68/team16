<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import type { PropType } from 'vue';
  import type { NewsPost } from '@/interfaces/news_post';
  import NewsCard from './NewsCard.vue';
  import { getAllNewsPosts } from '@/services/api/news_post';

  const props = defineProps({
    // Optional: pass news directly (for embedding in other pages)
    news: {
      type: Array as PropType<NewsPost[]>,
      default: null,
    },
    // Limit number of news to display
    limit: {
      type: Number as PropType<number>,
      default: 6,
    },
    // Show only public news (status_news_id === 1)
    onlyPublic: {
      type: Boolean as PropType<boolean>,
      default: true,
    },
    // Title for the section
    title: {
      type: String as PropType<string>,
      default: 'ข่าวประชาสัมพันธ์ล่าสุด',
    },
    // Show view all button
    showViewAll: {
      type: Boolean as PropType<boolean>,
      default: true,
    },
  });

  const emit = defineEmits<{
    (e: 'view-detail', news: NewsPost): void;
    (e: 'view-all'): void;
  }>();

  // State
  const localNews = ref<NewsPost[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // Use passed news or load from API
  const displayNews = computed(() => {
    let list = props.news ?? localNews.value;
    
    // Filter only public if needed (assuming status_news_id 1 is Public)
    if (props.onlyPublic) {
      list = list.filter(s => s.status_news_id === 1);
    }
    
    // Sort by CreatedAt desc (optional, but good for "latest" news)
    // Assuming CreatedAt is string ISO date
    list = [...list].sort((a, b) => new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime());

    // Limit
    return list.slice(0, props.limit);
  });

  // โหลดข้อมูลจาก API
  async function loadNews() {
    if (props.news) return;
    
    loading.value = true;
    error.value = null;
    
    try {
      // Use requireAuth = false for public access
      const res = await getAllNewsPosts(false);
      localNews.value = res ?? [];
      
    } catch (err: any) {
      console.error('Failed to load news:', err);
      error.value = err?.message ?? 'ไม่สามารถโหลดข้อมูลข่าวได้';
      localNews.value = [];
      
    } finally {
      loading.value = false;
    }
  }

  function handleViewDetail(news: NewsPost) {
    emit('view-detail', news);
  }

  function handleViewAll() {
    emit('view-all');
  }

  onMounted(() => {
    loadNews();
  });
</script>

<template>
  <section class="news-list-section">
    <!-- Section Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <div class="w-1.5 h-8 bg-gradient-to-b from-[#1e3a8a] to-blue-400 rounded-full"></div>
        <h2 class="text-2xl font-bold text-gray-800">{{ title }}</h2>
      </div>
      
      <button 
        v-if="showViewAll && displayNews.length > 0"
        @click="handleViewAll"
        class="flex items-center gap-2 text-[#1e3a8a] hover:text-[#152c6f] 
               font-medium transition-colors duration-200"
      >
        <span>ดูทั้งหมด</span>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
        </svg>
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="i in 3" 
        :key="i" 
        class="bg-white rounded-2xl shadow-md border border-gray-100 overflow-hidden animate-pulse"
      >
        <div class="h-40 bg-gray-200"></div>
        <div class="p-5 space-y-3">
          <div class="h-6 bg-gray-200 rounded w-3/4"></div>
          <div class="h-4 bg-gray-200 rounded w-full"></div>
          <div class="h-4 bg-gray-200 rounded w-2/3"></div>
          <div class="h-10 bg-gray-200 rounded w-full mt-4"></div>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div 
      v-else-if="error" 
      class="bg-red-50 border border-red-200 rounded-xl p-6 text-center"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-red-400 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <p class="text-red-600 font-medium">{{ error }}</p>
      <button 
        @click="loadNews" 
        class="mt-3 text-sm text-red-500 hover:text-red-700 underline"
      >
        ลองใหม่อีกครั้ง
      </button>
    </div>

    <!-- Empty State -->
    <div 
      v-else-if="displayNews.length === 0" 
      class="bg-gray-50 border border-gray-200 rounded-xl p-10 text-center"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-gray-300 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
      </svg>
      <p class="text-gray-500 text-lg">ไม่มีข่าวประชาสัมพันธ์ในขณะนี้</p>
    </div>

    <!-- News Cards Grid -->
    <div 
      v-else 
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      <NewsCard
        v-for="item in displayNews"
        :key="item.ID"
        :news="item"
        @view-detail="handleViewDetail"
      />
    </div>
  </section>
</template>

<style scoped>
.news-list-section {
  width: 100%;
}
</style>
