<script setup lang="ts">
  import { computed } from 'vue';
  import type { PropType } from 'vue';
  import type { ScholarshipResponse } from '@/interfaces/scholarship';
  import { formatDateThai } from '@/utils/dateFormatter';

  const props = defineProps({
    scholarship: {
      type: Object as PropType<ScholarshipResponse>,
      required: true,
    },
    showButton: {
      type: Boolean as PropType<boolean>,
      default: true,
    },
  });

  const emit = defineEmits<{
    (e: 'view-detail', scholarship: ScholarshipResponse): void;
  }>();

  // ไอคอนประเภททุน (typeIcon)
  const typeIcon = computed(() => {
    const type = props.scholarship?.typescholarship?.type_name?.toLowerCase();
    if (type === 'full') {
      return `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" /></svg>`;
    }
    return `<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>`;
  });

  // สีพื้นหลัง (gradientClass)
  const gradientClass = computed(() => {
    const type = props.scholarship?.typescholarship?.type_name?.toLowerCase();
    if (type === 'full') {
      return 'from-amber-400 via-orange-400 to-orange-500';
    }
    return 'from-blue-400 via-indigo-400 to-purple-500';
  });

  // ตรวจสอบว่าทุนเปิดรับหรือไม่
  const isOpen = computed(() => {
    return props.scholarship?.statusscholarship?.status_name?.toLowerCase() === 'open';
  });

  // สีชิ้นส่วนสถานะ (statusBadgeClass)
  const statusBadgeClass = computed(() => {
    if (isOpen.value) {
      return 'bg-emerald-100 text-emerald-700 border-emerald-200';
    }
    return 'bg-gray-100 text-gray-600 border-gray-200';
  });

  // สีชิ้นส่วนประเภท (typeBadgeClass)
  const typeBadgeClass = computed(() => {
    const type = props.scholarship?.typescholarship?.type_name?.toLowerCase();
    if (type === 'full') {
      return 'bg-amber-100 text-amber-700 border-amber-200';
    }
    return 'bg-indigo-100 text-indigo-700 border-indigo-200';
  });

  function handleViewDetail() {
    emit('view-detail', props.scholarship);
  }
</script>

<template>
  <div 
    class="scholarship-card group bg-white rounded-2xl shadow-md hover:shadow-xl 
           border border-gray-100 overflow-hidden transition-all duration-300 
           hover:-translate-y-1 flex flex-col"
  >
    <!-- Header with gradient background -->
    <div 
      class="relative h-72 bg-gradient-to-br overflow-hidden"
      :class="gradientClass"
    >
      <!-- Decorative circles -->
      <div class="absolute -top-10 -right-10 w-32 h-32 bg-white/10 rounded-full"></div>
      <div class="absolute -bottom-10 -left-10 w-40 h-40 bg-white/10 rounded-full"></div>
      <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-20 h-20 bg-white/10 rounded-full"></div>
      
      <!-- Icon in center -->
      <div class="absolute inset-0 flex items-center justify-center">
        <div class="w-20 h-20 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center shadow-lg">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path d="M12 14l9-5-9-5-9 5 9 5z" />
            <path d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 14l9-5-9-5-9 5 9 5zm0 0l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14zm-4 6v-7.5l4-2.222" />
          </svg>
        </div>
      </div>

      <!-- Status badge -->
      <div class="absolute top-3 left-3">
        <span 
          class="px-3 py-1 text-xs font-semibold rounded-full border"
          :class="statusBadgeClass"
        >
          {{ scholarship.statusscholarship?.status_name || 'Unknown' }}
        </span>
      </div>

      <!-- Type badge -->
      <div class="absolute top-3 right-3">
        <span 
          class="px-3 py-1 text-xs font-semibold rounded-full border flex items-center gap-1"
          :class="typeBadgeClass"
        >
          <span v-html="typeIcon" class="w-4 h-4"></span>
          {{ scholarship.typescholarship?.type_name || 'Unknown' }}
        </span>
      </div>
    </div>

    <!-- Content -->
    <div class="p-5 flex-1 flex flex-col">
      <!-- Title -->
      <h3 class="text-base font-bold text-gray-800 mb-2 line-clamp-2 group-hover:text-[#1e3a8a] transition-colors">
        {{ scholarship.scholarship_name }}
      </h3>

      <!-- Description -->
      <div class="flex items-start gap-2 mb-3 text-gray-600 text-sm">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-orange-500 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
        </svg>
        <span class="line-clamp-2">{{ scholarship.sponsor?.company_name || 'ไม่ระบุผู้สนับสนุน' }}</span>
      </div>

      <!-- Date range -->
      <div class="flex items-center gap-2 text-gray-500 text-sm mb-2">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-[#1e3a8a] flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <span>{{ formatDateThai(scholarship.open_date) }} - {{ formatDateThai(scholarship.close_date) }}</span>
      </div>

      <!-- Spacer to push button to bottom -->
      <div class="flex-1"></div>

      <!-- View Detail Button - slides up on hover -->
      <button 
        v-if="showButton"
        @click="handleViewDetail"
        class="w-full py-3 px-4 bg-[#800020] text-white font-medium rounded-xl
               hover:bg-[#600018] shadow-md hover:shadow-lg
               flex items-center justify-center gap-2
               transition-all duration-300 ease-out
               opacity-0 translate-y-4 group-hover:opacity-100 group-hover:translate-y-0"
        :disabled="!isOpen"
        :class="{ '!opacity-50 cursor-not-allowed !translate-y-0': !isOpen }"
      >
        <span>ดูรายละเอียด</span>
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.scholarship-card {
  min-height: 420px;
}

/* Orange shadow on hover */
.scholarship-card:hover {
  box-shadow: 
    0 20px 40px -10px rgba(242, 101, 34, 0.25),
    0 10px 20px -5px rgba(242, 101, 34, 0.15);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
