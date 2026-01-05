<script setup lang="ts">
  import { computed } from 'vue'

  // --- Types ---
  export interface StatItem {
    title: string
    value: number | string
    description?: string
    icon?: IconType
    color?: ColorType
  }

  export type ColorType = 'blue' | 'green' | 'orange' | 'purple' | 'slate' | 'amber' | 'red'
  export type IconType = 'box' | 'users' | 'clock' | 'check' | 'chart' | 'calendar' | 'building' | 'clipboard' | 'star' | 'award'

  // --- Props ---
  const props = withDefaults(defineProps<{
    stats: StatItem[]
    columns?: number | 'auto'  // Number of columns or 'auto' for responsive
  }>(), {
    columns: 'auto'
  })

  // --- Color Configuration ---
  const colorMap: Record<ColorType, { value: string; bg: string; iconBg: string }> = {
    blue: {
      value: 'text-[#1e3a8a]',
      bg: 'bg-blue-50',
      iconBg: 'text-[#1e3a8a]'
    },
    green: {
      value: 'text-emerald-700',
      bg: 'bg-green-50',
      iconBg: 'text-emerald-700'
    },
    orange: {
      value: 'text-orange-500',
      bg: 'bg-orange-50',
      iconBg: 'text-orange-500'
    },
    purple: {
      value: 'text-purple-700',
      bg: 'bg-purple-50',
      iconBg: 'text-purple-700'
    },
    slate: {
      value: 'text-slate-700',
      bg: 'bg-slate-50',
      iconBg: 'text-slate-700'
    },
    amber: {
      value: 'text-amber-600',
      bg: 'bg-amber-50',
      iconBg: 'text-amber-600'
    },
    red: {
      value: 'text-red-600',
      bg: 'bg-red-50',
      iconBg: 'text-red-600'
    }
  }

  // --- Computed ---
  const gridClass = computed(() => {
    if (props.columns === 'auto') {
      // Responsive: 1 col on mobile, 2 on md, then based on stats count
      const count = props.stats.length
      if (count <= 2) return `grid-cols-1 md:grid-cols-${count}`
      if (count === 3) return 'grid-cols-1 md:grid-cols-3'
      if (count === 4) return 'grid-cols-1 md:grid-cols-2 xl:grid-cols-4'
      return 'grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4'
    }
    return `grid-cols-1 md:grid-cols-${Math.min(props.columns, props.stats.length)}`
  })

  const getColor = (color?: ColorType) => colorMap[color ?? 'blue']

  // --- Description Text Color Mapping ---
  const descColorMap: Record<ColorType, string> = {
    blue: 'text-blue-600',
    green: 'text-emerald-600',
    orange: 'text-orange-500',
    purple: 'text-purple-600',
    slate: 'text-slate-500',
    amber: 'text-amber-600',
    red: 'text-red-500'
  }

  const getDescColor = (color?: ColorType) => descColorMap[color ?? 'slate']
</script>

<template>
  <div 
    class="grid bg-white shadow rounded-2xl border border-gray-100 w-full divide-y md:divide-y-0 md:divide-x divide-gray-100"
    :class="gridClass"
  >
    <div 
      v-for="(stat, index) in stats" 
      :key="index"
      class="p-4 flex flex-row items-center justify-between"
    >
      <!-- Text Content -->
      <div>
        <div class="text-slate-500 text-sm mb-1">{{ stat.title }}</div>
        <div 
          class="text-3xl font-bold"
          :class="getColor(stat.color).value"
        >
          {{ typeof stat.value === 'number' ? stat.value.toLocaleString() : stat.value }}
        </div>
        <div 
          v-if="stat.description"
          class="text-xs mt-1 font-medium"
          :class="getDescColor(stat.color)"
        >
          {{ stat.description }}
        </div>
      </div>

      <!-- Icon -->
      <div 
        class="p-3 rounded-full"
        :class="getColor(stat.color).bg"
      >
        <!-- Box / Package Icon -->
        <svg v-if="stat.icon === 'box'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
        </svg>

        <!-- Users Icon -->
        <svg v-else-if="stat.icon === 'users'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/>
        </svg>

        <!-- Clock Icon -->
        <svg v-else-if="stat.icon === 'clock'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>

        <!-- Check Icon -->
        <svg v-else-if="stat.icon === 'check'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>

        <!-- Chart Icon -->
        <svg v-else-if="stat.icon === 'chart'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
        </svg>

        <!-- Calendar Icon -->
        <svg v-else-if="stat.icon === 'calendar'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
        </svg>

        <!-- Building Icon -->
        <svg v-else-if="stat.icon === 'building'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
        </svg>

        <!-- Clipboard Icon -->
        <svg v-else-if="stat.icon === 'clipboard'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/>
        </svg>

        <!-- Star Icon -->
        <svg v-else-if="stat.icon === 'star'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"/>
        </svg>

        <!-- Award Icon -->
        <svg v-else-if="stat.icon === 'award'" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/>
        </svg>

        <!-- Default Icon (Chart) -->
        <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" 
          class="inline-block w-8 h-8 stroke-current" :class="getColor(stat.color).iconBg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
        </svg>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Custom grid column classes for dynamic columns */
.grid-cols-2 { grid-template-columns: repeat(2, minmax(0, 1fr)); }
.grid-cols-3 { grid-template-columns: repeat(3, minmax(0, 1fr)); }
.grid-cols-4 { grid-template-columns: repeat(4, minmax(0, 1fr)); }
.grid-cols-5 { grid-template-columns: repeat(5, minmax(0, 1fr)); }
.grid-cols-6 { grid-template-columns: repeat(6, minmax(0, 1fr)); }
</style>
