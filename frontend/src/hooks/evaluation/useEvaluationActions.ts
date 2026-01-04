import { ref } from 'vue'
import type { Ref, ComputedRef } from 'vue'
import type { EvaluationResponse } from '@/interfaces/evaluation'
import { EvaluationService } from '@/services/evaluation/evaluation'
import Swal from 'sweetalert2'

export type EvaluationDecision = 'approved' | 'rejected' | 'waitlist'

const DECISION_LABELS: Record<EvaluationDecision, string> = {
  approved: 'อนุมัติ',
  rejected: 'ไม่อนุมัติ',
  waitlist: 'รอพิจารณา',
}

const DECISION_COLORS: Record<EvaluationDecision, string> = {
  approved: '#10b981',
  rejected: '#ef4444',
  waitlist: '#f59e0b',
}

// Hook สำหรับจัดการ Actions การประเมิน
// - การ Complete Evaluation พร้อม confirm dialog
// - การตัดสินใจ (approved/rejected/waitlist)
export function useEvaluationActions(
  evaluation: Ref<EvaluationResponse | null>,
  currentTotalScore: ComputedRef<number>,
  canComplete: ComputedRef<boolean>
) {
  const completing = ref(false)

  async function completeEvaluation(
    decision: EvaluationDecision,
    onSuccess?: () => void
  ): Promise<boolean> {
    if (!canComplete.value) {
      await Swal.fire({
        icon: 'warning',
        title: 'กรุณากรอกคะแนนให้ครบ',
        text: 'ต้องกรอกคะแนนทุกเกณฑ์ก่อนจึงจะสามารถสรุปผลได้',
      })
      return false
    }

    const result = await Swal.fire({
      icon: 'question',
      title: 'ยืนยันการสรุปผล',
      html: `
        <p>คะแนนรวม: <strong>${currentTotalScore.value.toFixed(2)}</strong> คะแนน</p>
        <p>ผลการพิจารณา: <strong>${DECISION_LABELS[decision]}</strong></p>
      `,
      showCancelButton: true,
      confirmButtonColor: DECISION_COLORS[decision],
      confirmButtonText: 'ยืนยัน',
      cancelButtonText: 'ยกเลิก',
    })

    if (!result.isConfirmed || !evaluation.value) {
      return false
    }

    completing.value = true
    try {
      await EvaluationService.complete(evaluation.value.ID, { final_decision: decision })

      await Swal.fire({
        icon: 'success',
        title: 'บันทึกผลสำเร็จ',
        timer: 1500,
        showConfirmButton: false,
      })

      onSuccess?.()
      return true
    } catch (err: any) {
      console.error('Failed to complete evaluation:', err)
      await Swal.fire({
        icon: 'error',
        title: 'เกิดข้อผิดพลาด',
        text: err?.response?.data?.error || 'กรุณาลองใหม่',
      })
      return false
    } finally {
      completing.value = false
    }
  }

  function getDecisionLabel(decision: EvaluationDecision): string {
    return DECISION_LABELS[decision]
  }

  function getDecisionColor(decision: EvaluationDecision): string {
    return DECISION_COLORS[decision]
  }

  return {
    // State
    completing,
    // Actions
    completeEvaluation,
    // Helpers
    getDecisionLabel,
    getDecisionColor,
  }
}
