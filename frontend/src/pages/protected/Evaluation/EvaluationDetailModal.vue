<script setup lang="ts">
  import { ref, watch, computed } from 'vue'
  import type { PropType } from 'vue'
  import type { EvaluationResponse, InterviewRoundCriterionResponse } from '@/interfaces/evaluation'
  import { EvaluationService, InterviewRoundCriteriaService } from '@/services/evaluation/evaluation'
  import { useEvaluationHelpers } from '@/hooks/evaluation'
  import { X, User, Award, CheckCircle, MessageSquare, Calendar, UserCheck } from 'lucide-vue-next'

  // Props & Emits
  const props = defineProps({
    isOpen: { type: Boolean as PropType<boolean>, default: false },
    evaluationId: { type: Number as PropType<number | null>, default: null },
  })
  const emit = defineEmits(['update:isOpen', 'close'])

  // Helpers
  const { getStatusLabel, getStatusColors, getScoreTypeLabel } = useEvaluationHelpers()

  // State
  const evaluation = ref<EvaluationResponse | null>(null)
  const roundCriteria = ref<InterviewRoundCriterionResponse[]>([])
  const loading = ref(false)

  // Computed
  const studentInfo = computed(() => {
    if (!evaluation.value?.application_scholarship?.application?.student_profile) return null
    const profile = evaluation.value.application_scholarship.application.student_profile
    return {
      name: `${profile.first_name_th || ''} ${profile.last_name_th || ''}`.trim() || '-',
      studentId: profile.student_id || '-',
      major: profile.major?.major_name || '-',
      gpa: profile.gpax ?? '-',
    }
  })

  const scholarshipInfo = computed(() => {
    if (!evaluation.value?.application_scholarship?.scholarship) return null
    return {
      name: evaluation.value.application_scholarship.scholarship.scholarship_name || '-',
    }
  })

  const adminInfo = computed(() => {
    if (!evaluation.value?.admin_profile) return null
    return {
      name: `${evaluation.value.admin_profile.first_name || ''} ${evaluation.value.admin_profile.last_name || ''}`.trim() || 'แอดมิน',
    }
  })

  // Get score for a criterion
  function getScoreForCriterion(criterionId: number) {
    const score = evaluation.value?.evaluation_scores?.find(s => s.evaluation_criterion_id === criterionId)
    return score || null
  }

  function getScoreDisplay(score: number | undefined, scoreType: string, maxScore: number): string {
    if (score === undefined || score === null) return '-'
    if (scoreType === 'pass_fail') {
      return score >= 1 ? 'ผ่าน' : 'ไม่ผ่าน'
    }
    return `${score}/${maxScore}`
  }

  function getScoreClass(score: number | undefined, scoreType: string): string {
    if (score === undefined || score === null) return 'text-gray-400'
    if (scoreType === 'pass_fail') {
      return score >= 1 ? 'text-green-600' : 'text-red-600'
    }
    return 'text-gray-900'
  }

  // Fetch data
  async function fetchData() {
    if (!props.evaluationId) return

    loading.value = true
    try {
      evaluation.value = await EvaluationService.getById(props.evaluationId)
      
      if (evaluation.value.interview_round_id) {
        roundCriteria.value = await InterviewRoundCriteriaService.getByRoundId(evaluation.value.interview_round_id)
      }
    } catch (err) {
      console.error('Failed to load evaluation:', err)
      evaluation.value = null
    } finally {
      loading.value = false
    }
  }

  function close() {
    emit('update:isOpen', false)
    emit('close')
  }

  function formatDate(dateStr: string | undefined): string {
    if (!dateStr) return '-'
    return new Date(dateStr).toLocaleDateString('th-TH', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  // Watch isOpen
  watch(
    () => props.isOpen,
    (open) => {
      if (open && props.evaluationId) {
        fetchData()
      }
    }
  )
</script>

<template>
  <teleport to="body">
    <div
      v-if="isOpen"
      class="fixed inset-0 z-[300] flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
      @click.self="close"
    >
      <div 
        class="bg-white w-full max-w-2xl max-h-[85vh] rounded-2xl shadow-2xl flex flex-col overflow-hidden animate-pop-in"
      >
        <!-- Header -->
        <div class="px-6 py-4 border-b flex items-center justify-between">
          <div>
            <h2 class="text-xl font-bold text-[#1e3a8a]">รายละเอียดการประเมิน</h2>
            <p v-if="studentInfo" class="text-sm text-gray-500">{{ studentInfo.name }}</p>
          </div>
          <button 
            @click="close" 
            class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="flex-1 flex items-center justify-center p-8">
          <div class="text-center">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-indigo-500 border-r-transparent"></div>
            <p class="mt-2 text-gray-500">กำลังโหลดข้อมูล...</p>
          </div>
        </div>

        <!-- Content -->
        <div v-else-if="evaluation" class="flex-1 overflow-y-auto p-6 space-y-6">
          <!-- Info Cards -->
          <div class="grid grid-cols-2 gap-4">
            <!-- Student Card -->
            <div class="bg-gray-50 rounded-xl p-4">
              <div class="flex items-center gap-2 mb-3">
                <div class="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center">
                  <User class="w-4 h-4 text-blue-600" />
                </div>
                <span class="text-sm font-medium text-gray-700">ข้อมูลนักศึกษา</span>
              </div>
              <div v-if="studentInfo" class="space-y-1 text-sm">
                <p><span class="text-gray-500">ชื่อ:</span> <span class="font-medium">{{ studentInfo.name }}</span></p>
                <p><span class="text-gray-500">รหัส:</span> <span class="font-medium">{{ studentInfo.studentId }}</span></p>
                <p><span class="text-gray-500">สาขา:</span> <span class="font-medium">{{ studentInfo.major }}</span></p>
              </div>
            </div>

            <!-- Scholarship Card -->
            <div class="bg-gray-50 rounded-xl p-4">
              <div class="flex items-center gap-2 mb-3">
                <div class="w-8 h-8 bg-purple-100 rounded-lg flex items-center justify-center">
                  <Award class="w-4 h-4 text-purple-600" />
                </div>
                <span class="text-sm font-medium text-gray-700">ทุนการศึกษา</span>
              </div>
              <div v-if="scholarshipInfo" class="text-sm">
                <p class="font-medium">{{ scholarshipInfo.name }}</p>
              </div>
            </div>
          </div>

          <!-- Status & Score Summary -->
          <div class="bg-gradient-to-r from-indigo-50 to-blue-50 rounded-xl p-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div>
                  <p class="text-sm text-gray-500">สถานะ</p>
                  <span 
                    class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium mt-1"
                    :class="[getStatusColors(evaluation.status).bg, getStatusColors(evaluation.status).text]"
                  >
                    {{ getStatusLabel(evaluation.status) }}
                  </span>
                </div>
              </div>
              <div class="text-right">
                <p class="text-sm text-gray-500">คะแนนรวม</p>
                <p class="text-3xl font-bold text-indigo-600">{{ evaluation.total_score?.toFixed(1) || '-' }}</p>
              </div>
            </div>
          </div>

          <!-- Scores Table -->
          <div>
            <h3 class="text-sm font-semibold text-gray-700 mb-3 flex items-center gap-2">
              <CheckCircle class="w-4 h-4 text-gray-400" />
              คะแนนรายเกณฑ์
            </h3>
            <div class="bg-white border border-gray-200 rounded-xl overflow-hidden">
              <table class="w-full">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase">เกณฑ์</th>
                    <th class="px-4 py-3 text-center text-xs font-semibold text-gray-500 uppercase">ประเภท</th>
                    <th class="px-4 py-3 text-center text-xs font-semibold text-gray-500 uppercase">คะแนน</th>
                    <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase">ความคิดเห็น</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100">
                  <tr v-for="criterion in roundCriteria" :key="criterion.ID" class="hover:bg-gray-50">
                    <td class="px-4 py-3">
                      <p class="font-medium text-gray-900">{{ criterion.evaluation_criterion?.name }}</p>
                      <p class="text-xs text-gray-500">น้ำหนัก: {{ criterion.weight }}</p>
                    </td>
                    <td class="px-4 py-3 text-center">
                      <span class="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded">
                        {{ getScoreTypeLabel(criterion.evaluation_criterion?.score_type || '') }}
                      </span>
                    </td>
                    <td class="px-4 py-3 text-center">
                      <span 
                        class="font-semibold"
                        :class="getScoreClass(
                          getScoreForCriterion(criterion.evaluation_criterion_id)?.score_value,
                          criterion.evaluation_criterion?.score_type || 'numeric'
                        )"
                      >
                        {{ getScoreDisplay(
                          getScoreForCriterion(criterion.evaluation_criterion_id)?.score_value,
                          criterion.evaluation_criterion?.score_type || 'numeric',
                          criterion.evaluation_criterion?.max_score || 100
                        ) }}
                      </span>
                    </td>
                    <td class="px-4 py-3 text-sm text-gray-600">
                      {{ getScoreForCriterion(criterion.evaluation_criterion_id)?.comment || '-' }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Remark -->
          <div v-if="evaluation.remark">
            <h3 class="text-sm font-semibold text-gray-700 mb-2 flex items-center gap-2">
              <MessageSquare class="w-4 h-4 text-gray-400" />
              หมายเหตุ
            </h3>
            <p class="text-sm text-gray-600 bg-yellow-50 border border-yellow-100 rounded-lg p-3">
              {{ evaluation.remark }}
            </p>
          </div>

          <!-- Meta Info -->
          <div class="flex items-center justify-between text-xs text-gray-400 pt-4 border-t">
            <div class="flex items-center gap-1">
              <UserCheck class="w-3 h-3" />
              <span>ประเมินโดย: {{ adminInfo?.name || '-' }}</span>
            </div>
            <div class="flex items-center gap-1">
              <Calendar class="w-3 h-3" />
              <span>{{ formatDate(evaluation.UpdatedAt) }}</span>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-6 py-4 border-t bg-gray-50 flex justify-end">
          <button 
            @click="close"
            class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg font-medium transition-colors"
          >
            ปิด
          </button>
        </div>
      </div>
    </div>
  </teleport>
</template>

<style scoped>
@keyframes pop-in {
  from { opacity: 0; transform: translateY(10px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
.animate-pop-in {
  animation: pop-in 0.2s ease-out;
}
</style>
