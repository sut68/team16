<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ScholarshipAPI } from '@/services/api/scholarship'
import { StatusscholarshipAPI } from '@/services/api/statusscholarship'
import { TypescholarshipAPI } from '@/services/api/typescholarship'

import type {
  ScholarshipResponse,
  ScholarshipStatusResponse,
  ScholarshipTypeResponse
} from '@/interfaces'

// Components
import ScholarshipModal from './ScholarshipModal.vue'
import FeatureScholarshipModal from './featurescholarshipmodal.vue'
import ScholarshipDetailModal from './scholarshipdetailmodal.vue'
import ScholarshipListCard from '@/components/ui/ScholarshipListCard.vue'
import StatsGrid from '@/components/ui/StatsGrid.vue'
import type { StatItem } from '@/components/ui/StatsGrid.vue'
import SearchBar from '@/components/ui/SearchBar.vue'
import Pagination from '@/components/ui/Pagination.vue'

// Icons
import { Plus, ChevronDown, RefreshCw, GraduationCap } from 'lucide-vue-next'
import Swal from 'sweetalert2'

// =====================
// State
// =====================
const scholarships = ref<ScholarshipResponse[]>([])
const statuses = ref<ScholarshipStatusResponse[]>([])
const types = ref<ScholarshipTypeResponse[]>([])

const isLoading = ref(true)
const error = ref<string | null>(null)

const searchQuery = ref('')
const sortOrder = ref<'newest' | 'oldest'>('newest')
const filterStatus = ref('all')
const filterType = ref('all')

// Pagination
const currentPage = ref(1)
const itemsPerPage = ref(10)

// Modal States
const isModalOpen = ref(false)
const selectedScholarship = ref<ScholarshipResponse | null>(null)

const isDetailModalOpen = ref(false)
const detailScholarship = ref<ScholarshipResponse | null>(null)

const isFeatureModalOpen = ref(false)
const currentScholarshipId = ref<number | null>(null)

// =====================
// Fetch Data
// =====================
const fetchScholarships = async () => {
  isLoading.value = true
  error.value = null
  try {
    scholarships.value = await ScholarshipAPI.getAll()
  } catch (err) {
    error.value = 'โหลดข้อมูลทุนไม่สำเร็จ'
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

const fetchStatuses = async () => {
  try {
    statuses.value = await StatusscholarshipAPI.getAll()
  } catch (err) {
    console.error('โหลดสถานะไม่สำเร็จ', err)
  }
}

const fetchTypes = async () => {
  try {
    types.value = await TypescholarshipAPI.getAll()
  } catch (err) {
    console.error('โหลดประเภทไม่สำเร็จ', err)
  }
}

onMounted(() => {
  fetchScholarships()
  fetchStatuses()
  fetchTypes()
})

// =====================
// Computed: Stats
// =====================
const stats = computed(() => {
  const all = scholarships.value
  const open = all.filter(s => 
    s.statusscholarship?.status_name?.toLowerCase().includes('open')
  ).length
  const closed = all.filter(s => 
    s.statusscholarship?.status_name?.toLowerCase().includes('closed')
  ).length
  const sponsors = new Set(all.map(s => s.sponsor?.ID).filter(Boolean)).size
  
  return { total: all.length, open, closed, sponsors }
})

const scholarshipStats = computed<StatItem[]>(() => [
  { 
    title: 'ทุนทั้งหมด', 
    value: stats.value.total, 
    description: 'รายการ',
    icon: 'clipboard', 
    color: 'blue' 
  },
  { 
    title: 'เปิดรับสมัคร', 
    value: stats.value.open, 
    description: 'ทุนที่กำลังเปิด',
    icon: 'check', 
    color: 'green' 
  },
  { 
    title: 'ปิดรับสมัคร', 
    value: stats.value.closed, 
    description: 'ทุนที่ปิดแล้ว',
    icon: 'clock', 
    color: 'orange' 
  },
  { 
    title: 'บริษัทผู้สนับสนุน', 
    value: stats.value.sponsors, 
    description: 'บริษัท',
    icon: 'building', 
    color: 'purple' 
  },
])

// =====================
// Computed: Filtering
// =====================
const filteredScholarships = computed(() => {
  let data = [...scholarships.value]

  // Filter by status
  if (filterStatus.value !== 'all') {
    data = data.filter(
      s => s.statusscholarship?.status_name === filterStatus.value
    )
  }

  // Filter by type
  if (filterType.value !== 'all') {
    data = data.filter(
      s => s.typescholarship?.type_name === filterType.value
    )
  }

  // Filter by search
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase()
    data = data.filter(
      s =>
        s.scholarship_name.toLowerCase().includes(q) ||
        s.description.toLowerCase().includes(q) ||
        s.sponsor?.company_name?.toLowerCase().includes(q)
    )
  }

  // Sort
  data.sort((a, b) =>
    sortOrder.value === 'newest'
      ? new Date(b.open_date).getTime() - new Date(a.open_date).getTime()
      : new Date(a.open_date).getTime() - new Date(b.open_date).getTime()
  )

  return data
})

// =====================
// Computed: Pagination
// =====================
const totalItems = computed(() => filteredScholarships.value.length)
const totalPages = computed(() => Math.ceil(totalItems.value / itemsPerPage.value))

const paginatedScholarships = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredScholarships.value.slice(start, end)
})

const pages = computed(() => {
  const arr: number[] = []
  for (let i = 1; i <= totalPages.value; i++) {
    arr.push(i)
  }
  return arr
})

// Reset page when filters change
watch([searchQuery, filterStatus, filterType], () => {
  currentPage.value = 1
})

// =====================
// Modal Controls
// =====================
const openAddModal = () => {
  selectedScholarship.value = null
  isModalOpen.value = true
}

const openEditModal = (item: ScholarshipResponse) => {
  selectedScholarship.value = item
  isModalOpen.value = true
}

const openDetailModal = (item: ScholarshipResponse) => {
  detailScholarship.value = item
  isDetailModalOpen.value = true
}

const handleDelete = async (item: ScholarshipResponse) => {
  const result = await Swal.fire({
    title: 'ยืนยันการลบ',
    text: `ต้องการลบทุน "${item.scholarship_name}" ใช่หรือไม่?`,
    icon: 'warning',
    showCancelButton: true,
    confirmButtonColor: '#dc2626',
    cancelButtonColor: '#6b7280',
    confirmButtonText: 'ลบ',
    cancelButtonText: 'ยกเลิก'
  })
  
  if (result.isConfirmed) {
    try {
      await ScholarshipAPI.delete(item.ID)
      await Swal.fire({
        icon: 'success',
        title: 'ลบสำเร็จ',
        timer: 1500,
        showConfirmButton: false
      })
      fetchScholarships()
    } catch (err) {
      Swal.fire({
        icon: 'error',
        title: 'ไม่สามารถลบได้',
        text: 'เกิดข้อผิดพลาด'
      })
    }
  }
}

const handleSaved = () => {
  isModalOpen.value = false
  fetchScholarships()
}

// =====================
// Feature Scholarship
// =====================
const openAddFeatureModal = (scholarshipId: number | null) => {
  currentScholarshipId.value = scholarshipId
  isFeatureModalOpen.value = true
}

const handleFeatureSaved = () => {
  isFeatureModalOpen.value = false
}

// =====================
// Pagination Controls
// =====================
function prevPage() {
  if (currentPage.value > 1) currentPage.value--
}

function nextPage() {
  if (currentPage.value < totalPages.value) currentPage.value++
}

function gotoPage(p: number) {
  currentPage.value = p
}

function clearSearch() {
  searchQuery.value = ''
}
</script>

<template>
  <!-- Page Container -->
  <div 
    class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow" 
    data-theme="light"
  >
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold text-gray-900">จัดการทุนการศึกษา</h1>

      <div class="flex items-center gap-3 flex-wrap">
        <!-- Search Bar -->
        <SearchBar 
          v-model="searchQuery" 
          placeholder="ค้นหาชื่อทุน / บริษัท..." 
          @clear="clearSearch" 
        />

        <!-- Status Filter -->
        <div class="relative">
          <select 
            v-model="filterStatus"
            class="appearance-none bg-white border border-gray-300 rounded-full px-4 py-2 pr-10 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">สถานะทั้งหมด</option>
            <option v-for="st in statuses" :key="st.ID" :value="st.status_name">
              {{ st.status_name }}
            </option>
          </select>
          <ChevronDown class="w-4 h-4 text-gray-400 absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none" />
        </div>

        <!-- Type Filter -->
        <div class="relative">
          <select 
            v-model="filterType"
            class="appearance-none bg-white border border-gray-300 rounded-full px-4 py-2 pr-10 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="all">ทุกประเภท</option>
            <option v-for="tp in types" :key="tp.ID" :value="tp.type_name">
              {{ tp.type_name }}
            </option>
          </select>
          <ChevronDown class="w-4 h-4 text-gray-400 absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none" />
        </div>

        <!-- Sort -->
        <div class="relative">
          <select 
            v-model="sortOrder"
            class="appearance-none bg-white border border-gray-300 rounded-full px-4 py-2 pr-10 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="newest">ใหม่ล่าสุด</option>
            <option value="oldest">เก่าที่สุด</option>
          </select>
          <ChevronDown class="w-4 h-4 text-gray-400 absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none" />
        </div>

        <!-- Refresh -->
        <button 
          @click="fetchScholarships"
          :disabled="isLoading"
          class="p-2 bg-gray-100 hover:bg-gray-200 rounded-full transition-colors"
          title="รีเฟรช"
        >
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': isLoading }" />
        </button>

        <!-- Add Button -->
        <button 
          @click="openAddModal" 
          class="btn-ghost-rounded"
        >
          <Plus class="w-4 h-4" />
          <span class="font-medium">เพิ่มทุนใหม่</span>
        </button>
      </div>
    </div>

    <!-- Stats Grid -->
    <StatsGrid 
      :stats="scholarshipStats" 
      :columns="4"
      class="mb-6"
    />

    <!-- Loading State -->
    <div v-if="isLoading" class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-r-transparent"></div>
        <p class="mt-2 text-gray-500">กำลังโหลดข้อมูล...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="p-4 bg-red-100 text-red-700 rounded-lg">
      {{ error }}
    </div>

    <!-- Content -->
    <div v-else class="flex-1 min-h-0 flex flex-col gap-4">
      <!-- Card List -->
      <div class="overflow-y-auto flex-1 space-y-4 pr-2">
        <TransitionGroup name="list" tag="div" class="space-y-4">
          <ScholarshipListCard
            v-for="item in paginatedScholarships"
            :key="item.ID"
            :scholarship="item"
            @view="openDetailModal"
            @edit="openEditModal"
            @manage-features="openAddFeatureModal"
            @delete="handleDelete"
          />
        </TransitionGroup>

        <!-- Empty State -->
        <div 
          v-if="paginatedScholarships.length === 0" 
          class="flex flex-col items-center justify-center py-16 text-gray-400"
        >
          <div class="bg-gray-100 p-5 rounded-full mb-4">
            <GraduationCap class="w-12 h-12 opacity-50" />
          </div>
          <p v-if="searchQuery" class="text-gray-500">
            ไม่พบข้อมูลที่ตรงกับ "<strong class="text-gray-700">{{ searchQuery }}</strong>"
          </p>
          <p v-else class="text-gray-500">ยังไม่มีทุนการศึกษาในระบบ</p>
          <button 
            @click="openAddModal" 
            class="btn-ghost-rounded mt-4"
          >
            <Plus class="w-4 h-4" />
            <span>เพิ่มทุนแรก</span>
          </button>
        </div>
      </div>

      <!-- Pagination -->
      <Pagination
        v-if="totalItems > 0"
        :page="currentPage"
        :total-pages="totalPages"
        :pages="pages"
        :total="totalItems"
        :per-page="itemsPerPage"
        :paged-length="paginatedScholarships.length"
        @prev="prevPage"
        @next="nextPage"
        @goto="gotoPage"
      />
    </div>
  </div>

  <!-- Modals -->
  <ScholarshipModal
    v-if="isModalOpen"
    :isOpen="isModalOpen"
    :scholarship="selectedScholarship"
    @close="isModalOpen = false"
    @saved="handleSaved"
  />

  <FeatureScholarshipModal
    v-if="isFeatureModalOpen"
    :isOpen="isFeatureModalOpen"
    :scholarshipId="currentScholarshipId"
    @close="isFeatureModalOpen = false"
    @saved="handleFeatureSaved"
  />

  <ScholarshipDetailModal
    :isOpen="isDetailModalOpen"
    :scholarship="detailScholarship"
    @close="isDetailModalOpen = false"
  />
</template>

<style scoped>
/* Ghost Rounded Button (matching SponsorList) */
.btn-ghost-rounded {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0 1.25rem;
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

/* List transition */
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

/* Scrollbar */
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: rgba(0,0,0,0.15);
  border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: rgba(0,0,0,0.25);
}
</style>
