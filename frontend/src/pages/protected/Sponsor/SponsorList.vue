<script setup lang="ts">
import { onMounted, watch, computed } from 'vue'
import type { SponsorPayload } from '@/interfaces/sponsor'
import { SponsorService } from '@/services/sponsor/sponsor'

// Hooks
import { useSponsorList } from '@/hooks/sponsor/useSponsorList'
import { usePagination } from '@/hooks/sponsor/usePagination'
import { useSearch } from '@/hooks/sponsor/useSearch'
import { useSponsorModals } from '@/hooks/sponsor/useSponsorModal'

// Components
import SponsorForm from './SponsorForm.vue'
import SponsorEdit from './SponsorEdit.vue'
import SponsorContact from './SponsorContact.vue'
import SponsorActionMenu from '@/components/ui/SponsorActionMenu.vue'
import StatsGrid from '@/components/ui/StatsGrid.vue'
import type { StatItem } from '@/components/ui/StatsGrid.vue'
import SearchBar from '@/components/ui/SearchBar.vue'
import Pagination from '@/components/ui/Pagination.vue'

// Icons
import { ExternalLink, Plus } from 'lucide-vue-next'
import Swal from 'sweetalert2'

// ========== Hooks ==========
const {
  sponsors,
  loading,
  error,
  loadingToggle,
  totalSponsors,
  activeSponsors,
  totalContacts,
  fetchSponsors,
  toggleStatus,
  removeOne,
  addSponsor,
  updateContactsCount,
} = useSponsorList()

const { searchQuery, q, filtered, clearSearch } = useSearch(
  sponsors,
  (s, term) =>
    (s.company_name || '').toLowerCase().includes(term) ||
    (s.website || '').toLowerCase().includes(term) ||
    (s.industry_name || '').toLowerCase().includes(term)
)

const {
  page,
  perPage,
  total,
  totalPages,
  paged,
  pages,
  prevPage,
  nextPage,
  resetPage,
} = usePagination(filtered)

const {
  isCreateSponsorOpen,
  isEditSponsorOpen,
  isContactsOpen,
  creating,
  updating,
  currentSponsor,
  currentSponsorId,
  currentContacts,
  serverErrors,
  openCreateForm,
  openEditForm,
  openContacts,
} = useSponsorModals()

// Reset page when search changes
watch(q, () => resetPage())

// Computed stats for StatsGrid
const sponsorStats = computed<StatItem[]>(() => [
  { 
    title: 'Total Sponsors', 
    value: totalSponsors.value, 
    description: 'บริษัททั้งหมด',
    icon: 'building', 
    color: 'blue' 
  },
  { 
    title: 'Active Sponsors', 
    value: activeSponsors.value, 
    description: 'กำลังดำเนินการ',
    icon: 'check', 
    color: 'green' 
  },
  { 
    title: 'Total Contacts', 
    value: totalContacts.value, 
    description: 'ผู้ติดต่อทั้งหมด',
    icon: 'users', 
    color: 'orange' 
  },
])

// ========== Event Handlers ==========
function handleGotoPage(p: number) {
  page.value = p
}

async function onCreateSponsor(payload: SponsorPayload) {
  creating.value = true
  try {
    const created = await SponsorService.create(payload)
    addSponsor(created)
    await Swal.fire({ icon: 'success', title: 'สร้างบริษัทสำเร็จ' })
    isCreateSponsorOpen.value = false
    resetPage()
  } catch (err: any) {
    console.error('onCreateSponsor error:', err)
    const resp = err?.response?.data
    if (resp?.errors) {
      await Swal.fire({
        icon: 'warning',
        title: 'ข้อมูลไม่ผ่านการตรวจสอบ',
        text: resp.message ?? 'กรุณาตรวจสอบข้อผิดพลาดในฟอร์ม',
      })
    } else {
      await Swal.fire({
        icon: 'error',
        title: 'ผิดพลาด',
        text: resp?.message ?? err?.message ?? 'เกิดข้อผิดพลาด',
      })
    }
  } finally {
    creating.value = false
  }
}

async function onUpdateSponsor(payload: SponsorPayload) {
  if (!currentSponsor.value) return
  updating.value = true
  serverErrors.value = null

  try {
    const sponsorFields = {
      company_name: payload.company_name,
      website: payload.website,
      industry_id: payload.industry_id,
      status: payload.status,
      description: payload.description,
    }
    await SponsorService.update(currentSponsor.value.ID, sponsorFields)

    // Contacts upsert/delete logic
    const orig = currentSponsor.value.contacts ?? []
    const curr = payload.contacts ?? []
    const origIDs = orig.map(c => c.ID)
    const currIDs = curr.filter(c => (c as any).ID).map(c => (c as any).ID)
    const upsert = curr.map(c => ({
      ...((c as any).ID ? { ID: (c as any).ID } : {}),
      name: c.name,
      email: c.email,
      phone: c.phone,
      position: c.position,
    }))
    const delete_ids = origIDs.filter(id => !currIDs.includes(id))

    if (upsert.length || delete_ids.length) {
      await SponsorService.updateContacts(currentSponsor.value.ID, { upsert, delete_ids })
    }

    await Swal.fire({ icon: 'success', title: 'อัปเดตสำเร็จ' })
    await fetchSponsors()
    isEditSponsorOpen.value = false
  } catch (err: any) {
    console.error(err)
    const resp = err?.response?.data
    if (resp?.errors) {
      serverErrors.value = resp.errors
      Swal.fire({
        icon: 'warning',
        title: 'ข้อมูลไม่ผ่านการตรวจสอบ',
        text: resp.message ?? 'ตรวจสอบข้อมูลที่กรอก',
      })
    } else {
      Swal.fire({ icon: 'error', title: 'ผิดพลาด', text: resp?.message ?? 'เกิดข้อผิดพลาด' })
    }
  } finally {
    updating.value = false
  }
}

async function onContactsSaved(newContacts: any[]) {
  isContactsOpen.value = false
  if (currentSponsor.value?.ID === currentSponsorId.value) {
    currentSponsor.value.contacts = newContacts
  }
  if (currentSponsorId.value) {
    updateContactsCount(currentSponsorId.value, newContacts.length)
  }
  await Swal.fire({ icon: 'success', title: 'บันทึกผู้ติดต่อเรียบร้อย' })
}

function onToggleStatus(id: number, status?: string) {
  toggleStatus(id, status)
}

// ========== Lifecycle ==========
onMounted(fetchSponsors)
</script>

<template>
  <div 
    class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow" 
    data-theme="light"
    data-testid="sponsor-list-page"
  >

    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6" data-testid="sponsor-list-header">
      <h1 class="text-2xl font-bold text-gray-900" data-testid="sponsor-list-title">บริษัทที่ให้ทุน</h1>

      <div class="flex items-center gap-3 w-full md:w-auto">
        <!-- Search Bar -->
        <SearchBar 
          v-model="searchQuery" 
          placeholder="ค้นหาชื่อ / เว็บไซต์..." 
          @clear="clearSearch" 
        />

        <!-- Add Button -->
        <button 
          @click="openCreateForm" 
          class="btn-ghost-rounded w-full md:w-auto"
          data-testid="btn-add-sponsor"
        >
          <Plus class="w-4 h-4"/>
          <span class="font-medium">เพิ่มบริษัท</span>
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <StatsGrid 
      :stats="sponsorStats" 
      :columns="3"
      class="mb-6"
      data-testid="sponsor-stats-container"
    />

    <!-- Loading State -->
    <div v-if="loading" class="p-6 w-full" data-testid="sponsor-loading-state">
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
    <div v-if="error" class="p-4 alert alert-error mb-4" data-testid="sponsor-error-state">{{ error }}</div>

    <!-- Content -->
    <div v-if="!loading && !error" class="flex-1 min-h-0 flex flex-col gap-4" data-testid="sponsor-content">

      <!-- Table -->
      <div
        class="overflow-x-auto overflow-y-auto flex-1 min-h-[400px] rounded-xl border border-slate-200 bg-white shadow-sm"
        role="region"
        aria-label="Sponsors table"
        data-testid="sponsor-table-container"
      >
        <table class="table table-sm w-full" data-testid="sponsor-table">
          <!-- Header -->
          <thead>
            <tr class="bg-slate-800 text-white">
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">#</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">ชื่อบริษัท</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 hidden md:table-cell">เว็บไซต์</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 hidden lg:table-cell">อุตสาหกรรม</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">สถานะ</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-center hidden sm:table-cell">Contacts</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-right">Actions</th>
            </tr>
          </thead>

          <tbody class="divide-y divide-slate-100">
            <tr 
              v-for="(s, index) in paged" 
              :key="s.ID" 
              class="transition-all duration-200 hover:bg-blue-50/50 hover:shadow-sm"
              :class="index % 2 === 0 ? 'bg-white' : 'bg-slate-50/50'"
              :data-testid="`sponsor-row-${s.ID}`"
            >
              <td class="py-3 px-4 text-sm text-slate-500 font-mono">{{ s.ID }}</td>

              <td class="py-3 px-4">
                <router-link 
                  :to="{ name: 'SponsorsProfile', params: { id: s.ID }}"
                  class="font-semibold text-slate-800 hover:text-blue-600 transition-colors duration-150 block truncate max-w-xs"
                  :title="s.company_name"
                  :data-testid="`sponsor-link-${s.ID}`"
                >
                  {{ s.company_name }}
                </router-link>
              </td>

              <td class="py-3 px-4 text-sm hidden md:table-cell">
                <div class="truncate max-w-[200px]" :title="s.website ?? '-'">
                  <a
                    v-if="s.website"
                    :href="s.website"
                    target="_blank"
                    rel="noopener"
                    class="inline-flex items-center gap-1 text-blue-600 hover:text-blue-800 hover:underline"
                  >
                    {{ s.website }}
                    <ExternalLink class="w-3.5 h-3.5 opacity-60 flex-shrink-0"/> 
                  </a>
                  <span v-else class="text-slate-400">-</span>
                </div>
              </td>

              <td class="py-3 px-4 text-sm text-slate-600 truncate max-w-[150px] hidden lg:table-cell" :title="s.industry_name ?? '-'">
                {{ s.industry_name ?? '-' }}
              </td>

              <td class="py-3 px-4">
                <span
                  :class="[
                    'inline-flex items-center px-2.5 py-1 text-xs font-semibold rounded-md shadow-sm',
                    s.status === 'active'
                      ? 'text-emerald-800 bg-emerald-100 ring-1 ring-emerald-200'
                      : 'text-rose-800 bg-rose-100 ring-1 ring-rose-200'
                  ]"
                >
                  <span 
                    class="w-1.5 h-1.5 rounded-full mr-1.5 animate-pulse"
                    :class="s.status === 'active' ? 'bg-emerald-500' : 'bg-rose-500'"
                  ></span>
                  {{ s.status === 'active' ? 'Active' : 'Inactive' }}
                </span>
              </td>

              <td class="py-3 px-4 text-center text-sm hidden sm:table-cell">
                <span class="inline-flex items-center justify-center min-w-[2rem] px-2.5 py-1 bg-slate-100 text-slate-700 rounded-md text-xs font-semibold">
                  {{ s.contacts_count ?? 0 }}
                </span>
              </td>

              <td class="py-3 px-4 text-right">
                <SponsorActionMenu 
                  :id="s.ID"
                  :status="s.status"
                  :disabled="loadingToggle === s.ID"
                  @edit="openEditForm"
                  @edit-contacts="openContacts"
                  @toggle-status="() => onToggleStatus(s.ID, s.status)"
                  @delete="removeOne"
                />
              </td>
            </tr>

            <!-- No data in current page -->
            <tr v-if="paged.length === 0 && total > 0">
              <td colspan="7" class="py-12 text-center">
                <div class="text-slate-400">
                  <svg class="w-12 h-12 mx-auto mb-3 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p class="text-sm">หน้า {{ page }} ไม่มีรายการ</p>
                  <p class="text-xs mt-1">ลองไปหน้าก่อนหน้านี้</p>
                </div>
              </td>
            </tr>

            <!-- Empty State Row -->
            <tr v-if="total === 0" data-testid="sponsor-empty-state" class="border-0 hover:bg-transparent">
              <td colspan="7" class="h-[400px] border-0 align-middle">
                <div class="flex flex-col items-center justify-center h-full text-gray-400">
                  <div class="bg-gray-100 p-5 rounded-full mb-4">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 opacity-50" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                        d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                    </svg>
                  </div>
                  <p v-if="searchQuery" class="text-gray-500">
                    ไม่พบข้อมูลที่ตรงกับ "<strong class="text-gray-700">{{ searchQuery }}</strong>"
                  </p>
                  <p v-else class="text-gray-500">ยังไม่มีบริษัทในระบบ</p>
                  <button 
                    @click="openCreateForm" 
                    class="btn-ghost-rounded mt-4"
                    data-testid="btn-add-first-sponsor"
                  >
                    <Plus class="w-4 h-4"/>
                    <span>เพิ่มบริษัทแรก</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- Pagination -->
      <Pagination
        v-if="total > 0"
        :page="page"
        :total-pages="totalPages"
        :pages="pages"
        :total="total"
        :per-page="perPage"
        :paged-length="paged.length"
        @prev="prevPage"
        @next="nextPage"
        @goto="handleGotoPage"
      />
    </div>
  </div>

  <!-- Modals -->
  <SponsorForm
    v-model:isOpen="isCreateSponsorOpen"
    :loading="creating"
    @create="onCreateSponsor"
    data-testid="sponsor-create-modal"
  />

  <SponsorEdit
    v-model:isOpen="isEditSponsorOpen"
    :loading="updating"
    :initialData="currentSponsor"
    :serverErrors="serverErrors"
    @update="onUpdateSponsor"
    @delete="removeOne"
    data-testid="sponsor-edit-modal"
  />

  <SponsorContact
    v-if="currentSponsorId !== null"
    v-model:isOpen="isContactsOpen"
    :sponsorId="currentSponsorId"
    :sponsorName="currentSponsor?.company_name ?? ''"
    :initialContacts="currentContacts"
    @saved="onContactsSaved"
    @close="() => { isContactsOpen = false }"
    data-testid="sponsor-contact-modal"
  />
</template>

<style scoped>
*,
*::before,
*::after {
  box-sizing: border-box;
}

/* Button Style */
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

.btn-ghost-rounded.btn-primary {
  background-color: #1e3a8a;
  color: white;
  border-color: #1e3a8a;
}

/* Table Styles */
.table-scroll {
  scrollbar-gutter: stable both-edges;
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

td, th {
  padding-top: 0.75rem;
  padding-bottom: 0.75rem;
  padding-left: 1rem;
  padding-right: 1rem;
}

/* Scrollbar */
.table-scroll::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.table-scroll::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.table-scroll::-webkit-scrollbar-thumb {
  background: rgba(0,0,0,0.15);
  border-radius: 4px;
}

.table-scroll::-webkit-scrollbar-thumb:hover {
  background: rgba(0,0,0,0.25);
}
</style>