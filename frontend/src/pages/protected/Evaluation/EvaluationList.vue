<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { EvaluationResponse } from '@/interfaces/evaluation'
import type { ApplicationScholarshipResponse } from '@/interfaces'
import { EvaluationService } from '@/services/evaluation/evaluation'
import { getAllApplicationScholarships } from '@/services/api/application'
import Swal from 'sweetalert2'
import EvaluationFormModal from './EvaluationFormModal.vue'
import EvaluationDetailModal from './EvaluationDetailModal.vue'
import StatsGrid from '@/components/ui/StatsGrid.vue'
import type { StatItem } from '@/components/ui/StatsGrid.vue'

// Icons
import { 
  Search, X, RefreshCw, Eye, FileEdit, 
  Users, CheckCircle, XCircle, Clock, Award, ChevronDown, Plus, UserPlus
} from 'lucide-vue-next'

// ========== State ==========
const evaluations = ref<EvaluationResponse[]>([])
const qualifiedApplicants = ref<ApplicationScholarshipResponse[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const searchQuery = ref('')
const activeTab = ref<'evaluations' | 'qualified'>('evaluations')
const creatingEvaluation = ref(false)

// Filter State
const selectedScholarshipId = ref<number | null>(null)
const selectedRoundId = ref<number | null>(null)
const statusFilter = ref<string>('')

// Modal State
const isFormModalOpen = ref(false)
const isDetailModalOpen = ref(false)
const selectedEvaluationId = ref<number | null>(null)

// Inline Decision State
const pendingDecisions = ref<Record<number, string>>({})

// ========== Computed: Filters Options ==========
const scholarships = computed(() => {
  const map = new Map<number, string>()
  evaluations.value.forEach(e => {
    const scholarship = e.application_scholarship?.scholarship
    if (scholarship?.ID) {
      map.set(scholarship.ID, scholarship.scholarship_name || `ทุน #${scholarship.ID}`)
    }
  })
  return Array.from(map, ([id, name]) => ({ id, name }))
})

const interviewRounds = computed(() => {
  const map = new Map<number, { id: number; name: string; scholarshipId: number }>()
  evaluations.value.forEach(e => {
    const round = e.interview_round
    const scholarshipId = e.application_scholarship?.scholarship?.ID
    if (round?.ID) {
      map.set(round.ID, {
        id: round.ID,
        name: round.name || `รอบ #${round.ID}`,
        scholarshipId: scholarshipId || 0
      })
    }
  })
  // Filter by selected scholarship
  let result = Array.from(map.values())
  if (selectedScholarshipId.value) {
    result = result.filter(r => r.scholarshipId === selectedScholarshipId.value)
  }
  return result
})

// ========== Computed: Filtered Data ==========
const filteredEvaluations = computed(() => {
  let result = evaluations.value

  // Filter by scholarship
  if (selectedScholarshipId.value) {
    result = result.filter(e => 
      e.application_scholarship?.scholarship?.ID === selectedScholarshipId.value
    )
  }

  // Filter by round
  if (selectedRoundId.value) {
    result = result.filter(e => e.interview_round_id === selectedRoundId.value)
  }

  // Filter by status
  if (statusFilter.value) {
    result = result.filter(e => e.status === statusFilter.value)
  }

  // Filter by search
  const term = searchQuery.value.toLowerCase().trim()
  if (term) {
    result = result.filter(e => {
      const name = getStudentName(e).toLowerCase()
      const studentId = getStudentId(e).toLowerCase()
      return name.includes(term) || studentId.includes(term)
    })
  }

  // Sort by total_score (descending)
  return [...result].sort((a, b) => (b.total_score || 0) - (a.total_score || 0))
})

// ========== Computed: Stats ==========
const stats = computed(() => {
  const data = filteredEvaluations.value
  return {
    total: data.length,
    approved: data.filter(e => e.status === 'approved').length,
    rejected: data.filter(e => e.status === 'rejected').length,
    completed: data.filter(e => e.status === 'completed').length,
    pending: data.filter(e => e.status === 'pending' || e.status === 'in_progress').length,
  }
})

// Computed stats for StatsGrid
const evaluationStats = computed<StatItem[]>(() => [
  { 
    title: 'ผู้สมัครทั้งหมด', 
    value: stats.value.total, 
    description: 'รายการ',
    icon: 'users', 
    color: 'blue' 
  },
  { 
    title: 'อนุมัติ', 
    value: stats.value.approved, 
    description: 'ผ่านเกณฑ์',
    icon: 'check', 
    color: 'green' 
  },
  { 
    title: 'รอตัดสิน', 
    value: stats.value.completed, 
    description: 'ประเมินเสร็จแล้ว',
    icon: 'award', 
    color: 'purple' 
  },
  { 
    title: 'รอประเมิน', 
    value: stats.value.pending, 
    description: 'รอดำเนินการ',
    icon: 'clock', 
    color: 'orange' 
  },
])

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

// Fetch qualified applicants (พร้อมประเมิน)
// ดึงทั้ง 'qualified' (รอจอง) และ 'interview_scheduled' (จองแล้ว รอประเมิน)
async function fetchQualifiedApplicants() {
  try {
    // Fetch both statuses
    const [qualified, scheduled] = await Promise.all([
      getAllApplicationScholarships('qualified'),
      getAllApplicationScholarships('interview_scheduled')
    ])
    const allApplicants = [...qualified, ...scheduled]
    
    // Filter out those who already have evaluations
    const evaluatedIds = new Set(evaluations.value.map(e => e.application_scholarship_id))
    qualifiedApplicants.value = allApplicants.filter(app => !evaluatedIds.has(app.ID))
  } catch (err: any) {
    console.error('Failed to fetch qualified applicants:', err)
  }
}

// Create evaluation for a qualified applicant
async function createEvaluationForApplicant(app: ApplicationScholarshipResponse) {
  if (creatingEvaluation.value) return
  
  // Get interview_round_id from the applicant's booking
  let interviewRoundId: number | null = null
  if (app.interviewe_bookings && app.interviewe_bookings.length > 0) {
    // Get from the latest booking's slot -> interview_round
    const latestBooking = app.interviewe_bookings[app.interviewe_bookings.length - 1]
    if (latestBooking?.slot?.interview_round_id) {
      interviewRoundId = latestBooking.slot.interview_round_id
    }
  }
  
  // If no booking found, use default or ask user
  if (!interviewRoundId) {
    // Use first available interview round from existing evaluations
    const firstRound = interviewRounds.value[0]
    if (firstRound) {
      interviewRoundId = firstRound.id
    } else {
      interviewRoundId = 1 // Fallback to 1
    }
  }
  
  const result = await Swal.fire({
    title: 'สร้างการประเมิน',
    text: `ต้องการสร้างการประเมินสำหรับ ${app.application?.student_profile?.first_name_th || 'ผู้สมัคร'} หรือไม่?`,
    icon: 'question',
    showCancelButton: true,
    confirmButtonText: 'สร้าง',
    cancelButtonText: 'ยกเลิก',
    confirmButtonColor: '#1e3a8a'
  })
  
  if (!result.isConfirmed) return
  
  creatingEvaluation.value = true
  try {
    await EvaluationService.create({
      application_scholarship_id: app.ID,
      interview_round_id: interviewRoundId,
      admin_id: 1 // TODO: Should be current admin ID from auth
    })
    
    await Swal.fire({
      icon: 'success',
      title: 'สร้างการประเมินสำเร็จ',
      timer: 1500,
      showConfirmButton: false
    })
    
    // Refresh data
    await fetchEvaluations()
    await fetchQualifiedApplicants()
    activeTab.value = 'evaluations'
  } catch (err: any) {
    Swal.fire({
      icon: 'error',
      title: 'ไม่สามารถสร้างการประเมินได้',
      text: err?.response?.data?.error || 'เกิดข้อผิดพลาด'
    })
  } finally {
    creatingEvaluation.value = false
  }
}

function getStudentName(e: EvaluationResponse): string {
  const profile = e.application_scholarship?.application?.student_profile
  if (!profile) return '-'
  return `${profile.first_name_th || ''} ${profile.last_name_th || ''}`.trim() || '-'
}

function getStudentId(e: EvaluationResponse): string {
  return e.application_scholarship?.application?.student_profile?.student_id || '-'
}

function getScholarshipName(e: EvaluationResponse): string {
  return e.application_scholarship?.scholarship?.scholarship_name || '-'
}

function getRoundName(e: EvaluationResponse): string {
  return e.interview_round?.name || '-'
}

function getStatusInfo(status: string) {
  switch (status) {
    case 'pending':
      return { label: 'รอประเมิน', bg: 'bg-amber-100', text: 'text-amber-700', icon: Clock }
    case 'in_progress':
      return { label: 'กำลังประเมิน', bg: 'bg-blue-100', text: 'text-blue-700', icon: RefreshCw }
    case 'completed':
      return { label: 'รอตัดสิน', bg: 'bg-purple-100', text: 'text-purple-700', icon: Award }
    case 'approved':
      return { label: 'อนุมัติ', bg: 'bg-emerald-100', text: 'text-emerald-700', icon: CheckCircle }
    case 'rejected':
      return { label: 'ไม่อนุมัติ', bg: 'bg-red-100', text: 'text-red-700', icon: XCircle }
    default:
      return { label: status, bg: 'bg-gray-100', text: 'text-gray-700', icon: Clock }
  }
}

function getScorePercentage(score: number | undefined, max: number = 100): number {
  if (!score) return 0
  return Math.min((score / max) * 100, 100)
}

function getProgressColor(percentage: number): string {
  if (percentage >= 80) return 'bg-emerald-500'
  if (percentage >= 60) return 'bg-blue-500'
  if (percentage >= 40) return 'bg-amber-500'
  return 'bg-red-500'
}

function getRankClass(index: number): string {
  if (index === 0) return 'bg-amber-400 text-white' // Gold
  if (index === 1) return 'bg-gray-400 text-white' // Silver
  if (index === 2) return 'bg-amber-600 text-white' // Bronze
  return 'bg-gray-100 text-gray-600'
}

// Modal handlers
function openFormModal(e: EvaluationResponse) {
  selectedEvaluationId.value = e.ID
  isFormModalOpen.value = true
}

function openDetailModal(e: EvaluationResponse) {
  selectedEvaluationId.value = e.ID
  isDetailModalOpen.value = true
}

function onFormCompleted() {
  fetchEvaluations()
}

// Inline Decision handlers
function getDecision(evaluationId: number): string {
  return pendingDecisions.value[evaluationId] || ''
}

function setDecision(evaluationId: number, decision: string) {
  pendingDecisions.value[evaluationId] = decision
}

async function saveDecision(e: EvaluationResponse) {
  const decision = pendingDecisions.value[e.ID] as 'approved' | 'rejected' | 'waitlist'
  if (!decision) return

  try {
    await EvaluationService.complete(e.ID, { final_decision: decision })
    await Swal.fire({ icon: 'success', title: 'บันทึกผลสำเร็จ', timer: 1500, showConfirmButton: false })
    delete pendingDecisions.value[e.ID]
    await fetchEvaluations()
  } catch (err: any) {
    Swal.fire({
      icon: 'error',
      title: 'ไม่สามารถบันทึกได้',
      text: err?.response?.data?.error || 'เกิดข้อผิดพลาด',
    })
  }
}

// Reset round when scholarship changes
watch(selectedScholarshipId, () => {
  selectedRoundId.value = null
})

// ========== Lifecycle ==========
onMounted(async () => {
  await fetchEvaluations()
  await fetchQualifiedApplicants()
})
</script>

<template>
  <div 
    class="evaluation-list-wrapper w-full mx-auto flex flex-col h-full p-6 bg-white" 
    data-theme="light"
  >
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold text-[#1e3a8a]">การประเมินผู้สมัครทุน</h1>

      <div class="flex items-center gap-3 flex-wrap">

        <!-- Filters Row -->
        <!-- Scholarship Filter -->
        <div class="relative">
          <select 
            v-model="selectedScholarshipId"
            class="appearance-none bg-white border border-gray-300 rounded-full px-4 py-2 pr-10 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option :value="null">ทุกทุนการศึกษา</option>
            <option v-for="s in scholarships" :key="s.id" :value="s.id">{{ s.name }}</option>
          </select>
          <ChevronDown class="w-4 h-4 text-gray-400 absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none" />
        </div>

        <!-- Round Filter -->
        <div class="relative">
          <select 
            v-model="selectedRoundId"
            class="appearance-none bg-white border border-gray-300 rounded-full px-4 py-2 pr-10 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            :disabled="interviewRounds.length === 0"
          >
            <option :value="null">ทุกรอบ</option>
            <option v-for="r in interviewRounds" :key="r.id" :value="r.id">{{ r.name }}</option>
          </select>
          <ChevronDown class="w-4 h-4 text-gray-400 absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none" />
        </div>

        <!-- Status Filter -->
        <select 
          v-model="statusFilter"
          class="px-4 py-2 border border-gray-300 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">ทุกสถานะ</option>
          <option value="pending">รอประเมิน</option>
          <option value="in_progress">กำลังประเมิน</option>
          <option value="completed">รอตัดสิน</option>
          <option value="approved">อนุมัติ</option>
          <option value="rejected">ไม่อนุมัติ</option>
        </select>

        <!-- Search Bar -->
        <div class="relative flex-1 md:w-64">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="ค้นหาชื่อ / รหัส..."
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

    <!-- Tab Navigation -->
    <div class="flex gap-2 mb-6">
      <button
        @click="activeTab = 'evaluations'"
        class="px-4 py-2 rounded-full text-sm font-medium transition-colors"
        :class="activeTab === 'evaluations' 
          ? 'bg-[#1e3a8a] text-white' 
          : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
      >
        <Award class="w-4 h-4 inline-block mr-1" />
        การประเมิน ({{ filteredEvaluations.length }})
      </button>
      <button
        @click="activeTab = 'qualified'"
        class="px-4 py-2 rounded-full text-sm font-medium transition-colors flex items-center gap-2 relative"
        :class="activeTab === 'qualified' 
          ? 'bg-emerald-600 text-white' 
          : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
      >
        <UserPlus class="w-4 h-4" />
        ผู้สมัครพร้อมประเมิน
        <!-- Badge -->
        <span 
          v-if="qualifiedApplicants.length > 0 && activeTab !== 'qualified'"
          class="absolute -top-1 -right-1 min-w-5 h-5 px-1.5 bg-red-500 text-white text-xs font-bold rounded-full flex items-center justify-center animate-pulse"
        >
          {{ qualifiedApplicants.length }}
        </span>
        <span v-else class="text-xs opacity-70">({{ qualifiedApplicants.length }})</span>
      </button>
    </div>

    <!-- Stats Cards -->
    <StatsGrid 
      v-if="activeTab === 'evaluations'"
      :stats="evaluationStats" 
      :columns="4"
      class="mb-6"
    />

    <!-- Loading State -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-r-transparent"></div>
        <p class="mt-2 text-gray-500">กำลังโหลดข้อมูล...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="p-6">
      <div class="p-4 bg-red-100 text-red-700 rounded-lg">{{ error }}</div>
    </div>

    <!-- Qualified Applicants Table -->
    <div v-else-if="activeTab === 'qualified'" class="flex-1 min-h-0 flex flex-col gap-4">
      <div class="overflow-x-auto overflow-y-auto flex-1 min-h-[400px] rounded-xl border border-slate-200 bg-white shadow-sm">
        <table class="table table-sm w-full" data-theme="light">
          <thead>
            <tr class="bg-emerald-50 border-b border-emerald-200">
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-emerald-700 text-left">#</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-emerald-700 text-left">ชื่อนักศึกษา</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-emerald-700 text-left hidden md:table-cell">ทุนการศึกษา</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-emerald-700 text-center">สถานะ</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-emerald-700 text-center">การจัดการ</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr 
              v-for="(app, index) in qualifiedApplicants" 
              :key="app.ID"
              class="hover:bg-emerald-50/50 transition-colors"
            >
              <td class="py-3 px-4">
                <span class="inline-flex items-center justify-center w-8 h-8 rounded-full bg-emerald-100 text-emerald-600 text-sm font-bold">
                  {{ index + 1 }}
                </span>
              </td>
              <td class="py-3 px-4">
                <p class="font-semibold text-gray-900">
                  {{ app.application?.student_profile?.first_name_th || '' }} {{ app.application?.student_profile?.last_name_th || '' }}
                </p>
                <p class="text-sm text-gray-500">{{ app.application?.student_profile?.student_id || '-' }}</p>
              </td>
              <td class="py-3 px-4 text-sm text-gray-600 hidden md:table-cell">
                {{ app.scholarship?.scholarship_name || '-' }}
              </td>
              <td class="py-3 px-4 text-center">
                <span class="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-semibold rounded-full bg-emerald-100 text-emerald-700">
                  <CheckCircle class="w-3.5 h-3.5" />
                  พร้อมประเมิน
                </span>
              </td>
              <td class="py-3 px-4 text-center">
                <button
                  @click="createEvaluationForApplicant(app)"
                  :disabled="creatingEvaluation"
                  class="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-sm font-medium transition-colors inline-flex items-center gap-1 disabled:opacity-50"
                >
                  <Plus class="w-4 h-4" />
                  สร้างการประเมิน
                </button>
              </td>
            </tr>
            <tr v-if="qualifiedApplicants.length === 0">
              <td colspan="5" class="h-[300px] text-center align-middle">
                <div class="flex flex-col items-center justify-center text-gray-400">
                  <CheckCircle class="w-12 h-12 mb-3 opacity-50" />
                  <p class="text-gray-500">ไม่มีผู้สมัครที่พร้อมประเมิน</p>
                  <p class="text-sm text-gray-400">ผู้สมัครทุกคนได้รับการประเมินแล้ว</p>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Evaluations Content -->
    <div v-else-if="!loading && !error && activeTab === 'evaluations'" class="flex-1 min-h-0 flex flex-col gap-4">
      <!-- Table -->
      <div class="overflow-x-auto overflow-y-auto flex-1 min-h-[400px] rounded-xl border border-slate-200 bg-white shadow-sm">
        <table class="table table-sm w-full" data-theme="light">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-200">
                <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-left">#</th>
                <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-left">ชื่อนักศึกษา</th>
                <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-left hidden md:table-cell">ทุนการศึกษา</th>
                <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-left hidden lg:table-cell">รอบสัมภาษณ์</th>
                <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-center">คะแนนรวม</th>
                <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-center">สถานะ</th>
                <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-center">การจัดการ</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr 
                v-for="(e, index) in filteredEvaluations" 
                :key="e.ID"
                class="hover:bg-blue-50/50 transition-colors"
                :class="index % 2 === 0 ? 'bg-white' : 'bg-slate-50/30'"
              >
                <!-- Rank -->
                <td class="py-3 px-4">
                  <span 
                    class="inline-flex items-center justify-center w-8 h-8 rounded-full text-sm font-bold"
                    :class="getRankClass(index)"
                  >
                    {{ index + 1 }}
                  </span>
                </td>

                <!-- Student Name -->
                <td class="py-3 px-4">
                  <p class="font-semibold text-gray-900">{{ getStudentName(e) }}</p>
                  <p class="text-sm text-gray-500">{{ getStudentId(e) }}</p>
                </td>

                <!-- Scholarship (hidden on mobile) -->
                <td class="py-3 px-4 text-sm text-gray-600 hidden md:table-cell">
                  {{ getScholarshipName(e) }}
                </td>

                <!-- Round (hidden on tablet) -->
                <td class="py-3 px-4 text-sm text-gray-600 hidden lg:table-cell">
                  {{ getRoundName(e) }}
                </td>

                <!-- Score -->
                <td class="py-3 px-4 text-center">
                  <div class="flex flex-col items-center">
                    <span 
                      class="font-bold text-lg"
                      :class="e.total_score >= 70 ? 'text-emerald-600' : e.total_score >= 50 ? 'text-amber-600' : 'text-gray-600'"
                      v-if="e.total_score"
                    >
                      {{ e.total_score.toFixed(1) }}
                    </span>
                    <span v-else class="text-gray-400">-</span>
                    <div v-if="e.total_score" class="w-16 h-1.5 bg-gray-200 rounded-full mt-1 overflow-hidden">
                      <div 
                        class="h-full rounded-full transition-all"
                        :class="getProgressColor(getScorePercentage(e.total_score))"
                        :style="{ width: `${getScorePercentage(e.total_score)}%` }"
                      ></div>
                    </div>
                  </div>
                </td>

                <!-- Status -->
                <td class="py-3 px-4 text-center">
                  <span 
                    class="inline-flex items-center gap-1 px-2.5 py-1 text-xs font-semibold rounded-full"
                    :class="[getStatusInfo(e.status).bg, getStatusInfo(e.status).text]"
                  >
                    <component :is="getStatusInfo(e.status).icon" class="w-3.5 h-3.5" />
                    {{ getStatusInfo(e.status).label }}
                  </span>
                </td>

                <!-- Actions -->
                <td class="py-3 px-4">
                  <div class="flex items-center justify-center gap-2">
                    <!-- View Button (always show) -->
                    <button
                      @click="openDetailModal(e)"
                      class="p-2 text-gray-500 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                      title="ดูรายละเอียด"
                    >
                      <Eye class="w-4 h-4" />
                    </button>

                    <!-- Evaluate Button (for pending/in_progress) -->
                    <button 
                      v-if="e.status === 'pending' || e.status === 'in_progress'"
                      @click="openFormModal(e)"
                      class="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                      title="ประเมิน"
                    >
                      <FileEdit class="w-4 h-4" />
                    </button>

                    <!-- Decision Dropdown (for completed) -->
                    <template v-if="e.status === 'completed'">
                      <select
                        :value="getDecision(e.ID)"
                        @change="setDecision(e.ID, ($event.target as HTMLSelectElement).value)"
                        class="text-xs border border-gray-300 rounded-lg px-2 py-1.5 focus:ring-2 focus:ring-blue-500"
                      >
                        <option value="">ตัดสิน...</option>
                        <option value="approved">อนุมัติ</option>
                        <option value="rejected">ไม่อนุมัติ</option>
                        <option value="waitlist">รอพิจารณา</option>
                      </select>
                      <button
                        v-if="getDecision(e.ID)"
                        @click="saveDecision(e)"
                        class="px-2 py-1 bg-emerald-600 hover:bg-emerald-700 text-white text-xs font-medium rounded-lg transition-colors"
                      >
                        บันทึก
                      </button>
                    </template>

                    <!-- Status text (for approved/rejected) -->
                    <span v-if="e.status === 'approved'" class="text-xs font-medium text-emerald-600 px-2">
                      ✓ อนุมัติแล้ว
                    </span>
                    <span v-if="e.status === 'rejected'" class="text-xs font-medium text-red-600 px-2">
                      ✗ ไม่ผ่าน
                    </span>
                  </div>
                </td>
              </tr>

              <!-- Empty State -->
              <tr v-if="filteredEvaluations.length === 0">
                <td colspan="7" class="h-[300px] text-center align-middle">
                  <div class="flex flex-col items-center justify-center text-gray-400">
                    <Users class="w-12 h-12 mb-3 opacity-50" />
                    <p class="text-gray-500">ไม่พบข้อมูลการประเมิน</p>
                    <p class="text-sm text-gray-400">ลองเปลี่ยนเงื่อนไขการกรอง</p>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

    <!-- Modals -->
    <EvaluationFormModal
      v-model:isOpen="isFormModalOpen"
      :evaluationId="selectedEvaluationId"
      @completed="onFormCompleted"
    />

    <EvaluationDetailModal
      v-model:isOpen="isDetailModalOpen"
      :evaluationId="selectedEvaluationId"
      @close="isDetailModalOpen = false"
    />
  </div>
</template>
