<script setup lang="ts">
  import { 
    Building2, CheckCircle, Users, BarChart3,
    ClipboardList, Clock, RefreshCw, CheckCircle2
  } from 'lucide-vue-next'

  defineProps<{
    label: string
    value: number | string
    icon?: string
    trend?: number
    color?: 'purple' | 'green' | 'blue' | 'orange' | 'slate' | 'amber'
  }>()

  const colorMap = {
    purple: {
      gradient: 'from-purple-500 to-indigo-600',
      bgGradient: 'from-white via-purple-50/50 to-indigo-100/60',
      bg: 'bg-purple-50',
      text: 'text-purple-600',
      iconBg: 'bg-purple-100',
      ring: 'ring-purple-200',
    },
    green: {
      gradient: 'from-emerald-500 to-teal-600',
      bgGradient: 'from-white via-emerald-50/50 to-teal-100/60',
      bg: 'bg-emerald-50',
      text: 'text-emerald-600',
      iconBg: 'bg-emerald-100',
      ring: 'ring-emerald-200',
    },
    blue: {
      gradient: 'from-blue-500 to-cyan-600',
      bgGradient: 'from-white via-blue-50/50 to-cyan-100/60',
      bg: 'bg-blue-50',
      text: 'text-blue-600',
      iconBg: 'bg-blue-100',
      ring: 'ring-blue-200',
    },
    orange: {
      gradient: 'from-orange-500 to-amber-600',
      bgGradient: 'from-white via-orange-50/50 to-amber-100/60',
      bg: 'bg-orange-50',
      text: 'text-orange-600',
      iconBg: 'bg-orange-100',
      ring: 'ring-orange-200',
    },
    slate: {
      gradient: 'from-slate-500 to-slate-700',
      bgGradient: 'from-white via-slate-50/50 to-slate-100/60',
      bg: 'bg-slate-50',
      text: 'text-slate-600',
      iconBg: 'bg-slate-100',
      ring: 'ring-slate-200',
    },
    amber: {
      gradient: 'from-amber-500 to-orange-600',
      bgGradient: 'from-white via-amber-50/50 to-orange-100/60',
      bg: 'bg-amber-50',
      text: 'text-amber-600',
      iconBg: 'bg-amber-100',
      ring: 'ring-amber-200',
    },
  }

  const iconMap: Record<string, any> = {
    // Sponsor icons
    sponsors: Building2,
    active: CheckCircle,
    contacts: Users,
    default: BarChart3,
    // Evaluation icons
    total: ClipboardList,
    pending: Clock,
    inProgress: RefreshCw,
    completed: CheckCircle2,
  }
</script>

<template>
  <div 
    class="stat-card group relative overflow-hidden rounded-2xl p-5
           border border-gray-100 shadow-sm
           hover:shadow-lg hover:shadow-gray-200/50
           transition-all duration-300 ease-out
           hover:-translate-y-1
           bg-gradient-to-br"
    :class="colorMap[color ?? 'purple'].bgGradient"
  >
    <!-- Gradient accent bar -->
    <div 
      class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r opacity-90"
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
        <component 
          :is="iconMap[icon ?? 'default'] ?? iconMap.default" 
          class="w-6 h-6" 
          :class="colorMap[color ?? 'purple'].text" 
        />
      </div>
    </div>

    <!-- Decorative blur circle -->
    <div 
      class="absolute -bottom-8 -right-8 w-24 h-24 rounded-full opacity-15 blur-2xl
             transition-opacity duration-300 group-hover:opacity-25"
      :class="colorMap[color ?? 'purple'].gradient.replace('from-', 'bg-').split(' ')[0]"
    ></div>
  </div>
</template>

<style scoped>
.stat-card {
  backdrop-filter: blur(10px);
}
</style>
