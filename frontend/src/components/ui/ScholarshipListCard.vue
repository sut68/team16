<script setup lang="ts">
import { ref, computed } from 'vue'
import { Building2, Calendar, FolderOpen, CalendarX } from 'lucide-vue-next'
import ScholarshipActionMenu from './ScholarshipActionMenu.vue'
import type { ScholarshipResponse } from '@/interfaces'

// Props
const props = defineProps<{
  scholarship: ScholarshipResponse
}>()

// State for menu open tracking
const isMenuOpen = ref(false)

// Emits
const emit = defineEmits<{
  view: [scholarship: ScholarshipResponse]
  edit: [scholarship: ScholarshipResponse]
  manageFeatures: [scholarshipId: number]
  delete: [scholarship: ScholarshipResponse]
}>()

// Computed
const companyInitials = computed(() => {
  const name = props.scholarship.sponsor?.company_name || 'N/A'
  return name.substring(0, 2).toUpperCase()
})

const statusInfo = computed(() => {
  const status = props.scholarship.statusscholarship?.status_name?.toLowerCase() || ''
  
  if (status.includes('open') || status === 'เปิดรับสมัคร') {
    return {
      label: 'เปิดรับสมัคร',
      bg: 'bg-emerald-100',
      text: 'text-emerald-700',
      dot: 'bg-emerald-500',
      pulse: true
    }
  } else if (status.includes('closed') || status === 'ปิดรับสมัคร') {
    return {
      label: 'ปิดรับสมัคร',
      bg: 'bg-red-100',
      text: 'text-red-700',
      dot: 'bg-red-500',
      pulse: false
    }
  } else if (status.includes('pending') || status === 'รอเปิดรับ') {
    return {
      label: 'รอเปิดรับ',
      bg: 'bg-amber-100',
      text: 'text-amber-700',
      dot: 'bg-amber-500',
      pulse: false
    }
  }
  
  return {
    label: props.scholarship.statusscholarship?.status_name || 'ไม่ระบุ',
    bg: 'bg-gray-100',
    text: 'text-gray-700',
    dot: 'bg-gray-500',
    pulse: false
  }
})

const typeInfo = computed(() => {
  const typeName = props.scholarship.typescholarship?.type_name || 'ไม่ระบุประเภท'
  return typeName
})

// Format date helper
function formatDate(dateStr: string | undefined): string {
  if (!dateStr) return '-'
  try {
    const date = new Date(dateStr)
    return date.toLocaleDateString('th-TH', {
      day: 'numeric',
      month: 'short',
      year: '2-digit'
    })
  } catch {
    return dateStr
  }
}

// Handlers
function handleCardClick() {
  emit('view', props.scholarship)
}

function handleView(id: number) {
  emit('view', props.scholarship)
}

function handleEdit(id: number) {
  emit('edit', props.scholarship)
}

function handleManageFeatures(id: number) {
  emit('manageFeatures', id)
}

function handleDelete(id: number) {
  emit('delete', props.scholarship)
}
</script>

<template>
  <div
    @click="handleCardClick"
    class="scholarship-list-card relative bg-white rounded-2xl border border-gray-200 shadow-sm hover:shadow-md hover:border-blue-200 transition-all duration-200 p-5 cursor-pointer group"
    :class="{ 'z-50': isMenuOpen, 'z-0': !isMenuOpen }"
  >
    <!-- Row 1: Header -->
    <div class="flex items-start justify-between gap-4">
      <!-- Left: Avatar + Title -->
      <div class="flex items-start gap-4 flex-1 min-w-0">
        <!-- Avatar -->
        <div 
          class="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white font-bold text-sm shrink-0 shadow-sm"
        >
          {{ companyInitials }}
        </div>

        <!-- Title & Company -->
        <div class="min-w-0 flex-1">
          <h3 class="text-lg font-bold text-gray-900 truncate group-hover:text-blue-700 transition-colors">
            {{ scholarship.scholarship_name }}
          </h3>
          <p class="text-sm text-gray-500 flex items-center gap-1.5 mt-0.5">
            <Building2 class="w-3.5 h-3.5" />
            <span class="truncate">{{ scholarship.sponsor?.company_name || 'ไม่ระบุบริษัท' }}</span>
          </p>
        </div>
      </div>

      <!-- Right: Status Badge -->
      <div class="flex items-center gap-2 shrink-0">
        <span 
          class="inline-flex items-center gap-1.5 px-3 py-1 text-xs font-semibold rounded-full"
          :class="[statusInfo.bg, statusInfo.text]"
        >
          <span 
            class="w-1.5 h-1.5 rounded-full"
            :class="[statusInfo.dot, { 'animate-pulse': statusInfo.pulse }]"
          ></span>
          {{ statusInfo.label }}
        </span>
      </div>
    </div>

    <!-- Row 2: Description -->
    <p class="text-sm text-gray-600 mt-4 line-clamp-2 leading-relaxed">
      {{ scholarship.description || 'ไม่มีรายละเอียด' }}
    </p>

    <!-- Divider -->
    <div class="border-t border-dashed border-gray-200 my-4"></div>

    <!-- Row 3: Meta Info & Actions -->
    <div class="flex items-center justify-between">
      <!-- Info Pills -->
      <div class="flex items-center gap-2 flex-wrap">
        <!-- Type -->
        <span class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-blue-50 text-blue-700 rounded-full text-xs font-medium">
          <FolderOpen class="w-3.5 h-3.5" />
          {{ typeInfo }}
        </span>

        <!-- Open Date -->
        <span class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-green-50 text-green-700 rounded-full text-xs font-medium">
          <Calendar class="w-3.5 h-3.5" />
          เปิด: {{ formatDate(scholarship.open_date) }}
        </span>

        <!-- Close Date -->
        <span class="inline-flex items-center gap-1.5 px-3 py-1.5 bg-red-50 text-red-700 rounded-full text-xs font-medium">
          <CalendarX class="w-3.5 h-3.5" />
          ปิด: {{ formatDate(scholarship.close_date) }}
        </span>
      </div>

      <!-- Action Menu -->
      <ScholarshipActionMenu
        :scholarship-id="scholarship.ID"
        :scholarship-name="scholarship.scholarship_name"
        @view="handleView"
        @edit="handleEdit"
        @manage-features="handleManageFeatures"
        @delete="handleDelete"
        @menu-open="isMenuOpen = true"
        @menu-close="isMenuOpen = false"
      />
    </div>
  </div>
</template>

<style scoped>
.scholarship-list-card {
  will-change: transform, box-shadow;
}

.scholarship-list-card:hover {
  transform: translateY(-2px);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
