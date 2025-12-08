<script setup lang="ts">
  import { ref, computed, onMounted, watch } from 'vue'
  import { SponsorService } from '../../../services/sponsor/sponsor';
  import type { SponsorPayload, SponsorResponse, SponsorView } from '../../../interfaces/sponsor';
  import SponsorForm from './SponsorForm.vue';
  import { ExternalLink, Plus, ChevronLeft, ChevronRight } from "lucide-vue-next";
  import SponsorActionMenu from '../../../components/ui/SponsorActionMenu.vue';

  const sponsors = ref<SponsorView[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const isCreateSponsorOpen = ref(false)
  const creating = ref(false)
  const loadingToggle = ref<number | null>(null);

  // search and page
  const q = ref("")
  const searchQuery = ref("")
  const page = ref(1)
  const perPage = ref<number>(10)

  // debounce timer
  let _searchTimer: ReturnType<typeof setTimeout> | null = null
  watch(searchQuery, (val) => {
    if (_searchTimer) clearTimeout(_searchTimer)
    _searchTimer = setTimeout(() => {
      q.value = val.trim().toLowerCase()
      page.value = 1
      _searchTimer = null
    }, 300)
  })

  watch(q, () => {
  page.value = 1
  })

  function clearSearch() {
    searchQuery.value = ''
    q.value = ''
  }

  // map API response
  function toSponsorView(s: SponsorResponse): SponsorView {
    return {
      ID:               s.ID,
      company_name:     s.company_name,
      website:          s.website ?? null,
      industry_name:    s.industry?.name ?? null,
      industry_id:      s.industry_id ?? s.industry?.ID ?? null,
      status:           s.status,
      contacts_count:   Array.isArray(s.contacts) ? s.contacts.length : 0,
    }
  }

  const filtered = computed(() => {
    const term = (q.value || '').trim().toLowerCase()
    if (!term) return sponsors.value
    return sponsors.value.filter(s =>
      (s.company_name || '').toString().toLowerCase().includes(term) ||
      (s.website || '').toString().toLowerCase().includes(term) ||
      (s.industry_name || '').toString().toLowerCase().includes(term)
    )
  })

  const total = computed(() => filtered.value.length)
  const totalPages = computed(() => Math.max(1, Math.ceil(total.value / perPage.value)))

  const paged = computed(() => {
    const start = (page.value - 1) * perPage.value
    return filtered.value.slice(start, start + perPage.value)
  })

  const totalSponsors = computed(() => sponsors.value.length)

  const activeSponsors = computed(() => 
    sponsors.value.filter(s => s.status === "active").length
  )

  const totalContacts = computed(() =>
    sponsors.value.reduce((sum, s) => sum + (s.contacts_count ?? 0), 0)
  )

  watch(totalPages, (tp) => {
    if (page.value > tp) page.value = tp
  })

  watch(perPage, () => {
    page.value = 1
  })

  const pages = computed(() => {
    const tp = totalPages.value
    const maxButtons = 7
    if (tp <= maxButtons) return Array.from({length: tp}, (_, i) => i + 1)

    const half = Math.floor(maxButtons / 2)
    let start = page.value - half
    let end = page.value + half

    if (start < 1) {
      start = 1
      end = maxButtons
    } else if (end > tp) {
      end = tp
      start = tp - maxButtons + 1
    }

    return Array.from({ length: end - start + 1 }, (_, i) => start + i)
  })

  async function fetchSponsors() {
    loading.value = true
    error.value = null
    try {
      const res = await SponsorService.getAll()
      sponsors.value = Array.isArray(res) ? res.map((s: SponsorResponse) => toSponsorView(s)) : []
      page.value = 1

    } catch (err: any) {
      console.error(err)
      error.value = err?.response?.data?.message || err?.message || "โหลดข้อมูล Sponsors ไม่สำเร็จ"

    } finally {
      loading.value = false

    }
  }

  async function onToggleStatus(id:number, currentStatus?: string) {
    const idx = sponsors.value.findIndex(s => s.ID === id);
    if (idx === -1) return;

    const prev = sponsors.value[idx];
    if (!prev) return;
    const prevStatus = prev.status;
    const newStatus = currentStatus === 'active' ? 'inactive' : 'active';

    sponsors.value[idx] = { ...prev, status: newStatus };
    loadingToggle.value = id;

    try {
      const payload = { status: newStatus };
      const updated = await SponsorService.update(id, payload);
      if (updated && updated.status) {
        sponsors.value[idx] = { ...sponsors.value[idx], status: updated.status };
      }
    } catch (err) {
      sponsors.value[idx] = { ...sponsors.value[idx], status: prevStatus };
      console.error('Toggle status failed', err);
      alert('เปลี่ยนสถานะไม่สำเร็จ ลองใหม่อีกครั้ง');
    } finally {
      loadingToggle.value = null;
    }
  }

  async function removeOne(id: number) {
    if (!confirm('ลบบริษัทนี้จริงหรือไม่? การกระทำนี้ไม่สามารถย้อนกลับได้')) return
    loading.value = true
    error.value = null
    const snapshot = [...sponsors.value]
    sponsors.value = sponsors.value.filter(s => s.ID !== id)
    
    try {
      await SponsorService.delete(id)
      alert('ลบสำเร็จ')

      await fetchSponsors()

    } catch (err: any) {
      console.error(err)
      sponsors.value = snapshot

      alert(err?.response?.data?.message || 'ลบไม่สำเร็จ')

    } finally {
      loading.value = false

    }
  }

  // เปิด Sponsor Form
  function openCreateForm() {
    isCreateSponsorOpen.value = true
  }

  async function onCreateSponsor(payload: SponsorPayload) {
    creating.value = true
    try {
      const created = await SponsorService.create(payload)

      const view = toSponsorView(created)

      sponsors.value.unshift(view)

      isCreateSponsorOpen.value = false

      page.value = 1
    } catch (err: any) {
      console.error(err)
      alert(err?.response?.data?.message || 'เพิ่มบริษัทไม่สำเร็จ')
    } finally {
      creating.value = false
    }
  }

  function prevPage() {
    if (page.value > 1) page.value--
  }
  function nextPage() {
    if (page.value < totalPages.value) page.value++
  }

  onMounted(() => {
    fetchSponsors()
  })
</script>

<template>
  <div class="w-full mx-auto flex flex-col h-full p-6 bg-white rounded-tl-[30px] shadow" data-theme="light">

    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <h1 class="text-2xl font-bold text-[#1e3a8a]">บริษัทที่ให้ทุน</h1>

      <div class="flex items-center gap-3 w-full md:w-auto">
        <!-- Search Bar -->
        <div class="relative flex-1 md:flex-none w-full md:w-96">
          <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400 pointer-events-none">
            <!-- icon -->
            <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </span>

          <input
            type="text"
            v-model="searchQuery"
            placeholder="ค้นหาชื่อ / เว็บไซต์..."
            class="input input-bordered input-sm h-10 w-full rounded-full pl-10 bg-white border-gray-300 shadow-sm focus:border-[#1e3a8a] focus:ring-1 focus:ring-[#1e3a8a] text-sm"
            aria-label="Search sponsors"
          />

          <button
            v-if="searchQuery"
            @click="clearSearch"
            class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-700"
            aria-label="Clear search"
          >
            <!-- icon -->
            <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- ปุ่มที่ใช้ไปหน้าเพิ่ม Sponsors -->
        <button
          @click="openCreateForm"
          class="btn btn-sm bg-white border border-gray-300 text-gray-700
                hover:bg-gray-100 hover:border-gray-400
                flex items-center justify-center gap-2
                rounded-full px-5 h-10 w-full md:w-auto shadow-sm
                transition-all duration-150"
        >
          <Plus class="w-4 h-4"/>
          <span class="font-medium">เพิ่มบริษัท</span>
        </button>
      </div>

    </div>

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
      <div class="stats shadow">
        <!-- Total Sponsors -->
        <div class="stat rounded-2xl border bg-white">
          <div class="text-sm text-gray-500 flex items-center gap-2">
            <span class="w-2 h-2 bg-purple-400 rounded-full"></span>
            Total Sponsors
          </div>
          <div class="text-2xl font-semibold mt-2">
            {{ totalSponsors }}
          </div>
        </div>
      </div>

      <!-- Active Sponsors -->
      <div class="stats shadow">
        <div class="stat rounded-2xl border bg-white">
          <div class="text-sm text-gray-500 flex items-center gap-2">
            <span class="w-2 h-2 bg-green-500 rounded-full"></span>
            Active Sponsor
          </div>
          <div class="text-2xl font-semibold mt-2">
            {{ activeSponsors }}
          </div>
        </div>
      </div>

      <!-- Total Contacts -->
      <div class="stats shadow">
        <div class="stat rounded-2xl border bg-white">
          <div class="text-sm text-gray-500 flex items-center gap-2">
            <span class="w-2 h-2 bg-blue-500 rounded-full"></span>
            Total Contacts
          </div>
          <div class="text-2xl font-semibold mt-2">
            {{ totalContacts }}
          </div>
        </div>
      </div>

    </div>

    <!-- Loading หรือ Error -->
    <div v-if="loading" class="p-6 w-full">
      <!-- simple skeleton -->
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

    <div v-if="error" class="p-4 alert alert-error mb-4">{{ error }}</div>

    <div v-if="!loading && !error" class="flex-1 min-h-0 flex flex-col gap-4">

      <!-- Table wrapper -->
      <div
        class="table-scroll overflow-x-auto overflow-y-auto min-h-0 bg-white rounded flex-1"
        role="region"
        aria-label="Sponsors table"
        tabindex="0"
      >
        <table class="table w-full">

          <!-- Header Table -->
          <colgroup>
            <col style="width:4%" />
            <col style="width:30%" />
            <col style="width:24%" />
            <col style="width:16%" />
            <col style="width:10%" />
            <col style="width:8%" />
            <col style="width:8%" />
          </colgroup>

          <thead class="bg-white/95">
            <tr class="sticky top-0 z-10" style="box-shadow: 0 2px 6px rgba(0,0,0,0.04);">
              <th class="px-4 py-3 text-left">#</th>
              <th class="px-4 py-3 text-left">ชื่อบริษัท</th>
              <th class="px-4 py-3 text-left">เว็บไซต์</th>
              <th class="px-4 py-3 text-left">อุตสาหกรรม</th>
              <th class="px-4 py-3 text-left">สถานะ</th>
              <th class="px-4 py-3 text-center">Contacts</th>
              <th class="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>

          <!-- Body Table -->
          <tbody>
            <tr v-for="s in paged" :key="s.ID" class="odd:bg-gray-50 hover:bg-gray-100 focus-within:bg-gray-100">
              <td class="px-4 py-3 text-sm text-gray-600">{{ s.ID }}</td>

              <td class="px-4 py-3">
                <router-link 
                  :to="`/sponsors/${s.ID}/profile`" 
                  class="font-semibold text-gray-800 hover:text-blue-600 transition-colors duration-150 block truncate"
                  :title="s.company_name"
                >
                  {{ s.company_name }}
                </router-link>
              </td>

              <td class="px-4 py-3 text-sm">
                <div class="truncate" :title="s.website ?? '-'">
                  <a
                    v-if="s.website"
                    :href="s.website"
                    target="_blank"
                    rel="noopener"
                    class="inline-flex items-center gap-1 text-blue-600 hover:text-blue-800 hover:underline"
                  >
                    {{ s.website }}
                    <ExternalLink class="w-4 h-4 opacity-60"/> 
                  </a>
                  <span v-else class="text-sm text-gray-500">-</span>
                </div>
              </td>

              <td class="px-4 py-3 text-sm text-gray-700 truncate" :title="s.industry_name ?? '-'">
                {{ s.industry_name ?? '-' }}
              </td>

              <td class="px-4 py-3">
                <span
                  :class="[
                    'px-2 text-sm font-medium rounded-md',
                    s.status === 'active'
                      ? 'text-green-600 bg-green-100'
                      : 'text-red-600 bg-red-100'
                  ]"
                >
                  {{ s.status }}
                </span>
              </td>

              <td class="px-4 py-3 text-center text-sm">{{ s.contacts_count ?? 0 }}</td>

              <td class="px-4 py-3 text-right whitespace-nowrap">
                <SponsorActionMenu 
                  :id="s.ID"
                  :status="s.status"
                  :disabled="loadingToggle === s.ID"
                  @edit="(id) => $router.push(`/sponsor/${id}/edit`)"
                  @toggle-status="() => onToggleStatus(s.ID, s.status)"
                  @delete="removeOne"
                />
              </td>
            </tr>

            <tr v-if="paged.length === 0 && total > 0">
              <td colspan="7" class="px-4 py-8 text-center text-sm text-gray-500">
                หน้า {{ page }} ไม่มีรายการ — ลองไปหน้าก่อนหน้านี้
              </td>
            </tr>
          </tbody>

        </table>
      </div>
      
      <!-- empty state -->
      <div v-if="total === 0" :key="'empty'" class="flex flex-col items-center py-12 text-gray-400">
        <div class="bg-white/50 p-4 rounded-full mb-3">
          <!-- icon -->
          <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 opacity-50" fill="none" viewBox="0 0 24 24"
            stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>

        <p v-if="searchQuery">ไม่พบข้อมูลที่ตรงกับ "{{ searchQuery }}"</p>
        <p v-else>ไม่พบรายการในรอบที่เลือก</p>
      </div>

      <div class="flex flex-col md:flex-row items-center justify-between gap-3 mt-2">
        
        <div class="text-sm">
          แสดง
          <strong>{{ (page-1)*perPage + (paged.length ? 1 : 0) }}</strong>
          -
          <strong>{{ Math.min(page*perPage, total) }}</strong>
          จาก <strong>{{ total }}</strong> รายการ
        </div>
  
        <nav class="flex items-center gap-2" aria-label="Pagination">
          <button 
            class="btn btn-sm bg-white border border-gray-300 text-gray-700
                hover:bg-gray-100 hover:border-gray-400
                flex items-center justify-center gap-2
                rounded-full px-5 h-10 w-full md:w-auto shadow-sm
                transition-all duration-150"
            :disabled="page <= 1" @click="prevPage" 
            aria-label="Previous page"
          >  
            <ChevronLeft class="w-4 h-4"/>
            <span class="font-medium">Prev</span>
          </button>

          <template v-for="p in pages" :key="p">
            <button
              class="btn btn-sm bg-white border border-gray-300 text-gray-700
                hover:bg-gray-100 hover:border-gray-400
                flex items-center justify-center gap-2
                rounded-full px-5 h-10 w-full md:w-auto shadow-sm
                transition-all duration-150"
              :class="p === page ? 'btn-primary' : ''"
              @click="page = p"
              :aria-current="p === page ? 'page' : undefined"
            >{{ p }}</button>
          </template>

          <button 
            class="btn btn-sm bg-white border border-gray-300 text-gray-700
                hover:bg-gray-100 hover:border-gray-400
                flex items-center justify-center gap-2
                rounded-full px-5 h-10 w-full md:w-auto shadow-sm
                transition-all duration-150" 
            :disabled="page >= totalPages" 
            @click="nextPage" 
            aria-label="Next page"
          >
            <span class="font-medium">Next</span>
            <ChevronRight class="w-4 h-4"/>
          </button>
        </nav>
      </div>

    </div>
  </div>

  <SponsorForm
    v-model:isOpen="isCreateSponsorOpen"
    :loading="creating"
    @create="onCreateSponsor"
  />
</template>

<style scoped>
  *,
  *::before,
  *::after {
    box-sizing: border-box;
  }

  .table-scroll {
    scrollbar-gutter: stable both-edges;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
  }

  thead th {
    position: sticky;
    top: 0;
    background: white;
    z-index: 10;
    box-shadow: 0 2px 6px rgba(0,0,0,0.04);
  }

  td, th {
    padding-top: 0.75rem;
    padding-bottom: 0.75rem;
    padding-left: 0.75rem;
    padding-right: 0.75rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  tbody tr:hover {
    background-color: #f5f7fa;
  }

  .table-scroll::-webkit-scrollbar {
    width: 10px;
    height: 10px;
  }
  .table-scroll::-webkit-scrollbar-thumb {
    background: rgba(0,0,0,0.12);
    border-radius: 999px;
  }

  @media (max-width: 768px) {
    table {
      table-layout: auto;
    }
    th, td {
      white-space: normal;
      word-break: break-word;
    }
    .truncate { white-space: normal; text-overflow: unset; overflow: visible; }
  }

  @media (max-width: 640px) {
    th:nth-child(6), td:nth-child(6),
    th:nth-child(7), td:nth-child(7) {
      display: none;
    }
  }

</style>
