<script setup lang="ts">
  import SponsorService, { type BatchContactsPayload } from '@/services/sponsor/sponsor';
  import type { ContactPayload, ContactResponse, SponsorResponse } from '../../../interfaces/sponsor';
  import { computed, ref, watch } from 'vue';
  import { useRoute } from 'vue-router';
  import { Building2, Globe } from 'lucide-vue-next';

  const route = useRoute();
  const sponsorId = computed<number | null>(() => {
    const param = route.params.id
    if (!param) return null

    const value = Array.isArray(param) ? param[0] : param
    const id = Number(value)

    return Number.isNaN(id) ? null : id
  })

  const sponsor = ref<SponsorResponse | null>(null)
  const loading = ref(true)
  const error = ref<string | null>(null)

  // Data Fetching ขอข้อมูลจาก Backend
  const fetchSponsor = async () => {
    if (sponsorId.value === null) {
      error.value = "Invalid Sponsor ID in URL."
      return
    }

    loading.value = true
    error.value = null

    try {
      sponsor.value = await SponsorService.getById(sponsorId.value)
    } catch (err) {
      error.value = "Failed to load sponsor data."
    } finally {
      loading.value = false
    }
  }

  watch(sponsorId, (id) => {
    if (id === null) return
    fetchSponsor()
  }, { immediate: true })

  const savingContacts = ref(false)
  const updateSponsorContacts = async (payload: BatchContactsPayload) => {
    if (!sponsor.value) return;
    savingContacts.value = true

    try {
      const { contacts } = await SponsorService.updateContacts(
        sponsor.value.ID,
        payload
      )
      sponsor.value.contacts = contacts
      window.alert("บันทึกผู้ติดต่อสำเร็จ")

    } catch {
      window.alert("บันทึกผู้ติดต่อไม่สำเร็จ")

    } finally {
      savingContacts.value = false

    }
  }

  // handlers
  const onAddContact = () => {
    const name = window.prompt('ชื่อผู้ติดต่อใหม่:');
    if (!name) return;
    const newContactPayload: ContactPayload = {
      name,
      email: window.prompt('อีเมล:', 'new@example.com') || '',
      phone: window.prompt('เบอร์โทรศัพท์:', 'N/A') || '',
      position: window.prompt('ตำแหน่ง:', 'General Contact'),
    };

    updateSponsorContacts({
      upsert: [newContactPayload]
    });
  }

  const onDeleteContact = (contact: ContactResponse) => {
    const ok = window.confirm(`ลบผู้ติดต่อ ${contact.name} จริงหรือไม่?`);
    if (!ok) return;

    updateSponsorContacts({
      delete_ids: [contact.ID]
    });
  };

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
            <div class="text-2xl font-semibold text-gray-900">5</div>
            <div class="text-xs text-gray-500">จำนวนทุน</div>
          </div>
        </div>
      </section>

      <!-- content -->
      <div class="grid grid-cols-1 lg:grid-cols-[2fr_1fr] gap-6 items-stretch">
        <!-- main -->
        <section class="bg-white rounded-xl shadow-sm ring-1 ring-gray-200 p-6 flex flex-col h-full">
          <h2 class="text-lg font-semibold text-gray-900 mb-3">
            ข้อมูลทุนทั้งหมด
          </h2>

          <p class="text-gray-700 leading-relaxed">
            {{ sponsor.description || 'ไม่มีคำอธิบาย...' }}
          </p>

          <div class="mt-4 text-sm text-gray-400">
            (ส่วนนี้จะแสดงรายการทุนที่เกี่ยวข้องกับผู้สนับสนุนรายนี้)
          </div>
        </section>

        <!-- contact -->
        <section class="bg-white rounded-xl shadow-sm ring-1 ring-gray-200 p-6 flex flex-col h-full">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-gray-900">
              ผู้ติดต่อบริษัท
            </h2>

            <button
              @click="onAddContact"
              class="inline-flex items-center gap-1 rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700"
            >
              + เพิ่มผู้ติดต่อ
            </button>
          </div>

          <div class="flex-1">

            <div v-if="!sponsor.contacts?.length" class="text-gray-500 text-sm">
              ยังไม่มีผู้ติดต่อ
            </div>
  
            <div
              v-for="c in sponsor.contacts"
              :key="c.ID"
              class="py-3 border-t first:border-t-0"
            >
              <div class="flex items-start justify-between">
                <div class="space-y-0.5">
                  <div class="font-medium text-gray-900">
                    {{ c.name }}
                  </div>
                  <div class="text-sm text-gray-500">
                    {{ c.position || 'N/A' }}
                  </div>
                  <div class="text-sm text-gray-500">
                    📧 {{ c.email }}
                  </div>
                  <div class="text-sm text-gray-500">
                    📞 {{ c.phone }}
                  </div>
                </div>
  
                <div class="flex gap-3 text-xs">
                  <button class="text-blue-600 hover:underline">
                    แก้ไข
                  </button>
                  <button
                    @click="onDeleteContact(c)"
                    class="text-red-600 hover:underline"
                  >
                    ลบ
                  </button>
                </div>
              </div>
            </div>
          </div>
        </section>

      </div>
    </div>
  </div>
</template>