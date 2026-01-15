<script setup lang="ts">
import { computed } from 'vue';
import type { PropType } from 'vue';
import type { NewsPost } from '@/interfaces/news_post';
import { formatDateThai } from '@/utils/dateFormatter';

const props = defineProps({
  news: {
    type: Object as PropType<NewsPost>,
    required: true,
  },
  showButton: {
    type: Boolean as PropType<boolean>,
    default: true,
  },
});

const emit = defineEmits<{
  (e: 'view-detail', news: NewsPost): void;
}>();

// Helper to get image URL
const imageUrl = computed(() => {
  if (props.news.file_path) {
    if (props.news.file_path.startsWith('http')) {
      return props.news.file_path;
    }
    // อ่าน Base URL จาก ENV (ตัด /api ออกสำหรับ static files)
    const baseUrl = (import.meta.env.VITE_API_URL || 'http://localhost:8080/api').replace('/api', '');
    return `${baseUrl}/${props.news.file_path}`;
  }
  return null;
});


// Gradient background (random or based on id)
const gradientClass = computed(() => {
  // Simple logic to alternate colors based on ID
  if ((props.news.ID || 0) % 2 === 0) {
    return 'from-blue-400 via-indigo-400 to-purple-500';
  }
  return 'from-amber-400 via-orange-400 to-orange-500';
});

function handleViewDetail() {
  emit('view-detail', props.news);
}
</script>

<template>
  <div 
    class="news-card group bg-white rounded-2xl shadow-md hover:shadow-xl 
           border border-gray-100 overflow-hidden transition-all duration-300 
           hover:-translate-y-1 flex flex-col"
  >
    <!-- Header / Image -->
    <div class="relative h-64 bg-gray-200 overflow-hidden group">
        <!-- If we have an image, show it -->
        <img 
            v-if="imageUrl" 
            :src="imageUrl" 
            alt="News Image" 
            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
        />
        
        <!-- Fallback Gradient if no image -->
        <div 
            v-else 
            class="w-full h-full bg-gradient-to-br flex items-center justify-center p-8"
            :class="gradientClass"
        >
             <div class="w-16 h-16 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center shadow-lg">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
                </svg>
            </div>
        </div>

        <!-- Date Badge -->
        <div class="absolute top-3 left-3">
             <span class="px-3 py-1 text-xs font-semibold rounded-full bg-white/90 text-gray-700 shadow-sm backdrop-blur-sm">
                {{ formatDateThai(news.CreatedAt) }}
             </span>
        </div>
    </div>

    <!-- Content -->
    <div class="p-5 flex-1 flex flex-col">
      <!-- Title -->
      <h3 class="text-lg font-bold text-gray-800 mb-2 line-clamp-2 group-hover:text-[#1e3a8a] transition-colors">
        {{ news.title }}
      </h3>

      <!-- Detail (truncated) -->
      <div 
        class="text-gray-600 text-sm mb-4 line-clamp-3 prose prose-sm max-w-none" 
        v-html="news.post_detail"
      ></div>

      <!-- Spacer to push button to bottom -->
      <div class="flex-1"></div>
      
      <!-- Sponsor/Scholarship info if available -->
      <div v-if="news.scholarship?.sponsor?.company_name" class="flex items-center gap-2 mb-3 text-xs text-gray-500">
         <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-orange-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
             <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
         </svg>
         <span class="truncate">{{ news.scholarship.sponsor.company_name }}</span>
      </div>


      <!-- View Detail Button - slides up on hover -->
      <button 
        v-if="showButton"
        @click="handleViewDetail"
        class="w-full py-3 px-4 bg-[#1e3a8a] text-white font-medium rounded-xl
               hover:bg-[#152c6f] shadow-md hover:shadow-lg
               flex items-center justify-center gap-2
               transition-all duration-300 ease-out
               opacity-0 translate-y-4 group-hover:opacity-100 group-hover:translate-y-0"
      >
        <span>อ่านเพิ่มเติม</span>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.news-card {
  min-height: 420px;
}

/* Base Hover Shadow */
.news-card:hover {
  box-shadow: 
    0 20px 40px -10px rgba(30, 58, 138, 0.25),
    0 10px 20px -5px rgba(30, 58, 138, 0.15);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
