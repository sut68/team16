<script setup lang="ts">
import { ref, watch, toRef } from 'vue'
import type { PropType } from 'vue'
import Swal from 'sweetalert2'

// Hooks
import {
  useEvaluationData,
  useScoreManagement,
  useEvaluationActions,
  useEvaluationHelpers
} from '@/hooks/evaluation'

// Icons
import {
  Save, CheckCircle, XCircle, Clock,
  Award, User, BookOpen, MessageSquare, X
} from 'lucide-vue-next'

// ========== Props / Emits ==========
const props = defineProps({
  isOpen: { type: Boolean as PropType<boolean>, default: false },
  evaluationId: { type: Number as PropType<number | null>, default: null },
})
const emit = defineEmits(['update:isOpen', 'close', 'completed'])

// ========== Convert props to ref for hooks ==========
const evaluationIdRef = toRef(props, 'evaluationId')

// ========== Use Hooks ==========
const {
  evaluation,
  roundCriteria,
  loading,
  fetchData,
  studentInfo,
  scholarshipInfo,
  reset: resetData
} = useEvaluationData(evaluationIdRef)

const {
  remark,
  saving,
  currentTotalScore,
  canComplete,
  initializeScoreInputs,
  getScoreInput,
  setScore,
  setComment,
  saveAllScores,
  getScoreErrors,
  VALIDATION_RULES,
  reset: resetScores
} = useScoreManagement(evaluation, roundCriteria)

const {
  completing,
  completeEvaluation: doCompleteEvaluation
} = useEvaluationActions(evaluation, currentTotalScore, canComplete)

const { getScoreTypeLabel } = useEvaluationHelpers()

// ========== Watch isOpen ==========
watch(
  () => props.isOpen,
  async (open) => {
    if (open && props.evaluationId) {
      try {
        await fetchData()
        initializeScoreInputs()
      } catch {
        Swal.fire({ icon: 'error', title: 'ไม่สามารถโหลดข้อมูลได้' })
        close()
      }
    }
  }
)

// ========== Methods =========
function close() {
  resetData()
  resetScores()
  emit('update:isOpen', false)
  emit('close')
}

async function handleSaveAllScores() {
  try {
    await saveAllScores()
    await Swal.fire({
      icon: 'success',
      title: 'บันทึกสำเร็จ',
      timer: 1500,
      showConfirmButton: false,
    })
    // Refresh data
    await fetchData()
    initializeScoreInputs()
  } catch (err: any) {
    Swal.fire({
      icon: 'error',
      title: 'บันทึกไม่สำเร็จ',
      text: err?.response?.data?.error || 'กรุณาลองใหม่',
    })
  }
}

async function handleCompleteEvaluation(decision: 'approved' | 'rejected' | 'waitlist') {
  const success = await doCompleteEvaluation(decision, () => {
    emit('completed')
    close()
  })
  if (!success) return
}

// Combined saving state
const isBusy = ref(false)
watch([saving, completing], ([s, c]) => {
  isBusy.value = s || c
})
</script>

<template>
  <Teleport to="body">
    <div
      v-if="props.isOpen"
      class="fixed inset-0 z-[200] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
      @click.self="close"
      data-theme="light"
    >
      <div
        role="dialog"
        aria-modal="true"
        class="bg-white w-full max-w-6xl h-[90vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in"
      >
        <!-- Header -->
        <div class="px-6 py-4 border-b flex items-center justify-between bg-gradient-to-r from-indigo-600 to-blue-600">
          <div>
            <h2 class="text-xl font-bold text-white">ประเมินผู้สมัคร</h2>
            <p class="text-sm text-white/80" v-if="evaluation">
              รหัสการประเมิน: #{{ evaluation.ID }}
            </p>
          </div>

          <button
            class="p-2 hover:bg-white/20 rounded-full transition-colors"
            @click="close"
            aria-label="ปิด"
          >
            <X class="w-6 h-6 text-white" />
          </button>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="flex-1 flex items-center justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>

        <!-- Body -->
        <div v-else class="flex-1 overflow-y-auto p-6 bg-gray-50">
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

            <!-- Left Column - Info Cards -->
            <div class="lg:col-span-1 space-y-4">
              <!-- Student Info Card -->
              <div class="bg-white rounded-xl shadow-sm p-5 border border-gray-100">
                <div class="flex items-center gap-3 mb-4">
                  <div class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                    <User class="w-5 h-5 text-blue-600" />
                  </div>
                  <h3 class="font-semibold text-gray-800">ข้อมูลนักศึกษา</h3>
                </div>
                <div v-if="studentInfo" class="space-y-2 text-sm">
                  <div>
                    <span class="text-gray-500">ชื่อ-นามสกุล:</span>
                    <p class="font-semibold text-gray-800">{{ studentInfo.name }}</p>
                  </div>
                  <div>
                    <span class="text-gray-500">รหัสนักศึกษา:</span>
                    <p class="font-semibold text-gray-800">{{ studentInfo.studentId }}</p>
                  </div>
                  <div>
                    <span class="text-gray-500">สาขา:</span>
                    <p class="font-semibold text-gray-800">{{ studentInfo.major }}</p>
                  </div>
                  <div>
                    <span class="text-gray-500">เกรดเฉลี่ย:</span>
                    <p class="font-semibold text-gray-800">{{ studentInfo.gpa }}</p>
                  </div>
                </div>
              </div>

              <!-- Scholarship Info Card -->
              <div class="bg-white rounded-xl shadow-sm p-5 border border-gray-100">
                <div class="flex items-center gap-3 mb-4">
                  <div class="w-10 h-10 bg-purple-100 rounded-full flex items-center justify-center">
                    <BookOpen class="w-5 h-5 text-purple-600" />
                  </div>
                  <h3 class="font-semibold text-gray-800">ทุนการศึกษา</h3>
                </div>
                <div v-if="scholarshipInfo" class="text-sm">
                  <p class="font-semibold text-gray-800">{{ scholarshipInfo.name }}</p>
                </div>
              </div>

              <!-- Score Summary Card -->
              <div class="bg-gradient-to-br from-indigo-600 to-blue-700 rounded-xl shadow-lg p-5 text-white">
                <div class="flex items-center gap-3 mb-4">
                  <div class="w-10 h-10 bg-white/20 rounded-full flex items-center justify-center">
                    <Award class="w-5 h-5" />
                  </div>
                  <h3 class="font-semibold">คะแนนรวม</h3>
                </div>
                <div class="text-4xl font-bold text-center my-4">
                  {{ currentTotalScore.toFixed(1) }}
                </div>
                <div class="text-center text-sm opacity-80">จากคะแนนเต็ม 100</div>
              </div>
            </div>

            <!-- Right Column - Scoring Form -->
            <div class="lg:col-span-2 space-y-4">
              <!-- Criteria Scoring -->
              <div class="bg-white rounded-xl shadow-sm p-5 border border-gray-100">
                <h3 class="font-semibold text-gray-800 mb-4">เกณฑ์การประเมิน</h3>

                <div class="space-y-3">
                  <div
                    v-for="criterion in roundCriteria"
                    :key="criterion.ID"
                    class="border border-gray-200 rounded-lg p-4 hover:border-blue-300 transition-colors"
                  >
                    <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
                      <div class="flex-1">
                        <div class="flex items-center gap-2 flex-wrap">
                          <h4 class="font-medium text-gray-800">
                            {{ criterion.evaluation_criterion?.name }}
                          </h4>
                          <span class="text-xs px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full">
                            น้ำหนัก: {{ criterion.weight }}
                          </span>
                        </div>
                        <p class="text-sm text-gray-500 mt-1">
                          {{ criterion.evaluation_criterion?.description || '-' }}
                        </p>
                        <div class="text-xs text-gray-400 mt-1">
                          ประเภท: {{ getScoreTypeLabel(criterion.evaluation_criterion?.score_type || '') }} |
                          คะแนนเต็ม: {{ criterion.evaluation_criterion?.max_score }}
                        </div>
                      </div>

                      <div class="flex items-center gap-3">
                        <div class="flex flex-col items-end">
                          <label class="text-xs text-gray-500 mb-1">คะแนน (ต้อง > 0)</label>
                          <input
                            :value="getScoreInput(criterion.evaluation_criterion_id).score"
                            @input="setScore(criterion.evaluation_criterion_id, Number(($event.target as HTMLInputElement).value))"
                            type="number"
                            :min="0.01"
                            :max="criterion.evaluation_criterion?.max_score || 100"
                            step="0.1"
                            class="w-24 px-3 py-2 border rounded-lg text-center font-semibold focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                            :class="getScoreErrors(criterion.evaluation_criterion_id).score ? 'border-red-500' : 'border-gray-300'"
                          />
                          <p v-if="getScoreErrors(criterion.evaluation_criterion_id).score" class="mt-1 text-xs text-red-500">
                            {{ getScoreErrors(criterion.evaluation_criterion_id).score }}
                          </p>
                        </div>
                      </div>
                    </div>

                    <!-- Comment -->
                    <div class="mt-3">
                      <div class="flex items-center justify-between mb-1">
                        <span class="text-xs text-gray-500">ความคิดเห็น (ไม่บังคับ)</span>
                        <span 
                          class="text-xs" 
                          :class="(getScoreInput(criterion.evaluation_criterion_id).comment?.length || 0) > VALIDATION_RULES.comment.maxLength ? 'text-red-500' : 'text-gray-400'"
                        >
                          {{ getScoreInput(criterion.evaluation_criterion_id).comment?.length || 0 }}/{{ VALIDATION_RULES.comment.maxLength }}
                        </span>
                      </div>
                      <textarea
                        :value="getScoreInput(criterion.evaluation_criterion_id).comment"
                        @input="setComment(criterion.evaluation_criterion_id, ($event.target as HTMLTextAreaElement).value)"
                        rows="2"
                        :maxlength="VALIDATION_RULES.comment.maxLength"
                        class="w-full px-3 py-2 border rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        :class="getScoreErrors(criterion.evaluation_criterion_id).comment ? 'border-red-500' : 'border-gray-200'"
                        placeholder="ความคิดเห็น (ถ้ามี)..."
                      ></textarea>
                      <p v-if="getScoreErrors(criterion.evaluation_criterion_id).comment" class="mt-1 text-xs text-red-500">
                        {{ getScoreErrors(criterion.evaluation_criterion_id).comment }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Remark -->
              <div class="bg-white rounded-xl shadow-sm p-5 border border-gray-100">
                <div class="flex items-center gap-3 mb-3">
                  <MessageSquare class="w-5 h-5 text-gray-600" />
                  <h3 class="font-semibold text-gray-800">หมายเหตุ</h3>
                </div>
                <textarea
                  v-model="remark"
                  rows="3"
                  class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="หมายเหตุเพิ่มเติม..."
                ></textarea>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-6 py-4 border-t bg-slate-50 flex flex-col sm:flex-row items-center justify-between gap-4">
          <button
            @click="handleSaveAllScores"
            :disabled="isBusy"
            class="w-full sm:w-auto px-6 py-2.5 bg-gray-600 hover:bg-gray-700 text-white rounded-full font-medium transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
          >
            <Save class="w-4 h-4" />
            <span>บันทึกคะแนน</span>
          </button>

          <div class="flex flex-col sm:flex-row items-center gap-2 w-full sm:w-auto">
            <button
              @click="handleCompleteEvaluation('rejected')"
              :disabled="isBusy || !canComplete"
              class="w-full sm:w-auto px-5 py-2.5 bg-red-500 hover:bg-red-600 text-white rounded-full font-medium transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
            >
              <XCircle class="w-4 h-4" />
              <span>ไม่อนุมัติ</span>
            </button>

            <button
              @click="handleCompleteEvaluation('waitlist')"
              :disabled="isBusy || !canComplete"
              class="w-full sm:w-auto px-5 py-2.5 bg-amber-500 hover:bg-amber-600 text-white rounded-full font-medium transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
            >
              <Clock class="w-4 h-4" />
              <span>รอพิจารณา</span>
            </button>

            <button
              @click="handleCompleteEvaluation('approved')"
              :disabled="isBusy || !canComplete"
              class="w-full sm:w-auto px-5 py-2.5 bg-emerald-500 hover:bg-emerald-600 text-white rounded-full font-medium transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
            >
              <CheckCircle class="w-4 h-4" />
              <span>อนุมัติ</span>
            </button>
          </div>
        </div>

        <!-- Warning if not complete -->
        <div v-if="!canComplete && !loading" class="px-6 pb-3 bg-slate-50 text-center">
          <p class="text-sm text-amber-600">
            กรุณากรอกคะแนนให้ครบทุกเกณฑ์ก่อนสรุปผล
          </p>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
@keyframes pop-in {
  from { opacity: 0; transform: translateY(6px) scale(.995); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.animate-pop-in {
  animation: pop-in .12s cubic-bezier(.2,.9,.3,1);
}

input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
  -webkit-appearance: none;
  appearance: none;
  margin: 0;
}

input[type="number"] {
  -moz-appearance: textfield;
  appearance: textfield;
}
</style>
