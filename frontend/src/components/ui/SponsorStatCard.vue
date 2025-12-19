<script setup lang="ts">
defineProps<{
  label: string
  value: number | string
  icon?: string // 'sponsors' | 'active' | 'contacts'
  trend?: number // เช่น +5 หรือ -3 (optional)
  color?: 'purple' | 'green' | 'blue' | 'orange'
}>()

const colorMap = {
  purple: {
    gradient: 'from-purple-500 to-indigo-600',
    bg: 'bg-purple-50',
    text: 'text-purple-600',
    iconBg: 'bg-purple-100',
    ring: 'ring-purple-200',
  },
  green: {
    gradient: 'from-emerald-500 to-teal-600',
    bg: 'bg-emerald-50',
    text: 'text-emerald-600',
    iconBg: 'bg-emerald-100',
    ring: 'ring-emerald-200',
  },
  blue: {
    gradient: 'from-blue-500 to-cyan-600',
    bg: 'bg-blue-50',
    text: 'text-blue-600',
    iconBg: 'bg-blue-100',
    ring: 'ring-blue-200',
  },
  orange: {
    gradient: 'from-orange-500 to-amber-600',
    bg: 'bg-orange-50',
    text: 'text-orange-600',
    iconBg: 'bg-orange-100',
    ring: 'ring-orange-200',
  },
}
</script>

<template>
  <div 
    class="stat-card group relative overflow-hidden rounded-2xl bg-white p-5
           border border-gray-100 shadow-sm
           hover:shadow-lg hover:shadow-gray-200/50
           transition-all duration-300 ease-out
           hover:-translate-y-1"
  >
    <!-- Gradient accent bar -->
    <div 
      class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r opacity-80"
      :class="colorMap[color ?? 'purple'].gradient"
    ></div>

    <!-- Content -->
    <div class="flex items-start justify-between">
      <div class="space-y-2">
        <!-- Label -->
        <p class="text-sm font-medium text-gray-500 tracking-wide">
          {{ label }}
        </p>
        
        <!-- Value -->
        <div class="flex items-baseline gap-2">
          <span class="text-3xl font-bold text-gray-900 tabular-nums">
            {{ typeof value === 'number' ? value.toLocaleString() : value }}
          </span>
          
          <!-- Trend indicator (optional) -->
          <span 
            v-if="trend !== undefined"
            class="inline-flex items-center gap-0.5 text-xs font-medium px-1.5 py-0.5 rounded-full"
            :class="trend >= 0 ? 'text-emerald-700 bg-emerald-100' : 'text-red-700 bg-red-100'"
          >
            <svg 
              class="w-3 h-3" 
              :class="trend >= 0 ? '' : 'rotate-180'"
              fill="none" viewBox="0 0 24 24" stroke="currentColor"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18"/>
            </svg>
            {{ Math.abs(trend) }}%
          </span>
        </div>
      </div>

      <!-- Icon -->
      <div 
        class="flex items-center justify-center w-12 h-12 rounded-xl
               transition-transform duration-300 group-hover:scale-110"
        :class="colorMap[color ?? 'purple'].iconBg"
      >
        <!-- Sponsors icon -->
        <svg v-if="icon === 'sponsors'" class="w-6 h-6" :class="colorMap[color ?? 'purple'].text" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
        </svg>
        
        <!-- Active icon -->
        <svg v-else-if="icon === 'active'" class="w-6 h-6" :class="colorMap[color ?? 'purple'].text" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        
        <!-- Contacts icon -->
        <svg v-else-if="icon === 'contacts'" class="w-6 h-6" :class="colorMap[color ?? 'purple'].text" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
        </svg>
        
        <!-- Default icon -->
        <svg v-else class="w-6 h-6" :class="colorMap[color ?? 'purple'].text" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
        </svg>
      </div>
    </div>

    <!-- Decorative blur circle -->
    <div 
      class="absolute -bottom-8 -right-8 w-24 h-24 rounded-full opacity-10 blur-2xl
             transition-opacity duration-300 group-hover:opacity-20"
      :class="colorMap[color ?? 'purple'].gradient.replace('from-', 'bg-').split(' ')[0]"
    ></div>
  </div>
</template>

<style scoped>
.stat-card {
  backdrop-filter: blur(10px);
}
</style>