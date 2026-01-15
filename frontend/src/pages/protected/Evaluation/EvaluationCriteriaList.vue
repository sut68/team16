<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import type { EvaluationCriterionResponse } from '@/interfaces/evaluation'
  import { EvaluationCriteriaService } from '@/services/evaluation/evaluation'
  import Swal from 'sweetalert2'
  import StatsGrid from '@/components/ui/StatsGrid.vue'
  import type { StatItem } from '@/components/ui/StatsGrid.vue'
  import CriteriaFormModal from './CriteriaFormModal.vue'

  // Icons
  import { Plus, Pencil, Trash2, Search, X, Filter, CheckCircle, XCircle } from 'lucide-vue-next'

  // ========== State ==========
  const criteria = ref<EvaluationCriterionResponse[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const searchQuery = ref('')

  // Modal State
  const isFormModalOpen = ref(false)
  const selectedCriterion = ref<EvaluationCriterionResponse | null>(null)

  // ========== Computed ==========
  const filteredCriteria = computed(() => {
    const term = searchQuery.value.toLowerCase().trim()
    if (!term) return criteria.value
    return criteria.value.filter(c => 
      c.name.toLowerCase().includes(term) ||
      c.description?.toLowerCase().includes(term)
    )
  })

  const totalCriteria = computed(() => criteria.value.length)
  const activeCriteria = computed(() => criteria.value.filter(c => c.is_active).length)
  const inactiveCriteria = computed(() => criteria.value.filter(c => !c.is_active).length)

  // Computed stats for StatsGrid
  const criteriaStats = computed<StatItem[]>(() => [
    { 
      title: 'เกณฑ์ทั้งหมด', 
      value: totalCriteria.value, 
      description: 'รายการ',
      icon: 'clipboard', 
      color: 'blue' 
    },
    { 
      title: 'ใช้งานอยู่', 
      value: activeCriteria.value, 
      description: 'เปิดใช้งาน',
      icon: 'check', 
      color: 'green' 
    },
    { 
      title: 'ปิดใช้งาน', 
      value: inactiveCriteria.value, 
      description: 'ไม่ได้ใช้',
      icon: 'clock', 
      color: 'slate' 
    },
  ])

  // ========== Methods ==========
  async function fetchCriteria() {
    loading.value = true
    error.value = null
    try {
      criteria.value = await EvaluationCriteriaService.getAll()
    } catch (err: any) {
      console.error('Failed to fetch criteria:', err)
      error.value = err?.message || 'ไม่สามารถโหลดข้อมูลได้'
    } finally {
      loading.value = false
    }
  }

  function openCreateModal() {
    selectedCriterion.value = null
    isFormModalOpen.value = true
  }

  function openEditModal(criterion: EvaluationCriterionResponse) {
    selectedCriterion.value = criterion
    isFormModalOpen.value = true
  }

  function closeFormModal() {
    isFormModalOpen.value = false
    selectedCriterion.value = null
  }

  async function onFormSaved() {
    await fetchCriteria()
  }

  async function handleDelete(criterion: EvaluationCriterionResponse) {
    const result = await Swal.fire({
      icon: 'warning',
      title: 'ยืนยันการลบ',
      text: `ต้องการลบเกณฑ์ "${criterion.name}" ใช่หรือไม่?`,
      showCancelButton: true,
      confirmButtonColor: '#ef4444',
      confirmButtonText: 'ลบ',
      cancelButtonText: 'ยกเลิก',
    })

    if (result.isConfirmed) {
      try {
        await EvaluationCriteriaService.delete(criterion.ID)
        await Swal.fire({ icon: 'success', title: 'ลบสำเร็จ', timer: 1500, showConfirmButton: false })
        await fetchCriteria()
      } catch (err: any) {
        console.error('Delete error:', err)
        Swal.fire({
          icon: 'error',
          title: 'ไม่สามารถลบได้',
          text: err?.response?.data?.error || 'เกณฑ์นี้อาจมีการใช้งานอยู่',
        })
      }
    }
  }

  async function toggleActive(criterion: EvaluationCriterionResponse) {
    try {
      await EvaluationCriteriaService.update(criterion.ID, { is_active: !criterion.is_active })
      await fetchCriteria()
    } catch (err: any) {
      console.error('Toggle error:', err)
      Swal.fire({ icon: 'error', title: 'เกิดข้อผิดพลาด' })
    }
  }

  function getScoreTypeLabel(type: string) {
    switch (type) {
      case 'numeric': return 'คะแนนตัวเลข'
      case 'grade': return 'เกรด (A-F)'
      case 'pass_fail': return 'ผ่าน/ไม่ผ่าน'
      default: return type
    }
  }

  function getScoreTypeBadgeClass(type: string) {
    switch (type) {
      case 'numeric': return 'bg-blue-100 text-blue-800'
      case 'grade': return 'bg-purple-100 text-purple-800'
      case 'pass_fail': return 'bg-amber-100 text-amber-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  // ========== Lifecycle ==========
  onMounted(fetchCriteria)
</script>

<template>
  <div 
    class="evaluation-criteria-wrapper w-full mx-auto flex flex-col h-full p-6 bg-white" 
    data-theme="light"
  >
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
      <h2 class="text-xl font-bold text-[#1e3a8a]">เกณฑ์การประเมิน</h2>

      <div class="flex items-center gap-3 w-full md:w-auto">
        <!-- Search Bar -->
        <div class="relative flex-1 md:w-64">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="ค้นหาเกณฑ์..."
            class="w-full pl-10 pr-8 py-2 border border-gray-300 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <!-- Add Button -->
        <button 
          @click="openCreateModal" 
          class="btn-ghost-rounded"
        >
          <Plus class="w-4 h-4" />
          <span class="hidden sm:inline">เพิ่มเกณฑ์</span>
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <StatsGrid 
      :stats="criteriaStats" 
      :columns="3"
      class="mb-6"
    />

    <!-- Loading State -->
    <div v-if="loading" class="p-6 w-full">
      <div class="animate-pulse space-y-3">
        <div class="h-6 bg-gray-200 rounded w-1/4"></div>
        <div class="bg-white rounded shadow overflow-hidden">
          <div class="p-4">
            <div class="h-3 bg-gray-200 rounded w-full mb-2"></div>
            <div class="h-3 bg-gray-200 rounded w-3/4"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error" class="p-4 bg-red-100 text-red-700 rounded-lg mb-4">{{ error }}</div>

    <!-- Content -->
    <div v-if="!loading && !error" class="flex-1 min-h-0 flex flex-col gap-4">
      <!-- Table -->
      <div class="overflow-x-auto overflow-y-auto flex-1 rounded-xl border border-slate-200 bg-white shadow-sm">
        <table class="table table-sm w-full">
          <thead>
            <tr class="bg-slate-800 text-white">
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">#</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">ชื่อเกณฑ์</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 hidden md:table-cell">รายละเอียด</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">ประเภท</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-center">คะแนนเต็ม</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-center">น้ำหนัก</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-center">สถานะ</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-right">การจัดการ</th>
            </tr>
          </thead>

          <tbody class="divide-y divide-slate-100">
            <tr 
              v-for="(c, index) in filteredCriteria" 
              :key="c.ID" 
              class="transition-all duration-200 hover:bg-blue-50/50"
              :class="index % 2 === 0 ? 'bg-white' : 'bg-slate-50/50'"
            >
              <td class="py-3 px-4 text-sm text-slate-500 font-mono">{{ c.ID }}</td>
              
              <td class="py-3 px-4">
                <span class="font-semibold text-slate-800">{{ c.name }}</span>
              </td>

              <td class="py-3 px-4 text-sm text-slate-600 hidden md:table-cell">
                <span class="truncate block max-w-xs" :title="c.description">
                  {{ c.description || '-' }}
                </span>
              </td>

              <td class="py-3 px-4">
                <span 
                  class="inline-flex items-center px-2.5 py-1 text-xs font-semibold rounded-md"
                  :class="getScoreTypeBadgeClass(c.score_type)"
                >
                  {{ getScoreTypeLabel(c.score_type) }}
                </span>
              </td>

              <td class="py-3 px-4 text-center">
                <span class="font-mono text-sm">{{ c.max_score }}</span>
              </td>

              <td class="py-3 px-4 text-center">
                <span class="font-mono text-sm">{{ c.weight }}</span>
              </td>

              <td class="py-3 px-4 text-center">
                <button 
                  @click="toggleActive(c)"
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-semibold cursor-pointer transition-colors"
                  :class="c.is_active 
                    ? 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200' 
                    : 'bg-slate-100 text-slate-500 hover:bg-slate-200'"
                >
                  <CheckCircle v-if="c.is_active" class="w-3.5 h-3.5" />
                  <XCircle v-else class="w-3.5 h-3.5" />
                  {{ c.is_active ? 'Active' : 'Inactive' }}
                </button>
              </td>

              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button 
                    @click="openEditModal(c)"
                    class="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                    title="แก้ไข"
                  >
                    <Pencil class="w-4 h-4" />
                  </button>
                  <button 
                    @click="handleDelete(c)"
                    class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    title="ลบ"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-if="filteredCriteria.length === 0">
              <td colspan="8" class="h-[300px] text-center align-middle">
                <div class="flex flex-col items-center justify-center text-gray-400">
                  <Filter class="w-12 h-12 mb-3 opacity-50" />
                  <p v-if="searchQuery" class="text-gray-500">
                    ไม่พบเกณฑ์ที่ตรงกับ "<strong class="text-gray-700">{{ searchQuery }}</strong>"
                  </p>
                  <p v-else class="text-gray-500">ยังไม่มีเกณฑ์ในระบบ</p>
                  <button 
                    v-if="!searchQuery"
                    @click="openCreateModal" 
                    class="btn-primary mt-4"
                  >
                    <Plus class="w-4 h-4" />
                    <span>เพิ่มเกณฑ์แรก</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Criteria Form Modal -->
    <CriteriaFormModal
      :is-open="isFormModalOpen"
      :criterion="selectedCriterion"
      @close="closeFormModal"
      @saved="onFormSaved"
    />
  </div>
</template>

<style scoped>
.btn-ghost-rounded {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding-left: 1rem;
  padding-right: 1rem;
  height: 2.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
  background-color: white;
  border: 1px solid #d1d5db;
  border-radius: 9999px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  transition: all 150ms ease;
  cursor: pointer;
}

.btn-ghost-rounded:hover:not(:disabled) {
  background-color: #f3f4f6;
  border-color: #9ca3af;
}

.btn-ghost-rounded:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead th {
  position: sticky;
  top: 0;
  background: white;
  z-index: 10;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #6b7280;
}
</style>
