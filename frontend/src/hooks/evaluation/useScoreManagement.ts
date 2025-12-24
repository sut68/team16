import { ref, computed } from 'vue'
import type { Ref } from 'vue'
import type { EvaluationResponse, InterviewRoundCriterionResponse } from '@/interfaces/evaluation'
import { EvaluationScoreService, EvaluationService } from '@/services/evaluation/evaluation'

export interface ScoreInput {
  score: number
  comment: string
  scoreRecordId?: number
}

// Hook สำหรับจัดการคะแนนการประเมิน
// - จัดการ input คะแนนแต่ละเกณฑ์
// - คำนวณคะแนนรวม
// - บันทึกคะแนน
export function useScoreManagement(
  evaluation: Ref<EvaluationResponse | null>,
  roundCriteria: Ref<InterviewRoundCriterionResponse[]>
) {
  // Validation Rules (ตาม Backend)
  const VALIDATION_RULES = {
    scoreValue: { min: 0.01 }, // ต้อง > 0 (ประมาณ 0.01 เป็นค่าต่ำสุด)
    comment: { maxLength: 500 },
  }

  // State
  const scoreInputs = ref<Record<number, ScoreInput>>({})
  const scoreErrors = ref<Record<number, { score?: string; comment?: string }>>({})
  const remark = ref('')
  const saving = ref(false)

  // Computed - คะแนนรวมถ่วงน้ำหนัก
  const currentTotalScore = computed(() => {
    let totalScore = 0
    let totalWeight = 0

    for (const criterion of roundCriteria.value) {
      const input = scoreInputs.value[criterion.evaluation_criterion_id]
      if (input && criterion.evaluation_criterion) {
        const maxScore = criterion.evaluation_criterion.max_score || 100
        const normalizedScore = (input.score / maxScore) * 100
        const weightedScore = normalizedScore * criterion.weight
        totalScore += weightedScore
        totalWeight += criterion.weight
      }
    }

    return totalWeight > 0 ? totalScore / totalWeight : 0
  })

  // Computed - ตรวจสอบว่ากรอกคะแนนครบทุกเกณฑ์ และคะแนนต้อง > 0
  const canComplete = computed(() => {
    return roundCriteria.value.every(c => {
      const input = scoreInputs.value[c.evaluation_criterion_id]
      return input && input.score !== undefined && input.score !== null && input.score > 0
    })
  })

  // Computed - จำนวนเกณฑ์ที่กรอกแล้ว
  const filledCriteriaCount = computed(() => {
    return roundCriteria.value.filter(c => {
      const input = scoreInputs.value[c.evaluation_criterion_id]
      return input && input.score !== undefined && input.score !== null && input.score > 0
    }).length
  })

  // Actions
  function initializeScoreInputs() {
    scoreInputs.value = {}

    // โหลดคะแนนที่บันทึกไว้แล้ว
    if (evaluation.value?.evaluation_scores) {
      for (const score of evaluation.value.evaluation_scores) {
        scoreInputs.value[score.evaluation_criterion_id] = {
          score: score.score_value,
          comment: score.comment || '',
          scoreRecordId: score.ID,
        }
      }
    }

    // เพิ่มเกณฑ์ที่ยังไม่มีคะแนน
    for (const criterion of roundCriteria.value) {
      if (!scoreInputs.value[criterion.evaluation_criterion_id]) {
        scoreInputs.value[criterion.evaluation_criterion_id] = {
          score: 0,
          comment: '',
        }
      }
    }

    // โหลด remark
    remark.value = evaluation.value?.remark || ''
  }

  function getScoreInput(criterionId: number): ScoreInput {
    return scoreInputs.value[criterionId] || { score: 0, comment: '' }
  }

  function setScore(criterionId: number, value: number) {
    if (scoreInputs.value[criterionId]) {
      scoreInputs.value[criterionId].score = value

      // Real-time validation สำหรับคะแนน
      if (!scoreErrors.value[criterionId]) {
        scoreErrors.value[criterionId] = {}
      }

      if (value === null || value === undefined || isNaN(value) || value <= 0) {
        scoreErrors.value[criterionId].score = 'คะแนนต้องมากกว่า 0'
      } else {
        delete scoreErrors.value[criterionId].score
        // ลบ object ถ้าไม่มี error เหลือ
        if (Object.keys(scoreErrors.value[criterionId]).length === 0) {
          delete scoreErrors.value[criterionId]
        }
      }
    }
  }

  function setComment(criterionId: number, value: string) {
    if (scoreInputs.value[criterionId]) {
      scoreInputs.value[criterionId].comment = value
      // Validate comment length
      if (value.length > VALIDATION_RULES.comment.maxLength) {
        if (!scoreErrors.value[criterionId]) {
          scoreErrors.value[criterionId] = {}
        }
        scoreErrors.value[criterionId].comment = `ความคิดเห็นต้องไม่เกิน ${VALIDATION_RULES.comment.maxLength} ตัวอักษร`
      } else {
        if (scoreErrors.value[criterionId]) {
          delete scoreErrors.value[criterionId].comment
        }
      }
    }
  }

  // Validate all scores before saving
  function validateScores(): boolean {
    scoreErrors.value = {}
    let isValid = true

    for (const criterion of roundCriteria.value) {
      const input = scoreInputs.value[criterion.evaluation_criterion_id]
      if (!input) continue

      const errors: { score?: string; comment?: string } = {}

      // Score validation: ต้อง > 0
      if (!input.score || input.score <= 0) {
        errors.score = 'คะแนนต้องมากกว่า 0'
        isValid = false
      }

      // Comment validation: max 500
      if (input.comment && input.comment.length > VALIDATION_RULES.comment.maxLength) {
        errors.comment = `ความคิดเห็นต้องไม่เกิน ${VALIDATION_RULES.comment.maxLength} ตัวอักษร`
        isValid = false
      }

      if (Object.keys(errors).length > 0) {
        scoreErrors.value[criterion.evaluation_criterion_id] = errors
      }
    }

    return isValid
  }

  // Get errors for a specific criterion
  function getScoreErrors(criterionId: number): { score?: string; comment?: string } {
    return scoreErrors.value[criterionId] || {}
  }

  async function saveScore(criterionId: number) {
    const input = scoreInputs.value[criterionId]
    if (!input || !evaluation.value) return

    try {
      if (input.scoreRecordId) {
        // Update existing score
        await EvaluationScoreService.update(input.scoreRecordId, {
          score_value: input.score,
          comment: input.comment,
        })
      } else {
        // Create new score
        const result = await EvaluationScoreService.addScore(evaluation.value.ID, {
          evaluation_criterion_id: criterionId,
          score_value: input.score,
          comment: input.comment,
          scoring_admin_id: evaluation.value.admin_id,
        })
        input.scoreRecordId = result.ID
      }
    } catch (err: any) {
      console.error('Failed to save score:', err)
      throw err
    }
  }

  async function saveAllScores(): Promise<void> {
    // Validate before saving
    if (!validateScores()) {
      throw new Error('กรุณาตรวจสอบข้อมูลให้ถูกต้อง')
    }

    saving.value = true
    try {
      // บันทึกคะแนนทุกเกณฑ์
      for (const criterion of roundCriteria.value) {
        await saveScore(criterion.evaluation_criterion_id)
      }

      // บันทึก remark
      if (evaluation.value) {
        await EvaluationService.update(evaluation.value.ID, { remark: remark.value })
      }
    } catch (err: any) {
      console.error('Failed to save all scores:', err)
      throw err
    } finally {
      saving.value = false
    }
  }

  function reset() {
    scoreInputs.value = {}
    scoreErrors.value = {}
    remark.value = ''
  }

  return {
    // State
    scoreInputs,
    scoreErrors,
    remark,
    saving,
    // Computed
    currentTotalScore,
    canComplete,
    filledCriteriaCount,
    // Validation
    VALIDATION_RULES,
    validateScores,
    getScoreErrors,
    // Actions
    initializeScoreInputs,
    getScoreInput,
    setScore,
    setComment,
    saveScore,
    saveAllScores,
    reset,
  }
}
