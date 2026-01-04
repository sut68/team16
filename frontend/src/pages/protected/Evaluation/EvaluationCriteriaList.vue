<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import type { EvaluationCriterionResponse, EvaluationCriterionPayload } from '@/interfaces/evaluation'
  import { EvaluationCriteriaService } from '@/services/evaluation/evaluation'
  import Swal from 'sweetalert2'
  import StatCard from '@/components/ui/StatCard.vue'

  // Icons
  import { Plus, Pencil, Trash2, Search, X, Filter, CheckCircle, XCircle } from 'lucide-vue-next'

  // ========== State ==========
  const criteria = ref<EvaluationCriterionResponse[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const searchQuery = ref('')

  // Modal State
  const isFormOpen = ref(false)
  const isEditing = ref(false)
  const formLoading = ref(false)
  const currentCriterion = ref<EvaluationCriterionResponse | null>(null)

  // Form Data
  const formData = ref<EvaluationCriterionPayload>({
    name: '',
    description: '',
    score_type: 'numeric',
    max_score: 100,
    weight: 1.0,
    is_active: true,
  })

  // Validation Errors
  const formErrors = ref<{
    name?: string
    description?: string
    max_score?: string
    weight?: string
  }>({})

  // ========== Validation Rules (ตาม Backend) ==========
  const VALIDATION_RULES = {
    name: { minLength: 2, maxLength: 100 },
    description: { maxLength: 500 },
    max_score: { min: 0, max: 1000 },
    weight: { min: 0, max: 10 },
  }

  // ========== Computed ==========
  const filteredCriteria = computed(() => {
    const term = searchQuery.value.toLowerCase().trim()
    if (!term) return criteria.value
    return criteria.value.filter(c => 
      c.name.toLowerCase().includes(term) ||
      c.description?.toLowerCase().includes(term)
    )
  })

  const totalCriteria = computed(() => criteria.value.length)
  const activeCriteria = computed(() => criteria.value.filter(c => c.is_active).length)
  const inactiveCriteria = computed(() => criteria.value.filter(c => !c.is_active).length)

  // ========== Methods ==========
  async function fetchCriteria() {
    loading.value = true
    error.value = null
    try {
      criteria.value = await EvaluationCriteriaService.getAll()
    } catch (err: any) {
      console.error('Failed to fetch criteria:', err)
      error.value = err?.message || 'ไม่สามารถโหลดข้อมูลได้'
    } finally {
      loading.value = false
    }
  }

  function openCreateForm() {
    isEditing.value = false
    currentCriterion.value = null
    formData.value = {
      name: '',
      description: '',
      score_type: 'numeric',
      max_score: 100,
      weight: 1.0,
      is_active: true,
    }
    formErrors.value = {}
    isFormOpen.value = true
  }

  function openEditForm(criterion: EvaluationCriterionResponse) {
    isEditing.value = true
    currentCriterion.value = criterion
    formData.value = {
      name: criterion.name,
      description: criterion.description,
      score_type: criterion.score_type,
      max_score: criterion.max_score,
      weight: criterion.weight,
      is_active: criterion.is_active,
    }
    formErrors.value = {}
    isFormOpen.value = true
  }

  function closeForm() {
    isFormOpen.value = false
    currentCriterion.value = null
    formErrors.value = {}
  }

  // ========== Form Validation ==========
  function validateForm(): boolean {
    formErrors.value = {}
    let isValid = true

    // Name validation: required, min 2, max 100
    const name = formData.value.name?.trim() || ''
    if (!name) {
      formErrors.value.name = 'กรุณากรอกชื่อเกณฑ์'
      isValid = false
    } else if (name.length < VALIDATION_RULES.name.minLength) {
      formErrors.value.name = `ชื่อเกณฑ์ต้องมีอย่างน้อย ${VALIDATION_RULES.name.minLength} ตัวอักษร`
      isValid = false
    } else if (name.length > VALIDATION_RULES.name.maxLength) {
      formErrors.value.name = `ชื่อเกณฑ์ต้องไม่เกิน ${VALIDATION_RULES.name.maxLength} ตัวอักษร`
      isValid = false
    }

    // Description validation: max 500
    const description = formData.value.description || ''
    if (description.length > VALIDATION_RULES.description.maxLength) {
      formErrors.value.description = `รายละเอียดต้องไม่เกิน ${VALIDATION_RULES.description.maxLength} ตัวอักษร`
      isValid = false
    }

    // MaxScore validation: 0-1000
    const maxScore = formData.value.max_score ?? 0
    if (maxScore < VALIDATION_RULES.max_score.min || maxScore > VALIDATION_RULES.max_score.max) {
      formErrors.value.max_score = `คะแนนเต็มต้องอยู่ระหว่าง ${VALIDATION_RULES.max_score.min} - ${VALIDATION_RULES.max_score.max}`
      isValid = false
    }

    // Weight validation: 0-10
    const weight = formData.value.weight ?? 0
    if (weight < VALIDATION_RULES.weight.min || weight > VALIDATION_RULES.weight.max) {
      formErrors.value.weight = `น้ำหนักต้องอยู่ระหว่าง ${VALIDATION_RULES.weight.min} - ${VALIDATION_RULES.weight.max}`
      isValid = false
    }

    return isValid
  }

  // ========== Real-time Validation Functions ==========
  function validateNameRealtime() {
    const name = formData.value.name?.trim() || ''
    if (!name) {
      formErrors.value.name = 'กรุณากรอกชื่อเกณฑ์'
    } else if (name.length < VALIDATION_RULES.name.minLength) {
      formErrors.value.name = `ชื่อเกณฑ์ต้องมีอย่างน้อย ${VALIDATION_RULES.name.minLength} ตัวอักษร`
    } else if (name.length > VALIDATION_RULES.name.maxLength) {
      formErrors.value.name = `ชื่อเกณฑ์ต้องไม่เกิน ${VALIDATION_RULES.name.maxLength} ตัวอักษร`
    } else {
      delete formErrors.value.name
    }
  }

  function validateDescriptionRealtime() {
    const description = formData.value.description || ''
    if (description.length > VALIDATION_RULES.description.maxLength) {
      formErrors.value.description = `รายละเอียดต้องไม่เกิน ${VALIDATION_RULES.description.maxLength} ตัวอักษร`
    } else {
      delete formErrors.value.description
    }
  }

  function validateMaxScoreRealtime() {
    const maxScore = formData.value.max_score
    if (maxScore === null || maxScore === undefined || isNaN(maxScore)) {
      formErrors.value.max_score = 'กรุณากรอกคะแนนเต็ม'
    } else if (maxScore < VALIDATION_RULES.max_score.min || maxScore > VALIDATION_RULES.max_score.max) {
      formErrors.value.max_score = `คะแนนเต็มต้องอยู่ระหว่าง ${VALIDATION_RULES.max_score.min} - ${VALIDATION_RULES.max_score.max}`
    } else {
      delete formErrors.value.max_score
    }
  }

  function validateWeightRealtime() {
    const weight = formData.value.weight
    if (weight === null || weight === undefined || isNaN(weight)) {
      formErrors.value.weight = 'กรุณากรอกน้ำหนัก'
    } else if (weight < VALIDATION_RULES.weight.min || weight > VALIDATION_RULES.weight.max) {
      formErrors.value.weight = `น้ำหนักต้องอยู่ระหว่าง ${VALIDATION_RULES.weight.min} - ${VALIDATION_RULES.weight.max}`
    } else {
      delete formErrors.value.weight
    }
  }

  const isFormValid = computed(() => {
    const name = formData.value.name?.trim() || ''
    const description = formData.value.description || ''
    const maxScore = formData.value.max_score ?? 0
    const weight = formData.value.weight ?? 0

    return (
      name.length >= VALIDATION_RULES.name.minLength &&
      name.length <= VALIDATION_RULES.name.maxLength &&
      description.length <= VALIDATION_RULES.description.maxLength &&
      maxScore >= VALIDATION_RULES.max_score.min &&
      maxScore <= VALIDATION_RULES.max_score.max &&
      weight >= VALIDATION_RULES.weight.min &&
      weight <= VALIDATION_RULES.weight.max
    )
  })

  async function handleSubmit() {
    // Validate before submit
    if (!validateForm()) {
      return
    }

    formLoading.value = true
    try {
      if (isEditing.value && currentCriterion.value) {
        await EvaluationCriteriaService.update(currentCriterion.value.ID, formData.value)
        await Swal.fire({ icon: 'success', title: 'อัปเดตเกณฑ์สำเร็จ', timer: 1500, showConfirmButton: false })
      } else {
        await EvaluationCriteriaService.create(formData.value)
        await Swal.fire({ icon: 'success', title: 'สร้างเกณฑ์สำเร็จ', timer: 1500, showConfirmButton: false })
      }
      closeForm()
      await fetchCriteria()
    } catch (err: any) {
      console.error('Submit error:', err)
      const resp = err?.response?.data
      Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: resp?.error || resp?.message || 'กรุณาลองใหม่อีกครั้ง',
      })
    } finally {
      formLoading.value = false
    }
  }

  async function handleDelete(criterion: EvaluationCriterionResponse) {
    const result = await Swal.fire({
      icon: 'warning',
      title: 'ยืนยันการลบ',
      text: `ต้องการลบเกณฑ์ "${criterion.name}" ใช่หรือไม่?`,
      showCancelButton: true,
      confirmButtonColor: '#ef4444',
      confirmButtonText: 'ลบ',
      cancelButtonText: 'ยกเลิก',
    })

    if (result.isConfirmed) {
      try {
        await EvaluationCriteriaService.delete(criterion.ID)
        await Swal.fire({ icon: 'success', title: 'ลบสำเร็จ', timer: 1500, showConfirmButton: false })
        await fetchCriteria()
      } catch (err: any) {
        console.error('Delete error:', err)
        Swal.fire({
          icon: 'error',
          title: 'ไม่สามารถลบได้',
          text: err?.response?.data?.error || 'เกณฑ์นี้อาจมีการใช้งานอยู่',
        })
      }
    }
  }

  async function toggleActive(criterion: EvaluationCriterionResponse) {
    try {
      await EvaluationCriteriaService.update(criterion.ID, { is_active: !criterion.is_active })
      await fetchCriteria()
    } catch (err: any) {
      console.error('Toggle error:', err)
      Swal.fire({ icon: 'error', title: 'เกิดข้อผิดพลาด' })
    }
  }

  function getScoreTypeLabel(type: string) {
    switch (type) {
      case 'numeric': return 'คะแนนตัวเลข'
      case 'grade': return 'เกรด (A-F)'
      case 'pass_fail': return 'ผ่าน/ไม่ผ่าน'
      default: return type
    }
  }

  function getScoreTypeBadgeClass(type: string) {
    switch (type) {
      case 'numeric': return 'bg-blue-100 text-blue-800'
      case 'grade': return 'bg-purple-100 text-purple-800'
      case 'pass_fail': return 'bg-amber-100 text-amber-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  // ========== Lifecycle ==========
  onMounted(fetchCriteria)
</script>

<template>
  <div 
    class="evaluation-criteria-wrapper w-full mx-auto flex flex-col h-full p-6 bg-white" 
    data-theme="light"
  >
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold text-[#1e3a8a]">เกณฑ์การประเมิน</h1>

      <div class="flex items-center gap-3 w-full md:w-auto">
        <!-- Search Bar -->
        <div class="relative flex-1 md:w-64">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="ค้นหาเกณฑ์..."
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

        <!-- Add Button -->
        <button 
          @click="openCreateForm" 
          class="btn-primary"
        >
          <Plus class="w-4 h-4" />
          <span class="hidden sm:inline">เพิ่มเกณฑ์</span>
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
      <StatCard
        label="เกณฑ์ทั้งหมด"
        :value="totalCriteria"
        icon="total"
        color="blue"
      />
      <StatCard
        label="ใช้งานอยู่"
        :value="activeCriteria"
        icon="active"
        color="green"
      />
      <StatCard
        label="ปิดใช้งาน"
        :value="inactiveCriteria"
        icon="default"
        color="slate"
      />
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="p-6 w-full">
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
    <div v-if="error" class="p-4 bg-red-100 text-red-700 rounded-lg mb-4">{{ error }}</div>

    <!-- Content -->
    <div v-if="!loading && !error" class="flex-1 min-h-0 flex flex-col gap-4">
      <!-- Table -->
      <div class="overflow-x-auto overflow-y-auto flex-1 min-h-[400px] rounded-xl border border-slate-200 bg-white shadow-sm">
        <table class="table table-sm w-full">
          <thead>
            <tr class="bg-slate-800 text-white">
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">#</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">ชื่อเกณฑ์</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 hidden md:table-cell">รายละเอียด</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200">ประเภท</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-center">คะแนนเต็ม</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-center">น้ำหนัก</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-center">สถานะ</th>
              <th class="py-4 px-4 text-xs font-semibold uppercase tracking-wider text-slate-200 text-right">การจัดการ</th>
            </tr>
          </thead>

          <tbody class="divide-y divide-slate-100">
            <tr 
              v-for="(c, index) in filteredCriteria" 
              :key="c.ID" 
              class="transition-all duration-200 hover:bg-blue-50/50"
              :class="index % 2 === 0 ? 'bg-white' : 'bg-slate-50/50'"
            >
              <td class="py-3 px-4 text-sm text-slate-500 font-mono">{{ c.ID }}</td>
              
              <td class="py-3 px-4">
                <span class="font-semibold text-slate-800">{{ c.name }}</span>
              </td>

              <td class="py-3 px-4 text-sm text-slate-600 hidden md:table-cell">
                <span class="truncate block max-w-xs" :title="c.description">
                  {{ c.description || '-' }}
                </span>
              </td>

              <td class="py-3 px-4">
                <span 
                  class="inline-flex items-center px-2.5 py-1 text-xs font-semibold rounded-md"
                  :class="getScoreTypeBadgeClass(c.score_type)"
                >
                  {{ getScoreTypeLabel(c.score_type) }}
                </span>
              </td>

              <td class="py-3 px-4 text-center">
                <span class="font-mono text-sm">{{ c.max_score }}</span>
              </td>

              <td class="py-3 px-4 text-center">
                <span class="font-mono text-sm">{{ c.weight }}</span>
              </td>

              <td class="py-3 px-4 text-center">
                <button 
                  @click="toggleActive(c)"
                  class="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-semibold cursor-pointer transition-colors"
                  :class="c.is_active 
                    ? 'bg-emerald-100 text-emerald-700 hover:bg-emerald-200' 
                    : 'bg-slate-100 text-slate-500 hover:bg-slate-200'"
                >
                  <CheckCircle v-if="c.is_active" class="w-3.5 h-3.5" />
                  <XCircle v-else class="w-3.5 h-3.5" />
                  {{ c.is_active ? 'Active' : 'Inactive' }}
                </button>
              </td>

              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button 
                    @click="openEditForm(c)"
                    class="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                    title="แก้ไข"
                  >
                    <Pencil class="w-4 h-4" />
                  </button>
                  <button 
                    @click="handleDelete(c)"
                    class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    title="ลบ"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>

            <!-- Empty State -->
            <tr v-if="filteredCriteria.length === 0">
              <td colspan="8" class="h-[300px] text-center align-middle">
                <div class="flex flex-col items-center justify-center text-gray-400">
                  <Filter class="w-12 h-12 mb-3 opacity-50" />
                  <p v-if="searchQuery" class="text-gray-500">
                    ไม่พบเกณฑ์ที่ตรงกับ "<strong class="text-gray-700">{{ searchQuery }}</strong>"
                  </p>
                  <p v-else class="text-gray-500">ยังไม่มีเกณฑ์ในระบบ</p>
                  <button 
                    v-if="!searchQuery"
                    @click="openCreateForm" 
                    class="btn-primary mt-4"
                  >
                    <Plus class="w-4 h-4" />
                    <span>เพิ่มเกณฑ์แรก</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form -->
    <Teleport to="body">
      <div 
        v-if="isFormOpen" 
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        @click.self="closeForm"
      >
        <div class="bg-white rounded-2xl shadow-2xl w-full max-w-lg mx-4 overflow-hidden">
          <!-- Header -->
          <div class="bg-gradient-to-r from-indigo-600 to-blue-600 px-6 py-4">
            <h2 class="text-xl font-bold text-white">
              {{ isEditing ? 'แก้ไขเกณฑ์' : 'เพิ่มเกณฑ์ใหม่' }}
            </h2>
          </div>

          <!-- Form -->
          <form @submit.prevent="handleSubmit" class="p-6 space-y-4">
            <!-- Name -->
            <div>
              <div class="flex items-center justify-between mb-1">
                <label class="block text-sm font-medium text-gray-700">ชื่อเกณฑ์ *</label>
                <span class="text-xs" :class="(formData.name?.length || 0) > 100 ? 'text-red-500' : 'text-gray-400'">
                  {{ formData.name?.length || 0 }}/100
                </span>
              </div>
              <input
                v-model="formData.name"
                @input="validateNameRealtime"
                type="text"
                required
                minlength="2"
                maxlength="100"
                class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="formErrors.name ? 'border-red-500' : 'border-gray-300'"
                placeholder="เช่น ผลการเรียน, ทักษะการสื่อสาร"
              />
              <p v-if="formErrors.name" class="mt-1 text-xs text-red-500">{{ formErrors.name }}</p>
            </div>

            <!-- Description -->
            <div>
              <div class="flex items-center justify-between mb-1">
                <label class="block text-sm font-medium text-gray-700">รายละเอียด</label>
                <span class="text-xs" :class="(formData.description?.length || 0) > 500 ? 'text-red-500' : 'text-gray-400'">
                  {{ formData.description?.length || 0 }}/500
                </span>
              </div>
              <textarea
                v-model="formData.description"
                @input="validateDescriptionRealtime"
                rows="2"
                maxlength="500"
                class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :class="formErrors.description ? 'border-red-500' : 'border-gray-300'"
                placeholder="อธิบายเกณฑ์การประเมิน..."
              ></textarea>
              <p v-if="formErrors.description" class="mt-1 text-xs text-red-500">{{ formErrors.description }}</p>
            </div>

            <!-- Score Type -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">ประเภทคะแนน</label>
              <select
                v-model="formData.score_type"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              >
                <option value="numeric">คะแนนตัวเลข</option>
                <option value="grade">เกรด (A-F)</option>
                <option value="pass_fail">ผ่าน/ไม่ผ่าน</option>
              </select>
            </div>

            <!-- Max Score & Weight -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">คะแนนเต็ม (0-1000)</label>
                <input
                  v-model.number="formData.max_score"
                  @input="validateMaxScoreRealtime"
                  type="number"
                  min="0"
                  max="100"
                  step="0.01"
                  class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  :class="formErrors.max_score ? 'border-red-500' : 'border-gray-300'"
                />
                <p v-if="formErrors.max_score" class="mt-1 text-xs text-red-500">{{ formErrors.max_score }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">น้ำหนัก (0-10)</label>
                <input
                  v-model.number="formData.weight"
                  @input="validateWeightRealtime"
                  type="number"
                  min="0"
                  max="10"
                  step="0.1"
                  class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  :class="formErrors.weight ? 'border-red-500' : 'border-gray-300'"
                />
                <p v-if="formErrors.weight" class="mt-1 text-xs text-red-500">{{ formErrors.weight }}</p>
              </div>
            </div>

            <!-- Is Active -->
            <div class="flex items-center gap-3">
              <input
                v-model="formData.is_active"
                type="checkbox"
                id="is_active"
                class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
              />
              <label for="is_active" class="text-sm text-gray-700">เปิดใช้งาน</label>
            </div>

            <!-- Actions -->
            <div class="flex justify-end gap-3 pt-4 border-t">
              <button 
                type="button"
                @click="closeForm"
                class="px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
              >
                ยกเลิก
              </button>
              <button 
                type="submit"
                :disabled="formLoading || !isFormValid"
                class="px-6 py-2 text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors disabled:opacity-50"
              >
                <span v-if="formLoading">กำลังบันทึก...</span>
                <span v-else>{{ isEditing ? 'อัปเดต' : 'สร้าง' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.btn-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.5rem 1.25rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: white;
  background: linear-gradient(135deg, #3b82f6, #1e40af);
  border-radius: 9999px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 150ms ease;
  cursor: pointer;
  border: none;
}

.btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
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
</style>
