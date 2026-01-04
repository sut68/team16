<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import type { EvaluationResponse } from '@/interfaces/evaluation'
  import { EvaluationService } from '@/services/evaluation/evaluation'
  import Swal from 'sweetalert2'
  import EvaluationFormModal from './EvaluationFormModal.vue'
  import StatCard from '@/components/ui/StatCard.vue'

  // Icons
  import { 
    Search, X, Filter, ClipboardCheck, 
    Clock, CheckCircle, XCircle, AlertCircle, RefreshCw, Trash2
  } from 'lucide-vue-next'

  // ========== State ==========
  const evaluations = ref<EvaluationResponse[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const searchQuery = ref('')

  // Modal State
  const isFormModalOpen = ref(false)
  const selectedEvaluationId = ref<number | null>(null)
  const statusFilter = ref<string>('')

  // ========== Computed ==========
  const filteredEvaluations = computed(() => {
    let result = evaluations.value

    // Filter by status
    if (statusFilter.value) {
      result = result.filter(e => e.status === statusFilter.value)
    }

    // Filter by search
    const term = searchQuery.value.toLowerCase().trim()
    if (term) {
      result = result.filter(e => {
        const studentName = e.application_scholarship?.application?.student_profile?.firstname || ''
        const scholarshipName = e.application_scholarship?.scholarship?.name || ''
        const roundName = e.interview_round?.name || ''
        return studentName.toLowerCase().includes(term) ||
          scholarshipName.toLowerCase().includes(term) ||
          roundName.toLowerCase().includes(term)
      })
    }

    return result
  })

  // Stats
  const totalEvaluations = computed(() => evaluations.value.length)
  const pendingCount = computed(() => evaluations.value.filter(e => e.status === 'pending').length)
  const inProgressCount = computed(() => evaluations.value.filter(e => e.status === 'in_progress').length)
  const completedCount = computed(() => evaluations.value.filter(e => ['completed', 'approved', 'rejected'].includes(e.status)).length)

  // ========== Methods ==========
  async function fetchEvaluations() {
    loading.value = true
    error.value = null
    try {
      evaluations.value = await EvaluationService.getAll()
    } catch (err: any) {
      console.error('Failed to fetch evaluations:', err)
      error.value = err?.message || 'ไม่สามารถโหลดข้อมูลได้'
    } finally {
      loading.value = false
    }
  }

  function getStudentName(evaluation: EvaluationResponse): string {
    const profile = evaluation.application_scholarship?.application?.student_profile
    if (!profile) return '-'
    return `${profile.firstname || ''} ${profile.lastname || ''}`.trim() || '-'
  }

  function getScholarshipName(evaluation: EvaluationResponse): string {
    return evaluation.application_scholarship?.scholarship?.name || '-'
  }

  function getRoundName(evaluation: EvaluationResponse): string {
    return evaluation.interview_round?.name || '-'
  }

  function getStatusInfo(status: string) {
    switch (status) {
      case 'pending':
        return { label: 'รอประเมิน', class: 'bg-amber-100 text-amber-700', icon: Clock }
      case 'in_progress':
        return { label: 'กำลังประเมิน', class: 'bg-blue-100 text-blue-700', icon: RefreshCw }
      case 'completed':
        return { label: 'ประเมินเสร็จ', class: 'bg-slate-100 text-slate-700', icon: AlertCircle }
      case 'approved':
        return { label: 'อนุมัติ', class: 'bg-emerald-100 text-emerald-700', icon: CheckCircle }
      case 'rejected':
        return { label: 'ไม่อนุมัติ', class: 'bg-red-100 text-red-700', icon: XCircle }
      default:
        return { label: status, class: 'bg-gray-100 text-gray-700', icon: AlertCircle }
    }
  }

  function openEvaluationForm(evaluation: EvaluationResponse) {
    selectedEvaluationId.value = evaluation.ID
    isFormModalOpen.value = true
  }

  function onFormModalCompleted() {
    fetchEvaluations()
  }

  async function handleDelete(evaluation: EvaluationResponse) {
    const result = await Swal.fire({
      icon: 'warning',
      title: 'ยืนยันการลบ',
      text: `ต้องการลบการประเมินของ "${getStudentName(evaluation)}" ใช่หรือไม่?`,
      showCancelButton: true,
      confirmButtonColor: '#ef4444',
      confirmButtonText: 'ลบ',
      cancelButtonText: 'ยกเลิก',
    })

    if (result.isConfirmed) {
      try {
        await EvaluationService.delete(evaluation.ID)
        await Swal.fire({ icon: 'success', title: 'ลบสำเร็จ', timer: 1500, showConfirmButton: false })
        await fetchEvaluations()
      } catch (err: any) {
        Swal.fire({
          icon: 'error',
          title: 'ไม่สามารถลบได้',
          text: err?.response?.data?.error || 'เกิดข้อผิดพลาด',
        })
      }
    }
  }

  // ========== Lifecycle ==========
  onMounted(fetchEvaluations)
</script>

<template>
  <div 
    class="evaluation-list-wrapper w-full mx-auto flex flex-col h-full p-6 bg-white" 
    data-theme="light"
  >
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold text-[#1e3a8a]">รายการประเมินผู้สมัคร</h1>

      <div class="flex items-center gap-3 flex-wrap">
        <!-- Status Filter -->
        <select
          v-model="statusFilter"
          class="px-4 py-2 border border-gray-300 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">ทุกสถานะ</option>
          <option value="pending">รอประเมิน</option>
          <option value="in_progress">กำลังประเมิน</option>
          <option value="completed">ประเมินเสร็จ</option>
          <option value="approved">อนุมัติ</option>
          <option value="rejected">ไม่อนุมัติ</option>
        </select>

        <!-- Search Bar -->
        <div class="relative flex-1 md:w-64">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="ค้นหาชื่อ / ทุน..."
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

        <!-- Refresh Button -->
        <button 
          @click="fetchEvaluations"
          :disabled="loading"
          class="p-2 bg-gray-100 hover:bg-gray-200 rounded-full transition-colors"
          title="รีเฟรช"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-6">
      <StatCard
        label="ทั้งหมด"
        :value="totalEvaluations"
        icon="total"
        color="slate"
      />
      <StatCard
        label="รอประเมิน"
        :value="pendingCount"
        icon="pending"
        color="amber"
      />
      <StatCard
        label="กำลังประเมิน"
        :value="inProgressCount"
        icon="inProgress"
        color="blue"
      />
      <StatCard
        label="เสร็จสิ้น"
        :value="completedCount"
        icon="completed"
        color="green"
      />
    </div>

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
      <div class="overflow-x-auto overflow-y-auto flex-1 min-h-[400px] rounded-xl border border-slate-200 bg-white shadow-sm">
        <table class="table table-sm w-full">
          <thead>
            <tr class="bg-slate-800 text-white">
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">#</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">ชื่อนักศึกษา</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 hidden md:table-cell">ทุนการศึกษา</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 hidden lg:table-cell">รอบสัมภาษณ์</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-center">คะแนนรวม</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-center">สถานะ</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-right">การจัดการ</th>
            </tr>
          </thead>

          <tbody class="divide-y divide-slate-100">
            <tr 
              v-for="(e, index) in filteredEvaluations" 
              :key="e.ID" 
              class="transition-all duration-200 hover:bg-blue-50/50"
              :class="index % 2 === 0 ? 'bg-white' : 'bg-slate-50/50'"
            >
              <td class="py-3 px-4 text-sm text-slate-500 font-mono">{{ e.ID }}</td>
              
              <td class="py-3 px-4">
                <span class="font-semibold text-slate-800">{{ getStudentName(e) }}</span>
              </td>

              <td class="py-3 px-4 text-sm text-slate-600 hidden md:table-cell">
                {{ getScholarshipName(e) }}
              </td>

              <td class="py-3 px-4 text-sm text-slate-600 hidden lg:table-cell">
                {{ getRoundName(e) }}
              </td>

              <td class="py-3 px-4 text-center">
                <span 
                  class="font-bold text-lg"
                  :class="e.total_score >= 70 ? 'text-emerald-600' : e.total_score >= 50 ? 'text-amber-600' : 'text-slate-600'"
                >
                  {{ e.total_score.toFixed(1) }}
                </span>
              </td>

              <td class="py-3 px-4 text-center">
                <span 
                  class="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-semibold rounded-full"
                  :class="getStatusInfo(e.status).class"
                >
                  <component :is="getStatusInfo(e.status).icon" class="w-3.5 h-3.5" />
                  {{ getStatusInfo(e.status).label }}
                </span>
              </td>

              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-2">

                  <!-- Evaluate -->
                  <button 
                    v-if="e.status === 'pending' || e.status === 'in_progress'"
                    @click="openEvaluationForm(e)"
                    class="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                    title="ประเมิน"
                  >
                    <ClipboardCheck class="w-4 h-4" />
                  </button>

                  <!-- Delete -->
                  <button 
                    @click="handleDelete(e)"
                    class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    title="ลบ"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-if="filteredEvaluations.length === 0">
              <td colspan="7" class="h-[300px] text-center align-middle">
                <div class="flex flex-col items-center justify-center text-gray-400">
                  <Filter class="w-12 h-12 mb-3 opacity-50" />
                  <p v-if="searchQuery || statusFilter" class="text-gray-500">
                    ไม่พบการประเมินที่ตรงกับเงื่อนไข
                  </p>
                  <p v-else class="text-gray-500">ยังไม่มีการประเมินในระบบ</p>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
    </div>
    </div>

    <!-- Evaluation Form Modal -->
    <EvaluationFormModal
      v-model:isOpen="isFormModalOpen"
      :evaluationId="selectedEvaluationId"
      @completed="onFormModalCompleted"
    />
  </div>
</template>

<style scoped>
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
