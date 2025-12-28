<script setup lang="ts">
import { ref, onMounted } from 'vue';

// Stats data - can be fetched from API later
const stats = ref([
  {
    id: 'scholarships',
    value: 15,
    label: 'ทุนที่เปิดรับ',
    suffix: 'ทุน',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path d="M12 14l9-5-9-5-9 5 9 5z" /><path d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 14l9-5-9-5-9 5 9 5zm0 0l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14zm-4 6v-7.5l4-2.222" /></svg>`,
    color: 'from-orange-500 to-amber-500',
    bgColor: 'bg-orange-50',
    textColor: 'text-orange-600',
  },
  {
    id: 'recipients',
    value: 500,
    label: 'ผู้รับทุนสะสม',
    suffix: 'คน',
    prefix: '+',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" /></svg>`,
    color: 'from-blue-500 to-indigo-500',
    bgColor: 'bg-blue-50',
    textColor: 'text-blue-600',
  },
  {
    id: 'totalValue',
    value: 5.2,
    label: 'มูลค่าทุนรวม',
    suffix: 'ล้านบาท',
    prefix: '',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>`,
    color: 'from-emerald-500 to-teal-500',
    bgColor: 'bg-emerald-50',
    textColor: 'text-emerald-600',
  },
]);

// Animated counter
const displayValues = ref<number[]>([]);

function animateCounters() {
  stats.value.forEach((stat, index) => {
    const target = stat.value;
    const duration = 2000; // 2 seconds
    const steps = 60;
    const increment = target / steps;
    let current = 0;
    
    displayValues.value[index] = 0;
    
    const timer = setInterval(() => {
      current += increment;
      if (current >= target) {
        displayValues.value[index] = target;
        clearInterval(timer);
      } else {
        displayValues.value[index] = Math.floor(current * 10) / 10;
      }
    }, duration / steps);
  });
}

onMounted(() => {
  // Initialize display values
  displayValues.value = stats.value.map(() => 0);
  
  // Start animation after a short delay
  setTimeout(animateCounters, 300);
});

function formatNumber(value: number, hasDecimal: boolean = false): string {
  if (hasDecimal) {
    return value.toFixed(1);
  }
  return Math.floor(value).toLocaleString('th-TH');
}
</script>

<template>
  <section class="py-12">
    <div class="container mx-auto px-4">
      <!-- Section Title -->
      <div class="text-center mb-10">
        <h2 class="text-2xl md:text-3xl font-bold text-gray-800 mb-2">
          สถิติทุนการศึกษา <span class="text-[#F26522]">มทส.</span>
        </h2>
        <p class="text-gray-500">ข้อมูลทุนการศึกษาล่าสุด</p>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div 
          v-for="(stat, index) in stats" 
          :key="stat.id"
          class="stat-card group relative bg-white rounded-2xl p-6 shadow-md hover:shadow-xl 
                 border border-gray-100 transition-all duration-300 hover:-translate-y-1 overflow-hidden"
        >
          <!-- Background gradient on hover -->
          <div 
            class="absolute inset-0 opacity-0 group-hover:opacity-5 transition-opacity duration-300 bg-gradient-to-br"
            :class="stat.color"
          ></div>

          <!-- Icon -->
          <div 
            class="w-14 h-14 rounded-xl flex items-center justify-center mb-4 transition-transform duration-300 group-hover:scale-110"
            :class="[stat.bgColor, stat.textColor]"
          >
            <span v-html="stat.icon"></span>
          </div>

          <!-- Value -->
          <div class="flex items-baseline gap-1 mb-1">
            <span v-if="stat.prefix" class="text-2xl font-bold text-gray-800">{{ stat.prefix }}</span>
            <span class="text-4xl font-bold text-gray-800">
              {{ stat.id === 'totalValue' ? formatNumber(displayValues[index] || 0, true) : formatNumber(displayValues[index] || 0) }}
            </span>
            <span class="text-lg font-medium text-gray-500 ml-1">{{ stat.suffix }}</span>
          </div>

          <!-- Label -->
          <p class="text-gray-600 font-medium">{{ stat.label }}</p>

          <!-- Decorative corner -->
          <div 
            class="absolute -bottom-4 -right-4 w-20 h-20 rounded-full opacity-10"
            :class="stat.bgColor"
          ></div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.stat-card {
  position: relative;
}
</style>
