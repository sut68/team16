<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import type { PropType } from 'vue';
  import type { ScholarshipResponse } from '@/interfaces/scholarship';
  import ScholarshipCard from './ScholarshipCard.vue';
  import { ScholarshipAPI } from '@/services/api/scholarship';

  const props = defineProps({
    // Optional: pass scholarships directly (for embedding in other pages)
    scholarships: {
      type: Array as PropType<ScholarshipResponse[]>,
      default: null,
    },
    // Limit number of scholarships to display
    limit: {
      type: Number as PropType<number>,
      default: 6,
    },
    // Show only open scholarships
    onlyOpen: {
      type: Boolean as PropType<boolean>,
      default: true,
    },
    // Title for the section
    title: {
      type: String as PropType<string>,
      default: 'ทุนการศึกษาที่เปิดรับสมัคร',
    },
    // Show view all button
    showViewAll: {
      type: Boolean as PropType<boolean>,
      default: true,
    },
  });

  const emit = defineEmits<{
    (e: 'view-detail', scholarship: ScholarshipResponse): void;
    (e: 'view-all'): void;
  }>();

  // State
  const localScholarships = ref<ScholarshipResponse[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  // Use passed scholarships or load from API
  const displayScholarships = computed(() => {
    let list = props.scholarships ?? localScholarships.value;
    
    // Filter only open if needed
    if (props.onlyOpen) {
      list = list.filter(s => s.statusscholarship?.status_name?.toLowerCase() === 'open');
    }
    
    // Limit
    return list.slice(0, props.limit);
  });

  // โหลดข้อมูลจาก API
  async function loadScholarships() {
    if (props.scholarships) return;
    
    loading.value = true;
    error.value = null;
    
    try {
      const res = await ScholarshipAPI.getAll();
      localScholarships.value = res ?? [];
      
    } catch (err: any) {
      console.error('Failed to load scholarships:', err);
      error.value = err?.message ?? 'ไม่สามารถโหลดข้อมูลทุนได้';
      localScholarships.value = [];
      
    } finally {
      loading.value = false;
    }
  }

  function handleViewDetail(scholarship: ScholarshipResponse) {
    emit('view-detail', scholarship);
  }

  function handleViewAll() {
    emit('view-all');
  }

  onMounted(() => {
    loadScholarships();
  });
</script>

<template>
  <section class="scholarship-list-section">
    <!-- Section Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <div class="w-1.5 h-8 bg-gradient-to-b from-[#1e3a8a] to-blue-400 rounded-full"></div>
        <h2 class="text-2xl font-bold text-gray-800">{{ title }}</h2>
      </div>
      
      <button 
        v-if="showViewAll && displayScholarships.length > 0"
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
          <div class="h-4 bg-gray-200 rounded w-1/2"></div>
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
        @click="loadScholarships" 
        class="mt-3 text-sm text-red-500 hover:text-red-700 underline"
      >
        ลองใหม่อีกครั้ง
      </button>
    </div>

    <!-- Empty State -->
    <div 
      v-else-if="displayScholarships.length === 0" 
      class="bg-gray-50 border border-gray-200 rounded-xl p-10 text-center"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-gray-300 mx-auto mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path d="M12 14l9-5-9-5-9 5 9 5z" />
        <path d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 14l9-5-9-5-9 5 9 5zm0 0l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14zm-4 6v-7.5l4-2.222" />
      </svg>
      <p class="text-gray-500 text-lg">ไม่มีทุนการศึกษาที่เปิดรับสมัครในขณะนี้</p>
    </div>

    <!-- Scholarship Cards Grid -->
    <div 
      v-else 
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      <ScholarshipCard
        v-for="scholarship in displayScholarships"
        :key="scholarship.ID"
        :scholarship="scholarship"
        @view-detail="handleViewDetail"
      />
    </div>
  </section>
</template>

<style scoped>
.scholarship-list-section {
  width: 100%;
}
</style>
