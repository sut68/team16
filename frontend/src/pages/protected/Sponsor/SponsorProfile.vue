<script setup lang="ts">
  import SponsorService from '@/services/sponsor/sponsor';
  import type { ContactResponse, SponsorResponse, SponsorScholarshipResponse } from '../../../interfaces/sponsor';
  import { computed, ref, watch } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { Building2, Globe, GraduationCap, Calendar, Tag, ArrowLeft, User, Mail, Phone, Pencil } from 'lucide-vue-next';
  import SponsorContact from './SponsorContact.vue';
  import Swal from 'sweetalert2';

  const route = useRoute();
  const router = useRouter();

  const goBack = () => {
    router.push('/admin/sponsors')
  }

  const sponsorId = computed<number | null>(() => {
    const param = route.params.id
    if (!param) return null

    const value = Array.isArray(param) ? param[0] : param
    const id = Number(value)

    return Number.isNaN(id) ? null : id
  })

  const sponsor = ref<SponsorResponse | null>(null)
  const scholarships = ref<SponsorScholarshipResponse[]>([])
  const loading = ref(true)
  const loadingScholarships = ref(false)
  const error = ref<string | null>(null)

  // Data Fetching ขอข้อมูลจาก Backend
  const fetchSponsor = async () => {
    if (sponsorId.value === null) {
      error.value = "Invalid Sponsor ID in URL."
      await Swal.fire({ icon: 'error', title: 'ข้อผิดพลาด', text: 'ไม่พบ Sponsor ID ใน URL' })
      return
    }

    loading.value = true
    error.value = null

    try {
      sponsor.value = await SponsorService.getById(sponsorId.value)
      // ดึงข้อมูลทุนของ Sponsor
      await fetchScholarships()
    } catch (err: any) {
      error.value = "Failed to load sponsor data."
      await Swal.fire({ 
        icon: 'error', 
        title: 'โหลดข้อมูลไม่สำเร็จ', 
        text: err?.response?.data?.message || 'ไม่สามารถโหลดข้อมูลบริษัทได้' 
      })
    } finally {
      loading.value = false
    }
  }

  // ดึงข้อมูลทุนของ Sponsor
  const fetchScholarships = async () => {
    if (sponsorId.value === null) return

    loadingScholarships.value = true
    try {
      scholarships.value = await SponsorService.getScholarships(sponsorId.value)
    } catch (err: any) {
      console.error('Failed to load scholarships:', err)
      scholarships.value = []
      // แสดง toast แบบไม่ block
      Swal.fire({ 
        icon: 'warning', 
        title: 'โหลดทุนไม่สำเร็จ', 
        text: 'ไม่สามารถโหลดข้อมูลทุนการศึกษาได้',
        timer: 3000,
        showConfirmButton: false
      })
    } finally {
      loadingScholarships.value = false
    }
  }

  // Helper functions
  const formatDate = (dateStr: string): string => {
    if (!dateStr) return '-'
    const date = new Date(dateStr)
    return date.toLocaleDateString('th-TH', { year: 'numeric', month: 'short', day: 'numeric' })
  }

  const getStatusColor = (status: string): string => {
    const statusLower = status?.toLowerCase() || ''
    if (statusLower.includes('open') || statusLower.includes('เปิด')) return 'bg-green-100 text-green-700'
    if (statusLower.includes('close') || statusLower.includes('ปิด')) return 'bg-red-100 text-red-700'
    return 'bg-gray-100 text-gray-700'
  }

  watch(sponsorId, (id) => {
    if (id === null) return
    fetchSponsor()
  }, { immediate: true })

  // Contact Modal State
  const isContactModalOpen = ref(false)

  const openContactModal = () => {
    isContactModalOpen.value = true
  }

  const onContactsSaved = (contacts: ContactResponse[]) => {
    if (sponsor.value) {
      sponsor.value.contacts = contacts
    }
  }

</script>

<template>
  <div class="space-y-10 min-h-screen">
    <div v-if="loading" class="text-center p-10">
      Loading Sponsor Data...
    </div>
    <div v-else-if="error" class="text-center p-10 text-red-600">
      Error: {{ error }}
    </div>

    <div v-else-if="sponsor">
      
      <!-- header -->
      <section class="bg-white mb-6 rounded-xl rounded-tl-[30px] shadow-sm ring-1 ring-gray-200 p-6">

        <!-- back button -->
        <button 
          @click="goBack"
          class="inline-flex items-center gap-2 text-gray-600 hover:text-gray-800 mb-4"
        >
          <ArrowLeft class="w-6 h-6" />
        </button>

        <div class="flex items-start gap-6">
          <div class="w-20 h-20 flex items-center justify-center shrink-0">
            <Building2 class="w-full h-full text-[#F37021]" />
          </div>

          <div class="flex-1 space-y-1">
            <h1 class="text-2xl font-semibold text-gray-900">
              {{ sponsor.company_name }}
            </h1>

            <p class="text-m text-gray-500">
              {{ sponsor.industry?.name || '—' }}
            </p>

            <div class="flex flex-wrap items-center gap-4 mt-2 text-m">
              <span v-if="sponsor.website" class="inline-flex items-center gap-1 text-blue-600 hover:underline">
                <Globe class="text-gray-400"/>
                <a :href="sponsor.website" target="_blank">
                  {{ sponsor.website }}
                </a>
              </span>
              <span v-else class="text-gray-400"><Globe class="text-gray-400"/> N/A</span>

              <span class="text-gray-600">
                Status:
                <span class="font-medium text-green-600 capitalize">
                  {{ sponsor.status }}
                </span>
              </span>
            </div>
          </div>

          <div class="text-right shrink-0">
            <div class="text-2xl font-semibold text-gray-900">{{ scholarships.length }}</div>
            <div class="text-xs text-gray-500">จำนวนทุน</div>
          </div>
        </div>
      </section>

      <!-- content -->
      <div class="grid grid-cols-1 lg:grid-cols-[2fr_1fr] gap-6 items-stretch">
        <!-- main: Scholarships List -->
        <section class="bg-white rounded-xl shadow-sm ring-1 ring-gray-200 p-6 flex flex-col h-full">
          <div class="flex items-center gap-2 mb-4">
            <GraduationCap class="w-5 h-5 text-[#F37021]" />
            <h2 class="text-lg font-semibold text-gray-900">
              ทุนการศึกษาทั้งหมด ({{ scholarships.length }})
            </h2>
          </div>

          <!-- Description -->
          <p v-if="sponsor.description" class="text-gray-600 text-sm mb-4 pb-4 border-b">
            {{ sponsor.description }}
          </p>

          <!-- Loading -->
          <div v-if="loadingScholarships" class="flex-1 flex items-center justify-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-[#F37021]"></div>
          </div>

          <!-- No Scholarships -->
          <div v-else-if="!scholarships.length" class="flex-1 flex flex-col items-center justify-center py-8 text-gray-400">
            <GraduationCap class="w-12 h-12 mb-2 opacity-50" />
            <p>ยังไม่มีทุนการศึกษา</p>
          </div>

          <!-- Scholarships List -->
          <div v-else class="flex-1 space-y-3 overflow-auto">
            <div
              v-for="scholarship in scholarships"
              :key="scholarship.ID"
              class="p-4 rounded-lg border border-gray-100 hover:border-[#F37021]/30 hover:bg-orange-50/30 transition-all cursor-pointer"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="flex-1 min-w-0">
                  <h3 class="font-medium text-gray-900 truncate">
                    {{ scholarship.scholarship_name }}
                  </h3>
                  <p class="text-sm text-gray-500 line-clamp-2 mt-1">
                    {{ scholarship.description }}
                  </p>
                  
                  <div class="flex flex-wrap items-center gap-3 mt-2 text-xs">
                    <!-- Type -->
                    <span class="inline-flex items-center gap-1 text-gray-500">
                      <Tag class="w-3 h-3" />
                      {{ scholarship.typescholarship?.type_name || '-' }}
                    </span>
                    <!-- Date -->
                    <span class="inline-flex items-center gap-1 text-gray-500">
                      <Calendar class="w-3 h-3" />
                      {{ formatDate(scholarship.open_date) }} - {{ formatDate(scholarship.close_date) }}
                    </span>
                  </div>
                </div>

                <!-- Status Badge -->
                <span 
                  class="shrink-0 px-2 py-1 rounded-full text-xs font-medium"
                  :class="getStatusColor(scholarship.statusscholarship?.status_name)"
                >
                  {{ scholarship.statusscholarship?.status_name || '-' }}
                </span>
              </div>
            </div>
          </div>
        </section>

        <!-- contact -->
        <section class="bg-white rounded-xl shadow-sm ring-1 ring-gray-200 p-6 flex flex-col h-full">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-gray-900">
              ผู้ติดต่อบริษัท ({{ sponsor.contacts?.length || 0 }})
            </h2>

            <button
              @click="openContactModal"
              class="p-2 text-gray-500 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
              title="จัดการผู้ติดต่อ"
            >
              <Pencil class="w-5 h-5" />
            </button>
          </div>

          <div class="flex-1 space-y-3">
            <!-- Empty State -->
            <div v-if="!sponsor.contacts?.length" class="flex flex-col items-center justify-center py-8 text-gray-400">
              <User class="w-10 h-10 mb-2 opacity-50" />
              <p class="text-sm">ยังไม่มีผู้ติดต่อ</p>
              <button 
                @click="openContactModal"
                class="mt-2 text-blue-600 hover:underline text-sm"
              >
                + เพิ่มผู้ติดต่อ
              </button>
            </div>

            <!-- Contact List with Avatar -->
            <div
              v-for="c in sponsor.contacts"
              :key="c.ID"
              class="flex items-center gap-4 p-3 rounded-lg hover:bg-gray-50 transition-colors"
            >
              <!-- Avatar -->
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white font-semibold text-lg shrink-0">
                {{ c.name?.charAt(0)?.toUpperCase() || '?' }}
              </div>

              <!-- Info -->
              <div class="flex-1 min-w-0">
                <div class="font-medium text-gray-900 truncate">
                  {{ c.name }}
                </div>
                <div class="text-sm text-gray-500">
                  {{ c.position || '-' }}
                </div>
                <div class="flex flex-wrap items-center gap-3 mt-1 text-xs text-gray-500">
                  <span class="inline-flex items-center gap-1">
                    <Mail class="w-3 h-3" />
                    {{ c.email }}
                  </span>
                  <span class="inline-flex items-center gap-1">
                    <Phone class="w-3 h-3" />
                    {{ c.phone }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </section>

        <!-- Contact Modal -->
        <SponsorContact
          v-if="sponsor"
          v-model:isOpen="isContactModalOpen"
          :sponsorId="sponsor.ID"
          :sponsorName="sponsor.company_name"
          :initialContacts="sponsor.contacts || []"
          @saved="onContactsSaved"
        />
      </div>
    </div>
  </div>
</template>