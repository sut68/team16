// useIndustries.ts
// Hook สำหรับโหลดและจัดการข้อมูล Industries
// ใช้ร่วมกันใน SponsorForm.vue และ SponsorEdit.vue

import { ref, onMounted } from 'vue'
import type { IndustryResponse } from '@/interfaces/sponsor'
import { IndustryService } from '@/services/sponsor/industry'

export interface UseIndustriesOptions {
  // โหลดอัตโนมัติเมื่อ mount
  autoLoad?: boolean
  // industries ที่ส่งมาจาก parent (ถ้ามีจะไม่โหลดใหม่)
  initialIndustries?: IndustryResponse[] | null
}

export function useIndustries(options: UseIndustriesOptions = {}) {
  const { autoLoad = true, initialIndustries = null } = options

  // State
  const industries = ref<IndustryResponse[]>(initialIndustries ?? [])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const loaded = ref(initialIndustries !== null && initialIndustries.length > 0)

  // โหลดข้อมูล Industries
  async function load(force = false): Promise<void> {
    // ถ้ามี initialIndustries แล้ว และไม่ force reload ก็ไม่ต้องโหลด
    if (!force && loaded.value && industries.value.length > 0) {
      return
    }

    // ถ้ากำลังโหลดอยู่ก็ไม่ต้องโหลดซ้ำ
    if (loading.value) {
      return
    }

    loading.value = true
    error.value = null

    try {
      const res = await IndustryService.getAll()
      industries.value = res ?? []
      loaded.value = true
    } catch (err: any) {
      console.error('โหลดอุตสาหกรรมผิดพลาด:', err)
      error.value = err?.message ?? 'โหลดอุตสาหกรรมไม่สำเร็จ'
      industries.value = []
    } finally {
      loading.value = false
    }
  }

  // ตั้งค่า industries จาก parent
  function setIndustries(data: IndustryResponse[]) {
    industries.value = data
    loaded.value = true
  }

  // Reset state
  function reset() {
    industries.value = []
    loading.value = false
    error.value = null
    loaded.value = false
  }

  // Auto load on mount if enabled
  if (autoLoad) {
    onMounted(() => {
      if (!loaded.value) {
        load()
      }
    })
  }

  return {
    // State
    industries,
    loading,
    error,
    loaded,
    // Actions
    load,
    setIndustries,
    reset,
  }
}
