// Utility functions สำหรับระบบประเมิน

export type ScoreType = 'numeric' | 'grade' | 'pass_fail'
export type EvaluationStatus = 'pending' | 'in_progress' | 'completed' | 'approved' | 'rejected'

const SCORE_TYPE_LABELS: Record<ScoreType, string> = {
  numeric: 'คะแนนตัวเลข',
  grade: 'เกรด',
  pass_fail: 'ผ่าน/ไม่ผ่าน',
}

const STATUS_LABELS: Record<EvaluationStatus, string> = {
  pending: 'รอประเมิน',
  in_progress: 'กำลังประเมิน',
  completed: 'ประเมินเสร็จ',
  approved: 'อนุมัติ',
  rejected: 'ไม่อนุมัติ',
}

const STATUS_COLORS: Record<EvaluationStatus, { bg: string; text: string }> = {
  pending: { bg: 'bg-gray-100', text: 'text-gray-600' },
  in_progress: { bg: 'bg-blue-100', text: 'text-blue-600' },
  completed: { bg: 'bg-purple-100', text: 'text-purple-600' },
  approved: { bg: 'bg-green-100', text: 'text-green-600' },
  rejected: { bg: 'bg-red-100', text: 'text-red-600' },
}

// Hook สำหรับ Helper functions ในระบบประเมิน
export function useEvaluationHelpers() {
  function getScoreTypeLabel(type: string): string {
    return SCORE_TYPE_LABELS[type as ScoreType] || type
  }

  function getStatusLabel(status: string): string {
    return STATUS_LABELS[status as EvaluationStatus] || status
  }

  function getStatusColors(status: string): { bg: string; text: string } {
    return STATUS_COLORS[status as EvaluationStatus] || { bg: 'bg-gray-100', text: 'text-gray-600' }
  }

  function formatScore(score: number, maxScore: number = 100): string {
    return `${score.toFixed(1)}/${maxScore}`
  }

  function calculatePercentage(score: number, maxScore: number = 100): number {
    if (maxScore === 0) return 0
    return (score / maxScore) * 100
  }

  function getScoreColorClass(percentage: number): string {
    if (percentage >= 80) return 'text-green-600'
    if (percentage >= 60) return 'text-yellow-600'
    if (percentage >= 40) return 'text-orange-600'
    return 'text-red-600'
  }

  return {
    getScoreTypeLabel,
    getStatusLabel,
    getStatusColors,
    formatScore,
    calculatePercentage,
    getScoreColorClass,
  }
}
