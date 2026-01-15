<script setup lang="ts">
import { ref, computed, onMounted, inject } from 'vue'
import type { ApplicationScholarshipResponse } from '@/interfaces'
import { getAllApplicationScholarships } from '@/services/api/application'
import { EvaluationService } from '@/services/evaluation/evaluation'
import { Get } from '@/services/api/https'
import Swal from 'sweetalert2'

// Icons
import { 
  Search, X, RefreshCw, Plus, CheckCircle, UserPlus, ChevronDown
} from 'lucide-vue-next'

// Emit for parent to refresh count
const emit = defineEmits(['created'])

// Get refresh function from parent
const refreshParentCount = inject<() => void>('refreshQualifiedCount', () => {})

// ========== State ==========
const qualifiedApplicants = ref<ApplicationScholarshipResponse[]>([])
const evaluatedIds = ref<Set<number>>(new Set())
const loading = ref(false)
const creatingEvaluation = ref(false)
const adminId = ref<number>(1)
const searchQuery = ref('')
const selectedScholarship = ref<string>('all')

// ========== Computed ==========
const scholarshipOptions = computed(() => {
  const map = new Map<number, string>()
  qualifiedApplicants.value.forEach(app => {
    if (app.scholarship?.ID) {
      map.set(app.scholarship.ID, app.scholarship.scholarship_name || `ทุน #${app.scholarship.ID}`)
    }
  })
  return Array.from(map, ([id, name]) => ({ id, name }))
})

const filteredApplicants = computed(() => {
  let result = qualifiedApplicants.value
  
  // Filter out already evaluated
  result = result.filter(app => !evaluatedIds.value.has(app.ID))
  
  // Filter by scholarship
  if (selectedScholarship.value !== 'all') {
    const scholarshipId = parseInt(selectedScholarship.value)
    result = result.filter(app => app.scholarship?.ID === scholarshipId)
  }
  
  // Filter by search
  const term = searchQuery.value.toLowerCase().trim()
  if (term) {
    result = result.filter(app => {
      const firstName = app.application?.student_profile?.first_name_th || ''
      const lastName = app.application?.student_profile?.last_name_th || ''
      const studentId = app.application?.student_profile?.student_id || ''
      const fullName = `${firstName} ${lastName}`.toLowerCase()
      return fullName.includes(term) || studentId.toLowerCase().includes(term)
    })
  }
  
  return result
})

// ========== Methods ==========
async function fetchData() {
  loading.value = true
  try {
    // Fetch qualified applicants
    const [qualified, scheduled] = await Promise.all([
      getAllApplicationScholarships('qualified'),
      getAllApplicationScholarships('interview_scheduled')
    ])
    qualifiedApplicants.value = [...qualified, ...scheduled]
    
    // Fetch existing evaluations to filter out
    const evaluations = await EvaluationService.getAll()
    evaluatedIds.value = new Set(evaluations.map(e => e.application_scholarship_id))
  } catch (e) {
    console.error('Error fetching data:', e)
  } finally {
    loading.value = false
  }
}

async function createEvaluationForApplicant(app: ApplicationScholarshipResponse) {
  if (creatingEvaluation.value) return
  
  // Get interview_round_id from the applicant's booking
  let interviewRoundId: number | null = null
  if (app.interviewe_bookings && app.interviewe_bookings.length > 0) {
    const latestBooking = app.interviewe_bookings[app.interviewe_bookings.length - 1]
    if (latestBooking?.slot?.interview_round_id) {
      interviewRoundId = latestBooking.slot.interview_round_id
    }
  }
  
  if (!interviewRoundId) {
    interviewRoundId = 1 // Fallback
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
      admin_id: adminId.value
    })
    
    await Swal.fire({
      icon: 'success',
      title: 'สร้างการประเมินสำเร็จ',
      timer: 1500,
      showConfirmButton: false
    })
    
    // Refresh data
    await fetchData()
    refreshParentCount()
    emit('created')
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

// ========== Lifecycle ==========
onMounted(async () => {
  // Get admin profile
  try {
    const profileRes: any = await Get('/profile/me')
    if (profileRes && profileRes.role === 'admin' && profileRes.data) {
      adminId.value = profileRes.data.ID
    }
  } catch (e) {
    console.error('Error fetching admin profile:', e)
  }
  
  await fetchData()
})
</script>

<template>
  <div 
    class="w-full mx-auto flex flex-col h-full p-6 bg-white" 
    data-theme="light"
  >
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
      <h2 class="text-xl font-bold text-[#1e3a8a]">ผู้สมัครพร้อมสร้างการประเมิน</h2>

      <div class="flex items-center gap-3 flex-wrap">
        <!-- Scholarship Filter -->
        <div class="relative">
          <select 
            v-model="selectedScholarship"
            class="appearance-none bg-white border border-gray-300 rounded-full px-4 py-2 pr-10 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
          >
            <option value="all">ทุกทุนการศึกษา</option>
            <option v-for="s in scholarshipOptions" :key="s.id" :value="s.id.toString()">{{ s.name }}</option>
          </select>
          <ChevronDown class="w-4 h-4 text-gray-400 absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none" />
        </div>

        <!-- Search Bar -->
        <div class="relative flex-1 md:w-64">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="ค้นหาชื่อ / รหัส..."
            class="w-full pl-10 pr-8 py-2 border border-gray-300 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
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
          @click="fetchData"
          :disabled="loading"
          class="p-2 bg-gray-100 hover:bg-gray-200 rounded-full transition-colors"
          title="รีเฟรช"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-emerald-500 border-r-transparent"></div>
        <p class="mt-2 text-gray-500">กำลังโหลดข้อมูล...</p>
      </div>
    </div>

    <!-- Content -->
    <div v-else class="flex-1 min-h-0 flex flex-col gap-4">
      <div class="overflow-x-auto overflow-y-auto flex-1 rounded-xl border border-slate-200 bg-white shadow-sm">
        <table class="table table-sm w-full" data-theme="light">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-200">
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-left">#</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-left">ชื่อนักศึกษา</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-left hidden md:table-cell">ทุนการศึกษา</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-center">สถานะ</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-gray-600 text-center">การจัดการ</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr 
              v-for="(app, index) in filteredApplicants" 
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
            <tr v-if="filteredApplicants.length === 0">
              <td colspan="5" class="h-[300px] text-center align-middle">
                <div class="flex flex-col items-center justify-center text-gray-400">
                  <UserPlus class="w-12 h-12 mb-3 opacity-50" />
                  <p class="text-gray-500">ไม่มีผู้สมัครที่พร้อมประเมิน</p>
                  <p class="text-sm text-gray-400">ผู้สมัครทุกคนได้รับการประเมินแล้ว</p>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
