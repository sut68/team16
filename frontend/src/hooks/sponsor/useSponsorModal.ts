import { ref } from 'vue'
import { SponsorService } from '@/services/sponsor/sponsor'
import type { SponsorResponse } from '@/interfaces/sponsor'
import Swal from 'sweetalert2'

export function useSponsorModals() {
  // Modal States
  const isCreateSponsorOpen = ref(false)
  const isEditSponsorOpen = ref(false)
  const isContactsOpen = ref(false)

  // Loading States
  const creating = ref(false)
  const updating = ref(false)

  // Current Data
  const currentSponsor = ref<SponsorResponse | null>(null)
  const currentSponsorId = ref<number | null>(null)
  const currentContacts = ref<any[]>([])
  const serverErrors = ref<Record<string, any> | null>(null)

  // Create Modal
  function openCreateForm() {
    isCreateSponsorOpen.value = true
  }

  function closeCreateForm() {
    isCreateSponsorOpen.value = false
  }

  // Edit Modal
  async function openEditForm(id: number) {
    try {
      const sponsor = await SponsorService.getById(id)
      currentSponsor.value = sponsor
      serverErrors.value = null
      isEditSponsorOpen.value = true
    } catch (err) {
      console.error(err)
      Swal.fire({ icon: 'error', title: 'ไม่พบข้อมูลบริษัท' })
    }
  }

  function closeEditForm() {
    isEditSponsorOpen.value = false
    currentSponsor.value = null
  }

  // Contacts Modal
  async function openContacts(id: number) {
    currentSponsorId.value = id

    // Reuse cached sponsor if available
    if (currentSponsor.value?.ID === id && Array.isArray(currentSponsor.value.contacts)) {
      currentContacts.value = currentSponsor.value.contacts
      isContactsOpen.value = true
      return
    }

    try {
      const sponsor = await SponsorService.getById(id)
      currentContacts.value = sponsor?.contacts ?? []
      currentSponsor.value = sponsor
      isContactsOpen.value = true
    } catch (err: any) {
      console.error('openContacts: failed to fetch sponsor', err)
      Swal.fire({ icon: 'error', title: 'ไม่สามารถโหลดผู้ติดต่อได้' })
    }
  }

  function closeContacts() {
    isContactsOpen.value = false
  }

  return {
    // Modal States
    isCreateSponsorOpen,
    isEditSponsorOpen,
    isContactsOpen,
    // Loading
    creating,
    updating,
    // Data
    currentSponsor,
    currentSponsorId,
    currentContacts,
    serverErrors,
    // Actions
    openCreateForm,
    closeCreateForm,
    openEditForm,
    closeEditForm,
    openContacts,
    closeContacts,
  }
}