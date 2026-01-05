import { ref, computed } from 'vue'
import type { Ref } from 'vue'
import type { EvaluationResponse, InterviewRoundCriterionResponse } from '@/interfaces/evaluation'
import {
  EvaluationService,
  InterviewRoundCriteriaService
} from '@/services/evaluation/evaluation'

// Hook สำหรับจัดการข้อมูลการประเมิน
// โหลดข้อมูล Evaluation และ Criteria
// แปลงข้อมูลเป็น computed ที่ใช้งานง่าย
export function useEvaluationData(evaluationId: Ref<number | null>) {
  // State
  const evaluation = ref<EvaluationResponse | null>(null)
  const roundCriteria = ref<InterviewRoundCriterionResponse[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Computed - Student Info
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

  // Computed - Scholarship Info
  const scholarshipInfo = computed(() => {
    if (!evaluation.value?.application_scholarship?.scholarship) return null
    const scholarship = evaluation.value.application_scholarship.scholarship
    return {
      name: scholarship.scholarship_name || '-',
      description: scholarship.description || '',
    }
  })

  // Computed - Interview Round Info
  const interviewRoundInfo = computed(() => {
    if (!evaluation.value?.interview_round) return null
    return evaluation.value.interview_round
  })

  // Actions
  async function fetchData() {
    if (!evaluationId.value) return

    loading.value = true
    error.value = null

    try {
      evaluation.value = await EvaluationService.getById(evaluationId.value)

      if (evaluation.value.interview_round_id) {
        roundCriteria.value = await InterviewRoundCriteriaService.getByRoundId(
          evaluation.value.interview_round_id
        )
      }
    } catch (err: any) {
      console.error('Failed to fetch evaluation data:', err)
      error.value = err?.response?.data?.error || 'ไม่สามารถโหลดข้อมูลได้'
      throw err
    } finally {
      loading.value = false
    }
  }

  function reset() {
    evaluation.value = null
    roundCriteria.value = []
    error.value = null
  }

  return {
    // State
    evaluation,
    roundCriteria,
    loading,
    error,
    // Computed
    studentInfo,
    scholarshipInfo,
    interviewRoundInfo,
    // Actions
    fetchData,
    reset,
  }
}
