import { ref, computed } from 'vue'
import { SponsorService } from '@/services/sponsor/sponsor'
import type { SponsorResponse, SponsorView } from '@/interfaces/sponsor'
import Swal from 'sweetalert2'

/**
 * แปลง API response เป็น view model
 */
function toSponsorView(s: SponsorResponse): SponsorView {
  return {
    ID: s.ID,
    company_name: s.company_name,
    website: s.website ?? null,
    industry_name: s.industry?.name ?? null,
    industry_id: s.industry_id ?? s.industry?.ID ?? null,
    status: s.status,
    contacts_count: Array.isArray(s.contacts) ? s.contacts.length : 0,
  }
}

export function useSponsorList() {
  // State
  const sponsors = ref<SponsorView[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const loadingToggle = ref<number | null>(null)

  // Computed Stats
  const totalSponsors = computed(() => sponsors.value.length)

  const activeSponsors = computed(() =>
    sponsors.value.filter(s => s.status === 'active').length
  )

  const totalContacts = computed(() =>
    sponsors.value.reduce((sum, s) => sum + (s.contacts_count ?? 0), 0)
  )

  // Actions
  async function fetchSponsors() {
    loading.value = true
    error.value = null
    try {
      const res = await SponsorService.getAll()
      sponsors.value = Array.isArray(res)
        ? res.map((s: SponsorResponse) => toSponsorView(s))
        : []
    } catch (err: any) {
      console.error(err)
      error.value = err?.response?.data?.message || err?.message || 'โหลดข้อมูล Sponsors ไม่สำเร็จ'
    } finally {
      loading.value = false
    }
  }

  async function toggleStatus(id: number, currentStatus?: string) {
    const idx = sponsors.value.findIndex(s => s.ID === id)
    if (idx === -1) return

    const prev = sponsors.value[idx]
    if (!prev) return

    const prevStatus = prev.status
    const newStatus = currentStatus === 'active' ? 'inactive' : 'active'

    // Optimistic update
    sponsors.value[idx] = { ...prev, status: newStatus }
    loadingToggle.value = id

    try {
      const payload = { status: newStatus }
      const updated = await SponsorService.update(id, payload)
      if (updated?.status) {
        sponsors.value[idx] = { ...sponsors.value[idx], status: updated.status }
      }
    } catch (err) {
      // Rollback
      sponsors.value[idx] = { ...sponsors.value[idx], status: prevStatus }
      console.error('Toggle status failed', err)
      Swal.fire({ icon: 'error', title: 'ผิดพลาด', text: 'เปลี่ยนสถานะไม่สำเร็จ ลองใหม่อีกครั้ง' })
    } finally {
      loadingToggle.value = null
    }
  }

  async function removeOne(id: number) {
    loading.value = true
    error.value = null
    const snapshot = [...sponsors.value]
    sponsors.value = sponsors.value.filter(s => s.ID !== id)

    try {
      await SponsorService.delete(id)
      Swal.fire({ icon: 'success', title: 'ลบสำเร็จ' })
      await fetchSponsors()
    } catch (err: any) {
      console.error(err)
      sponsors.value = snapshot
      Swal.fire({ icon: 'error', title: 'ลบไม่สำเร็จ', text: err?.response?.data?.message || 'เกิดข้อผิดพลาด' })
    } finally {
      loading.value = false
    }
  }

  /**
   * เพิ่ม sponsor ใหม่เข้าไปใน list (optimistic)
   */
  function addSponsor(created: SponsorResponse) {
    const view = toSponsorView(created)
    if (!sponsors.value.some(s => s.ID === view.ID)) {
      sponsors.value.unshift(view)
    } else {
      const idx = sponsors.value.findIndex(s => s.ID === view.ID)
      if (idx !== -1) sponsors.value.splice(idx, 1, view)
    }
  }

  /**
   * อัปเดต contacts_count สำหรับ sponsor
   */
  function updateContactsCount(id: number, count: number) {
    const idx = sponsors.value.findIndex(s => s.ID === id)
    if (idx !== -1 && sponsors.value[idx]) {
      sponsors.value[idx].contacts_count = count
    }
  }

  return {
    // State
    sponsors,
    loading,
    error,
    loadingToggle,
    // Computed
    totalSponsors,
    activeSponsors,
    totalContacts,
    // Actions
    fetchSponsors,
    toggleStatus,
    removeOne,
    addSponsor,
    updateContactsCount,
    toSponsorView,
  }
}