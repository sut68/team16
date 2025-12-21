<script setup lang="ts">
import { ChevronsLeft, ChevronLeft, ChevronRight, ChevronsRight } from 'lucide-vue-next'

defineProps<{
  page: number
  totalPages: number
  pages: number[]
  total: number
  perPage: number
  pagedLength: number
}>()

defineEmits<{
  prev: []
  next: []
  first: []
  last: []
  goto: [page: number]
}>()
</script>

<template>
  <div 
    class="flex flex-col md:flex-row items-center justify-between gap-4 pt-2 pb-0 pl-2 pr-2"
    data-testid="pagination-container"
  >
    <!-- Info text -->
    <div class="text-sm text-slate-500" data-testid="pagination-info">
      แสดง
      <span class="font-semibold text-slate-700">{{ (page - 1) * perPage + (pagedLength ? 1 : 0) }}</span>
      -
      <span class="font-semibold text-slate-700">{{ Math.min(page * perPage, total) }}</span>
      จาก <span class="font-semibold text-slate-800">{{ total }}</span> รายการ
    </div>

    <!-- Pagination Controls -->
    <div class="flex items-center gap-1">
      <!-- First Page -->
      <button 
        class="pagination-btn"
        :class="{ 'pagination-btn-disabled': page <= 1 }"
        :disabled="page <= 1"
        @click="$emit('goto', 1)"
        title="หน้าแรก"
        data-testid="pagination-first"
      >
        <ChevronsLeft class="w-4 h-4" />
      </button>

      <!-- Previous Page -->
      <button 
        class="pagination-btn"
        :class="{ 'pagination-btn-disabled': page <= 1 }"
        :disabled="page <= 1"
        @click="$emit('prev')"
        title="หน้าก่อนหน้า"
        data-testid="pagination-prev"
      >
        <ChevronLeft class="w-4 h-4" />
      </button>

      <!-- Page Indicator -->
      <div class="px-4 py-2 text-sm font-medium text-slate-600" data-testid="pagination-current">
        <span class="text-slate-800 font-semibold" data-testid="pagination-current-page">{{ page }}</span>
        <span class="text-slate-400 mx-1">/</span>
        <span data-testid="pagination-total-pages">{{ totalPages }}</span>
      </div>

      <!-- Next Page -->
      <button 
        class="pagination-btn"
        :class="{ 'pagination-btn-disabled': page >= totalPages }"
        :disabled="page >= totalPages"
        @click="$emit('next')"
        title="หน้าถัดไป"
        data-testid="pagination-next"
      >
        <ChevronRight class="w-4 h-4" />
      </button>

      <!-- Last Page -->
      <button 
        class="pagination-btn"
        :class="{ 'pagination-btn-disabled': page >= totalPages }"
        :disabled="page >= totalPages"
        @click="$emit('goto', totalPages)"
        title="หน้าสุดท้าย"
        data-testid="pagination-last"
      >
        <ChevronsRight class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.pagination-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 0.5rem;
  border: 1px solid #e2e8f0;
  background: white;
  color: #475569;
  transition: all 0.15s ease;
  cursor: pointer;
}

.pagination-btn:hover:not(:disabled) {
  background: #f8fafc;
  border-color: #cbd5e1;
  color: #1e293b;
}

.pagination-btn:active:not(:disabled) {
  background: #f1f5f9;
  transform: scale(0.95);
}

.pagination-btn-disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>