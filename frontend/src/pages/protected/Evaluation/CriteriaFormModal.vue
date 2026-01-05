<script setup lang="ts">
  import { ref, computed, watch } from 'vue'
  import type { EvaluationCriterionResponse, EvaluationCriterionPayload } from '@/interfaces/evaluation'
  import { EvaluationCriteriaService } from '@/services/evaluation/evaluation'
  import Swal from 'sweetalert2'

  // ========== Validators ==========
  import {
    validateCriterionForm,
    validateCriterionName,
    validateCriterionDescription,
    validateCriterionMaxScore,
    validateCriterionWeight,
    isCriterionFormValid,
    type CriterionFormErrors,
  } from '@/validators/evaluation_validator'

  // ========== Props & Emits ==========
  const props = defineProps<{
    isOpen: boolean
    criterion?: EvaluationCriterionResponse | null
  }>()

  const emit = defineEmits<{
    (e: 'close'): void
    (e: 'saved'): void
  }>()

  // ========== State ==========
  const formLoading = ref(false)
  const formErrors = ref<CriterionFormErrors>({})

  // Form Data
  const formData = ref<EvaluationCriterionPayload>({
    name: '',
    description: '',
    score_type: 'numeric',
    max_score: 100,
    weight: 1.0,
    is_active: true,
  })

  // ========== Computed ==========
  const isEditing = computed(() => !!props.criterion)
  const isFormValid = computed(() => isCriterionFormValid(formData.value))

  // ========== Watch ==========
  watch(() => props.isOpen, (open) => {
    if (open) {
      if (props.criterion) {
        // Edit mode
        formData.value = {
          name: props.criterion.name,
          description: props.criterion.description,
          score_type: props.criterion.score_type,
          max_score: props.criterion.max_score,
          weight: props.criterion.weight,
          is_active: props.criterion.is_active,
        }
      } else {
        // Create mode
        formData.value = {
          name: '',
          description: '',
          score_type: 'numeric',
          max_score: 100,
          weight: 1.0,
          is_active: true,
        }
      }
      formErrors.value = {}
    }
  })

  // ========== Methods ==========
  function closeModal() {
    emit('close')
  }

  function validateForm(): boolean {
    const result = validateCriterionForm(formData.value)
    formErrors.value = result.errors
    return result.valid
  }

  // Real-time Validation Functions
  function validateNameRealtime() {
    const error = validateCriterionName(formData.value.name)
    if (error) {
      formErrors.value.name = error
    } else {
      delete formErrors.value.name
    }
  }

  function validateDescriptionRealtime() {
    const error = validateCriterionDescription(formData.value.description)
    if (error) {
      formErrors.value.description = error
    } else {
      delete formErrors.value.description
    }
  }

  function validateMaxScoreRealtime() {
    const error = validateCriterionMaxScore(formData.value.max_score)
    if (error) {
      formErrors.value.max_score = error
    } else {
      delete formErrors.value.max_score
    }
  }

  function validateWeightRealtime() {
    const error = validateCriterionWeight(formData.value.weight)
    if (error) {
      formErrors.value.weight = error
    } else {
      delete formErrors.value.weight
    }
  }

  async function handleSubmit() {
    if (!validateForm()) {
      return
    }

    formLoading.value = true
    try {
      if (isEditing.value && props.criterion) {
        await EvaluationCriteriaService.update(props.criterion.ID, formData.value)
        await Swal.fire({ icon: 'success', title: 'อัปเดตเกณฑ์สำเร็จ', timer: 1500, showConfirmButton: false })
      } else {
        await EvaluationCriteriaService.create(formData.value)
        await Swal.fire({ icon: 'success', title: 'สร้างเกณฑ์สำเร็จ', timer: 1500, showConfirmButton: false })
      }
      emit('saved')
      closeModal()
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
</script>

<template>
  <Teleport to="body">
    <div 
      v-if="isOpen" 
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      @click.self="closeModal"
      data-theme="light"
    >
      <div class="bg-white rounded-2xl shadow-2xl w-full max-w-lg mx-4 overflow-hidden animate-pop-in">
        <!-- Header -->
        <div class="px-6 py-4 border-b flex items-center justify-between bg-slate-50">
          <div>
            <h2 class="text-xl font-bold text-[#1e3a8a]">
              {{ isEditing ? 'แก้ไขเกณฑ์' : 'เพิ่มเกณฑ์ใหม่' }}
            </h2>
          </div>

          <button class="btn btn-circle btn-ghost btn-sm" @click="closeModal" aria-label="ปิด">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24"
              stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
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
                max="1000"
                step="1"
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
              @click="closeModal"
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
</template>

<style scoped>
.animate-pop-in {
  animation: pop-in 0.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes pop-in {
  0% {
    opacity: 0;
    transform: scale(0.95);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
